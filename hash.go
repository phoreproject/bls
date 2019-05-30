package bls

import (
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

func HashSecretKey(b [32]byte) *FR {
	// this implements hash_to_field as defined in the IETF standard with:
	// msg = b
	// ctr = 0
	// m = 1
	// hash_fn = SHA256
	// hash_reps = 2

	h := sha256.New()
	h.Write(b[:])
	msgHash := h.Sum(nil)
	msgPrime := append(msgHash, 0x00)

	t := []byte{}
	for j := uint16(1); j <= uint16(2); j++ {
		h := sha256.New()
		h.Write(msgPrime)
		h.Write([]byte{0x01})
		var jByte [2]byte
		binary.BigEndian.PutUint16(jByte[:], j)
		h.Write(jByte[1:])
		out := h.Sum(nil)
		t = append(t, out...)
	}

	tBig := new(big.Int)
	tBig.SetBytes(t)
	tBig.Mod(tBig, RFieldModulus.ToBig())
	tFR, _ := FRReprFromBigInt(tBig)
	return FRReprToFR(tFR)
}

func hp(msg []byte, ctr uint8) FQ {
	// this implements hash_to_field as defined in the IETF standard with:
	// msg = msg
	// ctr = ctr
	// m = 1
	// hash_fn = SHA256
	// hash_reps = 2
	h := sha256.New()
	h.Write(msg)
	msgHash := h.Sum(nil)
	var ctrBytes [2]byte
	binary.BigEndian.PutUint16(ctrBytes[:], uint16(ctr))
	msgPrime := append(msgHash, ctrBytes[1:]...)

	t := []byte{}
	for j := uint16(1); j <= uint16(2); j++ {
		h := sha256.New()
		h.Write(msgPrime)
		h.Write([]byte{0x01})
		var jByte [2]byte
		binary.BigEndian.PutUint16(jByte[:], j)
		h.Write(jByte[1:])
		out := h.Sum(nil)
		t = append(t, out...)
	}

	tBig := new(big.Int)
	tBig.SetBytes(t)
	tBig.Mod(tBig, QFieldModulus.ToBig())
	tFQ, _ := FQReprFromBigInt(tBig)
	return FQReprToFQ(tFQ)
}

func hp2(msg []byte, ctr uint8) FQ2 {
	// this implements hash_to_field as defined in the IETF standard with:
	// msg = msg
	// ctr = ctr
	// m = 2
	// hash_fn = SHA256
	// hash_reps = 2
	h := sha256.New()
	h.Write(msg)
	msgHash := h.Sum(nil)
	var ctrBytes [2]byte
	binary.BigEndian.PutUint16(ctrBytes[:], uint16(ctr))
	msgPrime := append(msgHash, ctrBytes[1:]...)

	var fqs [2]FQ
	for i := uint16(1); i <= 2; i++ {
		t := []byte{}

		for j := uint16(1); j <= uint16(2); j++ {
			h := sha256.New()
			h.Write(msgPrime)
			var iByte [2]byte
			binary.BigEndian.PutUint16(iByte[:], i)
			h.Write(iByte[1:])
			var jByte [2]byte
			binary.BigEndian.PutUint16(jByte[:], j)
			h.Write(jByte[1:])
			out := h.Sum(nil)
			t = append(t, out...)
		}

		tBig := new(big.Int)
		tBig.SetBytes(t)
		tBig.Mod(tBig, QFieldModulus.ToBig())
		tFQ, _ := FQReprFromBigInt(tBig)
		fqs[i-1] = FQReprToFQ(tFQ)
	}

	return NewFQ2(fqs[0], fqs[1])
}

// 11-isogeny from Ell1' to Ell1 map coefficients
var xNum11 = []FQ{
	FQReprToFQ(fqReprFromHexUnchecked("11a05f2b1e833340b809101dd99815856b303e88a2d7005ff2627b56cdb4e2c85610c2d5f2e62d6eaeac1662734649b7")),
	FQReprToFQ(fqReprFromHexUnchecked("17294ed3e943ab2f0588bab22147a81c7c17e75b2f6a8417f565e33c70d1e86b4838f2a6f318c356e834eef1b3cb83bb")),
	FQReprToFQ(fqReprFromHexUnchecked("d54005db97678ec1d1048c5d10a9a1bce032473295983e56878e501ec68e25c958c3e3d2a09729fe0179f9dac9edcb0")),
	FQReprToFQ(fqReprFromHexUnchecked("1778e7166fcc6db74e0609d307e55412d7f5e4656a8dbf25f1b33289f1b330835336e25ce3107193c5b388641d9b6861")),
	FQReprToFQ(fqReprFromHexUnchecked("e99726a3199f4436642b4b3e4118e5499db995a1257fb3f086eeb65982fac18985a286f301e77c451154ce9ac8895d9")),
	FQReprToFQ(fqReprFromHexUnchecked("1630c3250d7313ff01d1201bf7a74ab5db3cb17dd952799b9ed3ab9097e68f90a0870d2dcae73d19cd13c1c66f652983")),
	FQReprToFQ(fqReprFromHexUnchecked("d6ed6553fe44d296a3726c38ae652bfb11586264f0f8ce19008e218f9c86b2a8da25128c1052ecaddd7f225a139ed84")),
	FQReprToFQ(fqReprFromHexUnchecked("17b81e7701abdbe2e8743884d1117e53356de5ab275b4db1a682c62ef0f2753339b7c8f8c8f475af9ccb5618e3f0c88e")),
	FQReprToFQ(fqReprFromHexUnchecked("80d3cf1f9a78fc47b90b33563be990dc43b756ce79f5574a2c596c928c5d1de4fa295f296b74e956d71986a8497e317")),
	FQReprToFQ(fqReprFromHexUnchecked("169b1f8e1bcfa7c42e0c37515d138f22dd2ecb803a0c5c99676314baf4bb1b7fa3190b2edc0327797f241067be390c9e")),
	FQReprToFQ(fqReprFromHexUnchecked("10321da079ce07e272d8ec09d2565b0dfa7dccdde6787f96d50af36003b14866f69b771f8c285decca67df3f1605fb7b")),
	FQReprToFQ(fqReprFromHexUnchecked("6e08c248e260e70bd1e962381edee3d31d79d7e22c837bc23c0bf1bc24c6b68c24b1b80b64d391fa9c8ba2e8ba2d229")),
}

var xDen11 = []FQ{
	FQReprToFQ(fqReprFromHexUnchecked("8ca8d548cff19ae18b2e62f4bd3fa6f01d5ef4ba35b48ba9c9588617fc8ac62b558d681be343df8993cf9fa40d21b1c")),
	FQReprToFQ(fqReprFromHexUnchecked("12561a5deb559c4348b4711298e536367041e8ca0cf0800c0126c2588c48bf5713daa8846cb026e9e5c8276ec82b3bff")),
	FQReprToFQ(fqReprFromHexUnchecked("b2962fe57a3225e8137e629bff2991f6f89416f5a718cd1fca64e00b11aceacd6a3d0967c94fedcfcc239ba5cb83e19")),
	FQReprToFQ(fqReprFromHexUnchecked("3425581a58ae2fec83aafef7c40eb545b08243f16b1655154cca8abc28d6fd04976d5243eecf5c4130de8938dc62cd8")),
	FQReprToFQ(fqReprFromHexUnchecked("13a8e162022914a80a6f1d5f43e7a07dffdfc759a12062bb8d6b44e833b306da9bd29ba81f35781d539d395b3532a21e")),
	FQReprToFQ(fqReprFromHexUnchecked("e7355f8e4e667b955390f7f0506c6e9395735e9ce9cad4d0a43bcef24b8982f7400d24bc4228f11c02df9a29f6304a5")),
	FQReprToFQ(fqReprFromHexUnchecked("772caacf16936190f3e0c63e0596721570f5799af53a1894e2e073062aede9cea73b3538f0de06cec2574496ee84a3a")),
	FQReprToFQ(fqReprFromHexUnchecked("14a7ac2a9d64a8b230b3f5b074cf01996e7f63c21bca68a81996e1cdf9822c580fa5b9489d11e2d311f7d99bbdcc5a5e")),
	FQReprToFQ(fqReprFromHexUnchecked("a10ecf6ada54f825e920b3dafc7a3cce07f8d1d7161366b74100da67f39883503826692abba43704776ec3a79a1d641")),
	FQReprToFQ(fqReprFromHexUnchecked("95fc13ab9e92ad4476d6e3eb3a56680f682b4ee96f7d03776df533978f31c1593174e4b4b7865002d6384d168ecdd0a")),
	FQReprToFQ(fqReprFromHexUnchecked("1")),
}

var yNum11 = []FQ{
	FQReprToFQ(fqReprFromHexUnchecked("90d97c81ba24ee0259d1f094980dcfa11ad138e48a869522b52af6c956543d3cd0c7aee9b3ba3c2be9845719707bb33")),
	FQReprToFQ(fqReprFromHexUnchecked("134996a104ee5811d51036d776fb46831223e96c254f383d0f906343eb67ad34d6c56711962fa8bfe097e75a2e41c696")),
	FQReprToFQ(fqReprFromHexUnchecked("cc786baa966e66f4a384c86a3b49942552e2d658a31ce2c344be4b91400da7d26d521628b00523b8dfe240c72de1f6")),
	FQReprToFQ(fqReprFromHexUnchecked("1f86376e8981c217898751ad8746757d42aa7b90eeb791c09e4a3ec03251cf9de405aba9ec61deca6355c77b0e5f4cb")),
	FQReprToFQ(fqReprFromHexUnchecked("8cc03fdefe0ff135caf4fe2a21529c4195536fbe3ce50b879833fd221351adc2ee7f8dc099040a841b6daecf2e8fedb")),
	FQReprToFQ(fqReprFromHexUnchecked("16603fca40634b6a2211e11db8f0a6a074a7d0d4afadb7bd76505c3d3ad5544e203f6326c95a807299b23ab13633a5f0")),
	FQReprToFQ(fqReprFromHexUnchecked("4ab0b9bcfac1bbcb2c977d027796b3ce75bb8ca2be184cb5231413c4d634f3747a87ac2460f415ec961f8855fe9d6f2")),
	FQReprToFQ(fqReprFromHexUnchecked("987c8d5333ab86fde9926bd2ca6c674170a05bfe3bdd81ffd038da6c26c842642f64550fedfe935a15e4ca31870fb29")),
	FQReprToFQ(fqReprFromHexUnchecked("9fc4018bd96684be88c9e221e4da1bb8f3abd16679dc26c1e8b6e6a1f20cabe69d65201c78607a360370e577bdba587")),
	FQReprToFQ(fqReprFromHexUnchecked("e1bba7a1186bdb5223abde7ada14a23c42a0ca7915af6fe06985e7ed1e4d43b9b3f7055dd4eba6f2bafaaebca731c30")),
	FQReprToFQ(fqReprFromHexUnchecked("19713e47937cd1be0dfd0b8f1d43fb93cd2fcbcb6caf493fd1183e416389e61031bf3a5cce3fbafce813711ad011c132")),
	FQReprToFQ(fqReprFromHexUnchecked("18b46a908f36f6deb918c143fed2edcc523559b8aaf0c2462e6bfe7f911f643249d9cdf41b44d606ce07c8a4d0074d8e")),
	FQReprToFQ(fqReprFromHexUnchecked("b182cac101b9399d155096004f53f447aa7b12a3426b08ec02710e807b4633f06c851c1919211f20d4c04f00b971ef8")),
	FQReprToFQ(fqReprFromHexUnchecked("245a394ad1eca9b72fc00ae7be315dc757b3b080d4c158013e6632d3c40659cc6cf90ad1c232a6442d9d3f5db980133")),
	FQReprToFQ(fqReprFromHexUnchecked("5c129645e44cf1102a159f748c4a3fc5e673d81d7e86568d9ab0f5d396a7ce46ba1049b6579afb7866b1e715475224b")),
	FQReprToFQ(fqReprFromHexUnchecked("15e6be4e990f03ce4ea50b3b42df2eb5cb181d8f84965a3957add4fa95af01b2b665027efec01c7704b456be69c8b604")),
}

var yDen11 = []FQ{
	FQReprToFQ(fqReprFromHexUnchecked("16112c4c3a9c98b252181140fad0eae9601a6de578980be6eec3232b5be72e7a07f3688ef60c206d01479253b03663c1")),
	FQReprToFQ(fqReprFromHexUnchecked("1962d75c2381201e1a0cbd6c43c348b885c84ff731c4d59ca4a10356f453e01f78a4260763529e3532f6102c2e49a03d")),
	FQReprToFQ(fqReprFromHexUnchecked("58df3306640da276faaae7d6e8eb15778c4855551ae7f310c35a5dd279cd2eca6757cd636f96f891e2538b53dbf67f2")),
	FQReprToFQ(fqReprFromHexUnchecked("16b7d288798e5395f20d23bf89edb4d1d115c5dbddbcd30e123da489e726af41727364f2c28297ada8d26d98445f5416")),
	FQReprToFQ(fqReprFromHexUnchecked("be0e079545f43e4b00cc912f8228ddcc6d19c9f0f69bbb0542eda0fc9dec916a20b15dc0fd2ededda39142311a5001d")),
	FQReprToFQ(fqReprFromHexUnchecked("8d9e5297186db2d9fb266eaac783182b70152c65550d881c5ecd87b6f0f5a6449f38db9dfa9cce202c6477faaf9b7ac")),
	FQReprToFQ(fqReprFromHexUnchecked("166007c08a99db2fc3ba8734ace9824b5eecfdfa8d0cf8ef5dd365bc400a0051d5fa9c01a58b1fb93d1a1399126a775c")),
	FQReprToFQ(fqReprFromHexUnchecked("16a3ef08be3ea7ea03bcddfabba6ff6ee5a4375efa1f4fd7feb34fd206357132b920f5b00801dee460ee415a15812ed9")),
	FQReprToFQ(fqReprFromHexUnchecked("1866c8ed336c61231a1be54fd1d74cc4f9fb0ce4c6af5920abc5750c4bf39b4852cfe2f7bb9248836b233d9d55535d4a")),
	FQReprToFQ(fqReprFromHexUnchecked("167a55cda70a6e1cea820597d94a84903216f763e13d87bb5308592e7ea7d4fbc7385ea3d529b35e346ef48bb8913f55")),
	FQReprToFQ(fqReprFromHexUnchecked("4d2f259eea405bd48f010a01ad2911d9c6dd039bb61a6290e591b36e636a5c871a5c29f4f83060400f8b49cba8f6aa8")),
	FQReprToFQ(fqReprFromHexUnchecked("accbb67481d033ff5852c1e48c50c477f94ff8aefce42d28c0f9a88cea7913516f968986f7ebbea9684b529e2561092")),
	FQReprToFQ(fqReprFromHexUnchecked("ad6b9514c767fe3c3613144b45f1496543346d98adf02267d5ceef9a00d9b8693000763e3b90ac11e99b138573345cc")),
	FQReprToFQ(fqReprFromHexUnchecked("2660400eb2e4f3b628bdd0d53cd76f2bf565b94e72927c1cb748df27942480e420517bd8714cc80d1fadc1326ed06f7")),
	FQReprToFQ(fqReprFromHexUnchecked("e0fa1d816ddc03e6b24255e0d7819c171c40f65e273b853324efcd6356caa205ca2f570f13497804415473a1d634b8f")),
	FQReprToFQ(fqReprFromHexUnchecked("1")),
}

var mapCoeffs11 = [][]FQ{xNum11, xDen11, yNum11, yDen11}

func iso11(p *G1Affine) *G1Affine {
	x := p.x
	y := p.y
	mapVals := make([]FQ, 4)

	for idx, coeffs := range mapCoeffs11 {
		mapVals[idx] = coeffs[len(coeffs)-1]
		for coeffIdx := len(coeffs) - 2; coeffIdx >= 0; coeffIdx-- {
			mapVals[idx].MulAssign(x)
			mapVals[idx].AddAssign(coeffs[coeffIdx])
		}
	}

	newX := mapVals[0]
	newX.DivAssign(mapVals[1])

	newY := y.Copy()
	newY.MulAssign(mapVals[2])
	newY.DivAssign(mapVals[3])

	return NewG1Affine(newX, newY)
}

var xNum3 = []FQ2{
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("5c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97d6")),
		FQReprToFQ(fqReprFromHexUnchecked("5c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97d6")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("0")),
		FQReprToFQ(fqReprFromHexUnchecked("11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71a")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71e")),
		FQReprToFQ(fqReprFromHexUnchecked("8ab05f8bdd54cde190937e76bc3e447cc27c3d6fbd7063fcd104635a790520c0a395554e5c6aaaa9354ffffffffe38d")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("171d6541fa38ccfaed6dea691f5fb614cb14b4e7f4e810aa22d6108f142b85757098e38d0f671c7188e2aaaaaaaa5ed1")),
		FQReprToFQ(fqReprFromHexUnchecked("0")),
	),
}

var xDen3 = []FQ2{
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("0")),
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa63")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("c")),
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa9f")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("1")),
		FQReprToFQ(fqReprFromHexUnchecked("0")),
	),
}

var yNum3 = []FQ2{
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("1530477c7ab4113b59a4c18b076d11930f7da5d4a07f649bf54439d87d27e500fc8c25ebf8c92f6812cfc71c71c6d706")),
		FQReprToFQ(fqReprFromHexUnchecked("1530477c7ab4113b59a4c18b076d11930f7da5d4a07f649bf54439d87d27e500fc8c25ebf8c92f6812cfc71c71c6d706")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("0")),
		FQReprToFQ(fqReprFromHexUnchecked("5c759507e8e333ebb5b7a9a47d7ed8532c52d39fd3a042a88b58423c50ae15d5c2638e343d9c71c6238aaaaaaaa97be")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("11560bf17baa99bc32126fced787c88f984f87adf7ae0c7f9a208c6b4f20a4181472aaa9cb8d555526a9ffffffffc71c")),
		FQReprToFQ(fqReprFromHexUnchecked("8ab05f8bdd54cde190937e76bc3e447cc27c3d6fbd7063fcd104635a790520c0a395554e5c6aaaa9354ffffffffe38f")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("124c9ad43b6cf79bfbf7043de3811ad0761b0f37a1e26286b0e977c69aa274524e79097a56dc4bd9e1b371c71c718b10")),
		FQReprToFQ(fqReprFromHexUnchecked("0")),
	),
}

var yDen3 = []FQ2{
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa8fb")),
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa8fb")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("0")),
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffa9d3")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("12")),
		FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaa99")),
	),
	NewFQ2(
		FQReprToFQ(fqReprFromHexUnchecked("1")),
		FQReprToFQ(fqReprFromHexUnchecked("0")),
	),
}

var mapCoeffs3 = [][]FQ2{xNum3, xDen3, yNum3, yDen3}

func iso3(p *G2Affine) *G2Affine {
	x := p.x
	y := p.y
	mapVals := make([]FQ2, 4)

	for idx, coeffs := range mapCoeffs3 {
		mapVals[idx] = coeffs[len(coeffs)-1]
		for coeffIdx := len(coeffs) - 2; coeffIdx >= 0; coeffIdx-- {
			mapVals[idx].MulAssign(x)
			mapVals[idx].AddAssign(coeffs[coeffIdx])
		}
	}

	newX := mapVals[0]
	newX.DivAssign(mapVals[1])

	newY := y.Copy()
	newY.MulAssign(mapVals[2])
	newY.DivAssign(mapVals[3])

	return NewG2Affine(newX, newY)
}

// ClearH clears the cofactor for Ell1.
func ClearH(p *G1Affine) *G1Affine {
	xP := p.Mul(NewFQRepr(0xd201000000010000))
	return xP.AddAffine(p).ToAffine()
}

func optimizedSWUMap(t1 *FQ, t2 *FQ) *G1Affine {
	Pp := optimizedSWUMapHelper(*t1)

	if t2 != nil {
		Pp2 := optimizedSWUMapHelper(*t2)

		Pp = Pp.ToProjective().AddAffine(Pp2).ToAffine()
	}
	Pp = iso11(Pp)
	return ClearH(Pp)
}

const cipherSuite = 0x01

// HashG1 converts a message to a point on the G2 curve.
func HashG1(msg []byte) *G1Affine {
	cipherSuiteAndMessage := append([]byte{cipherSuite}, msg...)
	t1 := hp(cipherSuiteAndMessage, 0)
	t2 := hp(cipherSuiteAndMessage, 1)
	return optimizedSWUMap(&t1, &t2)
}

var iwsc = NewFQ2(
	FQReprToFQ(fqReprFromHexUnchecked("d0088f51cbff34d258dd3db21a5d66bb23ba5c279c2895fb39869507b587b120f55ffff58a9ffffdcff7fffffffd556")),
	FQReprToFQ(fqReprFromHexUnchecked("d0088f51cbff34d258dd3db21a5d66bb23ba5c279c2895fb39869507b587b120f55ffff58a9ffffdcff7fffffffd555")),
)

var kQiX = FQReprToFQ(fqReprFromHexUnchecked("1a0111ea397fe699ec02408663d4de85aa0d857d89759ad4897d29650fb85f9b409427eb4f49fffd8bfd00000000aaad"))
var kQiY = FQReprToFQ(fqReprFromHexUnchecked("6af0e0437ff400b6831e36d6bd17ffe48395dabc2d3435e77f76e17009241c5ee67992f72ec05f4c81084fbede3cc09"))

func psi(g *G2Affine) *G2Affine {
	newX := fq2nqr.Copy()
	qiX := iwsc.Copy()
	qiX.MulAssign(g.x)
	qiX.c0.MulAssign(kQiX)
	qiX.c1.MulAssign(kQiX)
	qiX.c1.NegAssign()
	newX.MulAssign(qiX)

	newY := fq2nqr.Copy()
	qiY := iwsc.Copy()
	qiY.MulAssign(g.y)

	y0y1 := qiY.c0.Copy()
	y0y1.AddAssign(qiY.c1)
	y0y1.MulAssign(kQiY)

	y0SubY1 := qiY.c0.Copy()
	y0SubY1.SubAssign(qiY.c1)
	y0SubY1.MulAssign(kQiY)

	qiY = NewFQ2(y0y1, y0SubY1)
	newY.MulAssign(qiY)

	return NewG2Affine(newX, newY)
}

func clearH2(p *G2Affine) *G2Affine {
	work := p.Mul(NewFQRepr(0xd201000000010000))

	work = work.AddAffine(p)

	minusPsiP := psi(p)
	minusPsiP.NegAssign()

	work = work.AddAffine(minusPsiP)

	work = work.Mul(NewFQRepr(0xd201000000010000))
	work = work.AddAffine(minusPsiP)

	negP := p.Copy()
	negP.NegAssign()

	work = work.AddAffine(negP)
	p2 := p.ToProjective().Double().ToAffine()
	psiPsi2P := psi(psi(p2))
	work = work.AddAffine(psiPsi2P)
	return work.ToAffine()
}

func optimizedSWUMap2(t1 *FQ2, t2 *FQ2) *G2Affine {
	Pp := OptimizedSWU2MapHelper(*t1)

	if t2 != nil {
		Pp2 := OptimizedSWU2MapHelper(*t2)

		Pp = Pp.ToProjective().AddAffine(Pp2).ToAffine()
	}
	Pp = iso3(Pp)

	return clearH2(Pp)
}

// HashG2 converts a message to a point on the G2 curve.
func HashG2(msg []byte) *G2Affine {
	cipherSuiteAndMessage := append([]byte{cipherSuite}, msg...)
	t1 := hp2(cipherSuiteAndMessage, 0)
	t2 := hp2(cipherSuiteAndMessage, 1)
	h := optimizedSWUMap2(&t1, &t2)
	return h
}
