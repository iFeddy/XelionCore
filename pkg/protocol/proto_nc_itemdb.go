package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHARSAVE_2WLDMAN_ITEMACTIONCOOLTIME_CMD (size = 6 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_ITEMACTIONCOOLTIME_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	ActionCoolTime                   uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD (size = 9 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	PacketOrder                      uint8                // offset: 6, size: 1 bytes
	NPartMark                        uint8                // offset: 7, size: 1 bytes
	Itemstr                          uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK (size = 2 bytes)
type PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK_DATA struct {
	Num                              uint16               // offset: 0, size: 2 bytes
	Group                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_REQ (size = 2 bytes)
type PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_REQ_DATA struct {
	Num                              uint16               // offset: 0, size: 2 bytes
	Group                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_CHARSAVE_SELL_ITEM_INFO_CMD (size = 7 bytes)
type PROTO_NC_CHARSAVE_SELL_ITEM_INFO_CMD_DATA struct {
	NWorldHandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	NumOfSellItem                    uint8                // offset: 6, size: 1 bytes
	SellItemList                     [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_DB_CMD (size = 78 bytes)
type PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_DB_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	UseItemMinimonInfo               [50]byte             // offset: 4, size: 50 bytes
	DelSlotList_Normal               [12]byte             // offset: 54, size: 12 bytes
	DelSlotList_Charged              [12]byte             // offset: 66, size: 12 bytes
}

// PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_WORLD_CMD (size = 56 bytes)
type PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_WORLD_CMD_DATA struct {
	NWorldHandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	UseItemMinimonInfo               [50]byte             // offset: 6, size: 50 bytes
}

// PROTO_NC_CHAR_CLIENT_ITEM_CMD (size = 3 bytes)
type PROTO_NC_CHAR_CLIENT_ITEM_CMD_DATA struct {
	Numofitem                        uint8                // offset: 0, size: 1 bytes
	Box                              uint8                // offset: 1, size: 1 bytes
	Flag                             uint8                // offset: 2, size: 1 bytes
	ItemArray                        [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_ACK (size = 16 bytes)
type PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenType                       uint8                // offset: 6, size: 1 bytes
	NOwner                           uint32               // offset: 7, size: 4 bytes
	Error                            uint16               // offset: 11, size: 2 bytes
	ItemCmd                          [3]byte              // offset: 13, size: 3 bytes
}

// PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_ACK (size = 16 bytes)
type PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenType                       uint8                // offset: 6, size: 1 bytes
	NOwner                           uint32               // offset: 7, size: 4 bytes
	Error                            uint16               // offset: 11, size: 2 bytes
	ItemCmd                          [3]byte              // offset: 13, size: 3 bytes
}

// PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_REQ (size = 16 bytes)
type PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenType                       uint8                // offset: 6, size: 1 bytes
	NOwner                           uint32               // offset: 7, size: 4 bytes
	SlotORnumber                     uint8                // offset: 11, size: 1 bytes
	From                             uint16               // offset: 12, size: 2 bytes
	To                               uint16               // offset: 14, size: 2 bytes
}

// PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_REQ (size = 11 bytes)
type PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenType                       uint8                // offset: 6, size: 1 bytes
	NOwner                           uint32               // offset: 7, size: 4 bytes
}

// PROTO_NC_CHAR_ITEM_CMD (size = 3 bytes)
type PROTO_NC_CHAR_ITEM_CMD_DATA struct {
	PacketOrder                      uint8                // offset: 0, size: 1 bytes
	NPartMark                        uint8                // offset: 1, size: 1 bytes
	ItemCmd                          uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_CHAR_SELL_ITEM_INFO_ZONE_CMD (size = 1 bytes)
type PROTO_NC_CHAR_SELL_ITEM_INFO_ZONE_CMD_DATA struct {
	NumOfSellItem                    uint8                // offset: 0, size: 1 bytes
	SellItemList                     [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_CHAR_TUTORIAL_MAKE_ITEM_ACK (size = 9 bytes)
type PROTO_NC_CHAR_TUTORIAL_MAKE_ITEM_ACK_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	RErrorType                       uint16               // offset: 6, size: 2 bytes
	BIsSkip                          uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_CHAR_TUTORIAL_MAKE_ITEM_REQ (size = 12 bytes)
type PROTO_NC_CHAR_TUTORIAL_MAKE_ITEM_REQ_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NCreateWorldNo                   uint8                // offset: 2, size: 1 bytes
	NCharNo                          uint32               // offset: 3, size: 4 bytes
	NUserNo                          uint32               // offset: 7, size: 4 bytes
	BIsSkip                          uint8                // offset: 11, size: 1 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ (size = 1 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ_DATA struct {}

// PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ (size = 24 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ_DATA struct {
	ChargedItemList                  [12]uint16           // offset: 0, size: 24 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD (size = 50 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD_DATA struct {
	UseItemMinimonInfo               [50]byte             // offset: 0, size: 50 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CMD (size = 56 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CMD_DATA struct {
	NWorldHandle                     uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	UseItemMinimonInfo               [50]byte             // offset: 6, size: 50 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_INFO_ZONE_CMD (size = 56 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_INFO_ZONE_CMD_DATA struct {
	NSessionHandle                   uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	UseItemMinimonInfo               [50]byte             // offset: 6, size: 50 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ (size = 1 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ_DATA struct {}

// PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ (size = 24 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ_DATA struct {
	NormalItemList                   [12]uint16           // offset: 0, size: 24 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_NOTICE_CMD (size = 4 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_NOTICE_CMD_DATA struct {
	NErr                             uint16               // offset: 0, size: 2 bytes
	NItemID                          uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD (size = 3 bytes)
type PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD_DATA struct {
	NCharHandle                      uint16               // offset: 0, size: 2 bytes
	BUse                             uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK (size = 18 bytes)
type PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NCoin                            uint64               // offset: 6, size: 8 bytes
	NLockIndex                       uint16               // offset: 14, size: 2 bytes
	NError                           uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_REQ (size = 24 bytes)
type PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_REQ_DATA struct {
	NHeader                          NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemReg                         uint64               // offset: 6, size: 8 bytes
	NCoin                            uint64               // offset: 14, size: 8 bytes
	NLockIndex                       uint16               // offset: 22, size: 2 bytes
}

// PROTO_NC_GAMBLE_COIN_USE_COINITEM_MESSAGE_CMD (size = 8 bytes)
type PROTO_NC_GAMBLE_COIN_USE_COINITEM_MESSAGE_CMD_DATA struct {
	NAddCoin                         uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_ITEMDB_ADMINCREATEFAIL_ACK (size = 6 bytes)
type PROTO_NC_ITEMDB_ADMINCREATEFAIL_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Owner                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_ITEMDB_ADMINCREATESUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_ADMINCREATESUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_ADMINCREATE_REQ (size = 40 bytes)
type PROTO_NC_ITEMDB_ADMINCREATE_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Create                           [38]byte             // offset: 2, size: 38 bytes
}

// PROTO_NC_ITEMDB_BOOTHTRADE_ACK (size = 32 bytes)
type PROTO_NC_ITEMDB_BOOTHTRADE_ACK_DATA struct {
	HeaderSeller                     NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	HeaderBuyer                      NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	BoothlockindexSeller             uint16               // offset: 12, size: 2 bytes
	BoothlockindexBuyer              uint16               // offset: 14, size: 2 bytes
	NBoothTradeType                  uint32               // offset: 16, size: 4 bytes
	NItemKey                         uint64               // offset: 20, size: 8 bytes
	NLotTraded                       uint16               // offset: 28, size: 2 bytes
	Err                              uint16               // offset: 30, size: 2 bytes
}

// PROTO_NC_ITEMDB_BOOTHTRADE_ALL_REQ (size = 68 bytes)
type PROTO_NC_ITEMDB_BOOTHTRADE_ALL_REQ_DATA struct {
	HeaderSeller                     NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	HeaderBuyer                      NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	BoothSeller                      uint32               // offset: 12, size: 4 bytes
	BoothBuyer                       uint32               // offset: 16, size: 4 bytes
	BoothlockindexSeller             uint16               // offset: 20, size: 2 bytes
	BoothlockindexBuyer              uint16               // offset: 22, size: 2 bytes
	IsSellBooth                      uint8                // offset: 24, size: 1 bytes
	Itemregnum                       uint64               // offset: 25, size: 8 bytes
	Itemid                           uint16               // offset: 33, size: 2 bytes
	InvenslotBuyer                   uint8                // offset: 35, size: 1 bytes
	NLot                             uint64               // offset: 36, size: 8 bytes
	Price                            uint64               // offset: 44, size: 8 bytes
	BoothSellerResultMoney           uint64               // offset: 52, size: 8 bytes
	BoothBuyerResultMoney            uint64               // offset: 60, size: 8 bytes
}

// PROTO_NC_ITEMDB_BOOTHTRADE_LOT_REQ (size = 103 bytes)
type PROTO_NC_ITEMDB_BOOTHTRADE_LOT_REQ_DATA struct {
	HeaderSeller                     NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	HeaderBuyer                      NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	BoothSeller                      uint32               // offset: 12, size: 4 bytes
	BoothBuyer                       uint32               // offset: 16, size: 4 bytes
	BoothlockindexSeller             uint16               // offset: 20, size: 2 bytes
	BoothlockindexBuyer              uint16               // offset: 22, size: 2 bytes
	IsSellBooth                      uint8                // offset: 24, size: 1 bytes
	Itemid                           uint16               // offset: 25, size: 2 bytes
	Lotseller                        uint64               // offset: 27, size: 8 bytes
	Lotcreate                        [38]byte             // offset: 35, size: 38 bytes
	Createlot                        uint16               // offset: 73, size: 2 bytes
	Price                            uint64               // offset: 75, size: 8 bytes
	BoothSellerResultLot             uint16               // offset: 83, size: 2 bytes
	BoothBuyerResultLot              uint16               // offset: 85, size: 2 bytes
	BoothSellerResultMoney           uint64               // offset: 87, size: 8 bytes
	BoothBuyerResultMoney            uint64               // offset: 95, size: 8 bytes
}

// PROTO_NC_ITEMDB_BOOTHTRADE_MERGE_REQ (size = 74 bytes)
type PROTO_NC_ITEMDB_BOOTHTRADE_MERGE_REQ_DATA struct {
	HeaderSeller                     NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	HeaderBuyer                      NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	BoothSeller                      uint32               // offset: 12, size: 4 bytes
	BoothBuyer                       uint32               // offset: 16, size: 4 bytes
	BoothlockindexSeller             uint16               // offset: 20, size: 2 bytes
	BoothlockindexBuyer              uint16               // offset: 22, size: 2 bytes
	IsSellBooth                      uint8                // offset: 24, size: 1 bytes
	Itemid                           uint16               // offset: 25, size: 2 bytes
	Lotseller                        uint64               // offset: 27, size: 8 bytes
	Lotbuyer                         uint64               // offset: 35, size: 8 bytes
	Mergelot                         uint16               // offset: 43, size: 2 bytes
	Price                            uint64               // offset: 45, size: 8 bytes
	Allmove                          uint8                // offset: 53, size: 1 bytes
	BoothSellerResultLot             uint16               // offset: 54, size: 2 bytes
	BoothBuyerResultLot              uint16               // offset: 56, size: 2 bytes
	BoothSellerResultMoney           uint64               // offset: 58, size: 8 bytes
	BoothBuyerResultMoney            uint64               // offset: 66, size: 8 bytes
}

// PROTO_NC_ITEMDB_BRACELET_UPGRADE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_BRACELET_UPGRADE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_BRACELET_UPGRADE_REQ (size = 33 bytes)
type PROTO_NC_ITEMDB_BRACELET_UPGRADE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Upgrade_item_regnum              uint64               // offset: 8, size: 8 bytes
	Upgrade_item_id                  uint16               // offset: 16, size: 2 bytes
	Upgrade_item_adj_level           uint8                // offset: 18, size: 1 bytes
	Upgrade_item_destroy             uint8                // offset: 19, size: 1 bytes
	Upgrade_item_new_level           uint8                // offset: 20, size: 1 bytes
	Upgrade_item_fail_count          uint8                // offset: 21, size: 1 bytes
	Raw_regnum                       uint64               // offset: 22, size: 8 bytes
	Raw_item_id                      uint16               // offset: 30, size: 2 bytes
	Raw_lot                          uint8                // offset: 32, size: 1 bytes
}

// PROTO_NC_ITEMDB_BUYALLSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_BUYALLSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_BUYALL_REQ (size = 62 bytes)
type PROTO_NC_ITEMDB_BUYALL_REQ_DATA struct {
	Itembuy                          [38]byte             // offset: 0, size: 38 bytes
	Lockindex                        uint16               // offset: 38, size: 2 bytes
	Price_cen                        uint64               // offset: 40, size: 8 bytes
	Price_fame                       uint32               // offset: 48, size: 4 bytes
	Price_gamblecoin                 uint32               // offset: 52, size: 4 bytes
	UserNo                           uint32               // offset: 56, size: 4 bytes
	Vendor                           uint16               // offset: 60, size: 2 bytes
	Attr                             [0]byte              // offset: 62, flexible array member
}

// PROTO_NC_ITEMDB_BUYCAPSULE_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_BUYCAPSULE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NBuyNpcHandle                    uint16               // offset: 6, size: 2 bytes
	NLockIndex                       uint16               // offset: 8, size: 2 bytes
	NError                           uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_BUYCAPSULE_REQ (size = 252 bytes)
type PROTO_NC_ITEMDB_BUYCAPSULE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NBuyNpcHandle                    uint16               // offset: 6, size: 2 bytes
	NLockIndex                       uint16               // offset: 8, size: 2 bytes
	NMap                             uint16               // offset: 10, size: 2 bytes
	NPrice                           uint64               // offset: 12, size: 8 bytes
	NFame                            uint32               // offset: 20, size: 4 bytes
	NGbCoin                          uint32               // offset: 24, size: 4 bytes
	CapsuleLocation                  uint16               // offset: 28, size: 2 bytes
	CapsuleInfo                      [111]byte            // offset: 30, size: 111 bytes
	ContentInfo                      [111]byte            // offset: 141, size: 111 bytes
}

// PROTO_NC_ITEMDB_BUYLOTFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_BUYLOTFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_BUYLOTSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_BUYLOTSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_BUYLOT_REQ (size = 38 bytes)
type PROTO_NC_ITEMDB_BUYLOT_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Iteminfo                         [12]byte             // offset: 2, size: 12 bytes
	Itemid                           uint16               // offset: 14, size: 2 bytes
	Price_cen                        uint64               // offset: 16, size: 8 bytes
	Price_fame                       uint32               // offset: 24, size: 4 bytes
	Price_gamblecoin                 uint32               // offset: 28, size: 4 bytes
	UserNo                           uint32               // offset: 32, size: 4 bytes
	Vendor                           uint16               // offset: 36, size: 2 bytes
}

// PROTO_NC_ITEMDB_CAPSULEITEM_ACK (size = 129 bytes)
type PROTO_NC_ITEMDB_CAPSULEITEM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	NItemReg                         uint64               // offset: 8, size: 8 bytes
	NLocation                        uint16               // offset: 16, size: 2 bytes
	ContentItemInfo                  [111]byte            // offset: 18, size: 111 bytes
}

// PROTO_NC_ITEMDB_CAPSULEITEM_REQ (size = 16 bytes)
type PROTO_NC_ITEMDB_CAPSULEITEM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemReg                         uint64               // offset: 6, size: 8 bytes
	NLocation                        uint16               // offset: 14, size: 2 bytes
}

// PROTO_NC_ITEMDB_CHARGED_LIST_ACK (size = 11 bytes)
type PROTO_NC_ITEMDB_CHARGED_LIST_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ErrorCode                        uint16               // offset: 6, size: 2 bytes
	NPartMark                        uint8                // offset: 8, size: 1 bytes
	NumOfChargedItem                 uint16               // offset: 9, size: 2 bytes
	ChargedItemInfoList              [0]byte              // offset: 11, flexible array member
}

// PROTO_NC_ITEMDB_CHARGED_LIST_REQ (size = 14 bytes)
type PROTO_NC_ITEMDB_CHARGED_LIST_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NUserNo                          uint32               // offset: 6, size: 4 bytes
	From                             uint16               // offset: 10, size: 2 bytes
	To                               uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_ITEMDB_CHARGED_WITHDRAW_ACK (size = 26 bytes)
type PROTO_NC_ITEMDB_CHARGED_WITHDRAW_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	ChargedItemInfo                  [16]byte             // offset: 8, size: 16 bytes
	ErrorCode                        uint16               // offset: 24, size: 2 bytes
}

// PROTO_NC_ITEMDB_CHARGED_WITHDRAW_REQ (size = 29 bytes)
type PROTO_NC_ITEMDB_CHARGED_WITHDRAW_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NUserNo                          uint32               // offset: 8, size: 4 bytes
	ChargedItemInfo                  [16]byte             // offset: 12, size: 16 bytes
	NumOfItemCreate                  uint8                // offset: 28, size: 1 bytes
	ItemCreate                       [0]byte              // offset: 29, flexible array member
}

// PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_ACK (size = 23 bytes)
type PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NLockIndex                       uint16               // offset: 2, size: 2 bytes
	NCharRegNum                      uint32               // offset: 4, size: 4 bytes
	NItemKey                         uint64               // offset: 8, size: 8 bytes
	NItemID                          uint16               // offset: 16, size: 2 bytes
	EChatColorItemType               uint32               // offset: 18, size: 4 bytes
	NChangeColor                     uint8                // offset: 22, size: 1 bytes
}

// PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_REQ (size = 25 bytes)
type PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NItemKey                         uint64               // offset: 6, size: 8 bytes
	NItemID                          uint16               // offset: 14, size: 2 bytes
	NItemLot                         uint32               // offset: 16, size: 4 bytes
	EChatColorItemType               uint32               // offset: 20, size: 4 bytes
	NChangeColor                     uint8                // offset: 24, size: 1 bytes
}

// PROTO_NC_ITEMDB_CHESTITEM_ACK (size = 19 bytes)
type PROTO_NC_ITEMDB_CHESTITEM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Chest                            uint64               // offset: 6, size: 8 bytes
	Location                         uint16               // offset: 14, size: 2 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
	Item                             uint8                // offset: 18, size: 1 bytes
}

// PROTO_NC_ITEMDB_CHESTITEM_REQ (size = 16 bytes)
type PROTO_NC_ITEMDB_CHESTITEM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Chest                            uint64               // offset: 6, size: 8 bytes
	Location                         uint16               // offset: 14, size: 2 bytes
}

// PROTO_NC_ITEMDB_CLASS_CHANGE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_CLASS_CHANGE_ACK_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NHandleNum                       uint16               // offset: 4, size: 2 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_CLASS_CHANGE_REQ (size = 41 bytes)
type PROTO_NC_ITEMDB_CLASS_CHANGE_REQ_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NHandleNum                       uint16               // offset: 4, size: 2 bytes
	NStuffKey                        uint64               // offset: 6, size: 8 bytes
	NStuffID                         uint16               // offset: 14, size: 2 bytes
	NStuffLot                        uint32               // offset: 16, size: 4 bytes
	NMaxStat                         uint8                // offset: 20, size: 1 bytes
	NFreeStat                        [5]byte              // offset: 21, size: 5 bytes
	NNewClass                        uint8                // offset: 26, size: 1 bytes
	NTotalPrice                      uint64               // offset: 27, size: 8 bytes
	Lockindex                        uint16               // offset: 35, size: 2 bytes
	NDeleteSkillCount                uint32               // offset: 37, size: 4 bytes
	NDeleteSkillID                   [0]byte              // offset: 41, flexible array member
}

// PROTO_NC_ITEMDB_CREATEITEMLISTFAIL_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_CREATEITEMLISTFAIL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Error                            uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_CREATEITEMLISTSUC_ACK (size = 8 bytes)
type PROTO_NC_ITEMDB_CREATEITEMLISTSUC_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_ITEMDB_CREATEITEMLIST_REQ (size = 23 bytes)
type PROTO_NC_ITEMDB_CREATEITEMLIST_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Owner                            uint32               // offset: 6, size: 4 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	MapField                         uint16               // offset: 12, size: 2 bytes
	Price                            uint64               // offset: 14, size: 8 bytes
	Itemlist                         uint8                // offset: 22, size: 1 bytes
}

// PROTO_NC_ITEMDB_DEPOSIT_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_DEPOSIT_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Err                              uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_DEPOSIT_REQ (size = 20 bytes)
type PROTO_NC_ITEMDB_DEPOSIT_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Usrregnum                        uint32               // offset: 8, size: 4 bytes
	Cen                              uint64               // offset: 12, size: 8 bytes
}

// PROTO_NC_ITEMDB_DESTROY_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_DESTROY_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_DISMANTLE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_DISMANTLE_REQ (size = 26 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Raw                              [17]byte             // offset: 8, size: 17 bytes
	Producenumber                    uint8                // offset: 25, size: 1 bytes
	Produce                          [0]byte              // offset: 26, flexible array member
}

// PROTO_NC_ITEMDB_DROPALLSUC_ACK (size = 111 bytes)
type PROTO_NC_ITEMDB_DROPALLSUC_ACK_DATA struct {
	Dropitemstr                      [111]byte            // offset: 0, size: 111 bytes
}

// PROTO_NC_ITEMDB_DROPALL_REQ (size = 123 bytes)
type PROTO_NC_ITEMDB_DROPALL_REQ_DATA struct {
	Itemid                           uint16               // offset: 0, size: 2 bytes
	Dropitem                         [10]byte             // offset: 2, size: 10 bytes
	Dropitemstr                      [111]byte            // offset: 12, size: 111 bytes
}

// PROTO_NC_ITEMDB_DROPLOT_REQ (size = 30 bytes)
type PROTO_NC_ITEMDB_DROPLOT_REQ_DATA struct {
	Itemid                           uint16               // offset: 0, size: 2 bytes
	Iteminfo                         [12]byte             // offset: 2, size: 12 bytes
	Dropitem                         [12]byte             // offset: 14, size: 12 bytes
	NCurLot                          uint32               // offset: 26, size: 4 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	ItemLockIndex                    uint16               // offset: 8, size: 2 bytes
	Equipment                        [6]byte              // offset: 10, size: 6 bytes
	Gem                              uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ (size = 34 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemLockIndex                    uint16               // offset: 6, size: 2 bytes
	Equipment                        [12]byte             // offset: 8, size: 12 bytes
	Gem                              [12]byte             // offset: 20, size: 12 bytes
	GemSocketIndex                   uint8                // offset: 32, size: 1 bytes
	GemLot                           uint8                // offset: 33, size: 1 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK (size = 23 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	ItemLockIndex                    uint16               // offset: 8, size: 2 bytes
	Equipment                        [3]byte              // offset: 10, size: 3 bytes
	Drill                            [10]byte             // offset: 13, size: 10 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ (size = 32 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemLockIndex                    uint16               // offset: 6, size: 2 bytes
	Equipment                        [12]byte             // offset: 8, size: 12 bytes
	Drill                            [12]byte             // offset: 20, size: 12 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK (size = 26 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	ItemLockIndex                    uint16               // offset: 8, size: 2 bytes
	Equipment                        uint16               // offset: 10, size: 2 bytes
	Remover                          [10]byte             // offset: 12, size: 10 bytes
	GemSlot                          [3]byte              // offset: 22, size: 3 bytes
	NeedItem                         uint8                // offset: 25, size: 1 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ (size = 36 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemLockIndex                    uint16               // offset: 6, size: 2 bytes
	Equipment                        [12]byte             // offset: 8, size: 12 bytes
	Remover                          [12]byte             // offset: 20, size: 12 bytes
	GemSlot                          [3]byte              // offset: 32, size: 3 bytes
	NeedItem                         uint8                // offset: 35, size: 1 bytes
}

// PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ (size = 24 bytes)
type PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemRegistNumber                 uint64               // offset: 6, size: 8 bytes
	SocketCount                      uint8                // offset: 14, size: 1 bytes
	Sockets                          [9]byte              // offset: 15, size: 9 bytes
}

// PROTO_NC_ITEMDB_EQUIPFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_EQUIPFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_EQUIPSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_EQUIPSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_ACK (size = 34 bytes)
type PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
	NBelongedItem_ItemKey            uint64               // offset: 10, size: 8 bytes
	NBelongedItem_ItemID             uint16               // offset: 18, size: 2 bytes
	NCancelItem_ItemKey              uint64               // offset: 20, size: 8 bytes
	NCancelItem_ItemID               uint16               // offset: 28, size: 2 bytes
	NCancelItem_LotCnt               uint32               // offset: 30, size: 4 bytes
}

// PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_REQ (size = 32 bytes)
type PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NBelongedItem_ItemKey            uint64               // offset: 8, size: 8 bytes
	NBelongedItem_ItemID             uint16               // offset: 16, size: 2 bytes
	NCancelItem_ItemKey              uint64               // offset: 18, size: 8 bytes
	NCancelItem_ItemID               uint16               // offset: 26, size: 2 bytes
	NCancelItem_LotCnt               uint32               // offset: 28, size: 4 bytes
}

// PROTO_NC_ITEMDB_EQUIP_REQ (size = 29 bytes)
type PROTO_NC_ITEMDB_EQUIP_REQ_DATA struct {
	RelocA                           [10]byte             // offset: 0, size: 10 bytes
	RelocA_ItemID                    uint16               // offset: 10, size: 2 bytes
	RelocA_IsPutOnBelonged           uint32               // offset: 12, size: 4 bytes
	RelocB                           [10]byte             // offset: 16, size: 10 bytes
	Lockindex                        uint16               // offset: 26, size: 2 bytes
	Unequipnum                       uint8                // offset: 28, size: 1 bytes
	Unequip                          [0]byte              // offset: 29, flexible array member
}

// PROTO_NC_ITEMDB_EXCHANGEFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_EXCHANGEFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_EXCHANGESUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_EXCHANGESUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_EXCHANGE_REQ (size = 42 bytes)
type PROTO_NC_ITEMDB_EXCHANGE_REQ_DATA struct {
	OwnerA                           uint32               // offset: 0, size: 4 bytes
	OwnerB                           uint32               // offset: 4, size: 4 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	RelocA                           [10]byte             // offset: 10, size: 10 bytes
	RelocB                           [10]byte             // offset: 20, size: 10 bytes
	RelocA_ID                        uint16               // offset: 30, size: 2 bytes
	RelocB_ID                        uint16               // offset: 32, size: 2 bytes
	RelocA_Lot                       uint32               // offset: 34, size: 4 bytes
	RelocB_Lot                       uint32               // offset: 38, size: 4 bytes
}

// PROTO_NC_ITEMDB_FAIL_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_FAIL_ACK struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_FURNITURE_ENDURE_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_FURNITURE_ENDURE_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	NFurnitureHnd                    uint16               // offset: 2, size: 2 bytes
	NItemKey                         uint64               // offset: 4, size: 8 bytes
	DNewEndureTime                   uint32               // offset: 12, size: 4 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_FURNITURE_ENDURE_REQ (size = 35 bytes)
type PROTO_NC_ITEMDB_FURNITURE_ENDURE_REQ_DATA struct {
	NFurnitureHnd                    uint16               // offset: 0, size: 2 bytes
	NItemKey                         uint64               // offset: 2, size: 8 bytes
	DNewEndureTime                   uint32               // offset: 10, size: 4 bytes
	CommonData                       [21]byte             // offset: 14, size: 21 bytes
}

// PROTO_NC_ITEMDB_GETFROMCAPSULE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_GETFROMCAPSULE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_GETFROMCAPSULE_REQ (size = 30 bytes)
type PROTO_NC_ITEMDB_GETFROMCAPSULE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NOwner                           uint32               // offset: 6, size: 4 bytes
	NLockIndex                       uint16               // offset: 10, size: 2 bytes
	NItmeReg                         uint64               // offset: 12, size: 8 bytes
	NReloc                           [10]byte             // offset: 20, size: 10 bytes
}

// PROTO_NC_ITEMDB_GETFROMCHESTFAIL_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_GETFROMCHESTFAIL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_GETFROMCHESTSUC_ACK (size = 8 bytes)
type PROTO_NC_ITEMDB_GETFROMCHESTSUC_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_ITEMDB_GETFROMCHEST_REQ (size = 21 bytes)
type PROTO_NC_ITEMDB_GETFROMCHEST_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Owner                            uint32               // offset: 6, size: 4 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	Chest                            uint64               // offset: 12, size: 8 bytes
	Relocnum                         uint8                // offset: 20, size: 1 bytes
	Reloc                            [0]byte              // offset: 21, flexible array member
}

// PROTO_NC_ITEMDB_INC_DEC_MONEY_REQ (size = 9 bytes)
type PROTO_NC_ITEMDB_INC_DEC_MONEY_REQ_DATA struct {
	Owner                            uint32               // offset: 0, size: 4 bytes
	Cen                              uint32               // offset: 4, size: 4 bytes
	Positive                         uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_ITEMDB_ITEMBREAKFAIL_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_ITEMBREAKFAIL_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	ErrorCode                        uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMBREAKSUC_ACK (size = 8 bytes)
type PROTO_NC_ITEMDB_ITEMBREAKSUC_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMBREAK_REQ (size = 25 bytes)
type PROTO_NC_ITEMDB_ITEMBREAK_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	BreakItemKey                     uint64               // offset: 8, size: 8 bytes
	Owner                            uint32               // offset: 16, size: 4 bytes
	NID                              uint16               // offset: 20, size: 2 bytes
	NLot                             uint16               // offset: 22, size: 2 bytes
	NDeleteChestItem                 uint8                // offset: 24, size: 1 bytes
	DeleteChestItem                  [0]byte              // offset: 25, flexible array member
}

// PROTO_NC_ITEMDB_ITEMCHANGE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_ITEMCHANGE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMCHANGE_REQ (size = 16 bytes)
type PROTO_NC_ITEMDB_ITEMCHANGE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Itemkey                          uint64               // offset: 8, size: 8 bytes
	Itemstruct                       [0]byte              // offset: 16, flexible array member
}

// PROTO_NC_ITEMDB_ITEMINFO_ACK (size = 121 bytes)
type PROTO_NC_ITEMDB_ITEMINFO_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	FailProtocol                     uint16               // offset: 8, size: 2 bytes
	ItemInfo                         [111]byte            // offset: 10, size: 111 bytes
}

// PROTO_NC_ITEMDB_ITEMINFO_REQ (size = 16 bytes)
type PROTO_NC_ITEMDB_ITEMINFO_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	ItemKey                          uint64               // offset: 6, size: 8 bytes
	FailProtocol                     uint16               // offset: 14, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ (size = 34 bytes)
type PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NVendor                          uint16               // offset: 8, size: 2 bytes
	NUserNo                          uint32               // offset: 10, size: 4 bytes
	NItemMoneyKey                    uint64               // offset: 14, size: 8 bytes
	NItemMoneyID                     uint16               // offset: 22, size: 2 bytes
	NDecreaseLot                     uint32               // offset: 24, size: 4 bytes
	NMoneyTotalLot                   uint32               // offset: 28, size: 4 bytes
	NDeleteItemCount                 uint8                // offset: 32, size: 1 bytes
	NCreateItemCount                 uint8                // offset: 33, size: 1 bytes
	CreateItem                       [0]byte              // offset: 34, flexible array member
}

// PROTO_NC_ITEMDB_ITEMMONEY_BUYLOT_REQ (size = 47 bytes)
type PROTO_NC_ITEMDB_ITEMMONEY_BUYLOT_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NVendor                          uint16               // offset: 8, size: 2 bytes
	NUserNo                          uint32               // offset: 10, size: 4 bytes
	NItemKey                         uint64               // offset: 14, size: 8 bytes
	NItemID                          uint16               // offset: 22, size: 2 bytes
	NAddLot                          uint32               // offset: 24, size: 4 bytes
	NMoneyTotalLot                   uint32               // offset: 28, size: 4 bytes
	NItemMoneyKey                    uint64               // offset: 32, size: 8 bytes
	NItemMoneyID                     uint16               // offset: 40, size: 2 bytes
	NDecreaseLot                     uint32               // offset: 42, size: 4 bytes
	NDeleteItemCount                 uint8                // offset: 46, size: 1 bytes
	NDeleteItemInfo                  [0]byte              // offset: 47, flexible array member
}

// PROTO_NC_ITEMDB_ITEMMONEY_BUY_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_ITEMMONEY_BUY_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMREBUILD_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_ITEMREBUILD_ACK_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMREBUILD_REQ (size = 54 bytes)
type PROTO_NC_ITEMDB_ITEMREBUILD_REQ_DATA struct {
	Clientregnum                     uint32               // offset: 0, size: 4 bytes
	Clientzonehandle                 uint16               // offset: 4, size: 2 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	EraseItemKey                     uint64               // offset: 8, size: 8 bytes
	CreateItem                       [38]byte             // offset: 16, size: 38 bytes
}

// PROTO_NC_ITEMDB_ITEMTOTALINFORM_ACK (size = 111 bytes)
type PROTO_NC_ITEMDB_ITEMTOTALINFORM_ACK_DATA struct {
	Itemtotalinfo                    [111]byte            // offset: 0, size: 111 bytes
}

// PROTO_NC_ITEMDB_ITEMTOTALINFORM_REQ (size = 8 bytes)
type PROTO_NC_ITEMDB_ITEMTOTALINFORM_REQ_DATA struct {
	Itemregnum                       uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_ITEMDB_ITEMTRADEFAIL_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_ITEMTRADEFAIL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMTRADESUC_ACK (size = 8 bytes)
type PROTO_NC_ITEMDB_ITEMTRADESUC_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_ITEMDB_ITEMTRADE_REQ (size = 42 bytes)
type PROTO_NC_ITEMDB_ITEMTRADE_REQ_DATA struct {
	A                                NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	B                                NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	LockindexA                       uint16               // offset: 12, size: 2 bytes
	LockindexB                       uint16               // offset: 14, size: 2 bytes
	Cen                              [25]byte             // offset: 16, size: 25 bytes
	Tradenumber                      uint8                // offset: 41, size: 1 bytes
	Tradeitem                        [0]byte              // offset: 42, flexible array member
}

// PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_ACK (size = 28 bytes)
type PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	MapName                          Name3                // offset: 8, size: 12 bytes
	SkillID                          uint16               // offset: 20, size: 2 bytes
	SkillCoolTime                    uint32               // offset: 22, size: 4 bytes
	NError                           uint16               // offset: 26, size: 2 bytes
}

// PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_REQ (size = 40 bytes)
type PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	MapName                          Name3                // offset: 8, size: 12 bytes
	SkillID                          uint16               // offset: 20, size: 2 bytes
	SkillCoolTime                    uint32               // offset: 22, size: 4 bytes
	NStuffKey                        uint64               // offset: 26, size: 8 bytes
	NStuffID                         uint16               // offset: 34, size: 2 bytes
	NStuffLot                        uint32               // offset: 36, size: 4 bytes
}

// PROTO_NC_ITEMDB_MERGEFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_MERGEFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_MERGESUC_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_MERGESUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Lotmain                          uint64               // offset: 2, size: 8 bytes
	MergeResultLot                   uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_MERGE_REQ (size = 36 bytes)
type PROTO_NC_ITEMDB_MERGE_REQ_DATA struct {
	OwnerA                           uint32               // offset: 0, size: 4 bytes
	OwnerB                           uint32               // offset: 4, size: 4 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	Itemid                           uint16               // offset: 10, size: 2 bytes
	Lotmain                          uint64               // offset: 12, size: 8 bytes
	Lotsub                           uint64               // offset: 20, size: 8 bytes
	NMainLot                         uint32               // offset: 28, size: 4 bytes
	NSubLot                          uint32               // offset: 32, size: 4 bytes
}

// PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NEffectHandle                    uint16               // offset: 8, size: 2 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_REQ (size = 32 bytes)
type PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NDemandItemRegNum                uint64               // offset: 8, size: 8 bytes
	NDemandItemID                    uint16               // offset: 16, size: 2 bytes
	NNeedItemLot                     uint32               // offset: 18, size: 4 bytes
	NDemandMoney                     uint64               // offset: 22, size: 8 bytes
	NEffectHandle                    uint16               // offset: 30, size: 2 bytes
}

// PROTO_NC_ITEMDB_MINIHOUSE_PORTAL_EFFECT_DEMANDGOOD_REQ (size = 32 bytes)
type PROTO_NC_ITEMDB_MINIHOUSE_PORTAL_EFFECT_DEMANDGOOD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NDemandItemRegNum                uint64               // offset: 8, size: 8 bytes
	NDemandItemID                    uint16               // offset: 16, size: 2 bytes
	NNeedItemLot                     uint32               // offset: 18, size: 4 bytes
	NDemandMoney                     uint64               // offset: 22, size: 8 bytes
	NEffectHandle                    uint16               // offset: 30, size: 2 bytes
}

// PROTO_NC_ITEMDB_MIX_ITEM_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_MIX_ITEM_ACK_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NHandleNum                       uint16               // offset: 4, size: 2 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_MIX_ITEM_REQ (size = 163 bytes)
type PROTO_NC_ITEMDB_MIX_ITEM_REQ_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NHandleNum                       uint16               // offset: 4, size: 2 bytes
	NSubRawKey                       uint64               // offset: 6, size: 8 bytes
	NSubRawID                        uint16               // offset: 14, size: 2 bytes
	NSubRawLot                       uint32               // offset: 16, size: 4 bytes
	NRawLeftKey                      uint64               // offset: 20, size: 8 bytes
	NRawLeftID                       uint16               // offset: 28, size: 2 bytes
	NRawLeftLot                      uint32               // offset: 30, size: 4 bytes
	NRawRightKey                     uint64               // offset: 34, size: 8 bytes
	NRawRightID                      uint16               // offset: 42, size: 2 bytes
	NRawRightLot                     uint32               // offset: 44, size: 4 bytes
	MixItemRoc                       uint16               // offset: 48, size: 2 bytes
	MixItemTotalInfo                 [111]byte            // offset: 50, size: 111 bytes
	NLockIndex                       uint16               // offset: 161, size: 2 bytes
}

// PROTO_NC_ITEMDB_MOVER_RAREMOVER_ACK (size = 114 bytes)
type PROTO_NC_ITEMDB_MOVER_RAREMOVER_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
	NSlot                            uint8                // offset: 10, size: 1 bytes
	NItemID                          uint16               // offset: 11, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 13, size: 101 bytes
}

// PROTO_NC_ITEMDB_MOVER_RAREMOVER_REQ (size = 343 bytes)
type PROTO_NC_ITEMDB_MOVER_RAREMOVER_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NMainMoverItemKey                uint64               // offset: 8, size: 8 bytes
	NMainMoverItemID                 uint16               // offset: 16, size: 2 bytes
	MainMoverAttr                    [101]byte            // offset: 18, size: 101 bytes
	NSubMoverItemKey                 uint64               // offset: 119, size: 8 bytes
	NSubMoverItemID                  uint16               // offset: 127, size: 2 bytes
	SubMoverAttr                     [101]byte            // offset: 129, size: 101 bytes
	NRareMoverItemKey                uint64               // offset: 230, size: 8 bytes
	NRareMoverItemID                 uint16               // offset: 238, size: 2 bytes
	RareMoverLoc                     uint16               // offset: 240, size: 2 bytes
	RareMoverAttr                    [101]byte            // offset: 242, size: 101 bytes
}

// PROTO_NC_ITEMDB_MOVER_UPGRADE_ACK (size = 115 bytes)
type PROTO_NC_ITEMDB_MOVER_UPGRADE_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
	BUpgrade                         uint8                // offset: 10, size: 1 bytes
	NSlot                            uint8                // offset: 11, size: 1 bytes
	NItemID                          uint16               // offset: 12, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 14, size: 101 bytes
}

// PROTO_NC_ITEMDB_MOVER_UPGRADE_REQ (size = 233 bytes)
type PROTO_NC_ITEMDB_MOVER_UPGRADE_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	BUpgrade                         uint8                // offset: 8, size: 1 bytes
	NSubMoverItemKey                 uint64               // offset: 9, size: 8 bytes
	NSubMoverItemID                  uint16               // offset: 17, size: 2 bytes
	SubMoverAttr                     [101]byte            // offset: 19, size: 101 bytes
	NMainMoverItemKey                uint64               // offset: 120, size: 8 bytes
	NMainMoverItemID                 uint16               // offset: 128, size: 2 bytes
	MainMoverLoc                     uint16               // offset: 130, size: 2 bytes
	MainMoverAttr                    [101]byte            // offset: 132, size: 101 bytes
}

// PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemLockIndex                    uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_REQ (size = 19 bytes)
type PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	ItemLockIndex                    uint16               // offset: 6, size: 2 bytes
	UseVaultItemRegNo                uint64               // offset: 8, size: 8 bytes
	UseVaultItemID                   uint16               // offset: 16, size: 2 bytes
	Itemcreatecount                  uint8                // offset: 18, size: 1 bytes
	Itemcreate                       [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_ITEMDB_NEW_UPGRADE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_NEW_UPGRADE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_NEW_UPGRADE_REQ (size = 36 bytes)
type PROTO_NC_ITEMDB_NEW_UPGRADE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Upgrade_item_regnum              uint64               // offset: 8, size: 8 bytes
	Upgrade_item_id                  uint16               // offset: 16, size: 2 bytes
	Upgrade_item_adj_level           uint8                // offset: 18, size: 1 bytes
	Upgrade_item_new_level           uint8                // offset: 19, size: 1 bytes
	Upgrade_item_fail_count          uint8                // offset: 20, size: 1 bytes
	Raw_regnum                       uint64               // offset: 21, size: 8 bytes
	Raw_item_id                      uint16               // offset: 29, size: 2 bytes
	Raw_lot                          uint8                // offset: 31, size: 1 bytes
	Is_randomoption_enable           uint8                // offset: 32, size: 1 bytes
	Upgrade_randomoption             [3]byte              // offset: 33, size: 3 bytes
}

// PROTO_NC_ITEMDB_OPENSTORAGE_ACK (size = 17 bytes)
type PROTO_NC_ITEMDB_OPENSTORAGE_ACK_DATA struct {
	Zonepackheader                   NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Curpage                          uint8                // offset: 6, size: 1 bytes
	NOpenType                        uint8                // offset: 7, size: 1 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Itemcounter                      uint8                // offset: 16, size: 1 bytes
	Itemarray                        [0]byte              // offset: 17, flexible array member
}

// PROTO_NC_ITEMDB_OPENSTORAGE_FAIL_ACK (size = 8 bytes)
type PROTO_NC_ITEMDB_OPENSTORAGE_FAIL_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_ITEMDB_OPENSTORAGE_REQ (size = 12 bytes)
type PROTO_NC_ITEMDB_OPENSTORAGE_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Page                             uint8                // offset: 6, size: 1 bytes
	Usrregnum                        uint32               // offset: 7, size: 4 bytes
	NOpenType                        uint8                // offset: 11, size: 1 bytes
}

// PROTO_NC_ITEMDB_PICKMERGE_REQ (size = 26 bytes)
type PROTO_NC_ITEMDB_PICKMERGE_REQ_DATA struct {
	Itemid                           uint16               // offset: 0, size: 2 bytes
	Itemonfield                      uint64               // offset: 2, size: 8 bytes
	NCellItemLot                     uint32               // offset: 10, size: 4 bytes
	Itempick                         [12]byte             // offset: 14, size: 12 bytes
}

// PROTO_NC_ITEMDB_PICKMONEY_REQ (size = 16 bytes)
type PROTO_NC_ITEMDB_PICKMONEY_REQ_DATA struct {
	Owner                            uint32               // offset: 0, size: 4 bytes
	Itemonfield                      uint64               // offset: 4, size: 8 bytes
	Cen                              uint32               // offset: 12, size: 4 bytes
}

// PROTO_NC_ITEMDB_PRODUCE_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_PRODUCE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Product                          uint16               // offset: 8, size: 2 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_PRODUCE_REQ (size = 50 bytes)
type PROTO_NC_ITEMDB_PRODUCE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Product                          uint16               // offset: 8, size: 2 bytes
	Produce                          [39]byte             // offset: 10, size: 39 bytes
	Rawnumber                        uint8                // offset: 49, size: 1 bytes
	Raw                              [0]byte              // offset: 50, flexible array member
}

// PROTO_NC_ITEMDB_QUESTREWARD_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_QUESTREWARD_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NQuestID                         uint16               // offset: 8, size: 2 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_QUESTREWARD_REQ (size = 45 bytes)
type PROTO_NC_ITEMDB_QUESTREWARD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	NQuestID                         uint16               // offset: 8, size: 2 bytes
	Money                            uint32               // offset: 10, size: 4 bytes
	Exp                              uint32               // offset: 14, size: 4 bytes
	Fame                             uint32               // offset: 18, size: 4 bytes
	Current_money                    uint64               // offset: 22, size: 8 bytes
	Current_exp                      uint64               // offset: 30, size: 8 bytes
	Current_fame                     uint32               // offset: 38, size: 4 bytes
	Ct_type                          uint8                // offset: 42, size: 1 bytes
	Ct_elementno                     uint8                // offset: 43, size: 1 bytes
	Itemcreatecount                  uint8                // offset: 44, size: 1 bytes
	Itemcreate                       [0]byte              // offset: 45, flexible array member
}

// PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_ACK (size = 21 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
	NItemKey                         uint64               // offset: 10, size: 8 bytes
	NItemID                          uint16               // offset: 18, size: 2 bytes
	NItemSlotNo                      uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_REQ (size = 134 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NStuffKey                        uint64               // offset: 8, size: 8 bytes
	NStuffID                         uint16               // offset: 16, size: 2 bytes
	NStuffLot                        uint32               // offset: 18, size: 4 bytes
	NItemKey                         uint64               // offset: 22, size: 8 bytes
	NItemID                          uint16               // offset: 30, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 32, size: 101 bytes
	NItemSlotNo                      uint8                // offset: 133, size: 1 bytes
}

// PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_REQ (size = 220 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NItemKey                         uint64               // offset: 8, size: 8 bytes
	NItemID                          uint16               // offset: 16, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 18, size: 101 bytes
	ItemAttrOrig                     [101]byte            // offset: 119, size: 101 bytes
}

// PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ (size = 146 bytes)
type PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NItemKey                         uint64               // offset: 8, size: 8 bytes
	NItemID                          uint16               // offset: 16, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 18, size: 101 bytes
	NROCount_before                  uint32               // offset: 119, size: 4 bytes
	NROCount_after                   uint32               // offset: 123, size: 4 bytes
	NStuffID                         uint16               // offset: 127, size: 2 bytes
	NStuffTotalDecreaseLot           uint32               // offset: 129, size: 4 bytes
	NStuffKey                        uint64               // offset: 133, size: 8 bytes
	NDecreaseLot                     uint32               // offset: 141, size: 4 bytes
	NDeleteStuffCount                uint8                // offset: 145, size: 1 bytes
	NDeleteStuffKey                  [0]byte              // offset: 146, flexible array member
}

// PROTO_NC_ITEMDB_REINFORCEUNEQUIPSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_REINFORCEUNEQUIPSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_REINFORCEUNEQUIP_REQ (size = 14 bytes)
type PROTO_NC_ITEMDB_REINFORCEUNEQUIP_REQ_DATA struct {
	Unequ                            [10]byte             // offset: 0, size: 10 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	Equloc                           uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_ITEMDB_RELOCFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_RELOCFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_RELOCSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_RELOCSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_RELOC_REQ (size = 28 bytes)
type PROTO_NC_ITEMDB_RELOC_REQ_DATA struct {
	OwnerA                           uint32               // offset: 0, size: 4 bytes
	OwnerB                           uint32               // offset: 4, size: 4 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	Item                             [10]byte             // offset: 10, size: 10 bytes
	Inven_from                       uint16               // offset: 20, size: 2 bytes
	NFromItemID                      uint16               // offset: 22, size: 2 bytes
	NFromItemLot                     uint32               // offset: 24, size: 4 bytes
}

// PROTO_NC_ITEMDB_REPURCHASE_ALL_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_REPURCHASE_ALL_ACK_DATA struct {
	NHandleNum                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	NLockIndex                       uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_REPURCHASE_ALL_REQ (size = 64 bytes)
type PROTO_NC_ITEMDB_REPURCHASE_ALL_REQ_DATA struct {
	NHandleNum                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	RepurchaseItem                   [38]byte             // offset: 8, size: 38 bytes
	NPrice                           uint32               // offset: 46, size: 4 bytes
	NSellItemKey                     uint64               // offset: 50, size: 8 bytes
	UserNo                           uint32               // offset: 58, size: 4 bytes
	Vendor                           uint16               // offset: 62, size: 2 bytes
}

// PROTO_NC_ITEMDB_REPURCHASE_LOT_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_REPURCHASE_LOT_ACK_DATA struct {
	NHandleNum                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	NLockIndex                       uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_REPURCHASE_LOT_REQ (size = 40 bytes)
type PROTO_NC_ITEMDB_REPURCHASE_LOT_REQ_DATA struct {
	NHandleNum                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NLockIndex                       uint16               // offset: 6, size: 2 bytes
	NItemID                          uint16               // offset: 8, size: 2 bytes
	NItemAttrChange                  [12]byte             // offset: 10, size: 12 bytes
	NPrice                           uint32               // offset: 22, size: 4 bytes
	NSellItemKey                     uint64               // offset: 26, size: 8 bytes
	UserNo                           uint32               // offset: 34, size: 4 bytes
	Vendor                           uint16               // offset: 38, size: 2 bytes
}

// PROTO_NC_ITEMDB_RESET_SCROLL_LINK_MAP_INFO_CMD (size = 4 bytes)
type PROTO_NC_ITEMDB_RESET_SCROLL_LINK_MAP_INFO_CMD_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_ITEMDB_SELLALLSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_SELLALLSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_SELLALL_REQ (size = 33 bytes)
type PROTO_NC_ITEMDB_SELLALL_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Itemsell                         [10]byte             // offset: 2, size: 10 bytes
	Price                            uint32               // offset: 12, size: 4 bytes
	Itemid                           uint16               // offset: 16, size: 2 bytes
	UserNo                           uint32               // offset: 18, size: 4 bytes
	Vendor                           uint16               // offset: 22, size: 2 bytes
	NLot                             uint64               // offset: 24, size: 8 bytes
	NDeleteChestItem                 uint8                // offset: 32, size: 1 bytes
	DeleteChestItem                  [0]byte              // offset: 33, flexible array member
}

// PROTO_NC_ITEMDB_SELLLOTFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_SELLLOTFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_SELLLOTSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_SELLLOTSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_SELLLOT_REQ (size = 26 bytes)
type PROTO_NC_ITEMDB_SELLLOT_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Iteminfo                         [12]byte             // offset: 2, size: 12 bytes
	Itemid                           uint16               // offset: 14, size: 2 bytes
	Price                            uint32               // offset: 16, size: 4 bytes
	UserNo                           uint32               // offset: 20, size: 4 bytes
	Vendor                           uint16               // offset: 24, size: 2 bytes
}

// PROTO_NC_ITEMDB_SHIELDENDURESET_CMD (size = 12 bytes)
type PROTO_NC_ITEMDB_SHIELDENDURESET_CMD_DATA struct {
	ShieldKey                        uint64               // offset: 0, size: 8 bytes
	ChargedNumber                    uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_ACK (size = 24 bytes)
type PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	ShieldKey                        uint64               // offset: 8, size: 8 bytes
	ShieldSlot                       uint16               // offset: 16, size: 2 bytes
	ChargedNumber                    uint32               // offset: 18, size: 4 bytes
	Error                            uint16               // offset: 22, size: 2 bytes
}

// PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_REQ (size = 35 bytes)
type PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_REQ_DATA struct {
	ShieldKey                        uint64               // offset: 0, size: 8 bytes
	ShieldSlot                       uint16               // offset: 8, size: 2 bytes
	ChargedNumber                    uint32               // offset: 10, size: 4 bytes
	CommonData                       [21]byte             // offset: 14, size: 21 bytes
}

// PROTO_NC_ITEMDB_SOULSTONEBUYFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_SOULSTONEBUYFAIL_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_SOULSTONEBUYSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_SOULSTONEBUYSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_SOULSTONEBUY_REQ (size = 17 bytes)
type PROTO_NC_ITEMDB_SOULSTONEBUY_REQ_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Flag                             uint8                // offset: 4, size: 1 bytes
	Number                           uint16               // offset: 5, size: 2 bytes
	Price                            uint64               // offset: 7, size: 8 bytes
	Lockindex                        uint16               // offset: 15, size: 2 bytes
}

// PROTO_NC_ITEMDB_SPLITFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_SPLITFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_SPLITSUC_ACK (size = 12 bytes)
type PROTO_NC_ITEMDB_SPLITSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Lotmain                          uint64               // offset: 2, size: 8 bytes
	RemainLot                        uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_SPLIT_N_MERGEFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_SPLIT_N_MERGEFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_SPLIT_N_MERGESUC_ACK (size = 24 bytes)
type PROTO_NC_ITEMDB_SPLIT_N_MERGESUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	FromItemKey                      uint64               // offset: 2, size: 8 bytes
	ToItemKey                        uint64               // offset: 10, size: 8 bytes
	MoveLot                          uint16               // offset: 18, size: 2 bytes
	ResultLotFrom                    uint16               // offset: 20, size: 2 bytes
	ResultLotTo                      uint16               // offset: 22, size: 2 bytes
}

// PROTO_NC_ITEMDB_SPLIT_N_MERGE_REQ (size = 38 bytes)
type PROTO_NC_ITEMDB_SPLIT_N_MERGE_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Itemid                           uint16               // offset: 2, size: 2 bytes
	MoveLot                          uint16               // offset: 4, size: 2 bytes
	FromOwner                        uint32               // offset: 6, size: 4 bytes
	FromItemKey                      uint64               // offset: 10, size: 8 bytes
	ToOwner                          uint32               // offset: 18, size: 4 bytes
	ToItemKey                        uint64               // offset: 22, size: 8 bytes
	NFromItemLot                     uint32               // offset: 30, size: 4 bytes
	NToItemLot                       uint32               // offset: 34, size: 4 bytes
}

// PROTO_NC_ITEMDB_SPLIT_REQ (size = 62 bytes)
type PROTO_NC_ITEMDB_SPLIT_REQ_DATA struct {
	OwnerA                           uint32               // offset: 0, size: 4 bytes
	OwnerB                           uint32               // offset: 4, size: 4 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	Itemid                           uint16               // offset: 10, size: 2 bytes
	Lotmain                          uint64               // offset: 12, size: 8 bytes
	Lotcreate                        [38]byte             // offset: 20, size: 38 bytes
	Splitlot                         uint16               // offset: 58, size: 2 bytes
	MainLot                          uint16               // offset: 60, size: 2 bytes
}

// PROTO_NC_ITEMDB_TERMEXTEND_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_TERMEXTEND_ACK_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NError                           uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_TERMEXTEND_REQ (size = 141 bytes)
type PROTO_NC_ITEMDB_TERMEXTEND_REQ_DATA struct {
	NLockIndex                       uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NHandleNum                       uint16               // offset: 6, size: 2 bytes
	NExtendItemKey                   uint64               // offset: 8, size: 8 bytes
	NExtendItemID                    uint16               // offset: 16, size: 2 bytes
	NExtendItemLot                   uint32               // offset: 18, size: 4 bytes
	NTermItemKey                     uint64               // offset: 22, size: 8 bytes
	NTermItemID                      uint16               // offset: 30, size: 2 bytes
	TermItemAttr                     [101]byte            // offset: 32, size: 101 bytes
	OrigDelDateTime                  uint32               // offset: 133, size: 4 bytes
	NExtendTime                      uint32               // offset: 137, size: 4 bytes
}

// PROTO_NC_ITEMDB_UNEQUIPFAIL_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_UNEQUIPFAIL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_UNEQUIPSUC_ACK (size = 2 bytes)
type PROTO_NC_ITEMDB_UNEQUIPSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_ITEMDB_UNEQUIP_REQ (size = 14 bytes)
type PROTO_NC_ITEMDB_UNEQUIP_REQ_DATA struct {
	Unequ                            [10]byte             // offset: 0, size: 10 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	Equloc                           uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_ITEMDB_UPGRADE_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_UPGRADE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_ITEMDB_UPGRADE_REQ (size = 78 bytes)
type PROTO_NC_ITEMDB_UPGRADE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Upgrade_item_regnum              uint64               // offset: 8, size: 8 bytes
	Upgrade_item_id                  uint16               // offset: 16, size: 2 bytes
	Upgrade_item_adj_level           uint8                // offset: 18, size: 1 bytes
	Upgrade_item_destroy             uint8                // offset: 19, size: 1 bytes
	Upgrade_item_new_level           uint8                // offset: 20, size: 1 bytes
	Upgrade_item_fail_count          uint8                // offset: 21, size: 1 bytes
	Raw_regnum                       uint64               // offset: 22, size: 8 bytes
	Raw_left_regnum                  uint64               // offset: 30, size: 8 bytes
	Raw_right_regnum                 uint64               // offset: 38, size: 8 bytes
	Raw_middle_regnum                uint64               // offset: 46, size: 8 bytes
	Gift_money                       uint64               // offset: 54, size: 8 bytes
	Raw_item_id                      uint16               // offset: 62, size: 2 bytes
	Raw_left_item_id                 uint16               // offset: 64, size: 2 bytes
	Raw_right_item_id                uint16               // offset: 66, size: 2 bytes
	Raw_lot                          uint8                // offset: 68, size: 1 bytes
	Raw_left_lot                     uint8                // offset: 69, size: 1 bytes
	Raw_right_lot                    uint8                // offset: 70, size: 1 bytes
	Raw_middle_item_id               uint16               // offset: 71, size: 2 bytes
	Raw_middle_lot                   uint8                // offset: 73, size: 1 bytes
	Is_randomoption_enable           uint8                // offset: 74, size: 1 bytes
	Upgrade_randomoption             [3]byte              // offset: 75, size: 3 bytes
}

// PROTO_NC_ITEMDB_USEALL_REQ (size = 12 bytes)
type PROTO_NC_ITEMDB_USEALL_REQ_DATA struct {
	Attritemregnum                   uint64               // offset: 0, size: 8 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	Itemid                           uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_ITEMDB_USELOT_REQ (size = 14 bytes)
type PROTO_NC_ITEMDB_USELOT_REQ_DATA struct {
	Attr                             [12]byte             // offset: 0, size: 12 bytes
	Lockindex                        uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_ITEMDB_USE_ACK (size = 4 bytes)
type PROTO_NC_ITEMDB_USE_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_ITEMDB_USE_VARIATION_ITEM_REQ (size = 21 bytes)
type PROTO_NC_ITEMDB_USE_VARIATION_ITEM_REQ struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Attr                             [12]byte             // offset: 6, size: 12 bytes
	Lockindex                        uint16               // offset: 18, size: 2 bytes
	IsClear                          uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_ITEMDB_WEAPONENDURESET_CMD (size = 12 bytes)
type PROTO_NC_ITEMDB_WEAPONENDURESET_CMD_DATA struct {
	WeaponKey                        uint64               // offset: 0, size: 8 bytes
	ChargedNumber                    uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_ACK (size = 24 bytes)
type PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	WeaponKey                        uint64               // offset: 8, size: 8 bytes
	WeaponSlot                       uint16               // offset: 16, size: 2 bytes
	ChargedNumber                    uint32               // offset: 18, size: 4 bytes
	Error                            uint16               // offset: 22, size: 2 bytes
}

// PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_REQ (size = 35 bytes)
type PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_REQ_DATA struct {
	WeaponKey                        uint64               // offset: 0, size: 8 bytes
	WeaponSlot                       uint16               // offset: 8, size: 2 bytes
	ChargedNumber                    uint32               // offset: 10, size: 4 bytes
	CommonData                       [21]byte             // offset: 14, size: 21 bytes
}

// PROTO_NC_ITEMDB_WITHDRAW_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_WITHDRAW_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Err                              uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_WITHDRAW_REQ (size = 20 bytes)
type PROTO_NC_ITEMDB_WITHDRAW_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Usrregnum                        uint32               // offset: 8, size: 4 bytes
	Cen                              uint64               // offset: 12, size: 8 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD_DATA struct {
	BItemInfoClose                   uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ_DATA struct {
	BItemInfoClose                   uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MISC_ITEMSHOP_URL_ACK (size = 4 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NLen                             uint16               // offset: 2, size: 2 bytes
	SURL                             [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK (size = 8 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NError                           uint16               // offset: 4, size: 2 bytes
	NLen                             uint16               // offset: 6, size: 2 bytes
	SURL                             [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ (size = 8 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NUserNo                          uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_MISC_ITEMSHOP_URL_REQ (size = 1 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_REQ_DATA struct {}

// PROTO_NC_PET_USE_ITEM_FAIL_ACK (size = 2 bytes)
type PROTO_NC_PET_USE_ITEM_FAIL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PROMOTION_REWARD_ITEM_CMD (size = 1 bytes)
type PROTO_NC_PROMOTION_REWARD_ITEM_CMD_DATA struct {
	NNumOfItems                      uint8                // offset: 0, size: 1 bytes
	RewardList                       [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD (size = 2 bytes)
type PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD (size = 6 bytes)
type PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
	NSelectedItemIndex               uint32               // offset: 2, size: 4 bytes
}
