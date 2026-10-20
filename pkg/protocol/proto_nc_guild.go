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
