package entities

import (
	"math"

	"gopher-budokai-panic/internal/audio"
	"gopher-budokai-panic/internal/gfx"
)

type FighterState int

const (
	StateIdle FighterState = iota
	StateMove
	StateDragonDash
	StateChargeKi
	StateMelee
	StateBlast
	StateBeam
	StateVanish
	StateHurt
	StateKnockback
)

type Fighter struct {
	ID          int
	Type        gfx.FighterType
	X, Y        float64
	VX, VY      float64
	FacingLeft  bool
	Health      float64
	MaxHealth   float64
	Ki          float64
	MaxKi       float64
	IsSparking  bool
	SparkTimer  float64
	State       FighterState
	StateTimer  float64
	ComboStep   int
	ComboTimer  float64
	HoverPhase  float64
	HitStop     float64
	Panic       PanicState
	ActiveBeam  *SuperBeam
	LastHitByX  float64
}

func NewFighter(id int, fType gfx.FighterType, startX, startY float64, facingLeft bool) *Fighter {
	return &Fighter{
		ID:         id,
		Type:       fType,
		X:          startX,
		Y:          startY,
		FacingLeft: facingLeft,
		Health:     1000.0,
		MaxHealth:  1000.0,
		Ki:         40.0,
		MaxKi:      100.0,
		State:      StateIdle,
	}
}

func (f *Fighter) Update(dt float64, oppX, oppY float64, arenaW, arenaH float64) {
	// Hit-stop freeze frames on heavy impact
	if f.HitStop > 0 {
		f.HitStop -= dt
		return
	}

	f.HoverPhase += 4.0 * dt

	// Face opponent automatically unless dashing
	if f.State != StateDragonDash && f.State != StateKnockback {
		f.FacingLeft = f.X > oppX
	}

	// Panic state tick
	if f.Panic.Update(dt) {
		// Just recovered naturally
		audio.Get().PlayRecoverKiai()
	}

	// Sparking mode countdown
	if f.IsSparking {
		f.SparkTimer -= dt
		if f.SparkTimer <= 0 {
			f.IsSparking = false
			f.Ki = 30.0
		}
	}

	// State machine update
	switch f.State {
	case StateIdle:
		f.VX *= 0.85
		f.VY *= 0.85
		// Gentle hovering floating bob
		f.Y += math.Sin(f.HoverPhase) * 0.4

	case StateMove:
		f.X += f.VX * dt
		f.Y += f.VY * dt
		f.VX *= 0.88
		f.VY *= 0.88
		if math.Abs(f.VX) < 5 && math.Abs(f.VY) < 5 {
			f.State = StateIdle
		}

	case StateDragonDash:
		f.X += f.VX * dt
		f.Y += f.VY * dt
		f.StateTimer -= dt
		// Consume Ki continuously during Dragon Dash
		f.Ki -= 18.0 * dt
		if f.Ki <= 0 {
			f.Ki = 0
			// Overheated dash triggers panic risk!
			f.Panic.AddPanic(15.0)
			f.State = StateIdle
		}
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}

	case StateChargeKi:
		f.VX = 0
		f.VY = 0
		chargeRate := 35.0
		if f.IsSparking {
			chargeRate = 15.0
		}
		f.Ki += chargeRate * dt
		audio.Get().PlayCharge()
		if f.Ki >= f.MaxKi {
			f.Ki = f.MaxKi
			if !f.IsSparking {
				f.IsSparking = true
				f.SparkTimer = 10.0 // 10 seconds of Sparking!
			}
		}

	case StateMelee:
		f.StateTimer -= dt
		f.X += f.VX * dt
		f.Y += f.VY * dt
		f.VX *= 0.82
		f.VY *= 0.82
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}

	case StateBlast:
		f.StateTimer -= dt
		f.VX *= 0.8
		f.VY *= 0.8
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}

	case StateBeam:
		f.VX = 0
		f.VY = 0
		if f.ActiveBeam != nil {
			f.ActiveBeam.Update(dt, f.X, f.Y)
			if f.ActiveBeam.State == BeamStateDone {
				f.ActiveBeam = nil
				f.State = StateIdle
			}
		} else {
			f.State = StateIdle
		}

	case StateVanish:
		f.StateTimer -= dt
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}

	case StateHurt:
		f.StateTimer -= dt
		f.X += f.VX * dt
		f.Y += f.VY * dt
		f.VX *= 0.85
		f.VY *= 0.85
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}

	case StateKnockback:
		f.StateTimer -= dt
		f.X += f.VX * dt
		f.Y += f.VY * dt
		f.VX *= 0.94
		f.VY *= 0.94
		if f.StateTimer <= 0 {
			f.State = StateIdle
		}
	}

	// Arena boundaries clamping
	margin := 24.0
	if f.X < margin {
		f.X = margin
		if f.State == StateKnockback {
			// Wall slam panic increase!
			f.Panic.AddPanic(15.0)
			f.VX = -f.VX * 0.5
		}
	}
	if f.X > arenaW-margin {
		f.X = arenaW - margin
		if f.State == StateKnockback {
			f.Panic.AddPanic(15.0)
			f.VX = -f.VX * 0.5
		}
	}
	if f.Y < margin {
		f.Y = margin
	}
	if f.Y > arenaH-margin*1.5 {
		f.Y = arenaH - margin*1.5
	}
}

// StartMove initiates flight in 8 directions.
func (f *Fighter) StartMove(dirX, dirY float64) {
	if f.State != StateIdle && f.State != StateMove {
		return
	}
	speed := 210.0
	if f.IsSparking {
		speed = 280.0
	}
	if f.Panic.IsPanicked {
		// Shaky frantic movement during panic
		speed *= 1.3
	}
	f.VX = dirX * speed
	f.VY = dirY * speed
	f.State = StateMove
}

// StartDragonDash executes the high-speed supersonic flight toward the target.
func (f *Fighter) StartDragonDash(targetX, targetY float64) {
	if f.State == StateHurt || f.State == StateKnockback || f.Ki < 10.0 {
		return
	}
	dx := targetX - f.X
	dy := targetY - f.Y
	dist := math.Hypot(dx, dy)
	if dist == 0 {
		return
	}
	speed := 520.0
	if f.IsSparking {
		speed = 680.0
	}
	f.VX = (dx / dist) * speed
	f.VY = (dy / dist) * speed
	f.State = StateDragonDash
	f.StateTimer = 0.8
	f.Ki -= 8.0
}

// StartMelee initiates punch/kick combo strings.
func (f *Fighter) StartMelee(oppX, oppY float64) bool {
	if f.State == StateHurt || f.State == StateKnockback || f.State == StateBeam {
		return false
	}
	dist := math.Hypot(oppX-f.X, oppY-f.Y)
	if dist > 55.0 {
		return false
	}

	f.ComboStep++
	if f.ComboStep > 3 {
		f.ComboStep = 1
	}

	f.State = StateMelee
	f.StateTimer = 0.22

	// Lunge forward slightly
	dir := 1.0
	if f.FacingLeft {
		dir = -1.0
	}
	f.VX = dir * 140.0
	f.VY = 0
	return true
}

// StartKiBlast shoots a rapid energy pellet.
func (f *Fighter) StartKiBlast(targetX, targetY float64) *KiBlast {
	if f.State == StateHurt || f.State == StateKnockback || f.State == StateBeam || f.Ki < 5.0 {
		return nil
	}
	dx := targetX - f.X
	dy := targetY - f.Y
	dist := math.Hypot(dx, dy)
	if dist == 0 {
		dist = 1
	}
	dirX := dx / dist
	dirY := dy / dist

	f.Ki -= 5.0
	f.State = StateBlast
	f.StateTimer = 0.16
	audio.Get().PlayBlast()

	spawnX := f.X + dirX*20.0
	spawnY := f.Y + dirY*10.0
	return NewKiBlast(spawnX, spawnY, dirX, dirY, f.ID, f.IsSparking)
}

// StartSuperBeam fires the colossal beam (Kamehameha / Final Flash).
func (f *Fighter) StartSuperBeam() *SuperBeam {
	if f.State == StateHurt || f.State == StateKnockback || f.ActiveBeam != nil || f.Ki < 35.0 {
		return nil
	}
	f.Ki -= 35.0
	f.State = StateBeam
	dirX := 1.0
	if f.FacingLeft {
		dirX = -1.0
	}
	beam := NewSuperBeam(f.ID, f.X, f.Y, dirX, f.IsSparking)
	f.ActiveBeam = beam
	audio.Get().PlayBeam()
	return beam
}

// PerformVanish executes Instant Transmission to appear behind the opponent.
func (f *Fighter) PerformVanish(targetX, targetY float64, oppFacingLeft bool) {
	if f.Ki < 15.0 && !f.Panic.IsPanicked {
		return
	}
	if !f.Panic.IsPanicked {
		f.Ki -= 15.0
	}
	// Teleport behind opponent
	offset := 38.0
	if oppFacingLeft {
		f.X = targetX + offset
	} else {
		f.X = targetX - offset
	}
	f.Y = targetY
	f.FacingLeft = !oppFacingLeft
	f.State = StateVanish
	f.StateTimer = 0.12
	f.VX = 0
	f.VY = 0
	audio.Get().PlayVanish()
}

// TakeDamage applies damage with Panic modifiers and triggers hurt states.
func (f *Fighter) TakeDamage(amount float64, isKnockback bool, fromX float64) bool {
	// Check Z-Counter window
	if f.Panic.ZCounterWindow > 0 {
		return false // Evaded / Countered!
	}

	// In PANIC! state, incoming damage is 1.5x!
	if f.Panic.IsPanicked {
		amount *= 1.5
	}

	f.Health -= amount
	if f.Health < 0 {
		f.Health = 0
	}

	f.LastHitByX = fromX
	f.HitStop = 0.08 // Hit stop impact freeze

	// Add panic meter on being hit
	panicGain := amount * 0.45
	f.Panic.AddPanic(panicGain)
	if f.Panic.IsPanicked {
		audio.Get().PlayPanicAlert()
	}

	// Knockback physics
	dir := 1.0
	if fromX > f.X {
		dir = -1.0
	}
	if isKnockback {
		f.State = StateKnockback
		f.StateTimer = 0.55
		f.VX = dir * 360.0
		f.VY = -80.0
	} else {
		f.State = StateHurt
		f.StateTimer = 0.22
		f.VX = dir * 100.0
	}

	audio.Get().PlayHit()
	return true
}

func (f *Fighter) TriggerRecoverKiai(opp *Fighter) bool {
	if !f.Panic.IsPanicked {
		return false
	}
	f.Panic.ForceRecover()
	audio.Get().PlayRecoverKiai()

	// Push away opponent with massive explosive wave shockwave
	if opp != nil {
		dx := opp.X - f.X
		dist := math.Max(math.Abs(dx), 10.0)
		dir := dx / dist
		opp.TakeDamage(30.0, true, f.X)
		opp.VX = dir * 420.0
	}
	return true
}

func (f *Fighter) GetPose() gfx.Pose {
	if f.Panic.IsPanicked {
		return gfx.PosePanic
	}
	switch f.State {
	case StateKnockback:
		return gfx.PoseKnockback
	case StateHurt:
		return gfx.PoseHurt
	case StateDragonDash:
		return gfx.PoseDash
	case StateChargeKi:
		return gfx.PoseCharge
	case StateMelee:
		if f.ComboStep%2 == 0 {
			return gfx.PoseMelee2
		}
		return gfx.PoseMelee1
	case StateBeam:
		if f.ActiveBeam != nil && f.ActiveBeam.State == BeamStateCharging {
			return gfx.PoseBeamPrep
		}
		return gfx.PoseBeamFire
	case StateMove:
		return gfx.PoseMove
	default:
		return gfx.PoseIdle
	}
}
