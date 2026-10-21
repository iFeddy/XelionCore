package item

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// ItemStatOption represents a single 3-byte stat option in ItemAttributeClass (type + value).
type ItemStatOption struct {
	Type  uint8
	Value uint16
}

// CalcAttrSizeForClass calculates the dynamic binary attribute size in bytes emulating
// iac_fill2buffer in Zone.exe based on ItemClass and the number of options/inner items (numOpts).
func CalcAttrSizeForClass(itemClass uint8, numOpts int) int {
	if numOpts < 0 {
		numOpts = 0
	}
	switch itemClass {
	case 6, 7, 8, 38: // Armor, Shield, Boot, Bracelet: 13 + (3 * numOpts)
		return 13 + (3 * numOpts)
	case 5: // Weapon: 65 + (3 * numOpts)
		return 65 + (3 * numOpts)
	case 4: // Amulet / Ring: 38 + (3 * numOpts)
		return 38 + (3 * numOpts)
	case 15: // Chest: 1 + (8 * items_adentro)
		return 1 + (8 * numOpts)
	case 13: // BindItem: 1 + (10 * count)
		return 1 + (10 * numOpts)
	case 10: // Decoration: 8 bytes fixed
		return 8
	case 0, 11, 12, 14, 16, 19, 20, 22, 25, 28, 33, 34: // Simple 1-byte classes (ByteLot, Scrolls, Gems, etc.)
		return 1
	case 1, 3, 21, 35: // 2-byte classes (WordLot, QuestLot, KQStep, ActiveSkill)
		return 2
	case 2, 18, 24, 26, 27, 37: // 4-byte classes (DwrdLot, MiniHouseSkin, Amount, CostumWeapon, ActionItem, CostumShield)
		return 4
	case 9: // Furniture
		return 36
	case 17, 32: // KingdomQuestItem / NoEffect (0-byte payload in iac_fill2buffer)
		return 0
	case 23: // Riding
		return 19
	case 29: // Capsule
		return 12
	case 30: // CardCollect_Close
		return 9
	case 31: // CardCollect
		return 5
	case 36: // Pet
		return 26
	default:
		return 4
	}
}

// GetAttrSizeForClass returns the base (0-option) attribute payload size in bytes for a given ItemClass.
func GetAttrSizeForClass(itemClass uint8) int {
	return CalcAttrSizeForClass(itemClass, 0)
}

// BuildItemAttrBuffer serializes the dynamic item attribute slice (attr) emulating iac_fill2buffer.
// The returned byte slice has the exact dynamic length required by PROTO_NC_CHAR_CLIENT_ITEM_CMD without trailing zero padding.
func BuildItemAttrBuffer(itemClass uint8, amount uint32, opts []ItemStatOption, innerCount int) []byte {
	numOpts := len(opts)
	if itemClass == 13 || itemClass == 15 {
		numOpts = innerCount
	}
	size := CalcAttrSizeForClass(itemClass, numOpts)
	if size == 0 {
		return nil
	}
	buf := make([]byte, size)

	switch itemClass {
	case 6, 7, 8, 38: // Armor, Shield, Boot, Bracelet: offset 12 holds (numOpts << 1) | 1
		buf[12] = byte((numOpts << 1) | 1)
		writeStatOptions(buf[13:], opts)
	case 5: // Weapon: offset 64 holds (numOpts << 1) | 1
		buf[64] = byte((numOpts << 1) | 1)
		writeStatOptions(buf[65:], opts)
	case 4: // Amulet / Ring: offset 37 holds (numOpts << 1) | 1
		buf[37] = byte((numOpts << 1) | 1)
		writeStatOptions(buf[38:], opts)
	case 15: // Chest: 1 byte count + (8 * items_adentro)
		buf[0] = byte(numOpts)
	case 13: // BindItem: 1 byte count + (10 * count)
		buf[0] = byte(numOpts)
	case 10: // Decoration: 8 bytes fixed
		// Initialized to 8 clean bytes
	case 0, 11, 12, 14, 16, 19, 20, 22, 25, 28, 33, 34: // 1 byte fixed
		if amount == 0 {
			amount = 1
		}
		buf[0] = byte(amount)
	case 1, 3, 21, 35: // 2 bytes fixed
		if amount == 0 {
			amount = 1
		}
		binary.LittleEndian.PutUint16(buf[0:2], uint16(amount))
	case 2, 18, 24, 26, 27, 37: // 4 bytes fixed
		if amount == 0 {
			amount = 1
		}
		binary.LittleEndian.PutUint32(buf[0:4], amount)
	}

	return buf
}

func writeStatOptions(dst []byte, opts []ItemStatOption) {
	for i, opt := range opts {
		off := i * 3
		if off+3 <= len(dst) {
			dst[off] = opt.Type
			binary.LittleEndian.PutUint16(dst[off+1:off+3], opt.Value)
		}
	}
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

// ItemInfo holds the official item definition parsed from ItemInfo.shn.
type ItemInfo struct {
	ID        uint16
	InxName   string
	Name      string
	SortOrder int
	Type      uint32
	Class     uint8
	MaxLot    uint32
	Equip     uint8
	TwoHand   bool
	AtkSpeed  uint32
	DemandLv  uint32
	Grade     uint32
	MinWC     uint32
	MaxWC     uint32
	AC        uint32
	MinMA     uint32
	MaxMA     uint32
	MR        uint32
	TH        uint32
	TB        uint32
	UseClass  uint32
	BuyPrice  uint32
	SellPrice uint32
	UpLimit   uint8
}

// ItemUseEffect holds the official consumable recovery/buff effect parsed from ItemUseEffect.shn.
type ItemUseEffect struct {
	ItemIndex       string
	UseEffectA      uint32
	UseValueA       uint16
	UseEffectB      uint32
	UseValueB       uint16
	UseEffectC      uint32
	UseValueC       uint16
	UseAbstateName  string
	UseAbstateIndex uint32
}

// Store holds the parsed item configurations for ZoneService.
type Store struct {
	mu            sync.RWMutex
	options       []ItemOptionEntry
	useFunctions  map[string]ItemUseFunction
	randomOptions map[string]RandomOptionEntry
	containers    map[string]TreasureContainer
	itemClasses   map[uint16]uint8
	itemsByID     map[uint16]*ItemInfo
	itemsByName   map[string]*ItemInfo
	useEffects    map[string]*ItemUseEffect
	useClasses    map[uint32]map[string]bool
	classPrefixes map[byte]string
}

// NewStore initializes a new empty Item store.
func NewStore() *Store {
	return &Store{
		useFunctions:  make(map[string]ItemUseFunction),
		randomOptions: make(map[string]RandomOptionEntry),
		containers:    make(map[string]TreasureContainer),
		itemClasses:   make(map[uint16]uint8),
		itemsByID:     make(map[uint16]*ItemInfo),
		itemsByName:   make(map[string]*ItemInfo),
		useEffects:    make(map[string]*ItemUseEffect),
		useClasses:    make(map[uint32]map[string]bool),
		classPrefixes: make(map[byte]string),
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
	byID := make(map[uint16]*ItemInfo, file.RecordCount)
	byName := make(map[string]*ItemInfo, file.RecordCount)

	for _, row := range file.Rows {
		id := row.GetUint16("ID")
		inxName := strings.TrimSpace(row.GetString("InxName"))
		cls := uint8(row.GetInt("Class"))
		info := &ItemInfo{
			ID:        id,
			InxName:   inxName,
			Name:      strings.TrimSpace(row.GetString("Name")),
			SortOrder: -1,
			Type:      uint32(row.GetInt("Type")),
			Class:     cls,
			MaxLot:    uint32(row.GetInt("MaxLot")),
			Equip:     uint8(row.GetInt("Equip")),
			TwoHand:   row.GetByte("TwoHand") != 0,
			AtkSpeed:  uint32(row.GetInt("AtkSpeed")),
			DemandLv:  uint32(row.GetInt("DemandLv")),
			Grade:     uint32(row.GetInt("Grade")),
			MinWC:     uint32(row.GetInt("MinWC")),
			MaxWC:     uint32(row.GetInt("MaxWC")),
			AC:        uint32(row.GetInt("AC")),
			MinMA:     uint32(row.GetInt("MinMA")),
			MaxMA:     uint32(row.GetInt("MaxMA")),
			MR:        uint32(row.GetInt("MR")),
			TH:        uint32(row.GetInt("TH")),
			TB:        uint32(row.GetInt("TB")),
			UseClass:  uint32(row.GetInt("UseClass")),
			BuyPrice:  uint32(row.GetInt("BuyPrice")),
			SellPrice: uint32(row.GetInt("SellPrice")),
			UpLimit:   row.GetByte("UpLimit"),
		}
		classes[id] = cls
		byID[id] = info
		if inxName != "" {
			byName[strings.ToLower(inxName)] = info
		}
	}

	s.mu.Lock()
	s.itemClasses = classes
	s.itemsByID = byID
	s.itemsByName = byName
	s.mu.Unlock()

	return nil
}

// LoadInventoryDataSHN loads official inventory sorting and class eligibility data.
func (s *Store) LoadInventoryDataSHN(itemInfoServerPath, itemSortPath, useClassTypeInfoPath, classNamePath string) error {
	itemSort, err := shn.ReadFile(itemSortPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", itemSortPath, err)
	}
	sortOrder := make(map[string]int, itemSort.RecordCount)
	for i, row := range itemSort.Rows {
		name := strings.ToLower(strings.TrimSpace(row.GetString("ItemSort_Index")))
		if name != "" {
			sortOrder[name] = i
		}
	}

	itemInfoServer, err := shn.ReadFile(itemInfoServerPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", itemInfoServerPath, err)
	}
	itemOrders := make(map[uint16]int, itemInfoServer.RecordCount)
	for _, row := range itemInfoServer.Rows {
		category := strings.ToLower(strings.TrimSpace(row.GetString("ItemSort_Index")))
		order, ok := sortOrder[category]
		if ok {
			itemOrders[row.GetUint16("ID")] = order
		}
	}

	useClassInfo, err := shn.ReadFile(useClassTypeInfoPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", useClassTypeInfoPath, err)
	}
	classColumns := make([]string, 0, len(useClassInfo.Columns)-1)
	for _, column := range useClassInfo.Columns {
		if column.Name != "" && column.Name != "UseClass" {
			classColumns = append(classColumns, column.Name)
		}
	}
	useClasses := make(map[uint32]map[string]bool, useClassInfo.RecordCount)
	for _, row := range useClassInfo.Rows {
		useClass := uint32(row.GetInt("UseClass"))
		allowed := make(map[string]bool, len(classColumns))
		for _, class := range classColumns {
			allowed[strings.ToLower(class)] = row.GetByte(class) != 0
		}
		useClasses[useClass] = allowed
	}

	classNameInfo, err := shn.ReadFile(classNamePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", classNamePath, err)
	}
	classPrefixes := make(map[byte]string, classNameInfo.RecordCount)
	for _, row := range classNameInfo.Rows {
		classID := row.GetByte("ClassID")
		prefix := strings.TrimSpace(row.GetString("acPrefix"))
		if prefix != "" && prefix != "-" {
			classPrefixes[classID] = prefix
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for id, order := range itemOrders {
		if info, ok := s.itemsByID[id]; ok {
			info.SortOrder = order
		}
	}
	s.useClasses = useClasses
	s.classPrefixes = classPrefixes
	return nil
}

// CanUseItemByClassID checks eligibility using the official ClassName class ID.
func (s *Store) CanUseItemByClassID(itemID uint16, classID byte) (bool, bool) {
	if s == nil {
		return false, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.itemsByID[itemID]
	if !ok {
		return false, false
	}
	prefix, ok := s.classPrefixes[classID]
	if !ok {
		return false, false
	}
	classes, ok := s.useClasses[info.UseClass]
	if !ok {
		return false, false
	}
	allowed, ok := classes[strings.ToLower(prefix)]
	return allowed, ok
}

// LoadItemUseEffectSHN loads official consumable recovery/buff definitions from ItemUseEffect.shn.
func (s *Store) LoadItemUseEffectSHN(filePath string) error {
	file, err := shn.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	effects := make(map[string]*ItemUseEffect, file.RecordCount)
	for _, row := range file.Rows {
		idx := strings.TrimSpace(row.GetString("ItemIndex"))
		if idx == "" {
			continue
		}
		effects[strings.ToLower(idx)] = &ItemUseEffect{
			ItemIndex:       idx,
			UseEffectA:      uint32(row.GetInt("UseEffectA")),
			UseValueA:       row.GetUint16("UseValueA"),
			UseEffectB:      uint32(row.GetInt("UseEffectB")),
			UseValueB:       row.GetUint16("UseValueB"),
			UseEffectC:      uint32(row.GetInt("UseEffectC")),
			UseValueC:       row.GetUint16("UseValueC"),
			UseAbstateName:  strings.TrimSpace(row.GetString("UseAbstateName")),
			UseAbstateIndex: uint32(row.GetInt("UseAbstateIndex")),
		}
	}

	s.mu.Lock()
	s.useEffects = effects
	s.mu.Unlock()

	return nil
}

// GetItemInfo returns the full ItemInfo for a given itemID.
func (s *Store) GetItemInfo(itemID uint16) (*ItemInfo, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.itemsByID[itemID]
	return info, ok
}

// GetItemInfoByName returns the full ItemInfo for a given InxName (case-insensitive).
func (s *Store) GetItemInfoByName(inxName string) (*ItemInfo, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.itemsByName[strings.ToLower(strings.TrimSpace(inxName))]
	return info, ok
}

// GetItemUseEffect returns the consumable recovery/buff effect for an item InxName (case-insensitive).
func (s *Store) GetItemUseEffect(inxName string) (*ItemUseEffect, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	eff, ok := s.useEffects[strings.ToLower(strings.TrimSpace(inxName))]
	return eff, ok
}

// GetItemClass returns the official ItemClass for a given itemID.
func (s *Store) GetItemClass(itemID uint16) (uint8, bool) {
	if s == nil {
		return 0, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	cls, ok := s.itemClasses[itemID]
	return cls, ok
}

// GetItemAttrSize returns the base binary attribute payload size in bytes for an item (0 options).
func (s *Store) GetItemAttrSize(itemID uint16) int {
	cls, ok := s.GetItemClass(itemID)
	if !ok {
		return 4 // default fallback: Amount (4 bytes)
	}
	return CalcAttrSizeForClass(cls, 0)
}

// BuildItemAttr serializes the dynamic attribute byte slice for an item emulating iac_fill2buffer.
func (s *Store) BuildItemAttr(itemID uint16, amount uint32, opts []ItemStatOption, innerCount int) []byte {
	cls, ok := s.GetItemClass(itemID)
	if !ok {
		cls = 24 // default fallback: Class 24 (Amount, 4 bytes)
	}
	return BuildItemAttrBuffer(cls, amount, opts, innerCount)
}
