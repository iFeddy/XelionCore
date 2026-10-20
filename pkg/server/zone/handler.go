package zone

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"
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
	case protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ:
		h.handleRewardInvenOpen(client, pkt)
	case protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ:
		h.handleChargedInvenOpen(client, pkt)
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
	}
}
