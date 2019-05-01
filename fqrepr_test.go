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
		current, carry = bls.MACWithCarry(current, a.Uint64(), b.Uint64(), carry)

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
		current, carry = bls.AddWithCarry(current, a.Uint64(), carry)

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

func TestRandomBytesFromBytes(t *testing.T) {
	r := NewXORShift(200)

	for i := 0; i < TestSamples; i++ {
		a, _ := bls.RandFQ(r)

		repr := a.ToRepr()

		reprBytes := repr.Bytes()
		newRepr := bls.FQReprFromBytes(reprBytes)

		if !newRepr.Equals(repr) {
			t.Fatal("FQRepr serialization/deserialization failed")
		}
	}
}

func TestMontReduce(t *testing.T) {
	hi := bls.FQRepr{1839480447532984087, 3926924786351924057, 16772763484671214791, 603559161728877186, 17550439636508622399, 43104129123962762}
	lo := bls.FQRepr{3464588527451676294, 7118675286620626963, 8738819137066355118, 12891471781342038311, 14889512650477226038, 12078494706333498109}

	expected := bls.FQRepr{13451288730302620273, 10097742279870053774, 15949884091978425806, 5885175747529691540, 1016841820992199104, 845620083434234474}

	out := bls.FQRepr(bls.MontReduce(hi, lo))

	if !out.Equals(expected) {
		t.Fatal("mont reduce returning incorrect values")
	}
}
