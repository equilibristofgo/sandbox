package main

import (
	"fmt"
	"math"
)

type Layer struct {
	Weights []float64
	Bias    []float64
}

type Network struct {
	Layers []Layer
}

func (n *Network) Forward(input []float64) []float64 {
	output := input
	for i, layer := range n.Layers {
		output = applyLinear(output, layer.Weights, layer.Bias)
		output = applyReLU(output)
		fmt.Printf("Capa %d procesada. Salida: %v\n", i+1, output)
	}
	return output
}

func applyLinear(input, weights, bias []float64) []float64 {
	n_out := len(bias)
	n_in := len(input)
	result := make([]float64, n_out)

	for i := 0; i < n_out; i++ {
		sum := 0.0
		for j := 0; j < n_in; j++ {
			// El índice debe estar dentro de len(weights)
			idx := i*n_in + j
			sum += input[j] * weights[idx]
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
	network := Network{
		Layers: []Layer{
			{
				// Capa 1: Entrada 3 -> Salida 2 (Necesita 3*2 = 6 pesos)
				Weights: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6},
				Bias:    []float64{0.0, 0.1},
			},
			{
				// Capa 2: Entrada 2 -> Salida 2 (Necesita 2*2 = 4 pesos)
				// Cambiado de 6 a 4 elementos para que encaje
				Weights: []float64{0.7, 0.8, 0.9, 1.0},
				Bias:    []float64{0.2, 0.3},
			},
		},
	}

	input := []float64{1.0, 2.0, 3.0}
	output := network.Forward(input)

	fmt.Println("🚀 Resultado final:", output)
}
