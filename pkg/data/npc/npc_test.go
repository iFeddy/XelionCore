package npc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadNPC(t *testing.T) {
	npcDir := filepath.Join("..", "..", "..", "data", "World", "NPC")

	store, err := Load(npcDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// 1. Check maps count matches Maps/*.json on disk
	mapEntries, err := os.ReadDir(filepath.Join(npcDir, "Maps"))
	if err != nil {
		t.Fatalf("ReadDir(Maps) failed: %v", err)
	}
	wantMaps := 0
	for _, e := range mapEntries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			wantMaps++
		}
	}
	if count := store.MapCount(); count == 0 || count != wantMaps {
		t.Errorf("expected %d maps (matching Maps/*.json), got %d", wantMaps, count)
	}

	// 2. Check total and enabled NPCs
	allNPCs := store.GetAllNPCs(false)
	if len(allNPCs) == 0 {
		t.Errorf("expected non-zero total NPCs")
	}

	enabledNPCs := store.GetAllNPCs(true)
	if len(enabledNPCs) == 0 || len(enabledNPCs) > len(allNPCs) {
		t.Errorf("invalid enabled NPCs count: %d (total=%d)", len(enabledNPCs), len(allNPCs))
	}

	// 3. Test map query (RouN)
	rounAll := store.GetNPCs("RouN", false)
	rounEnabled := store.GetNPCs("RouN", true)
	if len(rounAll) == 0 {
		t.Errorf("expected NPCs in RouN, got 0")
	}
	if len(rounEnabled) > len(rounAll) {
		t.Errorf("enabled NPCs cannot exceed total NPCs")
	}

	// Verify specific NPC
	var foundJames bool
	for _, n := range rounEnabled {
		if n.MobName == "RouSmithJames" {
			foundJames = true
			if n.X != 5645 || n.Y != 8824 || n.Direct != -62 || n.Role != "Merchant" || n.RoleArg != "Weapon" {
				t.Errorf("RouSmithJames fields mismatch: %+v", n)
			}
		}
	}
	if !foundJames {
		t.Errorf("RouSmithJames not found in RouN enabled NPCs")
	}

	// 4. Check links
	allLinks := store.GetLinks(false)
	if len(allLinks) != 268 {
		t.Errorf("expected 268 total links, got %d", len(allLinks))
	}

	enabledLinks := store.GetLinks(true)
	if len(enabledLinks) != 233 {
		t.Errorf("expected 233 enabled links, got %d", len(enabledLinks))
	}

	// Check specific link
	link, ok := store.GetLinkByArg("GateRou1")
	if !ok {
		t.Errorf("GateRou1 not found in LinkTable")
	} else {
		if link.MapServer != "RouCos02" || link.X != 17418 || link.Y != 13832 || !link.Enabled {
			t.Errorf("GateRou1 mismatch: %+v", link)
		}
	}

	// 5. Check NPCAction
	conditions := store.GetConditions()
	if len(conditions) != 20 {
		t.Errorf("expected 20 conditions, got %d", len(conditions))
	}

	cond, ok := store.GetConditionByID(1)
	if !ok || cond.ConditionA != "Default" {
		t.Errorf("Condition 1 mismatch: %+v", cond)
	}

	actions := store.GetActions()
	if len(actions) != 50 {
		t.Errorf("expected 50 actions, got %d", len(actions))
	}

	act, ok := store.GetActionByID(1)
	if !ok || act.ActionIndex != "Link" {
		t.Errorf("Action 1 mismatch: %+v", act)
	}
}
