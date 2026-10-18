package npc

import (
	"path/filepath"
	"testing"
)

func TestLoadNPC(t *testing.T) {
	npcDir := filepath.Join("..", "..", "..", "data", "World", "NPC")

	store, err := Load(npcDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// 1. Check maps count
	if count := store.MapCount(); count != 95 {
		t.Errorf("expected 95 maps, got %d", count)
	}

	// 2. Check total NPCs
	allNPCs := store.GetAllNPCs(false)
	if len(allNPCs) != 657 {
		t.Errorf("expected 657 total NPCs, got %d", len(allNPCs))
	}

	enabledNPCs := store.GetAllNPCs(true)
	if len(enabledNPCs) != 530 {
		t.Errorf("expected 530 enabled NPCs, got %d", len(enabledNPCs))
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
