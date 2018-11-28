package bls

import (
	"fmt"
	"io"
	"math/big"
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

// MulByNonresidue multiplies by quadratic nonresidue v.
func (f FQ6) MulByNonresidue() *FQ6 {
	out := NewFQ6(f.c2.Copy(), f.c0.Copy(), f.c1.Copy())
	out.c0.MultiplyByNonresidueAssign()
	return out
}

// MulBy1 multiplies the FQ6 by an FQ2.
func (f FQ6) MulBy1(c1 *FQ2) *FQ6 {
	b := f.c1.Mul(c1)
	tmp := f.c1.Add(f.c2)
	t1 := c1.Mul(tmp).Sub(b).MultiplyByNonresidue()
	tmp = f.c0.Add(f.c1)
	t2 := c1.Mul(tmp).Sub(b)
	return NewFQ6(t1, t2, b)
}

// MulBy01 multiplies by c0 and c1.
func (f FQ6) MulBy01(c0 *FQ2, c1 *FQ2) *FQ6 {
	a := f.c0.Mul(c0)
	b := f.c1.Mul(c1)

	tmp := f.c1.Add(f.c2)
	t1 := c1.Mul(tmp)
	t1.SubAssign(b)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(a)
	tmp = f.c0.Add(f.c2)
	t3 := c0.Mul(tmp)
	t3.SubAssign(a)
	t3.AddAssign(b)
	tmp = f.c0.Add(f.c1)
	t2 := c0.Add(c1)
	t2.MulAssign(tmp)
	t2.SubAssign(a)
	t2.SubAssign(b)

	return NewFQ6(
		t1,
		t2,
		t3,
	)
}

// MulBy01Assign multiplies by c0 and c1.
func (f *FQ6) MulBy01Assign(c0 *FQ2, c1 *FQ2) {
	a := f.c0.Mul(c0)
	b := f.c1.Mul(c1)

	tmp := f.c1.Add(f.c2)
	t1 := c1.Mul(tmp)
	t1.SubAssign(b)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(a)
	tmp = f.c0.Add(f.c2)
	t3 := c0.Mul(tmp)
	t3.SubAssign(a)
	t3.AddAssign(b)
	tmp = f.c0.Add(f.c1)
	t2 := c0.Add(c1)
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

// Double doubles the coefficients of the FQ6 element.
func (f FQ6) Double() *FQ6 {
	return NewFQ6(
		f.c0.Double(),
		f.c1.Double(),
		f.c2.Double(),
	)
}

// DoubleAssign doubles the coefficients of the FQ6 element.
func (f FQ6) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
	f.c2.DoubleAssign()
}

// Neg negates the coefficients of the FQ6 element.
func (f FQ6) Neg() *FQ6 {
	return NewFQ6(
		f.c0.Neg(),
		f.c1.Neg(),
		f.c2.Neg(),
	)
}

// NegAssign negates the coefficients of the FQ6 element.
func (f FQ6) NegAssign() {
	f.c0.NegAssign()
	f.c1.NegAssign()
	f.c2.NegAssign()
}

// Add adds the coefficients of the FQ6 element to another.
func (f FQ6) Add(other *FQ6) *FQ6 {
	return NewFQ6(
		f.c0.Add(other.c0),
		f.c1.Add(other.c1),
		f.c2.Add(other.c2),
	)
}

// AddAssign the coefficients of the FQ6 element to another.
func (f FQ6) AddAssign(other *FQ6) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
	f.c2.AddAssign(other.c2)
}

// Sub subtracts the coefficients of the FQ6 element from another.
func (f FQ6) Sub(other *FQ6) *FQ6 {
	return NewFQ6(
		f.c0.Sub(other.c0),
		f.c1.Sub(other.c1),
		f.c2.Sub(other.c2),
	)
}

// SubAssign subtracts the coefficients of the FQ6 element from another.
func (f FQ6) SubAssign(other *FQ6) {
	f.c0.SubAssign(other.c0)
	f.c1.SubAssign(other.c1)
	f.c2.SubAssign(other.c2)
}

func getFrobExpMinus1Over3(power int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Sub(new(big.Int).Exp(QFieldModulus, big.NewInt(power), nil), bigOne), bigThree)
}

func get2FrobExpMinus2Over3(power int64) *big.Int {
	qPow := new(big.Int).Exp(QFieldModulus, big.NewInt(power), nil)
	qPowDouble := new(big.Int).Mul(qPow, bigTwo)
	qPowMinusTwo := new(big.Int).Sub(qPowDouble, bigTwo)
	return new(big.Int).Div(qPowMinusTwo, bigThree)
}

var bigThree = big.NewInt(3)

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

// FrobeniusMap runs the frobenius map algorithm with a certain power.
func (f FQ6) FrobeniusMap(power uint8) *FQ6 {
	n0 := f.c0.FrobeniusMap(power)
	n1 := f.c1.FrobeniusMap(power)
	n2 := f.c2.FrobeniusMap(power)
	return NewFQ6(
		n0,
		n1.Mul(frobeniusCoeffFQ6c1[power%6]),
		n2.Mul(frobeniusCoeffFQ6c2[power%6]),
	)
}

// FrobeniusMapAssign runs the frobenius map algorithm with a certain power.
func (f FQ6) FrobeniusMapAssign(power uint8) {
	f.c0.FrobeniusMapAssign(power)
	f.c1.FrobeniusMapAssign(power)
	f.c1.MulAssign(frobeniusCoeffFQ6c1[power%6])
	f.c2.FrobeniusMapAssign(power)
	f.c2.MulAssign(frobeniusCoeffFQ6c2[power%6])
}

// Square squares the FQ6 element.
func (f FQ6) Square() *FQ6 {
	s0 := f.c0.Square()
	ab := f.c0.Mul(f.c1)
	s1 := ab.Double()
	s2 := f.c0.Sub(f.c1).Add(f.c2).Square()
	bc := f.c1.Mul(f.c2)
	s3 := bc.Double()
	s4 := f.c2.Square()

	return NewFQ6(
		s3.MultiplyByNonresidue().Add(s0),
		s4.MultiplyByNonresidue().Add(s1),
		s1.Add(s2).Add(s3).Sub(s0).Sub(s4),
	)
}

// Mul multiplies two FQ6 elements together.
func (f FQ6) Mul(other *FQ6) *FQ6 {
	aa := f.c0.Mul(other.c0)
	bb := f.c1.Mul(other.c1)
	cc := f.c2.Mul(other.c2)

	tmp := f.c1.Add(f.c2)
	t1 := other.c1.Add(other.c2)
	t1.MulAssign(tmp)
	t1.SubAssign(bb)
	t1.SubAssign(cc)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(aa)

	tmp = f.c0.Add(f.c2)
	t3 := other.c0.Add(other.c2)
	t3.MulAssign(tmp)
	t3.SubAssign(aa)
	t3.AddAssign(bb)
	t3.SubAssign(cc)
	tmp = f.c0.Add(f.c1)
	t2 := other.c0.Add(other.c1)
	t2.MulAssign(tmp)
	t2.SubAssign(aa)
	t2.SubAssign(bb)
	cc.MultiplyByNonresidueAssign()
	t2.AddAssign(cc)

	return NewFQ6(
		t1,
		t2,
		t3,
	)
}

// MulAssign multiplies two FQ6 elements together.
func (f *FQ6) MulAssign(other *FQ6) {
	aa := f.c0.Mul(other.c0)
	bb := f.c1.Mul(other.c1)
	cc := f.c2.Mul(other.c2)

	tmp := f.c1.Add(f.c2)
	t1 := other.c1.Add(other.c2)
	t1.MulAssign(tmp)
	t1.SubAssign(bb)
	t1.SubAssign(cc)
	t1.MultiplyByNonresidueAssign()
	t1.AddAssign(aa)

	tmp = f.c0.Add(f.c2)
	f.c2 = other.c0.Add(other.c2)
	f.c2.MulAssign(tmp)
	f.c2.SubAssign(aa)
	f.c2.AddAssign(bb)
	f.c2.SubAssign(cc)
	tmp = f.c0.Add(f.c1)
	f.c1 = other.c0.Add(other.c1)
	f.c1.MulAssign(tmp)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
	cc.MultiplyByNonresidueAssign()
	f.c1.AddAssign(cc)

	f.c0 = t1
}

// Inverse finds the inverse of the FQ6 element.
func (f FQ6) Inverse() *FQ6 {
	c0 := f.c2.MultiplyByNonresidue().Mul(f.c1).Neg()
	c0 = c0.Add(f.c0.Square())
	c1 := f.c2.Square().MultiplyByNonresidue()
	c1 = c1.Sub(f.c0.Mul(f.c1))
	c2 := f.c1.Square().Sub(f.c0.Mul(f.c2))

	tmp := f.c2.Mul(c1).Add(f.c1.Mul(c2)).MultiplyByNonresidue().Add(f.c0.Mul(c0))
	tmpInverse := tmp.Inverse()
	if tmpInverse == nil {
		return nil
	}
	return NewFQ6(tmpInverse.Mul(c0), tmpInverse.Mul(c1), tmpInverse.Mul(c2))
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
