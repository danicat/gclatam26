package app

import (
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"

	"panic-recover/internal/game"
	"panic-recover/internal/sound"
)

const (
	windowWidth  = 960
	windowHeight = 540
	fixedDT      = time.Second / 60
)

var (
	colorBackground = color.RGBA{R: 10, G: 12, B: 24, A: 255}
	colorGrid       = color.RGBA{R: 25, G: 31, B: 53, A: 255}
	colorText       = color.RGBA{R: 220, G: 238, B: 228, A: 255}
	colorCalm       = color.RGBA{R: 72, G: 224, B: 126, A: 255}
	colorPanic      = color.RGBA{R: 244, G: 72, B: 82, A: 255}
	colorRecover    = color.RGBA{R: 62, G: 220, B: 244, A: 255}
	colorBug        = color.RGBA{R: 238, G: 76, B: 190, A: 255}
)

type App struct {
	model            *game.Model
	sprites          spriteAssets
	particles        *ParticleSystem
	audio            *soundSystem
	audioAttempted   bool
	frameTime        float64
	cycleText        string
	bugsText         string
	stateText        string
	lastPhase        game.Phase
	lastScene        game.Scene
	lastCycle        int
	lastEliminations int
}

func New() *App {
	a := &App{model: game.NewModel(), particles: NewParticleSystem(128)}
	if sprites, err := loadEmbeddedAssets(); err == nil {
		a.sprites = sprites
	}
	a.updateHUD()
	a.lastPhase = a.model.Phase
	a.lastScene = a.model.Scene
	a.lastCycle = a.model.Cycle
	return a
}

func (a *App) Layout(_, _ int) (int, int) {
	return game.VirtualWidth, game.VirtualHeight
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	input := game.Input{
		Move: game.Vec2{
			X: axis(ebiten.KeyD, ebiten.KeyRight) - axis(ebiten.KeyA, ebiten.KeyLeft),
			Y: axis(ebiten.KeyS, ebiten.KeyDown) - axis(ebiten.KeyW, ebiten.KeyUp),
		},
		PanicPressed:   inpututil.IsKeyJustPressed(ebiten.KeySpace),
		StartPressed:   inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
		RestartPressed: inpututil.IsKeyJustPressed(ebiten.KeyR),
	}
	a.model.Update(input, fixedDT)
	a.frameTime += fixedDT.Seconds()
	a.updateAudio()
	a.updateParticles()
	a.updateHUD()
	return nil
}

func axis(positive, negative ebiten.Key) float64 {
	var result float64
	if ebiten.IsKeyPressed(positive) {
		result++
	}
	if ebiten.IsKeyPressed(negative) {
		result--
	}
	return result
}

func (a *App) updateHUD() {
	a.cycleText = "CYCLE " + strconv.Itoa(a.model.Cycle+1) + "/" + strconv.Itoa(len(a.model.Config.Cycles))
	quota := 0
	if a.model.Scene == game.ScenePlaying && a.model.Cycle < len(a.model.Config.Cycles) {
		quota = a.model.Config.Cycles[a.model.Cycle].Quota
	}
	a.bugsText = "BUGS " + strconv.Itoa(a.model.Eliminations) + "/" + strconv.Itoa(quota)
	switch a.model.Scene {
	case game.SceneTitle:
		a.stateText = "PRESS ENTER TO START"
	case game.SceneVictory:
		a.stateText = "RECOVER SUCCESSFUL — PRESS R"
	case game.SceneGameOver:
		a.stateText = "STABILITY EXHAUSTED — PRESS R"
	default:
		switch a.model.Phase {
		case game.PhaseCalm:
			a.stateText = "CALM"
		case game.PhasePanic:
			a.stateText = "PANIC"
		case game.PhaseRecoverAvailable:
			a.stateText = "RECOVER()"
		}
	}
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)
	a.drawGrid(screen)

	switch a.model.Scene {
	case game.SceneTitle:
		a.drawTitle(screen)
	case game.ScenePlaying:
		a.drawPlayfield(screen)
	case game.SceneVictory, game.SceneGameOver:
		a.drawPlayfield(screen)
		a.drawResult(screen)
	}
}

func (a *App) drawGrid(screen *ebiten.Image) {
	for x := 0; x <= game.VirtualWidth; x += 16 {
		vector.StrokeLine(screen, float32(x), 0, float32(x), game.VirtualHeight, 1, colorGrid, false)
	}
	for y := 0; y <= game.VirtualHeight; y += 16 {
		vector.StrokeLine(screen, 0, float32(y), game.VirtualWidth, float32(y), 1, colorGrid, false)
	}
}

func (a *App) drawTitle(screen *ebiten.Image) {
	text.Draw(screen, "NEON BREATHER", basicfont.Face7x13, 106, 58, colorCalm)
	text.Draw(screen, "MOVE: WASD / ARROWS", basicfont.Face7x13, 90, 91, colorText)
	text.Draw(screen, "PANIC: SPACE — DESTROY BUGS", basicfont.Face7x13, 90, 106, colorPanic)
	text.Draw(screen, "REACH RECOVER BEFORE STABILITY HITS ZERO", basicfont.Face7x13, 53, 121, colorRecover)
	text.Draw(screen, a.stateText, basicfont.Face7x13, 111, 151, colorText)
}

func (a *App) drawPlayfield(screen *ebiten.Image) {
	if a.model.Recover.Active {
		pulse := float32(1 + 0.2*math.Sin(a.frameTime*8))
		vector.DrawFilledCircle(screen, float32(a.model.Recover.Position.X), float32(a.model.Recover.Position.Y), float32(a.model.Recover.Radius)*pulse, colorRecover, false)
	}
	for _, bug := range a.model.Bugs {
		if bug.Alive {
			if a.sprites.bug != nil {
				drawSprite(screen, a.sprites.bug, bug.Position, 15)
			} else {
				vector.DrawFilledCircle(screen, float32(bug.Position.X), float32(bug.Position.Y), float32(bug.Radius), colorBug, false)
			}
		}
	}
	playerColor := colorCalm
	if a.model.Phase == game.PhasePanic {
		playerColor = colorPanic
	}
	if a.sprites.gopher != nil {
		if a.model.Phase == game.PhasePanic {
			drawSpriteWithTint(screen, a.sprites.gopher, a.model.Player.Position, 18, color.RGBA{R: 255, G: 96, B: 96, A: 255})
		} else {
			drawSprite(screen, a.sprites.gopher, a.model.Player.Position, 18)
		}
	} else {
		vector.DrawFilledCircle(screen, float32(a.model.Player.Position.X), float32(a.model.Player.Position.Y), float32(a.model.Player.Radius), playerColor, false)
	}
	if a.particles != nil {
		a.particles.Draw(screen)
	}
	text.Draw(screen, a.cycleText, basicfont.Face7x13, 8, 14, colorText)
	text.Draw(screen, a.bugsText, basicfont.Face7x13, 252, 14, colorText)
	a.drawStability(screen)
	text.Draw(screen, a.stateText, basicfont.Face7x13, 8, 174, playerColor)
}

func (a *App) updateParticles() {
	if a.particles == nil {
		return
	}
	a.particles.Update(fixedDT.Seconds())
	if a.model.Phase != a.lastPhase {
		particleColor := colorCalm
		if a.model.Phase == game.PhasePanic {
			particleColor = colorPanic
		} else if a.model.Phase == game.PhaseRecoverAvailable {
			particleColor = colorRecover
		}
		a.spawnBurst(a.model.Player.Position, particleColor, 12, 16)
	}
	if a.model.Eliminations > a.lastEliminations {
		a.spawnBurst(a.model.Player.Position, colorRecover, 8, 12)
	}
	if a.model.Scene != a.lastScene {
		if a.model.Scene == game.SceneVictory {
			a.spawnBurst(a.model.Player.Position, colorRecover, 20, 22)
		} else if a.model.Scene == game.SceneGameOver {
			a.spawnBurst(a.model.Player.Position, colorPanic, 20, 22)
		}
	}
	a.lastPhase = a.model.Phase
	a.lastScene = a.model.Scene
	a.lastCycle = a.model.Cycle
	a.lastEliminations = a.model.Eliminations
}

func (a *App) updateAudio() {
	phaseChanged := a.model.Phase != a.lastPhase
	sceneChanged := a.model.Scene != a.lastScene
	if phaseChanged || sceneChanged || a.model.Eliminations > a.lastEliminations {
		a.ensureAudio()
	}
	if a.audio == nil {
		return
	}
	if phaseChanged && a.model.Phase == game.PhasePanic {
		forced := a.model.PanicRemaining < a.model.Config.PanicDuration
		_ = a.audio.play(effectForPanicTransition(forced))
	}
	if a.model.Eliminations > a.lastEliminations {
		_ = a.audio.play(sound.EffectElimination)
	}
	if phaseChanged && a.model.Phase == game.PhaseRecoverAvailable {
		_ = a.audio.play(sound.EffectRecover)
	}
	if sceneChanged {
		switch a.model.Scene {
		case game.SceneVictory:
			_ = a.audio.play(sound.EffectVictory)
		case game.SceneGameOver:
			_ = a.audio.play(sound.EffectGameOver)
		}
	}
}

func (a *App) ensureAudio() {
	if a.audioAttempted {
		return
	}
	a.audioAttempted = true
	system, err := newSoundSystem(a.sprites)
	if err == nil {
		a.audio = system
	}
}

func (a *App) spawnBurst(position game.Vec2, particleColor color.RGBA, count int, speed float64) {
	for i := 0; i < count; i++ {
		angle := float64(i) * 2 * math.Pi / float64(count)
		a.particles.Spawn(position, game.Vec2{X: math.Cos(angle) * speed, Y: math.Sin(angle) * speed}, 0.45, particleColor)
	}
}

func (a *App) drawStability(screen *ebiten.Image) {
	const barWidth = 90
	vector.DrawFilledRect(screen, 115, 7, barWidth, 6, colorGrid, false)
	fraction := float32(0)
	if a.model.Phase == game.PhasePanic || a.model.Phase == game.PhaseRecoverAvailable {
		fraction = float32(a.model.PanicRemaining.Seconds() / a.model.Config.PanicDuration.Seconds())
	}
	if fraction > 0 {
		vector.DrawFilledRect(screen, 115, 7, barWidth*fraction, 6, colorPanic, false)
	}
}

func (a *App) drawResult(screen *ebiten.Image) {
	resultColor := colorPanic
	if a.model.Scene == game.SceneVictory {
		resultColor = colorRecover
	}
	vector.DrawFilledRect(screen, 58, 65, 204, 49, colorBackground, false)
	text.Draw(screen, a.stateText, basicfont.Face7x13, 71, 93, resultColor)
}
