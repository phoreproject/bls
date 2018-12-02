package bls_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/phoreproject/bls"
)

const TestSamples = 1000

func TestFQFromString(t *testing.T) {
	r := NewXORShift(1)

	for i := 0; i < TestSamples; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())
		s := n.String()

		f, err := bls.FQReprFromString(s, 10)
		if err != nil {
			t.Fatal(err)
		}
		if f.ToBig().Cmp(n) != 0 {
			t.Fatalf("big number does not match FQ (expected %x, got %s)", n, f)
		}
	}
}

func TestFQOneZero(t *testing.T) {
	if bls.FQOne.ToRepr().ToBig().Cmp(big.NewInt(1)) != 0 {
		t.Errorf("one does not equal 1. (expected: 1, actual: %s)", bls.FQOne.ToRepr().ToBig())
	}
	if bls.FQZero.ToRepr().ToBig().Cmp(big.NewInt(0)) != 0 {
		t.Errorf("one does not equal 0. (expected: 0, actual: %s)", bls.FQZero.ToRepr().ToBig())
	}
}

func TestFQCopy(t *testing.T) {
	r := NewXORShift(1)

	for i := 0; i < TestSamples; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		fCopy := f.Copy()
		fCopy.Div2()

		if f.Equals(fCopy) {
			t.Fatal("copy doesn't work")
		}
	}
}

func TestIsValid(t *testing.T) {
	f := bigOne.Copy()
	f.Lsh(383)
	f.AddNoCarry(bigOne)
	fmt.Println(f)
	fmt.Println(bls.QFieldModulus)
	if bls.FQReprToFQ(f) != nil {
		t.Fatal("2^383-1 should not be valid")
	}
}

var QFieldModulusBig = bls.QFieldModulus.ToBig()

func TestAddAssign(t *testing.T) {
	r := NewXORShift(1)
	total := big.NewInt(0)
	totalFQ := bls.FQReprToFQ(bls.NewFQRepr(0))

	for i := 0; i < TestSamples; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		totalFQ.AddAssign(bls.FQReprToFQ(f))
		total.Add(total, n)
		total.Mod(total, QFieldModulusBig)

		if totalFQ.ToRepr().ToBig().Cmp(total) != 0 {
			t.Error("addition totals do not match between big int and FQ")
		}
	}
}

func TestMulAssign(t *testing.T) {
	r := NewXORShift(1)
	total := big.NewInt(0)
	totalFQ := bls.FQReprToFQ(bls.NewFQRepr(0))

	for i := 0; i < TestSamples; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		totalFQ.MulAssign(bls.FQReprToFQ(f))
		total.Mul(total, n)
		total.Mod(total, QFieldModulusBig)

		if totalFQ.ToRepr().ToBig().Cmp(total) != 0 {
			t.Error("multiplication totals do not match between big int and FQ")
		}
	}
}

func TestSubAssign(t *testing.T) {
	r := NewXORShift(1)
	total := big.NewInt(0)
	totalFQ := bls.FQReprToFQ(bls.NewFQRepr(0))

	for i := 0; i < TestSamples; i++ {
		n, _ := rand.Int(r, bls.QFieldModulus.ToBig())

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		totalFQ.SubAssign(bls.FQReprToFQ(f))
		total.Sub(total, n)
		total.Mod(total, QFieldModulusBig)

		if totalFQ.ToRepr().ToBig().Cmp(total) != 0 {
			t.Error("multiplication totals do not match between big int and FQ")
		}
	}
}

func TestSquare(t *testing.T) {
	// r := NewXORShift(1)
	total := big.NewInt(2398928)
	totalFQ := bls.FQReprToFQ(bls.NewFQRepr(2398928))

	for i := 0; i < TestSamples; i++ {
		totalFQ.SquareAssign()
		total.Mul(total, total)
		total.Mod(total, QFieldModulusBig)

		if totalFQ.ToRepr().ToBig().Cmp(total) != 0 {
			t.Fatal("exp totals do not match between big int and FQ")
		}
	}
}

func TestExp(t *testing.T) {
	r := NewXORShift(1)
	total := big.NewInt(2)
	totalFQ := bls.FQReprToFQ(bls.NewFQRepr(2))

	for i := 0; i < 1; i++ {
		n, _ := rand.Int(r, big.NewInt(100))

		f, err := bls.FQReprFromBigInt(n)
		if err != nil {
			t.Fatal(err)
		}

		totalFQ = totalFQ.Exp(f)
		total.Exp(total, n, QFieldModulusBig)

		if totalFQ.ToRepr().ToBig().Cmp(total) != 0 {
			t.Fatal("exp totals do not match between big int and FQ")
		}
	}
}

func BenchmarkFQAddAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
		f2 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		f2, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.AddAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQSubAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
		f2 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		f2, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.SubAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQMulAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
		f2 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		f2, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		g := inData[count].f1.Copy()
		g.MulAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQMul2(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.DoubleAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQSquare(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.SquareAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQInverse(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Inverse()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQNegate(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.NegAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQSqrt(b *testing.B) {
	type addData struct {
		f1 *bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Sqrt()
		count = (count + 1) % g1MulAssignSamples
	}
}
