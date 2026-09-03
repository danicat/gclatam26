package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-at-the-disco/internal/game"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Panic! At The Disco: Saturday Night Flee-ver")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game loop terminated unexpectedly: %v", err)
	}
}
