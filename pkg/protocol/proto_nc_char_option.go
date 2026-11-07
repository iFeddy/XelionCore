package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHAR_OPTION_DB_GET_GAME_ACK (size = 67 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_GAME_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [64]byte             // offset: 3, size: 64 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_GAME_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_GAME_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_KEYMAPPING_ACK (size = 311 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_KEYMAPPING_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [308]byte            // offset: 3, size: 308 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_KEYMAPPING_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_KEYMAPPING_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTDATA_ACK (size = 1027 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTDATA_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [1024]byte           // offset: 3, size: 1024 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTDATA_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTDATA_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTSIZE_ACK (size = 27 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTSIZE_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [24]byte             // offset: 3, size: 24 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTSIZE_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SHORTCUTSIZE_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SOUND_ACK (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SOUND_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_SOUND_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_SOUND_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_VIDEO_ACK (size = 63 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_VIDEO_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [60]byte             // offset: 3, size: 60 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_VIDEO_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_VIDEO_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_WINDOWPOS_ACK (size = 395 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_WINDOWPOS_ACK_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	BSuccess                         uint8                // offset: 2, size: 1 bytes
	Data                             [392]byte            // offset: 3, size: 392 bytes
}

// PROTO_NC_CHAR_OPTION_DB_GET_WINDOWPOS_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_DB_GET_WINDOWPOS_REQ_DATA struct {
	WorldManagerHandle               uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_GAME_CMD (size = 68 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_GAME_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [64]byte             // offset: 4, size: 64 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_KEYMAPPING_CMD (size = 312 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_KEYMAPPING_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [308]byte            // offset: 4, size: 308 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_SHORTCUTDATA_CMD (size = 1028 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_SHORTCUTDATA_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [1024]byte           // offset: 4, size: 1024 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_SHORTCUTSIZE_CMD (size = 28 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_SHORTCUTSIZE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [24]byte             // offset: 4, size: 24 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_SOUND_CMD (size = 5 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_SOUND_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_VIDEO_CMD (size = 64 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_VIDEO_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [60]byte             // offset: 4, size: 60 bytes
}

// PROTO_NC_CHAR_OPTION_DB_SET_WINDOWPOS_CMD (size = 396 bytes)
type PROTO_NC_CHAR_OPTION_DB_SET_WINDOWPOS_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	Data                             [392]byte            // offset: 4, size: 392 bytes
}

// PROTO_NC_CHAR_OPTION_GAME (size = 64 bytes)
type PROTO_NC_CHAR_OPTION_GAME struct {
	Data                             [64]byte             // offset: 0, size: 64 bytes
}

// PROTO_NC_CHAR_OPTION_GET_GAME_ACK (size = 65 bytes)
type PROTO_NC_CHAR_OPTION_GET_GAME_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [64]byte             // offset: 1, size: 64 bytes
}

// PROTO_NC_CHAR_OPTION_GET_GAME_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_GAME_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK (size = 309 bytes)
type PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [308]byte            // offset: 1, size: 308 bytes
}

// PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK (size = 1025 bytes)
type PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [1024]byte           // offset: 1, size: 1024 bytes
}

// PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK (size = 25 bytes)
type PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [24]byte             // offset: 1, size: 24 bytes
}

// PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_SOUND_ACK (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_GET_SOUND_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_GET_SOUND_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_SOUND_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_VIDEO_ACK (size = 61 bytes)
type PROTO_NC_CHAR_OPTION_GET_VIDEO_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [60]byte             // offset: 1, size: 60 bytes
}

// PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK (size = 393 bytes)
type PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK_DATA struct {
	BSuccess                         uint8                // offset: 0, size: 1 bytes
	Data                             [392]byte            // offset: 1, size: 392 bytes
}

// PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_DATA_TYPE_CMD (size = 8 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_DATA_TYPE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DataType                         uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_GAMEOPTION_CMD (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_GAMEOPTION_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBGameOptionData                 uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_KEYMAP_CMD (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_KEYMAP_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBKeyMapData                     uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_SHORTCUTDATA_CMD (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_GET_SHORTCUTDATA_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBShortCutData                   uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_GAMEOPTION_ACK (size = 8 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_GAMEOPTION_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	DBGameOptionData                 uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_GAMEOPTION_REQ (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_GAMEOPTION_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_KEYMAP_ACK (size = 8 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_KEYMAP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	DBKeyMapData                     uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_KEYMAP_REQ (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_INIT_KEYMAP_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_GAMEOPTION_ACK (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_GAMEOPTION_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_GAMEOPTION_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_GAMEOPTION_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBGameOptionData                 uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_KEYMAP_ACK (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_KEYMAP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_KEYMAP_REQ (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_KEYMAP_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBKeyMapData                     uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_SHORTCUTDATA_ACK (size = 6 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_SHORTCUTDATA_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_SHORTCUTDATA_REQ (size = 5 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_DB_SET_SHORTCUTDATA_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBShortCutData                   uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD_DATA struct {
	NGameOptionDataCnt               uint16               // offset: 0, size: 2 bytes
	GameOptionData                   [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD_DATA struct {
	NKeyMapDataCnt                   uint16               // offset: 0, size: 2 bytes
	KeyMapData                       [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD_DATA struct {
	NShortCutDataCnt                 uint16               // offset: 0, size: 2 bytes
	ShortCutData                     [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_ACK (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	DBGameOptionData                 uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_ACK (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	DBKeyMapData                     uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ_DATA struct {}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_ACK (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ_DATA struct {
	NGameOptionDataCnt               uint16               // offset: 0, size: 2 bytes
	GameOptionData                   [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_ACK (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ_DATA struct {
	NKeyMapDataCnt                   uint16               // offset: 0, size: 2 bytes
	KeyMapData                       [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK (size = 2 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ_DATA struct {
	NShortCutDataCnt                 uint8                // offset: 0, size: 1 bytes
	ShortCutData                     [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_CHAR_OPTION_KEYMAPPING (size = 308 bytes)
type PROTO_NC_CHAR_OPTION_KEYMAPPING struct {
	Data                             [308]byte            // offset: 0, size: 308 bytes
}

// PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR_REQ (size = 4 bytes)
type PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_CHAR_OPTION_SET_GAME_CMD (size = 64 bytes)
type PROTO_NC_CHAR_OPTION_SET_GAME_CMD_DATA struct {
	Data                             [64]byte             // offset: 0, size: 64 bytes
}

// PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD (size = 308 bytes)
type PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD_DATA struct {
	Data                             [308]byte            // offset: 0, size: 308 bytes
}

// PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD (size = 1024 bytes)
type PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD_DATA struct {
	Data                             [1024]byte           // offset: 0, size: 1024 bytes
}

// PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD (size = 24 bytes)
type PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD_DATA struct {
	Data                             [24]byte             // offset: 0, size: 24 bytes
}

// PROTO_NC_CHAR_OPTION_SET_SOUND_CMD (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_SET_SOUND_CMD_DATA struct {
	Data                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_SET_VIDEO_CMD (size = 60 bytes)
type PROTO_NC_CHAR_OPTION_SET_VIDEO_CMD_DATA struct {
	Data                             [60]byte             // offset: 0, size: 60 bytes
}

// PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD (size = 392 bytes)
type PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD_DATA struct {
	Data                             [392]byte            // offset: 0, size: 392 bytes
}

// PROTO_NC_CHAR_OPTION_SHORTCUTDATA (size = 1024 bytes)
type PROTO_NC_CHAR_OPTION_SHORTCUTDATA struct {
	Data                             [1024]byte           // offset: 0, size: 1024 bytes
}

// PROTO_NC_CHAR_OPTION_SHORTCUTSIZE (size = 24 bytes)
type PROTO_NC_CHAR_OPTION_SHORTCUTSIZE struct {
	Data                             [24]byte             // offset: 0, size: 24 bytes
}

// PROTO_NC_CHAR_OPTION_SOUND (size = 1 bytes)
type PROTO_NC_CHAR_OPTION_SOUND struct {
	Data                             [1]byte              // offset: 0, size: 1 bytes
}

// PROTO_NC_CHAR_OPTION_VIDEO (size = 60 bytes)
type PROTO_NC_CHAR_OPTION_VIDEO struct {
	Data                             [60]byte             // offset: 0, size: 60 bytes
}

// PROTO_NC_CHAR_OPTION_WINDOWPOS (size = 392 bytes)
type PROTO_NC_CHAR_OPTION_WINDOWPOS struct {
	Data                             [392]byte            // offset: 0, size: 392 bytes
}
