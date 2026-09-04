package entities

import (
	"testing"

	"gopher-budokai-panic/internal/gfx"
)

func TestPanicStateTransitions(t *testing.T) {
	var ps PanicState
	if ps.IsPanicked {
		t.Fatalf("expected not panicked initially")
	}

	// Add panic below threshold
	ps.AddPanic(40.0)
	if ps.IsPanicked || ps.Meter != 40.0 {
		t.Fatalf("expected 40 meter and not panicked, got meter=%.1f panicked=%v", ps.Meter, ps.IsPanicked)
	}

	// Cross 100 threshold into PANIC!
	ps.AddPanic(70.0)
	if !ps.IsPanicked || ps.Meter != 100.0 {
		t.Fatalf("expected panicked with 100 meter")
	}

	// Test mash recovery effort build
	recovered := ps.TryMashRecover() // 28%
	if recovered {
		t.Fatalf("expected not recovered after single mash")
	}
	ps.TryMashRecover() // 56%
	ps.TryMashRecover() // 84%
	recovered = ps.TryMashRecover() // 112% -> recovered!
	if !recovered || ps.IsPanicked || ps.Meter != 0 {
		t.Fatalf("expected recovery to succeed, got recovered=%v panicked=%v", recovered, ps.IsPanicked)
	}
}

func TestFighterCombat(t *testing.T) {
	p := NewFighter(0, gfx.FighterPlayer, 100, 100, false)
	cpu := NewFighter(1, gfx.FighterCPU, 200, 100, true)

	// Test damage application
	p.TakeDamage(100.0, false, cpu.X)
	if p.Health != 900.0 {
		t.Fatalf("expected 900 health, got %.1f", p.Health)
	}
	if p.Panic.Meter <= 0 {
		t.Fatalf("expected panic meter increase on damage")
	}

	// Test beam firing
	p.State = StateIdle
	p.Ki = 50.0
	beam := p.StartSuperBeam()
	if beam == nil {
		t.Fatalf("expected super beam to fire")
	}
	if p.State != StateBeam {
		t.Fatalf("expected fighter state to be StateBeam")
	}
}
