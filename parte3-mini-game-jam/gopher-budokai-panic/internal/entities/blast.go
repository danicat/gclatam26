package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type KiBlast struct {
	X, Y       float64
	VX, VY     float64
	Radius     float32
	OwnerID    int
	Damage     float64
	Color      color.RGBA
	Active     bool
	TailLength float32
}

func NewKiBlast(x, y, dirX, dirY float64, ownerID int, isSparking bool) *KiBlast {
	speed := 480.0
	col := color.RGBA{R: 255, G: 235, B: 60, A: 255} // Gold
	dmg := 22.0
	rad := float32(6.0)

	if ownerID != 0 {
		// CPU / Vegeta style cyan/purple blast
		col = color.RGBA{R: 80, G: 210, B: 255, A: 255}
	}
	if isSparking {
		dmg *= 1.4
		rad *= 1.3
		speed *= 1.2
	}

	return &KiBlast{
		X:          x,
		Y:          y,
		VX:         dirX * speed,
		VY:         dirY * speed,
		Radius:     rad,
		OwnerID:    ownerID,
		Damage:     dmg,
		Color:      col,
		Active:     true,
		TailLength: 14.0,
	}
}

func (kb *KiBlast) Update(dt float64, screenW, screenH float64) {
	if !kb.Active {
		return
	}
	kb.X += kb.VX * dt
	kb.Y += kb.VY * dt

	if kb.X < -50 || kb.X > screenW+50 || kb.Y < -50 || kb.Y > screenH+50 {
		kb.Active = false
	}
}

func (kb *KiBlast) Draw(screen *ebiten.Image) {
	if !kb.Active {
		return
	}
	x := float32(kb.X)
	y := float32(kb.Y)

	// Glowing outer halo
	haloCol := color.RGBA{
		R: kb.Color.R,
		G: kb.Color.G,
		B: kb.Color.B,
		A: 140,
	}
	vector.DrawFilledCircle(screen, x, y, kb.Radius*1.8, haloCol, true)

	// Bright core
	vector.DrawFilledCircle(screen, x, y, kb.Radius, kb.Color, true)
	// Pure white hot center
	white := color.RGBA{R: 255, G: 255, B: 255, A: 240}
	vector.DrawFilledCircle(screen, x, y, kb.Radius*0.45, white, true)

	// Trailing speed streak
	streakCol := color.RGBA{R: kb.Color.R, G: kb.Color.G, B: kb.Color.B, A: 120}
	vector.StrokeLine(screen, x, y, x-float32(kb.VX*0.025), y-float32(kb.VY*0.025), kb.Radius*0.9, streakCol, true)
}
