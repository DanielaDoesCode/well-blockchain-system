//go:build integration

// Integration test for the tamper-evident ledger: inserts rows through the
// wrapped, hash-chained inserts, seals batches, anchors them on anvil via
// LogAnchor, then verifies the full ledger — including tamper detection.
//
// Requires anvil on :8545 and postgres (well_db) as the other integration
// tests do. Run from well-backend:
//
//	go test -tags integration ./internal/ledger/
package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"Thesis/well-blockchain/well-backend/internal/db"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func dialAnvil(url string) (*ethclient.Client, error) {
	var lastErr error
	for i := 0; i < 30; i++ {
		client, err := ethclient.Dial(url)
		if err == nil {
			if _, err := client.ChainID(context.Background()); err == nil {
				return client, nil
			}
			client.Close()
			lastErr = err
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("anvil not ready after 30s: %w", lastErr)
}

func TestLedgerIntegration(t *testing.T) {
	anvilURL := "http://localhost:8545"
	dbURL := "postgres://postgres:postgres@localhost:5432/well_db?sslmode=disable"
	if v := os.Getenv("LEDGER_TEST_ANVIL"); v != "" {
		anvilURL = v
	}

	ctx := context.Background()
	client, err := dialAnvil(anvilURL)
	if err != nil {
		t.Fatalf("anvil: %v", err)
	}
	defer client.Close()

	database, err := db.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer database.Close()

	schema, err := os.ReadFile("../../../db/init.sql")
	if err != nil {
		t.Fatalf("read init.sql: %v", err)
	}
	if _, err := database.Pool.Exec(ctx, string(schema)); err != nil {
		t.Fatalf("exec init.sql: %v", err)
	}
	if _, err := database.Pool.Exec(ctx,
		"TRUNCATE event_logs, access_logs, batched_logs RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("clean ledger tables: %v", err)
	}

	// ── deploy LogAnchor ──
	chainID, _ := client.ChainID(ctx)
	keyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	deployerKey, _ := crypto.HexToECDSA(keyHex)
	auth, _ := bind.NewKeyedTransactorWithChainID(deployerKey, chainID)
	auth.GasLimit = 5_000_000
	auth.GasPrice = big.NewInt(1_000_000_000)

	raw, err := os.ReadFile("../../../out/LogAnchor.sol/LogAnchor.json")
	if err != nil {
		t.Fatalf("read LogAnchor artifact (run `forge build`): %v", err)
	}
	var art struct {
		ABI      json.RawMessage `json:"abi"`
		Bytecode struct {
			Object string `json:"object"`
		} `json:"bytecode"`
	}
	if err := json.Unmarshal(raw, &art); err != nil {
		t.Fatalf("unmarshal artifact: %v", err)
	}
	parsed, err := abi.JSON(strings.NewReader(string(art.ABI)))
	if err != nil {
		t.Fatalf("parse ABI: %v", err)
	}
	anchorAddr, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(art.Bytecode.Object), client)
	if err != nil {
		t.Fatalf("deploy LogAnchor: %v", err)
	}
	if _, err := bind.WaitMined(ctx, client, tx); err != nil {
		t.Fatalf("wait LogAnchor deploy: %v", err)
	}
	t.Logf("deployed LogAnchor → %s", anchorAddr.Hex())

	// ── ledger (anchoring operator) ──
	lgr, err := New(ctx, anvilURL, database, anchorAddr.Hex(), keyHex, 100, 30*time.Second)
	if err != nil {
		t.Fatalf("new ledger: %v", err)
	}
	defer lgr.Close()

	// Verify-only ledger (no operator key) used later for tamper checks.
	verifyOnly, err := New(ctx, anvilURL, database, anchorAddr.Hex(), "", 0, 0)
	if err != nil {
		t.Fatalf("new verify-only ledger: %v", err)
	}
	defer verifyOnly.Close()

	// ── seed rows through the chained inserts ──
	for i := 0; i < 5; i++ {
		if err := database.InsertAccessLog(ctx, db.AccessLog{
			AccessorAddress: "0x1111111111111111111111111111111111111111",
			PatientAddress:  fmt.Sprintf("0x222222222222222222222222222222222222222%d", i),
			RecordID:        fmt.Sprintf("record-%d", i),
			AccessType:      "VIEW",
			PolicyMatched:   "allowed",
			WasEmergency:    false,
			TxHash:          "0x" + strings.Repeat(fmt.Sprintf("%x", i+1), 64),
		}); err != nil {
			t.Fatalf("insert access log %d: %v", i, err)
		}
	}
	for i := 0; i < 3; i++ {
		if err := database.InsertEventLog(ctx, db.EventLog{
			EventType:       "AccessGranted",
			ContractAddress: "0x3333333333333333333333333333333333333333",
			TxHash:          "0x" + strings.Repeat(fmt.Sprintf("%x", i+10), 64),
			BlockNumber:     1000 + uint64(i),
			Data:            map[string]any{"subject": fmt.Sprintf("id-%d", i), "note": "import"},
		}); err != nil {
			t.Fatalf("insert event log %d: %v", i, err)
		}
	}

	// ── seal + anchor a batch per table ──
	lgr.SealPending(ctx)

	full, err := verifyOnly.VerifyLedger(ctx)
	if err != nil {
		t.Fatalf("verify failed with error: %v", err)
	}
	if !full.Pass {
		t.Fatalf("expected clean ledger, got issues: %v", full.Issues)
	}
	if full.RowsChecked != 8 || full.BatchesChecked != 2 || full.AnchorsChecked != 2 {
		t.Fatalf("unexpected counts: %+v", full)
	}

	// every row verifies individually too
	for seq := int64(1); seq <= 5; seq++ {
		rep, err := verifyOnly.VerifySingle(ctx, db.TableAccessLogs, seq)
		if err != nil {
			t.Fatalf("verify single %s:%d: %v", db.TableAccessLogs, seq, err)
		}
		if !rep.Pass {
			t.Fatalf("verify single %s:%d: %+v", db.TableAccessLogs, seq, rep.Issues)
		}
	}
	for seq := int64(1); seq <= 3; seq++ {
		rep, err := verifyOnly.VerifySingle(ctx, db.TableEventLogs, seq)
		if err != nil {
			t.Fatalf("verify single %s:%d: %v", db.TableEventLogs, seq, err)
		}
		if !rep.Pass {
			t.Fatalf("verify single %s:%d: %+v", db.TableEventLogs, seq, rep.Issues)
		}
	}

	// ── extend the chain with a second batch and re-verify ──
	for i := 5; i < 7; i++ {
		if err := database.InsertAccessLog(ctx, db.AccessLog{
			AccessorAddress: "0x1111111111111111111111111111111111111111",
			PatientAddress:  fmt.Sprintf("0x222222222222222222222222222222222222222%d", i),
			RecordID:        fmt.Sprintf("record-%d", i),
			AccessType:      "VIEW",
			PolicyMatched:   "allowed",
			WasEmergency:    false,
			TxHash:          "0x" + strings.Repeat(fmt.Sprintf("%x", 9), 64),
		}); err != nil {
			t.Fatalf("insert access log %d: %v", i, err)
		}
	}
	lgr.SealPending(ctx)
	full, err = verifyOnly.VerifyLedger(ctx)
	if err != nil {
		t.Fatalf("verify failed with error: %v", err)
	}
	if !full.Pass {
		t.Fatalf("expected clean ledger after second batch: %v", full.Issues)
	}

	// ── tamper detection: mutate a stored row ──
	if _, err := database.Pool.Exec(ctx,
		`UPDATE access_logs SET access_type = 'WRITE_SENSITIVE' WHERE seq = 2`); err != nil {
		t.Fatalf("tamper: %v", err)
	}
	full, err = verifyOnly.VerifyLedger(ctx)
	if err != nil {
		t.Fatalf("verify after tamper: %v", err)
	}
	if full.Pass {
		t.Fatal("expected tamper detection, ledger still verified clean")
	}
	found := false
	for _, issue := range full.Issues {
		if strings.Contains(issue, "row 2") && strings.Contains(issue, "row_hash mismatch") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a row-level hash mismatch issue, got: %v", full.Issues)
	}

	// single-row verify of the tampered row must also fail
	single, err := verifyOnly.VerifySingle(ctx, db.TableAccessLogs, 2)
	if err != nil {
		t.Fatalf("verify single tampered: %v", err)
	}
	if single.Pass {
		t.Fatalf("expected tampered row verify to fail: %+v", single)
	}

	// ── restore → ledger clean again ──
	if _, err := database.Pool.Exec(ctx,
		`UPDATE access_logs SET access_type = 'VIEW' WHERE seq = 2`); err != nil {
		t.Fatalf("restore: %v", err)
	}
	full, err = verifyOnly.VerifyLedger(ctx)
	if err != nil {
		t.Fatalf("verify after restore: %v", err)
	}
	if !full.Pass {
		t.Fatalf("expected clean ledger after restore: %v", full.Issues)
	}

	// ── tamper detection: mid-chain row deletion ──
	if _, err := database.Pool.Exec(ctx,
		`DELETE FROM event_logs WHERE seq = 2`); err != nil {
		t.Fatalf("delete row: %v", err)
	}
	full, err = verifyOnly.VerifyLedger(ctx)
	if err != nil {
		t.Fatalf("verify after deletion: %v", err)
	}
	if full.Pass {
		t.Fatal("expected deletion to be detected, ledger still verified clean")
	}

	t.Log("ledger integration test passed (insert → seal → anchor → verify → tamper detected)")
}