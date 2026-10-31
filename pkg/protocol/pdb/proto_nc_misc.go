package pdb

// ============================================================================
// Estructuras oficiales de Misceláneos, S2S y Eventos
// Departamento 02: NC_MISC
// Extraídas directamente de docs/Servicios/Account.md y AccountLog.md
// Fidelidad 100% C++ PDB (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_MISC_SEED_ACK (OpCode 0x0807, Size: 2 bytes)
type PROTO_NC_MISC_SEED_ACK struct {
	Seed uint16 // Semilla criptográfica inicial
}

// PROTO_NC_MISC_S2SCONNECTION_RDY_SEND (OpCode 0x0801, Size: 3 bytes)
type PROTO_NC_MISC_S2SCONNECTION_RDY_SEND struct {
	Netpacketheader uint16 // 2B
	Dummy           uint8  // 1B
}

// PROTO_NC_MISC_S2SCONNECTION_REQ (OpCode 0x0802, Size: 7 bytes)
type PROTO_NC_MISC_S2SCONNECTION_REQ struct {
	EchoData          uint8  // Dato de eco (1B)
	ServerToID        uint8  // ID de servidor destino (1B)
	ServerFromID      uint8  // ID de servidor origen (1B)
	ServerFromWorldNum uint8  // Número de mundo origen (1B)
	ServerFromZoneNum  uint8  // Número de zona origen (1B)
	Key               uint16 // Clave de autenticación S2S (2B)
}

// PROTO_NC_MISC_S2SCONNECTION_ACK (OpCode 0x0803, Size: 3 bytes)
type PROTO_NC_MISC_S2SCONNECTION_ACK struct {
	EchoData uint8  // Dato de eco recibido (1B)
	Error    uint16 // Código de error S2S (2B)
}

// PROTO_NC_MISC_HEARTBEAT_ACK_SEND (OpCode 0x0805, Size: 3 bytes)
type PROTO_NC_MISC_HEARTBEAT_ACK_SEND struct {
	Netpacketheader uint16 // 2B
	Dummy           uint8  // 1B
}

// PROTO_NC_MISC_USER_COUNT_CMD (OpCode 0x0808, Size: 4 bytes)
type PROTO_NC_MISC_USER_COUNT_CMD struct {
	NumOfUser uint16 // Usuarios concurrentes actuales (2B)
	MaxOfUser uint16 // Capacidad máxima (2B)
}

// PROTO_NC_MISC_CONNECTER_REQ (OpCode 0x0809, Size: 1 byte)
type PROTO_NC_MISC_CONNECTER_REQ struct {
	Dummy uint8 // 1B
}

// PROTO_NC_MISC_CONNECTER_ACK (OpCode 0x080A, Size: 131 bytes)
type PROTO_NC_MISC_CONNECTER_ACK struct {
	Worldno        uint8      // 1B
	Connectercount uint16     // 2B
	NEachClass     [64]uint16 // 128B
}

// PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD (OpCode 0x080B, Size: 2 bytes)
type PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD struct {
	PacketSize uint16 // Tamaño del paquete encapsulado (2B)
}

// PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD (OpCode 0x080C, Size: 2 bytes)
type PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD struct {
	PacketSize uint16 // Tamaño del paquete encapsulado (2B)
}

// PROTO_NC_MISC_WEB_DB_KEY_REQ (OpCode 0x081D, Size: 6 bytes)
type PROTO_NC_MISC_WEB_DB_KEY_REQ struct {
	NWorldSessionHandle uint16 // Handle de sesión en WorldManager (2B)
	NUserNo             uint32 // nUserNo (4B)
}

// PROTO_NC_MISC_WEB_DB_KEY_ACK (Size: 40 bytes)
type PROTO_NC_MISC_WEB_DB_KEY_ACK struct {
	Error               uint16   // Código de retorno (2B)
	NWorldSessionHandle uint16   // Handle de sesión en WorldManager (2B)
	NUserNo             uint32   // nUserNo (4B)
	WebKey              [32]byte // Token de 32 bytes web (32B)
}

// PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS (OpCode 0x081F, Size: 21 bytes)
type PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS struct {
	WorldNo     uint8  // nWorldNo (1B)
	UserNo      uint32 // nUserNo (4B)
	CharNo      uint32 // nCharNo (4B)
	CharClass   uint8  // Clase de personaje (1B)
	CharLevel   uint8  // Nivel alcanzado (1B)
	CharGetExp  uint64 // Experiencia ganada (8B)
	CharPlayMin uint16 // Minutos jugados en la sesión (2B)
}

// PROTO_NC_MISC_EVENT_HIT3_ADD_CASH (OpCode 0x0820, Size: 28 bytes)
type PROTO_NC_MISC_EVENT_HIT3_ADD_CASH struct {
	UserID  [20]byte // Name5 Usuario (20B)
	AddCash uint64   // Cantidad de Cash a acreditar (8B)
}

// PROTO_NC_MISC_CONNECTFROMWHERE_REQ (Size: 10 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_REQ struct {
	Netpacketzoneheader [6]byte // 6B
	Connectip           [4]byte // 4B
}

// PROTO_NC_MISC_CONNECTFROMWHERE_ACK (Size: 7 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_ACK struct {
	Netpacketzoneheader [6]byte // 6B
	Connectfrom         uint8   // 1B
}

// PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ (OpCode 0x082F, Size: 12 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ struct {
	HZone       uint16                             // 2B
	FromZoneReq PROTO_NC_MISC_CONNECTFROMWHERE_REQ // 10B
}

// PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK (Size: 9 bytes)
type PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK struct {
	HZone     uint16                             // 2B
	ToZoneAck PROTO_NC_MISC_CONNECTFROMWHERE_ACK // 7B
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC (OpCode 0x0838, Size: 13 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC struct {
	NUserNo  uint32 // 4B
	NCharNo  uint32 // 4B
	NWorldNo uint8  // 1B
	HZone    uint16 // 2B
	NQuestID uint16 // 2B
}

// PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM (Size: 15 bytes)
type PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM struct {
	NUserNo  uint32 // 4B
	NCharNo  uint32 // 4B
	NWorldNo uint8  // 1B
	NError   uint16 // 2B
	HZone    uint16 // 2B
	NQuestID uint16 // 2B
}

// PROTO_NC_MISC_EVENT_L20_DB_REQ (OpCode 0x082C, Size: 272 bytes)
type PROTO_NC_MISC_EVENT_L20_DB_REQ struct {
	NWMHandle   uint16    // 2B
	SUserID     [256]byte // Name256Byte (256B)
	NUserNo     uint32    // 4B
	NEventNo    uint8     // 1B
	NWorldNo    uint8     // 1B
	NCharNo     uint32    // 4B
	NEventPoint uint32    // 4B
}

// PROTO_NC_MISC_EVENT_L20_DB_ACK (OpCode 0x082D, Size: 13 bytes)
type PROTO_NC_MISC_EVENT_L20_DB_ACK struct {
	NWMHandle uint16 // 2B
	NUserNo   uint32 // 4B
	NEventNo  uint8  // 1B
	NCharNo   uint32 // 4B
	NError    uint16 // 2B
}

// PROTO_NC_MISC_HACK_SCAN_STORE_CMD (Size: 2 bytes)
type PROTO_NC_MISC_HACK_SCAN_STORE_CMD struct {
	NDescSize uint16 // Tamaño de la descripción (2B)
}

// PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD (OpCode 0x0835, Size: 51 bytes)
type PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD struct {
	NUserNo  uint32                            // 4B
	NWorldNo uint8                             // 1B
	NCharNo  uint32                            // 4B
	SIP      [20]byte                          // 20B
	SCharID  [20]byte                          // 20B
	ScanInfo PROTO_NC_MISC_HACK_SCAN_STORE_CMD // 2B
}

// PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD (OpCode 0x0843, Size: 6 bytes)
type PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD struct {
	NUserNo   uint32 // 4B
	NPlayTime uint16 // Minutos jugados acumulados (2B)
}

// PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ (Size: 8 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ struct {
	NCharNo uint32 // 4B
	NUserNo uint32 // 4B
}

// PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK (Size: 8 bytes)
type PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK struct {
	NCharNo uint32 // 4B
	NError  uint16 // 2B
	NLen    uint16 // 2B
}
