//go:build riscv64

package main

// VectorAddRVV ejecuta dst[i] += src[i] en ensamblador RVV 1.0 (vector_riscv64.s).
func VectorAddRVV(dst, src []float32)

// DotProductRVV calcula el producto escalar a·b usando VFMACC.VV y VFREDUSUM.VS.
func DotProductRVV(a, b []float32) float32

const EngineBackend = "RISC-V 64 RVV 1.0 (Ensamblador Nativo T-Head C908 / VPU 128-bit)"

