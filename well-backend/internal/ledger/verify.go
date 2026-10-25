package ledger

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/ledger/hash"

	"github.com/jackc/pgx/v5"
)

// RowVerifyReport is the per-row result of VerifySingle.
type RowVerifyReport struct {
	Table        string   `json:"table"`
	Seq          int64    `json:"seq"`
	Found        bool     `json:"found"`
	Batched      bool     `json:"batched"`
	RowHashOK    bool     `json:"rowHashOK"`
	PrevLinkOK   bool     `json:"prevLinkOK"`
	BatchRootOK  bool     `json:"batchRootOK"`
	ChainHashOK  bool     `json:"chainHashOK"`
	Anchored     bool     `json:"anchored"`
	AnchorOK     bool     `json:"anchorOK"`
	AnchorTxHash string   `json:"anchorTxHash,omitempty"`
	Pass         bool     `json:"pass"`
	Issues       []string `json:"issues,omitempty"`
}

// LedgerReport is the aggregate result of VerifyLedger.
type LedgerReport struct {
	TablesChecked       int      `json:"tablesChecked"`
	RowsChecked         int      `json:"rowsChecked"`
	BatchesChecked      int      `json:"batchesChecked"`
	AnchorsChecked      int      `json:"anchorsChecked"`
	UnexpectedFailures  int      `json:"unexpectedFailures"`
	Issues              []string `json:"issues,omitempty"`
	Pass                bool     `json:"pass"`
}

// VerifySingle recomputes the hashes for one row and cross-checks them against
// both the stored ledger columns and (when anchored) the on-chain anchor.
func (l *Ledger) VerifySingle(ctx context.Context, table string, seq int64) (*RowVerifyReport, error) {
	rep := &RowVerifyReport{Table: table, Seq: seq, Pass: true}

	var rowHash, prevHash []byte
	var content *big.Int
	var batchID *string

	switch table {
	case db.TableAccessLogs:
		r, err := l.database.GetAccessLogRow(ctx, seq)
		if err != nil {
			if isNoRows(err) {
				return rep, nil // not found
			}
			return nil, err
		}
		rep.Found = true
		rep.Batched = r.BatchID != nil
		batchID = r.BatchID
		rowHash = r.RowHash
		prevHash = r.PrevHash
		content = hash.AccessLogContentHash(
			r.AccessedAt, r.AccessorAddress, r.PatientAddress, r.RecordID,
			r.AccessType, r.PolicyMatched, r.WasEmergency, r.TxHash)
	case db.TableEventLogs:
		r, err := l.database.GetEventLogRow(ctx, seq)
		if err != nil {
			if isNoRows(err) {
				return rep, nil
			}
			return nil, err
		}
		rep.Found = true
		rep.Batched = r.BatchID != nil
		batchID = r.BatchID
		rowHash = r.RowHash
		prevHash = r.PrevHash
		content = hash.EventLogContentHash(
			r.EventType, r.ContractAddress, r.TxHash, r.BlockNumber, r.Data, r.LoggedAt)
	default:
		return nil, fmt.Errorf("unknown log table %q", table)
	}

	if !rep.Found {
		rep.Pass = false
		rep.Issues = append(rep.Issues, fmt.Sprintf("row %d not found in %s", seq, table))
		return rep, nil
	}

	// 1. Row hash: stored row_hash must equal Poseidon(stored prev_hash, content).
	recomputed := hash.RowHash(prevHash, content)
	rep.RowHashOK = bytes.Equal(hash.FieldBytes(recomputed), rowHash)
	if !rep.RowHashOK {
		rep.Pass = false
		rep.Issues = append(rep.Issues, fmt.Sprintf("row_hash mismatch: stored 0x%s, recomputed 0x%s",
			hex.EncodeToString(rowHash), hex.EncodeToString(hash.FieldBytes(recomputed))))
	}

	// 2. Prev link: stored prev_hash must equal the previous chain row's hash.
	prevActual, err := l.database.PrevRowHash(ctx, table, seq)
	if err != nil {
		return nil, err
	}
	expectedPrev := hash.FieldBytes(hash.Zero)
	if len(prevActual) > 0 {
		expectedPrev = prevActual
	}
	rep.PrevLinkOK = bytes.Equal(expectedPrev, zeroIfNil(prevHash))
	if !rep.PrevLinkOK {
		rep.Pass = false
		rep.Issues = append(rep.Issues, fmt.Sprintf("prev_hash breaks the chain: stored 0x%s, expected 0x%s",
			hex.EncodeToString(zeroIfNil(prevHash)), hex.EncodeToString(expectedPrev)))
	}

	// 3. Batch + chain + anchor.
	if !rep.Batched {
		rep.Issues = append(rep.Issues, "row not yet sealed into a batch")
		return rep, nil
	}

	if err := l.checkBatch(ctx, table, *batchID, rep, seq, rowHash); err != nil {
		return nil, err
	}
	return rep, nil
}

// checkBatch verifies a row's enclosing batch: Merkle root, chain hash and
// (if anchored) the value recorded on LogAnchor.
func (l *Ledger) checkBatch(ctx context.Context, table, batchID string, rep *RowVerifyReport, seq int64, rowHash []byte) error {
	batch, err := l.database.BatchByID(ctx, batchID)
	if err != nil {
		return err
	}

	leafHashes, err := l.database.BatchRowsForTree(ctx, table, batch.StartSeq, batch.EndSeq)
	if err != nil {
		return err
	}
	leaves := make([]*big.Int, len(leafHashes))
	for i, h := range leafHashes {
		leaves[i] = hash.BytesToField32(h)
	}
	root := hash.MerkleRoot(leaves)

	rootHex := hex.EncodeToString(hash.FieldBytes(root))
	rep.BatchRootOK = rootHex == batch.MerkleRoot
	if !rep.BatchRootOK {
		rep.Issues = append(rep.Issues, fmt.Sprintf("batch %s merkle root mismatch: stored %s, recomputed %s",
			batch.BatchID, batch.MerkleRoot, rootHex))
	}

	prevF := hash.Zero
	if len(batch.PrevChainHash) > 0 {
		prevF = hash.BytesToField(batch.PrevChainHash)
	}
	chain := hash.Poseidon2(prevF, root)
	rep.ChainHashOK = bytes.Equal(hash.FieldBytes(chain), batch.ChainHash)
	if !rep.ChainHashOK {
		rep.Issues = append(rep.Issues, fmt.Sprintf("batch %s chain hash mismatch: stored 0x%s, recomputed 0x%s",
			batch.BatchID, hex.EncodeToString(batch.ChainHash), hex.EncodeToString(hash.FieldBytes(chain))))
	}

	// Merkle membership of this row within its batch (row_seq -> leaf index).
	if idx := seq - batch.StartSeq; idx >= 0 && idx < int64(batch.LeafCount) {
		if sibs, ok := hash.MerkleProof(leaves, int(idx)); ok {
			if !hash.VerifyMerkle(hash.BytesToField32(rowHash), int(idx), sibs, root) {
				rep.Issues = append(rep.Issues, fmt.Sprintf("row %d is not a member of batch %s", seq, batch.BatchID))
			}
		}
	}

	rep.Anchored = batch.Anchored
	if batch.AnchorTxHash != "" {
		rep.AnchorTxHash = batch.AnchorTxHash
	}
	if batch.Anchored {
		if l.anchor == nil {
			rep.AnchorOK = false
			rep.Issues = append(rep.Issues, "batch anchored but on-chain anchor not available (no LOG_ANCHOR_ADDRESS)")
		} else {
			onChain, err := l.anchor.Anchors(nil, big.NewInt(batch.AnchorBatchID))
			if err != nil {
				return fmt.Errorf("read anchor %d: %w", batch.AnchorBatchID, err)
			}
			rep.AnchorOK = onChain.ChainHash.Cmp(hash.BytesToField(batch.ChainHash)) == 0
			if !rep.AnchorOK {
				rep.Issues = append(rep.Issues, fmt.Sprintf(
					"on-chain anchor %d mismatch: chain 0x%s, anchored 0x%s",
					batch.AnchorBatchID,
					hex.EncodeToString(batch.ChainHash),
					hex.EncodeToString(hash.FieldBytes(onChain.ChainHash))))
			}
		}
	}

	rep.Pass = rep.RowHashOK && rep.PrevLinkOK && rep.BatchRootOK && rep.ChainHashOK &&
		(!rep.Anchored || rep.AnchorOK)
	return nil
}

// VerifyLedger replays the entire ledger: every row hash, every batch root /
// chain link, and every on-chain anchor of both tables.
func (l *Ledger) VerifyLedger(ctx context.Context) (*LedgerReport, error) {
	rep := &LedgerReport{Pass: true}

	for _, table := range db.LedgerTables {
		if err := l.verifyTable(ctx, table, rep); err != nil {
			return nil, err
		}
	}
	return rep, nil
}

func (l *Ledger) verifyTable(ctx context.Context, table string, rep *LedgerReport) error {
	// ── Rows: content recomputation + chain continuity ──
	var prevHash []byte
	var rowsChecked int
	switch table {
	case db.TableAccessLogs:
		rows, err := l.database.AllAccessLogRows(ctx)
		if err != nil {
			return err
		}
		for _, r := range rows {
			rowsChecked++
			content := hash.AccessLogContentHash(
				r.AccessedAt, r.AccessorAddress, r.PatientAddress, r.RecordID,
				r.AccessType, r.PolicyMatched, r.WasEmergency, r.TxHash)
			recomputed := hash.RowHash(r.PrevHash, content)
			if !bytes.Equal(hash.FieldBytes(recomputed), r.RowHash) {
				rep.addIssue(table, "row %d: row_hash mismatch (content or prev_hash tampered)", r.Seq)
			}
			if !bytes.Equal(zeroIfNil(r.PrevHash), zeroIfNil(prevHash)) {
				rep.addIssue(table, "row %d: prev_hash breaks chain continuity", r.Seq)
			}
			prevHash = r.RowHash
		}
	case db.TableEventLogs:
		rows, err := l.database.AllEventLogRows(ctx)
		if err != nil {
			return err
		}
		for _, r := range rows {
			rowsChecked++
			content := hash.EventLogContentHash(
				r.EventType, r.ContractAddress, r.TxHash, r.BlockNumber, r.Data, r.LoggedAt)
			recomputed := hash.RowHash(r.PrevHash, content)
			if !bytes.Equal(hash.FieldBytes(recomputed), r.RowHash) {
				rep.addIssue(table, "row %d: row_hash mismatch (content or prev_hash tampered)", r.Seq)
			}
			if !bytes.Equal(zeroIfNil(r.PrevHash), zeroIfNil(prevHash)) {
				rep.addIssue(table, "row %d: prev_hash breaks chain continuity", r.Seq)
			}
			prevHash = r.RowHash
		}
	}
	rep.RowsChecked += rowsChecked

	// ── Batches: merkle roots, chain links, anchors ──
	batches, err := l.database.AllBatches(ctx, table)
	if err != nil {
		return err
	}
	rep.TablesChecked++

	var prevChain []byte
	for i, b := range batches {
		rep.BatchesChecked++
		leafHashes, err := l.database.BatchRowsForTree(ctx, b.LogTable, b.StartSeq, b.EndSeq)
		if err != nil {
			return err
		}
		leaves := make([]*big.Int, len(leafHashes))
		for i, h := range leafHashes {
			leaves[i] = hash.BytesToField32(h)
		}
		root := hash.MerkleRoot(leaves)
		rootHex := hex.EncodeToString(hash.FieldBytes(root))
		if rootHex != b.MerkleRoot {
			rep.addIssue(table, "batch %s (%d): merkle root mismatch: stored %s, recomputed %s",
				b.BatchID, b.BatchSeq, b.MerkleRoot, rootHex)
		}

		if !bytes.Equal(zeroIfNil(b.PrevChainHash), zeroIfNil(prevChain)) {
			rep.addIssue(table, "batch %d: prev_chain_hash breaks chain continuity", b.BatchSeq)
		}

		prevF := hash.Zero
		if len(b.PrevChainHash) > 0 {
			prevF = hash.BytesToField(b.PrevChainHash)
		}
		chain := hash.Poseidon2(prevF, root)
		if !bytes.Equal(hash.FieldBytes(chain), b.ChainHash) {
			rep.addIssue(table, "batch %d: chain hash mismatch", b.BatchSeq)
		}
		prevChain = b.ChainHash

		if b.Anchored {
			rep.AnchorsChecked++
			if l.anchor == nil {
				rep.addIssue(table, "batch %d: anchored but on-chain anchor unavailable", b.BatchSeq)
				continue
			}
			onChain, err := l.anchor.Anchors(nil, big.NewInt(b.AnchorBatchID))
			if err != nil {
				return fmt.Errorf("read anchor %d: %w", b.AnchorBatchID, err)
			}
			if onChain.ChainHash.Cmp(hash.BytesToField(b.ChainHash)) != 0 {
				rep.addIssue(table, "batch %d: on-chain anchor %d hash mismatch", b.BatchSeq, b.AnchorBatchID)
			}
		} else if i != len(batches)-1 && len(batches) > 0 {
			rep.addIssue(table, "batch %d: not anchored (only the newest batch may be pending)", b.BatchSeq)
		}
	}

	return nil
}

func (rep *LedgerReport) addIssue(table, format string, args ...any) {
	rep.Issues = append(rep.Issues, fmt.Sprintf("[%s] "+format, append([]any{table}, args...)...))
	rep.Pass = false
	rep.UnexpectedFailures++
}

func zeroIfNil(b []byte) []byte {
	if b == nil {
		return hash.FieldBytes(hash.Zero)
	}
	return b
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}