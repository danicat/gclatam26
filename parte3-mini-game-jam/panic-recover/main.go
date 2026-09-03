package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"panic-recover/internal/app"
)

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Panic Recover")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(app.New()); err != nil {
		log.Fatal(err)
	}
}
