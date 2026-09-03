package gfx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ParticleKind int

const (
	ParticleMirrorShard ParticleKind = iota
	ParticleSpark
	ParticleSweat
	ParticleDashGhost
	ParticleSmoke
)

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life, Age  float64
	StartSize  float64
	EndSize    float64
	Color      color.RGBA
	Kind       ParticleKind
	Active     bool
}

type ParticleSystem struct {
	pool []Particle
}

func NewParticleSystem(capacity int) *ParticleSystem {
	return &ParticleSystem{
		pool: make([]Particle, capacity),
	}
}

func (ps *ParticleSystem) Emit(x, y, vx, vy, life, startSize, endSize float64, col color.RGBA, kind ParticleKind) {
	for i := range ps.pool {
		if !ps.pool[i].Active {
			ps.pool[i] = Particle{
				X:         x,
				Y:         y,
				VX:        vx,
				VY:        vy,
				Life:      life,
				Age:       0,
				StartSize: startSize,
				EndSize:   endSize,
				Color:     col,
				Kind:      kind,
				Active:    true,
			}
			return
		}
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.pool {
		p := &ps.pool[i]
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

		// Apply gravity or deceleration depending on particle kind
		switch p.Kind {
		case ParticleMirrorShard:
			p.VY += 300.0 * dt // Falling shard gravity
		case ParticleSweat:
			p.VY += 150.0 * dt
		case ParticleSmoke:
			p.VY -= 20.0 * dt // Rising smoke
			p.VX *= 0.95
		case ParticleDashGhost:
			p.VX *= 0.90
			p.VY *= 0.90
		case ParticleSpark:
			p.VX *= 0.92
			p.VY *= 0.92
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for i := range ps.pool {
		p := &ps.pool[i]
		if !p.Active {
			continue
		}

		progress := p.Age / p.Life
		if progress > 1.0 {
			progress = 1.0
		}
		size := p.StartSize + (p.EndSize-p.StartSize)*progress
		if size <= 0.5 {
			size = 0.5
		}

		// Fade alpha over lifetime
		alphaFactor := 1.0 - progress
		c := p.Color
		c.A = uint8(float64(c.A) * alphaFactor)
		c.R = uint8(float64(c.R) * alphaFactor)
		c.G = uint8(float64(c.G) * alphaFactor)
		c.B = uint8(float64(c.B) * alphaFactor)

		switch p.Kind {
		case ParticleDashGhost:
			// Rectangular ghost slice
			vector.DrawFilledRect(screen, float32(p.X-size/2), float32(p.Y-size), float32(size), float32(size*2), c, false)
		case ParticleMirrorShard:
			// Diamond / angled shard
			vector.DrawFilledRect(screen, float32(p.X-size/2), float32(p.Y-size/2), float32(size), float32(size), c, false)
		default:
			vector.DrawFilledCircle(screen, float32(p.X), float32(p.Y), float32(size), c, false)
		}
	}
}

func (ps *ParticleSystem) Reset() {
	for i := range ps.pool {
		ps.pool[i].Active = false
	}
}
