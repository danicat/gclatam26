package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"panic-recover/internal/game"
)

type Particle struct {
	Position game.Vec2
	Velocity game.Vec2
	Life     float64
	Age      float64
	Color    color.RGBA
	Active   bool
}

type ParticleSystem struct {
	particles []Particle
}

func NewParticleSystem(capacity int) *ParticleSystem {
	if capacity < 0 {
		capacity = 0
	}
	return &ParticleSystem{particles: make([]Particle, capacity)}
}

func (ps *ParticleSystem) Spawn(position, velocity game.Vec2, life float64, particleColor color.RGBA) {
	for i := range ps.particles {
		if ps.particles[i].Active {
			continue
		}
		ps.particles[i] = Particle{
			Position: position,
			Velocity: velocity,
			Life:     life,
			Color:    particleColor,
			Active:   life > 0,
		}
		return
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.particles {
		particle := &ps.particles[i]
		if !particle.Active {
			continue
		}
		particle.Age += dt
		if particle.Age >= particle.Life {
			particle.Active = false
			continue
		}
		particle.Position.X += particle.Velocity.X * dt
		particle.Position.Y += particle.Velocity.Y * dt
	}
}

func (ps *ParticleSystem) ActiveCount() int {
	count := 0
	for i := range ps.particles {
		if ps.particles[i].Active {
			count++
		}
	}
	return count
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for i := range ps.particles {
		particle := &ps.particles[i]
		if !particle.Active {
			continue
		}
		alpha := float64(particle.Color.A) * (1 - particle.Age/particle.Life)
		if alpha < 0 {
			alpha = 0
		}
		particleColor := particle.Color
		particleColor.A = uint8(alpha)
		vector.DrawFilledCircle(screen, float32(particle.Position.X), float32(particle.Position.Y), 1.5, particleColor, false)
	}
}
