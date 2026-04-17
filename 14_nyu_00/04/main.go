package main

import (
	"fmt"
	"math"
)

// --- DEFINICIONES ---
// Estas se guardan en la memoria del Kernel para otras celdas

type Layer struct {
	Weights []float64
	Bias    []float64
}

type Network struct {
	Layers []Layer
}

func (n *Network) Forward(input []float64) []float64 {
	output := input
	for _, layer := range n.Layers {
		output = applyLinear(output, layer.Weights, layer.Bias)
		output = applyReLU(output)
	}
	return output
}

func applyLinear(input, weights, bias []float64) []float64 {
	n_out := len(bias) // Número de neuronas en esta capa
	n_in := len(input) // Tamaño de la entrada
	result := make([]float64, n_out)

	for i := 0; i < n_out; i++ {
		sum := 0.0
		for j := 0; j < n_in; j++ {
			// weights[i*n_in+j] asume que los pesos están en un array plano
			sum += input[j] * weights[i*n_in+j]
		}
		result[i] = sum + bias[i]
	}
	return result
}

func applyReLU(input []float64) []float64 {
	result := make([]float64, len(input))
	for i := range input {
		result[i] = math.Max(input[i], 0.0)
	}
	return result
}

func main() {

	// --- EJECUCIÓN ---
	// Esto es lo que antes estaba dentro de func main()

	network := Network{
		Layers: []Layer{
			{
				Weights: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6},
				Bias:    []float64{0.0, 0.1},
			},
			{
				Weights: []float64{0.7, 0.8, 0.9, 1.0, 1.1, 1.2}, // Pesos para 2 neuronas de entrada
				Bias:    []float64{0.2, 0.3},
			},
		},
	}

	input := []float64{1.0, 2.0, 3.0}
	output := network.Forward(input)

	fmt.Println("🚀 Salida de la Red Neuronal:", output)
}
