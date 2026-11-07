package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHAR_MASPUP_CMD (size = 1 bytes)
type PROTO_NC_CHAR_MASPUP_CMD_DATA struct {}

// PROTO_NC_HOLY_PROMISE_ADD_CMD (size = 28 bytes)
type PROTO_NC_HOLY_PROMISE_ADD_CMD_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	MemberInfo                       [26]byte             // offset: 2, size: 26 bytes
}

// PROTO_NC_HOLY_PROMISE_CENTRANSFER_RNG (size = 35 bytes)
type PROTO_NC_HOLY_PROMISE_CENTRANSFER_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Upper                            uint32               // offset: 7, size: 4 bytes
	Cen                              uint32               // offset: 11, size: 4 bytes
	Lower                            Name5                // offset: 15, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_CLIENT_GET_REMAIN_MONEY_CMD (size = 8 bytes)
type PROTO_NC_HOLY_PROMISE_CLIENT_GET_REMAIN_MONEY_CMD_DATA struct {
	NMoney                           uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_CHAR_ACK (size = 6 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_CHAR_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBError                          uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_CHAR_REQ (size = 4 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_CHAR_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_DOWN_ACK (size = 30 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_DOWN_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NDownCharNo                      uint32               // offset: 4, size: 4 bytes
	DownCharID                       Name5                // offset: 8, size: 20 bytes
	DBError                          uint16               // offset: 28, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_DOWN_REQ (size = 28 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_DOWN_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NDownCharNo                      uint32               // offset: 4, size: 4 bytes
	DownCharID                       Name5                // offset: 8, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_UP_ACK (size = 11 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_UP_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NUpCharNo                        uint32               // offset: 4, size: 4 bytes
	NDelCmdType                      uint8                // offset: 8, size: 1 bytes
	DBError                          uint16               // offset: 9, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_DEL_UP_REQ (size = 9 bytes)
type PROTO_NC_HOLY_PROMISE_DB_DEL_UP_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NUpCharNo                        uint32               // offset: 4, size: 4 bytes
	NDelCmdType                      uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK (size = 16 bytes)
type PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	NMoney                           uint64               // offset: 8, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_REQ (size = 6 bytes)
type PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_GET_MEMBER_ACK (size = 8 bytes)
type PROTO_NC_HOLY_PROMISE_DB_GET_MEMBER_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBError                          uint16               // offset: 4, size: 2 bytes
	NPart                            uint8                // offset: 6, size: 1 bytes
	NMemberCount                     uint8                // offset: 7, size: 1 bytes
	MemberInfo                       [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_HOLY_PROMISE_DB_GET_REMAIN_MONEY_CMD (size = 12 bytes)
type PROTO_NC_HOLY_PROMISE_DB_GET_REMAIN_MONEY_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NMoney                           uint64               // offset: 4, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_GET_UP_ACK (size = 48 bytes)
type PROTO_NC_HOLY_PROMISE_DB_GET_UP_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	DBError                          uint16               // offset: 4, size: 2 bytes
	UpMemberInfo                     [42]byte             // offset: 6, size: 42 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_REWARD_ACK (size = 9 bytes)
type PROTO_NC_HOLY_PROMISE_DB_REWARD_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NClass                           uint8                // offset: 4, size: 1 bytes
	NLevel                           uint8                // offset: 5, size: 1 bytes
	DBError                          uint16               // offset: 6, size: 2 bytes
	NNumOfItems                      uint8                // offset: 8, size: 1 bytes
	RewardList                       [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_HOLY_PROMISE_DB_REWARD_REQ (size = 6 bytes)
type PROTO_NC_HOLY_PROMISE_DB_REWARD_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NClass                           uint8                // offset: 4, size: 1 bytes
	NLevel                           uint8                // offset: 5, size: 1 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_SET_DATE_ACK (size = 14 bytes)
type PROTO_NC_HOLY_PROMISE_DB_SET_DATE_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	RejoinableDate                   uint32               // offset: 4, size: 4 bytes
	MemberAcceptableDate             uint32               // offset: 8, size: 4 bytes
	DBError                          uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_SET_DATE_REQ (size = 12 bytes)
type PROTO_NC_HOLY_PROMISE_DB_SET_DATE_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	RejoinableDate                   uint32               // offset: 4, size: 4 bytes
	MemberAcceptableDate             uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_SET_UP_ACK (size = 86 bytes)
type PROTO_NC_HOLY_PROMISE_DB_SET_UP_ACK_DATA struct {
	ReqInfo                          [42]byte             // offset: 0, size: 42 bytes
	UpInfo                           [42]byte             // offset: 42, size: 42 bytes
	DBError                          uint16               // offset: 84, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_SET_UP_REQ (size = 84 bytes)
type PROTO_NC_HOLY_PROMISE_DB_SET_UP_REQ_DATA struct {
	ReqInfo                          [42]byte             // offset: 0, size: 42 bytes
	UpInfo                           [42]byte             // offset: 42, size: 42 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK (size = 16 bytes)
type PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	NMoney                           uint64               // offset: 8, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_REQ (size = 6 bytes)
type PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_HOLY_PROMISE_DEL_CMD (size = 21 bytes)
type PROTO_NC_HOLY_PROMISE_DEL_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	NDelCmdType                      uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_HOLY_PROMISE_DEL_DOWN_ACK (size = 22 bytes)
type PROTO_NC_HOLY_PROMISE_DEL_DOWN_ACK_DATA struct {
	DownCharID                       Name5                // offset: 0, size: 20 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DEL_DOWN_REQ (size = 20 bytes)
type PROTO_NC_HOLY_PROMISE_DEL_DOWN_REQ_DATA struct {
	DownCharID                       Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_DEL_UP_ACK (size = 3 bytes)
type PROTO_NC_HOLY_PROMISE_DEL_UP_ACK_DATA struct {
	NDelCmdType                      uint8                // offset: 0, size: 1 bytes
	Error                            uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_HOLY_PROMISE_DEL_UP_REQ (size = 1 bytes)
type PROTO_NC_HOLY_PROMISE_DEL_UP_REQ_DATA struct {}

// PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK (size = 10 bytes)
type PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NMoney                           uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_LEVEL_CMD (size = 21 bytes)
type PROTO_NC_HOLY_PROMISE_LEVEL_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Level                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_HOLY_PROMISE_LIST_CMD (size = 29 bytes)
type PROTO_NC_HOLY_PROMISE_LIST_CMD_DATA struct {
	UpInfo                           [26]byte             // offset: 0, size: 26 bytes
	NPart                            uint8                // offset: 26, size: 1 bytes
	MemberCount                      uint16               // offset: 27, size: 2 bytes
	MemberInfo                       [0]byte              // offset: 29, flexible array member
}

// PROTO_NC_HOLY_PROMISE_LOGIN_CMD (size = 20 bytes)
type PROTO_NC_HOLY_PROMISE_LOGIN_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_LOGOUT_CMD (size = 20 bytes)
type PROTO_NC_HOLY_PROMISE_LOGOUT_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_MYUPPER_ACK (size = 10 bytes)
type PROTO_NC_HOLY_PROMISE_MYUPPER_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NUpCharNo                        uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_HOLY_PROMISE_MYUPPER_REQ (size = 6 bytes)
type PROTO_NC_HOLY_PROMISE_MYUPPER_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_HOLY_PROMISE_MY_UP_ZONE (size = 8 bytes)
type PROTO_NC_HOLY_PROMISE_MY_UP_ZONE_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NUpCharNo                        uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_HOLY_PROMISE_REWARD_ITEM_CMD (size = 1 bytes)
type PROTO_NC_HOLY_PROMISE_REWARD_ITEM_CMD_DATA struct {
	NNumOfItems                      uint8                // offset: 0, size: 1 bytes
	RewardList                       [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_HOLY_PROMISE_REWARD_MONEY_CMD (size = 28 bytes)
type PROTO_NC_HOLY_PROMISE_REWARD_MONEY_CMD_DATA struct {
	MemberCharID                     Name5                // offset: 0, size: 20 bytes
	NMoney                           uint64               // offset: 20, size: 8 bytes
}

// PROTO_NC_HOLY_PROMISE_SET_UP_ACK (size = 68 bytes)
type PROTO_NC_HOLY_PROMISE_SET_UP_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	UpInfo                           [26]byte             // offset: 2, size: 26 bytes
	ResetableDateTime                int32                // offset: 28, size: 4 bytes
	Tm_ResetableDateTime             [36]byte             // offset: 32, size: 36 bytes
}

// PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ACK (size = 41 bytes)
type PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ACK_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	UpCharID                         Name5                // offset: 20, size: 20 bytes
	BAccept                          uint8                // offset: 40, size: 1 bytes
}

// PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ING (size = 20 bytes)
type PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ING_DATA struct {
	UpCharID                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_REQ (size = 40 bytes)
type PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_REQ_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	UpCharID                         Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_SET_UP_REQ (size = 40 bytes)
type PROTO_NC_HOLY_PROMISE_SET_UP_REQ_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	UpCharID                         Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK (size = 10 bytes)
type PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NMoney                           uint64               // offset: 2, size: 8 bytes
}
