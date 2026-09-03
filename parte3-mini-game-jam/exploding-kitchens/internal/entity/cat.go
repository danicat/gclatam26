package entity

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// CatState represents what the mischievous cat is currently doing.
type CatState int

const (
	CatWandering CatState = iota
	CatTargetingStation
	CatSittingOnStation
	CatShooedFleeing
	CatResting
)

// Cat is an autonomous chaos agent roaming the kitchen.
type Cat struct {
	X, Y        float64
	W, H        float64
	State       CatState
	TargetX     float64
	TargetY     float64
	TargetStn   *Station
	RestTimer   float64
	StateTimer  float64
	WalkAnim    float64
	FacingRight bool
	drawOpts    ebiten.DrawImageOptions
}

// NewCat spawns a mischievous orange cat at the designated position.
func NewCat(x, y float64) *Cat {
	return &Cat{
		X:           x,
		Y:           y,
		W:           12,
		H:           10,
		State:       CatWandering,
		TargetX:     x,
		TargetY:     y,
		FacingRight: true,
	}
}

// Update handles cat autonomous decision-making and movements.
func (c *Cat) Update(dt float64, stations []*Station, roomMinX, roomMaxX, roomMinY, roomMaxY float64) (shooed bool) {
	c.WalkAnim += dt * 8.0
	c.StateTimer += dt

	switch c.State {
	case CatResting:
		c.RestTimer -= dt
		if c.RestTimer <= 0 {
			c.State = CatWandering
			c.pickRandomTarget(roomMinX, roomMaxX, roomMinY, roomMaxY)
		}

	case CatWandering:
		if c.moveTowards(c.TargetX, c.TargetY, 35.0, dt) {
			// Arrived at destination; decide whether to target a station
			if len(stations) > 0 && rand.Float64() < 0.60 {
				c.State = CatTargetingStation
				c.TargetStn = stations[rand.Intn(len(stations))]
			} else {
				c.pickRandomTarget(roomMinX, roomMaxX, roomMinY, roomMaxY)
			}
		}

	case CatTargetingStation:
		if c.TargetStn == nil || c.TargetStn.State == StateExploded {
			c.State = CatWandering
			c.pickRandomTarget(roomMinX, roomMaxX, roomMinY, roomMaxY)
			break
		}

		stnX := c.TargetStn.X + c.TargetStn.W/2 - c.W/2
		stnY := c.TargetStn.Y + c.TargetStn.H/2 - c.H/2
		if c.moveTowards(stnX, stnY, 45.0, dt) {
			c.State = CatSittingOnStation
			c.TargetStn.CatBoost = true
			c.StateTimer = 0
		}

	case CatSittingOnStation:
		if c.TargetStn == nil || c.TargetStn.State == StateExploded {
			c.detachStation()
			c.Shoo(roomMinX, roomMinY)
			return true
		}
		// Sits for up to 6 seconds if not disturbed
		if c.StateTimer >= 6.0 {
			c.detachStation()
			c.State = CatWandering
			c.pickRandomTarget(roomMinX, roomMaxX, roomMinY, roomMaxY)
		}

	case CatShooedFleeing:
		if c.moveTowards(c.TargetX, c.TargetY, 80.0, dt) {
			c.State = CatResting
			c.RestTimer = 10.0 + rand.Float64()*5.0 // Calm for 10-15s
		}
	}

	return false
}

// Shoo startles the cat and sends it fleeing to a quiet corner.
func (c *Cat) Shoo(rugX, rugY float64) {
	c.detachStation()
	c.State = CatShooedFleeing
	c.TargetX = rugX + 4.0
	c.TargetY = rugY + 4.0
}

func (c *Cat) detachStation() {
	if c.TargetStn != nil {
		c.TargetStn.CatBoost = false
		c.TargetStn = nil
	}
}

func (c *Cat) pickRandomTarget(minX, maxX, minY, maxY float64) {
	c.TargetX = minX + rand.Float64()*(maxX-minX)
	c.TargetY = minY + rand.Float64()*(maxY-minY)
}

func (c *Cat) moveTowards(tx, ty float64, speed float64, dt float64) bool {
	dx := tx - c.X
	dy := ty - c.Y
	dist := math.Hypot(dx, dy)

	if dist < 2.0 {
		c.X = tx
		c.Y = ty
		return true
	}

	c.FacingRight = dx >= 0
	c.X += (dx / dist) * speed * dt
	c.Y += (dy / dist) * speed * dt
	return false
}

// Draw renders the animated orange tabby cat.
func (c *Cat) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	isMoving := c.State == CatWandering || c.State == CatTargetingStation || c.State == CatShooedFleeing
	tailOffset := math.Sin(c.WalkAnim) * 2.0

	// Orange tabby palette
	bodyCol := color.RGBA{240, 130, 40, 255}
	stripeCol := color.RGBA{190, 80, 20, 255}
	eyeCol := color.RGBA{70, 220, 80, 255}
	whiteCol := color.RGBA{255, 245, 235, 255}

	// Body
	c.drawRect(screen, pixelImg, c.X+2, c.Y+3, 8, 5, bodyCol)
	// Stripes
	c.drawRect(screen, pixelImg, c.X+4, c.Y+3, 1, 4, stripeCol)
	c.drawRect(screen, pixelImg, c.X+6, c.Y+3, 1, 4, stripeCol)

	// Head
	hx := c.X + 7
	if !c.FacingRight {
		hx = c.X + 1
	}
	c.drawRect(screen, pixelImg, hx, c.Y+1, 5, 5, bodyCol)
	// Ears
	c.drawRect(screen, pixelImg, hx+1, c.Y-1, 1, 2, stripeCol)
	c.drawRect(screen, pixelImg, hx+3, c.Y-1, 1, 2, stripeCol)
	// Eyes
	eyeX := hx + 3
	if !c.FacingRight {
		eyeX = hx + 1
	}
	c.drawRect(screen, pixelImg, eyeX, c.Y+2, 1, 1, eyeCol)

	// Paws / Legs
	legBob := 0.0
	if isMoving {
		legBob = math.Abs(math.Sin(c.WalkAnim)) * 1.5
	}
	c.drawRect(screen, pixelImg, c.X+3, c.Y+8-legBob, 2, 2, whiteCol)
	c.drawRect(screen, pixelImg, c.X+7, c.Y+8+legBob, 2, 2, whiteCol)

	// Tail
	tailX := c.X
	if !c.FacingRight {
		tailX = c.X + 10
	}
	c.drawRect(screen, pixelImg, tailX, c.Y+2+tailOffset, 2, 2, bodyCol)
}

func (c *Cat) drawRect(screen *ebiten.Image, pixelImg *ebiten.Image, x, y, w, h float64, col color.RGBA) {
	c.drawOpts.GeoM.Reset()
	c.drawOpts.GeoM.Scale(w, h)
	c.drawOpts.GeoM.Translate(x, y)

	af := float32(col.A) / 255.0
	c.drawOpts.ColorScale.Reset()
	c.drawOpts.ColorScale.Scale(
		(float32(col.R)/255.0)*af,
		(float32(col.G)/255.0)*af,
		(float32(col.B)/255.0)*af,
		af,
	)
	screen.DrawImage(pixelImg, &c.drawOpts)
}
