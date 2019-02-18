package bls

import (
	"bytes"
	"io"
	"log"
	"sort"
)

// Signature is a message signature.
type Signature struct {
	s *G1Projective
}

// Serialize serializes a signature in compressed form.
func (s *Signature) Serialize() [48]byte {
	return CompressG1(s.s.ToAffine())
}

// DeserializeSignature deserializes a signature from bytes.
func DeserializeSignature(b [48]byte) (*Signature, error) {
	a, err := DecompressG1(b)
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
	p *G2Projective
}

func (p PublicKey) String() string {
	return p.p.String()
}

// Serialize serializes a public key to bytes.
func (p PublicKey) Serialize() [96]byte {
	return CompressG2(p.p.ToAffine())
}

func concatAppend(slices [][]byte) []byte {
	var tmp []byte
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

// SerializeBig serializes a public key uncompressed.
func (p PublicKey) SerializeBig() [193]byte {
	affine := p.p.ToAffine()
	out := [193]byte{}
	infinity := affine.infinity
	copy(out[1:49], affine.x.c0.n.ToBig().Bytes())
	copy(out[49:97], affine.x.c1.n.ToBig().Bytes())
	copy(out[97:145], affine.y.c0.n.ToBig().Bytes())
	copy(out[145:193], affine.y.c1.n.ToBig().Bytes())

	if infinity {
		out[0] = 1
	}

	return out
}

// DeserializePublicKeyBig deserializes a public key uncompressed.
func DeserializePublicKeyBig(bigPublicKey [193]byte) *PublicKey {
	g := G2Affine{}
	if bigPublicKey[0] == 1 {
		g.infinity = true
		return &PublicKey{p: g.ToProjective()}
	}
	g.x = &FQ2{}
	g.y = &FQ2{}
	var xC0Bytes [48]byte
	var xC1Bytes [48]byte
	var yC0Bytes [48]byte
	var yC1Bytes [48]byte
	copy(xC0Bytes[:], bigPublicKey[1:49])
	copy(xC1Bytes[:], bigPublicKey[49:97])
	copy(yC0Bytes[:], bigPublicKey[97:145])
	copy(yC1Bytes[:], bigPublicKey[145:193])
	xC0 := FQReprFromBytes(xC0Bytes)
	xC1 := FQReprFromBytes(xC1Bytes)
	yC0 := FQReprFromBytes(yC0Bytes)
	yC1 := FQReprFromBytes(yC1Bytes)
	g.x.c0 = &FQ{n: xC0}
	g.x.c1 = &FQ{n: xC1}
	g.y.c0 = &FQ{n: yC0}
	g.y.c1 = &FQ{n: yC1}
	return &PublicKey{p: g.ToProjective()}
}

// Equals checks if two public keys are equal
func (p PublicKey) Equals(other PublicKey) bool {
	return p.p.Equal(other.p)
}

// DeserializePublicKey deserializes a public key from
// bytes.
func DeserializePublicKey(b [96]byte) (*PublicKey, error) {
	a, err := DecompressG2(b)
	if err != nil {
		return nil, err
	}

	return &PublicKey{p: a.ToProjective()}, nil
}

// SecretKey represents a BLS private key.
type SecretKey struct {
	f *FR
}

func (s SecretKey) String() string {
	return s.f.String()
}

// Serialize serializes a secret key to bytes.
func (s SecretKey) Serialize() [32]byte {
	return s.f.n.Bytes()
}

// DeserializeSecretKey deserializes a secret key from
// bytes.
func DeserializeSecretKey(b [32]byte) *SecretKey {
	return &SecretKey{&FR{FRReprFromBytes(b)}}
}

// Sign signs a message with a secret key.
func Sign(message []byte, key *SecretKey, domain uint64) *Signature {
	h := HashG1(message, domain).MulFR(key.f.n)
	return &Signature{s: h}
}

// PrivToPub converts the private key into a public key.
func PrivToPub(k *SecretKey) *PublicKey {
	return &PublicKey{p: G2AffineOne.MulFR(k.f.n)}
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

// KeyFromFQRepr returns a new key based on a FQRepr in
// FR.
func KeyFromFQRepr(i *FRRepr) *SecretKey {
	return &SecretKey{f: FRReprToFR(i)}
}

// Verify verifies a signature against a message and a public key.
func Verify(m []byte, pub *PublicKey, sig *Signature, domain uint64) bool {
	h := HashG1(m, domain)
	lhs := Pairing(sig.s, G2ProjectiveOne)
	rhs := Pairing(h, pub.p)
	return lhs.Equals(rhs)
}

// AggregateSignatures adds up all of the signatures.
func AggregateSignatures(s []*Signature) *Signature {
	newSig := &Signature{s: G1ProjectiveZero.Copy()}
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
	newPub := &PublicKey{p: G2ProjectiveZero.Copy()}
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
	return &Signature{s: G1ProjectiveZero.Copy()}
}

// NewAggregatePubkey creates a blank public key.
func NewAggregatePubkey() *PublicKey {
	return &PublicKey{p: G2ProjectiveZero.Copy()}
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
func (s *Signature) VerifyAggregate(pubKeys []*PublicKey, msgs [][]byte, domain uint64) bool {
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
		h := HashG1(msgs[i], domain)
		rhs.MulAssign(Pairing(h, pubKeys[i].p))
	}
	return lhs.Equals(rhs)
}

// VerifyAggregateCommon verifies each public key against a message.
// This is vulnerable to rogue public-key attack. Each user must
// provide a proof-of-knowledge of the public key.
func (s *Signature) VerifyAggregateCommon(pubKeys []*PublicKey, msg []byte, domain uint64) bool {
	h := HashG1(msg, domain)
	lhs := Pairing(s.s, G2ProjectiveOne)
	rhs := FQ12One.Copy()
	for _, p := range pubKeys {
		rhs.MulAssign(Pairing(h, p.p))
	}
	return lhs.Equals(rhs)
}
