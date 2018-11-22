package bls

import (
	"fmt"
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
	return NewFQ6(f.c2.Copy().MultiplyByNonresidue(), f.c0.Copy(), f.c1.Copy())
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
	t1 := c1.Mul(tmp).Sub(b).MultiplyByNonresidue().Add(a)
	tmp = f.c0.Add(f.c2)
	t3 := c0.Mul(tmp).Sub(a).Add(b)
	tmp = f.c0.Add(f.c1)
	t2 := c0.Add(c1).Mul(tmp).Sub(a).Sub(b)

	return NewFQ6(
		t1,
		t2,
		t3,
	)
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

// Neg negates the coefficients of the FQ6 element.
func (f FQ6) Neg() *FQ6 {
	return NewFQ6(
		f.c0.Neg(),
		f.c1.Neg(),
		f.c2.Neg(),
	)
}

// Add adds the coefficients of the FQ6 element to another.
func (f FQ6) Add(other *FQ6) *FQ6 {
	return NewFQ6(
		f.c0.Add(other.c0),
		f.c1.Add(other.c1),
		f.c2.Add(other.c2),
	)
}

// Sub subtracts the coefficients of the FQ6 element from another.
func (f FQ6) Sub(other *FQ6) *FQ6 {
	return NewFQ6(
		f.c0.Sub(other.c0),
		f.c1.Sub(other.c1),
		f.c2.Sub(other.c2),
	)
}

var fq6c10, _ = new(big.Int).SetString("3380320199399472671518931668520476396067793891014375699959770179129436917079669831430077592723774664465579537268733", 10)
var fq6c11, _ = new(big.Int).SetString("3838308620157845674988348254171675463841990158783295489149568424364307477627266222040030802976299176456994858070129", 10)
var fq6c12, _ = new(big.Int).SetString("786190290886016440328299728779656453203981590080344581554777668754318906274739675415266862557957487153214149780712", 10)
var fq6c21, _ = new(big.Int).SetString("3216219264335650953089490096956247703352901229858663303777280467369712744216098189027420766571058176884680122779075", 10)
var fq6c20, _ = new(big.Int).SetString("622089355822194721898858157215427760489088928924632185372287956994594733411168033012610036405240999572314735291054", 10)
var fq6c24, _ = new(big.Int).SetString("786190290886016440328299728779656453203981590080344581554777668754318906274739675415266862557957487153214149780712", 10)
var fq6c25, _ = new(big.Int).SetString("164100935063821718429441571564228692714892661155712396182489711759724172863571642402656826152716487580899414489658", 10)

var frobeniusCoeffFQ6c1 = [6]*FQ2{
	{
		c0: NewFQ(fq6c10),
		c1: FQZero,
	},
	{
		c0: FQZero,
		c1: NewFQ(fq6c11),
	},
	{
		c0: NewFQ(fq6c12),
		c1: FQZero,
	},
	{
		c0: FQZero,
		c1: NewFQ(fq6c10),
	},
	{
		c0: NewFQ(fq6c24),
		c1: FQZero,
	},
	{
		c0: FQZero,
		c1: NewFQ(fq6c12),
	},
}

var frobeniusCoeffFQ6c2 = [6]*FQ2{
	{
		c0: NewFQ(fq6c10),
		c1: FQZero,
	},
	{
		c0: NewFQ(fq6c21),
		c1: FQZero,
	},
	{
		c0: NewFQ(fq6c11),
		c1: FQZero,
	},
	{
		c0: NewFQ(fq6c20),
		c1: FQZero,
	},
	{
		c0: NewFQ(fq6c11),
		c1: FQZero,
	},
	{
		c0: NewFQ(fq6c25),
		c1: FQZero,
	},
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
	t1 := other.c1.Add(other.c2).Mul(tmp).Sub(bb).Sub(cc).MultiplyByNonresidue().Add(aa)
	tmp = f.c0.Add(f.c2)
	t3 := other.c0.Add(other.c2).Mul(tmp).Sub(aa).Add(bb).Sub(cc)
	tmp = f.c0.Add(f.c1)
	t2 := other.c0.Add(other.c1).Mul(tmp).Sub(aa).Sub(bb).Add(cc.MultiplyByNonresidue())

	return NewFQ6(
		t1,
		t2,
		t3,
	)
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

func RandFQ6() (*FQ6, error) {
	c0, err := RandFQ2()
	if err != nil {
		return nil, err
	}
	c1, err := RandFQ2()
	if err != nil {
		return nil, err
	}
	c2, err := RandFQ2()
	if err != nil {
		return nil, err
	}
	return NewFQ6(c0, c1, c2), nil
}
