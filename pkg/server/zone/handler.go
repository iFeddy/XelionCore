package zone

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
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
	classStore  *classparam.Store
	fieldStore  *field.Store
	chkEngine   *shn.ChecksumEngine
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
	case protocol.PROTO_NC_ACT_STOP_REQ:
		h.handleStop(client, pkt)
	case protocol.PROTO_NC_ACT_JUMP_CMD:
		h.handleJump(client, pkt)
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
		_ = h.charRepo.UpdateCharacter(player.Char)

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

func (h *Handler) sendInGameChunk(p *ZonePlayer) {
	// 1. PROTO_NC_CHAR_CLIENT_BASE_CMD (0x1038, wire len = 107, payload = 105)
	infoPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD)
	p.WriteDetailedInfo(infoPkt)
	_ = p.Client.SendPacket(infoPkt)

	// 2. PROTO_NC_CHAR_CLIENT_SHAPE_CMD (0x1039, wire len = 6, payload = 4)
	lookPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD)
	p.WriteLook(lookPkt)
	_ = p.Client.SendPacket(lookPkt)

	// 3. PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD (0x103A, wire len = 8, payload = 6)
	// struct PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD {
	//     unsigned long chrregnum;       // 4 bytes
	//     unsigned char bNeedClear;      // 1 byte
	//     unsigned char nNumOfDoingQuest; // 1 byte
	// }
	questsBusy := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)
	questsBusy.WriteUint32(uint32(p.Char.ID))
	questsBusy.WriteByte(1) // bNeedClear = 1
	questsBusy.WriteByte(0) // nNumOfDoingQuest = 0
	_ = p.Client.SendPacket(questsBusy)

	// 4. PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD (0x103B, wire len = 14, payload = 12)
	// struct PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD {
	//     unsigned long  chrregnum;           // 4 bytes
	//     unsigned short nTotalDoneQuest;     // 2 bytes
	//     unsigned short nTotalDoneQuestSize; // 2 bytes
	//     unsigned short nDoneQuestCount;     // 2 bytes
	//     unsigned short nIndex;              // 2 bytes
	// }
	questsDone := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD)
	questsDone.WriteUint32(uint32(p.Char.ID))
	questsDone.WriteUint16(0) // nTotalDoneQuest
	questsDone.WriteUint16(0) // nTotalDoneQuestSize
	questsDone.WriteUint16(0) // nDoneQuestCount
	questsDone.WriteUint16(0) // nIndex
	_ = p.Client.SendPacket(questsDone)

	// 5. PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD (0x10CE, wire len = 8, payload = 6)
	// struct PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD {
	//     unsigned long  chrregnum;       // 4 bytes
	//     unsigned short nNumOfReadQuest; // 2 bytes
	// }
	questsRead := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD)
	questsRead.WriteUint32(uint32(p.Char.ID))
	questsRead.WriteUint16(0) // nNumOfReadQuest = 0
	_ = p.Client.SendPacket(questsRead)

	// 6. PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD (0x10D7, wire len = 8, payload = 6)
	// struct PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD {
	//     unsigned long  chrregnum;         // 4 bytes
	//     unsigned short nNumOfRepeatQuest; // 2 bytes
	// }
	questsRepeat := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD)
	questsRepeat.WriteUint32(uint32(p.Char.ID))
	questsRepeat.WriteUint16(0) // nNumOfRepeatQuest = 0
	_ = p.Client.SendPacket(questsRepeat)

	// 7. PROTO_NC_CHAR_CLIENT_PASSIVE_CMD (0x103E, wire len = 4, payload = 2)
	// struct PROTO_NC_CHAR_CLIENT_PASSIVE_CMD {
	//     unsigned short number; // 2 bytes
	// }
	passiveSkills := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)
	passiveSkills.WriteUint16(0) // number = 0
	_ = p.Client.SendPacket(passiveSkills)

	// 8. PROTO_NC_CHAR_CLIENT_SKILL_CMD (0x103D):
	// En el binario oficial Zone.exe (so_charskilldata), si el personaje no tiene
	// skills aprendidas (nMaxNum == 0), el servidor oficial NO envía 0x103D.
	// Si tuviera skills activas, se enviaría aquí.

	// 9. PROTO_NC_CHAR_CLIENT_GAME_CMD (0x1048, wire len = 6, payload = 4)
	// struct PROTO_NC_CHAR_CLIENT_GAME_CMD {
	//     unsigned short partyno; // 2 bytes (0xFFFF = sin party)
	//     unsigned short nRaidNo; // 2 bytes (0xFFFF = sin raid)
	// }
	gamePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD)
	gamePkt.WriteUint32(0xFFFFFFFF) // partyno = 0xFFFF, nRaidNo = 0xFFFF (Sin party ni raid)
	_ = p.Client.SendPacket(gamePkt)

	// 10-13. PROTO_NC_CHAR_CLIENT_ITEM_CMD (0x1047, wire len = 5, payload = 3 each)
	// struct PROTO_NC_CHAR_CLIENT_ITEM_CMD {
	//     unsigned char numofitem; // 1 byte
	//     unsigned char box;       // 1 byte
	//     unsigned char flag;      // 1 byte (1)
	//     PROTO_ITEMPACKET_INFORM ItemArray[numofitem];
	// }
	// The 4 official boxes sent on login: 8 (Equip), 9 (Inventory), 12 (House), 15 (Emotions)
	charItems, _ := h.charRepo.GetCharacterItems(p.Char.ID)
	var invenItems []*database.Item
	for _, it := range charItems {
		if it.StorageType == 9 {
			invenItems = append(invenItems, it)
		}
	}

	boxes := []uint8{8, 9, 12, 15}
	for _, box := range boxes {
		itemPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD)
		if box == 9 && len(invenItems) > 0 {
			itemPkt.WriteByte(byte(len(invenItems))) // numofitem
			itemPkt.WriteByte(9)                     // box 9
			itemPkt.WriteByte(1)                     // flag = 1
			for _, it := range invenItems {
				itemPkt.WriteByte(8)                // datasize: 2 (location) + 2 (itemid) + 4 (amount) = 8
				itemPkt.WriteByte(byte(it.Storage)) // Pos (slot 0..23)
				itemPkt.WriteByte(0x24)             // Type (0x24 = normal inventory item)
				itemPkt.WriteUint16(it.ItemID)      // ItemID
				amt := uint32(it.Flags)
				if amt == 0 {
					amt = 1
				}
				itemPkt.WriteUint32(amt)            // Amount
			}
		} else {
			itemPkt.WriteByte(0)   // numofitem = 0
			itemPkt.WriteByte(box) // box
			itemPkt.WriteByte(1)   // flag = 1
		}
		_ = p.Client.SendPacket(itemPkt)
	}

	// 14. PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD (0x1049, wire len = 8, payload = 6)
	// struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD {
	//     unsigned char  CurrentTitle;        // 1 byte
	//     unsigned char  CurrentTitleElement; // 1 byte
	//     unsigned short CurrentTitleMobID;   // 2 bytes
	//     unsigned short NumOfTitle;          // 2 bytes
	// }
	titles := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD)
	titles.WriteByte(0)   // CurrentTitle = 0
	titles.WriteByte(0)   // CurrentTitleElement = 0
	titles.WriteUint16(0) // CurrentTitleMobID = 0
	titles.WriteUint16(0) // NumOfTitle = 0
	_ = p.Client.SendPacket(titles)

	// 15. PROTO_NC_CHAR_MAPLOGIN_ACK (0x1802, wire len = 244, payload = 242)
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
