package audio

import (
	"testing"
)

func TestSynthesizedAudioBuffers(t *testing.T) {
	keyClick := synthKeyClick()
	if len(keyClick) == 0 {
		t.Errorf("synthKeyClick produced empty buffer")
	}

	tone := synthTone(440, 0.05, "sine", 0.5)
	if len(tone) == 0 {
		t.Errorf("synthTone produced empty buffer")
	}

	doubleTone := synthDoubleTone(440, 880, 0.05, 0.5)
	if len(doubleTone) == 0 {
		t.Errorf("synthDoubleTone produced empty buffer")
	}

	siren := synthSquareSiren(880, 440, 0.1, 0.5)
	if len(siren) == 0 {
		t.Errorf("synthSquareSiren produced empty buffer")
	}

	fanfare := synthRecoverFanfare()
	if len(fanfare) == 0 {
		t.Errorf("synthRecoverFanfare produced empty buffer")
	}

	crash := synthPanicCrash()
	if len(crash) == 0 {
		t.Errorf("synthPanicCrash produced empty buffer")
	}

	bgm := synthCyberBGM()
	if len(bgm) == 0 {
		t.Errorf("synthCyberBGM produced empty buffer")
	}
	// Check that buffer length corresponds to reasonable duration
	// 44100 samples/sec * 2 channels * 2 bytes/sample = 176400 bytes/sec
	bytesPerSec := 44100 * 2 * 2
	expectedMinBytes := int(float64(bytesPerSec) * 2.0)
	if len(bgm) < expectedMinBytes {
		t.Errorf("bgm buffer too short: got %d bytes, expected >= %d", len(bgm), expectedMinBytes)
	}
}
