package login

import (
	"net"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestNA2026_LoginFlow(t *testing.T) {
	accRepo := database.NewMemoryAccountRepository()
	transfers := database.NewTransferManager()
	wm := NewWorldManager()
	wm.AddWorld(&World{
		ID:     0,
		Name:   "Xelion",
		IP:     "192.168.0.4",
		Port:   9310,
		Status: protocol.WorldStatus_Low,
	})

	cfg := &Config{
		BindAddr: "127.0.0.1:19026",
		Port:     19026,
	}
	srv := NewServer(cfg, accRepo, wm, transfers)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start login server: %v", err)
	}
	defer srv.Stop()

	// Crear cuenta de prueba
	_, _ = accRepo.CreateAccount("admin", "admin")

	conn, err := net.DialTimeout("tcp", "127.0.0.1:19026", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 1. Recibir Handshake (0x0807)
	hs, err := readServerPacket(conn)
	if err != nil || hs.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("failed to receive seed: %v", err)
	}
	seed, _ := hs.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(seed)
	if err != nil {
		t.Fatalf("crypto error: %v", err)
	}

	// 2. Enviar Version Check NA2026 (0x0C2C)
	vPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_REQ)
	vPkt.WriteString("3dca91289ab564b3ce1f5c8ca1b362a6", 32)
	if err := sendClientPacket(conn, vPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send version check: %v", err)
	}

	// Recibir 0x0C2D (3 bytes)
	vAck, err := readServerPacket(conn)
	if err != nil || vAck.OpCode != protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_ACK {
		t.Fatalf("expected 0x0C2D, got 0x%04X (err: %v)", vAck.OpCode, err)
	}
	t.Logf("Paso 1: Versión NA2026 validada exitosamente (0x0C2D)")

	// 3. Enviar Credenciales NA2026 (0x0C01, 351 bytes payload)
	loginPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NA2026_LOGIN_REQ)
	loginPkt.Fill(32, 0) // 32B padding inicial
	loginPkt.WriteString("admin", 260)
	loginPkt.WriteString(md5Hash("admin"), 32)
	loginPkt.WriteString("Original", 23)

	if err := sendClientPacket(conn, loginPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send login: %v", err)
	}

	// Recibir secuencia oficial NA2026:
	// a. 0x0C40 (Auth OK)
	ack1, err := readServerPacket(conn)
	if err != nil || ack1.OpCode != protocol.PROTO_NC_USER_NA2026_AUTH_SUCCESS_ACK {
		t.Fatalf("expected 0x0C40, got 0x%04X (err: %v)", ack1.OpCode, err)
	}

	// b. 0x0C47 (Pre-worlds)
	ack2, err := readServerPacket(conn)
	if err != nil || ack2.OpCode != protocol.PROTO_NC_USER_NA2026_PRE_WORLDS_ACK {
		t.Fatalf("expected 0x0C47, got 0x%04X (err: %v)", ack2.OpCode, err)
	}

	// c. 0x0C06 (World list)
	ack3, err := readServerPacket(conn)
	if err != nil || ack3.OpCode != protocol.PROTO_NC_USER_NA2026_LOGIN_ACK {
		t.Fatalf("expected 0x0C06, got 0x%04X (err: %v)", ack3.OpCode, err)
	}
	numWorlds, _ := ack3.Reader.ReadByte()
	if numWorlds != 16 {
		t.Fatalf("expected 16 worlds in NA2026 official table, got %d", numWorlds)
	}

	// d. 0x0C3E (Account token)
	ack4, err := readServerPacket(conn)
	if err != nil || ack4.OpCode != protocol.PROTO_NC_USER_NA2026_ACCOUNT_TOKEN_ACK {
		t.Fatalf("expected 0x0C3E, got 0x%04X (err: %v)", ack4.OpCode, err)
	}
	t.Logf("Paso 2: Secuencia oficial NA2026 de login recibida completa (0x0C40, 0x0C47, 0x0C06, 0x0C3E)")

	// 4. Probar selección de mundo NA2026 con World 8 (0x0C0A)
	selPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NA2026_WORLDSELECT_REQ)
	selPkt.WriteByte(8) // World 8 (Enid ES / Xelion)
	if err := sendClientPacket(conn, selPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send world select: %v", err)
	}

	// Recibir 0x0C0B
	selAck, err := readServerPacket(conn)
	if err != nil || selAck.OpCode != protocol.PROTO_NC_USER_NA2026_WORLDSELECT_ACK {
		t.Fatalf("expected 0x0C0B, got 0x%04X (err: %v)", selAck.OpCode, err)
	}
	wStatus, _ := selAck.Reader.ReadByte()
	wIP, _ := selAck.Reader.ReadString(16)
	wPort, _ := selAck.Reader.ReadUint16()
	t.Logf("Paso 3: Transferencia NA2026 a World 8 autorizada (Status=%d, Destino=%s:%d)", wStatus, wIP, wPort)
}

func TestLogin_VersionValidationMismatch(t *testing.T) {
	accRepo := database.NewMemoryAccountRepository()
	transfers := database.NewTransferManager()
	wm := NewWorldManager()
	wm.AddWorld(&World{
		ID:     0,
		Name:   "Xelion",
		Status: protocol.WorldStatus_Low,
		IP:     "127.0.0.1",
		Port:   9014,
	})

	// Servidor configurado con versión estricta '10022024000000'
	srvCfg := &Config{
		BindAddr:      "127.0.0.1:19027",
		Port:          19027,
		ClientVersion: "10022024000000",
	}
	srv := NewServer(srvCfg, accRepo, wm, transfers)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer srv.Stop()

	// 1. Cliente NA2026 conectándose con hash '3dca91289ab564b3ce1f5c8ca1b362a6' que NO coincide con '10022024000000'
	conn, err := net.Dial("tcp", "127.0.0.1:19027")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	handshake, err := readServerPacket(conn)
	if err != nil || handshake.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("expected seed: %v", err)
	}
	seed, _ := handshake.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(seed)
	if err != nil {
		t.Fatalf("crypto error: %v", err)
	}

	verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_REQ)
	verPkt.WriteString("3dca91289ab564b3ce1f5c8ca1b362a6", 32)
	_ = sendClientPacket(conn, verPkt, clientCrypto)

	// El servidor DEBE responder con 0x0C2D con status 1 (error de versión)
	ack, err := readServerPacket(conn)
	if err != nil || ack.OpCode != protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_ACK {
		t.Fatalf("expected 0x0C2D, got 0x%04X (err: %v)", ack.OpCode, err)
	}
	status, _ := ack.Reader.ReadByte()
	if status == 0 {
		t.Fatalf("se esperaba rechazo de versión (status != 0), pero el servidor aceptó con status 0")
	}
	t.Logf("Verificación estricta exitosa: Cliente NA2026 rechazado con status=%d por no coincidir con CLIENT_VERSION en .ENV", status)

	// 2. Cliente Legacy enviando versión incorrecta
	conn2, err := net.Dial("tcp", "127.0.0.1:19027")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn2.Close()

	handshake2, err := readServerPacket(conn2)
	if err != nil || handshake2.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("expected seed: %v", err)
	}
	seed2, _ := handshake2.Reader.ReadInt16()
	clientCrypto2, err := crypto.NewNetCrypto(seed2)
	if err != nil {
		t.Fatalf("crypto error: %v", err)
	}

	verPkt2 := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt2.WriteString("VERSION_ERRONEA_TEST", 64)
	_ = sendClientPacket(conn2, verPkt2, clientCrypto2)

	// Servidor debe rechazar con 0x0C66 (PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK)
	ack2, err := readServerPacket(conn2)
	if err != nil || ack2.OpCode != protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK {
		t.Fatalf("expected 0x0C66 (PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK), got 0x%04X (err: %v)", ack2.OpCode, err)
	}
	t.Logf("Verificación estricta exitosa: Cliente legacy rechazado con 0x0C66 por no coincidir con CLIENT_VERSION en .ENV")
}

func readServerPacket(conn net.Conn) (*protocol.Packet, error) {
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var lenBuf [1]byte
	if _, err := conn.Read(lenBuf[:]); err != nil {
		return nil, err
	}
	var length int
	if lenBuf[0] == 0 {
		var extBuf [2]byte
		if _, err := conn.Read(extBuf[:]); err != nil {
			return nil, err
		}
		length = int(extBuf[0]) | (int(extBuf[1]) << 8)
	} else {
		length = int(lenBuf[0])
	}
	payload := make([]byte, length)
	readBytes := 0
	for readBytes < length {
		n, err := conn.Read(payload[readBytes:])
		if err != nil {
			return nil, err
		}
		readBytes += n
	}
	return protocol.NewPacket(payload), nil
}

func sendClientPacket(conn net.Conn, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error {
	raw := pw.Data()
	cryptBuf := make([]byte, len(raw))
	copy(cryptBuf, raw)
	nc.Crypt(cryptBuf)

	var wire []byte
	if len(cryptBuf) <= 255 {
		wire = make([]byte, 1+len(cryptBuf))
		wire[0] = byte(len(cryptBuf))
		copy(wire[1:], cryptBuf)
	} else {
		wire = make([]byte, 3+len(cryptBuf))
		wire[0] = 0
		wire[1] = byte(len(cryptBuf) & 0xFF)
		wire[2] = byte((len(cryptBuf) >> 8) & 0xFF)
		copy(wire[3:], cryptBuf)
	}
	_, err := conn.Write(wire)
	return err
}
