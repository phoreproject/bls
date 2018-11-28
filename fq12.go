package bls

import (
	"fmt"
	"io"
	"math/big"
)

// FQ12 is an element of Fq12, represented by c0 + c1 * w.
type FQ12 struct {
	c0 *FQ6
	c1 *FQ6
}

// NewFQ12 creates a new FQ12 element from two FQ6 elements.
func NewFQ12(c0 *FQ6, c1 *FQ6) *FQ12 {
	return &FQ12{
		c0: c0,
		c1: c1,
	}
}

func (f *FQ12) String() string {
	return fmt.Sprintf("Fq12(%s + %s * w)", f.c0, f.c1)
}

// Conjugate returns the conjugate of the FQ12 element.
func (f *FQ12) Conjugate() *FQ12 {
	return NewFQ12(f.c0, f.c1.Neg())
}

// MulBy014 multiplies FQ12 element by 3 FQ2 elements.
func (f *FQ12) MulBy014(c0 *FQ2, c1 *FQ2, c4 *FQ2) *FQ12 {
	aa := f.c0.MulBy01(c0, c1)
	bb := f.c1.MulBy1(c4)
	o := c1.Add(c4)
	return NewFQ12(
		bb.MulByNonresidue().Add(aa),
		f.c1.Add(f.c0).MulBy01(c0, o).Sub(aa).Sub(bb),
	)
}

// FQ12Zero is the zero element of FQ12.
var FQ12Zero = NewFQ12(FQ6Zero, FQ6Zero)

// FQ12One is the one element of FQ12.
var FQ12One = NewFQ12(FQ6One, FQ6Zero)

// Equals checks if two FQ12 elements are equal.
func (f FQ12) Equals(other *FQ12) bool {
	return f.c0.Equals(other.c0) && f.c1.Equals(other.c1)
}

// Double doubles each coefficient in an FQ12 element.
func (f FQ12) Double() *FQ12 {
	return NewFQ12(f.c0.Double(), f.c1.Double())
}

// Neg negates each coefficient in an FQ12 element.
func (f FQ12) Neg() *FQ12 {
	return NewFQ12(f.c0.Neg(), f.c1.Neg())
}

// Add adds two FQ12 elements together.
func (f FQ12) Add(other *FQ12) *FQ12 {
	return NewFQ12(
		f.c0.Add(other.c0),
		f.c1.Add(other.c1),
	)
}

// Sub subtracts one FQ12 element from another.
func (f FQ12) Sub(other *FQ12) *FQ12 {
	return NewFQ12(
		f.c0.Sub(other.c0),
		f.c1.Sub(other.c1),
	)
}

// RandFQ12 generates a random FQ12 element.
func RandFQ12(reader io.Reader) (*FQ12, error) {
	a, err := RandFQ6(reader)
	if err != nil {
		return nil, err
	}
	b, err := RandFQ6(reader)
	if err != nil {
		return nil, err
	}
	return NewFQ12(a, b), nil
}

// Copy returns a copy of the FQ12 element.
func (f FQ12) Copy() *FQ12 {
	return NewFQ12(f.c0.Copy(), f.c1.Copy())
}

// Exp raises the element ot a specific power.
func (f FQ12) Exp(n *big.Int) *FQ12 {
	if n.Cmp(bigZero) == 0 {
		return FQ12One.Copy()
	} else if n.Cmp(bigOne) == 0 {
		return f.Copy()
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigZero) == 0 {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo))
	} else {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo)).Mul(&f)
	}
}

var bigSix = big.NewInt(6)

func getFrobExpMinus1Over6(power int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Sub(new(big.Int).Exp(QFieldModulus, big.NewInt(power), nil), bigOne), bigSix)
}

var frobeniusCoeffFQ12c1 = [12]*FQ2{
	FQ2One,
	fq2nqr.Exp(getFrobExpMinus1Over6(1)),
	fq2nqr.Exp(getFrobExpMinus1Over6(2)),
	fq2nqr.Exp(getFrobExpMinus1Over6(3)),
	fq2nqr.Exp(getFrobExpMinus1Over6(4)),
	fq2nqr.Exp(getFrobExpMinus1Over6(5)),
	fq2nqr.Exp(getFrobExpMinus1Over6(6)),
	fq2nqr.Exp(getFrobExpMinus1Over6(7)),
	fq2nqr.Exp(getFrobExpMinus1Over6(8)),
	fq2nqr.Exp(getFrobExpMinus1Over6(9)),
	fq2nqr.Exp(getFrobExpMinus1Over6(10)),
	fq2nqr.Exp(getFrobExpMinus1Over6(11)),
}

// FrobeniusMap calculates the frobenius map of an FQ12 element.
func (f FQ12) FrobeniusMap(power uint8) *FQ12 {
	c1 := f.c1.FrobeniusMap(power)
	return NewFQ12(
		f.c0.FrobeniusMap(power),
		NewFQ6(
			c1.c0.Mul(frobeniusCoeffFQ12c1[power%12]),
			c1.c1.Mul(frobeniusCoeffFQ12c1[power%12]),
			c1.c2.Mul(frobeniusCoeffFQ12c1[power%12]),
		),
	)
}

// Square calculates the square of the FQ12 element.
func (f FQ12) Square() *FQ12 {
	ab := f.c0.Mul(f.c1)
	c0c1 := f.c0.Add(f.c1)
	c0 := f.c1.MulByNonresidue().Add(f.c0).Mul(c0c1).Sub(ab)
	c1 := ab.Add(ab)
	ab = ab.MulByNonresidue()
	return NewFQ12(
		c0.Sub(ab),
		c1,
	)
}

// Mul multiplies two FQ12 elements together.
func (f FQ12) Mul(other *FQ12) *FQ12 {
	aa := f.c0.Mul(other.c0)
	bb := f.c1.Mul(other.c1)
	o := other.c0.Add(other.c1)
	return NewFQ12(
		bb.MulByNonresidue().Add(aa),
		f.c1.Add(f.c0).Mul(o).Sub(aa).Sub(bb),
	)
}

// MulAssign multiplies two FQ12 elements together.
func (f FQ12) MulAssign(other *FQ12) {
	aa := f.c0.Mul(other.c0)
	bb := f.c1.Mul(other.c1)
	o := other.c0.Add(other.c1)

	f.c0 = bb.Copy()
	f.c0.MulByNonresidueAssign()
	f.c0.AddAssign(aa)

	f.c1.AddAssign(f.c0)
	f.c1.MulAssign(o)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
}

// Inverse finds the inverse of an FQ12
func (f FQ12) Inverse() *FQ12 {
	c1s := f.c1.Square().MulByNonresidue()
	c0s := f.c0.Square().Sub(c1s)

	i := c0s.Inverse()
	if i == nil {
		return nil
	}

	return NewFQ12(i.Mul(f.c0), i.Mul(f.c1).Neg())
}
