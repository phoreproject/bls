package bls

import (
	"errors"
	"math/big"
)

var curveOrder, _ = new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

// B is the first field.
var B = NewFQ(big.NewInt(3), FieldModulus)

var b2First = NewFQ2([]*FQ{
	NewFQ(big.NewInt(3), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
})

var b2Second = NewFQ2([]*FQ{
	NewFQ(big.NewInt(9), FieldModulus),
	NewFQ(big.NewInt(1), FieldModulus),
})

// B2 is the quadratic field extension of B.
var B2 = b2First.Div(b2Second)

// B12 is the 12th-degree field extension of B.
var B12 = NewFQ12([]*FQ{
	NewFQ(big.NewInt(3), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
	NewFQ(big.NewInt(0), FieldModulus),
})

// G1 is the B1 generator point.
var G1 = [3]*FQ{
	NewFQ(big.NewInt(1), FieldModulus),
	NewFQ(big.NewInt(2), FieldModulus),
	NewFQ(big.NewInt(1), FieldModulus),
}

var g211, _ = new(big.Int).SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781", 10)
var g212, _ = new(big.Int).SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634", 10)
var g221, _ = new(big.Int).SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930", 10)
var g222, _ = new(big.Int).SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531", 10)

var g21 = NewFQ2([]*FQ{
	NewFQ(g211, FieldModulus),
	NewFQ(g212, FieldModulus),
})

var g22 = NewFQ2([]*FQ{
	NewFQ(g221, FieldModulus),
	NewFQ(g222, FieldModulus),
})

// G2 is the B2 generator point.
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

	y.Exp(big.NewInt(2))

	return y.Exp(big.NewInt(2)).Mul(z).Sub(x.Exp(big.NewInt(3))).Equals(b.Mul(z.Exp(big.NewInt(3)))), nil
}

// DoubleFQ performs EC doubling.
func DoubleFQ(pt [3]*FQ) [3]*FQ {
	eight := big.NewInt(8)
	four := big.NewInt(4)

	x, y, z := pt[0], pt[1], pt[2]
	w := x.Mul(x).Mul(NewFQ(big.NewInt(3), FieldModulus))
	s := y.Mul(z)
	b := x.Mul(y).Mul(s)
	h := w.Mul(w).Sub(b.Mul(NewFQ(eight, FieldModulus)))
	sSquared := s.Mul(s)
	newX := NewFQ(bigTwo, FieldModulus).Mul(h).Mul(s)
	newY := w.Mul(NewFQ(four, FieldModulus).Mul(b).Sub(h)).Sub(NewFQ(eight, FieldModulus).Mul(y).Mul(y).Mul(sSquared))
	newZ := NewFQ(eight, FieldModulus).Mul(s).Mul(sSquared)
	return [3]*FQ{newX, newY, newZ}
}

// DoubleFQP performs EC doubling.
func DoubleFQP(pt [3]*FQP) [3]*FQP {
	eight := big.NewInt(8)
	four := big.NewInt(4)

	x, y, z := pt[0], pt[1], pt[2]
	w := x.Mul(x).MulScalar(NewFQ(big.NewInt(3), FieldModulus))
	s := y.Mul(z)
	b := x.Mul(y).Mul(s)
	h := w.Mul(w).Sub(b.MulScalar(NewFQ(eight, FieldModulus)))
	sSquared := s.Mul(s)
	newX := h.MulScalar(NewFQ(bigTwo, FieldModulus)).Mul(s)
	newY := w.Mul(b.MulScalar(NewFQ(four, FieldModulus)).Sub(h)).Sub(y.MulScalar(NewFQ(eight, FieldModulus)).Mul(y).Mul(sSquared))
	newZ := s.MulScalar(NewFQ(eight, FieldModulus)).Mul(sSquared)
	return [3]*FQP{newX, newY, newZ}
}

// AddFQ performs EC addition.
func AddFQ(pt1 [3]*FQ, pt2 [3]*FQ) [3]*FQ {
	one, zero := NewFQ(bigOne, FieldModulus), NewFQ(bigZero, FieldModulus)
	two := NewFQ(bigTwo, FieldModulus)
	if pt2[2].Equals(zero) {
		return pt1
	} else if pt1[2].Equals(zero) {
		return pt2
	}

	x1, y1, z1 := pt1[0], pt1[1], pt1[2]
	x2, y2, z2 := pt2[0], pt2[1], pt2[2]
	u1 := y2.Mul(z1)
	u2 := y1.Mul(z2)
	v1 := x2.Mul(z1)
	v2 := x1.Mul(z2)
	if v1.Equals(v2) && u1.Equals(u2) {
		return DoubleFQ(pt1)
	} else if v1.Equals(v2) {
		return [3]*FQ{one, one, zero}
	}
	u := u1.Sub(u2)
	v := v1.Sub(v2)
	vSquared := v.Mul(v)
	vSquaredTimesV2 := vSquared.Mul(v2)
	vCubed := vSquared.Mul(v)
	w := z1.Mul(z2)
	a := u.Mul(u).Mul(w).Sub(vCubed).Sub(vSquaredTimesV2.Mul(two))
	newX := v.Mul(a)
	newY := u.Mul(vSquaredTimesV2.Sub(a)).Sub(vCubed.Mul(u2))
	newZ := vCubed.Mul(w)
	return [3]*FQ{newX, newY, newZ}
}

// AddFQP performs EC addition.
func AddFQP(pt1 [3]*FQP, pt2 [3]*FQP) [3]*FQP {
	one, _ := FQPOne(pt1[0])
	zero, _ := FQPZero(pt1[0])
	two := NewFQ(bigTwo, FieldModulus)
	if pt2[2].Equals(zero) {
		return pt1
	} else if pt1[2].Equals(zero) {
		return pt2
	}

	x1, y1, z1 := pt1[0], pt1[1], pt1[2]
	x2, y2, z2 := pt2[0], pt2[1], pt2[2]
	u1 := y2.Mul(z1)
	u2 := y1.Mul(z2)
	v1 := x2.Mul(z1)
	v2 := x1.Mul(z2)
	if v1.Equals(v2) && u1.Equals(u2) {
		return DoubleFQP(pt1)
	} else if v1.Equals(v2) {
		return [3]*FQP{one, one, zero}
	}
	u := u1.Sub(u2)
	v := v1.Sub(v2)
	vSquared := v.Mul(v)
	vSquaredTimesV2 := vSquared.Mul(v2)
	vCubed := vSquared.Mul(v)
	w := z1.Mul(z2)
	a := u.Mul(u).Mul(w).Sub(vCubed).Sub(vSquaredTimesV2.MulScalar(two))
	newX := v.Mul(a)
	newY := u.Mul(vSquaredTimesV2.Sub(a)).Sub(vCubed.Mul(u2))
	newZ := vCubed.Mul(w)
	return [3]*FQP{newX, newY, newZ}
}

// MultiplyFQ performs EC multiplication.
func MultiplyFQ(point [3]*FQ, n *big.Int) [3]*FQ {
	if n.Cmp(bigZero) == 0 {
		one, zero := NewFQ(bigOne, FieldModulus), NewFQ(bigZero, FieldModulus)
		return [3]*FQ{one, one, zero}
	} else if n.Cmp(bigOne) == 0 {
		return point
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigOne) != 0 {
		return MultiplyFQ(DoubleFQ(point), new(big.Int).Div(n, bigTwo))
	} else {
		return AddFQ(MultiplyFQ(DoubleFQ(point), new(big.Int).Rsh(n, 1)), point)
	}
}

// MultiplyFQP performs EC multiplication.
func MultiplyFQP(point [3]*FQP, n *big.Int) [3]*FQP {
	if n.Cmp(bigZero) == 0 {
		one, _ := FQPOne(point[0])
		zero, _ := FQPZero(point[0])
		return [3]*FQP{one, one, zero}
	} else if n.Cmp(bigOne) == 0 {
		return point
	} else if new(big.Int).Mod(n, bigTwo).Cmp(bigOne) != 0 {
		return MultiplyFQP(DoubleFQP(point), new(big.Int).Div(n, bigTwo))
	} else {
		return AddFQP(MultiplyFQP(DoubleFQP(point), new(big.Int).Rsh(n, 1)), point)
	}
}

// FQPEqual checks if two points are equal?
// TODO: update these docs. something about DDH?
func FQPEqual(pt1 [3]*FQP, pt2 [3]*FQP) bool {
	x1, y1, z1 := pt1[0], pt1[1], pt1[2]
	x2, y2, z2 := pt2[0], pt2[1], pt2[2]
	return x1.Mul(z2).Equals(x2.Mul(z1)) && y1.Mul(z2).Equals(y2.Mul(z1))
}

// FQEqual checks if two points are equal?
// TODO: update these docs. something about DDH?
func FQEqual(pt1 [3]*FQ, pt2 [3]*FQ) bool {
	x1, y1, z1 := pt1[0], pt1[1], pt1[2]
	x2, y2, z2 := pt2[0], pt2[1], pt2[2]
	return x1.Mul(z2).Equals(x2.Mul(z1)) && y1.Mul(z2).Equals(y2.Mul(z1))
}

var w = NewFQ12([]*FQ{
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigOne, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
	NewFQ(bigZero, FieldModulus),
})

// NegFQ converts P to -P
func NegFQ(pt [3]*FQ) [3]*FQ {
	return [3]*FQ{pt[0].Copy(), pt[1].Neg(), pt[2].Copy()}
}

// NegFQP converts P to -P
func NegFQP(pt [3]*FQP) [3]*FQP {
	return [3]*FQP{pt[0].Copy(), pt[1].Neg(), pt[2].Copy()}
}

// Twist twists a field from Z[p] / x**2 to Z[p] / x**2 - 18*x + 82
func Twist(pt [3]*FQP) [3]*FQP {
	x, y, z := pt[0], pt[1], pt[2]
	nine := NewFQ(big.NewInt(9), FieldModulus)
	bigThree := big.NewInt(3)

	xCoeffs0 := x.elements[0].Sub(x.elements[1].Mul(nine))
	xCoeffs1 := x.elements[1]
	yCoeffs0 := y.elements[0].Sub(y.elements[1].Mul(nine))
	yCoeffs1 := y.elements[1]
	zCoeffs0 := z.elements[0].Sub(z.elements[1].Mul(nine))
	zCoeffs1 := z.elements[1]

	nxCoeffs := make([]*FQ, 12)
	nxCoeffs[0] = xCoeffs0
	nxCoeffs[6] = xCoeffs1
	for i, f := range nxCoeffs {
		if f == nil {
			nxCoeffs[i] = NewFQ(bigZero, FieldModulus)
		} else {
			f.fieldModulus = FieldModulus
		}
	}
	nx := NewFQ12(nxCoeffs)

	nyCoeffs := make([]*FQ, 12)
	nyCoeffs[0] = yCoeffs0
	nyCoeffs[6] = yCoeffs1
	for i, f := range nyCoeffs {
		if f == nil {
			nyCoeffs[i] = NewFQ(bigZero, FieldModulus)
		} else {
			f.fieldModulus = FieldModulus
		}
	}
	ny := NewFQ12(nyCoeffs)

	nzCoeffs := make([]*FQ, 12)
	nzCoeffs[0] = zCoeffs0
	nzCoeffs[6] = zCoeffs1
	for i, f := range nzCoeffs {
		if f == nil {
			nzCoeffs[i] = NewFQ(bigZero, FieldModulus)
		} else {
			f.fieldModulus = FieldModulus
		}
	}
	nz := NewFQ12(nzCoeffs)

	return [3]*FQP{nx.Mul(w.Exp(bigTwo)), ny.Mul(w.Exp(bigThree)), nz}
}

// G12 is the B12 generator point.
var G12 = Twist(G2)
