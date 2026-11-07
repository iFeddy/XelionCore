package protocol

import (
	"encoding/binary"
	"time"
)

// Packet builders oficiales para el sistema de Gremios y Academia (Guild & Guild Academy).
// Nomenclatura y offsets 100% idénticos a WorldManager.pdb / WorldManager.md (Dept 29 & 38).

// packLogonInfo empaqueta el bitfield oficial de 4 bytes de SHINE_GUILD_MEMBER_LOGON_INFO:
// bit 0 (len 1): login (1 si online, 0 si offline)
// bits 1..11 (len 11): year
// bits 12..15 (len 4): month (1..12)
// bits 16..20 (len 5): day (1..31)
// bits 21..25 (len 5): hour (0..23)
func packLogonInfo(t time.Time, isOnline bool) uint32 {
	var v uint32
	if isOnline {
		v |= 1
	}
	y := uint32(t.Year()) & 0x7FF
	m := uint32(t.Month()) & 0xF
	d := uint32(t.Day()) & 0x1F
	h := uint32(t.Hour()) & 0x1F

	v |= (y << 1)
	v |= (m << 12)
	v |= (d << 16)
	v |= (h << 21)
	return v
}

// packTm escribe los 36 bytes de struct tm (9 campos int32 de MSVC x86):
// tm_sec, tm_min, tm_hour, tm_mday, tm_mon, tm_year, tm_wday, tm_yday, tm_isdst
func packTm(t time.Time) []byte {
	buf := make([]byte, 36)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(t.Second()))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(t.Minute()))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(t.Hour()))
	binary.LittleEndian.PutUint32(buf[12:16], uint32(t.Day()))
	binary.LittleEndian.PutUint32(buf[16:20], uint32(t.Month()-1))      // 0-indexed
	binary.LittleEndian.PutUint32(buf[20:24], uint32(t.Year()-1900))   // años desde 1900
	binary.LittleEndian.PutUint32(buf[24:28], uint32(t.Weekday()))     // 0 = Domingo
	binary.LittleEndian.PutUint32(buf[28:32], uint32(t.YearDay()-1))   // 0-indexed
	binary.LittleEndian.PutUint32(buf[32:36], 0)                       // isdst
	return buf
}

// ============================================================================
// ESTRUCTURAS DE DATOS OFICIALES (DEPT 29 - GUILD)
// ============================================================================

// GuildListItemEntry modela SHINE_GUILD_LIST_NEW (46 bytes) para la lista pública de gremios.
type GuildListItemEntry struct {
	GuildNo      uint32 // 4 bytes
	GuildName    string // Name4 (16 bytes)
	MasterCharID string // Name5 (20 bytes)
	Type         byte   // 1 byte
	Grade        byte   // 1 byte
	MemberCnt    uint16 // 2 bytes
	MaxMembers   uint16 // 2 bytes
}

// New_PROTO_NC_GUILD_LIST_ACK construye el paquete oficial de listado de gremios (OpCode 0x7404).
func New_PROTO_NC_GUILD_LIST_ACK(errCode uint16, partMask byte, totalCount, listCount uint16, guilds []GuildListItemEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_LIST_ACK)
	pw.WriteUint16(errCode)
	pw.WriteByte(partMask)
	pw.WriteUint16(totalCount)
	pw.WriteUint16(listCount)
	for _, g := range guilds {
		pw.WriteUint32(g.GuildNo)
		pw.WriteString(g.GuildName, 16)
		pw.WriteString(g.MasterCharID, 20)
		pw.WriteByte(g.Type)
		pw.WriteByte(g.Grade)
		pw.WriteUint16(g.MemberCnt)
		pw.WriteUint16(g.MaxMembers)
	}
	return pw
}

// New_PROTO_NC_GUILD_MAKE_ACK construye la respuesta oficial de creación de gremio (OpCode 0x7406).
// Struct oficial: PROTO_NC_GUILD_MAKE_ACK (27 bytes: Error[2], No[4], Name[16], Type[1], dCreateDate[4])
func New_PROTO_NC_GUILD_MAKE_ACK(errCode uint16, guildNo uint32, guildName string, guildType byte, createDate time.Time) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MAKE_ACK)
	pw.WriteUint16(errCode)
	pw.WriteUint32(guildNo)
	pw.WriteString(guildName, 16)
	pw.WriteByte(guildType)
	pw.WriteInt32(int32(createDate.Unix()))
	return pw
}

// New_PROTO_NC_GUILD_DISMISS_ACK responde a la solicitud de disolución de gremio (OpCode 0x7408).
// Struct oficial: PROTO_NC_GUILD_DISMISS_ACK (2 bytes: Error[2])
func New_PROTO_NC_GUILD_DISMISS_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_DISMISS_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_DISMISS_CMD notifica la disolución de un gremio a sus miembros (OpCode 0x742F).
// Struct oficial: PROTO_NC_GUILD_DISMISS_CMD (57 bytes: Name[16], nDismissStatus[1], dDismissDate[4], tm_dDismissDate[36])
func New_PROTO_NC_GUILD_DISMISS_CMD(guildName string, dismissStatus byte, dismissDate time.Time) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_DISMISS_CMD)
	pw.WriteString(guildName, 16)
	pw.WriteByte(dismissStatus)
	pw.WriteInt32(int32(dismissDate.Unix()))
	pw.WriteBytes(packTm(dismissDate))
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_INVITE_REQ solicita unirse a un jugador (OpCode 0x7409).
func New_PROTO_NC_GUILD_MEMBER_INVITE_REQ(targetName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_INVITE_REQ)
	pw.WriteString(targetName, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_INVITE_ACK responde al invitador si la invitación fue enviada (OpCode 0x740A).
// Struct oficial: PROTO_NC_GUILD_MEMBER_INVITE_ACK (22 bytes: ToCharID[20], Error[2])
func New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_INVITE_ACK)
	pw.WriteString(targetName, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_JOIN_REQ envía la invitación emergente al jugador destino (OpCode 0x740B).
// Struct oficial: PROTO_NC_GUILD_MEMBER_JOIN_REQ (36 bytes: GuildName[16], FromCharID[20])
func New_PROTO_NC_GUILD_MEMBER_JOIN_REQ(guildName, fromCharID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_JOIN_REQ)
	pw.WriteString(guildName, 16)
	pw.WriteString(fromCharID, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_JOIN_ACK es la respuesta del cliente aceptando o rechazando (OpCode 0x740C).
// Struct oficial: PROTO_NC_GUILD_MEMBER_JOIN_ACK (17 bytes: GuildName[16], JoinAck[1])
func New_PROTO_NC_GUILD_MEMBER_JOIN_ACK(guildName string, joinAck byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_JOIN_ACK)
	pw.WriteString(guildName, 16)
	pw.WriteByte(joinAck)
	return pw
}

// GuildMemberClientEntry modela GUILD_MEMBER_CLIENT (110 bytes).
type GuildMemberClientEntry struct {
	Name      string    // Name5 (20 bytes)
	Grade     byte      // 1 byte (0=Master, 1=SubMaster, 2=Officer, 3=Member, 4=Newbie)
	Flags     uint32    // 4 bytes
	ExpRatio  uint16    // 2 bytes
	Intro     string    // 64 bytes
	LogonTime time.Time // 4 bytes (bitfield)
	IsOnline  bool
	ClassID   byte      // 1 byte
	Level     byte      // 1 byte
	IsParty   byte      // 1 byte
	Map       string    // Name3 (12 bytes)
}

// writeGuildMemberClient escribe exactamente los 110 bytes de GUILD_MEMBER_CLIENT.
func writeGuildMemberClient(pw *PacketWriter, m GuildMemberClientEntry) {
	pw.WriteString(m.Name, 20)
	pw.WriteByte(m.Grade)
	pw.WriteUint32(m.Flags)
	pw.WriteUint16(m.ExpRatio)
	pw.WriteString(m.Intro, 64)
	pw.WriteUint32(packLogonInfo(m.LogonTime, m.IsOnline))
	pw.WriteByte(m.ClassID)
	pw.WriteByte(m.Level)
	pw.WriteByte(m.IsParty)
	pw.WriteString(m.Map, 12)
}

// New_PROTO_NC_GUILD_MEMBER_JOIN_CMD anuncia un nuevo miembro a todo el gremio (OpCode 0x7436).
// Struct oficial: PROTO_NC_GUILD_MEMBER_JOIN_CMD (110 bytes: GUILD_MEMBER_CLIENT Member)
func New_PROTO_NC_GUILD_MEMBER_JOIN_CMD(member GuildMemberClientEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_JOIN_CMD)
	writeGuildMemberClient(pw, member)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK responde a la solicitud voluntaria de salida de gremio (OpCode 0x741D).
func New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LEAVE_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LEAVE_CMD notifica la salida de un miembro a todo el gremio (OpCode 0x7438).
// Struct oficial: PROTO_NC_GUILD_MEMBER_LEAVE_CMD (20 bytes: CharID[20])
func New_PROTO_NC_GUILD_MEMBER_LEAVE_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LEAVE_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_VANISH_ACK responde a la expulsión de un miembro (OpCode 0x740F).
// Struct oficial: PROTO_NC_GUILD_MEMBER_VANISH_ACK (22 bytes: CharID[20], Error[2])
func New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(charID string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_VANISH_ACK)
	pw.WriteString(charID, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_VANISH_CMD notifica la expulsión a todo el gremio (OpCode 0x7437).
// Struct oficial: PROTO_NC_GUILD_MEMBER_VANISH_CMD (20 bytes: CharID[20])
func New_PROTO_NC_GUILD_MEMBER_VANISH_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_VANISH_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_GRADE_ACK responde al cambio de rango (OpCode 0x7417).
// Struct oficial: PROTO_NC_GUILD_MEMBER_GRADE_ACK (23 bytes: CharID[20], MemberGrade[1], Error[2])
func New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(charID string, grade byte, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_GRADE_ACK)
	pw.WriteString(charID, 20)
	pw.WriteByte(grade)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_GRADE_CMD difunde el cambio de rango a todo el gremio (OpCode 0x7439).
// Struct oficial: PROTO_NC_GUILD_MEMBER_GRADE_CMD (21 bytes: CharID[20], Grade[1])
func New_PROTO_NC_GUILD_MEMBER_GRADE_CMD(charID string, grade byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_GRADE_CMD)
	pw.WriteString(charID, 20)
	pw.WriteByte(grade)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LOGON_CMD avisa conexión de miembro (OpCode 0x743D).
func New_PROTO_NC_GUILD_MEMBER_LOGON_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LOGON_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LOGOFF_CMD avisa desconexión de miembro (OpCode 0x743E).
func New_PROTO_NC_GUILD_MEMBER_LOGOFF_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LOGOFF_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_MAP_CMD notifica cambio de mapa (OpCode 0x7441).
func New_PROTO_NC_GUILD_MEMBER_MAP_CMD(charID string, mapName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_MAP_CMD)
	pw.WriteString(charID, 20)
	pw.WriteString(mapName, 12)
	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LEVEL_CMD notifica subida de nivel de miembro (OpCode 0x7440).
func New_PROTO_NC_GUILD_MEMBER_LEVEL_CMD(charID string, level byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LEVEL_CMD)
	pw.WriteString(charID, 20)
	pw.WriteByte(level)
	return pw
}

// New_PROTO_NC_GUILD_NOTIFY_ACK responde al guardado de anuncio del gremio (OpCode 0x7411).
// Struct oficial: PROTO_NC_GUILD_NOTIFY_ACK (2 bytes: Error[2])
func New_PROTO_NC_GUILD_NOTIFY_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_NOTIFY_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_NOTIFY_CMD difunde el anuncio de gremio a los miembros (OpCode 0x742D).
// Struct oficial: PROTO_NC_GUILD_NOTIFY_CMD (62 + N bytes: dNotifyDate[4], tm[36], sNotifyCharID[20], nNotifyLength[2], sNotify[N])
func New_PROTO_NC_GUILD_NOTIFY_CMD(notifyDate time.Time, notifyCharID string, notify string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_NOTIFY_CMD)
	pw.WriteInt32(int32(notifyDate.Unix()))
	pw.WriteBytes(packTm(notifyDate))
	pw.WriteString(notifyCharID, 20)
	nBytes := []byte(notify)
	pw.WriteUint16(uint16(len(nBytes)))
	pw.WriteBytes(nBytes)
	return pw
}

// GuildClientDetails representa la estructura binaria exacta de GUILD_CLIENT (1067 bytes).
type GuildClientDetails struct {
	No            uint32
	Name          string
	Money         uint64
	Type          byte
	Grade         byte
	Fame          uint32
	StoneLevel    uint16
	Exp           uint64
	CreateDate    time.Time
	NumMembers    uint16
	MaxMembers    uint16
	WarStatus     byte
	DismissStatus byte
	DismissDate   time.Time
	Intro         string
	NotifyDate    time.Time
	NotifyCharID  string
	Notify        string
}

// New_PROTO_NC_GUILD_MY_GUILD_INFO_ACK construye la información detallada del gremio del usuario (OpCode 0x7419).
// Struct oficial: PROTO_NC_GUILD_MY_GUILD_INFO_ACK (2 bytes Error + 1067 bytes GUILD_CLIENT = 1069 bytes).
func New_PROTO_NC_GUILD_MY_GUILD_INFO_ACK(errCode uint16, g *GuildClientDetails) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MY_GUILD_INFO_ACK)
	pw.WriteUint16(errCode)
	if g == nil || errCode != 0 {
		return pw
	}

	pw.WriteUint32(g.No)
	pw.WriteString(g.Name, 16)
	pw.WriteUint64(g.Money)
	pw.WriteByte(g.Type)
	pw.WriteByte(g.Grade)
	pw.WriteUint32(g.Fame)
	pw.WriteUint16(g.StoneLevel)
	pw.WriteUint64(g.Exp)
	pw.WriteInt32(int32(g.CreateDate.Unix()))
	pw.WriteBytes(packTm(g.CreateDate))
	pw.WriteUint16(g.NumMembers)
	pw.WriteUint16(g.MaxMembers)
	pw.WriteByte(g.WarStatus)

	// War dates (request, start, end) & tm structs
	zeroTime := time.Unix(0, 0)
	pw.WriteInt32(0) // dWarRequestDate
	pw.WriteInt32(0) // dWarStartDate
	pw.WriteInt32(0) // dWarEndDate
	pw.WriteBytes(packTm(zeroTime))
	pw.WriteBytes(packTm(zeroTime))
	pw.WriteBytes(packTm(zeroTime))

	// War enemy & scores
	pw.WriteUint32(0)          // nWarEnemy
	pw.WriteString("", 16)     // sWarEnemyName
	pw.WriteByte(0)            // nWarEnemyGrade
	pw.WriteBytes(make([]byte, 42)) // MyScore (SHINE_GUILD_SCORE)
	pw.WriteBytes(make([]byte, 42)) // EnemyScore (SHINE_GUILD_SCORE)
	pw.WriteUint32(0)          // nWarWinCount
	pw.WriteUint32(0)          // nWarLoseCount
	pw.WriteUint32(0)          // nWarDrawCount

	// Dismiss info
	pw.WriteByte(g.DismissStatus)
	pw.WriteInt32(int32(g.DismissDate.Unix()))
	pw.WriteBytes(packTm(g.DismissDate))

	// Intro
	pw.WriteString(g.Intro, 128)

	// Notify info
	pw.WriteInt32(int32(g.NotifyDate.Unix()))
	pw.WriteBytes(packTm(g.NotifyDate))
	pw.WriteString(g.NotifyCharID, 20)
	pw.WriteString(g.Notify, 512)

	return pw
}

// New_PROTO_NC_GUILD_MEMBER_LIST_ACK construye la lista de miembros de gremio (OpCode 0x741B).
// Struct oficial: PROTO_NC_GUILD_MEMBER_LIST_ACK (6 bytes: Total[2], Index[2], Count[2] + Count * GUILD_MEMBER_CLIENT[110])
func New_PROTO_NC_GUILD_MEMBER_LIST_ACK(totalCount, memberIndex, count uint16, members []GuildMemberClientEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_MEMBER_LIST_ACK)
	pw.WriteUint16(totalCount)
	pw.WriteUint16(memberIndex)
	pw.WriteUint16(count)
	for _, m := range members {
		writeGuildMemberClient(pw, m)
	}
	return pw
}

// New_PROTO_NC_GUILD_CHAT_ACK responde al envío de mensaje en chat de gremio (OpCode 0x7475).
func New_PROTO_NC_GUILD_CHAT_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_CHAT_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_CHAT_CMD transmite un mensaje del chat de gremio a sus miembros (OpCode 0x7474).
// Struct oficial: PROTO_NC_GUILD_CHAT_CMD (26 + N bytes: nGuildNo[4], talker[20], itemLinkDataCount[1], len[1], content[N])
func New_PROTO_NC_GUILD_CHAT_CMD(guildNo uint32, talker string, msg string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_CHAT_CMD)
	pw.WriteUint32(guildNo)
	pw.WriteString(talker, 20)
	msgBytes := []byte(msg)
	pw.WriteByte(0)                  // itemLinkDataCount
	pw.WriteByte(byte(len(msgBytes))) // len
	pw.WriteBytes(msgBytes)          // content
	return pw
}

// New_PROTO_NC_GUILD_NAME_ACK responde con el nombre de un gremio (OpCode 0x7477).
func New_PROTO_NC_GUILD_NAME_ACK(guildNo uint32, guildName string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_NAME_ACK)
	pw.WriteUint32(guildNo)
	pw.WriteString(guildName, 16)
	return pw
}

// ============================================================================
// ESTRUCTURAS DE DATOS OFICIALES (DEPT 38 - GUILD ACADEMY)
// ============================================================================

// AcademyListItemEntry modela SHINE_GUILD_ACADEMY_LIST (42 bytes).
type AcademyListItemEntry struct {
	GuildName           string // Name4 (16 bytes)
	GuildMasterName     string // Name5 (20 bytes)
	NumOfGuildMembers   uint16 // 2 bytes
	NumOfAcademyMembers uint16 // 2 bytes
	NumOfGraduate       uint16 // 2 bytes
}

// New_PROTO_NC_GUILD_ACADEMY_LIST_ACK construye el directorio de academias (OpCode 0x980C).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_LIST_ACK (11 bytes: Error[2], PartMask[1], Total[4], Count[4] + Array)
func New_PROTO_NC_GUILD_ACADEMY_LIST_ACK(errCode uint16, partMask byte, totalCount, count uint32, list []AcademyListItemEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_LIST_ACK)
	pw.WriteUint16(errCode)
	pw.WriteByte(partMask)
	pw.WriteUint32(totalCount)
	pw.WriteUint32(count)
	for _, item := range list {
		pw.WriteString(item.GuildName, 16)
		pw.WriteString(item.GuildMasterName, 20)
		pw.WriteUint16(item.NumOfGuildMembers)
		pw.WriteUint16(item.NumOfAcademyMembers)
		pw.WriteUint16(item.NumOfGraduate)
	}
	return pw
}

// AcademyClientDetails representa GUILD_ACADEMY_CLIENT (736 bytes).
type AcademyClientDetails struct {
	AcademyMasterName string
	NumAcademyMembers uint16
	MaxAcademyMembers uint16
	AcademyPoint      uint32
	AcademyRank       uint32
	BuffUntilTime     int32
	Intro             string
	NotifyDate        time.Time
	NotifyCharID      string
	Notify            string
}

// New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK construye la información de la academia propia (OpCode 0x9808).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK (2 bytes Error + 736 bytes GUILD_ACADEMY_CLIENT).
func New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK(errCode uint16, a *AcademyClientDetails) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK)
	pw.WriteUint16(errCode)
	if a == nil || errCode != 0 {
		return pw
	}

	pw.WriteString(a.AcademyMasterName, 20)
	pw.WriteUint16(a.NumAcademyMembers)
	pw.WriteUint16(a.MaxAcademyMembers)
	pw.WriteUint32(a.AcademyPoint)
	pw.WriteUint32(a.AcademyRank)
	pw.WriteInt32(a.BuffUntilTime)
	pw.WriteString(a.Intro, 128)
	pw.WriteInt32(int32(a.NotifyDate.Unix()))
	pw.WriteBytes(packTm(a.NotifyDate))
	pw.WriteString(a.NotifyCharID, 20)
	pw.WriteString(a.Notify, 512)

	return pw
}

// AcademyMemberClientEntry modela GUILD_ACADEMY_MEMBER_CLIENT (109 bytes).
type AcademyMemberClientEntry struct {
	Name      string
	Grade     byte
	Intro     string
	LogonTime time.Time
	IsOnline  bool
	ClassID   byte
	Level     byte
	IsParty   byte
	Map       string
	JoinTime  time.Time
	JoinLevel byte
}

func writeAcademyMemberClient(pw *PacketWriter, m AcademyMemberClientEntry) {
	pw.WriteString(m.Name, 20)
	pw.WriteByte(m.Grade)
	pw.WriteString(m.Intro, 64)
	pw.WriteUint32(packLogonInfo(m.LogonTime, m.IsOnline))
	pw.WriteByte(m.ClassID)
	pw.WriteByte(m.Level)
	pw.WriteByte(m.IsParty)
	pw.WriteString(m.Map, 12)
	pw.WriteUint32(packLogonInfo(m.JoinTime, false)) // dJoinInfo bitfield
	pw.WriteByte(m.JoinLevel)
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK construye la lista de miembros de la academia (OpCode 0x980E).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK (6 bytes: Total[2], Index[2], Count[2] + Count * GUILD_ACADEMY_MEMBER_CLIENT[109])
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK(totalCount, memberIndex, count uint16, members []AcademyMemberClientEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK)
	pw.WriteUint16(totalCount)
	pw.WriteUint16(memberIndex)
	pw.WriteUint16(count)
	for _, m := range members {
		writeAcademyMemberClient(pw, m)
	}
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK responde al ingreso de academia (OpCode 0x9812).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK (18 bytes: GuildName[16], Error[2])
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK)
	pw.WriteString(guildName, 16)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD difunde el nuevo alumno a la academia (OpCode 0x9813).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD (109 bytes: GUILD_ACADEMY_MEMBER_CLIENT Member)
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD(member AcademyMemberClientEntry) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD)
	writeAcademyMemberClient(pw, member)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK responde a la salida voluntaria de academia (OpCode 0x981C).
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD notifica la salida de un alumno (OpCode 0x9860).
// Struct oficial: PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD (20 bytes: CharID[20])
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK responde a la expulsión de academia (OpCode 0x9817).
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK(charID string, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK)
	pw.WriteString(charID, 20)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD difunde la expulsión de un alumno (OpCode 0x9818).
func New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD(charID string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD)
	pw.WriteString(charID, 20)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK responde al guardado de anuncio de academia (OpCode 0x9825).
func New_PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD difunde el anuncio de la academia (OpCode 0x9826).
func New_PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD(notifyDate time.Time, notifyCharID, notify string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD)
	pw.WriteInt32(int32(notifyDate.Unix()))
	pw.WriteBytes(packTm(notifyDate))
	pw.WriteString(notifyCharID, 20)
	nBytes := []byte(notify)
	pw.WriteUint16(uint16(len(nBytes)))
	pw.WriteBytes(nBytes)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK responde al chat de academia (OpCode 0x986A).
func New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_CHAT_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_GUILD_ACADEMY_CHAT_CMD difunde un mensaje de chat de academia (OpCode 0x9869).
func New_PROTO_NC_GUILD_ACADEMY_CHAT_CMD(guildNo uint32, talker string, msg string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_GUILD_ACADEMY_CHAT_CMD)
	pw.WriteUint32(guildNo)
	pw.WriteString(talker, 20)
	msgBytes := []byte(msg)
	pw.WriteByte(0)                  // itemLinkDataCount
	pw.WriteByte(byte(len(msgBytes))) // len
	pw.WriteBytes(msgBytes)          // content
	return pw
}

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_AVATAR_GUILD_DATA_ACK (size = 9 bytes)
type PROTO_NC_AVATAR_GUILD_DATA_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	Chrregnum                        uint32               // offset: 4, size: 4 bytes
	Slot                             uint8                // offset: 8, size: 1 bytes
	Guild                            [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_AVATAR_GUILD_DATA_REQ (size = 7 bytes)
type PROTO_NC_AVATAR_GUILD_DATA_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	Slot                             uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_AVATAR_GUILD_DATA_REQ_SEND (size = 10 bytes)
type PROTO_NC_AVATAR_GUILD_DATA_REQ_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Req                              [7]byte              // offset: 3, size: 7 bytes
}

// PROTO_NC_AVATAR_GUILD_MEMBER_DATA_ACK (size = 14 bytes)
type PROTO_NC_AVATAR_GUILD_MEMBER_DATA_ACK_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Slot                             uint8                // offset: 4, size: 1 bytes
	NNo                              uint32               // offset: 5, size: 4 bytes
	Error                            uint16               // offset: 9, size: 2 bytes
	PartMask                         uint8                // offset: 11, size: 1 bytes
	NumOfMember                      uint16               // offset: 12, size: 2 bytes
	MemberList                       [0]byte              // offset: 14, flexible array member
}

// PROTO_NC_AVATAR_GUILD_MEMBER_DATA_REQ (size = 11 bytes)
type PROTO_NC_AVATAR_GUILD_MEMBER_DATA_REQ_DATA struct {
	Netpacketheader                  NETPACKETHEADER      // offset: 0, size: 2 bytes
	Chrregnum                        uint32               // offset: 2, size: 4 bytes
	NNo                              uint32               // offset: 6, size: 4 bytes
	Slot                             uint8                // offset: 10, size: 1 bytes
}

// PROTO_NC_CHAR_GUILD_CMD (size = 4 bytes)
type PROTO_NC_CHAR_GUILD_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	Guild                            [0]byte              // offset: 4, flexible array member
}

// PROTO_NC_CHAR_SOMEONEGUILDACADEMYCHANGE_CMD (size = 6 bytes)
type PROTO_NC_CHAR_SOMEONEGUILDACADEMYCHANGE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Guild                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_SOMEONEGUILDCHANGE_CMD (size = 6 bytes)
type PROTO_NC_CHAR_SOMEONEGUILDCHANGE_CMD_DATA struct {
	Handle                           uint16               // offset: 0, size: 2 bytes
	Guild                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_CHAR_ZONE_GUILD_CMD (size = 4 bytes)
type PROTO_NC_CHAR_ZONE_GUILD_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_CHAT_ACK (size = 2 bytes)
type PROTO_NC_GUILD_CHAT_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_CHAT_CMD (size = 26 bytes)
type PROTO_NC_GUILD_CHAT_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	Talker                           Name5                // offset: 4, size: 20 bytes
	Chat                             uint16               // offset: 24, size: 2 bytes
}

// PROTO_NC_GUILD_CHAT_REQ (size = 2 bytes)
type PROTO_NC_GUILD_CHAT_REQ_DATA struct {
	ItemLinkDataCount                uint8                // offset: 0, size: 1 bytes
	Len                              uint8                // offset: 1, size: 1 bytes
	Content                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_DB_ACK (size = 1091 bytes)
type PROTO_NC_GUILD_DB_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	Error                            uint16               // offset: 2, size: 2 bytes
	GuildClient                      [1067]byte           // offset: 4, size: 1067 bytes
	GuildServer                      [20]byte             // offset: 1071, size: 20 bytes
}

// PROTO_NC_GUILD_DB_ALL_ACK (size = 5 bytes)
type PROTO_NC_GUILD_DB_ALL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NumOfGuild                       uint16               // offset: 3, size: 2 bytes
	GuildDB                          [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_DB_ALL_REQ (size = 1 bytes)
type PROTO_NC_GUILD_DB_ALL_REQ_DATA struct {}

// PROTO_NC_GUILD_DB_DELETE_ACK (size = 6 bytes)
type PROTO_NC_GUILD_DB_DELETE_ACK_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_GUILD_DB_DELETE_CANCEL_ACK (size = 8 bytes)
type PROTO_NC_GUILD_DB_DELETE_CANCEL_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	HTool                            uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_GUILD_DB_DELETE_CANCEL_REQ (size = 6 bytes)
type PROTO_NC_GUILD_DB_DELETE_CANCEL_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	HTool                            uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_GUILD_DB_DELETE_REQ (size = 20 bytes)
type PROTO_NC_GUILD_DB_DELETE_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	SGuildName                       [16]byte             // offset: 4, size: 16 bytes
}

// PROTO_NC_GUILD_DB_DISMISS_ACK (size = 13 bytes)
type PROTO_NC_GUILD_DB_DISMISS_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NDismissStatus                   uint8                // offset: 6, size: 1 bytes
	DDismissDate                     int32                // offset: 7, size: 4 bytes
	Error                            uint16               // offset: 11, size: 2 bytes
}

// PROTO_NC_GUILD_DB_DISMISS_REQ (size = 11 bytes)
type PROTO_NC_GUILD_DB_DISMISS_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NDismissStatus                   uint8                // offset: 6, size: 1 bytes
	DDismissDate                     int32                // offset: 7, size: 4 bytes
}

// PROTO_NC_GUILD_DB_INTRO_ACK (size = 12 bytes)
type PROTO_NC_GUILD_DB_INTRO_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	Error                            uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MAKE_ACK (size = 45 bytes)
type PROTO_NC_GUILD_DB_MAKE_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	Make                             [29]byte             // offset: 6, size: 29 bytes
	DCreateDate                      int32                // offset: 35, size: 4 bytes
	NNo                              uint32               // offset: 39, size: 4 bytes
	Error                            uint16               // offset: 43, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MAKE_REQ (size = 55 bytes)
type PROTO_NC_GUILD_DB_MAKE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	SCharID                          Name5                // offset: 6, size: 20 bytes
	Make                             [29]byte             // offset: 26, size: 29 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_ACK (size = 9 bytes)
type PROTO_NC_GUILD_DB_MEMBER_ACK_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	Error                            uint16               // offset: 4, size: 2 bytes
	PartMask                         uint8                // offset: 6, size: 1 bytes
	NumOfMember                      uint16               // offset: 7, size: 2 bytes
	MemberList                       [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_GUILD_DB_MEMBER_GRADE_ACK (size = 37 bytes)
type PROTO_NC_GUILD_DB_MEMBER_GRADE_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
	MemberGrade                      uint8                // offset: 30, size: 1 bytes
	Error                            uint16               // offset: 31, size: 2 bytes
	NSetCharNo                       uint32               // offset: 33, size: 4 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_GRADE_REQ (size = 36 bytes)
type PROTO_NC_GUILD_DB_MEMBER_GRADE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
	MemberGrade                      uint8                // offset: 30, size: 1 bytes
	NSetCharNo                       uint32               // offset: 31, size: 4 bytes
	MemberGrade_Old                  uint8                // offset: 35, size: 1 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_INTRO_ACK (size = 12 bytes)
type PROTO_NC_GUILD_DB_MEMBER_INTRO_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	Error                            uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_JOIN_ACK (size = 13 bytes)
type PROTO_NC_GUILD_DB_MEMBER_JOIN_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	NGrade                           uint8                // offset: 10, size: 1 bytes
	Error                            uint16               // offset: 11, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_JOIN_REQ (size = 31 bytes)
type PROTO_NC_GUILD_DB_MEMBER_JOIN_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	SCharID                          Name5                // offset: 10, size: 20 bytes
	NGrade                           uint8                // offset: 30, size: 1 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_LEAVE_ACK (size = 12 bytes)
type PROTO_NC_GUILD_DB_MEMBER_LEAVE_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	Error                            uint16               // offset: 10, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_LEAVE_REQ (size = 30 bytes)
type PROTO_NC_GUILD_DB_MEMBER_LEAVE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	SCharID                          Name5                // offset: 10, size: 20 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_REQ (size = 4 bytes)
type PROTO_NC_GUILD_DB_MEMBER_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_VANISH_ACK (size = 32 bytes)
type PROTO_NC_GUILD_DB_MEMBER_VANISH_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
	Error                            uint16               // offset: 30, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MEMBER_VANISH_REQ (size = 30 bytes)
type PROTO_NC_GUILD_DB_MEMBER_VANISH_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	SCharID                          Name5                // offset: 2, size: 20 bytes
	NNo                              uint32               // offset: 22, size: 4 bytes
	NCharNo                          uint32               // offset: 26, size: 4 bytes
}

// PROTO_NC_GUILD_DB_MONEY_ADD_ACK (size = 36 bytes)
type PROTO_NC_GUILD_DB_MONEY_ADD_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	NMoneyAdd                        uint64               // offset: 10, size: 8 bytes
	NMemberMoney                     uint64               // offset: 18, size: 8 bytes
	NTotalMoney                      uint64               // offset: 26, size: 8 bytes
	Error                            uint16               // offset: 34, size: 2 bytes
}

// PROTO_NC_GUILD_DB_MONEY_SUB_ACK (size = 36 bytes)
type PROTO_NC_GUILD_DB_MONEY_SUB_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NCharNo                          uint32               // offset: 6, size: 4 bytes
	NMoneySub                        uint64               // offset: 10, size: 8 bytes
	NMemberMoney                     uint64               // offset: 18, size: 8 bytes
	NTotalMoney                      uint64               // offset: 26, size: 8 bytes
	Error                            uint16               // offset: 34, size: 2 bytes
}

// PROTO_NC_GUILD_DB_NOTIFY_ACK (size = 34 bytes)
type PROTO_NC_GUILD_DB_NOTIFY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NNo                              uint32               // offset: 4, size: 4 bytes
	DNotifyDate                      int32                // offset: 8, size: 4 bytes
	SNotifyCharID                    Name5                // offset: 12, size: 20 bytes
	NNotifyLength                    uint16               // offset: 32, size: 2 bytes
	SNotify                          [0]byte              // offset: 34, flexible array member
}

// PROTO_NC_GUILD_DB_NOTIFY_REQ (size = 28 bytes)
type PROTO_NC_GUILD_DB_NOTIFY_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	SNotifyCharID                    Name5                // offset: 6, size: 20 bytes
	NNotifyLength                    uint16               // offset: 26, size: 2 bytes
	SNotify                          [0]byte              // offset: 28, flexible array member
}

// PROTO_NC_GUILD_DB_RENAME_ACK (size = 24 bytes)
type PROTO_NC_GUILD_DB_RENAME_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	SGuildName                       [16]byte             // offset: 6, size: 16 bytes
	Err                              uint16               // offset: 22, size: 2 bytes
}

// PROTO_NC_GUILD_DB_RENAME_REQ (size = 22 bytes)
type PROTO_NC_GUILD_DB_RENAME_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	SGuildName                       [16]byte             // offset: 6, size: 16 bytes
}

// PROTO_NC_GUILD_DB_RESULT_WRITE_CMD (size = 40 bytes)
type PROTO_NC_GUILD_DB_RESULT_WRITE_CMD_DATA struct {
	Guild1                           [20]byte             // offset: 0, size: 20 bytes
	Guild2                           [20]byte             // offset: 20, size: 20 bytes
}

// PROTO_NC_GUILD_DB_RETYPE_ACK (size = 9 bytes)
type PROTO_NC_GUILD_DB_RETYPE_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NType                            uint8                // offset: 6, size: 1 bytes
	Err                              uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_GUILD_DB_RETYPE_REQ (size = 7 bytes)
type PROTO_NC_GUILD_DB_RETYPE_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NType                            uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_GUILD_DB_TOKEN_ALL_ACK (size = 5 bytes)
type PROTO_NC_GUILD_DB_TOKEN_ALL_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NumOfGuild                       uint16               // offset: 3, size: 2 bytes
	GuildTokenDB                     [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_DB_TOKEN_ALL_REQ (size = 1 bytes)
type PROTO_NC_GUILD_DB_TOKEN_ALL_REQ_DATA struct {}

// PROTO_NC_GUILD_DB_WAR_ACK (size = 24 bytes)
type PROTO_NC_GUILD_DB_WAR_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	WarReq                           [22]byte             // offset: 2, size: 22 bytes
}

// PROTO_NC_GUILD_DB_WAR_REQ (size = 22 bytes)
type PROTO_NC_GUILD_DB_WAR_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	NWarEnemy                        uint32               // offset: 6, size: 4 bytes
	DWarRequestDate                  int32                // offset: 10, size: 4 bytes
	DWarStartDate                    int32                // offset: 14, size: 4 bytes
	DWarEndDate                      int32                // offset: 18, size: 4 bytes
}

// PROTO_NC_GUILD_DELETE_CMD (size = 16 bytes)
type PROTO_NC_GUILD_DELETE_CMD_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
}

// PROTO_NC_GUILD_DISMISS_ACK (size = 2 bytes)
type PROTO_NC_GUILD_DISMISS_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_DISMISS_CMD (size = 57 bytes)
type PROTO_NC_GUILD_DISMISS_CMD_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
	NDismissStatus                   uint8                // offset: 16, size: 1 bytes
	DDismissDate                     int32                // offset: 17, size: 4 bytes
	Tm_dDismissDate                  [36]byte             // offset: 21, size: 36 bytes
}

// PROTO_NC_GUILD_DISMISS_REQ (size = 28 bytes)
type PROTO_NC_GUILD_DISMISS_REQ_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
	Password                         Name3                // offset: 16, size: 12 bytes
}

// PROTO_NC_GUILD_EMBLEM_CHECK_AVAILABILITY_ACK (size = 2 bytes)
type PROTO_NC_GUILD_EMBLEM_CHECK_AVAILABILITY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_EMBLEM_INFO_CMD (size = 3 bytes)
type PROTO_NC_GUILD_EMBLEM_INFO_CMD_DATA struct {
	EmblemInfo                       uint16               // offset: 0, size: 2 bytes
	EmblemState                      uint8                // offset: 2, size: 1 bytes
}

// PROTO_NC_GUILD_EMBLEM_INFO_DB_ACK (size = 3 bytes)
type PROTO_NC_GUILD_EMBLEM_INFO_DB_ACK_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	NNumOfEmblem                     uint16               // offset: 1, size: 2 bytes
	EmblemInfo                       [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GUILD_EMBLEM_INFO_DB_REQ (size = 1 bytes)
type PROTO_NC_GUILD_EMBLEM_INFO_DB_REQ_DATA struct {}

// PROTO_NC_GUILD_EMBLEM_INFO_NOTICE_CMD (size = 6 bytes)
type PROTO_NC_GUILD_EMBLEM_INFO_NOTICE_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	EmblemInfo                       uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_GUILD_EMBLEM_LEVELUP_CMD (size = 4 bytes)
type PROTO_NC_GUILD_EMBLEM_LEVELUP_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_EMBLEM_OFF_MSG_CMD (size = 2 bytes)
type PROTO_NC_GUILD_EMBLEM_OFF_MSG_CMD_DATA struct {
	NMonth                           uint8                // offset: 0, size: 1 bytes
	NDay                             uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_GUILD_EMBLEM_SAVE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_EMBLEM_SAVE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_EMBLEM_SAVE_DB_ACK (size = 10 bytes)
type PROTO_NC_GUILD_EMBLEM_SAVE_DB_ACK_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	EmblemInfo                       uint16               // offset: 4, size: 2 bytes
	EmblemCreateDate                 int32                // offset: 6, size: 4 bytes
}

// PROTO_NC_GUILD_EMBLEM_SAVE_DB_REQ (size = 14 bytes)
type PROTO_NC_GUILD_EMBLEM_SAVE_DB_REQ_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
	NGuildNo                         uint32               // offset: 4, size: 4 bytes
	EmblemInfo                       uint16               // offset: 8, size: 2 bytes
	EmblemCreateDate                 int32                // offset: 10, size: 4 bytes
}

// PROTO_NC_GUILD_EMBLEM_SAVE_REQ (size = 2 bytes)
type PROTO_NC_GUILD_EMBLEM_SAVE_REQ_DATA struct {
	EmblemInfo                       uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_EMBLEM_STATE_DB_ACK (size = 9 bytes)
type PROTO_NC_GUILD_EMBLEM_STATE_DB_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	EmblemDeleteDate                 int32                // offset: 4, size: 4 bytes
	EmblemState                      uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GUILD_EMBLEM_STATE_DB_REQ (size = 9 bytes)
type PROTO_NC_GUILD_EMBLEM_STATE_DB_REQ_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	EmblemDeleteDate                 int32                // offset: 4, size: 4 bytes
	EmblemState                      uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GUILD_GRADE_ACK (size = 3 bytes)
type PROTO_NC_GUILD_GRADE_ACK_DATA struct {
	RequestGrade                     uint8                // offset: 0, size: 1 bytes
	Error                            uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_GUILD_GRADE_CMD (size = 1 bytes)
type PROTO_NC_GUILD_GRADE_CMD_DATA struct {
	Grade                            uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_GUILDWARCONFIRM_ACK (size = 16 bytes)
type PROTO_NC_GUILD_GUILDWARCONFIRM_ACK_DATA struct {
	GuildwarA                        uint64               // offset: 0, size: 8 bytes
	GuildwarB                        uint64               // offset: 8, size: 8 bytes
}

// PROTO_NC_GUILD_GUILDWARCONFIRM_REQ (size = 8 bytes)
type PROTO_NC_GUILD_GUILDWARCONFIRM_REQ_DATA struct {
	GuildA                           uint32               // offset: 0, size: 4 bytes
	GuildB                           uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_GUILDWARSTATUS_ACK (size = 2 bytes)
type PROTO_NC_GUILD_GUILDWARSTATUS_ACK_DATA struct {
	Warnumber                        uint16               // offset: 0, size: 2 bytes
	War                              [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_GUILDWARSTATUS_REQ (size = 1 bytes)
type PROTO_NC_GUILD_GUILDWARSTATUS_REQ_DATA struct {}

// PROTO_NC_GUILD_HISTORY_DB_LIST_ACK (size = 10 bytes)
type PROTO_NC_GUILD_HISTORY_DB_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	TotalHistoryCount                uint16               // offset: 4, size: 2 bytes
	HistoryCount                     uint16               // offset: 6, size: 2 bytes
	HistoryIndex                     uint16               // offset: 8, size: 2 bytes
	HistoryList                      [0]byte              // offset: 10, flexible array member
}

// PROTO_NC_GUILD_HISTORY_DB_LIST_REQ (size = 6 bytes)
type PROTO_NC_GUILD_HISTORY_DB_LIST_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_HISTORY_LIST_ACK (size = 8 bytes)
type PROTO_NC_GUILD_HISTORY_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	TotalHistoryCount                uint16               // offset: 2, size: 2 bytes
	HistoryCount                     uint16               // offset: 4, size: 2 bytes
	HistoryIndex                     uint16               // offset: 6, size: 2 bytes
	HistoryList                      [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_GUILD_INFO_ACK (size = 1 bytes)
type PROTO_NC_GUILD_INFO_ACK_DATA struct {}

// PROTO_NC_GUILD_INFO_REQ (size = 1 bytes)
type PROTO_NC_GUILD_INFO_REQ_DATA struct {}

// PROTO_NC_GUILD_LIST_ACK (size = 7 bytes)
type PROTO_NC_GUILD_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	TotalNumOfGuild                  uint16               // offset: 3, size: 2 bytes
	NGuildListCnt                    uint16               // offset: 5, size: 2 bytes
	GuildList                        [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_GUILD_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_MAKE_ACK (size = 27 bytes)
type PROTO_NC_GUILD_MAKE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	No                               uint32               // offset: 2, size: 4 bytes
	Name                             [16]byte             // offset: 6, size: 16 bytes
	TypeField                        uint8                // offset: 22, size: 1 bytes
	DCreateDate                      int32                // offset: 23, size: 4 bytes
}

// PROTO_NC_GUILD_MAKE_REQ (size = 29 bytes)
type PROTO_NC_GUILD_MAKE_REQ_DATA struct {
	Name                             [16]byte             // offset: 0, size: 16 bytes
	Password                         Name3                // offset: 16, size: 12 bytes
	TypeField                        uint8                // offset: 28, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_CLASS_CMD (size = 21 bytes)
type PROTO_NC_GUILD_MEMBER_CLASS_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	ClassID                          uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_EXP_RATIO_CMD (size = 22 bytes)
type PROTO_NC_GUILD_MEMBER_EXP_RATIO_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	ExpRatio                         uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_FLAGS_CMD (size = 24 bytes)
type PROTO_NC_GUILD_MEMBER_FLAGS_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Flags                            uint32               // offset: 20, size: 4 bytes
}

// PROTO_NC_GUILD_MEMBER_GRADE_ACK (size = 23 bytes)
type PROTO_NC_GUILD_MEMBER_GRADE_ACK_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	MemberGrade                      uint8                // offset: 20, size: 1 bytes
	Error                            uint16               // offset: 21, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_GRADE_CMD (size = 21 bytes)
type PROTO_NC_GUILD_MEMBER_GRADE_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Grade                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_GRADE_REQ (size = 21 bytes)
type PROTO_NC_GUILD_MEMBER_GRADE_REQ_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	MemberGrade                      uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_INTRO_ACK (size = 2 bytes)
type PROTO_NC_GUILD_MEMBER_INTRO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_INTRO_CMD (size = 22 bytes)
type PROTO_NC_GUILD_MEMBER_INTRO_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	NMemberIntroLength               uint16               // offset: 20, size: 2 bytes
	SMemberIntro                     [0]byte              // offset: 22, flexible array member
}

// PROTO_NC_GUILD_MEMBER_INTRO_REQ (size = 2 bytes)
type PROTO_NC_GUILD_MEMBER_INTRO_REQ_DATA struct {
	NMemberIntroLength               uint16               // offset: 0, size: 2 bytes
	SMemberIntro                     [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_MEMBER_INVITE_ACK (size = 22 bytes)
type PROTO_NC_GUILD_MEMBER_INVITE_ACK_DATA struct {
	ToCharID                         Name5                // offset: 0, size: 20 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_INVITE_REQ (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_INVITE_REQ_DATA struct {
	ToCharID                         Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_JOIN_ACK (size = 17 bytes)
type PROTO_NC_GUILD_MEMBER_JOIN_ACK_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	JoinAck                          uint8                // offset: 16, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_JOIN_CMD (size = 110 bytes)
type PROTO_NC_GUILD_MEMBER_JOIN_CMD_DATA struct {
	Member                           [110]byte            // offset: 0, size: 110 bytes
}

// PROTO_NC_GUILD_MEMBER_JOIN_ERR (size = 18 bytes)
type PROTO_NC_GUILD_MEMBER_JOIN_ERR_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	Error                            uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_JOIN_REQ (size = 36 bytes)
type PROTO_NC_GUILD_MEMBER_JOIN_REQ_DATA struct {
	GuildName                        [16]byte             // offset: 0, size: 16 bytes
	FromCharID                       Name5                // offset: 16, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_LEAVE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_MEMBER_LEAVE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_LEAVE_CMD (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_LEAVE_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_LEAVE_REQ (size = 1 bytes)
type PROTO_NC_GUILD_MEMBER_LEAVE_REQ_DATA struct {}

// PROTO_NC_GUILD_MEMBER_LEVEL_CMD (size = 21 bytes)
type PROTO_NC_GUILD_MEMBER_LEVEL_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Level                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_LIST_ACK (size = 6 bytes)
type PROTO_NC_GUILD_MEMBER_LIST_ACK_DATA struct {
	TotalMemberCount                 uint16               // offset: 0, size: 2 bytes
	MemberIndex                      uint16               // offset: 2, size: 2 bytes
	MemberCount                      uint16               // offset: 4, size: 2 bytes
	MemberList                       [0]byte              // offset: 6, flexible array member
}

// PROTO_NC_GUILD_MEMBER_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_MEMBER_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_MEMBER_LOGOFF_CMD (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_LOGOFF_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_LOGON_CMD (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_LOGON_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_MAP_CMD (size = 32 bytes)
type PROTO_NC_GUILD_MEMBER_MAP_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	SMap                             Name3                // offset: 20, size: 12 bytes
}

// PROTO_NC_GUILD_MEMBER_PARTY_CMD (size = 21 bytes)
type PROTO_NC_GUILD_MEMBER_PARTY_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Party                            uint8                // offset: 20, size: 1 bytes
}

// PROTO_NC_GUILD_MEMBER_VANISH_ACK (size = 22 bytes)
type PROTO_NC_GUILD_MEMBER_VANISH_ACK_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	Error                            uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_MEMBER_VANISH_CMD (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_VANISH_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MEMBER_VANISH_REQ (size = 20 bytes)
type PROTO_NC_GUILD_MEMBER_VANISH_REQ_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_MOBGUILD_CMD (size = 6 bytes)
type PROTO_NC_GUILD_MOBGUILD_CMD_DATA struct {
	Mobhandle                        uint16               // offset: 0, size: 2 bytes
	Mobguild                         uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_MONEY_ADD_ACK (size = 10 bytes)
type PROTO_NC_GUILD_MONEY_ADD_ACK_DATA struct {
	GuildMoneyAdd                    uint64               // offset: 0, size: 8 bytes
	Error                            uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GUILD_MONEY_ADD_CMD (size = 36 bytes)
type PROTO_NC_GUILD_MONEY_ADD_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	GuildMoneyAdd                    uint64               // offset: 20, size: 8 bytes
	GuildTotalMoney                  uint64               // offset: 28, size: 8 bytes
}

// PROTO_NC_GUILD_MONEY_ADD_REQ (size = 8 bytes)
type PROTO_NC_GUILD_MONEY_ADD_REQ_DATA struct {
	GuildMoneyAdd                    uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_GUILD_MONEY_SUB_ACK (size = 10 bytes)
type PROTO_NC_GUILD_MONEY_SUB_ACK_DATA struct {
	GuildMoneySub                    uint64               // offset: 0, size: 8 bytes
	Error                            uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GUILD_MONEY_SUB_CMD (size = 36 bytes)
type PROTO_NC_GUILD_MONEY_SUB_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
	GuildMoneySub                    uint64               // offset: 20, size: 8 bytes
	GuildTotalMoney                  uint64               // offset: 28, size: 8 bytes
}

// PROTO_NC_GUILD_MY_GUILD_INFO_ACK (size = 2 bytes)
type PROTO_NC_GUILD_MY_GUILD_INFO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	GuildInfo                        [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_MY_GUILD_INFO_REQ (size = 1 bytes)
type PROTO_NC_GUILD_MY_GUILD_INFO_REQ_DATA struct {}

// PROTO_NC_GUILD_MY_GUILD_TOKEN_INFO_ACK (size = 10 bytes)
type PROTO_NC_GUILD_MY_GUILD_TOKEN_INFO_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGuildToken                      uint64               // offset: 2, size: 8 bytes
}

// PROTO_NC_GUILD_MY_GUILD_TOKEN_INFO_REQ (size = 1 bytes)
type PROTO_NC_GUILD_MY_GUILD_TOKEN_INFO_REQ_DATA struct {}

// PROTO_NC_GUILD_MY_GUILD_TOURNAMENT_MATCH_TIME_ACK (size = 45 bytes)
type PROTO_NC_GUILD_MY_GUILD_TOURNAMENT_MATCH_TIME_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	GTCondInfo                       [43]byte             // offset: 2, size: 43 bytes
}

// PROTO_NC_GUILD_MY_GUILD_TOURNAMENT_MATCH_TIME_REQ (size = 1 bytes)
type PROTO_NC_GUILD_MY_GUILD_TOURNAMENT_MATCH_TIME_REQ_DATA struct {}

// PROTO_NC_GUILD_NAME_ACK (size = 22 bytes)
type PROTO_NC_GUILD_NAME_ACK_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	Name                             [16]byte             // offset: 4, size: 16 bytes
	EmblemInfo                       uint16               // offset: 20, size: 2 bytes
}

// PROTO_NC_GUILD_NAME_REQ (size = 4 bytes)
type PROTO_NC_GUILD_NAME_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_NOTIFY_ACK (size = 2 bytes)
type PROTO_NC_GUILD_NOTIFY_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_NOTIFY_CMD (size = 62 bytes)
type PROTO_NC_GUILD_NOTIFY_CMD_DATA struct {
	DNotifyDate                      int32                // offset: 0, size: 4 bytes
	Tm_dNotifyDate                   [36]byte             // offset: 4, size: 36 bytes
	SNotifyCharID                    Name5                // offset: 40, size: 20 bytes
	NNotifyLength                    uint16               // offset: 60, size: 2 bytes
	SNotify                          [0]byte              // offset: 62, flexible array member
}

// PROTO_NC_GUILD_NOTIFY_REQ (size = 2 bytes)
type PROTO_NC_GUILD_NOTIFY_REQ_DATA struct {
	NNotifyLength                    uint16               // offset: 0, size: 2 bytes
	SNotify                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_RENAME_CMD (size = 20 bytes)
type PROTO_NC_GUILD_RENAME_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	SGuildName                       [16]byte             // offset: 4, size: 16 bytes
}

// PROTO_NC_GUILD_RENAME_CMD_SEND (size = 23 bytes)
type PROTO_NC_GUILD_RENAME_CMD_SEND struct {
	Size                             uint8                // offset: 0, size: 1 bytes
	Netcmd                           uint16               // offset: 1, size: 2 bytes
	Cmd                              [20]byte             // offset: 3, size: 20 bytes
}

// PROTO_NC_GUILD_RENAME_REQ (size = 17 bytes)
type PROTO_NC_GUILD_RENAME_REQ_DATA struct {
	NItemSlot                        uint8                // offset: 0, size: 1 bytes
	SGuildName                       [16]byte             // offset: 1, size: 16 bytes
}

// PROTO_NC_GUILD_RETYPE_CMD (size = 5 bytes)
type PROTO_NC_GUILD_RETYPE_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NType                            uint8                // offset: 4, size: 1 bytes
}

// PROTO_NC_GUILD_RETYPE_REQ (size = 2 bytes)
type PROTO_NC_GUILD_RETYPE_REQ_DATA struct {
	NItemSlot                        uint8                // offset: 0, size: 1 bytes
	NType                            uint8                // offset: 1, size: 1 bytes
}

// PROTO_NC_GUILD_STORAGEOPEN_ACK (size = 19 bytes)
type PROTO_NC_GUILD_STORAGEOPEN_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	Cen                              uint64               // offset: 2, size: 8 bytes
	NToken                           uint64               // offset: 10, size: 8 bytes
	Itemcounter                      uint8                // offset: 18, size: 1 bytes
	Itemarray                        [0]byte              // offset: 19, flexible array member
}

// PROTO_NC_GUILD_STORAGEWITHDRAW_CMD (size = 33 bytes)
type PROTO_NC_GUILD_STORAGEWITHDRAW_CMD_DATA struct {
	TypeField                        uint8                // offset: 0, size: 1 bytes
	IsBuy                            uint8                // offset: 1, size: 1 bytes
	Charid                           Name5                // offset: 2, size: 20 bytes
	Item                             uint16               // offset: 22, size: 2 bytes
	Slot                             uint8                // offset: 24, size: 1 bytes
	Lot                              [0]byte              // offset: 25, size: 0 bytes
	Cen                              uint64               // offset: 25, size: 8 bytes
}

// PROTO_NC_GUILD_STORAGEWITHDRAW_RNG (size = 44 bytes)
type PROTO_NC_GUILD_STORAGEWITHDRAW_RNG_DATA struct {
	Start                            [7]byte              // offset: 0, size: 7 bytes
	Gldregnum                        uint32               // offset: 7, size: 4 bytes
	TypeField                        uint8                // offset: 11, size: 1 bytes
	Slot                             uint8                // offset: 12, size: 1 bytes
	IsBuy                            uint8                // offset: 13, size: 1 bytes
	Charid                           Name5                // offset: 14, size: 20 bytes
	Item                             uint16               // offset: 34, size: 2 bytes
	Lot                              [0]byte              // offset: 36, size: 0 bytes
	Cen                              uint64               // offset: 36, size: 8 bytes
}

// PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_ACK (size = 9 bytes)
type PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Err                              uint16               // offset: 6, size: 2 bytes
	Grade                            uint8                // offset: 8, size: 1 bytes
}

// PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_REQ (size = 6 bytes)
type PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
}

// PROTO_NC_GUILD_TOURNAMENTSTART_CMD (size = 14 bytes)
type PROTO_NC_GUILD_TOURNAMENTSTART_CMD_DATA struct {
	NGuildA                          uint32               // offset: 0, size: 4 bytes
	NGuildB                          uint32               // offset: 4, size: 4 bytes
	NGuildTourNamentType             uint16               // offset: 8, size: 2 bytes
	NGuildMyRegNumber                uint32               // offset: 10, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENTSTOPMSG_CMD (size = 20 bytes)
type PROTO_NC_GUILD_TOURNAMENTSTOPMSG_CMD_DATA struct {
	WinnerGuild                      uint32               // offset: 0, size: 4 bytes
	LoseGuild                        uint32               // offset: 4, size: 4 bytes
	AGuildRegNumber                  uint32               // offset: 8, size: 4 bytes
	AGuildScore                      uint16               // offset: 12, size: 2 bytes
	BGuildRegNumber                  uint32               // offset: 14, size: 4 bytes
	BGuildScore                      uint16               // offset: 18, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_FINAL_SELECTION_ACK (size = 7 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_FINAL_SELECTION_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	NFSGuildNum                      uint8                // offset: 6, size: 1 bytes
	NFSGuildNo                       [0]byte              // offset: 7, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_DB_FINAL_SELECTION_REQ (size = 7 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_FINAL_SELECTION_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	NMinLevel                        uint8                // offset: 4, size: 1 bytes
	NMinMember                       uint16               // offset: 5, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_GET_ACK (size = 204 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_GET_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	NMatchType                       uint8                // offset: 6, size: 1 bytes
	MatchTime                        [9]uint32            // offset: 7, size: 36 bytes
	TournamentTree                   [155]byte            // offset: 43, size: 155 bytes
	NLastWinnerGuildNo               uint32               // offset: 198, size: 4 bytes
	NGuildTournamentJoinCnt          uint16               // offset: 202, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_GET_REQ (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_GET_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_ACK (size = 14 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_ACK_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	NTeamPosition                    uint16               // offset: 6, size: 2 bytes
	NGuildNo                         uint32               // offset: 8, size: 4 bytes
	Error                            uint16               // offset: 12, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_LIST_ACK (size = 11 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	HWMChar                          uint16               // offset: 3, size: 2 bytes
	NGTNo                            uint32               // offset: 5, size: 4 bytes
	NJoinGuildNum                    uint16               // offset: 9, size: 2 bytes
	JoinGuildList                    [0]byte              // offset: 11, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_LIST_REQ (size = 6 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_LIST_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_ACK (size = 35 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Error                            uint16               // offset: 8, size: 2 bytes
	HWMChar                          uint16               // offset: 10, size: 2 bytes
	NCharNo                          uint32               // offset: 12, size: 4 bytes
	NGTNo                            uint32               // offset: 16, size: 4 bytes
	NGuildNo                         uint32               // offset: 20, size: 4 bytes
	NMinLevel                        uint8                // offset: 24, size: 1 bytes
	NMinMember                       uint16               // offset: 25, size: 2 bytes
	NProposalMoney                   uint64               // offset: 27, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_REQ (size = 33 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	HWMChar                          uint16               // offset: 8, size: 2 bytes
	NCharNo                          uint32               // offset: 10, size: 4 bytes
	NGTNo                            uint32               // offset: 14, size: 4 bytes
	NGuildNo                         uint32               // offset: 18, size: 4 bytes
	NMinLevel                        uint8                // offset: 22, size: 1 bytes
	NMinMember                       uint16               // offset: 23, size: 2 bytes
	NProposalMoney                   uint64               // offset: 25, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_ACK (size = 28 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NAddMoney                        uint32               // offset: 4, size: 4 bytes
	NAddFame                         uint32               // offset: 8, size: 4 bytes
	NAddExp                          uint32               // offset: 12, size: 4 bytes
	NAddToken                        uint32               // offset: 16, size: 4 bytes
	NAddMoneyDivision                uint64               // offset: 20, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ (size = 12 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	OwnerGuildNum                    uint32               // offset: 4, size: 4 bytes
	NItem                            int32                // offset: 8, size: 4 bytes
	ItemOpt                          [0]byte              // offset: 12, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_ACK (size = 6 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_REQ (size = 161 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_REQ_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	TournamentTree                   [155]byte            // offset: 6, size: 155 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_TIME_ACK (size = 6 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_TIME_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_TIME_REQ (size = 42 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_MATCH_TIME_REQ_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	MatchTime                        [9]uint32            // offset: 6, size: 36 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_RESULT_ACK (size = 9 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_RESULT_ACK_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	NTeamPosition                    uint16               // offset: 4, size: 2 bytes
	NResult                          uint8                // offset: 6, size: 1 bytes
	Error                            uint16               // offset: 7, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_RESULT_REQ (size = 11 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_RESULT_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	NTeamPosition                    uint16               // offset: 4, size: 2 bytes
	NResult                          uint8                // offset: 6, size: 1 bytes
	NGuildNo                         uint32               // offset: 7, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_TYPE_ACK (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_TYPE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NGTNo                            uint32               // offset: 2, size: 4 bytes
	NMatchType                       uint8                // offset: 6, size: 1 bytes
	NFinalSelectionGuildNum          uint8                // offset: 7, size: 1 bytes
	FinalSelectionGuild              [0]byte              // offset: 8, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_DB_SET_TYPE_REQ (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DB_SET_TYPE_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	NMatchType                       uint8                // offset: 4, size: 1 bytes
	NMinLevel                        uint8                // offset: 5, size: 1 bytes
	NMinMemberCnt                    uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DICEGAME_BEFORE_END_TIME_MSG_CMD (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DICEGAME_BEFORE_END_TIME_MSG_CMD_DATA struct {
	UnEndTime                        uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DICEGAME_START_CMD (size = 52 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DICEGAME_START_CMD_DATA struct {
	AGuildRegNumber                  uint32               // offset: 0, size: 4 bytes
	AGuildDelegate                   uint16               // offset: 4, size: 2 bytes
	AGuildDelegateName               Name5                // offset: 6, size: 20 bytes
	BGuildRegNumber                  uint32               // offset: 26, size: 4 bytes
	BGuildDelegate                   uint16               // offset: 30, size: 2 bytes
	BGuildDelegateName               Name5                // offset: 32, size: 20 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_CMD (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_CMD_DATA struct {
	GuildRegNumber                   uint32               // offset: 0, size: 4 bytes
	UnDiceGameResult                 uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ_DATA struct {}

// PROTO_NC_GUILD_TOURNAMENT_END_CMD (size = 12 bytes)
type PROTO_NC_GUILD_TOURNAMENT_END_CMD_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NTournamentMapNo                 uint16               // offset: 1, size: 2 bytes
	NStatus                          uint8                // offset: 3, size: 1 bytes
	NWinnerGuildNo                   uint32               // offset: 4, size: 4 bytes
	NLoseGuildNo                     uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ENTER_LIST_DB_GET_ACK (size = 5 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ENTER_LIST_DB_GET_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NEnterGuildNum                   uint16               // offset: 3, size: 2 bytes
	EnterGuildList                   [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_ENTER_LIST_DB_GET_REQ (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ENTER_LIST_DB_GET_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_CMD (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_CMD_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	HPlayer                          uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ_DATA struct {
	HTargetFlag                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_RESULT_CMD (size = 48 bytes)
type PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_RESULT_CMD_DATA struct {
	CIDPlayer                        Name5                // offset: 0, size: 20 bytes
	CIDFlag                          Name5                // offset: 20, size: 20 bytes
	GRegGuildNo                      uint32               // offset: 40, size: 4 bytes
	GRegFlagCaptureGuild             uint32               // offset: 44, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_GOLD_REFUND_ZONE_CMD (size = 3 bytes)
type PROTO_NC_GUILD_TOURNAMENT_GOLD_REFUND_ZONE_CMD_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	NEnterGuildNum                   uint16               // offset: 1, size: 2 bytes
	EnterGuildList                   [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_ITEM_EFFECT_CMD (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ITEM_EFFECT_CMD_DATA struct {
	GTIiemID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ITEM_FLAG_CMD (size = 5 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ITEM_FLAG_CMD_DATA struct {
	FlagScoreMultiple                uint8                // offset: 0, size: 1 bytes
	NGuildRegNo                      uint32               // offset: 1, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ITEM_PICK_CMD (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ITEM_PICK_CMD_DATA struct {
	GTIiemID                         uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD_DATA struct {
	NNum                             uint16               // offset: 0, size: 2 bytes
	Members                          [0]byte              // offset: 2, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_JOIN_ACK (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_JOIN_ACK_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NTeamPosition                    uint16               // offset: 4, size: 2 bytes
	NError                           uint16               // offset: 6, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_JOIN_LIST_ACK (size = 9 bytes)
type PROTO_NC_GUILD_TOURNAMENT_JOIN_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NGTNo                            uint32               // offset: 3, size: 4 bytes
	NJoinGuildNum                    uint16               // offset: 7, size: 2 bytes
	JoinGuildList                    [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_JOIN_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_JOIN_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_TOURNAMENT_JOIN_REQ (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_JOIN_REQ_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_LAST_WINNER_CMD (size = 5 bytes)
type PROTO_NC_GUILD_TOURNAMENT_LAST_WINNER_CMD_DATA struct {
	BisNewWinnerGuild                uint8                // offset: 0, size: 1 bytes
	NLastWinnerGuildNo               uint32               // offset: 1, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_LIST_ACK (size = 1014 bytes)
type PROTO_NC_GUILD_TOURNAMENT_LIST_ACK_DATA struct {
	NError                           uint16               // offset: 0, size: 2 bytes
	Time_Start                       int32                // offset: 2, size: 4 bytes
	Time_Practic                     int32                // offset: 6, size: 4 bytes
	Time_PracticEnd                  int32                // offset: 10, size: 4 bytes
	Time_Match_161                   int32                // offset: 14, size: 4 bytes
	Time_Match_162                   int32                // offset: 18, size: 4 bytes
	Time_Match_8                     int32                // offset: 22, size: 4 bytes
	Time_Match_4                     int32                // offset: 26, size: 4 bytes
	Time_Match_2                     int32                // offset: 30, size: 4 bytes
	Time_Match_End                   int32                // offset: 34, size: 4 bytes
	Tm_Time_Start                    [36]byte             // offset: 38, size: 36 bytes
	Tm_Time_Practic                  [36]byte             // offset: 74, size: 36 bytes
	Tm_Time_PracticEnd               [36]byte             // offset: 110, size: 36 bytes
	Tm_Time_Match_161                [36]byte             // offset: 146, size: 36 bytes
	Tm_Time_Match_162                [36]byte             // offset: 182, size: 36 bytes
	Tm_Time_Match_8                  [36]byte             // offset: 218, size: 36 bytes
	Tm_Time_Match_4                  [36]byte             // offset: 254, size: 36 bytes
	Tm_Time_Match_2                  [36]byte             // offset: 290, size: 36 bytes
	Tm_Time_Match_End                [36]byte             // offset: 326, size: 36 bytes
	NMatchType                       uint8                // offset: 362, size: 1 bytes
	TournamentList                   [651]byte            // offset: 363, size: 651 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_TOURNAMENT_MANAGERUSERMSG_CMD (size = 20 bytes)
type PROTO_NC_GUILD_TOURNAMENT_MANAGERUSERMSG_CMD_DATA struct {
	CharID                           Name5                // offset: 0, size: 20 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_NOTIFY_CMD (size = 21 bytes)
type PROTO_NC_GUILD_TOURNAMENT_NOTIFY_CMD_DATA struct {
	TypeField                        uint32               // offset: 0, size: 4 bytes
	NRanking                         uint8                // offset: 4, size: 1 bytes
	SGuildName                       [16]byte             // offset: 5, size: 16 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK (size = 11 bytes)
type PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NTournamentSeed                  uint8                // offset: 2, size: 1 bytes
	NGuildRed                        uint32               // offset: 3, size: 4 bytes
	NGuildBlue                       uint32               // offset: 7, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_PLAYERDIEMSG_CMD (size = 44 bytes)
type PROTO_NC_GUILD_TOURNAMENT_PLAYERDIEMSG_CMD_DATA struct {
	VictimCharID                     Name5                // offset: 0, size: 20 bytes
	MurderCharID                     Name5                // offset: 20, size: 20 bytes
	VictimGuildGrade                 int32                // offset: 40, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_PLAYERKILLMSG_CMD (size = 44 bytes)
type PROTO_NC_GUILD_TOURNAMENT_PLAYERKILLMSG_CMD_DATA struct {
	VictimCharID                     Name5                // offset: 0, size: 20 bytes
	MurderCharID                     Name5                // offset: 20, size: 20 bytes
	VictimGuildGrade                 int32                // offset: 40, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_CMD (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_CMD_DATA struct {}

// PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ_DATA struct {
	Ack                              uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_REFUND_NOTICE_CMD (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_REFUND_NOTICE_CMD_DATA struct {
	NRefundMoney                     uint64               // offset: 0, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_REFUND_NOTICE_DB_SET_REQ (size = 8 bytes)
type PROTO_NC_GUILD_TOURNAMENT_REFUND_NOTICE_DB_SET_REQ_DATA struct {
	NGTNo                            uint32               // offset: 0, size: 4 bytes
	NNo                              uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_RESULT_CMD (size = 12 bytes)
type PROTO_NC_GUILD_TOURNAMENT_RESULT_CMD_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NTournamentMapNo                 uint16               // offset: 1, size: 2 bytes
	NStatus                          uint8                // offset: 3, size: 1 bytes
	NWinnerGuildNo                   uint32               // offset: 4, size: 4 bytes
	NLoseGuildNo                     uint32               // offset: 8, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_SCORE_CMD (size = 17 bytes)
type PROTO_NC_GUILD_TOURNAMENT_SCORE_CMD_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	RegNumber                        uint32               // offset: 1, size: 4 bytes
	UnAllScore                       uint32               // offset: 5, size: 4 bytes
	UnCaptureFlagCount               uint32               // offset: 9, size: 4 bytes
	UnPlayerNum                      uint32               // offset: 13, size: 4 bytes
	PlayerScore                      [0]byte              // offset: 17, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_SKILLPOINT_CMD (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_SKILLPOINT_CMD_DATA struct {
	WSkillPoint                      uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_STARTMSG_CMD (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_STARTMSG_CMD_DATA struct {
	NStartBeforeTime                 int32                // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_START_CMD (size = 17 bytes)
type PROTO_NC_GUILD_TOURNAMENT_START_CMD_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NTournamentMapNo                 uint16               // offset: 1, size: 2 bytes
	NExploerTimeMin                  uint16               // offset: 3, size: 2 bytes
	NWaitPlayTimeSec                 uint16               // offset: 5, size: 2 bytes
	NPlayTimeMin                     uint16               // offset: 7, size: 2 bytes
	NNo1                             uint32               // offset: 9, size: 4 bytes
	NNo2                             uint32               // offset: 13, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_TYPE_CMD (size = 1 bytes)
type PROTO_NC_GUILD_TOURNAMENT_TYPE_CMD_DATA struct {
	NMatchType                       uint8                // offset: 0, size: 1 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK (size = 4 bytes)
type PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	WUseSkillNumber                  uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_USESKILL_CMD (size = 6 bytes)
type PROTO_NC_GUILD_TOURNAMENT_USESKILL_CMD_DATA struct {
	RegGuildNumber                   uint32               // offset: 0, size: 4 bytes
	WUseSkillNumber                  uint16               // offset: 4, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ (size = 2 bytes)
type PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ_DATA struct {
	WUseSkillNumber                  uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_ENTER_CMD (size = 9 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_ENTER_CMD_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NGuildNo                         uint32               // offset: 1, size: 4 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_OUT_CMD (size = 9 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_OUT_CMD_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NGuildNo                         uint32               // offset: 1, size: 4 bytes
	NCharNo                          uint32               // offset: 5, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_ACK (size = 27 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	HWMChar                          uint16               // offset: 2, size: 2 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NGTNo                            uint32               // offset: 8, size: 4 bytes
	NGuildNo                         uint32               // offset: 12, size: 4 bytes
	NMinLevel                        uint8                // offset: 16, size: 1 bytes
	NMinMember                       uint16               // offset: 17, size: 2 bytes
	NProposalMoney                   uint64               // offset: 19, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_REQ (size = 25 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_REQ_DATA struct {
	HWMChar                          uint16               // offset: 0, size: 2 bytes
	NCharNo                          uint32               // offset: 2, size: 4 bytes
	NGTNo                            uint32               // offset: 6, size: 4 bytes
	NGuildNo                         uint32               // offset: 10, size: 4 bytes
	NMinLevel                        uint8                // offset: 14, size: 1 bytes
	NMinMember                       uint16               // offset: 15, size: 2 bytes
	NProposalMoney                   uint64               // offset: 17, size: 8 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_ACK (size = 9 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	NTournamentSeed                  uint8                // offset: 2, size: 1 bytes
	NGuildNo                         uint32               // offset: 3, size: 4 bytes
	NGuildMemberCnt                  uint16               // offset: 7, size: 2 bytes
	MemberGradeList                  [0]byte              // offset: 9, flexible array member
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_REQ (size = 5 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_REQ_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NGuildNo                         uint32               // offset: 1, size: 4 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_ACK (size = 5 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_ACK_DATA struct {
	Err                              uint16               // offset: 0, size: 2 bytes
	NTournamentSeed                  uint8                // offset: 2, size: 1 bytes
	NPlayerHnd                       uint16               // offset: 3, size: 2 bytes
}

// PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_REQ (size = 3 bytes)
type PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_REQ_DATA struct {
	NTournamentSeed                  uint8                // offset: 0, size: 1 bytes
	NPlayerHnd                       uint16               // offset: 1, size: 2 bytes
}

// PROTO_NC_GUILD_WAR_ABLE_LIST_ACK (size = 5 bytes)
type PROTO_NC_GUILD_WAR_ABLE_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NumOfWarInfo                     uint16               // offset: 3, size: 2 bytes
	WarInfoList                      [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_WAR_ABLE_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_WAR_ABLE_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_WAR_ACK (size = 2 bytes)
type PROTO_NC_GUILD_WAR_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
}

// PROTO_NC_GUILD_WAR_COOLDOWN_DONE_CMD (size = 1 bytes)
type PROTO_NC_GUILD_WAR_COOLDOWN_DONE_CMD_DATA struct {}

// PROTO_NC_GUILD_WAR_END_CMD (size = 34 bytes)
type PROTO_NC_GUILD_WAR_END_CMD_DATA struct {
	EnemyGuildName                   [16]byte             // offset: 0, size: 16 bytes
	RequestNo                        uint32               // offset: 16, size: 4 bytes
	EnemyGuildFameAdd                uint32               // offset: 20, size: 4 bytes
	GuildFameAdd                     uint32               // offset: 24, size: 4 bytes
	GuildTotalFame                   uint32               // offset: 28, size: 4 bytes
	WarResult                        uint8                // offset: 32, size: 1 bytes
	WarStatus                        uint8                // offset: 33, size: 1 bytes
}

// PROTO_NC_GUILD_WAR_LIST_ACK (size = 5 bytes)
type PROTO_NC_GUILD_WAR_LIST_ACK_DATA struct {
	Error                            uint16               // offset: 0, size: 2 bytes
	PartMask                         uint8                // offset: 2, size: 1 bytes
	NumOfWarInfo                     uint16               // offset: 3, size: 2 bytes
	WarInfoList                      [0]byte              // offset: 5, flexible array member
}

// PROTO_NC_GUILD_WAR_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_WAR_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_WAR_REQ (size = 16 bytes)
type PROTO_NC_GUILD_WAR_REQ_DATA struct {
	TargetName                       [16]byte             // offset: 0, size: 16 bytes
}

// PROTO_NC_GUILD_WAR_SCORE_CMD (size = 116 bytes)
type PROTO_NC_GUILD_WAR_SCORE_CMD_DATA struct {
	GuildName1                       [16]byte             // offset: 0, size: 16 bytes
	GuildScore1                      [42]byte             // offset: 16, size: 42 bytes
	GuildName2                       [16]byte             // offset: 58, size: 16 bytes
	GuildScore2                      [42]byte             // offset: 74, size: 42 bytes
}

// PROTO_NC_GUILD_WAR_START_CMD (size = 60 bytes)
type PROTO_NC_GUILD_WAR_START_CMD_DATA struct {
	WarGuildName                     [16]byte             // offset: 0, size: 16 bytes
	RequestNo                        uint32               // offset: 16, size: 4 bytes
	EndDate                          int32                // offset: 20, size: 4 bytes
	Tm_EndDate                       [36]byte             // offset: 24, size: 36 bytes
}

// PROTO_NC_GUILD_WAR_TARGET_CMD (size = 234 bytes)
type PROTO_NC_GUILD_WAR_TARGET_CMD_DATA struct {
	WarInfo                          [234]byte            // offset: 0, size: 234 bytes
}

// PROTO_NC_GUILD_WORLD_RENAME_ACK (size = 25 bytes)
type PROTO_NC_GUILD_WORLD_RENAME_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemSlot                        uint8                // offset: 6, size: 1 bytes
	SGuildName                       [16]byte             // offset: 7, size: 16 bytes
	Err                              uint16               // offset: 23, size: 2 bytes
}

// PROTO_NC_GUILD_WORLD_RENAME_CMD (size = 22 bytes)
type PROTO_NC_GUILD_WORLD_RENAME_CMD_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	SGuildName                       [16]byte             // offset: 6, size: 16 bytes
}

// PROTO_NC_GUILD_WORLD_RENAME_REQ (size = 23 bytes)
type PROTO_NC_GUILD_WORLD_RENAME_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemSlot                        uint8                // offset: 6, size: 1 bytes
	SGuildName                       [16]byte             // offset: 7, size: 16 bytes
}

// PROTO_NC_GUILD_WORLD_RETYPE_ACK (size = 10 bytes)
type PROTO_NC_GUILD_WORLD_RETYPE_ACK_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemSlot                        uint8                // offset: 6, size: 1 bytes
	NType                            uint8                // offset: 7, size: 1 bytes
	Err                              uint16               // offset: 8, size: 2 bytes
}

// PROTO_NC_GUILD_WORLD_RETYPE_CMD (size = 7 bytes)
type PROTO_NC_GUILD_WORLD_RETYPE_CMD_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NType                            uint8                // offset: 6, size: 1 bytes
}

// PROTO_NC_GUILD_WORLD_RETYPE_REQ (size = 8 bytes)
type PROTO_NC_GUILD_WORLD_RETYPE_REQ_DATA struct {
	Netpacketzoneheader              NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	NItemSlot                        uint8                // offset: 6, size: 1 bytes
	NType                            uint8                // offset: 7, size: 1 bytes
}

// PROTO_NC_GUILD_WORLD_SET_GUILD_TOKEN_CMD (size = 16 bytes)
type PROTO_NC_GUILD_WORLD_SET_GUILD_TOKEN_CMD_DATA struct {
	NDecreaseGuildToken              uint64               // offset: 0, size: 8 bytes
	NRestGuildToken                  uint64               // offset: 8, size: 8 bytes
}

// PROTO_NC_GUILD_WORLD_USE_GUILD_TOKEN_CMD (size = 20 bytes)
type PROTO_NC_GUILD_WORLD_USE_GUILD_TOKEN_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NDecreaseGuildToken              uint64               // offset: 4, size: 8 bytes
	NRestGuildToken                  uint64               // offset: 12, size: 8 bytes
}

// PROTO_NC_GUILD_ZONE_DELETE_CMD (size = 4 bytes)
type PROTO_NC_GUILD_ZONE_DELETE_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ZONE_LIST_ACK (size = 3 bytes)
type PROTO_NC_GUILD_ZONE_LIST_ACK_DATA struct {
	PartMask                         uint8                // offset: 0, size: 1 bytes
	NumOfList                        uint16               // offset: 1, size: 2 bytes
	GuildZoneList                    [0]byte              // offset: 3, flexible array member
}

// PROTO_NC_GUILD_ZONE_LIST_REQ (size = 1 bytes)
type PROTO_NC_GUILD_ZONE_LIST_REQ_DATA struct {}

// PROTO_NC_GUILD_ZONE_MAKE_CMD (size = 53 bytes)
type PROTO_NC_GUILD_ZONE_MAKE_CMD_DATA struct {
	NGuildMasterCharNo               uint32               // offset: 0, size: 4 bytes
	GuildZone                        [49]byte             // offset: 4, size: 49 bytes
}

// PROTO_NC_GUILD_ZONE_MEMBER_JOIN_CMD (size = 8 bytes)
type PROTO_NC_GUILD_ZONE_MEMBER_JOIN_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_ZONE_MEMBER_LEAVE_CMD (size = 4 bytes)
type PROTO_NC_GUILD_ZONE_MEMBER_LEAVE_CMD_DATA struct {
	NCharNo                          uint32               // offset: 0, size: 4 bytes
}

// PROTO_NC_GUILD_ZONE_USE_GUILD_TOKEN_CMD (size = 12 bytes)
type PROTO_NC_GUILD_ZONE_USE_GUILD_TOKEN_CMD_DATA struct {
	NGuildNo                         uint32               // offset: 0, size: 4 bytes
	NDecreaseGuildToken              uint64               // offset: 4, size: 8 bytes
}

// PROTO_NC_GUILD_ZONE_WAR_END_CMD (size = 8 bytes)
type PROTO_NC_GUILD_ZONE_WAR_END_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	NEnemyNo                         uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_GUILD_ZONE_WAR_KILL_CMD (size = 4 bytes)
type PROTO_NC_GUILD_ZONE_WAR_KILL_CMD_DATA struct {
	HWMKiller                        uint16               // offset: 0, size: 2 bytes
	HWMTarget                        uint16               // offset: 2, size: 2 bytes
}

// PROTO_NC_GUILD_ZONE_WAR_START_CMD (size = 8 bytes)
type PROTO_NC_GUILD_ZONE_WAR_START_CMD_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	NEnemyNo                         uint32               // offset: 4, size: 4 bytes
}

// PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_ACK (size = 18 bytes)
type PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_ACK_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Cen                              uint64               // offset: 8, size: 8 bytes
	Err                              uint16               // offset: 16, size: 2 bytes
}

// PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_REQ (size = 20 bytes)
type PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Lockindex                        uint16               // offset: 6, size: 2 bytes
	Guildregnum                      uint32               // offset: 8, size: 4 bytes
	Cen                              uint64               // offset: 12, size: 8 bytes
}

// PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALLSUC_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALLSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NPriceGuildToken                 uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALL_REQ (size = 48 bytes)
type PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALL_REQ_DATA struct {
	NNo                              uint32               // offset: 0, size: 4 bytes
	ItemBuy                          [38]byte             // offset: 4, size: 38 bytes
	Lockindex                        uint16               // offset: 42, size: 2 bytes
	NPriceGuildToken                 uint32               // offset: 44, size: 4 bytes
	Attr                             [0]byte              // offset: 48, flexible array member
}

// PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOTSUC_ACK (size = 10 bytes)
type PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOTSUC_ACK_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	NGuildNo                         uint32               // offset: 2, size: 4 bytes
	NPriceGuildToken                 uint32               // offset: 6, size: 4 bytes
}

// PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOT_REQ (size = 24 bytes)
type PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOT_REQ_DATA struct {
	Lockindex                        uint16               // offset: 0, size: 2 bytes
	NNo                              uint32               // offset: 2, size: 4 bytes
	ItemInfo                         [12]byte             // offset: 6, size: 12 bytes
	NItemID                          uint16               // offset: 18, size: 2 bytes
	NPriceGuildToken                 uint32               // offset: 20, size: 4 bytes
}

// PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_ACK (size = 23 bytes)
type PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_ACK_DATA struct {
	Zonepackheader                   NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Cen                              uint64               // offset: 6, size: 8 bytes
	NToken                           uint64               // offset: 14, size: 8 bytes
	Itemcounter                      uint8                // offset: 22, size: 1 bytes
	Itemarray                        [0]byte              // offset: 23, flexible array member
}

// PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_REQ (size = 10 bytes)
type PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_REQ_DATA struct {
	Packheader                       NETPACKETZONEHEADER  // offset: 0, size: 6 bytes
	Guildregnum                      uint32               // offset: 6, size: 4 bytes
}
