package protocol

import (
	"bytes"
	"testing"
)

func TestPacketSerializationAndCrypto(t *testing.T) {
	pw := NewPacketWriter(SH3Type_VersionAllowed)
	pw.WriteInt16(1)

	wire := pw.BuildWirePacket()

	// Simulate Client receiving wire packet from Server (Server packets are unencrypted)
	// Check length byte
	if wire[0] != 4 { // 2 bytes opcode + 2 bytes int16 = 4 bytes
		t.Fatalf("expected wire length byte to be 4, got %d", wire[0])
	}

	payload := make([]byte, wire[0])
	copy(payload, wire[1:])

	packet := NewPacket(payload)
	if packet == nil {
		t.Fatalf("failed to parse packet")
	}

	if packet.OpCode != SH3Type_VersionAllowed {
		t.Fatalf("expected opcode %d, got %d", SH3Type_VersionAllowed, packet.OpCode)
	}

	val, err := packet.Reader.ReadInt16()
	if err != nil {
		t.Fatalf("failed to read int16: %v", err)
	}

	if val != 1 {
		t.Fatalf("expected value 1, got %d", val)
	}
}

func TestOpCodePacking(t *testing.T) {
	header := uint8(3)
	typ := uint16(90) // CH3Type_Login

	opCode := MakeOpCode(header, typ)
	if ExtractHeader(opCode) != header {
		t.Fatalf("expected header %d, got %d", header, ExtractHeader(opCode))
	}
	if ExtractType(opCode) != typ {
		t.Fatalf("expected type %d, got %d", typ, ExtractType(opCode))
	}
}

func TestStringPadding(t *testing.T) {
	pw := NewPacketWriter(100)
	pw.WriteString("TestUser", 16)

	reader := NewPacketReader(pw.Data()[2:]) // skip opcode
	str, err := reader.ReadString(16)
	if err != nil {
		t.Fatalf("failed to read string: %v", err)
	}
	if str != "TestUser" {
		t.Fatalf("expected 'TestUser', got '%s'", str)
	}
	if !bytes.Equal(pw.Data()[2+8:], make([]byte, 8)) {
		t.Fatalf("expected remaining 8 bytes to be null padded")
	}
}
