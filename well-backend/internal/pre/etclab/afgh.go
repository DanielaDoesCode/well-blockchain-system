// Package etclab is a vendored copy of github.com/etclab/pre (four commits,
// one star, unreleased), carrying local changes. Original BSD 3-Clause,
// Copyright (c) 2025 Stephen Herwig. See ./LICENSE and ./AUTHORS.
//
// Local changes, all deliberate and specified in docs/pre-design.md:
//
//   - Package renamed pre -> etclab to avoid colliding with our own package.
//   - The github.com/etclab/mu dependency is removed. It existed only to serve
//     mu.Panicf, and a panic on a marshal failure is unacceptable on a
//     request path. Callers must validate before calling in here.
//   - The library's KdfGtToAes256 is removed. It passes nil for both salt and
//     info to hkdf.New, giving no domain separation. Our package derives
//     content keys instead.
//
// Everything below the package clause is otherwise upstream code. The security
// properties of that code, including its lack of input validation, are not
// fixed here; see docs/pre-design.md section 4.
package etclab

import (
	bls "github.com/cloudflare/circl/ecc/bls12381"
)

type PublicParams struct {
	G1 *bls.G1
	G2 *bls.G2
	Z  *bls.Gt
}

func NewPublicParams() *PublicParams {
	pp := new(PublicParams)

	pp.G1 = bls.G1Generator()
	pp.G2 = bls.G2Generator()
	pp.Z = bls.Pair(pp.G1, pp.G2)

	return pp
}

type PublicKey struct {
	G1toA *bls.G1
	G2toA *bls.G2
}

type SecretKey struct {
	A *bls.Scalar
}

type ReEncryptionKey struct {
	RK *bls.G2
}

type KeyPair struct {
	PK *PublicKey
	SK *SecretKey
}

func KeyGen(pp *PublicParams) *KeyPair {
	sk := new(SecretKey)
	sk.A = randomScalar()

	pk := new(PublicKey)
	pk.G1toA = new(bls.G1)
	pk.G1toA.ScalarMult(sk.A, pp.G1)
	pk.G2toA = new(bls.G2)
	pk.G2toA.ScalarMult(sk.A, pp.G2)

	return &KeyPair{
		PK: pk,
		SK: sk,
	}
}

func ReEncryptionKeyGen(_ *PublicParams, aliceSK *SecretKey, bobPK *PublicKey) *ReEncryptionKey {
	aInv := new(bls.Scalar)
	aInv.Inv(aliceSK.A)

	rk := new(bls.G2)
	rk.ScalarMult(aInv, bobPK.G2toA)

	return &ReEncryptionKey{
		RK: rk,
	}
}

// Ciphertext1 is the paper's second-level ciphertext, (g^{a r}, m Z^r).
// Decryptable by the delegator (Decrypt1, the paper's D2) and re-encryptable by
// anyone holding the matching re-encryption key.
type Ciphertext1 struct {
	Alpha *bls.Gt
	Beta  *bls.G1
}

func Encrypt(pp *PublicParams, m *bls.Gt, pk *PublicKey) *Ciphertext1 {
	r := randomScalar()
	alpha := new(bls.Gt)
	alpha.Exp(pp.Z, r)
	alpha.Mul(alpha, m)

	beta := new(bls.G1)
	beta.ScalarMult(r, pk.G1toA)

	return &Ciphertext1{
		Alpha: alpha,
		Beta:  beta,
	}
}

// Ciphertext2 is a re-encrypted ciphertext, (Z^{b r}, m Z^r). Decryptable only
// by the delegatee (Decrypt2, the paper's D1).
type Ciphertext2 struct {
	Alpha *bls.Gt
	Beta  *bls.Gt
}

func ReEncrypt(pp *PublicParams, rk *ReEncryptionKey, ct1 *Ciphertext1) *Ciphertext2 {
	beta := bls.Pair(ct1.Beta, rk.RK)

	return &Ciphertext2{
		Alpha: ct1.Alpha,
		Beta:  beta,
	}
}

// Decrypt1 is the paper's D2, despite the library's naming. It inverts the
// pairing, so the first component of the ciphertext must be in G1.
func Decrypt1(pp *PublicParams, ct1 *Ciphertext1, sk *SecretKey) *bls.Gt {
	aInv := new(bls.Scalar)
	aInv.Inv(sk.A)
	tmp1 := new(bls.G2)
	tmp1.ScalarMult(aInv, pp.G2)

	tmp2 := bls.Pair(ct1.Beta, tmp1)
	tmp2.Inv(tmp2)

	m := new(bls.Gt)
	m.Mul(ct1.Alpha, tmp2)

	return m
}

// Decrypt2 is the paper's D1. Both components are already in Gt, so no pairing
// is needed.
func Decrypt2(pp *PublicParams, ct2 *Ciphertext2, sk *SecretKey) *bls.Gt {
	bInv := new(bls.Scalar)
	bInv.Inv(sk.A)
	tmp := new(bls.Gt)
	tmp.Exp(ct2.Beta, bInv)
	tmp.Inv(tmp)

	m := new(bls.Gt)
	m.Mul(ct2.Alpha, tmp)

	return m
}
