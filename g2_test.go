package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

func BenchmarkG2MulAssign(b *testing.B) {
	type mulData struct {
		g *bls.G2Projective
		f *bls.FR
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]mulData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		gx, _ := bls.RandFQ2(r)
		gy, _ := bls.RandFQ2(r)
		gz, _ := bls.RandFQ2(r)
		randFR, _ := bls.RandFR(r)
		inData[i] = mulData{
			g: bls.NewG2Projective(gx, gy, gz),
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

func BenchmarkG2AddAssign(b *testing.B) {
	type addData struct {
		g1 *bls.G2Projective
		g2 *bls.G2Projective
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		g1x, _ := bls.RandFQ2(r)
		g1y, _ := bls.RandFQ2(r)
		g1z, _ := bls.RandFQ2(r)
		g2x, _ := bls.RandFQ2(r)
		g2y, _ := bls.RandFQ2(r)
		g2z, _ := bls.RandFQ2(r)
		inData[i] = addData{
			g1: bls.NewG2Projective(g1x, g1y, g1z),
			g2: bls.NewG2Projective(g2x, g2y, g2z),
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].g1.Add(inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkG2AddAssignMixed(b *testing.B) {
	type addData struct {
		g1 *bls.G2Projective
		g2 *bls.G2Affine
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		g1x, _ := bls.RandFQ2(r)
		g1y, _ := bls.RandFQ2(r)
		g1z, _ := bls.RandFQ2(r)
		g2x, _ := bls.RandFQ2(r)
		g2y, _ := bls.RandFQ2(r)
		inData[i] = addData{
			g1: bls.NewG2Projective(g1x, g1y, g1z),
			g2: bls.NewG2Affine(g2x, g2y),
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].g1.AddAffine(inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}
