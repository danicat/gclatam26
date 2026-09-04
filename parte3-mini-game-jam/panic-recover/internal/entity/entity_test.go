package entity

import (
	"testing"
)

func TestPlayerTakeDamageAndPanicSurge(t *testing.T) {
	ps := NewParticleSystem(10)
	p := NewPlayer(100, 100)

	// Normal damage
	panicked := p.TakeDamage(30, ps)
	if panicked {
		t.Errorf("expected normal damage not to panic, but got panicked = true")
	}
	if p.HP != 70 {
		t.Errorf("expected HP = 70, got %f", p.HP)
	}

	// Fatal damage triggers panic surge
	panicked = p.TakeDamage(80, ps)
	if !panicked {
		t.Errorf("expected fatal damage to trigger panic, got panicked = false")
	}
	if !p.InPanic {
		t.Errorf("expected player.InPanic = true")
	}
	if p.PanicTimer != 5.0 {
		t.Errorf("expected PanicTimer = 5.0, got %f", p.PanicTimer)
	}
	if p.IsDead {
		t.Errorf("expected player to survive during panic surge")
	}

	// During panic, damage should be ignored (invulnerable)
	p.TakeDamage(50, ps)
	if p.HP != 1 {
		t.Errorf("expected HP to stay at 1 during panic, got %f", p.HP)
	}
}

func TestPlayerPanicTimeout(t *testing.T) {
	ps := NewParticleSystem(10)
	p := NewPlayer(100, 100)

	// Enter panic
	p.TakeDamage(100, ps)

	// Simulate 5.1 seconds of update
	p.Update(5.1, 640, 360, ps)
	if !p.IsDead {
		t.Errorf("expected player to die after panic timer expires")
	}
}

func TestPlayerRecoverFromPanic(t *testing.T) {
	ps := NewParticleSystem(10)
	bm := NewBulletManager(10)
	p := NewPlayer(100, 100)

	// Enter panic
	p.TakeDamage(100, ps)

	// Collect recover drop
	pickup := &Pickup{Type: PickupTypeRecover, Active: true}
	p.CollectPickup(pickup, ps, bm)

	if p.InPanic {
		t.Errorf("expected player to no longer be in panic")
	}
	if p.HP != 60 {
		t.Errorf("expected HP to recover to 60, got %f", p.HP)
	}
	if p.IsDead {
		t.Errorf("expected player to be alive")
	}
}

func TestPlayerMutexShieldAndDrones(t *testing.T) {
	ps := NewParticleSystem(10)
	bm := NewBulletManager(10)
	p := NewPlayer(100, 100)

	// Collect Mutex pickup
	p.CollectPickup(&Pickup{Type: PickupTypeMutex, Active: true}, ps, bm)
	if p.ShieldTimer <= 0 {
		t.Errorf("expected shield timer > 0, got %f", p.ShieldTimer)
	}

	// Damage while shield is active should deal 0 damage
	p.TakeDamage(50, ps)
	if p.HP != 100 {
		t.Errorf("expected HP = 100 with shield, got %f", p.HP)
	}

	// Collect Worker drones
	p.CollectPickup(&Pickup{Type: PickupTypeWorker, Active: true}, ps, bm)
	if p.DroneCount != 1 {
		t.Errorf("expected 1 drone, got %d", p.DroneCount)
	}
	p.CollectPickup(&Pickup{Type: PickupTypeWorker, Active: true}, ps, bm)
	if p.DroneCount != 2 {
		t.Errorf("expected 2 drones, got %d", p.DroneCount)
	}
	// Max 2 drones
	p.CollectPickup(&Pickup{Type: PickupTypeWorker, Active: true}, ps, bm)
	if p.DroneCount != 2 {
		t.Errorf("expected drone count capped at 2, got %d", p.DroneCount)
	}
}

func TestPlayerBoundaryClamping(t *testing.T) {
	ps := NewParticleSystem(10)
	p := NewPlayer(0, 0)
	p.Update(0.1, 640, 360, ps)

	if p.X < 20 || p.Y < 30 {
		t.Errorf("expected player clamped inside boundaries, got X=%f, Y=%f", p.X, p.Y)
	}
}

func TestBossSpawnAndDamage(t *testing.T) {
	boss := NewBoss()
	boss.Spawn(1, 640)
	if !boss.Active {
		t.Fatalf("expected boss to be active after spawn")
	}
	if boss.HP <= 0 {
		t.Fatalf("expected boss HP > 0, got %f", boss.HP)
	}

	ps := NewParticleSystem(10)
	// Normal damage
	defeated := boss.TakeDamage(50, ps)
	if defeated {
		t.Errorf("expected boss not to be defeated by 50 dmg")
	}

	// Fatal damage
	defeated = boss.TakeDamage(1000, ps)
	if !defeated {
		t.Errorf("expected boss to be defeated by 1000 dmg")
	}
	if boss.Active {
		t.Errorf("expected boss to be inactive after defeat")
	}
}

func TestEnemyManagerTierScaling(t *testing.T) {
	em := NewEnemyManager(10)
	em.SetTier(2)
	em.Spawn(EnemyTypeNilPointer, 100, 100, false)

	enemies := em.Enemies()
	if !enemies[0].Active {
		t.Fatalf("expected spawned enemy to be active")
	}
	// At tier 2, hpMultiplier is 1.25, base HP for NilPointer is 30 -> 37.5
	if enemies[0].MaxHP <= 30 {
		t.Errorf("expected scaled HP > 30, got %f", enemies[0].MaxHP)
	}
}
