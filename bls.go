package bls

import "math/big"

// MillerLoopItem are the inputs to the miller loop.
type MillerLoopItem struct {
	p *G1Affine
	q *G2Prepared
}

type pairingItem struct {
	p      *G1Affine
	q      [][3]*FQ2
	qIndex int
}

// MillerLoop runs the miller loop algorithm.
func MillerLoop(items []MillerLoopItem) *FQ12 {
	pairs := make([]pairingItem, len(items))
	for i, item := range items {
		if !item.p.IsZero() && !item.q.IsZero() {
			pairs[i] = pairingItem{
				p:      item.p.Copy(),
				q:      item.q.coeffs,
				qIndex: 0,
			}
		}
	}

	ell := func(f *FQ12, coeffs [3]*FQ2, p *G1Affine) {
		c0 := coeffs[0]
		c1 := coeffs[1]

		c0.c0 = c0.c0.Mul(p.y)
		c0.c1 = c0.c0.Mul(p.y)
		c1.c0 = c0.c0.Mul(p.x)
		c1.c1 = c0.c0.Mul(p.x)

		f.MulBy014(coeffs[2], c1, c0)
	}

	f := FQ12One.Copy()

	foundOne := false
	blsXRsh1 := new(big.Int).Rsh(blsX, 1)
	bs := blsXRsh1.Bytes()
	for i := uint(0); i < uint(len(bs)*8); i++ {
		segment := i / 8
		bit := 7 - i%8
		set := bs[segment]&(1<<bit) > 0
		if !foundOne {
			foundOne = set
			continue
		}

		for i, pair := range pairs {
			ell(f, pair.q[pair.qIndex], pair.p.Copy())
			pairs[i].qIndex++
		}
		if set {
			for i, pair := range pairs {
				ell(f, pair.q[pair.qIndex], pair.p.Copy())
				pairs[i].qIndex++
			}
		}

		f = f.Square()
	}
	for i, pair := range pairs {
		ell(f, pair.q[pair.qIndex], pair.p.Copy())
		pairs[i].qIndex++
	}

	if blsIsNegative {
		f = f.Conjugate()
	}
	return f
}

// FinalExponentiation performs the final exponentiation on the
// FQ12 element.
func FinalExponentiation(r *FQ12) *FQ12 {
	f1 := r.Conjugate()
	f2 := r.Inverse()
	if f1 == nil {
		return nil
	}
	r = f1.Mul(f2)
	f2 = r
	r = r.FrobeniusMap(2).Mul(f2)

	ExpByX := func(f *FQ12, x *big.Int) *FQ12 {
		newf := f.Exp(x)
		if blsIsNegative {
			newf = newf.Conjugate()
		}
		return newf
	}

	x := new(big.Int).Set(blsX)

	y0 := r.Square()
	y1 := ExpByX(y0, x)
	x.Rsh(x, 1)
	y2 := ExpByX(y1, x)
	x.Lsh(x, 1)
	y3 := r.Conjugate()
	y1 = y1.Mul(y3).Conjugate().Mul(y2)
	y2 = ExpByX(y1, x)
	y3 = ExpByX(y2, x)
	y1 = y1.Conjugate()
	y3 = y3.Mul(y1)
	y1 = y1.Conjugate().FrobeniusMap(3)
	y2 = y2.FrobeniusMap(2)
	y1 = y1.Mul(y2)
	y2 = ExpByX(y3, x).Mul(y0).Mul(r).Mul(y2)
	y2 = y3.FrobeniusMap(1)
	y1 = y1.Mul(y2)
	return y1.Mul(y2)
}

// Pairing performs a pairing given the G1 and G2 elements.
func Pairing(p *G1Projective, q *G2Projective) *FQ12 {
	return FinalExponentiation(MillerLoop([]MillerLoopItem{
		{p.ToAffine(), G2AffineToPrepared(q.ToAffine())},
	}))
}
