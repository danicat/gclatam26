package assets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Palette definitions
var (
	// Hero: Good, healthy, blue Gopher
	ColorGopherCyan = color.RGBA{0x00, 0xAD, 0xD8, 0xFF}
	ColorGopherDark = color.RGBA{0x00, 0x8C, 0xB0, 0xFF}
	ColorEarPink    = color.RGBA{0xF5, 0xA9, 0xB8, 0xFF}
	ColorWhite      = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	ColorBlack      = color.RGBA{0x10, 0x14, 0x1A, 0xFF}
	ColorShieldBlue = color.RGBA{0x61, 0xAF, 0xEF, 0xFF}
	ColorGold       = color.RGBA{0xE5, 0xC0, 0x7B, 0xFF}
	ColorDeferGreen = color.RGBA{0x98, 0xC3, 0x79, 0xFF}

	// Enemies: Bad, corrupted, dark-red Gophers
	ColorBadGopherRed   = color.RGBA{0x9E, 0x1A, 0x24, 0xFF} // Dark red fur
	ColorBadGopherDark  = color.RGBA{0x5E, 0x0C, 0x12, 0xFF} // Deep crimson shading
	ColorBadGopherHorn  = color.RGBA{0x3B, 0x05, 0x08, 0xFF} // Spiky horns
	ColorEvilEyeYellow  = color.RGBA{0xFF, 0xCC, 0x00, 0xFF} // Menacing glowing eyes
	ColorEvilEyeRed     = color.RGBA{0xFF, 0x33, 0x44, 0xFF} // Blood red pupils
	ColorCorruptedCore  = color.RGBA{0xDC, 0x26, 0x26, 0xFF} // Bright danger red
)

type Sprites struct {
	Player         *ebiten.Image
	PlayerShielded *ebiten.Image
	InvaderNil     [2]*ebiten.Image
	InvaderIndex   [2]*ebiten.Image
	InvaderDivide  [2]*ebiten.Image
	UFO            *ebiten.Image
	Boss           *ebiten.Image
	BulletHero     *ebiten.Image
	BulletPanic    *ebiten.Image
	PowerupMutex   *ebiten.Image
	PowerupChan    *ebiten.Image
	PowerupTimeout *ebiten.Image
	PowerupBadge   *ebiten.Image
}

var LoadedSprites *Sprites

func InitSprites() {
	LoadedSprites = &Sprites{
		Player:         createPlayerSprite(false),
		PlayerShielded: createPlayerSprite(true),
		InvaderNil:     [2]*ebiten.Image{createBadGopherScout(0), createBadGopherScout(1)},
		InvaderIndex:   [2]*ebiten.Image{createBadGopherSpiky(0), createBadGopherSpiky(1)},
		InvaderDivide:  [2]*ebiten.Image{createBadGopherBrute(0), createBadGopherBrute(1)},
		UFO:            createBadGopherUFO(),
		Boss:           createBadGopherBoss(),
		BulletHero:     createHeroBullet(),
		BulletPanic:    createPanicBullet(),
		PowerupMutex:   createPowerupSprite(ColorShieldBlue, "M"),
		PowerupChan:    createPowerupSprite(ColorGold, "C"),
		PowerupTimeout: createPowerupSprite(ColorDeferGreen, "T"),
		PowerupBadge:   createPowerupSprite(ColorWhite, "*"),
	}
}

func parsePixelArt(grid []string, palette map[rune]color.Color) *ebiten.Image {
	h := len(grid)
	w := len(grid[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range grid {
		for x, char := range row {
			if col, ok := palette[char]; ok && col != nil {
				img.Set(x, y, col)
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

// Healthy, good, blue Gopher hero with cute big eyes and teeth
func createPlayerSprite(shielded bool) *ebiten.Image {
	grid := []string{
		"......GGGGGGGG......",
		"....GGGGGGGGGGGG....",
		"...GGEEGGGGGGGEEGG...",
		"...GEEEGGGGGGGEEEG...",
		"..GGEEGGGGGGGGGEEGG..",
		"..GGGGWWGGGGWWGGGG..",
		"..GGGWKKWGGWKKWGGG..",
		"..GGGWKKWGGWKKWGGG..",
		"..GGGGWWGGGGWWGGGG..",
		"...GGGGGGKKGGGGGG...",
		"....GGGGWWWWGGGG....",
		"....GGGGWWWWGGGG....",
		"...GGGGGGGGGGGGGG...",
		"..GGGG..CCCC..GGGG..",
		".GGGGG..CCCC..GGGGG.",
		"GGGGGG..CCCC..GGGGGG",
	}
	pal := map[rune]color.Color{
		'G': ColorGopherCyan,
		'E': ColorEarPink,
		'W': ColorWhite,
		'K': ColorBlack,
		'C': color.RGBA{0x56, 0xB6, 0xC2, 0xFF},
	}
	if shielded {
		pal['G'] = ColorShieldBlue
		pal['C'] = ColorGold
		pal['E'] = ColorGold
	}
	return parsePixelArt(grid, pal)
}

// Bad Gopher Scout (Tier 1: panic("nil pointer dereference"))
func createBadGopherScout(frame int) *ebiten.Image {
	var grid []string
	if frame == 0 {
		grid = []string{
			"..H..........H..",
			".HH..RRRRRR..HH.",
			".H..RRRRRRRR..H.",
			"...RRRRRRRRRR...",
			"...RRYKRRRRYKR..",
			"...RRYKRRRRYKR..",
			"...RRRRKKRRRR...",
			"....RRWWWWWRR...",
			"....RRW.W.WRR...",
			"...RRRRRRRRRR...",
			"..RR..RRRR..RR..",
			"..R..........R..",
		}
	} else {
		grid = []string{
			"..H..........H..",
			".HH..RRRRRR..HH.",
			".H..RRRRRRRR..H.",
			"...RRRRRRRRRR...",
			"...RRYKRRRRYKR..",
			"...RRYKRRRRYKR..",
			"...RRRRKKRRRR...",
			"....RRWWWWWRR...",
			"....RRW.W.WRR...",
			"...RRRRRRRRRR...",
			"...R..RRRR..R...",
			"....RR....RR....",
		}
	}
	pal := map[rune]color.Color{
		'R': ColorBadGopherRed,
		'H': ColorBadGopherHorn,
		'Y': ColorEvilEyeYellow,
		'K': ColorBlack,
		'W': ColorWhite,
	}
	return parsePixelArt(grid, pal)
}

// Bad Gopher Spiky (Tier 2: panic("index out of range"))
func createBadGopherSpiky(frame int) *ebiten.Image {
	var grid []string
	if frame == 0 {
		grid = []string{
			".H............H.",
			".HH..RRRRRR..HH.",
			"..H.RRRRRRRR.H..",
			"...RRRRRRRRRR...",
			"..RRRYKRRRYKRR..",
			"..RRRYKRRRYKRR..",
			"..RRRRRKKRRRRR..",
			"...RRRWWWWWRR...",
			"...RRRW.W.WRR...",
			"..RRRRRRRRRRRR..",
			".RR..RR..RR..RR.",
			".R............R.",
		}
	} else {
		grid = []string{
			".H............H.",
			".HH..RRRRRR..HH.",
			"..H.RRRRRRRR.H..",
			"...RRRRRRRRRR...",
			"..RRRYKRRRYKRR..",
			"..RRRYKRRRYKRR..",
			"..RRRRRKKRRRRR..",
			"...RRRWWWWWRR...",
			"...RRRW.W.WRR...",
			"..RRRRRRRRRRRR..",
			"..RR.RR..RR.RR..",
			"...R........R...",
		}
	}
	pal := map[rune]color.Color{
		'R': ColorBadGopherRed,
		'H': ColorBadGopherHorn,
		'Y': ColorEvilEyeRed,
		'K': ColorBlack,
		'W': ColorWhite,
	}
	return parsePixelArt(grid, pal)
}

// Bad Gopher Brute (Tier 3: panic("integer divide by zero"))
func createBadGopherBrute(frame int) *ebiten.Image {
	var grid []string
	if frame == 0 {
		grid = []string{
			"HH............HH",
			"HHH..RRRRRR..HHH",
			".HH.RRRRRRRR.HH.",
			"...RRRRRRRRRR...",
			"..RRRYKRRRYKRR..",
			"..RRRYKRRRYKRR..",
			"..RRRRRKKRRRRR..",
			"...RRWWWWWRRR...",
			"...RRW.W.WRRR...",
			"..RRRRRRRRRRRR..",
			".RRRRRRRRRRRRRR.",
			"RR............RR",
		}
	} else {
		grid = []string{
			"HH............HH",
			"HHH..RRRRRR..HHH",
			".HH.RRRRRRRR.HH.",
			"...RRRRRRRRRR...",
			"..RRRYKRRRYKRR..",
			"..RRRYKRRRYKRR..",
			"..RRRRRKKRRRRR..",
			"...RRWWWWWRRR...",
			"...RRW.W.WRRR...",
			"..RRRRRRRRRRRR..",
			"..RRRRRRRRRRRR..",
			"...RR......RR...",
		}
	}
	pal := map[rune]color.Color{
		'R': ColorBadGopherDark,
		'H': ColorBadGopherHorn,
		'Y': ColorEvilEyeYellow,
		'K': ColorBlack,
		'W': ColorWhite,
	}
	return parsePixelArt(grid, pal)
}

// Bad Gopher UFO Drone (Flying dark-red UFO with bad gopher head)
func createBadGopherUFO() *ebiten.Image {
	grid := []string{
		".......DHRRHD.......",
		"....DDHRRRRRRHDD....",
		"...DHRRYKRRRYKRD....",
		"..DHRRRYKRRRYKRRHD..",
		".DDHRRRWWWWWWRRHDDD.",
		"DDDDDDDDDDDDDDDDDDDD",
		"RRRRRRRRRRRRRRRRRRRR",
		".RR..YY..YY..YY..RR.",
		"..R..............R..",
		"....................",
	}
	pal := map[rune]color.Color{
		'R': ColorBadGopherRed,
		'D': ColorBadGopherDark,
		'H': ColorBadGopherHorn,
		'Y': ColorEvilEyeYellow,
		'K': ColorBlack,
		'W': ColorWhite,
	}
	return parsePixelArt(grid, pal)
}

// Corrupted Mega-Bad-Gopher Boss
func createBadGopherBoss() *ebiten.Image {
	grid := []string{
		"HHH..........................................HHH",
		"HHHH........................................HHHH",
		"HHHHH........DDDDDDDDDDDDDDDDDDDD..........HHHHH",
		".HHHHH.....DDDDDDDDDDDDDDDDDDDDDDDD.......HHHHH.",
		"..HHHH...DDDRRRRRRRRRRRRRRRRRRRRRRDDD.....HHHH..",
		"...HH...DDRRRRRRRRRRRRRRRRRRRRRRRRRRDD.....HH...",
		".......DDRRRRRRRRRRRRRRRRRRRRRRRRRRRRDD.........",
		"......DDRRRYYYKKRRRRRRRRRRRRRRRYYYKKRRDD........",
		"......DDRRRYYYKKRRRRRRRRRRRRRRRYYYKKRRDD........",
		"......DDRRRKKKKKRRRRRRRRRRRRRRRKKKKKRRDD........",
		"......DDRRRRRRRRRRRRRKKKKRRRRRRRRRRRRRDD........",
		"......DDRRRRRRRRRRRRKKKKKKRRRRRRRRRRRRDD........",
		"......DDRRRRRRRWWWWWWWWWWWWWWWWWRRRRRRDD........",
		"......DDRRRRRRRW..W..W..W..W..WWRRRRRRDD........",
		".......DDRRRRRRW..W..W..W..W..WWRRRRRRDD........",
		"........DDRRRRRWWWWWWWWWWWWWWWWWRRRRRDD.........",
		".........DDRRRRRRRRRRRRRRRRRRRRRRRRRDD..........",
		"..........DDRRRRRRRRCCCCCCRRRRRRRRRDD...........",
		".........DDDDRRRRRRRCCCCCCRRRRRRRRDDDD..........",
		"........DD..DDRRRRRRCCCCCCRRRRRRDD..DD..........",
		".......DD....DDDRRRRRRRRRRRRRRDDD....DD.........",
		"......DD......DDDDDDDDDDDDDDDDDD......DD........",
		".....DD........DD............DD........DD.......",
		"....DD..........D............D..........DD......",
	}
	pal := map[rune]color.Color{
		'R': ColorBadGopherRed,
		'D': ColorBadGopherDark,
		'H': ColorBadGopherHorn,
		'Y': ColorEvilEyeYellow,
		'K': ColorBlack,
		'W': ColorWhite,
		'C': ColorCorruptedCore,
	}
	return parsePixelArt(grid, pal)
}

func createHeroBullet() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 3, 10))
	for y := 0; y < 10; y++ {
		img.Set(1, y, ColorWhite)
		img.Set(0, y, ColorGopherCyan)
		img.Set(2, y, ColorGopherCyan)
	}
	return ebiten.NewImageFromImage(img)
}

// Corrupted dark-red panic spikes
func createPanicBullet() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 3, 8))
	for y := 0; y < 8; y++ {
		c := ColorBadGopherRed
		if y%2 == 0 {
			c = ColorCorruptedCore
		}
		img.Set(1, y, ColorWhite)
		img.Set(0, y, c)
		img.Set(2, y, c)
	}
	return ebiten.NewImageFromImage(img)
}

func createPowerupSprite(c color.Color, label string) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 12, 12))
	for y := 0; y < 12; y++ {
		for x := 0; x < 12; x++ {
			if x == 0 || x == 11 || y == 0 || y == 11 {
				img.Set(x, y, c)
			} else {
				img.Set(x, y, ColorBlack)
			}
		}
	}
	// Center marker
	img.Set(5, 5, c)
	img.Set(6, 5, c)
	img.Set(5, 6, c)
	img.Set(6, 6, c)
	return ebiten.NewImageFromImage(img)
}
