package bls_test

import (
	"crypto/rand"
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
			t.Fatal("subtraction totals do not match between big int and FQ")
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
		n, _ := rand.Int(r, QFieldModulusBig)

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

func TestSqrt(t *testing.T) {
	r := NewXORShift(1)

	for i := 0; i < 1000; i++ {
		f, err := bls.RandFQ(r)
		if err != nil {
			t.Fatal(err)
		}

		a, success := f.Sqrt()
		if !success {
			continue
		}
		a.SquareAssign()

		if !a.Equals(f) {
			t.Fatal("sqrt(a)^2 != a")
		}
	}
}

func TestInverse(t *testing.T) {
	// r := NewXORShift(1)

	for i := 0; i < 1; i++ {
		fRepr, err := bls.FQReprFromString("08aad39fba5b1d27bd5706262b1e2ee6c3da7dff5974ecbb0bee2bd75d4bc10973d8e59fd31f225247a335deb379592c", 16)
		if err != nil {
			t.Fatal(err)
		}

		f := bls.FQReprToFQ(fRepr)

		fInv, _ := f.Inverse()
		f.MulAssign(fInv)

		if !f.Equals(bls.FQOne) {
			t.Fatal("a*a^-1 != 1")
		}
	}
}

func BenchmarkFQAddAssign(b *testing.B) {
	type addData struct {
		f1 bls.FQ
		f2 bls.FQ
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
	type subData struct {
		f1 bls.FQ
		f2 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]subData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		f2, _ := bls.RandFQ(r)
		inData[i] = subData{
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
	type mulData struct {
		f1 bls.FQ
		f2 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]mulData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		f2, _ := bls.RandFQ(r)
		inData[i] = mulData{
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
	type doubleData struct {
		f1 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]doubleData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = doubleData{
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
	type squareData struct {
		f1 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]squareData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = squareData{
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
	type invData struct {
		f1 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]invData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = invData{
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
	type negData struct {
		f1 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]negData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = negData{
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
	type sqrtData struct {
		f1 bls.FQ
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]sqrtData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ(r)
		inData[i] = sqrtData{
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
