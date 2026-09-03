package audio

import (
	"testing"
)

func TestAudioSynthesis(t *testing.T) {
	laser := synthLaser()
	if len(laser) == 0 {
		t.Fatalf("expected synthLaser to generate non-empty buffer")
	}

	explosion := synthExplosion()
	if len(explosion) == 0 {
		t.Fatalf("expected synthExplosion to generate non-empty buffer")
	}

	chime := synthRecoverChime()
	if len(chime) == 0 {
		t.Fatalf("expected synthRecoverChime to generate non-empty buffer")
	}

	bgm := synthBGM()
	if len(bgm) == 0 {
		t.Fatalf("expected synthBGM to generate non-empty buffer")
	}
}

func TestInfiniteLoopReader(t *testing.T) {
	sampleData := []byte{1, 2, 3, 4}
	r := &infiniteLoopReader{data: sampleData, offset: 0}

	readBuf := make([]byte, 10)
	n, err := r.Read(readBuf)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if n != 10 {
		t.Fatalf("expected 10 bytes read, got %d", n)
	}
	// Check looping content
	expected := []byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2}
	for i := range expected {
		if readBuf[i] != expected[i] {
			t.Errorf("byte %d: expected %d, got %d", i, expected[i], readBuf[i])
		}
	}
}
