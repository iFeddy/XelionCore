package protocol

import (
	"time"
)

// Packet builders oficiales para el sistema de Kingdom Quests (KQ).
// Nomenclatura y offsets 100% idénticos a WorldManager.pdb / WorldManager.md (Dept 22).

// ProtoKQInfoClient modela PROTO_KQ_INFO_CLIENT (141 bytes).
type ProtoKQInfoClient struct {
	Handle             uint32
	Status             byte   // 0=Scheduled, 1=Recruiting, 2=Running, 3=Ended
	NumOfJoiner        uint16
	ID                 uint16
	Title              string // char[64]
	LimitTime          uint16 // seconds
	StartTime          time.Time
	StartWaitTime      uint16
	MinLevel           byte
	MaxLevel           byte
	MinPlayers         uint16
	MaxPlayers         uint16
	PlayerRepeatMode   byte
	PlayerRepeatCount  uint16
	PlayerRevivalMode  byte
	PlayerRevivalCount byte
	DemandQuest        uint16
	DemandItem         uint16
	DemandClass        int64
	DemandGender       byte
}

func writeProtoKQInfoClient(pw *PacketWriter, q ProtoKQInfoClient) {
	pw.WriteUint32(q.Handle)
	pw.WriteByte(q.Status)
	pw.WriteUint16(q.NumOfJoiner)
	pw.WriteUint16(q.ID)
	pw.WriteString(q.Title, 64)
	pw.WriteUint16(q.LimitTime)
	pw.WriteInt32(int32(q.StartTime.Unix()))
	pw.WriteBytes(packTm(q.StartTime))
	pw.WriteUint16(q.StartWaitTime)
	pw.WriteByte(q.MinLevel)
	pw.WriteByte(q.MaxLevel)
	pw.WriteUint16(q.MinPlayers)
	pw.WriteUint16(q.MaxPlayers)
	pw.WriteByte(q.PlayerRepeatMode)
	pw.WriteUint16(q.PlayerRepeatCount)
	pw.WriteByte(q.PlayerRevivalMode)
	pw.WriteByte(q.PlayerRevivalCount)
	pw.WriteUint16(q.DemandQuest)
	pw.WriteUint16(q.DemandItem)
	pw.WriteInt64(q.DemandClass)
	pw.WriteByte(q.DemandGender)
}

// New_PROTO_NC_KQ_LIST_TIME_ACK construye el reloj oficial de sincronización de KQ (OpCode 0x581C).
// Struct oficial: PROTO_NC_KQ_LIST_TIME_ACK (40 bytes: ServerTime[4], tm_ServerTime[36])
func New_PROTO_NC_KQ_LIST_TIME_ACK(serverTime time.Time) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_LIST_TIME_ACK)
	pw.WriteInt32(int32(serverTime.Unix()))
	pw.WriteBytes(packTm(serverTime))
	return pw
}

// New_PROTO_NC_KQ_LIST_ACK construye la lista completa de Kingdom Quests (OpCode 0x5802).
// Struct oficial: PROTO_NC_KQ_LIST_ACK (50 bytes + N * 141 PROTO_KQ_INFO_CLIENT)
func New_PROTO_NC_KQ_LIST_ACK(serverTime time.Time, startHandle, endHandle uint32, quests []ProtoKQInfoClient) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_LIST_ACK)
	pw.WriteInt32(int32(serverTime.Unix()))
	pw.WriteBytes(packTm(serverTime))
	pw.WriteUint32(startHandle)
	pw.WriteUint32(endHandle)
	pw.WriteUint16(uint16(len(quests)))
	for _, q := range quests {
		writeProtoKQInfoClient(pw, q)
	}
	return pw
}

// New_PROTO_NC_KQ_LIST_ADD_ACK añade una o más Kingdom Quests al cliente (OpCode 0x581D).
// Struct oficial: PROTO_NC_KQ_LIST_ADD_ACK (2 bytes: NumOfAdd[2] + N * 141 PROTO_KQ_INFO_CLIENT)
func New_PROTO_NC_KQ_LIST_ADD_ACK(quests []ProtoKQInfoClient) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_LIST_ADD_ACK)
	pw.WriteUint16(uint16(len(quests)))
	for _, q := range quests {
		writeProtoKQInfoClient(pw, q)
	}
	return pw
}

// KQUpdateItem modela KQ_UPDATE_ITEMS (7 bytes: Handle[4], Status[1], NumOfJoiner[2]).
type KQUpdateItem struct {
	Handle      uint32
	Status      byte
	NumOfJoiner uint16
}

// New_PROTO_NC_KQ_LIST_UPDATE_ACK actualiza el estado de las KQs activas (OpCode 0x581F).
func New_PROTO_NC_KQ_LIST_UPDATE_ACK(updates []KQUpdateItem) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_LIST_UPDATE_ACK)
	pw.WriteUint16(uint16(len(updates)))
	for _, u := range updates {
		pw.WriteUint32(u.Handle)
		pw.WriteByte(u.Status)
		pw.WriteUint16(u.NumOfJoiner)
	}
	return pw
}

// New_PROTO_NC_KQ_LIST_DELETE_ACK elimina una o más KQs del cliente (OpCode 0x581E).
func New_PROTO_NC_KQ_LIST_DELETE_ACK(deleteHandles []uint32) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_LIST_DELETE_ACK)
	pw.WriteUint16(uint16(len(deleteHandles)))
	for _, h := range deleteHandles {
		pw.WriteUint32(h)
	}
	return pw
}

// New_PROTO_NC_KQ_STATUS_ACK responde el estado y la lista de jugadores inscritos (OpCode 0x5804).
// Struct oficial: PROTO_NC_KQ_STATUS_ACK (7 bytes: Handle[4], Status[1], NumOfJoiner[2] + N * Name5[20])
func New_PROTO_NC_KQ_STATUS_ACK(handle uint32, status byte, joiners []string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_STATUS_ACK)
	pw.WriteUint32(handle)
	pw.WriteByte(status)
	pw.WriteUint16(uint16(len(joiners)))
	for _, name := range joiners {
		pw.WriteString(name, 20)
	}
	return pw
}

// New_PROTO_NC_KQ_JOIN_ACK responde al registro del jugador en una Kingdom Quest (OpCode 0x5806).
// Struct oficial: PROTO_NC_KQ_JOIN_ACK (6 bytes: Handle[4], Error[2])
func New_PROTO_NC_KQ_JOIN_ACK(handle uint32, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_JOIN_ACK)
	pw.WriteUint32(handle)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_KQ_JOIN_CANCEL_ACK responde a la cancelación del registro (OpCode 0x5808).
// Struct oficial: PROTO_NC_KQ_JOIN_CANCEL_ACK (6 bytes: Handle[4], Error[2])
func New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle uint32, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_JOIN_CANCEL_ACK)
	pw.WriteUint32(handle)
	pw.WriteUint16(errCode)
	return pw
}

// KQJoinCharInfo modela KQ_JOIN_CHAR_INFO (23 bytes).
type KQJoinCharInfo struct {
	Level    byte   // 1 byte
	Class    byte   // 1 byte
	CharName string // Name5 (20 bytes)
	TeamType byte   // 1 byte
}

// New_PROTO_NC_KQ_JOIN_LIST_ACK entrega la lista detallada de participantes (OpCode 0x5832).
// Struct oficial: PROTO_NC_KQ_JOIN_LIST_ACK (3 bytes: Error[2], Num[1] + N * 23 bytes)
func New_PROTO_NC_KQ_JOIN_LIST_ACK(errCode uint16, chars []KQJoinCharInfo) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_JOIN_LIST_ACK)
	pw.WriteUint16(errCode)
	pw.WriteByte(byte(len(chars)))
	for _, c := range chars {
		pw.WriteByte(c.Level)
		pw.WriteByte(c.Class)
		pw.WriteString(c.CharName, 20)
		pw.WriteByte(c.TeamType)
	}
	return pw
}

// New_PROTO_NC_KQ_TEAM_SELECT_ACK responde al cambio o selección de equipo (OpCode 0x5838).
// Struct oficial: PROTO_NC_KQ_TEAM_SELECT_ACK (3 bytes: Error[2], TeamType[1])
func New_PROTO_NC_KQ_TEAM_SELECT_ACK(errCode uint16, teamType byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_TEAM_SELECT_ACK)
	pw.WriteUint16(errCode)
	pw.WriteByte(teamType)
	return pw
}

// New_PROTO_NC_KQ_TEAM_SELECT_CMD notifica la selección de equipo a la sala (OpCode 0x5839).
// Struct oficial: PROTO_NC_KQ_TEAM_SELECT_CMD (21 bytes: CharID[20], TeamType[1])
func New_PROTO_NC_KQ_TEAM_SELECT_CMD(charName string, teamType byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_TEAM_SELECT_CMD)
	pw.WriteString(charName, 20)
	pw.WriteByte(teamType)
	return pw
}

// New_PROTO_NC_KQ_JOINING_ALARM_CMD emite la alarma de reclutamiento activo a los jugadores (OpCode 0x5824).
// Struct oficial: PROTO_NC_KQ_JOINING_ALARM_CMD (8 bytes info + 1 byte len + N bytes msg)
func New_PROTO_NC_KQ_JOINING_ALARM_CMD(handle uint32, id uint16, minLev, maxLev byte, msg string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_JOINING_ALARM_CMD)
	pw.WriteUint32(handle)
	pw.WriteUint16(id)
	pw.WriteByte(minLev)
	pw.WriteByte(maxLev)
	msgBytes := []byte(msg)
	pw.WriteByte(byte(len(msgBytes)))
	pw.WriteBytes(msgBytes)
	return pw
}

// New_PROTO_NC_KQ_JOINING_ALARM_END_CMD notifica el fin del periodo de reclutamiento (OpCode 0x5825).
// Struct oficial: PROTO_NC_KQ_JOINING_ALARM_END_CMD (6 bytes: Handle[4], ID[2])
func New_PROTO_NC_KQ_JOINING_ALARM_END_CMD(handle uint32, id uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_JOINING_ALARM_END_CMD)
	pw.WriteUint32(handle)
	pw.WriteUint16(id)
	return pw
}

// New_PROTO_NC_KQ_NOTIFY_CMD envía un anuncio oficial sobre la Kingdom Quest (OpCode 0x580B).
func New_PROTO_NC_KQ_NOTIFY_CMD(msg string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_KQ_NOTIFY_CMD)
	msgBytes := []byte(msg)
	pw.WriteByte(byte(len(msgBytes)))
	pw.WriteBytes(msgBytes)
	return pw
}

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_CHAR_KQMAP_CMD (size = 28 bytes)
type PROTO_NC_CHAR_KQMAP_CMD_DATA struct {
	NKQHandle                        uint32               // offset: 0, size: 4 bytes
	SKQMapName                       Name3                // offset: 4, size: 12 bytes
	NKQCoord                         uint64               // offset: 16, size: 8 bytes
	DKQDate                          uint32               // offset: 24, size: 4 bytes
}

// PROTO_NC_KQ_JOINER (size = 5 bytes)
type PROTO_NC_KQ_JOINER struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NTeamType                        uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_KQ_JOINING_ALARM_CMD (size = 9 bytes)
type PROTO_NC_KQ_JOINING_ALARM_CMD_DATA struct {
	KQInfo                           uint64               // offset: 0, size: 8 bytes
	MsgLength                        uint8                // offset: 8, size: 1 bytes
	Msg                              [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_KQ_JOINING_ALARM_END_CMD (size = 6 bytes)
type PROTO_NC_KQ_JOINING_ALARM_END_CMD_DATA struct {
	NHandle                          uint32               // offset: 0, size: 4 bytes
	NID                              uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_KQ_JOINING_ALARM_LIST (size = 2 bytes)
type PROTO_NC_KQ_JOINING_ALARM_LIST_DATA struct {
	NNumOfKQ                         uint16               // offset: 0, size: 2 bytes
	KQList                           [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_KQ_JOIN_ACK (size = 6 bytes)
type PROTO_NC_KQ_JOIN_ACK_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_KQ_JOIN_CANCEL_ACK (size = 6 bytes)
type PROTO_NC_KQ_JOIN_CANCEL_ACK_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_KQ_JOIN_CANCEL_REQ (size = 4 bytes)
type PROTO_NC_KQ_JOIN_CANCEL_REQ_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_JOIN_LIST_ACK (size = 3 bytes)
type PROTO_NC_KQ_JOIN_LIST_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NNum                             uint8                // offset: 2, size: 1 bytes
	CharInfoList                     [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_KQ_JOIN_LIST_REQ (size = 4 bytes)
type PROTO_NC_KQ_JOIN_LIST_REQ_DATA struct {
	NHandle                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_JOIN_REQ (size = 4 bytes)
type PROTO_NC_KQ_JOIN_REQ_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD (size = 52 bytes)
type PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD_DATA struct {
	NCharNum                         uint32               // offset: 0, size: 4 bytes
	SMapName                         [48]byte             // offset: 4, size: 48 bytes
}

// PROTO_NC_KQ_LIST_ACK (size = 50 bytes)
type PROTO_NC_KQ_LIST_ACK_DATA struct {
	ServerTime                       int32                // offset: 0, size: 4 bytes
	Tm_ServerTime                    [36]byte             // offset: 4, size: 36 bytes
	NewStartHandle                   uint32               // offset: 40, size: 4 bytes
	NewEndHandle                     uint32               // offset: 44, size: 4 bytes
	NumOfNewQuest                    uint16               // offset: 48, size: 2 bytes
	NewQuestArray                    [0]byte              // offset: 50, flexible array member
}

// PROTO_NC_KQ_LIST_ADD_ACK (size = 2 bytes)
type PROTO_NC_KQ_LIST_ADD_ACK_DATA struct {
	NumOfAdd                         uint16               // offset: 0, size: 2 bytes
	AddQuestArray                    [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_KQ_LIST_DELETE_ACK (size = 2 bytes)
type PROTO_NC_KQ_LIST_DELETE_ACK_DATA struct {
	NumOfDelete                      uint16               // offset: 0, size: 2 bytes
	DeleteQuestArray                 [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_KQ_LIST_REFRESH_REQ (size = 1 bytes)
type PROTO_NC_KQ_LIST_REFRESH_REQ_DATA struct {}

// PROTO_NC_KQ_LIST_REQ (size = 8 bytes)
type PROTO_NC_KQ_LIST_REQ_DATA struct {
	StartHandle                      uint32               // offset: 0, size: 4 bytes
	EndHandle                        uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_KQ_LIST_TIME_ACK (size = 40 bytes)
type PROTO_NC_KQ_LIST_TIME_ACK_DATA struct {
	ServerTime                       int32                // offset: 0, size: 4 bytes
	Tm_ServerTime                    [36]byte             // offset: 4, size: 36 bytes
}

// PROTO_NC_KQ_LIST_UPDATE_ACK (size = 2 bytes)
type PROTO_NC_KQ_LIST_UPDATE_ACK_DATA struct {
	NumOfUpdate                      uint16               // offset: 0, size: 2 bytes
	UpdateQuestArray                 [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_KQ_NOTIFY_CMD (size = 1 bytes)
type PROTO_NC_KQ_NOTIFY_CMD_DATA struct {
	MsgLength                        uint8                // offset: 0, size: 1 bytes
	Msg                              [0]byte              // offset: 1, flexible array member
}

// PROTO_NC_KQ_PLAYER_DISJOIN_CMD (size = 8 bytes)
type PROTO_NC_KQ_PLAYER_DISJOIN_CMD_DATA struct {
	NHandle                          uint32               // offset: 0, size: 4 bytes
	NCharNum                         uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_KQ_SCHEDULE_ACK (size = 10 bytes)
type PROTO_NC_KQ_SCHEDULE_ACK_DATA struct {
	NewStartHandle                   uint32               // offset: 0, size: 4 bytes
	NewEndHandle                     uint32               // offset: 4, size: 4 bytes
	NumOfNewQuest                    uint16               // offset: 8, size: 2 bytes
	NewQuestArray                    [0]byte              // offset: 10, flexible array member
}

// PROTO_NC_KQ_SCHEDULE_REQ (size = 8 bytes)
type PROTO_NC_KQ_SCHEDULE_REQ_DATA struct {
	StartHandle                      uint32               // offset: 0, size: 4 bytes
	EndHandle                        uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_KQ_STATUS_ACK (size = 7 bytes)
type PROTO_NC_KQ_STATUS_ACK_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
	Status                           uint8                // offset: 4, size: 1 bytes
	NumOfJoiner                      uint16               // offset: 5, size: 2 bytes
	JoinerNameList                   [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_KQ_STATUS_REQ (size = 4 bytes)
type PROTO_NC_KQ_STATUS_REQ_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_TEAM_SELECT_ACK (size = 3 bytes)
type PROTO_NC_KQ_TEAM_SELECT_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NTeamType                        uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_KQ_TEAM_SELECT_CMD (size = 21 bytes)
type PROTO_NC_KQ_TEAM_SELECT_CMD_DATA struct {
	SCharID                          Name5                // offset: 0, size: 20 bytes
	NTeamType                        uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_KQ_TEAM_SELECT_REQ (size = 1 bytes)
type PROTO_NC_KQ_TEAM_SELECT_REQ_DATA struct {
	NTeamType                        uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_KQ_VOTE_BAN_MSG_CMD (size = 4 bytes)
type PROTO_NC_KQ_VOTE_BAN_MSG_CMD_DATA struct {
	NVoteRate                        uint8                // offset: 0, size: 1 bytes
	NYesNum                          uint8                // offset: 1, size: 1 bytes
	NNoNum                           uint8                // offset: 2, size: 1 bytes
	NCancelNum                       uint8                // offset: 3, size: 1 bytes
}

// PROTO_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD (size = 1 bytes)
type PROTO_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD_DATA struct {}

// PROTO_NC_KQ_VOTE_CANCEL_CMD (size = 20 bytes)
type PROTO_NC_KQ_VOTE_CANCEL_CMD_DATA struct {
	STarget                          Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_KQ_VOTE_RESULT_FAIL_CMD (size = 23 bytes)
type PROTO_NC_KQ_VOTE_RESULT_FAIL_CMD_DATA struct {
	STarget                          Name5                // offset: 0, size: 20 bytes
	NYesNum                          uint8                // offset: 20, size: 1 bytes
	NNoNum                           uint8                // offset: 21, size: 1 bytes
	NCancelNum                       uint8                // offset: 22, size: 1 bytes
}

// PROTO_NC_KQ_VOTE_RESULT_SUC_CMD (size = 24 bytes)
type PROTO_NC_KQ_VOTE_RESULT_SUC_CMD_DATA struct {
	STarget                          Name5                // offset: 0, size: 20 bytes
	NVoteRate                        uint8                // offset: 20, size: 1 bytes
	NYesNum                          uint8                // offset: 21, size: 1 bytes
	NNoNum                           uint8                // offset: 22, size: 1 bytes
	NCancelNum                       uint8                // offset: 23, size: 1 bytes
}

// PROTO_NC_KQ_VOTE_START_ACK (size = 2 bytes)
type PROTO_NC_KQ_VOTE_START_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_KQ_VOTE_START_CHECK_ACK (size = 2 bytes)
type PROTO_NC_KQ_VOTE_START_CHECK_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_KQ_VOTE_START_CHECK_REQ (size = 1 bytes)
type PROTO_NC_KQ_VOTE_START_CHECK_REQ_DATA struct {}

// PROTO_NC_KQ_VOTE_START_REQ (size = 22 bytes)
type PROTO_NC_KQ_VOTE_START_REQ_DATA struct {
	STarget                          Name5                // offset: 0, size: 20 bytes
	NVoteType                        uint8                // offset: 20, size: 1 bytes
	NLen                             uint8                // offset: 21, size: 1 bytes
	SContents                        [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_KQ_VOTE_VOTING_ACK (size = 2 bytes)
type PROTO_NC_KQ_VOTE_VOTING_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_KQ_VOTE_VOTING_CMD (size = 78 bytes)
type PROTO_NC_KQ_VOTE_VOTING_CMD_DATA struct {
	SStarter                         Name5                // offset: 0, size: 20 bytes
	STarget                          Name5                // offset: 20, size: 20 bytes
	NVoteType                        uint8                // offset: 40, size: 1 bytes
	TEndTime                         [36]byte             // offset: 41, size: 36 bytes
	NLen                             uint8                // offset: 77, size: 1 bytes
	SContents                        [0]byte              // offset: 78, flexible array member
}

// PROTO_NC_KQ_VOTE_VOTING_REQ (size = 4 bytes)
type PROTO_NC_KQ_VOTE_VOTING_REQ_DATA struct {
	EVoteType                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_W2Z_DESTROY_CMD (size = 4 bytes)
type PROTO_NC_KQ_W2Z_DESTROY_CMD_DATA struct {
	Hquest                           uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_W2Z_MAKE_REQ (size = 377 bytes)
type PROTO_NC_KQ_W2Z_MAKE_REQ_DATA struct {
	KQInfo                           [377]byte            // offset: 0, size: 377 bytes
}

// PROTO_NC_KQ_W2Z_START_CMD (size = 379 bytes)
type PROTO_NC_KQ_W2Z_START_CMD_DATA struct {
	KQInfo                           [377]byte            // offset: 0, size: 377 bytes
	NumOfJoiner                      uint16               // offset: 377, size: 2 bytes
	JoinerArray                      [0]byte              // offset: 379, flexible array member
}

// PROTO_NC_KQ_Z2W_END_CMD (size = 4 bytes)
type PROTO_NC_KQ_Z2W_END_CMD_DATA struct {
	Handle                           uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_KQ_Z2W_MAKE_ACK (size = 6 bytes)
type PROTO_NC_KQ_Z2W_MAKE_ACK_DATA struct {
	Hquest                           uint32               // offset: 0, size: 4 bytes
	Err                              uint16               // offset: 4, size: 2 bytes
}
