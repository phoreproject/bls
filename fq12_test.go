package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/phoreproject/bls"
)

func TestFQ12MulBy014(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c0, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		c1, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		c5, err := bls.RandFQ2(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		a, err := bls.RandFQ12(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		b := a.Copy()
		b.MulAssign(bls.NewFQ12(
			bls.NewFQ6(c0, c1, bls.FQ2Zero),
			bls.NewFQ6(bls.FQ2Zero, c5, bls.FQ2Zero),
		))
		a.MulBy014Assign(c0, c1, c5)

		if !a.Equals(b) {
			t.Fatal("MulBy014 is broken.")
		}
	}
}

func TestFQ12Squaring(t *testing.T) {
	x := NewXORShift(2)
	a, _ := bls.RandFQ12(x)

	b := a.Copy()
	a.SquareAssign()
	b.MulAssign(b)

	if !a.Equals(b) {
		t.Fatal("squaring didn't work")
	}
}

func TestFQ12RandomFrobenius(t *testing.T) {
	x := NewXORShift(2)

	for i := 0; i < 10; i++ {
		for j := 0; j < 14; j++ {
			a, _ := bls.RandFQ12(x)
			b := a.Copy()

			for k := 0; k < j; k++ {
				a = a.Exp(bls.QFieldModulus)
			}
			b.FrobeniusMapAssign(uint8(j))

			if !a.Equals(b) {
				t.Fatal("frobenius map does not match exponent")
			}
		}
	}
}

func TestFQ12RandomMultiplication(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b, _ := bls.RandFQ12(x)
		c, _ := bls.RandFQ12(x)

		t0 := a.Copy()
		t0.MulAssign(b)
		t0.MulAssign(c)

		t1 := a.Copy()
		t1.MulAssign(c)
		t1.MulAssign(b)

		t2 := b.Copy()
		t2.MulAssign(c)
		t2.MulAssign(a)

		if !t0.Equals(t1) {
			t.Fatal("expected (a*b)*c == (a*c)*b")
		}

		if !t1.Equals(t2) {
			t.Fatal("expected (a*c)*b == (b*c)*a")
		}
	}
}

func TestFQ12RandomAddition(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b, _ := bls.RandFQ12(x)
		c, _ := bls.RandFQ12(x)

		t0 := a.Copy()
		t0.AddAssign(b)
		t0.AddAssign(c)

		t1 := a.Copy()
		t1.AddAssign(c)
		t1.AddAssign(b)

		t2 := b.Copy()
		t2.AddAssign(c)
		t2.AddAssign(a)

		if !t0.Equals(t1) {
			t.Fatal("expected (a+b)+c == (a+c)+b")
		}

		if !t1.Equals(t2) {
			t.Fatal("expected (a+c)+b == (b+c)+a")
		}
	}
}

func TestFQ12RandomSubtraction(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b, _ := bls.RandFQ12(x)

		t0 := a.Copy()
		t0.SubAssign(b)

		t1 := b.Copy()
		t1.SubAssign(a)

		t2 := t0.Copy()
		t2.AddAssign(t1)

		if !t2.IsZero() {
			t.Fatal("expected (a - b) + (b - a) = 0")
		}
	}
}

func TestFQ12RandomNegation(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b := a.Copy()

		b.NegAssign()
		b.AddAssign(a)

		if !b.IsZero() {
			t.Fatal("expected (a + ~a) = 0")
		}
	}
}

func TestFQ12RandomDoubling(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b := a.Copy()

		a.AddAssign(b)
		b.DoubleAssign()

		if !a.Equals(b) {
			t.Fatal("expected 2a = 2a")
		}
	}
}

func TestFQ12RandomSquaring(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b := a.Copy()

		a.MulAssign(b)
		b.SquareAssign()

		if !a.Equals(b) {
			t.Fatal("expected a^2 = a^2")
		}
	}
}

func TestFQ12RandomInversion(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b := a.Copy()

		b.InverseAssign()

		a.MulAssign(b)

		if !a.Equals(bls.FQ12One) {
			t.Fatal("expected a * a^-1 = 1")
		}
	}
}

func TestFQ12RandomExpansion(t *testing.T) {
	x := NewXORShift(3)

	for i := 0; i < 10000; i++ {
		a, _ := bls.RandFQ12(x)
		b, _ := bls.RandFQ12(x)
		c, _ := bls.RandFQ12(x)
		d, _ := bls.RandFQ12(x)

		t0 := a.Copy()
		t0.AddAssign(b)
		t1 := c.Copy()
		t1.AddAssign(d)
		t0.MulAssign(t1)

		t2 := a.Copy()
		t2.MulAssign(c)
		t3 := b.Copy()
		t3.MulAssign(c)
		t4 := a.Copy()
		t4.MulAssign(d)
		t5 := b.Copy()
		t5.MulAssign(d)

		t2.AddAssign(t3)
		t2.AddAssign(t4)
		t2.AddAssign(t5)

		if !t2.Equals(t0) {
			t.Fatal("(a + b)(c + d) should = (a*c + b*c + a*d + b*d)")
		}
	}
}

func BenchmarkFQ12MulAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
		f2 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
		f2, _ := bls.RandFQ12(r)
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

func BenchmarkFQ12SquareAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
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

func BenchmarkFQ12InverseAssign(b *testing.B) {
	type addData struct {
		f1 *bls.FQ12
	}

	r := NewXORShift(1)
	inData := [g1MulAssignSamples]addData{}
	for i := 0; i < g1MulAssignSamples; i++ {
		f1, _ := bls.RandFQ12(r)
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
