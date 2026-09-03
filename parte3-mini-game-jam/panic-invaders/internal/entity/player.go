package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"panic-invaders/internal/assets"
	"panic-invaders/internal/audio"
)

type Player struct {
	X            float64
	Y            float64
	Speed        float64
	Width        float64
	Height       float64
	Lives        int
	Score        int
	ShootCooldown int
	ShieldTimer  int
	ChanTimer    int
	TimeoutTimer int
	InvulnTimer  int
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:            x,
		Y:            y,
		Speed:        3.5,
		Width:        22,
		Height:       16,
		Lives:        3,
		Score:        0,
		ShootCooldown: 0,
		ShieldTimer:  0,
		ChanTimer:    0,
		TimeoutTimer: 0,
		InvulnTimer:  0,
	}
}

func (p *Player) Update(bullets *[]*Bullet) {
	if p.ShootCooldown > 0 {
		p.ShootCooldown--
	}
	if p.ShieldTimer > 0 {
		p.ShieldTimer--
	}
	if p.ChanTimer > 0 {
		p.ChanTimer--
	}
	if p.TimeoutTimer > 0 {
		p.TimeoutTimer--
	}
	if p.InvulnTimer > 0 {
		p.InvulnTimer--
	}

	// Movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= p.Speed
		if p.X < 16 {
			p.X = 16
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += p.Speed
		if p.X > 640-p.Width-16 {
			p.X = 640 - p.Width - 16
		}
	}

	// Shoot recover() beam
	canShoot := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyJ)
	if p.ChanTimer > 0 {
		// Auto fire in chan mode
		canShoot = ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyJ)
	}

	if canShoot && p.ShootCooldown == 0 {
		p.Shoot(bullets)
	}
}

func (p *Player) Shoot(bullets *[]*Bullet) {
	if p.ChanTimer > 0 {
		p.ShootCooldown = 12
		// Triple shot
		*bullets = append(*bullets,
			NewHeroBullet(p.X+p.Width/2-1, p.Y-8),
			&Bullet{X: p.X, Y: p.Y - 8, Vy: -6.5, IsEnemy: false, Active: true, Width: 3, Height: 10},
			&Bullet{X: p.X + p.Width - 3, Y: p.Y - 8, Vy: -6.5, IsEnemy: false, Active: true, Width: 3, Height: 10},
		)
	} else {
		p.ShootCooldown = 18
		*bullets = append(*bullets, NewHeroBullet(p.X+p.Width/2-1, p.Y-8))
	}
	audio.GlobalAudio.PlayLaser()
}

func (p *Player) TakeDamage() bool {
	if p.ShieldTimer > 0 || p.InvulnTimer > 0 {
		return false
	}
	p.Lives--
	p.InvulnTimer = 90 // 1.5 seconds invulnerability
	audio.GlobalAudio.PlayExplosion()
	return true
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Blinking during invulnerability
	if p.InvulnTimer > 0 && (p.InvulnTimer/6)%2 == 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)

	if p.ShieldTimer > 0 {
		screen.DrawImage(assets.LoadedSprites.PlayerShielded, op)
	} else {
		screen.DrawImage(assets.LoadedSprites.Player, op)
	}
}
