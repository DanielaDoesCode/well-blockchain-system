package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Thesis/well-blockchain/well-backend/internal/ledger/hash"

	"github.com/jackc/pgx/v5"
)

type AccessLog struct {
	AccessorAddress string
	PatientAddress  string
	RecordID        string
	AccessType      string
	PolicyMatched   string
	WasEmergency    bool
	TxHash          string
}

type EventLog struct {
	EventType       string
	ContractAddress string
	TxHash          string
	BlockNumber     uint64
	Data            map[string]any
}

// InsertAccessLog appends a row to the access_logs ledger. The row is hashed
// into the chain (row_hash = Poseidon(prev_hash, content)) in the same
// transaction that stores it, so a stale or tampered row is undetectable only
// if the whole anchored chain is rewritten.
func (db *DB) InsertAccessLog(ctx context.Context, log AccessLog) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert access log: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", ledgerAdvisoryKey(TableAccessLogs)); err != nil {
		return fmt.Errorf("insert access log: lock: %w", err)
	}

	prevHash, err := ledgerTip(ctx, tx, TableAccessLogs)
	if err != nil {
		return fmt.Errorf("insert access log: %w", err)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	content := hash.AccessLogContentHash(
		now,
		log.AccessorAddress,
		log.PatientAddress,
		log.RecordID,
		log.AccessType,
		log.PolicyMatched,
		log.WasEmergency,
		log.TxHash,
	)
	rowHash := hash.FieldBytes(hash.RowHash(prevHash, content))
	if prevHash == nil {
		prevHash = hash.FieldBytes(hash.Zero)
	}

	query := `
		INSERT INTO access_logs (
			accessed_at,
			accessor_address,
			patient_address,
			record_id,
			access_type,
			policy_matched,
			was_emergency,
			tx_hash,
			prev_hash,
			row_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = tx.Exec(ctx, query,
		now,
		log.AccessorAddress,
		log.PatientAddress,
		log.RecordID,
		log.AccessType,
		log.PolicyMatched,
		log.WasEmergency,
		log.TxHash,
		prevHash,
		rowHash,
	)
	if err != nil {
		return fmt.Errorf("insert access log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("insert access log: commit: %w", err)
	}
	return nil
}

// InsertEventLog appends a row to the event_logs ledger (see InsertAccessLog).
func (db *DB) InsertEventLog(ctx context.Context, log EventLog) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert event log: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", ledgerAdvisoryKey(TableEventLogs)); err != nil {
		return fmt.Errorf("insert event log: lock: %w", err)
	}

	prevHash, err := ledgerTip(ctx, tx, TableEventLogs)
	if err != nil {
		return fmt.Errorf("insert event log: %w", err)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	dataJSON, err := json.Marshal(log.Data)
	if err != nil {
		return fmt.Errorf("insert event log: marshal data: %w", err)
	}
	content := hash.EventLogContentHash(
		log.EventType,
		log.ContractAddress,
		log.TxHash,
		log.BlockNumber,
		dataJSON,
		now,
	)
	rowHash := hash.FieldBytes(hash.RowHash(prevHash, content))
	if prevHash == nil {
		prevHash = hash.FieldBytes(hash.Zero)
	}

	query := `
		INSERT INTO event_logs (
			event_type,
			contract_address,
			tx_hash,
			block_number,
			data,
			logged_at,
			prev_hash,
			row_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.Exec(ctx, query,
		log.EventType,
		log.ContractAddress,
		log.TxHash,
		log.BlockNumber,
		log.Data,
		now,
		prevHash,
		rowHash,
	)
	if err != nil {
		return fmt.Errorf("insert event log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("insert event log: commit: %w", err)
	}
	return nil
}

// marshalData canonicalises a data map to the same byte form used for hashing,
// so the hash is computed over exactly what JSON the row carries.

// InsertEventLogTx appends an event log inside an existing transaction WITHOUT
// ledger chaining. It is kept for callers that manage their own transaction and
// is not part of the tamper-evident ledger.
func (db *DB) InsertEventLogTx(ctx context.Context, tx pgx.Tx, log EventLog) error {
	query := `
		INSERT INTO event_logs (
			event_type,
			contract_address,
			tx_hash,
			block_number,
			data,
			logged_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.Exec(ctx, query,
		log.EventType,
		log.ContractAddress,
		log.TxHash,
		log.BlockNumber,
		log.Data,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert event log: %w", err)
	}

	return nil
}
