package levels

import (
	"testing"
)

func TestAllLevelsDefinition(t *testing.T) {
	if len(AllLevels) != 10 {
		t.Fatalf("expected 10 levels, got %d", len(AllLevels))
	}

	for _, l := range AllLevels {
		if l.Title == "" {
			t.Errorf("level %d has empty title", l.ID)
		}
		if len(l.CodeLines) == 0 {
			t.Errorf("level %d has empty code lines", l.ID)
		}
		if l.TargetLineIndex < 0 || l.TargetLineIndex >= len(l.CodeLines) {
			t.Errorf("level %d target line index %d out of bounds (len=%d)", l.ID, l.TargetLineIndex, len(l.CodeLines))
		}
		if l.TimeLimit <= 0 {
			t.Errorf("level %d has invalid time limit %f", l.ID, l.TimeLimit)
		}
		if l.Validate == nil {
			t.Errorf("level %d missing validate function", l.ID)
		}
	}
}

func TestLevelValidationSolutions(t *testing.T) {
	testSolutions := []struct {
		levelID  int
		lineIdx  int
		solution string
		valid    bool
	}{
		{1, 1, "gopher := &Gopher{}", true},
		{1, 2, "gopher = &Gopher{Score: 100}", true},
		{1, 2, "random broken code", false},
		{2, 2, "if count != 0 { return total / count }", true},
		{2, 2, "ratio := total / count", false},
		{3, 2, "last := items[n-1]", true},
		{3, 2, "last := items[len(items)-1]", true},
		{4, 1, "inventory := make(map[string]int)", true},
		{5, 1, "ch <- msg; close(ch)", true},
		{6, 3, "// close(ch)", true},
		{7, 1, "if n <= 0 { return }", true},
		{8, 1, "name, ok := data.(string)", true},
		{9, 2, "mu.Lock(); hits[key]++; mu.Unlock()", true},
		{10, 1, "defer func() { recover() }()", true},
	}

	for _, tc := range testSolutions {
		lvl := AllLevels[tc.levelID-1]
		got := lvl.Validate(tc.lineIdx, tc.solution)
		if got != tc.valid {
			t.Errorf("Level %d line %d validation for '%s': expected %v, got %v",
				tc.levelID, tc.lineIdx, tc.solution, tc.valid, got)
		}
	}
}
