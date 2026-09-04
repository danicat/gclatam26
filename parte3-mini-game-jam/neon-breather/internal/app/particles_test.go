package app

import (
	"image/color"
	"testing"

	"panic-recover/internal/game"
)

func TestParticleSystemUsesFixedPoolAndExpiresParticles(t *testing.T) {
	t.Parallel()

	particles := NewParticleSystem(1)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	particles.Spawn(game.Vec2{X: 10, Y: 10}, game.Vec2{X: 4, Y: 0}, 0.5, white)
	particles.Spawn(game.Vec2{X: 20, Y: 10}, game.Vec2{}, 1, white)

	if got := particles.ActiveCount(); got != 1 {
		t.Fatalf("ActiveCount after full pool = %d, want 1", got)
	}
	particles.Update(0.5)
	if got := particles.ActiveCount(); got != 0 {
		t.Fatalf("ActiveCount after expiration = %d, want 0", got)
	}
}
