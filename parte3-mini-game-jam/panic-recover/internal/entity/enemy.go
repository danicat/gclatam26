package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/art"
	"panic-recover/internal/audio"
)

type EnemyType int

const (
	EnemyTypeNilPointer EnemyType = iota
	EnemyTypeConcurrent
	EnemyTypeDeadlock
	EnemyTypeMemoryLeak
	EnemyTypeGoroutine
)

type Enemy struct {
	Type        EnemyType
	X, Y        float64
	VX, VY      float64
	HP          float64
	MaxHP       float64
	Points      int
	Radius      float64
	Age         float64
	ShootTimer  float64
	Active      bool
	IsMiniSplit bool // For memory leak split children
}

type EnemyManager struct {
	enemies         []Enemy
	spawnTimer      float64
	waveTimer       float64
	tier            int
	speedMultiplier float64
	hpMultiplier    float64
}

func NewEnemyManager(capacity int) *EnemyManager {
	return &EnemyManager{
		enemies:         make([]Enemy, capacity),
		tier:            1,
		speedMultiplier: 1.0,
		hpMultiplier:    1.0,
	}
}

func (em *EnemyManager) SetTier(tier int) {
	em.tier = tier
	em.speedMultiplier = 1.0 + float64(tier-1)*0.18
	em.hpMultiplier = 1.0 + float64(tier-1)*0.25
}

func (em *EnemyManager) Spawn(eType EnemyType, x, y float64, isMini bool) {
	for i := range em.enemies {
		if !em.enemies[i].Active {
			e := &em.enemies[i]
			e.Active = true
			e.Type = eType
			e.X = x
			e.Y = y
			e.Age = 0
			e.IsMiniSplit = isMini

			switch eType {
			case EnemyTypeNilPointer:
				e.MaxHP = 30
				e.HP = e.MaxHP
				e.Points = 150
				e.Radius = 12
				e.VY = 80.0
				e.VX = 60.0
				e.ShootTimer = 0.8 + rand.Float64()*0.6

			case EnemyTypeConcurrent:
				e.MaxHP = 50
				e.HP = e.MaxHP
				e.Points = 250
				e.Radius = 14
				e.VY = 55.0
				e.VX = 0
				e.ShootTimer = 1.2 + rand.Float64()*0.8

			case EnemyTypeDeadlock:
				e.MaxHP = 120
				e.HP = e.MaxHP
				e.Points = 500
				e.Radius = 16
				e.VY = 30.0
				e.VX = 0
				e.ShootTimer = 1.5 + rand.Float64()*0.5

			case EnemyTypeMemoryLeak:
				if isMini {
					e.MaxHP = 20
					e.HP = e.MaxHP
					e.Points = 100
					e.Radius = 8
					e.VY = 65.0
					e.VX = (rand.Float64() - 0.5) * 80.0
					e.ShootTimer = 999.0 // Mini leaks don't shoot
				} else {
					e.MaxHP = 60
					e.HP = e.MaxHP
					e.Points = 300
					e.Radius = 14
					e.VY = 45.0
					e.VX = 0
					e.ShootTimer = 1.8 + rand.Float64()*0.8
				}

			case EnemyTypeGoroutine:
				e.MaxHP = 18
				e.HP = e.MaxHP
				e.Points = 80
				e.Radius = 9
				e.VY = 120.0
				e.VX = 0
				e.ShootTimer = 999.0 // Dive-bombers
			}

			// Apply difficulty scaling to all spawned enemy types
			e.MaxHP *= em.hpMultiplier
			e.HP = e.MaxHP
			e.VY *= em.speedMultiplier
			e.VX *= em.speedMultiplier
			e.Points = int(float64(e.Points) * (1.0 + float64(em.tier-1)*0.25))
			break
		}
	}
}

func (em *EnemyManager) Update(dt float64, screenW, screenH float64, bm *BulletManager, pm *PickupManager, ps *ParticleSystem, playerX, playerY float64, inPanic bool, bossActive bool) int {
	scoreGained := 0

	// Wave Spawner
	em.spawnTimer += dt
	em.waveTimer += dt

	spawnInterval := (1.2 - math.Min(0.7, em.waveTimer*0.005)) / em.speedMultiplier
	if bossActive {
		spawnInterval *= 2.8 // Throttle small enemy waves during boss fight
	}
	if em.spawnTimer >= spawnInterval {
		em.spawnTimer = 0
		r := rand.Float64()
		spawnX := 40.0 + rand.Float64()*(screenW-80.0)

		if r < 0.35 {
			em.Spawn(EnemyTypeNilPointer, spawnX, -20, false)
		} else if r < 0.60 {
			em.Spawn(EnemyTypeGoroutine, spawnX, -20, false)
			// Spawn small formation
			if rand.Float64() < 0.5 {
				em.Spawn(EnemyTypeGoroutine, spawnX-25, -35, false)
				em.Spawn(EnemyTypeGoroutine, spawnX+25, -35, false)
			}
		} else if r < 0.80 {
			em.Spawn(EnemyTypeConcurrent, spawnX, -20, false)
		} else if r < 0.92 {
			em.Spawn(EnemyTypeMemoryLeak, spawnX, -20, false)
		} else {
			em.Spawn(EnemyTypeDeadlock, spawnX, -30, false)
		}
	}

	// Update existing enemies
	for i := range em.enemies {
		e := &em.enemies[i]
		if !e.Active {
			continue
		}

		e.Age += dt

		// Movement behaviors
		switch e.Type {
		case EnemyTypeNilPointer:
			// Zigzag pattern
			e.X += math.Sin(e.Age*5.0) * e.VX * dt
			e.Y += e.VY * dt
		case EnemyTypeConcurrent:
			// Wave sweep
			e.X += math.Cos(e.Age*3.0) * 80.0 * dt
			e.Y += e.VY * dt
		case EnemyTypeDeadlock:
			e.Y += e.VY * dt
		case EnemyTypeMemoryLeak:
			e.X += e.VX * dt
			e.Y += e.VY * dt
		case EnemyTypeGoroutine:
			e.Y += e.VY * dt
		}

		// Shooting logic
		e.ShootTimer -= dt
		if e.ShootTimer <= 0 {
			switch e.Type {
			case EnemyTypeNilPointer:
				e.ShootTimer = 1.4 + rand.Float64()*0.8
				bm.Spawn(e.X, e.Y+10, 0, 160.0, 15, false, false)
				audio.PlayEnemyLaser()
			case EnemyTypeConcurrent:
				e.ShootTimer = 1.6 + rand.Float64()*0.8
				bm.Spawn(e.X-8, e.Y+10, -40, 150.0, 15, false, false)
				bm.Spawn(e.X+8, e.Y+10, 40, 150.0, 15, false, false)
				audio.PlayEnemyLaser()
			case EnemyTypeDeadlock:
				e.ShootTimer = 2.0 + rand.Float64()*0.5
				// 3-way spread
				bm.Spawn(e.X, e.Y+12, 0, 140.0, 20, false, false)
				bm.Spawn(e.X, e.Y+12, -60, 130.0, 20, false, false)
				bm.Spawn(e.X, e.Y+12, 60, 130.0, 20, false, false)
				audio.PlayEnemyLaser()
			case EnemyTypeMemoryLeak:
				e.ShootTimer = 2.2
				// Aimed shot toward player
				dx := playerX - e.X
				dy := playerY - e.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					speed := 130.0
					bm.Spawn(e.X, e.Y+10, (dx/dist)*speed, (dy/dist)*speed, 15, false, false)
					audio.PlayEnemyLaser()
				}
			}
		}

		// Offscreen cleanup
		if e.Y > screenH+40 || e.X < -50 || e.X > screenW+50 {
			e.Active = false
		}
	}

	return scoreGained
}

func (em *EnemyManager) HandleBulletCollisions(bm *BulletManager, pm *PickupManager, ps *ParticleSystem, inPanic bool) int {
	scoreGained := 0
	bullets := bm.Bullets()

	for bIdx := range bullets {
		b := &bullets[bIdx]
		if !b.Active || !b.IsPlayer {
			continue
		}

		for eIdx := range em.enemies {
			e := &em.enemies[eIdx]
			if !e.Active {
				continue
			}

			// Circle/Box overlap test
			dx := b.X - e.X
			dy := b.Y - e.Y
			distSq := dx*dx + dy*dy
			hitRadius := e.Radius + b.Width*0.5

			if distSq <= hitRadius*hitRadius {
				// Bullet hit enemy!
				e.HP -= b.Damage
				b.Active = false

				// Hit sparks
				ps.EmitExplosion(b.X, b.Y, 4, color.RGBA{255, 255, 200, 255})

				if e.HP <= 0 {
					// Enemy destroyed!
					e.Active = false
					scoreGained += e.Points
					audio.PlayExplosion()

					// Death explosion
					var expColor color.RGBA
					switch e.Type {
					case EnemyTypeNilPointer:
						expColor = color.RGBA{200, 50, 255, 255}
					case EnemyTypeConcurrent:
						expColor = color.RGBA{255, 50, 50, 255}
					case EnemyTypeDeadlock:
						expColor = color.RGBA{255, 200, 50, 255}
					case EnemyTypeMemoryLeak:
						expColor = color.RGBA{50, 255, 80, 255}
					case EnemyTypeGoroutine:
						expColor = color.RGBA{255, 130, 20, 255}
					}
					ps.EmitExplosion(e.X, e.Y, 20, expColor)

					// Memory leak special: splits into 2 mini-leaks
					if e.Type == EnemyTypeMemoryLeak && !e.IsMiniSplit {
						em.Spawn(EnemyTypeMemoryLeak, e.X-10, e.Y, true)
						em.Spawn(EnemyTypeMemoryLeak, e.X+10, e.Y, true)
					}

					// DROPS
					if inPanic {
						// CRITICAL HOOK: IN PANIC MODE, GUARANTEE RECOVER DROP!
						pm.Spawn(e.X, e.Y, PickupTypeRecover)
					} else {
						// Normal drops
						roll := rand.Float64()
						if roll < 0.12 {
							pm.Spawn(e.X, e.Y, PickupTypeRecover)
						} else if roll < 0.22 {
							pm.Spawn(e.X, e.Y, PickupTypeMutex)
						} else if roll < 0.30 {
							pm.Spawn(e.X, e.Y, PickupTypeWorker)
						}
					}
				}
				break // Bullet consumed
			}
		}
	}
	return scoreGained
}

func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for i := range em.enemies {
		e := &em.enemies[i]
		if !e.Active {
			continue
		}

		var tex *ebiten.Image
		switch e.Type {
		case EnemyTypeNilPointer:
			tex = art.EnemyNilPointer
		case EnemyTypeConcurrent:
			tex = art.EnemyConcurrent
		case EnemyTypeDeadlock:
			tex = art.EnemyDeadlock
		case EnemyTypeMemoryLeak:
			tex = art.EnemyMemoryLeak
		case EnemyTypeGoroutine:
			tex = art.EnemyGoroutine
		}

		if tex == nil {
			continue
		}

		bounds := tex.Bounds()
		tw := float64(bounds.Dx())
		th := float64(bounds.Dy())

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(-tw/2.0, -th/2.0)

		if e.IsMiniSplit {
			op.GeoM.Scale(0.65, 0.65)
		}

		op.GeoM.Translate(e.X, e.Y)
		screen.DrawImage(tex, &op)
	}
}

func (em *EnemyManager) Enemies() []Enemy {
	return em.enemies
}
