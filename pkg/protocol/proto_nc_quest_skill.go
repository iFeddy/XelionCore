package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_REQ (size = 21 bytes)
type PROTO_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Restatitemkey                    uint64               // offset: 8, size: 8 bytes
	Restatitemid                     uint16               // offset: 16, size: 2 bytes
	Isdelete                         uint8                // offset: 18, size: 1 bytes
	Attr                             uint16               // offset: 19, size: 2 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD (size = 8 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Questdoing                       [6]byte              // offset: 2, size: 6 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD (size = 14 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Questdone                        [12]byte             // offset: 2, size: 12 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD (size = 8 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Questread                        [6]byte              // offset: 2, size: 6 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD (size = 8 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	Questrepeat                      [6]byte              // offset: 2, size: 6 bytes
}

// PROTO_NC_CHARSAVE_2WLDMAN_SKILL_CMD (size = 13 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_SKILL_CMD_DATA struct {
	Wldmanhandle                     uint16               // offset: 0, size: 2 bytes
	PartMark                         uint8                // offset: 2, size: 1 bytes
	NMaxNum                          uint16               // offset: 3, size: 2 bytes
	Skill                            uint64               // offset: 5, size: 8 bytes
}

// PROTO_NC_CHARSAVE_QUEST_DOING_REQ (size = 8 bytes)
type PROTO_NC_CHARSAVE_QUEST_DOING_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	BNeedClear                       uint8                // offset: 6, size: 1 bytes
	NNumOfDoingQuest                 uint8                // offset: 7, size: 1 bytes
	QuestDoingArray                  [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_CHARSAVE_SKILL_REQ (size = 8 bytes)
type PROTO_NC_CHARSAVE_SKILL_REQ_DATA struct {
	Handle                           [6]byte              // offset: 0, size: 6 bytes
	Number                           uint16               // offset: 6, size: 2 bytes
	Skill                            [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_CHAR_CLIENT_SKILL_CMD (size = 10 bytes)
type PROTO_NC_CHAR_CLIENT_SKILL_CMD_DATA struct {
	Restempow                        uint8                // offset: 0, size: 1 bytes
	PartMark                         uint8                // offset: 1, size: 1 bytes
	NMaxNum                          uint16               // offset: 2, size: 2 bytes
	Skill                            [6]byte              // offset: 4, size: 6 bytes
}

// PROTO_NC_CHAR_QUEST_DOING_CMD (size = 6 bytes)
type PROTO_NC_CHAR_QUEST_DOING_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	BNeedClear                       uint8                // offset: 4, size: 1 bytes
	NNumOfDoingQuest                 uint8                // offset: 5, size: 1 bytes
	QuestDoingArray                  [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_QUEST_DONE_CMD (size = 12 bytes)
type PROTO_NC_CHAR_QUEST_DONE_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NTotalDoneQuest                  uint16               // offset: 4, size: 2 bytes
	NTotalDoneQuestSize              uint16               // offset: 6, size: 2 bytes
	NDoneQuestCount                  uint16               // offset: 8, size: 2 bytes
	NIndex                           uint16               // offset: 10, size: 2 bytes
	QuestDoneArray                   [0]byte              // offset: 12, flexible array member
}

// PROTO_NC_CHAR_QUEST_READ_CMD (size = 6 bytes)
type PROTO_NC_CHAR_QUEST_READ_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NNumOfReadQuest                  uint16               // offset: 4, size: 2 bytes
	QuestReadIDArray                 [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_QUEST_REPEAT_CMD (size = 6 bytes)
type PROTO_NC_CHAR_QUEST_REPEAT_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NNumOfRepeatQuest                uint16               // offset: 4, size: 2 bytes
	QuestRepeatArray                 [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_SKILLCLIENT_CMD (size = 6 bytes)
type PROTO_NC_CHAR_SKILLCLIENT_CMD struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Number                           uint16               // offset: 4, size: 2 bytes
	Skill                            [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_CHAR_SKILL_CMD (size = 9 bytes)
type PROTO_NC_CHAR_SKILL_CMD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	PartMark                         uint8                // offset: 4, size: 1 bytes
	NMaxNum                          uint16               // offset: 5, size: 2 bytes
	Number                           uint16               // offset: 7, size: 2 bytes
	Skill                            [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_CHAR_SKILL_PASSIVE_CMD (size = 2 bytes)
type PROTO_NC_CHAR_SKILL_PASSIVE_CMD_DATA struct {
	Number                           uint16               // offset: 0, size: 2 bytes
	Passive                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK (size = 2 bytes)
type PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK_DATA struct {
	NScenarioID                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ (size = 2 bytes)
type PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ_DATA struct {
	NScenarioID                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_DB_CLEAR_ACK (size = 10 bytes)
type PROTO_NC_QUEST_DB_CLEAR_ACK_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQuestID                         uint16               // offset: 6, size: 2 bytes
	ErrorType                        uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_QUEST_DB_CLEAR_REQ (size = 8 bytes)
type PROTO_NC_QUEST_DB_CLEAR_REQ_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQuestID                         uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_QUEST_DB_DONE_ACK (size = 43 bytes)
type PROTO_NC_QUEST_DB_DONE_ACK_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQSC                             uint8                // offset: 6, size: 1 bytes
	QuestInfo                        [32]byte             // offset: 7, size: 32 bytes
	NLockIndex                       uint16               // offset: 39, size: 2 bytes
	ErrorType                        uint16               // offset: 41, size: 2 bytes
}

// PROTO_NC_QUEST_DB_DONE_REQ (size = 83 bytes)
type PROTO_NC_QUEST_DB_DONE_REQ_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	UserNo                           uint32               // offset: 6, size: 4 bytes
	NPCID                            uint16               // offset: 10, size: 2 bytes
	NQSC                             uint8                // offset: 12, size: 1 bytes
	QuestInfo                        [32]byte             // offset: 13, size: 32 bytes
	NMoney                           uint32               // offset: 45, size: 4 bytes
	NExp                             uint32               // offset: 49, size: 4 bytes
	NFame                            uint32               // offset: 53, size: 4 bytes
	NCur_Money                       uint64               // offset: 57, size: 8 bytes
	NCur_Exp                         uint64               // offset: 65, size: 8 bytes
	NCur_Fame                        uint32               // offset: 73, size: 4 bytes
	CT_Type                          uint8                // offset: 77, size: 1 bytes
	CT_ElementNo                     uint8                // offset: 78, size: 1 bytes
	NCur_CT_ElementValue             uint8                // offset: 79, size: 1 bytes
	NLockIndex                       uint16               // offset: 80, size: 2 bytes
	NItemCreateCount                 uint8                // offset: 82, size: 1 bytes
	ItemCreate                       [0]byte              // offset: 83, flexible array member
}

// PROTO_NC_QUEST_DB_GIVE_UP_ACK (size = 10 bytes)
type PROTO_NC_QUEST_DB_GIVE_UP_ACK_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQuestID                         uint16               // offset: 6, size: 2 bytes
	ErrorType                        uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_QUEST_DB_GIVE_UP_REQ (size = 8 bytes)
type PROTO_NC_QUEST_DB_GIVE_UP_REQ_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQuestID                         uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_QUEST_DB_READ_ACK (size = 10 bytes)
type PROTO_NC_QUEST_DB_READ_ACK_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	NQuestID                         uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_QUEST_DB_READ_REQ (size = 8 bytes)
type PROTO_NC_QUEST_DB_READ_REQ_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQuestID                         uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_QUEST_DB_SET_INFO_ACK (size = 41 bytes)
type PROTO_NC_QUEST_DB_SET_INFO_ACK_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQSC                             uint8                // offset: 6, size: 1 bytes
	QuestInfo                        [32]byte             // offset: 7, size: 32 bytes
	ErrorType                        uint16               // offset: 39, size: 2 bytes
}

// PROTO_NC_QUEST_DB_SET_INFO_REQ (size = 39 bytes)
type PROTO_NC_QUEST_DB_SET_INFO_REQ_DATA struct {
	ZoneHeader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NQSC                             uint8                // offset: 6, size: 1 bytes
	QuestInfo                        [32]byte             // offset: 7, size: 32 bytes
}

// PROTO_NC_QUEST_GIVE_UP_ACK (size = 4 bytes)
type PROTO_NC_QUEST_GIVE_UP_ACK_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
	ErrorCode                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_QUEST_GIVE_UP_REQ (size = 2 bytes)
type PROTO_NC_QUEST_GIVE_UP_REQ_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_JOBDUNGEON_FIND_RNG (size = 115 bytes)
type PROTO_NC_QUEST_JOBDUNGEON_FIND_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	NError                           uint16               // offset: 7, size: 2 bytes
	ClientMapName                    [33]byte             // offset: 9, size: 33 bytes
	ServerMapName                    [33]byte             // offset: 42, size: 33 bytes
	ScriptName                       [33]byte             // offset: 75, size: 33 bytes
	QuestScript                      [7]byte              // offset: 108, size: 7 bytes
}

// PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD (size = 1 bytes)
type PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD_DATA struct {
	NumOfActionQuest                 uint8                // offset: 0, size: 1 bytes
	Mobofquest                       [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_QUEST_READ_ACK (size = 4 bytes)
type PROTO_NC_QUEST_READ_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NQuestID                         uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_QUEST_READ_REQ (size = 2 bytes)
type PROTO_NC_QUEST_READ_REQ_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_RESET_TIME_CMD (size = 16 bytes)
type PROTO_NC_QUEST_RESET_TIME_CMD_DATA struct {
	TResetYearQuest                  int32                // offset: 0, size: 4 bytes
	TResetMonthQuest                 int32                // offset: 4, size: 4 bytes
	TResetWeekQuest                  int32                // offset: 8, size: 4 bytes
	TResetDayQuest                   int32                // offset: 12, size: 4 bytes
}

// PROTO_NC_QUEST_RESET_TIME_ZONE_CMD (size = 20 bytes)
type PROTO_NC_QUEST_RESET_TIME_ZONE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	TResetYearQuest                  int32                // offset: 4, size: 4 bytes
	TResetMonthQuest                 int32                // offset: 8, size: 4 bytes
	TResetWeekQuest                  int32                // offset: 12, size: 4 bytes
	TResetDayQuest                   int32                // offset: 16, size: 4 bytes
}

// PROTO_NC_QUEST_SCENARIO_RUN_CMD (size = 2 bytes)
type PROTO_NC_QUEST_SCENARIO_RUN_CMD_DATA struct {
	NScenarioID                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_SCRIPT_CMD_ACK (size = 7 bytes)
type PROTO_NC_QUEST_SCRIPT_CMD_ACK_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
	NQSC                             uint8                // offset: 2, size: 1 bytes
	NResult                          uint32               // offset: 3, size: 4 bytes
}

// PROTO_NC_QUEST_SCRIPT_CMD_REQ (size = 103 bytes)
type PROTO_NC_QUEST_SCRIPT_CMD_REQ_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
	Command                          [101]byte            // offset: 2, size: 101 bytes
}

// PROTO_NC_QUEST_SELECT_START_ACK (size = 6 bytes)
type PROTO_NC_QUEST_SELECT_START_ACK_DATA struct {
	NNPCID                           uint16               // offset: 0, size: 2 bytes
	NQuestID                         uint16               // offset: 2, size: 2 bytes
	ErrorType                        uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_QUEST_SELECT_START_REQ (size = 4 bytes)
type PROTO_NC_QUEST_SELECT_START_REQ_DATA struct {
	NNPCID                           uint16               // offset: 0, size: 2 bytes
	NQuestID                         uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_QUEST_START_ACK (size = 2 bytes)
type PROTO_NC_QUEST_START_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_QUEST_START_REQ (size = 2 bytes)
type PROTO_NC_QUEST_START_REQ_DATA struct {
	NQuestID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_SKILL_COOLTIME_CMD (size = 6 bytes)
type PROTO_NC_SKILL_COOLTIME_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Sec                              uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_SKILL_EMPOWALLOC_ACK (size = 4 bytes)
type PROTO_NC_SKILL_EMPOWALLOC_ACK_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_SKILL_EMPOWALLOC_DB_ACK (size = 13 bytes)
type PROTO_NC_SKILL_EMPOWALLOC_DB_ACK_DATA struct {
	Zonehandle                       uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Skill                            uint16               // offset: 8, size: 2 bytes
	Err                              uint16               // offset: 10, size: 2 bytes
	BUseItem                         uint8                // offset: 12, size: 1 bytes
}

// PROTO_NC_SKILL_EMPOWALLOC_DB_REQ (size = 40 bytes)
type PROTO_NC_SKILL_EMPOWALLOC_DB_REQ_DATA struct {
	Zonehandle                       uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Skill                            uint16               // offset: 8, size: 2 bytes
	BUseItem                         uint8                // offset: 10, size: 1 bytes
	Empower                          [6]byte              // offset: 11, size: 6 bytes
	Redistitem                       [23]byte             // offset: 17, size: 23 bytes
}

// PROTO_NC_SKILL_EMPOWALLOC_REQ (size = 6 bytes)
type PROTO_NC_SKILL_EMPOWALLOC_REQ_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Plus                             uint16               // offset: 2, size: 2 bytes
	Minus                            uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_SKILL_EMPOWPOINT_CMD (size = 1 bytes)
type PROTO_NC_SKILL_EMPOWPOINT_CMD_DATA struct {
	Restempow                        uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_SKILL_EMPOW_RESET_DB_REQ (size = 6 bytes)
type PROTO_NC_SKILL_EMPOW_RESET_DB_REQ_DATA struct {
	Zonehandle                       uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_SKILL_EMPOW_RESET_SUC_ACK (size = 1 bytes)
type PROTO_NC_SKILL_EMPOW_RESET_SUC_ACK_DATA struct {
	Totalempow                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_SKILL_ERASE_ACK (size = 10 bytes)
type PROTO_NC_SKILL_ERASE_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Skill                            uint16               // offset: 6, size: 2 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_SKILL_ERASE_REQ (size = 8 bytes)
type PROTO_NC_SKILL_ERASE_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Skill                            uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_SKILL_JUMP_CMD (size = 21 bytes)
type PROTO_NC_SKILL_JUMP_CMD_DATA struct {
	NObject                          uint16               // offset: 0, size: 2 bytes
	From                             uint64               // offset: 2, size: 8 bytes
	To                               uint64               // offset: 10, size: 8 bytes
	NDirect                          uint8                // offset: 18, size: 1 bytes
	SkillID                          uint16               // offset: 19, size: 2 bytes
}

// PROTO_NC_SKILL_PASSIVESKILL_LEARN_CMD (size = 8 bytes)
type PROTO_NC_SKILL_PASSIVESKILL_LEARN_CMD_DATA struct {
	Wmshandle                        uint16               // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Passiveskillid                   uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_SKILL_PRODUCTFIELD_ACK (size = 4 bytes)
type PROTO_NC_SKILL_PRODUCTFIELD_ACK_DATA struct {
	Mainskill                        uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_SKILL_PRODUCTFIELD_REQ (size = 2 bytes)
type PROTO_NC_SKILL_PRODUCTFIELD_REQ_DATA struct {
	Mainskill                        uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_SKILL_REVIVE_CMD (size = 12 bytes)
type PROTO_NC_SKILL_REVIVE_CMD_DATA struct {
	Hp                               uint32               // offset: 0, size: 4 bytes
	Sp                               uint32               // offset: 4, size: 4 bytes
	Lp                               uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_SKILL_SKILLEXP_CLIENT_CMD (size = 6 bytes)
type PROTO_NC_SKILL_SKILLEXP_CLIENT_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Mastery                          uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_SKILL_SKILLTEACHSUC_ACK (size = 3 bytes)
type PROTO_NC_SKILL_SKILLTEACHSUC_ACK_DATA struct {
	Skillid                          uint16               // offset: 0, size: 2 bytes
	Ispassive                        uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_SKILL_SKILLTEACH_REQ (size = 9 bytes)
type PROTO_NC_SKILL_SKILLTEACH_REQ_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Skillid                          uint16               // offset: 4, size: 2 bytes
	Ispassive                        uint8                // offset: 6, size: 1 bytes
	Removeid                         uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_SKILL_SOMEONEREVIVE_CMD (size = 14 bytes)
type PROTO_NC_SKILL_SOMEONEREVIVE_CMD_DATA struct {
	Reviver                          uint16               // offset: 0, size: 2 bytes
	Hp                               uint32               // offset: 2, size: 4 bytes
	Sp                               uint32               // offset: 6, size: 4 bytes
	Lp                               uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_SKILL_UNLEARN_ACK (size = 4 bytes)
type PROTO_NC_SKILL_UNLEARN_ACK_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_SKILL_UNLEARN_REQ (size = 3 bytes)
type PROTO_NC_SKILL_UNLEARN_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
	Skill                            uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_SKILL_WARP_CMD (size = 20 bytes)
type PROTO_NC_SKILL_WARP_CMD_DATA struct {
	Object                           uint16               // offset: 0, size: 2 bytes
	From                             uint64               // offset: 2, size: 8 bytes
	To                               uint64               // offset: 10, size: 8 bytes
	Skill                            uint16               // offset: 18, size: 2 bytes
}
