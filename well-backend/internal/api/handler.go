package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/pre"
	"Thesis/well-blockchain/well-backend/internal/vault"
)

type RecordStore interface {
	InsertEncryptedKey(ctx context.Context, r *db.EncryptedKeyRecord) error
	GetEncryptedKey(ctx context.Context, recordAddress string) (*db.EncryptedKeyRecord, error)
	InsertReEncryptedKey(ctx context.Context, r *db.ReEncryptedKeyRecord) error
	GetReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) ([]byte, error)
	DeleteReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) error
	UpsertReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string, rk []byte) error
	GetReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string) ([]byte, error)

	UpsertIdentity(ctx context.Context, e *db.Identity) error
	GetIdentityByNationalID(ctx context.Context, nationalID string) (*db.Identity, error)
	GetIdentitiesByName(ctx context.Context, name string) ([]*db.Identity, error)
	DeactivateIdentity(ctx context.Context, nationalID string) error
	RecordLookup(ctx context.Context, l *db.IdentityLookup) error
}

type IPFSClient interface {
	Add(ctx context.Context, data []byte) (string, error)
	Cat(ctx context.Context, cid string) ([]byte, error)
}

type VaultClient interface {
	EnsurePatientKey(ctx context.Context, patientAddr string) error
	WrapKeyWithPatientKey(ctx context.Context, patientAddr string, contentKey []byte) ([]byte, error)
	Rekey(ctx context.Context, patientAddr, doctorAddr string, ciphertext []byte) ([]byte, error)

	// Off-chain attribute issuance (Vault KV).
	StoreAttributes(ctx context.Context, addr string, attrs *vault.Attributes) error
	GetAttributes(ctx context.Context, addr string) (*vault.Attributes, error)
	IssueUUID(ctx context.Context, kind, name string) (string, error)
	ResolveUUID(ctx context.Context, kind, name string) (string, error)

	// Emergency recovery escrow (Vault KV).
	StoreRecoveryKey(ctx context.Context, summaryAddr string, recoveryKey []byte) error
	GetRecoveryKey(ctx context.Context, summaryAddr string) ([]byte, error)
}

// EmergencyResolver reads owner/access state of an EmergencySummary on-chain.
type EmergencyResolver interface {
	OwnerOf(ctx context.Context, summaryAddr string) (string, error)
	HasAccess(ctx context.Context, summaryAddr string, caller string) (bool, error)
}

// RoleProvider resolves an address's on-chain role (UserRegistry.getRoleByAddress).
// Roles: Patient=0, Doctor=1, FirstResponder=2, Admin=3, Hospital=4.
type RoleProvider interface {
	RoleOf(ctx context.Context, addr string) (uint8, error)
}

const (
	RolePatient        uint8 = 0
	RoleDoctor         uint8 = 1
	RoleFirstResponder uint8 = 2
	RoleAdmin          uint8 = 3
	RoleHospital       uint8 = 4
)

type Handler struct {
	store     RecordStore
	ipfs      IPFSClient
	vault     VaultClient
	roles     RoleProvider
	sessions  *SessionStore
	emergency EmergencyResolver
}

type uploadResponse struct {
	IPFSCid     string `json:"ipfs_cid"`
	ContentHash string `json:"content_hash"`
}

type keyStoreRequest struct {
	RecordAddress string `json:"record_address"`
	IPFSCid       string `json:"ipfs_cid"`
	ContentHash   string `json:"content_hash"`
	EncryptedKey  string `json:"encrypted_key"`
	PatientAddr   string `json:"patient_address"`
	// The following two fields select the PRE path. When both are present,
	// EncryptedKey is the patient's canonical PRE lockbox (C1) and is stored
	// verbatim after validation. When absent, the legacy Vault-wrapped path is
	// used for backwards compatibility.
	LockboxSalt      string `json:"lockbox_salt"`
	PatientPublicKey string `json:"patient_public_key"`
}

type delegateRequest struct {
	PatientAddress  string `json:"patient_address"`
	DoctorAddress   string `json:"doctor_address"`
	ReEncryptionKey string `json:"re_encryption_key"`
}

type reencryptRequest struct {
	DoctorAddress string `json:"doctor_address"`
}

type reencryptResponse struct {
	ReEncryptedKey string `json:"reencrypted_key"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type issueRequest struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type issueResponse struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type assignRequest struct {
	DoctorAddress string `json:"doctor_address"`
	Department    string `json:"department"`
	Institution   string `json:"institution"`
	ValidUntil    string `json:"valid_until"`
}

type attributesResponse struct {
	DoctorAddress string `json:"doctor_address"`
	Department    string `json:"department"`
	Institution   string `json:"institution"`
	ValidUntil    string `json:"valid_until"`
}

type storeRecoveryKeyRequest struct {
	EmergencySummary string `json:"emergency_summary"`
	ContentKey       string `json:"content_key"`
}

type storeRecoveryKeyResponse struct {
	RecoveryKey string `json:"recovery_key"`
}

type getRecoveryKeyResponse struct {
	RecoveryKey string `json:"recovery_key"`
}

type reencryptRecoveryKeyRequest struct {
	EmergencySummary string `json:"emergency_summary"`
	Recipient        string `json:"recipient"`
}

func NewHandler(store RecordStore, ipfsClient IPFSClient, vaultClient VaultClient, roles RoleProvider, sessions *SessionStore, emergency EmergencyResolver) *Handler {
	return &Handler{
		store:     store,
		ipfs:      ipfsClient,
		vault:     vaultClient,
		roles:     roles,
		sessions:  sessions,
		emergency: emergency,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHandler := h.requireAuth

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/records/upload":
		authHandler(h.handleUpload)(w, r)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/records/") && strings.HasSuffix(r.URL.Path, "/key"):
		authHandler(h.handleStoreKey)(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/records/") && strings.HasSuffix(r.URL.Path, "/key"):
		authHandler(h.handleGetKey)(w, r)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/records/") && strings.HasSuffix(r.URL.Path, "/delegate"):
		authHandler(h.handleDelegate)(w, r)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/records/") && strings.HasSuffix(r.URL.Path, "/reencrypt"):
		authHandler(h.handleReEncrypt)(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/records/") && strings.HasSuffix(r.URL.Path, "/download"):
		authHandler(h.handleDownload)(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/attributes/issue":
		authHandler(h.handleIssueAttributes)(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/attributes/assign":
		authHandler(h.handleAssignAttributes)(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/attributes/registry/"):
		authHandler(h.handleRegistryLookup)(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/attributes/"):
		authHandler(h.handleGetAttributes)(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/emergency/recovery-key":
		authHandler(h.handleStoreRecoveryKey)(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/emergency/recovery-key":
		authHandler(h.handleGetRecoveryKey)(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/emergency/recovery-key/re-encrypt":
		authHandler(h.handleReencryptRecoveryKey)(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/identities":
		authHandler(h.handleProvisionIdentity)(w, r)
	case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/api/v1/identities/"):
		authHandler(h.handleDeactivateIdentity)(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/identities":
		authHandler(h.handleLookupIdentity)(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) requireAdmin(ctx context.Context) error {
	caller, ok := ctx.Value(AuthCtxKey).(string)
	if !ok || caller == "" {
		return errors.New("unauthenticated")
	}
	role, err := h.roles.RoleOf(ctx, caller)
	if err != nil {
		return errors.New("unregistered user")
	}
	if role != RoleAdmin {
		return errors.New("admin role required")
	}
	return nil
}

func (h *Handler) requireAuth(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Address")
		token := r.Header.Get("X-Session-Token")

		if user == "" || token == "" {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing X-User-Address or X-Session-Token"})
			return
		}

		if !h.sessions.IsValid(r.Context(), user, token) {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid or expired session"})
			return
		}

		ctx := context.WithValue(r.Context(), AuthCtxKey, user)
		next(w, r.WithContext(ctx))
	}
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid multipart form"})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing file field"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "read file: " + err.Error()})
		return
	}

	cid, err := h.ipfs.Add(r.Context(), data)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "ipfs upload: " + err.Error()})
		return
	}

	hash := sha256.Sum256(data)
	contentHash := "0x" + hex.EncodeToString(hash[:])

	log.Printf("record uploaded: cid=%s hash=%s size=%d", cid, contentHash, len(data))

	writeJSON(w, http.StatusOK, uploadResponse{
		IPFSCid:     cid,
		ContentHash: contentHash,
	})
}

func (h *Handler) handleStoreKey(w http.ResponseWriter, r *http.Request) {
	var req keyStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}

	if req.RecordAddress == "" || req.IPFSCid == "" || req.EncryptedKey == "" || req.PatientAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing fields"})
		return
	}

	patientAddr := strings.ToLower(req.PatientAddr)

	rec := &db.EncryptedKeyRecord{
		RecordAddress:  strings.ToLower(req.RecordAddress),
		IPFSCid:        req.IPFSCid,
		ContentHash:    req.ContentHash,
		PatientAddress: patientAddr,
	}

	if req.LockboxSalt != "" || req.PatientPublicKey != "" {
		// PRE lockbox path: the client seals the content key and sends the
		// canonical delegator ciphertext C1, the HKDF salt and the patient
		// public key. The backend validates and stores them; it never sees the
		// patient secret key or the plaintext content key.
		if req.LockboxSalt == "" || req.PatientPublicKey == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "lockbox_salt and patient_public_key must be provided together"})
			return
		}

		// Only the record's owner may install their lockbox.
		caller, err := ExtractUser(r.Context())
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
			return
		}
		if strings.ToLower(caller) != patientAddr {
			writeJSON(w, http.StatusForbidden, errorResponse{Error: "only the patient may store their lockbox"})
			return
		}

		lockbox, err := hex.DecodeString(strings.TrimPrefix(req.EncryptedKey, "0x"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid encrypted_key hex"})
			return
		}
		if _, err := pre.ParseDelegatorCiphertext(lockbox); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid lockbox: " + err.Error()})
			return
		}

		salt, err := hex.DecodeString(strings.TrimPrefix(req.LockboxSalt, "0x"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid lockbox_salt hex"})
			return
		}

		pubKey, err := hex.DecodeString(strings.TrimPrefix(req.PatientPublicKey, "0x"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid patient_public_key hex"})
			return
		}
		if _, err := pre.ParsePublicKey(pubKey); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid patient_public_key: " + err.Error()})
			return
		}

		rec.EncryptedKey = lockbox
		rec.LockboxSalt = salt
		rec.PatientPublicKey = pubKey
	} else {
		// Legacy path: the client sends a raw content key which the backend
		// wraps with the patient's Vault transit key.
		if err := h.vault.EnsurePatientKey(r.Context(), patientAddr); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault ensure key: " + err.Error()})
			return
		}

		rawKey, err := hex.DecodeString(strings.TrimPrefix(req.EncryptedKey, "0x"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid encrypted_key hex"})
			return
		}

		wrappedKey, err := h.vault.WrapKeyWithPatientKey(r.Context(), patientAddr, rawKey)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault wrap: " + err.Error()})
			return
		}
		rec.EncryptedKey = wrappedKey
	}

	if err := h.store.InsertEncryptedKey(r.Context(), rec); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "stored"})
}

// handleDelegate stores a patient-supplied PRE re-encryption key rk_{A->B}.
// The patient derives it from their PRE secret key; the backend only receives
// the transform key and can never open a lockbox.
func (h *Handler) handleDelegate(w http.ResponseWriter, r *http.Request) {
	recordAddr := strings.ToLower(extractRecordAddress(r.URL.Path, "/delegate"))
	if recordAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid record address"})
		return
	}

	var req delegateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}
	if !addrRe.MatchString(req.PatientAddress) || !addrRe.MatchString(req.DoctorAddress) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid patient_address or doctor_address"})
		return
	}
	if req.ReEncryptionKey == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing re_encryption_key"})
		return
	}

	rkBytes, err := hex.DecodeString(strings.TrimPrefix(req.ReEncryptionKey, "0x"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid re_encryption_key hex"})
		return
	}
	if _, err := pre.ParseReEncryptionKey(rkBytes); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid re_encryption_key: " + err.Error()})
		return
	}

	// The caller must be the patient that owns the record's lockbox.
	caller, err := ExtractUser(r.Context())
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	patientAddr := strings.ToLower(req.PatientAddress)
	if strings.ToLower(caller) != patientAddr {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "only the patient may delegate"})
		return
	}

	encKeyRec, err := h.store.GetEncryptedKey(r.Context(), recordAddr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "record key not found"})
		return
	}
	if encKeyRec.PatientAddress != patientAddr {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "patient_address does not own record"})
		return
	}

	if err := h.store.UpsertReEncryptionKey(r.Context(), patientAddr, strings.ToLower(req.DoctorAddress), rkBytes); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "delegated"})
}

func (h *Handler) handleGetKey(w http.ResponseWriter, r *http.Request) {
	recordAddr := strings.ToLower(extractRecordAddress(r.URL.Path, "/key"))
	if recordAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid record address"})
		return
	}

	doctorAddr := strings.ToLower(r.URL.Query().Get("doctor_address"))
	if doctorAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing doctor_address query param"})
		return
	}

	key, err := h.store.GetReEncryptedKey(r.Context(), recordAddr, doctorAddr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "reencrypted key not found"})
		return
	}

	writeJSON(w, http.StatusOK, reencryptResponse{
		ReEncryptedKey: hex.EncodeToString(key),
	})
}

func (h *Handler) handleReEncrypt(w http.ResponseWriter, r *http.Request) {
	recordAddr := extractRecordAddress(r.URL.Path, "/reencrypt")
	if recordAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid record address"})
		return
	}

	var req reencryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}
	if !addrRe.MatchString(req.DoctorAddress) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid doctor_address"})
		return
	}
	doctorAddr := strings.ToLower(req.DoctorAddress)

	encKeyRec, err := h.store.GetEncryptedKey(r.Context(), recordAddr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "record key not found"})
		return
	}
	if len(encKeyRec.LockboxSalt) == 0 || len(encKeyRec.PatientPublicKey) == 0 {
		writeJSON(w, http.StatusConflict, errorResponse{Error: "record predates PRE; patient must re-seal it before delegation"})
		return
	}

	rkBytes, err := h.store.GetReEncryptionKey(r.Context(), encKeyRec.PatientAddress, doctorAddr)
	if err != nil {
		writeJSON(w, http.StatusConflict, errorResponse{Error: "no re-encryption key; patient must delegate first"})
		return
	}

	lockbox, err := pre.ParseDelegatorCiphertext(encKeyRec.EncryptedKey)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "stored lockbox invalid: " + err.Error()})
		return
	}
	rk, err := pre.ParseReEncryptionKey(rkBytes)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "stored re-encryption key invalid: " + err.Error()})
		return
	}

	reCT, err := pre.Reencrypt(lockbox, rk)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "re-encrypt: " + err.Error()})
		return
	}
	reBytes, err := reCT.MarshalBinary()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "encode re-encrypted ciphertext: " + err.Error()})
		return
	}

	if err := h.store.InsertReEncryptedKey(r.Context(), &db.ReEncryptedKeyRecord{
		RecordAddress:  recordAddr,
		DoctorAddress:  doctorAddr,
		ReEncryptedKey: reBytes,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, reencryptResponse{
		ReEncryptedKey: hex.EncodeToString(reBytes),
	})
}

func (h *Handler) handleDownload(w http.ResponseWriter, r *http.Request) {
	recordAddr := extractRecordAddress(r.URL.Path, "/download")
	if recordAddr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid record address"})
		return
	}

	encKeyRec, err := h.store.GetEncryptedKey(r.Context(), recordAddr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "record not found"})
		return
	}

	data, err := h.ipfs.Cat(r.Context(), encKeyRec.IPFSCid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "ipfs: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.bin"`, encKeyRec.IPFSCid))
	w.Write(data)
}

func extractRecordAddress(path, suffix string) string {
	trimmed := strings.TrimPrefix(path, "/api/v1/records/")
	trimmed = strings.TrimSuffix(trimmed, suffix)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	return strings.ToLower(parts[0])
}

var (
	uuidRe = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)
	addrRe = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
)

func (h *Handler) handleIssueAttributes(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r.Context()); err != nil {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
		return
	}

	var req issueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	if req.Kind == "" || req.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "kind and name are required"})
		return
	}

	uuid, err := h.vault.IssueUUID(r.Context(), req.Kind, req.Name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, issueResponse{Kind: req.Kind, Name: req.Name, UUID: uuid})
}

func (h *Handler) handleAssignAttributes(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r.Context()); err != nil {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
		return
	}

	var req assignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	if !addrRe.MatchString(req.DoctorAddress) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid doctor_address"})
		return
	}
	if !uuidRe.MatchString(req.Department) || !uuidRe.MatchString(req.Institution) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "department and institution must be issued 256-bit UUIDs"})
		return
	}
	if req.ValidUntil == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "valid_until is required"})
		return
	}

	attrs := &vault.Attributes{
		Department:  strings.ToLower(req.Department),
		Institution: strings.ToLower(req.Institution),
		ValidUntil:  req.ValidUntil,
	}
	if err := h.vault.StoreAttributes(r.Context(), strings.ToLower(req.DoctorAddress), attrs); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, attributesResponse{
		DoctorAddress: strings.ToLower(req.DoctorAddress),
		Department:    attrs.Department,
		Institution:   attrs.Institution,
		ValidUntil:    attrs.ValidUntil,
	})
}

func (h *Handler) handleGetAttributes(w http.ResponseWriter, r *http.Request) {
	addr := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/v1/attributes/"))
	if !addrRe.MatchString(addr) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid address"})
		return
	}

	caller, _ := r.Context().Value(AuthCtxKey).(string)
	if !strings.EqualFold(caller, addr) {
		if err := h.requireAdmin(r.Context()); err != nil {
			writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
			return
		}
	}

	attrs, err := h.vault.GetAttributes(r.Context(), addr)
	if err != nil {
		if errors.Is(err, vault.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "attributes not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, attributesResponse{
		DoctorAddress: addr,
		Department:    attrs.Department,
		Institution:   attrs.Institution,
		ValidUntil:    attrs.ValidUntil,
	})
}

func (h *Handler) handleRegistryLookup(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/attributes/registry/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "expected /api/v1/attributes/registry/{kind}/{name}"})
		return
	}

	uuid, err := h.vault.ResolveUUID(r.Context(), parts[0], parts[1])
	if err != nil {
		if errors.Is(err, vault.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "uuid not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, issueResponse{Kind: parts[0], Name: parts[1], UUID: uuid})
}

// resolveEmergencyOwner resolves the on-chain owner of an EmergencySummary.
func (h *Handler) resolveEmergencyOwner(ctx context.Context, summaryAddr string) (string, error) {
	if h.emergency == nil {
		return "", errors.New("emergency resolver not configured")
	}
	owner, err := h.emergency.OwnerOf(ctx, strings.ToLower(summaryAddr))
	if err != nil {
		return "", err
	}
	return strings.ToLower(owner), nil
}

// handleStoreRecoveryKey escrows the emergency data key: the backend wraps it
// with the patient's Vault transit key, stores it in Vault KV, and returns the
// wrapped blob so the owner can also persist it on-chain.
func (h *Handler) handleStoreRecoveryKey(w http.ResponseWriter, r *http.Request) {
	var req storeRecoveryKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}
	if !addrRe.MatchString(req.EmergencySummary) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid emergency_summary"})
		return
	}
	if req.ContentKey == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "content_key is required"})
		return
	}

	caller, err := ExtractUser(r.Context())
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	caller = strings.ToLower(caller)

	owner, err := h.resolveEmergencyOwner(r.Context(), req.EmergencySummary)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "resolve owner: " + err.Error()})
		return
	}
	if caller != owner {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "only the summary owner may escrow the recovery key"})
		return
	}

	rawKey, err := hex.DecodeString(strings.TrimPrefix(req.ContentKey, "0x"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid content_key hex"})
		return
	}

	if err := h.vault.EnsurePatientKey(r.Context(), owner); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault ensure key: " + err.Error()})
		return
	}

	wrappedKey, err := h.vault.WrapKeyWithPatientKey(r.Context(), owner, rawKey)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault wrap: " + err.Error()})
		return
	}

	if err := h.vault.StoreRecoveryKey(r.Context(), strings.ToLower(req.EmergencySummary), wrappedKey); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault store: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, storeRecoveryKeyResponse{
		RecoveryKey: hex.EncodeToString(wrappedKey),
	})
}

// handleGetRecoveryKey returns the escrowed, wrapped data key to the owner.
func (h *Handler) handleGetRecoveryKey(w http.ResponseWriter, r *http.Request) {
	summaryAddr := strings.ToLower(r.URL.Query().Get("emergency_summary"))
	if !addrRe.MatchString(summaryAddr) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "expected emergency_summary query param"})
		return
	}

	caller, err := ExtractUser(r.Context())
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	caller = strings.ToLower(caller)

	owner, err := h.resolveEmergencyOwner(r.Context(), summaryAddr)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "resolve owner: " + err.Error()})
		return
	}
	if caller != owner {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "only the summary owner may read the recovery key"})
		return
	}

	blob, err := h.vault.GetRecoveryKey(r.Context(), summaryAddr)
	if err != nil {
		if errors.Is(err, vault.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "recovery key not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault read: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, getRecoveryKeyResponse{
		RecoveryKey: hex.EncodeToString(blob),
	})
}

// handleReencryptRecoveryKey yields the escaped data key wrapped for a recipient.
// Authorized for the summary owner or any caller with an active on-chain grant
// (RiskAssessment emergencyGrant or an approved access request).
func (h *Handler) handleReencryptRecoveryKey(w http.ResponseWriter, r *http.Request) {
	var req reencryptRecoveryKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}
	if !addrRe.MatchString(req.EmergencySummary) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid emergency_summary"})
		return
	}
	if !addrRe.MatchString(req.Recipient) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid recipient"})
		return
	}

	caller, err := ExtractUser(r.Context())
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	caller = strings.ToLower(caller)
	summaryAddr := strings.ToLower(req.EmergencySummary)

	owner, err := h.resolveEmergencyOwner(r.Context(), summaryAddr)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "resolve owner: " + err.Error()})
		return
	}

	authorized := caller == owner
	if !authorized {
		if h.emergency == nil {
			writeJSON(w, http.StatusForbidden, errorResponse{Error: "no emergency resolver configured"})
			return
		}
		ok, err := h.emergency.HasAccess(r.Context(), summaryAddr, caller)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "check access: " + err.Error()})
			return
		}
		authorized = ok
	}
	if !authorized {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "owner or active grant required to re-encrypt"})
		return
	}

	blob, err := h.vault.GetRecoveryKey(r.Context(), summaryAddr)
	if err != nil {
		if errors.Is(err, vault.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "recovery key not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault read: " + err.Error()})
		return
	}

	recipient := strings.ToLower(req.Recipient)
	if err := h.vault.EnsurePatientKey(r.Context(), recipient); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault ensure recipient key: " + err.Error()})
		return
	}

	rekeyed, err := h.vault.Rekey(r.Context(), owner, recipient, blob)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "vault rekey: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, reencryptResponse{
		ReEncryptedKey: hex.EncodeToString(rekeyed),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
