package bls_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/phoreproject/bls"
)

func TestEvenOdd(t *testing.T) {
	in := bls.NewFQRepr(1)
	if in.IsEven() {
		t.Error("1 should not be even")
	}
	if !in.IsOdd() {
		t.Error("1 should be odd")
	}
	in = bls.NewFQRepr(2)
	if !in.IsEven() {
		t.Error("2 should be even")
	}
	if in.IsOdd() {
		t.Error("2 should not be odd")
	}
}

func TestIsZero(t *testing.T) {
	in := bls.NewFQRepr(0)
	if !in.IsZero() {
		t.Error("0.IsZero() should be true")
	}
}

func TestRsh(t *testing.T) {
	in := bls.NewFQRepr(100)
	in.Rsh(1)
	if !in.Equals(bls.NewFQRepr(50)) {
		t.Error("100 >> 1 should equal 50")
	}
	i, _ := bls.FQReprFromString("10000000000000000000000000", 10)
	iRsh1, _ := bls.FQReprFromString("5000000000000000000000000", 10)
	i.Rsh(1)
	if !i.Equals(iRsh1) {
		t.Error("10000000000000000000000000 >> 1 should equal 5000000000000000000000000")
	}
}

func TestLsh(t *testing.T) {
	in := bls.NewFQRepr(100)
	in.Lsh(1)
	if !in.Equals(bls.NewFQRepr(200)) {
		t.Error("100 << 1 should equal 200")
	}
	i, _ := bls.FQReprFromString("10000000000000000000000000", 10)
	iRsh1, _ := bls.FQReprFromString("20000000000000000000000000", 10)
	i.Lsh(1)
	if !i.Equals(iRsh1) {
		t.Error("10000000000000000000000000 << 1 should equal 20000000000000000000000000")
	}
}

func TestRandomMACWithCarry(t *testing.T) {
	carry := uint64(0)
	carryBig := big.NewInt(0)
	current := uint64(0)
	currentBig := big.NewInt(0)

	r := NewXORShift(200)

	for i := 0; i < TestSamples; i++ {
		a, _ := rand.Int(r, new(big.Int).SetUint64(0xffffffffffffffff))
		b, _ := rand.Int(r, new(big.Int).SetUint64(0xffffffffffffffff))
		current = bls.MACWithCarry(current, a.Uint64(), b.Uint64(), &carry)

		a.Mul(a, b)
		currentBig.Add(currentBig, a)
		currentBig.Add(currentBig, carryBig)
		carryBig = new(big.Int).Rsh(currentBig, 64)
		currentBig.And(currentBig, new(big.Int).SetUint64(0xffffffffffffffff))

		if current != currentBig.Uint64() {
			t.Fatal("current != currentBig")
		}

		if carry != carryBig.Uint64() {
			t.Fatal("current != currentBig")
		}
	}
}

func TestRandomAddWithCarry(t *testing.T) {
	carry := uint64(0)
	carryBig := big.NewInt(0)
	current := uint64(0)
	currentBig := big.NewInt(0)

	r := NewXORShift(200)

	for i := 0; i < TestSamples; i++ {
		a, _ := rand.Int(r, new(big.Int).SetUint64(0xffffffffffffffff))
		current = bls.AddWithCarry(current, a.Uint64(), &carry)

		currentBig.Add(currentBig, a)
		currentBig.Add(currentBig, carryBig)
		carryBig = new(big.Int).Rsh(currentBig, 64)
		currentBig.And(currentBig, new(big.Int).SetUint64(0xffffffffffffffff))

		if current != currentBig.Uint64() {
			t.Fatal("current != currentBig")
		}

		if carry != carryBig.Uint64() {
			t.Fatal("current != currentBig")
		}
	}
}
