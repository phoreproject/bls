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
		b := a.Mul(nqr)
		a = a.MulByNonresidue()
		if !a.Equals(b) {
			t.Error("FQ6.MulByNonresidue not working properly")
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
		b := a.Mul(bls.NewFQ6(bls.FQ2Zero, c1, bls.FQ2Zero))
		a = a.MulBy1(c1)

		if !a.Equals(b) {
			t.Error("FQ6.MulBy1 not working")
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
		b := a.Mul(bls.NewFQ6(c0, c1, bls.FQ2Zero))
		a = a.MulBy01(c0, c1)

		if !a.Equals(b) {
			t.Error("FQ6.MulBy1 not working")
		}
	}
}

func BenchmarkFQ6Add(b *testing.B) {
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

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Add(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ6Sub(b *testing.B) {
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

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Sub(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ6Mul(b *testing.B) {
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

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Mul(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ6Square(b *testing.B) {
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

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Square()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ6Inverse(b *testing.B) {
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

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.Inverse()
		count = (count + 1) % g1MulAssignSamples
	}
}
