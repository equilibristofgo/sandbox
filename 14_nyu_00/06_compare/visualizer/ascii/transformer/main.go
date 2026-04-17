package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	matrixSize = 8
	delay      = 800 * time.Millisecond
)

type Stage struct {
	Name    string
	MatrixA string
	MatrixB string
	MatrixC string
	Status  string
}

func main() {
	fmt.Print("\033[H\033[2J") // Clear
	fmt.Println("=== TRANSFORMER STAGE & KV CACHE SIMULATOR ===")
	fmt.Println("Simulando el flujo de datos y el 'Cache de Multiplicaciones'")
	fmt.Println("----------------------------------------------------------")

	stages := []Stage{
		{
			Name:    "Stage 1: Q, K, V Projections",
			MatrixA: "Tokens",
			MatrixB: "Weights_Proj",
			MatrixC: "Q, K, V Vectors",
			Status:  "Calculando proyecciones iniciales...",
		},
		{
			Name:    "Stage 2: Self-Attention (Q @ K^T)",
			MatrixA: "Query",
			MatrixB: "Key (Cache)",
			MatrixC: "Attn Scores",
			Status:  "¡HIT! Usando K de la caché del token anterior.",
		},
		{
			Name:    "Stage 3: MLP / Feed-Forward",
			MatrixA: "Attn Output",
			MatrixB: "Weights_MLP",
			MatrixC: "Next States",
			Status:  "Procesamiento final de la capa.",
		},
	}

	for _, s := range stages {
		renderStage(s)
		time.Sleep(delay * 2)
	}

	fmt.Println("\n\n=== EXPLICACIÓN DEL KV CACHE ===")
	fmt.Println("1. En cada token nuevo, NO recalculamos los Key/Value de los anteriores.")
	fmt.Println("2. Los guardamos en VRAM (KV Cache).")
	fmt.Println("3. Esto ahorra trillones de multiplicaciones en secuencias largas.")
	fmt.Println("----------------------------------------------------------")
}

func renderStage(s Stage) {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("\n>>> %s\n", s.Name)
	fmt.Printf("Estado: %s\n\n", s.Status)

	// Dibujar esquema visual
	gridA := createGrid("█")
	gridB := createGrid("▓")
	gridC := createGrid("░")

	if strings.Contains(s.Status, "HIT") {
		gridB = createGrid("✔") // Representar cache hit
	}

	for i := 0; i < matrixSize; i++ {
		// Matrix A
		fmt.Printf("  %s  ", gridA[i])
		if i == matrixSize/2 {
			fmt.Print("  @  ")
		} else {
			fmt.Print("     ")
		}

		// Matrix B
		fmt.Printf("  %s  ", gridB[i])

		if i == matrixSize/2 {
			fmt.Print("  =  ")
		} else {
			fmt.Print("     ")
		}

		// Matrix C (poblándose)
		fmt.Printf(" [%s]\n", gridC[i])
	}

	fmt.Printf("\n  %-15s   %-15s   %-15s\n", s.MatrixA, s.MatrixB, s.MatrixC)
}

func createGrid(char string) []string {
	res := make([]string, matrixSize)
	for i := 0; i < matrixSize; i++ {
		res[i] = strings.Repeat(char, matrixSize)
	}
	return res
}
