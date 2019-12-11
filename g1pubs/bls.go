package g1pubs

import (
	"bytes"
	"io"
	"log"
	"sort"

	"github.com/phoreproject/bls"
)

// Signature is a message signature.
type Signature struct {
	s *bls.G2Projective
}

// Serialize serializes a signature in compressed form.
func (s *Signature) Serialize() [96]byte {
	return bls.CompressG2(s.s.ToAffine())
}

func (s *Signature) String() string {
	return s.s.String()
}

// GetPoint gets the G2 point associated with the signature.
func (s *Signature) GetPoint() *bls.G2Projective {
	return s.s.Copy()
}

// NewSignatureFromG2 creates a new signature from a G2
// element.
func NewSignatureFromG2(g2 *bls.G2Affine) *Signature {
	return &Signature{g2.ToProjective()}
}

// DeserializeSignature deserializes a signature from bytes.
func DeserializeSignature(b [96]byte) (*Signature, error) {
	a, err := bls.DecompressG2(b)
	if err != nil {
		return nil, err
	}

	return &Signature{s: a.ToProjective()}, nil
}

// Copy returns a copy of the signature.
func (s *Signature) Copy() *Signature {
	return &Signature{s.s.Copy()}
}

// PublicKey is a public key.
type PublicKey struct {
	p *bls.G1Projective
}

func (p PublicKey) String() string {
	return p.p.String()
}

// GetPoint gets the G1 point associated with the public key.
func (p *PublicKey) GetPoint() *bls.G1Projective {
	return p.p.Copy()
}

// Serialize serializes a public key to bytes.
func (p PublicKey) Serialize() [48]byte {
	return bls.CompressG1(p.p.ToAffine())
}

// NewPublicKeyFromG1 creates a new public key from a G1 element.
func NewPublicKeyFromG1(g1 *bls.G1Affine) *PublicKey {
	return &PublicKey{g1.ToProjective()}
}

func concatAppend(slices [][]byte) []byte {
	var tmp []byte
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

// Equals checks if two public keys are equal
func (p PublicKey) Equals(other PublicKey) bool {
	return p.p.Equal(other.p)
}

// DeserializePublicKey deserializes a public key from
// bytes.
func DeserializePublicKey(b [48]byte) (*PublicKey, error) {
	a, err := bls.DecompressG1(b)
	if err != nil {
		return nil, err
	}

	return &PublicKey{p: a.ToProjective()}, nil
}

// SecretKey represents a BLS private key.
type SecretKey struct {
	f *bls.FR
}

// GetFRElement gets the underlying FR element.
func (s SecretKey) GetFRElement() *bls.FR {
	return s.f
}

func (s SecretKey) String() string {
	return s.f.String()
}

// Serialize serializes a secret key to bytes.
func (s SecretKey) Serialize() [32]byte {
	return s.f.Bytes()
}

// DeserializeSecretKey deserializes a secret key from
// bytes.
func DeserializeSecretKey(b [32]byte) *SecretKey {
	return &SecretKey{bls.FRReprToFR(bls.FRReprFromBytes(b))}
}

// DeriveSecretKey derives a secret key from
// bytes.
func DeriveSecretKey(b [32]byte) *SecretKey {
	return &SecretKey{bls.HashSecretKey(b)}
}

// Sign signs a message with a secret key.
func Sign(message []byte, key *SecretKey) *Signature {
	h := bls.HashG2(message).MulFR(key.f.ToRepr())
	return &Signature{s: h}
}

// SignWithDomain signs a message with a secret key and its domain.
func SignWithDomain(message [32]byte, key *SecretKey, domain [8]byte) *Signature {
	h := bls.HashG2WithDomain(message, domain).MulFR(key.f.ToRepr())
	return &Signature{s: h}
}

// PrivToPub converts the private key into a public key.
func PrivToPub(k *SecretKey) *PublicKey {
	return &PublicKey{p: bls.G1AffineOne.MulFR(k.f.ToRepr())}
}

// RandKey generates a random secret key.
func RandKey(r io.Reader) (*SecretKey, error) {
	k, err := bls.RandFR(r)
	if err != nil {
		return nil, err
	}
	s := &SecretKey{f: k}
	return s, nil
}

// KeyFromFQRepr returns a new key based on a FQRepr in
// FR.
func KeyFromFQRepr(i *bls.FRRepr) *SecretKey {
	return &SecretKey{f: bls.FRReprToFR(i)}
}

// Verify verifies a signature against a message and a public key.
func Verify(m []byte, pub *PublicKey, sig *Signature) bool {
	h := bls.HashG2(m)
	return bls.CompareTwoPairings(bls.G1ProjectiveOne, sig.s, pub.p, h.ToProjective())
}

// VerifyWithDomain verifies a signature against a message and a public key and a domain
func VerifyWithDomain(m [32]byte, pub *PublicKey, sig *Signature, domain [8]byte) bool {
	h := bls.HashG2WithDomain(m, domain)
	return bls.CompareTwoPairings(bls.G1ProjectiveOne, sig.s, pub.p, h.ToAffine().ToProjective())
}

// AggregateSignatures adds up all of the signatures.
func AggregateSignatures(s []*Signature) *Signature {
	newSig := &Signature{s: bls.G2ProjectiveZero.Copy()}
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

// AggregatePublicKeys adds public keys together.
func AggregatePublicKeys(p []*PublicKey) *PublicKey {
	newPub := &PublicKey{p: bls.G1ProjectiveZero.Copy()}
	for _, pub := range p {
		newPub.Aggregate(pub)
	}
	return newPub
}

// Aggregate adds two public keys together.
func (p *PublicKey) Aggregate(other *PublicKey) {
	newP := p.p.Add(other.p)
	p.p = newP
}

// Copy copies the public key and returns it.
func (p *PublicKey) Copy() *PublicKey {
	return &PublicKey{p: p.p.Copy()}
}

// NewAggregateSignature creates a blank aggregate signature.
func NewAggregateSignature() *Signature {
	return &Signature{s: bls.G2ProjectiveZero.Copy()}
}

// NewAggregatePubkey creates a blank public key.
func NewAggregatePubkey() *PublicKey {
	return &PublicKey{p: bls.G1ProjectiveZero.Copy()}
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

	msgsCopy := make([][]byte, len(msgs))

	for i, m := range msgs {
		msgsCopy[i] = make([]byte, len(m))
		copy(msgsCopy[i], m)
	}

	msgsSorted := sortByteArrays(msgsCopy)
	lastMsg := []byte(nil)

	// check for duplicates
	for _, m := range msgsSorted {
		if bytes.Equal(m, lastMsg) {
			return false
		}
		lastMsg = m
	}

	lhs := bls.Pairing(bls.G1ProjectiveOne, s.s)
	rhs := bls.FQ12One.Copy()
	for i := range pubKeys {
		h := bls.HashG2(msgs[i])
		rhs.MulAssign(bls.Pairing(pubKeys[i].p, h.ToProjective()))
	}
	return lhs.Equals(rhs)
}

// VerifyAggregateCommon verifies each public key against a message.
// This is vulnerable to rogue public-key attack. Each user must
// provide a proof-of-knowledge of the public key.
func (s *Signature) VerifyAggregateCommon(pubKeys []*PublicKey, msg []byte) bool {
	aggPub := AggregatePublicKeys(pubKeys)
	return Verify(msg, aggPub, s)
}

// VerifyAggregateCommonWithDomain verifies each public key against a message and
// its domain.
func (s *Signature) VerifyAggregateCommonWithDomain(pubKeys []*PublicKey, msg [32]byte, domain [8]byte) bool {
	aggPub := AggregatePublicKeys(pubKeys)
	return VerifyWithDomain(msg, aggPub, s, domain)
}

// VerifyAggregateWithDomain verifies each public key against each message and its domain.
func (s *Signature) VerifyAggregateWithDomain(pubKeys []*PublicKey, msgs [][32]byte, domain [8]byte) bool {
	if len(pubKeys) != len(msgs) {
		return false
	}
	lhs := bls.Pairing(bls.G1ProjectiveOne, s.s)
	rhs := bls.FQ12One.Copy()
	for i := range pubKeys {
		h := bls.HashG2WithDomain(msgs[i], domain)
		rhs.MulAssign(bls.Pairing(pubKeys[i].p, h.ToAffine().ToProjective()))
	}
	return lhs.Equals(rhs)
}
