package entity

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
)

type PickupType int

const (
	PickupTypeRecover PickupType = iota
	PickupTypeMutex
	PickupTypeWorker
)

type Pickup struct {
	Type   PickupType
	X, Y   float64
	VY     float64
	Active bool
	BobAge float64
}

type PickupManager struct {
	pickups []Pickup
}

func NewPickupManager(capacity int) *PickupManager {
	return &PickupManager{
		pickups: make([]Pickup, capacity),
	}
}

func (pm *PickupManager) Spawn(x, y float64, pType PickupType) {
	for i := range pm.pickups {
		if !pm.pickups[i].Active {
			p := &pm.pickups[i]
			p.Active = true
			p.Type = pType
			p.X = x
			p.Y = y
			p.VY = 45.0
			p.BobAge = 0
			break
		}
	}
}

func (pm *PickupManager) Update(dt float64, playerX, playerY float64, screenH float64) {
	for i := range pm.pickups {
		p := &pm.pickups[i]
		if p.Active {
			p.BobAge += dt
			p.Y += p.VY * dt

			// Slight magnetic attraction if player is within 90 pixels
			dx := playerX - p.X
			dy := playerY - p.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 90.0 && dist > 1.0 {
				pull := 90.0 * (1.0 - dist/90.0)
				p.X += (dx / dist) * pull * dt
				p.Y += (dy / dist) * pull * dt
			}

			if p.Y > screenH+30 {
				p.Active = false
			}
		}
	}
}

func (pm *PickupManager) Draw(screen *ebiten.Image) {
	for i := range pm.pickups {
		p := &pm.pickups[i]
		if !p.Active {
			continue
		}

		var tex *ebiten.Image
		switch p.Type {
		case PickupTypeRecover:
			tex = art.PickupRecover
		case PickupTypeMutex:
			tex = art.PickupMutex
		case PickupTypeWorker:
			tex = art.PickupWorker
		}

		if tex == nil {
			continue
		}

		// Floating bob
		bobOffset := math.Sin(p.BobAge*6.0) * 2.0

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(-9, -9) // Center pivot for 18x18
		op.GeoM.Translate(p.X, p.Y+bobOffset)
		screen.DrawImage(tex, &op)
	}
}

func (pm *PickupManager) Pickups() []Pickup {
	return pm.pickups
}

func (pm *PickupManager) Deactivate(index int) {
	if index >= 0 && index < len(pm.pickups) {
		pm.pickups[index].Active = false
	}
}
