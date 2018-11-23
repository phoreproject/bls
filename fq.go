package bls

import (
	"fmt"
	"math/big"
)

// FQ is an element in a field.
type FQ struct {
	n *big.Int
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)
var bigTwo = big.NewInt(2)

// FieldModulus is the modulus of the field.
var FieldModulus, _ = new(big.Int).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", 10)

// R = 2**384 % FieldModulus
var R, _ = new(big.Int).SetString("3380320199399472671518931668520476396067793891014375699959770179129436917079669831430077592723774664465579537268733", 10)

// R2 = R**2 % FieldModulus
var R2, _ = new(big.Int).SetString("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750", 10)

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
		lm, low, hm, high = nm, new, lm, low
	}
	return new(big.Int).Mod(lm, n)
}

// NewFQ creates a new field element.
func NewFQ(n *big.Int) *FQ {
	outN := new(big.Int).Mod(n, FieldModulus)
	return &FQ{n: outN}
}

// Copy creates a copy of the field element.
func (f FQ) Copy() *FQ {
	return &FQ{n: new(big.Int).Set(f.n)}
}

// Add adds two field elements together.
func (f FQ) Add(other *FQ) *FQ {
	out := new(big.Int).Add(f.n, other.n)
	out.Mod(out, FieldModulus)
	return &FQ{n: out}
}

// Mul multiplies two field elements together.
func (f FQ) Mul(other *FQ) *FQ {
	out := new(big.Int).Mul(f.n, other.n)
	out.Mod(out, FieldModulus)
	return &FQ{n: out}
}

// Sub subtracts one field element from the other.
func (f FQ) Sub(other *FQ) *FQ {
	out := new(big.Int).Sub(f.n, other.n)
	out.Mod(out, FieldModulus)
	return &FQ{n: out}
}

// Div divides one field element by another.
func (f FQ) Div(other *FQ) *FQ {
	otherInverse := &FQ{n: primeFieldInv(other.n, FieldModulus)}
	return f.Mul(otherInverse)
}

// Exp exponentiates the field element to the given power.
func (f FQ) Exp(n *big.Int) *FQ {
	if n.Cmp(bigZero) == 0 {
		return &FQ{n: new(big.Int).Set(bigOne)}
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
	return NewFQ(new(big.Int).Neg(f.n))
}

func (f FQ) String() string {
	return f.n.String()
}

// Cmp compares this field element to another.
func (f FQ) Cmp(other *FQ) int {
	return f.n.Cmp(other.n)
}

// Double doubles the
func (f FQ) Double() *FQ {
	return NewFQ(new(big.Int).Mul(f.n, bigTwo))
}

// IsZero checks if the field element is zero.
func (f FQ) IsZero() bool {
	return f.n.Cmp(bigZero) == 0
}

// Square squares a field element.
func (f FQ) Square() *FQ {
	return NewFQ(new(big.Int).Mul(f.n, f.n))
}

// Sqrt calculates the square root of the field element.
func (f FQ) Sqrt() *FQ {
	// Shank's algorithm for q mod 4 = 3
	// https://eprint.iacr.org/2012/685.pdf (page 9, algorithm 2)

	a1 := f.Exp(qMinus3Over4)
	a0 := a1.Square().Mul(&f)

	if a0.Equals(NewFQ(negativeOne)) {
		return nil
	}
	return a1.Mul(&f)
}

func isEven(b *big.Int) bool {
	return new(big.Int).Mod(b, bigTwo).Cmp(bigZero) == 0
}

// Inverse finds the inverse of the field element.
func (f FQ) Inverse() *FQ {
	if f.IsZero() {
		return nil
	}
	u := new(big.Int).Set(f.n)
	v := new(big.Int).Set(FieldModulus)
	b := NewFQ(new(big.Int).Set(bigOne))
	c := FQZero.Copy()

	for u.Cmp(bigOne) != 0 && v.Cmp(bigOne) != 0 {
		fmt.Println(b, c)
		for isEven(u) {
			u.Div(u, bigTwo)
			if isEven(b.n) {
				b.n.Div(b.n, bigTwo)
			} else {
				b.n.Add(b.n, FieldModulus)
				b.n.Div(b.n, bigTwo)
			}
		}

		for isEven(v) {
			v.Div(v, bigTwo)
			if isEven(c.n) {
				c.n.Div(c.n, bigTwo)
			} else {
				c.n.Add(c.n, FieldModulus)
				c.n.Div(c.n, bigTwo)
			}
		}

		if u.Cmp(v) >= 0 {
			u.Sub(u, v)
			b = b.Sub(c)
		} else {
			v.Sub(v, u)
			c = c.Sub(b)
		}
	}
	if u.Cmp(bigOne) == 0 {
		return b
	}
	return c
}

var qMinus1Over2, _ = new(big.Int).SetString("2001204777610833696708894912867952078278441409969503942666029068062015825245418932221343814564507832018947136279893", 10)

// LegendreSymbol is the legendre symbol of an element.
type LegendreSymbol uint8

const (
	// LegendreZero is the legendre symbol of zero.
	LegendreZero = LegendreSymbol(iota)

	// LegendreQuadraticResidue is the legendre symbol of quadratic residue.
	LegendreQuadraticResidue

	// LegendreQuadraticNonResidue is the legendre symbol of quadratic non-residue.
	LegendreQuadraticNonResidue
)

// Legendre gets the legendre symbol of the element.
func (f *FQ) Legendre() LegendreSymbol {
	o := f.Exp(qMinus1Over2)
	if o.IsZero() {
		return LegendreZero
	} else if o.Equals(FQOne) {
		return LegendreQuadraticResidue
	} else {
		return LegendreQuadraticNonResidue
	}
}

// FQZero is the FQ at 0.
var FQZero = NewFQ(bigZero)

// FQOne is the FQ at 1.
var FQOne = NewFQ(bigOne)
