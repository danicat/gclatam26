package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	VirtualWidth  = 480
	VirtualHeight = 270
)

type GameState int

const (
	StatePlaying GameState = iota
	StateLevelClear
	StateFainted
	StateGameWon
)

type Game struct {
	levels       []LevelDef
	currentIdx   int
	currentLevel *LevelState
	state        GameState

	offscreen     *ebiten.Image
	shaderOptions *ebiten.DrawRectShaderOptions
	drawOptions   *ebiten.DrawImageOptions

	particles      *ParticleSystem
	shakeIntensity float64
	shakeX, shakeY float64

	clockHandAngle float64
	totalTime      float32
}

func NewGame() (*Game, error) {
	initAssets()
	initAudio()
	if err := initShader(); err != nil {
		return nil, fmt.Errorf("failed to compile shader: %w", err)
	}

	levels := GetAllLevels()
	g := &Game{
		levels:        levels,
		currentIdx:    0,
		offscreen:     ebiten.NewImage(VirtualWidth, VirtualHeight),
		shaderOptions: &ebiten.DrawRectShaderOptions{},
		drawOptions:   &ebiten.DrawImageOptions{},
		particles:     NewParticleSystem(),
	}
	g.loadLevel(0)
	return g, nil
}

func (g *Game) loadLevel(idx int) {
	g.currentIdx = idx
	g.currentLevel = NewLevelState(g.levels[idx])
	g.state = StatePlaying
	g.shakeIntensity = 0
	g.clockHandAngle = 0
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0
	g.totalTime += float32(dt)

	// Screen shake decay (Ease-out)
	if g.shakeIntensity > 0.05 {
		g.shakeX = (rand.Float64()*2 - 1) * g.shakeIntensity
		g.shakeY = (rand.Float64()*2 - 1) * g.shakeIntensity
		g.shakeIntensity *= 0.85
	} else {
		g.shakeIntensity = 0
		g.shakeX = 0
		g.shakeY = 0
	}

	// Update smooth visual positions for entities
	g.currentLevel.UpdateVisuals(dt)
	g.particles.Update(dt)

	// Smoothly advance clock hand to match panic progress
	targetAngle := g.currentLevel.PanicPercent() * 2.0 * math.Pi
	g.clockHandAngle += (targetAngle - g.clockHandAngle) * 0.28

	// Toggle Fullscreen (F11)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Instant restart on 'R'
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.loadLevel(g.currentIdx)
		return nil
	}

	lvl := g.currentLevel
	roomPixelW := lvl.Width * TileSize
	roomPixelH := lvl.Height * TileSize
	originX := float64((VirtualWidth-roomPixelW)/2) + g.shakeX
	originY := float64(((VirtualHeight-roomPixelH)/2)+12) + g.shakeY

	// Particle ambience: Void mist from open holes
	if rand.Float64() < 0.15 {
		for y := 0; y < lvl.Height; y++ {
			for x := 0; x < lvl.Width; x++ {
				if lvl.Tiles[y][x] == TileHole && rand.Float64() < 0.4 {
					g.particles.EmitVoidMist(originX+float64(x*TileSize), originY+float64(y*TileSize))
				}
			}
		}
	}

	// Particle ambience: Panic wisps around Gopher when in panic zone
	if lvl.PanicPercent() >= 0.80 && rand.Float64() < 0.4 {
		g.particles.EmitPanicWisps(originX+lvl.PlayerVisualX, originY+lvl.PlayerVisualY)
	}

	switch g.state {
	case StatePlaying:
		var dx, dy int
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			dy = -1
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			dy = 1
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			dx = -1
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			dx = 1
		}

		if dx != 0 || dy != 0 {
			moved := g.currentLevel.Move(dx, dy)
			if moved {
				// Spawn VFX particles based on action event
				switch g.currentLevel.LastEvent {
				case EventPush:
					g.particles.EmitDust(originX+g.currentLevel.EventX, originY+g.currentLevel.EventY, 6)
					g.shakeIntensity = 2.0
				case EventHoleFilled:
					g.particles.EmitDust(originX+g.currentLevel.EventX, originY+g.currentLevel.EventY, 14)
					g.shakeIntensity = 4.5
				case EventRecover:
					g.particles.EmitSparkles(originX+g.currentLevel.EventX, originY+g.currentLevel.EventY, 16)
				case EventWin:
					g.particles.EmitSparkles(originX+g.currentLevel.EventX, originY+g.currentLevel.EventY, 20)
				case EventFaint:
					g.shakeIntensity = 6.0
				}
			}

			if g.currentLevel.Cleared {
				if g.currentIdx+1 >= len(g.levels) {
					g.state = StateGameWon
				} else {
					g.state = StateLevelClear
				}
			} else if g.currentLevel.Fainted {
				g.state = StateFainted
			}
		}

	case StateLevelClear:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.loadLevel(g.currentIdx + 1)
		}

	case StateFainted:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.loadLevel(g.currentIdx)
		}

	case StateGameWon:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.loadLevel(0)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Render all game content to the offscreen buffer
	g.offscreen.Fill(color.RGBA{R: 14, G: 15, B: 22, A: 255})
	g.drawGameWorld(g.offscreen)
	g.particles.Draw(g.offscreen)
	g.drawHUD(g.offscreen)

	// 2. Post-processing with Chromatic Aberration Kage shader (toned down)
	panicPercent := g.currentLevel.PanicPercent()
	var intensity float32 = 0.0

	// Chromatic aberration triggers gently at >= 80% panic
	if panicPercent >= 0.80 {
		raw := float32((panicPercent - 0.80) / 0.20)
		if raw > 1.0 {
			raw = 1.0
		}
		// Smooth soft curve
		intensity = raw * 0.8
	}

	if aberrationShader != nil && intensity > 0.001 {
		g.shaderOptions.Images[0] = g.offscreen
		g.shaderOptions.Uniforms = map[string]any{
			"Intensity": intensity,
			"Time":      g.totalTime,
		}
		screen.DrawRectShader(VirtualWidth, VirtualHeight, aberrationShader, g.shaderOptions)
	} else {
		g.drawOptions.GeoM.Reset()
		g.drawOptions.ColorScale.Reset()
		screen.DrawImage(g.offscreen, g.drawOptions)
	}
}

func (g *Game) drawGameWorld(dst *ebiten.Image) {
	lvl := g.currentLevel
	roomPixelW := lvl.Width * TileSize
	roomPixelH := lvl.Height * TileSize
	originX := float64((VirtualWidth-roomPixelW)/2) + g.shakeX
	originY := float64(((VirtualHeight-roomPixelH)/2)+12) + g.shakeY

	// Render Base Tiles
	for y := 0; y < lvl.Height; y++ {
		for x := 0; x < lvl.Width; x++ {
			px := originX + float64(x*TileSize)
			py := originY + float64(y*TileSize)
			t := lvl.Tiles[y][x]

			g.drawOptions.GeoM.Reset()
			g.drawOptions.ColorScale.Reset()
			g.drawOptions.GeoM.Translate(px, py)

			switch t {
			case TileWall:
				dst.DrawImage(imgWall, g.drawOptions)
			case TileHole:
				dst.DrawImage(imgHole, g.drawOptions)
			case TileHoleFilled:
				dst.DrawImage(imgHoleFilled, g.drawOptions)
			case TileClock:
				dst.DrawImage(imgFloor, g.drawOptions)
				dst.DrawImage(imgClock, g.drawOptions)
			case TileArtifact:
				dst.DrawImage(imgFloor, g.drawOptions)
				dst.DrawImage(imgArtifact, g.drawOptions)
			default:
				dst.DrawImage(imgFloor, g.drawOptions)
			}
		}
	}

	// Render Boulders with smooth visual interpolation
	for _, b := range lvl.Boulders {
		px := originX + b.VisualX
		py := originY + b.VisualY
		g.drawOptions.GeoM.Reset()
		g.drawOptions.ColorScale.Reset()
		g.drawOptions.GeoM.Translate(px, py)
		dst.DrawImage(imgBoulder, g.drawOptions)
	}

	// Render Gopher (Player) with smooth visual interpolation
	px := originX + lvl.PlayerVisualX
	py := originY + lvl.PlayerVisualY
	g.drawOptions.GeoM.Reset()
	g.drawOptions.ColorScale.Reset()
	g.drawOptions.GeoM.Translate(px, py)
	dst.DrawImage(imgGopher, g.drawOptions)
}

func (g *Game) drawHUD(dst *ebiten.Image) {
	lvl := g.currentLevel

	// Top Title Bar
	drawText(dst, lvl.Def.Name, 12, 10, color.RGBA{R: 240, G: 240, B: 255, A: 255}, 1)
	drawText(dst, lvl.Def.SubTitle, 12, 20, color.RGBA{R: 140, G: 150, B: 180, A: 255}, 1)

	// Clock Sprite ticking with time limit
	clockX := float64(VirtualWidth - 172)
	clockY := 5.0

	panicP := lvl.PanicPercent()

	// Pocket-watch body with subtle nervous pulsing when in panic mode
	g.drawOptions.GeoM.Reset()
	g.drawOptions.ColorScale.Reset()
	if panicP >= 0.80 {
		scale := 1.0 + float64(math.Sin(float64(g.totalTime*12.0)))*0.06
		g.drawOptions.GeoM.Translate(-12.0, -13.0)
		g.drawOptions.GeoM.Scale(scale, scale)
		g.drawOptions.GeoM.Translate(12.0, 13.0)
	}
	g.drawOptions.GeoM.Translate(clockX, clockY)
	dst.DrawImage(imgHUDClockFace, g.drawOptions)

	// Rotating Clock Hand (advances from 12 o'clock clockwise to doom)
	var jitter float64
	if panicP >= 0.80 {
		jitter = math.Sin(float64(g.totalTime*28.0)) * 0.1
	}
	g.drawOptions.GeoM.Reset()
	g.drawOptions.ColorScale.Reset()
	g.drawOptions.GeoM.Translate(-1.5, -8.5)
	g.drawOptions.GeoM.Rotate(g.clockHandAngle + jitter)
	g.drawOptions.GeoM.Translate(clockX+12.0, clockY+13.0)
	dst.DrawImage(imgHUDClockHand, g.drawOptions)

	// Turn Counter (Right side)
	turnsStr := fmt.Sprintf("TURNS: %d / %d", lvl.TurnsLeft, lvl.MaxTurns)
	drawText(dst, turnsStr, VirtualWidth-142, 10, color.RGBA{R: 255, G: 255, B: 255, A: 255}, 1)

	// Panic Bar Gauge
	barW := 125.0
	barH := 8.0
	barX := float64(VirtualWidth - 142)
	barY := 22.0

	// Bar Background via GPU rect
	drawRectGPU(dst, g.drawOptions, barX-1, barY-1, barW+2, barH+2, 0.14, 0.15, 0.20, 1.0)

	fillW := barW * panicP
	if fillW > barW {
		fillW = barW
	}

	// Bar Color depending on sanity state
	var cr, cg, cb float32 = 0.24, 0.78, 0.55 // Calm Green
	if panicP >= 0.80 {
		// Pulsing Panic Red
		pulse := float32(math.Sin(float64(g.totalTime*12.0))*0.15 + 0.85)
		cr, cg, cb = pulse, 0.12, 0.24
	} else if panicP >= 0.50 {
		cr, cg, cb = 0.94, 0.75, 0.16 // Nervous Yellow
	}

	if fillW > 0 {
		drawRectGPU(dst, g.drawOptions, barX, barY, fillW, barH, cr, cg, cb, 1.0)
	}

	// Panic percentage label or warning
	if panicP >= 0.80 {
		drawText(dst, "!! PANIC !!", int(barX)+30, int(barY)+1, color.RGBA{R: 255, G: 255, B: 255, A: 255}, 1)
	} else {
		panicLabel := fmt.Sprintf("PANIC: %d%%", int(panicP*100))
		drawText(dst, panicLabel, int(barX)+34, int(barY)+1, color.RGBA{R: 20, G: 20, B: 30, A: 255}, 1)
	}

	// Bottom Hotkey Reminders
	drawText(dst, "[WASD/ARROWS] MOVE/PUSH   [R] RESTART", 12, VirtualHeight-14, color.RGBA{R: 110, G: 115, B: 140, A: 255}, 1)

	// Overlay Modals
	switch g.state {
	case StateLevelClear:
		g.drawModal(dst, "ARTIFACT RECOVERED!", "PRESS [SPACE] FOR NEXT CHAMBER", 0.16, 0.71, 0.43)
	case StateFainted:
		g.drawModal(dst, "MADNESS OVERTAKES YOU!", "PRESS [R] OR [SPACE] TO RECOVER", 0.86, 0.18, 0.18)
	case StateGameWon:
		g.drawModal(dst, "ALL ARTIFACTS RECOVERED!", "THE GOPHER REMAINS SANE! PRESS [SPACE]", 0.94, 0.78, 0.16)
	}
}

func (g *Game) drawModal(dst *ebiten.Image, title, prompt string, ar, ag, ab float32) {
	modalW := 320.0
	modalH := 60.0
	mx := float64((VirtualWidth - int(modalW)) / 2)
	my := float64((VirtualHeight - int(modalH)) / 2)

	// Border outline
	drawRectGPU(dst, g.drawOptions, mx-2, my-2, modalW+4, modalH+4, ar, ag, ab, 1.0)
	// Dark translucent fill
	drawRectGPU(dst, g.drawOptions, mx, my, modalW, modalH, 0.07, 0.08, 0.12, 0.96)

	accentCol := color.RGBA{R: uint8(ar * 255), G: uint8(ag * 255), B: uint8(ab * 255), A: 255}
	tw := getTextWidth(title, 2)
	drawText(dst, title, (VirtualWidth-tw)/2, int(my)+12, accentCol, 2)

	pw := getTextWidth(prompt, 1)
	drawText(dst, prompt, (VirtualWidth-pw)/2, int(my)+38, color.RGBA{R: 220, G: 225, B: 240, A: 255}, 1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return VirtualWidth, VirtualHeight
}
