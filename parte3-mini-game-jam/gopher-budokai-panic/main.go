package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"gopher-budokai-panic/internal/game"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Gopher Budokai: Panic & Recover! - GopherCon LATAM 2026 Game Jam")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game loop crashed: %v", err)
	}
}
