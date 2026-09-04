package entity

import (
	"testing"
)

func TestBulletUpdate(t *testing.T) {
	b := NewHeroBullet(100, 100)
	initialY := b.Y
	b.Update()
	if b.Y >= initialY {
		t.Errorf("expected hero bullet Y to decrease, got %f -> %f", initialY, b.Y)
	}

	b.Y = -20
	b.Update()
	if b.Active {
		t.Errorf("expected bullet outside screen bounds to be deactivated")
	}
}

func TestBarrierDestruction(t *testing.T) {
	bar := NewBarrier(100, 100, "defer test()")
	bullet := NewPanicBullet(101, 101)

	hit := bar.CheckBulletCollision(bullet)
	if !hit {
		t.Fatalf("expected bullet to collide with barrier")
	}
	if bullet.Active {
		t.Errorf("expected bullet to be deactivated after hitting barrier")
	}
	if bar.Grid[0][0].Alive {
		t.Errorf("expected barrier block [0][0] to be destroyed")
	}
}

func TestPlayerDamageAndShield(t *testing.T) {
	p := NewPlayer(100, 100)
	p.ShieldTimer = 100

	damaged := p.TakeDamage()
	if damaged {
		t.Errorf("expected player with shield not to take damage")
	}
	if p.Lives != 3 {
		t.Errorf("expected 3 lives, got %d", p.Lives)
	}

	p.ShieldTimer = 0
	p.InvulnTimer = 0
	damaged = p.TakeDamage()
	if !damaged {
		t.Errorf("expected unshielded player to take damage")
	}
	if p.Lives != 2 {
		t.Errorf("expected 2 lives, got %d", p.Lives)
	}
}
