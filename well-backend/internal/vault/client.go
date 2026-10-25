package vault

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Attributes are the off-chain provider attributes issued by an admin and
// stored in Vault KV — never on-chain. The doctor's ZK proof is generated
// from these (Vault -> backend API -> proof generator).
type Attributes struct {
	Department  string `json:"department"`
	Institution string `json:"institution"`
	ValidUntil  string `json:"valid_until"`
}

var ErrNotFound = errors.New("not found")

type Client struct {
	cl *vaultapi.Client
}

func New(addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	cl, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	cl.SetToken(token)

	// Enable the transit secrets engine if not already enabled
	c := &Client{cl: cl}
	c.enableTransitEngine()
	c.enableKVEngine()

	return c, nil
}

func (c *Client) enableTransitEngine() {
	mounts, err := c.cl.Sys().ListMounts()
	if err != nil {
		log.Printf("vault: list mounts: %v (transit may not be available)", err)
		return
	}
	if _, ok := mounts["transit/"]; ok {
		return
	}
	if err := c.cl.Sys().Mount("transit", &vaultapi.MountInput{
		Type: "transit",
	}); err != nil {
		log.Printf("vault: enable transit: %v", err)
	}
}

func (c *Client) keyName(patientAddr string) string {
	return fmt.Sprintf("well-patient-%s", patientAddr)
}

func (c *Client) StorePatientKey(ctx context.Context, patientAddr string, privateKeyHex string) error {
	path := fmt.Sprintf("transit/keys/%s", c.keyName(patientAddr))
	data := map[string]any{
		"type":  "ecdsa-p256",
		"derived": false,
	}

	_, err := c.cl.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return fmt.Errorf("vault create key %s: %w", path, err)
	}

	importPath := fmt.Sprintf("transit/keys/%s/import", c.keyName(patientAddr))
	importData := map[string]any{
		"key": privateKeyHex,
	}
	_, err = c.cl.Logical().WriteWithContext(ctx, importPath, importData)
	return err
}

func (c *Client) EnsurePatientKey(ctx context.Context, patientAddr string) error {
	name := c.keyName(patientAddr)
	path := fmt.Sprintf("transit/keys/%s", name)
	secret, err := c.cl.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return fmt.Errorf("check key %s: %w", name, err)
	}
	if secret != nil {
		return nil // already exists
	}
	data := map[string]any{
		"type": "aes256-gcm96",
	}
	_, err = c.cl.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return fmt.Errorf("create key %s: %w", name, err)
	}
	return nil
}

func (c *Client) WrapKeyWithPatientKey(ctx context.Context, patientAddr string, contentKey []byte) ([]byte, error) {
	path := fmt.Sprintf("transit/encrypt/%s", c.keyName(patientAddr))
	data := map[string]any{
		"plaintext": hex.EncodeToString(contentKey),
	}

	secret, err := c.cl.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("vault encrypt: %w", err)
	}

	ct, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return nil, fmt.Errorf("vault: ciphertext not in response")
	}

	return []byte(ct), nil
}

func (c *Client) UnwrapKeyWithPatientKey(ctx context.Context, patientAddr string, ciphertext []byte) ([]byte, error) {
	path := fmt.Sprintf("transit/decrypt/%s", c.keyName(patientAddr))
	data := map[string]any{
		"ciphertext": string(ciphertext),
	}

	secret, err := c.cl.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("vault decrypt: %w", err)
	}

	pt, ok := secret.Data["plaintext"].(string)
	if !ok {
		return nil, fmt.Errorf("vault: plaintext not in response")
	}

	return hex.DecodeString(pt)
}

func (c *Client) WrapKeyWithDoctorKey(ctx context.Context, doctorAddr string, contentKey []byte) ([]byte, error) {
	path := fmt.Sprintf("transit/encrypt/%s", c.keyName(doctorAddr))
	data := map[string]any{
		"plaintext": hex.EncodeToString(contentKey),
	}

	secret, err := c.cl.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("vault encrypt for doctor: %w", err)
	}

	ct, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return nil, fmt.Errorf("vault: ciphertext not in response")
	}

	return []byte(ct), nil
}

func (c *Client) Rekey(ctx context.Context, patientAddr, doctorAddr string, ciphertext []byte) ([]byte, error) {
	pt, err := c.UnwrapKeyWithPatientKey(ctx, patientAddr, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("rekey unwrap: %w", err)
	}

	ct, err := c.WrapKeyWithDoctorKey(ctx, doctorAddr, pt)
	if err != nil {
		return nil, fmt.Errorf("rekey wrap: %w", err)
	}

	return ct, nil
}

func (c *Client) enableKVEngine() {
	mounts, err := c.cl.Sys().ListMounts()
	if err != nil {
		log.Printf("vault: list mounts: %v (kv may not be available)", err)
		return
	}
	if _, ok := mounts["secret/"]; ok {
		return
	}
	if err := c.cl.Sys().Mount("secret", &vaultapi.MountInput{
		Type: "kv-v2",
	}); err != nil {
		log.Printf("vault: enable kv-v2: %v", err)
	}
}

func (c *Client) attributesPath(addr string) string {
	return fmt.Sprintf("secret/data/well/attributes/%s", strings.ToLower(addr))
}

func (c *Client) registryPath(kind, name string) string {
	return fmt.Sprintf("secret/data/well/registry/%s/%s", kind, name)
}

// StoreAttributes upserts a doctor's issued attributes in Vault KV.
func (c *Client) StoreAttributes(ctx context.Context, addr string, attrs *Attributes) error {
	data := map[string]any{
		"data": map[string]any{
			"department":  attrs.Department,
			"institution": attrs.Institution,
			"valid_until": attrs.ValidUntil,
		},
	}
	if _, err := c.cl.Logical().WriteWithContext(ctx, c.attributesPath(addr), data); err != nil {
		return fmt.Errorf("vault store attributes %s: %w", addr, err)
	}
	return nil
}

// GetAttributes returns a doctor's issued attributes, or ErrNotFound.
func (c *Client) GetAttributes(ctx context.Context, addr string) (*Attributes, error) {
	secret, err := c.cl.Logical().ReadWithContext(ctx, c.attributesPath(addr))
	if err != nil {
		return nil, fmt.Errorf("vault read attributes %s: %w", addr, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, ErrNotFound
	}
	raw, ok := secret.Data["data"].(map[string]any)
	if !ok || raw == nil {
		return nil, ErrNotFound
	}
	attrs := &Attributes{}
	if v, ok := raw["department"].(string); ok {
		attrs.Department = v
	}
	if v, ok := raw["institution"].(string); ok {
		attrs.Institution = v
	}
	if v, ok := raw["valid_until"].(string); ok {
		attrs.ValidUntil = v
	}
	if attrs.Department == "" && attrs.Institution == "" {
		return nil, ErrNotFound
	}
	return attrs, nil
}

// IssueUUID mints a fresh 256-bit UUID for a department/institution and
// records it in the registry so it can be resolved by name.
func (c *Client) IssueUUID(ctx context.Context, kind, name string) (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("vault random uuid: %w", err)
	}
	uuid := "0x" + hex.EncodeToString(b[:])

	data := map[string]any{
		"data": map[string]any{
			"uuid":      uuid,
			"kind":      kind,
			"issued_at": time.Now().UTC().Format(time.RFC3339),
		},
	}
	if _, err := c.cl.Logical().WriteWithContext(ctx, c.registryPath(kind, name), data); err != nil {
		return "", fmt.Errorf("vault store registry %s/%s: %w", kind, name, err)
	}
	return uuid, nil
}

// ResolveUUID looks up an issued UUID by kind/name, or ErrNotFound.
func (c *Client) ResolveUUID(ctx context.Context, kind, name string) (string, error) {
	secret, err := c.cl.Logical().ReadWithContext(ctx, c.registryPath(kind, name))
	if err != nil {
		return "", fmt.Errorf("vault read registry %s/%s: %w", kind, name, err)
	}
	if secret == nil || secret.Data == nil {
		return "", ErrNotFound
	}
	raw, ok := secret.Data["data"].(map[string]any)
	if !ok || raw == nil {
		return "", ErrNotFound
	}
	uuid, ok := raw["uuid"].(string)
	if !ok || uuid == "" {
		return "", ErrNotFound
	}
	return uuid, nil
}

func (c *Client) emergencyPath(summaryAddr string) string {
	return fmt.Sprintf("secret/data/well/emergency/%s", strings.ToLower(summaryAddr))
}

// ContentKey is the escrowed, wrapped emergency data key stored in Vault KV.
// The blob is produced by the transit engine (WrapKeyWithPatientKey), so Vault
// is the only party able to unwrap it and re-encrypt for a recipient.
type ContentKey struct {
	RecoveryKey string `json:"recovery_key"`
}

// StoreRecoveryKey upserts the wrapped emergency data key in Vault KV.
func (c *Client) StoreRecoveryKey(ctx context.Context, summaryAddr string, recoveryKey []byte) error {
	data := map[string]any{
		"data": map[string]any{
			"recovery_key": hex.EncodeToString(recoveryKey),
			"stored_at":    time.Now().UTC().Format(time.RFC3339),
		},
	}
	if _, err := c.cl.Logical().WriteWithContext(ctx, c.emergencyPath(summaryAddr), data); err != nil {
		return fmt.Errorf("vault store recovery key %s: %w", summaryAddr, err)
	}
	return nil
}

// GetRecoveryKey returns the escrowed, wrapped emergency data key, or ErrNotFound.
func (c *Client) GetRecoveryKey(ctx context.Context, summaryAddr string) ([]byte, error) {
	secret, err := c.cl.Logical().ReadWithContext(ctx, c.emergencyPath(summaryAddr))
	if err != nil {
		return nil, fmt.Errorf("vault read recovery key %s: %w", summaryAddr, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, ErrNotFound
	}
	raw, ok := secret.Data["data"].(map[string]any)
	if !ok || raw == nil {
		return nil, ErrNotFound
	}
	recoveryKey, ok := raw["recovery_key"].(string)
	if !ok || recoveryKey == "" {
		return nil, ErrNotFound
	}
	return hex.DecodeString(recoveryKey)
}
