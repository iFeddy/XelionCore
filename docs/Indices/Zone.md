# Ãndice: ZoneService

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. Protocolos en pkg/protocol.

## pkg/server/zone/character.go
- L39: `func JobToClassName(job byte) string`
- L59: `type ZonePlayer struct`
- L91: `func NewZonePlayer(client *network.Client, ch *database.Character, classStore *classparam.Store) *ZonePlayer`
- L171: `func (p *ZonePlayer) RecalculateEquipmentStats(itemStore *item.Store)`
- L203: `func (p *ZonePlayer) WriteLook(pw *protocol.PacketWriter)`
- L216: `func (p *ZonePlayer) WriteEquipment(pw *protocol.PacketWriter)`
- L241: `func (p *ZonePlayer) WriteDetailedInfo(pw *protocol.PacketWriter)`
- L286: `func (p *ZonePlayer) WriteDetailedInfoExtra(pw *protocol.PacketWriter, levelUp bool)`
- L373: `func (p *ZonePlayer) WriteCharacterDisplay(pw *protocol.PacketWriter)`

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
- L142: `func (h *Handler) OnClientDisconnect(client *network.Client)`
- L185: `func (h *Handler) getPlayer(client *network.Client) (*ZonePlayer, bool)`
- L192: `func (h *Handler) handleTransferKey(client *network.Client, pkt *protocol.Packet)`
- L292: `func (h *Handler) sendInGameChunk(p *ZonePlayer)`
- L460: `func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet)`
- L473: `func (h *Handler) handleWalk(client *network.Client, pkt *protocol.Packet)`
- L490: `func (h *Handler) handleRun(client *network.Client, pkt *protocol.Packet)`
- L509: `func (h *Handler) handleMoveWalk(client *network.Client, pkt *protocol.Packet)`
- L532: `func (h *Handler) handleMoveRun(client *network.Client, pkt *protocol.Packet)`
- L560: `func (h *Handler) handleStop(client *network.Client, pkt *protocol.Packet)`
- L582: `func (h *Handler) handleJump(client *network.Client, pkt *protocol.Packet)`
- L593: `func (h *Handler) handleEmoticon(client *network.Client, pkt *protocol.Packet)`
- L609: `func (h *Handler) handleEmoticonStop(client *network.Client, pkt *protocol.Packet)`
- L627: `func (h *Handler) handleChatNormal(client *network.Client, pkt *protocol.Packet)`
- L654: `func (h *Handler) handleBeginRest(client *network.Client, pkt *protocol.Packet)`
- L666: `func (h *Handler) handleEndRest(client *network.Client, pkt *protocol.Packet)`
- L680: `func (h *Handler) SetItemStore(store *item.Store)`
- L686: `func (h *Handler) SetAccountRepo(repo *database.MySQLAccountRepository)`
- L693: `func (h *Handler) ActivePlayersCount() int`
- L700: `func (h *Handler) ActivePlayersSnapshot() []*ZonePlayer`
- L712: `func (h *Handler) handleReportAdd(client *network.Client, pkt *protocol.Packet)`
- L774: `func (h *Handler) handleReportCancel(client *network.Client, pkt *protocol.Packet)`
- L792: `func (h *Handler) handleReportGet(client *network.Client, pkt *protocol.Packet)`
- L822: `func (h *Handler) handleRewardInvenOpen(client *network.Client, pkt *protocol.Packet)`
- L873: `func (h *Handler) handleChargedInvenOpen(client *network.Client, pkt *protocol.Packet)`
- L913: `func (h *Handler) SetDataStore(ds *DataStore)`
- L941: `func (h *Handler) handleNPCClick(client *network.Client, pkt *protocol.Packet)`
- L975: `func (h *Handler) handleNPCMenuOpenAck(client *network.Client, pkt *protocol.Packet)`
- L986: `func (h *Handler) handleNPCActionCmd(client *network.Client, pkt *protocol.Packet)`
- L999: `func (h *Handler) handleNPCMenuCmd(client *network.Client, pkt *protocol.Packet)`
- L1009: `func (h *Handler) handleUseGate(client *network.Client, pkt *protocol.Packet)`
- L1028: `func (h *Handler) triggerGateWarp(player *ZonePlayer, gateArg string)`
- L1050: `func (h *Handler) WarpPlayer(player *ZonePlayer, targetMapID string, x, y int, rotation byte)`
- L1083: `func (h *Handler) handleMultyLinkSelect(client *network.Client, pkt *protocol.Packet)`
- L1113: `func (h *Handler) handleTownPortal(client *network.Client, pkt *protocol.Packet)`
- L1127: `func (h *Handler) handleWingFly(client *network.Client, pkt *protocol.Packet)`
- L1145: `func (h *Handler) sendCellChange(client *network.Client, box uint8, slot int16, it *database.Item, exchange uint16)`
- L1184: `func (h *Handler) findItemAt(charID int, box uint8, slot int16) *database.Item`
- L1204: `func (h *Handler) findFreeInvenSlot(charID int) (int16, bool)`
- L1228: `func (h *Handler) syncPlayerEquipsAndStats(player *ZonePlayer)`
- L1253: `func (h *Handler) broadcastEquipAppearance(player *ZonePlayer, equipSlot byte, itemID uint16, upgrade byte)`
- L1277: `func (h *Handler) broadcastUnequipAppearance(player *ZonePlayer, equipSlot byte)`
- L1293: `func (h *Handler) handleItemReloc(client *network.Client, pkt *protocol.Packet)`
- L1403: `func (h *Handler) handleItemSplit(client *network.Client, pkt *protocol.Packet)`
- L1462: `func (h *Handler) handleItemEquip(client *network.Client, pkt *protocol.Packet)`
- L1495: `func (h *Handler) handleItemRingEquip(client *network.Client, pkt *protocol.Packet)`
- L1521: `func (h *Handler) equipItemToSlot(player *ZonePlayer, srcItem *database.Item, invenSlot byte, equipSlot byte)`
- L1560: `func (h *Handler) canEquipItem(player *ZonePlayer, item *database.Item) bool`
- L1579: `func sortInventoryItems(items []*database.Item, store *item.Store)`
- L1614: `func (h *Handler) handleItemUnequip(client *network.Client, pkt *protocol.Packet)`
- L1668: `func (h *Handler) handleItemDrop(client *network.Client, pkt *protocol.Packet)`
- L1726: `func (h *Handler) handleItemPick(client *network.Client, pkt *protocol.Packet)`
- L1786: `func (h *Handler) handleItemUse(client *network.Client, pkt *protocol.Packet)`
- L1905: `func (h *Handler) handleAutoArrangeInven(client *network.Client, pkt *protocol.Packet)`
- L1955: `func (h *Handler) handleSoulStoneHPBuy(client *network.Client, pkt *protocol.Packet)`
- L2015: `func (h *Handler) handleSoulStoneSPBuy(client *network.Client, pkt *protocol.Packet)`
- L2078: `func (h *Handler) handleSoulStoneHPUse(client *network.Client, pkt *protocol.Packet)`
- L2120: `func (h *Handler) handleSoulStoneSPUse(client *network.Client, pkt *protocol.Packet)`

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
- L171: `func (mm *MapManager) GetOrCreateMap(mapID string) *ZoneMap`
- L223: `func (m *ZoneMap) AssignObjectID(p *ZonePlayer)`
- L233: `func (m *ZoneMap) AddPlayer(p *ZonePlayer)`
- L283: `func (m *ZoneMap) SpawnDroppedItem(itemID uint16, lot uint32, flags int, x, y int32, dropMobHandle uint16, ownerCharID int) *ZoneDroppedItem`
- L318: `func (m *ZoneMap) ConsumeDroppedItem(handle uint16) (*ZoneDroppedItem, bool)`
- L338: `func (m *ZoneMap) SendDroppedItems(p *ZonePlayer)`
- L360: `func (m *ZoneMap) SendNPCs(p *ZonePlayer)`
- L408: `func (m *ZoneMap) GetNPC(handle uint16) (*ZoneNPC, bool)`
- L416: `func (m *ZoneMap) GetAllNPCs() []*ZoneNPC`
- L427: `func (m *ZoneMap) GetGates() []*ZoneNPC`
- L436: `func (m *ZoneMap) FindGateNear(x, y int, maxDist float64) (*ZoneNPC, bool)`
- L452: `func (m *ZoneMap) RemovePlayer(p *ZonePlayer)`
- L466: `func (m *ZoneMap) Broadcast(pw *protocol.PacketWriter, excludeID uint16)`
- L479: `func (m *ZoneMap) MovePlayer(p *ZonePlayer, oldX, oldY, newX, newY int, walk bool)`
- L500: `func (m *ZoneMap) MovePlayerCmd(p *ZonePlayer, fromX, fromY, toX, toY int, walk bool)`
- L523: `func (m *ZoneMap) StopPlayer(p *ZonePlayer, x, y int)`
- L535: `func (m *ZoneMap) JumpPlayer(p *ZonePlayer)`
- L543: `func (m *ZoneMap) EmoticonPlayer(p *ZonePlayer, emoticonID byte)`
- L552: `func (m *ZoneMap) EmoticonStopPlayer(p *ZonePlayer, emoticonID byte, emoticonFrame uint16)`
- L561: `func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte)`

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
- L20: `func readServerPacket(r io.Reader) (*protocol.Packet, error)`
- L43: `func sendClientPacket(w io.Writer, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error`
- L67: `func TestZoneServer_Phase1_LifecycleAndHandshake(t *testing.T)`
- L123: `func TestZoneServer_Phase1_TransferKeyVerification(t *testing.T)`
- L220: `func TestZoneServer_Phase1_InvalidTransferRejection(t *testing.T)`
- L272: `func TestZoneServer_Phase1_ChecksumVerification(t *testing.T)`
- L359: `func TestZoneServer_ReportsAndLogoutLog(t *testing.T)`
- L486: `func TestZoneServer_Phase1_NPCSpawnsAndGateWarps(t *testing.T)`

## cmd/ZoneService/main.go
- L19: `func main()`

