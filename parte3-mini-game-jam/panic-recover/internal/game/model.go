package game

import "time"

const (
	VirtualWidth  = 320
	VirtualHeight = 180
)

type Scene uint8

const (
	SceneTitle Scene = iota
	ScenePlaying
	SceneVictory
	SceneGameOver
)

type Phase uint8

const (
	PhaseCalm Phase = iota
	PhasePanic
	PhaseRecoverAvailable
)

type Vec2 struct {
	X float64
	Y float64
}

type Player struct {
	Position Vec2
	Velocity Vec2
	Radius   float64
}

type Bug struct {
	Position Vec2
	Radius   float64
	Speed    float64
	Alive    bool
}

type RecoverZone struct {
	Position Vec2
	Radius   float64
	Active   bool
}

type CycleSpec struct {
	BugCount int
	Quota    int
	BugSpeed float64
}

type Config struct {
	CalmDuration        time.Duration
	PanicDuration       time.Duration
	ForcedPanicFraction float64
	Cycles              []CycleSpec
}

type Model struct {
	Scene          Scene
	Phase          Phase
	Config         Config
	Cycle          int
	CalmRemaining  time.Duration
	PanicRemaining time.Duration
	Eliminations   int
	Player         Player
	Bugs           []Bug
	Recover        RecoverZone
}

func NewModel() *Model {
	return &Model{
		Scene: SceneTitle,
		Config: Config{
			CalmDuration:        5 * time.Second,
			PanicDuration:       12 * time.Second,
			ForcedPanicFraction: 0.7,
			Cycles: []CycleSpec{
				{BugCount: 5, Quota: 3, BugSpeed: 22},
				{BugCount: 8, Quota: 5, BugSpeed: 27},
				{BugCount: 11, Quota: 7, BugSpeed: 32},
			},
		},
	}
}

func (m *Model) Start() {
	m.Scene = ScenePlaying
	m.Cycle = 0
	m.startCycle()
}

func (m *Model) startCycle() {
	spec := m.Config.Cycles[m.Cycle]
	m.Phase = PhaseCalm
	m.CalmRemaining = m.Config.CalmDuration
	m.PanicRemaining = m.Config.PanicDuration
	m.Eliminations = 0
	m.Player = Player{
		Position: Vec2{X: VirtualWidth / 2, Y: VirtualHeight / 2},
		Radius:   7,
	}
	m.Recover = RecoverZone{Radius: 11}
	m.Bugs = make([]Bug, spec.BugCount)
	positions := [...]Vec2{
		{X: 28, Y: 28}, {X: 160, Y: 24}, {X: 292, Y: 28},
		{X: 292, Y: 90}, {X: 292, Y: 152}, {X: 160, Y: 156},
		{X: 28, Y: 152}, {X: 28, Y: 90}, {X: 92, Y: 36},
		{X: 228, Y: 144}, {X: 228, Y: 36},
	}
	for i := range m.Bugs {
		m.Bugs[i] = Bug{
			Position: positions[i],
			Radius:   6,
			Speed:    spec.BugSpeed,
			Alive:    true,
		}
	}
}
