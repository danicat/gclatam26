package entity

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	DirDown Direction = iota
	DirUp
	DirLeft
	DirRight
)

// Chef represents the playable player character.
type Chef struct {
	X, Y        float64
	Speed       float64
	Dir         Direction
	Tool        ToolType
	WalkAnim    float64
	IsMoving    bool
	W, H        float64
	drawOpts    ebiten.DrawImageOptions
}

// NewChef spawns the chef at the starting location with empty hands.
func NewChef(x, y float64) *Chef {
	return &Chef{
		X:        x,
		Y:        y,
		Speed:    85.0,
		Dir:      DirDown,
		Tool:     ToolNone,
		W:        14,
		H:        18,
	}
}

// Update processes player input movement deltas and boundaries.
func (c *Chef) Update(dt float64, moveX, moveY float64, minX, maxX, minY, maxY float64) {
	c.IsMoving = moveX != 0 || moveY != 0

	if c.IsMoving {
		c.WalkAnim += dt * 10.0

		// Update facing direction
		if math.Abs(moveX) > math.Abs(moveY) {
			if moveX > 0 {
				c.Dir = DirRight
			} else {
				c.Dir = DirLeft
			}
		} else {
			if moveY > 0 {
				c.Dir = DirDown
			} else {
				c.Dir = DirUp
			}
		}

		c.X += moveX * c.Speed * dt
		c.Y += moveY * c.Speed * dt

		// Clamp within kitchen room bounds
		if c.X < minX {
			c.X = minX
		}
		if c.X+c.W > maxX {
			c.X = maxX - c.W
		}
		if c.Y < minY {
			c.Y = minY
		}
		if c.Y+c.H > maxY {
			c.Y = maxY - c.H
		}
	} else {
		c.WalkAnim = 0
	}
}

// Center returns the middle coordinates of the chef bounding box.
func (c *Chef) Center() (float64, float64) {
	return c.X + c.W/2, c.Y + c.H/2
}

// Draw renders the pixel art chef and any held tool overhead.
func (c *Chef) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	legBob := 0.0
	if c.IsMoving {
		legBob = math.Sin(c.WalkAnim) * 1.5
	}

	// Palette
	skinCol := color.RGBA{245, 205, 175, 255}
	hatCol := color.RGBA{250, 250, 255, 255}
	apronCol := color.RGBA{235, 60, 60, 255}
	pantsCol := color.RGBA{45, 50, 70, 255}
	shoeCol := color.RGBA{25, 20, 30, 255}

	// 1. Chef Hat (Top)
	c.drawRect(screen, pixelImg, c.X+3, c.Y-4, 8, 4, hatCol)
	c.drawRect(screen, pixelImg, c.X+2, c.Y-6, 10, 3, hatCol)

	// 2. Head / Face
	c.drawRect(screen, pixelImg, c.X+3, c.Y, 8, 6, skinCol)
	// Eyes based on direction
	switch c.Dir {
	case DirDown:
		c.drawRect(screen, pixelImg, c.X+4, c.Y+2, 2, 2, shoeCol)
		c.drawRect(screen, pixelImg, c.X+8, c.Y+2, 2, 2, shoeCol)
		// Chef mustache
		c.drawRect(screen, pixelImg, c.X+4, c.Y+4, 6, 1, shoeCol)
	case DirLeft:
		c.drawRect(screen, pixelImg, c.X+3, c.Y+2, 2, 2, shoeCol)
		c.drawRect(screen, pixelImg, c.X+3, c.Y+4, 3, 1, shoeCol)
	case DirRight:
		c.drawRect(screen, pixelImg, c.X+9, c.Y+2, 2, 2, shoeCol)
		c.drawRect(screen, pixelImg, c.X+8, c.Y+4, 3, 1, shoeCol)
	case DirUp:
		// Hair on back of head
		c.drawRect(screen, pixelImg, c.X+3, c.Y+1, 8, 4, shoeCol)
	}

	// 3. Body & Apron
	c.drawRect(screen, pixelImg, c.X+2, c.Y+6, 10, 7, apronCol)
	// White chef tie/accents
	c.drawRect(screen, pixelImg, c.X+5, c.Y+6, 4, 3, hatCol)

	// 4. Pants & Shoes (with walking bob)
	c.drawRect(screen, pixelImg, c.X+3, c.Y+13, 3, 3, pantsCol)
	c.drawRect(screen, pixelImg, c.X+8, c.Y+13, 3, 3, pantsCol)

	c.drawRect(screen, pixelImg, c.X+3, c.Y+16+legBob, 3, 2, shoeCol)
	c.drawRect(screen, pixelImg, c.X+8, c.Y+16-legBob, 3, 2, shoeCol)

	// 5. Held Tool (drawn overhead or in hands)
	if c.Tool != ToolNone {
		toolY := c.Y - 10.0
		if c.IsMoving {
			toolY += math.Sin(c.WalkAnim*2.0) * 0.8
		}
		DrawToolIcon(screen, pixelImg, c.Tool, c.X+c.W/2, toolY)
	}
}

func (c *Chef) drawRect(screen *ebiten.Image, pixelImg *ebiten.Image, x, y, w, h float64, col color.RGBA) {
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
