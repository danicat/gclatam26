package levels

import (
	"image/color"
	"math/rand"

	"panic-at-the-disco/internal/entities"
	"panic-at-the-disco/internal/gfx"
)

type ZoneID int

const (
	ZoneDanceFloor ZoneID = 1
	ZoneVIPLounge  ZoneID = 2
	ZoneBackstage  ZoneID = 3
)

type LevelConfig struct {
	ID             ZoneID
	Name           string
	Duration       float64 // Seconds before roof collapses
	DiscoFloorCols int
	DiscoFloorRows int
	BPM            float64
	HazardInterval float64
	PlayerStartX   float64
	PlayerStartY   float64
	ExitX          float64
	ExitY          float64
	ExitW          float64
	ExitH          float64
}

// GetLevelConfig returns the configuration for a given zone.
func GetLevelConfig(id ZoneID) LevelConfig {
	switch id {
	case ZoneDanceFloor:
		return LevelConfig{
			ID:             ZoneDanceFloor,
			Name:           "ZONE 1: THE MAIN DANCE FLOOR",
			Duration:       45.0,
			DiscoFloorCols: 14,
			DiscoFloorRows: 7,
			BPM:            120.0,
			HazardInterval: 1.8,
			PlayerStartX:   70.0,
			PlayerStartY:   180.0,
			ExitX:          570.0,
			ExitY:          50.0,
			ExitW:          45.0,
			ExitH:          45.0,
		}
	case ZoneVIPLounge:
		return LevelConfig{
			ID:             ZoneVIPLounge,
			Name:           "ZONE 2: THE VIP LOUNGE & BAR",
			Duration:       40.0,
			DiscoFloorCols: 12,
			DiscoFloorRows: 7,
			BPM:            126.0,
			HazardInterval: 1.5,
			PlayerStartX:   60.0,
			PlayerStartY:   260.0,
			ExitX:          570.0,
			ExitY:          160.0,
			ExitW:          45.0,
			ExitH:          45.0,
		}
	case ZoneBackstage:
		return LevelConfig{
			ID:             ZoneBackstage,
			Name:           "ZONE 3: BACKSTAGE & FIRE ALLEY",
			Duration:       35.0,
			DiscoFloorCols: 10,
			DiscoFloorRows: 6,
			BPM:            132.0,
			HazardInterval: 1.2,
			PlayerStartX:   60.0,
			PlayerStartY:   180.0,
			ExitX:          295.0,
			ExitY:          45.0,
			ExitW:          50.0,
			ExitH:          45.0,
		}
	default:
		return GetLevelConfig(ZoneDanceFloor)
	}
}

// SetupZoneEntities populates the level with stage-specific static/dynamic props.
func SetupZoneEntities(cfg LevelConfig) (*gfx.DiscoFloor, *entities.ExitDoor, []*entities.DrinkPuddle, []*entities.PanickedClubber) {
	// Disco floor
	df := gfx.NewDiscoFloor(35.0, 45.0, 570.0, 275.0, cfg.DiscoFloorCols, cfg.DiscoFloorRows, cfg.BPM)

	// Exit Door
	exit := entities.NewExitDoor(cfg.ExitX, cfg.ExitY, cfg.ExitW, cfg.ExitH)

	// Puddles (mainly in VIP lounge and backstage)
	var puddles []*entities.DrinkPuddle
	rnd := rand.New(rand.NewSource(int64(cfg.ID * 999)))

	puddleCount := 0
	if cfg.ID == ZoneVIPLounge {
		puddleCount = 6
	} else if cfg.ID == ZoneBackstage {
		puddleCount = 3
	}

	puddleColors := []color.RGBA{
		{R: 255, G: 165, B: 0, A: 160}, // Amber whiskey
		{R: 220, G: 20, B: 60, A: 160},  // Cosmopolitan red
		{R: 50, G: 205, B: 50, A: 160},  // Midori green
	}

	for i := 0; i < puddleCount; i++ {
		px := 120.0 + rnd.Float64()*380.0
		py := 80.0 + rnd.Float64()*190.0
		rad := 16.0 + rnd.Float64()*14.0
		col := puddleColors[i%len(puddleColors)]
		puddles = append(puddles, entities.NewDrinkPuddle(px, py, rad, col))
	}

	// Panicked Clubbers
	var clubbers []*entities.PanickedClubber
	clubberCount := 4
	if cfg.ID == ZoneDanceFloor {
		clubberCount = 6
	} else if cfg.ID == ZoneVIPLounge {
		clubberCount = 4
	} else {
		clubberCount = 2
	}

	outfitColors := []color.RGBA{
		{R: 255, G: 215, B: 0, A: 255},  // Gold jumpsuit
		{R: 0, G: 255, B: 255, A: 255},  // Cyan leisure suit
		{R: 255, G: 105, B: 180, A: 255}, // Hot pink sequins
		{R: 147, G: 112, B: 219, A: 255}, // Purple velvet
	}

	for i := 0; i < clubberCount; i++ {
		cx := 100.0 + rnd.Float64()*420.0
		cy := 70.0 + rnd.Float64()*200.0
		col := outfitColors[i%len(outfitColors)]
		clubbers = append(clubbers, entities.NewPanickedClubber(cx, cy, col))
	}

	return df, exit, puddles, clubbers
}
