package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/phoreproject/bls"
)

func TestFQ6MultiplyByNonresidue(t *testing.T) {
	nqr := bls.NewFQ6(bls.FQ2Zero, bls.FQ2One, bls.FQ2Zero)

	for i := 0; i < 1000; i++ {
		a, err := bls.RandFQ6(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		b.MulAssign(nqr)
		a.MulByNonresidueAssign()
		if !a.Equals(b) {
			t.Fatal("FQ6.MulByNonresidue not working properly")
		}
	}
}

func TestFQ6MultiplyBy1(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c1, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ6(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		b.MulAssign(bls.NewFQ6(bls.FQ2Zero, c1, bls.FQ2Zero))
		a.MulBy1Assign(c1)

		if !a.Equals(b) {
			t.Fatal("FQ6.MulBy1 not working")
		}
	}
}

func TestFQ6MultiplyBy01(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c0, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		c1, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ6(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		b.MulAssign(bls.NewFQ6(c0, c1, bls.FQ2Zero))
		a.MulBy01Assign(c0, c1)

		if !a.Equals(b) {
			t.Fatal("FQ6.MulBy1 not working")
		}
	}
}

func BenchmarkFQ6AddAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ6
		f2 *bls.FQ6
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ6(r)
		f2, _ := bls.RandFQ6(r)
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

func BenchmarkFQ6SubAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ6
		f2 *bls.FQ6
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ6(r)
		f2, _ := bls.RandFQ6(r)
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

func BenchmarkFQ6MulAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ6
		f2 *bls.FQ6
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ6(r)
		f2, _ := bls.RandFQ6(r)
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

func BenchmarkFQ6SquareAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ6
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ6(r)
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

func BenchmarkFQ6InverseAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ6
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ6(r)
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
