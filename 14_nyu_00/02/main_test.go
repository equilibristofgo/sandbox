package main

import (
	"00_02_matrix/matrix_avx"
	"fmt"
	"testing"
)

func VecAddGo(a, b, res []float32) {
	for i := range a {
		res[i] = a[i] + b[i]
	}
}

func BenchmarkVecAdd(b *testing.B) {
	sizes := []int{1024, 65536, 1048576}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Go-%d", size), func(b *testing.B) {
			vecA := make([]float32, size)
			vecB := make([]float32, size)
			res := make([]float32, size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				VecAddGo(vecA, vecB, res)
			}
		})
		b.Run(fmt.Sprintf("AVX-%d", size), func(b *testing.B) {
			vecA := make([]float32, size)
			vecB := make([]float32, size)
			res := make([]float32, size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matrix_avx.VecAddAVX(vecA, vecB, res)
			}
		})
	}
}
