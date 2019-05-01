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
		rhs.AddAssign(bls.BCoeff)

		y, success := rhs.Sqrt()

		if success {
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

		i++
		x.AddAssign(bls.FQOne)
	}
}

func TestG1DoublingCorrectness(t *testing.T) {
	p := bls.NewG1Projective(
		bls.FQReprToFQ(bls.FQRepr{0x47fd1f891d6e8bbf, 0x79a3b0448f31a2aa, 0x81f3339e5f9968f, 0x485e77d50a5df10d, 0x4c6fcac4b55fd479, 0x86ed4d9906fb064}),
		bls.FQReprToFQ(bls.FQRepr{0xd25ee6461538c65, 0x9f3bbb2ecd3719b9, 0xa06fd3f1e540910d, 0xcefca68333c35288, 0x570c8005f8573fa6, 0x152ca696fe034442}),
		bls.FQOne,
	)

	newP := p.Double()

	expectedP := bls.NewG1Affine(
		bls.FQReprToFQ(bls.FQRepr{0xf939ddfe0ead7018, 0x3b03942e732aecb, 0xce0e9c38fdb11851, 0x4b914c16687dcde0, 0x66c8baf177d20533, 0xaf960cff3d83833}),
		bls.FQReprToFQ(bls.FQRepr{0x3f0675695f5177a8, 0x2b6d82ae178a1ba0, 0x9096380dd8e51b11, 0x1771a65b60572f4e, 0x8b547c1313b27555, 0x135075589a687b1e}),
	)

	if !newP.ToAffine().Equals(expectedP) {
		t.Fatal("doubling is incorrect")
	}
}

func TestG1AdditionCorrectness(t *testing.T) {
	p1 := bls.NewG1Projective(
		bls.FQReprToFQ(bls.FQRepr{0x47fd1f891d6e8bbf, 0x79a3b0448f31a2aa, 0x81f3339e5f9968f, 0x485e77d50a5df10d, 0x4c6fcac4b55fd479, 0x86ed4d9906fb064}),
		bls.FQReprToFQ(bls.FQRepr{0xd25ee6461538c65, 0x9f3bbb2ecd3719b9, 0xa06fd3f1e540910d, 0xcefca68333c35288, 0x570c8005f8573fa6, 0x152ca696fe034442}),
		bls.FQOne,
	)

	p2 := bls.NewG1Projective(
		bls.FQReprToFQ(bls.FQRepr{0xeec78f3096213cbf, 0xa12beb1fea1056e6, 0xc286c0211c40dd54, 0x5f44314ec5e3fb03, 0x24e8538737c6e675, 0x8abd623a594fba8}),
		bls.FQReprToFQ(bls.FQRepr{0x6b0528f088bb7044, 0x2fdeb5c82917ff9e, 0x9a5181f2fac226ad, 0xd65104c6f95a872a, 0x1f2998a5a9c61253, 0xe74846154a9e44}),
		bls.FQOne,
	)

	newP := p1.Add(p2).ToAffine()

	expectedP := bls.NewG1Affine(
		bls.FQReprToFQ(bls.FQRepr{0x6dd3098f22235df, 0xe865d221c8090260, 0xeb96bb99fa50779f, 0xc4f9a52a428e23bb, 0xd178b28dd4f407ef, 0x17fb8905e9183c69}),
		bls.FQReprToFQ(bls.FQRepr{0xd0de9d65292b7710, 0xf6a05f2bcf1d9ca7, 0x1040e27012f20b64, 0xeec8d1a5b7466c58, 0x4bc362649dce6376, 0x430cbdc5455b00a}),
	)

	if !newP.Equals(expectedP) {
		t.Fatal("addition is incorrect")
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

const g1MulAssignSamples = 200

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
