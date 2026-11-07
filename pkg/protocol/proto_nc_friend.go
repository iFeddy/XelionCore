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

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_CHAR_FRIEND_CMD (size = 2 bytes)
type PROTO_NC_CHAR_FRIEND_CMD_DATA struct {
	Numberoffriend                   uint16               // offset: 0, size: 2 bytes
	Friend_array                     [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_FRIEND_CMD_INFO (size = 25 bytes)
type PROTO_NC_CHAR_FRIEND_CMD_INFO struct {
	Charno                           uint32               // offset: 0, size: 4 bytes
	Charid                           Name5                // offset: 4, size: 20 bytes
	Flag                             uint8                // offset: 24, size: 1 bytes
}

// PROTO_NC_CT_ADD_FRIEND_CMD (size = 8 bytes)
type PROTO_NC_CT_ADD_FRIEND_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NumberOfFriends                  uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_FRIEND_ADD_CMD (size = 72 bytes)
type PROTO_NC_FRIEND_ADD_CMD_DATA struct {
	Friendinfo                       [72]byte             // offset: 0, size: 72 bytes
}

// PROTO_NC_FRIEND_CLASS_CHANGE_CMD (size = 21 bytes)
type PROTO_NC_FRIEND_CLASS_CHANGE_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	New_class                        uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_FRIEND_DB_DEL_ACK (size = 48 bytes)
type PROTO_NC_FRIEND_DB_DEL_ACK_DATA struct {
	Handle_charno                    uint16               // offset: 0, size: 2 bytes
	Charno                           uint32               // offset: 2, size: 4 bytes
	Charid                           Name5                // offset: 6, size: 20 bytes
	Friendid                         Name5                // offset: 26, size: 20 bytes
	Err                              uint16               // offset: 46, size: 2 bytes
}

// PROTO_NC_FRIEND_DB_DEL_REQ (size = 46 bytes)
type PROTO_NC_FRIEND_DB_DEL_REQ_DATA struct {
	Handle_charno                    uint16               // offset: 0, size: 2 bytes
	Charno                           uint32               // offset: 2, size: 4 bytes
	Charid                           Name5                // offset: 6, size: 20 bytes
	Friendid                         Name5                // offset: 26, size: 20 bytes
}

// PROTO_NC_FRIEND_DB_GET_ACK (size = 84 bytes)
type PROTO_NC_FRIEND_DB_GET_ACK_DATA struct {
	Handle_char                      uint16               // offset: 0, size: 2 bytes
	Charno                           uint32               // offset: 2, size: 4 bytes
	Friendno                         uint32               // offset: 6, size: 4 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
	Friendinfo                       [72]byte             // offset: 12, size: 72 bytes
}

// PROTO_NC_FRIEND_DB_GET_REQ (size = 10 bytes)
type PROTO_NC_FRIEND_DB_GET_REQ_DATA struct {
	Handle_char                      uint16               // offset: 0, size: 2 bytes
	Charno                           uint32               // offset: 2, size: 4 bytes
	Friendno                         uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_FRIEND_DB_POINT (size = 8 bytes)
type PROTO_NC_FRIEND_DB_POINT struct {
	Charregnum                       uint32               // offset: 0, size: 4 bytes
	Friendpoint                      uint16               // offset: 4, size: 2 bytes
	Diffpoint                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_FRIEND_DB_POINT_CMD (size = 2 bytes)
type PROTO_NC_FRIEND_DB_POINT_CMD_DATA struct {
	Numoffriend                      uint16               // offset: 0, size: 2 bytes
	Pointinfo                        [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_FRIEND_DB_SET_ACK (size = 14 bytes)
type PROTO_NC_FRIEND_DB_SET_ACK_DATA struct {
	Handle_charno                    uint16               // offset: 0, size: 2 bytes
	Handle_friendno                  uint16               // offset: 2, size: 2 bytes
	Charno                           uint32               // offset: 4, size: 4 bytes
	Friendno                         uint32               // offset: 8, size: 4 bytes
	Err                              uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_FRIEND_DB_SET_REQ (size = 13 bytes)
type PROTO_NC_FRIEND_DB_SET_REQ_DATA struct {
	Handle_charno                    uint16               // offset: 0, size: 2 bytes
	Handle_friendno                  uint16               // offset: 2, size: 2 bytes
	Charno                           uint32               // offset: 4, size: 4 bytes
	Friendno                         uint32               // offset: 8, size: 4 bytes
	Flag                             uint8                // offset: 12, size: 1 bytes
}

// PROTO_NC_FRIEND_DEL_ACK (size = 42 bytes)
type PROTO_NC_FRIEND_DEL_ACK_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
	Err                              uint16               // offset: 40, size: 2 bytes
}

// PROTO_NC_FRIEND_DEL_CMD (size = 20 bytes)
type PROTO_NC_FRIEND_DEL_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_FRIEND_DEL_REQ (size = 40 bytes)
type PROTO_NC_FRIEND_DEL_REQ_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_FRIEND_FIND_FRIENDS_ACK (size = 4 bytes)
type PROTO_NC_FRIEND_FIND_FRIENDS_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Numoffriend                      uint16               // offset: 2, size: 2 bytes
	Friendinfo                       [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_FRIEND_GET_DIFF_FRIEND_POINT_CMD (size = 2 bytes)
type PROTO_NC_FRIEND_GET_DIFF_FRIEND_POINT_CMD_DATA struct {
	Diffpoint                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_FRIEND_LEVEL_CMD (size = 21 bytes)
type PROTO_NC_FRIEND_LEVEL_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Level                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_FRIEND_LIST_CMD (size = 1 bytes)
type PROTO_NC_FRIEND_LIST_CMD_DATA struct {
	Numoffriend                      uint8                // offset: 0, size: 1 bytes
	Friendinfo                       [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_FRIEND_LOGIN_CMD (size = 32 bytes)
type PROTO_NC_FRIEND_LOGIN_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	MapField                         Name3                // offset: 20, size: 12 bytes
}

// PROTO_NC_FRIEND_LOGOUT_CMD (size = 20 bytes)
type PROTO_NC_FRIEND_LOGOUT_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_FRIEND_MAP_CMD (size = 32 bytes)
type PROTO_NC_FRIEND_MAP_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	MapField                         Name3                // offset: 20, size: 12 bytes
}

// PROTO_NC_FRIEND_POINT_ACK (size = 2 bytes)
type PROTO_NC_FRIEND_POINT_ACK_DATA struct {
	Friendpoint                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_FRIEND_REFUSE_CMD (size = 20 bytes)
type PROTO_NC_FRIEND_REFUSE_CMD_DATA struct {
	Friendid                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_FRIEND_SET_ACK (size = 42 bytes)
type PROTO_NC_FRIEND_SET_ACK_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
	Err                              uint16               // offset: 40, size: 2 bytes
}

// PROTO_NC_FRIEND_SET_CONFIRM_ACK (size = 41 bytes)
type PROTO_NC_FRIEND_SET_CONFIRM_ACK_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
	Accept_friend                    uint8                // offset: 40, size: 1 bytes
}

// PROTO_NC_FRIEND_SET_CONFIRM_REQ (size = 40 bytes)
type PROTO_NC_FRIEND_SET_CONFIRM_REQ_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_FRIEND_SET_REQ (size = 40 bytes)
type PROTO_NC_FRIEND_SET_REQ_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Friendid                         Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_WORLD_CMD (size = 20 bytes)
type PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_WORLD_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_ZONE_CMD (size = 20 bytes)
type PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_ZONE_CMD_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_ITEMDB_UES_FRIEND_POINT_ACK (size = 15 bytes)
type PROTO_NC_ITEMDB_UES_FRIEND_POINT_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	ItemID                           uint16               // offset: 8, size: 2 bytes
	Friendpoint                      uint16               // offset: 10, size: 2 bytes
	Lockindex                        uint16               // offset: 12, size: 2 bytes
	Numofitem                        uint8                // offset: 14, size: 1 bytes
}

// PROTO_NC_ITEMDB_UES_FRIEND_POINT_REQ (size = 123 bytes)
type PROTO_NC_ITEMDB_UES_FRIEND_POINT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemInven                        uint16               // offset: 6, size: 2 bytes
	FriendUsePoint                   uint16               // offset: 8, size: 2 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	Iti                              [111]byte            // offset: 12, size: 111 bytes
}
