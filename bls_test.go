package bls_test

import (
	"math/big"
)

var bigOne = big.NewInt(1)
var oneLsh256Minus1 = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 256), bigOne)

// func TestHashToG2(t *testing.T) {
// 	o := bls.Blake([]byte("hello there..."))
// 	bls.HashToG2(o)
// 	t.Fatal()
// }

// func TestAcceptance(t *testing.T) {
// 	msgToSign := []byte("hello there...")

// 	privKey, err := rand.Int(rand.Reader, oneLsh256Minus1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	pubKey := bls.PrivToPub(privKey)

// 	signature, err := bls.Sign(bls.Blake(msgToSign), privKey)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	valid, err := bls.Verify(bls.Blake(msgToSign), pubKey, signature)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !valid {
// 		t.Fatal("signature was not valid")
// 	}
// }
