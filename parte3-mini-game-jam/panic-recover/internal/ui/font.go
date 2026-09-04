package ui

import (
	"bytes"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
)

var (
	fontOnce   sync.Once
	fontSource *text.GoTextFaceSource
)

// Init initializes the Go Mono font source.
func Init() {
	fontOnce.Do(func() {
		var err error
		fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
		if err != nil {
			panic(err)
		}
	})
}

// FontSource returns the cached text face source.
func FontSource() *text.GoTextFaceSource {
	if fontSource == nil {
		Init()
	}
	return fontSource
}

// DrawText draws text at (x, y) with a given size and color.
func DrawText(dst *ebiten.Image, str string, x, y float64, size float64, clr color.Color) {
	if fontSource == nil {
		Init()
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)

	face := &text.GoTextFace{
		Source: fontSource,
		Size:   size,
	}
	text.Draw(dst, str, face, op)
}

// DrawCenteredText draws horizontally centered text around xCenter.
func DrawCenteredText(dst *ebiten.Image, str string, xCenter, y float64, size float64, clr color.Color) {
	if fontSource == nil {
		Init()
	}
	face := &text.GoTextFace{
		Source: fontSource,
		Size:   size,
	}
	w, _ := text.Measure(str, face, 0)
	x := xCenter - w/2.0

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, str, face, op)
}
