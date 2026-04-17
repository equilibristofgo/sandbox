package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxVRAM       = 12.0 // GB (Arc B580)
	hiddenDim     = 4096 // Llama-2/3 7B style
	numLayers     = 32
	bytesPerParam = 0.5 // INT4 (4 bits = 0.5 bytes)
	kvBytes       = 4   // (2 for K, 2 for V in FP16)
)

func main() {
	fmt.Print("\033[H\033[2J") // Clear
	fmt.Println("=== VRAM & CONTEXT GROWTH SIMULATOR (12GB ARC B580) ===")
	fmt.Println("Simulando carga de pesos + crecimiento de KV Cache por Token")
	fmt.Println("------------------------------------------------------------")

	// Peso fijo del modelo (Weights)
	// 7B parámetros * 0.5 bytes (INT4) = ~3.5GB
	weightMem := 3.5

	fmt.Printf("1. Cargando Pesos del Modelo (INT4, 7B)... %.2f GB\n", weightMem)
	time.Sleep(1 * time.Second)

	for seqLen := 1; seqLen <= 4096; seqLen += 64 {
		// Cálculo de KV Cache:
		// Formula: 2 (K+V) * num_layers * hidden_dim * seq_len * precision_bytes
		// Simplificado para la simulación:
		kvMem := (float64(seqLen) * float64(hiddenDim) * float64(numLayers) * 2 * 2) / (1024 * 1024 * 1024)
		totalMem := weightMem + kvMem

		fmt.Print("\033[H")
		fmt.Printf("\n=== MONITOR DE VRAM (GPU ARC B580) ===\n")
		fmt.Printf("Contexto actual: %d tokens\n", seqLen)
		fmt.Printf("Pesos: %.2f GB | KV Cache: %.2f GB\n", weightMem, kvMem)
		fmt.Printf("Uso Total: %.2f / %.2f GB (%.1f%%)\n\n", totalMem, maxVRAM, (totalMem/maxVRAM)*100)

		drawTank(totalMem, maxVRAM)

		if totalMem > maxVRAM {
			fmt.Println("\n\n[!!!] CRASH: OUT OF MEMORY (OOM) [!!!]")
			fmt.Println("El contexto ha excedido la capacidad física de la GPU.")
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func drawTank(current, total float64) {
	width := 50
	filled := int((current / total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	fmt.Printf("[%s]\n", bar)

	// Marcadores
	fmt.Print("0GB")
	fmt.Print(strings.Repeat(" ", width/2-3))
	fmt.Print("6GB")
	fmt.Print(strings.Repeat(" ", width/2-3))
	fmt.Println("12GB")
}
