package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func main() {

	// 1. Definir entrada (Vector columna de 3 elementos)
	x := mat.NewVecDense(3, []float64{0, 1, 1})

	// 2. Definir pesos (Matriz 2 filas, 3 columnas)
	weights := mat.NewDense(1, 3, []float64{
		-10, 10, 5,
	})

	// 3. Definir sesgo
	bias := mat.NewVecDense(1, []float64{-5})

	// 4. PREPARAR EL RESULTADO
	// En Gonum, Mul() no devuelve una matriz nueva,
	// necesita una ya existente para escribir el resultado.
	result := mat.NewVecDense(1, nil)

	// 5. OPERACIÓN: result = weights * x
	// (2x3) * (3x1) = (2x1)
	result.MulVec(weights, x)

	// 6. OPERACIÓN: output = result + bias
	output := mat.NewVecDense(1, nil)
	output.AddVec(result, bias)

	fmt.Printf("🚀 Salida con Gonum: %v\n", mat.Formatted(output))
}
