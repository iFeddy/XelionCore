package item

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// ItemAttrSizes maps official ItemClass (0 to 38) from BoothItemInfo::GetAttrSize to network attribute payload size.
var ItemAttrSizes = [39]int{
	0: 1,   // ByteLot
	1: 2,   // WordLot
	2: 4,   // DwrdLot
	3: 2,   // QuestLot
	4: 38,  // Amulet (trimmed from 62)
	5: 65,  // Weapon (trimmed from 89)
	6: 37,  // Armor
	7: 37,  // Shield
	8: 37,  // Boot
	9: 36,  // Furniture
	10: 8,  // Decoration
	11: 1,  // SkillScroll
	12: 1,  // RecallScroll
	13: 101,// BindItem
	14: 1,  // UpSource
	15: 65, // ItemChest
	16: 1,  // WeaponTitleLicence
	17: 1,  // KingdomQuestItem
	18: 4,  // MiniHouseSkin
	19: 1,  // UpgradeRedGem
	20: 1,  // UpgradeBlueGem
	21: 2,  // KQStep
	22: 1,  // Feed
	23: 19, // Riding
	24: 4,  // Amount
	25: 1,  // UpgradeGoldGem
	26: 4,  // CostumWeapon
	27: 4,  // ActionItem
	28: 1,  // Enchant
	29: 12, // Capsule
	30: 9,  // CardCollect_Close
	31: 5,  // CardCollect
	32: 1,  // NoEffect
	33: 1,  // NoEffect
	34: 1,  // GBCoin
	35: 2,  // ActiveSkill
	36: 26, // Pet
	37: 4,  // CostumShield
	38: 37, // Bracelet (trimmed)
}

// GetAttrSizeForClass returns the exact attribute payload size in bytes for a given ItemClass.
func GetAttrSizeForClass(itemClass uint8) int {
	if int(itemClass) < len(ItemAttrSizes) {
		sz := ItemAttrSizes[itemClass]
		if sz > 0 {
			return sz
		}
	}
	return 4 // default fallback: Amount
}

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
	itemClasses   map[uint16]uint8
}

// NewStore initializes a new empty Item store.
func NewStore() *Store {
	return &Store{
		useFunctions:  make(map[string]ItemUseFunction),
		randomOptions: make(map[string]RandomOptionEntry),
		containers:    make(map[string]TreasureContainer),
		itemClasses:   make(map[uint16]uint8),
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


// LoadItemInfoSHN loads official item definitions from ItemInfo.shn.
func (s *Store) LoadItemInfoSHN(filePath string) error {
	file, err := shn.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	classes := make(map[uint16]uint8, file.RecordCount)
	for _, row := range file.Rows {
		id := row.GetUint16("ID")
		if id > 0 {
			classes[id] = uint8(row.GetInt("Class"))
		}
	}

	s.mu.Lock()
	s.itemClasses = classes
	s.mu.Unlock()

	return nil
}

// GetItemClass returns the official ItemClass for a given itemID.
func (s *Store) GetItemClass(itemID uint16) (uint8, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cls, ok := s.itemClasses[itemID]
	return cls, ok
}

// GetItemAttrSize returns the exact binary attribute payload size in bytes for an item.
func (s *Store) GetItemAttrSize(itemID uint16) int {
	if s == nil {
		return 4
	}
	s.mu.RLock()
	cls, ok := s.itemClasses[itemID]
	s.mu.RUnlock()
	if !ok {
		return 4 // default fallback: Amount (4 bytes)
	}
	return GetAttrSizeForClass(cls)
}
