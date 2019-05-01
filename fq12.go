package bls

import (
	"fmt"
	"io"
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
func (f *FQ12) MulBy014Assign(c0 FQ2, c1 FQ2, c4 FQ2) {
	aa := f.c0.Copy()
	aa.MulBy01Assign(c0, c1)

	bb := f.c1.Copy()
	bb.MulBy1Assign(c4)

	c1.AddAssign(c4)
	f.c1.AddAssign(f.c0)
	f.c1.MulBy01Assign(c0, c1)
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
func (f *FQ12) DoubleAssign() {
	f.c0.DoubleAssign()
	f.c1.DoubleAssign()
}

// NegAssign negates each coefficient in an FQ12 element.
func (f *FQ12) NegAssign() {
	f.c1.NegAssign()
	f.c0.NegAssign()
}

// AddAssign adds two FQ12 elements together.
func (f *FQ12) AddAssign(other *FQ12) {
	f.c0.AddAssign(other.c0)
	f.c1.AddAssign(other.c1)
}

// IsZero returns if the FQ12 element is zero.
func (f *FQ12) IsZero() bool {
	return f.c0.IsZero() && f.c1.IsZero()
}

// SubAssign subtracts one FQ12 element from another.
func (f *FQ12) SubAssign(other *FQ12) {
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
func (f FQ12) Exp(n FQRepr) *FQ12 {
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

var frobeniusCoeffFQ12c1 = [12]FQ2{
	FQ2One,
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb}),
		FQReprToFQRaw(FQRepr{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0xecfb361b798dba3a, 0xc100ddb891865a2c, 0xec08ff1232bda8e, 0xd5c13cc6f1ca4721, 0x47222a47bf7b5c04, 0x110f184e51c5f59}),
		FQReprToFQRaw(FQRepr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8}),
		FQReprToFQRaw(FQRepr{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x51ba4ab241b6160}),
		FQReprToFQRaw(FQRepr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd}),
		FQReprToFQRaw(FQRepr{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x43f5fffffffcaaae, 0x32b7fff2ed47fffd, 0x7e83a49a2e99d69, 0xeca8f3318332bb7a, 0xef148d1ea0f4c069, 0x40ab3263eff0206}),
		FQReprToFQRaw(FQRepr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0xb2f66aad4ce5d646, 0x5842a06bfc497cec, 0xcf4895d42599d394, 0xc11b9cba40a8e8d0, 0x2e3813cbe5a0de89, 0x110eefda88847faf}),
		FQReprToFQRaw(FQRepr{0x7089552b319d465, 0xc6695f92b50a8313, 0x97e83cccd117228f, 0xa35baecab2dc29ee, 0x1ce393ea5daace4d, 0x8f2220fb0fb66eb}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x3f97d6e83d050d2, 0x18f0206554638741}),
		FQReprToFQRaw(FQRepr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x7bcfa7a25aa30fda, 0xdc17dec12a927e7c, 0x2f088dd86b4ebef1, 0xd1ca2087da74d4a7, 0x2da2596696cebc1d, 0xe2b7eedbbfd87d2}),
		FQReprToFQRaw(FQRepr{0x3e2f585da55c9ad1, 0x4294213d86c18183, 0x382844c88b623732, 0x92ad2afd19103e18, 0x1d794e4fac7cf0b9, 0xbd592fc7d825ec8}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x890dc9e4867545c3, 0x2af322533285a5d5, 0x50880866309b7e2c, 0xa20d1b8c7e881024, 0x14e4f04fe2db9068, 0x14e56d3f1564853a}),
		FQReprToFQRaw(FQRepr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
	),
	NewFQ2(
		FQReprToFQRaw(FQRepr{0x82d83cf50dbce43f, 0xa2813e53df9d018f, 0xc6f0caa53c65e181, 0x7525cf528d50fe95, 0x4a85ed50f4798a6b, 0x171da0fd6cf8eebd}),
		FQReprToFQRaw(FQRepr{0x3726c30af242c66c, 0x7c2ac1aad1b6fe70, 0xa04007fbba4b14a2, 0xef517c3266341429, 0x95ba654ed2226b, 0x2e370eccc86f7dd}),
	),
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
	c0c1 := f.c0.Copy()
	c0c1.AddAssign(f.c1)
	c0 := f.c1.Copy()
	c0.MulByNonresidueAssign()
	c0.AddAssign(f.c0)
	c0.MulAssign(c0c1)
	c0.SubAssign(ab)
	f.c1 = ab.Copy()
	f.c1.AddAssign(ab)
	ab.MulByNonresidueAssign()
	c0.SubAssign(ab)
	f.c0 = c0
}

// MulAssign multiplies two FQ12 elements together.
func (f *FQ12) MulAssign(other *FQ12) {
	aa := f.c0.Copy()
	aa.MulAssign(other.c0)
	bb := f.c1.Copy()
	bb.MulAssign(other.c1)
	o := other.c0.Copy()
	o.AddAssign(other.c1)

	f.c1.AddAssign(f.c0)
	f.c1.MulAssign(o)
	f.c1.SubAssign(aa)
	f.c1.SubAssign(bb)
	f.c0 = bb.Copy()
	f.c0.MulByNonresidueAssign()
	f.c0.AddAssign(aa)
}

// InverseAssign finds the inverse of an FQ12
func (f *FQ12) InverseAssign() bool {
	c0s := f.c0.Copy()
	c0s.SquareAssign()
	c1s := f.c1.Copy()
	c1s.SquareAssign()
	c1s.MulByNonresidueAssign()
	c0s.SubAssign(c1s)

	if !c0s.InverseAssign() {
		return false
	}

	tmp := NewFQ12(c0s.Copy(), c0s.Copy())
	tmp.c0.MulAssign(f.c0)
	tmp.c1.MulAssign(f.c1)
	tmp.c1.NegAssign()

	f.c0 = tmp.c0
	f.c1 = tmp.c1

	return true
}
