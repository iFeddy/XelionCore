package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_briefinfo.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD (sizeof = 14 bytes)
type PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Info                             [12]byte             // offset: 2, size: 12 bytes
}

func (s *PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD (sizeof = 3 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD_DATA)

// PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD (sizeof = 2 bytes)
type PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD_DATA struct {
	Hnd                              uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BRIEFINFO_BUILDDOOR_CMD (sizeof = 48 bytes)
type PROTO_NC_BRIEFINFO_BUILDDOOR_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Mobid                            uint16               // offset: 2, size: 2 bytes
	Coord                            [9]byte              // offset: 4, size: 9 bytes
	Doorstate                        uint8                // offset: 13, size: 1 bytes
	Blockindex                       Name8                // offset: 14, size: 32 bytes
	Scale                            uint16               // offset: 46, size: 2 bytes
}

func (s *PROTO_NC_BRIEFINFO_BUILDDOOR_CMD_DATA) FixedWireSize() int { return 48 }

// PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD (sizeof = 5 bytes)
type PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Item                             uint16               // offset: 2, size: 2 bytes
	NSlotNum                         uint8                // offset: 4, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD (sizeof = 6 bytes)
type PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Item                             uint16               // offset: 2, size: 2 bytes
	Upgrade                          uint8                // offset: 4, size: 1 bytes
	NSlotNum                         uint8                // offset: 5, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD (sizeof = 9 bytes)
type PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD_DATA struct {
	Upgradeinfo                      [6]byte              // offset: 0, size: 6 bytes
	Currentmobid                     uint16               // offset: 6, size: 2 bytes
	Currentkilllevel                 uint8                // offset: 8, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD_DATA) FixedWireSize() int { return 9 }

// PROTO_NC_BRIEFINFO_CHARACTER_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_CHARACTER_CMD_DATA struct {
	Charnum                          uint8                // offset: 0, size: 1 bytes
	Chars                            [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_CHARACTER_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_DOOR_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_DOOR_CMD_DATA struct {
	Doornum                          uint8                // offset: 0, size: 1 bytes
	Doors                            [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_DOOR_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (sizeof = 15 bytes)
type PROTO_NC_BRIEFINFO_DROPEDITEM_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Itemid                           uint16               // offset: 2, size: 2 bytes
	Location                         SHINE_XY_TYPE        // offset: 4, size: 8 bytes
	Dropmobhandle                    uint16               // offset: 12, size: 2 bytes
	Attr                             uint8                // offset: 14, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_DROPEDITEM_CMD_DATA) FixedWireSize() int { return 15 }

// PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD (sizeof = 48 bytes)
type PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Effectname                       Name8                // offset: 2, size: 32 bytes
	Coord                            [9]byte              // offset: 34, size: 9 bytes
	Detach                           uint16               // offset: 43, size: 2 bytes
	Scale                            uint16               // offset: 45, size: 2 bytes
	Flag                             uint8                // offset: 47, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD_DATA) FixedWireSize() int { return 48 }

// PROTO_NC_BRIEFINFO_EFFECT_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_EFFECT_CMD_DATA struct {
	Effectnum                        uint8                // offset: 0, size: 1 bytes
	Effects                          [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_EFFECT_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_INFORM_CMD (sizeof = 6 bytes)
type PROTO_NC_BRIEFINFO_INFORM_CMD_DATA struct {
	NMyHnd                           uint16               // offset: 0, size: 2 bytes
	ReceiveNetCommand                uint16               // offset: 2, size: 2 bytes
	Hnd                              uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BRIEFINFO_INFORM_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD_DATA struct {
	Itemnum                          uint8                // offset: 0, size: 1 bytes
	Items                            [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD (sizeof = 235 bytes)
type PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Charid                           Name5                // offset: 2, size: 20 bytes
	Coord                            [9]byte              // offset: 22, size: 9 bytes
	Mode                             uint8                // offset: 31, size: 1 bytes
	Chrclass                         uint8                // offset: 32, size: 1 bytes
	Shape                            uint32               // offset: 33, size: 4 bytes
	Shapedata                        [45]byte             // offset: 37, size: 45 bytes
	Polymorph                        uint16               // offset: 82, size: 2 bytes
	Emoticon                         [3]byte              // offset: 84, size: 3 bytes
	Chartitle                        uint32               // offset: 87, size: 4 bytes
	Abstatebit                       [99]byte             // offset: 91, size: 99 bytes
	Myguild                          uint32               // offset: 190, size: 4 bytes
	TypeField                        uint8                // offset: 194, size: 1 bytes
	IsGuildAcademyMember             uint8                // offset: 195, size: 1 bytes
	IsAutoPick                       uint8                // offset: 196, size: 1 bytes
	Level                            uint8                // offset: 197, size: 1 bytes
	SAnimation                       [32]byte             // offset: 198, size: 32 bytes
	NMoverHnd                        uint16               // offset: 230, size: 2 bytes
	NMoverSlot                       uint8                // offset: 232, size: 1 bytes
	NKQTeamType                      uint8                // offset: 233, size: 1 bytes
	IsUseItemMinimon                 uint8                // offset: 234, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD_DATA) FixedWireSize() int { return 235 }

// PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD_DATA struct {
	Magicfields                      uint8                // offset: 0, size: 1 bytes
	Magicfield                       [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD (sizeof = 16 bytes)
type PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Caster                           uint16               // offset: 2, size: 2 bytes
	Skillid                          uint16               // offset: 4, size: 2 bytes
	Location                         SHINE_XY_TYPE        // offset: 6, size: 8 bytes
	Radious                          uint16               // offset: 14, size: 2 bytes
}

func (s *PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_BRIEFINFO_MINIHOUSEBUILD_CMD (sizeof = 64 bytes)
type PROTO_NC_BRIEFINFO_MINIHOUSEBUILD_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Camp                             [12]byte             // offset: 2, size: 12 bytes
	Charid                           Name5                // offset: 14, size: 20 bytes
	Coord                            [9]byte              // offset: 34, size: 9 bytes
	Title                            [21]byte             // offset: 43, size: 21 bytes
}

func (s *PROTO_NC_BRIEFINFO_MINIHOUSEBUILD_CMD_DATA) FixedWireSize() int { return 64 }

// PROTO_NC_BRIEFINFO_MINIHOUSE_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_MINIHOUSE_CMD_DATA struct {
	Minihouses                       uint8                // offset: 0, size: 1 bytes
	Minihouse                        [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_MINIHOUSE_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_MOB_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_MOB_CMD_DATA struct {
	Mobnum                           uint8                // offset: 0, size: 1 bytes
	Mobs                             [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_MOB_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_MOVER_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_MOVER_CMD_DATA struct {
	NMoverNum                        uint8                // offset: 0, size: 1 bytes
	Movers                           [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_MOVER_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_PET_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_PET_CMD_DATA struct {
	NPetCnt                          uint8                // offset: 0, size: 1 bytes
	Pets                             [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_PET_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_PLAYER_INFO_APPEAR_CMD (sizeof = 30 bytes)
type PROTO_NC_BRIEFINFO_PLAYER_INFO_APPEAR_CMD struct {
	NHnd                             uint16               // offset: 0, size: 2 bytes
	SID                              Name5                // offset: 2, size: 20 bytes
	CharTitle                        uint32               // offset: 22, size: 4 bytes
	NGuildID                         uint32               // offset: 26, size: 4 bytes
}

func (s *PROTO_NC_BRIEFINFO_PLAYER_INFO_APPEAR_CMD) FixedWireSize() int { return 30 }

// PROTO_NC_BRIEFINFO_PLAYER_LIST_INFO_APPEAR_CMD (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_PLAYER_LIST_INFO_APPEAR_CMD_DATA struct {
	NPlayerCnt                       uint8                // offset: 0, size: 1 bytes
	PlayerInfo                       [0]byte              // offset: 1, flexible array member
}

func (s *PROTO_NC_BRIEFINFO_PLAYER_LIST_INFO_APPEAR_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_REGENMOB_CMD (sizeof = 149 bytes)
type PROTO_NC_BRIEFINFO_REGENMOB_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Mode                             uint8                // offset: 2, size: 1 bytes
	Mobid                            uint16               // offset: 3, size: 2 bytes
	Coord                            [9]byte              // offset: 5, size: 9 bytes
	Flagstate                        uint8                // offset: 14, size: 1 bytes
	Flag                             [99]byte             // offset: 15, size: 99 bytes
	SAnimation                       [32]byte             // offset: 114, size: 32 bytes
	NAnimationLevel                  uint8                // offset: 146, size: 1 bytes
	NKQTeamType                      uint8                // offset: 147, size: 1 bytes
	BRegenAni                        uint8                // offset: 148, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_REGENMOB_CMD_DATA) FixedWireSize() int { return 149 }

// PROTO_NC_BRIEFINFO_REGENMOVER_CMD (sizeof = 139 bytes)
type PROTO_NC_BRIEFINFO_REGENMOVER_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NID                              uint32               // offset: 2, size: 4 bytes
	NHP                              uint32               // offset: 6, size: 4 bytes
	NCoord                           [9]byte              // offset: 10, size: 9 bytes
	AbstateBit                       [99]byte             // offset: 19, size: 99 bytes
	NGrade                           uint8                // offset: 118, size: 1 bytes
	NSlotHandle                      [10]uint16           // offset: 119, size: 20 bytes
}

func (s *PROTO_NC_BRIEFINFO_REGENMOVER_CMD_DATA) FixedWireSize() int { return 139 }

// PROTO_NC_BRIEFINFO_REGENPET_CMD (sizeof = 63 bytes)
type PROTO_NC_BRIEFINFO_REGENPET_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NID                              uint32               // offset: 2, size: 4 bytes
	XYDir                            [9]byte              // offset: 6, size: 9 bytes
	SPetName                         [16]byte             // offset: 15, size: 16 bytes
	SAnimationIndex                  [32]byte             // offset: 31, size: 32 bytes
}

func (s *PROTO_NC_BRIEFINFO_REGENPET_CMD_DATA) FixedWireSize() int { return 63 }

// PROTO_NC_BRIEFINFO_UNEQUIP_CMD (sizeof = 3 bytes)
type PROTO_NC_BRIEFINFO_UNEQUIP_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Slot                             uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_UNEQUIP_CMD_DATA) FixedWireSize() int { return 3 }
