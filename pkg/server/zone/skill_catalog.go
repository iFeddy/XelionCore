package zone

import (
	"encoding/binary"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// SkillAbStateEntry representa uno de los 4 estados anormales (StaNameA..D) que una skill puede aplicar.
type SkillAbStateEntry struct {
	Name     string
	Strength int
	SucRate  int
}

// SkillSpecialEntry representa un par SpecialIndex/SpecialValue de ActiveSkill.shn (A..E).
type SkillSpecialEntry struct {
	Index uint32
	Value int
}

// SkillDef representa una habilidad activa cargada desde ActiveSkill.shn + ActiveSkillInfoServer.shn
// con paridad de campos respecto a SkillDataBox::SkillDataIndex en Zone.exe.
type SkillDef struct {
	ID               uint16
	InxName          string
	Name             string
	DemandType       uint32
	SP               int
	SPRate           int
	HP               int
	HPRate           int
	LP               int
	Range            int
	First            uint32
	Last             uint32
	IsMovingSkill    bool
	UsableDegree     uint16
	DirectionRotate  uint16
	SkillDegree      uint16
	SkillTargetState uint32
	CastTimeMs       int
	DlyTimeMs        int
	DlyGroupNum      int
	DlyTimeGroupMs   int

	// Escalado físico (RulesOfEngagementPhisycalSkill::roe_AttackPower 0x00506B40)
	MinWC     int
	MinWCRate int
	MaxWC     int
	MaxWCRate int

	// Escalado mágico (RulesOfEngagementMagicalSkill::roe_AttackPower 0x00506ED0)
	MinMA     int
	MinMARate int
	MaxMA     int
	MaxMARate int

	Area         int
	TargetNumber int
	UseClass     uint32

	AbStates [4]SkillAbStateEntry

	// Tablas de empoderamiento (ActiveSkill.shn nT0..nT3, 5 niveles: índice 0..4 => empower 1..5)
	// Track 0 (nT0): Damage / Heal bonus (Zone.exe +0x1BB)
	// Track 1 (nT1): SP reduction (Zone.exe +0x1CF)
	// Track 2 (nT2): AbState duration bonus in ms (Zone.exe +0x1E3)
	// Track 3 (nT3): Cooldown reduction in ms (Zone.exe +0x1F7)
	EmpowerDamage   [5]int
	EmpowerSP       [5]int
	EmpowerKeepTime [5]int
	EmpowerCoolTime [5]int

	EffectType   uint32
	Specials     [5]SkillSpecialEntry
	CannotInside bool
	DemandSoul   uint8
	HitID        uint16

	// Campos de ActiveSkillInfoServer.shn
	PyHitRate     int
	MaHitRate     int
	SkillHitType  uint32
	ItemUseSkill  bool
	AggroPerDmg   int
	AbsoluteAggro int
	AddSoul       uint8
}

// IsMagical indica si la skill utiliza escala mágica (MinMA/MaxMA o MinMARate/MaxMARate).
func (s *SkillDef) IsMagical() bool {
	if s == nil {
		return false
	}
	return s.MinMA > 0 || s.MaxMA > 0 || s.MinMARate > 0 || s.MaxMARate > 0
}

// EmpowerDamageBonus retorna el bono de daño/curación del track 0 (nivel 1..5).
func (s *SkillDef) EmpowerDamageBonus(level uint8) int {
	if s == nil || level == 0 {
		return 0
	}
	if level > 5 {
		level = 5
	}
	return s.EmpowerDamage[level-1]
}

// EmpowerSPReduction retorna la reducción de SP del track 1 (nivel 1..5).
func (s *SkillDef) EmpowerSPReduction(level uint8) int {
	if s == nil || level == 0 {
		return 0
	}
	if level > 5 {
		level = 5
	}
	return s.EmpowerSP[level-1]
}

// EmpowerKeepTimeBonus retorna el bono de duración de AbState del track 2 (nivel 1..5).
func (s *SkillDef) EmpowerKeepTimeBonus(level uint8) int {
	if s == nil || level == 0 {
		return 0
	}
	if level > 5 {
		level = 5
	}
	return s.EmpowerKeepTime[level-1]
}

// EmpowerCoolTimeReduction retorna la reducción de cooldown del track 3 (nivel 1..5).
func (s *SkillDef) EmpowerCoolTimeReduction(level uint8) int {
	if s == nil || level == 0 {
		return 0
	}
	if level > 5 {
		level = 5
	}
	return s.EmpowerCoolTime[level-1]
}

// HealAmount retorna el valor base de curación de HP si la skill tiene efecto de curación.
// En ActiveSkill.shn, SpecialIndex == 1 corresponde a HP Heal (SpecialValue = cantidad base).
func (s *SkillDef) HealAmount() int {
	if s == nil {
		return 0
	}
	for _, sp := range s.Specials {
		if sp.Index == 1 && sp.Value > 0 {
			return sp.Value
		}
	}
	return 0
}

// IsTeleportSkill indica si la skill desplaza/teletransporta al lanzador hacia el objetivo o coordenada.
func (s *SkillDef) IsTeleportSkill() bool {
	if s == nil {
		return false
	}
	for _, sp := range s.Specials {
		// En Fiesta SHN: índices de teleport/blink/charge (ej. 12 = MoveToTarget / Teleport)
		if sp.Index == 12 {
			return true
		}
	}
	return false
}

// SkillCatalog indexa en memoria todas las skills activas desde ActiveSkill.shn + ActiveSkillInfoServer.shn.
type SkillCatalog struct {
	ByID  map[uint16]*SkillDef
	ByInx map[string]*SkillDef
}

// Get retorna la definición oficial de una skill activa por ID.
func (c *SkillCatalog) Get(id uint16) *SkillDef {
	if c == nil || c.ByID == nil {
		return nil
	}
	return c.ByID[id]
}

// GetByInx retorna la definición oficial de una skill activa por InxName (case-insensitive).
func (c *SkillCatalog) GetByInx(inx string) *SkillDef {
	if c == nil || c.ByInx == nil {
		return nil
	}
	return c.ByInx[strings.ToLower(strings.TrimSpace(inx))]
}

func readSHNIntAt(row *shn.Row, idx int) int {
	b := row.GetBytesAt(idx)
	switch len(b) {
	case 4:
		return int(int32(binary.LittleEndian.Uint32(b)))
	case 2:
		return int(int16(binary.LittleEndian.Uint16(b)))
	case 1:
		return int(b[0])
	default:
		return 0
	}
}

// BuildSkillCatalog construye SkillCatalog reutilizando los *shn.File ya cargados en DataStore.
func BuildSkillCatalog(ds *DataStore) *SkillCatalog {
	cat := &SkillCatalog{
		ByID:  make(map[uint16]*SkillDef),
		ByInx: make(map[string]*SkillDef),
	}
	if ds == nil {
		return cat
	}
	activeSHN := ds.ActiveSkillSHN
	if activeSHN == nil {
		activeSHN = ds.GetSHN("ActiveSkill")
	}
	if activeSHN == nil {
		return cat
	}

	// Localizar índices base de las tablas de empower (nT0, nT1, nT2, nT3: 5 columnas consecutivas cada una)
	idxT0 := activeSHN.ColumnIndex("nT0")
	idxT1 := activeSHN.ColumnIndex("nT1")
	idxT2 := activeSHN.ColumnIndex("nT2")
	idxT3 := activeSHN.ColumnIndex("nT3")

	for i := range activeSHN.Rows {
		row := &activeSHN.Rows[i]
		id := row.GetUint16("ID")
		inx := strings.TrimSpace(row.GetString("InxName"))
		if id == 0 && inx == "" {
			continue
		}

		def := &SkillDef{
			ID:               id,
			InxName:          inx,
			Name:             strings.TrimSpace(row.GetString("Name")),
			DemandType:       row.GetUint32("DemandType"),
			SP:               row.GetInt("SP"),
			SPRate:           row.GetInt("SPRate"),
			HP:               row.GetInt("HP"),
			HPRate:           row.GetInt("HPRate"),
			LP:               row.GetInt("LP"),
			Range:            row.GetInt("Range"),
			First:            row.GetUint32("First"),
			Last:             row.GetUint32("Last"),
			IsMovingSkill:    row.GetByte("IsMovingSkill") != 0,
			UsableDegree:     row.GetUint16("UsableDegree"),
			DirectionRotate:  row.GetUint16("DirectionRotate"),
			SkillDegree:      row.GetUint16("SkillDegree"),
			SkillTargetState: row.GetUint32("SkillTargetState"),
			CastTimeMs:       row.GetInt("CastTime"),
			DlyTimeMs:        row.GetInt("DlyTime"),
			DlyGroupNum:      row.GetInt("DlyGroupNum"),
			DlyTimeGroupMs:   row.GetInt("DlyTimeGroup"),
			MinWC:            row.GetInt("MinWC"),
			MinWCRate:        row.GetInt("MinWCRate"),
			MaxWC:            row.GetInt("MaxWC"),
			MaxWCRate:        row.GetInt("MaxWCRate"),
			MinMA:            row.GetInt("MinMA"),
			MinMARate:        row.GetInt("MinMARate"),
			MaxMA:            row.GetInt("MaxMA"),
			MaxMARate:        row.GetInt("MaxMARate"),
			Area:             row.GetInt("Area"),
			TargetNumber:     row.GetInt("TargetNumber"),
			UseClass:         row.GetUint32("UseClass"),
			EffectType:       row.GetUint32("EffectType"),
			CannotInside:     row.GetByte("CannotInside") != 0,
			DemandSoul:       row.GetByte("DemandSoul"),
			HitID:            row.GetUint16("HitID"),
			PyHitRate:        950,
			MaHitRate:        950,
		}

		def.AbStates = [4]SkillAbStateEntry{
			{Name: strings.TrimSpace(row.GetString("StaNameA")), Strength: row.GetInt("StaStrengthA"), SucRate: row.GetInt("StaSucRateA")},
			{Name: strings.TrimSpace(row.GetString("StaNameB")), Strength: row.GetInt("StaStrengthB"), SucRate: row.GetInt("StaSucRateB")},
			{Name: strings.TrimSpace(row.GetString("StaNameC")), Strength: row.GetInt("StaStrengthC"), SucRate: row.GetInt("StaSucRateC")},
			{Name: strings.TrimSpace(row.GetString("StaNameD")), Strength: row.GetInt("StaStrengthD"), SucRate: row.GetInt("StaSucRateD")},
		}

		def.Specials = [5]SkillSpecialEntry{
			{Index: row.GetUint32("SpecialIndexA"), Value: row.GetInt("SpecialValueA")},
			{Index: row.GetUint32("SpecialIndexB"), Value: row.GetInt("SpecialValueB")},
			{Index: row.GetUint32("SpecialIndexC"), Value: row.GetInt("SpecialValueC")},
			{Index: row.GetUint32("SpecialIndexD"), Value: row.GetInt("SpecialValueD")},
			{Index: row.GetUint32("SpecialIndexE"), Value: row.GetInt("SpecialValueE")},
		}

		if idxT0 >= 0 {
			for k := 0; k < 5; k++ {
				def.EmpowerDamage[k] = readSHNIntAt(row, idxT0+k)
			}
		}
		if idxT1 >= 0 {
			for k := 0; k < 5; k++ {
				def.EmpowerSP[k] = readSHNIntAt(row, idxT1+k)
			}
		}
		if idxT2 >= 0 {
			for k := 0; k < 5; k++ {
				def.EmpowerKeepTime[k] = readSHNIntAt(row, idxT2+k)
			}
		}
		if idxT3 >= 0 {
			for k := 0; k < 5; k++ {
				def.EmpowerCoolTime[k] = readSHNIntAt(row, idxT3+k)
			}
		}

		cat.ByID[id] = def
		if inx != "" {
			cat.ByInx[strings.ToLower(inx)] = def
		}
	}

	if srvSHN := ds.GetSHN("ActiveSkillInfoServer"); srvSHN != nil {
		for i := range srvSHN.Rows {
			row := &srvSHN.Rows[i]
			id := row.GetUint16("ID")
			def := cat.ByID[id]
			if def == nil {
				inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
				def = cat.ByInx[inx]
			}
			if def == nil {
				continue
			}
			if p := row.GetInt("SkilPyHitRate"); p > 0 {
				def.PyHitRate = p
			}
			if m := row.GetInt("SkilMaHitRate"); m > 0 {
				def.MaHitRate = m
			}
			def.SkillHitType = row.GetUint32("SkillHitType")
			def.ItemUseSkill = row.GetByte("ItremUseSkill") != 0
			def.AggroPerDmg = row.GetInt("AggroPerDamage")
			def.AbsoluteAggro = row.GetInt("AbsoluteAggro")
			def.AddSoul = row.GetByte("AddSoul")
		}
	}

	return cat
}
