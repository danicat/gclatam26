package app

import (
	"testing"

	"panic-recover/internal/game"
)

func TestLayoutUsesVirtualResolution(t *testing.T) {
	t.Parallel()

	a := &App{}
	width, height := a.Layout(1920, 1080)

	if width != game.VirtualWidth || height != game.VirtualHeight {
		t.Fatalf("Layout() = (%d, %d), want (%d, %d)", width, height, game.VirtualWidth, game.VirtualHeight)
	}
}
