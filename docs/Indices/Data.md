# Ãndice: pkg/data (cargadores de datos)

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. Datos JSON en data/.

## pkg/data/abstate/abstate.go
- L12: `type LevelConfig struct`
- L29: `type StateConfig struct`
- L39: `type Database struct`
- L45: `func LoadFromDir(dirPath string) (*Database, error)`
- L90: `func LoadFromFile(path string) (*Database, error)`
- L132: `func (db *Database) GetByName(name string) (*StateConfig, bool)`
- L141: `func (db *Database) GetByID(id uint16) (*StateConfig, bool)`
- L147: `func (s *StateConfig) GetLevelConfig(level int) (*LevelConfig, error)`

## pkg/data/abstate/abstate_test.go
- L10: `func TestLoadFromFile(t *testing.T)`

## pkg/data/classparam/classparam.go
- L13: `type LevelParam struct`
- L51: `type Store struct`
- L57: `func NewStore() *Store`
- L64: `func Load(dir string) (*Store, error)`
- L73: `func (s *Store) LoadDir(dir string) error`
- L108: `func (s *Store) GetParam(className string, level int) (*LevelParam, bool)`
- L120: `func (s *Store) GetClass(className string) ([]LevelParam, bool)`
- L134: `func (s *Store) Classes() []string`
- L146: `func (s *Store) ClassCount() int`

## pkg/data/classparam/classparam_test.go
- L10: `func TestClassParamStore_LoadAll(t *testing.T)`

## pkg/data/combat/combat.go
- L12: `type CommonConfig struct`
- L31: `type ExpLevelEntry struct`
- L39: `type FreeStatTable struct`
- L68: `type AngleRate struct`
- L74: `type DamageByAngle struct`
- L80: `type SoulDemandRate struct`
- L86: `type PartyMemberBonus struct`
- L92: `type LevelDiffBonus struct`
- L98: `type StandardDeviationHandicap struct`
- L104: `type ExpRecalculation struct`
- L111: `type Store struct`
- L122: `func NewStore() *Store`
- L127: `func (s *Store) LoadDir(dir string) error`
- L189: `func (s *Store) Common() CommonConfig`
- L196: `func (s *Store) GetNextExp(level int) (int, bool)`
- L206: `func (s *Store) MaxLevel() int`
- L213: `func (s *Store) GetPartyExpBonus(memberCount int) int`
- L228: `func (s *Store) GetAngleDamageRate(isMonster bool, angle int) int`

## pkg/data/combat/combat_test.go
- L10: `func TestCombatStore_LoadAll(t *testing.T)`

## pkg/data/drop/drop.go
- L13: `type DropEntry struct`
- L23: `type MoneyDrop struct`
- L30: `type TradeBoxDrop struct`
- L36: `type MobDropTable struct`
- L48: `type ItemGroupEntry struct`
- L56: `type Store struct`
- L63: `func NewStore() *Store`
- L71: `func (s *Store) LoadDir(dir string) error`
- L116: `func (s *Store) GetMobDrop(mobID string) (*MobDropTable, bool)`
- L127: `func (s *Store) GetGroup(groupName string) ([]ItemGroupEntry, bool)`
- L140: `func (s *Store) MobCount() int`
- L147: `func (s *Store) GroupCount() int`

## pkg/data/drop/drop_test.go
- L10: `func TestDropStore_LoadAll(t *testing.T)`

## pkg/data/field/field.go
- L13: `type CombatFlags struct`
- L27: `type TeleportFlags struct`
- L34: `type Coord struct`
- L40: `type CoordDir struct`
- L47: `type RespawnConfig struct`
- L54: `type PermissionFlags struct`
- L70: `type ExpLossConfig struct`
- L76: `type FieldEntry struct`
- L99: `type InstanceDungeonEntry struct`
- L117: `type TournamentDoor struct`
- L125: `type GuildTournamentEntry struct`
- L136: `type TutorialMapEntry struct`
- L143: `type Store struct`
- L154: `func NewStore() *Store`
- L163: `func Load(dir string) (*Store, error)`
- L172: `func (s *Store) LoadDir(dir string) error`
- L243: `func (s *Store) GetFieldBySerial(serial int) (*FieldEntry, bool)`
- L251: `func (s *Store) GetField(mapID string) (*FieldEntry, bool)`
- L262: `func (s *Store) GetAllFieldsByMapID(mapID string) []*FieldEntry`
- L272: `func (s *Store) GetDungeonsByZone(zoneNumber int) []InstanceDungeonEntry`
- L285: `func (s *Store) FieldCount() int`

## pkg/data/field/field_test.go
- L10: `func TestFieldStore_LoadAll(t *testing.T)`

## pkg/data/item/item.go
- L16: `type ItemStatOption struct`
- L23: `func CalcAttrSizeForClass(itemClass uint8, numOpts int) int`
- L66: `func GetAttrSizeForClass(itemClass uint8) int`
- L72: `func BuildItemAttrBuffer(itemClass uint8, amount uint32, opts []ItemStatOption, innerCount int) []byte`
- L119: `func writeStatOptions(dst []byte, opts []ItemStatOption)`
- L130: `type ItemOptionEntry struct`
- L137: `type ItemUseFunction struct`
- L143: `type StatRange struct`
- L149: `type RandomOptionEntry struct`
- L161: `type TreasureRewardItem struct`
- L170: `type TreasureContainer struct`
- L179: `type ItemInfo struct`
- L207: `type ItemUseEffect struct`
- L220: `type Store struct`
- L235: `func NewStore() *Store`
- L250: `func (s *Store) LoadDir(dir string) error`
- L318: `func (s *Store) GetUseFunction(itemID string) (*ItemUseFunction, bool)`
- L329: `func (s *Store) GetRandomOption(dropItemIndex string) (*RandomOptionEntry, bool)`
- L340: `func (s *Store) GetContainer(containerID string) (*TreasureContainer, bool)`
- L351: `func (s *Store) Options() []ItemOptionEntry`
- L360: `func (s *Store) LoadItemInfoSHN(filePath string) error`
- L417: `func (s *Store) LoadInventoryDataSHN(itemInfoServerPath, itemSortPath, useClassTypeInfoPath, classNamePath string) error`
- L489: `func (s *Store) CanUseItemByClassID(itemID uint16, classID byte) (bool, bool)`
- L512: `func (s *Store) LoadItemUseEffectSHN(filePath string) error`
- L545: `func (s *Store) GetItemInfo(itemID uint16) (*ItemInfo, bool)`
- L556: `func (s *Store) GetItemInfoByName(inxName string) (*ItemInfo, bool)`
- L567: `func (s *Store) GetItemUseEffect(inxName string) (*ItemUseEffect, bool)`
- L578: `func (s *Store) GetItemClass(itemID uint16) (uint8, bool)`
- L589: `func (s *Store) GetItemAttrSize(itemID uint16) int`
- L598: `func (s *Store) BuildItemAttr(itemID uint16, amount uint32, opts []ItemStatOption, innerCount int) []byte`

## pkg/data/item/item_test.go
- L10: `func TestItemStore_LoadAll(t *testing.T)`
- L80: `func TestBuildItemAttrBuffer_DynamicSizes(t *testing.T)`

## pkg/data/misc/misc.go
- L13: `type RecallCoordEntry struct`
- L22: `type StealthMatrixConfig struct`
- L28: `type SkillBreedMob struct`
- L50: `type PolymorphEntry struct`
- L59: `type NeutralizationEntry struct`
- L69: `type AbStateResetEntry struct`
- L81: `type ExpandSkillEntry struct`
- L91: `type SkillMechanicsConfig struct`
- L100: `type MobChatEntry struct`
- L108: `type MobChatConfig struct`
- L119: `type Store struct`
- L128: `func NewStore() *Store`
- L135: `func (s *Store) LoadDir(dir string) error`
- L192: `func (s *Store) GetRecallCoord(itemIndex string) (*RecallCoordEntry, bool)`
- L204: `func (s *Store) CanSee(viewerLayer, targetLayer string) bool`
- L214: `func (s *Store) CanAttack(attackerLayer, targetLayer string) bool`
- L224: `func (s *Store) Mechanics() SkillMechanicsConfig`
- L231: `func (s *Store) MobChats() MobChatConfig`

## pkg/data/misc/misc_test.go
- L10: `func TestMiscStore_LoadAll(t *testing.T)`

## pkg/data/mobattack/mobattack.go
- L12: `type AttackOrder struct`
- L18: `type SkillChange struct`
- L26: `type MobAI struct`
- L33: `type Database struct`
- L38: `func LoadFromDir(dirPath string) (*Database, error)`
- L75: `func LoadFromFile(path string) (*Database, error)`
- L107: `func (db *Database) GetMobAI(name string) (*MobAI, bool)`

## pkg/data/mobattack/mobattack_test.go
- L10: `func TestLoadFromFile(t *testing.T)`

## pkg/data/mobregen/mobregen.go
- L12: `type MobSpawnEntry struct`
- L25: `type Group struct`
- L37: `type MapRegen struct`
- L44: `func LoadMapRegen(baseDir, mapName string) (*MapRegen, error)`
- L83: `func LoadAll(baseDir string) (map[string]*MapRegen, error)`
- L114: `func (m *MapRegen) GetActiveGroups() map[string]Group`

## pkg/data/mobregen/mobregen_test.go
- L10: `func TestLoadMapRegen_RouVal01(t *testing.T)`
- L41: `func TestLoadMapRegen_RouN_CommentedRecords(t *testing.T)`
- L75: `func TestLoadMapRegen_Subdirectory_KQ(t *testing.T)`
- L88: `func TestLoadMapRegen_SpaceGroupNames_BerKal01(t *testing.T)`

## pkg/data/mobroam/mobroam.go
- L12: `type Waypoint struct`
- L19: `type Route struct`
- L27: `type Database struct`
- L32: `func LoadFromDir(dirPath string) (*Database, error)`
- L69: `func LoadFromFile(path string) (*Database, error)`
- L102: `func (db *Database) GetRoute(name string) (*Route, bool)`

## pkg/data/mobroam/mobroam_test.go
- L10: `func TestLoadFromFile(t *testing.T)`

## pkg/data/mobsetting/mobsetting.go
- L12: `type ActionTrigger struct`
- L18: `type EntitySetting struct`
- L24: `type Database struct`
- L29: `func LoadFromDir(dirPath string) (*Database, error)`
- L62: `func LoadFromFile(path string) (*Database, error)`
- L95: `func (db *Database) GetTriggers(name string) []ActionTrigger`
- L100: `func (db *Database) HasTriggers(name string) bool`
- L106: `func (db *Database) TotalEntities() int`
- L111: `func (db *Database) TotalTriggers() int`

## pkg/data/mobsetting/mobsetting_test.go
- L8: `func TestLoadMobSetting(t *testing.T)`

## pkg/data/npc/npc.go
- L13: `type ShineNPC struct`
- L26: `type LinkTableEntry struct`
- L38: `type NPCCondition struct`
- L51: `type NPCAction struct`
- L60: `type NPCActionData struct`
- L66: `type Store struct`
- L77: `func Load(npcDir string) (*Store, error)`
- L149: `func (s *Store) GetNPCs(mapName string, onlyEnabled bool) []ShineNPC`
- L174: `func (s *Store) GetAllNPCs(onlyEnabled bool) []ShineNPC`
- L190: `func (s *Store) MapCount() int`
- L197: `func (s *Store) GetLinks(onlyEnabled bool) []LinkTableEntry`
- L217: `func (s *Store) GetLinkByArg(arg string) (LinkTableEntry, bool)`
- L225: `func (s *Store) GetConditions() []NPCCondition`
- L234: `func (s *Store) GetConditionByID(id int) (NPCCondition, bool)`
- L242: `func (s *Store) GetActions() []NPCAction`
- L251: `func (s *Store) GetActionByID(id int) (NPCAction, bool)`

## pkg/data/npc/npc_test.go
- L8: `func TestLoadNPC(t *testing.T)`

## pkg/data/npcitem/npcitem.go
- L12: `type ShopTab struct`
- L20: `type ShopItem struct`
- L30: `type Merchant struct`
- L36: `type Database struct`
- L41: `func LoadFromDir(dirPath string) (*Database, error)`
- L68: `func LoadMerchant(filePath string) (*Merchant, error)`
- L83: `func (db *Database) GetMerchant(name string) (*Merchant, bool)`
- L89: `func (db *Database) TotalMerchants() int`
- L94: `func (db *Database) TotalTabs() int`
- L103: `func (db *Database) TotalValidItems() int`
- L112: `func (m *Merchant) GetTab(index int) (*ShopTab, bool)`
- L120: `func (m *Merchant) GetItems(tabIndex int) []ShopItem`
- L146: `func (m *Merchant) GetAllItems() []ShopItem`
- L155: `func (m *Merchant) HasItem(itemInCode string) bool`
- L169: `func (m *Merchant) GetItemAt(tabIndex, row, col int) (string, bool)`

## pkg/data/npcitem/npcitem_test.go
- L8: `func TestLoadNPCItemList(t *testing.T)`

## pkg/data/questparser/questparser.go
- L12: `type Store struct`
- L19: `func NewStore() *Store`
- L27: `func (s *Store) LoadFile(filePath string) error`
- L55: `func (s *Store) GetID(className string) (int, bool)`
- L63: `func (s *Store) GetName(classID int) (string, bool)`
- L71: `func (s *Store) Count() int`

## pkg/data/questparser/questparser_test.go
- L10: `func TestQuestParserStore_Load(t *testing.T)`

## pkg/data/script/script.go
- L12: `type ScriptTable map[string]string`
- L16: `type Database struct`
- L21: `func LoadFromDir(dirPath string) (*Database, error)`
- L55: `func (db *Database) Get(scriptName, key string) (string, bool)`
- L65: `func (db *Database) TotalFiles() int`
- L70: `func (db *Database) TotalEntries() int`

## pkg/data/script/script_test.go
- L8: `func TestLoadScriptDatabase(t *testing.T)`

## pkg/data/shn/checksum.go
- L69: `type Name8 [32]byte`
- L72: `func (n Name8) String() string`
- L77: `type ChecksumEngine struct`
- L83: `func findFileCaseInsensitive(filename string, dirs []string) (string, error)`
- L101: `func NewChecksumEngine(dirs ...string) (*ChecksumEngine, error)`
- L152: `func (e *ChecksumEngine) VerifyDataFiles(clientChecksums [49]Name8) (bool, int)`
- L163: `func (e *ChecksumEngine) GetChecksums() [49]Name8`
- L168: `func (e *ChecksumEngine) GetRawHash(index int) string`
- L176: `func (e *ChecksumEngine) GetFileName(index int) string`

## pkg/data/shn/checksum_test.go
- L10: `func TestChecksumEngine_Verification(t *testing.T)`

## pkg/data/shn/shn.go
- L32: `func Decrypt(data []byte) []byte`
- L63: `type Column struct`
- L70: `type Row struct`
- L75: `func (r *Row) GetByte(colName string) byte`
- L84: `func (r *Row) GetUint16(colName string) uint16`
- L93: `func (r *Row) GetInt16(colName string) int16`
- L98: `func (r *Row) GetUint32(colName string) uint32`
- L107: `func (r *Row) GetInt(colName string) int`
- L125: `func (r *Row) GetInt64(colName string) int64`
- L134: `func (r *Row) GetString(colName string) string`
- L147: `func (r *Row) RawBytes(colName string) []byte`
- L152: `type File struct`
- L162: `func ReadFile(filePath string) (*File, error)`
- L171: `func Parse(data []byte) (*File, error)`
- L255: `func (f *File) FindRows(colName string, targetVal string) []*Row`
- L266: `func (f *File) FindRow(colName string, targetVal string) *Row`

## pkg/data/shn/shn_test.go
- L10: `func TestSHNParser_ClassName(t *testing.T)`
- L49: `func TestSHNParser_BadNameFilter(t *testing.T)`
- L65: `func TestSHNParser_MapInfo(t *testing.T)`
- L84: `func TestSHNParser_SingleData(t *testing.T)`
- L106: `func TestSHNParser_CharacterAppearance(t *testing.T)`

## pkg/script/lua/engine.go
- L16: `type Engine struct`
- L23: `type Config struct`
- L29: `func NewEngine(cfg Config) (*Engine, error)`
- L77: `func (e *Engine) Close()`
- L87: `func (e *Engine) DoFile(filename string) error`
- L101: `func (e *Engine) CallFunction(name string, args ...lua.LValue) ([]lua.LValue, error)`
- L130: `func (e *Engine) GetGlobal(name string) lua.LValue`
- L137: `func (e *Engine) GetGlobalTableField(tableName, fieldName string) (lua.LValue, error)`
- L151: `func (e *Engine) registerOfficialAPIs(debug bool)`

## pkg/script/lua/engine_test.go
- L11: `func TestLuaEngine_InitAndCommon(t *testing.T)`
- L65: `func TestMobBehavior_LoadAllLua(t *testing.T)`
- L97: `func TestMigratedScenarios_LoadAllLua(t *testing.T)`

