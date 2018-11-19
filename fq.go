package bls

import (
	"errors"
	"math/big"
)

// FQ is an element in a field.
type FQ struct {
	n            *big.Int
	fieldModulus *big.Int
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

var fieldModulus, _ = new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)

func primeFieldInv(a *big.Int, n *big.Int) *big.Int {
	if a.Cmp(bigZero) == 0 {
		return big.NewInt(0)
	}
	lm := big.NewInt(1)
	hm := big.NewInt(0)
	low := new(big.Int).Mod(a, n)
	high := new(big.Int).Set(n)
	for low.Cmp(bigOne) > 0 {
		r := new(big.Int).Div(high, low)
		nm := new(big.Int).Sub(hm, new(big.Int).Mul(lm, r))
		new := new(big.Int).Sub(high, new(big.Int).Mul(low, r))
		lm = nm
		low = new
		hm = lm
		high = low
	}
	return new(big.Int).Mod(lm, n)
}

// NewFQ creates a new field element.
func NewFQ(n *big.Int, fieldModulus *big.Int) *FQ {
	outN := new(big.Int).Mod(n, fieldModulus)
	return &FQ{n: outN, fieldModulus: fieldModulus}
}

// Copy creates a copy of the field element.
func (f FQ) Copy() *FQ {
	return &FQ{n: new(big.Int).Set(f.n), fieldModulus: new(big.Int).Set(f.fieldModulus)}
}

// Add adds two field elements together.
func (f FQ) Add(other *FQ) *FQ {
	out := new(big.Int).Add(f.n, other.n)
	out.Mod(out, f.fieldModulus)
	return &FQ{n: out, fieldModulus: f.fieldModulus}
}

// Mul multiplies two field elements together.
func (f FQ) Mul(other *FQ) *FQ {
	out := new(big.Int).Mul(f.n, other.n)
	out.Mod(out, f.fieldModulus)
	return &FQ{n: out, fieldModulus: f.fieldModulus}
}

// Sub subtracts one field element from the other.
func (f FQ) Sub(other *FQ) *FQ {
	out := new(big.Int).Sub(f.n, other.n)
	out.Mod(out, f.fieldModulus)
	return &FQ{n: out, fieldModulus: f.fieldModulus}
}

// Div divides one field element by another.
func (f FQ) Div(other *FQ) *FQ {
	otherInverse := &FQ{n: primeFieldInv(other.n, f.fieldModulus), fieldModulus: f.fieldModulus}
	return f.Mul(otherInverse)
}

// Exp exponentiates the field element to the given power.
func (f FQ) Exp(n *big.Int) *FQ {
	if n.Cmp(bigZero) == 0 {
		return &FQ{n: new(big.Int).Set(bigOne), fieldModulus: f.fieldModulus}
	} else if n.Cmp(bigOne) == 0 {
		return f.Copy()
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigZero) == 0 {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo))
	} else {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo)).Mul(&f)
	}
}

// Equals checks equality of two field elements.
func (f FQ) Equals(other *FQ) bool {
	return f.n.Cmp(other.n) == 0
}

// Neg gets the negative value of the field element mod fieldModulus.
func (f FQ) Neg() *FQ {
	return NewFQ(new(big.Int).Neg(f.n), f.fieldModulus)
}

func (f FQ) String() string {
	return f.n.String()
}

// Polynomial is a polynomial with certain coefficients.
type Polynomial []*big.Int

// Deg gets the degree of the polynomial.
func (p Polynomial) Deg() int {
	d := len(p) - 1
	for p[d].Cmp(bigZero) == 0 && d > 0 {
		d--
	}
	return d
}

// FQP is a polynomial with field element coefficients.
type FQP struct {
	elements            []*FQ
	modulusCoefficients []*big.Int
	mcs                 map[int]int
	degree              int
}

// NewFQP creates a new polynomial with field element coefficients.
func NewFQP(elements []*FQ, modulusCoefficients []*big.Int, mcs map[int]int) (*FQP, error) {
	return NewFQPWithDegree(elements, modulusCoefficients, mcs, len(modulusCoefficients))
}

// NewFQPWithDegree creates a new polynomial with the specified degree.
func NewFQPWithDegree(elements []*FQ, modulusCoefficients []*big.Int, mcs map[int]int, degree int) (*FQP, error) {
	if len(elements) == 0 {
		return nil, errors.New("FQP cannot have 0 elements")
	}
	return &FQP{elements: elements, mcs: mcs, modulusCoefficients: modulusCoefficients, degree: degree}, nil

}

// Copy creates a copy of the FQP provided.
func (f FQP) Copy() *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i, e := range f.elements {
		newElements[i] = e.Copy()
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// Add adds two FQp's together.
func (f FQP) Add(other *FQP) *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i, e := range f.elements {
		newElements[i] = e.Add(other.elements[i])
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// Sub subtracts one FQP from another.
func (f FQP) Sub(other *FQP) *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i, e := range f.elements {
		newElements[i] = e.Sub(other.elements[i])
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// MulScalar multiplies each element in an FQP by a scalar.
func (f FQP) MulScalar(scalar *FQ) *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i, e := range f.elements {
		newElements[i] = e.Mul(scalar)
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// Mul multiplies two polynomials together.
func (f FQP) Mul(other *FQP) *FQP {
	newElements := make([]*FQ, f.degree*2-1)
	for i, eli := range f.elements {
		for j, elj := range other.elements {
			toAdd := eli.Mul(elj)
			if newElements[i+j] == nil {
				newElements[i+j] = toAdd
			} else {
				newElements[i+j] = newElements[i+j].Add(toAdd)
			}
		}
	}

	for exp := f.degree - 2; exp > -1; exp-- {
		top := newElements[len(newElements)-1]
		newElements = newElements[:len(newElements)-1]
		for i, c := range f.mcs {
			newElements[exp+i] = newElements[exp+i].Sub(top.Mul(&FQ{n: big.NewInt(int64(c)), fieldModulus: fieldModulus}))
		}
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// Exp exponentiates the polynomial to a power.
func (f FQP) Exp(other *big.Int) *FQP {
	o := &FQP{elements: append([]*FQ{
		{
			n:            big.NewInt(1),
			fieldModulus: fieldModulus,
		},
	}, ZerosFQ(f.degree-1)...), mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
	t := &f
	otherCopy := new(big.Int).Set(other)
	for otherCopy.Cmp(bigZero) > 0 {
		if new(big.Int).And(otherCopy, bigOne).Cmp(bigZero) > 0 {
			o = o.Mul(t) // t.Deg() = 3 o.Deg() = 2
		}
		otherCopy.Rsh(otherCopy, 1)
		t = t.Mul(t)
	}
	return o
}

// DivScalar multiplies each element by the prime field inverse of the scalar.
func (f FQP) DivScalar(scalar *FQ) *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i, e := range f.elements {
		newElements[i] = e.Mul(&FQ{n: primeFieldInv(scalar.n, f.elements[0].fieldModulus), fieldModulus: fieldModulus})
	}
	return &FQP{elements: newElements, mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}
}

// Div multiplies the polynomial by the inverse of the argument.
func (f FQP) Div(other *FQP) *FQP {
	return f.Mul(other.Inv())
}

// Zeros puts a bunch of big.Int zeros in an array.
func Zeros(num int) []*big.Int {
	out := make([]*big.Int, num)
	for i := 0; i < num; i++ {
		out[i] = new(big.Int).Set(bigZero)
	}
	return out
}

// ZerosFQ puts a bunch of FQ zeros in an array.
func ZerosFQ(num int) []*FQ {
	out := make([]*FQ, num)
	for i := 0; i < num; i++ {
		out[i] = &FQ{n: new(big.Int).Set(bigZero), fieldModulus: fieldModulus}
	}
	return out
}

func polyRoundedDiv(a []*big.Int, b []*big.Int, mod *big.Int) []*big.Int {
	degA := Polynomial(a).Deg()
	degB := Polynomial(b).Deg()
	temp := make([]*big.Int, len(a))
	o := make([]*big.Int, len(a))
	for i, x := range a {
		temp[i] = new(big.Int).Set(x)
		o[i] = new(big.Int).Set(bigZero)
	}

	for i := degA - degB; i > -1; i-- {
		o[i].Add(o[i], new(big.Int).Mul(temp[degB+i], primeFieldInv(b[degB], mod)))
		for c := 0; c < degB+1; c++ {
			temp[c+i].Sub(temp[c+i], o[c])
		}
	}
	o = o[:Polynomial(o).Deg()+1]
	for i := range o {
		o[i].Mod(o[i], mod)
	}
	return o
}

// Inv uses the extended euclidean algorithm to find the modular inverse.
func (f FQP) Inv() *FQP {
	lm := append([]*big.Int{new(big.Int).Set(bigOne)}, Zeros(f.degree)...)
	hm := Zeros(f.degree + 1)
	low := make([]*big.Int, len(f.elements)+1)
	high := make([]*big.Int, len(f.elements)+1)
	for i := range f.elements {
		low[i] = new(big.Int).Set(f.elements[i].n)
		high[i] = new(big.Int).Set(f.modulusCoefficients[i])
	}
	high[len(high)-1] = new(big.Int).Set(bigOne)
	low[len(low)-1] = new(big.Int).Set(bigZero)

	for Polynomial(low).Deg() > 0 {
		r := polyRoundedDiv(high, low, fieldModulus)
		r = append(r, Zeros(f.degree+1-len(r))...)
		nm := make([]*big.Int, len(hm))
		for i := range hm {
			nm[i] = new(big.Int).Set(hm[i])
		}
		n := make([]*big.Int, len(high))
		for i := range high {
			n[i] = new(big.Int).Set(high[i])
		}
		for i := 0; i < f.degree+1; i++ {
			for j := 0; j < f.degree+1-i; j++ {
				nm[i+j].Sub(nm[i+j], new(big.Int).Mul(lm[i], r[j]))
				n[i+j].Sub(n[i+j], new(big.Int).Mul(low[i], r[j]))
			}
		}
		for i := range nm {
			nm[i].Mod(nm[i], fieldModulus)
		}
		for i := range n {
			n[i].Mod(n[i], fieldModulus)
		}
		lm, low, hm, high = nm, n, lm, low
	}

	lmFQ := make([]*FQ, len(lm))
	for i := range lmFQ {
		lmFQ[i] = &FQ{n: lm[i], fieldModulus: fieldModulus}
	}

	out := FQP{elements: lmFQ[:f.degree], mcs: f.mcs, modulusCoefficients: f.modulusCoefficients, degree: f.degree}.DivScalar(&FQ{n: low[0], fieldModulus: fieldModulus})
	return out
}

// Equals checks if two FQPs' coefficients are equal.
func (f FQP) Equals(other *FQP) bool {
	if len(f.elements) != len(other.elements) {
		return false
	}
	for i, e := range f.elements {
		if !e.Equals(other.elements[i]) {
			return false
		}
	}
	return true
}

// Neg negates each coefficient in the FQP.
func (f FQP) Neg() *FQP {
	newElements := make([]*FQ, len(f.elements))
	for i := range newElements {
		newElements[i] = newElements[i].Neg()
	}
	newF, _ := NewFQPWithDegree(newElements, f.modulusCoefficients, f.mcs, f.degree)
	return newF
}

func (f FQP) String() string {
	s := "FQP["
	for _, i := range f.elements {
		s += i.n.String() + ", "
	}
	return s + "]"
}

// NewFQ2 creates a new quadratic field extension.
func NewFQ2(coeffs []*FQ) (*FQP, error) {
	mcTuples := make(map[int]int)
	mcTuples[0] = 1
	if len(coeffs) != 2 {
		return nil, errors.New("wrong number of coefficients for quadratic polynomial")
	}
	modulusCoefficients := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
	}
	return NewFQP(coeffs, modulusCoefficients, mcTuples)
}

// FQ2One returns the 1-value FQ2
func FQ2One() *FQP {
	f, _ := NewFQ2([]*FQ{
		&FQ{n: big.NewInt(1), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
	})
	return f
}

// FQ2Zero returns the 0-value FQ2
func FQ2Zero() *FQP {
	f, _ := NewFQ2([]*FQ{
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
	})
	return f
}

// FQ12Zero returns the 0-value FQ12
func FQ12Zero() *FQP {
	f, _ := NewFQ12([]*FQ{
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
		&FQ{n: big.NewInt(0), fieldModulus: fieldModulus},
	})
	return f
}

// NewFQ12 creates a new 12th-degree field extension.
func NewFQ12(coeffs []*FQ) (*FQP, error) {
	mcTuples := make(map[int]int)
	mcTuples[0] = 1
	if len(coeffs) != 12 {
		return nil, errors.New("wrong number of coefficients for quadratic polynomial")
	}
	modulusCoefficients := []*big.Int{
		big.NewInt(82),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(-18),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
	}
	return NewFQP(coeffs, modulusCoefficients, mcTuples)
}
