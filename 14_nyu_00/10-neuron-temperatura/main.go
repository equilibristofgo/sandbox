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
	p.Title.Text = "Neuron Output vs Hour of Day"
	p.X.Label.Text = "Hour of Day (0 to 24)"
	p.Y.Label.Text = "Predicted Temperature (0 to 1)"

	// Generate input values from 0 to 24
	inputs := make([]float64, 25)
	for i := 0; i < 25; i++ {
		inputs[i] = float64(i)
	}

	// Generate neuron output values with bias = 0 and bias = 0.5
	neuron1 := NewNeuron([]float64{0.1}, 0.0) // Bias = 0
	neuron2 := NewNeuron([]float64{0.1}, 0.5) // Bias = 0.5

	// Generate outputs for both neurons
	neuron1Output := make([]float64, 25)
	neuron2Output := make([]float64, 25)

	for i := 0; i < 25; i++ {
		neuron1Output[i] = neuron1.Compute([]float64{inputs[i]})
		neuron2Output[i] = neuron2.Compute([]float64{inputs[i]})
	}

	// Create XY points for both functions
	neuron1Points := make(plotter.XYs, len(neuron1Output))
	neuron2Points := make(plotter.XYs, len(neuron2Output))

	for i := range neuron1Output {
		neuron1Points[i].X = inputs[i]
		neuron1Points[i].Y = neuron1Output[i]
		neuron2Points[i].X = inputs[i]
		neuron2Points[i].Y = neuron2Output[i]
	}

	// Create a line plotter for the first neuron (bias = 0)
	neuron1Line, err := plotter.NewLine(neuron1Points)
	if err != nil {
		panic(err)
	}
	neuron1Line.Color = color.RGBA{255, 0, 0, 255} // Red

	// Create a line plotter for the second neuron (bias = 0.5)
	neuron2Line, err := plotter.NewLine(neuron2Points)
	if err != nil {
		panic(err)
	}
	neuron2Line.Color = color.RGBA{0, 0, 255, 255} // Blue

	// Add both lines to the plot
	p.Add(neuron1Line)
	p.Add(neuron2Line)

	// Save the plot to a file
	if err := p.Save(800, 600, "temperature_prediction.png"); err != nil {
		panic(err)
	}
}
