package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Thesis/well-blockchain/well-backend/internal/db"
	"Thesis/well-blockchain/well-backend/internal/pre"
	"Thesis/well-blockchain/well-backend/internal/vault"
)

type mockIPFS struct{}

func (m *mockIPFS) Add(ctx context.Context, data []byte) (string, error) {
	return "QmTest123", nil
}

func (m *mockIPFS) Cat(ctx context.Context, cid string) ([]byte, error) {
	return []byte("encrypted-blob-data"), nil
}

type mockStore struct {
	encKey     *db.EncryptedKeyRecord
	reKeys     map[string][]byte
	rks        map[string][]byte
	identities map[string]*db.Identity
	lookups    []db.IdentityLookup
}

func (m *mockStore) InsertEncryptedKey(ctx context.Context, r *db.EncryptedKeyRecord) error {
	m.encKey = r
	return nil
}

func (m *mockStore) GetEncryptedKey(ctx context.Context, recordAddress string) (*db.EncryptedKeyRecord, error) {
	if m.encKey == nil {
		return nil, io.EOF
	}
	return m.encKey, nil
}

func (m *mockStore) InsertReEncryptedKey(ctx context.Context, r *db.ReEncryptedKeyRecord) error {
	if m.reKeys == nil {
		m.reKeys = make(map[string][]byte)
	}
	m.reKeys[r.RecordAddress+"|"+r.DoctorAddress] = r.ReEncryptedKey
	return nil
}

func (m *mockStore) GetReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) ([]byte, error) {
	if m.reKeys == nil {
		return nil, io.EOF
	}
	key, ok := m.reKeys[recordAddress+"|"+doctorAddress]
	if !ok {
		return nil, io.EOF
	}
	return key, nil
}

func (m *mockStore) DeleteReEncryptedKey(ctx context.Context, recordAddress, doctorAddress string) error {
	delete(m.reKeys, recordAddress+"|"+doctorAddress)
	return nil
}

func (m *mockStore) UpsertReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string, rk []byte) error {
	if m.rks == nil {
		m.rks = make(map[string][]byte)
	}
	m.rks[patientAddress+"|"+doctorAddress] = rk
	return nil
}

func (m *mockStore) GetReEncryptionKey(ctx context.Context, patientAddress, doctorAddress string) ([]byte, error) {
	if m.rks == nil {
		return nil, io.EOF
	}
	rk, ok := m.rks[patientAddress+"|"+doctorAddress]
	if !ok {
		return nil, io.EOF
	}
	return rk, nil
}

func (m *mockStore) UpsertIdentity(ctx context.Context, e *db.Identity) error {
	if m.identities == nil {
		m.identities = make(map[string]*db.Identity)
	}
	m.identities[e.NationalID] = e
	return nil
}

func (m *mockStore) GetIdentityByNationalID(ctx context.Context, nationalID string) (*db.Identity, error) {
	e, ok := m.identities[nationalID]
	if !ok {
		return nil, db.ErrIdentityNotFound
	}
	return e, nil
}

func (m *mockStore) GetIdentitiesByName(ctx context.Context, name string) ([]*db.Identity, error) {
	out := []*db.Identity{}
	for _, e := range m.identities {
		if e.DisplayName == name && e.Status == "active" {
			out = append(out, e)
		}
	}
	return out, nil
}

func (m *mockStore) DeactivateIdentity(ctx context.Context, nationalID string) error {
	e, ok := m.identities[nationalID]
	if !ok {
		return db.ErrIdentityNotFound
	}
	e.Status = "revoked"
	return nil
}

func (m *mockStore) RecordLookup(ctx context.Context, l *db.IdentityLookup) error {
	m.lookups = append(m.lookups, *l)
	return nil
}

type mockVaultClient struct {
	attrs    map[string]*vault.Attributes
	registry map[string]string
	recovery map[string][]byte
}

func (m *mockVaultClient) StoreRecoveryKey(ctx context.Context, summaryAddr string, recoveryKey []byte) error {
	if m.recovery == nil {
		m.recovery = make(map[string][]byte)
	}
	m.recovery[strings.ToLower(summaryAddr)] = recoveryKey
	return nil
}

func (m *mockVaultClient) GetRecoveryKey(ctx context.Context, summaryAddr string) ([]byte, error) {
	if m.recovery == nil {
		return nil, vault.ErrNotFound
	}
	key, ok := m.recovery[strings.ToLower(summaryAddr)]
	if !ok {
		return nil, vault.ErrNotFound
	}
	return key, nil
}

func (m *mockVaultClient) EnsurePatientKey(ctx context.Context, patientAddr string) error {
	return nil
}

func (m *mockVaultClient) WrapKeyWithPatientKey(ctx context.Context, patientAddr string, contentKey []byte) ([]byte, error) {
	return contentKey, nil
}

func (m *mockVaultClient) Rekey(ctx context.Context, patientAddr, doctorAddr string, ciphertext []byte) ([]byte, error) {
	return []byte("reencrypted-key-data"), nil
}

func (m *mockVaultClient) StoreAttributes(ctx context.Context, addr string, attrs *vault.Attributes) error {
	if m.attrs == nil {
		m.attrs = make(map[string]*vault.Attributes)
	}
	m.attrs[strings.ToLower(addr)] = attrs
	return nil
}

func (m *mockVaultClient) GetAttributes(ctx context.Context, addr string) (*vault.Attributes, error) {
	if m.attrs == nil {
		return nil, vault.ErrNotFound
	}
	attrs, ok := m.attrs[strings.ToLower(addr)]
	if !ok {
		return nil, vault.ErrNotFound
	}
	return attrs, nil
}

func (m *mockVaultClient) IssueUUID(ctx context.Context, kind, name string) (string, error) {
	uuid := "0x" + strings.Repeat("1", 64)
	if m.registry == nil {
		m.registry = make(map[string]string)
	}
	m.registry[kind+"/"+name] = uuid
	return uuid, nil
}

func (m *mockVaultClient) ResolveUUID(ctx context.Context, kind, name string) (string, error) {
	if m.registry == nil {
		return "", vault.ErrNotFound
	}
	uuid, ok := m.registry[kind+"/"+name]
	if !ok {
		return "", vault.ErrNotFound
	}
	return uuid, nil
}

type mockRoleProvider struct {
	roles map[string]uint8
}

func (m *mockRoleProvider) RoleOf(ctx context.Context, addr string) (uint8, error) {
	if m.roles == nil {
		return 0, errors.New("unregistered user")
	}
	role, ok := m.roles[strings.ToLower(addr)]
	if !ok {
		return 0, errors.New("unregistered user")
	}
	return role, nil
}

type mockEmergencyResolver struct {
	owner        map[string]string
	granted      map[string]bool
}

func (m *mockEmergencyResolver) OwnerOf(ctx context.Context, summaryAddr string) (string, error) {
	if m.owner == nil {
		return "", errors.New("summary not found")
	}
	owner, ok := m.owner[strings.ToLower(summaryAddr)]
	if !ok {
		return "", errors.New("summary not found")
	}
	return owner, nil
}

func (m *mockEmergencyResolver) HasAccess(ctx context.Context, summaryAddr string, caller string) (bool, error) {
	if m.granted == nil {
		return false, nil
	}
	key := strings.ToLower(summaryAddr) + "|" + strings.ToLower(caller)
	return m.granted[key], nil
}

func TestUploadHandler(t *testing.T) {
	ss := NewSessionStore()
	ss.Add("0xdoctor", "valid-token", time.Now().Add(time.Hour))

	handler := NewHandler(&mockStore{}, &mockIPFS{}, nil, nil, ss, &mockEmergencyResolver{})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("file", "test.bin")
	fw.Write([]byte("encrypted-medical-record-data"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/records/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-User-Address", "0xdoctor")
	req.Header.Set("X-Session-Token", "valid-token")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.IPFSCid != "QmTest123" {
		t.Errorf("expected cid QmTest123, got %s", result.IPFSCid)
	}

	if result.ContentHash == "" {
		t.Error("expected content hash")
	}
}

func TestUploadHandlerUnauthenticated(t *testing.T) {
	ss := NewSessionStore()

	handler := NewHandler(&mockStore{}, &mockIPFS{}, nil, nil, ss, &mockEmergencyResolver{})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("file", "test.bin")
	fw.Write([]byte("data"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/records/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestStoreAndGetKey(t *testing.T) {
	patientAddr := "0x" + strings.Repeat("d", 40)
	doctorAddr := "0x" + strings.Repeat("e", 40)
	recordAddr := "0x" + strings.Repeat("c", 40)

	ss := NewSessionStore()
	ss.Add(patientAddr, "patient-token", time.Now().Add(time.Hour))
	ss.Add(doctorAddr, "doctor-token", time.Now().Add(time.Hour))

	patientPK, patientSK, err := pre.KeyGen()
	if err != nil {
		t.Fatalf("patient keygen: %v", err)
	}
	doctorPK, doctorSK, err := pre.KeyGen()
	if err != nil {
		t.Fatalf("doctor keygen: %v", err)
	}

	lockbox, err := pre.SealRandom(patientPK)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	want, err := pre.DelegateeSelfDecrypt(lockbox, patientSK)
	if err != nil {
		t.Fatalf("self decrypt: %v", err)
	}
	lockboxBytes, err := lockbox.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal lockbox: %v", err)
	}
	pkBytes, err := patientPK.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal patient pk: %v", err)
	}
	salt, err := pre.NewSalt()
	if err != nil {
		t.Fatalf("salt: %v", err)
	}
	rk, err := pre.DeriveReEncryptionKey(patientSK, doctorPK)
	if err != nil {
		t.Fatalf("derive rk: %v", err)
	}
	rkBytes, err := rk.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal rk: %v", err)
	}

	store := &mockStore{}
	handler := NewHandler(store, &mockIPFS{}, &mockVaultClient{}, nil, ss, &mockEmergencyResolver{})

	storeReq := keyStoreRequest{
		RecordAddress:    recordAddr,
		IPFSCid:          "QmCid",
		ContentHash:      "0xhash",
		EncryptedKey:     hex.EncodeToString(lockboxBytes),
		PatientAddr:      patientAddr,
		LockboxSalt:      hex.EncodeToString(salt),
		PatientPublicKey: hex.EncodeToString(pkBytes),
	}
	body, _ := json.Marshal(storeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/records/"+recordAddr+"/key", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Address", patientAddr)
	req.Header.Set("X-Session-Token", "patient-token")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 on store key, got %d: %s", resp.Code, resp.Body.String())
	}

	delegateBody, _ := json.Marshal(delegateRequest{
		PatientAddress:  patientAddr,
		DoctorAddress:   doctorAddr,
		ReEncryptionKey: hex.EncodeToString(rkBytes),
	})
	delegateReq := httptest.NewRequest(http.MethodPost, "/api/v1/records/"+recordAddr+"/delegate", bytes.NewReader(delegateBody))
	delegateReq.Header.Set("Content-Type", "application/json")
	delegateReq.Header.Set("X-User-Address", patientAddr)
	delegateReq.Header.Set("X-Session-Token", "patient-token")
	delegateResp := httptest.NewRecorder()
	handler.ServeHTTP(delegateResp, delegateReq)
	if delegateResp.Code != http.StatusOK {
		t.Fatalf("expected 200 on delegate, got %d: %s", delegateResp.Code, delegateResp.Body.String())
	}

	reencryptBody, _ := json.Marshal(reencryptRequest{DoctorAddress: doctorAddr})
	reencryptReq := httptest.NewRequest(http.MethodPost, "/api/v1/records/"+recordAddr+"/reencrypt", bytes.NewReader(reencryptBody))
	reencryptReq.Header.Set("Content-Type", "application/json")
	reencryptReq.Header.Set("X-User-Address", doctorAddr)
	reencryptReq.Header.Set("X-Session-Token", "doctor-token")
	reencryptResp := httptest.NewRecorder()
	handler.ServeHTTP(reencryptResp, reencryptReq)
	if reencryptResp.Code != http.StatusOK {
		t.Fatalf("expected 200 on reencrypt, got %d: %s", reencryptResp.Code, reencryptResp.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/records/"+recordAddr+"/key?doctor_address="+doctorAddr, nil)
	getReq.Header.Set("X-User-Address", doctorAddr)
	getReq.Header.Set("X-Session-Token", "doctor-token")
	getResp := httptest.NewRecorder()
	handler.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected 200 on get, got %d: %s", getResp.Code, getResp.Body.String())
	}

	var keyResp reencryptResponse
	if err := json.NewDecoder(getResp.Body).Decode(&keyResp); err != nil {
		t.Fatalf("decode key response: %v", err)
	}
	reBytes, err := hex.DecodeString(keyResp.ReEncryptedKey)
	if err != nil {
		t.Fatalf("decode re-encrypted key hex: %v", err)
	}
	reCT, err := pre.ParseReencryptedCiphertext(reBytes)
	if err != nil {
		t.Fatalf("parse re-encrypted ciphertext: %v", err)
	}
	got, err := pre.DelegateeDecrypt(reCT, doctorSK)
	if err != nil {
		t.Fatalf("delegatee decrypt: %v", err)
	}

	wantBytes, err := want.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}
	gotBytes, err := got.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	if !bytes.Equal(wantBytes, gotBytes) {
		t.Fatal("doctor did not recover the sealed message through the server PRE flow")
	}
}

func TestReEncryptWithoutDelegation(t *testing.T) {
	patientAddr := "0x" + strings.Repeat("d", 40)
	doctorAddr := "0x" + strings.Repeat("e", 40)
	recordAddr := "0x" + strings.Repeat("c", 40)

	ss := NewSessionStore()
	ss.Add(doctorAddr, "doctor-token", time.Now().Add(time.Hour))

	patientPK, _, err := pre.KeyGen()
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	lockbox, err := pre.SealRandom(patientPK)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	lockboxBytes, _ := lockbox.MarshalBinary()
	pkBytes, _ := patientPK.MarshalBinary()
	salt, _ := pre.NewSalt()

	store := &mockStore{
		encKey: &db.EncryptedKeyRecord{
			RecordAddress:    recordAddr,
			IPFSCid:          "QmCid",
			ContentHash:      "0xhash",
			EncryptedKey:     lockboxBytes,
			PatientAddress:   patientAddr,
			LockboxSalt:      salt,
			PatientPublicKey: pkBytes,
		},
	}
	handler := NewHandler(store, &mockIPFS{}, &mockVaultClient{}, nil, ss, &mockEmergencyResolver{})

	body, _ := json.Marshal(reencryptRequest{DoctorAddress: doctorAddr})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/records/"+recordAddr+"/reencrypt", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Address", doctorAddr)
	req.Header.Set("X-Session-Token", "doctor-token")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusConflict {
		t.Fatalf("expected 409 without a delegation, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestDownload(t *testing.T) {
	ss := NewSessionStore()
	ss.Add("0xdoctor", "valid-token", time.Now().Add(time.Hour))

	store := &mockStore{
		encKey: &db.EncryptedKeyRecord{
			RecordAddress:  "0xabc",
			IPFSCid:        "QmCid",
			ContentHash:    "0xhash",
			PatientAddress: "0xpatient",
		},
	}
	handler := NewHandler(store, &mockIPFS{}, nil, nil, ss, &mockEmergencyResolver{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/records/0xabc/download", nil)
	req.Header.Set("X-User-Address", "0xdoctor")
	req.Header.Set("X-Session-Token", "valid-token")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "encrypted-blob-data" {
		t.Errorf("unexpected body: got %s", string(body))
	}
}

func TestExtractRecordAddress(t *testing.T) {
	tests := []struct {
		path   string
		suffix string
		want   string
	}{
		{"/api/v1/records/0xabc/key", "/key", "0xabc"},
		{"/api/v1/records/0xabc/download", "/download", "0xabc"},
		{"/api/v1/records/0xabc/reencrypt", "/reencrypt", "0xabc"},
		{"/api/v1/records/0xABC/key", "/key", "0xabc"},
		{"/api/v1/records//key", "/key", ""},
		{"", "/key", ""},
	}

	for _, tt := range tests {
		got := extractRecordAddress(tt.path, tt.suffix)
		if got != tt.want {
			t.Errorf("extractRecordAddress(%q, %q) = %q, want %q", tt.path, tt.suffix, got, tt.want)
		}
	}
}

const (
	testAdmin  = "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	testDoctor = "0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"
	testPatient = "0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
)

func newAttrsTestHandler(roles map[string]uint8) *Handler {
	ss := NewSessionStore()
	ss.Add(testAdmin, "admin-token", time.Now().Add(time.Hour))
	ss.Add(testDoctor, "doctor-token", time.Now().Add(time.Hour))
	ss.Add(testPatient, "patient-token", time.Now().Add(time.Hour))
	ss.Add("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", "cc-token", time.Now().Add(time.Hour))
	normalized := make(map[string]uint8, len(roles))
	for k, v := range roles {
		normalized[strings.ToLower(k)] = v
	}
	return NewHandler(&mockStore{}, &mockIPFS{}, &mockVaultClient{}, &mockRoleProvider{roles: normalized}, ss, &mockEmergencyResolver{})
}

func authenticatedRequest(method, path, body string, user, token string) *http.Request {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Address", user)
	req.Header.Set("X-Session-Token", token)
	return req
}

func TestIssueAttributes(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin, testDoctor: RoleDoctor})

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/attributes/issue",
		`{"kind":"department","name":"cardiology"}`, testAdmin, "admin-token"))

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.Code, resp.Body.String())
	}
	var out issueResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !uuidRe.MatchString(out.UUID) {
		t.Errorf("expected 256-bit hex uuid, got %q", out.UUID)
	}

	lookup := httptest.NewRecorder()
	handler.ServeHTTP(lookup, authenticatedRequest(http.MethodGet, "/api/v1/attributes/registry/department/cardiology",
		"", testAdmin, "admin-token"))
	if lookup.Code != http.StatusOK {
		t.Fatalf("expected 200 on lookup, got %d: %s", lookup.Code, lookup.Body.String())
	}
	var lu issueResponse
	if err := json.NewDecoder(lookup.Body).Decode(&lu); err != nil {
		t.Fatalf("decode lookup: %v", err)
	}
	if lu.UUID != out.UUID {
		t.Errorf("lookup uuid %q != minted uuid %q", lu.UUID, out.UUID)
	}
}

func TestIssueAttributesNonAdmin(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin, testDoctor: RoleDoctor})

	for _, tc := range []struct {
		name string
		user string
		token string
	}{
		{"doctor", testDoctor, "doctor-token"},
		{"patient", testPatient, "patient-token"},
		{"unregistered", "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", "cc-token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/attributes/issue",
				`{"kind":"department","name":"cardiology"}`, tc.user, tc.token))
			if resp.Code != http.StatusForbidden {
				t.Fatalf("expected 403, got %d: %s", resp.Code, resp.Body.String())
			}
		})
	}
}

func TestAssignAndGetAttributes(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin, testDoctor: RoleDoctor})
	deptUUID := "0x" + strings.Repeat("a", 64)
	instUUID := "0x" + strings.Repeat("b", 64)

	body := `{"doctor_address":"` + testDoctor + `","department":"` + deptUUID +
		`","institution":"` + instUUID + `","valid_until":"1810000000"}`
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/attributes/assign", body, testAdmin, "admin-token"))
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	got := httptest.NewRecorder()
	handler.ServeHTTP(got, authenticatedRequest(http.MethodGet, "/api/v1/attributes/"+testDoctor, "", testDoctor, "doctor-token"))
	if got.Code != http.StatusOK {
		t.Fatalf("expected 200 on owner get, got %d: %s", got.Code, got.Body.String())
	}
	var out attributesResponse
	if err := json.NewDecoder(got.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Department != deptUUID || out.Institution != instUUID || out.ValidUntil != "1810000000" {
		t.Errorf("unexpected attributes: %+v", out)
	}

	adminGet := httptest.NewRecorder()
	handler.ServeHTTP(adminGet, authenticatedRequest(http.MethodGet, "/api/v1/attributes/"+testDoctor, "", testAdmin, "admin-token"))
	if adminGet.Code != http.StatusOK {
		t.Fatalf("expected 200 on admin get, got %d: %s", adminGet.Code, adminGet.Body.String())
	}

	patientGet := httptest.NewRecorder()
	handler.ServeHTTP(patientGet, authenticatedRequest(http.MethodGet, "/api/v1/attributes/"+testDoctor, "", testPatient, "patient-token"))
	if patientGet.Code != http.StatusForbidden {
		t.Fatalf("expected 403 on non-owner non-admin get, got %d: %s", patientGet.Code, patientGet.Body.String())
	}
}

func TestAssignAttributesNonAdmin(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin, testDoctor: RoleDoctor})
	deptUUID := "0x" + strings.Repeat("a", 64)
	instUUID := "0x" + strings.Repeat("b", 64)

	body := `{"doctor_address":"` + testDoctor + `","department":"` + deptUUID +
		`","institution":"` + instUUID + `","valid_until":"1810000000"}`
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/attributes/assign", body, testDoctor, "doctor-token"))
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestAssignAttributesInvalid(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin})

	for _, tc := range []struct {
		name string
		body string
	}{
		{"bad address", `{"doctor_address":"0x123","department":"` + "0x" + strings.Repeat("a", 64) + `","institution":"` + "0x" + strings.Repeat("b", 64) + `","valid_until":"1810000000"}`},
		{"bad department", `{"doctor_address":"` + testDoctor + `","department":"not-a-uuid","institution":"` + "0x" + strings.Repeat("b", 64) + `","valid_until":"1810000000"}`},
		{"missing valid_until", `{"doctor_address":"` + testDoctor + `","department":"` + "0x" + strings.Repeat("a", 64) + `","institution":"` + "0x" + strings.Repeat("b", 64) + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/attributes/assign", tc.body, testAdmin, "admin-token"))
			if resp.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", resp.Code, resp.Body.String())
			}
		})
	}
}

func TestGetAttributesNotFound(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin, testDoctor: RoleDoctor})
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet, "/api/v1/attributes/"+testDoctor, "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestAttributesRequireAuth(t *testing.T) {
	handler := newAttrsTestHandler(map[string]uint8{testAdmin: RoleAdmin})
	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/attributes/issue"},
		{http.MethodPost, "/api/v1/attributes/assign"},
		{http.MethodGet, "/api/v1/attributes/" + testDoctor},
		{http.MethodGet, "/api/v1/attributes/registry/department/cardiology"},
	} {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: expected 401, got %d", tc.method, tc.path, resp.Code)
		}
	}
}

const testSummaryAddr = "0x9999999999999999999999999999999999999999"
const testGrantedCaller = "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"

func newEmergencyHandler() (*Handler, *mockVaultClient) {
	ss := NewSessionStore()
	ss.Add(testPatient, "patient-token", time.Now().Add(time.Hour))
	ss.Add(testGrantedCaller, "granted-token", time.Now().Add(time.Hour))
	ss.Add("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", "doctor-token", time.Now().Add(time.Hour))

	resolver := &mockEmergencyResolver{
		owner: map[string]string{
			testSummaryAddr: testPatient,
		},
		granted: map[string]bool{
			strings.ToLower(testSummaryAddr) + "|" + strings.ToLower(testGrantedCaller): true,
		},
	}
	vault := &mockVaultClient{}
	handler := NewHandler(&mockStore{}, &mockIPFS{}, vault, nil, ss, resolver)
	return handler, vault
}

func TestStoreRecoveryKey(t *testing.T) {
	handler, vault := newEmergencyHandler()

	contentKey := hex.EncodeToString([]byte("raw-content-key-12345"))
	body, _ := json.Marshal(storeRecoveryKeyRequest{
		EmergencySummary: testSummaryAddr,
		ContentKey:       contentKey,
	})

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/emergency/recovery-key",
		string(body), testPatient, "patient-token"))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var out storeRecoveryKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// mock WrapKeyWithPatientKey returns rawKey unchanged
	if out.RecoveryKey != contentKey {
		t.Errorf("expected recovery_key == contentKey %s, got %s", contentKey, out.RecoveryKey)
	}

	// Verify Vault got the blob
	blob, err := vault.GetRecoveryKey(nil, testSummaryAddr)
	if err != nil {
		t.Fatalf("vault GetRecoveryKey: %v", err)
	}
	expected, _ := hex.DecodeString(contentKey)
	if len(blob) != len(expected) {
		t.Errorf("vault blob length mismatch: %d vs %d", len(blob), len(expected))
	}
}

func TestStoreRecoveryKeyForbidden(t *testing.T) {
	handler, _ := newEmergencyHandler()

	body, _ := json.Marshal(storeRecoveryKeyRequest{
		EmergencySummary: testSummaryAddr,
		ContentKey:       hex.EncodeToString([]byte("key")),
	})

	// Doctor (non-owner)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/emergency/recovery-key",
		string(body), "0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", "doctor-token"))
	if resp.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-owner, got %d", resp.Code)
	}

	// Unauthenticated
	resp2 := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/emergency/recovery-key", strings.NewReader(string(body)))
	handler.ServeHTTP(resp2, req)
	if resp2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated, got %d", resp2.Code)
	}
}

func TestGetRecoveryKey(t *testing.T) {
	handler, vault := newEmergencyHandler()
	vault.recovery = map[string][]byte{
		testSummaryAddr: []byte("stored-blob-bytes"),
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/emergency/recovery-key?emergency_summary="+testSummaryAddr,
		"", testPatient, "patient-token"))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var out getRecoveryKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.RecoveryKey != hex.EncodeToString(vault.recovery[testSummaryAddr]) {
		t.Errorf("unexpected recovery_key")
	}
}

func TestGetRecoveryKeyNotFound(t *testing.T) {
	handler, _ := newEmergencyHandler()

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/emergency/recovery-key?emergency_summary="+testSummaryAddr,
		"", testPatient, "patient-token"))
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}

func TestReencryptRecoveryKey(t *testing.T) {
	handler, vault := newEmergencyHandler()
	vault.recovery = map[string][]byte{
		testSummaryAddr: []byte("stored-blob-bytes"),
	}
	recipient := "0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"

	body, _ := json.Marshal(reencryptRecoveryKeyRequest{
		EmergencySummary: testSummaryAddr,
		Recipient:        recipient,
	})

	// Owner should succeed
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/emergency/recovery-key/re-encrypt",
		string(body), testPatient, "patient-token"))
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 for owner, got %d: %s", resp.Code, resp.Body.String())
	}
	var out reencryptResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.ReEncryptedKey != hex.EncodeToString([]byte("reencrypted-key-data")) {
		t.Errorf("unexpected reencrypted key: %s", out.ReEncryptedKey)
	}

	// Granted caller should succeed
	resp2 := httptest.NewRecorder()
	handler.ServeHTTP(resp2, authenticatedRequest(http.MethodPost, "/api/v1/emergency/recovery-key/re-encrypt",
		string(body), testGrantedCaller, "granted-token"))
	if resp2.Code != http.StatusOK {
		t.Fatalf("expected 200 for granted caller, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

func TestReencryptRecoveryKeyForbidden(t *testing.T) {
	handler, vault := newEmergencyHandler()
	vault.recovery = map[string][]byte{
		testSummaryAddr: []byte("stored-blob-bytes"),
	}

	body, _ := json.Marshal(reencryptRecoveryKeyRequest{
		EmergencySummary: testSummaryAddr,
		Recipient:        "0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD",
	})

	// Doctor (non-owner, non-granted)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/emergency/recovery-key/re-encrypt",
		string(body), "0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", "doctor-token"))
	if resp.Code != http.StatusForbidden {
		t.Errorf("expected 403 for unauthorized caller, got %d", resp.Code)
	}
}

func TestEmergencyRequireAuth(t *testing.T) {
	handler, _ := newEmergencyHandler()

	for _, tc := range []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodPost, "/api/v1/emergency/recovery-key", `{}`},
		{http.MethodGet, "/api/v1/emergency/recovery-key?emergency_summary=" + testSummaryAddr, ""},
		{http.MethodPost, "/api/v1/emergency/recovery-key/re-encrypt", `{}`},
	} {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: expected 401, got %d", tc.method, tc.path, resp.Code)
		}
	}
}

func newIdentityTestHandler(store *mockStore, roles map[string]uint8) *Handler {
	ss := NewSessionStore()
	ss.Add(testAdmin, "admin-token", time.Now().Add(time.Hour))
	ss.Add(testDoctor, "doctor-token", time.Now().Add(time.Hour))
	ss.Add(testPatient, "patient-token", time.Now().Add(time.Hour))
	normalized := make(map[string]uint8, len(roles))
	for k, v := range roles {
		normalized[strings.ToLower(k)] = v
	}
	return NewHandler(store, &mockIPFS{}, &mockVaultClient{}, &mockRoleProvider{roles: normalized}, ss, &mockEmergencyResolver{})
}

func TestIdentityProvisionAndLookup(t *testing.T) {
	store := &mockStore{}
	handler := newIdentityTestHandler(store, map[string]uint8{
		testAdmin:   RoleAdmin,
		testDoctor:  RoleDoctor,
		testPatient: RolePatient,
	})

	// Non-admin cannot provision.
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/identities",
		`{"national_id":"PAT-1","display_name":"Alice","role":"patient","address":"0xaAaAaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"}`,
		testDoctor, "doctor-token"))
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-admin provision, got %d: %s", resp.Code, resp.Body.String())
	}

	// Admin provisions a patient.
	const aliceAddr = "0x3333333333333333333333333333333333333333"
	provisionBody := `{"national_id":"PAT-1","display_name":"Alice Patient","role":"patient",
		"address":"` + aliceAddr + `",
		"commitment":"0x` + strings.Repeat("ab", 32) + `","internal_id":"0x` + strings.Repeat("cd", 32) + `"}`
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPost, "/api/v1/identities",
		provisionBody, testAdmin, "admin-token"))
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201 on provision, got %d: %s", resp.Code, resp.Body.String())
	}

	// Doctor resolves by national ID.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?national_id=PAT-1", "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 on lookup, got %d: %s", resp.Code, resp.Body.String())
	}
	var out struct {
		Identity identityResponse `json:"identity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode lookup: %v", err)
	}
	if out.Identity.NationalID != "PAT-1" || out.Identity.Address != aliceAddr {
		t.Errorf("unexpected identity: %+v", out.Identity)
	}

	// Doctor resolves by name.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?name=Alice%20Patient", "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 on name lookup, got %d: %s", resp.Code, resp.Body.String())
	}

	// Unknown national ID -> 404.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?national_id=UNKNOWN", "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown id, got %d: %s", resp.Code, resp.Body.String())
	}

	// Both query params -> 400.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?national_id=PAT-1&name=Alice", "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for both params, got %d", resp.Code)
	}

	// Patient role cannot resolve.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?national_id=PAT-1", "", testPatient, "patient-token"))
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for patient lookup, got %d", resp.Code)
	}

	// Admin deactivates; lookup then 404s.
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodPatch,
		"/api/v1/identities/PAT-1", "", testAdmin, "admin-token"))
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 on deactivate, got %d: %s", resp.Code, resp.Body.String())
	}
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, authenticatedRequest(http.MethodGet,
		"/api/v1/identities?national_id=PAT-1", "", testDoctor, "doctor-token"))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deactivate, got %d", resp.Code)
	}

	if len(store.lookups) == 0 {
		t.Error("expected every lookup to be audited")
	}
}
