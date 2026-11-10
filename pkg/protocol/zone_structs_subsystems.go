package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_subsystems.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_ANNOUNCE_Z2W_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ANNOUNCE_Z2W_CMD_DATA)

// PROTO_NC_BOOTH_BUYREFRESH_CMD (sizeof = 3 bytes)
type PROTO_NC_BOOTH_BUYREFRESH_CMD_DATA struct {
	Boothslot                        uint8                // offset: 0, size: 1 bytes
	Restlot                          uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_BUYREFRESH_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BOOTH_ENTRY_BUY_ACK (sizeof = 5 bytes)
type PROTO_NC_BOOTH_ENTRY_BUY_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Boothhandle                      uint16               // offset: 2, size: 2 bytes
	Numofitem                        uint8                // offset: 4, size: 1 bytes
	Items                            [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_BOOTH_ENTRY_BUY_ACK_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BOOTH_ENTRY_REQ (sizeof = 2 bytes)
type PROTO_NC_BOOTH_ENTRY_REQ_DATA struct {
	Booth                            uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_ENTRY_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BOOTH_ENTRY_SELL_ACK (sizeof = 5 bytes)
type PROTO_NC_BOOTH_ENTRY_SELL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Boothhandle                      uint16               // offset: 2, size: 2 bytes
	Numofitem                        uint8                // offset: 4, size: 1 bytes
	Items                            [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_BOOTH_ENTRY_SELL_ACK_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BOOTH_ITEMTRADE_REQ (sizeof = 5 bytes)
type PROTO_NC_BOOTH_ITEMTRADE_REQ_DATA struct {
	Booth                            uint16               // offset: 0, size: 2 bytes
	Slot                             uint8                // offset: 2, size: 1 bytes
	Itemnum                          uint16               // offset: 3, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_ITEMTRADE_REQ_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BOOTH_OPEN_REQ (sizeof = 31 bytes)
type PROTO_NC_BOOTH_OPEN_REQ_DATA struct {
	Signboard                        [30]byte             // offset: 0, size: 30 bytes
	Flag                             uint8                // offset: 30, size: 1 bytes
	Items                            [0]byte              // offset: 31, flexible array member
}

func (s *PROTO_NC_BOOTH_OPEN_REQ_DATA) FixedWireSize() int { return 31 }

// PROTO_NC_BOOTH_REFRESH_REQ (sizeof = 2 bytes)
type PROTO_NC_BOOTH_REFRESH_REQ_DATA struct {
	Booth                            uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_REFRESH_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD (sizeof = 2 bytes)
type PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD_DATA struct {
	NClosedBoothOwnerHandle          uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK (sizeof = 12 bytes)
type PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NTargetBoothOwnerHandle          uint16               // offset: 2, size: 2 bytes
	BoothPosition                    SHINE_XY_TYPE        // offset: 4, size: 8 bytes
}

func (s *PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ (sizeof = 6 bytes)
type PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ_DATA struct {
	NBoothItemHandle                 uint32               // offset: 0, size: 4 bytes
	NTargetBoothOwnerHandle          uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK (sizeof = 3 bytes)
type PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NNumOfItems                      uint8                // offset: 2, size: 1 bytes
	ItemList                         [0]byte              // offset: 3, flexible array member
}

func (s *PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ (sizeof = 6 bytes)
type PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ_DATA struct {
	Category                         uint32               // offset: 0, size: 4 bytes
	BReqOnlyMyClass                  uint8                // offset: 4, size: 1 bytes
	BIsOpeningReq                    uint8                // offset: 5, size: 1 bytes
}

func (s *PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BOOTH_SOMEONECLOSE_CMD (sizeof = 45 bytes)
type PROTO_NC_BOOTH_SOMEONECLOSE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Shape                            [43]byte             // offset: 2, size: 43 bytes
}

func (s *PROTO_NC_BOOTH_SOMEONECLOSE_CMD_DATA) FixedWireSize() int { return 45 }

// PROTO_NC_BOOTH_SOMEONEINTERIORSTART_CMD (sizeof = 32 bytes)
type PROTO_NC_BOOTH_SOMEONEINTERIORSTART_CMD_DATA struct {
	Boothowner                       uint16               // offset: 0, size: 2 bytes
	Signboard                        [30]byte             // offset: 2, size: 30 bytes
}

func (s *PROTO_NC_BOOTH_SOMEONEINTERIORSTART_CMD_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_BOOTH_SOMEONEOPEN_CMD (sizeof = 45 bytes)
type PROTO_NC_BOOTH_SOMEONEOPEN_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Tent                             [12]byte             // offset: 2, size: 12 bytes
	Issell                           uint8                // offset: 14, size: 1 bytes
	Signboard                        [30]byte             // offset: 15, size: 30 bytes
}

func (s *PROTO_NC_BOOTH_SOMEONEOPEN_CMD_DATA) FixedWireSize() int { return 45 }

// PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK_DATA)

// PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_ACK (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_ACK_DATA)

// PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_REQ (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_REQ_DATA)

// PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ_DATA)

// PROTO_NC_COLLECT_CARDOPEN_CMD (sizeof = 106 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_CARDOPEN_CMD_DATA)

// PROTO_NC_COLLECT_CARDREGIST_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_CARDREGIST_ACK_DATA)

// PROTO_NC_COLLECT_CARDREGIST_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_CARDREGIST_REQ_DATA)

// PROTO_NC_COLLECT_PEEPING_COLLECT_ACK (sizeof = 28 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_COLLECT_ACK_DATA)

// PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK (sizeof = 28 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK_DATA)

// PROTO_NC_COLLECT_PEEPING_DB_COLLECT_ACK (sizeof = 34 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_DB_COLLECT_ACK_DATA)

// PROTO_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK (sizeof = 34 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK_DATA)

// PROTO_NC_COLLECT_PEEPING_DB_REQ (sizeof = 26 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_DB_REQ_DATA)

// PROTO_NC_COLLECT_PEEPING_REQ (sizeof = 20 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_PEEPING_REQ_DATA)

// PROTO_NC_COLLECT_REGIST_ACK (sizeof = 18 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REGIST_ACK_DATA)

// PROTO_NC_COLLECT_REGIST_REQ (sizeof = 26 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REGIST_REQ_DATA)

// PROTO_NC_COLLECT_REWARD_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REWARD_ACK_DATA)

// PROTO_NC_COLLECT_REWARD_DB_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REWARD_DB_ACK_DATA)

// PROTO_NC_COLLECT_REWARD_DB_REQ (sizeof = 119 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REWARD_DB_REQ_DATA)

// PROTO_NC_COLLECT_REWARD_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_COLLECT_REWARD_REQ_DATA)

// PROTO_NC_CT_ADD_FRIEND_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_CT_ADD_FRIEND_CMD_DATA)

// PROTO_NC_CT_CHARTTING_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_CT_CHARTTING_CMD_DATA)

// PROTO_NC_CT_DB_SET_CMD (sizeof = 8 bytes)
type PROTO_NC_CT_DB_SET_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Info                             uint16               // offset: 6, size: 2 bytes
}

func (s *PROTO_NC_CT_DB_SET_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_CT_LUASCRIPT_SET_WORLD_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_CT_LUASCRIPT_SET_WORLD_CMD_DATA)

// PROTO_NC_CT_LUASCRIPT_SET_ZONE_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_CT_LUASCRIPT_SET_ZONE_CMD_DATA)

// PROTO_NC_CT_SET_CMD (sizeof = 2 bytes)
type PROTO_NC_CT_SET_CMD_DATA struct {
	Info                             uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_CT_SET_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_CT_SET_CURRENT_ACK (sizeof = 6 bytes)
type PROTO_NC_CT_SET_CURRENT_ACK_DATA struct {
	TypeField                        uint8                // offset: 0, size: 1 bytes
	ElementNo                        uint8                // offset: 1, size: 1 bytes
	MobID                            uint16               // offset: 2, size: 2 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_CT_SET_CURRENT_ACK_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_CT_SET_CURRENT_DB_CMD (sizeof = 14 bytes)
type PROTO_NC_CT_SET_CURRENT_DB_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	CurrentTitle                     uint8                // offset: 6, size: 1 bytes
	CurrentTitleElement              uint8                // offset: 7, size: 1 bytes
	CurrentTitleMobID                uint16               // offset: 8, size: 2 bytes
	UseTitle                         uint16               // offset: 10, size: 2 bytes
	UseAbleTitle                     uint16               // offset: 12, size: 2 bytes
}

func (s *PROTO_NC_CT_SET_CURRENT_DB_CMD_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_CT_SET_CURRENT_REQ (sizeof = 4 bytes)
type PROTO_NC_CT_SET_CURRENT_REQ_DATA struct {
	TypeField                        uint8                // offset: 0, size: 1 bytes
	ElementNo                        uint8                // offset: 1, size: 1 bytes
	MobID                            uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_CT_SET_CURRENT_REQ_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_CT_SET_SOMEONECHANGE_CMD (sizeof = 6 bytes)
type PROTO_NC_CT_SET_SOMEONECHANGE_CMD_DATA struct {
	Hnd                              uint16               // offset: 0, size: 2 bytes
	TypeField                        uint8                // offset: 2, size: 1 bytes
	ElementNo                        uint8                // offset: 3, size: 1 bytes
	MobID                            uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_CT_SET_SOMEONECHANGE_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_EVENT_ADD_UPDATE_EVENT_CMD (sizeof = 1363 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_EVENT_ADD_UPDATE_EVENT_CMD_DATA)

// PROTO_NC_EVENT_DEL_EVENT_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_EVENT_DEL_EVENT_CMD_DATA)

// PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_ZONE_CMD (sizeof = 20 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_ZONE_CMD_DATA)

// PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK (sizeof = 7 bytes)
type PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Friendpoint                      uint16               // offset: 2, size: 2 bytes
	Itemid                           uint16               // offset: 4, size: 2 bytes
	Numofitem                        uint8                // offset: 6, size: 1 bytes
}

func (s *PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_HOLY_PROMISE_CENTRANSFER_RNG (sizeof = 35 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_CENTRANSFER_RNG_DATA)

// PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK (sizeof = 16 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK_DATA)

// PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_REQ_DATA)

// PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK (sizeof = 16 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK_DATA)

// PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_REQ_DATA)

// PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK_DATA)

// PROTO_NC_HOLY_PROMISE_MYUPPER_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_MYUPPER_ACK_DATA)

// PROTO_NC_HOLY_PROMISE_MY_UP_ZONE (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_MY_UP_ZONE_DATA)

// PROTO_NC_HOLY_PROMISE_REWARD_MONEY_CMD (sizeof = 28 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_REWARD_MONEY_CMD_DATA)

// PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK_DATA)

// PROTO_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD (sizeof = 9 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD_DATA)

// PROTO_NC_INSTANCE_DUNGEON_EMPTY_DUNGEON_CMD (sizeof = 5 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_EMPTY_DUNGEON_CMD_DATA)

// PROTO_NC_INSTANCE_DUNGEON_FIND_RNG (sizeof = 70 bytes)
type PROTO_NC_INSTANCE_DUNGEON_FIND_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	InstanceCategory                 uint32               // offset: 7, size: 4 bytes
	IDRegisterNumber                 uint32               // offset: 11, size: 4 bytes
	Argument                         [20]byte             // offset: 15, size: 20 bytes
	ServerMapName                    [33]byte             // offset: 35, size: 33 bytes
	Error                            uint16               // offset: 68, size: 2 bytes
}

func (s *PROTO_NC_INSTANCE_DUNGEON_FIND_RNG_DATA) FixedWireSize() int { return 70 }

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_ACK (sizeof = 13 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_ACK_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ (sizeof = 35 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_REQ (sizeof = 36 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_REQ_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_ACK (sizeof = 25 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_ACK_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ (sizeof = 39 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_REQ (sizeof = 40 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_REQ_DATA)

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK (sizeof = 6 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK_DATA struct {
	NNPCHandle                       uint16               // offset: 0, size: 2 bytes
	ELevelType                       uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_REQ (sizeof = 15 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_REQ_DATA struct {
	NNPCHandle                       uint16               // offset: 0, size: 2 bytes
	SMapIndex                        Name3                // offset: 2, size: 12 bytes
	NMenuCount                       uint8                // offset: 14, size: 1 bytes
	ELevelType                       [0]byte              // offset: 15, flexible array member
}

func (s *PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_REQ_DATA) FixedWireSize() int { return 15 }

// PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD_DATA)

// PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK (sizeof = 2 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK_DATA struct {
	NMapCount                        uint16               // offset: 0, size: 2 bytes
	SIndunInfos                      [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD (sizeof = 5 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD_DATA struct {
	NPlayTimeHour                    uint8                // offset: 0, size: 1 bytes
	NPlayTimeMinute                  uint8                // offset: 1, size: 1 bytes
	NPlayTimeSecond                  uint8                // offset: 2, size: 1 bytes
	NRankCount                       uint16               // offset: 3, size: 2 bytes
	Ranks                            [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_DB_CMD (sizeof = 19 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_DB_CMD_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
	NPlayTime                        uint32               // offset: 13, size: 4 bytes
	NRankCount                       uint16               // offset: 17, size: 2 bytes
	Ranks                            [0]byte              // offset: 19, flexible array member
}

func (s *PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_DB_CMD_DATA) FixedWireSize() int { return 19 }

// PROTO_NC_INSTANCE_DUNGEON_RESET_COUNTDOWN_CMD (sizeof = 5 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_INSTANCE_DUNGEON_RESET_COUNTDOWN_CMD_DATA)

// PROTO_NC_KQ_ENTRYRESPONCE_ACK (sizeof = 5 bytes)
type PROTO_NC_KQ_ENTRYRESPONCE_ACK_DATA struct {
	Reply                            uint8                // offset: 0, size: 1 bytes
	EncHandle                        uint32               // offset: 1, size: 4 bytes
}

func (s *PROTO_NC_KQ_ENTRYRESPONCE_ACK_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_KQ_JOINER (sizeof = 5 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_JOINER)

// PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD (sizeof = 52 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD_DATA)

// PROTO_NC_KQ_MOBKILLNUMBER_CMD (sizeof = 4 bytes)
type PROTO_NC_KQ_MOBKILLNUMBER_CMD_DATA struct {
	CurrentMobKill                   uint16               // offset: 0, size: 2 bytes
	DemandMobKill                    uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_KQ_MOBKILLNUMBER_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_KQ_PLAYER_DISJOIN_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_PLAYER_DISJOIN_CMD_DATA)

// PROTO_NC_KQ_RESTDEADNUM_CMD (sizeof = 1 bytes)
type PROTO_NC_KQ_RESTDEADNUM_CMD_DATA struct {
	Number                           uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_KQ_RESTDEADNUM_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_KQ_REWARDFAIL_ACK (sizeof = 10 bytes)
type PROTO_NC_KQ_REWARDFAIL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Error                            uint16               // offset: 8, size: 2 bytes
}

func (s *PROTO_NC_KQ_REWARDFAIL_ACK_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_KQ_REWARDSUC_ACK (sizeof = 8 bytes)
type PROTO_NC_KQ_REWARDSUC_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

func (s *PROTO_NC_KQ_REWARDSUC_ACK_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_KQ_REWARD_REQ (sizeof = 35 bytes)
type PROTO_NC_KQ_REWARD_REQ_DATA struct {
	Fame                             uint32               // offset: 0, size: 4 bytes
	Cen                              uint64               // offset: 4, size: 8 bytes
	Itmlst                           [23]byte             // offset: 12, size: 23 bytes
}

func (s *PROTO_NC_KQ_REWARD_REQ_DATA) FixedWireSize() int { return 35 }

// PROTO_NC_KQ_SCORE_BOARD_INFO_CMD (sizeof = 6 bytes)
type PROTO_NC_KQ_SCORE_BOARD_INFO_CMD_DATA struct {
	BUseRound                        uint8                // offset: 0, size: 1 bytes
	NRound                           uint8                // offset: 1, size: 1 bytes
	Red                              uint16               // offset: 2, size: 2 bytes
	Blue                             uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_KQ_SCORE_BOARD_INFO_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_KQ_SCORE_CMD (sizeof = 5 bytes)
type PROTO_NC_KQ_SCORE_CMD_DATA struct {
	Num                              uint8                // offset: 0, size: 1 bytes
	Base                             uint32               // offset: 1, size: 4 bytes
	Array                            [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_KQ_SCORE_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_KQ_SCORE_INFO_CMD (sizeof = 8 bytes)
type PROTO_NC_KQ_SCORE_INFO_CMD_DATA struct {
	NScore                           [2]uint32            // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_KQ_SCORE_INFO_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_KQ_SCORE_SIMPLE_CMD (sizeof = 18 bytes)
type PROTO_NC_KQ_SCORE_SIMPLE_CMD_DATA struct {
	Num                              uint8                // offset: 0, size: 1 bytes
	Base                             uint32               // offset: 1, size: 4 bytes
	Self                             [13]byte             // offset: 5, size: 13 bytes
	Array                            [0]byte              // offset: 18, flexible array member
}

func (s *PROTO_NC_KQ_SCORE_SIMPLE_CMD_DATA) FixedWireSize() int { return 18 }

// PROTO_NC_KQ_TEAM_TYPE_CMD (sizeof = 1 bytes)
type PROTO_NC_KQ_TEAM_TYPE_CMD_DATA struct {
	NTeamType                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_KQ_TEAM_TYPE_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_KQ_W2Z_DESTROY_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_W2Z_DESTROY_CMD_DATA)

// PROTO_NC_KQ_W2Z_MAKE_REQ (sizeof = 377 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_W2Z_MAKE_REQ_DATA)

// PROTO_NC_KQ_W2Z_START_CMD (sizeof = 379 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_W2Z_START_CMD_DATA)

// PROTO_NC_KQ_WINTER_EVENT_2014_SCORE_CMD (sizeof = 4 bytes)
type PROTO_NC_KQ_WINTER_EVENT_2014_SCORE_CMD_DATA struct {
	Red                              uint16               // offset: 0, size: 2 bytes
	Blue                             uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_KQ_WINTER_EVENT_2014_SCORE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_KQ_Z2W_END_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_Z2W_END_CMD_DATA)

// PROTO_NC_KQ_Z2W_MAKE_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_KQ_Z2W_MAKE_ACK_DATA)

// PROTO_NC_MENU_OPENSTORAGE_CMD (sizeof = 12 bytes)
type PROTO_NC_MENU_OPENSTORAGE_CMD_DATA struct {
	Cen                              uint64               // offset: 0, size: 8 bytes
	Maxpage                          uint8                // offset: 8, size: 1 bytes
	Curpage                          uint8                // offset: 9, size: 1 bytes
	NOpenType                        uint8                // offset: 10, size: 1 bytes
	Itemcounter                      uint8                // offset: 11, size: 1 bytes
	Itemarray                        [0]byte              // offset: 12, flexible array member
}

func (s *PROTO_NC_MENU_OPENSTORAGE_CMD_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_MENU_SERVERMENU_ACK (sizeof = 1 bytes)
type PROTO_NC_MENU_SERVERMENU_ACK_DATA struct {
	Reply                            uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_MENU_SERVERMENU_ACK_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_MENU_SERVERMENU_REQ (sizeof = 142 bytes)
type PROTO_NC_MENU_SERVERMENU_REQ_DATA struct {
	Title                            [128]byte            // offset: 0, size: 128 bytes
	Priority                         uint8                // offset: 128, size: 1 bytes
	NpcHandle                        uint16               // offset: 129, size: 2 bytes
	NpcPosition                      SHINE_XY_TYPE        // offset: 131, size: 8 bytes
	LimitRange                       uint16               // offset: 139, size: 2 bytes
	Menunum                          uint8                // offset: 141, size: 1 bytes
	Menu                             [0]byte              // offset: 142, flexible array member
}

func (s *PROTO_NC_MENU_SERVERMENU_REQ_DATA) FixedWireSize() int { return 142 }

// PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD (sizeof = 24 bytes)
type PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD_DATA struct {
	Hp                               [12]byte             // offset: 0, size: 12 bytes
	Sp                               [12]byte             // offset: 12, size: 12 bytes
}

func (s *PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD_DATA) FixedWireSize() int { return 24 }

// PROTO_NC_MENU_SHOPOPENTABLE_CMD (sizeof = 4 bytes)
type PROTO_NC_MENU_SHOPOPENTABLE_CMD struct {
	Itemnum                          uint16               // offset: 0, size: 2 bytes
	Npc                              uint16               // offset: 2, size: 2 bytes
	Itemlist                         [0]byte              // offset: 4, flexible array member
}

func (s *PROTO_NC_MENU_SHOPOPENTABLE_CMD) FixedWireSize() int { return 4 }

// PROTO_NC_MENU_SHOPOPEN_CMD (sizeof = 4 bytes)
type PROTO_NC_MENU_SHOPOPEN_CMD struct {
	Itemnum                          uint16               // offset: 0, size: 2 bytes
	Npc                              uint16               // offset: 2, size: 2 bytes
	Itemlist                         [0]byte              // offset: 4, flexible array member
}

func (s *PROTO_NC_MENU_SHOPOPEN_CMD) FixedWireSize() int { return 4 }

// PROTO_NC_MINIHOUSE_ACTIV_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_ACTIV_REQ_DATA)

// PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ_DATA)

// PROTO_NC_MINIHOUSE_BUILDING_CMD (sizeof = 66 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_BUILDING_CMD_DATA)

// PROTO_NC_MINIHOUSE_BUILDING_REQ (sizeof = 133 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_BUILDING_REQ_DATA)

// PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK_DATA)

// PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD_DATA)

// PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ_DATA)

// PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ_DATA)

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD_DATA)

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ_DATA)

// PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_ACK (sizeof = 17 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_ACK_DATA)

// PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_REQ (sizeof = 25 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_REQ_DATA)

// PROTO_NC_MINIHOUSE_DB_PORTAL_DEL_ACK (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_PORTAL_DEL_ACK_DATA)

// PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_ACK_DATA)

// PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_REQ_DATA)

// PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_ACK (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_ACK_DATA)

// PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_REQ (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_REQ_DATA)

// PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ_DATA)

// PROTO_NC_MINIHOUSE_FUNICHER_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FUNICHER_CMD_DATA)

// PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ_DATA)

// PROTO_NC_MINIHOUSE_FURNITURE_ENDURE_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_FURNITURE_ENDURE_CMD_DATA)

// PROTO_NC_MINIHOUSE_KICKOUT_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_KICKOUT_CMD_DATA)

// PROTO_NC_MINIHOUSE_KICKOUT_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_KICKOUT_REQ_DATA)

// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK (sizeof = 140 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK_DATA)

// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD_DATA)

// PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD (sizeof = 23 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ (sizeof = 9 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ_DATA)

// PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ (sizeof = 21 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_ADD_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_ADD_CMD_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD_DATA)

// PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ_DATA)

// PROTO_NC_MINIHOUSE_REARRANGE_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_REARRANGE_REQ_DATA)

// PROTO_NC_MINIHOUSE_VISITREADY_ACK (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_VISITREADY_ACK_DATA)

// PROTO_NC_MINIHOUSE_VISITREADY_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_VISITREADY_REQ_DATA)

// PROTO_NC_MINIHOUSE_VISIT_REQ (sizeof = 11 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MINIHOUSE_VISIT_REQ_DATA)

// PROTO_NC_MOVER_HUNGRY_CMD (sizeof = 2 bytes)
type PROTO_NC_MOVER_HUNGRY_CMD_DATA struct {
	NHungry                          uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_MOVER_HUNGRY_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_MOVER_MOVESPEED_CMD (sizeof = 6 bytes)
type PROTO_NC_MOVER_MOVESPEED_CMD_DATA struct {
	NMoverHandle                     uint16               // offset: 0, size: 2 bytes
	NWalk                            uint16               // offset: 2, size: 2 bytes
	NRun                             uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_MOVER_MOVESPEED_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_MOVER_RIDE_ON_CMD (sizeof = 4 bytes)
type PROTO_NC_MOVER_RIDE_ON_CMD_DATA struct {
	NMoverHandle                     uint16               // offset: 0, size: 2 bytes
	NSlot                            uint8                // offset: 2, size: 1 bytes
	NGrade                           uint8                // offset: 3, size: 1 bytes
}

func (s *PROTO_NC_MOVER_RIDE_ON_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_MOVER_RIDE_ON_REQ (sizeof = 3 bytes)
type PROTO_NC_MOVER_RIDE_ON_REQ_DATA struct {
	NMoverHandle                     uint16               // offset: 0, size: 2 bytes
	NSlot                            uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_MOVER_RIDE_ON_REQ_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ (sizeof = 10 bytes)
type PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ_DATA struct {
	NSkillID                         uint16               // offset: 0, size: 2 bytes
	NTargetLoc                       SHINE_XY_TYPE        // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ (sizeof = 4 bytes)
type PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ_DATA struct {
	NSkillID                         uint16               // offset: 0, size: 2 bytes
	NTarget                          uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD (sizeof = 2 bytes)
type PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD (sizeof = 6 bytes)
type PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NMoverHandle                     uint16               // offset: 2, size: 2 bytes
	NSlot                            uint8                // offset: 4, size: 1 bytes
	NGrade                           uint8                // offset: 5, size: 1 bytes
}

func (s *PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_PARTY_CHANGEMASTER_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_CHANGEMASTER_CMD_DATA)

// PROTO_NC_PARTY_DISMISS_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_DISMISS_CMD_DATA)

// PROTO_NC_PARTY_FUNDAMENTAL_CMD (sizeof = 66 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_FUNDAMENTAL_CMD_DATA)

// PROTO_NC_PARTY_ITEM_LOOTING_ZONE_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_ITEM_LOOTING_ZONE_CMD_DATA)

// PROTO_NC_PARTY_JOIN_CMD (sizeof = 34 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_JOIN_CMD_DATA)

// PROTO_NC_PARTY_KICKOFF_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_KICKOFF_CMD_DATA)

// PROTO_NC_PARTY_LEAVE_CMD (sizeof = 33 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_LEAVE_CMD_DATA)

// PROTO_NC_PARTY_LOGIN_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_LOGIN_CMD_DATA)

// PROTO_NC_PARTY_LOGOUT_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_LOGOUT_CMD_DATA)

// PROTO_NC_PARTY_MEMBERCLASS_CMD (sizeof = 1 bytes)
type PROTO_NC_PARTY_MEMBERCLASS_CMD_DATA struct {
	Membernum                        uint8                // offset: 0, size: 1 bytes
	Member                           [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_PARTY_MEMBERCLASS_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_PARTY_MEMBERINFORM_ACK (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_MEMBERINFORM_ACK_DATA)

// PROTO_NC_PARTY_MEMBERINFORM_CMD (sizeof = 1 bytes)
type PROTO_NC_PARTY_MEMBERINFORM_CMD_DATA struct {
	Membernum                        uint8                // offset: 0, size: 1 bytes
	Member                           [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_PARTY_MEMBERINFORM_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_PARTY_MEMBERINFORM_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_MEMBERINFORM_REQ_DATA)

// PROTO_NC_PARTY_MEMBERLOCATION_CMD (sizeof = 1 bytes)
type PROTO_NC_PARTY_MEMBERLOCATION_CMD_DATA struct {
	Membernum                        uint8                // offset: 0, size: 1 bytes
	Member                           [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_PARTY_MEMBERLOCATION_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_PARTY_MEMBERMAPOUT (sizeof = 20 bytes)
type PROTO_NC_PARTY_MEMBERMAPOUT_DATA struct {
	Memberid                         Name5                // offset: 0, size: 20 bytes
}

func (s *PROTO_NC_PARTY_MEMBERMAPOUT_DATA) FixedWireSize() int { return 20 }

// PROTO_NC_PARTY_ZONE_JOIN_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_ZONE_JOIN_CMD_DATA)

// PROTO_NC_PARTY_ZONE_LEAVE_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_ZONE_LEAVE_CMD_DATA)

// PROTO_NC_PARTY_ZONE_SET_LOOTER_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PARTY_ZONE_SET_LOOTER_CMD_DATA)

// PROTO_NC_PET_ASK_NEW_NAME_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_ASK_NEW_NAME_ACK_DATA)

// PROTO_NC_PET_ASK_NEW_NAME_REQ (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_ASK_NEW_NAME_REQ_DATA)

// PROTO_NC_PET_CREATE_DB_ACK (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_CREATE_DB_ACK_DATA)

// PROTO_NC_PET_CREATE_DB_REQ (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_CREATE_DB_REQ_DATA)

// PROTO_NC_PET_LINK_RESUMMON_CMD (sizeof = 23 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_LINK_RESUMMON_CMD_DATA)

// PROTO_NC_PET_LOAD_INFO_DB_ACK (sizeof = 26 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_LOAD_INFO_DB_ACK_DATA)

// PROTO_NC_PET_LOAD_INFO_DB_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_LOAD_INFO_DB_REQ_DATA)

// PROTO_NC_PET_SET_NAME_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_ACK_DATA)

// PROTO_NC_PET_SET_NAME_CANCEL_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_CANCEL_ACK_DATA)

// PROTO_NC_PET_SET_NAME_CANCEL_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_CANCEL_REQ_DATA)

// PROTO_NC_PET_SET_NAME_CMD (sizeof = 18 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_CMD_DATA)

// PROTO_NC_PET_SET_NAME_DB_ACK (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_DB_ACK_DATA)

// PROTO_NC_PET_SET_NAME_DB_REQ (sizeof = 22 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_DB_REQ_DATA)

// PROTO_NC_PET_SET_NAME_REQ (sizeof = 18 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_NAME_REQ_DATA)

// PROTO_NC_PET_SET_TENDENCY_DB_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_SET_TENDENCY_DB_REQ_DATA)

// PROTO_NC_PET_USE_ITEM_FAIL_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PET_USE_ITEM_FAIL_ACK_DATA)

// PROTO_NC_PRISON_END_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_PRISON_END_ACK_DATA)

// PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK_DATA)

// PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ_DATA)

// PROTO_NC_QUEST_DB_CLEAR_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_CLEAR_ACK_DATA)

// PROTO_NC_QUEST_DB_CLEAR_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_CLEAR_REQ_DATA)

// PROTO_NC_QUEST_DB_DONE_ACK (sizeof = 43 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_DONE_ACK_DATA)

// PROTO_NC_QUEST_DB_DONE_REQ (sizeof = 83 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_DONE_REQ_DATA)

// PROTO_NC_QUEST_DB_GIVE_UP_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_GIVE_UP_ACK_DATA)

// PROTO_NC_QUEST_DB_GIVE_UP_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_GIVE_UP_REQ_DATA)

// PROTO_NC_QUEST_DB_READ_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_READ_ACK_DATA)

// PROTO_NC_QUEST_DB_READ_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_READ_REQ_DATA)

// PROTO_NC_QUEST_DB_SET_INFO_ACK (sizeof = 41 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_SET_INFO_ACK_DATA)

// PROTO_NC_QUEST_DB_SET_INFO_REQ (sizeof = 39 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_DB_SET_INFO_REQ_DATA)

// PROTO_NC_QUEST_GIVE_UP_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_GIVE_UP_ACK_DATA)

// PROTO_NC_QUEST_GIVE_UP_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_GIVE_UP_REQ_DATA)

// PROTO_NC_QUEST_JOBDUNGEON_FIND_RNG (sizeof = 115 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_JOBDUNGEON_FIND_RNG_DATA)

// PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD_DATA)

// PROTO_NC_QUEST_READ_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_READ_ACK_DATA)

// PROTO_NC_QUEST_READ_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_READ_REQ_DATA)

// PROTO_NC_QUEST_RESET_TIME_CMD (sizeof = 16 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_RESET_TIME_CMD_DATA)

// PROTO_NC_QUEST_RESET_TIME_ZONE_CMD (sizeof = 20 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_RESET_TIME_ZONE_CMD_DATA)

// PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD_DATA)

// PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD_DATA)

// PROTO_NC_QUEST_SCENARIO_RUN_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_SCENARIO_RUN_CMD_DATA)

// PROTO_NC_QUEST_SCRIPT_CMD_ACK (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_SCRIPT_CMD_ACK_DATA)

// PROTO_NC_QUEST_SCRIPT_CMD_REQ (sizeof = 103 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_SCRIPT_CMD_REQ_DATA)

// PROTO_NC_QUEST_SELECT_START_ACK (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_SELECT_START_ACK_DATA)

// PROTO_NC_QUEST_SELECT_START_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_SELECT_START_REQ_DATA)

// PROTO_NC_QUEST_START_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_START_ACK_DATA)

// PROTO_NC_QUEST_START_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_QUEST_START_REQ_DATA)

// PROTO_NC_RAID_LOCATION_CMD (sizeof = 1 bytes)
type PROTO_NC_RAID_LOCATION_CMD_DATA struct {
	NNum                             uint8                // offset: 0, size: 1 bytes
	RaidMembers                      [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_RAID_LOCATION_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_RAID_MAPOUT_CMD (sizeof = 20 bytes)
type PROTO_NC_RAID_MAPOUT_CMD_DATA struct {
	SCharID                          Name5                // offset: 0, size: 20 bytes
}

func (s *PROTO_NC_RAID_MAPOUT_CMD_DATA) FixedWireSize() int { return 20 }

// PROTO_NC_RAID_MEMBERCLASS_CMD (sizeof = 1 bytes)
type PROTO_NC_RAID_MEMBERCLASS_CMD_DATA struct {
	NNum                             uint8                // offset: 0, size: 1 bytes
	RaidMembers                      [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_RAID_MEMBERCLASS_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_RAID_MEMBERINFORM_CMD (sizeof = 16 bytes)
type PROTO_NC_RAID_MEMBERINFORM_CMD_DATA struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	Hp                               uint32               // offset: 4, size: 4 bytes
	Sp                               uint32               // offset: 8, size: 4 bytes
	Lp                               uint32               // offset: 12, size: 4 bytes
}

func (s *PROTO_NC_RAID_MEMBERINFORM_CMD_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_RAID_ZONE_DEL_SUB_LEADER_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_DEL_SUB_LEADER_CMD_DATA)

// PROTO_NC_RAID_ZONE_DISMISS_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_DISMISS_CMD_DATA)

// PROTO_NC_RAID_ZONE_LEAVE_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_LEAVE_CMD_DATA)

// PROTO_NC_RAID_ZONE_LOOTING_TYPE_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_LOOTING_TYPE_CMD_DATA)

// PROTO_NC_RAID_ZONE_MAKE_CMD (sizeof = 32 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_MAKE_CMD_DATA)

// PROTO_NC_RAID_ZONE_MEMBER_JOIN_CMD (sizeof = 42 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_MEMBER_JOIN_CMD_DATA)

// PROTO_NC_RAID_ZONE_SET_LEADER_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_SET_LEADER_CMD_DATA)

// PROTO_NC_RAID_ZONE_SET_LOOTER_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_SET_LOOTER_CMD_DATA)

// PROTO_NC_RAID_ZONE_SET_SUB_LEADER_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_SET_SUB_LEADER_CMD_DATA)

// PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD (sizeof = 50 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_DATA)

// PROTO_NC_RAID_ZONE_VANISH_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_RAID_ZONE_VANISH_CMD_DATA)

// PROTO_NC_SCENARIO_ANIMATE_CMD (sizeof = 34 bytes)
type PROTO_NC_SCENARIO_ANIMATE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Anicode                          Name8                // offset: 2, size: 32 bytes
}

func (s *PROTO_NC_SCENARIO_ANIMATE_CMD_DATA) FixedWireSize() int { return 34 }

// PROTO_NC_SCENARIO_AREAENTRY_ACK (sizeof = 32 bytes)
type PROTO_NC_SCENARIO_AREAENTRY_ACK_DATA struct {
	Areaindex                        Name8                // offset: 0, size: 32 bytes
}

func (s *PROTO_NC_SCENARIO_AREAENTRY_ACK_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_SCENARIO_AREAENTRY_REQ (sizeof = 32 bytes)
type PROTO_NC_SCENARIO_AREAENTRY_REQ_DATA struct {
	Areaindex                        Name8                // offset: 0, size: 32 bytes
}

func (s *PROTO_NC_SCENARIO_AREAENTRY_REQ_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_SCENARIO_AREALEAVE_ACK (sizeof = 32 bytes)
type PROTO_NC_SCENARIO_AREALEAVE_ACK_DATA struct {
	Areaindex                        Name8                // offset: 0, size: 32 bytes
}

func (s *PROTO_NC_SCENARIO_AREALEAVE_ACK_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_SCENARIO_AREALEAVE_REQ (sizeof = 32 bytes)
type PROTO_NC_SCENARIO_AREALEAVE_REQ_DATA struct {
	Areaindex                        Name8                // offset: 0, size: 32 bytes
}

func (s *PROTO_NC_SCENARIO_AREALEAVE_REQ_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_SCENARIO_BRIEFEFFECT_CMD (sizeof = 1 bytes)
type PROTO_NC_SCENARIO_BRIEFEFFECT_CMD_DATA struct {
	Effectcode                       uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_SCENARIO_BRIEFEFFECT_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_SCENARIO_CAMERA_MOVE_CMD (sizeof = 15 bytes)
type PROTO_NC_SCENARIO_CAMERA_MOVE_CMD_DATA struct {
	MoveType                         uint8                // offset: 0, size: 1 bytes
	ViewCoord                        SHINE_XY_TYPE        // offset: 1, size: 8 bytes
	AngleXZ                          uint16               // offset: 9, size: 2 bytes
	AngleY                           uint16               // offset: 11, size: 2 bytes
	Dist                             uint16               // offset: 13, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_CAMERA_MOVE_CMD_DATA) FixedWireSize() int { return 15 }

// PROTO_NC_SCENARIO_CHATWIN_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_CHATWIN_CMD_DATA struct {
	Mob                              uint16               // offset: 0, size: 2 bytes
	Len                              uint8                // offset: 2, size: 1 bytes
	Content                          [0]byte              // offset: 3, flexible array member
}

func (s *PROTO_NC_SCENARIO_CHATWIN_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_DEL_DIRECTIONAL_ARROW_CMD (sizeof = 1 bytes)
type PROTO_NC_SCENARIO_DEL_DIRECTIONAL_ARROW_CMD_DATA struct {}

func (s *PROTO_NC_SCENARIO_DEL_DIRECTIONAL_ARROW_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_SCENARIO_DIALOG_CMD (sizeof = 34 bytes)
type PROTO_NC_SCENARIO_DIALOG_CMD_DATA struct {
	SFacecutIndex                    [33]byte             // offset: 0, size: 33 bytes
	NLen                             uint8                // offset: 33, size: 1 bytes
	SContents                        [0]byte              // offset: 34, flexible array member
}

func (s *PROTO_NC_SCENARIO_DIALOG_CMD_DATA) FixedWireSize() int { return 34 }

// PROTO_NC_SCENARIO_DIRECTIONAL_ARROW_CMD (sizeof = 8 bytes)
type PROTO_NC_SCENARIO_DIRECTIONAL_ARROW_CMD_DATA struct {
	DirectionalPos                   SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_SCENARIO_DIRECTIONAL_ARROW_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_SCENARIO_DOORSTATE_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_DOORSTATE_CMD_DATA struct {
	Door                             uint16               // offset: 0, size: 2 bytes
	Doorstate                        uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_SCENARIO_DOORSTATE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_EFFECT_TIMER_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_EFFECT_TIMER_CMD_DATA struct {
	NTimerFunction                   uint8                // offset: 0, size: 1 bytes
	NSecond                          uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_EFFECT_TIMER_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_FOG_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_FOG_CMD_DATA struct {
	Fog                              uint8                // offset: 0, size: 1 bytes
	Depth                            uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_FOG_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_HIDE_OTHER_PLAYER_CMD (sizeof = 1 bytes)
type PROTO_NC_SCENARIO_HIDE_OTHER_PLAYER_CMD_DATA struct {
	BIsHide                          uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_SCENARIO_HIDE_OTHER_PLAYER_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_SCENARIO_LIGHT_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_LIGHT_CMD_DATA struct {
	Red                              uint8                // offset: 0, size: 1 bytes
	Green                            uint8                // offset: 1, size: 1 bytes
	Blue                             uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_SCENARIO_LIGHT_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_MAPMARK_CMD (sizeof = 2 bytes)
type PROTO_NC_SCENARIO_MAPMARK_CMD_DATA struct {
	NNum                             uint16               // offset: 0, size: 2 bytes
	IconList                         [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_SCENARIO_MAPMARK_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SCENARIO_MAP_OBJECT_CONTROL_CMD (sizeof = 37 bytes)
type PROTO_NC_SCENARIO_MAP_OBJECT_CONTROL_CMD_DATA struct {
	SObjectIndex                     [32]byte             // offset: 0, size: 32 bytes
	IsEnableObject                   uint8                // offset: 32, size: 1 bytes
	EControlType                     uint32               // offset: 33, size: 4 bytes
}

func (s *PROTO_NC_SCENARIO_MAP_OBJECT_CONTROL_CMD_DATA) FixedWireSize() int { return 37 }

// PROTO_NC_SCENARIO_MESSAGENOTICE_CMD (sizeof = 4 bytes)
type PROTO_NC_SCENARIO_MESSAGENOTICE_CMD_DATA struct {
	Message                          uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_SCENARIO_MESSAGENOTICE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_SCENARIO_MESSAGE_CMD (sizeof = 17 bytes)
type PROTO_NC_SCENARIO_MESSAGE_CMD_DATA struct {
	DialogIndex                      [17]byte             // offset: 0, size: 17 bytes
}

func (s *PROTO_NC_SCENARIO_MESSAGE_CMD_DATA) FixedWireSize() int { return 17 }

// PROTO_NC_SCENARIO_NPCCHAT_CMD (sizeof = 35 bytes)
type PROTO_NC_SCENARIO_NPCCHAT_CMD_DATA struct {
	DialogIndex                      [33]byte             // offset: 0, size: 33 bytes
	NPCHandle                        uint16               // offset: 33, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_NPCCHAT_CMD_DATA) FixedWireSize() int { return 35 }

// PROTO_NC_SCENARIO_OBJTYPECHANGE_CMD (sizeof = 3 bytes)
type PROTO_NC_SCENARIO_OBJTYPECHANGE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	TypeField                        uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_SCENARIO_OBJTYPECHANGE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_SCENARIO_RANKING_LIST_CMD (sizeof = 1000 bytes)
type PROTO_NC_SCENARIO_RANKING_LIST_CMD_DATA struct {
	PreviousRanking                  [500]byte            // offset: 0, size: 500 bytes
	CurrentRanking                   [500]byte            // offset: 500, size: 500 bytes
}

func (s *PROTO_NC_SCENARIO_RANKING_LIST_CMD_DATA) FixedWireSize() int { return 1000 }

// PROTO_NC_SCENARIO_SCRIPTMESSAGE_CMD (sizeof = 37 bytes)
type PROTO_NC_SCENARIO_SCRIPTMESSAGE_CMD_DATA struct {
	DialogIndex                      [32]byte             // offset: 0, size: 32 bytes
	Nlen                             [5]byte              // offset: 32, size: 5 bytes
	SContent                         [0]byte              // offset: 37, flexible array member
}

func (s *PROTO_NC_SCENARIO_SCRIPTMESSAGE_CMD_DATA) FixedWireSize() int { return 37 }

// PROTO_NC_SCENARIO_SYSTEM_MESSAGE_CMD (sizeof = 1 bytes)
type PROTO_NC_SCENARIO_SYSTEM_MESSAGE_CMD_DATA struct {
	NLen                             uint8                // offset: 0, size: 1 bytes
	SMessage                         [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_SCENARIO_SYSTEM_MESSAGE_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_SCENARIO_TIMER_CMD (sizeof = 2 bytes)
type PROTO_NC_SCENARIO_TIMER_CMD_DATA struct {
	NSecond                          uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_TIMER_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SCENARIO_TIMER_END_CMD (sizeof = 4 bytes)
type PROTO_NC_SCENARIO_TIMER_END_CMD_DATA struct {
	NDeleteSec                       uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_SCENARIO_TIMER_END_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_SCENARIO_TIMER_START_CMD (sizeof = 1 bytes)
type PROTO_NC_SCENARIO_TIMER_START_CMD_DATA struct {}

func (s *PROTO_NC_SCENARIO_TIMER_START_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_SCENARIO_TOPVIEW_CMD (sizeof = 13 bytes)
type PROTO_NC_SCENARIO_TOPVIEW_CMD_DATA struct {
	NMoveType                        uint8                // offset: 0, size: 1 bytes
	CenterPos                        SHINE_XY_TYPE        // offset: 1, size: 8 bytes
	NDistance                        uint16               // offset: 9, size: 2 bytes
	NAngle                           uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_SCENARIO_TOPVIEW_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_SKILL_COOLTIME_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_COOLTIME_CMD_DATA)

// PROTO_NC_SKILL_EMPOWALLOC_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOWALLOC_ACK_DATA)

// PROTO_NC_SKILL_EMPOWALLOC_DB_ACK (sizeof = 13 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOWALLOC_DB_ACK_DATA)

// PROTO_NC_SKILL_EMPOWALLOC_DB_REQ (sizeof = 40 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOWALLOC_DB_REQ_DATA)

// PROTO_NC_SKILL_EMPOWALLOC_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOWALLOC_REQ_DATA)

// PROTO_NC_SKILL_EMPOWPOINT_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOWPOINT_CMD_DATA)

// PROTO_NC_SKILL_EMPOW_RESET_DB_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOW_RESET_DB_REQ_DATA)

// PROTO_NC_SKILL_EMPOW_RESET_SUC_ACK (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_EMPOW_RESET_SUC_ACK_DATA)

// PROTO_NC_SKILL_ERASE_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_ERASE_ACK_DATA)

// PROTO_NC_SKILL_ERASE_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_ERASE_REQ_DATA)

// PROTO_NC_SKILL_JUMP_CMD (sizeof = 21 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_JUMP_CMD_DATA)

// PROTO_NC_SKILL_PASSIVESKILL_LEARN_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_PASSIVESKILL_LEARN_CMD_DATA)

// PROTO_NC_SKILL_PRODUCTFIELD_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_PRODUCTFIELD_ACK_DATA)

// PROTO_NC_SKILL_PRODUCTFIELD_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_PRODUCTFIELD_REQ_DATA)

// PROTO_NC_SKILL_REVIVE_CMD (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_REVIVE_CMD_DATA)

// PROTO_NC_SKILL_SKILLEXP_CLIENT_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_SKILLEXP_CLIENT_CMD_DATA)

// PROTO_NC_SKILL_SKILLTEACHSUC_ACK (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_SKILLTEACHSUC_ACK_DATA)

// PROTO_NC_SKILL_SKILLTEACH_REQ (sizeof = 9 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_SKILLTEACH_REQ_DATA)

// PROTO_NC_SKILL_SOMEONEREVIVE_CMD (sizeof = 14 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_SOMEONEREVIVE_CMD_DATA)

// PROTO_NC_SKILL_UNLEARN_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_UNLEARN_ACK_DATA)

// PROTO_NC_SKILL_UNLEARN_REQ (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_UNLEARN_REQ_DATA)

// PROTO_NC_SKILL_WARP_CMD (sizeof = 20 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_SKILL_WARP_CMD_DATA)

// PROTO_NC_SOULSTONE_HP_BUY_REQ (sizeof = 2 bytes)
type PROTO_NC_SOULSTONE_HP_BUY_REQ_DATA struct {
	Number                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SOULSTONE_HP_BUY_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD (sizeof = 2 bytes)
type PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SOULSTONE_SP_BUY_ACK (sizeof = 2 bytes)
type PROTO_NC_SOULSTONE_SP_BUY_ACK_DATA struct {
	Totalnumber                      uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SOULSTONE_SP_BUY_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SOULSTONE_SP_BUY_REQ (sizeof = 2 bytes)
type PROTO_NC_SOULSTONE_SP_BUY_REQ_DATA struct {
	Number                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SOULSTONE_SP_BUY_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD (sizeof = 2 bytes)
type PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_TRADE_CENBOARDING_ACK (sizeof = 8 bytes)
type PROTO_NC_TRADE_CENBOARDING_ACK_DATA struct {
	Cen                              uint64               // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_TRADE_CENBOARDING_ACK_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_TRADE_DOWNBOARD_ACK (sizeof = 1 bytes)
type PROTO_NC_TRADE_DOWNBOARD_ACK_DATA struct {
	Slotboard                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_TRADE_DOWNBOARD_ACK_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_TRADE_DOWNBOARD_REQ (sizeof = 1 bytes)
type PROTO_NC_TRADE_DOWNBOARD_REQ_DATA struct {
	Slotboard                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_TRADE_DOWNBOARD_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_TRADE_OPPOSITCENBOARDING_CMD (sizeof = 8 bytes)
type PROTO_NC_TRADE_OPPOSITCENBOARDING_CMD_DATA struct {
	Cen                              uint64               // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_TRADE_OPPOSITCENBOARDING_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_TRADE_OPPOSITDOWNBOARD_CMD (sizeof = 1 bytes)
type PROTO_NC_TRADE_OPPOSITDOWNBOARD_CMD_DATA struct {
	Slotboard                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_TRADE_OPPOSITDOWNBOARD_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_TRADE_OPPOSITUPBOARD_CMD (sizeof = 104 bytes)
type PROTO_NC_TRADE_OPPOSITUPBOARD_CMD_DATA struct {
	Slotboard                        uint8                // offset: 0, size: 1 bytes
	Iteminfo                         [103]byte            // offset: 1, size: 103 bytes
}

func (s *PROTO_NC_TRADE_OPPOSITUPBOARD_CMD_DATA) FixedWireSize() int { return 104 }

// PROTO_NC_TRADE_PROPOSE_ASK_REQ (sizeof = 2 bytes)
type PROTO_NC_TRADE_PROPOSE_ASK_REQ_DATA struct {
	Proposer                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_TRADE_PROPOSE_ASK_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_TRADE_PROPOSE_REQ (sizeof = 2 bytes)
type PROTO_NC_TRADE_PROPOSE_REQ_DATA struct {
	Proposee                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_TRADE_PROPOSE_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_TRADE_START_CMD (sizeof = 2 bytes)
type PROTO_NC_TRADE_START_CMD_DATA struct {
	Opposite                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_TRADE_START_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_TRADE_UPBOARD_ACK (sizeof = 2 bytes)
type PROTO_NC_TRADE_UPBOARD_ACK_DATA struct {
	Slotinven                        uint8                // offset: 0, size: 1 bytes
	Slotboard                        uint8                // offset: 1, size: 1 bytes
}

func (s *PROTO_NC_TRADE_UPBOARD_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_TRADE_UPBOARD_REQ (sizeof = 1 bytes)
type PROTO_NC_TRADE_UPBOARD_REQ_DATA struct {
	Slotinven                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_TRADE_UPBOARD_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_WT_GRADE_CMD (sizeof = 3 bytes)
type PROTO_NC_WT_GRADE_CMD_DATA struct {
	ItemInven                        ITEM_INVEN           // offset: 0, size: 2 bytes
	Grade                            uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_WT_GRADE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_WT_LICENSE_ACK (sizeof = 7 bytes)
type PROTO_NC_WT_LICENSE_ACK_DATA struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
	BSet                             uint8                // offset: 4, size: 1 bytes
	Err                              uint16               // offset: 5, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_ACK_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_WT_LICENSE_CLR_ACK (sizeof = 24 bytes)
type PROTO_NC_WT_LICENSE_CLR_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Echo_data                        uint32               // offset: 8, size: 4 bytes
	WeaponItemKey                    uint64               // offset: 12, size: 8 bytes
	MobID                            uint16               // offset: 20, size: 2 bytes
	ErrorCode                        uint16               // offset: 22, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_CLR_ACK_DATA) FixedWireSize() int { return 24 }

// PROTO_NC_WT_LICENSE_CLR_REQ (sizeof = 35 bytes)
type PROTO_NC_WT_LICENSE_CLR_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Echo_data                        uint32               // offset: 8, size: 4 bytes
	WeaponItemKey                    uint64               // offset: 12, size: 8 bytes
	MobID                            uint16               // offset: 20, size: 2 bytes
	LicenseClearItemKey              uint64               // offset: 22, size: 8 bytes
	LicenseClearItemID               uint16               // offset: 30, size: 2 bytes
	LicenseClearItemDelete           uint8                // offset: 32, size: 1 bytes
	LicenseClearAttr                 uint16               // offset: 33, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_CLR_REQ_DATA) FixedWireSize() int { return 35 }

// PROTO_NC_WT_LICENSE_REQ (sizeof = 7 bytes)
type PROTO_NC_WT_LICENSE_REQ_DATA struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
	Mobid                            uint16               // offset: 4, size: 2 bytes
	BSet                             uint8                // offset: 6, size: 1 bytes
}

func (s *PROTO_NC_WT_LICENSE_REQ_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_WT_LICENSE_SET_ACK (sizeof = 14 bytes)
type PROTO_NC_WT_LICENSE_SET_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	Echo_data                        uint32               // offset: 10, size: 4 bytes
}

func (s *PROTO_NC_WT_LICENSE_SET_ACK_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_WT_LICENSE_SET_REQ (sizeof = 54 bytes)
type PROTO_NC_WT_LICENSE_SET_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Chrregnum                        uint32               // offset: 8, size: 4 bytes
	Chrname                          Name5                // offset: 12, size: 20 bytes
	WeaponItemKey                    uint64               // offset: 32, size: 8 bytes
	LicenItemKey                     uint64               // offset: 40, size: 8 bytes
	LicenseMobID                     uint16               // offset: 48, size: 2 bytes
	Echo_data                        uint32               // offset: 50, size: 4 bytes
}

func (s *PROTO_NC_WT_LICENSE_SET_REQ_DATA) FixedWireSize() int { return 54 }

// PROTO_NC_WT_MOBINC_CMD (sizeof = 2 bytes)
type PROTO_NC_WT_MOBINC_CMD_DATA struct {
	Mob                              uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_WT_MOBINC_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_WT_MOB_KILLCOUNT_SET_ACK (sizeof = 22 bytes)
type PROTO_NC_WT_MOB_KILLCOUNT_SET_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemKey                          uint64               // offset: 6, size: 8 bytes
	MobID                            uint16               // offset: 14, size: 2 bytes
	MobKillCount                     uint32               // offset: 16, size: 4 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

func (s *PROTO_NC_WT_MOB_KILLCOUNT_SET_ACK_DATA) FixedWireSize() int { return 22 }

// PROTO_NC_WT_MOB_KILLCOUNT_SET_REQ (sizeof = 20 bytes)
type PROTO_NC_WT_MOB_KILLCOUNT_SET_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemKey                          uint64               // offset: 6, size: 8 bytes
	MobID                            uint16               // offset: 14, size: 2 bytes
	MobKillCount                     uint32               // offset: 16, size: 4 bytes
}

func (s *PROTO_NC_WT_MOB_KILLCOUNT_SET_REQ_DATA) FixedWireSize() int { return 20 }

// PROTO_NC_WT_TITLE_SET_ACK (sizeof = 8 bytes)
type PROTO_NC_WT_TITLE_SET_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
}

func (s *PROTO_NC_WT_TITLE_SET_ACK_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_ZONE_PRISON_END_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ZONE_PRISON_END_CMD_DATA)

// PROTO_NC_ZONE_PRISON_GO_ACK (sizeof = 374 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ZONE_PRISON_GO_ACK_DATA)

// PROTO_NC_ZONE_PRISON_GO_REQ (sizeof = 372 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ZONE_PRISON_GO_REQ_DATA)
