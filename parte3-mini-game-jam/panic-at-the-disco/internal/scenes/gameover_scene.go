package scenes

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"panic-at-the-disco/internal/audio"
	"panic-at-the-disco/internal/gfx"
	"panic-at-the-disco/internal/input"
	"panic-at-the-disco/internal/levels"
)

type GameOverScene struct {
	ae          *audio.AudioEngine
	reason      string
	score       int
	surviveTime float64
	timer       float64
}

func NewGameOverScene(reason string, score int, surviveTime float64) *GameOverScene {
	return &GameOverScene{
		ae:          audio.GetAudioEngine(),
		reason:      reason,
		score:       score,
		surviveTime: surviveTime,
	}
}

func (gos *GameOverScene) Enter() {
	gos.ae.StopBGM()
	gos.timer = 0
}

func (gos *GameOverScene) Exit() {}

func (gos *GameOverScene) Update(dt float64) SceneAction {
	gos.timer += dt
	in := input.Poll()

	if in.RestartJustDown || in.ConfirmJustDown {
		return SceneAction{
			Type:       ActionSwitchScene,
			NextScene:  ScenePlay,
			TargetZone: levels.ZoneDanceFloor,
			Lives:      3,
			Score:      0,
		}
	}

	if in.PauseJustDown {
		return SceneAction{
			Type:      ActionSwitchScene,
			NextScene: SceneTitle,
		}
	}

	return SceneAction{Type: ActionNone}
}

func (gos *GameOverScene) Draw(screen *ebiten.Image) {
	// Red strobe tint
	pulse := math.Sin(gos.timer * 4.0)
	bgR := uint8(40 + 20*pulse)
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{bgR, 8, 12, 255}, false)

	// "GAME OVER"
	title := "GAME OVER"
	tw := gfx.MeasureText(title, 3.4)
	gfx.DrawText(screen, title, (640-tw)/2, 85.0, 3.4, color.RGBA{255, 30, 30, 255}, true)

	// Reason
	rw := gfx.MeasureText(gos.reason, 1.4)
	gfx.DrawText(screen, gos.reason, (640-rw)/2, 145.0, 1.4, color.RGBA{255, 180, 50, 255}, true)

	// Stats
	scoreStr := fmt.Sprintf("FINAL SCORE: %06d", gos.score)
	sw := gfx.MeasureText(scoreStr, 1.3)
	gfx.DrawText(screen, scoreStr, (640-sw)/2, 195.0, 1.3, color.RGBA{240, 240, 250, 255}, true)

	timeStr := fmt.Sprintf("TIME SURVIVED: %04.1f SECONDS", gos.surviveTime)
	tmw := gfx.MeasureText(timeStr, 1.2)
	gfx.DrawText(screen, timeStr, (640-tmw)/2, 220.0, 1.2, color.RGBA{200, 200, 220, 255}, true)

	// Action Prompts
	retryStr := "[PRESS R OR ENTER TO RETRY]"
	retW := gfx.MeasureText(retryStr, 1.3)
	blink := math.Sin(gos.timer * 6.0)
	retryCol := color.RGBA{255, 220, 0, 255}
	if blink < 0 {
		retryCol = color.RGBA{180, 160, 0, 255}
	}
	gfx.DrawText(screen, retryStr, (640-retW)/2, 280.0, 1.3, retryCol, true)

	exitStr := "[PRESS ESC FOR TITLE SCREEN]"
	ew := gfx.MeasureText(exitStr, 1.0)
	gfx.DrawText(screen, exitStr, (640-ew)/2, 310.0, 1.0, color.RGBA{160, 160, 180, 255}, true)
}
