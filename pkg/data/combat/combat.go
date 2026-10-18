package combat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// CommonConfig defines global character constants and system settings.
type CommonConfig struct {
	RunSpeed         int `json:"runSpeed"`
	WalkSpeed        int `json:"walkSpeed"`
	AttackSpeed      int `json:"attackSpeed"`
	ShoutLevel       int `json:"shoutLevel"`
	ShoutDelay       int `json:"shoutDelay"`
	LevelLimit       int `json:"levelLimit"`
	MaxExpBonus      int `json:"maxExpBonus"`
	RestExpRate      int `json:"restExpRate"`
	MinNeedTime      int `json:"minNeedTime"`
	DefaultBonusTime int `json:"defaultBonusTime"`
	IntervalTime     int `json:"intervalTime"`
	AddBuffTime      int `json:"addBuffTime"`
	MaxBuffTime      int `json:"maxBuffTime"`
	LostExpLevel     int `json:"lostExpLevel"`
	MaxFreeStat      int `json:"maxFreeStat"`
}

// ExpLevelEntry defines level thresholds, exp required to advance, and stat distributions.
type ExpLevelEntry struct {
	Level           int            `json:"level"`
	NextExp         int            `json:"nextExp"`
	ExpLostAtPvP    int            `json:"expLostAtPvP"`
	ClassStatPoints map[string]int `json:"classStatPoints"`
}

// FreeStatBonus represents the stat bonuses earned from free stat point distribution.
type FreeStatTable struct {
	Str []struct {
		Points         int `json:"points"`
		PhysicalAttack int `json:"physicalAttack"`
	} `json:"str"`
	Int []struct {
		Points      int `json:"points"`
		MagicAttack int `json:"magicAttack"`
	} `json:"int"`
	Dex []struct {
		Points      int `json:"points"`
		AimRate     int `json:"aimRate"`
		EvasionRate int `json:"evasionRate"`
	} `json:"dex"`
	Con []struct {
		Points      int `json:"points"`
		PhysicalDef int `json:"physicalDef"`
		BlockRate   int `json:"blockRate"`
		MaxHP       int `json:"maxHp"`
	} `json:"con"`
	Men []struct {
		Points   int `json:"points"`
		MagicDef int `json:"magicDef"`
		CritRate int `json:"critRate"`
		MaxSP    int `json:"maxSp"`
	} `json:"men"`
}

// AngleRate defines a threshold angle and its associated damage multiplier (1000 = 100%).
type AngleRate struct {
	Angle      int `json:"angle"`
	DamageRate int `json:"damageRate"`
}

// DamageByAngle holds directional angle damage multipliers.
type DamageByAngle struct {
	Character []AngleRate `json:"character"`
	Monster   []AngleRate `json:"monster"`
}

// SoulDemandRate defines skill damage scaling when consuming accumulated soul stones.
type SoulDemandRate struct {
	DemandSoul int   `json:"demandSoul"`
	Rates      []int `json:"rates"`
}

// PartyMemberBonus defines party exp modifier.
type PartyMemberBonus struct {
	PartyMembers int `json:"partyMembers"`
	BonusRate    int `json:"bonusRate"`
}

// LevelDiffBonus defines exp modifier based on (playerLevel - mobLevel).
type LevelDiffBonus struct {
	LevelDiff int `json:"levelDiff"`
	BonusRate int `json:"bonusRate"`
}

// StandardDeviationHandicap defines party level deviation penalty.
type StandardDeviationHandicap struct {
	LowerBound   int `json:"lowerBound"`
	HandicapRate int `json:"handicapRate"`
}

// ExpRecalculation holds exp recalculation curves.
type ExpRecalculation struct {
	ByPartyMembers    []PartyMemberBonus          `json:"byPartyMembers"`
	ByLevelDiff       []LevelDiffBonus            `json:"byLevelDiff"`
	StandardDeviation []StandardDeviationHandicap `json:"standardDeviation"`
}

// Store holds the parsed combat and progression configurations for ZoneService.
type Store struct {
	mu               sync.RWMutex
	common           CommonConfig
	expTable         []ExpLevelEntry
	freeStat         FreeStatTable
	damageByAngle    DamageByAngle
	damageBySoul     []SoulDemandRate
	expRecalculation ExpRecalculation
}

// NewStore creates a new uninitialized Combat Store.
func NewStore() *Store {
	return &Store{}
}

// LoadDir loads all combat JSON files from the designated directory.
func (s *Store) LoadDir(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Common.json
	commonData, err := os.ReadFile(filepath.Join(dir, "Common.json"))
	if err != nil {
		return fmt.Errorf("read Common.json: %w", err)
	}
	if err := json.Unmarshal(commonData, &s.common); err != nil {
		return fmt.Errorf("unmarshal Common.json: %w", err)
	}

	// 2. ExpTable.json
	expData, err := os.ReadFile(filepath.Join(dir, "ExpTable.json"))
	if err != nil {
		return fmt.Errorf("read ExpTable.json: %w", err)
	}
	if err := json.Unmarshal(expData, &s.expTable); err != nil {
		return fmt.Errorf("unmarshal ExpTable.json: %w", err)
	}

	// 3. FreeStatTable.json
	freeStatData, err := os.ReadFile(filepath.Join(dir, "FreeStatTable.json"))
	if err != nil {
		return fmt.Errorf("read FreeStatTable.json: %w", err)
	}
	if err := json.Unmarshal(freeStatData, &s.freeStat); err != nil {
		return fmt.Errorf("unmarshal FreeStatTable.json: %w", err)
	}

	// 4. DamageByAngle.json
	angleData, err := os.ReadFile(filepath.Join(dir, "DamageByAngle.json"))
	if err != nil {
		return fmt.Errorf("read DamageByAngle.json: %w", err)
	}
	if err := json.Unmarshal(angleData, &s.damageByAngle); err != nil {
		return fmt.Errorf("unmarshal DamageByAngle.json: %w", err)
	}

	// 5. DamageBySoul.json
	soulData, err := os.ReadFile(filepath.Join(dir, "DamageBySoul.json"))
	if err != nil {
		return fmt.Errorf("read DamageBySoul.json: %w", err)
	}
	if err := json.Unmarshal(soulData, &s.damageBySoul); err != nil {
		return fmt.Errorf("unmarshal DamageBySoul.json: %w", err)
	}

	// 6. ExpRecalculation.json
	recalcData, err := os.ReadFile(filepath.Join(dir, "ExpRecalculation.json"))
	if err != nil {
		return fmt.Errorf("read ExpRecalculation.json: %w", err)
	}
	if err := json.Unmarshal(recalcData, &s.expRecalculation); err != nil {
		return fmt.Errorf("unmarshal ExpRecalculation.json: %w", err)
	}

	return nil
}

// Common returns the global character configuration constants.
func (s *Store) Common() CommonConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.common
}

// GetNextExp returns the required experience for a given level to reach the next level.
func (s *Store) GetNextExp(level int) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if level < 1 || level > len(s.expTable) {
		return 0, false
	}
	return s.expTable[level-1].NextExp, true
}

// MaxLevel returns the maximum supported level in the experience table.
func (s *Store) MaxLevel() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.expTable)
}

// GetPartyExpBonus returns the experience multiplier for party size (1000 = 100%).
func (s *Store) GetPartyExpBonus(memberCount int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, bonus := range s.expRecalculation.ByPartyMembers {
		if bonus.PartyMembers == memberCount {
			return bonus.BonusRate
		}
	}
	if memberCount >= 5 && len(s.expRecalculation.ByPartyMembers) > 0 {
		return s.expRecalculation.ByPartyMembers[len(s.expRecalculation.ByPartyMembers)-1].BonusRate
	}
	return 1000
}

// GetAngleDamageRate returns the damage multiplier based on attack angle (interpolated between thresholds).
func (s *Store) GetAngleDamageRate(isMonster bool, angle int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var records []AngleRate
	if isMonster {
		records = s.damageByAngle.Monster
	} else {
		records = s.damageByAngle.Character
	}

	if len(records) == 0 {
		return 1000
	}

	// Clamp angle to 0..180
	if angle < 0 {
		angle = -angle
	}
	if angle > 180 {
		angle = 360 - angle
	}
	if angle > 180 {
		angle = 180
	}

	for i := 0; i < len(records); i++ {
		if records[i].Angle == angle {
			return records[i].DamageRate
		}
		if i < len(records)-1 && angle > records[i].Angle && angle < records[i+1].Angle {
			// Linear interpolation
			a1, r1 := records[i].Angle, records[i].DamageRate
			a2, r2 := records[i+1].Angle, records[i+1].DamageRate
			return r1 + (r2-r1)*(angle-a1)/(a2-a1)
		}
	}

	return records[len(records)-1].DamageRate
}
