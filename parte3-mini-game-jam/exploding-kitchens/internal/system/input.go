package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Action represents an abstract player input action.
type Action int

const (
	ActionNone Action = iota
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionInteract // Space or E
	ActionDrop     // Q or Shift
	ActionPause    // Esc or P
	ActionConfirm  // Enter or Space
)

// InputState stores the instantaneous input status for the current frame.
type InputState struct {
	MoveX        float64
	MoveY        float64
	InteractJust bool
	InteractHeld bool
	DropJust     bool
	PauseJust    bool
	ConfirmJust  bool
	ToggleFull   bool
}

// InputManager queries physical devices and computes abstract input actions.
type InputManager struct{}

// NewInputManager creates a new InputManager.
func NewInputManager() *InputManager {
	return &InputManager{}
}

// Poll reads all active devices (Keyboard, Gamepad) and returns the consolidated InputState.
func (im *InputManager) Poll() InputState {
	var state InputState

	// Movement: Keyboard (WASD + Arrow Keys)
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		state.MoveY -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		state.MoveY += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		state.MoveX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		state.MoveX += 1.0
	}

	// Movement: Gamepad (First connected gamepad)
	gamepadIDs := ebiten.AppendGamepadIDs(nil)
	if len(gamepadIDs) > 0 {
		id := gamepadIDs[0]
		// Analog stick
		axisX := ebiten.GamepadAxisValue(id, 0)
		axisY := ebiten.GamepadAxisValue(id, 1)
		const deadzone = 0.25
		if axisX < -deadzone || axisX > deadzone {
			state.MoveX = axisX
		}
		if axisY < -deadzone || axisY > deadzone {
			state.MoveY = axisY
		}

		// D-Pad
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton0) {
			state.InteractHeld = true
		}
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton0) {
			state.InteractJust = true
			state.ConfirmJust = true
		}
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton2) {
			state.DropJust = true
		}
		if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton6) {
			state.PauseJust = true
		}
	}

	// Normalizing diagonal keyboard movement
	if state.MoveX != 0 && state.MoveY != 0 {
		const invSqrt2 = 0.70710678118
		state.MoveX *= invSqrt2
		state.MoveY *= invSqrt2
	}

	// Actions: Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyE) {
		state.InteractJust = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyE) {
		state.InteractHeld = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyShift) {
		state.DropJust = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		state.PauseJust = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.ConfirmJust = true
	}

	// Fullscreen toggle (F11 or Alt+Enter)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		state.ToggleFull = true
	}

	return state
}
