package entities

import (
	"testing"

	"panic-at-the-disco/internal/input"
)

func TestPlayerCreationAndBounds(t *testing.T) {
	p := NewPlayer(100.0, 150.0)

	if p.X != 100.0 || p.Y != 150.0 {
		t.Fatalf("Expected position (100, 150), got (%f, %f)", p.X, p.Y)
	}
	if p.Lives != 3 {
		t.Fatalf("Expected 3 lives, got %d", p.Lives)
	}

	bx, by, bw, bh := p.Bounds()
	if bw <= 0 || bh <= 0 {
		t.Fatalf("Invalid bounds (%f, %f, %f, %f)", bx, by, bw, bh)
	}
}

func TestPlayerBoundaryClamping(t *testing.T) {
	p := NewPlayer(50.0, 50.0)

	// Move hard left past fieldX
	in := input.InputState{MoveX: -1.0, MoveY: 0.0}
	fieldX, fieldY, fieldW, fieldH := 40.0, 40.0, 500.0, 200.0

	// Step multiple frames
	for i := 0; i < 60; i++ {
		p.Update(1.0/60.0, in, fieldX, fieldY, fieldW, fieldH, nil, nil)
	}

	minAllowedX := fieldX + p.Width/2
	if p.X < minAllowedX-0.01 {
		t.Fatalf("Player breached left boundary: got %f, min %f", p.X, minAllowedX)
	}
}

func TestPlayerDashActivation(t *testing.T) {
	p := NewPlayer(100.0, 100.0)
	p.GrooveMeter = 50.0

	in := input.InputState{DashJustDown: true, MoveX: 1.0}
	fieldX, fieldY, fieldW, fieldH := 0.0, 0.0, 640.0, 360.0

	// Dash should activate
	p.Update(1.0/60.0, in, fieldX, fieldY, fieldW, fieldH, nil, nil)

	if !p.IsDashing {
		t.Fatal("Expected player to be dashing")
	}
	if p.GrooveMeter >= 50.0 {
		t.Fatal("Expected groove meter to be consumed by dash")
	}
	if p.InvulnerableTimer <= 0 {
		t.Fatal("Expected invulnerability during dash")
	}
}

func TestPlayerDamageAndInvulnerability(t *testing.T) {
	p := NewPlayer(100.0, 100.0)
	p.Lives = 3

	// Take damage
	damaged := p.TakeDamage(1, nil, nil)
	if !damaged {
		t.Fatal("Expected damage to be applied")
	}
	if p.Lives != 2 {
		t.Fatalf("Expected 2 lives remaining, got %d", p.Lives)
	}
	if p.InvulnerableTimer <= 0 {
		t.Fatal("Expected invulnerability timer to be active")
	}

	// Immediate second damage attempt should be ignored during invulnerability
	secondDamage := p.TakeDamage(1, nil, nil)
	if secondDamage {
		t.Fatal("Player should be immune to damage while invulnerable")
	}
	if p.Lives != 2 {
		t.Fatalf("Expected lives to stay at 2, got %d", p.Lives)
	}
}

func TestPlayerSlipPhysics(t *testing.T) {
	p := NewPlayer(100.0, 100.0)
	p.ApplySlip(0.5, nil)

	if !p.IsSlipping {
		t.Fatal("Expected player to be slipping")
	}
	if p.SlipTimer <= 0 {
		t.Fatal("Expected positive slip timer")
	}
}
