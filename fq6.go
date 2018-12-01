package bls

import (
	"fmt"
	"io"
)

// FQ6 is an element of FQ6 represented by c0 + c1*v + v2*v**2
type FQ6 struct {
	c0 *FQ2
	c1 *FQ2
	c2 *FQ2
}

// NewFQ6 creates a new FQ6 element.
func NewFQ6(c0 *FQ2, c1 *FQ2, c2 *FQ2) *FQ6 {
	return &FQ6{
		c0: c0,
		c1: c1,
		c2: c2,
	}
}

func (f FQ6) String() string {
	return fmt.Sprintf("Fq6(%s + %s*v + %s*v^2)", f.c0, f.c1, f.c2)
}

// Copy creates a copy of the field element.
func (f FQ6) Copy() *FQ6 {
	return NewFQ6(f.c0.Copy(), f.c1.Copy(), f.c2.Copy())
}

// MulByNonresidueAssign multiplies by quadratic nonresidue v.
func (f FQ6) MulByNonresidueAssign() {
	f.c0, f.c1, f.c2 = f.c2, f.c0, f.c1
	f.c0.MultiplyByNonresidueAssign()
}

// MulBy1Assign multiplies the FQ6 by an FQ2.
func (f FQ6) MulBy1Assign(c1 *FQ2) {
	b := f.c1.Copy()
	b.MulAssign(c1)
	tmp := f.c1.Copy()
	tmp.AddAssign(f.c2)
	t1 := c1.Copy()
	t1.MulAssign(tmp)
	t1.SubAssign(b)
	t1.MultiplyByNonresidueAssign()
	tmp = f.c0.Copy()
	tmp.AddAssign(f.c1)
	t2 := c1.Copy()
	t2.MulAssign(tmp)
	t2.SubAssign(b)
	f.c0 = t1
	f.c1 = t2
	f.c2 = b
}

// MulBy01Assign multiplies by c0 and c1.
func (f *FQ6) MulBy01Assign(c0 *FQ2, c1 *FQ2) {
	a := f.c0.Copy()
	a.MulAssign(c0)
	b := f.c1.Copy()
	b.MulAssign(c1)

	tmp := f.c1.Copy()
	tmp.AddAssign(f.c2)
	t1 := c1.Copy()
	t1.MulAssign(tmp)
	t1.SubAssign(b)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(a)
	tmp = f.c0.Copy()
	tmp.AddAssign(f.c2)
	t3 := c0.Copy()
	t3.MulAssign(tmp)
	t3.SubAssign(a)
	t3.AddAssign(b)
	tmp = f.c0.Copy()
	tmp.AddAssign(f.c1)
	t2 := c0.Copy()
	t2.AddAssign(c1)
	t2.MulAssign(tmp)
	t2.SubAssign(a)
	t2.SubAssign(b)

	f.c0 = t1
	f.c1 = t2
	f.c2 = t3
}

// FQ6Zero represents the zero value of FQ6.
var FQ6Zero = NewFQ6(FQ2Zero, FQ2Zero, FQ2Zero)

// FQ6One represents the one value of FQ6.
var FQ6One = NewFQ6(FQ2One, FQ2Zero, FQ2Zero)

// Equals checks if two FQ6 elements are equal.
func (f FQ6) Equals(other *FQ6) bool {
	return f.c0.Equals(other.c0) && f.c1.Equals(other.c1) && f.c2.Equals(other.c2)
}

// IsZero checks if the FQ6 element is zero.
func (f FQ6) IsZero() bool {
	return f.Equals(FQ6Zero)
}

// DoubleAssign doubles the coefficients of the FQ6 element.
func (f FQ6) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
	f.c2.DoubleAssign()
}

// NegAssign negates the coefficients of the FQ6 element.
func (f FQ6) NegAssign() {
	f.c0.NegAssign()
	f.c1.NegAssign()
	f.c2.NegAssign()
}

// AddAssign the coefficients of the FQ6 element to another.
func (f FQ6) AddAssign(other *FQ6) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
	f.c2.AddAssign(other.c2)
}

// SubAssign subtracts the coefficients of the FQ6 element from another.
func (f FQ6) SubAssign(other *FQ6) {
	f.c0.SubAssign(other.c0)
	f.c1.SubAssign(other.c1)
	f.c2.SubAssign(other.c2)
}

func getFrobExpMinus1Over3(power int64) *FQRepr {
	f := FQReprToFQ(QFieldModulus.Copy())
	f = f.Exp(NewFQRepr(uint64(power)))
	f.SubAssign(FQOne)
	f.divAssign(bigThreeFQ)
	return f.n
}

func get2FrobExpMinus2Over3(power int64) *FQRepr {
	f := FQReprToFQ(QFieldModulus.Copy())
	f = f.Exp(NewFQRepr(uint64(power)))
	f.SubAssign(bigTwoFQ)
	f.divAssign(bigThreeFQ)
	return f.n
}

var bigThree = NewFQRepr(3)
var bigThreeFQ = FQReprToFQ(bigThree)

var fq2nqr = NewFQ2(
	FQOne,
	FQOne,
)

var frobeniusCoeffFQ6c1 = [6]*FQ2{
	// Fq2(u + 1)**(((q^0) - 1) / 3)
	FQ2One,
	// Fq2(u + 1)**(((q^1) - 1) / 3)
	fq2nqr.Exp(getFrobExpMinus1Over3(1)),
	// Fq2(u + 1)**(((q^2) - 1) / 3)
	fq2nqr.Exp(getFrobExpMinus1Over3(2)),
	// Fq2(u + 1)**(((q^3) - 1) / 3)
	fq2nqr.Exp(getFrobExpMinus1Over3(3)),
	// Fq2(u + 1)**(((q^4) - 1) / 3)
	fq2nqr.Exp(getFrobExpMinus1Over3(4)),
	// Fq2(u + 1)**(((q^5) - 1) / 3)
	fq2nqr.Exp(getFrobExpMinus1Over3(5)),
}

var frobeniusCoeffFQ6c2 = [6]*FQ2{
	// Fq2(u + 1)**(((2q^0) - 2) / 3)
	FQ2One,
	// Fq2(u + 1)**(((2q^1) - 2) / 3)
	fq2nqr.Exp(get2FrobExpMinus2Over3(1)),
	// Fq2(u + 1)**(((2q^2) - 2) / 3)
	fq2nqr.Exp(get2FrobExpMinus2Over3(2)),
	// Fq2(u + 1)**(((2q^3) - 2) / 3)
	fq2nqr.Exp(get2FrobExpMinus2Over3(3)),
	// Fq2(u + 1)**(((2q^4) - 2) / 3)
	fq2nqr.Exp(get2FrobExpMinus2Over3(4)),
	// Fq2(u + 1)**(((2q^5) - 2) / 3)
	fq2nqr.Exp(get2FrobExpMinus2Over3(5)),
}

// FrobeniusMapAssign runs the frobenius map algorithm with a certain power.
func (f FQ6) FrobeniusMapAssign(power uint8) {
	f.c0.FrobeniusMapAssign(power)
	f.c1.FrobeniusMapAssign(power)
	f.c1.MulAssign(frobeniusCoeffFQ6c1[power%6])
	f.c2.FrobeniusMapAssign(power)
	f.c2.MulAssign(frobeniusCoeffFQ6c2[power%6])
}

// SquareAssign squares the FQ6 element.
func (f FQ6) SquareAssign() {
	s0 := f.c0.Copy()
	s0.SquareAssign()
	ab := f.c0.Copy()
	ab.MulAssign(f.c1)
	s1 := ab.Copy()
	s1.DoubleAssign()
	s2 := f.c0.Copy()
	s2.SubAssign(f.c1)
	s2.AddAssign(f.c2)
	s2.SquareAssign()
	bc := f.c1.Copy()
	bc.MulAssign(f.c2)
	s3 := bc.Copy()
	s3.DoubleAssign()
	s4 := f.c2.Copy()
	s4.SquareAssign()

	f.c0 = s3.Copy()
	f.c0.MultiplyByNonresidueAssign()
	f.c0.AddAssign(s0)

	f.c1 = s4.Copy()
	f.c1.MultiplyByNonresidueAssign()
	f.c1.AddAssign(s1)

	f.c2 = s1.Copy()
	f.c2.AddAssign(s2)
	f.c2.AddAssign(s3)
	f.c2.SubAssign(s0)
	f.c2.SubAssign(s4)
}

// MulAssign multiplies two FQ6 elements together.
func (f *FQ6) MulAssign(other *FQ6) {
	aa := f.c0.Copy()
	aa.MulAssign(other.c0)
	bb := f.c1.Copy()
	bb.MulAssign(other.c1)
	cc := f.c2.Copy()
	cc.MulAssign(other.c2)

	tmp := f.c1.Copy()
	tmp.AddAssign(f.c2)
	t1 := other.c1.Copy()
	t1.AddAssign(other.c2)
	t1.MulAssign(tmp)
	t1.SubAssign(bb)
	t1.SubAssign(cc)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(aa)

	tmp = f.c0.Copy()
	tmp.AddAssign(f.c2)
	f.c2 = other.c0.Copy()
	f.c2.AddAssign(other.c2)
	f.c2.MulAssign(tmp)
	f.c2.SubAssign(aa)
	f.c2.AddAssign(bb)
	f.c2.SubAssign(cc)
	tmp = f.c0.Copy()
	tmp.AddAssign(f.c1)
	f.c1 = other.c0.Copy()
	f.c1.AddAssign(other.c1)
	f.c1.MulAssign(tmp)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
	cc.MultiplyByNonresidueAssign()
	f.c1.AddAssign(cc)

	f.c0 = t1
}

// InverseAssign finds the inverse of the FQ6 element.
func (f FQ6) InverseAssign() bool {
	c0 := f.c2.Copy()
	c0.MultiplyByNonresidueAssign()
	c0.MulAssign(f.c1)
	c0.NegAssign()
	c0s := f.c0.Copy()
	c0s.SquareAssign()
	c0.AddAssign(c0s)
	c1 := f.c2.Copy()
	c1.SquareAssign()
	c1.MultiplyByNonresidueAssign()
	c0c1 := f.c0.Copy()
	c0c1.MulAssign(f.c1)
	c0c2 := f.c0.Copy()
	c0c2.MulAssign(f.c2)
	c1c2 := f.c1.Copy()
	c1c2.MulAssign(f.c2)
	c1.SubAssign(c0c1)
	c2 := f.c1.Copy()
	c2.SquareAssign()
	c2.SubAssign(c0c2)

	tmp := f.c2.Copy()
	tmp.MulAssign(c1)
	tmp.AddAssign(c1c2)
	tmp.MultiplyByNonresidueAssign()
	c0c0 := f.c0.Copy()
	c0c0.SquareAssign()
	tmp.AddAssign(c0c0)
	tmpInverse := tmp.Copy()
	if !tmpInverse.InverseAssign() {
		return false
	}
	f.c0 = tmpInverse.Copy()
	f.c0.MulAssign(c0)
	f.c1 = tmpInverse.Copy()
	f.c1.MulAssign(c1)
	f.c2 = tmpInverse.Copy()
	f.c2.MulAssign(c2)
	return true
}

// RandFQ6 generates a random FQ6 element.
func RandFQ6(reader io.Reader) (*FQ6, error) {
	c0, err := RandFQ2(reader)
	if err != nil {
		return nil, err
	}
	c1, err := RandFQ2(reader)
	if err != nil {
		return nil, err
	}
	c2, err := RandFQ2(reader)
	if err != nil {
		return nil, err
	}
	return NewFQ6(c0, c1, c2), nil
}
