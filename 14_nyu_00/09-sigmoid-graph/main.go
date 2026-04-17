package main

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func main() {
	// Crear un gráfico
	p := plot.New()
	p.Title.Text = "Sigmoid Function"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "sigmoid(x)"

	// Generar datos
	z := make([]float64, 100)
	for i := 0; i < 100; i++ {
		z[i] = float64(i-50) / 5.0 // Rango de -10 a 10
	}

	// Crear puntos para graficar
	points := make(plotter.XYs, len(z))
	for i := range z {
		points[i].X = z[i]
		points[i].Y = sigmoid(z[i])
	}

	// Agregar puntos al gráfico
	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	// Mostrar el gráfico
	if err := p.Save(400, 400, "sigmoid.png"); err != nil {
		panic(err)
	}
}
