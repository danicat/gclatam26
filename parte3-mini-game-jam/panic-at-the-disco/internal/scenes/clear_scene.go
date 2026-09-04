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

type StageClearScene struct {
	ae         *audio.AudioEngine
	nextZone   levels.ZoneID
	score      int
	lives      int
	timer      float64
	particles  *gfx.ParticleSystem
}

func NewStageClearScene(nextZone levels.ZoneID, score, lives int) *StageClearScene {
	return &StageClearScene{
		ae:        audio.GetAudioEngine(),
		nextZone:  nextZone,
		score:     score,
		lives:     lives,
		particles: gfx.NewParticleSystem(150),
	}
}

func (sc *StageClearScene) Enter() {
	sc.timer = 0
	sc.ae.PlaySFXWin()
}

func (sc *StageClearScene) Exit() {}

func (sc *StageClearScene) Update(dt float64) SceneAction {
	sc.timer += dt
	sc.particles.Update(dt)

	// Burst victory confetti
	if sc.timer < 0.5 && math.Mod(sc.timer, 0.08) < dt {
		for i := 0; i < 8; i++ {
			ang := float64(i) * (2 * math.Pi / 8)
			sc.particles.Emit(
				320.0, 160.0,
				math.Cos(ang)*90.0, math.Sin(ang)*90.0,
				0.8, 3.0, 1.0,
				color.RGBA{255, 215, 0, 255},
				gfx.ParticleMirrorShard,
			)
		}
	}

	in := input.Poll()
	if in.ConfirmJustDown || in.DashJustDown || sc.timer >= 2.5 {
		return SceneAction{
			Type:       ActionSwitchScene,
			NextScene:  ScenePlay,
			TargetZone: sc.nextZone,
			Score:      sc.score,
			Lives:      sc.lives,
		}
	}

	return SceneAction{Type: ActionNone}
}

func (sc *StageClearScene) Draw(screen *ebiten.Image) {
	// Dark backdrop
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{10, 15, 25, 255}, false)

	sc.particles.Draw(screen)

	// Stage Cleared Banner
	title := "STAGE EVACUATED!"
	tw := gfx.MeasureText(title, 2.8)
	gfx.DrawText(screen, title, (640-tw)/2, 110.0, 2.8, color.RGBA{0, 255, 180, 255}, true)

	cfg := levels.GetLevelConfig(sc.nextZone)
	nextStr := fmt.Sprintf("UP NEXT: %s", cfg.Name)
	nw := gfx.MeasureText(nextStr, 1.4)
	gfx.DrawText(screen, nextStr, (640-nw)/2, 170.0, 1.4, color.RGBA{255, 215, 0, 255}, true)

	scoreStr := fmt.Sprintf("CURRENT SCORE: %06d", sc.score)
	sw := gfx.MeasureText(scoreStr, 1.3)
	gfx.DrawText(screen, scoreStr, (640-sw)/2, 205.0, 1.3, color.RGBA{240, 240, 250, 255}, true)

	// Prompt
	prompt := "PRESS SPACE TO ENTER NEXT ZONE"
	pw := gfx.MeasureText(prompt, 1.2)
	gfx.DrawText(screen, prompt, (640-pw)/2, 275.0, 1.2, color.RGBA{0, 240, 255, 255}, true)
}
