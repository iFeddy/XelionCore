package pdb

// ============================================================================
// Estructuras oficiales de Reportes GM y Wrappers DATA_SEND1
// Departamento 32: NC_DATA (Reportes GM)
// Extraídas directamente de docs/Servicios/AccountLog.md y WorldManager.md
// Fidelidad 100% C++ PDB (#pragma pack(push, 1))
// ============================================================================

// PROTO_NC_DATA_REPORT_ADD_REQ (OpCode 0x8007, Size: 2369 bytes)
type PROTO_NC_DATA_REPORT_ADD_REQ struct {
	DwReporter      uint32     // dwReporter (4B)
	HReporter       uint16     // hReporter (2B)
	NWorldNo        uint8      // nWorldNo (1B)
	SReporterCharID [20]byte   // Name5 (20B)
	SHarmer         [20]byte   // Name5 (20B)
	ByReportType    [16]byte   // char[16] (16B)
	ByReason        [256]byte  // char[256] (256B)
	NChatLogSize    uint16     // unsigned short (2B)
	ByChatLog       [2048]byte // unsigned char[2048] (2048B)
}

// PROTO_NC_DATA_REPORT_ADD_ACK (OpCode 0x8008, Size: 28 bytes)
type PROTO_NC_DATA_REPORT_ADD_ACK struct {
	Err        uint16   // unsigned short (2B)
	HReporter  uint16   // unsigned short (2B)
	SHarmer    [20]byte // Name5 (20B)
	NRegNumber int32    // int (4B)
}

// PROTO_NC_DATA_REPORT_CANCEL_REQ (OpCode 0x8009, Size: 6 bytes)
type PROTO_NC_DATA_REPORT_CANCEL_REQ struct {
	NRegNumber int32  // int (4B)
	HReporter  uint16 // unsigned short (2B)
}

// PROTO_NC_DATA_REPORT_CANCEL_ACK (OpCode 0x800A, Size: 8 bytes)
type PROTO_NC_DATA_REPORT_CANCEL_ACK struct {
	NRegNumber int32  // int (4B)
	HReporter  uint16 // unsigned short (2B)
	Err        uint16 // unsigned short (2B)
}

// PROTO_NC_DATA_REPORT_GET_REQ (OpCode 0x800B, Size: 7 bytes)
type PROTO_NC_DATA_REPORT_GET_REQ struct {
	HReporter  uint16 // unsigned short (2B)
	DwReporter uint32 // unsigned int (4B)
	NWorldNo   uint8  // unsigned char (1B)
}

// PROTO_NC_DATA_REPORT_GET_ACK (OpCode 0x800C, Size: 137 bytes)
type PROTO_NC_DATA_REPORT_GET_ACK struct {
	HReporter      uint16    // unsigned short (2B)
	Err            uint16    // unsigned short (2B)
	NRegNumber     int32     // int (4B)
	ByReportResult uint8     // unsigned char (1B)
	ByRemark       [128]byte // unsigned char[128] (128B)
}

// Wrappers PROTO_NC_DATA_SEND1 (header de 2B + dummy de 1B + struct empaquetada)

// PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_ADD_ACK (Size: 31 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_ADD_ACK struct {
	Netpacketheader uint16                       // 2B
	NDummy          uint8                        // 1B
	Data            PROTO_NC_DATA_REPORT_ADD_ACK // 28B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_CANCEL_ACK (Size: 11 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_CANCEL_ACK struct {
	Netpacketheader uint16                          // 2B
	NDummy          uint8                           // 1B
	Data            PROTO_NC_DATA_REPORT_CANCEL_ACK // 8B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_GET_ACK (Size: 140 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_GET_ACK struct {
	Netpacketheader uint16                       // 2B
	NDummy          uint8                        // 1B
	Data            PROTO_NC_DATA_REPORT_GET_ACK // 137B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_L20_DB_ACK (Size: 16 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_L20_DB_ACK struct {
	Netpacketheader uint16                     // 2B
	NDummy          uint8                      // 1B
	Data            PROTO_NC_MISC_EVENT_L20_DB_ACK // 13B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_USER_GER_IS_IP_BLOCK_ACK (Size: 293 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_USER_GER_IS_IP_BLOCK_ACK struct {
	Netpacketheader uint16                            // 2B
	NDummy          uint8                             // 1B
	Data            PROTO_NC_USER_GER_IS_IP_BLOCK_ACK // 290B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_USER_IS_IP_BLOCK_ACK (Size: 278 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_USER_IS_IP_BLOCK_ACK struct {
	Netpacketheader uint16                        // 2B
	NDummy          uint8                         // 1B
	Data            PROTO_NC_USER_IS_IP_BLOCK_ACK // 275B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_USER_TW_IS_IP_BLOCK_ACK (Size: 14 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_USER_TW_IS_IP_BLOCK_ACK struct {
	Netpacketheader uint16                           // 2B
	NDummy          uint8                            // 1B
	Data            PROTO_NC_USER_TW_IS_IP_BLOCK_ACK // 11B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK (Size: 12 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK struct {
	Netpacketheader uint16                                // 2B
	NDummy          uint8                                 // 1B
	Data            PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK // 9B
}

// PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM (Size: 18 bytes)
type PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM struct {
	Netpacketheader uint16                                  // 2B
	NDummy          uint8                                   // 1B
	Data            PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM // 15B
}
