package entities

import (
	"image/color"
	"math"
	"math/rand"
)

type ParticleType int

const (
	PartDust ParticleType = iota
	PartStar
	PartWater
	PartSmoke
	PartConfetti
)

// Particle representa um elemento visual efêmero na tela.
type Particle struct {
	Active   bool
	Type     ParticleType
	X, Y, Z  float64
	Vx, Vy, Vz float64
	Color    color.RGBA
	Life     float64
	MaxLife  float64
	Size     float64
}

// ParticleSystem gerencia um pool fixo de partículas pré-alocadas sem GC overhead.
type ParticleSystem struct {
	particles []Particle
}

// NewParticleSystem cria o pool com 250 partículas pré-alocadas.
func NewParticleSystem(capacity int) *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, capacity),
	}
}

// Spawn adiciona uma partícula no pool reutilizando slots inativos.
func (ps *ParticleSystem) Spawn(pType ParticleType, x, y, z float64, vx, vy, vz float64, col color.RGBA, life float64, size float64) {
	for i := range ps.particles {
		if !ps.particles[i].Active {
			ps.particles[i] = Particle{
				Active:  true,
				Type:    pType,
				X:       x,
				Y:       y,
				Z:       z,
				Vx:      vx,
				Vy:      vy,
				Vz:      vz,
				Color:   col,
				Life:    life,
				MaxLife: life,
				Size:    size,
			}
			return
		}
	}
}

// BurstStars emite uma explosão estelar comemorativa dourada (para entregas perfeitas).
func (ps *ParticleSystem) BurstStars(x, y, z float64, count int) {
	goldColors := []color.RGBA{
		{255, 230, 40, 255},
		{255, 215, 0, 255},
		{30, 110, 240, 255}, // Azul Express
		{255, 255, 255, 255},
	}

	for i := 0; i < count; i++ {
		angle := rand.Float64() * 6.283
		spd := 30.0 + rand.Float64()*70.0
		col := goldColors[rand.Intn(len(goldColors))]
		ps.Spawn(
			PartStar,
			x, y, z,
			math.Cos(angle)*spd,
			math.Sin(angle)*spd,
			20.0+rand.Float64()*40.0,
			col,
			0.6+rand.Float64()*0.4,
			4.0,
		)
	}
}

// SpawnDust emite poeira nas rodas da bicicleta/veículo.
func (ps *ParticleSystem) SpawnDust(x, y, z float64) {
	ps.Spawn(
		PartDust,
		x-4.0+rand.Float64()*8.0,
		y-4.0+rand.Float64()*8.0,
		z,
		-10.0+rand.Float64()*20.0,
		-15.0-rand.Float64()*15.0,
		5.0+rand.Float64()*10.0,
		color.RGBA{180, 175, 170, 180},
		0.4+rand.Float64()*0.2,
		3.0,
	)
}

// Update avança a física das partículas ativas.
func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.particles {
		p := &ps.particles[i]
		if !p.Active {
			continue
		}

		p.Life -= dt
		if p.Life <= 0 {
			p.Active = false
			continue
		}

		p.X += p.Vx * dt
		p.Y += p.Vy * dt
		p.Z += p.Vz * dt
		p.Vz -= 90.0 * dt // Gravidade suave

		if p.Z < 0 {
			p.Z = 0
			p.Vz = -p.Vz * 0.4 // Quique leve
		}
	}
}

// GetActiveParticles retorna as partículas ativas para renderização.
func (ps *ParticleSystem) GetActiveParticles() []Particle {
	return ps.particles
}
