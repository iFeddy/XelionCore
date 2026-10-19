package zone

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler processes incoming packets for the Zone Server.
type Handler struct {
	mu         sync.RWMutex
	charRepo   database.CharacterRepository
	transfers  *database.TransferManager
	classStore *classparam.Store
	fieldStore *field.Store
	players    map[*network.Client]*ZonePlayer
	roumenMap  *ZoneMap
}

// NewHandler initializes a new Zone packet handler.
func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.Store) *Handler {
	return &Handler{
		charRepo:   charRepo,
		transfers:  transfers,
		classStore: classStore,
		fieldStore: fieldStore,
		players:    make(map[*network.Client]*ZonePlayer),
		roumenMap:  NewZoneMap(0, "Roumen", "RouN"),
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
		h.roumenMap.RemovePlayer(player)
		_ = h.charRepo.UpdateCharacter(player.Char)
	}
}

func (h *Handler) getPlayer(client *network.Client) (*ZonePlayer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	p, ok := h.players[client]
	return p, ok
}

func (h *Handler) handleTransferKey(client *network.Client, pkt *protocol.Packet) {
	randomID, err := pkt.Reader.ReadUint16()
	if err != nil {
		log.Printf("[Zone][%s] Error reading randomID: %v", client.RemoteAddr(), err)
		return
	}

	charName, err := pkt.Reader.ReadString(16)
	if err != nil {
		log.Printf("[Zone][%s] Error reading character name: %v", client.RemoteAddr(), err)
		return
	}
	charName = strings.TrimSpace(charName)

	log.Printf("[Zone][%s] TransferKey received for character '%s' (RandomID: %d)",
		client.RemoteAddr(), charName, randomID)

	transfer, ok := h.transfers.ConsumeZoneTransfer(charName, randomID, pkt.Raw)
	if !ok {
		log.Printf("[Zone][%s] Transfer verification failed for '%s' (RandomID: %d)",
			client.RemoteAddr(), charName, randomID)
		errPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(4, 0) // RequestedCharacterIDNotMatching
		_ = client.SendPacket(errPkt)
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
		return
	}

	player := NewZonePlayer(client, char, h.classStore)
	h.roumenMap.AssignObjectID(player)
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

	// 5. PROTO_NC_CHAR_CLIENT_SKILL_CMD (0x103D, wire len = 12, payload = 10)
	// struct PROTO_NC_CHAR_CLIENT_SKILL_CMD {
	//     unsigned char  restempow; // 1 byte
	//     unsigned char  PartMark;  // 1 byte (3 = all parts/mark)
	//     unsigned short nMaxNum;   // 2 bytes
	//     unsigned long  chrregnum; // 4 bytes
	//     unsigned short number;    // 2 bytes
	// }
	activeSkills := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SKILL_CMD)
	activeSkills.WriteByte(0)                  // restempow
	activeSkills.WriteByte(3)                  // PartMark (3 = all parts)
	activeSkills.WriteUint16(0)                // nMaxNum
	activeSkills.WriteUint32(uint32(p.Char.ID)) // chrregnum
	activeSkills.WriteUint16(0)                // number
	_ = p.Client.SendPacket(activeSkills)

	// 6. PROTO_NC_CHAR_CLIENT_PASSIVE_CMD (0x103E, wire len = 4, payload = 2)
	// struct PROTO_NC_CHAR_CLIENT_PASSIVE_CMD {
	//     unsigned short number; // 2 bytes
	// }
	passiveSkills := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)
	passiveSkills.WriteUint16(0) // number = 0
	_ = p.Client.SendPacket(passiveSkills)

	// 7-10. PROTO_NC_CHAR_CLIENT_ITEM_CMD (0x1047, wire len = 5, payload = 3 each)
	// struct PROTO_NC_CHAR_CLIENT_ITEM_CMD {
	//     unsigned char numofitem; // 1 byte
	//     unsigned char box;       // 1 byte
	//     unsigned char flag;      // 1 byte (1)
	// }
	// The 4 official boxes sent on login: 8 (Equip), 9 (Inventory), 12 (House), 15 (Emotions)
	boxes := []uint8{8, 9, 12, 15}
	for _, box := range boxes {
		itemPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD)
		itemPkt.WriteByte(0)   // numofitem = 0
		itemPkt.WriteByte(box) // box
		itemPkt.WriteByte(1)   // flag = 1
		_ = p.Client.SendPacket(itemPkt)
	}

	// 11. PROTO_NC_CHAR_CLIENT_GAME_CMD (0x1048, wire len = 6, payload = 4)
	// struct PROTO_NC_CHAR_CLIENT_GAME_CMD {
	//     unsigned long party_handle; // 4 bytes (0 = no party)
	// }
	gamePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD)
	gamePkt.WriteUint32(0) // party_handle = 0
	_ = p.Client.SendPacket(gamePkt)

	// 12. PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD (0x1049, wire len = 8, payload = 6)
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

	// 13. PROTO_NC_CHAR_MAPLOGIN_ACK (0x1802, wire len = 244, payload = 242)
	extraPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_MAPLOGIN_ACK)
	p.WriteDetailedInfoExtra(extraPkt, false)
	_ = p.Client.SendPacket(extraPkt)
}

func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	log.Printf("[Zone][%s] ClientReady received from '%s'. Spawning into Roumen...",
		client.RemoteAddr(), player.Char.Name)

	// 1. Usable Points
	pts := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_POINTS_ACK)
	pts.WriteByte(0)
	_ = client.SendPacket(pts)

	// 2. Add player to map and broadcast spawn
	h.roumenMap.AddPlayer(player)
}

func (h *Handler) handleWalk(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	h.roumenMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), true)
}

func (h *Handler) handleRun(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	h.roumenMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), false)
}

func (h *Handler) handleStop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	newX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	newY, _ := pkt.Reader.ReadInt32()

	h.roumenMap.StopPlayer(player, int(newX), int(newY))
}

func (h *Handler) handleJump(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	h.roumenMap.JumpPlayer(player)
}

func (h *Handler) handleChatNormal(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
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

	log.Printf("[Zone][Chat][%s] %s: %s", h.roumenMap.Name, player.Char.Name, text)

	color := byte(0x2A) // Normal chat color
	if player.Client.Admin > 0 {
		color = 0x03 // GM gold color
	}

	h.roumenMap.NormalChat(player, text, color)
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
