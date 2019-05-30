package g2pubs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/phoreproject/bls"
	"github.com/phoreproject/bls/g2pubs"
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
	for i := 0; i < loopCount; i++ {
		priv, _ := g2pubs.RandKey(r)
		pub := g2pubs.PrivToPub(priv)
		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		sig := g2pubs.Sign(msg, priv)
		if !g2pubs.Verify(msg, pub, sig) {
			return errors.New("sig did not verify")
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessage(loopCount int) error {
	r := NewXORShift(2)
	pubkeys := make([]*g2pubs.PublicKey, 0, 1000)
	sigs := make([]*g2pubs.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := g2pubs.RandKey(r)
		pub := g2pubs.PrivToPub(priv)
		sig := g2pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		if i < 10 || i > (loopCount-5) {
			newSig := g2pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregateCommon(pubkeys, msg) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessageMissingSig(loopCount int) error {
	r := NewXORShift(3)
	skippedSig := loopCount / 2
	pubkeys := make([]*g2pubs.PublicKey, 0, 1000)
	sigs := make([]*g2pubs.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := g2pubs.RandKey(r)
		pub := g2pubs.PrivToPub(priv)
		sig := g2pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		if i != skippedSig {
			sigs = append(sigs, sig)
		}
		if i < 10 || i > (loopCount-5) {
			newSig := g2pubs.AggregateSignatures(sigs)
			if newSig.VerifyAggregateCommon(pubkeys, msg) != (i < skippedSig) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func AggregateSignatures(loopCount int) error {
	r := NewXORShift(4)
	pubkeys := make([]*g2pubs.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := make([]*g2pubs.Signature, 0, 1000)
	for i := 0; i < loopCount; i++ {
		priv, _ := g2pubs.RandKey(r)
		pub := g2pubs.PrivToPub(priv)
		msg := []byte(fmt.Sprintf(">16 character identical message %d", i))
		sig := g2pubs.Sign(msg, priv)
		pubkeys = append(pubkeys, pub)
		msgs = append(msgs, msg)
		sigs = append(sigs, sig)

		if i < 10 || i > (loopCount-5) {
			newSig := g2pubs.AggregateSignatures(sigs)
			if !newSig.VerifyAggregate(pubkeys, msgs) {
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

func TestSignVerifyAggregateCommon(t *testing.T) {
	err := SignVerifyAggregateCommonMessage(10)
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

func TestAggregateSignaturesDuplicatedMessages(t *testing.T) {
	r := NewXORShift(5)

	pubkeys := make([]*g2pubs.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := g2pubs.NewAggregateSignature()

	key, _ := g2pubs.RandKey(r)
	pub := g2pubs.PrivToPub(key)
	message := []byte(">16 char first message")
	sig := g2pubs.Sign(message, key)
	pubkeys = append(pubkeys, pub)
	msgs = append(msgs, message)
	sigs.Aggregate(sig)

	if !sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature does not verify")
	}

	key2, _ := g2pubs.RandKey(r)
	pub2 := g2pubs.PrivToPub(key2)
	message2 := []byte(">16 char second message")
	sig2 := g2pubs.Sign(message2, key2)
	pubkeys = append(pubkeys, pub2)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig2)

	if !sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature does not verify")
	}

	key3, _ := g2pubs.RandKey(r)
	pub3 := g2pubs.PrivToPub(key3)
	sig3 := g2pubs.Sign(message2, key3)
	pubkeys = append(pubkeys, pub3)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig3)

	if sigs.VerifyAggregate(pubkeys, msgs) {
		t.Fatal("signature verifies with duplicate message")
	}
}

func TestAggregateSigsSeparate(t *testing.T) {
	x := NewXORShift(20)
	priv1, _ := g2pubs.RandKey(x)
	priv2, _ := g2pubs.RandKey(x)
	priv3, _ := g2pubs.RandKey(x)

	pub1 := g2pubs.PrivToPub(priv1)
	pub2 := g2pubs.PrivToPub(priv2)
	pub3 := g2pubs.PrivToPub(priv3)

	msg := []byte("test 1")
	sig1 := g2pubs.Sign(msg, priv1)
	sig2 := g2pubs.Sign(msg, priv2)
	sig3 := g2pubs.Sign(msg, priv3)

	aggSigs := g2pubs.AggregateSignatures([]*g2pubs.Signature{sig1, sig2, sig3})

	aggPubs := g2pubs.NewAggregatePubkey()
	aggPubs.Aggregate(pub1)
	aggPubs.Aggregate(pub2)
	aggPubs.Aggregate(pub3)

	valid := g2pubs.Verify(msg, aggPubs, aggSigs)
	if !valid {
		t.Fatal("expected aggregate signature to be valid")
	}
}

func BenchmarkBLSAggregateSignature(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := g2pubs.RandKey(r)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g2pubs.Sign(msg, priv)

	s := g2pubs.NewAggregateSignature()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Aggregate(sig)
	}
}

func BenchmarkBLSSign(b *testing.B) {
	r := NewXORShift(5)
	privs := make([]*g2pubs.SecretKey, b.N)
	for i := range privs {
		privs[i], _ = g2pubs.RandKey(r)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		g2pubs.Sign(msg, privs[i])
		// if !g2pubs.Verify(msg, pub, sig) {
		// 	return errors.New("sig did not verify")
		// }
	}
}

func BenchmarkBLSVerify(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := g2pubs.RandKey(r)
	pub := g2pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g2pubs.Sign(msg, priv)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g2pubs.Verify(msg, pub, sig)
	}
}

func TestSignatureSerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := g2pubs.RandKey(r)
	pub := g2pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g2pubs.Sign(msg, priv)

	if !g2pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization")
	}

	sigSer := sig.Serialize()
	sigDeser, err := g2pubs.DeserializeSignature(sigSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g2pubs.Verify(msg, pub, sigDeser) {
		t.Fatal("message did not verify after serialization/deserialization")
	}
}

func TestPubkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := g2pubs.RandKey(r)
	pub := g2pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g2pubs.Sign(msg, priv)

	if !g2pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization of pubkey")
	}

	pubSer := pub.Serialize()
	pubDeser, err := g2pubs.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g2pubs.Verify(msg, pubDeser, sig) {
		t.Fatal("message did not verify after serialization/deserialization of pubkey")
	}
}

func TestSecretkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(3)
	priv, _ := g2pubs.RandKey(r)
	privSer := priv.Serialize()
	privNew := g2pubs.DeserializeSecretKey(privSer)
	pub := g2pubs.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := g2pubs.Sign(msg, privNew)

	if !g2pubs.Verify(msg, pub, sig) {
		t.Fatal("message did not verify before serialization/deserialization of secret")
	}

	pubSer := pub.Serialize()
	pubDeser, err := g2pubs.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !g2pubs.Verify(msg, pubDeser, sig) {
		t.Fatal("message did not verify after serialization/deserialization of secret")
	}
}

func TestDeriveSecretKey(t *testing.T) {
	var secKeyIn [32]byte
	copy(secKeyIn[:], []byte("11223344556677889900112233445566"))
	k := g2pubs.DeriveSecretKey(secKeyIn)

	expectedElement, _ := bls.FRReprFromString("414e2c2a330cf94edb70e1c88efa851e80fe5eb14ff08fe5b7e588b4fe9899e4", 16)
	expectedFRElement := bls.FRReprToFR(expectedElement)

	if !expectedFRElement.Equals(k.GetFRElement()) {
		t.Fatal("expected secret key to match")
	}
}
