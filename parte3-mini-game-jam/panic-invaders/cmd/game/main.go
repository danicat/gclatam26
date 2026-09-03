package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-invaders/internal/game"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Panic Invaders: In recover() We Trust - GopherCon LATAM 2026")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
