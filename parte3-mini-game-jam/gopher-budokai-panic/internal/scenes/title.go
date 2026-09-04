package scenes

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gopher-budokai-panic/internal/audio"
	"gopher-budokai-panic/internal/gfx"
)

type TitleScene struct {
	arena       *gfx.Arena
	animPhase   float64
	spriteCache *gfx.SpriteCache
}

func NewTitleScene() *TitleScene {
	return &TitleScene{
		arena:       gfx.NewArena(),
		spriteCache: gfx.InitSprites(),
	}
}

func (ts *TitleScene) Enter() {
	audio.Get().PlayBGM("title")
}

func (ts *TitleScene) Exit() {}

func (ts *TitleScene) Update(dt float64) Scene {
	ts.animPhase += dt
	ts.arena.Update(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return NewBattleScene()
	}
	return nil
}

func (ts *TitleScene) Draw(screen *ebiten.Image) {
	w, h := 640.0, 360.0
	ts.arena.Draw(screen, w, h)

	// Draw decorative Goku & Vegeta faceoff in center
	p1Sprite := ts.spriteCache.GetSprite(gfx.FighterPlayer, gfx.PoseIdle)
	p2Sprite := ts.spriteCache.GetSprite(gfx.FighterCPU, gfx.PoseIdle)

	bob1 := math.Sin(ts.animPhase*3.0) * 4.0
	bob2 := math.Cos(ts.animPhase*3.0) * 4.0

	gfx.DrawFighter(screen, p1Sprite, 250, 190+bob1, 1.8, 1.8, 0, false)
	gfx.DrawFighter(screen, p2Sprite, 390, 190+bob2, 1.8, 1.8, 0, true)

	// Title Card Banner
	bannerW := float32(440.0)
	bannerH := float32(70.0)
	bx := float32(w)/2 - bannerW/2
	by := float32(35.0)

	vector.DrawFilledRect(screen, bx-3, by-3, bannerW+6, bannerH+6, color.RGBA{R: 255, G: 200, B: 20, A: 255}, false)
	vector.DrawFilledRect(screen, bx, by, bannerW, bannerH, color.RGBA{R: 15, G: 15, B: 25, A: 240}, false)

	// Title Logos
	ebitenutil.DebugPrintAt(screen, "GOPHER BUDOKAI TENKAICHI", int(bx+110), int(by+12))
	ebitenutil.DebugPrintAt(screen, "★ THEME: PANIC!!! (& RECOVER?) ★", int(bx+90), int(by+32))
	ebitenutil.DebugPrintAt(screen, "GopherCon LATAM 2026 Mini Game Jam", int(bx+115), int(by+50))

	// Press Start Prompt with pulsing color
	pulse := uint8(200 + math.Sin(ts.animPhase*6.0)*55)
	startCol := color.RGBA{R: 255, G: pulse, B: 50, A: 255}
	_ = startCol
	ebitenutil.DebugPrintAt(screen, ">> PRESS [SPACE] OR [ENTER] TO START BATTLE <<", 185, 275)

	// Controls Overview
	vector.DrawFilledRect(screen, 60, 305, 520, 42, color.RGBA{R: 10, G: 10, B: 15, A: 210}, false)
	ebitenutil.DebugPrintAt(screen, "CONTROLS: WASD / Arrows: Move  |  J: Melee Combo  |  K: Ki Blast  |  L: Dragon Dash", 70, 312)
	ebitenutil.DebugPrintAt(screen, "SPACE: Ki Charge (Hold) / MASH TO RECOVER!  |  I: Super Beam  |  Shift: Vanish", 70, 328)
}
