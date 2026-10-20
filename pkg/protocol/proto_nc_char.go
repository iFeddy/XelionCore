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

// New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK construye la respuesta oficial para volver a la selección de mundos (OpCode 0x0C34).
// Struct oficial PDB: PROTO_NC_USER_WILL_WORLD_SELECT_ACK (34 bytes: unsigned short nError, Name8 sOTP)
func New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_WILL_WORLD_SELECT_ACK)
	pw.WriteUint16(0) // nError = 0 (éxito)
	pw.Fill(32, 0)    // sOTP (vacío)
	return pw
}

// New_PROTO_NC_USER_BACKTOWORLDLIST_ACK es un alias retrocompatible para New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK (OpCode 0x0C34).
func New_PROTO_NC_USER_BACKTOWORLDLIST_ACK() *PacketWriter {
	return New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK()
}

// New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK construye la respuesta de accesos directos (OpCode 0x7003).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK (1025 bytes: 1 byte bSuccess + 1024 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 1024 {
			pw.WriteBytes(data[:1024])
		} else {
			pw.WriteBytes(data)
			pw.Fill(1024-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(1024, 0)
	}
	return pw
}

// New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK construye la respuesta del tamaño de quickbar (OpCode 0x7005).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK (25 bytes: 1 byte bSuccess + 24 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 24 {
			pw.WriteBytes(data[:24])
		} else {
			pw.WriteBytes(data)
			pw.Fill(24-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(24, 0)
	}
	return pw
}

// New_PROTO_NC_CHAR_OPTION_GET_GAME_ACK construye la respuesta de opciones de juego (OpCode 0x7009 / 0x700B).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_GAME_ACK (65 bytes: 1 byte bSuccess + 64 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_GAME_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_GAME_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 64 {
			pw.WriteBytes(data[:64])
		} else {
			pw.WriteBytes(data)
			pw.Fill(64-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(64, 0)
	}
	return pw
}

// New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK construye la respuesta de posición de ventanas (OpCode 0x700B / 0x700D).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK (137 bytes: 1 byte bSuccess + 136 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 136 {
			pw.WriteBytes(data[:136])
		} else {
			pw.WriteBytes(data)
			pw.Fill(136-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(136, 0)
	}
	return pw
}

// New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK construye la respuesta de mapeo de teclas (OpCode 0x700D / 0x700F).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK (53 bytes: 1 byte bSuccess + 52 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 52 {
			pw.WriteBytes(data[:52])
		} else {
			pw.WriteBytes(data)
			pw.Fill(52-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(52, 0)
	}
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
