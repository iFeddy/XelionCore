package classparam_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
)

func TestClassParamStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	store := classparam.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	if count := store.ClassCount(); count != 27 {
		t.Fatalf("Expected 27 classes loaded, got %d", count)
	}

	// Verify Fighter Level 1 stats
	fighterLvl1, ok := store.GetParam("Fighter", 1)
	if !ok {
		t.Fatalf("Fighter level 1 not found")
	}
	if fighterLvl1.Strength != 6 || fighterLvl1.Constitution != 5 || fighterLvl1.MaxHP != 52 {
		t.Errorf("Unexpected Fighter Lvl 1 stats: %+v", fighterLvl1)
	}

	// Verify Fighter Level 150 stats
	fighterLvl150, ok := store.GetParam("fighter", 150)
	if !ok {
		t.Fatalf("Fighter level 150 not found")
	}
	if fighterLvl150.Level != 150 {
		t.Errorf("Expected level 150, got %d", fighterLvl150.Level)
	}

	// Verify non-existent class or level bounds
	if _, ok := store.GetParam("UnknownClass", 1); ok {
		t.Errorf("Expected false for non-existent class")
	}
	if _, ok := store.GetParam("Fighter", 0); ok {
		t.Errorf("Expected false for level 0")
	}
	if _, ok := store.GetParam("Fighter", 151); ok {
		t.Errorf("Expected false for level 151")
	}

	// Verify Mage Level 1 stats
	mageLvl1, ok := store.GetParam("Mage", 1)
	if !ok {
		t.Fatalf("Mage level 1 not found")
	}
	if mageLvl1.Intelligence != 10 || mageLvl1.Strength != 1 {
		t.Errorf("Unexpected Mage Lvl 1 stats: %+v", mageLvl1)
	}
}
