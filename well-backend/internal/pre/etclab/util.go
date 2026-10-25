package etclab

import (
	"crypto/rand"
	"io"

	bls "github.com/cloudflare/circl/ecc/bls12381"
)

// randomScalar draws a scalar from crypto/rand. It may return zero, which the
// caller must handle; see docs/pre-design.md section 4.2.
func randomScalar() *bls.Scalar {
	z := new(bls.Scalar)
	z.Random(randReader)
	return z
}

// RandomGt returns a random element of Gt, or the identity if either draw is
// zero. Callers must check IsIdentity on the result and redraw; see
// docs/pre-design.md section 9.
func RandomGt() *bls.Gt {
	a := randomScalar()
	b := randomScalar()

	g1 := bls.G1Generator()
	g2 := bls.G2Generator()

	g1.ScalarMult(a, g1)
	g2.ScalarMult(b, g2)

	z := bls.Pair(g1, g2)
	return z
}

// randReader is indirected so tests can force degenerate draws. Production
// code must not reassign it.
var randReader io.Reader = rand.Reader
