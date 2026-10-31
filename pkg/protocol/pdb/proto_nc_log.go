package pdb

// ============================================================================
// Estructuras oficiales de Auditoría de Cuentas (AccountLog)
// Departamento 01: NC_LOG
// Extraídas directamente de docs/Servicios/AccountLog.md
// Fidelidad 100% C++ PDB (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_LOG_USER_LOGIN (OpCode 0x0721, Size: 9 bytes)
type PROTO_NC_LOG_USER_LOGIN struct {
	Usrregnum uint32  // dwUserNo (4B)
	Worldnum  uint8   // nWorldNo (1B)
	IP        [4]byte // IP binaria v4 (4B)
}

// PROTO_NC_LOG_USER_LOGINFAIL (OpCode 0x0720, Size: 278 bytes)
type PROTO_NC_LOG_USER_LOGINFAIL struct {
	Userid       [256]byte // Name256Byte ID intentado (256B)
	Userpassword [16]byte  // Name4 Contraseña/Hash intentado (16B)
	Userip       [4]byte   // IP binaria (4B)
	Error        uint16    // Código oficial de error (2B)
}

// PROTO_NC_LOG_USER_LOGOUT (OpCode 0x0722, Size: 9 bytes)
type PROTO_NC_LOG_USER_LOGOUT struct {
	Usrregnum uint32 // dwUserNo (4B)
	Worldnum  uint8  // nWorldNo (1B)
	Playmin   uint16 // Minutos jugados acumulados (2B)
	Error     uint16 // Código de error / causa de desconexión (2B)
}
