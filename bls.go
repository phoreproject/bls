package bls

import (
	"bytes"
	"io"
	"log"
	"math/big"
	"sort"
)

// Signature is a message signature.
type Signature struct {
	s *G1Projective
}

// Serialize serializes a signature in compressed form.
func (s *Signature) Serialize() []byte {
	return CompressG1(s.s.ToAffine()).Bytes()
}

// DeserializeSignature deserializes a signature from bytes.
func DeserializeSignature(b []byte) (*Signature, error) {
	a, err := DecompressG1(new(big.Int).SetBytes(b))
	if err != nil {
		return nil, err
	}

	return &Signature{s: a.ToProjective()}, nil
}

// PublicKey is a public key.
type PublicKey struct {
	p *G2Projective
}

func (g PublicKey) String() string {
	return g.p.String()
}

// Serialize serializes a public key to bytes.
func (g PublicKey) Serialize() []byte {
	return CompressG2(g.p.ToAffine()).Bytes()
}

// DeserializePublicKey deserializes a public key from
// bytes.
func DeserializePublicKey(b []byte) (*PublicKey, error) {
	a, err := DecompressG2(new(big.Int).SetBytes(b))
	if err != nil {
		return nil, err
	}

	return &PublicKey{p: a.ToProjective()}, nil
}

// SecretKey represents a BLS private key.
type SecretKey struct {
	f *FR
}

func (g SecretKey) String() string {
	return g.f.String()
}

// Serialize serializes a secret key to bytes.
func (s SecretKey) Serialize() []byte {
	return s.f.ToBig().Bytes()
}

// DeserializeSecretKey deserializes a secret key from
// bytes.
func DeserializeSecretKey(b []byte) *SecretKey {
	return &SecretKey{NewFR(new(big.Int).SetBytes(b))}
}

// Sign signs a message with a secret key.
func Sign(message []byte, key *SecretKey) *Signature {
	h := HashG1(message).Mul(key.f.n)
	return &Signature{s: h}
}

// PrivToPub converts the private key into a public key.
func PrivToPub(k *SecretKey) *PublicKey {
	return &PublicKey{p: G2AffineOne.Mul(k.f.n)}
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
	lhs := Pairing(sig.s, G2ProjectiveOne)
	rhs := Pairing(h, pub.p)
	return lhs.Equals(rhs)
}

// AggregateSignatures adds up all of the signatures.
func AggregateSignatures(s []*Signature) *Signature {
	newSig := &Signature{s: G1ProjectiveZero}
	for _, sig := range s {
		newSig.Aggregate(sig)
	}
	return newSig
}

// Aggregate adds one signature to another
func (s *Signature) Aggregate(other *Signature) {
	newS := s.s.Add(other.s)
	s.s = newS
}

// NewAggregateSignature creates a blank aggregate signature.
func NewAggregateSignature() *Signature {
	return &Signature{s: G1ProjectiveZero.Copy()}
}

// implement `Interface` in sort package.
type sortableByteArray [][]byte

func (b sortableByteArray) Len() int {
	return len(b)
}

func (b sortableByteArray) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortableByteArray) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

func sortByteArrays(src [][]byte) [][]byte {
	sorted := sortableByteArray(src)
	sort.Sort(sorted)
	return sorted
}

// VerifyAggregate verifies each public key against each message.
func (s *Signature) VerifyAggregate(pubKeys []*PublicKey, msgs [][]byte) bool {
	if len(pubKeys) != len(msgs) {
		return false
	}

	// messages must be distinct
	msgsSorted := sortByteArrays(msgs)
	lastMsg := []byte(nil)

	// check for duplicates
	for _, m := range msgsSorted {
		if bytes.Equal(m, lastMsg) {
			return false
		}
		lastMsg = m
	}

	lhs := Pairing(s.s, G2ProjectiveOne)
	rhs := FQ12One.Copy()
	for i := range pubKeys {
		h := HashG1(msgs[i])
		rhs.MulAssign(Pairing(h, pubKeys[i].p))
	}
	return lhs.Equals(rhs)
}

// VerifyAggregateCommon verifies each public key against a message.
// This is vulnerable to rogue public-key attack. Each user must
// provide a proof-of-knowledge of the public key.
func (s *Signature) VerifyAggregateCommon(pubKeys []*PublicKey, msg []byte) bool {
	h := HashG1(msg)
	lhs := Pairing(s.s, G2ProjectiveOne)
	rhs := FQ12One.Copy()
	for _, p := range pubKeys {
		rhs.MulAssign(Pairing(h, p.p))
	}
	return lhs.Equals(rhs)
}
