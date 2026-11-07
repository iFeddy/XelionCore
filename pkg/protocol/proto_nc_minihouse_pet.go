package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_CHARSAVE_2WLDMAN_PET_LINK_RESUMMON_CMD (size = 23 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_PET_LINK_RESUMMON_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	BDoResummon                      uint8                // offset: 4, size: 1 bytes
	XResummonCode                    uint64               // offset: 5, size: 8 bytes
	NSummonSlot                      uint16               // offset: 13, size: 2 bytes
	NId                              uint32               // offset: 15, size: 4 bytes
	NRegNum                          uint32               // offset: 19, size: 4 bytes
}

// PROTO_NC_CHAR_PET_CMD (size = 1 bytes)
type PROTO_NC_CHAR_PET_CMD_DATA struct {}

// PROTO_NC_MINIHOUSE_ACTIV_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_ACTIV_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ_DATA struct {
	ArrangeMode                      uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_BUILDING_CMD (size = 66 bytes)
type PROTO_NC_MINIHOUSE_BUILDING_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	Minihouseinfo                    [64]byte             // offset: 2, size: 64 bytes
}

// PROTO_NC_MINIHOUSE_BUILDING_REQ (size = 133 bytes)
type PROTO_NC_MINIHOUSE_BUILDING_REQ_DATA struct {
	Maxvisiter                       uint8                // offset: 0, size: 1 bytes
	Password                         [9]byte              // offset: 1, size: 9 bytes
	Title                            [21]byte             // offset: 10, size: 21 bytes
	BItemInfoClose                   uint8                // offset: 31, size: 1 bytes
	SNotify                          [101]byte            // offset: 32, size: 101 bytes
}

// PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK (size = 2 bytes)
type PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD (size = 10 bytes)
type PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD_DATA struct {
	PreHandle                        uint16               // offset: 0, size: 2 bytes
	NextHandle                       uint16               // offset: 2, size: 2 bytes
	NFurnitureHnd                    uint16               // offset: 4, size: 2 bytes
	ActorHnd                         uint16               // offset: 6, size: 2 bytes
	ActorIdx                         uint8                // offset: 8, size: 1 bytes
	ActeeIdx                         uint8                // offset: 9, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ (size = 4 bytes)
type PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	NFurnitureHnd                    uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ (size = 10 bytes)
type PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ_DATA struct {
	Hnd                              uint16               // offset: 0, size: 2 bytes
	Location                         uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD_DATA struct {
	Furniturecount                   uint8                // offset: 0, size: 1 bytes
	Rearrange                        [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ_DATA struct {
	Furniturecount                   uint8                // offset: 0, size: 1 bytes
	Rearrange                        [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_ACK (size = 17 bytes)
type PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NPortalItemID                    uint16               // offset: 6, size: 2 bytes
	BType                            uint8                // offset: 8, size: 1 bytes
	DEndDate                         int32                // offset: 9, size: 4 bytes
	Lockindex                        uint16               // offset: 13, size: 2 bytes
	Error                            uint16               // offset: 15, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_REQ (size = 25 bytes)
type PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NPortalItemID                    uint16               // offset: 6, size: 2 bytes
	BType                            uint8                // offset: 8, size: 1 bytes
	DEndDate                         int32                // offset: 9, size: 4 bytes
	Lockindex                        uint16               // offset: 13, size: 2 bytes
	Item                             [10]byte             // offset: 15, size: 10 bytes
}

// PROTO_NC_MINIHOUSE_DB_PORTAL_DEL_ACK (size = 8 bytes)
type PROTO_NC_MINIHOUSE_DB_PORTAL_DEL_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_ACK (size = 10 bytes)
type PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
	NPortalCnt                       uint16               // offset: 8, size: 2 bytes
	PortalInfo                       [0]byte              // offset: 10, flexible array member
}

// PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_REQ (size = 6 bytes)
type PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_ACK (size = 8 bytes)
type PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_ACK_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Error                            uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_REQ (size = 10 bytes)
type PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_REQ_DATA struct {
	Header                           NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NVisiterCharNo                   uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ_DATA struct {
	Furniturecount                   uint8                // offset: 0, size: 1 bytes
	Furniturehandle                  [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_FUNICHER_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_FUNICHER_CMD_DATA struct {
	Furniturecount                   uint8                // offset: 0, size: 1 bytes
	Next                             int8                 // offset: 1, size: 1 bytes
	Array                            [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ_DATA struct {
	Funitureinfocomplete             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK (size = 4 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK_DATA struct {
	NEffectHandle                    uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD (size = 4 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD_DATA struct {
	NEffectHandle                    uint16               // offset: 0, size: 2 bytes
	NApplyPlayerCnt                  uint16               // offset: 2, size: 2 bytes
	NApplyPlayerHandle               [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ (size = 4 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ_DATA struct {
	NFurnitureHnd                    uint16               // offset: 0, size: 2 bytes
	NEffectHandle                    uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK (size = 6 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK_DATA struct {
	NFurnitureHnd                    uint16               // offset: 0, size: 2 bytes
	NPosition                        uint8                // offset: 2, size: 1 bytes
	NEmotionNo                       uint8                // offset: 3, size: 1 bytes
	NError                           uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK (size = 2 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD_DATA struct {
	NCancelEmotionCnt                uint16               // offset: 0, size: 2 bytes
	NCancelEmotion                   [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ_DATA struct {}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD_DATA struct {
	NEmotionCnt                      uint16               // offset: 0, size: 2 bytes
	Emotion                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ (size = 4 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ_DATA struct {
	NFurnitureHnd                    uint16               // offset: 0, size: 2 bytes
	NPosition                        uint8                // offset: 2, size: 1 bytes
	NEmotionNo                       uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_FURNITURE_ENDURE_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_FURNITURE_ENDURE_CMD_DATA struct {
	NEndureFurnitureCnt              uint8                // offset: 0, size: 1 bytes
	EndureFurniture                  [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_KICKOUT_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_KICKOUT_CMD_DATA struct {
	Kick                             uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_KICKOUT_REQ (size = 2 bytes)
type PROTO_NC_MINIHOUSE_KICKOUT_REQ_DATA struct {
	Hnd                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK (size = 140 bytes)
type PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK_DATA struct {
	Error                            int32                // offset: 0, size: 4 bytes
	Maxvisiter                       uint8                // offset: 4, size: 1 bytes
	BItemInfoClose                   uint8                // offset: 5, size: 1 bytes
	SNotify                          [101]byte            // offset: 6, size: 101 bytes
	TargetKey                        [32]byte             // offset: 107, size: 32 bytes
	Playernum                        uint8                // offset: 139, size: 1 bytes
	Player                           [0]byte              // offset: 140, flexible array member
}

// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD (size = 3 bytes)
type PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD_DATA struct {
	Level                            uint8                // offset: 0, size: 1 bytes
	Hnd                              uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD_DATA struct {
	Hnd                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ_DATA struct {
	Maxvisiter                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD_DATA struct {
	NNotifyLen                       uint8                // offset: 0, size: 1 bytes
	SNotify                          [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ_DATA struct {
	NNotifyLen                       uint8                // offset: 0, size: 1 bytes
	SNotify                          [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD (size = 23 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD_DATA struct {
	Mh_hnd                           uint16               // offset: 0, size: 2 bytes
	Title                            [21]byte             // offset: 2, size: 21 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD (size = 1 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD_DATA struct {
	IsPublic                         uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ (size = 9 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ_DATA struct {
	Password                         [9]byte              // offset: 0, size: 9 bytes
}

// PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ (size = 21 bytes)
type PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ_DATA struct {
	Title                            [21]byte             // offset: 0, size: 21 bytes
}

// PROTO_NC_MINIHOUSE_PORTAL_ADD_CMD (size = 6 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_ADD_CMD_DATA struct {
	NewPortal                        [6]byte              // offset: 0, size: 6 bytes
}

// PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ_DATA struct {}

// PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD_DATA struct {
	NPortalItemID                    uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK (size = 2 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD (size = 4 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD_DATA struct {
	NEffectHandle                    uint16               // offset: 0, size: 2 bytes
	NApplyPlayerCnt                  uint16               // offset: 2, size: 2 bytes
	NApplyPlayerHandle               [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ_DATA struct {}

// PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD (size = 2 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD_DATA struct {
	NPortalCnt                       uint16               // offset: 0, size: 2 bytes
	PortalInfo                       [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ (size = 2 bytes)
type PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ_DATA struct {
	NPortalItemID                    uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_REARRANGE_REQ (size = 1 bytes)
type PROTO_NC_MINIHOUSE_REARRANGE_REQ_DATA struct {
	Furniturecount                   uint8                // offset: 0, size: 1 bytes
	Rearrange                        [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_MINIHOUSE_VISITREADY_ACK (size = 7 bytes)
type PROTO_NC_MINIHOUSE_VISITREADY_ACK_DATA struct {
	Error                            int32                // offset: 0, size: 4 bytes
	IsPublic                         uint8                // offset: 4, size: 1 bytes
	Visitercnt                       uint8                // offset: 5, size: 1 bytes
	Maxvisiter                       uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_MINIHOUSE_VISITREADY_REQ (size = 2 bytes)
type PROTO_NC_MINIHOUSE_VISITREADY_REQ_DATA struct {
	Mh_hnd                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_MINIHOUSE_VISIT_REQ (size = 11 bytes)
type PROTO_NC_MINIHOUSE_VISIT_REQ_DATA struct {
	Mh_hnd                           uint16               // offset: 0, size: 2 bytes
	Password                         [9]byte              // offset: 2, size: 9 bytes
}

// PROTO_NC_PET_ASK_NEW_NAME_ACK (size = 4 bytes)
type PROTO_NC_PET_ASK_NEW_NAME_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NPetHandle                       uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_PET_ASK_NEW_NAME_REQ (size = 7 bytes)
type PROTO_NC_PET_ASK_NEW_NAME_REQ_DATA struct {
	NPetHandle                       uint16               // offset: 0, size: 2 bytes
	NPetID                           uint32               // offset: 2, size: 4 bytes
	NFirstSummon                     uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_PET_CREATE_DB_ACK (size = 8 bytes)
type PROTO_NC_PET_CREATE_DB_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NMasterHandle                    uint16               // offset: 2, size: 2 bytes
	NPetRegNum                       uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_PET_CREATE_DB_REQ (size = 24 bytes)
type PROTO_NC_PET_CREATE_DB_REQ_DATA struct {
	NMasterHandle                    uint16               // offset: 0, size: 2 bytes
	NMasterRegNo                     uint32               // offset: 2, size: 4 bytes
	SPetName                         [16]byte             // offset: 6, size: 16 bytes
	NPetMind                         uint8                // offset: 22, size: 1 bytes
	NPetStress                       uint8                // offset: 23, size: 1 bytes
}

// PROTO_NC_PET_LINK_RESUMMON_CMD (size = 23 bytes)
type PROTO_NC_PET_LINK_RESUMMON_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	BDoResummon                      uint8                // offset: 4, size: 1 bytes
	XResummonCode                    uint64               // offset: 5, size: 8 bytes
	NSummonSlot                      uint16               // offset: 13, size: 2 bytes
	NId                              uint32               // offset: 15, size: 4 bytes
	NRegNum                          uint32               // offset: 19, size: 4 bytes
}

// PROTO_NC_PET_LOAD_INFO_DB_ACK (size = 26 bytes)
type PROTO_NC_PET_LOAD_INFO_DB_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NMasterHandle                    uint16               // offset: 2, size: 2 bytes
	NPetRegNum                       uint32               // offset: 4, size: 4 bytes
	SPetName                         [16]byte             // offset: 8, size: 16 bytes
	NPetMind                         uint8                // offset: 24, size: 1 bytes
	NPetStress                       uint8                // offset: 25, size: 1 bytes
}

// PROTO_NC_PET_LOAD_INFO_DB_REQ (size = 6 bytes)
type PROTO_NC_PET_LOAD_INFO_DB_REQ_DATA struct {
	NMasterHandle                    uint16               // offset: 0, size: 2 bytes
	NPetRegNum                       uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PET_SET_NAME_ACK (size = 6 bytes)
type PROTO_NC_PET_SET_NAME_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NPetID                           uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PET_SET_NAME_CANCEL_ACK (size = 6 bytes)
type PROTO_NC_PET_SET_NAME_CANCEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NPetID                           uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PET_SET_NAME_CANCEL_REQ (size = 2 bytes)
type PROTO_NC_PET_SET_NAME_CANCEL_REQ_DATA struct {
	NPetHandle                       uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_PET_SET_NAME_CMD (size = 18 bytes)
type PROTO_NC_PET_SET_NAME_CMD_DATA struct {
	NPetHandle                       uint16               // offset: 0, size: 2 bytes
	SPetName                         [16]byte             // offset: 2, size: 16 bytes
}

// PROTO_NC_PET_SET_NAME_DB_ACK (size = 8 bytes)
type PROTO_NC_PET_SET_NAME_DB_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NMasterHandle                    uint16               // offset: 2, size: 2 bytes
	NPetRegNum                       uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_PET_SET_NAME_DB_REQ (size = 22 bytes)
type PROTO_NC_PET_SET_NAME_DB_REQ_DATA struct {
	NMasterHandle                    uint16               // offset: 0, size: 2 bytes
	NPetRegNum                       uint32               // offset: 2, size: 4 bytes
	SPetName                         [16]byte             // offset: 6, size: 16 bytes
}

// PROTO_NC_PET_SET_NAME_REQ (size = 18 bytes)
type PROTO_NC_PET_SET_NAME_REQ_DATA struct {
	NPetHandle                       uint16               // offset: 0, size: 2 bytes
	SPetName                         [16]byte             // offset: 2, size: 16 bytes
}

// PROTO_NC_PET_SET_TENDENCY_DB_REQ (size = 8 bytes)
type PROTO_NC_PET_SET_TENDENCY_DB_REQ_DATA struct {
	NMasterHandle                    uint16               // offset: 0, size: 2 bytes
	NPetRegNum                       uint32               // offset: 2, size: 4 bytes
	NPetMind                         uint8                // offset: 6, size: 1 bytes
	NPetStress                       uint8                // offset: 7, size: 1 bytes
}
