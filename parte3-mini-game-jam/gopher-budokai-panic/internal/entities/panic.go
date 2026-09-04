package entities

// PanicState manages the "Panic!!! (& recover?)" mechanics for a fighter.
type PanicState struct {
	Meter          float64 // 0.0 to 100.0%
	IsPanicked     bool    // True when Meter reaches 100.0
	PanicTimer     float64 // Duration remaining in panic state
	RecoverEffort  float64 // 0.0 to 100.0 (accumulated by mashing recover)
	ZCounterWindow float64 // Active parry / instant transmission counter window
	RecoverFlash   float64 // Visual shockwave flash indicator
}

const (
	MaxPanicDuration = 3.5 // Seconds panic lasts if not recovered
	RecoverPerMash   = 28.0 // Effort granted per mash input
	PanicDecayRate   = 6.0  // Natural panic decrease per second when not taking damage
)

func (ps *PanicState) AddPanic(amount float64) {
	if ps.IsPanicked {
		return
	}
	ps.Meter += amount
	if ps.Meter >= 100.0 {
		ps.Meter = 100.0
		ps.IsPanicked = true
		ps.PanicTimer = MaxPanicDuration
		ps.RecoverEffort = 0.0
	}
}

func (ps *PanicState) Update(dt float64) bool {
	if ps.RecoverFlash > 0 {
		ps.RecoverFlash -= dt
		if ps.RecoverFlash < 0 {
			ps.RecoverFlash = 0
		}
	}

	if ps.ZCounterWindow > 0 {
		ps.ZCounterWindow -= dt
		if ps.ZCounterWindow < 0 {
			ps.ZCounterWindow = 0
		}
	}

	if ps.IsPanicked {
		ps.PanicTimer -= dt
		// Naturally recovers after timer expires
		if ps.PanicTimer <= 0 {
			ps.ForceRecover()
			return true
		}
		return false
	}

	// Gradual calm-down decay
	if ps.Meter > 0 {
		ps.Meter -= PanicDecayRate * dt
		if ps.Meter < 0 {
			ps.Meter = 0
		}
	}
	return false
}

// TryMashRecover processes a mash keypress to build recovery effort.
// Returns true if the recovery threshold was met and panic was broken!
func (ps *PanicState) TryMashRecover() bool {
	if !ps.IsPanicked {
		return false
	}
	ps.RecoverEffort += RecoverPerMash
	if ps.RecoverEffort >= 100.0 {
		ps.ForceRecover()
		return true
	}
	return false
}

// ForceRecover clears panic status and triggers recovery flash.
func (ps *PanicState) ForceRecover() {
	ps.IsPanicked = false
	ps.Meter = 0.0
	ps.PanicTimer = 0.0
	ps.RecoverEffort = 0.0
	ps.RecoverFlash = 0.4
}

// ActivateZCounter initiates a brief window (e.g. 0.2s) where any incoming hit
// triggers an Instant Transmission counter-attack!
func (ps *PanicState) ActivateZCounter() {
	ps.ZCounterWindow = 0.22
}
