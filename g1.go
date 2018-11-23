package bls

import (
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

var bCoeff, _ = new(big.Int).SetString("1514052131932888505822357196874193114600527104240479143842906308145652716846165732392247483508051665748635331395571", 10)

// G1AffineOne represents the point at 1 on G1.
var G1AffineOne = &G1Affine{NewFQ(g1GeneratorX), NewFQ(g1GeneratorY), false}

func (g G1Affine) String() string {
	if g.infinity {
		return fmt.Sprintf("G1: infinity")
	}
	return fmt.Sprintf("G1: (%s, %s)", g.x, g.y)
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
	for i := uint(0); i < uint(b.BitLen()); i++ {
		o := bs[i/8]&(1<<(i%8)) > 0
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
	x3b := g.x.Square().Mul(g.x).Add(NewFQ(bCoeff))

	return y2.Equals(x3b)
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
	return fmt.Sprintf("G1: (%s, %s, %s)", g.x, g.y, g.z)
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
	for i := uint(0); i < uint(b.BitLen()); i++ {
		o := bs[i/8]&(1<<(i%8)) > 0
		res = res.Double()
		if o {
			res = res.Add(&g)
		}
	}
	return res
}
