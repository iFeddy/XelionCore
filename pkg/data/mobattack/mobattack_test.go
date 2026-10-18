package mobattack_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/mobattack"
)

func TestLoadFromFile(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "MobAttackSequence")

	db, err := mobattack.LoadFromFile(dirPath)
	if err != nil {
		t.Fatalf("LoadFromFile from dir failed: %v", err)
	}

	if len(db.Mobs) != 297 {
		t.Errorf("Expected 297 mobs, got %d", len(db.Mobs))
	}

	// 1. Verificar Boss Karen
	karen, ok := db.GetMobAI("Karen")
	if !ok {
		t.Fatalf("Mob 'Karen' not found")
	}
	if len(karen.Sequence) != 28 {
		t.Errorf("Expected Karen sequence to have 28 turns, got %d", len(karen.Sequence))
	}
	if len(karen.SkillChanges) != 7 {
		t.Errorf("Expected Karen to have 7 skill change rules, got %d", len(karen.SkillChanges))
	}

	// 2. Verificar Boss Eglack
	eglack, ok := db.GetMobAI("AdlFH_Eglack")
	if !ok {
		t.Fatalf("Mob 'AdlFH_Eglack' not found")
	}
	if len(eglack.Sequence) != 20 {
		t.Errorf("Expected Eglack sequence to have 20 turns, got %d", len(eglack.Sequence))
	}
	if len(eglack.SkillChanges) != 21 {
		t.Errorf("Expected Eglack to have 21 skill change rules, got %d", len(eglack.SkillChanges))
	}

	// 3. Verificar mob básico FKnuckleman
	knuckle, ok := db.GetMobAI("AdlFH_Fknuckleman")
	if !ok {
		t.Fatalf("Mob 'AdlFH_Fknuckleman' not found")
	}
	if len(knuckle.Sequence) != 5 {
		t.Errorf("Expected FKnuckleman to have 5 turns, got %d", len(knuckle.Sequence))
	}
}
