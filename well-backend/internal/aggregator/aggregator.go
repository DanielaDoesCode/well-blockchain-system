package aggregator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/db"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Aggregator struct {
	client     *ethclient.Client
	database   *db.DB
	auth       *bind.TransactOpts

	researchContributer *contracts.ResearchContributer
	wellRegistry        *contracts.Registry

	threshold    int64
	pollInterval time.Duration

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
}

func New(
	ctx context.Context,
	rpcURL string,
	database *db.DB,
	privateKeyHex string,
	researchContributerAddr string,
	wellRegistryAddr string,
	threshold int64,
) (*Aggregator, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum node: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("create transactor: %w", err)
	}
	auth.GasLimit = 5_000_000
	auth.GasPrice = big.NewInt(1_000_000_000)

	contributer, err := contracts.NewResearchContributer(common.HexToAddress(researchContributerAddr), client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("bind research contributer: %w", err)
	}

	registry, err := contracts.NewRegistry(common.HexToAddress(wellRegistryAddr), client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("bind registry: %w", err)
	}

	return &Aggregator{
		client:              client,
		database:            database,
		auth:                auth,
		researchContributer: contributer,
		wellRegistry:        registry,
		threshold:           threshold,
		pollInterval:        10 * time.Second,
		stopCh:              make(chan struct{}),
	}, nil
}

func (a *Aggregator) Start(ctx context.Context) {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return
	}
	a.running = true
	a.mu.Unlock()

	go a.pollLoop(ctx)
	log.Printf("aggregator started: threshold=%d, pollInterval=%v", a.threshold, a.pollInterval)
}

func (a *Aggregator) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		close(a.stopCh)
		a.running = false
	}
}

func (a *Aggregator) Close() {
	a.Stop()
	a.client.Close()
}

func (a *Aggregator) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.checkAndSubmit(ctx)
		case <-a.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (a *Aggregator) checkAndSubmit(ctx context.Context) {
	counters, err := a.database.GetCounters(ctx)
	if err != nil {
		log.Printf("aggregator: get counters: %v", err)
		return
	}

	for _, c := range counters {
		pending := c.TotalCount - c.PublishedCount
		if pending < a.threshold {
			continue
		}

		log.Printf("aggregator: condition_group=%s has %d pending (threshold=%d), submitting...",
			c.ConditionGroup, pending, a.threshold)

		contribs, err := a.database.GetPendingContributions(ctx, c.ConditionGroup)
		if err != nil {
			log.Printf("aggregator: get pending: %v", err)
			continue
		}

		var submitted int64
		for _, contrib := range contribs {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := a.submitOne(ctx, contrib); err != nil {
				log.Printf("aggregator: submit contribution %s: %v", contrib.ID, err)
				break
			}
			submitted++
		}

		if submitted > 0 {
			if err := a.database.IncrementPublishedCount(ctx, c.ConditionGroup, submitted); err != nil {
				log.Printf("aggregator: increment published count: %v", err)
			}
		}
	}
}

func (a *Aggregator) submitOne(ctx context.Context, c *db.Contribution) error {
	var pA [2]*big.Int
	pA[0] = new(big.Int)
	pA[1] = new(big.Int)
	pA[0].SetString(c.PA0, 10)
	pA[1].SetString(c.PA1, 10)

	var pB [2][2]*big.Int
	pB[0][0] = new(big.Int)
	pB[0][1] = new(big.Int)
	pB[1][0] = new(big.Int)
	pB[1][1] = new(big.Int)
	pB[0][0].SetString(c.PB00, 10)
	pB[0][1].SetString(c.PB01, 10)
	pB[1][0].SetString(c.PB10, 10)
	pB[1][1].SetString(c.PB11, 10)

	var pC [2]*big.Int
	pC[0] = new(big.Int)
	pC[1] = new(big.Int)
	pC[0].SetString(c.PC0, 10)
	pC[1].SetString(c.PC1, 10)

	commitment := common.HexToHash(c.Commitment)
	nullifier  := common.HexToHash(c.Nullifier)

	// Parse the merkleRoot stored at submission time.
	// This is the root the proof was built against; it must still match
	// AttestationManager.getRoot() at tx execution time for the contract to accept it.
	merkleRoot := new(big.Int)
	merkleRoot.SetString(c.MerkleRoot, 10)

	tx, err := a.researchContributer.Contribute(a.auth, pA, pB, pC, commitment, nullifier, merkleRoot)
	if err != nil {
		return fmt.Errorf("contribute tx: %w", err)
	}

	log.Printf("aggregator: submitted contribution %s tx=%s", c.ID, tx.Hash().Hex())

	if err := a.database.MarkContributionSubmitted(ctx, c.ID); err != nil {
		return fmt.Errorf("mark submitted: %w", err)
	}

	return nil
}

func (a *Aggregator) VerifyProof(
	ctx context.Context,
	pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int,
	commitment common.Hash, nullifier common.Hash,
	merkleRoot *big.Int,
) (bool, error) {
	verifierAddr, err := a.researchContributer.ResearchVerifier(nil)
	if err != nil {
		return false, fmt.Errorf("get verifier address: %w", err)
	}

	verifier, err := contracts.NewResearchVerifier(verifierAddr, a.client)
	if err != nil {
		return false, fmt.Errorf("bind verifier: %w", err)
	}

	// Public signal order must match the circuit and ResearchVerifier.sol:
	//   [0] merkleRoot  (public input)
	//   [1] commitment  (public output)
	//   [2] nullifier   (public output)
	var pubSignals [3]*big.Int
	pubSignals[0] = merkleRoot
	pubSignals[1] = new(big.Int).SetBytes(commitment.Bytes())
	pubSignals[2] = new(big.Int).SetBytes(nullifier.Bytes())

	return verifier.VerifyProof(nil, pA, pB, pC, pubSignals)
}

func (a *Aggregator) IsRegistered(ctx context.Context, commitment common.Hash) (bool, error) {
	return a.wellRegistry.IsRegistered(nil, commitment)
}

func (a *Aggregator) Address() common.Address {
	return a.auth.From
}

var _ *ecdsa.PublicKey
