package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Logical input actions.
type Action int

const (
	ActionUp Action = iota
	ActionDown
	ActionLeft
	ActionRight
	ActionDash
	ActionConfirm
	ActionRestart
	ActionPause
)

// InputState holds the current frame's polled input state.
type InputState struct {
	MoveX        float64
	MoveY        float64
	DashJustDown bool
	ConfirmJustDown bool
	RestartJustDown bool
	PauseJustDown   bool
}

// Poll reads keyboard and gamepad inputs and returns an InputState.
func Poll() InputState {
	var state InputState

	var dx, dy float64

	// Keyboard movement
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1.0
	}

	// Action buttons
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		state.DashJustDown = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.ConfirmJustDown = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		state.RestartJustDown = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		state.PauseJustDown = true
	}

	// Gamepad polling
	gamepadIDs := ebiten.AppendGamepadIDs(nil)
	if len(gamepadIDs) > 0 {
		id := gamepadIDs[0]
		// D-Pad buttons
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton11) { // Up
			dy -= 1.0
		}
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton12) { // Right
			dx += 1.0
		}
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton13) { // Down
			dy += 1.0
		}
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton14) { // Left
			dx -= 1.0
		}

		// Analog Stick (Axis 0 = X, Axis 1 = Y)
		axesCount := ebiten.GamepadAxisCount(id)
		if axesCount >= 2 {
			ax := ebiten.GamepadAxisValue(id, 0)
			ay := ebiten.GamepadAxisValue(id, 1)
			const deadzone = 0.2
			if math.Abs(ax) > deadzone {
				dx += ax
			}
			if math.Abs(ay) > deadzone {
				dy += ay
			}
		}

		// Buttons: South (Button 0) = Dash / Confirm
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton0) {
			state.DashJustDown = true
			state.ConfirmJustDown = true
		}
		// Button West (Button 2) = Restart
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton2) {
			state.RestartJustDown = true
		}
		// Start (Button 9) = Pause
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton9) {
			state.PauseJustDown = true
		}
	}

	// Normalize movement vector to avoid faster diagonal movement
	lenSq := dx*dx + dy*dy
	if lenSq > 0 {
		l := math.Sqrt(lenSq)
		if l > 1.0 {
			dx /= l
			dy /= l
		}
	}

	state.MoveX = dx
	state.MoveY = dy
	return state
}
