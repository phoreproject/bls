package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestFQ6MultiplyByNonresidue(t *testing.T) {
	nqr := bls.NewFQ6(bls.FQ2Zero, bls.FQ2One, bls.FQ2Zero)

	for i := 0; i < 1000; i++ {
		a, err := bls.RandFQ6()
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
		c1, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ6()
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
		c0, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		c1, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ6()
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
