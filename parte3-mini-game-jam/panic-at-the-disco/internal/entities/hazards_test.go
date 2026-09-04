package entities

import (
	"image/color"
	"testing"
)

func TestFallingDiscoBallLifecycle(t *testing.T) {
	db := NewFallingDiscoBall(100.0, 100.0, 25.0, 0.2)
	p := NewPlayer(100.0, 100.0)

	if db.State != HazardTelegraph {
		t.Fatalf("Expected initial state HazardTelegraph, got %v", db.State)
	}

	// Step past telegraph duration (0.2s)
	db.Update(0.25, p, nil, nil)
	if db.State != HazardFalling {
		t.Fatalf("Expected state HazardFalling after telegraph, got %v", db.State)
	}

	// Step while falling until impact
	for i := 0; i < 60; i++ {
		db.Update(0.016, p, nil, nil)
		if db.State == HazardImpact || db.State == HazardFinished {
			break
		}
	}

	if db.State != HazardImpact && db.State != HazardFinished {
		t.Fatalf("Expected state HazardImpact or Finished, got %v", db.State)
	}
}

func TestExitDoorDetection(t *testing.T) {
	door := NewExitDoor(200.0, 50.0, 40.0, 40.0)
	pInside := NewPlayer(210.0, 60.0)
	pOutside := NewPlayer(100.0, 200.0)

	if !door.IsPlayerInside(pInside) {
		t.Fatal("Expected player to be detected inside exit door")
	}
	if door.IsPlayerInside(pOutside) {
		t.Fatal("Player outside should not trigger exit door")
	}
}

func TestDrinkPuddleSlip(t *testing.T) {
	puddle := NewDrinkPuddle(150.0, 150.0, 20.0, color.RGBA{255, 165, 0, 255})
	pOnPuddle := NewPlayer(150.0, 150.0)
	pFarAway := NewPlayer(300.0, 300.0)

	puddle.Update(pFarAway, nil)
	if pFarAway.IsSlipping {
		t.Fatal("Player far away should not slip")
	}

	puddle.Update(pOnPuddle, nil)
	if !pOnPuddle.IsSlipping {
		t.Fatal("Player standing on puddle should slip")
	}
}
