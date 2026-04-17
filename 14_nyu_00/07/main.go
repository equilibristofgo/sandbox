package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// 1. Datos de entrada
	x := mat.NewVecDense(3, []float64{1.0, 2.0, 3.0})

	// 2. Pesos (Matriz 2x3)
	weights := mat.NewDense(2, 3, []float64{
		0.1, 0.2, 0.3,
		0.4, 0.5, 0.6,
	})

	// 3. Sesgo
	bias := mat.NewVecDense(2, []float64{0.0, 0.1})

	// 4. Contenedor para el resultado (Debe ser del tamaño de la salida: 2)
	// Usamos NewVecDense(2, nil) para reservar el espacio
	result := mat.NewVecDense(2, nil)

	// 5. Multiplicación de Matriz por Vector: result = weights * x
	result.MulVec(weights, x)

	// 6. Suma de vectores: result = result + bias
	result.AddVec(result, bias)

	fmt.Printf("🚀 Resultado con Gonum:\n%v\n", mat.Formatted(result))
}
