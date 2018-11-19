package bls

import (
	"errors"
	"math/big"
)

var curveOrder, _ = new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

var B = NewFQ(big.NewInt(3), fieldModulus)

var b2First, _ = NewFQ2([]*FQ{
	NewFQ(big.NewInt(3), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
})

var b2Second, _ = NewFQ2([]*FQ{
	NewFQ(big.NewInt(9), fieldModulus),
	NewFQ(big.NewInt(1), fieldModulus),
})

var B2 = b2First.Div(b2Second)

var b12, _ = NewFQ12([]*FQ{
	NewFQ(big.NewInt(3), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
	NewFQ(big.NewInt(0), fieldModulus),
})

var G1 = [3]*FQ{
	NewFQ(big.NewInt(1), fieldModulus),
	NewFQ(big.NewInt(2), fieldModulus),
	NewFQ(big.NewInt(1), fieldModulus),
}

var g211, _ = new(big.Int).SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781", 10)
var g212, _ = new(big.Int).SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634", 10)
var g221, _ = new(big.Int).SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930", 10)
var g222, _ = new(big.Int).SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531", 10)

var g21, _ = NewFQ2([]*FQ{
	NewFQ(g211, fieldModulus),
	NewFQ(g212, fieldModulus),
})

var g22, _ = NewFQ2([]*FQ{
	NewFQ(g221, fieldModulus),
	NewFQ(g222, fieldModulus),
})

var G2 = [3]*FQP{
	g21,
	g22,
	FQ2One(),
}

var z1 = [3]*FQ{
	&FQ{n: big.NewInt(1)},
	&FQ{n: big.NewInt(1)},
	&FQ{n: big.NewInt(0)},
}

var z2 = [3]*FQP{
	FQ2One(),
	FQ2One(),
	FQ2Zero(),
}

// IsInfFQ checks if FQ is infinite.
func IsInfFQ(in []*FQ) bool {
	return in[len(in)-1].Equals(&FQ{n: big.NewInt(0)})
}

// IsInfFQP checks if FQ2 is infinite.
func IsInfFQP(in []*FQP) (bool, error) {
	if len(in) == 0 {
		return false, errors.New("the point is 0 dimensional")
	}
	if len(in[0].elements) == 2 {
		return in[len(in)-1].Equals(FQ2Zero()), nil
	} else if len(in[0].elements) == 12 {
		return in[len(in)-1].Equals(FQ12Zero()), nil
	}
	return false, errors.New("the FQP is not of degree 2 or 12")
}

// IsOnCurveFQ checks if the FQ point is on the curve.
func IsOnCurveFQ(pt [3]*FQ, b *FQ) bool {
	if IsInfFQ(pt[:]) {
		return true
	}

	x := pt[0]
	y := pt[1]
	z := pt[2]

	return y.Exp(big.NewInt(2)).Mul(z).Sub(x.Exp(big.NewInt(3))).Equals(b.Mul(z.Exp(big.NewInt(3))))
}

// IsOnCurveFQP checks if an FQP is on the given curve.
func IsOnCurveFQP(pt [3]*FQP, b *FQP) (bool, error) {
	inf, err := IsInfFQP(pt[:])
	if err != nil {
		return false, err
	}
	if inf {
		return false, nil
	}
	x := pt[0]
	y := pt[1]
	z := pt[2]

	return y.Exp(big.NewInt(2)).Mul(z).Sub(x.Exp(big.NewInt(3))).Equals(b.Mul(z.Exp(big.NewInt(3)))), nil
}
