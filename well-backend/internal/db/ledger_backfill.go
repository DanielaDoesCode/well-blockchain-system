package db

import (
	"context"
	"fmt"
	"time"

	"Thesis/well-blockchain/well-backend/internal/ledger/hash"
)

// AccessLogRow is a complete access_logs row as stored, used for verification
// and for the legacy backfill.
type AccessLogRow struct {
	Seq              int64
	AccessedAt       time.Time
	AccessorAddress  string
	PatientAddress   string
	RecordID         string
	AccessType       string
	PolicyMatched    string
	WasEmergency     bool
	TxHash           string
	PrevHash         []byte
	RowHash          []byte
	BatchID          *string
}

// EventLogRow is a complete event_logs row as stored.
type EventLogRow struct {
	Seq              int64
	EventType        string
	ContractAddress  string
	TxHash           string
	BlockNumber      uint64
	Data             []byte // canonical JSON text
	LoggedAt         time.Time
	PrevHash         []byte
	RowHash          []byte
	BatchID          *string
}

var accessLogRowCols = "seq, accessed_at, accessor_address, patient_address, record_id, access_type, " +
	"COALESCE(policy_matched, ''), was_emergency, tx_hash, prev_hash, row_hash, batch_id"

var eventLogRowCols = "seq, event_type, contract_address, tx_hash, block_number, data::text, " +
	"logged_at, prev_hash, row_hash, batch_id"

func scanAccessLogRow(row pgxRow) (AccessLogRow, error) {
	var r AccessLogRow
	err := row.Scan(&r.Seq, &r.AccessedAt, &r.AccessorAddress, &r.PatientAddress, &r.RecordID,
		&r.AccessType, &r.PolicyMatched, &r.WasEmergency, &r.TxHash, &r.PrevHash, &r.RowHash, &r.BatchID)
	return r, err
}

func scanEventLogRow(row pgxRow) (EventLogRow, error) {
	var r EventLogRow
	var data string
	err := row.Scan(&r.Seq, &r.EventType, &r.ContractAddress, &r.TxHash, &r.BlockNumber,
		&data, &r.LoggedAt, &r.PrevHash, &r.RowHash, &r.BatchID)
	r.Data = []byte(data)
	return r, err
}

type pgxRow interface {
	Scan(dest ...any) error
}

// GetAccessLogRow fetches one access_logs row by seq.
func (db *DB) GetAccessLogRow(ctx context.Context, seq int64) (AccessLogRow, error) {
	return scanAccessLogRow(db.Pool.QueryRow(ctx,
		fmt.Sprintf("SELECT %s FROM access_logs WHERE seq = $1", accessLogRowCols), seq))
}

// GetEventLogRow fetches one event_logs row by seq.
func (db *DB) GetEventLogRow(ctx context.Context, seq int64) (EventLogRow, error) {
	return scanEventLogRow(db.Pool.QueryRow(ctx,
		fmt.Sprintf("SELECT %s FROM event_logs WHERE seq = $1", eventLogRowCols), seq))
}

// AllAccessLogRows returns every access_logs row in seq order.
func (db *DB) AllAccessLogRows(ctx context.Context) ([]AccessLogRow, error) {
	rows, err := db.Pool.Query(ctx,
		fmt.Sprintf("SELECT %s FROM access_logs ORDER BY seq ASC", accessLogRowCols))
	if err != nil {
		return nil, fmt.Errorf("all access log rows: %w", err)
	}
	defer rows.Close()
	var out []AccessLogRow
	for rows.Next() {
		r, err := scanAccessLogRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan access log row: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// AllEventLogRows returns every event_logs row in seq order.
func (db *DB) AllEventLogRows(ctx context.Context) ([]EventLogRow, error) {
	rows, err := db.Pool.Query(ctx,
		fmt.Sprintf("SELECT %s FROM event_logs ORDER BY seq ASC", eventLogRowCols))
	if err != nil {
		return nil, fmt.Errorf("all event log rows: %w", err)
	}
	defer rows.Close()
	var out []EventLogRow
	for rows.Next() {
		r, err := scanEventLogRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan event log row: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// BackfillLedger hashes any rows created before the ledger feature existed
// (row_hash IS NULL), chaining them in seq order onto the existing chain.
// Idempotent — call on every boot before starting the listener.
func (db *DB) BackfillLedger(ctx context.Context) error {
	accessRows, err := db.AllAccessLogRows(ctx)
	if err != nil {
		return err
	}
	if err := db.backfillAccessRows(ctx, accessRows); err != nil {
		return err
	}

	eventRows, err := db.AllEventLogRows(ctx)
	if err != nil {
		return err
	}
	return db.backfillEventRows(ctx, eventRows)
}

func (db *DB) backfillAccessRows(ctx context.Context, rows []AccessLogRow) error {
	var prev []byte
	for _, r := range rows {
		if len(r.RowHash) == 0 {
			content := hash.AccessLogContentHash(
				r.AccessedAt, r.AccessorAddress, r.PatientAddress, r.RecordID,
				r.AccessType, r.PolicyMatched, r.WasEmergency, r.TxHash)
			h := hash.RowHash(prev, content)
			if _, err := db.Pool.Exec(ctx,
				`UPDATE access_logs SET prev_hash = $1, row_hash = $2 WHERE seq = $3 AND row_hash IS NULL`,
				prev, hash.FieldBytes(h), r.Seq); err != nil {
				return fmt.Errorf("backfill access row %d: %w", r.Seq, err)
			}
			prev = hash.FieldBytes(h)
			continue
		}
		// Already hashed (or a pre-backfill append): advance the chaining only
		// if it is a descendant of what we processed.
		if len(prev) > 0 && bytesEqual(prev, r.PrevHash) {
			prev = r.RowHash
		}
	}
	return nil
}

func (db *DB) backfillEventRows(ctx context.Context, rows []EventLogRow) error {
	var prev []byte
	for _, r := range rows {
		if len(r.RowHash) == 0 {
			content := hash.EventLogContentHash(
				r.EventType, r.ContractAddress, r.TxHash, r.BlockNumber, r.Data, r.LoggedAt)
			h := hash.RowHash(prev, content)
			if _, err := db.Pool.Exec(ctx,
				`UPDATE event_logs SET prev_hash = $1, row_hash = $2 WHERE seq = $3 AND row_hash IS NULL`,
				prev, hash.FieldBytes(h), r.Seq); err != nil {
				return fmt.Errorf("backfill event row %d: %w", r.Seq, err)
			}
			prev = hash.FieldBytes(h)
			continue
		}
		if len(prev) > 0 && bytesEqual(prev, r.PrevHash) {
			prev = r.RowHash
		}
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}