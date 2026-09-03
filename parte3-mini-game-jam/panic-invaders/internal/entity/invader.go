package entity

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-invaders/internal/assets"
	"panic-invaders/internal/audio"
)

type InvaderType int

const (
	InvaderNil InvaderType = iota
	InvaderIndex
	InvaderDivide
)

type Invader struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Type   InvaderType
	Alive  bool
}

type UFO struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Active bool
	Speed  float64
}

type Boss struct {
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Active   bool
	MaxHP    int
	HP       int
	Speed    float64
	ShootCD  int
}

type InvaderFleet struct {
	Invaders    []*Invader
	Direction   float64
	SpeedX      float64
	StepY       float64
	AnimFrame   int
	AnimTimer   int
	ShootTimer  int
	UFOSpawnCD  int
	UFO         *UFO
	Boss        *Boss
	Wave        int
}

func NewInvaderFleet(wave int) *InvaderFleet {
	fleet := &InvaderFleet{
		Direction:  1.0,
		SpeedX:     1.0 + float64(wave)*0.3,
		StepY:      12.0,
		AnimFrame:  0,
		AnimTimer:  0,
		ShootTimer: 60,
		UFOSpawnCD: 400 + rand.Intn(300),
		UFO:        &UFO{Width: 20, Height: 10, Active: false},
		Wave:       wave,
	}

	rows := 3
	cols := 7
	if wave > 1 {
		rows = 4
		cols = 8
	}

	startX := 100.0
	startY := 50.0
	spacingX := 45.0
	spacingY := 26.0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			iType := InvaderDivide
			if r == 0 {
				iType = InvaderNil
			} else if r == 1 {
				iType = InvaderIndex
			}

			fleet.Invaders = append(fleet.Invaders, &Invader{
				X:      startX + float64(c)*spacingX,
				Y:      startY + float64(r)*spacingY,
				Width:  16,
				Height: 12,
				Type:   iType,
				Alive:  true,
			})
		}
	}

	if wave >= 3 {
		fleet.Boss = &Boss{
			X:       280,
			Y:       35,
			Width:   48,
			Height:  24,
			Active:  true,
			MaxHP:   25,
			HP:      25,
			Speed:   1.5,
			ShootCD: 45,
		}
	}

	return fleet
}

func (f *InvaderFleet) AliveCount() int {
	count := 0
	for _, inv := range f.Invaders {
		if inv.Alive {
			count++
		}
	}
	return count
}

func (f *InvaderFleet) Update(bullets *[]*Bullet, powerups *[]*Powerup, timeoutActive bool) bool {
	// Animation cycle
	f.AnimTimer++
	if f.AnimTimer > 30 {
		f.AnimTimer = 0
		f.AnimFrame = 1 - f.AnimFrame
	}

	effectiveSpeed := f.SpeedX
	if timeoutActive {
		effectiveSpeed *= 0.4
	}

	// Speed up as invaders are eliminated
	alive := f.AliveCount()
	if alive > 0 && alive < 8 {
		effectiveSpeed *= 1.5
	} else if alive == 1 {
		effectiveSpeed *= 2.0
	}

	// Move invaders
	moveDown := false
	for _, inv := range f.Invaders {
		if !inv.Alive {
			continue
		}
		newX := inv.X + f.Direction*effectiveSpeed
		if newX < 20 || newX > 640-inv.Width-20 {
			moveDown = true
			break
		}
	}

	if moveDown {
		f.Direction *= -1
		for _, inv := range f.Invaders {
			if inv.Alive {
				inv.Y += f.StepY
				// Check if reached bottom
				if inv.Y >= 290 {
					return true // Stack bottom breach!
				}
			}
		}
	} else {
		for _, inv := range f.Invaders {
			if inv.Alive {
				inv.X += f.Direction * effectiveSpeed
			}
		}
	}

	// Shooting
	f.ShootTimer--
	if f.ShootTimer <= 0 {
		f.ShootTimer = 45 - f.Wave*5
		if f.ShootTimer < 20 {
			f.ShootTimer = 20
		}
		f.shootPanicBullet(bullets)
	}

	// UFO management
	f.updateUFO(bullets, powerups)

	// Boss management
	if f.Boss != nil && f.Boss.Active {
		f.updateBoss(bullets, powerups)
	}

	return false
}

func (f *InvaderFleet) shootPanicBullet(bullets *[]*Bullet) {
	// Collect bottom-most alive invaders in each column
	type colInvader struct {
		inv *Invader
	}
	cols := make(map[int]*colInvader)

	for _, inv := range f.Invaders {
		if !inv.Alive {
			continue
		}
		colKey := int(inv.X / 30.0)
		if existing, ok := cols[colKey]; !ok || inv.Y > existing.inv.Y {
			cols[colKey] = &colInvader{inv: inv}
		}
	}

	if len(cols) == 0 {
		return
	}

	// Pick random shooter
	list := make([]*Invader, 0, len(cols))
	for _, v := range cols {
		list = append(list, v.inv)
	}
	shooter := list[rand.Intn(len(list))]
	*bullets = append(*bullets, NewPanicBullet(shooter.X+shooter.Width/2-1, shooter.Y+shooter.Height))
}

func (f *InvaderFleet) updateUFO(bullets *[]*Bullet, powerups *[]*Powerup) {
	if f.UFO.Active {
		f.UFO.X += f.UFO.Speed
		if f.UFO.X > 660 || f.UFO.X < -30 {
			f.UFO.Active = false
		}
	} else {
		f.UFOSpawnCD--
		if f.UFOSpawnCD <= 0 {
			f.UFO.Active = true
			f.UFOSpawnCD = 500 + rand.Intn(400)
			if rand.Float64() < 0.5 {
				f.UFO.X = -20
				f.UFO.Speed = 2.0
			} else {
				f.UFO.X = 650
				f.UFO.Speed = -2.0
			}
			f.UFO.Y = 28
		}
	}
}

func (f *InvaderFleet) updateBoss(bullets *[]*Bullet, powerups *[]*Powerup) {
	b := f.Boss
	b.X += b.Speed
	if b.X < 50 || b.X > 640-b.Width-50 {
		b.Speed *= -1
	}

	b.ShootCD--
	if b.ShootCD <= 0 {
		b.ShootCD = 50
		// Boss fires a 3-bullet spread
		*bullets = append(*bullets,
			NewPanicBullet(b.X+b.Width/2-1, b.Y+b.Height),
			NewPanicBullet(b.X+10, b.Y+b.Height),
			NewPanicBullet(b.X+b.Width-10, b.Y+b.Height),
		)
	}
}

func (f *InvaderFleet) CheckBulletCollisions(bullet *Bullet, player *Player, powerups *[]*Powerup) bool {
	if !bullet.Active || bullet.IsEnemy {
		return false
	}

	bx := bullet.X
	by := bullet.Y
	bw := bullet.Width
	bh := bullet.Height

	// Check Invaders
	for _, inv := range f.Invaders {
		if !inv.Alive {
			continue
		}
		if bx+bw >= inv.X && bx <= inv.X+inv.Width &&
			by+bh >= inv.Y && by <= inv.Y+inv.Height {
			inv.Alive = false
			bullet.Active = false
			audio.GlobalAudio.PlayHit()

			// Score
			points := 10
			switch inv.Type {
			case InvaderNil:
				points = 30
			case InvaderIndex:
				points = 20
			}
			player.Score += points

			// Chance of power-up drop
			if rand.Float64() < 0.12 {
				pTypes := []PowerupType{PowerupMutex, PowerupChan, PowerupTimeout}
				pType := pTypes[rand.Intn(len(pTypes))]
				*powerups = append(*powerups, NewPowerup(inv.X+inv.Width/2, inv.Y, pType))
			}
			return true
		}
	}

	// Check UFO
	if f.UFO.Active {
		u := f.UFO
		if bx+bw >= u.X && bx <= u.X+u.Width &&
			by+bh >= u.Y && by <= u.Y+u.Height {
			u.Active = false
			bullet.Active = false
			player.Score += 250
			audio.GlobalAudio.PlayExplosion()
			// UFO always drops GopherCon LATAM Badge!
			*powerups = append(*powerups, NewPowerup(u.X, u.Y, PowerupBadge))
			return true
		}
	}

	// Check Boss
	if f.Boss != nil && f.Boss.Active {
		b := f.Boss
		if bx+bw >= b.X && bx <= b.X+b.Width &&
			by+bh >= b.Y && by <= b.Y+b.Height {
			bullet.Active = false
			b.HP--
			audio.GlobalAudio.PlayHit()

			if b.HP <= 0 {
				b.Active = false
				player.Score += 1000
				audio.GlobalAudio.PlayExplosion()
				*powerups = append(*powerups, NewPowerup(b.X+b.Width/2, b.Y+b.Height/2, PowerupBadge))
			}
			return true
		}
	}

	return false
}

func (f *InvaderFleet) Draw(screen *ebiten.Image) {
	for _, inv := range f.Invaders {
		if !inv.Alive {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(inv.X, inv.Y)

		switch inv.Type {
		case InvaderNil:
			screen.DrawImage(assets.LoadedSprites.InvaderNil[f.AnimFrame], op)
		case InvaderIndex:
			screen.DrawImage(assets.LoadedSprites.InvaderIndex[f.AnimFrame], op)
		case InvaderDivide:
			screen.DrawImage(assets.LoadedSprites.InvaderDivide[f.AnimFrame], op)
		}
	}

	if f.UFO.Active {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(f.UFO.X, f.UFO.Y)
		screen.DrawImage(assets.LoadedSprites.UFO, op)
	}

	if f.Boss != nil && f.Boss.Active {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(f.Boss.X, f.Boss.Y)
		screen.DrawImage(assets.LoadedSprites.Boss, op)

		// Boss HP Bar
		barW := 60.0
		barH := 4.0
		hpPct := float64(f.Boss.HP) / float64(f.Boss.MaxHP)
		hpW := barW * hpPct

		hpBg := ebiten.NewImage(int(barW), int(barH))
		hpBg.Fill(assets.ColorBlack)
		opBg := &ebiten.DrawImageOptions{}
		opBg.GeoM.Translate(f.Boss.X-6, f.Boss.Y-8)
		screen.DrawImage(hpBg, opBg)

		if hpW > 0 {
			hpFg := ebiten.NewImage(int(hpW), int(barH))
			hpFg.Fill(assets.ColorCorruptedCore)
			screen.DrawImage(hpFg, opBg)
		}
	}
}
