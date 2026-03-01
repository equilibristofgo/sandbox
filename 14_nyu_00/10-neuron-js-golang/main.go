package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Dataset XOR: El reto clásico de aprendizaje no lineal
var trainingData = []struct {
	inputs  []float64
	targets []float64
}{
	{[]float64{0, 0}, []float64{0}},
	{[]float64{0, 1}, []float64{1}},
	{[]float64{1, 0}, []float64{1}},
	{[]float64{1, 1}, []float64{0}},
}

// --- RED NEURONAL ---

type NeuralNetwork struct {
	inputLayer  []float64
	hiddenLayer []float64
	outputLayer []float64
	weightsIH   [][]float64 // Pesos Input -> Hidden
	weightsHO   [][]float64 // Pesos Hidden -> Output
	biasH       []float64
	biasO       []float64
}

func NewNeuralNetwork(in, hid, out int) *NeuralNetwork {
	nn := &NeuralNetwork{
		inputLayer:  make([]float64, in),
		hiddenLayer: make([]float64, hid),
		outputLayer: make([]float64, out),
		weightsIH:   make([][]float64, hid),
		weightsHO:   make([][]float64, out),
		biasH:       make([]float64, hid),
		biasO:       make([]float64, out),
	}
	nn.ResetWeights()
	return nn
}

func (nn *NeuralNetwork) ResetWeights() {
	for i := range nn.weightsIH {
		nn.weightsIH[i] = make([]float64, len(nn.inputLayer))
		for j := range nn.weightsIH[i] {
			nn.weightsIH[i][j] = rand.NormFloat64() * 0.5
		}
		nn.biasH[i] = rand.NormFloat64() * 0.5
	}
	for i := range nn.weightsHO {
		nn.weightsHO[i] = make([]float64, len(nn.hiddenLayer))
		for j := range nn.weightsHO[i] {
			nn.weightsHO[i][j] = rand.NormFloat64() * 0.5
		}
		nn.biasO[i] = rand.NormFloat64() * 0.5
	}
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func (nn *NeuralNetwork) Forward(inputs []float64) {
	copy(nn.inputLayer, inputs)
	// Input a Hidden (usamos Tanh para variedad interna)
	for i := 0; i < len(nn.hiddenLayer); i++ {
		sum := nn.biasH[i]
		for j := 0; j < len(nn.inputLayer); j++ {
			sum += nn.inputLayer[j] * nn.weightsIH[i][j]
		}
		nn.hiddenLayer[i] = math.Tanh(sum)
	}
	// Hidden a Output (usamos Sigmoide como en tu gráfica)
	for i := 0; i < len(nn.outputLayer); i++ {
		sum := nn.biasO[i]
		for j := 0; j < len(nn.hiddenLayer); j++ {
			sum += nn.hiddenLayer[j] * nn.weightsHO[i][j]
		}
		nn.outputLayer[i] = sigmoid(sum)
	}
}

func (nn *NeuralNetwork) Backpropagate(targets []float64, lr float64) float64 {
	outErr := make([]float64, len(nn.outputLayer))
	totalErr := 0.0
	for i := range outErr {
		outErr[i] = targets[i] - nn.outputLayer[i]
		totalErr += outErr[i] * outErr[i]
	}

	// Gradiente Output
	for i := range nn.weightsHO {
		gradient := outErr[i] * nn.outputLayer[i] * (1 - nn.outputLayer[i])
		for j := range nn.weightsHO[i] {
			nn.weightsHO[i][j] += lr * gradient * nn.hiddenLayer[j]
		}
		nn.biasO[i] += lr * gradient
	}

	// Error Hidden
	hidErr := make([]float64, len(nn.hiddenLayer))
	for j := range nn.hiddenLayer {
		sum := 0.0
		for i := range nn.outputLayer {
			sum += outErr[i] * nn.weightsHO[i][j]
		}
		hidErr[j] = sum
	}

	// Gradiente Hidden
	for i := range nn.weightsIH {
		gradient := hidErr[i] * (1 - nn.hiddenLayer[i]*nn.hiddenLayer[i])
		for j := range nn.weightsIH[i] {
			nn.weightsIH[i][j] += lr * gradient * nn.inputLayer[j]
		}
		nn.biasH[i] += lr * gradient
	}
	return totalErr
}

// --- JUEGO ---

type Game struct {
	nn         *NeuralNetwork
	lr         float64
	errHistory []float64
	currentIdx int
}

func (g *Game) Update() error {
	// Controles
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.nn.ResetWeights()
		g.errHistory = nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.lr += 0.002
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.lr = math.Max(0, g.lr-0.002)
	}

	// Entrenamiento
	var multiplier float64 = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		multiplier = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		multiplier = -1
	}

	if multiplier != 0 {
		for i := 0; i < 15; i++ {
			sample := trainingData[rand.Intn(len(trainingData))]
			g.nn.Forward(sample.inputs)
			err := g.nn.Backpropagate(sample.targets, g.lr*multiplier)
			if len(g.errHistory) > 250 {
				g.errHistory = g.errHistory[1:]
			}
			g.errHistory = append(g.errHistory, err)
		}
	}

	// Rotar visualización
	if time.Now().UnixNano()%30 == 0 {
		g.currentIdx = (g.currentIdx + 1) % len(trainingData)
	}
	g.nn.Forward(trainingData[g.currentIdx].inputs)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	inX, hidX, outX := 120.0, 320.0, 520.0

	// 1. Conexiones y Pesos Numéricos
	for i, weights := range g.nn.weightsIH {
		for j, w := range weights {
			drawEdge(screen, float32(inX), float32(100+j*150), float32(hidX), float32(60+i*90), w)
		}
	}
	for i, weights := range g.nn.weightsHO {
		for j, w := range weights {
			drawEdge(screen, float32(hidX), float32(60+j*90), float32(outX), float32(180+i*100), w)
		}
	}

	// 2. Nodos
	drawLayer(screen, inX, 100, 150, g.nn.inputLayer, color.RGBA{100, 100, 255, 255})
	drawLayer(screen, hidX, 60, 90, g.nn.hiddenLayer, color.RGBA{100, 255, 100, 255})
	drawLayer(screen, outX, 180, 100, g.nn.outputLayer, color.RGBA{255, 100, 100, 255})

	// 3. Gráfica de Error (Línea amarilla)
	drawGraph(screen, g.errHistory)

	// GUI
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"CONTROLES:\nArriba: Entrenar\nAbajo: Des-entrenar\nR: Reset Pesos\n\nLearning Rate: %.3f\nTarget: %.0f\nPred: %.4f",
		g.lr, trainingData[g.currentIdx].targets[0], g.nn.outputLayer[0]))
}

func drawEdge(screen *ebiten.Image, x1, y1, x2, y2 float32, w float64) {
	c := color.RGBA{200, 50, 50, 100}
	if w > 0 {
		c = color.RGBA{50, 200, 50, 100}
	}
	vector.StrokeLine(screen, x1, y1, x2, y2, float32(math.Abs(w)*2)+0.5, c, true)
	// Dibujar valor del peso
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f", w), int((x1+x2)/2)-10, int((y1+y2)/2)-5)
}

func drawLayer(screen *ebiten.Image, x, startY, space float64, nodes []float64, c color.RGBA) {
	for i, val := range nodes {
		py := float32(startY + float64(i)*space)
		alpha := uint8(math.Min(255, math.Max(50, math.Abs(val)*255)))
		vector.DrawFilledCircle(screen, float32(x), py, 18, color.RGBA{c.R, c.G, c.B, alpha}, true)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", val), int(x)-15, int(py)+20)
	}
}

func drawGraph(screen *ebiten.Image, hist []float64) {
	if len(hist) < 2 {
		return
	}
	var bx, by float32 = 350, 450
	vector.StrokeLine(screen, bx, by, bx+250, by, 1, color.White, true)
	for i := 0; i < len(hist)-1; i++ {
		vector.StrokeLine(screen, bx+float32(i), by-float32(hist[i]*100), bx+float32(i+1), by-float32(hist[i+1]*100), 2, color.RGBA{255, 255, 0, 255}, true)
	}
	ebitenutil.DebugPrintAt(screen, "HISTORIAL ERROR", int(bx), int(by)+5)
}

func (g *Game) Layout(w, h int) (int, int) { return 640, 480 }

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{nn: NewNeuralNetwork(2, 4, 1), lr: 0.1}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Simulador XOR: Aprendizaje y Olvido")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
