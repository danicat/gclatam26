package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const TileSize = 24

var (
	imgFloor        *ebiten.Image
	imgWall         *ebiten.Image
	imgHole         *ebiten.Image
	imgHoleFilled   *ebiten.Image
	imgBoulder      *ebiten.Image
	imgClock        *ebiten.Image
	imgArtifact     *ebiten.Image
	imgGopher       *ebiten.Image
	imgHUDClockFace *ebiten.Image
	imgHUDClockHand *ebiten.Image
)

func initAssets() {
	initWhitePixel()
	imgFloor = createFloorImage()
	imgWall = createWallImage()
	imgHole = createHoleImage()
	imgHoleFilled = createHoleFilledImage()
	imgBoulder = createBoulderImage()
	imgClock = createClockImage()
	imgArtifact = createArtifactImage()
	imgGopher = createGopherImage()
	imgHUDClockFace = createHUDClockFaceImage()
	imgHUDClockHand = createHUDClockHandImage()
}

func setPixel(img *ebiten.Image, x, y int, c color.Color) {
	if x >= 0 && x < TileSize && y >= 0 && y < TileSize {
		img.Set(x, y, c)
	}
}

func fillRect(img *ebiten.Image, x, y, w, h int, c color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			setPixel(img, x+dx, y+dy, c)
		}
	}
}

func createFloorImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	base := color.RGBA{R: 28, G: 30, B: 42, A: 255}
	border := color.RGBA{R: 20, G: 22, B: 32, A: 255}
	crack := color.RGBA{R: 38, G: 40, B: 56, A: 255}

	fillRect(img, 0, 0, TileSize, TileSize, base)
	// Outer border
	for i := 0; i < TileSize; i++ {
		setPixel(img, i, 0, border)
		setPixel(img, 0, i, border)
		setPixel(img, i, TileSize-1, border)
		setPixel(img, TileSize-1, i, border)
	}
	// Subtle floor grain / stone tiles
	setPixel(img, 6, 8, crack)
	setPixel(img, 7, 8, crack)
	setPixel(img, 8, 9, crack)
	setPixel(img, 18, 15, crack)
	setPixel(img, 19, 15, crack)
	return img
}

func createWallImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	wallDark := color.RGBA{R: 44, G: 46, B: 62, A: 255}
	wallLight := color.RGBA{R: 64, G: 68, B: 90, A: 255}
	mortar := color.RGBA{R: 18, G: 19, B: 28, A: 255}
	runeColor := color.RGBA{R: 120, G: 70, B: 180, A: 255}

	fillRect(img, 0, 0, TileSize, TileSize, wallDark)

	// Top highlight
	fillRect(img, 1, 1, TileSize-2, 2, wallLight)

	// Brick mortar lines
	for x := 0; x < TileSize; x++ {
		setPixel(img, x, 7, mortar)
		setPixel(img, x, 15, mortar)
		setPixel(img, x, TileSize-1, mortar)
	}
	for y := 0; y < 8; y++ {
		setPixel(img, 11, y, mortar)
	}
	for y := 8; y < 16; y++ {
		setPixel(img, 5, y, mortar)
		setPixel(img, 17, y, mortar)
	}
	for y := 16; y < TileSize; y++ {
		setPixel(img, 11, y, mortar)
	}

	// Tiny glowing rune carving on wall
	setPixel(img, 11, 11, runeColor)
	setPixel(img, 11, 12, runeColor)
	setPixel(img, 10, 11, runeColor)
	setPixel(img, 12, 11, runeColor)

	return img
}

func createHoleImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	voidEdge := color.RGBA{R: 12, G: 10, B: 22, A: 255}
	voidDeep := color.RGBA{R: 4, G: 3, B: 8, A: 255}
	mist := color.RGBA{R: 40, G: 15, B: 60, A: 255}

	// Void hole with jagged edges
	fillRect(img, 2, 2, TileSize-4, TileSize-4, voidDeep)
	// Outer dark rim
	for i := 2; i < TileSize-2; i++ {
		setPixel(img, i, 2, voidEdge)
		setPixel(img, 2, i, voidEdge)
		setPixel(img, i, TileSize-3, voidEdge)
		setPixel(img, TileSize-3, i, voidEdge)
	}
	// Cosmic mist in pit
	setPixel(img, 10, 10, mist)
	setPixel(img, 11, 10, mist)
	setPixel(img, 12, 11, mist)
	setPixel(img, 14, 15, mist)
	setPixel(img, 8, 14, mist)
	return img
}

func createHoleFilledImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	// Dark hole rim
	rim := color.RGBA{R: 15, G: 12, B: 24, A: 255}
	boulderSurface := color.RGBA{R: 90, G: 92, B: 110, A: 255}
	cracks := color.RGBA{R: 50, G: 52, B: 68, A: 255}

	fillRect(img, 1, 1, TileSize-2, TileSize-2, rim)
	fillRect(img, 3, 3, TileSize-6, TileSize-6, boulderSurface)

	// Cracks showing boulder wedged into void
	setPixel(img, 7, 7, cracks)
	setPixel(img, 8, 8, cracks)
	setPixel(img, 9, 8, cracks)
	setPixel(img, 10, 9, cracks)
	setPixel(img, 14, 12, cracks)
	setPixel(img, 15, 13, cracks)
	return img
}

func createBoulderImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	boulderBase := color.RGBA{R: 115, G: 118, B: 135, A: 255}
	boulderLight := color.RGBA{R: 150, G: 154, B: 175, A: 255}
	boulderDark := color.RGBA{R: 70, G: 72, B: 88, A: 255}
	runeGlow := color.RGBA{R: 80, G: 220, B: 200, A: 255}

	// Rounded stone body
	fillRect(img, 4, 3, 16, 18, boulderBase)
	fillRect(img, 3, 5, 18, 14, boulderBase)

	// Top-left highlight
	fillRect(img, 5, 4, 10, 3, boulderLight)
	fillRect(img, 4, 6, 3, 8, boulderLight)

	// Bottom-right shadow
	fillRect(img, 5, 19, 14, 2, boulderDark)
	fillRect(img, 19, 6, 2, 13, boulderDark)

	// Eldritch rune carved in stone
	setPixel(img, 11, 10, runeGlow)
	setPixel(img, 12, 10, runeGlow)
	setPixel(img, 11, 11, runeGlow)
	setPixel(img, 12, 11, runeGlow)
	setPixel(img, 11, 12, runeGlow)
	setPixel(img, 12, 13, runeGlow)
	setPixel(img, 10, 11, runeGlow)
	setPixel(img, 13, 11, runeGlow)
	return img
}

func createClockImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	gold := color.RGBA{R: 245, G: 195, B: 45, A: 255}
	goldLight := color.RGBA{R: 255, G: 235, B: 120, A: 255}
	face := color.RGBA{R: 250, G: 248, B: 240, A: 255}
	hand := color.RGBA{R: 40, G: 30, B: 20, A: 255}
	sparkle := color.RGBA{R: 100, G: 255, B: 220, A: 255}

	// Clock casing (circle-ish)
	fillRect(img, 5, 5, 14, 14, gold)
	fillRect(img, 6, 4, 12, 16, gold)
	fillRect(img, 4, 6, 16, 12, gold)

	// Inner white clock face
	fillRect(img, 7, 7, 10, 10, face)

	// Clock hands (pointing to 2 o'clock)
	setPixel(img, 11, 11, hand)
	setPixel(img, 11, 10, hand)
	setPixel(img, 11, 9, hand)
	setPixel(img, 12, 11, hand)
	setPixel(img, 13, 11, hand)

	// Ring top / loop
	setPixel(img, 11, 2, goldLight)
	setPixel(img, 12, 2, goldLight)
	setPixel(img, 10, 3, goldLight)
	setPixel(img, 13, 3, goldLight)

	// Calming sparkle particles
	setPixel(img, 2, 4, sparkle)
	setPixel(img, 21, 5, sparkle)
	setPixel(img, 20, 18, sparkle)
	return img
}

func createArtifactImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	obsidian := color.RGBA{R: 35, G: 20, B: 55, A: 255}
	goldTrim := color.RGBA{R: 220, G: 175, B: 40, A: 255}
	eldritchGlow := color.RGBA{R: 190, G: 60, B: 240, A: 255}
	core := color.RGBA{R: 255, G: 120, B: 255, A: 255}

	// Pedestal / Base
	fillRect(img, 5, 18, 14, 3, goldTrim)
	fillRect(img, 7, 16, 10, 2, obsidian)

	// Idol body (cosmic tentacled relic)
	fillRect(img, 8, 8, 8, 8, obsidian)
	fillRect(img, 9, 6, 6, 2, obsidian)
	fillRect(img, 10, 4, 4, 2, goldTrim)

	// Outer tentacle curls
	setPixel(img, 6, 10, obsidian)
	setPixel(img, 5, 9, obsidian)
	setPixel(img, 5, 8, goldTrim)
	setPixel(img, 17, 10, obsidian)
	setPixel(img, 18, 9, obsidian)
	setPixel(img, 18, 8, goldTrim)

	// Glowing pulsating eye / core
	fillRect(img, 10, 9, 4, 4, eldritchGlow)
	setPixel(img, 11, 10, core)
	setPixel(img, 12, 10, core)
	setPixel(img, 11, 11, core)
	setPixel(img, 12, 11, core)

	return img
}

func createGopherImage() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)

	gopherBlue := color.RGBA{R: 60, G: 175, B: 220, A: 255}
	gopherLight := color.RGBA{R: 110, G: 210, B: 245, A: 255}
	gopherDark := color.RGBA{R: 35, G: 130, B: 175, A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 20, G: 20, B: 25, A: 255}
	snout := color.RGBA{R: 245, G: 230, B: 205, A: 255}
	capeRed := color.RGBA{R: 220, G: 45, B: 45, A: 255}
	capeDark := color.RGBA{R: 160, G: 25, B: 25, A: 255}

	// Gopher Body (Pudgy, rounded)
	fillRect(img, 6, 7, 12, 13, gopherBlue)
	fillRect(img, 5, 9, 14, 9, gopherBlue)

	// Ears
	fillRect(img, 5, 4, 3, 4, gopherBlue)
	fillRect(img, 16, 4, 3, 4, gopherBlue)
	setPixel(img, 6, 5, snout)
	setPixel(img, 17, 5, snout)

	// Head highlight
	fillRect(img, 8, 6, 8, 2, gopherLight)

	// Red Explorer Scarf / Cape (around neck)
	fillRect(img, 6, 14, 12, 3, capeRed)
	setPixel(img, 5, 15, capeDark)
	setPixel(img, 5, 16, capeRed)
	setPixel(img, 4, 17, capeRed) // fluttering cape tail

	// Belly / Snout
	fillRect(img, 9, 10, 6, 4, snout)
	setPixel(img, 11, 10, black) // tiny nose
	setPixel(img, 12, 10, black)

	// Buck teeth!
	setPixel(img, 11, 12, white)
	setPixel(img, 12, 12, white)

	// Big Cute Expressive Eyes
	fillRect(img, 7, 7, 4, 4, white)
	fillRect(img, 13, 7, 4, 4, white)
	// Black Pupils looking forward/slightly wide
	setPixel(img, 8, 8, black)
	setPixel(img, 9, 8, black)
	setPixel(img, 8, 9, black)
	setPixel(img, 9, 9, black)
	setPixel(img, 14, 8, black)
	setPixel(img, 15, 8, black)
	setPixel(img, 14, 9, black)
	setPixel(img, 15, 9, black)

	// Little paws / feet
	fillRect(img, 6, 20, 3, 2, snout)
	fillRect(img, 15, 20, 3, 2, snout)

	// Shading on side
	for y := 9; y < 19; y++ {
		setPixel(img, 18, y, gopherDark)
	}

	return img
}

func createHUDClockFaceImage() *ebiten.Image {
	img := ebiten.NewImage(24, 24)
	goldDark := color.RGBA{R: 150, G: 105, B: 20, A: 255}
	gold := color.RGBA{R: 220, G: 175, B: 40, A: 255}
	goldLight := color.RGBA{R: 255, G: 225, B: 110, A: 255}
	darkDial := color.RGBA{R: 14, G: 16, B: 25, A: 255}
	tickCol := color.RGBA{R: 180, G: 190, B: 210, A: 255}
	tickRed := color.RGBA{R: 235, G: 45, B: 45, A: 255}

	// Pocket-watch body circle (radius ~8.5 centered at 12.0, 13.0)
	for y := 3; y < 23; y++ {
		for x := 2; x < 22; x++ {
			dx := float64(x) - 11.5
			dy := float64(y) - 12.5
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 9.2 && dist >= 7.0 {
				if dy < -3 || dx < -3 {
					setPixel(img, x, y, goldLight)
				} else if dy > 3 || dx > 3 {
					setPixel(img, x, y, goldDark)
				} else {
					setPixel(img, x, y, gold)
				}
			} else if dist < 7.0 {
				setPixel(img, x, y, darkDial)
			}
		}
	}

	// Top winder crown
	fillRect(img, 10, 0, 4, 3, goldLight)
	setPixel(img, 9, 1, goldDark)
	setPixel(img, 14, 1, goldDark)

	// 12 hour ticks
	setPixel(img, 11, 7, tickRed) // 12 o'clock (Madness deadline!)
	setPixel(img, 12, 7, tickRed)
	setPixel(img, 15, 8, tickCol)  // 1 o'clock
	setPixel(img, 17, 10, tickCol) // 2 o'clock
	setPixel(img, 17, 12, tickCol) // 3 o'clock
	setPixel(img, 17, 13, tickCol)
	setPixel(img, 16, 15, tickCol) // 4 o'clock
	setPixel(img, 14, 17, tickCol) // 5 o'clock
	setPixel(img, 11, 17, tickCol) // 6 o'clock
	setPixel(img, 12, 17, tickCol)
	setPixel(img, 9, 17, tickCol) // 7 o'clock
	setPixel(img, 7, 15, tickCol) // 8 o'clock
	setPixel(img, 6, 12, tickCol) // 9 o'clock
	setPixel(img, 6, 13, tickCol)
	setPixel(img, 7, 10, tickCol) // 10 o'clock
	setPixel(img, 9, 8, tickCol)  // 11 o'clock

	// Center brass pivot pin
	setPixel(img, 11, 12, goldLight)
	setPixel(img, 12, 12, goldLight)
	setPixel(img, 11, 13, goldDark)
	setPixel(img, 12, 13, goldDark)

	return img
}

func createHUDClockHandImage() *ebiten.Image {
	img := ebiten.NewImage(3, 10)
	handCol := color.RGBA{R: 240, G: 230, B: 200, A: 255}
	tipCol := color.RGBA{R: 245, G: 65, B: 65, A: 255} // Sharp red pointer

	// Slender needle pointing upwards from bottom pivot (1, 9) to tip (1, 1)
	setPixel(img, 1, 1, tipCol)
	setPixel(img, 1, 2, tipCol)
	setPixel(img, 0, 3, handCol)
	setPixel(img, 1, 3, handCol)
	setPixel(img, 2, 3, handCol)
	for y := 4; y <= 9; y++ {
		setPixel(img, 1, y, handCol)
	}

	return img
}
