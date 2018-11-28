package bls

import (
	"fmt"
	"io"
	"math/big"
)

var oneLsh384MinusOne = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 384), bigOne)

// FQ2 represents an element of Fq2, represented by c0 + c1 * u.
type FQ2 struct {
	c0 *FQ
	c1 *FQ
}

// NewFQ2 constructs a new FQ2 element given two FQ elements.
func NewFQ2(c0 *FQ, c1 *FQ) *FQ2 {
	return &FQ2{
		c0: c0,
		c1: c1,
	}
}

func (f FQ2) String() string {
	return fmt.Sprintf("Fq2(%s + %s * u)", f.c0, f.c1)
}

// Cmp compares two FQ2 elements.
func (f FQ2) Cmp(other *FQ2) int {
	cOut := f.c1.Cmp(other.c1)
	if cOut != 0 {
		return cOut
	}
	return f.c0.Cmp(other.c0)
}

// MultiplyByNonresidue multiplies this element by the cubic and quadratic
// nonresidue 1 + u.
func (f FQ2) MultiplyByNonresidue() *FQ2 {
	return &FQ2{
		c0: f.c0.Sub(f.c1),
		c1: f.c1.Add(f.c0),
	}
}

// MultiplyByNonresidueAssign multiplies this element by the cubic and quadratic
// nonresidue 1 + u.
func (f *FQ2) MultiplyByNonresidueAssign() {
	oldC0 := f.c0.Copy()
	f.c0.SubAssign(f.c1)
	f.c1.AddAssign(oldC0)
}

// Norm gets the norm of Fq2 as extension field in i over Fq.
func (f *FQ2) Norm() *FQ {
	t0 := f.c0.Copy()
	t1 := f.c1.Copy()
	t0.SquareAssign()
	t1.SquareAssign()
	t1.AddAssign(t0)
	return t1
}

// FQ2Zero gets the zero element of the field.
var FQ2Zero = &FQ2{
	c0: FQZero,
	c1: FQZero,
}

// FQ2One gets the one-element of the field.
var FQ2One = &FQ2{
	c0: FQOne,
	c1: FQZero,
}

// IsZero checks if the field element is zero.
func (f FQ2) IsZero() bool {
	return f.c0.IsZero() && f.c1.IsZero()
}

// Square squares the FQ2 element.
func (f FQ2) Square() *FQ2 {
	ab := f.c0.Mul(f.c1)
	c0c1 := f.c0.Copy()
	c0c1.AddAssign(f.c1)
	c0 := f.c1.Neg()
	c0.AddAssign(f.c0)
	c0.MulAssign(c0c1)
	c0.SubAssign(ab)
	c0.AddAssign(ab)
	ab.AddAssign(ab)
	return &FQ2{
		c0: c0,
		c1: ab,
	}
}

// SquareAssign squares the FQ2 element.
func (f *FQ2) SquareAssign() {
	ab := f.c0.Mul(f.c1)
	c0 := f.c1.Neg()
	c0.AddAssign(f.c0)
	f.c0.AddAssign(f.c1)
	c0.MulAssign(f.c0)
	c0.SubAssign(ab)
	c0.AddAssign(ab)
	ab.AddAssign(ab)
	f.c0 = c0
	f.c1 = ab
}

// Double doubles an FQ2 element.
func (f FQ2) Double() *FQ2 {
	return &FQ2{
		c0: f.c0.Double(),
		c1: f.c1.Double(),
	}
}

// DoubleAssign doubles an FQ2 element.
func (f *FQ2) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
}

// Neg negates a FQ2 element.
func (f FQ2) Neg() *FQ2 {
	return &FQ2{
		c0: f.c0.Neg(),
		c1: f.c1.Neg(),
	}
}

// NegAssign negates a FQ2 element.
func (f *FQ2) NegAssign() {
	f.c0.NegAssign()
	f.c1.NegAssign()
}

// Add adds two FQ2 elements together.
func (f FQ2) Add(other *FQ2) *FQ2 {
	return &FQ2{
		c0: f.c0.Add(other.c0),
		c1: f.c1.Add(other.c1),
	}
}

// AddAssign adds two FQ2 elements together.
func (f *FQ2) AddAssign(other *FQ2) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
}

// Sub subtracts one field element from another.
func (f FQ2) Sub(other *FQ2) *FQ2 {
	return &FQ2{
		c0: f.c0.Sub(other.c0),
		c1: f.c1.Sub(other.c1),
	}
}

// SubAssign subtracts one field element from another.
func (f *FQ2) SubAssign(other *FQ2) {
	f.c0.SubAssign(other.c0)
	f.c1.SubAssign(other.c1)
}

// Mul multiplies two FQ2 elements together.
func (f FQ2) Mul(other *FQ2) *FQ2 {
	aa := f.c0.Copy()
	aa.MulAssign(other.c0)
	bb := f.c1.Copy()
	bb.MulAssign(other.c1)
	o := other.c0.Add(other.c1)
	c1 := f.c1.Add(f.c0)
	c1.MulAssign(o)
	c1.SubAssign(aa)
	c1.SubAssign(bb)

	aa.SubAssign(bb)

	return &FQ2{
		c1: c1,
		c0: aa,
	}
}

// MulAssign multiplies two FQ2 elements together.
func (f *FQ2) MulAssign(other *FQ2) {
	aa := f.c0.Copy()
	aa.MulAssign(other.c0)
	bb := f.c1.Copy()
	bb.MulAssign(other.c1)
	o := other.c0.Copy()
	o.AddAssign(other.c1)
	f.c1.AddAssign(f.c0)
	f.c1.MulAssign(o)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)

	f.c0 = aa
	f.c0.SubAssign(bb)
}

// Inverse finds the inverse of the field element.
func (f FQ2) Inverse() *FQ2 {
	inv := f.c0.Square().Add(f.c1.Square()).Inverse()
	if inv == nil {
		return nil
	}
	return &FQ2{
		c0: f.c0.Mul(inv),
		c1: f.c1.Mul(inv).Neg(),
	}
}

// InverseAssign finds the inverse of the field element.
func (f *FQ2) InverseAssign() bool {
	inv := f.c0.Square()
	inv.AddAssign(f.c1.Square())
	inv.Inverse()
	if inv == nil {
		return false
	}
	f.c0.MulAssign(inv)
	f.c1.MulAssign(inv)
	f.c1.NegAssign()
	return true
}

var frobeniusCoeffFQ2c11 = NewFQ(bigOne).Neg().Exp(qMinus1Over2)

var frobeniusCoeffFQ2c1 = [2]*FQ{
	FQOne,
	frobeniusCoeffFQ2c11,
}

// FrobeniusMap multiplies the element by the Frobenius automorphism
// coefficient.
func (f FQ2) FrobeniusMap(power uint8) *FQ2 {
	return NewFQ2(f.c0, f.c1.Mul(frobeniusCoeffFQ2c1[power%2]))
}

// FrobeniusMapAssign multiplies the element by the Frobenius automorphism
// coefficient.
func (f FQ2) FrobeniusMapAssign(power uint8) {
	f.c1.MulAssign(frobeniusCoeffFQ2c1[power%2])
}

// Legendre gets the legendre symbol of the FQ2 element.
func (f FQ2) Legendre() LegendreSymbol {
	return f.Norm().Legendre()
}

var qMinus3Over4, _ = new(big.Int).SetString("1000602388805416848354447456433976039139220704984751971333014534031007912622709466110671907282253916009473568139946", 10)

// Exp raises the element ot a specific power.
func (f FQ2) Exp(n *big.Int) *FQ2 {
	if n.Cmp(bigZero) == 0 {
		return FQ2One.Copy()
	} else if n.Cmp(bigOne) == 0 {
		return f.Copy()
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigZero) == 0 {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo))
	} else {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo)).Mul(&f)
	}
}

// ExpAssign raises the element ot a specific power.
func (f *FQ2) ExpAssign(n *big.Int) {
	if n.Cmp(bigZero) == 0 {
		f = FQ2One.Copy()
	} else if n.Cmp(bigOne) == 0 {
		return
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigZero) == 0 {
		f.SquareAssign()
		n.Div(n, bigTwo)
		f.ExpAssign(n)
	} else {
		f0 := f.Copy()
		f.MulAssign(f)
		n.Div(n, bigTwo)
		f.ExpAssign(n)
		f.Mul(f0)
	}
}

// -(2**384 mod q) mod q
var negativeOne, _ = new(big.Int).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", 10)

// Equals checks if this FQ2 equals another one.
func (f FQ2) Equals(other *FQ2) bool {
	return f.Cmp(other) == 0
}

// Sqrt finds the sqrt of a field element.
func (f FQ2) Sqrt() *FQ2 {
	// Algorithm 9, https://eprint.iacr.org/2012/685.pdf
	if f.IsZero() {
		return FQ2Zero
	}
	a1 := f.Exp(qMinus3Over4)
	alpha := a1.Square().Mul(&f)
	a0 := alpha.FrobeniusMap(1).Mul(alpha)

	neg1 := &FQ2{
		c0: NewFQ(negativeOne),
		c1: FQZero,
	}

	if a0.Equals(neg1) {
		return nil
	}
	a1 = a1.Mul(&f)

	if alpha.Equals(neg1) {
		return a1.Mul(&FQ2{
			c0: FQZero,
			c1: FQOne,
		})
	}
	return alpha.Add(FQ2One).Exp(qMinus1Over2).Mul(a1)
}

// Copy returns a copy of the field element.
func (f *FQ2) Copy() *FQ2 {
	return &FQ2{
		c0: f.c0.Copy(),
		c1: f.c1.Copy(),
	}
}

// RandFQ2 generates a random FQ2 element.
func RandFQ2(reader io.Reader) (*FQ2, error) {
	i0, err := RandFQ(reader)
	if err != nil {
		return nil, err
	}
	i1, err := RandFQ(reader)
	if err != nil {
		return nil, err
	}
	return NewFQ2(
		i0,
		i1,
	), nil
}
