package main

import (
	"log"

	"box-boy/internal/config"
	"box-boy/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Erro fatal ao rodar BoxBoy: %v", err)
	}
}
