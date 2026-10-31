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

	db, err := sql.Open("mysql", cfg.DSN())
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

	// 1. Recibir Handshake inicial: PROTO_NC_MISC_SEED_ACK (0x0807)
	handshakePkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read initial handshake: %v", err)
	}
	if handshakePkt.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_MISC_SEED_ACK), got 0x%04X",
			protocol.PROTO_NC_MISC_SEED_ACK, handshakePkt.OpCode)
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

	// 2. Enviar versión activa en FOUS: PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ (0x0C65)
	activeVerKey, err := repo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("Failed to read FOUS.GameVersionKey from MySQL: %v", err)
	}
	verReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verReq.WriteString(activeVerKey, 64)
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

func TestLoginService_HeartbeatAndLogout(t *testing.T) {
	appCfg, err := config.Load("../../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", appCfg.DSN())
	if err != nil {
		t.Fatalf("Failed to open MySQL: %v", err)
	}
	defer db.Close()

	testPort := 19011
	cfg := &login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", testPort),
		Port:     testPort,
	}

	repo := database.NewMySQLAccountRepository(db)
	wm := login.NewWorldManager()
	transfers := database.NewTransferManager(nil)

	srv := login.NewServer(cfg, repo, wm, transfers)
	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start login server: %v", err)
	}
	defer srv.Stop()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Leer handshake inicial
	hs, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read handshake: %v", err)
	}
	off, _ := hs.Reader.ReadInt16()
	nc, _ := crypto.NewNetCrypto(off)
	client := &mockClient{conn: conn, crypto: nc}

	// 1. Probar Heartbeat (0x0804 -> 0x0805)
	hbReq := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HEARTBEAT_REQ)
	if err := client.sendPacket(hbReq); err != nil {
		t.Fatalf("Failed to send heartbeat req: %v", err)
	}

	hbAck, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("Failed to read heartbeat ack: %v", err)
	}
	if hbAck.OpCode != protocol.PROTO_NC_MISC_HEARTBEAT_ACK {
		t.Fatalf("Expected OpCode 0x%04X (PROTO_NC_MISC_HEARTBEAT_ACK), got 0x%04X",
			protocol.PROTO_NC_MISC_HEARTBEAT_ACK, hbAck.OpCode)
	}
	t.Logf("Heartbeat verificado exitosamente (OpCode: 0x%04X)", hbAck.OpCode)

	// 2. Probar Normal Logout (0x0C18)
	logoutCmd := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	if err := client.sendPacket(logoutCmd); err != nil {
		t.Fatalf("Failed to send logout cmd: %v", err)
	}

	// El servidor debe cerrar la conexión
	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, err = readServerPacket(conn)
	if err == nil {
		t.Errorf("Expected connection to be closed after PROTO_NC_USER_NORMALLOGOUT_CMD")
	} else {
		t.Logf("Desconexión limpia verificada tras PROTO_NC_USER_NORMALLOGOUT_CMD: %v", err)
	}
}

func TestLoginService_MultiRegion_OTP_WrongVersion_DynamicWorlds(t *testing.T) {
	cfg, err := config.Load("../../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		t.Fatalf("Failed to open MySQL: %v", err)
	}
	defer db.Close()

	repo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)

	// WorldManager con lectura dinámica desde MySQL tWorldServer
	wm := login.NewWorldManager(repo)

	activeVerKey, err := repo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("Failed to read FOUS.GameVersionKey from MySQL: %v", err)
	}

	testPort := 19012
	loginCfg := &login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", testPort),
		Port:     testPort,
	}

	srv := login.NewServer(loginCfg, repo, wm, transfers)
	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start login server: %v", err)
	}
	defer srv.Stop()

	// ------------------------------------------------------------------------
	// Test 1: Versión incorrecta rechazada con 0x0C02 y cierre de conexión
	// ------------------------------------------------------------------------
	t.Run("WrongClientVersion_0x0C02", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		hs, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read handshake: %v", err)
		}
		off, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(off)
		client := &mockClient{conn: conn, crypto: nc}

		// Enviar versión errónea
		wrongVerPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		wrongVerPkt.WriteString("INVALID_VER_999", 64)
		if err := client.sendPacket(wrongVerPkt); err != nil {
			t.Fatalf("Failed to send version pkt: %v", err)
		}

		// Debe responder 0x0C02 (PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK)
		ack, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read version ack: %v", err)
		}
		if ack.OpCode != protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK {
			t.Fatalf("Expected OpCode 0x%04X (WRONGVERSION_CHECK_ACK), got 0x%04X",
				protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK, ack.OpCode)
		}
		t.Logf("Escenario Versión Errónea Verificado: Servidor respondió 0x%04X", ack.OpCode)

		// Debe cerrar conexión inmediatamente
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, err = readServerPacket(conn)
		if err == nil {
			t.Errorf("Expected EOF/close after wrong version, but connection remained open")
		}
	})

	// ------------------------------------------------------------------------
	// Test 2: Login Estándar Internacional (0x0C06, 272B) y lectura dinámica de mundos
	// ------------------------------------------------------------------------
	t.Run("StandardLogin_0x0C06", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		hs, _ := readServerPacket(conn)
		off, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(off)
		client := &mockClient{conn: conn, crypto: nc}

		// Version Check correcto
		verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		verPkt.WriteString(activeVerKey, 64)
		_ = client.sendPacket(verPkt)
		_, _ = readServerPacket(conn) // 0x0C67

		// Enviar PROTO_NC_USER_LOGIN_REQ (0x0C06, 272B: User[256], Password[16])
		loginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGIN_REQ)
		loginReq.WriteString("admin", 256)
		loginReq.WriteString("admin", 16)
		_ = client.sendPacket(loginReq)

		loginAck, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read login ack: %v", err)
		}
		if loginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
			t.Fatalf("Expected OpCode 0x%04X (LOGIN_ACK), got 0x%04X",
				protocol.PROTO_NC_USER_LOGIN_ACK, loginAck.OpCode)
		}

		numWorlds, _ := loginAck.Reader.ReadByte()
		if numWorlds == 0 {
			t.Errorf("Expected at least 1 world from MySQL tWorldServer, got 0")
		}
		t.Logf("Login Estándar 0x0C06 exitoso. Mundos dinámicos recibidos: %d", numWorlds)
	})

	// ------------------------------------------------------------------------
	// Test 3: Login con OTP (0x0C37, 32B) y emisión de token de mundo (64B)
	// ------------------------------------------------------------------------
	t.Run("LoginWithOTP_0x0C37", func(t *testing.T) {
		// Esperar a que la desconexión del test anterior finalice
		time.Sleep(50 * time.Millisecond)

		// Generar un OTP de 32 bytes para admin (UserNo = 1)
		var testOTP [32]byte
		copy(testOTP[:], "0123456789abcdef0123456789abcdef")

		err := repo.CreateOTP(1, testOTP)
		if err != nil {
			t.Fatalf("Failed to create test OTP: %v", err)
		}

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		hs, _ := readServerPacket(conn)
		off, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(off)
		client := &mockClient{conn: conn, crypto: nc}

		// Version Check
		verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		verPkt.WriteString(activeVerKey, 64)
		_ = client.sendPacket(verPkt)
		_, _ = readServerPacket(conn) // 0x0C67

		// Enviar PROTO_NC_USER_LOGIN_WITH_OTP_REQ (0x0C37, 32B)
		otpReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ)
		otpReq.WriteBytes(testOTP[:])
		_ = client.sendPacket(otpReq)

		loginAck, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read login ack via OTP: %v", err)
		}
		if loginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
			t.Fatalf("Expected OpCode 0x%04X (LOGIN_ACK via OTP), got 0x%04X",
				protocol.PROTO_NC_USER_LOGIN_ACK, loginAck.OpCode)
		}
		t.Logf("Login con OTP (0x0C37) exitoso")

		// Seleccionar mundo 0 (Xelion)
		selPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
		selPkt.WriteByte(0)
		_ = client.sendPacket(selPkt)

		selAck, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read worldselect ack: %v", err)
		}
		if selAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
			t.Fatalf("Expected OpCode 0x%04X (WORLDSELECT_ACK), got 0x%04X",
				protocol.PROTO_NC_USER_WORLDSELECT_ACK, selAck.OpCode)
		}

		status, _ := selAck.Reader.ReadByte()
		ip, _ := selAck.Reader.ReadString(16)
		port, _ := selAck.Reader.ReadUint16()
		tokenBytes, _ := selAck.Reader.ReadBytes(64)

		if len(tokenBytes) != 64 {
			t.Errorf("Expected 64 bytes token in WORLDSELECT_ACK, got %d bytes", len(tokenBytes))
		}
		t.Logf("WORLDSELECT_ACK verificado con token criptográfico de 64B: Status=%d, %s:%d", status, ip, port)
	})
}
