package game

import (
	"math"
)

type TileType int

const (
	TileFloor TileType = iota
	TileWall
	TileHole
	TileHoleFilled
	TileClock
	TileArtifact
)

type Boulder struct {
	X, Y             int
	VisualX, VisualY float64
}

type LevelDef struct {
	Name              string
	SubTitle          string
	Width             int
	Height            int
	MaxTurns          int
	ClockRecoverTurns int
	PlayerStartX      int
	PlayerStartY      int
	Layout            []string
}

type ActionEvent int

const (
	EventNone ActionEvent = iota
	EventStep
	EventPush
	EventHoleFilled
	EventRecover
	EventWin
	EventFaint
)

type LevelState struct {
	Def               LevelDef
	Width             int
	Height            int
	MaxTurns          int
	ClockRecoverTurns int
	TurnsLeft         int

	PlayerX       int
	PlayerY       int
	PlayerVisualX float64
	PlayerVisualY float64

	Tiles     [][]TileType
	Boulders  []Boulder
	Fainted   bool
	Cleared   bool
	StepCount int

	LastEvent      ActionEvent
	EventX, EventY float64
}

func NewLevelState(def LevelDef) *LevelState {
	ls := &LevelState{
		Def:               def,
		Width:             def.Width,
		Height:            def.Height,
		MaxTurns:          def.MaxTurns,
		ClockRecoverTurns: def.ClockRecoverTurns,
		TurnsLeft:         def.MaxTurns,
		PlayerX:           def.PlayerStartX,
		PlayerY:           def.PlayerStartY,
		PlayerVisualX:     float64(def.PlayerStartX * TileSize),
		PlayerVisualY:     float64(def.PlayerStartY * TileSize),
		Boulders:          make([]Boulder, 0),
		Tiles:             make([][]TileType, def.Height),
	}

	for y := 0; y < def.Height; y++ {
		ls.Tiles[y] = make([]TileType, def.Width)
		rowStr := ""
		if y < len(def.Layout) {
			rowStr = def.Layout[y]
		}
		for x := 0; x < def.Width; x++ {
			var ch byte = ' '
			if x < len(rowStr) {
				ch = rowStr[x]
			}
			switch ch {
			case '#':
				ls.Tiles[y][x] = TileWall
			case 'O':
				ls.Tiles[y][x] = TileHole
			case 'B':
				ls.Tiles[y][x] = TileFloor
				ls.Boulders = append(ls.Boulders, Boulder{
					X:       x,
					Y:       y,
					VisualX: float64(x * TileSize),
					VisualY: float64(y * TileSize),
				})
			case 'C':
				ls.Tiles[y][x] = TileClock
			case 'A':
				ls.Tiles[y][x] = TileArtifact
			case 'P':
				ls.Tiles[y][x] = TileFloor
				ls.PlayerX = x
				ls.PlayerY = y
				ls.PlayerVisualX = float64(x * TileSize)
				ls.PlayerVisualY = float64(y * TileSize)
			default:
				ls.Tiles[y][x] = TileFloor
			}
		}
	}

	return ls
}

func (ls *LevelState) UpdateVisuals(dt float64) {
	// Smooth lerp for Player movement (ease-out)
	targetPX := float64(ls.PlayerX * TileSize)
	targetPY := float64(ls.PlayerY * TileSize)
	lerpSpeed := math.Min(1.0, dt*24.0)

	ls.PlayerVisualX += (targetPX - ls.PlayerVisualX) * lerpSpeed
	ls.PlayerVisualY += (targetPY - ls.PlayerVisualY) * lerpSpeed

	// Smooth lerp for Boulders
	for i := range ls.Boulders {
		b := &ls.Boulders[i]
		targetBX := float64(b.X * TileSize)
		targetBY := float64(b.Y * TileSize)
		b.VisualX += (targetBX - b.VisualX) * lerpSpeed
		b.VisualY += (targetBY - b.VisualY) * lerpSpeed
	}
}

func (ls *LevelState) PanicPercent() float64 {
	if ls.MaxTurns <= 0 {
		return 0.0
	}
	used := float64(ls.MaxTurns - ls.TurnsLeft)
	p := used / float64(ls.MaxTurns)
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	return p
}

func (ls *LevelState) findBoulder(x, y int) int {
	for i, b := range ls.Boulders {
		if b.X == x && b.Y == y {
			return i
		}
	}
	return -1
}

func (ls *LevelState) Move(dx, dy int) bool {
	if ls.Fainted || ls.Cleared {
		return false
	}

	nx := ls.PlayerX + dx
	ny := ls.PlayerY + dy

	if nx < 0 || nx >= ls.Width || ny < 0 || ny >= ls.Height {
		return false
	}

	ls.LastEvent = EventNone

	// 1. Check if pushing a Boulder
	bIdx := ls.findBoulder(nx, ny)
	if bIdx >= 0 {
		bx := nx + dx
		by := ny + dy
		if bx < 0 || bx >= ls.Width || by < 0 || by >= ls.Height {
			return false
		}
		if ls.Tiles[by][bx] == TileWall || ls.Tiles[by][bx] == TileArtifact || ls.Tiles[by][bx] == TileClock || ls.findBoulder(bx, by) >= 0 {
			return false
		}

		if ls.Tiles[by][bx] == TileHole {
			ls.Tiles[by][bx] = TileHoleFilled
			ls.Boulders = append(ls.Boulders[:bIdx], ls.Boulders[bIdx+1:]...)
			playSound(sndFillHole, 0.9)
			ls.LastEvent = EventHoleFilled
			ls.EventX = float64(bx * TileSize)
			ls.EventY = float64(by * TileSize)
		} else {
			ls.Boulders[bIdx].X = bx
			ls.Boulders[bIdx].Y = by
			playSound(sndPush, 0.7)
			ls.LastEvent = EventPush
			ls.EventX = float64(bx * TileSize)
			ls.EventY = float64(by * TileSize)
		}

		ls.PlayerX = nx
		ls.PlayerY = ny
		ls.consumeTurn()
		return true
	}

	targetTile := ls.Tiles[ny][nx]

	// 2. Wall or Open Hole -> Blocked
	if targetTile == TileWall || targetTile == TileHole {
		return false
	}

	// 3. Artifact Goal -> Level Complete!
	if targetTile == TileArtifact {
		ls.PlayerX = nx
		ls.PlayerY = ny
		ls.Cleared = true
		playSound(sndWin, 0.8)
		ls.LastEvent = EventWin
		ls.EventX = float64(nx * TileSize)
		ls.EventY = float64(ny * TileSize)
		return true
	}

	// 4. Clock / Recovery pickup
	if targetTile == TileClock {
		ls.PlayerX = nx
		ls.PlayerY = ny
		ls.Tiles[ny][nx] = TileFloor
		ls.TurnsLeft += ls.ClockRecoverTurns
		if ls.TurnsLeft > ls.MaxTurns {
			ls.TurnsLeft = ls.MaxTurns
		}
		playSound(sndRecover, 0.8)
		ls.LastEvent = EventRecover
		ls.EventX = float64(nx * TileSize)
		ls.EventY = float64(ny * TileSize)
		ls.consumeTurn()
		return true
	}

	// 5. Normal floor or filled hole
	ls.PlayerX = nx
	ls.PlayerY = ny
	playSound(sndStep, 0.5)
	ls.LastEvent = EventStep
	ls.consumeTurn()
	return true
}

func (ls *LevelState) consumeTurn() {
	ls.StepCount++
	ls.TurnsLeft--
	if ls.TurnsLeft <= 0 {
		ls.TurnsLeft = 0
		ls.Fainted = true
		ls.LastEvent = EventFaint
		playSound(sndFaint, 0.8)
	} else {
		playSound(sndTick, 0.4)
		if ls.PanicPercent() >= 0.80 {
			playSound(sndHeartbeat, 0.6)
		}
	}
}
