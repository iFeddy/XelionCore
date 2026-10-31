package pdb

// ============================================================================
// Estructuras oficiales de Auditoría de Sistema (Syslog)
// Departamento 50: NC_SYSLOG
// Extraídas directamente de docs/Servicios/AccountLog.md
// Fidelidad 100% C++ PDB (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS (OpCode 0xC801, Size: 40 bytes)
type PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS struct {
	UserNo    uint32   // nUserNo (4B)
	UserIP    [16]byte // IP string (16B)
	Spawnapps [20]byte // Name5 Aplicación origen (20B)
}

// PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE (OpCode 0xC802, Size: 296 bytes)
type PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE struct {
	UserName  [260]byte // Nombre de cuenta intentado (260B)
	UserIP    [16]byte  // IP cliente (16B)
	Spawnapps [20]byte  // Name5 Aplicación origen (20B)
}

// PROTO_NC_SYSLOG_ACCOUNT_LOGOUT (OpCode 0xC803, Size: 8 bytes)
type PROTO_NC_SYSLOG_ACCOUNT_LOGOUT struct {
	UserNo        uint32 // nUserNo (4B)
	LoginKeepTime uint32 // Segundos de sesión activa (4B)
}

// PROTO_NC_SYSLOG_SERVER_CCU (OpCode 0xC811, Size: 5 bytes)
type PROTO_NC_SYSLOG_SERVER_CCU struct {
	UserCount uint32 // Usuarios concurrentes reales (4B)
	WorldNo   uint8  // Número de mundo (1B)
}
