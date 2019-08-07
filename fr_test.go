package bls

import (
	"crypto/rand"
	"testing"
)

func TestFRInverse(t *testing.T) {
	one := FRReprToFR(&FRRepr{1, 0, 0, 0})
	for i := 0; i < 10; i++ {
		newFR, _ := RandFR(rand.Reader)
		inverse := newFR.Inverse()
		newFR.MulAssign(inverse)
		if !one.Equals(newFR) {
			t.Errorf("Multiplication with inverse must be one.")
		}
	}
}
