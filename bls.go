package bls

import (
	"io"
	"math/big"
)

// Signature is a message signature.
type Signature struct {
	s *G1Affine
}

// PublicKey is a public key.
type PublicKey struct {
	p *G2Affine
}

func (g PublicKey) String() string {
	return g.p.String()
}

// SecretKey represents a BLS private key.
type SecretKey struct {
	f *FR
}

func (g SecretKey) String() string {
	return g.f.String()
}

// Sign signs a message with a secret key.
func Sign(message []byte, key *SecretKey) *Signature {
	h := HashG1(message).Mul(key.f.n)
	return &Signature{s: h.ToAffine()}
}

// PrivToPub converts the private key into a public key.
func PrivToPub(k *SecretKey) *PublicKey {
	return &PublicKey{p: G2AffineOne.Copy().Mul(k.f.n).ToAffine()}
}

// RandKey generates a random secret key.
func RandKey(r io.Reader) (*SecretKey, error) {
	k, err := RandFR(r)
	if err != nil {
		return nil, err
	}
	s := &SecretKey{f: k}
	return s, nil
}

// KeyFromBig returns a new key based on a big int in
// FR.
func KeyFromBig(i *big.Int) *SecretKey {
	return &SecretKey{f: NewFR(i)}
}

// Verify verifies a signature against a message and a public key.
func Verify(m []byte, pub *PublicKey, sig *Signature) bool {
	h := HashG1(m)
	lhs := Pairing(sig.s.ToProjective(), G2ProjectiveOne.Copy())
	rhs := Pairing(h, pub.p.ToProjective())
	return lhs.Equals(rhs)
}
