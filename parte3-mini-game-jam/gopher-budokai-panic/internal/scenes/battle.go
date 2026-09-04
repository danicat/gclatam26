package scenes

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"gopher-budokai-panic/internal/audio"
	"gopher-budokai-panic/internal/entities"
	"gopher-budokai-panic/internal/gfx"
	"gopher-budokai-panic/internal/ui"
)

type BattleScene struct {
	player      *entities.Fighter
	cpu         *entities.Fighter
	ai          *entities.AIController
	hud         *ui.HUD
	arena       *gfx.Arena
	particles   *gfx.ParticlePool
	spriteCache *gfx.SpriteCache
	blasts      []*entities.KiBlast

	screenShake float64
	shakeAmp    float64
	stats       MatchStats
}

type MatchStats struct {
	RecoveriesCount int
	PanicsSuffered  int
	BeamsFired      int
	MaxCombo        int
}

func NewBattleScene() *BattleScene {
	bs := &BattleScene{
		player:      entities.NewFighter(0, gfx.FighterPlayer, 160, 200, false),
		cpu:         entities.NewFighter(1, gfx.FighterCPU, 480, 200, true),
		hud:         ui.NewHUD(),
		arena:       gfx.NewArena(),
		particles:   gfx.NewParticlePool(600),
		spriteCache: gfx.InitSprites(),
		blasts:      make([]*entities.KiBlast, 0, 32),
	}
	bs.ai = entities.NewAIController(bs.cpu)
	return bs
}

func (bs *BattleScene) Enter() {
	audio.Get().PlayBGM("battle")
}

func (bs *BattleScene) Exit() {
	audio.Get().StopBGM()
}

func (bs *BattleScene) Update(dt float64) Scene {
	w, h := 640.0, 360.0

	// Dynamic Screen Shake decay
	if bs.screenShake > 0 {
		bs.screenShake -= dt
		if bs.screenShake <= 0 {
			bs.screenShake = 0
			bs.shakeAmp = 0
		}
	}

	// Dynamic Soundtrack adaptation to the Game Jam Theme (Panic & Recover)
	if bs.player.Panic.IsPanicked || bs.cpu.Panic.IsPanicked {
		audio.Get().PlayBGM("panic")
	} else {
		audio.Get().PlayBGM("battle")
	}

	// 1. Player Input Processing
	bs.handlePlayerInput(dt)

	// 2. CPU AI Decision & Actions
	cpuBlast, _ := bs.ai.Update(dt, bs.player)
	if cpuBlast != nil {
		bs.blasts = append(bs.blasts, cpuBlast)
	}

	// 3. Update Fighters
	bs.player.Update(dt, bs.cpu.X, bs.cpu.Y, w, h)
	bs.cpu.Update(dt, bs.player.X, bs.player.Y, w, h)

	// 4. Update Ki Blasts
	for _, b := range bs.blasts {
		if b.Active {
			b.Update(dt, w, h)
			// Collision against Opponent
			var target *entities.Fighter
			if b.OwnerID == bs.player.ID {
				target = bs.cpu
			} else {
				target = bs.player
			}

			dist := math.Hypot(target.X-b.X, target.Y-b.Y)
			if dist < float64(b.Radius)+20.0 {
				b.Active = false
				target.TakeDamage(b.Damage, false, b.X)
				bs.particles.SpawnExplosion(b.X, b.Y, b.Color)
				bs.addShake(0.12, 2.5)
			}
		}
	}

	// 5. Check Super Beam Collisions & Beam Clash
	bs.handleBeamInteractions(dt)

	// 6. Spawn Visual Aura & Dash Particles
	bs.spawnFighterVFX(bs.player)
	bs.spawnFighterVFX(bs.cpu)

	// 7. Update Systems
	bs.arena.Update(dt)
	bs.particles.Update(dt)
	bs.hud.Update(dt)

	// 8. Clean dead blasts
	activeBlasts := bs.blasts[:0]
	for _, b := range bs.blasts {
		if b.Active {
			activeBlasts = append(activeBlasts, b)
		}
	}
	bs.blasts = activeBlasts

	// 9. Check Victory / Defeat Conditions
	if bs.player.Health <= 0 {
		return NewGameOverScene(false, bs.stats)
	}
	if bs.cpu.Health <= 0 {
		return NewGameOverScene(true, bs.stats)
	}

	return nil
}

func (bs *BattleScene) handlePlayerInput(dt float64) {
	p := bs.player

	// Flight Movement (WASD / Arrows)
	dirX, dirY := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dirX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dirX += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dirY -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dirY += 1.0
	}

	if dirX != 0 || dirY != 0 {
		mag := math.Hypot(dirX, dirY)
		p.StartMove(dirX/mag, dirY/mag)
	}

	// Dragon Dash (L)
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		p.StartDragonDash(bs.cpu.X, bs.cpu.Y)
	}

	// Melee Strike / Combo (J)
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		if p.StartMelee(bs.cpu.X, bs.cpu.Y) {
			// Check melee hit connection
			dist := math.Hypot(bs.cpu.X-p.X, bs.cpu.Y-p.Y)
			if dist < 50.0 {
				isKnockback := (p.ComboStep == 3)
				dmg := 35.0
				if isKnockback {
					dmg = 60.0
					bs.addShake(0.25, 6.0)
				} else {
					bs.addShake(0.08, 2.0)
				}
				bs.cpu.TakeDamage(dmg, isKnockback, p.X)
				bs.particles.SpawnHitSparks((p.X+bs.cpu.X)/2, (p.Y+bs.cpu.Y)/2, 8)
			}
		}
	}

	// Ki Blast (K)
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		blast := p.StartKiBlast(bs.cpu.X, bs.cpu.Y)
		if blast != nil {
			bs.blasts = append(bs.blasts, blast)
		}
	}

	// Space key: Multi-purpose (Charge Ki OR Mash to Recover!)
	if p.Panic.IsPanicked {
		// MASH TO RECOVER! (The Game Jam Mechanic)
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyShift) {
			if p.Panic.TryMashRecover() {
				// EXPLOSIVE KIAI RECOVERY!
				p.TriggerRecoverKiai(bs.cpu)
				bs.particles.SpawnRecoverWave(p.X, p.Y)
				bs.addShake(0.35, 8.0)
				bs.stats.RecoveriesCount++
			}
		}
	} else {
		// Normal Ki Charge when holding Space
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.State = entities.StateChargeKi
		}
	}

	// Super Beam (I or U)
	if inpututil.IsKeyJustPressed(ebiten.KeyI) || inpututil.IsKeyJustPressed(ebiten.KeyU) {
		beam := p.StartSuperBeam()
		if beam != nil {
			bs.stats.BeamsFired++
			bs.addShake(0.3, 4.0)
		}
	}

	// Vanish / Instant Transmission (Shift)
	if inpututil.IsKeyJustPressed(ebiten.KeyShift) && !p.Panic.IsPanicked {
		p.PerformVanish(bs.cpu.X, bs.cpu.Y, bs.cpu.FacingLeft)
		bs.particles.SpawnHitSparks(p.X, p.Y, 5)
	}
}

func (bs *BattleScene) handleBeamInteractions(dt float64) {
	pBeam := bs.player.ActiveBeam
	cBeam := bs.cpu.ActiveBeam

	// 1. Beam Clash Check
	if entities.CheckBeamClash(pBeam, cBeam) {
		// Both beams colliding!
		bs.addShake(0.05, 3.0)
		audio.Get().PlayClash()

		// Mash Space to push Beam Clash towards CPU!
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
			pBeam.ClashOffset += 14.0
			cBeam.ClashOffset += 14.0
		}
		// CPU AI counter-pushes
		pBeam.ClashOffset -= 40.0 * dt
		cBeam.ClashOffset -= 40.0 * dt

		// Check if player overwhelmed CPU in clash!
		if pBeam.ClashOffset > 140.0 {
			cBeam.State = entities.BeamStateDone
			bs.cpu.TakeDamage(180.0, true, bs.player.X)
			bs.particles.SpawnExplosion(bs.cpu.X, bs.cpu.Y, pBeam.Color)
			bs.addShake(0.4, 9.0)
		}
		return
	}

	// 2. Player Beam vs CPU
	if pBeam != nil && pBeam.State == entities.BeamStateFiring {
		beamY := pBeam.Y
		if math.Abs(bs.cpu.Y-beamY) < 32.0 {
			hitX := pBeam.X
			if (pBeam.DirX > 0 && bs.cpu.X > hitX) || (pBeam.DirX < 0 && bs.cpu.X < hitX) {
				bs.cpu.TakeDamage(120.0*dt, true, pBeam.X)
				bs.particles.SpawnHitSparks(bs.cpu.X, bs.cpu.Y, 2)
				bs.addShake(0.05, 2.0)
			}
		}
	}

	// 3. CPU Beam vs Player
	if cBeam != nil && cBeam.State == entities.BeamStateFiring {
		beamY := cBeam.Y
		if math.Abs(bs.player.Y-beamY) < 32.0 {
			hitX := cBeam.X
			if (cBeam.DirX > 0 && bs.player.X > hitX) || (cBeam.DirX < 0 && bs.player.X < hitX) {
				bs.player.TakeDamage(120.0*dt, true, cBeam.X)
				bs.particles.SpawnHitSparks(bs.player.X, bs.player.Y, 2)
				bs.addShake(0.05, 2.0)
			}
		}
	}
}

func (bs *BattleScene) spawnFighterVFX(f *entities.Fighter) {
	// Aura during Ki Charge or Sparking
	if f.State == entities.StateChargeKi || f.IsSparking {
		auraCol := color.RGBA{R: 255, G: 220, B: 50, A: 255} // Gold
		if f.Type == gfx.FighterCPU {
			auraCol = color.RGBA{R: 120, G: 180, B: 255, A: 255} // Cyan
		}
		bs.particles.SpawnKiAura(f.X, f.Y+12, auraCol, f.IsSparking)
	}

	// Dash trail
	if f.State == entities.StateDragonDash {
		trailCol := color.RGBA{R: 255, G: 240, B: 120, A: 200}
		if f.Type == gfx.FighterCPU {
			trailCol = color.RGBA{R: 100, G: 220, B: 255, A: 200}
		}
		bs.particles.SpawnDashTrail(f.X, f.Y, f.VX, f.VY, trailCol)
	}
}

func (bs *BattleScene) addShake(duration, amp float64) {
	bs.screenShake = duration
	if amp > bs.shakeAmp {
		bs.shakeAmp = amp
	}
}

func (bs *BattleScene) Draw(screen *ebiten.Image) {
	w, h := 640.0, 360.0

	// Dynamic camera screen shake
	var shakeOffX, shakeOffY float64
	if bs.screenShake > 0 {
		shakeOffX = (math.Sin(bs.screenShake*60.0)) * bs.shakeAmp
		shakeOffY = (math.Cos(bs.screenShake*50.0)) * bs.shakeAmp
	}

	// Draw Scenic Arena
	bs.arena.Draw(screen, w, h)

	// Draw Particles (Under fighters)
	bs.particles.Draw(screen)

	// Draw Active Super Beams
	if bs.player.ActiveBeam != nil {
		bs.player.ActiveBeam.Draw(screen)
	}
	if bs.cpu.ActiveBeam != nil {
		bs.cpu.ActiveBeam.Draw(screen)
	}

	// Draw Ki Blasts
	for _, b := range bs.blasts {
		b.Draw(screen)
	}

	// Draw Fighters
	p1Sprite := bs.spriteCache.GetSprite(bs.player.Type, bs.player.GetPose())
	p2Sprite := bs.spriteCache.GetSprite(bs.cpu.Type, bs.cpu.GetPose())

	gfx.DrawFighter(screen, p1Sprite, bs.player.X+shakeOffX, bs.player.Y+shakeOffY, 1.4, 1.4, 0, bs.player.FacingLeft)
	gfx.DrawFighter(screen, p2Sprite, bs.cpu.X+shakeOffX, bs.cpu.Y+shakeOffY, 1.4, 1.4, 0, bs.cpu.FacingLeft)

	// Draw Red Panic Alarm Vignette if player is panicked
	if bs.player.Panic.IsPanicked {
		progress := bs.player.Panic.PanicTimer / entities.MaxPanicDuration
		gfx.DrawPanicVignette(screen, w, h, progress)
	}

	// Draw HUD & UI Prompts
	bs.hud.Draw(screen, bs.player, bs.cpu, w, h)
}
