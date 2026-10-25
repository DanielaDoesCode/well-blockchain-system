// Package hash implements the pure-Go Poseidon hash used by the log ledger.
//
// It is byte-for-byte compatible with:
//   - the research ZK circuit (zk-circuits/research_contribute.circom),
//   - the deployed PoseidonT3 library (src/lib/PoseidonT3.sol) and the
//     AttestationManager Merkle tree,
//   - circomlibjs (the JS library already used in zk-circuits/gen_merkle_proof.js).
//
// The implementation mirrors the standard circomlibjs "poseidon" construction
// for exactly two inputs, i.e. a HadesPermutation with t=3, 8 full rounds and
// 57 partial rounds over the BN254 scalar field. Constants are generated from
// circomlibjs's poseidon_constants module (see constants.gen.go).
package hash

import (
	"math/big"
)

// FieldPrime is the BN254 scalar field modulus.
var FieldPrime, _ = new(big.Int).SetString(
	"21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

var (
	arkBig = parseArk()
	mdsBig = parseMDS()
)

func parseBig(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("ledger/hash: bad constant: " + s)
	}
	return n
}

func parseArk() *[65][3]*big.Int {
	var out [65][3]*big.Int
	for r := 0; r < 65; r++ {
		for i := 0; i < 3; i++ {
			out[r][i] = parseBig(ark[r][i])
		}
	}
	return &out
}

func parseMDS() *[3][3]*big.Int {
	var out [3][3]*big.Int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			out[i][j] = parseBig(mds[i][j])
		}
	}
	return &out
}

func mod(n *big.Int) *big.Int { return new(big.Int).Mod(n, FieldPrime) }

func fieldAdd(a, b *big.Int) *big.Int { return mod(new(big.Int).Add(a, b)) }

func fieldMul(a, b *big.Int) *big.Int { return mod(new(big.Int).Mul(a, b)) }

// pow5 computes x^5 mod p.
func pow5(x *big.Int) *big.Int {
	x2 := fieldMul(x, x) // x^2
	x4 := fieldMul(x2, x2)
	return fieldMul(x4, x)
}

// Poseidon2 returns Poseidon(left, right) over the field, using the same
// two-input construction as the deployed PoseidonT3 and the research circuit.
// Both inputs must already be reduced mod FieldPrime.
func Poseidon2(left, right *big.Int) *big.Int {
	state := [3]*big.Int{new(big.Int), left, right}

	const nRoundsF = 8
	const nRoundsP = 57
	total := nRoundsF + nRoundsP

	for r := 0; r < total; r++ {
		for i := 0; i < 3; i++ {
			state[i] = fieldAdd(state[i], arkBig[r][i])
		}

		if r < nRoundsF/2 || r >= nRoundsF/2+nRoundsP {
			// full round: S-box on every element
			for i := 0; i < 3; i++ {
				state[i] = pow5(state[i])
			}
		} else {
			// partial round: S-box on the first element only
			state[0] = pow5(state[0])
		}

		s0 := fieldAdd(fieldMul(mdsBig[0][0], state[0]),
			fieldAdd(fieldMul(mdsBig[0][1], state[1]), fieldMul(mdsBig[0][2], state[2])))
		s1 := fieldAdd(fieldMul(mdsBig[1][0], state[0]),
			fieldAdd(fieldMul(mdsBig[1][1], state[1]), fieldMul(mdsBig[1][2], state[2])))
		s2 := fieldAdd(fieldMul(mdsBig[2][0], state[0]),
			fieldAdd(fieldMul(mdsBig[2][1], state[1]), fieldMul(mdsBig[2][2], state[2])))
		state = [3]*big.Int{s0, s1, s2}
	}

	return state[0]
}

// HashFold folds an arbitrary number of field elements into a single hash:
//
//	h = 0
//	for each v: h = Poseidon2(h, v)
//
// This is the domain-separating, variable-length hash used for row content.
func HashFold(vals ...*big.Int) *big.Int {
	h := new(big.Int)
	for _, v := range vals {
		h = Poseidon2(h, fieldMod(v))
	}
	return h
}

// fieldMod reduces v mod the field prime (identity when already reduced).
func fieldMod(v *big.Int) *big.Int {
	if v.Sign() < 0 {
		return mod(v)
	}
	if v.Cmp(FieldPrime) < 0 {
		return new(big.Int).Set(v)
	}
	return mod(v)
}

// FieldBytes big-endian-encodes a field element as exactly 32 bytes.
func FieldBytes(v *big.Int) []byte {
	return v.FillBytes(make([]byte, 32))
}

// BytesToField interprets up-to-32-byte big-endian bytes (e.g. an address or
// hash) as a field element, reducing mod p.
func BytesToField(b []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(b)
	return mod(z)
}

// BytesToField32 interprets a 32-byte value as a field element (reduced mod p).
func BytesToField32(b []byte) *big.Int {
	return BytesToField(b)
}