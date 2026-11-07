package protocol

// ============================================================================
// Tipos primitivos y cabeceras oficiales C++ (#pragma pack(push, 1))
// Extraídos directamente de los TPI de WorldManager.pdb y Character.pdb
// ============================================================================

// NETPACKETHEADER representa la cabecera estándar de paquete de red interno (2 bytes).
type NETPACKETHEADER uint16

// NETPACKETZONEHEADER representa la cabecera de enlace entre Zone y World (6 bytes).
type NETPACKETZONEHEADER [6]byte

// Name3 representa una cadena fija de 12 bytes.
type Name3 [12]byte

// Name4 representa una cadena fija de 16 bytes (IPs, nombres de gremio cortos, mapas).
type Name4 [16]byte

// Name5 representa una cadena fija de 20 bytes (nombre oficial de avatar/personaje).
type Name5 [20]byte

// Name8 representa una cadena fija de 32 bytes (hashes MD5, OTP, contraseñas).
type Name8 [32]byte

// Name256Byte representa una cadena fija de 256 bytes (nombres de cuenta/usuario largos).
type Name256Byte [256]byte
