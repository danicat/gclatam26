package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"panic-recover/internal/art"
	"panic-recover/internal/audio"
	"panic-recover/internal/editor"
	"panic-recover/internal/levels"
)

const (
	VirtualWidth  = 640
	VirtualHeight = 360

	CodeStartX      = 120
	PanicThresholdY = 280
)

// GameState represents the current active scene in the game FSM.
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateLevelRecovered
	StatePanicCrash
	StateGameWin
)

// Game implements the ebiten.Game interface.
type Game struct {
	State GameState

	// Levels & Progress
	CurrentLevelIdx int
	CurrentLevel    levels.Level
	WorkingLines    []string
	Score           int
	TotalTimeUsed   float64

	// Gameplay Timers & Movement
	TimeRemaining float64
	MaxTime       float64
	CodeOffsetY   float64
	WarningPlayed bool
	FeedbackMsg   string
	FeedbackTimer float64
	ScreenShake   float64

	// Subsystems
	editor    *editor.Editor
	particles *art.ParticleSystem
	matrix    *art.MatrixRain
	audioSys  *audio.SoundSystem

	// Attract / Title timers
	titleBlinkTimer float64
	winConfetti     float64
}

// NewGame initializes the game instance.
func NewGame() *Game {
	g := &Game{
		State:     StateTitle,
		editor:    editor.NewEditor(),
		particles: art.NewParticleSystem(300),
		matrix:    art.NewMatrixRain(45, VirtualWidth, VirtualHeight),
		audioSys:  audio.GetSoundSystem(),
	}
	return g
}

// StartLevel loads and resets state for level with given index.
func (g *Game) StartLevel(idx int) {
	if idx >= len(levels.AllLevels) {
		g.State = StateGameWin
		g.audioSys.PlayRecover()
		return
	}

	g.CurrentLevelIdx = idx
	g.CurrentLevel = levels.AllLevels[idx]
	g.WorkingLines = make([]string, len(g.CurrentLevel.CodeLines))
	copy(g.WorkingLines, g.CurrentLevel.CodeLines)

	g.TimeRemaining = g.CurrentLevel.TimeLimit
	g.MaxTime = g.CurrentLevel.TimeLimit
	g.CodeOffsetY = 65.0
	g.WarningPlayed = false
	g.FeedbackMsg = ""
	g.FeedbackTimer = 0
	g.ScreenShake = 0

	g.editor.CancelEditing()
	g.editor.SelectedLineIndex = g.CurrentLevel.TargetLineIndex
	g.editor.HintVisible = false

	g.State = StatePlaying
	g.audioSys.StartBGM()
}

// Update handles logic tick at 60 TPS.
func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// Global hotkeys: Fullscreen & Mute
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.audioSys.ToggleMute()
	}

	// Update background visual effects
	g.particles.Update(dt)
	g.matrix.Update(dt, VirtualHeight)

	if g.ScreenShake > 0 {
		g.ScreenShake -= dt * 8.0
		if g.ScreenShake < 0 {
			g.ScreenShake = 0
		}
	}
	if g.FeedbackTimer > 0 {
		g.FeedbackTimer -= dt
		if g.FeedbackTimer <= 0 {
			g.FeedbackMsg = ""
		}
	}

	switch g.State {
	case StateTitle:
		g.updateTitle(dt)
	case StatePlaying:
		g.updatePlaying(dt)
	case StateLevelRecovered:
		g.updateRecovered(dt)
	case StatePanicCrash:
		g.updatePanicCrash(dt)
	case StateGameWin:
		g.updateGameWin(dt)
	}

	return nil
}

func (g *Game) updateTitle(dt float64) {
	g.titleBlinkTimer += dt
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.audioSys.PlaySelect()
		g.Score = 0
		g.TotalTimeUsed = 0
		g.StartLevel(0)
	}
}

func (g *Game) updatePlaying(dt float64) {
	g.TimeRemaining -= dt
	g.TotalTimeUsed += dt

	// Code falls down continuously
	g.CodeOffsetY += g.CurrentLevel.FallSpeed * dt

	// Warning siren pulse when time is critical (< 5s)
	if g.TimeRemaining <= 5.0 && g.TimeRemaining > 0 {
		if math.Mod(g.TimeRemaining, 1.0) < 0.15 && !g.WarningPlayed {
			g.audioSys.PlayWarning()
			g.WarningPlayed = true
		} else if math.Mod(g.TimeRemaining, 1.0) >= 0.15 {
			g.WarningPlayed = false
		}
	}

	// Calculate bottom line position
	numLines := len(g.WorkingLines)
	bottomY := g.CodeOffsetY + float64(numLines*art.LineHeight)

	// Panic loss condition: Timer expired or code touched the Panic Threshold
	if g.TimeRemaining <= 0 || bottomY >= PanicThresholdY {
		g.triggerPanic()
		return
	}

	// Handle navigation or inline editing
	if !g.editor.IsActive {
		navigated, enterEdit, toggleHint := g.editor.UpdateNavigation(numLines)
		if navigated {
			g.audioSys.PlaySelect()
		}
		if toggleHint {
			g.audioSys.PlaySelect()
		}
		if enterEdit {
			g.editor.StartEditing(g.WorkingLines[g.editor.SelectedLineIndex])
			g.audioSys.PlayEditMode()
		}
	} else {
		submittedText, submitted, cancelled, keyPressed := g.editor.UpdateEditing(dt)
		if keyPressed {
			g.audioSys.PlayKeyClick()
		}
		if cancelled {
			g.audioSys.PlaySelect()
		}
		if submitted {
			g.evaluateSubmission(submittedText)
		}
	}
}

func (g *Game) evaluateSubmission(text string) {
	idx := g.editor.SelectedLineIndex
	isCorrect := g.CurrentLevel.Validate(idx, text)

	if isCorrect {
		// Update working line
		g.WorkingLines[idx] = text

		// Reward score based on remaining time
		timeBonus := int(g.TimeRemaining * 50.0)
		levelClearBonus := (g.CurrentLevelIdx + 1) * 200
		g.Score += timeBonus + levelClearBonus

		// Spawn celebratory green particles
		lineY := g.CodeOffsetY + float64(idx*art.LineHeight)
		g.particles.SpawnSparks(VirtualWidth/2, lineY, 45, art.ColorGreenRecover)

		g.audioSys.PlayRecover()
		g.State = StateLevelRecovered
	} else {
		// Wrong fix: time penalty and error buzz
		g.TimeRemaining -= 2.0
		if g.TimeRemaining < 0 {
			g.TimeRemaining = 0
		}
		g.ScreenShake = 1.0
		g.FeedbackMsg = "COMPILER WARNING: Still panics! -2.0s penalty."
		g.FeedbackTimer = 2.5
		g.particles.SpawnSparks(float64(CodeStartX+100), g.CodeOffsetY+float64(idx*art.LineHeight), 15, art.ColorRedPanic)
		g.audioSys.PlayWarning()
	}
}

func (g *Game) triggerPanic() {
	g.ScreenShake = 2.0
	g.State = StatePanicCrash
	g.audioSys.PlayPanic()
	g.particles.SpawnSparks(VirtualWidth/2, PanicThresholdY, 80, art.ColorRedPanic)
}

func (g *Game) updateRecovered(dt float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.audioSys.PlaySelect()
		g.StartLevel(g.CurrentLevelIdx + 1)
	}
}

func (g *Game) updatePanicCrash(dt float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.audioSys.PlaySelect()
		// Retry current level
		g.StartLevel(g.CurrentLevelIdx)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.State = StateTitle
	}
}

func (g *Game) updateGameWin(dt float64) {
	g.winConfetti += dt
	if math.Mod(g.winConfetti, 0.2) < dt {
		g.particles.SpawnSparks(float64(VirtualWidth/4+int(g.winConfetti*60)%VirtualWidth/2), 60, 10, art.ColorCyanGlow)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = StateTitle
	}
}

// Draw renders the active scene.
func (g *Game) Draw(screen *ebiten.Image) {
	// Base terminal background
	screen.Fill(art.ColorBgDark)

	// Apply screen shake offset if active
	shakeX, shakeY := 0, 0
	if g.ScreenShake > 0 {
		shakeX = int((math.Sin(g.ScreenShake*40.0)) * g.ScreenShake * 4.0)
		shakeY = int((math.Cos(g.ScreenShake*40.0)) * g.ScreenShake * 4.0)
	}

	// Draw digital matrix rain
	g.matrix.Draw(screen)

	switch g.State {
	case StateTitle:
		g.drawTitle(screen)
	case StatePlaying:
		g.drawPlaying(screen, shakeX, shakeY)
	case StateLevelRecovered:
		g.drawRecovered(screen)
	case StatePanicCrash:
		g.drawPanicCrash(screen, shakeX, shakeY)
	case StateGameWin:
		g.drawGameWin(screen)
	}

	// Draw particles and retro scanlines on top of everything
	g.particles.Draw(screen)
	art.DrawScanlines(screen, VirtualWidth, VirtualHeight)
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	// Center title box
	art.DrawRect(screen, 60, 45, 520, 270, art.ColorPanelBg)
	art.DrawBorder(screen, 60, 45, 520, 270, art.ColorCyanGlow)

	// Decorative window title bar
	art.DrawRect(screen, 60, 45, 520, 18, art.ColorPanelBorder)
	art.DrawText(screen, "TERMINAL: GOPHER RUNTIME DEBUGGER", 70, 58, art.ColorCyanGlow)

	// Window control dots
	art.DrawRect(screen, 545, 50, 8, 8, art.ColorRedPanic)
	art.DrawRect(screen, 557, 50, 8, 8, art.ColorYellowAlert)
	art.DrawRect(screen, 569, 50, 8, 8, art.ColorGreenRecover)

	// Main Title
	art.DrawText(screen, "==================================================", 80, 90, art.ColorCyanGlow)
	art.DrawText(screen, "         P A N I C ! ! !   ( &   r e c o v e r ? )", 80, 110, art.ColorRedPanic)
	art.DrawText(screen, "==================================================", 80, 130, art.ColorCyanGlow)

	art.DrawText(screen, "GopherCon LATAM 2026 Mini Game Jam Edition", 175, 155, art.ColorYellowAlert)
	art.DrawText(screen, "Defeat the falling bugs before the runtime crashes!", 145, 175, art.ColorTextWhite)

	// Controls list
	art.DrawText(screen, "CONTROLS:", 120, 210, art.ColorCyanGlow)
	art.DrawText(screen, "- [UP / DOWN] or [W / S] : Select code line", 140, 228, art.ColorTextWhite)
	art.DrawText(screen, "- [ENTER / SPACE]       : Edit line & submit fix", 140, 244, art.ColorTextWhite)
	art.DrawText(screen, "- [H]                   : Toggle diagnostic hint", 140, 260, art.ColorTextWhite)
	art.DrawText(screen, "- [M]                   : Mute / Unmute audio", 140, 276, art.ColorTextWhite)

	// Blinking start prompt
	if int(g.titleBlinkTimer*2)%2 == 0 {
		art.DrawText(screen, ">>> PRESS [ENTER] OR [SPACE] TO START <<<", 170, 302, art.ColorGreenRecover)
	}
}

func (g *Game) drawPlaying(screen *ebiten.Image, sx, sy int) {
	// Top HUD Bar
	art.DrawRect(screen, 0, 0, VirtualWidth, 42, art.ColorPanelBg)
	art.DrawRect(screen, 0, 42, VirtualWidth, 1, art.ColorCyanGlow)

	// HUD Details
	lvlStr := fmt.Sprintf("LEVEL: %02d/10", g.CurrentLevelIdx+1)
	scoreStr := fmt.Sprintf("SCORE: %05d", g.Score)
	art.DrawText(screen, lvlStr, 16+sx, 16+sy, art.ColorCyanGlow)
	art.DrawText(screen, g.CurrentLevel.Title, 130+sx, 16+sy, art.ColorYellowAlert)
	art.DrawText(screen, scoreStr, 530+sx, 16+sy, art.ColorGreenRecover)

	// Panic Timer Bar
	timerPct := g.TimeRemaining / g.MaxTime
	if timerPct < 0 {
		timerPct = 0
	}
	barW := int(timerPct * 400.0)
	barColor := art.ColorGreenRecover
	if timerPct < 0.5 {
		barColor = art.ColorYellowAlert
	}
	if timerPct < 0.25 {
		barColor = art.ColorRedPanic
	}
	art.DrawRect(screen, 130+sx, 24+sy, 380, 10, art.ColorPanelBorder)
	art.DrawRect(screen, 130+sx, 24+sy, int(float64(barW)*(380.0/400.0)), 10, barColor)
	timeStr := fmt.Sprintf("PANIC IN: %04.1fs", g.TimeRemaining)
	art.DrawText(screen, timeStr, 520+sx, 33+sy, barColor)

	// Draw Laser Panic Horizon at bottom
	laserPulse := math.Sin(g.TotalTimeUsed * 8.0)
	laserColor := art.ColorLaserRed
	if laserPulse > 0 {
		laserColor = art.ColorYellowAlert
	}
	art.DrawRect(screen, 0, PanicThresholdY, VirtualWidth, 2, laserColor)
	art.DrawText(screen, "▲▲▲ RUNTIME PANIC HORIZON (CRASH ZONE) ▲▲▲", 185, PanicThresholdY+12, laserColor)

	// Draw Code Lines Window Container
	startY := int(g.CodeOffsetY) + sy
	numLines := len(g.WorkingLines)
	codeBoxW := 500
	codeBoxH := numLines*art.LineHeight + 10
	boxX := CodeStartX - 45 + sx
	boxY := startY - 12
	gutterW := 38

	// Code editor backdrop panel
	art.DrawRect(screen, boxX, boxY, codeBoxW, codeBoxH, art.ColorCodeBg)
	art.DrawBorder(screen, boxX, boxY, codeBoxW, codeBoxH, art.ColorCodeBorder)

	// Gutter column on left with vertical separator line
	art.DrawRect(screen, boxX, boxY, gutterW, codeBoxH, art.ColorGutterBg)
	art.DrawRect(screen, boxX+gutterW, boxY, 1, codeBoxH, art.ColorGutterBorder)

	for i := 0; i < numLines; i++ {
		lineY := startY + i*art.LineHeight
		isSelected := (i == g.editor.SelectedLineIndex)
		selX := boxX + gutterW + 1
		selW := codeBoxW - gutterW - 1

		// Draw selection / edit background highlight
		if isSelected {
			if g.editor.IsActive {
				// EDIT MODE: deep warm amber background with vivid gold border & left accent
				art.DrawRect(screen, selX, lineY-11, selW, art.LineHeight, art.ColorEditBar)
				art.DrawBorder(screen, selX, lineY-11, selW, art.LineHeight, art.ColorEditBorder)
				art.DrawRect(screen, selX, lineY-11, 4, art.LineHeight, art.ColorEditAccent)
				art.DrawRect(screen, boxX, lineY-11, gutterW, art.LineHeight, color.RGBA{64, 48, 16, 255})
			} else {
				// NAVIGATE MODE: deep solid navy background with electric blue border & left accent
				art.DrawRect(screen, selX, lineY-11, selW, art.LineHeight, art.ColorSelectBar)
				art.DrawBorder(screen, selX, lineY-11, selW, art.LineHeight, art.ColorSelectBorder)
				art.DrawRect(screen, selX, lineY-11, 4, art.LineHeight, art.ColorSelectAccent)
				art.DrawRect(screen, boxX, lineY-11, gutterW, art.LineHeight, color.RGBA{28, 45, 80, 255})
			}
		}

		// Line number and active pointer in gutter
		var lineNumColor color.Color = art.ColorLineNumber
		pointer := "  "
		if isSelected {
			if g.editor.IsActive {
				lineNumColor = art.ColorLineNumEdit
				pointer = "✎ "
			} else {
				lineNumColor = art.ColorLineNumActive
				pointer = "> "
			}
		}
		lineNumStr := fmt.Sprintf("%s%02d", pointer, i+1)
		art.DrawText(screen, lineNumStr, boxX+4, lineY, lineNumColor)

		// Draw line text
		textX := boxX + gutterW + 10
		if isSelected && g.editor.IsActive {
			// In edit mode: show current buffer text with cursor and [EDITING] badge
			bufStr := g.editor.CurrentBufferString()
			art.DrawHighlightedLine(screen, bufStr, textX, lineY)

			// Draw blinking terminal cursor (high-contrast pure white)
			if g.editor.CursorVisible {
				cursorX := textX + g.editor.CursorPos*art.CharWidth
				art.DrawRect(screen, cursorX, lineY-10, 7, 12, art.ColorCursor)
			}

			// Right-aligned edit badge
			art.DrawText(screen, "[EDITING]", boxX+codeBoxW-72, lineY, art.ColorEditAccent)
		} else {
			// In navigation mode: show normal syntax highlighted line
			art.DrawHighlightedLine(screen, g.WorkingLines[i], textX, lineY)
		}
	}

	// Draw Diagnostic Hint Panel or Status Footer
	art.DrawRect(screen, 0, 312, VirtualWidth, 48, art.ColorPanelBg)
	art.DrawRect(screen, 0, 312, VirtualWidth, 1, art.ColorPanelBorder)

	if g.FeedbackMsg != "" {
		art.DrawText(screen, g.FeedbackMsg, 20, 335, art.ColorRedPanic)
	} else if g.editor.IsActive {
		art.DrawText(screen, "[EDITING LINE] Type replacement fix, press [ENTER] to submit, [ESC] to cancel", 20, 335, art.ColorYellowAlert)
	} else if g.editor.HintVisible {
		art.DrawText(screen, g.CurrentLevel.Hint, 20, 332, art.ColorGreenRecover)
		art.DrawText(screen, fmt.Sprintf("Error: %s", g.CurrentLevel.PanicMessage), 20, 348, art.ColorYellowAlert)
	} else {
		modeText := "[NAVIGATE] [UP/DOWN] Select line | [ENTER] Edit line | [H] Diagnostic Hint"
		art.DrawText(screen, modeText, 20, 335, art.ColorComment)
	}
}

func (g *Game) drawRecovered(screen *ebiten.Image) {
	art.DrawRect(screen, 80, 60, 480, 240, art.ColorPanelBg)
	art.DrawBorder(screen, 80, 60, 480, 240, art.ColorGreenRecover)

	art.DrawText(screen, "==================================================", 100, 90, art.ColorGreenRecover)
	art.DrawText(screen, "       * * *   P A N I C   R E C O V E R E D !   * * *", 100, 110, art.ColorGreenRecover)
	art.DrawText(screen, "==================================================", 100, 130, art.ColorGreenRecover)

	art.DrawText(screen, fmt.Sprintf("Bug Cleared: %s", g.CurrentLevel.Title), 100, 160, art.ColorYellowAlert)
	art.DrawText(screen, g.CurrentLevel.Explanation, 100, 180, art.ColorTextWhite)

	bonusStr := fmt.Sprintf("Score: %05d  (+%d Time Bonus)", g.Score, int(g.TimeRemaining*50))
	art.DrawText(screen, bonusStr, 100, 215, art.ColorCyanGlow)

	art.DrawText(screen, ">>> PRESS [ENTER] OR [SPACE] FOR NEXT LEVEL <<<", 140, 265, art.ColorGreenRecover)
}

func (g *Game) drawPanicCrash(screen *ebiten.Image, sx, sy int) {
	art.DrawRect(screen, 50, 40, 540, 280, color.RGBA{30, 0, 0, 245})
	art.DrawBorder(screen, 50, 40, 540, 280, art.ColorRedPanic)

	art.DrawText(screen, "FATAL CRASH DUMP: RUNTIME PANIC!", 70+sx, 65+sy, art.ColorRedPanic)
	art.DrawText(screen, "----------------------------------------------------", 70+sx, 80+sy, art.ColorRedPanic)

	art.DrawText(screen, g.CurrentLevel.PanicMessage, 70+sx, 105+sy, art.ColorYellowAlert)
	art.DrawText(screen, "goroutine 1 [running]:", 70+sx, 125+sy, art.ColorComment)
	art.DrawText(screen, "main.executeHandler(0x104b2, 0x14000108040)", 90+sx, 145+sy, art.ColorComment)
	art.DrawText(screen, "    /gclatam2026/panic-recover/main.go:88 +0x3c", 110+sx, 165+sy, art.ColorComment)
	art.DrawText(screen, "exit status 2 (SIGABRT)", 70+sx, 195+sy, art.ColorRedPanic)

	art.DrawText(screen, fmt.Sprintf("Hint: %s", g.CurrentLevel.Hint), 70+sx, 230+sy, art.ColorGreenRecover)

	art.DrawText(screen, "[ENTER] Retry Level         [ESC] Return to Title", 130+sx, 285+sy, art.ColorTextWhite)
}

func (g *Game) drawGameWin(screen *ebiten.Image) {
	art.DrawRect(screen, 60, 40, 520, 280, art.ColorPanelBg)
	art.DrawBorder(screen, 60, 40, 520, 280, art.ColorCyanGlow)

	art.DrawText(screen, "==================================================", 80, 75, art.ColorCyanGlow)
	art.DrawText(screen, "   S Y S T E M   S T A B I L I Z E D !   ( 1 0 / 1 0 )", 80, 95, art.ColorGreenRecover)
	art.DrawText(screen, "==================================================", 80, 115, art.ColorCyanGlow)

	art.DrawText(screen, "Congratulations, Gopher! You recovered all 10 fatal panics!", 90, 145, art.ColorYellowAlert)
	art.DrawText(screen, "Production is safe from nil pointers, race conditions and leaks.", 80, 168, art.ColorTextWhite)

	finalScoreStr := fmt.Sprintf("FINAL SCORE: %06d", g.Score)
	art.DrawText(screen, finalScoreStr, 220, 205, art.ColorGreenRecover)

	rank := "Senior Go Runtime Engineer"
	if g.Score > 15000 {
		rank = "Staff Runtime Architect (Gopher Deity)"
	}
	art.DrawText(screen, fmt.Sprintf("EARNED TITLE: %s", rank), 150, 230, art.ColorCyanGlow)

	art.DrawText(screen, ">>> PRESS [ENTER] TO RETURN TO TITLE <<<", 170, 285, art.ColorYellowAlert)
}

// Layout implements ebiten.Game virtual canvas resolution (16:9 640x360).
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}
