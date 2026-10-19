package zone

import (
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func readServerPacket(r io.Reader) (*protocol.Packet, error) {
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}

	length := int(lenBuf[0])
	if length == 0 {
		extBuf := make([]byte, 2)
		if _, err := io.ReadFull(r, extBuf); err != nil {
			return nil, err
		}
		length = int(binary.LittleEndian.Uint16(extBuf))
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return protocol.NewPacket(payload), nil
}

func sendClientPacket(w io.Writer, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error {
	bufCopy := make([]byte, len(pw.Data()))
	copy(bufCopy, pw.Data())
	if nc != nil {
		nc.Crypt(bufCopy)
	}

	payloadLen := len(bufCopy)
	var wire []byte
	if payloadLen <= 255 {
		wire = make([]byte, 1+payloadLen)
		wire[0] = byte(payloadLen)
		copy(wire[1:], bufCopy)
	} else {
		wire = make([]byte, 3+payloadLen)
		wire[0] = 0x00
		binary.LittleEndian.PutUint16(wire[1:3], uint16(payloadLen))
		copy(wire[3:], bufCopy)
	}

	_, err := w.Write(wire)
	return err
}

func TestZoneServer_Phase1_LifecycleAndHandshake(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, err := classparam.Load(classDir)
	if err != nil {
		t.Fatalf("failed to load classparam: %v", err)
	}

	fieldStore, err := field.Load(mapDir)
	if err != nil {
		t.Fatalf("failed to load field: %v", err)
	}

	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19210",
		Port:     19210,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	// 1. Conectar cliente TCP
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19210", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to zone server: %v", err)
	}
	defer conn.Close()

	// 2. Leer Handshake de semilla XOR (PROTO_NC_MISC_SEED_CMD, OpCode 0x0807)
	pkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake packet: %v", err)
	}

	if pkt.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("expected PROTO_NC_MISC_SEED_CMD (0x0807), got 0x%04X", pkt.OpCode)
	}

	xorPos, err := pkt.Reader.ReadInt16()
	if err != nil {
		t.Fatalf("failed to read xor offset: %v", err)
	}
	if xorPos < 0 {
		t.Fatalf("invalid negative xor offset: %d", xorPos)
	}
}

func TestZoneServer_Phase1_TransferKeyVerification(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, _ := classparam.Load(classDir)
	fieldStore, _ := field.Load(mapDir)
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19211",
		Port:     19211,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	// Crear personaje en el repo
	char := &database.Character{
		ID:        1,
		AccountID: 1,
		Name:      "ZoneHero",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
		XPos:      5600,
		YPos:      8800,
		CurHP:     52,
		CurSP:     10,
	}
	if err := charRepo.CreateCharacter(char); err != nil {
		t.Fatalf("failed to create test char: %v", err)
	}

	// Registrar transferencia válida desde WorldService
	transfers.CreateZoneTransfer(1, "TestUser", char.Name, char.ID, 12345, 0, "127.0.0.1")

	// Conectar cliente
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19211", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Leer handshake inicial y configurar NetCrypto
	hsPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		t.Fatalf("failed to create client crypto: %v", err)
	}

	// Enviar PROTO_NC_MAP_LOGIN_REQ con credenciales válidas
	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(12345) // randomID
	req.WriteString(char.Name, 16)

	if err := sendClientPacket(conn, req, clientCrypto); err != nil {
		t.Fatalf("failed to send PROTO_NC_MAP_LOGIN_REQ: %v", err)
	}

	// Leer respuesta del servidor (primer paquete del InGame chunk: PROTO_NC_CHAR_BASE_CMD)
	respPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if respPkt.OpCode != protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD {
		t.Errorf("expected PROTO_NC_CHAR_CLIENT_BASE_CMD (0x%04X), got 0x%04X", protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD, respPkt.OpCode)
	}
}

func TestZoneServer_Phase1_InvalidTransferRejection(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, _ := classparam.Load(classDir)
	fieldStore, _ := field.Load(mapDir)
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19212",
		Port:     19212,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	conn, err := net.DialTimeout("tcp", "127.0.0.1:19212", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Descartar handshake y configurar NetCrypto
	hsPkt, _ := readServerPacket(conn)
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, _ := crypto.NewNetCrypto(xorPos)

	// Enviar transferencia inexistente
	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(9999) // RandomID inválido
	req.WriteString("GhostChar", 16)

	if err := sendClientPacket(conn, req, clientCrypto); err != nil {
		t.Fatalf("failed to send packet: %v", err)
	}

	respPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read rejection: %v", err)
	}

	if respPkt.OpCode != protocol.PROTO_NC_MAP_LOGINFAIL_ACK {
		t.Errorf("expected PROTO_NC_MAP_LOGINFAIL_ACK (0x%04X), got 0x%04X", protocol.PROTO_NC_MAP_LOGINFAIL_ACK, respPkt.OpCode)
	}
}
