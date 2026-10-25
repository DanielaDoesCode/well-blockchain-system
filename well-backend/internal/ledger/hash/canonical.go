package hash

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"time"
)

// Domain prefixes separate the two ledgers and version the rows.
const (
	DomainAccessLogs = "well.access_logs.v1"
	DomainEventLogs  = "well.event_logs.v1"
)

// Zero is the field zero used for the genesis chain state and empty Merkle
// leaves (matching AttestationManager's zero-value leaves).
var Zero = new(big.Int)

// StrField encodes a string as a field element (bytes of the string, reduced
// mod p; empty string is zero).
func StrField(s string) *big.Int {
	return BytesToField([]byte(s))
}

// HexField decodes a 0x-prefixed hex string (address/hash) into a field element.
// CHAR(n) columns right-pad short values with spaces, so whitespace is trimmed
// first and odd-length hex is left-padded with a nibble — this makes the field
// deterministic whether it is hashed before (raw input) or after (padded read-back).
func HexField(s string) *big.Int {
	s = strings.TrimSpace(trim0x(s))
	if s == "" {
		return Zero
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		// Not valid hex (e.g. record_id like "record-0"): fall back to hashing
		// the trimmed literal characters so insert and read-back stay identical.
		return StrField(s)
	}
	return BytesToField(b)
}

// UintField encodes a uint64 as a field element.
func UintField(u uint64) *big.Int {
	return new(big.Int).SetUint64(u)
}

// BoolField encodes a boolean as 0 or 1.
func BoolField(b bool) *big.Int {
	if b {
		return big.NewInt(1)
	}
	return Zero
}

// timeToMicros renders a timestamp as whole microseconds since the Unix epoch,
// truncated to the microsecond precision that PostgreSQL TIMESTAMP retains.
func timeToMicros(t time.Time) int64 {
	t = t.UTC().Truncate(time.Microsecond)
	return t.Unix()*1_000_000 + int64(t.Nanosecond())/1000
}

// TimeField encodes a timestamp (microsecond precision) as a field element.
func TimeField(t time.Time) *big.Int {
	return new(big.Int).SetInt64(timeToMicros(t))
}

// JSONBytesField hashes canonical JSON bytes. The input is normalised first
// (key-sorted, no whitespace, numbers preserved verbatim via json.Number) so
// the result is stable regardless of how PostgreSQL emits the JSONB text.
func JSONBytesField(raw []byte) *big.Int {
	return foldBytesChunks(canonicalJSONBytes(raw))
}

// DataField hashes a Go data value the same way JSONBytesField hashes its
// marshalled form (used at insert time, before the value reaches PostgreSQL).
func DataField(v any) *big.Int {
	b, err := json.Marshal(v)
	if err != nil {
		b = []byte("{}")
	}
	return JSONBytesField(b)
}

// canonicalJSONBytes parses JSON and re-marshals it deterministically:
// map keys sorted, no insignificant whitespace, numbers preserved as-is.
func canonicalJSONBytes(raw []byte) []byte {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		// Unparsable JSON never reaches the DB (it is jsonb); hash raw bytes.
		return raw
	}
	out, err := json.Marshal(v)
	if err != nil {
		return raw
	}
	return out
}

// foldBytesChunks folds byte data into a single field element, chunked into
// ≤31-byte pieces so every chunk is < the field prime.
func foldBytesChunks(b []byte) *big.Int {
	if len(b) == 0 {
		return Zero
	}
	var vals []*big.Int
	for off := 0; off < len(b); off += 31 {
		end := off + 31
		if end > len(b) {
			end = len(b)
		}
		vals = append(vals, BytesToField(b[off:end]))
	}
	return HashFold(vals...)
}

func trim0x(s string) string {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return s[2:]
	}
	return s
}

// AccessLogContentHash hashes the content of an access_logs row (excluding the
// self-referential prev_hash/row_hash/seq columns), in a fixed column order.
func AccessLogContentHash(
	accessedAt time.Time,
	accessorAddress, patientAddress, recordID, accessType, policyMatched string,
	wasEmergency bool,
	txHash string,
) *big.Int {
	return HashFold(
		new(big.Int).Set(Zero),
		StrField(DomainAccessLogs),
		TimeField(accessedAt),
		HexField(accessorAddress),
		HexField(patientAddress),
		HexField(recordID),
		StrField(accessType),
		StrField(policyMatched),
		BoolField(wasEmergency),
		HexField(txHash),
	)
}

// EventLogContentHash hashes the content of an event_logs row (excluding the
// self-referential columns), in a fixed column order. data is the JSONB value
// as text read back from PostgreSQL (or the marshalled Go value).
func EventLogContentHash(
	eventType, contractAddress, txHash string,
	blockNumber uint64,
	data []byte,
	loggedAt time.Time,
) *big.Int {
	return HashFold(
		new(big.Int).Set(Zero),
		StrField(DomainEventLogs),
		StrField(eventType),
		HexField(contractAddress),
		HexField(txHash),
		UintField(blockNumber),
		JSONBytesField(data),
		TimeField(loggedAt),
	)
}

// RowHash chains a row onto its predecessor: Poseidon(prevHash, contentHash).
// prev is the previous row hash as 32-byte big-endian (nil ⇒ genesis, zero).
func RowHash(prev []byte, contentHash *big.Int) *big.Int {
	prevF := Zero
	if len(prev) > 0 {
		prevF = BytesToField(prev)
	}
	return Poseidon2(prevF, contentHash)
}