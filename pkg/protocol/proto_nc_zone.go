package protocol

import "strings"

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

// PROTO_NC_CHAR_ZONE_CHARDATA_REQ_DATA representa la cabecera de datos del personaje enviada en MAP_LOGIN_REQ.
// Struct oficial PDB (22 bytes):
// - wldmanhandle: unsigned short (uint16, 2 bytes)
// - charid: Name5 (char[20], 20 bytes)
type PROTO_NC_CHAR_ZONE_CHARDATA_REQ_DATA struct {
	Wldmanhandle uint16
	Charid       string
}

// PROTO_NC_MAP_LOGIN_REQ_DATA representa la petición oficial de entrada a mapa desde el cliente (OpCode 0x1801).
// Struct oficial PDB (1590 bytes):
// - chardata: PROTO_NC_CHAR_ZONE_CHARDATA_REQ (22 bytes)
// - checksum: Name8[49] (1568 bytes = 49 hashes MD5 de 32 bytes c/u)
type PROTO_NC_MAP_LOGIN_REQ_DATA struct {
	Chardata     PROTO_NC_CHAR_ZONE_CHARDATA_REQ_DATA
	Checksum     [49][32]byte
	HasChecksums bool
}

// Read_PROTO_NC_MAP_LOGIN_REQ deserializa la estructura binaria oficial PROTO_NC_MAP_LOGIN_REQ.
func Read_PROTO_NC_MAP_LOGIN_REQ(pr *PacketReader) (*PROTO_NC_MAP_LOGIN_REQ_DATA, error) {
	req := &PROTO_NC_MAP_LOGIN_REQ_DATA{}
	handle, err := pr.ReadUint16()
	if err != nil {
		return nil, err
	}
	req.Chardata.Wldmanhandle = handle

	var nameLen int
	if pr.Remaining() >= 20 {
		nameLen = 20
	} else {
		nameLen = pr.Remaining()
	}

	name, err := pr.ReadString(nameLen)
	if err != nil {
		return nil, err
	}
	req.Chardata.Charid = strings.TrimRight(name, "\x00")

	// Si hay suficientes bytes en el paquete para los 49 hashes (1568 bytes)
	if pr.Remaining() >= 49*32 {
		for i := 0; i < 49; i++ {
			hashBytes, err := pr.ReadBytes(32)
			if err != nil {
				return nil, err
			}
			copy(req.Checksum[i][:], hashBytes)
		}
		req.HasChecksums = true
	}

	return req, nil
}

