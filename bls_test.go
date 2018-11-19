package bls_test

import (
	"math/big"
	"testing"

	"github.com/phoreproject/bls"
)

var bigOne = big.NewInt(1)
var oneLsh256Minus1 = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 256), bigOne)

func TestHashToG2(t *testing.T) {
	o := bls.Blake([]byte("hello there..."))
	out := bls.HashToG2(o)

	px1, _ := new(big.Int).SetString("3024970319033122561273152708816972627135445766175496293797634841850801341960", 10)
	px2, _ := new(big.Int).SetString("17888274449893510155473035301137869599433912919703800262498960849824302963733", 10)
	py1, _ := new(big.Int).SetString("3921864413939991736637795721853123277524416918819216099882730681704583058617", 10)
	py2, _ := new(big.Int).SetString("12794533939718669739591797131633398411194657834562331177287006583450870753356", 10)
	pz1, _ := new(big.Int).SetString("18129289412360544431963621391702356710127202697964354762328354420214632253681", 10)
	pz2, _ := new(big.Int).SetString("21160333175316459947332528067163588917037477177638280831314217472833430057962", 10)

	expected := [3]*bls.FQP{
		bls.NewFQ2([]*bls.FQ{
			bls.NewFQ(px1, bls.FieldModulus),
			bls.NewFQ(px2, bls.FieldModulus),
		}),
		bls.NewFQ2([]*bls.FQ{
			bls.NewFQ(py1, bls.FieldModulus),
			bls.NewFQ(py2, bls.FieldModulus),
		}),
		bls.NewFQ2([]*bls.FQ{
			bls.NewFQ(pz1, bls.FieldModulus),
			bls.NewFQ(pz2, bls.FieldModulus),
		}),
	}

	if !bls.FQPEqual(out, expected) {
		t.Fatal("hash to G2 does not return expected result")
	}
}

func TestPrivToPub(t *testing.T) {
	expected, _ := new(big.Int).SetString("9366015879375004571250438303432407971238053874512316318402267084951246439740", 10)
	actual := bls.PrivToPub(big.NewInt(31))

	if expected.Cmp(actual) != 0 {
		t.Log(actual)
		t.Fatal("generate pub key does not match expected")
	}
}

func TestAcceptance(t *testing.T) {
	msgToSign := []byte("hello there...")

	privKey := big.NewInt(31)

	pubKey := bls.PrivToPub(privKey)

	signature, err := bls.Sign(bls.Blake(msgToSign), privKey)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := bls.Verify(bls.Blake(msgToSign), pubKey, signature)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("signature was not valid")
	}
}
