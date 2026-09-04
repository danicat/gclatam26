package scenes

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"panic-at-the-disco/internal/audio"
	"panic-at-the-disco/internal/entities"
	"panic-at-the-disco/internal/gfx"
	"panic-at-the-disco/internal/input"
	"panic-at-the-disco/internal/levels"
)

type PlayScene struct {
	ae               *audio.AudioEngine
	levelCfg         levels.LevelConfig
	floor            *gfx.DiscoFloor
	player           *entities.Player
	exitDoor         *entities.ExitDoor
	puddles          []*entities.DrinkPuddle
	clubbers         []*entities.PanickedClubber
	discoBalls       []*entities.FallingDiscoBall
	trusses          []*entities.FallingTruss
	particles        *gfx.ParticleSystem
	collapseTimer    float64
	hazardSpawnTimer float64
	cameraShake      float64
	score            int
	timeSurvived     float64
	isFastBGMActive  bool
}

func NewPlayScene(zone levels.ZoneID, initialLives, initialScore int) *PlayScene {
	cfg := levels.GetLevelConfig(zone)
	floor, exit, puddles, clubbers := levels.SetupZoneEntities(cfg)

	player := entities.NewPlayer(cfg.PlayerStartX, cfg.PlayerStartY)
	if initialLives > 0 {
		player.Lives = initialLives
	}

	return &PlayScene{
		ae:            audio.GetAudioEngine(),
		levelCfg:      cfg,
		floor:         floor,
		player:        player,
		exitDoor:      exit,
		puddles:       puddles,
		clubbers:      clubbers,
		discoBalls:    make([]*entities.FallingDiscoBall, 0, 10),
		trusses:       make([]*entities.FallingTruss, 0, 5),
		particles:     gfx.NewParticleSystem(500),
		collapseTimer: cfg.Duration,
		score:         initialScore,
	}
}

func (ps *PlayScene) Enter() {
	ps.ae.PlayBGM(false)
	ps.isFastBGMActive = false
}

func (ps *PlayScene) Exit() {
	ps.ae.PauseBGM()
}

func (ps *PlayScene) Update(dt float64) SceneAction {
	ps.timeSurvived += dt
	ps.collapseTimer -= dt
	ps.score += int(dt * 10.0) // Passive survival points

	in := input.Poll()

	// 1. Check loss conditions
	if ps.player.Lives <= 0 {
		return SceneAction{
			Type:        ActionSwitchScene,
			NextScene:   SceneGameOver,
			LossReason:  "CRUSHED BY COLLAPSING ROOF!",
			Score:       ps.score,
			SurviveTime: ps.timeSurvived,
		}
	}
	if ps.collapseTimer <= 0 {
		return SceneAction{
			Type:        ActionSwitchScene,
			NextScene:   SceneGameOver,
			LossReason:  "TIME OUT! THE CLUB ROOF CAVED IN!",
			Score:       ps.score,
			SurviveTime: ps.timeSurvived,
		}
	}

	// 2. Fast tempo & sirens when under 12 seconds
	if ps.collapseTimer < 12.0 && !ps.isFastBGMActive {
		ps.ae.PlayBGM(true)
		ps.isFastBGMActive = true
	}

	// 3. Camera shake decay
	if ps.cameraShake > 0 {
		ps.cameraShake = math.Max(0, ps.cameraShake-dt*18.0)
	}

	// 4. Update Subsystems
	ps.floor.Update(dt)
	ps.particles.Update(dt)
	ps.exitDoor.Update(dt)

	// Playable field bounds
	fieldX, fieldY := 35.0, 45.0
	fieldW, fieldH := 570.0, 275.0

	// Player update
	ps.player.Update(dt, in, fieldX, fieldY, fieldW, fieldH, ps.particles, ps.ae)

	// Floor hazard check
	if ps.floor.IsHazardAt(ps.player.X, ps.player.Y) {
		if ps.player.TakeDamage(1, ps.particles, ps.ae) {
			ps.cameraShake = 10.0
		}
	}

	// Puddles update
	for _, pud := range ps.puddles {
		pud.Update(ps.player, ps.ae)
	}

	// Panicked clubbers update
	for _, c := range ps.clubbers {
		c.Update(dt, fieldX, fieldY, fieldW, fieldH, ps.player)
	}

	// 5. Hazard Spawner
	ps.hazardSpawnTimer += dt
	if ps.hazardSpawnTimer >= ps.levelCfg.HazardInterval {
		ps.hazardSpawnTimer = 0
		ps.spawnCeilingHazard()
	}

	// Update active disco balls
	activeBalls := ps.discoBalls[:0]
	for _, db := range ps.discoBalls {
		db.Update(dt, ps.player, ps.particles, ps.ae)
		if db.State == entities.HazardImpact && !db.DealtDamage {
			ps.cameraShake = 12.0
		}
		if db.State != entities.HazardFinished {
			activeBalls = append(activeBalls, db)
		}
	}
	ps.discoBalls = activeBalls

	// Update active trusses
	activeTrusses := ps.trusses[:0]
	for _, tr := range ps.trusses {
		tr.Update(dt, ps.player, ps.particles, ps.ae)
		if tr.State == entities.HazardImpact && !tr.DealtDamage {
			ps.cameraShake = 14.0
		}
		if tr.State != entities.HazardFinished {
			activeTrusses = append(activeTrusses, tr)
		}
	}
	ps.trusses = activeTrusses

	// 6. Check Win / Door reached
	if ps.exitDoor.IsPlayerInside(ps.player) {
		ps.ae.PlaySFXWin()
		ps.score += int(ps.collapseTimer * 20.0) // Bonus points for remaining time

		if ps.levelCfg.ID < levels.ZoneBackstage {
			// Next zone
			nextZone := ps.levelCfg.ID + 1
			return SceneAction{
				Type:        ActionSwitchScene,
				NextScene:   SceneClear,
				TargetZone:  nextZone,
				Score:       ps.score,
				Lives:       ps.player.Lives,
				SurviveTime: ps.timeSurvived,
			}
		}

		// Escaped Zone 3 -> Grand Victory!
		return SceneAction{
			Type:        ActionSwitchScene,
			NextScene:   SceneVictory,
			Score:       ps.score,
			Lives:       ps.player.Lives,
			SurviveTime: ps.timeSurvived,
		}
	}

	return SceneAction{Type: ActionNone}
}

func (ps *PlayScene) spawnCeilingHazard() {
	rnd := rand.New(rand.NewSource(int64(ps.player.X*100 + ps.player.Y*10 + ps.timeSurvived*1000)))

	// 1. Falling Disco Ball (Targets near player with slight spread)
	targetX := ps.player.X + (rnd.Float64()-0.5)*140.0
	targetY := ps.player.Y + (rnd.Float64()-0.5)*120.0
	// Clamp to floor
	targetX = math.Max(60.0, math.Min(580.0, targetX))
	targetY = math.Max(60.0, math.Min(290.0, targetY))

	warningTime := 1.2
	if ps.collapseTimer < 15.0 {
		warningTime = 0.85
	}
	radius := 26.0 + rnd.Float64()*12.0

	ps.discoBalls = append(ps.discoBalls, entities.NewFallingDiscoBall(targetX, targetY, radius, warningTime))

	// 2. In Zone 2 and 3: Occasionally spawn falling trusses
	if ps.levelCfg.ID >= levels.ZoneVIPLounge && rnd.Float64() < 0.45 {
		trussX := 50.0 + rnd.Float64()*440.0
		trussY := 60.0 + rnd.Float64()*220.0
		trussW := 90.0 + rnd.Float64()*50.0
		ps.trusses = append(ps.trusses, entities.NewFallingTruss(trussX, trussY, trussW, 14.0, warningTime*1.1))
	}

	// 3. Trigger floor hazard tiles
	col := rnd.Intn(ps.floor.Cols)
	row := rnd.Intn(ps.floor.Rows)
	ps.floor.TriggerHazardTile(col, row, warningTime)
}

func (ps *PlayScene) Draw(screen *ebiten.Image) {
	// Screen shake translation
	var shakeX, shakeY float32
	if ps.cameraShake > 0 {
		shakeX = float32((rand.Float64() - 0.5) * ps.cameraShake * 1.5)
		shakeY = float32((rand.Float64() - 0.5) * ps.cameraShake * 1.5)
	}

	// Draw to temporary surface or offset
	_ = shakeX
	_ = shakeY

	// 1. Outer club walls
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{10, 8, 18, 255}, false)

	// 2. Disco dance floor
	ps.floor.Draw(screen)

	// 3. Spilled drink puddles
	for _, pud := range ps.puddles {
		pud.Draw(screen)
	}

	// 4. Emergency Exit Door
	ps.exitDoor.Draw(screen)

	// 5. Panicked Clubbers
	for _, c := range ps.clubbers {
		c.Draw(screen)
	}

	// 6. Falling Trusses
	for _, tr := range ps.trusses {
		tr.Draw(screen)
	}

	// 7. Falling Disco Balls
	for _, db := range ps.discoBalls {
		db.Draw(screen)
	}

	// 8. Player
	ps.player.Draw(screen)

	// 9. Particle VFX
	ps.particles.Draw(screen)

	// 10. HUD Overlay
	ps.drawHUD(screen)
}

func (ps *PlayScene) drawHUD(screen *ebiten.Image) {
	// Top status bar banner
	vector.DrawFilledRect(screen, 0, 0, 640, 38, color.RGBA{15, 12, 28, 240}, false)
	vector.StrokeLine(screen, 0, 38, 640, 38, 2.0, color.RGBA{255, 0, 128, 255}, false)

	// Zone Name
	gfx.DrawText(screen, ps.levelCfg.Name, 12.0, 6.0, 1.2, color.RGBA{0, 240, 255, 255}, true)

	// Lives (Disco Hearts)
	livesText := "LIVES: "
	for i := 0; i < ps.player.Lives; i++ {
		livesText += "<3 "
	}
	gfx.DrawText(screen, livesText, 12.0, 22.0, 1.2, color.RGBA{255, 50, 100, 255}, true)

	// Collapse Countdown Timer in Center
	timerCol := color.RGBA{255, 220, 0, 255}
	if ps.collapseTimer < 12.0 {
		timerCol = color.RGBA{255, 30, 30, 255} // Emergency red
	}
	timerStr := fmt.Sprintf("COLLAPSE IN: %04.1fS", math.Max(0, ps.collapseTimer))
	tw := gfx.MeasureText(timerStr, 1.6)
	gfx.DrawText(screen, timerStr, (640-tw)/2, 11.0, 1.6, timerCol, true)

	// Score at Right
	scoreStr := fmt.Sprintf("SCORE: %06d", ps.score)
	sw := gfx.MeasureText(scoreStr, 1.2)
	gfx.DrawText(screen, scoreStr, 628-sw, 6.0, 1.2, color.RGBA{255, 215, 0, 255}, true)

	// Panic & Groove Meters at Right
	panicStr := fmt.Sprintf("PANIC: %02d%%", int(ps.player.PanicLevel))
	gfx.DrawText(screen, panicStr, 628-sw, 22.0, 1.1, color.RGBA{255, 100, 100, 255}, true)

	// Groove Dash Indicator
	dashStr := "[DASH READY]"
	dashCol := color.RGBA{0, 255, 200, 255}
	if ps.player.GrooveMeter < 30.0 {
		dashStr = "[CHARGING...]"
		dashCol = color.RGBA{140, 140, 150, 200}
	}
	dw := gfx.MeasureText(dashStr, 1.1)
	gfx.DrawText(screen, dashStr, 510-dw, 22.0, 1.1, dashCol, true)

	// Bottom warning when timer is critical
	if ps.collapseTimer < 12.0 {
		blink := math.Sin(ps.timeSurvived * 14.0)
		if blink > 0 {
			warnText := "! ! ! CEILING IS COLLAPSING ! ! !"
			ww := gfx.MeasureText(warnText, 1.4)
			gfx.DrawText(screen, warnText, (640-ww)/2, 335.0, 1.4, color.RGBA{255, 20, 20, 255}, true)
		}
	}
}
