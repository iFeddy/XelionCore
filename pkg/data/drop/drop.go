package drop

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// DropEntry represents an individual item drop rule and rate for a monster.
type DropEntry struct {
	Item       string `json:"item"`
	Rate       int    `json:"rate"` // Rate out of 1,000,000
	MinUpgrade int    `json:"minUpgrade"`
	MaxUpgrade int    `json:"maxUpgrade"`
	Rule       string `json:"rule"`
	Num        int    `json:"num"`
}

// MoneyDrop represents the copper/money drop range and probability.
type MoneyDrop struct {
	Min  int `json:"min"`
	Max  int `json:"max"`
	Rate int `json:"rate"` // Rate out of 1,000,000
}

// TradeBoxDrop represents a reward box drop chance.
type TradeBoxDrop struct {
	Box  string `json:"box"`
	Rate int    `json:"rate"`
}

// MobDropTable defines the complete drop profile for a specific monster.
type MobDropTable struct {
	MapArea       string         `json:"mapArea"`
	MinLevel      int            `json:"minLevel"`
	MaxLevel      int            `json:"maxLevel"`
	AbStateCnt    int            `json:"abStateCnt"`
	Money         MoneyDrop      `json:"money"`
	TradeBoxes    []TradeBoxDrop `json:"tradeBoxes"`
	Drops         []DropEntry    `json:"drops"`
	ExcludedItems []string       `json:"excludedItems"`
}

// ItemGroupEntry represents a member of an item drop group pool.
type ItemGroupEntry struct {
	ItemID       string `json:"itemId"`
	MinQuantity  int    `json:"minQuantity"`
	MaxQuantity  int    `json:"maxQuantity"`
	UpgradeRates []int  `json:"upgradeRates"` // Probabilities for +0 to +15 upgrades
}

// Store holds the parsed drop tables and item groups for ZoneService.
type Store struct {
	mu     sync.RWMutex
	mobs   map[string]MobDropTable
	groups map[string][]ItemGroupEntry
}

// NewStore initializes a new Drop store.
func NewStore() *Store {
	return &Store{
		mobs:   make(map[string]MobDropTable),
		groups: make(map[string][]ItemGroupEntry),
	}
}

// LoadDir loads ItemDropTable.json and ItemDropGroup.json from the specified directory.
func (s *Store) LoadDir(dir string) error {
	// 1. Load ItemDropTable.json
	tablePath := filepath.Join(dir, "ItemDropTable.json")
	tableData, err := os.ReadFile(tablePath)
	if err != nil {
		return fmt.Errorf("read ItemDropTable.json: %w", err)
	}

	loadedMobs := make(map[string]MobDropTable)
	if err := json.Unmarshal(tableData, &loadedMobs); err != nil {
		return fmt.Errorf("unmarshal ItemDropTable.json: %w", err)
	}

	// 2. Load ItemDropGroup.json
	groupPath := filepath.Join(dir, "ItemDropGroup.json")
	groupData, err := os.ReadFile(groupPath)
	if err != nil {
		return fmt.Errorf("read ItemDropGroup.json: %w", err)
	}

	loadedGroups := make(map[string][]ItemGroupEntry)
	if err := json.Unmarshal(groupData, &loadedGroups); err != nil {
		return fmt.Errorf("unmarshal ItemDropGroup.json: %w", err)
	}

	// Case-insensitive lookup normalization
	normalizedMobs := make(map[string]MobDropTable, len(loadedMobs))
	for k, v := range loadedMobs {
		normalizedMobs[strings.ToLower(k)] = v
	}

	normalizedGroups := make(map[string][]ItemGroupEntry, len(loadedGroups))
	for k, v := range loadedGroups {
		normalizedGroups[strings.ToLower(k)] = v
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.mobs = normalizedMobs
	s.groups = normalizedGroups

	return nil
}

// GetMobDrop returns the drop table for a specific mob by ID (case-insensitive).
func (s *Store) GetMobDrop(mobID string) (*MobDropTable, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	table, ok := s.mobs[strings.ToLower(mobID)]
	if !ok {
		return nil, false
	}
	return &table, true
}

// GetGroup returns the items contained in an item drop group pool (case-insensitive).
func (s *Store) GetGroup(groupName string) ([]ItemGroupEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.groups[strings.ToLower(groupName)]
	if !ok {
		return nil, false
	}
	res := make([]ItemGroupEntry, len(entries))
	copy(res, entries)
	return res, true
}

// MobCount returns the total number of mob drop tables loaded.
func (s *Store) MobCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.mobs)
}

// GroupCount returns the total number of item drop groups loaded.
func (s *Store) GroupCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.groups)
}
