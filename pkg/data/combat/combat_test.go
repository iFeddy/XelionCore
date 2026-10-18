package combat_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
)

func TestCombatStore_LoadAll(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "data", "World", "Combat")
	store := combat.NewStore()

	if err := store.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	// 1. Verify CommonConfig
	common := store.Common()
	if common.RunSpeed != 127 || common.WalkSpeed != 33 || common.AttackSpeed != 13 {
		t.Errorf("Unexpected common config speeds: %+v", common)
	}
	if common.LevelLimit != 125 {
		t.Errorf("Expected LevelLimit 125, got %d", common.LevelLimit)
	}
	if common.MaxFreeStat != 180 {
		t.Errorf("Expected MaxFreeStat 180, got %d", common.MaxFreeStat)
	}

	// 2. Verify ExpTable
	if maxLvl := store.MaxLevel(); maxLvl != 150 {
		t.Errorf("Expected 150 levels in expTable, got %d", maxLvl)
	}
	lvl1Exp, ok := store.GetNextExp(1)
	if !ok || lvl1Exp != 10 {
		t.Errorf("Expected level 1 exp 10, got %d (ok=%v)", lvl1Exp, ok)
	}
	lvl10Exp, ok := store.GetNextExp(10)
	if !ok || lvl10Exp != 2350 {
		t.Errorf("Expected level 10 exp 2350, got %d (ok=%v)", lvl10Exp, ok)
	}

	// 3. Verify Party Exp Bonus
	if bonus := store.GetPartyExpBonus(1); bonus != 1000 {
		t.Errorf("Expected 1-player bonus 1000, got %d", bonus)
	}
	if bonus := store.GetPartyExpBonus(5); bonus != 1800 {
		t.Errorf("Expected 5-player bonus 1800, got %d", bonus)
	}

	// 4. Verify Angle Damage Rates
	front := store.GetAngleDamageRate(false, 0)
	if front != 1000 {
		t.Errorf("Expected front attack rate 1000, got %d", front)
	}
	back := store.GetAngleDamageRate(false, 180)
	if back != 1200 {
		t.Errorf("Expected back attack rate 1200, got %d", back)
	}
	flank := store.GetAngleDamageRate(false, 90)
	if flank != 1100 {
		t.Errorf("Expected 90 degree rate 1100, got %d", flank)
	}
}
