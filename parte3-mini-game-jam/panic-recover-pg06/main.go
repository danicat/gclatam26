package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"panic-recover/internal/game"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Panic!!! (& recover?) - GopherCon LATAM 2026")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game exited with error: %v", err)
	}
}
