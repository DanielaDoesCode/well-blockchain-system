package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrIdentityNotFound = errors.New("identity not found")

// Identity is one row of the off-chain identity index. The PII fields
// (national ID, display name) are written by an admin/issuer at enrollment;
// the on-chain facts (address, commitment, summary, internalID) are refreshed
// by the chain-event listener.
type Identity struct {
	NationalID     string
	DisplayName    string
	Role           string
	Address        string
	SummaryAddress string
	Commitment     string
	InternalID     string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type IdentityLookup struct {
	CallerAddress string
	QueryKey      string
	QueryValue    string
	ResultCount   int
	LookedAt      time.Time
}

// orNil maps an empty string to SQL NULL so partial upserts preserve existing values.
func orNil(v string) any {
	if v == "" {
		return nil
	}
	return v
}

func (db *DB) UpsertIdentity(ctx context.Context, e *Identity) error {
	query := `
		INSERT INTO identities (national_id, display_name, role, address, summary_address, commitment, internal_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8, 'active'))
		ON CONFLICT (address) DO UPDATE SET
			national_id = COALESCE(EXCLUDED.national_id, identities.national_id),
			display_name = COALESCE(EXCLUDED.display_name, identities.display_name),
			role = COALESCE(EXCLUDED.role, identities.role),
			summary_address = COALESCE(EXCLUDED.summary_address, identities.summary_address),
			commitment = COALESCE(EXCLUDED.commitment, identities.commitment),
			internal_id = COALESCE(EXCLUDED.internal_id, identities.internal_id),
			status = COALESCE(EXCLUDED.status, identities.status),
			updated_at = NOW()
	`
	status := e.Status
	if status == "" {
		status = "active"
	}
	_, err := db.Pool.Exec(ctx, query,
		e.NationalID, e.DisplayName, e.Role, e.Address,
		orNil(e.SummaryAddress), orNil(e.Commitment), orNil(e.InternalID), status,
	)
	if err != nil {
		return fmt.Errorf("upsert identity: %w", err)
	}
	return nil
}

// AttachCommitment refreshes the commitment known for an enrolled address.
func (db *DB) AttachCommitment(ctx context.Context, address, commitment string) error {
	if _, err := db.Pool.Exec(ctx, `
		UPDATE identities SET commitment = $1, updated_at = NOW() WHERE address = $2
	`, commitment, address); err != nil {
		return fmt.Errorf("attach commitment: %w", err)
	}
	return nil
}

// AttachSummary refreshes the deployed EmergencySummary address for an enrolled patient.
func (db *DB) AttachSummary(ctx context.Context, address, summary string) error {
	if _, err := db.Pool.Exec(ctx, `
		UPDATE identities SET summary_address = $1, updated_at = NOW() WHERE address = $2
	`, summary, address); err != nil {
		return fmt.Errorf("attach summary: %w", err)
	}
	return nil
}

// AttachRegisteredIdentity binds the on-chain commitment/internalID/role to an
// already-enrolled identity (UserRegistry.Registered).
func (db *DB) AttachRegisteredIdentity(ctx context.Context, commitment, internalID, role string) error {
	if _, err := db.Pool.Exec(ctx, `
		UPDATE identities
		SET internal_id = COALESCE($1, internal_id),
		    commitment = COALESCE($2, commitment),
		    role = COALESCE($3, role),
		    updated_at = NOW()
		WHERE commitment = $2 OR internal_id = $1
	`, orNil(internalID), commitment, role); err != nil {
		return fmt.Errorf("attach registered identity: %w", err)
	}
	return nil
}

// TombstoneByCommitment marks an identity revoked (UserRegistry.Revoked).
func (db *DB) TombstoneByCommitment(ctx context.Context, commitment string) error {
	if _, err := db.Pool.Exec(ctx, `
		UPDATE identities SET status = 'revoked', updated_at = NOW()
		WHERE commitment = $1 AND status <> 'revoked'
	`, commitment); err != nil {
		return fmt.Errorf("tombstone identity: %w", err)
	}
	return nil
}

// GetIdentityByNationalID resolves the pseudo-random address of an active identity.
func (db *DB) GetIdentityByNationalID(ctx context.Context, nationalID string) (*Identity, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT national_id, display_name, role, address,
		       COALESCE(summary_address, ''), COALESCE(commitment, ''),
		       COALESCE(internal_id, ''), status, created_at, updated_at
		FROM identities WHERE national_id = $1
	`, nationalID)
	return scanIdentity(row)
}

// GetIdentitiesByName resolves potential matches by display name (active only).
func (db *DB) GetIdentitiesByName(ctx context.Context, name string) ([]*Identity, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT national_id, display_name, role, address,
		       COALESCE(summary_address, ''), COALESCE(commitment, ''),
		       COALESCE(internal_id, ''), status, created_at, updated_at
		FROM identities
		WHERE display_name = $1 AND status = 'active'
		ORDER BY created_at
	`, name)
	if err != nil {
		return nil, fmt.Errorf("query identities by name: %w", err)
	}
	defer rows.Close()

	out := []*Identity{}
	for rows.Next() {
		e := &Identity{}
		if err := rows.Scan(&e.NationalID, &e.DisplayName, &e.Role, &e.Address,
			&e.SummaryAddress, &e.Commitment, &e.InternalID, &e.Status,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan identity: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (db *DB) DeactivateIdentity(ctx context.Context, nationalID string) error {
	tag, err := db.Pool.Exec(ctx, `
		UPDATE identities SET status = 'revoked', updated_at = NOW() WHERE national_id = $1
	`, nationalID)
	if err != nil {
		return fmt.Errorf("deactivate identity: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrIdentityNotFound
	}
	return nil
}

func (db *DB) RecordLookup(ctx context.Context, l *IdentityLookup) error {
	if _, err := db.Pool.Exec(ctx, `
		INSERT INTO identity_lookups (caller_address, query_key, query_value, result_count)
		VALUES ($1, $2, $3, $4)
	`, l.CallerAddress, l.QueryKey, l.QueryValue, l.ResultCount); err != nil {
		return fmt.Errorf("record identity lookup: %w", err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanIdentity(row rowScanner) (*Identity, error) {
	e := &Identity{}
	err := row.Scan(&e.NationalID, &e.DisplayName, &e.Role, &e.Address,
		&e.SummaryAddress, &e.Commitment, &e.InternalID, &e.Status,
		&e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIdentityNotFound
		}
		return nil, fmt.Errorf("scan identity: %w", err)
	}
	return e, nil
}