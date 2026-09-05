package game

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// 4x6 Pixel Font for retro UI rendering without external TTF dependencies
var fontGlyphs = map[byte][]string{
	'A': {" ## ", "#  #", "####", "#  #", "#  #", "#  #"},
	'B': {"### ", "#  #", "### ", "#  #", "#  #", "### "},
	'C': {" ###", "#   ", "#   ", "#   ", "#   ", " ###"},
	'D': {"##  ", "# # ", "#  #", "#  #", "# # ", "##  "},
	'E': {"####", "#   ", "### ", "#   ", "#   ", "####"},
	'F': {"####", "#   ", "### ", "#   ", "#   ", "#   "},
	'G': {" ###", "#   ", "# ##", "#  #", "#  #", " ###"},
	'H': {"#  #", "#  #", "####", "#  #", "#  #", "#  #"},
	'I': {"###", " # ", " # ", " # ", " # ", "###"},
	'J': {"  ##", "   #", "   #", "   #", "#  #", " ## "},
	'K': {"#  #", "# # ", "##  ", "# # ", "#  #", "#  #"},
	'L': {"#   ", "#   ", "#   ", "#   ", "#   ", "####"},
	'M': {"#  #", "####", "#  #", "#  #", "#  #", "#  #"},
	'N': {"#  #", "## #", "# ##", "#  #", "#  #", "#  #"},
	'O': {" ## ", "#  #", "#  #", "#  #", "#  #", " ## "},
	'P': {"### ", "#  #", "### ", "#   ", "#   ", "#   "},
	'Q': {" ## ", "#  #", "#  #", "# ##", " ###", "   #"},
	'R': {"### ", "#  #", "### ", "# # ", "#  #", "#  #"},
	'S': {" ###", "#   ", " ## ", "   #", "   #", "### "},
	'T': {"###", " # ", " # ", " # ", " # ", " # "},
	'U': {"#  #", "#  #", "#  #", "#  #", "#  #", " ## "},
	'V': {"#  #", "#  #", "#  #", "#  #", " ## ", "  # "},
	'W': {"#  #", "#  #", "#  #", "#  #", "####", "#  #"},
	'X': {"#  #", " ## ", "  # ", " ## ", "#  #", "#  #"},
	'Y': {"#  #", "#  #", " ## ", "  # ", "  # ", "  # "},
	'Z': {"####", "   #", "  # ", " #  ", "#   ", "####"},
	'0': {" ## ", "#  #", "#  #", "#  #", "#  #", " ## "},
	'1': {" # ", "## ", " # ", " # ", " # ", "###"},
	'2': {"### ", "   #", " ###", "#   ", "#   ", "####"},
	'3': {"### ", "   #", " ## ", "   #", "   #", "### "},
	'4': {"#  #", "#  #", "####", "   #", "   #", "   #"},
	'5': {"####", "#   ", "### ", "   #", "   #", "### "},
	'6': {" ###", "#   ", "### ", "#  #", "#  #", " ###"},
	'7': {"####", "   #", "  # ", "  # ", " #  ", " #  "},
	'8': {" ## ", "#  #", " ## ", "#  #", "#  #", " ## "},
	'9': {" ###", "#  #", " ###", "   #", "   #", "### "},
	':': {"", " #", "", "", " #", ""},
	'-': {"", "", "####", "", "", ""},
	'/': {"   #", "  # ", "  # ", " #  ", " #  ", "#   "},
	'%': {"#  #", "  # ", " #  ", " #  ", "#  #", "#  #"},
	'!': {" # ", " # ", " # ", " # ", "   ", " # "},
	'.': {"", "", "", "", "##", "##"},
	',': {"", "", "", "", "##", " #"},
	'(': {" #", "# ", "# ", "# ", "# ", " #"},
	')': {"# ", " #", " #", " #", " #", "# "},
	' ': {"    ", "    ", "    ", "    ", "    ", "    "},
}

func drawText(dst *ebiten.Image, str string, x, y int, c color.Color, scale int) {
	if scale < 1 {
		scale = 1
	}
	upper := strings.ToUpper(str)
	cursorX := x

	for i := 0; i < len(upper); i++ {
		ch := upper[i]
		rows, ok := fontGlyphs[ch]
		if !ok {
			cursorX += 4 * scale
			continue
		}

		maxGlyphW := 0
		for gy, row := range rows {
			if len(row) > maxGlyphW {
				maxGlyphW = len(row)
			}
			for gx := 0; gx < len(row); gx++ {
				if row[gx] == '#' {
					// Draw pixel scaled
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							px := cursorX + (gx * scale) + sx
							py := y + (gy * scale) + sy
							if px >= 0 && px < dst.Bounds().Dx() && py >= 0 && py < dst.Bounds().Dy() {
								dst.Set(px, py, c)
							}
						}
					}
				}
			}
		}

		if maxGlyphW == 0 {
			maxGlyphW = 3
		}
		cursorX += (maxGlyphW + 1) * scale
	}
}

func getTextWidth(str string, scale int) int {
	if scale < 1 {
		scale = 1
	}
	upper := strings.ToUpper(str)
	w := 0
	for i := 0; i < len(upper); i++ {
		ch := upper[i]
		rows, ok := fontGlyphs[ch]
		if !ok {
			w += 4 * scale
			continue
		}
		maxGlyphW := 0
		for _, row := range rows {
			if len(row) > maxGlyphW {
				maxGlyphW = len(row)
			}
		}
		if maxGlyphW == 0 {
			maxGlyphW = 3
		}
		w += (maxGlyphW + 1) * scale
	}
	return w
}
