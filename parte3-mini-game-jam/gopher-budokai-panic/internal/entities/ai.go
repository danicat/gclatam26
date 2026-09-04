package entities

import (
	"math"
	"math/rand"
)

type AIController struct {
	fighter    *Fighter
	rng        *rand.Rand
	actionTime float64
	mashTimer  float64
}

func NewAIController(f *Fighter) *AIController {
	return &AIController{
		fighter: f,
		rng:     rand.New(rand.NewSource(99)),
	}
}

func (ai *AIController) Update(dt float64, player *Fighter) (spawnBlast *KiBlast, spawnBeam *SuperBeam) {
	f := ai.fighter
	if f == nil || player == nil {
		return nil, nil
	}

	dist := math.Hypot(player.X-f.X, player.Y-f.Y)
	ai.actionTime -= dt
	ai.mashTimer -= dt

	// 1. If CPU is in PANIC!: Rapidly mash recover!
	if f.Panic.IsPanicked {
		if ai.mashTimer <= 0 {
			ai.mashTimer = 0.12 // Mash every 120ms
			if f.Panic.TryMashRecover() {
				f.TriggerRecoverKiai(player)
				return nil, nil
			}
		}
		// If cornered while panicked, try emergency vanish
		if ai.rng.Float64() < 0.05 {
			f.PerformVanish(player.X, player.Y, player.FacingLeft)
		}
		return nil, nil
	}

	// 2. If PLAYER is in PANIC!: Capitalize and attack aggressively!
	if player.Panic.IsPanicked {
		if dist > 80.0 {
			// Super beam finisher if CPU has enough Ki
			if f.Ki >= 35.0 && f.ActiveBeam == nil && ai.rng.Float64() < 0.04 {
				return nil, f.StartSuperBeam()
			}
			// Dragon Dash straight in to punish
			if f.Ki >= 10.0 && f.State != StateDragonDash {
				f.StartDragonDash(player.X, player.Y)
			}
		} else {
			// Close range melee barrage
			f.StartMelee(player.X, player.Y)
		}
		return nil, nil
	}

	// 3. Normal BT3 Tactical Combat
	if ai.actionTime > 0 {
		return nil, nil
	}
	ai.actionTime = 0.20 + ai.rng.Float64()*0.25

	// Vanish counter if player is lunging in melee
	if player.State == StateMelee && dist < 60.0 && ai.rng.Float64() < 0.35 {
		f.PerformVanish(player.X, player.Y, player.FacingLeft)
		return nil, nil
	}

	// Long range actions
	if dist > 260.0 {
		roll := ai.rng.Float64()
		if f.Ki < 40.0 || roll < 0.40 {
			// Charge Ki
			f.State = StateChargeKi
			ai.actionTime = 0.50
		} else if roll < 0.70 && f.Ki >= 10.0 {
			// Dragon Dash in
			f.StartDragonDash(player.X, player.Y)
		} else if roll < 0.90 && f.Ki >= 5.0 {
			// Rapid Ki Blasts
			return f.StartKiBlast(player.X, player.Y), nil
		} else if f.Ki >= 35.0 && f.ActiveBeam == nil {
			// Super Beam
			return nil, f.StartSuperBeam()
		}
		return nil, nil
	}

	// Mid range actions (80px - 260px)
	if dist > 80.0 {
		roll := ai.rng.Float64()
		if roll < 0.35 && f.Ki >= 10.0 {
			f.StartDragonDash(player.X, player.Y)
		} else if roll < 0.65 && f.Ki >= 5.0 {
			return f.StartKiBlast(player.X, player.Y), nil
		} else {
			// Maneuver
			dirX := 0.0
			if f.X < player.X {
				dirX = 1.0
			} else {
				dirX = -1.0
			}
			dirY := 0.0
			if f.Y < player.Y {
				dirY = 0.8
			} else {
				dirY = -0.8
			}
			f.StartMove(dirX, dirY)
		}
		return nil, nil
	}

	// Close range actions (< 80px)
	roll := ai.rng.Float64()
	if roll < 0.60 {
		f.StartMelee(player.X, player.Y)
	} else if roll < 0.80 && f.Ki >= 15.0 {
		f.PerformVanish(player.X, player.Y, player.FacingLeft)
	} else {
		// Back away slightly
		dirX := 1.0
		if f.X > player.X {
			dirX = 1.0
		} else {
			dirX = -1.0
		}
		f.StartMove(dirX, 0)
	}

	return nil, nil
}
