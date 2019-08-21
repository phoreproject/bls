package g1pubs_test

import (
	"crypto/rand"
	"fmt"
	"github.com/phoreproject/bls/g1pubs"
	"testing"
)

func ToBytes32(v []byte) (out [32]byte) {
	copy(out[:], v)
	return
}

func BenchmarkVerifyWithDomain(b *testing.B) {
	sk, err := g1pubs.RandKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	msg := ToBytes32([]byte("Some msg"))
	domain := [8]byte{42}
	sig := g1pubs.SignWithDomain(msg, sk, domain)

	pk := g1pubs.PrivToPub(sk)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g1pubs.VerifyWithDomain(msg, pk, sig, domain) {
			b.Fatal("could not verify sig")
		}
	}
}

func BenchmarkVerifyAggregateCommonWithDomain(b *testing.B) {
	sigN := 128
	msg := ToBytes32([]byte("Some message"))
	domain := [8]byte{42}

	var sigs []*g1pubs.Signature
	var pks []*g1pubs.PublicKey
	for i := 0; i < sigN; i++ {
		sk, err := g1pubs.RandKey(rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
		sigs = append(sigs, g1pubs.SignWithDomain(msg, sk, domain))
		pks = append(pks, g1pubs.PrivToPub(sk))
	}

	aggSig := g1pubs.AggregateSignatures(sigs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !aggSig.VerifyAggregateCommonWithDomain(pks, msg, domain) {
			b.Fatal("could not verify aggregate sig")
		}
	}
}

func BenchmarkVerifyAggregateMultipleWithDomain(b *testing.B) {
	sigN := 128
	var msg [32]byte
	var msgs [][32]byte
	domain := [8]byte{42}

	var sigs []*g1pubs.Signature
	var pks []*g1pubs.PublicKey
	for i := 0; i < sigN; i++ {
		sk, err := g1pubs.RandKey(rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
		copy(msg[:], fmt.Sprintf("Some message %d", i))
		sigs = append(sigs, g1pubs.SignWithDomain(msg, sk, domain))
		pks = append(pks, g1pubs.PrivToPub(sk))
		msgs = append(msgs, msg)
	}

	aggSig := g1pubs.AggregateSignatures(sigs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !aggSig.VerifyAggregateWithDomain(pks, msgs, domain) {
			b.Fatal("could not verify aggregate sig")
		}
	}
}
