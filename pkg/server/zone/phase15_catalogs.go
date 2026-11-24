package zone

import (
	"encoding/binary"
	"log"
	"path/filepath"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// loadSHNWithFallback intenta cargar una tabla desde DataStore RAM y, si no está en los 180
// archivos base, la lee desde disco (data/ o data/View/) y la cachea inmediatamente en ds.SHNTables.
func (ds *DataStore) loadSHNWithFallback(dataDir string, names ...string) *shn.File {
	for _, name := range names {
		if ds != nil {
			if f := ds.GetSHN(name); f != nil {
				return f
			}
		}
		if dataDir != "" {
			for _, sub := range []string{name + ".shn", filepath.Join("View", name+".shn")} {
				if f, err := shn.ReadFile(filepath.Join(dataDir, sub)); err == nil {
					if ds != nil {
						if ds.SHNTables == nil {
							ds.SHNTables = make(map[string]*shn.File)
						}
						base := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
						ds.SHNTables[base] = f
						ds.SHNTables[strings.ToLower(base)] = f
					}
					return f
				}
			}
		}
	}
	return nil
}

// QuestMobReq define un requisito de cacería de una misión (QuestData.shn).
type QuestMobReq struct {
	Enabled        bool
	MobOrSpeciesID uint16
	IsSpecies      bool
	Count          uint8
}

// QuestItemReq define un requisito de recolección de ítem (QuestData.shn).
type QuestItemReq struct {
	Enabled bool
	ItemID  uint16
	Count   uint16
}

// QuestDropDef define un drop especial de misión cuando un mob muere.
type QuestDropDef struct {
	MobID    uint16
	ItemID   uint16
	DropRate uint16
}

// QuestRewardDef define un ítem de recompensa fijo (Mode=1) o seleccionable (Mode=2).
type QuestRewardDef struct {
	Mode   uint8 // 1=Fixed, 2=Select
	Type   uint8 // 2=Item
	ItemID uint16
	Count  uint16
}

// QuestDef representa una misión completa parseada desde QuestData.shn + QuestDialog.shn + QuestScript.shn.
type QuestDef struct {
	QuestID        uint16
	TitleDialogID  uint32
	DescDialogID   uint32
	Title          string
	Desc           string
	CheckLevel     bool
	MinLevel       uint8
	MaxLevel       uint8
	Enabled        bool
	StartNPCID     uint16
	CheckPrevQuest bool
	PrevQuestID    uint16
	EndNPCID       uint16
	MobKills       []QuestMobReq
	ItemCollects   []QuestItemReq
	Drops          []QuestDropDef
	ExpEnabled     bool
	RewardExp      uint32
	MoneyEnabled   bool
	RewardMoney    uint32
	Rewards        []QuestRewardDef
	ScriptID       uint16
}

// QuestCatalog almacena las misiones, diálogos, scripts QSC y especies de QuestSpecies.shn.
type QuestCatalog struct {
	ByID         map[uint16]*QuestDef
	Dialogs      map[uint32]string
	Scripts      map[uint16]string
	SpeciesMobs  map[uint16][]uint16 // SpeciesID -> []MobID
	MobToSpecies map[uint16][]uint16 // MobID -> []SpeciesID
}

// GetQuest retorna la definición oficial de una misión.
func (c *QuestCatalog) GetQuest(id uint16) *QuestDef {
	if c == nil {
		return nil
	}
	return c.ByID[id]
}

// GetDialog retorna un diálogo oficial de QuestDialog.shn.
func (c *QuestCatalog) GetDialog(id uint32) string {
	if c == nil {
		return ""
	}
	return c.Dialogs[id]
}

// GetScript retorna el script QSC oficial de QuestScript.shn.
func (c *QuestCatalog) GetScript(id uint16) string {
	if c == nil {
		return ""
	}
	return c.Scripts[id]
}

// BuildQuestCatalog indexa QuestData.shn, QuestDialog.shn, QuestScript.shn y QuestSpecies.shn.
func BuildQuestCatalog(ds *DataStore, dataDir string) *QuestCatalog {
	cat := &QuestCatalog{
		ByID:         make(map[uint16]*QuestDef),
		Dialogs:      make(map[uint32]string),
		Scripts:      make(map[uint16]string),
		SpeciesMobs:  make(map[uint16][]uint16),
		MobToSpecies: make(map[uint16][]uint16),
	}

	if f := ds.loadSHNWithFallback(dataDir, "QuestDialog"); f != nil {
		for _, r := range f.Rows {
			id := r.GetUint32At(0)
			txt := r.GetString("Dialog")
			if txt == "" && r.ColumnCount() > 1 {
				txt = strings.TrimRight(string(r.GetBytesAt(1)), "\x00")
			}
			if id > 0 {
				cat.Dialogs[id] = txt
			}
		}
	}

	if f := ds.loadSHNWithFallback(dataDir, "QuestScript"); f != nil {
		for _, r := range f.Rows {
			id := uint16(r.GetIntAt(0))
			scr := r.GetString("Script")
			if scr == "" && r.ColumnCount() > 1 {
				scr = strings.TrimRight(string(r.GetBytesAt(1)), "\x00")
			}
			if id > 0 {
				cat.Scripts[id] = scr
			}
		}
	}

	if f := ds.loadSHNWithFallback(dataDir, "QuestSpecies"); f != nil {
		for _, r := range f.Rows {
			spID := uint16(r.GetIntAt(0))
			for col := 1; col < r.ColumnCount(); col++ {
				mobName := strings.TrimSpace(strings.TrimRight(string(r.GetBytesAt(col)), "\x00"))
				if mobName == "" || mobName == "-" {
					continue
				}
				var mobID uint16
				if ds.MobCatalog != nil {
					if tmpl := ds.MobCatalog.LookupTemplate(mobName); tmpl != nil {
						mobID = tmpl.ID
					}
				}
				if mobID == 0 {
					mobID = uint16(r.GetIntAt(col))
				}
				if mobID > 0 {
					cat.SpeciesMobs[spID] = append(cat.SpeciesMobs[spID], mobID)
					cat.MobToSpecies[mobID] = append(cat.MobToSpecies[mobID], spID)
				}
			}
		}
	}

	if f := ds.loadSHNWithFallback(dataDir, "QuestData"); f != nil {
		for _, r := range f.Rows {
			hdr := r.GetBytes("HeaderData")
			if len(hdr) < 680 {
				continue
			}
			qid := uint16(binary.LittleEndian.Uint32(hdr[0x04:0x08]))
			if qid == 0 {
				continue
			}
			titleID := binary.LittleEndian.Uint32(hdr[0x08:0x0C])
			descID := binary.LittleEndian.Uint32(hdr[0x0C:0x10])
			q := &QuestDef{
				QuestID:        qid,
				TitleDialogID:  titleID,
				DescDialogID:   descID,
				Title:          cat.Dialogs[titleID],
				Desc:           cat.Dialogs[descID],
				CheckLevel:     hdr[0x1A] != 0,
				MinLevel:       hdr[0x1B],
				MaxLevel:       hdr[0x1C],
				Enabled:        hdr[0x1D] != 0,
				StartNPCID:     binary.LittleEndian.Uint16(hdr[0x1E:0x20]),
				CheckPrevQuest: binary.LittleEndian.Uint16(hdr[0x38:0x3A]) != 0,
				PrevQuestID:    binary.LittleEndian.Uint16(hdr[0x3A:0x3C]),
				EndNPCID:       binary.LittleEndian.Uint16(hdr[0x5E:0x60]),
				ExpEnabled:     binary.LittleEndian.Uint16(hdr[0x204:0x206]) != 0,
				RewardExp:      binary.LittleEndian.Uint32(hdr[0x208:0x20C]),
				MoneyEnabled:   hdr[0x210] != 0,
				RewardMoney:    binary.LittleEndian.Uint32(hdr[0x214:0x218]),
				ScriptID:       qid,
			}
			// MobKill slots (5 slots a 6 bytes desde 0x64)
			for i := 0; i < 5; i++ {
				off := 0x64 + i*6
				if off+6 <= len(hdr) {
					en := binary.LittleEndian.Uint16(hdr[off:off+2]) != 0
					mobOrSp := binary.LittleEndian.Uint16(hdr[off+2 : off+4])
					isSp := hdr[off+4] != 0
					cnt := hdr[off+5]
					if en && mobOrSp > 0 && cnt > 0 {
						q.MobKills = append(q.MobKills, QuestMobReq{
							Enabled:        en,
							MobOrSpeciesID: mobOrSp,
							IsSpecies:      isSp,
							Count:          cnt,
						})
					}
				}
			}
			// ItemCollect slots (5 slots a 6 bytes desde 0x84)
			for i := 0; i < 5; i++ {
				off := 0x84 + i*6
				if off+6 <= len(hdr) {
					en := binary.LittleEndian.Uint16(hdr[off:off+2]) != 0
					itemID := binary.LittleEndian.Uint16(hdr[off+2 : off+4])
					cnt := binary.LittleEndian.Uint16(hdr[off+4 : off+6])
					if en && itemID > 0 && cnt > 0 {
						q.ItemCollects = append(q.ItemCollects, QuestItemReq{
							Enabled: en,
							ItemID:  itemID,
							Count:   cnt,
						})
					}
				}
			}
			// QuestDrop slots (hasta 10 slots de 8 bytes desde 0xC0)
			for i := 0; i < 10; i++ {
				off := 0xC0 + i*8
				if off+6 <= 0x1E0 {
					mobID := binary.LittleEndian.Uint16(hdr[off : off+2])
					itemID := binary.LittleEndian.Uint16(hdr[off+2 : off+4])
					rate := binary.LittleEndian.Uint16(hdr[off+4 : off+6])
					if mobID > 0 && itemID > 0 {
						q.Drops = append(q.Drops, QuestDropDef{MobID: mobID, ItemID: itemID, DropRate: rate})
					}
				}
			}
			// Recompensas (stride 12 bytes desde 0x21C)
			for off := 0x21C; off+8 <= len(hdr); off += 12 {
				mode := hdr[off]
				typ := hdr[off+1]
				itemID := binary.LittleEndian.Uint16(hdr[off+4 : off+6])
				cnt := binary.LittleEndian.Uint16(hdr[off+6 : off+8])
				if mode != 0 && itemID > 0 {
					if cnt == 0 {
						cnt = 1
					}
					q.Rewards = append(q.Rewards, QuestRewardDef{
						Mode:   mode,
						Type:   typ,
						ItemID: itemID,
						Count:  cnt,
					})
				}
			}
			cat.ByID[qid] = q
		}
	}
	return cat
}

// KQDef define una Kingdom Quest desde KingdomQuest.shn y KingdomQuestMap.shn.
type KQDef struct {
	ID            uint16
	Title         string
	LimitTimeSec  int
	StartWaitSec  int
	MinLevel      uint8
	MaxLevel      uint8
	MinPlayers    int
	MaxPlayers    int
	DemandQuest   uint16
	DemandItem    uint16
	DemandMobKill uint16
	RewardIndex   string
	ScriptName    string
	BaseMap       string
	Maps          []string
}

// KQRewardItem define una recompensa posible dentro de KingdomQuestRew.shn.
type KQRewardItem struct {
	ItemInx string
	Rate    int
}

// KQRewardGroup agrupa las recompensas de KingdomQuestRew.shn por IndexString.
type KQRewardGroup struct {
	IndexString string
	BoxItemInx  string
	Items       []KQRewardItem
}

// KQCatalog almacena las definiciones de Kingdom Quest y sus recompensas.
type KQCatalog struct {
	ByID        map[uint16]*KQDef
	ByMap       map[string]*KQDef
	RewardByKey map[string]*KQRewardGroup
}

// BuildKQCatalog carga KingdomQuest.shn, KingdomQuestMap.shn y KingdomQuestRew.shn.
func BuildKQCatalog(ds *DataStore, dataDir string) *KQCatalog {
	cat := &KQCatalog{
		ByID:        make(map[uint16]*KQDef),
		ByMap:       make(map[string]*KQDef),
		RewardByKey: make(map[string]*KQRewardGroup),
	}

	mapsByID := make(map[uint16][]string)
	if fm := ds.loadSHNWithFallback(dataDir, "KingdomQuestMap"); fm != nil {
		for _, r := range fm.Rows {
			id := uint16(r.GetIntAt(0))
			if id == 0 {
				id = r.GetUint16("ID")
			}
			var mList []string
			for _, col := range fm.Columns {
				if strings.HasPrefix(col.Name, "Map") || col.Name == "BaseMap" {
					mName := strings.TrimSpace(r.GetString(col.Name))
					if mName != "" && mName != "-" {
						mList = append(mList, mName)
					}
				}
			}
			if id > 0 && len(mList) > 0 {
				mapsByID[id] = mList
			}
		}
	}

	if fkq := ds.loadSHNWithFallback(dataDir, "KingdomQuest"); fkq != nil {
		for _, r := range fkq.Rows {
			id := uint16(r.GetInt("ID"))
			if id == 0 {
				id = uint16(r.GetIntAt(0))
			}
			if id == 0 {
				continue
			}
			def := &KQDef{
				ID:            id,
				Title:         strings.TrimSpace(r.GetString("Title")),
				LimitTimeSec:  r.GetInt("LimitTime"),
				StartWaitSec:  r.GetInt("StartWaitTime"),
				MinLevel:      uint8(r.GetInt("MinLevel")),
				MaxLevel:      uint8(r.GetInt("MaxLevel")),
				MinPlayers:    r.GetInt("MinPlayers"),
				MaxPlayers:    r.GetInt("MaxPlayers"),
				DemandQuest:   uint16(r.GetInt("DemandQuest")),
				DemandItem:    uint16(r.GetInt("DemandItem")),
				DemandMobKill: uint16(r.GetInt("DemandMobKill")),
				RewardIndex:   strings.TrimSpace(r.GetString("RewardIndex")),
				ScriptName:    strings.TrimSpace(r.GetString("ScriptLanguage")),
				Maps:          mapsByID[id],
			}
			if def.LimitTimeSec <= 0 {
				def.LimitTimeSec = 1800
			}
			if len(def.Maps) > 0 {
				def.BaseMap = def.Maps[0]
				for _, m := range def.Maps {
					cat.ByMap[strings.ToLower(m)] = def
				}
			}
			if def.ScriptName != "" && !strings.Contains(def.ScriptName, "/") {
				cat.ByMap[strings.ToLower(def.ScriptName)] = def
			}
			cat.ByID[id] = def
		}
	}

	if fr := ds.loadSHNWithFallback(dataDir, "KingdomQuestRew"); fr != nil {
		for _, r := range fr.Rows {
			idxStr := strings.TrimSpace(r.GetString("IndexString"))
			if idxStr == "" {
				continue
			}
			grp := &KQRewardGroup{
				IndexString: idxStr,
				BoxItemInx:  strings.TrimSpace(r.GetString("KQBoxItemIDX")),
			}
			for _, col := range fr.Columns {
				if strings.HasPrefix(col.Name, "Reward") && !strings.HasPrefix(col.Name, "RewardRate") {
					itemInx := strings.TrimSpace(r.GetString(col.Name))
					if itemInx != "" && itemInx != "-" {
						suffix := strings.TrimPrefix(col.Name, "Reward")
						rate := r.GetInt("RewardRate" + suffix)
						if rate <= 0 {
							rate = 100
						}
						grp.Items = append(grp.Items, KQRewardItem{ItemInx: itemInx, Rate: rate})
					}
				}
			}
			cat.RewardByKey[strings.ToLower(idxStr)] = grp
		}
	}
	return cat
}

// GTSkillDef define un skill especial del Guild Tournament (GuildTournamentSkill.shn).
type GTSkillDef struct {
	Index      uint32
	DeathPoint uint32
	StaName    string
	DelayTime  uint32
}

// GTRewardDef define una recompensa de Guild Tournament (GuildTournamentReward.shn).
type GTRewardDef struct {
	Rank        uint8
	RewardGroup uint16
	RewardType  uint16
	Value1      uint32
	Value2      uint32
	Value3      uint32
}

// GTCatalog almacena todas las tablas GuildTournament*.shn.
type GTCatalog struct {
	MinLv       uint8
	MinMem      uint8
	JoinMoney   uint64
	OccupySec   uint32
	OccupyScore uint32
	Skills      map[uint32]*GTSkillDef
	Rewards     []GTRewardDef
}

// BuildGTCatalog carga las tablas GuildTournament*.shn.
func BuildGTCatalog(ds *DataStore, dataDir string) *GTCatalog {
	cat := &GTCatalog{
		MinLv:       1,
		MinMem:      1,
		OccupySec:   10,
		OccupyScore: 100,
		Skills:      make(map[uint32]*GTSkillDef),
	}
	if fr := ds.loadSHNWithFallback(dataDir, "GuildTournamentRequire"); fr != nil && len(fr.Rows) > 0 {
		r := fr.Rows[0]
		if v := r.GetInt("MinLv"); v > 0 {
			cat.MinLv = uint8(v)
		}
		if v := r.GetInt("MinMem"); v > 0 {
			cat.MinMem = uint8(v)
		}
		cat.JoinMoney = uint64(r.GetInt("JoinMoney"))
	}
	if fo := ds.loadSHNWithFallback(dataDir, "GuildTournamentOccupy"); fo != nil && len(fo.Rows) > 0 {
		r := fo.Rows[0]
		if v := r.GetInt("OccupyTime"); v > 0 {
			cat.OccupySec = uint32(v)
		}
		if v := r.GetInt("Score"); v > 0 {
			cat.OccupyScore = uint32(v)
		}
	}
	if fs := ds.loadSHNWithFallback(dataDir, "GuildTournamentSkill"); fs != nil {
		for _, r := range fs.Rows {
			idx := uint32(r.GetInt("Index"))
			if idx == 0 {
				idx = uint32(r.GetIntAt(0))
			}
			if idx > 0 {
				cat.Skills[idx] = &GTSkillDef{
					Index:      idx,
					DeathPoint: uint32(r.GetInt("DeathPoint")),
					StaName:    strings.TrimSpace(r.GetString("StaName")),
					DelayTime:  uint32(r.GetInt("DlyTime")),
				}
			}
		}
	}
	if frw := ds.loadSHNWithFallback(dataDir, "GuildTournamentReward"); frw != nil {
		for _, r := range frw.Rows {
			cat.Rewards = append(cat.Rewards, GTRewardDef{
				Rank:        uint8(r.GetInt("Rank")),
				RewardGroup: uint16(r.GetInt("RewardGroup")),
				RewardType:  uint16(r.GetInt("RewardType")),
				Value1:      uint32(r.GetInt("Value1")),
				Value2:      uint32(r.GetInt("Value2")),
				Value3:      uint32(r.GetInt("Value3")),
			})
		}
	}
	return cat
}

// MiniHouseDef define un tipo de MiniHouse desde MiniHouse.shn.
type MiniHouseDef struct {
	Handle     uint16
	ItemID     uint16
	KeepHours  uint32
	HPTick     uint32
	SPTick     uint32
	HPRecovery uint32
	SPRecovery uint32
	Casting    uint32
	Slots      uint16
}

// MiniHouseFurnitureDef define un mueble de MiniHouseFurniture.shn.
type MiniHouseFurnitureDef struct {
	ItemID      uint16
	Endure      uint32
	ObjEffectID uint32
}

// MiniHouseCatalog almacena definiciones de MiniHouse.shn, Furniture, Endure, ObjEffect y Ani.
type MiniHouseCatalog struct {
	ByItemID    map[uint16]*MiniHouseDef
	ByHandle    map[uint16]*MiniHouseDef
	Furniture   map[uint16]*MiniHouseFurnitureDef
	EndureByID  map[uint16]uint32
	ValidEmotes map[uint16]bool
}

// BuildMiniHouseCatalog indexa todas las tablas MiniHouse*.shn.
func BuildMiniHouseCatalog(ds *DataStore, dataDir string) *MiniHouseCatalog {
	cat := &MiniHouseCatalog{
		ByItemID:    make(map[uint16]*MiniHouseDef),
		ByHandle:    make(map[uint16]*MiniHouseDef),
		Furniture:   make(map[uint16]*MiniHouseFurnitureDef),
		EndureByID:  make(map[uint16]uint32),
		ValidEmotes: make(map[uint16]bool),
	}
	resolveItemInx := func(r shn.Row) uint16 {
		itemInx := strings.TrimSpace(r.GetString("ItemID"))
		if itemInx != "" && ds.ItemStore != nil {
			if info, ok := ds.ItemStore.GetItemInfoByName(itemInx); ok && info != nil {
				return info.ID
			}
		}
		if v := uint16(r.GetInt("ItemID")); v > 0 {
			return v
		}
		return 0
	}
	if f := ds.loadSHNWithFallback(dataDir, "MiniHouse"); f != nil {
		for _, r := range f.Rows {
			hnd := uint16(r.GetInt("Handle"))
			itemID := resolveItemInx(r)
			def := &MiniHouseDef{
				Handle:     hnd,
				ItemID:     itemID,
				KeepHours:  uint32(r.GetInt("KeepTime_Hour")),
				HPTick:     uint32(r.GetInt("HPTick")),
				SPTick:     uint32(r.GetInt("SPTick")),
				HPRecovery: uint32(r.GetInt("HPRecovery")),
				SPRecovery: uint32(r.GetInt("SPRecovery")),
				Casting:    uint32(r.GetInt("Casting")),
				Slots:      uint16(r.GetInt("Slot")),
			}
			if def.ItemID > 0 {
				cat.ByItemID[def.ItemID] = def
			}
			cat.ByHandle[def.Handle] = def
		}
	}
	if fe := ds.loadSHNWithFallback(dataDir, "MiniHouseEndure"); fe != nil {
		for _, r := range fe.Rows {
			id := resolveItemInx(r)
			if id == 0 {
				id = uint16(r.GetIntAt(0))
			}
			val := uint32(r.GetInt("Endure"))
			if val == 0 {
				val = uint32(r.GetInt("MaxEndure"))
			}
			if id > 0 {
				cat.EndureByID[id] = val
			}
		}
	}
	if ff := ds.loadSHNWithFallback(dataDir, "MiniHouseFurniture"); ff != nil {
		for _, r := range ff.Rows {
			id := resolveItemInx(r)
			if id == 0 {
				id = uint16(r.GetIntAt(0))
			}
			if id > 0 {
				endure := uint32(r.GetInt("Endure"))
				if endure == 0 {
					endure = cat.EndureByID[id]
				}
				cat.Furniture[id] = &MiniHouseFurnitureDef{
					ItemID:      id,
					Endure:      endure,
					ObjEffectID: uint32(r.GetInt("EffectID")),
				}
			}
		}
	}
	if fa := ds.loadSHNWithFallback(dataDir, "MiniHouseObjAni"); fa != nil {
		for _, r := range fa.Rows {
			id := uint16(r.GetIntAt(0))
			if id > 0 {
				cat.ValidEmotes[id] = true
			}
		}
	}
	return cat
}

// PetActionDef representa una acción de mascota en GTIServer.shn.
type PetActionDef struct {
	ID          uint32
	InxName     string
	ActionType  uint8
	TargetValue int
}

// PetCatalog indexa GTIServer.shn, GTIView.shn, GTIBreedSubject.shn y GTIGetRate.shn.
type PetCatalog struct {
	ActionsByInx map[string]*PetActionDef
	ViewNames    map[uint32]string
	GetRates     map[uint32]int
}

// BuildPetCatalog carga las tablas GTI*.shn.
func BuildPetCatalog(ds *DataStore, dataDir string) *PetCatalog {
	cat := &PetCatalog{
		ActionsByInx: make(map[string]*PetActionDef),
		ViewNames:    make(map[uint32]string),
		GetRates:     make(map[uint32]int),
	}
	if f := ds.loadSHNWithFallback(dataDir, "GTIServer"); f != nil {
		for _, r := range f.Rows {
			id := uint32(r.GetInt("ID"))
			inx := strings.TrimSpace(r.GetString("InxName"))
			if inx != "" {
				cat.ActionsByInx[strings.ToLower(inx)] = &PetActionDef{
					ID:          id,
					InxName:     inx,
					ActionType:  uint8(r.GetInt("GTIActionType")),
					TargetValue: r.GetInt("Value"),
				}
			}
		}
	}
	if fv := ds.loadSHNWithFallback(dataDir, "GTIView"); fv != nil {
		for _, r := range fv.Rows {
			id := uint32(r.GetInt("ID"))
			name := strings.TrimSpace(r.GetString("Name"))
			if id > 0 && name != "" {
				cat.ViewNames[id] = name
			}
		}
	}
	if fg := ds.loadSHNWithFallback(dataDir, "GTIGetRate"); fg != nil {
		for _, r := range fg.Rows {
			id := uint32(r.GetIntAt(0))
			if id > 0 {
				cat.GetRates[id] = r.GetIntAt(1)
			}
		}
	}
	return cat
}

// WeaponTitleDef define un rango de título de arma en WeaponTitleData.shn.
type WeaponTitleDef struct {
	MobID        uint16
	Level        uint8
	Prefix       uint16
	Suffix       uint16
	MobKillCount uint32
}

// WeaponTitleCatalog almacena WeaponTitleData.shn.
type WeaponTitleCatalog struct {
	ByKey  map[uint32]*WeaponTitleDef // (MobID<<8)|Level
	Rows   []*WeaponTitleDef
}

// BuildWeaponTitleCatalog carga WeaponTitleData.shn.
func BuildWeaponTitleCatalog(ds *DataStore, dataDir string) *WeaponTitleCatalog {
	cat := &WeaponTitleCatalog{
		ByKey: make(map[uint32]*WeaponTitleDef),
	}
	if f := ds.loadSHNWithFallback(dataDir, "WeaponTitleData"); f != nil {
		for _, r := range f.Rows {
			def := &WeaponTitleDef{
				MobID:        uint16(r.GetInt("MobID")),
				Level:        uint8(r.GetInt("Level")),
				Prefix:       uint16(r.GetInt("Prefix")),
				Suffix:       uint16(r.GetInt("Suffix")),
				MobKillCount: uint32(r.GetInt("MobKillCount")),
			}
			key := (uint32(def.MobID) << 8) | uint32(def.Level)
			cat.ByKey[key] = def
			cat.Rows = append(cat.Rows, def)
		}
	}
	return cat
}

// DiceGameDef define una regla de apuesta de DiceGame / GBDiceGame.
type DiceGameDef struct {
	BetType  uint8
	MinBet   uint32
	MaxBet   uint32
	RateMult float64
}

// SlotBetCoinDef define las monedas permitidas en GBSMBetCoin.shn.
type SlotBetCoinDef struct {
	CoinTier uint8
	Coins    uint32
	Mult     uint32
}

// CasinoCatalog almacena Dice*.shn, BelongDice.shn y todas las tablas GB*.shn (Taisai + SlotMachine).
type CasinoCatalog struct {
	DiceGames      []DiceGameDef
	TaxRatePct     int
	MaxExchange    uint32
	SlotBetCoins   []SlotBetCoinDef
	SlotAllRate    uint32
	SlotLineRate   uint32
	SlotCenterRate uint32
	JackpotBasePct uint32
}

// BuildCasinoCatalog carga DiceGame.shn, BelongDice.shn y GB*.shn.
func BuildCasinoCatalog(ds *DataStore, dataDir string) *CasinoCatalog {
	cat := &CasinoCatalog{
		TaxRatePct:     5,
		MaxExchange:    100000,
		SlotAllRate:    50,
		SlotLineRate:   10,
		SlotCenterRate: 25,
		JackpotBasePct: 5,
	}
	if fd := ds.loadSHNWithFallback(dataDir, "DiceGame", "GBDiceGame"); fd != nil {
		for _, r := range fd.Rows {
			betType := uint8(r.GetIntAt(0))
			minB := uint32(r.GetIntAt(1))
			maxB := uint32(r.GetIntAt(2))
			mult := float64(r.GetIntAt(3))
			if mult <= 0 {
				mult = 2.0
			}
			cat.DiceGames = append(cat.DiceGames, DiceGameDef{
				BetType:  betType,
				MinBet:   minB,
				MaxBet:   maxB,
				RateMult: mult,
			})
		}
	}
	if ft := ds.loadSHNWithFallback(dataDir, "GBTaxRate"); ft != nil && len(ft.Rows) > 0 {
		if v := ft.Rows[0].GetIntAt(0); v >= 0 && v <= 50 {
			cat.TaxRatePct = v
		}
	}
	if fm := ds.loadSHNWithFallback(dataDir, "GBExchangeMaxCoin"); fm != nil && len(fm.Rows) > 0 {
		if v := uint32(fm.Rows[0].GetIntAt(0)); v > 0 {
			cat.MaxExchange = v
		}
	}
	if fb := ds.loadSHNWithFallback(dataDir, "GBSMBetCoin"); fb != nil {
		for _, r := range fb.Rows {
			c := uint32(r.GetIntAt(1))
			if c == 0 {
				c = uint32(r.GetIntAt(0))
			}
			if c > 0 {
				cat.SlotBetCoins = append(cat.SlotBetCoins, SlotBetCoinDef{
					CoinTier: uint8(len(cat.SlotBetCoins)),
					Coins:    c,
					Mult:     uint32(len(cat.SlotBetCoins) + 1),
				})
			}
		}
	}
	if fa := ds.loadSHNWithFallback(dataDir, "GBSMAll"); fa != nil && len(fa.Rows) > 0 {
		if v := uint32(fa.Rows[0].GetIntAt(1)); v > 0 {
			cat.SlotAllRate = v
		}
	}
	if fl := ds.loadSHNWithFallback(dataDir, "GBSMLine"); fl != nil && len(fl.Rows) > 0 {
		if v := uint32(fl.Rows[0].GetIntAt(1)); v > 0 {
			cat.SlotLineRate = v
		}
	}
	if fc := ds.loadSHNWithFallback(dataDir, "GBSMCenter"); fc != nil && len(fc.Rows) > 0 {
		if v := uint32(fc.Rows[0].GetIntAt(1)); v > 0 {
			cat.SlotCenterRate = v
		}
	}
	if fj := ds.loadSHNWithFallback(dataDir, "GBSMJPRate"); fj != nil && len(fj.Rows) > 0 {
		if v := uint32(fj.Rows[0].GetIntAt(0)); v > 0 {
			cat.JackpotBasePct = v
		}
	}
	return cat
}

// BuildPhase15Catalogs inicializa todos los catálogos de la Fase 15 en DataStore.
func (ds *DataStore) BuildPhase15Catalogs(dataDir string) {
	if ds == nil {
		return
	}
	ds.QuestCatalog = BuildQuestCatalog(ds, dataDir)
	ds.KQCatalog = BuildKQCatalog(ds, dataDir)
	ds.GTCatalog = BuildGTCatalog(ds, dataDir)
	ds.MiniHouseCatalog = BuildMiniHouseCatalog(ds, dataDir)
	ds.PetCatalog = BuildPetCatalog(ds, dataDir)
	ds.WeaponTitleCatalog = BuildWeaponTitleCatalog(ds, dataDir)
	ds.CasinoCatalog = BuildCasinoCatalog(ds, dataDir)
	log.Printf("[Zone][Data] Phase 15 Catalogs loaded: %d quests, %d dialogs, %d scripts, %d KQs, %d MiniHouses, %d WeaponTitles.",
		len(ds.QuestCatalog.ByID), len(ds.QuestCatalog.Dialogs), len(ds.QuestCatalog.Scripts),
		len(ds.KQCatalog.ByID), len(ds.MiniHouseCatalog.ByItemID), len(ds.WeaponTitleCatalog.Rows))
}
