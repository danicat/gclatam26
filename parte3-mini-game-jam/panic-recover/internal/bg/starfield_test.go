package bg

import (
	"testing"
)

func TestStarfieldUpdate(t *testing.T) {
	sf := NewStarfield(640, 360)
	if len(sf.stars) == 0 {
		t.Fatalf("expected stars to be initialized")
	}

	initialOffset := sf.gridOffset
	sf.Update(0.1)

	if sf.UpdateCount() != 1 {
		t.Errorf("expected update count 1, got %d", sf.UpdateCount())
	}

	if sf.gridOffset == initialOffset {
		t.Errorf("expected grid offset to change after update")
	}
}
