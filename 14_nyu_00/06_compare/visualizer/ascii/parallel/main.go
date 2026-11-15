package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	matrixSize = 10
	delay      = 100 * time.Millisecond
)

func main() {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println("=== GPU Parallelism Terminal Visualizer ===")
	fmt.Println("Representación del proceso de MatMul en paralelo")
	fmt.Println("A (Input) @ B (Weights) = C (Output)")
	fmt.Println("-------------------------------------------")

	// Inicializar matrices con símbolos
	matrixA := make([][]string, matrixSize)
	matrixB := make([][]string, matrixSize)
	matrixC := make([][]string, matrixSize)

	for i := 0; i < matrixSize; i++ {
		matrixA[i] = make([]string, matrixSize)
		matrixB[i] = make([]string, matrixSize)
		matrixC[i] = make([]string, matrixSize)
		for j := 0; j < matrixSize; j++ {
			matrixA[i][j] = "░"
			matrixB[i][j] = "░"
			matrixC[i][j] = " "
		}
	}

	// Simular el lanzamiento de hilos
	for step := 0; step < matrixSize*matrixSize; step += 3 {
		fmt.Print("\033[H") // Cursor al inicio
		fmt.Printf("\nLanzando hilos de GPU... Paso %d/%d\n\n", step, matrixSize*matrixSize)

		// Activar 3 hilos aleatorios por "ciclo"
		for k := 0; k < 3; k++ {
			r := rand.Intn(matrixSize)
			c := rand.Intn(matrixSize)
			matrixC[r][c] = "█" // Resultado calculado
		}

		// Pintar estado
		printMatrices(matrixA, matrixB, matrixC)
		time.Sleep(delay)
	}

	fmt.Println("\n\nCálculo completado en la GPU!")
}

func printMatrices(a, b, c [][]string) {
	for i := 0; i < matrixSize; i++ {
		// Fila de A
		fmt.Print(strings.Join(a[i], ""))
		if i == matrixSize/2 {
			fmt.Print("  @  ")
		} else {
			fmt.Print("     ")
		}

		// Fila de B
		fmt.Print(strings.Join(b[i], ""))

		if i == matrixSize/2 {
			fmt.Print("  =  ")
		} else {
			fmt.Print("     ")
		}

		// Fila de C
		fmt.Print("[")
		fmt.Print(strings.Join(c[i], ""))
		fmt.Println("]")
	}
}
