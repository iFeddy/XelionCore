package login_test

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
)

type mockClient struct {
	conn   net.Conn
	crypto *crypto.NetCrypto
}

func readServerPacket(conn net.Conn) (*protocol.Packet, error) {
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}

	var payloadLen int
	if lenBuf[0] == 0 {
		extBuf := make([]byte, 2)
		if _, err := io.ReadFull(conn, extBuf); err != nil {
			return nil, err
		}
		payloadLen = int(binary.LittleEndian.Uint16(extBuf))
	} else {
		payloadLen = int(lenBuf[0])
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	return protocol.NewPacket(payload), nil
}

func (c *mockClient) sendPacket(pw *protocol.PacketWriter) error {
	bufCopy := make([]byte, len(pw.Data()))
	copy(bufCopy, pw.Data())
	c.crypto.Crypt(bufCopy)

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

	_, err := c.conn.Write(wire)
	return err
}

func TestE2ELoginServiceFlow(t *testing.T) {
	cfg, err := config.Load("../../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping MySQL: %v", err)
	}

	repo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager()
	wm := login.NewWorldManager()
	wm.AddWorld(&login.World{
		ID:     0,
		Name:   "Xelion",
		Status: protocol.WorldStatus_Low,
		IP:     "127.0.0.1",
		Port:   9014,
	})

	testPort := 19010
	loginCfg := &login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", testPort),
		Port:     testPort,
	}

	server := login.NewServer(loginCfg, repo, wm, transfers)
	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start LoginServer: %v", err)
	}
	defer server.Stop()

	time.Sleep(100 * time.Millisecond)

	// =========================================================================
	// ESCENARIO 1: Autenticación exitosa de 'admin' con flujo oficial
	// =========================================================================
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatalf("Failed to connect to LoginServer: %v", err)
	}
	defer conn.Close()

	// 1. Recibir Handshake inicial: PROTO_NC_MISC_SEED_CMD (0x0807)
	handshakePkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read initial handshake: %v", err)
	}
	if handshakePkt.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_MISC_SEED_CMD), got 0x%04X",
			protocol.PROTO_NC_MISC_SEED_CMD, handshakePkt.OpCode)
	}
	xorOffset, err := handshakePkt.Reader.ReadInt16()
	if err != nil {
		t.Fatalf("Failed to read xor offset: %v", err)
	}
	t.Logf("Paso 1: Handshake recibido con XOR Offset = %d", xorOffset)

	nc, err := crypto.NewNetCrypto(xorOffset)
	if err != nil {
		t.Fatalf("Failed to initialize NetCrypto: %v", err)
	}
	client := &mockClient{conn: conn, crypto: nc}

	// 2. Enviar versión: PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ (0x0C65)
	verReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verReq.WriteString("10022024000000", 64)
	if err := client.sendPacket(verReq); err != nil {
		t.Fatalf("Failed to send version check req: %v", err)
	}

	// 3. Recibir: PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK (0x0C67)
	verAck, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read version check ack: %v", err)
	}
	if verAck.OpCode != protocol.PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK {
		t.Fatalf("Expected OpCode 0x%04X (RIGHTVERSION_CHECK_ACK), got 0x%04X",
			protocol.PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK, verAck.OpCode)
	}
	t.Logf("Paso 2: Versión validada exitosamente (0x%04X)", verAck.OpCode)

	// 4. Enviar Login oficial: PROTO_NC_USER_US_LOGIN_REQ (0x0C5A, 316 bytes)
	loginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	loginReq.WriteString("admin", 260)
	loginReq.WriteString("admin", 36)
	loginReq.WriteString("Original", 20)
	if err := client.sendPacket(loginReq); err != nil {
		t.Fatalf("Failed to send login req: %v", err)
	}

	// 5. Recibir: PROTO_NC_USER_LOGIN_ACK (0x0C0A)
	loginAck, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read login ack: %v", err)
	}
	if loginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_USER_LOGIN_ACK), got 0x%04X",
			protocol.PROTO_NC_USER_LOGIN_ACK, loginAck.OpCode)
	}
	numWorlds, _ := loginAck.Reader.ReadByte()
	t.Logf("Paso 3: Login exitoso, recibido listado con %d slots de mundos (0x%04X)", numWorlds, loginAck.OpCode)

	// 6. Enviar selección de mundo: PROTO_NC_USER_WORLDSELECT_REQ (0x0C0B)
	selectReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	selectReq.WriteByte(0) // World ID 0 (Xelion)
	if err := client.sendPacket(selectReq); err != nil {
		t.Fatalf("Failed to send world select req: %v", err)
	}

	// 7. Recibir: PROTO_NC_USER_WORLDSELECT_ACK (0x0C0C, 83 bytes)
	selectAck, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read world select ack: %v", err)
	}
	if selectAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_USER_WORLDSELECT_ACK), got 0x%04X",
			protocol.PROTO_NC_USER_WORLDSELECT_ACK, selectAck.OpCode)
	}
	status, _ := selectAck.Reader.ReadByte()
	ip, _ := selectAck.Reader.ReadString(16)
	port, _ := selectAck.Reader.ReadUint16()
	token, _ := selectAck.Reader.ReadString(32)

	t.Logf("Paso 4: Transferencia autorizada hacia %s:%d (Status: %d, Token: %s)", ip, port, status, token)

	// 8. Verificar que la auditoría en tAccountLog se persistió en MySQL
	var successLogs int
	err = db.QueryRow("SELECT COUNT(*) FROM tAccountLog WHERE sUserID = 'admin' AND nType = 1").Scan(&successLogs)
	if err != nil || successLogs == 0 {
		t.Fatalf("Expected audit log for admin in tAccountLog (nType=1), got count=%d, err=%v", successLogs, err)
	}
	t.Logf("Paso 5: Auditoría oficial verificada en MySQL tAccountLog (Filas éxito: %d)", successLogs)

	// =========================================================================
	// ESCENARIO 2: Intento con contraseña incorrecta
	// =========================================================================
	conn2, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	hs2, _ := readServerPacket(conn2)
	off2, _ := hs2.Reader.ReadInt16()
	nc2, _ := crypto.NewNetCrypto(off2)
	client2 := &mockClient{conn: conn2, crypto: nc2}

	badLoginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	badLoginReq.WriteString("admin", 260)
	badLoginReq.WriteString("wrong_password", 36)
	badLoginReq.WriteString("Original", 20)
	if err := client2.sendPacket(badLoginReq); err != nil {
		t.Fatalf("Failed to send bad login req: %v", err)
	}

	failAck, err := readServerPacket(conn2)
	if err != nil {
		t.Fatalf("Failed to read fail ack: %v", err)
	}
	if failAck.OpCode != protocol.PROTO_NC_USER_LOGINFAIL_ACK {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_USER_LOGINFAIL_ACK), got 0x%04X",
			protocol.PROTO_NC_USER_LOGINFAIL_ACK, failAck.OpCode)
	}
	errCode, _ := failAck.Reader.ReadUint16()
	if errCode != uint16(protocol.ErrInvalidCredentials) {
		t.Fatalf("Expected error code 0x%04X (ErrInvalidCredentials), got 0x%04X",
			protocol.ErrInvalidCredentials, errCode)
	}
	t.Logf("Escenario 2 Verificado: Servidor rechazó contraseña errónea con PROTO_NC_USER_LOGINFAIL_ACK (0x%04X)", errCode)
}
