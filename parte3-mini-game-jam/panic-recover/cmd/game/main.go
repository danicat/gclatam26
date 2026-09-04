package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-recover/internal/game"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Panic Recover: Runtime Defender")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("game exited with error: %v", err)
	}
}
