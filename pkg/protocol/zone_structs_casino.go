package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_casino.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_DICE_TAISAI_BETTING_ACK (sizeof = 4 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NBettingIndex                    uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_ACK_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_DICE_TAISAI_BETTING_CANCEL_ACK (sizeof = 4 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_CANCEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NBettingCancelIndex              uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_CANCEL_ACK_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_DICE_TAISAI_BETTING_CANCEL_CMD (sizeof = 10 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_CANCEL_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
	BettingCancelInfo                uint64               // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_CANCEL_CMD_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ (sizeof = 10 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ_DATA struct {
	BettingCancelInfo                uint64               // offset: 0, size: 8 bytes
	NBettingCancelIndex              uint16               // offset: 8, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_DICE_TAISAI_BETTING_CMD (sizeof = 10 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
	BettingInfo                      uint64               // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_CMD_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_DICE_TAISAI_BETTING_REQ (sizeof = 10 bytes)
type PROTO_NC_DICE_TAISAI_BETTING_REQ_DATA struct {
	BettingInfo                      uint64               // offset: 0, size: 8 bytes
	NBettingIndex                    uint16               // offset: 8, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BETTING_REQ_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_DICE_TAISAI_BET_START_CMD (sizeof = 7 bytes)
type PROTO_NC_DICE_TAISAI_BET_START_CMD_DATA struct {
	BModeType                        uint32               // offset: 0, size: 4 bytes
	CurrentRollingDice               [3]byte              // offset: 4, size: 3 bytes
}

func (s *PROTO_NC_DICE_TAISAI_BET_START_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_DICE_TAISAI_BOARD_FOLD_CMD (sizeof = 1 bytes)
type PROTO_NC_DICE_TAISAI_BOARD_FOLD_CMD_DATA struct {}

func (s *PROTO_NC_DICE_TAISAI_BOARD_FOLD_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK (sizeof = 22 bytes)
type PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NGameNo                          int32                // offset: 2, size: 4 bytes
	Header                           NETPACKETZONEHEADER  // offset: 6, size: 6 bytes
	NActiveDiceTaiSaiKey             uint64               // offset: 12, size: 8 bytes
	NResultCnt                       uint16               // offset: 20, size: 2 bytes
	DiceRollResult                   [0]byte              // offset: 22, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK_DATA) FixedWireSize() int { return 22 }

// PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD (sizeof = 20 bytes)
type PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NActiveDiceTaiSaiKey             uint64               // offset: 6, size: 8 bytes
	NGameNo                          int32                // offset: 14, size: 4 bytes
	NResultLogCnt                    uint16               // offset: 18, size: 2 bytes
	DiceRollResultLog                [0]byte              // offset: 20, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD_DATA) FixedWireSize() int { return 20 }

// PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ (sizeof = 19 bytes)
type PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NActiveDiceTaiSaiKey             uint64               // offset: 6, size: 8 bytes
	DiceNum                          [3]byte              // offset: 14, size: 3 bytes
	NResultCnt                       uint16               // offset: 17, size: 2 bytes
	DiceRollResult                   [0]byte              // offset: 19, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ_DATA) FixedWireSize() int { return 19 }

// PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK (sizeof = 30 bytes)
type PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 2, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 8, size: 2 bytes
	NTaiSaiItemKey                   uint64               // offset: 10, size: 8 bytes
	NGameCnt                         uint32               // offset: 18, size: 4 bytes
	NExpectIncomeMoney               int64                // offset: 22, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK_DATA) FixedWireSize() int { return 30 }

// PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_REQ (sizeof = 26 bytes)
type PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 6, size: 2 bytes
	NTaiSaiItemKey                   uint64               // offset: 8, size: 8 bytes
	NRewardMoney                     uint64               // offset: 16, size: 8 bytes
	GetMasterRate                    uint16               // offset: 24, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_REQ_DATA) FixedWireSize() int { return 26 }

// PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK (sizeof = 32 bytes)
type PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 2, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 8, size: 2 bytes
	Lockindex                        uint16               // offset: 10, size: 2 bytes
	NDiceTaiSaiKey                   uint64               // offset: 12, size: 8 bytes
	NGameItemID                      uint16               // offset: 20, size: 2 bytes
	NMoneyCardItemID                 uint16               // offset: 22, size: 2 bytes
	NRewardMoney                     uint64               // offset: 24, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_REQ (sizeof = 32 bytes)
type PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 6, size: 2 bytes
	Lockindex                        uint16               // offset: 8, size: 2 bytes
	NDiceTaiSaiKey                   uint64               // offset: 10, size: 8 bytes
	NGameItemID                      uint16               // offset: 18, size: 2 bytes
	NMoneyCardItemID                 uint16               // offset: 20, size: 2 bytes
	NRewardMoney                     uint64               // offset: 22, size: 8 bytes
	GetMasterRate                    uint16               // offset: 30, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_REQ_DATA) FixedWireSize() int { return 32 }

// PROTO_NC_DICE_TAISAI_DB_RANK_ACK (sizeof = 13 bytes)
type PROTO_NC_DICE_TAISAI_DB_RANK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 2, size: 6 bytes
	NReqCharNo                       uint32               // offset: 8, size: 4 bytes
	NRankCnt                         uint8                // offset: 12, size: 1 bytes
	Rank                             [0]byte              // offset: 13, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_DB_RANK_ACK_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_DICE_TAISAI_DB_RANK_REQ (sizeof = 18 bytes)
type PROTO_NC_DICE_TAISAI_DB_RANK_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NReqCharNo                       uint32               // offset: 6, size: 4 bytes
	NTaiSaiItemKey                   uint64               // offset: 10, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_RANK_REQ_DATA) FixedWireSize() int { return 18 }

// PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK (sizeof = 28 bytes)
type PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Lockindex                        uint16               // offset: 2, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 4, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 10, size: 2 bytes
	NTaiSaiItemKey                   uint64               // offset: 12, size: 8 bytes
	NReceiptIncomeMoney              uint64               // offset: 20, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK_DATA) FixedWireSize() int { return 28 }

// PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_REQ (sizeof = 18 bytes)
type PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 2, size: 6 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 8, size: 2 bytes
	NTaiSaiItemKey                   uint64               // offset: 10, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_REQ_DATA) FixedWireSize() int { return 18 }

// PROTO_NC_DICE_TAISAI_DICE_ROLL_ACK (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_DICE_ROLL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DICE_ROLL_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_DICE_ROLL_CMD (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_DICE_ROLL_CMD_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DICE_ROLL_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD (sizeof = 5 bytes)
type PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NDiceInfo                        [3]byte              // offset: 2, size: 3 bytes
}

func (s *PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD (sizeof = 1 bytes)
type PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD_DATA struct {
	NGamerCnt                        uint8                // offset: 0, size: 1 bytes
	Emotion                          [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_ACK (sizeof = 16 bytes)
type PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 2, size: 2 bytes
	NGameCnt                         uint32               // offset: 4, size: 4 bytes
	NExpectIncomeMoney               uint64               // offset: 8, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_ACK_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ_DATA struct {
	NInvenSlot                       ITEM_INVEN           // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_GAME_BETTING_INFO_CMD (sizeof = 3 bytes)
type PROTO_NC_DICE_TAISAI_GAME_BETTING_INFO_CMD_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	NBettingCnt                      uint16               // offset: 1, size: 2 bytes
	Betting                          [0]byte              // offset: 3, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_GAME_BETTING_INFO_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK (sizeof = 50 bytes)
type PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	DiceHistory                      [21]byte             // offset: 2, size: 21 bytes
	NGameStatus                      uint32               // offset: 23, size: 4 bytes
	BModeType                        uint32               // offset: 27, size: 4 bytes
	NTimer                           uint16               // offset: 31, size: 2 bytes
	NInterval                        uint32               // offset: 33, size: 4 bytes
	NGetSysRate                      uint16               // offset: 37, size: 2 bytes
	NMaxBetMoney                     uint64               // offset: 39, size: 8 bytes
	NUseMinLv                        uint16               // offset: 47, size: 2 bytes
	NGamerCnt                        uint8                // offset: 49, size: 1 bytes
	NGamerHnd                        [0]byte              // offset: 50, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK_DATA) FixedWireSize() int { return 50 }

// PROTO_NC_DICE_TAISAI_GAME_JOIN_CMD (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_GAME_JOIN_CMD_DATA struct {
	NPlayerHnd                       uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_GAME_JOIN_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_GAME_LEAVE_ACK (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_GAME_LEAVE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_GAME_LEAVE_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_GAME_LEAVE_CMD (sizeof = 3 bytes)
type PROTO_NC_DICE_TAISAI_GAME_LEAVE_CMD_DATA struct {
	BReserve                         uint8                // offset: 0, size: 1 bytes
	NPlayerHnd                       uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_GAME_LEAVE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_ACK (sizeof = 6 bytes)
type PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	BModeType                        uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_ACK_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_CMD (sizeof = 4 bytes)
type PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_CMD_DATA struct {
	BModeType                        uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_CMD (sizeof = 1 bytes)
type PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_CMD_DATA struct {
	NGamerCnt                        uint8                // offset: 0, size: 1 bytes
	LargeAmount                      [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_DICE_TAISAI_RANK_ACK (sizeof = 484 bytes)
type PROTO_NC_DICE_TAISAI_RANK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NRankHistoryCnt                  uint8                // offset: 2, size: 1 bytes
	RankHistory                      [480]byte            // offset: 3, size: 480 bytes
	NRankCurrCnt                     uint8                // offset: 483, size: 1 bytes
	RankCurr                         [0]byte              // offset: 484, flexible array member
}

func (s *PROTO_NC_DICE_TAISAI_RANK_ACK_DATA) FixedWireSize() int { return 484 }

// PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_ACK (sizeof = 12 bytes)
type PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NInvenSlot                       ITEM_INVEN           // offset: 2, size: 2 bytes
	ReceiptIncomeMoney               uint64               // offset: 4, size: 8 bytes
}

func (s *PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_ACK_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ (sizeof = 2 bytes)
type PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ_DATA struct {
	NInvenSlot                       ITEM_INVEN           // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_DICE_TAISAI_TIMER_CMD (sizeof = 7 bytes)
type PROTO_NC_DICE_TAISAI_TIMER_CMD_DATA struct {
	BIsStart                         uint8                // offset: 0, size: 1 bytes
	NTimer                           uint16               // offset: 1, size: 2 bytes
	NInterval                        uint32               // offset: 3, size: 4 bytes
}

func (s *PROTO_NC_DICE_TAISAI_TIMER_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_GAMBLE_COIN_BUY_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_BUY_REQ_DATA)

// PROTO_NC_GAMBLE_COIN_CHANGE_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_CHANGE_CMD_DATA)

// PROTO_NC_GAMBLE_COIN_DB_BUY_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_BUY_ACK_DATA)

// PROTO_NC_GAMBLE_COIN_DB_BUY_REQ (sizeof = 32 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_BUY_REQ_DATA)

// PROTO_NC_GAMBLE_COIN_DB_SELL_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_SELL_ACK_DATA)

// PROTO_NC_GAMBLE_COIN_DB_SELL_REQ (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_SELL_REQ_DATA)

// PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK (sizeof = 18 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK_DATA)

// PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_REQ (sizeof = 24 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_REQ_DATA)

// PROTO_NC_GAMBLE_COIN_USE_COINITEM_MESSAGE_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_USE_COINITEM_MESSAGE_CMD_DATA)

// PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_BET_START_CMD (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_BET_START_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK (sizeof = 9 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_REQ (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK (sizeof = 14 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ (sizeof = 11 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD (sizeof = 5 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_BETTING_INFO_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_BETTING_INFO_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK (sizeof = 50 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD_DATA)

// PROTO_NC_GAMBLE_DICE_TAISAI_TIMER_CMD (sizeof = 7 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_DICE_TAISAI_TIMER_CMD_DATA)

// PROTO_NC_GAMBLE_ENTER_PLAYER_DIRECT_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_ENTER_PLAYER_DIRECT_CMD_DATA)

// PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD_DATA)

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK_DATA)

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_REQ (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_REQ_DATA)

// PROTO_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD_DATA)

// PROTO_NC_GAMBLE_PLAYERACT_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_PLAYERACT_CMD_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK (sizeof = 8 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_REQ (sizeof = 36 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_REQ_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD (sizeof = 28 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK (sizeof = 11 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK_DATA)

// PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK (sizeof = 1 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK_DATA)

// PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK (sizeof = 9 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK_DATA)

// PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ_DATA)

// PROTO_NC_GAMBLE_WORLD_PREVMAPNAME_CMD (sizeof = 18 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_WORLD_PREVMAPNAME_CMD_DATA)

// PROTO_NC_GAMBLE_ZONE_PREVMAPNAME_CMD (sizeof = 12 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_GAMBLE_ZONE_PREVMAPNAME_CMD_DATA)
