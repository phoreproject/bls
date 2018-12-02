package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/phoreproject/bls"
)

func TestFQ12MulBy014(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c0, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		c1, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		c5, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ12(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		b.MulAssign(bls.NewFQ12(
			bls.NewFQ6(c0, c1, bls.FQ2Zero),
			bls.NewFQ6(bls.FQ2Zero, c5, bls.FQ2Zero),
		))
		a.MulBy014Assign(c0, c1, c5)

		if !a.Equals(b) {
			t.Fatal("MulBy014 is broken.")
		}
	}
}

func BenchmarkFQ12MulAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
		f2 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
		f2, _ := bls.RandFQ12(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.MulAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ12SquareAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
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

func BenchmarkFQ12InverseAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
		inData[i] = addData{
			f1: f1,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.InverseAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}
