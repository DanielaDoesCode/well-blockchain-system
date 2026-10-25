package hash

import "math/big"

// MerkleRoot computes the Merkle root of a balanced binary Poseidon tree over
// the given leaves (field elements). The leaf set is zero-padded up to the next
// power of two with the field zero, mirroring AttestationManager's zero leaves.
//
// A single leaf returns that leaf as the root.
func MerkleRoot(leaves []*big.Int) *big.Int {
	n := nextPow2(len(leaves))
	level := make([]*big.Int, n)
	copy(level, leaves)
	for i := len(leaves); i < n; i++ {
		level[i] = new(big.Int).Set(Zero)
	}
	for n > 1 {
		next := make([]*big.Int, n/2)
		for i := 0; i < n/2; i++ {
			next[i] = Poseidon2(level[2*i], level[2*i+1])
		}
		level = next
		n = n / 2
	}
	return level[0]
}

// MerkleProof returns the sibling hashes on the path from leaf at position idx
// (0-based) up to the root. The tree is rebuilt from the full leaf set.
func MerkleProof(leaves []*big.Int, idx int) (siblings []*big.Int, ok bool) {
	if idx < 0 || idx >= len(leaves) {
		return nil, false
	}
	n := nextPow2(len(leaves))
	level := make([]*big.Int, n)
	copy(level, leaves)
	for i := len(leaves); i < n; i++ {
		level[i] = new(big.Int).Set(Zero)
	}

	var out []*big.Int
	i := idx
	for n > 1 {
		sib := i ^ 1
		out = append(out, new(big.Int).Set(level[sib]))
		next := make([]*big.Int, n/2)
		for k := 0; k < n/2; k++ {
			next[k] = Poseidon2(level[2*k], level[2*k+1])
		}
		level = next
		i = i / 2
		n = n / 2
	}
	return out, true
}

// VerifyMerkle recomputes the root from the leaf and its sibling path and
// compares it with the expected root.
func VerifyMerkle(leaf *big.Int, idx int, siblings []*big.Int, root *big.Int) bool {
	h := new(big.Int).Set(leaf)
	if idx < 0 {
		return false
	}
	i := idx
	for _, sib := range siblings {
		if i&1 == 0 {
			h = Poseidon2(h, sib)
		} else {
			h = Poseidon2(sib, h)
		}
		i = i / 2
	}
	return h.Cmp(root) == 0
}

func nextPow2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}