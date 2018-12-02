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

// ConjugateAssign returns the conjugate of the FQ12 element.
func (f *FQ12) ConjugateAssign() {
	f.c1.NegAssign()
}

// MulBy014Assign multiplies FQ12 element by 3 FQ2 elements.
func (f *FQ12) MulBy014Assign(c0 *FQ2, c1 *FQ2, c4 *FQ2) {
	aa := f.c0.Copy()
	aa.MulBy01Assign(c0, c1)
	bb := f.c1.Copy()
	bb.MulBy1Assign(c4)
	o := c1.Copy()
	c1.AddAssign(c4)
	f.c1.AddAssign(f.c0)
	f.c1.MulBy01Assign(c0, o)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
	f.c0 = bb.Copy()
	f.c0.MulByNonresidueAssign()
	f.c0.AddAssign(aa)
}

// FQ12Zero is the zero element of FQ12.
var FQ12Zero = NewFQ12(FQ6Zero, FQ6Zero)

// FQ12One is the one element of FQ12.
var FQ12One = NewFQ12(FQ6One, FQ6Zero)

// Equals checks if two FQ12 elements are equal.
func (f FQ12) Equals(other *FQ12) bool {
	return f.c0.Equals(other.c0) && f.c1.Equals(other.c1)
}

// DoubleAssign doubles each coefficient in an FQ12 element.
func (f FQ12) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
}

// NegAssign negates each coefficient in an FQ12 element.
func (f *FQ12) NegAssign() {
	f.c1.NegAssign()
	f.c0.NegAssign()
}

// AddAssign adds two FQ12 elements together.
func (f FQ12) AddAssign(other *FQ12) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
}

// SubAssign subtracts one FQ12 element from another.
func (f FQ12) SubAssign(other *FQ12) {
	f.c0.SubAssign(other.c0)
	f.c1.SubAssign(other.c1)
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
func (f FQ12) Exp(n *FQRepr) *FQ12 {
	nCopy := n.Copy()
	res := FQ12One.Copy()
	fi := f.Copy()
	for nCopy.Cmp(bigZero) != 0 {
		if !isEven(nCopy) {
			res.MulAssign(fi)
		}
		fi.MulAssign(fi)
		nCopy.Rsh(1)
	}
	return res
}

var bigSix = NewFQRepr(6)

var qFieldModulusBig = QFieldModulus.ToBig()

func getFrobExpMinus1Over6(power int64) *FQRepr {
	out := new(big.Int).Exp(qFieldModulusBig, big.NewInt(power), qFieldModulusBig)
	out.Sub(out, big.NewInt(1))
	out.Div(out, big.NewInt(6))
	out.Mod(out, qFieldModulusBig)
	o, _ := FQReprFromBigInt(out)
	return o
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

// FrobeniusMapAssign calculates the frobenius map of an FQ12 element.
func (f *FQ12) FrobeniusMapAssign(power uint8) {
	f.c0.FrobeniusMapAssign(power)
	f.c1.FrobeniusMapAssign(power)
	f.c1.c0.MulAssign(frobeniusCoeffFQ12c1[power%12])
	f.c1.c1.MulAssign(frobeniusCoeffFQ12c1[power%12])
	f.c1.c2.MulAssign(frobeniusCoeffFQ12c1[power%12])
}

// SquareAssign squares the FQ2 element.
func (f *FQ12) SquareAssign() {
	ab := f.c0.Copy()
	ab.MulAssign(f.c1)
	c0 := f.c1.Copy()
	f.c1.NegAssign()
	c0.AddAssign(f.c0)
	f.c0.AddAssign(f.c1)
	c0.MulAssign(f.c0)
	c0.SubAssign(ab)
	c0.AddAssign(ab)
	ab.AddAssign(ab)
	f.c0 = c0
	f.c1 = ab
}

// MulAssign multiplies two FQ12 elements together.
func (f *FQ12) MulAssign(other *FQ12) {
	aa := f.c0.Copy()
	aa.MulAssign(other.c0)
	bb := f.c1.Copy()
	bb.MulAssign(other.c1)
	o := other.c0.Copy()
	o.AddAssign(other.c1)

	f.c1 = f.c1.Copy()
	f.c1.AddAssign(f.c0)
	f.c1.MulAssign(o)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
	f.c0 = bb.Copy()
	f.c0.MulByNonresidueAssign()
	f.c0.AddAssign(aa)
}

// InverseAssign finds the inverse of an FQ12
func (f FQ12) InverseAssign() bool {
	c1s := f.c1.Copy()
	c1s.SquareAssign()
	c1s.Copy()
	c1s.MulByNonresidueAssign()
	c0s := f.c0.Copy()
	c0s.SquareAssign()
	c0s.SubAssign(c1s)

	i := c0s.Copy()

	if i.InverseAssign() {
		return false
	}

	f.c0.MulAssign(i)
	f.c1.MulAssign(i)
	f.c1.NegAssign()
	return true
}
