// Package ledger turns the off-chain Postgres audit tables (access_logs,
// event_logs) into a blockchain-style, tamper-evident ledger:
//
//   - every row is hash-chained onto its predecessor at insert time
//     (row_hash = Poseidon(prev_hash, content)),
//   - a background loop seals rows into Merkle-tree batches ("blocks"), each
//     with chain_hash = Poseidon(prev_chain_hash, merkle_root),
//   - each batch's chain hash is anchored on-chain via LogAnchor.anchor(),
//     making the whole history verifiable against the EVM chain.
//
// Verification (VerifySingle / VerifyLedger) recomputes hashes from stored
// data and compares against the anchored chain, detecting any tampering.
package ledger

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/ledger/hash"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Ledger owns the batching + on-chain anchoring loop.
type Ledger struct {
	client    *ethclient.Client
	database  *db.DB
	anchor    *contracts.LogAnchor
	anchorF   *contracts.LogAnchorFilterer
	anchorAt  common.Address
	auth      *bind.TransactOpts
	batchSize int
	interval  time.Duration

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}

	// sealMu serialises all batch creation + anchoring so the poll loop and a
	// shutdown flush (or two loops) can never seal/claim the same rows twice.
	sealMu sync.Mutex
}

// New builds a Ledger. anchorAddr may be empty for verification-only use; when
// privateKeyHex is empty the ledger can verify but not anchor.
func New(
	ctx context.Context,
	rpcURL string,
	database *db.DB,
	anchorAddr, privateKeyHex string,
	batchSize int,
	interval time.Duration,
) (*Ledger, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("ledger: dial ethereum node: %w", err)
	}

	l := &Ledger{
		client:    client,
		database:  database,
		batchSize: batchSize,
		interval:  interval,
		stopCh:    make(chan struct{}),
	}

	if anchorAddr != "" {
		addr := common.HexToAddress(anchorAddr)
		l.anchorAt = addr
		anchor, err := contracts.NewLogAnchor(addr, client)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("ledger: bind log anchor: %w", err)
		}
		l.anchor = anchor
		l.anchorF, err = contracts.NewLogAnchorFilterer(addr, client)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("ledger: bind log anchor filterer: %w", err)
		}
	}

	if privateKeyHex != "" {
		chainID, err := client.ChainID(ctx)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("ledger: get chain ID: %w", err)
		}
		privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("ledger: parse private key: %w", err)
		}
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("ledger: create transactor: %w", err)
		}
		auth.GasLimit = 5_000_000
		auth.GasPrice = big.NewInt(1_000_000_000)
		l.auth = auth
	}

	return l, nil
}

// Start launches the batching/anchor poll loop.
func (l *Ledger) Start(ctx context.Context) {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return
	}
	l.running = true
	l.mu.Unlock()

	go l.pollLoop(ctx)
	log.Printf("ledger started: batchSize=%d, interval=%v", l.batchSize, l.interval)
}

// Stop halts the poll loop.
func (l *Ledger) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.running {
		close(l.stopCh)
		l.running = false
	}
}

// Close stops the loop and the RPC client.
func (l *Ledger) Close() {
	l.Stop()
	l.client.Close()
}

func (l *Ledger) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.SealPending(ctx)
		case <-l.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

// SealPending seals up to batchSize unbatched rows per table into batches and
// anchors each batch's chain hash. It is also called on demand after backfill.
func (l *Ledger) SealPending(ctx context.Context) {
	if l.anchor == nil || l.auth == nil {
		return
	}
	for _, table := range db.LedgerTables {
		select {
		case <-ctx.Done():
			return
		default:
		}
		sealed, err := l.sealBatch(ctx, table)
		if err != nil {
			log.Printf("ledger: seal batch (%s): %v", table, err)
		} else if sealed {
			log.Printf("ledger: sealed a batch for %s", table)
		}
		// Recover batches that were stored but whose anchor transaction failed
		// (or is still pending from a prior tick): anchor in batch_seq order.
		if err := l.anchorUnanchored(ctx, table); err != nil {
			log.Printf("ledger: anchor pending batches (%s): %v", table, err)
		}
	}
}

// sealBatch takes the oldest unbatched rows of a table and wraps them into a
// batch, then anchors the resulting chain hash.
func (l *Ledger) sealBatch(ctx context.Context, table string) (bool, error) {
	l.sealMu.Lock()
	defer l.sealMu.Unlock()

	rows, err := l.database.GetUnbatchedRows(ctx, table, l.batchSize)
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return false, nil
	}

	prevChain, err := l.database.LatestChainHash(ctx, table)
	if err != nil {
		return false, err
	}

	leaves := make([]*big.Int, len(rows))
	for i, r := range rows {
		if len(r.RowHash) == 0 {
			return false, fmt.Errorf("row %d has no row_hash (run BackfillLedger first)", r.Seq)
		}
		leaves[i] = hash.BytesToField32(r.RowHash)
	}
	root := hash.MerkleRoot(leaves)

	prevF := hash.Zero
	prevChainBytes := hash.FieldBytes(hash.Zero)
	if len(prevChain) > 0 {
		prevF = hash.BytesToField(prevChain)
		prevChainBytes = prevChain
	}
	chain := hash.Poseidon2(prevF, root)

	batchID, _, err := l.database.InsertBatch(
		ctx, table,
		rows[0].Seq, rows[len(rows)-1].Seq,
		len(rows),
		hex.EncodeToString(hash.FieldBytes(root)),
		hash.FieldBytes(root),
		prevChainBytes,
		hash.FieldBytes(chain),
	)
	if err != nil {
		return false, err
	}

	anchorBatchID, txHash, err := l.anchorChainHash(ctx, chain)
	if err != nil {
		// Batch is stored; the caller's anchor-retry pass recovers it.
		return true, fmt.Errorf("batch %s stored but not anchored: %w", batchID, err)
	}
	if err := l.database.MarkBatchAnchored(ctx, batchID, anchorBatchID, txHash); err != nil {
		return true, fmt.Errorf("mark anchored: %w", err)
	}

	log.Printf("ledger: table=%s batch=%s seqs=%d..%d rows=%d root=0x%s chain=0x%s anchor=%d tx=%s",
		table, batchID, rows[0].Seq, rows[len(rows)-1].Seq, len(rows),
		hex.EncodeToString(hash.FieldBytes(root)),
		hex.EncodeToString(hash.FieldBytes(chain)),
		anchorBatchID, txHash)
	return true, nil
}

// anchorChainHash records chain on LogAnchor and returns the on-chain batch id.
func (l *Ledger) anchorChainHash(ctx context.Context, chain *big.Int) (uint64, string, error) {
	tx, err := l.anchor.Anchor(l.auth, chain)
	if err != nil {
		return 0, "", fmt.Errorf("anchor tx: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, l.client, tx)
	if err != nil {
		return 0, "", fmt.Errorf("wait mined: %w", err)
	}

	var batchID uint64
	found := false
	for _, lg := range receipt.Logs {
		if lg.Address != l.anchorAt {
			continue
		}
		ev, perr := l.anchorF.ParseAnchored(*lg)
		if perr != nil {
			continue
		}
		batchID = ev.BatchId.Uint64()
		found = true
		break
	}
	if receipt.Status == 0 {
		return 0, "", fmt.Errorf("anchor tx reverted")
	}
	if !found {
		return 0, "", fmt.Errorf("anchored event not found in receipt")
	}
	return batchID, tx.Hash().Hex(), nil
}

// Flush seals every unbatched row into batches and anchors every batch of both
// tables (including any stored-but-unanchored ones), leaving the ledger fully
// anchored up to the current tip. Called at graceful shutdown so a stopped
// backend leaves a complete, verifiable ledger behind.
func (l *Ledger) Flush(ctx context.Context) error {
	if l.anchor == nil || l.auth == nil {
		return nil
	}
	var errs []error
	for _, table := range db.LedgerTables {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := l.flushTable(ctx, table); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", table, err))
			continue
		}
		if err := l.anchorUnanchored(ctx, table); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", table, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("ledger flush: %v", errs)
	}
	return nil
}

// flushTable seals rows into batches until none remain for the table.
func (l *Ledger) flushTable(ctx context.Context, table string) error {
	for {
		rows, err := l.database.GetUnbatchedRows(ctx, table, l.batchSize)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		if _, err := l.sealBatch(ctx, table); err != nil {
			return err
		}
	}
}

// anchorUnanchored anchors every batch of a table that has no on-chain anchor
// yet, in batch_seq order so the on-chain LogAnchor sequence matches the chain.
func (l *Ledger) anchorUnanchored(ctx context.Context, table string) error {
	for {
		batch, err := l.database.OldestUnanchoredBatch(ctx, table)
		if err != nil {
			return err
		}
		if batch == nil {
			return nil
		}
		if err := l.anchorBatch(ctx, *batch); err != nil {
			return err
		}
	}
}

// anchorBatch sends the batch's chain hash to LogAnchor and records the on-chain
// anchor back into batched_logs.
func (l *Ledger) anchorBatch(ctx context.Context, b db.BatchStatus) error {
	l.sealMu.Lock()
	defer l.sealMu.Unlock()

	anchorBatchID, txHash, err := l.anchorChainHash(ctx, hash.BytesToField(b.ChainHash))
	if err != nil {
		return fmt.Errorf("anchor batch %s (seq=%d): %w", b.BatchID, b.BatchSeq, err)
	}
	if err := l.database.MarkBatchAnchored(ctx, b.BatchID, anchorBatchID, txHash); err != nil {
		return fmt.Errorf("mark anchored %s: %w", b.BatchID, err)
	}
	log.Printf("ledger: table=%s batch=%s seq=%d anchored chain=0x%s anchor=%d tx=%s",
		b.LogTable, b.BatchID, b.BatchSeq,
		hex.EncodeToString(b.ChainHash), anchorBatchID, txHash)
	return nil
}
