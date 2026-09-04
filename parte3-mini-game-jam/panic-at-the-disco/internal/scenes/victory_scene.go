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

type VictoryScene struct {
	ae          *audio.AudioEngine
	score       int
	lives       int
	surviveTime float64
	timer       float64
	particles   *gfx.ParticleSystem
}

func NewVictoryScene(score, lives int, surviveTime float64) *VictoryScene {
	return &VictoryScene{
		ae:          audio.GetAudioEngine(),
		score:       score,
		lives:       lives,
		surviveTime: surviveTime,
		particles:   gfx.NewParticleSystem(300),
	}
}

func (vs *VictoryScene) Enter() {
	vs.ae.PlayBGM(false)
	vs.ae.PlaySFXWin()
	vs.timer = 0
}

func (vs *VictoryScene) Exit() {
	vs.ae.StopBGM()
}

func (vs *VictoryScene) Update(dt float64) SceneAction {
	vs.timer += dt
	vs.particles.Update(dt)

	// Continuous festive confetti
	if math.Mod(vs.timer, 0.05) < dt {
		for i := 0; i < 4; i++ {
			x := 40.0 + math.Mod(vs.timer*240.0+float64(i*130), 560.0)
			vs.particles.Emit(
				x, -10.0,
				(math.Sin(vs.timer*8.0+float64(i)))*40.0, 70.0+float64(i*20),
				1.8, 3.5, 1.0,
				color.RGBA{R: uint8(150 + i*30), G: 255, B: uint8(100 + i*40), A: 255},
				gfx.ParticleMirrorShard,
			)
		}
	}

	in := input.Poll()
	if in.ConfirmJustDown || in.DashJustDown {
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

func (vs *VictoryScene) Draw(screen *ebiten.Image) {
	// Sunrise gradient background (Deep violet to golden sunrise)
	for y := 0; y < 360; y += 4 {
		prog := float64(y) / 360.0
		r := uint8(20 + 80*prog)
		g := uint8(10 + 40*prog)
		b := uint8(50 - 20*prog)
		vector.DrawFilledRect(screen, 0, float32(y), 640, 4, color.RGBA{r, g, b, 255}, false)
	}

	// Draw confetti
	vs.particles.Draw(screen)

	// Outdoor getaway taxi / street silhouette at bottom
	vector.DrawFilledRect(screen, 0, 290, 640, 70, color.RGBA{15, 12, 20, 255}, false)
	vector.StrokeLine(screen, 0, 290, 640, 290, 2.0, color.RGBA{255, 215, 0, 200}, false)

	// Triumphant afro dancer posing in safety
	gfx.DrawPlayer(screen, 320.0, 275.0, 1.0, vs.timer, false, false, 0.0)

	// Victory Header
	title := "YOU ESCAPED THE DISCO INFERNO!"
	tw := gfx.MeasureText(title, 2.2)
	pulse := math.Sin(vs.timer * 5.0)
	titleCol := color.RGBA{
		R: 255,
		G: uint8(215 + 40*pulse),
		B: 0,
		A: 255,
	}
	gfx.DrawText(screen, title, (640-tw)/2, 60.0, 2.2, titleCol, true)

	sub := "SATURDAY NIGHT FEVER NEVER DIES!"
	sw := gfx.MeasureText(sub, 1.4)
	gfx.DrawText(screen, sub, (640-sw)/2, 100.0, 1.4, color.RGBA{0, 255, 220, 255}, true)

	// Results Card
	vector.DrawFilledRect(screen, 160, 130, 320, 95, color.RGBA{20, 15, 35, 210}, false)
	vector.StrokeRect(screen, 160, 130, 320, 95, 2.0, color.RGBA{0, 240, 255, 220}, false)

	scoreStr := fmt.Sprintf("TOTAL SCORE: %06d", vs.score)
	gfx.DrawText(screen, scoreStr, 185.0, 145.0, 1.3, color.RGBA{255, 255, 255, 255}, true)

	livesStr := fmt.Sprintf("LIVES REMAINING: %d", vs.lives)
	gfx.DrawText(screen, livesStr, 185.0, 170.0, 1.3, color.RGBA{255, 100, 150, 255}, true)

	timeStr := fmt.Sprintf("TIME ELAPSED: %04.1f SECONDS", vs.surviveTime)
	gfx.DrawText(screen, timeStr, 185.0, 195.0, 1.3, color.RGBA{255, 220, 100, 255}, true)

	// Prompt
	prompt := "[PRESS ENTER OR SPACE TO DANCE AGAIN]"
	pw := gfx.MeasureText(prompt, 1.3)
	blink := math.Sin(vs.timer * 6.0)
	promptCol := color.RGBA{0, 255, 150, 255}
	if blink < 0 {
		promptCol = color.RGBA{0, 180, 100, 255}
	}
	gfx.DrawText(screen, prompt, (640-pw)/2, 330.0, 1.3, promptCol, true)
}
