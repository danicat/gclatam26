package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"panic-at-the-disco/internal/scenes"
)

const (
	VirtualWidth  = 640
	VirtualHeight = 360
)

type Game struct {
	currentScene scenes.Scene
}

func NewGame() *Game {
	g := &Game{}
	g.currentScene = scenes.NewTitleScene()
	g.currentScene.Enter()
	return g
}

func (g *Game) Update() error {
	// Fullscreen toggle hotkey (F11 or Alt+Enter) per ebitengineer standard
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	const dt = 1.0 / 60.0

	// Update active scene
	action := g.currentScene.Update(dt)
	if action.Type == scenes.ActionSwitchScene {
		g.switchScene(action)
	}

	return nil
}

func (g *Game) switchScene(action scenes.SceneAction) {
	if g.currentScene != nil {
		g.currentScene.Exit()
	}

	switch action.NextScene {
	case scenes.SceneTitle:
		g.currentScene = scenes.NewTitleScene()
	case scenes.ScenePlay:
		g.currentScene = scenes.NewPlayScene(action.TargetZone, action.Lives, action.Score)
	case scenes.SceneClear:
		g.currentScene = scenes.NewStageClearScene(action.TargetZone, action.Score, action.Lives)
	case scenes.SceneGameOver:
		g.currentScene = scenes.NewGameOverScene(action.LossReason, action.Score, action.SurviveTime)
	case scenes.SceneVictory:
		g.currentScene = scenes.NewVictoryScene(action.Score, action.Lives, action.SurviveTime)
	}

	if g.currentScene != nil {
		g.currentScene.Enter()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return VirtualWidth, VirtualHeight
}
