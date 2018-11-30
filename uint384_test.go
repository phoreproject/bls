package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func BenchmarkMACWithCarry(b *testing.B) {
	carry := uint64(0)
	for i := 0; i < b.N; i++ {
		bls.MACWithCarry(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, 0x1234567812345678, &carry)
	}
}

func BenchmarkSubWithCarry(b *testing.B) {
	borrow := uint64(0)
	for i := 0; i < b.N; i++ {
		bls.SubWithBorrow(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, &borrow)
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
		out := bls.SubWithBorrow(c.a, c.b, &borrow)
		if out != c.out {
			t.Fatalf("%d - %d - %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.borrow, out, c.out)
		}
		if borrow != c.outBorrow {
			t.Fatalf("%d - %d - %d is giving incorrect borrow of %d instead of %d", c.a, c.b, c.borrow, borrow, c.borrow)
		}
	}
}

func BenchmarkAddWithCarry(b *testing.B) {
	borrow := uint64(0)
	for i := 0; i < b.N; i++ {
		bls.AddWithCarry(0xFFFFFFFF00000000, 0x00000000FFFFFFFF, &borrow)
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
	}

	for _, c := range cases {
		carry := c.carry
		out := bls.AddWithCarry(c.a, c.b, &carry)
		if out != c.out {
			t.Fatalf("%d + %d + %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.carry, out, c.out)
		}
		if carry != c.outCarry {
			t.Fatalf("%d + %d + %d is giving incorrect carry of %d instead of %d", c.a, c.b, c.carry, carry, c.outCarry)
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
		carry := c.carry
		out := bls.MACWithCarry(c.a, c.b, c.c, &carry)
		if out != c.out {
			t.Fatalf("%d + %d * %d + %d is giving incorrect answer of %d instead of %d", c.a, c.b, c.c, c.carry, out, c.out)
		}
		if carry != c.outCarry {
			t.Fatalf("%d + %d * %d + %d is giving incorrect carry of %d instead of %d", c.a, c.b, c.c, c.carry, carry, c.outCarry)
		}
	}
}
