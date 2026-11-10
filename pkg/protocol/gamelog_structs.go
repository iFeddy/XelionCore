package protocol

import (
	"bytes"
	"encoding/binary"
)

// ============================================================================
// Estructuras Oficiales C++ (#pragma pack(push, 1)) — gamelog_structs.go
// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.
// ============================================================================

// PROTO_NC_LOG_GAME_DATA_TYPE_0 (sizeof = 68 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
	NTargetID                        uint32               // offset: 36, size: 4 bytes
	NItemKey                         uint64               // offset: 40, size: 8 bytes
	NInt1                            uint32               // offset: 48, size: 4 bytes
	NInt2                            uint32               // offset: 52, size: 4 bytes
	NInt3                            uint32               // offset: 56, size: 4 bytes
	NBigint1                         uint64               // offset: 60, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA) FixedWireSize() int { return 68 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_1 (sizeof = 32 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_1_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_1_DATA) FixedWireSize() int { return 32 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_1_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_1_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_2 (sizeof = 40 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_2_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NInt1                            uint32               // offset: 32, size: 4 bytes
	NInt2                            uint32               // offset: 36, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_2_DATA) FixedWireSize() int { return 40 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_2_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_2_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_3 (sizeof = 44 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_3_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
	NTargetID                        uint32               // offset: 36, size: 4 bytes
	NInt1                            uint32               // offset: 40, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_3_DATA) FixedWireSize() int { return 44 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_3_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_3_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_4 (sizeof = 44 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_4_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NInt1                            uint32               // offset: 32, size: 4 bytes
	NInt2                            uint32               // offset: 36, size: 4 bytes
	NInt3                            uint32               // offset: 40, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_4_DATA) FixedWireSize() int { return 44 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_4_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_4_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_5 (sizeof = 16 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_5_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NInt1                            uint32               // offset: 8, size: 4 bytes
	NInt2                            uint32               // offset: 12, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_5_DATA) FixedWireSize() int { return 16 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_5_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_5_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_6 (sizeof = 44 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_6_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NInt1                            uint32               // offset: 32, size: 4 bytes
	NBigint1                         uint64               // offset: 36, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_6_DATA) FixedWireSize() int { return 44 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_6_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_6_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_7 (sizeof = 44 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_7_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
	NInt1                            uint32               // offset: 36, size: 4 bytes
	NInt2                            uint32               // offset: 40, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_7_DATA) FixedWireSize() int { return 44 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_7_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_7_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_8 (sizeof = 36 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_8_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NInt1                            uint32               // offset: 32, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_8_DATA) FixedWireSize() int { return 36 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_8_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_8_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_9 (sizeof = 48 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_9_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NItemKey                         uint64               // offset: 32, size: 8 bytes
	NInt1                            uint32               // offset: 40, size: 4 bytes
	NInt2                            uint32               // offset: 44, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_9_DATA) FixedWireSize() int { return 48 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_9_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_9_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_A (sizeof = 36 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_A_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 8, size: 4 bytes
	NTargetID                        uint32               // offset: 12, size: 4 bytes
	NItemKey                         uint64               // offset: 16, size: 8 bytes
	NInt1                            uint32               // offset: 24, size: 4 bytes
	NInt2                            uint32               // offset: 28, size: 4 bytes
	NInt3                            uint32               // offset: 32, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_A_DATA) FixedWireSize() int { return 36 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_A_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_A_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_B (sizeof = 60 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_B_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
	NTargetID                        uint32               // offset: 36, size: 4 bytes
	NItemKey                         uint64               // offset: 40, size: 8 bytes
	NInt1                            uint32               // offset: 48, size: 4 bytes
	NInt2                            uint32               // offset: 52, size: 4 bytes
	NInt3                            uint32               // offset: 56, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_B_DATA) FixedWireSize() int { return 60 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_B_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_B_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_C (sizeof = 56 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_C_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetID                        uint32               // offset: 32, size: 4 bytes
	NItemKey                         uint64               // offset: 36, size: 8 bytes
	NInt1                            uint32               // offset: 44, size: 4 bytes
	NInt2                            uint32               // offset: 48, size: 4 bytes
	NInt3                            uint32               // offset: 52, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_C_DATA) FixedWireSize() int { return 56 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_C_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_C_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_D (sizeof = 24 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_D_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	NItemKey                         uint64               // offset: 8, size: 8 bytes
	NInt1                            uint32               // offset: 16, size: 4 bytes
	NInt2                            uint32               // offset: 20, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_D_DATA) FixedWireSize() int { return 24 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_D_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_D_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_E (sizeof = 48 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_E_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	SMap                             Name3                // offset: 4, size: 12 bytes
	NMapX                            uint32               // offset: 16, size: 4 bytes
	NMapY                            uint32               // offset: 20, size: 4 bytes
	NMapZ                            uint32               // offset: 24, size: 4 bytes
	NTargetID                        uint32               // offset: 28, size: 4 bytes
	NItemKey                         uint64               // offset: 32, size: 8 bytes
	NInt1                            uint32               // offset: 40, size: 4 bytes
	NInt2                            uint32               // offset: 44, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_E_DATA) FixedWireSize() int { return 48 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_E_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_E_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_F (sizeof = 64 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_F_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetID                        uint32               // offset: 32, size: 4 bytes
	NItemKey                         uint64               // offset: 36, size: 8 bytes
	NInt1                            uint32               // offset: 44, size: 4 bytes
	NInt2                            uint32               // offset: 48, size: 4 bytes
	NInt3                            uint32               // offset: 52, size: 4 bytes
	NBigint1                         uint64               // offset: 56, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_F_DATA) FixedWireSize() int { return 64 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_F_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_F_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_G (sizeof = 60 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_G_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetID                        uint32               // offset: 32, size: 4 bytes
	NItemKey                         uint64               // offset: 36, size: 8 bytes
	NInt1                            uint32               // offset: 44, size: 4 bytes
	NInt2                            uint32               // offset: 48, size: 4 bytes
	NBigint1                         uint64               // offset: 52, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_G_DATA) FixedWireSize() int { return 60 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_G_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_G_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_H (sizeof = 36 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_H_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_H_DATA) FixedWireSize() int { return 36 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_H_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_H_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_I (sizeof = 40 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_I_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetCharNo                    uint32               // offset: 32, size: 4 bytes
	NTargetID                        uint32               // offset: 36, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_I_DATA) FixedWireSize() int { return 40 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_I_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_I_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_J (sizeof = 52 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_J_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NTargetID                        uint32               // offset: 32, size: 4 bytes
	NItemKey                         uint64               // offset: 36, size: 8 bytes
	NInt1                            uint32               // offset: 44, size: 4 bytes
	NInt2                            uint32               // offset: 48, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_J_DATA) FixedWireSize() int { return 52 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_J_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_J_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_DATA_TYPE_K (sizeof = 48 bytes)
type PROTO_NC_LOG_GAME_DATA_TYPE_K_DATA struct {
	NType                            uint32               // offset: 0, size: 4 bytes
	NCharNo                          uint32               // offset: 4, size: 4 bytes
	SMap                             Name3                // offset: 8, size: 12 bytes
	NMapX                            uint32               // offset: 20, size: 4 bytes
	NMapY                            uint32               // offset: 24, size: 4 bytes
	NMapZ                            uint32               // offset: 28, size: 4 bytes
	NInt1                            uint32               // offset: 32, size: 4 bytes
	NInt2                            uint32               // offset: 36, size: 4 bytes
	NBigint1                         uint64               // offset: 40, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_K_DATA) FixedWireSize() int { return 48 }

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_K_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_DATA_TYPE_K_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO (sizeof = 32 bytes)
type PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Itemkey                          uint64               // offset: 4, size: 8 bytes
	Mappos                           MAPPOS               // offset: 12, size: 20 bytes
}

func (s *PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO_DATA) FixedWireSize() int { return 32 }

func (s *PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED (sizeof = 16 bytes)
type PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Itemkey                          uint64               // offset: 4, size: 8 bytes
	Iteminvenfrom                    ITEM_INVEN           // offset: 12, size: 2 bytes
	Iteminvento                      ITEM_INVEN           // offset: 14, size: 2 bytes
}

func (s *PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED_DATA) FixedWireSize() int { return 16 }

func (s *PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_ITEM_UPGRADE (sizeof = 43 bytes)
type PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	Itemkey                          uint64               // offset: 24, size: 8 bytes
	Itemoldlevel                     uint8                // offset: 32, size: 1 bytes
	Itemnewlevel                     uint8                // offset: 33, size: 1 bytes
	Left_item_id                     uint16               // offset: 34, size: 2 bytes
	Right_item_id                    uint16               // offset: 36, size: 2 bytes
	Middle_item_id                   uint16               // offset: 38, size: 2 bytes
	Randomoption_elemental           [3]byte              // offset: 40, size: 3 bytes
}

func (s *PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA) FixedWireSize() int { return 43 }

func (s *PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_KQ_ENTER (sizeof = 30 bytes)
type PROTO_NC_LOG_GAME_KQ_ENTER_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	KQ_ID                            uint16               // offset: 24, size: 2 bytes
	KQ_Handle                        uint32               // offset: 26, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_KQ_ENTER_DATA) FixedWireSize() int { return 30 }

func (s *PROTO_NC_LOG_GAME_KQ_ENTER_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_KQ_ENTER_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_KQ_LEAVE (sizeof = 31 bytes)
type PROTO_NC_LOG_GAME_KQ_LEAVE_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	KQ_ID                            uint16               // offset: 24, size: 2 bytes
	KQ_Handle                        uint32               // offset: 26, size: 4 bytes
	Leave_Status                     uint8                // offset: 30, size: 1 bytes
}

func (s *PROTO_NC_LOG_GAME_KQ_LEAVE_DATA) FixedWireSize() int { return 31 }

func (s *PROTO_NC_LOG_GAME_KQ_LEAVE_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_KQ_LEAVE_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_LEVEL_UP (sizeof = 33 bytes)
type PROTO_NC_LOG_GAME_LEVEL_UP_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	Level                            uint8                // offset: 24, size: 1 bytes
	Currentexp                       uint64               // offset: 25, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_LEVEL_UP_DATA) FixedWireSize() int { return 33 }

func (s *PROTO_NC_LOG_GAME_LEVEL_UP_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_LEVEL_UP_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING (sizeof = 24 bytes)
type PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
}

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING_DATA) FixedWireSize() int { return 24 }

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MINIHOUSE_VISIT (sizeof = 28 bytes)
type PROTO_NC_LOG_GAME_MINIHOUSE_VISIT_DATA struct {
	Owner_chrregnum                  uint32               // offset: 0, size: 4 bytes
	Visitor_chrregnum                uint32               // offset: 4, size: 4 bytes
	Mappos                           MAPPOS               // offset: 8, size: 20 bytes
}

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_VISIT_DATA) FixedWireSize() int { return 28 }

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_VISIT_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MINIHOUSE_VISIT_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MK2_DEAD (sizeof = 26 bytes)
type PROTO_NC_LOG_GAME_MK2_DEAD_DATA struct {
	Mobid                            uint16               // offset: 0, size: 2 bytes
	Mappos                           MAPPOS               // offset: 2, size: 20 bytes
	NCharNo                          uint32               // offset: 22, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_MK2_DEAD_DATA) FixedWireSize() int { return 26 }

func (s *PROTO_NC_LOG_GAME_MK2_DEAD_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MK2_DEAD_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MK2_START (sizeof = 28 bytes)
type PROTO_NC_LOG_GAME_MK2_START_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	NLevel                           uint8                // offset: 4, size: 1 bytes
	NClass                           uint8                // offset: 5, size: 1 bytes
	Mobid                            uint16               // offset: 6, size: 2 bytes
	Mappos                           MAPPOS               // offset: 8, size: 20 bytes
}

func (s *PROTO_NC_LOG_GAME_MK2_START_DATA) FixedWireSize() int { return 28 }

func (s *PROTO_NC_LOG_GAME_MK2_START_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MK2_START_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MKED (sizeof = 34 bytes)
type PROTO_NC_LOG_GAME_MKED_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	Mobid                            uint16               // offset: 24, size: 2 bytes
	Exp                              uint64               // offset: 26, size: 8 bytes
}

func (s *PROTO_NC_LOG_GAME_MKED_DATA) FixedWireSize() int { return 34 }

func (s *PROTO_NC_LOG_GAME_MKED_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MKED_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MK_DROP_ITEM (sizeof = 34 bytes)
type PROTO_NC_LOG_GAME_MK_DROP_ITEM_DATA struct {
	NItemID                          uint16               // offset: 0, size: 2 bytes
	NLot                             uint16               // offset: 2, size: 2 bytes
	NItemKey                         uint64               // offset: 4, size: 8 bytes
	Mobid                            uint16               // offset: 12, size: 2 bytes
	Mappos                           MAPPOS               // offset: 14, size: 20 bytes
}

func (s *PROTO_NC_LOG_GAME_MK_DROP_ITEM_DATA) FixedWireSize() int { return 34 }

func (s *PROTO_NC_LOG_GAME_MK_DROP_ITEM_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MK_DROP_ITEM_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_MK_GETEXPINFIELD (sizeof = 12 bytes)
type PROTO_NC_LOG_GAME_MK_GETEXPINFIELD_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mobid                            uint16               // offset: 4, size: 2 bytes
	Number                           uint16               // offset: 6, size: 2 bytes
	Exp                              uint32               // offset: 8, size: 4 bytes
}

func (s *PROTO_NC_LOG_GAME_MK_GETEXPINFIELD_DATA) FixedWireSize() int { return 12 }

func (s *PROTO_NC_LOG_GAME_MK_GETEXPINFIELD_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_MK_GETEXPINFIELD_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GAME_PKED (sizeof = 33 bytes)
type PROTO_NC_LOG_GAME_PKED_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	Killerdiechrrgnum                uint32               // offset: 24, size: 4 bytes
	Pkcount                          uint32               // offset: 28, size: 4 bytes
	Pkyellowtime                     uint8                // offset: 32, size: 1 bytes
}

func (s *PROTO_NC_LOG_GAME_PKED_DATA) FixedWireSize() int { return 33 }

func (s *PROTO_NC_LOG_GAME_PKED_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GAME_PKED_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO (sizeof = 41 bytes)
type PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO_DATA struct {
	Chrregnum                        uint32               // offset: 0, size: 4 bytes
	Mappos                           MAPPOS               // offset: 4, size: 20 bytes
	StoreType                        uint8                // offset: 24, size: 1 bytes
	CharExp                          uint64               // offset: 25, size: 8 bytes
	CharMoney                        uint64               // offset: 33, size: 8 bytes
}

func (s *PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO_DATA) FixedWireSize() int { return 41 }

func (s *PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_REGENLOCATESAVE_CMD (sizeof = 38 bytes)
type PROTO_NC_LOG_REGENLOCATESAVE_CMD_DATA struct {
	ChrRegNum                        uint32               // offset: 0, size: 4 bytes
	SaveMap                          Name3                // offset: 4, size: 12 bytes
	SaveLocate                       SHINE_XY_TYPE        // offset: 16, size: 8 bytes
	KQHandle                         uint32               // offset: 24, size: 4 bytes
	InvalidLocate                    SHINE_XY_TYPE        // offset: 28, size: 8 bytes
	Flag                             uint16               // offset: 36, size: 2 bytes
}

func (s *PROTO_NC_LOG_REGENLOCATESAVE_CMD_DATA) FixedWireSize() int { return 38 }

func (s *PROTO_NC_LOG_REGENLOCATESAVE_CMD_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_REGENLOCATESAVE_CMD_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_WEDDING_HALL_CANCEL (sizeof = 8 bytes)
type PROTO_NC_LOG_WEDDING_HALL_CANCEL_DATA struct {
	CharRegisterNumber               uint32               // offset: 0, size: 4 bytes
	PartnerRegisterNumber            uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_LOG_WEDDING_HALL_CANCEL_DATA) FixedWireSize() int { return 8 }

func (s *PROTO_NC_LOG_WEDDING_HALL_CANCEL_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_WEDDING_HALL_CANCEL_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_WEDDING_HALL_RESERVE (sizeof = 8 bytes)
type PROTO_NC_LOG_WEDDING_HALL_RESERVE_DATA struct {
	CharRegisterNumber               uint32               // offset: 0, size: 4 bytes
	PartnerRegisterNumber            uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_LOG_WEDDING_HALL_RESERVE_DATA) FixedWireSize() int { return 8 }

func (s *PROTO_NC_LOG_WEDDING_HALL_RESERVE_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_WEDDING_HALL_RESERVE_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_LOG_WEDDING_HALL_START (sizeof = 8 bytes)
type PROTO_NC_LOG_WEDDING_HALL_START_DATA struct {
	CharRegisterNumber               uint32               // offset: 0, size: 4 bytes
	PartnerRegisterNumber            uint32               // offset: 4, size: 4 bytes
}

func (s *PROTO_NC_LOG_WEDDING_HALL_START_DATA) FixedWireSize() int { return 8 }

func (s *PROTO_NC_LOG_WEDDING_HALL_START_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_LOG_WEDDING_HALL_START_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_SYSLOG_CHAR_DEATH (sizeof = 22 bytes)
type PROTO_NC_SYSLOG_CHAR_DEATH_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	ZoneNo                           uint8                // offset: 8, size: 1 bytes
	AttackerType                     uint8                // offset: 9, size: 1 bytes
	AttackerID                       uint32               // offset: 10, size: 4 bytes
	DeathCoord                       SHINE_XY_TYPE        // offset: 14, size: 8 bytes
}

func (s *PROTO_NC_SYSLOG_CHAR_DEATH_DATA) FixedWireSize() int { return 22 }

func (s *PROTO_NC_SYSLOG_CHAR_DEATH_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_SYSLOG_CHAR_DEATH_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_SYSLOG_CHAR_QUEST_STARTED (sizeof = 13 bytes)
type PROTO_NC_SYSLOG_CHAR_QUEST_STARTED_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	QuestID                          uint16               // offset: 8, size: 2 bytes
	ZoneNo                           uint8                // offset: 10, size: 1 bytes
	QuestStartNPC                    uint16               // offset: 11, size: 2 bytes
}

func (s *PROTO_NC_SYSLOG_CHAR_QUEST_STARTED_DATA) FixedWireSize() int { return 13 }

func (s *PROTO_NC_SYSLOG_CHAR_QUEST_STARTED_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_SYSLOG_CHAR_QUEST_STARTED_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_SYSLOG_CHAR_VICTORY (sizeof = 22 bytes)
type PROTO_NC_SYSLOG_CHAR_VICTORY_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	ZoneNo                           uint8                // offset: 8, size: 1 bytes
	TargetType                       uint8                // offset: 9, size: 1 bytes
	TargetID                         uint32               // offset: 10, size: 4 bytes
	KillCoord                        SHINE_XY_TYPE        // offset: 14, size: 8 bytes
}

func (s *PROTO_NC_SYSLOG_CHAR_VICTORY_DATA) FixedWireSize() int { return 22 }

func (s *PROTO_NC_SYSLOG_CHAR_VICTORY_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_SYSLOG_CHAR_VICTORY_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}

// PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION (sizeof = 11 bytes)
type PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION_DATA struct {
	UserNo                           uint32               // offset: 0, size: 4 bytes
	CharNo                           uint32               // offset: 4, size: 4 bytes
	DestZoneNo                       uint8                // offset: 8, size: 1 bytes
	SrcZoneNo                        uint8                // offset: 9, size: 1 bytes
	CharLev                          uint8                // offset: 10, size: 1 bytes
}

func (s *PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION_DATA) FixedWireSize() int { return 11 }

func (s *PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION_DATA) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION_DATA) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, s)
}
