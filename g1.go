package bls

import (
	"errors"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/blake2b"
)

// G1Affine is an affine point on the G1 curve.
type G1Affine struct {
	x        *FQ
	y        *FQ
	infinity bool
}

// NewG1Affine constructs a new G1Affine point.
func NewG1Affine(x *FQ, y *FQ) *G1Affine {
	return &G1Affine{x: x, y: y, infinity: false}
}

// G1AffineZero represents the point at infinity on G1.
var G1AffineZero = &G1Affine{FQZero, FQOne, true}

var g1GeneratorX, _ = new(big.Int).SetString("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507", 10)
var g1GeneratorY, _ = new(big.Int).SetString("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569", 10)

// BCoeff of the G1 curve.
var BCoeff = big.NewInt(4)

// G1AffineOne represents the point at 1 on G1.
var G1AffineOne = &G1Affine{NewFQ(g1GeneratorX), NewFQ(g1GeneratorY), false}

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

// Neg negates the point.
func (g G1Affine) Neg() *G1Affine {
	if !g.IsZero() {
		return NewG1Affine(g.x, g.y.Neg())
	}
	return g.Copy()
}

// NegAssign negates the point.
func (g G1Affine) NegAssign() {
	if !g.IsZero() {
		g.y.NegAssign()
	}
}

// ToProjective converts an affine point to a projective one.
func (g G1Affine) ToProjective() *G1Projective {
	if g.IsZero() {
		return G1ProjectiveZero
	}
	return NewG1Projective(g.x, g.y, FQOne)
}

// Mul performs a EC multiply operation on the point.
func (g G1Affine) Mul(b *big.Int) *G1Projective {
	bs := b.Bytes()
	res := G1ProjectiveZero.Copy()
	for i := uint(0); i < uint((b.BitLen()+7)/8)*8; i++ {
		part := i / 8
		bit := 7 - i%8
		o := bs[part]&(1<<(bit)) > 0
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
	y2 := g.y.Square()
	x3b := g.x.Square().Mul(g.x).Add(NewFQ(BCoeff))

	return y2.Equals(x3b)
}

// GetG1PointFromX attempts to reconstruct an affine point given
// an x-coordinate. The point is not guaranteed to be in the subgroup.
// If and only if `greatest` is set will the lexicographically
// largest y-coordinate be selected.
func GetG1PointFromX(x *FQ, greatest bool) *G1Affine {
	x3b := x.Square().Mul(x).Add(NewFQ(BCoeff))

	y := x3b.Sqrt()

	if y == nil {
		return nil
	}

	negY := y.Neg()

	yVal := negY
	if (y.Cmp(negY) < 0) != greatest {
		yVal = y
	}
	return NewG1Affine(x, yVal)
}

var frChar, _ = new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)

// IsInCorrectSubgroupAssumingOnCurve checks if the point multiplied by the
// field characteristic equals zero.
func (g G1Affine) IsInCorrectSubgroupAssumingOnCurve() bool {
	return g.Mul(frChar).IsZero()
}

// G1 cofactor = (x - 1)^2 / 3  = 76329603384216526031706109802092473003
var g1Cofactor, _ = new(big.Int).SetString("76329603384216526031706109802092473003", 10)

// ScaleByCofactor scales the G1Affine point by the cofactor.
func (g G1Affine) ScaleByCofactor() *G1Projective {
	return g.Mul(g1Cofactor)
}

// Equals checks if two affine points are equal.
func (g G1Affine) Equals(other *G1Affine) bool {
	return (g.infinity == other.infinity) || (g.x.Equals(other.x) && g.y.Equals(other.y))
}

// DecompressG1 decompresses the big int into an affine point and checks
// if it is in the correct prime group.
func DecompressG1(b *big.Int) (*G1Affine, error) {
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
func DecompressG1Unchecked(b *big.Int) (*G1Affine, error) {
	copy := new(big.Int).Set(b)

	copyBytes := copy.Bytes()

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

	x := NewFQ(new(big.Int).SetBytes(copyBytes))

	return GetG1PointFromX(x, greatest), nil
}

// CompressG1 compresses a G1 point into an int.
func CompressG1(affine *G1Affine) *big.Int {
	res := [48]byte{}
	if affine.IsZero() {
		res[0] |= 1 << 6
	} else {
		out0 := new(big.Int).Set(affine.x.n).Bytes()
		copy(res[:], out0)

		negY := affine.y.Neg()

		if affine.y.Cmp(negY) > 0 {
			res[0] |= 1 << 5
		}
	}

	res[0] |= 1 << 7
	return new(big.Int).SetBytes(res[:])
}

// G1Projective is a projective point on the G1 curve.
type G1Projective struct {
	x *FQ
	y *FQ
	z *FQ
}

// NewG1Projective creates a new G1Projective point.
func NewG1Projective(x *FQ, y *FQ, z *FQ) *G1Projective {
	return &G1Projective{x, y, z}
}

// G1ProjectiveZero is the point at infinity where Z = 0.
var G1ProjectiveZero = &G1Projective{FQZero, FQOne, FQZero}

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

	z1 := g.z.Square()
	z2 := other.z.Square()

	tmp1 := g.x.Mul(z2)
	tmp2 := other.x.Mul(z1)
	if !tmp1.Equals(tmp2) {
		return false
	}

	return z1.Mul(g.z).Mul(other.y).Equals(z2.Mul(other.z).Mul(g.y))
}

// ToAffine converts a G1Projective point to affine form.
func (g G1Projective) ToAffine() *G1Affine {
	if g.IsZero() {
		return G1AffineZero
	} else if g.z.IsZero() {
		return NewG1Affine(g.x, g.y)
	}

	// nonzero so must have an inverse
	zInv := g.z.Inverse()
	zInvSquared := zInv.Square()

	return NewG1Affine(g.x.Mul(zInvSquared), g.y.Mul(zInvSquared).Mul(zInv))
}

// Double performs EC doubling on the point.
func (g G1Projective) Double() *G1Projective {
	if g.IsZero() {
		return g.Copy()
	}

	// A = x1^2
	a := g.x.Square()

	// B = y1^2
	b := g.y.Square()

	// C = B^2
	c := b.Square()

	// D = 2*((X1+B)^2-A-C)
	d := g.x.Add(b)
	d.SquareAssign()
	d.SubAssign(a)
	d.SubAssign(c)
	d.DoubleAssign()

	// E = 3*A
	e := a.Double()
	e.AddAssign(a)

	// F = E^2
	f := e.Square()

	// z3 = 2*Y1*Z1
	newZ := g.z.Mul(g.y)
	newZ.DoubleAssign()

	// x3 = F-2*D
	newX := f.Sub(d)
	newX.SubAssign(d)

	c.DoubleAssign()
	c.DoubleAssign()
	c.DoubleAssign()

	newY := d.Sub(newX)
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
	z1z1 := g.z.Square()

	// Z2Z2 = Z2^2
	z2z2 := other.z.Square()

	// U1 = X1*Z2Z2
	u1 := g.x.Mul(z2z2)

	// U2 = x2*Z1Z1
	u2 := other.x.Mul(z1z1)

	// S1 = Y1*Z2*Z2Z2
	s1 := g.y.Mul(other.z)
	s1.MulAssign(z2z2)

	// S2 = Y2*Z1*Z1Z1
	s2 := other.y.Mul(g.z)
	s2.MulAssign(z1z1)

	if u1.Equals(u2) && s1.Equals(s2) {
		// points are equal
		return g.Double()
	}

	// H = U2-U1
	h := u2.Sub(u1)

	// I = (2*H)^2
	i := h.Double()
	i.SquareAssign()

	// J = H * I
	j := h.Mul(i)

	// r = 2*(S2-S1)
	s2.SubAssign(s1)
	s2.DoubleAssign()

	// U1 = U1*I
	u1.MulAssign(i)

	// X3 = r^2 - J - 2*V
	newX := s2.Square()
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
	newZ := g.z.Add(other.z)
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
	z1z1 := g.z.Square()

	// U2 = x2*Z1Z1
	u2 := other.x.Mul(z1z1)

	// S2 = Y2*Z1*Z1Z1
	s2 := other.y.Mul(g.z).Mul(z1z1)

	if g.x.Equals(u2) && g.y.Equals(s2) {
		// points are equal
		return g.Double()
	}

	// H = U2-X1
	u2.SubAssign(g.x)

	// HH = H^2
	hh := u2.Square()

	// I = 4*HH
	i := hh.Double()
	i.DoubleAssign()

	// J = H * I
	j := u2.Mul(i)

	// r = 2*(S2-Y1)
	s2.SubAssign(g.y)
	s2.DoubleAssign()

	// v = X1*I
	v := g.x.Mul(i)

	// X3 = r^2 - J - 2*V
	newX := s2.Square()
	newX.SubAssign(j)
	newX.SubAssign(v)
	newX.SubAssign(v)

	// Y3 = r*(V - X3) - 2*Y1*J
	newY := v.Sub(newX)
	newY.MulAssign(s2)
	i0 := g.y.Mul(j)
	i0.DoubleAssign()
	newY.SubAssign(i0)

	// Z3 = (Z1+H)^2 - Z1Z1 - HH
	newZ := g.z.Add(u2)
	newZ.SquareAssign()
	newZ.SubAssign(z1z1)
	newZ.SubAssign(hh)

	return NewG1Projective(newX, newY, newZ)
}

// Mul performs a EC multiply operation on the point.
func (g G1Projective) Mul(b *big.Int) *G1Projective {
	bs := b.Bytes()
	res := G1ProjectiveZero.Copy()
	for i := uint(0); i < uint((b.BitLen()+7)/8)*8; i++ {
		part := i / 8
		bit := 7 - i%8
		o := bs[part]&(1<<(bit)) > 0
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
		p := GetG1PointFromX(f, greatest)
		if p == nil {
			continue
		}
		p1 := p.ScaleByCofactor()
		if !p.IsZero() {
			return p1, nil
		}
	}
}

// SWEncodeG1 implements the Shallue-van de Woestijne encoding.
func SWEncodeG1(t *FQ) *G1Affine {
	if t.IsZero() {
		return G1AffineZero.Copy()
	}

	parity := t.Parity()

	w := t.Square()
	w.AddAssign(NewFQ(BCoeff))
	w.AddAssign(FQOne)

	if w.IsZero() {
		ret := G1AffineOne.Copy()
		if parity {
			ret.Neg()
		}
		return ret
	}

	w = w.Inverse()
	w.MulAssign(NewFQ(swencSqrtNegThree))
	w.MulAssign(t)

	x1 := w.Mul(t)
	x1.NegAssign()
	x1.AddAssign(NewFQ(swencSqrtNegThreeMinusOneDivTwo))
	if p := GetG1PointFromX(x1, parity); p != nil {
		return p
	}

	x2 := x1.Neg()
	x2.SubAssign(FQOne)
	if p := GetG1PointFromX(x2, parity); p != nil {
		return p
	}

	x3 := w.Square()
	x3 = x3.Inverse()
	x3.AddAssign(FQOne)
	return GetG1PointFromX(x3, parity)
}

// HashG1 converts a message to a point on the G2 curve.
func HashG1(msg []byte) *G1Projective {
	hasher0, _ := blake2b.New(64, nil)
	hasher0.Write(msg)
	hasher0.Write([]byte("G2_0"))
	hasher1, _ := blake2b.New(64, nil)
	hasher0.Write(msg)
	hasher1.Write([]byte("G2_1"))
	t0 := HashFQ(hasher0)
	t0Affine := SWEncodeG1(t0)
	t1 := HashFQ(hasher1)
	t1Affine := SWEncodeG1(t1)

	res := t0Affine.ToProjective()
	res = res.AddAffine(t1Affine)
	return res.ToAffine().ScaleByCofactor()
}
