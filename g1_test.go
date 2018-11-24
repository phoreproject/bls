package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func TestG1Generator(t *testing.T) {
	x := bls.FQZero.Copy()
	i := 0

	for {
		// y^2 = x^3 + b
		rhs := x.Square().Mul(x).Add(bls.NewFQ(bls.BCoeff))

		y := rhs.Sqrt()

		if y != nil {
			negY := y.Neg()
			pY := negY

			if y.Cmp(negY) < 0 {
				pY = y
			}

			p := bls.NewG1Affine(x, pY)

			if p.IsInCorrectSubgroupAssumingOnCurve() {
				t.Fatal("new point should be in subgroup")
			}

			g1 := p.ScaleByCofactor()

			if !g1.IsZero() {
				if i != 4 {
					t.Fatal("non-zero point should be 4th point")
				}

				g1 := g1.ToAffine()

				if !g1.IsInCorrectSubgroupAssumingOnCurve() {
					t.Fatal("point is not in correct subgroup")
				}

				if !g1.Equals(bls.G1AffineOne) {
					t.Fatal("point is not equal to generator point")
				}
				break
			}
		}

		i += 1
		x = x.Add(bls.FQOne)
	}
}
