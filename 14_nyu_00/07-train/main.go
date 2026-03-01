package main

import (
	"fmt"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

// --- FASE DE DEFINICIÓN SENSORIAL ---

func sigmoid(x float64) float64 { return 1.0 / (1.0 + math.Exp(-x)) }

// La derivada de la sigmoide: El "sensor de cambio"
func sigmoidPrime(x float64) float64 {
	s := sigmoid(x)
	return s * (1.0 - s)
}

func main() {
	// 1. INPUTS (El Plano 2D): 4 patrones de XOR
	inputs := mat.NewDense(4, 2, []float64{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
	})

	// 2. TARGETS (La Verdad): Lo que esperamos
	targets := mat.NewDense(4, 1, []float64{0, 1, 1, 0})

	// 3. ARQUITECTURA (Los Pesos): Inicialización aleatoria (Ruido inicial)
	// Capa Oculta: 2 neuronas | Capa Salida: 1 neurona
	wHidden := mat.NewDense(2, 2, nil)
	wOutput := mat.NewDense(2, 1, nil)
	applyRand := func(i, j int, v float64) float64 { return rand.Float64() }
	wHidden.Apply(applyRand, wHidden)
	wOutput.Apply(applyRand, wOutput)

	learningRate := 0.5

	// --- FASE DE ENTRENAMIENTO (REENTRENAMIENTO QLoRA) ---
	for epoch := 0; epoch < 10000; epoch++ {

		// FORWARD PASS: La Alucinación del Sistema 1
		var hiddenLayerInput, hiddenLayerOutput mat.Dense
		hiddenLayerInput.Mul(inputs, wHidden)
		hiddenLayerOutput.Apply(func(i, j int, v float64) float64 { return sigmoid(v) }, &hiddenLayerInput)

		var finalInput, finalOutput mat.Dense
		finalInput.Mul(&hiddenLayerOutput, wOutput)
		finalOutput.Apply(func(i, j int, v float64) float64 { return sigmoid(v) }, &finalInput)

		// BACKPROPAGATION: El "Error de Predicción" de Barrett
		// Calcular error de salida
		var outputError mat.Dense
		outputError.Sub(targets, &finalOutput)

		// Gradiente de salida (Cúanto "duele" el error)
		var dOutput mat.Dense
		finalOutput.Apply(func(i, j int, v float64) float64 { return v * (1 - v) }, &finalOutput)
		dOutput.MulElem(&outputError, &finalOutput)

		// Error de la capa oculta (Propagando la culpa hacia atrás)
		var hiddenError mat.Dense
		hiddenError.Mul(&dOutput, wOutput.T())

		var dHidden mat.Dense
		hiddenLayerOutput.Apply(func(i, j int, v float64) float64 { return v * (1 - v) }, &hiddenLayerOutput)
		dHidden.MulElem(&hiddenError, &hiddenLayerOutput)

		// ACTUALIZACIÓN DE PESOS: El Adaptador QLoRA mental
		var adjOutput, adjHidden mat.Dense
		adjOutput.Mul(hiddenLayerOutput.T(), &dOutput)
		adjOutput.Scale(learningRate, &adjOutput)
		wOutput.Add(wOutput, &adjOutput)

		adjHidden.Mul(inputs.T(), &dHidden)
		adjHidden.Scale(learningRate, &adjHidden)
		wHidden.Add(wHidden, &adjHidden)

		if epoch%2000 == 0 {
			fmt.Printf("Epoch %d - Error sum: %.4f\n", epoch, mat.Sum(&outputError))
		}
	}

	// RESULTADO FINAL
	fmt.Println("\nPredicciones finales:")
	// (Aquí repetirías un Forward Pass final para ver los resultados cercanos a 0 y 1)
}
