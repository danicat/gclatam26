package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"gopher-budokai-panic/internal/scenes"
)

const (
	VirtualWidth  = 640
	VirtualHeight = 360
)

type Game struct {
	currentScene scenes.Scene
	lastTime     time.Time
}

func NewGame() *Game {
	g := &Game{
		lastTime: time.Now(),
	}
	g.SwitchScene(scenes.NewTitleScene())
	return g
}

func (g *Game) SwitchScene(next scenes.Scene) {
	if g.currentScene != nil {
		g.currentScene.Exit()
	}
	g.currentScene = next
	if g.currentScene != nil {
		g.currentScene.Enter()
	}
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	// Cap delta time to prevent spiraling after pauses
	if dt > 0.05 {
		dt = 0.05
	}

	// Fullscreen toggle (F11)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if g.currentScene != nil {
		next := g.currentScene.Update(dt)
		if next != nil {
			g.SwitchScene(next)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return VirtualWidth, VirtualHeight
}
