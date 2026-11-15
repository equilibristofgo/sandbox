package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lmatrix_ze -Wl,-rpath,. -lze_loader
#include "matrix_ze.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"time"
)

func matrixMultiplyZE(A, B []float32, N int) []float32 {
	res := make([]float32, N*N)
	C.matrix_multiply_ze((*C.float)(&A[0]), (*C.float)(&B[0]), (*C.float)(&res[0]), C.int(N))
	return res
}

func matrixMultiplyCPU(A, B []float32, N int) []float32 {
	res := make([]float32, N*N)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			var sum float32
			for k := 0; k < N; k++ {
				sum += A[i*N+k] * B[k*N+j]
			}
			res[i*N+j] = sum
		}
	}
	return res
}

func main() {
	const N = 2048 // Aumentado para comparar con SYCL
	fmt.Printf("Level Zero Matrix Size: %dx%d\n", N, N)

	A := make([]float32, N*N)
	B := make([]float32, N*N)

	for i := range A {
		A[i] = rand.Float32()
		B[i] = rand.Float32()
	}

	fmt.Println("Starting GPU (Level Zero) multiplication...")
	start := time.Now()
	resGPU := matrixMultiplyZE(A, B, N)
	gpuTime := time.Since(start)
	fmt.Printf("GPU Time: %v\n", gpuTime)

	fmt.Println("Starting CPU (Go) multiplication...")
	start = time.Now()
	resCPU := matrixMultiplyCPU(A, B, N)
	cpuTime := time.Since(start)
	fmt.Printf("CPU Time: %v\n", cpuTime)

	// Verification
	match := true
	for i := range resGPU {
		diff := resGPU[i] - resCPU[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > 1e-3 {
			fmt.Printf("Mismatch at index %d: GPU=%f, CPU=%f\n", i, resGPU[i], resCPU[i])
			match = false
			break
		}
	}

	if match {
		fmt.Println("Verification SUCCESS: Matrices match!")
	} else {
		fmt.Println("Verification FAILED: Matrices do not match.")
	}

	speedup := float64(cpuTime) / float64(gpuTime)
	fmt.Printf("Speedup: %.2fx\n", speedup)
}
