package protocol

import "encoding/binary"

// Packet represents a decoded Fiesta network packet.
type Packet struct {
	OpCode uint16
	Header uint8
	Type   uint16
	Reader *PacketReader
	Raw    []byte
}

// NewPacket parses a decrypted byte slice into a Packet.
func NewPacket(data []byte) *Packet {
	if len(data) < 2 {
		return nil
	}
	opCode := binary.LittleEndian.Uint16(data[0:2])
	return &Packet{
		OpCode: opCode,
		Header: ExtractHeader(opCode),
		Type:   ExtractType(opCode),
		Reader: NewPacketReader(data[2:]), // Reader starts after the 2-byte OpCode
		Raw:    data,
	}
}
