package field

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CombatFlags defines PvP and attack rules on a map.
type CombatFlags struct {
	ItemDropByKilled int  `json:"itemDropByKilled"`
	Pker             int  `json:"pker"`
	EnemyGuild       int  `json:"enemyGuild"`
	Party            int  `json:"party"`
	Monster          int  `json:"monster"`
	Summoned         int  `json:"summoned"`
	IsPkKq           bool `json:"isPkKq"`
	IsFreePk         bool `json:"isFreePk"`
	IsPartyBattle    bool `json:"isPartyBattle"`
	NameHide         bool `json:"nameHide"`
}

// TeleportFlags defines routing and system map status.
type TeleportFlags struct {
	LinkIn    bool `json:"linkIn"`
	LinkOut   bool `json:"linkOut"`
	SystemMap bool `json:"systemMap"`
}

// Coord represents a 2D position.
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// CoordDir represents a 2D position with orientation.
type CoordDir struct {
	X   int `json:"x"`
	Y   int `json:"y"`
	Dir int `json:"dir"`
}

// RespawnConfig defines return city and respawn coordinates upon death.
type RespawnConfig struct {
	City      string  `json:"city"`
	SpotCount int     `json:"spotCount"`
	Spots     []Coord `json:"spots"`
}

// PermissionFlags defines allowed actions and player features on a map.
type PermissionFlags struct {
	CanRestart   bool `json:"canRestart"`
	CanTrade     bool `json:"canTrade"`
	CanMiniHouse bool `json:"canMiniHouse"`
	CanItem      bool `json:"canItem"`
	CanSkill     bool `json:"canSkill"`
	Chat         bool `json:"chat"`
	Shout        bool `json:"shout"`
	CanBooth     bool `json:"canBooth"`
	CanProduce   bool `json:"canProduce"`
	CanRide      bool `json:"canRide"`
	CanStone     bool `json:"canStone"`
	CanParty     bool `json:"canParty"`
}

// ExpLossConfig defines experience penalty upon death.
type ExpLossConfig struct {
	ByMonster int `json:"byMonster"`
	ByPlayer  int `json:"byPlayer"`
}

// FieldEntry defines the complete configuration for a world map or sub-layer.
type FieldEntry struct {
	Serial      int             `json:"serial"`
	MapID       string          `json:"mapId"`
	SubFrom     string          `json:"subFrom"`
	SubTo       string          `json:"subTo"`
	Name        string          `json:"name"`
	IsKQ        bool            `json:"isKq"`
	KQ          int             `json:"kq"`
	Width       int             `json:"width"`
	Height      int             `json:"height"`
	ImmortalSec int             `json:"immortalSec"`
	ScriptName  string          `json:"scriptName"`
	Combat      CombatFlags     `json:"combat"`
	Teleport    TeleportFlags   `json:"teleport"`
	Respawn     RespawnConfig   `json:"respawn"`
	Permissions PermissionFlags `json:"permissions"`
	ExpLoss     ExpLossConfig   `json:"expLoss"`
	UsrSubLayer int             `json:"usrSubLayer"`
	CheckSum    int             `json:"checkSum"`
	Fiesta      int             `json:"fiesta"`
}

// InstanceDungeonEntry defines instance dungeon properties, zone process assignment, and requirements.
type InstanceDungeonEntry struct {
	Argument     string `json:"argument"`
	IDNo         int    `json:"idNo"`
	ZoneNumber   int    `json:"zoneNumber"`
	MapID        string `json:"mapId"`
	ScriptName   string `json:"scriptName"`
	ModeIDLv     int    `json:"modeIdLv"`
	EntranceType int    `json:"entranceType"`
	Guild        int    `json:"guild"`
	Individual   int    `json:"individual"`
	NeedQuest    int    `json:"needQuest"`
	NeedItem     string `json:"needItem"`
	Consume      int    `json:"consume"`
	LimitTime    int    `json:"limitTime"`
	MaxUseChr    int    `json:"maxUseChr"`
}

// TournamentDoor defines a guild arena gate spawn point.
type TournamentDoor struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Dir   int    `json:"dir"`
	Block string `json:"block"`
}

// GuildTournamentEntry defines a guild tournament arena.
type GuildTournamentEntry struct {
	MapID        string           `json:"mapId"`
	GTMapNo      int              `json:"gtMapNo"`
	Stone        CoordDir         `json:"stone"`
	GuildARegion Coord            `json:"guildARegion"`
	GuildBRegion Coord            `json:"guildBRegion"`
	DoorsA       []TournamentDoor `json:"doorsA"`
	DoorsB       []TournamentDoor `json:"doorsB"`
}

// TutorialMapEntry defines the new character tutorial map.
type TutorialMapEntry struct {
	MapID      string `json:"mapId"`
	MaxChr     int    `json:"maxChr"`
	ScriptName string `json:"scriptName"`
}

// Store holds the parsed world map tables for ZoneService.
type Store struct {
	mu                sync.RWMutex
	fields            []FieldEntry
	bySerial          map[int]*FieldEntry
	byMapID           map[string][]*FieldEntry
	instanceDungeons  []InstanceDungeonEntry
	guildTournaments  map[string]*GuildTournamentEntry
	tutorialMap       *TutorialMapEntry
}

// NewStore initializes a new empty Field store.
func NewStore() *Store {
	return &Store{
		bySerial:         make(map[int]*FieldEntry),
		byMapID:          make(map[string][]*FieldEntry),
		guildTournaments: make(map[string]*GuildTournamentEntry),
	}
}

// LoadDir loads all 4 map JSON files from the designated directory.
func (s *Store) LoadDir(dir string) error {
	// 1. FieldList.json
	fData, err := os.ReadFile(filepath.Join(dir, "FieldList.json"))
	if err != nil {
		return fmt.Errorf("read FieldList.json: %w", err)
	}
	var loadedFields []FieldEntry
	if err := json.Unmarshal(fData, &loadedFields); err != nil {
		return fmt.Errorf("unmarshal FieldList.json: %w", err)
	}

	// 2. InstanceDungeon.json
	idData, err := os.ReadFile(filepath.Join(dir, "InstanceDungeon.json"))
	if err != nil {
		return fmt.Errorf("read InstanceDungeon.json: %w", err)
	}
	var loadedID []InstanceDungeonEntry
	if err := json.Unmarshal(idData, &loadedID); err != nil {
		return fmt.Errorf("unmarshal InstanceDungeon.json: %w", err)
	}

	// 3. GuildTournament.json
	gtData, err := os.ReadFile(filepath.Join(dir, "GuildTournament.json"))
	if err != nil {
		return fmt.Errorf("read GuildTournament.json: %w", err)
	}
	var loadedGT []GuildTournamentEntry
	if err := json.Unmarshal(gtData, &loadedGT); err != nil {
		return fmt.Errorf("unmarshal GuildTournament.json: %w", err)
	}

	// 4. TutorialMap.json
	tutData, err := os.ReadFile(filepath.Join(dir, "TutorialMap.json"))
	if err != nil {
		return fmt.Errorf("read TutorialMap.json: %w", err)
	}
	var loadedTut []TutorialMapEntry
	if err := json.Unmarshal(tutData, &loadedTut); err != nil {
		return fmt.Errorf("unmarshal TutorialMap.json: %w", err)
	}

	bySerial := make(map[int]*FieldEntry, len(loadedFields))
	byMapID := make(map[string][]*FieldEntry)
	for i := range loadedFields {
		f := &loadedFields[i]
		bySerial[f.Serial] = f
		normID := strings.ToLower(f.MapID)
		byMapID[normID] = append(byMapID[normID], f)
	}

	gtMap := make(map[string]*GuildTournamentEntry, len(loadedGT))
	for i := range loadedGT {
		gt := &loadedGT[i]
		gtMap[strings.ToLower(gt.MapID)] = gt
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.fields = loadedFields
	s.bySerial = bySerial
	s.byMapID = byMapID
	s.instanceDungeons = loadedID
	s.guildTournaments = gtMap
	if len(loadedTut) > 0 {
		s.tutorialMap = &loadedTut[0]
	}

	return nil
}

// GetFieldBySerial returns a map by its unique serial integer ID.
func (s *Store) GetFieldBySerial(serial int) (*FieldEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.bySerial[serial]
	return f, ok
}

// GetField returns the default/primary field entry for a given map ID (case-insensitive).
func (s *Store) GetField(mapID string) (*FieldEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fields, ok := s.byMapID[strings.ToLower(mapID)]
	if !ok || len(fields) == 0 {
		return nil, false
	}
	return fields[0], true
}

// GetAllFieldsByMapID returns all sub-layers/variants for a map ID.
func (s *Store) GetAllFieldsByMapID(mapID string) []*FieldEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fields := s.byMapID[strings.ToLower(mapID)]
	res := make([]*FieldEntry, len(fields))
	copy(res, fields)
	return res
}

// GetDungeonsByZone returns all instance dungeons allocated to a specific Zone server number.
func (s *Store) GetDungeonsByZone(zoneNumber int) []InstanceDungeonEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []InstanceDungeonEntry
	for _, d := range s.instanceDungeons {
		if d.ZoneNumber == zoneNumber {
			res = append(res, d)
		}
	}
	return res
}

// FieldCount returns total registered fields in FieldList.
func (s *Store) FieldCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.fields)
}
