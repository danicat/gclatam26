package entity

import (
	"testing"
)

func TestStationLifecycleAndDetonation(t *testing.T) {
	stn := NewStation(10, 10, StationPressureCooker)
	stn.MaxTime = 10.0
	stn.Timer = 0.0

	// Initially cooking
	stn.Update(1.0)
	if stn.State != StateCooking {
		t.Fatalf("expected StateCooking, got %v", stn.State)
	}

	// Advance to warning (55% = 5.5s)
	stn.Update(5.0)
	if stn.State != StateWarning {
		t.Fatalf("expected StateWarning, got %v", stn.State)
	}

	// Advance to panic (80% = 8.0s)
	stn.Update(2.0)
	if stn.State != StatePanic {
		t.Fatalf("expected StatePanic, got %v", stn.State)
	}

	// In early panic (8.0s / 10s = 80%), not clutch yet
	if stn.IsClutch() {
		t.Fatalf("expected IsClutch to be false at 80%% progress")
	}

	// Advance to clutch (8.6s / 10s = 86%)
	stn.Update(0.6)
	if !stn.IsClutch() {
		t.Fatalf("expected IsClutch to be true at 86%% progress")
	}

	// Advance to detonation (10.0s)
	exploded := stn.Update(1.5)
	if !exploded {
		t.Fatalf("expected exploded to be true when reaching 100%%")
	}
	if stn.State != StateExploded {
		t.Fatalf("expected StateExploded, got %v", stn.State)
	}
}

func TestStationRecoveryToolMatching(t *testing.T) {
	stnFryer := NewStation(0, 0, StationDeepFryer)
	stnFryer.State = StatePanic

	// Deep fryer requires fire extinguisher
	if stnFryer.CanRecover(ToolNone) {
		t.Errorf("fryer should not recover with ToolNone")
	}
	if stnFryer.CanRecover(ToolIce) {
		t.Errorf("fryer should not recover with ToolIce")
	}
	if !stnFryer.CanRecover(ToolExtinguisher) {
		t.Errorf("fryer should recover with ToolExtinguisher")
	}

	// Exploded station requires wrench
	stnFryer.State = StateExploded
	if stnFryer.CanRecover(ToolExtinguisher) {
		t.Errorf("exploded fryer should not recover with ToolExtinguisher")
	}
	if !stnFryer.CanRecover(ToolWrench) {
		t.Errorf("exploded fryer should recover with ToolWrench")
	}

	// Repair station
	stnFryer.Repair()
	if stnFryer.State != StateCooking {
		t.Errorf("expected StateCooking after repair, got %v", stnFryer.State)
	}
}

func TestCatBoostSpeed(t *testing.T) {
	stn := NewStation(0, 0, StationStoveTop)
	stn.MaxTime = 100.0
	stn.Timer = 0.0

	// Normal tick
	stn.Update(1.0)
	normalElapsed := stn.Timer

	// Boosted tick with cat
	stn.CatBoost = true
	stn.Update(1.0)
	boostedElapsed := stn.Timer - normalElapsed

	if boostedElapsed <= normalElapsed*2.0 {
		t.Errorf("expected cat boost to at least double timer progress: got %f vs %f", boostedElapsed, normalElapsed)
	}
}
