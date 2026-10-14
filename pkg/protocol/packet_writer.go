package protocol

import (
	"encoding/binary"

	"github.com/iFeddy/XelionCore/pkg/crypto"
)

// PacketWriter constructs packet buffers (Little Endian).
type PacketWriter struct {
	opCode uint16
	buf    []byte
}

// NewPacketWriter creates a PacketWriter initialized with an OpCode.
func NewPacketWriter(opCode uint16) *PacketWriter {
	pw := &PacketWriter{
		opCode: opCode,
		buf:    make([]byte, 0, 64),
	}
	pw.WriteUint16(opCode)
	return pw
}

// OpCode returns the packet's OpCode.
func (pw *PacketWriter) OpCode() uint16 {
	return pw.opCode
}

// WriteByte appends a single byte.
func (pw *PacketWriter) WriteByte(b byte) {
	pw.buf = append(pw.buf, b)
}

// WriteBool appends 1 for true, 0 for false.
func (pw *PacketWriter) WriteBool(v bool) {
	if v {
		pw.WriteByte(1)
	} else {
		pw.WriteByte(0)
	}
}

// WriteUint16 appends an unsigned 16-bit integer in little-endian.
func (pw *PacketWriter) WriteUint16(v uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	pw.buf = append(pw.buf, b...)
}

// WriteInt16 appends a signed 16-bit integer in little-endian.
func (pw *PacketWriter) WriteInt16(v int16) {
	pw.WriteUint16(uint16(v))
}

// WriteUint32 appends an unsigned 32-bit integer in little-endian.
func (pw *PacketWriter) WriteUint32(v uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	pw.buf = append(pw.buf, b...)
}

// WriteInt32 appends a signed 32-bit integer in little-endian.
func (pw *PacketWriter) WriteInt32(v int32) {
	pw.WriteUint32(uint32(v))
}

// WriteUint64 appends an unsigned 64-bit integer in little-endian.
func (pw *PacketWriter) WriteUint64(v uint64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	pw.buf = append(pw.buf, b...)
}

// WriteInt64 appends a signed 64-bit integer in little-endian.
func (pw *PacketWriter) WriteInt64(v int64) {
	pw.WriteUint64(uint64(v))
}

// WriteBytes appends raw bytes.
func (pw *PacketWriter) WriteBytes(b []byte) {
	pw.buf = append(pw.buf, b...)
}

// WriteString appends a string padded with null bytes up to fixedLen.
func (pw *PacketWriter) WriteString(s string, fixedLen int) {
	strBytes := []byte(s)
	if len(strBytes) >= fixedLen {
		pw.buf = append(pw.buf, strBytes[:fixedLen]...)
		return
	}
	pw.buf = append(pw.buf, strBytes...)
	padding := make([]byte, fixedLen-len(strBytes))
	pw.buf = append(pw.buf, padding...)
}

// Fill appends count bytes of value b.
func (pw *PacketWriter) Fill(count int, b byte) {
	for i := 0; i < count; i++ {
		pw.buf = append(pw.buf, b)
	}
}

// Data returns the raw unencrypted payload.
func (pw *PacketWriter) Data() []byte {
	return pw.buf
}

// BuildWirePacket encrypts the payload with the given NetCrypto and adds Fiesta framing.
// Header format:
// - If len(payload) < 255: [1-byte length] + encrypted payload
// - If len(payload) >= 255: [0x00, uint16 length] + encrypted payload
func (pw *PacketWriter) BuildWirePacket(c *crypto.NetCrypto) []byte {
	payloadLen := len(pw.buf)
	encrypted := make([]byte, payloadLen)
	copy(encrypted, pw.buf)
	if c != nil {
		c.Crypt(encrypted)
	}

	if payloadLen < 255 {
		wire := make([]byte, 1+payloadLen)
		wire[0] = byte(payloadLen)
		copy(wire[1:], encrypted)
		return wire
	}

	wire := make([]byte, 3+payloadLen)
	wire[0] = 0x00
	binary.LittleEndian.PutUint16(wire[1:3], uint16(payloadLen))
	copy(wire[3:], encrypted)
	return wire
}
