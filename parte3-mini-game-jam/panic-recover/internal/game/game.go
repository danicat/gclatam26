package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"panic-recover/internal/art"
	"panic-recover/internal/audio"
	"panic-recover/internal/bg"
	"panic-recover/internal/entity"
	"panic-recover/internal/ui"
)

type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

const (
	ScreenWidth  = 640
	ScreenHeight = 360
)

type Game struct {
	state          GameState
	starfield      *bg.Starfield
	player         *entity.Player
	bulletMgr      *entity.BulletManager
	enemyMgr       *entity.EnemyManager
	pickupMgr      *entity.PickupManager
	particleSys    *entity.ParticleSystem
	boss           *entity.Boss
	difficultyTier int
	nextBossScore  int
	tierBannerText string
	tierBannerTime float64

	screenShake float64
	panicFlash  float64
	pixelImg    *ebiten.Image
}

func NewGame() *Game {
	// Initialize core systems
	art.Init()
	audio.Init()
	ui.Init()

	dot := ebiten.NewImage(1, 1)
	dot.Fill(color.White)

	g := &Game{
		state:          StateTitle,
		starfield:      bg.NewStarfield(ScreenWidth, ScreenHeight),
		bulletMgr:      entity.NewBulletManager(250),
		enemyMgr:       entity.NewEnemyManager(80),
		pickupMgr:      entity.NewPickupManager(50),
		particleSys:    entity.NewParticleSystem(300),
		boss:           entity.NewBoss(),
		difficultyTier: 1,
		nextBossScore:  1500,
		pixelImg:       dot,
	}
	g.resetPlayState()
	return g
}

func (g *Game) resetPlayState() {
	g.player = entity.NewPlayer(ScreenWidth/2.0, ScreenHeight-60.0)
	g.bulletMgr = entity.NewBulletManager(250)
	g.enemyMgr = entity.NewEnemyManager(80)
	g.pickupMgr = entity.NewPickupManager(50)
	g.particleSys = entity.NewParticleSystem(300)
	g.boss = entity.NewBoss()
	g.difficultyTier = 1
	g.nextBossScore = 1500
	g.tierBannerText = ""
	g.tierBannerTime = 0
	g.screenShake = 0
	g.panicFlash = 0
	g.starfield.SetSpeedMultiplier(1.0)
	g.enemyMgr.SetTier(1)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// Fullscreen toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Update background unless paused
	if g.state != StatePaused {
		g.starfield.Update(dt)
	}

	switch g.state {
	case StateTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.resetPlayState()
			g.state = StatePlaying
		}

	case StatePlaying:
		// ESC pauses the game
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePaused
			return nil
		}
		g.updatePlaying(dt)

	case StatePaused:
		// ESC resumes the game
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying
			return nil
		}
		// R restarts the run
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.resetPlayState()
			g.state = StatePlaying
			return nil
		}
		// Q returns to title
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.resetPlayState()
			g.state = StateTitle
			return nil
		}

	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.resetPlayState()
			g.state = StatePlaying
		}
	}

	return nil
}

func (g *Game) updatePlaying(dt float64) {
	// Screen shake decay
	if g.screenShake > 0 {
		g.screenShake -= dt * 18.0
		if g.screenShake < 0 {
			g.screenShake = 0
		}
	}

	if g.tierBannerTime > 0 {
		g.tierBannerTime -= dt
	}

	// Trigger Boss Spawn
	if !g.boss.Active && g.player.Score >= g.nextBossScore {
		g.boss.Spawn(g.difficultyTier, ScreenWidth)
		g.screenShake = 14.0
		audio.PlayPanicSiren()
		g.tierBannerText = fmt.Sprintf("⚠️ CRITICAL ALERT: %s DETECTED! ⚠️", g.boss.Name)
		g.tierBannerTime = 3.5
	}

	// Handle Player Input
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1
	}

	// Gamepad support
	gamepadIDs := ebiten.AppendGamepadIDs(nil)
	if len(gamepadIDs) > 0 {
		gpID := gamepadIDs[0]
		if ebiten.IsStandardGamepadAxisAvailable(gpID, ebiten.StandardGamepadAxisLeftStickHorizontal) {
			ax := ebiten.StandardGamepadAxisValue(gpID, ebiten.StandardGamepadAxisLeftStickHorizontal)
			ay := ebiten.StandardGamepadAxisValue(gpID, ebiten.StandardGamepadAxisLeftStickVertical)
			if math.Abs(ax) > 0.2 {
				dx += ax
			}
			if math.Abs(ay) > 0.2 {
				dy += ay
			}
		}
		if ebiten.IsStandardGamepadButtonPressed(gpID, ebiten.StandardGamepadButtonLeftLeft) {
			dx -= 1
		}
		if ebiten.IsStandardGamepadButtonPressed(gpID, ebiten.StandardGamepadButtonLeftRight) {
			dx += 1
		}
		if ebiten.IsStandardGamepadButtonPressed(gpID, ebiten.StandardGamepadButtonLeftTop) {
			dy -= 1
		}
		if ebiten.IsStandardGamepadButtonPressed(gpID, ebiten.StandardGamepadButtonLeftBottom) {
			dy += 1
		}
	}

	g.player.Move(dx, dy, dt)

	// Fire weapons (Autofire when holding button)
	isFiring := ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyJ) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if len(gamepadIDs) > 0 && ebiten.IsStandardGamepadButtonPressed(gamepadIDs[0], ebiten.StandardGamepadButtonRightBottom) {
		isFiring = true
	}

	if isFiring {
		g.player.TryShoot(g.bulletMgr)
	}

	// Deploy Recover EMP Bomb
	useBomb := inpututil.IsKeyJustPressed(ebiten.KeyK) ||
		inpututil.IsKeyJustPressed(ebiten.KeyE) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	if len(gamepadIDs) > 0 && inpututil.IsStandardGamepadButtonJustPressed(gamepadIDs[0], ebiten.StandardGamepadButtonRightRight) {
		useBomb = true
	}

	if useBomb {
		if g.player.UseRecoverBomb(g.particleSys, g.bulletMgr) {
			g.screenShake = 12.0
		}
	}

	// Update Player
	g.player.Update(dt, ScreenWidth, ScreenHeight, g.particleSys)
	if g.player.InPanic {
		g.panicFlash += dt * 8.0
		g.screenShake = math.Max(g.screenShake, 3.0)
	}

	if g.player.IsDead {
		g.state = StateGameOver
		return
	}

	// Update Bullets
	g.bulletMgr.Update(dt, ScreenWidth, ScreenHeight)

	// Update Boss
	if g.boss.Active {
		g.boss.Update(dt, ScreenWidth, g.bulletMgr, g.enemyMgr, g.particleSys, g.player.X, g.player.Y)

		// Boss ramming damage vs Player
		distToBoss := math.Hypot(g.player.X-g.boss.X, g.player.Y-g.boss.Y)
		if distToBoss <= g.player.HitRadius+g.boss.Radius {
			enteredPanic := g.player.TakeDamage(40, g.particleSys)
			g.screenShake = 16.0
			if enteredPanic {
				g.screenShake = 22.0
			}
		}
	}

	// Update Pickups
	g.pickupMgr.Update(dt, g.player.X, g.player.Y, ScreenHeight)

	// Update Enemies & Spawns (throttled when boss is active)
	_ = g.enemyMgr.Update(dt, ScreenWidth, ScreenHeight, g.bulletMgr, g.pickupMgr, g.particleSys, g.player.X, g.player.Y, g.player.InPanic, g.boss.Active)

	// Update Particles
	g.particleSys.Update(dt)

	// Bullet Collisions against Boss
	if g.boss.Active {
		bullets := g.bulletMgr.Bullets()
		for i := range bullets {
			b := &bullets[i]
			if !b.Active || !b.IsPlayer {
				continue
			}
			dist := math.Hypot(b.X-g.boss.X, b.Y-g.boss.Y)
			if dist <= g.boss.Radius+b.Width*0.5 {
				b.Active = false
				defeated := g.boss.TakeDamage(b.Damage, g.particleSys)
				if defeated {
					audio.PlayExplosion()
					g.screenShake = 20.0
					for k := 0; k < 8; k++ {
						g.particleSys.EmitExplosion(
							g.boss.X+(rand.Float64()-0.5)*40,
							g.boss.Y+(rand.Float64()-0.5)*30,
							25,
							color.RGBA{255, uint8(100 + rand.Intn(155)), 40, 255},
						)
					}
					// Guaranteed epic drops
					g.pickupMgr.Spawn(g.boss.X-25, g.boss.Y, entity.PickupTypeRecover)
					g.pickupMgr.Spawn(g.boss.X+25, g.boss.Y, entity.PickupTypeRecover)
					g.pickupMgr.Spawn(g.boss.X, g.boss.Y-15, entity.PickupTypeMutex)
					g.pickupMgr.Spawn(g.boss.X, g.boss.Y+15, entity.PickupTypeWorker)

					g.player.Score += 2500 * g.difficultyTier

					// Progression: harder and faster!
					g.difficultyTier++
					g.enemyMgr.SetTier(g.difficultyTier)
					g.starfield.SetSpeedMultiplier(1.0 + float64(g.difficultyTier-1)*0.20)
					g.nextBossScore = g.player.Score + 2000 + g.difficultyTier*1000
					g.tierBannerText = fmt.Sprintf("SYSTEM OVERCLOCKED! TIER %d: SPEED +%d%%, ENEMIES BUFFED!", g.difficultyTier, (g.difficultyTier-1)*18)
					g.tierBannerTime = 4.0
				}
				break
			}
		}
	}

	// Bullet Collisions against Enemies
	scoreGained := g.enemyMgr.HandleBulletCollisions(g.bulletMgr, g.pickupMgr, g.particleSys, g.player.InPanic)
	g.player.Score += scoreGained

	// Enemy Bullets vs Player
	bullets := g.bulletMgr.Bullets()
	for i := range bullets {
		b := &bullets[i]
		if !b.Active || b.IsPlayer {
			continue
		}

		dx := b.X - g.player.X
		dy := b.Y - g.player.Y
		distSq := dx*dx + dy*dy
		hitDist := g.player.HitRadius + b.Width*0.5

		if distSq <= hitDist*hitDist {
			b.Active = false
			enteredPanic := g.player.TakeDamage(b.Damage, g.particleSys)
			g.screenShake = 8.0
			if enteredPanic {
				g.screenShake = 16.0
			}
			break
		}
	}

	// Enemy Ship Ramming vs Player
	enemies := g.enemyMgr.Enemies()
	for i := range enemies {
		e := &enemies[i]
		if !e.Active {
			continue
		}

		dx := e.X - g.player.X
		dy := e.Y - g.player.Y
		distSq := dx*dx + dy*dy
		hitDist := g.player.HitRadius + e.Radius

		if distSq <= hitDist*hitDist {
			// Ramming deals 30 damage and destroys smaller enemies
			enteredPanic := g.player.TakeDamage(30, g.particleSys)
			g.screenShake = 10.0
			if enteredPanic {
				g.screenShake = 18.0
			}
			if e.Type != entity.EnemyTypeDeadlock {
				e.Active = false
				g.particleSys.EmitExplosion(e.X, e.Y, 15, color.RGBA{255, 100, 30, 255})
			}
		}
	}

	// Pickups vs Player
	pickups := g.pickupMgr.Pickups()
	for i := range pickups {
		p := &pickups[i]
		if !p.Active {
			continue
		}

		dx := p.X - g.player.X
		dy := p.Y - g.player.Y
		distSq := dx*dx + dy*dy
		if distSq <= (g.player.HitRadius+12)*(g.player.HitRadius+12) {
			p.Active = false
			g.player.CollectPickup(p, g.particleSys, g.bulletMgr)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Screen shake transform
	var shakeX, shakeY float64
	if g.screenShake > 0 {
		shakeX = (rand.Float64() - 0.5) * g.screenShake
		shakeY = (rand.Float64() - 0.5) * g.screenShake
	}

	screen.Clear()

	// Draw starfield
	g.starfield.Draw(screen, ui.FontSource())

	switch g.state {
	case StateTitle:
		g.drawTitle(screen)
	case StatePlaying:
		g.drawPlaying(screen, shakeX, shakeY)
	case StatePaused:
		g.drawPlaying(screen, 0, 0)
		g.drawPaused(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	}

	// Retro CRT scanline overlay
	if art.ScanlineOverlay != nil {
		var op ebiten.DrawImageOptions
		screen.DrawImage(art.ScanlineOverlay, &op)
	}
}

func (g *Game) drawPlaying(screen *ebiten.Image, shakeX, shakeY float64) {
	// Game world container with shake
	worldImg := ebiten.NewImage(ScreenWidth, ScreenHeight)

	// 1. Pickups
	g.pickupMgr.Draw(worldImg)

	// 2. Enemies
	g.enemyMgr.Draw(worldImg)

	// 2b. Boss
	g.boss.Draw(worldImg)

	// 3. Bullets
	g.bulletMgr.Draw(worldImg)

	// 4. Player & Drones
	g.player.Draw(worldImg)

	// 5. Particles
	g.particleSys.Draw(worldImg)

	// Render world with screen shake
	var worldOp ebiten.DrawImageOptions
	worldOp.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(worldImg, &worldOp)

	// Panic Mode Red Screen Vignette Flash
	if g.player.InPanic {
		alpha := uint8((math.Sin(g.panicFlash)*0.5 + 0.5) * 80.0)
		flashImg := ebiten.NewImage(ScreenWidth, ScreenHeight)
		flashImg.Fill(color.RGBA{255, 20, 20, alpha})
		screen.DrawImage(flashImg, nil)
	}

	// 6. Draw HUD
	g.drawHUD(screen)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	// Top Header Bar Background
	headerBar := ebiten.NewImage(ScreenWidth, 26)
	headerBar.Fill(color.RGBA{10, 15, 28, 220})
	screen.DrawImage(headerBar, nil)

	// Score
	scoreStr := fmt.Sprintf("SCORE: %06d", g.player.Score)
	ui.DrawText(screen, scoreStr, 12, 6, 11, color.RGBA{180, 230, 255, 255})

	// Tier
	tierStr := fmt.Sprintf("TIER %d", g.difficultyTier)
	ui.DrawText(screen, tierStr, 125, 6, 11, color.RGBA{255, 215, 80, 255})

	// Health Bar
	ui.DrawText(screen, "HP:", 185, 6, 11, color.RGBA{200, 220, 255, 255})
	hpBarWidth := 85.0
	hpRatio := math.Max(0.0, g.player.HP/g.player.MaxHP)

	// Bar border
	borderImg := ebiten.NewImage(int(hpBarWidth)+4, 12)
	borderImg.Fill(color.RGBA{40, 60, 90, 255})
	var bOp ebiten.DrawImageOptions
	bOp.GeoM.Translate(210, 7)
	screen.DrawImage(borderImg, &bOp)

	// Bar fill
	if hpRatio > 0 {
		fillWidth := int(hpBarWidth * hpRatio)
		if fillWidth > 0 {
			barFill := ebiten.NewImage(fillWidth, 8)
			var hpColor color.RGBA
			if hpRatio > 0.5 {
				hpColor = color.RGBA{50, 225, 120, 255}
			} else if hpRatio > 0.25 {
				hpColor = color.RGBA{240, 200, 40, 255}
			} else {
				hpColor = color.RGBA{240, 50, 50, 255}
			}
			barFill.Fill(hpColor)
			var fOp ebiten.DrawImageOptions
			fOp.GeoM.Translate(212, 9)
			screen.DrawImage(barFill, &fOp)
		}
	}

	// Recover Bombs stocked
	recStr := fmt.Sprintf("RECOVERS: %d", g.player.RecoverStock)
	ui.DrawText(screen, recStr, 312, 6, 11, color.RGBA{80, 255, 160, 255})

	// Active Shield / Drones
	if g.player.ShieldTimer > 0 {
		shieldStr := fmt.Sprintf("MUTEX: %.1fs", g.player.ShieldTimer)
		ui.DrawText(screen, shieldStr, 435, 6, 11, color.RGBA{255, 215, 50, 255})
	} else if g.player.DroneCount > 0 {
		droneStr := fmt.Sprintf("WORKERS: %d", g.player.DroneCount)
		ui.DrawText(screen, droneStr, 435, 6, 11, color.RGBA{40, 220, 255, 255})
	}

	// Pause shortcut hint
	ui.DrawText(screen, "[ESC: PAUSE]", 540, 6, 11, color.RGBA{160, 180, 210, 255})

	// Boss Health Bar UI
	if g.boss.Active {
		bossBarW := 240.0
		bossBarX := (ScreenWidth - bossBarW) / 2.0
		bossBarY := 30.0

		// Name
		ui.DrawCenteredText(screen, g.boss.Name, ScreenWidth/2.0, bossBarY, 10, color.RGBA{255, 110, 110, 255})

		// Frame
		frameImg := ebiten.NewImage(int(bossBarW)+4, 10)
		frameImg.Fill(color.RGBA{55, 15, 20, 230})
		var fOp ebiten.DrawImageOptions
		fOp.GeoM.Translate(bossBarX-2, bossBarY+14)
		screen.DrawImage(frameImg, &fOp)

		// Bar
		bRatio := math.Max(0.0, g.boss.HP/g.boss.MaxHP)
		if bRatio > 0 {
			fW := int(bossBarW * bRatio)
			if fW > 0 {
				bossFill := ebiten.NewImage(fW, 6)
				bossFill.Fill(color.RGBA{255, 45, 55, 255})
				var fillOp ebiten.DrawImageOptions
				fillOp.GeoM.Translate(bossBarX, bossBarY+16)
				screen.DrawImage(bossFill, &fillOp)
			}
		}
	}

	// Tier Notification / Alert Banner
	if g.tierBannerTime > 0 {
		bannerY := 48.0
		if g.boss.Active {
			bannerY = 60.0
		}
		ui.DrawCenteredText(screen, g.tierBannerText, ScreenWidth/2.0, bannerY, 11, color.RGBA{255, 230, 80, 255})
	}

	// PANIC MODE BANNER (CRITICAL UI)
	if g.player.InPanic {
		bannerY := 32.0
		if g.boss.Active {
			bannerY = 64.0
		}
		bannerBg := ebiten.NewImage(ScreenWidth, 30)
		bannerBg.Fill(color.RGBA{180, 15, 15, 230})
		var bOp ebiten.DrawImageOptions
		bOp.GeoM.Translate(0, bannerY)
		screen.DrawImage(bannerBg, &bOp)

		panicMsg := fmt.Sprintf("⚠️ PANIC: FATAL ERROR! KILL ENEMY FOR recover()! [ %.1fs ] ⚠️", g.player.PanicTimer)
		ui.DrawCenteredText(screen, panicMsg, ScreenWidth/2.0, bannerY+7, 13, color.RGBA{255, 255, 160, 255})
	}
}

func (g *Game) drawPaused(screen *ebiten.Image) {
	// Semi-transparent backdrop to freeze world in place visually
	dim := ebiten.NewImage(ScreenWidth, ScreenHeight)
	dim.Fill(color.RGBA{5, 10, 22, 215})
	screen.DrawImage(dim, nil)

	// Border outline
	border := ebiten.NewImage(504, 264)
	border.Fill(color.RGBA{0, 173, 216, 120})
	var bOp ebiten.DrawImageOptions
	bOp.GeoM.Translate(68, 48)
	screen.DrawImage(border, &bOp)

	// Cyber modal container
	panel := ebiten.NewImage(500, 260)
	panel.Fill(color.RGBA{15, 24, 40, 245})
	var pOp ebiten.DrawImageOptions
	pOp.GeoM.Translate(70, 50)
	screen.DrawImage(panel, &pOp)

	// Header
	ui.DrawCenteredText(screen, "=== RUNTIME PAUSED: EXECUTION HALTED ===", ScreenWidth/2.0, 68, 14, color.RGBA{80, 220, 255, 255})
	ui.DrawCenteredText(screen, "COMMAND & CONTROL PROTOCOL", ScreenWidth/2.0, 88, 10, color.RGBA{140, 175, 215, 255})

	// Command list
	y := 112.0
	spacing := 19.0

	cmds := []struct {
		Action string
		Key    string
	}{
		{"Mover Nave Gopher", "W, A, S, D  ou  Setas"},
		{"Disparo Laser (Autofire)", "Espaço / J / Clique Esquerdo"},
		{"Ativar Recover EMP (Bomba)", "K / E / Clique Direito"},
		{"Pausar / Retomar Jogo", "ESC"},
		{"Tela Cheia (Fullscreen)", "F11 / Alt + Enter"},
		{"Reiniciar Execução", "R"},
		{"Retornar ao Menu Principal", "Q"},
	}

	for _, c := range cmds {
		ui.DrawText(screen, c.Action, 95, y, 10, color.RGBA{220, 235, 250, 255})
		ui.DrawText(screen, c.Key, 345, y, 10, color.RGBA{255, 215, 80, 255})
		y += spacing
	}

	// Bottom action bar
	ui.DrawCenteredText(screen, "[ ESC: RETOMAR ]     [ R: REINICIAR ]     [ Q: MENU PRINCIPAL ]", ScreenWidth/2.0, 278, 11, color.RGBA{80, 255, 160, 255})
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	// Dark translucent backdrop panel
	panel := ebiten.NewImage(540, 280)
	panel.Fill(color.RGBA{12, 18, 32, 230})
	var pOp ebiten.DrawImageOptions
	pOp.GeoM.Translate(50, 40)
	screen.DrawImage(panel, &pOp)

	// Title
	ui.DrawCenteredText(screen, "PANIC RECOVER: RUNTIME DEFENDER", ScreenWidth/2.0, 55, 18, color.RGBA{70, 210, 255, 255})
	ui.DrawCenteredText(screen, "GopherCon LATAM 2026 Mini Game Jam Edition", ScreenWidth/2.0, 80, 11, color.RGBA{160, 190, 220, 255})

	// Enemy roster preview list
	ui.DrawText(screen, "ENEMY THREAT ROSTER (GO RUNTIME BUGS):", 75, 110, 11, color.RGBA{255, 215, 80, 255})
	ui.DrawText(screen, "• nil pointer dereference : Ultra-swift zigzag darts", 85, 130, 10, color.RGBA{210, 140, 255, 255})
	ui.DrawText(screen, "• concurrent map writes   : Dual-hull interceptors with crossfire", 85, 146, 10, color.RGBA{255, 120, 120, 255})
	ui.DrawText(screen, "• deadlock                : Armored bronze bunker with 3-way spread", 85, 162, 10, color.RGBA{255, 200, 100, 255})
	ui.DrawText(screen, "• memory leak             : Expands & divides into mini-leaks on death", 85, 178, 10, color.RGBA{120, 255, 140, 255})
	ui.DrawText(screen, "• goroutine leak          : Fast dive-bombing drone swarms", 85, 194, 10, color.RGBA{255, 160, 80, 255})

	// Core twist explanation
	ui.DrawCenteredText(screen, "SURGE HOOK: Taking fatal damage triggers 5s PANIC MODE!", ScreenWidth/2.0, 225, 11, color.RGBA{255, 80, 80, 255})
	ui.DrawCenteredText(screen, "Defeat an enemy or grab a recover() drop before time expires to survive!", ScreenWidth/2.0, 240, 10, color.RGBA{180, 255, 200, 255})

	// Prompt to start
	flash := int(g.starfield.UpdateCount()) / 30 // blink
	_ = flash
	ui.DrawCenteredText(screen, "[ PRESS SPACE OR ENTER TO LAUNCH RUNTIME ]", ScreenWidth/2.0, 280, 13, color.RGBA{255, 255, 255, 255})
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Dim background
	dim := ebiten.NewImage(ScreenWidth, ScreenHeight)
	dim.Fill(color.RGBA{10, 5, 15, 230})
	screen.DrawImage(dim, nil)

	// Authentic Go Stack Trace Crash Screen!
	ui.DrawText(screen, "panic: runtime error: fatal error: call stack unwound (0 recovers left)", 35, 45, 11, color.RGBA{255, 70, 70, 255})
	ui.DrawText(screen, "goroutine 1 [running]:", 35, 70, 10, color.RGBA{180, 180, 180, 255})
	ui.DrawText(screen, "runtime/panic.go:838 +0x24a", 55, 86, 10, color.RGBA{140, 160, 200, 255})
	ui.DrawText(screen, "github.com/gclatam26/panic-recover/ship.UnwindStack()", 35, 106, 10, color.RGBA{180, 180, 180, 255})
	ui.DrawText(screen, "    /workspace/panic-recover/internal/entity/player.go:42 +0x18b", 55, 122, 10, color.RGBA{140, 160, 200, 255})
	ui.DrawText(screen, "main.DefendRuntime()", 35, 142, 10, color.RGBA{180, 180, 180, 255})
	ui.DrawText(screen, "    /workspace/panic-recover/cmd/game/main.go:28 +0x5f", 55, 158, 10, color.RGBA{140, 160, 200, 255})

	// Score Summary
	scoreMsg := fmt.Sprintf("DEFENSE SCORE: %d", g.player.Score)
	ui.DrawCenteredText(screen, scoreMsg, ScreenWidth/2.0, 220, 18, color.RGBA{255, 220, 80, 255})

	ui.DrawCenteredText(screen, "[ PRESS 'R' OR SPACE TO RECOVER & RESTART ]", ScreenWidth/2.0, 270, 14, color.RGBA{80, 255, 160, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
