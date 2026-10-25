package pre

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"testing"

	bls "github.com/cloudflare/circl/ecc/bls12381"
)

func newPair(t *testing.T) (*PublicKey, *SecretKey) {
	t.Helper()
	pk, sk, err := KeyGen()
	if err != nil {
		t.Fatalf("KeyGen: %v", err)
	}
	return pk, sk
}

func TestFullDelegationRoundTrip(t *testing.T) {
	patientPK, patientSK := newPair(t)
	doctorPK, doctorSK := newPair(t)

	ct, err := SealRandom(patientPK)
	if err != nil {
		t.Fatalf("SealRandom: %v", err)
	}

	// Patient self-recovery: the paper's D2.
	selfMsg, err := DelegateeSelfDecrypt(ct, patientSK)
	if err != nil {
		t.Fatalf("DelegateeSelfDecrypt: %v", err)
	}

	rk, err := DeriveReEncryptionKey(patientSK, doctorPK)
	if err != nil {
		t.Fatalf("DeriveReEncryptionKey: %v", err)
	}
	rct, err := Reencrypt(ct, rk)
	if err != nil {
		t.Fatalf("Reencrypt: %v", err)
	}

	// Doctor: the paper's D1.
	doctorMsg, err := DelegateeDecrypt(rct, doctorSK)
	if err != nil {
		t.Fatalf("DelegateeDecrypt: %v", err)
	}

	// Both paths must agree. This pins the paper/library naming inversion.
	if !selfMsg.IsEqual(doctorMsg) {
		t.Errorf("D2 (self) and D1 (delegatee) disagree")
	}
	if isGtIdentity(selfMsg) {
		t.Error("recovered message is the identity")
	}
}

func TestUnrelatedDoctorCannotDecrypt(t *testing.T) {
	patientPK, patientSK := newPair(t)
	doctorPK, _ := newPair(t)
	strangerPK, strangerSK := newPair(t)

	ct, _ := SealRandom(patientPK)
	rk, _ := DeriveReEncryptionKey(patientSK, doctorPK)
	rct, err := Reencrypt(ct, rk)
	if err != nil {
		t.Fatalf("Reencrypt: %v", err)
	}

	wrong, err := DelegateeDecrypt(rct, strangerSK)
	if err != nil {
		// Rejecting outright is also an acceptable outcome.
		return
	}
	trueMsg, err := DelegateeSelfDecrypt(ct, patientSK)
	if err != nil {
		t.Fatalf("DelegateeSelfDecrypt: %v", err)
	}
	if wrong.IsEqual(trueMsg) {
		t.Fatal("an unrelated secret key recovered the message")
	}
	_ = strangerPK
}

func TestWrongSecretGivesWrongMessage(t *testing.T) {
	// The scheme has no CCA security: a wrong key yields a wrong plaintext
	// rather than an error. Document that behaviour so callers do not rely on
	// decryption returning an error.
	patientPK, patientSK := newPair(t)
	_, strangerSK := newPair(t)

	ct, _ := SealRandom(patientPK)
	trueMsg, _ := DelegateeSelfDecrypt(ct, patientSK)
	wrongMsg, err := DelegateeSelfDecrypt(ct, strangerSK)
	if err != nil {
		return
	}
	if wrongMsg.IsEqual(trueMsg) {
		t.Fatal("wrong secret recovered the message")
	}
}

// --- serialisation round trips ---

func TestPublicKeyRoundTrip(t *testing.T) {
	pk, _ := newPair(t)
	b, err := pk.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if len(b) != SizeG1+SizeG2 {
		t.Fatalf("public key length = %d, want %d", len(b), SizeG1+SizeG2)
	}
	parsed, err := ParsePublicKey(b)
	if err != nil {
		t.Fatalf("ParsePublicKey: %v", err)
	}
	if !parsed.g1.IsEqual(pk.g1) || !parsed.g2.IsEqual(pk.g2) {
		t.Error("public key did not round trip")
	}
}

func TestReEncryptionKeyRoundTrip(t *testing.T) {
	_, patientSK := newPair(t)
	doctorPK, _ := newPair(t)

	rk, _ := DeriveReEncryptionKey(patientSK, doctorPK)
	b, err := rk.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if len(b) != SizeG2 {
		t.Fatalf("re-encryption key length = %d, want %d", len(b), SizeG2)
	}
	parsed, err := ParseReEncryptionKey(b)
	if err != nil {
		t.Fatalf("ParseReEncryptionKey: %v", err)
	}
	if !parsed.rk.IsEqual(rk.rk) {
		t.Error("re-encryption key did not round trip")
	}
}

func TestCiphertextRoundTrip(t *testing.T) {
	patientPK, patientSK := newPair(t)
	doctorPK, doctorSK := newPair(t)

	ct, _ := SealRandom(patientPK)
	ctBytes, err := ct.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if len(ctBytes) != SizeGt+SizeG1 {
		t.Errorf("delegator ciphertext length = %d, want %d", len(ctBytes), SizeGt+SizeG1)
	}
	parsedCT, err := ParseDelegatorCiphertext(ctBytes)
	if err != nil {
		t.Fatalf("ParseDelegatorCiphertext: %v", err)
	}
	msgA, _ := DelegateeSelfDecrypt(parsedCT, patientSK)
	msgB, _ := DelegateeSelfDecrypt(ct, patientSK)
	if !msgA.IsEqual(msgB) {
		t.Error("delegator ciphertext did not round trip")
	}

	rk, _ := DeriveReEncryptionKey(patientSK, doctorPK)
	rct, _ := Reencrypt(ct, rk)
	rctBytes, err := rct.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary re-encrypted: %v", err)
	}
	if len(rctBytes) != 2*SizeGt {
		t.Errorf("re-encrypted length = %d, want %d", len(rctBytes), 2*SizeGt)
	}
	parsedRCT, err := ParseReencryptedCiphertext(rctBytes)
	if err != nil {
		t.Fatalf("ParseReencryptedCiphertext: %v", err)
	}
	msgC, _ := DelegateeDecrypt(parsedRCT, doctorSK)
	if !msgC.IsEqual(msgA) {
		t.Error("re-encrypted ciphertext did not round trip")
	}
}

// --- validation ---

func TestGtIdentityRejected(t *testing.T) {
	// The documented gap: circl accepts an all-zero Gt encoding with no error.
	// We must reject it.
	ident := make([]byte, SizeGt)
	gt := new(bls.Gt)
	if err := gt.UnmarshalBinary(ident); err != nil {
		t.Fatalf("expected circl to accept all-zero Gt, got %v", err)
	}
	if !isGtIdentity(gt) {
		t.Fatal("expected all-zero Gt to be the identity")
	}
	if _, err := decodeGt(ident); !errors.Is(err, ErrIdentity) {
		t.Errorf("decodeGt(all-zero) error = %v, want ErrIdentity", err)
	}
	// And it must be rejected in a real lockbox slot too.
	_, err := ParseDelegatorCiphertext(make([]byte, SizeGt+SizeG1))
	if err == nil {
		t.Error("ParseDelegatorCiphertext accepted an all-zero Gt slot")
	}
}

func TestIdentityPointsRejected(t *testing.T) {
	// The canonical identity encoding has the infinity bit set and all other
	// bits clear. It is 48 bytes compressed, 96 uncompressed.
	g1 := new(bls.G1)
	zeroG1 := g1.BytesCompressed()
	if _, err := decodeG1(zeroG1); !errors.Is(err, ErrIdentity) {
		t.Errorf("decodeG1(identity) error = %v, want ErrIdentity", err)
	}
	// Uncompressed identity is rejected on length first.
	uncG1 := g1.Bytes()
	if _, err := decodeG1(uncG1); !errors.Is(err, ErrMalformed) {
		t.Errorf("decodeG1(uncompressed identity) error = %v, want ErrMalformed", err)
	}

	zeroG2 := new(bls.G2).BytesCompressed()
	if _, err := decodeG2(zeroG2); !errors.Is(err, ErrIdentity) {
		t.Errorf("decodeG2(identity) error = %v, want ErrIdentity", err)
	}
}

func TestUncompressedEncodingRejected(t *testing.T) {
	// circl accepts both compressed and uncompressed for the same point. We
	// must reject the uncompressed form so a point has one representation.
	g1 := bls.G1Generator()
	unc := g1.Bytes()
	if err := new(bls.G1).SetBytes(unc); err != nil {
		t.Fatalf("expected circl to accept uncompressed G1, got %v", err)
	}
	if _, err := decodeG1(unc); !errors.Is(err, ErrMalformed) {
		t.Errorf("decodeG1(uncompressed) error = %v, want ErrMalformed (wrong length)", err)
	}
	// Sanity: the compressed form is accepted.
	if _, err := decodeG1(g1.BytesCompressed()); err != nil {
		t.Errorf("decodeG1(compressed) unexpected error: %v", err)
	}
}

func TestGtNonCanonicalInputRejected(t *testing.T) {
	gt := bls.Pair(bls.G1Generator(), bls.G2Generator())
	canonical, err := gt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if _, err := decodeGt(canonical); err != nil {
		t.Errorf("decodeGt(canonical) unexpected error: %v", err)
	}

	// A field element at or above the modulus is non-canonical. circl rejects
	// this at the field level with "value out of range".
	outOfRange := make([]byte, len(canonical))
	for i := range outOfRange {
		outOfRange[i] = 0xff
	}
	if _, err := decodeGt(outOfRange); err == nil {
		t.Error("decodeGt accepted an out-of-range field element")
	}
}

func TestGtBitFlipYieldsDifferentButValidElement(t *testing.T) {
	// Worth pinning: a bit flip inside a valid Fp element produces a
	// *different but perfectly canonical* Gt element. It is not a parsing
	// failure and our decoder must not claim otherwise. Correctness here rests
	// on the AEAD tag over the content, not on the lockbox encoding. This is
	// the practical face of the scheme's lack of CCA security.
	gt := bls.Pair(bls.G1Generator(), bls.G2Generator())
	b, err := gt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	mangled := make([]byte, len(b))
	copy(mangled, b)
	mangled[300] ^= 0x01

	parsed, err := decodeGt(mangled)
	if err != nil {
		t.Skipf("bit flip happened to land outside the field range: %v", err)
	}
	if parsed.IsEqual(gt) {
		t.Error("bit flip did not change the element; test is not exercising anything")
	}
}

func TestWrongLengthsRejected(t *testing.T) {
	cases := []struct {
		name string
		fn   func([]byte) error
		size int
	}{
		{"g1 short", func(b []byte) error { _, err := decodeG1(b); return err }, SizeG1},
		{"g1 long", func(b []byte) error { _, err := decodeG1(b); return err }, SizeG1},
		{"g2 short", func(b []byte) error { _, err := decodeG2(b); return err }, SizeG2},
		{"gt short", func(b []byte) error { _, err := decodeGt(b); return err }, SizeGt},
		{"gt long", func(b []byte) error { _, err := decodeGt(b); return err }, SizeGt},
	}
	for _, tc := range cases {
		for _, delta := range []int{-1, +1} {
			b := make([]byte, tc.size+delta)
			// Fill with a valid generator encoding where possible.
			if tc.size == SizeG1 {
				copy(b, bls.G1Generator().BytesCompressed())
			} else if tc.size == SizeG2 {
				copy(b, bls.G2Generator().BytesCompressed())
			} else {
				gt := bls.Pair(bls.G1Generator(), bls.G2Generator())
				v, _ := gt.MarshalBinary()
				copy(b, v)
			}
			if err := tc.fn(b); err == nil {
				t.Errorf("%s (%d bytes): expected rejection", tc.name, len(b))
			}
		}
	}
}

func TestCiphertextTypeConfusionRejected(t *testing.T) {
	patientPK, _ := newPair(t)
	ct, _ := SealRandom(patientPK)
	ctBytes, _ := ct.MarshalBinary()

	// Feed a re-encrypted ciphertext to the delegator parser and vice versa.
	rk, err := DeriveReEncryptionKey(mustPair(t).sk, patientPK)
	if err != nil {
		t.Fatalf("DeriveReEncryptionKey: %v", err)
	}
	rct, _ := Reencrypt(ct, rk)
	rctBytes, _ := rct.MarshalBinary()

	// 1152 bytes into a 624-byte parser.
	if _, err := ParseDelegatorCiphertext(rctBytes); err == nil {
		t.Error("ParseDelegatorCiphertext accepted a re-encrypted ciphertext")
	}
	// And the reverse.
	if _, err := ParseReencryptedCiphertext(ctBytes); err == nil {
		t.Error("ParseReencryptedCiphertext accepted a delegator ciphertext")
	}
}

type pair struct {
	pk *PublicKey
	sk *SecretKey
}

func mustPair(t *testing.T) pair {
	t.Helper()
	pk, sk := newPair(t)
	return pair{pk, sk}
}

func TestNilInputsRejected(t *testing.T) {
	ct, _ := SealRandom(mustPair(t).pk)
	rk, _ := DeriveReEncryptionKey(mustPair(t).sk, mustPair(t).pk)

	if _, err := SealRandom(nil); err == nil {
		t.Error("SealRandom(nil pk) accepted")
	}
	if _, err := Reencrypt(nil, rk); err == nil {
		t.Error("Reencrypt(nil ct) accepted")
	}
	if _, err := Reencrypt(ct, nil); err == nil {
		t.Error("Reencrypt(nil rk) accepted")
	}
	if _, err := DelegateeDecrypt(nil, mustPair(t).sk); err == nil {
		t.Error("DelegateeDecrypt(nil ct) accepted")
	}
	if _, err := DelegateeSelfDecrypt(nil, mustPair(t).sk); err == nil {
		t.Error("DelegateeSelfDecrypt(nil ct) accepted")
	}
	if _, err := DeriveReEncryptionKey(nil, mustPair(t).pk); err == nil {
		t.Error("DeriveReEncryptionKey(nil sk) accepted")
	}
}

// --- degenerate randomness ---

func TestZeroScalarIsResampled(t *testing.T) {
	old := randSource
	defer func() { randSource = old }()

	// First draw zero, then behave normally.
	zeroed := &oneShotZeroReader{after: rand.Reader}
	randSource = zeroed

	pk, sk, err := KeyGen()
	if err != nil {
		t.Fatalf("KeyGen with a zero first draw: %v", err)
	}
	if sk.a.IsZero() == 1 {
		t.Error("KeyGen produced a zero secret")
	}
	if pk.g1.IsIdentity() || pk.g2.IsIdentity() {
		t.Error("KeyGen produced an identity public key")
	}
	if zeroed.calls < 2 {
		t.Errorf("expected a resample after the zero draw, calls = %d", zeroed.calls)
	}
}

type oneShotZeroReader struct {
	after io.Reader
	calls int
}

func (o *oneShotZeroReader) Read(p []byte) (int, error) {
	o.calls++
	if o.calls == 1 {
		for i := range p {
			p[i] = 0
		}
		return len(p), nil
	}
	return o.after.Read(p)
}

func TestRandomMessageIsNotIdentity(t *testing.T) {
	for i := 0; i < 200; i++ {
		ct, err := SealRandom(mustPair(t).pk)
		if err != nil {
			t.Fatalf("SealRandom: %v", err)
		}
		if isGtIdentity(ct.mGt) {
			t.Fatal("SealRandom produced an identity message")
		}
	}
}

func TestParseSecretKeyRoundTrip(t *testing.T) {
	pk, sk := newPair(t)
	parsed, err := ParseSecretKey(sk.Bytes())
	if err != nil {
		t.Fatalf("ParseSecretKey: %v", err)
	}
	if !bytes.Equal(parsed.Bytes(), sk.Bytes()) {
		t.Fatal("secret key did not round-trip")
	}

	// A secret restored from bytes must reproduce the public key and must be
	// able to open a lockbox sealed to pk.
	rePK := &PublicKey{g1: new(bls.G1), g2: new(bls.G2)}
	rePK.g1.ScalarMult(parsed.a, genG1)
	rePK.g2.ScalarMult(parsed.a, genG2)
	if !rePK.g1.IsEqual(pk.g1) || !rePK.g2.IsEqual(pk.g2) {
		t.Fatal("restored secret does not reproduce the public key")
	}

	ct, err := SealRandom(pk)
	if err != nil {
		t.Fatalf("SealRandom: %v", err)
	}
	if _, err := DelegateeSelfDecrypt(ct, parsed); err != nil {
		t.Fatalf("self decrypt with restored key: %v", err)
	}

	if _, err := ParseSecretKey(make([]byte, 32)); !errors.Is(err, ErrDegenerate) {
		t.Fatalf("zero secret: want ErrDegenerate, got %v", err)
	}
	if _, err := ParseSecretKey([]byte{1, 2, 3}); !errors.Is(err, ErrMalformed) {
		t.Fatalf("short secret: want ErrMalformed, got %v", err)
	}
}
