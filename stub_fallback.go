// +build !amd64 !gc

package bls

import (
	"math/big"
	"math/bits"
)

// MultiplyFQRepr multiplies two FQRepr values together.
func MultiplyFQRepr(a, b [6]uint64) (hi [6]uint64, lo [6]uint64) {
	carry := uint64(0)
	lo[0], carry = MACWithCarry(0, a[0], b[0], 0)
	lo[1], carry = MACWithCarry(0, a[0], b[1], carry)
	lo[2], carry = MACWithCarry(0, a[0], b[2], carry)
	lo[3], carry = MACWithCarry(0, a[0], b[3], carry)
	lo[4], carry = MACWithCarry(0, a[0], b[4], carry)
	lo[5], carry = MACWithCarry(0, a[0], b[5], carry)
	hi[0] = carry
	lo[1], carry = MACWithCarry(lo[1], a[1], b[0], 0)
	lo[2], carry = MACWithCarry(lo[2], a[1], b[1], carry)
	lo[3], carry = MACWithCarry(lo[3], a[1], b[2], carry)
	lo[4], carry = MACWithCarry(lo[4], a[1], b[3], carry)
	lo[5], carry = MACWithCarry(lo[5], a[1], b[4], carry)
	hi[0], carry = MACWithCarry(hi[0], a[1], b[5], carry)
	hi[1] = carry
	lo[2], carry = MACWithCarry(lo[2], a[2], b[0], 0)
	lo[3], carry = MACWithCarry(lo[3], a[2], b[1], carry)
	lo[4], carry = MACWithCarry(lo[4], a[2], b[2], carry)
	lo[5], carry = MACWithCarry(lo[5], a[2], b[3], carry)
	hi[0], carry = MACWithCarry(hi[0], a[2], b[4], carry)
	hi[1], carry = MACWithCarry(hi[1], a[2], b[5], carry)
	hi[2] = carry
	lo[3], carry = MACWithCarry(lo[3], a[3], b[0], 0)
	lo[4], carry = MACWithCarry(lo[4], a[3], b[1], carry)
	lo[5], carry = MACWithCarry(lo[5], a[3], b[2], carry)
	hi[0], carry = MACWithCarry(hi[0], a[3], b[3], carry)
	hi[1], carry = MACWithCarry(hi[1], a[3], b[4], carry)
	hi[2], carry = MACWithCarry(hi[2], a[3], b[5], carry)
	hi[3] = carry
	lo[4], carry = MACWithCarry(lo[4], a[4], b[0], 0)
	lo[5], carry = MACWithCarry(lo[5], a[4], b[1], carry)
	hi[0], carry = MACWithCarry(hi[0], a[4], b[2], carry)
	hi[1], carry = MACWithCarry(hi[1], a[4], b[3], carry)
	hi[2], carry = MACWithCarry(hi[2], a[4], b[4], carry)
	hi[3], carry = MACWithCarry(hi[3], a[4], b[5], carry)
	hi[4] = carry
	lo[5], carry = MACWithCarry(lo[5], a[5], b[0], 0)
	hi[0], carry = MACWithCarry(hi[0], a[5], b[1], carry)
	hi[1], carry = MACWithCarry(hi[1], a[5], b[2], carry)
	hi[2], carry = MACWithCarry(hi[2], a[5], b[3], carry)
	hi[3], carry = MACWithCarry(hi[3], a[5], b[4], carry)
	hi[4], carry = MACWithCarry(hi[4], a[5], b[5], carry)
	hi[5] = carry

	return hi, lo
}

const montInvFQ = uint64(0x89f3fffcfffcfffd)

func MontReduce(hi, lo [6]uint64) [6]uint64 {
	k := lo[0] * montInvFQ
	_, carry := MACWithCarry(lo[0], k, QFieldModulus[0], 0)
	lo[1], carry = MACWithCarry(lo[1], k, QFieldModulus[1], carry)
	lo[2], carry = MACWithCarry(lo[2], k, QFieldModulus[2], carry)
	lo[3], carry = MACWithCarry(lo[3], k, QFieldModulus[3], carry)
	lo[4], carry = MACWithCarry(lo[4], k, QFieldModulus[4], carry)
	lo[5], carry = MACWithCarry(lo[5], k, QFieldModulus[5], carry)
	hi[0], carry = AddWithCarry(hi[0], 0, carry)
	carry2 := carry
	k = lo[1] * montInvFQ
	_, carry = MACWithCarry(lo[1], k, QFieldModulus[0], 0)
	lo[2], carry = MACWithCarry(lo[2], k, QFieldModulus[1], carry)
	lo[3], carry = MACWithCarry(lo[3], k, QFieldModulus[2], carry)
	lo[4], carry = MACWithCarry(lo[4], k, QFieldModulus[3], carry)
	lo[5], carry = MACWithCarry(lo[5], k, QFieldModulus[4], carry)
	hi[0], carry = MACWithCarry(hi[0], k, QFieldModulus[5], carry)
	hi[1], carry = AddWithCarry(hi[1], carry2, carry)
	carry2 = carry
	k = lo[2] * montInvFQ
	_, carry = MACWithCarry(lo[2], k, QFieldModulus[0], 0)
	lo[3], carry = MACWithCarry(lo[3], k, QFieldModulus[1], carry)
	lo[4], carry = MACWithCarry(lo[4], k, QFieldModulus[2], carry)
	lo[5], carry = MACWithCarry(lo[5], k, QFieldModulus[3], carry)
	hi[0], carry = MACWithCarry(hi[0], k, QFieldModulus[4], carry)
	hi[1], carry = MACWithCarry(hi[1], k, QFieldModulus[5], carry)
	hi[2], carry = AddWithCarry(hi[2], carry2, carry)
	carry2 = carry
	k = lo[3] * montInvFQ
	_, carry = MACWithCarry(lo[3], k, QFieldModulus[0], 0)
	lo[4], carry = MACWithCarry(lo[4], k, QFieldModulus[1], carry)
	lo[5], carry = MACWithCarry(lo[5], k, QFieldModulus[2], carry)
	hi[0], carry = MACWithCarry(hi[0], k, QFieldModulus[3], carry)
	hi[1], carry = MACWithCarry(hi[1], k, QFieldModulus[4], carry)
	hi[2], carry = MACWithCarry(hi[2], k, QFieldModulus[5], carry)
	hi[3], carry = AddWithCarry(hi[3], carry2, carry)
	carry2 = carry
	k = lo[4] * montInvFQ
	_, carry = MACWithCarry(lo[4], k, QFieldModulus[0], 0)
	lo[5], carry = MACWithCarry(lo[5], k, QFieldModulus[1], carry)
	hi[0], carry = MACWithCarry(hi[0], k, QFieldModulus[2], carry)
	hi[1], carry = MACWithCarry(hi[1], k, QFieldModulus[3], carry)
	hi[2], carry = MACWithCarry(hi[2], k, QFieldModulus[4], carry)
	hi[3], carry = MACWithCarry(hi[3], k, QFieldModulus[5], carry)
	hi[4], carry = AddWithCarry(hi[4], carry2, carry)
	carry2 = carry
	k = lo[5] * montInvFQ
	_, carry = MACWithCarry(lo[5], k, QFieldModulus[0], 0)
	hi[0], carry = MACWithCarry(hi[0], k, QFieldModulus[1], carry)
	hi[1], carry = MACWithCarry(hi[1], k, QFieldModulus[2], carry)
	hi[2], carry = MACWithCarry(hi[2], k, QFieldModulus[3], carry)
	hi[3], carry = MACWithCarry(hi[3], k, QFieldModulus[4], carry)
	hi[4], carry = MACWithCarry(hi[4], k, QFieldModulus[5], carry)
	hi[5], carry = AddWithCarry(hi[5], carry2, carry)
	return hi
}

//
func AddNoCarry(a, b [6]uint64) (out [6]uint64) {
	carry := uint64(0)
	for i := 0; i < 6; i++ {
		out[i], carry = AddWithCarry(a[i], b[i], carry)
	}

	return out
}

func SubNoBorrow(a, b [6]uint64) (out [6]uint64) {
	borrow := uint64(0)
	for i := 0; i < 6; i++ {
		out[i], borrow = SubWithBorrow(a[i], b[i], borrow)
	}
	return out
}

func AddWithCarry(a, b, carry uint64) (uint64, uint64) {
	out, outCarry := bits.Add64(a, b, 0)
	out, outCarry2 := bits.Add64(out, carry, 0)
	return out, outCarry + outCarry2
}

var oneLsh64 = new(big.Int).Add(new(big.Int).SetUint64(0xffffffffffffffff), big.NewInt(1))

func SubWithBorrow(a, b, borrow uint64) (uint64, uint64) {
	o, c := bits.Sub64(a, b, borrow)
	return o, c
}

func MACWithCarry(a, b, c, carry uint64) (out uint64, newCarry uint64) {
	carryOut, bc := bits.Mul64(b, c)
	abc, carryOut2 := bits.Add64(bc, a, 0)
	abcc, carryOut3 := bits.Add64(abc, carry, 0)

	return abcc, carryOut + carryOut2 + carryOut3
}
