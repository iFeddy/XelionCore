package protocol

// New_PROTO_NC_MAP_LOGINFAIL_ACK construye la respuesta oficial de fallo al transferir de zona (OpCode 0x1803).
// Struct oficial PDB: PROTO_NC_MAP_LOGINFAIL_ACK (3 bytes)
// - err: unsigned short (uint16)
// - nWrongDataFileIndex: unsigned char (byte)
func New_PROTO_NC_MAP_LOGINFAIL_ACK(errCode uint16, wrongDataFileIndex byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_MAP_LOGINFAIL_ACK)
	pw.WriteUint16(errCode)
	pw.WriteByte(wrongDataFileIndex)
	return pw
}

// New_PROTO_NC_CHAR_INGAME_ACK construye la confirmación oficial de ingreso al mundo (OpCode 0x581C).
// Enviado cuando el avatar está completamente listo para aparecer en el mapa.
func New_PROTO_NC_CHAR_INGAME_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_CHAR_INGAME_ACK)
	return pw
}

// New_PROTO_NC_CHAR_POINTS_ACK envía los puntos de atributo disponibles al cliente (OpCode 0x105B / Dept: 4, Cmd: 91).
func New_PROTO_NC_CHAR_POINTS_ACK(usablePoints byte) *PacketWriter {
	pw := NewPacketWriter(MakeOpCode(4, 91))
	pw.WriteByte(usablePoints)
	return pw
}

// ShineCharStatVar representa la estructura oficial de stat con base y cambio (PDB: SHINE_CHAR_STATVAR, 8 bytes).
type ShineCharStatVar struct {
	Base   uint32
	Change uint32
}

// WriteToPacket serializa SHINE_CHAR_STATVAR en el flujo binario.
func (s ShineCharStatVar) WriteToPacket(pw *PacketWriter) {
	pw.WriteUint32(s.Base)
	pw.WriteUint32(s.Change)
}
