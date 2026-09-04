package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BeamState int

const (
	BeamStateCharging BeamState = iota
	BeamStateFiring
	BeamStateDone
)

type SuperBeam struct {
	OwnerID     int
	X, Y        float64
	DirX        float64
	State       BeamState
	ChargeTimer float64
	FireTimer   float64
	Width       float32
	Length      float32
	Color       color.RGBA
	Damage      float64
	IsClashing  bool
	ClashOffset float64 // Moves during button mashing
}

func NewSuperBeam(ownerID int, x, y, dirX float64, isSparking bool) *SuperBeam {
	col := color.RGBA{R: 90, G: 200, B: 255, A: 255} // Blue Kamehameha
	if ownerID != 0 {
		col = color.RGBA{R: 255, G: 220, B: 40, A: 255} // Gold Final Flash
	}
	w := float32(26.0)
	if isSparking {
		w = 38.0
	}

	return &SuperBeam{
		OwnerID:     ownerID,
		X:           x,
		Y:           y,
		DirX:        dirX,
		State:       BeamStateCharging,
		ChargeTimer: 0.65, // Charging buildup time
		FireTimer:   1.40, // Sustained beam duration
		Width:       w,
		Length:      640.0,
		Color:       col,
		Damage:      380.0, // Massive super damage
	}
}

func (b *SuperBeam) Update(dt float64, ownerX, ownerY float64) {
	b.X = ownerX
	b.Y = ownerY

	switch b.State {
	case BeamStateCharging:
		b.ChargeTimer -= dt
		if b.ChargeTimer <= 0 {
			b.State = BeamStateFiring
		}
	case BeamStateFiring:
		b.FireTimer -= dt
		if b.FireTimer <= 0 {
			b.State = BeamStateDone
		}
	}
}

func (b *SuperBeam) Draw(screen *ebiten.Image) {
	if b.State == BeamStateDone {
		return
	}

	x := float32(b.X)
	y := float32(b.Y)

	if b.State == BeamStateCharging {
		// Charging ki sphere in palms
		chargeProgress := float32(1.0 - (b.ChargeTimer / 0.65))
		rad := b.Width * 0.7 * chargeProgress
		if rad < 3 {
			rad = 3
		}
		halo := color.RGBA{R: b.Color.R, G: b.Color.G, B: b.Color.B, A: 160}
		vector.DrawFilledCircle(screen, x, y, rad*1.7, halo, true)
		vector.DrawFilledCircle(screen, x, y, rad, b.Color, true)
		vector.DrawFilledCircle(screen, x, y, rad*0.5, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)
		return
	}

	// BeamStateFiring: Huge horizontal beam cylinder
	endX := x + float32(b.DirX*float64(b.Length))
	if b.IsClashing {
		// Shortens to clash meeting point
		endX = float32(320.0 + b.ClashOffset)
	}

	// Outer energy corona
	coronaCol := color.RGBA{R: b.Color.R, G: b.Color.G, B: b.Color.B, A: 120}
	vector.StrokeLine(screen, x, y, endX, y, b.Width*1.7, coronaCol, false)

	// Mid energy beam
	vector.StrokeLine(screen, x, y, endX, y, b.Width, b.Color, false)

	// Piercing white core
	white := color.RGBA{R: 255, G: 255, B: 255, A: 250}
	vector.StrokeLine(screen, x, y, endX, y, b.Width*0.45, white, false)

	// Expanding burst muzzle at origin
	vector.DrawFilledCircle(screen, x, y, b.Width*1.2, b.Color, true)
	vector.DrawFilledCircle(screen, x, y, b.Width*0.7, white, true)

	// Impact head at beam tip
	vector.DrawFilledCircle(screen, endX, y, b.Width*1.1, b.Color, true)
	vector.DrawFilledCircle(screen, endX, y, b.Width*0.6, white, true)
}

// CheckBeamClash determines if two active firing beams in opposing directions clash!
func CheckBeamClash(b1, b2 *SuperBeam) bool {
	if b1 == nil || b2 == nil {
		return false
	}
	if b1.State != BeamStateFiring || b2.State != BeamStateFiring {
		return false
	}
	// Facing towards each other
	if (b1.DirX > 0 && b2.DirX < 0 && b1.X < b2.X) ||
		(b1.DirX < 0 && b2.DirX > 0 && b1.X > b2.X) {
		// Y alignment check
		if math.Abs(b1.Y-b2.Y) < 45.0 {
			b1.IsClashing = true
			b2.IsClashing = true
			return true
		}
	}
	return false
}
