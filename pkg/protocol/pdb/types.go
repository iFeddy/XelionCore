package pdb

// ============================================================================
// Tipos primitivos y cabeceras oficiales C++ (#pragma pack(push, 1))
// Extraídos directamente de los TPI de Login.pdb, Account.pdb y AccountLog.pdb
// ============================================================================

type NETPACKETHEADER uint16
type NETPACKETZONEHEADER [6]byte
type Name4 [16]byte
type Name5 [20]byte
type Name8 [32]byte
type Name256Byte [256]byte
