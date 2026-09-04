package game

import (
	"testing"
	"time"
)

func TestNewModelDefinesJamRules(t *testing.T) {
	t.Parallel()

	m := NewModel()

	if m.Scene != SceneTitle {
		t.Fatalf("Scene = %v, want %v", m.Scene, SceneTitle)
	}
	if got, want := m.Config.CalmDuration, 5*time.Second; got != want {
		t.Fatalf("CalmDuration = %v, want %v", got, want)
	}
	if got, want := m.Config.PanicDuration, 12*time.Second; got != want {
		t.Fatalf("PanicDuration = %v, want %v", got, want)
	}
	if got, want := m.Config.ForcedPanicFraction, 0.7; got != want {
		t.Fatalf("ForcedPanicFraction = %v, want %v", got, want)
	}

	wantCycles := []CycleSpec{
		{BugCount: 5, Quota: 3, BugSpeed: 22},
		{BugCount: 8, Quota: 5, BugSpeed: 27},
		{BugCount: 11, Quota: 7, BugSpeed: 32},
	}
	if len(m.Config.Cycles) != len(wantCycles) {
		t.Fatalf("len(Cycles) = %d, want %d", len(m.Config.Cycles), len(wantCycles))
	}
	for i, want := range wantCycles {
		if got := m.Config.Cycles[i]; got != want {
			t.Errorf("Cycles[%d] = %+v, want %+v", i, got, want)
		}
	}
}

func TestStartBeginsFirstCalmCycle(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.Start()

	if m.Scene != ScenePlaying {
		t.Fatalf("Scene = %v, want %v", m.Scene, ScenePlaying)
	}
	if m.Phase != PhaseCalm {
		t.Fatalf("Phase = %v, want %v", m.Phase, PhaseCalm)
	}
	if m.Cycle != 0 {
		t.Fatalf("Cycle = %d, want 0", m.Cycle)
	}
	if m.CalmRemaining != 5*time.Second {
		t.Fatalf("CalmRemaining = %v, want 5s", m.CalmRemaining)
	}
	if len(m.Bugs) != 5 {
		t.Fatalf("len(Bugs) = %d, want 5", len(m.Bugs))
	}
}
