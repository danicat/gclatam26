package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"panic-invaders/internal/assets"
)

type Block struct {
	Alive bool
}

type Barrier struct {
	X      float64
	Y      float64
	Label  string
	Rows   int
	Cols   int
	BlockW float64
	BlockH float64
	Grid   [][]Block
	img    *ebiten.Image
}

func NewBarrier(x, y float64, label string) *Barrier {
	rows := 5
	cols := 10
	bw := 4.0
	bh := 3.0

	grid := make([][]Block, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]Block, cols)
		for c := 0; c < cols; c++ {
			// Cut an arch at the bottom center
			if r >= 3 && (c >= 3 && c <= 6) {
				grid[r][c] = Block{Alive: false}
			} else {
				grid[r][c] = Block{Alive: true}
			}
		}
	}

	b := &Barrier{
		X:      x,
		Y:      y,
		Label:  label,
		Rows:   rows,
		Cols:   cols,
		BlockW: bw,
		BlockH: bh,
		Grid:   grid,
	}
	b.renderImage()
	return b
}

func (b *Barrier) renderImage() {
	w := int(float64(b.Cols) * b.BlockW)
	h := int(float64(b.Rows) * b.BlockH)
	if b.img == nil {
		b.img = ebiten.NewImage(w, h)
	} else {
		b.img.Clear()
	}

	blockImg := ebiten.NewImage(int(b.BlockW), int(b.BlockH))
	blockImg.Fill(assets.ColorDeferGreen)

	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Grid[r][c].Alive {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(c)*b.BlockW, float64(r)*b.BlockH)
				b.img.DrawImage(blockImg, op)
			}
		}
	}
}

func (b *Barrier) CheckBulletCollision(bullet *Bullet) bool {
	if !bullet.Active {
		return false
	}
	bx := bullet.X
	by := bullet.Y
	bw := bullet.Width
	bh := bullet.Height

	barrierW := float64(b.Cols) * b.BlockW
	barrierH := float64(b.Rows) * b.BlockH

	// AABB bounds check
	if bx+bw < b.X || bx > b.X+barrierW || by+bh < b.Y || by > b.Y+barrierH {
		return false
	}

	// Check individual blocks
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if !b.Grid[r][c].Alive {
				continue
			}
			blockX := b.X + float64(c)*b.BlockW
			blockY := b.Y + float64(r)*b.BlockH

			if bx+bw >= blockX && bx <= blockX+b.BlockW &&
				by+bh >= blockY && by <= blockY+b.BlockH {
				b.Grid[r][c].Alive = false
				bullet.Active = false
				b.renderImage()
				return true
			}
		}
	}
	return false
}

func (b *Barrier) Draw(screen *ebiten.Image) {
	if b.img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.X, b.Y)
		screen.DrawImage(b.img, op)
	}
}
