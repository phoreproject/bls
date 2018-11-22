package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestFQ12MulBy014(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c0, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		c1, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		c5, err := bls.RandFQ2()
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ12()
		if err != nil {
			t.Fatal(err)
		}
		b := a.Mul(bls.NewFQ12(
			bls.NewFQ6(c0, c1, bls.FQ2Zero),
			bls.NewFQ6(bls.FQ2Zero, c5, bls.FQ2Zero),
		))
		a = a.MulBy014(c0, c1, c5)

		if !a.Equals(b) {
			t.Error("MulBy014 is broken.")
		}
	}
}
