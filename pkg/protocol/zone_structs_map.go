package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_map.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD (sizeof = 1 bytes)
type PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD_DATA struct {
	BCanUseReviveItem                uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_MAP_EXPBONUS_RNG (sizeof = 9 bytes)
type PROTO_NC_MAP_EXPBONUS_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Bonus                            uint16               // offset: 7, size: 2 bytes
}

func (s *PROTO_NC_MAP_EXPBONUS_RNG_DATA) FixedWireSize() int { return 9 }

// PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD (sizeof = 4 bytes)
type PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD_DATA struct {
	EFieldMapType                    uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_MAP_INDUN_LEVEL_VIEW_CMD (sizeof = 4 bytes)
type PROTO_NC_MAP_INDUN_LEVEL_VIEW_CMD_DATA struct {
	ELevelType                       uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_MAP_INDUN_LEVEL_VIEW_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_MAP_ITEMBONUS_RNG (sizeof = 9 bytes)
type PROTO_NC_MAP_ITEMBONUS_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Bonus                            uint16               // offset: 7, size: 2 bytes
}

func (s *PROTO_NC_MAP_ITEMBONUS_RNG_DATA) FixedWireSize() int { return 9 }

// PROTO_NC_MAP_LINKEND_CMD (sizeof = 47 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_LINKEND_CMD_DATA)

// PROTO_NC_MAP_LINKRESERVE_ACK (sizeof = 36 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_LINKRESERVE_ACK_DATA)

// PROTO_NC_MAP_LINKRESERVE_REQ (sizeof = 59 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_LINKRESERVE_REQ_DATA)

// PROTO_NC_MAP_LINKSTART_CMD (sizeof = 6 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_LINKSTART_CMD_DATA)

// PROTO_NC_MAP_LOGINFAIL_ACK (sizeof = 3 bytes)
type PROTO_NC_MAP_LOGINFAIL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NWrongDataFileIndex              uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_MAP_LOGINFAIL_ACK_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_MAP_LOGIN_REQ (sizeof = 1590 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_LOGIN_REQ_DATA)

// PROTO_NC_MAP_LOGOUT_CMD (sizeof = 2 bytes)
type PROTO_NC_MAP_LOGOUT_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_MAP_LOGOUT_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_MAP_MULTY_LINK_CMD (sizeof = 73 bytes)
type PROTO_NC_MAP_MULTY_LINK_CMD_DATA struct {
	NpcHandle                        uint16               // offset: 0, size: 2 bytes
	NpcPosition                      SHINE_XY_TYPE        // offset: 2, size: 8 bytes
	LimitRange                       uint16               // offset: 10, size: 2 bytes
	Num                              uint8                // offset: 12, size: 1 bytes
	LinkMapName                      [60]byte             // offset: 13, size: 60 bytes
}

func (s *PROTO_NC_MAP_MULTY_LINK_CMD_DATA) FixedWireSize() int { return 73 }

// PROTO_NC_MAP_MULTY_LINK_SELECT_ACK (sizeof = 2 bytes)
type PROTO_NC_MAP_MULTY_LINK_SELECT_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_MAP_MULTY_LINK_SELECT_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_MAP_MULTY_LINK_SELECT_REQ (sizeof = 12 bytes)
type PROTO_NC_MAP_MULTY_LINK_SELECT_REQ_DATA struct {
	LinkMapName                      Name3                // offset: 0, size: 12 bytes
}

func (s *PROTO_NC_MAP_MULTY_LINK_SELECT_REQ_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_MAP_REGIST_CMD (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_MAP_REGIST_CMD_DATA)

// PROTO_NC_MAP_TOWNPORTAL_REQ (sizeof = 1 bytes)
type PROTO_NC_MAP_TOWNPORTAL_REQ_DATA struct {
	Portalindex                      uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_MAP_TOWNPORTAL_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_MAP_WING_FLY_REQ (sizeof = 1 bytes)
type PROTO_NC_MAP_WING_FLY_REQ_DATA struct {
	Where                            uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_MAP_WING_FLY_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_MAP_WING_SAVE_REQ (sizeof = 1 bytes)
type PROTO_NC_MAP_WING_SAVE_REQ_DATA struct {
	Where                            uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_MAP_WING_SAVE_REQ_DATA) FixedWireSize() int { return 1 }
