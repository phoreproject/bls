package bls

import (
	"crypto/rand"
	"fmt"
	"hash"
	"io"
)

// FQ is an element in a field.
type FQ struct {
	n *FQRepr
}

var bigZero = NewFQRepr(0)
var bigOne = NewFQRepr(1)
var bigTwo = NewFQRepr(2)

// FQZero is the zero FQ element
var FQZero = FQReprToFQRaw(bigZero)

// FQOne is the one FQ element
var FQOne = FQReprToFQ(bigOne)
var bigTwoFQ = FQReprToFQ(bigTwo)

// QFieldModulus is the modulus of the field.
var QFieldModulus, _ = FQReprFromString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", 10)

// FQR is 2**384 % Q used for moving numbers into Montgomery form.
var FQR = &FQRepr{0x760900000002fffd, 0xebf4000bc40c0002, 0x5f48985753c758ba, 0x77ce585370525745, 0x5c071a97a256ec6d, 0x15f65ec3fa80e493}

// FQR2 is R^2 % Q.
var FQR2, _ = FQReprFromString("2708263910654730174793787626328176511836455197166317677006154293982164122222515399004018013397331347120527951271750", 10)

// Copy creates a copy of the field element.
func (f *FQ) Copy() *FQ {
	return &FQ{f.n.Copy()}
}

// IsValid checks if the element is valid.
func (f *FQ) IsValid() bool {
	return f.n[5]&0xf000000000000000 == 0 || f.n.Cmp(QFieldModulus) < 0
}

func (f *FQ) reduceAssign() {
	if !f.IsValid() {
		f.n.SubNoBorrow(QFieldModulus)
	}
}

// FQReprToFQ gets a pointer to a FQ given a pointer
// to an FQRepr
func FQReprToFQ(o *FQRepr) *FQ {
	r := &FQ{n: o.Copy()}
	if r.IsValid() {
		r.MulAssign(&FQ{FQR2})
		return r
	}
	return nil
}

// FQReprToFQRaw gets a pointer to a FQ without converting
// to montgomery form.
func FQReprToFQRaw(o *FQRepr) *FQ {
	return &FQ{n: o}
}

// AddAssign multiplies a field element by this one.
func (f *FQ) AddAssign(other *FQ) {
	f.n.AddNoCarry(other.n)
	f.reduceAssign()
}

const montInvFQ = uint64(0x89f3fffcfffcfffd)

func (f *FQ) montReduce(r0 uint64, r1 uint64, r2 uint64, r3 uint64, r4 uint64, r5 uint64, r6 uint64, r7 uint64, r8 uint64, r9 uint64, r10 uint64, r11 uint64) {
	k := r0 * montInvFQ
	_, carry := MACWithCarry(r0, k, QFieldModulus[0], 0)
	r1, carry = MACWithCarry(r1, k, QFieldModulus[1], carry)
	r2, carry = MACWithCarry(r2, k, QFieldModulus[2], carry)
	r3, carry = MACWithCarry(r3, k, QFieldModulus[3], carry)
	r4, carry = MACWithCarry(r4, k, QFieldModulus[4], carry)
	r5, carry = MACWithCarry(r5, k, QFieldModulus[5], carry)
	r6, carry = AddWithCarry(r6, 0, carry)
	carry2 := carry
	k = r1 * montInvFQ
	_, carry = MACWithCarry(r1, k, QFieldModulus[0], 0)
	r2, carry = MACWithCarry(r2, k, QFieldModulus[1], carry)
	r3, carry = MACWithCarry(r3, k, QFieldModulus[2], carry)
	r4, carry = MACWithCarry(r4, k, QFieldModulus[3], carry)
	r5, carry = MACWithCarry(r5, k, QFieldModulus[4], carry)
	r6, carry = MACWithCarry(r6, k, QFieldModulus[5], carry)
	r7, carry = AddWithCarry(r7, carry2, carry)
	carry2 = carry
	k = r2 * montInvFQ
	_, carry = MACWithCarry(r2, k, QFieldModulus[0], 0)
	r3, carry = MACWithCarry(r3, k, QFieldModulus[1], carry)
	r4, carry = MACWithCarry(r4, k, QFieldModulus[2], carry)
	r5, carry = MACWithCarry(r5, k, QFieldModulus[3], carry)
	r6, carry = MACWithCarry(r6, k, QFieldModulus[4], carry)
	r7, carry = MACWithCarry(r7, k, QFieldModulus[5], carry)
	r8, carry = AddWithCarry(r8, carry2, carry)
	carry2 = carry
	k = r3 * montInvFQ
	_, carry = MACWithCarry(r3, k, QFieldModulus[0], 0)
	r4, carry = MACWithCarry(r4, k, QFieldModulus[1], carry)
	r5, carry = MACWithCarry(r5, k, QFieldModulus[2], carry)
	r6, carry = MACWithCarry(r6, k, QFieldModulus[3], carry)
	r7, carry = MACWithCarry(r7, k, QFieldModulus[4], carry)
	r8, carry = MACWithCarry(r8, k, QFieldModulus[5], carry)
	r9, carry = AddWithCarry(r9, carry2, carry)
	carry2 = carry
	k = r4 * montInvFQ
	_, carry = MACWithCarry(r4, k, QFieldModulus[0], 0)
	r5, carry = MACWithCarry(r5, k, QFieldModulus[1], carry)
	r6, carry = MACWithCarry(r6, k, QFieldModulus[2], carry)
	r7, carry = MACWithCarry(r7, k, QFieldModulus[3], carry)
	r8, carry = MACWithCarry(r8, k, QFieldModulus[4], carry)
	r9, carry = MACWithCarry(r9, k, QFieldModulus[5], carry)
	r10, carry = AddWithCarry(r10, carry2, carry)
	carry2 = carry
	k = r5 * montInvFQ
	_, carry = MACWithCarry(r5, k, QFieldModulus[0], 0)
	r6, carry = MACWithCarry(r6, k, QFieldModulus[1], carry)
	r7, carry = MACWithCarry(r7, k, QFieldModulus[2], carry)
	r8, carry = MACWithCarry(r8, k, QFieldModulus[3], carry)
	r9, carry = MACWithCarry(r9, k, QFieldModulus[4], carry)
	r10, carry = MACWithCarry(r10, k, QFieldModulus[5], carry)
	r11, carry = AddWithCarry(r11, carry2, carry)
	f.n[0] = r6
	f.n[1] = r7
	f.n[2] = r8
	f.n[3] = r9
	f.n[4] = r10
	f.n[5] = r11
	f.reduceAssign()
}

// MulAssign multiplies a field element by this one.
func (f FQ) MulAssign(other *FQ) {
	r0, carry := MACWithCarry(0, f.n[0], other.n[0], 0)
	r1, carry := MACWithCarry(0, f.n[0], other.n[1], carry)
	r2, carry := MACWithCarry(0, f.n[0], other.n[2], carry)
	r3, carry := MACWithCarry(0, f.n[0], other.n[3], carry)
	r4, carry := MACWithCarry(0, f.n[0], other.n[4], carry)
	r5, carry := MACWithCarry(0, f.n[0], other.n[5], carry)
	r6 := carry
	r1, carry = MACWithCarry(r1, f.n[1], other.n[0], 0)
	r2, carry = MACWithCarry(r2, f.n[1], other.n[1], carry)
	r3, carry = MACWithCarry(r3, f.n[1], other.n[2], carry)
	r4, carry = MACWithCarry(r4, f.n[1], other.n[3], carry)
	r5, carry = MACWithCarry(r5, f.n[1], other.n[4], carry)
	r6, carry = MACWithCarry(r6, f.n[1], other.n[5], carry)
	r7 := carry
	r2, carry = MACWithCarry(r2, f.n[2], other.n[0], 0)
	r3, carry = MACWithCarry(r3, f.n[2], other.n[1], carry)
	r4, carry = MACWithCarry(r4, f.n[2], other.n[2], carry)
	r5, carry = MACWithCarry(r5, f.n[2], other.n[3], carry)
	r6, carry = MACWithCarry(r6, f.n[2], other.n[4], carry)
	r7, carry = MACWithCarry(r7, f.n[2], other.n[5], carry)
	r8 := carry
	r3, carry = MACWithCarry(r3, f.n[3], other.n[0], 0)
	r4, carry = MACWithCarry(r4, f.n[3], other.n[1], carry)
	r5, carry = MACWithCarry(r5, f.n[3], other.n[2], carry)
	r6, carry = MACWithCarry(r6, f.n[3], other.n[3], carry)
	r7, carry = MACWithCarry(r7, f.n[3], other.n[4], carry)
	r8, carry = MACWithCarry(r8, f.n[3], other.n[5], carry)
	r9 := carry
	r4, carry = MACWithCarry(r4, f.n[4], other.n[0], 0)
	r5, carry = MACWithCarry(r5, f.n[4], other.n[1], carry)
	r6, carry = MACWithCarry(r6, f.n[4], other.n[2], carry)
	r7, carry = MACWithCarry(r7, f.n[4], other.n[3], carry)
	r8, carry = MACWithCarry(r8, f.n[4], other.n[4], carry)
	r9, carry = MACWithCarry(r9, f.n[4], other.n[5], carry)
	r10 := carry
	r5, carry = MACWithCarry(r5, f.n[5], other.n[0], 0)
	r6, carry = MACWithCarry(r6, f.n[5], other.n[1], carry)
	r7, carry = MACWithCarry(r7, f.n[5], other.n[2], carry)
	r8, carry = MACWithCarry(r8, f.n[5], other.n[3], carry)
	r9, carry = MACWithCarry(r9, f.n[5], other.n[4], carry)
	r10, carry = MACWithCarry(r10, f.n[5], other.n[5], carry)
	r11 := carry
	f.montReduce(r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11)
}

// SubAssign subtracts a field element from this one.
func (f *FQ) SubAssign(other *FQ) {
	if other.n.Cmp(f.n) > 0 {
		f.n.AddNoCarry(QFieldModulus)
	}
	f.n.SubNoBorrow(other.n)
}

// divAssign is a slow divide.
func (f *FQ) divAssign(other *FQ) {
	a := f.n.ToBig()
	a.Div(a, other.n.ToBig())
	fOut, _ := FQReprFromBigInt(a)
	fqOut := FQReprToFQ(fOut)
	*f = *fqOut
}

// Exp raises the element to a specific power.
func (f *FQ) Exp(n *FQRepr) *FQ {
	iter := NewBitIterator(n[:])
	res := FQOne.Copy()
	foundOne := false
	next, done := iter.Next()
	for !done {
		if foundOne {
			res.SquareAssign()
		} else {
			foundOne = next
		}
		if next {
			res.MulAssign(f)
		}
		next, done = iter.Next()
	}
	return res
}

// Equals checks equality of two field elements.
func (f FQ) Equals(other *FQ) bool {
	return f.n.Equals(other.n)
}

// NegAssign gets the negative value of the field element mod QFieldModulus.
func (f *FQ) NegAssign() {
	if !f.IsZero() {
		tmp := QFieldModulus.Copy()
		tmp.SubNoBorrow(f.n)
		f.n = tmp
	}
}

func (f FQ) String() string {
	return fmt.Sprintf("Fq(0x%s)", f.ToRepr().String())
}

// Cmp compares this field element to another.
func (f FQ) Cmp(other *FQ) int {
	return f.ToRepr().Cmp(other.ToRepr())
}

// DoubleAssign doubles the element
func (f *FQ) DoubleAssign() {
	f.n.Mul2()
	f.reduceAssign()
}

// IsZero checks if the field element is zero.
func (f FQ) IsZero() bool {
	return f.n.Cmp(bigZero) == 0
}

// SquareAssign squares a field element.
func (f *FQ) SquareAssign() {
	r1, carry := MACWithCarry(0, f.n[0], f.n[1], 0)
	r2, carry := MACWithCarry(0, f.n[0], f.n[2], carry)
	r3, carry := MACWithCarry(0, f.n[0], f.n[3], carry)
	r4, carry := MACWithCarry(0, f.n[0], f.n[4], carry)
	r5, carry := MACWithCarry(0, f.n[0], f.n[5], carry)
	r6 := carry
	r3, carry = MACWithCarry(r3, f.n[1], f.n[2], 0)
	r4, carry = MACWithCarry(r4, f.n[1], f.n[3], carry)
	r5, carry = MACWithCarry(r5, f.n[1], f.n[4], carry)
	r6, carry = MACWithCarry(r6, f.n[1], f.n[5], carry)
	r7 := carry
	r5, carry = MACWithCarry(r5, f.n[2], f.n[3], 0)
	r6, carry = MACWithCarry(r6, f.n[2], f.n[4], carry)
	r7, carry = MACWithCarry(r7, f.n[2], f.n[5], carry)
	r8 := carry
	r7, carry = MACWithCarry(r7, f.n[3], f.n[4], 0)
	r8, carry = MACWithCarry(r8, f.n[3], f.n[5], carry)
	r9 := carry
	r9, carry = MACWithCarry(r9, f.n[4], f.n[5], 0)
	r10 := carry
	r11 := r10 >> 63
	r10 = (r10 << 1) | (r9 >> 63)
	r9 = (r9 << 1) | (r8 >> 63)
	r8 = (r8 << 1) | (r7 >> 63)
	r7 = (r7 << 1) | (r6 >> 63)
	r6 = (r6 << 1) | (r5 >> 63)
	r5 = (r5 << 1) | (r4 >> 63)
	r4 = (r4 << 1) | (r3 >> 63)
	r3 = (r3 << 1) | (r2 >> 63)
	r2 = (r2 << 1) | (r1 >> 63)
	r1 = r1 << 1

	carry = 0
	r0, carry := MACWithCarry(0, f.n[0], f.n[0], carry)
	r1, carry = AddWithCarry(r1, 0, carry)
	r2, carry = MACWithCarry(r2, f.n[1], f.n[1], carry)
	r3, carry = AddWithCarry(r3, 0, carry)
	r4, carry = MACWithCarry(r4, f.n[2], f.n[2], carry)
	r5, carry = AddWithCarry(r5, 0, carry)
	r6, carry = MACWithCarry(r6, f.n[3], f.n[3], carry)
	r7, carry = AddWithCarry(r7, 0, carry)
	r8, carry = MACWithCarry(r8, f.n[4], f.n[4], carry)
	r9, carry = AddWithCarry(r9, 0, carry)
	r10, carry = MACWithCarry(r10, f.n[5], f.n[5], carry)
	r11, carry = AddWithCarry(r11, 0, carry)
	f.montReduce(r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11)
}

var negativeOneFQ = FQReprToFQ(negativeOne)

// Sqrt calculates the square root of the field element.
func (f FQ) Sqrt() *FQ {
	// Shank's algorithm for q mod 4 = 3
	// https://eprint.iacr.org/2012/685.pdf (page 9, algorithm 2)

	a1 := f.Exp(qMinus3Over4)
	a0 := a1.Copy()
	a0.SquareAssign()
	a0.MulAssign(&f)

	if a0.Equals(negativeOneFQ) {
		return nil
	}
	a1.MulAssign(&f)
	return a1
}

func isEven(b *FQRepr) bool {
	return b.IsEven()
}

// Inverse finds the inverse of the field element.
func (f FQ) Inverse() *FQ {
	if f.IsZero() {
		return nil
	}
	u := f.n.Copy()
	v := QFieldModulus.Copy()
	b := FQReprToFQRaw(FQR2.Copy())
	c := FQZero.Copy()

	for u.Cmp(bigOne) != 0 && v.Cmp(bigOne) != 0 {
		for isEven(u) {
			u.Div2()
			if isEven(b.n) {
				b.n.Div2()
			} else {
				b.n.AddNoCarry(QFieldModulus)
				b.n.Div2()
			}
		}

		for isEven(v) {
			v.Div2()
			if isEven(c.n) {
				c.n.Div2()
			} else {
				c.n.AddNoCarry(QFieldModulus)
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
	if u.Cmp(bigOne) == 0 {
		return b
	}
	return c
}

// Parity checks if the point is greater than the point negated.
func (f FQ) Parity() bool {
	neg := f.Copy()
	neg.NegAssign()
	return f.Cmp(neg) > 0
}

// MulBits multiplies the number by a big number.
func (f FQ) MulBits(b *FQRepr) *FQ {
	res := FQZero.Copy()
	for i := uint(0); i < b.BitLen(); i++ {
		res.DoubleAssign()
		if b.Bit(i) {
			res.AddAssign(&f)
		}
	}
	return res
}

// MulBytes multiplies the number by some bytes.
func (f FQ) MulBytes(b []byte) *FQ {
	res := FQZero.Copy()
	for i := uint(0); i < uint(len(b)*8); i++ {
		res.DoubleAssign()
		if b[i/8]&(1<<(i%8)) != 0 {
			res.AddAssign(&f)
		}
	}
	return res
}

// HashFQ calculates a new FQ2 value based on a hash.
func HashFQ(hasher hash.Hash) *FQ {
	digest := hasher.Sum(nil)
	return FQOne.MulBytes(digest)
}

var qMinus1Over2 = &FQRepr{0xdcff7fffffffd555, 0xf55ffff58a9ffff, 0xb39869507b587b12, 0xb23ba5c279c2895f, 0x258dd3db21a5d66b, 0xd0088f51cbff34d}

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

// ToRepr gets the 256-bit representation of the field element.
func (f *FQ) ToRepr() *FQRepr {
	out := f.Copy()
	out.montReduce(
		f.n[0],
		f.n[1],
		f.n[2],
		f.n[3],
		f.n[4],
		f.n[5],
		0,
		0,
		0,
		0,
		0,
		0,
	)
	return out.n
}

// RandFQ generates a random FQ element.
func RandFQ(reader io.Reader) (*FQ, error) {
	r, err := rand.Int(reader, QFieldModulus.ToBig())
	if err != nil {
		return nil, err
	}
	b, _ := FQReprFromBigInt(r)
	return FQReprToFQ(b), nil
}
