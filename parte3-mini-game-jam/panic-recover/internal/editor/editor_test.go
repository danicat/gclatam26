package editor

import (
	"testing"
)

func TestEditorBufferOperations(t *testing.T) {
	ed := NewEditor()
	ed.StartEditing("gopher.Score = 100")

	if !ed.IsActive {
		t.Fatalf("expected editor to be active")
	}
	if ed.CurrentBufferString() != "gopher.Score = 100" {
		t.Errorf("expected initial text 'gopher.Score = 100', got '%s'", ed.CurrentBufferString())
	}
	if ed.CursorPos != len("gopher.Score = 100") {
		t.Errorf("expected cursor at end (%d), got %d", len("gopher.Score = 100"), ed.CursorPos)
	}

	ed.CancelEditing()
	if ed.IsActive {
		t.Errorf("expected editor to be inactive after cancel")
	}
}
