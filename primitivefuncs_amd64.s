#include "textflag.h"

TEXT ·MACWithCarry(SB),NOSPLIT,$0
	MOVQ b+8(FP), AX
    MOVQ c+16(FP), BX
    MULQ BX
    MOVQ a+0(FP), BX
    ADDQ BX, AX
    ADCQ $0, DX
    MOVQ a+24(FP), CX
    ADDQ (CX), AX
    ADCQ $0, DX
	MOVQ AX, ret+32(FP)
    MOVQ DX, (CX)
	RET

TEXT ·SubWithBorrow(SB),NOSPLIT,$0
	MOVQ a+0(FP), AX
    MOVQ a+16(FP), CX
    MOVQ (CX), DX
    MOVQ b+8(FP), BX
    SUBQ BX, AX
    SETCS BX
    SUBQ DX, AX
    SETCS DX
    ORQ DX, BX
    MOVQ BX, (CX)
    MOVQ AX, ret+24(FP)
	RET

TEXT ·AddWithCarry(SB),NOSPLIT,$0
	MOVQ a+0(FP), AX
    MOVQ a+16(FP), CX
    MOVQ (CX), DX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    SETCS BX
    ADDQ DX, AX
    SETCS DX
    ORQ DX, BX
    MOVQ BX, (CX)
    MOVQ AX, ret+24(FP)
	RET
