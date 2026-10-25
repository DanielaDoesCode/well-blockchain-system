package pre

import (
	"bytes"
	"errors"
	"testing"
)

func TestContentKeyDeterminism(t *testing.T) {
	salt, err := NewSalt()
	if err != nil {
		t.Fatalf("NewSalt: %v", err)
	}
	msg := bytes.Repeat([]byte{0x2a}, SizeGt)

	k1, err := DeriveContentKey(msg, salt, "0xrecord")
	if err != nil {
		t.Fatalf("DeriveContentKey: %v", err)
	}
	k2, _ := DeriveContentKey(msg, salt, "0xrecord")
	if !bytes.Equal(k1, k2) {
		t.Error("content key is not deterministic for fixed inputs")
	}
	if len(k1) != ContentKeySize {
		t.Errorf("content key length = %d, want %d", len(k1), ContentKeySize)
	}
}

func TestContentKeyDomainSeparation(t *testing.T) {
	salt, _ := NewSalt()
	msg := bytes.Repeat([]byte{0x2a}, SizeGt)

	base, _ := DeriveContentKey(msg, salt, "0xrecord")

	// A different record address must give a different key.
	otherRecord, _ := DeriveContentKey(msg, salt, "0xother")
	if bytes.Equal(base, otherRecord) {
		t.Error("record address is not bound into the key")
	}

	// A different salt must give a different key.
	salt2, _ := NewSalt()
	otherSalt, _ := DeriveContentKey(msg, salt2, "0xrecord")
	if bytes.Equal(base, otherSalt) {
		t.Error("salt is not bound into the key")
	}

	// A different message must give a different key.
	msg2 := bytes.Repeat([]byte{0x2b}, SizeGt)
	otherMsg, _ := DeriveContentKey(msg2, salt, "0xrecord")
	if bytes.Equal(base, otherMsg) {
		t.Error("message is not bound into the key")
	}
}

func TestContentKeyRejectsBadInputs(t *testing.T) {
	msg := bytes.Repeat([]byte{0x2a}, SizeGt)
	salt, _ := NewSalt()

	if _, err := DeriveContentKey(msg, salt[:8], "0xrecord"); err == nil {
		t.Error("short salt accepted")
	}
	if _, err := DeriveContentKey(nil, salt, "0xrecord"); err == nil {
		t.Error("empty message accepted")
	}
	if _, err := DeriveContentKey(msg, salt, ""); err == nil {
		t.Error("empty record address accepted")
	}
}

func TestContentEncryptRoundTrip(t *testing.T) {
	salt, _ := NewSalt()
	key, _ := DeriveContentKey(bytes.Repeat([]byte{1}, SizeGt), salt, "0xrec")

	nonce := make([]byte, 12)
	for i := range nonce {
		nonce[i] = byte(i)
	}
	plaintext := []byte("patient record contents")

	sealed, err := EncryptContent(key, plaintext, nonce, "0xrec")
	if err != nil {
		t.Fatalf("EncryptContent: %v", err)
	}
	if bytes.Contains(sealed, plaintext) {
		t.Error("plaintext appears verbatim in the sealed output")
	}

	got, err := DecryptContent(key, sealed, "0xrec")
	if err != nil {
		t.Fatalf("DecryptContent: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Error("content did not round trip")
	}
}

func TestContentDecryptBindsRecordAddress(t *testing.T) {
	salt, _ := NewSalt()
	key, _ := DeriveContentKey(bytes.Repeat([]byte{1}, SizeGt), salt, "0xrec")
	nonce := bytes.Repeat([]byte{7}, 12)
	sealed, _ := EncryptContent(key, []byte("secret"), nonce, "0xrec")

	// Same key, same bytes, different associated data must fail.
	if _, err := DecryptContent(key, sealed, "0xother"); !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("cross-record decryption error = %v, want ErrDecryptFailed", err)
	}
}

func TestContentDecryptWrongKeyFails(t *testing.T) {
	salt, _ := NewSalt()
	key, _ := DeriveContentKey(bytes.Repeat([]byte{1}, SizeGt), salt, "0xrec")
	wrong, _ := DeriveContentKey(bytes.Repeat([]byte{2}, SizeGt), salt, "0xrec")

	nonce := bytes.Repeat([]byte{7}, 12)
	sealed, _ := EncryptContent(key, []byte("secret"), nonce, "0xrec")

	if _, err := DecryptContent(wrong, sealed, "0xrec"); !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("wrong-key decryption error = %v, want ErrDecryptFailed", err)
	}
}

func TestContentDecryptTamperedFails(t *testing.T) {
	salt, _ := NewSalt()
	key, _ := DeriveContentKey(bytes.Repeat([]byte{1}, SizeGt), salt, "0xrec")
	nonce := bytes.Repeat([]byte{7}, 12)
	sealed, _ := EncryptContent(key, []byte("secret"), nonce, "0xrec")

	tampered := append([]byte{}, sealed...)
	tampered[len(tampered)-1] ^= 0x01
	if _, err := DecryptContent(key, tampered, "0xrec"); !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("tampered decryption error = %v, want ErrDecryptFailed", err)
	}

	if _, err := DecryptContent(key, sealed[:5], "0xrec"); !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("truncated decryption error = %v, want ErrDecryptFailed", err)
	}
}

// TestLockboxSubstitutionIsCaughtByAEAD is the end-to-end statement of the
// scheme's practical integrity model: the PRE layer cannot detect a
// substituted lockbox, because a substituted lockbox decrypts to a different
// but valid group element. Only the AEAD tag over the content catches it.
func TestLockboxSubstitutionIsCaughtByAEAD(t *testing.T) {
	patientPK, patientSK := newPair(t)

	salt, _ := NewSalt()
	nonce := bytes.Repeat([]byte{9}, 12)
	record := "0xabc"

	// Legitimate record.
	ct, err := SealRandom(patientPK)
	if err != nil {
		t.Fatalf("SealRandom: %v", err)
	}
	m, err := DelegateeSelfDecrypt(ct, patientSK)
	if err != nil {
		t.Fatalf("DelegateeSelfDecrypt: %v", err)
	}
	key, err := DeriveContentKeyFromMessage(m, salt, record)
	if err != nil {
		t.Fatalf("DeriveContentKeyFromMessage: %v", err)
	}
	sealed, _ := EncryptContent(key, []byte("diagnosis"), nonce, record)

	// An attacker with database write access swaps the lockbox for another
	// valid one under the same public key.
	evil, err := SealRandom(patientPK)
	if err != nil {
		t.Fatalf("SealRandom evil: %v", err)
	}
	evilM, err := DelegateeSelfDecrypt(evil, patientSK)
	if err != nil {
		t.Fatalf("DelegateeSelfDecrypt evil: %v", err)
	}
	evilKey, err := DeriveContentKeyFromMessage(evilM, salt, record)
	if err != nil {
		t.Fatalf("DeriveContentKeyFromMessage evil: %v", err)
	}

	// The PRE layer sees nothing wrong: both keys are well-formed.
	if bytes.Equal(key, evilKey) {
		t.Fatal("two distinct lockboxes produced the same content key")
	}
	// The AEAD tag is what fails.
	if _, err := DecryptContent(evilKey, sealed, record); !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("substituted lockbox error = %v, want ErrDecryptFailed", err)
	}
}

func TestNewSaltIsFresh(t *testing.T) {
	a, err := NewSalt()
	if err != nil {
		t.Fatalf("NewSalt: %v", err)
	}
	b, _ := NewSalt()
	if len(a) != SaltSize {
		t.Errorf("salt length = %d, want %d", len(a), SaltSize)
	}
	if bytes.Equal(a, b) {
		t.Error("two salts were identical")
	}
}
