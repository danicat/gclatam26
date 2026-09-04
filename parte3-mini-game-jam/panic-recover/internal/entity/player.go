package entity

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
	"panic-recover/internal/audio"
)

type Player struct {
	X, Y          float64
	Speed         float64
	HP            float64
	MaxHP         float64
	InPanic       bool
	PanicTimer    float64 // 5.0 seconds frenzy countdown
	ShieldTimer   float64 // sync.Mutex invulnerability
	DroneCount    int     // go worker companions (0, 1, or 2)
	DroneAngle    float64
	RecoverStock  int     // Stored recover EMP bombs
	ShootCooldown float64
	Score         int
	HitRadius     float64
	InvulnFlash   float64
	IsDead        bool
}

func NewPlayer(startX, startY float64) *Player {
	return &Player{
		X:             startX,
		Y:             startY,
		Speed:         240.0,
		HP:            100.0,
		MaxHP:         100.0,
		HitRadius:     7.0, // Forgiving arcade hitbox
		DroneCount:    0,
		RecoverStock:  0,
		ShootCooldown: 0,
	}
}

func (p *Player) Update(dt float64, screenW, screenH float64, ps *ParticleSystem) {
	if p.IsDead {
		return
	}

	p.DroneAngle += 3.5 * dt
	if p.ShootCooldown > 0 {
		p.ShootCooldown -= dt
	}
	if p.InvulnFlash > 0 {
		p.InvulnFlash -= dt
	}

	// Update Shield
	if p.ShieldTimer > 0 {
		p.ShieldTimer -= dt
	}

	// Update Panic state
	if p.InPanic {
		p.PanicTimer -= dt
		// Thruster particles in panic mode
		ps.EmitThruster(p.X-6, p.Y+12, true)
		ps.EmitThruster(p.X+6, p.Y+12, true)

		if p.PanicTimer <= 0 {
			// Call stack completely unwound! Crash!
			p.IsDead = true
			audio.PlayExplosion()
			ps.EmitExplosion(p.X, p.Y, 40, color.RGBA{255, 30, 20, 255})
		}
	} else {
		// Normal engine thruster particles
		ps.EmitThruster(p.X-6, p.Y+12, false)
		ps.EmitThruster(p.X+6, p.Y+12, false)
	}

	// Clamp boundaries
	if p.X < 20 {
		p.X = 20
	}
	if p.X > screenW-20 {
		p.X = screenW - 20
	}
	if p.Y < 30 {
		p.Y = 30
	}
	if p.Y > screenH-25 {
		p.Y = screenH - 25
	}
}

func (p *Player) Move(dx, dy float64, dt float64) {
	if p.IsDead {
		return
	}
	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		inv := 1.0 / math.Sqrt(2.0)
		dx *= inv
		dy *= inv
	}
	p.X += dx * p.Speed * dt
	p.Y += dy * p.Speed * dt
}

func (p *Player) TryShoot(bm *BulletManager) {
	if p.IsDead || p.ShootCooldown > 0 {
		return
	}

	if p.InPanic {
		// OVERDRIVE: 3x fire rate erratic burst
		p.ShootCooldown = 0.055
		audio.PlayLaser()

		// 3-way heavy flaming bullets
		bm.Spawn(p.X-8, p.Y-10, -25, -420, 28, true, true)
		bm.Spawn(p.X, p.Y-14, 0, -450, 32, true, true)
		bm.Spawn(p.X+8, p.Y-10, 25, -420, 28, true, true)
	} else {
		// Normal twin cyan lasers
		p.ShootCooldown = 0.13
		audio.PlayLaser()
		bm.Spawn(p.X-7, p.Y-10, 0, -420, 18, true, false)
		bm.Spawn(p.X+7, p.Y-10, 0, -420, 18, true, false)
	}

	// Drones fire auxiliary lasers
	if p.DroneCount > 0 {
		d1x, d1y := p.GetDronePos(0)
		bm.Spawn(d1x, d1y-6, 0, -400, 12, true, false)
	}
	if p.DroneCount > 1 {
		d2x, d2y := p.GetDronePos(1)
		bm.Spawn(d2x, d2y-6, 0, -400, 12, true, false)
	}
}

func (p *Player) TakeDamage(dmg float64, ps *ParticleSystem) bool {
	if p.IsDead || p.ShieldTimer > 0 {
		return false
	}

	if p.InPanic {
		// Invulnerable during the 5s panic surge!
		return false
	}

	p.HP -= dmg
	p.InvulnFlash = 0.2
	ps.EmitExplosion(p.X, p.Y, 8, color.RGBA{255, 120, 50, 255})

	if p.HP <= 0 {
		// TRIGGER NEAR-DEATH SURGE: ENTER PANIC MODE!
		p.HP = 1
		p.InPanic = true
		p.PanicTimer = 5.0
		audio.PlayPanicSiren()
		ps.EmitShockwave(p.X, p.Y, 24)
		return true
	}
	return false
}

func (p *Player) RecoverFromPanic(ps *ParticleSystem, bm *BulletManager) {
	p.InPanic = false
	p.HP = 60.0
	p.PanicTimer = 0
	audio.PlayRecoverChime()
	ps.EmitShockwave(p.X, p.Y, 32)
	bm.ClearEnemyBullets(ps)
}

func (p *Player) UseRecoverBomb(ps *ParticleSystem, bm *BulletManager) bool {
	if p.RecoverStock <= 0 {
		return false
	}
	p.RecoverStock--
	audio.PlayRecoverChime()
	ps.EmitShockwave(p.X, p.Y, 36)
	bm.ClearEnemyBullets(ps)
	return true
}

func (p *Player) CollectPickup(pickup *Pickup, ps *ParticleSystem, bm *BulletManager) {
	audio.PlayPickup()

	switch pickup.Type {
	case PickupTypeRecover:
		if p.InPanic {
			p.RecoverFromPanic(ps, bm)
		} else {
			p.HP = math.Min(p.MaxHP, p.HP+40)
			p.RecoverStock++
			p.Score += 250
		}
	case PickupTypeMutex:
		p.ShieldTimer = 6.0
		p.Score += 150
	case PickupTypeWorker:
		if p.DroneCount < 2 {
			p.DroneCount++
		}
		p.Score += 200
	}
}

func (p *Player) GetDronePos(index int) (float64, float64) {
	angle := p.DroneAngle
	if index == 1 {
		angle += math.Pi
	}
	orbitRadius := 26.0
	dx := math.Cos(angle) * orbitRadius
	dy := math.Sin(angle) * orbitRadius * 0.5 // elliptical orbit
	return p.X + dx, p.Y + dy
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.IsDead {
		return
	}

	// Flicker effect when in panic or recently hit
	if p.InvulnFlash > 0 {
		if int(p.InvulnFlash*40)%2 == 0 {
			return
		}
	}

	// Draw Gopher Vanguard Ship
	shipTex := art.PlayerShip
	if shipTex != nil {
		var op ebiten.DrawImageOptions
		bounds := shipTex.Bounds()
		op.GeoM.Translate(-float64(bounds.Dx())/2.0, -float64(bounds.Dy())/2.0)
		op.GeoM.Translate(p.X, p.Y)

		if p.InPanic {
			// Red-orange pulsing tint in panic mode
			op.ColorScale.Scale(1.4, 0.5, 0.4, 1.0)
		}

		screen.DrawImage(shipTex, &op)
	}

	// Draw Companion Drones
	droneTex := art.DroneShip
	if droneTex != nil {
		for i := 0; i < p.DroneCount; i++ {
			dx, dy := p.GetDronePos(i)
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(-8, -8) // 16x16 center
			op.GeoM.Translate(dx, dy)
			screen.DrawImage(droneTex, &op)
		}
	}

	// Draw sync.Mutex Shield Aura if active
	if p.ShieldTimer > 0 {
		shieldTex := art.ShieldAura
		if shieldTex != nil {
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(-24, -24)
			op.GeoM.Translate(p.X, p.Y)
			op.Blend = ebiten.BlendLighter
			screen.DrawImage(shieldTex, &op)
		}
	}
}
