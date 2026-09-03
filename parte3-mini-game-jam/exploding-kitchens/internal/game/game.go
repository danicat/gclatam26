package game

import (
	"image/color"

	"exploding-kitchens/internal/system"
	"exploding-kitchens/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	VirtualWidth  = 320
	VirtualHeight = 180
	TargetFPS     = 60
)

// Game implements the ebiten.Game interface with scene state management and virtual resolution scaling.
type Game struct {
	context      *GameContext
	inputManager *system.InputManager
	scenes       map[SceneID]Scene
	currentID    SceneID
	currentScene Scene
}

// NewGame initializes the engine, audio, fonts, input system, and scene map.
func NewGame(titleScene, playScene, gameOverScene Scene) *Game {
	// Pre-allocate 1x1 white pixel image for fast procedural drawing
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)

	font := ui.NewPixelFont()
	audioMgr := system.NewAudioManager()

	ctx := &GameContext{
		Audio:    audioMgr,
		Font:     font,
		PixelImg: pixel,
	}

	scenes := map[SceneID]Scene{
		SceneTitle:    titleScene,
		ScenePlay:     playScene,
		SceneGameOver: gameOverScene,
	}

	g := &Game{
		context:      ctx,
		inputManager: system.NewInputManager(),
		scenes:       scenes,
		currentID:    SceneTitle,
		currentScene: scenes[SceneTitle],
	}

	g.currentScene.Enter(g.context)
	return g
}

// Update handles input polling, fullscreen toggles, delta-time progression, and scene transitions.
func (g *Game) Update() error {
	in := g.inputManager.Poll()

	// Fullscreen toggle (F11 or Alt+Enter)
	if in.ToggleFull {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Fixed delta time for 60 TPS cycle timing
	const dt = 1.0 / float64(TargetFPS)

	nextID, err := g.currentScene.Update(dt, in)
	if err != nil {
		return err
	}

	// Handle Scene Transitions
	if nextID != SceneKeepCurrent && nextID != g.currentID {
		if nextScene, exists := g.scenes[nextID]; exists {
			g.currentScene.Exit()
			g.currentID = nextID
			g.currentScene = nextScene
			g.currentScene.Enter(g.context)
		}
	}

	return nil
}

// Draw renders the active scene onto the virtual 320x180 resolution canvas with zero allocations.
func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

// Layout dictates the fixed internal virtual pixel resolution regardless of outside window size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}
