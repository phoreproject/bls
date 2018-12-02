#include "textflag.h"

TEXT ·MACWithCarry(SB),NOSPLIT,$0
	MOVQ b+8(FP), AX
    MOVQ c+16(FP), BX
    MULQ BX
    MOVQ a+0(FP), BX
    ADDQ BX, AX
    ADCQ $0, DX
    MOVQ a+24(FP), CX
    ADDQ CX, AX
    ADCQ $0, DX
    MOVQ DX, carry+40(FP)
	MOVQ AX, ret+32(FP)
	RET

TEXT ·SubWithBorrow(SB),NOSPLIT,$0
	MOVQ a+0(FP), AX
    MOVQ a+16(FP), DX
    MOVQ b+8(FP), BX
    SUBQ BX, AX
    MOVQ $0, BX
    SETCS BX
    SUBQ DX, AX
    MOVQ $0, DX
    SETCS DX
    ORQ DX, BX
    MOVQ BX, carry+32(FP)
    MOVQ AX, ret+24(FP)
	RET

TEXT ·AddWithCarry(SB),NOSPLIT,$0
	MOVQ a+0(FP), AX
    MOVQ a+16(FP), DX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ $0, BX
    SETCS BX
    ADDQ DX, AX
    MOVQ $0, DX
    SETCS DX
    ORQ DX, BX
    MOVQ BX, carry+32(FP)
    MOVQ AX, ret+24(FP)
	RET
