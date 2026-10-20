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

// Handler despacha los paquetes de red del World Server (WorldManager + Character unificados).
type Handler struct {
	charRepo  database.CharacterRepository
	transfers *database.TransferManager
	dataStore *WorldDataStore
	partyMgr  *PartyManager
	zoneIP    string
	zonePort  uint16
}

// NewHandler inicializa un nuevo despachador de paquetes para WorldService.
func NewHandler(
	charRepo database.CharacterRepository,
	transfers *database.TransferManager,
	dataStore *WorldDataStore,
	zoneIP string,
	zonePort uint16,
) *Handler {
	if zoneIP == "" {
		zoneIP = "127.0.0.1"
	}
	if zonePort == 0 {
		zonePort = 9210
	}
	return &Handler{
		charRepo:  charRepo,
		transfers: transfers,
		dataStore: dataStore,
		partyMgr:  NewPartyManager(),
		zoneIP:    zoneIP,
		zonePort:  zonePort,
	}
}

// HandlePacket despacha paquetes basándose exclusivamente en OpCodes oficiales PROTO_NC_*.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_MISC_HEARTBEAT_ACK:
		// Client heartbeat acknowledge

	case protocol.PROTO_NC_MISC_GAMETIME_REQ:
		h.handleGameTime(client, pkt)

	case protocol.PROTO_NC_MISC_WORLD_CLIENT_KEY_REQ:
		h.handleWorldClientKey(client, pkt)

	case protocol.PROTO_NC_USER_BACKTOWORLDLIST_REQ:
		h.handleBackToWorldSelect(client, pkt)

	case protocol.PROTO_NC_CHAR_LOGIN_REQ:
		h.handleCharSelect(client, pkt)

	case protocol.PROTO_NC_CHAR_CREATE_REQ:
		h.handleCreateCharacter(client, pkt)

	case protocol.PROTO_NC_CHAR_DELETE_REQ:
		h.handleDeleteCharacter(client, pkt)

	case protocol.PROTO_NC_CHAR_INGAME_REQ:
		h.handleGotIngame(client, pkt)

	case protocol.PROTO_NC_CHAR_GET_QUICKBAR_REQ:
		h.handleGetQuickBar(client, pkt)

	case protocol.PROTO_NC_CHAR_GET_QUICKBAR_STATE_REQ:
		h.handleGetQuickBarState(client, pkt)

	case protocol.PROTO_NC_CHAR_GET_GAME_SETTINGS_REQ:
		h.handleGetGameSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_GET_CLIENT_SETTINGS_REQ:
		h.handleGetClientSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_GET_SHORTCUTS_REQ:
		h.handleGetShortCuts(client, pkt)

	case protocol.PROTO_NC_CHAR_SAVE_QUICKBAR_REQ,
		protocol.PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ,
		protocol.PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ,
		protocol.PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ,
		protocol.PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ:
		h.handleSaveSettings(client, pkt)

	case protocol.PROTO_NC_PRISON_GET_REQ:
		h.handlePrisonGet(client, pkt)

	// Subsistema Oficial de Grupos (Party - Dept: 14)
	case protocol.PROTO_NC_PARTY_JOIN_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.partyMgr.InviteMember(client, targetName)

	case protocol.PROTO_NC_PARTY_JOINPROPOSE_ALLOW_ACK:
		h.partyMgr.AcceptInvite(client)

	case protocol.PROTO_NC_PARTY_JOINPROPOSE_REJECT_ACK:
		h.partyMgr.RejectInvite(client)

	case protocol.PROTO_NC_PARTY_LEAVE_REQ:
		h.partyMgr.LeaveParty(client)

	case protocol.PROTO_NC_PARTY_KICKOFF_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.partyMgr.KickMember(client, targetName)

	case protocol.PROTO_NC_PARTY_CHANGEMASTER_REQ:
		newMaster, _ := pkt.Reader.ReadString(20)
		h.partyMgr.ChangeMaster(client, newMaster)

	case protocol.PROTO_NC_PARTY_ITEM_LOOTING_SET:
		lootType, _ := pkt.Reader.ReadByte()
		itemGrade, _ := pkt.Reader.ReadByte()
		h.partyMgr.SetLooting(client, lootType, itemGrade)

	case protocol.PROTO_NC_PARTY_FINDER_ADD_REQ:
		purpose, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())
		h.partyMgr.AddFinderAd(client, purpose)

	case protocol.PROTO_NC_PARTY_FINDER_DELETE_REQ:
		h.partyMgr.DeleteFinderAd(client)

	case protocol.PROTO_NC_PARTY_FINDER_LIST_REQ:
		h.partyMgr.ListFinderAds(client)

	default:
		log.Printf("[World][%s] Unhandled OpCode: 0x%04X (Dept=%d, Cmd=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

// handleGameTime responde con la hora del servidor (PROTO_NC_MISC_GAMETIME_ACK - 0x080E).
func (h *Handler) handleGameTime(client *network.Client, pkt *protocol.Packet) {
	now := time.Now()
	resp := protocol.New_PROTO_NC_MISC_GAMETIME_ACK(byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	_ = client.SendPacket(resp)
}

// handleWorldClientKey valida el token de traspaso proveniente de LoginService
// y autentica la sesión del cliente en el WorldServer.
// Soporta la estructura oficial PROTO_NC_USER_LOGINWORLD_REQ (320 bytes payload):
//   - char user[256] /* Name256Byte */
//   - unsigned short validate_new[32] (64 bytes)
func (h *Handler) handleWorldClientKey(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Received PROTO_NC_USER_LOGINWORLD_REQ (len %d)",
		client.RemoteAddr(), len(pkt.Raw))

	var username string
	var token string

	if pkt.Reader.Remaining() >= 256 {
		uBytes, _ := pkt.Reader.ReadBytes(256)
		username = strings.TrimRight(string(uBytes), "\x00")
		username = strings.TrimSpace(username)

		if pkt.Reader.Remaining() >= 32 {
			tBytes, _ := pkt.Reader.ReadBytes(32)
			token = strings.TrimRight(string(tBytes), "\x00")
			token = strings.TrimSpace(token)
		}
	} else if pkt.Reader.Remaining() >= 18+32 {
		// Fallback para paquetes heredados/mockeados en tests
		_ = pkt.Reader.Skip(18)
		token, _ = pkt.Reader.ReadString(32)
		token = strings.TrimSpace(token)
	}

	transfer, ok := h.transfers.ConsumeTransfer(token, pkt.Raw, username)
	if !ok {
		log.Printf("[World][%s] Transfer token '%s' (user: '%s') not found or expired",
			client.RemoteAddr(), token, username)
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

	// Enviar lista de personajes disponibles
	h.sendCharacterList(client)
}

// sendCharacterList despacha la lista oficial de avatares (PROTO_NC_CHAR_CLIENT_CHARLIST_ACK - 0x0C14).
func (h *Handler) sendCharacterList(client *network.Client) {
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		log.Printf("[World][%s] Error loading characters for account %d: %v",
			client.RemoteAddr(), client.AccountID, err)
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARLIST_ACK)
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

// handleCreateCharacter procesa la creación de personaje validando contra BadNameFilter.shn
// y asignando valores base desde default_character_data.json.
func (h *Handler) handleCreateCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[World][%s] Unauthenticated CreateCharacter request rejected", client.RemoteAddr())
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

	isMale := ((jobGender >> 7) & 0x01) == 1
	classID := byte((jobGender >> 2) & 0x1F)
	job := protocol.Job(classID)

	switch job {
	case protocol.Job_Fighter, protocol.Job_Cleric, protocol.Job_Archer, protocol.Job_Mage, protocol.Job_Trickster, protocol.Job_Crusader:
		// Clase válida soportada por el cliente oficial
	default:
		log.Printf("[World][%s] Invalid job ID %d requested", client.RemoteAddr(), classID)
		h.sendCharCreationError(client, protocol.CreateCharError_WrongClass)
		return
	}

	// 1. Validación de apariencia oficial (FaceInfo, HairInfo, HairColorInfo y límites de clase)
	if h.dataStore != nil && !h.dataStore.ValidateAppearance(classID, hair, color, style) {
		log.Printf("[World][%s] Character appearance rejected by limits (Class=%d, Hair=%d, Color=%d, Face=%d)",
			client.RemoteAddr(), classID, hair, color, style)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	// 2. Validación contra BadNameFilter.shn oficial
	if h.dataStore != nil && h.dataStore.BadNames != nil && h.dataStore.BadNames.IsBadName(name) {
		log.Printf("[World][%s] Character name '%s' rejected by BadNameFilter", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	// 3. Validación de unicidad de nombre en la base de datos
	if existing, _ := h.charRepo.GetCharacterByName(name); existing != nil {
		log.Printf("[World][%s] Character name '%s' already exists in database", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_NameTaken)
		return
	}

	// 4. Cargar parámetros base desde default_character_data.json
	startMap := "RouN"
	startX := 6888
	startY := 7045
	hp := 100
	sp := 100
	hpSoul := int16(15)
	spSoul := int16(7)
	money := int64(0)

	if h.dataStore != nil && h.dataStore.DefaultData != nil {
		if classData, exists := h.dataStore.DefaultData.GetClassData(classID); exists {
			if classData.StartMap != "" {
				startMap = classData.StartMap
			}
			if classData.StartX != 0 {
				startX = classData.StartX
			}
			if classData.StartY != 0 {
				startY = classData.StartY
			}
			if classData.HP != 0 {
				hp = classData.HP
			}
			if classData.SP != 0 {
				sp = classData.SP
			}
			if classData.HPSoul != 0 {
				hpSoul = classData.HPSoul
			}
			if classData.SPSoul != 0 {
				spSoul = classData.SPSoul
			}
			money = classData.Money
		}
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
		TownName:  startMap,
		XPos:      startX,
		YPos:      startY,
		CurHP:     hp,
		CurSP:     sp,
		Equips:    map[byte]uint16{wSlot: wItem},
	}
	newChar.TCharacterState.Money = money
	newChar.TCharacterState.HPStones = hpSoul
	newChar.TCharacterState.SPStones = spSoul

	if h.dataStore != nil && h.dataStore.DefaultData != nil {
		if classData, exists := h.dataStore.DefaultData.GetClassData(classID); exists {
			var invenSlot int16 = 0
			for _, sItem := range classData.StarterItems {
				if uint16(sItem.ItemID) == wItem {
					continue // Ya se inserta como arma equipada en slot de arma
				}
				amount := sItem.Amount
				if amount <= 0 {
					amount = 1
				}
				newChar.Items = append(newChar.Items, &database.Item{
					StorageType: 9, // Inventario / Mochila (Box 9)
					Storage:     invenSlot,
					ItemID:      uint16(sItem.ItemID),
					Flags:       int32(amount),
				})
				invenSlot++
			}

			for _, sc := range classData.Shortcuts {
				newChar.Shortcuts = append(newChar.Shortcuts, &database.Shortcut{
					SlotNo:   sc.Slot,
					CodeType: 1,
					CodeNo:   int16(sc.CodeNo),
					Value:    int32(sc.Value),
				})
			}

			for _, sk := range classData.StarterSkills {
				newChar.Skills = append(newChar.Skills, &database.CharacterSkill{
					SkillNo:    uint16(sk),
					SkillLevel: 1,
				})
			}
		}
	}

	if err := h.charRepo.CreateCharacter(newChar); err != nil {
		log.Printf("[World][%s] Failed to create character '%s': %v", client.RemoteAddr(), name, err)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	log.Printf("[World][%s] Created character '%s' (ID: %d, Job: %d, Slot: %d) for user '%s'",
		client.RemoteAddr(), name, newChar.ID, classID, slot, client.Username)

	// Enviar respuesta oficial de confirmación PROTO_NC_CHAR_CREATE_ACK (0x1406)
	resp := protocol.New_PROTO_NC_CHAR_CREATE_ACK(func(pw *protocol.PacketWriter) {
		WriteBasicCharInfo(newChar, pw)
	})
	_ = client.SendPacket(resp)
}

// handleDeleteCharacter ejecuta el soft-delete oficial liberando el slot y nickname.
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

	log.Printf("[World][%s] Soft-deleted character in slot %d for user '%s'", client.RemoteAddr(), slot, client.Username)

	resp := protocol.New_PROTO_NC_CHAR_DELETE_ACK(slot)
	_ = client.SendPacket(resp)
}

// handleBackToWorldSelect procesa la solicitud de volver a la selección de mundos.
func (h *Handler) handleBackToWorldSelect(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.New_PROTO_NC_USER_BACKTOWORLDLIST_ACK()
	_ = client.SendPacket(resp)
}

// handleCharSelect procesa la selección de un personaje y entrega el traspaso oficial hacia ZoneServer.
func (h *Handler) handleCharSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		log.Printf("[World][%s] Invalid character slot in CharSelect: %v", client.RemoteAddr(), err)
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
		log.Printf("[World][%s] Character in slot %d not found for account %d",
			client.RemoteAddr(), slot, client.AccountID)
		failResp := protocol.New_PROTO_NC_CHAR_LOGINFAIL_ACK(uint16(protocol.ErrInvalidCredentials))
		_ = client.SendPacket(failResp)
		return
	}

	client.CharacterID = selected.ID
	client.Slot = slot
	h.partyMgr.RegisterOnline(client, selected.ID, selected.Name, selected.Job, selected.Level)

	randomID := client.RandomID
	if randomID == 0 {
		randomID = uint16(rand.Intn(65534) + 1)
		client.RandomID = randomID
	}
	clientIP := client.RemoteIP()

	h.transfers.CreateZoneTransfer(
		client.AccountID,
		client.Username,
		selected.Name,
		selected.ID,
		randomID,
		client.Admin,
		clientIP,
	)

	log.Printf("[World][%s] Character '%s' (ID: %d) selected. Zone transfer token generated for %s:%d",
		client.RemoteAddr(), selected.Name, selected.ID, h.zoneIP, h.zonePort)

	// Enviar respuesta oficial PROTO_NC_CHAR_LOGIN_ACK (0x1003) con IP y Puerto de Zone
	resp := protocol.New_PROTO_NC_CHAR_LOGIN_ACK(h.zoneIP, h.zonePort)
	_ = client.SendPacket(resp)
}

func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINFAIL_ACK)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) sendCharCreationError(client *network.Client, errCode protocol.CreateCharError) {
	resp := protocol.New_PROTO_NC_CHAR_CREATE_FAIL_ACK(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) getActiveCharNo(client *network.Client) int {
	if client.CharacterID > 0 {
		return client.CharacterID
	}
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err == nil && len(chars) > 0 {
		return chars[0].ID
	}
	return 0
}

func (h *Handler) handleGetQuickBar(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.ShortCutData) > 0 {
			data = opts.ShortCutData
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK(data))
}

func (h *Handler) handleGetQuickBarState(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.ShortCutSize) > 0 {
			data = opts.ShortCutSize
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK(data))
}

func (h *Handler) handleGetGameSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.Game) > 0 {
			data = opts.Game
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_GAME_ACK(data))
}

func (h *Handler) handleGetClientSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.WindowsPos) > 0 {
			data = opts.WindowsPos
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK(data))
}

func (h *Handler) handleGetShortCuts(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.KeyMapping) > 0 {
			data = opts.KeyMapping
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK(data))
}

func (h *Handler) handleSaveSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	if charNo <= 0 {
		log.Printf("[World][%s] Cannot save settings: No active character selected", client.RemoteAddr())
		return
	}

	payload, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
	var optType string

	switch pkt.OpCode {
	case protocol.PROTO_NC_CHAR_SAVE_QUICKBAR_REQ:
		optType = "shortcut_data"
	case protocol.PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ:
		optType = "shortcut_size"
	case protocol.PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ:
		optType = "game"
	case protocol.PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ:
		optType = "window_pos"
	case protocol.PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ:
		optType = "key_mapping"
	default:
		log.Printf("[World][%s] Unknown save settings OpCode 0x%04X", client.RemoteAddr(), pkt.OpCode)
		return
	}

	if err := h.charRepo.SaveCharacterOption(charNo, optType, payload); err != nil {
		log.Printf("[World][%s] Error saving %s for char %d: %v", client.RemoteAddr(), optType, charNo, err)
		return
	}

	log.Printf("[World][%s] Saved %s (%d bytes) in tCharacterOptions for char %d",
		client.RemoteAddr(), optType, len(payload), charNo)
}

func (h *Handler) handlePrisonGet(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_GET_ACK)
	resp.WriteByte(0) // bInPrison = 0 (libre / no encarcelado)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGotIngame(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Player '%s' entering in-game zone (0x%04X). Dispatching official in-game ACKs...",
		client.RemoteAddr(), client.Username, pkt.OpCode)

	// 1. Packet PROTO_NC_KQ_LIST_TIME_ACK (0x581C) con timestamp Unix dinámico del servidor
	p1 := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_LIST_TIME_ACK)
	p1.WriteUint32(uint32(time.Now().Unix()))
	_ = client.SendPacket(p1)

	// 2. Packet PROTO_NC_KQ_LIST_COUNT_ACK (0x581D) - 0 Kingdom Quests activas iniciales
	p2 := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_LIST_COUNT_ACK)
	p2.WriteInt16(0)
	_ = client.SendPacket(p2)

	// 3. Packet PROTO_NC_FRIEND_LIST_ACK (21, 7) - 0 amigos
	p3 := protocol.NewPacketWriter(protocol.PROTO_NC_FRIEND_LIST_ACK)
	p3.WriteByte(0)
	_ = client.SendPacket(p3)

	// 4. Packet PROTO_NC_MISC_TIME_SYNC_ACK (0x0845 / Dept:2, Cmd:69) - Sincronización oficial del reloj
	p4 := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_TIME_SYNC_ACK)
	now := time.Now()
	millis := uint32(now.Hour()*3600000 + now.Minute()*60000 + now.Second()*1000 + now.Nanosecond()/1000000)
	p4.WriteUint32(millis)
	_ = client.SendPacket(p4)

	log.Printf("[World][%s] Official in-game ACK sequence sent successfully for '%s'",
		client.RemoteAddr(), client.Username)
}
