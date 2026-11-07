package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_AVATAR_CREATEDATAFAIL_ACK (size = 22 bytes)
type PROTO_NC_AVATAR_CREATEDATAFAIL_ACK_DATA struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_AVATAR_CREATEDATASUC_ACK (size = 1 bytes)
type PROTO_NC_AVATAR_CREATEDATASUC_ACK_DATA struct {
	Numofavatar                      uint8                // offset: 0, size: 1 bytes
	Avatar                           [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_AVATAR_CREATEDATA_REQ (size = 33 bytes)
type PROTO_NC_AVATAR_CREATEDATA_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Usrregnum                        uint32               // offset: 2, size: 4 bytes
	Createworld                      uint8                // offset: 6, size: 1 bytes
	Adminlevel                       uint8                // offset: 7, size: 1 bytes
	Clientdata                       [25]byte             // offset: 8, size: 25 bytes
}

// PROTO_NC_AVATAR_CREATEDATA_REQ_SEND (size = 36 bytes)
type PROTO_NC_AVATAR_CREATEDATA_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [33]byte             // offset: 3, size: 33 bytes
}

// PROTO_NC_AVATAR_CREATEFAIL_ACK (size = 2 bytes)
type PROTO_NC_AVATAR_CREATEFAIL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_AVATAR_CREATEFAIL_ACK_SEND (size = 5 bytes)
type PROTO_NC_AVATAR_CREATEFAIL_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_AVATAR_CREATESUCC_ACK (size = 131 bytes)
type PROTO_NC_AVATAR_CREATESUCC_ACK_DATA struct {
	Numofavatar                      uint8                // offset: 0, size: 1 bytes
	Avatar                           [130]byte            // offset: 1, size: 130 bytes
}

// PROTO_NC_AVATAR_CREATESUCC_ACK_SEND (size = 134 bytes)
type PROTO_NC_AVATAR_CREATESUCC_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              [131]byte            // offset: 3, size: 131 bytes
}

// PROTO_NC_AVATAR_CREATE_REQ (size = 25 bytes)
type PROTO_NC_AVATAR_CREATE_REQ_DATA struct {
	Slotnum                          uint8                // offset: 0, size: 1 bytes
	Name                             Name5                // offset: 1, size: 20 bytes
	Char_shape                       uint32               // offset: 21, size: 4 bytes
}

// PROTO_NC_AVATAR_ERASEDATAFAIL_ACK (size = 4 bytes)
type PROTO_NC_AVATAR_ERASEDATAFAIL_ACK_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_AVATAR_ERASEDATASUC_ACK (size = 15 bytes)
type PROTO_NC_AVATAR_ERASEDATASUC_ACK_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Usrregnum                        uint32               // offset: 2, size: 4 bytes
	Avatar                           uint32               // offset: 6, size: 4 bytes
	Slot                             uint8                // offset: 10, size: 1 bytes
	NNo                              uint32               // offset: 11, size: 4 bytes
}

// PROTO_NC_AVATAR_ERASEDATA_REQ (size = 15 bytes)
type PROTO_NC_AVATAR_ERASEDATA_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Usrregnum                        uint32               // offset: 2, size: 4 bytes
	Avatar                           uint32               // offset: 6, size: 4 bytes
	Slot                             uint8                // offset: 10, size: 1 bytes
	NNo                              uint32               // offset: 11, size: 4 bytes
}

// PROTO_NC_AVATAR_ERASEDATA_REQ_SEND (size = 18 bytes)
type PROTO_NC_AVATAR_ERASEDATA_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [15]byte             // offset: 3, size: 15 bytes
}

// PROTO_NC_AVATAR_ERASEFAIL_ACK (size = 2 bytes)
type PROTO_NC_AVATAR_ERASEFAIL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_AVATAR_ERASEFAIL_ACK_SEND (size = 5 bytes)
type PROTO_NC_AVATAR_ERASEFAIL_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_AVATAR_ERASESUC_ACK (size = 1 bytes)
type PROTO_NC_AVATAR_ERASESUC_ACK struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_AVATAR_ERASESUC_ACK_SEND (size = 4 bytes)
type PROTO_NC_AVATAR_ERASESUC_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_AVATAR_ERASE_REQ (size = 1 bytes)
type PROTO_NC_AVATAR_ERASE_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_AVATAR_RENAME_ACK (size = 23 bytes)
type PROTO_NC_AVATAR_RENAME_ACK_DATA struct {
	SlotNo                           uint8                // offset: 0, size: 1 bytes
	NewName                          Name5                // offset: 1, size: 20 bytes
	Error                            uint16               // offset: 21, size: 2 bytes
}

// PROTO_NC_AVATAR_RENAME_DB_ACK (size = 29 bytes)
type PROTO_NC_AVATAR_RENAME_DB_ACK_DATA struct {
	NetPacketHeader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	SlotNo                           uint8                // offset: 6, size: 1 bytes
	NewName                          Name5                // offset: 7, size: 20 bytes
	Error                            uint16               // offset: 27, size: 2 bytes
}

// PROTO_NC_AVATAR_RENAME_DB_REQ (size = 27 bytes)
type PROTO_NC_AVATAR_RENAME_DB_REQ_DATA struct {
	NetPacketHeader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	SlotNo                           uint8                // offset: 6, size: 1 bytes
	NewName                          Name5                // offset: 7, size: 20 bytes
}

// PROTO_NC_AVATAR_RENAME_REQ (size = 21 bytes)
type PROTO_NC_AVATAR_RENAME_REQ_DATA struct {
	SlotNo                           uint8                // offset: 0, size: 1 bytes
	NewName                          Name5                // offset: 1, size: 20 bytes
}

// PROTO_NC_LOG_GAME_CREATE_AVATAR (size = 4 bytes)
type PROTO_NC_LOG_GAME_CREATE_AVATAR_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_LOG_GAME_CREATE_AVATAR_SEND (size = 7 bytes)
type PROTO_NC_LOG_GAME_CREATE_AVATAR_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Log                              uint32               // offset: 3, size: 4 bytes
}

// PROTO_NC_LOG_GAME_DELETE_AVATAR (size = 4 bytes)
type PROTO_NC_LOG_GAME_DELETE_AVATAR_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_LOG_GAME_DELETE_AVATAR_SEND (size = 7 bytes)
type PROTO_NC_LOG_GAME_DELETE_AVATAR_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Log                              uint32               // offset: 3, size: 4 bytes
}

// PROTO_NC_USER_AVATARINFO_ACK (size = 1 bytes)
type PROTO_NC_USER_AVATARINFO_ACK_DATA struct {
	Numofavatar                      uint8                // offset: 0, size: 1 bytes
	Avatar                           [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_USER_AVATARINFO_REQ (size = 6 bytes)
type PROTO_NC_USER_AVATARINFO_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Usrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_USER_AVATARINFO_REQ_SEND (size = 9 bytes)
type PROTO_NC_USER_AVATARINFO_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [6]byte              // offset: 3, size: 6 bytes
}

// PROTO_NC_USER_AVATAR_LIST_REQ (size = 1 bytes)
type PROTO_NC_USER_AVATAR_LIST_REQ_DATA struct {}
