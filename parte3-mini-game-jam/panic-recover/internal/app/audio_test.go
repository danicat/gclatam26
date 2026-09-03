package app

import (
	"testing"

	"panic-recover/internal/game"
	"panic-recover/internal/sound"
)

func TestEffectForPanicTransitionDistinguishesForcedEntry(t *testing.T) {
	t.Parallel()

	if got := effectForPanicTransition(false); got != sound.EffectPanic {
		t.Fatalf("normal transition = %q, want %q", got, sound.EffectPanic)
	}
	if got := effectForPanicTransition(true); got != sound.EffectForcedPanic {
		t.Fatalf("forced transition = %q, want %q", got, sound.EffectForcedPanic)
	}
	if got := effectForPhase(game.PhaseRecoverAvailable); got != sound.EffectRecover {
		t.Fatalf("recover phase = %q, want %q", got, sound.EffectRecover)
	}
}
