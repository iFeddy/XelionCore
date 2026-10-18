package misc_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/misc"
)

func TestMiscStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "Misc")
	store := misc.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	// 1. Verify RecallCoord
	rouRecall, ok := store.GetRecallCoord("StageRouRecallScroll")
	if !ok {
		t.Fatalf("StageRouRecallScroll not found")
	}
	if rouRecall.MapName != "RouN" || rouRecall.X != 6445 || rouRecall.Y != 8630 {
		t.Errorf("Unexpected Rou recall coords: %+v", rouRecall)
	}

	// 2. Verify StealthMatrix
	if !store.CanSee("Base", "Base") {
		t.Errorf("Expected Base to see Base")
	}
	if store.CanSee("Base", "AdminHide") {
		t.Errorf("Expected Base NOT to see AdminHide")
	}
	if !store.CanAttack("Base", "Base") {
		t.Errorf("Expected Base can attack Base")
	}
	if store.CanAttack("Observer", "Base") {
		t.Errorf("Expected Observer cannot attack Base")
	}

	// 3. Verify SkillMechanics
	mech := store.Mechanics()
	if len(mech.SkillBreedMobs) != 79 {
		t.Errorf("Expected 79 breed mobs, got %d", len(mech.SkillBreedMobs))
	}
	if len(mech.Polymorph) != 1 {
		t.Errorf("Expected 1 polymorph entry, got %d", len(mech.Polymorph))
	}
	if len(mech.AbStateReset) != 7 {
		t.Errorf("Expected 7 abstate reset entries, got %d", len(mech.AbStateReset))
	}

	// 4. Verify MobChats
	chats := store.MobChats()
	if len(chats.Peace) != 330 {
		t.Errorf("Expected 330 peace chats, got %d", len(chats.Peace))
	}
	if len(chats.Attack) != 330 {
		t.Errorf("Expected 330 attack chats, got %d", len(chats.Attack))
	}
}
