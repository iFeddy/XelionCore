package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_act.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_ACT_ACTIONBYITEM_ACK (sizeof = 3 bytes)
type PROTO_NC_ACT_ACTIONBYITEM_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	Emoticonid                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ACT_ACTIONBYITEM_ACK_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_ACT_ACTIONBYITEM_REQ (sizeof = 1 bytes)
type PROTO_NC_ACT_ACTIONBYITEM_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_ACTIONBYITEM_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_ANIMATION_LEVEL_CHANGE_CMD (sizeof = 3 bytes)
type PROTO_NC_ACT_ANIMATION_LEVEL_CHANGE_CMD_DATA struct {
	NMonsterHandle                   uint16               // offset: 0, size: 2 bytes
	NAnimationLevel                  uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ACT_ANIMATION_LEVEL_CHANGE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_ACT_ANIMATION_START_CMD (sizeof = 34 bytes)
type PROTO_NC_ACT_ANIMATION_START_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	SAnimationIndex                  [32]byte             // offset: 2, size: 32 bytes
}

func (s *PROTO_NC_ACT_ANIMATION_START_CMD_DATA) FixedWireSize() int { return 34 }

// PROTO_NC_ACT_ANIMATION_STOP_CMD (sizeof = 34 bytes)
type PROTO_NC_ACT_ANIMATION_STOP_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	SAnimationIndex                  [32]byte             // offset: 2, size: 32 bytes
}

func (s *PROTO_NC_ACT_ANIMATION_STOP_CMD_DATA) FixedWireSize() int { return 34 }

// PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK (sizeof = 2 bytes)
type PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_CHANGEMODE_REQ (sizeof = 1 bytes)
type PROTO_NC_ACT_CHANGEMODE_REQ_DATA struct {
	Mode                             uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_CHANGEMODE_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_CHAT_REQ (sizeof = 2 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ACT_CHAT_REQ_DATA)

// PROTO_NC_ACT_CREATECASTBAR (sizeof = 2 bytes)
type PROTO_NC_ACT_CREATECASTBAR_DATA struct {
	Millisec                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_CREATECASTBAR_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_EFFECT_MESSAGE_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_EFFECT_MESSAGE_CMD_DATA struct {
	NEffectEnum                      uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_ACT_EFFECT_MESSAGE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_EMOTICON_CMD (sizeof = 1 bytes)
type PROTO_NC_ACT_EMOTICON_CMD_DATA struct {
	Emoticonid                       uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_EMOTICON_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_EVENT_CODE_ACTION_CMD (sizeof = 6 bytes)
type PROTO_NC_ACT_EVENT_CODE_ACTION_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NActionEventCode                 uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_ACT_EVENT_CODE_ACTION_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_ACT_GATHERCOMPLETE_ACK (sizeof = 4 bytes)
type PROTO_NC_ACT_GATHERCOMPLETE_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Object                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_GATHERCOMPLETE_ACK_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_GATHERSTART_ACK (sizeof = 6 bytes)
type PROTO_NC_ACT_GATHERSTART_ACK_DATA struct {
	Picker                           uint16               // offset: 0, size: 2 bytes
	Object                           uint16               // offset: 2, size: 2 bytes
	Err                              uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_ACT_GATHERSTART_ACK_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_ACT_GATHERSTART_REQ (sizeof = 2 bytes)
type PROTO_NC_ACT_GATHERSTART_REQ_DATA struct {
	Objhandle                        uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_GATHERSTART_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_MOVEFAIL_CMD (sizeof = 8 bytes)
type PROTO_NC_ACT_MOVEFAIL_CMD_DATA struct {
	Back                             SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_ACT_MOVEFAIL_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_ACT_MOVESPEED_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_MOVESPEED_CMD_DATA struct {
	Walkspeed                        uint16               // offset: 0, size: 2 bytes
	Runspeed                         uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_MOVESPEED_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_MOVEWALK_CMD (sizeof = 16 bytes)
type PROTO_NC_ACT_MOVEWALK_CMD_DATA struct {
	From                             SHINE_XY_TYPE        // offset: 0, size: 8 bytes
	To                               SHINE_XY_TYPE        // offset: 8, size: 8 bytes
}

func (s *PROTO_NC_ACT_MOVEWALK_CMD_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_ACT_NPCCLICK_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_NPCCLICK_CMD_DATA struct {
	Npchandle                        uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_NPCCLICK_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_NPCMENUOPEN_ACK (sizeof = 1 bytes)
type PROTO_NC_ACT_NPCMENUOPEN_ACK_DATA struct {
	Ack                              uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_NPCMENUOPEN_ACK_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_NPCMENUOPEN_REQ (sizeof = 2 bytes)
type PROTO_NC_ACT_NPCMENUOPEN_REQ_DATA struct {
	Mobid                            uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_NPCMENUOPEN_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_NPC_ACTION_CMD (sizeof = 7 bytes)
type PROTO_NC_ACT_NPC_ACTION_CMD_DATA struct {
	NType                            uint8                // offset: 0, size: 1 bytes
	NNPCHandle                       uint16               // offset: 1, size: 2 bytes
	NECode                           uint32               // offset: 3, size: 4 bytes
}

func (s *PROTO_NC_ACT_NPC_ACTION_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_ACT_NPC_MENU_CMD (sizeof = 1 bytes)
type PROTO_NC_ACT_NPC_MENU_CMD_DATA struct {
	Ack                              uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_NPC_MENU_CMD_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_OBJECT_EFFECT_CMD (sizeof = 6 bytes)
type PROTO_NC_ACT_OBJECT_EFFECT_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	NEffectID                        uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_ACT_OBJECT_EFFECT_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_ACT_OBJECT_SOUND_CMD (sizeof = 34 bytes)
type PROTO_NC_ACT_OBJECT_SOUND_CMD_DATA struct {
	NHandle                          uint16               // offset: 0, size: 2 bytes
	SSoundFileName                   [32]byte             // offset: 2, size: 32 bytes
}

func (s *PROTO_NC_ACT_OBJECT_SOUND_CMD_DATA) FixedWireSize() int { return 34 }

// PROTO_NC_ACT_PLAY_SOUND_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_PLAY_SOUND_CMD_DATA struct {
	FileNameLen                      uint16               // offset: 0, size: 2 bytes
	FileName                         [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_ACT_PLAY_SOUND_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_PRODUCE_CAST_REQ (sizeof = 2 bytes)
type PROTO_NC_ACT_PRODUCE_CAST_REQ_DATA struct {
	Produceskill                     uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_PRODUCE_CAST_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_PRODUCE_MAKE_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_PRODUCE_MAKE_CMD_DATA struct {
	Item                             uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_PRODUCE_MAKE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_REINFORCEMOVEBYPATH_CMD (sizeof = 6 bytes)
type PROTO_NC_ACT_REINFORCEMOVEBYPATH_CMD_DATA struct {
	Mover                            uint16               // offset: 0, size: 2 bytes
	Speed                            uint16               // offset: 2, size: 2 bytes
	Direct                           uint8                // offset: 4, size: 1 bytes
	Pathlength                       uint8                // offset: 5, size: 1 bytes
	Path                             [0]byte              // offset: 6, flexible array member
}

func (s *PROTO_NC_ACT_REINFORCEMOVEBYPATH_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_ACT_REINFORCE_RUN_CMD (sizeof = 8 bytes)
type PROTO_NC_ACT_REINFORCE_RUN_CMD_DATA struct {
	Xy                               SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_ACT_REINFORCE_RUN_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_ACT_REINFORCE_STOP_CMD (sizeof = 8 bytes)
type PROTO_NC_ACT_REINFORCE_STOP_CMD_DATA struct {
	Loc                              SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_ACT_REINFORCE_STOP_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_ACT_RIDE_FEEDING_REQ (sizeof = 1 bytes)
type PROTO_NC_ACT_RIDE_FEEDING_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_RIDE_FEEDING_REQ_DATA) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_RIDE_HUNGRY_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_RIDE_HUNGRY_CMD_DATA struct {
	Hungry                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_RIDE_HUNGRY_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_RIDE_ON_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_RIDE_ON_CMD_DATA struct {
	Ride                             uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_RIDE_ON_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_ROAR_REQ (sizeof = 2 bytes)
type PROTO_NC_ACT_ROAR_REQ_DATA struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_ACT_ROAR_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SCRIPT_MSG_CMD (sizeof = 39 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ACT_SCRIPT_MSG_CMD_DATA)

// PROTO_NC_ACT_SCRIPT_MSG_WORLD_CMD (sizeof = 39 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_ACT_SCRIPT_MSG_WORLD_CMD_DATA)

// PROTO_NC_ACT_SETITEMHEALEFFECT (sizeof = 5 bytes)
type PROTO_NC_ACT_SETITEMHEALEFFECT_DATA struct {
	Who                              uint16               // offset: 0, size: 2 bytes
	Flag                             uint8                // offset: 2, size: 1 bytes
	Hpchange                         uint16               // offset: 3, size: 2 bytes
}

func (s *PROTO_NC_ACT_SETITEMHEALEFFECT_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_ACT_SHOUT_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_SHOUT_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_ACT_SHOUT_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SHOW_CINEMATIC_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_SHOW_CINEMATIC_CMD struct {
	PTextLen                         uint16               // offset: 0, size: 2 bytes
	PText                            [0]byte              // offset: 2, flexible array member
}

func (s *PROTO_NC_ACT_SHOW_CINEMATIC_CMD) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SOMEEONEJUMP_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_SOMEEONEJUMP_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEEONEJUMP_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SOMEONECHANGEMODE_CMD (sizeof = 3 bytes)
type PROTO_NC_ACT_SOMEONECHANGEMODE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Mode                             uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ACT_SOMEONECHANGEMODE_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_ACT_SOMEONECHAT_CMD (sizeof = 7 bytes)
type PROTO_NC_ACT_SOMEONECHAT_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Handle                           uint16               // offset: 1, size: 2 bytes
	Len                              uint8                // offset: 3, size: 1 bytes
	Flag                             uint8                // offset: 4, size: 1 bytes
	NChatFontColorID                 uint8                // offset: 5, size: 1 bytes
	NChatBalloonColorID              uint8                // offset: 6, size: 1 bytes
	Content                          [0]byte              // offset: 7, flexible array member
}

func (s *PROTO_NC_ACT_SOMEONECHAT_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD (sizeof = 5 bytes)
type PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Emoticon                         [3]byte              // offset: 2, size: 3 bytes
}

func (s *PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_ACT_SOMEONEEMOTICON_CMD (sizeof = 3 bytes)
type PROTO_NC_ACT_SOMEONEEMOTICON_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Emoticonid                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ACT_SOMEONEEMOTICON_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_ACT_SOMEONEFOLDTENT_CMD (sizeof = 45 bytes)
type PROTO_NC_ACT_SOMEONEFOLDTENT_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Shape                            [43]byte             // offset: 2, size: 43 bytes
}

func (s *PROTO_NC_ACT_SOMEONEFOLDTENT_CMD_DATA) FixedWireSize() int { return 45 }

// PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	Object                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	Object                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_SOMEONEGATHERSTART_CMD (sizeof = 6 bytes)
type PROTO_NC_ACT_SOMEONEGATHERSTART_CMD_DATA struct {
	Player                           uint16               // offset: 0, size: 2 bytes
	Object                           uint16               // offset: 2, size: 2 bytes
	Picker                           uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEGATHERSTART_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_ACT_SOMEONEMOVEWALK_CMD (sizeof = 22 bytes)
type PROTO_NC_ACT_SOMEONEMOVEWALK_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	From                             SHINE_XY_TYPE        // offset: 2, size: 8 bytes
	To                               SHINE_XY_TYPE        // offset: 10, size: 8 bytes
	Speed                            uint16               // offset: 18, size: 2 bytes
	Moveattr                         uint16               // offset: 20, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEMOVEWALK_CMD_DATA) FixedWireSize() int { return 22 }

// PROTO_NC_ACT_SOMEONEPITCHTENT_CMD (sizeof = 14 bytes)
type PROTO_NC_ACT_SOMEONEPITCHTENT_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Camp                             [12]byte             // offset: 2, size: 12 bytes
}

func (s *PROTO_NC_ACT_SOMEONEPITCHTENT_CMD_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_ACT_SOMEONEPRODUCE_CASTCUT_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_SOMEONEPRODUCE_CASTCUT_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEPRODUCE_CASTCUT_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
	Item                             uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
	Item                             uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_SOMEONERIDE_OFF_CMD (sizeof = 2 bytes)
type PROTO_NC_ACT_SOMEONERIDE_OFF_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONERIDE_OFF_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SOMEONERIDE_ON_CMD (sizeof = 4 bytes)
type PROTO_NC_ACT_SOMEONERIDE_ON_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Ride                             uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONERIDE_ON_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_ACT_SOMEONESHOUT_CMD (sizeof = 23 bytes)
type PROTO_NC_ACT_SOMEONESHOUT_CMD_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Speaker                          [20]byte             // offset: 1, size: 20 bytes
	Flag                             uint8                // offset: 21, size: 1 bytes
	Len                              uint8                // offset: 22, size: 1 bytes
	Content                          [0]byte              // offset: 23, flexible array member
}

func (s *PROTO_NC_ACT_SOMEONESHOUT_CMD_DATA) FixedWireSize() int { return 23 }

// PROTO_NC_ACT_SOMEONESTOP_CMD (sizeof = 10 bytes)
type PROTO_NC_ACT_SOMEONESTOP_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Loc                              SHINE_XY_TYPE        // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_ACT_SOMEONESTOP_CMD_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_ACT_SOMEONEWALK_CMD (sizeof = 18 bytes)
type PROTO_NC_ACT_SOMEONEWALK_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	From                             SHINE_XY_TYPE        // offset: 2, size: 8 bytes
	To                               SHINE_XY_TYPE        // offset: 10, size: 8 bytes
}

func (s *PROTO_NC_ACT_SOMEONEWALK_CMD_DATA) FixedWireSize() int { return 18 }

// PROTO_NC_ACT_STOP_REQ (sizeof = 8 bytes)
type PROTO_NC_ACT_STOP_REQ_DATA struct {
	Loc                              SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_ACT_STOP_REQ_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_ACT_WALK_REQ (sizeof = 16 bytes)
type PROTO_NC_ACT_WALK_REQ_DATA struct {
	From                             SHINE_XY_TYPE        // offset: 0, size: 8 bytes
	To                               SHINE_XY_TYPE        // offset: 8, size: 8 bytes
}

func (s *PROTO_NC_ACT_WALK_REQ_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_ACT_WEDDING_COUPLE_ENTRANCE_RNG (sizeof = 16 bytes)
type PROTO_NC_ACT_WEDDING_COUPLE_ENTRANCE_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Groom                            uint32               // offset: 7, size: 4 bytes
	Bride                            uint32               // offset: 11, size: 4 bytes
	RecallCount                      uint8                // offset: 15, size: 1 bytes
}

func (s *PROTO_NC_ACT_WEDDING_COUPLE_ENTRANCE_RNG_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK (sizeof = 2 bytes)
type PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_ACK (sizeof = 174 bytes)
type PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	GroomID                          Name5                // offset: 2, size: 20 bytes
	BrideID                          Name5                // offset: 22, size: 20 bytes
	EnterStart                       int64                // offset: 42, size: 8 bytes
	WeddingStart                     int64                // offset: 50, size: 8 bytes
	WeddingEnd                       int64                // offset: 58, size: 8 bytes
	Tm_EnterStart                    [36]byte             // offset: 66, size: 36 bytes
	Tm_WeddingStart                  [36]byte             // offset: 102, size: 36 bytes
	Tm_WeddingEnd                    [36]byte             // offset: 138, size: 36 bytes
}

func (s *PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_ACK_DATA) FixedWireSize() int { return 174 }

// PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK (sizeof = 46 bytes)
type PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	ReservedTime                     int64                // offset: 2, size: 8 bytes
	Tm_ReservedTime                  [36]byte             // offset: 10, size: 36 bytes
}

func (s *PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK_DATA) FixedWireSize() int { return 46 }

// PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK (sizeof = 40 bytes)
type PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK_DATA struct {
	From                             uint16               // offset: 0, size: 2 bytes
	Reply                            uint8                // offset: 2, size: 1 bytes
	Response_word                    [37]byte             // offset: 3, size: 37 bytes
}

func (s *PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK_DATA) FixedWireSize() int { return 40 }

// PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ (sizeof = 39 bytes)
type PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ_DATA struct {
	From                             uint16               // offset: 0, size: 2 bytes
	Propose_word                     [37]byte             // offset: 2, size: 37 bytes
}

func (s *PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ_DATA) FixedWireSize() int { return 39 }

// PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK (sizeof = 41 bytes)
type PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK_DATA struct {
	To                               uint16               // offset: 0, size: 2 bytes
	Err                              uint16               // offset: 2, size: 2 bytes
	Response_word                    [37]byte             // offset: 4, size: 37 bytes
}

func (s *PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK_DATA) FixedWireSize() int { return 41 }

// PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ (sizeof = 39 bytes)
type PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ_DATA struct {
	To                               uint16               // offset: 0, size: 2 bytes
	Propose_word                     [37]byte             // offset: 2, size: 37 bytes
}

func (s *PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ_DATA) FixedWireSize() int { return 39 }
