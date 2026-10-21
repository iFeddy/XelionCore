package zone

import (
	"log"
	"path/filepath"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
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
	ActiveSkillSHN  *shn.File
	PassiveSkillSHN *shn.File
	MobInfoSHN      *shn.File
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
	}

	// 18. PassiveSkill.shn
	passiveSkillPath := filepath.Join(dataDir, "PassiveSkill.shn")
	ds.PassiveSkillSHN, err = shn.ReadFile(passiveSkillPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: PassiveSkill.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded PassiveSkill.shn (%d passive skills).", ds.PassiveSkillSHN.RecordCount)
	}

	// 19. MobInfo.shn
	mobInfoPath := filepath.Join(dataDir, "MobInfo.shn")
	ds.MobInfoSHN, err = shn.ReadFile(mobInfoPath)
	if err != nil {
		log.Printf("[Zone][Data] Warning: MobInfo.shn load error: %v", err)
	} else {
		log.Printf("[Zone][Data] Loaded MobInfo.shn (%d monsters).", ds.MobInfoSHN.RecordCount)
	}

	return ds, nil
}
