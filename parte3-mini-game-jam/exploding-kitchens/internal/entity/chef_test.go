package entity

import (
	"testing"
)

func TestChefMovementAndBounds(t *testing.T) {
	chef := NewChef(100, 100)

	// Move right
	chef.Update(1.0, 1.0, 0.0, 0, 200, 0, 200)
	if chef.Dir != DirRight {
		t.Errorf("expected DirRight, got %v", chef.Dir)
	}
	if chef.X <= 100 {
		t.Errorf("expected chef to move right, got X = %f", chef.X)
	}

	// Move out of bounds left
	chef.Update(1.0, -10.0, 0.0, 50, 200, 50, 200)
	if chef.X < 50 {
		t.Errorf("expected chef X to be clamped to minX=50, got %f", chef.X)
	}

	cx, cy := chef.Center()
	if cx != chef.X+chef.W/2 || cy != chef.Y+chef.H/2 {
		t.Errorf("center mismatch")
	}
}

func TestCatShooAndFlee(t *testing.T) {
	cat := NewCat(50, 50)
	stn := NewStation(100, 100, StationStoveTop)
	cat.TargetStn = stn
	stn.CatBoost = true
	cat.State = CatSittingOnStation

	// Shoo cat
	cat.Shoo(10, 10)
	if cat.State != CatShooedFleeing {
		t.Errorf("expected CatShooedFleeing, got %v", cat.State)
	}
	if stn.CatBoost {
		t.Errorf("cat boost should be false after shoo")
	}
	if cat.TargetStn != nil {
		t.Errorf("target station should be detached after shoo")
	}
}
