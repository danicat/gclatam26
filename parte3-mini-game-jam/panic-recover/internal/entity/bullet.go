package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
)

type Bullet struct {
	X, Y       float64
	VX, VY     float64
	Damage     float64
	IsPlayer   bool
	IsPanic    bool
	Active     bool
	Width      float64
	Height     float64
}

type BulletManager struct {
	bullets []Bullet
}

func NewBulletManager(capacity int) *BulletManager {
	return &BulletManager{
		bullets: make([]Bullet, capacity),
	}
}

func (bm *BulletManager) Spawn(x, y, vx, vy, damage float64, isPlayer, isPanic bool) {
	for i := range bm.bullets {
		if !bm.bullets[i].Active {
			b := &bm.bullets[i]
			b.Active = true
			b.X = x
			b.Y = y
			b.VX = vx
			b.VY = vy
			b.Damage = damage
			b.IsPlayer = isPlayer
			b.IsPanic = isPanic
			if isPlayer {
				if isPanic {
					b.Width = 8
					b.Height = 18
				} else {
					b.Width = 6
					b.Height = 14
				}
			} else {
				b.Width = 8
				b.Height = 8
			}
			break
		}
	}
}

func (bm *BulletManager) Update(dt float64, screenW, screenH float64) {
	for i := range bm.bullets {
		b := &bm.bullets[i]
		if b.Active {
			b.X += b.VX * dt
			b.Y += b.VY * dt
			if b.Y < -30 || b.Y > screenH+30 || b.X < -30 || b.X > screenW+30 {
				b.Active = false
			}
		}
	}
}

func (bm *BulletManager) ClearEnemyBullets(ps *ParticleSystem) {
	for i := range bm.bullets {
		b := &bm.bullets[i]
		if b.Active && !b.IsPlayer {
			b.Active = false
			if ps != nil {
				ps.EmitThruster(b.X, b.Y, false)
			}
		}
	}
}

func (bm *BulletManager) Draw(screen *ebiten.Image) {
	for i := range bm.bullets {
		b := &bm.bullets[i]
		if !b.Active {
			continue
		}

		var tex *ebiten.Image
		if b.IsPlayer {
			if b.IsPanic {
				tex = art.BulletPanic
			} else {
				tex = art.BulletPlayer
			}
		} else {
			tex = art.BulletEnemy
		}

		if tex == nil {
			continue
		}

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(-b.Width/2.0, -b.Height/2.0)
		op.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(tex, &op)
	}
}

func (bm *BulletManager) Bullets() []Bullet {
	return bm.bullets
}

func (bm *BulletManager) Deactivate(index int) {
	if index >= 0 && index < len(bm.bullets) {
		bm.bullets[index].Active = false
	}
}
