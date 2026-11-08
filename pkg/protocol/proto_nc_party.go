package protocol

// Packet builders oficiales para el sistema de grupos (Party)
// Nomenclatura y offsets 100% idénticos a WorldManager.pdb / WorldManager.md (Dept: 14)

// PROTO_NC_PARTY_JOINPROPOSE_REQ construye la invitación a grupo recibida por el invitado (OpCode 0x3803).
// Struct oficial PDB: PROTO_NC_PARTY_JOINPROPOSE_REQ (20 bytes: Name5 mastername)
func New_PROTO_NC_PARTY_JOINPROPOSE_REQ(masterName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_JOINPROPOSE_REQ)
	pw.WriteString(masterName, 20)
	return pw
}

// New_PROTO_NC_PARTY_JOIN_ACK construye la respuesta de invitación a grupo (OpCode 0x3807).
// Struct oficial PDB: PROTO_NC_PARTY_JOIN_ACK (23 bytes: Name5 memberid, unsigned char charclass, unsigned short err)
func New_PROTO_NC_PARTY_JOIN_ACK(memberID string, charClass byte, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_JOIN_ACK)
	pw.WriteString(memberID, 20)
	pw.WriteByte(charClass)
	pw.WriteUint16(errCode)
	return pw
}

// PartyMemberInfo representa cada miembro dentro de PROTO_NC_PARTY_MEMBER_LIST_CMD.
// Struct oficial PDB: PROTO_NC_PARTY_MEMBER_LIST_CMD::MEMBER_INFO (22 bytes)
type PartyMemberInfo struct {
	MemberName   string // Name5 (20 bytes)
	MemberStatus byte   // 1 = Master, 0 = Member
	CharClass    byte   // Class ID
}

// New_PROTO_NC_PARTY_MEMBER_LIST_CMD construye la lista oficial de miembros del grupo (OpCode 0x3809).
// Struct oficial PDB: PROTO_NC_PARTY_MEMBER_LIST_CMD (1 byte membercount + membercount * 22 bytes)
func New_PROTO_NC_PARTY_MEMBER_LIST_CMD(members []PartyMemberInfo) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_MEMBER_LIST_CMD)
	pw.WriteByte(byte(len(members)))
	for _, m := range members {
		pw.WriteString(m.MemberName, 20)
		pw.WriteByte(m.MemberStatus)
		pw.WriteByte(m.CharClass)
	}
	return pw
}

// New_PROTO_NC_PARTY_LEAVE_ACK construye la confirmación de salida de grupo (OpCode 0x380B).
func New_PROTO_NC_PARTY_LEAVE_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_LEAVE_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_PARTY_LEAVE_CMD notifica la salida de un miembro a todo el grupo (OpCode 0x380C).
// Struct oficial PDB: PROTO_NC_PARTY_LEAVE_CMD (6 bytes: unsigned short partyregnum, unsigned int leavememberregnum)
func New_PROTO_NC_PARTY_LEAVE_CMD(partyID uint16, charID uint32) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_LEAVE_CMD)
	pw.WriteUint16(partyID)
	pw.WriteUint32(charID)
	return pw
}

// New_PROTO_NC_PARTY_KICKOFF_ACK construye la confirmación de expulsión de un miembro (OpCode 0x3815).
// Struct oficial PDB: PROTO_NC_PARTY_KICKOFF_ACK (22 bytes: Name5 memberid, unsigned short err)
func New_PROTO_NC_PARTY_KICKOFF_ACK(targetName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_KICKOFF_ACK)
	pw.WriteString(targetName, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_PARTY_DISMISS_CMD notifica la disolución del grupo a sus miembros (OpCode 0x381F).
// Struct oficial PDB: PROTO_NC_PARTY_DISMISS_CMD (2 bytes: unsigned short partyregnum)
func New_PROTO_NC_PARTY_DISMISS_CMD(partyID uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_DISMISS_CMD)
	pw.WriteUint16(partyID)
	return pw
}

// New_PROTO_NC_PARTY_CHANGEMASTER_ACK construye la confirmación de cambio de líder (OpCode 0x3829).
func New_PROTO_NC_PARTY_CHANGEMASTER_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_CHANGEMASTER_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_PARTY_ITEM_LOOTING_CMD notifica la regla de botín a los miembros (OpCode 0x384C).
// Struct oficial PDB: PROTO_NC_PARTY_ITEM_LOOTING_CMD (2 bytes: unsigned char lootingtype, unsigned char nItemGrade)
func New_PROTO_NC_PARTY_ITEM_LOOTING_CMD(lootType, itemGrade byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_ITEM_LOOTING_CMD)
	pw.WriteByte(lootType)
	pw.WriteByte(itemGrade)
	return pw
}

// New_PROTO_NC_PARTY_SET_LOOTER_ACK construye la respuesta de asignación de recolector (OpCode 0x385B).
func New_PROTO_NC_PARTY_SET_LOOTER_ACK(looterName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_SET_LOOTER_ACK)
	pw.WriteString(looterName, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_PARTY_SET_LOOTER_BROAD_CMD transmite a toda la party el nuevo recolector (OpCode 0x385D).
func New_PROTO_NC_PARTY_SET_LOOTER_BROAD_CMD(commanderName, looterName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_SET_LOOTER_BROAD_CMD)
	pw.WriteString(commanderName, 20)
	pw.WriteString(looterName, 20)
	return pw
}

// PartyFinderEntry representa una entrada en el buscador de grupos (0x3854 / 0x3855).
type PartyFinderEntry struct {
	MasterName string
	Level      byte
	MemberCnt  byte
	Purpose    string
}

// New_PROTO_NC_PARTY_FINDER_LIST_ACK construye la lista del buscador de grupos (OpCode 0x3855).
func New_PROTO_NC_PARTY_FINDER_LIST_ACK(entries []PartyFinderEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_PARTY_FINDER_LIST_ACK)
	pw.WriteByte(byte(len(entries)))
	for _, e := range entries {
		pw.WriteString(e.MasterName, 20)
		pw.WriteByte(e.Level)
		pw.WriteByte(e.MemberCnt)
		pw.WriteString(e.Purpose, 64)
	}
	return pw
}

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_ACT_PARTYCHAT_ACK (size = 2 bytes)
type PROTO_NC_ACT_PARTYCHAT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ACT_PARTYCHAT_ACK_SEND (size = 5 bytes)
type PROTO_NC_ACT_PARTYCHAT_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_ACT_PARTYCHAT_CMD (size = 22 bytes)
type PROTO_NC_ACT_PARTYCHAT_CMD_DATA struct {
	Talker                           Name5                // offset: 0, size: 20 bytes
	Chat                             uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_ACT_PARTYCHAT_REQ (size = 2 bytes)
type PROTO_NC_ACT_PARTYCHAT_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_FRIEND_PARTY_CMD (size = 21 bytes)
type PROTO_NC_FRIEND_PARTY_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Is_party                         uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_PARTY_CHANGEMASTER_ACK (size = 22 bytes)
type PROTO_NC_PARTY_CHANGEMASTER_ACK_DATA struct {
	Newmaster                        Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PARTY_CHANGEMASTER_ACK_SEND (size = 25 bytes)
type PROTO_NC_PARTY_CHANGEMASTER_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [22]byte             // offset: 3, size: 22 bytes
}

// PROTO_NC_PARTY_CHANGEMASTER_CMD (size = 6 bytes)
type PROTO_NC_PARTY_CHANGEMASTER_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Newmaster_chrregnum              uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PARTY_CHANGEMASTER_CMD_SEND (size = 9 bytes)
type PROTO_NC_PARTY_CHANGEMASTER_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [6]byte              // offset: 3, size: 6 bytes
}

// PROTO_NC_PARTY_CHANGEMASTER_REQ (size = 20 bytes)
type PROTO_NC_PARTY_CHANGEMASTER_REQ_DATA struct {
	Newmaster                        Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_DISMISS_ACK_SEND (size = 3 bytes)
type PROTO_NC_PARTY_DISMISS_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_PARTY_DISMISS_CMD (size = 2 bytes)
type PROTO_NC_PARTY_DISMISS_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PARTY_DISMISS_CMD_SEND (size = 5 bytes)
type PROTO_NC_PARTY_DISMISS_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_PARTY_FINDER_ADD_ACK (size = 2 bytes)
type PROTO_NC_PARTY_FINDER_ADD_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PARTY_FINDER_ADD_REQ (size = 128 bytes)
type PROTO_NC_PARTY_FINDER_ADD_REQ_DATA struct {
	SMsg                             [128]byte            // offset: 0, size: 128 bytes
}

// PROTO_NC_PARTY_FINDER_DELETE_ACK (size = 2 bytes)
type PROTO_NC_PARTY_FINDER_DELETE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PARTY_FINDER_DELETE_REQ (size = 1 bytes)
type PROTO_NC_PARTY_FINDER_DELETE_REQ_DATA struct {}

// PROTO_NC_PARTY_FINDER_DELETE_YOUR_MSG_CMD (size = 1 bytes)
type PROTO_NC_PARTY_FINDER_DELETE_YOUR_MSG_CMD_DATA struct {}

// PROTO_NC_PARTY_FINDER_INFO (size = 154 bytes)
type PROTO_NC_PARTY_FINDER_INFO struct {
	NPastSecond                      uint16               // offset: 0, size: 2 bytes
	DRegistTime                      int32                // offset: 2, size: 4 bytes
	SCharID                          Name5                // offset: 6, size: 20 bytes
	SMsg                             [128]byte            // offset: 26, size: 128 bytes
}

// PROTO_NC_PARTY_FINDER_LIST_ACK (size = 8 bytes)
type PROTO_NC_PARTY_FINDER_LIST_ACK_DATA struct {
	NPageNo                          uint16               // offset: 0, size: 2 bytes
	NCountOfOnePage                  uint8                // offset: 2, size: 1 bytes
	BYouAreRegisted                  uint8                // offset: 3, size: 1 bytes
	NMaxPages                        uint16               // offset: 4, size: 2 bytes
	NNumOfList                       uint16               // offset: 6, size: 2 bytes
	InfoList                         [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_PARTY_FINDER_LIST_REQ (size = 3 bytes)
type PROTO_NC_PARTY_FINDER_LIST_REQ_DATA struct {
	NPageNo                          uint16               // offset: 0, size: 2 bytes
	NCountOfOnePage                  uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_PARTY_FUNDAMENTAL_CMD (size = 66 bytes)
type PROTO_NC_PARTY_FUNDAMENTAL_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Master                           [32]byte             // offset: 2, size: 32 bytes
	Member                           [32]byte             // offset: 34, size: 32 bytes
}

// PROTO_NC_PARTY_FUNDAMENTAL_CMD_SEND (size = 69 bytes)
type PROTO_NC_PARTY_FUNDAMENTAL_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [66]byte             // offset: 3, size: 66 bytes
}

// PROTO_NC_PARTY_ITEM_JOIN_LOOTING_CMD (size = 22 bytes)
type PROTO_NC_PARTY_ITEM_JOIN_LOOTING_CMD_DATA struct {
	NItemGrade                       uint8                // offset: 0, size: 1 bytes
	NLootingType                     uint8                // offset: 1, size: 1 bytes
	LooterName                       Name5                // offset: 2, size: 20 bytes
}

// PROTO_NC_PARTY_ITEM_LOOTING_CMD (size = 2 bytes)
type PROTO_NC_PARTY_ITEM_LOOTING_CMD_DATA struct {
	LootType                         uint8                // offset: 0, size: 1 bytes
	ItemGrade                        uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_PARTY_ITEM_LOOTING_SET (size = 2 bytes)
type PROTO_NC_PARTY_ITEM_LOOTING_SET_DATA struct {
	Lootingtype                      uint8                // offset: 0, size: 1 bytes
	NItemGrade                       uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_PARTY_ITEM_LOOTING_ZONE_CMD (size = 8 bytes)
type PROTO_NC_PARTY_ITEM_LOOTING_ZONE_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	NItemGrade                       uint8                // offset: 2, size: 1 bytes
	NLootingType                     uint8                // offset: 3, size: 1 bytes
	LooterNo                         uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_PARTY_JOINPROPOSE_REQ (size = 20 bytes)
type PROTO_NC_PARTY_JOINPROPOSE_REQ_DATA struct {
	Mastername                       Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_JOINPROPOSE_REQ_SEND (size = 23 bytes)
type PROTO_NC_PARTY_JOINPROPOSE_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [20]byte             // offset: 3, size: 20 bytes
}

// PROTO_NC_PARTY_JOIN_ACK (size = 23 bytes)
type PROTO_NC_PARTY_JOIN_ACK_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
	Charclass                        uint8                // offset: 20, size: 1 bytes
	Err                              uint16               // offset: 21, size: 2 bytes
}

// PROTO_NC_PARTY_JOIN_ACK_SEND (size = 26 bytes)
type PROTO_NC_PARTY_JOIN_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [23]byte             // offset: 3, size: 23 bytes
}

// PROTO_NC_PARTY_JOIN_CMD (size = 34 bytes)
type PROTO_NC_PARTY_JOIN_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Joiner                           [32]byte             // offset: 2, size: 32 bytes
}

// PROTO_NC_PARTY_JOIN_CMD_SEND (size = 37 bytes)
type PROTO_NC_PARTY_JOIN_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [34]byte             // offset: 3, size: 34 bytes
}

// PROTO_NC_PARTY_JOIN_REQ (size = 20 bytes)
type PROTO_NC_PARTY_JOIN_REQ_DATA struct {
	Target                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_KICKOFF_ACK (size = 22 bytes)
type PROTO_NC_PARTY_KICKOFF_ACK_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PARTY_KICKOFF_ACK_SEND (size = 25 bytes)
type PROTO_NC_PARTY_KICKOFF_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [22]byte             // offset: 3, size: 22 bytes
}

// PROTO_NC_PARTY_KICKOFF_CMD (size = 6 bytes)
type PROTO_NC_PARTY_KICKOFF_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Kickoffmemberregnum              uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PARTY_KICKOFF_CMD_SEND (size = 9 bytes)
type PROTO_NC_PARTY_KICKOFF_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [6]byte              // offset: 3, size: 6 bytes
}

// PROTO_NC_PARTY_KICKOFF_REQ (size = 20 bytes)
type PROTO_NC_PARTY_KICKOFF_REQ_DATA struct {
	Member                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_LEAVE_ACK (size = 22 bytes)
type PROTO_NC_PARTY_LEAVE_ACK_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PARTY_LEAVE_ACK_SEND (size = 25 bytes)
type PROTO_NC_PARTY_LEAVE_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [22]byte             // offset: 3, size: 22 bytes
}

// PROTO_NC_PARTY_LEAVE_CMD (size = 33 bytes)
type PROTO_NC_PARTY_LEAVE_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Leavemember                      [31]byte             // offset: 2, size: 31 bytes
}

// PROTO_NC_PARTY_LEAVE_CMD_SEND (size = 36 bytes)
type PROTO_NC_PARTY_LEAVE_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [33]byte             // offset: 3, size: 33 bytes
}

// PROTO_NC_PARTY_LEAVE_REQ (size = 1 bytes)
type PROTO_NC_PARTY_LEAVE_REQ_DATA struct {}

// PROTO_NC_PARTY_LOGININFO_CMD (size = 20 bytes)
type PROTO_NC_PARTY_LOGININFO_CMD_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_LOGININFO_CMD_SEND (size = 23 bytes)
type PROTO_NC_PARTY_LOGININFO_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [20]byte             // offset: 3, size: 20 bytes
}

// PROTO_NC_PARTY_LOGIN_CMD (size = 6 bytes)
type PROTO_NC_PARTY_LOGIN_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PARTY_LOGIN_CMD_SEND (size = 9 bytes)
type PROTO_NC_PARTY_LOGIN_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [6]byte              // offset: 3, size: 6 bytes
}

// PROTO_NC_PARTY_LOGOUTINFO_CMD (size = 20 bytes)
type PROTO_NC_PARTY_LOGOUTINFO_CMD_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_LOGOUTINFO_CMD_SEND (size = 23 bytes)
type PROTO_NC_PARTY_LOGOUTINFO_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [20]byte             // offset: 3, size: 20 bytes
}

// PROTO_NC_PARTY_LOGOUT_CMD (size = 6 bytes)
type PROTO_NC_PARTY_LOGOUT_CMD_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PARTY_LOGOUT_CMD_SEND (size = 9 bytes)
type PROTO_NC_PARTY_LOGOUT_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [6]byte              // offset: 3, size: 6 bytes
}

// PROTO_NC_PARTY_MEMBERINFORM_ACK (size = 3 bytes)
type PROTO_NC_PARTY_MEMBERINFORM_ACK_DATA struct {
	PartyRegnum                      uint16               // offset: 0, size: 2 bytes
	NumberOfMember                   uint8                // offset: 2, size: 1 bytes
	Member                           [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_PARTY_MEMBERINFORM_REQ (size = 2 bytes)
type PROTO_NC_PARTY_MEMBERINFORM_REQ_DATA struct {
	Partyregnum                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PARTY_MEMBER_LIST_CMD (size = 1 bytes)
type PROTO_NC_PARTY_MEMBER_LIST_CMD_DATA struct {
	Membercount                      uint8                // offset: 0, size: 1 bytes
	Memberarray                      [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_PARTY_MEMBER_LIST_CMD_SEND (size = 114 bytes)
type PROTO_NC_PARTY_MEMBER_LIST_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Membercount                      uint8                // offset: 3, size: 1 bytes
	Memberarray                      [110]byte            // offset: 4, size: 110 bytes
}

// PROTO_NC_PARTY_SET_LOOTER_ACK (size = 22 bytes)
type PROTO_NC_PARTY_SET_LOOTER_ACK_DATA struct {
	SNewLooterCharID                 Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PARTY_SET_LOOTER_BROAD_CMD (size = 40 bytes)
type PROTO_NC_PARTY_SET_LOOTER_BROAD_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
	SNewLooterCharID                 Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_PARTY_SET_LOOTER_CMD (size = 20 bytes)
type PROTO_NC_PARTY_SET_LOOTER_CMD_DATA struct {
	SCommandCharID                   Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_SET_LOOTER_REQ (size = 20 bytes)
type PROTO_NC_PARTY_SET_LOOTER_REQ_DATA struct {
	SNewLooterCharID                 Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_PARTY_ZONE_JOIN_CMD (size = 8 bytes)
type PROTO_NC_PARTY_ZONE_JOIN_CMD_DATA struct {
	NMasterCharNo                    uint32               // offset: 0, size: 4 bytes
	NJoinCharNo                      uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_PARTY_ZONE_LEAVE_CMD (size = 4 bytes)
type PROTO_NC_PARTY_ZONE_LEAVE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_PARTY_ZONE_SET_LOOTER_CMD (size = 6 bytes)
type PROTO_NC_PARTY_ZONE_SET_LOOTER_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NNewLooterCharNo                 uint32               // offset: 2, size: 4 bytes
}
