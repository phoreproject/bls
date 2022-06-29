package bls

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

type xorShift struct {
	state uint64
}

func newXorShift(state uint64) *xorShift {
	return &xorShift{state}
}

func (xor *xorShift) Read(b []byte) (int, error) {
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

var expectedG1X, _ = FQReprFromString("157b33991dfbf91105d15160da770421049539e3708dbe7d0fac0bcdcf70cde4e0ba392e05d248ebeea681629653501e", 16)
var expectedG1Y, _ = FQReprFromString("186b975a6b00f9ed4f4415de4a7df651ed0d1ac4a2e50d145028b596b6493ec4ec6a3259d4897676c2c2c67af5f581d5", 16)

var expectedG1Hash = NewG1Affine(
	FQReprToFQ(expectedG1X),
	FQReprToFQ(expectedG1Y),
)

func TestHashG1(t *testing.T) {
	actualHash := HashG1([]byte("the message to be signed"))

	if !actualHash.Equals(expectedG1Hash) {
		t.Fatal("expected hash to match other implementations")
	}
}

func BenchmarkHashG1(t *testing.B) {
	data := make([][]byte, 100)
	r := newXorShift(2)

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
		HashG1(data[i%len(data)])
	}
}

var expectedG2c0X, _ = FQReprFromString("cac64370233bfc0a5cb46981969ef1aec583bb661084c7940581cda548f5bf015c74bf23ccdb87816a3cd96ebdb2bfd", 16)
var expectedG2c1X, _ = FQReprFromString("11f6e0fdfbc31b55c04cda9c896099c9f135c2eb2c504cfe0b98e98ef0e511583a9b1eb47b4c19904d820c2eab6608a9", 16)
var expectedG2c0Y, _ = FQReprFromString("19ffc47d113a320834d1b3979a932c1224195be2cd83b7cd70e1800b56d9a8b94a3ce0e303cdda31e191ecc48223906f", 16)
var expectedG2c1Y, _ = FQReprFromString("c75d7aae0477efd4219f9d01950fa9e83378162d7befe1e9df0e506503f04d7a49250d6a85f09dffa245a8fe583d3fd", 16)

var expectedG2Hash = NewG2Affine(
	NewFQ2(
		FQReprToFQ(expectedG2c0X),
		FQReprToFQ(expectedG2c1X),
	),
	NewFQ2(
		FQReprToFQ(expectedG2c0Y),
		FQReprToFQ(expectedG2c1Y),
	),
)

func TestHashG2(t *testing.T) {
	actualHash := HashG2([]byte("the message to be signed"))

	if !actualHash.Equals(expectedG2Hash) {
		t.Fatal("expected hash to match other implementations")
	}
}

var expectedSerializedG2, _ = hex.DecodeString("a6ef29e7241e1a1cc60fee328e3290c023d55a6701db500eefab7f91391a8b8726fd0024121e64637281f907137fe268187b4baca36388e96194b73a7d532f6eea6bc098778dbfd3404584613b5ba9da97d5602e31fdbe9270b863876529b254")

func TestHashG2WithDomain(t *testing.T) {
	actualHash := HashG2WithDomain([32]byte{}, [8]byte{})

	compressedPoint := CompressG2(actualHash.ToAffine())

	if !bytes.Equal(expectedSerializedG2, compressedPoint[:]) {
		t.Fatal("expected hash to match test")
	}
}

func BenchmarkHashG2(t *testing.B) {
	data := make([][]byte, 100)
	r := newXorShift(2)

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
		HashG2(data[i%len(data)])
	}
}

func decodeHexOrDie(hexStr string) []byte {
	out, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return out
}

func TestExpandMessage(t *testing.T) {
	// tests from IETF spec
	// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#appendix-K.1
	tests := []struct {
		msg        []byte
		dst        []byte
		lenInBytes uint16
		outBytes   []byte
	}{
		{
			msg:        []byte{},
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x20,
			outBytes:   decodeHexOrDie("68a985b87eb6b46952128911f2a4412bbc302a9d759667f87f7a21d803f07235"),
		},
		{
			msg:        []byte("abc"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x20,
			outBytes:   decodeHexOrDie("d8ccab23b5985ccea865c6c97b6e5b8350e794e603b4b97902f53a8a0d605615"),
		},
		{
			msg:        []byte("abcdef0123456789"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x20,
			outBytes:   decodeHexOrDie("eff31487c770a893cfb36f912fbfcbff40d5661771ca4b2cb4eafe524333f5c1"),
		},
		{
			msg:        []byte("q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x20,
			outBytes:   decodeHexOrDie("b23a1d2b4d97b2ef7785562a7e8bac7eed54ed6e97e29aa51bfe3f12ddad1ff9"),
		},
		{
			msg:        []byte("a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x20,
			outBytes:   decodeHexOrDie("4623227bcc01293b8c130bf771da8c298dede7383243dc0993d2d94823958c4c"),
		},
		{
			msg:        []byte(""),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x80,
			outBytes:   decodeHexOrDie("af84c27ccfd45d41914fdff5df25293e221afc53d8ad2ac06d5e3e29485dadbee0d121587713a3e0dd4d5e69e93eb7cd4f5df4cd103e188cf60cb02edc3edf18eda8576c412b18ffb658e3dd6ec849469b979d444cf7b26911a08e63cf31f9dcc541708d3491184472c2c29bb749d4286b004ceb5ee6b9a7fa5b646c993f0ced"),
		},
		{
			msg:        []byte("abc"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x80,
			outBytes:   decodeHexOrDie("abba86a6129e366fc877aab32fc4ffc70120d8996c88aee2fe4b32d6c7b6437a647e6c3163d40b76a73cf6a5674ef1d890f95b664ee0afa5359a5c4e07985635bbecbac65d747d3d2da7ec2b8221b17b0ca9dc8a1ac1c07ea6a1e60583e2cb00058e77b7b72a298425cd1b941ad4ec65e8afc50303a22c0f99b0509b4c895f40"),
		},
		{
			msg:        []byte("abcdef0123456789"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x80,
			outBytes:   decodeHexOrDie("ef904a29bffc4cf9ee82832451c946ac3c8f8058ae97d8d629831a74c6572bd9ebd0df635cd1f208e2038e760c4994984ce73f0d55ea9f22af83ba4734569d4bc95e18350f740c07eef653cbb9f87910d833751825f0ebefa1abe5420bb52be14cf489b37fe1a72f7de2d10be453b2c9d9eb20c7e3f6edc5a60629178d9478df"),
		},
		{
			msg:        []byte("q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x80,
			outBytes:   decodeHexOrDie("80be107d0884f0d881bb460322f0443d38bd222db8bd0b0a5312a6fedb49c1bbd88fd75d8b9a09486c60123dfa1d73c1cc3169761b17476d3c6b7cbbd727acd0e2c942f4dd96ae3da5de368d26b32286e32de7e5a8cb2949f866a0b80c58116b29fa7fabb3ea7d520ee603e0c25bcaf0b9a5e92ec6a1fe4e0391d1cdbce8c68a"),
		},
		{
			msg:        []byte("a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			dst:        []byte("QUUX-V01-CS02-with-expander-SHA256-128"),
			lenInBytes: 0x80,
			outBytes:   decodeHexOrDie("546aff5444b5b79aa6148bd81728704c32decb73a3ba76e9e75885cad9def1d06d6792f8a7d12794e90efed817d96920d728896a4510864370c207f99bd4a608ea121700ef01ed879745ee3e4ceef777eda6d9e5e38b90c86ea6fb0b36504ba4a45d22e86f6db5dd43d98a294bebb9125d5b794e9d2a81181066eb954966a487"),
		},
	}

	for _, test := range tests {
		hash := expandMessageXmd(test.msg, test.dst, test.lenInBytes)

		if !bytes.Equal(hash, test.outBytes) {
			t.Fatalf("expected ExpandMessageXmd(%s, %s, %d) output of %s, but got %s", hex.EncodeToString(test.msg), hex.EncodeToString(test.dst), test.lenInBytes, hex.EncodeToString(test.outBytes), hex.EncodeToString(hash))
		}
	}
}
