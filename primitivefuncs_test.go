package bls_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/phoreproject/bls"
)

func BenchmarkMACWithCarry(b *testing.B) {
	carry := uint64(0)
	for i := 0; i < b.N; i++ {
		_, carry = bls.MACWithCarry(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, 0x1234567812345678, carry)
	}
}

func BenchmarkSubWithCarry(b *testing.B) {
	borrow := uint64(0)
	for i := 0; i < b.N; i++ {
		_, borrow = bls.SubWithBorrow(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, borrow)
	}
}

func TestSubWithCarry(t *testing.T) {
	cases := []struct {
		a         uint64
		b         uint64
		borrow    uint64
		out       uint64
		outBorrow uint64
	}{
		{
			a:         0,
			b:         3,
			borrow:    0,
			out:       18446744073709551613,
			outBorrow: 1,
		},
		{
			a:         0,
			b:         2,
			borrow:    0,
			out:       18446744073709551614,
			outBorrow: 1,
		},
		{
			a:         0,
			b:         2,
			borrow:    1,
			out:       18446744073709551613,
			outBorrow: 1,
		},
		{
			a:         2,
			b:         0,
			borrow:    0,
			out:       2,
			outBorrow: 0,
		},
		{
			a:         2,
			b:         0,
			borrow:    1,
			out:       1,
			outBorrow: 0,
		},
		{
			a:         2,
			b:         1,
			borrow:    0,
			out:       1,
			outBorrow: 0,
		},
		{
			a:         2,
			b:         1,
			borrow:    1,
			out:       0,
			outBorrow: 0,
		},
		{
			a:         0,
			b:         0,
			borrow:    1,
			out:       18446744073709551615,
			outBorrow: 1,
		},
	}

	for _, c := range cases {
		borrow := c.borrow
		out, borrow := bls.SubWithBorrow(c.a, c.b, borrow)
		if out != c.out {
			t.Fatalf("%d - %d - %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.borrow, out, c.out)
		}
		if borrow != c.outBorrow {
			t.Fatalf("%d - %d - %d is giving incorrect borrow of %d instead of %d", c.a, c.b, c.borrow, borrow, c.outBorrow)
		}
	}
}

func BenchmarkAddWithCarry(b *testing.B) {
	borrow := uint64(0)
	for i := 0; i < b.N; i++ {
		_, borrow = bls.AddWithCarry(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, borrow)
	}
}

func TestAddWithCarry(t *testing.T) {
	cases := []struct {
		a        uint64
		b        uint64
		carry    uint64
		out      uint64
		outCarry uint64
	}{
		{
			1,
			1,
			0,
			2,
			0,
		},
		{
			1,
			0,
			0,
			1,
			0,
		},
		{
			0,
			1,
			0,
			1,
			0,
		},
		{
			1,
			1,
			1,
			3,
			0,
		},
		{
			18446744073709551615,
			1,
			0,
			0,
			1,
		},
		{
			18446744073709551615,
			0,
			1,
			0,
			1,
		},
		{
			0,
			0,
			0,
			0,
			0,
		},
		{
			0,
			4043378133346814763,
			0,
			4043378133346814763,
			0,
		},
	}

	for _, c := range cases {
		carry := c.carry
		out, carry := bls.AddWithCarry(c.a, c.b, carry)
		if out != c.out {
			t.Errorf("%d + %d + %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.carry, out, c.out)
		}
		if carry != c.outCarry {
			t.Errorf("%d + %d + %d is giving incorrect carry of %d instead of %d", c.a, c.b, c.carry, carry, c.outCarry)
		}
	}
}

func TestMACWithCarry(t *testing.T) {
	cases := []struct {
		a        uint64
		b        uint64
		c        uint64
		carry    uint64
		out      uint64
		outCarry uint64
	}{
		{
			0,
			1,
			1,
			0,
			1,
			0,
		},
		{
			0,
			4294967296,
			4294967296,
			0,
			0,
			1,
		},
		{
			0,
			4294967296,
			4294967296,
			1,
			1,
			1,
		},
		{
			5,
			4294967296,
			4294967296,
			1,
			6,
			1,
		},
	}

	for _, c := range cases {
		out, carry := bls.MACWithCarry(c.a, c.b, c.c, c.carry)
		if out != c.out {
			t.Fatalf("%d + %d * %d + %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.c, c.carry, out, c.out)
		}
		if carry != c.outCarry {
			t.Fatalf("%d + %d * %d + %d is giving incorrect carry of %d instead of %d", c.a, c.b, c.c, c.carry, carry, c.outCarry)
		}
	}
}

var oneLsh384MinusOne = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 384), big.NewInt(1))

func TestMultiplyFQReprOverflow(t *testing.T) {
	f0 := bls.FQRepr{4276637899304358534, 4043378133346814763, 8835052805473178628, 2680116066972705497, 18387885609531466875, 90398708109242637}
	f1 := bls.FQRepr{6568974633585825615, 15677163513955518067, 16490785605261833339, 9784757811163378176, 10803760609847905278, 1860524254683672351}

	expectedLo, _ := new(big.Int).SetString("11236684981931288970748045803123803013547375290961326241645975575468319186166860317135684626985730476857897168732506", 10)
	expectedHi, _ := new(big.Int).SetString("19474954426416495774628754584220018623553762741680760536161274135926885656938068449644545696043757020826963123558", 10)

	hi, lo := bls.MultiplyFQRepr(f0, f1)
	loBig := bls.FQRepr(lo).ToBig()
	hiBig := bls.FQRepr(hi).ToBig()

	if loBig.Cmp(expectedLo) != 0 {
		t.Fatalf("expected lo bits to equal %x, got: %x", expectedLo, loBig)
	}

	if hiBig.Cmp(expectedHi) != 0 {
		t.Fatalf("expected hi bits to equal %x, got: %x", expectedHi, hiBig)
	}
}

func TestRandomMultiplyFQRepr(t *testing.T) {
	r := NewXORShift(1)
	total := big.NewInt(1)
	totalFQ := bls.NewFQRepr(1)

	for i := 0; i < 1000000; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		_, totalFQ = bls.MultiplyFQRepr(totalFQ, f)
		total.Mul(total, n)
		total.And(total, oneLsh384MinusOne)

		if bls.FQRepr(totalFQ).ToBig().Cmp(total) != 0 {
			t.Fatal("multiplication totals do not match between big int and FQ")
		}
	}
}
