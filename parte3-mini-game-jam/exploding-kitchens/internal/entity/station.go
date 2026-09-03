package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// StationKind differentiates the appliance function and visual style.
type StationKind int

const (
	StationPressureCooker StationKind = iota
	StationDeepFryer
	StationMicrowave
	StationStoveTop
)

func (k StationKind) Name() string {
	switch k {
	case StationPressureCooker:
		return "PRESSURE COOKER"
	case StationDeepFryer:
		return "DEEP FRYER"
	case StationMicrowave:
		return "MICROWAVE"
	case StationStoveTop:
		return "STOVE TOP"
	default:
		return "STATION"
	}
}

// StationState represents the operational phase of an appliance.
type StationState int

const (
	StateIdle StationState = iota
	StateCooking
	StateWarning
	StatePanic
	StateExploded
)

// Station is an interactive cooking station prone to overheating and catastrophic detonation.
type Station struct {
	X, Y        float64
	W, H        float64
	Kind        StationKind
	State       StationState
	Timer       float64
	MaxTime     float64
	CatBoost    bool // True if a cat is currently sitting on this appliance
	drawOpts    ebiten.DrawImageOptions
	pulseTimer  float64
	sparkTimer  float64
}

// NewStation creates an appliance station with preset cook times.
func NewStation(x, y float64, kind StationKind) *Station {
	baseTime := 14.0 + rand.Float64()*6.0
	switch kind {
	case StationPressureCooker:
		baseTime = 12.0
	case StationDeepFryer:
		baseTime = 15.0
	case StationMicrowave:
		baseTime = 10.0
	case StationStoveTop:
		baseTime = 16.0
	}

	return &Station{
		X:       x,
		Y:       y,
		W:       26,
		H:       22,
		Kind:    kind,
		State:   StateCooking,
		Timer:   rand.Float64() * 3.0, // Stagger initial timers
		MaxTime: baseTime,
	}
}

// Progress returns the danger percentage (0.0 to 1.0).
func (s *Station) Progress() float64 {
	if s.MaxTime <= 0 {
		return 0
	}
	p := s.Timer / s.MaxTime
	if p > 1.0 {
		p = 1.0
	}
	return p
}

// IsClutch returns true if the appliance is in the final 15% before detonation.
func (s *Station) IsClutch() bool {
	return s.State == StatePanic && s.Progress() >= 0.85
}

// Update advances station timers and updates state.
func (s *Station) Update(dt float64) (exploded bool) {
	s.pulseTimer += dt * 8.0
	s.sparkTimer += dt

	if s.State == StateExploded {
		return false
	}

	// Cats speed up cook time by 2.5x!
	rate := 1.0
	if s.CatBoost {
		rate = 2.5
	}
	s.Timer += dt * rate

	progress := s.Progress()

	if progress >= 1.0 {
		s.State = StateExploded
		return true // Trigger explosion event!
	} else if progress >= 0.80 {
		s.State = StatePanic
	} else if progress >= 0.55 {
		s.State = StateWarning
	} else {
		s.State = StateCooking
	}

	return false
}

// CanRecover checks if the provided tool can defuse or repair the station.
func (s *Station) CanRecover(tool ToolType) bool {
	if s.State == StateExploded {
		return tool == ToolWrench
	}

	if s.State == StateWarning || s.State == StatePanic {
		switch s.Kind {
		case StationPressureCooker:
			return tool == ToolIce || tool == ToolNone // Can vent with bare hands or chill with ice
		case StationDeepFryer:
			return tool == ToolExtinguisher
		case StationMicrowave:
			return tool == ToolNone || tool == ToolIce
		case StationStoveTop:
			return tool == ToolExtinguisher || tool == ToolIce
		}
	}
	return false
}

// ResetCooking returns the station to normal cooking or idle state after recovery.
func (s *Station) ResetCooking() {
	s.State = StateCooking
	s.Timer = 0
	s.MaxTime = 12.0 + rand.Float64()*8.0
}

// Repair restores an exploded station back to operational state.
func (s *Station) Repair() {
	s.State = StateCooking
	s.Timer = 0
	s.MaxTime = 14.0 + rand.Float64()*6.0
}

// Draw renders the station appliance, status lights, and progress meter.
func (s *Station) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	// Base Counter / Cabinet
	s.drawRect(screen, pixelImg, s.X, s.Y, s.W, s.H, color.RGBA{50, 45, 65, 255})
	s.drawRect(screen, pixelImg, s.X+1, s.Y+1, s.W-2, s.H-2, color.RGBA{80, 75, 95, 255})

	if s.State == StateExploded {
		// Charred / Burned wreckage
		s.drawRect(screen, pixelImg, s.X+2, s.Y+2, s.W-4, s.H-4, color.RGBA{30, 25, 30, 255})
		// Broken sparks / X mark
		s.drawRect(screen, pixelImg, s.X+6, s.Y+6, s.W-12, s.H-12, color.RGBA{80, 40, 40, 255})
		DrawToolIcon(screen, pixelImg, ToolWrench, s.X+s.W/2, s.Y+s.H/2)
		return
	}

	// Appliance specific top graphics
	switch s.Kind {
	case StationPressureCooker:
		// Silver cylindrical pot
		s.drawRect(screen, pixelImg, s.X+4, s.Y+4, s.W-8, s.H-8, color.RGBA{190, 200, 215, 255})
		// Pressure valve / gauge
		s.drawRect(screen, pixelImg, s.X+s.W/2-2, s.Y+2, 4, 3, color.RGBA{220, 160, 40, 255})

	case StationDeepFryer:
		// Oil vat with bubbling yellow/orange
		s.drawRect(screen, pixelImg, s.X+3, s.Y+3, s.W-6, s.H-6, color.RGBA{150, 155, 170, 255})
		s.drawRect(screen, pixelImg, s.X+5, s.Y+5, s.W-10, s.H-10, color.RGBA{230, 160, 20, 255})

	case StationMicrowave:
		// Dark glass window + keypad
		s.drawRect(screen, pixelImg, s.X+3, s.Y+3, s.W-12, s.H-6, color.RGBA{35, 40, 55, 255})
		s.drawRect(screen, pixelImg, s.X+s.W-8, s.Y+3, 5, s.H-6, color.RGBA{110, 120, 140, 255})

	case StationStoveTop:
		// 2 round burner coils
		s.drawRect(screen, pixelImg, s.X+4, s.Y+5, 6, 6, color.RGBA{40, 35, 45, 255})
		s.drawRect(screen, pixelImg, s.X+s.W-10, s.Y+5, 6, 6, color.RGBA{40, 35, 45, 255})
		if s.State != StateIdle {
			s.drawRect(screen, pixelImg, s.X+5, s.Y+6, 4, 4, color.RGBA{240, 70, 30, 255})
			s.drawRect(screen, pixelImg, s.X+s.W-9, s.Y+6, 4, 4, color.RGBA{240, 70, 30, 255})
		}
	}

	// Danger Outline Pulse
	if s.State == StatePanic {
		alpha := uint8(160 + math.Sin(s.pulseTimer)*80)
		s.drawOutline(screen, pixelImg, s.X-1, s.Y-1, s.W+2, s.H+2, color.RGBA{255, 30, 30, alpha})
	} else if s.State == StateWarning {
		s.drawOutline(screen, pixelImg, s.X-1, s.Y-1, s.W+2, s.H+2, color.RGBA{255, 200, 30, 180})
	}

	// Floating Danger Progress Bar
	barW := s.W
	barH := 3.0
	barX := s.X
	barY := s.Y - 5.0

	// Bar background
	s.drawRect(screen, pixelImg, barX, barY, barW, barH, color.RGBA{20, 15, 30, 200})

	// Fill color based on progress
	prog := s.Progress()
	fillW := barW * prog
	fillCol := color.RGBA{70, 200, 90, 255} // Green
	if prog >= 0.80 {
		fillCol = color.RGBA{255, 40, 40, 255} // Critical Red
	} else if prog >= 0.55 {
		fillCol = color.RGBA{255, 200, 30, 255} // Warning Yellow
	}
	s.drawRect(screen, pixelImg, barX, barY, fillW, barH, fillCol)
}

func (s *Station) drawRect(screen *ebiten.Image, pixelImg *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	s.drawOpts.GeoM.Reset()
	s.drawOpts.GeoM.Scale(w, h)
	s.drawOpts.GeoM.Translate(x, y)

	af := float32(c.A) / 255.0
	s.drawOpts.ColorScale.Reset()
	s.drawOpts.ColorScale.Scale(
		(float32(c.R)/255.0)*af,
		(float32(c.G)/255.0)*af,
		(float32(c.B)/255.0)*af,
		af,
	)
	screen.DrawImage(pixelImg, &s.drawOpts)
}

func (s *Station) drawOutline(screen *ebiten.Image, pixelImg *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	s.drawRect(screen, pixelImg, x, y, w, 1, c)
	s.drawRect(screen, pixelImg, x, y+h-1, w, 1, c)
	s.drawRect(screen, pixelImg, x, y, 1, h, c)
	s.drawRect(screen, pixelImg, x+w-1, y, 1, h, c)
}
