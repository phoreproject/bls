package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func BenchmarkFQAdd(b *testing.B) {
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
		inData[count].f1.Add(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQSub(b *testing.B) {
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
		inData[count].f1.Sub(inData[count].f2)
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
		inData[count].f1.Double()
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
		inData[count].f1.Square()
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
		inData[count].f1.Neg()
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
