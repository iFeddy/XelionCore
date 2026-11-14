package zone

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/mobregen"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// MobTemplate unifica la definición oficial de un monstruo desde MobInfo.shn,
// MobInfoServer.shn, MobSpecies.shn y MobWeapon.shn.
type MobTemplate struct {
	ID               uint16
	InxName          string
	Name             string
	Level            uint8
	MaxHP            int
	MaxSP            int
	WalkSpeed        uint16
	RunSpeed         uint16
	IsNPC            bool
	Size             int
	GradeType        int
	MobType          int
	Visible          bool
	AC               uint32
	TB               uint32
	MR               uint32
	MB               uint32
	EnemyDetectType  int
	MobKillInx       int
	MonEXP           uint32
	EXPRange         int
	DetectRange      int
	FollowRange      int
	PeaceHPRecovDly  int
	PeaceHPRecov     int
	MobRaceType      int
	Rank             uint8
	FamilyArea       int
	FamilyRescueArea int
	FamilyRescueCnt  int
	ReturnToRegen    bool
	IsRoaming        bool
	RoamDistance     int
	RoamRestTimeMs   int
	AllCanLoot       bool
	RegenIntervalSec int
	// Arma básica (MobWeapon.shn)
	MinWC         uint32
	MaxWC         uint32
	TH            uint32
	MinMA         uint32
	MaxMA         uint32
	MH            uint32
	AttackRange   int
	AttackDelayMs int
	StaName       string
	StaStrength   int
	StaRate       int
	SpeciesParent string
	IsKillable    bool
}

// MobSkillDef almacena los atributos de una skill activa usable por monstruos desde ActiveSkill.shn + ActiveSkillInfoServer.shn.
type MobSkillDef struct {
	ID          uint16
	InxName     string
	MinWC       uint32
	MaxWC       uint32
	MinMA       uint32
	MaxMA       uint32
	Range       int
	Area        int
	DelayMs     int
	PyHitRate   int
	MaHitRate   int
	StaName     string
	StaStrength int
	StaSucRate  int
}

// MobAutoActionEntry representa una regla automática de MobAutoAction.shn.
type MobAutoActionEntry struct {
	MobInx     string
	Attack     int
	Target     int
	ActionType int
	StateInx   string
	Strength   int
	EffectRate int
	Range      int
}

// MobConditionEntry representa condiciones de MobConditionServer.shn.
type MobConditionEntry struct {
	MobInx      string
	AniLv       uint8
	MCType      int
	MCValueMin  int
	MCValueMax  int
}

// MobCatalog indexa en memoria todas las tablas .shn de monstruos y habilidades de ataque.
type MobCatalog struct {
	ByID         map[uint16]*MobTemplate
	ByInx        map[string]*MobTemplate
	KillAnnounce map[uint16]uint32 // MobID -> TextIndex
	AutoActions  map[string][]MobAutoActionEntry
	Conditions   map[string][]MobConditionEntry
	SkillsByInx  map[string]*MobSkillDef
	SkillsByID   map[uint16]*MobSkillDef
}

// BuildMobCatalog construye el índice completo de monstruos a partir de DataStore.
func BuildMobCatalog(ds *DataStore) *MobCatalog {
	cat := &MobCatalog{
		ByID:         make(map[uint16]*MobTemplate),
		ByInx:        make(map[string]*MobTemplate),
		KillAnnounce: make(map[uint16]uint32),
		AutoActions:  make(map[string][]MobAutoActionEntry),
		Conditions:   make(map[string][]MobConditionEntry),
		SkillsByInx:  make(map[string]*MobSkillDef),
		SkillsByID:   make(map[uint16]*MobSkillDef),
	}
	if ds == nil {
		return cat
	}

	// 1. MobInfo.shn
	mobInfo := ds.MobInfoSHN
	if mobInfo == nil {
		mobInfo = ds.GetSHN("MobInfo")
	}
	if mobInfo != nil {
		for _, row := range mobInfo.Rows {
			id := row.GetUint16("ID")
			inx := strings.TrimSpace(row.GetString("InxName"))
			if inx == "" {
				continue
			}
			t := &MobTemplate{
				ID:               id,
				InxName:          inx,
				Name:             strings.TrimSpace(row.GetString("Name")),
				Level:            uint8(row.GetInt("Level")),
				MaxHP:            row.GetInt("MaxHP"),
				WalkSpeed:        uint16(row.GetInt("WalkSpeed")),
				RunSpeed:         uint16(row.GetInt("RunSpeed")),
				IsNPC:            row.GetByte("IsNPC") != 0,
				Size:             row.GetInt("Size"),
				GradeType:        row.GetInt("GradeType"),
				MobType:          row.GetInt("Type"),
				MaxSP:            100,
				AttackRange:      120,
				AttackDelayMs:    1300,
				DetectRange:      264,
				FollowRange:      1500,
				EXPRange:         1500,
				RegenIntervalSec: 25,
				ReturnToRegen:    true,
				Visible:          true,
				IsKillable:       true,
			}
			if t.WalkSpeed == 0 {
				t.WalkSpeed = 30
			}
			if t.RunSpeed == 0 {
				t.RunSpeed = 105
			}
			cat.ByID[id] = t
			cat.ByInx[strings.ToLower(inx)] = t
		}
	}

	// 2. MobKillAble.shn
	killAbleMap := make(map[int]bool)
	if f := ds.GetSHN("MobKillAble"); f != nil {
		for _, row := range f.Rows {
			killInx := row.GetInt("MobKillInx")
			userKill := row.GetByte("User") != 0
			killAbleMap[killInx] = userKill
		}
	}

	// 3. MobInfoServer.shn
	if f := ds.GetSHN("MobInfoServer"); f != nil {
		for _, row := range f.Rows {
			id := uint16(row.GetInt("ID"))
			inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
			t := cat.ByInx[inx]
			if t == nil {
				t = cat.ByID[id]
			}
			if t == nil {
				continue
			}
			t.Visible = row.GetByte("Visible") != 0
			t.AC = uint32(row.GetUint16("AC"))
			t.TB = uint32(row.GetUint16("TB"))
			t.MR = uint32(row.GetUint16("MR"))
			t.MB = uint32(row.GetUint16("MB"))
			t.EnemyDetectType = row.GetInt("EnemyDetectType")
			t.MobKillInx = row.GetInt("MobKillInx")
			if uk, ok := killAbleMap[t.MobKillInx]; ok {
				t.IsKillable = uk
			}
			t.MonEXP = uint32(row.GetInt("MonEXP"))
			if r := int(row.GetUint16("EXPRange")); r > 0 {
				t.EXPRange = r
			}
			if d := int(row.GetUint16("DetectCha")); d > 0 {
				t.DetectRange = d
			}
			if fCha := row.GetInt("FollowCha"); fCha > 0 {
				t.FollowRange = fCha
			}
			t.PeaceHPRecovDly = int(row.GetUint16("PceHPRcvDly"))
			t.PeaceHPRecov = int(row.GetUint16("PceHPRcv"))
			t.MobRaceType = row.GetInt("MobRaceType")
			t.Rank = row.GetByte("Rank")
			t.FamilyArea = row.GetInt("FamilyArea")
			t.FamilyRescueArea = row.GetInt("FamilyRescArea")
			t.FamilyRescueCnt = int(row.GetByte("FamilyRescCount"))
			t.ReturnToRegen = row.GetByte("Return2Regen") != 0
			t.IsRoaming = row.GetByte("IsRoaming") != 0
			t.RoamDistance = int(row.GetUint16("RoamingDistance"))
			t.RoamRestTimeMs = int(row.GetUint16("RoamingRestTime"))
			if sp := int(row.GetUint16("MaxSP")); sp > 0 {
				t.MaxSP = sp
			}
			t.AllCanLoot = row.GetByte("AllCanLoot") != 0
			if ri := int(row.GetUint16("RegenInterval")); ri > 0 {
				t.RegenIntervalSec = ri
			}
		}
	}

	// 4. MobWeapon.shn (primera arma por mob InxName)
	if f := ds.GetSHN("MobWeapon"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
			t := cat.ByInx[inx]
			if t == nil {
				continue
			}
			minWC := uint32(row.GetInt("MinWC"))
			maxWC := uint32(row.GetInt("MaxWC"))
			minMA := uint32(row.GetInt("MinMA"))
			maxMA := uint32(row.GetInt("MaxMA"))
			if minWC > 0 || maxWC > 0 || (t.MinWC == 0 && t.MaxWC == 0) {
				t.MinWC = minWC
				t.MaxWC = maxWC
				t.TH = uint32(row.GetUint16("TH"))
				t.MinMA = minMA
				t.MaxMA = maxMA
				t.MH = uint32(row.GetUint16("MH"))
				if atkDly := int(row.GetUint16("AtkDly")); atkDly > 0 {
					t.AttackDelayMs = atkDly
				}
				if rng := int(row.GetUint16("Range")); rng > 0 {
					// En MobWeapon.shn, Range viene en unidades reducidas o directas (ej. 10 -> ~100 unidades mundo)
					if rng < 50 {
						t.AttackRange = rng * 12
					} else {
						t.AttackRange = rng
					}
				}
				sta := strings.TrimSpace(row.GetString("StaName"))
				if sta != "" && sta != "-" {
					t.StaName = sta
					t.StaStrength = int(row.GetUint16("StaStrength"))
					t.StaRate = int(row.GetUint16("StaRate"))
				}
			}
		}
	}

	// 5. MobKillAnnounce.shn
	if f := ds.GetSHN("MobKillAnnounce"); f != nil {
		for _, row := range f.Rows {
			mobID := uint16(row.GetInt("MobID"))
			textIdx := uint32(row.GetInt("TextIndex"))
			cat.KillAnnounce[mobID] = textIdx
		}
	}

	// 6. MobAutoAction.shn
	if f := ds.GetSHN("MobAutoAction"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("MobInx")))
			if inx == "" {
				continue
			}
			cat.AutoActions[inx] = append(cat.AutoActions[inx], MobAutoActionEntry{
				MobInx:     inx,
				Attack:     row.GetInt("Attack"),
				Target:     row.GetInt("Target"),
				ActionType: row.GetInt("ActionType"),
				StateInx:   strings.TrimSpace(row.GetString("StateInx")),
				Strength:   int(row.GetByte("Strength")),
				EffectRate: int(row.GetUint16("EffectRate")),
				Range:      int(row.GetUint16("Range")),
			})
		}
	}

	// 7. MobConditionServer.shn
	if f := ds.GetSHN("MobConditionServer"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("MobInx")))
			if inx == "" {
				continue
			}
			cat.Conditions[inx] = append(cat.Conditions[inx], MobConditionEntry{
				MobInx:     inx,
				AniLv:      row.GetByte("AniLv"),
				MCType:     row.GetInt("MC_Type"),
				MCValueMin: row.GetInt("MC_ValueMin"),
				MCValueMax: row.GetInt("MC_ValueMax"),
			})
		}
	}

	// 8. ActiveSkill.shn + ActiveSkillInfoServer.shn para skills de monstruos
	activeSHN := ds.ActiveSkillSHN
	if activeSHN == nil {
		activeSHN = ds.GetSHN("ActiveSkill")
	}
	if activeSHN != nil {
		for _, row := range activeSHN.Rows {
			id := row.GetUint16("ID")
			inx := strings.TrimSpace(row.GetString("InxName"))
			if inx == "" {
				continue
			}
			sk := &MobSkillDef{
				ID:          id,
				InxName:     inx,
				MinWC:       uint32(row.GetInt("MinWC")),
				MaxWC:       uint32(row.GetInt("MaxWC")),
				MinMA:       uint32(row.GetInt("MinMA")),
				MaxMA:       uint32(row.GetInt("MaxMA")),
				Range:       row.GetInt("Range"),
				Area:        row.GetInt("Area"),
				DelayMs:     row.GetInt("DlyTime"),
				StaName:     strings.TrimSpace(row.GetString("StaNameA")),
				StaStrength: row.GetInt("StaStrengthA"),
				StaSucRate:  row.GetInt("StaSucRateA"),
				PyHitRate:   950,
				MaHitRate:   950,
			}
			cat.SkillsByID[id] = sk
			cat.SkillsByInx[strings.ToLower(inx)] = sk
		}
	}
	if f := ds.GetSHN("ActiveSkillInfoServer"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
			if sk := cat.SkillsByInx[inx]; sk != nil {
				if p := row.GetInt("SkilPyHitRate"); p > 0 {
					sk.PyHitRate = p
				}
				if m := row.GetInt("SkilMaHitRate"); m > 0 {
					sk.MaHitRate = m
				}
			}
		}
	}

	return cat
}

// LookupTemplate busca una plantilla oficial por InxName (o nombre legible).
func (c *MobCatalog) LookupTemplate(nameOrInx string) *MobTemplate {
	if c == nil {
		return nil
	}
	k := strings.ToLower(strings.TrimSpace(nameOrInx))
	return c.ByInx[k]
}

// MobRegenGroup mantiene el estado vivo de un grupo de data/MobRegen/<Map>.json.
type MobRegenGroup struct {
	GroupID string
	Def     mobregen.Group
}

// SetDataStore vincula el DataStore de datos oficiales al MapManager e inicializa el MobCatalog.
func (mm *MapManager) SetDataStore(ds *DataStore) {
	mm.mu.Lock()
	mm.dataStore = ds
	if ds != nil {
		if ds.MobCatalog == nil {
			ds.MobCatalog = BuildMobCatalog(ds)
		}
		if ds.TitleCatalog == nil {
			ds.TitleCatalog = BuildTitleCatalog(ds)
		}
		if ds.CollectCatalog == nil {
			ds.CollectCatalog = BuildCollectCatalog(ds)
		}
	}
	maps := make([]*ZoneMap, 0, len(mm.maps))
	for _, zMap := range mm.maps {
		maps = append(maps, zMap)
	}
	mm.mu.Unlock()

	for _, zMap := range maps {
		zMap.SetDataStore(ds)
	}
}

// SetHandler registra la referencia al Handler principal para callbacks de drops/combat/AI.
func (mm *MapManager) SetHandler(h *Handler) {
	mm.mu.Lock()
	mm.handler = h
	maps := make([]*ZoneMap, 0, len(mm.maps))
	for _, zMap := range mm.maps {
		maps = append(maps, zMap)
	}
	mm.mu.Unlock()

	for _, zMap := range maps {
		zMap.SetHandler(h)
	}
}

// SetDataStore configura los catálogos en un ZoneMap e instancia los spawns iniciales de MobRegen.
func (m *ZoneMap) SetDataStore(ds *DataStore) {
	m.mu.Lock()
	m.dataStore = ds
	if ds != nil && ds.MobCatalog != nil {
		m.mobCatalog = ds.MobCatalog
	}
	alreadyLoaded := len(m.regenGroups) > 0
	m.mu.Unlock()

	if ds != nil && !alreadyLoaded {
		m.LoadMobsFromRegen(ds)
	}
}

// SetHandler registra el Handler propietario en el ZoneMap.
func (m *ZoneMap) SetHandler(h *Handler) {
	m.mu.Lock()
	m.handler = h
	m.mu.Unlock()
}

// LoadMobsFromRegen carga todos los grupos de data/MobRegen/<MapID>.json y genera las instancias vivas de ZoneMonster.
func (m *ZoneMap) LoadMobsFromRegen(ds *DataStore) int {
	if m == nil || ds == nil || ds.MobRegen == nil {
		return 0
	}
	var mr *mobregen.MapRegen
	for k, v := range ds.MobRegen {
		if strings.EqualFold(k, m.MapID) {
			mr = v
			break
		}
	}
	if mr == nil || len(mr.Groups) == 0 {
		return 0
	}

	if ds.MobCatalog == nil {
		ds.MobCatalog = BuildMobCatalog(ds)
	}

	m.mu.Lock()
	m.dataStore = ds
	m.mobCatalog = ds.MobCatalog
	m.regenGroups = make(map[string]*MobRegenGroup, len(mr.Groups))
	m.mu.Unlock()

	spawnedCount := 0
	for groupID, grp := range mr.Groups {
		rg := &MobRegenGroup{GroupID: groupID, Def: grp}
		m.mu.Lock()
		m.regenGroups[groupID] = rg
		m.mu.Unlock()

		if !grp.Enabled {
			continue
		}

		for _, entry := range grp.Mobs {
			if !entry.Enabled || entry.Count <= 0 {
				continue
			}
			tmpl := ds.MobCatalog.LookupTemplate(entry.Mob)
			if tmpl == nil {
				continue
			}

			// Calcular intervalo de regeneración en segundos desde entrada de regen o MobInfoServer.
			regenSec := entry.RegStandard
			if regenSec <= 0 {
				regenSec = tmpl.RegenIntervalSec
			}
			if regenSec <= 0 {
				regenSec = 25
			}

			for i := 0; i < entry.Count; i++ {
				sx, sy := m.pickWalkableCoordInGroup(grp)
				dir := int16(rand.Intn(256))
				mob := m.SpawnMobFromTemplate(tmpl, sx, sy, dir, groupID, regenSec, false)
				if mob != nil {
					if grp.IsFamily && mob.FamilyRescueCnt == 0 {
						mob.FamilyRescueCnt = 2
						if mob.FamilyRescueArea == 0 {
							mob.FamilyRescueArea = 500
						}
					}
					spawnedCount++
				}
			}
		}
	}

	// Inicializar Lua Behavior/AIScript bridge del mapa si aplica.
	m.initMapLuaAI()

	log.Printf("[Zone][Map %s] Spawned %d live monsters from %d MobRegen groups.", m.MapID, spawnedCount, len(mr.Groups))
	return spawnedCount
}

func (m *ZoneMap) pickWalkableCoordInGroup(grp mobregen.Group) (int32, int32) {
	halfW := grp.Width / 2
	halfH := grp.Height / 2
	for attempt := 0; attempt < 8; attempt++ {
		x := grp.CenterX
		y := grp.CenterY
		if halfW > 0 {
			x += rand.Intn(halfW*2+1) - halfW
		}
		if halfH > 0 {
			y += rand.Intn(halfH*2+1) - halfH
		}
		if m.IsPositionWalkable(float64(x), float64(y)) {
			return int32(x), int32(y)
		}
	}
	if m.blockInfo != nil {
		tx, ty := blockinfo.WorldToTile(float64(grp.CenterX), float64(grp.CenterY))
		if nx, ny, ok := m.blockInfo.FindNearestWalkable(tx, ty, 6); ok {
			wx, wy := blockinfo.TileToWorld(nx, ny)
			return int32(wx), int32(wy)
		}
	}
	return int32(grp.CenterX), int32(grp.CenterY)
}

// IsPositionWalkable valida unas coordenadas mundo contra .shbd y puertas cerradas de .sbi.
func (m *ZoneMap) IsPositionWalkable(wx, wy float64) bool {
	if wx < 0 || wy < 0 {
		return false
	}
	if m.blockInfo != nil {
		if m.blockInfo.IsWorldMoveBlocked(wx, wy) {
			return false
		}
	}
	if m.doors != nil {
		tx, ty := blockinfo.WorldToTile(wx, wy)
		if m.doors.IsDoorMoveBlocked(tx, ty) {
			return false
		}
	}
	return true
}

// SpawnMobFromTemplate crea e inserta un ZoneMonster desde MobTemplate, enlazando MobRoam, MobSetting y MobAttackSequence.
func (m *ZoneMap) SpawnMobFromTemplate(tmpl *MobTemplate, x, y int32, dir int16, groupID string, regenSec int, broadcast bool) *ZoneMonster {
	if m == nil || tmpl == nil {
		return nil
	}

	m.mu.Lock()
	handle := m.nextObjID
	m.nextObjID++
	ds := m.dataStore
	m.mu.Unlock()

	minWC := tmpl.MinWC
	maxWC := tmpl.MaxWC
	if maxWC == 0 {
		maxWC = uint32(tmpl.Level)*4 + 6
		minWC = maxWC * 7 / 10
	}
	hp := tmpl.MaxHP
	if hp <= 0 {
		hp = int(tmpl.Level)*20 + 25
	}

	mob := NewZoneMonster(
		handle,
		tmpl.ID,
		tmpl.Name,
		m.MapID,
		x, y,
		dir,
		tmpl.Level,
		hp,
		tmpl.MaxSP,
		tmpl.AC,
		tmpl.MR,
		minWC,
		maxWC,
		tmpl.TH,
		tmpl.TB,
		tmpl.MonEXP,
	)
	mob.InxName = tmpl.InxName
	mob.MinMA = tmpl.MinMA
	mob.MaxMA = tmpl.MaxMA
	mob.RegenGroupID = groupID
	if regenSec > 0 {
		mob.RegenIntervalSec = regenSec
	} else {
		mob.RegenIntervalSec = tmpl.RegenIntervalSec
	}
	mob.EXPRange = tmpl.EXPRange
	mob.WalkSpeed = tmpl.WalkSpeed
	mob.RunSpeed = tmpl.RunSpeed
	mob.DetectRange = tmpl.DetectRange
	mob.FollowRange = tmpl.FollowRange
	mob.AttackRange = tmpl.AttackRange
	mob.AttackDelayMs = tmpl.AttackDelayMs
	mob.EnemyDetectType = tmpl.EnemyDetectType
	mob.ReturnToRegen = tmpl.ReturnToRegen
	mob.AllCanLoot = tmpl.AllCanLoot
	mob.Visible = tmpl.Visible
	mob.IsKillableByUser = tmpl.IsKillable
	mob.Rank = tmpl.Rank
	mob.GradeType = tmpl.GradeType
	mob.MobType = tmpl.MobType
	mob.RaceType = tmpl.MobRaceType
	mob.FamilyArea = tmpl.FamilyArea
	mob.FamilyRescueArea = tmpl.FamilyRescueArea
	mob.FamilyRescueCnt = tmpl.FamilyRescueCnt
	mob.PeaceHPRecovDly = tmpl.PeaceHPRecovDly
	mob.PeaceHPRecov = tmpl.PeaceHPRecov
	mob.IsRoaming = tmpl.IsRoaming
	mob.RoamDistance = tmpl.RoamDistance
	mob.RoamRestTimeMs = tmpl.RoamRestTimeMs

	if ds != nil {
		if ds.MobAttack != nil {
			if ai, ok := ds.MobAttack.GetMobAI(tmpl.InxName); ok {
				mob.MobAI = ai
			}
		}
		if ds.MobSetting != nil {
			if trigs := ds.MobSetting.GetTriggers(tmpl.InxName); len(trigs) > 0 {
				mob.Triggers = trigs
			}
		}
		if ds.MobRoam != nil {
			if rt, ok := ds.MobRoam.GetRoute(tmpl.InxName); ok && len(rt.Points) > 0 {
				mob.Route = rt
				mob.AIState = MobAIRoam
			}
		}
	}
	if mob.Route == nil && mob.IsRoaming {
		mob.AIState = MobAIRoam
	}

	m.AddMonster(mob)
	if broadcast {
		m.BroadcastMobSpawn(mob, true)
	}
	return mob
}

// SpawnMonsterByInx genera un monstruo dinámico (por script Lua, regen de test, o evento).
func (m *ZoneMap) SpawnMonsterByInx(inxName string, x, y int32, broadcast bool) (*ZoneMonster, error) {
	m.mu.RLock()
	cat := m.mobCatalog
	ds := m.dataStore
	m.mu.RUnlock()

	if cat == nil && ds != nil {
		cat = BuildMobCatalog(ds)
	}
	if cat == nil {
		return nil, fmt.Errorf("mob catalog unavailable")
	}
	tmpl := cat.LookupTemplate(inxName)
	if tmpl == nil {
		return nil, fmt.Errorf("monster '%s' not found in MobInfo.shn", inxName)
	}
	mob := m.SpawnMobFromTemplate(tmpl, x, y, 0, "", tmpl.RegenIntervalSec, broadcast)
	return mob, nil
}

// BroadcastMobSpawn difunde PROTO_NC_BRIEFINFO_REGENMOB_CMD (0x1C08, 149B) al área AOI.
func (m *ZoneMap) BroadcastMobSpawn(mob *ZoneMonster, regenAni bool) {
	if m == nil || mob == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD)
	var ani byte
	if regenAni {
		ani = 1
	}
	mob.WriteRegenMob(pw, ani)
	mx, my := mob.GetPosition()
	m.grid.BroadcastAOI(pw, mx, my, 0)
}

// SendMobs envía un lote de monstruos al cliente mediante PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09), máx 28 (0x1C) por paquete.
func (m *ZoneMap) SendMobs(p *ZonePlayer, mobs []*ZoneMonster) {
	if p == nil || p.Client == nil || len(mobs) == 0 {
		return
	}
	for i := 0; i < len(mobs); i += BriefInfoMaxBatch {
		end := i + BriefInfoMaxBatch
		if end > len(mobs) {
			end = len(mobs)
		}
		chunk := mobs[i:end]
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_MOB_CMD)
		pw.WriteByte(byte(len(chunk)))
		for _, mob := range chunk {
			mob.WriteRegenMob(pw, 0)
		}
		_ = p.Client.SendPacket(pw)
	}
}
