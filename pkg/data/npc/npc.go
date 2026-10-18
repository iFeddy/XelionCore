package npc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ShineNPC represents an NPC placement record in a world map.
type ShineNPC struct {
	MobName string `json:"mobName"`
	Map     string `json:"map"`
	X       int32  `json:"x"`
	Y       int32  `json:"y"`
	Direct  int16  `json:"direct"`
	Menu    uint8  `json:"menu"`
	Role    string `json:"role"`
	RoleArg string `json:"roleArg"`
	Enabled bool   `json:"enabled"`
}

// LinkTableEntry represents a portal / map link record.
type LinkTableEntry struct {
	Argument  string `json:"argument"`
	MapServer string `json:"mapServer"`
	MapClient string `json:"mapClient"`
	X         int32  `json:"x"`
	Y         int32  `json:"y"`
	Direct    int16  `json:"direct"`
	Party     uint8  `json:"party"`
	Enabled   bool   `json:"enabled"`
}

// NPCCondition represents a condition entry for NPC behavior.
type NPCCondition struct {
	ID         int    `json:"id"`
	ConditionA string `json:"conditionA"`
	TypeA      string `json:"typeA"`
	AX         int32  `json:"ax"`
	AY         int32  `json:"ay"`
	ConditionB string `json:"conditionB"`
	TypeB      string `json:"typeB"`
	BX         int32  `json:"bx"`
	BY         int32  `json:"by"`
}

// NPCAction represents an action entry for NPC behavior.
type NPCAction struct {
	ID          int    `json:"id"`
	ActionIndex string `json:"actionIndex"`
	Type        string `json:"type"`
	X           int32  `json:"x"`
	Y           int32  `json:"y"`
}

// NPCActionData holds both conditions and actions.
type NPCActionData struct {
	Conditions []NPCCondition `json:"conditions"`
	Actions    []NPCAction    `json:"actions"`
}

// Store holds the parsed NPC data with fast indexed lookups.
type Store struct {
	mu         sync.RWMutex
	npcsByMap  map[string][]ShineNPC
	links      []LinkTableEntry
	linksByArg map[string]LinkTableEntry
	actions    NPCActionData
	condByID   map[int]NPCCondition
	actByID    map[int]NPCAction
}

// Load reads all map JSON files, LinkTable.json, and NPCAction.json from npcDir.
func Load(npcDir string) (*Store, error) {
	store := &Store{
		npcsByMap:  make(map[string][]ShineNPC),
		linksByArg: make(map[string]LinkTableEntry),
		condByID:   make(map[int]NPCCondition),
		actByID:    make(map[int]NPCAction),
	}

	// 1. Load Maps/
	mapsDir := filepath.Join(npcDir, "Maps")
	entries, err := os.ReadDir(mapsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read maps directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		mapName := strings.TrimSuffix(entry.Name(), ".json")
		filePath := filepath.Join(mapsDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read map file %s: %w", entry.Name(), err)
		}

		var npcs []ShineNPC
		if err := json.Unmarshal(data, &npcs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", entry.Name(), err)
		}

		store.npcsByMap[mapName] = npcs
	}

	// 2. Load LinkTable.json
	linkPath := filepath.Join(npcDir, "LinkTable.json")
	if linkData, err := os.ReadFile(linkPath); err == nil {
		var links []LinkTableEntry
		if err := json.Unmarshal(linkData, &links); err != nil {
			return nil, fmt.Errorf("failed to unmarshal LinkTable.json: %w", err)
		}
		store.links = links
		for _, l := range links {
			store.linksByArg[l.Argument] = l
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read LinkTable.json: %w", err)
	}

	// 3. Load NPCAction.json
	actionPath := filepath.Join(npcDir, "NPCAction.json")
	if actionData, err := os.ReadFile(actionPath); err == nil {
		var ad NPCActionData
		if err := json.Unmarshal(actionData, &ad); err != nil {
			return nil, fmt.Errorf("failed to unmarshal NPCAction.json: %w", err)
		}
		store.actions = ad
		for _, c := range ad.Conditions {
			store.condByID[c.ID] = c
		}
		for _, a := range ad.Actions {
			store.actByID[a.ID] = a
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read NPCAction.json: %w", err)
	}

	return store, nil
}

// GetNPCs returns the NPCs for a given map. If onlyEnabled is true, only active NPCs are returned.
func (s *Store) GetNPCs(mapName string, onlyEnabled bool) []ShineNPC {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list, ok := s.npcsByMap[mapName]
	if !ok {
		return nil
	}

	if !onlyEnabled {
		res := make([]ShineNPC, len(list))
		copy(res, list)
		return res
	}

	var res []ShineNPC
	for _, n := range list {
		if n.Enabled {
			res = append(res, n)
		}
	}
	return res
}

// GetAllNPCs returns all NPCs across all maps.
func (s *Store) GetAllNPCs(onlyEnabled bool) []ShineNPC {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []ShineNPC
	for _, list := range s.npcsByMap {
		for _, n := range list {
			if !onlyEnabled || n.Enabled {
				res = append(res, n)
			}
		}
	}
	return res
}

// MapCount returns the number of distinct map files loaded.
func (s *Store) MapCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.npcsByMap)
}

// GetLinks returns all links / portal entries.
func (s *Store) GetLinks(onlyEnabled bool) []LinkTableEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !onlyEnabled {
		res := make([]LinkTableEntry, len(s.links))
		copy(res, s.links)
		return res
	}

	var res []LinkTableEntry
	for _, l := range s.links {
		if l.Enabled {
			res = append(res, l)
		}
	}
	return res
}

// GetLinkByArg returns a link entry by its argument name.
func (s *Store) GetLinkByArg(arg string) (LinkTableEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.linksByArg[arg]
	return l, ok
}

// GetConditions returns all NPC conditions.
func (s *Store) GetConditions() []NPCCondition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]NPCCondition, len(s.actions.Conditions))
	copy(res, s.actions.Conditions)
	return res
}

// GetConditionByID retrieves an NPCCondition by ID.
func (s *Store) GetConditionByID(id int) (NPCCondition, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.condByID[id]
	return c, ok
}

// GetActions returns all NPC actions.
func (s *Store) GetActions() []NPCAction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]NPCAction, len(s.actions.Actions))
	copy(res, s.actions.Actions)
	return res
}

// GetActionByID retrieves an NPCAction by ID.
func (s *Store) GetActionByID(id int) (NPCAction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.actByID[id]
	return a, ok
}
