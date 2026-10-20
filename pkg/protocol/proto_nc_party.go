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
