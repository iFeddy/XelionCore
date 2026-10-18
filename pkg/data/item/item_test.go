package item_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/item"
)

func TestItemStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "Item")
	store := item.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	// 1. Verify ItemOptions
	opts := store.Options()
	if len(opts) != 18 {
		t.Fatalf("Expected 18 options, got %d", len(opts))
	}
	if opts[0].Degree != 1 || opts[0].Type != "Number" || len(opts[0].Weights) != 9 {
		t.Errorf("Unexpected first option: %+v", opts[0])
	}

	// 2. Verify ItemUseFunction
	fn, ok := store.GetUseFunction("SilverWings")
	if !ok {
		t.Fatalf("SilverWings use function not found")
	}
	if fn.UseFunction != "UseSilverWings" || fn.Broadcast {
		t.Errorf("Unexpected SilverWings function: %+v", fn)
	}

	fnKQ, ok := store.GetUseFunction("kq_soulcrystal")
	if !ok {
		t.Fatalf("KQ_SoulCrystal not found (case-insensitive)")
	}
	if fnKQ.UseFunction != "BossMobDetect" || !fnKQ.Broadcast {
		t.Errorf("Unexpected KQ_SoulCrystal function: %+v", fnKQ)
	}

	// 3. Verify RandomOptionTable
	randW1, ok := store.GetRandomOption("RandomNamedW01")
	if !ok {
		t.Fatalf("RandomNamedW01 not found")
	}
	if randW1.Unidentified || randW1.MinOptions != 1 || randW1.MaxOptions != 3 {
		t.Errorf("Unexpected RandomNamedW01 options: %+v", randW1)
	}
	if randW1.STR.Min != 1 || randW1.STR.Max != 2 {
		t.Errorf("Unexpected RandomNamedW01 STR range: %+v", randW1.STR)
	}

	// 4. Verify TreasureReward
	tbag, ok := store.GetContainer("Tbag06")
	if !ok {
		t.Fatalf("Container Tbag06 not found")
	}
	if tbag.DummyID != "DTbag06" || len(tbag.Rewards) != 9 {
		t.Errorf("Unexpected Tbag06 container: %+v", tbag)
	}
	if tbag.Rewards[0].ItemID != "GrimjawMask00" || tbag.Rewards[0].Rate != 50000 {
		t.Errorf("Unexpected first reward in Tbag06: %+v", tbag.Rewards[0])
	}

	// 5. Non-existent lookups
	if _, ok := store.GetUseFunction("NonExistentItem123"); ok {
		t.Errorf("Expected false for non-existent item function")
	}
	if _, ok := store.GetRandomOption("NonExistentDrop123"); ok {
		t.Errorf("Expected false for non-existent random option")
	}
	if _, ok := store.GetContainer("NonExistentContainer123"); ok {
		t.Errorf("Expected false for non-existent container")
	}
}
