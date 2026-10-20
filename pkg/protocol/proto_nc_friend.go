package protocol

// Packet builders oficiales para el sistema de amigos (Friend)
// Nomenclatura y offsets 100% idénticos a WorldManager.pdb / WorldManager.md (Dept: 21)

// PROTO_NC_FRIEND_SET_CONFIRM_REQ envía la solicitud de amistad al jugador objetivo (OpCode 0x5403).
// Struct oficial PDB: PROTO_NC_FRIEND_SET_CONFIRM_REQ (20 bytes: Name5 requestername)
func New_PROTO_NC_FRIEND_SET_CONFIRM_REQ(requesterName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_SET_CONFIRM_REQ)
	pw.WriteString(requesterName, 20)
	return pw
}

// New_PROTO_NC_FRIEND_SET_ACK responde el resultado del intento de agregar amigo (OpCode 0x5402).
// Struct oficial PDB: PROTO_NC_FRIEND_SET_ACK (22 bytes: Name5 targetname, unsigned short err)
func New_PROTO_NC_FRIEND_SET_ACK(targetName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_SET_ACK)
	pw.WriteString(targetName, 20)
	pw.WriteUint16(errCode)
	return pw
}

// FriendEntry representa la información de cada amigo en PROTO_NC_FRIEND_LIST_CMD.
// Struct oficial PDB: PROTO_FRIEND_INFO (charid[20], bIsLogin[1], map[12], charclass[1], level[1])
type FriendEntry struct {
	CharName string // Name5 (20 bytes)
	IsOnline bool   // 1 byte
	Map      string // Name3 (12 bytes)
	Class    byte   // 1 byte
	Level    byte   // 1 byte
}

// New_PROTO_NC_FRIEND_LIST_CMD construye la lista oficial de amigos del jugador (OpCode 0x5407).
func New_PROTO_NC_FRIEND_LIST_CMD(friends []FriendEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_LIST_CMD)
	pw.WriteByte(byte(len(friends)))
	for _, f := range friends {
		pw.WriteString(f.CharName, 20)
		pw.WriteBool(f.IsOnline)
		pw.WriteString(f.Map, 12)
		pw.WriteByte(f.Class)
		pw.WriteByte(f.Level)
	}
	return pw
}

// New_PROTO_NC_FRIEND_ADD_CMD notifica la incorporación de un nuevo amigo a la lista (OpCode 0x5408).
func New_PROTO_NC_FRIEND_ADD_CMD(f FriendEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_ADD_CMD)
	pw.WriteString(f.CharName, 20)
	pw.WriteBool(f.IsOnline)
	pw.WriteString(f.Map, 12)
	pw.WriteByte(f.Class)
	pw.WriteByte(f.Level)
	return pw
}

// New_PROTO_NC_FRIEND_DEL_ACK confirma la eliminación de un amigo (OpCode 0x5406).
func New_PROTO_NC_FRIEND_DEL_ACK(friendName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_DEL_ACK)
	pw.WriteString(friendName, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_FRIEND_LOGIN_CMD notifica a los amigos que un jugador se conectó (OpCode 0x5409).
// Struct oficial PDB: PROTO_NC_FRIEND_LOGIN_CMD (32 bytes: Name5 charid, Name3 map)
func New_PROTO_NC_FRIEND_LOGIN_CMD(charName, mapName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_LOGIN_CMD)
	pw.WriteString(charName, 20)
	pw.WriteString(mapName, 12)
	return pw
}

// New_PROTO_NC_FRIEND_LOGOUT_CMD notifica a los amigos que un jugador se desconectó (OpCode 0x540A).
// Struct oficial PDB: PROTO_NC_FRIEND_LOGOUT_CMD (20 bytes: Name5 charid)
func New_PROTO_NC_FRIEND_LOGOUT_CMD(charName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_LOGOUT_CMD)
	pw.WriteString(charName, 20)
	return pw
}

// New_PROTO_NC_FRIEND_POINT_ACK responde la consulta de puntos de amistad (OpCode 0x541E).
func New_PROTO_NC_FRIEND_POINT_ACK(points uint32) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_FRIEND_POINT_ACK)
	pw.WriteUint32(points)
	return pw
}
