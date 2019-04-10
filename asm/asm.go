// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/gotypes"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// muladdc multiplies a and b,
func muladdc(loRef, a, b, carry Register) {
	MOVQ(a, RAX)
	MULQ(b)

	ADDQ(loRef, RAX)
	ADCQ(Imm(0), RDX)
	ADDQ(carry, RAX)
	ADCQ(Imm(0), RDX)

	MOVQ(RDX, carry)
	MOVQ(RAX, loRef)
}

// muladdc multiplies a and b,
func muladdcConst(loRef Register, a Op, b Register, carry Register) {
	MOVQ(a, RAX)
	MULQ(b)

	ADDQ(loRef, RAX)
	ADCQ(Imm(0), RDX)
	ADDQ(carry, RAX)
	ADCQ(Imm(0), RDX)

	MOVQ(RDX, carry)
	MOVQ(RAX, loRef)
}

var QFieldModulus = [6]uint64{0xb9feffffffffaaab, 0x1eabfffeb153ffff, 0x6730d2a0f6b0f624, 0x64774b84f38512bf, 0x4b1ba7b6434bacd7, 0x1a0111ea397fe69a}

const montInvFQ = uint64(0x89f3fffcfffcfffd)

func r(i int) Component {
	if i >= 6 {
		return Param("hi").Index(i % 6)
	} else {
		return Param("lo").Index(i)
	}
}

func main() {
	Package("github.com/phoreproject/bls")
	Implement("MACWithCarry")
	Doc("Finds a + b * c + carry and returns the result and the carry.")
	b := Load(Param("b"), GP64())
	c := Load(Param("c"), RAX)
	Comment("Multiply b and c")
	MULQ(b)
	a := Load(Param("a"), GP64())
	Comment("Add a")
	ADDQ(a, c)
	Comment("Add to result carry if needed")
	ADCQ(Imm(0), RDX)
	carry := Load(Param("carry"), GP64())
	Comment("Add input carry to running result")
	ADDQ(carry, c)
	Comment("Add to result carry if needed")
	ADCQ(Imm(0), RDX)
	Store(c, ReturnIndex(0))
	Store(RDX, ReturnIndex(1))
	RET()

	Implement("SubWithBorrow")
	Doc("Finds a - b - borrow and returns the result and the borrow.")
	a = Load(Param("a"), GP64())
	b = Load(Param("b"), GP64())
	newBorrow1 := GP64()
	Comment("a = a - b")
	XORQ(newBorrow1, newBorrow1)
	SUBQ(b, a)
	Comment("Zero out borrow1 and set if overflowed")
	SETCS(newBorrow1.As8())
	borrow := Load(Param("borrow"), GP64())
	Comment("a = a - borrow")
	newBorrow2 := GP64()
	XORQ(newBorrow2, newBorrow2)
	SUBQ(borrow, a)
	Comment("Zero out borrow2 and set if overflowed")
	SETCS(newBorrow2.As8())
	Comment("borrow2 = borrow2 | borrow1")
	ORQ(newBorrow1, newBorrow2)
	Store(a, ReturnIndex(0))
	Store(newBorrow2, ReturnIndex(1))
	RET()

	Implement("AddWithCarry")
	Doc("Finds a + b + carry and returns the result and the borrow.")
	a = Load(Param("a"), GP64())
	b = Load(Param("b"), GP64())
	carry = Load(Param("carry"), GP64())
	newCarry := GP64()
	Comment("Zero out new carry")
	XORQ(newCarry, newCarry)
	Comment("Add a + b")
	ADDQ(b, a)
	Comment("Add to new carry if needed")
	ADCQ(Imm(0), newCarry)
	Comment("Add old carry")
	ADDQ(carry, a)
	Comment("Add to new carry if needed")
	ADCQ(Imm(0), newCarry)
	Store(a, ReturnIndex(0))
	Store(newCarry, ReturnIndex(1))
	RET()

	/*
												 a  b  c  d  e  f
											x  g  h  i  j  k  l
											-------------------
												al bl cl dl el fl
										 ak bk ck dk ek fk
									aj bj cj dj ej fj
						   ai bi ci di ei fi
						ah bh ch dh eh fh
			 + ag bg cg dg eg fg
			 ----------------------------------
		  11 10 9  8  7  6  5  4  3  2  1  0

			1. r0, carry = MUL(f, l)
			2. r1, carry = ADD(MUL(e, l), carry)
			3. r2, carry = ADD(MUL(d, l), carry)
			4. r3, carry = ADD(MUL(c, l), carry)
			5. r4, carry = ADD(MUL(b, l), carry)
			6. r5, carry = ADD(MUL(a, l), carry)

			r1, carry = ADD(r1, MUL(f, k))
			r2, carry = ADD(r2, MUL(e, k), carry)
			r3, carry = ADD(r3, MUL(e, k), carry)
			r4, carry = ADD(r4, MUL(e, k), carry)
			r5, carry = ADD(r5, MUL(e, k), carry)
			r6, carry = ADD(r2, MUL(e, k), carry)
	*/

	Implement("MultiplyFQRepr")
	registers := make([]Register, 12)
	registersUsed := make([]bool, 12)
	for i := range registers {
		registers[i] = GP64()
		registersUsed[i] = false
	}

	for i := 0; i < 6; i++ {
		carry := GP64()
		Comment("carry = 0")
		XORQ(carry, carry)
		j := 0
		for j < 6 {
			ai := Load(Param("a").Index(i), RAX)
			bj := Load(Param("b").Index(j), RBX)
			if !registersUsed[i+j] {
				Commentf("registers[%d] = 0", i+j)
				XORQ(registers[i+j], registers[i+j])
				registersUsed[i+j] = true
			}
			Commentf("carry = ((registers[%d] + a[%d] * b[%d] + carry) >> 64) & 0xFFFFFFFFFFFFFFFF", i+j, i, j)
			Commentf("registers[%d] = (registers[%d] + a[%d] * b[%d] + carry) & 0xFFFFFFFFFFFFFFFF", i+j, i+j, i, j)
			muladdc(registers[i+j], ai, bj, carry)
			j++
		}
		Commentf("registers[%d] = carry", i+j)
		MOVQ(carry, registers[i+j])
		registersUsed[i+j] = true
	}

	for i := 0; i < 6; i++ {
		Commentf("lo[%d] = registers[%d]", i, i)
		Store(registers[i], Return("lo").Index(i))
	}
	for i := 0; i < 6; i++ {
		Commentf("hi[%d] = registers[%d]", i, i+6)
		Store(registers[i+6], Return("hi").Index(i))
	}
	RET()

	Implement("MontReduce")

	Commentf("reg = [0] * 12")
	regs := AllocLocal(8 * 12)
	tempReg := GP64()
	for i := 0; i < 12; i++ {
		Commentf("temp = r[%d]", i)
		Load(r(i), tempReg)
		Commentf("reg[%d] = temp", i)
		MOVQ(tempReg, regs.Offset(8*i))
	}

	carryOver := GP64()
	Comment("carryOver = 0")
	XORQ(carryOver, carryOver)
	newCarry = GP64()
	lastReg := GP64()
	k := GP64()
	carry = GP64()

	for i := 0; i < 6; i++ {
		Commentf("rax = %d", montInvFQ)
		MOVQ(Imm(montInvFQ), RAX)
		Commentf("k = reg[%d]", i)
		MOVQ(regs.Offset(8*i), k)
		Commentf("rax = (rax * k) & 0xFFFFFFFFFFFFFFFF")
		MULQ(k)
		Commentf("k = rax")
		MOVQ(RAX, k)
		Commentf("carry = 0")
		XORQ(carry, carry)
		j := 0
		for j < 6 {
			regValue := GP64()
			Commentf("carryTemp = ((reg[%d] + QFieldModulus[%d] * k + carry) >> 64) & 0xFFFFFFFFFFFFFFFF", i+j, j)
			MOVQ(regs.Offset(8*(i+j)), regValue)
			muladdcConst(regValue, Imm(QFieldModulus[j]), k, carry)
			if j != 0 {
				Commentf("reg[%d] = (reg[%d] + QFieldModulus[%d] * k + carry) & 0xFFFFFFFFFFFFFFFF", i+j, i+j, j)
				MOVQ(regValue, regs.Offset(8*(i+j)))
			}
			Comment("carry = carryTemp")
			j++
		}

		// this should calculate lastReg + carryOver + carry
		Commentf("newCarry = 0")
		XORQ(newCarry, newCarry)
		Commentf("lastReg = reg[%d]", i+j)
		MOVQ(regs.Offset(8*(i+j)), lastReg)
		Comment("newCarry = ((lastReg + carry + carryOver) >> 64) & 0xFFFFFFFFFFFFFFFF")
		Comment("lastReg = (lastReg + carry + carryOver) & 0xFFFFFFFFFFFFFFFF")
		ADDQ(carry, lastReg)
		ADCQ(Imm(0), newCarry)
		ADDQ(carryOver, lastReg)
		ADCQ(Imm(0), newCarry)
		Commentf("carryOver = newCarry")
		MOVQ(newCarry, carryOver)
		Commentf("reg[%d] = lastReg", i+j)
		MOVQ(lastReg, regs.Offset(8*(i+j)))
	}
	tempReg = GP64()
	for i := 0; i < 6; i++ {
		MOVQ(regs.Offset(8*(i+6)), tempReg)
		Commentf("out[%d] = reg[%d]", i, i+6)
		Store(tempReg, Return("out").Index(i))
	}
	RET()
	Generate()
}
