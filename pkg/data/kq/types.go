// Package kq models and loads the official Kingdom Quest data tables and wire structures.
// All struct layouts, sizes, and status enums match WorldManager.pdb, Zone.pdb, and the
// disassembled instructions of WorldManager.exe and Zone.exe.
package kq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// Official C++ sizeof constants from WorldManager.pdb / Zone.pdb.
const (
	SizeOfKingdomQuest        = 180 // KINGDOM_QUEST (WorldManager.pdb, WorldManager.exe @ 0x00451B46)
	SizeOfKingdomQuestMap     = 363 // KINGDOM_QUEST_MAP (WorldManager.pdb, WorldManager.exe @ 0x00453E77)
	SizeOfKQTeam              = 22  // KQTeam (WorldManager.pdb, WorldManager.exe @ 0x00451D86)
	SizeOfKQIsVote            = 3   // KQIsVote (WorldManager.pdb, WorldManager.exe @ 0x00451B9A)
	SizeOfKQVoteMajorityRate  = 1   // KQVoteMajorityRate (WorldManager.pdb, WorldManager.exe @ 0x00451B66)
	SizeOfKingdomQuestRew     = 128 // KINGDOM_QUEST_REW (Zone.pdb)
	SizeOfKQItem              = 38  // KQItem (Zone.pdb)
	SizeOfShineReward         = 66  // ShineReward (Zone.pdb, Zone.exe @ 0x0064C830)
	SizeOfProtoKQMapInfo      = 26  // PROTO_KQ_MAP_INFO (WorldManager.pdb / Zone.pdb)
	SizeOfProtoKQInfoClient   = 141 // PROTO_KQ_INFO_CLIENT (WorldManager.pdb)
	SizeOfProtoKQInfo         = 377 // PROTO_KQ_INFO (WorldManager.pdb / Zone.pdb)
	SizeOfCTimeTM             = 36  // struct tm (MSVC x86)
	SizeOfShineXYType         = 8   // SHINE_XY_TYPE
)

// KQStatus represents the official byte state of a Kingdom Quest room in PROTO_KQ_INFO::Status
// (offset +0x04 in PROTO_KQ_INFO, +0x08 in CKQ entry stride 0x179).
// Every value and transition below is extracted directly from WorldManager.exe.
type KQStatus uint8

const (
	// KQStatusScheduled (0x00): initial state set in CKQServer::AddNewScheduleList (exe: 0x004556DE),
	// polled by CKQServer::DoSetMakeRoom (exe: 0x00455998).
	KQStatusScheduled KQStatus = 0x00

	// KQStatusMakingRoom (0x01): set in CKQServer::DoSetMakeRoom after AllocMapLink +
	// Send_NC_KQ_W2Z_MAKE_REQ (exe: 0x004559BA); checked in DoSetNoMapBF (exe: 0x004522A0)
	// and SetJoining (exe: 0x004541EE).
	KQStatusMakingRoom KQStatus = 0x01

	// KQStatusJoining (0x02): set in CKQServer::SetJoining (exe: 0x004542A4) on
	// fc_NC_KQ_Z2W_MAKE_ACK == 0x981 (exe: 0x0042AA0E); checked in CKQ::IsJoinable (exe: 0x0045CF20),
	// DoSetStart (exe: 0x00456046), Send_NC_KQ_JOINING_ALARM_CMD (exe: 0x0045296F),
	// Send_NC_KQ_JOINING_ALARM_END_CMD (exe: 0x00452AAB), Send_NC_KQ_JOINING_ALARM_LIST (exe: 0x00452BDB),
	// and Recv_NC_KQ_TEAM_SELECT_REQ (exe: 0x00453822).
	KQStatusJoining KQStatus = 0x02

	// KQStatusCountdown (0x03): set in CKQServer::DoSetStart when NumOfJoiner >= MaxPlayers
	// (exe: 0x0045608B) or after StartTime+StartWaitTime with MinPlayers & KQTeam_CanKQStart
	// (exe: 0x00456142), initializing RunCounter=0x0A (exe: 0x0045608F, 0x00456146).
	KQStatusCountdown KQStatus = 0x03

	// KQStatusRunning (0x04): set in CKQServer::DoSetStart when Status==3 && RunCounter==0
	// (exe: 0x00455FE1); checked in InKQStatusRunning (exe: 0x00452D98), LinkInKQ (exe: 0x00454E88),
	// LinkOutKQ (exe: 0x00454E29), Recv_NC_KQ_VOTE_START_REQ (exe: 0x00456415), and
	// Recv_NC_KQ_VOTE_START_CHECK_REQ (exe: 0x004533C2). Visible list filter is Status <= 4
	// (fc_NC_KQ_LIST_REQ @ 0x004139E2 / 0x00413A89, Ack_NC_KQ_LIST_REFRESH @ 0x00455A90).
	KQStatusRunning KQStatus = 0x04

	// KQStatusDone (0x05): set in CKQServer::SetDone (exe: 0x00454892) via fc_NC_KQ_Z2W_END_CMD
	// (exe: 0x0042AB33). Covered by finished range 5..10 in DelOldShceduleList (exe: 0x004540B2).
	KQStatusDone KQStatus = 0x05

	// KQStatusDoneSkip (0x06): set in CKQServer::SetDoneSkip (exe: 0x004548F8) from DoSetStart
	// (exe: 0x00456124, reason=2) or KQTeam_CanKQStart (exe: 0x004553FC reason=2, 0x00455443 reason=3).
	KQStatusDoneSkip KQStatus = 0x06

	// KQStatusReserved7 (0x07): falls within finished range 5..10 in DelOldShceduleList
	// (exe: 0x004540B2..0x004540B6), though no mov instruction writes 0x07 in WorldManager.exe.
	KQStatusReserved7 KQStatus = 0x07

	// KQStatusNoMap (0x08): set in DoSetMakeRoom when AllocMapLink fails (exe: 0x004559C4) or
	// in SetNoMapBF (exe: 0x004524AE) via fc_NC_KQ_Z2W_MAKE_ACK != 0x981 (exe: 0x0042AA22).
	KQStatusNoMap KQStatus = 0x08

	// KQStatusMakeTimeout (0x09): set in DoSetNoMapBF (exe: 0x00452300) when Status == 1 and
	// ScheduleTime + 60s <= CurrentTime without Z2W_MAKE_ACK.
	KQStatusMakeTimeout KQStatus = 0x09

	// KQStatusNotFound (0x0A): returned in fc_NC_KQ_STATUS_REQ (exe: 0x00413C6A) when
	// CKQ::Get(handle) is NULL; also within finished range 5..10 in DelOldShceduleList (exe: 0x004540B2).
	KQStatusNotFound KQStatus = 0x0A

	// KQStatusDeleting (0x0B): set in DelOldShceduleList (exe: 0x004540FD) and swept at 0x00454126,
	// in SetJoining on invalid prior status (exe: 0x00454227), and in fc_NC_OPTOOL_KQ_DELETE_CMD
	// (exe: 0x00425DF5).
	KQStatusDeleting KQStatus = 0x0B
)

// KQTeamType models enum KQ_TEAM_TYPE (WorldManager.pdb / Zone.pdb).
type KQTeamType uint8

const (
	KQTeamRed  KQTeamType = 0 // KQTT_RED
	KQTeamBlue KQTeamType = 1 // KQTT_BLUE
	KQTeamMax  KQTeamType = 2 // KQTT_MAX
)

// KQTeamDivideType models enum KQTeamDivideType (WorldManager.pdb).
type KQTeamDivideTypeVal uint16

const (
	KQTeamDivideNone       KQTeamDivideTypeVal = 0 // KQTD_NONE
	KQTeamDivideRandom     KQTeamDivideTypeVal = 1 // KQTD_RANDOM
	KQTeamDivideUserSelect KQTeamDivideTypeVal = 2 // KQTD_USERSELECT
	KQTeamDivideMax        KQTeamDivideTypeVal = 3 // MAX_KQTEAMDIVIDETYPE
)

// KQVotingType models enum KQ_VOTING_TYPE (WorldManager.pdb).
type KQVotingType uint32

const (
	KQVotingCancel KQVotingType = 0 // KVT_CANCEL
	KQVotingYes    KQVotingType = 1 // KVT_YES
	KQVotingNo     KQVotingType = 2 // KVT_NO
	KQVotingMax    KQVotingType = 3 // KVT_MAX
)

// KQMedal models enum KQ_Medal (Zone.pdb).
type KQMedal uint8

const (
	KQMedalNone   KQMedal = 0 // KQ_NONE
	KQMedalGold   KQMedal = 1 // KQ_GOLD
	KQMedalSilver KQMedal = 2 // KQ_SILVER
	KQMedalCopper KQMedal = 3 // KQ_COUPER
)

// CTimeTM models MSVC x86 struct tm (36 bytes).
type CTimeTM struct {
	TmSec   int32
	TmMin   int32
	TmHour  int32
	TmMday  int32
	TmMon   int32
	TmYear  int32
	TmWday  int32
	TmYday  int32
	TmIsdst int32
}

// CTimeTMFromTime converts a Go time.Time to MSVC struct tm.
func CTimeTMFromTime(t time.Time) CTimeTM {
	if t.IsZero() {
		return CTimeTM{}
	}
	return CTimeTM{
		TmSec:   int32(t.Second()),
		TmMin:   int32(t.Minute()),
		TmHour:  int32(t.Hour()),
		TmMday:  int32(t.Day()),
		TmMon:   int32(t.Month() - 1),
		TmYear:  int32(t.Year() - 1900),
		TmWday:  int32(t.Weekday()),
		TmYday:  int32(t.YearDay() - 1),
		TmIsdst: 0,
	}
}

// ToTime converts CTimeTM to Go time.Time in loc (defaults to time.Local if nil).
func (tm CTimeTM) ToTime(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(
		int(tm.TmYear)+1900,
		time.Month(tm.TmMon)+1,
		int(tm.TmMday),
		int(tm.TmHour),
		int(tm.TmMin),
		int(tm.TmSec),
		0,
		loc,
	)
}

func (tm CTimeTM) marshalInto(dst []byte) {
	binary.LittleEndian.PutUint32(dst[0:4], uint32(tm.TmSec))
	binary.LittleEndian.PutUint32(dst[4:8], uint32(tm.TmMin))
	binary.LittleEndian.PutUint32(dst[8:12], uint32(tm.TmHour))
	binary.LittleEndian.PutUint32(dst[12:16], uint32(tm.TmMday))
	binary.LittleEndian.PutUint32(dst[16:20], uint32(tm.TmMon))
	binary.LittleEndian.PutUint32(dst[20:24], uint32(tm.TmYear))
	binary.LittleEndian.PutUint32(dst[24:28], uint32(tm.TmWday))
	binary.LittleEndian.PutUint32(dst[28:32], uint32(tm.TmYday))
	binary.LittleEndian.PutUint32(dst[32:36], uint32(tm.TmIsdst))
}

func unmarshalCTimeTM(src []byte) CTimeTM {
	return CTimeTM{
		TmSec:   int32(binary.LittleEndian.Uint32(src[0:4])),
		TmMin:   int32(binary.LittleEndian.Uint32(src[4:8])),
		TmHour:  int32(binary.LittleEndian.Uint32(src[8:12])),
		TmMday:  int32(binary.LittleEndian.Uint32(src[12:16])),
		TmMon:   int32(binary.LittleEndian.Uint32(src[16:20])),
		TmYear:  int32(binary.LittleEndian.Uint32(src[20:24])),
		TmWday:  int32(binary.LittleEndian.Uint32(src[24:28])),
		TmYday:  int32(binary.LittleEndian.Uint32(src[28:32])),
		TmIsdst: int32(binary.LittleEndian.Uint32(src[32:36])),
	}
}

func trimCString(b []byte) string {
	if idx := bytes.IndexByte(b, 0x00); idx >= 0 {
		return string(b[:idx])
	}
	return string(b)
}

func copyCString(dst []byte, s string) {
	clear(dst)
	copy(dst, s)
}

// SHINE_XY_TYPE models C++ SHINE_XY_TYPE (8 bytes).
type SHINE_XY_TYPE struct {
	X uint32
	Y uint32
}

// KINGDOM_QUEST models C++ struct KINGDOM_QUEST (180 bytes, KingdomQuest.shn).
type KINGDOM_QUEST struct {
	ID                 uint16    // +0x0000
	Title              [64]byte  // +0x0002
	LimitTime          uint16    // +0x0042
	STYear             uint8     // +0x0044
	STMonth            uint8     // +0x0045
	STDay              uint8     // +0x0046
	STHour             uint8     // +0x0047
	STMinute           uint8     // +0x0048
	STSecond           uint8     // +0x0049
	StartWaitTime      uint16    // +0x004a
	NextStartMode      uint8     // +0x004c
	NextStartDeleyMin  uint16    // +0x004d (official PDB spelling: NextStartDeleyMin)
	RepeatMode         uint8     // +0x004f
	RepeatCount        uint16    // +0x0050
	MinLevel           uint8     // +0x0052
	MaxLevel           uint8     // +0x0053
	MinPlayers         uint16    // +0x0054
	MaxPlayers         uint16    // +0x0056
	PlayerRepeatMode   uint8     // +0x0058
	PlayerRepeatCount  uint16    // +0x0059
	PlayerRevivalMode  uint8     // +0x005b
	PlayerRevivalCount uint8     // +0x005c
	DemandQuest        uint16    // +0x005d
	DemandItem         uint16    // +0x005f
	DemandMobKill      uint8     // +0x0061
	RewardIndex        uint32    // +0x0062 (enum KINGDOMQUEST_REW_ID)
	MapLink            [4]uint16 // +0x0066 (0xFFFF = unused)
	ScriptLanguage     [32]byte  // +0x006e
	InitValue          [32]byte  // +0x008e
	UseClass           uint32    // +0x00ae (enum UseClassType)
	DemandGender       [2]uint8  // +0x00b2
}

// Size returns the official C++ sizeof(KINGDOM_QUEST) = 180.
func (KINGDOM_QUEST) Size() int { return SizeOfKingdomQuest }

// TitleString returns Title as a trimmed Go string.
func (k KINGDOM_QUEST) TitleString() string { return trimCString(k.Title[:]) }

// ScriptLanguageString returns ScriptLanguage as a trimmed Go string.
func (k KINGDOM_QUEST) ScriptLanguageString() string { return trimCString(k.ScriptLanguage[:]) }

// InitValueString returns InitValue as a trimmed Go string.
func (k KINGDOM_QUEST) InitValueString() string { return trimCString(k.InitValue[:]) }

// MarshalBinary serializes KINGDOM_QUEST in packed little-endian format (180 bytes).
func (k KINGDOM_QUEST) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKingdomQuest)
	binary.LittleEndian.PutUint16(buf[0x00:0x02], k.ID)
	copy(buf[0x02:0x42], k.Title[:])
	binary.LittleEndian.PutUint16(buf[0x42:0x44], k.LimitTime)
	buf[0x44] = k.STYear
	buf[0x45] = k.STMonth
	buf[0x46] = k.STDay
	buf[0x47] = k.STHour
	buf[0x48] = k.STMinute
	buf[0x49] = k.STSecond
	binary.LittleEndian.PutUint16(buf[0x4a:0x4c], k.StartWaitTime)
	buf[0x4c] = k.NextStartMode
	binary.LittleEndian.PutUint16(buf[0x4d:0x4f], k.NextStartDeleyMin)
	buf[0x4f] = k.RepeatMode
	binary.LittleEndian.PutUint16(buf[0x50:0x52], k.RepeatCount)
	buf[0x52] = k.MinLevel
	buf[0x53] = k.MaxLevel
	binary.LittleEndian.PutUint16(buf[0x54:0x56], k.MinPlayers)
	binary.LittleEndian.PutUint16(buf[0x56:0x58], k.MaxPlayers)
	buf[0x58] = k.PlayerRepeatMode
	binary.LittleEndian.PutUint16(buf[0x59:0x5b], k.PlayerRepeatCount)
	buf[0x5b] = k.PlayerRevivalMode
	buf[0x5c] = k.PlayerRevivalCount
	binary.LittleEndian.PutUint16(buf[0x5d:0x5f], k.DemandQuest)
	binary.LittleEndian.PutUint16(buf[0x5f:0x61], k.DemandItem)
	buf[0x61] = k.DemandMobKill
	binary.LittleEndian.PutUint32(buf[0x62:0x66], k.RewardIndex)
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint16(buf[0x66+i*2:0x68+i*2], k.MapLink[i])
	}
	copy(buf[0x6e:0x8e], k.ScriptLanguage[:])
	copy(buf[0x8e:0xae], k.InitValue[:])
	binary.LittleEndian.PutUint32(buf[0xae:0xb2], k.UseClass)
	buf[0xb2] = k.DemandGender[0]
	buf[0xb3] = k.DemandGender[1]
	return buf, nil
}

// UnmarshalBinary deserializes 180 bytes into KINGDOM_QUEST.
func (k *KINGDOM_QUEST) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKingdomQuest {
		return fmt.Errorf("KINGDOM_QUEST buffer too short: got %d, want %d", len(data), SizeOfKingdomQuest)
	}
	k.ID = binary.LittleEndian.Uint16(data[0x00:0x02])
	copy(k.Title[:], data[0x02:0x42])
	k.LimitTime = binary.LittleEndian.Uint16(data[0x42:0x44])
	k.STYear = data[0x44]
	k.STMonth = data[0x45]
	k.STDay = data[0x46]
	k.STHour = data[0x47]
	k.STMinute = data[0x48]
	k.STSecond = data[0x49]
	k.StartWaitTime = binary.LittleEndian.Uint16(data[0x4a:0x4c])
	k.NextStartMode = data[0x4c]
	k.NextStartDeleyMin = binary.LittleEndian.Uint16(data[0x4d:0x4f])
	k.RepeatMode = data[0x4f]
	k.RepeatCount = binary.LittleEndian.Uint16(data[0x50:0x52])
	k.MinLevel = data[0x52]
	k.MaxLevel = data[0x53]
	k.MinPlayers = binary.LittleEndian.Uint16(data[0x54:0x56])
	k.MaxPlayers = binary.LittleEndian.Uint16(data[0x56:0x58])
	k.PlayerRepeatMode = data[0x58]
	k.PlayerRepeatCount = binary.LittleEndian.Uint16(data[0x59:0x5b])
	k.PlayerRevivalMode = data[0x5b]
	k.PlayerRevivalCount = data[0x5c]
	k.DemandQuest = binary.LittleEndian.Uint16(data[0x5d:0x5f])
	k.DemandItem = binary.LittleEndian.Uint16(data[0x5f:0x61])
	k.DemandMobKill = data[0x61]
	k.RewardIndex = binary.LittleEndian.Uint32(data[0x62:0x66])
	for i := 0; i < 4; i++ {
		k.MapLink[i] = binary.LittleEndian.Uint16(data[0x66+i*2 : 0x68+i*2])
	}
	copy(k.ScriptLanguage[:], data[0x6e:0x8e])
	copy(k.InitValue[:], data[0x8e:0xae])
	k.UseClass = binary.LittleEndian.Uint32(data[0xae:0xb2])
	k.DemandGender[0] = data[0xb2]
	k.DemandGender[1] = data[0xb3]
	return nil
}

// KINGDOM_QUEST_MAP models C++ struct KINGDOM_QUEST_MAP (363 bytes, KingdomQuestMap.shn).
type KINGDOM_QUEST_MAP struct {
	NumOfMap uint8        // +0x0000
	BaseMap  [32]byte     // +0x0001
	Map      [10][32]byte // +0x0021
	Clear    [10]uint8    // +0x0161
}

// Size returns the official C++ sizeof(KINGDOM_QUEST_MAP) = 363.
func (KINGDOM_QUEST_MAP) Size() int { return SizeOfKingdomQuestMap }

// BaseMapString returns BaseMap as a trimmed Go string.
func (m KINGDOM_QUEST_MAP) BaseMapString() string { return trimCString(m.BaseMap[:]) }

// MapString returns Map[i] as a trimmed Go string ("" if out of range).
func (m KINGDOM_QUEST_MAP) MapString(i int) string {
	if i < 0 || i >= len(m.Map) {
		return ""
	}
	return trimCString(m.Map[i][:])
}

// MarshalBinary serializes KINGDOM_QUEST_MAP in packed format (363 bytes).
func (m KINGDOM_QUEST_MAP) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKingdomQuestMap)
	buf[0x00] = m.NumOfMap
	copy(buf[0x01:0x21], m.BaseMap[:])
	for i := 0; i < 10; i++ {
		off := 0x21 + i*32
		copy(buf[off:off+32], m.Map[i][:])
	}
	copy(buf[0x161:0x16b], m.Clear[:])
	return buf, nil
}

// UnmarshalBinary deserializes 363 bytes into KINGDOM_QUEST_MAP.
func (m *KINGDOM_QUEST_MAP) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKingdomQuestMap {
		return fmt.Errorf("KINGDOM_QUEST_MAP buffer too short: got %d, want %d", len(data), SizeOfKingdomQuestMap)
	}
	m.NumOfMap = data[0x00]
	copy(m.BaseMap[:], data[0x01:0x21])
	for i := 0; i < 10; i++ {
		off := 0x21 + i*32
		copy(m.Map[i][:], data[off:off+32])
	}
	copy(m.Clear[:], data[0x161:0x16b])
	return nil
}

// KQTeam models C++ struct KQTeam (22 bytes, KQTeam.shn).
type KQTeam struct {
	ID               uint16 // +0x0000
	MaxMemberGap     uint8  // +0x0002
	IsTeamPVP        uint8  // +0x0003
	KQTeamDivideType uint16 // +0x0004
	RegenXRed        uint32 // +0x0006
	RegenYRed        uint32 // +0x000a
	RegenXBlue       uint32 // +0x000e
	RegenYBlue       uint32 // +0x0012
}

// Size returns the official C++ sizeof(KQTeam) = 22.
func (KQTeam) Size() int { return SizeOfKQTeam }

// MarshalBinary serializes KQTeam in packed little-endian format (22 bytes).
func (t KQTeam) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKQTeam)
	binary.LittleEndian.PutUint16(buf[0x00:0x02], t.ID)
	buf[0x02] = t.MaxMemberGap
	buf[0x03] = t.IsTeamPVP
	binary.LittleEndian.PutUint16(buf[0x04:0x06], t.KQTeamDivideType)
	binary.LittleEndian.PutUint32(buf[0x06:0x0a], t.RegenXRed)
	binary.LittleEndian.PutUint32(buf[0x0a:0x0e], t.RegenYRed)
	binary.LittleEndian.PutUint32(buf[0x0e:0x12], t.RegenXBlue)
	binary.LittleEndian.PutUint32(buf[0x12:0x16], t.RegenYBlue)
	return buf, nil
}

// UnmarshalBinary deserializes 22 bytes into KQTeam.
func (t *KQTeam) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKQTeam {
		return fmt.Errorf("KQTeam buffer too short: got %d, want %d", len(data), SizeOfKQTeam)
	}
	t.ID = binary.LittleEndian.Uint16(data[0x00:0x02])
	t.MaxMemberGap = data[0x02]
	t.IsTeamPVP = data[0x03]
	t.KQTeamDivideType = binary.LittleEndian.Uint16(data[0x04:0x06])
	t.RegenXRed = binary.LittleEndian.Uint32(data[0x06:0x0a])
	t.RegenYRed = binary.LittleEndian.Uint32(data[0x0a:0x0e])
	t.RegenXBlue = binary.LittleEndian.Uint32(data[0x0e:0x12])
	t.RegenYBlue = binary.LittleEndian.Uint32(data[0x12:0x16])
	return nil
}

// KQIsVote models C++ struct KQIsVote (3 bytes, KQIsVote.shn).
type KQIsVote struct {
	ID     uint16 // +0x0000
	IsVote uint8  // +0x0002
}

// Size returns the official C++ sizeof(KQIsVote) = 3.
func (KQIsVote) Size() int { return SizeOfKQIsVote }

// MarshalBinary serializes KQIsVote (3 bytes).
func (v KQIsVote) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKQIsVote)
	binary.LittleEndian.PutUint16(buf[0:2], v.ID)
	buf[2] = v.IsVote
	return buf, nil
}

// UnmarshalBinary deserializes 3 bytes into KQIsVote.
func (v *KQIsVote) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKQIsVote {
		return fmt.Errorf("KQIsVote buffer too short: got %d, want %d", len(data), SizeOfKQIsVote)
	}
	v.ID = binary.LittleEndian.Uint16(data[0:2])
	v.IsVote = data[2]
	return nil
}

// KQVoteMajorityRate models C++ struct KQVoteMajorityRate (1 byte, KQVoteMajorityRate.shn).
type KQVoteMajorityRate struct {
	VoteAgreeRate uint8 // +0x0000
}

// Size returns the official C++ sizeof(KQVoteMajorityRate) = 1.
func (KQVoteMajorityRate) Size() int { return SizeOfKQVoteMajorityRate }

// MarshalBinary serializes KQVoteMajorityRate (1 byte).
func (r KQVoteMajorityRate) MarshalBinary() ([]byte, error) {
	return []byte{r.VoteAgreeRate}, nil
}

// UnmarshalBinary deserializes 1 byte into KQVoteMajorityRate.
func (r *KQVoteMajorityRate) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKQVoteMajorityRate {
		return fmt.Errorf("KQVoteMajorityRate buffer too short: got %d, want %d", len(data), SizeOfKQVoteMajorityRate)
	}
	r.VoteAgreeRate = data[0]
	return nil
}

// KINGDOM_QUEST_REW models C++ struct KINGDOM_QUEST_REW (128 bytes, KingdomQuestRew.shn).
type KINGDOM_QUEST_REW struct {
	ID           uint32     // +0x0000 (enum KINGDOMQUEST_REW_ID)
	IndexString  [32]byte   // +0x0004
	KQBoxItemIDX [32]byte   // +0x0024
	Reward       [15]uint16 // +0x0044
	RewardRate   [15]uint16 // +0x0062
}

// Size returns the official C++ sizeof(KINGDOM_QUEST_REW) = 128.
func (KINGDOM_QUEST_REW) Size() int { return SizeOfKingdomQuestRew }

// IndexName returns IndexString as a trimmed Go string.
func (r KINGDOM_QUEST_REW) IndexName() string { return trimCString(r.IndexString[:]) }

// BoxItemIndex returns KQBoxItemIDX as a trimmed Go string.
func (r KINGDOM_QUEST_REW) BoxItemIndex() string { return trimCString(r.KQBoxItemIDX[:]) }

// MarshalBinary serializes KINGDOM_QUEST_REW in packed little-endian format (128 bytes).
func (r KINGDOM_QUEST_REW) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKingdomQuestRew)
	binary.LittleEndian.PutUint32(buf[0x00:0x04], r.ID)
	copy(buf[0x04:0x24], r.IndexString[:])
	copy(buf[0x24:0x44], r.KQBoxItemIDX[:])
	for i := 0; i < 15; i++ {
		binary.LittleEndian.PutUint16(buf[0x44+i*2:0x46+i*2], r.Reward[i])
		binary.LittleEndian.PutUint16(buf[0x62+i*2:0x64+i*2], r.RewardRate[i])
	}
	return buf, nil
}

// UnmarshalBinary deserializes 128 bytes into KINGDOM_QUEST_REW.
func (r *KINGDOM_QUEST_REW) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKingdomQuestRew {
		return fmt.Errorf("KINGDOM_QUEST_REW buffer too short: got %d, want %d", len(data), SizeOfKingdomQuestRew)
	}
	r.ID = binary.LittleEndian.Uint32(data[0x00:0x04])
	copy(r.IndexString[:], data[0x04:0x24])
	copy(r.KQBoxItemIDX[:], data[0x24:0x44])
	for i := 0; i < 15; i++ {
		r.Reward[i] = binary.LittleEndian.Uint16(data[0x44+i*2 : 0x46+i*2])
		r.RewardRate[i] = binary.LittleEndian.Uint16(data[0x62+i*2 : 0x64+i*2])
	}
	return nil
}

// ShineRewardType mirrors the switch on ShineReward::RewardType in ShinePlayer::sp_KQReward
// (Zone.exe @ 0x0052DF3C..0x0052DFDB, jump table @ 0x0052E348).
type ShineRewardType uint8

const (
	ShineRewardTypeNone ShineRewardType = 0 // exe: 0x0052DF50
	ShineRewardTypeItem ShineRewardType = 1 // exe: 0x0052DF7D (tcm_ItemMake)
	ShineRewardTypeExp  ShineRewardType = 2 // exe: 0x0052DFA3 (totalExp += [edi+0x24])
	ShineRewardTypeCen  ShineRewardType = 3 // exe: 0x0052DFBF (totalMoney += [edi+0x24])
	ShineRewardTypeFame ShineRewardType = 4 // exe: 0x0052DFDB (totalFame += [edi+0x24])
)

// SHINE_REWARD models C++ struct ShineReward (66 bytes, ShineReward.shn).
// Loaded by RewardData::rd_Load (Zone.exe @ 0x0064C830).
type SHINE_REWARD struct {
	RewardHandle uint16          // +0x0000
	RewardType   ShineRewardType // +0x0002
	Argument     [33]byte        // +0x0003 (item InxName or ItemGroup name)
	Quantity     uint32          // +0x0024 (Exp / Cen / Fame / item count)
	Upgrade      [10]uint16      // +0x0028 (upgrade probability weights)
	OptionDegree uint16          // +0x003C
	TitleDegree  uint32          // +0x003E
}

// Size returns the official C++ sizeof(ShineReward) = 66.
func (SHINE_REWARD) Size() int { return SizeOfShineReward }

// ArgumentString returns Argument as a trimmed Go string.
func (s SHINE_REWARD) ArgumentString() string { return trimCString(s.Argument[:]) }

// MarshalBinary serializes SHINE_REWARD in packed little-endian format (66 bytes).
func (s SHINE_REWARD) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfShineReward)
	binary.LittleEndian.PutUint16(buf[0x00:0x02], s.RewardHandle)
	buf[0x02] = byte(s.RewardType)
	copy(buf[0x03:0x24], s.Argument[:])
	binary.LittleEndian.PutUint32(buf[0x24:0x28], s.Quantity)
	for i := 0; i < 10; i++ {
		binary.LittleEndian.PutUint16(buf[0x28+i*2:0x2A+i*2], s.Upgrade[i])
	}
	binary.LittleEndian.PutUint16(buf[0x3C:0x3E], s.OptionDegree)
	binary.LittleEndian.PutUint32(buf[0x3E:0x42], s.TitleDegree)
	return buf, nil
}

// UnmarshalBinary deserializes 66 bytes into SHINE_REWARD.
func (s *SHINE_REWARD) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfShineReward {
		return fmt.Errorf("SHINE_REWARD buffer too short: got %d, want %d", len(data), SizeOfShineReward)
	}
	s.RewardHandle = binary.LittleEndian.Uint16(data[0x00:0x02])
	s.RewardType = ShineRewardType(data[0x02])
	copy(s.Argument[:], data[0x03:0x24])
	s.Quantity = binary.LittleEndian.Uint32(data[0x24:0x28])
	for i := 0; i < 10; i++ {
		s.Upgrade[i] = binary.LittleEndian.Uint16(data[0x28+i*2 : 0x2A+i*2])
	}
	s.OptionDegree = binary.LittleEndian.Uint16(data[0x3C:0x3E])
	s.TitleDegree = binary.LittleEndian.Uint32(data[0x3E:0x42])
	return nil
}

// KQItem models C++ struct KQItem (38 bytes, KQItem.shn).
type KQItem struct {
	ItemIndex      [32]byte // +0x0000
	MoveSpdRate    uint16   // +0x0020
	AbsoluteAttack uint16   // +0x0022
	PickupLimit    uint16   // +0x0024
}

// Size returns the official C++ sizeof(KQItem) = 38.
func (KQItem) Size() int { return SizeOfKQItem }

// ItemIndexString returns ItemIndex as a trimmed Go string.
func (k KQItem) ItemIndexString() string { return trimCString(k.ItemIndex[:]) }

// MarshalBinary serializes KQItem in packed little-endian format (38 bytes).
func (k KQItem) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfKQItem)
	copy(buf[0x00:0x20], k.ItemIndex[:])
	binary.LittleEndian.PutUint16(buf[0x20:0x22], k.MoveSpdRate)
	binary.LittleEndian.PutUint16(buf[0x22:0x24], k.AbsoluteAttack)
	binary.LittleEndian.PutUint16(buf[0x24:0x26], k.PickupLimit)
	return buf, nil
}

// UnmarshalBinary deserializes 38 bytes into KQItem.
func (k *KQItem) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfKQItem {
		return fmt.Errorf("KQItem buffer too short: got %d, want %d", len(data), SizeOfKQItem)
	}
	copy(k.ItemIndex[:], data[0x00:0x20])
	k.MoveSpdRate = binary.LittleEndian.Uint16(data[0x20:0x22])
	k.AbsoluteAttack = binary.LittleEndian.Uint16(data[0x22:0x24])
	k.PickupLimit = binary.LittleEndian.Uint16(data[0x24:0x26])
	return nil
}

// PROTO_KQ_MAP_INFO models C++ struct PROTO_KQ_MAP_INFO (26 bytes, WorldManager.pdb / Zone.pdb).
// Matches CKQServer::AllocMapLink (WorldManager.exe @ 0x00453EB0).
type PROTO_KQ_MAP_INFO struct {
	MapIndex uint8    // +0x0000
	MapBase  [12]byte // +0x0001 (Name3)
	MapName  [12]byte // +0x000d (Name3)
	MapClear uint8    // +0x0019
}

// Size returns the official C++ sizeof(PROTO_KQ_MAP_INFO) = 26.
func (PROTO_KQ_MAP_INFO) Size() int { return SizeOfProtoKQMapInfo }

// MapBaseString returns MapBase as a trimmed Go string.
func (m PROTO_KQ_MAP_INFO) MapBaseString() string { return trimCString(m.MapBase[:]) }

// MapNameString returns MapName as a trimmed Go string.
func (m PROTO_KQ_MAP_INFO) MapNameString() string { return trimCString(m.MapName[:]) }

// SetMapBase copies s into MapBase (zero-padded 12 bytes).
func (m *PROTO_KQ_MAP_INFO) SetMapBase(s string) { copyCString(m.MapBase[:], s) }

// SetMapName copies s into MapName (zero-padded 12 bytes).
func (m *PROTO_KQ_MAP_INFO) SetMapName(s string) { copyCString(m.MapName[:], s) }

// MarshalBinary serializes PROTO_KQ_MAP_INFO (26 bytes).
func (m PROTO_KQ_MAP_INFO) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfProtoKQMapInfo)
	buf[0x00] = m.MapIndex
	copy(buf[0x01:0x0d], m.MapBase[:])
	copy(buf[0x0d:0x19], m.MapName[:])
	buf[0x19] = m.MapClear
	return buf, nil
}

// UnmarshalBinary deserializes 26 bytes into PROTO_KQ_MAP_INFO.
func (m *PROTO_KQ_MAP_INFO) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfProtoKQMapInfo {
		return fmt.Errorf("PROTO_KQ_MAP_INFO buffer too short: got %d, want %d", len(data), SizeOfProtoKQMapInfo)
	}
	m.MapIndex = data[0x00]
	copy(m.MapBase[:], data[0x01:0x0d])
	copy(m.MapName[:], data[0x0d:0x19])
	m.MapClear = data[0x19]
	return nil
}

// PROTO_KQ_INFO_CLIENT models C++ struct PROTO_KQ_INFO_CLIENT (141 bytes, WorldManager.pdb).
type PROTO_KQ_INFO_CLIENT struct {
	Handle             uint32   // +0x0000
	Status             KQStatus // +0x0004
	NumOfJoiner        uint16   // +0x0005
	ID                 uint16   // +0x0007
	Title              [64]byte // +0x0009
	LimitTime          uint16   // +0x0049
	StartTime          int32    // +0x004b
	TmStartTime        CTimeTM  // +0x004f (36 bytes)
	StartWaitTime      uint16   // +0x0073
	MinLevel           uint8    // +0x0075
	MaxLevel           uint8    // +0x0076
	MinPlayers         uint16   // +0x0077
	MaxPlayers         uint16   // +0x0079
	PlayerRepeatMode   uint8    // +0x007b
	PlayerRepeatCount  uint16   // +0x007c
	PlayerRevivalMode  uint8    // +0x007e
	PlayerRevivalCount uint8    // +0x007f
	DemandQuest        uint16   // +0x0080
	DemandItem         uint16   // +0x0082
	DemandClass        int64    // +0x0084
	DemandGender       uint8    // +0x008c
}

// Size returns the official C++ sizeof(PROTO_KQ_INFO_CLIENT) = 141.
func (PROTO_KQ_INFO_CLIENT) Size() int { return SizeOfProtoKQInfoClient }

// TitleString returns Title as a trimmed Go string.
func (c PROTO_KQ_INFO_CLIENT) TitleString() string { return trimCString(c.Title[:]) }

// SetTitle copies s into Title (zero-padded 64 bytes).
func (c *PROTO_KQ_INFO_CLIENT) SetTitle(s string) { copyCString(c.Title[:], s) }

// MarshalBinary serializes PROTO_KQ_INFO_CLIENT (141 bytes).
func (c PROTO_KQ_INFO_CLIENT) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfProtoKQInfoClient)
	binary.LittleEndian.PutUint32(buf[0x00:0x04], c.Handle)
	buf[0x04] = byte(c.Status)
	binary.LittleEndian.PutUint16(buf[0x05:0x07], c.NumOfJoiner)
	binary.LittleEndian.PutUint16(buf[0x07:0x09], c.ID)
	copy(buf[0x09:0x49], c.Title[:])
	binary.LittleEndian.PutUint16(buf[0x49:0x4b], c.LimitTime)
	binary.LittleEndian.PutUint32(buf[0x4b:0x4f], uint32(c.StartTime))
	c.TmStartTime.marshalInto(buf[0x4f:0x73])
	binary.LittleEndian.PutUint16(buf[0x73:0x75], c.StartWaitTime)
	buf[0x75] = c.MinLevel
	buf[0x76] = c.MaxLevel
	binary.LittleEndian.PutUint16(buf[0x77:0x79], c.MinPlayers)
	binary.LittleEndian.PutUint16(buf[0x79:0x7b], c.MaxPlayers)
	buf[0x7b] = c.PlayerRepeatMode
	binary.LittleEndian.PutUint16(buf[0x7c:0x7e], c.PlayerRepeatCount)
	buf[0x7e] = c.PlayerRevivalMode
	buf[0x7f] = c.PlayerRevivalCount
	binary.LittleEndian.PutUint16(buf[0x80:0x82], c.DemandQuest)
	binary.LittleEndian.PutUint16(buf[0x82:0x84], c.DemandItem)
	binary.LittleEndian.PutUint64(buf[0x84:0x8c], uint64(c.DemandClass))
	buf[0x8c] = c.DemandGender
	return buf, nil
}

// UnmarshalBinary deserializes 141 bytes into PROTO_KQ_INFO_CLIENT.
func (c *PROTO_KQ_INFO_CLIENT) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfProtoKQInfoClient {
		return fmt.Errorf("PROTO_KQ_INFO_CLIENT buffer too short: got %d, want %d", len(data), SizeOfProtoKQInfoClient)
	}
	c.Handle = binary.LittleEndian.Uint32(data[0x00:0x04])
	c.Status = KQStatus(data[0x04])
	c.NumOfJoiner = binary.LittleEndian.Uint16(data[0x05:0x07])
	c.ID = binary.LittleEndian.Uint16(data[0x07:0x09])
	copy(c.Title[:], data[0x09:0x49])
	c.LimitTime = binary.LittleEndian.Uint16(data[0x49:0x4b])
	c.StartTime = int32(binary.LittleEndian.Uint32(data[0x4b:0x4f]))
	c.TmStartTime = unmarshalCTimeTM(data[0x4f:0x73])
	c.StartWaitTime = binary.LittleEndian.Uint16(data[0x73:0x75])
	c.MinLevel = data[0x75]
	c.MaxLevel = data[0x76]
	c.MinPlayers = binary.LittleEndian.Uint16(data[0x77:0x79])
	c.MaxPlayers = binary.LittleEndian.Uint16(data[0x79:0x7b])
	c.PlayerRepeatMode = data[0x7b]
	c.PlayerRepeatCount = binary.LittleEndian.Uint16(data[0x7c:0x7e])
	c.PlayerRevivalMode = data[0x7e]
	c.PlayerRevivalCount = data[0x7f]
	c.DemandQuest = binary.LittleEndian.Uint16(data[0x80:0x82])
	c.DemandItem = binary.LittleEndian.Uint16(data[0x82:0x84])
	c.DemandClass = int64(binary.LittleEndian.Uint64(data[0x84:0x8c]))
	c.DemandGender = data[0x8c]
	return nil
}

// PROTO_KQ_INFO models C++ struct PROTO_KQ_INFO (377 bytes, WorldManager.pdb / Zone.pdb).
// Used as the element of CKQ (stride 0x179) and the payload of W2Z_MAKE_REQ / W2Z_START_CMD.
type PROTO_KQ_INFO struct {
	PROTO_KQ_INFO_CLIENT                   // +0x0000 (141 bytes)
	NextStartMode        uint8             // +0x008d
	NextStartDelayMin    uint16            // +0x008e
	RepeatMode           uint8             // +0x0090
	RepeatCount          uint16            // +0x0091
	RewardIndex          uint16            // +0x0093
	DemandMobKill        uint16            // +0x0095
	ScheduleTime         int32             // +0x0097
	TmScheduleTime       CTimeTM           // +0x009b (36 bytes)
	RunCounter           uint8             // +0x00bf
	MapLink              [4]PROTO_KQ_MAP_INFO // +0x00c0 (104 bytes)
	ScriptLanguage       [32]byte          // +0x0128
	ScriptInitValue      [32]byte          // +0x0148
	IsTeamPVP            uint8             // +0x0168
	TeamRegenXY          [2]SHINE_XY_TYPE  // +0x0169 (16 bytes)
}

// Size returns the official C++ sizeof(PROTO_KQ_INFO) = 377.
func (PROTO_KQ_INFO) Size() int { return SizeOfProtoKQInfo }

// ScriptLanguageString returns ScriptLanguage as a trimmed Go string.
func (q PROTO_KQ_INFO) ScriptLanguageString() string { return trimCString(q.ScriptLanguage[:]) }

// ScriptInitValueString returns ScriptInitValue as a trimmed Go string.
func (q PROTO_KQ_INFO) ScriptInitValueString() string { return trimCString(q.ScriptInitValue[:]) }

// SetScriptLanguage copies s into ScriptLanguage (zero-padded 32 bytes).
func (q *PROTO_KQ_INFO) SetScriptLanguage(s string) { copyCString(q.ScriptLanguage[:], s) }

// SetScriptInitValue copies s into ScriptInitValue (zero-padded 32 bytes).
func (q *PROTO_KQ_INFO) SetScriptInitValue(s string) { copyCString(q.ScriptInitValue[:], s) }

// MarshalBinary serializes PROTO_KQ_INFO in packed little-endian format (377 bytes).
func (q PROTO_KQ_INFO) MarshalBinary() ([]byte, error) {
	buf := make([]byte, SizeOfProtoKQInfo)
	clientBytes, err := q.PROTO_KQ_INFO_CLIENT.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(buf[0x00:0x8d], clientBytes)
	buf[0x8d] = q.NextStartMode
	binary.LittleEndian.PutUint16(buf[0x8e:0x90], q.NextStartDelayMin)
	buf[0x90] = q.RepeatMode
	binary.LittleEndian.PutUint16(buf[0x91:0x93], q.RepeatCount)
	binary.LittleEndian.PutUint16(buf[0x93:0x95], q.RewardIndex)
	binary.LittleEndian.PutUint16(buf[0x95:0x97], q.DemandMobKill)
	binary.LittleEndian.PutUint32(buf[0x97:0x9b], uint32(q.ScheduleTime))
	q.TmScheduleTime.marshalInto(buf[0x9b:0xbf])
	buf[0xbf] = q.RunCounter
	for i := 0; i < 4; i++ {
		mb, err := q.MapLink[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		off := 0xc0 + i*SizeOfProtoKQMapInfo
		copy(buf[off:off+SizeOfProtoKQMapInfo], mb)
	}
	copy(buf[0x128:0x148], q.ScriptLanguage[:])
	copy(buf[0x148:0x168], q.ScriptInitValue[:])
	buf[0x168] = q.IsTeamPVP
	for i := 0; i < 2; i++ {
		off := 0x169 + i*SizeOfShineXYType
		binary.LittleEndian.PutUint32(buf[off:off+4], q.TeamRegenXY[i].X)
		binary.LittleEndian.PutUint32(buf[off+4:off+8], q.TeamRegenXY[i].Y)
	}
	return buf, nil
}

// UnmarshalBinary deserializes 377 bytes into PROTO_KQ_INFO.
func (q *PROTO_KQ_INFO) UnmarshalBinary(data []byte) error {
	if len(data) < SizeOfProtoKQInfo {
		return fmt.Errorf("PROTO_KQ_INFO buffer too short: got %d, want %d", len(data), SizeOfProtoKQInfo)
	}
	if err := q.PROTO_KQ_INFO_CLIENT.UnmarshalBinary(data[0x00:0x8d]); err != nil {
		return err
	}
	q.NextStartMode = data[0x8d]
	q.NextStartDelayMin = binary.LittleEndian.Uint16(data[0x8e:0x90])
	q.RepeatMode = data[0x90]
	q.RepeatCount = binary.LittleEndian.Uint16(data[0x91:0x93])
	q.RewardIndex = binary.LittleEndian.Uint16(data[0x93:0x95])
	q.DemandMobKill = binary.LittleEndian.Uint16(data[0x95:0x97])
	q.ScheduleTime = int32(binary.LittleEndian.Uint32(data[0x97:0x9b]))
	q.TmScheduleTime = unmarshalCTimeTM(data[0x9b:0xbf])
	q.RunCounter = data[0xbf]
	for i := 0; i < 4; i++ {
		off := 0xc0 + i*SizeOfProtoKQMapInfo
		if err := q.MapLink[i].UnmarshalBinary(data[off : off+SizeOfProtoKQMapInfo]); err != nil {
			return err
		}
	}
	copy(q.ScriptLanguage[:], data[0x128:0x148])
	copy(q.ScriptInitValue[:], data[0x148:0x168])
	q.IsTeamPVP = data[0x168]
	for i := 0; i < 2; i++ {
		off := 0x169 + i*SizeOfShineXYType
		q.TeamRegenXY[i].X = binary.LittleEndian.Uint32(data[off : off+4])
		q.TeamRegenXY[i].Y = binary.LittleEndian.Uint32(data[off+4 : off+8])
	}
	return nil
}
