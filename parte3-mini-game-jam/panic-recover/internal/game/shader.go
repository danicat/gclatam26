package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Chromatic Aberration Kage shader
// When panic is >= 80%, this shader applies a subtle, atmospheric separation of RGB channels
// and gentle psychological distortion without excessive visual strain.
const chromaticAberrationShaderSource = `//kage:unit pixels
package main

var Intensity float
var Time float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	if Intensity <= 0.001 {
		return imageSrc0At(srcPos) * color
	}

	// Toned-down subtle chromatic offset (max ~1.8 to 2.2 pixels)
	wobble := sin(Time*3.5 + srcPos.y*0.04) * (Intensity * 0.6)
	offset := vec2(Intensity*2.0 + wobble, 0.0)

	// Sample separated color channels with gentle dispersion
	r := imageSrc0At(srcPos + offset).r
	g := imageSrc0At(srcPos).g
	b := imageSrc0At(srcPos - offset).b
	a := imageSrc0At(srcPos).a

	// Gentle, atmospheric edge vignette (subtle darkening)
	center := vec2(240.0, 135.0) // Virtual screen center (480x270)
	dist := distance(srcPos, center) / 280.0
	vignette := 1.0 - (dist * dist * Intensity * 0.16)

	return vec4(r * vignette, g * vignette, b * vignette, a) * color
}
`

var aberrationShader *ebiten.Shader

func initShader() error {
	var err error
	aberrationShader, err = ebiten.NewShader([]byte(chromaticAberrationShaderSource))
	return err
}
