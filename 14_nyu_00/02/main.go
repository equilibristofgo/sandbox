package main

import (
	"00_02_matrix/matrix_avx"
	"fmt"
)

func main() {
	// Definición de un "Tensor" manual (matriz 2x2)
	matrixA := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	fmt.Println("🧠 Mi primera matriz en Go para el curso de NYU:")
	for _, row := range matrixA {
		fmt.Println(row)
	}

	fmt.Println("\n🚀 Probando AVX (Suma de Vectores):")

	a := []float32{1, 2, 3, 4, 5, 6, 7, 8}
	b := []float32{10, 20, 30, 40, 50, 60, 70, 80}
	res := make([]float32, 8)

	matrix_avx.VecAddAVX(a, b, res)

	fmt.Printf("Vector A:   %v\n", a)
	fmt.Printf("Vector B:   %v\n", b)
	fmt.Printf("Resultado:  %v\n", res)
}
