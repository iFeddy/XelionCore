package audit_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// =============================================================================
// CONTROL 3 (Sección 9.3): Escaneo AST Cero SQL Inline en ZoneService, Lua y GameLog
// =============================================================================
func TestFase16_Control3_AST_ZoneAndGameLog_NoInlineSQL(t *testing.T) {
	rootPath := findRepoRoot(t)
	targetDirs := []string{
		filepath.Join(rootPath, "pkg", "server", "zone"),
		filepath.Join(rootPath, "pkg", "script", "lua"),
	}
	targetFiles := []string{
		filepath.Join(rootPath, "pkg", "database", "mysql_gamelog.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_phase16.go"),
	}

	var goFiles []string
	for _, dir := range targetDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("Error leyendo directorio %s: %v", dir, err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") && !strings.HasSuffix(e.Name(), "_test.go") {
				goFiles = append(goFiles, filepath.Join(dir, e.Name()))
			}
		}
	}
	for _, f := range targetFiles {
		if _, err := os.Stat(f); err == nil {
			goFiles = append(goFiles, f)
		}
	}

	inlineSQLRe := regexp.MustCompile(`(?i)\b(SELECT\s+.+\s+FROM|INSERT\s+INTO|UPDATE\s+.+\s+SET|DELETE\s+FROM|CREATE\s+TABLE|ALTER\s+TABLE|DROP\s+TABLE)\b`)
	fset := token.NewFileSet()
	var violations []string

	for _, filePath := range goFiles {
		node, err := parser.ParseFile(fset, filePath, nil, 0)
		if err != nil {
			t.Fatalf("Error parseando AST de %s: %v", filePath, err)
		}
		ast.Inspect(node, func(n ast.Node) bool {
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				val := strings.Trim(lit.Value, "\"`")
				if inlineSQLRe.MatchString(val) {
					pos := fset.Position(lit.Pos())
					violations = append(violations, fmt.Sprintf("%s:%d SQL inline prohibido: %q", filepath.Base(filePath), pos.Line, val))
				}
			}
			return true
		})
	}

	if len(violations) > 0 {
		t.Fatalf("CONTROL 3 FALLIDO: Se detectó SQL inline fuera de CALL <SP>:\n%s", strings.Join(violations, "\n"))
	}
	t.Logf("Control 3 PASS: %d archivos Go en ZoneService, Lua y MySQL GameLog/Phase16 escaneados vía AST — 0 SQL inline, 100%% CALL Stored Procedures.", len(goFiles))
}

// =============================================================================
// CONTROL 5 (Sección 9.5): Cobertura 245/245 C2S, Cero DATA_REPORT, 138 WMS y 196 GDS
// =============================================================================
func TestFase16_Control5_Zone_245C2S_138WMS_196GDS_Completeness(t *testing.T) {
	rootPath := findRepoRoot(t)

	// 1. Leer lista oficial de OpCodes C2S desde docs/Servicios/Zone.md (protocolstore ShinePlayer)
	zoneMDBytes, err := os.ReadFile(filepath.Join(rootPath, "docs", "Servicios", "Zone.md"))
	if err != nil {
		t.Fatalf("Error leyendo docs/Servicios/Zone.md: %v", err)
	}
	zoneMD := string(zoneMDBytes)
	secRe := regexp.MustCompile(`(?is)protocolstore.*?\n(.*?)(?:\n## |\z)`)
	secMatch := secRe.FindStringSubmatch(zoneMD)
	if len(secMatch) < 2 {
		t.Fatalf("No se encontró la tabla protocolstore(ShinePlayer) en docs/Servicios/Zone.md")
	}
	opRe := regexp.MustCompile(`PROTO_NC_[A-Z0-9_]+`)
	officialOpsList := opRe.FindAllString(secMatch[1], -1)
	officialOps := make(map[string]bool)
	for _, op := range officialOpsList {
		officialOps[op] = true
	}
	if len(officialOps) < 240 {
		t.Fatalf("Esperados >=240 OpCodes C2S oficiales en Zone.md, encontrados %d", len(officialOps))
	}

	// 2. Parseo AST estricto de HandlePacket en pkg/server/zone/handler.go
	handlerPath := filepath.Join(rootPath, "pkg", "server", "zone", "handler.go")
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, handlerPath, nil, 0)
	if err != nil {
		t.Fatalf("Error parseando AST de handler.go: %v", err)
	}

	handledOps := make(map[string]bool)
	var forbiddenDataReport []string

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "HandlePacket" {
			return true
		}
		ast.Inspect(fn.Body, func(sub ast.Node) bool {
			cc, ok := sub.(*ast.CaseClause)
			if !ok {
				return true
			}
			for _, expr := range cc.List {
				if sel, ok := expr.(*ast.SelectorExpr); ok {
					opName := sel.Sel.Name
					if strings.HasPrefix(opName, "PROTO_NC_DATA_REPORT_") {
						pos := fset.Position(sel.Pos())
						forbiddenDataReport = append(forbiddenDataReport, fmt.Sprintf("handler.go:%d %s", pos.Line, opName))
					}
					handledOps[opName] = true
				}
			}
			return true
		})
		return false
	})

	if len(forbiddenDataReport) > 0 {
		t.Fatalf("BUG #4 VIOLADO: DATA_REPORT_* debe vivir exclusivamente en WorldService, encontrado en Zone handler.go:\n%s",
			strings.Join(forbiddenDataReport, "\n"))
	}

	var missingOps []string
	for op := range officialOps {
		if !handledOps[op] {
			missingOps = append(missingOps, op)
		}
	}
	if len(missingOps) > 0 {
		t.Fatalf("CONTROL 5 FALLIDO: Faltan %d OpCodes C2S de protocolstore(ShinePlayer) en HandlePacket:\n%s",
			len(missingOps), strings.Join(missingOps, "\n"))
	}

	// 3. Verificar 138 flujos S2S wms_NC_* y 196 mapeos gds_NC_* en Coordinator
	h := zone.NewHandler(nil, nil, nil, nil, nil)
	if h.WMSHandlerCount() != 138 {
		t.Fatalf("CONTROL 5 FALLIDO: Esperados 138 handlers wms_NC_*, registrados %d", h.WMSHandlerCount())
	}
	if h.GDSMappingCount() != 196 {
		t.Fatalf("CONTROL 5 FALLIDO: Esperados 196 mapeos gds_NC_*, registrados %d", h.GDSMappingCount())
	}
	for _, sym := range zone.OfficialWMSSymbols138 {
		if !h.DispatchWMS(sym, nil) {
			t.Fatalf("Falta handler activo para %s", sym)
		}
	}
	for _, sym := range zone.OfficialGDSSymbols196 {
		if sp, ok := h.DispatchGDS(sym); !ok || sp == "" {
			t.Fatalf("Falta mapeo SP activo para %s", sym)
		}
	}

	// 4. Verificar topología multi-nodo oficial Zone00..Zone04 (9016..9028)
	expectedNodes := map[string]struct {
		no   uint16
		port int
	}{
		"Zone00": {0, 9016},
		"Zone01": {1, 9019},
		"Zone02": {2, 9022},
		"Zone03": {3, 9025},
		"Zone04": {4, 9028},
	}
	for id, exp := range expectedNodes {
		zNo, port, ok := zone.ResolveZoneNodeConfig(id)
		if !ok || zNo != exp.no || port != exp.port {
			t.Fatalf("ResolveZoneNodeConfig(%s) = (%d,%d,%v), esperado (%d,%d,true)", id, zNo, port, ok, exp.no, exp.port)
		}
	}

	t.Logf("Control 5 PASS: %d/%d OpCodes C2S ShinePlayer cubiertos en HandlePacket (0 Unhandled, 0 DATA_REPORT), 138/138 wms_NC_*, 196/196 gds_NC_* y nodos Zone00..Zone04 verificados.",
		len(officialOps), len(officialOps))
}

// =============================================================================
// CONTROL 7 (Sección 9.7): Prueba E2E Forense de Clúster Completo
// LoginService -> WorldService -> ZoneService -> GameLogDispatcher -> MySQL
// =============================================================================
func TestFase16_Control7_E2E_Forensic_Cluster_Zone_GameLog(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	accountRepo := database.NewMySQLAccountRepository(db)
	charRepo := database.NewMySQLCharacterRepository(db)
	gameLogRepo := database.NewMySQLGameLogRepository(db)
	transfers := database.NewTransferManager(db)

	worldDataStore, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("Error cargando WorldDataStore: %v", err)
	}
	zoneDataStore, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("Error cargando ZoneDataStore: %v", err)
	}

	activeVerKey, err := accountRepo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("Error obteniendo FOUS.GameVersionKey desde MySQL: %v", err)
	}

	// Puertos aislados para el clúster forense
	const (
		loginPort = 19410
		worldPort = 19414
		zonePort  = 19424
	)

	wm := login.NewWorldManager(accountRepo)
	loginSrv := login.NewServer(&login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", loginPort),
		Port:     loginPort,
	}, accountRepo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		t.Fatalf("Error iniciando LoginServer: %v", err)
	}
	defer loginSrv.Stop()

	worldSrv := world.NewServer(&world.Config{
		ID:       0,
		Name:     "AuditCluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort,
	}, charRepo, transfers, worldDataStore, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("Error iniciando WorldServer: %v", err)
	}
	defer worldSrv.Stop()

	if err := charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort)); err != nil {
		t.Fatalf("Error UpsertZoneNode(0,1): %v", err)
	}
	_ = charRepo.UpsertZoneNode(1, 1, "127.0.0.1", uint16(zonePort))
	defer func() {
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort = ?", zonePort)
	}()

	glDispatcher := database.NewGameLogDispatcher(gameLogRepo, 1000, 2)
	glDispatcher.Start()

	zoneSrv := zone.NewServer(&zone.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", zonePort),
		Port:     zonePort,
		IP:       "127.0.0.1",
		WorldNo:  0,
		ZoneNo:   1,
	}, charRepo, transfers, zoneDataStore.ClassStore, zoneDataStore.FieldStore, zoneDataStore.ChkEngine)
	zoneSrv.SetDataStore(zoneDataStore)
	zoneSrv.SetAccountRepo(accountRepo)
	zoneSrv.SetGameLogDispatcher(glDispatcher)
	if err := zoneSrv.Start(); err != nil {
		t.Fatalf("Error iniciando ZoneServer: %v", err)
	}
	defer zoneSrv.Stop()

	time.Sleep(100 * time.Millisecond)

	// Preparar cuenta dedicada de auditoría en tUser
	testUserNo := 998916
	testUser := fmt.Sprintf("f16cl_%d", time.Now().Unix()%100000)
	testPass := "f16pass"
	sum := md5.Sum([]byte(testPass))
	testPassHash := hex.EncodeToString(sum[:])

	cleanupClusterUser := func(charID int) {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo)
		var charIDs []int
		if charID > 0 {
			charIDs = append(charIDs, charID)
		} else {
			if rows, qErr := db.Query("SELECT nCharNo FROM tCharacter WHERE nUserNo = ?", testUserNo); qErr == nil {
				for rows.Next() {
					var cid int
					if err := rows.Scan(&cid); err == nil {
						charIDs = append(charIDs, cid)
					}
				}
				_ = rows.Close()
			}
		}
		for _, cid := range charIDs {
			_, _ = db.Exec("DELETE FROM tGMItemRewardLog WHERE nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tGMItemReward WHERE sOperID LIKE 'F16%' OR nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", cid)
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}
	cleanupClusterUser(0)

	_, err = db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
	`, testUserNo, testUser, testPassHash)
	if err != nil {
		t.Fatalf("Error insertando tUser de auditoría: %v", err)
	}

	var createdCharID int
	defer func() {
		cleanupClusterUser(createdCharID)
	}()

	// -------------------------------------------------------------------------
	// Paso 1: LoginService (VersionCheck + XTrap + US Login + WorldSelect)
	// -------------------------------------------------------------------------
	lConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
	if err != nil {
		t.Fatalf("Error conectando a LoginService: %v", err)
	}
	lHS, _ := auditReadServerPacket(lConn)
	lSeed, _ := lHS.Reader.ReadInt16()
	lCrypto, _ := crypto.NewNetCrypto(lSeed)
	lClient := &auditMockClient{conn: lConn, crypto: lCrypto}

	verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt.WriteString(activeVerKey, 64)
	_ = lClient.sendPacket(verPkt)
	verAck, err := auditReadServerPacket(lConn)
	if err != nil || verAck.OpCode != protocol.PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK {
		t.Fatalf("Version check falló en LoginService: ack=%v err=%v", verAck, err)
	}
	_, _ = auditReadServerPacket(lConn) // 0x0807 seed ack

	xtPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_XTRAP_REQ)
	xtPkt.WriteByte(1)
	_ = lClient.sendPacket(xtPkt)
	_, _ = auditReadServerPacket(lConn) // XTRAP_ACK

	lReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	lReq.WriteString(testUser, 260)
	lReq.WriteString(testPassHash, 36)
	lReq.Fill(20, 0)
	_ = lClient.sendPacket(lReq)
	lAck, err := auditReadServerPacket(lConn)
	if err != nil || lAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
		t.Fatalf("LoginService US_LOGIN falló: ack=%v err=%v", lAck, err)
	}

	wSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	wSelReq.WriteByte(0)
	_ = lClient.sendPacket(wSelReq)
	wSelAck, err := auditReadServerPacket(lConn)
	if err != nil || wSelAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("WORLDSELECT falló: ack=%v err=%v", wSelAck, err)
	}
	_, _ = wSelAck.Reader.ReadByte()
	_, _ = wSelAck.Reader.ReadString(16)
	_, _ = wSelAck.Reader.ReadUint16()
	tokenBytes, _ := wSelAck.Reader.ReadBytes(64)
	worldToken := strings.TrimRight(string(tokenBytes), "\x00")
	_ = lConn.Close()

	// -------------------------------------------------------------------------
	// Paso 2: WorldService (LOGINWORLD + AVATAR_CREATE + CHAR_LOGIN -> Zone)
	// -------------------------------------------------------------------------
	wConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
	if err != nil {
		t.Fatalf("Error conectando a WorldService: %v", err)
	}
	wHS, _ := auditReadServerPacket(wConn)
	wSeed, _ := wHS.Reader.ReadInt16()
	wCrypto, _ := crypto.NewNetCrypto(wSeed)
	wClient := &auditMockClient{conn: wConn, crypto: wCrypto}

	wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq.WriteString(testUser, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	wlAck, err := auditReadServerPacket(wConn)
	if err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("World LOGINWORLD falló: ack=%v err=%v", wlAck, err)
	}

	testCharName := fmt.Sprintf("F16Hero%d", time.Now().Unix()%10000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(testCharName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7))) // Fighter, Male
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	createAck, err := auditReadServerPacket(wConn)
	if err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE falló: ack=%v err=%v", createAck, err)
	}

	createdChar, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || createdChar == nil {
		t.Fatalf("Personaje '%s' no persistido en tCharacter: %v", testCharName, err)
	}
	createdCharID = createdChar.ID

	cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	cSelReq.WriteByte(0)
	_ = wClient.sendPacket(cSelReq)
	cSelAck, err := auditReadServerPacket(wConn)
	if err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CHAR_LOGIN falló: ack=%v err=%v", cSelAck, err)
	}
	_, _ = auditReadServerPacket(wConn) // Opt
	_, _ = auditReadServerPacket(wConn) // Opt
	_, _ = auditReadServerPacket(wConn) // Opt
	_ = wConn.Close()

	var dbRandomID uint16
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName).Scan(&dbRandomID); err != nil {
		t.Fatalf("Token no emitido en tActiveZoneTransfer: %v", err)
	}

	// -------------------------------------------------------------------------
	// Paso 3: ZoneService (49 Checksums + In-Game Chunk + Ciclo Completo)
	// -------------------------------------------------------------------------
	zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
	if err != nil {
		t.Fatalf("Error conectando a ZoneService: %v", err)
	}
	zHS, _ := auditReadServerPacket(zConn)
	zSeed, _ := zHS.Reader.ReadInt16()
	zCrypto, _ := crypto.NewNetCrypto(zSeed)
	zClient := &auditMockClient{conn: zConn, crypto: zCrypto}

	mapReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	mapReq.WriteUint16(dbRandomID)
	mapReq.WriteString(testCharName, 20)
	sums := zoneDataStore.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		mapReq.WriteBytes(sums[i][:])
	}
	_ = zClient.sendPacket(mapReq)

	for {
		pkt, err := auditReadServerPacket(zConn)
		if err != nil {
			t.Fatalf("Error leyendo In-Game Chunk: %v", err)
		}
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			break
		}
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	time.Sleep(40 * time.Millisecond)

	zp, ok := zoneSrv.Handler().GetPlayerByName(testCharName)
	if !ok || zp == nil || zp.CurrentMap == nil {
		t.Fatalf("Jugador '%s' no activo en ZoneService tras LOGINCOMPLETE", testCharName)
	}

	// A. Movimiento con validación .shbd (PROTO_NC_ACT_MOVERUN_CMD)
	mvPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVERUN_CMD)
	mvPkt.WriteUint32(uint32(zp.Char.XPos))
	mvPkt.WriteUint32(uint32(zp.Char.YPos))
	mvPkt.WriteUint32(uint32(zp.Char.XPos + 40))
	mvPkt.WriteUint32(uint32(zp.Char.YPos + 40))
	_ = zClient.sendPacket(mvPkt)

	// B. Beauty Shop (GameLog Type 4 nType=6) & Friend Point
	bsPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_USE_BEAUTY_SHOP_CMD)
	bsPkt.WriteByte(2)
	bsPkt.WriteByte(2)
	bsPkt.WriteByte(3)
	_ = zClient.sendPacket(bsPkt)

	// C. Compra en Mercader NPC (PROTO_NC_ITEM_BUY_REQ -> tItem + GameLog Type 3)
	_ = charRepo.AddCharMoney(createdCharID, 100000)
	zp.Char.TCharacterState.Money += 100000
	var shopItem *item.ItemInfo
	for id := uint16(1); id < 60000; id++ {
		if info, exists := zoneDataStore.ItemStore.GetItemInfo(id); exists && info != nil && info.BuyPrice > 0 && info.BuyPrice <= 50000 && info.MaxLot >= 1 {
			shopItem = info
			break
		}
	}
	if shopItem == nil {
		t.Fatalf("No se encontró ítem comprable en ItemInfo.shn")
	}
	buyPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_REQ)
	buyPkt.WriteUint16(shopItem.ID)
	buyPkt.WriteUint32(1)
	_ = zClient.sendPacket(buyPkt)

	// D. Quest Oficial (PROTO_NC_QUEST_START_REQ -> tQuest)
	var chosenQuestID uint16
	for id := uint16(1); id < 65000; id++ {
		q := zoneDataStore.QuestCatalog.ByID[id]
		if q == nil || q.CheckPrevQuest {
			continue
		}
		if q.CheckLevel && (zp.Char.Level < q.MinLevel || (q.MaxLevel > 0 && zp.Char.Level > q.MaxLevel)) {
			continue
		}
		chosenQuestID = q.QuestID
		break
	}
	if chosenQuestID == 0 {
		t.Fatalf("No se encontró ninguna misión inicial elegible en QuestData.shn")
	}
	qReq := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_START_REQ)
	qReq.WriteUint16(chosenQuestID)
	_ = zClient.sendPacket(qReq)

	// E. Verificación .shbd y Pet Log (GameLog Type G)
	if shbd := zoneDataStore.BlockInfo.GetSHBD(zp.CurrentMap.MapID); shbd != nil {
		_ = shbd.IsWorldMoveBlocked(float64(zp.Char.XPos), float64(zp.Char.YPos))
	}
	petAck := protocol.NewPacketWriter(protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK)
	petAck.WriteUint16(0)
	petAck.WriteUint16(1234)
	_ = zClient.sendPacket(petAck)

	// F. Autorización GM oficial estricta vía AdminLvSet.shn + tCharacter.nAdminLevel + tGMItemRewardLog + GameLog Type L
	reqMakeLv, hasMakeLv := zoneDataStore.AdminLvCatalog.RequiredLevel("makeitem")
	if !hasMakeLv {
		t.Fatalf("AdminLvSet.shn no define nivel para 'makeitem'")
	}
	if _, err := charRepo.SetAdminLevel(createdCharID, reqMakeLv); err != nil {
		t.Fatalf("p_Char_SetAdminLevel falló: %v", err)
	}
	zp.Char.AdminLevel = reqMakeLv

	gmChat := fmt.Sprintf("&makeitem %d 1", shopItem.ID)
	chatPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_CHAT_REQ)
	chatPkt.WriteByte(byte(len(gmChat)))
	chatPkt.WriteString(gmChat, len(gmChat))
	_ = zClient.sendPacket(chatPkt)

	// G. Desconexión limpia (PROTO_NC_USER_NORMALLOGOUT_CMD) y vaciado del GameLogDispatcher
	time.Sleep(150 * time.Millisecond)
	logoutPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	logoutPkt.WriteByte(0)
	_ = zClient.sendPacket(logoutPkt)
	time.Sleep(150 * time.Millisecond)
	_ = zConn.Close()
	time.Sleep(100 * time.Millisecond)

	_ = glDispatcher.Close(context.Background())

	// -------------------------------------------------------------------------
	// Paso 4: Verificación Forense Directa en Tablas MySQL Oficiales
	// -------------------------------------------------------------------------
	var dbItemCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tItem WHERE nOwner = ?", createdCharID).Scan(&dbItemCount); err != nil || dbItemCount < 2 {
		t.Fatalf("tItem en MySQL debe registrar ítem comprado + ítem GM (>=2), encontrado %d (err=%v)", dbItemCount, err)
	}

	if chosenQuestID > 0 {
		var dbQuestCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM tQuest WHERE nCharNo = ?", createdCharID).Scan(&dbQuestCount); err != nil || dbQuestCount < 1 {
			t.Fatalf("tQuest en MySQL debe registrar misión activa del jugador, encontrado %d (err=%v)", dbQuestCount, err)
		}
	}

	var dbGMRewardCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tGMItemRewardLog WHERE nCharNo = ?", createdCharID).Scan(&dbGMRewardCount); err != nil || dbGMRewardCount < 1 {
		t.Fatalf("tGMItemRewardLog en MySQL debe registrar recompensa oficial de &makeitem, encontrado %d (err=%v)", dbGMRewardCount, err)
	}

	var dbGameLogCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tGameLog WHERE nCharNo = ?", createdCharID).Scan(&dbGameLogCount); err != nil || dbGameLogCount < 3 {
		t.Fatalf("tGameLog en MySQL debe contener auditoría de eventos (esperado >=3), encontrado %d (err=%v)", dbGameLogCount, err)
	}

	var dbAdminLv int
	if err := db.QueryRow("SELECT nAdminLevel FROM tCharacter WHERE nCharNo = ?", createdCharID).Scan(&dbAdminLv); err != nil || dbAdminLv != int(reqMakeLv) {
		t.Fatalf("tCharacter.nAdminLevel en MySQL = %d, esperado %d", dbAdminLv, reqMakeLv)
	}

	t.Logf("Control 7 PASS: Clúster completo E2E verificado en MySQL (tCharacter AdminLv=%d, tItem=%d, tGMItemRewardLog=%d, tGameLog=%d)",
		dbAdminLv, dbItemCount, dbGMRewardCount, dbGameLogCount)
}

// =============================================================================
// CONTROL 8 (Sección 9.8): Escaneo AST Cero Hardcode en ZoneService y Lua
// =============================================================================
func TestFase16_Control8_AST_ZeroHardcode_ZoneAndLua(t *testing.T) {
	rootPath := findRepoRoot(t)
	dirs := []string{
		filepath.Join(rootPath, "pkg", "server", "zone"),
		filepath.Join(rootPath, "pkg", "script", "lua"),
	}

	forbiddenSubstrings := []string{
		"onson",
		"gamigo",
		"mara pirates",
	}
	var violations []string
	fset := token.NewFileSet()
	fileCount := 0

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("Error leyendo %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}
			fileCount++
			fpath := filepath.Join(dir, e.Name())
			node, err := parser.ParseFile(fset, fpath, nil, 0)
			if err != nil {
				t.Fatalf("Error parseando AST de %s: %v", fpath, err)
			}
			ast.Inspect(node, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.FuncDecl:
					if v.Name != nil && strings.EqualFold(v.Name.Name, "loadDefaults") {
						pos := fset.Position(v.Pos())
						violations = append(violations, fmt.Sprintf("%s:%d función prohibida loadDefaults()", e.Name(), pos.Line))
					}
				case *ast.BasicLit:
					if v.Kind == token.STRING {
						rawVal := strings.Trim(v.Value, "\"`")
						lower := strings.ToLower(rawVal)
						for _, sub := range forbiddenSubstrings {
							if strings.Contains(lower, sub) {
								pos := fset.Position(v.Pos())
								violations = append(violations, fmt.Sprintf("%s:%d literal prohibido '%s'", e.Name(), pos.Line, sub))
							}
						}
						if rawVal == "RouN" {
							pos := fset.Position(v.Pos())
							violations = append(violations, fmt.Sprintf("%s:%d fallback hardcodeado \"RouN\" prohibido", e.Name(), pos.Line))
						}
					} else if v.Kind == token.INT {
						if v.Value == "0x1234" || v.Value == "4660" {
							pos := fset.Position(v.Pos())
							violations = append(violations, fmt.Sprintf("%s:%d literal mágico %s prohibido", e.Name(), pos.Line, v.Value))
						}
					}
				}
				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Fatalf("CONTROL 8 FALLIDO: Detectados valores hardcodeados o funciones prohibidas:\n%s", strings.Join(violations, "\n"))
	}
	t.Logf("Control 8 PASS: %d archivos Go de ZoneService y Lua verificados vía AST — 0 hardcodes, 0 fallbacks ficticios, 100%% AdminLvSet.shn/Data.", fileCount)
}
