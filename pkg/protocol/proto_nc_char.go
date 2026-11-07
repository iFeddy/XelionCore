package protocol

// New_PROTO_NC_AVATAR_CREATEFAIL_ACK construye el error oficial de creación de personaje (OpCode 0x1404).
// Struct oficial PDB: PROTO_NC_AVATAR_CREATEFAIL_ACK (2 bytes)
// - err: unsigned short (uint16)
func New_PROTO_NC_AVATAR_CREATEFAIL_ACK(errorCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_AVATAR_CREATEFAIL_ACK)
	pw.WriteUint16(errorCode)
	return pw
}

// New_PROTO_NC_AVATAR_CREATESUCC_ACK construye la confirmación oficial de creación de personaje (OpCode 0x1406).
// Struct oficial PDB: PROTO_NC_AVATAR_CREATESUCC_ACK (131 bytes)
// - numofavatar: unsigned char (1 byte)
// - avatar: PROTO_AVATARINFORMATION (130 bytes)
func New_PROTO_NC_AVATAR_CREATESUCC_ACK(charInfoWriter func(pw *PacketWriter)) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_AVATAR_CREATESUCC_ACK)
	pw.WriteByte(1)
	if charInfoWriter != nil {
		charInfoWriter(pw)
	}
	return pw
}

// New_PROTO_NC_AVATAR_ERASESUCC_ACK construye la confirmación oficial de borrado de personaje (OpCode 0x140C).
// Struct oficial PDB: PROTO_NC_AVATAR_ERASESUC_ACK (1 byte)
// - slot: unsigned char
func New_PROTO_NC_AVATAR_ERASESUCC_ACK(slot byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_AVATAR_ERASESUCC_ACK)
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
func New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK(otp [32]byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_WILL_WORLD_SELECT_ACK)
	pw.WriteUint16(0) // nError = 0 (éxito)
	pw.WriteBytes(otp[:])
	return pw
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

// New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK construye la respuesta de posición de ventanas (OpCode 0x700D).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK (393 bytes / 0x189: 1 byte bSuccess + 392 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 392 {
			pw.WriteBytes(data[:392])
		} else {
			pw.WriteBytes(data)
			pw.Fill(392-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(392, 0)
	}
	return pw
}

// New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK construye la respuesta de mapeo de teclas (OpCode 0x700F).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK (309 bytes / 0x135: 1 byte bSuccess + 308 bytes Data).
func New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK)
	if len(data) > 0 {
		pw.WriteByte(1)
		if len(data) >= 308 {
			pw.WriteBytes(data[:308])
		} else {
			pw.WriteBytes(data)
			pw.Fill(308-len(data), 0)
		}
	} else {
		pw.WriteByte(0)
		pw.Fill(308, 0)
	}
	return pw
}

// SHORT_CUT_DATA representa un bloque individual de acceso directo del sistema mejorado NA2016 (7 bytes).
// Struct oficial PDB: SHORT_CUT_DATA { unsigned char nSlotNo; unsigned short nCodeNo; int nValue; }
type SHORT_CUT_DATA struct {
	SlotNo byte
	CodeNo uint16
	Value  int32
}

// KEY_MAP_DATA representa un bloque individual de mapeo de tecla del sistema mejorado NA2016 (4 bytes).
// Struct oficial PDB: KEY_MAP_DATA { unsigned short nFunctionNo; unsigned char nExtendKey; unsigned char nASCIICode; }
type KEY_MAP_DATA struct {
	FunctionNo uint16
	ExtendKey  byte
	ASCIICode  byte
}

// GAME_OPTION_DATA representa un bloque individual de opción de juego del sistema mejorado NA2016 (3 bytes).
// Struct oficial PDB: GAME_OPTION_DATA { unsigned short nOptionNo; unsigned char nValue; }
type GAME_OPTION_DATA struct {
	OptionNo uint16
	Value    byte
}

// New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD construye el comando de accesos directos mejorados (OpCode 0x7032).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD (2 + 7*N bytes):
// - nShortCutDataCnt: unsigned short (2 bytes)
// - ShortCutData: SHORT_CUT_DATA[nShortCutDataCnt] (7 bytes cada uno)
func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD(entries []SHORT_CUT_DATA, rawData []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD)
	if len(entries) > 0 {
		pw.WriteUint16(uint16(len(entries)))
		for _, e := range entries {
			pw.WriteByte(e.SlotNo)
			pw.WriteUint16(e.CodeNo)
			pw.WriteInt32(e.Value)
		}
		return pw
	}
	if len(rawData) >= 2 {
		pw.WriteBytes(rawData)
		return pw
	}
	pw.WriteUint16(0) // nShortCutDataCnt = 0 (estructura vacía válida)
	return pw
}

// New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD construye el comando de mapeo de teclas mejorado (OpCode 0x7033).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD (2 + 4*N bytes):
// - nKeyMapDataCnt: unsigned short (2 bytes)
// - KeyMapData: KEY_MAP_DATA[nKeyMapDataCnt] (4 bytes cada uno)
func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD(entries []KEY_MAP_DATA, rawData []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD)
	if len(entries) > 0 {
		pw.WriteUint16(uint16(len(entries)))
		for _, e := range entries {
			pw.WriteUint16(e.FunctionNo)
			pw.WriteByte(e.ExtendKey)
			pw.WriteByte(e.ASCIICode)
		}
		return pw
	}
	if len(rawData) >= 2 {
		pw.WriteBytes(rawData)
		return pw
	}
	pw.WriteUint16(0) // nKeyMapDataCnt = 0 (estructura vacía válida)
	return pw
}

// New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD construye el comando de opciones de juego mejoradas (OpCode 0x7034).
// Struct oficial PDB: PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD (2 + 3*N bytes):
// - nGameOptionDataCnt: unsigned short (2 bytes)
// - GameOptionData: GAME_OPTION_DATA[nGameOptionDataCnt] (3 bytes cada uno)
func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD(entries []GAME_OPTION_DATA, rawData []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD)
	if len(entries) > 0 {
		pw.WriteUint16(uint16(len(entries)))
		for _, e := range entries {
			pw.WriteUint16(e.OptionNo)
			pw.WriteByte(e.Value)
		}
		return pw
	}
	if len(rawData) >= 2 {
		pw.WriteBytes(rawData)
		return pw
	}
	pw.WriteUint16(0) // nGameOptionDataCnt = 0 (estructura vacía válida)
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

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_ANNOUNCE_W2C_CMD (size = 2 bytes)
type PROTO_NC_ANNOUNCE_W2C_CMD_DATA struct {
	AnnounceMessage                  uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ANNOUNCE_Z2W_CMD (size = 2 bytes)
type PROTO_NC_ANNOUNCE_Z2W_CMD_DATA struct {
	AnnounceMessage                  uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_BAT_CLIENT_MOB_KILL_ANNOUNCE_CMD (size = 4 bytes)
type PROTO_NC_BAT_CLIENT_MOB_KILL_ANNOUNCE_CMD_DATA struct {
	NTextIndex                       uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_BAT_PKINPKFIELD_WMS_CMD (size = 10 bytes)
type PROTO_NC_BAT_PKINPKFIELD_WMS_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Pkcount                          uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD (size = 4 bytes)
type PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD_DATA struct {
	NTextIndex                       uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD (size = 3 bytes)
type PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Count                            uint8                // offset: 2, size: 1 bytes
	InfoList                         [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_CHARGED_BOOTHSLOTSIZE_CMD (size = 1 bytes)
type PROTO_NC_CHARGED_BOOTHSLOTSIZE_CMD_DATA struct {
	Boothsize                        uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHARGED_BUFFTERMINATE_CMD (size = 4 bytes)
type PROTO_NC_CHARGED_BUFFTERMINATE_CMD_DATA struct {
	ChargedBuffKey                   uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHARGED_DELETEWEAPONTITLE_CMD (size = 2 bytes)
type PROTO_NC_CHARGED_DELETEWEAPONTITLE_CMD_DATA struct {
	Licenseitem                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHARGED_RESETBUFF_CMD (size = 14 bytes)
type PROTO_NC_CHARGED_RESETBUFF_CMD_DATA struct {
	CharRegistNumber                 uint32               // offset: 0, size: 4 bytes
	ChargedBuffKey                   uint32               // offset: 4, size: 4 bytes
	ChargedBuffID                    uint16               // offset: 8, size: 2 bytes
	NUserNo                          uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_CHARGED_SETBUFF_CMD (size = 32 bytes)
type PROTO_NC_CHARGED_SETBUFF_CMD_DATA struct {
	CharRegistNumber                 uint32               // offset: 0, size: 4 bytes
	ItemKey                          uint64               // offset: 4, size: 8 bytes
	ItemId                           uint16               // offset: 12, size: 2 bytes
	ChargedBuff                      [14]byte             // offset: 14, size: 14 bytes
	NUserNo                          uint32               // offset: 28, size: 4 bytes
}

// PROTO_NC_CHARGED_STAT_INITIALIZE_DB_FAIL_ACK (size = 10 bytes)
type PROTO_NC_CHARGED_STAT_INITIALIZE_DB_FAIL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_CHARGED_STAT_INITIALIZE_DB_SUC_ACK (size = 8 bytes)
type PROTO_NC_CHARGED_STAT_INITIALIZE_DB_SUC_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHARGED_STAT_INITIALIZE_SUC_CMD (size = 1 bytes)
type PROTO_NC_CHARGED_STAT_INITIALIZE_SUC_CMD_DATA struct {
	Remain                           uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_ABSTATE_CMD (size = 10 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_ABSTATE_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Abstate                          uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_ALL_CMD (size = 107 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_ALL_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	All                              [105]byte            // offset: 2, size: 105 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_CHESTINFO_CMD (size = 10 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_CHESTINFO_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Chest                            uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_COININFO_CMD (size = 22 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_COININFO_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	NCoin                            uint64               // offset: 2, size: 8 bytes
	NExchangedCoin                   uint64               // offset: 10, size: 8 bytes
	NCharNo                          uint32               // offset: 18, size: 4 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_GAME_CMD (size = 11 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_GAME_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	Game                             [5]byte              // offset: 6, size: 5 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_LINK_FROM_CMD (size = 26 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_LINK_FROM_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Linkfrom                         [24]byte             // offset: 2, size: 24 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_MISC_CMD (size = 18 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_MISC_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Cen                              uint64               // offset: 2, size: 8 bytes
	Fame                             uint32               // offset: 10, size: 4 bytes
	NCharNo                          uint32               // offset: 14, size: 4 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_SINGLE_OPTION_CMD (size = 17 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_SINGLE_OPTION_CMD_DATA struct {
	WmHandle                         uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Resttime                         int32                // offset: 6, size: 4 bytes
	Autopick                         uint8                // offset: 10, size: 1 bytes
	Mystery_vault                    uint8                // offset: 11, size: 1 bytes
	NChatBan                         uint32               // offset: 12, size: 4 bytes
	BFirstJoin                       uint8                // offset: 16, size: 1 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_TITLE_CMD (size = 1037 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_TITLE_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	CTData                           [1031]byte           // offset: 6, size: 1031 bytes
}

// PROTO_NC_CHARSAVE_ABSTATE_REQ (size = 8 bytes)
type PROTO_NC_CHARSAVE_ABSTATE_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	Abstatenum                       uint16               // offset: 6, size: 2 bytes
	Abstate                          [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_CHARSAVE_ALL_ACK (size = 6 bytes)
type PROTO_NC_CHARSAVE_ALL_ACK_DATA struct {
	Bufferlocation                   uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHARSAVE_ALL_REQ (size = 105 bytes)
type PROTO_NC_CHARSAVE_ALL_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	Level                            [17]byte             // offset: 6, size: 17 bytes
	Location                         [48]byte             // offset: 23, size: 48 bytes
	Stat                             [34]byte             // offset: 71, size: 34 bytes
}

// PROTO_NC_CHARSAVE_CHARSTAT_CMD (size = 34 bytes)
type PROTO_NC_CHARSAVE_CHARSTAT_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	CurPwrStone                      uint32               // offset: 4, size: 4 bytes
	CurGrdStone                      uint32               // offset: 8, size: 4 bytes
	CurHPStone                       uint32               // offset: 12, size: 4 bytes
	CurSPStone                       uint32               // offset: 16, size: 4 bytes
	CurHP                            uint32               // offset: 20, size: 4 bytes
	CurSP                            uint32               // offset: 24, size: 4 bytes
	Statdistribute                   [6]byte              // offset: 28, size: 6 bytes
}

// PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_ACK (size = 8 bytes)
type PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_REQ (size = 7 bytes)
type PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	BtMysteryType                    uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_CHARSAVE_LEVEL_CMD (size = 17 bytes)
type PROTO_NC_CHARSAVE_LEVEL_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Level                            uint8                // offset: 4, size: 1 bytes
	Exp                              uint64               // offset: 5, size: 8 bytes
	Fame                             uint32               // offset: 13, size: 4 bytes
}

// PROTO_NC_CHARSAVE_LOCATION_CMD (size = 48 bytes)
type PROTO_NC_CHARSAVE_LOCATION_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	MapField                         Name3                // offset: 4, size: 12 bytes
	Coord                            uint64               // offset: 16, size: 8 bytes
	Kqhandle                         uint32               // offset: 24, size: 4 bytes
	Map_kq                           Name3                // offset: 28, size: 12 bytes
	Coord_kq                         uint64               // offset: 40, size: 8 bytes
}

// PROTO_NC_CHARSAVE_PKCOUNT_CMD (size = 8 bytes)
type PROTO_NC_CHARSAVE_PKCOUNT_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Pkcount                          uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHARSAVE_REST_EXP_LAST_EXEC_TIME_SAVE_REQ (size = 6 bytes)
type PROTO_NC_CHARSAVE_REST_EXP_LAST_EXEC_TIME_SAVE_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
}

// PROTO_NC_CHARSAVE_TITLE_REQ (size = 8 bytes)
type PROTO_NC_CHARSAVE_TITLE_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	CTData                           uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK (size = 2 bytes)
type PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ (size = 1 bytes)
type PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ_DATA struct {
	BtMysteryType                    uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_ABSTATE_CMD (size = 6 bytes)
type PROTO_NC_CHAR_ABSTATE_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Number                           uint16               // offset: 4, size: 2 bytes
	Abstate                          [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD (size = 1 bytes)
type PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD_DATA struct {
	NAdminLevel                      uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_ARENA_CMD (size = 1 bytes)
type PROTO_NC_CHAR_ARENA_CMD_DATA struct {}

// PROTO_NC_CHAR_BASEPARAMCHANGE_CMD (size = 1 bytes)
type PROTO_NC_CHAR_BASEPARAMCHANGE_CMD_DATA struct {
	Changenum                        uint8                // offset: 0, size: 1 bytes
	Param                            [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_CHAR_BASE_CMD (size = 105 bytes)
type PROTO_NC_CHAR_BASE_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Charid                           Name5                // offset: 4, size: 20 bytes
	Slotno                           uint8                // offset: 24, size: 1 bytes
	Level                            uint8                // offset: 25, size: 1 bytes
	Experience                       uint64               // offset: 26, size: 8 bytes
	CurPwrStone                      uint16               // offset: 34, size: 2 bytes
	CurGrdStone                      uint16               // offset: 36, size: 2 bytes
	CurHPStone                       uint16               // offset: 38, size: 2 bytes
	CurSPStone                       uint16               // offset: 40, size: 2 bytes
	CurHP                            uint32               // offset: 42, size: 4 bytes
	CurSP                            uint32               // offset: 46, size: 4 bytes
	CurLP                            uint32               // offset: 50, size: 4 bytes
	Fame                             uint32               // offset: 54, size: 4 bytes
	Cen                              uint64               // offset: 58, size: 8 bytes
	Logininfo                        [21]byte             // offset: 66, size: 21 bytes
	Statdistribute                   [6]byte              // offset: 87, size: 6 bytes
	Pkyellowtime                     uint8                // offset: 93, size: 1 bytes
	Pkcount                          uint32               // offset: 94, size: 4 bytes
	Prisonmin                        uint16               // offset: 98, size: 2 bytes
	Adminlevel                       uint8                // offset: 100, size: 1 bytes
	Flags                            uint32               // offset: 101, size: 4 bytes
}

// PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Number                           uint16               // offset: 4, size: 2 bytes
	ViewList                         [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_CARDCOLLECT_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CARDCOLLECT_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Number                           uint16               // offset: 4, size: 2 bytes
	CardList                         [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_CARDCOLLECT_REWARD_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CARDCOLLECT_REWARD_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Number                           uint16               // offset: 4, size: 2 bytes
	RewardID                         [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_CENCHANGE_CMD (size = 8 bytes)
type PROTO_NC_CHAR_CENCHANGE_CMD_DATA struct {
	Cen                              uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_CHAR_CHANGEBYCONDITION_PARAM_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CHANGEBYCONDITION_PARAM_CMD_DATA struct {
	NSkillID                         uint16               // offset: 0, size: 2 bytes
	NChangeRate                      uint16               // offset: 2, size: 2 bytes
	NParamNum                        uint16               // offset: 4, size: 2 bytes
	AParam                           [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_CHARDATAFAIL_ACK (size = 4 bytes)
type PROTO_NC_CHAR_CHARDATAFAIL_ACK_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAR_CHARDATA_ACK (size = 2 bytes)
type PROTO_NC_CHAR_CHARDATA_ACK_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_CHARDATA_REQ (size = 14 bytes)
type PROTO_NC_CHAR_CHARDATA_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	NUserNo                          uint32               // offset: 6, size: 4 bytes
	TExchangedCoinSumDate            int32                // offset: 10, size: 4 bytes
}

// PROTO_NC_CHAR_CHARGEDBUFF_CMD (size = 2 bytes)
type PROTO_NC_CHAR_CHARGEDBUFF_CMD_DATA struct {
	NumOfChargedBuff                 uint16               // offset: 0, size: 2 bytes
	ChargedBuff                      [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ (size = 4 bytes)
type PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ_DATA struct {
	NChargedBuffKey                  uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_CHESTINFO_CMD (size = 2 bytes)
type PROTO_NC_CHAR_CHESTINFO_CMD_DATA struct {
	Itemnum                          uint16               // offset: 0, size: 2 bytes
	Location                         [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_CHESTINFO_REQ (size = 8 bytes)
type PROTO_NC_CHAR_CHESTINFO_REQ struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	Number                           uint16               // offset: 6, size: 2 bytes
	Loc                              [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_CHAR_CLASSCHANGE_ACK (size = 7 bytes)
type PROTO_NC_CHAR_CLASSCHANGE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Newclass                         uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_CHAR_CLASSCHANGE_CMD (size = 9 bytes)
type PROTO_NC_CHAR_CLASSCHANGE_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Wldmanhandle                     uint16               // offset: 6, size: 2 bytes
	Newclass                         uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_CHAR_CLASSCHANGE_REQ (size = 8 bytes)
type PROTO_NC_CHAR_CLASSCHANGE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Newclass                         uint8                // offset: 6, size: 1 bytes
	Level                            uint8                // offset: 7, size: 1 bytes
}

// PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK (size = 3 bytes)
type PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	BEnable                          uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD (size = 3 bytes)
type PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	BEnable                          uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ (size = 1 bytes)
type PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ_DATA struct {
	BEnable                          uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD_DATA struct {
	CurrentTitle                     uint8                // offset: 0, size: 1 bytes
	CurrentTitleElement              uint8                // offset: 1, size: 1 bytes
	CurrentTitleMobID                uint16               // offset: 2, size: 2 bytes
	NumOfTitle                       uint16               // offset: 4, size: 2 bytes
	TitleArray                       [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_COININFO_CMD (size = 16 bytes)
type PROTO_NC_CHAR_COININFO_CMD_DATA struct {
	NCoin                            uint64               // offset: 0, size: 8 bytes
	NExchangedCoin                   uint64               // offset: 8, size: 8 bytes
}

// PROTO_NC_CHAR_DB_AUTO_PICK_CMD (size = 5 bytes)
type PROTO_NC_CHAR_DB_AUTO_PICK_CMD struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	BEnable                          uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_CHAR_DB_NEWBIE_GUIDE_VIEW_SET_CMD (size = 5 bytes)
type PROTO_NC_CHAR_DB_NEWBIE_GUIDE_VIEW_SET_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NGuideView                       uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_CHAR_DB_REST_EXP_LAST_EXEC_TIME_CMD (size = 8 bytes)
type PROTO_NC_CHAR_DB_REST_EXP_LAST_EXEC_TIME_CMD struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DRestExpLastExecTime             int32                // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_DEADMENU_CMD (size = 9 bytes)
type PROTO_NC_CHAR_DEADMENU_CMD_DATA struct {
	Second                           uint32               // offset: 0, size: 4 bytes
	Priority                         uint8                // offset: 4, size: 1 bytes
	EMenuType                        uint32               // offset: 5, size: 4 bytes
}

// PROTO_NC_CHAR_EMBLEM_CMD (size = 1 bytes)
type PROTO_NC_CHAR_EMBLEM_CMD_DATA struct {}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_CHANGE_DAY_CMD (size = 1 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_CHANGE_DAY_CMD_DATA struct {
	NChangeIndex                     uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_CMD (size = 1 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_CMD_DATA struct {}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_DB_ACK (size = 6 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_DB_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NError                           uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_DB_REQ (size = 40 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_CHECK_DB_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TCheckTime                       [36]byte             // offset: 4, size: 36 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_CLIENT_CMD (size = 66 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_CLIENT_CMD_DATA struct {
	TEventStart                      [36]byte             // offset: 0, size: 36 bytes
	NCurIndex                        uint8                // offset: 36, size: 1 bytes
	AttendanceArray                  [28]byte             // offset: 37, size: 28 bytes
	RewardNum                        uint8                // offset: 65, size: 1 bytes
	RewardList                       [0]byte              // offset: 66, flexible array member
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_DB_ACK (size = 1020 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_DB_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NError                           uint16               // offset: 4, size: 2 bytes
	TConnect                         int32                // offset: 6, size: 4 bytes
	AttendanceNum                    uint8                // offset: 10, size: 1 bytes
	AttendanceArray                  [1008]byte           // offset: 11, size: 1008 bytes
	RewardNum                        uint8                // offset: 1019, size: 1 bytes
	RewardList                       [0]byte              // offset: 1020, flexible array member
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_DB_REQ (size = 80 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_LIST_DB_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TmEventStartDate                 [36]byte             // offset: 4, size: 36 bytes
	TmEventEndDate                   [36]byte             // offset: 40, size: 36 bytes
	TConnect                         int32                // offset: 76, size: 4 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK (size = 4 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK_DATA struct {
	NIndex                           uint8                // offset: 0, size: 1 bytes
	NRewardID                        uint8                // offset: 1, size: 1 bytes
	NError                           uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_DB_ACK (size = 8 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_DB_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NIndex                           uint8                // offset: 4, size: 1 bytes
	NRewardID                        uint8                // offset: 5, size: 1 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_DB_REQ (size = 114 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_DB_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NIndex                           uint8                // offset: 4, size: 1 bytes
	TmEventDay                       [36]byte             // offset: 5, size: 36 bytes
	NRewardID                        uint8                // offset: 41, size: 1 bytes
	TmCheckStartDay                  [36]byte             // offset: 42, size: 36 bytes
	TmCheckEndDay                    [36]byte             // offset: 78, size: 36 bytes
}

// PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ (size = 2 bytes)
type PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ_DATA struct {
	NIndex                           uint8                // offset: 0, size: 1 bytes
	NRewardID                        uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_CHAR_EXP_CHANGED_CMD (size = 14 bytes)
type PROTO_NC_CHAR_EXP_CHANGED_CMD_DATA struct {
	Wmhandle                         uint16               // offset: 0, size: 2 bytes
	CharNo                           uint32               // offset: 2, size: 4 bytes
	CurrentExp                       uint64               // offset: 6, size: 8 bytes
}

// PROTO_NC_CHAR_FAMECHANGE_CMD (size = 4 bytes)
type PROTO_NC_CHAR_FAMECHANGE_CMD_DATA struct {
	Fame                             uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_FAMESAVE_CMD (size = 8 bytes)
type PROTO_NC_CHAR_FAMESAVE_CMD_DATA struct {
	Charregnum                       uint32               // offset: 0, size: 4 bytes
	Fame                             uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_FREESTAT_SET_DB_ACK (size = 15 bytes)
type PROTO_NC_CHAR_FREESTAT_SET_DB_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NFreeStat                        [5]byte              // offset: 6, size: 5 bytes
	NLockIndex                       uint16               // offset: 11, size: 2 bytes
	NErrorType                       uint16               // offset: 13, size: 2 bytes
}

// PROTO_NC_CHAR_FREESTAT_SET_DB_REQ (size = 15 bytes)
type PROTO_NC_CHAR_FREESTAT_SET_DB_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NMaxStat                         uint8                // offset: 6, size: 1 bytes
	NFreeStat                        [5]byte              // offset: 7, size: 5 bytes
	NLockIndex                       uint16               // offset: 12, size: 2 bytes
	NItemNum                         uint8                // offset: 14, size: 1 bytes
	DecItemInfo                      [0]byte              // offset: 15, flexible array member
}

// PROTO_NC_CHAR_HOUSE_CMD (size = 1 bytes)
type PROTO_NC_CHAR_HOUSE_CMD_DATA struct {}

// PROTO_NC_CHAR_KICKPLAYEROUT_RNG (size = 28 bytes)
type PROTO_NC_CHAR_KICKPLAYEROUT_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Tmp                              uint8                // offset: 7, size: 1 bytes
	Charid                           Name5                // offset: 8, size: 20 bytes
}

// PROTO_NC_CHAR_LEVEL_CHANGED_CMD (size = 7 bytes)
type PROTO_NC_CHAR_LEVEL_CHANGED_CMD_DATA struct {
	Wmhandle                         uint16               // offset: 0, size: 2 bytes
	CharNo                           uint32               // offset: 2, size: 4 bytes
	CharNewLevel                     uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_CHAR_LINKFROM_CMD (size = 24 bytes)
type PROTO_NC_CHAR_LINKFROM_CMD struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	MapField                         Name3                // offset: 4, size: 12 bytes
	Coord                            uint64               // offset: 16, size: 8 bytes
}

// PROTO_NC_CHAR_LOGINFAIL_ACK (size = 2 bytes)
type PROTO_NC_CHAR_LOGINFAIL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_LOGIN_ACK (size = 18 bytes)
type PROTO_NC_CHAR_LOGIN_ACK_DATA struct {
	Zoneip                           [16]byte             // offset: 0, size: 16 bytes
	Zoneport                         uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_CHAR_LOGIN_DB (size = 56 bytes)
type PROTO_NC_CHAR_LOGIN_DB_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           [20]byte             // offset: 4, size: 20 bytes
	Cen                              uint64               // offset: 24, size: 8 bytes
	Usrregnum                        uint32               // offset: 32, size: 4 bytes
	Spawnapps                        Name5                // offset: 36, size: 20 bytes
}

// PROTO_NC_CHAR_LOGIN_REQ (size = 1 bytes)
type PROTO_NC_CHAR_LOGIN_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_LOGOUT_DB (size = 35 bytes)
type PROTO_NC_CHAR_LOGOUT_DB_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           [20]byte             // offset: 4, size: 20 bytes
	Playmin                          uint16               // offset: 24, size: 2 bytes
	State                            uint8                // offset: 26, size: 1 bytes
	Cen                              uint64               // offset: 27, size: 8 bytes
}

// PROTO_NC_CHAR_MAPLOGIN_ACK (size = 242 bytes)
type PROTO_NC_CHAR_MAPLOGIN_ACK struct {
	Charhandle                       uint16               // offset: 0, size: 2 bytes
	Param                            [232]byte            // offset: 2, size: 232 bytes
	Logincoord                       uint64               // offset: 234, size: 8 bytes
}

// PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD (size = 1 bytes)
type PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD_DATA struct {
	Mystery_vault                    uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_LIST_CMD (size = 1 bytes)
type PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_LIST_CMD_DATA struct {
	NGuideViewNum                    uint8                // offset: 0, size: 1 bytes
	NGuideViewList                   [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD (size = 1 bytes)
type PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD_DATA struct {
	NGuideView                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_PLAYERBANNED_RNG (size = 27 bytes)
type PROTO_NC_CHAR_PLAYERBANNED_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Charid                           Name5                // offset: 7, size: 20 bytes
}

// PROTO_NC_CHAR_PLAYERFOUND_BY_NORMAL_USER_RNG (size = 49 bytes)
type PROTO_NC_CHAR_PLAYERFOUND_BY_NORMAL_USER_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	MapField                         Name3                // offset: 7, size: 12 bytes
	Coord                            uint64               // offset: 19, size: 8 bytes
	Charid                           Name5                // offset: 27, size: 20 bytes
	CanRecall                        uint8                // offset: 47, size: 1 bytes
	ModeIDLv                         uint8                // offset: 48, size: 1 bytes
}

// PROTO_NC_CHAR_PLAYERFOUND_RNG (size = 47 bytes)
type PROTO_NC_CHAR_PLAYERFOUND_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	MapField                         Name3                // offset: 7, size: 12 bytes
	Coord                            uint64               // offset: 19, size: 8 bytes
	Charid                           Name5                // offset: 27, size: 20 bytes
}

// PROTO_NC_CHAR_PLAYERSEARCH_BY_NORMAL_USER_RNG (size = 27 bytes)
type PROTO_NC_CHAR_PLAYERSEARCH_BY_NORMAL_USER_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Charid                           Name5                // offset: 7, size: 20 bytes
}

// PROTO_NC_CHAR_PLAYERSEARCH_RNG (size = 27 bytes)
type PROTO_NC_CHAR_PLAYERSEARCH_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Charid                           Name5                // offset: 7, size: 20 bytes
}

// PROTO_NC_CHAR_PLAYERSUMMON_RNG (size = 47 bytes)
type PROTO_NC_CHAR_PLAYERSUMMON_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Charid                           Name5                // offset: 7, size: 20 bytes
	MapField                         Name3                // offset: 27, size: 12 bytes
	Coord                            uint64               // offset: 39, size: 8 bytes
}

// PROTO_NC_CHAR_POLYMORPH_CMD (size = 8 bytes)
type PROTO_NC_CHAR_POLYMORPH_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	Polymorph                        uint16               // offset: 2, size: 2 bytes
	AbstateIndex                     uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_PROMOTE_ACK (size = 1 bytes)
type PROTO_NC_CHAR_PROMOTE_ACK_DATA struct {
	Newclass                         uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_REGISTNUMBER_ACK (size = 30 bytes)
type PROTO_NC_CHAR_REGISTNUMBER_ACK_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
	Chrregnum                        uint32               // offset: 4, size: 4 bytes
	Loginreq                         [22]byte             // offset: 8, size: 22 bytes
}

// PROTO_NC_CHAR_REGISTNUMBER_ACK_SEND (size = 39 bytes)
type PROTO_NC_CHAR_REGISTNUMBER_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 3, size: 6 bytes
	Ack                              [30]byte             // offset: 9, size: 30 bytes
}

// PROTO_NC_CHAR_REGISTNUMBER_REQ (size = 22 bytes)
type PROTO_NC_CHAR_REGISTNUMBER_REQ_DATA struct {
	Loginreq                         [22]byte             // offset: 0, size: 22 bytes
}

// PROTO_NC_CHAR_REGNUM_VARIFICATION_ACK (size = 7 bytes)
type PROTO_NC_CHAR_REGNUM_VARIFICATION_ACK_DATA struct {
	Zonehandle                       uint16               // offset: 0, size: 2 bytes
	Charregnum                       uint32               // offset: 2, size: 4 bytes
	Iscorrect                        uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_CHAR_REGNUM_VARIFICATION_REQ (size = 26 bytes)
type PROTO_NC_CHAR_REGNUM_VARIFICATION_REQ_DATA struct {
	Zonehandle                       uint16               // offset: 0, size: 2 bytes
	Zone_charid                      Name5                // offset: 2, size: 20 bytes
	Charregnum                       uint32               // offset: 22, size: 4 bytes
}

// PROTO_NC_CHAR_REVIVEOTHER_CMD (size = 30 bytes)
type PROTO_NC_CHAR_REVIVEOTHER_CMD_DATA struct {
	Link                             [10]byte             // offset: 0, size: 10 bytes
	Sock                             [18]byte             // offset: 10, size: 18 bytes
	Wldmanhandle                     uint16               // offset: 28, size: 2 bytes
}

// PROTO_NC_CHAR_REVIVESAME_CMD (size = 10 bytes)
type PROTO_NC_CHAR_REVIVESAME_CMD_DATA struct {
	Mapid                            uint16               // offset: 0, size: 2 bytes
	Location                         uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_CHAR_SEAWAR_CMD (size = 1 bytes)
type PROTO_NC_CHAR_SEAWAR_CMD_DATA struct {}

// PROTO_NC_CHAR_SET_STYLE_ACK (size = 2 bytes)
type PROTO_NC_CHAR_SET_STYLE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_SET_STYLE_DB_ACK (size = 12 bytes)
type PROTO_NC_CHAR_SET_STYLE_DB_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	WMHandle                         uint16               // offset: 2, size: 2 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	Info                             uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_CHAR_SET_STYLE_DB_REQ (size = 10 bytes)
type PROTO_NC_CHAR_SET_STYLE_DB_REQ_DATA struct {
	WMHandle                         uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	Info                             uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_CHAR_SET_STYLE_GET_INFO_ACK (size = 26 bytes)
type PROTO_NC_CHAR_SET_STYLE_GET_INFO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Items                            [24]byte             // offset: 2, size: 24 bytes
}

// PROTO_NC_CHAR_SET_STYLE_GET_INFO_DB_ACK (size = 32 bytes)
type PROTO_NC_CHAR_SET_STYLE_GET_INFO_DB_ACK_DATA struct {
	WMHandle                         uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	Items                            [24]byte             // offset: 8, size: 24 bytes
}

// PROTO_NC_CHAR_SET_STYLE_GET_INFO_DB_REQ (size = 6 bytes)
type PROTO_NC_CHAR_SET_STYLE_GET_INFO_DB_REQ_DATA struct {
	WMHandle                         uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ (size = 1 bytes)
type PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ_DATA struct {}

// PROTO_NC_CHAR_SET_STYLE_REQ (size = 4 bytes)
type PROTO_NC_CHAR_SET_STYLE_REQ_DATA struct {
	Info                             uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_SINGLE_OPTION_CMD (size = 15 bytes)
type PROTO_NC_CHAR_SINGLE_OPTION_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Resttime                         int32                // offset: 4, size: 4 bytes
	Autopick                         uint8                // offset: 8, size: 1 bytes
	Mystery_vault                    uint8                // offset: 9, size: 1 bytes
	NChatBan                         uint32               // offset: 10, size: 4 bytes
	BFirstJoin                       uint8                // offset: 14, size: 1 bytes
}

// PROTO_NC_CHAR_SOMEONEPROMOTE_CMD (size = 3 bytes)
type PROTO_NC_CHAR_SOMEONEPROMOTE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Newclass                         uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_DECPOINTFAIL_DB_ACK (size = 10 bytes)
type PROTO_NC_CHAR_STAT_DECPOINTFAIL_DB_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_CHAR_STAT_DECPOINTSUC_ACK (size = 1 bytes)
type PROTO_NC_CHAR_STAT_DECPOINTSUC_ACK_DATA struct {
	Stat                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_DECPOINTSUC_DB_ACK (size = 9 bytes)
type PROTO_NC_CHAR_STAT_DECPOINTSUC_DB_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Stat                             uint8                // offset: 6, size: 1 bytes
	Lockindex                        uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_CHAR_STAT_DECPOINT_REQ (size = 1 bytes)
type PROTO_NC_CHAR_STAT_DECPOINT_REQ_DATA struct {
	Stat                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_INCPOINTFAIL_DB_ACK (size = 8 bytes)
type PROTO_NC_CHAR_STAT_INCPOINTFAIL_DB_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHAR_STAT_INCPOINTSUC_ACK (size = 1 bytes)
type PROTO_NC_CHAR_STAT_INCPOINTSUC_ACK_DATA struct {
	Stat                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_INCPOINTSUC_DB_ACK (size = 7 bytes)
type PROTO_NC_CHAR_STAT_INCPOINTSUC_DB_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Stat                             uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_INCPOINT_REQ (size = 1 bytes)
type PROTO_NC_CHAR_STAT_INCPOINT_REQ_DATA struct {
	Stat                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_STAT_REMAINPOINT_CMD (size = 1 bytes)
type PROTO_NC_CHAR_STAT_REMAINPOINT_CMD_DATA struct {
	Remain                           uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_SUPPORT_REWARD_CMD (size = 8 bytes)
type PROTO_NC_CHAR_SUPPORT_REWARD_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	EUserType                        uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_ACK (size = 6 bytes)
type PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NError                           uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_REQ (size = 10 bytes)
type PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NMaxStat                         uint8                // offset: 4, size: 1 bytes
	NFreeStat                        [5]byte              // offset: 5, size: 5 bytes
}

// PROTO_NC_CHAR_TUTORIAL_INFO_WORLD_CMD (size = 9 bytes)
type PROTO_NC_CHAR_TUTORIAL_INFO_WORLD_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TutorialInfo                     [5]byte              // offset: 4, size: 5 bytes
}

// PROTO_NC_CHAR_TUTORIAL_INFO_ZONE_CMD (size = 11 bytes)
type PROTO_NC_CHAR_TUTORIAL_INFO_ZONE_CMD_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	TutorialInfo                     [5]byte              // offset: 6, size: 5 bytes
}

// PROTO_NC_CHAR_TUTORIAL_POPUP_ACK (size = 1 bytes)
type PROTO_NC_CHAR_TUTORIAL_POPUP_ACK_DATA struct {
	BIsSkip                          uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_TUTORIAL_POPUP_REQ (size = 1 bytes)
type PROTO_NC_CHAR_TUTORIAL_POPUP_REQ_DATA struct {}

// PROTO_NC_CHAR_TUTORIAL_STEP_ACK (size = 1 bytes)
type PROTO_NC_CHAR_TUTORIAL_STEP_ACK_DATA struct {
	NTutorialStep                    uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_TUTORIAL_STEP_REQ (size = 1 bytes)
type PROTO_NC_CHAR_TUTORIAL_STEP_REQ_DATA struct {
	NTutorialStep                    uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK (size = 9 bytes)
type PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TutorialInfo                     [5]byte              // offset: 4, size: 5 bytes
}

// PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_CMD (size = 9 bytes)
type PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TutorialInfo                     [5]byte              // offset: 4, size: 5 bytes
}

// PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_REQ (size = 13 bytes)
type PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NUserNo                          uint32               // offset: 4, size: 4 bytes
	TutorialInfo                     [5]byte              // offset: 8, size: 5 bytes
}

// PROTO_NC_CHAR_WEDDINGDATA_ACK (size = 43 bytes)
type PROTO_NC_CHAR_WEDDINGDATA_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	Status                           uint8                // offset: 8, size: 1 bytes
	When                             uint32               // offset: 9, size: 4 bytes
	SubWhen                          [0]byte              // offset: 13, size: 0 bytes
	ApplicationLimit                 [0]byte              // offset: 13, size: 0 bytes
	DelayLimit                       [0]byte              // offset: 13, size: 0 bytes
	DivorceLimit                     [0]byte              // offset: 13, size: 0 bytes
	NewlyLimit                       uint32               // offset: 13, size: 4 bytes
	Partner                          [24]byte             // offset: 17, size: 24 bytes
	WeddingCount                     uint8                // offset: 41, size: 1 bytes
	RefusalCount                     uint8                // offset: 42, size: 1 bytes
}

// PROTO_NC_CHAR_WEDDINGDATA_REQ (size = 6 bytes)
type PROTO_NC_CHAR_WEDDINGDATA_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_CHAR_WEDDING_CANCEL_WEDDING (size = 4 bytes)
type PROTO_NC_CHAR_WEDDING_CANCEL_WEDDING_DATA struct {
	RegisterNumber                   uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_ACK (size = 16 bytes)
type PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Success                          int32                // offset: 6, size: 4 bytes
	MySubStatus                      uint8                // offset: 10, size: 1 bytes
	Partner                          uint32               // offset: 11, size: 4 bytes
	PartnerSubStatus                 uint8                // offset: 15, size: 1 bytes
}

// PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_REQ (size = 11 bytes)
type PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Partner                          uint32               // offset: 6, size: 4 bytes
	IncreseCount                     uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_CHAR_WEDDING_DIVORCE_DO_REQ (size = 11 bytes)
type PROTO_NC_CHAR_WEDDING_DIVORCE_DO_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Partner                          uint32               // offset: 6, size: 4 bytes
	Reason                           uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK (size = 14 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	PartnerInfo                      uint64               // offset: 6, size: 8 bytes
}

// PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ (size = 10 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Partner                          uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG (size = 15 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	PartnerInfo                      uint64               // offset: 7, size: 8 bytes
}

// PROTO_NC_CHAR_WEDDING_PARTNER_SUMMON_RNG (size = 35 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_SUMMON_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	RegisterNumber                   uint32               // offset: 7, size: 4 bytes
	MapField                         Name3                // offset: 11, size: 12 bytes
	Coord                            uint64               // offset: 23, size: 8 bytes
	From                             uint8                // offset: 31, size: 1 bytes
	To                               uint8                // offset: 32, size: 1 bytes
	Err                              uint16               // offset: 33, size: 2 bytes
}

// PROTO_NC_CHAR_WEDDING_PROPOSE_ACK (size = 14 bytes)
type PROTO_NC_CHAR_WEDDING_PROPOSE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Success                          int32                // offset: 6, size: 4 bytes
	Partner                          uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_CHAR_WEDDING_PROPOSE_REQ (size = 10 bytes)
type PROTO_NC_CHAR_WEDDING_PROPOSE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Partner                          uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_CHAR_WEDDING_REFRESH_INFO_RNG (size = 7 bytes)
type PROTO_NC_CHAR_WEDDING_REFRESH_INFO_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
}

// PROTO_NC_CHAR_ZONE_AUTO_PICK_CMD (size = 9 bytes)
type PROTO_NC_CHAR_ZONE_AUTO_PICK_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Wldmanhandle                     uint16               // offset: 6, size: 2 bytes
	BEnable                          uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_CHAR_ZONE_CHARDATAFAIL_ACK_SEND (size = 11 bytes)
type PROTO_NC_CHAR_ZONE_CHARDATAFAIL_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 3, size: 6 bytes
	Ack                              uint16               // offset: 9, size: 2 bytes
}

// PROTO_NC_CHAR_ZONE_CHARDATA_ACK (size = 2 bytes)
type PROTO_NC_CHAR_ZONE_CHARDATA_ACK_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_ZONE_CHARDATA_ACK_SEND (size = 11 bytes)
type PROTO_NC_CHAR_ZONE_CHARDATA_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 3, size: 6 bytes
	Ack                              uint16               // offset: 9, size: 2 bytes
}

// PROTO_NC_CHAR_ZONE_CHARDATA_REQ (size = 22 bytes)
type PROTO_NC_CHAR_ZONE_CHARDATA_REQ_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Charid                           Name5                // offset: 2, size: 20 bytes
}

// PROTO_NC_CHER_EVENT_ATTENDANCE_CHANGE_START_CMD (size = 36 bytes)
type PROTO_NC_CHER_EVENT_ATTENDANCE_CHANGE_START_CMD_DATA struct {
	TEventStart                      [36]byte             // offset: 0, size: 36 bytes
}

// PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK (size = 6 bytes)
type PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK_DATA struct {
	BookMarkSlot                     uint16               // offset: 0, size: 2 bytes
	CardID                           uint16               // offset: 2, size: 2 bytes
	Err                              uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_ACK (size = 12 bytes)
type PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	BookMarkSlot                     uint16               // offset: 6, size: 2 bytes
	CardID                           uint16               // offset: 8, size: 2 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_REQ (size = 10 bytes)
type PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	BookMarkSlot                     uint16               // offset: 6, size: 2 bytes
	CardID                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ (size = 4 bytes)
type PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ_DATA struct {
	BookMarkSlot                     uint16               // offset: 0, size: 2 bytes
	CardID                           uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_COLLECT_CARDOPEN_CMD (size = 106 bytes)
type PROTO_NC_COLLECT_CARDOPEN_CMD_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Slot                             uint8                // offset: 2, size: 1 bytes
	CardInform                       [103]byte            // offset: 3, size: 103 bytes
}

// PROTO_NC_COLLECT_CARDREGIST_ACK (size = 10 bytes)
type PROTO_NC_COLLECT_CARDREGIST_ACK_DATA struct {
	Invenslot                        uint8                // offset: 0, size: 1 bytes
	CardID                           uint16               // offset: 1, size: 2 bytes
	Star                             uint8                // offset: 3, size: 1 bytes
	Serial                           uint32               // offset: 4, size: 4 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_COLLECT_CARDREGIST_REQ (size = 1 bytes)
type PROTO_NC_COLLECT_CARDREGIST_REQ_DATA struct {
	Invenslot                        uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_COLLECT_PEEPING_COLLECT_ACK (size = 28 bytes)
type PROTO_NC_COLLECT_PEEPING_COLLECT_ACK_DATA struct {
	Charname                         Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
	Collection                       [6]byte              // offset: 22, size: 6 bytes
}

// PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK (size = 28 bytes)
type PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK_DATA struct {
	Charname                         Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
	Collectionview                   [6]byte              // offset: 22, size: 6 bytes
}

// PROTO_NC_COLLECT_PEEPING_DB_COLLECT_ACK (size = 34 bytes)
type PROTO_NC_COLLECT_PEEPING_DB_COLLECT_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Targetcharname                   Name5                // offset: 6, size: 20 bytes
	Err                              uint16               // offset: 26, size: 2 bytes
	Collection                       [6]byte              // offset: 28, size: 6 bytes
}

// PROTO_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK (size = 34 bytes)
type PROTO_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Targetcharname                   Name5                // offset: 6, size: 20 bytes
	Err                              uint16               // offset: 26, size: 2 bytes
	Collectionview                   [6]byte              // offset: 28, size: 6 bytes
}

// PROTO_NC_COLLECT_PEEPING_DB_REQ (size = 26 bytes)
type PROTO_NC_COLLECT_PEEPING_DB_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Targetcharname                   Name5                // offset: 6, size: 20 bytes
}

// PROTO_NC_COLLECT_PEEPING_REQ (size = 20 bytes)
type PROTO_NC_COLLECT_PEEPING_REQ_DATA struct {
	Charname                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_COLLECT_REGIST_ACK (size = 18 bytes)
type PROTO_NC_COLLECT_REGIST_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	CardInfo                         [10]byte             // offset: 8, size: 10 bytes
}

// PROTO_NC_COLLECT_REGIST_REQ (size = 26 bytes)
type PROTO_NC_COLLECT_REGIST_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	CardKey                          uint64               // offset: 8, size: 8 bytes
	CardInfo                         [10]byte             // offset: 16, size: 10 bytes
}

// PROTO_NC_COLLECT_REWARD_ACK (size = 4 bytes)
type PROTO_NC_COLLECT_REWARD_ACK_DATA struct {
	RewardID                         uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_COLLECT_REWARD_DB_ACK (size = 10 bytes)
type PROTO_NC_COLLECT_REWARD_DB_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	RewardID                         uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_COLLECT_REWARD_DB_REQ (size = 119 bytes)
type PROTO_NC_COLLECT_REWARD_DB_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	RewardID                         uint16               // offset: 6, size: 2 bytes
	RewardItem                       [111]byte            // offset: 8, size: 111 bytes
}

// PROTO_NC_COLLECT_REWARD_REQ (size = 4 bytes)
type PROTO_NC_COLLECT_REWARD_REQ_DATA struct {
	RewardID                         uint16               // offset: 0, size: 2 bytes
	RewardItem                       uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CT_CHARTTING_CMD (size = 6 bytes)
type PROTO_NC_CT_CHARTTING_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_CT_LUASCRIPT_SET_WORLD_CMD (size = 6 bytes)
type PROTO_NC_CT_LUASCRIPT_SET_WORLD_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	CTInfo                           uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CT_LUASCRIPT_SET_ZONE_CMD (size = 6 bytes)
type PROTO_NC_CT_LUASCRIPT_SET_ZONE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	CTInfo                           uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_DATA_REPORT_ADD_ACK (size = 28 bytes)
type PROTO_NC_DATA_REPORT_ADD_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	HReporter                        uint16               // offset: 2, size: 2 bytes
	SHarmer                          Name5                // offset: 4, size: 20 bytes
	NRegNumber                       int32                // offset: 24, size: 4 bytes
}

// PROTO_NC_DATA_REPORT_ADD_REQ (size = 2369 bytes)
type PROTO_NC_DATA_REPORT_ADD_REQ_DATA struct {
	DwReporter                       uint32               // offset: 0, size: 4 bytes
	HReporter                        uint16               // offset: 4, size: 2 bytes
	NWorldNo                         uint8                // offset: 6, size: 1 bytes
	SReporterCharID                  Name5                // offset: 7, size: 20 bytes
	SHarmer                          Name5                // offset: 27, size: 20 bytes
	ByReportType                     [16]byte             // offset: 47, size: 16 bytes
	ByReason                         [256]byte            // offset: 63, size: 256 bytes
	NChatLogSize                     uint16               // offset: 319, size: 2 bytes
	ByChatLog                        [2048]byte           // offset: 321, size: 2048 bytes
}

// PROTO_NC_DATA_REPORT_CANCEL_ACK (size = 8 bytes)
type PROTO_NC_DATA_REPORT_CANCEL_ACK_DATA struct {
	NRegNumber                       int32                // offset: 0, size: 4 bytes
	HReporter                        uint16               // offset: 4, size: 2 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_DATA_REPORT_CANCEL_REQ (size = 6 bytes)
type PROTO_NC_DATA_REPORT_CANCEL_REQ_DATA struct {
	NRegNumber                       int32                // offset: 0, size: 4 bytes
	HReporter                        uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_DATA_REPORT_GET_ACK (size = 137 bytes)
type PROTO_NC_DATA_REPORT_GET_ACK_DATA struct {
	HReporter                        uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
	NRegNumber                       int32                // offset: 4, size: 4 bytes
	ByReportResult                   uint8                // offset: 8, size: 1 bytes
	ByRemark                         [128]byte            // offset: 9, size: 128 bytes
}

// PROTO_NC_DATA_REPORT_GET_REQ (size = 7 bytes)
type PROTO_NC_DATA_REPORT_GET_REQ_DATA struct {
	HReporter                        uint16               // offset: 0, size: 2 bytes
	DwReporter                       uint32               // offset: 2, size: 4 bytes
	NWorldNo                         uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_EVENT_ADD_EVENT_ACK (size = 1365 bytes)
type PROTO_NC_EVENT_ADD_EVENT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	EventData                        [1363]byte           // offset: 2, size: 1363 bytes
}

// PROTO_NC_EVENT_ADD_EVENT_REQ (size = 1363 bytes)
type PROTO_NC_EVENT_ADD_EVENT_REQ_DATA struct {
	EventData                        [1363]byte           // offset: 0, size: 1363 bytes
}

// PROTO_NC_EVENT_ADD_UPDATE_EVENT_CMD (size = 1363 bytes)
type PROTO_NC_EVENT_ADD_UPDATE_EVENT_CMD_DATA struct {
	EventData                        [1363]byte           // offset: 0, size: 1363 bytes
}

// PROTO_NC_EVENT_DEL_EVENT_ACK (size = 4 bytes)
type PROTO_NC_EVENT_DEL_EVENT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NEventNo                         uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_EVENT_DEL_EVENT_CMD (size = 2 bytes)
type PROTO_NC_EVENT_DEL_EVENT_CMD_DATA struct {
	NEventNo                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_EVENT_DEL_EVENT_REQ (size = 2 bytes)
type PROTO_NC_EVENT_DEL_EVENT_REQ_DATA struct {
	NEventNo                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_EVENT_GET_ALL_EVENT_INFO_ACK (size = 1370 bytes)
type PROTO_NC_EVENT_GET_ALL_EVENT_INFO_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NAllEventCnt                     uint16               // offset: 2, size: 2 bytes
	NPacketOrder                     uint16               // offset: 4, size: 2 bytes
	EventData                        [1363]byte           // offset: 6, size: 1363 bytes
	BIsFlagAllReady                  uint8                // offset: 1369, size: 1 bytes
}

// PROTO_NC_EVENT_GET_ALL_EVENT_INFO_REQ (size = 1 bytes)
type PROTO_NC_EVENT_GET_ALL_EVENT_INFO_REQ_DATA struct {}

// PROTO_NC_EVENT_SET_ALL_READY_ACK (size = 4 bytes)
type PROTO_NC_EVENT_SET_ALL_READY_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NEventNo                         uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_EVENT_SET_ALL_READY_REQ (size = 3 bytes)
type PROTO_NC_EVENT_SET_ALL_READY_REQ_DATA struct {
	NEventNo                         uint16               // offset: 0, size: 2 bytes
	BIsFlagAllReady                  uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_EVENT_UPDATE_EVENT_ACK (size = 1365 bytes)
type PROTO_NC_EVENT_UPDATE_EVENT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	EventData                        [1363]byte           // offset: 2, size: 1363 bytes
}

// PROTO_NC_EVENT_UPDATE_EVENT_REQ (size = 1363 bytes)
type PROTO_NC_EVENT_UPDATE_EVENT_REQ_DATA struct {
	EventData                        [1363]byte           // offset: 0, size: 1363 bytes
}

// PROTO_NC_GAMBLE_COIN_BUY_REQ (size = 8 bytes)
type PROTO_NC_GAMBLE_COIN_BUY_REQ_DATA struct {
	NCoin                            uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_GAMBLE_COIN_CHANGE_CMD (size = 8 bytes)
type PROTO_NC_GAMBLE_COIN_CHANGE_CMD_DATA struct {
	NCoin                            uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_BUY_ACK (size = 10 bytes)
type PROTO_NC_GAMBLE_COIN_DB_BUY_ACK_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_BUY_REQ (size = 32 bytes)
type PROTO_NC_GAMBLE_COIN_DB_BUY_REQ_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NCoin                            uint64               // offset: 8, size: 8 bytes
	NMoney                           uint64               // offset: 16, size: 8 bytes
	NCanExchangeMaxCoin              uint32               // offset: 24, size: 4 bytes
	TExchangedCoinDate               int32                // offset: 28, size: 4 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_SELL_ACK (size = 10 bytes)
type PROTO_NC_GAMBLE_COIN_DB_SELL_ACK_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_SELL_REQ (size = 24 bytes)
type PROTO_NC_GAMBLE_COIN_DB_SELL_REQ_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NCoin                            uint64               // offset: 8, size: 8 bytes
	NMoney                           uint64               // offset: 16, size: 8 bytes
}

// PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ (size = 1 bytes)
type PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ_DATA struct {
	NInvenSlot                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK (size = 4 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NBettingIndex                    uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK (size = 4 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NBettingCancelIndex              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD (size = 10 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
	BettingCancelInfo                uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ (size = 10 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ_DATA struct {
	BettingCancelInfo                uint64               // offset: 0, size: 8 bytes
	NBettingCancelIndex              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD (size = 10 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
	BettingInfo                      uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ (size = 10 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ_DATA struct {
	BettingInfo                      uint64               // offset: 0, size: 8 bytes
	NBettingIndex                    uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_BET_START_CMD (size = 7 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_BET_START_CMD_DATA struct {
	BModeType                        uint32               // offset: 0, size: 4 bytes
	CurrentRollingDice               [3]byte              // offset: 4, size: 3 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK (size = 14 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NGameNo                          int32                // offset: 2, size: 4 bytes
	Header                           NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	NResultCnt                       uint16               // offset: 12, size: 2 bytes
	DiceRollResult                   [0]byte              // offset: 14, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD (size = 12 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NGameNo                          int32                // offset: 6, size: 4 bytes
	NResultLogCnt                    uint16               // offset: 10, size: 2 bytes
	DiceRollResultLog                [0]byte              // offset: 12, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ (size = 11 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	DiceNum                          [3]byte              // offset: 6, size: 3 bytes
	NResultCnt                       uint16               // offset: 9, size: 2 bytes
	DiceRollResult                   [0]byte              // offset: 11, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_CMD (size = 2 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_CMD_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD (size = 5 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NDiceInfo                        [3]byte              // offset: 2, size: 3 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_BETTING_INFO_CMD (size = 3 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_BETTING_INFO_CMD_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	NBettingCnt                      uint16               // offset: 1, size: 2 bytes
	Betting                          [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK (size = 50 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	DiceHistory                      [21]byte             // offset: 2, size: 21 bytes
	NGameStatus                      uint32               // offset: 23, size: 4 bytes
	BModeType                        uint32               // offset: 27, size: 4 bytes
	NTimer                           uint16               // offset: 31, size: 2 bytes
	NInterval                        uint32               // offset: 33, size: 4 bytes
	NGetSysRate                      uint16               // offset: 37, size: 2 bytes
	NMaxBetMoney                     uint64               // offset: 39, size: 8 bytes
	NUseMinLv                        uint16               // offset: 47, size: 2 bytes
	NGamerCnt                        uint8                // offset: 49, size: 1 bytes
	NGamerHnd                        [0]byte              // offset: 50, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD (size = 2 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ (size = 4 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ_DATA struct {
	NNpcHandle                       uint16               // offset: 0, size: 2 bytes
	NWhereStand                      uint8                // offset: 2, size: 1 bytes
	NDirect                          uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK (size = 2 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD (size = 3 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD_DATA struct {
	BReserve                         uint8                // offset: 0, size: 1 bytes
	NPlayerHnd                       uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD (size = 1 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD_DATA struct {
	NGamerCnt                        uint8                // offset: 0, size: 1 bytes
	LargeAmount                      [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_TIMER_CMD (size = 7 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_TIMER_CMD_DATA struct {
	BIsStart                         uint8                // offset: 0, size: 1 bytes
	NTimer                           uint16               // offset: 1, size: 2 bytes
	NInterval                        uint32               // offset: 3, size: 4 bytes
}

// PROTO_NC_GAMBLE_ENTER_PLAYER_DIRECT_CMD (size = 3 bytes)
type PROTO_NC_GAMBLE_ENTER_PLAYER_DIRECT_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NDirect                          uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD (size = 8 bytes)
type PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD_DATA struct {
	NExchangeCoin                    uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK (size = 12 bytes)
type PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NExchangedCoin                   uint64               // offset: 4, size: 8 bytes
}

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_REQ (size = 8 bytes)
type PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DChackDate                       int32                // offset: 4, size: 4 bytes
}

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD (size = 4 bytes)
type PROTO_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD_DATA struct {
	DChackDate                       int32                // offset: 0, size: 4 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK (size = 8 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK_DATA struct {
	Machine                          uint8                // offset: 0, size: 1 bytes
	NError                           uint16               // offset: 1, size: 2 bytes
	JackPotPrice                     uint32               // offset: 3, size: 4 bytes
	RankNumber                       uint8                // offset: 7, size: 1 bytes
	RankList                         [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ (size = 1 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ_DATA struct {
	NSlotMachineNum                  uint8                // offset: 0, size: 1 bytes
	MachineIndex                     [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK (size = 10 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_REQ (size = 36 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	PlayerRegistNumber               uint32               // offset: 8, size: 4 bytes
	CardDisplay                      [10]byte             // offset: 12, size: 10 bytes
	IsJackPot                        uint8                // offset: 22, size: 1 bytes
	JackPotAccumulate                uint32               // offset: 23, size: 4 bytes
	BetUnit                          uint16               // offset: 27, size: 2 bytes
	BetLine                          uint8                // offset: 29, size: 1 bytes
	NJackpotRate                     uint16               // offset: 30, size: 2 bytes
	WinMoney                         uint32               // offset: 32, size: 4 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK (size = 4 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK_DATA struct {
	BettingUnit                      uint16               // offset: 0, size: 2 bytes
	NError                           uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ (size = 4 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ_DATA struct {
	NNpcHandle                       uint16               // offset: 0, size: 2 bytes
	NWhereStand                      uint8                // offset: 2, size: 1 bytes
	NDirect                          uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK (size = 2 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD (size = 1 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD_DATA struct {
	JackPotNum                       uint8                // offset: 0, size: 1 bytes
	JackPotList                      [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD (size = 28 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Price                            uint64               // offset: 20, size: 8 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK (size = 2 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ (size = 4 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ_DATA struct {
	Betting                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK (size = 11 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Screen                           [9]byte              // offset: 2, size: 9 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK (size = 10 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Bet                              uint16               // offset: 2, size: 2 bytes
	Win                              uint32               // offset: 4, size: 4 bytes
	WinLine                          uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK (size = 9 bytes)
type PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NNpcHandle                       uint16               // offset: 2, size: 2 bytes
	NGambleType                      uint32               // offset: 4, size: 4 bytes
	NWhereStand                      uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ (size = 2 bytes)
type PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ_DATA struct {
	NNpcHandle                       uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GAMBLE_WORLD_PREVMAPNAME_CMD (size = 18 bytes)
type PROTO_NC_GAMBLE_WORLD_PREVMAPNAME_CMD_DATA struct {
	NWldHandle                       uint16               // offset: 0, size: 2 bytes
	SPrevMapName                     Name3                // offset: 2, size: 12 bytes
	NCharNo                          uint32               // offset: 14, size: 4 bytes
}

// PROTO_NC_GAMBLE_ZONE_PREVMAPNAME_CMD (size = 12 bytes)
type PROTO_NC_GAMBLE_ZONE_PREVMAPNAME_CMD_DATA struct {
	SPrevMapName                     Name3                // offset: 0, size: 12 bytes
}

// PROTO_NC_LOG_GAME_LINK (size = 41 bytes)
type PROTO_NC_LOG_GAME_LINK_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           [20]byte             // offset: 4, size: 20 bytes
	StoreType                        uint8                // offset: 24, size: 1 bytes
	CharExp                          uint64               // offset: 25, size: 8 bytes
	CharMoney                        uint64               // offset: 33, size: 8 bytes
}

// PROTO_NC_MAP_LINKEND_CMD (size = 47 bytes)
type PROTO_NC_MAP_LINKEND_CMD_DATA struct {
	Wmshandle                        uint16               // offset: 0, size: 2 bytes
	Zonehandle                       uint16               // offset: 2, size: 2 bytes
	Chrregnum                        uint32               // offset: 4, size: 4 bytes
	MapField                         Name3                // offset: 8, size: 12 bytes
	Zoneid                           uint8                // offset: 20, size: 1 bytes
	Location                         uint64               // offset: 21, size: 8 bytes
	StoreType                        uint8                // offset: 29, size: 1 bytes
	CharExp                          uint64               // offset: 30, size: 8 bytes
	CharMoney                        uint64               // offset: 38, size: 8 bytes
	IsPartyRelease                   uint8                // offset: 46, size: 1 bytes
}

// PROTO_NC_MAP_LINKRESERVE_ACK (size = 36 bytes)
type PROTO_NC_MAP_LINKRESERVE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Linkother                        [30]byte             // offset: 6, size: 30 bytes
}

// PROTO_NC_MAP_LINKRESERVE_ACK_SEND (size = 39 bytes)
type PROTO_NC_MAP_LINKRESERVE_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [36]byte             // offset: 3, size: 36 bytes
}

// PROTO_NC_MAP_LINKRESERVE_REQ (size = 59 bytes)
type PROTO_NC_MAP_LINKRESERVE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Wldmanhandle                     uint16               // offset: 6, size: 2 bytes
	Newmechineid                     uint8                // offset: 8, size: 1 bytes
	Linkto                           [20]byte             // offset: 9, size: 20 bytes
	Linkother                        [30]byte             // offset: 29, size: 30 bytes
}

// PROTO_NC_MAP_LINKSTART_CMD (size = 6 bytes)
type PROTO_NC_MAP_LINKSTART_CMD_DATA struct {
	Wmshandle                        uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_MAP_REGIST_CMD (size = 2 bytes)
type PROTO_NC_MAP_REGIST_CMD_DATA struct {
	Numofmap                         uint16               // offset: 0, size: 2 bytes
	Maparray                         [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS (size = 21 bytes)
type PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS_DATA struct {
	WorldNo                          uint8                // offset: 0, size: 1 bytes
	UserNo                           uint32               // offset: 1, size: 4 bytes
	CharNo                           uint32               // offset: 5, size: 4 bytes
	CharClass                        uint8                // offset: 9, size: 1 bytes
	CharLevel                        uint8                // offset: 10, size: 1 bytes
	CharGetExp                       uint64               // offset: 11, size: 8 bytes
	CharPlayMin                      uint16               // offset: 19, size: 2 bytes
}

// PROTO_NC_PROMOTION_DB_REWARD_ACK (size = 14 bytes)
type PROTO_NC_PROMOTION_DB_REWARD_ACK_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NWorldNo                         uint8                // offset: 4, size: 1 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
	NClass                           uint8                // offset: 9, size: 1 bytes
	NLevel                           uint8                // offset: 10, size: 1 bytes
	DBError                          uint16               // offset: 11, size: 2 bytes
	NNumOfItems                      uint8                // offset: 13, size: 1 bytes
	RewardList                       [0]byte              // offset: 14, flexible array member
}

// PROTO_NC_REPORT_ADD_ACK (size = 26 bytes)
type PROTO_NC_REPORT_ADD_ACK_DATA struct {
	SHarmer                          Name5                // offset: 0, size: 20 bytes
	NRegNumber                       uint32               // offset: 20, size: 4 bytes
	Err                              uint16               // offset: 24, size: 2 bytes
}

// PROTO_NC_REPORT_ADD_REQ (size = 2342 bytes)
type PROTO_NC_REPORT_ADD_REQ_DATA struct {
	SHarmer                          Name5                // offset: 0, size: 20 bytes
	ByReportType                     [16]byte             // offset: 20, size: 16 bytes
	ByReason                         [256]byte            // offset: 36, size: 256 bytes
	NChatLogSize                     uint16               // offset: 292, size: 2 bytes
	ByChatLog                        [2048]byte           // offset: 294, size: 2048 bytes
}

// PROTO_NC_REPORT_CANCEL_ACK (size = 6 bytes)
type PROTO_NC_REPORT_CANCEL_ACK_DATA struct {
	NRegNumber                       uint32               // offset: 0, size: 4 bytes
	Err                              uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_REPORT_CANCEL_REQ (size = 4 bytes)
type PROTO_NC_REPORT_CANCEL_REQ_DATA struct {
	NRegNumber                       uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_REPORT_GET_ACK (size = 135 bytes)
type PROTO_NC_REPORT_GET_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NRegNumber                       uint32               // offset: 2, size: 4 bytes
	ByReportResult                   uint8                // offset: 6, size: 1 bytes
	ByRemark                         [128]byte            // offset: 7, size: 128 bytes
}

// PROTO_NC_REPORT_GET_REQ (size = 1 bytes)
type PROTO_NC_REPORT_GET_REQ_DATA struct {}

// PROTO_NC_SYSLOG_ACCOUNT_LOGOUT (size = 8 bytes)
type PROTO_NC_SYSLOG_ACCOUNT_LOGOUT_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	LoginKeepTime                    uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_SYSLOG_CHAR_CREATED (size = 9 bytes)
type PROTO_NC_SYSLOG_CHAR_CREATED_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	CharClass                        uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_SYSLOG_CHAR_DELETED (size = 9 bytes)
type PROTO_NC_SYSLOG_CHAR_DELETED_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	CharClass                        uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_SYSLOG_CHAR_ENTER_GAME (size = 8 bytes)
type PROTO_NC_SYSLOG_CHAR_ENTER_GAME_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_SYSLOG_CHAR_LEAVE_GAME (size = 12 bytes)
type PROTO_NC_SYSLOG_CHAR_LEAVE_GAME_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	PlayTime                         uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_SYSLOG_CHAR_LEVEL_UP (size = 9 bytes)
type PROTO_NC_SYSLOG_CHAR_LEVEL_UP_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	CharLev                          uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_SYSLOG_SERVER_CCU (size = 5 bytes)
type PROTO_NC_SYSLOG_SERVER_CCU_DATA struct {
	UserCount                        uint32               // offset: 0, size: 4 bytes
	WorldNo                          uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACADEMY_MASTER_CMD (size = 20 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACADEMY_MASTER_CMD_DATA struct {
	SMasterID                        Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK (size = 23 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NSlotNo                          uint8                // offset: 2, size: 1 bytes
	SNewID                           Name5                // offset: 3, size: 20 bytes
}

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_CMD (size = 41 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_CMD_DATA struct {
	BInit                            uint8                // offset: 0, size: 1 bytes
	SOldID                           Name5                // offset: 1, size: 20 bytes
	SNewID                           Name5                // offset: 21, size: 20 bytes
}

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ (size = 21 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ_DATA struct {
	NSlotNo                          uint8                // offset: 0, size: 1 bytes
	SNewID                           Name5                // offset: 1, size: 20 bytes
}

// PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ACK (size = 32 bytes)
type PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NError                           uint16               // offset: 2, size: 2 bytes
	NUserNo                          uint32               // offset: 4, size: 4 bytes
	NCharNo                          uint32               // offset: 8, size: 4 bytes
	SNewID                           Name5                // offset: 12, size: 20 bytes
}

// PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_REQ (size = 35 bytes)
type PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NRowNo                           uint32               // offset: 2, size: 4 bytes
	NUserNo                          uint32               // offset: 6, size: 4 bytes
	NCharNo                          uint32               // offset: 10, size: 4 bytes
	SNewID                           Name5                // offset: 14, size: 20 bytes
	BInit                            uint8                // offset: 34, size: 1 bytes
}

// PROTO_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD (size = 45 bytes)
type PROTO_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	BInit                            uint8                // offset: 4, size: 1 bytes
	SOldID                           Name5                // offset: 5, size: 20 bytes
	SNewID                           Name5                // offset: 25, size: 20 bytes
}
