package scene

import (
	"image/color"
	"math"
	"math/rand"

	"exploding-kitchens/internal/entity"
	"exploding-kitchens/internal/game"
	"exploding-kitchens/internal/system"
	"exploding-kitchens/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

// PlayScene implements the primary 2-minute "Panic & Recover" kitchen gameplay.
type PlayScene struct {
	ctx        *game.GameContext
	chef       *entity.Chef
	stations   []*entity.Station
	toolRacks  []*entity.ToolRack
	cat        *entity.Cat
	particles  *system.ParticlePool
	hud        *ui.HUD
	drawOpts   ebiten.DrawImageOptions

	// Room bounds
	roomMinX, roomMaxX float64
	roomMinY, roomMaxY float64
	catRugX, catRugY   float64

	// Gameplay stats
	timeLeft       float64
	chaos          float64
	score          int
	clutches       int
	explosions     int
	isPaused       bool
	shakeTimer     float64
	shakeIntensity float64
	alarmTickTimer float64
}

// NewPlayScene creates a fresh instance of PlayScene.
func NewPlayScene() *PlayScene {
	return &PlayScene{}
}

// Enter initializes or resets all game entities for a new shift.
func (ps *PlayScene) Enter(ctx *game.GameContext) {
	ps.ctx = ctx
	ps.particles = system.NewParticlePool(300)
	ps.hud = ui.NewHUD(ctx.Font, ctx.PixelImg)

	// Room dimensions (Virtual 320x180)
	ps.roomMinX = 24.0
	ps.roomMaxX = 296.0
	ps.roomMinY = 24.0
	ps.roomMaxY = 160.0

	// Chef in center of room
	ps.chef = entity.NewChef(150, 90)

	// Kitchen Stations along top wall
	ps.stations = []*entity.Station{
		entity.NewStation(36, 26, entity.StationPressureCooker),
		entity.NewStation(96, 26, entity.StationStoveTop),
		entity.NewStation(160, 26, entity.StationDeepFryer),
		entity.NewStation(224, 26, entity.StationMicrowave),
	}

	// Tool Racks along bottom wall
	ps.toolRacks = []*entity.ToolRack{
		entity.NewToolRack(45, 140, entity.ToolIce),
		entity.NewToolRack(105, 140, entity.ToolExtinguisher),
		entity.NewToolRack(180, 140, entity.ToolWrench),
		entity.NewToolRack(240, 140, entity.ToolCatToy),
	}

	// Cat Resting Rug (bottom left)
	ps.catRugX = 26.0
	ps.catRugY = 85.0
	ps.cat = entity.NewCat(100, 100)

	// Reset shift variables
	ps.timeLeft = 120.0 // 2 minutes
	ps.chaos = 0.0
	ps.score = 0
	ps.clutches = 0
	ps.explosions = 0
	ps.isPaused = false
	ps.shakeTimer = 0
	ps.shakeIntensity = 0
	ps.alarmTickTimer = 0

	ctx.Audio.StartBGM()
}

// Update handles game logic, interaction detection, and state progression.
func (ps *PlayScene) Update(dt float64, in system.InputState) (game.SceneID, error) {
	// Attract Demo Mode interrupt
	if ps.ctx.DemoMode {
		if in.InteractJust || in.ConfirmJust || in.PauseJust || in.MoveX != 0 || in.MoveY != 0 {
			ps.ctx.DemoMode = false
			return game.SceneTitle, nil
		}
		ps.runDemoAI(dt)
	}

	// Pause toggle
	if in.PauseJust {
		ps.isPaused = !ps.isPaused
	}
	if ps.isPaused {
		return game.SceneKeepCurrent, nil
	}

	// 1. Shift Timer
	ps.timeLeft -= dt
	if ps.timeLeft <= 0 {
		ps.timeLeft = 0
		ps.saveSessionStats(true)
		return game.SceneGameOver, nil
	}

	// 2. Chef Movement & Dropping
	if !ps.ctx.DemoMode {
		ps.chef.Update(dt, in.MoveX, in.MoveY, ps.roomMinX, ps.roomMaxX, ps.roomMinY, ps.roomMaxY)

		if in.DropJust && ps.chef.Tool != entity.ToolNone {
			ps.hud.AddPopup("DROPPED", ps.chef.X-4, ps.chef.Y-8, color.RGBA{200, 200, 220, 255})
			ps.chef.Tool = entity.ToolNone
			ps.ctx.Audio.PlayPickup()
		}

		// 3. Player Interaction (Space / E)
		if in.InteractJust {
			ps.handleInteraction()
		}
	}

	// 4. Update Stations & Handle Detonations
	hasPanic := false
	hasWarning := false
	for _, stn := range ps.stations {
		exploded := stn.Update(dt)
		if exploded {
			ps.explosions++
			ps.chaos += 25.0
			ps.triggerScreenShake(0.4, 5.0)
			ps.ctx.Audio.PlayExplosion()
			ps.particles.EmitExplosion(stn.X+stn.W/2, stn.Y+stn.H/2, 35)
			ps.hud.AddPopup("EXPLOSION! +CHAOS", stn.X-10, stn.Y-8, color.RGBA{255, 40, 40, 255})
		}

		if stn.State == entity.StatePanic {
			hasPanic = true
			ps.particles.EmitSmoke(stn.X+stn.W/2, stn.Y+stn.H/2, 1)
			ps.particles.EmitFire(stn.X+stn.W/2, stn.Y+stn.H/2, 1)
		} else if stn.State == entity.StateWarning {
			hasWarning = true
			ps.particles.EmitSteam(stn.X+stn.W/2, stn.Y+stn.H/2, 1)
		}
	}

	// Periodic Alarm sound
	ps.alarmTickTimer += dt
	if hasPanic && ps.alarmTickTimer >= 0.4 {
		ps.ctx.Audio.PlayAlarm()
		ps.alarmTickTimer = 0
	} else if hasWarning && ps.alarmTickTimer >= 0.8 {
		ps.ctx.Audio.PlayWarning()
		ps.alarmTickTimer = 0
	}

	// 5. Update Cat & Chef collision with Cat
	ps.cat.Update(dt, ps.stations, ps.roomMinX+10, ps.roomMaxX-10, ps.roomMinY+15, ps.roomMaxY-15)

	cx, cy := ps.chef.Center()
	catDist := math.Hypot(cx-(ps.cat.X+ps.cat.W/2), cy-(ps.cat.Y+ps.cat.H/2))
	if catDist < 16.0 && ps.cat.State == entity.CatSittingOnStation {
		// Chef nudges cat off appliance
		ps.shooCat()
	}

	// 6. Update Screen Shake
	if ps.shakeTimer > 0 {
		ps.shakeTimer -= dt
		if ps.shakeTimer <= 0 {
			ps.shakeIntensity = 0
		}
	}

	// 7. Update Particles & HUD
	ps.particles.Update(dt)
	ps.hud.Update(dt)

	// 8. Chaos Bleed & Meltdown Condition
	if hasPanic {
		ps.chaos += 2.0 * dt // Unattended panic causes chaos to steadily rise
	}
	if ps.chaos > 100.0 {
		ps.chaos = 100.0
	}

	if ps.chaos >= 100.0 {
		ps.saveSessionStats(false)
		return game.SceneGameOver, nil
	}

	// 9. Adaptive Audio
	ps.ctx.Audio.SetPanicMode(hasPanic || ps.chaos >= 60.0)

	return game.SceneKeepCurrent, nil
}

func (ps *PlayScene) handleInteraction() {
	cx, cy := ps.chef.Center()

	// 1. Check Tool Racks
	for _, rack := range ps.toolRacks {
		rx := rack.X + rack.W/2
		ry := rack.Y + rack.H/2
		if math.Hypot(cx-rx, cy-ry) < 20.0 {
			ps.chef.Tool = rack.Tool
			ps.ctx.Audio.PlayPickup()
			ps.hud.AddPopup("EQUIPPED "+rack.Tool.Name(), ps.chef.X-15, ps.chef.Y-10, color.RGBA{180, 240, 255, 255})
			return
		}
	}

	// 2. Check Cat
	catDist := math.Hypot(cx-(ps.cat.X+ps.cat.W/2), cy-(ps.cat.Y+ps.cat.H/2))
	if catDist < 24.0 && (ps.cat.State == entity.CatSittingOnStation || ps.cat.State == entity.CatWandering) {
		ps.shooCat()
		return
	}

	// 3. Check Stations
	for _, stn := range ps.stations {
		sx := stn.X + stn.W/2
		sy := stn.Y + stn.H/2
		if math.Hypot(cx-sx, cy-sy) < 26.0 {
			if stn.CanRecover(ps.chef.Tool) {
				if stn.State == entity.StateExploded {
					stn.Repair()
					ps.score += 150
					ps.ctx.Audio.PlaySteam()
					ps.particles.EmitSteam(sx, sy, 10)
					ps.hud.AddPopup("REPAIRED! +150", sx-15, sy-8, color.RGBA{140, 220, 255, 255})
					return
				}

				if stn.IsClutch() {
					// CLUTCH RECOVERY!
					ps.clutches++
					ps.score += 500
					ps.chaos = math.Max(0, ps.chaos-20.0)
					ps.ctx.Audio.PlayClutch()
					ps.particles.EmitClutch(sx, sy, 25)
					ps.hud.AddPopup("CLUTCH DEFUSE! +500", sx-25, sy-10, color.RGBA{255, 240, 80, 255})
				} else {
					// Normal Recovery
					ps.score += 100
					ps.chaos = math.Max(0, ps.chaos-10.0)
					if ps.chef.Tool == entity.ToolExtinguisher {
						ps.ctx.Audio.PlayExtinguish()
					} else {
						ps.ctx.Audio.PlaySteam()
					}
					ps.particles.EmitSteam(sx, sy, 12)
					ps.hud.AddPopup("DEFUSED! +100", sx-15, sy-8, color.RGBA{100, 255, 120, 255})
				}

				stn.ResetCooking()
				return
			} else {
				// Tool mismatch notification
				ps.hud.AddPopup("WRONG TOOL!", sx-15, sy-8, color.RGBA{255, 100, 100, 255})
			}
		}
	}
}

func (ps *PlayScene) shooCat() {
	ps.cat.Shoo(ps.catRugX, ps.catRugY)
	ps.score += 50
	ps.ctx.Audio.PlayCat()
	ps.hud.AddPopup("CAT SHOOED! +50", ps.cat.X-15, ps.cat.Y-8, color.RGBA{255, 180, 80, 255})
}

func (ps *PlayScene) triggerScreenShake(duration, intensity float64) {
	ps.shakeTimer = duration
	ps.shakeIntensity = intensity
}

func (ps *PlayScene) runDemoAI(dt float64) {
	// Autonomous heuristic for Attract Demo mode
	var targetStn *entity.Station
	maxProg := -1.0
	for _, s := range ps.stations {
		if s.State != entity.StateExploded && s.Progress() > maxProg {
			maxProg = s.Progress()
			targetStn = s
		}
	}

	if targetStn != nil {
		// Run towards required tool or station
		tx := targetStn.X + targetStn.W/2
		ty := targetStn.Y + targetStn.H/2 + 15

		cx, cy := ps.chef.Center()
		dx := tx - cx
		dy := ty - cy
		dist := math.Hypot(dx, dy)

		if dist > 8.0 {
			ps.chef.Update(dt, dx/dist, dy/dist, ps.roomMinX, ps.roomMaxX, ps.roomMinY, ps.roomMaxY)
		} else {
			// In range: simulate defuse
			targetStn.ResetCooking()
			ps.particles.EmitClutch(targetStn.X+targetStn.W/2, targetStn.Y+targetStn.H/2, 10)
		}
	}
}

func (ps *PlayScene) saveSessionStats(survived bool) {
	ps.ctx.LastScore = ps.score
	ps.ctx.LastClutches = ps.clutches
	ps.ctx.LastExplosions = ps.explosions
	ps.ctx.LastSurvived = survived

	if ps.score > ps.ctx.HighScore {
		ps.ctx.HighScore = ps.score
	}
}

// Draw renders the kitchen room, furniture, entities, particles, and HUD.
func (ps *PlayScene) Draw(screen *ebiten.Image) {
	// Screen Shake Offset
	var offsetX, offsetY float64
	if ps.shakeTimer > 0 {
		offsetX = (rand.Float64()*2 - 1) * ps.shakeIntensity
		offsetY = (rand.Float64()*2 - 1) * ps.shakeIntensity
	}

	// 1. Kitchen Floor (Checkerboard pattern)
	tileSize := 16.0
	for y := 0.0; y < 180; y += tileSize {
		for x := 0.0; x < 320; x += tileSize {
			isEven := (int(x/tileSize) + int(y/tileSize)) % 2 == 0
			tileCol := color.RGBA{75, 70, 90, 255}
			if isEven {
				tileCol = color.RGBA{95, 90, 115, 255}
			}
			ps.drawRect(screen, x+offsetX, y+offsetY, tileSize, tileSize, tileCol)
		}
	}

	// 2. Kitchen Perimeter Walls
	wallCol := color.RGBA{35, 30, 48, 255}
	trimCol := color.RGBA{55, 48, 72, 255}
	// Top wall
	ps.drawRect(screen, 0+offsetX, 0+offsetY, 320, 24, wallCol)
	ps.drawRect(screen, 0+offsetX, 23+offsetY, 320, 1, trimCol)
	// Left wall
	ps.drawRect(screen, 0+offsetX, 24+offsetY, 20, 156, wallCol)
	ps.drawRect(screen, 19+offsetX, 24+offsetY, 1, 156, trimCol)
	// Right wall
	ps.drawRect(screen, 300+offsetX, 24+offsetY, 20, 156, wallCol)
	ps.drawRect(screen, 300+offsetX, 24+offsetY, 1, 156, trimCol)
	// Bottom wall
	ps.drawRect(screen, 0+offsetX, 164+offsetY, 320, 16, wallCol)
	ps.drawRect(screen, 0+offsetX, 164+offsetY, 320, 1, trimCol)

	// 3. Cat Rug (bottom left corner)
	ps.drawRect(screen, ps.catRugX+offsetX, ps.catRugY+offsetY, 20, 16, color.RGBA{180, 80, 100, 255})
	ps.drawRect(screen, ps.catRugX+2+offsetX, ps.catRugY+2+offsetY, 16, 12, color.RGBA{220, 110, 130, 255})

	// 4. Tool Racks
	for _, tr := range ps.toolRacks {
		tr.Draw(screen, ps.ctx.PixelImg)
	}

	// 5. Stations
	for _, stn := range ps.stations {
		stn.Draw(screen, ps.ctx.PixelImg)
	}

	// 6. Cat
	ps.cat.Draw(screen, ps.ctx.PixelImg)

	// 7. Chef
	ps.chef.Draw(screen, ps.ctx.PixelImg)

	// 8. VFX Particles
	ps.particles.Draw(screen)

	// 9. HUD Overlay
	ps.hud.Draw(screen, ps.score, ps.clutches, ps.timeLeft, ps.chaos, ps.chef.Tool)

	// 10. Pause or Demo Overlay
	if ps.isPaused {
		ps.drawRect(screen, 0, 0, 320, 180, color.RGBA{10, 10, 20, 180})
		ps.ctx.Font.DrawText(screen, "GAME PAUSED", 110, 80, 1.5, color.RGBA{255, 240, 120, 255}, true)
		ps.ctx.Font.DrawText(screen, "PRESS ESC OR P TO RESUME", 85, 105, 1.0, color.RGBA{200, 200, 220, 255}, true)
	}

	if ps.ctx.DemoMode {
		ps.ctx.Font.DrawText(screen, "DEMO MODE - PRESS ANY KEY TO PLAY", 55, 168, 1.0, color.RGBA{255, 255, 100, 255}, true)
	}
}

func (ps *PlayScene) drawRect(screen *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	ps.drawOpts.GeoM.Reset()
	ps.drawOpts.GeoM.Scale(w, h)
	ps.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	ps.drawOpts.ColorScale.Reset()
	ps.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(ps.ctx.PixelImg, &ps.drawOpts)
}

// Exit stops or fades active sound channels upon leaving the scene.
func (ps *PlayScene) Exit() {
	ps.ctx.Audio.SetPanicMode(false)
}
