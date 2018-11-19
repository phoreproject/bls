package bls

import (
	"errors"
	"math/big"

	blake2b "github.com/minio/blake2b-simd"
)

var hexRootX, _ = new(big.Int).SetString("21573744529824266246521972077326577680729363968861965890554801909984373949499", 10)
var hexRootY, _ = new(big.Int).SetString("16854739155576650954933913186877292401521110422362946064090026408937773542853", 10)

var hexRoot = NewFQ2([]*FQ{
	NewFQ(hexRootX, fieldModulus),
	NewFQ(hexRootY, fieldModulus),
})

var oneLsh255 = new(big.Int).Lsh(bigOne, 255)

// CompressG1 compresses a point into a single 256-bit number.
func CompressG1(pt [3]*FQ) *big.Int {
	pt2 := Normalize(pt)
	return new(big.Int).Add(pt2[0].n, new(big.Int).Mul(new(big.Int).Mod(pt[1].n, bigTwo), oneLsh255))
}

// DecompressG1 decompresses the number into a point.
func DecompressG1(pt *big.Int) [3]*FQ {
	if pt.Cmp(bigZero) == 0 {
		return [3]*FQ{
			NewFQ(bigOne, fieldModulus),
			NewFQ(bigOne, fieldModulus),
			NewFQ(bigZero, fieldModulus),
		}
	}
	x := new(big.Int).Mod(pt, oneLsh255)
	yMod2 := new(big.Int).Div(pt, oneLsh255)
	yBase := new(big.Int).Mod(new(big.Int).Add(new(big.Int).Exp(x, big.NewInt(3), nil), B.n), fieldModulus)
	yPow := new(big.Int).Div(new(big.Int).Add(fieldModulus, bigOne), big.NewInt(4))
	y := new(big.Int).Exp(yBase, yPow, fieldModulus)
	if new(big.Int).Mod(y, bigTwo).Cmp(yMod2) != 0 {
		y = new(big.Int).Sub(fieldModulus, y)
	}
	return [3]*FQ{
		NewFQ(x, fieldModulus),
		NewFQ(y, fieldModulus),
		NewFQ(bigOne, fieldModulus),
	}
}

// CompressG2 takes a FQ2 point and compresses it to a single
// big.Int.
func CompressG2(pt [3]*FQP) ([2]*big.Int, error) {
	onCurve, err := IsOnCurveFQP(pt, B2)
	if err != nil {
		return [2]*big.Int{}, err
	}
	if !onCurve {
		return [2]*big.Int{}, errors.New("point is not on B2")
	}
	ptNorm := NormalizeFQP(pt)
	x, y := ptNorm[0], ptNorm[1]
	xPart := new(big.Int).Mul(oneLsh255, new(big.Int).Mod(y.elements[0].n, bigTwo))
	xPart.Add(xPart, x.elements[0].n)
	return [2]*big.Int{xPart, x.elements[1].n}, nil
}

var sqrtFQ2Exponent = new(big.Int).Div(new(big.Int).Add(new(big.Int).Exp(fieldModulus, bigTwo, nil), big.NewInt(15)), big.NewInt(32))

// SqrtFQ2 square-roots an FQ2 element.
func SqrtFQ2(x *FQP) *FQP {
	y := x.Exp(sqrtFQ2Exponent)
	for !y.Exp(bigTwo).Equals(x) {
		y = y.Mul(hexRoot)
	}
	return y
}

// DecompressG2 takes a compressed point and converts it to
// a point with elements FQ2.
func DecompressG2(compressed [2]*big.Int) ([3]*FQP, error) {
	x1 := new(big.Int).Mod(compressed[0], oneLsh255)
	y1Mod2 := new(big.Int).Div(compressed[0], oneLsh255)
	x2 := compressed[1]
	x := NewFQ2([]*FQ{
		NewFQ(x1, fieldModulus),
		NewFQ(x2, fieldModulus),
	})
	if x.Equals(FQ2Zero()) {
		return [3]*FQP{
			FQ2One(),
			FQ2One(),
			FQ2Zero(),
		}, nil
	}
	y := SqrtFQ2(x.Exp(big.NewInt(3)).Add(B2))
	if new(big.Int).Mod(y.elements[0].n, bigTwo).Cmp(y1Mod2) != 0 {
		y = y.MulScalar(NewFQ(big.NewInt(-1), fieldModulus))
	}
	out := [3]*FQP{
		x,
		y,
		FQ2One(),
	}
	onCurve, err := IsOnCurveFQP(out, B2)
	if err != nil {
		return [3]*FQP{}, err
	}
	if !onCurve {
		return [3]*FQP{}, errors.New("decompressed point is not on curve b2")
	}
	return out, nil
}

// Blake calculates the blake2b hash of the given bytes.
func Blake(x []byte) Hash {
	return Hash(blake2b.Sum256(x))
}

var hashToG2Exponent = new(big.Int).Div(new(big.Int).Add(new(big.Int).Exp(fieldModulus, bigTwo, nil), bigOne), bigTwo)

// HashToG2 converts a 256-bit hash into a point on G2.
func HashToG2(m Hash) [3]*FQP {
	// TODO: spec not complete yet
	one := FQ2One()
	k2 := m
	var xcb *FQP
	var x *FQP
	for {
		k1 := Blake(k2[:])
		k2 := Blake(k1[:])
		x1 := new(big.Int).SetBytes(k1[:])
		x2 := new(big.Int).SetBytes(k2[:])
		x = NewFQ2([]*FQ{
			NewFQ(x1, fieldModulus),
			NewFQ(x2, fieldModulus),
		})
		xcb = x.Exp(big.NewInt(3)).Add(B2)
		if xcb.Exp(hashToG2Exponent).Equals(one) {
			break
		}
	}
	y := SqrtFQ2(xcb)
	originalPoint := [3]*FQP{
		x,
		y,
		FQ2One(),
	}
	mulFactor := new(big.Int).Sub(new(big.Int).Mul(bigTwo, fieldModulus), curveOrder)
	o := MultiplyFQP(originalPoint, mulFactor)
	return o
}

// Hash represents any 256-bit hash output.
type Hash [32]byte

// Sign signs a message with a private key.
func Sign(m Hash, k *big.Int) ([2]*big.Int, error) {
	return CompressG2(MultiplyFQP(HashToG2(m), k))
}

// PrivToPub converts a private key to a public key.
func PrivToPub(priv *big.Int) *big.Int {
	return CompressG1(MultiplyFQ(G1, priv))
}

// Verify verifies e(sig, g) = e(H(m), g^x).
func Verify(m Hash, pub *big.Int, sig [2]*big.Int) (bool, error) {
	signatureG2, err := DecompressG2(sig)
	if err != nil {
		return false, err
	}
	left, err := Pairing(signatureG2, G1, false)
	if err != nil {
		return false, err
	}
	right, err := Pairing(HashToG2(m), NegFQ(DecompressG1(pub)), false)
	if err != nil {
		return false, err
	}
	finalExponentiation := FinalExponentiateFQP(left.Mul(right))
	return finalExponentiation.Equals(FQ12One()), nil
}

// AggregatePubs aggregates multiple public keys together.
func AggregatePubs(pubs []*big.Int) *big.Int {
	o := z1
	for _, p := range pubs {
		o = AddFQ(o, DecompressG1(p))
	}
	return CompressG1(o)
}

// AggregateSigs aggregates multiple signatures into a
// single signature.
func AggregateSigs(sigs [][2]*big.Int) ([2]*big.Int, error) {
	o := z2
	for _, s := range sigs {
		pt, err := DecompressG2(s)
		if err != nil {
			return [2]*big.Int{}, err
		}
		o = AddFQP(o, pt)
	}
	return CompressG2(o)
}
