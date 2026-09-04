package gfx

import (
	"image/color"
	"testing"
)

func TestParticlePool(t *testing.T) {
	ps := NewParticleSystem(10)

	// Verify all inactive
	for i, p := range ps.pool {
		if p.Active {
			t.Fatalf("Particle %d should be inactive initially", i)
		}
	}

	// Emit 5 particles
	for i := 0; i < 5; i++ {
		ps.Emit(10.0, 10.0, 50.0, 50.0, 0.5, 4.0, 1.0, color.RGBA{255, 255, 255, 255}, ParticleSpark)
	}

	activeCount := 0
	for _, p := range ps.pool {
		if p.Active {
			activeCount++
		}
	}
	if activeCount != 5 {
		t.Fatalf("Expected 5 active particles, got %d", activeCount)
	}

	// Update past life duration (0.5s)
	ps.Update(0.6)

	activeCountAfter := 0
	for _, p := range ps.pool {
		if p.Active {
			activeCountAfter++
		}
	}
	if activeCountAfter != 0 {
		t.Fatalf("Expected 0 active particles after expiration, got %d", activeCountAfter)
	}
}

func TestParticlePoolReset(t *testing.T) {
	ps := NewParticleSystem(10)
	ps.Emit(0, 0, 0, 0, 1.0, 1.0, 1.0, color.RGBA{}, ParticleSpark)
	ps.Reset()

	for i, p := range ps.pool {
		if p.Active {
			t.Fatalf("Particle %d still active after Reset()", i)
		}
	}
}
