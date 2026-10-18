package misc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RecallCoordEntry defines the warp target map and coordinates for a recall scroll.
type RecallCoordEntry struct {
	ItemIndex string `json:"itemIndex"`
	ItemIdent int    `json:"itemIdent"`
	MapName   string `json:"mapName"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
}

// StealthMatrixConfig defines layer visibility and attack permissions.
type StealthMatrixConfig struct {
	Visibility map[string]map[string]bool `json:"visibility"`
	CanAttack  map[string]map[string]bool `json:"canAttack"`
}

// SkillBreedMob represents a summoned totem, tornado, trap, or entity spawned by a skill.
type SkillBreedMob struct {
	Serial      int    `json:"serial"`
	SkillIndex  string `json:"skillIndex"`
	MobIndex    string `json:"mobIndex"`
	Num         int    `json:"num"`
	MoveType    string `json:"moveType"`
	Target      string `json:"target"`
	Range       int    `json:"range"`
	Shape       string `json:"shape"`
	LifeTime    int    `json:"lifeTime"`
	AbState     string `json:"abState"`
	Filter      string `json:"filter"`
	Aggro       int    `json:"aggro"`
	Interval    int    `json:"interval"`
	ActiveTime  int    `json:"activeTime"`
	SubSkill    string `json:"subSkill"`
	ApplyTarget string `json:"applyTarget"`
	Disappear   string `json:"disappear"`
	Option      string `json:"option"`
}

// PolymorphEntry defines transformation skills.
type PolymorphEntry struct {
	Serial    int    `json:"serial"`
	Skill     string `json:"skill"`
	MobIndex  string `json:"mobIndex"`
	CanMove   bool   `json:"canMove"`
	CanAttack bool   `json:"canAttack"`
}

// NeutralizationEntry defines equipment disarm effects.
type NeutralizationEntry struct {
	Serial int    `json:"serial"`
	Skill  string `json:"skill"`
	Weapon int    `json:"weapon"`
	Shield int    `json:"shield"`
	Body   int    `json:"body"`
	Leg    int    `json:"leg"`
}

// AbStateResetEntry defines actions that break or cancel active buffs/stealth.
type AbStateResetEntry struct {
	Serial     int    `json:"serial"`
	AbState    string `json:"abState"`
	Run        bool   `json:"run"`
	Walk       bool   `json:"walk"`
	Attack     bool   `json:"attack"`
	Attacked   bool   `json:"attacked"`
	EquipWhere string `json:"equipWhere"`
	EquipClass string `json:"equipClass"`
}

// ExpandSkillEntry defines conditional skill combo damage and bonus effects.
type ExpandSkillEntry struct {
	Serial    int    `json:"serial"`
	Skill     string `json:"skill"`
	Condition string `json:"condition"`
	Dmg       int    `json:"dmg"`
	AbState   string `json:"abState"`
	Critical  int    `json:"critical"`
}

// SkillMechanicsConfig holds the 5 skill and combat mechanics tables.
type SkillMechanicsConfig struct {
	SkillBreedMobs []SkillBreedMob       `json:"skillBreedMobs"`
	Polymorph      []PolymorphEntry      `json:"polymorph"`
	Neutralization []NeutralizationEntry `json:"neutralization"`
	AbStateReset   []AbStateResetEntry   `json:"abStateReset"`
	ExpandSkill    []ExpandSkillEntry    `json:"expandSkill"`
}

// MobChatEntry defines dialogue lines and trigger rates for a monster.
type MobChatEntry struct {
	Mob      string   `json:"mob"`
	Rate0    int      `json:"rate0"`
	Rate1    int      `json:"rate1"`
	Messages []string `json:"messages"`
}

// MobChatConfig holds dialogue lines across all 7 trigger categories.
type MobChatConfig struct {
	Peace      []MobChatEntry `json:"peace"`
	Attack     []MobChatEntry `json:"attack"`
	Damaged    []MobChatEntry `json:"damaged"`
	Dead       []MobChatEntry `json:"dead"`
	HelpMain   []MobChatEntry `json:"helpMain"`
	HelpSub    []MobChatEntry `json:"helpSub"`
	RaceRescue []MobChatEntry `json:"raceRescue"`
}

// Store holds the parsed misc world data for ZoneService.
type Store struct {
	mu           sync.RWMutex
	recallCoords map[string]RecallCoordEntry
	stealth      StealthMatrixConfig
	mechanics    SkillMechanicsConfig
	mobChats     MobChatConfig
}

// NewStore initializes a new empty Misc store.
func NewStore() *Store {
	return &Store{
		recallCoords: make(map[string]RecallCoordEntry),
	}
}

// LoadDir loads all 4 misc JSON files from the designated directory.
func (s *Store) LoadDir(dir string) error {
	// 1. RecallCoord.json
	rcData, err := os.ReadFile(filepath.Join(dir, "RecallCoord.json"))
	if err != nil {
		return fmt.Errorf("read RecallCoord.json: %w", err)
	}
	var loadedRC []RecallCoordEntry
	if err := json.Unmarshal(rcData, &loadedRC); err != nil {
		return fmt.Errorf("unmarshal RecallCoord.json: %w", err)
	}

	// 2. StealthMatrix.json
	smData, err := os.ReadFile(filepath.Join(dir, "StealthMatrix.json"))
	if err != nil {
		return fmt.Errorf("read StealthMatrix.json: %w", err)
	}
	var loadedSM StealthMatrixConfig
	if err := json.Unmarshal(smData, &loadedSM); err != nil {
		return fmt.Errorf("unmarshal StealthMatrix.json: %w", err)
	}

	// 3. SkillMechanics.json
	smecData, err := os.ReadFile(filepath.Join(dir, "SkillMechanics.json"))
	if err != nil {
		return fmt.Errorf("read SkillMechanics.json: %w", err)
	}
	var loadedMechanics SkillMechanicsConfig
	if err := json.Unmarshal(smecData, &loadedMechanics); err != nil {
		return fmt.Errorf("unmarshal SkillMechanics.json: %w", err)
	}

	// 4. MobChat.json
	mcData, err := os.ReadFile(filepath.Join(dir, "MobChat.json"))
	if err != nil {
		return fmt.Errorf("read MobChat.json: %w", err)
	}
	var loadedMC MobChatConfig
	if err := json.Unmarshal(mcData, &loadedMC); err != nil {
		return fmt.Errorf("unmarshal MobChat.json: %w", err)
	}

	normRC := make(map[string]RecallCoordEntry, len(loadedRC))
	for _, rc := range loadedRC {
		normRC[strings.ToLower(rc.ItemIndex)] = rc
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.recallCoords = normRC
	s.stealth = loadedSM
	s.mechanics = loadedMechanics
	s.mobChats = loadedMC

	return nil
}

// GetRecallCoord returns warp destination coordinates for a recall scroll (case-insensitive).
func (s *Store) GetRecallCoord(itemIndex string) (*RecallCoordEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rc, ok := s.recallCoords[strings.ToLower(itemIndex)]
	if !ok {
		return nil, false
	}
	return &rc, true
}

// CanSee checks if a viewer on viewerLayer can detect an entity on targetLayer.
// In Fiesta's SendMyBrief table, targetLayer must send its brief to viewerLayer.
func (s *Store) CanSee(viewerLayer, targetLayer string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if row, ok := s.stealth.Visibility[targetLayer]; ok {
		return row[viewerLayer]
	}
	return false
}

// CanAttack checks if an entity on attackerLayer can attack an entity on targetLayer.
func (s *Store) CanAttack(attackerLayer, targetLayer string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if row, ok := s.stealth.CanAttack[attackerLayer]; ok {
		return row[targetLayer]
	}
	return false
}

// Mechanics returns skill mechanics, dummy mob summons, and buff reset rules.
func (s *Store) Mechanics() SkillMechanicsConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mechanics
}

// MobChats returns all monster dialogue configurations.
func (s *Store) MobChats() MobChatConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mobChats
}
