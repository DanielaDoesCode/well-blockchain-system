package pre

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Content key parameters.
const (
	// ContentKeySize is the AES-256 key length in bytes.
	ContentKeySize = 32
	// SaltSize is the length of the per-record KDF salt.
	SaltSize = 32
)

// KDFLabel is the scheme version bound into every derived key. Changing the
// scheme requires changing this, so that keys cannot be reused across
// versions.
const KDFLabel = "well-pre-v1"

// Errors from content key handling.
var (
	ErrDecryptFailed = errors.New("pre: content decryption failed")
)

func hkdfBytes(ikm, salt []byte, info []byte, length int) ([]byte, error) {
	if len(ikm) == 0 {
		return nil, fmt.Errorf("%w: empty input key material", ErrMalformed)
	}
	if length <= 0 {
		return nil, fmt.Errorf("%w: non-positive output length %d", ErrMalformed, length)
	}
	out := make([]byte, length)
	r := hkdf.New(sha256.New, ikm, salt, info)
	if _, err := io.ReadFull(r, out); err != nil {
		return nil, fmt.Errorf("%w: hkdf: %v", ErrMalformed, err)
	}
	return out, nil
}

// NewSalt returns a fresh random KDF salt for one record.
func NewSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(randSource, salt); err != nil {
		return nil, fmt.Errorf("pre: read salt: %w", err)
	}
	return salt, nil
}

// ContentInfo is the HKDF info string for a record. It binds the derived key
// to one record and one purpose, so a lockbox cannot be moved between records
// and a content key cannot be reused for metadata.
func ContentInfo(recordAddress string) []byte {
	return []byte(KDFLabel + "\x00content\x00" + recordAddress)
}

// DeriveContentKey turns a recovered Gt message into an AES-256 key.
//
// The library's own KdfGtToAes256 passes nil for both salt and info, so it
// provides no domain separation at all. This version requires both, and binds
// the record address in.
//
// The recovered message is the only secret input. No salt choice can rescue a
// guessable message.
func DeriveContentKey(message []byte, salt []byte, recordAddress string) ([]byte, error) {
	if len(salt) != SaltSize {
		return nil, fmt.Errorf("%w: salt must be %d bytes, got %d", ErrMalformed, SaltSize, len(salt))
	}
	if recordAddress == "" {
		return nil, fmt.Errorf("%w: empty record address", ErrMalformed)
	}
	return hkdfBytes(message, salt, ContentInfo(recordAddress), ContentKeySize)
}

// DeriveContentKeyFromMessage is DeriveContentKey for a recovered group
// element, serialising it canonically first.
func DeriveContentKeyFromMessage(m interface{ MarshalBinary() ([]byte, error) }, salt []byte, recordAddress string) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("%w: nil message", ErrMalformed)
	}
	raw, err := m.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: encode message: %v", ErrMalformed, err)
	}
	if isZero(raw) {
		return nil, fmt.Errorf("%w: message encoded to all zeros", ErrIdentity)
	}
	return DeriveContentKey(raw, salt, recordAddress)
}

func isZero(b []byte) bool {
	for _, x := range b {
		if x != 0 {
			return false
		}
	}
	return true
}

// EncryptContent seals plaintext with AES-256-GCM under the given key. The
// record address is bound as associated data, so a ciphertext cannot be moved
// between records even if both are readable by the same doctor.
//
// Returns the nonce followed by the sealed bytes.
func EncryptContent(key, plaintext, nonce []byte, recordAddress string) ([]byte, error) {
	aead, err := newContentAEAD(key)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("%w: nonce must be %d bytes, got %d", ErrMalformed, aead.NonceSize(), len(nonce))
	}
	sealed := aead.Seal(nil, nonce, plaintext, []byte(recordAddress))
	return append(append([]byte{}, nonce...), sealed...), nil
}

// DecryptContent opens bytes produced by EncryptContent.
//
// Every failure returns ErrDecryptFailed, deliberately without distinguishing
// a bad tag from a bad key or a bad nonce. The scheme has no CCA security, so
// a substituted lockbox produces a wrong key rather than an error, and the
// only thing that detects that is this tag.
func DecryptContent(key, sealed []byte, recordAddress string) ([]byte, error) {
	aead, err := newContentAEAD(key)
	if err != nil {
		return nil, err
	}
	if len(sealed) < aead.NonceSize() {
		return nil, ErrDecryptFailed
	}
	nonce := sealed[:aead.NonceSize()]
	body := sealed[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, body, []byte(recordAddress))
	if err != nil {
		return nil, ErrDecryptFailed
	}
	return plaintext, nil
}

func newContentAEAD(key []byte) (cipher.AEAD, error) {
	if len(key) != ContentKeySize {
		return nil, fmt.Errorf("%w: content key must be %d bytes, got %d", ErrMalformed, ContentKeySize, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: aes: %v", ErrMalformed, err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: gcm: %v", ErrMalformed, err)
	}
	return aead, nil
}
