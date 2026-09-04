package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
)

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life       float64
	Age        float64
	StartSize  float64
	EndSize    float64
	StartColor color.RGBA
	EndColor   color.RGBA
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

func (ps *ParticleSystem) EmitExplosion(x, y float64, count int, mainColor color.RGBA) {
	spawned := 0
	for i := range ps.pool {
		if !ps.pool[i].Active {
			p := &ps.pool[i]
			p.Active = true
			p.X = x
			p.Y = y
			angle := rand.Float64() * 2.0 * math.Pi
			speed := 40.0 + rand.Float64()*180.0
			p.VX = math.Cos(angle) * speed
			p.VY = math.Sin(angle) * speed
			p.Life = 0.35 + rand.Float64()*0.35
			p.Age = 0
			p.StartSize = 1.4 + rand.Float64()*0.8
			p.EndSize = 0.1
			p.StartColor = mainColor
			p.EndColor = color.RGBA{mainColor.R / 4, mainColor.G / 4, mainColor.B / 4, 0}

			spawned++
			if spawned >= count {
				break
			}
		}
	}
}

func (ps *ParticleSystem) EmitThruster(x, y float64, isPanic bool) {
	for i := range ps.pool {
		if !ps.pool[i].Active {
			p := &ps.pool[i]
			p.Active = true
			p.X = x + (rand.Float64()-0.5)*4.0
			p.Y = y
			p.VX = (rand.Float64() - 0.5) * 20.0
			p.VY = 60.0 + rand.Float64()*60.0
			p.Life = 0.15 + rand.Float64()*0.1
			p.Age = 0
			p.StartSize = 0.8
			p.EndSize = 0.1

			if isPanic {
				p.StartColor = color.RGBA{255, 60, 20, 240}
				p.EndColor = color.RGBA{255, 200, 40, 0}
			} else {
				p.StartColor = color.RGBA{80, 200, 255, 240}
				p.EndColor = color.RGBA{30, 80, 255, 0}
			}
			break
		}
	}
}

func (ps *ParticleSystem) EmitShockwave(x, y float64, count int) {
	spawned := 0
	for i := range ps.pool {
		if !ps.pool[i].Active {
			p := &ps.pool[i]
			p.Active = true
			p.X = x
			p.Y = y
			angle := (float64(spawned) / float64(count)) * 2.0 * math.Pi
			speed := 220.0
			p.VX = math.Cos(angle) * speed
			p.VY = math.Sin(angle) * speed
			p.Life = 0.45
			p.Age = 0
			p.StartSize = 1.8
			p.EndSize = 0.2
			p.StartColor = color.RGBA{80, 255, 160, 255}
			p.EndColor = color.RGBA{20, 150, 80, 0}

			spawned++
			if spawned >= count {
				break
			}
		}
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.pool {
		p := &ps.pool[i]
		if p.Active {
			p.Age += dt
			if p.Age >= p.Life {
				p.Active = false
				continue
			}
			p.X += p.VX * dt
			p.Y += p.VY * dt
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	tex := art.ParticleGlow
	if tex == nil {
		return
	}

	for i := range ps.pool {
		p := &ps.pool[i]
		if !p.Active {
			continue
		}
		prog := p.Age / p.Life
		size := p.StartSize + prog*(p.EndSize-p.StartSize)

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(-8, -8) // center pivot
		op.GeoM.Scale(size, size)
		op.GeoM.Translate(p.X, p.Y)

		curColor := lerpRGBA(p.StartColor, p.EndColor, prog)
		op.ColorScale.ScaleWithColor(curColor)
		op.Blend = ebiten.BlendLighter

		screen.DrawImage(tex, &op)
	}
}

func lerpRGBA(c1, c2 color.RGBA, t float64) color.RGBA {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}
	return color.RGBA{
		R: uint8(float64(c1.R) + t*(float64(c2.R)-float64(c1.R))),
		G: uint8(float64(c1.G) + t*(float64(c2.G)-float64(c1.G))),
		B: uint8(float64(c1.B) + t*(float64(c2.B)-float64(c1.B))),
		A: uint8(float64(c1.A) + t*(float64(c2.A)-float64(c1.A))),
	}
}
