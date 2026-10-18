package mobregen_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/mobregen"
)

func TestLoadMapRegen_RouVal01(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "data", "MobRegen")

	reg, err := mobregen.LoadMapRegen(baseDir, "RouVal01")
	if err != nil {
		t.Fatalf("LoadMapRegen failed for RouVal01: %v", err)
	}

	if len(reg.Groups) != 91 {
		t.Errorf("Expected 91 groups in RouVal01, got %d", len(reg.Groups))
	}

	// Verificar grupo RVal01_01
	g1, ok := reg.Groups["RVal01_01"]
	if !ok {
		t.Fatalf("Group RVal01_01 not found")
	}
	if !g1.Enabled {
		t.Errorf("Expected RVal01_01 to be enabled")
	}
	if g1.CenterX != 18950 || g1.CenterY != 20240 || g1.RangeDegree != 1130 {
		t.Errorf("Coords mismatch in RVal01_01: (%d, %d, %d)", g1.CenterX, g1.CenterY, g1.RangeDegree)
	}
	if len(g1.Mobs) != 3 {
		t.Fatalf("Expected 3 mobs in RVal01_01, got %d", len(g1.Mobs))
	}
	if g1.Mobs[0].Mob != "LittleLizardMan" || g1.Mobs[0].Count != 3 {
		t.Errorf("Mob mismatch: %s count %d", g1.Mobs[0].Mob, g1.Mobs[0].Count)
	}
}

func TestLoadMapRegen_RouN_CommentedRecords(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "data", "MobRegen")

	reg, err := mobregen.LoadMapRegen(baseDir, "RouN")
	if err != nil {
		t.Fatalf("LoadMapRegen failed for RouN: %v", err)
	}

	if len(reg.Groups) != 4 {
		t.Errorf("Expected 4 groups in RouN, got %d", len(reg.Groups))
	}

	// Verificar árbol de navidad (comentado / disabled)
	xtree, ok := reg.Groups["XTree"]
	if !ok {
		t.Fatalf("XTree not found in RouN")
	}
	if xtree.Enabled {
		t.Errorf("Expected XTree to be disabled (commented in original txt)")
	}
	if len(xtree.Mobs) != 1 || xtree.Mobs[0].Mob != "ChristmasTree" {
		t.Errorf("Expected ChristmasTree in XTree mob list")
	}
	if xtree.Mobs[0].Enabled {
		t.Errorf("Expected ChristmasTree mob to be disabled")
	}

	// Active groups en RouN deben ser 0
	active := reg.GetActiveGroups()
	if len(active) != 0 {
		t.Errorf("Expected 0 active groups in RouN, got %d", len(active))
	}
}

func TestLoadMapRegen_Subdirectory_KQ(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "data", "MobRegen")

	reg, err := mobregen.LoadMapRegen(baseDir, "KDEddyHill")
	if err != nil {
		t.Fatalf("LoadMapRegen failed for KDEddyHill in KingdomQuest subdirectory: %v", err)
	}

	if len(reg.Groups) != 56 {
		t.Errorf("Expected 56 groups in KDEddyHill, got %d", len(reg.Groups))
	}
}

func TestLoadMapRegen_SpaceGroupNames_BerKal01(t *testing.T) {
	baseDir := filepath.Join("..", "..", "..", "data", "MobRegen")

	reg, err := mobregen.LoadMapRegen(baseDir, "BerKal01")
	if err != nil {
		t.Fatalf("LoadMapRegen failed for BerKal01: %v", err)
	}

	// Verificar grupo 'Collection 1' con espacios
	col1, ok := reg.Groups["Collection 1"]
	if !ok {
		t.Fatalf("Group 'Collection 1' not found in BerKal01")
	}
	if col1.CenterX != 6624 || col1.CenterY != 8492 {
		t.Errorf("Coords mismatch in Collection 1: (%d, %d)", col1.CenterX, col1.CenterY)
	}
	if len(col1.Mobs) != 4 {
		t.Errorf("Expected 4 gathering resources in Collection 1, got %d", len(col1.Mobs))
	}
}
