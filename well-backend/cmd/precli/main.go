// Command precli is a small reference PRE client used by the end-to-end
// scripts. It performs the operations that a real client must do itself and
// that the backend deliberately cannot: generating PRE key pairs, sealing a
// lockbox, deriving a re-encryption key, and opening a re-encrypted ciphertext.
//
// It is not part of the server. The state file it writes contains secret keys
// and must never be uploaded anywhere.
//
// Usage:
//
//	precli setup  --state FILE
//	precli self   --state FILE --lockbox HEX
//	precli verify --state FILE --ct HEX
//	precli bench  [--reps N] [--verify N]
//
// `bench` times the individual AFGH operations in-process (no HTTP, no
// database, no chain) and prints a JSON report on stdout: latency
// distributions per operation, wire sizes, and round-trip correctness counts.
// It is the primitive-level measurement that the HTTP benchmarks and the
// Locust load runs cannot separate out from their own stack.
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"Thesis/well-blockchain/well-backend/internal/pre"
)

type clientState struct {
	PatientSecret string `json:"patient_secret"`
	DoctorSecret  string `json:"doctor_secret"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: precli <setup|self|verify|bench> [flags]")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "setup":
		setup(os.Args[2:])
	case "self":
		selfDecrypt(os.Args[2:])
	case "verify":
		verify(os.Args[2:])
	case "bench":
		bench(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func setup(args []string) {
	fs := flag.NewFlagSet("setup", flag.ExitOnError)
	statePath := fs.String("state", "", "path to write the client key state (contains secrets)")
	fs.Parse(args)
	if *statePath == "" {
		fatal("setup: --state is required")
	}

	patientPK, patientSK, err := pre.KeyGen()
	if err != nil {
		fatal("keygen patient: %v", err)
	}
	doctorPK, doctorSK, err := pre.KeyGen()
	if err != nil {
		fatal("keygen doctor: %v", err)
	}

	ct, err := pre.SealRandom(patientPK)
	if err != nil {
		fatal("seal: %v", err)
	}
	salt, err := pre.NewSalt()
	if err != nil {
		fatal("salt: %v", err)
	}
	rk, err := pre.DeriveReEncryptionKey(patientSK, doctorPK)
	if err != nil {
		fatal("derive rk: %v", err)
	}

	lockbox, err := ct.MarshalBinary()
	if err != nil {
		fatal("marshal lockbox: %v", err)
	}
	pk, err := patientPK.MarshalBinary()
	if err != nil {
		fatal("marshal public key: %v", err)
	}
	rkBytes, err := rk.MarshalBinary()
	if err != nil {
		fatal("marshal rk: %v", err)
	}

	state := clientState{
		PatientSecret: hex.EncodeToString(patientSK.Bytes()),
		DoctorSecret:  hex.EncodeToString(doctorSK.Bytes()),
	}
	blob, err := json.Marshal(state)
	if err != nil {
		fatal("marshal state: %v", err)
	}
	if err := os.WriteFile(*statePath, blob, 0o600); err != nil {
		fatal("write state: %v", err)
	}

	emit(map[string]string{
		"lockbox":            hex.EncodeToString(lockbox),
		"lockbox_salt":       hex.EncodeToString(salt),
		"patient_public_key": hex.EncodeToString(pk),
		"re_encryption_key":  hex.EncodeToString(rkBytes),
	})
}

func selfDecrypt(args []string) {
	fs := flag.NewFlagSet("self", flag.ExitOnError)
	statePath := fs.String("state", "", "path to the client key state")
	lockboxHex := fs.String("lockbox", "", "lockbox (delegator ciphertext) hex")
	fs.Parse(args)

	state := loadState(*statePath)
	sk, err := pre.ParseSecretKey(mustHex(state.PatientSecret, "state.patient_secret"))
	if err != nil {
		fatal("parse patient secret: %v", err)
	}
	lockbox, err := pre.ParseDelegatorCiphertext(mustHex(*lockboxHex, "--lockbox"))
	if err != nil {
		fatal("parse lockbox: %v", err)
	}
	msg, err := pre.DelegateeSelfDecrypt(lockbox, sk)
	if err != nil {
		fatal("self decrypt: %v", err)
	}
	msgBytes, err := msg.MarshalBinary()
	if err != nil {
		fatal("marshal message: %v", err)
	}
	emit(map[string]string{"message": hex.EncodeToString(msgBytes)})
}

func verify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	statePath := fs.String("state", "", "path to the client key state")
	ctHex := fs.String("ct", "", "re-encrypted ciphertext hex")
	fs.Parse(args)

	state := loadState(*statePath)
	sk, err := pre.ParseSecretKey(mustHex(state.DoctorSecret, "state.doctor_secret"))
	if err != nil {
		fatal("parse doctor secret: %v", err)
	}
	ct, err := pre.ParseReencryptedCiphertext(mustHex(*ctHex, "--ct"))
	if err != nil {
		fatal("parse ciphertext: %v", err)
	}
	msg, err := pre.DelegateeDecrypt(ct, sk)
	if err != nil {
		fatal("delegatee decrypt: %v", err)
	}
	msgBytes, err := msg.MarshalBinary()
	if err != nil {
		fatal("marshal message: %v", err)
	}
	emit(map[string]string{"message": hex.EncodeToString(msgBytes)})
}

// opStat is a latency distribution in milliseconds for one operation.
type opStat struct {
	N    int     `json:"n"`
	Min  float64 `json:"minMs"`
	P50  float64 `json:"p50Ms"`
	P95  float64 `json:"p95Ms"`
	P99  float64 `json:"p99Ms"`
	Max  float64 `json:"maxMs"`
	Mean float64 `json:"meanMs"`
}

// bench times the AFGH primitives in-process. Everything here is pure CPU on
// BLS12-381: no HTTP, no Postgres, no chain. The numbers isolate the
// cryptographic cost that the HTTP and load benchmarks fold in with everything
// else.
func bench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	reps := fs.Int("reps", 1000, "timed iterations per operation")
	verifyReps := fs.Int("verify", 100, "round-trip correctness checks")
	fs.Parse(args)
	if *reps < 1 {
		fatal("bench: --reps must be >= 1")
	}

	patientPK, patientSK, err := pre.KeyGen()
	if err != nil {
		fatal("keygen patient: %v", err)
	}
	doctorPK, doctorSK, err := pre.KeyGen()
	if err != nil {
		fatal("keygen doctor: %v", err)
	}
	ct, err := pre.SealRandom(patientPK)
	if err != nil {
		fatal("seal: %v", err)
	}
	rk, err := pre.DeriveReEncryptionKey(patientSK, doctorPK)
	if err != nil {
		fatal("derive rk: %v", err)
	}
	ct2, err := pre.Reencrypt(ct, rk)
	if err != nil {
		fatal("reencrypt: %v", err)
	}

	lockbox, err := ct.MarshalBinary()
	if err != nil {
		fatal("marshal lockbox: %v", err)
	}
	pk, err := patientPK.MarshalBinary()
	if err != nil {
		fatal("marshal public key: %v", err)
	}
	rkBytes, err := rk.MarshalBinary()
	if err != nil {
		fatal("marshal rk: %v", err)
	}
	ct2Bytes, err := ct2.MarshalBinary()
	if err != nil {
		fatal("marshal reencrypted: %v", err)
	}
	salt, err := pre.NewSalt()
	if err != nil {
		fatal("salt: %v", err)
	}

	// The request-path simulation: hex decode + parse both inputs, transform,
	// marshal + hex encode the result. This is what handleReEncrypt actually
	// spends, and what the `record.reencrypt` HTTP metric wraps around.
	lockboxHex := hex.EncodeToString(lockbox)
	rkHex := hex.EncodeToString(rkBytes)

	ops := map[string]opStat{}
	run := func(name string, fn func() error) {
		ops[name] = timed(name, *reps, fn)
	}

	run("keygen", func() error {
		_, _, err := pre.KeyGen()
		return err
	})
	run("seal", func() error {
		_, err := pre.SealRandom(patientPK)
		return err
	})
	run("deriveRk", func() error {
		_, err := pre.DeriveReEncryptionKey(patientSK, doctorPK)
		return err
	})
	run("reencrypt", func() error {
		_, err := pre.Reencrypt(ct, rk)
		return err
	})
	run("reencrypt_endpoint", func() error {
		c, err := pre.ParseDelegatorCiphertext(mustHex(lockboxHex, "lockbox"))
		if err != nil {
			return err
		}
		k, err := pre.ParseReEncryptionKey(mustHex(rkHex, "rk"))
		if err != nil {
			return err
		}
		out, err := pre.Reencrypt(c, k)
		if err != nil {
			return err
		}
		b, err := out.MarshalBinary()
		if err != nil {
			return err
		}
		_ = hex.EncodeToString(b)
		return nil
	})
	run("delegateeDecrypt", func() error {
		_, err := pre.DelegateeDecrypt(ct2, doctorSK)
		return err
	})
	run("selfDecrypt", func() error {
		_, err := pre.DelegateeSelfDecrypt(ct, patientSK)
		return err
	})

	// Correctness: seal -> reencrypt -> delegatee decrypt must equal the
	// patient's own self-decrypt of the same lockbox (both equal the sealed
	// message), for --verify fresh ciphertexts.
	failures := 0
	nVerify := *verifyReps
	for i := 0; i < nVerify; i++ {
		c, err := pre.SealRandom(patientPK)
		if err != nil {
			failures++
			continue
		}
		r, err := pre.Reencrypt(c, rk)
		if err != nil {
			failures++
			continue
		}
		viaDelegatee, err := pre.DelegateeDecrypt(r, doctorSK)
		if err != nil {
			failures++
			continue
		}
		viaSelf, err := pre.DelegateeSelfDecrypt(c, patientSK)
		if err != nil {
			failures++
			continue
		}
		a, err := viaDelegatee.MarshalBinary()
		if err != nil {
			failures++
			continue
		}
		b, err := viaSelf.MarshalBinary()
		if err != nil {
			failures++
			continue
		}
		if string(a) != string(b) {
			failures++
		}
	}

	emit(map[string]any{
		"reps":   *reps,
		"verify": nVerify,
		"ops":    ops,
		"sizes": map[string]int{
			"lockbox_bytes":            len(lockbox),
			"re_encryption_key_bytes":  len(rkBytes),
			"reencrypted_bytes":        len(ct2Bytes),
			"patient_public_key_bytes": len(pk),
			"lockbox_salt_bytes":       len(salt),
		},
		"correctness": map[string]int{"roundtrips": nVerify, "failures": failures},
	})
	if failures > 0 {
		fatal("bench: %d/%d round-trips failed", failures, nVerify)
	}
}

// timed runs fn n times after a short warm-up and summarises the durations.
func timed(name string, n int, fn func() error) opStat {
	warm := n
	if warm > 10 {
		warm = 10
	}
	for i := 0; i < warm; i++ {
		if err := fn(); err != nil {
			fatal("bench %s (warm-up): %v", name, err)
		}
	}
	durs := make([]float64, 0, n)
	var sum float64
	for i := 0; i < n; i++ {
		t0 := time.Now()
		if err := fn(); err != nil {
			fatal("bench %s: %v", name, err)
		}
		us := float64(time.Since(t0).Microseconds()) / 1000.0
		durs = append(durs, us)
		sum += us
	}
	sort.Float64s(durs)
	return opStat{
		N:    n,
		Min:  durs[0],
		P50:  pct(durs, 50),
		P95:  pct(durs, 95),
		P99:  pct(durs, 99),
		Max:  durs[len(durs)-1],
		Mean: sum / float64(n),
	}
}

func pct(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted))*p/100.0 + 0.5)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func loadState(path string) clientState {
	if path == "" {
		fatal("--state is required")
	}
	blob, err := os.ReadFile(path)
	if err != nil {
		fatal("read state: %v", err)
	}
	var state clientState
	if err := json.Unmarshal(blob, &state); err != nil {
		fatal("parse state: %v", err)
	}
	return state
}

func mustHex(s, what string) []byte {
	if s == "" {
		fatal("%s is empty", what)
	}
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		fatal("%s is not hex: %v", what, err)
	}
	return b
}

func emit(v any) {
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(v); err != nil {
		fatal("write output: %v", err)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
