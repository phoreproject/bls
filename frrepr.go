package bls

import (
	"encoding/binary"
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
	for i := 0; i < 4; i++ {
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
		for i := 0; i < 4; i++ {
			t, f[i] = f[i], t
		}
		n -= 64
	}

	if n > 0 {
		t := uint64(0)
		for i := 0; i < 4; i++ {
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
	for i := 0; i < 4; i++ {
		f[i], carry = AddWithCarry(f[i], g[i], carry)
	}
}

// SubNoBorrow subtracts two FRReprs from another and does not handle
// borrow.
func (f *FRRepr) SubNoBorrow(g *FRRepr) {
	borrow := uint64(0)
	for i := 0; i < 4; i++ {
		f[i], borrow = SubWithBorrow(f[i], g[i], borrow)
	}
}

// Equals checks if two FRRepr's are equal.
func (f *FRRepr) Equals(g *FRRepr) bool {
	return f[0] == g[0] && f[1] == g[1] && f[2] == g[2] && f[3] == g[3]
}

// Cmp compares two FRRepr's
func (f *FRRepr) Cmp(g *FRRepr) int {
	for i := 3; i >= 0; i-- {
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
	newf := FRRepr(newBytes)
	return &newf
}

// ToString converts the FRRepr to a string.
func (f FRRepr) String() string {
	return fmt.Sprintf("%016x%016x%016x%016x", f[3], f[2], f[1], f[0])
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
func FRReprFromBytes(b [32]byte) *FRRepr {
	m0 := binary.BigEndian.Uint64(b[0:8])
	m1 := binary.BigEndian.Uint64(b[8:16])
	m2 := binary.BigEndian.Uint64(b[16:24])
	m3 := binary.BigEndian.Uint64(b[24:32])
	return &FRRepr{m3, m2, m1, m0}
}

// Bytes gets the bytes used for an FRRepr.
func (f FRRepr) Bytes() [32]byte {
	var out [32]byte
	binary.BigEndian.PutUint64(out[0:8], f[3])
	binary.BigEndian.PutUint64(out[8:16], f[2])
	binary.BigEndian.PutUint64(out[16:24], f[1])
	binary.BigEndian.PutUint64(out[24:32], f[0])
	return out
}

// Bit checks if a bit is set (little-endian)
func (f FRRepr) Bit(n uint) bool {
	return f[n/64]&(1<<(n%64)) != 0
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
	for i := 3; i >= 0; i-- {
		out.Add(out, new(big.Int).SetUint64(f[i]))
		if i != 0 {
			out.Lsh(out, 64)
		}
	}
	return out
}

// FRReprFromBigInt create a FRRepr from a big.Int.
func FRReprFromBigInt(n *big.Int) (*FRRepr, error) {
	if n.BitLen() > 256 || n.Sign() == -1 {
		return nil, errors.New("invalid input string")
	}

	out := new(big.Int).Set(n)

	newf := NewFRRepr(0)
	i := 0
	for out.Cmp(bigIntZero) != 0 {
		o := new(big.Int).And(out, oneLsh64MinusOne)
		newf[i] = o.Uint64()
		i++
		out.Rsh(out, 64)
	}

	return newf, nil
}

// ToFQ converts an FRRepr to an FQ.
func (f *FRRepr) ToFQ() FQRepr {
	newf := NewFQRepr(f[0])
	newf[1] = f[1]
	newf[2] = f[2]
	newf[3] = f[3]
	return newf
}
