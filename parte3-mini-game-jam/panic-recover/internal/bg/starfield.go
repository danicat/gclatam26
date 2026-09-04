package bg

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Star struct {
	X, Y       float64
	Speed      float64
	Brightness uint8
	Size       float64
}

type CodeSnippet struct {
	Text  string
	X, Y  float64
	Speed float64
	Alpha float64
}

type Starfield struct {
	width, height int
	stars         []Star
	codeSnippets  []CodeSnippet
	gridOffset    float64
	gridSpeed     float64
	dotTexture    *ebiten.Image
	updateCount   int64
	speedMult     float64
}

func (sf *Starfield) UpdateCount() int64 {
	return sf.updateCount
}

func (sf *Starfield) SetSpeedMultiplier(mult float64) {
	if mult < 0.5 {
		mult = 0.5
	}
	sf.speedMult = mult
}

var sampleSnippets = []string{
	"goroutine 42 [running]:",
	"0x00c00008e010",
	"runtime.panic()",
	"defer recover()",
	"sync.Mutex.Lock()",
	"select { case <-ch: }",
	"make(chan error, 1)",
	"type assertion (*Gopher)",
	"0x7ffeefbff568",
	"go worker(runtime.CPU)",
}

func NewStarfield(width, height int) *Starfield {
	numStars := 120
	sf := &Starfield{
		width:      width,
		height:     height,
		stars:      make([]Star, numStars),
		gridSpeed:  45.0,
		dotTexture: ebiten.NewImage(2, 2),
	}
	sf.dotTexture.Fill(color.White)

	for i := 0; i < numStars; i++ {
		sf.stars[i] = Star{
			X:          rand.Float64() * float64(width),
			Y:          rand.Float64() * float64(height),
			Speed:      25.0 + rand.Float64()*110.0,
			Brightness: uint8(100 + rand.Intn(155)),
			Size:       1.0 + rand.Float64()*1.5,
		}
	}

	// Initialize 6 drifting code snippets
	for i := 0; i < 6; i++ {
		sf.codeSnippets = append(sf.codeSnippets, CodeSnippet{
			Text:  sampleSnippets[rand.Intn(len(sampleSnippets))],
			X:     20.0 + rand.Float64()*float64(width-180),
			Y:     rand.Float64() * float64(height),
			Speed: 20.0 + rand.Float64()*25.0,
			Alpha: 0.25 + rand.Float64()*0.25,
		})
	}

	sf.speedMult = 1.0
	return sf
}

func (sf *Starfield) Update(dt float64) {
	sf.updateCount++
	effectiveDt := dt * sf.speedMult
	// Scroll grid
	sf.gridOffset += sf.gridSpeed * effectiveDt
	if sf.gridOffset >= 40.0 {
		sf.gridOffset -= 40.0
	}

	// Update stars
	for i := range sf.stars {
		s := &sf.stars[i]
		s.Y += s.Speed * effectiveDt
		if s.Y > float64(sf.height) {
			s.Y = 0
			s.X = rand.Float64() * float64(sf.width)
		}
	}

	// Update code snippets
	for i := range sf.codeSnippets {
		cs := &sf.codeSnippets[i]
		cs.Y += cs.Speed * effectiveDt
		if cs.Y > float64(sf.height)+20 {
			cs.Y = -20
			cs.X = 20.0 + rand.Float64()*float64(sf.width-180)
			cs.Text = sampleSnippets[rand.Intn(len(sampleSnippets))]
		}
	}
}

func (sf *Starfield) Draw(screen *ebiten.Image, monoSource *text.GoTextFaceSource) {
	// 1. Draw dark cyber-grid lines
	gridColor := color.RGBA{15, 25, 45, 180}
	gridStep := 40.0

	// Vertical grid lines
	for x := 0.0; x < float64(sf.width); x += gridStep {
		var op ebiten.DrawImageOptions
		op.GeoM.Scale(1.0, float64(sf.height))
		op.GeoM.Translate(x, 0)
		op.ColorScale.ScaleWithColor(gridColor)
		screen.DrawImage(sf.dotTexture, &op)
	}

	// Horizontal scrolling grid lines
	for y := sf.gridOffset - 40.0; y < float64(sf.height); y += gridStep {
		if y < 0 {
			continue
		}
		var op ebiten.DrawImageOptions
		op.GeoM.Scale(float64(sf.width), 1.0)
		op.GeoM.Translate(0, y)
		op.ColorScale.ScaleWithColor(gridColor)
		screen.DrawImage(sf.dotTexture, &op)
	}

	// 2. Draw stars
	for _, s := range sf.stars {
		var op ebiten.DrawImageOptions
		op.GeoM.Scale(s.Size*0.5, s.Size*0.5)
		op.GeoM.Translate(s.X, s.Y)
		c := color.RGBA{s.Brightness / 2, s.Brightness, s.Brightness, 255}
		op.ColorScale.ScaleWithColor(c)
		screen.DrawImage(sf.dotTexture, &op)
	}

	// 3. Draw drifting code snippets if font is provided
	if monoSource != nil {
		for _, cs := range sf.codeSnippets {
			drawTextSimple(screen, cs.Text, cs.X, cs.Y, color.RGBA{40, 140, 200, uint8(cs.Alpha * 255)}, monoSource, 11)
		}
	}
}

func drawTextSimple(dst *ebiten.Image, msg string, x, y float64, clr color.Color, src *text.GoTextFaceSource, size float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	f := &text.GoTextFace{
		Source: src,
		Size:   size,
	}
	text.Draw(dst, msg, f, op)
}
