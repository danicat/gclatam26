package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gopher-budokai-panic/internal/entities"
)

type HUD struct {
	pulsePhase float64
}

func NewHUD() *HUD {
	return &HUD{}
}

func (h *HUD) Update(dt float64) {
	h.pulsePhase += 10.0 * dt
}

func (h *HUD) Draw(screen *ebiten.Image, p1, p2 *entities.Fighter, screenW, screenH float64) {
	if p1 == nil || p2 == nil {
		return
	}

	sw := float32(screenW)
	_ = screenH

	// P1 HUD (Left side)
	h.drawFighterHUD(screen, p1, 20, 16, 200, false, "GOKU (SSJ)")

	// P2 HUD (Right side)
	h.drawFighterHUD(screen, p2, sw-220, 16, 200, true, "VEGETA")

	// Center VS & Match Timer
	ebitenutil.DebugPrintAt(screen, "VS", int(sw/2)-8, 18)

	// Interactive PANIC / RECOVER prompt if Player is in PANIC!
	if p1.Panic.IsPanicked {
		h.drawPanicRecoverPrompt(screen, p1, sw/2, 110)
	}

	// Sparking Announcement
	if p1.IsSparking {
		sparkPulse := uint8(200 + math.Sin(h.pulsePhase)*55)
		sparkCol := color.RGBA{R: 255, G: 220, B: sparkPulse, A: 255}
		_ = sparkCol
		ebitenutil.DebugPrintAt(screen, "★ SPARKING! ★", 20, 68)
	}
	if p2.IsSparking {
		ebitenutil.DebugPrintAt(screen, "★ SPARKING! ★", int(sw)-110, 68)
	}
}

func (h *HUD) drawFighterHUD(screen *ebiten.Image, f *entities.Fighter, x, y, width float32, alignRight bool, name string) {
	borderCol := color.RGBA{R: 30, G: 30, B: 40, A: 240}
	bgCol := color.RGBA{R: 50, G: 15, B: 15, A: 200}

	// 1. Fighter Name
	nameX := int(x)
	if alignRight {
		nameX = int(x + width - float32(len(name)*7))
	}
	ebitenutil.DebugPrintAt(screen, name, nameX, int(y)-12)

	// 2. Health Bar (Multi-layered BT3 style)
	hBarH := float32(14.0)
	vector.DrawFilledRect(screen, x-2, y-2, width+4, hBarH+4, borderCol, false)
	vector.DrawFilledRect(screen, x, y, width, hBarH, bgCol, false)

	hpRatio := float32(f.Health / f.MaxHealth)
	if hpRatio > 0 {
		var hpCol color.RGBA
		if hpRatio > 0.5 {
			hpCol = color.RGBA{R: 50, G: 220, B: 80, A: 255} // Green
		} else if hpRatio > 0.25 {
			hpCol = color.RGBA{R: 240, G: 200, B: 40, A: 255} // Yellow
		} else {
			// Red pulsing low HP
			pulse := uint8(200 + math.Sin(h.pulsePhase)*55)
			hpCol = color.RGBA{R: pulse, G: 30, B: 30, A: 255}
		}

		fillW := width * hpRatio
		fillX := x
		if alignRight {
			fillX = x + width - fillW
		}
		vector.DrawFilledRect(screen, fillX, y, fillW, hBarH, hpCol, false)
		// Top highlight sheen
		sheenCol := color.RGBA{R: 255, G: 255, B: 255, A: 90}
		vector.DrawFilledRect(screen, fillX, y, fillW, 3, sheenCol, false)
	}

	// 3. Ki Bar
	kiY := y + hBarH + 4
	kiH := float32(8.0)
	vector.DrawFilledRect(screen, x-1, kiY-1, width*0.75+2, kiH+2, borderCol, false)
	vector.DrawFilledRect(screen, x, kiY, width*0.75, kiH, color.RGBA{R: 15, G: 25, B: 45, A: 200}, false)

	kiRatio := float32(f.Ki / f.MaxKi)
	if kiRatio > 0 {
		kiCol := color.RGBA{R: 40, G: 160, B: 255, A: 255}
		if f.IsSparking {
			kiCol = color.RGBA{R: 255, G: 230, B: 50, A: 255} // Gold during Sparking
		}
		fillKiW := width * 0.75 * kiRatio
		fillKiX := x
		if alignRight {
			fillKiX = x + width*0.75 - fillKiW
		}
		vector.DrawFilledRect(screen, fillKiX, kiY, fillKiW, kiH, kiCol, false)
	}

	// 4. PANIC METER (The Game Jam Theme!)
	panicY := kiY + kiH + 4
	panicH := float32(6.0)
	vector.DrawFilledRect(screen, x-1, panicY-1, width*0.6+2, panicH+2, borderCol, false)
	vector.DrawFilledRect(screen, x, panicY, width*0.6, panicH, color.RGBA{R: 35, G: 30, B: 30, A: 180}, false)

	panicRatio := float32(f.Panic.Meter / 100.0)
	if panicRatio > 0 {
		var panicCol color.RGBA
		if f.Panic.IsPanicked {
			// Flashing frantic red
			pPulse := uint8(200 + math.Sin(h.pulsePhase*1.5)*55)
			panicCol = color.RGBA{R: 255, G: pPulse, B: 0, A: 255}
		} else if panicRatio > 0.7 {
			panicCol = color.RGBA{R: 255, G: 120, B: 20, A: 255} // Orange
		} else {
			panicCol = color.RGBA{R: 230, G: 210, B: 180, A: 255} // Light
		}

		fillPW := width * 0.6 * panicRatio
		fillPX := x
		if alignRight {
			fillPX = x + width*0.6 - fillPW
		}
		vector.DrawFilledRect(screen, fillPX, panicY, fillPW, panicH, panicCol, false)
	}

	// Panic status tag
	if f.Panic.IsPanicked {
		pTagX := int(x)
		if alignRight {
			pTagX = int(x + width*0.6 - 55)
		}
		ebitenutil.DebugPrintAt(screen, "!PANIC!", pTagX, int(panicY+8))
	}
}

// drawPanicRecoverPrompt renders the high-tension interactive recovery UI
func (h *HUD) drawPanicRecoverPrompt(screen *ebiten.Image, p1 *entities.Fighter, centerX, centerY float32) {
	// Shaky offset
	shakeX := float32(math.Sin(h.pulsePhase*2.0) * 3.0)
	shakeY := float32(math.Cos(h.pulsePhase*2.0) * 2.0)

	boxW := float32(280.0)
	boxH := float32(48.0)
	boxX := centerX - boxW/2 + shakeX
	boxY := centerY + shakeY

	// Flashing warning banner
	borderCol := color.RGBA{R: 255, G: 30, B: 30, A: 240}
	bgCol := color.RGBA{R: 40, G: 0, B: 0, A: 220}
	vector.DrawFilledRect(screen, boxX-2, boxY-2, boxW+4, boxH+4, borderCol, false)
	vector.DrawFilledRect(screen, boxX, boxY, boxW, boxH, bgCol, false)

	// Text alerts
	alertText := "!!! PANIC !!! MASH [SPACE] TO RECOVER!"
	ebitenutil.DebugPrintAt(screen, alertText, int(boxX+14), int(boxY+6))

	// Recovery Effort Progress Bar
	barX := boxX + 20
	barY := boxY + 26
	barW := boxW - 40
	barH := float32(12.0)

	vector.DrawFilledRect(screen, barX-1, barY-1, barW+2, barH+2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{R: 20, G: 20, B: 30, A: 255}, false)

	effortRatio := float32(p1.Panic.RecoverEffort / 100.0)
	if effortRatio > 1.0 {
		effortRatio = 1.0
	}
	if effortRatio > 0 {
		vector.DrawFilledRect(screen, barX, barY, barW*effortRatio, barH, color.RGBA{R: 80, G: 240, B: 255, A: 255}, false)
	}

	pctText := fmt.Sprintf("RECOVERY: %d%%", int(effortRatio*100))
	ebitenutil.DebugPrintAt(screen, pctText, int(boxX+barW/2-25), int(barY-1))
}
