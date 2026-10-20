package world

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// PartyMember representa un jugador dentro de un grupo oficial.
type PartyMember struct {
	CharID    int
	CharName  string
	Class     byte
	Level     byte
	IsMaster  bool
	Client    *network.Client
}

// Party representa un grupo de Fiesta Online administrado por WorldManager (CParty).
type Party struct {
	ID          uint16
	MasterName  string
	MasterID    int
	LootingType byte // 0: Free-for-all, 1: Round-Robin, 2: Leader-Loot
	ItemGrade   byte
	Members     []*PartyMember
}

// PartyFinderAd representa un anuncio en el buscador de grupos (0x3850).
type PartyFinderAd struct {
	MasterName string
	PartyID    uint16
	Level      byte
	Purpose    string
}

// PartyManager administra la coordinación global de grupos en WorldService.
type PartyManager struct {
	mu            sync.RWMutex
	nextPartyID   uint16
	parties       map[uint16]*Party
	charToParty   map[int]uint16
	onlineChars   map[string]*PartyMember // CharName (lowercase) -> PartyMember
	pendingInvites map[int]int            // TargetCharID -> InviterCharID
	finderAds     map[uint16]*PartyFinderAd
}

// NewPartyManager inicializa el administrador de grupos.
func NewPartyManager() *PartyManager {
	return &PartyManager{
		nextPartyID:    1,
		parties:        make(map[uint16]*Party),
		charToParty:    make(map[int]uint16),
		onlineChars:    make(map[string]*PartyMember),
		pendingInvites: make(map[int]int),
		finderAds:      make(map[uint16]*PartyFinderAd),
	}
}

// RegisterOnline registra la sesión de un personaje en WorldManager.
func (pm *PartyManager) RegisterOnline(client *network.Client, charID int, charName string, classID, level byte) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(charName))
	pm.onlineChars[key] = &PartyMember{
		CharID:   charID,
		CharName: charName,
		Class:    classID,
		Level:    level,
		Client:   client,
	}
}

// UnregisterOnline remueve al personaje de la sesión y gestiona su salida del grupo si correspondía.
func (pm *PartyManager) UnregisterOnline(charID int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var charName string
	for k, m := range pm.onlineChars {
		if m.CharID == charID {
			charName = k
			delete(pm.onlineChars, k)
			break
		}
	}
	delete(pm.pendingInvites, charID)

	partyID, inParty := pm.charToParty[charID]
	if !inParty {
		return
	}

	delete(pm.charToParty, charID)
	party, exists := pm.parties[partyID]
	if !exists {
		return
	}

	// Remover miembro
	newMembers := make([]*PartyMember, 0, len(party.Members)-1)
	for _, m := range party.Members {
		if m.CharID != charID {
			newMembers = append(newMembers, m)
		}
	}
	party.Members = newMembers

	// Si queda 1 o 0 miembros, disolver grupo
	if len(party.Members) <= 1 {
		pm.dissolvePartyLocked(party)
		return
	}

	// Si el que salió era el líder, asignar el siguiente como líder
	if party.MasterID == charID && len(party.Members) > 0 {
		party.Members[0].IsMaster = true
		party.MasterID = party.Members[0].CharID
		party.MasterName = party.Members[0].CharName
	}

	pm.broadcastMemberListLocked(party)
	_ = charName
}

// InviteMember procesa la solicitud de invitación a grupo (0x3802 PROTO_NC_PARTY_JOIN_REQ).
func (pm *PartyManager) InviteMember(inviter *network.Client, targetName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	inviterID := inviter.CharacterID
	targetKey := strings.ToLower(strings.TrimSpace(targetName))
	target, exists := pm.onlineChars[targetKey]
	if !exists || target.Client == nil {
		log.Printf("[Party] Target '%s' is not online", targetName)
		_ = inviter.SendPacket(protocol.New_PROTO_NC_PARTY_JOIN_ACK(targetName, 0, 1)) // 1: Target offline / not found
		return
	}

	if target.CharID == inviterID {
		log.Printf("[Party] Cannot invite oneself")
		return
	}

	// Verificar si el invitado ya está en un grupo
	if _, inParty := pm.charToParty[target.CharID]; inParty {
		log.Printf("[Party] Target '%s' is already in a party", targetName)
		_ = inviter.SendPacket(protocol.New_PROTO_NC_PARTY_JOIN_ACK(targetName, target.Class, 2)) // 2: Already in party
		return
	}

	// Verificar grupo del invitador
	var party *Party
	partyID, inParty := pm.charToParty[inviterID]
	if inParty {
		party = pm.parties[partyID]
		if party.MasterID != inviterID {
			log.Printf("[Party] Inviter %d is not party master", inviterID)
			return
		}
		if len(party.Members) >= 5 {
			log.Printf("[Party] Party %d is full", partyID)
			_ = inviter.SendPacket(protocol.New_PROTO_NC_PARTY_JOIN_ACK(targetName, target.Class, 3)) // 3: Party is full
			return
		}
	}

	// Registrar invitación pendiente
	pm.pendingInvites[target.CharID] = inviterID

	// Obtener nombre del invitador
	var inviterName string
	for _, m := range pm.onlineChars {
		if m.CharID == inviterID {
			inviterName = m.CharName
			break
		}
	}
	if inviterName == "" {
		inviterName = inviter.Username
	}

	log.Printf("[Party] '%s' invited '%s' to party", inviterName, target.CharName)
	_ = target.Client.SendPacket(protocol.New_PROTO_NC_PARTY_JOINPROPOSE_REQ(inviterName))
}

// AcceptInvite procesa la aceptación de invitación (0x3804 PROTO_NC_PARTY_JOINPROPOSE_ALLOW_ACK).
func (pm *PartyManager) AcceptInvite(target *network.Client) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	targetID := target.CharacterID
	inviterID, hasInvite := pm.pendingInvites[targetID]
	if !hasInvite {
		return
	}
	delete(pm.pendingInvites, targetID)

	var targetMember *PartyMember
	for _, m := range pm.onlineChars {
		if m.CharID == targetID {
			targetMember = m
			break
		}
	}
	if targetMember == nil {
		return
	}

	var party *Party
	partyID, inParty := pm.charToParty[inviterID]
	if !inParty {
		// Crear nuevo grupo
		var inviterMember *PartyMember
		for _, m := range pm.onlineChars {
			if m.CharID == inviterID {
				inviterMember = m
				break
			}
		}
		if inviterMember == nil {
			return
		}

		newID := pm.nextPartyID
		pm.nextPartyID++

		inviterMember.IsMaster = true
		targetMember.IsMaster = false

		party = &Party{
			ID:         newID,
			MasterName: inviterMember.CharName,
			MasterID:   inviterMember.CharID,
			Members:    []*PartyMember{inviterMember, targetMember},
		}

		pm.parties[newID] = party
		pm.charToParty[inviterID] = newID
		pm.charToParty[targetID] = newID
	} else {
		party = pm.parties[partyID]
		if len(party.Members) >= 5 {
			return
		}
		targetMember.IsMaster = false
		party.Members = append(party.Members, targetMember)
		pm.charToParty[targetID] = partyID
	}

	log.Printf("[Party] '%s' joined party %d (Total: %d members)",
		targetMember.CharName, party.ID, len(party.Members))

	// Notificar a todos los miembros con la lista actualizada
	pm.broadcastMemberListLocked(party)
}

// RejectInvite procesa el rechazo de invitación (0x3805 PROTO_NC_PARTY_JOINPROPOSE_REJECT_ACK).
func (pm *PartyManager) RejectInvite(target *network.Client) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	targetID := target.CharacterID
	inviterID, hasInvite := pm.pendingInvites[targetID]
	if !hasInvite {
		return
	}
	delete(pm.pendingInvites, targetID)

	for _, m := range pm.onlineChars {
		if m.CharID == inviterID && m.Client != nil {
			_ = m.Client.SendPacket(protocol.New_PROTO_NC_PARTY_JOIN_ACK("", 0, 4)) // 4: Rejected
			break
		}
	}
}

// LeaveParty procesa la salida voluntaria de un miembro (0x380A PROTO_NC_PARTY_LEAVE_REQ).
func (pm *PartyManager) LeaveParty(client *network.Client) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	charID := client.CharacterID
	partyID, inParty := pm.charToParty[charID]
	if !inParty {
		_ = client.SendPacket(protocol.New_PROTO_NC_PARTY_LEAVE_ACK(1))
		return
	}

	delete(pm.charToParty, charID)
	party, exists := pm.parties[partyID]
	if !exists {
		return
	}

	newMembers := make([]*PartyMember, 0, len(party.Members)-1)
	for _, m := range party.Members {
		if m.CharID != charID {
			newMembers = append(newMembers, m)
		}
	}
	party.Members = newMembers

	// Confirmación al que se retira
	_ = client.SendPacket(protocol.New_PROTO_NC_PARTY_LEAVE_ACK(0))

	// Si queda 1 o 0 miembros, disolver grupo
	if len(party.Members) <= 1 {
		pm.dissolvePartyLocked(party)
		return
	}

	// Si salió el líder, promover al primero
	if party.MasterID == charID && len(party.Members) > 0 {
		party.Members[0].IsMaster = true
		party.MasterID = party.Members[0].CharID
		party.MasterName = party.Members[0].CharName
	}

	pm.broadcastMemberListLocked(party)
}

// KickMember procesa la expulsión de un miembro por parte del líder (0x3814 PROTO_NC_PARTY_KICKOFF_REQ).
func (pm *PartyManager) KickMember(master *network.Client, targetName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	masterID := master.CharacterID
	partyID, inParty := pm.charToParty[masterID]
	if !inParty {
		return
	}
	party := pm.parties[partyID]
	if party.MasterID != masterID {
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(targetName))
	var kickedMember *PartyMember
	newMembers := make([]*PartyMember, 0, len(party.Members))
	for _, m := range party.Members {
		if strings.ToLower(m.CharName) == targetKey {
			kickedMember = m
		} else {
			newMembers = append(newMembers, m)
		}
	}

	if kickedMember == nil {
		_ = master.SendPacket(protocol.New_PROTO_NC_PARTY_KICKOFF_ACK(targetName, 1))
		return
	}

	delete(pm.charToParty, kickedMember.CharID)
	party.Members = newMembers

	if kickedMember.Client != nil {
		_ = kickedMember.Client.SendPacket(protocol.New_PROTO_NC_PARTY_LEAVE_ACK(0))
	}
	_ = master.SendPacket(protocol.New_PROTO_NC_PARTY_KICKOFF_ACK(targetName, 0))

	if len(party.Members) <= 1 {
		pm.dissolvePartyLocked(party)
		return
	}

	pm.broadcastMemberListLocked(party)
}

// ChangeMaster transfiere el liderazgo a otro miembro (0x3828 PROTO_NC_PARTY_CHANGEMASTER_REQ).
func (pm *PartyManager) ChangeMaster(master *network.Client, newMasterName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	masterID := master.CharacterID
	partyID, inParty := pm.charToParty[masterID]
	if !inParty {
		return
	}
	party := pm.parties[partyID]
	if party.MasterID != masterID {
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(newMasterName))
	var newMaster *PartyMember
	for _, m := range party.Members {
		if strings.ToLower(m.CharName) == targetKey {
			newMaster = m
			break
		}
	}

	if newMaster == nil {
		_ = master.SendPacket(protocol.New_PROTO_NC_PARTY_CHANGEMASTER_ACK(1))
		return
	}

	for _, m := range party.Members {
		m.IsMaster = (m.CharID == newMaster.CharID)
	}
	party.MasterID = newMaster.CharID
	party.MasterName = newMaster.CharName

	_ = master.SendPacket(protocol.New_PROTO_NC_PARTY_CHANGEMASTER_ACK(0))
	pm.broadcastMemberListLocked(party)
}

// SetLooting actualiza la regla de distribución de botín (0x384B PROTO_NC_PARTY_ITEM_LOOTING_SET).
func (pm *PartyManager) SetLooting(master *network.Client, lootType, itemGrade byte) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	masterID := master.CharacterID
	partyID, inParty := pm.charToParty[masterID]
	if !inParty {
		return
	}
	party := pm.parties[partyID]
	if party.MasterID != masterID {
		return
	}

	party.LootingType = lootType
	party.ItemGrade = itemGrade

	cmd := protocol.New_PROTO_NC_PARTY_ITEM_LOOTING_CMD(lootType, itemGrade)
	for _, m := range party.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(cmd)
		}
	}
}

// AddFinderAd publica un anuncio en el buscador de grupos (0x3850).
func (pm *PartyManager) AddFinderAd(client *network.Client, purpose string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	charID := client.CharacterID
	partyID, inParty := pm.charToParty[charID]
	if !inParty {
		return
	}
	party := pm.parties[partyID]
	if party.MasterID != charID {
		return
	}

	level := byte(1)
	for _, m := range party.Members {
		if m.CharID == charID {
			level = m.Level
			break
		}
	}

	pm.finderAds[partyID] = &PartyFinderAd{
		MasterName: party.MasterName,
		PartyID:    partyID,
		Level:      level,
		Purpose:    purpose,
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PARTY_FINDER_ADD_ACK)
	resp.WriteByte(0) // Success
	_ = client.SendPacket(resp)
}

// DeleteFinderAd elimina el anuncio del buscador de grupos (0x3852).
func (pm *PartyManager) DeleteFinderAd(client *network.Client) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	charID := client.CharacterID
	partyID, inParty := pm.charToParty[charID]
	if inParty {
		delete(pm.finderAds, partyID)
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PARTY_FINDER_DELETE_ACK)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

// ListFinderAds despacha todos los anuncios activos del buscador de grupos (0x3854).
func (pm *PartyManager) ListFinderAds(client *network.Client) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	entries := make([]protocol.PartyFinderEntry, 0, len(pm.finderAds))
	for partyID, ad := range pm.finderAds {
		party, exists := pm.parties[partyID]
		cnt := byte(1)
		if exists {
			cnt = byte(len(party.Members))
		}
		entries = append(entries, protocol.PartyFinderEntry{
			MasterName: ad.MasterName,
			Level:      ad.Level,
			MemberCnt:  cnt,
			Purpose:    ad.Purpose,
		})
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_PARTY_FINDER_LIST_ACK(entries))
}

func (pm *PartyManager) broadcastMemberListLocked(party *Party) {
	memberInfos := make([]protocol.PartyMemberInfo, 0, len(party.Members))
	for _, m := range party.Members {
		status := byte(0)
		if m.IsMaster {
			status = 1
		}
		memberInfos = append(memberInfos, protocol.PartyMemberInfo{
			MemberName:   m.CharName,
			MemberStatus: status,
			CharClass:    m.Class,
		})
	}

	cmd := protocol.New_PROTO_NC_PARTY_MEMBER_LIST_CMD(memberInfos)
	for _, m := range party.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(cmd)
		}
	}
}

func (pm *PartyManager) dissolvePartyLocked(party *Party) {
	dismissCmd := protocol.New_PROTO_NC_PARTY_DISMISS_CMD(party.ID)
	for _, m := range party.Members {
		delete(pm.charToParty, m.CharID)
		if m.Client != nil {
			_ = m.Client.SendPacket(dismissCmd)
		}
	}
	delete(pm.parties, party.ID)
	delete(pm.finderAds, party.ID)
	log.Printf("[Party] Party %d dissolved", party.ID)
}
