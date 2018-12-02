package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/phoreproject/bls"
)

var (
	bigZero = bls.NewFQRepr(0)
	bigOne  = bls.NewFQRepr(1)
	bigTwo  = bls.NewFQRepr(2)
)

func TestFQ2Ordering(t *testing.T) {
	a := bls.NewFQ2(bls.FQZero.Copy(), bls.FQZero.Copy())
	b := a.Copy()

	if a.Cmp(b) != 0 {
		t.Error("a != b after cloning a to b")
	}
	b.AddAssign(bls.NewFQ2(bls.FQOne, bls.FQZero))
	if a.Cmp(b) >= 0 {
		t.Error("a >= b after adding to b")
	}
	a.AddAssign(bls.NewFQ2(bls.FQOne, bls.FQZero))
	if a.Cmp(b) != 0 {
		t.Error("a != b after adding to a and b")
	}
	b.AddAssign(bls.NewFQ2(bls.FQZero, bls.FQOne))
	if a.Cmp(b) >= 0 {
		t.Error("a >= b after adding to b.c1")
	}
	a.AddAssign(bls.NewFQ2(bls.FQOne, bls.FQZero))
	if a.Cmp(b) >= 0 {
		t.Error("c0 is taking precedent over c1")
	}
	a.AddAssign(bls.NewFQ2(bls.FQZero, bls.FQOne))
	if a.Cmp(b) <= 0 {
		t.Error("FQ2(2, 1) <= FQ2(1, 1)")
	}
	b.AddAssign(bls.NewFQ2(bls.FQOne, bls.FQZero))
	if a.Cmp(b) != 0 {
		t.Error("FQ2(2, 1) != FQ2(2, 1)")
	}
}

func TestFQ2Basics(t *testing.T) {
	f := bls.NewFQ2(bls.FQZero, bls.FQZero)
	if !f.Equals(bls.FQ2Zero) {
		t.Error("FQ2Zero != FQ2(0, 0)")
	}

	f = bls.NewFQ2(bls.FQOne, bls.FQZero)
	if !f.Equals(bls.FQ2One) {
		t.Error("FQ2One != FQ2(1, 0)")
	}

	if bls.FQ2One.IsZero() {
		t.Error("FQ2One.IsZero() == true")
	}
	if !bls.FQ2Zero.IsZero() {
		t.Error("FQ2Zero.IsZero() != true")
	}
	f = bls.NewFQ2(bls.FQZero, bls.FQOne)
	if f.IsZero() {
		t.Error("FQ2(0, 1).IsZero() == true")
	}
}

func TestFQ2Squaring(t *testing.T) {
	a := bls.NewFQ2(bls.FQOne, bls.FQOne)
	a.SquareAssign()
	expected := bls.NewFQ2(bls.FQZero, bls.FQReprToFQ(bigTwo))
	if !a.Equals(expected) {
		t.Error("FQ(1, 1).Square() != FQ(0, 2)")
	}

	a = bls.NewFQ2(bls.FQZero, bls.FQOne)
	a.SquareAssign()
	neg1 := bls.FQOne.Copy()
	neg1.NegAssign()
	expected = bls.NewFQ2(neg1, bls.FQZero)
	if !a.Equals(expected) {
		t.Error("FQ(0, 1).Square() != FQ(-1, 0)")
	}

	a0, _ := bls.FQReprFromString("7080c5fa1d8e04241b76dcc1c3fbe5ef7f295a94e58ae7c90e34aab6fb6a6bd4eef5c946536f6029c2c6309bbf8b598", 16)
	a1, _ := bls.FQReprFromString("10d1615e75250a21fc58a7b7be815407bfb99020604137a0dac5a4c911a4353e6ad3291177c8c7e538f473b3c870a4ab", 16)
	a = bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	a.SquareAssign()
	expected0, _ := bls.FQReprFromString("7eac81369c433614cf17b5893c3d327cb674157618da1760dc46ab8fad67ae0b9f2a66eae1073baf262c28c538bcf68", 16)
	expected1, _ := bls.FQReprFromString("1542a61c8a8db994739c983042779a6538d0d7275a9689e1e75138bce4cec7aaa23eb7e12dd54d98c1579cf58e980cf8", 16)
	expected = bls.NewFQ2(bls.FQReprToFQ(expected0), bls.FQReprToFQ(expected1))
	if !a.Equals(expected) {
		t.Error("adding FQ2s together not giving expected result")
	}
}

func TestFQ2Mul(t *testing.T) {
	a0, _ := bls.FQReprFromString("787155314249811040945591987890197698923521563000754191852630526325910580973833767165585004431929413598958418534147", 10)
	a1, _ := bls.FQReprFromString("3701300508704531706796752084315366662585685126320813849768310480839238335689933757701592694115590705658584012352975", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	b0, _ := bls.FQReprFromString("428812050294605776242613635112226565923171104978712194146076278968219164256021918826122932548851868482795586147198", 10)
	b1, _ := bls.FQReprFromString("3983398787459011889399209724009438146317670502764434060579092960070359230210191684510008045453936255762435213464276", 10)
	b := bls.NewFQ2(bls.FQReprToFQ(b0), bls.FQReprToFQ(b1))
	o0, _ := bls.FQReprFromString("3589136112018482378111625145859346905296349062557304864057134187313531783880863512355763193523645312504462844151780", 10)
	o1, _ := bls.FQReprFromString("297668880544309439220785639618566999376775158493905442186353189401526493115249312794616596981510685741636811493754", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.MulAssign(b)
	if !a.Equals(o) {
		t.Error("FQ2 Mul not working properly")
	}
}

func TestFQ2Inverse(t *testing.T) {
	invZero := bls.FQ2Zero.Copy()
	if invZero.InverseAssign() {
		t.Error("inverse of zero is returning a non-nil value")
	}

	a0, _ := bls.FQReprFromString("787155314249811040945591987890197698923521563000754191852630526325910580973833767165585004431929413598958418534147", 10)
	a1, _ := bls.FQReprFromString("3701300508704531706796752084315366662585685126320813849768310480839238335689933757701592694115590705658584012352975", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	o0, _ := bls.FQReprFromString("2973628342543885151746258559273164780300946315949957485732289400047711856966246507651218645021769655299729012680084", 10)
	o1, _ := bls.FQReprFromString("2498621830671500058873354492670308521035465162737562150317014095994530700756690595853342128497677859200322409409716", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.InverseAssign()
	if !a.Equals(o) {
		t.Error("FQ2 Inv not working properly")
	}
}

func TestFQ2Addition(t *testing.T) {
	a0, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	a1, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	b0, _ := bls.FQReprFromString("3048378325167389389927636537438764659078261027448144976696012669013445082263630955837031612987742627153430205370098", 10)
	b1, _ := bls.FQReprFromString("2745732935355827219217491212445536316258315122004495924168966599922007114897473294472413581011479900113332636006841", 10)
	b := bls.NewFQ2(bls.FQReprToFQ(b0), bls.FQReprToFQ(b1))
	o0, _ := bls.FQReprFromString("3197977163940080762958009789663025044787365462715437957461204574073501724672681252378504259284666769425426845208249", 10)
	o1, _ := bls.FQReprFromString("2926742669024788240740727831667279901593764337752150720908161880718213702916397600862029174885519475199268964768404", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.AddAssign(b)
	if !a.Equals(o) {
		t.Error("FQ2 add not working properly")
	}
}

func TestFQ2Subtraction(t *testing.T) {
	a0, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	a1, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	b0, _ := bls.FQReprFromString("3048378325167389389927636537438764659078261027448144976696012669013445082263630955837031612987742627153430205370098", 10)
	b1, _ := bls.FQReprFromString("2745732935355827219217491212445536316258315122004495924168966599922007114897473294472413581011479900113332636006841", 10)
	b := bls.NewFQ2(bls.FQReprToFQ(b0), bls.FQReprToFQ(b1))
	o0, _ := bls.FQReprFromString("1103630068826969376520526540521399883187726227758155889401237372170643210636257205147128662438197179156460707027840", 10)
	o1, _ := bls.FQReprFromString("1437686353534801195723535232512111425634016913682166757902286816998231123612288876359889641991575339010497965314509", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.SubAssign(b)
	if !a.Equals(o) {
		t.Error("FQ2 sub not working properly")
	}
}

func TestFQ2Negation(t *testing.T) {
	a0, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	a1, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	o0, _ := bls.FQReprFromString("3852810716448976020387416573511643770847778384671714904566866231063975008081787567901214982832091521765897632721636", 10)
	o1, _ := bls.FQReprFromString("3821399821552706371894553206514160571221433604191353088592862855327825062471913558053072035254976088951957943798224", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.NegAssign()
	if !a.Equals(o) {
		t.Error("FQ2 negation not working properly")
	}
}

func TestFQ2Doubling(t *testing.T) {
	a0, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	a1, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	o0, _ := bls.FQReprFromString("299197677545382746060746504448520771418208870534585961530383810120113284818100593082945292593848284543993279676302", 10)
	o1, _ := bls.FQReprFromString("362019467337922043046473238443487170670898431495309593478390561592413176037848612779231187748079150171872657523126", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	a.DoubleAssign()
	if !a.Equals(o) {
		t.Error("FQ2 double not working properly")
	}
}

func TestFQ2FrobeniusMap(t *testing.T) {

	a0, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	a1, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	o00, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	o01, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	o0 := bls.NewFQ2(bls.FQReprToFQ(o00), bls.FQReprToFQ(o01))
	a.FrobeniusMapAssign(0)
	if !a.Equals(o0) {
		t.Error("FQ2 frobenius map not working properly")
	}
	o10, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	o11, _ := bls.FQReprFromString("3821399821552706371894553206514160571221433604191353088592862855327825062471913558053072035254976088951957943798224", 10)
	o1 := bls.NewFQ2(bls.FQReprToFQ(o10), bls.FQReprToFQ(o11))
	a.FrobeniusMapAssign(1)
	if !a.Equals(o1) {
		t.Error("FQ2 frobenius map not working properly")
	}
	o20, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	o21, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	o2 := bls.NewFQ2(bls.FQReprToFQ(o20), bls.FQReprToFQ(o21))
	a.FrobeniusMapAssign(1)
	if !a.Equals(o2) {
		t.Error("FQ2 frobenius map not working properly")
	}
	o30, _ := bls.FQReprFromString("149598838772691373030373252224260385709104435267292980765191905060056642409050296541472646296924142271996639838151", 10)
	o31, _ := bls.FQReprFromString("181009733668961021523236619221743585335449215747654796739195280796206588018924306389615593874039575085936328761563", 10)
	o3 := bls.NewFQ2(bls.FQReprToFQ(o30), bls.FQReprToFQ(o31))
	a.FrobeniusMapAssign(2)
	if !a.Equals(o3) {
		t.Error("FQ2 frobenius map not working properly")
	}
}

func TestFQ2Sqrt(t *testing.T) {
	a0, _ := bls.FQReprFromString("1199141494453035944462437764039889211157529566989597369793908631238550963749870321413115974274142904836460239708711", 10)
	a1, _ := bls.FQReprFromString("2275726519012059234888434151906712303065680739818771469342471201468687921783347469538129328635084833206431018046147", 10)
	a := bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQReprToFQ(a1))
	o0, _ := bls.FQReprFromString("2610880034315515246135455194558180576521075490930523809216746350127301781195091608845140759868139294094360490170565", 10)
	o1, _ := bls.FQReprFromString("1361064970384811010963395494397997021963412519774828006120608106648478662746078959841022528364671568910239381710674", 10)
	o := bls.NewFQ2(bls.FQReprToFQ(o0), bls.FQReprToFQ(o1))
	if !a.Sqrt().Equals(o) {
		t.Error("FQ2 sqrt not working properly")
	}
	a0, _ = bls.FQReprFromString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015661931409599199851", 10)
	a = bls.NewFQ2(bls.FQReprToFQ(a0), bls.FQZero)
	o1, _ = bls.FQReprFromString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894226663331", 10)
	o = bls.NewFQ2(bls.FQZero, bls.FQReprToFQ(o1))
	if !a.Sqrt().Equals(o) {
		t.Error("FQ2 sqrt not working properly")
	}
}

func TestFQ2Legendre(t *testing.T) {
	if bls.LegendreZero != bls.FQ2Zero.Legendre() {
		t.Error("legendre of zero field element does not equal LegendreZero")
	}

	m1 := bls.FQ2One.Copy()
	m1.NegAssign()
	if bls.LegendreQuadraticResidue != m1.Legendre() {
		t.Error("sqrt(-1) is not quadratic residue")
	}
	m1.MultiplyByNonresidueAssign()
	if bls.LegendreQuadraticNonResidue != m1.Legendre() {
		t.Error("1.Neg().MulByNonresidue() is quadratic residue")
	}
}

func TestFQ2MulNonresidue(t *testing.T) {
	nqr := bls.NewFQ2(bls.FQOne, bls.FQOne)

	for i := 0; i < 1000; i++ {
		a, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		a.MultiplyByNonresidueAssign()
		b.MulAssign(nqr)

		if !a.Equals(b) {
			t.Error("a != b")
		}
	}
}

func BenchmarkFQ2AddAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ2
		f2 *bls.FQ2
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ2(r)
		f2, _ := bls.RandFQ2(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.AddAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ2SubAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ2
		f2 *bls.FQ2
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ2(r)
		f2, _ := bls.RandFQ2(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.SubAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ2MulAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ2
		f2 *bls.FQ2
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ2(r)
		f2, _ := bls.RandFQ2(r)
		inData[i] = addData{
			f1: f1,
			f2: f2,
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.MulAssign(inData[count].f2)
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ2SquareAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ2
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ2(r)
		inData[i] = addData{
			f1: f1,
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.SquareAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}

func BenchmarkFQ2InverseAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ2
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ2(r)
		inData[i] = addData{
			f1: f1,
		}
	}
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		inData[count].f1.InverseAssign()
		count = (count + 1) % g1MulAssignSamples
	}
}
