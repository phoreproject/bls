package bls_test

import (
	"encoding/hex"
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

func TestAggregateSigsSeparate(t *testing.T) {
	x := NewXORShift(20)
	priv1, _ := bls.RandKey(x)
	priv2, _ := bls.RandKey(x)
	priv3, _ := bls.RandKey(x)

	pub1 := bls.PrivToPub(priv1)
	pub2 := bls.PrivToPub(priv2)
	pub3 := bls.PrivToPub(priv3)

	msg := []byte("test 1")
	sig1 := bls.Sign(msg, priv1, 0)
	sig2 := bls.Sign(msg, priv2, 0)
	sig3 := bls.Sign(msg, priv3, 0)

	aggSigs := bls.AggregateSignatures([]*bls.Signature{sig1, sig2, sig3})

	aggPubs := bls.NewAggregatePubkey()
	aggPubs.Aggregate(pub1)
	aggPubs.Aggregate(pub2)
	aggPubs.Aggregate(pub3)

	valid := bls.Verify(msg, aggPubs, aggSigs, 0)
	if !valid {
		t.Fatal("expected aggregate signature to be valid")
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

func TestSignatureAggregationSorting(t *testing.T) {
	pub0xc0, _ := bls.FQReprFromString("004f797e1f17917e8e16047c9f092c1739607f3650726f81805797d1998354400fa08378bf2bd5b840cb34a8a6af949c", 16)
	pub0xc1, _ := bls.FQReprFromString("16d279303ed244852fa687db9187cbfc7b2f4a96464393b298f2fab94e973422d02b38e0899b5bf5660ca2b3b4a6c8f7", 16)
	pub0yc0, _ := bls.FQReprFromString("097ca937ba99f1a4201b13b527029a8fe3527ddb02f9e607a1fb2b509c6ce1bac8b8ced37530c3a10962335be82568de", 16)
	pub0yc1, _ := bls.FQReprFromString("1672f96d77e3c9ef82c0d7768fb20c9e38fe36f46939f6679d60c9b5099f2bdb796499d0a016ca0b79ec1a56907b31a8", 16)
	pub0 := bls.NewPublicKeyFromG1(bls.NewG2Affine(
		bls.NewFQ2(
			bls.FQReprToFQ(pub0xc0),
			bls.FQReprToFQ(pub0xc1)),
		bls.NewFQ2(
			bls.FQReprToFQ(pub0yc0),
			bls.FQReprToFQ(pub0yc1))))

	pub1xc0, _ := bls.FQReprFromString("0881d52a11b70e812cc948a50757903739a48a2242588daa1f5d505037940c365ef34864c9aea388d2e04fd32e7603de", 16)
	pub1xc1, _ := bls.FQReprFromString("0cb0c79c3c1fd4958e41f8f4a4f977c924f9942f87ff6321385fbbccc2aa60f57b2ea0e208ce5fe9bb0f38186b953260", 16)
	pub1yc0, _ := bls.FQReprFromString("11a1db660a2d8d5d87c248c01b9bc12ee06f8bcb505cb810c8fc8ee031a39c1b795793daa8300a56a02bdd3f67cf9059", 16)
	pub1yc1, _ := bls.FQReprFromString("11fa167610bec1cd1e8e2ecddb7eb1bf10c31b55094cca40c03dbcdd1da4f9f36c48195769a36e538283a86b0372a3c9", 16)
	pub1 := bls.NewPublicKeyFromG1(bls.NewG2Affine(
		bls.NewFQ2(
			bls.FQReprToFQ(pub1xc0),
			bls.FQReprToFQ(pub1xc1)),
		bls.NewFQ2(
			bls.FQReprToFQ(pub1yc0),
			bls.FQReprToFQ(pub1yc1))))

	var msg0 [32]byte
	var msg1 [32]byte

	msg0Bytes, _ := hex.DecodeString("a6b157433a3f9477e08ebf1de817b5443d60486044b8b9295418c3b656146d67")
	msg1Bytes, _ := hex.DecodeString("59d1c0f3a6bd317f142ac7e9ae68233aa76b595db6865ec1637cd8c303eab77a")

	copy(msg0[:], msg0Bytes)
	copy(msg1[:], msg1Bytes)

	sig0x, _ := bls.FQReprFromString("0b148de9dc19c3710d0849dd566d593b4f8d5383cf1cab7be65350c0e988ee79466174be2e707d378a3b5bdb9569833d", 16)
	sig0y, _ := bls.FQReprFromString("17bd938d081276d700e27455951b12abe00456e9a359692c9740523a5adc71f2399c5a76530e720af2c4b60e02e31067", 16)
	sig1x, _ := bls.FQReprFromString("0c2d61978d38e6e0043a5efac8c337e02cb9c81692f46a1fd60a7ac3a728d93a5f1f692b9c6fd349a7df5c70c09db9cf", 16)
	sig1y, _ := bls.FQReprFromString("11439b6da88eb86f225313e6e0d51aa3d344d99e8cf999487c092c16f563b8bdfc71ba669b8d1a805b04da55afb30750", 16)

	sig0 := bls.NewSignatureFromG1(bls.NewG1Affine(
		bls.FQReprToFQ(sig0x),
		bls.FQReprToFQ(sig0y)))

	sig1 := bls.NewSignatureFromG1(bls.NewG1Affine(
		bls.FQReprToFQ(sig1x),
		bls.FQReprToFQ(sig1y)))

	if !bls.Verify(msg0[:], pub0, sig0, 0) {
		t.Fatal("sig0 did not verify")
	}

	if !bls.Verify(msg1[:], pub1, sig1, 0) {
		t.Fatal("sig1 did not verify")
	}

	aggSig := bls.NewAggregateSignature()
	aggSig.Aggregate(sig0)
	aggSig.Aggregate(sig1)

	if !aggSig.VerifyAggregate([]*bls.PublicKey{pub0, pub1}, [][]byte{msg0[:], msg1[:]}, 0) {
		t.Fatal("aggregate signature did not verify")
	}
}
