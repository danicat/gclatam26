package art

import (
	"image/color"
	"math"
	"math/rand"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// Predefined cyber/retro color palette
var (
	ColorBgDark        = color.RGBA{8, 10, 16, 255}
	ColorPanelBg       = color.RGBA{16, 22, 34, 245}
	ColorPanelBorder   = color.RGBA{38, 50, 72, 255}
	ColorCodeBg        = color.RGBA{12, 16, 26, 240} // Floating code canvas
	ColorCodeBorder    = color.RGBA{45, 60, 90, 255}
	ColorGutterBg      = color.RGBA{16, 20, 32, 255} // Line section gutter
	ColorGutterBorder  = color.RGBA{55, 75, 110, 255} // Vertical separator
	ColorLineNumber    = color.RGBA{130, 165, 210, 255} // Crisp readable steel blue
	ColorLineNumActive = color.RGBA{0, 255, 255, 255} // Glowing cyan when line is selected
	ColorLineNumEdit   = color.RGBA{255, 215, 64, 255} // Glowing gold when line is in edit

	ColorCyanGlow      = color.RGBA{0, 235, 255, 255}
	ColorGreenRecover  = color.RGBA{105, 255, 174, 255}
	ColorRedPanic      = color.RGBA{255, 61, 0, 255}
	ColorYellowAlert   = color.RGBA{255, 215, 64, 255}
	ColorPurpleToken   = color.RGBA{230, 120, 255, 255}
	ColorComment       = color.RGBA{115, 140, 165, 255}
	ColorTextWhite     = color.RGBA{248, 250, 252, 255}

	// High contrast line selection (Navigation mode)
	ColorSelectBar     = color.RGBA{24, 48, 96, 245} // Rich deep navy fill
	ColorSelectBorder  = color.RGBA{64, 196, 255, 255} // Electric blue border
	ColorSelectAccent  = color.RGBA{0, 235, 255, 255} // Bright left indicator pip

	// High contrast line editing (Edit mode - when selected to change it)
	ColorEditBar       = color.RGBA{48, 36, 12, 250} // Deep warm amber-obsidian fill
	ColorEditBorder    = color.RGBA{255, 193, 7, 255} // Vibrant gold border
	ColorEditAccent    = color.RGBA{255, 235, 59, 255} // Neon amber left indicator
	ColorCursor        = color.RGBA{255, 255, 255, 255} // Pure white blinking cursor for max contrast

	ColorLaserRed      = color.RGBA{255, 23, 68, 255}
	ColorScanline      = color.RGBA{0, 0, 0, 35}
)

var (
	DefaultFace font.Face = basicfont.Face7x13
	CharWidth             = 7
	LineHeight            = 15

	whitePixel *ebiten.Image
)

func init() {
	whitePixel = ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
}

// Particle represents a single visual effect particle in memory.
type Particle struct {
	X, Y   float64
	Vx, Vy float64
	Color  color.RGBA
	Life   float64
	MaxLife float64
	Size   float64
}

// ParticleSystem maintains pre-allocated particle pools for zero GC allocations.
type ParticleSystem struct {
	pool []Particle
	rnd  *rand.Rand
}

func NewParticleSystem(capacity int) *ParticleSystem {
	return &ParticleSystem{
		pool: make([]Particle, 0, capacity),
		rnd:  rand.New(rand.NewSource(1337)),
	}
}

// SpawnSparks bursts particles of given color at (x, y).
func (ps *ParticleSystem) SpawnSparks(x, y float64, count int, clr color.RGBA) {
	for i := 0; i < count; i++ {
		angle := ps.rnd.Float64() * 2.0 * math.Pi
		speed := 40.0 + ps.rnd.Float64()*120.0
		life := 0.4 + ps.rnd.Float64()*0.6
		p := Particle{
			X:       x,
			Y:       y,
			Vx:      math.Cos(angle) * speed,
			Vy:      math.Sin(angle) * speed,
			Color:   clr,
			Life:    life,
			MaxLife: life,
			Size:    2.0 + ps.rnd.Float64()*2.0,
		}
		ps.pool = append(ps.pool, p)
	}
}

// Update advances all particles by dt.
func (ps *ParticleSystem) Update(dt float64) {
	activeIdx := 0
	for i := 0; i < len(ps.pool); i++ {
		p := &ps.pool[i]
		p.Life -= dt
		if p.Life > 0 {
			p.X += p.Vx * dt
			p.Y += p.Vy * dt
			// Gravity or drag
			p.Vy += 60.0 * dt
			ps.pool[activeIdx] = *p
			activeIdx++
		}
	}
	ps.pool = ps.pool[:activeIdx]
}

// Draw renders particles onto the screen.
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, p := range ps.pool {
		alpha := p.Life / p.MaxLife
		if alpha <= 0 {
			continue
		}
		c := color.RGBA{
			R: p.Color.R,
			G: p.Color.G,
			B: p.Color.B,
			A: uint8(float64(p.Color.A) * alpha),
		}
		DrawRect(screen, int(p.X), int(p.Y), int(p.Size), int(p.Size), c)
	}
}

// MatrixRain manages retro digital matrix rain background lines.
type MatrixDrop struct {
	X     float64
	Y     float64
	Speed float64
	Char  rune
	Alpha float64
}

type MatrixRain struct {
	drops []MatrixDrop
	rnd   *rand.Rand
	timer float64
}

func NewMatrixRain(count int, screenW, screenH int) *MatrixRain {
	mr := &MatrixRain{
		drops: make([]MatrixDrop, count),
		rnd:   rand.New(rand.NewSource(42)),
	}
	chars := []rune("0101{};:=_&*%#@")
	for i := range mr.drops {
		mr.drops[i] = MatrixDrop{
			X:     float64(mr.rnd.Intn(screenW)),
			Y:     float64(mr.rnd.Intn(screenH)),
			Speed: 15.0 + mr.rnd.Float64()*45.0,
			Char:  chars[mr.rnd.Intn(len(chars))],
			Alpha: 0.15 + mr.rnd.Float64()*0.25,
		}
	}
	return mr
}

func (mr *MatrixRain) Update(dt float64, screenH int) {
	mr.timer += dt
	chars := []rune("0101{};:=_&*%#@")
	for i := range mr.drops {
		mr.drops[i].Y += mr.drops[i].Speed * dt
		if mr.drops[i].Y > float64(screenH) {
			mr.drops[i].Y = -10
			mr.drops[i].Char = chars[mr.rnd.Intn(len(chars))]
		}
	}
}

func (mr *MatrixRain) Draw(screen *ebiten.Image) {
	for _, d := range mr.drops {
		clr := color.RGBA{0, 255, 120, uint8(d.Alpha * 255)}
		text.Draw(screen, string(d.Char), DefaultFace, int(d.X), int(d.Y), clr)
	}
}

// ============================================================================
// Drawing Helpers
// ============================================================================

// DrawRect draws a filled rectangle.
func DrawRect(dst *ebiten.Image, x, y, w, h int, clr color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(h))
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	dst.DrawImage(whitePixel, op)
}

// DrawBorder draws a 1-pixel hollow rectangle border.
func DrawBorder(dst *ebiten.Image, x, y, w, h int, clr color.Color) {
	DrawRect(dst, x, y, w, 1, clr)         // Top
	DrawRect(dst, x, y+h-1, w, 1, clr)     // Bottom
	DrawRect(dst, x, y, 1, h, clr)         // Left
	DrawRect(dst, x+w-1, y, 1, h, clr)     // Right
}

// DrawText draws single colored text with basicfont face.
func DrawText(dst *ebiten.Image, str string, x, y int, clr color.Color) {
	text.Draw(dst, str, DefaultFace, x, y, clr)
}

// SyntaxToken represents a colored fragment of a code line.
type SyntaxToken struct {
	Text  string
	Color color.RGBA
}

var (
	keywordRegex = regexp.MustCompile(`\b(package|func|var|type|struct|return|if|else|for|go|defer|chan|make|close|panic|recover|import|select|case|default|switch|len|println)\b`)
	typeRegex    = regexp.MustCompile(`\b(int|string|bool|error|any|byte|float64|uint|uint8|int64|Gopher|User|map)\b`)
	numberRegex  = regexp.MustCompile(`\b\d+\b`)
)

// HighlightGoLine performs lightweight regex-based tokenization for Go code syntax coloring.
func HighlightGoLine(line string) []SyntaxToken {
	if strings.HasPrefix(strings.TrimSpace(line), "//") {
		return []SyntaxToken{{Text: line, Color: ColorComment}}
	}

	// Simple tokenizer: split by words and symbols while preserving spacing
	tokens := make([]SyntaxToken, 0)
	i := 0
	runes := []rune(line)
	n := len(runes)

	for i < n {
		// Check string literal
		if runes[i] == '"' {
			start := i
			i++
			for i < n && runes[i] != '"' {
				if runes[i] == '\\' && i+1 < n {
					i += 2
				} else {
					i++
				}
			}
			if i < n && runes[i] == '"' {
				i++
			}
			tokens = append(tokens, SyntaxToken{Text: string(runes[start:i]), Color: ColorGreenRecover})
			continue
		}

		// Check comment till end of line
		if i+1 < n && runes[i] == '/' && runes[i+1] == '/' {
			tokens = append(tokens, SyntaxToken{Text: string(runes[i:]), Color: ColorComment})
			break
		}

		// Check word token (keywords, types, identifiers)
		if isAlpha(runes[i]) || runes[i] == '_' {
			start := i
			for i < n && (isAlphaNum(runes[i]) || runes[i] == '_') {
				i++
			}
			word := string(runes[start:i])
			clr := ColorTextWhite
			if keywordRegex.MatchString(word) {
				clr = ColorCyanGlow
			} else if typeRegex.MatchString(word) {
				clr = ColorYellowAlert
			}
			tokens = append(tokens, SyntaxToken{Text: word, Color: clr})
			continue
		}

		// Check numbers
		if isDigit(runes[i]) {
			start := i
			for i < n && isDigit(runes[i]) {
				i++
			}
			tokens = append(tokens, SyntaxToken{Text: string(runes[start:i]), Color: ColorPurpleToken})
			continue
		}

		// Symbols, operators, spaces
		tokens = append(tokens, SyntaxToken{Text: string(runes[i]), Color: ColorTextWhite})
		i++
	}

	return tokens
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

// DrawHighlightedLine draws a tokenized syntax-colored code line.
func DrawHighlightedLine(dst *ebiten.Image, line string, x, y int) {
	tokens := HighlightGoLine(line)
	curX := x
	for _, tok := range tokens {
		text.Draw(dst, tok.Text, DefaultFace, curX, y, tok.Color)
		curX += len(tok.Text) * CharWidth
	}
}

// DrawScanlines draws subtle retro CRT horizontal scanlines across the screen.
func DrawScanlines(dst *ebiten.Image, screenW, screenH int) {
	for y := 0; y < screenH; y += 2 {
		DrawRect(dst, 0, y, screenW, 1, ColorScanline)
	}
}
