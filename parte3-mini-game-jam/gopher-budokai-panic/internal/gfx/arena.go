package gfx

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Arena struct {
	cloudOffset float64
	debrisPhase float64
}

func NewArena() *Arena {
	return &Arena{}
}

func (a *Arena) Update(dt float64) {
	a.cloudOffset += 12.0 * dt
	if a.cloudOffset >= 640.0 {
		a.cloudOffset -= 640.0
	}
	a.debrisPhase += 3.0 * dt
}

func (a *Arena) Draw(screen *ebiten.Image, width, height float64) {
	w := float32(width)
	h := float32(height)

	// 1. Sky Gradient (Warm DBZ Anime Blue to Cyan Horizon)
	skySteps := 12
	stepHeight := h * 0.70 / float32(skySteps)
	for i := 0; i < skySteps; i++ {
		t := float32(i) / float32(skySteps)
		r := uint8(20 + t*100)
		g := uint8(80 + t*120)
		b := uint8(190 + t*50)
		y := float32(i) * stepHeight
		vector.DrawFilledRect(screen, 0, y, w, stepHeight+1, color.RGBA{R: r, G: g, B: b, A: 255}, false)
	}

	// 2. Parallax Moving Anime Clouds
	cloudCol := color.RGBA{R: 255, G: 255, B: 255, A: 160}
	for i := 0; i < 4; i++ {
		cx := float32(math.Mod(float64(i)*200.0-a.cloudOffset, float64(w+120.0))) - 60.0
		cy := float32(40 + (i%2)*30)
		vector.DrawFilledCircle(screen, cx, cy, 26, cloudCol, true)
		vector.DrawFilledCircle(screen, cx+25, cy-8, 34, cloudCol, true)
		vector.DrawFilledCircle(screen, cx+55, cy+2, 28, cloudCol, true)
		vector.DrawFilledRect(screen, cx, cy, 60, 20, cloudCol, false)
	}

	// 3. Distant Rocky Mountains (Classic DBZ Wasteland Pinnacles)
	mountainCol := color.RGBA{R: 160, G: 110, B: 65, A: 255}
	mountainDark := color.RGBA{R: 120, G: 75, B: 40, A: 255}
	pinnacles := []struct{ x, w, topY float32 }{
		{40, 70, 160},
		{150, 90, 140},
		{280, 110, 150},
		{430, 80, 135},
		{540, 100, 155},
	}
	for _, p := range pinnacles {
		// Left slope
		vector.DrawFilledRect(screen, p.x, p.topY, p.w*0.5, h*0.75-p.topY, mountainCol, false)
		// Right shade slope
		vector.DrawFilledRect(screen, p.x+p.w*0.5, p.topY, p.w*0.5, h*0.75-p.topY, mountainDark, false)
	}

	// 4. Ground Arena Floor (Martial Arts Platform & Rocky Terrain)
	groundY := h * 0.75
	groundH := h - groundY

	// Ground base
	vector.DrawFilledRect(screen, 0, groundY, w, groundH, color.RGBA{R: 195, G: 160, B: 110, A: 255}, false)
	// Arena Ring Tile Inset
	ringMargin := float32(50.0)
	vector.DrawFilledRect(screen, ringMargin, groundY+6, w-ringMargin*2, groundH-12, color.RGBA{R: 215, G: 190, B: 145, A: 255}, false)

	// Tile grid lines
	gridCol := color.RGBA{R: 155, G: 120, B: 80, A: 255}
	for tx := ringMargin; tx <= w-ringMargin; tx += 45 {
		vector.StrokeLine(screen, tx, groundY+6, tx, h-6, 1.5, gridCol, false)
	}
	vector.StrokeLine(screen, ringMargin, groundY+25, w-ringMargin, groundY+25, 1.5, gridCol, false)
	vector.StrokeLine(screen, ringMargin, groundY+50, w-ringMargin, groundY+50, 1.5, gridCol, false)

	// Floating energy rocks / debris when battle is intense
	for i := 0; i < 6; i++ {
		rx := float32(80 + i*90)
		ry := float32(groundY-15) + float32(math.Sin(a.debrisPhase+float64(i)*1.2))*6.0
		vector.DrawFilledRect(screen, rx, ry, 6, 5, mountainDark, false)
	}
}

// DrawPanicVignette renders the red flashing alarm vignette when the player is in PANIC!
func DrawPanicVignette(screen *ebiten.Image, width, height float64, panicProgress float64) {
	if panicProgress <= 0 {
		return
	}
	w := float32(width)
	h := float32(height)

	alpha := float32(panicProgress) * 0.45
	redCol := color.RGBA{R: 255, G: 20, B: 20, A: uint8(alpha * 255)}

	// Red border frame
	border := float32(14.0 * panicProgress)
	vector.DrawFilledRect(screen, 0, 0, w, border, redCol, false)
	vector.DrawFilledRect(screen, 0, h-border, w, border, redCol, false)
	vector.DrawFilledRect(screen, 0, 0, border, h, redCol, false)
	vector.DrawFilledRect(screen, w-border, 0, border, h, redCol, false)

	// Corner pulse wedges
	cornerSize := float32(60.0 * panicProgress)
	vector.DrawFilledRect(screen, 0, 0, cornerSize, cornerSize*0.5, redCol, false)
	vector.DrawFilledRect(screen, w-cornerSize, 0, cornerSize, cornerSize*0.5, redCol, false)
	vector.DrawFilledRect(screen, 0, h-cornerSize*0.5, cornerSize, cornerSize*0.5, redCol, false)
	vector.DrawFilledRect(screen, w-cornerSize, h-cornerSize*0.5, cornerSize, cornerSize*0.5, redCol, false)
}
