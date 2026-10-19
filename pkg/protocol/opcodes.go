package protocol

// MakeOpCode combines a 6-bit Department (Header) and a 10-bit Command (Type) into a 16-bit OpCode.
// Formula oficial: (Department << 10) | (Command & 0x3FF)
func MakeOpCode(dept uint8, cmd uint16) uint16 {
	return (uint16(dept) << 10) | (cmd & 0x03FF)
}

// ExtractHeader extracts the 6-bit department from a 16-bit OpCode.
func ExtractHeader(opCode uint16) uint8 {
	return uint8(opCode >> 10)
}

// ExtractType extracts the 10-bit command from a 16-bit OpCode.
func ExtractType(opCode uint16) uint16 {
	return opCode & 0x03FF
}

// Protocol OpCodes (100% oficial de Fiesta Online PDBs: Zone.pdb, WorldManager.pdb, Login.pdb, Character.pdb)
const (
	// Department 2: MISC / System / Handshake
	PROTO_NC_MISC_HEARTBEAT_REQ  uint16 = (2 << 10) | 4  // 0x0804
	PROTO_NC_MISC_HEARTBEAT_ACK  uint16 = (2 << 10) | 5  // 0x0805
	PROTO_NC_MISC_SEED_CMD       uint16 = (2 << 10) | 7  // 0x0807 (Handshake seed offset)
	PROTO_NC_MISC_GAMETIME_REQ   uint16 = (2 << 10) | 13 // 0x080D
	PROTO_NC_MISC_GAMETIME_ACK   uint16 = (2 << 10) | 14 // 0x080E
	PROTO_NC_MISC_TIME_SYNC_ACK  uint16 = (2 << 10) | 69 // 0x0845
	PROTO_NC_MISC_UPDATE_TIME    uint16 = (2 << 10) | 73 // 0x0849

	// Department 3: USER / Login & World Selection
	PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK uint16 = (3 << 10) | 2   // 0x0C02
	PROTO_NC_USER_FILE_CHECK_REQ                uint16 = (3 << 10) | 4   // 0x0C04
	PROTO_NC_USER_FILE_CHECK_ACK                uint16 = (3 << 10) | 5   // 0x0C05
	PROTO_NC_USER_LOGIN_REQ                     uint16 = (3 << 10) | 6   // 0x0C06
	PROTO_NC_USER_LOGINFAIL_ACK                 uint16 = (3 << 10) | 9   // 0x0C09
	PROTO_NC_USER_LOGIN_ACK                     uint16 = (3 << 10) | 10  // 0x0C0A
	PROTO_NC_USER_WORLDSELECT_REQ               uint16 = (3 << 10) | 11  // 0x0C0B
	PROTO_NC_USER_WORLDSELECT_ACK               uint16 = (3 << 10) | 12  // 0x0C0C
	PROTO_NC_USER_LOGINWORLD_REQ                uint16 = (3 << 10) | 15  // 0x0C0F
	PROTO_NC_MISC_WORLD_CLIENT_KEY_REQ          uint16 = (3 << 10) | 15  // 0x0C0F
	PROTO_NC_USER_LOGINWORLD_ACK                uint16 = (3 << 10) | 16  // 0x0C10
	PROTO_NC_MISC_WORLD_CLIENT_KEY_ACK          uint16 = (3 << 10) | 16  // 0x0C10
	PROTO_NC_CHAR_CLIENT_CHARLIST_ACK           uint16 = (3 << 10) | 20  // 0x0C14
	PROTO_NC_USER_NORMALLOGOUT_CMD              uint16 = (3 << 10) | 24  // 0x0C18
	PROTO_NC_USER_WORLDSTATUS_REQ               uint16 = (3 << 10) | 27  // 0x0C1B
	PROTO_NC_USER_WORLDSTATUS_ACK               uint16 = (3 << 10) | 28  // 0x0C1C
	PROTO_NC_USER_BACKTOWORLDLIST_REQ           uint16 = (3 << 10) | 51  // 0x0C33
	PROTO_NC_USER_BACKTOWORLDLIST_ACK           uint16 = (3 << 10) | 52  // 0x0C34
	PROTO_NC_USER_LOGIN_WITH_OTP_REQ            uint16 = (3 << 10) | 55  // 0x0C37
	PROTO_NC_USER_LOGIN_WITH_OTP_LEGACY         uint16 = (3 << 10) | 56  // 0x0C38
	PROTO_NC_USER_US_LOGIN_REQ                  uint16 = (3 << 10) | 90  // 0x0C5A (NA/US client login)
	PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ      uint16 = (3 << 10) | 101 // 0x0C65
	PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK uint16 = (3 << 10) | 103 // 0x0C67

	// Department 4: CHAR / Character Selection & Zone Character Data
	PROTO_NC_CHAR_LOGIN_REQ           uint16 = (4 << 10) | 1  // 0x1001
	PROTO_NC_CHAR_LOGINFAIL_ACK       uint16 = (4 << 10) | 2  // 0x1002
	PROTO_NC_CHAR_LOGIN_ACK           uint16 = (4 << 10) | 3  // 0x1003
	PROTO_NC_CHAR_SERVER_IP_ACK       uint16 = (4 << 10) | 3  // 0x1003
	PROTO_NC_CHAR_BASE_CMD            uint16 = (4 << 10) | 56 // 0x1038 (105 bytes)
	PROTO_NC_CHAR_LOOK_CMD            uint16 = (4 << 10) | 57 // 0x1039 (4 bytes)
	PROTO_NC_QUEST_BUSY_LIST_CMD      uint16 = (4 << 10) | 58 // 0x103A
	PROTO_NC_QUEST_DONE_LIST_CMD      uint16 = (4 << 10) | 59 // 0x103B
	PROTO_NC_SKILL_ACTIVE_LIST_CMD    uint16 = (4 << 10) | 61 // 0x103D
	PROTO_NC_SKILL_PASSIVE_LIST_CMD   uint16 = (4 << 10) | 62 // 0x103E
	PROTO_NC_ITEM_LIST_CMD            uint16 = (4 << 10) | 71 // 0x1047
	PROTO_NC_CHAR_INFO_END_CMD        uint16 = (4 << 10) | 72 // 0x1048
	PROTO_NC_TITLE_LIST_CMD           uint16 = (4 << 10) | 73 // 0x1049
	PROTO_NC_CHAR_POINTS_ACK          uint16 = (4 << 10) | 91 // 0x105B

	// Department 5: CHAR Management
	PROTO_NC_CHAR_CREATE_REQ      uint16 = (5 << 10) | 1  // 0x1401
	PROTO_NC_CHAR_DELETE_REQ      uint16 = (5 << 10) | 7  // 0x1407 (C2S Character Delete)
	PROTO_NC_CHAR_CREATE_FAIL_ACK uint16 = (5 << 10) | 4  // 0x1404
	PROTO_NC_CHAR_CREATE_ACK      uint16 = (5 << 10) | 6  // 0x1406
	PROTO_NC_CHAR_DELETE_ACK      uint16 = (5 << 10) | 12 // 0x140C

	// Department 6: MAP / Zone Handshake & Navigation
	PROTO_NC_MAP_LOGIN_REQ        uint16 = (6 << 10) | 1  // 0x1801 (C2S TransferKey)
	PROTO_NC_CHAR_MAPLOGIN_ACK    uint16 = (6 << 10) | 2  // 0x1802 (S2C Map Login 242 bytes)
	PROTO_NC_MAP_LOGINFAIL_ACK    uint16 = (6 << 10) | 3  // 0x1803 (S2C Transfer failure)
	PROTO_NC_MAP_CLIENT_READY_REQ uint16 = (6 << 10) | 3  // 0x1803 (C2S ClientReady)
	PROTO_NC_MAP_CHANGE_CMD       uint16 = (6 << 10) | 9  // 0x1809
	PROTO_NC_MAP_CHANGE_ZONE_CMD  uint16 = (6 << 10) | 10 // 0x180A
	PROTO_NC_MAP_TELEPORT_REQ     uint16 = (6 << 10) | 26 // 0x181A
	PROTO_NC_MAP_TELEPORT_ACK     uint16 = (6 << 10) | 27 // 0x181B

	// Department 7: BRIEFINFO / Spawns & Visibilidad
	PROTO_NC_BRIEFINFO_UNSHARE_EQUIP_ACK  uint16 = (7 << 10) | 4  // 0x1C04
	PROTO_NC_BRIEFINFO_SHARE_EQUIP_ACK    uint16 = (7 << 10) | 5  // 0x1C05
	PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD uint16 = (7 << 10) | 6  // 0x1C06 (235 bytes)
	PROTO_NC_BRIEFINFO_CHARACTER_ACK      uint16 = (7 << 10) | 7  // 0x1C07 (Multi-player spawn)
	PROTO_NC_BRIEFINFO_LEAVE_CMD          uint16 = (7 << 10) | 14 // 0x1C0E (Despawn)

	// Department 8: ACT / Movimiento, Combate & Chat
	PROTO_NC_ACT_CHAT_REQ        uint16 = (8 << 10) | 1  // 0x2001
	PROTO_NC_ACT_CHAT_ACK        uint16 = (8 << 10) | 2  // 0x2002
	PROTO_NC_ACT_GM_NOTICE_ACK   uint16 = (8 << 10) | 17 // 0x2011
	PROTO_NC_ACT_STOP_REQ        uint16 = (8 << 10) | 18 // 0x2012
	PROTO_NC_ACT_STOP_ACK        uint16 = (8 << 10) | 19 // 0x2013
	PROTO_NC_ACT_WALK_REQ        uint16 = (8 << 10) | 23 // 0x2017
	PROTO_NC_ACT_WALK_ACK        uint16 = (8 << 10) | 24 // 0x2018
	PROTO_NC_ACT_RUN_REQ         uint16 = (8 << 10) | 25 // 0x2019
	PROTO_NC_ACT_RUN_ACK         uint16 = (8 << 10) | 26 // 0x201A
	PROTO_NC_ACT_TELEPORT_CMD    uint16 = (8 << 10) | 27 // 0x201B
	PROTO_NC_ACT_EMOTE_REQ       uint16 = (8 << 10) | 32 // 0x2020
	PROTO_NC_ACT_JUMP_REQ        uint16 = (8 << 10) | 36 // 0x2024
	PROTO_NC_ACT_JUMP_ACK        uint16 = (8 << 10) | 37 // 0x2025
	PROTO_NC_ACT_REST_START_REQ  uint16 = (8 << 10) | 39 // 0x2027
	PROTO_NC_ACT_REST_START_ACK  uint16 = (8 << 10) | 40 // 0x2028
	PROTO_NC_ACT_REST_STOP_REQ   uint16 = (8 << 10) | 42 // 0x202A
	PROTO_NC_ACT_REST_STOP_ACK   uint16 = (8 << 10) | 43 // 0x202B
	PROTO_NC_ACT_BLOCK_WALK_ACK  uint16 = (8 << 10) | 74 // 0x204A

	// Department 12: ITEM / Inventario
	PROTO_NC_ITEM_PREMIUM_LIST_ACK uint16 = (12 << 10) | 33 // 0x3021
	PROTO_NC_ITEM_REWARD_LIST_ACK  uint16 = (12 << 10) | 45 // 0x302D

	// Department 21: FRIEND / Social
	PROTO_NC_FRIEND_LIST_ACK uint16 = (21 << 10) | 7 // 0x5407

	// Department 22: CHAR In-Game / Kingdom Quest
	PROTO_NC_CHAR_INGAME_REQ   uint16 = (22 << 10) | 27 // 0x581B
	PROTO_NC_CHAR_INGAME_ACK   uint16 = (22 << 10) | 28 // 0x581C
	PROTO_NC_KQ_LIST_COUNT_ACK uint16 = (22 << 10) | 29 // 0x581D

	// Department 28: Settings & QuickBar
	PROTO_NC_CHAR_GET_QUICKBAR_REQ         uint16 = (28 << 10) | 2  // 0x7002
	PROTO_NC_CHAR_GET_QUICKBAR_ACK         uint16 = (28 << 10) | 3  // 0x7003
	PROTO_NC_CHAR_GET_QUICKBAR_STATE_REQ   uint16 = (28 << 10) | 4  // 0x7004
	PROTO_NC_CHAR_GET_QUICKBAR_STATE_ACK   uint16 = (28 << 10) | 5  // 0x7005
	PROTO_NC_CHAR_GET_GAME_SETTINGS_REQ    uint16 = (28 << 10) | 10 // 0x700A
	PROTO_NC_CHAR_GET_GAME_SETTINGS_ACK    uint16 = (28 << 10) | 11 // 0x700B
	PROTO_NC_CHAR_GET_CLIENT_SETTINGS_REQ  uint16 = (28 << 10) | 12 // 0x700C
	PROTO_NC_CHAR_GET_CLIENT_SETTINGS_ACK  uint16 = (28 << 10) | 13 // 0x700D
	PROTO_NC_CHAR_GET_SHORTCUTS_REQ        uint16 = (28 << 10) | 14 // 0x700E
	PROTO_NC_CHAR_GET_SHORTCUTS_ACK        uint16 = (28 << 10) | 15 // 0x700F
	PROTO_NC_CHAR_SAVE_QUICKBAR_REQ        uint16 = (28 << 10) | 16 // 0x7010
	PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ  uint16 = (28 << 10) | 17 // 0x7011
	PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ   uint16 = (28 << 10) | 20 // 0x7014
	PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ uint16 = (28 << 10) | 21 // 0x7015
	PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ       uint16 = (28 << 10) | 22 // 0x7016

	// Department 31: SYSTEM
	PROTO_NC_SYSTEM_UNKNOWN_REQ uint16 = (31 << 10) | 6 // 0x7C06
	PROTO_NC_SYSTEM_UNKNOWN_ACK uint16 = (31 << 10) | 7 // 0x7C07
)

// ServerError codes returned in PROTO_NC_USER_LOGINFAIL_ACK packets.
type ServerError uint16

const (
	ErrException          ServerError = 0x42 // Unknown Error / Exception
	ErrDatabase           ServerError = 0x43 // DB Error
	ErrWrongRegion        ServerError = 0x44 // Auth Failed / Wrong Region
	ErrInvalidCredentials ServerError = 0x45 // Please check ID or password
	ErrDatabaseError      ServerError = 0x46 // DB Error
	ErrBlocked            ServerError = 0x47 // The ID has been blocked
	ErrServerMaintenance  ServerError = 0x48 // World server maintenance
	ErrTimeout            ServerError = 0x49 // Timeout
	ErrLoginFailed        ServerError = 0x4A // Login Failed
	ErrAgreementMissing   ServerError = 0x4B // Please accept the agreement
)

// WorldStatus states for server selection display.
type WorldStatus uint8

const (
	WorldStatus_Offline         WorldStatus = 0
	WorldStatus_Maintenance     WorldStatus = 1
	WorldStatus_EmptyServer     WorldStatus = 2
	WorldStatus_Reserved        WorldStatus = 3
	WorldStatus_OfflineUnkError WorldStatus = 4
	WorldStatus_Full            WorldStatus = 5
	WorldStatus_Low             WorldStatus = 6
	WorldStatus_Medium          WorldStatus = 9
	WorldStatus_High            WorldStatus = 10
)

// CreateCharError codes returned in PROTO_NC_CHAR_CREATE_FAIL_ACK packets.
type CreateCharError uint16

const (
	CreateCharError_FailedToCreate CreateCharError = 130 // Failed to create the char
	CreateCharError_WrongClass     CreateCharError = 131 // You have chosen wrong class for your avatar
	CreateCharError_NameTaken      CreateCharError = 132 // The name is already taken by another user
	CreateCharError_ErrorInMaxSlot CreateCharError = 133 // There is an error in the Max slot Number
	CreateCharError_NameInUse      CreateCharError = 385 // The name is already in use
)

// Job identifies character classes in Fiesta Online.
type Job byte

const (
	Job_Fighter   Job = 1
	Job_Cleric    Job = 6
	Job_Archer    Job = 11
	Job_Mage      Job = 16
	Job_Trickster Job = 21
	Job_Crusader  Job = 26
)
