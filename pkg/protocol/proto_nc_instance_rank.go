package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHAR_EMPTY_INSTANCE_DUNGEON_RNG (size = 11 bytes)
type PROTO_NC_CHAR_EMPTY_INSTANCE_DUNGEON_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Category                         uint32               // offset: 7, size: 4 bytes
}

// PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK (size = 3 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NRankAllCnt                      uint8                // offset: 2, size: 1 bytes
	RankAll                          [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ (size = 1 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ_DATA struct {}

// PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK (size = 3 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NRankCurrCnt                     uint8                // offset: 2, size: 1 bytes
	RankCurr                         [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ (size = 1 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ_DATA struct {}

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK (size = 9 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Header                           NETPACKETZONEHEADER  // offset: 2, size: 6 bytes
	NRankAllCnt                      uint8                // offset: 8, size: 1 bytes
	RankAll                          [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_REQ (size = 6 bytes)
type PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK (size = 1 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK_DATA struct {
	ListLen                          uint8                // offset: 0, size: 1 bytes
	RankList                         [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK (size = 1 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK_DATA struct {
	ListLen                          uint8                // offset: 0, size: 1 bytes
	RankList                         [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD (size = 9 bytes)
type PROTO_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD_DATA struct {
	Category                         uint8                // offset: 0, size: 1 bytes
	RegisterNumber                   uint32               // offset: 1, size: 4 bytes
	ReaderRegisterNumber             uint32               // offset: 5, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_EMPTY_DUNGEON_CMD (size = 5 bytes)
type PROTO_NC_INSTANCE_DUNGEON_EMPTY_DUNGEON_CMD_DATA struct {
	Category                         uint8                // offset: 0, size: 1 bytes
	RegisterNumber                   uint32               // offset: 1, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK (size = 12 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NNPCHandle                       uint16               // offset: 4, size: 2 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	ELevelType                       uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_ACK (size = 13 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_ACK_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NCharRegNum                      uint32               // offset: 1, size: 4 bytes
	NNPCHandle                       uint16               // offset: 5, size: 2 bytes
	NError                           uint16               // offset: 7, size: 2 bytes
	ELevelType                       uint32               // offset: 9, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ (size = 35 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NCharRegNum                      uint32               // offset: 1, size: 4 bytes
	NInDunRegNum                     uint32               // offset: 5, size: 4 bytes
	IndunCategory                    uint32               // offset: 9, size: 4 bytes
	NNPCHandle                       uint16               // offset: 13, size: 2 bytes
	Argument                         [20]byte             // offset: 15, size: 20 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_REQ (size = 36 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_REQ_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NDestZoneNum                     uint8                // offset: 1, size: 1 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NInDunRegNum                     uint32               // offset: 6, size: 4 bytes
	IndunCategory                    uint32               // offset: 10, size: 4 bytes
	NNPCHandle                       uint16               // offset: 14, size: 2 bytes
	Argument                         [20]byte             // offset: 16, size: 20 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK (size = 24 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NNPCHandle                       uint16               // offset: 4, size: 2 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
	ELevelType                       uint32               // offset: 8, size: 4 bytes
	SMapName                         Name3                // offset: 12, size: 12 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_ACK (size = 25 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_ACK_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NCharRegNum                      uint32               // offset: 1, size: 4 bytes
	NNPCHandle                       uint16               // offset: 5, size: 2 bytes
	NError                           uint16               // offset: 7, size: 2 bytes
	ELevelType                       uint32               // offset: 9, size: 4 bytes
	SMapName                         Name3                // offset: 13, size: 12 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ (size = 39 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NCharRegNum                      uint32               // offset: 1, size: 4 bytes
	NInDunRegNum                     uint32               // offset: 5, size: 4 bytes
	IndunCategory                    uint32               // offset: 9, size: 4 bytes
	NNPCHandle                       uint16               // offset: 13, size: 2 bytes
	Argument                         [20]byte             // offset: 15, size: 20 bytes
	ELevelType                       uint32               // offset: 35, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_REQ (size = 40 bytes)
type PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_REQ_DATA struct {
	NStartZoneNum                    uint8                // offset: 0, size: 1 bytes
	NDestZoneNum                     uint8                // offset: 1, size: 1 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
	NInDunRegNum                     uint32               // offset: 6, size: 4 bytes
	IndunCategory                    uint32               // offset: 10, size: 4 bytes
	NNPCHandle                       uint16               // offset: 14, size: 2 bytes
	Argument                         [20]byte             // offset: 16, size: 20 bytes
	ELevelType                       uint32               // offset: 36, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD (size = 2 bytes)
type PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD_DATA struct {
	NMapCount                        uint16               // offset: 0, size: 2 bytes
	SIndunInfos                      [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_ACK (size = 1972 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_ACK_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
	DamageDealtShapeInfo1st          uint32               // offset: 13, size: 4 bytes
	DamageDealtShapeInfo2nd          uint32               // offset: 17, size: 4 bytes
	DamageDealtShapeInfo3th          uint32               // offset: 21, size: 4 bytes
	DamageTakenShapeInfo1st          uint32               // offset: 25, size: 4 bytes
	DamageTakenShapeInfo2nd          uint32               // offset: 29, size: 4 bytes
	DamageTakenShapeInfo3th          uint32               // offset: 33, size: 4 bytes
	HealingDoneShapeInfo1st          uint32               // offset: 37, size: 4 bytes
	HealingDoneShapeInfo2nd          uint32               // offset: 41, size: 4 bytes
	HealingDoneShapeInfo3th          uint32               // offset: 45, size: 4 bytes
	NDamageDealtCount                uint8                // offset: 49, size: 1 bytes
	NDamageTakenCount                uint8                // offset: 50, size: 1 bytes
	NHealingDoneCount                uint8                // offset: 51, size: 1 bytes
	DamageDealtRanks                 [640]byte            // offset: 52, size: 640 bytes
	DamageTakenRanks                 [640]byte            // offset: 692, size: 640 bytes
	HealingDoneRanks                 [640]byte            // offset: 1332, size: 640 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_MYRANK_CMD (size = 6 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_MYRANK_CMD_DATA struct {
	NCharRegNum                      uint32               // offset: 0, size: 4 bytes
	NMyRankCount                     uint16               // offset: 4, size: 2 bytes
	MyRank                           [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_NEW_RANK_MSG (size = 38 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_NEW_RANK_MSG_DATA struct {
	NCharID                          Name5                // offset: 0, size: 20 bytes
	SIndunInfo                       [13]byte             // offset: 20, size: 13 bytes
	RankType                         uint8                // offset: 33, size: 1 bytes
	NRank                            uint32               // offset: 34, size: 4 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_REQ (size = 2 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_REQ_DATA struct {
	NMapCount                        uint16               // offset: 0, size: 2 bytes
	SIndunInfos                      [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_WORLD_CMD (size = 19 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_WORLD_CMD_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
	NPlayTime                        uint32               // offset: 13, size: 4 bytes
	NRankCount                       uint16               // offset: 17, size: 2 bytes
	Ranks                            [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_ACK (size = 19 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_ACK_DATA struct {
	ShapeInfo1st                     uint32               // offset: 0, size: 4 bytes
	ShapeInfo2nd                     uint32               // offset: 4, size: 4 bytes
	ShapeInfo3th                     uint32               // offset: 8, size: 4 bytes
	NErr                             uint16               // offset: 12, size: 2 bytes
	NDamageDealtRank                 uint32               // offset: 14, size: 4 bytes
	NDamageDealtCount                uint8                // offset: 18, size: 1 bytes
	DamageDealtRanks                 [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_REQ (size = 13 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_REQ_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_ACK (size = 19 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_ACK_DATA struct {
	ShapeInfo1st                     uint32               // offset: 0, size: 4 bytes
	ShapeInfo2nd                     uint32               // offset: 4, size: 4 bytes
	ShapeInfo3th                     uint32               // offset: 8, size: 4 bytes
	NErr                             uint16               // offset: 12, size: 2 bytes
	NDamageTakenRank                 uint32               // offset: 14, size: 4 bytes
	NDamageTakenCount                uint8                // offset: 18, size: 1 bytes
	DamageTakenRanks                 [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_REQ (size = 13 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_REQ_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_ACK (size = 19 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_ACK_DATA struct {
	ShapeInfo1st                     uint32               // offset: 0, size: 4 bytes
	ShapeInfo2nd                     uint32               // offset: 4, size: 4 bytes
	ShapeInfo3th                     uint32               // offset: 8, size: 4 bytes
	NErr                             uint16               // offset: 12, size: 2 bytes
	NHealingDoneRank                 uint32               // offset: 14, size: 4 bytes
	NHealingDoneCount                uint8                // offset: 18, size: 1 bytes
	HealingDoneRanks                 [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_REQ (size = 13 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_REQ_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_ACK (size = 62 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_ACK_DATA struct {
	NErr                             uint16               // offset: 0, size: 2 bytes
	NPrevDamageDealt                 uint64               // offset: 2, size: 8 bytes
	NPrevDamageTaken                 uint64               // offset: 10, size: 8 bytes
	NPrevHealingDone                 uint64               // offset: 18, size: 8 bytes
	NPrevPlayTimeHour                uint8                // offset: 26, size: 1 bytes
	NPrevPlayTimeMinute              uint8                // offset: 27, size: 1 bytes
	NPrevPlayTimeSecond              uint8                // offset: 28, size: 1 bytes
	NPrevPlayTimeYear                uint8                // offset: 29, size: 1 bytes
	NPrevPlayTimeMonth               uint8                // offset: 30, size: 1 bytes
	NPrevPlayTimeDay                 uint8                // offset: 31, size: 1 bytes
	NCurrDamageDealt                 uint64               // offset: 32, size: 8 bytes
	NCurrDamageTaken                 uint64               // offset: 40, size: 8 bytes
	NCurrHealingDone                 uint64               // offset: 48, size: 8 bytes
	NCurrPlayTimeHour                uint8                // offset: 56, size: 1 bytes
	NCurrPlayTimeMinute              uint8                // offset: 57, size: 1 bytes
	NCurrPlayTimeSecond              uint8                // offset: 58, size: 1 bytes
	NCurrPlayTimeYear                uint8                // offset: 59, size: 1 bytes
	NCurrPlayTimeMonth               uint8                // offset: 60, size: 1 bytes
	NCurrPlayTimeDay                 uint8                // offset: 61, size: 1 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_REQ (size = 13 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_REQ_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_UPDATE_DB_ACK (size = 15 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_UPDATE_DB_ACK_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
	CharListCount                    uint16               // offset: 13, size: 2 bytes
	CharList                         [0]byte              // offset: 15, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RANK_UPDATE_DB_REQ (size = 15 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RANK_UPDATE_DB_REQ_DATA struct {
	SIndunInfo                       [13]byte             // offset: 0, size: 13 bytes
	CharListCount                    uint16               // offset: 13, size: 2 bytes
	CharList                         [0]byte              // offset: 15, flexible array member
}

// PROTO_NC_INSTANCE_DUNGEON_RESET_COUNTDOWN_CMD (size = 5 bytes)
type PROTO_NC_INSTANCE_DUNGEON_RESET_COUNTDOWN_CMD_DATA struct {
	Category                         uint8                // offset: 0, size: 1 bytes
	RegisterNumber                   uint32               // offset: 1, size: 4 bytes
}
