// Command logverify audits the tamper-evident audit ledger: it recomputes row
// hashes, batch Merkle roots and chain hashes from the database and compares
// them against the chain hashes anchored on the EVM chain via LogAnchor.
//
// Usage:
//
//	logverify single --table access_logs|event_logs --seq N
//	logverify full
//
// It reads DB_URL / ETH_RPC / LOG_ANCHOR_ADDRESS from the environment.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/ledger"
)

func main() {
	log.SetFlags(0)

	sub := ""
	if len(os.Args) > 1 {
		sub = os.Args[1]
	}
	fs := flag.NewFlagSet("logverify "+sub, flag.ExitOnError)
	table := fs.String("table", "", "one of access_logs, event_logs")
	seq := fs.Int64("seq", 0, "row sequence number to verify")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "logverify verifies the tamper-evident audit ledger")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  logverify single --table access_logs|event_logs --seq N")
		fmt.Fprintln(os.Stderr, "  logverify full")
		os.Exit(2)
	}

	env := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("missing %s", k)
		}
		return v
	}

	ctx := context.Background()

	switch sub {
	case "single":
		if err := fs.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		runSingle(ctx, env, *table, *seq)
	case "full":
		runFull(ctx, env)
	default:
		fs.Usage()
	}
}

func runSingle(ctx context.Context, env func(string) string, table string, seq int64) {
	if table == "" || seq <= 0 {
		log.Fatalf("single requires --table and --seq")
	}

	database, err := db.New(ctx, env("DB_URL"))
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	l, err := ledger.New(ctx, env("ETH_RPC"), database, env("LOG_ANCHOR_ADDRESS"), "", 0, 0)
	if err != nil {
		log.Fatalf("ledger: %v", err)
	}
	defer l.Close()

	rep, err := l.VerifySingle(ctx, table, seq)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	out, _ := json.MarshalIndent(rep, "", "  ")
	fmt.Println(string(out))
	if !rep.Pass {
		os.Exit(1)
	}
}

func runFull(ctx context.Context, env func(string) string) {
	database, err := db.New(ctx, env("DB_URL"))
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	l, err := ledger.New(ctx, env("ETH_RPC"), database, env("LOG_ANCHOR_ADDRESS"), "", 0, 0)
	if err != nil {
		log.Fatalf("ledger: %v", err)
	}
	defer l.Close()

	rep, err := l.VerifyLedger(ctx)
	if err != nil {
		log.Fatalf("verify: %v", err)
	}
	out, _ := json.MarshalIndent(rep, "", "  ")
	fmt.Println(string(out))
	if !rep.Pass {
		os.Exit(1)
	}
	fmt.Println("OK: ledger verifies against the on-chain anchors")
	os.Exit(0)
}