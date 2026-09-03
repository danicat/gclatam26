package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// PixelFont renders crisp retro 4x5 pixel typography without any external TTF dependencies.
type PixelFont struct {
	pixelImg *ebiten.Image
	drawOpts ebiten.DrawImageOptions
}

// NewPixelFont creates a font renderer with a cached 1x1 white pixel texture.
func NewPixelFont() *PixelFont {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return &PixelFont{pixelImg: img}
}

// 4x5 character matrix definitions (bitmask: 20 bits per char, 4 bits per row, 5 rows)
// Row 0: bits 16-19, Row 1: bits 12-15, Row 2: bits 8-11, Row 3: bits 4-7, Row 4: bits 0-3
var glyphMap = map[rune]uint32{
	'A': 0b0110_1001_1111_1001_1001,
	'B': 0b1110_1001_1110_1001_1110,
	'C': 0b0111_1000_1000_1000_0111,
	'D': 0b1110_1001_1001_1001_1110,
	'E': 0b1111_1000_1110_1000_1111,
	'F': 0b1111_1000_1110_1000_1000,
	'G': 0b0111_1000_1011_1001_0111,
	'H': 0b1001_1001_1111_1001_1001,
	'I': 0b1110_0100_0100_0100_1110,
	'J': 0b0011_0001_0001_1001_0110,
	'K': 0b1001_1010_1100_1010_1001,
	'L': 0b1000_1000_1000_1000_1111,
	'M': 0b1001_1111_1111_1001_1001,
	'N': 0b1001_1101_1011_1001_1001,
	'O': 0b0110_1001_1001_1001_0110,
	'P': 0b1110_1001_1110_1000_1000,
	'Q': 0b0110_1001_1001_1011_0111,
	'R': 0b1110_1001_1110_1010_1001,
	'S': 0b0111_1000_0110_0001_1110,
	'T': 0b1111_0100_0100_0100_0100,
	'U': 0b1001_1001_1001_1001_0110,
	'V': 0b1001_1001_1001_0110_0100,
	'W': 0b1001_1001_1111_1111_1001,
	'X': 0b1001_0110_0100_0110_1001,
	'Y': 0b1001_1001_0110_0100_0100,
	'Z': 0b1111_0001_0110_1000_1111,

	'0': 0b0110_1001_1001_1001_0110,
	'1': 0b0100_1100_0100_0100_1110,
	'2': 0b1110_0001_0110_1000_1111,
	'3': 0b1110_0001_0110_0001_1110,
	'4': 0b1001_1001_1111_0001_0001,
	'5': 0b1111_1000_1110_0001_1110,
	'6': 0b0111_1000_1110_1001_0110,
	'7': 0b1111_0001_0010_0100_0100,
	'8': 0b0110_1001_0110_1001_0110,
	'9': 0b0110_1001_0111_0001_1110,

	':': 0b0000_0100_0000_0100_0000,
	'!': 0b0100_0100_0100_0000_0100,
	'?': 0b1110_0001_0110_0000_0100,
	'-': 0b0000_0000_1111_0000_0000,
	'+': 0b0000_0100_1110_0100_0000,
	'/': 0b0001_0010_0100_1000_0000,
	'%': 0b1001_0010_0100_0100_1001,
	'.': 0b0000_0000_0000_0000_0100,
	',': 0b0000_0000_0000_0100_1000,
	'\'': 0b0100_0100_0000_0000_0000,
	'[': 0b0110_0100_0100_0100_0110,
	']': 0b0110_0010_0010_0010_0110,
	'(': 0b0010_0100_0100_0100_0010,
	')': 0b0100_0010_0010_0010_0100,
	'<': 0b0010_0100_1000_0100_0010,
	'>': 0b0100_0010_0001_0010_0100,
	'*': 0b0000_1010_0100_1010_0000,
	'x': 0b0000_1001_0110_0110_1001,
	' ': 0b0000_0000_0000_0000_0000,
}

// DrawText renders string text with shadow, scaling, and custom color.
func (f *PixelFont) DrawText(screen *ebiten.Image, text string, x, y float64, scale float64, c color.RGBA, shadow bool) {
	text = strings.ToUpper(text)

	if shadow {
		shadowCol := color.RGBA{15, 12, 28, 200}
		f.drawTextInternal(screen, text, x+scale, y+scale, scale, shadowCol)
	}
	f.drawTextInternal(screen, text, x, y, scale, c)
}

func (f *PixelFont) drawTextInternal(screen *ebiten.Image, text string, startX, startY float64, scale float64, c color.RGBA) {
	curX := startX
	curY := startY

	f.drawOpts.ColorScale.Reset()
	f.drawOpts.ColorScale.Scale(
		float32(c.R)/255.0*(float32(c.A)/255.0),
		float32(c.G)/255.0*(float32(c.A)/255.0),
		float32(c.B)/255.0*(float32(c.A)/255.0),
		float32(c.A)/255.0,
	)

	charWidth := 4.0 * scale
	charHeight := 5.0 * scale
	spacing := 1.0 * scale

	for _, ch := range text {
		if ch == '\n' {
			curX = startX
			curY += charHeight + 3*scale
			continue
		}

		mask, exists := glyphMap[ch]
		if !exists {
			mask = glyphMap['?']
		}

		for row := 0; row < 5; row++ {
			for col := 0; col < 4; col++ {
				shift := uint(19 - (row*4 + col))
				if (mask>>shift)&1 == 1 {
					f.drawOpts.GeoM.Reset()
					f.drawOpts.GeoM.Scale(scale, scale)
					f.drawOpts.GeoM.Translate(curX+float64(col)*scale, curY+float64(row)*scale)
					screen.DrawImage(f.pixelImg, &f.drawOpts)
				}
			}
		}
		curX += charWidth + spacing
	}
}

// MeasureText calculates the width and height of a string in virtual pixels.
func (f *PixelFont) MeasureText(text string, scale float64) (float64, float64) {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	w := float64(maxLen) * (5.0 * scale)
	h := float64(len(lines)) * (8.0 * scale)
	return w, h
}
