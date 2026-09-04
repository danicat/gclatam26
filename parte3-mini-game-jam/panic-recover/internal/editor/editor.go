package editor

import (
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Editor manages user navigation and inline line editing.
type Editor struct {
	IsActive          bool
	SelectedLineIndex int
	Buffer            []rune
	CursorPos         int
	BlinkTimer        float64
	CursorVisible     bool
	HintVisible       bool

	// Key repeat tracking for backspace / arrow navigation
	backspaceTimer float64
	leftTimer      float64
	rightTimer     float64
	inputChars     []rune
}

// NewEditor creates a new Editor instance.
func NewEditor() *Editor {
	return &Editor{
		CursorVisible: true,
		inputChars:    make([]rune, 0, 16),
	}
}

// SetLineCount adjusts selection if lines change.
func (e *Editor) SetLineCount(count int) {
	if count <= 0 {
		e.SelectedLineIndex = 0
		return
	}
	if e.SelectedLineIndex >= count {
		e.SelectedLineIndex = count - 1
	}
}

// StartEditing initiates inline editing for the currently selected line.
func (e *Editor) StartEditing(initialText string) {
	e.IsActive = true
	e.Buffer = []rune(initialText)
	e.CursorPos = len(e.Buffer)
	e.BlinkTimer = 0
	e.CursorVisible = true
}

// CancelEditing aborts the current edit without submitting.
func (e *Editor) CancelEditing() {
	e.IsActive = false
	e.Buffer = nil
	e.CursorPos = 0
}

// UpdateNavigation handles up/down line selection.
// Returns (navigated bool, enterEdit bool, toggleHint bool).
func (e *Editor) UpdateNavigation(numLines int) (bool, bool, bool) {
	if e.IsActive {
		return false, false, false
	}

	navigated := false
	toggleHint := false

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if e.SelectedLineIndex > 0 {
			e.SelectedLineIndex--
			navigated = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if e.SelectedLineIndex < numLines-1 {
			e.SelectedLineIndex++
			navigated = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) || inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		e.HintVisible = !e.HintVisible
		toggleHint = true
	}

	enterEdit := inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyE)

	return navigated, enterEdit, toggleHint
}

// UpdateEditing handles typing, cursor movements, and submission.
// Returns (submittedText string, submitted bool, cancelled bool, keyPressed bool).
func (e *Editor) UpdateEditing(dt float64) (string, bool, bool, bool) {
	if !e.IsActive {
		return "", false, false, false
	}

	// Update cursor blink
	e.BlinkTimer += dt
	if e.BlinkTimer >= 0.45 {
		e.BlinkTimer = 0
		e.CursorVisible = !e.CursorVisible
	}

	keyPressed := false

	// Cancel on Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		e.CancelEditing()
		return "", false, true, true
	}

	// Submit on Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		submitted := string(e.Buffer)
		e.CancelEditing()
		return submitted, true, false, true
	}

	// Home / End
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		e.CursorPos = 0
		e.CursorVisible = true
		keyPressed = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		e.CursorPos = len(e.Buffer)
		e.CursorVisible = true
		keyPressed = true
	}

	// Left Arrow (with repeat)
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if e.CursorPos > 0 {
			e.CursorPos--
			e.CursorVisible = true
			keyPressed = true
		}
		e.leftTimer = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		e.leftTimer += dt
		if e.leftTimer > 0.35 {
			if e.CursorPos > 0 {
				e.CursorPos--
				e.CursorVisible = true
				keyPressed = true
			}
			e.leftTimer = 0.28
		}
	} else {
		e.leftTimer = 0
	}

	// Right Arrow (with repeat)
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if e.CursorPos < len(e.Buffer) {
			e.CursorPos++
			e.CursorVisible = true
			keyPressed = true
		}
		e.rightTimer = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		e.rightTimer += dt
		if e.rightTimer > 0.35 {
			if e.CursorPos < len(e.Buffer) {
				e.CursorPos++
				e.CursorVisible = true
				keyPressed = true
			}
			e.rightTimer = 0.28
		}
	} else {
		e.rightTimer = 0
	}

	// Backspace (with repeat)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if e.CursorPos > 0 {
			e.Buffer = append(e.Buffer[:e.CursorPos-1], e.Buffer[e.CursorPos:]...)
			e.CursorPos--
			e.CursorVisible = true
			keyPressed = true
		}
		e.backspaceTimer = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		e.backspaceTimer += dt
		if e.backspaceTimer > 0.35 {
			if e.CursorPos > 0 {
				e.Buffer = append(e.Buffer[:e.CursorPos-1], e.Buffer[e.CursorPos:]...)
				e.CursorPos--
				e.CursorVisible = true
				keyPressed = true
			}
			e.backspaceTimer = 0.28
		}
	} else {
		e.backspaceTimer = 0
	}

	// Delete key
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		if e.CursorPos < len(e.Buffer) {
			e.Buffer = append(e.Buffer[:e.CursorPos], e.Buffer[e.CursorPos+1:]...)
			e.CursorVisible = true
			keyPressed = true
		}
	}

	// Tab key (insert 4 spaces)
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		spaces := []rune("    ")
		e.Buffer = append(e.Buffer[:e.CursorPos], append(spaces, e.Buffer[e.CursorPos:]...)...)
		e.CursorPos += 4
		e.CursorVisible = true
		keyPressed = true
	}

	// Handle typed characters
	e.inputChars = ebiten.AppendInputChars(e.inputChars[:0])
	for _, r := range e.inputChars {
		// Filter out non-printable runes
		if unicode.IsPrint(r) {
			e.Buffer = append(e.Buffer[:e.CursorPos], append([]rune{r}, e.Buffer[e.CursorPos:]...)...)
			e.CursorPos++
			e.CursorVisible = true
			keyPressed = true
		}
	}

	return "", false, false, keyPressed
}

// CurrentBufferString returns string representation of current edit buffer.
func (e *Editor) CurrentBufferString() string {
	return string(e.Buffer)
}
