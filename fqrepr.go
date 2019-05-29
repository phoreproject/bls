package bls

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/bits"
)

//go:generate go run asm/asm.go -out primitivefuncs_amd64.s

// FQRepr represents a uint384. The least significant bits are first.
type FQRepr [6]uint64

// IsOdd checks if the FQRepr is odd.
func (f FQRepr) IsOdd() bool {
	return f[0]&1 == 1
}

// IsEven checks if the FQRepr is even.
func (f FQRepr) IsEven() bool {
	return f[0]&1 == 0
}

// IsZero checks if the FQRepr is zero.
func (f FQRepr) IsZero() bool {
	for _, f := range f {
		if f != 0 {
			return false
		}
	}
	return true
}

// NewFQRepr creates a new number given a uint64.
func NewFQRepr(n uint64) FQRepr {
	return FQRepr{n, 0, 0, 0, 0, 0}
}

// Rsh shifts the FQRepr right by a certain number of bits.
func (f *FQRepr) Rsh(n uint) {
	if n >= 64*6 {
		f[0] = 0
		f[1] = 0
		f[2] = 0
		f[3] = 0
		f[4] = 0
		f[5] = 0
		return
	}

	for n >= 64 {
		t := uint64(0)
		for i := 5; i >= 0; i-- {
			t, f[i] = f[i], t
		}
		n -= 64
	}

	if n > 0 {
		t := uint64(0)
		for i := 5; i >= 0; i-- {
			t2 := f[i] << (64 - n)
			f[i] >>= n
			f[i] |= t
			t = t2
		}
	}
}

// Div2 divides the FQRepr by 2.
func (f *FQRepr) Div2() {
	t := uint64(0)
	for i := 5; i >= 0; i-- {
		t2 := f[i] << 63
		f[i] >>= 1
		f[i] |= t
		t = t2
	}
}

// Mul2 multiplies the FQRepr by 2.
func (f *FQRepr) Mul2() {
	last := uint64(0)
	for i := 0; i < 6; i++ {
		tmp := f[i] >> 63
		f[i] <<= 1
		f[i] |= last
		last = tmp
	}
}

// Lsh shifts the FQRepr left by a certain number of bits.
func (f *FQRepr) Lsh(n uint) {
	if n >= 64*6 {
		f[0] = 0
		f[1] = 0
		f[2] = 0
		f[3] = 0
		f[4] = 0
		f[5] = 0
		return
	}

	for n >= 64 {
		t := uint64(0)
		for i := 0; i < 6; i++ {
			t, f[i] = f[i], t
		}
		n -= 64
	}

	if n > 0 {
		t := uint64(0)
		for i := 0; i < 6; i++ {
			t2 := f[i] >> (64 - n)
			f[i] <<= n
			f[i] |= t
			t = t2
		}
	}
}

// AddNoCarry adds two FQReprs to another and does not handle
// carry.
func (f *FQRepr) AddNoCarry(g FQRepr) {
	*f = AddNoCarry(*f, g)
}

// SubNoBorrow subtracts two FQReprs from another and does not handle
// borrow.
func (f *FQRepr) SubNoBorrow(g FQRepr) {
	*f = SubNoBorrow(*f, g)
}

// Equals checks if two FQRepr's are equal.
func (f *FQRepr) Equals(g FQRepr) bool {
	return f[0] == g[0] && f[1] == g[1] && f[2] == g[2] && f[3] == g[3] && f[4] == g[4] && f[5] == g[5]
}

// Cmp compares two FQRepr's
func (f *FQRepr) Cmp(g FQRepr) int {
	for i := 5; i >= 0; i-- {
		if f[i] == g[i] {
			continue
		}
		if f[i] > g[i] {
			return 1
		} else if f[i] < g[i] {
			return -1
		}
	}
	return 0
}

// Copy copies a FQRepr to a new instance and returns it.
func (f *FQRepr) Copy() FQRepr {
	return *f
}

// ToString converts the FQRepr to a string.
func (f FQRepr) String() string {
	return fmt.Sprintf("%016x%016x%016x%016x%016x%016x", f[5], f[4], f[3], f[2], f[1], f[0])
}

// BitLen counts the number of bits the number is.
func (f FQRepr) BitLen() uint {
	ret := uint(6 * 64)
	for i := 5; i >= 0; i-- {
		leading := uint(bits.LeadingZeros64(f[i]))
		ret -= leading
		if leading != 64 {
			break
		}
	}

	return ret
}

// FQReprFromBytes gets a new FQRepr from big-endian bytes.
func FQReprFromBytes(b [48]byte) FQRepr {
	m0 := binary.BigEndian.Uint64(b[0:8])
	m1 := binary.BigEndian.Uint64(b[8:16])
	m2 := binary.BigEndian.Uint64(b[16:24])
	m3 := binary.BigEndian.Uint64(b[24:32])
	m4 := binary.BigEndian.Uint64(b[32:40])
	m5 := binary.BigEndian.Uint64(b[40:48])
	return FQRepr{m5, m4, m3, m2, m1, m0}
}

// Bytes gets the bytes used for an FQRepr.
func (f FQRepr) Bytes() [48]byte {
	var out [48]byte
	binary.BigEndian.PutUint64(out[0:8], f[5])
	binary.BigEndian.PutUint64(out[8:16], f[4])
	binary.BigEndian.PutUint64(out[16:24], f[3])
	binary.BigEndian.PutUint64(out[24:32], f[2])
	binary.BigEndian.PutUint64(out[32:40], f[1])
	binary.BigEndian.PutUint64(out[40:48], f[0])
	return out
}

// Bit checks if a bit is set (little-endian)
func (f FQRepr) Bit(n uint) bool {
	return f[n/64]&(1<<(n%64)) != 0
}

// FQReprFromString creates a FQRepr from a string.
func FQReprFromString(s string, b uint) (FQRepr, error) {
	out, valid := new(big.Int).SetString(s, int(b))
	if !valid {
		return FQRepr{}, errors.New("FQRepr not valid")
	}
	return FQReprFromBigInt(out)
}

func fqReprFromHexUnchecked(s string) FQRepr {
	out, _ := new(big.Int).SetString(s, 16)
	return fqReprFromBigIntUnchecked(out)
}

func fqReprFromStringUnchecked(s string, b uint) FQRepr {
	out, _ := new(big.Int).SetString(s, int(b))
	return fqReprFromBigIntUnchecked(out)
}

// ToBig gets the big.Int representation of the FQRepr.
func (f FQRepr) ToBig() *big.Int {
	out := big.NewInt(0)
	for i := 5; i >= 0; i-- {
		out.Add(out, new(big.Int).SetUint64(f[i]))
		if i != 0 {
			out.Lsh(out, 64)
		}
	}
	return out
}

var bigIntZero = big.NewInt(0)
var oneLsh64MinusOne = new(big.Int).SetUint64(0xffffffffffffffff)

// FQReprFromBigInt create a FQRepr from a big.Int.
func FQReprFromBigInt(n *big.Int) (FQRepr, error) {
	if n.BitLen() > 384 || n.Sign() == -1 {
		return FQRepr{}, errors.New("invalid input string")
	}

	out := new(big.Int).Set(n)

	newf := NewFQRepr(0)
	i := 0
	for out.Cmp(bigIntZero) != 0 {
		o := new(big.Int).And(out, oneLsh64MinusOne)
		newf[i] = o.Uint64()
		i++
		out.Rsh(out, 64)
	}

	return newf, nil
}

// FQReprFromBigInt create a FQRepr from a big.Int.
func fqReprFromBigIntUnchecked(n *big.Int) FQRepr {
	out := new(big.Int).Set(n)

	newf := NewFQRepr(0)
	i := 0
	for out.Cmp(bigIntZero) != 0 {
		o := new(big.Int).And(out, oneLsh64MinusOne)
		newf[i] = o.Uint64()
		i++
		out.Rsh(out, 64)
	}

	return newf
}
