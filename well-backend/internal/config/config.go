package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBURL                   string
	EthRPC                  string
	RegistryAddress         string
	UserRegistryAddress     string
	HistoryFactoryAddress   string
	ZKAuthVerifierAddress   string
	EmergencySummaryFactory string
	IncidentManagerAddress  string
	RiskAssessmentAddress   string
	ResearchContributerAddr string
	AggregatorPrivateKey    string
	AggregatorPort          string
	AggregatorThreshold     int64

	LogAnchorAddress     string
	AnchorPrivateKey     string
	LedgerBatchSize      int
	LedgerAnchorInterval int

	IPFSURL    string
	VaultAddr  string
	VaultToken string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	ethRPC := os.Getenv("ETH_RPC")
	if ethRPC == "" {
		return nil, fmt.Errorf("ETH_RPC environment variable is required")
	}

	registryAddress := os.Getenv("REGISTRY_ADDRESS")
	if registryAddress == "" {
		return nil, fmt.Errorf("REGISTRY_ADDRESS environment variable is required")
	}

	userRegistryAddress := os.Getenv("USER_REGISTRY_ADDRESS")
	if userRegistryAddress == "" {
		return nil, fmt.Errorf("USER_REGISTRY_ADDRESS environment variable is required")
	}

	historyFactoryAddress := os.Getenv("HISTORY_FACTORY_ADDRESS")
	if historyFactoryAddress == "" {
		return nil, fmt.Errorf("HISTORY_FACTORY_ADDRESS environment variable is required")
	}

	cfg := &Config{
		DBURL:                   dbURL,
		EthRPC:                  ethRPC,
		RegistryAddress:         registryAddress,
		UserRegistryAddress:     userRegistryAddress,
		HistoryFactoryAddress:   historyFactoryAddress,
		ZKAuthVerifierAddress:   os.Getenv("ZKAUTH_VERIFIER_ADDRESS"),
		EmergencySummaryFactory: os.Getenv("EMERGENCY_SUMMARY_FACTORY_ADDRESS"),
		IncidentManagerAddress:  os.Getenv("INCIDENT_MANAGER_ADDRESS"),
		RiskAssessmentAddress:   os.Getenv("RISK_ASSESSMENT_ADDRESS"),
		ResearchContributerAddr: os.Getenv("RESEARCH_CONTRIBUTER_ADDRESS"),
		AggregatorPrivateKey:    os.Getenv("AGGREGATOR_PRIVATE_KEY"),
		AggregatorPort:          os.Getenv("AGGREGATOR_PORT"),
		AggregatorThreshold:     3,
		IPFSURL:                 os.Getenv("IPFS_API_URL"),
		VaultAddr:               os.Getenv("VAULT_ADDR"),
		VaultToken:              os.Getenv("VAULT_TOKEN"),
		LogAnchorAddress:        os.Getenv("LOG_ANCHOR_ADDRESS"),
		AnchorPrivateKey:        os.Getenv("ANCHOR_PRIVATE_KEY"),
		LedgerBatchSize:         100,
		LedgerAnchorInterval:    60,
	}

	if v := os.Getenv("AGGREGATOR_THRESHOLD"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &cfg.AggregatorThreshold); err != nil {
			return nil, fmt.Errorf("invalid AGGREGATOR_THRESHOLD: %w", err)
		}
	}
	if cfg.AggregatorPort == "" {
		cfg.AggregatorPort = "8080"
	}
	if cfg.IPFSURL == "" {
		cfg.IPFSURL = "http://localhost:5001"
	}
	if cfg.VaultAddr == "" {
		cfg.VaultAddr = "http://localhost:8200"
	}

	if v := os.Getenv("LEDGER_BATCH_SIZE"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &cfg.LedgerBatchSize); err != nil || cfg.LedgerBatchSize < 1 {
			return nil, fmt.Errorf("invalid LEDGER_BATCH_SIZE: %w", err)
		}
	}
	if v := os.Getenv("LEDGER_ANCHOR_INTERVAL_SECS"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &cfg.LedgerAnchorInterval); err != nil || cfg.LedgerAnchorInterval < 1 {
			return nil, fmt.Errorf("invalid LEDGER_ANCHOR_INTERVAL_SECS: %w", err)
		}
	}

	return cfg, nil
}
