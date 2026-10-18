package field_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/field"
)

func TestFieldStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "Map")
	store := field.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	if count := store.FieldCount(); count != 153 {
		t.Fatalf("Expected 153 fields, got %d", count)
	}

	// 1. Verify Roumen (Serial 0)
	rou, ok := store.GetFieldBySerial(0)
	if !ok {
		t.Fatalf("Serial 0 not found")
	}
	if rou.MapID != "Rou" || rou.Name != "Roumen" || rou.Width != 256 || rou.Height != 256 {
		t.Errorf("Unexpected Roumen config: %+v", rou)
	}
	if rou.Respawn.City != "Rou" || len(rou.Respawn.Spots) != 4 || rou.Respawn.Spots[0].X != 4134 {
		t.Errorf("Unexpected Roumen respawn: %+v", rou.Respawn)
	}
	if !rou.Permissions.CanTrade || !rou.Permissions.CanMiniHouse {
		t.Errorf("Expected true for trade and minihouse in Roumen")
	}

	// 2. Lookup by MapID case-insensitive
	eld, ok := store.GetField("eld")
	if !ok {
		t.Fatalf("Elderine (eld) not found")
	}
	if eld.Name != "Elderine" || eld.Width != 512 {
		t.Errorf("Unexpected Elderine: %+v", eld)
	}

	// 3. Multi-layer map lookup
	rouLayers := store.GetAllFieldsByMapID("Rou")
	if len(rouLayers) != 3 {
		t.Errorf("Expected 3 layers for Rou, got %d", len(rouLayers))
	}

	// 4. Instance Dungeons by Zone
	zone3Dungeons := store.GetDungeonsByZone(3)
	if len(zone3Dungeons) == 0 {
		t.Errorf("Expected instance dungeons for zone 3")
	}
	// Check UrgDragon / GraveYard
	foundGrave := false
	for _, d := range zone3Dungeons {
		if d.MapID == "UrgDragon" && d.ScriptName == "ID/GraveYard/GraveYard" {
			foundGrave = true
			break
		}
	}
	if !foundGrave {
		t.Errorf("Expected UrgDragon / GraveYard dungeon in zone 3")
	}

	// 5. Non-existent lookups
	if _, ok := store.GetFieldBySerial(9999); ok {
		t.Errorf("Expected false for non-existent serial")
	}
	if _, ok := store.GetField("NonExistentMapXYZ"); ok {
		t.Errorf("Expected false for non-existent map ID")
	}
}
