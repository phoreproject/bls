package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestBitIterator(t *testing.T) {
	a := bls.NewBitIterator([]uint64{0xa953d79b83f6ab59, 0x6dea2059e200bd39})
	expected := "01101101111010100010000001011001111000100000000010111101001110011010100101010011110101111001101110000011111101101010101101011001"

	for _, e := range expected {
		bit, done := a.Next()
		if done {
			t.Fatal("iterator finished too soon")
		}
		if bit != (e == '1') {
			t.Fatal("invalid bit")
		}
	}

	if _, done := a.Next(); done != true {
		t.Fatal("iterator did not finish in time")
	}
}
