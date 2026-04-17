package main

import (
	"fmt"
	"math"
)

// Sigmoid returns the logistic sigmoid of z
func Sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

func main() {
	// Ejemplo de uso de la función Sigmoid
	values := []float64{-10, -5, 0, 5, 10}
	fmt.Println("z\t\tSigmoid(z)")
	for _, z := range values {
		result := Sigmoid(z)
		fmt.Printf("%.2f\t%.6f\n", z, result)
	}
}
