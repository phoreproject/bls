package bls_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/phoreproject/bls"
)

func SignVerify(loopCount int) error {
	r := NewXORShift(1)
	for i := 0; i < loopCount; i++ {
		priv, _ := bls.RandKey(r)
		pub := bls.PrivToPub(priv)
		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		sig := bls.Sign(msg, priv, 0)
		if !bls.Verify(msg, pub, sig, 0) {
			return errors.New("sig did not verify")
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessage(loopCount int) error {
	r := NewXORShift(2)
	pubkeys := make([]*bls.PublicKey, 0, 1000)
	sigs := make([]*bls.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := bls.RandKey(r)
		pub := bls.PrivToPub(priv)
		sig := bls.Sign(msg, priv, 0)
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		if i < 10 || i > (loopCount-5) {
			newSig := bls.AggregateSignatures(sigs)
			if !newSig.VerifyAggregateCommon(pubkeys, msg, 0) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func SignVerifyAggregateCommonMessageMissingSig(loopCount int) error {
	r := NewXORShift(3)
	skippedSig := loopCount / 2
	pubkeys := make([]*bls.PublicKey, 0, 1000)
	sigs := make([]*bls.Signature, 0, 1000)
	msg := []byte(">16 character identical message")
	for i := 0; i < loopCount; i++ {
		priv, _ := bls.RandKey(r)
		pub := bls.PrivToPub(priv)
		sig := bls.Sign(msg, priv, 0)
		pubkeys = append(pubkeys, pub)
		if i != skippedSig {
			sigs = append(sigs, sig)
		}
		if i < 10 || i > (loopCount-5) {
			newSig := bls.AggregateSignatures(sigs)
			if newSig.VerifyAggregateCommon(pubkeys, msg, 0) != (i < skippedSig) {
				return errors.New("sig did not verify")
			}
		}
	}
	return nil
}

func AggregateSignatures(loopCount int) error {
	r := NewXORShift(4)
	pubkeys := make([]*bls.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := make([]*bls.Signature, 0, 1000)
	for i := 0; i < loopCount; i++ {
		priv, _ := bls.RandKey(r)
		pub := bls.PrivToPub(priv)
		msg := []byte(fmt.Sprintf(">16 character identical message %d", i))
		sig := bls.Sign(msg, priv, 0)
		pubkeys = append(pubkeys, pub)
		msgs = append(msgs, msg)
		sigs = append(sigs, sig)

		if i < 10 || i > (loopCount-5) {
			newSig := bls.AggregateSignatures(sigs)
			if !newSig.VerifyAggregate(pubkeys, msgs, 0) {
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

	pubkeys := make([]*bls.PublicKey, 0, 1000)
	msgs := make([][]byte, 0, 1000)
	sigs := bls.NewAggregateSignature()

	key, _ := bls.RandKey(r)
	pub := bls.PrivToPub(key)
	message := []byte(">16 char first message")
	sig := bls.Sign(message, key, 0)
	pubkeys = append(pubkeys, pub)
	msgs = append(msgs, message)
	sigs.Aggregate(sig)

	if !sigs.VerifyAggregate(pubkeys, msgs, 0) {
		t.Fatal("signature does not verify")
	}

	key2, _ := bls.RandKey(r)
	pub2 := bls.PrivToPub(key2)
	message2 := []byte(">16 char second message")
	sig2 := bls.Sign(message2, key2, 0)
	pubkeys = append(pubkeys, pub2)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig2)

	if !sigs.VerifyAggregate(pubkeys, msgs, 0) {
		t.Fatal("signature does not verify")
	}

	key3, _ := bls.RandKey(r)
	pub3 := bls.PrivToPub(key3)
	sig3 := bls.Sign(message2, key3, 0)
	pubkeys = append(pubkeys, pub3)
	msgs = append(msgs, message2)
	sigs.Aggregate(sig3)

	if sigs.VerifyAggregate(pubkeys, msgs, 0) {
		t.Fatal("signature verifies with duplicate message")
	}
}

func BenchmarkBLSAggregateSignature(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := bls.RandKey(r)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, priv, 0)

	s := bls.NewAggregateSignature()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Aggregate(sig)
	}
}

func BenchmarkBLSSign(b *testing.B) {
	r := NewXORShift(5)
	privs := make([]*bls.SecretKey, b.N)
	for i := range privs {
		privs[i], _ = bls.RandKey(r)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		msg := []byte(fmt.Sprintf("Hello world! 16 characters %d", i))
		bls.Sign(msg, privs[i], 0)
		// if !bls.Verify(msg, pub, sig) {
		// 	return errors.New("sig did not verify")
		// }
	}
}

func BenchmarkBLSVerify(b *testing.B) {
	r := NewXORShift(5)
	priv, _ := bls.RandKey(r)
	pub := bls.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, priv, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bls.Verify(msg, pub, sig, 0)
	}
}

func TestSignatureSerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := bls.RandKey(r)
	pub := bls.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, priv, 0)

	if !bls.Verify(msg, pub, sig, 0) {
		t.Fatal("message did not verify before serialization/deserialization")
	}

	sigSer := sig.Serialize()
	sigDeser, err := bls.DeserializeSignature(sigSer)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.Verify(msg, pub, sigDeser, 0) {
		t.Fatal("message did not verify after serialization/deserialization")
	}
}

func TestPubkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := bls.RandKey(r)
	pub := bls.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, priv, 0)

	if !bls.Verify(msg, pub, sig, 0) {
		t.Fatal("message did not verify before serialization/deserialization of pubkey")
	}

	pubSer := pub.Serialize()
	pubDeser, err := bls.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.Verify(msg, pubDeser, sig, 0) {
		t.Fatal("message did not verify after serialization/deserialization of pubkey")
	}
}

func TestSecretkeySerializeDeserialize(t *testing.T) {
	r := NewXORShift(3)
	priv, _ := bls.RandKey(r)
	privSer := priv.Serialize()
	privNew := bls.DeserializeSecretKey(privSer)
	pub := bls.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, privNew, 0)

	if !bls.Verify(msg, pub, sig, 0) {
		t.Fatal("message did not verify before serialization/deserialization of secret")
	}

	pubSer := pub.Serialize()
	pubDeser, err := bls.DeserializePublicKey(pubSer)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.Verify(msg, pubDeser, sig, 0) {
		t.Fatal("message did not verify after serialization/deserialization of secret")
	}
}

func TestPubkeySerializeDeserializeBig(t *testing.T) {
	r := NewXORShift(1)
	priv, _ := bls.RandKey(r)
	pub := bls.PrivToPub(priv)
	msg := []byte(fmt.Sprintf(">16 character identical message"))
	sig := bls.Sign(msg, priv, 0)

	if !bls.Verify(msg, pub, sig, 0) {
		t.Fatal("message did not verify before serialization/deserialization of uncompressed pubkey")
	}

	pubSer := pub.SerializeBig()
	pubDeser := bls.DeserializePublicKeyBig(pubSer)
	if !bls.Verify(msg, pubDeser, sig, 0) {
		t.Fatal("message did not verify after serialization/deserialization of uncompressed pubkey")
	}
}
