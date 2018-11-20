package bls

import (
	"errors"
	"math/big"
)

var ateLoopCount, _ = new(big.Int).SetString("29793968203157093288", 10)
var logAteLoopCount = 63
var pseudoBinaryEncoding = [65]int{
	0, 0, 0, 1, 0, 1, 0, -1, 0, 0, 1, -1, 0, 0, 1, 0,
	0, 1, 1, 0, -1, 0, 0, 1, 0, -1, 0, 0, 0, 0, 1, 1,
	1, 0, 0, -1, 0, 0, 1, 0, 0, 0, 0, 0, -1, 0, 0, 1,
	1, 0, 0, -1, 0, 0, 0, 1, 1, 0, -1, 0, 0, 1, 0, 1,
}

var finalExponentiationPower = new(big.Int).Div(new(big.Int).Sub(new(big.Int).Exp(FieldModulus, big.NewInt(12), nil), bigOne), curveOrder)

// LineFuncFQ creates a function representing a line between p1
// and p2 and evaluates it at T. Returns a numerator and
// denominator.
func LineFuncFQ(p1 [3]*FQ, p2 [3]*FQ, t [3]*FQ) (*FQ, *FQ) {
	zero := NewFQ(bigZero, FieldModulus)
	three := NewFQ(big.NewInt(3), FieldModulus)
	x1, y1, z1 := p1[0], p1[1], p1[2]
	x2, y2, z2 := p2[0], p2[1], p2[2]
	xt, yt, zt := t[0], t[1], t[2]

	mNumerator := y2.Mul(z1).Sub(y1.Mul(z2))
	mDenominator := x2.Mul(z1).Sub(x1.Mul(z2))
	if !mDenominator.Equals(zero) {
		return mNumerator.Mul(xt.Mul(z1).Sub(x1.Mul(zt))).Sub(mDenominator.Mul(yt.Mul(z1).Sub(y1.Mul(zt)))),
			mDenominator.Mul(zt).Mul(z1)
	}
	mNumerator = x1.Mul(x1).Mul(three)
	mDenominator = y1.Mul(z1).Mul(NewFQ(bigTwo, FieldModulus))
	return mNumerator.Mul(xt.Mul(z1).Sub(x1.Mul(zt))).Sub(mDenominator.Mul(yt.Mul(z1).Sub(y1.Mul(zt)))),
		mDenominator.Mul(zt).Mul(z1)
}

// LineFuncFQP creates a function representing a line between p1
// and p2 and evaluates it at T. Returns a numerator and
// denominator.
func LineFuncFQP(p1 [3]*FQP, p2 [3]*FQP, t [3]*FQP) (*FQP, *FQP) {
	zero, _ := FQPZero(p1[0])
	three := NewFQ(big.NewInt(3), FieldModulus)
	x1, y1, z1 := p1[0], p1[1], p1[2]
	x2, y2, z2 := p2[0], p2[1], p2[2]
	xt, yt, zt := t[0], t[1], t[2]

	mNumerator := y2.Mul(z1).Sub(y1.Mul(z2))
	mDenominator := x2.Mul(z1).Sub(x1.Mul(z2))
	if !mDenominator.Equals(zero) {
		return mNumerator.Mul(xt.Mul(z1).Sub(x1.Mul(zt))).Sub(mDenominator.Mul(yt.Mul(z1).Sub(y1.Mul(zt)))),
			mDenominator.Mul(zt).Mul(z1)
	}
	mNumerator = x1.Mul(x1).MulScalar(three)
	mDenominator = y1.Mul(z1).MulScalar(NewFQ(bigTwo, FieldModulus))
	return mNumerator.Mul(xt.Mul(z1).Sub(x1.Mul(zt))).Sub(mDenominator.Mul(yt.Mul(z1).Sub(y1.Mul(zt)))),
		mDenominator.Mul(zt).Mul(z1)
}

// CastFQToFQ12 converts an field element to FQ12 form
func CastFQToFQ12(i *FQ) *FQP {
	f := NewFQ12([]*FQ{
		i,
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
		NewFQ(bigZero, FieldModulus),
	})
	return f
}

// CastPointToFQ12 casts a point to FQ12 form
func CastPointToFQ12(pt [3]*FQ) [3]*FQP {
	x, y, z := pt[0], pt[1], pt[2]
	return [3]*FQP{CastFQToFQ12(x), CastFQToFQ12(y), CastFQToFQ12(z)}
}

// Normalize normalizes the point.
func Normalize(p [3]*FQ) [2]*FQ {
	return [2]*FQ{
		p[0].Div(p[2]),
		p[1].Div(p[2]),
	}
}

// NormalizeFQP normalizes the point.
func NormalizeFQP(p [3]*FQP) [2]*FQP {
	return [2]*FQP{
		p[0].Div(p[2]),
		p[1].Div(p[2]),
	}
}

// MillerLoop is the main miller loop.
func MillerLoop(q [3]*FQP, p [3]*FQP, finalExponentiate bool) *FQP {
	r := q
	fNum, fDen := FQ12One(), FQ12One()
	for i := 63; i > 0; i-- {
		v := pseudoBinaryEncoding[i]
		n, d := LineFuncFQP(r, r, p)
		fNum = fNum.Mul(fNum).Mul(n)
		fDen = fDen.Mul(fDen).Mul(d)
		r = DoubleFQP(r)
		if v == 1 {
			n, d := LineFuncFQP(r, q, p)
			fNum = fNum.Mul(n)
			fDen = fDen.Mul(d)
			r = AddFQP(r, q)
		} else if v == -1 {
			nQ := NegFQP(q)
			n, d := LineFuncFQP(r, nQ, p)
			fNum = fNum.Mul(n)
			fDen = fDen.Mul(d)
			r = AddFQP(r, nQ)
		}
	}
	q1 := [3]*FQP{q[0].Exp(FieldModulus), q[1].Exp(FieldModulus), q[2].Exp(FieldModulus)}
	onCurve, err := IsOnCurveFQP(q1, B12)
	if err != nil {
		panic(err)
	}
	if !onCurve {
		panic("q1 is not on b12")
	}
	nQ2 := [3]*FQP{q1[0].Exp(FieldModulus), q1[1].Exp(FieldModulus).Neg(), q1[2].Exp(FieldModulus)}
	onCurve, err = IsOnCurveFQP(nQ2, B12)
	if err != nil {
		panic(err)
	}
	if !onCurve {
		panic("nQ2 is not on b12")
	}
	n1, d1 := LineFuncFQP(r, q1, p)
	r = AddFQP(r, q1)
	n2, d2 := LineFuncFQP(r, nQ2, p)
	f := fNum.Mul(n1).Mul(n2).Div(fDen.Mul(d1).Mul(d2))
	if finalExponentiate {
		return f.Exp(finalExponentiationPower)
	}
	return f
}

// Pairing calculates the pairing given the two base elements.
func Pairing(q [3]*FQP, p [3]*FQ, finalExponentiate bool) (*FQP, error) {
	onCurve, err := IsOnCurveFQP(q, B2)
	if err != nil {
		return nil, err
	}
	if !onCurve {
		return nil, errors.New("q is not on the b2 curve")
	}
	onCurve = IsOnCurveFQ(p, B)
	if !onCurve {
		return nil, errors.New("p is not on the b curve")
	}

	qZero, _ := FQPZero(q[2])

	if p[2].Equals(NewFQ(bigZero, FieldModulus)) || q[2].Equals(qZero) {
		return FQ12One(), nil
	}
	return MillerLoop(Twist(q), CastPointToFQ12(p), finalExponentiate), nil
}

// FinalExponentiateFQP performs the final exponentiation of the pairing.
func FinalExponentiateFQP(p *FQP) *FQP {
	return p.Exp(finalExponentiationPower)
}
