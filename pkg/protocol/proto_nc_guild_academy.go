package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHAR_GUILD_ACADEMY_CMD (size = 5 bytes)
type PROTO_NC_CHAR_GUILD_ACADEMY_CMD_DATA struct {
	NGuildAcademyNo                  uint32               // offset: 0, size: 4 bytes
	IsGuildAcademyMember             uint8                // offset: 4, size: 1 bytes
	GuildAcademy                     [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_CHAR_GUILD_ACADEMY_ZONE_CMD (size = 4 bytes)
type PROTO_NC_CHAR_GUILD_ACADEMY_ZONE_CMD_DATA struct {
	NGuildAcademyNo                  uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_CHAT_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_CHAT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_CHAT_CMD (size = 26 bytes)
type PROTO_NC_GUILD_ACADEMY_CHAT_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	Talker                           Name5                // offset: 4, size: 20 bytes
	Chat                             uint16               // offset: 24, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_CHAT_REQ (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_CHAT_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_ITEM_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_ITEM_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_ITEM_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_ITEM_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_MONEY_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_MONEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_MONEY_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_MONEY_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_ACK (size = 61 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	NCharNo                          uint32               // offset: 8, size: 4 bytes
	NGuildNo                         uint32               // offset: 12, size: 4 bytes
	NWantMoney                       uint64               // offset: 16, size: 8 bytes
	NWantItemCnt                     uint32               // offset: 24, size: 4 bytes
	WantItem                         [16]byte             // offset: 28, size: 16 bytes
	NSubAcademyMoney                 uint64               // offset: 44, size: 8 bytes
	NStorageMoney                    uint64               // offset: 52, size: 8 bytes
	NModifyItemCnt                   uint8                // offset: 60, size: 1 bytes
	ModifyItem                       [0]byte              // offset: 61, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_REQ (size = 15 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	NGuildNo                         uint32               // offset: 10, size: 4 bytes
	NLevel                           uint8                // offset: 14, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_ITEM_ACK (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_ITEM_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_ITEM_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_ITEM_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_MONEY_ACK (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_MONEY_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_MONEY_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_CLEAR_REWARD_MONEY_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_GRADUATE_ACK (size = 69 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_GRADUATE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	PartMask                         uint8                // offset: 4, size: 1 bytes
	GraduateData                     [64]byte             // offset: 5, size: 64 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_GRADUATE_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_GRADUATE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_RANKING_LIST_ACK (size = 53 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_RANKING_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	PartMask                         uint8                // offset: 4, size: 1 bytes
	RankData                         [48]byte             // offset: 5, size: 48 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_RANKING_LIST_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_RANKING_LIST_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NRoundTime                       uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_ITEM_ACK (size = 5 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_ITEM_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NumOfRewardGood                  uint8                // offset: 4, size: 1 bytes
	RewardItem                       [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_ITEM_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_ITEM_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_MONEY_ACK (size = 5 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_MONEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NumOfRewardGood                  uint8                // offset: 4, size: 1 bytes
	RewardMoney                      [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_MONEY_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GET_REWARD_MONEY_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GRADUATE_JOIN_ACK (size = 8 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GRADUATE_JOIN_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_GRADUATE_JOIN_REQ (size = 32 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_GRADUATE_JOIN_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	NGuildNo                         uint32               // offset: 6, size: 4 bytes
	SName                            Name5                // offset: 10, size: 20 bytes
	NClass                           uint8                // offset: 30, size: 1 bytes
	NGraduateLevel                   uint8                // offset: 31, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_LEVEL_UP_ACK (size = 16 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_LEVEL_UP_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NGuildNo                         uint32               // offset: 4, size: 4 bytes
	NAddAcademyPoint                 uint32               // offset: 8, size: 4 bytes
	DAddBuffEndTime                  int32                // offset: 12, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_LEVEL_UP_REQ (size = 14 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_LEVEL_UP_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NAddAcademyPoint                 uint32               // offset: 6, size: 4 bytes
	DAddBuffEndTime                  int32                // offset: 10, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_ACK (size = 9 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_ACK_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
	PartMask                         uint8                // offset: 6, size: 1 bytes
	NumOfAcademyMember               uint16               // offset: 7, size: 2 bytes
	AcademyMemberList                [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_GUILD_JOIN_ACK (size = 33 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_GUILD_JOIN_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	SInviteCharID                    Name5                // offset: 6, size: 20 bytes
	NTargetCharNo                    uint32               // offset: 26, size: 4 bytes
	NGrade                           uint8                // offset: 30, size: 1 bytes
	Error                            uint16               // offset: 31, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_GUILD_JOIN_REQ (size = 51 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_GUILD_JOIN_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	SInviteCharID                    Name5                // offset: 6, size: 20 bytes
	NTargetCharNo                    uint32               // offset: 26, size: 4 bytes
	STargetCharID                    Name5                // offset: 30, size: 20 bytes
	NGrade                           uint8                // offset: 50, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_JOIN_ACK (size = 18 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_JOIN_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	NGrade                           uint8                // offset: 10, size: 1 bytes
	DJoinInfo                        uint32               // offset: 11, size: 4 bytes
	NJoinLevel                       uint8                // offset: 15, size: 1 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_JOIN_REQ (size = 32 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_JOIN_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	SCharID                          Name5                // offset: 10, size: 20 bytes
	NJoinLevel                       uint8                // offset: 30, size: 1 bytes
	NGrade                           uint8                // offset: 31, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_LEAVE_ACK (size = 12 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_LEAVE_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	Error                            uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_LEAVE_REQ (size = 30 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_LEAVE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	SCharID                          Name5                // offset: 10, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_REQ (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_VANISH_ACK (size = 32 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_VANISH_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
	Error                            uint16               // offset: 30, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_MEMBER_VANISH_REQ (size = 30 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_MEMBER_VANISH_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_NOTIFY_ACK (size = 14 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_NOTIFY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NNo                              uint32               // offset: 4, size: 4 bytes
	DNotifyDate                      int32                // offset: 8, size: 4 bytes
	NNotifyLength                    uint16               // offset: 12, size: 2 bytes
	SNotify                          [0]byte              // offset: 14, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_NOTIFY_REQ (size = 8 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_NOTIFY_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NNotifyLength                    uint16               // offset: 6, size: 2 bytes
	SNotify                          [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_RANK_BALANCE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_RANK_BALANCE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_RANK_BALANCE_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_RANK_BALANCE_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_DB_REWARD_WANT_CMD (size = 22 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_REWARD_WANT_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NGuildNo                         uint32               // offset: 4, size: 4 bytes
	NWant_Cen                        uint64               // offset: 8, size: 8 bytes
	NWant_ItemID                     uint16               // offset: 16, size: 2 bytes
	NWant_ItemLot                    uint32               // offset: 18, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_ACK (size = 36 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SName                            Name5                // offset: 8, size: 20 bytes
	NGuildNo                         uint32               // offset: 28, size: 4 bytes
	NSetCharNo                       uint32               // offset: 32, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_BY_LEAVE_ACK (size = 54 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_BY_LEAVE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NOldAcademyMasterCharNo          uint32               // offset: 6, size: 4 bytes
	SOldAcademyMasterName            Name5                // offset: 10, size: 20 bytes
	NNewAcademyMasterCharNo          uint32               // offset: 30, size: 4 bytes
	SNewAcademyMasterName            Name5                // offset: 34, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_BY_LEAVE_REQ (size = 52 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_BY_LEAVE_REQ_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NOldAcademyMasterCharNo          uint32               // offset: 4, size: 4 bytes
	SOldAcademyMasterName            Name5                // offset: 8, size: 20 bytes
	NNewAcademyMasterCharNo          uint32               // offset: 28, size: 4 bytes
	SNewAcademyMasterName            Name5                // offset: 32, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_REQ (size = 34 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_MASTER_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	SName                            Name5                // offset: 6, size: 20 bytes
	NGuildNo                         uint32               // offset: 26, size: 4 bytes
	NSetCharNo                       uint32               // offset: 30, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_ITEM_ACK (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_ITEM_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_ITEM_REQ (size = 7 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_ITEM_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NumOfRewardGood                  uint8                // offset: 6, size: 1 bytes
	RewardItem                       [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_MONEY_ACK (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_MONEY_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_MONEY_REQ (size = 7 bytes)
type PROTO_NC_GUILD_ACADEMY_DB_SET_REWARD_MONEY_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NumOfRewardGood                  uint8                // offset: 6, size: 1 bytes
	RewardMoney                      [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_DELETE_CMD (size = 16 bytes)
type PROTO_NC_GUILD_ACADEMY_DELETE_CMD_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
}

// PROTO_NC_GUILD_ACADEMY_DISMISS_CMD (size = 57 bytes)
type PROTO_NC_GUILD_ACADEMY_DISMISS_CMD_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
	NDismissStatus                   uint8                // offset: 16, size: 1 bytes
	DDismissDate                     int32                // offset: 17, size: 4 bytes
	Tm_dDismissDate                  [36]byte             // offset: 21, size: 36 bytes
}

// PROTO_NC_GUILD_ACADEMY_GET_GRADUATE_ACK (size = 67 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_GRADUATE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	GraduateData                     [64]byte             // offset: 3, size: 64 bytes
}

// PROTO_NC_GUILD_ACADEMY_GET_GRADUATE_REQ (size = 16 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_GRADUATE_REQ_DATA struct {
	SGuildName                       [16]byte             // offset: 0, size: 16 bytes
}

// PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_ACK (size = 51 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	RankData                         [48]byte             // offset: 3, size: 48 bytes
}

// PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_REQ (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_REQ_DATA struct {
	NRoundTime                       uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_GET_REWARD_ITEM_ACK (size = 3 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_REWARD_ITEM_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NumOfRewardGood                  uint8                // offset: 2, size: 1 bytes
	RewardItem                       [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_GET_REWARD_ITEM_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_REWARD_ITEM_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_GET_REWARD_MONEY_ACK (size = 3 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_REWARD_MONEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NumOfRewardGood                  uint8                // offset: 2, size: 1 bytes
	RewardMoney                      [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_GET_REWARD_MONEY_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_GET_REWARD_MONEY_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_HISTORY_DB_LIST_ACK (size = 10 bytes)
type PROTO_NC_GUILD_ACADEMY_HISTORY_DB_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	TotalHistoryCount                uint16               // offset: 4, size: 2 bytes
	HistoryCount                     uint16               // offset: 6, size: 2 bytes
	HistoryIndex                     uint16               // offset: 8, size: 2 bytes
	HistoryList                      [0]byte              // offset: 10, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_HISTORY_DB_LIST_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_HISTORY_DB_LIST_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildAcademyNo                  uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_HISTORY_LIST_ACK (size = 8 bytes)
type PROTO_NC_GUILD_ACADEMY_HISTORY_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	TotalHistoryCount                uint16               // offset: 2, size: 2 bytes
	HistoryCount                     uint16               // offset: 4, size: 2 bytes
	HistoryIndex                     uint16               // offset: 6, size: 2 bytes
	HistoryList                      [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_LIST_ACK (size = 11 bytes)
type PROTO_NC_GUILD_ACADEMY_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	TotalNumOfGuildAcademy           uint32               // offset: 3, size: 4 bytes
	NumOfGuildAcademy                uint32               // offset: 7, size: 4 bytes
	GuildAcademyList                 [0]byte              // offset: 11, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_MASTER_TELEPORT_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_MASTER_TELEPORT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MASTER_TELEPORT_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MASTER_TELEPORT_REQ_DATA struct {
	STargetCharID                    Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_CMD (size = 40 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_CMD_DATA struct {
	SChatBanCancelSourceCharID       Name5                // offset: 0, size: 20 bytes
	SChatBanCancelTargetCharID       Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_REQ_DATA struct {
	SChatBanCancelTargetCharID       Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CMD (size = 40 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CMD_DATA struct {
	SChatBanSourceCharID             Name5                // offset: 0, size: 20 bytes
	SChatBanTargetCharID             Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_REQ_DATA struct {
	SChatBanTargetCharID             Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_CLASS_CMD (size = 21 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_CLASS_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	ClassID                          uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_ACK (size = 18 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NBanCharNo                       uint32               // offset: 8, size: 4 bytes
	NRestTime                        uint32               // offset: 12, size: 4 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_CANCEL_ACK (size = 14 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_CANCEL_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NBanCancelCharNo                 uint32               // offset: 8, size: 4 bytes
	Error                            uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_CANCEL_REQ (size = 12 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_CANCEL_REQ_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NBanCancelCharNo                 uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_REQ (size = 16 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_DB_CHAT_BAN_REQ_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NBanCharNo                       uint32               // offset: 8, size: 4 bytes
	NRestTime                        uint32               // offset: 12, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_DB_SAVE_CHAT_BAN_TIME_CMD (size = 12 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_DB_SAVE_CHAT_BAN_TIME_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NRestTime                        uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_INVITE_ACK (size = 22 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_INVITE_ACK_DATA struct {
	SInviteTargetCharID              Name5                // offset: 0, size: 20 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_INVITE_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_INVITE_REQ_DATA struct {
	SInviteTargetCharID              Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_ACK (size = 37 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_ACK_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	SInviteSourceCharID              Name5                // offset: 16, size: 20 bytes
	JoinAck                          uint8                // offset: 36, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_CMD (size = 130 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_CMD_DATA struct {
	SInviteSourceCharID              Name5                // offset: 0, size: 20 bytes
	Member                           [110]byte            // offset: 20, size: 110 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_ERR (size = 18 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_ERR_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_REQ (size = 36 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_REQ_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	SInviteSourceCharID              Name5                // offset: 16, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_INTRO_CMD (size = 22 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_INTRO_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	NMemberIntroLength               uint16               // offset: 20, size: 2 bytes
	SMemberIntro                     [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK (size = 18 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD (size = 109 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD_DATA struct {
	Member                           [109]byte            // offset: 0, size: 109 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_REQ (size = 16 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_REQ_DATA struct {
	SGuildName                       [16]byte             // offset: 0, size: 16 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LEVEL_CMD (size = 21 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LEVEL_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Level                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK_DATA struct {
	TotalAcademyMemberCount          uint16               // offset: 0, size: 2 bytes
	AcademyMemberIndex               uint16               // offset: 2, size: 2 bytes
	AcademyMemberCount               uint16               // offset: 4, size: 2 bytes
	AcademyMemberList                [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LOGOFF_CMD (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LOGOFF_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_LOGON_CMD (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_LOGON_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_MAP_CMD (size = 32 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_MAP_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	SMap                             Name3                // offset: 20, size: 12 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_PARTY_CMD (size = 21 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_PARTY_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Party                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK (size = 22 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_REQ_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_MY_ACADEMY_RANK_INFO_ACK (size = 14 bytes)
type PROTO_NC_GUILD_ACADEMY_MY_ACADEMY_RANK_INFO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NAcademyPoint                    uint32               // offset: 2, size: 4 bytes
	NAcademyRank                     uint32               // offset: 6, size: 4 bytes
	DAcademyBuffUntilTime            int32                // offset: 10, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_MY_ACADEMY_RANK_INFO_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_MY_ACADEMY_RANK_INFO_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	GuildAcademyInfo                 [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD (size = 62 bytes)
type PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD_DATA struct {
	DNotifyDate                      int32                // offset: 0, size: 4 bytes
	Tm_dNotifyDate                   [36]byte             // offset: 4, size: 36 bytes
	SNotifyCharID                    Name5                // offset: 40, size: 20 bytes
	NNotifyLength                    uint16               // offset: 60, size: 2 bytes
	SNotify                          [0]byte              // offset: 62, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_NOTIFY_REQ (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_NOTIFY_REQ_DATA struct {
	NNotifyLength                    uint16               // offset: 0, size: 2 bytes
	SNotify                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_GRADE_INFO_CMD (size = 9 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_GRADE_INFO_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NGuildNo                         uint32               // offset: 4, size: 4 bytes
	Grade                            uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_ITEM_INFO_ZONE_RNG (size = 176 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_ITEM_INFO_ZONE_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	NGldregnum                       uint32               // offset: 7, size: 4 bytes
	Who                              uint8                // offset: 11, size: 1 bytes
	Charid                           Name5                // offset: 12, size: 20 bytes
	TypeField                        uint8                // offset: 32, size: 1 bytes
	From                             [16]byte             // offset: 33, size: 16 bytes
	To                               [16]byte             // offset: 49, size: 16 bytes
	ChangeIteminfo                   [111]byte            // offset: 65, size: 111 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_REWARD_ZONE_RNG (size = 48 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_REWARD_ZONE_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	NGuildNo                         uint32               // offset: 7, size: 4 bytes
	SCharID                          Name5                // offset: 11, size: 20 bytes
	Reward_Cen                       uint64               // offset: 31, size: 8 bytes
	Storage_Cen                      uint64               // offset: 39, size: 8 bytes
	Reward_ItemNum                   uint8                // offset: 47, size: 1 bytes
	Reward_ItemInfo                  [0]byte              // offset: 48, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD_DATA struct {
	NCloseType                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_CMD (size = 41 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_CMD_DATA struct {
	Who                              uint8                // offset: 0, size: 1 bytes
	TypeField                        uint8                // offset: 1, size: 1 bytes
	Charid                           Name5                // offset: 2, size: 20 bytes
	Item                             uint16               // offset: 22, size: 2 bytes
	Slot                             uint8                // offset: 24, size: 1 bytes
	Lot                              [0]byte              // offset: 25, size: 0 bytes
	Cen                              uint64               // offset: 25, size: 8 bytes
	Storage_cen                      uint64               // offset: 33, size: 8 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_RNG (size = 53 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Gldregnum                        uint32               // offset: 7, size: 4 bytes
	Why                              uint8                // offset: 11, size: 1 bytes
	Who                              uint8                // offset: 12, size: 1 bytes
	TypeField                        uint8                // offset: 13, size: 1 bytes
	Slot                             uint8                // offset: 14, size: 1 bytes
	Charid                           Name5                // offset: 15, size: 20 bytes
	Item                             uint16               // offset: 35, size: 2 bytes
	Lot                              [0]byte              // offset: 37, size: 0 bytes
	Cen                              uint64               // offset: 37, size: 8 bytes
	Storage_cen                      uint64               // offset: 45, size: 8 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK (size = 11 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Cen                              uint64               // offset: 2, size: 8 bytes
	Itemcounter                      uint8                // offset: 10, size: 1 bytes
	Itemarray                        [0]byte              // offset: 11, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ_DATA struct {
	NOpenType                        uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_CMD (size = 41 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_CMD_DATA struct {
	Who                              uint8                // offset: 0, size: 1 bytes
	TypeField                        uint8                // offset: 1, size: 1 bytes
	Charid                           Name5                // offset: 2, size: 20 bytes
	Item                             uint16               // offset: 22, size: 2 bytes
	Slot                             uint8                // offset: 24, size: 1 bytes
	Lot                              [0]byte              // offset: 25, size: 0 bytes
	Cen                              uint64               // offset: 25, size: 8 bytes
	Storage_cen                      uint64               // offset: 33, size: 8 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_ACK (size = 9 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	Grade                            uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_REQ (size = 6 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_RNG (size = 53 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Gldregnum                        uint32               // offset: 7, size: 4 bytes
	Why                              uint8                // offset: 11, size: 1 bytes
	Who                              uint8                // offset: 12, size: 1 bytes
	TypeField                        uint8                // offset: 13, size: 1 bytes
	Slot                             uint8                // offset: 14, size: 1 bytes
	Charid                           Name5                // offset: 15, size: 20 bytes
	Item                             uint16               // offset: 35, size: 2 bytes
	Lot                              [0]byte              // offset: 37, size: 0 bytes
	Cen                              uint64               // offset: 37, size: 8 bytes
	Storage_cen                      uint64               // offset: 45, size: 8 bytes
}

// PROTO_NC_GUILD_ACADEMY_REWARD_WANT_CMD (size = 14 bytes)
type PROTO_NC_GUILD_ACADEMY_REWARD_WANT_CMD_DATA struct {
	NWant_Cen                        uint64               // offset: 0, size: 8 bytes
	NWant_ItemID                     uint16               // offset: 8, size: 2 bytes
	NWant_ItemLot                    uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_MASTER_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_MASTER_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_MASTER_BY_LEAVE_CMD (size = 40 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_MASTER_BY_LEAVE_CMD_DATA struct {
	SOldAcademyMasterName            Name5                // offset: 0, size: 20 bytes
	SNewAcademyMasterName            Name5                // offset: 20, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_MASTER_CMD (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_MASTER_CMD_DATA struct {
	SName                            Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_MASTER_REQ (size = 20 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_MASTER_REQ_DATA struct {
	SName                            Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_REWARD_ITEM_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_REWARD_ITEM_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_REWARD_ITEM_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_REWARD_ITEM_REQ_DATA struct {
	NumOfRewardGood                  uint8                // offset: 0, size: 1 bytes
	RewardItem                       [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_SET_REWARD_MONEY_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_REWARD_MONEY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_SET_REWARD_MONEY_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_SET_REWARD_MONEY_REQ_DATA struct {
	NumOfRewardGood                  uint8                // offset: 0, size: 1 bytes
	RewardMoney                      [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_START_DB_ALL_ACK (size = 5 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_ALL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NumOfGuildAcademy                uint16               // offset: 3, size: 2 bytes
	GuildAcademyDB                   [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_START_DB_ALL_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_ALL_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_START_DB_GET_RANKING_LIST_ACK (size = 11 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_GET_RANKING_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	RankData                         uint64               // offset: 3, size: 8 bytes
}

// PROTO_NC_GUILD_ACADEMY_START_DB_GET_RANKING_LIST_REQ (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_GET_RANKING_LIST_REQ_DATA struct {
	NRoundTime                       uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_START_DB_RANK_BALANCE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_RANK_BALANCE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_ACADEMY_START_DB_RANK_BALANCE_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ACADEMY_START_DB_RANK_BALANCE_REQ_DATA struct {}

// PROTO_NC_GUILD_ACADEMY_WAR_END_CMD (size = 5 bytes)
type PROTO_NC_GUILD_ACADEMY_WAR_END_CMD_DATA struct {
	EnemyGuildNo                     uint32               // offset: 0, size: 4 bytes
	WarStatus                        uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_WAR_INFO_LIST_CMD (size = 2 bytes)
type PROTO_NC_GUILD_ACADEMY_WAR_INFO_LIST_CMD_DATA struct {
	NumOfWarInfo                     uint16               // offset: 0, size: 2 bytes
	WarInfoList                      [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_WAR_START_CMD (size = 5 bytes)
type PROTO_NC_GUILD_ACADEMY_WAR_START_CMD_DATA struct {
	EnemyGuildNo                     uint32               // offset: 0, size: 4 bytes
	WarStatus                        uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_ACADEMY_MASTER_BUFF_CMD (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_ACADEMY_MASTER_BUFF_CMD_DATA struct {
	NumOfMaster                      uint32               // offset: 0, size: 4 bytes
	Master                           [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_BUFF_CMD (size = 36 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_BUFF_CMD_DATA struct {
	BuffName                         [32]byte             // offset: 0, size: 32 bytes
	NumOfGuild                       uint32               // offset: 32, size: 4 bytes
	Guild                            [0]byte              // offset: 36, flexible array member
}

// PROTO_NC_GUILD_ACADEMY_ZONE_MASTER_TELEPORT_CMD (size = 24 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_MASTER_TELEPORT_CMD_DATA struct {
	NMasterCharNo                    uint32               // offset: 0, size: 4 bytes
	STargetCharID                    Name5                // offset: 4, size: 20 bytes
}

// PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_GUILD_JOIN_CMD (size = 8 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_GUILD_JOIN_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_JOIN_CMD (size = 8 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_JOIN_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_LEAVE_CMD (size = 4 bytes)
type PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_LEAVE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ (size = 20 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Guildregnum                      uint32               // offset: 8, size: 4 bytes
	Cen                              uint64               // offset: 12, size: 8 bytes
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK (size = 16 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK_DATA struct {
	Zonepackheader                   NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Cen                              uint64               // offset: 6, size: 8 bytes
	NOpenType                        uint8                // offset: 14, size: 1 bytes
	Itemcounter                      uint8                // offset: 15, size: 1 bytes
	Itemarray                        [0]byte              // offset: 16, flexible array member
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ (size = 11 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Guildregnum                      uint32               // offset: 6, size: 4 bytes
	NOpenType                        uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ (size = 20 bytes)
type PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Guildregnum                      uint32               // offset: 8, size: 4 bytes
	Cen                              uint64               // offset: 12, size: 8 bytes
}
