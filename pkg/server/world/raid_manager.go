package world

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// RaidMember representa un miembro dentro de una banda (Raid).
type RaidMember struct {
	CharID     int
	CharName   string
	Class      byte
	Level      byte
	PartyIndex byte // 0..3
	SlotIndex  byte // 0..4
	IsLeader   bool
	IsSubLeader bool
	Client     *network.Client
}

// Raid representa una banda oficial de hasta 20 jugadores (4 grupos x 5 miembros).
type Raid struct {
	ID          uint16
	LeaderID    int
	LeaderName  string
	LooterName  string
	SubLeaders  map[string]bool
	LootingType byte
	ItemGrade   byte
	PartyNos    [4]uint16
	Members     map[int]*RaidMember // CharID -> RaidMember
}

// RaidManager administra los grupos de banda (Dept 44) en WorldService.
type RaidManager struct {
	mu             sync.RWMutex
	charRepo       database.CharacterRepository
	partyMgr       *PartyManager
	nextRaidID     uint16
	raids          map[uint16]*Raid
	charToRaid     map[int]uint16
	pendingInvites map[int]uint16 // TargetCharID -> RaidID
}

// NewRaidManager inicializa el administrador de bandas de WorldService.
func NewRaidManager(charRepo database.CharacterRepository, partyMgr *PartyManager) *RaidManager {
	return &RaidManager{
		charRepo:       charRepo,
		partyMgr:       partyMgr,
		nextRaidID:     1,
		raids:          make(map[uint16]*Raid),
		charToRaid:     make(map[int]uint16),
		pendingInvites: make(map[int]uint16),
	}
}

// HandleMake crea una nueva banda a partir de una party o el líder (0xB001).
func (rm *RaidManager) HandleMake(client *network.Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	if _, inRaid := rm.charToRaid[charID]; inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_MAKE_ACK(0, 1))
		return
	}

	raidID := rm.nextRaidID
	rm.nextRaidID++

	raid := &Raid{
		ID:         raidID,
		LeaderID:   charID,
		LeaderName: client.Username,
		SubLeaders: make(map[string]bool),
		PartyNos:   [4]uint16{raidID * 10, raidID*10 + 1, raidID*10 + 2, raidID*10 + 3},
		Members:    make(map[int]*RaidMember),
	}

	leaderMember := &RaidMember{
		CharID:     charID,
		CharName:   client.Username,
		PartyIndex: 0,
		SlotIndex:  0,
		IsLeader:   true,
		Client:     client,
	}
	raid.Members[charID] = leaderMember
	rm.raids[raidID] = raid
	rm.charToRaid[charID] = raidID

	var firstPartyChars [5]string
	firstPartyChars[0] = client.Username

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_MAKE_ACK(raidID, 0))
	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_MAKE_BROAD_CMD(raidID, raid.PartyNos, firstPartyChars))
	log.Printf("[Raid] Raid %d created by '%s'", raidID, client.Username)
}

// HandleLootingType configura la regla de botín de la banda (0xB008).
func (rm *RaidManager) HandleLootingType(client *network.Client, itemGrade, lootingType byte) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_LOOTING_TYPE_ACK(1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_LOOTING_TYPE_ACK(2))
		return
	}

	raid.ItemGrade = itemGrade
	raid.LootingType = lootingType

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_LOOTING_TYPE_ACK(0))
	broadCmd := protocol.New_PROTO_NC_RAID_LOOTING_TYPE_BROAD_CMD(itemGrade, lootingType)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleInvitation invita a un jugador a unirse a la banda (0xB00C).
func (rm *RaidManager) HandleInvitation(client *network.Client, targetName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID && !raid.SubLeaders[strings.ToLower(client.Username)] {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 2))
		return
	}

	if len(raid.Members) >= 20 {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 3))
		return
	}

	// Buscar target en partyMgr online
	var targetClient *network.Client
	var targetID int
	if rm.partyMgr != nil {
		targetKey := strings.ToLower(strings.TrimSpace(targetName))
		if m, exists := rm.partyMgr.onlineChars[targetKey]; exists {
			targetClient = m.Client
			targetID = m.CharID
		}
	}

	if targetClient == nil || targetID == 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 4)) // Offline
		return
	}

	if _, alreadyInRaid := rm.charToRaid[targetID]; alreadyInRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 5)) // Already in raid
		return
	}

	rm.pendingInvites[targetID] = raidID
	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ACK(targetName, 0))
	_ = targetClient.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_CMD(client.Username))
}

// HandleInvitationAnswer procesa la respuesta a la invitación de banda (0xB00F).
func (rm *RaidManager) HandleInvitationAnswer(client *network.Client, isAllow byte, masterName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, hasInvite := rm.pendingInvites[charID]
	if !hasInvite {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ANSWER_ACK(1))
		return
	}
	delete(rm.pendingInvites, charID)

	if isAllow == 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ANSWER_ACK(0))
		return
	}

	raid, exists := rm.raids[raidID]
	if !exists || len(raid.Members) >= 20 {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ANSWER_ACK(2))
		return
	}

	partyIdx := byte(len(raid.Members) / 5)
	slotIdx := byte(len(raid.Members) % 5)

	newMember := &RaidMember{
		CharID:     charID,
		CharName:   client.Username,
		PartyIndex: partyIdx,
		SlotIndex:  slotIdx,
		Client:     client,
	}
	raid.Members[charID] = newMember
	rm.charToRaid[charID] = raidID

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_INVITATION_ANSWER_ACK(0))

	broadCmd := protocol.New_PROTO_NC_RAID_MEMBER_JOIN_BROAD_CMD(raid.PartyNos[partyIdx], partyIdx, slotIdx, client.Username)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
	log.Printf("[Raid] '%s' joined Raid %d (Slot %d:%d)", client.Username, raidID, partyIdx, slotIdx)
}

// HandleLeave remueve al personaje de la banda (0xB014).
func (rm *RaidManager) HandleLeave(client *network.Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_LEAVE_ACK(1))
		return
	}
	raid, exists := rm.raids[raidID]
	if !exists {
		return
	}

	member, hasMember := raid.Members[charID]
	if !hasMember {
		return
	}

	delete(raid.Members, charID)
	delete(rm.charToRaid, charID)

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_LEAVE_ACK(0))

	broadCmd := protocol.New_PROTO_NC_RAID_LEAVE_BROAD_CMD(raid.PartyNos[member.PartyIndex], member.PartyIndex, member.SlotIndex, member.CharName)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}

	if len(raid.Members) == 0 {
		delete(rm.raids, raidID)
		log.Printf("[Raid] Raid %d dissolved (empty)", raidID)
	} else if raid.LeaderID == charID {
		// Asignar primer miembro restante como nuevo líder
		for _, m := range raid.Members {
			raid.LeaderID = m.CharID
			raid.LeaderName = m.CharName
			m.IsLeader = true
			setLeaderBroad := protocol.New_PROTO_NC_RAID_SET_LEADER_BROAD_CMD("System", m.CharName, raid.PartyNos[m.PartyIndex], m.PartyIndex, m.SlotIndex)
			for _, om := range raid.Members {
				if om.Client != nil {
					_ = om.Client.SendPacket(setLeaderBroad)
				}
			}
			break
		}
	}
}

// HandleVanish expulsa a un miembro de la banda (0xB018).
func (rm *RaidManager) HandleVanish(client *network.Client, vanishName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_VANISH_ACK(vanishName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID && !raid.SubLeaders[strings.ToLower(client.Username)] {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_VANISH_ACK(vanishName, 2))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(vanishName))
	var targetMember *RaidMember
	for _, m := range raid.Members {
		if strings.ToLower(m.CharName) == targetKey {
			targetMember = m
			break
		}
	}

	if targetMember == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_VANISH_ACK(vanishName, 3))
		return
	}

	delete(raid.Members, targetMember.CharID)
	delete(rm.charToRaid, targetMember.CharID)

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_VANISH_ACK(targetMember.CharName, 0))
	if targetMember.Client != nil {
		_ = targetMember.Client.SendPacket(protocol.New_PROTO_NC_RAID_LEAVE_ACK(0))
	}

	broadCmd := protocol.New_PROTO_NC_RAID_VANISH_BROAD_CMD(raid.PartyNos[targetMember.PartyIndex], targetMember.PartyIndex, targetMember.SlotIndex, client.Username, targetMember.CharName)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleSlotMove intercambia o mueve miembros entre casillas de la banda (0xB027).
func (rm *RaidManager) HandleSlotMove(client *network.Client, prePos, postPos uint32, preChar, postChar string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SLOT_MOVE_ACK(preChar, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID && !raid.SubLeaders[strings.ToLower(client.Username)] {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SLOT_MOVE_ACK(preChar, 2))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SLOT_MOVE_ACK(preChar, 0))
	broadCmd := protocol.New_PROTO_NC_RAID_SLOT_MOVE_BROAD_CMD(client.Username, prePos, postPos, preChar, postChar)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleSetLeader transfiere el liderazgo de la banda (0xB02B).
func (rm *RaidManager) HandleSetLeader(client *network.Client, newLeaderName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LEADER_ACK(newLeaderName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LEADER_ACK(newLeaderName, 2))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(newLeaderName))
	var newLeader *RaidMember
	for _, m := range raid.Members {
		if strings.ToLower(m.CharName) == targetKey {
			newLeader = m
			break
		}
	}

	if newLeader == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LEADER_ACK(newLeaderName, 3))
		return
	}

	for _, m := range raid.Members {
		m.IsLeader = (m.CharID == newLeader.CharID)
	}
	raid.LeaderID = newLeader.CharID
	raid.LeaderName = newLeader.CharName

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LEADER_ACK(newLeader.CharName, 0))
	broadCmd := protocol.New_PROTO_NC_RAID_SET_LEADER_BROAD_CMD(client.Username, newLeader.CharName, raid.PartyNos[newLeader.PartyIndex], newLeader.PartyIndex, newLeader.SlotIndex)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleSetSubLeader asigna un subjefe de banda (0xB030).
func (rm *RaidManager) HandleSetSubLeader(client *network.Client, subLeaderName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_SUB_LEADER_ACK(subLeaderName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_SUB_LEADER_ACK(subLeaderName, 2))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(subLeaderName))
	var target *RaidMember
	for _, m := range raid.Members {
		if strings.ToLower(m.CharName) == targetKey {
			target = m
			break
		}
	}
	if target == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_SUB_LEADER_ACK(subLeaderName, 3))
		return
	}

	target.IsSubLeader = true
	raid.SubLeaders[targetKey] = true

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_SUB_LEADER_ACK(target.CharName, 0))
	broadCmd := protocol.New_PROTO_NC_RAID_SET_SUB_LEADER_BROAD_CMD(client.Username, target.CharName, raid.PartyNos[target.PartyIndex], target.PartyIndex, target.SlotIndex)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleDelSubLeader remueve a un subjefe de banda (0xB035).
func (rm *RaidManager) HandleDelSubLeader(client *network.Client, delSubLeaderName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_DEL_SUB_LEADER_ACK(delSubLeaderName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_DEL_SUB_LEADER_ACK(delSubLeaderName, 2))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(delSubLeaderName))
	var target *RaidMember
	for _, m := range raid.Members {
		if strings.ToLower(m.CharName) == targetKey {
			target = m
			break
		}
	}
	if target == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_DEL_SUB_LEADER_ACK(delSubLeaderName, 3))
		return
	}

	target.IsSubLeader = false
	delete(raid.SubLeaders, targetKey)

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_DEL_SUB_LEADER_ACK(target.CharName, 0))
	broadCmd := protocol.New_PROTO_NC_RAID_DEL_SUB_LEADER_BROAD_CMD(client.Username, target.CharName, raid.PartyNos[target.PartyIndex], target.PartyIndex, target.SlotIndex)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleSetLooter designa al recolector de banda (0xB03A).
func (rm *RaidManager) HandleSetLooter(client *network.Client, looterName string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LOOTER_ACK(looterName, 1))
		return
	}
	raid := rm.raids[raidID]
	if raid.LeaderID != charID {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LOOTER_ACK(looterName, 2))
		return
	}

	targetKey := strings.ToLower(strings.TrimSpace(looterName))
	var target *RaidMember
	for _, m := range raid.Members {
		if strings.ToLower(m.CharName) == targetKey {
			target = m
			break
		}
	}
	if target == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LOOTER_ACK(looterName, 3))
		return
	}

	raid.LooterName = target.CharName

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_SET_LOOTER_ACK(target.CharName, 0))
	broadCmd := protocol.New_PROTO_NC_RAID_SET_LOOTER_BROAD_CMD(client.Username, target.CharName, raid.PartyNos[target.PartyIndex], target.PartyIndex, target.SlotIndex)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleChat difunde un mensaje a todos los miembros de la banda (0xB03F).
func (rm *RaidManager) HandleChat(client *network.Client, chatText string) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_CHAT_ACK(1))
		return
	}
	raid, exists := rm.raids[raidID]
	if !exists {
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_CHAT_ACK(0))
	broadCmd := protocol.New_PROTO_NC_RAID_CHAT_BROAD_CMD(client.Username, chatText)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}

// HandleWarning emite un aviso destacado del líder a la banda (0xB042).
func (rm *RaidManager) HandleWarning(client *network.Client, warningText string) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	charID := client.CharacterID
	raidID, inRaid := rm.charToRaid[charID]
	if !inRaid {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_WARNING_ACK(1))
		return
	}
	raid, exists := rm.raids[raidID]
	if !exists {
		return
	}
	if raid.LeaderID != charID && !raid.SubLeaders[strings.ToLower(client.Username)] {
		_ = client.SendPacket(protocol.New_PROTO_NC_RAID_WARNING_ACK(2))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_RAID_WARNING_ACK(0))
	broadCmd := protocol.New_PROTO_NC_RAID_WARNING_BROAD_CMD(client.Username, warningText)
	for _, m := range raid.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(broadCmd)
		}
	}
}
