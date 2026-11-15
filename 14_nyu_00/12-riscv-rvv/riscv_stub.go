//go:build !riscv64

package main

import "runtime"

const EngineBackend = "Simulación Go Fallback (Host " + runtime.GOARCH + " — Cross-compilar con GOARCH=riscv64 para RVV 1.0)"

// VectorAddRVV simula en CPU host el bucle Vector-Length Agnostic de RVV 1.0.
func VectorAddRVV(dst, src []float32) {
	n := len(dst)
	if len(src) < n {
		n = len(src)
	}
	for i := 0; i < n; i++ {
		dst[i] += src[i]
	}
}

// DotProductRVV calcula el producto escalar a·b en Go estándar cuando no corre en riscv64.
func DotProductRVV(a, b []float32) float32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var sum float32
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	return sum
}

