package world

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler despacha los paquetes de red del World Server (WorldManager + Character unificados).
type Handler struct {
	worldID       byte
	charRepo      database.CharacterRepository
	accountRepo   database.AccountRepository
	shopRepo      database.AccountShopRepository
	sessionRepo   database.LoginSessionRepository
	logRepo       database.AccountLogRepository
	transfers     *database.TransferManager
	dataStore     *WorldDataStore
	partyMgr      *PartyManager
	friendMgr     *FriendManager
	guildMgr      *GuildManager
	kqMgr         *KQManager
	zoneIP        string
	zonePort      uint16
	activeClients map[int]*network.Client
	mu            sync.Mutex
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
	pm := NewPartyManager()
	fm := NewFriendManager(charRepo, pm)
	gm := NewGuildManager(charRepo)
	var km *KQManager
	if dataStore != nil {
		km = NewKQManager(dataStore.KingdomQuest)
	} else {
		km = NewKQManager(nil)
	}
	return &Handler{
		charRepo:      charRepo,
		transfers:     transfers,
		dataStore:     dataStore,
		partyMgr:      pm,
		friendMgr:     fm,
		guildMgr:      gm,
		kqMgr:         km,
		zoneIP:        zoneIP,
		zonePort:      zonePort,
		activeClients: make(map[int]*network.Client),
	}
}

// SetAccountRepo vincula los repositorios de cuenta, logs y sesiones.
func (h *Handler) SetAccountRepo(repo database.AccountRepository) {
	h.accountRepo = repo
	if s, ok := repo.(database.LoginSessionRepository); ok {
		h.sessionRepo = s
	}
	if l, ok := repo.(database.AccountLogRepository); ok {
		h.logRepo = l
	}
	if sh, ok := repo.(database.AccountShopRepository); ok {
		h.shopRepo = sh
	}
}

// SetWorldID asigna el ID de mundo correspondiente a este WorldServer.
func (h *Handler) SetWorldID(id byte) {
	h.worldID = id
}

// ActivePlayersCount retorna la cantidad de clientes activos en WorldService.
func (h *Handler) ActivePlayersCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.activeClients)
}

// OnClientDisconnect limpia la sesión activa en el mapa de concurrencia.
func (h *Handler) OnClientDisconnect(client *network.Client) {
	if client.AccountID > 0 {
		h.mu.Lock()
		if h.activeClients[client.AccountID] == client {
			delete(h.activeClients, client.AccountID)
		}
		h.mu.Unlock()
	}
}

// HandlePacket despacha paquetes basándose exclusivamente en OpCodes oficiales PROTO_NC_*.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_MISC_HEARTBEAT_ACK:
		// Client heartbeat acknowledge

	case protocol.PROTO_NC_MISC_GAMETIME_REQ:
		h.handleGameTime(client, pkt)

	case protocol.PROTO_NC_USER_LOGINWORLD_REQ:
		h.handleWorldClientKey(client, pkt)

	case protocol.PROTO_NC_USER_WILL_WORLD_SELECT_REQ:
		h.handleBackToWorldSelect(client, pkt)

	case protocol.PROTO_NC_USER_RETURN_CHECK_REQ:
		h.handleReturnCheck(client, pkt)

	case protocol.PROTO_NC_USER_SET_RETURN_CMD:
		h.handleSetReturn(client, pkt)

	case protocol.PROTO_NC_USER_LOGOUT_LAST_DAY_REQ:
		h.handleLogoutLastDay(client, pkt)

	case protocol.PROTO_NC_MISC_WEB_DB_KEY_REQ:
		h.handleWebDBKey(client, pkt)

	case protocol.PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ:
		h.handleConnectFromWhere(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_HIT3_ADD_CASH:
		h.handleEventHit3AddCash(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC:
		h.handleEventDoneMunsang(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_L20_DB_REQ:
		h.handleEventL20(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD:
		h.handleEventLongTimePlay(client, pkt)

	case protocol.PROTO_NC_CHAR_LOGIN_REQ:
		h.handleCharSelect(client, pkt)

	case protocol.PROTO_NC_AVATAR_CREATE_REQ:
		h.handleCreateCharacter(client, pkt)

	case protocol.PROTO_NC_AVATAR_ERASE_REQ:
		h.handleDeleteCharacter(client, pkt)

	case protocol.PROTO_NC_KQ_LIST_REFRESH_REQ:
		h.handleGotIngame(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ:
		h.handleGetQuickBar(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ:
		h.handleGetQuickBarState(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_GAME_REQ:
		h.handleGetGameSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ:
		h.handleGetClientSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ:
		h.handleGetShortCuts(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_GAME_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD:
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

	// Subsistema Oficial de Amigos (Friend - Dept: 21)
	case protocol.PROTO_NC_FRIEND_SET_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.friendMgr.RequestFriend(client, targetName)

	case protocol.PROTO_NC_FRIEND_SET_CONFIRM_ACK:
		h.friendMgr.AcceptFriend(client)

	case protocol.PROTO_NC_FRIEND_DEL_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.friendMgr.DeleteFriend(client, targetName)

	case protocol.PROTO_NC_FRIEND_POINT_REQ:
		h.friendMgr.SendFriendPoints(client)

	// Subsistema Oficial de Gremios (Guild - Dept: 29)
	case protocol.PROTO_NC_GUILD_LIST_REQ:
		h.guildMgr.HandleGuildList(client, 0)

	case protocol.PROTO_NC_GUILD_MAKE_REQ:
		name, _ := pkt.Reader.ReadString(16)
		password, _ := pkt.Reader.ReadString(12)
		guildType, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleGuildMake(client, name, password, guildType)

	case protocol.PROTO_NC_GUILD_MY_GUILD_INFO_REQ:
		h.guildMgr.HandleMyGuildInfo(client)

	case protocol.PROTO_NC_GUILD_MEMBER_LIST_REQ:
		h.guildMgr.HandleMemberList(client)

	case protocol.PROTO_NC_GUILD_MEMBER_INVITE_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleInvite(client, targetName)

	case protocol.PROTO_NC_GUILD_MEMBER_JOIN_ACK:
		guildName, _ := pkt.Reader.ReadString(16)
		joinAck, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleJoinAck(client, guildName, joinAck)

	case protocol.PROTO_NC_GUILD_MEMBER_LEAVE_REQ:
		h.guildMgr.HandleLeave(client)

	case protocol.PROTO_NC_GUILD_MEMBER_VANISH_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleVanish(client, targetName)

	case protocol.PROTO_NC_GUILD_MEMBER_GRADE_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		grade, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleGrade(client, targetName, grade)

	case protocol.PROTO_NC_GUILD_NOTIFY_REQ:
		nLen, _ := pkt.Reader.ReadUint16()
		notifyBytes, _ := pkt.Reader.ReadBytes(int(nLen))
		h.guildMgr.HandleNotify(client, string(notifyBytes))

	case protocol.PROTO_NC_GUILD_DISMISS_REQ:
		h.guildMgr.HandleDismiss(client)

	case protocol.PROTO_NC_GUILD_CHAT_REQ:
		_, _ = pkt.Reader.ReadByte() // itemLinkCount
		lenByte, _ := pkt.Reader.ReadByte()
		msgBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.guildMgr.HandleChat(client, string(msgBytes))

	// Subsistema Oficial de Academia de Gremios (Guild Academy - Dept: 38)
	case protocol.PROTO_NC_GUILD_ACADEMY_LIST_REQ:
		h.guildMgr.HandleAcademyList(client, 0)

	case protocol.PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_REQ:
		h.guildMgr.HandleAcademyMyInfo(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_REQ:
		h.guildMgr.HandleAcademyMemberList(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_REQ:
		guildName, _ := pkt.Reader.ReadString(16)
		h.guildMgr.HandleAcademyJoin(client, guildName)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_REQ:
		h.guildMgr.HandleAcademyLeave(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_CHAT_REQ:
		_, _ = pkt.Reader.ReadByte() // itemLinkCount
		lenByte, _ := pkt.Reader.ReadByte()
		msgBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.guildMgr.HandleAcademyChat(client, string(msgBytes))

	// Subsistema Oficial de Kingdom Quests (KQ - Dept: 22)
	case protocol.PROTO_NC_KQ_LIST_REQ:
		h.kqMgr.HandleList(client)

	case protocol.PROTO_NC_KQ_STATUS_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleStatus(client, handle)

	case protocol.PROTO_NC_KQ_JOIN_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		var level, class byte = 1, 1
		ch, err := h.charRepo.GetCharacterByID(client.CharacterID)
		if err == nil && ch != nil {
			level = ch.Level
			class = ch.Job
		}
		h.kqMgr.HandleJoin(client, handle, level, class)

	case protocol.PROTO_NC_KQ_JOIN_CANCEL_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleJoinCancel(client, handle)

	case protocol.PROTO_NC_KQ_JOIN_LIST_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleJoinList(client, handle)

	case protocol.PROTO_NC_KQ_TEAM_SELECT_REQ:
		teamType, _ := pkt.Reader.ReadByte()
		h.kqMgr.HandleTeamSelect(client, teamType)

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
	var token64 [64]byte

	if pkt.Reader.Remaining() >= 320 {
		uBytes, _ := pkt.Reader.ReadBytes(256)
		username = strings.TrimRight(string(uBytes), "\x00")
		username = strings.TrimSpace(username)

		tBytes, _ := pkt.Reader.ReadBytes(64)
		copy(token64[:], tBytes)
	} else if pkt.Reader.Remaining() >= 64 {
		tBytes, _ := pkt.Reader.ReadBytes(64)
		copy(token64[:], tBytes)
	} else if pkt.Reader.Remaining() >= 18+32 {
		_ = pkt.Reader.Skip(18)
		tBytes, _ := pkt.Reader.ReadBytes(32)
		copy(token64[:32], tBytes)
	}

	var authenticated bool
	var userNo int
	var admin byte

	// 1. Validar token oficial de 64 bytes contra MySQL (usp_Session_ConsumeWorldToken)
	if h.sessionRepo != nil {
		uNo, uID, err := h.sessionRepo.ConsumeWorldToken(h.worldID, token64, client.RemoteIP())
		if err == nil && uNo > 0 {
			authenticated = true
			userNo = uNo
			if username == "" {
				username = uID
			}
			if h.accountRepo != nil {
				acc, err := h.accountRepo.GetAccountByUsername(username)
				if err == nil && acc != nil {
					admin = acc.Admin
				}
			}
		}
	}

	// 2. Fallback de soporte para tests en memoria (TransferManager)
	if !authenticated && h.transfers != nil {
		tokenStr := strings.TrimRight(string(token64[:32]), "\x00")
		tokenStr = strings.TrimSpace(tokenStr)
		transfer, ok := h.transfers.ConsumeTransfer(tokenStr, token64[:], username, h.worldID)
		if ok && transfer != nil {
			authenticated = true
			userNo = transfer.AccountID
			username = transfer.Username
			admin = transfer.Admin
		}
	}

	if !authenticated {
		log.Printf("[World][%s] World token verification failed for user '%s'", client.RemoteAddr(), username)
		h.sendError(client, protocol.ErrInvalidCredentials)
		client.Close()
		return
	}

	// 3. Expulsión limpia de sesión previa duplicada (KICKOFFFROMWORLD_CMD / CONNECTCUT_CMD)
	h.mu.Lock()
	if prevClient, exists := h.activeClients[userNo]; exists && prevClient != nil && prevClient != client {
		log.Printf("[World][%s] Sesión previa activa detectada para UserNo %d, expulsando de WorldServer",
			client.RemoteAddr(), userNo)
		_ = prevClient.SendPacket(protocol.New_PROTO_NC_USER_KICKOFFFROMWORLD_CMD(uint32(userNo)))
		_ = prevClient.SendPacket(protocol.New_PROTO_NC_USER_CONNECTCUT_CMD(1))
		prevClient.Close()
		delete(h.activeClients, userNo)
	}
	h.activeClients[userNo] = client
	h.mu.Unlock()

	client.AccountID = userNo
	client.Username = username
	client.Admin = admin
	client.IsAuthenticated = true

	log.Printf("[World][%s] Authenticated user '%s' (AccountID: %d)",
		client.RemoteAddr(), client.Username, client.AccountID)

	// Enviar lista de personajes disponibles
	h.sendCharacterList(client)
}

// sendCharacterList despacha la lista oficial de avatares (PROTO_NC_USER_LOGINWORLD_ACK - 0x0C14).
func (h *Handler) sendCharacterList(client *network.Client) {
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		log.Printf("[World][%s] Error loading characters for account %d: %v",
			client.RemoteAddr(), client.AccountID, err)
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_ACK)
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

	// Enviar respuesta oficial de confirmación PROTO_NC_AVATAR_CREATESUCC_ACK (0x1406)
	resp := protocol.New_PROTO_NC_AVATAR_CREATESUCC_ACK(func(pw *protocol.PacketWriter) {
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

	resp := protocol.New_PROTO_NC_AVATAR_ERASESUCC_ACK(slot)
	_ = client.SendPacket(resp)
}

// handleBackToWorldSelect procesa la solicitud de volver a la selección de mundos
// generando un One-Time Password criptográfico oficial de 32 bytes (Name8 sOTP)
// y emitiendo PROTO_NC_USER_WILL_WORLD_SELECT_ACK (0x0C34).
func (h *Handler) handleBackToWorldSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.AccountID <= 0 {
		resp := protocol.New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK([32]byte{})
		_ = client.SendPacket(resp)
		return
	}

	var otp [32]byte
	const hexChars = "0123456789abcdef"
	for i := 0; i < 32; i++ {
		otp[i] = hexChars[rand.Intn(len(hexChars))]
	}

	if h.sessionRepo != nil {
		err := h.sessionRepo.CreateOTP(client.AccountID, otp)
		if err != nil {
			log.Printf("[World][%s] Error persistiendo sOTP en tActiveAccountSession para UserNo %d: %v",
				client.RemoteAddr(), client.AccountID, err)
		} else {
			log.Printf("[World][%s] Generado sOTP de 32 bytes para UserNo %d", client.RemoteAddr(), client.AccountID)
		}
	}

	resp := protocol.New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK(otp)
	_ = client.SendPacket(resp)
}

// handleReturnCheck procesa la consulta de usuario que retorna (0x0C60).
func (h *Handler) handleReturnCheck(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = u
		}
	}

	var isReturn byte
	var lastSec int32
	if h.logRepo != nil && userNo > 0 {
		diff, err := h.logRepo.GetDateDiffLastLogout(int(userNo))
		if err == nil && diff >= 30*24*3600 {
			isReturn = 1
			lastSec = int32(diff)
		}
	}

	resp := protocol.New_PROTO_NC_USER_RETURN_CHECK_ACK(handle, userNo, isReturn, lastSec)
	_ = client.SendPacket(resp)
}

// handleSetReturn procesa la confirmación de reclamo de paquete de retorno (0x0C64).
func (h *Handler) handleSetReturn(client *network.Client, pkt *protocol.Packet) {
	userNo := client.AccountID
	if pkt.Reader.Remaining() >= 4 {
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = int(u)
		}
	}

	if h.accountRepo != nil && userNo > 0 {
		_ = h.accountRepo.SetReturnUserItem(userNo)
		log.Printf("[World][%s] Paquete de retorno acreditado para UserNo %d", client.RemoteAddr(), userNo)
	}
}

// handleLogoutLastDay procesa la consulta de días desde la última desconexión (0x0C62).
func (h *Handler) handleLogoutLastDay(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = u
		}
	}

	var lastSec int32
	if h.logRepo != nil && userNo > 0 {
		diff, err := h.logRepo.GetDateDiffLastLogout(int(userNo))
		if err == nil {
			lastSec = int32(diff)
		}
	}

	resp := protocol.New_PROTO_NC_USER_LOGOUT_LAST_DAY_ACK(handle, userNo, lastSec)
	_ = client.SendPacket(resp)
}

// handleWebDBKey procesa la solicitud de token SSO para ItemShop Web (0x081D).
func (h *Handler) handleWebDBKey(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, _ := pkt.Reader.ReadUint32()
		if u > 0 {
			userNo = u
		}
	}

	var key string
	var errCode uint16 = 0
	if h.accountRepo != nil && userNo > 0 {
		k, err := h.accountRepo.GetOrIssueWebKey(int(userNo), "")
		if err != nil {
			errCode = 1
		} else {
			key = k
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_WEB_DB_KEY_ACK)
	pw.WriteUint16(errCode)
	pw.WriteUint16(handle)
	pw.WriteUint32(userNo)
	pw.WriteString(key, 32)
	_ = client.SendPacket(pw)
}

// handleConnectFromWhere procesa la verificación de origen de conexión y PC Bang (0x082F).
func (h *Handler) handleConnectFromWhere(client *network.Client, pkt *protocol.Packet) {
	var hZone uint16
	var zoneHeader []byte
	var connectIPStr string = client.RemoteIP()

	rem := pkt.Reader.Remaining()
	if rem >= 12 {
		hZone, _ = pkt.Reader.ReadUint16()
		zoneHeader, _ = pkt.Reader.ReadBytes(6)
		ipBytes, _ := pkt.Reader.ReadBytes(4)
		if len(ipBytes) == 4 && (ipBytes[0] != 0 || ipBytes[1] != 0 || ipBytes[2] != 0 || ipBytes[3] != 0) {
			connectIPStr = fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		}
	} else if rem >= 10 {
		zoneHeader, _ = pkt.Reader.ReadBytes(6)
		ipBytes, _ := pkt.Reader.ReadBytes(4)
		if len(ipBytes) == 4 && (ipBytes[0] != 0 || ipBytes[1] != 0 || ipBytes[2] != 0 || ipBytes[3] != 0) {
			connectIPStr = fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		}
	}

	if len(zoneHeader) < 6 {
		zoneHeader = make([]byte, 6)
	}

	var isPcbang bool
	if h.accountRepo != nil {
		pc, _ := h.accountRepo.CheckPcbangIP(connectIPStr)
		isPcbang = pc
	}

	var connectFrom uint8 = 0
	if isPcbang {
		connectFrom = 1
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK)
	pw.WriteUint16(hZone)
	pw.WriteBytes(zoneHeader)
	pw.WriteByte(connectFrom)
	_ = client.SendPacket(pw)
}

// handleEventHit3AddCash procesa la acreditación de cash por evento (0x0820).
func (h *Handler) handleEventHit3AddCash(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	if rem >= 28 {
		userIDStr, _ := pkt.Reader.ReadString(20)
		addCash, _ := pkt.Reader.ReadUint64()
		userIDStr = strings.TrimSpace(userIDStr)
		if h.shopRepo != nil && userIDStr != "" && addCash > 0 {
			_ = h.shopRepo.AddEventCashByUserID(userIDStr, int64(addCash))
		}
	} else if rem >= 8 {
		userNo, _ := pkt.Reader.ReadUint32()
		cash, _ := pkt.Reader.ReadUint32()
		if h.shopRepo != nil && userNo > 0 && cash > 0 {
			_ = h.shopRepo.AddCash(int(userNo), 0, int(cash), "EventHit3")
		}
	}
}

// handleEventDoneMunsang procesa la validación y canje de cupón cultural Munsang (0x0838).
func (h *Handler) handleEventDoneMunsang(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var (
		userNo  uint32
		charNo  uint32
		worldNo uint8
		hZone   uint16
		questID uint16
		nError  uint16
	)

	if rem >= 13 {
		userNo, _ = pkt.Reader.ReadUint32()
		charNo, _ = pkt.Reader.ReadUint32()
		worldNo, _ = pkt.Reader.ReadByte()
		hZone, _ = pkt.Reader.ReadUint16()
		questID, _ = pkt.Reader.ReadUint16()
	} else if rem >= 4 {
		userNo, _ = pkt.Reader.ReadUint32()
	}

	if h.shopRepo != nil && userNo > 0 {
		errCode, err := h.shopRepo.RedeemMunsangEvent(int(userNo), questID)
		if err == nil {
			nError = errCode
		} else {
			nError = 5445 // Código de fallo oficial de Account.exe (0x1545)
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM)
	pw.WriteUint32(userNo)
	pw.WriteUint32(charNo)
	pw.WriteByte(worldNo)
	pw.WriteUint16(nError)
	pw.WriteUint16(hZone)
	pw.WriteUint16(questID)
	_ = client.SendPacket(pw)
}

// handleEventL20 procesa el registro y verificación del evento Nivel 20 (0x082C).
func (h *Handler) handleEventL20(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var (
		wmHandle   uint16
		userID     string
		userNo     uint32
		eventNo    uint8 = 1
		worldNo    uint8 = h.worldID
		charNo     uint32
		eventPoint uint32
		nError     uint16
	)

	if rem >= 272 {
		wmHandle, _ = pkt.Reader.ReadUint16()
		uID, _ := pkt.Reader.ReadString(256)
		userID = strings.TrimSpace(uID)
		userNo, _ = pkt.Reader.ReadUint32()
		eventNo, _ = pkt.Reader.ReadByte()
		worldNo, _ = pkt.Reader.ReadByte()
		charNo, _ = pkt.Reader.ReadUint32()
		eventPoint, _ = pkt.Reader.ReadUint32()
	} else if rem >= 4 {
		u, _ := pkt.Reader.ReadUint32()
		userNo = u
	} else if client.AccountID > 0 {
		userNo = uint32(client.AccountID)
	}

	if h.logRepo != nil && userNo > 0 {
		err := h.logRepo.EventL20SetFull(int(userNo), userID, eventNo, worldNo, int(charNo), int(eventPoint))
		if err != nil {
			nError = 1
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_EVENT_L20_DB_ACK)
	pw.WriteUint16(wmHandle)
	pw.WriteUint32(userNo)
	pw.WriteByte(eventNo)
	pw.WriteUint32(charNo)
	pw.WriteUint16(nError)
	_ = client.SendPacket(pw)
}

// handleEventLongTimePlay procesa el otorgamiento de boletos por tiempo de juego (0x0843).
func (h *Handler) handleEventLongTimePlay(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	if rem >= 6 {
		userNo, _ := pkt.Reader.ReadUint32()
		playTime, _ := pkt.Reader.ReadUint16()
		if h.logRepo != nil && userNo > 0 {
			_ = h.logRepo.EventiPadSetTicket(int(userNo), playTime)
		}
	} else if rem >= 4 {
		userNo, _ := pkt.Reader.ReadUint32()
		if h.logRepo != nil && userNo > 0 {
			_ = h.logRepo.EventiPadSetTicket(int(userNo))
		}
	}
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

	// Enviar configuración de interfaz del personaje (SendToClient_CharOptionData @ 0x00431680, DataType == 2 NA2016):
	// 1. 0x7032 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD)
	// 2. 0x7033 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD)
	// 3. 0x7034 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD)
	h.sendCharOptionDataImprove(client, selected.ID)
}

// sendCharOptionDataImprove envía automáticamente los paquetes del sistema mejorado de opciones (DataType == 2)
// justo después de PROTO_NC_CHAR_LOGIN_ACK (0x1003), replicando CWMClientSession::SendToClient_CharOptionData.
func (h *Handler) sendCharOptionDataImprove(client *network.Client, charNo int) {
	var shortcutEntries []protocol.SHORT_CUT_DATA
	if shortcuts, err := h.charRepo.GetCharacterShortcuts(charNo); err == nil && len(shortcuts) > 0 {
		shortcutEntries = make([]protocol.SHORT_CUT_DATA, 0, len(shortcuts))
		for _, sc := range shortcuts {
			if sc == nil {
				continue
			}
			shortcutEntries = append(shortcutEntries, protocol.SHORT_CUT_DATA{
				SlotNo: sc.SlotNo,
				CodeNo: uint16(sc.CodeNo),
				Value:  sc.Value,
			})
		}
	}

	var rawShortcutData, rawKeyMapData, rawGameOptionData []byte
	if opts, err := h.charRepo.GetCharacterOptions(charNo); err == nil && opts != nil {
		if len(shortcutEntries) == 0 && len(opts.ShortCutData) >= 2 && (len(opts.ShortCutData)-2)%7 == 0 {
			rawShortcutData = opts.ShortCutData
		}
		if len(opts.KeyMapping) >= 2 && (len(opts.KeyMapping)-2)%4 == 0 {
			rawKeyMapData = opts.KeyMapping
		}
		if len(opts.Game) >= 2 && (len(opts.Game)-2)%3 == 0 {
			rawGameOptionData = opts.Game
		}
	}

	// 1. 0x7032: PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
	pktShortcut := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD(shortcutEntries, rawShortcutData)
	_ = client.SendPacket(pktShortcut)

	// 2. 0x7033: PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
	pktKeyMap := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD(nil, rawKeyMapData)
	_ = client.SendPacket(pktKeyMap)

	// 3. 0x7034: PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
	pktGameOpt := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD(nil, rawGameOptionData)
	_ = client.SendPacket(pktGameOpt)
}

func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINFAIL_ACK)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) sendCharCreationError(client *network.Client, errCode protocol.CreateCharError) {
	resp := protocol.New_PROTO_NC_AVATAR_CREATEFAIL_ACK(uint16(errCode))
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
	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD:
		optType = "shortcut_data"
	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD:
		optType = "shortcut_size"
	case protocol.PROTO_NC_CHAR_OPTION_SET_GAME_CMD:
		optType = "game"
	case protocol.PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD:
		optType = "window_pos"
	case protocol.PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD:
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

	// 1. Registrar en KQManager y despachar reloj y lista oficial de Kingdom Quests
	h.kqMgr.RegisterOnline(client, client.CharacterID)
	h.kqMgr.HandleList(client)

	// 2. Despachar lista oficial de amigos del jugador desde tFriend (0x5407)
	h.friendMgr.SendFriendList(client)
	h.friendMgr.BroadcastLogin(client, client.Username, "RouN")

	// 3. Registrar en GuildManager y enviar información de gremio si pertenece a uno
	h.guildMgr.RegisterOnline(client, client.CharacterID, client.Username)
	h.guildMgr.HandleMyGuildInfo(client)
	h.guildMgr.HandleMemberList(client)

	// 4. Packet PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD (0x0845 / Dept:2, Cmd:69) - Sincronización oficial del reloj
	p4 := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD)
	now := time.Now()
	millis := uint32(now.Hour()*3600000 + now.Minute()*60000 + now.Second()*1000 + now.Nanosecond()/1000000)
	p4.WriteUint32(millis)
	_ = client.SendPacket(p4)

	log.Printf("[World][%s] Official in-game ACK sequence sent successfully for '%s'",
		client.RemoteAddr(), client.Username)
}
