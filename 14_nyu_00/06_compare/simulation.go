package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== COMPARATIVA DE RENDIMIENTO (LIVE) ===")
	runBenchmark("SYCL", "../06_oneapi")
	runBenchmark("Level Zero", "../06_level_zero")

	fmt.Println("\n=== SIMULACIÓN TRANSFORMER & VRAM (12GB ARC B580) ===")
	simulateTransformerLimits()
}

func runBenchmark(name, dir string) {
	fmt.Printf("[%s] Ejecutando...\n", name)

	// Convertimos path a absoluto para evitar líos con el cd
	absDir, _ := filepath.Abs(dir)

	cmd := exec.Command("bash", "-c", "cd "+absDir+" && source /opt/intel/oneapi/setvars.sh && ./matrix_go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("  Error corriendo %s: %v\n", name, err)
		return
	}

	// Extraer solo las líneas de tiempo para limpiar la salida
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "GPU Time") || strings.Contains(line, "Speedup") || strings.Contains(line, "Device") {
			fmt.Println("  " + strings.TrimSpace(line))
		}
	}
}

func simulateTransformerLimits() {
	const (
		VRAMTotalGB = 12.0
		OSReserved  = 1.0 // Reserva típica del sistema
		VRAMUsable  = VRAMTotalGB - OSReserved
	)

	fmt.Printf("VRAM Usable estimada: %.1f GB\n", VRAMUsable)

	// Tamaños típicos (Hidden Dimension)
	hSizes := []int{4096, 8192, 12288}

	for _, h := range hSizes {
		fmt.Printf("\n--- Configuración Hidden Dimension: %d ---\n", h)

		// 1. Pesos de una capa Transformer (Estimación simplificada)
		// Q, K, V, O (4 * h^2) + MLP (8 * h^2 en Llama/SwiGLU) = 12 * h^2 per layer
		layerWeightsFP16 := float64(12*h*h*2) / (1024 * 1024 * 1024)
		layerWeightsINT4 := float64(12*h*h) * 0.5 / (1024 * 1024 * 1024)

		fmt.Printf("  Memoria pesos por capa (FP16): %.3f GB\n", layerWeightsFP16)
		fmt.Printf("  Memoria pesos por capa (INT4): %.3f GB\n", layerWeightsINT4)

		// 2. ¿Cuántas capas caben en 12GB? (Dedicando todo el espacio a pesos)
		maxLayersINT4 := int(VRAMUsable / layerWeightsINT4)
		fmt.Printf("  Max capas INT4 posibles (sin KV Cache): %d\n", maxLayersINT4)

		// 3. Impacto del KV Cache
		// Memoria = 2 * n_layers * n_heads * d_head * seq_len * precision
		// Simplificado: 2 * n_layers * hidden_size * seq_len * 2 (FP16)
		seqLen := 4096.0
		layersRealistic := 32 // Un modelo típico (ej. Llama 7B)
		kvCacheSize := (2.0 * float64(layersRealistic) * float64(h) * seqLen * 2.0) / (1024 * 1024 * 1024)

		fmt.Printf("  KV Cache (32 capas, context %.0f, FP16): %.2f GB\n", seqLen, kvCacheSize)

		totalMem7B_INT4 := (float64(layersRealistic) * layerWeightsINT4) + kvCacheSize
		fmt.Printf("  Total Estimado para 7B (32 capas) en INT4: %.2f GB\n", totalMem7B_INT4)
		if totalMem7B_INT4 < VRAMUsable {
			fmt.Println("  [CHECK] ¡Este modelo CABE en una sola GPU!")
		} else {
			fmt.Println("  [WARNING] Excede la VRAM de una sola GPU.")
		}
	}
}
