package zone

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler processes incoming packets for the Zone Server.
type Handler struct {
	mu          sync.RWMutex
	charRepo    database.CharacterRepository
	accountRepo *database.MySQLAccountRepository
	transfers   *database.TransferManager
	dataStore   *DataStore
	classStore  *classparam.Store
	fieldStore  *field.Store
	npcStore    *npc.Store
	mobInfoSHN  *shn.File
	chkEngine   *shn.ChecksumEngine
	itemStore   *item.Store
	players     map[*network.Client]*ZonePlayer
	mapManager  *MapManager
}

// NewHandler initializes a new Zone packet handler.
func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.Store, chkEngine *shn.ChecksumEngine) *Handler {
	return &Handler{
		charRepo:   charRepo,
		transfers:  transfers,
		classStore: classStore,
		fieldStore: fieldStore,
		chkEngine:  chkEngine,
		players:    make(map[*network.Client]*ZonePlayer),
		mapManager: NewMapManager(fieldStore),
	}
}

// HandlePacket routes packets to their specific handlers using official PROTO_NC_* opcodes.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_MAP_LOGIN_REQ:
		h.handleTransferKey(client, pkt)
	case protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD, protocol.PROTO_NC_CHAR_INGAME_REQ:
		h.handleClientReady(client, pkt)
	case protocol.PROTO_NC_ACT_WALK_REQ:
		h.handleWalk(client, pkt)
	case protocol.PROTO_NC_ACT_RUN_REQ:
		h.handleRun(client, pkt)
	case protocol.PROTO_NC_ACT_MOVEWALK_CMD:
		h.handleMoveWalk(client, pkt)
	case protocol.PROTO_NC_ACT_MOVERUN_CMD:
		h.handleMoveRun(client, pkt)
	case protocol.PROTO_NC_ACT_STOP_REQ:
		h.handleStop(client, pkt)
	case protocol.PROTO_NC_ACT_JUMP_CMD:
		h.handleJump(client, pkt)
	case protocol.PROTO_NC_ACT_EMOTICON_CMD:
		h.handleEmoticon(client, pkt)
	case protocol.PROTO_NC_ACT_EMOTICONSTOP_CMD:
		h.handleEmoticonStop(client, pkt)
	case protocol.PROTO_NC_ACT_CHAT_REQ:
		h.handleChatNormal(client, pkt)
	case protocol.PROTO_NC_ACT_REST_START_REQ:
		h.handleBeginRest(client, pkt)
	case protocol.PROTO_NC_ACT_REST_STOP_REQ:
		h.handleEndRest(client, pkt)
	case protocol.PROTO_NC_DATA_REPORT_ADD_REQ:
		h.handleReportAdd(client, pkt)
	case protocol.PROTO_NC_DATA_REPORT_CANCEL_REQ:
		h.handleReportCancel(client, pkt)
	case protocol.PROTO_NC_DATA_REPORT_GET_REQ:
		h.handleReportGet(client, pkt)
	case protocol.PROTO_NC_ACT_ENDOFTRADE_CMD:
		// 0x200B (sp_NC_ACT_ENDOFTRADE_CMD en Zone.pdb): si el jugador no está en un trade activo, se ignora silenciosamente.
	case protocol.PROTO_NC_ACT_NPCCLICK_CMD:
		h.handleNPCClick(client, pkt)
	case protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK:
		h.handleNPCMenuOpenAck(client, pkt)
	case protocol.PROTO_NC_ACT_NPC_ACTION_CMD:
		h.handleNPCActionCmd(client, pkt)
	case protocol.PROTO_NC_ACT_NPC_MENU_CMD:
		h.handleNPCMenuCmd(client, pkt)
	case protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ:
		h.handleUseGate(client, pkt)
	case protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ:
		h.handleMultyLinkSelect(client, pkt)
	case protocol.PROTO_NC_MAP_TOWNPORTAL_REQ:
		h.handleTownPortal(client, pkt)
	case protocol.PROTO_NC_MAP_WING_FLY_REQ:
		h.handleWingFly(client, pkt)
	case protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ:
		h.handleRewardInvenOpen(client, pkt)
	case protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ:
		h.handleChargedInvenOpen(client, pkt)
	case protocol.PROTO_NC_ITEM_RELOC_REQ:
		h.handleItemReloc(client, pkt)
	case protocol.PROTO_NC_ITEM_SPLIT_REQ:
		h.handleItemSplit(client, pkt)
	case protocol.PROTO_NC_ITEM_EQUIP_REQ:
		h.handleItemEquip(client, pkt)
	case protocol.PROTO_NC_ITEM_RINGEQUIP_REQ:
		h.handleItemRingEquip(client, pkt)
	case protocol.PROTO_NC_ITEM_UNEQUIP_REQ:
		h.handleItemUnequip(client, pkt)
	case protocol.PROTO_NC_ITEM_USE_REQ:
		h.handleItemUse(client, pkt)
	case protocol.PROTO_NC_ITEM_DROP_REQ:
		h.handleItemDrop(client, pkt)
	case protocol.PROTO_NC_ITEM_PICK_REQ:
		h.handleItemPick(client, pkt)
	case protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ:
		h.handleAutoArrangeInven(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ:
		h.handleSoulStoneHPBuy(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_SP_BUY_REQ:
		h.handleSoulStoneSPBuy(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_HP_USE_REQ:
		h.handleSoulStoneHPUse(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_SP_USE_REQ:
		h.handleSoulStoneSPUse(client, pkt)
	case protocol.PROTO_NC_MISC_HEARTBEAT_ACK, protocol.PROTO_NC_MISC_GAMETIME_REQ:
		// Heartbeat pong received, ignore
	default:
		log.Printf("[Zone][%s] Unhandled OpCode: 0x%04X (Header=%d, Type=%d, Len=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type, len(pkt.Raw))
	}
}

// OnClientDisconnect cleans up player resources on TCP disconnect.
func (h *Handler) OnClientDisconnect(client *network.Client) {
	h.mu.Lock()
	player, exists := h.players[client]
	delete(h.players, client)
	h.mu.Unlock()

	if exists && player != nil {
		if player.CurrentMap != nil {
			player.CurrentMap.RemovePlayer(player)
		}

		if player.Char != nil {
			mapID := player.Char.TownName
			if player.CurrentMap != nil && player.CurrentMap.MapID != "" {
				mapID = player.CurrentMap.MapID
			}
			_ = h.charRepo.SaveCharLocation(player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, player.Rotation)
			_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)
			_ = h.charRepo.UpdateCharacter(player.Char)
		}

		// Persistir registro oficial de desconexión y tiempo jugado en tLogout (usp_Log_setLogout)
		if h.accountRepo != nil && player.Char != nil {
			playMin := int16(time.Since(player.LoginTime).Minutes())
			if playMin < 1 {
				playMin = 1
			}
			worldNo := uint8(0)
			_ = h.accountRepo.LogLogout(
				worldNo,
				player.Char.AccountID,
				player.Char.ID,
				player.Char.Job,
				player.Char.Level,
				player.Char.TCharacterState.Exp,
				playMin,
			)
			log.Printf("[Zone][%s] Logout registrado en tLogout: '%s' (User: %d, Char: %d, PlayTime: %d min)",
				client.RemoteAddr(), player.Char.Name, player.Char.AccountID, player.Char.ID, playMin)
		}
	}
}

func (h *Handler) getPlayer(client *network.Client) (*ZonePlayer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	p, ok := h.players[client]
	return p, ok
}

func (h *Handler) handleTransferKey(client *network.Client, pkt *protocol.Packet) {
	req, err := protocol.Read_PROTO_NC_MAP_LOGIN_REQ(pkt.Reader)
	if err != nil {
		log.Printf("[Zone][%s] Error reading PROTO_NC_MAP_LOGIN_REQ: %v", client.RemoteAddr(), err)
		return
	}

	charName := strings.TrimSpace(req.Chardata.Charid)
	randomID := req.Chardata.Wldmanhandle

	log.Printf("[Zone][%s] PROTO_NC_MAP_LOGIN_REQ received for character '%s' (RandomID: %d, HasChecksums: %v)",
		client.RemoteAddr(), charName, randomID, req.HasChecksums)

	// Verificación oficial binaria de DataFileCheck (Zone.exe 0x0045264b)
	if h.chkEngine != nil && req.HasChecksums {
		var clientChecksums [49]shn.Name8
		for i := 0; i < 49; i++ {
			clientChecksums[i] = req.Checksum[i]
		}
		valid, wrongIdx := h.chkEngine.VerifyDataFiles(clientChecksums)
		if !valid {
			fileName := h.chkEngine.GetFileName(wrongIdx)
			expectedHash := h.chkEngine.GetRawHash(wrongIdx)
			clientHash := string(req.Checksum[wrongIdx][:])
			log.Printf("[Zone][%s] DataFileCheck FAILED at index %d (%s): client hash='%s', server expected='%s'",
				client.RemoteAddr(), wrongIdx, fileName, clientHash, expectedHash)

			// Error oficial en Zone.exe (0x00452727): código 0x0147 (327 decimal) + nWrongDataFileIndex (1 byte)
			failPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(0x0147, byte(wrongIdx))
			_ = client.SendPacket(failPkt)
			client.Close()
			return
		}
		log.Printf("[Zone][%s] DataFileCheck PASSED for all 49 .shn files", client.RemoteAddr())
	}

	transfer, ok := h.transfers.ConsumeZoneTransfer(charName, randomID, pkt.Raw)
	if !ok {
		log.Printf("[Zone][%s] Transfer verification failed for '%s' (RandomID: %d)",
			client.RemoteAddr(), charName, randomID)
		errPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(4, 0) // RequestedCharacterIDNotMatching
		_ = client.SendPacket(errPkt)
		client.Close()
		return
	}

	char, err := h.charRepo.GetCharacterByName(transfer.CharName)
	if err != nil {
		char, err = h.charRepo.GetCharacterByID(transfer.CharID)
	}

	if err != nil || char == nil {
		log.Printf("[Zone][%s] Character '%s' not found in database: %v", client.RemoteAddr(), transfer.CharName, err)
		errPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(4, 0)
		_ = client.SendPacket(errPkt)
		client.Close()
		return
	}

	player := NewZonePlayer(client, char, h.classStore)
	mapID := char.TownName
	if mapID == "" {
		mapID = "RouN"
	}
	zMap := h.mapManager.GetOrCreateMap(mapID)
	player.CurrentMap = zMap
	zMap.AssignObjectID(player)
	h.mu.Lock()
	h.players[client] = player
	h.mu.Unlock()

	client.Username = transfer.Username
	client.AccountID = transfer.AccountID
	client.Admin = transfer.Admin
	client.IsAuthenticated = true

	log.Printf("[Zone][%s] Character '%s' authenticated successfully. Dispatching In-Game Chunk...",
		client.RemoteAddr(), char.Name)

	h.sendInGameChunk(player)
}

// sendInGameChunk despacha la secuencia estricta de paquetes de entrada a Zona (cliente NA2016):
//  1. 0x1038 (PROTO_NC_CHAR_CLIENT_BASE_CMD)             - Siempre
//  2. 0x1039 (PROTO_NC_CHAR_CLIENT_SHAPE_CMD)            - Siempre
//  3. 0x103A (PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)      - Siempre
//  4. 0x103B (PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD)       - Siempre
//  5. 0x10CE (PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD)       - Siempre
//  6. 0x10D7 (PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD)     - Siempre
//  7. 0x103E (PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)          - Siempre
//  8. 0x104A (PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD)      - Siempre (incondicional, 0 si está vacío)
//  9. 0x103D (PROTO_NC_CHAR_CLIENT_SKILL_CMD)            - Solo si len(activeSkills) > 0 (PartMark = 3)
// 10. 0x1048 (PROTO_NC_CHAR_CLIENT_GAME_CMD)             - Siempre
// 11. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 8)      - Solo si len > 0 (Equip)
// 12. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 9)      - Solo si len > 0 (Inven)
// 13. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 12)     - Solo si len > 0 (House)
// 14. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 15)     - Siempre (ActionItem, incluso vacío)
// 15. 0x1049 (PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD)        - Siempre
// 16. 0x3801/0x3802/0x3803 (Card Regist)                 - Solo si tiene cartas registradas
// 17. 0x1802 (PROTO_NC_CHAR_MAPLOGIN_ACK)                - Siempre (último paquete del bloque)
func (h *Handler) sendInGameChunk(p *ZonePlayer) {
	// 1. 0x1038 - PROTO_NC_CHAR_CLIENT_BASE_CMD (Siempre, wire len = 107, payload = 105)
	infoPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD)
	p.WriteDetailedInfo(infoPkt)
	_ = p.Client.SendPacket(infoPkt)

	// 2. 0x1039 - PROTO_NC_CHAR_CLIENT_SHAPE_CMD (Siempre, wire len = 6, payload = 4)
	lookPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD)
	p.WriteLook(lookPkt)
	_ = p.Client.SendPacket(lookPkt)

	// 3. 0x103A - PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD (Siempre, wire len = 8, payload = 6)
	questsBusy := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)
	questsBusy.WriteUint32(uint32(p.Char.ID))
	questsBusy.WriteByte(1) // bNeedClear = 1
	questsBusy.WriteByte(0) // nNumOfDoingQuest = 0
	_ = p.Client.SendPacket(questsBusy)

	// 4. 0x103B - PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD (Siempre, wire len = 14, payload = 12)
	questsDone := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD)
	questsDone.WriteUint32(uint32(p.Char.ID))
	questsDone.WriteUint16(0) // nTotalDoneQuest
	questsDone.WriteUint16(0) // nTotalDoneQuestSize
	questsDone.WriteUint16(0) // nDoneQuestCount
	questsDone.WriteUint16(0) // nIndex
	_ = p.Client.SendPacket(questsDone)

	// 5. 0x10CE - PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD (Siempre, wire len = 8, payload = 6)
	questsRead := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD)
	questsRead.WriteUint32(uint32(p.Char.ID))
	questsRead.WriteUint16(0) // nNumOfReadQuest = 0
	_ = p.Client.SendPacket(questsRead)

	// 6. 0x10D7 - PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD (Siempre, wire len = 8, payload = 6)
	questsRepeat := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD)
	questsRepeat.WriteUint32(uint32(p.Char.ID))
	questsRepeat.WriteUint16(0) // nNumOfRepeatQuest = 0
	_ = p.Client.SendPacket(questsRepeat)

	// 7. 0x103E - PROTO_NC_CHAR_CLIENT_PASSIVE_CMD (Siempre, wire len = 4, payload = 2)
	skills, _ := h.charRepo.GetCharacterSkills(p.Char.ID)
	var activeSkills []*database.CharacterSkill
	var passiveCount uint16
	for _, sk := range skills {
		activeSkills = append(activeSkills, sk)
	}

	passiveSkills := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)
	passiveSkills.WriteUint16(passiveCount)
	_ = p.Client.SendPacket(passiveSkills)

	// 8. 0x104A - PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD (CRÍTICO: Siempre, incondicional)
	type chargedBuffEntry struct {
		ID       uint16
		TimeLeft uint32
		KeepTime uint32
		Value    uint32
	}
	var chargedBuffs []chargedBuffEntry
	buffPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD)
	buffPkt.WriteUint16(uint16(len(chargedBuffs))) // 0 si no tiene buffs (2 bytes de payload)
	for _, buff := range chargedBuffs {
		buffPkt.WriteUint16(buff.ID)
		buffPkt.WriteUint32(buff.TimeLeft)
		buffPkt.WriteUint32(buff.KeepTime)
		buffPkt.WriteUint32(buff.Value)
	}
	_ = p.Client.SendPacket(buffPkt)

	// 9. 0x103D - PROTO_NC_CHAR_CLIENT_SKILL_CMD (Solo si tiene skills activas, con PartMark = 3)
	if len(activeSkills) > 0 {
		skillPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SKILL_CMD)
		skillPkt.WriteByte(0)                           // restempow (1 byte)
		skillPkt.WriteByte(3)                           // PartMark = 3 (first=1 | last=1, 1 byte)
		skillPkt.WriteUint16(uint16(len(activeSkills))) // nMaxNum (2 bytes)
		skillPkt.WriteUint32(uint32(p.Char.ID))         // chrregnum (4 bytes)
		skillPkt.WriteUint16(uint16(len(activeSkills))) // number (2 bytes)
		for _, sk := range activeSkills {
			skillPkt.WriteUint16(sk.SkillNo) // skillid (2 bytes)
			skillPkt.WriteUint32(0)          // cooltime (4 bytes)
			skillPkt.WriteUint16(0)          // empow (2 bytes)
			skillPkt.WriteUint32(0)          // mastery (4 bytes)
		}
		_ = p.Client.SendPacket(skillPkt)
	}

	// 10. 0x1048 - PROTO_NC_CHAR_CLIENT_GAME_CMD (Siempre, wire len = 6, payload = 4)
	gamePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD)
	gamePkt.WriteUint32(0xFFFFFFFF) // partyno = 0xFFFF, nRaidNo = 0xFFFF (Sin party ni raid)
	_ = p.Client.SendPacket(gamePkt)

	// 11-14. 0x1047 - PROTO_NC_CHAR_CLIENT_ITEM_CMD:
	// - Box 8 (Equip):      Solo si len > 0
	// - Box 9 (Inven):      Solo si len > 0
	// - Box 12 (House):     Solo si len > 0
	// - Box 15 (ActionItem): Siempre (incluso vacío: numofitem=0, box=15, flag=1)
	charItems, _ := h.charRepo.GetCharacterItems(p.Char.ID)
	itemsByBox := make(map[uint8][]*database.Item)
	for _, it := range charItems {
		boxType := it.StorageType
		if boxType == 0 {
			boxType = 8
		}
		itemsByBox[boxType] = append(itemsByBox[boxType], it)
		if boxType == 8 {
			if p.Char.Equips == nil {
				p.Char.Equips = make(map[byte]uint16)
			}
			p.Char.Equips[byte(it.Storage)] = it.ItemID
		}
	}

	sendItemBox := func(box uint8, items []*database.Item, sendIfEmpty bool) {
		if len(items) == 0 && !sendIfEmpty {
			return
		}
		itemPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD)
		if len(items) > 0 {
			itemPkt.WriteByte(byte(len(items))) // numofitem
			itemPkt.WriteByte(box)              // box
			itemPkt.WriteByte(1)                // flag = 1
			for _, it := range items {
				loc := (uint16(box) << 10) | (uint16(it.Storage) & 0x03FF)
				attrBuf := h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
				dataSize := byte(2 + 2 + len(attrBuf))
				itemPkt.WriteByte(dataSize)    // datasize: 2 (location) + 2 (itemid) + len(attr)
				itemPkt.WriteUint16(loc)       // location: ITEM_INVEN (uint16 bitfield: (box << 10) | slot)
				itemPkt.WriteUint16(it.ItemID) // ItemID
				itemPkt.WriteBytes(attrBuf)    // attr: dynamic payload emulating iac_fill2buffer
			}
		} else {
			itemPkt.WriteByte(0)   // numofitem = 0
			itemPkt.WriteByte(box) // box = 15 (ActionItem)
			itemPkt.WriteByte(1)   // flag = 1
		}
		_ = p.Client.SendPacket(itemPkt)
	}

	// 11. 0x1047 (Box 8: Equip) - Solo si len > 0
	sendItemBox(8, itemsByBox[8], false)
	// 12. 0x1047 (Box 9: Inven) - Solo si len > 0
	sendItemBox(9, itemsByBox[9], false)
	// 13. 0x1047 (Box 12: House) - Solo si len > 0
	sendItemBox(12, itemsByBox[12], false)
	// 14. 0x1047 (Box 15: ActionItem) - Siempre (incluso vacío)
	sendItemBox(15, itemsByBox[15], true)

	// 15. 0x1049 - PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD (Siempre, wire len = 8, payload = 6)
	titles := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD)
	titles.WriteByte(0)   // CurrentTitle = 0
	titles.WriteByte(0)   // CurrentTitleElement = 0
	titles.WriteUint16(0) // CurrentTitleMobID = 0
	titles.WriteUint16(0) // NumOfTitle = 0
	_ = p.Client.SendPacket(titles)

	// 16. 0x3801 / 0x3802 / 0x3803 (Card Regist) - Solo si tiene cartas registradas
	var cardRegistPackets []*protocol.PacketWriter
	for _, cardPkt := range cardRegistPackets {
		_ = p.Client.SendPacket(cardPkt)
	}

	// 17. 0x1802 - PROTO_NC_CHAR_MAPLOGIN_ACK (CRÍTICO: Absolutamente el último paquete del bloque)
	p.RecalculateEquipmentStats(h.itemStore)
	extraPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_MAPLOGIN_ACK)
	p.WriteDetailedInfoExtra(extraPkt, false)
	_ = p.Client.SendPacket(extraPkt)
}

func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	log.Printf("[Zone][%s] ClientReady received from '%s'. Spawning into '%s' (MapID: %s, Serial: %d)...",
		client.RemoteAddr(), player.Char.Name, player.CurrentMap.Name, player.CurrentMap.MapID, player.CurrentMap.ID)

	// Add player to map and broadcast spawn
	player.CurrentMap.AddPlayer(player)
}

func (h *Handler) handleWalk(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), true)
}

func (h *Handler) handleRun(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), false)
}

// handleMoveWalk implementa sp_NC_ACT_MOVEWALK_CMD (0x2017) según Zone.pdb (RVA 0x56C60):
// Lee PROTO_NC_ACT_MOVEWALK_CMD (from: SHINE_XY_TYPE 8B, to: SHINE_XY_TYPE 8B) y transmite 0x2018 (PROTO_NC_ACT_SOMEONEMOVEWALK_CMD).
func (h *Handler) handleMoveWalk(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	fromX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	fromY, _ := pkt.Reader.ReadInt32()
	toX, _ := pkt.Reader.ReadInt32()
	toY, _ := pkt.Reader.ReadInt32()

	if fromX == toX && fromY == toY {
		return
	}

	player.CurrentMap.MovePlayerCmd(player, int(fromX), int(fromY), int(toX), int(toY), true)
}

// handleMoveRun implementa sp_NC_ACT_MOVERUN_CMD (0x2019) según Zone.pdb (RVA 0x56FD0):
// Lee PROTO_NC_ACT_MOVERUN_CMD (from: SHINE_XY_TYPE 8B, to: SHINE_XY_TYPE 8B) y transmite 0x201A (PROTO_NC_ACT_SOMEONEMOVERUN_CMD).
func (h *Handler) handleMoveRun(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	fromX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	fromY, _ := pkt.Reader.ReadInt32()
	toX, _ := pkt.Reader.ReadInt32()
	toY, _ := pkt.Reader.ReadInt32()

	if fromX == toX && fromY == toY {
		return
	}

	player.CurrentMap.MovePlayerCmd(player, int(fromX), int(fromY), int(toX), int(toY), false)

	// Verificar activación de Gate si el jugador corrió hacia el portal
	if time.Since(player.LastWarpTime) >= 3*time.Second {
		if gate, found := player.CurrentMap.FindGateNear(int(toX), int(toY), 180.0); found && gate.GateToWhere != "" {
			h.triggerGateWarp(player, gate.GateToWhere)
		}
	}
}

func (h *Handler) handleStop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	newX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.StopPlayer(player, int(newX), int(newY))

	// Verificar si el jugador pisó la zona de activación de un portal/gate
	if time.Since(player.LastWarpTime) >= 3*time.Second {
		if gate, found := player.CurrentMap.FindGateNear(int(newX), int(newY), 180.0); found && gate.GateToWhere != "" {
			h.triggerGateWarp(player, gate.GateToWhere)
		}
	}
}

func (h *Handler) handleJump(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	player.CurrentMap.JumpPlayer(player)
}

// handleEmoticon maneja 0x2020 (PROTO_NC_ACT_EMOTICON_CMD -> sp_NC_ACT_EMOTICON_CMD / sp_DoEmotion en Zone.pdb).
// Lee PROTO_NC_ACT_EMOTICON_CMD (emoticonid: uint8, 1 byte) y transmite 0x2021 (PROTO_NC_ACT_SOMEONEEMOTICON_CMD).
func (h *Handler) handleEmoticon(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	emoticonID, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	player.CurrentMap.EmoticonPlayer(player, emoticonID)
}

// handleEmoticonStop maneja 0x2022 (PROTO_NC_ACT_EMOTICONSTOP_CMD -> sp_NC_ACT_EMOTICONSTOP_CMD en Zone.pdb).
// Lee STOPEMOTICON_DESCRIPT (emoticonid: uint8 1B, emoticonframe: uint16 2B) y transmite 0x2023 (PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD).
func (h *Handler) handleEmoticonStop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	emoticonID, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	emoticonFrame, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	player.CurrentMap.EmoticonStopPlayer(player, emoticonID, emoticonFrame)
}

func (h *Handler) handleChatNormal(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	textLen, err := pkt.Reader.ReadByte()
	if err != nil || textLen == 0 {
		return
	}

	text, err := pkt.Reader.ReadString(int(textLen))
	if err != nil {
		return
	}
	text = strings.TrimSpace(text)

	log.Printf("[Zone][Chat][%s] %s: %s", player.CurrentMap.Name, player.Char.Name, text)

	color := byte(0x2A) // Normal chat color
	if player.Client.Admin > 0 {
		color = 0x03 // GM gold color
	}

	player.CurrentMap.NormalChat(player, text, color)
}

func (h *Handler) handleBeginRest(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	player.State = 4 // Resting
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_REST_START_ACK)
	pw.WriteUint16(0x0A81) // OK to rest
	_ = client.SendPacket(pw)
}

func (h *Handler) handleEndRest(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	player.State = 1 // Normal
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_REST_STOP_ACK)
	pw.WriteUint16(0x0A81)
	_ = client.SendPacket(pw)
}

// SetAccountRepo configura el repositorio de cuentas y logs oficial (tLogout, tReport, etc.).
// SetItemStore sets the item data store for equipment and attribute serialization.
func (h *Handler) SetItemStore(store *item.Store) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.itemStore = store
}

func (h *Handler) SetAccountRepo(repo *database.MySQLAccountRepository) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.accountRepo = repo
}

// ActivePlayersCount retorna la cantidad actual de jugadores conectados en la zona.
func (h *Handler) ActivePlayersCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.players)
}

// ActivePlayersSnapshot retorna una instantánea de los jugadores para telemetría de clases.
func (h *Handler) ActivePlayersSnapshot() []*ZonePlayer {
	h.mu.RLock()
	defer h.mu.RUnlock()
	list := make([]*ZonePlayer, 0, len(h.players))
	for _, p := range h.players {
		if p != nil {
			list = append(list, p)
		}
	}
	return list
}

func (h *Handler) handleReportAdd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	// Estructura oficial PROTO_NC_DATA_REPORT_ADD_REQ (2369 bytes)
	// dwReporter(4), hReporter(2), nWorldNo(1), sReporterCharID(20), sHarmer(20), byReportType(16), byReason(256), nChatLogSize(2), byChatLog(...)
	_ = pkt.Reader.Skip(4) // dwReporter
	hReporter, _ := pkt.Reader.ReadUint16()
	nWorldNo, _ := pkt.Reader.ReadByte()
	sReporterCharID, _ := pkt.Reader.ReadString(20)
	sHarmer, _ := pkt.Reader.ReadString(20)
	byReportType, _ := pkt.Reader.ReadString(16)
	byReason, _ := pkt.Reader.ReadString(256)
	chatLogSize, _ := pkt.Reader.ReadUint16()
	var chatLog string
	if chatLogSize > 0 && pkt.Reader.Remaining() >= int(chatLogSize) {
		chatLog, _ = pkt.Reader.ReadString(int(chatLogSize))
	} else if pkt.Reader.Remaining() > 0 {
		chatLog, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
	}

	sReporterCharID = strings.TrimSpace(sReporterCharID)
	if sReporterCharID == "" {
		sReporterCharID = player.Char.Name
	}
	sHarmer = strings.TrimSpace(sHarmer)
	byReportType = strings.TrimSpace(byReportType)
	byReason = strings.TrimSpace(byReason)

	log.Printf("[Zone][%s] Denuncia oficial recibida de '%s' contra '%s' (Tipo: %s, Razón: %s)",
		client.RemoteAddr(), sReporterCharID, sHarmer, byReportType, byReason)

	var regNo int
	var err error
	if h.accountRepo != nil {
		rep := &database.Report{
			Reporter:       player.Char.AccountID,
			ReporterCharID: sReporterCharID,
			Harmer:         sHarmer,
			WorldNo:        nWorldNo,
			ReportType:     byReportType,
			Reason:         byReason,
			ChatLog:        chatLog,
		}
		regNo, err = h.accountRepo.AddReport(rep)
	}

	var errCode uint16 = 0
	if err != nil {
		errCode = 1 // Error interno
	} else if regNo == -1 {
		errCode = 2 // Cooldown 1 hora por reportero
	} else if regNo == -2 {
		errCode = 3 // Cooldown 24 horas contra el mismo objetivo
	}

	ack := protocol.New_PROTO_NC_DATA_REPORT_ADD_ACK(errCode, hReporter, sHarmer, int32(regNo))
	_ = client.SendPacket(ack)
}

func (h *Handler) handleReportCancel(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	regNumber, _ := pkt.Reader.ReadInt32()
	hReporter, _ := pkt.Reader.ReadUint16()

	var errCode int16 = 0
	if h.accountRepo != nil {
		errCode, _ = h.accountRepo.DeleteReport(int(regNumber))
	}

	ack := protocol.New_PROTO_NC_DATA_REPORT_CANCEL_ACK(regNumber, hReporter, uint16(errCode))
	_ = client.SendPacket(ack)
}

func (h *Handler) handleReportGet(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	hReporter, _ := pkt.Reader.ReadUint16()
	_ = pkt.Reader.Skip(4) // dwReporter
	nWorldNo, _ := pkt.Reader.ReadByte()

	if h.accountRepo != nil {
		summaries, err := h.accountRepo.GetReports(player.Char.AccountID, nWorldNo)
		if err == nil && len(summaries) > 0 {
			for _, s := range summaries {
				ack := protocol.New_PROTO_NC_DATA_REPORT_GET_ACK(hReporter, 0, int32(s.RegNumber), s.Result, s.Remark)
				_ = client.SendPacket(ack)
			}
			return
		}
	}

	ack := protocol.New_PROTO_NC_DATA_REPORT_GET_ACK(hReporter, 1, 0, 0, "")
	_ = client.SendPacket(ack)
}

// handleRewardInvenOpen implementa sp_NC_ITEM_REWARDINVENOPEN_REQ (0x302C) y so_ply_FillRewardInven (0x302D)
// según Zone.pdb (RewardInven.cpp, RVA 0xFB060 y 0xFB310):
//   - Lee uint16 page (24 slots por página: from = page * 24, to = page * 24 + 23)
//   - Consulta tItem con nStorageType = 2 (Reward Inventory)
//   - Responde con 0x302D (PROTO_NC_ITEM_REWARDINVENOPEN_ACK): itemcounter (1B) + PROTO_ITEMPACKET_INFORM[itemcounter]
func (h *Handler) handleRewardInvenOpen(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	page, err := pkt.Reader.ReadUint16()
	if err != nil {
		page = 0
	}
	fromSlot := int16(page * 24)
	toSlot := int16(page*24 + 23)

	var pageItems []*database.Item
	if h.charRepo != nil {
		allRewardItems, err := h.charRepo.GetItemsByStorage(player.Char.ID, 2) // nStorageType = 2 (Reward Storage en Zone.exe)
		if err == nil {
			for _, it := range allRewardItems {
				if it != nil && it.Storage >= fromSlot && it.Storage <= toSlot {
					pageItems = append(pageItems, it)
				}
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REWARDINVENOPEN_ACK)
	ack.WriteByte(byte(len(pageItems))) // itemcounter (1 byte)
	for _, it := range pageItems {
		loc := (uint16(2) << 10) | (uint16(it.Storage) & 0x03FF)
		var attrBuf []byte
		if h.itemStore != nil {
			attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
		} else {
			attrBuf = []byte{1}
		}
		dataSize := byte(2 + 2 + len(attrBuf))
		ack.WriteByte(dataSize)    // datasize
		ack.WriteUint16(loc)       // location: ITEM_INVEN
		ack.WriteUint16(it.ItemID) // ItemID
		ack.WriteBytes(attrBuf)    // attr (iac_fill2buffer)
	}
	_ = client.SendPacket(ack)
}

// handleChargedInvenOpen implementa sp_NC_ITEM_CHARGEDINVENOPEN_REQ (0x3020) y so_ply_FillChargedInven (0x3021)
// según Zone.pdb (RVA 0x913C0 y 0x91F30):
//   - Lee uint16 page (24 slots por página: from = page * 24, to = page * 24 + 23)
//   - Consulta tChargeItem (usp_Charge_ItemSelectList / usp_ChargeItem_GetList) filtrando por userNo e isDraw = 0
//   - Responde con 0x3021 (PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK):
//     ErrorCode (uint16 = 0x1041), nPartMark (uint8 = 3), NumOfChargedItem (uint16),
//     seguido de PROTO_CHARGED_ITEM_INFO[NumOfChargedItem] (16 bytes cada uno: ItemOrderNo, ItemCode, ItemAmount, ItemRegisterDate).
func (h *Handler) handleChargedInvenOpen(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	page, err := pkt.Reader.ReadUint16()
	if err != nil {
		page = 0
	}

	var pageItems []*database.ChargedItem
	if h.accountRepo != nil {
		allCharged, err := h.accountRepo.GetChargedItems(player.Char.AccountID)
		if err == nil && len(allCharged) > 0 {
			startIdx := int(page) * 24
			endIdx := startIdx + 24
			if startIdx < len(allCharged) {
				if endIdx > len(allCharged) {
					endIdx = len(allCharged)
				}
				pageItems = allCharged[startIdx:endIdx]
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK)
	ack.WriteUint16(0x1041)               // ErrorCode = 0x1041 (4161 = éxito en so_ply_FillChargedInven, Zone.exe 0x492F8E)
	ack.WriteByte(3)                      // nPartMark = 3 (first=1 | last=1, Zone.exe 0x493007)
	ack.WriteUint16(uint16(len(pageItems))) // NumOfChargedItem
	for _, ci := range pageItems {
		ack.WriteUint32(ci.RowNo)                                  // 0x00: ItemOrderNo (4 bytes)
		ack.WriteUint32(ci.GoodsNo)                                // 0x04: ItemCode (4 bytes)
		ack.WriteUint32(ci.Amount)                                 // 0x08: ItemAmount (4 bytes)
		ack.WriteUint32(database.PackShineDateTime(ci.RegisterDate)) // 0x0C: ItemRegisterDate (ShineDateTime, 4 bytes)
	}
	_ = client.SendPacket(ack)
}

// SetDataStore sets the entire loaded game data store for the zone handler.
func (h *Handler) SetDataStore(ds *DataStore) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.dataStore = ds
	if ds != nil {
		if ds.ClassStore != nil {
			h.classStore = ds.ClassStore
		}
		if ds.FieldStore != nil {
			h.fieldStore = ds.FieldStore
		}
		if ds.ItemStore != nil {
			h.itemStore = ds.ItemStore
		}
		if ds.ChkEngine != nil {
			h.chkEngine = ds.ChkEngine
		}
		if ds.NPCStore != nil {
			h.npcStore = ds.NPCStore
		}
		if ds.MobInfoSHN != nil {
			h.mobInfoSHN = ds.MobInfoSHN
		}
		h.mapManager.SetNPCStore(ds.NPCStore, ds.MobInfoSHN)
	}
}

// handleNPCClick procesa 0x200A (PROTO_NC_ACT_NPCCLICK_CMD -> sp_NC_ACT_NPCCLICK_CMD en Zone.pdb).
func (h *Handler) handleNPCClick(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	npcHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	npc, found := player.CurrentMap.GetNPC(npcHandle)
	if !found || npc == nil {
		log.Printf("[Zone][%s] Clicked unknown NPC handle: %d", client.RemoteAddr(), npcHandle)
		return
	}

	log.Printf("[Zone][%s] Player '%s' clicked NPC '%s' (Handle: %d, MobID: %d, Role: %s, RoleArg: %s)",
		client.RemoteAddr(), player.Char.Name, npc.Name, npc.Handle, npc.MobID, npc.Role, npc.RoleArg)

	// Si es un Gate / Portal, ejecutar la transición
	if npc.IsGate && npc.GateToWhere != "" {
		h.triggerGateWarp(player, npc.GateToWhere)
		return
	}

	// PROTO_NC_ACT_NPCMENUOPEN_REQ (0x201C, wire len = 4, payload = 2)
	// struct PROTO_NC_ACT_NPCMENUOPEN_REQ { unsigned short mobid; };
	openReq := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ)
	openReq.WriteUint16(npc.MobID)
	_ = client.SendPacket(openReq)
}

// handleNPCMenuOpenAck procesa 0x201D (PROTO_NC_ACT_NPCMENUOPEN_ACK -> sp_NC_ACT_NPCMENUOPEN_ACK en Zone.pdb).
func (h *Handler) handleNPCMenuOpenAck(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Player '%s' opened NPC menu (Ack: %d)",
		client.RemoteAddr(), player.Char.Name, ack)
}

// handleNPCActionCmd procesa 0x2072 (PROTO_NC_ACT_NPC_ACTION_CMD -> sp_NC_ACT_NPC_ACTION_CMD en Zone.pdb).
func (h *Handler) handleNPCActionCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	nType, _ := pkt.Reader.ReadByte()
	nNPCHandle, _ := pkt.Reader.ReadUint16()
	nECode, _ := pkt.Reader.ReadUint32()
	log.Printf("[Zone][%s] NPC action command received: Type=%d, NPCHandle=%d, ECode=%d",
		client.RemoteAddr(), nType, nNPCHandle, nECode)
}

// handleNPCMenuCmd procesa 0x2073 (PROTO_NC_ACT_NPC_MENU_CMD -> sp_NC_ACT_NPC_MENU_CMD en Zone.pdb).
func (h *Handler) handleNPCMenuCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] NPC menu response received: Ack=%d", client.RemoteAddr(), ack)
}

// handleUseGate procesa 0x2070 (PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ -> sp_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ en Zone.pdb).
func (h *Handler) handleUseGate(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	// PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK (0x2071, nError: 0)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK)
	ack.WriteUint16(0) // OK
	_ = client.SendPacket(ack)

	// Buscar el gate más cercano
	gate, found := player.CurrentMap.FindGateNear(player.Char.XPos, player.Char.YPos, 350.0)
	if found && gate != nil && gate.GateToWhere != "" {
		h.triggerGateWarp(player, gate.GateToWhere)
	}
}

// triggerGateWarp consulta LinkTable.json y ejecuta la teletransportación oficial entre mapas.
func (h *Handler) triggerGateWarp(player *ZonePlayer, gateArg string) {
	if player == nil || player.Client == nil || h.dataStore == nil || h.dataStore.NPCStore == nil {
		return
	}

	if time.Since(player.LastWarpTime) < 3*time.Second {
		return
	}

	link, found := h.dataStore.NPCStore.GetLinkByArg(gateArg)
	if !found {
		log.Printf("[Zone][%s] Gate '%s' not found in LinkTable", player.Client.RemoteAddr(), gateArg)
		return
	}

	log.Printf("[Zone][%s] Gate '%s' activated -> Warping to %s (%d, %d, Dir: %d)",
		player.Client.RemoteAddr(), gateArg, link.MapServer, link.X, link.Y, link.Direct)

	h.WarpPlayer(player, link.MapServer, int(link.X), int(link.Y), byte(link.Direct))
}

// WarpPlayer traslada al jugador entre mapas oficiales, actualizando mapas, coordenadas y persistencia.
func (h *Handler) WarpPlayer(player *ZonePlayer, targetMapID string, x, y int, rotation byte) {
	if player == nil || player.Client == nil {
		return
	}

	targetMapID = strings.TrimSpace(targetMapID)
	if targetMapID == "" {
		targetMapID = "RouN"
	}

	oldMap := player.CurrentMap
	if oldMap != nil {
		oldMap.RemovePlayer(player)
	}

	newMap := h.mapManager.GetOrCreateMap(targetMapID)
	player.CurrentMap = newMap
	player.Char.TownName = newMap.MapID
	player.Char.XPos = x
	player.Char.YPos = y
	player.Rotation = rotation
	player.LastWarpTime = time.Now()

	_ = h.charRepo.SaveCharLocation(player.Char.ID, newMap.MapID, x, y, rotation)

	log.Printf("[Zone][%s] Player '%s' entering new map '%s' (%d, %d, Dir: %d)",
		player.Client.RemoteAddr(), player.Char.Name, newMap.MapID, x, y, rotation)

	newMap.AssignObjectID(player)
	h.sendInGameChunk(player)
}

// handleMultyLinkSelect procesa 0x181F (PROTO_NC_MAP_MULTY_LINK_SELECT_REQ -> sp_NC_MAP_MULTY_LINK_SELECT_REQ en Zone.pdb).
func (h *Handler) handleMultyLinkSelect(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	linkMapName, err := pkt.Reader.ReadString(12)
	if err != nil {
		return
	}
	destMap := strings.TrimRight(linkMapName, "\x00")

	// PROTO_NC_MAP_MULTY_LINK_SELECT_ACK (0x1820, wire len = 4, payload = 2)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK)
	ack.WriteUint16(0) // OK
	_ = client.SendPacket(ack)

	spawnX := 6208
	spawnY := 6405
	if h.fieldStore != nil {
		if f, ok := h.fieldStore.GetField(destMap); ok && f != nil && len(f.Respawn.Spots) > 0 {
			spawnX = f.Respawn.Spots[0].X
			spawnY = f.Respawn.Spots[0].Y
		}
	}

	h.WarpPlayer(player, destMap, spawnX, spawnY, 0x55)
}

// handleTownPortal procesa 0x181A (PROTO_NC_MAP_TOWNPORTAL_REQ -> sp_NC_MAP_TOWNPORTAL_REQ en Zone.pdb).
func (h *Handler) handleTownPortal(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	portalIndex, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Town portal request (Index: %d) from '%s'",
		client.RemoteAddr(), portalIndex, player.Char.Name)

	h.WarpPlayer(player, "RouN", 6208, 6405, 0x55)
}

// handleWingFly procesa 0x1816 (PROTO_NC_MAP_WING_FLY_REQ -> sp_NC_MAP_WING_FLY_REQ en Zone.pdb).
func (h *Handler) handleWingFly(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	where, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Wing fly request (Where: %d) from '%s'",
		client.RemoteAddr(), where, player.Char.Name)

	h.WarpPlayer(player, "RouN", 6208, 6405, 0x55)
}

// sendCellChange implementa so_ply_ChangeCell (Zone.exe VA 0x518260):
//   - Si box == 8 ((location & 0xFC00) == 0x2000): envía PROTO_NC_ITEM_EQUIPCHANGE_CMD (0x3002)
//     (exchange: uint16, location: uint8, itemid: uint16 + attr si itemid != 0xFFFF)
//   - Si box != 8: envía PROTO_NC_ITEM_CELLCHANGE_CMD (0x3001)
//     (exchange: uint16, location: uint16, itemid: uint16 + attr si itemid != 0xFFFF)
func (h *Handler) sendCellChange(client *network.Client, box uint8, slot int16, it *database.Item, exchange uint16) {
	if box == 8 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD)
		pw.WriteUint16(exchange)
		pw.WriteByte(byte(slot))
		if it == nil || it.ItemID == 0xFFFF {
			pw.WriteUint16(0xFFFF)
		} else {
			pw.WriteUint16(it.ItemID)
			var attrBuf []byte
			if h.itemStore != nil {
				attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
			} else {
				attrBuf = []byte{byte(it.Flags)}
			}
			pw.WriteBytes(attrBuf)
		}
		_ = client.SendPacket(pw)
	} else {
		loc := (uint16(box) << 10) | (uint16(slot) & 0x03FF)
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CELLCHANGE_CMD)
		pw.WriteUint16(exchange)
		pw.WriteUint16(loc)
		if it == nil || it.ItemID == 0xFFFF {
			pw.WriteUint16(0xFFFF)
		} else {
			pw.WriteUint16(it.ItemID)
			var attrBuf []byte
			if h.itemStore != nil {
				attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
			} else {
				attrBuf = []byte{byte(it.Flags)}
			}
			pw.WriteBytes(attrBuf)
		}
		_ = client.SendPacket(pw)
	}
}

func (h *Handler) findItemAt(charID int, box uint8, slot int16) *database.Item {
	if h.charRepo == nil {
		return nil
	}
	items, err := h.charRepo.GetCharacterItems(charID)
	if err != nil {
		return nil
	}
	for _, it := range items {
		st := it.StorageType
		if st == 0 {
			st = 8
		}
		if st == box && it.Storage == slot {
			return it
		}
	}
	return nil
}

func (h *Handler) findFreeInvenSlot(charID int) (int16, bool) {
	if h.charRepo == nil {
		return 0, false
	}
	items, err := h.charRepo.GetCharacterItems(charID)
	if err != nil {
		return 0, false
	}
	used := make(map[int16]bool)
	for _, it := range items {
		if it.StorageType == 9 {
			used[it.Storage] = true
		}
	}
	for slot := int16(0); slot < 96; slot++ {
		if !used[slot] {
			return slot, true
		}
	}
	return 0, false
}

// syncPlayerEquipsAndStats actualiza el mapa de equipos del jugador, recalcula stats de equipo
// y envía PROTO_NC_CHAR_STATE_CHANGE_CMD (0x104E) si los stats cambiaron.
func (h *Handler) syncPlayerEquipsAndStats(player *ZonePlayer) {
	if player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	items, err := h.charRepo.GetCharacterItems(player.Char.ID)
	if err != nil {
		return
	}
	equips := make(map[byte]uint16)
	for _, it := range items {
		st := it.StorageType
		if st == 0 {
			st = 8
		}
		if st == 8 {
			equips[byte(it.Storage)] = it.ItemID
		}
	}
	player.Char.Equips = equips
	player.RecalculateEquipmentStats(h.itemStore)
}

// broadcastEquipAppearance transmite el cambio visual de equipamiento al mapa según Zone.exe:
//   - Slot 12 (Weapon): PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD (0x1C05, 9B)
//   - Otros slots:      PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD (0x1C03, 6B)
func (h *Handler) broadcastEquipAppearance(player *ZonePlayer, equipSlot byte, itemID uint16, upgrade byte) {
	if player == nil || player.CurrentMap == nil {
		return
	}
	if equipSlot == slotWeapon {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD)
		pw.WriteUint16(player.MapObjectID) // handle
		pw.WriteUint16(itemID)          // item
		pw.WriteByte(upgrade)           // upgrade
		pw.WriteByte(equipSlot)         // nSlotNum
		pw.WriteUint16(0xFFFF)          // currentmobid
		pw.WriteByte(0)                 // currentkilllevel
		player.CurrentMap.Broadcast(pw, 0)
	} else {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD)
		pw.WriteUint16(player.MapObjectID) // handle
		pw.WriteUint16(itemID)          // item
		pw.WriteByte(upgrade)           // upgrade
		pw.WriteByte(equipSlot)         // nSlotNum
		player.CurrentMap.Broadcast(pw, 0)
	}
}

// broadcastUnequipAppearance transmite PROTO_NC_BRIEFINFO_UNEQUIP_CMD (0x1C04, 3B) al mapa.
func (h *Handler) broadcastUnequipAppearance(player *ZonePlayer, equipSlot byte) {
	if player == nil || player.CurrentMap == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_UNEQUIP_CMD)
	pw.WriteUint16(player.MapObjectID) // handle
	pw.WriteByte(equipSlot)         // slot
	player.CurrentMap.Broadcast(pw, 0)
}

// handleItemReloc implementa sp_NC_ITEM_RELOC_REQ (0x300B) y gds_NC_ITEMDB_RELOCSUC_ACK / EXCHANGESUC_ACK (Zone.exe VA 0x533070 / 0x533350):
//   - Lee from (ITEM_INVEN uint16) y to (ITEM_INVEN uint16)
//   - Si el destino está vacío: mueve el ítem (RelocateItem), envía 0x3001/0x3002 para origen (vacío) y destino (ocupado)
//   - Si el destino tiene un ítem del mismo ID y apilable (< MaxLot): combina stacks
//   - Si el destino tiene otro ítem: intercambia posiciones (ExchangeItem)
//   - Envía PROTO_NC_ITEM_RELOC_ACK (0x300C) con err = 0x0241 (577)
func (h *Handler) handleItemReloc(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	fromLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	toLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	fromBox := uint8(fromLoc >> 10)
	fromSlot := int16(fromLoc & 0x03FF)
	toBox := uint8(toLoc >> 10)
	toSlot := int16(toLoc & 0x03FF)

	if fromBox == toBox && fromSlot == toSlot {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
		ack.WriteUint16(0x0241)
		_ = client.SendPacket(ack)
		return
	}

	srcItem := h.findItemAt(player.Char.ID, fromBox, fromSlot)
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
		ack.WriteUint16(0x0242) // Reloc fail
		_ = client.SendPacket(ack)
		return
	}

	dstItem := h.findItemAt(player.Char.ID, toBox, toSlot)
	if dstItem == nil {
		// Mover a celda vacía
		_ = h.charRepo.SetItemStorage(srcItem.ItemKey, toBox, toSlot)
		srcItem.StorageType = toBox
		srcItem.Storage = toSlot

		h.sendCellChange(client, fromBox, fromSlot, nil, 0xFFFF)
		h.sendCellChange(client, toBox, toSlot, srcItem, 0xFFFF)
	} else if dstItem.ItemID == srcItem.ItemID && toBox == 9 && fromBox == 9 {
		// Verificar si es apilable (MaxLot > 1)
		maxLot := uint32(1)
		if h.itemStore != nil {
			if info, ok := h.itemStore.GetItemInfo(srcItem.ItemID); ok && info != nil && info.MaxLot > 1 {
				maxLot = info.MaxLot
			}
		}
		srcLot := uint32(srcItem.Flags)
		if srcLot == 0 {
			srcLot = 1
		}
		dstLot := uint32(dstItem.Flags)
		if dstLot == 0 {
			dstLot = 1
		}

		if maxLot > 1 && dstLot < maxLot {
			total := srcLot + dstLot
			if total <= maxLot {
				dstItem.Flags = int32(total)
				_ = h.charRepo.UpdateItemFlags(dstItem.ItemKey, int(dstItem.Flags))
				_ = h.charRepo.DeleteItem(srcItem.ItemKey)
				h.sendCellChange(client, fromBox, fromSlot, nil, 0xFFFF)
				h.sendCellChange(client, toBox, toSlot, dstItem, 0xFFFF)
			} else {
				dstItem.Flags = int32(maxLot)
				srcItem.Flags = int32(total - maxLot)
				_ = h.charRepo.UpdateItemFlags(dstItem.ItemKey, int(dstItem.Flags))
				_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
				h.sendCellChange(client, fromBox, fromSlot, srcItem, 0xFFFF)
				h.sendCellChange(client, toBox, toSlot, dstItem, 0xFFFF)
			}
		} else {
			_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, toBox, toSlot, dstItem.ItemKey, fromBox, fromSlot)
			srcItem.StorageType = toBox
			srcItem.Storage = toSlot
			dstItem.StorageType = fromBox
			dstItem.Storage = fromSlot
			h.sendCellChange(client, toBox, toSlot, srcItem, fromLoc)
			h.sendCellChange(client, fromBox, fromSlot, dstItem, toLoc)
		}
	} else {
		// Intercambio de celdas (Exchange)
		_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, toBox, toSlot, dstItem.ItemKey, fromBox, fromSlot)
		srcItem.StorageType = toBox
		srcItem.Storage = toSlot
		dstItem.StorageType = fromBox
		dstItem.Storage = fromSlot
		h.sendCellChange(client, toBox, toSlot, srcItem, fromLoc)
		h.sendCellChange(client, fromBox, fromSlot, dstItem, toLoc)
	}

	if fromBox == 8 || toBox == 8 {
		h.syncPlayerEquipsAndStats(player)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
	ack.WriteUint16(0x0241) // Éxito oficial en Zone.exe (0x0241)
	_ = client.SendPacket(ack)
}

// handleItemSplit implementa sp_NC_ITEM_SPLIT_REQ (0x300D) y gds_NC_ITEMDB_SPLITSUC_ACK (Zone.exe VA 0x539740):
//   - Lee from (uint16), to (uint16), lot (uint8)
//   - Reduce el stack origen en `lot` y crea un nuevo ítem en `to` con `Flags = lot`
//   - Envía 0x3001 para ambas celdas y PROTO_NC_ITEM_SPLIT_ACK (0x300E) con err = 0x0441 (1089)
func (h *Handler) handleItemSplit(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	fromLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	toLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	splitLot, err := pkt.Reader.ReadByte()
	if err != nil || splitLot == 0 {
		return
	}

	fromBox := uint8(fromLoc >> 10)
	fromSlot := int16(fromLoc & 0x03FF)
	toBox := uint8(toLoc >> 10)
	toSlot := int16(toLoc & 0x03FF)

	srcItem := h.findItemAt(player.Char.ID, fromBox, fromSlot)
	dstItem := h.findItemAt(player.Char.ID, toBox, toSlot)
	if srcItem == nil || dstItem != nil || srcItem.Flags <= int32(splitLot) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
		ack.WriteUint16(0x0442) // Split fail
		_ = client.SendPacket(ack)
		return
	}

	srcItem.Flags -= int32(splitLot)
	_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))

	newItem := &database.Item{StorageType: toBox, Owner: player.Char.ID, Storage: int16(toSlot), ItemID: srcItem.ItemID, Flags: int32(int(splitLot))}
	err = h.charRepo.CreateItem(newItem)
	if err != nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
		ack.WriteUint16(0x0442)
		_ = client.SendPacket(ack)
		return
	}


	h.sendCellChange(client, fromBox, fromSlot, srcItem, 0xFFFF)
	h.sendCellChange(client, toBox, toSlot, newItem, 0xFFFF)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
	ack.WriteUint16(0x0441) // Éxito oficial en Zone.exe (0x0441)
	_ = client.SendPacket(ack)
}

// handleItemEquip implementa sp_NC_ITEM_EQUIP_REQ (0x300F) y gds_NC_ITEMDB_EQUIPSUC_ACK (Zone.exe VA 0x52CBA0):
//   - Lee invenslot (uint8) del inventario (Box 9)
//   - Determina el slot de equipamiento desde ItemInfo.shn (info.Equip)
//   - Mueve o intercambia el ítem con Box 8 (Equip), recalcula stats y transmite cambio visual al mapa
//   - Envía PROTO_NC_ITEM_EQUIP_ACK (0x3011) con err = 0x0281 (641)
func (h *Handler) handleItemEquip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	equipSlot := slotWeapon
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(srcItem.ItemID); ok && info != nil {
			if info.Equip > 0 && info.Equip != 255 {
				equipSlot = byte(info.Equip)
			}
		}
	}

	h.equipItemToSlot(player, srcItem, invenSlot, equipSlot)
}

// handleItemRingEquip implementa sp_NC_ITEM_RINGEQUIP_REQ (0x3010):
//   - Lee invenslot (uint8) y ringslot (uint8)
func (h *Handler) handleItemRingEquip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	ringSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	h.equipItemToSlot(player, srcItem, invenSlot, ringSlot)
}

func (h *Handler) equipItemToSlot(player *ZonePlayer, srcItem *database.Item, invenSlot byte, equipSlot byte) {
	client := player.Client
	if !h.canEquipItem(player, srcItem) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	dstEquip := h.findItemAt(player.Char.ID, 8, int16(equipSlot))
	invenLoc := (uint16(9) << 10) | (uint16(invenSlot) & 0x03FF)
	equipLoc := (uint16(8) << 10) | (uint16(equipSlot) & 0x03FF)

	if dstEquip == nil {
		_ = h.charRepo.SetItemStorage(srcItem.ItemKey, 8, int16(equipSlot))
		srcItem.StorageType = 8
		srcItem.Storage = int16(equipSlot)

		h.sendCellChange(client, 9, int16(invenSlot), nil, 0xFFFF)
		h.sendCellChange(client, 8, int16(equipSlot), srcItem, 0xFFFF)
	} else {
		_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, 8, int16(equipSlot), dstEquip.ItemKey, 9, int16(invenSlot))
		srcItem.StorageType = 8
		srcItem.Storage = int16(equipSlot)
		dstEquip.StorageType = 9
		dstEquip.Storage = int16(invenSlot)

		h.sendCellChange(client, 8, int16(equipSlot), srcItem, invenLoc)
		h.sendCellChange(client, 9, int16(invenSlot), dstEquip, equipLoc)
	}

	h.syncPlayerEquipsAndStats(player)
	h.broadcastEquipAppearance(player, equipSlot, srcItem.ItemID, 0)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
	ack.WriteUint16(0x0281) // Éxito oficial en Zone.exe (0x0281)
	_ = client.SendPacket(ack)
}

func (h *Handler) canEquipItem(player *ZonePlayer, item *database.Item) bool {
	if player == nil || player.Char == nil || item == nil {
		return false
	}
	if h.itemStore == nil {
		if player.Client != nil {
			log.Printf("[Zone][%s] cannot validate equipment item %d: official item data is not loaded",
				player.Client.RemoteAddr(), item.ItemID)
		}
		return false
	}
	info, ok := h.itemStore.GetItemInfo(item.ItemID)
	if !ok || info == nil || uint32(player.Char.Level) < info.DemandLv {
		return false
	}
	allowed, known := h.itemStore.CanUseItemByClassID(item.ItemID, player.Char.Job)
	return known && allowed
}

func sortInventoryItems(items []*database.Item, store *item.Store) {
	sort.SliceStable(items, func(i, j int) bool {
		leftOrder, leftGrade, leftID := int(^uint(0)>>1), uint32(0), items[i].ItemID
		rightOrder, rightGrade, rightID := int(^uint(0)>>1), uint32(0), items[j].ItemID
		if store != nil {
			if info, ok := store.GetItemInfo(items[i].ItemID); ok && info != nil {
				if info.SortOrder >= 0 {
					leftOrder = info.SortOrder
				}
				leftGrade = info.Grade
			}
			if info, ok := store.GetItemInfo(items[j].ItemID); ok && info != nil {
				if info.SortOrder >= 0 {
					rightOrder = info.SortOrder
				}
				rightGrade = info.Grade
			}
		}
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		if leftGrade != rightGrade {
			return leftGrade < rightGrade
		}
		if leftID != rightID {
			return leftID < rightID
		}
		return items[i].Storage < items[j].Storage
	})
}

// handleItemUnequip implementa sp_NC_ITEM_UNEQUIP_REQ (0x3012) y gds_NC_ITEMDB_UNEQUIPSUC_ACK (Zone.exe VA 0x53B0A0):
//   - Lee equipslot (uint8) e invenslot (uint8)
//   - Mueve el ítem de Box 8 a Box 9, recalcula stats, transmite PROTO_NC_BRIEFINFO_UNEQUIP_CMD (0x1C04)
//   - Envía PROTO_NC_ITEM_EQUIP_ACK (0x3011) con err = 0x0281 (641)
func (h *Handler) handleItemUnequip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcEquip := h.findItemAt(player.Char.ID, 8, int16(equipSlot))
	if srcEquip == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	targetSlot := int16(invenSlot)
	if h.findItemAt(player.Char.ID, 9, targetSlot) != nil {
		if freeSlot, found := h.findFreeInvenSlot(player.Char.ID); found {
			targetSlot = freeSlot
		} else {
			ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
			ack.WriteUint16(0x0282)
			_ = client.SendPacket(ack)
			return
		}
	}

	_ = h.charRepo.SetItemStorage(srcEquip.ItemKey, 9, targetSlot)
	srcEquip.StorageType = 9
	srcEquip.Storage = targetSlot

	h.sendCellChange(client, 9, targetSlot, srcEquip, 0xFFFF)
	h.sendCellChange(client, 8, int16(equipSlot), nil, 0xFFFF)

	h.syncPlayerEquipsAndStats(player)
	h.broadcastUnequipAppearance(player, equipSlot)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
	ack.WriteUint16(0x0281) // Éxito oficial en Zone.exe (0x0281)
	_ = client.SendPacket(ack)
}

// handleItemDrop implementa sp_NC_ITEM_DROP_REQ (0x3007) y gds_NC_ITEMDB_DROPALLSUC_ACK (Zone.exe VA 0x52AB80):
//   - Lee slot (uint8), lot (uint8), loc (SHINE_XY_TYPE: x int32, y int32)
//   - Elimina o reduce el ítem en Box 9, envía 0x3001, spawnea el ítem en el suelo del mapa (0x1C0A)
//   - Envía PROTO_NC_ITEM_DROP_ACK (0x3008) con err = 0x0301
func (h *Handler) handleItemDrop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	dropLot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	dropX, err := pkt.Reader.ReadInt32()
	if err != nil {
		dropX = int32(player.Char.XPos)
	}
	dropY, err := pkt.Reader.ReadInt32()
	if err != nil {
		dropY = int32(player.Char.YPos)
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(slot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_ACK)
		ack.WriteUint16(0x0304) // Drop fail en Zone.exe
		_ = client.SendPacket(ack)
		return
	}

	curLot := int(srcItem.Flags)
	if curLot <= 0 {
		curLot = 1
	}
	actualDrop := int(dropLot)
	if actualDrop <= 0 || actualDrop >= curLot {
		actualDrop = curLot
		_ = h.charRepo.DeleteItem(srcItem.ItemKey)
		h.sendCellChange(client, 9, int16(slot), nil, 0xFFFF)
	} else {
		srcItem.Flags = int32(curLot - actualDrop)
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
		h.sendCellChange(client, 9, int16(slot), srcItem, 0xFFFF)
	}

	player.CurrentMap.SpawnDroppedItem(srcItem.ItemID, uint32(actualDrop), int(srcItem.Flags), dropX, dropY, player.MapObjectID, player.Char.ID)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_ACK)
	ack.WriteUint16(0x0301) // Éxito oficial en Zone.exe (0x0301)
	_ = client.SendPacket(ack)
}

// handleItemPick implementa sp_NC_ITEM_PICK_REQ (0x3009) y so_ply_itempickresult (Zone.exe VA 0x5308D0):
//   - Lee itemhandle (uint16)
//   - Consume el ítem del mapa (enviando 0x1C0E PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
//   - Inserta en Box 9 del jugador, envía 0x3001 PROTO_NC_ITEM_CELLCHANGE_CMD
//   - Envía PROTO_NC_ITEM_PICK_ACK (0x300A, 10 bytes): itemid(2), lot(4), error(2 = 0x0341), itemhandle(2)
func (h *Handler) handleItemPick(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil {
		return
	}

	itemHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	freeSlot, hasFree := h.findFreeInvenSlot(player.Char.ID)
	if !hasFree {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
		ack.WriteUint16(0xFFFF)
		ack.WriteUint32(0)
		ack.WriteUint16(0x0344) // Inventory full
		ack.WriteUint16(itemHandle)
		_ = client.SendPacket(ack)
		return
	}

	dropped, found := player.CurrentMap.ConsumeDroppedItem(itemHandle)
	if !found || dropped == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
		ack.WriteUint16(0xFFFF)
		ack.WriteUint32(0)
		ack.WriteUint16(0x0342) // Item not found on field
		ack.WriteUint16(itemHandle)
		_ = client.SendPacket(ack)
		return
	}

	lot := int(dropped.Lot)
	if lot <= 0 {
		lot = 1
	}

	newItem := &database.Item{StorageType: 9, Owner: player.Char.ID, Storage: int16(freeSlot), ItemID: dropped.ItemID, Flags: int32(lot)}
	err = h.charRepo.CreateItem(newItem)
	if err != nil {
		return
	}


	h.sendCellChange(client, 9, freeSlot, newItem, 0xFFFF)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
	ack.WriteUint16(dropped.ItemID) // itemid (2B)
	ack.WriteUint32(dropped.Lot)    // lot (4B)
	ack.WriteUint16(0x0341)         // error = 0x0341 (éxito oficial en so_ply_itempickresult)
	ack.WriteUint16(itemHandle)     // itemhandle (2B)
	_ = client.SendPacket(ack)
}

// handleItemUse implementa sp_NC_ITEM_USE_REQ (0x3015), so_NC_ITEMDB_USE_ACK (VA 0x53FB50) y so_ply_ItemEffect (VA 0x59C730):
//   - Lee slot (uint8) e invenType (uint8)
//   - Aplica los efectos de ItemUseEffect.shn (0 = HP, 1 = SP, 2 = Full HP/SP) y envía 0x240E / 0x240F
//   - Consume 1 unidad del stack (o destruye el ítem si Lot == 1) y envía 0x3001 PROTO_NC_ITEM_CELLCHANGE_CMD
//   - Transmite PROTO_NC_ITEM_SOMEONEUSE_CMD (0x3029, 12B) al mapa y envía PROTO_NC_ITEM_USE_ACK (0x3016, 5B) con error = 0x0700
func (h *Handler) handleItemUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	invenType, err := pkt.Reader.ReadByte()
	if err != nil {
		invenType = 9
	}
	if invenType == 0 {
		invenType = 9
	}

	srcItem := h.findItemAt(player.Char.ID, invenType, int16(slot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACK)
		ack.WriteUint16(0x0701) // Use fail
		ack.WriteUint16(0xFFFF)
		ack.WriteByte(invenType)
		_ = client.SendPacket(ack)
		return
	}

	usedItemID := srcItem.ItemID
	usedItemName := ""
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(usedItemID); ok && info != nil {
			usedItemName = info.InxName
		}
	}

	// Aplicar efectos desde ItemUseEffect.shn (so_ply_ItemEffect VA 0x59C730)
	hpChanged := false
	spChanged := false
	if h.itemStore != nil {
		if eff, hasEff := h.itemStore.GetItemUseEffect(usedItemName); hasEff && eff != nil {
			applyEff := func(effectType uint32, v uint16) {
				val := int(v)
				if val == 0 {
					return
				}
				switch effectType {
				case 0: // Recover HP
					player.Char.CurHP += val
					if player.Char.CurHP > player.MaxHP {
						player.Char.CurHP = player.MaxHP
					}
					hpChanged = true
				case 1: // Recover SP
					player.Char.CurSP += val
					if player.Char.CurSP > player.MaxSP {
						player.Char.CurSP = player.MaxSP
					}
					spChanged = true
				case 2: // Recover Full HP & SP (99999)
					player.Char.CurHP = player.MaxHP
					player.Char.CurSP = player.MaxSP
					hpChanged = true
					spChanged = true
				}
			}
			applyEff(eff.UseEffectA, eff.UseValueA)
			applyEff(eff.UseEffectB, eff.UseValueB)
			applyEff(eff.UseEffectC, eff.UseValueC)
		}
	}

	if hpChanged {
		hpPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpPkt.WriteUint32(uint32(player.Char.CurHP))
		hpPkt.WriteUint16(0) // hpchangeorder
		_ = client.SendPacket(hpPkt)
	}
	if spChanged {
		spPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
		spPkt.WriteUint32(uint32(player.Char.CurSP))
		_ = client.SendPacket(spPkt)
	}
	if hpChanged || spChanged {
		_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)
	}

	// Consumir 1 unidad del ítem
	curLot := srcItem.Flags
	if curLot <= 1 {
		_ = h.charRepo.DeleteItem(srcItem.ItemKey)
		h.sendCellChange(client, invenType, int16(slot), nil, 0xFFFF)
	} else {
		srcItem.Flags = curLot - 1
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
		h.sendCellChange(client, invenType, int16(slot), srcItem, 0xFFFF)
	}

	// Broadcast PROTO_NC_ITEM_SOMEONEUSE_CMD (0x3029, 12B)
	if player.CurrentMap != nil {
		someUse := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SOMEONEUSE_CMD)
		someUse.WriteUint16(player.MapObjectID)
		someUse.WriteInt32(int32(player.Char.XPos))
		someUse.WriteInt32(int32(player.Char.YPos))
		someUse.WriteUint16(usedItemID)
		player.CurrentMap.Broadcast(someUse, 0)
	}

	// Enviar PROTO_NC_ITEM_USE_ACK (0x3016, 5B)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACK)
	ack.WriteUint16(0x0700)   // error = 0x0700 (éxito oficial en so_ply_ItemUseResultPacket)
	ack.WriteUint16(usedItemID)
	ack.WriteByte(invenType)
	_ = client.SendPacket(ack)
}

// handleAutoArrangeInven implementa sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ (0x304A, Zone.exe VA 0x51C360):
//   - Ordena los ítems de Box 9 por la categoría oficial ItemSort.shn y los criterios de ItemInfo.shn
//   - Compacta las posiciones desde el slot 0 y envía el ACK oficial 0x30A8
func (h *Handler) handleAutoArrangeInven(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	items, err := h.charRepo.GetCharacterItems(player.Char.ID)
	if err != nil {
		log.Printf("[Zone][%s] auto-arrange: failed to load inventory for character %d: %v",
			client.RemoteAddr(), player.Char.ID, err)
		return
	}
	inventory := make([]*database.Item, 0, len(items))
	for _, it := range items {
		if it != nil && it.StorageType == 9 {
			inventory = append(inventory, it)
		}
	}
	sortInventoryItems(inventory, h.itemStore)

	itemKeys := make([]int64, len(inventory))
	oldSlots := make([]int16, len(inventory))
	for i, it := range inventory {
		itemKeys[i] = it.ItemKey
		oldSlots[i] = it.Storage
	}
	if err := h.charRepo.ArrangeInventory(player.Char.ID, itemKeys); err != nil {
		log.Printf("[Zone][%s] auto-arrange: failed to persist inventory for character %d: %v",
			client.RemoteAddr(), player.Char.ID, err)
		return
	}

	for _, slot := range oldSlots {
		h.sendCellChange(client, 9, slot, nil, 0xFFFF)
	}
	for slot, it := range inventory {
		it.Storage = int16(slot)
		it.StorageType = 9
		h.sendCellChange(client, 9, int16(slot), it, 0xFFFF)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK)
	ack.WriteUint16(0x30A8) // Código oficial en sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ (VA 0x51C360)
	_ = client.SendPacket(ack)
}

// handleSoulStoneHPBuy implementa sp_NC_SOULSTONE_HP_BUY_REQ (0x5001) y gds_NC_ITEMDB_BUYHPSOULSTONESUC_ACK (VA 0x52D820):
//   - Lee number (uint16)
//   - Descuenta el precio en Cen (10 Cen por piedra en Nivel 1 o según Level), incrementa StonesHP
//   - Envía PROTO_NC_SOULSTONE_HP_BUY_ACK (0x5003), PROTO_NC_ITEM_BUY_ACK (0x3004 = 0x0201) y PROTO_NC_CHAR_CENCHANGE_CMD (0x1033)
func (h *Handler) handleSoulStoneHPBuy(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	buyNum, err := pkt.Reader.ReadUint16()
	if err != nil || buyNum == 0 {
		return
	}

	maxStone := int(player.MaxStonesHP)
	if player.LevelParam.MAXSoulHP > 0 {
		maxStone = player.LevelParam.MAXSoulHP
	}
	if int(player.StonesHP)+int(buyNum) > maxStone {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0742) // Max stones exceeded
		_ = client.SendPacket(fail)
		return
	}

	// Precio oficial por piedra desde ClassParam (LevelParam.PriceHPStone)
	totalCost := int64(player.LevelParam.PriceHPStone) * int64(buyNum)
	if player.Char.TCharacterState.Money < totalCost {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743) // Insufficient money
		_ = client.SendPacket(fail)
		return
	}
	if okSub, err := h.charRepo.SubCharMoney(player.Char.ID, totalCost); err != nil || !okSub {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743)
		_ = client.SendPacket(fail)
		return
	}

	player.Char.TCharacterState.Money -= totalCost
	player.StonesHP += int16(buyNum)
	player.Char.TCharacterState.HPStones = player.StonesHP

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	// 1. PROTO_NC_SOULSTONE_HP_BUY_ACK (0x5003, 2B: totalnumber)
	buyAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_ACK)
	buyAck.WriteUint16(uint16(player.StonesHP))
	_ = client.SendPacket(buyAck)

	// 2. PROTO_NC_ITEM_BUY_ACK (0x3004, 2B: err = 0x0201)
	itemBuyAck := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_ACK)
	itemBuyAck.WriteUint16(0x0201)
	_ = client.SendPacket(itemBuyAck)

	// 3. PROTO_NC_CHAR_CENCHANGE_CMD (0x1033, 8B: cen uint64)
	cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
	cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
	_ = client.SendPacket(cenCmd)
}

// handleSoulStoneSPBuy implementa sp_NC_SOULSTONE_SP_BUY_REQ (0x5002) y gds_NC_ITEMDB_BUYSPSOULSTONESUC_ACK (VA 0x52ED00):
func (h *Handler) handleSoulStoneSPBuy(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	buyNum, err := pkt.Reader.ReadUint16()
	if err != nil || buyNum == 0 {
		return
	}

	maxStone := int(player.MaxStonesSP)
	if player.LevelParam.MAXSoulSP > 0 {
		maxStone = player.LevelParam.MAXSoulSP
	}
	if int(player.StonesSP)+int(buyNum) > maxStone {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0742) // Max stones exceeded
		_ = client.SendPacket(fail)
		return
	}

	// Precio oficial por piedra desde ClassParam (LevelParam.PriceSPStone)
	totalCost := int64(player.LevelParam.PriceSPStone) * int64(buyNum)
	if player.Char.TCharacterState.Money < totalCost {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743) // Insufficient money
		_ = client.SendPacket(fail)
		return
	}
	if okSub, err := h.charRepo.SubCharMoney(player.Char.ID, totalCost); err != nil || !okSub {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743)
		_ = client.SendPacket(fail)
		return
	}

	player.Char.TCharacterState.Money -= totalCost
	player.StonesSP += int16(buyNum)
	player.Char.TCharacterState.SPStones = player.StonesSP

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	// 1. PROTO_NC_SOULSTONE_SP_BUY_ACK (0x5004, 2B: totalnumber)
	buyAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_BUY_ACK)
	buyAck.WriteUint16(uint16(player.StonesSP))
	_ = client.SendPacket(buyAck)

	// 2. PROTO_NC_ITEM_BUY_ACK (0x3004, 2B: err = 0x0201)
	itemBuyAck := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_ACK)
	itemBuyAck.WriteUint16(0x0201)
	_ = client.SendPacket(itemBuyAck)

	// 3. PROTO_NC_CHAR_CENCHANGE_CMD (0x1033, 8B: cen uint64)
	cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
	cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
	_ = client.SendPacket(cenCmd)
}

// handleSoulStoneHPUse implementa sp_NC_SOULSTONE_HP_USE_REQ (0x5007, VA 0x588A40):
//   - Verifica StonesHP > 0, decrementa StonesHP en 1, recupera LevelParam.SoulHP de HP
//   - Transmite PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD (0x500B, 2B) al mapa
//   - Envía PROTO_NC_SOULSTONE_HP_USESUC_ACK (0x5008, 0B) y PROTO_NC_BAT_HPCHANGE_CMD (0x240E, 6B)
func (h *Handler) handleSoulStoneHPUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	if player.StonesHP <= 0 {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK)
		_ = client.SendPacket(fail)
		return
	}

	player.StonesHP--
	player.Char.TCharacterState.HPStones = player.StonesHP

	recoverHP := player.LevelParam.SoulHP
	player.Char.CurHP += recoverHP
	if player.Char.CurHP > player.MaxHP {
		player.Char.CurHP = player.MaxHP
	}

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	if player.CurrentMap != nil {
		someCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD)
		someCmd.WriteUint16(player.MapObjectID)
		player.CurrentMap.Broadcast(someCmd, 0)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK)
	_ = client.SendPacket(sucAck)

	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0) // hpchangeorder
	_ = client.SendPacket(hpCmd)
}

// handleSoulStoneSPUse implementa sp_NC_SOULSTONE_SP_USE_REQ (0x5009, VA 0x589180):
//   - Verifica StonesSP > 0, decrementa StonesSP en 1, recupera LevelParam.SoulSP de SP
//   - Transmite PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD (0x500C, 2B) al mapa
//   - Envía PROTO_NC_SOULSTONE_SP_USESUC_ACK (0x500A, 0B) y PROTO_NC_BAT_SPCHANGE_CMD (0x240F, 4B)
func (h *Handler) handleSoulStoneSPUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	if player.StonesSP <= 0 {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK)
		_ = client.SendPacket(fail)
		return
	}

	player.StonesSP--
	player.Char.TCharacterState.SPStones = player.StonesSP

	recoverSP := player.LevelParam.SoulSP
	player.Char.CurSP += recoverSP
	if player.Char.CurSP > player.MaxSP {
		player.Char.CurSP = player.MaxSP
	}

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	if player.CurrentMap != nil {
		someCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD)
		someCmd.WriteUint16(player.MapObjectID)
		player.CurrentMap.Broadcast(someCmd, 0)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_USESUC_ACK)
	_ = client.SendPacket(sucAck)

	spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
	spCmd.WriteUint32(uint32(player.Char.CurSP))
	_ = client.SendPacket(spCmd)
}


const slotWeapon byte = 12 // ItemEquipEnum::ITEMEQUIP_RIGHTHAND (Slot_Weapon)
