package bls_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/phoreproject/bls"
)

var expectedG1X, _ = bls.FQReprFromString("157b33991dfbf91105d15160da770421049539e3708dbe7d0fac0bcdcf70cde4e0ba392e05d248ebeea681629653501e", 16)
var expectedG1Y, _ = bls.FQReprFromString("186b975a6b00f9ed4f4415de4a7df651ed0d1ac4a2e50d145028b596b6493ec4ec6a3259d4897676c2c2c67af5f581d5", 16)

var expectedG1Hash = bls.NewG1Affine(
	bls.FQReprToFQ(expectedG1X),
	bls.FQReprToFQ(expectedG1Y),
)

func TestHashG1(t *testing.T) {
	actualHash := bls.HashG1([]byte("the message to be signed"))

	if !actualHash.Equals(expectedG1Hash) {
		t.Fatal("expected hash to match other implementations")
	}
}

func BenchmarkHashG1(t *testing.B) {
	data := make([][]byte, 100)
	r := NewXORShift(2)

	h := sha256.New()
	for i := range data {
		h.Reset()
		randBytes := make([]byte, 32)
		r.Read(randBytes)
		h.Write(randBytes)
		data[i] = h.Sum(nil)
	}

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		bls.HashG1(data[i%len(data)])
	}
}

var expectedG2c0X, _ = bls.FQReprFromString("cac64370233bfc0a5cb46981969ef1aec583bb661084c7940581cda548f5bf015c74bf23ccdb87816a3cd96ebdb2bfd", 16)
var expectedG2c1X, _ = bls.FQReprFromString("11f6e0fdfbc31b55c04cda9c896099c9f135c2eb2c504cfe0b98e98ef0e511583a9b1eb47b4c19904d820c2eab6608a9", 16)
var expectedG2c0Y, _ = bls.FQReprFromString("19ffc47d113a320834d1b3979a932c1224195be2cd83b7cd70e1800b56d9a8b94a3ce0e303cdda31e191ecc48223906f", 16)
var expectedG2c1Y, _ = bls.FQReprFromString("c75d7aae0477efd4219f9d01950fa9e83378162d7befe1e9df0e506503f04d7a49250d6a85f09dffa245a8fe583d3fd", 16)

var expectedG2Hash = bls.NewG2Affine(
	bls.NewFQ2(
		bls.FQReprToFQ(expectedG2c0X),
		bls.FQReprToFQ(expectedG2c1X),
	),
	bls.NewFQ2(
		bls.FQReprToFQ(expectedG2c0Y),
		bls.FQReprToFQ(expectedG2c1Y),
	),
)

func TestHashG2(t *testing.T) {
	actualHash := bls.HashG2([]byte("the message to be signed"))

	if !actualHash.Equals(expectedG2Hash) {
		t.Fatal("expected hash to match other implementations")
	}
}

var expectedSerializedG2, _ = hex.DecodeString("a6ef29e7241e1a1cc60fee328e3290c023d55a6701db500eefab7f91391a8b8726fd0024121e64637281f907137fe268187b4baca36388e96194b73a7d532f6eea6bc098778dbfd3404584613b5ba9da97d5602e31fdbe9270b863876529b254")

func TestHashG2WithDomain(t *testing.T) {
	actualHash := bls.HashG2WithDomain([32]byte{}, [8]byte{})

	compressedPoint := bls.CompressG2(actualHash.ToAffine())

	if !bytes.Equal(expectedSerializedG2, compressedPoint[:]) {
		t.Fatal("expected hash to match test")
	}
}

func BenchmarkHashG2(t *testing.B) {
	data := make([][]byte, 100)
	r := NewXORShift(2)

	h := sha256.New()
	for i := range data {
		h.Reset()
		randBytes := make([]byte, 32)
		r.Read(randBytes)
		h.Write(randBytes)
		data[i] = h.Sum(nil)
	}

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		bls.HashG2(data[i%len(data)])
	}
}
