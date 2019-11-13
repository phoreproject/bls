package bls

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"math/big"
)

// G2Affine is an affine point on the G2 curve.
type G2Affine struct {
	x        FQ2
	y        FQ2
	infinity bool
}

// NewG2Affine constructs a new G2Affine point.
func NewG2Affine(x FQ2, y FQ2) *G2Affine {
	return &G2Affine{x: x, y: y, infinity: false}
}

// G2AffineZero represents the point at infinity on G2.
var G2AffineZero = &G2Affine{FQ2Zero, FQ2One, true}

var g2GeneratorXC1, _ = FQReprFromString("13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e", 16)
var g2GeneratorXC0, _ = FQReprFromString("24aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8", 16)
var g2GeneratorYC1, _ = FQReprFromString("606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be", 16)
var g2GeneratorYC0, _ = FQReprFromString("ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801", 16)

// BCoeffFQ2 of the G2 curve.
var BCoeffFQ2 = NewFQ2(BCoeff, BCoeff)

// G2AffineOne represents the point at 1 on G2.
var G2AffineOne = &G2Affine{
	x: NewFQ2(
		FQReprToFQ(g2GeneratorXC0),
		FQReprToFQ(g2GeneratorXC1),
	),
	y: NewFQ2(
		FQReprToFQ(g2GeneratorYC0),
		FQReprToFQ(g2GeneratorYC1),
	), infinity: false}

func (g G2Affine) String() string {
	if g.infinity {
		return fmt.Sprintf("G2(Infinity)")
	}
	return fmt.Sprintf("G2(x=%s, y=%s)", g.x, g.y)
}

// Copy returns a copy of the G2Affine point.
func (g G2Affine) Copy() *G2Affine {
	return &G2Affine{g.x.Copy(), g.y.Copy(), g.infinity}
}

// IsZero checks if the point is infinity.
func (g G2Affine) IsZero() bool {
	return g.infinity
}

// NegAssign negates the point.
func (g *G2Affine) NegAssign() {
	if !g.IsZero() {
		g.y.NegAssign()
	}
}

// ToProjective converts an affine point to a projective one.
func (g G2Affine) ToProjective() *G2Projective {
	if g.IsZero() {
		return G2ProjectiveZero.Copy()
	}

	return NewG2Projective(g.x, g.y, FQ2One)
}

// Mul performs a EC multiply operation on the point.
func (g G2Affine) Mul(b FQRepr) *G2Projective {
	res := G2ProjectiveZero.Copy()
	for i := uint(0); i < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.AddAffine(&g)
		}
	}
	return res
}

// MulFR performs a EC multiply operation on the point.
func (g G2Affine) MulFR(b *FRRepr) *G2Projective {
	res := G2ProjectiveZero.Copy()
	for i := uint(0); i < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.AddAffine(&g)
		}
	}
	return res
}

// MulBig performs a EC multiply operation on the point.
func (g G2Affine) MulBig(b big.Int) *G2Projective {
	res := G2ProjectiveZero.Copy()
	for i := 0; i < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o == 1 {
			res = res.AddAffine(&g)
		}
	}
	return res
}

// IsOnCurve checks if a point is on the G2 curve.
func (g G2Affine) IsOnCurve() bool {
	if g.infinity {
		return true
	}
	y2 := g.y.Copy()
	y2.SquareAssign()
	x3b := g.x.Copy()
	x3b.SquareAssign()
	x3b.MulAssign(g.x)
	x3b.AddAssign(BCoeffFQ2)

	return y2.Equals(x3b)
}

// G2 cofactor = (x^8 - 4 x^7 + 5 x^6) - (4 x^4 + 6 x^3 - 4 x^2 - 4 x + 13) // 9
var g2Cofactor, _ = new(big.Int).SetString("5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5", 16)

// ScaleByCofactor scales the G2Affine point by the cofactor.
func (g G2Affine) ScaleByCofactor() *G2Projective {
	return g.MulBig(*g2Cofactor)
}

// Equals checks if two affine points are equal.
func (g G2Affine) Equals(other *G2Affine) bool {
	return (g.infinity == other.infinity) || (g.x.Equals(other.x) && g.y.Equals(other.y))
}

// GetG2PointFromX attempts to reconstruct an affine point given
// an x-coordinate. The point is not guaranteed to be in the subgroup.
// If and only if `greatest` is set will the lexicographically
// largest y-coordinate be selected.
func GetG2PointFromX(x FQ2, greatest bool) (*G2Affine, error) {
	x3b := x.Copy()
	x3b.SquareAssign()
	x3b.MulAssign(x)
	x3b.AddAssign(BCoeffFQ2)

	y, success := x3b.Sqrt()

	if !success {
		return nil, errors.New("point not on curve")
	}

	negY := y.Copy()
	negY.NegAssign()

	yVal := negY
	if (y.Cmp(negY) < 0) != greatest {
		yVal = y
	}
	return NewG2Affine(x, yVal), nil
}

// SerializeBytes returns the serialized bytes for the points represented.
func (g *G2Affine) SerializeBytes() [192]byte {
	out := [192]byte{}

	xC0Bytes := g.x.c0.ToRepr().Bytes()
	xC1Bytes := g.x.c1.ToRepr().Bytes()
	yC0Bytes := g.y.c0.ToRepr().Bytes()
	yC1Bytes := g.y.c1.ToRepr().Bytes()

	copy(out[0:48], xC0Bytes[:])
	copy(out[48:96], xC1Bytes[:])
	copy(out[96:144], yC0Bytes[:])
	copy(out[144:192], yC1Bytes[:])

	return out
}

// SetRawBytes sets the coords given the serialized bytes.
func (g *G2Affine) SetRawBytes(uncompressed [192]byte) error {

	var xC0Bytes [48]byte
	var xC1Bytes [48]byte
	var yC0Bytes [48]byte
	var yC1Bytes [48]byte

	copy(xC0Bytes[:], uncompressed[0:48])
	copy(xC1Bytes[:], uncompressed[48:96])
	copy(yC0Bytes[:], uncompressed[96:144])
	copy(yC1Bytes[:], uncompressed[144:192])

	xc0FQ := FQReprToFQ(FQReprFromBytes(xC0Bytes))
	xc1FQ := FQReprToFQ(FQReprFromBytes(xC1Bytes))
	yc0FQ := FQReprToFQ(FQReprFromBytes(yC0Bytes))
	yc1FQ := FQReprToFQ(FQReprFromBytes(yC1Bytes))

	g.x = FQ2{
		c0: xc0FQ,
		c1: xc1FQ,
	}
	g.y = FQ2{
		c0: yc0FQ,
		c1: yc1FQ,
	}
	return nil
}

// DecompressG2 decompresses a G2 point from a big int and checks
// if it is in the correct subgroup.
func DecompressG2(c [96]byte) (*G2Affine, error) {
	affine, err := DecompressG2Unchecked(c)
	if err != nil {
		return nil, err
	}

	if !affine.IsInCorrectSubgroupAssumingOnCurve() {
		return nil, errors.New("point is not in correct subgroup")
	}
	return affine, nil
}

// DecompressG2Unchecked decompresses a G2 point from a big int.
func DecompressG2Unchecked(c [96]byte) (*G2Affine, error) {
	if c[0]&(1<<7) == 0 {
		return nil, errors.New("unexpected compression mode")
	}

	if c[0]&(1<<6) != 0 {
		c[0] &= 0x3f

		for _, b := range c {
			if b != 0 {
				return nil, errors.New("unexpected information in infinity point on G2")
			}
		}
		return G2AffineZero.Copy(), nil
	}
	greatest := c[0]&(1<<5) != 0

	c[0] &= 0x1f

	var xC0Bytes [48]byte
	var xC1Bytes [48]byte

	copy(xC1Bytes[:], c[:48])
	copy(xC0Bytes[:], c[48:])

	xC0 := FQReprFromBytes(xC0Bytes)
	xC0FQ := FQReprToFQ(xC0)
	xC1 := FQReprFromBytes(xC1Bytes)
	xC1FQ := FQReprToFQ(xC1)

	x := NewFQ2(xC0FQ, xC1FQ)

	return GetG2PointFromX(x, greatest)
}

// CompressG2 compresses a G2 point into an int.
func CompressG2(affine *G2Affine) [96]byte {
	res := [96]byte{}
	if affine.IsZero() {
		res[0] |= 1 << 6
	} else {
		out0 := affine.x.c1.ToRepr().Bytes()
		out1 := affine.x.c0.ToRepr().Bytes()
		copy(res[:48], out0[:])
		copy(res[48:], out1[:])

		negY := affine.y.Copy()
		negY.NegAssign()

		if affine.y.Cmp(negY) > 0 {
			res[0] |= 1 << 5
		}
	}

	res[0] |= 1 << 7

	return res
}

// IsInCorrectSubgroupAssumingOnCurve checks if the point multiplied by the
// field characteristic equals zero.
func (g G2Affine) IsInCorrectSubgroupAssumingOnCurve() bool {
	return g.MulFR(frChar).IsZero()
}

// G2Projective is a projective point on the G2 curve.
type G2Projective struct {
	x FQ2
	y FQ2
	z FQ2
}

// NewG2Projective creates a new G2Projective point.
func NewG2Projective(x FQ2, y FQ2, z FQ2) *G2Projective {
	return &G2Projective{x, y, z}
}

// G2ProjectiveZero is the point at infinity where Z = 0.
var G2ProjectiveZero = &G2Projective{FQ2Zero.Copy(), FQ2One.Copy(), FQ2Zero.Copy()}

// G2ProjectiveOne is the generator point on G2.
var G2ProjectiveOne = G2AffineOne.ToProjective()

func (g G2Projective) String() string {
	if g.IsZero() {
		return "G2: Infinity"
	}
	return g.ToAffine().String()
}

// Copy returns a copy of the G2Projective point.
func (g G2Projective) Copy() *G2Projective {
	return NewG2Projective(g.x.Copy(), g.y.Copy(), g.z.Copy())
}

// IsZero checks if the G2Projective point is zero.
func (g G2Projective) IsZero() bool {
	return g.z.IsZero()
}

// Equal checks if two projective points are equal.
func (g G2Projective) Equals(other *G2Projective) bool {
	if g.IsZero() {
		return other.IsZero()
	}
	if other.IsZero() {
		return false
	}

	z1 := g.z.Copy()
	z1.SquareAssign()
	z2 := other.z.Copy()
	z2.SquareAssign()

	tmp1 := g.x.Copy()
	tmp1.MulAssign(z2)
	tmp2 := other.x.Copy()
	tmp2.MulAssign(z1)
	if !tmp1.Equals(tmp2) {
		return false
	}
	lhs := z1
	lhs.MulAssign(g.z)
	lhs.MulAssign(other.y)

	rhs := z2
	rhs.MulAssign(other.z)
	rhs.MulAssign(g.y)

	return lhs.Equals(rhs)
}

// ToAffine converts a G2Projective point to affine form.
func (g G2Projective) ToAffine() *G2Affine {
	if g.IsZero() {
		return G2AffineZero
	} else if g.z.Equals(FQ2One) {
		return NewG2Affine(g.x, g.y)
	}

	// nonzero so must have an inverse
	zInv := g.z.Copy()
	zInv.InverseAssign()
	zInvSquared := zInv.Copy()
	zInvSquared.SquareAssign()

	x := g.x.Copy()
	x.MulAssign(zInvSquared)

	y := g.y.Copy()
	y.MulAssign(zInvSquared)
	y.MulAssign(zInv)

	return NewG2Affine(x, y)
}

// Double performs EC doubling on the point.
func (g G2Projective) Double() *G2Projective {
	if g.IsZero() {
		return g.Copy()
	}

	// A = x1^2
	a := g.x.Copy()
	a.SquareAssign()

	// B = y1^2
	b := g.y.Copy()
	b.SquareAssign()

	// C = B^2
	c := b.Copy()
	c.SquareAssign()

	// D = 2*((X1+B)^2-A-C)
	d := g.x.Copy()
	d.AddAssign(b)
	d.SquareAssign()
	d.SubAssign(a)
	d.SubAssign(c)
	d.DoubleAssign()

	// E = 3*A
	e := a.Copy()
	e.DoubleAssign()
	e.AddAssign(a)

	// F = E^2
	f := e.Copy()
	f.SquareAssign()

	// z3 = 2*Y1*Z1
	newZ := g.z.Copy()
	newZ.MulAssign(g.y)
	newZ.DoubleAssign()

	// x3 = F-2*D
	newX := f.Copy()
	newX.SubAssign(d)
	newX.SubAssign(d)

	c.DoubleAssign()
	c.DoubleAssign()
	c.DoubleAssign()

	newY := d.Copy()
	newY.SubAssign(newX)
	newY.MulAssign(e)
	newY.SubAssign(c)

	return NewG2Projective(newX, newY, newZ)
}

// Add performs an EC Add operation with another point.
func (g G2Projective) Add(other *G2Projective) *G2Projective {
	if g.IsZero() {
		return other.Copy()
	}
	if other.IsZero() {
		return g.Copy()
	}

	// Z1Z1 = Z1^2
	z1z1 := g.z.Copy()
	z1z1.SquareAssign()

	// Z2Z2 = Z2^2
	z2z2 := other.z.Copy()
	z2z2.SquareAssign()

	// U1 = X1*Z2Z2
	u1 := g.x.Copy()
	u1.MulAssign(z2z2)

	// U2 = x2*Z1Z1
	u2 := other.x.Copy()
	u2.MulAssign(z1z1)

	// S1 = Y1*Z2*Z2Z2
	s1 := g.y.Copy()
	s1.MulAssign(other.z)
	s1.MulAssign(z2z2)

	// S2 = Y2*Z1*Z1Z1
	s2 := other.y.Copy()
	s2.MulAssign(g.z)
	s2.MulAssign(z1z1)

	if u1.Equals(u2) && s1.Equals(s2) {
		// points are equal
		return g.Double()
	}

	// H = U2-U1
	h := u2.Copy()
	h.SubAssign(u1)

	// I = (2*H)^2
	i := h.Copy()
	i.DoubleAssign()
	i.SquareAssign()

	// J = H * I
	j := h.Copy()
	j.MulAssign(i)

	// r = 2*(S2-S1)
	r := s2.Copy()
	r.SubAssign(s1)
	r.DoubleAssign()

	// v = U1*I
	u1.MulAssign(i)

	// X3 = r^2 - J - 2*V
	newX := r.Copy()
	newX.SquareAssign()
	newX.SubAssign(j)
	newX.SubAssign(u1)
	newX.SubAssign(u1)

	// Y3 = r*(V - X3) - 2*S1*J
	u1.SubAssign(newX)
	u1.MulAssign(r)
	s1.MulAssign(j)
	s1.DoubleAssign()
	u1.SubAssign(s1)

	// Z3 = ((Z1+Z2)^2 - Z1Z1 - Z2Z2)*H
	newZ := g.z.Copy()
	newZ.AddAssign(other.z)
	newZ.SquareAssign()
	newZ.SubAssign(z1z1)
	newZ.SubAssign(z2z2)
	newZ.MulAssign(h)

	return NewG2Projective(newX, u1, newZ)
}

// AddAffine performs an EC Add operation with an affine point.
func (g G2Projective) AddAffine(other *G2Affine) *G2Projective {
	if g.IsZero() {
		return other.ToProjective()
	}
	if other.IsZero() {
		return g.Copy()
	}

	// Z1Z1 = Z1^2
	z1z1 := g.z.Copy()
	z1z1.SquareAssign()

	// U2 = x2*Z1Z1
	u2 := other.x.Copy()
	u2.MulAssign(z1z1)

	// S2 = Y2*Z1*Z1Z1
	s2 := other.y.Copy()
	s2.MulAssign(g.z)
	s2.MulAssign(z1z1)

	if g.x.Equals(u2) && g.y.Equals(s2) {
		// points are equal
		return g.Double()
	}

	// H = U2-X1
	u2.SubAssign(g.x)

	// HH = H^2
	hh := u2.Copy()
	hh.SquareAssign()

	// I = 4*HH
	i := hh.Copy()
	i.DoubleAssign()
	i.DoubleAssign()

	// J = H * I
	j := u2.Copy()
	j.MulAssign(i)

	// r = 2*(S2-Y1)
	s2.SubAssign(g.y)
	s2.DoubleAssign()

	// v = X1*I
	v := g.x.Copy()
	v.MulAssign(i)

	// X3 = r^2 - J - 2*V
	newX := s2.Copy()
	newX.SquareAssign()
	newX.SubAssign(j)
	newX.SubAssign(v)
	newX.SubAssign(v)

	// Y3 = r*(V - X3) - 2*Y1*J
	newY := v.Copy()
	newY.SubAssign(newX)
	newY.MulAssign(s2)
	i0 := g.y.Copy()
	i0.MulAssign(j)
	i0.DoubleAssign()
	newY.SubAssign(i0)

	// Z3 = (Z1+H)^2 - Z1Z1 - HH
	newZ := g.z.Copy()
	newZ.AddAssign(u2)
	newZ.SquareAssign()
	newZ.SubAssign(z1z1)
	newZ.SubAssign(hh)

	return NewG2Projective(newX, newY, newZ)
}

// Mul performs a EC multiply operation on the point.
func (g G2Projective) Mul(b FQRepr) *G2Projective {
	res := G2ProjectiveZero.Copy()
	for i := uint(0); i < uint(b.BitLen()); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.Add(&g)
		}
	}
	return res
}

// MulFR performs a EC multiply operation on the point.
func (g G2Projective) MulFR(b *FRRepr) *G2Projective {
	res := G2ProjectiveZero.Copy()
	for i := uint(0); i < uint(b.BitLen()); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.Add(&g)
		}
	}
	return res
}

var blsX, _ = FQReprFromString("d201000000010000", 16)

const blsIsNegative = true

// G2Prepared is a prepared G2 point multiplication by blsX.
type G2Prepared struct {
	coeffs   [][3]FQ2
	infinity bool
}

// IsZero checks if the point is at infinity.
func (g G2Prepared) IsZero() bool {
	return g.infinity
}

// G2AffineToPrepared performs multiplication of the affine point by blsX.
func G2AffineToPrepared(q *G2Affine) *G2Prepared {
	if q.IsZero() {
		return &G2Prepared{infinity: true}
	}

	doublingStep := func(r *G2Projective) (FQ2, FQ2, FQ2) {
		tmp0 := r.x.Copy()
		tmp0.SquareAssign()
		tmp1 := r.y.Copy()
		tmp1.SquareAssign()
		tmp2 := tmp1.Copy()
		tmp2.SquareAssign()
		tmp3 := tmp1.Copy()
		tmp3.AddAssign(r.x)
		tmp3.SquareAssign()
		tmp3.SubAssign(tmp0)
		tmp3.SubAssign(tmp2)
		tmp3.DoubleAssign()
		tmp4 := tmp0.Copy()
		tmp4.DoubleAssign()
		tmp4.AddAssign(tmp0)
		tmp6 := r.x.Copy()
		tmp6.AddAssign(tmp4)
		tmp5 := tmp4.Copy()
		tmp5.SquareAssign()
		zSquared := r.z.Copy()
		zSquared.SquareAssign()
		r.x = tmp5.Copy()
		r.x.SubAssign(tmp3)
		r.x.SubAssign(tmp3)
		r.z = r.z.Copy()
		r.z.AddAssign(r.y)
		r.z.SquareAssign()
		r.z.SubAssign(tmp1)
		r.z.SubAssign(zSquared)
		r.y = tmp3.Copy()
		r.y.SubAssign(r.x)
		r.y.MulAssign(tmp4)
		tmp2.DoubleAssign()
		tmp2.DoubleAssign()
		tmp2.DoubleAssign()
		r.y = r.y.Copy()
		r.y.SubAssign(tmp2)
		tmp3 = tmp4.Copy()
		tmp3.MulAssign(zSquared)
		tmp3.DoubleAssign()
		tmp3.NegAssign()
		tmp6 = tmp6.Copy()
		tmp6.SquareAssign()
		tmp6.SubAssign(tmp0)
		tmp6.SubAssign(tmp5)
		tmp1.DoubleAssign()
		tmp1.DoubleAssign()
		tmp6.SubAssign(tmp1)
		tmp0 = r.z.Copy()
		tmp0.MulAssign(zSquared)
		tmp0.DoubleAssign()
		return tmp0, tmp3, tmp6
	}

	additionStep := func(r *G2Projective, q *G2Affine) (FQ2, FQ2, FQ2) {
		zSquared := r.z.Copy()
		zSquared.SquareAssign()
		ySquared := q.y.Copy()
		ySquared.SquareAssign()
		t0 := zSquared.Copy()
		t0.MulAssign(q.x)
		t1 := q.y.Copy()
		t1.AddAssign(r.z)
		t1.SquareAssign()
		t1.SubAssign(ySquared)
		t1.SubAssign(zSquared)
		t1.MulAssign(zSquared)
		t2 := t0.Copy()
		t2.SubAssign(r.x)
		t3 := t2.Copy()
		t3.SquareAssign()
		t4 := t3.Copy()
		t4.DoubleAssign()
		t4.DoubleAssign()
		t5 := t4.Copy()
		t5.MulAssign(t2)
		t6 := t1.Copy()
		t6.SubAssign(r.y)
		t6.SubAssign(r.y)
		t9 := t6.Copy()
		t9.MulAssign(q.x)
		t7 := t4.Copy()
		t7.MulAssign(r.x)
		r.x = t6.Copy()
		r.x.SquareAssign()
		r.x.SubAssign(t5)
		r.x.SubAssign(t7)
		r.x.SubAssign(t7)
		r.z = r.z.Copy()
		r.z.AddAssign(t2)
		r.z.SquareAssign()
		r.z.SubAssign(zSquared)
		r.z.SubAssign(t3)
		t10 := q.y.Copy()
		t10.AddAssign(r.z)
		t8 := t7.Copy()
		t8.SubAssign(r.x)
		t8.MulAssign(t6)
		t0 = r.y.Copy()
		t0.MulAssign(t5)
		t0.DoubleAssign()
		r.y = t8.Copy()
		r.y.SubAssign(t0)
		t10.SquareAssign()
		t10.SubAssign(ySquared)
		zSquared = r.z.Copy()
		zSquared.SquareAssign()
		t10.SubAssign(zSquared)
		t9.DoubleAssign()
		t9.SubAssign(t10)
		t10 = r.z.Copy()
		t10.DoubleAssign()
		t6.NegAssign()
		t6.DoubleAssign()

		return t10, t6, t9
	}

	var coeffs [][3]FQ2
	r := q.ToProjective()

	foundOne := false
	blsXRsh1 := blsX.Copy()
	blsXRsh1.Rsh(1)
	for i := uint(0); i <= blsXRsh1.BitLen(); i++ {
		set := blsXRsh1.Bit(blsXRsh1.BitLen() - i)
		if !foundOne {
			foundOne = set
			continue
		}

		o0, o1, o2 := doublingStep(r)

		coeffs = append(coeffs, [3]FQ2{o0, o1, o2})

		if set {
			o0, o1, o2 := additionStep(r, q)
			coeffs = append(coeffs, [3]FQ2{o0, o1, o2})
		}
	}
	o0, o1, o2 := doublingStep(r)

	coeffs = append(coeffs, [3]FQ2{o0, o1, o2})

	return &G2Prepared{coeffs, false}
}

// RandG2 generates a random G2 element.
func RandG2(r io.Reader) (*G2Projective, error) {
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			return nil, err
		}
		greatest := false
		if b[0]%2 == 0 {
			greatest = true
		}
		f, err := RandFQ2(r)
		if err != nil {
			return nil, err
		}
		p, err := GetG2PointFromX(f, greatest)
		if err != nil {
			continue
		}
		p1 := p.ScaleByCofactor()
		if !p.IsZero() {
			return p1, nil
		}
	}
}

var swencSqrtNegThree, _ = FQReprFromString("1586958781458431025242759403266842894121773480562120986020912974854563298150952611241517463240701", 10)
var swencSqrtNegThreeMinusOneDivTwo, _ = FQReprFromString("793479390729215512621379701633421447060886740281060493010456487427281649075476305620758731620350", 10)
var swencSqrtNegThreeFQ2 = NewFQ2(FQReprToFQ(swencSqrtNegThree), FQZero.Copy())
var swencSqrtNegThreeMinusOneDivTwoFQ2 = NewFQ2(FQReprToFQ(swencSqrtNegThreeMinusOneDivTwo), FQZero.Copy())

// SWEncodeG2 implements the Shallue-van de Woestijne encoding.
func SWEncodeG2(t FQ2) *G2Affine {
	if t.IsZero() {
		return G2AffineZero.Copy()
	}

	parity := t.Parity()

	w := t.Copy()
	w.SquareAssign()
	w.AddAssign(BCoeffFQ2)
	w.AddAssign(FQ2One)

	if w.IsZero() {
		ret := G2AffineOne.Copy()
		if parity {
			ret.NegAssign()
		}
		return ret
	}

	w.InverseAssign()
	w.MulAssign(swencSqrtNegThreeFQ2)
	w.MulAssign(t)

	x1 := w.Copy()
	x1.MulAssign(t)
	x1.NegAssign()
	x1.AddAssign(swencSqrtNegThreeMinusOneDivTwoFQ2)
	if p, err := GetG2PointFromX(x1, parity); err == nil {
		return p
	}

	x2 := x1.Copy()
	x2.NegAssign()
	x2.SubAssign(FQ2One)
	if p, err := GetG2PointFromX(x2, parity); err == nil {
		return p
	}

	x3 := w.Copy()
	x3.SquareAssign()
	x3.InverseAssign()
	x3.AddAssign(FQ2One)
	point, _ := GetG2PointFromX(x3, parity)
	return point
}

var ell2pA = NewFQ2(
	FQReprToFQ(NewFQRepr(0)),
	FQReprToFQ(NewFQRepr(240)),
)

var ell2pB = NewFQ2(
	FQReprToFQ(NewFQRepr(1012)),
	FQReprToFQ(NewFQRepr(1012)),
)

var qSquaredPlus7Over16 = fqReprFromHexUnchecked("1380cd6fab1fecf3b854bdc8b278c1a18b5978a3b6a3ce0f8d649df8b904b89b1700ffff04feffffcb7f3fffffffc000")

var rv1 = FQReprToFQ(fqReprFromHexUnchecked("6af0e0437ff400b6831e36d6bd17ffe48395dabc2d3435e77f76e17009241c5ee67992f72ec05f4c81084fbede3cc09"))
var qMinusRV1 = FQReprToFQ(fqReprFromHexUnchecked("135203e60180a68ee2e9c448d77a2cd91c3dedd930b1cf60ef396489f61eb45e304466cf3e67fa0af1ee7b04121bdea2"))

var rootsOfUnity = []FQ2{
	NewFQ2(FQOne, FQZero),
	NewFQ2(FQZero, FQOne),
	NewFQ2(rv1, rv1),
	NewFQ2(rv1, qMinusRV1),
}

var ev1 = FQReprToFQ(fqReprFromHexUnchecked("2c4a7244a026bd3e305cc456ad9e235ed85f8b53954258ec8186bb3d4eccef7c4ee7b8d4b9e063a6c88d0aa3e03ba01"))
var ev2 = FQReprToFQ(fqReprFromHexUnchecked("85fa8cd9105715e641892a0f9a4bb2912b58b8d32f26594c60679cc7973076dc6638358daf3514d6426a813ae01f51a"))
var qMinusEv2 = FQReprToFQ(fqReprFromHexUnchecked("11a1691ca87a753be703151549a6f1ae51c1bff7c092ad2aa12a58d47d3deeb658487ca5d660aeb255d857ec51fdb591"))
var etas = []FQ2{
	NewFQ2(ev1, FQZero),
	NewFQ2(FQZero, ev1),
	NewFQ2(ev2, ev2),
	NewFQ2(ev2, qMinusEv2),
}

// returns -1 if x is larger than -x otherwise 1
func signFQ2(f FQ2) int {
	thresh := FQReprToFQ(qMinus1Over2)
	if f.c1.Cmp(thresh) > 0 {
		return -1
	}
	if f.c1.Cmp(FQZero) > 0 {
		return 1
	}
	if f.c0.Cmp(thresh) > 0 {
		return -1
	}
	if f.c0.Cmp(FQZero) > 0 {
		return 1
	}
	return 1
}

func OptimizedSWU2MapHelper(t FQ2) *G2Affine {
	numDenCommon := fq2nqr.Copy()
	numDenCommon.SquareAssign()

	tSquared := t.Copy()
	tSquared.SquareAssign()

	t4 := tSquared.Copy()
	t4.SquareAssign()

	numDenCommon.MulAssign(t4)

	negOneTimesTSquared := fq2nqr.Copy()
	negOneTimesTSquared.MulAssign(tSquared)

	numDenCommon.AddAssign(negOneTimesTSquared)

	var x0 FQ2
	if numDenCommon.Equals(FQ2Zero) {
		xiA := fq2nqr.Copy()
		xiA.MulAssign(ell2pA)
		x0 = ell2pB.Copy()
		x0.DivAssign(xiA)
	} else {
		ell2paTimesCommon := ell2pA.Copy()
		ell2paTimesCommon.MulAssign(numDenCommon)
		negEll2pb := ell2pB.Copy()
		negEll2pb.NegAssign()
		numDenCommon.AddAssign(FQ2One)
		negEll2pb.MulAssign(numDenCommon)
		x0 = negEll2pb
		x0.DivAssign(ell2paTimesCommon)
	}

	x0Cubed := x0.Copy()
	x0Cubed.SquareAssign()
	x0Cubed.MulAssign(x0)

	ellPAX0 := ell2pA.Copy()
	ellPAX0.MulAssign(x0)

	gx0 := x0Cubed
	gx0.AddAssign(ellPAX0)
	gx0.AddAssign(ell2pB)

	sqrtGX0, found := gx0.Sqrt()
	if found {
		y0Squared := sqrtGX0.Copy()
		y0Squared.SquareAssign()

		if y0Squared.Equals(gx0) {
			signT := signFQ2(t)
			signYT := signFQ2(sqrtGX0)

			if signT != signYT {
				sqrtGX0.NegAssign()
			}

			// g(x0) = y^2 as required
			return NewG2Affine(x0, sqrtGX0)
		}
	}

	tCubed := tSquared.Copy()
	tCubed.MulAssign(t)

	t6 := tCubed.Copy()
	t6.SquareAssign()

	x1 := fq2nqr.Copy()
	x1.MulAssign(tSquared)
	x1.MulAssign(x0)

	gx1 := fq2nqr.Copy()
	gx1.SquareAssign()
	gx1.MulAssign(fq2nqr)
	gx1.MulAssign(t6)
	gx1.MulAssign(gx0)

	y1, found := gx1.Sqrt()
	if !found {
		panic("This should never happen!")
	}

	y12 := y1.Copy()
	y12.SquareAssign()
	if y12.Equals(gx1) {
		signT := signFQ2(t)
		signYT := signFQ2(y1)

		if signT != signYT {
			y1.NegAssign()
		}

		return NewG2Affine(x1, y1)
	}

	return nil
}

// hashFunc returns the SHA-256 hash of the input
func hashFunc(in []byte) []byte {
	h := sha256.New()
	h.Write(in)
	return h.Sum(nil)
}

// Warning: See BLS standardisation process. This implementation is known to be unsecure.
func HashG2WithDomain(messageHash [32]byte, domain [8]byte) *G2Projective {

	xReBytes := make([]byte, 0, 32+8+1)
	xReBytes = append(xReBytes, messageHash[:]...)
	xReBytes = append(xReBytes, domain[:]...)
	xReBytes = append(xReBytes, 0x01)

	xImBytes := make([]byte, 0, 32+8+1)
	xImBytes = append(xImBytes, messageHash[:]...)
	xImBytes = append(xImBytes, domain[:]...)
	xImBytes = append(xImBytes, 0x02)

	xRe := new(big.Int)
	xRe.SetBytes(hashFunc(xReBytes))

	xIm := new(big.Int)
	xIm.SetBytes(hashFunc(xImBytes))

	// hash function is only 256 bits so this will never overflow
	xReFQ, _ := FQReprFromBigInt(xRe)
	xImFQ, _ := FQReprFromBigInt(xIm)

	x0 := NewFQ2(FQReprToFQ(xReFQ), FQReprToFQ(xImFQ))

	for {
		gx0 := x0.Copy()
		gx0.SquareAssign()
		gx0.MulAssign(x0)

		gx0.AddAssign(BCoeffFQ2)

		y0, found := gx0.Sqrt()

		if found {
			// favor the lower y value
			if !y0.Parity() {
				y0.NegAssign()
			}

			return NewG2Affine(x0, y0).ScaleByCofactor()
		}

		x0.AddAssign(FQ2One)
	}
}
