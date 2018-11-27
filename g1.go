package bls

import (
	"errors"
	"fmt"
	"math/big"
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
	d := g.x.Add(b).Square().Sub(a).Sub(c).Double()

	// E = 3*A
	e := a.Double().Add(a)

	// F = E^2
	f := e.Square()

	// z3 = 2*Y1*Z1
	newZ := g.z.Mul(g.y).Double()

	// x3 = F-2*D
	newX := f.Sub(d).Sub(d)

	newY := d.Sub(newX).Mul(e).Sub(c.Double().Double().Double())

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
	s1 := g.y.Mul(other.z).Mul(z2z2)

	// S2 = Y2*Z1*Z1Z1
	s2 := other.y.Mul(g.z).Mul(z1z1)

	if u1.Equals(u2) && s1.Equals(s2) {
		// points are equal
		return g.Double()
	}

	// H = U2-U1
	h := u2.Sub(u1)

	// I = (2*H)^2
	i := h.Double().Square()

	// J = H * I
	j := h.Mul(i)

	// r = 2*(S2-S1)
	r := s2.Sub(s1).Double()

	// v = U1*I
	v := u1.Mul(i)

	// X3 = r^2 - J - 2*V
	newX := r.Square().Sub(j).Sub(v).Sub(v)

	// Y3 = r*(V - X3) - 2*S1*J
	newY := v.Sub(newX).Mul(r).Sub(s1.Mul(j).Double())

	// Z3 = ((Z1+Z2)^2 - Z1Z1 - Z2Z2)*H
	newZ := g.z.Add(other.z).Square().Sub(z1z1).Sub(z2z2).Mul(h)

	return NewG1Projective(newX, newY, newZ)
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
	h := u2.Sub(g.x)

	// HH = H^2
	hh := h.Square()

	// I = 4*HH
	i := hh.Double().Double()

	// J = H * I
	j := h.Mul(i)

	// r = 2*(S2-Y1)
	r := s2.Sub(g.y).Double()

	// v = X1*I
	v := g.x.Mul(i)

	// X3 = r^2 - J - 2*V
	newX := r.Square().Sub(j).Sub(v).Sub(v)

	// Y3 = r*(V - X3) - 2*Y1*J
	newY := v.Sub(newX).Mul(r).Sub(g.y.Mul(j).Double())

	// Z3 = (Z1+H)^2 - Z1Z1 - HH
	newZ := g.z.Add(h).Square().Sub(z1z1).Sub(hh)

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
