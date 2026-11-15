package main

import (
	"fmt"
	"math"
	"runtime"
)

func sigmoid32(x float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(float64(-x))))
}

// NeuronRVV evalúa y = sigmoid(W·x + b) delegando el producto escalar en RVV 1.0 (en riscv64)
// o en el fallback portable (en amd64/arm64).
type NeuronRVV struct {
	Weights []float32
	Bias    float32
}

func (n *NeuronRVV) Forward(inputs []float32) float32 {
	z := DotProductRVV(n.Weights, inputs) + n.Bias
	return sigmoid32(z)
}

func main() {
	fmt.Println("=== Bloque 5: Edge AI en Silicio Libre (RISC-V + RVV 1.0) ===")
	fmt.Printf("Arquitectura : %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Motor activo : %s\n\n", EngineBackend)

	// 1. Demostración de suma vectorial (e.g. actualización de sesgos / residuales)
	dst := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	src := []float32{0.5, 0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 1.0}
	VectorAddRVV(dst, src)
	fmt.Printf("VectorAddRVV resultado: %v\n", dst)

	// 2. Red neuronal de 2 capas (XOR) evaluada con DotProductRVV
	hidden1 := NeuronRVV{Weights: []float32{20.0, 20.0}, Bias: -10.0}  // OR
	hidden2 := NeuronRVV{Weights: []float32{-20.0, -20.0}, Bias: 30.0} // NAND
	output := NeuronRVV{Weights: []float32{20.0, 20.0}, Bias: -30.0}   // AND

	cases := [][]float32{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}
	expected := []float32{0, 1, 1, 0}

	fmt.Println("\nInferencia XOR sobre capa vectorial:")
	for i, x := range cases {
		h := []float32{hidden1.Forward(x), hidden2.Forward(x)}
		y := output.Forward(h)
		fmt.Printf("  Input [%.0f, %.0f] -> Esperado: %.0f | Predicción: %.4f\n", x[0], x[1], expected[i], y)
	}

	fmt.Println("\nPara compilar binario nativo CanMV-K230 (RVV 1.0):")
	fmt.Println("  CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 GORISCV64=rva22u64 go build -ldflags=\"-s -w\" -o neuron_k230 .")
}
