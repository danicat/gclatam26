package ui

import (
	"fmt"
	"image/color"
	"math"

	"exploding-kitchens/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
)

// PopupMessage displays a transient floating notification (e.g. "CLUTCH! +500").
type PopupMessage struct {
	Text      string
	X, Y      float64
	Age, Life float64
	Color     color.RGBA
	Active    bool
}

// HUD renders score, countdown timer, chaos meter, equipped tool, and feedback popups.
type HUD struct {
	font       *PixelFont
	pixelImg   *ebiten.Image
	drawOpts   ebiten.DrawImageOptions
	popups     [8]PopupMessage
	pulseTimer float64
}

// NewHUD creates a new HUD renderer.
func NewHUD(font *PixelFont, pixelImg *ebiten.Image) *HUD {
	return &HUD{
		font:     font,
		pixelImg: pixelImg,
	}
}

// AddPopup triggers a floating score/status notification.
func (h *HUD) AddPopup(text string, x, y float64, col color.RGBA) {
	for i := range h.popups {
		if !h.popups[i].Active {
			h.popups[i] = PopupMessage{
				Text:   text,
				X:      x,
				Y:      y,
				Life:   1.2,
				Color:  col,
				Active: true,
			}
			return
		}
	}
}

// Update advances active popup animations.
func (h *HUD) Update(dt float64) {
	h.pulseTimer += dt * 6.0
	for i := range h.popups {
		p := &h.popups[i]
		if !p.Active {
			continue
		}
		p.Age += dt
		p.Y -= 15.0 * dt // Float upward
		if p.Age >= p.Life {
			p.Active = false
		}
	}
}

// Draw renders all HUD components onto the screen.
func (h *HUD) Draw(screen *ebiten.Image, score int, clutches int, timeLeft float64, chaos float64, tool entity.ToolType) {
	// Top Header Bar Background (320x16)
	h.drawRect(screen, 0, 0, 320, 16, color.RGBA{15, 12, 25, 230})
	h.drawRect(screen, 0, 16, 320, 1, color.RGBA{60, 50, 85, 255})

	// 1. Top-Left: Score and Clutch badge
	scoreStr := fmt.Sprintf("SCORE %06d", score)
	h.font.DrawText(screen, scoreStr, 6, 5, 1.0, color.RGBA{255, 240, 150, 255}, true)

	if clutches > 0 {
		clutchStr := fmt.Sprintf("CLUTCH x%d", clutches)
		h.font.DrawText(screen, clutchStr, 80, 5, 1.0, color.RGBA{80, 240, 255, 255}, true)
	}

	// 2. Top-Center: Countdown Timer (MM:SS)
	minutes := int(timeLeft) / 60
	seconds := int(timeLeft) % 60
	timeColor := color.RGBA{240, 240, 255, 255}
	if timeLeft <= 30.0 {
		// Blink red when under 30 seconds
		if int(h.pulseTimer)%2 == 0 {
			timeColor = color.RGBA{255, 60, 60, 255}
		}
	}
	timeStr := fmt.Sprintf("TIME %02d:%02d", minutes, seconds)
	h.font.DrawText(screen, timeStr, 138, 5, 1.0, timeColor, true)

	// 3. Top-Right: Chaos Meter (0% to 100%)
	chaosBarW := 50.0
	chaosBarH := 6.0
	chaosX := 260.0
	chaosY := 5.0

	h.font.DrawText(screen, "CHAOS", chaosX-26, 5, 1.0, color.RGBA{255, 180, 180, 255}, true)

	// Meter Frame
	h.drawRect(screen, chaosX-1, chaosY-1, chaosBarW+2, chaosBarH+2, color.RGBA{30, 20, 35, 255})
	h.drawRect(screen, chaosX, chaosY, chaosBarW, chaosBarH, color.RGBA{60, 30, 40, 255})

	// Meter Fill
	fillW := chaosBarW * (chaos / 100.0)
	fillCol := color.RGBA{240, 70, 70, 255}
	if chaos >= 75.0 {
		// Pulsing intense red
		alpha := uint8(180 + math.Sin(h.pulseTimer*1.5)*75)
		fillCol = color.RGBA{255, 20, 20, alpha}
	} else if chaos < 40.0 {
		fillCol = color.RGBA{240, 180, 40, 255}
	}
	h.drawRect(screen, chaosX, chaosY, fillW, chaosBarH, fillCol)

	// 4. Bottom-Center: Equipped Tool Display
	toolW := 90.0
	toolH := 12.0
	toolX := 160.0 - toolW/2
	toolY := 166.0

	h.drawRect(screen, toolX-1, toolY-1, toolW+2, toolH+2, color.RGBA{15, 12, 25, 200})
	h.drawRect(screen, toolX, toolY, toolW, toolH, color.RGBA{35, 30, 50, 220})
	entity.DrawToolIcon(screen, h.pixelImg, tool, toolX+8, toolY+6)
	h.font.DrawText(screen, tool.Name(), toolX+18, toolY+4, 1.0, color.RGBA{220, 230, 245, 255}, true)

	// 5. Draw active floating popups
	for i := range h.popups {
		p := &h.popups[i]
		if !p.Active {
			continue
		}
		prog := p.Age / p.Life
		alpha := uint8(255.0 * (1.0 - prog))
		c := p.Color
		c.A = alpha
		h.font.DrawText(screen, p.Text, p.X, p.Y, 1.0, c, true)
	}

	// 6. Emergency Strobe Border when Chaos is Critical (>= 75%)
	if chaos >= 75.0 {
		alpha := uint8(math.Abs(math.Sin(h.pulseTimer)) * 60)
		strobeCol := color.RGBA{255, 20, 20, alpha}
		// Red border vignette
		h.drawRect(screen, 0, 16, 320, 2, strobeCol)
		h.drawRect(screen, 0, 178, 320, 2, strobeCol)
		h.drawRect(screen, 0, 16, 2, 164, strobeCol)
		h.drawRect(screen, 318, 16, 2, 164, strobeCol)
	}
}

func (h *HUD) drawRect(screen *ebiten.Image, x, y, w, hDim float64, c color.RGBA) {
	h.drawOpts.GeoM.Reset()
	h.drawOpts.GeoM.Scale(w, hDim)
	h.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	h.drawOpts.ColorScale.Reset()
	h.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(h.pixelImg, &h.drawOpts)
}
