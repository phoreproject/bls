package bls_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/phoreproject/bls"
)

var bigZero = big.NewInt(0)

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

func TestPointAddition(t *testing.T) {
	fqZero := [3]*bls.FQ{
		bls.NewFQ(bigZero, bls.FieldModulus),
		bls.NewFQ(bigZero, bls.FieldModulus),
		bls.NewFQ(bigZero, bls.FieldModulus),
	}

	addZeros := bls.AddFQ(fqZero, fqZero)

	if !bls.FQEqual(addZeros, fqZero) {
		t.Fatal("AddFQ does not = 0 when adding two 0's")
	}

	fq1x, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208491", 10)
	fq1y, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208572", 10)

	fq1 := [3]*bls.FQ{
		bls.NewFQ(fq1x, bls.FieldModulus),
		bls.NewFQ(fq1y, bls.FieldModulus),
		bls.NewFQ(big.NewInt(64), bls.FieldModulus),
	}

	added := bls.AddFQ(fq1, bls.G1)

	expectedX, _ := new(big.Int).SetString("299200512", 10)
	expectedY, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894644690165063", 10)
	expectedZ, _ := new(big.Int).SetString("242970624", 10)

	expected := [3]*bls.FQ{
		bls.NewFQ(expectedX, bls.FieldModulus),
		bls.NewFQ(expectedY, bls.FieldModulus),
		bls.NewFQ(expectedZ, bls.FieldModulus),
	}

	if !bls.FQEqual(added, expected) {
		t.Fatal("AddFQ does not produce the expected result")
	}
}

func TestPointDoubling(t *testing.T) {
	fqOne := [3]*bls.FQ{
		bls.NewFQ(bigOne, bls.FieldModulus),
		bls.NewFQ(bigOne, bls.FieldModulus),
		bls.NewFQ(bigZero, bls.FieldModulus),
	}

	doubled := bls.DoubleFQ(fqOne)

	expectedY, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208556", 10)

	expected := [3]*bls.FQ{
		bls.NewFQ(bigZero, bls.FieldModulus),
		bls.NewFQ(expectedY, bls.FieldModulus),
		bls.NewFQ(bigZero, bls.FieldModulus),
	}

	if !bls.FQEqual(doubled, expected) {
		t.Fatal("DoubleFQ does not produce the expected result")
	}

	doubledG1 := bls.DoubleFQ(bls.G1)

	expectedX, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208491", 10)
	expectedY, _ = new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208572", 10)
	expectedZ, _ := new(big.Int).SetString("64", 10)

	expected = [3]*bls.FQ{
		bls.NewFQ(expectedX, bls.FieldModulus),
		bls.NewFQ(expectedY, bls.FieldModulus),
		bls.NewFQ(expectedZ, bls.FieldModulus),
	}

	if !bls.FQEqual(doubledG1, expected) {
		t.Fatal("DoubleFQ does not produce the expected result")
	}
}

func TestPointMultiplication(t *testing.T) {
	mulled := bls.MultiplyFQ(bls.G1, big.NewInt(3))

	expectedX, _ := new(big.Int).SetString("299200512", 10)
	expectedY, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894644690165063", 10)
	expectedZ, _ := new(big.Int).SetString("242970624", 10)

	expected := [3]*bls.FQ{
		bls.NewFQ(expectedX, bls.FieldModulus),
		bls.NewFQ(expectedY, bls.FieldModulus),
		bls.NewFQ(expectedZ, bls.FieldModulus),
	}

	fmt.Println(mulled)

	if !bls.FQEqual(mulled, expected) {
		t.Fatal("MultiplyFQ does not produce the expected result")
	}
}
