package app

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"panic-recover/internal/game"
)

//go:embed assets/gopher.png assets/bug.png assets/panic-recover.mp3
var embeddedAssets embed.FS

type spriteAssets struct {
	gopher *ebiten.Image
	bug    *ebiten.Image
	music  []byte
}

func loadEmbeddedAssets() (spriteAssets, error) {
	gopher, err := decodeEmbeddedImage("assets/gopher.png")
	if err != nil {
		return spriteAssets{}, err
	}
	bug, err := decodeEmbeddedImage("assets/bug.png")
	if err != nil {
		return spriteAssets{}, err
	}
	music, err := embeddedAssets.ReadFile("assets/panic-recover.mp3")
	if err != nil {
		return spriteAssets{}, fmt.Errorf("read embedded music: %w", err)
	}
	return spriteAssets{gopher: gopher, bug: bug, music: music}, nil
}

func decodeEmbeddedImage(name string) (*ebiten.Image, error) {
	data, err := embeddedAssets.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open embedded asset %s: %w", name, err)
	}
	defer data.Close()
	decoded, _, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("decode embedded asset %s: %w", name, err)
	}
	return ebiten.NewImageFromImage(decoded), nil
}

func drawSprite(screen, sprite *ebiten.Image, position game.Vec2, diameter float64) {
	if sprite == nil {
		return
	}
	bounds := sprite.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())
	options := &ebiten.DrawImageOptions{}
	options.Filter = ebiten.FilterNearest
	options.GeoM.Scale(diameter/width, diameter/height)
	options.GeoM.Translate(position.X-diameter/2, position.Y-diameter/2)
	screen.DrawImage(sprite, options)
}
