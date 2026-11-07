package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_LOG_USER_LOGOUT (size = 9 bytes)
type PROTO_NC_LOG_USER_LOGOUT_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
	Worldnum                         uint8                // offset: 4, size: 1 bytes
	Playmin                          uint16               // offset: 5, size: 2 bytes
	Error                            uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_LOG_USER_LOGOUT_SEND (size = 12 bytes)
type PROTO_NC_LOG_USER_LOGOUT_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Log                              [9]byte              // offset: 3, size: 9 bytes
}

// PROTO_NC_MISC_CLIENT_DEBUG_MSG_CMD (size = 4 bytes)
type PROTO_NC_MISC_CLIENT_DEBUG_MSG_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Msg_size                         uint16               // offset: 2, size: 2 bytes
	Msg                              [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_MISC_CLIENT_LOADING_BUG_DETECT_CMD (size = 64 bytes)
type PROTO_NC_MISC_CLIENT_LOADING_BUG_DETECT_CMD_DATA struct {
	SLoadingMapName                  Name3                // offset: 0, size: 12 bytes
	NFailType                        uint16               // offset: 12, size: 2 bytes
	NLoadingFailCnt                  uint16               // offset: 14, size: 2 bytes
	SFromMapName                     Name3                // offset: 16, size: 12 bytes
	SFromZoneIP                      [16]byte             // offset: 28, size: 16 bytes
	SToZoneIP                        [16]byte             // offset: 44, size: 16 bytes
	NToZonePort                      uint16               // offset: 60, size: 2 bytes
	NErrorDataCnt                    uint16               // offset: 62, size: 2 bytes
	ErrorData                        [0]byte              // offset: 64, flexible array member
}

// PROTO_NC_MISC_CONNECTER_ACK (size = 131 bytes)
type PROTO_NC_MISC_CONNECTER_ACK_DATA struct {
	Worldno                          uint8                // offset: 0, size: 1 bytes
	Connectercount                   uint16               // offset: 1, size: 2 bytes
	NEachClass                       [64]uint16           // offset: 3, size: 128 bytes
}

// PROTO_NC_MISC_CONNECTFROMWHERE_ACK (size = 7 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Connectfrom                      uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK (size = 9 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK_DATA struct {
	HZone                            uint16               // offset: 0, size: 2 bytes
	ToZoneAck                        [7]byte              // offset: 2, size: 7 bytes
}

// PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ (size = 12 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ_DATA struct {
	HZone                            uint16               // offset: 0, size: 2 bytes
	FromZoneReq                      [10]byte             // offset: 2, size: 10 bytes
}

// PROTO_NC_MISC_CONNECTFROMWHERE_REQ (size = 10 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Connectip                        [4]byte              // offset: 6, size: 4 bytes
}

// PROTO_NC_MISC_CS_ACK (size = 2 bytes)
type PROTO_NC_MISC_CS_ACK_DATA struct {
	DataSize                         uint16               // offset: 0, size: 2 bytes
	Data                             [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MISC_DB_CLIENT_LOADING_BUG_DETECT_CMD (size = 72 bytes)
type PROTO_NC_MISC_DB_CLIENT_LOADING_BUG_DETECT_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NWorldTime                       int32                // offset: 4, size: 4 bytes
	SLoadingMapName                  Name3                // offset: 8, size: 12 bytes
	NFailType                        uint16               // offset: 20, size: 2 bytes
	NLoadingFailCnt                  uint16               // offset: 22, size: 2 bytes
	SFromMapName                     Name3                // offset: 24, size: 12 bytes
	SFromZoneIP                      [16]byte             // offset: 36, size: 16 bytes
	SToZoneIP                        [16]byte             // offset: 52, size: 16 bytes
	NToZonePort                      uint16               // offset: 68, size: 2 bytes
	NErrorDataCnt                    uint16               // offset: 70, size: 2 bytes
	ErrorData                        [0]byte              // offset: 72, flexible array member
}

// PROTO_NC_MISC_EVENTNPC_STANDEND_CLIENT_CMD (size = 32 bytes)
type PROTO_NC_MISC_EVENTNPC_STANDEND_CLIENT_CMD_DATA struct {
	SMobIndex                        Name8                // offset: 0, size: 32 bytes
}

// PROTO_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD (size = 32 bytes)
type PROTO_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD_DATA struct {
	SMobIndex                        Name8                // offset: 0, size: 32 bytes
}

// PROTO_NC_MISC_EVENTNPC_STANDSTART_CLIENT_CMD (size = 32 bytes)
type PROTO_NC_MISC_EVENTNPC_STANDSTART_CLIENT_CMD_DATA struct {
	SMobIndex                        Name8                // offset: 0, size: 32 bytes
}

// PROTO_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD (size = 32 bytes)
type PROTO_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD_DATA struct {
	SMobIndex                        Name8                // offset: 0, size: 32 bytes
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM (size = 15 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NWorldNo                         uint8                // offset: 8, size: 1 bytes
	NError                           uint16               // offset: 9, size: 2 bytes
	HZone                            uint16               // offset: 11, size: 2 bytes
	NQuestID                         uint16               // offset: 13, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC (size = 13 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NWorldNo                         uint8                // offset: 8, size: 1 bytes
	HZone                            uint16               // offset: 9, size: 2 bytes
	NQuestID                         uint16               // offset: 11, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2Z (size = 13 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2Z_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NWorldNo                         uint8                // offset: 8, size: 1 bytes
	NError                           uint16               // offset: 9, size: 2 bytes
	NQuestID                         uint16               // offset: 11, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_Z2WM (size = 11 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_Z2WM_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NWorldNo                         uint8                // offset: 8, size: 1 bytes
	NQuestID                         uint16               // offset: 9, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_L20_CMD (size = 3 bytes)
type PROTO_NC_MISC_EVENT_L20_CMD_DATA struct {
	NEventNo                         uint8                // offset: 0, size: 1 bytes
	NError                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_L20_DB_ACK (size = 13 bytes)
type PROTO_NC_MISC_EVENT_L20_DB_ACK_DATA struct {
	NWMHandle                        uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	NEventNo                         uint8                // offset: 6, size: 1 bytes
	NCharNo                          uint32               // offset: 7, size: 4 bytes
	NError                           uint16               // offset: 11, size: 2 bytes
}

// PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD (size = 6 bytes)
type PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NPlayTime                        uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_MISC_GAMETIME_ACK (size = 3 bytes)
type PROTO_NC_MISC_GAMETIME_ACK_DATA struct {
	Hour                             uint8                // offset: 0, size: 1 bytes
	Minute                           uint8                // offset: 1, size: 1 bytes
	Second                           uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_MISC_HACK_SCAN_STORE_CMD (size = 2 bytes)
type PROTO_NC_MISC_HACK_SCAN_STORE_CMD_DATA struct {
	NDescSize                        uint16               // offset: 0, size: 2 bytes
	SDesc                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD (size = 51 bytes)
type PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
	NWorldNo                         uint8                // offset: 4, size: 1 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
	SIP                              [20]byte             // offset: 9, size: 20 bytes
	SCharID                          Name5                // offset: 29, size: 20 bytes
	ScanInfo                         uint16               // offset: 49, size: 2 bytes
}

// PROTO_NC_MISC_HEARTBEAT_ACK_SEND (size = 3 bytes)
type PROTO_NC_MISC_HEARTBEAT_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_HEARTBEAT_REQ_SEND (size = 3 bytes)
type PROTO_NC_MISC_HEARTBEAT_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB (size = 130 bytes)
type PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB_DATA struct {
	ID                               uint32               // offset: 0, size: 4 bytes
	HTool                            uint16               // offset: 4, size: 2 bytes
	ToolSendToWM                     uint32               // offset: 6, size: 4 bytes
	ToolRecvFromWM                   uint32               // offset: 10, size: 4 bytes
	TargetCharName                   Name5                // offset: 14, size: 20 bytes
	HWMClient                        uint16               // offset: 34, size: 2 bytes
	WMSendClient                     uint32               // offset: 36, size: 4 bytes
	WMRecvFromClient                 uint32               // offset: 40, size: 4 bytes
	ClientSendZone                   uint32               // offset: 44, size: 4 bytes
	ClientRecvZone                   uint32               // offset: 48, size: 4 bytes
	HZoneClient                      uint16               // offset: 52, size: 2 bytes
	ZoneSendToDB                     uint32               // offset: 54, size: 4 bytes
	ZoneRecvFromDB                   uint32               // offset: 58, size: 4 bytes
	ZoneSendToClient                 uint32               // offset: 62, size: 4 bytes
	ZoneRecvFromClient               uint32               // offset: 66, size: 4 bytes
	DBRecvFromZone                   uint32               // offset: 70, size: 4 bytes
	DBSendToZone                     uint32               // offset: 74, size: 4 bytes
	DBSendToODBC                     uint32               // offset: 78, size: 4 bytes
	DBRecvFromODBC                   uint32               // offset: 82, size: 4 bytes
	Header                           NETPACKETZONEHEADER  // offset: 86, size: 6 bytes
	Buylot                           [38]byte             // offset: 92, size: 38 bytes
}

// PROTO_NC_MISC_PINGTEST_TOOL_WM_DB (size = 38 bytes)
type PROTO_NC_MISC_PINGTEST_TOOL_WM_DB_DATA struct {
	ID                               uint32               // offset: 0, size: 4 bytes
	HTool                            uint16               // offset: 4, size: 2 bytes
	ToolSendToWM                     uint32               // offset: 6, size: 4 bytes
	ToolRecvFromWM                   uint32               // offset: 10, size: 4 bytes
	WMSendDB                         uint32               // offset: 14, size: 4 bytes
	WMRecvFromDB                     uint32               // offset: 18, size: 4 bytes
	DBRecvFromWM                     uint32               // offset: 22, size: 4 bytes
	DBSendToWM                       uint32               // offset: 26, size: 4 bytes
	DBSendToODBC                     uint32               // offset: 30, size: 4 bytes
	DBRecvFromODBC                   uint32               // offset: 34, size: 4 bytes
}

// PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE (size = 34 bytes)
type PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DATA struct {
	ID                               uint32               // offset: 0, size: 4 bytes
	HTool                            uint16               // offset: 4, size: 2 bytes
	ToolSendToWM                     uint32               // offset: 6, size: 4 bytes
	ToolRecvFromWM                   uint32               // offset: 10, size: 4 bytes
	WMTargetZoneNo                   uint32               // offset: 14, size: 4 bytes
	WMSendZone                       uint32               // offset: 18, size: 4 bytes
	WMRecvFromZone                   uint32               // offset: 22, size: 4 bytes
	ZoneRecvFromWMS                  uint32               // offset: 26, size: 4 bytes
	ZoneSendToWMS                    uint32               // offset: 30, size: 4 bytes
}

// PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB (size = 58 bytes)
type PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB_DATA struct {
	ID                               uint32               // offset: 0, size: 4 bytes
	HTool                            uint16               // offset: 4, size: 2 bytes
	ToolSendToWM                     uint32               // offset: 6, size: 4 bytes
	ToolRecvFromWM                   uint32               // offset: 10, size: 4 bytes
	WMTargetZoneNo                   uint32               // offset: 14, size: 4 bytes
	WMSendZone                       uint32               // offset: 18, size: 4 bytes
	WMRecvFromZone                   uint32               // offset: 22, size: 4 bytes
	ZoneSendToDB                     uint32               // offset: 26, size: 4 bytes
	ZoneRecvFromDB                   uint32               // offset: 30, size: 4 bytes
	ZoneSendToWMS                    uint32               // offset: 34, size: 4 bytes
	ZoneRecvFromWMS                  uint32               // offset: 38, size: 4 bytes
	DBRecvFromZone                   uint32               // offset: 42, size: 4 bytes
	DBSendToZone                     uint32               // offset: 46, size: 4 bytes
	DBSendToODBC                     uint32               // offset: 50, size: 4 bytes
	DBRecvFromODBC                   uint32               // offset: 54, size: 4 bytes
}

// PROTO_NC_MISC_RESTMINUTE_CMD (size = 3 bytes)
type PROTO_NC_MISC_RESTMINUTE_CMD_DATA struct {
	Flag                             uint8                // offset: 0, size: 1 bytes
	Second                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_RESTMINUTE_CMD_SEND (size = 6 bytes)
type PROTO_NC_MISC_RESTMINUTE_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [3]byte              // offset: 3, size: 3 bytes
}

// PROTO_NC_MISC_S2SCONNECTION_ACK (size = 3 bytes)
type PROTO_NC_MISC_S2SCONNECTION_ACK_DATA struct {
	Echo_data                        uint8                // offset: 0, size: 1 bytes
	Error                            uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_S2SCONNECTION_RDY_SEND (size = 3 bytes)
type PROTO_NC_MISC_S2SCONNECTION_RDY_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MISC_S2SCONNECTION_REQ (size = 7 bytes)
type PROTO_NC_MISC_S2SCONNECTION_REQ_DATA struct {
	Echo_data                        uint8                // offset: 0, size: 1 bytes
	Server_to_id                     uint8                // offset: 1, size: 1 bytes
	Server_from_id                   uint8                // offset: 2, size: 1 bytes
	Server_from_world_num            uint8                // offset: 3, size: 1 bytes
	Server_from_zone_num             uint8                // offset: 4, size: 1 bytes
	Key                              uint16               // offset: 5, size: 2 bytes
}

// PROTO_NC_MISC_S2SCONNECTION_REQ_SEND (size = 10 bytes)
type PROTO_NC_MISC_S2SCONNECTION_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [7]byte              // offset: 3, size: 7 bytes
}

// PROTO_NC_MISC_SEED_ACK (size = 2 bytes)
type PROTO_NC_MISC_SEED_ACK_DATA struct {
	Seed                             uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD (size = 37 bytes)
type PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD_DATA struct {
	DCurrentTM                       [36]byte             // offset: 0, size: 36 bytes
	NTimeZone                        int8                 // offset: 36, size: 1 bytes
}

// PROTO_NC_MISC_TIMEFROMWORLD_CMD (size = 36 bytes)
type PROTO_NC_MISC_TIMEFROMWORLD_CMD_DATA struct {
	TimeStr                          [36]byte             // offset: 0, size: 36 bytes
}

// PROTO_NC_MISC_USER_COUNT_CMD (size = 4 bytes)
type PROTO_NC_MISC_USER_COUNT_CMD_DATA struct {
	NumOfUser                        uint16               // offset: 0, size: 2 bytes
	MaxOfUser                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_MISC_WEB_DB_KEY_ACK (size = 40 bytes)
type PROTO_NC_MISC_WEB_DB_KEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NWorldSessionHandle              uint16               // offset: 2, size: 2 bytes
	NUserNo                          uint32               // offset: 4, size: 4 bytes
	WebKey                           [32]byte             // offset: 8, size: 32 bytes
}

// PROTO_NC_MISC_WEB_DB_KEY_REQ (size = 6 bytes)
type PROTO_NC_MISC_WEB_DB_KEY_REQ_DATA struct {
	NWorldSessionHandle              uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_MISC_WEB_KEY_ACK (size = 34 bytes)
type PROTO_NC_MISC_WEB_KEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	WebKey                           [32]byte             // offset: 2, size: 32 bytes
}

// PROTO_NC_MISC_WEB_KEY_REQ (size = 1 bytes)
type PROTO_NC_MISC_WEB_KEY_REQ_DATA struct {}

// PROTO_NC_MISC_WHSHANDLEFIX_CMD (size = 9 bytes)
type PROTO_NC_MISC_WHSHANDLEFIX_CMD_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Whshandle                        uint16               // offset: 6, size: 2 bytes
	Flag                             uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_MISC_WHSHANDLEREPAIR_CMD (size = 11 bytes)
type PROTO_NC_MISC_WHSHANDLEREPAIR_CMD_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Whshandle_old                    uint16               // offset: 6, size: 2 bytes
	Whshandle_new                    uint16               // offset: 8, size: 2 bytes
	Flag                             uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_MISC_XTRAP2_CLIENT_DATA_CMD (size = 2 bytes)
type PROTO_NC_MISC_XTRAP2_CLIENT_DATA_CMD_DATA struct {
	DataSize                         uint16               // offset: 0, size: 2 bytes
	Data                             [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_USER_CONNECTCUT2WORLDMANAGER_CMD (size = 8 bytes)
type PROTO_NC_USER_CONNECTCUT2WORLDMANAGER_CMD_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
	Wldmanhandle                     uint16               // offset: 4, size: 2 bytes
	Reason                           uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_USER_CONNECTCUT2ZONE_CMD (size = 4 bytes)
type PROTO_NC_USER_CONNECTCUT2ZONE_CMD_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_USER_CONNECTCUT2ZONE_CMD_SEND (size = 13 bytes)
type PROTO_NC_USER_CONNECTCUT2ZONE_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 3, size: 6 bytes
	Cmd                              uint32               // offset: 9, size: 4 bytes
}

// PROTO_NC_USER_CONNECTCUT_CMD (size = 2 bytes)
type PROTO_NC_USER_CONNECTCUT_CMD_DATA struct {
	Reason                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_USER_CONNECTCUT_CMD_SEND (size = 5 bytes)
type PROTO_NC_USER_CONNECTCUT_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_USER_CREATE_OTP_ACK (size = 36 bytes)
type PROTO_NC_USER_CREATE_OTP_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	HWM                              uint16               // offset: 2, size: 2 bytes
	SOTP                             Name8                // offset: 4, size: 32 bytes
}

// PROTO_NC_USER_CREATE_OTP_REQ (size = 266 bytes)
type PROTO_NC_USER_CREATE_OTP_REQ_DATA struct {
	HWM                              uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	SUserID                          [256]byte            // offset: 6, size: 256 bytes
	NClientIP                        uint32               // offset: 262, size: 4 bytes
}

// PROTO_NC_USER_KICKOFFFROMWORLD_CMD (size = 4 bytes)
type PROTO_NC_USER_KICKOFFFROMWORLD_CMD_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_USER_KICKOFFFROMWORLD_CMD_SEND (size = 13 bytes)
type PROTO_NC_USER_KICKOFFFROMWORLD_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 3, size: 6 bytes
	Cmd                              uint32               // offset: 9, size: 4 bytes
}

// PROTO_NC_USER_LOGINWORLDFAIL_ACK (size = 2 bytes)
type PROTO_NC_USER_LOGINWORLDFAIL_ACK_DATA struct {
	Errorcode                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_USER_LOGINWORLDFAIL_ACK_SEND (size = 5 bytes)
type PROTO_NC_USER_LOGINWORLDFAIL_ACK_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Ack                              uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_USER_LOGINWORLD_ACK (size = 3 bytes)
type PROTO_NC_USER_LOGINWORLD_ACK_DATA struct {
	Worldmanager                     uint16               // offset: 0, size: 2 bytes
	Numofavatar                      uint8                // offset: 2, size: 1 bytes
	Avatar                           [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_USER_LOGINWORLD_REQ (size = 320 bytes)
type PROTO_NC_USER_LOGINWORLD_REQ_DATA struct {
	User                             [256]byte            // offset: 0, size: 256 bytes
	Validate_new                     [32]uint16           // offset: 256, size: 64 bytes
}

// PROTO_NC_USER_LOGOUT_DB (size = 9 bytes)
type PROTO_NC_USER_LOGOUT_DB_DATA struct {
	Usrregnum                        uint32               // offset: 0, size: 4 bytes
	Worldnum                         uint8                // offset: 4, size: 1 bytes
	Playmin                          uint16               // offset: 5, size: 2 bytes
	Error                            uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_USER_LOGOUT_DB_SEND (size = 12 bytes)
type PROTO_NC_USER_LOGOUT_DB_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Db                               [9]byte              // offset: 3, size: 9 bytes
}

// PROTO_NC_USER_LOGOUT_LAST_DAY_ACK (size = 10 bytes)
type PROTO_NC_USER_LOGOUT_LAST_DAY_ACK_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	TLastSec                         int32                // offset: 6, size: 4 bytes
}

// PROTO_NC_USER_LOGOUT_LAST_DAY_REQ (size = 6 bytes)
type PROTO_NC_USER_LOGOUT_LAST_DAY_REQ_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_USER_NORMALLOGOUT_CMD (size = 1 bytes)
type PROTO_NC_USER_NORMALLOGOUT_CMD_DATA struct {
	LogoutType                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_USER_POSSIBLE_NEW_CONNECT_CMD (size = 6 bytes)
type PROTO_NC_USER_POSSIBLE_NEW_CONNECT_CMD_DATA struct {
	NWorldManagerHandle              uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_USER_RETURN_CHECK_ACK (size = 11 bytes)
type PROTO_NC_USER_RETURN_CHECK_ACK_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
	BReturnUser                      uint8                // offset: 6, size: 1 bytes
	TLastSec                         int32                // offset: 7, size: 4 bytes
}

// PROTO_NC_USER_RETURN_CHECK_REQ (size = 6 bytes)
type PROTO_NC_USER_RETURN_CHECK_REQ_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NUserNo                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_USER_SET_RETURN_CMD (size = 4 bytes)
type PROTO_NC_USER_SET_RETURN_CMD_DATA struct {
	NUserNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_USER_TEENAGER_ACK (size = 313 bytes)
type PROTO_NC_USER_TEENAGER_ACK_DATA struct {
	SUserID                          [256]byte            // offset: 0, size: 256 bytes
	NHandle                          uint16               // offset: 256, size: 2 bytes
	Data                             [55]byte             // offset: 258, size: 55 bytes
}

// PROTO_NC_USER_TEENAGER_CMD (size = 55 bytes)
type PROTO_NC_USER_TEENAGER_CMD_DATA struct {
	SUserKey                         [50]byte             // offset: 0, size: 50 bytes
	NType                            uint8                // offset: 50, size: 1 bytes
	NMinute                          uint32               // offset: 51, size: 4 bytes
}

// PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD (size = 1 bytes)
type PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD_DATA struct {
	BCanSkip                         uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_USER_WILLLOGIN_ACK (size = 348 bytes)
type PROTO_NC_USER_WILLLOGIN_ACK_DATA struct {
	Will_login_req                   [347]byte            // offset: 0, size: 347 bytes
	Ok                               uint8                // offset: 347, size: 1 bytes
}

// PROTO_NC_USER_WILLLOGIN_REQ (size = 347 bytes)
type PROTO_NC_USER_WILLLOGIN_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Willlogin_worldno                uint8                // offset: 2, size: 1 bytes
	Usrregnum                        uint32               // offset: 3, size: 4 bytes
	Userid                           [256]byte            // offset: 7, size: 256 bytes
	Validate_new                     [32]uint16           // offset: 263, size: 64 bytes
	Spawnapps                        Name5                // offset: 327, size: 20 bytes
}

// PROTO_NC_USER_WILL_WORLD_SELECT_ACK (size = 34 bytes)
type PROTO_NC_USER_WILL_WORLD_SELECT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	SOTP                             Name8                // offset: 2, size: 32 bytes
}

// PROTO_NC_USER_WILL_WORLD_SELECT_REQ (size = 1 bytes)
type PROTO_NC_USER_WILL_WORLD_SELECT_REQ_DATA struct {}
