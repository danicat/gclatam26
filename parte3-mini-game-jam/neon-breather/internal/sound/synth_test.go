package sound

import (
	"encoding/binary"
	"testing"
)

func TestGenerateEffectProducesStereoPCMWithoutClipping(t *testing.T) {
	t.Parallel()

	pcm := GenerateEffect(EffectPanic)
	if len(pcm) == 0 {
		t.Fatal("GenerateEffect returned empty PCM")
	}
	if len(pcm)%4 != 0 {
		t.Fatalf("PCM length = %d, want a multiple of 4 bytes", len(pcm))
	}
	for offset := 0; offset < len(pcm); offset += 2 {
		sample := int16(binary.LittleEndian.Uint16(pcm[offset : offset+2]))
		if sample < -32768 || sample > 32767 {
			t.Fatalf("sample at offset %d = %d, outside int16 range", offset, sample)
		}
	}
}

func TestAllGameplayEffectsHaveAudio(t *testing.T) {
	t.Parallel()

	for _, effect := range []Effect{
		EffectPanic,
		EffectForcedPanic,
		EffectElimination,
		EffectRecover,
		EffectCritical,
		EffectVictory,
		EffectGameOver,
	} {
		if pcm := GenerateEffect(effect); len(pcm) == 0 {
			t.Errorf("GenerateEffect(%q) returned empty PCM", effect)
		}
	}
}
