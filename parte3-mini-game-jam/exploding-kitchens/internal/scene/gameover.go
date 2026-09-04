package scene

import (
	"fmt"
	"image/color"
	"math"

	"exploding-kitchens/internal/game"
	"exploding-kitchens/internal/system"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameOverScene presents the post-shift performance debrief, star rating, and high scores.
type GameOverScene struct {
	ctx       *game.GameContext
	animTimer float64
	drawOpts  ebiten.DrawImageOptions
}

// NewGameOverScene creates a new GameOverScene.
func NewGameOverScene() *GameOverScene {
	return &GameOverScene{}
}

// Enter plays the conclusion fanfare or alert.
func (gos *GameOverScene) Enter(ctx *game.GameContext) {
	gos.ctx = ctx
	gos.animTimer = 0

	if ctx.LastSurvived {
		ctx.Audio.PlayClutch()
	} else {
		ctx.Audio.PlayExplosion()
	}
}

// Update listens for replay or menu navigation.
func (gos *GameOverScene) Update(dt float64, in system.InputState) (game.SceneID, error) {
	gos.animTimer += dt * 3.0

	// Restart immediate gameplay
	if in.ConfirmJust || in.InteractJust {
		gos.ctx.Audio.PlayPickup()
		return game.ScenePlay, nil
	}

	// Return to title
	if in.PauseJust {
		gos.ctx.Audio.PlayPickup()
		return game.SceneTitle, nil
	}

	return game.SceneKeepCurrent, nil
}

// Draw renders the detailed shift debrief modal.
func (gos *GameOverScene) Draw(screen *ebiten.Image) {
	// Dark backdrop with kitchen silhouette
	gos.drawRect(screen, 0, 0, 320, 180, color.RGBA{18, 14, 28, 255})

	// Header banner
	headerCol := color.RGBA{255, 60, 60, 255}
	headerText := "KITCHEN MELTDOWN!"
	subText := "THE RESTAURANT DETONATED IN FLAMES!"

	if gos.ctx.LastSurvived {
		headerCol = color.RGBA{80, 240, 160, 255}
		headerText = "SHIFT COMPLETED! VICTORY!"
		subText = "YOU SURVIVED THE PANIC & RECOVER RUSH!"
	}

	headerBob := math.Sin(gos.animTimer) * 1.5
	gos.ctx.Font.DrawText(screen, headerText, 62, 22+headerBob, 1.4, headerCol, true)
	gos.ctx.Font.DrawText(screen, subText, 48, 40+headerBob, 0.9, color.RGBA{220, 220, 240, 255}, true)

	// Summary Card
	cardX := 50.0
	cardY := 56.0
	cardW := 220.0
	cardH := 82.0

	gos.drawRect(screen, cardX-1, cardY-1, cardW+2, cardH+2, color.RGBA{35, 28, 48, 255})
	gos.drawRect(screen, cardX, cardY, cardW, cardH, color.RGBA{28, 22, 40, 255})

	// Debrief Metrics
	scoreLine := fmt.Sprintf("FINAL SCORE     : %06d", gos.ctx.LastScore)
	clutchLine := fmt.Sprintf("CLUTCHES PULLED : %d", gos.ctx.LastClutches)
	explLine := fmt.Sprintf("DETONATIONS     : %d", gos.ctx.LastExplosions)
	recordLine := fmt.Sprintf("TOP RECORD      : %06d", gos.ctx.HighScore)

	gos.ctx.Font.DrawText(screen, scoreLine, cardX+15, cardY+12, 1.0, color.RGBA{255, 240, 140, 255}, true)
	gos.ctx.Font.DrawText(screen, clutchLine, cardX+15, cardY+26, 1.0, color.RGBA{100, 240, 255, 255}, true)
	gos.ctx.Font.DrawText(screen, explLine, cardX+15, cardY+40, 1.0, color.RGBA{255, 120, 120, 255}, true)
	gos.ctx.Font.DrawText(screen, recordLine, cardX+15, cardY+54, 1.0, color.RGBA{180, 255, 180, 255}, true)

	// Stars Rating
	stars := "★☆☆"
	if gos.ctx.LastSurvived && gos.ctx.LastScore >= 3000 {
		stars = "★★★"
	} else if gos.ctx.LastScore >= 1500 {
		stars = "★★☆"
	}
	gos.ctx.Font.DrawText(screen, "CHEF RATING: "+stars, cardX+15, cardY+68, 1.0, color.RGBA{255, 215, 60, 255}, true)

	// Flashing Action Prompts
	if int(gos.animTimer*2)%2 == 0 {
		gos.ctx.Font.DrawText(screen, "PRESS SPACE OR ENTER TO PLAY AGAIN", 58, 150, 1.0, color.RGBA{255, 245, 120, 255}, true)
	}
	gos.ctx.Font.DrawText(screen, "PRESS ESC FOR MAIN MENU", 98, 164, 0.9, color.RGBA{160, 160, 190, 255}, false)
}

func (gos *GameOverScene) drawRect(screen *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	gos.drawOpts.GeoM.Reset()
	gos.drawOpts.GeoM.Scale(w, h)
	gos.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	gos.drawOpts.ColorScale.Reset()
	gos.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(gos.ctx.PixelImg, &gos.drawOpts)
}

// Exit cleans up scene resources.
func (gos *GameOverScene) Exit() {}
