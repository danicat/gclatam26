package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"panic-invaders/internal/assets"
)

type Bullet struct {
	X       float64
	Y       float64
	Vy      float64
	IsEnemy bool
	Active  bool
	Width   float64
	Height  float64
}

func NewHeroBullet(x, y float64) *Bullet {
	return &Bullet{
		X:       x,
		Y:       y,
		Vy:      -7.0,
		IsEnemy: false,
		Active:  true,
		Width:   3,
		Height:  10,
	}
}

func NewPanicBullet(x, y float64) *Bullet {
	return &Bullet{
		X:       x,
		Y:       y,
		Vy:      3.5,
		IsEnemy: true,
		Active:  true,
		Width:   3,
		Height:  8,
	}
}

func (b *Bullet) Update() {
	if !b.Active {
		return
	}
	b.Y += b.Vy
	if b.Y < -10 || b.Y > 370 {
		b.Active = false
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if !b.Active {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X, b.Y)
	if b.IsEnemy {
		screen.DrawImage(assets.LoadedSprites.BulletPanic, op)
	} else {
		screen.DrawImage(assets.LoadedSprites.BulletHero, op)
	}
}
