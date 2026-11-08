package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_RAID_CHAT_ACK (size = 2 bytes)
type PROTO_NC_RAID_CHAT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_CHAT_BROAD_CMD (size = 22 bytes)
type PROTO_NC_RAID_CHAT_BROAD_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	STalkerCharID                    Name5                // offset: 1, size: 20 bytes
	NChatLen                         uint8                // offset: 21, size: 1 bytes
	SChat                            [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_RAID_CHAT_REQ (size = 2 bytes)
type PROTO_NC_RAID_CHAT_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	NChatLen                         uint8                // offset: 1, size: 1 bytes
	SChat                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_RAID_DEL_SUB_LEADER_ACK (size = 22 bytes)
type PROTO_NC_RAID_DEL_SUB_LEADER_ACK_DATA struct {
	SDelSubLeaderCharID              Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_DEL_SUB_LEADER_BROAD_CMD (size = 44 bytes)
type PROTO_NC_RAID_DEL_SUB_LEADER_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	SDelSubLeaderCharID              Name5                // offset: 20, size: 20 bytes
	NPartyNo                         uint16               // offset: 40, size: 2 bytes
	NPartyIndex                      uint8                // offset: 42, size: 1 bytes
	NSlotIndex                       uint8                // offset: 43, size: 1 bytes
}

// PROTO_NC_RAID_DEL_SUB_LEADER_CMD (size = 20 bytes)
type PROTO_NC_RAID_DEL_SUB_LEADER_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_DEL_SUB_LEADER_REQ (size = 20 bytes)
type PROTO_NC_RAID_DEL_SUB_LEADER_REQ_DATA struct {
	SDelSubLeaderCharID              Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_DISMISS_BROAD_CMD (size = 1 bytes)
type PROTO_NC_RAID_DISMISS_BROAD_CMD_DATA struct {}

// PROTO_NC_RAID_INVITATION_ACK (size = 22 bytes)
type PROTO_NC_RAID_INVITATION_ACK_DATA struct {
	SInvitationCharID                Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_INVITATION_ANSWER_ACK (size = 2 bytes)
type PROTO_NC_RAID_INVITATION_ANSWER_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_INVITATION_ANSWER_CMD (size = 21 bytes)
type PROTO_NC_RAID_INVITATION_ANSWER_CMD_DATA struct {
	SInvitationCharID                Name5                // offset: 0, size: 20 bytes
	BIsAllow                         uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_RAID_INVITATION_ANSWER_REQ (size = 21 bytes)
type PROTO_NC_RAID_INVITATION_ANSWER_REQ_DATA struct {
	BIsAllow                         uint8                // offset: 0, size: 1 bytes
	SMasterCharID                    Name5                // offset: 1, size: 20 bytes
}

// PROTO_NC_RAID_INVITATION_CMD (size = 20 bytes)
type PROTO_NC_RAID_INVITATION_CMD_DATA struct {
	SMasterCharID                    Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_INVITATION_REQ (size = 20 bytes)
type PROTO_NC_RAID_INVITATION_REQ_DATA struct {
	SInvitationCharID                Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_LEAVE_ACK (size = 2 bytes)
type PROTO_NC_RAID_LEAVE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_LEAVE_BROAD_CMD (size = 24 bytes)
type PROTO_NC_RAID_LEAVE_BROAD_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	SCharID                          Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_RAID_LOGININFO_CMD (size = 24 bytes)
type PROTO_NC_RAID_LOGININFO_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	SCharID                          Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_RAID_LOGIN_CMD (size = 6 bytes)
type PROTO_NC_RAID_LOGIN_CMD struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_RAID_LOGOUTINFO_CMD (size = 24 bytes)
type PROTO_NC_RAID_LOGOUTINFO_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	SCharID                          Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_RAID_LOGOUT_CMD (size = 6 bytes)
type PROTO_NC_RAID_LOGOUT_CMD struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_RAID_LOOTING_TYPE_ACK (size = 2 bytes)
type PROTO_NC_RAID_LOOTING_TYPE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_LOOTING_TYPE_BROAD_CMD (size = 2 bytes)
type PROTO_NC_RAID_LOOTING_TYPE_BROAD_CMD_DATA struct {
	NItemGrade                       uint8                // offset: 0, size: 1 bytes
	NLootingType                     uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_RAID_LOOTING_TYPE_REQ (size = 2 bytes)
type PROTO_NC_RAID_LOOTING_TYPE_REQ_DATA struct {
	NItemGrade                       uint8                // offset: 0, size: 1 bytes
	NLootingType                     uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_RAID_MAKE_ACK (size = 4 bytes)
type PROTO_NC_RAID_MAKE_ACK_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NError                           uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_RAID_MAKE_BROAD_CMD (size = 110 bytes)
type PROTO_NC_RAID_MAKE_BROAD_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NRaidPartyNo                     [4]uint16            // offset: 2, size: 8 bytes
	SFirstPartyCharID                [100]byte            // offset: 10, size: 100 bytes
}

// PROTO_NC_RAID_MAKE_REQ (size = 1 bytes)
type PROTO_NC_RAID_MAKE_REQ_DATA struct {}

// PROTO_NC_RAID_MEMBER_JOIN_BROAD_CMD (size = 24 bytes)
type PROTO_NC_RAID_MEMBER_JOIN_BROAD_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	SCharID                          Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_RAID_MEMBER_LIST_CMD (size = 542 bytes)
type PROTO_NC_RAID_MEMBER_LIST_CMD_DATA struct {
	Leaderid                         Name5                // offset: 0, size: 20 bytes
	SubLeaderid                      [80]byte             // offset: 20, size: 80 bytes
	Looterid                         Name5                // offset: 100, size: 20 bytes
	NItemGrade                       uint8                // offset: 120, size: 1 bytes
	NLootingType                     uint8                // offset: 121, size: 1 bytes
	Memberarray                      [420]byte            // offset: 122, size: 420 bytes
}

// PROTO_NC_RAID_SET_LEADER_ACK (size = 22 bytes)
type PROTO_NC_RAID_SET_LEADER_ACK_DATA struct {
	SNewLeaderCharID                 Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_SET_LEADER_BROAD_CMD (size = 44 bytes)
type PROTO_NC_RAID_SET_LEADER_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	SNewLeaderCharID                 Name5                // offset: 20, size: 20 bytes
	NPartyNo                         uint16               // offset: 40, size: 2 bytes
	NPartyIndex                      uint8                // offset: 42, size: 1 bytes
	NSlotIndex                       uint8                // offset: 43, size: 1 bytes
}

// PROTO_NC_RAID_SET_LEADER_CMD (size = 20 bytes)
type PROTO_NC_RAID_SET_LEADER_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SET_LEADER_REQ (size = 20 bytes)
type PROTO_NC_RAID_SET_LEADER_REQ_DATA struct {
	SNewLeaderCharID                 Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SET_LOOTER_ACK (size = 22 bytes)
type PROTO_NC_RAID_SET_LOOTER_ACK_DATA struct {
	SNewLooterCharID                 Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_SET_LOOTER_BROAD_CMD (size = 44 bytes)
type PROTO_NC_RAID_SET_LOOTER_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	SNewLooterCharID                 Name5                // offset: 20, size: 20 bytes
	NPartyNo                         uint16               // offset: 40, size: 2 bytes
	NPartyIndex                      uint8                // offset: 42, size: 1 bytes
	NSlotIndex                       uint8                // offset: 43, size: 1 bytes
}

// PROTO_NC_RAID_SET_LOOTER_CMD (size = 20 bytes)
type PROTO_NC_RAID_SET_LOOTER_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SET_LOOTER_REQ (size = 20 bytes)
type PROTO_NC_RAID_SET_LOOTER_REQ_DATA struct {
	SNewLooterCharID                 Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SET_SUB_LEADER_ACK (size = 22 bytes)
type PROTO_NC_RAID_SET_SUB_LEADER_ACK_DATA struct {
	SNewSubLeaderCharID              Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_SET_SUB_LEADER_BROAD_CMD (size = 44 bytes)
type PROTO_NC_RAID_SET_SUB_LEADER_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	SNewSubLeaderCharID              Name5                // offset: 20, size: 20 bytes
	NPartyNo                         uint16               // offset: 40, size: 2 bytes
	NPartyIndex                      uint8                // offset: 42, size: 1 bytes
	NSlotIndex                       uint8                // offset: 43, size: 1 bytes
}

// PROTO_NC_RAID_SET_SUB_LEADER_CMD (size = 20 bytes)
type PROTO_NC_RAID_SET_SUB_LEADER_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SET_SUB_LEADER_REQ (size = 20 bytes)
type PROTO_NC_RAID_SET_SUB_LEADER_REQ_DATA struct {
	SNewSubLeaderCharID              Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_SLOT_MOVE_ACK (size = 22 bytes)
type PROTO_NC_RAID_SLOT_MOVE_ACK_DATA struct {
	SSlotMoveCharID                  Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_SLOT_MOVE_BROAD_CMD (size = 228 bytes)
type PROTO_NC_RAID_SLOT_MOVE_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	PrePosition                      uint32               // offset: 20, size: 4 bytes
	PostPosition                     uint32               // offset: 24, size: 4 bytes
	SPreCharID                       [100]byte            // offset: 28, size: 100 bytes
	SPostCharID                      [100]byte            // offset: 128, size: 100 bytes
}

// PROTO_NC_RAID_SLOT_MOVE_REQ (size = 48 bytes)
type PROTO_NC_RAID_SLOT_MOVE_REQ_DATA struct {
	PrePosition                      uint32               // offset: 0, size: 4 bytes
	PostPosition                     uint32               // offset: 4, size: 4 bytes
	SPreSlotMoveCharID               Name5                // offset: 8, size: 20 bytes
	SPostSlotMoveCharID              Name5                // offset: 28, size: 20 bytes
}

// PROTO_NC_RAID_VANISH_ACK (size = 22 bytes)
type PROTO_NC_RAID_VANISH_ACK_DATA struct {
	SVanishCharID                    Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_RAID_VANISH_BROAD_CMD (size = 44 bytes)
type PROTO_NC_RAID_VANISH_BROAD_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	SCommandCharID                   Name5                // offset: 4, size: 20 bytes
	SVanishCharID                    Name5                // offset: 24, size: 20 bytes
}

// PROTO_NC_RAID_VANISH_CMD (size = 20 bytes)
type PROTO_NC_RAID_VANISH_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_VANISH_REQ (size = 20 bytes)
type PROTO_NC_RAID_VANISH_REQ_DATA struct {
	SVanishCharID                    Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_RAID_WARNING_ACK (size = 2 bytes)
type PROTO_NC_RAID_WARNING_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_WARNING_BROAD_CMD (size = 22 bytes)
type PROTO_NC_RAID_WARNING_BROAD_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	STalkerCharID                    Name5                // offset: 1, size: 20 bytes
	NChatLen                         uint8                // offset: 21, size: 1 bytes
	SChat                            [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_RAID_WARNING_REQ (size = 2 bytes)
type PROTO_NC_RAID_WARNING_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	NChatLen                         uint8                // offset: 1, size: 1 bytes
	SChat                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_RAID_ZONE_DEL_SUB_LEADER_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_DEL_SUB_LEADER_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NDelLeaderCharNo                 uint32               // offset: 2, size: 4 bytes
	NPartyNo                         uint16               // offset: 6, size: 2 bytes
	NPartyIndex                      uint8                // offset: 8, size: 1 bytes
	NSlotIndex                       uint8                // offset: 9, size: 1 bytes
}

// PROTO_NC_RAID_ZONE_DISMISS_CMD (size = 2 bytes)
type PROTO_NC_RAID_ZONE_DISMISS_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_RAID_ZONE_LEAVE_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_LEAVE_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NPartyNo                         uint16               // offset: 2, size: 2 bytes
	NPartyIndex                      uint8                // offset: 4, size: 1 bytes
	NSlotIndex                       uint8                // offset: 5, size: 1 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_RAID_ZONE_LOOTING_TYPE_CMD (size = 4 bytes)
type PROTO_NC_RAID_ZONE_LOOTING_TYPE_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NItemGrade                       uint8                // offset: 2, size: 1 bytes
	NLootingType                     uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_RAID_ZONE_MAKE_CMD (size = 32 bytes)
type PROTO_NC_RAID_ZONE_MAKE_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NRaidNo                          uint16               // offset: 2, size: 2 bytes
	NRaidPartyNo                     [4]uint16            // offset: 4, size: 8 bytes
	FirstPartyCharRegNo              [5]uint32            // offset: 12, size: 20 bytes
}

// PROTO_NC_RAID_ZONE_MEMBER_JOIN_CMD (size = 42 bytes)
type PROTO_NC_RAID_ZONE_MEMBER_JOIN_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NPartyNo                         uint16               // offset: 2, size: 2 bytes
	NPartyIndex                      uint8                // offset: 4, size: 1 bytes
	NSlotIndex                       uint8                // offset: 5, size: 1 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	TInfo                            [32]byte             // offset: 10, size: 32 bytes
}

// PROTO_NC_RAID_ZONE_SET_LEADER_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_SET_LEADER_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NNewLeaderCharNo                 uint32               // offset: 2, size: 4 bytes
	NPartyNo                         uint16               // offset: 6, size: 2 bytes
	NPartyIndex                      uint8                // offset: 8, size: 1 bytes
	NSlotIndex                       uint8                // offset: 9, size: 1 bytes
}

// PROTO_NC_RAID_ZONE_SET_LOOTER_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_SET_LOOTER_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NNewLooterCharNo                 uint32               // offset: 2, size: 4 bytes
	NPartyNo                         uint16               // offset: 6, size: 2 bytes
	NPartyIndex                      uint8                // offset: 8, size: 1 bytes
	NSlotIndex                       uint8                // offset: 9, size: 1 bytes
}

// PROTO_NC_RAID_ZONE_SET_SUB_LEADER_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_SET_SUB_LEADER_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NNewSubLeaderCharNo              uint32               // offset: 2, size: 4 bytes
	NPartyNo                         uint16               // offset: 6, size: 2 bytes
	NPartyIndex                      uint8                // offset: 8, size: 1 bytes
	NSlotIndex                       uint8                // offset: 9, size: 1 bytes
}

// PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD (size = 50 bytes)
type PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	PrePosition                      [24]byte             // offset: 2, size: 24 bytes
	PostPosition                     [24]byte             // offset: 26, size: 24 bytes
}

// PROTO_NC_RAID_ZONE_VANISH_CMD (size = 10 bytes)
type PROTO_NC_RAID_ZONE_VANISH_CMD_DATA struct {
	NRaidNo                          uint16               // offset: 0, size: 2 bytes
	NPartyNo                         uint16               // offset: 2, size: 2 bytes
	NPartyIndex                      uint8                // offset: 4, size: 1 bytes
	NSlotIndex                       uint8                // offset: 5, size: 1 bytes
	NVanishCharNo                    uint32               // offset: 6, size: 4 bytes
}

// ============================================================================
// Packet Builders Oficiales para el Sistema de Raid (Dept 44)
// ============================================================================

func New_PROTO_NC_RAID_MAKE_ACK(raidNo, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_MAKE_ACK)
	pw.WriteUint16(raidNo)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_MAKE_BROAD_CMD(raidNo uint16, partyNos [4]uint16, firstPartyChars [5]string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_MAKE_BROAD_CMD)
	pw.WriteUint16(raidNo)
	for i := 0; i < 4; i++ {
		pw.WriteUint16(partyNos[i])
	}
	for i := 0; i < 5; i++ {
		pw.WriteString(firstPartyChars[i], 20)
	}
	return pw
}

func New_PROTO_NC_RAID_LOOTING_TYPE_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_LOOTING_TYPE_ACK)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_LOOTING_TYPE_BROAD_CMD(itemGrade, lootType byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_LOOTING_TYPE_BROAD_CMD)
	pw.WriteByte(itemGrade)
	pw.WriteByte(lootType)
	return pw
}

func New_PROTO_NC_RAID_INVITATION_ACK(invitationCharID string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_INVITATION_ACK)
	pw.WriteString(invitationCharID, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_INVITATION_CMD(masterCharID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_INVITATION_CMD)
	pw.WriteString(masterCharID, 20)
	return pw
}

func New_PROTO_NC_RAID_INVITATION_ANSWER_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_INVITATION_ANSWER_ACK)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_MEMBER_JOIN_BROAD_CMD(partyNo uint16, partyIndex, slotIndex byte, charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_MEMBER_JOIN_BROAD_CMD)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	pw.WriteString(charID, 20)
	return pw
}

func New_PROTO_NC_RAID_LEAVE_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_LEAVE_ACK)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_LEAVE_BROAD_CMD(partyNo uint16, partyIndex, slotIndex byte, charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_LEAVE_BROAD_CMD)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	pw.WriteString(charID, 20)
	return pw
}

func New_PROTO_NC_RAID_VANISH_ACK(vanishCharID string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_VANISH_ACK)
	pw.WriteString(vanishCharID, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_VANISH_BROAD_CMD(partyNo uint16, partyIndex, slotIndex byte, cmdCharID, vanishCharID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_VANISH_BROAD_CMD)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	pw.WriteString(cmdCharID, 20)
	pw.WriteString(vanishCharID, 20)
	return pw
}

func New_PROTO_NC_RAID_SLOT_MOVE_ACK(charID string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SLOT_MOVE_ACK)
	pw.WriteString(charID, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_SLOT_MOVE_BROAD_CMD(cmdCharID string, prePos, postPos uint32, preChar, postChar string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SLOT_MOVE_BROAD_CMD)
	pw.WriteString(cmdCharID, 20)
	pw.WriteUint32(prePos)
	pw.WriteUint32(postPos)
	pw.WriteString(preChar, 100)
	pw.WriteString(postChar, 100)
	return pw
}

func New_PROTO_NC_RAID_SET_LEADER_ACK(newLeader string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_LEADER_ACK)
	pw.WriteString(newLeader, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_SET_LEADER_BROAD_CMD(cmdChar, newLeader string, partyNo uint16, partyIndex, slotIndex byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_LEADER_BROAD_CMD)
	pw.WriteString(cmdChar, 20)
	pw.WriteString(newLeader, 20)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	return pw
}

func New_PROTO_NC_RAID_SET_SUB_LEADER_ACK(newSubLeader string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_SUB_LEADER_ACK)
	pw.WriteString(newSubLeader, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_SET_SUB_LEADER_BROAD_CMD(cmdChar, newSubLeader string, partyNo uint16, partyIndex, slotIndex byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_SUB_LEADER_BROAD_CMD)
	pw.WriteString(cmdChar, 20)
	pw.WriteString(newSubLeader, 20)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	return pw
}

func New_PROTO_NC_RAID_DEL_SUB_LEADER_ACK(delSubLeader string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_DEL_SUB_LEADER_ACK)
	pw.WriteString(delSubLeader, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_DEL_SUB_LEADER_BROAD_CMD(cmdChar, delSubLeader string, partyNo uint16, partyIndex, slotIndex byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_DEL_SUB_LEADER_BROAD_CMD)
	pw.WriteString(cmdChar, 20)
	pw.WriteString(delSubLeader, 20)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	return pw
}

func New_PROTO_NC_RAID_SET_LOOTER_ACK(newLooter string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_LOOTER_ACK)
	pw.WriteString(newLooter, 20)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_SET_LOOTER_BROAD_CMD(cmdChar, newLooter string, partyNo uint16, partyIndex, slotIndex byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_SET_LOOTER_BROAD_CMD)
	pw.WriteString(cmdChar, 20)
	pw.WriteString(newLooter, 20)
	pw.WriteUint16(partyNo)
	pw.WriteByte(partyIndex)
	pw.WriteByte(slotIndex)
	return pw
}

func New_PROTO_NC_RAID_CHAT_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_CHAT_ACK)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_CHAT_BROAD_CMD(talker string, chat string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_CHAT_BROAD_CMD)
	pw.WriteByte(0) // ItemLinkDataCount
	pw.WriteString(talker, 20)
	pw.WriteByte(byte(len(chat)))
	pw.WriteBytes([]byte(chat))
	return pw
}

func New_PROTO_NC_RAID_WARNING_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_WARNING_ACK)
	pw.WriteUint16(errCode)
	return pw
}

func New_PROTO_NC_RAID_WARNING_BROAD_CMD(talker string, warning string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_RAID_WARNING_BROAD_CMD)
	pw.WriteByte(0) // ItemLinkDataCount
	pw.WriteString(talker, 20)
	pw.WriteByte(byte(len(warning)))
	pw.WriteBytes([]byte(warning))
	return pw
}
