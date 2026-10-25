package pre

import (
	"testing"
)

// The decode entry points are the trust boundary. They must never panic and
// never allocate without bound, regardless of input. This is the main defence
// against a remote crash, given that the vendored library panics on some
// malformed input and we must not let unvalidated bytes reach it.

func FuzzParsePublicKey(f *testing.F) {
	pk, _, _ := KeyGen()
	seed, _ := pk.MarshalBinary()
	f.Add(seed)
	f.Add([]byte{})
	f.Add(make([]byte, SizeG1+SizeG2))

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on input of length %d: %v", len(b), r)
			}
		}()
		p, err := ParsePublicKey(b)
		if err != nil {
			return
		}
		// Anything accepted must re-encode to exactly the input.
		out, err := p.MarshalBinary()
		if err != nil {
			t.Fatalf("accepted input failed to re-encode: %v", err)
		}
		if len(out) != len(b) {
			t.Fatalf("accepted input re-encoded to length %d, want %d", len(out), len(b))
		}
		for i := range out {
			if out[i] != b[i] {
				t.Fatalf("accepted non-canonical input at byte %d", i)
			}
		}
	})
}

func FuzzParseDelegatorCiphertext(f *testing.F) {
	pk, _, _ := KeyGen()
	ct, _ := SealRandom(pk)
	seed, _ := ct.MarshalBinary()
	f.Add(seed)
	f.Add([]byte{})
	f.Add(make([]byte, SizeGt+SizeG1))

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on input of length %d: %v", len(b), r)
			}
		}()
		c, err := ParseDelegatorCiphertext(b)
		if err != nil {
			return
		}
		out, err := c.MarshalBinary()
		if err != nil {
			t.Fatalf("accepted input failed to re-encode: %v", err)
		}
		if len(out) != len(b) {
			t.Fatalf("non-canonical: re-encoded to %d bytes, want %d", len(out), len(b))
		}
		for i := range out {
			if out[i] != b[i] {
				t.Fatalf("accepted non-canonical input at byte %d", i)
			}
		}
	})
}

func FuzzParseReencryptedCiphertext(f *testing.F) {
	pk, _, _ := KeyGen()
	ct, _ := SealRandom(pk)
	_, delegatorSK, _ := KeyGen()
	rk, _ := DeriveReEncryptionKey(delegatorSK, pk)
	rct, _ := Reencrypt(ct, rk)
	seed, _ := rct.MarshalBinary()
	f.Add(seed)
	f.Add([]byte{})
	f.Add(make([]byte, 2*SizeGt))

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on input of length %d: %v", len(b), r)
			}
		}()
		c, err := ParseReencryptedCiphertext(b)
		if err != nil {
			return
		}
		out, err := c.MarshalBinary()
		if err != nil {
			t.Fatalf("accepted input failed to re-encode: %v", err)
		}
		if len(out) != len(b) {
			t.Fatalf("non-canonical: re-encoded to %d bytes, want %d", len(out), len(b))
		}
		for i := range out {
			if out[i] != b[i] {
				t.Fatalf("accepted non-canonical input at byte %d", i)
			}
		}
	})
}

func FuzzParseReEncryptionKey(f *testing.F) {
	_, sk, _ := KeyGen()
	pk2, _, _ := KeyGen()
	rk, _ := DeriveReEncryptionKey(sk, pk2)
	seed, _ := rk.MarshalBinary()
	f.Add(seed)
	f.Add([]byte{})
	f.Add(make([]byte, SizeG2))

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on input of length %d: %v", len(b), r)
			}
		}()
		rk, err := ParseReEncryptionKey(b)
		if err != nil {
			return
		}
		out, err := rk.MarshalBinary()
		if err != nil {
			t.Fatalf("accepted input failed to re-encode: %v", err)
		}
		if len(out) != len(b) {
			t.Fatalf("non-canonical: re-encoded to %d bytes, want %d", len(out), len(b))
		}
	})
}
