package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_misc.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_MISC_CHAT_BLOCK_SPAMER_BLOCKTIME_CMD (sizeof = 4 bytes)
type PROTO_NC_MISC_CHAT_BLOCK_SPAMER_BLOCKTIME_CMD_DATA struct {
	NBlockTime                       uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_MISC_CHAT_BLOCK_SPAMER_BLOCKTIME_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_MISC_CONNECTFROMWHERE_ACK (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_CONNECTFROMWHERE_ACK_DATA)

// PROTO_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD (sizeof = 32 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD_DATA)

// PROTO_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD (sizeof = 32 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD_DATA)

// PROTO_NC_MISC_GM_CHAT_COLOR_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_GM_CHAT_COLOR_REQ_DATA)

// PROTO_NC_MISC_HEARTBEAT_REQ_SEND (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_HEARTBEAT_REQ_SEND)

// PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB (sizeof = 90 bytes)
type PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB_DATA struct {
	ID                               uint32               // offset: 0, size: 4 bytes
	ClientSendZone                   uint32               // offset: 4, size: 4 bytes
	ClientRecvZone                   uint32               // offset: 8, size: 4 bytes
	HZoneClient                      uint16               // offset: 12, size: 2 bytes
	ZoneSendToDB                     uint32               // offset: 14, size: 4 bytes
	ZoneRecvFromDB                   uint32               // offset: 18, size: 4 bytes
	ZoneSendToClient                 uint32               // offset: 22, size: 4 bytes
	ZoneRecvFromClient               uint32               // offset: 26, size: 4 bytes
	DBRecvFromZone                   uint32               // offset: 30, size: 4 bytes
	DBSendToZone                     uint32               // offset: 34, size: 4 bytes
	DBSendToODBC                     uint32               // offset: 38, size: 4 bytes
	DBRecvFromODBC                   uint32               // offset: 42, size: 4 bytes
	Header                           NETPACKETZONEHEADER  // offset: 46, size: 6 bytes
	Buylot                           [38]byte             // offset: 52, size: 38 bytes
}

func (s *PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB_DATA) FixedWireSize() int { return 90 }

// PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB (sizeof = 130 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB_DATA)

// PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE (sizeof = 34 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DATA)

// PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB (sizeof = 58 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB_DATA)

// PROTO_NC_MISC_RESTMINUTE_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_RESTMINUTE_CMD_DATA)

// PROTO_NC_MISC_S2SCONNECTION_ACK (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_S2SCONNECTION_ACK_DATA)

// PROTO_NC_MISC_S2SCONNECTION_REQ (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_S2SCONNECTION_REQ_DATA)

// PROTO_NC_MISC_SEED_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_SEED_ACK_DATA)

// PROTO_NC_MISC_SERVERPARAMETER_ACK (sizeof = 3 bytes)
type PROTO_NC_MISC_SERVERPARAMETER_ACK_DATA struct {
	Shoutlevel                       uint8                // offset: 0, size: 1 bytes
	Shoutdelay                       uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_MISC_SERVERPARAMETER_ACK_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_REQ (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_REQ_DATA)

// PROTO_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD_DATA)

// PROTO_NC_MISC_TIMEFROMWORLD_CMD (sizeof = 36 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_TIMEFROMWORLD_CMD_DATA)

// PROTO_NC_MISC_WHSHANDLEREPAIR_CMD (sizeof = 11 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MISC_WHSHANDLEREPAIR_CMD_DATA)

// PROTO_NC_MISC_ZONERINGLINKTEST_RNG (sizeof = 8 bytes)
type PROTO_NC_MISC_ZONERINGLINKTEST_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Mechines                         uint8                // offset: 7, size: 1 bytes
	Addr                             [0]byte              // offset: 8, flexible array member
}

func (s *PROTO_NC_MISC_ZONERINGLINKTEST_RNG_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_USER_CONNECTCUT2WORLDMANAGER_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_USER_CONNECTCUT2WORLDMANAGER_CMD_DATA)

// PROTO_NC_USER_CONNECTCUT2ZONE_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_USER_CONNECTCUT2ZONE_CMD_DATA)

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK (sizeof = 4 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NUsedItemID                      uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ (sizeof = 1 bytes)
type PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ_DATA struct {
	NInvenSlot                       uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK (sizeof = 12 bytes)
type PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	NUsedItemID                      uint16               // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_REQ (sizeof = 38 bytes)
type PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_REQ_DATA struct {
	NCharIDChangeItemKey             uint64               // offset: 0, size: 8 bytes
	NCharIDChangeItemSlot            ITEM_INVEN           // offset: 8, size: 2 bytes
	BInit                            uint8                // offset: 10, size: 1 bytes
	NUsedItemID                      uint16               // offset: 11, size: 2 bytes
	NUserNo                          uint32               // offset: 13, size: 4 bytes
	CommonData                       [21]byte             // offset: 17, size: 21 bytes
}

func (s *PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_REQ_DATA) FixedWireSize() int { return 38 }

// PROTO_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD (sizeof = 45 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD_DATA)

// PROTO_NC_USER_NORMALLOGOUT_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_USER_NORMALLOGOUT_CMD_DATA)

// PROTO_NC_USER_POSSIBLE_NEW_CONNECT_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_USER_POSSIBLE_NEW_CONNECT_CMD_DATA)

// PROTO_NC_USER_REGISENUMBER_ACK (sizeof = 4 bytes)
type PROTO_NC_USER_REGISENUMBER_ACK_DATA struct {
	Userregnum                       uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_USER_REGISENUMBER_ACK_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_USER_USE_BEAUTY_SHOP_CMD (sizeof = 1 bytes)
type PROTO_NC_USER_USE_BEAUTY_SHOP_CMD_DATA struct {}

func (s *PROTO_NC_USER_USE_BEAUTY_SHOP_CMD_DATA) FixedWireSize() int { return 1 }
