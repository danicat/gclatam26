package gfx

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TileState int

const (
	TileNormal TileState = iota
	TileWarning
	TileHazard
)

type FloorTile struct {
	Col, Row int
	X, Y     float64
	W, H     float64
	State    TileState
	Timer    float64
	BaseHue  float64
}

type DiscoFloor struct {
	Cols      int
	Rows      int
	TileW     float64
	TileH     float64
	OriginX   float64
	OriginY   float64
	Tiles     []FloorTile
	BeatTime  float64
	BPM       float64
	StepCount int
}

func NewDiscoFloor(originX, originY, w, h float64, cols, rows int, bpm float64) *DiscoFloor {
	df := &DiscoFloor{
		Cols:      cols,
		Rows:      rows,
		TileW:     w / float64(cols),
		TileH:     h / float64(rows),
		OriginX:   originX,
		OriginY:   originY,
		Tiles:     make([]FloorTile, cols*rows),
		BPM:       bpm,
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			df.Tiles[idx] = FloorTile{
				Col:     c,
				Row:     r,
				X:       originX + float64(c)*df.TileW,
				Y:       originY + float64(r)*df.TileH,
				W:       df.TileW,
				H:       df.TileH,
				State:   TileNormal,
				BaseHue: float64((c*37 + r*53) % 360),
			}
		}
	}
	return df
}

func (df *DiscoFloor) Update(dt float64) {
	df.BeatTime += dt
	beatDuration := 60.0 / df.BPM
	if df.BeatTime >= beatDuration {
		df.BeatTime -= beatDuration
		df.StepCount++
	}

	for i := range df.Tiles {
		t := &df.Tiles[i]
		if t.Timer > 0 {
			t.Timer -= dt
			if t.Timer <= 0 {
				if t.State == TileWarning {
					t.State = TileHazard
					t.Timer = beatDuration * 1.5 // Hazard lasts 1.5 beats
				} else {
					t.State = TileNormal
					t.Timer = 0
				}
			}
		}
	}
}

// TriggerHazardTile sets random floor tiles to warning then hazard on beat.
func (df *DiscoFloor) TriggerHazardTile(col, row int, warningDuration float64) {
	if col >= 0 && col < df.Cols && row >= 0 && row < df.Rows {
		idx := row*df.Cols + col
		df.Tiles[idx].State = TileWarning
		df.Tiles[idx].Timer = warningDuration
	}
}

// IsHazardAt checks if an entity at world coordinates (x, y) is on an active hazard tile.
func (df *DiscoFloor) IsHazardAt(x, y float64) bool {
	if x < df.OriginX || y < df.OriginY {
		return false
	}
	c := int((x - df.OriginX) / df.TileW)
	r := int((y - df.OriginY) / df.TileH)
	if c >= 0 && c < df.Cols && r >= 0 && r < df.Rows {
		idx := r*df.Cols + c
		return df.Tiles[idx].State == TileHazard
	}
	return false
}

func (df *DiscoFloor) Draw(screen *ebiten.Image) {
	beatFrac := df.BeatTime / (60.0 / df.BPM)
	beatPulse := math.Sin(beatFrac * math.Pi)

	for i := range df.Tiles {
		t := &df.Tiles[i]
		var fillCol color.RGBA

		switch t.State {
		case TileWarning:
			// Rapid flashing amber/yellow
			blink := math.Sin(df.BeatTime * 25.0)
			if blink > 0 {
				fillCol = color.RGBA{R: 255, G: 190, B: 0, A: 220}
			} else {
				fillCol = color.RGBA{R: 200, G: 80, B: 0, A: 180}
			}
		case TileHazard:
			// Blazing electrical red/crimson
			blink := math.Sin(df.BeatTime * 35.0)
			if blink > 0 {
				fillCol = color.RGBA{R: 255, G: 30, B: 30, A: 240}
			} else {
				fillCol = color.RGBA{R: 180, G: 0, B: 20, A: 200}
			}
		default:
			// Groovy cycling 70s neon palette
			pattern := (t.Col + t.Row + df.StepCount) % 4
			baseIntensity := 0.65 + 0.35*beatPulse

			switch pattern {
			case 0: // Electric Cyan
				fillCol = color.RGBA{R: 0, G: uint8(220 * baseIntensity), B: uint8(255 * baseIntensity), A: 200}
			case 1: // Hot Magenta
				fillCol = color.RGBA{R: uint8(255 * baseIntensity), G: 0, B: uint8(160 * baseIntensity), A: 200}
			case 2: // Golden Disco
				fillCol = color.RGBA{R: uint8(255 * baseIntensity), G: uint8(215 * baseIntensity), B: 0, A: 200}
			case 3: // Deep Velvet Purple
				fillCol = color.RGBA{R: uint8(160 * baseIntensity), G: 0, B: uint8(240 * baseIntensity), A: 200}
			}
		}

		// Draw inner tile
		inset := float32(1.5)
		vector.DrawFilledRect(screen,
			float32(t.X)+inset, float32(t.Y)+inset,
			float32(t.W)-inset*2, float32(t.H)-inset*2,
			fillCol, false)

		// Tile border (dark neon grid lines)
		borderColor := color.RGBA{R: 20, G: 10, B: 30, A: 255}
		vector.StrokeRect(screen, float32(t.X), float32(t.Y), float32(t.W), float32(t.H), 1, borderColor, false)
	}
}
