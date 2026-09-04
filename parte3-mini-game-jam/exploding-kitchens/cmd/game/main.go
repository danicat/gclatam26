package main

import (
	"log"

	"exploding-kitchens/internal/game"
	"exploding-kitchens/internal/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Exploding Kitchens - Panic and Recover")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.TargetFPS)

	titleScene := scene.NewTitleScene()
	playScene := scene.NewPlayScene()
	gameOverScene := scene.NewGameOverScene()

	g := game.NewGame(titleScene, playScene, gameOverScene)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
