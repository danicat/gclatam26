package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// ToolType identifies the active physical tool carried by the chef.
type ToolType int

const (
	ToolNone ToolType = iota
	ToolExtinguisher
	ToolIce
	ToolWrench
	ToolCatToy
)

func (t ToolType) Name() string {
	switch t {
	case ToolExtinguisher:
		return "EXTINGUISHER"
	case ToolIce:
		return "ICE BUCKET"
	case ToolWrench:
		return "WRENCH"
	case ToolCatToy:
		return "CAT TOY"
	default:
		return "EMPTY HANDS"
	}
}

// ToolRack is a physical station where the chef can pick up or swap tools.
type ToolRack struct {
	X, Y     float64
	W, H     float64
	Tool     ToolType
	drawOpts ebiten.DrawImageOptions
}

func NewToolRack(x, y float64, tool ToolType) *ToolRack {
	return &ToolRack{
		X:    x,
		Y:    y,
		W:    18,
		H:    18,
		Tool: tool,
	}
}

// Draw renders the tool rack shelf and its hosted tool icon.
func (tr *ToolRack) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	// 1. Draw shelf table/frame
	tr.drawRect(screen, pixelImg, tr.X, tr.Y, tr.W, tr.H, color.RGBA{45, 40, 60, 255})
	tr.drawRect(screen, pixelImg, tr.X+1, tr.Y+1, tr.W-2, tr.H-2, color.RGBA{85, 75, 105, 255})
	tr.drawRect(screen, pixelImg, tr.X+2, tr.Y+2, tr.W-4, tr.H-4, color.RGBA{115, 105, 135, 255})

	// 2. Draw tool icon centered on rack
	DrawToolIcon(screen, pixelImg, tr.Tool, tr.X+tr.W/2, tr.Y+tr.H/2)
}

func (tr *ToolRack) drawRect(screen *ebiten.Image, pixelImg *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	tr.drawOpts.GeoM.Reset()
	tr.drawOpts.GeoM.Scale(w, h)
	tr.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	tr.drawOpts.ColorScale.Reset()
	tr.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(pixelImg, &tr.drawOpts)
}

// DrawToolIcon renders a crisp 8x8 procedural tool graphic at (cx, cy).
func DrawToolIcon(screen *ebiten.Image, pixelImg *ebiten.Image, tool ToolType, cx, cy float64) {
	var opts ebiten.DrawImageOptions
	drawPix := func(ox, oy, w, h float64, c color.RGBA) {
		opts.GeoM.Reset()
		opts.GeoM.Scale(w, h)
		opts.GeoM.Translate(cx+ox, cy+oy)
		af := float32(c.A) / 255.0
		opts.ColorScale.Reset()
		opts.ColorScale.Scale((float32(c.R)/255.0)*af, (float32(c.G)/255.0)*af, (float32(c.B)/255.0)*af, af)
		screen.DrawImage(pixelImg, &opts)
	}

	switch tool {
	case ToolExtinguisher:
		// Red tank body
		drawPix(-3, -4, 6, 8, color.RGBA{220, 40, 35, 255})
		// Tank highlight
		drawPix(-1, -3, 2, 6, color.RGBA{255, 110, 95, 255})
		// Black nozzle / handle
		drawPix(-2, -6, 4, 2, color.RGBA{30, 25, 40, 255})
		drawPix(1, -5, 3, 2, color.RGBA{60, 55, 75, 255})

	case ToolIce:
		// Blue bucket
		drawPix(-4, -2, 8, 6, color.RGBA{45, 110, 210, 255})
		// Ice cubes poking out
		drawPix(-3, -5, 3, 3, color.RGBA{180, 235, 255, 255})
		drawPix(1, -4, 3, 3, color.RGBA{220, 245, 255, 255})
		// Bucket rim
		drawPix(-5, -2, 10, 1.5, color.RGBA{110, 175, 255, 255})

	case ToolWrench:
		// Steel wrench diagonal
		drawPix(-4, -4, 3, 3, color.RGBA{200, 210, 225, 255})
		drawPix(-2, -2, 4, 3, color.RGBA{140, 150, 170, 255})
		drawPix(1, 0, 3, 3, color.RGBA{100, 110, 130, 255})
		drawPix(2, 2, 3, 3, color.RGBA{200, 210, 225, 255})

	case ToolCatToy:
		// Golden bell / yarn ball
		drawPix(-3, -3, 6, 6, color.RGBA{255, 205, 50, 255})
		drawPix(-1, -1, 2, 2, color.RGBA{255, 245, 140, 255})
		// Ribbon
		drawPix(-2, 3, 4, 2, color.RGBA{255, 70, 120, 255})
	}
}
