package main

import "fmt"

type Perceptron struct {
	Weights []float64
	Bias    float64
}

func (p *Perceptron) Forward(inputs []float64) float64 {
	sum := p.Bias
	for i, input := range inputs {
		sum += input * p.Weights[i]
	}
	return sum // Aquí faltaría la función de activación, ¡la verás en el curso!
}

func main() {
	p := Perceptron{Weights: []float64{0.5, -0.2}, Bias: 0.1}
	resultado := p.Forward([]float64{1.0, 1.0})
	fmt.Printf("🚀 Resultado del Perceptrón en Go: %v\n", resultado)
}
