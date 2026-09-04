package system

import (
	"testing"
)

func TestParticlePoolLifecycle(t *testing.T) {
	pool := NewParticlePool(10)

	// Emit steam (3 particles)
	pool.EmitSteam(50, 50, 3)

	activeCount := 0
	for _, p := range pool.particles {
		if p.Active {
			activeCount++
		}
	}
	if activeCount != 3 {
		t.Fatalf("expected 3 active particles, got %d", activeCount)
	}

	// Update past lifespan
	pool.Update(2.0)

	activeCountAfter := 0
	for _, p := range pool.particles {
		if p.Active {
			activeCountAfter++
		}
	}
	if activeCountAfter != 0 {
		t.Fatalf("expected 0 active particles after expiration, got %d", activeCountAfter)
	}
}
