package listener

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"Thesis/well-blockchain/well-backend/contracts"
	"Thesis/well-blockchain/well-backend/internal/api"
	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/pre"
	"Thesis/well-blockchain/well-backend/internal/vault"
)

type Listener struct {
	client           *ethclient.Client
	db               *db.DB
	vault            *vault.Client
	sessionStore     *api.SessionStore
	registry         *contracts.Registry
	userRegistry     *contracts.UserRegistry
	historyFactory   *contracts.HistoryFactory
	zkVerifier       *contracts.ZKAuthVerifier
	emergencyFactory *contracts.EmergencySummaryFactory
	incidentManager  *contracts.IncidentManager
	riskAssessment   *contracts.RiskAssessment

	mu            sync.Mutex
	historySubs   map[common.Address]*contracts.History
	recordSubs    map[common.Address]*contracts.Record
	emergencySubs map[common.Address]*contracts.EmergencySummary
}

func New(
	ctx context.Context,
	rpcURL string,
	database *db.DB,
	vaultClient *vault.Client,
	sessions *api.SessionStore,
	registryAddress string,
	userRegistryAddress string,
	historyFactoryAddress string,
	zkVerifierAddress string,
	emergencySummaryFactoryAddress string,
	incidentManagerAddress string,
	riskAssessmentAddress string,
) (*Listener, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect to ethereum node: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}
	log.Printf("connected to ethereum node, chain ID: %v", chainID)

	registry, err := contracts.NewRegistry(common.HexToAddress(registryAddress), client)
	if err != nil {
		return nil, fmt.Errorf("bind registry: %w", err)
	}

	userRegistry, err := contracts.NewUserRegistry(common.HexToAddress(userRegistryAddress), client)
	if err != nil {
		return nil, fmt.Errorf("bind user registry: %w", err)
	}

	historyFactory, err := contracts.NewHistoryFactory(common.HexToAddress(historyFactoryAddress), client)
	if err != nil {
		return nil, fmt.Errorf("bind history factory: %w", err)
	}

	l := &Listener{
		client:         client,
		db:             database,
		vault:          vaultClient,
		sessionStore:   sessions,
		registry:       registry,
		userRegistry:   userRegistry,
		historyFactory: historyFactory,
		historySubs:    make(map[common.Address]*contracts.History),
		recordSubs:     make(map[common.Address]*contracts.Record),
		emergencySubs:  make(map[common.Address]*contracts.EmergencySummary),
	}

	if zkVerifierAddress != "" {
		zkv, err := contracts.NewZKAuthVerifier(common.HexToAddress(zkVerifierAddress), client)
		if err != nil {
			return nil, fmt.Errorf("bind zk verifier: %w", err)
		}
		l.zkVerifier = zkv
		log.Printf("bound ZKAuthVerifier at %s", zkVerifierAddress)
	}

	if emergencySummaryFactoryAddress != "" {
		esf, err := contracts.NewEmergencySummaryFactory(common.HexToAddress(emergencySummaryFactoryAddress), client)
		if err != nil {
			return nil, fmt.Errorf("bind emergency summary factory: %w", err)
		}
		l.emergencyFactory = esf
		log.Printf("bound EmergencySummaryFactory at %s", emergencySummaryFactoryAddress)
	}

	if incidentManagerAddress != "" {
		im, err := contracts.NewIncidentManager(common.HexToAddress(incidentManagerAddress), client)
		if err != nil {
			return nil, fmt.Errorf("bind incident manager: %w", err)
		}
		l.incidentManager = im
		log.Printf("bound IncidentManager at %s", incidentManagerAddress)
	}

	if riskAssessmentAddress != "" {
		ra, err := contracts.NewRiskAssessment(common.HexToAddress(riskAssessmentAddress), client)
		if err != nil {
			return nil, fmt.Errorf("bind risk assessment: %w", err)
		}
		l.riskAssessment = ra
		log.Printf("bound RiskAssessment at %s", riskAssessmentAddress)
	}

	return l, nil
}

func (l *Listener) Start(ctx context.Context) error {
	//Registry subscriptions

	if err := l.subscribeToRegistry(ctx); err != nil {
		return fmt.Errorf("registry: %w", err)
	}

	//ZKAuthVerifier subscriptions
	if l.zkVerifier != nil {
		if err := l.subscribeToZKAuthVerifier(ctx); err != nil {
			return fmt.Errorf("zk verifier: %w", err)
		}
	}

	//EmergencySummary subscriptions (factory -> per-patient dynamic)
	if l.emergencyFactory != nil {
		if err := l.subscribeToEmergencySummaryFactory(ctx); err != nil {
			return fmt.Errorf("emergency summary: %w", err)
		}
	}

	//IncidentManager subscriptions
	if l.incidentManager != nil {
		if err := l.subscribeToIncidentManager(ctx); err != nil {
			return fmt.Errorf("incident manager: %w", err)
		}
	}

	//RiskAssessment subscriptions
	if l.riskAssessment != nil {
		if err := l.subscribeToRiskAssessment(ctx); err != nil {
			return fmt.Errorf("risk assessment: %w", err)
		}
	}

	return nil
}

func (l *Listener) Close() {
	l.client.Close()
}

func (l *Listener) storeRecordMeta(ctx context.Context, recordAddr common.Address, uploader common.Address) {
	record, err := contracts.NewRecord(recordAddr, l.client)
	if err != nil {
		log.Printf("storeRecordMeta: bind Record(%s): %v", recordAddr, err)
		return
	}

	contentHash, err := record.ContentHash(nil)
	if err != nil {
		log.Printf("storeRecordMeta: ContentHash(%s): %v", recordAddr, err)
		return
	}

	storageRef, err := record.StorageReference(nil)
	if err != nil {
		log.Printf("storeRecordMeta: StorageReference(%s): %v", recordAddr, err)
		return
	}

	owner, err := record.Owner(nil)
	if err != nil {
		log.Printf("storeRecordMeta: Owner(%s): %v", recordAddr, err)
		return
	}

	log.Printf("storeRecordMeta: record=%s cid=%s hash=%s owner=%s uploader=%s",
		recordAddr.Hex(), storageRef, common.Hash(contentHash).Hex(), owner.Hex(), uploader.Hex())
}

func (l *Listener) triggerPRE(ctx context.Context, recordAddr, doctorAddr common.Address) {
	if l.db == nil {
		log.Printf("triggerPRE: db not set, skipping")
		return
	}
	log.Printf("triggerPRE: record=%s doctor=%s", recordAddr.Hex(), doctorAddr.Hex())

	encRec, err := l.db.GetEncryptedKey(ctx, strings.ToLower(recordAddr.Hex()))
	if err != nil {
		log.Printf("triggerPRE: GetEncryptedKey(%s): %v (no lockbox stored yet, skipping)", recordAddr.Hex(), err)
		return
	}
	if len(encRec.LockboxSalt) == 0 || len(encRec.PatientPublicKey) == 0 {
		log.Printf("triggerPRE: record=%s predates PRE, skipping", recordAddr.Hex())
		return
	}

	rkBytes, err := l.db.GetReEncryptionKey(ctx, encRec.PatientAddress, strings.ToLower(doctorAddr.Hex()))
	if err != nil {
		log.Printf("triggerPRE: no re-encryption key for record=%s doctor=%s (patient must delegate first), skipping", recordAddr.Hex(), doctorAddr.Hex())
		return
	}

	lockbox, err := pre.ParseDelegatorCiphertext(encRec.EncryptedKey)
	if err != nil {
		log.Printf("triggerPRE: parse stored lockbox: %v", err)
		return
	}
	rk, err := pre.ParseReEncryptionKey(rkBytes)
	if err != nil {
		log.Printf("triggerPRE: parse stored re-encryption key: %v", err)
		return
	}

	reCT, err := pre.Reencrypt(lockbox, rk)
	if err != nil {
		log.Printf("triggerPRE: reencrypt: %v", err)
		return
	}
	reBytes, err := reCT.MarshalBinary()
	if err != nil {
		log.Printf("triggerPRE: encode re-encrypted ciphertext: %v", err)
		return
	}

	if err := l.db.InsertReEncryptedKey(ctx, &db.ReEncryptedKeyRecord{
		RecordAddress:  strings.ToLower(recordAddr.Hex()),
		DoctorAddress:  strings.ToLower(doctorAddr.Hex()),
		ReEncryptedKey: reBytes,
	}); err != nil {
		log.Printf("triggerPRE: InsertReEncryptedKey: %v", err)
		return
	}

	log.Printf("triggerPRE: complete record=%s doctor=%s", recordAddr.Hex(), doctorAddr.Hex())
}

//registry

func (l *Listener) subscribeToRegistry(ctx context.Context) error {
	//historyCreated(HistoryFactory)
	historyCreatedCh := make(chan *contracts.HistoryFactoryHistoryCreated)
	hfSub, err := l.historyFactory.WatchHistoryCreated(
		&bind.WatchOpts{Context: ctx},
		historyCreatedCh,
		nil, nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchHistoryCreated: %w", err)
	}

	//registered(UserRegistry)
	registeredCh := make(chan *contracts.UserRegistryRegistered)
	urRegSub, err := l.userRegistry.WatchRegistered(
		&bind.WatchOpts{Context: ctx},
		registeredCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchRegistered: %w", err)
	}

	//revoked(UserRegistry)
	revokedCh := make(chan *contracts.UserRegistryRevoked)
	urRevSub, err := l.userRegistry.WatchRevoked(
		&bind.WatchOpts{Context: ctx},
		revokedCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchRevoked: %w", err)
	}

	go func() {
		defer hfSub.Unsubscribe()
		defer urRegSub.Unsubscribe()
		defer urRevSub.Unsubscribe()

		for {
			select {
			case ev := <-historyCreatedCh:
				l.logEvent(ev.Raw, "HistoryFactory.HistoryCreated", map[string]any{
					"patient":         ev.Patient.Hex(),
					"historyContract": ev.HistoryContract.Hex(),
				})
				if err := l.db.AttachCommitment(ctx, ev.Patient.Hex(), common.Hash(ev.Commitment).Hex()); err != nil {
					log.Printf("identity index: attach commitment: %v", err)
				}
				l.subscribeToHistory(ctx, ev.HistoryContract)

			case ev := <-registeredCh:
				l.logEvent(ev.Raw, "UserRegistry.Registered", map[string]any{
					"commitment": common.Hash(ev.Commitment).Hex(),
					"internalID": common.Hash(ev.InternalID).Hex(),
					"role":       ev.Role,
				})
				if err := l.db.AttachRegisteredIdentity(
					ctx,
					common.Hash(ev.Commitment).Hex(),
					common.Hash(ev.InternalID).Hex(),
					roleName(ev.Role),
				); err != nil {
					log.Printf("identity index: attach registered: %v", err)
				}

			case ev := <-revokedCh:
				l.logEvent(ev.Raw, "UserRegistry.Revoked", map[string]any{
					"commitment": common.Hash(ev.Commitment).Hex(),
					"internalID": common.Hash(ev.InternalID).Hex(),
				})
				if err := l.db.TombstoneByCommitment(ctx, common.Hash(ev.Commitment).Hex()); err != nil {
					log.Printf("identity index: tombstone: %v", err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("subscribed to UserRegistry / HistoryFactory events")
	return nil
}

// ── ZKAuthVerifier (2 events) ──────────────────────────────

func (l *Listener) subscribeToZKAuthVerifier(ctx context.Context) error {
	authenticatedCh := make(chan *contracts.ZKAuthVerifierAuthenticated)
	authSub, err := l.zkVerifier.WatchAuthenticated(
		&bind.WatchOpts{Context: ctx},
		authenticatedCh,
		nil, nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchAuthenticated: %w", err)
	}

	sessionRevokedCh := make(chan *contracts.ZKAuthVerifierSessionRevoked)
	srSub, err := l.zkVerifier.WatchSessionRevoked(
		&bind.WatchOpts{Context: ctx},
		sessionRevokedCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchZKAuthVerifierSessionRevoked: %w", err)
	}

	go func() {
		defer authSub.Unsubscribe()
		defer srSub.Unsubscribe()

		for {
			select {
			case ev := <-authenticatedCh:
				l.logEvent(ev.Raw, "ZKAuthVerifier.Authenticated", map[string]any{
					"user":         ev.User.Hex(),
					"commitment":   common.Hash(ev.Commitment).Hex(),
					"sessionToken": common.Hash(ev.SessionToken).Hex(),
					"expiresAt":    ev.ExpiresAt.String(),
				})
				expiresAt := time.Unix(ev.ExpiresAt.Int64(), 0)
				l.sessionStore.Add(ev.User.Hex(), common.Hash(ev.SessionToken).Hex(), expiresAt)

			case ev := <-sessionRevokedCh:
				l.logEvent(ev.Raw, "ZKAuthVerifier.SessionRevoked", map[string]any{
					"sessionToken": common.Hash(ev.SessionToken).Hex(),
					"user":         ev.User.Hex(),
				})

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("subscribed to ZKAuthVerifier events")
	return nil
}

// ── EmergencySummaryFactory (1 event, dynamic) ────────────

func (l *Listener) subscribeToEmergencySummaryFactory(ctx context.Context) error {
	createdCh := make(chan *contracts.EmergencySummaryFactoryEmergencySummaryCreated)
	csSub, err := l.emergencyFactory.WatchEmergencySummaryCreated(
		&bind.WatchOpts{Context: ctx},
		createdCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchEmergencySummaryCreated: %w", err)
	}

	go func() {
		defer csSub.Unsubscribe()

		for {
			select {
			case ev := <-createdCh:
				l.logEvent(ev.Raw, "EmergencySummaryFactory.EmergencySummaryCreated", map[string]any{
					"commitment":       common.Hash(ev.Commitment).Hex(),
					"patient":          ev.Patient.Hex(),
					"emergencySummary": ev.EmergencySummary.Hex(),
				})
				if err := l.db.AttachSummary(ctx, ev.Patient.Hex(), ev.EmergencySummary.Hex()); err != nil {
					log.Printf("identity index: attach summary: %v", err)
				}
				l.subscribeToEmergencySummary(ctx, ev.EmergencySummary, ev.Patient)

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("subscribed to EmergencySummaryFactory events")
	return nil
}

// ── EmergencySummary (5 events, dynamic) ────────────────

func (l *Listener) subscribeToEmergencySummary(ctx context.Context, summaryAddress, patientAddress common.Address) {
	l.mu.Lock()
	if _, exists := l.emergencySubs[summaryAddress]; exists {
		l.mu.Unlock()
		return
	}
	summary, err := contracts.NewEmergencySummary(summaryAddress, l.client)
	if err != nil {
		l.mu.Unlock()
		log.Printf("EmergencySummary(%s): bind: %v", summaryAddress, err)
		return
	}
	l.emergencySubs[summaryAddress] = summary
	l.mu.Unlock()

	dataUpdatedCh := make(chan *contracts.EmergencySummaryDataUpdated)
	duSub, err := summary.WatchDataUpdated(
		&bind.WatchOpts{Context: ctx},
		dataUpdatedCh,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchDataUpdated: %v", summaryAddress, err)
		return
	}

	accessRequestedCh := make(chan *contracts.EmergencySummaryAccessRequested)
	arSub, err := summary.WatchAccessRequested(
		&bind.WatchOpts{Context: ctx},
		accessRequestedCh,
		nil, nil,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchAccessRequested: %v", summaryAddress, err)
		return
	}

	accessGrantedCh := make(chan *contracts.EmergencySummaryAccessGranted)
	agSub, err := summary.WatchAccessGranted(
		&bind.WatchOpts{Context: ctx},
		accessGrantedCh,
		nil,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchAccessGranted: %v", summaryAddress, err)
		return
	}

	accessRevokedCh := make(chan *contracts.EmergencySummaryAccessRevoked)
	arvSub, err := summary.WatchAccessRevoked(
		&bind.WatchOpts{Context: ctx},
		accessRevokedCh,
		nil,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchAccessRevoked: %v", summaryAddress, err)
		return
	}

	publicKeyRegisteredCh := make(chan *contracts.EmergencySummaryPublicKeyRegistered)
	pkrSub, err := summary.WatchPublicKeyRegistered(
		&bind.WatchOpts{Context: ctx},
		publicKeyRegisteredCh,
		nil,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchPublicKeyRegistered: %v", summaryAddress, err)
		return
	}

	recoveryKeyStoredCh := make(chan *contracts.EmergencySummaryRecoveryKeyStored)
	rkSub, err := summary.WatchRecoveryKeyStored(
		&bind.WatchOpts{Context: ctx},
		recoveryKeyStoredCh,
		nil,
	)
	if err != nil {
		log.Printf("EmergencySummary(%s).WatchRecoveryKeyStored: %v", summaryAddress, err)
		return
	}

	go func() {
		defer duSub.Unsubscribe()
		defer arSub.Unsubscribe()
		defer agSub.Unsubscribe()
		defer arvSub.Unsubscribe()
		defer pkrSub.Unsubscribe()
		defer rkSub.Unsubscribe()

		for {
			select {
			case ev := <-dataUpdatedCh:
				l.logEvent(ev.Raw, "EmergencySummary.DataUpdated", map[string]any{
					"version":   ev.Version.String(),
					"dataHash":  common.Hash(ev.DataHash).Hex(),
					"timestamp": ev.Timestamp.String(),
				})

			case ev := <-accessRequestedCh:
				l.logEvent(ev.Raw, "EmergencySummary.AccessRequested", map[string]any{
					"requestId": ev.RequestId.String(),
					"responder": ev.Responder.Hex(),
					"timestamp": ev.Timestamp.String(),
				})

			case ev := <-accessGrantedCh:
				l.logEvent(ev.Raw, "EmergencySummary.AccessGranted", map[string]any{
					"user":      ev.User.Hex(),
					"expiresAt": ev.ExpiresAt.String(),
				})
				if err := l.db.InsertAccessLog(ctx, db.AccessLog{
					AccessorAddress: ev.User.Hex(),
					PatientAddress:  patientAddress.Hex(),
					RecordID:        summaryAddress.Hex(),
					AccessType:      "emergency_grant",
					WasEmergency:    true,
					TxHash:          ev.Raw.TxHash.Hex(),
				}); err != nil {
					log.Printf("EmergencySummary.AccessGranted: insert access log: %v", err)
				}

			case ev := <-accessRevokedCh:
				l.logEvent(ev.Raw, "EmergencySummary.AccessRevoked", map[string]any{
					"user": ev.User.Hex(),
				})

			case ev := <-publicKeyRegisteredCh:
				l.logEvent(ev.Raw, "EmergencySummary.PublicKeyRegistered", map[string]any{
					"user": ev.User.Hex(),
				})

			case ev := <-recoveryKeyStoredCh:
				l.logEvent(ev.Raw, "EmergencySummary.RecoveryKeyStored", map[string]any{
					"version":   ev.Version.String(),
					"timestamp": ev.Timestamp.String(),
				})

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("subscribed: EmergencySummary(%s)", summaryAddress.Hex())
}

// ── IncidentManager (3 events) ───────────────────────────

func (l *Listener) subscribeToIncidentManager(ctx context.Context) error {
	incidentCreatedCh := make(chan *contracts.IncidentManagerIncidentCreated)
	icSub, err := l.incidentManager.WatchIncidentCreated(
		&bind.WatchOpts{Context: ctx},
		incidentCreatedCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchIncidentCreated: %w", err)
	}

	respondersAddedCh := make(chan *contracts.IncidentManagerRespondersAdded)
	raSub, err := l.incidentManager.WatchRespondersAdded(
		&bind.WatchOpts{Context: ctx},
		respondersAddedCh,
		nil,
	)
	if err != nil {
		return fmt.Errorf("WatchRespondersAdded: %w", err)
	}

	incidentDeactivatedCh := make(chan *contracts.IncidentManagerIncidentDeactivated)
	idSub, err := l.incidentManager.WatchIncidentDeactivated(
		&bind.WatchOpts{Context: ctx},
		incidentDeactivatedCh,
		nil,
	)
	if err != nil {
		return fmt.Errorf("WatchIncidentDeactivated: %w", err)
	}

	go func() {
		defer icSub.Unsubscribe()
		defer raSub.Unsubscribe()
		defer idSub.Unsubscribe()

		for {
			select {
			case ev := <-incidentCreatedCh:
				l.logEvent(ev.Raw, "IncidentManager.IncidentCreated", map[string]any{
					"id":        common.Hash(ev.Id).Hex(),
					"hospital":  ev.Hospital.Hex(),
					"expiresAt": ev.ExpiresAt.String(),
				})

			case ev := <-respondersAddedCh:
				addrs := make([]string, len(ev.Responders))
				for i, a := range ev.Responders {
					addrs[i] = a.Hex()
				}
				l.logEvent(ev.Raw, "IncidentManager.RespondersAdded", map[string]any{
					"id":         common.Hash(ev.Id).Hex(),
					"responders": addrs,
				})

			case ev := <-incidentDeactivatedCh:
				l.logEvent(ev.Raw, "IncidentManager.IncidentDeactivated", map[string]any{
					"id": common.Hash(ev.Id).Hex(),
				})

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("subscribed to IncidentManager events")
	return nil
}

// ── RiskAssessment (1 event) ─────────────────────────────

func (l *Listener) subscribeToRiskAssessment(ctx context.Context) error {
	emergencyGrantedCh := make(chan *contracts.RiskAssessmentEmergencyGranted)
	egSub, err := l.riskAssessment.WatchEmergencyGranted(
		&bind.WatchOpts{Context: ctx},
		emergencyGrantedCh,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("WatchEmergencyGranted: %w", err)
	}

	go func() {
		defer egSub.Unsubscribe()

		for {
			select {
			case ev := <-emergencyGrantedCh:
				l.logEvent(ev.Raw, "RiskAssessment.EmergencyGranted", map[string]any{
					"responder":  ev.Responder.Hex(),
					"incidentId": common.Hash(ev.IncidentId).Hex(),
					"expiresAt":  ev.ExpiresAt.String(),
				})

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("subscribed to RiskAssessment events")
	return nil
}

// ── History (3 events, dynamic) ────────────────────────────

func (l *Listener) subscribeToHistory(ctx context.Context, historyAddress common.Address) {
	l.mu.Lock()
	if _, exists := l.historySubs[historyAddress]; exists {
		l.mu.Unlock()
		return
	}
	history, err := contracts.NewHistory(historyAddress, l.client)
	if err != nil {
		l.mu.Unlock()
		log.Printf("History(%s): bind: %v", historyAddress, err)
		return
	}
	l.historySubs[historyAddress] = history
	l.mu.Unlock()

	// AccessGranted
	accessGrantedCh := make(chan *contracts.HistoryAccessGranted)
	agSub, err := history.WatchAccessGranted(
		&bind.WatchOpts{Context: ctx},
		accessGrantedCh,
		nil, nil,
	)
	if err != nil {
		log.Printf("History(%s).WatchAccessGranted: %v", historyAddress, err)
		return
	}

	// AccessRevoked
	accessRevokedCh := make(chan *contracts.HistoryAccessRevoked)
	arSub, err := history.WatchAccessRevoked(
		&bind.WatchOpts{Context: ctx},
		accessRevokedCh,
		nil, nil,
	)
	if err != nil {
		log.Printf("History(%s).WatchAccessRevoked: %v", historyAddress, err)
		return
	}

	// RecordCreated — triggers dynamic Record subscription
	recordCreatedCh := make(chan *contracts.HistoryRecordCreated)
	rcSub, err := history.WatchRecordCreated(
		&bind.WatchOpts{Context: ctx},
		recordCreatedCh,
		nil, nil,
	)
	if err != nil {
		log.Printf("History(%s).WatchRecordCreated: %v", historyAddress, err)
		return
	}

	go func() {
		defer agSub.Unsubscribe()
		defer arSub.Unsubscribe()
		defer rcSub.Unsubscribe()

		for {
			select {
			case ev := <-accessGrantedCh:
				patient, err := history.Owner(nil)
				if err != nil {
					log.Printf("History(%s).Owner(): %v", historyAddress, err)
					continue
				}
				if err := l.db.InsertAccessLog(ctx, db.AccessLog{
					AccessorAddress: ev.User.Hex(),
					PatientAddress:  patient.Hex(),
					RecordID:        ev.Record.Hex(),
					AccessType:      "grant",
					TxHash:          ev.Raw.TxHash.Hex(),
					WasEmergency:    false,
				}); err != nil {
					log.Printf("insert access log: %v", err)
				}

			case ev := <-accessRevokedCh:
				l.logEvent(ev.Raw, "History.AccessRevoked", map[string]any{
					"record": ev.Record.Hex(),
					"user":   ev.User.Hex(),
				})

			case ev := <-recordCreatedCh:
				log.Printf("RecordCreated: record=%s uploader=%s", ev.Record.Hex(), ev.Uploader.Hex())
				l.storeRecordMeta(ctx, ev.Record, ev.Uploader)
				l.subscribeToRecord(ctx, ev.Record)

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("subscribed: History(%s)", historyAddress.Hex())
}

// ── Record (4 events, dynamic) ─────────────────────────────

func (l *Listener) subscribeToRecord(ctx context.Context, recordAddress common.Address) {
	l.mu.Lock()
	if _, exists := l.recordSubs[recordAddress]; exists {
		l.mu.Unlock()
		return
	}
	record, err := contracts.NewRecord(recordAddress, l.client)
	if err != nil {
		l.mu.Unlock()
		log.Printf("Record(%s): bind: %v", recordAddress, err)
		return
	}
	l.recordSubs[recordAddress] = record
	l.mu.Unlock()

	accessGrantedCh := make(chan *contracts.RecordAccessGranted)
	agSub, err := record.WatchAccessGranted(
		&bind.WatchOpts{Context: ctx},
		accessGrantedCh,
		nil, nil,
	)
	if err != nil {
		log.Printf("Record(%s).WatchAccessGranted: %v", recordAddress, err)
		return
	}

	accessRevokedCh := make(chan *contracts.RecordAccessRevoked)
	arSub, err := record.WatchAccessRevoked(
		&bind.WatchOpts{Context: ctx},
		accessRevokedCh,
		nil,
	)
	if err != nil {
		log.Printf("Record(%s).WatchAccessRevoked: %v", recordAddress, err)
		return
	}

	policyUpdatedCh := make(chan *contracts.RecordPolicyUpdated)
	puSub, err := record.WatchPolicyUpdated(
		&bind.WatchOpts{Context: ctx},
		policyUpdatedCh,
	)
	if err != nil {
		log.Printf("Record(%s).WatchPolicyUpdated: %v", recordAddress, err)
		return
	}

	integrityHashUpdatedCh := make(chan *contracts.RecordIntegrityHashUpdated)
	ihuSub, err := record.WatchIntegrityHashUpdated(
		&bind.WatchOpts{Context: ctx},
		integrityHashUpdatedCh,
	)
	if err != nil {
		log.Printf("Record(%s).WatchIntegrityHashUpdated: %v", recordAddress, err)
		return
	}

	go func() {
		defer agSub.Unsubscribe()
		defer arSub.Unsubscribe()
		defer puSub.Unsubscribe()
		defer ihuSub.Unsubscribe()

		for {
			select {
			case ev := <-accessGrantedCh:
				l.logEvent(ev.Raw, "Record.AccessGranted", map[string]any{
					"doctor":    ev.Doctor.Hex(),
					"grantedBy": ev.GrantedBy.Hex(),
					"grantedAt": ev.GrantedAt.String(),
				})
				l.triggerPRE(ctx, recordAddress, ev.Doctor)

			case ev := <-accessRevokedCh:
				l.logEvent(ev.Raw, "Record.AccessRevoked", map[string]any{
					"doctor":    ev.Doctor.Hex(),
					"revokedAt": ev.RevokedAt.String(),
				})
				if err := l.db.DeleteReEncryptedKey(ctx, strings.ToLower(recordAddress.Hex()), strings.ToLower(ev.Doctor.Hex())); err != nil {
					log.Printf("Record.AccessRevoked: DeleteReEncryptedKey(%s, %s): %v", recordAddress.Hex(), ev.Doctor.Hex(), err)
				}

			case ev := <-policyUpdatedCh:
				l.logEvent(ev.Raw, "Record.PolicyUpdated", map[string]any{
					"requiredDeptCommitment": common.Hash(ev.RequiredDeptCommitment).Hex(),
					"requiredInstCommitment": common.Hash(ev.RequiredInstCommitment).Hex(),
				})

			case ev := <-integrityHashUpdatedCh:
				l.logEvent(ev.Raw, "Record.IntegrityHashUpdated", map[string]any{
					"newHash": common.Hash(ev.NewHash).Hex(),
				})

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("subscribed: Record(%s)", recordAddress.Hex())
}

// ── helpers ────────────────────────────────────────────────

// roleName maps a UserRegistry role number (Patient=0, Doctor=1,
// FirstResponder=2, Admin=3, Hospital=4) to the identity-index role string.
func roleName(role uint8) string {
	switch role {
	case 0:
		return "patient"
	case 1:
		return "doctor"
	case 2:
		return "first_responder"
	case 3:
		return "admin"
	case 4:
		return "hospital"
	default:
		return "patient"
	}
}

func (l *Listener) logEvent(raw types.Log, eventType string, data map[string]any) {
	log.Printf("event received: %s tx=%s block=%d", eventType, raw.TxHash.Hex(), raw.BlockNumber)
	if err := l.db.InsertEventLog(context.Background(), db.EventLog{
		EventType:       eventType,
		ContractAddress: raw.Address.Hex(),
		TxHash:          raw.TxHash.Hex(),
		BlockNumber:     raw.BlockNumber,
		Data:            data,
	}); err != nil {
		log.Printf("log event (%s): %v", eventType, err)
	} else {
		log.Printf("event logged: %s", eventType)
	}
}
