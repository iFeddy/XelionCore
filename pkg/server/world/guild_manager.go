package world

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// GuildOnlinePlayer representa la sesión activa de un jugador para difusión en tiempo real.
type GuildOnlinePlayer struct {
	CharID   int
	CharName string
	Client   *network.Client
}

// GuildPendingInvite almacena una invitación de gremio pendiente de confirmación.
type GuildPendingInvite struct {
	GuildID         int
	GuildName       string
	InviterCharID   int
	InviterCharName string
	ExpiresAt       time.Time
}

// GuildManager administra la lógica de negocios y la difusión de eventos de Gremios y Academia.
type GuildManager struct {
	mu             sync.RWMutex
	charRepo       database.CharacterRepository
	onlinePlayers  map[int]*GuildOnlinePlayer        // CharID -> Session
	onlineByName   map[string]*GuildOnlinePlayer     // lowercase(CharName) -> Session
	pendingInvites map[int]*GuildPendingInvite       // TargetCharID -> Invite
}

// NewGuildManager inicializa el administrador de Gremios y Academia.
func NewGuildManager(repo database.CharacterRepository) *GuildManager {
	return &GuildManager{
		charRepo:       repo,
		onlinePlayers:  make(map[int]*GuildOnlinePlayer),
		onlineByName:   make(map[string]*GuildOnlinePlayer),
		pendingInvites: make(map[int]*GuildPendingInvite),
	}
}

// RegisterOnline registra la sesión de un jugador y notifica su conexión al gremio si corresponde.
func (gm *GuildManager) RegisterOnline(client *network.Client, charID int, charName string) {
	gm.mu.Lock()
	player := &GuildOnlinePlayer{
		CharID:   charID,
		CharName: charName,
		Client:   client,
	}
	gm.onlinePlayers[charID] = player
	gm.onlineByName[strings.ToLower(strings.TrimSpace(charName))] = player
	gm.mu.Unlock()

	// Avisar al gremio que el miembro inició sesión
	guild, _, err := gm.charRepo.GetGuildByCharID(charID)
	if err == nil && guild != nil {
		gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_MEMBER_LOGON_CMD(charName), charID)
	}
}

// UnregisterOnline remueve la sesión activa y notifica desconexión a sus compañeros de gremio.
func (gm *GuildManager) UnregisterOnline(charID int) {
	gm.mu.Lock()
	player, ok := gm.onlinePlayers[charID]
	if ok {
		delete(gm.onlinePlayers, charID)
		delete(gm.onlineByName, strings.ToLower(strings.TrimSpace(player.CharName)))
		delete(gm.pendingInvites, charID)
	}
	gm.mu.Unlock()

	if player != nil {
		guild, _, err := gm.charRepo.GetGuildByCharID(charID)
		if err == nil && guild != nil {
			gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_MEMBER_LOGOFF_CMD(player.CharName), charID)
		}
	}
}

// BroadcastToGuild envía un paquete a todos los miembros online de un gremio (excluyendo opcionalmente a un CharID).
func (gm *GuildManager) BroadcastToGuild(guildID int, packet *protocol.PacketWriter, excludeCharID int) {
	members, err := gm.charRepo.GetGuildMembers(guildID)
	if err != nil {
		return
	}

	gm.mu.RLock()
	defer gm.mu.RUnlock()

	for _, m := range members {
		if m.CharID == excludeCharID {
			continue
		}
		if p, ok := gm.onlinePlayers[m.CharID]; ok && p.Client != nil {
			_ = p.Client.SendPacket(packet)
		}
	}
}

// BroadcastToAcademy envía un paquete a todos los alumnos online de una academia.
func (gm *GuildManager) BroadcastToAcademy(guildID int, packet *protocol.PacketWriter, excludeCharID int) {
	students, err := gm.charRepo.GetAcademyMembers(guildID)
	if err != nil {
		return
	}

	gm.mu.RLock()
	defer gm.mu.RUnlock()

	for _, s := range students {
		if s.CharID == excludeCharID {
			continue
		}
		if p, ok := gm.onlinePlayers[s.CharID]; ok && p.Client != nil {
			_ = p.Client.SendPacket(packet)
		}
	}
}

// ============================================================================
// HANDLERS OFICIALES DE GREMIO (DEPARTAMENTO 29)
// ============================================================================

// HandleGuildList devuelve la lista oficial de gremios activos (PROTO_NC_GUILD_LIST_ACK 0x7404).
func (gm *GuildManager) HandleGuildList(client *network.Client, partMask byte) {
	guilds, err := gm.charRepo.GetAllGuilds()
	if err != nil {
		log.Printf("[Guild] Error loading guild list: %v", err)
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_LIST_ACK(1, partMask, 0, 0, nil))
		return
	}

	var entries []protocol.GuildListItemEntry
	for _, g := range guilds {
		entries = append(entries, protocol.GuildListItemEntry{
			GuildNo:      g.GuildID,
			GuildName:    g.GuildName,
			MasterCharID: g.MasterCharName,
			Type:         g.Type,
			Grade:        g.Grade,
			MemberCnt:    g.MemberCount,
			MaxMembers:   g.MaxMembers,
		})
	}

	resp := protocol.New_PROTO_NC_GUILD_LIST_ACK(0, partMask, uint16(len(entries)), uint16(len(entries)), entries)
	_ = client.SendPacket(resp)
}

// HandleGuildMake crea un nuevo gremio e inicializa al creador como Maestro (PROTO_NC_GUILD_MAKE_ACK 0x7406).
func (gm *GuildManager) HandleGuildMake(client *network.Client, name, password string, guildType byte) {
	charID := client.CharacterID
	if charID <= 0 {
		return
	}

	cleanName := strings.TrimSpace(name)
	if len(cleanName) < 2 || len(cleanName) > 16 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MAKE_ACK(1, 0, cleanName, guildType, time.Now()))
		return
	}

	guildID, err := gm.charRepo.CreateGuild(charID, cleanName, password, guildType)
	if err != nil {
		log.Printf("[Guild] Failed to create guild '%s' by char %d: %v", cleanName, charID, err)
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MAKE_ACK(1, 0, cleanName, guildType, time.Now()))
		return
	}

	now := time.Now()
	log.Printf("[Guild] Guild '%s' (ID: %d) successfully created by char %d", cleanName, guildID, charID)
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MAKE_ACK(0, uint32(guildID), cleanName, guildType, now))

	// Enviar información detallada del gremio recién creado
	gm.HandleMyGuildInfo(client)
}

// HandleMyGuildInfo despacha los datos completos del gremio del jugador (PROTO_NC_GUILD_MY_GUILD_INFO_ACK 0x7419).
func (gm *GuildManager) HandleMyGuildInfo(client *network.Client) {
	charID := client.CharacterID
	if charID <= 0 {
		return
	}

	guild, _, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || guild.DismissStatus != 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MY_GUILD_INFO_ACK(1, nil))
		return
	}

	members, _ := gm.charRepo.GetGuildMembers(guild.ID)
	details := &protocol.GuildClientDetails{
		No:            uint32(guild.ID),
		Name:          guild.Name,
		Money:         guild.Money,
		Type:          guild.Type,
		Grade:         guild.Grade,
		Fame:          guild.Fame,
		StoneLevel:    guild.StoneLevel,
		Exp:           guild.Exp,
		CreateDate:    guild.CreateDate,
		NumMembers:    uint16(len(members)),
		MaxMembers:    guild.MaxMembers,
		DismissStatus: guild.DismissStatus,
		DismissDate:   guild.DismissDate,
		Intro:         guild.Intro,
		NotifyDate:    guild.NotifyDate,
		NotifyCharID:  guild.NotifyCharID,
		Notify:        guild.Notify,
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MY_GUILD_INFO_ACK(0, details))
}

// HandleMemberList envía la lista de todos los miembros del gremio (PROTO_NC_GUILD_MEMBER_LIST_ACK 0x741B).
func (gm *GuildManager) HandleMemberList(client *network.Client) {
	charID := client.CharacterID
	if charID <= 0 {
		return
	}

	guild, _, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LIST_ACK(0, 0, 0, nil))
		return
	}

	members, err := gm.charRepo.GetGuildMembers(guild.ID)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LIST_ACK(0, 0, 0, nil))
		return
	}

	gm.mu.RLock()
	var entries []protocol.GuildMemberClientEntry
	for _, m := range members {
		_, isOnline := gm.onlinePlayers[m.CharID]
		entries = append(entries, protocol.GuildMemberClientEntry{
			Name:      m.CharName,
			Grade:     m.Grade,
			Flags:     m.Flags,
			ExpRatio:  m.ExpRatio,
			Intro:     m.Intro,
			LogonTime: m.LastLogin,
			IsOnline:  isOnline,
			ClassID:   m.Class,
			Level:     m.Level,
			IsParty:   0,
			Map:       m.Map,
		})
	}
	gm.mu.RUnlock()

	cnt := uint16(len(entries))
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LIST_ACK(cnt, 0, cnt, entries))
}

// HandleInvite procesa la invitación de un jugador a unirse al gremio (OpCodes 0x7409 / 0x740B).
func (gm *GuildManager) HandleInvite(client *network.Client, targetName string) {
	charID := client.CharacterID
	guild, myMember, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || myMember == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName, 1))
		return
	}

	// Solo Maestro (0) o SubMaestro (1) pueden invitar
	if myMember.Grade > 1 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName, 1))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(targetName))
	gm.mu.Lock()
	targetPlayer, isOnline := gm.onlineByName[targetKey]
	if !isOnline || targetPlayer == nil {
		gm.mu.Unlock()
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName, 1))
		return
	}

	// Registrar invitación pendiente (expira en 30 segundos)
	gm.pendingInvites[targetPlayer.CharID] = &GuildPendingInvite{
		GuildID:         guild.ID,
		GuildName:       guild.Name,
		InviterCharID:   charID,
		InviterCharName: myMember.CharName,
		ExpiresAt:       time.Now().Add(30 * time.Second),
	}
	gm.mu.Unlock()

	// Enviar solicitud emergente al jugador invitado
	_ = targetPlayer.Client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_JOIN_REQ(guild.Name, myMember.CharName))
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName, 0))
}

// HandleJoinAck procesa la respuesta del jugador a la invitación de gremio (OpCode 0x740C).
func (gm *GuildManager) HandleJoinAck(client *network.Client, guildName string, joinAck byte) {
	charID := client.CharacterID
	gm.mu.Lock()
	invite, ok := gm.pendingInvites[charID]
	delete(gm.pendingInvites, charID)
	gm.mu.Unlock()

	if !ok || invite == nil || time.Now().After(invite.ExpiresAt) {
		return
	}

	if joinAck != 1 {
		log.Printf("[Guild] Target char %d rejected invitation to guild '%s'", charID, guildName)
		return
	}

	// Agregar a tGuildMember (Grade 4 = Newbie)
	err := gm.charRepo.AddGuildMember(invite.GuildID, charID, 4)
	if err != nil {
		log.Printf("[Guild] Error joining guild: %v", err)
		return
	}

	// Cargar datos del personaje para difusión
	ch, err := gm.charRepo.GetCharacterByID(charID)
	if err != nil || ch == nil {
		return
	}

	memberEntry := protocol.GuildMemberClientEntry{
		Name:      ch.Name,
		Grade:     4,
		Flags:     0,
		ExpRatio:  0,
		Intro:     "",
		LogonTime: time.Now(),
		IsOnline:  true,
		ClassID:   ch.Job,
		Level:     ch.Level,
		IsParty:   0,
		Map:       strconv.Itoa(int(ch.Map)),
	}

	// Anunciar a todo el gremio la incorporación del nuevo miembro
	gm.BroadcastToGuild(invite.GuildID, protocol.New_PROTO_NC_GUILD_MEMBER_JOIN_CMD(memberEntry), 0)

	// Actualizar al nuevo miembro
	gm.HandleMyGuildInfo(client)
	gm.HandleMemberList(client)
}

// HandleLeave procesa la salida voluntaria de un miembro del gremio (OpCode 0x741C).
func (gm *GuildManager) HandleLeave(client *network.Client) {
	charID := client.CharacterID
	guild, member, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || member == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(1))
		return
	}

	// El maestro no puede abandonar sin transferir o disolver
	if member.Grade == 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(1))
		return
	}

	err = gm.charRepo.RemoveGuildMember(guild.ID, charID)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(0))
	gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_MEMBER_LEAVE_CMD(member.CharName), charID)
}

// HandleVanish expulsa a un miembro del gremio (OpCode 0x740E).
func (gm *GuildManager) HandleVanish(client *network.Client, targetName string) {
	charID := client.CharacterID
	guild, myMember, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || myMember == nil || myMember.Grade > 1 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(targetName, 1))
		return
	}

	targetCh, err := gm.charRepo.GetCharacterByName(targetName)
	if err != nil || targetCh == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(targetName, 1))
		return
	}

	err = gm.charRepo.RemoveGuildMember(guild.ID, targetCh.ID)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(targetName, 1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(targetName, 0))
	gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_MEMBER_VANISH_CMD(targetName), 0)
}

// HandleGrade cambia el rango de un miembro (OpCode 0x7416).
func (gm *GuildManager) HandleGrade(client *network.Client, targetName string, newGrade byte) {
	charID := client.CharacterID
	guild, myMember, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || myMember == nil || myMember.Grade != 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(targetName, newGrade, 1))
		return
	}

	targetCh, err := gm.charRepo.GetCharacterByName(targetName)
	if err != nil || targetCh == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(targetName, newGrade, 1))
		return
	}

	err = gm.charRepo.SetGuildMemberGrade(guild.ID, targetCh.ID, newGrade)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(targetName, newGrade, 1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(targetName, newGrade, 0))
	gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_CMD(targetName, newGrade), 0)
}

// HandleNotify actualiza el mensaje de anuncio del gremio (OpCode 0x7410).
func (gm *GuildManager) HandleNotify(client *network.Client, notify string) {
	charID := client.CharacterID
	guild, myMember, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || myMember == nil || myMember.Grade > 1 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_NOTIFY_ACK(1))
		return
	}

	err = gm.charRepo.SetGuildNotify(guild.ID, myMember.CharName, notify)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_NOTIFY_ACK(1))
		return
	}

	now := time.Now()
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_NOTIFY_ACK(0))
	gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_NOTIFY_CMD(now, myMember.CharName, notify), 0)
}

// HandleDismiss disuelve el gremio por orden de su Maestro (OpCode 0x7407).
func (gm *GuildManager) HandleDismiss(client *network.Client) {
	charID := client.CharacterID
	guild, myMember, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || myMember == nil || myMember.Grade != 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_DISMISS_ACK(1))
		return
	}

	err = gm.charRepo.DismissGuild(guild.ID)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_DISMISS_ACK(1))
		return
	}

	now := time.Now()
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_DISMISS_ACK(0))
	gm.BroadcastToGuild(guild.ID, protocol.New_PROTO_NC_GUILD_DISMISS_CMD(guild.Name, 1, now), 0)
}

// HandleChat transmite un mensaje por el canal oficial de chat del gremio (OpCode 0x7473).
func (gm *GuildManager) HandleChat(client *network.Client, msg string) {
	charID := client.CharacterID
	guild, member, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil || member == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_CHAT_ACK(1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_CHAT_ACK(0))
	packet := protocol.New_PROTO_NC_GUILD_CHAT_CMD(uint32(guild.ID), member.CharName, msg)
	gm.BroadcastToGuild(guild.ID, packet, 0)
}

// ============================================================================
// HANDLERS OFICIALES DE ACADEMIA DE GREMIO (DEPARTAMENTO 38)
// ============================================================================

// HandleAcademyList retorna el directorio público de academias (OpCode 0x980B).
func (gm *GuildManager) HandleAcademyList(client *network.Client, partMask byte) {
	guilds, err := gm.charRepo.GetAllGuilds()
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_LIST_ACK(1, partMask, 0, 0, nil))
		return
	}

	var entries []protocol.AcademyListItemEntry
	for _, g := range guilds {
		students, _ := gm.charRepo.GetAcademyMembers(int(g.GuildID))
		entries = append(entries, protocol.AcademyListItemEntry{
			GuildName:           g.GuildName,
			GuildMasterName:     g.MasterCharName,
			NumOfGuildMembers:   g.MemberCount,
			NumOfAcademyMembers: uint16(len(students)),
			NumOfGraduate:       0,
		})
	}

	cnt := uint32(len(entries))
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_LIST_ACK(0, partMask, cnt, cnt, entries))
}

// HandleAcademyMyInfo envía la información detallada de la academia del gremio (OpCode 0x9807).
func (gm *GuildManager) HandleAcademyMyInfo(client *network.Client) {
	charID := client.CharacterID
	guild, _, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK(1, nil))
		return
	}

	academy, err := gm.charRepo.GetGuildAcademy(guild.ID)
	if err != nil || academy == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK(1, nil))
		return
	}

	students, _ := gm.charRepo.GetAcademyMembers(guild.ID)
	details := &protocol.AcademyClientDetails{
		AcademyMasterName: academy.MasterCharName,
		NumAcademyMembers: uint16(len(students)),
		MaxAcademyMembers: academy.MaxMember,
		AcademyPoint:      academy.Point,
		AcademyRank:       1,
		BuffUntilTime:     int32(academy.BuffEndTime.Unix()),
		Intro:             academy.Intro,
		NotifyDate:        academy.NotifyDate,
		NotifyCharID:      academy.NotifyCharID,
		Notify:            academy.Notify,
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK(0, details))
}

// HandleAcademyMemberList envía los aprendices inscritos en la academia (OpCode 0x980D).
func (gm *GuildManager) HandleAcademyMemberList(client *network.Client) {
	charID := client.CharacterID
	guild, _, err := gm.charRepo.GetGuildByCharID(charID)
	if err != nil || guild == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK(0, 0, 0, nil))
		return
	}

	students, err := gm.charRepo.GetAcademyMembers(guild.ID)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK(0, 0, 0, nil))
		return
	}

	gm.mu.RLock()
	var entries []protocol.AcademyMemberClientEntry
	for _, s := range students {
		_, isOnline := gm.onlinePlayers[s.CharID]
		entries = append(entries, protocol.AcademyMemberClientEntry{
			Name:      s.CharName,
			Grade:     s.Grade,
			Intro:     "",
			LogonTime: s.LastLogin,
			IsOnline:  isOnline,
			ClassID:   s.Class,
			Level:     s.Level,
			IsParty:   0,
			Map:       s.Map,
			JoinTime:  s.JoinDate,
			JoinLevel: s.JoinLevel,
		})
	}
	gm.mu.RUnlock()

	cnt := uint16(len(entries))
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK(cnt, 0, cnt, entries))
}

// HandleAcademyJoin inscribe a un jugador a la academia de un gremio (OpCode 0x9811).
func (gm *GuildManager) HandleAcademyJoin(client *network.Client, guildName string) {
	charID := client.CharacterID
	ch, err := gm.charRepo.GetCharacterByID(charID)
	if err != nil || ch == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName, 1))
		return
	}

	guilds, err := gm.charRepo.GetAllGuilds()
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName, 1))
		return
	}

	var targetGuild *database.GuildListItem
	for _, g := range guilds {
		if strings.EqualFold(g.GuildName, strings.TrimSpace(guildName)) {
			targetGuild = g
			break
		}
	}

	if targetGuild == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName, 1))
		return
	}

	err = gm.charRepo.AddAcademyMember(int(targetGuild.GuildID), charID, ch.Level)
	if err != nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName, 1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName, 0))

	studentEntry := protocol.AcademyMemberClientEntry{
		Name:      ch.Name,
		Grade:     0,
		Intro:     "",
		LogonTime: time.Now(),
		IsOnline:  true,
		ClassID:   ch.Job,
		Level:     ch.Level,
		IsParty:   0,
		Map:       strconv.Itoa(int(ch.Map)),
		JoinTime:  time.Now(),
		JoinLevel: ch.Level,
	}

	gm.BroadcastToAcademy(int(targetGuild.GuildID), protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD(studentEntry), 0)
}

// HandleAcademyLeave retira a un alumno de la academia (OpCode 0x981B).
func (gm *GuildManager) HandleAcademyLeave(client *network.Client) {
	charID := client.CharacterID
	ch, err := gm.charRepo.GetCharacterByID(charID)
	if err != nil || ch == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK(1))
		return
	}

	academy, _, _ := gm.charRepo.GetAcademyByCharID(charID)
	if academy == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK(1))
		return
	}

	_ = gm.charRepo.RemoveAcademyMember(academy.GuildID, charID)
	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK(0))
	gm.BroadcastToAcademy(academy.GuildID, protocol.New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD(ch.Name), charID)
}

// HandleAcademyChat difunde un mensaje por el chat de la academia (OpCode 0x9868).
func (gm *GuildManager) HandleAcademyChat(client *network.Client, msg string) {
	charID := client.CharacterID
	ch, err := gm.charRepo.GetCharacterByID(charID)
	if err != nil || ch == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(1))
		return
	}

	academy, _, _ := gm.charRepo.GetAcademyByCharID(charID)
	if academy == nil {
		// También permitir al maestro o miembros del gremio hablar en el chat de academia
		guild, _, _ := gm.charRepo.GetGuildByCharID(charID)
		if guild != nil {
			_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(0))
			packet := protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_CMD(uint32(guild.ID), ch.Name, msg)
			gm.BroadcastToAcademy(guild.ID, packet, 0)
			return
		}
		_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(1))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(0))
	packet := protocol.New_PROTO_NC_GUILD_ACADEMY_CHAT_CMD(uint32(academy.GuildID), ch.Name, msg)
	gm.BroadcastToAcademy(academy.GuildID, packet, 0)
}
