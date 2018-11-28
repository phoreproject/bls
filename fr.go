package bls

import (
	"crypto/rand"
	"io"
	"math/big"
)

// FR is an element in a field.
type FR struct {
	n *big.Int
}

// RFieldModulus is the modulus of the field.
var RFieldModulus, _ = new(big.Int).SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

// NewFR creates a new field element.
func NewFR(n *big.Int) *FR {
	outN := new(big.Int).Mod(n, RFieldModulus)
	return &FR{n: outN}
}

// Copy creates a copy of the field element.
func (f FR) Copy() *FR {
	return &FR{n: new(big.Int).Set(f.n)}
}

// Add adds two field elements together.
func (f FR) Add(other *FR) *FR {
	out := new(big.Int).Add(f.n, other.n)
	out.Mod(out, RFieldModulus)
	return &FR{n: out}
}

// AddAssign adds two field elements together.
func (f FR) AddAssign(other *FR) {
	f.n.Add(f.n, other.n)
	f.n.Mod(f.n, RFieldModulus)
}

// Mul multiplies two field elements together.
func (f FR) Mul(other *FR) *FR {
	out := new(big.Int).Mul(f.n, other.n)
	out.Mod(out, RFieldModulus)
	return &FR{n: out}
}

// MulAssign multiplies one field element by the other.
func (f FR) MulAssign(other *FR) {
	f.n.Mul(f.n, other.n)
	f.n.Mod(f.n, RFieldModulus)
}

// Sub subtracts one field element from the other.
func (f FR) Sub(other *FR) *FR {
	out := new(big.Int).Sub(f.n, other.n)
	out.Mod(out, RFieldModulus)
	return &FR{n: out}
}

// SubAssign subtracts one field element from the other.
func (f FR) SubAssign(other *FR) {
	f.n.Sub(f.n, other.n)
	f.n.Mod(f.n, RFieldModulus)
}

// Div divides one field element by another.
func (f FR) Div(other *FR) *FR {
	otherInverse := &FR{n: primeFieldInv(other.n, RFieldModulus)}
	return f.Mul(otherInverse)
}

// Exp exponentiates the field element to the given power.
func (f FR) Exp(n *big.Int) *FR {
	if n.Cmp(bigZero) == 0 {
		return &FR{n: new(big.Int).Set(bigOne)}
	} else if n.Cmp(bigOne) == 0 {
		return f.Copy()
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigZero) == 0 {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo))
	} else {
		return f.Mul(&f).Exp(new(big.Int).Div(n, bigTwo)).Mul(&f)
	}
}

// Equals checks equality of two field elements.
func (f FR) Equals(other *FR) bool {
	return f.n.Cmp(other.n) == 0
}

// Neg gets the negative value of the field element mod RFieldModulus.
func (f FR) Neg() *FR {
	return NewFR(new(big.Int).Neg(f.n))
}

func (f FR) String() string {
	return f.n.String()
}

// Cmp compares this field element to another.
func (f FR) Cmp(other *FR) int {
	return f.n.Cmp(other.n)
}

// Double doubles the
func (f FR) Double() *FR {
	return NewFR(new(big.Int).Lsh(f.n, 1))
}

// IsZero checks if the field element is zero.
func (f FR) IsZero() bool {
	return f.n.Cmp(bigZero) == 0
}

// Square squares a field element.
func (f FR) Square() *FR {
	return NewFR(new(big.Int).Mul(f.n, f.n))
}

// Sqrt calculates the square root of the field element.
func (f FR) Sqrt() *FR {
	// Shank's algorithm for q mod 4 = 3
	// https://eprint.iacr.org/2012/685.pdf (page 9, algorithm 2)

	a1 := f.Exp(qMinus3Over4)
	a0 := a1.Square().Mul(&f)

	if a0.Equals(NewFR(negativeOne)) {
		return nil
	}
	return a1.Mul(&f)
}

// Inverse finds the inverse of the field element.
func (f FR) Inverse() *FR {
	if f.IsZero() {
		return nil
	}
	u := new(big.Int).Set(f.n)
	v := new(big.Int).Set(RFieldModulus)
	b := NewFR(new(big.Int).Set(bigOne))
	c := FRZero.Copy()

	for u.Cmp(bigOne) != 0 && v.Cmp(bigOne) != 0 {
		for isEven(u) {
			u.Div(u, bigTwo)
			if isEven(b.n) {
				b.n.Div(b.n, bigTwo)
			} else {
				b.n.Add(b.n, RFieldModulus)
				b.n.Div(b.n, bigTwo)
			}
		}

		for isEven(v) {
			v.Div(v, bigTwo)
			if isEven(c.n) {
				c.n.Div(c.n, bigTwo)
			} else {
				c.n.Add(c.n, RFieldModulus)
				c.n.Div(c.n, bigTwo)
			}
		}

		if u.Cmp(v) >= 0 {
			u.Sub(u, v)
			b.SubAssign(c)
		} else {
			v.Sub(v, u)
			c.SubAssign(b)
		}
	}
	if u.Cmp(bigOne) == 0 {
		return b
	}
	return c
}

var rMinus1Over2, _ = new(big.Int).SetString("26217937587563095239723870254092982918845276250263818911301829349969290592256", 10)

// Legendre gets the legendre symbol of the element.
func (f *FR) Legendre() LegendreSymbol {
	o := f.Exp(rMinus1Over2)
	if o.IsZero() {
		return LegendreZero
	} else if o.Equals(FROne) {
		return LegendreQuadraticResidue
	} else {
		return LegendreQuadraticNonResidue
	}
}

// ToBig converts the FR element to the underlying big number.
func (f *FR) ToBig() *big.Int {
	return new(big.Int).Set(f.n)
}

// RandFR generates a random FR element.
func RandFR(reader io.Reader) (*FR, error) {
	r, err := rand.Int(reader, RFieldModulus)
	if err != nil {
		return nil, err
	}
	return NewFR(r), nil
}

// FRZero is the FR at 0.
var FRZero = NewFR(bigZero)

// FROne is the FR at 1.
var FROne = NewFR(bigOne)
