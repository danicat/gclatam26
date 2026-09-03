package scenes

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"panic-at-the-disco/internal/audio"
	"panic-at-the-disco/internal/gfx"
	"panic-at-the-disco/internal/input"
	"panic-at-the-disco/internal/levels"
)

type TitleScene struct {
	ae         *audio.AudioEngine
	floor      *gfx.DiscoFloor
	particles  *gfx.ParticleSystem
	timer      float64
	discoBallY float64
}

func NewTitleScene() *TitleScene {
	return &TitleScene{
		ae:         audio.GetAudioEngine(),
		floor:      gfx.NewDiscoFloor(0, 0, 640, 360, 16, 9, 120.0),
		particles:  gfx.NewParticleSystem(100),
		discoBallY: 65.0,
	}
}

func (ts *TitleScene) Enter() {
	ts.ae.PlayBGM(false)
	ts.timer = 0
}

func (ts *TitleScene) Exit() {
	// Keep BGM running or let PlayScene handle it
}

func (ts *TitleScene) Update(dt float64) SceneAction {
	ts.timer += dt
	ts.floor.Update(dt)
	ts.particles.Update(dt)

	// Sparkle particles around title
	if math.Mod(ts.timer, 0.15) < dt {
		ts.particles.Emit(
			320.0+(math.Sin(ts.timer*12.0)*180.0),
			ts.discoBallY+(math.Cos(ts.timer*15.0)*30.0),
			(math.Sin(ts.timer*8.0))*30.0,
			-20.0,
			0.4, 2.5, 0.5,
			color.RGBA{255, 255, 255, 255},
			gfx.ParticleMirrorShard,
		)
	}

	in := input.Poll()
	if in.ConfirmJustDown || in.DashJustDown {
		ts.ae.PlaySFXDash()
		return SceneAction{
			Type:       ActionSwitchScene,
			NextScene:  ScenePlay,
			TargetZone: levels.ZoneDanceFloor,
			Lives:      3,
			Score:      0,
		}
	}

	return SceneAction{Type: ActionNone}
}

func (ts *TitleScene) Draw(screen *ebiten.Image) {
	// 1. Background disco floor (dimmed)
	ts.floor.Draw(screen)

	// Dark overlay for readability
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{15, 10, 25, 200}, false)

	// 2. Sparkling Disco Ball in center
	gfx.DrawDiscoBall(screen, 320.0, ts.discoBallY, 32.0, ts.timer*2.0)
	ts.particles.Draw(screen)

	// 3. Main Title: "PANIC! AT THE DISCO"
	titlePulse := math.Sin(ts.timer * 4.0)
	titleCol := color.RGBA{
		R: uint8(255),
		G: uint8(200 + 55*titlePulse),
		B: uint8(50),
		A: 255,
	}
	title := "PANIC! AT THE DISCO"
	tw := gfx.MeasureText(title, 3.0)
	gfx.DrawText(screen, title, (640-tw)/2, 115.0, 3.0, titleCol, true)

	// Subtitle: "SATURDAY NIGHT FLEE-VER"
	subTitle := "SATURDAY NIGHT FLEE-VER"
	subCol := color.RGBA{0, 240, 255, 255}
	stw := gfx.MeasureText(subTitle, 2.0)
	gfx.DrawText(screen, subTitle, (640-stw)/2, 148.0, 2.0, subCol, true)

	// 4. Instructions Card
	cardX := float32(110.0)
	cardY := float32(178.0)
	cardW := float32(420.0)
	cardH := float32(115.0)

	vector.DrawFilledRect(screen, cardX, cardY, cardW, cardH, color.RGBA{20, 15, 35, 220}, false)
	vector.StrokeRect(screen, cardX, cardY, cardW, cardH, 2.0, color.RGBA{255, 0, 140, 240}, false)

	// Instruction text
	gfx.DrawText(screen, "HOW TO SURVIVE THE COLLAPSE:", 135.0, 190.0, 1.5, color.RGBA{255, 215, 0, 255}, true)
	gfx.DrawText(screen, "- [W,A,S,D] OR [ARROWS] : MOVE DANCER", 135.0, 212.0, 1.2, color.RGBA{240, 240, 250, 255}, true)
	gfx.DrawText(screen, "- [SPACE] OR [J]        : DISCO DASH (INVULNERABLE)", 135.0, 228.0, 1.2, color.RGBA{240, 240, 250, 255}, true)
	gfx.DrawText(screen, "- DODGE SHADOWS & HAZARDS BEFORE ROOF FALLS!", 135.0, 244.0, 1.2, color.RGBA{255, 120, 120, 255}, true)
	gfx.DrawText(screen, "- REACH THE GREEN EMERGENCY EXIT IN TIME!", 135.0, 260.0, 1.2, color.RGBA{100, 255, 120, 255}, true)

	// 5. Pulsing Start Prompt
	startAlpha := math.Sin(ts.timer * 6.0)
	if startAlpha > -0.2 {
		prompt := ">>> PRESS SPACE OR ENTER TO ESCAPE <<<"
		pw := gfx.MeasureText(prompt, 1.8)
		promptCol := color.RGBA{0, 255, 180, 255}
		gfx.DrawText(screen, prompt, (640-pw)/2, 312.0, 1.8, promptCol, true)
	}
}
