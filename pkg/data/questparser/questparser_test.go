package questparser_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/questparser"
)

func TestQuestParserStore_Load(t *testing.T) {
	filePath := filepath.Join("..", "..", "..", "data", "World", "Quest", "QuestParser.json")
	store := questparser.NewStore()

	if err := store.LoadFile(filePath); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if count := store.Count(); count != 25 {
		t.Fatalf("Expected 25 classes, got %d", count)
	}

	// 1. Verify Fighter -> 1
	fighterID, ok := store.GetID("Fighter")
	if !ok || fighterID != 1 {
		t.Errorf("Expected Fighter ID 1, got %d (ok=%v)", fighterID, ok)
	}

	// 2. Verify case-insensitivity: "assassin" -> 25
	assassinID, ok := store.GetID("assassin")
	if !ok || assassinID != 25 {
		t.Errorf("Expected Assassin ID 25, got %d (ok=%v)", assassinID, ok)
	}

	// 3. Verify reverse lookup: 16 -> "Mage"
	mageName, ok := store.GetName(16)
	if !ok || mageName != "Mage" {
		t.Errorf("Expected class 16 to be Mage, got %s (ok=%v)", mageName, ok)
	}

	// 4. Non-existent lookups
	if _, ok := store.GetID("UnknownClass"); ok {
		t.Errorf("Expected false for unknown class")
	}
	if _, ok := store.GetName(999); ok {
		t.Errorf("Expected false for unknown class ID")
	}
}
