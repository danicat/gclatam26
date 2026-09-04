package app

import "testing"

func TestEmbeddedSpritesAreLoadable(t *testing.T) {
	t.Parallel()

	assets, err := loadEmbeddedAssets()
	if err != nil {
		t.Fatalf("loadEmbeddedAssets() error = %v", err)
	}
	if assets.gopher == nil || assets.bug == nil {
		t.Fatal("embedded sprites must both be loaded")
	}
}
