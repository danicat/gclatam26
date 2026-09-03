package game

import (
	"testing"
)

type mockScene struct {
	entered bool
	exited  bool
}

func (m *mockScene) Enter(ctx *GameContext)                                  { m.entered = true }
func (m *mockScene) Update(dt float64, in systemMock) (SceneID, error)       { return SceneKeepCurrent, nil }
func (m *mockScene) Draw(screen mockScreen)                                  {}
func (m *mockScene) Exit()                                                   { m.exited = true }

type systemMock struct{}
type mockScreen struct{}

func TestGameLayoutConstantResolution(t *testing.T) {
	title := &mockScene{}
	play := &mockScene{}
	gameover := &mockScene{}

	// Verify constant 16:9 virtual pixel canvas scaling
	const expectedW = 320
	const expectedH = 180

	g := &Game{}
	w, h := g.Layout(1920, 1080)
	if w != expectedW || h != expectedH {
		t.Errorf("expected (%d, %d) on 1080p, got (%d, %d)", expectedW, expectedH, w, h)
	}

	w, h = g.Layout(3840, 2160)
	if w != expectedW || h != expectedH {
		t.Errorf("expected (%d, %d) on 4k, got (%d, %d)", expectedW, expectedH, w, h)
	}

	_ = title
	_ = play
	_ = gameover
}
