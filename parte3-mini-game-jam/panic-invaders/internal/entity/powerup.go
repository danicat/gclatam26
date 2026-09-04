package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"panic-invaders/internal/assets"
)

type PowerupType int

const (
	PowerupMutex PowerupType = iota
	PowerupChan
	PowerupTimeout
	PowerupBadge
)

type Powerup struct {
	X      float64
	Y      float64
	Vy     float64
	Type   PowerupType
	Active bool
	Width  float64
	Height float64
}

func NewPowerup(x, y float64, pType PowerupType) *Powerup {
	return &Powerup{
		X:      x,
		Y:      y,
		Vy:     1.8,
		Type:   pType,
		Active: true,
		Width:  12,
		Height: 12,
	}
}

func (p *Powerup) Update() {
	if !p.Active {
		return
	}
	p.Y += p.Vy
	if p.Y > 370 {
		p.Active = false
	}
}

func (p *Powerup) Draw(screen *ebiten.Image) {
	if !p.Active {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)

	switch p.Type {
	case PowerupMutex:
		screen.DrawImage(assets.LoadedSprites.PowerupMutex, op)
	case PowerupChan:
		screen.DrawImage(assets.LoadedSprites.PowerupChan, op)
	case PowerupTimeout:
		screen.DrawImage(assets.LoadedSprites.PowerupTimeout, op)
	case PowerupBadge:
		screen.DrawImage(assets.LoadedSprites.PowerupBadge, op)
	}
}
