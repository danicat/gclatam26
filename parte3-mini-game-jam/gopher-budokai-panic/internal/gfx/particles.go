package gfx

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ParticleType int

const (
	ParticleAura ParticleType = iota
	ParticleSpark
	ParticleSmoke
	ParticleShockwave
	ParticleDashLine
)

type Particle struct {
	Type       ParticleType
	X, Y       float64
	VX, VY     float64
	Life, Age  float64
	StartSize  float32
	EndSize    float32
	Color      color.RGBA
	Additive   bool
	Active     bool
}

type ParticlePool struct {
	pool []Particle
	rng  *rand.Rand
}

func NewParticlePool(capacity int) *ParticlePool {
	return &ParticlePool{
		pool: make([]Particle, capacity),
		rng:  rand.New(rand.NewSource(42)),
	}
}

func (pp *ParticlePool) spawn(p Particle) {
	for i := range pp.pool {
		if !pp.pool[i].Active {
			pp.pool[i] = p
			pp.pool[i].Active = true
			return
		}
	}
}

func (pp *ParticlePool) Update(dt float64) {
	for i := range pp.pool {
		p := &pp.pool[i]
		if !p.Active {
			continue
		}
		p.Age += dt
		if p.Age >= p.Life {
			p.Active = false
			continue
		}
		p.X += p.VX * dt
		p.Y += p.VY * dt
	}
}

func (pp *ParticlePool) Draw(screen *ebiten.Image) {
	for i := range pp.pool {
		p := &pp.pool[i]
		if !p.Active {
			continue
		}
		progress := float32(p.Age / p.Life)
		size := p.StartSize + (p.EndSize-p.StartSize)*progress
		if size <= 0 {
			continue
		}

		alpha := float32(1.0 - progress)
		c := color.RGBA{
			R: uint8(float32(p.Color.R) * alpha),
			G: uint8(float32(p.Color.G) * alpha),
			B: uint8(float32(p.Color.B) * alpha),
			A: uint8(float32(p.Color.A) * alpha),
		}

		switch p.Type {
		case ParticleShockwave:
			// Expanding circular wave outline
			vector.StrokeCircle(screen, float32(p.X), float32(p.Y), size, 2.5, c, true)
		case ParticleDashLine:
			// Elongated speed streak
			vector.StrokeLine(screen, float32(p.X), float32(p.Y), float32(p.X-p.VX*0.06), float32(p.Y-p.VY*0.06), size, c, true)
		default:
			// Solid energy particle
			vector.DrawFilledCircle(screen, float32(p.X), float32(p.Y), size, c, true)
		}
	}
}

// Spawners

func (pp *ParticlePool) SpawnKiAura(x, y float64, col color.RGBA, isSparking bool) {
	count := 2
	if isSparking {
		count = 5
	}
	for i := 0; i < count; i++ {
		offsetX := (pp.rng.Float64() - 0.5) * 24.0
		offsetY := pp.rng.Float64() * 20.0
		vy := -(40.0 + pp.rng.Float64()*60.0)
		if isSparking {
			vy *= 1.6
		}
		pp.spawn(Particle{
			Type:      ParticleAura,
			X:         x + offsetX,
			Y:         y + offsetY,
			VX:        (pp.rng.Float64() - 0.5) * 25.0,
			VY:        vy,
			Life:      0.25 + pp.rng.Float64()*0.25,
			StartSize: 3.5,
			EndSize:   0.5,
			Color:     col,
			Additive:  true,
		})
	}
}

func (pp *ParticlePool) SpawnDashTrail(x, y, vx, vy float64, col color.RGBA) {
	pp.spawn(Particle{
		Type:      ParticleDashLine,
		X:         x + (pp.rng.Float64()-0.5)*10.0,
		Y:         y + (pp.rng.Float64()-0.5)*10.0,
		VX:        vx,
		VY:        vy,
		Life:      0.15,
		StartSize: 2.0,
		EndSize:   0.5,
		Color:     col,
		Additive:  true,
	})
}

func (pp *ParticlePool) SpawnHitSparks(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := pp.rng.Float64() * 2.0 * math.Pi
		speed := 80.0 + pp.rng.Float64()*180.0
		col := color.RGBA{R: 255, G: 230, B: 80, A: 255}
		if pp.rng.Float64() < 0.3 {
			col = color.RGBA{R: 255, G: 120, B: 40, A: 255}
		}
		pp.spawn(Particle{
			Type:      ParticleSpark,
			X:         x,
			Y:         y,
			VX:        math.Cos(angle) * speed,
			VY:        math.Sin(angle) * speed,
			Life:      0.18 + pp.rng.Float64()*0.15,
			StartSize: 3.0,
			EndSize:   0.5,
			Color:     col,
			Additive:  true,
		})
	}
}

func (pp *ParticlePool) SpawnRecoverWave(x, y float64) {
	// Massive circular shockwaves expanding outward
	for i := 0; i < 3; i++ {
		pp.spawn(Particle{
			Type:      ParticleShockwave,
			X:         x,
			Y:         y,
			Life:      0.35 + float64(i)*0.08,
			StartSize: 5.0,
			EndSize:   90.0 + float32(i)*30.0,
			Color:     color.RGBA{R: 240, G: 250, B: 255, A: 255},
			Additive:  true,
		})
	}
	// Radial blast sparks
	for i := 0; i < 24; i++ {
		angle := float64(i) * (2.0 * math.Pi / 24.0)
		speed := 160.0 + pp.rng.Float64()*80.0
		pp.spawn(Particle{
			Type:      ParticleSpark,
			X:         x,
			Y:         y,
			VX:        math.Cos(angle) * speed,
			VY:        math.Sin(angle) * speed,
			Life:      0.30,
			StartSize: 4.0,
			EndSize:   0.5,
			Color:     color.RGBA{R: 120, G: 220, B: 255, A: 255},
			Additive:  true,
		})
	}
}

func (pp *ParticlePool) SpawnExplosion(x, y float64, col color.RGBA) {
	pp.spawn(Particle{
		Type:      ParticleShockwave,
		X:         x,
		Y:         y,
		Life:      0.25,
		StartSize: 2.0,
		EndSize:   45.0,
		Color:     col,
		Additive:  true,
	})
	for i := 0; i < 16; i++ {
		angle := pp.rng.Float64() * 2.0 * math.Pi
		speed := 60.0 + pp.rng.Float64()*120.0
		pp.spawn(Particle{
			Type:      ParticleSpark,
			X:         x,
			Y:         y,
			VX:        math.Cos(angle) * speed,
			VY:        math.Sin(angle) * speed,
			Life:      0.22,
			StartSize: 3.5,
			EndSize:   0.8,
			Color:     col,
			Additive:  true,
		})
	}
}
