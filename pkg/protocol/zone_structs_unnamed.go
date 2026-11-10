package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_unnamed.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_ACT_SETITEMHEALEFFECT::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_ACT_SETITEMHEALEFFECT_unnamed_type_flag struct {
	Isheal                           uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_SETITEMHEALEFFECT_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_SOMEONECHAT_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_ACT_SOMEONECHAT_CMD_unnamed_type_flag struct {
	GMColor                          [0]byte              // offset: 0, size: 0 bytes
	Chatwin                          uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_SOMEONECHAT_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr> (sizeof = 2 bytes)
type PROTO_NC_ACT_SOMEONEMOVEWALK_CMD_unnamed_type_moveattr struct {
	Direct                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ACT_SOMEONEMOVEWALK_CMD_unnamed_type_moveattr) FixedWireSize() int { return 2 }

// PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_ACT_SOMEONESHOUT_CMD_unnamed_type_flag struct {
	GMColor                          [0]byte              // offset: 0, size: 0 bytes
	IsMob                            uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ACT_SOMEONESHOUT_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker> (sizeof = 20 bytes)
type PROTO_NC_ACT_SOMEONESHOUT_CMD_unnamed_type_speaker struct {
	CharID                           [0]byte              // offset: 0, size: 0 bytes
	MobID                            [20]byte             // offset: 0, size: 20 bytes
}

func (s *PROTO_NC_ACT_SOMEONESHOUT_CMD_unnamed_type_speaker) FixedWireSize() int { return 20 }

// PROTO_NC_BAT_AREADOTDAMAGE_CMD::targetinfo (sizeof = 12 bytes)
type PROTO_NC_BAT_AREADOTDAMAGE_CMD_targetinfo struct {
	Object                           uint16               // offset: 0, size: 2 bytes
	Resthp                           uint32               // offset: 2, size: 4 bytes
	Damage                           uint32               // offset: 6, size: 4 bytes
	Hpchangeorder                    uint16               // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_BAT_AREADOTDAMAGE_CMD_targetinfo) FixedWireSize() int { return 12 }

// PROTO_NC_BAT_BASH_HITTED_CMD::<unnamed-type-damage> (sizeof = 2 bytes)
type PROTO_NC_BAT_BASH_HITTED_CMD_unnamed_type_damage struct {
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Damage                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_HITTED_CMD_unnamed_type_damage) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_BASH_HITTED_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BAT_BASH_HITTED_CMD_unnamed_type_flag struct {
	Actioncode                       [0]byte              // offset: 0, size: 0 bytes
	Isresist                         uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BAT_BASH_HITTED_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BAT_BASH_HIT_CMD::<unnamed-type-damage> (sizeof = 2 bytes)
type PROTO_NC_BAT_BASH_HIT_CMD_unnamed_type_damage struct {
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Damage                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_HIT_CMD_unnamed_type_damage) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_BASH_HIT_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BAT_BASH_HIT_CMD_unnamed_type_flag struct {
	Actioncode                       [0]byte              // offset: 0, size: 0 bytes
	Isresist                         [0]byte              // offset: 0, size: 0 bytes
	IsImmune                         uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BAT_BASH_HIT_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BAT_BASH_MISSED_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BAT_BASH_MISSED_CMD_unnamed_type_flag struct {
	Actioncode                       [0]byte              // offset: 0, size: 0 bytes
	Ismissed                         [0]byte              // offset: 0, size: 0 bytes
	Isshieldblock                    [0]byte              // offset: 0, size: 0 bytes
	Isresist                         [0]byte              // offset: 0, size: 0 bytes
	IsImmune                         uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BAT_BASH_MISSED_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BAT_BASH_MISS_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BAT_BASH_MISS_CMD_unnamed_type_flag struct {
	Actioncode                       [0]byte              // offset: 0, size: 0 bytes
	Ismissed                         [0]byte              // offset: 0, size: 0 bytes
	Isshieldblock                    [0]byte              // offset: 0, size: 0 bytes
	IsImmune                         uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BAT_BASH_MISS_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage (sizeof = 14 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_SkillDamage struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Flag                             uint16               // offset: 2, size: 2 bytes
	Hpchange                         uint32               // offset: 4, size: 4 bytes
	Resthp                           uint32               // offset: 8, size: 4 bytes
	Hpchangeorder                    uint16               // offset: 12, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_SkillDamage) FixedWireSize() int { return 14 }

// PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage::<unnamed-type-flag> (sizeof = 2 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_SkillDamage_unnamed_type_flag struct {
	Isdamage                         [0]byte              // offset: 0, size: 0 bytes
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Ismissed                         [0]byte              // offset: 0, size: 0 bytes
	Isshieldblock                    [0]byte              // offset: 0, size: 0 bytes
	Isheal                           [0]byte              // offset: 0, size: 0 bytes
	Isenchant                        [0]byte              // offset: 0, size: 0 bytes
	Isresist                         [0]byte              // offset: 0, size: 0 bytes
	IsCostumWeapon                   uint8                // offset: 0, size: 1 bytes
	IsDead                           [0]byte              // offset: 1, size: 0 bytes
	IsImmune                         [0]byte              // offset: 1, size: 0 bytes
	IsCostumShield                   uint8                // offset: 1, size: 1 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_SkillDamage_unnamed_type_flag) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SKILLSMASH_HIT_CMD::<unnamed-type-damage> (sizeof = 2 bytes)
type PROTO_NC_BAT_SKILLSMASH_HIT_CMD_unnamed_type_damage struct {
	Isdead                           [0]byte              // offset: 0, size: 0 bytes
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Damage                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLSMASH_HIT_CMD_unnamed_type_damage) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SMASH_HITTED_ACK::<unnamed-type-hitted> (sizeof = 3 bytes)
type PROTO_NC_BAT_SMASH_HITTED_ACK_unnamed_type_hitted struct {
	Damage                           uint16               // offset: 0, size: 2 bytes
	Attackcode                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_BAT_SMASH_HITTED_ACK_unnamed_type_hitted) FixedWireSize() int { return 3 }

// PROTO_NC_BAT_SMASH_HITTED_ACK::<unnamed-type-hitted>::<unnamed-type-damage> (sizeof = 2 bytes)
type PROTO_NC_BAT_SMASH_HITTED_ACK_unnamed_type_hitted_unnamed_type_damage struct {
	Isdead                           [0]byte              // offset: 0, size: 0 bytes
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Damage                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_SMASH_HITTED_ACK_unnamed_type_hitted_unnamed_type_damage) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SMASH_HIT_ACK::<unnamed-type-damage> (sizeof = 2 bytes)
type PROTO_NC_BAT_SMASH_HIT_ACK_unnamed_type_damage struct {
	Isdead                           [0]byte              // offset: 0, size: 0 bytes
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Damage                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_SMASH_HIT_ACK_unnamed_type_damage) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD_unnamed_type_flag struct {
	IsMissed                         [0]byte              // offset: 0, size: 0 bytes
	IsCostumCharged                  [0]byte              // offset: 0, size: 0 bytes
	IsCostumShieldCharged            uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BAT_SWING_DAMAGE_CMD::<unnamed-type-flag> (sizeof = 2 bytes)
type PROTO_NC_BAT_SWING_DAMAGE_CMD_unnamed_type_flag struct {
	Iscritical                       [0]byte              // offset: 0, size: 0 bytes
	Isresist                         [0]byte              // offset: 0, size: 0 bytes
	Ismissed                         [0]byte              // offset: 0, size: 0 bytes
	Isshieldblock                    [0]byte              // offset: 0, size: 0 bytes
	IsCostumCharged                  [0]byte              // offset: 0, size: 0 bytes
	IsDead                           [0]byte              // offset: 0, size: 0 bytes
	IsDamege2Heal                    [0]byte              // offset: 0, size: 0 bytes
	IsImmune                         uint8                // offset: 0, size: 1 bytes
	IsCostumShieldCharged            uint8                // offset: 1, size: 1 bytes
}

func (s *PROTO_NC_BAT_SWING_DAMAGE_CMD_unnamed_type_flag) FixedWireSize() int { return 2 }

// PROTO_NC_BOOTH_ENTRY_BUY_ACK::BoothItemList (sizeof = 13 bytes)
type PROTO_NC_BOOTH_ENTRY_BUY_ACK_BoothItemList struct {
	Boothslot                        uint8                // offset: 0, size: 1 bytes
	Item                             uint16               // offset: 1, size: 2 bytes
	Unitcost                         uint64               // offset: 3, size: 8 bytes
	Restquantity                     uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_ENTRY_BUY_ACK_BoothItemList) FixedWireSize() int { return 13 }

// PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList (sizeof = 113 bytes)
type PROTO_NC_BOOTH_ENTRY_SELL_ACK_BoothItemList struct {
	Datasize                         uint8                // offset: 0, size: 1 bytes
	Boothslot                        uint8                // offset: 1, size: 1 bytes
	Unitcost                         uint64               // offset: 2, size: 8 bytes
	Item                             [103]byte            // offset: 10, size: 103 bytes
}

func (s *PROTO_NC_BOOTH_ENTRY_SELL_ACK_BoothItemList) FixedWireSize() int { return 113 }

// PROTO_NC_BOOTH_OPEN_REQ::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BOOTH_OPEN_REQ_unnamed_type_flag struct {
	Issell                           [0]byte              // offset: 0, size: 0 bytes
	Itemnum                          uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BOOTH_OPEN_REQ_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BOOTH_OPEN_REQ::BoothItem (sizeof = 12 bytes)
type PROTO_NC_BOOTH_OPEN_REQ_BoothItem struct {
	Invenslot                        uint8                // offset: 0, size: 1 bytes
	Boardslot                        uint8                // offset: 1, size: 1 bytes
	Unitcost                         uint64               // offset: 2, size: 8 bytes
	Totallot                         uint16               // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_BOOTH_OPEN_REQ_BoothItem) FixedWireSize() int { return 12 }

// PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr> (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_DROPEDITEM_CMD_unnamed_type_attr struct {
	Rareness                         [0]byte              // offset: 0, size: 0 bytes
	Looting                          uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_DROPEDITEM_CMD_unnamed_type_attr) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD_unnamed_type_flag struct {
	Loop                             uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD::<unnamed-type-shapedata> (sizeof = 45 bytes)
type PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD_unnamed_type_shapedata struct {
	Notcamp                          [0]byte              // offset: 0, size: 0 bytes
	Camp                             [0]byte              // offset: 0, size: 0 bytes
	Booth                            [0]byte              // offset: 0, size: 0 bytes
	Ride                             [45]byte             // offset: 0, size: 45 bytes
}

func (s *PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD_unnamed_type_shapedata) FixedWireSize() int { return 45 }

// PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag> (sizeof = 99 bytes)
type PROTO_NC_BRIEFINFO_REGENMOB_CMD_unnamed_type_flag struct {
	Abstatebit                       [0]byte              // offset: 0, size: 0 bytes
	Gate2where                       [99]byte             // offset: 0, size: 99 bytes
}

func (s *PROTO_NC_BRIEFINFO_REGENMOB_CMD_unnamed_type_flag) FixedWireSize() int { return 99 }

// PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD::<unnamed-type-PacketOrder> (sizeof = 1 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD_unnamed_type_PacketOrder struct {
	NPacketOrderNum                  [0]byte              // offset: 0, size: 0 bytes
	NIsLastOfInven                   uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD_unnamed_type_PacketOrder) FixedWireSize() int { return 1 }

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD::tagQUESTDOING (sizeof = 6 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD_tagQUESTDOING struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	BNeedClear                       uint8                // offset: 4, size: 1 bytes
	NNumOfDoingQuest                 uint8                // offset: 5, size: 1 bytes
	QuestDoingArray                  [0]byte              // offset: 6, flexible array member
}

func (s *PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD_tagQUESTDOING) FixedWireSize() int { return 6 }

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD::tagQUESTDONE (sizeof = 12 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD_tagQUESTDONE struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NTotalDoneQuest                  uint16               // offset: 4, size: 2 bytes
	NTotalDoneQuestSize              uint16               // offset: 6, size: 2 bytes
	NDoneQuestCount                  uint16               // offset: 8, size: 2 bytes
	NIndex                           uint16               // offset: 10, size: 2 bytes
	QuestDoneArray                   [0]byte              // offset: 12, flexible array member
}

func (s *PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD_tagQUESTDONE) FixedWireSize() int { return 12 }

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD::tagQUESTREAD (sizeof = 6 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD_tagQUESTREAD struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NNumOfReadQuest                  uint16               // offset: 4, size: 2 bytes
	QuestReadIDArray                 [0]byte              // offset: 6, flexible array member
}

func (s *PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD_tagQUESTREAD) FixedWireSize() int { return 6 }

// PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD::tagQUESTREPEAT (sizeof = 6 bytes)
type PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD_tagQUESTREPEAT struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NNumOfRepeatQuest                uint16               // offset: 4, size: 2 bytes
	QuestRepeatArray                 [0]byte              // offset: 6, flexible array member
}

func (s *PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD_tagQUESTREPEAT) FixedWireSize() int { return 6 }

// PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags> (sizeof = 4 bytes)
type PROTO_NC_CHAR_BASE_CMD_unnamed_type_flags struct {
	Bin                              [0]byte              // offset: 0, size: 0 bytes
	Str                              uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_CHAR_BASE_CMD_unnamed_type_flags) FixedWireSize() int { return 4 }

// PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str> (sizeof = 4 bytes)
type PROTO_NC_CHAR_BASE_CMD_unnamed_type_flags_unnamed_type_str struct {
	Skillempower_can_reset           [0]byte              // offset: 0, size: 0 bytes
	Restunion                        uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_CHAR_BASE_CMD_unnamed_type_flags_unnamed_type_str) FixedWireSize() int { return 4 }

// PROTO_NC_CHAR_BASE_CMD::LoginLocation (sizeof = 21 bytes)
type PROTO_NC_CHAR_BASE_CMD_LoginLocation struct {
	Currentmap                       Name3                // offset: 0, size: 12 bytes
	Currentcoord                     [9]byte              // offset: 12, size: 9 bytes
}

func (s *PROTO_NC_CHAR_BASE_CMD_LoginLocation) FixedWireSize() int { return 21 }

// PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo (sizeof = 4 bytes)
type PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD_CardBookmarkInfo struct {
	Slot                             uint16               // offset: 0, size: 2 bytes
	CardID                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD_CardBookmarkInfo) FixedWireSize() int { return 4 }

// PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform (sizeof = 11 bytes)
type PROTO_NC_CHAR_CARDCOLLECT_CMD_CardInform struct {
	RegistTime                       int32                // offset: 0, size: 4 bytes
	SerialNumber                     uint32               // offset: 4, size: 4 bytes
	CardID                           uint16               // offset: 8, size: 2 bytes
	StarNumber                       uint8                // offset: 10, size: 1 bytes
}

func (s *PROTO_NC_CHAR_CARDCOLLECT_CMD_CardInform) FixedWireSize() int { return 11 }

// PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag> (sizeof = 1 bytes)
type PROTO_NC_CHAR_CLIENT_ITEM_CMD_unnamed_type_flag struct {
	Invenclear                       uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_CHAR_CLIENT_ITEM_CMD_unnamed_type_flag) FixedWireSize() int { return 1 }

// PROTO_NC_CHAR_FREESTAT_SET_DB_REQ::<unnamed-type-DecItemInfo> (sizeof = 12 bytes)
type PROTO_NC_CHAR_FREESTAT_SET_DB_REQ_unnamed_type_DecItemInfo struct {
	NIsDeleteItem                    uint8                // offset: 0, size: 1 bytes
	NItmeKey                         uint64               // offset: 1, size: 8 bytes
	NItemID                          uint16               // offset: 9, size: 2 bytes
	NItemLot                         uint8                // offset: 11, size: 1 bytes
}

func (s *PROTO_NC_CHAR_FREESTAT_SET_DB_REQ_unnamed_type_DecItemInfo) FixedWireSize() int { return 12 }

// PROTO_NC_CHAR_ITEM_CMD::<unnamed-type-PacketOrder> (sizeof = 1 bytes)
type PROTO_NC_CHAR_ITEM_CMD_unnamed_type_PacketOrder struct {
	NPacketOrderNum                  [0]byte              // offset: 0, size: 0 bytes
	NIsLastOfInven                   uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_CHAR_ITEM_CMD_unnamed_type_PacketOrder) FixedWireSize() int { return 1 }

// PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner> (sizeof = 24 bytes)
type PROTO_NC_CHAR_WEDDINGDATA_ACK_unnamed_type_Partner struct {
	Name                             Name5                // offset: 0, size: 20 bytes
	RegisterNumber                   uint32               // offset: 20, size: 4 bytes
}

func (s *PROTO_NC_CHAR_WEDDINGDATA_ACK_unnamed_type_Partner) FixedWireSize() int { return 24 }

// PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status> (sizeof = 1 bytes)
type PROTO_NC_CHAR_WEDDINGDATA_ACK_unnamed_type_Status struct {
	Main                             [0]byte              // offset: 0, size: 0 bytes
	Sub                              uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_CHAR_WEDDINGDATA_ACK_unnamed_type_Status) FixedWireSize() int { return 1 }

// PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo> (sizeof = 8 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK_unnamed_type_PartnerInfo struct {
	Level                            uint8                // offset: 0, size: 1 bytes
	Connected                        uint8                // offset: 1, size: 1 bytes
	LastConnect                      uint32               // offset: 2, size: 4 bytes
	CurrentMap                       uint16               // offset: 6, size: 2 bytes
}

func (s *PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK_unnamed_type_PartnerInfo) FixedWireSize() int { return 8 }

// PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG::<unnamed-type-PartnerInfo> (sizeof = 8 bytes)
type PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG_unnamed_type_PartnerInfo struct {
	Level                            uint8                // offset: 0, size: 1 bytes
	Connected                        uint8                // offset: 1, size: 1 bytes
	CurrentMap                       uint16               // offset: 2, size: 2 bytes
	RegisterNumber                   uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG_unnamed_type_PartnerInfo) FixedWireSize() int { return 8 }

// PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ::<unnamed-type-MachineIndex> (sizeof = 3 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ_unnamed_type_MachineIndex struct {
	Machine                          uint8                // offset: 0, size: 1 bytes
	BetUnits                         uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ_unnamed_type_MachineIndex) FixedWireSize() int { return 3 }

// PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD::<unnamed-type-JackPotList> (sizeof = 6 bytes)
type PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD_unnamed_type_JackPotList struct {
	BettingUnit                      uint16               // offset: 0, size: 2 bytes
	JackpotPrice                     uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD_unnamed_type_JackPotList) FixedWireSize() int { return 6 }

// PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD::Memberinfo (sizeof = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD_Memberinfo struct {
	Loc                              SHINE_XY_TYPE        // offset: 0, size: 8 bytes
}

func (s *PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD_Memberinfo) FixedWireSize() int { return 8 }

// PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-produce> (sizeof = 43 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_produce struct {
	Create                           uint8                // offset: 0, size: 1 bytes
	Product                          uint16               // offset: 1, size: 2 bytes
	Productuselot_log                uint16               // offset: 3, size: 2 bytes
	Iteminfo                         [38]byte             // offset: 5, size: 38 bytes
}

func (s *PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_produce) FixedWireSize() int { return 43 }

// PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-produce>::<unnamed-type-iteminfo> (sizeof = 38 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_produce_unnamed_type_iteminfo struct {
	Itemcreate                       [0]byte              // offset: 0, size: 0 bytes
	Itemchange                       [38]byte             // offset: 0, size: 38 bytes
}

func (s *PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_produce_unnamed_type_iteminfo) FixedWireSize() int { return 38 }

// PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-raw> (sizeof = 17 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_raw struct {
	Vanish                           uint8                // offset: 0, size: 1 bytes
	Rawid                            uint16               // offset: 1, size: 2 bytes
	Rawuselot_log                    uint16               // offset: 3, size: 2 bytes
	Iteminfo                         [12]byte             // offset: 5, size: 12 bytes
}

func (s *PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_raw) FixedWireSize() int { return 17 }

// PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-raw>::<unnamed-type-iteminfo> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_raw_unnamed_type_iteminfo struct {
	Itemvanish                       [0]byte              // offset: 0, size: 0 bytes
	Itemchange                       [12]byte             // offset: 0, size: 12 bytes
}

func (s *PROTO_NC_ITEMDB_DISMANTLE_REQ_unnamed_type_raw_unnamed_type_iteminfo) FixedWireSize() int { return 12 }

// PROTO_NC_ITEMDB_DROPLOT_REQ::<unnamed-type-dropitem> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_DROPLOT_REQ_unnamed_type_dropitem struct {
	Itemregnum                       uint64               // offset: 0, size: 8 bytes
	Lot                              uint32               // offset: 8, size: 4 bytes
}

func (s *PROTO_NC_ITEMDB_DROPLOT_REQ_unnamed_type_dropitem) FixedWireSize() int { return 12 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK::<unnamed-type-equipment> (sizeof = 6 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK_unnamed_type_equipment struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
	EnchantedSocketIndex             uint8                // offset: 2, size: 1 bytes
	GemItemIndex                     uint16               // offset: 3, size: 2 bytes
	RestGemCount                     uint8                // offset: 5, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK_unnamed_type_equipment) FixedWireSize() int { return 6 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK::<unnamed-type-gem> (sizeof = 2 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK_unnamed_type_gem struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK_unnamed_type_gem) FixedWireSize() int { return 2 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ::<unnamed-type-equipment> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ_unnamed_type_equipment struct {
	ItemRegistNumber                 uint64               // offset: 0, size: 8 bytes
	ItemID                           uint16               // offset: 8, size: 2 bytes
	ItemSlot                         ITEM_INVEN           // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ_unnamed_type_equipment) FixedWireSize() int { return 12 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK::<unnamed-type-drill> (sizeof = 10 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK_unnamed_type_drill struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
	ItemLot                          uint64               // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK_unnamed_type_drill) FixedWireSize() int { return 10 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK::<unnamed-type-equipment> (sizeof = 3 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK_unnamed_type_equipment struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
	AllSocketCount                   uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK_unnamed_type_equipment) FixedWireSize() int { return 3 }

// PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ::<unnamed-type-equipment> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ_unnamed_type_equipment struct {
	ItemRegistNumber                 uint64               // offset: 0, size: 8 bytes
	ItemID                           uint16               // offset: 8, size: 2 bytes
	ItemSlot                         ITEM_INVEN           // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ_unnamed_type_equipment) FixedWireSize() int { return 12 }

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-equipment> (sizeof = 2 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_equipment struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_equipment) FixedWireSize() int { return 2 }

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-gemSlot> (sizeof = 1 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_gemSlot struct {
	Deleted                          uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_gemSlot) FixedWireSize() int { return 1 }

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-remover> (sizeof = 10 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_remover struct {
	ItemSlot                         ITEM_INVEN           // offset: 0, size: 2 bytes
	ItemLot                          uint64               // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK_unnamed_type_remover) FixedWireSize() int { return 10 }

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ::<unnamed-type-equipment> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ_unnamed_type_equipment struct {
	ItemRegistNumber                 uint64               // offset: 0, size: 8 bytes
	ItemID                           uint16               // offset: 8, size: 2 bytes
	ItemSlot                         ITEM_INVEN           // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ_unnamed_type_equipment) FixedWireSize() int { return 12 }

// PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ::<unnamed-type-gemSlot> (sizeof = 1 bytes)
type PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ_unnamed_type_gemSlot struct {
	BeDeleted                        uint8                // offset: 0, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ_unnamed_type_gemSlot) FixedWireSize() int { return 1 }

// PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ::<unnamed-type-sockets> (sizeof = 3 bytes)
type PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ_unnamed_type_sockets struct {
	GemItemIndex                     uint16               // offset: 0, size: 2 bytes
	RestGemCount                     uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ_unnamed_type_sockets) FixedWireSize() int { return 3 }

// PROTO_NC_ITEMDB_EQUIP_REQ::<unnamed-type-unequip> (sizeof = 9 bytes)
type PROTO_NC_ITEMDB_EQUIP_REQ_unnamed_type_unequip struct {
	Item                             uint64               // offset: 0, size: 8 bytes
	Invenslot                        uint8                // offset: 8, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_EQUIP_REQ_unnamed_type_unequip) FixedWireSize() int { return 9 }

// PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ::CREATE_ITEM_INFO (sizeof = 113 bytes)
type PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ_CREATE_ITEM_INFO struct {
	NItemKey                         uint64               // offset: 0, size: 8 bytes
	NItemID                          uint16               // offset: 8, size: 2 bytes
	ItemLoc                          ITEM_INVEN           // offset: 10, size: 2 bytes
	ItemAttr                         [101]byte            // offset: 12, size: 101 bytes
}

func (s *PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ_CREATE_ITEM_INFO) FixedWireSize() int { return 113 }

// PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-cen> (sizeof = 25 bytes)
type PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_cen struct {
	From                             uint8                // offset: 0, size: 1 bytes
	Cen                              uint64               // offset: 1, size: 8 bytes
	Result_money_from                uint64               // offset: 9, size: 8 bytes
	Result_money_to                  uint64               // offset: 17, size: 8 bytes
}

func (s *PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_cen) FixedWireSize() int { return 25 }

// PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-tradeitem> (sizeof = 26 bytes)
type PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_tradeitem struct {
	InvenA                           [13]byte             // offset: 0, size: 13 bytes
	InvenB                           [13]byte             // offset: 13, size: 13 bytes
}

func (s *PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_tradeitem) FixedWireSize() int { return 26 }

// PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-tradeitem>::<unnamed-type-invenA> (sizeof = 13 bytes)
type PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_tradeitem_unnamed_type_invenA struct {
	Regnum                           uint64               // offset: 0, size: 8 bytes
	Id                               uint16               // offset: 8, size: 2 bytes
	Lot                              uint16               // offset: 10, size: 2 bytes
	Invenloc                         uint8                // offset: 12, size: 1 bytes
}

func (s *PROTO_NC_ITEMDB_ITEMTRADE_REQ_unnamed_type_tradeitem_unnamed_type_invenA) FixedWireSize() int { return 13 }

// PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-produce> (sizeof = 39 bytes)
type PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_produce struct {
	Create                           uint8                // offset: 0, size: 1 bytes
	Iteminfo                         [38]byte             // offset: 1, size: 38 bytes
}

func (s *PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_produce) FixedWireSize() int { return 39 }

// PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-produce>::<unnamed-type-iteminfo> (sizeof = 38 bytes)
type PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_produce_unnamed_type_iteminfo struct {
	Itemcreate                       [0]byte              // offset: 0, size: 0 bytes
	Itemchange                       [38]byte             // offset: 0, size: 38 bytes
}

func (s *PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_produce_unnamed_type_iteminfo) FixedWireSize() int { return 38 }

// PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-raw> (sizeof = 17 bytes)
type PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_raw struct {
	Vanish                           uint8                // offset: 0, size: 1 bytes
	Rawid                            uint16               // offset: 1, size: 2 bytes
	Rawuselot_log                    uint16               // offset: 3, size: 2 bytes
	Iteminfo                         [12]byte             // offset: 5, size: 12 bytes
}

func (s *PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_raw) FixedWireSize() int { return 17 }

// PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-raw>::<unnamed-type-iteminfo> (sizeof = 12 bytes)
type PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_raw_unnamed_type_iteminfo struct {
	Itemvanish                       [0]byte              // offset: 0, size: 0 bytes
	Itemchange                       [12]byte             // offset: 0, size: 12 bytes
}

func (s *PROTO_NC_ITEMDB_PRODUCE_REQ_unnamed_type_raw_unnamed_type_iteminfo) FixedWireSize() int { return 12 }

// PROTO_NC_KQ_SCORE_CMD::KQScore (sizeof = 38 bytes)
type PROTO_NC_KQ_SCORE_CMD_KQScore struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Flag                             uint16               // offset: 20, size: 2 bytes
	Killpoint                        uint32               // offset: 22, size: 4 bytes
	Downpoint                        uint32               // offset: 26, size: 4 bytes
	Contribute                       uint32               // offset: 30, size: 4 bytes
	Score                            uint32               // offset: 34, size: 4 bytes
}

func (s *PROTO_NC_KQ_SCORE_CMD_KQScore) FixedWireSize() int { return 38 }

// PROTO_NC_KQ_SCORE_CMD::KQScore::<unnamed-type-flag> (sizeof = 2 bytes)
type PROTO_NC_KQ_SCORE_CMD_KQScore_unnamed_type_flag struct {
	Medal                            [0]byte              // offset: 0, size: 0 bytes
	Gender                           uint8                // offset: 0, size: 1 bytes
	Chrcls                           uint8                // offset: 1, size: 1 bytes
}

func (s *PROTO_NC_KQ_SCORE_CMD_KQScore_unnamed_type_flag) FixedWireSize() int { return 2 }

// PROTO_NC_KQ_SCORE_SIMPLE_CMD::<unnamed-type-self> (sizeof = 13 bytes)
type PROTO_NC_KQ_SCORE_SIMPLE_CMD_unnamed_type_self struct {
	Rank                             uint8                // offset: 0, size: 1 bytes
	Kill                             uint32               // offset: 1, size: 4 bytes
	Down                             uint32               // offset: 5, size: 4 bytes
	Score                            uint32               // offset: 9, size: 4 bytes
}

func (s *PROTO_NC_KQ_SCORE_SIMPLE_CMD_unnamed_type_self) FixedWireSize() int { return 13 }

// PROTO_NC_KQ_SCORE_SIMPLE_CMD::KQScore (sizeof = 32 bytes)
type PROTO_NC_KQ_SCORE_SIMPLE_CMD_KQScore struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Kill                             uint32               // offset: 20, size: 4 bytes
	Down                             uint32               // offset: 24, size: 4 bytes
	Score                            uint32               // offset: 28, size: 4 bytes
}

func (s *PROTO_NC_KQ_SCORE_SIMPLE_CMD_KQScore) FixedWireSize() int { return 32 }

// PROTO_NC_MAP_LINKRESERVE_REQ::<unnamed-type-linkto> (sizeof = 20 bytes)
type PROTO_NC_MAP_LINKRESERVE_REQ_unnamed_type_linkto struct {
	Mapname                          Name3                // offset: 0, size: 12 bytes
	Location                         SHINE_XY_TYPE        // offset: 12, size: 8 bytes
}

func (s *PROTO_NC_MAP_LINKRESERVE_REQ_unnamed_type_linkto) FixedWireSize() int { return 20 }

// PROTO_NC_MAP_REGIST_CMD::<unnamed-type-maparray> (sizeof = 14 bytes)
type PROTO_NC_MAP_REGIST_CMD_unnamed_type_maparray struct {
	Canparty                         uint8                // offset: 0, size: 1 bytes
	Mapname                          Name3                // offset: 1, size: 12 bytes
	IsPartyReleaseMap                uint8                // offset: 13, size: 1 bytes
}

func (s *PROTO_NC_MAP_REGIST_CMD_unnamed_type_maparray) FixedWireSize() int { return 14 }

// PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD::SOULSTONEMENU (sizeof = 12 bytes)
type PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD_SOULSTONEMENU struct {
	Recover                          uint32               // offset: 0, size: 4 bytes
	Maxnum                           uint32               // offset: 4, size: 4 bytes
	Price                            uint32               // offset: 8, size: 4 bytes
}

func (s *PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD_SOULSTONEMENU) FixedWireSize() int { return 12 }

// PROTO_NC_MENU_SHOPOPENTABLE_CMD::MENUITEM (sizeof = 3 bytes)
type PROTO_NC_MENU_SHOPOPENTABLE_CMD_MENUITEM struct {
	Slot                             uint8                // offset: 0, size: 1 bytes
	Item                             uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_MENU_SHOPOPENTABLE_CMD_MENUITEM) FixedWireSize() int { return 3 }

// PROTO_NC_MENU_SHOPOPEN_CMD::MENUITEM (sizeof = 2 bytes)
type PROTO_NC_MENU_SHOPOPEN_CMD_MENUITEM struct {
	Item                             uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_MENU_SHOPOPEN_CMD_MENUITEM) FixedWireSize() int { return 2 }

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD::<unnamed-type-rearrange> (sizeof = 21 bytes)
type PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD_unnamed_type_rearrange struct {
	Furniturehandle                  uint16               // offset: 0, size: 2 bytes
	Furnitureid                      uint16               // offset: 2, size: 2 bytes
	Location                         [16]byte             // offset: 4, size: 16 bytes
	NEndureGrade                     uint8                // offset: 20, size: 1 bytes
}

func (s *PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD_unnamed_type_rearrange) FixedWireSize() int { return 21 }

// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ::<unnamed-type-rearrange> (sizeof = 18 bytes)
type PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ_unnamed_type_rearrange struct {
	Furniturehandle                  uint16               // offset: 0, size: 2 bytes
	Location                         [16]byte             // offset: 2, size: 16 bytes
}

func (s *PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ_unnamed_type_rearrange) FixedWireSize() int { return 18 }

// PROTO_NC_MINIHOUSE_FUNICHER_CMD::<unnamed-type-array> (sizeof = 21 bytes)
type PROTO_NC_MINIHOUSE_FUNICHER_CMD_unnamed_type_array struct {
	Furniturehandle                  uint16               // offset: 0, size: 2 bytes
	Furnitureid                      uint16               // offset: 2, size: 2 bytes
	Location                         [16]byte             // offset: 4, size: 16 bytes
	NEndureGrade                     uint8                // offset: 20, size: 1 bytes
}

func (s *PROTO_NC_MINIHOUSE_FUNICHER_CMD_unnamed_type_array) FixedWireSize() int { return 21 }

// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK::<unnamed-type-player> (sizeof = 3 bytes)
type PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK_unnamed_type_player struct {
	Level                            uint8                // offset: 0, size: 1 bytes
	Hnd                              uint16               // offset: 1, size: 2 bytes
}

func (s *PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK_unnamed_type_player) FixedWireSize() int { return 3 }

// PROTO_NC_MINIHOUSE_REARRANGE_REQ::<unnamed-type-rearrange> (sizeof = 18 bytes)
type PROTO_NC_MINIHOUSE_REARRANGE_REQ_unnamed_type_rearrange struct {
	Furniturehandle                  uint16               // offset: 0, size: 2 bytes
	Location                         [16]byte             // offset: 2, size: 16 bytes
}

func (s *PROTO_NC_MINIHOUSE_REARRANGE_REQ_unnamed_type_rearrange) FixedWireSize() int { return 18 }

// PROTO_NC_PARTY_MEMBERLOCATION_CMD::<unnamed-type-member> (sizeof = 28 bytes)
type PROTO_NC_PARTY_MEMBERLOCATION_CMD_unnamed_type_member struct {
	Charid                           Name5                // offset: 0, size: 20 bytes
	Loc                              SHINE_XY_TYPE        // offset: 20, size: 8 bytes
}

func (s *PROTO_NC_PARTY_MEMBERLOCATION_CMD_unnamed_type_member) FixedWireSize() int { return 28 }

// PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest (sizeof = 4 bytes)
type PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD_MobOfQuest struct {
	NMobID                           uint16               // offset: 0, size: 2 bytes
	NQuestID                         uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD_MobOfQuest) FixedWireSize() int { return 4 }

// PROTO_NC_RAID_LOCATION_CMD::Memberinfo (sizeof = 12 bytes)
type PROTO_NC_RAID_LOCATION_CMD_Memberinfo struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	Loc                              SHINE_XY_TYPE        // offset: 4, size: 8 bytes
}

func (s *PROTO_NC_RAID_LOCATION_CMD_Memberinfo) FixedWireSize() int { return 12 }

// PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD::<unnamed-type-PostPosition> (sizeof = 24 bytes)
type PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_unnamed_type_PostPosition struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	NPartySeqCharNo                  [5]uint32            // offset: 4, size: 20 bytes
}

func (s *PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_unnamed_type_PostPosition) FixedWireSize() int { return 24 }

// PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD::<unnamed-type-PrePosition> (sizeof = 24 bytes)
type PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_unnamed_type_PrePosition struct {
	NPartyNo                         uint16               // offset: 0, size: 2 bytes
	NPartyIndex                      uint8                // offset: 2, size: 1 bytes
	NSlotIndex                       uint8                // offset: 3, size: 1 bytes
	NPartySeqCharNo                  [5]uint32            // offset: 4, size: 20 bytes
}

func (s *PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD_unnamed_type_PrePosition) FixedWireSize() int { return 24 }

// PROTO_NC_SCENARIO_MAPMARK_CMD::IconInfo (sizeof = 15 bytes)
type PROTO_NC_SCENARIO_MAPMARK_CMD_IconInfo struct {
	NGroup                           uint16               // offset: 0, size: 2 bytes
	NCoord                           SHINE_XY_TYPE        // offset: 2, size: 8 bytes
	EMapIcon                         uint8                // offset: 10, size: 1 bytes
	NKeepTime                        uint32               // offset: 11, size: 4 bytes
}

func (s *PROTO_NC_SCENARIO_MAPMARK_CMD_IconInfo) FixedWireSize() int { return 15 }

// PROTO_NC_WT_LICENSE_CLR_ACK::echo_data (sizeof = 4 bytes)
type PROTO_NC_WT_LICENSE_CLR_ACK_echo_data struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_CLR_ACK_echo_data) FixedWireSize() int { return 4 }

// PROTO_NC_WT_LICENSE_CLR_REQ::echo_data (sizeof = 4 bytes)
type PROTO_NC_WT_LICENSE_CLR_REQ_echo_data struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_CLR_REQ_echo_data) FixedWireSize() int { return 4 }

// PROTO_NC_WT_LICENSE_SET_ACK::echo_data (sizeof = 4 bytes)
type PROTO_NC_WT_LICENSE_SET_ACK_echo_data struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_SET_ACK_echo_data) FixedWireSize() int { return 4 }

// PROTO_NC_WT_LICENSE_SET_REQ::echo_data (sizeof = 4 bytes)
type PROTO_NC_WT_LICENSE_SET_REQ_echo_data struct {
	Licenseitem                      ITEM_INVEN           // offset: 0, size: 2 bytes
	Weaponitem                       ITEM_INVEN           // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_WT_LICENSE_SET_REQ_echo_data) FixedWireSize() int { return 4 }
