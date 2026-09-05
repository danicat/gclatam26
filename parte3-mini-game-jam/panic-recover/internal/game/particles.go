package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const maxParticles = 128

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life, Age  float64
	Size       float64
	R, G, B, A float32
	Active     bool
}

type ParticleSystem struct {
	pool [maxParticles]Particle
	opts *ebiten.DrawImageOptions
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		opts: &ebiten.DrawImageOptions{},
	}
}

func (ps *ParticleSystem) spawn() *Particle {
	for i := 0; i < maxParticles; i++ {
		if !ps.pool[i].Active {
			p := &ps.pool[i]
			p.Active = true
			p.Age = 0
			return p
		}
	}
	return nil
}

// EmitDust emits dust puffs when a boulder is pushed or falls into a hole
func (ps *ParticleSystem) EmitDust(x, y float64, count int) {
	for i := 0; i < count; i++ {
		p := ps.spawn()
		if p == nil {
			return
		}
		angle := rand.Float64() * 2 * math.Pi
		speed := 15.0 + rand.Float64()*35.0
		p.X = x + float64(TileSize)/2
		p.Y = y + float64(TileSize)/2
		p.VX = math.Cos(angle) * speed
		p.VY = math.Sin(angle) * speed
		p.Life = 0.25 + rand.Float64()*0.2
		p.Size = 2.0 + rand.Float64()*2.0
		p.R = 0.65
		p.G = 0.68
		p.B = 0.75
		p.A = 0.9
	}
}

// EmitSparkles emits magical golden recovery sparkles
func (ps *ParticleSystem) EmitSparkles(x, y float64, count int) {
	for i := 0; i < count; i++ {
		p := ps.spawn()
		if p == nil {
			return
		}
		angle := rand.Float64() * 2 * math.Pi
		speed := 20.0 + rand.Float64()*40.0
		p.X = x + float64(TileSize)/2
		p.Y = y + float64(TileSize)/2
		p.VX = math.Cos(angle) * speed
		p.VY = math.Sin(angle)*speed - 15.0
		p.Life = 0.4 + rand.Float64()*0.3
		p.Size = 2.0 + rand.Float64()*2.5
		p.R = 1.0
		p.G = 0.85
		p.B = 0.2
		p.A = 1.0
	}
}

// EmitVoidMist emits subtle cosmic dark-purple mist from void holes
func (ps *ParticleSystem) EmitVoidMist(x, y float64) {
	p := ps.spawn()
	if p == nil {
		return
	}
	p.X = x + 4 + rand.Float64()*float64(TileSize-8)
	p.Y = y + float64(TileSize) - 4
	p.VX = (rand.Float64() - 0.5) * 8.0
	p.VY = -(8.0 + rand.Float64()*12.0)
	p.Life = 0.6 + rand.Float64()*0.4
	p.Size = 2.0 + rand.Float64()*2.0
	p.R = 0.5
	p.G = 0.15
	p.B = 0.7
	p.A = 0.6
}

// EmitPanicWisps emits red/purple cosmic distortion sparks when in >= 80% panic
func (ps *ParticleSystem) EmitPanicWisps(x, y float64) {
	p := ps.spawn()
	if p == nil {
		return
	}
	p.X = x + float64(TileSize)/2 + (rand.Float64()-0.5)*16
	p.Y = y + float64(TileSize)/2 + (rand.Float64()-0.5)*16
	p.VX = (rand.Float64() - 0.5) * 20.0
	p.VY = -15.0 - rand.Float64()*25.0
	p.Life = 0.3 + rand.Float64()*0.2
	p.Size = 1.5 + rand.Float64()*2.0
	p.R = 0.95
	p.G = 0.1
	p.B = 0.25
	p.A = 0.8
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := 0; i < maxParticles; i++ {
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
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for i := 0; i < maxParticles; i++ {
		p := &ps.pool[i]
		if !p.Active {
			continue
		}
		lifeRatio := float32(1.0 - (p.Age / p.Life))
		alpha := p.A * lifeRatio

		ps.opts.GeoM.Reset()
		ps.opts.GeoM.Scale(p.Size, p.Size)
		ps.opts.GeoM.Translate(p.X, p.Y)
		ps.opts.ColorScale.Reset()
		ps.opts.ColorScale.Scale(p.R, p.G, p.B, alpha)

		screen.DrawImage(imgWhitePixel, ps.opts)
	}
}

// Global 1x1 white pixel for zero-allocation rendering of rectangles & particles
var imgWhitePixel *ebiten.Image

func initWhitePixel() {
	if imgWhitePixel == nil {
		imgWhitePixel = ebiten.NewImage(1, 1)
		imgWhitePixel.Fill(color.White)
	}
}

// Zero-allocation helper to draw filled colored rectangles on screen via GPU
func drawRectGPU(dst *ebiten.Image, opts *ebiten.DrawImageOptions, x, y, w, h float64, r, g, b, a float32) {
	if w <= 0 || h <= 0 {
		return
	}
	opts.GeoM.Reset()
	opts.GeoM.Scale(w, h)
	opts.GeoM.Translate(x, y)
	opts.ColorScale.Reset()
	opts.ColorScale.Scale(r, g, b, a)
	dst.DrawImage(imgWhitePixel, opts)
}
