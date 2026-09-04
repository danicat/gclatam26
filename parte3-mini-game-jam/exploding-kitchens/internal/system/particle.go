package system

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// ParticleKind identifies the visual style and physics of a particle.
type ParticleKind int

const (
	ParticleSteam ParticleKind = iota
	ParticleSmoke
	ParticleFire
	ParticleSpark
	ParticleExplosion
	ParticleClutch
)

// Particle represents a single active or recycled visual effect.
type Particle struct {
	X, Y        float64
	VX, VY      float64
	Age, Life   float64
	Size        float64
	EndSize     float64
	Color       color.RGBA
	EndColor    color.RGBA
	Kind        ParticleKind
	Additive    bool
	Active      bool
}

// ParticlePool manages pre-allocated particles with zero heap allocation during gameplay.
type ParticlePool struct {
	particles []Particle
	dotImg    *ebiten.Image
	drawOpts  ebiten.DrawImageOptions
}

// NewParticlePool allocates a fixed particle pool buffer.
func NewParticlePool(capacity int) *ParticlePool {
	// Create a 2x2 white square texture for rendering scaled particles
	dot := ebiten.NewImage(2, 2)
	dot.Fill(color.White)

	return &ParticlePool{
		particles: make([]Particle, capacity),
		dotImg:    dot,
	}
}

// Spawn emits a single particle into the pool if a slot is available.
func (pp *ParticlePool) Spawn(p Particle) {
	for i := range pp.particles {
		if !pp.particles[i].Active {
			pp.particles[i] = p
			pp.particles[i].Active = true
			pp.particles[i].Age = 0
			return
		}
	}
}

// EmitSteam creates a burst of rising white/translucent steam puffs.
func (pp *ParticlePool) EmitSteam(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*0.8
		speed := 15.0 + rand.Float64()*25.0
		pp.Spawn(Particle{
			X:        x + (rand.Float64()-0.5)*6,
			Y:        y + (rand.Float64()-0.5)*4,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     0.4 + rand.Float64()*0.4,
			Size:     1.5,
			EndSize:  4.5,
			Color:    color.RGBA{220, 230, 245, 180},
			EndColor: color.RGBA{180, 200, 220, 0},
			Kind:     ParticleSteam,
		})
	}
}

// EmitSmoke creates thick billowing smoke clouds.
func (pp *ParticlePool) EmitSmoke(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*1.2
		speed := 20.0 + rand.Float64()*30.0
		pp.Spawn(Particle{
			X:        x + (rand.Float64()-0.5)*8,
			Y:        y + (rand.Float64()-0.5)*6,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     0.5 + rand.Float64()*0.5,
			Size:     2.0,
			EndSize:  6.0,
			Color:    color.RGBA{50, 45, 60, 220},
			EndColor: color.RGBA{20, 15, 30, 0},
			Kind:     ParticleSmoke,
		})
	}
}

// EmitFire creates flickering hot flames and sparks.
func (pp *ParticlePool) EmitFire(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*1.0
		speed := 25.0 + rand.Float64()*40.0
		pp.Spawn(Particle{
			X:        x + (rand.Float64()-0.5)*8,
			Y:        y + (rand.Float64()-0.5)*4,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     0.3 + rand.Float64()*0.3,
			Size:     2.0,
			EndSize:  0.8,
			Color:    color.RGBA{255, 210, 40, 255},
			EndColor: color.RGBA{220, 40, 20, 0},
			Kind:     ParticleFire,
			Additive: true,
		})
	}
}

// EmitExplosion produces a violent spherical blast of debris and fireball shockwave.
func (pp *ParticlePool) EmitExplosion(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 30.0 + rand.Float64()*90.0
		pp.Spawn(Particle{
			X:        x,
			Y:        y,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     0.4 + rand.Float64()*0.4,
			Size:     3.0,
			EndSize:  8.0,
			Color:    color.RGBA{255, 240, 180, 255},
			EndColor: color.RGBA{200, 30, 10, 0},
			Kind:     ParticleExplosion,
			Additive: true,
		})
	}
}

// EmitClutch produces golden recovery stars when a last-second defusal is pulled off.
func (pp *ParticlePool) EmitClutch(x, y float64, count int) {
	for i := 0; i < count; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*2.0
		speed := 30.0 + rand.Float64()*50.0
		pp.Spawn(Particle{
			X:        x + (rand.Float64()-0.5)*10,
			Y:        y + (rand.Float64()-0.5)*6,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     0.6 + rand.Float64()*0.4,
			Size:     2.0,
			EndSize:  4.0,
			Color:    color.RGBA{255, 245, 80, 255},
			EndColor: color.RGBA{80, 230, 255, 0},
			Kind:     ParticleClutch,
			Additive: true,
		})
	}
}

// Update advances active particle lifespans and positions by dt.
func (pp *ParticlePool) Update(dt float64) {
	for i := range pp.particles {
		p := &pp.particles[i]
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

		// Drag/deceleration on explosion and sparks
		if p.Kind == ParticleExplosion || p.Kind == ParticleSpark {
			p.VX *= 0.94
			p.VY *= 0.94
		}
	}
}

// Draw renders all active particles with zero heap allocations.
func (pp *ParticlePool) Draw(screen *ebiten.Image) {
	for i := range pp.particles {
		p := &pp.particles[i]
		if !p.Active {
			continue
		}

		progress := p.Age / p.Life
		if progress > 1.0 {
			progress = 1.0
		}

		curSize := p.Size + (p.EndSize-p.Size)*progress
		scale := curSize / 2.0 // based on 2x2 base texture

		r := float64(p.Color.R) + float64(int(p.EndColor.R)-int(p.Color.R))*progress
		g := float64(p.Color.G) + float64(int(p.EndColor.G)-int(p.Color.G))*progress
		b := float64(p.Color.B) + float64(int(p.EndColor.B)-int(p.Color.B))*progress
		a := float64(p.Color.A) + float64(int(p.EndColor.A)-int(p.Color.A))*progress

		pp.drawOpts.GeoM.Reset()
		pp.drawOpts.GeoM.Scale(scale, scale)
		pp.drawOpts.GeoM.Translate(p.X-curSize/2, p.Y-curSize/2)

		// Premultiplied alpha color scale
		af := a / 255.0
		pp.drawOpts.ColorScale.Reset()
		pp.drawOpts.ColorScale.Scale(
			float32((r/255.0)*af),
			float32((g/255.0)*af),
			float32((b/255.0)*af),
			float32(af),
		)

		if p.Additive {
			pp.drawOpts.Blend = ebiten.BlendLighter
		} else {
			pp.drawOpts.Blend = ebiten.BlendSourceOver
		}

		screen.DrawImage(pp.dotImg, &pp.drawOpts)
	}
}
