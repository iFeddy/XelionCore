package zone

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ============================================================================
// Official Zone.exe KingdomQuest Constants (Sub-phase 27C — 1:1 PDB/EXE Port)
// Every literal below is backed by a verified Zone.exe virtual address.
// ============================================================================

const (
	kqContainerMaxElements   = 0x12C              // exe: 0x00499A6D (KingdomQuestContainer ctor push 0x12c = 300)
	kqEntranceRingSize       = 0x80               // exe: 0x004950F7 (KingdomQuestEntrance ring capacity = 128)
	kqEntranceRingMask       = 0x7F               // exe: 0x004961B7 (and edx, [esi+0xfe04] -> 128-1)
	kqEntranceMaxJoiners     = 0x64               // exe: 0x004961BD (kqe_newentrance cap = 100 joiners)
	kqMapLinkSlots           = 4                  // exe: 0x00497C89 / 0x004956A9 (4 MapLink entries per KQElement)
	kqTicksPerMinute         = 0x258              // exe: 0x00495449 (kqe_QuestStart imul edx, 0x258 = 600 ticks/min @ 10Hz)
	kqTicksPerSecond         = 10                 // exe: 0x00499699 (lea ecx, [ecx+ecx*4]; lea ecx, [esi+ecx*2] -> *10)
	kqSecondsPerMinute       = 60                 // exe: 0x00495449 (0x258 / 10Hz = 60s per minute)
	kqReviveUnlimited        = 0x1869F            // exe: 0x00496BEB (kqe_ReviveNumber & kqpil_StorePlayer_GetKnockDownNumber = 99999)
	kqBattleStatDefault      = 0x2710             // exe: 0x00498FF7 (kqpil_InsertPlayerInfo mov [eax+0x30], 0x2710 = 10000)
	kqClearObjectMaskMake    = 0x3B6              // exe: 0x004980D3 (wms_NC_KQ_W2Z_MAKE_REQ fm_ClearObject(0x3b6))
	kqClearObjectMaskDestroy = 0x3B0              // exe: 0x0049736C (wms_NC_KQ_W2Z_DESTROY_CMD fm_ClearObject(0x3b0))
	kqEntryStateValidMask    = uint16(0x56)       // exe: 0x00496A55 (sp_NC_KQ_ENTRYRESPONCE_ACK bitmask 0x56 = states 1,2,4,6)
	kqInvalidHandle          = uint32(0xFFFFFFFF) // exe: 0x00497383 (FieldMap+0x3C = -1 when unbound)
	kqInvalidID              = uint16(0xFFFF)     // exe: 0x00497386 (FieldMap+0x40 = 0xFFFF when unbound)
	kqTeamMax                = uint8(2)           // exe: 0x004966C6 / 0x0055D75D (KQTT_MAX = 2)

	// Official Roumen fallback destination for ban/escape.
	kqRoumenFallbackMap = "RouN" // exe: 0x004971C2 / 0x0043D04D
	kqRoumenFallbackX   = 0x192D // exe: 0x004971E0 / 0x0043D06E (6445)
	kqRoumenFallbackY   = 0x21B6 // exe: 0x004971E7 / 0x0043D078 (8630)

	// Official Z2W_MAKE_ACK error codes (PROTO_NC_KQ_Z2W_MAKE_ACK = 0x580E).
	kqErrMakeNoErr          = uint16(0x0981) // exe: 0x00498AE9 (ERR_KINGDOMQUEST_MAKE_NOERR)
	kqErrMakeDuplicate      = uint16(0x0982) // exe: 0x00497B9E (ERR_KINGDOMQUEST_MAKE_DUPLICATE)
	kqErrMakeTooMany        = uint16(0x0983) // exe: 0x00498813 (ERR_KINGDOMQUEST_MAKE_TOOMANY)
	kqErrMakeScriptNotFound = uint16(0x098C) // exe: 0x00497C60 (ERR_KINGDOMQUEST_MAKE_SCRIPTNOTFOUND)

	// Official Item pickup error code when KQItem.shn PickupLimit is exceeded.
	kqErrPickupLimit = uint16(0x034A) // exe: 0x0052FC87 (so_ply_PickupItem push 0x34a)

	// Official Inventory box IDs passed to sp_RemoveKQItem (exe: 0x00443E4F / 0x00443E58).
	kqRemoveBoxEquip = 8 // exe: 0x0050F978 (sub eax, 8 -> Equipped box)
	kqRemoveBoxInven = 9 // exe: 0x0050F97D (dec eax -> Inventory box)

	// Official AmpersandCommand::ac_Store admin level for KQ GM commands (0x64 = 100).
	kqGMAdminLevel = byte(0x64) // exe: 0x00426533 / 0x004265E9 / 0x0042660E / 0x00426633
)

// Official static lookup tables extracted directly from Zone.exe .data.
var (
	// exe: 0x0083C788 (4 x uint16: map not found per MapLink slot 0..3)
	kqErrMakeNoMap = [kqMapLinkSlots]uint16{0x0984, 0x0985, 0x0986, 0x0987}

	// exe: 0x0083C790 (4 x uint16: submap in use per MapLink slot 0..3)
	kqErrMakeMapInUse = [kqMapLinkSlots]uint16{0x0988, 0x0989, 0x098A, 0x098B}

	// exe: 0x0083C798 (5 x uint32: GateKeeper retry wave delays in seconds)
	kqGateKeeperDelaySec = [5]int{0, 2, 5, 7, 10}
)

// ============================================================================
// DataStore Official KQ Data (MapInfo.shn + KQItemDataBox + 7 KQ SHNs)
// ============================================================================

// KQMapInfoEntry holds the official MapInfo.shn fields read by Zone.exe KQ code.
// exe: 0x00497E23 (MapInfo BTree lookup), +0x32 RegenX, +0x36 RegenY, +0x3A KingdomMap.
type KQMapInfoEntry struct {
	ID         uint16
	MapName    string
	RegenX     uint32 // MapInfo.shn +0x32 (exe: 0x004995D4 / 0x00496B2C)
	RegenY     uint32 // MapInfo.shn +0x36 (exe: 0x004995DD / 0x00496B38)
	KingdomMap uint8  // MapInfo.shn +0x3A
}

// KQItemDataBox ports Zone.exe KQItemDataBox (gKQItemDataBox @ 0x13388848).
type KQItemDataBox struct {
	mu     sync.RWMutex
	byID   map[uint16]kq.KQItem
	byName map[string]kq.KQItem
}

// NewKQItemDataBox constructs an empty KQItemDataBox (exe: 0x0064C020).
func NewKQItemDataBox() *KQItemDataBox {
	return &KQItemDataBox{
		byID:   make(map[uint16]kq.KQItem),
		byName: make(map[string]kq.KQItem),
	}
}

// kqidb_Load ports KQItemDataBox::kqidb_Load (exe: 0x0064BD80).
// Indexes KQItem.shn records by ItemID (resolved via ItemInfo.shn InxName).
func (b *KQItemDataBox) kqidb_Load(items []kq.KQItem, lookupID func(inxName string) (uint16, bool)) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.byID = make(map[uint16]kq.KQItem, len(items))
	b.byName = make(map[string]kq.KQItem, len(items))
	loaded := 0
	for _, it := range items {
		inx := strings.TrimSpace(it.ItemIndexString())
		if inx == "" {
			continue
		}
		b.byName[strings.ToLower(inx)] = it
		if lookupID != nil {
			if id, ok := lookupID(inx); ok {
				b.byID[id] = it
				loaded++
			}
		}
	}
	return loaded
}

// kqidb_GetData ports KQItemDataBox::kqidb_GetData (exe: 0x0064BC60).
func (b *KQItemDataBox) kqidb_GetData(itemID uint16) (*kq.KQItem, bool) {
	if b == nil {
		return nil, false
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	it, ok := b.byID[itemID]
	if !ok {
		return nil, false
	}
	cp := it
	return &cp, true
}

// kqidb_CanPickUp ports KQItemDataBox::kqidb_CanPickUp (exe: 0x0064BC90).
// Returns true if itemID is not in KQItem.shn, PickupLimit == 0, or totalCount <= PickupLimit.
func (b *KQItemDataBox) kqidb_CanPickUp(itemID uint16, totalCount uint64) bool {
	if b == nil {
		return true
	}
	b.mu.RLock()
	it, ok := b.byID[itemID]
	b.mu.RUnlock()
	if !ok {
		return true // exe: 0x0064BCAF (not in map -> allow)
	}
	if it.PickupLimit == 0 {
		return true // exe: 0x0064BCBB..0x0064BCBE (limit 0 -> unlimited)
	}
	return totalCount <= uint64(it.PickupLimit) // exe: 0x0064BCC0..0x0064BCCD
}

// LoadKQOfficialData loads the 7 official KQ SHN tables, KQItemDataBox, and MapInfo.shn index.
func (ds *DataStore) LoadKQOfficialData(dataDir string) {
	if ds == nil {
		return
	}
	store, err := kq.LoadDir(dataDir)
	if err != nil {
		log.Printf("[Zone][KQ] Warning: kq.LoadDir(%s) failed: %v", dataDir, err)
	} else {
		ds.KQStore = store
		ds.KQItemBox = NewKQItemDataBox()
		ds.KQItemBox.kqidb_Load(store.Items, func(inxName string) (uint16, bool) {
			if ds.ItemStore == nil {
				return 0, false
			}
			if info, ok := ds.ItemStore.GetItemInfoByName(inxName); ok && info != nil {
				return info.ID, true
			}
			return 0, false
		})
	}

	ds.KQMapInfo = make(map[string]*KQMapInfoEntry)
	if mf := ds.GetSHN("MapInfo"); mf != nil {
		for i := range mf.Rows {
			r := &mf.Rows[i]
			name := strings.TrimSpace(r.GetString("MapName"))
			if name == "" {
				continue
			}
			ds.KQMapInfo[strings.ToLower(name)] = &KQMapInfoEntry{
				ID:         r.GetUint16("ID"),
				MapName:    name,
				RegenX:     uint32(r.GetInt("RegenX")),
				RegenY:     uint32(r.GetInt("RegenY")),
				KingdomMap: r.GetByte("KingdomMap"),
			}
		}
	}
}

// GetKQMapInfo looks up a MapInfo.shn entry by map name (case-insensitive).
func (ds *DataStore) GetKQMapInfo(mapName string) (*KQMapInfoEntry, bool) {
	if ds == nil || ds.KQMapInfo == nil {
		return nil, false
	}
	e, ok := ds.KQMapInfo[strings.ToLower(strings.TrimSpace(mapName))]
	return e, ok
}

// ============================================================================
// KingdomQuest::KQPlayerInform & KQPlayerInfoList (Zone.exe 0x00495250..0x004990F7)
// ============================================================================

// KQPlayerInform ports KingdomQuest::KQPlayerInform (sizeof = 0x36 in Zone.pdb).
type KQPlayerInform struct {
	RegNum         uint32 // +0x00 (nCharNo)
	KnockDownNum   int32  // +0x04
	BattleStatBase uint32 // +0x30 (initialized to 0x2710 = 10000, exe: 0x00498FF7)
	IsJoiner       uint8  // +0x34
	TeamType       uint8  // +0x35
}

// kqpi_KnockDownMinus ports KQPlayerInform::kqpi_KnockDownMinus (exe: 0x00496D70).
func (p *KQPlayerInform) kqpi_KnockDownMinus(_ uint32) int {
	p.KnockDownNum-- // exe: 0x00496DB2
	return int(p.KnockDownNum)
}

// kqpi_KnockDownPlus ports KQPlayerInform::kqpi_KnockDownPlus (exe: 0x00496DE0).
func (p *KQPlayerInform) kqpi_KnockDownPlus(_ uint32) int {
	p.KnockDownNum++ // exe: 0x00496E22
	return int(p.KnockDownNum)
}

// kqpi_KnockDownGet ports KQPlayerInform::kqpi_KnockDownGet (exe: 0x00496E50).
func (p *KQPlayerInform) kqpi_KnockDownGet(_ uint32) int {
	return int(p.KnockDownNum) // exe: 0x00496E8E
}

// kqpi_GetTeamType ports KQPlayerInform::kqpi_GetTeamType (exe: 0x00495260).
func (p *KQPlayerInform) kqpi_GetTeamType(_ uint32) int {
	return int(p.TeamType) // exe: 0x00495260 (movzx eax, byte [ecx+0x35])
}

// KQPlayerInfoList ports KingdomQuest::KQPlayerInfoList (exe: 0x004999E0).
type KQPlayerInfoList struct {
	mu      sync.RWMutex
	players map[uint32]*KQPlayerInform
}

func (l *KQPlayerInfoList) clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.players = make(map[uint32]*KQPlayerInform)
}

// kqpil_StorePlayer_GetKnockDownNumber ports KQPlayerInfoList::kqpil_StorePlayer_GetKnockDownNumber (exe: 0x00498BA0).
// Returns KnockDownNum if found, or 0x1869F (99999) if not found (exe: 0x00498C18).
func (l *KQPlayerInfoList) kqpil_StorePlayer_GetKnockDownNumber(regNum uint32) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.players != nil {
		if info, ok := l.players[regNum]; ok && info != nil {
			return int(info.KnockDownNum) // exe: 0x00498C13
		}
	}
	return kqReviveUnlimited // exe: 0x00498C18 (mov eax, 0x1869f)
}

// kqpil_PlayerInformFunc ports KQPlayerInfoList::kqpil_PlayerInformFunc (exe: 0x00498C40).
func (l *KQPlayerInfoList) kqpil_PlayerInformFunc(fn func(*KQPlayerInform, uint32) int, regNum uint32) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.players == nil {
		return 0
	}
	info, ok := l.players[regNum]
	if !ok || info == nil {
		return 0 // exe: 0x00498D11
	}
	return fn(info, regNum) // exe: 0x00498D16
}

// kqpil_IsKQJoiner ports KQPlayerInfoList::kqpil_IsKQJoiner (exe: 0x00498E90).
func (l *KQPlayerInfoList) kqpil_IsKQJoiner(regNum uint32) uint8 {
	if regNum == kqInvalidHandle { // exe: 0x00498EBB
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.players == nil {
		return 0
	}
	info, ok := l.players[regNum]
	if !ok || info == nil {
		return 0 // exe: 0x00498F0A
	}
	return info.IsJoiner // exe: 0x00498F0C
}

// kqpil_InsertPlayerInfo ports KQPlayerInfoList::kqpil_InsertPlayerInfo (exe: 0x00498F30).
func (l *KQPlayerInfoList) kqpil_InsertPlayerInfo(regNum uint32, teamType uint8, isJoiner uint8) bool {
	if regNum == kqInvalidHandle { // exe: 0x00498F60
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.players == nil {
		l.players = make(map[uint32]*KQPlayerInform)
	}
	if _, exists := l.players[regNum]; exists {
		return true // exe: 0x00498FB3 -> 0x00499004 (already present -> returns 1)
	}
	l.players[regNum] = &KQPlayerInform{
		RegNum:         regNum,
		KnockDownNum:   0,                   // exe: 0x00498FD9
		BattleStatBase: kqBattleStatDefault, // exe: 0x00498FF7
		IsJoiner:       isJoiner,            // exe: 0x00498FFE
		TeamType:       teamType,            // exe: 0x00499001
	}
	return true
}

// kqpil_DeletePlayerInfo ports KQPlayerInfoList::kqpil_DeletePlayerInfo (exe: 0x00499030).
func (l *KQPlayerInfoList) kqpil_DeletePlayerInfo(regNum uint32) bool {
	if regNum == kqInvalidHandle { // exe: 0x0049905E
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.players == nil {
		return false
	}
	if _, ok := l.players[regNum]; !ok {
		return false // exe: 0x004990B1 -> 0x004990D8
	}
	delete(l.players, regNum)
	return true // exe: 0x004990DC
}

// ============================================================================
// KingdomQuest::KQElement & KingdomQuestContainer (Zone.exe 0x004953E0..0x00499A40)
// ============================================================================

// KQElementMap ports KingdomQuest::KQElement::KQElementMap (sizeof = 0x24, exe: 0x004194D0).
type KQElementMap struct {
	MapInfo  kq.PROTO_KQ_MAP_INFO // +0x00 (26 bytes)
	FieldMap *ZoneMap             // +0x1C (KQElement+0x19C for slot 0)
	RegenX   uint32               // MapInfo+0x32 (exe: 0x004995D4)
	RegenY   uint32               // MapInfo+0x36 (exe: 0x004995DD)
}

// KQElement ports KingdomQuest::KQElement (sizeof = 0x42C, exe: 0x00422E40).
type KQElement struct {
	EndTick         uint64                       // +0x00 (exe: 0x0049545C)
	Info            kq.PROTO_KQ_INFO             // +0x04 (377 bytes)
	Maps            [kqMapLinkSlots]KQElementMap // +0x180 (4 * 0x24 = 0x90 bytes)
	ScriptLanguage  string                       // +0x210
	ScriptInitValue string                       // +0x310
	PlayerInfoList  KQPlayerInfoList             // +0x410
}

// kqe_QuestStart ports KQElement::kqe_QuestStart (exe: 0x004953E0).
func (e *KQElement) kqe_QuestStart(curTick uint64, now time.Time) {
	if e == nil {
		return
	}
	// exe: 0x00495445..0x0049545C: EndTick = curTick + LimitTime * 0x258
	e.EndTick = curTick + uint64(e.Info.LimitTime)*kqTicksPerMinute
	for i := 0; i < kqMapLinkSlots; i++ {
		fm := e.Maps[i].FieldMap
		if fm == nil {
			continue
		}
		// exe: 0x00495482 cc_DropFilm
		fm.cc_DropFilm(e.ScriptLanguage)
		// exe: 0x00495491 mda_CloseAllDoor
		fm.mda_CloseAllDoor()
		// exe: 0x004954AF cc_PlayFilm
		fm.cc_PlayFilm(e.ScriptLanguage, e.ScriptInitValue, now)
	}
}

// kqe_ReviveNumber ports KQElement::kqe_ReviveNumber (exe: 0x00496B90, jmptbl @ 0x00496C34).
// Returns max allowed knockdowns and outFlag (0 or 1) according to PlayerRevivalMode (1..4).
func (e *KQElement) kqe_ReviveNumber() (maxRevives int, outFlag uint8) {
	if e == nil {
		return 0, 0
	}
	// exe: 0x00496BD4..0x00496C15 (PlayerRevivalMode @ +0x82, PlayerRevivalCount @ +0x83)
	switch e.Info.PlayerRevivalMode {
	case 1: // exe: 0x00496BEB
		return kqReviveUnlimited, 0
	case 2: // exe: 0x00496BF5
		return int(e.Info.PlayerRevivalCount), 0
	case 3: // exe: 0x00496C01
		return kqReviveUnlimited, 1
	case 4: // exe: 0x00496C0A
		return int(e.Info.PlayerRevivalCount), 1
	default: // exe: 0x00496C15
		return 0, 0
	}
}

// kqe_GetKQLimitMin ports KQElement::kqe_GetKQLimitMin (exe: 0x005E4D20).
func (e *KQElement) kqe_GetKQLimitMin() uint16 {
	if e == nil {
		return 0
	}
	return e.Info.LimitTime // exe: 0x005E4D5E (word [esi+0x4D])
}

// KingdomQuestContainer ports KingdomQuest::KingdomQuestContainer (kqcontainer @ 0x0D624A00).
// Capacity is fixed at 0x12C (300) elements (exe: 0x00499A6D).
type KingdomQuestContainer struct {
	mu       sync.RWMutex
	slots    [kqContainerMaxElements]*KQElement
	capacity int
}

// NewKingdomQuestContainer constructs a KingdomQuestContainer with capacity 300 (exe: 0x00499A40).
func NewKingdomQuestContainer() *KingdomQuestContainer {
	return &KingdomQuestContainer{
		capacity: kqContainerMaxElements, // exe: 0x00499A6D
	}
}

// kqc_SearchHandle ports KingdomQuestContainer::kqc_SearchHandle (exe: 0x004954F0).
func (c *KingdomQuestContainer) kqc_SearchHandle(handle uint32) (*KQElement, uint16, bool) {
	if c == nil {
		return nil, kqInvalidID, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := 0; i < c.capacity; i++ {
		if el := c.slots[i]; el != nil && el.Info.Handle == handle { // exe: 0x00495594
			return el, uint16(i), true
		}
	}
	return nil, kqInvalidID, false
}

// kqc_MapUseCheck ports KingdomQuestContainer::kqc_MapUseCheck (exe: 0x00495710)
// and ListIteratorKingdomQuestMapFind::lid_Call (exe: 0x00495620).
// Returns true if mapName is NOT currently used by any active KQElement; false if already in use.
func (c *KingdomQuestContainer) kqc_MapUseCheck(mapName string) bool {
	if c == nil {
		return true
	}
	target := strings.TrimSpace(mapName)
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := 0; i < c.capacity; i++ {
		el := c.slots[i]
		if el == nil {
			continue
		}
		// lid_Call (exe: 0x00495680..0x004956AC): check all 4 KQElementMap slots
		for m := 0; m < kqMapLinkSlots; m++ {
			fm := el.Maps[m].FieldMap
			if fm != nil && strings.EqualFold(fm.MapID, target) {
				return false // exe: 0x004956FB (in use -> return 0)
			}
		}
	}
	return true // exe: 0x004956CA (free -> return 1)
}

func (c *KingdomQuestContainer) allocLocked(el *KQElement) (uint16, bool) {
	for i := 0; i < c.capacity; i++ {
		if c.slots[i] == nil {
			c.slots[i] = el
			return uint16(i), true
		}
	}
	return kqInvalidID, false
}

func (c *KingdomQuestContainer) freeSlot(slotIdx uint16) bool {
	if c == nil || int(slotIdx) >= c.capacity {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.slots[slotIdx] == nil {
		return false
	}
	c.slots[slotIdx] = nil
	return true
}

func (c *KingdomQuestContainer) count() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	n := 0
	for i := 0; i < c.capacity; i++ {
		if c.slots[i] != nil {
			n++
		}
	}
	return n
}

// ============================================================================
// KingdomQuestEntrance (Zone.exe 0x004950F0, 0x00495140, 0x004960A0, 0x004960F0, 0x00496150)
// ============================================================================

// KQJoinerEntry models PROTO_NC_KQ_JOINER (5 bytes: chrregnum[4], nTeamType[1]).
type KQJoinerEntry struct {
	CharRegNum uint32 // +0x00
	TeamType   uint8  // +0x04
}

type kqEntranceSlot struct {
	Total   int             // +0x00
	Entered int             // +0x04
	Joiners []KQJoinerEntry // +0x08 (max 0x64 = 100, exe: 0x004961BD)
}

// KingdomQuestEntrance ports Zone.exe KingdomQuestEntrance (kqentrance @ 0x0083C818).
type KingdomQuestEntrance struct {
	mu     sync.Mutex
	slots  [kqEntranceRingSize]kqEntranceSlot // 0x80 slots, stride 0x1FC
	cursor uint32                             // +0xFE00, masked with 0x7F (+0xFE04)
}

// kqe_Compare ports KingdomQuestEntrance::kqe_Compare (exe: 0x00495140).
// Returns int(b.CharRegNum) - int(a.CharRegNum) (descending order by chrregnum).
func kqe_Compare(a, b KQJoinerEntry) int {
	if b.CharRegNum > a.CharRegNum {
		return 1
	}
	if b.CharRegNum < a.CharRegNum {
		return -1
	}
	return 0
}

// kqe_newentrance ports KingdomQuestEntrance::operator[] / kqe_newentrance (exe: 0x00496150).
func (e *KingdomQuestEntrance) kqe_newentrance(joiners []KQJoinerEntry) uint16 {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cursor = (e.cursor + 1) & kqEntranceRingMask // exe: 0x004961B6..0x004961C2
	count := len(joiners)
	if count > kqEntranceMaxJoiners { // exe: 0x004961BD..0x004961EE
		count = kqEntranceMaxJoiners
	}
	cp := make([]KQJoinerEntry, count)
	copy(cp, joiners[:count])
	e.slots[e.cursor] = kqEntranceSlot{
		Total:   count, // exe: 0x00496202
		Entered: 0,     // exe: 0x00496211
		Joiners: cp,    // exe: 0x00496235
	}
	return uint16(e.cursor)
}

// kqe_entry ports KingdomQuestEntrance::kqe_entry (exe: 0x004950F0).
func (e *KingdomQuestEntrance) kqe_entry(entryIdx uint16) {
	if entryIdx >= kqEntranceRingSize { // exe: 0x004950F7..0x00495103
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.slots[entryIdx].Entered++ // exe: 0x00495125
}

// kqe_allentry ports the inlined KingdomQuestEntrance::kqe_allentry (exe: 0x00496392..0x004963CD).
func (e *KingdomQuestEntrance) kqe_allentry(entryIdx uint16) bool {
	if entryIdx >= kqEntranceRingSize { // exe: 0x00496396..0x004963A2
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	s := &e.slots[entryIdx]
	return s.Entered >= s.Total // exe: 0x004963C2..0x004963C8
}

// kqe_sort ports KingdomQuestEntrance::kqe_sort (exe: 0x004960A0).
func (e *KingdomQuestEntrance) kqe_sort(entryIdx uint16) {
	if entryIdx >= kqEntranceRingSize { // exe: 0x004960A7..0x004960B6
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sortLocked(entryIdx)
}

func (e *KingdomQuestEntrance) sortLocked(entryIdx uint16) {
	s := &e.slots[entryIdx]
	sort.Slice(s.Joiners, func(i, j int) bool {
		return kqe_Compare(s.Joiners[i], s.Joiners[j]) < 0
	})
}

// kqe_search ports KingdomQuestEntrance::kqe_search (exe: 0x004960F0).
// Binary-searches the descending-sorted joiners array using kqe_Compare.
func (e *KingdomQuestEntrance) kqe_search(charRegNum uint32, entryIdx uint16) (*KQJoinerEntry, bool) {
	if entryIdx >= kqEntranceRingSize { // exe: 0x004960F7..0x00496106
		return nil, false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	s := &e.slots[entryIdx]
	key := KQJoinerEntry{CharRegNum: charRegNum}
	lo, hi := 0, len(s.Joiners)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		cmp := kqe_Compare(key, s.Joiners[mid])
		if cmp == 0 {
			return &s.Joiners[mid], true
		}
		if cmp < 0 {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return nil, false
}

// ============================================================================
// KingdomQuestGateKeeper (Zone.exe 0x00495150, 0x00495160, 0x00496280)
// ============================================================================

// KingdomQuestGateKeeper ports KingdomQuestGate::KingdomQuestGateKeeper (sizeof = 0x70, exe: 0x00495150).
type KingdomQuestGateKeeper struct {
	EntryIdx  uint16    // +0x04
	MapName   string    // +0x1C
	MapBase   string    // +0x3D
	RegenX    uint32    // +0x5E
	RegenY    uint32    // +0x62
	KQID      uint16    // +0x6A
	KQHandle  uint32    // +0x6C
	TriggerAt time.Time // scheduled time in zonevent (exe: 0x00499697..0x004996A2)
}

// erb_CheckLinkMap ports KingdomQuestGateKeeper::erb_CheckLinkMap (exe: 0x00495160).
// Validates whether the player can leave their current map into the KQ map; if not,
// sends KQEntryMapErr0 and KQEntryMapErr1 notices and returns false.
func (g *KingdomQuestGateKeeper) erb_CheckLinkMap(h *Handler, player *ZonePlayer) bool {
	if player == nil || player.CurrentMap == nil {
		return false
	}
	if player.CurrentMap.FieldEntry != nil && !player.CurrentMap.FieldEntry.Teleport.LinkOut {
		// exe: 0x00495180 / 0x004951DC (KQEntryMapErr0, KQEntryMapErr1)
		h.sendPlayerNotice(player, "KQEntryMapErr0")
		h.sendPlayerNotice(player, "KQEntryMapErr1")
		return false // exe: 0x00495232
	}
	return true // exe: 0x0049523B
}

// erb_happen ports KingdomQuestGateKeeper::erb_happen (exe: 0x00496280).
func (g *KingdomQuestGateKeeper) erb_happen(h *Handler) int {
	if h == nil {
		return 0
	}
	st := h.kqState()
	// exe: 0x00496392..0x004963CD: kqe_allentry early exit
	if st.entrance.kqe_allentry(g.EntryIdx) {
		return 0
	}
	// exe: 0x00496472..0x004964B6: initial kqe_sort
	st.entrance.kqe_sort(g.EntryIdx)

	// exe: 0x004964B9..0x004964EC: verify KQElement still exists
	elem, _, ok := st.container.kqc_SearchHandle(g.KQHandle)
	if !ok || elem == nil {
		return 0
	}

	h.mu.RLock()
	online := make([]*ZonePlayer, 0, len(h.players))
	for _, p := range h.players {
		if p != nil && p.Char != nil {
			online = append(online, p)
		}
	}
	h.mu.RUnlock()

	enteredCount := 0
	for _, player := range online {
		// exe: 0x00496555..0x0049656D: skip mode == 0 or mode == 3 (Dead)
		if player.State == 0 || player.State == 3 {
			continue
		}
		regNum := uint32(player.Char.ID)
		if regNum == kqInvalidHandle { // exe: 0x00496581
			continue
		}
		// exe: 0x004965B2..0x004965D5: sort & binary search joiner
		st.entrance.kqe_sort(g.EntryIdx)
		joiner, found := st.entrance.kqe_search(regNum, g.EntryIdx)
		if !found || joiner == nil {
			continue
		}
		// exe: 0x004965DE: erb_CheckLinkMap
		if !g.erb_CheckLinkMap(h, player) {
			continue
		}
		// exe: 0x004965F7..0x00496674: mode == 5 (Booth) -> KQEntryInCamp0 / KQEntryInCamp1
		if player.State == 5 || player.ActiveBooth != nil {
			h.sendPlayerNotice(player, "KQEntryInCamp0")
			h.sendPlayerNotice(player, "KQEntryInCamp1")
			continue
		}
		// exe: 0x00496685..0x00496693: mode == 4 (Resting/House) -> stand up
		if player.State == 4 {
			player.State = 1
		}
		// exe: 0x0049669F: kqe_entry
		st.entrance.kqe_entry(g.EntryIdx)

		// exe: 0x004966B2: CT_KQJoin (title hook)
		h.UnlockOrAdvanceTitle(player, kqTitleTypeKQJoin, 1, false)
		// exe: 0x004966BD..0x00496722: determine warp destination (TeamRegenXY if teamType < 2)
		destX := int(g.RegenX)
		destY := int(g.RegenY)
		teamType := joiner.TeamType
		if teamType < kqTeamMax { // exe: 0x004966C6..0x00496722
			destX = int(elem.Info.TeamRegenXY[teamType].X)
			destY = int(elem.Info.TeamRegenXY[teamType].Y)
		}

		// Record previous map for rollback + GameLog 0x0441 (exe: 0x00496730..0x004967DA)
		prevMap := ""
		prevX, prevY := player.Char.XPos, player.Char.YPos
		if player.CurrentMap != nil {
			prevMap = player.CurrentMap.MapID
		}
		h.SavePreKQReturnPos(player, prevMap, prevX, prevY)

		// exe: 0x004966D8 / 0x00496722: MapLink into KQ map
		h.WarpPlayer(player, g.MapName, destX, destY, 0)
		// exe: 0x004542F0 / Character.exe 0x00459E90: persistir posición normal + KQMap
		h.so_SaveLocation(player)
		// exe: 0x00496740..0x004967DA: emitir PROTO_NC_LOG_GAME_KQ_ENTER (0x0441)
		h.EmitKQEnterGameLog(player, prevMap, prevX, prevY, g.KQID, g.KQHandle)

		// exe: 0x004967E6: mark joiner consumed (chrregnum = 0xFFFFFFFF)
		st.entrance.mu.Lock()
		joiner.CharRegNum = kqInvalidHandle
		st.entrance.mu.Unlock()
		enteredCount++
	}
	return enteredCount
}

// ============================================================================
// KQContributeList (Zone.exe 0x0049A770, 0x0049A500)
// ============================================================================

// KQContributeList ports KQContribute::KQContributeList (g_KQContributeList @ 0x0D67A048).
type KQContributeList struct {
	mu     sync.RWMutex
	active map[uint32]map[uint32]int // kqHandle -> (charRegNum -> mobKills)
}

func NewKQContributeList() *KQContributeList {
	return &KQContributeList{
		active: make(map[uint32]map[uint32]int),
	}
}

// kqcl_NewKQ ports KQContributeList::kqcl_NewKQ (exe: 0x0049A770).
func (c *KQContributeList) kqcl_NewKQ(handle uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active == nil {
		c.active = make(map[uint32]map[uint32]int)
	}
	if _, exists := c.active[handle]; !exists {
		c.active[handle] = make(map[uint32]int)
	}
}

// kqpl_EndOfKQ ports KQContributeList::kqpl_EndOfKQ (exe: 0x0049A500).
func (c *KQContributeList) kqpl_EndOfKQ(handle uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.active, handle)
}

// ============================================================================
// ZoneKQRuntime — Runtime container attached to zone.Handler
// ============================================================================

// ZoneKQRuntime holds all Zone.exe global KQ singletons scoped to a Handler instance.
type ZoneKQRuntime struct {
	mu            sync.Mutex
	container     *KingdomQuestContainer
	entrance      *KingdomQuestEntrance
	contribute    *KQContributeList
	rewardTracker *KQRewardTracker
	gateKeepers   []*KingdomQuestGateKeeper
	nowFunc       func() time.Time
	bootTime      time.Time
	reloaded      map[string]int // tracks &Reroadkingdomquest reloads for diagnostics/tests
}

func newZoneKQRuntime() *ZoneKQRuntime {
	now := time.Now()
	return &ZoneKQRuntime{
		container:     NewKingdomQuestContainer(),
		entrance:      &KingdomQuestEntrance{},
		contribute:    NewKQContributeList(),
		rewardTracker: newKQRewardTracker(),
		nowFunc:       time.Now,
		bootTime:      now,
		reloaded:      make(map[string]int),
	}
}

func (r *ZoneKQRuntime) now() time.Time {
	r.mu.Lock()
	fn := r.nowFunc
	r.mu.Unlock()
	if fn != nil {
		return fn()
	}
	return time.Now()
}

// curTick returns ZoneServer 10Hz tick count matching global curtick @ 0x14D41A70.
func (r *ZoneKQRuntime) curTick() uint64 {
	elapsed := r.now().Sub(r.bootTime)
	if elapsed < 0 {
		return 0
	}
	return uint64(elapsed.Milliseconds() * kqTicksPerSecond / 1000)
}

func (h *Handler) kqState() *ZoneKQRuntime {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.kqRuntime == nil {
		h.kqRuntime = newZoneKQRuntime()
	}
	return h.kqRuntime
}

// SetKQClock sets a deterministic clock function for KQ GateKeeper & tick testing.
func (h *Handler) SetKQClock(nowFn func() time.Time) {
	st := h.kqState()
	st.mu.Lock()
	st.nowFunc = nowFn
	if nowFn != nil {
		st.bootTime = nowFn()
	}
	st.mu.Unlock()
}

// ProcessKQGateKeepers fires all scheduled KingdomQuestGateKeeper events whose TriggerAt <= now.
// Mirrors ZoneServer event loop invoking KingdomQuestGateKeeper::erb_happen (exe: 0x00496280).
func (h *Handler) ProcessKQGateKeepers(now time.Time) int {
	st := h.kqState()
	st.mu.Lock()
	var due []*KingdomQuestGateKeeper
	var pending []*KingdomQuestGateKeeper
	for _, gk := range st.gateKeepers {
		if !gk.TriggerAt.After(now) {
			due = append(due, gk)
		} else {
			pending = append(pending, gk)
		}
	}
	st.gateKeepers = pending
	st.mu.Unlock()

	totalEntered := 0
	for _, gk := range due {
		totalEntered += gk.erb_happen(h)
	}
	return totalEntered
}

// ============================================================================
// ZoneMap KQ Helpers (fm_GetKQhandle, fm_ClearObject, cc_DropFilm/PlayFilm, mda_CloseAllDoor)
// ============================================================================

// fm_GetKQhandle ports FieldMap::fm_GetKQhandle (exe: 0x00419370).
func (m *ZoneMap) fm_GetKQhandle() uint32 {
	if m == nil {
		return kqInvalidHandle
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.KQHandle
}

// fm_ClearObject ports FieldMap::fm_ClearObject (exe: 0x00495A10).
// Called with mask 0x3B6 on W2Z_MAKE_REQ (exe: 0x004980D3) and 0x3B0 on W2Z_DESTROY_CMD (exe: 0x0049736C).
func (m *ZoneMap) fm_ClearObject(mask uint32) {
	if m == nil || mask == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for hnd := range m.mobs {
		if m.grid != nil {
			m.grid.RemoveEntity(hnd)
		}
		delete(m.mobs, hnd)
	}
	for hnd := range m.droppedItems {
		if m.grid != nil {
			m.grid.RemoveEntity(hnd)
		}
		delete(m.droppedItems, hnd)
	}
	if mask == kqClearObjectMaskMake {
		m.KQScoreRed = 0
		m.KQScoreBlue = 0
		m.KQWinRed = 0
		m.KQWinBlue = 0
	}
}

// mda_CloseAllDoor ports MapDoorArray::mda_CloseAllDoor (exe: 0x0049E1E0).
func (m *ZoneMap) mda_CloseAllDoor() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.doors != nil {
		m.doors.CloseAllDoors()
	}
}

// cc_DropFilm ports CinemaComplex::cc_DropFilm (exe: 0x00509150).
func (m *ZoneMap) cc_DropFilm(scriptName string) {
	if m == nil {
		return
	}
	shouldStop := false
	m.mu.Lock()
	if scriptName == "" || strings.EqualFold(m.KQScriptName, scriptName) {
		m.KQScriptActive = false
		m.KQScriptName = ""
		m.KQInitValue = ""
		shouldStop = true
	}
	m.mu.Unlock()
	if shouldStop {
		m.StopKQScript()
	}
}

// cc_PlayFilm ports CinemaComplex::cc_PlayFilm (exe: 0x005089A0).
func (m *ZoneMap) cc_PlayFilm(scriptName, initValue string, now time.Time) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.KQScriptActive = true
	m.KQScriptName = scriptName
	m.KQInitValue = initValue
	m.InstanceStartedAt = now
	m.mu.Unlock()
	_ = m.StartKQScript(scriptName, initValue)
}

// EnsureKQMap finds or initializes the sub-map ZoneMap (e.g. KDRockCan00) backed by baseMap (e.g. KDRockCan).
// Mirrors Zone.exe FieldContainer submap hosting from FieldList + MapInfo.shn.
func (mm *MapManager) EnsureKQMap(mapName, baseMap string) *ZoneMap {
	if mm == nil {
		return nil
	}
	normName := strings.TrimSpace(mapName)
	normBase := strings.TrimSpace(baseMap)
	if normName == "" {
		return nil
	}
	if existing := mm.GetMap(normName); existing != nil {
		return existing
	}
	if !mm.OwnsMap(normBase) && !mm.OwnsMap(normName) {
		return nil
	}

	mm.mu.Lock()
	defer mm.mu.Unlock()
	k := strings.ToLower(normName)
	if existing, ok := mm.maps[k]; ok {
		return existing
	}

	var serial uint16
	name := normName
	width, height := 20000, 20000
	var fieldEntry *field.FieldEntry
	if mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetField(normBase); ok && entry != nil {
			fieldEntry = entry
			serial = uint16(entry.Serial)
			name = entry.Name
		} else if entry, ok := mm.fieldStore.GetField(normName); ok && entry != nil {
			fieldEntry = entry
			serial = uint16(entry.Serial)
			name = entry.Name
		}
	}

	zm := NewZoneMap(serial, normName, name)
	zm.Width = width
	zm.Height = height
	zm.FieldEntry = fieldEntry
	zm.dataStore = mm.dataStore
	zm.handler = mm.handler
	if mm.dataStore != nil {
		zm.mobCatalog = mm.dataStore.MobCatalog
	}
	if mm.blockStore != nil {
		mbi := mm.blockStore.GetSHBD(normBase)
		doors := mm.blockStore.NewInstanceDoors(normBase)
		zm.SetCollisionData(mbi, doors)
	}
	mm.maps[k] = zm
	mm.loadMapNPCs(zm)
	return zm
}

// ============================================================================
// ScenarioBookShelf Script Validation (exe: 0x00507B90 sbs_GetScenarioBook)
// ============================================================================

// sbs_GetScenarioBook ports ScenarioBookShelf::sbs_GetScenarioBook (exe: 0x00507B90).
// Validates that ScriptLanguage resolves to an actual script under data/LuaScript/ or data/Script/.
func (h *Handler) sbs_GetScenarioBook(scriptLanguage string) (string, bool) {
	lang := strings.TrimSpace(strings.ReplaceAll(scriptLanguage, "\\", "/"))
	if lang == "" {
		return "", false
	}
	dataDir := ""
	if h.dataStore != nil && h.dataStore.DataDir != "" {
		dataDir = h.dataStore.DataDir
	} else {
		dataDir = resolveZoneDataDir()
	}

	parts := strings.Split(lang, "/")
	base := parts[len(parts)-1]
	baseNoDigits := strings.TrimRight(base, "0123456789")
	parent := ""
	if len(parts) > 1 {
		parent = strings.Join(parts[:len(parts)-1], "/")
	}

	candidates := []string{
		filepath.Join(dataDir, "LuaScript", filepath.FromSlash(lang)+".lua"),
		filepath.Join(dataDir, "LuaScript", filepath.FromSlash(lang), base+".lua"),
	}
	if baseNoDigits != "" && baseNoDigits != base {
		candidates = append(candidates,
			filepath.Join(dataDir, "LuaScript", filepath.FromSlash(parent), baseNoDigits, baseNoDigits+".lua"),
		)
	}
	candidates = append(candidates,
		filepath.Join(dataDir, "Script", base+".json"),
	)

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c, true
		}
	}
	return "", false
}

// ============================================================================
// S2S Payload Helpers & Outbound Z2W Packets
// ============================================================================

func decodeProtoKQInfoFromWMS(payload map[string]interface{}) (*kq.PROTO_KQ_INFO, error) {
	if payload == nil {
		return nil, errors.New("nil WMS payload")
	}
	if direct, ok := payload["kq_info"].(kq.PROTO_KQ_INFO); ok {
		return &direct, nil
	}
	if directPtr, ok := payload["kq_info"].(*kq.PROTO_KQ_INFO); ok && directPtr != nil {
		cp := *directPtr
		return &cp, nil
	}
	if rawBytes, ok := payload["kq_info_raw"].([]byte); ok && len(rawBytes) >= kq.SizeOfProtoKQInfo {
		var info kq.PROTO_KQ_INFO
		if err := info.UnmarshalBinary(rawBytes); err != nil {
			return nil, err
		}
		return &info, nil
	}
	b64 := payloadString(payload, "kq_info_b64")
	if b64 == "" {
		return nil, errors.New("missing kq_info_b64 in WMS payload")
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	var info kq.PROTO_KQ_INFO
	if err := info.UnmarshalBinary(raw); err != nil {
		return nil, err
	}
	return &info, nil
}

func decodeJoinersFromWMS(payload map[string]interface{}) []KQJoinerEntry {
	if payload == nil {
		return nil
	}
	raw, ok := payload["joiners"]
	if !ok || raw == nil {
		return nil
	}
	switch list := raw.(type) {
	case []KQJoinerEntry:
		out := make([]KQJoinerEntry, len(list))
		copy(out, list)
		return out
	case []map[string]interface{}:
		out := make([]KQJoinerEntry, 0, len(list))
		for _, item := range list {
			reg, _ := payloadInt(item, "char_reg_num")
			team, _ := payloadInt(item, "team_type")
			out = append(out, KQJoinerEntry{
				CharRegNum: uint32(reg),
				TeamType:   uint8(team),
			})
		}
		return out
	case []interface{}:
		out := make([]KQJoinerEntry, 0, len(list))
		for _, elem := range list {
			if m, ok := elem.(map[string]interface{}); ok {
				reg, _ := payloadInt(m, "char_reg_num")
				team, _ := payloadInt(m, "team_type")
				out = append(out, KQJoinerEntry{
					CharRegNum: uint32(reg),
					TeamType:   uint8(team),
				})
			}
		}
		return out
	default:
		return nil
	}
}

func decodeOutMapsFromWMS(payload map[string]interface{}) [kqMapLinkSlots]string {
	var out [kqMapLinkSlots]string
	if payload == nil {
		return out
	}
	raw, ok := payload["out_maps"]
	if !ok || raw == nil {
		return out
	}
	switch v := raw.(type) {
	case [kqMapLinkSlots]string:
		return v
	case []string:
		for i := 0; i < len(v) && i < kqMapLinkSlots; i++ {
			out[i] = strings.TrimSpace(v[i])
		}
	case []interface{}:
		for i := 0; i < len(v) && i < kqMapLinkSlots; i++ {
			if s, ok := v[i].(string); ok {
				out[i] = strings.TrimSpace(s)
			}
		}
	}
	return out
}

func (h *Handler) sendZ2WMakeAck(handle uint32, errCode uint16) {
	// exe: 0x00497AB6 (0x580E PROTO_NC_KQ_Z2W_MAKE_ACK: hquest[4], err[2])
	h.PublishOutboundWMS("wms_NC_KQ_Z2W_MAKE_ACK", map[string]interface{}{
		"opcode": protocol.PROTO_NC_KQ_Z2W_MAKE_ACK,
		"handle": handle,
		"err":    errCode,
	})
}

// wms_EndOfKQPacket ports WorldManagerSession::wms_EndOfKQPacket (exe: 0x005A4800).
// Emits PROTO_NC_KQ_Z2W_END_CMD (0x5810, 4-byte Handle) to WorldManager and clears contribute entry.
func (h *Handler) wms_EndOfKQPacket(handle uint32) {
	// exe: 0x005A486C (0x5810 PROTO_NC_KQ_Z2W_END_CMD: Handle[4])
	h.PublishOutboundWMS("wms_NC_KQ_Z2W_END_CMD", map[string]interface{}{
		"opcode": protocol.PROTO_NC_KQ_Z2W_END_CMD,
		"handle": handle,
	})
	// exe: 0x005A48B5: g_KQContributeList.kqpl_EndOfKQ(handle)
	h.kqState().contribute.kqpl_EndOfKQ(handle)
}

// ============================================================================
// S2S Handlers (wms_NC_KQ_W2Z_MAKE_REQ, W2Z_START_CMD, W2Z_DESTROY_CMD,
//               PLAYER_DISJOIN_CMD, LINK_TO_FORCE_BY_BAN_CMD)
// ============================================================================

// wms_NC_KQ_W2Z_MAKE_REQ ports WorldManagerSession::wms_NC_KQ_W2Z_MAKE_REQ (exe: 0x00497790).
func (h *Handler) wms_NC_KQ_W2Z_MAKE_REQ(ev WMSEvent) {
	info, err := decodeProtoKQInfoFromWMS(ev.Payload)
	if err != nil || info == nil {
		return
	}
	handle := info.Handle
	st := h.kqState()

	// 1. Duplicate handle check (exe: 0x00497ADE -> 0x00497B9E: 0x0982)
	if _, _, exists := st.container.kqc_SearchHandle(handle); exists {
		h.sendZ2WMakeAck(handle, kqErrMakeDuplicate)
		return
	}

	// 2. Scenario script lookup (exe: 0x00497C0F -> 0x00497C60: 0x098C)
	scriptLang := info.ScriptLanguageString()
	if _, ok := h.sbs_GetScenarioBook(scriptLang); !ok {
		h.sendZ2WMakeAck(handle, kqErrMakeScriptNotFound)
		return
	}

	// 3. Validate and bind up to 4 MapLink slots (exe: 0x00497D30..0x0049811D)
	var elemMaps [kqMapLinkSlots]KQElementMap
	localHosted := true // exe: 0x00497CA6 (mov byte [ebp-0x7db], 1)

	rollbackBoundMaps := func(upto int) {
		for j := 0; j < upto; j++ {
			if fm := elemMaps[j].FieldMap; fm != nil {
				fm.mu.Lock()
				if fm.KQHandle == handle {
					fm.KQHandle = kqInvalidHandle
					fm.KQID = kqInvalidID
				}
				fm.mu.Unlock()
			}
		}
	}

	for i := 0; i < kqMapLinkSlots; i++ {
		ml := info.MapLink[i]
		mapName := ml.MapNameString()
		mapBase := ml.MapBaseString()

		// Empty MapName check (exe: 0x00497D30..0x00497D51 & 0x004982E1)
		if mapName == "" {
			if i == 0 {
				// exe: 0x004982ED (errTableNoMap[0] = 0x0984)
				h.sendZ2WMakeAck(handle, kqErrMakeNoMap[0])
				return
			}
			// exe: 0x004982E7 (i > 0 -> break to allocation)
			break
		}

		// Submap already in use by another KQ? (exe: 0x00497D8C -> 0x00498138: 0x0988+i)
		if !st.container.kqc_MapUseCheck(mapName) {
			rollbackBoundMaps(i)
			h.sendZ2WMakeAck(handle, kqErrMakeMapInUse[i])
			return
		}

		// Base map exists in MapInfo.shn? (exe: 0x00497E23..0x00497E80 -> 0x0049853D: 0x0984+i)
		var mapInfo *KQMapInfoEntry
		if h.dataStore != nil {
			mapInfo, _ = h.dataStore.GetKQMapInfo(mapBase)
		}
		if mapInfo == nil {
			rollbackBoundMaps(i)
			h.sendZ2WMakeAck(handle, kqErrMakeNoMap[i])
			return
		}

		elemMaps[i].MapInfo = ml
		elemMaps[i].RegenX = mapInfo.RegenX
		elemMaps[i].RegenY = mapInfo.RegenY

		// Resolve or host FieldMap on this Zone node (exe: 0x00497D9C..0x004980DD)
		var fm *ZoneMap
		if h.mapManager != nil {
			fm = h.mapManager.EnsureKQMap(mapName, mapBase)
		}
		elemMaps[i].FieldMap = fm
		if fm == nil {
			// exe: 0x00497E86..0x00497E92: map not hosted on this Zone -> localHosted = 0
			localHosted = false
			continue
		}

		// Check if FieldMap is already bound to another KQ (exe: 0x004980A6 -> 0x00498715: 0x0988+i)
		fm.mu.Lock()
		if fm.KQHandle != kqInvalidHandle {
			fm.mu.Unlock()
			rollbackBoundMaps(i)
			h.sendZ2WMakeAck(handle, kqErrMakeMapInUse[i])
			return
		}
		localHosted = true // exe: 0x004980AA
		fm.KQHandle = handle
		fm.KQID = info.ID
		fm.IsInstance = true
		fm.InstanceKind = "KQ"
		fm.InstanceDefID = info.ID
		fm.InstanceLimitSec = int(info.LimitTime) * kqSecondsPerMinute
		fm.mu.Unlock()

		// Clear map objects if MapClear != 0 (exe: 0x004980CA..0x004980D8: fm_ClearObject(0x3B6))
		if ml.MapClear != 0 {
			fm.fm_ClearObject(kqClearObjectMaskMake)
		}
	}

	// 4. Allocate KQElement slot in KingdomQuestContainer (exe: 0x004987C4..0x00498818: 0x0983)
	el := &KQElement{
		Info:            *info,
		Maps:            elemMaps,
		ScriptLanguage:  scriptLang,
		ScriptInitValue: info.ScriptInitValueString(),
	}
	el.PlayerInfoList.clear() // exe: 0x00498952

	st.container.mu.Lock()
	_, allocated := st.container.allocLocked(el)
	st.container.mu.Unlock()
	if !allocated {
		rollbackBoundMaps(kqMapLinkSlots)
		h.sendZ2WMakeAck(handle, kqErrMakeTooMany) // exe: 0x00498813
		return
	}

	// 5. If hosted locally, reply Z2W_MAKE_ACK(handle, 0x0981) (exe: 0x004989DC..0x00498B28)
	if localHosted {
		h.sendZ2WMakeAck(handle, kqErrMakeNoErr)
	}

	// 6. Initialize KQ contribute table (exe: 0x00498B3B)
	st.contribute.kqcl_NewKQ(handle)
}

// wms_NC_KQ_W2Z_START_CMD ports WorldManagerSession::wms_NC_KQ_W2Z_START_CMD (exe: 0x00499240).
func (h *Handler) wms_NC_KQ_W2Z_START_CMD(ev WMSEvent) {
	info, err := decodeProtoKQInfoFromWMS(ev.Payload)
	if err != nil || info == nil {
		return
	}
	joiners := decodeJoinersFromWMS(ev.Payload)
	st := h.kqState()

	// 1. Lookup KQElement by handle (exe: 0x004993B1)
	elem, _, ok := st.container.kqc_SearchHandle(info.Handle)
	if !ok || elem == nil {
		return // exe: 0x004993C3 -> Fail: KQ not exist
	}

	now := st.now()
	// 2. If Map[0] is hosted locally, start quest & populate KQPlayerInfoList (exe: 0x0049951F..0x00499573)
	if elem.Maps[0].FieldMap != nil {
		elem.kqe_QuestStart(st.curTick(), now) // exe: 0x00499529
		for _, j := range joiners {
			elem.PlayerInfoList.kqpil_InsertPlayerInfo(j.CharRegNum, j.TeamType, 1) // exe: 0x00499555
		}
	}

	// 3. Register joiner list into KingdomQuestEntrance ring buffer (exe: 0x0049957B)
	entryIdx := st.entrance.kqe_newentrance(joiners)

	// 4. Schedule 5 KingdomQuestGateKeeper waves at t + {0, 2, 5, 7, 10}s (exe: 0x00499600..0x004996BA)
	mapName := elem.Maps[0].MapInfo.MapNameString()
	mapBase := elem.Maps[0].MapInfo.MapBaseString()
	regenX := elem.Maps[0].RegenX
	regenY := elem.Maps[0].RegenY

	st.mu.Lock()
	for _, delaySec := range kqGateKeeperDelaySec {
		st.gateKeepers = append(st.gateKeepers, &KingdomQuestGateKeeper{
			EntryIdx:  entryIdx,
			MapName:   mapName,
			MapBase:   mapBase,
			RegenX:    regenX,
			RegenY:    regenY,
			KQID:      info.ID,
			KQHandle:  info.Handle,
			TriggerAt: now.Add(time.Duration(delaySec) * time.Second),
		})
	}
	st.mu.Unlock()

	// Process wave 0 (delay = 0s) immediately at start time
	h.ProcessKQGateKeepers(now)
}

// wms_NC_KQ_W2Z_DESTROY_CMD ports WorldManagerSession::wms_NC_KQ_W2Z_DESTROY_CMD (exe: 0x004972C0).
func (h *Handler) wms_NC_KQ_W2Z_DESTROY_CMD(ev WMSEvent) {
	handleInt, ok := payloadInt(ev.Payload, "handle")
	if !ok {
		return
	}
	handle := uint32(handleInt)
	st := h.kqState()

	// exe: 0x00497334: kqc_SearchHandle
	elem, slotIdx, found := st.container.kqc_SearchHandle(handle)
	if !found || elem == nil {
		return
	}

	// exe: 0x0049734C: g_KQContributeList.kqpl_EndOfKQ
	st.contribute.kqpl_EndOfKQ(handle)

	// exe: 0x00497351..0x004973B3: release all 4 KQElementMap entries
	for i := 0; i < kqMapLinkSlots; i++ {
		fm := elem.Maps[i].FieldMap
		if fm == nil {
			continue
		}
		if elem.Maps[i].MapInfo.MapClear != 0 { // exe: 0x00497366
			fm.fm_ClearObject(kqClearObjectMaskDestroy) // 0x3B0, exe: 0x0049736C
		}
		fm.mu.Lock()
		fm.KQHandle = kqInvalidHandle // exe: 0x00497383 (-1)
		fm.KQID = kqInvalidID         // exe: 0x00497386 (0xFFFF)
		fm.mu.Unlock()

		fm.cc_DropFilm(elem.ScriptLanguage) // exe: 0x00497396
		fm.mda_CloseAllDoor()               // exe: 0x004973A5
	}

	// exe: 0x004973BE: kqcontainer.l_Free(slotIdx)
	st.container.freeSlot(slotIdx)
}

// wms_NC_KQ_PLAYER_DISJOIN_CMD ports WorldManagerSession::wms_NC_KQ_PLAYER_DISJOIN_CMD (exe: 0x00499100).
func (h *Handler) wms_NC_KQ_PLAYER_DISJOIN_CMD(ev WMSEvent) {
	handleInt, okH := payloadInt(ev.Payload, "handle")
	charInt, okC := payloadInt(ev.Payload, "char_reg_num")
	if !okH || !okC {
		return
	}
	elem, _, found := h.kqState().container.kqc_SearchHandle(uint32(handleInt))
	if !found || elem == nil {
		return
	}
	// exe: 0x00499127: kqpil_DeletePlayerInfo
	elem.PlayerInfoList.kqpil_DeletePlayerInfo(uint32(charInt))
}

// wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD ports WorldManagerSession::wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD (exe: 0x00497120).
func (h *Handler) wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD(ev WMSEvent) {
	charInt, ok := payloadInt(ev.Payload, "char_reg_num")
	if !ok {
		return
	}
	// exe: 0x00497193: som_FindPlayer
	player, found := h.GetPlayerByCharID(charInt)
	if !found || player == nil || player.State == 0 || player.CurrentMap == nil {
		return // exe: 0x0049719C / 0x004971B4 / 0x004971BB
	}

	outMaps := decodeOutMapsFromWMS(ev.Payload)
	curMapName := player.CurrentMap.MapID
	// exe: 0x004971EE..0x00497217: check if player is currently on any of the 4 KQ maps
	inBannedKQMap := false
	for i := 0; i < kqMapLinkSlots; i++ {
		if outMaps[i] != "" && strings.EqualFold(curMapName, outMaps[i]) {
			inBannedKQMap = true
			break
		}
	}
	if !inBannedKQMap {
		return
	}
	// exe: 0x004971C2..0x0049722D: force link to RouN (0x192D, 0x21B6) = (6445, 8630)
	h.WarpPlayer(player, kqRoumenFallbackMap, kqRoumenFallbackX, kqRoumenFallbackY, 0)
}

// ============================================================================
// ShinePlayer KQ Methods (sp_NC_KQ_ENTRYRESPONCE_ACK, sp_MyKQElement,
//   sp_KQElementFunction, sp_EscapeFromEmptyKQ, sp_2Roumen, sp_KQItemDrop,
//   sp_RemoveKQItem, sp_GetJoinKQTeamType, sp_SendKQTeamType, so_ply_IsKQJoiner,
//   so_IsTeamPVPKQ, so_smo_TargetIsEnemyKQTeam)
// ============================================================================

// handleKQEntryResponseAck ports ShinePlayer::sp_NC_KQ_ENTRYRESPONCE_ACK (exe: 0x00496990).
// Decodes kqHandle = uint32(player.MapObjectID) ^ EncHandle and warps into KQElement.Maps[0].
func (h *Handler) handleKQEntryResponseAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt == nil || pkt.Reader.Remaining() < 5 {
		return
	}
	// exe: 0x00496A6E..0x00496A90: state bitmask check 0x56 (states 1, 2, 4, 6)
	if (kqEntryStateValidMask & (uint16(1) << player.State)) == 0 {
		return
	}
	reply, _ := pkt.Reader.ReadByte()
	encHandle, _ := pkt.Reader.ReadUint32()
	// exe: 0x00496A95: cmp byte [edi+2], 0; je exit
	if reply == 0 {
		return
	}
	// exe: 0x00496AA6..0x00496AB4: if state == 4 (Resting), stand up
	if player.State == 4 {
		player.State = 1
	}

	// exe: 0x00496AB6..0x00496ACA: kqHandle = uint32(so_GetHandle()) ^ EncHandle
	kqHandle := uint32(player.MapObjectID) ^ encHandle

	// exe: 0x00496AE1: kqc_SearchHandle(kqHandle)
	elem, _, found := h.kqState().container.kqc_SearchHandle(kqHandle)
	if !found || elem == nil {
		return // exe: 0x00496AE8 -> exit without fabricating a KQ
	}

	mapName := elem.Maps[0].MapInfo.MapNameString()
	if mapName == "" {
		return
	}
	// exe: 0x00496B26..0x00496B52: MapLink to Maps[0] at MapInfo RegenX/RegenY
	prevMap := ""
	prevX, prevY := player.Char.XPos, player.Char.YPos
	if player.CurrentMap != nil {
		prevMap = player.CurrentMap.MapID
	}
	h.SavePreKQReturnPos(player, prevMap, prevX, prevY)
	h.WarpPlayer(player, mapName, int(elem.Maps[0].RegenX), int(elem.Maps[0].RegenY), 0)
	h.so_SaveLocation(player)
	h.EmitKQEnterGameLog(player, prevMap, prevX, prevY, elem.Info.ID, elem.Info.Handle)
}

// sp_MyKQElement ports ShinePlayer::sp_MyKQElement (exe: 0x00496C50).
func (h *Handler) sp_MyKQElement(player *ZonePlayer) *KQElement {
	if player == nil || player.CurrentMap == nil {
		return nil // exe: 0x00496CC0
	}
	kqHandle := player.CurrentMap.fm_GetKQhandle()
	if kqHandle == kqInvalidHandle { // exe: 0x00496CCB
		return nil
	}
	elem, _, found := h.kqState().container.kqc_SearchHandle(kqHandle) // exe: 0x00496CD9
	if !found {
		return nil
	}
	return elem
}

// sp_KQElementFunction ports ShinePlayer::sp_KQElementFunction (exe: 0x00499710).
func (h *Handler) sp_KQElementFunction(player *ZonePlayer, fn func(*KQPlayerInform, uint32) int) int {
	if player == nil || player.Char == nil || fn == nil {
		return 0
	}
	elem := h.sp_MyKQElement(player) // exe: 0x00499776
	if elem == nil {
		return 0 // exe: 0x0049977F
	}
	return elem.PlayerInfoList.kqpil_PlayerInformFunc(fn, uint32(player.Char.ID)) // exe: 0x004997C9
}

// sp_2Roumen ports ShinePlayer::sp_2Roumen (exe: 0x0043CF80).
// Warps player out to their current map's configured ReviveMap/Respawn city,
// falling back to RouN (0x192D, 0x21B6) = (6445, 8630) if FieldAttr is missing.
func (h *Handler) sp_2Roumen(player *ZonePlayer) {
	if player == nil || player.CurrentMap == nil {
		return // exe: 0x0043CFF5
	}
	destMap := kqRoumenFallbackMap
	destX := kqRoumenFallbackX
	destY := kqRoumenFallbackY

	if h.fieldStore != nil {
		curName := player.CurrentMap.MapID
		fa, ok := h.fieldStore.GetField(curName)
		if !ok && player.CurrentMap.FieldEntry != nil {
			fa = player.CurrentMap.FieldEntry
			ok = true
		}
		if ok && fa != nil {
			// exe: 0x0043D081..0x0043D0F2: resolve ReviveMap FieldAttr if present
			revCity := strings.TrimSpace(fa.Respawn.City)
			targetFA := fa
			if revCity != "" && revCity != "-" {
				if cityFA, okCity := h.fieldStore.GetField(revCity); okCity && cityFA != nil {
					targetFA = cityFA
				}
			}
			if targetFA.Respawn.City != "" && targetFA.Respawn.City != "-" {
				destMap = targetFA.Respawn.City
			}
			if len(targetFA.Respawn.Spots) > 0 && (targetFA.Respawn.Spots[0].X != 0 || targetFA.Respawn.Spots[0].Y != 0) {
				destX = targetFA.Respawn.Spots[0].X
				destY = targetFA.Respawn.Spots[0].Y
			}
		}
	}
	h.WarpPlayer(player, destMap, destX, destY, 0)
}

func (h *Handler) isPlayerOnKQMap(player *ZonePlayer) bool {
	if player == nil || player.CurrentMap == nil {
		return false
	}
	if player.CurrentMap.FieldEntry != nil && player.CurrentMap.FieldEntry.IsKQ {
		return true
	}
	if h.dataStore != nil {
		if mi, ok := h.dataStore.GetKQMapInfo(player.CurrentMap.MapID); ok && mi != nil && mi.KingdomMap == 1 {
			return true
		}
	}
	return player.CurrentMap.InstanceKind == "KQ"
}

// sp_EscapeFromEmptyKQ ports ShinePlayer::sp_EscapeFromEmptyKQ (exe: 0x0043FA40).
// Invoked during sp_NC_MAP_LOGINCOMPLETE_CMD (exe: 0x00444034).
func (h *Handler) sp_EscapeFromEmptyKQ(player *ZonePlayer) {
	if player == nil || player.Char == nil || player.CurrentMap == nil {
		return // exe: 0x0043FAAE
	}
	// exe: 0x0043FAC1: only applies on KQ maps (KingdomMap / IsKQ == 1)
	if !h.isPlayerOnKQMap(player) {
		return
	}

	kqHandle := player.CurrentMap.fm_GetKQhandle()
	if kqHandle == kqInvalidHandle { // exe: 0x0043FACA -> 0x0043FC2A
		player.CurrentMap.mu.RLock()
		hasFilm := player.CurrentMap.KQScriptActive
		player.CurrentMap.mu.RUnlock()
		// exe: 0x0043FC35..0x0043FC4C: if CinemaComplex has no running film -> sp_2Roumen
		if !hasFilm {
			h.sp_2Roumen(player)
		}
		return
	}

	// exe: 0x0043FAE0: kqc_SearchHandle(kqHandle)
	elem, _, found := h.kqState().container.kqc_SearchHandle(kqHandle)
	if !found || elem == nil {
		return // exe: 0x0043FAF6
	}

	regNum := uint32(player.Char.ID)
	// exe: 0x0043FB2B: kqpil_IsKQJoiner
	if elem.PlayerInfoList.kqpil_IsKQJoiner(regNum) == 1 {
		// exe: 0x0043FB4B: get knockdown count
		kd := elem.PlayerInfoList.kqpil_StorePlayer_GetKnockDownNumber(regNum)
		// exe: 0x0043FB6B: get team type & apply to player
		team := elem.PlayerInfoList.kqpil_PlayerInformFunc((*KQPlayerInform).kqpi_GetTeamType, regNum)
		player.KQTeam = byte(team) // exe: 0x0043FB7B (so_SetKQTeamType)

		// exe: 0x0043FB83: kqe_ReviveNumber
		maxRevives, _ := elem.kqe_ReviveNumber()
		if kd >= maxRevives { // exe: 0x0043FB8B
			h.sp_2Roumen(player) // exe: 0x0043FB91
			return
		}

		// exe: 0x0043FB98..0x0043FBF5: emit PROTO_NC_KQ_RESTDEADNUM_CMD (0x5818, 1 byte)
		var rest byte
		if maxRevives == kqReviveUnlimited { // exe: 0x0043FBA6
			rest = 0 // exe: 0x0043FBAD
		} else {
			rest = byte(maxRevives - kd) // exe: 0x0043FBB1
		}
		if player.Client != nil {
			pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_RESTDEADNUM_CMD)
			pw.WriteByte(rest)
			_ = player.Client.SendPacket(pw)
		}
		return
	}

	// Not a joiner (exe: 0x0043FBF9..0x0043FC14):
	// Admins stay with "Doing KQ" notice; normal players escape via sp_2Roumen.
	if player.Admin > 0 || player.Char.AdminLevel > 0 {
		h.sendPlayerNotice(player, "Doing KQ") // exe: 0x0043FC0F
		return
	}
	h.sp_2Roumen(player) // exe: 0x0043FC05 -> 0x0043FB91
}

// sp_GetJoinKQTeamType ports ShinePlayer::sp_GetJoinKQTeamType (exe: 0x00563990).
// Called in sp_NC_MAP_LOGINCOMPLETE_CMD (exe: 0x00443922).
func (h *Handler) sp_GetJoinKQTeamType(player *ZonePlayer) {
	if player == nil || player.Char == nil || player.CurrentMap == nil {
		return
	}
	if !h.isPlayerOnKQMap(player) { // exe: 0x00563A15
		player.KQTeam = kqTeamMax
		return
	}
	kqHandle := player.CurrentMap.fm_GetKQhandle()
	if kqHandle == kqInvalidHandle || kqHandle == uint32(kqInvalidID) { // exe: 0x00563A1E
		player.KQTeam = kqTeamMax
		return
	}
	elem, _, found := h.kqState().container.kqc_SearchHandle(kqHandle) // exe: 0x00563A30
	if !found || elem == nil {
		return
	}
	regNum := uint32(player.Char.ID)
	if elem.PlayerInfoList.kqpil_IsKQJoiner(regNum) == 0 { // exe: 0x00563A60
		return
	}
	team := elem.PlayerInfoList.kqpil_PlayerInformFunc((*KQPlayerInform).kqpi_GetTeamType, regNum) // exe: 0x00563A7D
	player.KQTeam = byte(team)                                                                     // exe: 0x00563A8D
}

// sp_SendKQTeamType ports ShinePlayer::sp_SendKQTeamType (exe: 0x00561560).
// Emits PROTO_NC_KQ_TEAM_TYPE_CMD (0x583A, 1 byte: nTeamType).
func (h *Handler) sp_SendKQTeamType(player *ZonePlayer) {
	if player == nil || player.Client == nil {
		return
	}
	// exe: 0x00561568 (0x583A, len=3 = 2B opcode + 1B teamType)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_TEAM_TYPE_CMD)
	pw.WriteByte(player.KQTeam)
	_ = player.Client.SendPacket(pw)
}

// so_ply_IsKQJoiner ports ShinePlayer::so_ply_IsKQJoiner (exe: 0x005615D0).
func (h *Handler) so_ply_IsKQJoiner(player *ZonePlayer) bool {
	if player == nil || player.Char == nil {
		return false
	}
	elem := h.sp_MyKQElement(player) // exe: 0x005615D4
	if elem == nil {
		return false
	}
	return elem.PlayerInfoList.kqpil_IsKQJoiner(uint32(player.Char.ID)) != 0
}

// so_IsTeamPVPKQ ports ShineObject::so_IsTeamPVPKQ (exe: 0x0055D6C0).
func (h *Handler) so_IsTeamPVPKQ(player *ZonePlayer) bool {
	elem := h.sp_MyKQElement(player) // exe: 0x0055D6D5..0x0055D6E4
	if elem == nil {
		return false
	}
	return elem.Info.IsTeamPVP != 0 // exe: 0x0055D6ED (byte [eax+0x16C] != 0)
}

// so_smo_TargetIsEnemyKQTeam ports ShineMobileObject::so_smo_TargetIsEnemyKQTeam (exe: 0x0055D700).
func (h *Handler) so_smo_TargetIsEnemyKQTeam(attacker, target *ZonePlayer) bool {
	if attacker == nil || target == nil {
		return false // exe: 0x0055D70A
	}
	if !h.so_IsTeamPVPKQ(attacker) { // exe: 0x0055D716..0x0055D740
		return false
	}
	if attacker.KQTeam == kqTeamMax || target.KQTeam == kqTeamMax { // exe: 0x0055D75D..0x0055D764
		return false
	}
	return attacker.KQTeam != target.KQTeam // exe: 0x0055D76F
}

// sp_KQItemDrop ports ShinePlayer::sp_KQItemDrop (exe: 0x00496F00).
// Iterates inventory items whose ItemInfoServer.KQ_Item_Drop (+0xDC) != 0 and drops them on the map.
// Called from ShinePlayer::so_Unmark (exe: 0x005589D1).
func (h *Handler) sp_KQItemDrop(player *ZonePlayer) int {
	if player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil || h.itemStore == nil {
		return 0
	}
	items, err := h.charRepo.GetItemsByStorage(player.Char.ID, kqRemoveBoxInven)
	if err != nil || len(items) == 0 {
		return 0
	}
	dropped := 0
	for _, it := range items {
		if it == nil {
			continue
		}
		info, ok := h.itemStore.GetItemInfo(it.ItemID)
		if !ok || info == nil || !info.KQItemDrop { // exe: 0x0049701F (byte [eax+0xDC] != 0)
			continue
		}
		lot := uint32(it.Flags)
		if lot == 0 {
			lot = 1
		}
		// Remove from player inventory and spawn on ground (so_IsDropping @ 0x004970AE)
		_ = h.charRepo.DeleteItem(it.ItemKey)
		if player.Client != nil {
			h.sendCellChangeFull(player.Client, kqRemoveBoxInven, it.Storage, nil, kqInvalidID)
		}
		player.CurrentMap.SpawnDroppedItem(it.ItemID, lot, int(it.Flags), int32(player.Char.XPos), int32(player.Char.YPos), player.MapObjectID, 0)
		dropped++
	}
	return dropped
}

// sp_RemoveKQItem ports ShinePlayer::sp_RemoveKQItem (exe: 0x0050F8C0).
// boxType: 8 = Equipped, 9 = Inventory (exe: 0x0050F978..0x0050F98C).
// Removes items with ItemInfoServer.KQItem (+0xDA) != 0; if player is on a KQ map,
// only removes those that ALSO have KQ_Item_Drop (+0xDC) != 0 (exe: 0x0050FA69..0x0050FA79).
func (h *Handler) sp_RemoveKQItem(player *ZonePlayer, boxType byte) int {
	if player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil || h.itemStore == nil {
		return 0
	}
	if boxType != kqRemoveBoxEquip && boxType != kqRemoveBoxInven { // exe: 0x0050F978..0x0050F97E
		return 0
	}
	items, err := h.charRepo.GetItemsByStorage(player.Char.ID, boxType)
	if err != nil || len(items) == 0 {
		return 0
	}
	onKQMap := h.isPlayerOnKQMap(player)
	removed := 0
	for _, it := range items {
		if it == nil {
			continue
		}
		info, ok := h.itemStore.GetItemInfo(it.ItemID)
		if !ok || info == nil || !info.KQItem { // exe: 0x0050FA5C (byte [eax+0xDA] != 0)
			continue
		}
		if onKQMap && !info.KQItemDrop { // exe: 0x0050FA6C..0x0050FA79
			continue
		}
		_ = h.charRepo.DeleteItem(it.ItemKey)
		if player.Client != nil {
			h.sendCellChangeFull(player.Client, boxType, it.Storage, nil, kqInvalidID)
		}
		removed++
	}
	return removed
}

// sp_GetItemInvenLot counts the total quantity of itemID currently in the player's inventory.
// Mirrors ShinePlayer::sp_GetItemInvenLot (exe: 0x005280D0, called at 0x0052FC5E).
func (h *Handler) sp_GetItemInvenLot(player *ZonePlayer, itemID uint16) uint64 {
	if player == nil || player.Char == nil || h.charRepo == nil {
		return 0
	}
	items, err := h.charRepo.GetItemsByStorage(player.Char.ID, kqRemoveBoxInven)
	if err != nil {
		return 0
	}
	var total uint64
	for _, it := range items {
		if it != nil && it.ItemID == itemID {
			if it.Flags > 0 {
				total += uint64(it.Flags)
			} else {
				total++
			}
		}
	}
	return total
}

// onPlayerKQKnockDown handles knockdown tracking on player death in a KQ map.
// Ports the KQ hook in ShinePlayer::so_KilledBy (exe: 0x00434630..0x00434637).
func (h *Handler) onPlayerKQKnockDown(player *ZonePlayer) {
	if h == nil || player == nil || player.Char == nil {
		return
	}
	_ = h.sp_KQElementFunction(player, (*KQPlayerInform).kqpi_KnockDownPlus)
}

// onPlayerMapLoginCompleteKQ executes the 4 official KQ hooks in sp_NC_MAP_LOGINCOMPLETE_CMD:
// 0x00443922 sp_GetJoinKQTeamType, 0x00443929 sp_SendKQTeamType,
// 0x00443E4F/0x00443E58 sp_RemoveKQItem(8/9), 0x00444034 sp_EscapeFromEmptyKQ.
func (h *Handler) onPlayerMapLoginCompleteKQ(player *ZonePlayer) {
	if h == nil || player == nil || player.Char == nil {
		return
	}
	h.sp_GetJoinKQTeamType(player)
	h.sp_SendKQTeamType(player)
	h.sp_RemoveKQItem(player, kqRemoveBoxEquip)
	h.sp_RemoveKQItem(player, kqRemoveBoxInven)
	h.sp_EscapeFromEmptyKQ(player)
}

// ============================================================================
// Official GM Commands (AmpersandCommand registrations @ 0x00426526..0x00426642)
// ============================================================================

// ac_KingdomQuest ports AmpersandCommand::ac_KingdomQuest (exe: 0x00423DD0, "&kingdomquest").
// Usage: &kingdomquest <mapName> <scriptName> [initValue]
func (h *Handler) ac_KingdomQuest(player *ZonePlayer, args []string) error {
	if player == nil || len(args) < 2 {
		return errors.New("Usage: &kingdomquest <mapName> <scriptName> [initValue]")
	}
	mapName := strings.TrimSpace(args[0])
	scriptName := strings.TrimSpace(args[1])
	initVal := ""
	if len(args) >= 3 {
		initVal = strings.TrimSpace(args[2])
	}

	// exe: 0x00423E67: fo_GetFieldAttr(mapName)
	var baseMap string
	if h.fieldStore != nil {
		if fa, ok := h.fieldStore.GetField(mapName); ok && fa != nil {
			baseMap = fa.MapID
		}
	}
	if baseMap == "" {
		// Also accept submap name (e.g. KDRockCan00 -> KDRockCan) via MapInfo.shn
		if h.dataStore != nil {
			if mi, ok := h.dataStore.GetKQMapInfo(mapName); ok && mi != nil {
				baseMap = mi.MapName
			}
		}
	}
	if baseMap == "" {
		h.sendPlayerNotice(player, "Invalid map name") // exe: 0x00423E78
		return nil
	}

	// exe: 0x00423EC0..0x00423F09: MapInfo lookup
	var mapInfo *KQMapInfoEntry
	if h.dataStore != nil {
		mapInfo, _ = h.dataStore.GetKQMapInfo(baseMap)
	}
	if mapInfo == nil {
		h.sendPlayerNotice(player, "Invalid map name") // exe: 0x00423EFE
		return nil
	}

	// exe: 0x00423F52..0x00423F73: warp GM to (mapName, mapInfo.RegenX, mapInfo.RegenY)
	var fm *ZoneMap
	if h.mapManager != nil {
		fm = h.mapManager.EnsureKQMap(mapName, baseMap)
	}
	h.WarpPlayer(player, mapName, int(mapInfo.RegenX), int(mapInfo.RegenY), 0)

	// exe: 0x00423FE4..0x0042406B: find FieldMap and call cc_PlayFilm(scriptName, initVal, fm)
	if fm != nil {
		fm.cc_PlayFilm(scriptName, initVal, h.kqState().now())
	}
	return nil
}

// ac_ReroadKingdomQuest ports AmpersandCommand::ac_ReroadKingdomQuest (exe: 0x0041DC80, "&Reroadkingdomquest").
// Usage: &Reroadkingdomquest <scriptName>
func (h *Handler) ac_ReroadKingdomQuest(player *ZonePlayer, args []string) error {
	if len(args) < 1 {
		return errors.New("Usage: &Reroadkingdomquest <scriptName>")
	}
	scriptName := strings.TrimSpace(args[0])
	// exe: 0x0041DD0A: sbs_GetScenarioBook(scriptName)
	if _, ok := h.sbs_GetScenarioBook(scriptName); !ok {
		return nil // exe: 0x0041DD11 (not found -> silent return)
	}
	// exe: 0x0041DD19: cm_AllDropFilm(scriptName) across all maps
	if h.mapManager != nil {
		for _, m := range h.mapManager.GetAllMaps() {
			m.cc_DropFilm(scriptName)
		}
	}
	// exe: 0x0041DD24: sbsReLoad(scriptName)
	st := h.kqState()
	st.mu.Lock()
	st.reloaded[scriptName]++
	st.mu.Unlock()
	return nil
}

// ac_EndOfKingdomQuest ports AmpersandCommand::ac_EndOfKingdomQuest (exe: 0x0041DBD0, "&endofkq").
// Usage: &endofkq <kqHandle>
func (h *Handler) ac_EndOfKingdomQuest(player *ZonePlayer, args []string) error {
	if len(args) < 1 {
		return errors.New("Usage: &endofkq <kqHandle>")
	}
	handleVal, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 32)
	if err != nil {
		return fmt.Errorf("invalid kqHandle: %w", err)
	}
	// exe: 0x0041DC4B: wms_EndOfKQPacket(handle)
	h.wms_EndOfKQPacket(uint32(handleVal))
	return nil
}

// ac_KQReward ports AmpersandCommand::ac_KQReward (exe: 0x00423CC0, "&kqreward").
// Usage: &kqreward <rewardIndex>
func (h *Handler) ac_KQReward(player *ZonePlayer, args []string) error {
	if player == nil || len(args) < 1 {
		return errors.New("Usage: &kqreward <rewardIndex>")
	}
	rewIdx, err := strconv.ParseUint(strings.TrimSpace(args[0]), 10, 16)
	if err != nil {
		return fmt.Errorf("invalid rewardIndex: %w", err)
	}
	// exe: 0x00423D5C..0x00423D6C: construct temp KQElement with RewardIndex and invoke so_ply_KQRewardStruct
	tmp := &KQElement{}
	tmp.Info.RewardIndex = uint16(rewIdx)
	h.so_ply_KQRewardStruct(player, tmp)
	return nil
}

// so_ply_KQRewardStruct está implementado 1:1 en kq_rewards_lua.go (exe: 0x0052E9C0).

func init() {
	// exe: 0x004265EE (&kingdomquest -> ac_KingdomQuest @ 0x00423DD0, level 0x64 @ 0x004265E9)
	registerGMCommand([]string{"kingdomquest"}, gmCommandSpec{
		MinArgs:           2,
		Usage:             "&kingdomquest <mapName> <scriptName> [initValue]",
		DefaultAdminLevel: kqGMAdminLevel,
		Exec: func(ctx *gmContext) error {
			return ctx.Handler.ac_KingdomQuest(ctx.Player, ctx.Args)
		},
	})

	// exe: 0x00426613 (&Reroadkingdomquest -> ac_ReroadKingdomQuest @ 0x0041DC80, level 0x64 @ 0x0042660E)
	registerGMCommand([]string{"reroadkingdomquest"}, gmCommandSpec{
		MinArgs:           1,
		Usage:             "&Reroadkingdomquest <scriptName>",
		DefaultAdminLevel: kqGMAdminLevel,
		Exec: func(ctx *gmContext) error {
			return ctx.Handler.ac_ReroadKingdomQuest(ctx.Player, ctx.Args)
		},
	})

	// exe: 0x00426638 (&endofkq -> ac_EndOfKingdomQuest @ 0x0041DBD0, level 0x64 @ 0x00426633)
	registerGMCommand([]string{"endofkq"}, gmCommandSpec{
		MinArgs:           1,
		Usage:             "&endofkq <kqHandle>",
		DefaultAdminLevel: kqGMAdminLevel,
		Exec: func(ctx *gmContext) error {
			return ctx.Handler.ac_EndOfKingdomQuest(ctx.Player, ctx.Args)
		},
	})

	// exe: 0x00426535 (&kqreward -> ac_KQReward @ 0x00423CC0, level 0x64 @ 0x00426533)
	registerGMCommand([]string{"kqreward"}, gmCommandSpec{
		MinArgs:           1,
		Usage:             "&kqreward <rewardIndex>",
		DefaultAdminLevel: kqGMAdminLevel,
		Exec: func(ctx *gmContext) error {
			return ctx.Handler.ac_KQReward(ctx.Player, ctx.Args)
		},
	})
}
