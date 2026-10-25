-- Each sealed batch is a "block": a Merkle root over the batch's row hashes,
-- stored together with the previous chain hash. The chained value
--   chain_hash_i = Poseidon2(chain_hash_{i-1}, merkle_root_i)
-- is anchored on-chain via LogAnchor.anchor().
CREATE TABLE IF NOT EXISTS batched_logs(
    batch_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    log_table VARCHAR NOT NULL,
    batch_seq BIGINT NOT NULL,
    merkle_root VARCHAR NOT NULL,
    prev_chain_hash BYTEA NOT NULL,
    chain_hash BYTEA NOT NULL,
    start_seq BIGINT NOT NULL,
    end_seq BIGINT NOT NULL,
    leaf_count INT NOT NULL,
    zk_proof BYTEA,
    anchored BOOLEAN NOT NULL DEFAULT FALSE,
    anchor_batch_id BIGINT,
    anchor_tx_hash CHAR(66),
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(log_table, batch_seq)
);

-- The tamper-evident ledger columns (seq / prev_hash / row_hash / batch_id)
-- are written by the backend in the same transaction as the row itself:
--   row_hash    = Poseidon2(prev_hash, contentHash(canonical row))
--   batch_id    = set by the batcher when the row is sealed into a Merkle batch.
-- seq is the per-table chain position (BIGSERIAL = monotonic, never reused).

CREATE TABLE IF NOT EXISTS access_logs (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accessed_at TIMESTAMP NOT NULL,
    accessor_address CHAR(42) NOT NULL,
    patient_address CHAR(42) NOT NULL,
    record_id CHAR(42) NOT NULL,
    access_type VARCHAR NOT NULL,
    policy_matched VARCHAR,
    was_emergency BOOLEAN NOT NULL,
    tx_hash CHAR(66) NOT NULL,
    seq BIGSERIAL NOT NULL UNIQUE,
    prev_hash BYTEA,
    row_hash BYTEA,
    batch_id UUID REFERENCES batched_logs(batch_id)
);

CREATE TABLE IF NOT EXISTS event_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR NOT NULL,
    contract_address CHAR(42) NOT NULL,
    tx_hash CHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    data JSONB NOT NULL,
    logged_at TIMESTAMP NOT NULL DEFAULT NOW(),
    seq BIGSERIAL NOT NULL UNIQUE,
    prev_hash BYTEA,
    row_hash BYTEA,
    batch_id UUID REFERENCES batched_logs(batch_id)
);

CREATE INDEX IF NOT EXISTS idx_event_logs_event_type ON event_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_event_logs_tx_hash ON event_logs(tx_hash);
CREATE INDEX IF NOT EXISTS idx_event_logs_block_number ON event_logs(block_number);

CREATE TABLE IF NOT EXISTS pending_contributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commitment CHAR(66) NOT NULL,
    nullifier CHAR(66) NOT NULL,
    -- The on-chain AttestationManager root at the time of submission.
    -- Stored as decimal string to match big.Int serialisation in Go.
    -- The aggregator re-uses this root when calling ResearchContributer.contribute().
    merkle_root VARCHAR NOT NULL,
    pA_0 VARCHAR NOT NULL,
    pA_1 VARCHAR NOT NULL,
    pB_0_0 VARCHAR NOT NULL,
    pB_0_1 VARCHAR NOT NULL,
    pB_1_0 VARCHAR NOT NULL,
    pB_1_1 VARCHAR NOT NULL,
    pC_0 VARCHAR NOT NULL,
    pC_1 VARCHAR NOT NULL,
    condition_group VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS research_counters (
    condition_group VARCHAR PRIMARY KEY,
    total_count BIGINT NOT NULL DEFAULT 0,
    published_count BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- PRE lockbox for each record. `encrypted_key` is the canonical AFGH
-- second-level delegator ciphertext C1 (Gt || G1, 624 bytes) that only the
-- patient can open. `lockbox_salt` is the per-record HKDF salt and
-- `patient_public_key` is the patient's PRE public key (compressed G1+G2).
-- Legacy rows written before the PRE migration hold a Vault transit blob in
-- `encrypted_key` and leave the two new columns NULL.
CREATE TABLE IF NOT EXISTS encrypted_content_keys (
    record_address  CHAR(42) PRIMARY KEY,
    ipfs_cid        TEXT NOT NULL,
    content_hash    CHAR(66) NOT NULL,
    encrypted_key   BYTEA NOT NULL,
    patient_address CHAR(42) NOT NULL,
    lockbox_salt    BYTEA,
    patient_public_key BYTEA,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Re-encrypted keys for each authorized doctor (C_doctor)
CREATE TABLE IF NOT EXISTS reencrypted_keys (
    id              SERIAL PRIMARY KEY,
    record_address  CHAR(42) NOT NULL REFERENCES encrypted_content_keys(record_address) ON DELETE CASCADE,
    doctor_address  CHAR(42) NOT NULL,
    reencrypted_key BYTEA NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(record_address, doctor_address)
);

CREATE INDEX IF NOT EXISTS idx_reencrypted_keys_record ON reencrypted_keys(record_address);
CREATE INDEX IF NOT EXISTS idx_reencrypted_keys_doctor ON reencrypted_keys(doctor_address);

-- Patient-to-doctor PRE re-encryption keys rk_{A->B}, supplied by the patient
-- (who alone holds their PRE secret key). The backend uses these to re-encrypt
-- a record's lockbox for the doctor; it can transform but never open the
-- lockbox. Deleting a row revokes future re-encryption but is not retroactive.
CREATE TABLE IF NOT EXISTS re_encryption_keys (
    patient_address   CHAR(42) NOT NULL,
    doctor_address    CHAR(42) NOT NULL,
    re_encryption_key BYTEA NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (patient_address, doctor_address)
);

CREATE INDEX IF NOT EXISTS idx_re_encryption_keys_doctor ON re_encryption_keys(doctor_address);

-- Off-chain identity index: binds a real-world national ID / display name to the
-- pseudonymous on-chain identity (address, commitment, internalID, summary).
-- Populated by (a) admin enrollment and (b) the chain-event listener.
CREATE TABLE IF NOT EXISTS identities (
    address         CHAR(42) PRIMARY KEY,
    national_id     VARCHAR NOT NULL UNIQUE,
    display_name    VARCHAR NOT NULL,
    role            VARCHAR NOT NULL,           -- patient | doctor | first_responder | hospital | admin
    summary_address CHAR(42),
    commitment      CHAR(66),
    internal_id     CHAR(66),
    status          VARCHAR NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_identities_name   ON identities(display_name);
CREATE INDEX IF NOT EXISTS idx_identities_status ON identities(status);
CREATE INDEX IF NOT EXISTS idx_identities_role   ON identities(role);

-- Audit trail: every identity index lookup (who looked up whom / what).
CREATE TABLE IF NOT EXISTS identity_lookups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caller_address CHAR(42) NOT NULL,
    query_key VARCHAR NOT NULL,
    query_value VARCHAR NOT NULL,
    result_count INT NOT NULL,
    looked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_identity_lookups_caller ON identity_lookups(caller_address);
CREATE INDEX IF NOT EXISTS idx_identity_lookups_looked_at ON identity_lookups(looked_at);

CREATE INDEX IF NOT EXISTS idx_access_logs_seq ON access_logs(seq);
CREATE INDEX IF NOT EXISTS idx_access_logs_batch ON access_logs(batch_id);
CREATE INDEX IF NOT EXISTS idx_event_logs_seq ON event_logs(seq);
CREATE INDEX IF NOT EXISTS idx_event_logs_batch ON event_logs(batch_id);
CREATE INDEX IF NOT EXISTS idx_batched_logs_anchored ON batched_logs(anchored);

-- ── Ledger Migrations ─────────────────────────────
-- init.sql is re-run on fresh dev databases; these ALTERs upgrade databases
-- created before the ledger feature existed.
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS log_table VARCHAR NOT NULL DEFAULT 'access_logs';
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS batch_seq BIGINT;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS merkle_root VARCHAR;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS prev_chain_hash BYTEA;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS chain_hash BYTEA;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS start_seq BIGINT;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS end_seq BIGINT;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS leaf_count INT;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS anchored BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS anchor_batch_id BIGINT;
ALTER TABLE batched_logs ADD COLUMN IF NOT EXISTS anchor_tx_hash CHAR(66);

ALTER TABLE access_logs ADD COLUMN IF NOT EXISTS seq BIGSERIAL;
ALTER TABLE access_logs ADD COLUMN IF NOT EXISTS prev_hash BYTEA;
ALTER TABLE access_logs ADD COLUMN IF NOT EXISTS row_hash BYTEA;
ALTER TABLE access_logs ADD COLUMN IF NOT EXISTS batch_id UUID REFERENCES batched_logs(batch_id);

ALTER TABLE event_logs ADD COLUMN IF NOT EXISTS seq BIGSERIAL;
ALTER TABLE event_logs ADD COLUMN IF NOT EXISTS prev_hash BYTEA;
ALTER TABLE event_logs ADD COLUMN IF NOT EXISTS row_hash BYTEA;
ALTER TABLE event_logs ADD COLUMN IF NOT EXISTS batch_id UUID REFERENCES batched_logs(batch_id);

-- PRE migration: upgrade encrypted_content_keys created before the lockbox
-- columns existed.
ALTER TABLE encrypted_content_keys ADD COLUMN IF NOT EXISTS lockbox_salt BYTEA;
ALTER TABLE encrypted_content_keys ADD COLUMN IF NOT EXISTS patient_public_key BYTEA;
