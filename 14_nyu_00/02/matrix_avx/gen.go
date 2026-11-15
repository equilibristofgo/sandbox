//go:build ignore
// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("VecAddAVX", NOSPLIT, "func(a, b, res []float32)")
	Doc("VecAddAVX suma los elementos de a y b usando registros YMM de 256 bits.")

	// Cargamos los punteros a los slices
	ptrA := Load(Param("a").Base(), GP64())
	ptrB := Load(Param("b").Base(), GP64())
	ptrRes := Load(Param("res").Base(), GP64())

	n := Load(Param("a").Len(), GP64())

	// Bucle unrolled 4x (procesa 32 floats por iteración)
	Label("loop_unrolled")
	CMPQ(n, Imm(32))
	JL(LabelRef("loop_8"))

	// Iteración de 32 elementos (4 registros YMM)
	y0, y1, y2, y3 := YMM(), YMM(), YMM(), YMM()
	yb0, yb1, yb2, yb3 := YMM(), YMM(), YMM(), YMM()

	VMOVUPS(Mem{Base: ptrA}, y0)
	VMOVUPS(Mem{Base: ptrA, Disp: 32}, y1)
	VMOVUPS(Mem{Base: ptrA, Disp: 64}, y2)
	VMOVUPS(Mem{Base: ptrA, Disp: 96}, y3)

	VMOVUPS(Mem{Base: ptrB}, yb0)
	VMOVUPS(Mem{Base: ptrB, Disp: 32}, yb1)
	VMOVUPS(Mem{Base: ptrB, Disp: 64}, yb2)
	VMOVUPS(Mem{Base: ptrB, Disp: 96}, yb3)

	VADDPS(y0, yb0, y0)
	VADDPS(y1, yb1, y1)
	VADDPS(y2, yb2, y2)
	VADDPS(y3, yb3, y3)

	VMOVUPS(y0, Mem{Base: ptrRes})
	VMOVUPS(y1, Mem{Base: ptrRes, Disp: 32})
	VMOVUPS(y2, Mem{Base: ptrRes, Disp: 64})
	VMOVUPS(y3, Mem{Base: ptrRes, Disp: 96})

	ADDQ(Imm(128), ptrA)
	ADDQ(Imm(128), ptrB)
	ADDQ(Imm(128), ptrRes)
	SUBQ(Imm(32), n)
	JMP(LabelRef("loop_unrolled"))

	// Bucle para bloques de 8 elementos (resto)
	Label("loop_8")
	CMPQ(n, Imm(8))
	JL(LabelRef("done"))

	y := YMM()
	yb := YMM()
	VMOVUPS(Mem{Base: ptrA}, y)
	VMOVUPS(Mem{Base: ptrB}, yb)
	VADDPS(y, yb, y)
	VMOVUPS(y, Mem{Base: ptrRes})

	ADDQ(Imm(32), ptrA)
	ADDQ(Imm(32), ptrB)
	ADDQ(Imm(32), ptrRes)
	SUBQ(Imm(8), n)
	JMP(LabelRef("loop_8"))

	Label("done")
	RET()
	Generate()
}
