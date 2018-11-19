package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestIsOnCurve(t *testing.T) {
	if !bls.IsOnCurveFQ(bls.G1, bls.B) {
		t.Fatal("generator 1 is not on curve")
	}

	onCurve, err := bls.IsOnCurveFQP(bls.G2, bls.B2)
	if err != nil {
		t.Fatal(err)
	}
	if !onCurve {
		t.Fatal("generator 2 is not on curve")
	}

	onCurve, err = bls.IsOnCurveFQP(bls.G12, bls.B12)
	if err != nil {
		t.Fatal(err)
	}
	if !onCurve {
		t.Fatal("generator 12 is not on curve")
	}
}
