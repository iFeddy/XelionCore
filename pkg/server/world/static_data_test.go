package world

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWorldDataStore_AllFiles(t *testing.T) {
	// Locate data directory relative to this test file (pkg/server/world -> root is ../../..)
	dataDir := filepath.Join("..", "..", "..", "data")
	if _, err := os.Stat(filepath.Join(dataDir, "BadNameFilter.shn")); os.IsNotExist(err) {
		t.Skip("data directory not found, skipping file loading test")
	}

	store, err := LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil WorldDataStore")
	}

	// Verify ChecksumEngine was initialized with 49 files
	if store.Checksums == nil {
		t.Fatal("expected non-nil ChecksumEngine")
	}
	checksums := store.Checksums.GetChecksums()
	if len(checksums) != 49 {
		t.Errorf("expected 49 checksum files loaded, got %d", len(checksums))
	}
	// Verify first and last checksums are not empty
	if checksums[0].String() == "" || checksums[48].String() == "" {
		t.Errorf("expected valid non-empty checksums at indices 0 and 48")
	}

	// Verify BadNameFilter
	if store.BadNames == nil {
		t.Fatal("expected non-nil BadNameFilter")
	}
	if !store.BadNames.IsBadName("admin") {
		t.Errorf("expected 'admin' to be flagged as bad name")
	}
	if store.BadNames.IsBadName("Legolas") {
		t.Errorf("expected 'Legolas' not to be flagged as bad name")
	}

	// Verify Default Character Data loaded for all 6 official base classes (Fighter, Cleric, Archer, Mage, Joker, Sentinel)
	if store.DefaultData == nil {
		t.Fatal("expected non-nil DefaultData")
	}
	classes := []byte{1, 6, 11, 16, 21, 26}
	for _, c := range classes {
		cd, ok := store.DefaultData.GetClassData(c)
		if !ok {
			t.Errorf("expected class data for class ID %d", c)
		}
		if cd.HP <= 0 || cd.SP <= 0 {
			t.Errorf("invalid HP/SP for class %d: HP=%d SP=%d", c, cd.HP, cd.SP)
		}
	}

	// Verify Appearance Validation against official FaceInfo, HairInfo, HairColorInfo
	// Joker with Hair=1, Color=12, Face=4 (official starter selection)
	if !store.ValidateAppearance(21, 1, 12, 4) {
		t.Errorf("expected valid appearance for Joker (Hair=1, Color=12, Face=4)")
	}
	// Invalid hair ID 250 not in HairInfo.shn
	if store.ValidateAppearance(1, 250, 1, 1) {
		t.Errorf("expected invalid hair=250 to fail appearance check")
	}
	// Premium face ID 13 (Grade 4) must be rejected on creation (requires coupon)
	if store.ValidateAppearance(1, 1, 0, 13) {
		t.Errorf("expected Premium face=13 (Grade 4) to be rejected at character creation")
	}
	// Premium color ID 3 (Grade 4) must be rejected on creation (requires coupon)
	if store.ValidateAppearance(1, 1, 3, 0) {
		t.Errorf("expected Premium color=3 (Grade 4) to be rejected at character creation")
	}

	// Verify Loaded SHN files count (22 WorldManager + 4 Character = 26)
	if len(store.LoadedFiles) < 26 {
		t.Errorf("expected at least 26 SHN files loaded in WorldDataStore, got %d", len(store.LoadedFiles))
	}
}
