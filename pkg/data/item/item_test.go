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

func TestBuildItemAttrBuffer_DynamicSizes(t *testing.T) {
	// Class 6, 7, 8, 38 (Armor/Shield/Boot/Bracelet): 13 + (3 * numOpts), offset 12 = (numOpts << 1) | 1
	for _, cls := range []uint8{6, 7, 8, 38} {
		buf0 := item.BuildItemAttrBuffer(cls, 1, nil, 0)
		if len(buf0) != 13 || buf0[12] != 0x01 {
			t.Errorf("Class %d (0 opts): expected len=13, buf[12]=0x01; got len=%d, buf[12]=0x%02X", cls, len(buf0), buf0[12])
		}
		opts := []item.ItemStatOption{{Type: 0, Value: 5}, {Type: 1, Value: 10}}
		buf2 := item.BuildItemAttrBuffer(cls, 1, opts, 0)
		if len(buf2) != 19 || buf2[12] != 0x05 {
			t.Errorf("Class %d (2 opts): expected len=19, buf[12]=0x05; got len=%d, buf[12]=0x%02X", cls, len(buf2), buf2[12])
		}
	}

	// Class 5 (Weapon): 65 + (3 * numOpts), offset 64 = (numOpts << 1) | 1
	w0 := item.BuildItemAttrBuffer(5, 1, nil, 0)
	if len(w0) != 65 || w0[64] != 0x01 {
		t.Errorf("Class 5 (0 opts): expected len=65, w0[64]=0x01; got len=%d, w0[64]=0x%02X", len(w0), w0[64])
	}
	w1 := item.BuildItemAttrBuffer(5, 1, []item.ItemStatOption{{Type: 2, Value: 7}}, 0)
	if len(w1) != 68 || w1[64] != 0x03 {
		t.Errorf("Class 5 (1 opt): expected len=68, w1[64]=0x03; got len=%d, w1[64]=0x%02X", len(w1), w1[64])
	}

	// Class 4 (Amulet/Ring): 38 + (3 * numOpts), offset 37 = (numOpts << 1) | 1
	a0 := item.BuildItemAttrBuffer(4, 1, nil, 0)
	if len(a0) != 38 || a0[37] != 0x01 {
		t.Errorf("Class 4 (0 opts): expected len=38, a0[37]=0x01; got len=%d, a0[37]=0x%02X", len(a0), a0[37])
	}

	// Class 15 (Chest): 1 + (8 * items_adentro)
	c0 := item.BuildItemAttrBuffer(15, 1, nil, 0)
	if len(c0) != 1 || c0[0] != 0 {
		t.Errorf("Class 15 (0 inner): expected len=1, c0[0]=0; got len=%d, c0[0]=%d", len(c0), c0[0])
	}
	c2 := item.BuildItemAttrBuffer(15, 1, nil, 2)
	if len(c2) != 17 || c2[0] != 2 {
		t.Errorf("Class 15 (2 inner): expected len=17, c2[0]=2; got len=%d, c2[0]=%d", len(c2), c2[0])
	}

	// Class 13 (BindItem): 1 + (10 * count)
	b0 := item.BuildItemAttrBuffer(13, 1, nil, 0)
	if len(b0) != 1 || b0[0] != 0 {
		t.Errorf("Class 13 (0 count): expected len=1, b0[0]=0; got len=%d, b0[0]=%d", len(b0), b0[0])
	}
	b3 := item.BuildItemAttrBuffer(13, 1, nil, 3)
	if len(b3) != 31 || b3[0] != 3 {
		t.Errorf("Class 13 (3 count): expected len=31, b3[0]=3; got len=%d, b3[0]=%d", len(b3), b3[0])
	}

	// Class 10 (Decoration): 8 bytes fixed
	d := item.BuildItemAttrBuffer(10, 1, nil, 0)
	if len(d) != 8 {
		t.Errorf("Class 10: expected len=8, got %d", len(d))
	}

	// Simple classes: 1B, 2B, 4B
	if s1 := item.BuildItemAttrBuffer(0, 5, nil, 0); len(s1) != 1 || s1[0] != 5 {
		t.Errorf("Class 0: expected len=1, val=5; got len=%d", len(s1))
	}
	if s2 := item.BuildItemAttrBuffer(1, 300, nil, 0); len(s2) != 2 {
		t.Errorf("Class 1: expected len=2; got len=%d", len(s2))
	}
	if s4 := item.BuildItemAttrBuffer(24, 1000, nil, 0); len(s4) != 4 {
		t.Errorf("Class 24: expected len=4; got len=%d", len(s4))
	}
}
