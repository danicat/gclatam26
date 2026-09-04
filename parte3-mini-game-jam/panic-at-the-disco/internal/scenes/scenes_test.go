package scenes

import (
	"testing"

	"panic-at-the-disco/internal/levels"
)

func TestGameOverSceneCreation(t *testing.T) {
	gos := NewGameOverScene("CRUSHED BY ROOF", 1500, 32.5)
	if gos.reason != "CRUSHED BY ROOF" {
		t.Fatalf("Expected reason 'CRUSHED BY ROOF', got %s", gos.reason)
	}
	if gos.score != 1500 {
		t.Fatalf("Expected score 1500, got %d", gos.score)
	}
}

func TestVictorySceneCreation(t *testing.T) {
	vs := NewVictoryScene(3500, 2, 75.2)
	if vs.score != 3500 {
		t.Fatalf("Expected score 3500, got %d", vs.score)
	}
	if vs.lives != 2 {
		t.Fatalf("Expected 2 lives, got %d", vs.lives)
	}
}

func TestStageClearSceneProgression(t *testing.T) {
	sc := NewStageClearScene(levels.ZoneVIPLounge, 1200, 3)

	// Step past timeout (2.5s) to trigger transition
	action := sc.Update(2.6)
	if action.Type != ActionSwitchScene {
		t.Fatalf("Expected ActionSwitchScene after timer expiration, got %v", action.Type)
	}
	if action.NextScene != ScenePlay {
		t.Fatalf("Expected NextScene ScenePlay, got %v", action.NextScene)
	}
	if action.TargetZone != levels.ZoneVIPLounge {
		t.Fatalf("Expected TargetZone ZoneVIPLounge, got %v", action.TargetZone)
	}
}
