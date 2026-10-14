package protocol

// MakeOpCode combines a 6-bit Header and a 10-bit Type into a 16-bit OpCode.
// Formula: (Header << 10) | (Type & 0x3FF)
func MakeOpCode(header uint8, typ uint16) uint16 {
	return (uint16(header) << 10) | (typ & 0x03FF)
}

// ExtractHeader extracts the 6-bit header from a 16-bit OpCode.
func ExtractHeader(opCode uint16) uint8 {
	return uint8(opCode >> 10)
}

// ExtractType extracts the 10-bit type from a 16-bit OpCode.
func ExtractType(opCode uint16) uint16 {
	return opCode & 0x03FF
}

// Server -> Client OpCodes (SH = Server Header)
var (
	// Header 2: System / Handshake
	SH2Type_Ping              = MakeOpCode(2, 4)
	SH2Type_SetXorKeyPosition = MakeOpCode(2, 7) // Handshake offset
	SH2Type_UpdateClientTime  = MakeOpCode(2, 73)

	// Header 3: Login / World Selection
	SH3Type_IncorrectVersion = MakeOpCode(3, 2)
	SH3Type_FilecheckAllow   = MakeOpCode(3, 5)
	SH3Type_Error            = MakeOpCode(3, 9)
	SH3Type_WorldlistNew     = MakeOpCode(3, 10)
	SH3Type_WorldServerIP    = MakeOpCode(3, 12)
	SH3Type_CharacterList    = MakeOpCode(3, 20)
	SH3Type_WorldistResend   = MakeOpCode(3, 28)
	SH3Type_VersionAllowed   = MakeOpCode(3, 103)
)

// Client -> Server OpCodes (CH = Client Header)
var (
	// Header 2: System / Handshake
	CH2Type_Pong = MakeOpCode(2, 5)

	// Header 3: Login / World Selection
	CH3Type_FileHash          = MakeOpCode(3, 4)
	CH3Type_WorldSelect       = MakeOpCode(3, 11)
	CH3Type_WorldClientKey    = MakeOpCode(3, 15)
	CH3Type_BackToCharSelect  = MakeOpCode(3, 24)
	CH3Type_WorldReRequest    = MakeOpCode(3, 27)
	CH3Type_BackToWorldSelect = MakeOpCode(3, 51)
	CH3Type_Login             = MakeOpCode(3, 90)  // Corrected for NA2016 client
	CH3Type_Version           = MakeOpCode(3, 101)
)

// ServerError codes returned in SH3Type_Error packets.
type ServerError uint16

const (
	ErrException          ServerError = 0x42 // Unknown Error / Exception
	ErrDatabase           ServerError = 0x43 // DB Error
	ErrWrongRegion        ServerError = 0x44 // Auth Failed / Wrong Region
	ErrInvalidCredentials ServerError = 0x45 // Please check ID or password
	ErrDatabaseError      ServerError = 0x46 // DB Error
	ErrBlocked            ServerError = 0x47 // The ID has been blocked
	ErrServerMaintenance  ServerError = 0x48 // World server maintenance
	ErrTimeout            ServerError = 0x49 // Timeout
	ErrLoginFailed        ServerError = 0x4A // Login Failed
	ErrAgreementMissing   ServerError = 0x4B // Please accept the agreement
)

// WorldStatus states for server selection display.
type WorldStatus uint8

const (
	WorldStatus_Offline         WorldStatus = 0
	WorldStatus_Maintenance     WorldStatus = 1
	WorldStatus_EmptyServer     WorldStatus = 2
	WorldStatus_Reserved        WorldStatus = 3
	WorldStatus_OfflineUnkError WorldStatus = 4
	WorldStatus_Full            WorldStatus = 5
	WorldStatus_Low             WorldStatus = 6
	WorldStatus_Medium          WorldStatus = 9
	WorldStatus_High            WorldStatus = 10
)
