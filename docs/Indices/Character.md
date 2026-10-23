# Ãndice: WorldService / Character

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. Equivale a WorldManager+Character.

## pkg/server/world/e2e_world_test.go
- L22: `func readServerPacket(r io.Reader) (*protocol.Packet, error)`
- L47: `func sendClientPacket(w io.Writer, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error`
- L70: `func TestE2EWorldServiceFlow(t *testing.T)`

## pkg/server/world/friend_manager.go
- L14: `type FriendManager struct`
- L22: `func NewFriendManager(repo database.CharacterRepository, pm *PartyManager) *FriendManager`
- L31: `func (fm *FriendManager) SendFriendList(client *network.Client)`
- L69: `func (fm *FriendManager) RequestFriend(requester *network.Client, targetName string)`
- L106: `func (fm *FriendManager) AcceptFriend(target *network.Client)`
- L142: `func (fm *FriendManager) DeleteFriend(client *network.Client, friendName string)`
- L170: `func (fm *FriendManager) SendFriendPoints(client *network.Client)`
- L180: `func (fm *FriendManager) BroadcastLogin(client *network.Client, charName, mapName string)`
- L198: `func (fm *FriendManager) BroadcastLogout(charID int, charName string)`

## pkg/server/world/friend_test.go
- L12: `func TestFriendManager_FullFlow(t *testing.T)`

## pkg/server/world/guild_manager.go
- L16: `type GuildOnlinePlayer struct`
- L23: `type GuildPendingInvite struct`
- L32: `type GuildManager struct`
- L41: `func NewGuildManager(repo database.CharacterRepository) *GuildManager`
- L51: `func (gm *GuildManager) RegisterOnline(client *network.Client, charID int, charName string)`
- L70: `func (gm *GuildManager) UnregisterOnline(charID int)`
- L89: `func (gm *GuildManager) BroadcastToGuild(guildID int, packet *protocol.PacketWriter, excludeCharID int)`
- L109: `func (gm *GuildManager) BroadcastToAcademy(guildID int, packet *protocol.PacketWriter, excludeCharID int)`
- L133: `func (gm *GuildManager) HandleGuildList(client *network.Client, partMask byte)`
- L159: `func (gm *GuildManager) HandleGuildMake(client *network.Client, name, password string, guildType byte)`
- L187: `func (gm *GuildManager) HandleMyGuildInfo(client *network.Client)`
- L224: `func (gm *GuildManager) HandleMemberList(client *network.Client)`
- L267: `func (gm *GuildManager) HandleInvite(client *network.Client, targetName string)`
- L306: `func (gm *GuildManager) HandleJoinAck(client *network.Client, guildName string, joinAck byte)`
- L358: `func (gm *GuildManager) HandleLeave(client *network.Client)`
- L383: `func (gm *GuildManager) HandleVanish(client *network.Client, targetName string)`
- L408: `func (gm *GuildManager) HandleGrade(client *network.Client, targetName string, newGrade byte)`
- L433: `func (gm *GuildManager) HandleNotify(client *network.Client, notify string)`
- L453: `func (gm *GuildManager) HandleDismiss(client *network.Client)`
- L473: `func (gm *GuildManager) HandleChat(client *network.Client, msg string)`
- L491: `func (gm *GuildManager) HandleAcademyList(client *network.Client, partMask byte)`
- L515: `func (gm *GuildManager) HandleAcademyMyInfo(client *network.Client)`
- L547: `func (gm *GuildManager) HandleAcademyMemberList(client *network.Client)`
- L586: `func (gm *GuildManager) HandleAcademyJoin(client *network.Client, guildName string)`
- L639: `func (gm *GuildManager) HandleAcademyLeave(client *network.Client)`
- L659: `func (gm *GuildManager) HandleAcademyChat(client *network.Client, msg string)`

## pkg/server/world/guild_test.go
- L14: `func TestGuildManager_FullFlow(t *testing.T)`
- L199: `func TestGuildProtocol_PacketStructures(t *testing.T)`

## pkg/server/world/handler.go
- L15: `type Handler struct`
- L28: `func NewHandler(`
- L64: `func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet)`
- L278: `func (h *Handler) handleGameTime(client *network.Client, pkt *protocol.Packet)`
- L289: `func (h *Handler) handleWorldClientKey(client *network.Client, pkt *protocol.Packet)`
- L335: `func (h *Handler) sendCharacterList(client *network.Client)`
- L363: `func (h *Handler) handleCreateCharacter(client *network.Client, pkt *protocol.Packet)`
- L549: `func (h *Handler) handleDeleteCharacter(client *network.Client, pkt *protocol.Packet)`
- L572: `func (h *Handler) handleBackToWorldSelect(client *network.Client, pkt *protocol.Packet)`
- L578: `func (h *Handler) handleCharSelect(client *network.Client, pkt *protocol.Packet)`
- L648: `func (h *Handler) sendCharOptionDataImprove(client *network.Client, charNo int)`
- L690: `func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError)`
- L696: `func (h *Handler) sendCharCreationError(client *network.Client, errCode protocol.CreateCharError)`
- L701: `func (h *Handler) getActiveCharNo(client *network.Client) int`
- L712: `func (h *Handler) handleGetQuickBar(client *network.Client, pkt *protocol.Packet)`
- L724: `func (h *Handler) handleGetQuickBarState(client *network.Client, pkt *protocol.Packet)`
- L736: `func (h *Handler) handleGetGameSettings(client *network.Client, pkt *protocol.Packet)`
- L748: `func (h *Handler) handleGetClientSettings(client *network.Client, pkt *protocol.Packet)`
- L760: `func (h *Handler) handleGetShortCuts(client *network.Client, pkt *protocol.Packet)`
- L772: `func (h *Handler) handleSaveSettings(client *network.Client, pkt *protocol.Packet)`
- L807: `func (h *Handler) handlePrisonGet(client *network.Client, pkt *protocol.Packet)`
- L813: `func (h *Handler) handleGotIngame(client *network.Client, pkt *protocol.Packet)`

## pkg/server/world/handler_test.go
- L14: `func TestWriteBasicCharInfo(t *testing.T)`
- L40: `func TestCharacterRepositoryOperations(t *testing.T)`
- L98: `func readPacketFromConn(r io.Reader) (uint16, []byte, error)`
- L123: `func setupTestHandler() (*Handler, *database.MemoryCharacterRepository, *database.TransferManager)`
- L156: `func TestHandler_GameTime(t *testing.T)`
- L189: `func TestHandler_CreateCharacter_BadNameRejected(t *testing.T)`
- L236: `func TestHandler_CreateCharacter_Success(t *testing.T)`
- L293: `func TestHandler_CharSelect_ZoneHandover(t *testing.T)`
- L379: `func TestHandler_DeleteCharacter_Flow(t *testing.T)`
- L437: `func TestHandler_BackToWorldSelect(t *testing.T)`

## pkg/server/world/kq_manager.go
- L15: `type KQParticipant struct`
- L25: `type KQInstance struct`
- L41: `type KQManager struct`
- L50: `func NewKQManager(kqSHN *shn.File) *KQManager`
- L67: `func (km *KQManager) loadFromSHN(f *shn.File)`
- L118: `func (km *KQManager) loadDefaults()`
- L156: `func (km *KQManager) RegisterOnline(client *network.Client, charID int)`
- L163: `func (km *KQManager) UnregisterOnline(charID int)`
- L178: `func (km *KQManager) BroadcastToOnline(packet *protocol.PacketWriter)`
- L189: `func (km *KQManager) HandleList(client *network.Client)`
- L237: `func (km *KQManager) HandleStatus(client *network.Client, handle uint32)`
- L256: `func (km *KQManager) HandleJoin(client *network.Client, handle uint32, charLevel, charClass byte)`
- L323: `func (km *KQManager) HandleJoinCancel(client *network.Client, handle uint32)`
- L355: `func (km *KQManager) HandleJoinList(client *network.Client, handle uint32)`
- L379: `func (km *KQManager) HandleTeamSelect(client *network.Client, teamType byte)`

## pkg/server/world/kq_test.go
- L13: `func TestKQManager_FullFlow(t *testing.T)`
- L79: `func TestKQProtocol_PacketStructures(t *testing.T)`

## pkg/server/world/packet_helper.go
- L30: `func WriteBasicCharInfo(ch *database.Character, pw *protocol.PacketWriter)`
- L111: `func GetBeginnerWeapon(job protocol.Job) (slot byte, itemID uint16)`

## pkg/server/world/party_manager.go
- L13: `type PartyMember struct`
- L23: `type Party struct`
- L33: `type PartyFinderAd struct`
- L41: `type PartyManager struct`
- L52: `func NewPartyManager() *PartyManager`
- L64: `func (pm *PartyManager) RegisterOnline(client *network.Client, charID int, charName string, classID, level byte)`
- L79: `func (pm *PartyManager) UnregisterOnline(charID int)`
- L131: `func (pm *PartyManager) InviteMember(inviter *network.Client, targetName string)`
- L192: `func (pm *PartyManager) AcceptInvite(target *network.Client)`
- L263: `func (pm *PartyManager) RejectInvite(target *network.Client)`
- L283: `func (pm *PartyManager) LeaveParty(client *network.Client)`
- L328: `func (pm *PartyManager) KickMember(master *network.Client, targetName string)`
- L375: `func (pm *PartyManager) ChangeMaster(master *network.Client, newMasterName string)`
- L414: `func (pm *PartyManager) SetLooting(master *network.Client, lootType, itemGrade byte)`
- L440: `func (pm *PartyManager) AddFinderAd(client *network.Client, purpose string)`
- L475: `func (pm *PartyManager) DeleteFinderAd(client *network.Client)`
- L491: `func (pm *PartyManager) ListFinderAds(client *network.Client)`
- L513: `func (pm *PartyManager) broadcastMemberListLocked(party *Party)`
- L535: `func (pm *PartyManager) dissolvePartyLocked(party *Party)`

## pkg/server/world/party_test.go
- L11: `func TestPartyManager_FullFlow(t *testing.T)`

## pkg/server/world/server.go
- L12: `type Config struct`
- L22: `type Server struct`
- L29: `func NewServer(`
- L65: `func (s *Server) Start() error`
- L72: `func (s *Server) Stop()`

## pkg/server/world/static_data.go
- L15: `type StarterItem struct`
- L22: `type StarterShortcut struct`
- L29: `type MaxShape struct`
- L36: `type DefaultClassData struct`
- L56: `type DefaultCharacterDataStore struct`
- L61: `func LoadDefaultCharacterData(filePath string) (*DefaultCharacterDataStore, error)`
- L76: `func (s *DefaultCharacterDataStore) GetClassData(classID byte) (DefaultClassData, bool)`
- L83: `type BadNameFilter struct`
- L89: `func LoadBadNameFilter(filePath string) (*BadNameFilter, error)`
- L111: `func (f *BadNameFilter) IsBadName(name string) bool`
- L134: `type WorldDataStore struct`
- L181: `func LoadWorldDataStore(dataDir string) (*WorldDataStore, error)`
- L300: `func (s *WorldDataStore) ValidateAppearance(classID byte, hair, hairColor, face byte) bool`

## pkg/server/world/static_data_test.go
- L9: `func TestLoadWorldDataStore_AllFiles(t *testing.T)`

## pkg/database/account.go
- L6: `type TUser struct`
- L20: `type TUserAuth struct`
- L27: `type TUserBlock struct`
- L34: `type Account struct`
- L46: `type Report struct`
- L61: `type ReportSummary struct`
- L68: `type LogoutLog struct`
- L81: `type ChargedItem struct`
- L94: `func PackShineDateTime(t time.Time) uint32`
- L112: `func (a *Account) Sync()`

## pkg/database/character.go
- L24: `type TCharacter struct`
- L41: `type TCharacterShape struct`
- L50: `type TCharacterState struct`
- L66: `type TCharacterKeyMap struct`
- L73: `type TCharacterShortCut struct`
- L79: `type TCharacterGameOption struct`
- L85: `type TCharacterSingleOption struct`
- L91: `type Character struct`
- L127: `type Item struct`
- L138: `type Shortcut struct`
- L147: `type CharacterSkill struct`
- L154: `func (c *Character) Sync()`
- L217: `func (c *Character) GetEquip(slot byte) uint16`
- L228: `type CharacterOptions struct`
- L240: `type CharacterRepository interface`
- L307: `type GuildData struct`
- L328: `type GuildMemberData struct`
- L345: `type GuildAcademyData struct`
- L360: `type GuildAcademyMemberData struct`
- L375: `type GuildListItem struct`
- L386: `type FriendInfo struct`
- L395: `type charStore struct`
- L401: `type MemoryCharacterRepository struct`
- L417: `func NewMemoryCharacterRepositoryWithFile(filePath string) *MemoryCharacterRepository`
- L436: `func NewMemoryCharacterRepository() *MemoryCharacterRepository`
- L440: `func (r *MemoryCharacterRepository) load()`
- L463: `func (r *MemoryCharacterRepository) save()`
- L483: `func (r *MemoryCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error)`
- L498: `func (r *MemoryCharacterRepository) GetCharacterByName(name string) (*Character, error)`
- L512: `func (r *MemoryCharacterRepository) GetCharacterByID(id int) (*Character, error)`
- L525: `func (r *MemoryCharacterRepository) CreateCharacter(ch *Character) error`
- L553: `func (r *MemoryCharacterRepository) UpdateCharacter(ch *Character) error`
- L569: `func (r *MemoryCharacterRepository) DeleteCharacter(accountID int, slot byte) error`
- L585: `func (r *MemoryCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error)`
- L597: `func (r *MemoryCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error)`
- L609: `func (r *MemoryCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error)`
- L621: `func (r *MemoryCharacterRepository) nextItemKeyLocked() int64`
- L633: `func (r *MemoryCharacterRepository) SaveCharacterItem(item *Item) error`
- L651: `func (r *MemoryCharacterRepository) CreateItem(item *Item) error`
- L669: `func (r *MemoryCharacterRepository) DeleteItem(itemKey int64) error`
- L684: `func (r *MemoryCharacterRepository) SetItemStorage(itemKey int64, storageType byte, storageSlot int16) error`
- L700: `func (r *MemoryCharacterRepository) ExchangeItemStorage(itemKey1 int64, storageType1 byte, storageSlot1 int16, itemKey2 int64, storageType2 `
- L726: `func (r *MemoryCharacterRepository) ArrangeInventory(charNo int, itemKeys []int64) error`
- L761: `func (r *MemoryCharacterRepository) UpdateItemFlags(itemKey int64, flags int) error`
- L776: `func (r *MemoryCharacterRepository) GetItemsByStorage(charNo int, storageType byte) ([]*Item, error)`
- L792: `func (r *MemoryCharacterRepository) AddCharMoney(charNo int, amount int64) error`
- L805: `func (r *MemoryCharacterRepository) SubCharMoney(charNo int, amount int64) (bool, error)`
- L821: `func (r *MemoryCharacterRepository) SaveCharLocation(charNo int, mapName string, x, y int, rotation byte) error`
- L839: `func (r *MemoryCharacterRepository) SaveCharStats(charNo int, hp, sp, hps, sps int16) error`
- L857: `func (r *MemoryCharacterRepository) SaveCharLevelExp(charNo int, level int16, exp int64) error`
- L872: `func (r *MemoryCharacterRepository) SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error`
- L885: `func (r *MemoryCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error`
- L897: `func (r *MemoryCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error`
- L909: `func (r *MemoryCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error)`
- L920: `func (r *MemoryCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error`
- L947: `func (r *MemoryCharacterRepository) GetFriends(charNo int) ([]*FriendInfo, error)`
- L967: `func (r *MemoryCharacterRepository) AddFriend(charNo, friendCharNo int) error`
- L976: `func (r *MemoryCharacterRepository) RemoveFriend(charNo, friendCharNo int) error`
- L995: `func (r *MemoryCharacterRepository) GetFriendPoints(charNo int) (int, error)`
- L1001: `func (r *MemoryCharacterRepository) SetFriendPoints(charNo int, points int) error`
- L1008: `func (r *MemoryCharacterRepository) CreateGuild(masterCharID int, name, password string, guildType byte) (int, error)`
- L1075: `func (r *MemoryCharacterRepository) GetGuild(guildID int) (*GuildData, error)`
- L1085: `func (r *MemoryCharacterRepository) GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error)`
- L1098: `func (r *MemoryCharacterRepository) GetGuildMembers(guildID int) ([]*GuildMemberData, error)`
- L1110: `func (r *MemoryCharacterRepository) GetAllGuilds() ([]*GuildListItem, error)`
- L1139: `func (r *MemoryCharacterRepository) AddGuildMember(guildID, charID int, grade byte) error`
- L1170: `func (r *MemoryCharacterRepository) RemoveGuildMember(guildID, charID int) error`
- L1184: `func (r *MemoryCharacterRepository) SetGuildMemberGrade(guildID, charID int, grade byte) error`
- L1196: `func (r *MemoryCharacterRepository) SetGuildNotify(guildID int, charName, notify string) error`
- L1209: `func (r *MemoryCharacterRepository) DismissGuild(guildID int) error`
- L1222: `func (r *MemoryCharacterRepository) CreateGuildAcademy(guildID, masterCharID int) error`
- L1244: `func (r *MemoryCharacterRepository) GetGuildAcademy(guildID int) (*GuildAcademyData, error)`
- L1254: `func (r *MemoryCharacterRepository) GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error)`
- L1267: `func (r *MemoryCharacterRepository) GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error)`
- L1279: `func (r *MemoryCharacterRepository) AddAcademyMember(guildID, charID int, joinLevel byte) error`
- L1311: `func (r *MemoryCharacterRepository) RemoveAcademyMember(guildID, charID int) error`
- L1325: `func (r *MemoryCharacterRepository) SetAcademyNotify(guildID int, charName, notify string) error`

## pkg/database/database.go
- L18: `type AccountRepository interface`
- L23: `type accountStore struct`
- L29: `type MemoryAccountRepository struct`
- L37: `func NewMemoryAccountRepository() *MemoryAccountRepository`
- L67: `func (r *MemoryAccountRepository) load()`
- L90: `func (r *MemoryAccountRepository) save()`
- L110: `func (r *MemoryAccountRepository) GetAccountByUsername(username string) (*Account, error)`
- L123: `func (r *MemoryAccountRepository) CreateAccount(username, password string) (*Account, error)`

## pkg/database/mysql_account.go
- L14: `type MySQLAccountRepository struct`
- L19: `func NewMySQLAccountRepository(db *sql.DB) *MySQLAccountRepository`
- L26: `func (r *MySQLAccountRepository) initAuxiliaryTables()`
- L136: `func (r *MySQLAccountRepository) GetAccountByUsername(username string) (*Account, error)`
- L215: `func (r *MySQLAccountRepository) CreateAccount(username, password string) (*Account, error)`
- L245: `func (r *MySQLAccountRepository) IsIPBlocked(ip string) (bool, error)`
- L270: `func (r *MySQLAccountRepository) LogAccount(userID string, logType uint8, userNo int, worldNo uint8, errorCode uint16, ip string) error`
- L280: `func (r *MySQLAccountRepository) UpdateLoginSession(userNo int, ip string) error`
- L292: `func (r *MySQLAccountRepository) GetUserDateDiffCreated(userNo int) (int64, error)`
- L307: `func (r *MySQLAccountRepository) GetReturnUserDateDiffLastLog(userNo int) (int64, error)`
- L322: `func (r *MySQLAccountRepository) SetReturnUserItem(userNo int, itemNo int) error`
- L333: `func (r *MySQLAccountRepository) AddCash(userNo int, cashType uint8, cash int, reference string) error`
- L361: `func (r *MySQLAccountRepository) GetCash(userNo int, cashType uint8) (int, error)`
- L375: `func (r *MySQLAccountRepository) SetSerialKey(userNo int, serialKey int64) error`
- L386: `func (r *MySQLAccountRepository) LogLogout(worldNo uint8, userNo, charNo int, classVal, level uint8, exp int64, playTime int16) error`
- L397: `func (r *MySQLAccountRepository) GetDateDiffLastLogout(userNo int) (int64, error)`
- L415: `func (r *MySQLAccountRepository) AddReport(report *Report) (int, error)`
- L462: `func (r *MySQLAccountRepository) GetReports(reporterUserNo int, worldNo uint8) ([]*ReportSummary, error)`
- L490: `func (r *MySQLAccountRepository) DeleteReport(regNumber int) (int16, error)`
- L522: `func (r *MySQLAccountRepository) AddSameTimeChar(worldNo uint8, count int) (int, error)`
- L536: `func (r *MySQLAccountRepository) AddSameTimeCharClass(worldNo uint8, classVal uint8, level uint16, count int) error`
- L547: `func (r *MySQLAccountRepository) GetChargedItems(userNo int) ([]*ChargedItem, error)`

## pkg/database/mysql_account_test.go
- L16: `func setupTestDB(t *testing.T) (*sql.DB, *database.MySQLAccountRepository)`
- L38: `func TestMySQLAccountRepository_GetAccountByUsername(t *testing.T)`
- L65: `func TestMySQLAccountRepository_CreateAccount(t *testing.T)`
- L89: `func TestMySQLAccountRepository_LogAccountAndIPBlock(t *testing.T)`
- L123: `func TestMySQLAccountRepository_AdditionalProcedures(t *testing.T)`
- L258: `func md5HashString(s string) string`

## pkg/database/mysql_character.go
- L14: `type MySQLCharacterRepository struct`
- L19: `func NewMySQLCharacterRepository(db *sql.DB) *MySQLCharacterRepository`
- L25: `func (r *MySQLCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error)`
- L155: `func (r *MySQLCharacterRepository) GetCharacterByName(name string) (*Character, error)`
- L271: `func (r *MySQLCharacterRepository) GetCharacterByID(id int) (*Character, error)`
- L388: `func (r *MySQLCharacterRepository) CreateCharacter(ch *Character) error`
- L581: `func (r *MySQLCharacterRepository) DeleteCharacter(accountID int, slot byte) error`
- L628: `func (r *MySQLCharacterRepository) UpdateCharacter(ch *Character) error`
- L641: `func (r *MySQLCharacterRepository) getCharacterEquips(charNo int) (map[byte]uint16, error)`
- L665: `func (r *MySQLCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error)`
- L689: `func (r *MySQLCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error)`
- L713: `func (r *MySQLCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error)`
- L737: `func (r *MySQLCharacterRepository) SaveCharacterItem(item *Item) error`
- L753: `func (r *MySQLCharacterRepository) CreateItem(item *Item) error`
- L780: `func (r *MySQLCharacterRepository) DeleteItem(itemKey int64) error`
- L789: `func (r *MySQLCharacterRepository) SetItemStorage(itemKey int64, storageType byte, storageSlot int16) error`
- L801: `func (r *MySQLCharacterRepository) ExchangeItemStorage(itemKey1 int64, storageType1 byte, storageSlot1 int16, itemKey2 int64, storageType2 b`
- L827: `func (r *MySQLCharacterRepository) ArrangeInventory(charNo int, itemKeys []int64) error`
- L921: `func (r *MySQLCharacterRepository) UpdateItemFlags(itemKey int64, flags int) error`
- L930: `func (r *MySQLCharacterRepository) GetItemsByStorage(charNo int, storageType byte) ([]*Item, error)`
- L954: `func (r *MySQLCharacterRepository) AddCharMoney(charNo int, amount int64) error`
- L963: `func (r *MySQLCharacterRepository) SubCharMoney(charNo int, amount int64) (bool, error)`
- L978: `func (r *MySQLCharacterRepository) SaveCharLocation(charNo int, mapName string, x, y int, rotation byte) error`
- L991: `func (r *MySQLCharacterRepository) SaveCharStats(charNo int, hp, sp, hps, sps int16) error`
- L1004: `func (r *MySQLCharacterRepository) SaveCharLevelExp(charNo int, level int16, exp int64) error`
- L1017: `func (r *MySQLCharacterRepository) SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error`
- L1030: `func (r *MySQLCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error`
- L1045: `func (r *MySQLCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error`
- L1061: `func (r *MySQLCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error)`
- L1103: `func (r *MySQLCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error`
- L1138: `func (r *MySQLCharacterRepository) GetFriends(charNo int) ([]*FriendInfo, error)`
- L1164: `func (r *MySQLCharacterRepository) AddFriend(charNo, friendCharNo int) error`
- L1187: `func (r *MySQLCharacterRepository) RemoveFriend(charNo, friendCharNo int) error`
- L1201: `func (r *MySQLCharacterRepository) GetFriendPoints(charNo int) (int, error)`
- L1215: `func (r *MySQLCharacterRepository) SetFriendPoints(charNo int, points int) error`
- L1222: `func (r *MySQLCharacterRepository) CreateGuild(masterCharID int, name, password string, guildType byte) (int, error)`
- L1290: `func (r *MySQLCharacterRepository) GetGuild(guildID int) (*GuildData, error)`
- L1322: `func (r *MySQLCharacterRepository) GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error)`
- L1372: `func (r *MySQLCharacterRepository) GetGuildMembers(guildID int) ([]*GuildMemberData, error)`
- L1411: `func (r *MySQLCharacterRepository) GetAllGuilds() ([]*GuildListItem, error)`
- L1443: `func (r *MySQLCharacterRepository) AddGuildMember(guildID, charID int, grade byte) error`
- L1453: `func (r *MySQLCharacterRepository) RemoveGuildMember(guildID, charID int) error`
- L1460: `func (r *MySQLCharacterRepository) SetGuildMemberGrade(guildID, charID int, grade byte) error`
- L1467: `func (r *MySQLCharacterRepository) SetGuildNotify(guildID int, charName, notify string) error`
- L1478: `func (r *MySQLCharacterRepository) DismissGuild(guildID int) error`
- L1495: `func (r *MySQLCharacterRepository) CreateGuildAcademy(guildID, masterCharID int) error`
- L1506: `func (r *MySQLCharacterRepository) GetGuildAcademy(guildID int) (*GuildAcademyData, error)`
- L1536: `func (r *MySQLCharacterRepository) GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error)`
- L1584: `func (r *MySQLCharacterRepository) GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error)`
- L1623: `func (r *MySQLCharacterRepository) AddAcademyMember(guildID, charID int, joinLevel byte) error`
- L1633: `func (r *MySQLCharacterRepository) RemoveAcademyMember(guildID, charID int) error`
- L1640: `func (r *MySQLCharacterRepository) SetAcademyNotify(guildID int, charName, notify string) error`

## pkg/database/mysql_character_test.go
- L14: `func setupTestCharDB(t *testing.T) (*sql.DB, *database.MySQLCharacterRepository)`
- L36: `func TestMySQLCharacterRepository_Lifecycle(t *testing.T)`
- L194: `func TestMySQLCharacterRepository_StarterItemsAndShortcuts(t *testing.T)`
- L290: `func TestMySQLCharacterRepository_ZonePersistence(t *testing.T)`

## pkg/database/mysql_test.go
- L12: `func TestMySQLConnection(t *testing.T)`

## pkg/database/transfer.go
- L15: `type ClientTransfer struct`
- L25: `type ZoneTransfer struct`
- L38: `type TransferManager struct`
- L47: `func NewTransferManager(db ...*sql.DB) *TransferManager`
- L59: `func (tm *TransferManager) CreateTransfer(accountID int, username string, admin byte, hostIP string) string`
- L97: `func (tm *TransferManager) ConsumeTransfer(token string, raw []byte, username ...string) (*ClientTransfer, bool)`
- L151: `func (tm *TransferManager) CreateZoneTransfer(accountID int, username, charName string, charID int, randomID uint16, admin byte, hostIP stri`
- L180: `func (tm *TransferManager) ConsumeZoneTransfer(charName string, randomID uint16, raw []byte) (*ZoneTransfer, bool)`

## cmd/WorldService/main.go
- L18: `func main()`

