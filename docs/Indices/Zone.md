# Ãndice: ZoneService

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. Protocolos en pkg/protocol.

## pkg/server/zone/character.go
- L39: `func JobToClassName(job byte) string`
- L59: `type ZonePlayer struct`
- L92: `func NewZonePlayer(client *network.Client, ch *database.Character, classStore *classparam.Store) *ZonePlayer`
- L172: `func (p *ZonePlayer) RecalculateEquipmentStats(itemStore *item.Store)`
- L204: `func (p *ZonePlayer) WriteLook(pw *protocol.PacketWriter)`
- L217: `func (p *ZonePlayer) WriteEquipment(pw *protocol.PacketWriter)`
- L242: `func (p *ZonePlayer) WriteDetailedInfo(pw *protocol.PacketWriter)`
- L287: `func (p *ZonePlayer) WriteDetailedInfoExtra(pw *protocol.PacketWriter, levelUp bool)`
- L374: `func (p *ZonePlayer) WriteCharacterDisplay(pw *protocol.PacketWriter)`

## pkg/server/zone/config.go
- L4: `type Config struct`

## pkg/server/zone/data_store.go
- L26: `type DataStore struct`
- L49: `func LoadAllData(dataDir string) (*DataStore, error)`

## pkg/server/zone/equipment_test.go
- L12: `func TestCanEquipItemChecksOfficialClassAndLevel(t *testing.T)`
- L34: `func TestSortInventoryItemsUsesOfficialCategoryOrder(t *testing.T)`
- L61: `func loadOfficialItemStore(t *testing.T) *item.Store`

## pkg/server/zone/handler.go
- L21: `type Handler struct`
- L38: `func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.`
- L51: `func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet)`
- L151: `func (h *Handler) OnClientDisconnect(client *network.Client)`
- L194: `func (h *Handler) getPlayer(client *network.Client) (*ZonePlayer, bool)`
- L201: `func (h *Handler) handleTransferKey(client *network.Client, pkt *protocol.Packet)`
- L301: `func (h *Handler) sendInGameChunk(p *ZonePlayer)`
- L469: `func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet)`
- L482: `func (h *Handler) handleWalk(client *network.Client, pkt *protocol.Packet)`
- L499: `func (h *Handler) handleRun(client *network.Client, pkt *protocol.Packet)`
- L518: `func (h *Handler) handleMoveWalk(client *network.Client, pkt *protocol.Packet)`
- L541: `func (h *Handler) handleMoveRun(client *network.Client, pkt *protocol.Packet)`
- L569: `func (h *Handler) handleStop(client *network.Client, pkt *protocol.Packet)`
- L591: `func (h *Handler) handleJump(client *network.Client, pkt *protocol.Packet)`
- L602: `func (h *Handler) handleEmoticon(client *network.Client, pkt *protocol.Packet)`
- L618: `func (h *Handler) handleEmoticonStop(client *network.Client, pkt *protocol.Packet)`
- L636: `func (h *Handler) handleChatNormal(client *network.Client, pkt *protocol.Packet)`
- L663: `func (h *Handler) handleBeginRest(client *network.Client, pkt *protocol.Packet)`
- L675: `func (h *Handler) handleEndRest(client *network.Client, pkt *protocol.Packet)`
- L689: `func (h *Handler) SetItemStore(store *item.Store)`
- L695: `func (h *Handler) SetAccountRepo(repo *database.MySQLAccountRepository)`
- L702: `func (h *Handler) ActivePlayersCount() int`
- L709: `func (h *Handler) ActivePlayersSnapshot() []*ZonePlayer`
- L721: `func (h *Handler) handleReportAdd(client *network.Client, pkt *protocol.Packet)`
- L783: `func (h *Handler) handleReportCancel(client *network.Client, pkt *protocol.Packet)`
- L801: `func (h *Handler) handleReportGet(client *network.Client, pkt *protocol.Packet)`
- L831: `func (h *Handler) handleRewardInvenOpen(client *network.Client, pkt *protocol.Packet)`
- L882: `func (h *Handler) handleChargedInvenOpen(client *network.Client, pkt *protocol.Packet)`
- L922: `func (h *Handler) SetDataStore(ds *DataStore)`
- L950: `func (h *Handler) handleNPCClick(client *network.Client, pkt *protocol.Packet)`
- L984: `func (h *Handler) handleNPCMenuOpenAck(client *network.Client, pkt *protocol.Packet)`
- L995: `func (h *Handler) handleNPCActionCmd(client *network.Client, pkt *protocol.Packet)`
- L1008: `func (h *Handler) handleNPCMenuCmd(client *network.Client, pkt *protocol.Packet)`
- L1018: `func (h *Handler) handleUseGate(client *network.Client, pkt *protocol.Packet)`
- L1037: `func (h *Handler) triggerGateWarp(player *ZonePlayer, gateArg string)`
- L1059: `func (h *Handler) WarpPlayer(player *ZonePlayer, targetMapID string, x, y int, rotation byte)`
- L1092: `func (h *Handler) handleMultyLinkSelect(client *network.Client, pkt *protocol.Packet)`
- L1122: `func (h *Handler) handleTownPortal(client *network.Client, pkt *protocol.Packet)`
- L1136: `func (h *Handler) handleWingFly(client *network.Client, pkt *protocol.Packet)`
- L1154: `func (h *Handler) sendCellChange(client *network.Client, box uint8, slot int16, it *database.Item, exchange uint16)`
- L1193: `func (h *Handler) findItemAt(charID int, box uint8, slot int16) *database.Item`
- L1213: `func (h *Handler) findFreeInvenSlot(charID int) (int16, bool)`
- L1237: `func (h *Handler) syncPlayerEquipsAndStats(player *ZonePlayer)`
- L1262: `func (h *Handler) broadcastEquipAppearance(player *ZonePlayer, equipSlot byte, itemID uint16, upgrade byte)`
- L1286: `func (h *Handler) broadcastUnequipAppearance(player *ZonePlayer, equipSlot byte)`
- L1302: `func (h *Handler) handleItemReloc(client *network.Client, pkt *protocol.Packet)`
- L1412: `func (h *Handler) handleItemSplit(client *network.Client, pkt *protocol.Packet)`
- L1471: `func (h *Handler) handleItemEquip(client *network.Client, pkt *protocol.Packet)`
- L1504: `func (h *Handler) handleItemRingEquip(client *network.Client, pkt *protocol.Packet)`
- L1530: `func (h *Handler) equipItemToSlot(player *ZonePlayer, srcItem *database.Item, invenSlot byte, equipSlot byte)`
- L1569: `func (h *Handler) canEquipItem(player *ZonePlayer, item *database.Item) bool`
- L1588: `func sortInventoryItems(items []*database.Item, store *item.Store)`
- L1623: `func (h *Handler) handleItemUnequip(client *network.Client, pkt *protocol.Packet)`
- L1677: `func (h *Handler) handleItemDrop(client *network.Client, pkt *protocol.Packet)`
- L1735: `func (h *Handler) handleItemPick(client *network.Client, pkt *protocol.Packet)`
- L1795: `func (h *Handler) handleItemUse(client *network.Client, pkt *protocol.Packet)`
- L1914: `func (h *Handler) handleAutoArrangeInven(client *network.Client, pkt *protocol.Packet)`
- L1964: `func (h *Handler) handleSoulStoneHPBuy(client *network.Client, pkt *protocol.Packet)`
- L2024: `func (h *Handler) handleSoulStoneSPBuy(client *network.Client, pkt *protocol.Packet)`
- L2087: `func (h *Handler) handleSoulStoneHPUse(client *network.Client, pkt *protocol.Packet)`
- L2129: `func (h *Handler) handleSoulStoneSPUse(client *network.Client, pkt *protocol.Packet)`
- L2171: `func (h *Handler) handleLogoutReady(client *network.Client, pkt *protocol.Packet)`
- L2183: `func (h *Handler) handleLogoutCancel(client *network.Client, pkt *protocol.Packet)`
- L2196: `func (h *Handler) handleNormalLogout(client *network.Client, pkt *protocol.Packet)`

## pkg/server/zone/map.go
- L16: `type ZoneNPC struct`
- L32: `type ZoneDroppedItem struct`
- L45: `type ZoneMap struct`
- L61: `func NewZoneMap(id uint16, mapID, name string) *ZoneMap`
- L75: `type MapManager struct`
- L84: `func NewMapManager(fieldStore *field.Store) *MapManager`
- L93: `func (mm *MapManager) SetNPCStore(npcStore *npc.Store, mobInfoSHN *shn.File)`
- L114: `func (mm *MapManager) resolveMobID(mobName string) uint16`
- L122: `func (mm *MapManager) loadMapNPCs(zMap *ZoneMap)`
- L174: `func (mm *MapManager) GetOrCreateMap(mapID string) *ZoneMap`
- L226: `func (m *ZoneMap) AssignObjectID(p *ZonePlayer)`
- L236: `func (m *ZoneMap) AddPlayer(p *ZonePlayer)`
- L286: `func (m *ZoneMap) SpawnDroppedItem(itemID uint16, lot uint32, flags int, x, y int32, dropMobHandle uint16, ownerCharID int) *ZoneDroppedItem`
- L321: `func (m *ZoneMap) ConsumeDroppedItem(handle uint16) (*ZoneDroppedItem, bool)`
- L341: `func (m *ZoneMap) SendDroppedItems(p *ZonePlayer)`
- L363: `func (m *ZoneMap) SendNPCs(p *ZonePlayer)`
- L413: `func (m *ZoneMap) GetNPC(handle uint16) (*ZoneNPC, bool)`
- L421: `func (m *ZoneMap) GetAllNPCs() []*ZoneNPC`
- L432: `func (m *ZoneMap) GetGates() []*ZoneNPC`
- L441: `func (m *ZoneMap) FindGateNear(x, y int, maxDist float64) (*ZoneNPC, bool)`
- L457: `func (m *ZoneMap) RemovePlayer(p *ZonePlayer)`
- L471: `func (m *ZoneMap) Broadcast(pw *protocol.PacketWriter, excludeID uint16)`
- L484: `func (m *ZoneMap) MovePlayer(p *ZonePlayer, oldX, oldY, newX, newY int, walk bool)`
- L505: `func (m *ZoneMap) MovePlayerCmd(p *ZonePlayer, fromX, fromY, toX, toY int, walk bool)`
- L528: `func (m *ZoneMap) StopPlayer(p *ZonePlayer, x, y int)`
- L540: `func (m *ZoneMap) JumpPlayer(p *ZonePlayer)`
- L548: `func (m *ZoneMap) EmoticonPlayer(p *ZonePlayer, emoticonID byte)`
- L557: `func (m *ZoneMap) EmoticonStopPlayer(p *ZonePlayer, emoticonID byte, emoticonFrame uint16)`
- L566: `func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte)`

## pkg/server/zone/server.go
- L30: `type Server struct`
- L44: `func NewServer(cfg *Config, charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldS`
- L90: `func (s *Server) Start() error`
- L104: `func (s *Server) Stop()`
- L117: `func (s *Server) SetDataStore(ds *DataStore)`
- L121: `func (s *Server) SetItemStore(store *item.Store)`
- L125: `func (s *Server) SetAccountRepo(repo *database.MySQLAccountRepository)`
- L135: `func (s *Server) Handler() *Handler`

## pkg/server/zone/zone_phase2_test.go
- L16: `func waitOpCode(t *testing.T, conn net.Conn, op uint16) *protocol.Packet`
- L32: `func TestZoneServer_Phase2_InventoryEquipConsumablesAndSoulstones(t *testing.T)`

## pkg/server/zone/zone_test.go
- L21: `func TestHandler_NormalLogoutReturnsToWorldSelect(t *testing.T)`
- L90: `func readServerPacket(r io.Reader) (*protocol.Packet, error)`
- L113: `func sendClientPacket(w io.Writer, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error`
- L137: `func TestZoneServer_Phase1_LifecycleAndHandshake(t *testing.T)`
- L193: `func TestZoneServer_Phase1_TransferKeyVerification(t *testing.T)`
- L290: `func TestZoneServer_Phase1_InvalidTransferRejection(t *testing.T)`
- L342: `func TestZoneServer_Phase1_ChecksumVerification(t *testing.T)`
- L429: `func TestZoneServer_ReportsAndLogoutLog(t *testing.T)`
- L556: `func TestZoneServer_Phase1_NPCSpawnsAndGateWarps(t *testing.T)`

## cmd/ZoneService/main.go
- L19: `func main()`

