package main

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Neuron represents a classical neuron with weights, bias, and activation function
type Neuron struct {
	weights []float64
	bias    float64
}

// NewNeuron creates a new neuron with given weights and bias
func NewNeuron(weights []float64, bias float64) *Neuron {
	return &Neuron{
		weights: weights,
		bias:    bias,
	}
}

// Compute returns the output of the neuron given input x
func (n *Neuron) Compute(x []float64) float64 {
	// Compute the weighted sum of inputs
	weightedSum := n.bias
	for i := 0; i < len(x); i++ {
		weightedSum += n.weights[i] * x[i]
	}

	// Apply the activation function (sigmoid)
	return 1 / (1 + math.Exp(-weightedSum))
}

// sigmoid is the activation function
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func main() {
	// Create a new plot
	p := plot.New()
	p.Title.Text = "Sigmoid Function vs Neuron Output"
	p.X.Label.Text = "Input (x)"
	p.Y.Label.Text = "Output"

	// Generate input values from -10 to 10
	inputs := make([]float64, 100)
	for i := 0; i < 100; i++ {
		inputs[i] = float64(i-50) / 5.0 // Rango de -10 a 10
	}

	// Generate sigmoid function values
	sigValues := make([]float64, 100)
	for i := 0; i < 100; i++ {
		sigValues[i] = sigmoid(inputs[i])
	}

	// Generate neuron output values
	neuron := NewNeuron([]float64{9.0}, 0.7) // Example: weight = 2, bias = 0.1
	neuronOutput := make([]float64, 100)
	for i := 0; i < 100; i++ {
		neuronOutput[i] = neuron.Compute([]float64{inputs[i]})
	}

	// Create XY points for both functions
	sigPoints := make(plotter.XYs, len(sigValues))
	neuronPoints := make(plotter.XYs, len(neuronOutput))

	for i := range sigValues {
		sigPoints[i].X = inputs[i]
		sigPoints[i].Y = sigValues[i]
		neuronPoints[i].X = inputs[i]
		neuronPoints[i].Y = neuronOutput[i]
	}

	// Create a line plotter for the sigmoid function (blue)
	sigLine, err := plotter.NewLine(sigPoints)
	if err != nil {
		panic(err)
	}
	sigLine.Color = color.RGBA{255, 0, 0, 255} // Rojo

	// Create a line plotter for the neuron output (red)
	neuronLine, err := plotter.NewLine(neuronPoints)
	if err != nil {
		panic(err)
	}
	neuronLine.Color = color.RGBA{0, 0, 255, 255} // Azul

	// Add both lines to the plot
	p.Add(sigLine)
	p.Add(neuronLine)

	// Save the plot to a file
	if err := p.Save(800, 600, "sigmoid_vs_neuron.png"); err != nil {
		panic(err)
	}
}
