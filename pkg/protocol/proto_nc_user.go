package protocol

// Constantes oficiales de paquetes de red extraídos directamente de los PDBs oficiales.
// Nomenclatura 100% idéntica a la documentación oficial (docs/Servicios/Login.md).


// PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ (0x0C65 - Dept: 3, Cmd: 101, Size: 64 bytes)
type PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ_DATA struct {
	VersionKey string // sVersionKey[64]
}

// PROTO_NC_USER_US_LOGIN_REQ (0x0C5A - Dept: 3, Cmd: 90, Size: 316 bytes)
type PROTO_NC_USER_US_LOGIN_REQ_DATA struct {
	UserName  string // sUserName[260]
	Password  string // sPassword[36] (32 bytes MD5 hex string + 4 bytes padding)
	SpawnApps string // spawnapps Name5 (20 bytes)
}

// PROTO_NC_USER_LOGIN_ACK_WorldInfo representa cada mundo en PROTO_NC_USER_LOGIN_ACK (18 bytes)
type PROTO_NC_USER_LOGIN_ACK_WorldInfo struct {
	WorldNo     uint8    // worldno (1 byte)
	WorldName   string   // worldname[16] (16 bytes)
	WorldStatus uint8    // worldstatus (1 byte)
}

// PROTO_NC_USER_LOGINFAIL_ACK (0x0C09 - Dept: 3, Cmd: 9, Size: 2 bytes)
type PROTO_NC_USER_LOGINFAIL_ACK_DATA struct {
	Error uint16 // error (2 bytes)
}

// PROTO_NC_USER_WORLDSELECT_REQ (0x0C0B - Dept: 3, Cmd: 11, Size: 1 byte)
type PROTO_NC_USER_WORLDSELECT_REQ_DATA struct {
	WorldNo uint8 // worldno (1 byte)
}

// PROTO_NC_USER_WORLDSELECT_ACK (0x0C0C - Dept: 3, Cmd: 12, Size: 83 bytes)
type PROTO_NC_USER_WORLDSELECT_ACK_DATA struct {
	WorldStatus uint8       // worldstatus (1 byte)
	IP          string      // ip[16] (16 bytes)
	Port        uint16      // port (2 bytes)
	Token       string      // validate_new[32] (64 bytes transfer token)
}

// Funciones constructoras de paquetes oficiales

// New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK construye la respuesta de versión correcta (0x0C67).
func New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK)
	pw.WriteByte(0) // XTrapServerKeyLength = 0 (sin XTrap activo)
	return pw
}

// New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK construye la respuesta de versión incorrecta (0x0C02).
func New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK)
	return pw
}

// New_PROTO_NC_USER_LOGINFAIL_ACK construye el paquete oficial de rechazo de login (0x0C09).
func New_PROTO_NC_USER_LOGINFAIL_ACK(errCode uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_LOGINFAIL_ACK)
	pw.WriteUint16(errCode)
	return pw
}

// New_PROTO_NC_USER_LOGIN_ACK construye el paquete oficial con el listado de mundos (0x0C0A o 0x0C1C).
func New_PROTO_NC_USER_LOGIN_ACK(opCode uint16, worlds []PROTO_NC_USER_LOGIN_ACK_WorldInfo) *PacketWriter {
	pw := NewPacketWriter(opCode)
	const maxWorldSlots = 11
	pw.WriteByte(uint8(maxWorldSlots))

	for _, w := range worlds {
		pw.WriteByte(w.WorldNo)
		pw.WriteString(w.WorldName, 16)
		pw.WriteByte(w.WorldStatus)
	}

	// Completar slots no utilizados hasta 11 (estándar de Fiesta)
	for i := len(worlds); i < maxWorldSlots; i++ {
		pw.WriteByte(uint8(i))
		pw.WriteString("", 16)
		pw.WriteByte(uint8(WorldStatus_Offline))
	}

	return pw
}

// New_PROTO_NC_USER_WORLDSELECT_ACK construye la transferencia al WorldServer (0x0C0C, 83 bytes).
func New_PROTO_NC_USER_WORLDSELECT_ACK(status uint8, ip string, port uint16, token string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_WORLDSELECT_ACK)
	pw.WriteByte(status)
	pw.WriteString(ip, 16)
	pw.WriteUint16(port)

	// validate_new: unsigned short[32] = 64 bytes
	pw.WriteString(token, 32)
	pw.Fill(32, 0)

	return pw
}

// New_PROTO_NC_USER_FILE_CHECK_ACK construye la respuesta de verificación de archivos (0x0C05).
func New_PROTO_NC_USER_FILE_CHECK_ACK(allow bool) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_FILE_CHECK_ACK)
	pw.WriteBool(allow)
	return pw
}
