package game

import (
	"exploding-kitchens/internal/system"
	"exploding-kitchens/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// SceneID identifies a specific game scene.
type SceneID int

const (
	SceneKeepCurrent SceneID = iota
	SceneTitle
	ScenePlay
	SceneGameOver
)

// Scene defines the strict lifecycle hooks required for scene state management.
type Scene interface {
	Enter(ctx *GameContext)
	Update(dt float64, in system.InputState) (SceneID, error)
	Draw(screen *ebiten.Image)
	Exit()
}

// GameContext holds shared dependencies and cross-scene statistics.
type GameContext struct {
	Audio    *system.AudioManager
	Font     *ui.PixelFont
	PixelImg *ebiten.Image

	HighScore      int
	LastScore      int
	LastClutches   int
	LastExplosions int
	LastSurvived   bool
	DemoMode       bool
}
