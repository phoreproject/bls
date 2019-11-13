package bls

import (
	"errors"
	"fmt"
	"io"
)

// G1Affine is an affine point on the G1 curve.
type G1Affine struct {
	x        FQ
	y        FQ
	infinity bool
}

// NewG1Affine constructs a new G1Affine point.
func NewG1Affine(x FQ, y FQ) *G1Affine {
	return &G1Affine{x: x, y: y, infinity: false}
}

// G1AffineZero represents the point at infinity on G1.
var G1AffineZero = &G1Affine{FQZero.Copy(), FQOne.Copy(), true}

// IETF BLS standard
var g1GeneratorX, _ = FQReprFromString("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507", 10)
var g1GeneratorY, _ = FQReprFromString("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569", 10)

// BCoeff of the G1 curve.
var BCoeff = FQReprToFQRaw(FQRepr{0xaa270000000cfff3, 0x53cc0032fc34000a, 0x478fe97a6b0a807f, 0xb1d37ebee6ba24d7, 0x8ec9733bbf78ab2f, 0x9d645513d83de7e})

// G1AffineOne represents the point at 1 on G1.
var G1AffineOne = &G1Affine{FQReprToFQ(g1GeneratorX), FQReprToFQ(g1GeneratorY), false}

func (g G1Affine) String() string {
	if g.infinity {
		return fmt.Sprintf("G1(infinity)")
	}
	return fmt.Sprintf("G1(x=%s, y=%s)", g.x, g.y)
}

// Copy returns a copy of the G1Affine point.
func (g G1Affine) Copy() *G1Affine {
	return &G1Affine{g.x.Copy(), g.y.Copy(), g.infinity}
}

// IsZero checks if the point is infinity.
func (g G1Affine) IsZero() bool {
	return g.infinity
}

// NegAssign negates the point.
func (g *G1Affine) NegAssign() {
	if !g.IsZero() {
		g.y.NegAssign()
	}
}

// ToProjective converts an affine point to a projective one.
func (g G1Affine) ToProjective() *G1Projective {
	if g.IsZero() {
		return G1ProjectiveZero.Copy()
	}
	return NewG1Projective(g.x, g.y, FQOne)
}

// Mul performs a EC multiply operation on the point.
func (g G1Affine) Mul(b FQRepr) *G1Projective {
	res := G1ProjectiveZero.Copy()
	for i := uint(0); uint(i) < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.AddAffine(&g)
		}
	}
	return res
}

// MulFR performs a EC multiply operation on the point.
func (g G1Affine) MulFR(b *FRRepr) *G1Projective {
	res := G1ProjectiveZero.Copy()
	for i := uint(0); uint(i) < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.AddAffine(&g)
		}
	}
	return res
}

// IsOnCurve checks if a point is on the G1 curve.
func (g G1Affine) IsOnCurve() bool {
	if g.infinity {
		return true
	}
	y2 := g.y.Copy()
	y2.SquareAssign()
	x3b := g.x.Copy()
	x3b.SquareAssign()
	x3b.MulAssign(g.x)
	x3b.AddAssign(BCoeff)

	return y2.Equals(x3b)
}

// GetG1PointFromX attempts to reconstruct an affine point given
// an x-coordinate. The point is not guaranteed to be in the subgroup.
// If and only if `greatest` is set will the lexicographically
// largest y-coordinate be selected.
func GetG1PointFromX(x FQ, greatest bool) (*G1Affine, error) {
	x3b := x.Copy()
	x3b.SquareAssign()
	x3b.MulAssign(x)
	x3b.AddAssign(BCoeff)

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
	return NewG1Affine(x, yVal), nil
}

var frChar, _ = FRReprFromString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

// IsInCorrectSubgroupAssumingOnCurve checks if the point multiplied by the
// field characteristic equals zero.
func (g G1Affine) IsInCorrectSubgroupAssumingOnCurve() bool {
	tmp := g.Copy()
	tmp.MulFR(frChar)
	return tmp.MulFR(frChar).IsZero()
}

// G1 cofactor = (x - 1)^2 / 3  = 76329603384216526031706109802092473003
var g1Cofactor, _ = FQReprFromString("76329603384216526031706109802092473003", 10)

// ScaleByCofactor scales the G1Affine point by the cofactor.
func (g G1Affine) ScaleByCofactor() *G1Projective {
	return g.Mul(g1Cofactor)
}

// Equals checks if two affine points are equal.
func (g G1Affine) Equals(other *G1Affine) bool {
	return (g.infinity == other.infinity && g.infinity == true) || (g.x.Equals(other.x) && g.y.Equals(other.y))
}

// SerializeBytes returns the serialized bytes for the point represented.
func (g *G1Affine) SerializeBytes() [96]byte {
	out := [96]byte{}

	xBytes := g.x.ToRepr().Bytes()
	yBytes := g.y.ToRepr().Bytes()

	copy(out[0:48], xBytes[:])
	copy(out[48:96], yBytes[:])

	return out
}

// SetRawBytes sets the coords given the serialized bytes.
func (g *G1Affine) SetRawBytes(uncompressed [96]byte) {
	var xBytes [48]byte
	var yBytes [48]byte
	copy(xBytes[:], uncompressed[0:48])
	copy(yBytes[:], uncompressed[48:96])
	g.x = FQ{
		n: FQReprFromBytes(xBytes),
	}
	g.y = FQ{
		n: FQReprFromBytes(yBytes),
	}
}

// DecompressG1 decompresses the big int into an affine point and checks
// if it is in the correct prime group.
func DecompressG1(b [48]byte) (*G1Affine, error) {
	affine, err := DecompressG1Unchecked(b)
	if err != nil {
		return nil, err
	}

	if !affine.IsInCorrectSubgroupAssumingOnCurve() {
		return nil, errors.New("not in correct subgroup")
	}
	return affine, nil
}

// DecompressG1Unchecked decompresses the big int into an affine point without
// checking if it's in the correct prime group.
func DecompressG1Unchecked(b [48]byte) (*G1Affine, error) {
	var copyBytes [48]byte
	copy(copyBytes[:], b[:])

	if copyBytes[0]&(1<<7) == 0 {
		return nil, errors.New("unexpected compression mode")
	}

	if copyBytes[0]&(1<<6) != 0 {
		// this is the point at infinity
		copyBytes[0] &= 0x3f

		for _, b := range copyBytes {
			if b != 0 {
				return nil, errors.New("unexpected information in compressed infinity")
			}
		}

		return G1AffineZero.Copy(), nil
	}
	greatest := copyBytes[0]&(1<<5) != 0

	copyBytes[0] &= 0x1f

	x := FQReprFromBytes(copyBytes)
	xFQ := FQReprToFQ(x)

	return GetG1PointFromX(xFQ, greatest)
}

// CompressG1 compresses a G1 point into an int.
func CompressG1(affine *G1Affine) [48]byte {
	res := [48]byte{}
	if affine.IsZero() {
		res[0] |= 1 << 6
	} else {
		out0 := affine.x.ToRepr().Bytes()

		copy(res[48-len(out0):], out0[:])

		negY := affine.y.Copy()
		negY.NegAssign()

		if affine.y.Cmp(negY) > 0 {
			res[0] |= 1 << 5
		}
	}

	res[0] |= 1 << 7
	return res
}

// G1Projective is a projective point on the G1 curve.
type G1Projective struct {
	x FQ
	y FQ
	z FQ
}

// NegAssign negates the point.
func (g *G1Projective) NegAssign() {
	g.y.NegAssign()
}

// NewG1Projective creates a new G1Projective point.
func NewG1Projective(x FQ, y FQ, z FQ) *G1Projective {
	return &G1Projective{x, y, z}
}

// G1ProjectiveZero is the point at infinity where Z = 0.
var G1ProjectiveZero = &G1Projective{FQZero.Copy(), FQOne.Copy(), FQZero.Copy()}

// G1ProjectiveOne is the generator point on G1.
var G1ProjectiveOne = G1AffineOne.ToProjective()

func (g G1Projective) String() string {
	if g.IsZero() {
		return "G1: Infinity"
	}
	return g.ToAffine().String()
}

// Copy returns a copy of the G1Projective point.
func (g G1Projective) Copy() *G1Projective {
	return NewG1Projective(g.x.Copy(), g.y.Copy(), g.z.Copy())
}

// IsZero checks if the G1Projective point is zero.
func (g G1Projective) IsZero() bool {
	return g.z.IsZero()
}

// Equal checks if two projective points are equal.
func (g G1Projective) Equal(other *G1Projective) bool {
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

	z1.MulAssign(g.z)
	z1.MulAssign(other.y)
	z2.MulAssign(other.z)
	z2.MulAssign(g.y)

	return z1.Equals(z2)
}

// ToAffine converts a G1Projective point to affine form.
func (g G1Projective) ToAffine() *G1Affine {
	if g.IsZero() {
		return G1AffineZero.Copy()
	} else if g.z.IsZero() {
		return NewG1Affine(g.x, g.y)
	}

	// nonzero so must have an inverse
	zInv, _ := g.z.Inverse()
	zInvSquared := zInv.Copy()
	zInvSquared.SquareAssign()
	x := g.x.Copy()
	x.MulAssign(zInvSquared)
	y := g.y.Copy()
	y.MulAssign(zInvSquared)
	y.MulAssign(zInv)

	return NewG1Affine(x, y)
}

// Double performs EC doubling on the point.
func (g G1Projective) Double() *G1Projective {
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

	return NewG1Projective(newX, newY, newZ)
}

// Add performs an EC Add operation with another point.
func (g G1Projective) Add(other *G1Projective) *G1Projective {
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
	s2.SubAssign(s1)
	s2.DoubleAssign()

	// U1 = U1*I
	u1.MulAssign(i)

	// X3 = r^2 - J - 2*V
	newX := s2.Copy()
	newX.SquareAssign()
	newX.SubAssign(j)
	newX.SubAssign(u1)
	newX.SubAssign(u1)

	// Y3 = r*(V - X3) - 2*S1*J
	u1.SubAssign(newX)
	u1.MulAssign(s2)
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

	return NewG1Projective(newX, u1, newZ)
}

// AddAffine performs an EC Add operation with an affine point.
func (g G1Projective) AddAffine(other *G1Affine) *G1Projective {
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

	return NewG1Projective(newX, newY, newZ)
}

// Mul performs a EC multiply operation on the point.
func (g G1Projective) Mul(b FQRepr) *G1Projective {
	res := G1ProjectiveZero.Copy()
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
func (g G1Projective) MulFR(b *FRRepr) *G1Projective {
	res := G1ProjectiveZero.Copy()
	for i := uint(0); i < b.BitLen(); i++ {
		o := b.Bit(b.BitLen() - i - 1)
		res = res.Double()
		if o {
			res = res.Add(&g)
		}
	}
	return res
}

// RandG1 generates a random G1 element.
func RandG1(r io.Reader) (*G1Projective, error) {
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
		f, err := RandFQ(r)
		if err != nil {
			return nil, err
		}
		p, err := GetG1PointFromX(f, greatest)
		if err != nil {
			continue
		}
		p1 := p.ScaleByCofactor()
		if !p.IsZero() {
			return p1, nil
		}
	}
}

var ellPARepr, _ = FQReprFromString("144698a3b8e9433d693a02c96d4982b0ea985383ee66a8d8e8981aefd881ac98936f8da0e0f97f5cf428082d584c1d", 16)
var ellPBRepr, _ = FQReprFromString("12e2908d11688030018b12e8753eee3b2016c1f0f24f4070a0b9c14fcef35ef55a23215a316ceaa5d1cc48e98e172be0", 16)

var ellPA = FQReprToFQ(ellPARepr)
var ellPB = FQReprToFQ(ellPBRepr)

// returns -1 if x is larger than -x otherwise 1
func signFQ(f FQ) FQ {
	if f.Cmp(FQReprToFQ(qMinus1Over2)) > 0 {
		return negativeOneFQ
	}
	return FQOne
}

func optimizedSWUMapHelper(t FQ) *G1Affine {
	numDenCommon := negativeOneFQ.Copy()
	numDenCommon.SquareAssign()

	tSquared := t.Copy()
	tSquared.SquareAssign()

	t2 := tSquared.Copy()
	t2.SquareAssign()

	numDenCommon.MulAssign(t2)

	negOneTimesTSquared := negativeOneFQ.Copy()
	negOneTimesTSquared.MulAssign(tSquared)

	numDenCommon.AddAssign(negOneTimesTSquared)

	var x0 FQ
	if numDenCommon.Equals(FQZero) {
		xiA := negativeOneFQ.Copy()
		xiA.MulAssign(ellPA)
		x0 = ellPB.Copy()
		x0.DivAssign(xiA)
	} else {
		ellPATimesCommon := ellPA.Copy()
		ellPATimesCommon.MulAssign(numDenCommon)
		numDenCommon.AddAssign(FQOne)
		negEllPB := ellPB.Copy()
		negEllPB.NegAssign()
		x0 = negEllPB
		x0.MulAssign(numDenCommon)
		x0.DivAssign(ellPATimesCommon)
	}

	x0Cubed := x0.Copy()
	x0Cubed.SquareAssign()
	x0Cubed.MulAssign(x0)

	ellPAX0 := ellPA.Copy()
	ellPAX0.MulAssign(x0)

	gx0 := x0Cubed
	gx0.AddAssign(ellPAX0)
	gx0.AddAssign(ellPB)

	sqrtGX0, found := gx0.Sqrt()

	var x FQ
	var y FQ

	if found {
		// g(x0) is square, so return it's sqrt
		x = x0
		y = sqrtGX0
	} else {
		// g(x0) is not square, so find g(x1) = xi * t^2 * x0
		x1 := negativeOneFQ.Copy()
		x1.MulAssign(tSquared)
		x1.MulAssign(x0)

		ellPAX1 := ellPA.Copy()
		ellPAX1.MulAssign(x1)

		gx1 := x1.Copy()
		gx1.SquareAssign()
		gx1.MulAssign(x1)
		gx1.AddAssign(ellPAX1)
		gx1.AddAssign(ellPB)

		sqrtGX1, found := gx1.Sqrt()
		if !found {
			panic("this should never happen")
		}

		x = x1
		y = sqrtGX1
	}

	signT := signFQ(t)
	signYT := signFQ(y)

	signYT.MulAssign(signT)

	y.MulAssign(signYT)

	return NewG1Affine(x, y)
}
