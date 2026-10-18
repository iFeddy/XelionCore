package abstate_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
)

func TestLoadFromFile(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "AbState")

	db, err := abstate.LoadFromFile(dirPath)
	if err != nil {
		t.Fatalf("LoadFromFile from dir failed: %v", err)
	}

	if len(db.States) != 10 {
		t.Errorf("Expected 10 states, got %d", len(db.States))
	}

	// Test StaHardSkin (ID 1)
	skin, ok := db.GetByName("StaHardSkin")
	if !ok {
		t.Fatalf("StaHardSkin not found")
	}
	if skin.ID != 1 {
		t.Errorf("Expected ID 1, got %d", skin.ID)
	}

	// Verify Level 1: 303000 ms, min_def: 4, max_def: 5
	lv1, err := skin.GetLevelConfig(1)
	if err != nil {
		t.Fatalf("GetLevelConfig(1) failed: %v", err)
	}
	if lv1.DurationMs != 303000 {
		t.Errorf("Expected 303000 ms, got %d", lv1.DurationMs)
	}
	if lv1.MinDef != 4 || lv1.MaxDef != 5 {
		t.Errorf("Expected Def 4-5, got %d-%d", lv1.MinDef, lv1.MaxDef)
	}

	// Test GetByID(13) -> StaShield
	shield, ok := db.GetByID(13)
	if !ok {
		t.Fatalf("StaShield (ID 13) not found")
	}
	if shield.Name != "StaShield" {
		t.Errorf("Expected StaShield, got %s", shield.Name)
	}

	// Verify Level 2: 306000 ms, min_shield: 5, max_shield: 5
	lv2, err := shield.GetLevelConfig(2)
	if err != nil {
		t.Fatalf("GetLevelConfig(2) failed: %v", err)
	}
	if lv2.DurationMs != 306000 {
		t.Errorf("Expected 306000 ms, got %d", lv2.DurationMs)
	}
	if lv2.MinShield != 5 || lv2.MaxShield != 5 {
		t.Errorf("Expected Shield 5-5, got %d-%d", lv2.MinShield, lv2.MaxShield)
	}

	// Test StaCleHold (ID 11) - Hold Cleric
	hold, ok := db.GetByID(11)
	if !ok {
		t.Fatalf("StaCleHold (ID 11) not found")
	}
	lvHold1, err := hold.GetLevelConfig(1)
	if err != nil {
		t.Fatalf("GetLevelConfig(1) failed: %v", err)
	}
	if lvHold1.DurationMs != 2010 {
		t.Errorf("Expected 2010 ms, got %d", lvHold1.DurationMs)
	}
}
