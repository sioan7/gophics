package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth    = 800
	screenHeight   = 800
	maxTime        = 2.
	initialTime    = -2.
	timeIncrement  = .0007
	numberOfPoints = 100000
	scaleFactor    = 300
)

type point struct {
	x     float64
	y     float64
	color color.Color
}

func nextColor(prev color.Color, ratioOfChannel float64) color.Color {
	r := uint8(rand.Intn(128) + 128)
	g := uint8(rand.Intn(128) + 128)
	b := uint8(rand.Intn(128) + 128)
	return color.RGBA{r, g, b, 255}
	// return color.RGBA{255, 215, 0, 255}
	// return color.RGBA{0, 128, 255, 255}
	// return color.RGBA{255, 255, uint8((1 - ratioOfChannel) * 255), 255}
}

func calcPointCoordinates(p point, t float64) (float64, float64) {
	// x := -p.y*p.y - p.x*t + p.y
	// y := p.x*p.x - p.x*p.y + t

	// x := math.Sin(-p.y*p.y - p.x*t + p.y)
	// y := math.Cos(p.x*p.x - p.x*p.y + t)

	// x := math.Sin(-p.x+p.x*t) * t
	// y := math.Cos(p.y*p.y*t) * t

	x := math.Sin(-p.y*p.y-p.x*t+p.y*t+t) * t
	y := math.Cos(p.x*p.x-p.x*p.y+t) * t

	// x := math.Sin(-p.y*p.y - p.x*t + p.y*t + t)
	// y := math.Cos(p.x*p.x - p.x*p.y + t)
	return x, y
}

func updatePoints(points []point, t float64) {
	for i := 1; i < len(points); i++ {
		x, y := calcPointCoordinates(points[i-1], t)
		points[i].x = x
		points[i].y = y
	}
}

func initPoints(points []point, initialPoint point) {
	points[0] = initialPoint
	for i := 1; i < len(points); i++ {
		points[i] = point{0, 0, nextColor(points[0].color, float64(i)/numberOfPoints)}
	}
}

func nextInstant(t float64) float64 {
	return t + timeIncrement
}

func toScreenCoordinate(a float64, dimension int) int {
	b := int(a*scaleFactor) + dimension/2
	if b < 0 {
		b = 0
	}
	if b >= dimension {
		b = dimension - 1
	}
	return b
}

func toScreenCoordinates(p point) (int, int) {
	return toScreenCoordinate(p.x, screenWidth), toScreenCoordinate(p.y, screenHeight)
}

var (
	t            = initialTime
	initialPoint = point{t, t, color.White}
	points       = make([]point, numberOfPoints)
	secondPoints = make([]point, numberOfPoints)
)

func update(screen *ebiten.Image) error {

	updatePoints(points, t)
	t = nextInstant(t)
	fmt.Println(t)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.Black)

	for _, p := range points {
		x, y := toScreenCoordinates(p)
		screen.Set(x, y, p.color)
	}

	return nil
}

func main() {
	initPoints(points, initialPoint)
	initPoints(secondPoints, initialPoint)
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "gophics"); err != nil {
		log.Fatal(err)
	}
}
