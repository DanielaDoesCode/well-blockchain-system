//go:build integration

package aggregator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/db"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type compiledArtifact struct {
	ABI      json.RawMessage `json:"abi"`
	Bytecode struct {
		Object string `json:"object"`
	} `json:"bytecode"`
}

func dialAnvil(url string) (*ethclient.Client, error) {
	var lastErr error
	for i := 0; i < 30; i++ {
		client, err := ethclient.Dial(url)
		if err == nil {
			if _, err := client.ChainID(context.Background()); err == nil {
				return client, nil
			} else {
				client.Close()
				lastErr = err
			}
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

type testProofData struct {
	PA0, PA1, PB00, PB01, PB10, PB11, PC0, PC1 string
	Commitment, Nullifier, ConditionGroup       string
}

func firstProofData() *testProofData {
	return &testProofData{
		PA0: "12008770116893986162536347804216010609604197142438772991629698146602677844413",
		PA1: "4115187077618860705872467988101494501958919832440735347232776070881127198440",
		PB00: "19627666940019467562725206405038585523179045705980990386483191186481471810455",
		PB01: "8587332350230284920725311018929398404255497048059423220408489920983937214125",
		PB10: "5701560962250042389648449452032487152614935793004559867509975120071596045388",
		PB11: "3803181176004050457297703739483606767218949402541643882049602686640781720194",
		PC0:  "4819162284881573788255921048766792259310465269016784493643700733804211153129",
		PC1:  "3452769194622951038303865641728458722225675992200018884311961933422807389369",
		Commitment:     "0x0fb849f7cf35865c838cef48782e803b2c38263e2f467799c87eff168eb4d897",
		Nullifier:      "0x2246d113bbdac3f6fc85512f40dedde2f923ac941795749e0e7136c65e5f1d38",
		ConditionGroup: "42",
	}
}

func secondProofData() *testProofData {
	// secret=555555555, nullifierSeed=444444444, conditionGroup=42
	return &testProofData{
		PA0: "9575277363530284711812354769683430112126859265527832401693191194431328609381",
		PA1: "12504049318146328848254147854465508002616777523581351660223727924708620955490",
		PB00: "5892163652893919714198600541993387352490148732441110622531403118805107131034",
		PB01: "4158080810435050212214455793402781330580551062858709785678514453636940201157",
		PB10: "3246460076262938763810996083706202558835138763274330438469545253783184483775",
		PB11: "18937191206519917088595321994642456509092586047707626873274515555563178262535",
		PC0:  "12992733662730462342110087332101724492658674584831511223077101340766048426622",
		PC1:  "8300648089261740850500475128880571795678385302124229875316334972401415696322",
		Commitment:     "0x17256f65a34a71a608b64eb582ca7c674ac52121bf1e48095476071f0b18d916",
		Nullifier:      "0x0306211c0083563122fb88ee5e82df1d095f85140c4e8791eefdbd1b7dfeac14",
		ConditionGroup: "42",
	}
}

func submitProof(t *testing.T, baseURL string, pd *testProofData) {
	t.Helper()
	body := map[string]any{
		"pA_0":            pd.PA0,
		"pA_1":            pd.PA1,
		"pB_0_0":          pd.PB00,
		"pB_0_1":          pd.PB01,
		"pB_1_0":          pd.PB10,
		"pB_1_1":          pd.PB11,
		"pC_0":            pd.PC0,
		"pC_1":            pd.PC1,
		"commitment":      pd.Commitment,
		"nullifier":       pd.Nullifier,
		"condition_group": pd.ConditionGroup,
	}
	b, _ := json.Marshal(body)
	resp, err := http.Post(baseURL+"/api/v1/contribute", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusOK || result["status"] != "accepted" {
		t.Fatalf("unexpected response: status=%d body=%v", resp.StatusCode, result)
	}
}

func TestAggregatorIntegration(t *testing.T) {
	anvilHTTPURL := "http://localhost:8545"
	anvilWSURL := "ws://localhost:8545"
	dbURL := "postgres://postgres:postgres@localhost:5432/well_db?sslmode=disable"

	t.Log("connecting to Anvil...")
	client, err := dialAnvil(anvilHTTPURL)
	if err != nil {
		t.Fatalf("anvil: %v", err)
	}
	defer client.Close()

	t.Log("connecting to Postgres...")
	database, err := dialDB(dbURL)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer database.Close()

	schema, err := os.ReadFile("../../../db/init.sql")
	if err != nil {
		t.Fatalf("read init.sql: %v", err)
	}
	if _, err := database.Pool.Exec(context.Background(), string(schema)); err != nil {
		t.Fatalf("exec init.sql: %v", err)
	}
	_, _ = database.Pool.Exec(context.Background(), "DELETE FROM pending_contributions")
	_, _ = database.Pool.Exec(context.Background(), "DELETE FROM research_counters")

	chainID := big.NewInt(31337)
	deployerKey, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	deployerAuth, _ := bind.NewKeyedTransactorWithChainID(deployerKey, chainID)
	deployerAuth.GasLimit = 8_000_000
	deployerAuth.GasPrice = big.NewInt(1_000_000_000)

	groth16Addr := deployContract(t, deployerAuth, client,
		"../../../out/Groth16Verifier.sol/Groth16Verifier.json")
	zKAuthAddr := deployContract(t, deployerAuth, client,
		"../../../out/ZKAuthVerifier.sol/ZKAuthVerifier.json", groth16Addr, big.NewInt(3600))
	_ = zKAuthAddr
	userRegistryAddr := deployContract(t, deployerAuth, client,
		"../../../out/UserRegistry.sol/UserRegistry.json")
	historyFactoryAddr := deployContract(t, deployerAuth, client,
		"../../../out/HistoryFactory.sol/HistoryFactory.json", userRegistryAddr)
	wellRegistryAddr := deployContract(t, deployerAuth, client,
		"../../../out/Registry.sol/WellRegistry.json", userRegistryAddr, historyFactoryAddr)

	hf, _ := contracts.NewHistoryFactory(historyFactoryAddr, client)
	tx, _ := hf.SetWellRegistry(deployerAuth, wellRegistryAddr)
	bind.WaitMined(context.Background(), client, tx)

	researchVerifierAddr := deployContract(t, deployerAuth, client,
		"../../../out/ResearchVerifier.sol/ResearchVerifier.json")
	researchContributerAddr := deployContract(t, deployerAuth, client,
		"../../../out/ResearchContributer.sol/ResearchContributer.json",
		researchVerifierAddr, wellRegistryAddr)

	// Register patient 1 with the first proof's commitment
	pd1 := firstProofData()
	registry, _ := contracts.NewRegistry(wellRegistryAddr, client)
	commit1 := common.HexToHash(pd1.Commitment)
	tx2, err := registry.Register(deployerAuth, commit1, common.HexToHash("0xdeadbeef01"), 0)
	if err != nil {
		t.Fatalf("register patient 1: %v", err)
	}
	bind.WaitMined(context.Background(), client, tx2)

	// Register patient 2 with the second proof's commitment (different address)
	pd2 := secondProofData()
	commit2 := common.HexToHash(pd2.Commitment)

	patient2Key, _ := crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	patient2Auth, _ := bind.NewKeyedTransactorWithChainID(patient2Key, chainID)
	patient2Auth.GasLimit = 8_000_000
	patient2Auth.GasPrice = big.NewInt(1_000_000_000)

	tx3, err := registry.Register(patient2Auth, commit2, common.HexToHash("0xdeadbeef02"), 0)
	if err != nil {
		t.Fatalf("register patient 2: %v", err)
	}
	bind.WaitMined(context.Background(), client, tx3)

	// Create aggregator (threshold=2)
	agg, err := New(context.Background(), anvilWSURL, database,
		"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		researchContributerAddr.Hex(), wellRegistryAddr.Hex(), 2)
	if err != nil {
		t.Fatalf("create aggregator: %v", err)
	}
	defer agg.Close()
	agg.Start(context.Background())

	// Start HTTP server on a fixed port
	handler := NewHandler(agg, database)
	srv := &http.Server{Addr: ":8099", Handler: handler}
	go srv.ListenAndServe()
	defer srv.Close()
	time.Sleep(500 * time.Millisecond)
	baseURL := "http://localhost:8099"

	// ── Test: unregistered commitment rejected ──
	t.Run("RejectUnregistered", func(t *testing.T) {
		bad := map[string]any{
			"commitment": "0x0000000000000000000000000000000000000000000000000000000000000001",
			"nullifier":  "0x0000000000000000000000000000000000000000000000000000000000000002",
			"pA_0": "1", "pA_1": "1",
			"pB_0_0": "1", "pB_0_1": "1", "pB_1_0": "1", "pB_1_1": "1",
			"pC_0": "1", "pC_1": "1",
			"condition_group": "99",
		}
		b, _ := json.Marshal(bad)
		resp, _ := http.Post(baseURL+"/api/v1/contribute", "application/json", bytes.NewReader(b))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	// ── Test: submit first contribution ──
	t.Run("SubmitFirst", func(t *testing.T) {
		submitProof(t, baseURL, pd1)
	})

	// ── Test: submit second contribution (triggers threshold) ──
	t.Run("SubmitSecond", func(t *testing.T) {
		submitProof(t, baseURL, pd2)
	})

	// ── Test: aggregator submits to chain when threshold reached ──
	t.Run("ChainSubmission", func(t *testing.T) {
		time.Sleep(15 * time.Second)
		var published int64
		err := database.Pool.QueryRow(context.Background(),
			"SELECT published_count FROM research_counters WHERE condition_group = '42'",
		).Scan(&published)
		if err != nil {
			t.Fatalf("query: %v", err)
		}
		if published < 2 {
			t.Fatalf("expected published_count >= 2, got %d", published)
		}
		t.Logf("published_count = %d", published)
	})

	// ── Test: query counters ──
	t.Run("QueryCounters", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/counters")
		if err != nil {
			t.Fatalf("GET: %v", err)
		}
		defer resp.Body.Close()
		var result struct {
			Counters []struct {
				ConditionGroup string `json:"condition_group"`
				TotalCount     int64  `json:"total_count"`
				PublishedCount int64  `json:"published_count"`
			} `json:"counters"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		if len(result.Counters) == 0 {
			t.Fatal("no counters returned")
		}
		for _, c := range result.Counters {
			t.Logf("counter: group=%s total=%d published=%d", c.ConditionGroup, c.TotalCount, c.PublishedCount)
			if c.ConditionGroup == "42" {
				if c.TotalCount < 2 || c.PublishedCount < 2 {
					t.Fatalf("expected >=2, got total=%d published=%d", c.TotalCount, c.PublishedCount)
				}
				return
			}
		}
		t.Fatal("counter for group 42 not found")
	})
}
