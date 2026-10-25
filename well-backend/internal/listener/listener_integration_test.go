//go:build integration

package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/db"
)

func dialAnvil(url string) (*ethclient.Client, error) {
	var lastErr error
	for i := 0; i < 30; i++ {
		client, err := ethclient.Dial(url)
		if err == nil {
			_, err := client.ChainID(context.Background())
			if err == nil {
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

func dialDB(url string) (*db.DB, error) {
	var lastErr error
	for i := 0; i < 30; i++ {
		d, err := db.New(context.Background(), url)
		if err == nil {
			return d, nil
		}
		lastErr = err
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("postgres not ready after 30s: %w", lastErr)
}

type compiledArtifact struct {
	ABI      json.RawMessage `json:"abi"`
	Bytecode struct {
		Object string `json:"object"`
	} `json:"bytecode"`
}

func deployContract(t *testing.T, auth *bind.TransactOpts, client *ethclient.Client, artifactPath string, args ...any) common.Address {
	t.Helper()

	raw, err := os.ReadFile(artifactPath)
	if err != nil {
		t.Fatalf("read artifact %s: %v", artifactPath, err)
	}

	var art compiledArtifact
	if err := json.Unmarshal(raw, &art); err != nil {
		t.Fatalf("unmarshal artifact: %v", err)
	}

	parsed, err := abi.JSON(strings.NewReader(string(art.ABI)))
	if err != nil {
		t.Fatalf("parse ABI: %v", err)
	}

	addr, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(art.Bytecode.Object), client, args...)
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}

	if _, err := bind.WaitMined(context.Background(), client, tx); err != nil {
		t.Fatalf("wait mined: %v", err)
	}

	t.Logf("deployed %s → %s", artifactPath, addr.Hex())
	return addr
}

func checkEventLog(t *testing.T, database *db.DB, eventType string) bool {
	t.Helper()
	var count int
	err := database.Pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM event_logs WHERE event_type = $1", eventType).Scan(&count)
	if err != nil {
		t.Logf("query event_logs: %v", err)
		return false
	}
	return count > 0
}

func checkAccessLog(t *testing.T, database *db.DB) bool {
	t.Helper()
	var count int
	err := database.Pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM access_logs").Scan(&count)
	if err != nil {
		t.Logf("query access_logs: %v", err)
		return false
	}
	return count > 0
}

func TestListenerIntegration(t *testing.T) {
	anvilHTTPURL := "http://localhost:8545"
	anvilWSURL := "ws://localhost:8545"
	dbURL := "postgres://postgres:postgres@localhost:5432/well_db?sslmode=disable"

	t.Log("connecting to Hardhat...")
	client, err := dialAnvil(anvilHTTPURL)
	if err != nil {
		t.Fatalf("hardhat: %v\n  ensure `docker compose up -d` is running", err)
	}
	defer client.Close()

	t.Log("connecting to Postgres...")
	database, err := dialDB(dbURL)
	if err != nil {
		t.Fatalf("postgres: %v\n  ensure `docker compose up -d` is running", err)
	}
	defer database.Close()

	// initialize schema
	schema, err := os.ReadFile("../../../db/init.sql")
	if err != nil {
		t.Fatalf("read init.sql: %v", err)
	}
	if _, err := database.Pool.Exec(context.Background(), string(schema)); err != nil {
		t.Fatalf("exec init.sql: %v", err)
	}

	_, _ = database.Pool.Exec(context.Background(), "DELETE FROM event_logs")
	_, _ = database.Pool.Exec(context.Background(), "DELETE FROM access_logs")

	chainID := big.NewInt(31337)
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		t.Fatalf("create auth: %v", err)
	}
	auth.GasLimit = 8_000_000
	auth.GasPrice = big.NewInt(1_000_000_000)

	//deploy contracts

	// Deploy Groth16Verifier first (required by ZKAuthVerifier)
	groth16VerifierAddr := deployContract(t, auth, client,
		"../../../out/Groth16Verifier.sol/Groth16Verifier.json",
	)

	zkVerifierAddr := deployContract(t, auth, client,
		"../../../out/ZKAuthVerifier.sol/ZKAuthVerifier.json",
		groth16VerifierAddr,
		big.NewInt(3600),
	)

	userRegistryAddr := deployContract(t, auth, client,
		"../../../out/UserRegistry.sol/UserRegistry.json",
	)

	historyFactoryAddr := deployContract(t, auth, client,
		"../../../out/HistoryFactory.sol/HistoryFactory.json",
		userRegistryAddr,
	)

	registryAddr := deployContract(t, auth, client,
		"../../../out/Registry.sol/WellRegistry.json",
		userRegistryAddr,
		historyFactoryAddr,
	)

	// Wire up: set WellRegistry as the accessControl on HistoryFactory
	historyFactory, err := contracts.NewHistoryFactory(historyFactoryAddr, client)
	if err != nil {
		t.Fatalf("bind history factory for wiring: %v", err)
	}
	tx, err := historyFactory.SetWellRegistry(auth, registryAddr)
	if err != nil {
		t.Fatalf("set well registry: %v", err)
	}
	if _, err := bind.WaitMined(context.Background(), client, tx); err != nil {
		t.Fatalf("wait mined set well registry: %v", err)
	}

	//listener

	lst, err := New(context.Background(), anvilWSURL, database, nil,
		registryAddr.Hex(),
		userRegistryAddr.Hex(),
		historyFactoryAddr.Hex(),
		zkVerifierAddr.Hex(),
		"",
		"",
		"",
	)
	if err != nil {
		t.Fatalf("create listener: %v", err)
	}
	defer lst.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := lst.Start(ctx); err != nil {
		t.Fatalf("start listener: %v", err)
	}

	time.Sleep(2 * time.Second)

	//binds the contracts deployed

	registry, err := contracts.NewRegistry(registryAddr, client)
	if err != nil {
		t.Fatalf("bind registry: %v", err)
	}

	historyFactoryForParse, err := contracts.NewHistoryFactory(historyFactoryAddr, client)
	if err != nil {
		t.Fatalf("bind history factory for parse: %v", err)
	}

	// test1: register: registering doctor
	t.Run("UserRegistry.Registered", func(t *testing.T) {
		commitment := common.HexToHash("0xaa")
		internalID := common.HexToHash("0xbb")
		role := uint8(1) // Doctor

		tx, err := registry.Register(auth, commitment, internalID, role)
		if err != nil {
			t.Fatalf("register doctor: %v", err)
		}
		if _, err := bind.WaitMined(ctx, client, tx); err != nil {
			t.Fatalf("wait mined: %v", err)
		}

		time.Sleep(2 * time.Second)

		if !checkEventLog(t, database, "UserRegistry.Registered") {
			t.Error("UserRegistry.Registered event not found in event_logs")
		}
	})

	// ═══════════════════════════════════════════════════════════
	// Test 2: UserRegistry.Registered + HistoryCreated (Patient)
	// ═══════════════════════════════════════════════════════════

	t.Run("HistoryFactory.HistoryCreated", func(t *testing.T) {
		commitment := common.HexToHash("0xcc")
		internalID := common.HexToHash("0xdd")
		role := uint8(0) // Patient

		// Register as a different address to avoid "already registered"
		patientKey, err := crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
		if err != nil {
			t.Fatalf("parse patient key: %v", err)
		}
		patientAuth, err := bind.NewKeyedTransactorWithChainID(patientKey, chainID)
		if err != nil {
			t.Fatalf("create patient auth: %v", err)
		}
		patientAuth.GasLimit = 5_000_000
		patientAuth.GasPrice = big.NewInt(1_000_000_000)

		patientRegistry, err := contracts.NewRegistry(registryAddr, client)
		if err != nil {
			t.Fatalf("bind registry for patient: %v", err)
		}

		tx, err := patientRegistry.Register(patientAuth, commitment, internalID, role)
		if err != nil {
			t.Fatalf("register patient: %v", err)
		}
		if _, err := bind.WaitMined(ctx, client, tx); err != nil {
			t.Fatalf("wait mined: %v", err)
		}

		time.Sleep(2 * time.Second)

		if !checkEventLog(t, database, "UserRegistry.Registered") {
			t.Error("UserRegistry.Registered event not found for patient")
		}
		if !checkEventLog(t, database, "HistoryFactory.HistoryCreated") {
			t.Error("HistoryFactory.HistoryCreated event not found")
		}
	})

	// ═══════════════════════════════════════════════════════════
	// Test 3: UserRegistry.Revoked
	// ═══════════════════════════════════════════════════════════

	t.Run("UserRegistry.Revoked", func(t *testing.T) {
		commitment := common.HexToHash("0xaa")

		tx, err := registry.Revoke(auth, commitment)
		if err != nil {
			t.Fatalf("revoke: %v", err)
		}
		if _, err := bind.WaitMined(ctx, client, tx); err != nil {
			t.Fatalf("wait mined: %v", err)
		}

		time.Sleep(2 * time.Second)

		if !checkEventLog(t, database, "UserRegistry.Revoked") {
			t.Error("UserRegistry.Revoked event not found")
		}
	})

	// history.AccessGranted
	t.Run("History.AccessGranted", func(t *testing.T) {
		// Register a new patient (creates a History contract)
		patientCommitment := common.HexToHash("0xee")
		patientInternalID := common.HexToHash("0xff")

		patientKey, err := crypto.HexToECDSA("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a")
		if err != nil {
			t.Fatalf("parse patient key: %v", err)
		}
		patientAuth, err := bind.NewKeyedTransactorWithChainID(patientKey, chainID)
		if err != nil {
			t.Fatalf("create patient auth: %v", err)
		}
		patientAuth.GasLimit = 5_000_000
		patientAuth.GasPrice = big.NewInt(1_000_000_000)

		patientRegistry, err := contracts.NewRegistry(registryAddr, client)
		if err != nil {
			t.Fatalf("bind registry: %v", err)
		}

		tx, err := patientRegistry.Register(patientAuth, patientCommitment, patientInternalID, 0)
		if err != nil {
			t.Fatalf("register patient: %v", err)
		}
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			t.Fatalf("wait mined: %v", err)
		}

		// Extract the HistoryCreated event using HistoryFactory's ParseHistoryCreated
		var historyAddr common.Address
		for _, log := range receipt.Logs {
			parsed, err := historyFactoryForParse.ParseHistoryCreated(*log)
			if err == nil {
				historyAddr = parsed.HistoryContract
				break
			}
		}
		if historyAddr == (common.Address{}) {
			t.Fatal("could not find HistoryCreated event in receipt")
		}

		time.Sleep(2 * time.Second) // wait for poll + DB insert

		// Create a record via registry.createRecord
		recordMetadata := contracts.TypesRecordMetadata{
			RecordType: 0, // LabResult
			CreatedAt:  big.NewInt(0),
			UpdatedAt:  big.NewInt(0),
		}
		contentHash := common.HexToHash("0x11")
		storageRef := "ipfs://test"

		tx2, err := patientRegistry.CreateRecord(patientAuth,
			patientAuth.From,
			recordMetadata,
			contentHash,
			storageRef,
		)
		if err != nil {
			t.Fatalf("create record: %v", err)
		}
		receipt2, err := bind.WaitMined(ctx, client, tx2)
		if err != nil {
			t.Fatalf("wait mined record: %v", err)
		}

		// Extract Record address from RecordCreated event using abigen's ParseRecordCreated
		var recordAddr common.Address
		historyForParse, err := contracts.NewHistory(historyAddr, client)
		if err != nil {
			t.Fatalf("bind history for parse: %v", err)
		}
		for _, log := range receipt2.Logs {
			parsed, err := historyForParse.ParseRecordCreated(*log)
			if err == nil {
				recordAddr = parsed.Record
				break
			}
		}
		if recordAddr == (common.Address{}) {
			t.Fatal("could not find RecordCreated event in receipt")
		}

		time.Sleep(2 * time.Second) // wait for poll + DB insert

		// Grant access via History (patient calls grantRecordAccess)
		history, err := contracts.NewHistory(historyAddr, client)
		if err != nil {
			t.Fatalf("bind history: %v", err)
		}

		doctorAddr := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
		tx3, err := history.GrantRecordAccess(patientAuth, recordAddr, doctorAddr)
		if err != nil {
			t.Fatalf("grant access: %v", err)
		}
		if _, err := bind.WaitMined(ctx, client, tx3); err != nil {
			t.Fatalf("wait mined grant: %v", err)
		}

		time.Sleep(2 * time.Second)

		if !checkAccessLog(t, database) {
			t.Error("History.AccessGranted not found in access_logs")
		}
		if !checkEventLog(t, database, "Record.AccessGranted") {
			t.Error("Record.AccessGranted event not found")
		}
		if !checkEventLog(t, database, "History.AccessRevoked") {
			// revoke to also test that path
			tx4, err := history.RevokeRecordAccess(patientAuth, recordAddr, doctorAddr)
			if err != nil {
				t.Fatalf("revoke access: %v", err)
			}
			if _, err := bind.WaitMined(ctx, client, tx4); err != nil {
				t.Fatalf("wait mined revoke: %v", err)
			}
			time.Sleep(2 * time.Second)

			if !checkEventLog(t, database, "History.AccessRevoked") {
				t.Error("History.AccessRevoked event not found")
			}
			if !checkEventLog(t, database, "Record.AccessRevoked") {
				t.Error("Record.AccessRevoked event not found")
			}
		}
	})
}
