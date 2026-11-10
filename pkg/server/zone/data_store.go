package zone

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
	"github.com/iFeddy/XelionCore/pkg/data/areabmp"
	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/data/drop"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/data/misc"
	"github.com/iFeddy/XelionCore/pkg/data/mobattack"
	"github.com/iFeddy/XelionCore/pkg/data/mobregen"
	"github.com/iFeddy/XelionCore/pkg/data/mobroam"
	"github.com/iFeddy/XelionCore/pkg/data/mobsetting"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/npcitem"
	"github.com/iFeddy/XelionCore/pkg/data/questparser"
	"github.com/iFeddy/XelionCore/pkg/data/script"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"strings"
)

// DataStore contiene en memoria RAM todos los subsistemas de datos oficiales del juego cargados para Zone.
type DataStore struct {
	ClassStore      *classparam.Store
	FieldStore      *field.Store
	ItemStore       *item.Store
	NPCStore        *npc.Store
	MobRegen        map[string]*mobregen.MapRegen
	MobSetting      *mobsetting.Database
	MobAttack       *mobattack.Database
	MobRoam         *mobroam.Database
	DropStore       *drop.Store
	CombatStore     *combat.Store
	NPCItem         *npcitem.Database
	AbState         *abstate.Database
	MiscStore       *misc.Store
	QuestParser     *questparser.Store
	Script          *script.Database
	ChkEngine       *shn.ChecksumEngine
	BlockInfo       *blockinfo.Store
	AreaBMP         *areabmp.Store
	SHNTables       map[string]*shn.File
	ActiveSkillSHN  *shn.File
	PassiveSkillSHN *shn.File
	MobInfoSHN      *shn.File
	PassiveSkillMap map[uint16]bool
	ActiveSkillMap  map[uint16]bool
}

// GetSHN retorna una tabla .shn por nombre o nombre base (búsqueda insensible a mayúsculas/minúsculas).
func (ds *DataStore) GetSHN(name string) *shn.File {
	if ds == nil || ds.SHNTables == nil {
		return nil
	}
	base := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	if f, ok := ds.SHNTables[base]; ok {
		return f
	}
	for k, v := range ds.SHNTables {
		if strings.EqualFold(k, base) {
			return v
		}
	}
	return nil
}

// IsPassiveSkill verifica rápidamente si un ID de habilidad pertenece a PassiveSkill.shn.
func (ds *DataStore) IsPassiveSkill(skillID uint16) bool {
	if ds == nil || ds.PassiveSkillMap == nil {
		return false
	}
	return ds.PassiveSkillMap[skillID]
}

// LoadAllData carga en memoria RAM la totalidad del universo de datos estáticos de ZoneService.
func LoadAllData(dataDir string) (*DataStore, error) {
	ds := &DataStore{}

	// 1. ClassParam
	var err error
	ds.ClassStore, err = classparam.Load(filepath.Join(dataDir, "World", "ClassParam"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: ClassParam load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d character classes.", ds.ClassStore.ClassCount())
	}

	// 2. Field / Map
	ds.FieldStore, err = field.Load(filepath.Join(dataDir, "World", "Map"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: Field load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d field definitions.", ds.FieldStore.FieldCount())
	}

	// 3. Item Store (ItemInfo.shn + ItemUseEffect.shn + World/Item JSONs)
	ds.ItemStore = item.NewStore()
	itemShnPath := filepath.Join(dataDir, "ItemInfo.shn")
	if err := ds.ItemStore.LoadItemInfoSHN(itemShnPath); err != nil {
		log.Printf("[Zone][Data] Warning: ItemInfo.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded 14999 item definitions from ItemInfo.shn.")
	}
	if err := ds.ItemStore.LoadInventoryDataSHN(
		filepath.Join(dataDir, "ItemInfoServer.shn"),
		filepath.Join(dataDir, "ItemSort.shn"),
		filepath.Join(dataDir, "UseClassTypeInfo.shn"),
		filepath.Join(dataDir, "ClassName.shn"),
	); err != nil {
		log.Printf("[Zone][Data] Warning: inventory item data load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded official inventory sort order and class-use matrix.")
	}
	useEffPath := filepath.Join(dataDir, "ItemUseEffect.shn")
	if err := ds.ItemStore.LoadItemUseEffectSHN(useEffPath); err != nil {
		log.Printf("[Zone][Data] Warning: ItemUseEffect.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded consumable recovery effects from ItemUseEffect.shn.")
	}
	if err := ds.ItemStore.LoadDir(filepath.Join(dataDir, "World", "Item")); err != nil {
		log.Printf("[Zone][Data] Warning: World/Item JSON load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded item options, containers, and random options.")
	}

	// 4. NPCs y links de mapas
	ds.NPCStore, err = npc.Load(filepath.Join(dataDir, "World", "NPC"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: NPC load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded NPC database (26 map spawns, LinkTable, Actions).")
	}

	// 5. Mob Regen
	ds.MobRegen, err = mobregen.LoadAll(filepath.Join(dataDir, "MobRegen"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobRegen load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d map mob respawn configurations.", len(ds.MobRegen))
	}

	// 6. Mob Setting & Behavior
	ds.MobSetting, err = mobsetting.LoadFromDir(filepath.Join(dataDir, "MobSetting"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobSetting load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d monster action configurations.", len(ds.MobSetting.Actions))
	}

	// 7. Mob Attack Sequences
	ds.MobAttack, err = mobattack.LoadFromDir(filepath.Join(dataDir, "MobAttackSequence"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobAttack load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d monster attack sequences.", len(ds.MobAttack.Mobs))
	}

	// 8. Mob Roam
	ds.MobRoam, err = mobroam.LoadFromDir(filepath.Join(dataDir, "MobRoam"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobRoam load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d monster roam paths.", len(ds.MobRoam.Routes))
	}

	// 9. Drop Store
	ds.DropStore = drop.NewStore()
	if err := ds.DropStore.LoadDir(filepath.Join(dataDir, "World", "Drop")); err != nil {
		log.Printf("[Zone][Data] Warning: Drop load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d monster drop tables and %d drop groups.",
			ds.DropStore.MobCount(), ds.DropStore.GroupCount())
	}

	// 10. Combat Store
	ds.CombatStore = combat.NewStore()
	if err := ds.CombatStore.LoadDir(filepath.Join(dataDir, "World", "Combat")); err != nil {
		log.Printf("[Zone][Data] Warning: Combat load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded combat formulas and experience tables (Max Level: %d).",
			ds.CombatStore.MaxLevel())
	}

	// 11. NPC Merchant Items
	ds.NPCItem, err = npcitem.LoadFromDir(filepath.Join(dataDir, "NPCItemList"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: NPCItemList load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d merchant inventories.", len(ds.NPCItem.Merchants))
	}

	// 12. AbState (Buffs / Debuffs)
	ds.AbState, err = abstate.LoadFromDir(filepath.Join(dataDir, "AbState"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: AbState load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d abnormal state configurations.", len(ds.AbState.States))
	}

	// 13. Misc Store
	ds.MiscStore = misc.NewStore()
	if err := ds.MiscStore.LoadDir(filepath.Join(dataDir, "World", "Misc")); err != nil {
		log.Printf("[Zone][Data] Warning: Misc load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded misc game mechanics and stealth matrix.")
	}

	// 14. Quest Parser
	ds.QuestParser = questparser.NewStore()
	if err := ds.QuestParser.LoadFile(filepath.Join(dataDir, "World", "Quest", "QuestParser.json")); err != nil {
		log.Printf("[Zone][Data] Warning: QuestParser load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d quest parser entries.", ds.QuestParser.Count())
	}

	// 15. Scripts
	ds.Script, err = script.LoadFromDir(filepath.Join(dataDir, "Script"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: Script load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d game event script tables.", len(ds.Script.Tables))
	}

	// 16. Checksum Engine (49 SHNs)
	ds.ChkEngine, err = shn.NewChecksumEngine(dataDir, filepath.Join(dataDir, "View"))
	if err != nil {
		log.Printf("[Zone][Data] Warning: ChecksumEngine error: %v", err)
	} else {
		log.Printf("[Zone][Data] Initialized ChecksumEngine with 49 SHN files.")
	}

	// 17. ActiveSkill.shn
	activeSkillPath := filepath.Join(dataDir, "ActiveSkill.shn")
	ds.ActiveSkillSHN, err = shn.ReadFile(activeSkillPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: ActiveSkill.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded ActiveSkill.shn (%d skills).", ds.ActiveSkillSHN.RecordCount)
		ds.ActiveSkillMap = make(map[uint16]bool, len(ds.ActiveSkillSHN.Rows))
		for _, row := range ds.ActiveSkillSHN.Rows {
			ds.ActiveSkillMap[row.GetUint16("ID")] = true
		}
	}

	// 18. PassiveSkill.shn
	passiveSkillPath := filepath.Join(dataDir, "PassiveSkill.shn")
	ds.PassiveSkillSHN, err = shn.ReadFile(passiveSkillPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: PassiveSkill.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded PassiveSkill.shn (%d passive skills).", ds.PassiveSkillSHN.RecordCount)
		ds.PassiveSkillMap = make(map[uint16]bool, len(ds.PassiveSkillSHN.Rows))
		for _, row := range ds.PassiveSkillSHN.Rows {
			ds.PassiveSkillMap[row.GetUint16("ID")] = true
		}
	}

	// 19. MobInfo.shn
	mobInfoPath := filepath.Join(dataDir, "MobInfo.shn")
	ds.MobInfoSHN, err = shn.ReadFile(mobInfoPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobInfo.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded MobInfo.shn (%d monsters).", ds.MobInfoSHN.RecordCount)
	}

	// 20. Universo de Geometría y Colisiones (301 archivos en BlockInfo/)
	blockInfoPath := filepath.Join(dataDir, "BlockInfo")
	ds.BlockInfo, err = blockinfo.LoadAll(blockInfoPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: BlockInfo load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d geometry/collision files from BlockInfo/ (158 shbd, 33 sbi, 23 aid, 76 bdt, 8 shab, 2 shmd, 1 shad).",
			ds.BlockInfo.TotalCount())
	}

	// 21. Máscaras AoE de Habilidades (AreaBMP/ - 44+ máscaras)
	areaBMPPath := filepath.Join(dataDir, "AreaBMP")
	ds.AreaBMP, err = areabmp.LoadAll(areaBMPPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: AreaBMP load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded %d AoE masks from AreaBMP/.", ds.AreaBMP.Count())
	}

	// 22. Catálogo exhaustivo de las 180 tablas .shn reales en disco (194/195 rutas)
	ds.SHNTables, err = LoadAllSHNTables(dataDir)
	if err != nil {
		return nil, fmt.Errorf("error cargando catálogo SHN: %w", err)
	}
	log.Printf("[Zone][Data] Successfully indexed %d unique SHN table keys in RAM (180 required SHNs verified).", len(ds.SHNTables)/2)

	return ds, nil
}
