package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 850
	screenHeight = 650
	matrixSize   = 12
	cellSize     = 18
	offsetX      = 50
	offsetY      = 120
)

type GPUThread struct {
	row, col int
	active   bool
	progress float32
}

type Stage struct {
	Name    string
	MatrixA string
	MatrixB string
	MatrixC string
	Status  string
}

type Game struct {
	matrixA      [matrixSize][matrixSize]float32
	matrixB      [matrixSize][matrixSize]float32
	matrixC      [matrixSize][matrixSize]float32
	threads      []GPUThread
	simStarted   bool
	lastUpdate   time.Time
	currentStage int
	stages       []Stage
	stageTimer   time.Time
	isTransition bool
}

func (g *Game) Update() error {
	if !g.simStarted {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.simStarted = true
			g.currentStage = 0
			g.startStage()
		}
		return nil
	}

	if g.isTransition {
		if time.Since(g.stageTimer) > 2*time.Second {
			g.isTransition = false
			g.currentStage++
			if g.currentStage >= len(g.stages) {
				g.simStarted = false // End of sim
				return nil
			}
			g.startStage()
		}
		return nil
	}

	// Check if current stage is done (all matrix C filled)
	done := true
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			if g.matrixC[i][j] == 0 {
				done = false
				break
			}
		}
	}

	if done && !g.isTransition {
		g.isTransition = true
		g.stageTimer = time.Now()
		return nil
	}

	// Simular paralelismo
	if time.Since(g.lastUpdate) > 80*time.Millisecond {
		for i := 0; i < 3; i++ {
			r := rand.Intn(matrixSize)
			c := rand.Intn(matrixSize)
			if g.matrixC[r][c] == 0 {
				found := false
				for _, t := range g.threads {
					if t.row == r && t.col == c && t.active {
						found = true
						break
					}
				}
				if !found {
					g.threads = append(g.threads, GPUThread{row: r, col: c, active: true, progress: 0})
				}
			}
		}
		g.lastUpdate = time.Now()
	}

	for i := range g.threads {
		if g.threads[i].active {
			g.threads[i].progress += 0.08
			if g.threads[i].progress >= 1.0 {
				g.threads[i].active = false
				g.matrixC[g.threads[i].row][g.threads[i].col] = 1.0
			}
		}
	}

	return nil
}

func (g *Game) startStage() {
	g.threads = nil
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			g.matrixC[i][j] = 0
		}
	}
	g.lastUpdate = time.Now()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.simStarted {
		ebitenutil.DebugPrint(screen, "TRANSFORMER GPU SIMULATOR\n\nPresiona ESPACIO para iniciar la simulación completa\n\nVerás las etapas de Proyección, Atención y MLP.")
		return
	}

	stage := g.stages[g.currentStage]
	header := fmt.Sprintf("ETAPA %d: %s\n%s\n", g.currentStage+1, stage.Name, stage.Status)
	if g.isTransition {
		header += "\n--- PREPARANDO SIGUIENTE ETAPA ---"
	}
	ebitenutil.DebugPrint(screen, header)

	// Dibujar Matrices con etiquetas dinámicas
	drawMatrix(screen, stage.MatrixA, offsetX, offsetY, g.matrixA, color.RGBA{100, 100, 255, 255})

	// Si es etapa de atención, resaltar "Cache Hit"
	bColor := color.RGBA{100, 255, 100, 255}
	if g.currentStage == 1 { // Atención
		bColor = color.RGBA{0, 200, 0, 255}
		ebitenutil.DebugPrintAt(screen, "KV CACHE ACTIVE", offsetX+(matrixSize+2)*cellSize, offsetY-40)
	}
	drawMatrix(screen, stage.MatrixB, offsetX+(matrixSize+2)*cellSize, offsetY, g.matrixB, bColor)

	drawMatrix(screen, stage.MatrixC, offsetX+(2*matrixSize+4)*cellSize, offsetY, g.matrixC, color.RGBA{255, 255, 255, 255})

	// Dibujar Hilos GPU
	for _, t := range g.threads {
		if t.active {
			vector.DrawFilledRect(screen, float32(offsetX), float32(offsetY+t.row*cellSize), float32(matrixSize*cellSize), float32(cellSize), color.RGBA{255, 165, 0, 80}, true)
			vector.DrawFilledRect(screen, float32(offsetX+(matrixSize+2)*cellSize+t.col*cellSize), float32(offsetY), float32(cellSize), float32(matrixSize*cellSize), color.RGBA{255, 165, 0, 80}, true)

			targetX := float32(offsetX + (2*matrixSize+4)*cellSize + t.col*cellSize)
			targetY := float32(offsetY + t.row*cellSize)
			vector.DrawFilledRect(screen, targetX, targetY, float32(cellSize), float32(cellSize), color.RGBA{255, 165, 0, 255}, true)
		}
	}
}

func drawMatrix(screen *ebiten.Image, title string, x, y int, m [matrixSize][matrixSize]float32, c color.Color) {
	ebitenutil.DebugPrintAt(screen, title, x, y-20)
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			rectColor := color.RGBA{40, 40, 40, 255}
			if m[i][j] > 0 {
				rectColor = c.(color.RGBA)
			}
			vector.StrokeRect(screen, float32(x+j*cellSize), float32(y+i*cellSize), float32(cellSize), float32(cellSize), 1, color.RGBA{80, 80, 80, 255}, true)
			vector.DrawFilledRect(screen, float32(x+j*cellSize+1), float32(y+i*cellSize+1), float32(cellSize-2), float32(cellSize-2), rectColor, true)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		stages: []Stage{
			{
				Name:    "Proyecciones Q, K, V",
				MatrixA: "Tokens Input",
				MatrixB: "Weights Proj",
				MatrixC: "QKV Vectors",
				Status:  "Multiplicando tokens por pesos de proyección.",
			},
			{
				Name:    "Self-Attention (Q @ K^T)",
				MatrixA: "Query",
				MatrixB: "Key (KV Cache)",
				MatrixC: "Attention Scores",
				Status:  "¡Cache Hit! Usando Keys de tokens anteriores.",
			},
			{
				Name:    "Feed-Forward (MLP)",
				MatrixA: "Attn Output",
				MatrixB: "Weights MLP",
				MatrixC: "Next States",
				Status:  "Procesando la salida de atención.",
			},
		},
	}

	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			g.matrixA[i][j] = rand.Float32()
			g.matrixB[i][j] = rand.Float32()
		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Transformer GPU Simulator (Ebiten)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
