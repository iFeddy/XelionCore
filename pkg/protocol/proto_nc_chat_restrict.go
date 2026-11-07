package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_DB_CMD (size = 12 bytes)
type PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_DB_CMD_DATA struct {
	CharNo                           uint32               // offset: 0, size: 4 bytes
	BlockTime                        uint32               // offset: 4, size: 4 bytes
	DelayTime                        uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_WM_CMD (size = 14 bytes)
type PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_WM_CMD_DATA struct {
	WmHandle                         uint16               // offset: 0, size: 2 bytes
	CharNo                           uint32               // offset: 2, size: 4 bytes
	BlockTime                        uint32               // offset: 6, size: 4 bytes
	DelayTime                        uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_CHAR_GET_CHAT_BLOCK_SPAMER_CMD (size = 8 bytes)
type PROTO_NC_CHAR_GET_CHAT_BLOCK_SPAMER_CMD_DATA struct {
	BlockTime                        uint32               // offset: 0, size: 4 bytes
	DelayTime                        uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAR_GET_CHAT_BLOCK_SPAMER_DB_CMD (size = 12 bytes)
type PROTO_NC_CHAR_GET_CHAT_BLOCK_SPAMER_DB_CMD_DATA struct {
	CharNo                           uint32               // offset: 0, size: 4 bytes
	BlockTime                        uint32               // offset: 4, size: 4 bytes
	DelayTime                        uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_CHAT_RESTRICT_ADD_ACK (size = 24 bytes)
type PROTO_NC_CHAT_RESTRICT_ADD_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	ChatRestrictData                 [22]byte             // offset: 2, size: 22 bytes
}

// PROTO_NC_CHAT_RESTRICT_ADD_REQ (size = 20 bytes)
type PROTO_NC_CHAT_RESTRICT_ADD_REQ_DATA struct {
	SChatRestrictCharID              Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_ADD_ACK (size = 30 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_ADD_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NSessionHandle                   uint16               // offset: 2, size: 2 bytes
	ChatRestrictData                 [22]byte             // offset: 4, size: 22 bytes
	NChatRestrictCharNo              uint32               // offset: 26, size: 4 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_ADD_REQ (size = 26 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_ADD_REQ_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	SChatRestrictCharID              Name5                // offset: 6, size: 20 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_DEL_ACK (size = 8 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_DEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NSessionHandle                   uint16               // offset: 2, size: 2 bytes
	NChatRestrictCharNo              uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_DEL_ALL_ACK (size = 4 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_DEL_ALL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NSessionHandle                   uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_DEL_ALL_REQ (size = 6 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_DEL_ALL_REQ_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_DEL_REQ (size = 10 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_DEL_REQ_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	NChatRestrictCharNo              uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_CHAT_RESTRICT_DB_LIST_CMD (size = 2 bytes)
type PROTO_NC_CHAT_RESTRICT_DB_LIST_CMD_DATA struct {
	NChatRestrictListCnt             uint16               // offset: 0, size: 2 bytes
	ChatRestrictData                 [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHAT_RESTRICT_DEL_ACK (size = 22 bytes)
type PROTO_NC_CHAT_RESTRICT_DEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	SChatRestrictCharID              Name5                // offset: 2, size: 20 bytes
}

// PROTO_NC_CHAT_RESTRICT_DEL_ALL_ACK (size = 2 bytes)
type PROTO_NC_CHAT_RESTRICT_DEL_ALL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_CHAT_RESTRICT_DEL_REQ (size = 20 bytes)
type PROTO_NC_CHAT_RESTRICT_DEL_REQ_DATA struct {
	SChatRestrictCharID              Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_CHAT_RESTRICT_LIST_CMD (size = 2 bytes)
type PROTO_NC_CHAT_RESTRICT_LIST_CMD_DATA struct {
	NChatRestrictListCnt             uint16               // offset: 0, size: 2 bytes
	ChatRestrictData                 [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MISC_GET_CHAT_BLOCK_SPAM_FILTER_CMD (size = 3 bytes)
type PROTO_NC_MISC_GET_CHAT_BLOCK_SPAM_FILTER_CMD_DATA struct {
	BIsSpamSpread                    uint8                // offset: 0, size: 1 bytes
	BIsSpamBlock                     uint8                // offset: 1, size: 1 bytes
	NSpamFilterRate                  uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_MISC_GET_CHAT_BLOCK_SPAM_FILTER_DB_CMD (size = 11 bytes)
type PROTO_NC_MISC_GET_CHAT_BLOCK_SPAM_FILTER_DB_CMD_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	BIsSpamSpread                    uint8                // offset: 8, size: 1 bytes
	BIsSpamBlock                     uint8                // offset: 9, size: 1 bytes
	NSpamFilterRate                  uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_MISC_SET_CHAT_BLOCK_SPAM_FILTER_CMD (size = 3 bytes)
type PROTO_NC_MISC_SET_CHAT_BLOCK_SPAM_FILTER_CMD_DATA struct {
	BIsSpamSpread                    uint8                // offset: 0, size: 1 bytes
	BIsSpamBlock                     uint8                // offset: 1, size: 1 bytes
	NSpamFilterRate                  uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_MISC_SET_CHAT_BLOCK_SPAM_FILTER_DB_CMD (size = 11 bytes)
type PROTO_NC_MISC_SET_CHAT_BLOCK_SPAM_FILTER_DB_CMD_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	BIsSpamSpread                    uint8                // offset: 8, size: 1 bytes
	BIsSpamBlock                     uint8                // offset: 9, size: 1 bytes
	NSpamFilterRate                  uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_MISC_SPAMMER_CHAT_BAN_ACK (size = 22 bytes)
type PROTO_NC_MISC_SPAMMER_CHAT_BAN_ACK_DATA struct {
	SSpammerName                     Name5                // offset: 0, size: 20 bytes
	NError                           uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_MISC_SPAMMER_CHAT_BAN_REQ (size = 20 bytes)
type PROTO_NC_MISC_SPAMMER_CHAT_BAN_REQ_DATA struct {
	SSpammerName                     Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_ACK (size = 2 bytes)
type PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_DB_ACK (size = 30 bytes)
type PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_DB_ACK_DATA struct {
	NGMCharNo                        uint32               // offset: 0, size: 4 bytes
	SSpammerName                     Name5                // offset: 4, size: 20 bytes
	NChatBan                         uint32               // offset: 24, size: 4 bytes
	NError                           uint16               // offset: 28, size: 2 bytes
}

// PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_DB_REQ (size = 28 bytes)
type PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_DB_REQ_DATA struct {
	NGMCharNo                        uint32               // offset: 0, size: 4 bytes
	SSpammerName                     Name5                // offset: 4, size: 20 bytes
	NChatBan                         uint32               // offset: 24, size: 4 bytes
}

// PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_REQ (size = 24 bytes)
type PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_REQ_DATA struct {
	NGMCharNo                        uint32               // offset: 0, size: 4 bytes
	SSpammerName                     Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_MISC_SPAMMER_REPORT_ACK (size = 2 bytes)
type PROTO_NC_MISC_SPAMMER_REPORT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MISC_SPAMMER_REPORT_REQ (size = 20 bytes)
type PROTO_NC_MISC_SPAMMER_REPORT_REQ_DATA struct {
	SSpammerName                     Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_MISC_SPAMMER_SET_DB_CHAT_BAN_ACK (size = 30 bytes)
type PROTO_NC_MISC_SPAMMER_SET_DB_CHAT_BAN_ACK_DATA struct {
	NGMCharNo                        uint32               // offset: 0, size: 4 bytes
	SSpammerName                     Name5                // offset: 4, size: 20 bytes
	NChatBan                         uint32               // offset: 24, size: 4 bytes
	NError                           uint16               // offset: 28, size: 2 bytes
}

// PROTO_NC_MISC_SPAMMER_SET_DB_CHAT_BAN_REQ (size = 28 bytes)
type PROTO_NC_MISC_SPAMMER_SET_DB_CHAT_BAN_REQ_DATA struct {
	NGMCharNo                        uint32               // offset: 0, size: 4 bytes
	SSpammerName                     Name5                // offset: 4, size: 20 bytes
	NChatBan                         uint32               // offset: 24, size: 4 bytes
}

// PROTO_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD (size = 24 bytes)
type PROTO_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD_DATA struct {
	SSpammerName                     Name5                // offset: 0, size: 20 bytes
	NChatBan                         uint32               // offset: 20, size: 4 bytes
}
