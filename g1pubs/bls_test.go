package g1pubs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/phoreproject/bls"
	"github.com/phoreproject/bls/g1pubs"
)

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

func SignVerify(loopCount int) error {
	r := NewXORShift(1)
	for i := 0; i < 1; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		sig := g1pubs.Sign(msg, priv)
		if !g1pubs.Verify(msg, pub, sig) {
			return errors.New("sig did not verify")
		}
	}
	return nil
}

func SignVerifyWithDomain(loopCount int) error {
	r := NewXORShift(1)
	for i := 0; i < 1; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		var hashedMsg [32]byte
		copy(hashedMsg[:], []byte(fmt.Sprintf("Hello world! 16 characters %d", i)))
		sig := g1pubs.SignWithDomain(hashedMsg, priv, [8]byte{1})
		if !g1pubs.VerifyWithDomain(hashedMsg, pub, sig, [8]byte{1}) {
			return errors.New("sig did not verify")
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessage(loopCount int) error {
	r := NewXORShift(2)
	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	sigs := make([]*g1pubs.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		sig := g1pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		if i < 10 || i > (loopCount-5) {
			newSig := g1pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregateCommon(pubkeys, msg) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessageWithDomain(loopCount int) error {
	var hashedMessage [32]byte
	domain := [8]byte{100}
	r := NewXORShift(2)
	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	sigs := make([]*g1pubs.Signature, 0, 1000)
	copy(hashedMessage[:], []byte(">16 character identical message"))

	for i := 0; i < loopCount; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		sig := g1pubs.SignWithDomain(hashedMessage, priv, domain)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		if i < 10 || i > (loopCount-5) {
			newSig := g1pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregateCommonWithDomain(pubkeys, hashedMessage, domain) {
				return fmt.Errorf("sig did not verify for loop %d", i)
			}
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessageMissingSig(loopCount int) error {
	r := NewXORShift(3)
	skippedSig := loopCount / 2
	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	sigs := make([]*g1pubs.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		sig := g1pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		if i != skippedSig {
			sigs = append(sigs, sig)
		}
		if i < 10 || i > (loopCount-5) {
			newSig := g1pubs.AggregateSignatures(sigs)
			if newSig.VerifyAggregateCommon(pubkeys, msg) != (i < skippedSig) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func AggregateSignatures(loopCount int) error {
	r := NewXORShift(4)
	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := make([]*g1pubs.Signature, 0, 1000)
	for i := 0; i < loopCount; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		msg := []byte(fmt.Sprintf(">16 character identical message %d", i))
		sig := g1pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		msgs = append(msgs, msg)
		sigs = append(sigs, sig)

		if i < 10 || i > (loopCount-5) {
			newSig := g1pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregate(pubkeys, msgs) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func AggregateSignaturesWithDomain(loopCount int) error {
	var msg [32]byte
	r := NewXORShift(4)
	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	msgs := make([][32]byte, 0, 1000)
	sigs := make([]*g1pubs.Signature, 0, 1000)
	domain := [8]byte{0xFF}
	for i := 0; i < loopCount; i++ {
		priv, _ := g1pubs.RandKey(r)
		pub := g1pubs.PrivToPub(priv)
		copy(msg[:], []byte(fmt.Sprintf(">16 character identical message %d", i)))
		sig := g1pubs.SignWithDomain(msg, priv, domain)
		pubkeys = append(pubkeys, pub)
		msgs = append(msgs, msg)
		sigs = append(sigs, sig)

		if i < 10 || i > (loopCount-5) {
			newSig := g1pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregateWithDomain(pubkeys, msgs, domain) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func TestSignVerify(t *testing.T) {
	err := SignVerify(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignVerifyWithDomain(t *testing.T) {
	err := SignVerifyWithDomain(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignVerifyAggregateCommon(t *testing.T) {
	err := SignVerifyAggregateCommonMessage(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignVerifyAggregateCommonWithDomain(t *testing.T) {
	err := SignVerifyAggregateCommonMessageWithDomain(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignVerifyAggregateCommonMissingSig(t *testing.T) {
	err := SignVerifyAggregateCommonMessageMissingSig(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignAggregateSigs(t *testing.T) {
	err := AggregateSignatures(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignAggregateSigsWithDomain(t *testing.T) {
	err := AggregateSignaturesWithDomain(10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregateSignaturesDuplicatedMessages(t *testing.T) {
	r := NewXORShift(5)

	pubkeys := make([]*g1pubs.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := g1pubs.NewAggregateSignature()

	key, _ := g1pubs.RandKey(r)
	pub := g1pubs.PrivToPub(key)
	message := []byte(">16 char first message")
	sig := g1pubs.Sign(message, key)
	pubkeys = append(pubkeys, pub)
	msgs = append(msgs, message)
	sigs.Aggregate(sig)

	if !sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature does not verify")
	}

	key2, _ := g1pubs.RandKey(r)
	pub2 := g1pubs.PrivToPub(key2)
	message2 := []byte(">16 char second message")
	sig2 := g1pubs.Sign(message2, key2)
	pubkeys = append(pubkeys, pub2)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig2)

	if !sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature does not verify")
	}

	key3, _ := g1pubs.RandKey(r)
	pub3 := g1pubs.PrivToPub(key3)
	sig3 := g1pubs.Sign(message2, key3)
	pubkeys = append(pubkeys, pub3)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig3)

	if sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature verifies with duplicate message")
	}
}

func TestAggregateSigsSeparate(t *testing.T) {
	x := NewXORShift(20)
	priv1, _ := g1pubs.RandKey(x)
	priv2, _ := g1pubs.RandKey(x)
	priv3, _ := g1pubs.RandKey(x)

	pub1 := g1pubs.PrivToPub(priv1)
	pub2 := g1pubs.PrivToPub(priv2)
	pub3 := g1pubs.PrivToPub(priv3)

	msg := []byte("test 1")
	sig1 := g1pubs.Sign(msg, priv1)
	sig2 := g1pubs.Sign(msg, priv2)
	sig3 := g1pubs.Sign(msg, priv3)

	aggSigs := g1pubs.AggregateSignatures([]*g1pubs.Signature{sig1, sig2, sig3})

	aggPubs := g1pubs.NewAggregatePubkey()
	aggPubs.Aggregate(pub1)
	aggPubs.Aggregate(pub2)
	aggPubs.Aggregate(pub3)

	valid := g1pubs.Verify(msg, aggPubs, aggSigs)
	if !valid {
		t.Fatal("expected aggregate signature to be valid")
	}
}

func BenchmarkBLSAggregateSignature(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := g1pubs.RandKey(r)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g1pubs.Sign(msg, priv)

	s := g1pubs.NewAggregateSignature()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Aggregate(sig)
	}
}

func BenchmarkBLSSign(b *testing.B) {
	r := NewXORShift(5)
	privs := make([]*g1pubs.SecretKey, b.N)
	for i := range privs {
		privs[i], _ = g1pubs.RandKey(r)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		g1pubs.Sign(msg, privs[i])
		// if !g1pubs.Verify(msg, pub, sig) {
		// 	return errors.New("sig did not verify")
		// }
	}
}

func BenchmarkBLSVerify(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := g1pubs.RandKey(r)
	pub := g1pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g1pubs.Sign(msg, priv)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g1pubs.Verify(msg, pub, sig)
	}
}

func TestSignatureSerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := g1pubs.RandKey(r)
	pub := g1pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g1pubs.Sign(msg, priv)

	if !g1pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization")
	}

	sigSer := sig.Serialize()
	sigDeser, err := g1pubs.DeserializeSignature(sigSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g1pubs.Verify(msg, pub, sigDeser) {
		t.Fatal("message did not verify after serialization/deserialization")
	}
}

func TestPubkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := g1pubs.RandKey(r)
	pub := g1pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g1pubs.Sign(msg, priv)

	if !g1pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization of pubkey")
	}

	pubSer := pub.Serialize()
	pubDeser, err := g1pubs.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g1pubs.Verify(msg, pubDeser, sig) {
		t.Fatal("message did not verify after serialization/deserialization of pubkey")
	}
}

func TestSecretkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(3)
	priv, _ := g1pubs.RandKey(r)
	privSer := priv.Serialize()
	privNew := g1pubs.DeserializeSecretKey(privSer)
	pub := g1pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g1pubs.Sign(msg, privNew)

	if !g1pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization of secret")
	}

	pubSer := pub.Serialize()
	pubDeser, err := g1pubs.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g1pubs.Verify(msg, pubDeser, sig) {
		t.Fatal("message did not verify after serialization/deserialization of secret")
	}
}

func TestDeriveSecretKey(t *testing.T) {
	var secKeyIn [32]byte
	copy(secKeyIn[:], []byte("11223344556677889900112233445566"))
	k := g1pubs.DeriveSecretKey(secKeyIn)

	expectedElement, _ := bls.FRReprFromString("414e2c2a330cf94edb70e1c88efa851e80fe5eb14ff08fe5b7e588b4fe9899e4", 16)
	expectedFRElement := bls.FRReprToFR(expectedElement)

	if !expectedFRElement.Equals(k.GetFRElement()) {
		t.Fatal("expected secret key to match")
	}
}
