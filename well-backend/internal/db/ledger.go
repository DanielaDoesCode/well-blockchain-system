package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Ledger table identifiers. Bathces in batched_logs are tagged with these so
// each table has its own independent hash chain.
const (
	TableAccessLogs = "access_logs"
	TableEventLogs  = "event_logs"
)

// LedgerTables enumerates the ledger tables in sealing order.
var LedgerTables = []string{TableAccessLogs, TableEventLogs}

// advisory keys serialise chain appends per table; appended by the wrapped
// insert functions so two concurrent writers can never fork the chain.
func ledgerAdvisoryKey(table string) int64 {
	if table == TableEventLogs {
		return 787194
	}
	return 787193
}

// ContentColumn lists the hash-relevant columns (in canonical order,
// matching ledger/hash) per table. first is the key used to read them back.
var accessLogCols = "accessed_at, accessor_address, patient_address, record_id, access_type, policy_matched, was_emergency, tx_hash"
var eventLogCols = "event_type, contract_address, tx_hash, block_number, data::text AS data, logged_at"

type LedgerRow struct {
	Seq      int64
	PrevHash []byte
	RowHash  []byte
}

// ledgerTip returns the previous row hash for the table chain inside a running
// transaction. Returns nil when the chain is empty.
func ledgerTip(ctx context.Context, tx pgx.Tx, table string) ([]byte, error) {
	query := fmt.Sprintf("SELECT row_hash FROM %s ORDER BY seq DESC LIMIT 1 FOR UPDATE", table)
	var prev []byte
	err := tx.QueryRow(ctx, query).Scan(&prev)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read ledger tip (%s): %w", table, err)
	}
	return prev, nil
}

// CountUnbatchedRows counts rows not yet sealed into a batch.
func (db *DB) CountUnbatchedRows(ctx context.Context, table string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE batch_id IS NULL", table)
	var n int64
	if err := db.Pool.QueryRow(ctx, query).Scan(&n); err != nil {
		return 0, fmt.Errorf("count unbatched rows (%s): %w", table, err)
	}
	return n, nil
}

// GetUnbatchedRows returns up to limit unbatched rows in seq order.
func (db *DB) GetUnbatchedRows(ctx context.Context, table string, limit int) ([]LedgerRow, error) {
	query := fmt.Sprintf(
		"SELECT seq, prev_hash, row_hash FROM %s WHERE batch_id IS NULL ORDER BY seq ASC LIMIT $1", table)
	rows, err := db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get unbatched rows (%s): %w", table, err)
	}
	defer rows.Close()

	var out []LedgerRow
	for rows.Next() {
		var r LedgerRow
		if err := rows.Scan(&r.Seq, &r.PrevHash, &r.RowHash); err != nil {
			return nil, fmt.Errorf("scan unbatched row (%s): %w", table, err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// InsertBatch seals start..end into a new batch (block): it stores the Merkle
// root and chain hash, then claims the rows. Both steps are one transaction.
// Returns the batch's UUID and its per-table batch sequence number.
func (db *DB) InsertBatch(
	ctx context.Context,
	table string,
	startSeq, endSeq int64,
	leafCount int,
	merkleRoot string,
	merkleRootBytes, prevChainHash, chainHash []byte,
) (batchID string, batchSeq int64, err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("begin batch tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := tx.QueryRow(ctx, `
		INSERT INTO batched_logs (
			log_table, batch_seq, merkle_root, prev_chain_hash, chain_hash,
			start_seq, end_seq, leaf_count, submitted_at
		) VALUES (
			$1::text, COALESCE((SELECT MAX(batch_seq) FROM batched_logs WHERE log_table = $1::text), 0) + 1,
			$2, $3, $4, $5, $6, $7, NOW()
		) RETURNING batch_id, batch_seq`,
		table,
		merkleRoot,
		prevChainHash,
		chainHash,
		startSeq,
		endSeq,
		leafCount,
	).Scan(&batchID, &batchSeq); err != nil {
		return "", 0, fmt.Errorf("insert batch (%s): %w", table, err)
	}

	// Claim the rows so they can appear in one batch only.
	if _, err := tx.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET batch_id = $1 WHERE seq BETWEEN $2 AND $3", table),
		batchID, startSeq, endSeq); err != nil {
		return "", 0, fmt.Errorf("claim rows (%s): %w", table, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", 0, fmt.Errorf("commit batch (%s): %w", table, err)
	}
	_ = merkleRootBytes // root is stored as hex string; bytes kept for tests
	return batchID, batchSeq, nil
}

// LatestChainHash returns the chain hash of the most recent batch for the
// table, or nil when no batch exists yet (chain starts at field zero).
func (db *DB) LatestChainHash(ctx context.Context, table string) ([]byte, error) {
	query := `SELECT chain_hash FROM batched_logs WHERE log_table = $1 ORDER BY batch_seq DESC LIMIT 1`
	var h []byte
	err := db.Pool.QueryRow(ctx, query, table).Scan(&h)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("latest chain hash (%s): %w", table, err)
	}
	return h, nil
}

// LatestAnchoredChain returns the chain hash and anchor batch id of the most
// recently anchored batch for a table (nil/zero when nothing anchored).
func (db *DB) LatestAnchoredChain(ctx context.Context, table string) ([]byte, uint64, error) {
	query := `SELECT chain_hash, anchor_batch_id FROM batched_logs
		WHERE log_table = $1 AND anchored = TRUE
		ORDER BY batch_seq DESC LIMIT 1`
	var h []byte
	var anchorID *int64
	err := db.Pool.QueryRow(ctx, query, table).Scan(&h, &anchorID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("latest anchored chain (%s): %w", table, err)
	}
	if anchorID == nil {
		return nil, 0, nil
	}
	return h, uint64(*anchorID), nil
}

// MarkBatchAnchored records the on-chain anchor for a batch.
func (db *DB) MarkBatchAnchored(ctx context.Context, batchID string, anchorBatchID uint64, txHash string) error {
	if _, err := db.Pool.Exec(ctx, `
		UPDATE batched_logs
		SET anchored = TRUE, anchor_batch_id = $2, anchor_tx_hash = $3
		WHERE batch_id = $1`,
		batchID, int64(anchorBatchID), txHash); err != nil {
		return fmt.Errorf("mark batch anchored: %w", err)
	}
	return nil
}

// OldestUnanchoredBatch returns the earliest batch (by batch_seq) that has not
// yet been anchored on-chain, or nil when none are outstanding.
func (db *DB) OldestUnanchoredBatch(ctx context.Context, table string) (*BatchStatus, error) {
	var b BatchStatus
	var anchorBatch *int64
	var anchorTx *string
	err := db.Pool.QueryRow(ctx, `
		SELECT batch_id, batch_seq, log_table, merkle_root, prev_chain_hash,
		       chain_hash, start_seq, end_seq, leaf_count, anchored, anchor_batch_id, anchor_tx_hash
		FROM batched_logs WHERE log_table = $1 AND anchored = FALSE
		ORDER BY batch_seq ASC LIMIT 1`,
		table).Scan(&b.BatchID, &b.BatchSeq, &b.LogTable, &b.MerkleRoot, &b.PrevChainHash,
		&b.ChainHash, &b.StartSeq, &b.EndSeq, &b.LeafCount, &b.Anchored, &anchorBatch, &anchorTx)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("oldest unanchored batch (%s): %w", table, err)
	}
	if anchorBatch != nil {
		b.AnchorBatchID = *anchorBatch
	}
	if anchorTx != nil {
		b.AnchorTxHash = *anchorTx
	}
	return &b, nil
}

type BatchStatus struct {
	BatchID       string
	BatchSeq      int64
	LogTable      string
	MerkleRoot    string
	PrevChainHash []byte
	ChainHash     []byte
	StartSeq      int64
	EndSeq        int64
	LeafCount     int
	Anchored      bool
	AnchorBatchID int64
	AnchorTxHash  string
}

// BatchByID returns batch metadata for a batch id.
func (db *DB) BatchByID(ctx context.Context, batchID string) (BatchStatus, error) {
	var b BatchStatus
	var anchorBatch *int64
	var anchorTx *string
	err := db.Pool.QueryRow(ctx, `
		SELECT batch_id, batch_seq, log_table, merkle_root, prev_chain_hash,
		       chain_hash, start_seq, end_seq, leaf_count, anchored, anchor_batch_id, anchor_tx_hash
		FROM batched_logs WHERE batch_id = $1`, batchID).
		Scan(&b.BatchID, &b.BatchSeq, &b.LogTable, &b.MerkleRoot, &b.PrevChainHash,
			&b.ChainHash, &b.StartSeq, &b.EndSeq, &b.LeafCount, &b.Anchored, &anchorBatch, &anchorTx)
	if err != nil {
		return b, fmt.Errorf("batch by id: %w", err)
	}
	if anchorBatch != nil {
		b.AnchorBatchID = *anchorBatch
	}
	if anchorTx != nil {
		b.AnchorTxHash = *anchorTx
	}
	return b, nil
}

// BatchRowsForTree returns the row hashes of every row claimed by a batch, in
// seq order — used to rebuild the Merkle tree during verification.
func (db *DB) BatchRowsForTree(ctx context.Context, table string, startSeq, endSeq int64) ([][]byte, error) {
	query := fmt.Sprintf(
		"SELECT row_hash FROM %s WHERE seq BETWEEN $1 AND $2 ORDER BY seq ASC", table)
	rows, err := db.Pool.Query(ctx, query, startSeq, endSeq)
	if err != nil {
		return nil, fmt.Errorf("batch rows (%s): %w", table, err)
	}
	defer rows.Close()

	var out [][]byte
	for rows.Next() {
		var h []byte
		if err := rows.Scan(&h); err != nil {
			return nil, fmt.Errorf("scan batch row hash: %w", err)
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

// BatchSeqByRow finds the batch sequence covering a row (via its batch_id),
// returning also the batch UUID. row not yet batched → (nil string, false).
func (db *DB) BatchSeqByRow(ctx context.Context, table string, seq int64) (string, int64, bool, error) {
	var batchID *string
	query := fmt.Sprintf("SELECT batch_id FROM %s WHERE seq = $1", table)
	if err := db.Pool.QueryRow(ctx, query, seq).Scan(&batchID); err != nil {
		return "", 0, false, fmt.Errorf("batch seq by row (%s %d): %w", table, seq, err)
	}
	if batchID == nil {
		return "", 0, false, nil
	}
	b, err := db.BatchByID(ctx, *batchID)
	if err != nil {
		return "", 0, false, err
	}
	return b.BatchID, b.BatchSeq, true, nil
}

// AllBatches returns every batch for a table in batch_seq order.
func (db *DB) AllBatches(ctx context.Context, table string) ([]BatchStatus, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT batch_id, batch_seq, log_table, merkle_root, prev_chain_hash,
		       chain_hash, start_seq, end_seq, leaf_count, anchored, anchor_batch_id, anchor_tx_hash
		FROM batched_logs WHERE log_table = $1 ORDER BY batch_seq ASC`, table)
	if err != nil {
		return nil, fmt.Errorf("all batches (%s): %w", table, err)
	}
	defer rows.Close()

	var out []BatchStatus
	for rows.Next() {
		var b BatchStatus
		var anchorBatch *int64
		var anchorTx *string
		if err := rows.Scan(&b.BatchID, &b.BatchSeq, &b.LogTable, &b.MerkleRoot, &b.PrevChainHash,
			&b.ChainHash, &b.StartSeq, &b.EndSeq, &b.LeafCount, &b.Anchored, &anchorBatch, &anchorTx); err != nil {
			return nil, fmt.Errorf("scan batch (%s): %w", table, err)
		}
		if anchorBatch != nil {
			b.AnchorBatchID = *anchorBatch
		}
		if anchorTx != nil {
			b.AnchorTxHash = *anchorTx
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// AllRows returns every row hash for a table in seq order.
func (db *DB) AllRows(ctx context.Context, table string) (map[int64][]byte, []int64, error) {
	query := fmt.Sprintf("SELECT seq, row_hash FROM %s ORDER BY seq ASC", table)
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("all rows (%s): %w", table, err)
	}
	defer rows.Close()

	bySeq := map[int64][]byte{}
	var order []int64
	for rows.Next() {
		var seq int64
		var h []byte
		if err := rows.Scan(&seq, &h); err != nil {
			return nil, nil, fmt.Errorf("scan row (%s): %w", table, err)
		}
		bySeq[seq] = h
		order = append(order, seq)
	}
	return bySeq, order, rows.Err()
}

// PrevRowHash returns the row hash of the row immediately before seq (by seq
// ordering), or nil when seq is the chain head.
func (db *DB) PrevRowHash(ctx context.Context, table string, seq int64) ([]byte, error) {
	query := fmt.Sprintf(
		"SELECT row_hash FROM %s WHERE seq = (SELECT MAX(seq) FROM %s WHERE seq < $1)", table, table)
	var h []byte
	err := db.Pool.QueryRow(ctx, query, seq).Scan(&h)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("prev row hash (%s %d): %w", table, seq, err)
	}
	return h, nil
}
