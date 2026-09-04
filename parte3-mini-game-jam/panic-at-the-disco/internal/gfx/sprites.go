package gfx

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawPlayer draws the procedural 70s disco dancer with afro, glasses, bell-bottoms, and flailing panic arms.
func DrawPlayer(screen *ebiten.Image, x, y float64, facingX float64, animTime float64, isMoving bool, isDashing bool, panicLevel float64) {
	// Base bobbing
	bob := 0.0
	legSwing := 0.0
	armFlail := 0.0
	if isMoving {
		bob = math.Sin(animTime*14.0) * 2.0
		legSwing = math.Sin(animTime * 14.0)
	}

	// High panic causes rapid jitter and wild arm flailing
	panicJitterX := 0.0
	panicJitterY := 0.0
	if panicLevel > 50.0 {
		jitterFactor := (panicLevel - 50.0) / 50.0
		panicJitterX = math.Sin(animTime*30.0) * 1.5 * jitterFactor
		panicJitterY = math.Cos(animTime*33.0) * 1.5 * jitterFactor
		armFlail = math.Sin(animTime*22.0) * 6.0 * jitterFactor
	}

	cx := float32(x + panicJitterX)
	cy := float32(y + bob + panicJitterY)

	// Shadow under feet
	vector.DrawFilledCircle(screen, cx, float32(y+16), 11, color.RGBA{0, 0, 0, 100}, false)

	// Direction flip factor
	dir := float32(1.0)
	if facingX < 0 {
		dir = -1.0
	}

	// --- 1. Bell-bottom pants (Legs) ---
	// Pants color: Bright 70s White / Silver or Groovy Pink
	pantsColor := color.RGBA{R: 240, G: 240, B: 250, A: 255}
	if isDashing {
		pantsColor = color.RGBA{R: 0, G: 255, B: 255, A: 255} // Glowing neon cyan during dash
	}
	shoeColor := color.RGBA{R: 40, G: 20, B: 10, A: 255}

	legOffset := float32(legSwing * 5.0)

	// Left Leg (Bell-bottom flare)
	lx1 := cx - 5*dir - legOffset
	ly1 := cy + 6
	vector.DrawFilledRect(screen, lx1-3, ly1, 6, 9, pantsColor, false)
	vector.DrawFilledRect(screen, lx1-5, ly1+8, 9, 4, pantsColor, false) // Bell flare
	vector.DrawFilledRect(screen, lx1-6, ly1+12, 10, 3, shoeColor, false) // Platform shoe

	// Right Leg
	rx1 := cx + 5*dir + legOffset
	ry1 := cy + 6
	vector.DrawFilledRect(screen, rx1-3, ry1, 6, 9, pantsColor, false)
	vector.DrawFilledRect(screen, rx1-4, ry1+8, 9, 4, pantsColor, false) // Bell flare
	vector.DrawFilledRect(screen, rx1-4, ry1+12, 10, 3, shoeColor, false) // Platform shoe

	// --- 2. Torso (70s open-collar disco shirt) ---
	shirtColor := color.RGBA{R: 255, G: 20, B: 147, A: 255} // Deep Pink / Magenta
	if isDashing {
		shirtColor = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	}
	vector.DrawFilledRect(screen, cx-7, cy-5, 14, 12, shirtColor, false)

	// Gold medallion necklace
	vector.DrawFilledCircle(screen, cx, cy-1, 2.5, color.RGBA{255, 215, 0, 255}, false)

	// --- 3. Arms (Flailing in panic or disco posture) ---
	skinColor := color.RGBA{R: 240, G: 185, B: 140, A: 255}
	armAngle := float32(armFlail)

	// Left arm
	vector.DrawFilledRect(screen, cx-12, cy-6-armAngle, 5, 8, shirtColor, false)
	vector.DrawFilledCircle(screen, cx-10, cy+3-armAngle, 3, skinColor, false)

	// Right arm (Pointing up/flailing)
	vector.DrawFilledRect(screen, cx+7, cy-6+armAngle, 5, 8, shirtColor, false)
	vector.DrawFilledCircle(screen, cx+10, cy+3+armAngle, 3, skinColor, false)

	// --- 4. Head & Face ---
	vector.DrawFilledCircle(screen, cx, cy-11, 7, skinColor, false)

	// Funky oversized sunglasses
	glassesColor := color.RGBA{R: 20, G: 20, B: 30, A: 255}
	lensColor := color.RGBA{R: 255, G: 165, B: 0, A: 220} // Amber tint
	vector.DrawFilledRect(screen, cx-5+2*dir, cy-13, 10, 4, glassesColor, false)
	vector.DrawFilledRect(screen, cx-4+2*dir, cy-12, 3, 2, lensColor, false)
	vector.DrawFilledRect(screen, cx+1+2*dir, cy-12, 3, 2, lensColor, false)

	// Panic expression (wide open mouth when panic > 40%)
	if panicLevel > 40.0 {
		mouthW := float32(3.0 + (panicLevel-40)/20.0)
		if mouthW > 6.0 {
			mouthW = 6.0
		}
		vector.DrawFilledRect(screen, cx-mouthW/2+1*dir, cy-8, mouthW, 3, color.RGBA{180, 20, 20, 255}, false)
	}

	// --- 5. Magnificent Afro Hairstyle ---
	afroColor := color.RGBA{R: 45, G: 25, B: 15, A: 255}
	vector.DrawFilledCircle(screen, cx, cy-16, 11, afroColor, false)
	vector.DrawFilledCircle(screen, cx-6, cy-13, 8, afroColor, false)
	vector.DrawFilledCircle(screen, cx+6, cy-13, 8, afroColor, false)
}

// DrawDiscoBall draws the giant falling mirror ball with specular facet grid.
func DrawDiscoBall(screen *ebiten.Image, x, y, radius float64, rotation float64) {
	cx := float32(x)
	cy := float32(y)
	r := float32(radius)

	// Base sphere gradient / dark edge
	vector.DrawFilledCircle(screen, cx, cy, r, color.RGBA{120, 140, 160, 255}, false)

	// Inner reflective body
	vector.DrawFilledCircle(screen, cx, cy, r*0.92, color.RGBA{210, 225, 240, 255}, false)

	// Mirror facet grid
	gridCols := 6
	gridRows := 6
	stepX := (r * 1.8) / float32(gridCols)
	stepY := (r * 1.8) / float32(gridRows)

	for i := 0; i < gridCols; i++ {
		for j := 0; j < gridRows; j++ {
			fx := (cx - r*0.9) + float32(i)*stepX + float32(math.Sin(rotation+float64(j)*0.5))*2.0
			fy := (cy - r*0.9) + float32(j)*stepY

			// Check if inside circle
			distSq := (fx-cx)*(fx-cx) + (fy-cy)*(fy-cy)
			if distSq < (r*0.85)*(r*0.85) {
				// Shimmering mirror tile facet
				shimmer := math.Sin(rotation*4.0 + float64(i*3+j*5))
				var facetCol color.RGBA
				if shimmer > 0.4 {
					facetCol = color.RGBA{255, 255, 255, 255} // Bright sparkle reflection
				} else if shimmer < -0.4 {
					facetCol = color.RGBA{130, 150, 180, 255} // Cool silver shadow
				} else {
					facetCol = color.RGBA{190, 210, 230, 255}
				}
				vector.DrawFilledRect(screen, fx, fy, stepX-1.2, stepY-1.2, facetCol, false)
			}
		}
	}

	// Bright highlight reflection dot at top-left
	vector.DrawFilledCircle(screen, cx-r*0.35, cy-r*0.35, r*0.2, color.RGBA{255, 255, 255, 230}, false)
}

// DrawTelegraphCircle draws the expanding warning shadow on the dance floor.
func DrawTelegraphCircle(screen *ebiten.Image, x, y, radius, progress float64) {
	cx := float32(x)
	cy := float32(y)
	r := float32(radius)

	// Outer faint danger zone
	vector.DrawFilledCircle(screen, cx, cy, r, color.RGBA{255, 50, 50, 50}, false)
	vector.StrokeCircle(screen, cx, cy, r, 2.0, color.RGBA{255, 80, 80, 180}, false)

	// Inner expanding indicator
	innerR := r * float32(progress)
	if innerR > r {
		innerR = r
	}
	vector.DrawFilledCircle(screen, cx, cy, innerR, color.RGBA{255, 20, 20, uint8(80 + progress*140)}, false)
	vector.StrokeCircle(screen, cx, cy, innerR, 2.5, color.RGBA{255, 240, 0, 255}, false)
}

// DrawTruss draws a fallen or falling steel lighting truss with danger stripes.
func DrawTruss(screen *ebiten.Image, x, y, w, h float64) {
	fx := float32(x)
	fy := float32(y)
	fw := float32(w)
	fh := float32(h)

	// Steel truss background
	vector.DrawFilledRect(screen, fx, fy, fw, fh, color.RGBA{50, 55, 65, 255}, false)
	vector.StrokeRect(screen, fx, fy, fw, fh, 2.0, color.RGBA{160, 170, 190, 255}, false)

	// Hazard warning stripes (Yellow & Black diagonal)
	stripeW := float32(10.0)
	for s := fx; s < fx+fw; s += stripeW * 2 {
		vector.DrawFilledRect(screen, s, fy+2, stripeW, fh-4, color.RGBA{255, 215, 0, 230}, false)
	}

	// Sparks anchor points
	vector.DrawFilledCircle(screen, fx+4, fy+fh/2, 3, color.RGBA{255, 100, 0, 255}, false)
	vector.DrawFilledCircle(screen, fx+fw-4, fy+fh/2, 3, color.RGBA{255, 100, 0, 255}, false)
}

// DrawDrinkPuddle draws a slippery spilled cocktail on the dance floor.
func DrawDrinkPuddle(screen *ebiten.Image, x, y, radius float64, col color.RGBA) {
	cx := float32(x)
	cy := float32(y)
	r := float32(radius)

	// Main liquid slick (irregular ellipse)
	vector.DrawFilledCircle(screen, cx, cy, r, col, false)
	vector.DrawFilledCircle(screen, cx+r*0.4, cy-r*0.2, r*0.6, col, false)
	vector.DrawFilledCircle(screen, cx-r*0.3, cy+r*0.3, r*0.5, col, false)

	// Dropped broken glass or cocktail umbrella
	vector.DrawFilledCircle(screen, cx+2, cy-3, 2, color.RGBA{255, 255, 255, 220}, false)
	vector.StrokeLine(screen, cx-4, cy+2, cx-1, cy+7, 1.5, color.RGBA{255, 50, 50, 255}, false)
}

// DrawPanickedClubber draws a wandering retro nightclub patron in panic.
func DrawPanickedClubber(screen *ebiten.Image, x, y float64, animTime float64, outfitColor color.RGBA) {
	cx := float32(x)
	cy := float32(y)

	bob := float32(math.Sin(animTime*10.0) * 2.0)
	flail := float32(math.Sin(animTime*16.0) * 4.0)

	// Shadow
	vector.DrawFilledCircle(screen, cx, cy+14, 9, color.RGBA{0, 0, 0, 90}, false)

	// Pants / Skirt
	vector.DrawFilledRect(screen, cx-6, cy+5+bob, 12, 9, color.RGBA{30, 30, 50, 255}, false)

	// Shirt / Outfit
	vector.DrawFilledRect(screen, cx-7, cy-5+bob, 14, 11, outfitColor, false)

	// Arms waving up in the air
	skinColor := color.RGBA{235, 190, 150, 255}
	vector.DrawFilledRect(screen, cx-11, cy-12+bob-flail, 4, 10, skinColor, false)
	vector.DrawFilledRect(screen, cx+7, cy-12+bob+flail, 4, 10, skinColor, false)

	// Head
	vector.DrawFilledCircle(screen, cx, cy-10+bob, 6, skinColor, false)

	// Hair
	vector.DrawFilledCircle(screen, cx, cy-14+bob, 7, color.RGBA{80, 40, 20, 255}, false)
}

// DrawExitPortal draws the illuminated exit door with glowing green neon sign.
func DrawExitPortal(screen *ebiten.Image, x, y, w, h float64, animTime float64) {
	fx := float32(x)
	fy := float32(y)
	fw := float32(w)
	fh := float32(h)

	// Door opening (dark gateway to safety)
	vector.DrawFilledRect(screen, fx, fy, fw, fh, color.RGBA{10, 25, 20, 255}, false)

	// Glowing green border
	pulse := float32(0.8 + 0.2*math.Sin(animTime*6.0))
	borderCol := color.RGBA{
		R: uint8(30 * pulse),
		G: uint8(255 * pulse),
		B: uint8(100 * pulse),
		A: 255,
	}
	vector.StrokeRect(screen, fx, fy, fw, fh, 3.0, borderCol, false)

	// Overhead EXIT sign box
	signY := fy - 16
	vector.DrawFilledRect(screen, fx+fw/2-24, signY, 48, 14, color.RGBA{20, 60, 30, 255}, false)
	vector.StrokeRect(screen, fx+fw/2-24, signY, 48, 14, 1.5, borderCol, false)

	// EXIT letters block-drawn
	// E
	vector.DrawFilledRect(screen, fx+fw/2-18, signY+3, 6, 8, borderCol, false)
	vector.DrawFilledRect(screen, fx+fw/2-16, signY+5, 4, 1, color.RGBA{20, 60, 30, 255}, false)
	vector.DrawFilledRect(screen, fx+fw/2-16, signY+8, 4, 1, color.RGBA{20, 60, 30, 255}, false)
	// X
	vector.StrokeLine(screen, fx+fw/2-8, signY+3, fx+fw/2-4, signY+11, 2.0, borderCol, false)
	vector.StrokeLine(screen, fx+fw/2-4, signY+3, fx+fw/2-8, signY+11, 2.0, borderCol, false)
	// I
	vector.DrawFilledRect(screen, fx+fw/2, signY+3, 2, 8, borderCol, false)
	// T
	vector.DrawFilledRect(screen, fx+fw/2+5, signY+3, 7, 2, borderCol, false)
	vector.DrawFilledRect(screen, fx+fw/2+7.5, signY+3, 2, 8, borderCol, false)

	// Animated chevron arrow pointing toward the exit
	arrowPulse := math.Mod(animTime*30.0, float64(fh))
	vector.StrokeLine(screen, fx+fw/2-10, fy+float32(arrowPulse), fx+fw/2, fy+float32(arrowPulse)-8, 2.5, borderCol, false)
	vector.StrokeLine(screen, fx+fw/2, fy+float32(arrowPulse)-8, fx+fw/2+10, fy+float32(arrowPulse), 2.5, borderCol, false)
}
