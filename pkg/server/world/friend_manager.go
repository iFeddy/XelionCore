package world

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// FriendManager administra las interacciones sociales y persistencia de amigos en WorldService.
type FriendManager struct {
	mu             sync.RWMutex
	charRepo       database.CharacterRepository
	partyMgr       *PartyManager
	pendingInvites map[int]int // TargetCharID -> RequesterCharID
}

// NewFriendManager inicializa el gestor de amigos.
func NewFriendManager(repo database.CharacterRepository, pm *PartyManager) *FriendManager {
	return &FriendManager{
		charRepo:       repo,
		partyMgr:       pm,
		pendingInvites: make(map[int]int),
	}
}

// SendFriendList despacha la lista oficial de amigos del jugador (PROTO_NC_FRIEND_LIST_CMD 0x5407).
func (fm *FriendManager) SendFriendList(client *network.Client) {
	charID := client.CharacterID
	if charID <= 0 {
		return
	}

	friends, err := fm.charRepo.GetFriends(charID)
	if err != nil {
		log.Printf("[Friend][%s] Error loading friends for char %d: %v", client.RemoteAddr(), charID, err)
		return
	}

	fm.partyMgr.mu.RLock()
	entries := make([]protocol.FriendEntry, 0, len(friends))
	for _, f := range friends {
		isOnline := false
		currentMap := f.Map
		if onlineMem, exists := fm.partyMgr.onlineChars[strings.ToLower(f.CharName)]; exists && onlineMem.Client != nil {
			isOnline = true
		}
		if currentMap == "" {
			currentMap = "RouN"
		}
		entries = append(entries, protocol.FriendEntry{
			CharName: f.CharName,
			IsOnline: isOnline,
			Map:      currentMap,
			Class:    f.Class,
			Level:    f.Level,
		})
	}
	fm.partyMgr.mu.RUnlock()

	log.Printf("[Friend][%s] Sent %d friend(s) to '%s'", client.RemoteAddr(), len(entries), client.Username)
	_ = client.SendPacket(protocol.New_PROTO_NC_FRIEND_LIST_CMD(entries))
}

// RequestFriend procesa la solicitud de agregar amigo (0x5401 PROTO_NC_FRIEND_SET_REQ).
func (fm *FriendManager) RequestFriend(requester *network.Client, targetName string) {
	fm.partyMgr.mu.RLock()
	targetMem, exists := fm.partyMgr.onlineChars[strings.ToLower(strings.TrimSpace(targetName))]
	fm.partyMgr.mu.RUnlock()

	if !exists || targetMem.Client == nil {
		log.Printf("[Friend] Target friend '%s' is not online", targetName)
		_ = requester.SendPacket(protocol.New_PROTO_NC_FRIEND_SET_ACK(targetName, 1)) // 1: Not found / offline
		return
	}

	if targetMem.CharID == requester.CharacterID {
		log.Printf("[Friend] Cannot add oneself as friend")
		return
	}

	// Registrar solicitud pendiente
	fm.mu.Lock()
	fm.pendingInvites[targetMem.CharID] = requester.CharacterID
	fm.mu.Unlock()

	// Obtener nombre del solicitante
	requesterName := requester.Username
	fm.partyMgr.mu.RLock()
	for _, m := range fm.partyMgr.onlineChars {
		if m.CharID == requester.CharacterID {
			requesterName = m.CharName
			break
		}
	}
	fm.partyMgr.mu.RUnlock()

	log.Printf("[Friend] '%s' requested friendship with '%s'", requesterName, targetMem.CharName)
	_ = targetMem.Client.SendPacket(protocol.New_PROTO_NC_FRIEND_SET_CONFIRM_REQ(requesterName))
}

// AcceptFriend procesa la aceptación de solicitud (0x5404 PROTO_NC_FRIEND_SET_CONFIRM_ACK).
func (fm *FriendManager) AcceptFriend(target *network.Client) {
	fm.mu.Lock()
	requesterID, hasInvite := fm.pendingInvites[target.CharacterID]
	if hasInvite {
		delete(fm.pendingInvites, target.CharacterID)
	}
	fm.mu.Unlock()

	if !hasInvite {
		return
	}

	targetID := target.CharacterID

	// Persistir relación en base de datos tFriend (p_Friend_Set)
	if err := fm.charRepo.AddFriend(requesterID, targetID); err != nil {
		log.Printf("[Friend] Error persisting friendship between %d and %d: %v", requesterID, targetID, err)
		return
	}

	log.Printf("[Friend] Friendship established between char %d and char %d", requesterID, targetID)

	// Refrescar lista de ambos
	fm.SendFriendList(target)

	fm.partyMgr.mu.RLock()
	for _, m := range fm.partyMgr.onlineChars {
		if m.CharID == requesterID && m.Client != nil {
			fm.SendFriendList(m.Client)
			break
		}
	}
	fm.partyMgr.mu.RUnlock()
}

// DeleteFriend procesa la eliminación de un amigo (0x5405 PROTO_NC_FRIEND_DEL_REQ).
func (fm *FriendManager) DeleteFriend(client *network.Client, friendName string) {
	charID := client.CharacterID
	friendChar, err := fm.charRepo.GetCharacterByName(friendName)
	if err != nil || friendChar == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_FRIEND_DEL_ACK(friendName, 1))
		return
	}

	if err := fm.charRepo.RemoveFriend(charID, friendChar.ID); err != nil {
		log.Printf("[Friend] Error removing friend %s for char %d: %v", friendName, charID, err)
		_ = client.SendPacket(protocol.New_PROTO_NC_FRIEND_DEL_ACK(friendName, 1))
		return
	}

	log.Printf("[Friend] '%s' removed friend '%s'", client.Username, friendName)
	_ = client.SendPacket(protocol.New_PROTO_NC_FRIEND_DEL_ACK(friendName, 0))

	// Notificar a ambos clientes actualizando su lista
	fm.SendFriendList(client)

	fm.partyMgr.mu.RLock()
	if fMem, exists := fm.partyMgr.onlineChars[strings.ToLower(friendName)]; exists && fMem.Client != nil {
		fm.SendFriendList(fMem.Client)
	}
	fm.partyMgr.mu.RUnlock()
}

// SendFriendPoints responde la consulta de puntos de amistad (0x541D PROTO_NC_FRIEND_POINT_REQ).
func (fm *FriendManager) SendFriendPoints(client *network.Client) {
	charID := client.CharacterID
	points, err := fm.charRepo.GetFriendPoints(charID)
	if err != nil {
		points = 0
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_FRIEND_POINT_ACK(uint32(points)))
}

// BroadcastLogin notifica a todos los amigos que el jugador se conectó (0x5409).
func (fm *FriendManager) BroadcastLogin(client *network.Client, charName, mapName string) {
	friends, err := fm.charRepo.GetFriends(client.CharacterID)
	if err != nil {
		return
	}

	cmd := protocol.New_PROTO_NC_FRIEND_LOGIN_CMD(charName, mapName)
	fm.partyMgr.mu.RLock()
	defer fm.partyMgr.mu.RUnlock()

	for _, f := range friends {
		if onlineMem, exists := fm.partyMgr.onlineChars[strings.ToLower(f.CharName)]; exists && onlineMem.Client != nil {
			_ = onlineMem.Client.SendPacket(cmd)
		}
	}
}

// BroadcastLogout notifica a todos los amigos que el jugador se desconectó (0x540A).
func (fm *FriendManager) BroadcastLogout(charID int, charName string) {
	friends, err := fm.charRepo.GetFriends(charID)
	if err != nil {
		return
	}

	cmd := protocol.New_PROTO_NC_FRIEND_LOGOUT_CMD(charName)
	fm.partyMgr.mu.RLock()
	defer fm.partyMgr.mu.RUnlock()

	for _, f := range friends {
		if onlineMem, exists := fm.partyMgr.onlineChars[strings.ToLower(f.CharName)]; exists && onlineMem.Client != nil {
			_ = onlineMem.Client.SendPacket(cmd)
		}
	}
}
