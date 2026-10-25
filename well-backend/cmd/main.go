package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/aggregator"
	"Thesis/well-blockchain/well-backend/internal/api"
	"Thesis/well-blockchain/well-backend/internal/config"
	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/ipfs"
	"Thesis/well-blockchain/well-backend/internal/ledger"
	"Thesis/well-blockchain/well-backend/internal/listener"
	"Thesis/well-blockchain/well-backend/internal/vault"
)

// chainRoleProvider resolves an address's on-chain role via UserRegistry.
type chainRoleProvider struct {
	userRegistry *contracts.UserRegistry
}

func (p *chainRoleProvider) RoleOf(ctx context.Context, addr string) (uint8, error) {
	return p.userRegistry.GetRoleByAddress(nil, common.HexToAddress(addr))
}

// emergencyResolver reads owner/access state of an EmergencySummary on-chain.
type emergencyResolver struct {
	client *ethclient.Client
}

func (r *emergencyResolver) OwnerOf(ctx context.Context, summaryAddr string) (string, error) {
	summary, err := contracts.NewEmergencySummary(common.HexToAddress(summaryAddr), r.client)
	if err != nil {
		return "", err
	}
	owner, err := summary.Owner(nil)
	if err != nil {
		return "", err
	}
	return owner.Hex(), nil
}

func (r *emergencyResolver) HasAccess(ctx context.Context, summaryAddr string, caller string) (bool, error) {
	summary, err := contracts.NewEmergencySummary(common.HexToAddress(summaryAddr), r.client)
	if err != nil {
		return false, err
	}
	return summary.HasAccess(nil, common.HexToAddress(caller))
}

func main() {
	var lgr *ledger.Ledger

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	log.Println("connected to database")

	log.Println("backfilling ledger chains (idempotent)...")
	if err := database.BackfillLedger(ctx); err != nil {
		log.Fatalf("backfill ledger: %v", err)
	}
	log.Println("ledger backfill complete")

	vaultClient, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		log.Fatalf("vault: %v", err)
	}
	log.Println("connected to vault")

	ipfsClient := ipfs.New(cfg.IPFSURL)
	log.Println("connected to ipfs")

	sessions := api.NewSessionStore()

	lst, err := listener.New(
		ctx,
		cfg.EthRPC,
		database,
		vaultClient,
		sessions,
		cfg.RegistryAddress,
		cfg.UserRegistryAddress,
		cfg.HistoryFactoryAddress,
		cfg.ZKAuthVerifierAddress,
		cfg.EmergencySummaryFactory,
		cfg.IncidentManagerAddress,
		cfg.RiskAssessmentAddress,
	)
	if err != nil {
		log.Fatalf("listener: %v", err)
	}
	defer lst.Close()

	if err := lst.Start(ctx); err != nil {
		log.Fatalf("start listener: %v", err)
	}

	log.Println("listener started")

	ethClient, err := ethclient.DialContext(ctx, cfg.EthRPC)
	if err != nil {
		log.Fatalf("ethclient: %v", err)
	}
	userRegistry, err := contracts.NewUserRegistry(common.HexToAddress(cfg.UserRegistryAddress), ethClient)
	if err != nil {
		log.Fatalf("bind user registry: %v", err)
	}
	roles := &chainRoleProvider{userRegistry: userRegistry}
	log.Println("role provider ready")

	apiHandler := api.NewHandler(database, ipfsClient, vaultClient, roles, sessions, &emergencyResolver{client: ethClient})
	mux := http.NewServeMux()
	mux.Handle("/api/v1/records/", apiHandler)
	mux.Handle("/api/v1/attributes/", apiHandler)
	mux.Handle("/api/v1/emergency/", apiHandler)
	mux.Handle("/api/v1/identities", apiHandler)
	mux.Handle("/api/v1/identities/", apiHandler)

	if cfg.ResearchContributerAddr != "" && cfg.AggregatorPrivateKey != "" {
		agg, err := aggregator.New(
			ctx,
			cfg.EthRPC,
			database,
			cfg.AggregatorPrivateKey,
			cfg.ResearchContributerAddr,
			cfg.RegistryAddress,
			cfg.AggregatorThreshold,
		)
		if err != nil {
			log.Fatalf("aggregator: %v", err)
		}
		defer agg.Close()

		agg.Start(ctx)

		aggHandler := aggregator.NewHandler(agg, database)
		mux.Handle("/api/v1/contribute", aggHandler)
		mux.Handle("/api/v1/counters", aggHandler)

		log.Println("aggregator started")
	} else {
		log.Println("aggregator disabled (set RESEARCH_CONTRIBUTER_ADDRESS and AGGREGATOR_PRIVATE_KEY)")
	}

	if cfg.LogAnchorAddress != "" && cfg.AnchorPrivateKey != "" {
		lg, err := ledger.New(
			ctx,
			cfg.EthRPC,
			database,
			cfg.LogAnchorAddress,
			cfg.AnchorPrivateKey,
			cfg.LedgerBatchSize,
			time.Duration(cfg.LedgerAnchorInterval)*time.Second,
		)
		if err != nil {
			log.Fatalf("ledger: %v", err)
		}
		defer lg.Close()
		lg.Start(ctx)
		lgr = lg
		log.Println("ledger started")
	} else {
		log.Println("ledger disabled (set LOG_ANCHOR_ADDRESS and ANCHOR_PRIVATE_KEY)")
	}

	go func() {
		log.Printf("HTTP server listening on :%s", cfg.AggregatorPort)
		if err := http.ListenAndServe(":"+cfg.AggregatorPort, mux); err != nil {
			log.Fatalf("http: %v", err)
		}
	}()

	log.Println("waiting...")
	<-ctx.Done()

	log.Println("shutdown signal received; flushing ledger")
	if lgr != nil {
		lgr.Stop()
		flushCtx, flushCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer flushCancel()
		if err := lgr.Flush(flushCtx); err != nil {
			log.Printf("ledger flush on shutdown: %v", err)
		} else {
			log.Println("ledger flushed: all batches sealed and anchored")
		}
	}
	log.Println("shutdown complete")
}
