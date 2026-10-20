with open('pkg/data/item/item.go', 'rb') as f:
    content = f.read().decode('utf-8')

# 1. Update imports
old_import = '''import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)'''

new_import = '''import (
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
}'''

# 2. Update Store struct and NewStore
old_store = '''type Store struct {
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
}'''

new_store = '''type Store struct {
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
}'''

# 3. Add methods at end of file
new_methods = '''

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
'''

norm = content.replace('\r\n', '\n')
assert old_import in norm, "old_import not found"
assert old_store in norm, "old_store not found"

norm = norm.replace(old_import, new_import)
norm = norm.replace(old_store, new_store)
norm = norm + new_methods

with open('pkg/data/item/item.go', 'wb') as f:
    f.write(norm.replace('\n', '\r\n').encode('utf-8'))

print("Updated pkg/data/item/item.go!")
