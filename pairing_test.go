package bls_test

import (
	"testing"

	"github.com/phoreproject/bls"
)

var c000, _ = bls.FQReprFromString("2819105605953691245277803056322684086884703000473961065716485506033588504203831029066448642358042597501014294104502", 10)
var c001, _ = bls.FQReprFromString("1323968232986996742571315206151405965104242542339680722164220900812303524334628370163366153839984196298685227734799", 10)
var c010, _ = bls.FQReprFromString("2987335049721312504428602988447616328830341722376962214011674875969052835043875658579425548512925634040144704192135", 10)
var c011, _ = bls.FQReprFromString("3879723582452552452538684314479081967502111497413076598816163759028842927668327542875108457755966417881797966271311", 10)
var c020, _ = bls.FQReprFromString("261508182517997003171385743374653339186059518494239543139839025878870012614975302676296704930880982238308326681253", 10)
var c021, _ = bls.FQReprFromString("231488992246460459663813598342448669854473942105054381511346786719005883340876032043606739070883099647773793170614", 10)
var c100, _ = bls.FQReprFromString("3993582095516422658773669068931361134188738159766715576187490305611759126554796569868053818105850661142222948198557", 10)
var c101, _ = bls.FQReprFromString("1074773511698422344502264006159859710502164045911412750831641680783012525555872467108249271286757399121183508900634", 10)
var c110, _ = bls.FQReprFromString("2727588299083545686739024317998512740561167011046940249988557419323068809019137624943703910267790601287073339193943", 10)
var c111, _ = bls.FQReprFromString("493643299814437640914745677854369670041080344349607504656543355799077485536288866009245028091988146107059514546594", 10)
var c120, _ = bls.FQReprFromString("734401332196641441839439105942623141234148957972407782257355060229193854324927417865401895596108124443575283868655", 10)
var c121, _ = bls.FQReprFromString("2348330098288556420918672502923664952620152483128593484301759394583320358354186482723629999370241674973832318248497", 10)

func TestResultAgainstRelic(t *testing.T) {
	out := bls.Pairing(bls.G1ProjectiveOne.Copy(), bls.G2ProjectiveOne.Copy())
	expected := bls.NewFQ12(
		bls.NewFQ6(
			bls.NewFQ2(
				bls.FQReprToFQ(c000),
				bls.FQReprToFQ(c001),
			),
			bls.NewFQ2(
				bls.FQReprToFQ(c010),
				bls.FQReprToFQ(c011),
			),
			bls.NewFQ2(
				bls.FQReprToFQ(c020),
				bls.FQReprToFQ(c021),
			),
		),
		bls.NewFQ6(
			bls.NewFQ2(
				bls.FQReprToFQ(c100),
				bls.FQReprToFQ(c101),
			),
			bls.NewFQ2(
				bls.FQReprToFQ(c110),
				bls.FQReprToFQ(c111),
			),
			bls.NewFQ2(
				bls.FQReprToFQ(c120),
				bls.FQReprToFQ(c121),
			),
		),
	)

	if !out.Equals(expected) {
		t.Fatal("pairing result is wrong")
	}
}

func BenchmarkG2Prepare(b *testing.B) {
	type addData struct {
		g2 *bls.G2Affine
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandG2(r)
		inData[i] = addData{
			g2: f1.ToAffine(),
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		bls.G2AffineToPrepared(inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkMillerLoop(b *testing.B) {
	type addData struct {
		p *bls.G1Affine
		q *bls.G2Prepared
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f2, _ := bls.RandG2(r)
		f1, _ := bls.RandG1(r)
		inData[i] = addData{
			q: bls.G2AffineToPrepared(f2.ToAffine()),
			p: f1.ToAffine(),
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		bls.MillerLoop([]bls.MillerLoopItem{{P: inData[count].p, Q: inData[count].q}})
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFinalExponentiation(b *testing.B) {
	r := NewXORShift(1)
	inData := [g1MulAssignSamples]*bls.FQ12{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f2, _ := bls.RandG2(r)
		f1, _ := bls.RandG1(r)
		inData[i] = bls.MillerLoop([]bls.MillerLoopItem{
			{
				Q: bls.G2AffineToPrepared(f2.ToAffine()),
				P: f1.ToAffine(),
			},
		})
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		bls.FinalExponentiation(inData[count])
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkPairing(b *testing.B) {
	type pairingData struct {
		g1 *bls.G1Projective
		g2 *bls.G2Projective
	}
	r := NewXORShift(1)
	inData := [g1MulAssignSamples]pairingData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f2, _ := bls.RandG2(r)
		f1, _ := bls.RandG1(r)
		inData[i] = pairingData{g1: f1, g2: f2}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		bls.Pairing(inData[count].g1, inData[count].g2)
		count = (count + 1) % g1MulAssignSamples
	}
}
