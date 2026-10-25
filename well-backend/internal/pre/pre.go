// Package pre implements the AFGH Second Attempt proxy re-encryption scheme
// (Green, Allen, Bgentz, Hall, Kilzer, "Proxy Re-Encryption Revisited", ePrint
// 2005/028, section 3.1) over BLS12-381, with the strict validation and
// domain separation that scheme requires in an adversarial setting.
//
// The design this implements is specified in docs/pre-design.md. That document
// is normative; this file is an implementation of it, and where the two
// disagree the document is right.
//
// # Naming
//
// The paper's D1 and D2 are numbered opposite to the vendored library's
// Decrypt1 and Decrypt2. This package uses the paper's names in identifiers:
//
//   - DelegateeSelfDecrypt is the paper's D2, using a pairing. It decrypts a
//     DelegatorCiphertext, the second-level form a delegator delegates from.
//   - DelegateeDecrypt is the paper's D1. It decrypts a ReencryptedCiphertext.
//
// # Trust boundary
//
// Every type in this package is a validated, canonical value. Nothing here
// accepts a raw byte slice without checking it, and nothing here panics on
// untrusted input. Decode functions are total: they either return a valid
// value or a typed error.
package pre

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"

	bls "github.com/cloudflare/circl/ecc/bls12381"

	"Thesis/well-blockchain/well-backend/internal/pre/etclab"
)

// Component sizes on BLS12-381, measured rather than assumed.
const (
	// SizeG1 is the length of a canonically compressed G1 point.
	SizeG1 = 48
	// SizeG2 is the length of a canonically compressed G2 point.
	SizeG2 = 96
	// SizeGt is the length of a Gt element in its canonical encoding.
	SizeGt = 576
)

// Errors returned by the decode and validation paths. Callers should test with
// errors.Is rather than comparing strings, so that wrapping stays safe.
var (
	ErrMalformed      = errors.New("pre: malformed encoding")
	ErrNotCanonical   = errors.New("pre: non-canonical encoding")
	ErrIdentity       = errors.New("pre: identity element rejected")
	ErrWrongType      = errors.New("pre: component in the wrong group")
	ErrDegenerate     = errors.New("pre: degenerate key material")
	ErrVersionUnknown = errors.New("pre: unknown envelope version")
)

// params holds the curve parameters. These are constants, so they are built
// once and never mutated.
var (
	params  *etclab.PublicParams
	genG1   *bls.G1
	genG2   *bls.G2
	zeroOne *bls.Scalar
)

func init() {
	params = etclab.NewPublicParams()
	genG1 = bls.G1Generator()
	genG2 = bls.G2Generator()
	zeroOne = new(bls.Scalar)
	zeroOne.SetOne()
}

// randSource is indirected so tests can force degenerate draws such as a zero
// scalar. Production code must not reassign it.
var randSource io.Reader = rand.Reader

// randomNonZeroScalar returns a uniformly random non-zero scalar. The library's
// randomScalar can return zero, and an inverse of zero is not meaningful, so
// key generation would silently produce an unusable key.
func randomNonZeroScalar() (*bls.Scalar, error) {
	for attempt := 0; attempt < 128; attempt++ {
		s := new(bls.Scalar)
		s.Random(randSource)
		if s.IsZero() == 0 {
			return s, nil
		}
	}
	return nil, fmt.Errorf("%w: could not draw a non-zero scalar in 128 attempts", ErrDegenerate)
}

// randomMessage returns a random non-identity element of Gt. The library's
// RandomGt returns the identity if either underlying draw is zero.
func randomMessage() (*bls.Gt, error) {
	for attempt := 0; attempt < 128; attempt++ {
		m := etclab.RandomGt()
		if !isGtIdentity(m) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%w: could not draw a non-identity Gt in 128 attempts", ErrDegenerate)
}

func isGtIdentity(g *bls.Gt) bool {
	var zero bls.Gt
	return g.IsEqual(&zero)
}

// decodeG1 decodes a canonically compressed G1 point.
//
// circl's SetBytes already enforces projectivity, curve membership and
// r-torsion membership, and rejects field elements at or above the field
// modulus. It also accepts the uncompressed form, which would let one point
// have two valid byte representations, so we re-encode and compare.
func decodeG1(b []byte) (*bls.G1, error) {
	if len(b) != SizeG1 {
		return nil, fmt.Errorf("%w: G1 must be %d bytes, got %d", ErrMalformed, SizeG1, len(b))
	}
	p := new(bls.G1)
	if err := p.SetBytes(b); err != nil {
		return nil, fmt.Errorf("%w: G1 rejected by decoder: %v", ErrMalformed, err)
	}
	if p.IsIdentity() {
		return nil, ErrIdentity
	}
	if subtle.ConstantTimeCompare(p.BytesCompressed(), b) != 1 {
		return nil, fmt.Errorf("%w: G1 is not in canonical compressed form", ErrNotCanonical)
	}
	return p, nil
}

// decodeG2 decodes a canonically compressed G2 point, with the same reasoning
// as decodeG1.
func decodeG2(b []byte) (*bls.G2, error) {
	if len(b) != SizeG2 {
		return nil, fmt.Errorf("%w: G2 must be %d bytes, got %d", ErrMalformed, SizeG2, len(b))
	}
	p := new(bls.G2)
	if err := p.SetBytes(b); err != nil {
		return nil, fmt.Errorf("%w: G2 rejected by decoder: %v", ErrMalformed, err)
	}
	if p.IsIdentity() {
		return nil, ErrIdentity
	}
	if subtle.ConstantTimeCompare(p.BytesCompressed(), b) != 1 {
		return nil, fmt.Errorf("%w: G2 is not in canonical compressed form", ErrNotCanonical)
	}
	return p, nil
}

// decodeGt decodes a Gt element.
//
// This is the one place where circl's validation is not sufficient. SetBytes
// equivalents for Gt accept an all-zero encoding with a nil error, and the
// result is the multiplicative identity, which re-serialises to all zeros. No
// range or format check catches it, so the identity is rejected explicitly.
func decodeGt(b []byte) (*bls.Gt, error) {
	if len(b) != SizeGt {
		return nil, fmt.Errorf("%w: Gt must be %d bytes, got %d", ErrMalformed, SizeGt, len(b))
	}
	g := new(bls.Gt)
	if err := g.UnmarshalBinary(b); err != nil {
		return nil, fmt.Errorf("%w: Gt rejected by decoder: %v", ErrMalformed, err)
	}
	if isGtIdentity(g) {
		return nil, ErrIdentity
	}
	out, err := g.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: Gt re-encode failed: %v", ErrMalformed, err)
	}
	if subtle.ConstantTimeCompare(out, b) != 1 {
		return nil, fmt.Errorf("%w: Gt is not in canonical form", ErrNotCanonical)
	}
	return g, nil
}

// SecretKey is a delegator or delegatee secret scalar, a in the paper.
type SecretKey struct {
	a *bls.Scalar
}

// Bytes returns the canonical 32-byte encoding of the secret scalar.
func (s *SecretKey) Bytes() []byte {
	b, _ := s.a.MarshalBinary()
	return b
}

// PublicKey is g^a, in both G1 and G2. The G2 component is what makes
// delegation possible: a re-encryption key is derived from it.
type PublicKey struct {
	g1 *bls.G1
	g2 *bls.G2
}

// MarshalBinary returns the canonical encoding: compressed G1 then compressed
// G2, 144 bytes total.
func (p *PublicKey) MarshalBinary() ([]byte, error) {
	if err := validatePublicKey(p); err != nil {
		return nil, err
	}
	out := make([]byte, 0, SizeG1+SizeG2)
	out = append(out, p.g1.BytesCompressed()...)
	out = append(out, p.g2.BytesCompressed()...)
	return out, nil
}

// ParsePublicKey decodes and validates a public key.
func ParsePublicKey(b []byte) (*PublicKey, error) {
	if len(b) != SizeG1+SizeG2 {
		return nil, fmt.Errorf("%w: public key must be %d bytes, got %d", ErrMalformed, SizeG1+SizeG2, len(b))
	}
	g1, err := decodeG1(b[:SizeG1])
	if err != nil {
		return nil, err
	}
	g2, err := decodeG2(b[SizeG1:])
	if err != nil {
		return nil, err
	}
	// The two components must be the same secret, otherwise this is not a
	// public key at all. Checking that requires a discrete log, which is
	// infeasible, so we cannot verify it here. The G1 component is used for
	// encryption and the G2 component for delegation, so a mismatched pair
	// yields a key that encrypts to nobody. This is a known limitation and is
	// why public keys must be registered through an authenticated channel.
	return &PublicKey{g1: g1, g2: g2}, nil
}

func validatePublicKey(p *PublicKey) error {
	if p == nil || p.g1 == nil || p.g2 == nil {
		return fmt.Errorf("%w: nil public key", ErrMalformed)
	}
	if p.g1.IsIdentity() || p.g2.IsIdentity() {
		return ErrIdentity
	}
	return nil
}

// KeyGen produces a fresh key pair. The secret scalar is drawn non-zero, so
// the resulting public key is never the identity and 1/a always exists.
func KeyGen() (*PublicKey, *SecretKey, error) {
	a, err := randomNonZeroScalar()
	if err != nil {
		return nil, nil, err
	}
	pk := &PublicKey{}
	pk.g1 = new(bls.G1)
	pk.g1.ScalarMult(a, genG1)
	pk.g2 = new(bls.G2)
	pk.g2.ScalarMult(a, genG2)
	return pk, &SecretKey{a: a}, nil
}

// ParseSecretKey decodes a secret scalar produced by SecretKey.Bytes. It lives
// on the client side; the backend never handles secret keys.
func ParseSecretKey(b []byte) (*SecretKey, error) {
	if len(b) != 32 {
		return nil, fmt.Errorf("%w: secret key must be 32 bytes, got %d", ErrMalformed, len(b))
	}
	a := new(bls.Scalar)
	if err := a.UnmarshalBinary(b); err != nil {
		return nil, fmt.Errorf("%w: decode secret key: %v", ErrMalformed, err)
	}
	if a.IsZero() == 1 {
		return nil, fmt.Errorf("%w: secret key is zero", ErrDegenerate)
	}
	return &SecretKey{a: a}, nil
}

// ReEncryptionKey is rk = (g^b)^{1/a}, the paper's rk_{A->B}. It lives in G2 so
// that it can be paired with a second-level ciphertext's G1 component.
//
// A re-encryption key is safe for the backend to hold. It converts a
// second-level ciphertext into one the delegatee can read, and cannot itself
// recover the message.
type ReEncryptionKey struct {
	rk *bls.G2
}

// MarshalBinary returns the 96-byte compressed encoding.
func (r *ReEncryptionKey) MarshalBinary() ([]byte, error) {
	if r == nil || r.rk == nil {
		return nil, fmt.Errorf("%w: nil re-encryption key", ErrMalformed)
	}
	if r.rk.IsIdentity() {
		return nil, ErrIdentity
	}
	return r.rk.BytesCompressed(), nil
}

// ParseReEncryptionKey decodes and validates a re-encryption key. This is on
// the backend's input path, since re-encryption keys arrive from patients.
func ParseReEncryptionKey(b []byte) (*ReEncryptionKey, error) {
	p, err := decodeG2(b)
	if err != nil {
		return nil, err
	}
	return &ReEncryptionKey{rk: p}, nil
}

// DeriveReEncryptionKey computes rk_{A->B} from the delegator's secret and the
// delegatee's public key. It runs on the patient side. The delegatee's secret
// is neither needed nor accepted.
func DeriveReEncryptionKey(delegatorSK *SecretKey, delegateePK *PublicKey) (*ReEncryptionKey, error) {
	if delegatorSK == nil || delegatorSK.a == nil {
		return nil, fmt.Errorf("%w: nil delegator secret", ErrMalformed)
	}
	if err := validatePublicKey(delegateePK); err != nil {
		return nil, err
	}
	if delegatorSK.a.IsZero() == 1 {
		return nil, fmt.Errorf("%w: delegator secret is zero", ErrDegenerate)
	}

	aInv := new(bls.Scalar)
	aInv.Inv(delegatorSK.a)
	if aInv.IsZero() == 1 {
		return nil, fmt.Errorf("%w: delegator secret is not invertible", ErrDegenerate)
	}

	rk := new(bls.G2)
	rk.ScalarMult(aInv, delegateePK.g2)
	if rk.IsIdentity() {
		return nil, fmt.Errorf("%w: derived re-encryption key is the identity", ErrDegenerate)
	}
	return &ReEncryptionKey{rk: rk}, nil
}

// DelegatorCiphertext is the paper's second-level ciphertext,
// (g^{a r}, m Z^r). It is what a delegator publishes and what they can
// themselves decrypt. A backend can convert it for a delegatee.
type DelegatorCiphertext struct {
	mGt *bls.Gt
	aGr *bls.G1
}

// MarshalBinary returns the canonical encoding: Gt then compressed G1, 624
// bytes.
func (c *DelegatorCiphertext) MarshalBinary() ([]byte, error) {
	if err := validateDelegatorCiphertext(c); err != nil {
		return nil, err
	}
	mBytes, err := c.mGt.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: encode Gt: %v", ErrMalformed, err)
	}
	out := make([]byte, 0, SizeGt+SizeG1)
	out = append(out, mBytes...)
	out = append(out, c.aGr.BytesCompressed()...)
	return out, nil
}

// ParseDelegatorCiphertext decodes and validates a second-level lockbox.
func ParseDelegatorCiphertext(b []byte) (*DelegatorCiphertext, error) {
	if len(b) != SizeGt+SizeG1 {
		return nil, fmt.Errorf("%w: delegator ciphertext must be %d bytes, got %d", ErrMalformed, SizeGt+SizeG1, len(b))
	}
	m, err := decodeGt(b[:SizeGt])
	if err != nil {
		return nil, err
	}
	gr, err := decodeG1(b[SizeGt:])
	if err != nil {
		return nil, err
	}
	return &DelegatorCiphertext{mGt: m, aGr: gr}, nil
}

func validateDelegatorCiphertext(c *DelegatorCiphertext) error {
	if c == nil || c.mGt == nil || c.aGr == nil {
		return fmt.Errorf("%w: nil delegator ciphertext", ErrMalformed)
	}
	if isGtIdentity(c.mGt) {
		return ErrIdentity
	}
	if c.aGr.IsIdentity() {
		return ErrIdentity
	}
	return nil
}

// ReencryptedCiphertext is the paper's re-encrypted ciphertext,
// (Z^{b r}, m Z^r). Only the delegatee can decrypt it.
type ReencryptedCiphertext struct {
	mZr *bls.Gt
	zBr *bls.Gt
}

// MarshalBinary returns the canonical encoding: two Gt elements, 1152 bytes.
func (c *ReencryptedCiphertext) MarshalBinary() ([]byte, error) {
	if err := validateReencryptedCiphertext(c); err != nil {
		return nil, err
	}
	first, err := c.mZr.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: encode Gt: %v", ErrMalformed, err)
	}
	second, err := c.zBr.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: encode Gt: %v", ErrMalformed, err)
	}
	out := make([]byte, 0, 2*SizeGt)
	out = append(out, first...)
	out = append(out, second...)
	return out, nil
}

// ParseReencryptedCiphertext decodes and validates a re-encrypted lockbox.
func ParseReencryptedCiphertext(b []byte) (*ReencryptedCiphertext, error) {
	if len(b) != 2*SizeGt {
		return nil, fmt.Errorf("%w: re-encrypted ciphertext must be %d bytes, got %d", ErrMalformed, 2*SizeGt, len(b))
	}
	first, err := decodeGt(b[:SizeGt])
	if err != nil {
		return nil, err
	}
	second, err := decodeGt(b[SizeGt:])
	if err != nil {
		return nil, err
	}
	return &ReencryptedCiphertext{mZr: first, zBr: second}, nil
}

func validateReencryptedCiphertext(c *ReencryptedCiphertext) error {
	if c == nil || c.mZr == nil || c.zBr == nil {
		return fmt.Errorf("%w: nil re-encrypted ciphertext", ErrMalformed)
	}
	if isGtIdentity(c.mZr) || isGtIdentity(c.zBr) {
		return ErrIdentity
	}
	return nil
}

// Seal encrypts a message under a delegator's public key, producing the
// second-level ciphertext. The message must be a non-identity element of Gt.
func Seal(m *bls.Gt, delegatorPK *PublicKey) (*DelegatorCiphertext, error) {
	if m == nil {
		return nil, fmt.Errorf("%w: nil message", ErrMalformed)
	}
	if isGtIdentity(m) {
		return nil, ErrIdentity
	}
	if err := validatePublicKey(delegatorPK); err != nil {
		return nil, err
	}
	ct := etclab.Encrypt(params, m, toEtclabPublicKey(delegatorPK))
	if ct == nil {
		return nil, fmt.Errorf("%w: encryption produced no ciphertext", ErrMalformed)
	}
	out := &DelegatorCiphertext{mGt: ct.Alpha, aGr: ct.Beta}
	if err := validateDelegatorCiphertext(out); err != nil {
		return nil, err
	}
	return out, nil
}

// SealRandom encrypts a freshly generated random non-identity message. This is
// the entry point for content keys: the caller never supplies the message.
func SealRandom(delegatorPK *PublicKey) (*DelegatorCiphertext, error) {
	m, err := randomMessage()
	if err != nil {
		return nil, err
	}
	return Seal(m, delegatorPK)
}

// Reencrypt converts a second-level ciphertext into one the delegatee can
// read. This runs on the backend, which learns nothing from it.
func Reencrypt(ct *DelegatorCiphertext, rk *ReEncryptionKey) (*ReencryptedCiphertext, error) {
	if err := validateDelegatorCiphertext(ct); err != nil {
		return nil, err
	}
	if rk == nil || rk.rk == nil {
		return nil, fmt.Errorf("%w: nil re-encryption key", ErrMalformed)
	}
	if rk.rk.IsIdentity() {
		return nil, ErrIdentity
	}
	out := etclab.ReEncrypt(params, &etclab.ReEncryptionKey{RK: rk.rk}, &etclab.Ciphertext1{
		Alpha: ct.mGt,
		Beta:  ct.aGr,
	})
	if out == nil {
		return nil, fmt.Errorf("%w: re-encryption produced no ciphertext", ErrMalformed)
	}
	res := &ReencryptedCiphertext{mZr: out.Alpha, zBr: out.Beta}
	if err := validateReencryptedCiphertext(res); err != nil {
		return nil, err
	}
	return res, nil
}

// DelegateeDecrypt is the paper's D1, despite the library's Decrypt2 name. It
// recovers the message from a re-encrypted ciphertext. Runs on the delegatee.
func DelegateeDecrypt(ct *ReencryptedCiphertext, delegateeSK *SecretKey) (*bls.Gt, error) {
	if err := validateReencryptedCiphertext(ct); err != nil {
		return nil, err
	}
	if delegateeSK == nil || delegateeSK.a == nil {
		return nil, fmt.Errorf("%w: nil delegatee secret", ErrMalformed)
	}
	if delegateeSK.a.IsZero() == 1 {
		return nil, fmt.Errorf("%w: delegatee secret is zero", ErrDegenerate)
	}
	m := etclab.Decrypt2(params, &etclab.Ciphertext2{Alpha: ct.mZr, Beta: ct.zBr}, &etclab.SecretKey{A: delegateeSK.a})
	if m == nil {
		return nil, fmt.Errorf("%w: decryption produced no message", ErrMalformed)
	}
	// A wrong key yields a wrong message rather than an error, because the
	// scheme has no CCA security. Callers must treat the result as unverified
	// until the AEAD tag over the content checks out. We reject the identity
	// only, since that is provably not a valid output.
	if isGtIdentity(m) {
		return nil, ErrIdentity
	}
	return m, nil
}

// DelegateeSelfDecrypt is the paper's D2, despite the library's Decrypt1 name.
// It recovers the message from the second-level ciphertext the delegator
// published, so a patient can always read their own record without going
// through a delegatee. Runs on the delegator.
func DelegateeSelfDecrypt(ct *DelegatorCiphertext, delegatorSK *SecretKey) (*bls.Gt, error) {
	if err := validateDelegatorCiphertext(ct); err != nil {
		return nil, err
	}
	if delegatorSK == nil || delegatorSK.a == nil {
		return nil, fmt.Errorf("%w: nil delegator secret", ErrMalformed)
	}
	if delegatorSK.a.IsZero() == 1 {
		return nil, fmt.Errorf("%w: delegator secret is zero", ErrDegenerate)
	}
	m := etclab.Decrypt1(params, &etclab.Ciphertext1{Alpha: ct.mGt, Beta: ct.aGr}, &etclab.SecretKey{A: delegatorSK.a})
	if m == nil {
		return nil, fmt.Errorf("%w: decryption produced no message", ErrMalformed)
	}
	if isGtIdentity(m) {
		return nil, ErrIdentity
	}
	return m, nil
}

func toEtclabPublicKey(p *PublicKey) *etclab.PublicKey {
	return &etclab.PublicKey{G1toA: p.g1, G2toA: p.g2}
}
