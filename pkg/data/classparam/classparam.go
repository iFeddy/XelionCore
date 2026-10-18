package classparam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LevelParam represents the complete official stat progression for a class at a specific level.
type LevelParam struct {
	Level          int `json:"level"`
	Strength       int `json:"strength"`
	Constitution   int `json:"constitution"`
	Intelligence   int `json:"intelligence"`
	Wizdom         int `json:"wizdom"`
	Dexterity      int `json:"dexterity"`
	MentalPower    int `json:"mentalPower"`
	SoulHP         int `json:"soulHp"`
	MAXSoulHP      int `json:"maxSoulHp"`
	PriceHPStone   int `json:"priceHpStone"`
	SoulSP         int `json:"soulSp"`
	MAXSoulSP      int `json:"maxSoulSp"`
	PriceSPStone   int `json:"priceSpStone"`
	AtkPerAP       int `json:"atkPerAp"`
	DmgPerAP       int `json:"dmgPerAp"`
	MaxPwrStone    int `json:"maxPwrStone"`
	NumPwrStone    int `json:"numPwrStone"`
	PricePwrStone  int `json:"pricePwrStone"`
	PwrStoneWC     int `json:"pwrStoneWc"`
	PwrStoneMA     int `json:"pwrStoneMa"`
	MaxGrdStone    int `json:"maxGrdStone"`
	NumGrdStone    int `json:"numGrdStone"`
	PriceGrdStone  int `json:"priceGrdStone"`
	GrdStoneAC     int `json:"grdStoneAc"`
	GrdStoneMR     int `json:"grdStoneMr"`
	PainRes        int `json:"painRes"`
	RestraintRes   int `json:"restraintRes"`
	CurseRes       int `json:"curseRes"`
	ShockRes       int `json:"shockRes"`
	MaxHP          int `json:"maxHp"`
	MaxSP          int `json:"maxSp"`
	CharTitlePt    int `json:"charTitlePt"`
	SkillPwrPt     int `json:"skillPwrPt"`
	JobChangeDmgUp int `json:"jobChangeDmgUp"`
}

// Store holds the parsed level statistics for all character classes.
type Store struct {
	mu      sync.RWMutex
	classes map[string][]LevelParam
}

// NewStore initializes a new empty ClassParam store.
func NewStore() *Store {
	return &Store{
		classes: make(map[string][]LevelParam),
	}
}

// LoadDir loads all class parameter JSON files from the specified directory.
func (s *Store) LoadDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read classparam dir: %w", err)
	}

	loaded := make(map[string][]LevelParam)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		className := strings.TrimSuffix(entry.Name(), ".json")
		filePath := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read %s: %w", filePath, err)
		}

		var levels []LevelParam
		if err := json.Unmarshal(data, &levels); err != nil {
			return fmt.Errorf("unmarshal %s: %w", filePath, err)
		}

		loaded[strings.ToLower(className)] = levels
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.classes = loaded
	return nil
}

// GetParam returns the LevelParam for a specific class and level (1-indexed).
func (s *Store) GetParam(className string, level int) (*LevelParam, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	levels, ok := s.classes[strings.ToLower(className)]
	if !ok || level < 1 || level > len(levels) {
		return nil, false
	}
	return &levels[level-1], true
}

// GetClass returns all levels for a given class.
func (s *Store) GetClass(className string) ([]LevelParam, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	levels, ok := s.classes[strings.ToLower(className)]
	if !ok {
		return nil, false
	}
	res := make([]LevelParam, len(levels))
	copy(res, levels)
	return res, true
}

// Classes returns the names of all registered classes.
func (s *Store) Classes() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.classes))
	for name := range s.classes {
		names = append(names, name)
	}
	return names
}

// ClassCount returns the number of loaded classes.
func (s *Store) ClassCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.classes)
}
