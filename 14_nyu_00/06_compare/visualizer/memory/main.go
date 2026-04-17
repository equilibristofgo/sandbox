package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	maxVRAM      = 12.0 // GB
	hiddenDim    = 4096
	numLayers    = 32
)

type Game struct {
	seqLen     int
	simStarted bool
	lastUpdate time.Time
}

func (g *Game) Update() error {
	if !g.simStarted {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.simStarted = true
			g.lastUpdate = time.Now()
		}
		return nil
	}

	totalMem := 3.5 + float64(g.seqLen*hiddenDim*numLayers*2*2)/(1024*1024*1024)
	if totalMem > maxVRAM {
		return nil // Stop at OOM
	}

	if time.Since(g.lastUpdate) > 50*time.Millisecond {
		g.seqLen += 32
		g.lastUpdate = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "VRAM & CONTEXT MONITOR - GPU ARC B580 (12GB)\nPress SPACE to start...")

	// Tanque de memoria (Contenedor)
	tankX, tankY := float32(200), float32(100)
	tankW, tankH := float32(200), float32(300)
	vector.StrokeRect(screen, tankX, tankY, tankW, tankH, 2, color.White, true)

	// Pesos (3.5GB) - Fijo
	weightMem := 3.5
	weightH := float32((weightMem / maxVRAM) * float64(tankH))
	vector.DrawFilledRect(screen, tankX+2, tankY+tankH-weightH-2, tankW-4, weightH, color.RGBA{100, 100, 255, 255}, true)
	ebitenutil.DebugPrintAt(screen, "WEIGHTS (7B INT4)", int(tankX)+210, int(tankY+tankH-weightH/2))

	// KV Cache (Dinámico)
	kvMem := float64(g.seqLen*hiddenDim*numLayers*2*2) / (1024 * 1024 * 1024)
	kvH := float32((kvMem / maxVRAM) * float64(tankH))

	kvColor := color.RGBA{255, 165, 0, 200}
	totalMem := weightMem + kvMem
	if totalMem > maxVRAM {
		kvColor = color.RGBA{255, 0, 0, 255}
		ebitenutil.DebugPrintAt(screen, "!!! OUT OF MEMORY !!!", int(tankX), int(tankY)-30)
	}

	vector.DrawFilledRect(screen, tankX+2, tankY+tankH-weightH-kvH-2, tankW-4, kvH, kvColor, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("KV CACHE (%d tokens)", g.seqLen), int(tankX)+210, int(tankY+tankH-weightH-kvH/2))

	// Stats
	stats := fmt.Sprintf("Tokens: %d\nPesos: %.2f GB\nKV Cache: %.2f GB\nTotal: %.2f / 12.0 GB", g.seqLen, weightMem, kvMem, totalMem)
	ebitenutil.DebugPrintAt(screen, stats, 50, 100)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("VRAM Consumption Simulator (Ebiten)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
