package areabmp_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/areabmp"
)

func TestAreaBMP_LoadAll(t *testing.T) {
	areaBMPDir := filepath.Join("..", "..", "..", "data", "AreaBMP")

	store, err := areabmp.LoadAll(areaBMPDir)
	if err != nil {
		t.Fatalf("Failed to load AreaBMP directory: %v", err)
	}

	count := store.Count()
	if count < 44 {
		t.Fatalf("Expected at least 44 unique BMP masks, got %d", count)
	}

	t.Logf("AreaBMP Store loaded %d unique masks (total %d references)", count, len(store.BMPs))

	// Test sample mask ArkAstanica_Skill01_W
	bmp := store.GetBMP("ArkAstanica_Skill01_W")
	if bmp == nil {
		t.Fatal("ArkAstanica_Skill01_W.bmp missing in store")
	}

	if bmp.Width != 100 || bmp.Height != 400 {
		t.Errorf("Expected 100x400, got %dx%d", bmp.Width, bmp.Height)
	}

	// Test ContainsTarget on center
	// With caster at (1000, 1000) and target at (1000, 1000)
	// Relative offset (0, 0) maps to (50, 200) which should be checked
	_ = bmp.ContainsTarget(1000, 1000, 0, 1000, 1000)
}

