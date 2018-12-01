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
		rhs := x.Copy()
		rhs.SquareAssign()
		rhs.MulAssign(x)
		rhs.AddAssign(bls.FQReprToFQ(bls.BCoeff))

		y := rhs.Sqrt()

		if y != nil {
			negY := y.Copy()
			negY.NegAssign()
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
		x.AddAssign(bls.FQOne)
	}
}

type XORShift struct {
	state uint64
}

func NewXORShift(state uint64) *XORShift {
	return &XORShift{state}
}

func (xor *XORShift) Read(b []byte) (int, error) {
	for i := range b {
		x := xor.state
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		b[i] = uint8(x)
		xor.state = x
	}
	return len(b), nil
}

const g1MulAssignSamples = 10

func BenchmarkG1MulAssign(b *testing.B) {
	type mulData struct {
		g *bls.G1Projective
		f *bls.FR
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]mulData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		gx, _ := bls.RandFQ(r)
		gy, _ := bls.RandFQ(r)
		gz, _ := bls.RandFQ(r)
		randFR, _ := bls.RandFR(r)
		inData[i] = mulData{
			g: bls.NewG1Projective(gx, gy, gz),
			f: randFR,
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].g.Mul(inData[count].f.ToRepr().ToFQ())
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkG1AddAssign(b *testing.B) {
	type addData struct {
		g1 *bls.G1Projective
		g2 *bls.G1Projective
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		g1x, _ := bls.RandFQ(r)
		g1y, _ := bls.RandFQ(r)
		g1z, _ := bls.RandFQ(r)
		g2x, _ := bls.RandFQ(r)
		g2y, _ := bls.RandFQ(r)
		g2z, _ := bls.RandFQ(r)
		inData[i] = addData{
			g1: bls.NewG1Projective(g1x, g1y, g1z),
			g2: bls.NewG1Projective(g2x, g2y, g2z),
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].g1.Add(inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkG1AddAssignMixed(b *testing.B) {
	type addData struct {
		g1 *bls.G1Projective
		g2 *bls.G1Affine
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		g1x, _ := bls.RandFQ(r)
		g1y, _ := bls.RandFQ(r)
		g1z, _ := bls.RandFQ(r)
		g2x, _ := bls.RandFQ(r)
		g2y, _ := bls.RandFQ(r)
		inData[i] = addData{
			g1: bls.NewG1Projective(g1x, g1y, g1z),
			g2: bls.NewG1Affine(g2x, g2y),
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].g1.AddAffine(inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}
