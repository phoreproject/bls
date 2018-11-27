package bls

import (
	"errors"
	"fmt"
	"math/big"
)

// G2Affine is an affine point on the G2 curve.
type G2Affine struct {
	x        *FQ2
	y        *FQ2
	infinity bool
}

// NewG2Affine constructs a new G2Affine point.
func NewG2Affine(x *FQ2, y *FQ2) *G2Affine {
	return &G2Affine{x: x, y: y, infinity: false}
}

// G2AffineZero represents the point at infinity on G2.
var G2AffineZero = &G2Affine{FQ2Zero, FQ2One, true}

var g2GeneratorXC1, _ = new(big.Int).SetString("3059144344244213709971259814753781636986470325476647558659373206291635324768958432433509563104347017837885763365758", 10)
var g2GeneratorXC0, _ = new(big.Int).SetString("352701069587466618187139116011060144890029952792775240219908644239793785735715026873347600343865175952761926303160", 10)
var g2GeneratorYC1, _ = new(big.Int).SetString("927553665492332455747201965776037880757740193453592970025027978793976877002675564980949289727957565575433344219582", 10)
var g2GeneratorYC0, _ = new(big.Int).SetString("1985150602287291935568054521177171638300868978215655730859378665066344726373823718423869104263333984641494340347905", 10)

// BCoeffFQ2 of the G2 curve.
var BCoeffFQ2 = NewFQ2(NewFQ(BCoeff), NewFQ(BCoeff))

// G2AffineOne represents the point at 1 on G2.
var G2AffineOne = &G2Affine{
	x: NewFQ2(
		NewFQ(g2GeneratorXC0),
		NewFQ(g2GeneratorXC1),
	),
	y: NewFQ2(
		NewFQ(g2GeneratorYC0),
		NewFQ(g2GeneratorYC1),
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

// Neg negates the point.
func (g G2Affine) Neg() *G2Affine {
	if !g.IsZero() {
		return NewG2Affine(g.x, g.y.Neg())
	}
	return g.Copy()
}

// ToProjective converts an affine point to a projective one.
func (g G2Affine) ToProjective() *G2Projective {
	if g.IsZero() {
		return G2ProjectiveZero
	}
	return NewG2Projective(g.x, g.y, FQ2One)
}

// Mul performs a EC multiply operation on the point.
func (g G2Affine) Mul(b *big.Int) *G2Projective {
	bs := b.Bytes()
	res := G2ProjectiveZero.Copy()
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

// IsOnCurve checks if a point is on the G2 curve.
func (g G2Affine) IsOnCurve() bool {
	if g.infinity {
		return true
	}
	y2 := g.y.Square()
	x3b := g.x.Square().Mul(g.x).Add(BCoeffFQ2)

	return y2.Equals(x3b)
}

// G2 cofactor = (x^8 - 4 x^7 + 5 x^6) - (4 x^4 + 6 x^3 - 4 x^2 - 4 x + 13) // 9
var g2Cofactor, _ = new(big.Int).SetString("5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5", 16)

// ScaleByCofactor scales the G2Affine point by the cofactor.
func (g G2Affine) ScaleByCofactor() *G2Projective {
	return g.Mul(g2Cofactor)
}

// Equals checks if two affine points are equal.
func (g G2Affine) Equals(other *G2Affine) bool {
	return (g.infinity == other.infinity) || (g.x.Equals(other.x) && g.y.Equals(other.y))
}

// GetG2PointFromX attempts to reconstruct an affine point given
// an x-coordinate. The point is not guaranteed to be in the subgroup.
// If and only if `greatest` is set will the lexicographically
// largest y-coordinate be selected.
func GetG2PointFromX(x *FQ2, greatest bool) *G2Affine {
	x3b := x.Square().Mul(x).Add(BCoeffFQ2)

	y := x3b.Sqrt()

	if y == nil {
		return nil
	}

	negY := y.Neg()

	yVal := negY
	if (y.Cmp(negY) < 0) != greatest {
		yVal = y
	}
	return NewG2Affine(x, yVal)
}

// DecompressG2 decompresses a G2 point from a big int and checks
// if it is in the correct subgroup.
func DecompressG2(b *big.Int) (*G2Affine, error) {
	affine, err := DecompressG2Unchecked(b)
	if err != nil {
		return nil, err
	}

	if !affine.IsInCorrectSubgroupAssumingOnCurve() {
		return nil, errors.New("point is not in correct subgroup")
	}
	return affine, nil
}

// DecompressG2Unchecked decompresses a G2 point from a big int.
func DecompressG2Unchecked(b *big.Int) (*G2Affine, error) {
	copy := b.Bytes()

	if copy[0]&(1<<7) == 0 {
		return nil, errors.New("unexpected compression mode")
	}

	if copy[0]&(1<<6) != 0 {
		copy[0] &= 0x3f

		for _, b := range copy {
			if b != 0 {
				return nil, errors.New("unexpected information in infinity point on G2")
			}
		}
		return G2AffineZero.Copy(), nil
	}
	greatest := copy[0]&(1<<5) != 0

	copy[0] &= 0x1f

	xC0 := NewFQ(new(big.Int).SetBytes(copy[:48]))

	xC1 := NewFQ(new(big.Int).SetBytes(copy[48:]))

	x := NewFQ2(xC0, xC1)

	return GetG2PointFromX(x, greatest), nil
}

// CompressG2 compresses a G2 point into an int.
func CompressG2(affine *G2Affine) *big.Int {
	res := [96]byte{}
	if affine.IsZero() {
		res[0] |= 1 << 6
	} else {
		out0 := new(big.Int).Set(affine.x.c0.n).Bytes()
		out1 := new(big.Int).Set(affine.x.c0.n).Bytes()
		copy(res[:48], out0)
		copy(res[48:], out1)

		negY := affine.y.Neg()

		if affine.y.Cmp(negY) > 0 {
			res[0] |= 1 << 5
		}
	}

	res[0] |= 1 << 7
	return new(big.Int).SetBytes(res[:])
}

// IsInCorrectSubgroupAssumingOnCurve checks if the point multiplied by the
// field characteristic equals zero.
func (g G2Affine) IsInCorrectSubgroupAssumingOnCurve() bool {
	return g.Mul(frChar).IsZero()
}

// G2Projective is a projective point on the G2 curve.
type G2Projective struct {
	x *FQ2
	y *FQ2
	z *FQ2
}

// NewG2Projective creates a new G2Projective point.
func NewG2Projective(x *FQ2, y *FQ2, z *FQ2) *G2Projective {
	return &G2Projective{x, y, z}
}

// G2ProjectiveZero is the point at infinity where Z = 0.
var G2ProjectiveZero = &G2Projective{FQ2Zero, FQ2One, FQ2Zero}

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
func (g G2Projective) Equal(other *G2Projective) bool {
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

// ToAffine converts a G2Projective point to affine form.
func (g G2Projective) ToAffine() *G2Affine {
	if g.IsZero() {
		return G2AffineZero
	} else if g.z.IsZero() {
		return NewG2Affine(g.x, g.y)
	}

	// nonzero so must have an inverse
	zInv := g.z.Inverse()
	zInvSquared := zInv.Square()

	return NewG2Affine(g.x.Mul(zInvSquared), g.y.Mul(zInvSquared).Mul(zInv))
}

// Double performs EC doubling on the point.
func (g G2Projective) Double() *G2Projective {
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

	return NewG2Projective(newX, newY, newZ)
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

	return NewG2Projective(newX, newY, newZ)
}

// Mul performs a EC multiply operation on the point.
func (g G2Projective) Mul(b *big.Int) *G2Projective {
	bs := b.Bytes()
	res := G2ProjectiveZero.Copy()
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

var blsX, _ = new(big.Int).SetString("d201000000010000", 16)

const blsIsNegative = true

// G2Prepared is a prepared G2 point multiplication by blsX.
type G2Prepared struct {
	coeffs   [][3]*FQ2
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

	doublingStep := func(r *G2Projective) (*FQ2, *FQ2, *FQ2) {
		tmp0 := r.x.Square()
		tmp1 := r.y.Square()
		tmp2 := tmp1.Square()
		tmp3 := tmp1.Add(r.x).Square().Sub(tmp0).Sub(tmp2).Double()
		tmp4 := tmp0.Double().Add(tmp0)
		tmp6 := r.x.Add(tmp4)
		tmp5 := tmp4.Square()
		zSquared := r.z.Square()
		r.x = tmp5.Sub(tmp3).Sub(tmp3)
		r.z = r.z.Add(r.y).Square().Sub(tmp1).Sub(zSquared)
		r.y = tmp3.Sub(r.x).Mul(tmp4)
		tmp2 = tmp2.Double().Double().Double()
		r.y = r.y.Sub(tmp2)
		tmp3 = tmp4.Mul(zSquared).Double().Neg()
		tmp6 = tmp6.Square().Sub(tmp0).Sub(tmp5)
		tmp1 = tmp1.Double().Double()
		tmp6 = tmp6.Sub(tmp1)
		tmp0 = r.z.Mul(zSquared).Double()
		return tmp0, tmp3, tmp6
	}

	additionStep := func(r *G2Projective, q *G2Affine) (*FQ2, *FQ2, *FQ2) {
		zSquared := r.z.Square()
		ySquared := q.y.Square()
		t0 := zSquared.Mul(q.x)
		t1 := q.y.Add(r.z).Square().Sub(ySquared).Sub(zSquared).Mul(zSquared)
		t2 := t0.Sub(r.x)
		t3 := t2.Square()
		t4 := t3.Double().Double()
		t5 := t4.Mul(t2)
		t6 := t1.Sub(r.y).Sub(r.y)
		t9 := t6.Mul(q.x)
		t7 := t4.Mul(r.x)
		r.x = t6.Square().Sub(t5).Sub(t7).Sub(t7)
		r.z = r.z.Add(t2).Square().Sub(zSquared).Sub(t3)
		t10 := q.y.Add(r.z)
		t8 := t7.Sub(r.x).Mul(t6)
		t0 = r.y.Mul(t5).Double()
		r.y = t8.Sub(t0)
		t10 = t10.Square().Sub(ySquared).Sub(r.z.Square())
		t9 = t9.Double().Sub(t10)
		t10 = r.z.Double()
		t6 = t6.Neg()
		t1 = t6.Double()

		return t10, t1, t9
	}

	coeffs := [][3]*FQ2{}
	r := q.ToProjective()

	foundOne := false
	blsXLsh1 := new(big.Int).Rsh(blsX, 1)
	bs := blsXLsh1.Bytes()
	for i := uint(0); i < uint(len(bs)*8); i++ {
		segment := i / 8
		bit := 7 - i%8
		set := bs[segment]&(1<<bit) > 0
		if !foundOne {
			foundOne = set
			continue
		}

		o0, o1, o2 := doublingStep(r)

		coeffs = append(coeffs, [3]*FQ2{o0, o1, o2})

		if set {
			o0, o1, o2 := additionStep(r, q)
			coeffs = append(coeffs, [3]*FQ2{o0, o1, o2})
		}
	}
	o0, o1, o2 := doublingStep(r)

	coeffs = append(coeffs, [3]*FQ2{o0, o1, o2})

	return &G2Prepared{coeffs, false}
}
