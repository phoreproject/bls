package bls

import (
	"fmt"
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
func RandFQ12() (*FQ12, error) {
	a, err := RandFQ6()
	if err != nil {
		return nil, err
	}
	b, err := RandFQ6()
	if err != nil {
		return nil, err
	}
	return NewFQ12(a, b), nil
}

var frobeniusCoeffFQ12c110, _ = new(big.Int).SetString("1376889598125376727959055341295674356654925039980005395128828212993454708588385020118431646457834669954221389501541", 10)
var frobeniusCoeffFQ12c111, _ = new(big.Int).SetString("2625519957096290665458734484440229799901957779959002490203229923130576941902452844324255982671180994083672883058246", 10)
var frobeniusCoeffFQ12c130, _ = new(big.Int).SetString("1821461487266245992767491788684378228087062278322214693001359809350238716280406307949636812899085786271837335624401", 10)
var frobeniusCoeffFQ12c131, _ = new(big.Int).SetString("2180948067955421400650298037051525928469820541616793192330698326773792934210431556493050816229929877766056936935386", 10)
var frobeniusCoeffFQ12c150, _ = new(big.Int).SetString("444571889140869264808436447388703871432137238342209297872531596356784007692021287831205166441251116317615946122860", 10)
var frobeniusCoeffFQ12c151, _ = new(big.Int).SetString("3557837666080798128609353378347200285124745581596798587459526539767247642798816576611482462687764547720278326436927", 10)

var frobeniusCoeffFQ12c1 = [12]*FQ2{
	{
		c0: NewFQ(fq6c10),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c110),
		c1: NewFQ(frobeniusCoeffFQ12c111),
	},
	{
		c0: NewFQ(fq6c25),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c130),
		c1: NewFQ(frobeniusCoeffFQ12c131),
	},
	{
		c0: NewFQ(fq6c12),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c150),
		c1: NewFQ(frobeniusCoeffFQ12c151),
	},
	{
		c0: NewFQ(fq6c20),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c111),
		c1: NewFQ(frobeniusCoeffFQ12c110),
	},
	{
		c0: NewFQ(fq6c11),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c131),
		c1: NewFQ(frobeniusCoeffFQ12c130),
	},
	{
		c0: NewFQ(fq6c21),
		c1: FQZero,
	},
	{
		c0: NewFQ(frobeniusCoeffFQ12c151),
		c1: NewFQ(frobeniusCoeffFQ12c150),
	},
}

// FrobeniusMap calculates the frobenius map of an FQ12 element.
func (f FQ12) FrobeniusMap(power uint8) *FQ12 {
	n := f.c1.FrobeniusMap(power)
	return NewFQ12(
		f.c0.FrobeniusMap(power),
		NewFQ6(
			n.c0.Mul(frobeniusCoeffFQ12c1[power%12]),
			n.c1.Mul(frobeniusCoeffFQ12c1[power%12]),
			n.c2.Mul(frobeniusCoeffFQ12c1[power%12]),
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
