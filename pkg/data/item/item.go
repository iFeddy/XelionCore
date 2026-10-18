package item

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ItemOptionEntry defines weight distribution for a bonus stat degree and type.
type ItemOptionEntry struct {
	Degree  int    `json:"degree"`
	Type    string `json:"type"`
	Weights []int  `json:"weights"`
}

// ItemUseFunction defines the server-side action to invoke when consuming an item.
type ItemUseFunction struct {
	Broadcast   bool   `json:"broadcast"`
	UseFunction string `json:"useFunction"`
}

// StatRange defines a minimum and maximum attribute bonus boundary.
type StatRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// RandomOptionEntry defines the green/blue stat generation rules for named equipment.
type RandomOptionEntry struct {
	Unidentified bool      `json:"unidentified"`
	MinOptions   int       `json:"minOptions"`
	MaxOptions   int       `json:"maxOptions"`
	STR          StatRange `json:"str"`
	CON          StatRange `json:"con"`
	DEX          StatRange `json:"dex"`
	INT          StatRange `json:"int"`
	MEN          StatRange `json:"men"`
}

// TreasureRewardItem represents an individual item drop within a container or chest.
type TreasureRewardItem struct {
	ItemID     string `json:"itemId"`
	Lot        int    `json:"lot"`
	UpgradeMin int    `json:"upgradeMin"`
	UpgradeMax int    `json:"upgradeMax"`
	Rate       int    `json:"rate"`
}

// TreasureContainer represents a treasure chest, gift bag, or reward box.
type TreasureContainer struct {
	ContainerID string               `json:"containerId"`
	DummyID     string               `json:"dummyId"`
	MinLot      int                  `json:"minLot"`
	MaxLot      int                  `json:"maxLot"`
	Rewards     []TreasureRewardItem `json:"rewards"`
}

// Store holds the parsed item configurations for ZoneService.
type Store struct {
	mu            sync.RWMutex
	options       []ItemOptionEntry
	useFunctions  map[string]ItemUseFunction
	randomOptions map[string]RandomOptionEntry
	containers    map[string]TreasureContainer
}

// NewStore initializes a new empty Item store.
func NewStore() *Store {
	return &Store{
		useFunctions:  make(map[string]ItemUseFunction),
		randomOptions: make(map[string]RandomOptionEntry),
		containers:    make(map[string]TreasureContainer),
	}
}

// LoadDir loads all 4 item JSON files from the designated directory.
func (s *Store) LoadDir(dir string) error {
	// 1. ItemOptions.json
	optData, err := os.ReadFile(filepath.Join(dir, "ItemOptions.json"))
	if err != nil {
		return fmt.Errorf("read ItemOptions.json: %w", err)
	}
	var loadedOptions []ItemOptionEntry
	if err := json.Unmarshal(optData, &loadedOptions); err != nil {
		return fmt.Errorf("unmarshal ItemOptions.json: %w", err)
	}

	// 2. ItemUseFunction.json
	useData, err := os.ReadFile(filepath.Join(dir, "ItemUseFunction.json"))
	if err != nil {
		return fmt.Errorf("read ItemUseFunction.json: %w", err)
	}
	loadedUse := make(map[string]ItemUseFunction)
	if err := json.Unmarshal(useData, &loadedUse); err != nil {
		return fmt.Errorf("unmarshal ItemUseFunction.json: %w", err)
	}

	// 3. RandomOptionTable.json
	randData, err := os.ReadFile(filepath.Join(dir, "RandomOptionTable.json"))
	if err != nil {
		return fmt.Errorf("read RandomOptionTable.json: %w", err)
	}
	loadedRand := make(map[string]RandomOptionEntry)
	if err := json.Unmarshal(randData, &loadedRand); err != nil {
		return fmt.Errorf("unmarshal RandomOptionTable.json: %w", err)
	}

	// 4. TreasureReward.json
	treasData, err := os.ReadFile(filepath.Join(dir, "TreasureReward.json"))
	if err != nil {
		return fmt.Errorf("read TreasureReward.json: %w", err)
	}
	var loadedContainers []TreasureContainer
	if err := json.Unmarshal(treasData, &loadedContainers); err != nil {
		return fmt.Errorf("unmarshal TreasureReward.json: %w", err)
	}

	// Normalize lookups to case-insensitive
	normUse := make(map[string]ItemUseFunction, len(loadedUse))
	for k, v := range loadedUse {
		normUse[strings.ToLower(k)] = v
	}

	normRand := make(map[string]RandomOptionEntry, len(loadedRand))
	for k, v := range loadedRand {
		normRand[strings.ToLower(k)] = v
	}

	normContainers := make(map[string]TreasureContainer, len(loadedContainers))
	for _, c := range loadedContainers {
		normContainers[strings.ToLower(c.ContainerID)] = c
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.options = loadedOptions
	s.useFunctions = normUse
	s.randomOptions = normRand
	s.containers = normContainers

	return nil
}

// GetUseFunction returns the use function configuration for an item (case-insensitive).
func (s *Store) GetUseFunction(itemID string) (*ItemUseFunction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn, ok := s.useFunctions[strings.ToLower(itemID)]
	if !ok {
		return nil, false
	}
	return &fn, true
}

// GetRandomOption returns the random stat rules for a drop item (case-insensitive).
func (s *Store) GetRandomOption(dropItemIndex string) (*RandomOptionEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.randomOptions[strings.ToLower(dropItemIndex)]
	if !ok {
		return nil, false
	}
	return &entry, true
}

// GetContainer returns a treasure container and its reward contents (case-insensitive).
func (s *Store) GetContainer(containerID string) (*TreasureContainer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.containers[strings.ToLower(containerID)]
	if !ok {
		return nil, false
	}
	return &c, true
}

// Options returns all option degrees and their distribution weights.
func (s *Store) Options() []ItemOptionEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]ItemOptionEntry, len(s.options))
	copy(res, s.options)
	return res
}
