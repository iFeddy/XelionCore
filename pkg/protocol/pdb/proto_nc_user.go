package pdb

// ============================================================================
// Estructuras oficiales de Autenticación de Cuentas y Selección de Mundo
// Departamento 03: NC_USER y Stored Procedures USP_USER
// Extraídas directamente de docs/Servicios/Login.md, Account.md y AccountLog.md
// Fidelidad 100% C++ PDB (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ (OpCode 0x0C65, Size: 64 bytes)
type PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ struct {
	SVersionKey [64]byte // Clave SHA / versión enviada por el cliente
}

// PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK (OpCode 0x0C02, Size: 1 byte C++ / empty)
type PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK struct {
	Dummy uint8
}

// PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK (OpCode 0x0C67, Size: 1 byte)
type PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK struct {
	XTrapServerKeyLength uint8 // Longitud de la clave XTrap del servidor (0 si desactivado)
}

// PROTO_NC_USER_XTRAP_REQ (OpCode 0x0C04, Size: 1 byte)
type PROTO_NC_USER_XTRAP_REQ struct {
	XTrapClientKeyLength uint8 // Longitud de la clave XTrap del cliente
}

// PROTO_NC_USER_XTRAP_ACK (OpCode 0x0C05, Size: 1 byte)
type PROTO_NC_USER_XTRAP_ACK struct {
	BSuccess uint8 // 1: autorizar, 0: rechazar
}

// PROTO_NC_USER_LOGIN_REQ (OpCode 0x0C06, Size: 272 bytes)
type PROTO_NC_USER_LOGIN_REQ struct {
	User     [256]byte // Name256Byte Usuario
	Password [16]byte  // Name4 Hash / Contraseña
}

// PROTO_NC_USER_PASSWORD_CHECK_REQ (OpCode 0x0C07, Size: 274 bytes)
type PROTO_NC_USER_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16    // 2 bytes cabecera interna S2S
	User            [256]byte // Name256Byte Usuario
	Password        [16]byte  // Name4 Contraseña
}

// PROTO_NC_USER_PASSWORD_CHECK_ACK (OpCode 0x0C08, Size: 287 bytes)
type PROTO_NC_USER_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16    // 2 bytes cabecera interna S2S
	ID              [256]byte // Name256Byte ID
	Password        [16]byte  // Name4 Contraseña
	Usrregnum       uint32    // nUserNo
	Authid          uint8     // nAuthID
	Block           uint32    // bIsBlock
	Loginable       uint32    // bIsLoginAble
}

// PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND (OpCode 0x0C08, Size: 290 bytes)
type PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND struct {
	Netpacketheader uint16                           // 2B
	NDummy          uint8                            // 1B
	Data            PROTO_NC_USER_PASSWORD_CHECK_ACK // 287B
}

// PROTO_NC_USER_LOGINFAIL_ACK (OpCode 0x0C09, Size: 2 bytes)
type PROTO_NC_USER_LOGINFAIL_ACK struct {
	Error uint16 // Código numérico oficial de error (69..74)
}

// PROTO_NC_USER_LOGIN_ACK (OpCode 0x0C0A, Size: 1 byte cabecera + N * 18 bytes)
type PROTO_NC_USER_LOGIN_ACK struct {
	Numofworld uint8 // Cantidad de mundos transmitidos
}

// PROTO_NC_USER_LOGIN_ACK_WorldInfo (Size: 18 bytes)
type PROTO_NC_USER_LOGIN_ACK_WorldInfo struct {
	Worldno     uint8    // Número de mundo (1..N)
	Worldname   [16]byte // Name4 Nombre del mundo
	Worldstatus uint8    // Estado operativo (0: normal, 1: concurrido, 2: lleno/mantenimiento)
}

// PROTO_NC_USER_WORLDSELECT_REQ (OpCode 0x0C0B, Size: 1 byte)
type PROTO_NC_USER_WORLDSELECT_REQ struct {
	Worldno uint8 // Mundo seleccionado por el cliente
}

// PROTO_NC_USER_WORLDSELECT_ACK (OpCode 0x0C0C, Size: 83 bytes)
type PROTO_NC_USER_WORLDSELECT_ACK struct {
	Worldstatus uint8      // Estado de aceptación del mundo
	IP          [16]byte   // Name4 IP del WorldServer destino
	Port        uint16     // Puerto del WorldServer
	ValidateNew [32]uint16 // Token criptográfico de transferencia (32 * 2 = 64 bytes)
}

// PROTO_NC_USER_WILLLOGIN_REQ (OpCode 0x0C0D, Size: 347 bytes)
type PROTO_NC_USER_WILLLOGIN_REQ struct {
	Netpacketheader  uint16     // 2 bytes cabecera interna S2S
	WillloginWorldno uint8      // Número de mundo
	Usrregnum        uint32     // nUserNo
	Userid           [256]byte  // Name256Byte ID
	ValidateNew      [32]uint16 // Token criptográfico (64 bytes)
	Spawnapps        [20]byte   // spawnapps Name5 (20 bytes)
}

// PROTO_NC_USER_WILLLOGIN_ACK (OpCode 0x0C0E, Size: 348 bytes)
type PROTO_NC_USER_WILLLOGIN_ACK struct {
	WillLoginReq PROTO_NC_USER_WILLLOGIN_REQ // 347 bytes
	Ok           uint8                       // 1 byte resultado
}

// PROTO_NC_USER_WORLD_LOGIN_REQ (OpCode 0x0C0F, Size: 83 bytes)
type PROTO_NC_USER_WORLD_LOGIN_REQ struct {
	UserNo      uint32     // nUserNo (4B)
	WorldNo     uint8      // nWorldNo (1B)
	ValidateNew [32]uint16 // Token criptográfico transfer (64B)
	ClientIP    [14]byte   // IP origen (14B)
}

// PROTO_NC_USER_WORLD_LOGIN_ACK (OpCode 0x0C10, Size: 26 bytes)
type PROTO_NC_USER_WORLD_LOGIN_ACK struct {
	BSuccess uint8    // 1B
	NWorldNo uint8    // 1B
	SZWorld  [20]byte // Name5 (20B)
	DwUserNo uint32   // 4B
}

// PROTO_NC_USER_LOGIN_DB (Size: 9 bytes)
type PROTO_NC_USER_LOGIN_DB struct {
	Usrregnum uint32  // nUserNo (4B)
	Worldnum  uint8   // nWorldNo (1B)
	IP        [4]byte // IP binaria (4B)
}

// PROTO_NC_USER_LOGOUT_DB (Size: 9 bytes)
type PROTO_NC_USER_LOGOUT_DB struct {
	Usrregnum uint32 // nUserNo (4B)
	Worldnum  uint8  // nWorldNo (1B)
	Playmin   uint16 // Minutos de juego (2B)
	Error     uint16 // Código de error (2B)
}

// PROTO_NC_USER_KICKOFFFROMWORLD_CMD (OpCode 0x0C17, Size: 4 bytes)
type PROTO_NC_USER_KICKOFFFROMWORLD_CMD struct {
	Usrregnum uint32 // nUserNo a expulsar
}

// PROTO_NC_USER_CONNECTCUT_CMD (OpCode 0x0C19, Size: 2 bytes)
type PROTO_NC_USER_CONNECTCUT_CMD struct {
	Reason uint16 // Código de razón de corte
}

// PROTO_NC_USER_WORLD_STATUS_REQ (OpCode 0x0C1B, Size: 1 byte)
type PROTO_NC_USER_WORLD_STATUS_REQ struct {
	Dummy uint8
}

// PROTO_NC_USER_WORLD_STATUS_ACK (OpCode 0x0C1C, Size: 1 byte)
type PROTO_NC_USER_WORLD_STATUS_ACK struct {
	NumOfWorld uint8
}

// PROTO_NC_USER_LOGIN_NETMARBLE_REQ (OpCode 0x0C1D, Size: 2 bytes + variable)
type PROTO_NC_USER_LOGIN_NETMARBLE_REQ struct {
	CPCookieLength uint16 // Longitud de la cookie de autenticación
}

// PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ (OpCode 0x0C1E, Size: 278 bytes)
type PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ struct {
	Netpacketheader uint16    // 2B
	UniID           [20]byte  // Name5 (20B)
	UserID          [256]byte // Name256Byte (256B)
}

// PROTO_NC_USER_LOGIN_OUTSPARK_REQ (OpCode 0x0C20, Size: 85 bytes)
type PROTO_NC_USER_LOGIN_OUTSPARK_REQ struct {
	LogonToken [65]byte // Token de sesión Outspark (65B)
	Spawnapps  [20]byte // Name5 (20B)
}

// PROTO_NC_USER_LOGIN_WITH_OTP_REQ (OpCode 0x0C37, Size: 32 bytes)
type PROTO_NC_USER_LOGIN_WITH_OTP_REQ struct {
	SOTP [32]byte // Token One-Time Password de 32 bytes
}

// PROTO_NC_USER_CREATE_OTP_REQ (OpCode 0x0C35, Size: 266 bytes)
type PROTO_NC_USER_CREATE_OTP_REQ struct {
	Netpacketheader uint16    // 2B
	SAccount        [260]byte // Nombre de cuenta (260B)
	NClientHandle   uint16    // Handle de cliente (2B)
	NWorldSession   uint16    // Handle de sesión de mundo (2B)
}

// PROTO_NC_USER_CREATE_OTP_ACK (OpCode 0x0C36, Size: 36 bytes)
type PROTO_NC_USER_CREATE_OTP_ACK struct {
	NClientHandle uint16   // Handle de cliente (2B)
	NError        uint16   // Código de error (2B)
	SOTP          [32]byte // Token OTP generado (32B)
}

// PROTO_NC_USER_WILL_WORLD_SELECT_REQ (OpCode 0x0C33, Size: 6 bytes)
type PROTO_NC_USER_WILL_WORLD_SELECT_REQ struct {
	NHandle uint16 // Handle de sesión de cliente (2B)
	NUserNo uint32 // nUserNo (4B)
}

// PROTO_NC_USER_WILL_WORLD_SELECT_ACK (OpCode 0x0C34, Size: 34 bytes)
type PROTO_NC_USER_WILL_WORLD_SELECT_ACK struct {
	NError uint16   // Código de retorno (2B)
	SOTP   [32]byte // OTP de 32 bytes para volver a LoginService (32B)
}

// PROTO_NC_USER_TEENAGER_REQ (OpCode 0x0C26, Size: 258 bytes)
type PROTO_NC_USER_TEENAGER_REQ struct {
	Netpacketheader uint16    // 2B
	SUserID         [256]byte // Name256Byte (256B)
}

// PROTO_NC_USER_TEENAGER_CMD (Size: 55 bytes)
type PROTO_NC_USER_TEENAGER_CMD struct {
	NUserNo    uint32   // nUserNo (4B)
	NRemainSec uint32   // Segundos restantes permitidos (4B)
	BShutDown  uint8    // Flag de cierre inmediato (1B)
	SReason    [46]byte // Motivo de corte (46B)
}

// PROTO_NC_USER_TEENAGER_ACK (OpCode 0x0C27, Size: 313 bytes)
type PROTO_NC_USER_TEENAGER_ACK struct {
	SUserID [256]byte                  // Name256Byte (256B)
	NHandle uint16                     // Handle cliente (2B)
	Data    PROTO_NC_USER_TEENAGER_CMD // 55B
}

// PROTO_NC_USER_TEENAGER_SET_CMD (OpCode 0x0C28, Size: 54 bytes)
type PROTO_NC_USER_TEENAGER_SET_CMD struct {
	Netpacketheader uint16 // 2B
	NUserNo         uint32 // 4B
	NRemainSec      uint32 // 4B
	NStateCode      uint32 // 4B
	NReserved2      uint32 // 4B
	NReserved3      uint32 // 4B
	NReserved4      uint32 // 4B
	NReserved5      uint32 // 4B
	NReserved6      uint32 // 4B
	NReserved7      uint32 // 4B
	NReserved8      uint32 // 4B
	NReserved9      uint32 // 4B
	NReserved10     uint32 // 4B
	NReserved11     uint32 // 4B
}

// PROTO_NC_USER_IS_IP_BLOCK_REQ (OpCode 0x0C2A, Size: 298 bytes)
type PROTO_NC_USER_IS_IP_BLOCK_REQ struct {
	LoginData     PROTO_NC_USER_LOGIN_REQ // 272B
	NUserNo       uint32                  // 4B
	NClientHandle uint16                  // 2B
	SUserIP       [20]byte                // 20B
}

// PROTO_NC_USER_IS_IP_BLOCK_ACK (OpCode 0x0C2B, Size: 275 bytes)
type PROTO_NC_USER_IS_IP_BLOCK_ACK struct {
	LoginData     PROTO_NC_USER_LOGIN_REQ // 272B
	NClientHandle uint16                  // 2B
	BUserIPBlock  uint8                   // 1B
}

// PROTO_NC_USER_GER_PASSWORD_CHECK_REQ (OpCode 0x0C39, Size: 274 bytes)
type PROTO_NC_USER_GER_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16    // 2B
	User            [256]byte // 256B
	Password        [16]byte  // 16B
}

// PROTO_NC_USER_GER_PASSWORD_CHECK_ACK (OpCode 0x0C3A, Size: 287 bytes)
type PROTO_NC_USER_GER_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16    // 2B
	ID              [256]byte // 256B
	Password        [16]byte  // 16B
	Usrregnum       uint32    // 4B
	Authid          uint8     // 1B
	Block           uint32    // 4B
	Loginable       uint32    // 4B
}

// PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND (OpCode 0x0C3A, Size: 290 bytes)
type PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND struct {
	Netpacketheader uint16                               // 2B
	NDummy          uint8                                // 1B
	Data            PROTO_NC_USER_GER_PASSWORD_CHECK_ACK // 287B
}

// PROTO_NC_USER_GER_IS_IP_BLOCK_REQ (OpCode 0x0C3B, Size: 313 bytes)
type PROTO_NC_USER_GER_IS_IP_BLOCK_REQ struct {
	LoginData     PROTO_NC_USER_GER_PASSWORD_CHECK_ACK // 287B
	NUserNo       uint32                               // 4B
	NClientHandle uint16                               // 2B
	SUserIP       [20]byte                             // 20B
}

// PROTO_NC_USER_GER_IS_IP_BLOCK_ACK (OpCode 0x0C3C, Size: 290 bytes)
type PROTO_NC_USER_GER_IS_IP_BLOCK_ACK struct {
	LoginData     PROTO_NC_USER_GER_PASSWORD_CHECK_ACK // 287B
	NClientHandle uint16                               // 2B
	BUserIPBlock  uint8                                // 1B
}

// PROTO_NC_USER_TW_PASSWORD_CHECK_REQ (OpCode 0x0C3E, Size: 102 bytes)
type PROTO_NC_USER_TW_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16   // 2B
	SUserName       [64]byte // 64B
	SPassword       [36]byte // 36B
}

// PROTO_NC_USER_TW_PASSWORD_CHECK_ACK (OpCode 0x0C3F, Size: 8 bytes)
type PROTO_NC_USER_TW_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16 // 2B
	NError          uint16 // 2B
	NUserNo         uint32 // 4B
}

// PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND (OpCode 0x0C3F, Size: 11 bytes)
type PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND struct {
	Netpacketheader uint16                              // 2B
	NDummy          uint8                               // 1B
	Data            PROTO_NC_USER_TW_PASSWORD_CHECK_ACK // 8B
}

// PROTO_NC_USER_TW_IS_IP_BLOCK_REQ (OpCode 0x0C40, Size: 34 bytes)
type PROTO_NC_USER_TW_IS_IP_BLOCK_REQ struct {
	LoginData     PROTO_NC_USER_TW_PASSWORD_CHECK_ACK // 8B
	NUserNo       uint32                              // 4B
	NClientHandle uint16                              // 2B
	SUserIP       [20]byte                            // 20B
}

// PROTO_NC_USER_TW_IS_IP_BLOCK_ACK (OpCode 0x0C41, Size: 11 bytes)
type PROTO_NC_USER_TW_IS_IP_BLOCK_ACK struct {
	LoginData     PROTO_NC_USER_TW_PASSWORD_CHECK_ACK // 8B
	NClientHandle uint16                              // 2B
	BUserIPBlock  uint8                               // 1B
}

// PROTO_NC_USER_JP_PASSWORD_CHECK_REQ (OpCode 0x0C4B, Size: 336 bytes)
type PROTO_NC_USER_JP_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16    // 2B
	SUserName       [24]byte  // 24B
	SPassword       [36]byte  // 36B
	SChannel        [254]byte // 254B
	SUserIP         [20]byte  // 20B
}

// PROTO_NC_USER_JP_PASSWORD_CHECK_ACK (OpCode 0x0C4C, Size: 306 bytes)
type PROTO_NC_USER_JP_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16    // 2B
	NError          uint16    // 2B
	NUserNo         uint32    // 4B
	SUserName       [24]byte  // 24B
	SChannel        [254]byte // 254B
	SUserIP         [20]byte  // 20B
}

// PROTO_NC_USER_JP_IS_IP_BLOCK_REQ (OpCode 0x0C4D, Size: 332 bytes)
type PROTO_NC_USER_JP_IS_IP_BLOCK_REQ struct {
	LoginData     PROTO_NC_USER_JP_PASSWORD_CHECK_ACK // 306B
	NUserNo       uint32                              // 4B
	NClientHandle uint16                              // 2B
	SUserIP       [20]byte                            // 20B
}

// PROTO_NC_USER_JP_IS_IP_BLOCK_ACK (OpCode 0x0C4E, Size: 309 bytes)
type PROTO_NC_USER_JP_IS_IP_BLOCK_ACK struct {
	LoginData     PROTO_NC_USER_JP_PASSWORD_CHECK_ACK // 306B
	NClientHandle uint16                              // 2B
	BUserIPBlock  uint8                               // 1B
}

// PROTO_NC_USER_CH_PASSWORD_CHECK_REQ (OpCode 0x0C53, Size: 274 bytes)
type PROTO_NC_USER_CH_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16    // 2B
	User            [256]byte // 256B
	Password        [16]byte  // 16B
}

// PROTO_NC_USER_CH_PASSWORD_CHECK_ACK (OpCode 0x0C54, Size: 287 bytes)
type PROTO_NC_USER_CH_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16    // 2B
	ID              [256]byte // 256B
	Password        [16]byte  // 16B
	Usrregnum       uint32    // 4B
	Authid          uint8     // 1B
	Block           uint32    // 4B
	Loginable       uint32    // 4B
}

// PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND (OpCode 0x0C54, Size: 290 bytes)
type PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND struct {
	Netpacketheader uint16                              // 2B
	NDummy          uint8                               // 1B
	Data            PROTO_NC_USER_CH_PASSWORD_CHECK_ACK // 287B
}

// PROTO_NC_USER_US_LOGIN_REQ (OpCode 0x0C5A, Size: 316 bytes)
type PROTO_NC_USER_US_LOGIN_REQ struct {
	SUserName [260]byte // 260B
	SPassword [36]byte  // 36B
	Spawnapps [20]byte  // 20B
}

// PROTO_NC_USER_US_PASSWORD_CHECK_REQ (OpCode 0x0C5B, Size: 338 bytes)
type PROTO_NC_USER_US_PASSWORD_CHECK_REQ struct {
	Netpacketheader uint16    // 2B
	SUserName       [260]byte // 260B
	SPassword       [36]byte  // 36B
	SUserIP         [20]byte  // 20B
	Spawnapps       [20]byte  // 20B
}

// PROTO_NC_USER_US_PASSWORD_CHECK_ACK (OpCode 0x0C5C, Size: 353 bytes)
type PROTO_NC_USER_US_PASSWORD_CHECK_ACK struct {
	Netpacketheader uint16    // 2B
	NError          uint16    // 2B
	NUserNo         uint32    // 4B
	SUserName       [260]byte // 260B
	SPassword       [36]byte  // 36B
	SUserIP         [20]byte  // 20B
	NAuthID         uint8     // 1B
	NBlock          uint32    // 4B
	BLoginAble      uint32    // 4B
	Spawnapps       [20]byte  // 20B
}

// PROTO_NC_USER_US_IS_IP_BLOCK_REQ (OpCode 0x0C5D, Size: 379 bytes)
type PROTO_NC_USER_US_IS_IP_BLOCK_REQ struct {
	LoginData     PROTO_NC_USER_US_PASSWORD_CHECK_ACK // 353B
	NUserNo       uint32                              // 4B
	NClientHandle uint16                              // 2B
	SUserIP       [20]byte                            // 20B
}

// PROTO_NC_USER_US_IS_IP_BLOCK_ACK (OpCode 0x0C5E, Size: 356 bytes)
type PROTO_NC_USER_US_IS_IP_BLOCK_ACK struct {
	LoginData     PROTO_NC_USER_US_PASSWORD_CHECK_ACK // 353B
	NClientHandle uint16                              // 2B
	BUserIPBlock  uint8                               // 1B
}

// PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD (OpCode 0x0C5F, Size: 1 byte)
type PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD struct {
	BCanSkip uint8 // 1: permite omitir tutorial, 0: obligatorio
}

// PROTO_NC_USER_RETURN_CHECK_REQ (OpCode 0x0C60, Size: 6 bytes)
type PROTO_NC_USER_RETURN_CHECK_REQ struct {
	NHandle uint16 // Handle de sesión cliente
	NUserNo uint32 // nUserNo
}

// PROTO_NC_USER_RETURN_CHECK_ACK (OpCode 0x0C61, Size: 11 bytes)
type PROTO_NC_USER_RETURN_CHECK_ACK struct {
	NHandle     uint16 // Handle de sesión cliente
	NUserNo     uint32 // nUserNo
	BReturnUser uint8  // 1: usuario de retorno habilitado, 0: no
	TLastSec    int32  // Segundos desde última conexión
}

// PROTO_NC_USER_LOGOUT_LAST_DAY_REQ (OpCode 0x0C62, Size: 6 bytes)
type PROTO_NC_USER_LOGOUT_LAST_DAY_REQ struct {
	NHandle uint16 // Handle de cliente
	NUserNo uint32 // nUserNo
}

// PROTO_NC_USER_LOGOUT_LAST_DAY_ACK (OpCode 0x0C63, Size: 10 bytes)
type PROTO_NC_USER_LOGOUT_LAST_DAY_ACK struct {
	NHandle  uint16 // Handle de cliente
	NUserNo  uint32 // nUserNo
	TLastSec int32  // Segundos transcurridos desde logout
}

// PROTO_NC_USER_SET_RETURN_CMD (OpCode 0x0C64, Size: 4 bytes)
type PROTO_NC_USER_SET_RETURN_CMD struct {
	NUserNo uint32 // nUserNo al que se le aplicó el paquete de retorno
}

// PROTO_NC_USP_USER_CHARACTER_DELETE (Size: 9 bytes)
type PROTO_NC_USP_USER_CHARACTER_DELETE struct {
	NUserNo  uint32 // nUserNo
	NWorldNo uint8  // nWorldNo
	NCharNo  uint32 // nCharNo
}

// PROTO_NC_USP_USER_CHARACTER_INSERT (Size: 50 bytes)
type PROTO_NC_USP_USER_CHARACTER_INSERT struct {
	NUserNo   uint32   // nUserNo
	NWorldNo  uint8    // nWorldNo
	NCharNo   uint32   // nCharNo
	SCharName [40]byte // Nombre del personaje
	NClass    uint8    // Clase
}
