package protocol

// New_PROTO_NC_CHAR_CREATE_FAIL_ACK construye el error oficial de creación de personaje (OpCode 0x1404).
// Struct oficial PDB: PROTO_NC_AVATAR_CREATEFAIL_ACK (2 bytes)
// - err: unsigned short (uint16)
func New_PROTO_NC_CHAR_CREATE_FAIL_ACK(errorCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_CREATE_FAIL_ACK)
	pw.WriteUint16(errorCode)
	return pw
}

// New_PROTO_NC_CHAR_CREATE_ACK construye la confirmación oficial de creación de personaje (OpCode 0x1406).
// Struct oficial PDB: PROTO_NC_AVATAR_CREATESUCC_ACK (131 bytes)
// - numofavatar: unsigned char (1 byte)
// - avatar: PROTO_AVATARINFORMATION (130 bytes)
func New_PROTO_NC_CHAR_CREATE_ACK(charInfoWriter func(pw *PacketWriter)) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_CREATE_ACK)
	pw.WriteByte(1)
	if charInfoWriter != nil {
		charInfoWriter(pw)
	}
	return pw
}

// New_PROTO_NC_CHAR_DELETE_ACK construye la confirmación oficial de borrado de personaje (OpCode 0x140C).
// Struct oficial PDB: PROTO_NC_AVATAR_ERASESUC_ACK (1 byte)
// - slot: unsigned char
func New_PROTO_NC_CHAR_DELETE_ACK(slot byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_DELETE_ACK)
	pw.WriteByte(slot)
	return pw
}

// New_PROTO_NC_CHAR_LOGIN_ACK construye la respuesta oficial de traspaso a ZoneServer (OpCode 0x1003).
// Struct oficial PDB: PROTO_NC_CHAR_LOGIN_ACK (18 bytes)
// - zoneip: char[16] (Name4)
// - zoneport: unsigned short (uint16)
func New_PROTO_NC_CHAR_LOGIN_ACK(ip string, port uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_LOGIN_ACK)
	pw.WriteString(ip, 16)
	pw.WriteUint16(port)
	return pw
}

// New_PROTO_NC_CHAR_LOGINFAIL_ACK construye la respuesta oficial de fallo al seleccionar personaje (OpCode 0x1002).
// Struct oficial PDB: PROTO_NC_CHAR_LOGINFAIL_ACK (2 bytes)
// - err: unsigned short (uint16)
func New_PROTO_NC_CHAR_LOGINFAIL_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_LOGINFAIL_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_USER_BACKTOWORLDLIST_ACK construye la respuesta para volver a la selección de mundos (OpCode 0x0C34).
func New_PROTO_NC_USER_BACKTOWORLDLIST_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_BACKTOWORLDLIST_ACK)
	pw.WriteBytes([]byte{0x58, 0x1E})
	pw.Fill(32, 0)
	return pw
}

// New_PROTO_NC_MISC_GAMETIME_ACK construye la respuesta horaria del servidor de juego (OpCode 0x080E).
// Struct oficial PDB: PROTO_NC_MISC_GAMETIME_ACK (3 bytes)
// - hour: unsigned char
// - minute: unsigned char
// - second: unsigned char
func New_PROTO_NC_MISC_GAMETIME_ACK(hour, minute, second byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_MISC_GAMETIME_ACK)
	pw.WriteByte(hour)
	pw.WriteByte(minute)
	pw.WriteByte(second)
	return pw
}
