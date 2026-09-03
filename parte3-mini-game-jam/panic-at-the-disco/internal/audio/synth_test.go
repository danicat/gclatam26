package audio

import (
	"testing"
)

func TestProceduralWaveformGenerators(t *testing.T) {
	// 1. Whistle SFX
	whistle := synthesizeWhistle()
	if len(whistle) == 0 {
		t.Fatal("Expected non-empty whistle buffer")
	}
	// PCM 16-bit stereo has 4 bytes per sample
	if len(whistle)%4 != 0 {
		t.Fatalf("Whistle buffer length %d not aligned to 4 bytes", len(whistle))
	}

	// 2. Crash SFX
	crash := synthesizeCrash()
	if len(crash) == 0 || len(crash)%4 != 0 {
		t.Fatal("Invalid crash buffer")
	}

	// 3. Dash SFX
	dash := synthesizeDash()
	if len(dash) == 0 || len(dash)%4 != 0 {
		t.Fatal("Invalid dash buffer")
	}

	// 4. Slip SFX
	slip := synthesizeSlip()
	if len(slip) == 0 || len(slip)%4 != 0 {
		t.Fatal("Invalid slip buffer")
	}

	// 5. Win SFX
	win := synthesizeWin()
	if len(win) == 0 || len(win)%4 != 0 {
		t.Fatal("Invalid win buffer")
	}

	// 6. Alarm SFX
	alarm := synthesizeAlarm()
	if len(alarm) == 0 || len(alarm)%4 != 0 {
		t.Fatal("Invalid alarm buffer")
	}
}

func TestBGMGeneration(t *testing.T) {
	// Test synthesizing normal 120-BPM track
	normalBGM := synthesizeDiscoBGM(120.0, false)
	if len(normalBGM) == 0 || len(normalBGM)%4 != 0 {
		t.Fatal("Invalid normal BGM buffer")
	}

	// Test synthesizing fast 144-BPM panic track
	fastBGM := synthesizeDiscoBGM(144.0, true)
	if len(fastBGM) == 0 || len(fastBGM)%4 != 0 {
		t.Fatal("Invalid fast panic BGM buffer")
	}

	// Normal BGM at 120 BPM (16 beats = 8.0s) should be longer in byte count than 144 BPM (16 beats = 6.66s)
	if len(normalBGM) <= len(fastBGM) {
		t.Fatalf("Expected 120 BPM buffer (%d bytes) to be longer than 144 BPM buffer (%d bytes)", len(normalBGM), len(fastBGM))
	}
}
