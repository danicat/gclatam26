package entity

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
	"panic-recover/internal/audio"
)

type Boss struct {
	Active      bool
	X, Y        float64
	TargetY     float64
	VX          float64
	HP          float64
	MaxHP       float64
	Name        string
	Radius      float64
	AttackTimer float64
	MoveTimer   float64
	HitFlash    float64
	Phase       int
	Tier        int
}

func NewBoss() *Boss {
	return &Boss{
		Radius: 26.0,
	}
}

func (b *Boss) Spawn(tier int, screenW float64) {
	b.Active = true
	b.Tier = tier
	b.X = screenW / 2.0
	b.Y = -60.0
	b.TargetY = 75.0
	b.VX = 70.0 + float64(tier)*15.0
	b.MaxHP = 400.0 + float64(tier-1)*250.0
	b.HP = b.MaxHP
	b.AttackTimer = 2.0
	b.MoveTimer = 0
	b.HitFlash = 0
	b.Phase = 1

	bossNames := []string{
		"SIGSEGV: SEGMENTATION FAULT",
		"OOM: OUT OF MEMORY KILLER",
		"DEADLOCK: CYCLIC MUTEX TITAN",
		"STACK OVERFLOW: RECURSION HYDRA",
	}
	idx := (tier - 1) % len(bossNames)
	b.Name = fmt.Sprintf("BOSS %d: %s", tier, bossNames[idx])
}

func (b *Boss) Update(dt float64, screenW float64, bm *BulletManager, em *EnemyManager, ps *ParticleSystem, playerX, playerY float64) {
	if !b.Active {
		return
	}

	if b.HitFlash > 0 {
		b.HitFlash -= dt
	}

	// Entrance descent
	if b.Y < b.TargetY {
		b.Y += 60.0 * dt
		return
	}

	// Horizontal patrol
	b.MoveTimer += dt
	b.X += b.VX * dt
	if b.X < 50.0 {
		b.X = 50.0
		b.VX = math.Abs(b.VX)
	} else if b.X > screenW-50.0 {
		b.X = screenW - 50.0
		b.VX = -math.Abs(b.VX)
	}

	// Attack patterns based on HP phase
	b.AttackTimer -= dt
	hpRatio := b.HP / b.MaxHP

	if hpRatio < 0.4 {
		b.Phase = 2 // Enraged
	}

	if b.AttackTimer <= 0 {
		if b.Phase == 1 {
			// Phase 1: Alternating Aimed Double Laser & 6-Way Ring
			b.AttackTimer = 1.3 - math.Min(0.5, float64(b.Tier)*0.1)
			if rand.Float64() < 0.6 {
				// Aimed shot toward player
				dx := playerX - b.X
				dy := playerY - b.Y
				dist := math.Hypot(dx, dy)
				if dist > 0 {
					speed := 160.0
					bm.Spawn(b.X-18, b.Y+20, (dx/dist)*speed-20, (dy/dist)*speed, 18, false, false)
					bm.Spawn(b.X+18, b.Y+20, (dx/dist)*speed+20, (dy/dist)*speed, 18, false, false)
					audio.PlayEnemyLaser()
				}
			} else {
				// 6-way radial star burst
				numBullets := 6
				for i := 0; i < numBullets; i++ {
					angle := (float64(i) / float64(numBullets)) * 2.0 * math.Pi
					speed := 140.0
					bm.Spawn(b.X, b.Y+15, math.Cos(angle)*speed, math.Sin(angle)*speed, 15, false, false)
				}
				audio.PlayEnemyLaser()
			}
		} else {
			// Phase 2: Enraged Barrage + Minion Call
			b.AttackTimer = 1.0 - math.Min(0.4, float64(b.Tier)*0.08)

			// 8-way rotating spiral burst
			numBullets := 8
			spiralOffset := b.MoveTimer * 2.0
			for i := 0; i < numBullets; i++ {
				angle := (float64(i)/float64(numBullets))*2.0*math.Pi + spiralOffset
				speed := 150.0
				bm.Spawn(b.X, b.Y+15, math.Cos(angle)*speed, math.Sin(angle)*speed, 16, false, false)
			}
			audio.PlayEnemyLaser()

			// 25% chance to call minion bugs
			if rand.Float64() < 0.35 {
				em.Spawn(EnemyTypeNilPointer, b.X-30, b.Y, false)
				em.Spawn(EnemyTypeGoroutine, b.X+30, b.Y, false)
			}
		}
	}
}

func (b *Boss) TakeDamage(dmg float64, ps *ParticleSystem) bool {
	if !b.Active {
		return false
	}
	b.HP -= dmg
	b.HitFlash = 0.1
	ps.EmitExplosion(b.X+(rand.Float64()-0.5)*30, b.Y+(rand.Float64()-0.5)*20, 3, color.RGBA{255, 200, 200, 255})

	if b.HP <= 0 {
		b.Active = false
		return true // Defeated
	}
	return false
}

func (b *Boss) Draw(screen *ebiten.Image) {
	if !b.Active {
		return
	}

	tex := art.BossSigsegv
	if tex == nil {
		return
	}

	bounds := tex.Bounds()
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(-float64(bounds.Dx())/2.0, -float64(bounds.Dy())/2.0)
	op.GeoM.Translate(b.X, b.Y)

	// Flash white/red when taking damage
	if b.HitFlash > 0 {
		op.ColorScale.Scale(1.8, 1.8, 1.8, 1.0)
	}

	screen.DrawImage(tex, &op)
}
