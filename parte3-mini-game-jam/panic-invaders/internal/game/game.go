package game

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"panic-invaders/internal/assets"
	"panic-invaders/internal/audio"
	"panic-invaders/internal/entity"
)

type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateGameOver
	StateVictory
)

type Star struct {
	X float64
	Y float64
	S float64
	C color.Color
}

type Game struct {
	State        GameState
	Player       *entity.Player
	Fleet        *entity.InvaderFleet
	Bullets      []*entity.Bullet
	Barriers     []*entity.Barrier
	Powerups     []*entity.Powerup
	Stars        []Star
	Wave         int
	MaxWaves     int
	LogMessage   string
	LogTimer     int
	ScreenW      int
	ScreenH      int
}

func NewGame() *Game {
	assets.InitSprites()
	audio.InitAudio()

	g := &Game{
		ScreenW:  640,
		ScreenH:  360,
		State:    StateTitle,
		MaxWaves: 3,
		Wave:     1,
	}
	g.initStars()
	return g
}

func (g *Game) initStars() {
	g.Stars = make([]Star, 70)
	colors := []color.Color{
		color.RGBA{0x40, 0x48, 0x59, 0xFF},
		color.RGBA{0x60, 0x70, 0x88, 0xFF},
		color.RGBA{0x00, 0xAD, 0xD8, 0x88},
		color.RGBA{0xC6, 0x78, 0xDD, 0x88},
	}
	for i := range g.Stars {
		g.Stars[i] = Star{
			X: rand.Float64() * float64(g.ScreenW),
			Y: rand.Float64() * float64(g.ScreenH),
			S: 0.3 + rand.Float64()*1.2,
			C: colors[rand.Intn(len(colors))],
		}
	}
}

func (g *Game) startNewGame() {
	g.Wave = 1
	g.Player = entity.NewPlayer(float64(g.ScreenW)/2-11, float64(g.ScreenH)-34)
	g.Fleet = entity.NewInvaderFleet(g.Wave)
	g.Bullets = nil
	g.Powerups = nil
	g.initBarriers()
	g.State = StatePlaying
	g.setLog("System initialized. defer recover() active.", 120)
}

func (g *Game) nextWave() {
	g.Wave++
	if g.Wave > g.MaxWaves {
		g.State = StateVictory
		return
	}
	g.Fleet = entity.NewInvaderFleet(g.Wave)
	g.Bullets = nil
	g.Powerups = nil
	g.initBarriers()
	g.setLog(fmt.Sprintf("WAVE %d: Call Stack unwinding intensified!", g.Wave), 150)
}

func (g *Game) initBarriers() {
	g.Barriers = []*entity.Barrier{
		entity.NewBarrier(120, 275, "defer log.Flush()"),
		entity.NewBarrier(280, 275, "defer db.Close()"),
		entity.NewBarrier(440, 275, "defer mu.Unlock()"),
	}
}

func (g *Game) setLog(msg string, duration int) {
	g.LogMessage = msg
	g.LogTimer = duration
}

func (g *Game) Update() error {
	// Update starfield
	for i := range g.Stars {
		g.Stars[i].Y += g.Stars[i].S
		if g.Stars[i].Y > float64(g.ScreenH) {
			g.Stars[i].Y = 0
			g.Stars[i].X = rand.Float64() * float64(g.ScreenW)
		}
	}

	if g.LogTimer > 0 {
		g.LogTimer--
	}

	switch g.State {
	case StateTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.startNewGame()
		}

	case StatePlaying:
		g.updatePlaying()

	case StateGameOver, StateVictory:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.startNewGame()
		}
	}

	return nil
}

func (g *Game) updatePlaying() {
	timeoutActive := g.Player.TimeoutTimer > 0

	// Update Player
	g.Player.Update(&g.Bullets)

	// Update Invaders
	stackBreached := g.Fleet.Update(&g.Bullets, &g.Powerups, timeoutActive)
	if stackBreached {
		g.State = StateGameOver
		audio.GlobalAudio.PlayExplosion()
		return
	}

	// Check if wave is cleared
	bossDefeated := (g.Fleet.Boss == nil || !g.Fleet.Boss.Active)
	if g.Fleet.AliveCount() == 0 && bossDefeated {
		g.nextWave()
		return
	}

	// Update Bullets
	for _, b := range g.Bullets {
		b.Update()
	}

	// Bullet Collisions with Barriers
	for _, b := range g.Bullets {
		if !b.Active {
			continue
		}
		for _, bar := range g.Barriers {
			if bar.CheckBulletCollision(b) {
				break
			}
		}
	}

	// Player Bullet Collisions with Invaders
	for _, b := range g.Bullets {
		if !b.Active || b.IsEnemy {
			continue
		}
		if g.Fleet.CheckBulletCollisions(b, g.Player, &g.Powerups) {
			g.setLog("handled: err != nil (panic caught by recover)", 60)
		}
	}

	// Enemy Bullet Collisions with Player
	for _, b := range g.Bullets {
		if !b.Active || !b.IsEnemy {
			continue
		}
		bx := b.X
		by := b.Y
		bw := b.Width
		bh := b.Height
		px := g.Player.X
		py := g.Player.Y
		pw := g.Player.Width
		ph := g.Player.Height

		if bx+bw >= px && bx <= px+pw && by+bh >= py && by <= py+ph {
			b.Active = false
			if g.Player.TakeDamage() {
				g.setLog("WARNING: Goroutine crashed! Life lost!", 90)
				if g.Player.Lives <= 0 {
					g.State = StateGameOver
					return
				}
			}
		}
	}

	// Clean inactive bullets
	activeBullets := g.Bullets[:0]
	for _, b := range g.Bullets {
		if b.Active {
			activeBullets = append(activeBullets, b)
		}
	}
	g.Bullets = activeBullets

	// Update Powerups
	for _, p := range g.Powerups {
		p.Update()
		if !p.Active {
			continue
		}
		// Check collision with player
		if p.X+p.Width >= g.Player.X && p.X <= g.Player.X+g.Player.Width &&
			p.Y+p.Height >= g.Player.Y && p.Y <= g.Player.Y+g.Player.Height {
			p.Active = false
			audio.GlobalAudio.PlayPowerup()

			switch p.Type {
			case entity.PowerupMutex:
				g.Player.ShieldTimer = 480 // 8s
				g.setLog("POWER-UP: sync.Mutex active! Exclusao mutua protegida.", 90)
			case entity.PowerupChan:
				g.Player.ChanTimer = 480 // 8s
				g.setLog("POWER-UP: chan struct{}! Rajada tripla assincrona.", 90)
			case entity.PowerupTimeout:
				g.Player.TimeoutTimer = 360 // 6s
				g.setLog("POWER-UP: context.WithTimeout! Invasores desacelerados.", 90)
			case entity.PowerupBadge:
				g.Player.Score += 500
				// Nuke all enemy bullets on screen!
				for _, eb := range g.Bullets {
					if eb.IsEnemy {
						eb.Active = false
					}
				}
				g.setLog("GOPHERCON LATAM BADGE! +500 PTS + Tela limpa!", 120)
			}
		}
	}

	// Clean inactive powerups
	activePowerups := g.Powerups[:0]
	for _, p := range g.Powerups {
		if p.Active {
			activePowerups = append(activePowerups, p)
		}
	}
	g.Powerups = activePowerups
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(assets.ColorBlack)

	// Stars
	for _, s := range g.Stars {
		ebitenutil.DrawRect(screen, s.X, s.Y, s.S, s.S, s.C)
	}

	switch g.State {
	case StateTitle:
		g.drawTitle(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	case StateVictory:
		g.drawVictory(screen)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "=======================================================================", 35, 30)
	ebitenutil.DebugPrintAt(screen, "                 GOPHERCON LATAM 2026: MINI GAME JAM                   ", 35, 45)
	ebitenutil.DebugPrintAt(screen, "                 PANIC INVADERS: IN RECOVER() WE TRUST                 ", 35, 60)
	ebitenutil.DebugPrintAt(screen, "=======================================================================", 35, 75)

	ebitenutil.DebugPrintAt(screen, "A tempestade de panic() ameaca derrubar o palco principal da GopherCon!", 45, 110)
	ebitenutil.DebugPrintAt(screen, "Assuma o controle de recover() na base da Call Stack e salve o evento!", 45, 125)

	// Hero preview
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 142)
	screen.DrawImage(assets.LoadedSprites.Player, op)
	ebitenutil.DebugPrintAt(screen, "= O Gopher Heroico (Azul/recover)       [HEROI]", 110, 146)

	// Bad Gophers preview
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 168)
	screen.DrawImage(assets.LoadedSprites.InvaderNil[0], op)
	ebitenutil.DebugPrintAt(screen, "= Bad Gopher: panic(\"nil pointer\")     [30 PTS]", 110, 168)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 192)
	screen.DrawImage(assets.LoadedSprites.InvaderIndex[0], op)
	ebitenutil.DebugPrintAt(screen, "= Bad Gopher: panic(\"index range\")     [20 PTS]", 110, 192)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 216)
	screen.DrawImage(assets.LoadedSprites.InvaderDivide[0], op)
	ebitenutil.DebugPrintAt(screen, "= Bad Gopher: panic(\"divide zero\")     [10 PTS]", 110, 216)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 240)
	screen.DrawImage(assets.LoadedSprites.UFO, op)
	ebitenutil.DebugPrintAt(screen, "= Bad Gopher Drone: panic(\"wifi\")     [250 PTS + DROP]", 110, 240)

	ebitenutil.DebugPrintAt(screen, "CONTROLES: [A] / [D] ou Setas = Mover  |  [ESPACO] = Disparar recover()", 55, 275)

	ebitenutil.DebugPrintAt(screen, ">>> PRESSIONE [ESPACO] OU [ENTER] PARA INICIAR A DEFESA <<<", 70, 315)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw Barriers
	for _, b := range g.Barriers {
		b.Draw(screen)
	}

	// Draw Fleet
	g.Fleet.Draw(screen)

	// Draw Powerups
	for _, p := range g.Powerups {
		p.Draw(screen)
	}

	// Draw Bullets
	for _, b := range g.Bullets {
		b.Draw(screen)
	}

	// Draw Player
	g.Player.Draw(screen)

	// Top HUD
	hudText := fmt.Sprintf("GOPHERCON LATAM 2026 | SCORE: %06d | UPTIME: 99.999%% | GOROUTINES: %d | WAVE: %d/%d",
		g.Player.Score, g.Player.Lives, g.Wave, g.MaxWaves)
	ebitenutil.DebugPrintAt(screen, hudText, 20, 8)

	// Power-up status badges
	statusX := 20
	if g.Player.ShieldTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[MUTEX SHIELD: %ds]", g.Player.ShieldTimer/60), statusX, 22)
		statusX += 140
	}
	if g.Player.ChanTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[CHAN TRIPLE: %ds]", g.Player.ChanTimer/60), statusX, 22)
		statusX += 130
	}
	if g.Player.TimeoutTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[SLOW-MO: %ds]", g.Player.TimeoutTimer/60), statusX, 22)
	}

	// Bottom log line (Call stack console feedback)
	if g.LogTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("> LOG: %s", g.LogMessage), 20, 342)
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 30, 40, 580, 280, color.RGBA{0x40, 0x10, 0x15, 0xEE})

	ebitenutil.DebugPrintAt(screen, "=======================================================================", 45, 55)
	ebitenutil.DebugPrintAt(screen, "                FATAL: ALL GOROUTINES ARE ASLEEP - CRASH               ", 45, 70)
	ebitenutil.DebugPrintAt(screen, "=======================================================================", 45, 85)

	ebitenutil.DebugPrintAt(screen, "exit status 2: panic() unhandled in main goroutine", 60, 115)
	ebitenutil.DebugPrintAt(screen, "goroutine 1 [running]:", 60, 130)
	ebitenutil.DebugPrintAt(screen, "  gophercon/latam2026/keynote.Broadcast(0x0, 0xbadf00d)", 60, 145)
	ebitenutil.DebugPrintAt(screen, "  gophercon/latam2026/main.main()", 60, 160)
	ebitenutil.DebugPrintAt(screen, "        /florianopolis/gophercon2026/stage.go:42 +0x1337", 60, 175)

	finalScore := 0
	if g.Player != nil {
		finalScore = g.Player.Score
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PONTUACAO FINAL DE ESTABILIDADE: %06d", finalScore), 60, 215)
	ebitenutil.DebugPrintAt(screen, "A transmissao da GopherCon caiu! Tente novamente com defer recover()!", 60, 235)

	ebitenutil.DebugPrintAt(screen, ">>> PRESSIONE [ESPACO] OU [R] PARA REINICIAR O RUNTIME <<<", 75, 280)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 30, 40, 580, 280, color.RGBA{0x0B, 0x38, 0x24, 0xEE})

	ebitenutil.DebugPrintAt(screen, "=======================================================================", 45, 55)
	ebitenutil.DebugPrintAt(screen, "           GOPHERCON LATAM 2026 SALVA! 100% UPTIME ALCANCADO!          ", 45, 70)
	ebitenutil.DebugPrintAt(screen, "=======================================================================", 45, 85)

	ebitenutil.DebugPrintAt(screen, "TODAS AS HORDAS DE PANIC() FORAM INTERCEPTADAS POR RECOVER()!", 60, 120)
	ebitenutil.DebugPrintAt(screen, "A Keynote de Abertura foi transmitida com sucesso absoluto!", 60, 140)
	ebitenutil.DebugPrintAt(screen, "O Auditório aplaude de pe o Gopher Defensor!", 60, 160)

	finalScore := 0
	if g.Player != nil {
		finalScore = g.Player.Score
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PONTUACAO FINAL DE HEROISMO GO: %06d", finalScore), 60, 200)
	ebitenutil.DebugPrintAt(screen, "Process finished with exit code 0 (Clean Shutdown)", 60, 230)

	ebitenutil.DebugPrintAt(screen, ">>> PRESSIONE [ESPACO] OU [ENTER] PARA JOGAR NOVAMENTE <<<", 80, 280)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenW, g.ScreenH
}
