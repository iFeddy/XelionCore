package zone

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler processes incoming packets for the Zone Server.
type Handler struct {
	mu        sync.RWMutex
	charRepo  database.CharacterRepository
	transfers *database.TransferManager
	players   map[*network.Client]*ZonePlayer
	roumenMap *ZoneMap
}

// NewHandler initializes a new Zone packet handler.
func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager) *Handler {
	return &Handler{
		charRepo:  charRepo,
		transfers: transfers,
		players:   make(map[*network.Client]*ZonePlayer),
		roumenMap: NewZoneMap(0, "Roumen", "Rou"),
	}
}

// HandlePacket routes packets to their specific handlers.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.CH6Type_TransferKey:
		h.handleTransferKey(client, pkt)
	case protocol.CH6Type_ClientReady:
		h.handleClientReady(client, pkt)
	case protocol.CH8Type_Walk:
		h.handleWalk(client, pkt)
	case protocol.CH8Type_Run:
		h.handleRun(client, pkt)
	case protocol.CH8Type_Stop:
		h.handleStop(client, pkt)
	case protocol.CH8Type_Jump:
		h.handleJump(client, pkt)
	case protocol.CH8Type_ChatNormal:
		h.handleChatNormal(client, pkt)
	case protocol.CH8Type_BeginRest:
		h.handleBeginRest(client, pkt)
	case protocol.CH8Type_EndRest:
		h.handleEndRest(client, pkt)
	case protocol.CH2Type_Pong, protocol.CH2Type_Unk1:
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
		errPkt := protocol.NewPacketWriter(protocol.SH4Type_ConnectError)
		errPkt.WriteUint16(4) // RequestedCharacterIDNotMatching
		_ = client.SendPacket(errPkt)
		return
	}

	char, err := h.charRepo.GetCharacterByName(transfer.CharName)
	if err != nil {
		char, err = h.charRepo.GetCharacterByID(transfer.CharID)
	}

	if err != nil || char == nil {
		log.Printf("[Zone][%s] Character '%s' not found in database: %v", client.RemoteAddr(), transfer.CharName, err)
		errPkt := protocol.NewPacketWriter(protocol.SH4Type_ConnectError)
		errPkt.WriteUint16(4)
		_ = client.SendPacket(errPkt)
		return
	}

	player := NewZonePlayer(client, char)
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
	// 1. Character Info
	infoPkt := protocol.NewPacketWriter(protocol.SH4Type_CharacterInfo)
	p.WriteDetailedInfo(infoPkt)
	_ = p.Client.SendPacket(infoPkt)

	// 2. Character Look
	lookPkt := protocol.NewPacketWriter(protocol.SH4Type_CharacterLook)
	p.WriteLook(lookPkt)
	_ = p.Client.SendPacket(lookPkt)

	// 3. Busy Quests
	questsBusy := protocol.NewPacketWriter(protocol.SH4Type_CharacterQuestsBusy)
	questsBusy.WriteInt32(int32(p.Char.ID))
	questsBusy.WriteByte(1) // enabled list
	questsBusy.WriteByte(0) // count
	_ = p.Client.SendPacket(questsBusy)

	// 4. Done Quests
	questsDone := protocol.NewPacketWriter(protocol.SH4Type_CharacterQuestsDone)
	questsDone.WriteInt32(int32(p.Char.ID))
	questsDone.WriteInt16(0) // count
	_ = p.Client.SendPacket(questsDone)

	// 5. Active Skills
	activeSkills := protocol.NewPacketWriter(protocol.SH4Type_CharacterActiveSkillList)
	activeSkills.WriteByte(0)
	activeSkills.WriteInt32(int32(p.Char.ID))
	activeSkills.WriteUint16(0)
	_ = p.Client.SendPacket(activeSkills)

	// 6. Passive Skills
	passiveSkills := protocol.NewPacketWriter(protocol.SH4Type_CharacterPassiveSkillList)
	passiveSkills.WriteUint16(0)
	_ = p.Client.SendPacket(passiveSkills)

	// 7. Equipped Items (inv = 0x08)
	equipped := protocol.NewPacketWriter(protocol.SH4Type_CharacterItemList)
	equipped.WriteByte(0)    // item count
	equipped.WriteByte(0x08) // inventory number
	equipped.WriteByte(115)  // UNK
	_ = p.Client.SendPacket(equipped)

	// 8. Inventory Items (inv = 0x09)
	inventory := protocol.NewPacketWriter(protocol.SH4Type_CharacterItemList)
	inventory.WriteByte(0)
	inventory.WriteByte(0x09)
	inventory.WriteByte(115)
	_ = p.Client.SendPacket(inventory)

	// 9. House Items (inv = 0x0C)
	house := protocol.NewPacketWriter(protocol.SH4Type_CharacterItemList)
	house.WriteByte(0)
	house.WriteByte(0x0C)
	house.WriteByte(0x35)
	_ = p.Client.SendPacket(house)

	// 10. Premium Emotions (inv = 0x0F)
	emotions := protocol.NewPacketWriter(protocol.SH4Type_CharacterItemList)
	emotions.WriteByte(0)
	emotions.WriteByte(0x0F)
	emotions.WriteByte(0xCB)
	_ = p.Client.SendPacket(emotions)

	// 11. Premium Item List
	premium := protocol.NewPacketWriter(protocol.SH4Type_CharacterItemList)
	premium.WriteUint16(0)
	_ = p.Client.SendPacket(premium)

	// 12. Titles
	titles := protocol.NewPacketWriter(protocol.SH4Type_CharacterTitles)
	titles.WriteInt32(0)
	titles.WriteInt16(0)
	_ = p.Client.SendPacket(titles)

	// 13. Chunk End
	endPkt := protocol.NewPacketWriter(protocol.SH4Type_CharacterInfoEnd)
	endPkt.WriteUint16(0xFFFF)
	_ = p.Client.SendPacket(endPkt)

	// 14. Detailed Info Extra
	extraPkt := protocol.NewPacketWriter(protocol.SH6Type_DetailedCharacterInfo)
	p.WriteDetailedInfoExtra(extraPkt, false)
	_ = p.Client.SendPacket(extraPkt)

	// 15. Premium Item List (SH12Type_SendPremiumItemList)
	premiumList := protocol.NewPacketWriter(protocol.SH12Type_SendPremiumItemList)
	premiumList.WriteUint16(0x1041)
	premiumList.WriteByte(1)  // UNK
	premiumList.WriteUint16(0) // Item count = 0
	_ = p.Client.SendPacket(premiumList)

	// 16. Reward List (SH12Type_SendRewardList)
	rewardList := protocol.NewPacketWriter(protocol.SH12Type_SendRewardList)
	rewardList.WriteByte(0)  // Item count = 0
	rewardList.WriteByte(90) // UNK
	_ = p.Client.SendPacket(rewardList)
}

func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	log.Printf("[Zone][%s] ClientReady received from '%s'. Spawning into Roumen...",
		client.RemoteAddr(), player.Char.Name)

	// 1. Usable Points
	pts := protocol.NewPacketWriter(protocol.SH4Type_CharacterPoints)
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
	pw := protocol.NewPacketWriter(protocol.SH8Type_BeginRest)
	pw.WriteUint16(0x0A81) // OK to rest
	_ = client.SendPacket(pw)
}

func (h *Handler) handleEndRest(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	player.State = 1 // Normal
	pw := protocol.NewPacketWriter(protocol.SH8Type_EndRest)
	pw.WriteUint16(0x0A81)
	_ = client.SendPacket(pw)
}
