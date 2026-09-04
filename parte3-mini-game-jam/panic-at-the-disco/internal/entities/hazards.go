package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"panic-at-the-disco/internal/audio"
	"panic-at-the-disco/internal/gfx"
)

type HazardState int

const (
	HazardTelegraph HazardState = iota
	HazardFalling
	HazardImpact
	HazardFinished
)

// FallingDiscoBall represents a giant collapsing disco ball.
type FallingDiscoBall struct {
	X, Y              float64
	Radius            float64
	TelegraphDuration float64
	Elapsed           float64
	DropY             float64
	Rotation          float64
	State             HazardState
	DealtDamage       bool
}

func NewFallingDiscoBall(x, y, radius, warningDuration float64) *FallingDiscoBall {
	return &FallingDiscoBall{
		X:                 x,
		Y:                 y,
		Radius:            radius,
		TelegraphDuration: warningDuration,
		DropY:             y - 280.0,
		State:             HazardTelegraph,
	}
}

func (db *FallingDiscoBall) Update(dt float64, p *Player, ps *gfx.ParticleSystem, ae *audio.AudioEngine) {
	db.Elapsed += dt
	db.Rotation += dt * 3.0

	switch db.State {
	case HazardTelegraph:
		if db.Elapsed >= db.TelegraphDuration {
			db.State = HazardFalling
			if ae != nil {
				ae.PlaySFXWhistle()
			}
		}
	case HazardFalling:
		// Accelerate drop toward target Y
		dropSpeed := 580.0
		db.DropY += dropSpeed * dt
		if db.DropY >= db.Y {
			db.DropY = db.Y
			db.State = HazardImpact
		}
	case HazardImpact:
		if !db.DealtDamage {
			db.DealtDamage = true
			if ae != nil {
				ae.PlaySFXCrash()
			}

			// Check damage to player
			dist := math.Hypot(p.X-db.X, p.Y-db.Y)
			if dist < db.Radius+p.HitboxRadius {
				p.TakeDamage(1, ps, ae)
				// Knockback
				angle := math.Atan2(p.Y-db.Y, p.X-db.X)
				p.VX += math.Cos(angle) * 320.0
				p.VY += math.Sin(angle) * 320.0
			}

			// Explode into 30 mirror shards
			if ps != nil {
				rnd := rand.New(rand.NewSource(int64(db.X*1000 + db.Y)))
				for i := 0; i < 30; i++ {
					ang := rnd.Float64() * math.Pi * 2.0
					speed := 80.0 + rnd.Float64()*180.0
					size := 2.5 + rnd.Float64()*3.5
					col := color.RGBA{R: 210, G: 235, B: 255, A: 255}
					if i%3 == 0 {
						col = color.RGBA{R: 255, G: 255, B: 255, A: 255}
					}
					ps.Emit(db.X, db.Y, math.Cos(ang)*speed, math.Sin(ang)*speed-60.0, 0.65, size, 1.0, col, gfx.ParticleMirrorShard)
				}
			}
		}
		// Brief impact persistence then disappear
		if db.Elapsed >= db.TelegraphDuration+0.5 {
			db.State = HazardFinished
		}
	}
}

func (db *FallingDiscoBall) Draw(screen *ebiten.Image) {
	switch db.State {
	case HazardTelegraph:
		prog := db.Elapsed / db.TelegraphDuration
		gfx.DrawTelegraphCircle(screen, db.X, db.Y, db.Radius, prog)
	case HazardFalling:
		gfx.DrawTelegraphCircle(screen, db.X, db.Y, db.Radius, 1.0)
		gfx.DrawDiscoBall(screen, db.X, db.DropY, db.Radius*0.75, db.Rotation)
	case HazardImpact:
		gfx.DrawDiscoBall(screen, db.X, db.Y, db.Radius*0.75, db.Rotation)
	}
}

// FallingTruss represents a steel lighting truss dropping from the ceiling.
type FallingTruss struct {
	X, Y              float64
	W, H              float64
	TelegraphDuration float64
	Elapsed           float64
	DropY             float64
	State             HazardState
	DealtDamage       bool
	SolidTimer        float64
}

func NewFallingTruss(x, y, w, h, warningDuration float64) *FallingTruss {
	return &FallingTruss{
		X:                 x,
		Y:                 y,
		W:                 w,
		H:                 h,
		TelegraphDuration: warningDuration,
		DropY:             y - 260.0,
		State:             HazardTelegraph,
		SolidTimer:        6.0, // Stays on floor for 6s
	}
}

func (ft *FallingTruss) Update(dt float64, p *Player, ps *gfx.ParticleSystem, ae *audio.AudioEngine) {
	ft.Elapsed += dt

	switch ft.State {
	case HazardTelegraph:
		if ft.Elapsed >= ft.TelegraphDuration {
			ft.State = HazardFalling
			if ae != nil {
				ae.PlaySFXWhistle()
			}
		}
	case HazardFalling:
		ft.DropY += 600.0 * dt
		if ft.DropY >= ft.Y {
			ft.DropY = ft.Y
			ft.State = HazardImpact
		}
	case HazardImpact:
		if !ft.DealtDamage {
			ft.DealtDamage = true
			if ae != nil {
				ae.PlaySFXCrash()
			}

			// Check damage
			if p.X >= ft.X && p.X <= ft.X+ft.W && p.Y >= ft.Y && p.Y <= ft.Y+ft.H {
				p.TakeDamage(1, ps, ae)
			}

			// Emit sparks along the truss length
			if ps != nil {
				for s := ft.X; s < ft.X+ft.W; s += 16.0 {
					ps.Emit(s, ft.Y+ft.H/2, (rand.Float64()-0.5)*80.0, -50.0-rand.Float64()*70.0, 0.4, 3.0, 1.0, color.RGBA{255, 220, 0, 255}, gfx.ParticleSpark)
				}
			}
		}

		// Push player out if walking into solid truss
		if p.X > ft.X-p.Width/2 && p.X < ft.X+ft.W+p.Width/2 && p.Y > ft.Y-p.Height/2 && p.Y < ft.Y+ft.H+p.Height/2 {
			if p.Y < ft.Y {
				p.Y = ft.Y - p.Height/2
			} else if p.Y > ft.Y+ft.H {
				p.Y = ft.Y + ft.H + p.Height/2
			}
		}

		ft.SolidTimer -= dt
		if ft.SolidTimer <= 0 {
			ft.State = HazardFinished
		}
	}
}

func (ft *FallingTruss) Draw(screen *ebiten.Image) {
	switch ft.State {
	case HazardTelegraph:
		prog := ft.Elapsed / ft.TelegraphDuration
		gfx.DrawTelegraphCircle(screen, ft.X+ft.W/2, ft.Y+ft.H/2, ft.W/2, prog)
	case HazardFalling:
		gfx.DrawTruss(screen, ft.X, ft.DropY, ft.W, ft.H)
	case HazardImpact:
		gfx.DrawTruss(screen, ft.X, ft.Y, ft.W, ft.H)
	}
}

// DrinkPuddle represents an alcohol spill on the floor.
type DrinkPuddle struct {
	X, Y   float64
	Radius float64
	Color  color.RGBA
}

func NewDrinkPuddle(x, y, radius float64, col color.RGBA) *DrinkPuddle {
	return &DrinkPuddle{
		X:      x,
		Y:      y,
		Radius: radius,
		Color:  col,
	}
}

func (dp *DrinkPuddle) Update(p *Player, ae *audio.AudioEngine) {
	dist := math.Hypot(p.X-dp.X, p.Y-dp.Y)
	if dist < dp.Radius {
		p.ApplySlip(0.75, ae)
	}
}

func (dp *DrinkPuddle) Draw(screen *ebiten.Image) {
	gfx.DrawDrinkPuddle(screen, dp.X, dp.Y, dp.Radius, dp.Color)
}

// PanickedClubber represents a wandering crowd patron.
type PanickedClubber struct {
	X, Y        float64
	VX, VY      float64
	OutfitColor color.RGBA
	AnimTime    float64
	ChangeTimer float64
}

func NewPanickedClubber(x, y float64, col color.RGBA) *PanickedClubber {
	return &PanickedClubber{
		X:           x,
		Y:           y,
		OutfitColor: col,
	}
}

func (pc *PanickedClubber) Update(dt float64, boundsX, boundsY, boundsW, boundsH float64, p *Player) {
	pc.AnimTime += dt
	pc.ChangeTimer -= dt

	if pc.ChangeTimer <= 0 {
		pc.ChangeTimer = 0.8 + rand.Float64()*1.2
		ang := rand.Float64() * math.Pi * 2.0
		speed := 50.0 + rand.Float64()*50.0
		pc.VX = math.Cos(ang) * speed
		pc.VY = math.Sin(ang) * speed
	}

	pc.X += pc.VX * dt
	pc.Y += pc.VY * dt

	// Bounds bouncing
	if pc.X < boundsX+20 || pc.X > boundsX+boundsW-20 {
		pc.VX = -pc.VX
	}
	if pc.Y < boundsY+20 || pc.Y > boundsY+boundsH-20 {
		pc.VY = -pc.VY
	}

	// Soft bump with player
	dist := math.Hypot(p.X-pc.X, p.Y-pc.Y)
	if dist < 16.0 && dist > 0.01 {
		overlap := 16.0 - dist
		nx := (p.X - pc.X) / dist
		ny := (p.Y - pc.Y) / dist
		p.X += nx * overlap * 0.5
		p.Y += ny * overlap * 0.5
		pc.X -= nx * overlap * 0.5
		pc.Y -= ny * overlap * 0.5
	}
}

func (pc *PanickedClubber) Draw(screen *ebiten.Image) {
	gfx.DrawPanickedClubber(screen, pc.X, pc.Y, pc.AnimTime, pc.OutfitColor)
}

// ExitDoor represents the illuminated emergency exit.
type ExitDoor struct {
	X, Y, W, H float64
	AnimTime   float64
}

func NewExitDoor(x, y, w, h float64) *ExitDoor {
	return &ExitDoor{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

func (ed *ExitDoor) Update(dt float64) {
	ed.AnimTime += dt
}

func (ed *ExitDoor) IsPlayerInside(p *Player) bool {
	return p.X >= ed.X && p.X <= ed.X+ed.W && p.Y >= ed.Y && p.Y <= ed.Y+ed.H
}

func (ed *ExitDoor) Draw(screen *ebiten.Image) {
	gfx.DrawExitPortal(screen, ed.X, ed.Y, ed.W, ed.H, ed.AnimTime)
}
