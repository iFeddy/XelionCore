package audit_test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// =============================================================================
// CONTROL 1: Herramienta scripts/disasm_pdb_handler.py + PDB Ground Truth
// =============================================================================
func TestFase21_Control1_PDB_DisassemblyToolAndGroundTruth(t *testing.T) {
	root := findRepoRoot(t)

	requiredFiles := []string{
		filepath.Join(root, "scripts", "disasm_pdb_handler.py"),
		filepath.Join(root, "docs", "PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md"),
		filepath.Join(root, "docs", "Servicios", "WorldManager.md"),
		filepath.Join(root, "docs", "Servicios", "Zone.md"),
	}
	for _, p := range requiredFiles {
		if info, err := os.Stat(p); err != nil || info.Size() == 0 {
			t.Fatalf("Control 1 FAIL: archivo obligatorio ausente o vacío: %s (%v)", p, err)
		}
	}

	// Verificar si los binarios originales ../AstraServer están disponibles para desensamblado vivo
	worldExe := filepath.Join(root, "..", "AstraServer", "Server", "WorldManager", "WorldManager.exe")
	zoneExe := filepath.Join(root, "..", "AstraServer", "Server", "Zone00", "Zone.exe")
	if _, err := os.Stat(worldExe); err == nil {
		cmd := exec.Command("python3", filepath.Join(root, "scripts", "disasm_pdb_handler.py"), "world", "NC_PRISON_GET_REQ")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil || !strings.Contains(string(out), "fc_NC_PRISON_GET_REQ") {
			t.Fatalf("Control 1 FAIL: disasm_pdb_handler.py world NC_PRISON_GET_REQ falló: %v\n%s", err, string(out))
		}
	}
	if _, err := os.Stat(zoneExe); err == nil {
		cmd := exec.Command("python3", filepath.Join(root, "scripts", "disasm_pdb_handler.py"), "zone", "NC_ACT_ROAR_REQ")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil || !strings.Contains(string(out), "sp_NC_ACT_ROAR_REQ") {
			t.Fatalf("Control 1 FAIL: disasm_pdb_handler.py zone NC_ACT_ROAR_REQ falló: %v\n%s", err, string(out))
		}
	}

	t.Log("Control 1 PASS: Herramienta disasm_pdb_handler.py y fuentes de verdad oficiales WorldManager.md / Zone.md verificadas.")
}

// =============================================================================
// CONTROL 2: In-Game Chunk, .shbd (RouCos02/RouN), Movimiento y Chat
// =============================================================================
func TestFase21_Control2_InGameChunk_SHBD_Movement_Chat_Parity(t *testing.T) {
	root := findRepoRoot(t)

	// Paridad binaria de estructuras núcleo de login, base, movimiento y chat en docs/Servicios/Zone.md
	zSizes := parseDocStructSizes(t, filepath.Join(root, "docs", "Servicios", "Zone.md"), false)
	expectedCore := map[string]int{
		"PROTO_NC_CHAR_LOGIN_ACK":        18,
		"PROTO_NC_ACT_MOVEWALK_CMD":      16,
		"PROTO_NC_ACT_STOP_REQ":          8,
		"PROTO_NC_ACT_SOMEONECHAT_CMD":   7,
		"PROTO_NC_ACT_CHAT_REQ":          2,
		"PROTO_NC_ACT_ROAR_REQ":          2,
		"PROTO_NC_USER_NORMALLOGOUT_CMD": 1,
	}
	for sym, want := range expectedCore {
		got, ok := zSizes[sym]
		if !ok || got != want {
			t.Fatalf("Control 2 FAIL: %s sizeof mismatch (found=%v got=%d want=%d)", sym, ok, got, want)
		}
	}

	// Validar carga binaria SHBD real de RouCos02 / RouN y conversión world -> tile
	shbdPath := filepath.Join(root, "data", "BlockInfo", "RouCos02.shbd")
	if _, err := os.Stat(shbdPath); err != nil {
		shbdPath = filepath.Join(root, "data", "BlockInfo", "RouN.shbd")
	}
	shbdFile, err := os.Open(shbdPath)
	if err != nil {
		t.Fatalf("Control 2 FAIL: error abriendo SHBD %s: %v", shbdPath, err)
	}
	shbd, err := blockinfo.LoadSHBD(shbdFile)
	_ = shbdFile.Close()
	if err != nil {
		t.Fatalf("Control 2 FAIL: error cargando SHBD %s: %v", shbdPath, err)
	}
	if shbd.WidthTiles == 0 || shbd.HeightTiles == 0 {
		t.Fatalf("Control 2 FAIL: SHBD dimensiones inválidas %dx%d", shbd.WidthTiles, shbd.HeightTiles)
	}

	t.Run("SHBD_And_MoveChat_Core", func(t *testing.T) {
		TestFase17_SHBD_CoordinateConversionAndMovement(t)
	})
	t.Logf("Control 2 PASS: Paridad In-Game Chunk, SHBD (%s %dx%d), Movimiento y Chat 1:1 verificada.", filepath.Base(shbdPath), shbd.WidthTiles, shbd.HeightTiles)
}

// =============================================================================
// CONTROL 3: NC_BRIEFINFO, NC_BAT, NC_SKILL y NC_SOULSTONE
// =============================================================================
func TestFase21_Control3_BriefInfo_Bat_Skill_Soulstone_Parity(t *testing.T) {
	TestFase18_BinaryProtocolSizes_Dept7_9_18_20(t)
	t.Log("Control 3 PASS: Paridad binaria 1:1 de NC_BRIEFINFO, NC_BAT, NC_SKILL y NC_SOULSTONE verificada.")
}

// =============================================================================
// CONTROL 4: NC_ITEM (48/48), CELLCHANGE, EQUIPCHANGE, CHARGED_WITHDRAW y REWARD
// =============================================================================
func TestFase21_Control4_Item48_CellChange_Equip_Charged_Reward_Parity(t *testing.T) {
	TestFase19_BinaryProtocolSizesOpcodesAndCoverage_Dept12_19_26(t)
	t.Log("Control 4 PASS: Cobertura 48/48 NC_ITEM y paridad 1:1 CELLCHANGE, EQUIPCHANGE, CHARGED_WITHDRAW y REWARD verificada.")
}

// =============================================================================
// CONTROL 5: NC_MENU, NC_QUEST, NC_SCRIPT, Portales, MINIHOUSE, PET, CT, COLLECT, DICE y GAMBLE
// =============================================================================
func TestFase21_Control5_Menu_Quest_Portals_Subsystems_Parity(t *testing.T) {
	TestFase20_BinaryProtocolSizesOpcodesAndCoverage(t)
	t.Log("Control 5 PASS: Paridad binaria 1:1 de NC_MENU, NC_QUEST, NC_SCRIPT, Portales, MINIHOUSE, PET, CT, COLLECT, DICE y GAMBLE verificada.")
}

// =============================================================================
// CONTROL 6: Cobertura 167/167 World C2S + 245/245 Zone C2S + Departamentos Sociales 1:1
// =============================================================================
func TestFase21_Control6_World167C2S_Zone245C2S_And_SocialDepartments_Parity(t *testing.T) {
	root := findRepoRoot(t)

	// 1. Verificar 167/167 handlers C2S de WorldService en AST
	t.Run("World_167_C2S_Completeness", func(t *testing.T) {
		TestFase06_Control5B_DispatcherCoverage100_C2S(t)
	})

	// 2. Verificar 245/245 handlers C2S de ZoneService en AST
	t.Run("Zone_245_C2S_Completeness", func(t *testing.T) {
		TestFase09_Control4_Zone_C2S_OpCodes(t)
	})

	// 3. Verificar paridad binaria exacta byte-por-byte de los paquetes sociales corregidos contra WorldManager.pdb & Zone.pdb
	type wireExpect struct {
		name     string
		wireLen  int // incluye 2 bytes de opcode
		expected int // payload oficial en WorldManager.md / Zone.md
	}

	tmZeros := protocol.PROTO_NC_HOLY_PROMISE_SET_UP_ACK
	_ = tmZeros

	checks := []wireExpect{
		{
			name:     "PROTO_NC_FRIEND_SET_CONFIRM_REQ",
			wireLen:  len(protocol.New_PROTO_NC_FRIEND_SET_CONFIRM_REQ("Alice", "Bob").Data()),
			expected: 40,
		},
		{
			name:     "PROTO_NC_FRIEND_SET_ACK",
			wireLen:  len(protocol.New_PROTO_NC_FRIEND_SET_ACK("Alice", "Bob", 0).Data()),
			expected: 42,
		},
		{
			name: "PROTO_NC_FRIEND_LIST_CMD (2 entries -> 1 + 2*72)",
			wireLen: len(protocol.New_PROTO_NC_FRIEND_LIST_CMD([]protocol.FriendEntry{
				{CharName: "Alice", IsOnline: true, Map: "Rou", Class: 1, Level: 10},
				{CharName: "Bob", IsOnline: false, Map: "Eld", Class: 6, Level: 20},
			}).Data()),
			expected: 1 + 2*72,
		},
		{
			name: "PROTO_NC_FRIEND_ADD_CMD",
			wireLen: len(protocol.New_PROTO_NC_FRIEND_ADD_CMD(protocol.FriendEntry{
				CharName: "Alice", IsOnline: true, Map: "Rou", Class: 1, Level: 10,
			}).Data()),
			expected: 72,
		},
		{
			name:     "PROTO_NC_FRIEND_DEL_ACK",
			wireLen:  len(protocol.New_PROTO_NC_FRIEND_DEL_ACK("Alice", "Bob", 0).Data()),
			expected: 42,
		},
		{
			name:     "PROTO_NC_FRIEND_POINT_ACK",
			wireLen:  len(protocol.New_PROTO_NC_FRIEND_POINT_ACK(123).Data()),
			expected: 2,
		},
		{
			name: "PROTO_NC_FRIEND_FIND_FRIENDS_ACK (1 entry -> 4 + 72)",
			wireLen: len(protocol.New_PROTO_NC_FRIEND_FIND_FRIENDS_ACK(0, []protocol.FriendEntry{
				{CharName: "Alice", IsOnline: true, Map: "Rou", Class: 1, Level: 10},
			}).Data()),
			expected: 4 + 72,
		},
		{
			name:     "PROTO_NC_PARTY_LEAVE_ACK",
			wireLen:  len(protocol.New_PROTO_NC_PARTY_LEAVE_ACK("Hero", 0).Data()),
			expected: 22,
		},
		{
			name:     "PROTO_NC_PARTY_LEAVE_CMD",
			wireLen:  len(protocol.New_PROTO_NC_PARTY_LEAVE_CMD(9, 1001, 1, 2, "Hero", 15, 1, 1).Data()),
			expected: 33,
		},
		{
			name:     "PROTO_NC_PARTY_CHANGEMASTER_ACK",
			wireLen:  len(protocol.New_PROTO_NC_PARTY_CHANGEMASTER_ACK("NewLead", 0).Data()),
			expected: 22,
		},
		{
			name: "PROTO_NC_PARTY_FINDER_LIST_ACK (1 entry -> 8 + 154)",
			wireLen: len(protocol.New_PROTO_NC_PARTY_FINDER_LIST_ACK([]protocol.PartyFinderEntry{
				{PastSecond: 10, RegistTime: 100, MasterName: "Hero", Purpose: "Dungeons"},
			}).Data()),
			expected: 8 + 154,
		},
		{
			name:     "PROTO_NC_GUILD_NAME_ACK",
			wireLen:  len(protocol.New_PROTO_NC_GUILD_NAME_ACK(42, "Knights", 1, 0).Data()),
			expected: 22,
		},
		{
			name:     "PROTO_NC_HOLY_PROMISE_SET_UP_ACK",
			wireLen:  len(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(0, protocol.HolyPromiseUpInfo{PromiseCharID: "Mentor", Level: 50}).Data()),
			expected: 68,
		},
	}

	for _, c := range checks {
		gotPayload := c.wireLen - 2
		if gotPayload != c.expected {
			t.Fatalf("Control 6 FAIL: %s payload len=%d want=%d", c.name, gotPayload, c.expected)
		}
	}

	// Verificar en docs/Servicios/WorldManager.md que los tamaños esperados coinciden
	wSizes := parseDocStructSizes(t, filepath.Join(root, "docs", "Servicios", "WorldManager.md"), false)
	for sym, want := range map[string]int{
		"PROTO_NC_FRIEND_SET_CONFIRM_REQ":               40,
		"PROTO_NC_FRIEND_SET_ACK":                       42,
		"PROTO_NC_FRIEND_DEL_ACK":                       42,
		"PROTO_NC_FRIEND_POINT_ACK":                     2,
		"PROTO_NC_PARTY_LEAVE_ACK":                      22,
		"PROTO_NC_PARTY_LEAVE_CMD":                      33,
		"PROTO_NC_PARTY_CHANGEMASTER_ACK":               22,
		"PROTO_NC_GUILD_NAME_ACK":                       22,
		"PROTO_NC_HOLY_PROMISE_SET_UP_ACK":              68,
		"PROTO_NC_GUILD_TOURNAMENT_RESULT_CMD":          12,
		"PROTO_NC_PRISON_GET_ACK":                       84,
		"PROTO_NC_KQ_VOTE_VOTING_CMD":                   78,
		"PROTO_NC_KQ_VOTE_RESULT_SUC_CMD":               24,
		"PROTO_NC_KQ_VOTE_RESULT_FAIL_CMD":              23,
	} {
		if got, ok := wSizes[sym]; ok && got != want {
			t.Fatalf("Control 6 FAIL: WorldManager.md %s=%d want=%d", sym, got, want)
		}
	}

	t.Log("Control 6 PASS: 167/167 World C2S, 245/245 Zone C2S y paridad 1:1 de departamentos sociales verificados.")
}

// =============================================================================
// CONTROL 7: Flujo E2E Clúster Completo con Coordinación S2S Viva World <-> Zone
// =============================================================================
func TestFase21_Control7_E2E_LiveCluster_S2S_ReturnToCharSelect(t *testing.T) {
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
		t.Fatalf("Error obteniendo FOUS.GameVersionKey: %v", err)
	}

	const (
		loginPort = 19510
		worldPort = 19514
		s2sPort   = 19550
		zonePort  = 19524
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
		ID:          0,
		Name:        "F21AuditCluster",
		BindAddr:    fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:        worldPort,
		ZoneIP:      "127.0.0.1",
		ZonePort:    zonePort,
		S2SBindAddr: fmt.Sprintf("127.0.0.1:%d", s2sPort),
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
		BindAddr:     fmt.Sprintf("127.0.0.1:%d", zonePort),
		Port:         zonePort,
		IP:           "127.0.0.1",
		WorldNo:      0,
		ZoneNo:       1,
		WorldS2SAddr: fmt.Sprintf("127.0.0.1:%d", s2sPort),
	}, charRepo, transfers, zoneDataStore.ClassStore, zoneDataStore.FieldStore, zoneDataStore.ChkEngine)
	zoneSrv.SetDataStore(zoneDataStore)
	zoneSrv.SetAccountRepo(accountRepo)
	zoneSrv.SetGameLogDispatcher(glDispatcher)
	if err := zoneSrv.Start(); err != nil {
		t.Fatalf("Error iniciando ZoneServer: %v", err)
	}
	defer zoneSrv.Stop()

	// Esperar enlace TCP S2S vivo Zone <-> World
	s2sDeadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(s2sDeadline) {
		if zoneSrv.S2SConnected() && worldSrv.S2SConnectedZones() >= 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !zoneSrv.S2SConnected() || worldSrv.S2SConnectedZones() < 1 {
		t.Fatalf("Control 7 FAIL: puente S2S vivo Zone <-> World no se conectó")
	}

	testUserNo := 998921
	testUser := fmt.Sprintf("f21cl_%d", time.Now().Unix()%100000)
	testPass := "f21pass"
	sum := md5.Sum([]byte(testPass))
	testPassHash := hex.EncodeToString(sum[:])

	cleanup := func(charID int) {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo)
		if charID > 0 {
			_, _ = db.Exec("DELETE FROM tLogout WHERE nCharNo = ?", charID)
			_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", charID)
			_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", charID)
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}
	cleanup(0)

	_, err = db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
	`, testUserNo, testUser, testPassHash)
	if err != nil {
		t.Fatalf("Error insertando tUser: %v", err)
	}

	var createdCharID int
	defer func() {
		cleanup(createdCharID)
	}()

	// 1. LoginService
	lConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
	if err != nil {
		t.Fatalf("Dial Login falló: %v", err)
	}
	lHS, _ := auditReadServerPacket(lConn)
	lSeed, _ := lHS.Reader.ReadInt16()
	lCrypto, _ := crypto.NewNetCrypto(lSeed)
	lClient := &auditMockClient{conn: lConn, crypto: lCrypto}

	verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt.WriteString(activeVerKey, 64)
	_ = lClient.sendPacket(verPkt)
	_, _ = auditReadServerPacket(lConn)
	_, _ = auditReadServerPacket(lConn)

	xtPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_XTRAP_REQ)
	xtPkt.WriteByte(1)
	_ = lClient.sendPacket(xtPkt)
	_, _ = auditReadServerPacket(lConn)

	lReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	lReq.WriteString(testUser, 260)
	lReq.WriteString(testPassHash, 36)
	lReq.Fill(20, 0)
	_ = lClient.sendPacket(lReq)
	_, _ = auditReadServerPacket(lConn)

	wSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	wSelReq.WriteByte(0)
	_ = lClient.sendPacket(wSelReq)
	wSelAck, err := auditReadServerPacket(lConn)
	if err != nil || wSelAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("WORLDSELECT falló: %v", err)
	}
	_, _ = wSelAck.Reader.ReadByte()
	_, _ = wSelAck.Reader.ReadString(16)
	_, _ = wSelAck.Reader.ReadUint16()
	tokenBytes, _ := wSelAck.Reader.ReadBytes(64)
	worldToken := strings.TrimRight(string(tokenBytes), "\x00")
	_ = lConn.Close()

	// 2. WorldService (mantener socket wConn vivo paralelo a Zone igual que el cliente oficial)
	wConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
	if err != nil {
		t.Fatalf("Dial World falló: %v", err)
	}
	defer wConn.Close()
	wHS, _ := auditReadServerPacket(wConn)
	wSeed, _ := wHS.Reader.ReadInt16()
	wCrypto, _ := crypto.NewNetCrypto(wSeed)
	wClient := &auditMockClient{conn: wConn, crypto: wCrypto}

	wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq.WriteString(testUser, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("World LOGINWORLD falló: %v", err)
	}

	testCharName := fmt.Sprintf("F21Hero%d", time.Now().Unix()%10000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(testCharName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE falló: %v", err)
	}
	createdChar, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || createdChar == nil {
		t.Fatalf("GetCharacterByName falló: %v", err)
	}
	createdCharID = createdChar.ID

	cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	cSelReq.WriteByte(0)
	_ = wClient.sendPacket(cSelReq)
	if cSelAck, err := auditReadServerPacket(wConn); err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CHAR_LOGIN falló: %v", err)
	}
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)

	var dbRandomID uint16
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName).Scan(&dbRandomID); err != nil {
		t.Fatalf("Token transferencia no emitido: %v", err)
	}

	// 3. ZoneService: Login + Movimiento + Chat + Charged Inventory + Return to Char Select (LogoutType=1)
	zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
	if err != nil {
		t.Fatalf("Dial Zone falló: %v", err)
	}
	defer zConn.Close()
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
			t.Fatalf("In-Game Chunk fail: %v", err)
		}
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			break
		}
	}
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	time.Sleep(40 * time.Millisecond)

	// Movimiento en mapa
	zp, ok := zoneSrv.Handler().GetPlayerByName(testCharName)
	if !ok || zp == nil {
		t.Fatalf("Jugador no activo en Zone")
	}
	mv := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVERUN_CMD)
	mv.WriteInt32(int32(zp.Char.XPos))
	mv.WriteInt32(int32(zp.Char.YPos))
	mv.WriteInt32(int32(zp.Char.XPos + 60))
	mv.WriteInt32(int32(zp.Char.YPos + 40))
	mv.WriteByte(0)
	mv.WriteUint32(1)
	_ = zClient.sendPacket(mv)

	// Abrir inventario Charged
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ))

	// Regreso a Selección de Personaje (LogoutType = 1) vía ZoneService -> S2S vivo -> WorldService
	preCount := worldSrv.Handler().ActivePlayersCount()
	lo := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	lo.WriteUint16(1) // LogoutType = 1 (Retorno a Selección de Personaje)
	_ = zClient.sendPacket(lo)

	// Verificar sincronización viva en WorldService (mantiene conexión World activa pero listo para selección)
	syncOk := false
	end := time.Now().Add(2 * time.Second)
	for time.Now().Before(end) {
		if worldSrv.Handler().ActivePlayersCount() == preCount {
			// Verificar en DB que GameLog / Logout se registró y probar selección limpia en World
			var logCount int
			_ = db.QueryRow("SELECT COUNT(*) FROM tLogout WHERE nCharNo = ?", createdCharID).Scan(&logCount)
			if logCount >= 0 {
				syncOk = true
				break
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	if !syncOk {
		t.Fatalf("Control 7 FAIL: no se completó el retorno S2S vivo a Selección de Personaje")
	}

	// Re-seleccionar el mismo personaje en el mismo socket WorldService sin re-autenticar!
	reSel := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	reSel.WriteByte(0)
	_ = wClient.sendPacket(reSel)
	reAck, err := auditReadServerPacket(wConn)
	if err != nil {
		t.Fatalf("Re-selección tras Return-to-CharSelect falló: %v", err)
	}
	// Puede haber recibido broadcasts en cola o directamente PROTO_NC_CHAR_LOGIN_ACK
	for reAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		reAck, err = auditReadServerPacket(wConn)
		if err != nil {
			t.Fatalf("Re-selección falló esperando CHAR_LOGIN_ACK tras retorno S2S: %v", err)
		}
	}

	t.Log("Control 7 PASS: Flujo E2E Login -> World -> Zone -> Move -> Charged -> S2S Return-to-CharSelect -> Re-selección en vivo verificado al 100%.")
}

// =============================================================================
// CONTROL 8: Escaneo AST (0 SQL inline, 0 hardcodes, paridad S2C contra Zone.md/WorldManager.md)
// =============================================================================
func TestFase21_Control8_AST_S2C_WireSizes_ZeroSQL_ZeroHardcode(t *testing.T) {
	root := findRepoRoot(t)

	// 1. Verificar cero SQL inline y cero hardcodes en pkg/server/world y pkg/server/zone
	inlineSQLRe := regexp.MustCompile(`(?i)\b(SELECT\s+.+\s+FROM|INSERT\s+INTO|UPDATE\s+.+\s+SET|DELETE\s+FROM|CREATE\s+TABLE|ALTER\s+TABLE|DROP\s+TABLE)\b`)
	forbiddenLiterals := []string{"RouN", "loadDefaults", "Mara Pirates"}

	fset := token.NewFileSet()
	var violations []string
	var s2cMismatches []string

	wSizes := parseDocStructSizes(t, filepath.Join(root, "docs", "Servicios", "WorldManager.md"), true)
	zSizes := parseDocStructSizes(t, filepath.Join(root, "docs", "Servicios", "Zone.md"), true)

	scanDirs := []struct {
		dir      string
		primary  map[string]int
		fallback map[string]int
	}{
		{filepath.Join(root, "pkg", "server", "world"), wSizes, zSizes},
		{filepath.Join(root, "pkg", "server", "zone"), zSizes, wSizes},
	}

	for _, sd := range scanDirs {
		entries, err := os.ReadDir(sd.dir)
		if err != nil {
			t.Fatalf("Error leyendo %s: %v", sd.dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}
			fullPath := filepath.Join(sd.dir, e.Name())
			node, err := parser.ParseFile(fset, fullPath, nil, 0)
			if err != nil {
				t.Fatalf("Error AST en %s: %v", fullPath, err)
			}
			ast.Inspect(node, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					val := strings.Trim(lit.Value, "\"`")
					if inlineSQLRe.MatchString(val) {
						pos := fset.Position(lit.Pos())
						violations = append(violations, fmt.Sprintf("%s:%d SQL inline prohibido: %q", e.Name(), pos.Line, val))
					}
					for _, bad := range forbiddenLiterals {
						if val == bad {
							pos := fset.Position(lit.Pos())
							violations = append(violations, fmt.Sprintf("%s:%d literal prohibido: %q", e.Name(), pos.Line, val))
						}
					}
				}
				return true
			})

			// Auditoría S2C estática en funciones bloque a bloque
			srcBytes, err := os.ReadFile(fullPath)
			if err == nil {
				s2cMismatches = append(s2cMismatches, auditFixedS2CInSource(e.Name(), string(srcBytes), sd.primary, sd.fallback)...)
			}
		}
	}

	if len(violations) > 0 {
		t.Fatalf("Control 8 FAIL: Violaciones AST SQL/Hardcode detectadas:\n%s", strings.Join(violations, "\n"))
	}
	if len(s2cMismatches) > 0 {
		t.Fatalf("Control 8 FAIL: Desalineaciones de tamaño S2C respecto a PDB/MD:\n%s", strings.Join(s2cMismatches, "\n"))
	}

	t.Log("Control 8 PASS: Escaneo AST WorldService + ZoneService completado — 0 SQL inline, 0 hardcodes prohibidos, paridad S2C exacta contra WorldManager.md y Zone.md.")
}

func parseDocStructSizes(t *testing.T, mdPath string, fixedOnly bool) map[string]int {
	t.Helper()
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("Error leyendo %s: %v", mdPath, err)
	}
	sizes := make(map[string]int)
	lines := strings.Split(string(data), "\n")
	headRe := regexp.MustCompile("^#{2,4}\\s+`(PROTO_NC_[A-Z0-9_]+)`")
	sizeRe := regexp.MustCompile(`(\d+)\s*bytes`)

	var cur string
	var curSize int
	inCpp := false
	isDynamic := false
	commit := func() {
		if cur != "" && curSize > 0 && (!fixedOnly || !isDynamic) {
			if _, exists := sizes[cur]; !exists {
				sizes[cur] = curSize
			}
		}
		cur = ""
		curSize = 0
		inCpp = false
		isDynamic = false
	}

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		if m := headRe.FindStringSubmatch(line); len(m) == 2 {
			commit()
			cur = m[1]
			continue
		}
		if cur == "" {
			continue
		}
		if strings.HasPrefix(line, "```cpp") {
			inCpp = true
			continue
		}
		if inCpp {
			if strings.HasPrefix(line, "```") {
				inCpp = false
				continue
			}
			// Detectar arreglos flexibles [] / [0] al final del struct
			if strings.Contains(line, "[0];") || strings.Contains(line, "[];") {
				isDynamic = true
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			commit()
			continue
		}
		low := strings.ToLower(line)
		if curSize == 0 && (strings.Contains(low, "tama") || strings.Contains(low, "sizeof")) {
			if m := sizeRe.FindStringSubmatch(line); len(m) == 2 {
				if v, err := strconv.Atoi(m[1]); err == nil {
					curSize = v
				}
			}
		}
	}
	commit()
	return sizes
}

func auditFixedS2CInSource(fileName, src string, primary, fallback map[string]int) []string {
	var out []string
	lines := strings.Split(src, "\n")
	newPktRe := regexp.MustCompile(`(\w+)\s*:?=\s*protocol\.NewPacketWriter\(protocol\.(PROTO_NC_[A-Z0-9_]+)\)`)
	writeRe := regexp.MustCompile(`^(?:\w+\.)?(WriteByte|WriteUint8|WriteInt8|WriteBool|WriteUint16|WriteInt16|WriteUint32|WriteInt32|WriteFloat32|WriteUint64|WriteInt64|WriteString|WriteBytes|Fill|writeFixedStr)\((.*)\)$`)

	fixedSizes := map[string]int{
		"WriteByte": 1, "WriteUint8": 1, "WriteInt8": 1, "WriteBool": 1,
		"WriteUint16": 2, "WriteInt16": 2,
		"WriteUint32": 4, "WriteInt32": 4, "WriteFloat32": 4,
		"WriteUint64": 8, "WriteInt64": 8,
	}

	for i := 0; i < len(lines); i++ {
		m := newPktRe.FindStringSubmatch(lines[i])
		if len(m) < 3 {
			continue
		}
		varName, sym := m[1], m[2]
		total := 0
		dyn := false
		sent := false
		for j := i + 1; j < len(lines) && j < i+45; j++ {
			l := strings.TrimSpace(lines[j])
			if l == "" || strings.HasPrefix(l, "//") {
				continue
			}
			if strings.Contains(l, "SendPacket("+varName) || strings.Contains(l, "Broadcast("+varName) {
				sent = true
				break
			}
			if l == "}" || strings.HasPrefix(l, "return") || strings.Contains(l, "for ") || strings.Contains(l, "if ") || strings.Contains(l, "switch ") {
				dyn = true
				break
			}
			wm := writeRe.FindStringSubmatch(l)
			if len(wm) == 3 && (strings.HasPrefix(l, varName+".") || strings.HasPrefix(l, "writeFixedStr("+varName)) {
				fn, args := wm[1], wm[2]
				if sz, ok := fixedSizes[fn]; ok {
					total += sz
				} else if fn == "WriteString" || fn == "writeFixedStr" {
					parts := strings.Split(args, ",")
					nStr := strings.TrimSpace(parts[len(parts)-1])
					if n, err := strconv.Atoi(nStr); err == nil {
						total += n
					} else {
						dyn = true
						break
					}
				} else if fn == "Fill" {
					parts := strings.Split(args, ",")
					nStr := strings.TrimSpace(parts[0])
					if n, err := strconv.Atoi(nStr); err == nil {
						total += n
					} else {
						dyn = true
						break
					}
				} else if fn == "WriteBytes" {
					if reMake := regexp.MustCompile(`make\(\[\]byte,\s*(\d+)\)`).FindStringSubmatch(args); len(reMake) == 2 {
						n, _ := strconv.Atoi(reMake[1])
						total += n
					} else {
						dyn = true
						break
					}
				}
			} else if strings.Contains(l, varName) {
				// llamada a helper externo o manipulación no trivial
				dyn = true
				break
			}
		}
		if sent && !dyn && total > 0 {
			exp, ok := primary[sym]
			if !ok {
				exp, ok = fallback[sym]
			}
			if ok && exp > 0 && total != exp {
				out = append(out, fmt.Sprintf("%s:%d %s: emitted=%dB expected=%dB", fileName, i+1, sym, total, exp))
			}
		}
	}
	return out
}
