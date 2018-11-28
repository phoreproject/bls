package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestBasicAcceptance(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := bls.RandKey(r)
	pub := bls.PrivToPub(priv)
	msg := []byte("Hello world!")
	sig := bls.Sign(msg, priv)
	if !bls.Verify(msg, pub, sig) {
		t.Fatal("sig did not verify")
	}
}
