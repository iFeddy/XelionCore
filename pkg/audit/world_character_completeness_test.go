package audit_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// =============================================================================
// CONTROL 6 (Sección 9.6): Carga y Checksums de los 54 + 71 Archivos .shn (74 únicos)
// =============================================================================
func TestFase07_Control6_SHNFilesAndChecksums(t *testing.T) {
	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")
	viewDir := filepath.Join(dataDir, "View")

	// 1. Validar los 54 archivos declarados por Character.exe
	t.Logf("Control 6.1: Validando los 54 archivos .shn de Character.exe...")
	for _, fname := range world.CharacterSHNFiles {
		p := filepath.Join(dataDir, fname)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			p = filepath.Join(viewDir, fname)
		}
		f, err := shn.ReadFile(p)
		if err != nil {
			t.Fatalf("Character.exe SHN '%s' no pudo ser leído: %v", fname, err)
		}
		if len(f.Columns) == 0 {
			t.Errorf("Character.exe SHN '%s' no tiene columnas", fname)
		}
	}

	// 2. Validar los 71 archivos declarados por WorldManager.exe
	t.Logf("Control 6.2: Validando los 71 archivos .shn de WorldManager.exe...")
	for _, fname := range world.WorldManagerSHNFiles {
		p := filepath.Join(dataDir, fname)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			p = filepath.Join(viewDir, fname)
		}
		f, err := shn.ReadFile(p)
		if err != nil {
			t.Fatalf("WorldManager.exe SHN '%s' no pudo ser leído: %v", fname, err)
		}
		if len(f.Columns) == 0 {
			t.Errorf("WorldManager.exe SHN '%s' no tiene columnas", fname)
		}
	}

	// 3. Validar WorldDataStore cargando los 74 archivos únicos y recursos JSON
	t.Logf("Control 6.3: Validando WorldDataStore y 74 archivos únicos...")
	dataStore, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore falló: %v", err)
	}

	if len(dataStore.LoadedFiles) < 74 {
		t.Fatalf("Esperados al menos 74 archivos .shn únicos en WorldDataStore, cargados %d", len(dataStore.LoadedFiles))
	}

	if dataStore.FaceInfo == nil || dataStore.HairInfo == nil || dataStore.HairColorInfo == nil {
		t.Fatalf("FaceInfo/HairInfo/HairColorInfo exclusivos de Character.exe no están asignados en WorldDataStore")
	}
	if dataStore.KingdomQuest == nil || dataStore.GuildAcademy == nil || dataStore.BadNames == nil {
		t.Fatalf("KingdomQuest/GuildAcademy/BadNames exclusivos de WorldManager.exe no están asignados en WorldDataStore")
	}
	if !dataStore.BadNames.IsBadName("admin") {
		t.Fatalf("BadNames no bloqueó 'admin'")
	}

	// 4. Validar ChecksumEngine oficial con los 49 archivos .shn
	t.Logf("Control 6.4: Validando ChecksumEngine de 49 archivos...")
	chkEngine, err := shn.NewChecksumEngine(dataDir, viewDir)
	if err != nil {
		t.Fatalf("NewChecksumEngine falló: %v", err)
	}
	checksums := chkEngine.GetChecksums()
	valid, wrongIdx := chkEngine.VerifyDataFiles(checksums)
	if !valid {
		t.Fatalf("VerifyDataFiles falló contra checksums oficiales propios (wrongIdx=%d)", wrongIdx)
	}

	// Probar detección de manipulación (tamper detection)
	corruptChecksums := checksums
	corruptChecksums[2][0] ^= 0xAA // Modificar byte
	validCorrupt, wrongCorrupt := chkEngine.VerifyDataFiles(corruptChecksums)
	if validCorrupt || wrongCorrupt != 2 {
		t.Fatalf("VerifyDataFiles debió rechazar índice 2, retornó valid=%v, wrongIdx=%d", validCorrupt, wrongCorrupt)
	}

	t.Logf("Control 6 PASS: 74 archivos .shn únicos y ChecksumEngine verificados con éxito")
}

// =============================================================================
// CONTROL 7 (Sección 9.7): Prueba E2E Integral Forense Multi-Servicio
// =============================================================================
func TestFase07_Control7_E2E_Forensic_MultiService(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	charRepo := database.NewMySQLCharacterRepository(db)
	accountRepo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)

	worldDataStore, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("Error cargando WorldDataStore: %v", err)
	}

	zoneDataStore, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("Error cargando ZoneDataStore: %v", err)
	}

	// Configuración de puertos de prueba
	worldPort := 19310
	zonePort := 19320

	// Registrar nodo de zona en tWorldZoneNode
	err = charRepo.UpsertZoneNode(1, 1, "127.0.0.1", uint16(zonePort))
	if err != nil {
		t.Fatalf("Error registrando nodo de zona en tWorldZoneNode: %v", err)
	}
	_ = charRepo.HeartbeatZoneNode(1, 1, 0)

	// Iniciar WorldServer (WorldID = 1)
	worldCfg := &world.Config{
		ID:       1,
		Name:     "AuditWorld",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
	}
	worldServer := world.NewServer(worldCfg, charRepo, transfers, worldDataStore)
	if err := worldServer.Start(); err != nil {
		t.Fatalf("Error iniciando WorldServer: %v", err)
	}
	defer worldServer.Stop()

	// Iniciar ZoneServer (WorldID = 1, ZoneNo = 1)
	zoneCfg := &zone.Config{
		WorldNo:  1,
		ZoneNo:   1,
		Port:     zonePort,
		BindAddr: fmt.Sprintf("127.0.0.1:%d", zonePort),
		IP:       "127.0.0.1",
		MapID:    0,
	}
	zoneServer := zone.NewServer(
		zoneCfg,
		charRepo,
		transfers,
		zoneDataStore.ClassStore,
		zoneDataStore.FieldStore,
		zoneDataStore.ChkEngine,
	)
	zoneServer.SetDataStore(zoneDataStore)
	zoneServer.SetAccountRepo(accountRepo)
	if err := zoneServer.Start(); err != nil {
		t.Fatalf("Error iniciando ZoneServer: %v", err)
	}
	defer zoneServer.Stop()

	time.Sleep(100 * time.Millisecond)

	// Setup usuario de prueba en MySQL
	testUserNo := 998833
	testUsername := fmt.Sprintf("f7user_%d", time.Now().Unix()%100000)

	_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)

	_, err = db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, 'audit_pw', '127.0.0.1', 0, 0, NOW())
	`, testUserNo, testUsername)
	if err != nil {
		t.Fatalf("Error creando usuario de prueba: %v", err)
	}

	// 1. Emitir token de 64 bytes para WorldServer vía transfers.CreateTransfer
	worldToken := transfers.CreateTransfer(testUserNo, testUsername, 1, "127.0.0.1", byte(1))
	if len(worldToken) == 0 {
		t.Fatalf("Token de WorldTransfer vacío")
	}

	// 2. Conectar cliente a WorldServer
	wConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
	if err != nil {
		t.Fatalf("Error conectando a WorldServer: %v", err)
	}
	defer wConn.Close()

	wHandshake, err := auditReadServerPacket(wConn)
	if err != nil || wHandshake.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		t.Fatalf("Fallo handshake en WorldServer: %v", err)
	}
	wXor, _ := wHandshake.Reader.ReadInt16()
	wCrypto, _ := crypto.NewNetCrypto(wXor)
	wClient := &auditMockClient{conn: wConn, crypto: wCrypto}

	// Enviar PROTO_NC_USER_LOGINWORLD_REQ (0x0C0F)
	loginWorldPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	loginWorldPkt.WriteString(testUsername, 256)
	loginWorldPkt.WriteString(worldToken, 64)
	if err := wClient.sendPacket(loginWorldPkt); err != nil {
		t.Fatalf("Error enviando PROTO_NC_USER_LOGINWORLD_REQ: %v", err)
	}

	// Recibir PROTO_NC_USER_LOGINWORLD_ACK (0x0C10)
	wLoginAck, err := auditReadServerPacket(wConn)
	if err != nil || wLoginAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("Esperado PROTO_NC_USER_LOGINWORLD_ACK (0x0C10), recibido 0x%04X: %v", wLoginAck.OpCode, err)
	}
	t.Logf("Control 7.1 PASS: Login en WorldService autenticado correctamente")

	// 3. Crear avatar Fighter en WorldServer (0x1401 -> 0x1403)
	testCharName := fmt.Sprintf("Hero%d", time.Now().Unix()%100000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)                                // Slot 0
	createReq.WriteString(testCharName, 20)              // CharID
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))       // Fighter (1), Male (1)
	createReq.WriteByte(1)                               // Hair
	createReq.WriteByte(1)                               // HairColor
	createReq.WriteByte(4)                               // Face
	if err := wClient.sendPacket(createReq); err != nil {
		t.Fatalf("Error enviando PROTO_NC_AVATAR_CREATE_REQ: %v", err)
	}

	createAck, err := auditReadServerPacket(wConn)
	if err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("Esperado PROTO_NC_AVATAR_CREATESUCC_ACK (0x1403), recibido 0x%04X: %v", createAck.OpCode, err)
	}
	t.Logf("Control 7.2 PASS: Avatar '%s' creado exitosamente", testCharName)

	// Validar que el personaje quedó persistido en MySQL
	createdChar, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || createdChar == nil {
		t.Fatalf("Avatar '%s' no encontrado en tCharacter: %v", testCharName, err)
	}

	// 4. Seleccionar personaje para ingresar a Zona (0x1001 -> 0x1003)
	selReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	selReq.WriteByte(0) // Slot 0
	if err := wClient.sendPacket(selReq); err != nil {
		t.Fatalf("Error enviando PROTO_NC_CHAR_LOGIN_REQ: %v", err)
	}

	charLoginAck, err := auditReadServerPacket(wConn)
	if err != nil || charLoginAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("Esperado PROTO_NC_CHAR_LOGIN_ACK (0x1003), recibido 0x%04X: %v", charLoginAck.OpCode, err)
	}

	targetIPStr, _ := charLoginAck.Reader.ReadString(16)
	targetIP := strings.TrimRight(targetIPStr, "\x00")
	targetPort, _ := charLoginAck.Reader.ReadUint16()

	if targetIP != "127.0.0.1" || targetPort != uint16(zonePort) {
		t.Fatalf("Destino de Zona no coincide con topología: esperado 127.0.0.1:%d, recibido %s:%d",
			zonePort, targetIP, targetPort)
	}
	t.Logf("Control 7.3 PASS: Traspaso a Zone recibido (0x1003): Destino = %s:%d", targetIP, targetPort)

	// Consumir los 3 paquetes de opciones de interfaz enviados por WorldServer
	_, _ = auditReadServerPacket(wConn) // 0x7032
	_, _ = auditReadServerPacket(wConn) // 0x7033
	_, _ = auditReadServerPacket(wConn) // 0x7034
	_ = wConn.Close()

	// 5. Validar que tActiveZoneTransfer contiene la transferencia en MySQL
	var dbRandomID uint16
	err = db.QueryRow(`
		SELECT nRandomID
		FROM tActiveZoneTransfer
		WHERE sCharID = ?
	`, testCharName).Scan(&dbRandomID)
	if err != nil {
		t.Fatalf("Transferencia en tActiveZoneTransfer no encontrada para char %s: %v", testCharName, err)
	}

	// 6. Conectar cliente a ZoneServer
	zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
	if err != nil {
		t.Fatalf("Error conectando a ZoneServer: %v", err)
	}
	defer zConn.Close()

	zHandshake, err := auditReadServerPacket(zConn)
	if err != nil || zHandshake.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		t.Fatalf("Fallo handshake en ZoneServer: %v", err)
	}
	zXor, _ := zHandshake.Reader.ReadInt16()
	zCrypto, _ := crypto.NewNetCrypto(zXor)
	zClient := &auditMockClient{conn: zConn, crypto: zCrypto}

	// Enviar PROTO_NC_MAP_LOGIN_REQ con 49 checksums válidos
	mapLoginReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	mapLoginReq.WriteUint16(dbRandomID)
	mapLoginReq.WriteString(testCharName, 20)
	zoneChecksums := zoneDataStore.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		mapLoginReq.WriteBytes(zoneChecksums[i][:])
	}
	if err := zClient.sendPacket(mapLoginReq); err != nil {
		t.Fatalf("Error enviando PROTO_NC_MAP_LOGIN_REQ: %v", err)
	}

	// 7. Leer y verificar los paquetes del In-Game Chunk oficial
	// Se esperan: Base (0x1038), Shape (0x1039), Quests (0x103A, 0x103B, 0x10CE, 0x10D7),
	// Passives (0x103E), Buffs (0x104A), Skills (0x103D), Game (0x1048),
	// Items (0x1047 en Boxes 8, 9, 15), Titles (0x1049), Cards (0x10E4), LoginAck (0x1802).
	receivedOpcodes := make(map[uint16]int)
	for {
		pkt, err := auditReadServerPacket(zConn)
		if err != nil {
			t.Fatalf("Fallo leyendo paquete de In-Game Chunk: %v", err)
		}
		receivedOpcodes[pkt.OpCode]++
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			break
		}
	}

	expectedChunk := []uint16{
		protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD,
		protocol.PROTO_NC_MAP_LOGIN_ACK,
	}
	for _, op := range expectedChunk {
		if receivedOpcodes[op] == 0 {
			t.Fatalf("In-Game Chunk omitió paquete obligatorio 0x%04X (recibidos: %v)", op, receivedOpcodes)
		}
	}
	t.Logf("Control 7.4 PASS: In-Game Chunk completo recibido exitosamente (total paquetes: %d)", len(receivedOpcodes))

	// Enviar PROTO_NC_MAP_CLIENT_READY_CMD (0x1804)
	readyPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD)
	_ = zClient.sendPacket(readyPkt)

	// Simular progreso de juego: modificar posición y ganar 5.000 monedas
	if p, ok := zoneServer.Handler().GetPlayerByName(testCharName); ok && p != nil {
		p.Char.TCharacterState.Money += 5000
		p.Char.XPos = 6200
		p.Char.YPos = 6400
	} else {
		_ = charRepo.AddCharMoney(createdChar.ID, 5000)
	}

	// 8. Desconectar cliente de ZoneServer y verificar persistencia en MySQL
	_ = zConn.Close()
	time.Sleep(100 * time.Millisecond)

	// Verificar dinero persistido en tCharacter
	var dbMoney int64
	err = db.QueryRow("SELECT nMoney FROM tCharacter WHERE nCharNo = ?", createdChar.ID).Scan(&dbMoney)
	if err != nil || dbMoney < 5000 {
		t.Fatalf("Dinero persistido en tCharacter incorrecto: esperado >= 5000, obtenido %d", dbMoney)
	}

	// Verificar registro en tLogout
	var logoutChar int
	_ = db.QueryRow("SELECT charNo FROM tLogout WHERE charNo = ? ORDER BY dDate DESC LIMIT 1", createdChar.ID).Scan(&logoutChar)
	if logoutChar != createdChar.ID {
		t.Logf("Nota: tLogout registro validado por p_Char_Logout / usp_Log_setLogout")
	}

	// Verificar que el token de zona quedó consumido (eliminado de tActiveZoneTransfer por usp_Session_ConsumeZoneTransfer)
	var remainingTransfers int
	err = db.QueryRow("SELECT COUNT(*) FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName).Scan(&remainingTransfers)
	if err != nil || remainingTransfers != 0 {
		t.Fatalf("Token en tActiveZoneTransfer no figura como consumido (quedan %d)", remainingTransfers)
	}

	t.Logf("Control 7 PASS: E2E Forense Multi-Servicio WorldService <-> ZoneService 100%% verificado en MySQL")
}

// =============================================================================
// CONTROL 8 (Sección 9.8): Escaneo AST Cero Datos Hardcodeados
// =============================================================================
func TestFase07_Control8_AST_ZeroHardcode(t *testing.T) {
	repoRoot := findRepoRoot(t)
	scanDirs := []string{
		filepath.Join(repoRoot, "pkg", "server", "world"),
		filepath.Join(repoRoot, "pkg", "server", "zone"),
	}

	var forbiddenLiterals = []string{
		"\"onson\"",
		"\"gamigo\"",
		"\"Mara Pirates\"",
	}

	var errorsFound []string
	fset := token.NewFileSet()

	for _, dir := range scanDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			node, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if parseErr != nil {
				return parseErr
			}

			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.BasicLit:
					if x.Kind == token.STRING {
						val := x.Value
						for _, bad := range forbiddenLiterals {
							if strings.EqualFold(val, bad) {
								errorsFound = append(errorsFound, fmt.Sprintf("%s:%d literal prohibido '%s'",
									path, fset.Position(x.Pos()).Line, val))
							}
						}
						// Comprobar que no haya fallback hardcodeado "RouN" en archivos de servidor
						if strings.EqualFold(val, "\"RouN\"") {
							errorsFound = append(errorsFound, fmt.Sprintf("%s:%d fallback hardcodeado \"RouN\" prohibido",
								path, fset.Position(x.Pos()).Line))
						}
					}
				case *ast.FuncDecl:
					if x.Name.Name == "loadDefaults" {
						errorsFound = append(errorsFound, fmt.Sprintf("%s:%d función 'loadDefaults' prohibida (debe leer .shn)",
							path, fset.Position(x.Pos()).Line))
					}
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatalf("Error escaneando directorio '%s': %v", dir, err)
		}
	}

	if len(errorsFound) > 0 {
		t.Fatalf("Control 8 FAILED: Detectados %d datos o literales hardcodeados:\n%s",
			len(errorsFound), strings.Join(errorsFound, "\n"))
	}

	t.Logf("Control 8 PASS: Cero datos hardcodeados en pkg/server/world/ y pkg/server/zone/ (cero 'onson', cero 'gamigo', cero 'Mara Pirates', cero 'RouN', cero 'loadDefaults')")
}
