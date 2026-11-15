//go:build riscv64

#include "textflag.h"

// Nota: Go 1.24 declara las constantes RVV 1.0 en cmd/internal/obj/riscv, pero aún
// no habilita su codificador mnemónico en el ensamblador Plan 9, por lo que
// emitimos los opcodes RVV 1.0 de 32 bits mediante WORD con su equivalencia exacta.

// func VectorAddRVV(dst, src []float32)
// Suma vectorial float32 en sitio: dst[i] += src[i] (RVV 1.0, e32/m1)
TEXT ·VectorAddRVV(SB), NOSPLIT, $0-48
	MOV  dst_base+0(FP), A0
	MOV  src_base+24(FP), A1
	MOV  dst_len+8(FP), A2
	BEQ  A2, ZERO, add_done

add_loop:
	// VSETVLI A2, E32, M1, TA, MA -> A3 (vl procesados)
	WORD $0x0d0676d7
	// VLE32.V V1, (A0)
	WORD $0x02056087
	// VLE32.V V2, (A1)
	WORD $0x0205e107
	// VFADD.VV V3, V1, V2
	WORD $0x021111d7
	// VSE32.V V3, (A0)
	WORD $0x020561a7

	SLLI $2, A3, T0
	ADD  T0, A0
	ADD  T0, A1
	SUB  A3, A2
	BNE  A2, ZERO, add_loop

add_done:
	RET

// func DotProductRVV(a, b []float32) float32
// Producto escalar vectorial (y = W·x) con Fused Multiply-Add RVV 1.0
TEXT ·DotProductRVV(SB), NOSPLIT, $0-52
	MOV    a_base+0(FP), A0
	MOV    b_base+24(FP), A1
	MOV    a_len+8(FP), A2
	MOV    A2, A4
	FMVWX  ZERO, F0
	BEQ    A2, ZERO, dot_ret

	// VSETVLI A4, E32, M1, TU, MU -> T0 (inicializa VLMAX para limpiar V3)
	WORD $0x010772d7
	// VFMV.V.F V3, F0 (acumulador V3 = 0.0)
	WORD $0x5e0051d7

dot_loop:
	// VSETVLI A2, E32, M1, TU, MU -> A3 (preserva cola de V3 en la última iteración)
	WORD $0x010676d7
	// VLE32.V V1, (A0)
	WORD $0x02056087
	// VLE32.V V2, (A1)
	WORD $0x0205e107
	// VFMACC.VV V3, V1, V2 (V3 += V1 * V2)
	WORD $0xb22091d7

	SLLI $2, A3, T0
	ADD  T0, A0
	ADD  T0, A1
	SUB  A3, A2
	BNE  A2, ZERO, dot_loop

	// Reducción horizontal de V3 a escalar en F10 (FA0)
	// VSETVLI A4, E32, M1, TU, MU -> T0
	WORD $0x010772d7
	// VFMV.V.F V4, F0 (identidad 0.0 en V4[0])
	WORD $0x5e005257
	// VFREDUSUM.VS V4, V3, V4
	WORD $0x06321257
	// VFMV.F.S F10, V4
	WORD $0x42401557
	MOVF F10, ret+48(FP)
	RET

dot_ret:
	MOVF F0, ret+48(FP)
	RET
