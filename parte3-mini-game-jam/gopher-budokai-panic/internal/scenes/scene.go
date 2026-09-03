package scenes

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Enter()
	Update(dt float64) Scene
	Draw(screen *ebiten.Image)
	Exit()
}
