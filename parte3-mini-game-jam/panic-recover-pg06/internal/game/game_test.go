package game

import (
	"testing"
)

func TestGameInitialization(t *testing.T) {
	g := NewGame()
	if g.State != StateTitle {
		t.Errorf("expected initial state StateTitle, got %v", g.State)
	}

	w, h := g.Layout(1920, 1080)
	if w != VirtualWidth || h != VirtualHeight {
		t.Errorf("expected layout (%d, %d), got (%d, %d)", VirtualWidth, VirtualHeight, w, h)
	}
}

func TestGameStartLevelAndRecover(t *testing.T) {
	g := NewGame()
	g.StartLevel(0)

	if g.State != StatePlaying {
		t.Fatalf("expected state StatePlaying, got %v", g.State)
	}
	if g.CurrentLevelIdx != 0 {
		t.Errorf("expected level index 0, got %d", g.CurrentLevelIdx)
	}
	if len(g.WorkingLines) != len(g.CurrentLevel.CodeLines) {
		t.Errorf("expected %d working lines, got %d", len(g.CurrentLevel.CodeLines), len(g.WorkingLines))
	}

	// Test submission of a correct fix
	g.editor.SelectedLineIndex = 1
	g.evaluateSubmission("gopher := &Gopher{}")

	if g.State != StateLevelRecovered {
		t.Errorf("expected state StateLevelRecovered, got %v", g.State)
	}
	if g.Score <= 0 {
		t.Errorf("expected positive score bonus, got %d", g.Score)
	}
}

func TestGamePanicTrigger(t *testing.T) {
	g := NewGame()
	g.StartLevel(0)

	g.triggerPanic()
	if g.State != StatePanicCrash {
		t.Errorf("expected state StatePanicCrash, got %v", g.State)
	}
}
