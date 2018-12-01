package bls

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits"
)

// FRRepr represents a uint256.
type FRRepr [4]uint64

// IsOdd checks if the FRRepr is odd.
func (f FRRepr) IsOdd() bool {
	return f[0]&1 == 1
}

// IsEven checks if the FRRepr is even.
func (f FRRepr) IsEven() bool {
	return f[0]&1 == 0
}

// IsZero checks if the FRRepr is zero.
func (f FRRepr) IsZero() bool {
	for _, f := range f {
		if f != 0 {
			return false
		}
	}
	return true
}

// NewFRRepr creates a new number given a uint64.
func NewFRRepr(n uint64) *FRRepr {
	return &FRRepr{n, 0, 0, 0}
}

// Rsh shifts the FRRepr right by a certain number of bits.
func (f *FRRepr) Rsh(n uint) {
	if n >= 64*4 {
		out := NewFRRepr(0)
		*f = *out
		return
	}

	for n >= 64 {
		t := uint64(0)
		for i := 3; i >= 0; i-- {
			t, f[i] = f[i], t
		}
		n -= 64
	}

	if n > 0 {
		t := uint64(0)
		for i := 3; i >= 0; i-- {
			t2 := f[i] << (64 - n)
			f[i] >>= n
			f[i] |= t
			t = t2
		}
	}
}

// Div2 divides the FRRepr by 2.
func (f *FRRepr) Div2() {
	t := uint64(0)
	for i := 3; i >= 0; i-- {
		t2 := f[i] << 63
		f[i] >>= 1
		f[i] |= t
		t = t2
	}
}

// Mul2 multiplies the FRRepr by 2.
func (f *FRRepr) Mul2() {
	last := uint64(0)
	for i := 0; i < 3; i++ {
		tmp := f[i] >> 63
		f[i] <<= 1
		f[i] |= last
		last = tmp
	}
}

// Lsh shifts the FRRepr left by a certain number of bits.
func (f *FRRepr) Lsh(n uint) {
	if n >= 64*4 {
		out := NewFRRepr(0)
		*f = *out
		return
	}

	for n >= 64 {
		t := uint64(0)
		for i := 0; i < 3; i++ {
			t, f[i] = f[i], t
		}
	}

	if n > 0 {
		t := uint64(0)
		for i := 0; i < 3; i++ {
			t2 := f[i] >> (64 - n)
			f[i] <<= n
			f[i] |= t
			t = t2
		}
	}
}

// AddNoCarry adds two FRReprs to another and does not handle
// carry.
func (f *FRRepr) AddNoCarry(g *FRRepr) {
	carry := uint64(0)
	for i := 0; i < 6; i++ {
		f[i] = AddWithCarry(f[i], g[i], &carry)
	}
}

// SubNoBorrow subtracts two FRReprs from another and does not handle
// borrow.
func (f *FRRepr) SubNoBorrow(g *FRRepr) {
	borrow := uint64(0)
	for i := 0; i < 3; i++ {
		f[i] = SubWithBorrow(f[i], g[i], &borrow)
	}
}

// Equals checks if two FRRepr's are equal.
func (f *FRRepr) Equals(g *FRRepr) bool {
	return f[0] == g[0] && f[1] == g[1] && f[2] == g[2] && f[3] == g[3]
}

// Cmp compares two FRRepr's
func (f *FRRepr) Cmp(g *FRRepr) int {
	for i := 0; i < 3; i++ {
		if f[i] > g[i] {
			return 1
		} else if f[i] < g[i] {
			return -1
		}
	}
	return 0
}

// Copy copies a FRRepr to a new instance and returns it.
func (f *FRRepr) Copy() *FRRepr {
	var newBytes [4]uint64
	copy(newBytes[:], f[:])
	out := FRRepr(newBytes)
	return &out
}

// ToString converts the FRRepr to a string.
func (f FRRepr) ToString() string {
	return fmt.Sprintf("%x%x%x%x", f[3], f[2], f[1], f[0])
}

// BitLen counts the number of bits the number is.
func (f FRRepr) BitLen() uint {
	ret := uint(4 * 64)
	for i := 3; i >= 0; i-- {
		leading := uint(bits.LeadingZeros64(f[i]))
		ret -= leading
		if leading != 64 {
			break
		}
	}

	return ret
}

// FRReprFromBytes gets a new FRRepr from big-endian bytes.
func FRReprFromBytes(b []byte) (*FRRepr, error) {
	return FRReprFromBigInt(new(big.Int).SetBytes(b))
}

// Bit checks if a bit is set (little-endian)
func (f FRRepr) Bit(n uint) bool {
	return f[n/8]&(1<<(n%8)) != 0
}

// FRReprFromString creates a FRRepr from a string.
func FRReprFromString(s string, b uint) (*FRRepr, error) {
	out, valid := new(big.Int).SetString(s, int(b))
	if !valid {
		return nil, errors.New("FRRepr not valid")
	}
	return FRReprFromBigInt(out)
}

// ToBig gets the big.Int representation of the FRRepr.
func (f FRRepr) ToBig() *big.Int {
	out := big.NewInt(0)
	for i := 0; i < 3; i++ {
		out.Add(out, new(big.Int).SetUint64(f[i]))
		out.Lsh(out, 64)
	}
	return out
}

// FRReprFromBigInt create a FRRepr from a big.Int.
func FRReprFromBigInt(out *big.Int) (*FRRepr, error) {
	if out.BitLen() > 256 || out.Sign() == -1 {
		return nil, errors.New("invalid input string")
	}

	newf := NewFRRepr(0)
	for out.BitLen() > 0 {
		newf.AddNoCarry(NewFRRepr(out.Uint64()))
		newf.Lsh(64)
		out.Rsh(out, 64)
	}

	return newf, nil
}

// ToFQ converts an FRRepr to an FQ.
func (f *FRRepr) ToFQ() *FQRepr {
	newf := NewFQRepr(f[0])
	newf[1] = f[1]
	newf[2] = f[2]
	newf[3] = f[3]
	return newf
}
