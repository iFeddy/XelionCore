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
