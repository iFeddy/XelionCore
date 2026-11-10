package blockinfo_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
)

func TestBlockInfo_LoadAll301Files(t *testing.T) {
	blockInfoDir := filepath.Join("..", "..", "..", "data", "BlockInfo")

	store, err := blockinfo.LoadAll(blockInfoDir)
	if err != nil {
		t.Fatalf("Failed to load BlockInfo directory: %v", err)
	}

	if len(store.SHBD) != 158 {
		t.Errorf("Expected 158 .shbd files, got %d", len(store.SHBD))
	}
	if len(store.SBI) != 33 {
		t.Errorf("Expected 33 .sbi files, got %d", len(store.SBI))
	}
	if len(store.AID) != 23 {
		t.Errorf("Expected 23 .aid files, got %d", len(store.AID))
	}
	if len(store.BDT) != 76 {
		t.Errorf("Expected 76 .bdt files, got %d", len(store.BDT))
	}
	if len(store.SHAB) != 8 {
		t.Errorf("Expected 8 .shab files, got %d", len(store.SHAB))
	}
	if len(store.SHMD) != 2 {
		t.Errorf("Expected 2 .shmd files, got %d", len(store.SHMD))
	}
	if len(store.SHAD) != 1 {
		t.Errorf("Expected 1 .shad files, got %d", len(store.SHAD))
	}

	total := store.TotalCount()
	if total != 301 {
		t.Fatalf("Expected exactly 301 files in data/BlockInfo/, got %d", total)
	}

	t.Logf("BlockInfo Store successfully loaded all %d files (%d shbd, %d sbi, %d aid, %d bdt, %d shab, %d shmd, %d shad)",
		total, len(store.SHBD), len(store.SBI), len(store.AID), len(store.BDT), len(store.SHAB), len(store.SHMD), len(store.SHAD))

	// Test specific maps: Rou.shbd, RouN.shbd, Adl.shbd
	rou := store.GetSHBD("Rou")
	if rou == nil {
		t.Fatal("Rou.shbd missing")
	}
	if rou.WidthTiles != 2048 || rou.HeightTiles != 2048 {
		t.Errorf("Expected Rou.shbd 2048x2048, got %dx%d", rou.WidthTiles, rou.HeightTiles)
	}

	// Test boundary blocking (outside map is blocked)
	if !rou.IsMoveBlocked(3000, 3000) {
		t.Error("Expected out of bounds tile (3000, 3000) to be blocked")
	}

	// Test Bresenham attack blocked between blocked tiles
	if !rou.IsAttackBlocked(0, 0, 10, 10) {
		t.Error("Expected ray across border (0,0)->(10,10) in Rou to be blocked")
	}

	// Test SBI: AdlF.sbi
	adlfSBI := store.GetSBI("AdlF")
	if adlfSBI == nil {
		t.Fatal("AdlF.sbi missing")
	}
	instDoors := store.NewInstanceDoors("AdlF")
	if instDoors == nil {
		t.Fatal("NewInstanceDoors(AdlF) failed")
	}
	if len(instDoors.DoorList) != 2 {
		t.Fatalf("Expected 2 doors in AdlF.sbi, got %d", len(instDoors.DoorList))
	}

	// DoorBlock01 starts closed
	d1 := instDoors.Doors["DoorBlock01"]
	if d1 == nil {
		t.Fatal("DoorBlock01 not found")
	}
	if d1.IsOpen {
		t.Error("DoorBlock01 should be closed initially")
	}

	// Open DoorBlock01
	if !instDoors.DoorAction("DoorBlock01", blockinfo.DoorActionOpen) {
		t.Error("DoorAction open failed")
	}
	if !d1.IsOpen {
		t.Error("DoorBlock01 should be open after DoorAction")
	}

	// CloseAllDoors
	instDoors.CloseAllDoors()
	if d1.IsOpen {
		t.Error("DoorBlock01 should be closed after CloseAllDoors")
	}

	// Test AID: RouN.aid
	rounAID := store.GetAID("RouN")
	if rounAID == nil {
		t.Fatal("RouN.aid missing")
	}
	// Contains check for MH_Zone1 in RouN: center is (6411.73, 6527.53)
	if !rounAID.Contains("MH_Zone1", 6411.0, 6527.0) {
		t.Error("Expected RouN.aid to contain MH_Zone1 at (6411, 6527)")
	}
}

