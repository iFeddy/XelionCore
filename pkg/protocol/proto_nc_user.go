package protocol

import "fmt"

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
	WorldNo     uint8  // worldno (1 byte)
	WorldName   string // worldname[16] (16 bytes)
	WorldStatus uint8  // worldstatus (1 byte)
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
	WorldStatus uint8  // worldstatus (1 byte)
	IP          string // ip[16] (16 bytes)
	Port        uint16 // port (2 bytes)
	Token       string // validate_new[32] (64 bytes transfer token)
}

// Funciones constructoras de paquetes oficiales

// New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK construye la respuesta de versión correcta (0x0C67).
func New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK)
	pw.WriteByte(0) // XTrapServerKeyLength = 0 (sin XTrap activo)
	return pw
}

// New_PROTO_NC_USER_NA2026_VERSION_CHECK_ACK construye la respuesta de versión para NA2026 (0x0C2D, 3 bytes wire).
func New_PROTO_NC_USER_NA2026_VERSION_CHECK_ACK(bSuccess byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_VERSION_CHECK_ACK)
	pw.WriteByte(bSuccess) // 0x00 = OK
	return pw
}

// New_PROTO_NC_USER_NA2026_AUTH_SUCCESS_ACK construye el ACK inicial de autenticación en NA2026 (0x0C40, 3 bytes).
func New_PROTO_NC_USER_NA2026_AUTH_SUCCESS_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_AUTH_SUCCESS_ACK)
	pw.WriteByte(0) // 0x00 = OK
	return pw
}

// New_PROTO_NC_USER_NA2026_PRE_WORLDS_ACK construye la señal previa a la lista de mundos (0x0C47, 3 bytes).
func New_PROTO_NC_USER_NA2026_PRE_WORLDS_ACK() *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_PRE_WORLDS_ACK)
	pw.WriteByte(0xFF)
	return pw
}

// New_PROTO_NC_USER_NA2026_ACCOUNT_TOKEN_ACK construye el paquete con tokens de sesión en NA2026 (0x0C3E, 66 bytes payload).
func New_PROTO_NC_USER_NA2026_ACCOUNT_TOKEN_ACK(accountID int) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_ACCOUNT_TOKEN_ACK)
	// Formato capturado oficial:
	// String identificador/sesión de 18 bytes + bytes de control y tokens
	strID := fmt.Sprintf("%017d", accountID)
	pw.WriteString(strID, 18)
	pw.WriteBytes([]byte{0x0C, 0x16, 0x00, 0x00, 0x00, 0x78, 0x3D, 0xE0, 0x0C, 0x65, 0xDB, 0xA3, 0x77, 0x12, 0x00, 0x00})
	pw.WriteBytes([]byte{0x00, 0x00, 0x00, 0x27, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	pw.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x01, 0x28, 0x00, 0x00, 0x00, 0xD0, 0x00})
	return pw
}

// New_PROTO_NC_USER_NA2026_SESSION_PONG construye la respuesta de ping de sesión (0x0C35, 10 bytes payload).
func New_PROTO_NC_USER_NA2026_SESSION_PONG(data []byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_SESSION_PONG)
	if len(data) >= 8 {
		pw.WriteBytes(data[:8])
	} else {
		pw.Fill(8, 0)
	}
	return pw
}

// New_PROTO_NC_USER_NA2026_WORLDSELECT_ACK construye la transferencia al WorldServer en NA2026 (0x0C0B, 86 bytes payload).
func New_PROTO_NC_USER_NA2026_WORLDSELECT_ACK(status uint8, ip string, port uint16, token string, worldNo byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_WORLDSELECT_ACK)
	pw.WriteByte(status)
	pw.WriteString(ip, 16)
	pw.WriteUint16(port)

	// Token de 64 bytes
	pw.WriteString(token, 64)
	pw.WriteByte(worldNo) // Byte final en NA2026
	return pw
}

// New_PROTO_NC_USER_NA2026_WORLD_LOGIN_ACK construye el paquete oficial unificado 0x0C0F en respuesta a 0x0C0E.
// Contiene la cabecera de 24 bytes (randomID, charCount, maxSlots, IP[16], Port) y opcionalmente los datos serializados de los personajes de la cuenta.
func New_PROTO_NC_USER_NA2026_WORLD_LOGIN_ACK(randomID uint16, charCount byte, ip string, port uint16, charData ...[]byte) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_WORLD_LOGIN_ACK)
	if randomID == 0 {
		randomID = 0x7B2A
	}
	pw.WriteUint16(randomID)
	pw.WriteByte(charCount)
	pw.WriteByte(0x0C) // MaxSlots = 12
	pw.WriteString(ip, 16)
	pw.WriteUint16(port)
	for _, cd := range charData {
		pw.WriteBytes(cd)
	}
	return pw
}

// New_PROTO_NC_USER_NA2026_BACK_TO_WORLD_ACK construye la respuesta oficial al presionar 'Atrás' en la pantalla de personajes (0x0C23, 36 bytes payload).
func New_PROTO_NC_USER_NA2026_BACK_TO_WORLD_ACK(token string) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_BACK_TO_WORLD_ACK)
	pw.WriteBytes([]byte{0x58, 0x1E}) // Código de estado / handle capturado en oficial
	pw.WriteString(token, 32)         // Token de retorno de 32 bytes para LoginService
	return pw
}

// New_PROTO_NC_USER_NA2026_PRE_CREATE_CHAR_ACK construye la respuesta oficial 0x0C32 al hacer click en 'Crear Personaje' (8 bytes payload).
func New_PROTO_NC_USER_NA2026_PRE_CREATE_CHAR_ACK(randomID uint16) *PacketWriter {
	pw := NewPacketWriter(PROTO_NC_USER_NA2026_PRE_CREATE_CHAR_ACK)
	pw.WriteUint16(randomID)
	pw.WriteUint32(2)
	return pw
}

// New_PROTO_NC_MISC_HEARTBEAT_ACK construye la respuesta oficial de keep-alive heartbeat (0x0805).
func New_PROTO_NC_MISC_HEARTBEAT_ACK() *PacketWriter {
	return NewPacketWriter(PROTO_NC_MISC_HEARTBEAT_ACK)
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

// NA2026OfficialWorlds contiene la plantilla oficial exacta de los 16 mundos recuperada de la captura oficial de NA2026.
var NA2026OfficialWorlds = []PROTO_NC_USER_LOGIN_ACK_WorldInfo{
	{WorldNo: 10, WorldName: "Malephar DE", WorldStatus: 0x06},
	{WorldNo: 1, WorldName: "Apoline DE", WorldStatus: 0x02},
	{WorldNo: 15, WorldName: "Khazul DE", WorldStatus: 0x06},
	{WorldNo: 3, WorldName: "Epith DE", WorldStatus: 0x02},
	{WorldNo: 4, WorldName: "Nature DE", WorldStatus: 0x02},
	{WorldNo: 5, WorldName: "Elga DE", WorldStatus: 0x02},
	{WorldNo: 6, WorldName: "Cypian FR", WorldStatus: 0x06},
	{WorldNo: 7, WorldName: "Jenira DE", WorldStatus: 0x02},
	{WorldNo: 8, WorldName: "Xelion", WorldStatus: 0x06}, // World 8: Enid ES en oficial, Xelion en XelionCore (Activo)
	{WorldNo: 9, WorldName: "Hednis FR", WorldStatus: 0x06},
	{WorldNo: 0, WorldName: "Teva DE", WorldStatus: 0x02},
	{WorldNo: 11, WorldName: "GFT-ITA0", WorldStatus: 0x02},
	{WorldNo: 12, WorldName: "GFT-ENG1", WorldStatus: 0x02},
	{WorldNo: 13, WorldName: "Iyzel DE", WorldStatus: 0x06},
	{WorldNo: 14, WorldName: "Albireo FR", WorldStatus: 0x06},
	{WorldNo: 16, WorldName: "Kebing DE", WorldStatus: 0x0A},
}

// BuildNA2026OfficialWorldList genera la lista de 16 mundos idéntica a la captura oficial de NA2026.
// Inserta y activa los mundos configurados (principalmente World 8 y World 0) con el nombre y estado activos.
func BuildNA2026OfficialWorldList(configured []PROTO_NC_USER_LOGIN_ACK_WorldInfo) []PROTO_NC_USER_LOGIN_ACK_WorldInfo {
	list := make([]PROTO_NC_USER_LOGIN_ACK_WorldInfo, len(NA2026OfficialWorlds))
	copy(list, NA2026OfficialWorlds)

	for _, cw := range configured {
		found := false
		for i := range list {
			if list[i].WorldNo == cw.WorldNo {
				list[i].WorldName = cw.WorldName
				list[i].WorldStatus = cw.WorldStatus
				found = true
				break
			}
		}
		// Si es el mundo principal por defecto (ID 0) y no se configuró explícitamente el World 8,
		// replicamos el nombre y estado activo en World 8 (usado por el cliente NA2026 para servidores ES/Xelion)
		if cw.WorldNo == 0 {
			for i := range list {
				if list[i].WorldNo == 8 {
					list[i].WorldName = cw.WorldName
					list[i].WorldStatus = cw.WorldStatus
					break
				}
			}
		}
		if !found && len(list) < 16 {
			list = append(list, cw)
		}
	}
	return list
}

// New_PROTO_NC_USER_LOGIN_ACK construye el paquete oficial con el listado de mundos (0x0C0A o 0x0C1C).
// Estructura oficial (PDB Login.pdb):
// struct PROTO_NC_USER_LOGIN_ACK { unsigned char numofworld; WorldInfo worldinfo[0]; };
func New_PROTO_NC_USER_LOGIN_ACK(opCode uint16, worlds []PROTO_NC_USER_LOGIN_ACK_WorldInfo) *PacketWriter {
	pw := NewPacketWriter(opCode)
	pw.WriteByte(uint8(len(worlds)))

	// Si es el opcode de NA2026 (0x0C06 o 0x0C19), lleva el prefijo de 4 bytes identificado en la captura oficial
	if opCode == PROTO_NC_USER_NA2026_LOGIN_ACK || opCode == PROTO_NC_USER_NA2026_WORLDSTATUS_ACK {
		pw.WriteBytes([]byte{0x00, 0x27, 0x01, 0x90})
	}

	for _, w := range worlds {
		pw.WriteByte(w.WorldNo)
		pw.WriteString(w.WorldName, 16)
		pw.WriteByte(w.WorldStatus)
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
