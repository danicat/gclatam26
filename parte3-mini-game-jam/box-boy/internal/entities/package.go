package entities

import (
	"box-boy/internal/config"
	"math"
)

// PackageProjectile representa uma caixa de encomenda arremessada no ar.
type PackageProjectile struct {
	StartX, StartY, StartZ float64
	TargetX, TargetY       float64
	CurrentX, CurrentY     float64
	CurrentZ               float64

	Progress   float64 // 0.0 (saída do ciclista) a 1.0 (chegada ao alvo)
	Duration   float64
	PeakHeight float64
	PackageType int // 0: Padrão Amarelo, 1: Frágil, 2: Especial Grande

	TargetHouse *House
	HasLanded   bool
	Success     bool
}

// NewPackageProjectile cria um novo arremesso balístico de encomenda em direção ao alvo.
func NewPackageProjectile(px, py, pz float64, tx, ty float64, pkgType int, house *House) *PackageProjectile {
	dx := tx - px
	dy := ty - py
	dist := math.Hypot(dx, dy)
	dur := dist / config.PackageThrowSpeed
	if dur < 0.25 {
		dur = 0.25
	}

	return &PackageProjectile{
		StartX:      px,
		StartY:      py,
		StartZ:      pz + 12.0, // Sai da altura das mãos do entregador
		TargetX:     tx,
		TargetY:     ty,
		CurrentX:    px,
		CurrentY:    py,
		CurrentZ:    pz + 12.0,
		Progress:    0,
		Duration:    dur,
		PeakHeight:  38.0,
		PackageType: pkgType,
		TargetHouse: house,
	}
}

// Update avança a trajetória balística 3D do pacote.
func (p *PackageProjectile) Update(dt float64) bool {
	if p.HasLanded {
		return true
	}

	p.Progress += dt / p.Duration
	if p.Progress >= 1.0 {
		p.Progress = 1.0
		p.HasLanded = true
		p.CurrentX = p.TargetX
		p.CurrentY = p.TargetY
		p.CurrentZ = 0
		return true
	}

	t := p.Progress
	p.CurrentX = p.StartX + t*(p.TargetX-p.StartX)
	p.CurrentY = p.StartY + t*(p.TargetY-p.StartY)

	// Parábola de altura Z
	arc := 4.0 * p.PeakHeight * t * (1.0 - t)
	p.CurrentZ = (1.0-t)*p.StartZ + arc

	return false
}
