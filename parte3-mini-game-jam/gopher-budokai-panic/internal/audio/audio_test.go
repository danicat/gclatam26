package audio

import (
	"testing"
)

func TestAudioSynthesis(t *testing.T) {
	// Test single note synthesis
	note := NoteDef{
		WaveType:  "sawtooth",
		Duration:  0.05,
		StartFreq: 440.0,
		EndFreq:   440.0,
		Volume:    0.5,
		Attack:    0.01,
		Decay:     0.01,
		Sustain:   0.8,
		Release:   0.01,
	}
	pcm := SynthesizeNote(note)
	if len(pcm) == 0 {
		t.Fatalf("expected non-empty PCM buffer")
	}

	// Test track mixing with clamp
	track1 := TrackDef{
		Notes: []NoteDef{note},
	}
	track2 := TrackDef{
		Notes: []NoteDef{note},
	}
	mixed := MixTracks([]TrackDef{track1, track2}, 1.0)
	if len(mixed) == 0 {
		t.Fatalf("expected non-empty mixed PCM")
	}
}

func TestSFXGeneration(t *testing.T) {
	if len(BuildBlastSFX()) == 0 {
		t.Errorf("blast sfx empty")
	}
	if len(BuildHitSFX()) == 0 {
		t.Errorf("hit sfx empty")
	}
	if len(BuildChargeSFX()) == 0 {
		t.Errorf("charge sfx empty")
	}
	if len(BuildBeamSFX()) == 0 {
		t.Errorf("beam sfx empty")
	}
	if len(BuildVanishSFX()) == 0 {
		t.Errorf("vanish sfx empty")
	}
	if len(BuildPanicAlertSFX()) == 0 {
		t.Errorf("panic alert sfx empty")
	}
	if len(BuildRecoverKiaiSFX()) == 0 {
		t.Errorf("recover kiai sfx empty")
	}
}
