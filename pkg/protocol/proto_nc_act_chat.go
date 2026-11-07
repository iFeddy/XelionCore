package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_ACT_CHAT_REQ (size = 2 bytes)
type PROTO_NC_ACT_CHAT_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_ACT_NOTICE_CMD_SEND (size = 5 bytes)
type PROTO_NC_ACT_NOTICE_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_ACT_SCRIPT_MSG_CMD (size = 39 bytes)
type PROTO_NC_ACT_SCRIPT_MSG_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	SScriptMsgIndex                  [32]byte             // offset: 2, size: 32 bytes
	Nlen                             [5]byte              // offset: 34, size: 5 bytes
	SContent                         [0]byte              // offset: 39, flexible array member
}

// PROTO_NC_ACT_SCRIPT_MSG_WORLD_CMD (size = 39 bytes)
type PROTO_NC_ACT_SCRIPT_MSG_WORLD_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	SScriptMsgIndex                  [32]byte             // offset: 2, size: 32 bytes
	Nlen                             [5]byte              // offset: 34, size: 5 bytes
	SContent                         [0]byte              // offset: 39, flexible array member
}

// PROTO_NC_ACT_SOMEONEWHISPER_CMD (size = 23 bytes)
type PROTO_NC_ACT_SOMEONEWHISPER_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Talker                           Name5                // offset: 1, size: 20 bytes
	Flag                             uint8                // offset: 21, size: 1 bytes
	Len                              uint8                // offset: 22, size: 1 bytes
	Content                          [0]byte              // offset: 23, flexible array member
}

// PROTO_NC_ACT_WHISPERFAIL_ACK (size = 22 bytes)
type PROTO_NC_ACT_WHISPERFAIL_ACK_DATA struct {
	ErrorType                        uint16               // offset: 0, size: 2 bytes
	Receiver                         Name5                // offset: 2, size: 20 bytes
}

// PROTO_NC_ACT_WHISPERSUCCESS_ACK (size = 22 bytes)
type PROTO_NC_ACT_WHISPERSUCCESS_ACK_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Receiver                         Name5                // offset: 1, size: 20 bytes
	Len                              uint8                // offset: 21, size: 1 bytes
	Content                          [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_ACT_WHISPER_REQ (size = 22 bytes)
type PROTO_NC_ACT_WHISPER_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Receiver                         Name5                // offset: 1, size: 20 bytes
	Len                              uint8                // offset: 21, size: 1 bytes
	Content                          [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_CHARSAVE_2WLDMAN_CHAT_COLOR_CMD (size = 6 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_CHAT_COLOR_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NChatFontColorID                 uint8                // offset: 4, size: 1 bytes
	NChatBalloonColorID              uint8                // offset: 5, size: 1 bytes
}

// PROTO_NC_CHARSAVE_CHAT_COLOR_CMD (size = 6 bytes)
type PROTO_NC_CHARSAVE_CHAT_COLOR_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NChatFontColorID                 uint8                // offset: 4, size: 1 bytes
	NChatBalloonColorID              uint8                // offset: 5, size: 1 bytes
}

// PROTO_NC_CHAR_CHAT_COLOR_CMD (size = 6 bytes)
type PROTO_NC_CHAR_CHAT_COLOR_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NChatFontColorID                 uint8                // offset: 4, size: 1 bytes
	NChatBalloonColorID              uint8                // offset: 5, size: 1 bytes
}

// PROTO_NC_CHAR_ZONE_CHAT_COLOR_CMD (size = 8 bytes)
type PROTO_NC_CHAR_ZONE_CHAT_COLOR_CMD_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	NChatFontColorID                 uint8                // offset: 6, size: 1 bytes
	NChatBalloonColorID              uint8                // offset: 7, size: 1 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD (size = 1 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD_DATA struct {
	NGamerCnt                        uint8                // offset: 0, size: 1 bytes
	Emotion                          [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GAMBLE_PLAYERACT_CMD (size = 6 bytes)
type PROTO_NC_GAMBLE_PLAYERACT_CMD_DATA struct {
	Playerhandle                     uint16               // offset: 0, size: 2 bytes
	Effect_body                      uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_MISC_GM_CHAT_COLOR_REQ (size = 4 bytes)
type PROTO_NC_MISC_GM_CHAT_COLOR_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_USP_USER_CHARACTER_DELETE (size = 9 bytes)
type PROTO_NC_USP_USER_CHARACTER_DELETE_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NWorldNo                         uint8                // offset: 4, size: 1 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
}

// PROTO_NC_USP_USER_CHARACTER_INSERT (size = 50 bytes)
type PROTO_NC_USP_USER_CHARACTER_INSERT_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NWorldNo                         uint8                // offset: 4, size: 1 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
	SCharName                        [40]byte             // offset: 9, size: 40 bytes
	NClass                           uint8                // offset: 49, size: 1 bytes
}
