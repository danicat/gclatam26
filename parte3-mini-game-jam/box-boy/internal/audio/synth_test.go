package audio_test

import (
	"testing"

	"box-boy/internal/audio"
)

func TestGenerateNotePCM(t *testing.T) {
	note := audio.NoteDef{
		WaveType:  "sine",
		Duration:  0.05,
		StartFreq: 440.0,
		EndFreq:   440.0,
		Volume:    0.2,
		Attack:    0.01,
		Decay:     0.01,
		Sustain:   0.8,
		Release:   0.01,
		Pan:       0.0,
	}

	pcm := audio.GenerateNotePCM(note)
	if len(pcm) == 0 {
		t.Fatalf("buffer PCM não deveria ser vazio")
	}

	// 44100 amostras/segundo * 0.05s * 4 bytes/amostra (16 bits stereo) = 8820 bytes
	expectedSamples := int(0.05 * 44100)
	expectedBytes := expectedSamples * 4
	if len(pcm) != expectedBytes {
		t.Errorf("tamanho do buffer PCM incorreto: esperava %d bytes, obteve %d", expectedBytes, len(pcm))
	}
}

func TestMixTracksPCM(t *testing.T) {
	n1 := audio.GenerateNotePCM(audio.NoteDef{
		WaveType: "square", Duration: 0.02, StartFreq: 220, EndFreq: 220, Volume: 0.1,
	})
	n2 := audio.GenerateNotePCM(audio.NoteDef{
		WaveType: "triangle", Duration: 0.02, StartFreq: 330, EndFreq: 330, Volume: 0.1,
	})

	mixed := audio.MixTracksPCM([][]byte{n1, n2})
	if len(mixed) != len(n1) {
		t.Errorf("mixagem deveria manter comprimento das faixas: esperava %d, obteve %d", len(n1), len(mixed))
	}
}
