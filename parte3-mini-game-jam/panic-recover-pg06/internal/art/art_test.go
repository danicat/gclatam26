package art

import (
	"testing"
)

func TestHighlightGoLine(t *testing.T) {
	line := `func main() { var x int = 42; println("Hello") }`
	tokens := HighlightGoLine(line)
	if len(tokens) == 0 {
		t.Fatalf("expected non-empty tokens")
	}

	hasFunc := false
	hasString := false
	for _, tok := range tokens {
		if tok.Text == "func" && tok.Color == ColorCyanGlow {
			hasFunc = true
		}
		if tok.Text == `"Hello"` && tok.Color == ColorGreenRecover {
			hasString = true
		}
	}

	if !hasFunc {
		t.Errorf("expected 'func' keyword highlighted in cyan")
	}
	if !hasString {
		t.Errorf("expected '\"Hello\"' string highlighted in green")
	}
}

func TestParticleSystem(t *testing.T) {
	ps := NewParticleSystem(100)
	ps.SpawnSparks(50, 50, 20, ColorGreenRecover)
	if len(ps.pool) != 20 {
		t.Errorf("expected 20 particles, got %d", len(ps.pool))
	}
	// Advance time
	ps.Update(0.1)
	if len(ps.pool) != 20 {
		t.Errorf("expected particles to still be alive after 0.1s, got %d", len(ps.pool))
	}
	// Advance large amount of time to expire particles
	ps.Update(2.0)
	if len(ps.pool) != 0 {
		t.Errorf("expected particles to expire after 2.0s, got %d", len(ps.pool))
	}
}
