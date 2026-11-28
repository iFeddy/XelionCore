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
	DemandLv     int  // Nivel requerido oficial (ItemInfo.shn scroll DemandLv / ActiveSkillView.shn uiDemandLv)
	HasScroll    bool // true si existe libro/scroll de aprendizaje (Class == 11) en ItemInfo.shn (Zone.exe ac_AllLearnSKill VA 0x004213AF)

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

// IsReviveSkill indica si la skill apunta a jugadores muertos (SkillTargetState == 1, ej. Revival01 ID 1620).
func (s *SkillDef) IsReviveSkill() bool {
	if s == nil {
		return false
	}
	return s.SkillTargetState == 1
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

// HasAnyAbState indica si la skill declara al menos un AbState (StaNameA..D).
func (s *SkillDef) HasAnyAbState() bool {
	if s == nil {
		return false
	}
	for _, st := range s.AbStates {
		if st.Name != "" && st.Name != "-" && st.SucRate > 0 {
			return true
		}
	}
	return false
}

// MagicFieldDurationMs retorna la duración en ms si la skill crea una zona persistente en el suelo
// (ActiveSkill.shn: SpecialIndex == 13 con duración en Value, o EffectType == 1 como Inferno01/FrostNova01).
func (s *SkillDef) MagicFieldDurationMs() int {
	if s == nil {
		return 0
	}
	for _, sp := range s.Specials {
		if sp.Index == 13 && sp.Value > 0 {
			return sp.Value
		}
	}
	if s.EffectType == 1 {
		return 3000
	}
	return 0
}

// CureDispelSpec retorna si la skill cura debuffs en aliados (SpecialIndex 2 o 9, EffectType 8)
// o disipa buffs en enemigos (SpecialIndex 8 o 26), junto con el máximo de estados a remover (0 = todos).
func (s *SkillDef) CureDispelSpec() (isCure bool, isDispelEnemy bool, maxRemove int) {
	if s == nil {
		return false, false, 0
	}
	for _, sp := range s.Specials {
		switch sp.Index {
		case 2, 9:
			return true, false, sp.Value
		case 8, 26:
			return false, true, sp.Value
		}
	}
	if s.EffectType == 8 {
		return true, false, 1
	}
	return false, false, 0
}

// IsPureBuff indica si la skill aplica únicamente AbStates o Cure/Dispel sin daño ni curación directa.
func (s *SkillDef) IsPureBuff() bool {
	if s == nil || s.HealAmount() > 0 {
		return false
	}
	if s.MinWC > 0 || s.MaxWC > 0 || s.MinWCRate > 0 || s.MaxWCRate > 0 ||
		s.MinMA > 0 || s.MaxMA > 0 || s.MinMARate > 0 || s.MaxMARate > 0 {
		return false
	}
	isCure, isDispel, _ := s.CureDispelSpec()
	return s.HasAnyAbState() || isCure || isDispel
}

// AbStateDef representa una entrada oficial de Abstate.shn (todos los 792 estados anormales del cliente).
// En Fiesta Online, el índice oficial de protocolo (enum ABSTATEINDEX en 0x2426/0x2428/0x2429 y bitfield 99B)
// es la columna AbStataIndex (0..791), NO la columna ID (que llega a 8880+).
type AbStateDef struct {
	ID          uint16 // Abstate.shn AbStataIndex (0..791, enum ABSTATEINDEX oficial)
	RawID       uint16 // Abstate.shn ID (identificador interno de fila SHN)
	InxName     string
	SubAbState  string
	PartyStates [5]string
	PartyRange  int
	DispelIndex int // 1 = buff removible, 2 = debuff
	SaveType    int // Abstate.shn AbStateSaveType (!= 0 => persistir en tCharacterState al desconectar)
}

// SubAbStateAction representa un par ActionIndex/ActionArg de SubAbstate.shn (A..D).
type SubAbStateAction struct {
	Index int
	Arg   int
}

// SubAbStateDef representa una fila de SubAbstate.shn indexada por (InxName, Strength).
type SubAbStateDef struct {
	InxName    string
	Strength   int
	KeepTimeMs int
	Actions    [4]SubAbStateAction
}

// PassiveSkillDef representa una habilidad pasiva oficial cargada desde PassiveSkill.shn + PassiveSkillView.shn + ItemInfo.shn.
type PassiveSkillDef struct {
	ID              uint16 // ID base de PassiveSkill.shn (0..842, enviado al cliente en 0x103E / 0x4804 / 0x4827)
	DBID            uint16 // ID almacenado en tSkill.nSkillNo = PassiveSkillIDOffset + ID (>= 30000)
	InxName         string
	Name            string
	UseClass        uint32
	DemandLv        int
	HasScroll       bool
	MstRt           [12]int // Permille de maestría por tipo de arma (1..11)
	MstPl           [12]int // Bono plano de ataque por tipo de arma (1..11)
	TB              int
	MaxSP           int
	WCRateUp        int
	MARateUp        int
	HealUpRate      int
	DMGMinusRate    int
	CastingTimeRate int // Permille (1000 = 100%)
}

// PassiveBonusSummary resume los bonificadores efectivos de todas las pasivas aprendidas por un jugador.
type PassiveBonusSummary struct {
	WeaponRatePermille int // MstRt del arma equipada (permille)
	WeaponFlatPlus     int // MstPl del arma equipada (plano)
	WCRateUp           int
	MARateUp           int
	HealUpRate         int
	DMGMinusRate       int
	TB                 int
	MaxSP              int
	CastingTimeRate    int // Permille (1000 = sin reducción)
}

func weaponTypeToMasterySlot(wt uint8) int {
	switch wt {
	case 1: // Sword1
		return 1
	case 13, 16: // Hammer1
		return 2
	case 21: // Sword2
		return 3
	case 4: // Axe2
		return 4
	case 5: // Mace1
		return 5
	case 2: // Bow2
		return 6
	case 10: // CrossBow2
		return 7
	case 11: // Wand2
		return 8
	case 3: // Staff2
		return 9
	case 17, 18: // Claw
		return 10
	case 19: // DSword
		return 11
	default:
		return 0
	}
}

// SkillCatalog indexa en memoria todas las skills activas y pasivas desde ActiveSkill.shn, ActiveSkillInfoServer.shn,
// ActiveSkillView.shn, PassiveSkill.shn, PassiveSkillView.shn, ItemInfo.shn y UseClassTypeInfo.shn.
type SkillCatalog struct {
	ByID            map[uint16]*SkillDef
	ByInx           map[string]*SkillDef
	PassivesByID    map[uint16]*PassiveSkillDef // índice: ID base de PassiveSkill.shn (0..842)
	ItemWeaponType  map[uint16]uint8            // índice: ItemInfo.shn ID -> WeaponType
	UseClassMatrix  map[uint32][28]bool         // índice 1..27 = Job ID oficial (Fig=1 .. Sav=27)
	AbStateByInx    map[string]*AbStateDef      // índice: lower(InxName) desde Abstate.shn
	AbStateByIndex  map[uint16]*AbStateDef      // índice: AbStataIndex (0..791) desde Abstate.shn
	SubAbStateByKey map[string]*SubAbStateDef   // índice: lower(InxName)#Strength desde SubAbstate.shn
}

// ComputePassiveBonuses calcula los bonificadores efectivos de las skills pasivas aprendidas por el jugador
// cruzando el tipo de arma equipada en ItemInfo.shn con PassiveSkill.shn.
func (c *SkillCatalog) ComputePassiveBonuses(p *ZonePlayer) PassiveBonusSummary {
	out := PassiveBonusSummary{CastingTimeRate: 1000}
	if c == nil || p == nil || len(p.LearnedSkills) == 0 || len(c.PassivesByID) == 0 {
		return out
	}

	// Determinar el slot de maestría del arma equipada
	masterySlot := 0
	if p.Char != nil && p.Char.Equips != nil && len(c.ItemWeaponType) > 0 {
		for _, itemID := range p.Char.Equips {
			if itemID == 0 || itemID == 0xFFFF {
				continue
			}
			if wt := c.ItemWeaponType[itemID]; wt > 0 {
				if slot := weaponTypeToMasterySlot(wt); slot > 0 {
					masterySlot = slot
					break
				}
			}
		}
	}

	for skillNo := range p.LearnedSkills {
		if skillNo < PassiveSkillIDOffset {
			continue
		}
		ps := c.PassivesByID[skillNo-PassiveSkillIDOffset]
		if ps == nil {
			continue
		}
		if masterySlot > 0 {
			if ps.MstRt[masterySlot] > out.WeaponRatePermille {
				out.WeaponRatePermille = ps.MstRt[masterySlot]
			}
			if ps.MstPl[masterySlot] > out.WeaponFlatPlus {
				out.WeaponFlatPlus = ps.MstPl[masterySlot]
			}
		}
		if ps.WCRateUp > out.WCRateUp {
			out.WCRateUp = ps.WCRateUp
		}
		if ps.MARateUp > out.MARateUp {
			out.MARateUp = ps.MARateUp
		}
		if ps.HealUpRate > out.HealUpRate {
			out.HealUpRate = ps.HealUpRate
		}
		if ps.DMGMinusRate > out.DMGMinusRate {
			out.DMGMinusRate = ps.DMGMinusRate
		}
		if ps.TB > out.TB {
			out.TB = ps.TB
		}
		if ps.MaxSP > out.MaxSP {
			out.MaxSP = ps.MaxSP
		}
		if ps.CastingTimeRate > 0 && ps.CastingTimeRate < out.CastingTimeRate {
			out.CastingTimeRate = ps.CastingTimeRate
		}
	}
	return out
}

func subAbStateKey(inx string, strength int) string {
	var buf [12]byte
	i := len(buf)
	v := strength
	if v == 0 {
		i--
		buf[i] = '0'
	} else {
		neg := v < 0
		if neg {
			v = -v
		}
		for v > 0 {
			i--
			buf[i] = byte('0' + v%10)
			v /= 10
		}
		if neg {
			i--
			buf[i] = '-'
		}
	}
	return strings.ToLower(strings.TrimSpace(inx)) + "#" + string(buf[i:])
}

// LookupAbState busca un AbState oficial por InxName (case-insensitive) en Abstate.shn.
func (c *SkillCatalog) LookupAbState(name string) *AbStateDef {
	if c == nil || c.AbStateByInx == nil {
		return nil
	}
	return c.AbStateByInx[strings.ToLower(strings.TrimSpace(name))]
}

// LookupAbStateByIndex busca un AbState oficial por AbStataIndex (0..791) en Abstate.shn.
func (c *SkillCatalog) LookupAbStateByIndex(idx uint16) *AbStateDef {
	if c == nil || c.AbStateByIndex == nil {
		return nil
	}
	return c.AbStateByIndex[idx]
}

// LookupSubAbState busca la entrada de SubAbstate.shn por (InxName, Strength) con fallback a 999 y 1.
func (c *SkillCatalog) LookupSubAbState(inxName string, strength int) *SubAbStateDef {
	if c == nil || c.SubAbStateByKey == nil || inxName == "" || inxName == "-" {
		return nil
	}
	if s := c.SubAbStateByKey[subAbStateKey(inxName, strength)]; s != nil {
		return s
	}
	if s := c.SubAbStateByKey[subAbStateKey(inxName, 999)]; s != nil {
		return s
	}
	return c.SubAbStateByKey[subAbStateKey(inxName, 1)]
}

// CanClassUseSkill verifica según UseClassTypeInfo.shn si una clase (job 1..27) puede usar/aprender una skill con dado UseClass.
func (c *SkillCatalog) CanClassUseSkill(useClass uint32, job byte) bool {
	if c == nil || job < 1 || job > 27 {
		return false
	}
	if len(c.UseClassMatrix) == 0 {
		return false
	}
	row, ok := c.UseClassMatrix[useClass]
	if !ok {
		return false
	}
	return row[job]
}

// IsRealClassActiveSkill retorna true solo para skills activas reales de clase (UseClass > 1, con scroll Class=11 en ItemInfo.shn)
// habilitadas para el job indicado. Excluye skills globales (UseClass=1: Ball*, Recipe_*, SummonMate01) y DT_*.
func (c *SkillCatalog) IsRealClassActiveSkill(sk *SkillDef, job byte) bool {
	if c == nil || sk == nil || sk.UseClass <= 1 || !sk.HasScroll {
		return false
	}
	return c.CanClassUseSkill(sk.UseClass, job)
}

// IsRealClassPassiveSkill retorna true solo para skills pasivas reales de clase (UseClass > 1, con scroll Class=11 en ItemInfo.shn)
// habilitadas para el job indicado. Excluye pasivas globales (UseClass=1 como PowerofLove01).
func (c *SkillCatalog) IsRealClassPassiveSkill(ps *PassiveSkillDef, job byte) bool {
	if c == nil || ps == nil || ps.UseClass <= 1 || !ps.HasScroll {
		return false
	}
	return c.CanClassUseSkill(ps.UseClass, job)
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
		ByID:            make(map[uint16]*SkillDef),
		ByInx:           make(map[string]*SkillDef),
		PassivesByID:    make(map[uint16]*PassiveSkillDef),
		ItemWeaponType:  make(map[uint16]uint8),
		UseClassMatrix:  make(map[uint32][28]bool),
		AbStateByInx:    make(map[string]*AbStateDef),
		AbStateByIndex:  make(map[uint16]*AbStateDef),
		SubAbStateByKey: make(map[string]*SubAbStateDef),
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

	// Cargar uiDemandLv desde ActiveSkillView.shn (reutilizado desde ds.SHNTables)
	if viewSHN := ds.GetSHN("ActiveSkillView"); viewSHN != nil {
		for i := range viewSHN.Rows {
			row := &viewSHN.Rows[i]
			id := row.GetUint16("ID")
			def := cat.ByID[id]
			if def == nil {
				inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
				def = cat.ByInx[inx]
			}
			if def != nil {
				def.DemandLv = row.GetInt("uiDemandLv")
			}
		}
	}

	// Cargar matriz oficial UseClassTypeInfo.shn (columnas 1..27 corresponden 1:1 a Job 1..27)
	if ucSHN := ds.GetSHN("UseClassTypeInfo"); ucSHN != nil {
		for i := range ucSHN.Rows {
			row := &ucSHN.Rows[i]
			uc := row.GetUint32("UseClass")
			var flags [28]bool
			for job := 1; job <= 27 && job < row.ColumnCount(); job++ {
				flags[job] = row.GetByteAt(job) != 0
			}
			cat.UseClassMatrix[uc] = flags
		}
	}

	// Indexar libros/scrolls de habilidades (Class == 11) y WeaponType desde ItemInfo.shn
	scrollDemandByInx := make(map[string]int)
	if itemSHN := ds.GetSHN("ItemInfo"); itemSHN != nil {
		for i := range itemSHN.Rows {
			row := &itemSHN.Rows[i]
			if wt := row.GetByte("WeaponType"); wt > 0 {
				cat.ItemWeaponType[row.GetUint16("ID")] = wt
			}
			if row.GetInt("Class") != 11 {
				continue
			}
			inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
			if inx == "" {
				continue
			}
			scrollDemandByInx[inx] = row.GetInt("DemandLv")
		}
	}
	if len(scrollDemandByInx) > 0 {
		for _, def := range cat.ByID {
			if def == nil || def.InxName == "" {
				continue
			}
			if dlv, ok := scrollDemandByInx[strings.ToLower(def.InxName)]; ok {
				def.HasScroll = true
				if dlv > 0 {
					def.DemandLv = dlv
				}
			}
		}
	}

	// Indexar habilidades pasivas reales desde PassiveSkill.shn + PassiveSkillView.shn
	passiveSHN := ds.PassiveSkillSHN
	if passiveSHN == nil {
		passiveSHN = ds.GetSHN("PassiveSkill")
	}
	if passiveSHN != nil {
		for i := range passiveSHN.Rows {
			row := &passiveSHN.Rows[i]
			id := row.GetUint16("ID")
			inx := strings.TrimSpace(row.GetString("InxName"))
			if id == 0 && inx == "" {
				continue
			}
			pdef := &PassiveSkillDef{
				ID:              id,
				DBID:            PassiveSkillIDOffset + id,
				InxName:         inx,
				Name:            strings.TrimSpace(row.GetString("Name")),
				TB:              row.GetInt("TB"),
				MaxSP:           row.GetInt("MaxSP"),
				WCRateUp:        row.GetInt("WCRateUp"),
				MARateUp:        row.GetInt("MARateUp"),
				HealUpRate:      row.GetInt("HealUPRate"),
				DMGMinusRate:    row.GetInt("DMG_MinusRate"),
				CastingTimeRate: row.GetInt("CastingTime"),
				MstRt: [12]int{
					0,
					row.GetInt("MstRtSword1"),
					row.GetInt("MstRtHammer1"),
					row.GetInt("MstRtSword2"),
					row.GetInt("MstRtAxe2"),
					row.GetInt("MstRtMace1"),
					row.GetInt("MstRtBow2"),
					row.GetInt("MstRtCrossBow2"),
					row.GetInt("MstRtWand2"),
					row.GetInt("MstRtStaff2"),
					row.GetInt("MstRtClaw"),
					row.GetInt("MstRtDSword"),
				},
				MstPl: [12]int{
					0,
					row.GetInt("MstPlSword1"),
					row.GetInt("MstPlHammer1"),
					row.GetInt("MstPlSword2"),
					row.GetInt("MstPlAxe2"),
					row.GetInt("MstPlMace1"),
					row.GetInt("MstPlBow2"),
					row.GetInt("MstPlCrossBow2"),
					row.GetInt("MstPlWand2"),
					row.GetInt("MstPlStaff2"),
					row.GetInt("MstPlClaw"),
					row.GetInt("MstPlDSword"),
				},
			}
			if dlv, ok := scrollDemandByInx[strings.ToLower(inx)]; ok {
				pdef.HasScroll = true
				pdef.DemandLv = dlv
			}
			cat.PassivesByID[id] = pdef
		}
	}
	if pvSHN := ds.GetSHN("PassiveSkillView"); pvSHN != nil {
		for i := range pvSHN.Rows {
			row := &pvSHN.Rows[i]
			id := row.GetUint16("ID")
			pdef := cat.PassivesByID[id]
			if pdef == nil {
				continue
			}
			pdef.UseClass = row.GetUint32("UseClass")
			if pdef.DemandLv == 0 {
				pdef.DemandLv = row.GetInt("DemandLv")
			}
		}
	}

	// Indexar todos los AbStates oficiales desde Abstate.shn + SubAbstate.shn
	if abSHN := ds.GetSHN("Abstate"); abSHN != nil {
		for i := range abSHN.Rows {
			row := &abSHN.Rows[i]
			inx := strings.TrimSpace(row.GetString("InxName"))
			if inx == "" || inx == "-" {
				continue
			}
			def := &AbStateDef{
				ID:          row.GetUint16("AbStataIndex"),
				RawID:       row.GetUint16("ID"),
				InxName:     inx,
				SubAbState:  strings.TrimSpace(row.GetString("SubAbState")),
				PartyRange:  row.GetInt("PartyRange"),
				DispelIndex: row.GetInt("DispelIndex"),
				SaveType:    row.GetInt("AbStateSaveType"),
				PartyStates: [5]string{
					strings.TrimSpace(row.GetString("PartyState1")),
					strings.TrimSpace(row.GetString("PartyState2")),
					strings.TrimSpace(row.GetString("PartyState3")),
					strings.TrimSpace(row.GetString("PartyState4")),
					strings.TrimSpace(row.GetString("PartyState5")),
				},
			}
			cat.AbStateByInx[strings.ToLower(inx)] = def
			cat.AbStateByIndex[def.ID] = def
		}
	}
	if subSHN := ds.GetSHN("SubAbstate"); subSHN != nil {
		for i := range subSHN.Rows {
			row := &subSHN.Rows[i]
			inx := strings.TrimSpace(row.GetString("InxName"))
			if inx == "" || inx == "-" {
				continue
			}
			str := row.GetInt("Strength")
			cat.SubAbStateByKey[subAbStateKey(inx, str)] = &SubAbStateDef{
				InxName:    inx,
				Strength:   str,
				KeepTimeMs: row.GetInt("KeepTime"),
				Actions: [4]SubAbStateAction{
					{Index: row.GetInt("ActionIndexA"), Arg: row.GetInt("ActionArgA")},
					{Index: row.GetInt("ActionIndexB"), Arg: row.GetInt("ActionArgB")},
					{Index: row.GetInt("ActionIndexC"), Arg: row.GetInt("ActionArgC")},
					{Index: row.GetInt("ActionIndexD"), Arg: row.GetInt("ActionArgD")},
				},
			}
		}
	}

	return cat
}
