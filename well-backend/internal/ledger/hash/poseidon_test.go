package hash

import (
	"math/big"
	"testing"
)

// Golden fixtures were generated with circomlibjs buildPoseidon (the same
// construction as the deployed PoseidonT3 and the research ZK circuit).
func TestPoseidon2Golden(t *testing.T) {
	cases := []struct{ a, b, want string }{
		{"1", "2", "7853200120776062878684798364095072458815029376092732009249414926327459813530"},
		{"0", "0", "14744269619966411208579211824598458697587494354926760081771325075741142829156"},
		{"5", "5", "14848575449521340934220251267929796113247500567202650801215929266644347883284"},
		{"12345", "67890", "11344094074881186137859743404234365978119253787583526441303892667757095072923"},
		{"7", "999999", "1514322093323659934860758432526924160104843628712743518852087058129719276855"},
		{"101", "202", "19311084779183973643093452456901609492733052887983755367564460719253928490098"},
	}
	for _, c := range cases {
		a, _ := new(big.Int).SetString(c.a, 10)
		b, _ := new(big.Int).SetString(c.b, 10)
		got := Poseidon2(a, b)
		if got.String() != c.want {
			t.Errorf("Poseidon2(%s,%s) = %s, want %s", c.a, c.b, got, c.want)
		}
	}
}

func TestHashFoldGolden(t *testing.T) {
	fold := func(vals ...string) *big.Int {
		in := make([]*big.Int, len(vals))
		for i, v := range vals {
			in[i], _ = new(big.Int).SetString(v, 10)
		}
		return HashFold(in...)
	}
	if got := fold(); got.Sign() != 0 {
		t.Errorf("HashFold() = %s, want 0", got)
	}
	if got := fold("42"); got.String() != "9904646155488355737762297645225334693069781832889131634543122060982625196787" {
		t.Errorf("HashFold([42]) = %s", got)
	}
	if got := fold("1", "2", "3"); got.String() != "20127075603631019434055928315203707068407414306847615530687456290565086592967" {
		t.Errorf("HashFold([1,2,3]) = %s", got)
	}
}

func TestMerkleRootGolden(t *testing.T) {
	// leaves [10,20,30] padded with zero → root golden from circomlibjs.
	leaves := []*big.Int{big.NewInt(10), big.NewInt(20), big.NewInt(30)}
	root := MerkleRoot(leaves)
	const want = "8773947246789296406810645929898101869102898271330577097575788594035597337658"
	if root.String() != want {
		t.Fatalf("MerkleRoot = %s, want %s", root, want)
	}

	sibs, ok := MerkleProof(leaves, 0)
	if !ok || len(sibs) != 2 {
		t.Fatalf("unexpected proof: len=%d ok=%v", len(sibs), ok)
	}
	if !VerifyMerkle(leaves[0], 0, sibs, root) {
		t.Error("VerifyMerkle(leaf 0) failed")
	}
	sibs2, _ := MerkleProof(leaves, 2)
	if !VerifyMerkle(leaves[2], 2, sibs2, root) {
		t.Error("VerifyMerkle(leaf 2) failed")
	}
	// Tampering with the leaf must break the proof.
	if VerifyMerkle(big.NewInt(31), 2, sibs2, root) {
		t.Error("VerifyMerkle accepted a tampered leaf")
	}
}

func TestRowHashChains(t *testing.T) {
	content := StrField("a-row")
	g1 := RowHash(nil, content) // genesis
	g2 := RowHash(hashBytes(g1), content)
	if g1.Cmp(g2) == 0 {
		t.Error("chained hashes must differ")
	}
	// A different predecessor must change the hash.
	g3 := RowHash(hashBytes(big.NewInt(1)), content)
	if g2.Cmp(g3) == 0 {
		t.Error("different predecessor must change the hash")
	}
}

func TestJSONCanonicalization(t *testing.T) {
	// The ledger data values are strings (hex hashes / decimal big-ints), so
	// number-text-preserving canonicalisation is deterministic across the
	// insert → PostgreSQL jsonb → read-back round trip.
	a := JSONBytesField([]byte(`{"requestId":"123","responder":"0xabc","b":1}`))
	b := JSONBytesField([]byte(`{ "b": 1,  "requestId": "123",  "responder":"0xabc"}`))
	if a.Cmp(b) != 0 {
		t.Error("canonical JSON hashing is not order/whitespace invariant")
	}
	// And DataField on the equivalent Go value agrees.
	m := map[string]any{"requestId": "123", "responder": "0xabc", "b": float64(1)}
	if DataField(m).Cmp(a) != 0 {
		t.Error("DataField(map) != JSONBytesField(raw)")
	}
	// Number text is preserved verbatim (1.50 stays 1.50), not float-normalised.
	c := JSONBytesField([]byte(`{"n":1.50}`))
	d := JSONBytesField([]byte(`{"n":1.5}`))
	if c.Cmp(d) == 0 {
		t.Error("number text must be preserved, not float-normalised")
	}
}

func hashBytes(v *big.Int) []byte { return FieldBytes(v) }