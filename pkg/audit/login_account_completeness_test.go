package audit_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

// =============================================================================
// CONTROL 6 (Sección 9.6): Verificación Real de los 49 Archivos .shn + tHackScanLog
// =============================================================================
func TestFase07_Control6_SHNChecksumsAndHackScanLog(t *testing.T) {
	// 1. Localizar la carpeta oficial de datos SHN
	candidatePaths := []string{
		"../AstraServer/Server/9Data/Shine",
		"../../AstraServer/Server/9Data/Shine",
		"/Users/fedemac/Desktop/Xelion Online/AstraServer/Server/9Data/Shine",
	}

	var shnDir string
	for _, p := range candidatePaths {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			if _, err := os.Stat(filepath.Join(p, "Abstate.shn")); err == nil {
				shnDir = p
				break
			}
		}
	}

	if shnDir == "" {
		t.Fatalf("No se encontró el directorio oficial de SHNs en ninguna de las rutas: %v", candidatePaths)
	}

	// 2. Instanciar ChecksumEngine oficial con los 49 archivos
	engine, err := shn.NewChecksumEngine(shnDir)
	if err != nil {
		t.Fatalf("Fallo inicializando ChecksumEngine con '%s': %v", shnDir, err)
	}

	// Comprobar que los 49 archivos generaron su hash MD5 Name8 de 32 bytes
	officialChecksums := engine.GetChecksums()
	for i := 0; i < 49; i++ {
		rawHash := engine.GetRawHash(i)
		fileName := engine.GetFileName(i)

		if len(rawHash) != 32 {
			t.Errorf("Índice %d (%s): hash MD5 tiene longitud %d, esperada 32", i, fileName, len(rawHash))
		}
		if rawHash == strings.Repeat("0", 32) {
			t.Errorf("Índice %d (%s): hash no generado (todos ceros)", i, fileName)
		}
	}
	t.Logf("Control 6.1: Verificados satisfactoriamente los 49 hashes oficiales .shn (índices 0..48)")

	// 3. Verificar que la comparación oficial contra hashes idénticos es válida
	valid, wrongIdx := engine.VerifyDataFiles(officialChecksums)
	if !valid {
		t.Fatalf("VerifyDataFiles reportó inválido contra los propios hashes oficiales (wrongIdx=%d)", wrongIdx)
	}

	// 4. Simular alteración intencional de 1 byte en el archivo 2 (CharacterTitleData.shn)
	corruptedChecksums := officialChecksums
	corruptedChecksums[2][0] ^= 0x55 // Flip de bits

	validCorrupted, wrongCorruptedIdx := engine.VerifyDataFiles(corruptedChecksums)
	if validCorrupted {
		t.Fatalf("VerifyDataFiles debió rechazar el checksum corrupto en el índice 2, pero retornó válido")
	}
	if wrongCorruptedIdx != 2 {
		t.Fatalf("VerifyDataFiles reportó índice incorrecto: esperado 2 (CharacterTitleData.shn), obtenido %d", wrongCorruptedIdx)
	}
	t.Logf("Control 6.2: Detección de alteración simulada exitosa en índice 2 (%s)", engine.GetFileName(2))

	// 5. Persistir el incidente en MySQL tHackScanLog mediante CALL p_HackScan_Add
	db := getTestDB(t)
	defer db.Close()

	testUserNo := 1
	testWorldNo := uint8(0)
	testCharNo := 100
	testDesc := fmt.Sprintf("SHN mismatch idx %d (%s)", wrongCorruptedIdx, engine.GetFileName(wrongCorruptedIdx))

	_, err = db.Exec("CALL p_HackScan_Add(?, ?, ?, 1, ?)", testUserNo, testCharNo, testWorldNo, testDesc)
	if err != nil {
		t.Fatalf("Fallo ejecutando CALL p_HackScan_Add en MySQL: %v", err)
	}

	// 6. Verificar que la fila quedó persistida correctamente en tHackScanLog
	var (
		dbUserNo  int
		dbWorldNo uint8
		dbCharNo  int
		dbDetail  string
	)
	query := `
		SELECT nUserNo, nWorldNo, nCharNo, sDetail
		FROM tHackScanLog
		WHERE sDetail = ?
		ORDER BY dDate DESC
		LIMIT 1
	`
	err = db.QueryRow(query, testDesc).Scan(&dbUserNo, &dbWorldNo, &dbCharNo, &dbDetail)
	if err != nil {
		t.Fatalf("No se encontró el registro de alteración en tHackScanLog: %v", err)
	}
	if dbUserNo != testUserNo || dbWorldNo != testWorldNo || dbCharNo != testCharNo || dbDetail != testDesc {
		t.Errorf("Discrepancia en tHackScanLog: esperado (%d, %d, %d, '%s'), obtenido (%d, %d, %d, '%s')",
			testUserNo, testWorldNo, testCharNo, testDesc, dbUserNo, dbWorldNo, dbCharNo, dbDetail)
	}

	t.Logf("Control 6 PASS: 49 SHNs de ChecksumEngine y persistencia en tHackScanLog verificados al 100%%")
}

// =============================================================================
// CONTROL 7 (Sección 9.7): Prueba E2E Integral con Verificación Forense en MySQL
// =============================================================================
func TestFase07_Control7_FullE2E_ForensicMySQL(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	accountRepo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)
	wm := login.NewWorldManager(accountRepo)

	loginPort := 19020
	loginCfg := &login.Config{
		BindAddr:      fmt.Sprintf("127.0.0.1:%d", loginPort),
		Port:          loginPort,
		ClientVersion: "10022024000000",
	}

	// Sincronizar versión oficial en FOUS
	_ = accountRepo.SyncClientVersionKey(loginCfg.ClientVersion)

	loginSrv := login.NewServer(loginCfg, accountRepo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		t.Fatalf("Fallo iniciando LoginServer: %v", err)
	}
	defer loginSrv.Stop()

	worldPort := 19024
	worldCfg := &world.Config{
		ID:       0,
		Name:     "Xelion",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: 19224,
	}

	charRepo := database.NewMySQLCharacterRepository(db)
	worldSrv := world.NewServer(worldCfg, charRepo, transfers, nil, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("Fallo iniciando WorldServer: %v", err)
	}
	defer worldSrv.Stop()

	time.Sleep(50 * time.Millisecond)

	// -------------------------------------------------------------------------
	// 7.1: Versión errónea rechazada con 0x0C02 y desconexión
	// -------------------------------------------------------------------------
	t.Run("7.1_WrongVersion_0x0C02", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
		if err != nil {
			t.Fatalf("Fallo conectando a LoginServer: %v", err)
		}
		defer conn.Close()

		hs, err := auditReadServerPacket(conn)
		if err != nil || hs.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
			t.Fatalf("Fallo recibiendo seed handshake: %v", err)
		}
		seed, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(seed)
		client := &auditMockClient{conn: conn, crypto: nc}

		req := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		req.WriteString("WRONG_VERSION_KEY", 64)
		_ = client.sendPacket(req)

		ack, err := auditReadServerPacket(conn)
		if err != nil {
			t.Fatalf("Fallo leyendo respuesta de versión: %v", err)
		}
		if ack.OpCode != protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK {
			t.Fatalf("Esperado OpCode 0x%04X (WRONGVERSION), obtenido 0x%04X",
				protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK, ack.OpCode)
		}
		t.Logf("Paso 7.1 PASS: Versión errónea rechazada oficialmente con 0x0C02")
	})

	// -------------------------------------------------------------------------
	// 7.2: Handshake oficial y login exitoso preservando tUser.dDate intacto
	// -------------------------------------------------------------------------
	var worldToken [64]byte

	t.Run("7.2_Login_PreservesDate_And_WorldToken", func(t *testing.T) {
		// Capturar dDate original de 'admin' en MySQL antes del login
		var initialDDate time.Time
		err := db.QueryRow("SELECT dDate FROM tUser WHERE sUserID = 'admin'").Scan(&initialDDate)
		if err != nil {
			t.Fatalf("Error consultando dDate inicial de 'admin': %v", err)
		}

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
		if err != nil {
			t.Fatalf("Fallo conectando: %v", err)
		}
		defer conn.Close()

		hs, _ := auditReadServerPacket(conn)
		seed, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(seed)
		client := &auditMockClient{conn: conn, crypto: nc}

		// 1. Version check válido
		verReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		verReq.WriteString("10022024000000", 64)
		_ = client.sendPacket(verReq)

		rightAck, _ := auditReadServerPacket(conn) // 0x0C67
		if rightAck.OpCode != protocol.PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK {
			t.Fatalf("Esperado 0x0C67, obtenido 0x%04X", rightAck.OpCode)
		}
		_, _ = rightAck.Reader.ReadByte()    // bXTrapCheck
		_, _ = auditReadServerPacket(conn) // 0x0807 Seed ACK

		// 2. XTrap Check
		xtrapReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_XTRAP_REQ)
		xtrapReq.WriteByte(1)
		_ = client.sendPacket(xtrapReq)
		xtrapAck, _ := auditReadServerPacket(conn)
		if xtrapAck.OpCode != protocol.PROTO_NC_USER_XTRAP_ACK {
			t.Fatalf("Esperado 0x0C05 (XTRAP_ACK), obtenido 0x%04X", xtrapAck.OpCode)
		}

		// 3. Login US oficial (0x0C5A, 316B)
		loginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
		loginReq.WriteString("admin", 260)
		loginReq.WriteString("admin", 36)
		loginReq.Fill(20, 0)
		_ = client.sendPacket(loginReq)

		loginAck, err := auditReadServerPacket(conn)
		if err != nil {
			t.Fatalf("Fallo leyendo respuesta de login: %v", err)
		}
		if loginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
			t.Fatalf("Esperado 0x0C0A (LOGIN_ACK), obtenido 0x%04X", loginAck.OpCode)
		}

		// 4. Verificación forense en MySQL: tUser.dDate NO fue modificado
		var postLoginDDate time.Time
		var postLoginIP string
		err = db.QueryRow("SELECT dDate, sUserIP FROM tUser WHERE sUserID = 'admin'").
			Scan(&postLoginDDate, &postLoginIP)
		if err != nil {
			t.Fatalf("Error consultando tUser tras login: %v", err)
		}

		if !initialDDate.Equal(postLoginDDate) {
			t.Errorf("VIOLACIÓN FORENSE (Bug #2): tUser.dDate fue alterado! Inicial: %v, Posterior: %v",
				initialDDate, postLoginDDate)
		}
		if postLoginIP != "127.0.0.1" {
			t.Errorf("tUser.sUserIP no se actualizó con la IP de conexión: obtenido '%s'", postLoginIP)
		}

		// 5. Verificación de auditoría en tAccountLog (nType = 1 éxito)
		var logCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM tAccountLog WHERE sUserID = 'admin' AND nType = 1").Scan(&logCount)
		if logCount == 0 {
			t.Errorf("No se registró el ingreso oficial (nType = 1) en tAccountLog")
		}

		// 6. Selección de Mundo 0 (0x0C0B -> 0x0C0C con token de 64 bytes)
		selPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
		selPkt.WriteByte(0)
		_ = client.sendPacket(selPkt)

		selAck, err := auditReadServerPacket(conn)
		if err != nil {
			t.Fatalf("Fallo leyendo WORLDSELECT_ACK: %v", err)
		}
		if selAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
			t.Fatalf("Esperado 0x0C0C, obtenido 0x%04X", selAck.OpCode)
		}

		_, _ = selAck.Reader.ReadByte()      // status
		_, _ = selAck.Reader.ReadString(16)  // ip
		_, _ = selAck.Reader.ReadUint16()     // port
		tokenBytes, err := selAck.Reader.ReadBytes(64)
		if err != nil || len(tokenBytes) != 64 {
			t.Fatalf("WORLDSELECT_ACK debe contener token criptográfico exacto de 64B, obtenido %d bytes", len(tokenBytes))
		}
		copy(worldToken[:], tokenBytes)

		t.Logf("Paso 7.2 PASS: Login preservó dDate intacto y emitió token de mundo de 64 bytes")
	})

	// -------------------------------------------------------------------------
	// 7.3: Consumo de token 64B en WorldService, emisión de OTP 32B y Re-login
	// -------------------------------------------------------------------------
	t.Run("7.3_WorldLogin_Consume64B_And_OTP32B", func(t *testing.T) {
		worldConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
		if err != nil {
			t.Fatalf("Fallo conectando a WorldServer: %v", err)
		}
		defer worldConn.Close()

		hs, _ := auditReadServerPacket(worldConn)
		seed, _ := hs.Reader.ReadInt16()
		nc, _ := crypto.NewNetCrypto(seed)
		worldClient := &auditMockClient{conn: worldConn, crypto: nc}

		// 1. Enviar PROTO_NC_USER_LOGINWORLD_REQ oficial (320 bytes payload: user[256] + validate_new[64])
		keyReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
		keyReq.WriteString("admin", 256)
		keyReq.WriteBytes(worldToken[:])
		_ = worldClient.sendPacket(keyReq)

		// 2. Recibir lista de personajes (PROTO_NC_CHAR_CLIENT_CHARLIST_ACK - 0x0C14)
		charListAck, err := auditReadServerPacket(worldConn)
		if err != nil {
			t.Fatalf("Fallo leyendo respuesta de WorldLogin: %v", err)
		}
		if charListAck.OpCode != protocol.PROTO_NC_CHAR_CLIENT_CHARLIST_ACK {
			t.Fatalf("Esperado 0x0C14 (CHARLIST_ACK), obtenido 0x%04X", charListAck.OpCode)
		}
		t.Logf("Traspaso a WorldServer validado mediante token de 64 bytes (usp_Session_ConsumeWorldToken)")

		// 3. Volver a la selección de mundos solicitando sOTP de 32 bytes (0x0C33 -> 0x0C34)
		willSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WILL_WORLD_SELECT_REQ)
		_ = worldClient.sendPacket(willSelReq)

		willSelAck, err := auditReadServerPacket(worldConn)
		if err != nil {
			t.Fatalf("Fallo leyendo WILL_WORLD_SELECT_ACK: %v", err)
		}
		if willSelAck.OpCode != protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK {
			t.Fatalf("Esperado 0x0C34 (WILL_WORLD_SELECT_ACK), obtenido 0x%04X", willSelAck.OpCode)
		}

		errCode, _ := willSelAck.Reader.ReadUint16()
		if errCode != 0 {
			t.Fatalf("WILL_WORLD_SELECT_ACK reportó error: %d", errCode)
		}
		otpBytes, err := willSelAck.Reader.ReadBytes(32)
		if err != nil || len(otpBytes) != 32 {
			t.Fatalf("WILL_WORLD_SELECT_ACK debe contener un sOTP de 32 bytes exactos, obtenido %d", len(otpBytes))
		}
		t.Logf("WorldServer emitió sOTP criptográfico de 32 bytes: '%s'", string(otpBytes[:16]))

		// 4. Reconexión a LoginService autenticando mediante sOTP (0x0C37, 32 bytes)
		loginConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
		if err != nil {
			t.Fatalf("Fallo conectando a LoginServer: %v", err)
		}
		defer loginConn.Close()

		hs2, _ := auditReadServerPacket(loginConn)
		seed2, _ := hs2.Reader.ReadInt16()
		nc2, _ := crypto.NewNetCrypto(seed2)
		loginClient := &auditMockClient{conn: loginConn, crypto: nc2}

		// Version check
		verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		verPkt.WriteString("10022024000000", 64)
		_ = loginClient.sendPacket(verPkt)
		_, _ = auditReadServerPacket(loginConn) // 0x0C67
		_, _ = auditReadServerPacket(loginConn) // 0x0807

		// Enviar OTP
		otpLoginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ)
		otpLoginReq.WriteBytes(otpBytes)
		_ = loginClient.sendPacket(otpLoginReq)

		otpLoginAck, err := auditReadServerPacket(loginConn)
		if err != nil {
			t.Fatalf("Fallo leyendo respuesta de OTP Login: %v", err)
		}
		if otpLoginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
			t.Fatalf("Reconexión con OTP debió autenticar (0x0C0A), obtenido 0x%04X", otpLoginAck.OpCode)
		}
		t.Logf("Reconexión con sOTP exitosa (0x0C37 -> 0x0C0A)")

		// 5. Reintento inmediato con el mismo OTP -> debe fallar con 0x0C09 (un solo uso)
		loginConn2, _ := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
		defer loginConn2.Close()
		hs3, _ := auditReadServerPacket(loginConn2)
		seed3, _ := hs3.Reader.ReadInt16()
		nc3, _ := crypto.NewNetCrypto(seed3)
		loginClient2 := &auditMockClient{conn: loginConn2, crypto: nc3}

		verPkt2 := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
		verPkt2.WriteString("10022024000000", 64)
		_ = loginClient2.sendPacket(verPkt2)
		_, _ = auditReadServerPacket(loginConn2)
		_, _ = auditReadServerPacket(loginConn2)

		reusePkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ)
		reusePkt.WriteBytes(otpBytes)
		_ = loginClient2.sendPacket(reusePkt)

		failAck, _ := auditReadServerPacket(loginConn2)
		if failAck.OpCode != protocol.PROTO_NC_USER_LOGINFAIL_ACK {
			t.Fatalf("Reutilización de OTP debió ser rechazada con 0x0C09, obtenido 0x%04X", failAck.OpCode)
		}
		t.Logf("Paso 7.3 PASS: OTP consumido de un solo uso en MySQL verificado")
	})

	// -------------------------------------------------------------------------
	// 7.4: Transacciones operativas de Cash, ChargeItem, ReturnUser y Logout
	// -------------------------------------------------------------------------
	t.Run("7.4_Transversal_Cash_Shop_Return_Logout", func(t *testing.T) {
		testUserNo := 1

		// Cash: status = 0
		err := accountRepo.AddCash(testUserNo, 0, 500, "AUDIT_REWARD")
		if err != nil {
			t.Fatalf("AddCash falló: %v", err)
		}
		cash, err := accountRepo.GetCash(testUserNo, 0)
		if err != nil || cash < 500 {
			t.Fatalf("GetCash reportó monto inválido: %d (err: %v)", cash, err)
		}
		err = accountRepo.UseCash(testUserNo, 0, 100, "AUDIT_BUY")
		if err != nil {
			t.Fatalf("UseCash falló: %v", err)
		}

		// ReturnUser
		err = accountRepo.SetReturnUserItem(testUserNo, 1)
		if err != nil {
			t.Fatalf("SetReturnUserItem falló: %v", err)
		}

		// ChargeItem
		_, err = accountRepo.InsertChargedItem(&database.ChargedItem{
			UserNo:  testUserNo,
			GoodsNo: 1001,
			Amount:  5,
		})
		if err != nil {
			t.Fatalf("InsertChargedItem falló: %v", err)
		}
		items, err := accountRepo.GetChargedItems(testUserNo)
		if err != nil || len(items) == 0 {
			t.Fatalf("GetChargedItems no retornó items: %v", err)
		}
		err = accountRepo.DrawChargedItem(int64(items[0].RowNo))
		if err != nil {
			t.Fatalf("DrawChargedItem falló: %v", err)
		}

		// Logout con playMin real
		err = accountRepo.LogLogout(0, testUserNo, 1, 1, 10, 5000, 15)
		if err != nil {
			t.Fatalf("LogLogout falló: %v", err)
		}
		playMin, err := accountRepo.CloseSessionAndGetPlayMin(testUserNo)
		if err != nil {
			t.Fatalf("CloseSessionAndGetPlayMin falló: %v", err)
		}
		_ = accountRepo.LogAccount("admin", 2, testUserNo, playMin, "", 0, 0, "127.0.0.1")
		_ = accountRepo.AddSyslog(0xC803, testUserNo, "admin", "127.0.0.1", 0, 0, 0)

		t.Logf("Paso 7.4 PASS: Operaciones transversales (Cash status=0, ChargeItem, ReturnUser, Logout) exitosas")
	})

	t.Logf("Control 7 PASS: E2E Integral con Verificación Forense en MySQL verificado al 100%%")
}

// Helpers para sockets en test E2E
type auditMockClient struct {
	conn   net.Conn
	crypto *crypto.NetCrypto
}

func auditReadServerPacket(conn net.Conn) (*protocol.Packet, error) {
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
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

func (c *auditMockClient) sendPacket(pw *protocol.PacketWriter) error {
	bufCopy := make([]byte, len(pw.Data()))
	copy(bufCopy, pw.Data())
	if c.crypto != nil {
		c.crypto.Crypt(bufCopy)
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

	_, err := c.conn.Write(wire)
	return err
}
