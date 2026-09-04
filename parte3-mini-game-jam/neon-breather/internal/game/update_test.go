package game

import (
	"math"
	"testing"
	"time"
)

func TestUpdateActivatesFullPanicOnCommand(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{PanicPressed: true}, time.Second/60)

	if m.Phase != PhasePanic {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhasePanic)
	}
	if got, want := m.PanicRemaining, 12*time.Second; got != want {
		t.Fatalf("PanicRemaining = %v, want %v", got, want)
	}
}

func TestUpdateActivatesFullPanicWhenCalmExpires(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{}, 5*time.Second)

	if m.Phase != PhasePanic {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhasePanic)
	}
	if got, want := m.PanicRemaining, 12*time.Second; got != want {
		t.Fatalf("PanicRemaining = %v, want %v", got, want)
	}
}

func TestUpdateNormalizesDiagonalMovement(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{Move: Vec2{X: 1, Y: 1}}, time.Second)

	wantVelocity := 45 / math.Sqrt2
	if math.Abs(m.Player.Velocity.X-wantVelocity) > 0.001 {
		t.Fatalf("Velocity.X = %f, want %f", m.Player.Velocity.X, wantVelocity)
	}
	if math.Abs(m.Player.Velocity.Y-wantVelocity) > 0.001 {
		t.Fatalf("Velocity.Y = %f, want %f", m.Player.Velocity.Y, wantVelocity)
	}
	if math.Abs(m.Player.Position.X-(VirtualWidth/2+wantVelocity)) > 0.001 {
		t.Fatalf("Position.X = %f, want %f", m.Player.Position.X, VirtualWidth/2+wantVelocity)
	}
}

func TestCalmCollisionForcesPanicAtSeventyPercent(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Bugs[0].Position = m.Player.Position
	m.Update(Input{}, time.Second/60)

	if m.Phase != PhasePanic {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhasePanic)
	}
	if got, want := m.PanicRemaining, 8400*time.Millisecond; got != want {
		t.Fatalf("PanicRemaining = %v, want %v", got, want)
	}
}

func TestPanicTimeoutEndsTheGame(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{PanicPressed: true}, time.Second/60)
	m.Update(Input{}, 12*time.Second)

	if m.Scene != SceneGameOver {
		t.Fatalf("Scene = %v, want %v", m.Scene, SceneGameOver)
	}
	if m.PanicRemaining != 0 {
		t.Fatalf("PanicRemaining = %v, want 0", m.PanicRemaining)
	}
}

func TestPanicEliminationsUnlockRecover(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{PanicPressed: true}, time.Second/60)
	for i := 0; i < 3; i++ {
		m.Bugs[i].Position = m.Player.Position
	}
	m.Update(Input{}, time.Millisecond)

	if m.Eliminations != 3 {
		t.Fatalf("Eliminations = %d, want 3", m.Eliminations)
	}
	if m.Phase != PhaseRecoverAvailable {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhaseRecoverAvailable)
	}
	if !m.Recover.Active {
		t.Fatal("Recover.Active = false, want true")
	}
	if overlaps(m.Player.Position, m.Player.Radius, m.Recover.Position, m.Recover.Radius) {
		t.Fatal("Recover overlaps player")
	}
	for _, bug := range m.Bugs {
		if bug.Alive && overlaps(bug.Position, bug.Radius, m.Recover.Position, m.Recover.Radius) {
			t.Fatal("Recover overlaps a live bug")
		}
	}
}

func TestTouchingRecoverStartsNextCycle(t *testing.T) {
	t.Parallel()

	m := modelAtRecover(t, 0)
	m.Update(Input{}, time.Millisecond)

	if m.Cycle != 1 {
		t.Fatalf("Cycle = %d, want 1", m.Cycle)
	}
	if m.Phase != PhaseCalm {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhaseCalm)
	}
	if len(m.Bugs) != 8 {
		t.Fatalf("len(Bugs) = %d, want 8", len(m.Bugs))
	}
}

func TestTouchingThirdRecoverWins(t *testing.T) {
	t.Parallel()

	m := modelAtRecover(t, 2)
	m.Update(Input{}, time.Millisecond)

	if m.Scene != SceneVictory {
		t.Fatalf("Scene = %v, want %v", m.Scene, SceneVictory)
	}
}

func TestBugSeeksPlayer(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Player.Position = Vec2{X: 160, Y: 90}
	m.Bugs[0].Position = Vec2{X: 20, Y: 90}
	startDistance := distanceSquared(m.Bugs[0].Position, m.Player.Position)
	m.Update(Input{}, 100*time.Millisecond)

	if got := distanceSquared(m.Bugs[0].Position, m.Player.Position); got >= startDistance {
		t.Fatalf("distance after update = %f, want less than %f", got, startDistance)
	}
}

func TestPanicMovementPreservesInputDirection(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()
	m.Update(Input{PanicPressed: true}, time.Second/60)
	m.Player.Velocity = Vec2{}
	m.Update(Input{Move: Vec2{X: -1}}, time.Second/60)

	if m.Player.Velocity.X >= 0 {
		t.Fatalf("Velocity.X = %f, want negative for left input", m.Player.Velocity.X)
	}
}

func modelAtRecover(t *testing.T, cycle int) *Model {
	t.Helper()
	m := NewModel()
	m.Start()
	m.Cycle = cycle
	m.startCycle()
	m.Phase = PhaseRecoverAvailable
	m.Recover.Active = true
	m.Recover.Position = m.Player.Position
	return m
}
