package bls

import (
	"fmt"
	"hash"
	"io"
	"math/big"
)

var oneLsh384MinusOne, _ = FQReprFromBigInt(new(big.Int).Sub(new(big.Int).Lsh(bigOne.ToBig(), 384), bigOne.ToBig()))

// FQ2 represents an element of Fq2, represented by c0 + c1 * u.
type FQ2 struct {
	c0 FQ
	c1 FQ
}

// NewFQ2 constructs a new FQ2 element given two FQ elements.
func NewFQ2(c0 FQ, c1 FQ) FQ2 {
	return FQ2{
		c0: c0,
		c1: c1,
	}
}

func (f FQ2) String() string {
	return fmt.Sprintf("Fq2(%s + %s * u)", f.c0, f.c1)
}

// Cmp compares two FQ2 elements.
func (f FQ2) Cmp(other FQ2) int {
	cOut := f.c1.Cmp(other.c1)
	if cOut != 0 {
		return cOut
	}
	return f.c0.Cmp(other.c0)
}

// MultiplyByNonresidueAssign multiplies this element by the cubic and quadratic
// nonresidue 1 + u.
func (f *FQ2) MultiplyByNonresidueAssign() {
	oldC0 := f.c0.Copy()
	f.c0.SubAssign(f.c1)
	f.c1.AddAssign(oldC0)
}

// Norm gets the norm of Fq2 as extension field in i over Fq.
func (f *FQ2) Norm() FQ {
	t0 := f.c0.Copy()
	t1 := f.c1.Copy()
	t0.SquareAssign()
	t1.SquareAssign()
	t1.AddAssign(t0)
	return t1
}

// FQ2Zero gets the zero element of the field.
var FQ2Zero = FQ2{
	c0: FQZero,
	c1: FQZero,
}

// FQ2One gets the one-element of the field.
var FQ2One = FQ2{
	c0: FQOne,
	c1: FQZero,
}

// IsZero checks if the field element is zero.
func (f FQ2) IsZero() bool {
	return f.c0.IsZero() && f.c1.IsZero()
}

// SquareAssign squares the FQ2 element.
func (f *FQ2) SquareAssign() {
	ab := f.c0.Copy()
	ab.MulAssign(f.c1)
	c0c1 := f.c0.Copy()
	c0c1.AddAssign(f.c1)
	c0 := f.c1.Copy()
	c0.NegAssign()
	c0.AddAssign(f.c0)
	c0.MulAssign(c0c1)
	c0.SubAssign(ab)
	c0.AddAssign(ab)
	ab.AddAssign(ab)
	f.c0 = c0
	f.c1 = ab
}

// DoubleAssign doubles an FQ2 element.
func (f *FQ2) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
}

// NegAssign negates a FQ2 element.
func (f *FQ2) NegAssign() {
	f.c0.NegAssign()
	f.c1.NegAssign()
}

// AddAssign adds two FQ2 elements together.
func (f *FQ2) AddAssign(other FQ2) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
}

// SubAssign subtracts one field element from another.
func (f *FQ2) SubAssign(other FQ2) {
	f.c0.SubAssign(other.c0)
	f.c1.SubAssign(other.c1)
}

// MulAssign multiplies two FQ2 elements together.
func (f *FQ2) MulAssign(other FQ2) {
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

// InverseAssign finds the inverse of the field element.
func (f *FQ2) InverseAssign() bool {
	t1 := f.c1.Copy()
	t1.SquareAssign()
	t0 := f.c0.Copy()
	t0.SquareAssign()
	t0.AddAssign(t1)
	t, success := t0.Inverse()
	if !success {
		return false
	}
	f.c0.MulAssign(t)
	f.c1.MulAssign(t)
	f.c1.NegAssign()
	return true
}

var frobeniusCoeffFQ2c1 = [2]FQ{
	FQOne,
	FQReprToFQRaw(FQRepr{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206}),
}

// FrobeniusMapAssign multiplies the element by the Frobenius automorphism
// coefficient.
func (f *FQ2) FrobeniusMapAssign(power uint8) {
	f.c1.MulAssign(frobeniusCoeffFQ2c1[power%2])
}

// Legendre gets the legendre symbol of the FQ2 element.
func (f FQ2) Legendre() LegendreSymbol {
	norm := f.Norm()

	return norm.Legendre()
}

var qMinus3Over4 = fqReprFromHexUnchecked("680447a8e5ff9a692c6e9ed90d2eb35d91dd2e13ce144afd9cc34a83dac3d8907aaffffac54ffffee7fbfffffffeaaa")

// Exp raises the element ot a specific power.
func (f FQ2) Exp(n FQRepr) FQ2 {
	iter := NewBitIterator(n[:])
	res := FQ2One.Copy()
	foundOne := false
	next, done := iter.Next()
	for !done {
		if foundOne {
			res.SquareAssign()
		} else {
			foundOne = next
		}
		if next {
			res.MulAssign(f)
		}
		next, done = iter.Next()
	}
	return res
}

// -(2**384 mod q) mod q
var negativeOne, _ = FQReprFromString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559786", 10)

// Equals checks if this FQ2 equals another one.
func (f FQ2) Equals(other FQ2) bool {
	return f.Cmp(other) == 0
}

// Sqrt finds the sqrt of a field element.
func (f FQ2) Sqrt() (FQ2, bool) {
	// Algorithm 9, https://eprint.iacr.org/2012/685.pdf
	if f.IsZero() {
		return FQ2Zero, true
	}
	a1 := f.Exp(qMinus3Over4)
	alpha := a1.Copy()
	alpha.SquareAssign()
	alpha.MulAssign(f)
	a0 := alpha.Copy()
	a0.FrobeniusMapAssign(1)
	a0.MulAssign(alpha)

	neg1 := FQ2{
		c0: negativeOneFQ,
		c1: FQZero,
	}

	if a0.Equals(neg1) {
		return FQ2{}, false
	}
	a1.MulAssign(f)

	if alpha.Equals(neg1) {
		a1.MulAssign(FQ2{
			c0: FQZero,
			c1: FQOne,
		})
		return a1, true
	}
	alpha.AddAssign(FQ2One)
	alpha = alpha.Exp(qMinus1Over2)
	alpha.MulAssign(a1)
	return alpha, true
}

// Copy returns a copy of the field element.
func (f *FQ2) Copy() FQ2 {
	return *f
}

// RandFQ2 generates a random FQ2 element.
func RandFQ2(reader io.Reader) (FQ2, error) {
	i0, err := RandFQ(reader)
	if err != nil {
		return FQ2{}, err
	}
	i1, err := RandFQ(reader)
	if err != nil {
		return FQ2{}, err
	}
	return NewFQ2(
		i0,
		i1,
	), nil
}

// Parity checks if the point is greater than the point negated.
func (f FQ2) Parity() bool {
	neg := f.Copy()
	neg.NegAssign()
	return f.Cmp(neg) > 0
}

// MulBits multiplies the number by a big number.
func (f FQ2) MulBits(b *big.Int) FQ2 {
	res := FQ2Zero
	for i := 0; i < b.BitLen(); i++ {
		res.DoubleAssign()
		if b.Bit(b.BitLen()-1-i) == 1 {
			res.AddAssign(f)
		}
	}
	return res
}

// DivAssign divides the FQ2 element by another FQ2 element.
func (f *FQ2) DivAssign(other FQ2) {
	other.InverseAssign()
	f.MulAssign(other)
}

// HashFQ2 calculates a new FQ2 value based on a hash.
func HashFQ2(hasher hash.Hash) FQ2 {
	digest := hasher.Sum(nil)
	newB := new(big.Int).SetBytes(digest)
	return FQ2One.MulBits(newB)
}
