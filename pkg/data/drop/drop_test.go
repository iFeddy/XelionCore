package drop_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/drop"
)

func TestDropStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "Drop")
	store := drop.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	if count := store.MobCount(); count != 1485 {
		t.Fatalf("Expected 1485 mobs, got %d", count)
	}

	if gCount := store.GroupCount(); gCount != 959 {
		t.Fatalf("Expected 959 groups, got %d", gCount)
	}

	// 1. Verify SlowSlime
	slimeDrop, ok := store.GetMobDrop("SlowSlime")
	if !ok {
		t.Fatalf("SlowSlime drop table not found")
	}
	if slimeDrop.MinLevel != 1 || slimeDrop.MaxLevel != 16 {
		t.Errorf("Unexpected SlowSlime level bounds: min=%d, max=%d", slimeDrop.MinLevel, slimeDrop.MaxLevel)
	}
	if len(slimeDrop.Drops) == 0 {
		t.Errorf("Expected non-empty drops for SlowSlime")
	}

	// Check Potion01 in SlowSlime drops
	foundPotion := false
	for _, d := range slimeDrop.Drops {
		if d.Item == "Potion01" && d.Rate == 47500 {
			foundPotion = true
			break
		}
	}
	if !foundPotion {
		t.Errorf("Expected Potion01 with rate 47500 in SlowSlime drops")
	}

	// 2. Verify Honeying
	honeyingDrop, ok := store.GetMobDrop("honeying")
	if !ok {
		t.Fatalf("Honeying drop table not found (case-insensitive)")
	}
	if len(honeyingDrop.Drops) == 0 {
		t.Errorf("Expected non-empty drops for Honeying")
	}

	// 3. Verify ItemDropGroup Weapon1
	w1, ok := store.GetGroup("Weapon1")
	if !ok {
		t.Fatalf("Group Weapon1 not found")
	}
	if len(w1) == 0 || w1[0].ItemID != "Weapon1" {
		t.Errorf("Unexpected group Weapon1 entry: %+v", w1)
	}

	// 4. Non-existent mob and group
	if _, ok := store.GetMobDrop("NonExistentMob123"); ok {
		t.Errorf("Expected false for non-existent mob")
	}
	if _, ok := store.GetGroup("NonExistentGroup123"); ok {
		t.Errorf("Expected false for non-existent group")
	}
}
