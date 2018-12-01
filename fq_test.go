package bls_test

import (
	"crypto/rand"
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
