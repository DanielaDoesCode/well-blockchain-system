package db

import (
	"context"
	"fmt"
	"time"
)

type EncryptedKeyRecord struct {
	RecordAddress  string
	IPFSCid        string
	ContentHash    string
	EncryptedKey   []byte
	PatientAddress string
	// LockboxSalt and PatientPublicKey are set for PRE records (lockbox in
	// EncryptedKey). They are nil for legacy Vault-wrapped rows.
	LockboxSalt      []byte
	PatientPublicKey []byte
	CreatedAt        time.Time
}

type ReEncryptedKeyRecord struct {
	ID             int
	RecordAddress  string
	DoctorAddress  string
	ReEncryptedKey []byte
	CreatedAt      time.Time
}

func (db *DB) InsertEncryptedKey(ctx context.Context, r *EncryptedKeyRecord) error {
	query := `
		INSERT INTO encrypted_content_keys (record_address, ipfs_cid, content_hash, encrypted_key, patient_address, lockbox_salt, patient_public_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (record_address) DO UPDATE SET
			ipfs_cid = EXCLUDED.ipfs_cid,
			content_hash = EXCLUDED.content_hash,
			encrypted_key = EXCLUDED.encrypted_key,
			patient_address = EXCLUDED.patient_address,
			lockbox_salt = EXCLUDED.lockbox_salt,
			patient_public_key = EXCLUDED.patient_public_key
	`
	_, err := db.Pool.Exec(ctx, query,
		r.RecordAddress, r.IPFSCid, r.ContentHash, r.EncryptedKey, r.PatientAddress, r.LockboxSalt, r.PatientPublicKey,
	)
	if err != nil {
		return fmt.Errorf("insert encrypted key: %w", err)
	}
	return nil
}

func (db *DB) GetEncryptedKey(ctx context.Context, recordAddress string) (*EncryptedKeyRecord, error) {
	query := `
		SELECT record_address, ipfs_cid, content_hash, encrypted_key, patient_address, lockbox_salt, patient_public_key, created_at
		FROM encrypted_content_keys
		WHERE record_address = $1
	`
	r := &EncryptedKeyRecord{}
	err := db.Pool.QueryRow(ctx, query, recordAddress).Scan(
		&r.RecordAddress, &r.IPFSCid, &r.ContentHash, &r.EncryptedKey, &r.PatientAddress, &r.LockboxSalt, &r.PatientPublicKey, &r.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get encrypted key: %w", err)
	}
	return r, nil
}

// UpsertReEncryptionKey stores rk_{patient->doctor} supplied by the patient.
func (db *DB) UpsertReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string, rk []byte) error {
	query := `
		INSERT INTO re_encryption_keys (patient_address, doctor_address, re_encryption_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (patient_address, doctor_address) DO UPDATE SET
			re_encryption_key = EXCLUDED.re_encryption_key,
			created_at = NOW()
	`
	if _, err := db.Pool.Exec(ctx, query, patientAddress, doctorAddress, rk); err != nil {
		return fmt.Errorf("upsert re-encryption key: %w", err)
	}
	return nil
}

// GetReEncryptionKey returns rk_{patient->doctor}, or a wrapped error if absent.
func (db *DB) GetReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string) ([]byte, error) {
	query := `
		SELECT re_encryption_key
		FROM re_encryption_keys
		WHERE patient_address = $1 AND doctor_address = $2
	`
	var rk []byte
	if err := db.Pool.QueryRow(ctx, query, patientAddress, doctorAddress).Scan(&rk); err != nil {
		return nil, fmt.Errorf("get re-encryption key: %w", err)
	}
	return rk, nil
}

// DeleteReEncryptionKey removes a delegation. It revokes future re-encryption
// of the patient's lockboxes for this doctor but does not retract copies.
func (db *DB) DeleteReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string) error {
	query := `
		DELETE FROM re_encryption_keys
		WHERE patient_address = $1 AND doctor_address = $2
	`
	if _, err := db.Pool.Exec(ctx, query, patientAddress, doctorAddress); err != nil {
		return fmt.Errorf("delete re-encryption key: %w", err)
	}
	return nil
}

func (db *DB) InsertReEncryptedKey(ctx context.Context, r *ReEncryptedKeyRecord) error {
	query := `
		INSERT INTO reencrypted_keys (record_address, doctor_address, reencrypted_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (record_address, doctor_address) DO UPDATE SET
			reencrypted_key = EXCLUDED.reencrypted_key,
			created_at = NOW()
	`
	_, err := db.Pool.Exec(ctx, query,
		r.RecordAddress, r.DoctorAddress, r.ReEncryptedKey,
	)
	if err != nil {
		return fmt.Errorf("insert reencrypted key: %w", err)
	}
	return nil
}

func (db *DB) GetReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) ([]byte, error) {
	query := `
		SELECT reencrypted_key
		FROM reencrypted_keys
		WHERE record_address = $1 AND doctor_address = $2
	`
	var key []byte
	err := db.Pool.QueryRow(ctx, query, recordAddress, doctorAddress).Scan(&key)
	if err != nil {
		return nil, fmt.Errorf("get reencrypted key: %w", err)
	}
	return key, nil
}

func (db *DB) DeleteReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) error {
	query := `
		DELETE FROM reencrypted_keys
		WHERE record_address = $1 AND doctor_address = $2
	`
	_, err := db.Pool.Exec(ctx, query, recordAddress, doctorAddress)
	if err != nil {
		return fmt.Errorf("delete reencrypted key: %w", err)
	}
	return nil
}

func (db *DB) DeleteAllReEncryptedKeys(ctx context.Context, recordAddress string) error {
	_, err := db.Pool.Exec(ctx,
		"DELETE FROM reencrypted_keys WHERE record_address = $1", recordAddress,
	)
	if err != nil {
		return fmt.Errorf("delete all reencrypted keys: %w", err)
	}
	return nil
}
