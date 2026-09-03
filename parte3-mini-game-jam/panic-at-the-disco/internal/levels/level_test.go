package levels

import (
	"testing"
)

func TestZoneConfigs(t *testing.T) {
	zones := []ZoneID{ZoneDanceFloor, ZoneVIPLounge, ZoneBackstage}

	for _, z := range zones {
		cfg := GetLevelConfig(z)
		if cfg.Duration <= 0 {
			t.Fatalf("Zone %d has invalid duration %f", z, cfg.Duration)
		}
		if cfg.ExitW <= 0 || cfg.ExitH <= 0 {
			t.Fatalf("Zone %d has invalid exit dimensions (%f, %f)", z, cfg.ExitW, cfg.ExitH)
		}
		if cfg.BPM <= 0 {
			t.Fatalf("Zone %d has invalid BPM %f", z, cfg.BPM)
		}

		floor, exit, puddles, clubbers := SetupZoneEntities(cfg)
		if floor == nil || exit == nil {
			t.Fatalf("Zone %d setup returned nil floor or exit", z)
		}
		_ = puddles
		if len(clubbers) == 0 {
			t.Fatalf("Zone %d expected non-zero clubbers", z)
		}
	}
}
