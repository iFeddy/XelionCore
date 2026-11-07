package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character
// ============================================================================

// PROTO_NC_DATA_PRISON_ADD_GM_ACK (size = 110 bytes)
type PROTO_NC_DATA_PRISON_ADD_GM_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NGmHandler                       uint16               // offset: 2, size: 2 bytes
	NImprisonRegNum                  uint32               // offset: 4, size: 4 bytes
	SImprison                        Name5                // offset: 8, size: 20 bytes
	NMinute                          uint16               // offset: 28, size: 2 bytes
	SReason                          [16]byte             // offset: 30, size: 16 bytes
	SRemark                          [64]byte             // offset: 46, size: 64 bytes
}

// PROTO_NC_DATA_PRISON_ADD_GM_REQ (size = 364 bytes)
type PROTO_NC_DATA_PRISON_ADD_GM_REQ_DATA struct {
	NGmHandler                       uint16               // offset: 0, size: 2 bytes
	SGmUserID                        [256]byte            // offset: 2, size: 256 bytes
	NImprisonRegNum                  uint32               // offset: 258, size: 4 bytes
	SImprison                        Name5                // offset: 262, size: 20 bytes
	NMinute                          uint16               // offset: 282, size: 2 bytes
	SReason                          [16]byte             // offset: 284, size: 16 bytes
	SRemark                          [64]byte             // offset: 300, size: 64 bytes
}

// PROTO_NC_DATA_PRISON_ALTER_GM_ACK (size = 110 bytes)
type PROTO_NC_DATA_PRISON_ALTER_GM_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NGmHandler                       uint16               // offset: 2, size: 2 bytes
	NImprisonRegNum                  uint32               // offset: 4, size: 4 bytes
	SImprison                        Name5                // offset: 8, size: 20 bytes
	NMinute                          uint16               // offset: 28, size: 2 bytes
	SReason                          [16]byte             // offset: 30, size: 16 bytes
	SRemark                          [64]byte             // offset: 46, size: 64 bytes
}

// PROTO_NC_DATA_PRISON_ALTER_GM_REQ (size = 364 bytes)
type PROTO_NC_DATA_PRISON_ALTER_GM_REQ_DATA struct {
	NGmHandler                       uint16               // offset: 0, size: 2 bytes
	SGmUserID                        [256]byte            // offset: 2, size: 256 bytes
	NImprisonRegNum                  uint32               // offset: 258, size: 4 bytes
	SImprison                        Name5                // offset: 262, size: 20 bytes
	NMinute                          uint16               // offset: 282, size: 2 bytes
	SReason                          [16]byte             // offset: 284, size: 16 bytes
	SRemark                          [64]byte             // offset: 300, size: 64 bytes
}

// PROTO_NC_DATA_PRISON_GET_ACK (size = 84 bytes)
type PROTO_NC_DATA_PRISON_GET_ACK_DATA struct {
	NCharHandle                      uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
	SReason                          [16]byte             // offset: 4, size: 16 bytes
	SRemark                          [64]byte             // offset: 20, size: 64 bytes
}

// PROTO_NC_DATA_PRISON_GET_REQ (size = 6 bytes)
type PROTO_NC_DATA_PRISON_GET_REQ_DATA struct {
	NCharHandle                      uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_DATA_PRISON_UPDATE_MIN_CMD (size = 6 bytes)
type PROTO_NC_DATA_PRISON_UPDATE_MIN_CMD_DATA struct {
	NImprisonRegNum                  uint32               // offset: 0, size: 4 bytes
	NMinute                          uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_PRISON_ADD_GM_ACK (size = 22 bytes)
type PROTO_NC_PRISON_ADD_GM_ACK_DATA struct {
	SImprison                        Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PRISON_ADD_GM_REQ (size = 102 bytes)
type PROTO_NC_PRISON_ADD_GM_REQ_DATA struct {
	SImprison                        Name5                // offset: 0, size: 20 bytes
	NMinute                          uint16               // offset: 20, size: 2 bytes
	SReason                          [16]byte             // offset: 22, size: 16 bytes
	SRemark                          [64]byte             // offset: 38, size: 64 bytes
}

// PROTO_NC_PRISON_ALTER_GM_ACK (size = 22 bytes)
type PROTO_NC_PRISON_ALTER_GM_ACK_DATA struct {
	SImprison                        Name5                // offset: 0, size: 20 bytes
	Err                              uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_PRISON_ALTER_GM_REQ (size = 102 bytes)
type PROTO_NC_PRISON_ALTER_GM_REQ_DATA struct {
	SImprison                        Name5                // offset: 0, size: 20 bytes
	NMin                             uint16               // offset: 20, size: 2 bytes
	SReason                          [16]byte             // offset: 22, size: 16 bytes
	SRemark                          [64]byte             // offset: 38, size: 64 bytes
}

// PROTO_NC_PRISON_END_ACK (size = 6 bytes)
type PROTO_NC_PRISON_END_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NPrisonSecond                    uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_PRISON_END_REQ (size = 1 bytes)
type PROTO_NC_PRISON_END_REQ_DATA struct {}

// PROTO_NC_PRISON_GET_ACK (size = 84 bytes)
type PROTO_NC_PRISON_GET_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NMinute                          uint16               // offset: 2, size: 2 bytes
	SReason                          [16]byte             // offset: 4, size: 16 bytes
	SRemark                          [64]byte             // offset: 20, size: 64 bytes
}

// PROTO_NC_PRISON_OK_CMD (size = 82 bytes)
type PROTO_NC_PRISON_OK_CMD_DATA struct {
	NMinute                          uint16               // offset: 0, size: 2 bytes
	SReason                          [16]byte             // offset: 2, size: 16 bytes
	SRemark                          [64]byte             // offset: 18, size: 64 bytes
}

// PROTO_NC_ZONE_PRISON_END_CMD (size = 6 bytes)
type PROTO_NC_ZONE_PRISON_END_CMD_DATA struct {
	NCharZoneHandle                  uint16               // offset: 0, size: 2 bytes
	NCharRegNum                      uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_ZONE_PRISON_GO_ACK (size = 374 bytes)
type PROTO_NC_ZONE_PRISON_GO_ACK_DATA struct {
	NGmHandler                       uint16               // offset: 0, size: 2 bytes
	NGmRegNum                        uint32               // offset: 2, size: 4 bytes
	SGmUserID                        [256]byte            // offset: 6, size: 256 bytes
	SCharID                          Name5                // offset: 262, size: 20 bytes
	NCharSessHandler                 uint16               // offset: 282, size: 2 bytes
	NCharZoneHandle                  uint16               // offset: 284, size: 2 bytes
	NCharRegNum                      uint32               // offset: 286, size: 4 bytes
	NMinute                          uint16               // offset: 290, size: 2 bytes
	SReason                          [16]byte             // offset: 292, size: 16 bytes
	SRemark                          [64]byte             // offset: 308, size: 64 bytes
	Err                              uint16               // offset: 372, size: 2 bytes
}

// PROTO_NC_ZONE_PRISON_GO_REQ (size = 372 bytes)
type PROTO_NC_ZONE_PRISON_GO_REQ_DATA struct {
	NGmHandler                       uint16               // offset: 0, size: 2 bytes
	NGmRegNum                        uint32               // offset: 2, size: 4 bytes
	SGmUserID                        [256]byte            // offset: 6, size: 256 bytes
	SCharID                          Name5                // offset: 262, size: 20 bytes
	NCharSessHandle                  uint16               // offset: 282, size: 2 bytes
	NCharZoneHandle                  uint16               // offset: 284, size: 2 bytes
	NCharRegNum                      uint32               // offset: 286, size: 4 bytes
	NMinute                          uint16               // offset: 290, size: 2 bytes
	SReason                          [16]byte             // offset: 292, size: 16 bytes
	SRemark                          [64]byte             // offset: 308, size: 64 bytes
}
