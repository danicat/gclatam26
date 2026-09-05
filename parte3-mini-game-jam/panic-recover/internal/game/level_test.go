package game

import (
	"testing"
)

func TestLevelInitAndMovement(t *testing.T) {
	def := LevelDef{
		Name:              "Test Room",
		Width:             5,
		Height:            5,
		MaxTurns:          6,
		ClockRecoverTurns: 4,
		PlayerStartX:      1,
		PlayerStartY:      1,
		Layout: []string{
			"#####",
			"#P..#",
			"#...#",
			"#..A#",
			"#####",
		},
	}

	state := NewLevelState(def)
	if state.PlayerX != 1 || state.PlayerY != 1 {
		t.Fatalf("expected player at (1,1), got (%d,%d)", state.PlayerX, state.PlayerY)
	}
	if state.TurnsLeft != 6 {
		t.Fatalf("expected 6 turns, got %d", state.TurnsLeft)
	}

	// Move right
	moved := state.Move(1, 0)
	if !moved {
		t.Fatal("expected move to succeed")
	}
	if state.PlayerX != 2 || state.PlayerY != 1 {
		t.Fatalf("expected player at (2,1), got (%d,%d)", state.PlayerX, state.PlayerY)
	}
	if state.TurnsLeft != 5 {
		t.Fatalf("expected 5 turns remaining, got %d", state.TurnsLeft)
	}
}

func TestBoulderIntoHole(t *testing.T) {
	def := LevelDef{
		Name:              "Push Test",
		Width:             6,
		Height:            3,
		MaxTurns:          10,
		ClockRecoverTurns: 4,
		PlayerStartX:      1,
		PlayerStartY:      1,
		Layout: []string{
			"######",
			"#PBOA#",
			"######",
		},
	}

	state := NewLevelState(def)
	if len(state.Boulders) != 1 {
		t.Fatalf("expected 1 boulder, got %d", len(state.Boulders))
	}

	// Push boulder into hole at (3,1)
	moved := state.Move(1, 0)
	if !moved {
		t.Fatal("expected push to succeed")
	}

	// Boulder should fall into hole, turning tile into TileHoleFilled and boulder removed
	if len(state.Boulders) != 0 {
		t.Fatalf("expected 0 boulders remaining, got %d", len(state.Boulders))
	}
	if state.Tiles[1][3] != TileHoleFilled {
		t.Fatalf("expected tile (3,1) to be TileHoleFilled, got %v", state.Tiles[1][3])
	}
	if state.PlayerX != 2 || state.PlayerY != 1 {
		t.Fatalf("expected player at (2,1), got (%d,%d)", state.PlayerX, state.PlayerY)
	}
}

func TestPanicThresholdAndRecovery(t *testing.T) {
	def := LevelDef{
		Name:              "Panic Test",
		Width:             6,
		Height:            3,
		MaxTurns:          5,
		ClockRecoverTurns: 4,
		PlayerStartX:      1,
		PlayerStartY:      1,
		Layout: []string{
			"######",
			"#P.C.#",
			"######",
		},
	}

	state := NewLevelState(def)
	// Step 1: turns left = 4, panic = 20%
	state.Move(1, 0)
	if state.PanicPercent() < 0.19 || state.PanicPercent() > 0.21 {
		t.Fatalf("expected ~20%% panic, got %f", state.PanicPercent())
	}

	// Step 2: onto clock tile (3, 1)
	// Turns was 4, recovers +4 -> 5 (clamped to MaxTurns), then consumes 1 turn -> 4!
	state.Move(1, 0)
	if state.TurnsLeft != 4 {
		t.Fatalf("expected 4 turns after recovery pickup, got %d", state.TurnsLeft)
	}
}

func TestPanicFaintAtZeroTurns(t *testing.T) {
	def := LevelDef{
		Name:              "Faint Test",
		Width:             5,
		Height:            3,
		MaxTurns:          1,
		ClockRecoverTurns: 0,
		PlayerStartX:      1,
		PlayerStartY:      1,
		Layout: []string{
			"#####",
			"#P..#",
			"#####",
		},
	}

	state := NewLevelState(def)
	state.Move(1, 0)
	if !state.Fainted {
		t.Fatal("expected player to faint when turns reach 0")
	}
	if state.TurnsLeft != 0 {
		t.Fatalf("expected 0 turns left, got %d", state.TurnsLeft)
	}
}
