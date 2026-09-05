package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/game"
)

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Panic & Recover: Eldritch Gopher")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g, err := game.NewGame()
	if err != nil {
		log.Fatalf("failed to initialize game: %v", err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("game exited with error: %v", err)
	}
}
