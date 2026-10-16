package world

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler dispatches incoming packets for the World Server.
type Handler struct {
	charRepo  database.CharacterRepository
	transfers *database.TransferManager
	zoneIP    string
	zonePort  uint16
}

// NewHandler creates a new World packet handler.
func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager, zoneIP string, zonePort uint16) *Handler {
	if zoneIP == "" {
		zoneIP = "127.0.0.1"
	}
	if zonePort == 0 {
		zonePort = 9210
	}
	return &Handler{
		charRepo:  charRepo,
		transfers: transfers,
		zoneIP:    zoneIP,
		zonePort:  zonePort,
	}
}

// HandlePacket dispatches packets based on OpCode.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.CH2Type_Pong:
		// Client heartbeat acknowledge
	case protocol.CH2Type_Unk1:
		h.handleUnk1(client, pkt)
	case protocol.CH3Type_WorldClientKey:
		h.handleWorldClientKey(client, pkt)
	case protocol.CH3Type_BackToWorldSelect:
		h.handleBackToWorldSelect(client, pkt)
	case protocol.CH4Type_CharSelect:
		h.handleCharSelect(client, pkt)
	case protocol.CH5Type_CreateCharacter:
		h.handleCreateCharacter(client, pkt)
	case protocol.CH5Type_DeleteCharacter:
		h.handleDeleteCharacter(client, pkt)
	case protocol.CH22Type_GotIngame:
		h.handleGotIngame(client, pkt)
	case protocol.CH28Type_GetQuickBar:
		h.handleGetQuickBar(client, pkt)
	case protocol.CH28Type_GetQuickBarState:
		h.handleGetQuickBarState(client, pkt)
	case protocol.CH28Type_GetGameSettings:
		h.handleGetGameSettings(client, pkt)
	case protocol.CH28Type_GetClientSettings:
		h.handleGetClientSettings(client, pkt)
	case protocol.CH28Type_GetShortCuts:
		h.handleGetShortCuts(client, pkt)
	case protocol.CH28Type_SaveQuickBar, protocol.CH28Type_SaveQuickBarState,
		protocol.CH28Type_SaveGameSettings, protocol.CH28Type_SaveClientSettings,
		protocol.CH28Type_SaveShortCuts:
		h.handleSaveSettings(client, pkt)
	case protocol.CH31Type_GetUnknown:
		h.handleGetUnknown31(client, pkt)
	default:
		log.Printf("[World][%s] Unhandled OpCode: 0x%04X (Header=%d, Type=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

func (h *Handler) handleUnk1(client *network.Client, pkt *protocol.Packet) {
	now := time.Now()
	resp := protocol.NewPacketWriter(protocol.SH2Type_Unk1)
	resp.WriteByte(byte(now.Hour()))
	resp.WriteByte(byte(now.Minute()))
	resp.WriteByte(byte(now.Second()))
	_ = client.SendPacket(resp)
}

func (h *Handler) handleWorldClientKey(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Raw WorldClientKey (len %d): %x (ASCII: %q)",
		client.RemoteAddr(), len(pkt.Raw), pkt.Raw, string(pkt.Raw))

	var token string
	if pkt.Reader.Remaining() >= 18+32 {
		_ = pkt.Reader.Skip(18)
		token, _ = pkt.Reader.ReadString(32)
		token = strings.TrimSpace(token)
	}

	transfer, ok := h.transfers.ConsumeTransfer(token, pkt.Raw)
	if !ok {
		log.Printf("[World][%s] Transfer token '%s' not found or expired", client.RemoteAddr(), token)
		h.sendError(client, protocol.ErrInvalidCredentials)
		client.Close()
		return
	}

	client.AccountID = transfer.AccountID
	client.Username = transfer.Username
	client.Admin = transfer.Admin
	client.IsAuthenticated = true

	log.Printf("[World][%s] Authenticated user '%s' (AccountID: %d)",
		client.RemoteAddr(), client.Username, client.AccountID)

	// Send Character List
	h.sendCharacterList(client)
}

func (h *Handler) sendCharacterList(client *network.Client) {
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		log.Printf("[World][%s] Error loading characters for account %d: %v",
			client.RemoteAddr(), client.AccountID, err)
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	pw := protocol.NewPacketWriter(protocol.SH3Type_CharacterList)
	if client.RandomID == 0 {
		client.RandomID = uint16(rand.Intn(65534) + 1)
	}
	pw.WriteUint16(client.RandomID)
	pw.WriteByte(byte(len(chars)))

	for _, ch := range chars {
		WriteBasicCharInfo(ch, pw)
	}

	log.Printf("[World][%s] Sent %d character(s) to '%s'",
		client.RemoteAddr(), len(chars), client.Username)

	_ = client.SendPacket(pw)
}

func (h *Handler) handleCreateCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[World][%s] Unauthenticated CreateCharacter request", client.RemoteAddr())
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	name, err := pkt.Reader.ReadString(20)
	if err != nil {
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}
	name = strings.TrimSpace(name)

	jobGender, _ := pkt.Reader.ReadByte()
	hair, _ := pkt.Reader.ReadByte()
	color, _ := pkt.Reader.ReadByte()
	style, _ := pkt.Reader.ReadByte()

	if slot > 5 {
		log.Printf("[World][%s] Invalid slot %d requested", client.RemoteAddr(), slot)
		h.sendCharCreationError(client, protocol.CreateCharError_ErrorInMaxSlot)
		return
	}

	if len(name) < 3 || len(name) > 16 {
		log.Printf("[World][%s] Invalid character name length '%s'", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	// Check if name is taken
	if existing, _ := h.charRepo.GetCharacterByName(name); existing != nil {
		log.Printf("[World][%s] Character name '%s' already taken", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_NameTaken)
		return
	}

	isMale := ((jobGender >> 7) & 0x01) == 1
	classID := byte((jobGender >> 2) & 0x1F)
	job := protocol.Job(classID)

	switch job {
	case protocol.Job_Fighter, protocol.Job_Cleric, protocol.Job_Archer, protocol.Job_Mage, protocol.Job_Trickster:
		// Valid job
	default:
		log.Printf("[World][%s] Invalid job ID %d requested", client.RemoteAddr(), classID)
		h.sendCharCreationError(client, protocol.CreateCharError_WrongClass)
		return
	}

	wSlot, wItem := GetBeginnerWeapon(job)

	newChar := &database.Character{
		AccountID: client.AccountID,
		Name:      name,
		Slot:      slot,
		Job:       classID,
		Level:     1,
		Male:      isMale,
		Hair:      hair,
		HairColor: color,
		Face:      style,
		TownName:  "Roumen",
		XPos:      7636,
		YPos:      4610,
		CurHP:     100,
		CurSP:     100,
		Equips:    map[byte]uint16{wSlot: wItem},
	}

	if err := h.charRepo.CreateCharacter(newChar); err != nil {
		log.Printf("[World][%s] Failed to save character '%s': %v", client.RemoteAddr(), name, err)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	log.Printf("[World][%s] Created character '%s' (Job: %d, Slot: %d) for user '%s'",
		client.RemoteAddr(), name, classID, slot, client.Username)

	// Send Character Creation OK response
	resp := protocol.NewPacketWriter(protocol.SH5Type_CharCreationOK)
	resp.WriteByte(1)
	WriteBasicCharInfo(newChar, resp)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleDeleteCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		return
	}

	if err := h.charRepo.DeleteCharacter(client.AccountID, slot); err != nil {
		log.Printf("[World][%s] Error deleting character in slot %d: %v", client.RemoteAddr(), slot, err)
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	log.Printf("[World][%s] Deleted character in slot %d for user '%s'", client.RemoteAddr(), slot, client.Username)

	resp := protocol.NewPacketWriter(protocol.SH5Type_CharDeleteOK)
	resp.WriteByte(slot)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleBackToWorldSelect(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH3Type_BackToWorldListFromChar)
	resp.WriteBytes([]byte{0x58, 0x1E})
	resp.Fill(32, 0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleCharSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		log.Printf("[World][%s] Invalid character slot: %v", client.RemoteAddr(), err)
		return
	}

	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	var selected *database.Character
	for _, ch := range chars {
		if ch.Slot == slot {
			selected = ch
			break
		}
	}

	if selected == nil {
		log.Printf("[World][%s] Character in slot %d not found for account %d", client.RemoteAddr(), slot, client.AccountID)
		return
	}

	randomID := client.RandomID
	if randomID == 0 {
		randomID = uint16(rand.Intn(65534) + 1)
		client.RandomID = randomID
	}
	clientIP := client.RemoteIP()

	h.transfers.CreateZoneTransfer(client.AccountID, client.Username, selected.Name, selected.ID, randomID, client.Admin, clientIP)

	log.Printf("[World][%s] Character '%s' (ID: %d) selected. Transfer registered with RandomID: %d. Dispatching SH4Type_ServerIP -> %s:%d",
		client.RemoteAddr(), selected.Name, selected.ID, randomID, h.zoneIP, h.zonePort)

	// Send SH4Type_ServerIP (Header 4, Type 3)
	resp := protocol.NewPacketWriter(protocol.SH4Type_ServerIP)
	resp.WriteString(h.zoneIP, 16)
	resp.WriteUint16(h.zonePort)
	_ = client.SendPacket(resp)
}

func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.NewPacketWriter(protocol.SH3Type_Error)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) sendCharCreationError(client *network.Client, errCode protocol.CreateCharError) {
	resp := protocol.NewPacketWriter(protocol.SH5Type_CharCreationError)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGetQuickBar(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH28Type_LoadQuickBar)
	resp.WriteByte(0) // hasData = false
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGetQuickBarState(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH28Type_LoadQuickBarState)
	resp.WriteByte(0)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGetGameSettings(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH28Type_LoadGameSettings)
	resp.WriteByte(0)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGetClientSettings(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH28Type_LoadClientSettings)
	resp.WriteByte(0)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGetShortCuts(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH28Type_LoadShortCuts)
	resp.WriteByte(0)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleSaveSettings(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Received settings save packet 0x%04X (Header=%d, Type=%d)",
		client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
}

func (h *Handler) handleGetUnknown31(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.SH31Type_LoadUnknown)
	resp.WriteInt32(0xBD1)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGotIngame(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Player '%s' entering in-game zone (CH22Type_GotIngame 0x%04X). Dispatching in-game ACKs...",
		client.RemoteAddr(), client.Username, pkt.OpCode)

	// 1. Packet 0x581C (Header 22, Type 28)
	p1 := protocol.NewPacketWriter(protocol.SH22Type_GotIngameAck)
	p1.WriteUint32(0x4d0bc167)
	_ = client.SendPacket(p1)

	// 2. Packet 0x581D (Header 22, Type 29) - zero active Kingdom Quests
	p2 := protocol.NewPacketWriter(protocol.SH22Type_KingdomQuestCount)
	p2.WriteInt16(0)
	_ = client.SendPacket(p2)

	// 3. Friends List (Header 21, Type 7) - Friend count = 0
	p3 := protocol.NewPacketWriter(protocol.SH21Type_FriendList)
	p3.WriteByte(0)
	_ = client.SendPacket(p3)

	// 4. Time Sync (Header 2, Type 69)
	p4 := protocol.NewPacketWriter(protocol.SH2Type_UnkTimePacket)
	p4.WriteInt16(256)
	_ = client.SendPacket(p4)

	log.Printf("[World][%s] In-game ACK sequence sent successfully for '%s'",
		client.RemoteAddr(), client.Username)
}
