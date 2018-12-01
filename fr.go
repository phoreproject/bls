package bls

import (
	"crypto/rand"
	"fmt"
	"hash"
	"io"
)

// FR is an element in a field.
type FR struct {
	n *FRRepr
}

// RFieldModulus is the modulus of the R field.
var RFieldModulus, _ = FRReprFromString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

// IsValid checks if the element is valid.
func (f *FR) IsValid() bool {
	return f.n.Cmp(RFieldModulus) < 0
}

func (f *FR) reduceAssign() {
	if !f.IsValid() {
		f.n.SubNoBorrow(RFieldModulus)
	}
}

// Copy copies an FR element.
func (f *FR) Copy() *FR {
	return &FR{f.n.Copy()}
}

// FRReprToFR gets a pointer to a FR given a pointer
// to an FQRepr
func FRReprToFR(o *FRRepr) *FR {
	return &FR{n: o}
}

// AddAssign multiplies a field element by this one.
func (f FR) AddAssign(other *FR) {
	f.n.AddNoCarry(other.n)
	f.reduceAssign()
}

const montInvFR = uint64(0xfffffffeffffffff)

func (f FR) montReduce(r0 uint64, r1 *uint64, r2 *uint64, r3 *uint64, r4 *uint64, r5 *uint64, r6 *uint64, r7 *uint64) {
	k := r0 * montInvFR
	carry := uint64(0)
	MACWithCarry(r0, k, RFieldModulus[0], &carry)
	*r1 = MACWithCarry(*r1, k, RFieldModulus[1], &carry)
	*r2 = MACWithCarry(*r2, k, RFieldModulus[2], &carry)
	*r3 = MACWithCarry(*r3, k, RFieldModulus[3], &carry)
	*r4 = AddWithCarry(*r4, 0, &carry)
	carry2 := carry
	k = *r1 * montInvFR
	carry = 0
	MACWithCarry(*r1, k, RFieldModulus[0], &carry)
	*r2 = MACWithCarry(*r2, k, RFieldModulus[1], &carry)
	*r3 = MACWithCarry(*r3, k, RFieldModulus[2], &carry)
	*r4 = MACWithCarry(*r4, k, RFieldModulus[3], &carry)
	*r5 = AddWithCarry(*r5, carry2, &carry)
	carry2 = carry
	k = *r2 * montInvFR
	carry = 0
	MACWithCarry(*r2, k, RFieldModulus[0], &carry)
	*r3 = MACWithCarry(*r3, k, RFieldModulus[1], &carry)
	*r4 = MACWithCarry(*r4, k, RFieldModulus[2], &carry)
	*r5 = MACWithCarry(*r5, k, RFieldModulus[3], &carry)
	*r6 = AddWithCarry(*r6, carry2, &carry)
	carry2 = carry
	k = *r3 * montInvFR
	carry = 0
	MACWithCarry(*r3, k, RFieldModulus[0], &carry)
	*r4 = MACWithCarry(*r4, k, RFieldModulus[1], &carry)
	*r5 = MACWithCarry(*r5, k, RFieldModulus[2], &carry)
	*r6 = MACWithCarry(*r6, k, RFieldModulus[3], &carry)
	*r7 = AddWithCarry(*r7, carry2, &carry)
	f.n[0] = *r4
	f.n[1] = *r5
	f.n[2] = *r6
	f.n[3] = *r7
	f.reduceAssign()
}

// MulAssign multiplies a field element by this one.
func (f FR) MulAssign(other *FR) {
	carry := uint64(0)
	r0 := MACWithCarry(0, f.n[0], other.n[0], &carry)
	r1 := MACWithCarry(0, f.n[0], other.n[1], &carry)
	r2 := MACWithCarry(0, f.n[0], other.n[2], &carry)
	r3 := MACWithCarry(0, f.n[0], other.n[3], &carry)
	r4 := carry
	carry = 0
	r1 = MACWithCarry(r1, f.n[1], other.n[0], &carry)
	r2 = MACWithCarry(r2, f.n[1], other.n[1], &carry)
	r3 = MACWithCarry(r3, f.n[1], other.n[2], &carry)
	r4 = MACWithCarry(r4, f.n[1], other.n[3], &carry)
	r5 := carry
	carry = 0
	r2 = MACWithCarry(r2, f.n[2], other.n[0], &carry)
	r3 = MACWithCarry(r3, f.n[2], other.n[1], &carry)
	r4 = MACWithCarry(r4, f.n[2], other.n[2], &carry)
	r5 = MACWithCarry(r5, f.n[2], other.n[3], &carry)
	r6 := carry
	carry = 0
	r3 = MACWithCarry(r3, f.n[3], other.n[0], &carry)
	r4 = MACWithCarry(r4, f.n[3], other.n[1], &carry)
	r5 = MACWithCarry(r5, f.n[3], other.n[2], &carry)
	r6 = MACWithCarry(r6, f.n[3], other.n[3], &carry)
	r7 := carry
	f.montReduce(r0, &r1, &r2, &r3, &r4, &r5, &r6, &r7)
}

// SubAssign subtracts a field element from this one.
func (f FR) SubAssign(other *FR) {
	if other.n.Cmp(f.n) > 0 {
		f.n.AddNoCarry(RFieldModulus)
	}
	f.n.SubNoBorrow(other.n)
}

var frOne = NewFRRepr(1)
var frZero = NewFRRepr(0)
var bigOneFR = FRReprToFR(frOne)
var bigZeroFR = FRReprToFR(frZero)

// Exp raises the element to a specific power.
func (f *FR) Exp(n *FRRepr) *FR {
	nCopy := n.Copy()
	fi := f.Copy()
	fNew := bigOneFR.Copy()
	for nCopy.Cmp(frZero) != 0 {
		if nCopy.IsOdd() {
			fNew.MulAssign(fi)
		}
		fi.MulAssign(fi)
		nCopy.Rsh(1)
	}
	return fNew
}

// Equals checks equality of two field elements.
func (f FR) Equals(other *FR) bool {
	return f.n.Equals(other.n)
}

// NegAssign gets the negative value of the field element mod RFieldModulus.
func (f *FR) NegAssign() {
	if !f.IsZero() {
		tmp := RFieldModulus.Copy()
		tmp.SubNoBorrow(f.n)
		f.n = tmp
	}
}

func (f FR) String() string {
	return fmt.Sprintf("Fr(0x%096x)", f.n)
}

// Cmp compares this field element to another.
func (f FR) Cmp(other *FR) int {
	return f.n.Cmp(other.n)
}

// Double doubles the element
func (f FR) Double() *FR {
	ret := f.Copy()
	ret.n.Mul2()
	ret.reduceAssign()
	return ret
}

// DoubleAssign doubles the element
func (f *FR) DoubleAssign() {
	f.n.Mul2()
	f.reduceAssign()
}

// IsZero checks if the field element is zero.
func (f FR) IsZero() bool {
	return f.n.Cmp(frZero) == 0
}

// SquareAssign squares a field element.
func (f FR) SquareAssign() {
	carry := uint64(0)
	r1 := MACWithCarry(0, f.n[0], f.n[1], &carry)
	r2 := MACWithCarry(0, f.n[0], f.n[2], &carry)
	r3 := MACWithCarry(0, f.n[0], f.n[3], &carry)
	r4 := carry
	carry = 0
	r3 = MACWithCarry(0, f.n[1], f.n[2], &carry)
	r4 = MACWithCarry(0, f.n[1], f.n[3], &carry)
	r5 := carry
	carry = 0
	r5 = MACWithCarry(0, f.n[2], f.n[3], &carry)
	r6 := carry
	r7 := r6 >> 63
	r6 = (r6 << 1) | (r5 >> 63)
	r5 = (r5 << 1) | (r4 >> 63)
	r4 = (r4 << 1) | (r3 >> 63)
	r3 = (r3 << 1) | (r2 >> 63)
	r2 = (r2 << 1) | (r1 >> 63)
	r1 = r1 << 1

	carry = 0
	r0 := MACWithCarry(0, f.n[0], f.n[0], &carry)
	r1 = AddWithCarry(r1, 0, &carry)
	r2 = MACWithCarry(r2, f.n[1], f.n[1], &carry)
	r3 = AddWithCarry(r3, 0, &carry)
	r4 = MACWithCarry(r4, f.n[2], f.n[2], &carry)
	r5 = AddWithCarry(r5, 0, &carry)
	r6 = MACWithCarry(r6, f.n[3], f.n[3], &carry)
	r7 = AddWithCarry(r7, 0, &carry)
	f.montReduce(r0, &r1, &r2, &r3, &r4, &r5, &r6, &r7)
}

// Sqrt calculates the square root of the field element.
func (f FR) Sqrt() *FR {
	// TODO: fixme
	return FRReprToFR(frOne)
}

// Inverse finds the inverse of the field element.
func (f FR) Inverse() *FR {
	if f.IsZero() {
		return nil
	}
	u := f.n.Copy()
	v := RFieldModulus.Copy()
	b := bigOneFR.Copy()
	c := bigZeroFR.Copy()

	for u.Cmp(frOne) != 0 && v.Cmp(frOne) != 0 {
		for u.IsEven() {
			u.Div2()
			if b.n.IsEven() {
				b.n.Div2()
			} else {
				b.n.AddNoCarry(RFieldModulus)
				b.n.Div2()
			}
		}

		for v.IsEven() {
			v.Div2()
			if c.n.IsEven() {
				c.n.Div2()
			} else {
				c.n.AddNoCarry(RFieldModulus)
				c.n.Div2()
			}
		}

		if u.Cmp(v) >= 0 {
			u.SubNoBorrow(v)
			b.SubAssign(c)
		} else {
			v.SubNoBorrow(u)
			c.SubAssign(b)
		}
	}
	if u.IsZero() {
		return b
	}
	return c
}

// Parity checks if the point is greater than the point negated.
func (f FR) Parity() bool {
	neg := f.Copy()
	neg.NegAssign()
	return f.Cmp(neg) > 0
}

// MulBits multiplies the number by a big number.
func (f FR) MulBits(b *FRRepr) *FR {
	res := bigZeroFR.Copy()
	for i := uint(0); i < b.BitLen(); i++ {
		res.DoubleAssign()
		if b.Bit(i) {
			res.AddAssign(&f)
		}
	}
	return res
}

// HashFR calculates a new FR value based on a hash.
func HashFR(hasher hash.Hash) *FR {
	digest := hasher.Sum(nil)
	newB, err := FRReprFromBytes(digest)
	if err != nil {
		panic(err)
	}
	return FROne.MulBits(newB)
}

var rMinus1Over2, _ = FRReprFromString("26217937587563095239723870254092982918845276250263818911301829349969290592256", 10)

// Legendre gets the legendre symbol of the element.
func (f *FR) Legendre() LegendreSymbol {
	o := f.Exp(rMinus1Over2)
	if o.IsZero() {
		return LegendreZero
	} else if o.Equals(bigOneFR) {
		return LegendreQuadraticResidue
	} else {
		return LegendreQuadraticNonResidue
	}
}

// ToRepr gets the 256-bit representation of the field element.
func (f *FR) ToRepr() *FRRepr {
	return f.n.Copy()
}

// RandFR generates a random FQ element.
func RandFR(reader io.Reader) (*FR, error) {
	r, err := rand.Int(reader, RFieldModulus.ToBig())
	if err != nil {
		return nil, err
	}
	b, _ := FRReprFromBigInt(r)
	return FRReprToFR(b), nil
}

// FRZero is the FR at 0.
var FRZero = FRReprToFR(NewFRRepr(0))

// FROne is the FR at 1.
var FROne = FRReprToFR(NewFRRepr(1))
