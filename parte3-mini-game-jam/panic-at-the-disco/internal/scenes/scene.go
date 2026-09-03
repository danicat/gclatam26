package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"panic-at-the-disco/internal/levels"
)

type SceneID int

const (
	SceneNone SceneID = iota
	SceneTitle
	ScenePlay
	SceneClear
	SceneGameOver
	SceneVictory
)

type ActionType int

const (
	ActionNone ActionType = iota
	ActionSwitchScene
)

type SceneAction struct {
	Type        ActionType
	NextScene   SceneID
	TargetZone  levels.ZoneID
	Score       int
	Lives       int
	LossReason  string
	SurviveTime float64
}

// Scene defines the strict lifecycle hooks for game states.
type Scene interface {
	Enter()
	Update(dt float64) SceneAction
	Draw(screen *ebiten.Image)
	Exit()
}
