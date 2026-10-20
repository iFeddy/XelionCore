package protocol

// PROTO_NC_DATA_REPORT_* contiene las estructuras y constructores oficiales para el
// sistema de denuncias y reportes a GameMasters (docs/Servicios/WorldManager.md).

// New_PROTO_NC_DATA_REPORT_ADD_ACK construye la respuesta de alta de denuncia (0x8008, 28 bytes).
func New_PROTO_NC_DATA_REPORT_ADD_ACK(errCode uint16, hReporter uint16, harmer string, regNumber int32) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_DATA_REPORT_ADD_ACK)
	pw.WriteUint16(errCode)
	pw.WriteUint16(hReporter)
	pw.WriteString(harmer, 20) // Name5 (20 bytes)
	pw.WriteInt32(regNumber)
	return pw
}

// New_PROTO_NC_DATA_REPORT_CANCEL_ACK construye la respuesta de cancelación de denuncia (0x800A, 8 bytes).
func New_PROTO_NC_DATA_REPORT_CANCEL_ACK(regNumber int32, hReporter uint16, errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_DATA_REPORT_CANCEL_ACK)
	pw.WriteInt32(regNumber)
	pw.WriteUint16(hReporter)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_DATA_REPORT_GET_ACK construye la respuesta de consulta de denuncias (0x800C, 137 bytes).
func New_PROTO_NC_DATA_REPORT_GET_ACK(hReporter uint16, errCode uint16, regNumber int32, result uint8, remark string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_DATA_REPORT_GET_ACK)
	pw.WriteUint16(hReporter)
	pw.WriteUint16(errCode)
	pw.WriteInt32(regNumber)
	pw.WriteByte(result)
	pw.WriteString(remark, 128)
	return pw
}
