package bls

// MillerLoopItem are the inputs to the miller loop.
type MillerLoopItem struct {
	P *G1Affine
	Q *G2Prepared
}

type pairingItem struct {
	p      *G1Affine
	q      [][3]FQ2
	qIndex int
}

// MillerLoop runs the miller loop algorithm.
func MillerLoop(items []MillerLoopItem) *FQ12 {
	pairs := make([]pairingItem, len(items))
	for i, item := range items {
		if !item.P.IsZero() && !item.Q.IsZero() {
			pairs[i] = pairingItem{
				p:      item.P.Copy(),
				q:      item.Q.coeffs,
				qIndex: 0,
			}
		}
	}

	ell := func(f *FQ12, coeffs [3]FQ2, p *G1Affine) {

		c0 := coeffs[0]
		c1 := coeffs[1]

		c0.c0.MulAssign(p.y)
		c0.c1.MulAssign(p.y)
		c1.c0.MulAssign(p.x)
		c1.c1.MulAssign(p.x)

		f.MulBy014Assign(coeffs[2], c1, c0)
	}

	f := FQ12One.Copy()

	foundOne := false
	blsXRsh1 := blsX.Copy()
	blsXRsh1.Rsh(1)
	for q := uint(0); q <= blsXRsh1.BitLen(); q++ {
		set := blsXRsh1.Bit(blsXRsh1.BitLen() - q)
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

		f.SquareAssign()
	}
	for i, pair := range pairs {
		ell(f, pair.q[pair.qIndex], pair.p.Copy())
		pairs[i].qIndex++
	}

	if blsIsNegative {
		f.ConjugateAssign()
	}
	return f
}

// FinalExponentiation performs the final exponentiation on the
// FQ12 element.
func FinalExponentiation(r *FQ12) *FQ12 {
	f1 := r.Copy()
	f1.ConjugateAssign()
	f2 := r.Copy()
	if !f2.InverseAssign() {
		return nil
	}
	r = f1.Copy()
	r.MulAssign(f2)
	f2 = r.Copy()
	r.FrobeniusMapAssign(2)
	r.MulAssign(f2)

	ExpByX := func(f *FQ12, x FQRepr) *FQ12 {
		newf := f.Exp(x)
		if blsIsNegative {
			newf.ConjugateAssign()
		}
		return newf
	}

	x := blsX.Copy()

	y0 := r.Copy()
	y0.SquareAssign()
	y1 := ExpByX(y0, x)
	x.Rsh(1)
	y2 := ExpByX(y1, x)
	x.Lsh(1)
	y3 := r.Copy()
	y3.ConjugateAssign()
	y1 = y1.Copy()
	y1.MulAssign(y3)
	y1.ConjugateAssign()
	y1.MulAssign(y2)
	y2 = ExpByX(y1, x)
	y3 = ExpByX(y2, x)
	y1.ConjugateAssign()
	y3.MulAssign(y1)
	y1.ConjugateAssign()
	y1.FrobeniusMapAssign(3)
	y2.FrobeniusMapAssign(2)
	y1.MulAssign(y2)
	y2 = ExpByX(y3, x)
	y2.MulAssign(y0)
	y2.MulAssign(r)
	y1.MulAssign(y2)
	y3.FrobeniusMapAssign(1)
	y1.MulAssign(y3)
	return y1
}

// Pairing performs a pairing given the G1 and G2 elements.
func Pairing(p *G1Projective, q *G2Projective) *FQ12 {
	return FinalExponentiation(MillerLoop([]MillerLoopItem{
		{p.ToAffine(), G2AffineToPrepared(q.ToAffine())},
	}))
}

// CompareTwoPairings checks e(P1, Q1) == e(P2, Q2)
// <=> FE(ML(P1, Q1)ML(-P2, Q2)) == 1
func CompareTwoPairings(P1 *G1Projective, Q1 *G2Projective, P2 *G1Projective, Q2 *G2Projective) bool {
	negP2 := P2.Copy()
	negP2.NegAssign()
	return FinalExponentiation(
		MillerLoop(
			[]MillerLoopItem{{P1.ToAffine(), G2AffineToPrepared(Q1.ToAffine())}, {negP2.ToAffine(), G2AffineToPrepared(Q2.ToAffine())}})).Equals(FQ12One)

}
