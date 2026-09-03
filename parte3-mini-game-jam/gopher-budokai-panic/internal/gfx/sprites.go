package gfx

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Pose int

const (
	PoseIdle Pose = iota
	PoseMove
	PoseDash
	PoseCharge
	PoseMelee1
	PoseMelee2
	PoseBeamPrep
	PoseBeamFire
	PoseHurt
	PoseKnockback
	PosePanic
	NumPoses
)

const (
	SpriteWidth  = 48
	SpriteHeight = 48
)

type FighterType int

const (
	FighterPlayer FighterType = iota // Super Saiyan (Gold hair, Orange Gi)
	FighterCPU                       // Saiyan Prince (Dark hair, Royal Armor)
	NumFighters
)

type SpriteCache struct {
	images [NumFighters][NumPoses]*ebiten.Image
}

var globalSpriteCache *SpriteCache

func InitSprites() *SpriteCache {
	if globalSpriteCache != nil {
		return globalSpriteCache
	}
	cache := &SpriteCache{}
	for f := 0; f < int(NumFighters); f++ {
		for p := 0; p < int(NumPoses); p++ {
			cache.images[f][p] = renderFighterPose(FighterType(f), Pose(p))
		}
	}
	globalSpriteCache = cache
	return cache
}

func (sc *SpriteCache) GetSprite(f FighterType, p Pose) *ebiten.Image {
	if int(f) >= int(NumFighters) || int(p) >= int(NumPoses) {
		return sc.images[0][0]
	}
	return sc.images[f][p]
}

func renderFighterPose(f FighterType, p Pose) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, SpriteWidth, SpriteHeight))

	// Base color palette
	var (
		hairCol   color.RGBA
		hairHi    color.RGBA
		clothCol  color.RGBA
		clothDark color.RGBA
		skinCol   color.RGBA = color.RGBA{R: 255, G: 210, B: 165, A: 255}
		skinDark  color.RGBA = color.RGBA{R: 220, G: 160, B: 120, A: 255}
		accentCol color.RGBA
		bootCol   color.RGBA
		outline   color.RGBA = color.RGBA{R: 20, G: 15, B: 30, A: 255}
	)

	if f == FighterPlayer {
		// Super Saiyan
		hairCol = color.RGBA{R: 255, G: 225, B: 50, A: 255}
		hairHi = color.RGBA{R: 255, G: 255, B: 180, A: 255}
		clothCol = color.RGBA{R: 240, G: 95, B: 10, A: 255}    // Orange Gi
		clothDark = color.RGBA{R: 190, G: 60, B: 5, A: 255}
		accentCol = color.RGBA{R: 20, G: 70, B: 190, A: 255}   // Blue undershirt & belt
		bootCol = color.RGBA{R: 20, G: 50, B: 140, A: 255}
	} else {
		// Saiyan Prince (Armor)
		hairCol = color.RGBA{R: 35, G: 20, B: 45, A: 255}      // Dark spiked hair
		hairHi = color.RGBA{R: 80, G: 50, B: 110, A: 255}
		clothCol = color.RGBA{R: 225, G: 225, B: 235, A: 255}  // White Armor
		clothDark = color.RGBA{R: 160, G: 165, B: 185, A: 255}
		accentCol = color.RGBA{R: 220, G: 160, B: 20, A: 255}  // Gold armor straps
		bootCol = color.RGBA{R: 230, G: 230, B: 240, A: 255}   // White boots
	}

	setPixel := func(x, y int, c color.RGBA) {
		if x >= 0 && x < SpriteWidth && y >= 0 && y < SpriteHeight {
			// Premultiply alpha
			a := float64(c.A) / 255.0
			pr := uint8(float64(c.R) * a)
			pg := uint8(float64(c.G) * a)
			pb := uint8(float64(c.B) * a)
			img.SetRGBA(x, y, color.RGBA{R: pr, G: pg, B: pb, A: c.A})
		}
	}

	drawRect := func(x0, y0, w, h int, c color.RGBA) {
		for y := y0; y < y0+h; y++ {
			for x := x0; x < x0+w; x++ {
				setPixel(x, y, c)
			}
		}
	}

	drawOutlineRect := func(x0, y0, w, h int, fill, out color.RGBA) {
		for y := y0; y < y0+h; y++ {
			for x := x0; x < x0+w; x++ {
				if x == x0 || x == x0+w-1 || y == y0 || y == y0+h-1 {
					setPixel(x, y, out)
				} else {
					setPixel(x, y, fill)
				}
			}
		}
	}

	// Center reference
	cx := 24
	cy := 24

	// Dynamic offset based on pose
	headY := cy - 8
	torsoY := cy + 1
	legsY := cy + 11

	switch p {
	case PosePanic:
		// Trembling offset
		headY += 1
		torsoY += 2
		legsY += 2
	case PoseDash:
		// Leaning forward
		headY += 2
		torsoY += 2
	case PoseCharge:
		// Crouched power stance
		headY += 2
		torsoY += 3
		legsY += 3
	case PoseKnockback:
		headY -= 2
	}

	// 1. Legs
	switch p {
	case PoseDash:
		// Trailing legs in flight
		drawOutlineRect(cx-10, legsY, 6, 7, clothDark, outline)
		drawOutlineRect(cx-14, legsY+3, 5, 5, bootCol, outline)
	case PoseCharge:
		// Wide power stance
		drawOutlineRect(cx-8, legsY, 6, 9, clothDark, outline)
		drawOutlineRect(cx+3, legsY, 6, 9, clothDark, outline)
		drawOutlineRect(cx-9, legsY+7, 6, 4, bootCol, outline)
		drawOutlineRect(cx+4, legsY+7, 6, 4, bootCol, outline)
	case PoseMelee2:
		// Flying kick!
		drawOutlineRect(cx-5, legsY, 5, 8, clothDark, outline)
		drawOutlineRect(cx+2, legsY-4, 11, 6, bootCol, outline) // Extended kick leg
	default:
		// Hovering floating legs
		drawOutlineRect(cx-6, legsY, 5, 9, clothDark, outline)
		drawOutlineRect(cx+2, legsY, 5, 9, clothDark, outline)
		drawOutlineRect(cx-6, legsY+7, 5, 4, bootCol, outline)
		drawOutlineRect(cx+2, legsY+7, 5, 4, bootCol, outline)
	}

	// 2. Torso
	drawOutlineRect(cx-6, torsoY, 13, 11, clothCol, outline)
	if f == FighterPlayer {
		// Blue chest undershirt V
		drawRect(cx-3, torsoY+1, 7, 4, accentCol)
		// Blue belt sash
		drawRect(cx-5, torsoY+8, 11, 3, accentCol)
	} else {
		// Armor shoulder gold straps and chest panel
		drawRect(cx-5, torsoY+1, 3, 7, accentCol)
		drawRect(cx+3, torsoY+1, 3, 7, accentCol)
		drawRect(cx-2, torsoY+2, 5, 6, clothDark)
	}

	// 3. Arms based on pose
	switch p {
	case PoseCharge:
		// Clenched fists at sides
		drawOutlineRect(cx-11, torsoY+1, 5, 9, skinCol, outline)
		drawOutlineRect(cx+7, torsoY+1, 5, 9, skinCol, outline)
	case PoseMelee1:
		// Punching forward with fist extended
		drawOutlineRect(cx-5, torsoY+2, 5, 6, skinCol, outline)
		drawOutlineRect(cx+4, torsoY+1, 10, 5, skinDark, outline) // Right straight punch
	case PoseBeamPrep:
		// Cupped hands at hip gathering ki
		drawOutlineRect(cx-10, torsoY+3, 6, 6, skinCol, outline)
		drawOutlineRect(cx-8, torsoY+4, 7, 7, color.RGBA{R: 200, G: 240, B: 255, A: 255}, outline)
	case PoseBeamFire:
		// Double palms thrust straight forward releasing the beam!
		drawOutlineRect(cx+4, torsoY+1, 12, 6, skinCol, outline)
		drawOutlineRect(cx+4, torsoY+3, 12, 6, skinDark, outline)
	case PosePanic:
		// Frantic flailing hands near head
		drawOutlineRect(cx-9, torsoY-2, 5, 7, skinCol, outline)
		drawOutlineRect(cx+5, torsoY-2, 5, 7, skinCol, outline)
	default:
		// Ready guard stance
		drawOutlineRect(cx-9, torsoY+1, 4, 8, skinCol, outline)
		drawOutlineRect(cx+6, torsoY+1, 4, 8, skinCol, outline)
	}

	// 4. Head & Face
	drawOutlineRect(cx-5, headY, 11, 10, skinCol, outline)
	// Eyes & Eyebrows
	if p == PoseHurt || p == PoseKnockback {
		// Shut tight pain eyes
		setPixel(cx-2, headY+4, outline)
		setPixel(cx-1, headY+3, outline)
		setPixel(cx+2, headY+4, outline)
		setPixel(cx+3, headY+3, outline)
	} else if p == PosePanic {
		// Wide frantic pupil dots
		setPixel(cx-3, headY+3, outline)
		setPixel(cx-2, headY+3, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		setPixel(cx-2, headY+4, outline)
		setPixel(cx+2, headY+3, outline)
		setPixel(cx+3, headY+3, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		setPixel(cx+3, headY+4, outline)
		// Open gasping mouth
		drawRect(cx-1, headY+7, 3, 2, outline)
	} else {
		// Sharp intense DBZ warrior eyes
		eyeCol := color.RGBA{R: 30, G: 220, B: 230, A: 255} // Cyan for SSJ
		if f == FighterCPU {
			eyeCol = outline
		}
		setPixel(cx-3, headY+4, eyeCol)
		setPixel(cx-2, headY+4, outline)
		setPixel(cx+2, headY+4, eyeCol)
		setPixel(cx+3, headY+4, outline)
		// Smirk / grit mouth
		setPixel(cx, headY+7, outline)
		setPixel(cx+1, headY+7, outline)
	}

	// 5. Spiked Anime Hair
	if f == FighterPlayer {
		// Golden Super Saiyan Spikes pointing up and back
		spikes := []struct{ x, y, w, h int }{
			{cx - 7, headY - 8, 5, 9},   // Left flare
			{cx - 3, headY - 11, 6, 12}, // Center towering spike
			{cx + 2, headY - 10, 5, 11}, // Right upper spike
			{cx + 6, headY - 6, 4, 8},   // Right side flare
			{cx - 5, headY - 2, 3, 4},   // Front bang
		}
		for _, sp := range spikes {
			drawOutlineRect(sp.x, sp.y, sp.w, sp.h, hairCol, outline)
			drawRect(sp.x+1, sp.y+1, sp.w-2, sp.h/2, hairHi)
		}
	} else {
		// Flame-style dark spiked hair (Vegeta)
		spikes := []struct{ x, y, w, h int }{
			{cx - 6, headY - 9, 4, 10},
			{cx - 2, headY - 12, 5, 13},
			{cx + 2, headY - 10, 5, 11},
			{cx + 5, headY - 7, 4, 8},
		}
		for _, sp := range spikes {
			drawOutlineRect(sp.x, sp.y, sp.w, sp.h, hairCol, outline)
			drawRect(sp.x+1, sp.y+1, sp.w-2, sp.h/3, hairHi)
		}
	}

	// Panic sweat drops!
	if p == PosePanic {
		sweatCol := color.RGBA{R: 100, G: 220, B: 255, A: 255}
		setPixel(cx-7, headY+1, sweatCol)
		setPixel(cx-7, headY+2, sweatCol)
		setPixel(cx+7, headY+1, sweatCol)
		setPixel(cx+7, headY+2, sweatCol)
	}

	return ebiten.NewImageFromImage(img)
}

// DrawFighter draws the fighter sprite with correct matrix transformation order:
// Pivot -> Scale -> Rotate -> World Translation
func DrawFighter(screen *ebiten.Image, sprite *ebiten.Image, x, y, scaleX, scaleY, angle float64, facingLeft bool) {
	if sprite == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}

	pivotX := float64(SpriteWidth) / 2.0
	pivotY := float64(SpriteHeight) / 2.0

	// 1. Pivot translation to center
	op.GeoM.Translate(-pivotX, -pivotY)

	// 2. Facing direction flip & scale
	sx := scaleX
	if facingLeft {
		sx = -sx
	}
	op.GeoM.Scale(sx, scaleY)

	// 3. Rotation
	if angle != 0 {
		op.GeoM.Rotate(angle)
	}

	// 4. World translation
	op.GeoM.Translate(x, y)

	screen.DrawImage(sprite, op)
}

// DrawEnergyAura draws an animated pulsing energy aura around the fighter.
func DrawEnergyAura(screen *ebiten.Image, x, y float64, radius float64, col color.RGBA, isSparking bool) {
	// Pulsing translucent rings
	ringAlpha := float32(0.35)
	if isSparking {
		ringAlpha = 0.65
	}
	c := color.RGBA{
		R: col.R,
		G: col.G,
		B: col.B,
		A: uint8(float32(col.A) * ringAlpha),
	}
	_ = c
	_ = radius
	// Aura drawing helper
}
