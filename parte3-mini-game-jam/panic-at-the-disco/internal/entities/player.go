package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-at-the-disco/internal/audio"
	"panic-at-the-disco/internal/gfx"
	"panic-at-the-disco/internal/input"
)

type Player struct {
	X, Y              float64
	VX, VY            float64
	FacingX           float64
	Width, Height     float64
	HitboxRadius      float64
	Lives             int
	PanicLevel        float64 // 0 to 100%
	GrooveMeter       float64 // 0 to 100%
	IsDashing         bool
	DashDuration      float64
	DashCooldown      float64
	IsSlipping        bool
	SlipTimer         float64
	InvulnerableTimer float64
	AnimTime          float64
	IsMoving          bool
	Score             int
	SweatTimer        float64
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:            x,
		Y:            y,
		FacingX:      1.0,
		Width:        22.0,
		Height:       32.0,
		HitboxRadius: 9.0,
		Lives:        3,
		PanicLevel:   20.0,
		GrooveMeter:  40.0,
		DashCooldown: 0.0,
	}
}

func (p *Player) Bounds() (float64, float64, float64, float64) {
	return p.X - p.Width/2, p.Y - p.Height/2, p.Width, p.Height
}

func (p *Player) Hitbox() (float64, float64, float64) {
	return p.X, p.Y, p.HitboxRadius
}

func (p *Player) Update(dt float64, in input.InputState, boundsX, boundsY, boundsW, boundsH float64, ps *gfx.ParticleSystem, ae *audio.AudioEngine) {
	p.AnimTime += dt

	// Update timers
	if p.InvulnerableTimer > 0 {
		p.InvulnerableTimer -= dt
	}
	if p.DashCooldown > 0 {
		p.DashCooldown -= dt
	}

	// 1. Dash handling
	if p.IsDashing {
		p.DashDuration -= dt
		if p.DashDuration <= 0 {
			p.IsDashing = false
		}
		// Emit neon afterimage ghost particles during dash
		if ps != nil {
			ps.Emit(p.X, p.Y, 0, 0, 0.25, 18.0, 10.0, color.RGBA{0, 240, 255, 200}, gfx.ParticleDashGhost)
		}
	} else if in.DashJustDown && p.DashCooldown <= 0 && p.GrooveMeter >= 30.0 {
		p.IsDashing = true
		p.DashDuration = 0.22
		p.DashCooldown = 0.65
		p.InvulnerableTimer = 0.25
		p.GrooveMeter = math.Max(0.0, p.GrooveMeter-30.0)
		p.PanicLevel = math.Max(0.0, p.PanicLevel-15.0) // Dashing relieves panic!

		// Dash burst speed
		dashSpeed := 380.0
		dx, dy := in.MoveX, in.MoveY
		if dx == 0 && dy == 0 {
			dx = p.FacingX
		}
		p.VX = dx * dashSpeed
		p.VY = dy * dashSpeed

		if ae != nil {
			ae.PlaySFXDash()
		}
	}

	// 2. Normal movement vs Slipping movement
	if !p.IsDashing {
		baseSpeed := 155.0
		if p.IsSlipping {
			p.SlipTimer -= dt
			if p.SlipTimer <= 0 {
				p.IsSlipping = false
			}
			// Slipping: low friction ice physics
			p.VX += in.MoveX * baseSpeed * 2.0 * dt
			p.VY += in.MoveY * baseSpeed * 2.0 * dt
			p.VX *= math.Pow(0.85, dt*60.0)
			p.VY *= math.Pow(0.85, dt*60.0)
		} else {
			// Crisp arcade responsive movement
			targetVX := in.MoveX * baseSpeed
			targetVY := in.MoveY * baseSpeed
			friction := 18.0
			p.VX += (targetVX - p.VX) * friction * dt
			p.VY += (targetVY - p.VY) * friction * dt
		}
	}

	// Determine if moving and update facing direction
	velMag := math.Sqrt(p.VX*p.VX + p.VY*p.VY)
	p.IsMoving = velMag > 15.0
	if in.MoveX != 0 {
		p.FacingX = in.MoveX
	}

	// 3. Position integration
	p.X += p.VX * dt
	p.Y += p.VY * dt

	// 4. Boundary clamping
	minX := boundsX + p.Width/2
	maxX := boundsX + boundsW - p.Width/2
	minY := boundsY + p.Height/2
	maxY := boundsY + boundsH - p.Height/2

	if p.X < minX {
		p.X = minX
		p.VX = 0
	} else if p.X > maxX {
		p.X = maxX
		p.VX = 0
	}
	if p.Y < minY {
		p.Y = minY
		p.VY = 0
	} else if p.Y > maxY {
		p.Y = maxY
		p.VY = 0
	}

	// 5. Panic and Groove dynamics
	// Panic slowly increases over time as ceiling rumbles
	p.PanicLevel = math.Min(100.0, p.PanicLevel+dt*1.2)
	// Groove slowly recharges
	p.GrooveMeter = math.Min(100.0, p.GrooveMeter+dt*7.0)

	// High panic causes sweat drops
	if p.PanicLevel > 50.0 {
		p.SweatTimer += dt
		sweatInterval := 0.4 - (p.PanicLevel-50.0)/200.0
		if p.SweatTimer >= sweatInterval {
			p.SweatTimer = 0
			if ps != nil {
				ps.Emit(p.X+float64(randSign())*7.0, p.Y-14.0, float64(randSign())*25.0, -40.0, 0.4, 3.0, 1.0, color.RGBA{120, 210, 255, 230}, gfx.ParticleSweat)
			}
		}
	}
}

func (p *Player) ApplySlip(duration float64, ae *audio.AudioEngine) {
	if !p.IsSlipping && !p.IsDashing {
		p.IsSlipping = true
		p.SlipTimer = duration
		p.PanicLevel = math.Min(100.0, p.PanicLevel+15.0)
		if ae != nil {
			ae.PlaySFXSlip()
		}
	}
}

func (p *Player) TakeDamage(dmg int, ps *gfx.ParticleSystem, ae *audio.AudioEngine) bool {
	if p.InvulnerableTimer > 0 || p.IsDashing {
		return false
	}
	p.Lives -= dmg
	p.InvulnerableTimer = 1.2
	p.PanicLevel = math.Min(100.0, p.PanicLevel+30.0)

	// Burst sparks around player
	if ps != nil {
		for i := 0; i < 16; i++ {
			ang := float64(i) * (2.0 * math.Pi / 16.0)
			speed := 80.0 + float64(i%5)*20.0
			ps.Emit(p.X, p.Y, math.Cos(ang)*speed, math.Sin(ang)*speed, 0.35, 4.0, 1.0, color.RGBA{255, 60, 60, 255}, gfx.ParticleSpark)
		}
	}
	if ae != nil {
		ae.PlaySFXCrash()
	}
	return true
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Flashing effect during invulnerability
	if p.InvulnerableTimer > 0 {
		flash := math.Sin(p.AnimTime * 35.0)
		if flash < 0 {
			return
		}
	}
	gfx.DrawPlayer(screen, p.X, p.Y, p.FacingX, p.AnimTime, p.IsMoving, p.IsDashing, p.PanicLevel)
}

func randSign() int {
	if rand.Float64() < 0.5 {
		return 1
	}
	return -1
}
