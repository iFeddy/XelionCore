package protocol

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — zone_structs_bat.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_BAT_ABSTATEINFORM_CMD (sizeof = 8 bytes)
type PROTO_NC_BAT_ABSTATEINFORM_CMD_DATA struct {
	Abstate                          uint32               // offset: 0, size: 4 bytes
	Keeptime_millisec                uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_BAT_ABSTATEINFORM_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_BAT_ABSTATERESET_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_ABSTATERESET_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Abstate                          uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_BAT_ABSTATERESET_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_ABSTATESET_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_ABSTATESET_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Abstate                          uint32               // offset: 2, size: 4 bytes
}

func (s *PROTO_NC_BAT_ABSTATESET_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_ABSTATE_ERASE_ACK (sizeof = 2 bytes)
type PROTO_NC_BAT_ABSTATE_ERASE_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_ABSTATE_ERASE_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_ABSTATE_ERASE_REQ (sizeof = 2 bytes)
type PROTO_NC_BAT_ABSTATE_ERASE_REQ_DATA struct {
	AbstateIndex                     uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_ABSTATE_ERASE_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_AREADOTDAMAGE_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_AREADOTDAMAGE_CMD_DATA struct {
	Abstate                          uint16               // offset: 0, size: 2 bytes
	Centerobject                     uint16               // offset: 2, size: 2 bytes
	Num                              uint8                // offset: 4, size: 1 bytes
	Target                           [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_BAT_AREADOTDAMAGE_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_ASSIST_ACK (sizeof = 2 bytes)
type PROTO_NC_BAT_ASSIST_ACK_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_ASSIST_ACK_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_ASSIST_REQ (sizeof = 2 bytes)
type PROTO_NC_BAT_ASSIST_REQ_DATA struct {
	Tanker                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_ASSIST_REQ_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_BASH_HITTED_CMD (sizeof = 13 bytes)
type PROTO_NC_BAT_BASH_HITTED_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Damage                           uint16               // offset: 2, size: 2 bytes
	Flag                             uint8                // offset: 4, size: 1 bytes
	Resthp                           uint32               // offset: 5, size: 4 bytes
	Attackspeed                      uint16               // offset: 9, size: 2 bytes
	Hpchangeorder                    uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_HITTED_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_BAT_BASH_HIT_CMD (sizeof = 13 bytes)
type PROTO_NC_BAT_BASH_HIT_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Damage                           uint16               // offset: 2, size: 2 bytes
	Flag                             uint8                // offset: 4, size: 1 bytes
	Resthp                           uint32               // offset: 5, size: 4 bytes
	Attackspeed                      uint16               // offset: 9, size: 2 bytes
	Hpchangeorder                    uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_HIT_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_BAT_BASH_MISSED_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_BASH_MISSED_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Flag                             uint8                // offset: 2, size: 1 bytes
	Attackspeed                      uint16               // offset: 3, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_MISSED_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_BASH_MISS_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_BASH_MISS_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Flag                             uint8                // offset: 2, size: 1 bytes
	Attackspeed                      uint16               // offset: 3, size: 2 bytes
}

func (s *PROTO_NC_BAT_BASH_MISS_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_CEASE_FIRE_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_CEASE_FIRE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_CEASE_FIRE_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_DOTDAMAGE_CMD (sizeof = 13 bytes)
type PROTO_NC_BAT_DOTDAMAGE_CMD_DATA struct {
	Object                           uint16               // offset: 0, size: 2 bytes
	Resthp                           uint32               // offset: 2, size: 4 bytes
	Damage                           uint16               // offset: 6, size: 2 bytes
	Abstate                          uint16               // offset: 8, size: 2 bytes
	Hpchangeorder                    uint16               // offset: 10, size: 2 bytes
	IsMissDamage                     uint8                // offset: 12, size: 1 bytes
}

func (s *PROTO_NC_BAT_DOTDAMAGE_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_BAT_EXPGAIN_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_EXPGAIN_CMD_DATA struct {
	Expgain                          uint32               // offset: 0, size: 4 bytes
	Mobhandle                        uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_EXPGAIN_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_EXPLOST_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_EXPLOST_CMD_DATA struct {
	Explost                          uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_BAT_EXPLOST_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_FAMEGAIN_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_FAMEGAIN_CMD_DATA struct {
	Famegain                         uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_BAT_FAMEGAIN_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_HPCHANGE_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_HPCHANGE_CMD_DATA struct {
	Hp                               uint32               // offset: 0, size: 4 bytes
	Hpchangeorder                    uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_HPCHANGE_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_LEVELUP_CMD (sizeof = 235 bytes)
type PROTO_NC_BAT_LEVELUP_CMD_DATA struct {
	Level                            uint8                // offset: 0, size: 1 bytes
	Mobhandle                        uint16               // offset: 1, size: 2 bytes
	Newparam                         [232]byte            // offset: 3, size: 232 bytes
}

func (s *PROTO_NC_BAT_LEVELUP_CMD_DATA) FixedWireSize() int { return 235 }

// PROTO_NC_BAT_LPCHANGE_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_LPCHANGE_CMD_DATA struct {
	NLP                              uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_BAT_LPCHANGE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_MOBSLAYER_CMD (sizeof = 3 bytes)
type PROTO_NC_BAT_MOBSLAYER_CMD_DATA struct {
	Killedmob                        uint16               // offset: 0, size: 2 bytes
	Number                           uint8                // offset: 2, size: 1 bytes
	Slayers                          [0]byte              // offset: 3, flexible array member
}

func (s *PROTO_NC_BAT_MOBSLAYER_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BAT_PKINPKFIELD_CLIENT_CMD (sizeof = 8 bytes)
type PROTO_NC_BAT_PKINPKFIELD_CLIENT_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Defender                         uint16               // offset: 2, size: 2 bytes
	Pkcount                          uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_BAT_PKINPKFIELD_CLIENT_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_BAT_PKINPKFIELD_WMS_CMD (sizeof = 10 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_BAT_PKINPKFIELD_WMS_CMD_DATA)

// PROTO_NC_BAT_REALLYKILL_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_REALLYKILL_CMD_DATA struct {
	Dead                             uint16               // offset: 0, size: 2 bytes
	Attacker                         uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_BAT_REALLYKILL_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_REFLECTIONDAMAGE_CMD (sizeof = 14 bytes)
type PROTO_NC_BAT_REFLECTIONDAMAGE_CMD_DATA struct {
	Abstate                          uint16               // offset: 0, size: 2 bytes
	ReflectionObject                 uint16               // offset: 2, size: 2 bytes
	TargetObject                     uint16               // offset: 4, size: 2 bytes
	Resthp                           uint32               // offset: 6, size: 4 bytes
	Damage                           uint16               // offset: 10, size: 2 bytes
	Hpchangeorder                    uint16               // offset: 12, size: 2 bytes
}

func (s *PROTO_NC_BAT_REFLECTIONDAMAGE_CMD_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ (sizeof = 10 bytes)
type PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Locate                           SHINE_XY_TYPE        // offset: 2, size: 8 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ_DATA) FixedWireSize() int { return 10 }

// PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD_DATA struct {
	Index                            uint16               // offset: 0, size: 2 bytes
	Caster                           uint16               // offset: 2, size: 2 bytes
	NSkillID                         uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_DATA struct {
	Index                            uint16               // offset: 0, size: 2 bytes
	Caster                           uint16               // offset: 2, size: 2 bytes
	Targetnum                        uint8                // offset: 4, size: 1 bytes
	Target                           [0]byte              // offset: 5, flexible array member
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD (sizeof = 12 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Targetloc                        SHINE_XY_TYPE        // offset: 2, size: 8 bytes
	Index                            uint16               // offset: 10, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Targetobj                        uint16               // offset: 2, size: 2 bytes
	Index                            uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ (sizeof = 4 bytes)
type PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_SKILLBLAST_LIGHTNINGWAVE_CMD (sizeof = 8 bytes)
type PROTO_NC_BAT_SKILLBLAST_LIGHTNINGWAVE_CMD_DATA struct {
	NBIndex                          uint16               // offset: 0, size: 2 bytes
	NSkillID                         uint16               // offset: 2, size: 2 bytes
	NShotHnd                         uint16               // offset: 4, size: 2 bytes
	NTargetHnd                       uint16               // offset: 6, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLBLAST_LIGHTNINGWAVE_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_BAT_SKILLENCHANT_REQ (sizeof = 4 bytes)
type PROTO_NC_BAT_SKILLENCHANT_REQ_DATA struct {
	Skillid                          uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_BAT_SKILLENCHANT_REQ_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_SKILLSMASH_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_SKILLSMASH_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Attackcode                       uint8                // offset: 4, size: 1 bytes
}

func (s *PROTO_NC_BAT_SKILLSMASH_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_SKILLSMASH_HITTED_CMD (sizeof = 9 bytes)
type PROTO_NC_BAT_SKILLSMASH_HITTED_CMD_DATA struct {
	Attack                           uint16               // offset: 0, size: 2 bytes
	Hitted                           [7]byte              // offset: 2, size: 7 bytes
}

func (s *PROTO_NC_BAT_SKILLSMASH_HITTED_CMD_DATA) FixedWireSize() int { return 9 }

// PROTO_NC_BAT_SKILLSMASH_HIT_CMD (sizeof = 7 bytes)
type PROTO_NC_BAT_SKILLSMASH_HIT_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Skill                            uint16               // offset: 2, size: 2 bytes
	Damage                           uint16               // offset: 4, size: 2 bytes
	Attackcode                       uint8                // offset: 6, size: 1 bytes
}

func (s *PROTO_NC_BAT_SKILLSMASH_HIT_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_BAT_SKILLSMASH_MISS_CMD (sizeof = 3 bytes)
type PROTO_NC_BAT_SKILLSMASH_MISS_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Attackcode                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_BAT_SKILLSMASH_MISS_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BAT_SMASH_HITTED_ACK (sizeof = 5 bytes)
type PROTO_NC_BAT_SMASH_HITTED_ACK_DATA struct {
	Attack                           uint16               // offset: 0, size: 2 bytes
	Hitted                           [3]byte              // offset: 2, size: 3 bytes
}

func (s *PROTO_NC_BAT_SMASH_HITTED_ACK_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_SMASH_HIT_ACK (sizeof = 5 bytes)
type PROTO_NC_BAT_SMASH_HIT_ACK_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Damage                           uint16               // offset: 2, size: 2 bytes
	Attackcode                       uint8                // offset: 4, size: 1 bytes
}

func (s *PROTO_NC_BAT_SMASH_HIT_ACK_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_SMASH_MISS_ACK (sizeof = 3 bytes)
type PROTO_NC_BAT_SMASH_MISS_ACK_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Attackcode                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_BAT_SMASH_MISS_ACK_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BAT_SOMEONEBASH_HIT_CMD (sizeof = 13 bytes)
type PROTO_NC_BAT_SOMEONEBASH_HIT_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Actioncode                       uint8                // offset: 4, size: 1 bytes
	Resthp                           uint32               // offset: 5, size: 4 bytes
	Attackspeed                      uint16               // offset: 9, size: 2 bytes
	Hpchangeorder                    uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_BAT_SOMEONEBASH_HIT_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_BAT_SOMEONEBASH_MISS_CMD (sizeof = 7 bytes)
type PROTO_NC_BAT_SOMEONEBASH_MISS_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Actioncode                       uint8                // offset: 4, size: 1 bytes
	Attackspeed                      uint16               // offset: 5, size: 2 bytes
}

func (s *PROTO_NC_BAT_SOMEONEBASH_MISS_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_BAT_SOMEONEDAMAGED_LARGE_CMD (sizeof = 7 bytes)
type PROTO_NC_BAT_SOMEONEDAMAGED_LARGE_CMD_DATA struct {
	Attack                           uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Damage                           uint16               // offset: 4, size: 2 bytes
	Actionnumber                     uint8                // offset: 6, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOMEONEDAMAGED_LARGE_CMD_DATA) FixedWireSize() int { return 7 }

// PROTO_NC_BAT_SOMEONEDAMAGED_SMALL_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_SOMEONEDAMAGED_SMALL_CMD_DATA struct {
	Attack                           uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Damage                           uint8                // offset: 4, size: 1 bytes
	Actionnumber                     uint8                // offset: 5, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOMEONEDAMAGED_SMALL_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD (sizeof = 12 bytes)
type PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Caster                           uint16               // offset: 2, size: 2 bytes
	Loc                              SHINE_XY_TYPE        // offset: 4, size: 8 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD_DATA) FixedWireSize() int { return 12 }

// PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD (sizeof = 14 bytes)
type PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
	Castinfo                         [12]byte             // offset: 2, size: 12 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD_DATA) FixedWireSize() int { return 14 }

// PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD (sizeof = 8 bytes)
type PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
	Castinfo                         [6]byte              // offset: 2, size: 6 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD_DATA) FixedWireSize() int { return 8 }

// PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Caster                           uint16               // offset: 2, size: 2 bytes
	Target                           uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_SOMEONESKILLSMASH_DEAD_CMD (sizeof = 6 bytes)
type PROTO_NC_BAT_SOMEONESKILLSMASH_DEAD_CMD_DATA struct {
	Skill                            uint16               // offset: 0, size: 2 bytes
	Attack                           uint16               // offset: 2, size: 2 bytes
	Target                           uint16               // offset: 4, size: 2 bytes
}

func (s *PROTO_NC_BAT_SOMEONESKILLSMASH_DEAD_CMD_DATA) FixedWireSize() int { return 6 }

// PROTO_NC_BAT_SOMEONESMASH_DAMAGED_CMD (sizeof = 5 bytes)
type PROTO_NC_BAT_SOMEONESMASH_DAMAGED_CMD_DATA struct {
	Attack                           uint16               // offset: 0, size: 2 bytes
	Target                           uint16               // offset: 2, size: 2 bytes
	Actionnumber                     uint8                // offset: 4, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOMEONESMASH_DAMAGED_CMD_DATA) FixedWireSize() int { return 5 }

// PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD (sizeof = 13 bytes)
type PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Defender                         uint16               // offset: 2, size: 2 bytes
	Flag                             uint8                // offset: 4, size: 1 bytes
	Resthp                           uint32               // offset: 5, size: 4 bytes
	Hpchangeorder                    uint16               // offset: 9, size: 2 bytes
	Damageindex                      uint8                // offset: 11, size: 1 bytes
	Attacksequence                   uint8                // offset: 12, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD_DATA) FixedWireSize() int { return 13 }

// PROTO_NC_BAT_SOULCOLLECT_CMD (sizeof = 3 bytes)
type PROTO_NC_BAT_SOULCOLLECT_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
	Soulnumber                       uint8                // offset: 2, size: 1 bytes
}

func (s *PROTO_NC_BAT_SOULCOLLECT_CMD_DATA) FixedWireSize() int { return 3 }

// PROTO_NC_BAT_SPCHANGE_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_SPCHANGE_CMD_DATA struct {
	Sp                               uint32               // offset: 0, size: 4 bytes
}

func (s *PROTO_NC_BAT_SPCHANGE_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_SUMEONELEVELUP_CMD (sizeof = 4 bytes)
type PROTO_NC_BAT_SUMEONELEVELUP_CMD_DATA struct {
	Ply                              uint16               // offset: 0, size: 2 bytes
	Mobhandle                        uint16               // offset: 2, size: 2 bytes
}

func (s *PROTO_NC_BAT_SUMEONELEVELUP_CMD_DATA) FixedWireSize() int { return 4 }

// PROTO_NC_BAT_SUMEONESKILLCUT_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_SUMEONESKILLCUT_CMD_DATA struct {
	Caster                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_SUMEONESKILLCUT_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_SWING_DAMAGE_CMD (sizeof = 16 bytes)
type PROTO_NC_BAT_SWING_DAMAGE_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Defender                         uint16               // offset: 2, size: 2 bytes
	Flag                             uint16               // offset: 4, size: 2 bytes
	Damage                           uint16               // offset: 6, size: 2 bytes
	Resthp                           uint32               // offset: 8, size: 4 bytes
	Hpchangeorder                    uint16               // offset: 12, size: 2 bytes
	Damageindex                      uint8                // offset: 14, size: 1 bytes
	Attacksequence                   uint8                // offset: 15, size: 1 bytes
}

func (s *PROTO_NC_BAT_SWING_DAMAGE_CMD_DATA) FixedWireSize() int { return 16 }

// PROTO_NC_BAT_SWING_START_CMD (sizeof = 9 bytes)
type PROTO_NC_BAT_SWING_START_CMD_DATA struct {
	Attacker                         uint16               // offset: 0, size: 2 bytes
	Defender                         uint16               // offset: 2, size: 2 bytes
	Actioncode                       uint8                // offset: 4, size: 1 bytes
	Attackspeed                      uint16               // offset: 5, size: 2 bytes
	Damageindex                      uint8                // offset: 7, size: 1 bytes
	Attacksequence                   uint8                // offset: 8, size: 1 bytes
}

func (s *PROTO_NC_BAT_SWING_START_CMD_DATA) FixedWireSize() int { return 9 }

// PROTO_NC_BAT_TARGETCHANGE_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_TARGETCHANGE_CMD_DATA struct {
	Target                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_TARGETCHANGE_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_TARGETINFO_CMD (sizeof = 30 bytes)
type PROTO_NC_BAT_TARGETINFO_CMD_DATA struct {
	Order                            uint8                // offset: 0, size: 1 bytes
	Targethandle                     uint16               // offset: 1, size: 2 bytes
	Targethp                         uint32               // offset: 3, size: 4 bytes
	Targetmaxhp                      uint32               // offset: 7, size: 4 bytes
	Targetsp                         uint32               // offset: 11, size: 4 bytes
	Targetmaxsp                      uint32               // offset: 15, size: 4 bytes
	Targetlp                         uint32               // offset: 19, size: 4 bytes
	Targetmaxlp                      uint32               // offset: 23, size: 4 bytes
	Targetlevel                      uint8                // offset: 27, size: 1 bytes
	Hpchangeorder                    uint16               // offset: 28, size: 2 bytes
}

func (s *PROTO_NC_BAT_TARGETINFO_CMD_DATA) FixedWireSize() int { return 30 }

// PROTO_NC_BAT_TARGET_REQ (sizeof = 2 bytes)
type PROTO_NC_BAT_TARGET_REQ struct {
	Target                           uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_TARGET_REQ) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_TOGGLESKILL_OFF_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_TOGGLESKILL_OFF_CMD_DATA struct {
	NSkillID                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_TOGGLESKILL_OFF_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_TOGGLESKILL_ON_CMD (sizeof = 2 bytes)
type PROTO_NC_BAT_TOGGLESKILL_ON_CMD_DATA struct {
	NSkillID                         uint16               // offset: 0, size: 2 bytes
}

func (s *PROTO_NC_BAT_TOGGLESKILL_ON_CMD_DATA) FixedWireSize() int { return 2 }

// PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD (sizeof = 4 bytes)
// Declarada previamente en el paquete protocol (PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD_DATA)
