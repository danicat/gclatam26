package scenes

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gopher-budokai-panic/internal/audio"
	"gopher-budokai-panic/internal/gfx"
)

type GameOverScene struct {
	won         bool
	stats       MatchStats
	arena       *gfx.Arena
	animPhase   float64
	spriteCache *gfx.SpriteCache
}

func NewGameOverScene(won bool, stats MatchStats) *GameOverScene {
	return &GameOverScene{
		won:         won,
		stats:       stats,
		arena:       gfx.NewArena(),
		spriteCache: gfx.InitSprites(),
	}
}

func (gos *GameOverScene) Enter() {
	if gos.won {
		audio.Get().PlayBGM("title")
	} else {
		audio.Get().StopBGM()
	}
}

func (gos *GameOverScene) Exit() {}

func (gos *GameOverScene) Update(dt float64) Scene {
	gos.animPhase += dt
	gos.arena.Update(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return NewBattleScene() // Rematch!
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return NewTitleScene()
	}
	return nil
}

func (gos *GameOverScene) Draw(screen *ebiten.Image) {
	w, h := 640.0, 360.0
	gos.arena.Draw(screen, w, h)

	// Dim overlay
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), color.RGBA{R: 0, G: 0, B: 0, A: 160}, false)

	// Result Card
	cardW := float32(400.0)
	cardH := float32(230.0)
	cx := float32(w)/2 - cardW/2
	cy := float32(h)/2 - cardH/2 - 10

	var borderColor color.RGBA
	if gos.won {
		borderColor = color.RGBA{R: 255, G: 215, B: 0, A: 255} // Gold
	} else {
		borderColor = color.RGBA{R: 255, G: 50, B: 50, A: 255} // Red
	}

	vector.DrawFilledRect(screen, cx-3, cy-3, cardW+6, cardH+6, borderColor, false)
	vector.DrawFilledRect(screen, cx, cy, cardW, cardH, color.RGBA{R: 20, G: 20, B: 30, A: 245}, false)

	// Title Announcement
	pulse := uint8(200 + math.Sin(gos.animPhase*5.0)*55)
	_ = pulse
	if gos.won {
		ebitenutil.DebugPrintAt(screen, "★ ★ ★  VICTORY !  ★ ★ ★", int(cx+120), int(cy+25))
		ebitenutil.DebugPrintAt(screen, "Goku prevailed through the panic and conquered the battle!", int(cx+30), int(cy+55))
	} else {
		ebitenutil.DebugPrintAt(screen, "! ! !  DEFEAT  ! ! !", int(cx+135), int(cy+25))
		ebitenutil.DebugPrintAt(screen, "Overwhelmed by panic! Train harder to master recovery!", int(cx+35), int(cy+55))
	}

	// Match Recap Stats (Panic & Recover Metrics)
	vector.DrawFilledRect(screen, cx+25, cy+85, cardW-50, 65, color.RGBA{R: 10, G: 10, B: 20, A: 200}, false)
	rec1 := fmt.Sprintf("Kiai Panic Recoveries : %d", gos.stats.RecoveriesCount)
	rec2 := fmt.Sprintf("Super Beams Unleashed : %d", gos.stats.BeamsFired)
	ebitenutil.DebugPrintAt(screen, rec1, int(cx+45), int(cy+98))
	ebitenutil.DebugPrintAt(screen, rec2, int(cx+45), int(cy+120))

	// Rematch Action Prompt
	ebitenutil.DebugPrintAt(screen, ">> PRESS [SPACE] OR [ENTER] TO REMATCH <<", int(cx+55), int(cy+170))
	ebitenutil.DebugPrintAt(screen, "PRESS [ESC] FOR MAIN MENU", int(cx+120), int(cy+195))
}
