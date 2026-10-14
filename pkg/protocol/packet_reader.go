package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// PacketReader reads primitive types from a packet payload (Little Endian).
type PacketReader struct {
	data   []byte
	cursor int
}

// NewPacketReader creates a PacketReader for a byte slice.
func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{
		data:   data,
		cursor: 0,
	}
}

// Remaining returns the number of unread bytes.
func (r *PacketReader) Remaining() int {
	return len(r.data) - r.cursor
}

// ReadByte reads a single byte.
func (r *PacketReader) ReadByte() (byte, error) {
	if r.cursor >= len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	b := r.data[r.cursor]
	r.cursor++
	return b, nil
}

// ReadBool reads a boolean byte (0 = false, non-zero = true).
func (r *PacketReader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

// ReadUint16 reads an unsigned 16-bit integer in little-endian.
func (r *PacketReader) ReadUint16() (uint16, error) {
	if r.Remaining() < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	val := binary.LittleEndian.Uint16(r.data[r.cursor : r.cursor+2])
	r.cursor += 2
	return val, nil
}

// ReadInt16 reads a signed 16-bit integer in little-endian.
func (r *PacketReader) ReadInt16() (int16, error) {
	val, err := r.ReadUint16()
	return int16(val), err
}

// ReadUint32 reads an unsigned 32-bit integer in little-endian.
func (r *PacketReader) ReadUint32() (uint32, error) {
	if r.Remaining() < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	val := binary.LittleEndian.Uint32(r.data[r.cursor : r.cursor+4])
	r.cursor += 4
	return val, nil
}

// ReadInt32 reads a signed 32-bit integer in little-endian.
func (r *PacketReader) ReadInt32() (int32, error) {
	val, err := r.ReadUint32()
	return int32(val), err
}

// ReadUint64 reads an unsigned 64-bit integer in little-endian.
func (r *PacketReader) ReadUint64() (uint64, error) {
	if r.Remaining() < 8 {
		return 0, io.ErrUnexpectedEOF
	}
	val := binary.LittleEndian.Uint64(r.data[r.cursor : r.cursor+8])
	r.cursor += 8
	return val, nil
}

// ReadInt64 reads a signed 64-bit integer in little-endian.
func (r *PacketReader) ReadInt64() (int64, error) {
	val, err := r.ReadUint64()
	return int64(val), err
}

// ReadBytes reads exact n bytes.
func (r *PacketReader) ReadBytes(n int) ([]byte, error) {
	if n < 0 || r.Remaining() < n {
		return nil, io.ErrUnexpectedEOF
	}
	res := make([]byte, n)
	copy(res, r.data[r.cursor:r.cursor+n])
	r.cursor += n
	return res, nil
}

// ReadString reads a fixed-length string, stripping any trailing null (0x00) bytes.
func (r *PacketReader) ReadString(fixedLen int) (string, error) {
	buf, err := r.ReadBytes(fixedLen)
	if err != nil {
		return "", err
	}
	idx := bytes.IndexByte(buf, 0x00)
	if idx != -1 {
		buf = buf[:idx]
	}
	return string(buf), nil
}

// Skip advances the cursor by n bytes.
func (r *PacketReader) Skip(n int) error {
	if n < 0 || r.Remaining() < n {
		return errors.New("cannot skip beyond buffer length")
	}
	r.cursor += n
	return nil
}
