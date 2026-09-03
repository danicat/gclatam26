package scene

import (
	"fmt"
	"image/color"
	"math"

	"exploding-kitchens/internal/entity"
	"exploding-kitchens/internal/game"
	"exploding-kitchens/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

// TitleScene displays the main menu, controls, high score, and manages the 10s Attract Demo timeout.
type TitleScene struct {
	ctx         *game.GameContext
	idleTimer   float64
	animTimer   float64
	drawOpts    ebiten.DrawImageOptions
	catAnimX    float64
	catFacing   bool
}

// NewTitleScene creates a new TitleScene.
func NewTitleScene() *TitleScene {
	return &TitleScene{}
}

// Enter resets idle timers and plays the menu theme.
func (ts *TitleScene) Enter(ctx *game.GameContext) {
	ts.ctx = ctx
	ts.idleTimer = 0
	ts.animTimer = 0
	ts.catAnimX = 40.0
	ts.catFacing = true
	ctx.Audio.StartBGM()
}

// Update listens for start triggers and monitors idle timeout for Attract Demo mode.
func (ts *TitleScene) Update(dt float64, in system.InputState) (game.SceneID, error) {
	ts.animTimer += dt * 3.0

	// Animated decorative cat walking across bottom
	if ts.catFacing {
		ts.catAnimX += 30.0 * dt
		if ts.catAnimX > 270 {
			ts.catFacing = false
		}
	} else {
		ts.catAnimX -= 30.0 * dt
		if ts.catAnimX < 40 {
			ts.catFacing = true
		}
	}

	// User pressed Start / Confirm
	if in.ConfirmJust || in.InteractJust {
		ts.ctx.Audio.PlayPickup()
		ts.ctx.DemoMode = false
		return game.ScenePlay, nil
	}

	// Any movement resets idle timer
	if in.MoveX != 0 || in.MoveY != 0 || in.PauseJust || in.DropJust {
		ts.idleTimer = 0
	} else {
		ts.idleTimer += dt
	}

	// 10-second idle timeout triggers Attract Demo Mode
	if ts.idleTimer >= 10.0 {
		ts.ctx.DemoMode = true
		return game.ScenePlay, nil
	}

	return game.SceneKeepCurrent, nil
}

// Draw renders the title screen layout with animated retro typography.
func (ts *TitleScene) Draw(screen *ebiten.Image) {
	// Background gradient / tiles
	ts.drawRect(screen, 0, 0, 320, 180, color.RGBA{22, 18, 35, 255})

	// Checkerboard lower banner
	for x := 0.0; x < 320; x += 16 {
		for y := 130.0; y < 180; y += 16 {
			isEven := (int(x/16) + int(y/16)) % 2 == 0
			c := color.RGBA{35, 30, 52, 255}
			if isEven {
				c = color.RGBA{45, 40, 68, 255}
			}
			ts.drawRect(screen, x, y, 16, 16, c)
		}
	}

	// Decorative Top & Bottom Borders
	ts.drawRect(screen, 0, 12, 320, 2, color.RGBA{240, 160, 40, 255})
	ts.drawRect(screen, 0, 128, 320, 2, color.RGBA{240, 160, 40, 255})

	// Title Animation (Slight floating bob)
	titleBob := math.Sin(ts.animTimer) * 2.0
	ts.ctx.Font.DrawText(screen, "EXPLODING KITCHENS", 50, 28+titleBob, 1.8, color.RGBA{255, 75, 65, 255}, true)
	ts.ctx.Font.DrawText(screen, "- PANIC AND RECOVER -", 88, 52+titleBob, 1.0, color.RGBA{255, 220, 90, 255}, true)

	// High Score
	if ts.ctx.HighScore > 0 {
		hsStr := fmt.Sprintf("TOP CHEF RECORD: %06d", ts.ctx.HighScore)
		ts.ctx.Font.DrawText(screen, hsStr, 96, 68, 1.0, color.RGBA{110, 240, 255, 255}, true)
	}

	// Controls Summary Box
	boxX := 52.0
	boxY := 80.0
	boxW := 216.0
	boxH := 42.0

	ts.drawRect(screen, boxX-1, boxY-1, boxW+2, boxH+2, color.RGBA{10, 8, 18, 220})
	ts.drawRect(screen, boxX, boxY, boxW, boxH, color.RGBA{30, 25, 45, 230})

	ts.ctx.Font.DrawText(screen, "WASD / ARROWS : MOVE CHEF", boxX+10, boxY+7, 1.0, color.RGBA{220, 230, 245, 255}, false)
	ts.ctx.Font.DrawText(screen, "SPACE / E     : DEFUSE / PICKUP / SHOO CAT", boxX+10, boxY+18, 1.0, color.RGBA{220, 230, 245, 255}, false)
	ts.ctx.Font.DrawText(screen, "Q / SHIFT     : DROP HELD TOOL", boxX+10, boxY+29, 1.0, color.RGBA{220, 230, 245, 255}, false)

	// Flashing "PRESS SPACE TO START"
	if int(ts.animTimer*2)%2 == 0 {
		ts.ctx.Font.DrawText(screen, "PRESS SPACE OR ENTER TO START", 68, 142, 1.0, color.RGBA{255, 245, 120, 255}, true)
	}

	// Attract Demo timer hint
	demoSec := int(10.0 - ts.idleTimer)
	if demoSec < 5 && demoSec > 0 {
		demoStr := fmt.Sprintf("DEMO IN %d SEC", demoSec)
		ts.ctx.Font.DrawText(screen, demoStr, 126, 160, 0.9, color.RGBA{150, 160, 190, 200}, false)
	}

	// Animated cat running across bottom
	catY := 158.0
	entity.DrawToolIcon(screen, ts.ctx.PixelImg, entity.ToolExtinguisher, 28, catY-2)
	entity.DrawToolIcon(screen, ts.ctx.PixelImg, entity.ToolIce, 292, catY-2)
}

func (ts *TitleScene) drawRect(screen *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	ts.drawOpts.GeoM.Reset()
	ts.drawOpts.GeoM.Scale(w, h)
	ts.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	ts.drawOpts.ColorScale.Reset()
	ts.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(ts.ctx.PixelImg, &ts.drawOpts)
}

// Exit stops title background loops.
func (ts *TitleScene) Exit() {}
