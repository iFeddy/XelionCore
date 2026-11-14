package audit_test

import (
	"crypto/md5"
	"encoding/binary"
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

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// TestFase19_BinaryProtocolSizesOpcodesAndCoverage_Dept12_19_26 verifica:
//  1. La fórmula exacta OpCode = (department << 10) | command para NC_ITEM (12),
//     NC_MENU (15), NC_CHARSAVE (16), NC_TRADE (19), NC_BOOTH (26), NC_GUILD (29) y NC_GUILD_ACADEMY (38).
//  2. El tamaño empaquetado binary.Size / FixedWireSize de structs NC_ITEM contra Zone.md y Zone.pdb.
//  3. Cobertura AST 100% en Zone Handler.HandlePacket de 48/48 C2S NC_ITEM,
//     10/10 NC_TRADE, 8/8 NC_BOOTH, 1/1 NC_CHARSAVE (0x4025) y 1/1 NC_ITEMDB (0x74A0).
func TestFase19_BinaryProtocolSizesOpcodesAndCoverage_Dept12_19_26(t *testing.T) {
	opcodeChecks := []struct {
		name string
		got  uint16
		dept uint16
		cmd  uint16
	}{
		{"PROTO_NC_CHAR_CLIENT_ITEM_CMD", protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD, 4, 71},
		{"PROTO_NC_ITEM_CELLCHANGE_CMD", protocol.PROTO_NC_ITEM_CELLCHANGE_CMD, 12, 1},
		{"PROTO_NC_ITEM_EQUIPCHANGE_CMD", protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD, 12, 2},
		{"PROTO_NC_ITEM_BUY_REQ", protocol.PROTO_NC_ITEM_BUY_REQ, 12, 3},
		{"PROTO_NC_ITEM_BUY_ACK", protocol.PROTO_NC_ITEM_BUY_ACK, 12, 4},
		{"PROTO_NC_ITEM_SELL_ACK", protocol.PROTO_NC_ITEM_SELL_ACK, 12, 5},
		{"PROTO_NC_ITEM_SELL_REQ", protocol.PROTO_NC_ITEM_SELL_REQ, 12, 6},
		{"PROTO_NC_ITEM_DROP_REQ", protocol.PROTO_NC_ITEM_DROP_REQ, 12, 7},
		{"PROTO_NC_ITEM_DROP_ACK", protocol.PROTO_NC_ITEM_DROP_ACK, 12, 8},
		{"PROTO_NC_ITEM_PICK_REQ", protocol.PROTO_NC_ITEM_PICK_REQ, 12, 9},
		{"PROTO_NC_ITEM_PICK_ACK", protocol.PROTO_NC_ITEM_PICK_ACK, 12, 10},
		{"PROTO_NC_ITEM_RELOC_REQ", protocol.PROTO_NC_ITEM_RELOC_REQ, 12, 11},
		{"PROTO_NC_ITEM_RELOC_ACK", protocol.PROTO_NC_ITEM_RELOC_ACK, 12, 12},
		{"PROTO_NC_ITEM_SPLIT_REQ", protocol.PROTO_NC_ITEM_SPLIT_REQ, 12, 13},
		{"PROTO_NC_ITEM_SPLIT_ACK", protocol.PROTO_NC_ITEM_SPLIT_ACK, 12, 14},
		{"PROTO_NC_ITEM_EQUIP_REQ", protocol.PROTO_NC_ITEM_EQUIP_REQ, 12, 15},
		{"PROTO_NC_ITEM_RINGEQUIP_REQ", protocol.PROTO_NC_ITEM_RINGEQUIP_REQ, 12, 16},
		{"PROTO_NC_ITEM_EQUIP_ACK", protocol.PROTO_NC_ITEM_EQUIP_ACK, 12, 17},
		{"PROTO_NC_ITEM_UNEQUIP_REQ", protocol.PROTO_NC_ITEM_UNEQUIP_REQ, 12, 18},
		{"PROTO_NC_ITEM_UNEQUIP_ACK", protocol.PROTO_NC_ITEM_UNEQUIP_ACK, 12, 19},
		{"PROTO_NC_ITEM_USE_REQ", protocol.PROTO_NC_ITEM_USE_REQ, 12, 21},
		{"PROTO_NC_ITEM_USE_ACK", protocol.PROTO_NC_ITEM_USE_ACK, 12, 22},
		{"PROTO_NC_ITEM_UPGRADE_REQ", protocol.PROTO_NC_ITEM_UPGRADE_REQ, 12, 23},
		{"PROTO_NC_ITEM_UPGRADE_ACK", protocol.PROTO_NC_ITEM_UPGRADE_ACK, 12, 24},
		{"PROTO_NC_ITEM_USECOMPLETE_CMD", protocol.PROTO_NC_ITEM_USECOMPLETE_CMD, 12, 26},
		{"PROTO_NC_ITEM_DEPOSIT_REQ", protocol.PROTO_NC_ITEM_DEPOSIT_REQ, 12, 28},
		{"PROTO_NC_ITEM_DEPOSIT_ACK", protocol.PROTO_NC_ITEM_DEPOSIT_ACK, 12, 29},
		{"PROTO_NC_ITEM_WITHDRAW_REQ", protocol.PROTO_NC_ITEM_WITHDRAW_REQ, 12, 30},
		{"PROTO_NC_ITEM_WITHDRAW_ACK", protocol.PROTO_NC_ITEM_WITHDRAW_ACK, 12, 31},
		{"PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ", protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ, 12, 32},
		{"PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK", protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK, 12, 33},
		{"PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ", protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ, 12, 34},
		{"PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK", protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK, 12, 35},
		{"PROTO_NC_ITEM_DISMANTLE_REQ", protocol.PROTO_NC_ITEM_DISMANTLE_REQ, 12, 42},
		{"PROTO_NC_ITEM_DISMANTLE_ACK", protocol.PROTO_NC_ITEM_DISMANTLE_ACK, 12, 43},
		{"PROTO_NC_ITEM_REWARDINVENOPEN_REQ", protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ, 12, 44},
		{"PROTO_NC_ITEM_REWARDINVENOPEN_ACK", protocol.PROTO_NC_ITEM_REWARDINVENOPEN_ACK, 12, 45},
		{"PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ", protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ, 12, 74},
		{"PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK", protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK, 12, 75},
		{"PROTO_NC_ITEM_MOVER_UPGRADE_ACK", protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK, 12, 82},
		{"PROTO_NC_ITEM_REPURCHASE_REQ", protocol.PROTO_NC_ITEM_REPURCHASE_REQ, 12, 89},
		{"PROTO_NC_ITEM_REPURCHASE_ACK", protocol.PROTO_NC_ITEM_REPURCHASE_ACK, 12, 90},
		{"PROTO_NC_ITEM_MAPLINK_SCROLL_REQ", protocol.PROTO_NC_ITEM_MAPLINK_SCROLL_REQ, 12, 103},
		{"PROTO_NC_ITEM_MAPLINK_SCROLL_ACK", protocol.PROTO_NC_ITEM_MAPLINK_SCROLL_ACK, 12, 104},
		{"PROTO_NC_ITEM_MIX_ITEM_REQ", protocol.PROTO_NC_ITEM_MIX_ITEM_REQ, 12, 105},
		{"PROTO_NC_ITEM_MIX_ITEM_ACK", protocol.PROTO_NC_ITEM_MIX_ITEM_ACK, 12, 106},
		{"PROTO_NC_MENU_OPENSTORAGE_CMD", protocol.PROTO_NC_MENU_OPENSTORAGE_CMD, 15, 8},
		{"PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ", protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ, 16, 37},
		{"PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK", protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK, 16, 38},
		{"PROTO_NC_TRADE_PROPOSE_REQ", protocol.PROTO_NC_TRADE_PROPOSE_REQ, 19, 1},
		{"PROTO_NC_TRADE_UPBOARD_REQ", protocol.PROTO_NC_TRADE_UPBOARD_REQ, 19, 13},
		{"PROTO_NC_TRADE_OPPOSITUPBOARD_CMD", protocol.PROTO_NC_TRADE_OPPOSITUPBOARD_CMD, 19, 16},
		{"PROTO_NC_TRADE_TRADECOMPLETE_CMD", protocol.PROTO_NC_TRADE_TRADECOMPLETE_CMD, 19, 36},
		{"PROTO_NC_BOOTH_OPEN_REQ", protocol.PROTO_NC_BOOTH_OPEN_REQ, 26, 1},
		{"PROTO_NC_BOOTH_ENTRY_BUY_ACK", protocol.PROTO_NC_BOOTH_ENTRY_BUY_ACK, 26, 9},
		{"PROTO_NC_BOOTH_REFRESH_SELL_ACK", protocol.PROTO_NC_BOOTH_REFRESH_SELL_ACK, 26, 11},
		{"PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK", protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK, 26, 20},
		{"PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK", protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK, 26, 22},
		{"PROTO_NC_GUILD_STORAGEOPEN_ACK", protocol.PROTO_NC_GUILD_STORAGEOPEN_ACK, 29, 146},
		{"PROTO_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD", protocol.PROTO_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD, 29, 160},
		{"PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK", protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK, 38, 110},
	}
	for _, tc := range opcodeChecks {
		want := (tc.dept << 10) | tc.cmd
		if tc.got != want {
			t.Errorf("%s = 0x%04X, esperado=(%d<<10)|%d=0x%04X", tc.name, tc.got, tc.dept, tc.cmd, want)
		}
	}

	// Verificar sizeof empaquetado de structs NC_ITEM en pkg/protocol
	sizeChecks := []struct {
		name string
		got  int
		want int
	}{
		{"PROTO_NC_ITEM_CELLCHANGE_CMD", binary.Size(protocol.PROTO_NC_ITEM_CELLCHANGE_CMD_DATA{}), 6},
		{"PROTO_NC_ITEM_EQUIPCHANGE_CMD", binary.Size(protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD_DATA{}), 5},
		{"PROTO_NC_ITEM_BUY_REQ", binary.Size(protocol.PROTO_NC_ITEM_BUY_REQ_DATA{}), 6},
		{"PROTO_NC_ITEM_SELL_REQ", binary.Size(protocol.PROTO_NC_ITEM_SELL_REQ_DATA{}), 5},
		{"PROTO_NC_ITEM_REPURCHASE_REQ", binary.Size(protocol.PROTO_NC_ITEM_REPURCHASE_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_DROP_REQ", binary.Size(protocol.PROTO_NC_ITEM_DROP_REQ_DATA{}), 14},
		{"PROTO_NC_ITEM_PICK_REQ", binary.Size(protocol.PROTO_NC_ITEM_PICK_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_PICK_ACK", binary.Size(protocol.PROTO_NC_ITEM_PICK_ACK_DATA{}), 10},
		{"PROTO_NC_ITEM_RELOC_REQ", binary.Size(protocol.PROTO_NC_ITEM_RELOC_REQ_DATA{}), 4},
		{"PROTO_NC_ITEM_SPLIT_REQ", binary.Size(protocol.PROTO_NC_ITEM_SPLIT_REQ_DATA{}), 8},
		{"PROTO_NC_ITEM_EQUIP_REQ", binary.Size(protocol.PROTO_NC_ITEM_EQUIP_REQ_DATA{}), 1},
		{"PROTO_NC_ITEM_RINGEQUIP_REQ", binary.Size(protocol.PROTO_NC_ITEM_RINGEQUIP_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_UNEQUIP_REQ", binary.Size(protocol.PROTO_NC_ITEM_UNEQUIP_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_DISMANTLE_REQ", binary.Size(protocol.PROTO_NC_ITEM_DISMANTLE_REQ_DATA{}), 1},
		{"PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK", binary.Size(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK_DATA{}), 2},
		{"PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ", binary.Size(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK_fixed", binary.Size(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK_DATA{}), 5},
		{"PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ", binary.Size(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ_DATA{}), 4},
		{"PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK", binary.Size(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK_DATA{}), 6},
		{"PROTO_NC_ITEM_REWARDINVENOPEN_REQ", binary.Size(protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ_DATA{}), 2},
		{"PROTO_NC_ITEM_REWARDINVENOPEN_ACK_fixed", binary.Size(protocol.PROTO_NC_ITEM_REWARDINVENOPEN_ACK_DATA{}), 1},
		{"PROTO_NC_ITEM_MOVER_UPGRADE_ACK", binary.Size(protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK_DATA{}), 108},
	}
	for _, sc := range sizeChecks {
		if sc.got != sc.want {
			t.Errorf("%s binary.Size = %d, esperado=%d", sc.name, sc.got, sc.want)
		}
	}

	// Verificar por AST que HandlePacket despacha el 100% de C2S oficiales de Zone.md para NC_ITEM, NC_TRADE, NC_BOOTH, NC_CHARSAVE y NC_ITEMDB
	rootPath := findRepoRoot(t)
	zoneMDBytes, err := os.ReadFile(filepath.Join(rootPath, "docs", "Servicios", "Zone.md"))
	if err != nil {
		t.Fatalf("Error leyendo docs/Servicios/Zone.md: %v", err)
	}
	secRe := regexp.MustCompile(`(?is)protocolstore.*?\n(.*?)(?:\n## |\z)`)
	secMatch := secRe.FindStringSubmatch(string(zoneMDBytes))
	if len(secMatch) < 2 {
		t.Fatalf("No se encontró protocolstore en docs/Servicios/Zone.md")
	}
	opRe := regexp.MustCompile(`PROTO_NC_[A-Z0-9_]+`)
	officialOps := make(map[string]bool)
	for _, op := range opRe.FindAllString(secMatch[1], -1) {
		officialOps[op] = true
	}

	handlerPath := filepath.Join(rootPath, "pkg", "server", "zone", "handler.go")
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, handlerPath, nil, 0)
	if err != nil {
		t.Fatalf("Error parseando AST handler.go: %v", err)
	}
	handledOps := make(map[string]bool)
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
					handledOps[sel.Sel.Name] = true
				}
			}
			return true
		})
		return false
	})

	prefixes := []string{"PROTO_NC_ITEM_", "PROTO_NC_TRADE_", "PROTO_NC_BOOTH_", "PROTO_NC_CHARSAVE_", "PROTO_NC_ITEMDB_"}
	var missing []string
	counts := make(map[string]int)
	for op := range officialOps {
		for _, p := range prefixes {
			if strings.HasPrefix(op, p) {
				counts[p]++
				if !handledOps[op] {
					missing = append(missing, op)
				}
			}
		}
	}
	if len(missing) > 0 {
		t.Fatalf("Faltan handlers C2S en Zone HandlePacket: %v", missing)
	}
	if counts["PROTO_NC_ITEM_"] != 48 || counts["PROTO_NC_TRADE_"] != 10 || counts["PROTO_NC_BOOTH_"] != 8 || counts["PROTO_NC_CHARSAVE_"] != 1 || counts["PROTO_NC_ITEMDB_"] != 1 {
		t.Fatalf("Conteo oficial inesperado: %+v", counts)
	}
}

// TestFase19_Inventory_Charged_Reward_Economy_Booth_ClusterE2E valida en clúster TCP real
// (LoginService + WorldService + ZoneService + MySQL) la fidelidad binaria 1:1 de Zone.exe:
//  1) Serialización PROTO_NC_CHAR_CLIENT_ITEM_CMD (0x1047), CELLCHANGE_CMD (0x3001), EQUIPCHANGE_CMD (0x3002),
//     RELOC (0x300B), SPLIT (0x300D, 8B uint32 lot), EQUIP/UNEQUIP + 0x1034, DROP (0x3007, 14B) + PICK (0x3009, 0x0341),
//     AUTO_ARRANGE_INVEN (0x304A -> 0x304B 0x30A8) y DISMANTLE (0x302A -> 0x302B).
//  2) Charged Inventory: 0x3020 -> 0x3021 (0x1041) y 0x3022 -> Box 9 tItem + 0x3001 CELLCHANGE + 0x3023 (0x1041),
//     ítem retirable 100% movible vía RELOC, y rechazo 0x1045 si Box 9 está lleno.
//  3) Reward Inventory: 0x302C -> 0x302D (Box 2) + extracción Box 2->9 vía RELOC (y bloqueo Box *->2 con err=0x0243).
//  4) Economía / Forja / SafeTrade / Booth: BUY (0x0201), SELL (0x0381), REPURCHASE (0x3668), OPENSTORAGE_CMD (0x3C08),
//     UI_STATE_SAVE (0x4026), TRADE_OPPOSITUPBOARD_CMD (0x4C10), BOOTH_ENTRY/REFRESH (0x0A41), BOOTH_SEARCH (0x31C8/0x31D3).
func TestFase19_Inventory_Charged_Reward_Economy_Booth_ClusterE2E(t *testing.T) {
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

	const (
		loginPort = 19512
		worldPort = 19516
		zonePort  = 19526
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
		Name:     "Fase19Cluster",
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
		t.Fatalf("Error UpsertZoneNode: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort = ?", zonePort)
	}()

	glDispatcher := database.NewGameLogDispatcher(gameLogRepo, 1000, 2)
	glDispatcher.Start()
	defer glDispatcher.Close(t.Context())

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

	const (
		userNo1 = 998919
		userNo2 = 998920
	)
	ts := time.Now().Unix() % 100000
	pass := "f19pass"
	sum := md5.Sum([]byte(pass))
	passHash := hex.EncodeToString(sum[:])

	cleanupUser := func(uNo int) {
		_, _ = db.Exec("DELETE FROM tChargeItem WHERE userNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tLogout WHERE nUserNo = ?", uNo)
		if rows, qErr := db.Query("SELECT nCharNo FROM tCharacter WHERE nUserNo = ?", uNo); qErr == nil {
			var charIDs []int
			for rows.Next() {
				var cid int
				if rows.Scan(&cid) == nil {
					charIDs = append(charIDs, cid)
				}
			}
			_ = rows.Close()
			for _, cid := range charIDs {
				_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", cid)
				_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", cid)
			}
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", uNo)
	}
	cleanupUser(userNo1)
	cleanupUser(userNo2)
	defer func() {
		cleanupUser(userNo1)
		cleanupUser(userNo2)
	}()

	connectPlayerToZone := func(uNo int, uName, cName string) (net.Conn, *auditMockClient, int, []*protocol.Packet) {
		if _, err := db.Exec(`
			INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
			VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
		`, uNo, uName, passHash); err != nil {
			t.Fatalf("Error insertando tUser %d: %v", uNo, err)
		}

		lConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
		if err != nil {
			t.Fatalf("Dial LoginService: %v", err)
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
		lReq.WriteString(uName, 260)
		lReq.WriteString(passHash, 36)
		lReq.Fill(20, 0)
		_ = lClient.sendPacket(lReq)
		_, _ = auditReadServerPacket(lConn)

		wSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
		wSelReq.WriteByte(0)
		_ = lClient.sendPacket(wSelReq)
		wSelAck, err := auditReadServerPacket(lConn)
		if err != nil || wSelAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
			t.Fatalf("WORLDSELECT falló: %v / %v", wSelAck, err)
		}
		_, _ = wSelAck.Reader.ReadByte()
		_, _ = wSelAck.Reader.ReadString(16)
		_, _ = wSelAck.Reader.ReadUint16()
		tokenBytes, _ := wSelAck.Reader.ReadBytes(64)
		worldToken := strings.TrimRight(string(tokenBytes), "\x00")
		_ = lConn.Close()

		wConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
		if err != nil {
			t.Fatalf("Dial WorldService: %v", err)
		}
		defer wConn.Close()
		wHS, _ := auditReadServerPacket(wConn)
		wSeed, _ := wHS.Reader.ReadInt16()
		wCrypto, _ := crypto.NewNetCrypto(wSeed)
		wClient := &auditMockClient{conn: wConn, crypto: wCrypto}

		wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
		wlReq.WriteString(uName, 256)
		wlReq.WriteString(worldToken, 64)
		_ = wClient.sendPacket(wlReq)
		if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
			t.Fatalf("LOGINWORLD falló: %v / %v", wlAck, err)
		}

		createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
		createReq.WriteByte(0)
		createReq.WriteString(cName, 20)
		createReq.WriteByte(byte((1 << 2) | (1 << 7)))
		createReq.WriteByte(1)
		createReq.WriteByte(1)
		createReq.WriteByte(4)
		_ = wClient.sendPacket(createReq)
		if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
			t.Fatalf("AVATAR_CREATE falló: %v / %v", createAck, err)
		}

		createdChar, err := charRepo.GetCharacterByName(cName)
		if err != nil || createdChar == nil {
			t.Fatalf("Personaje '%s' no encontrado: %v", cName, err)
		}
		_ = charRepo.AddCharMoney(createdChar.ID, 500000)

		cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
		cSelReq.WriteByte(0)
		_ = wClient.sendPacket(cSelReq)
		if cSelAck, err := auditReadServerPacket(wConn); err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
			t.Fatalf("CHAR_LOGIN falló: %v / %v", cSelAck, err)
		}
		_, _ = auditReadServerPacket(wConn)
		_, _ = auditReadServerPacket(wConn)
		_, _ = auditReadServerPacket(wConn)

		var dbRandomID uint16
		if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", cName).Scan(&dbRandomID); err != nil {
			t.Fatalf("Token tActiveZoneTransfer no encontrado: %v", err)
		}

		zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
		if err != nil {
			t.Fatalf("Dial ZoneService: %v", err)
		}
		zHS, _ := auditReadServerPacket(zConn)
		zSeed, _ := zHS.Reader.ReadInt16()
		zCrypto, _ := crypto.NewNetCrypto(zSeed)
		zClient := &auditMockClient{conn: zConn, crypto: zCrypto}

		mapReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
		mapReq.WriteUint16(dbRandomID)
		mapReq.WriteString(cName, 20)
		sums := zoneDataStore.ChkEngine.GetChecksums()
		for i := 0; i < 49; i++ {
			mapReq.WriteBytes(sums[i][:])
		}
		_ = zClient.sendPacket(mapReq)

		var chunkPkts []*protocol.Packet
		for {
			pkt, rErr := auditReadServerPacket(zConn)
			if rErr != nil {
				t.Fatalf("Error leyendo MAP_LOGIN: %v", rErr)
			}
			chunkPkts = append(chunkPkts, pkt)
			if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
				break
			}
		}
		_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
		return zConn, zClient, createdChar.ID, chunkPkts
	}

	collectPackets := func(zConn net.Conn, wait time.Duration) []*protocol.Packet {
		var pkts []*protocol.Packet
		_ = zConn.SetReadDeadline(time.Now().Add(wait))
		for {
			pkt, rErr := auditReadServerPacket(zConn)
			if rErr != nil {
				break
			}
			pkts = append(pkts, pkt)
		}
		_ = zConn.SetReadDeadline(time.Time{})
		return pkts
	}

	findAllPkts := func(pkts []*protocol.Packet, op uint16) []*protocol.Packet {
		var out []*protocol.Packet
		for _, p := range pkts {
			if p.OpCode == op {
				out = append(out, p)
			}
		}
		return out
	}

	findPkt := func(pkts []*protocol.Packet, op uint16) *protocol.Packet {
		for _, p := range pkts {
			if p.OpCode == op {
				return p
			}
		}
		return nil
	}

	// Seleccionar ítems reales desde ItemInfo.shn para cubrir consumibles apilables y equipo Fighter
	var potionItemID, swordItemID uint16
	for id := 1; id < 65000; id++ {
		info, ok := zoneDataStore.ItemStore.GetItemInfo(uint16(id))
		if !ok || info == nil {
			continue
		}
		if potionItemID == 0 && (info.Class == 0 || info.Class == 1 || info.Class == 2) && info.MaxLot >= 10 {
			potionItemID = info.ID
		}
		if swordItemID == 0 && (info.Equip == 11 || info.Equip == 10) && info.DemandLv <= 5 && (info.UseClass&(1<<1)) != 0 {
			swordItemID = info.ID
		}
		if potionItemID != 0 && swordItemID != 0 {
			break
		}
	}
	if potionItemID == 0 || swordItemID == 0 {
		t.Fatalf("No se encontraron ítems de prueba en ItemInfo.shn (potion=%d sword=%d)", potionItemID, swordItemID)
	}

	// Conectar Jugador 1
	uName1 := fmt.Sprintf("f19u1_%d", ts)
	cName1 := fmt.Sprintf("F19HeroA%d", ts%1000)
	zConn1, zClient1, charID1, chunkPkts1 := connectPlayerToZone(userNo1, uName1, cName1)
	defer zConn1.Close()
	_ = collectPackets(zConn1, 250*time.Millisecond)

	// 1. Verificar formato PROTO_NC_CHAR_CLIENT_ITEM_CMD (0x1047) en In-Game Chunk
	itemChunkPkts := findAllPkts(chunkPkts1, protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD)
	if len(itemChunkPkts) == 0 {
		t.Fatalf("In-Game Chunk no emitió PROTO_NC_CHAR_CLIENT_ITEM_CMD (0x1047)")
	}
	for _, cp := range itemChunkPkts {
		pl := cp.Raw[2:]
		if len(pl) < 3 {
			t.Fatalf("0x1047 payload muy corto: %dB", len(pl))
		}
		numItems := int(pl[0])
		off := 3
		for i := 0; i < numItems; i++ {
			if off >= len(pl) {
				t.Fatalf("0x1047 truncado en ítem %d (off=%d len=%d)", i, off, len(pl))
			}
			ds := int(pl[off])
			if ds < 4 || off+1+ds > len(pl) {
				t.Fatalf("0x1047 datasize=%d inválido en ítem %d (off=%d len=%d)", ds, i, off, len(pl))
			}
			off += 1 + ds
		}
		if off != len(pl) {
			t.Fatalf("0x1047 longitud residual inválida: parsed=%d total=%d", off, len(pl))
		}
	}

	// 2. Dar ítems al Jugador 1 y probar Charged Inventory (0x3020 -> 0x3022) + Reward Inventory (0x302C -> 0x300B)
	if _, err := accountRepo.InsertChargedItem(&database.ChargedItem{
		UserNo:  userNo1,
		OrderNo: uint32(ts*100 + 1),
		GoodsNo: uint32(potionItemID),
		Amount:  5,
	}); err != nil {
		t.Fatalf("InsertChargedItem falló: %v", err)
	}

	// Abrir Charged Inventory (0x3020) -> 0x3021 err=0x1041
	ciOpen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ)
	ciOpen.WriteUint16(0)
	_ = zClient1.sendPacket(ciOpen)
	ciOpenPkts := collectPackets(zConn1, 150*time.Millisecond)
	ciOpenAck := findPkt(ciOpenPkts, protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK)
	if ciOpenAck == nil || len(ciOpenAck.Raw)-2 < 5 {
		t.Fatalf("CHARGEDINVENOPEN_ACK ausente o truncado: %v", ciOpenAck)
	}
	ciPl := ciOpenAck.Raw[2:]
	ciErr := binary.LittleEndian.Uint16(ciPl[0:2])
	ciCount := binary.LittleEndian.Uint16(ciPl[3:5])
	if ciErr != 0x1041 || ciCount < 1 {
		t.Fatalf("CHARGEDINVENOPEN_ACK err=0x%04X count=%d, esperado err=0x1041 count>=1", ciErr, ciCount)
	}
	orderNo := binary.LittleEndian.Uint32(ciPl[5:9])

	// Retirar ítem del Charged Inventory (0x3022) -> debe entregar Box 9 + 0x3001 CELLCHANGE + 0x3023 (0x1041)
	cwReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ)
	cwReq.WriteUint32(orderNo)
	_ = zClient1.sendPacket(cwReq)
	cwPkts := collectPackets(zConn1, 200*time.Millisecond)
	cwCell := findPkt(cwPkts, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD)
	cwAck := findPkt(cwPkts, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)
	if cwCell == nil {
		t.Fatalf("CHARGED_WITHDRAW no emitió PROTO_NC_ITEM_CELLCHANGE_CMD (0x3001)")
	}
	if cwAck == nil || len(cwAck.Raw)-2 != 6 || binary.LittleEndian.Uint16(cwAck.Raw[6:8]) != 0x1041 {
		t.Fatalf("CHARGED_WITHDRAW_ACK inválido (esperado 6B err=0x1041): %v", cwAck)
	}
	cwCellPl := cwCell.Raw[2:]
	withdrawnLoc := binary.LittleEndian.Uint16(cwCellPl[2:4])
	withdrawnBox := withdrawnLoc >> 10
	withdrawnSlot := withdrawnLoc & 0x03FF
	if withdrawnBox != 9 {
		t.Fatalf("Charged Withdraw ubicó ítem en box=%d, esperado box=9", withdrawnBox)
	}
	if len(cwCellPl) < 6 {
		t.Fatalf("0x3001 CELLCHANGE demasiado corto: %dB", len(cwCellPl))
	}

	// Confirmar que el ítem retirado del Charged Inventory es movible con PROTO_NC_ITEM_RELOC_REQ (0x300B: from, to)
	destSlot := (withdrawnSlot + 1) % 24
	relocReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocReq.WriteUint16(withdrawnLoc)
	relocReq.WriteUint16((9 << 10) | destSlot)
	_ = zClient1.sendPacket(relocReq)
	relocPkts := collectPackets(zConn1, 150*time.Millisecond)
	relocAck := findPkt(relocPkts, protocol.PROTO_NC_ITEM_RELOC_ACK)
	if relocAck == nil || binary.LittleEndian.Uint16(relocAck.Raw[2:4]) != 0x0241 {
		t.Fatalf("RELOC_REQ de ítem Charged falló: %v", relocAck)
	}
	if findPkt(relocPkts, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD) == nil {
		t.Fatalf("RELOC_REQ de ítem Charged no emitió CELLCHANGE_CMD (0x3001)")
	}

	// Probar que mover hacia Box 2 (RewardInven) es rechazado con 0x0243
	relocRewardReject := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocRewardReject.WriteUint16((9 << 10) | destSlot)
	relocRewardReject.WriteUint16((2 << 10) | 0)
	_ = zClient1.sendPacket(relocRewardReject)
	rejPkts := collectPackets(zConn1, 120*time.Millisecond)
	if rejAck := findPkt(rejPkts, protocol.PROTO_NC_ITEM_RELOC_ACK); rejAck == nil || binary.LittleEndian.Uint16(rejAck.Raw[2:4]) != 0x0243 {
		t.Fatalf("RELOC hacia Reward Box 2 debería fallar con 0x0243, obtenido: %v", rejAck)
	}

	// Reward Inventory (Box 2): insertar ítem en Box 2 -> 0x302C -> 0x302D -> extraer vía RELOC Box 2->9
	rewardItem := &database.Item{
		ItemKey:     ts*100000 + 101,
		StorageType: 2,
		Owner:       charID1,
		Storage:     0,
		ItemID:      potionItemID,
		Flags:       2,
		Date:        time.Now(),
	}
	if err := charRepo.CreateItem(rewardItem); err != nil {
		t.Fatalf("CreateItem Box 2 falló: %v", err)
	}
	_ = zClient1.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ))
	rwPkts := collectPackets(zConn1, 150*time.Millisecond)
	rwAck := findPkt(rwPkts, protocol.PROTO_NC_ITEM_REWARDINVENOPEN_ACK)
	if rwAck == nil || len(rwAck.Raw)-2 < 1 || rwAck.Raw[2] < 1 {
		t.Fatalf("REWARDINVENOPEN_ACK ausente o vacío: %v", rwAck)
	}
	extractSlot := (destSlot + 2) % 24
	relocExtract := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocExtract.WriteUint16((2 << 10) | 0)
	relocExtract.WriteUint16((9 << 10) | extractSlot)
	_ = zClient1.sendPacket(relocExtract)
	extPkts := collectPackets(zConn1, 150*time.Millisecond)
	if extAck := findPkt(extPkts, protocol.PROTO_NC_ITEM_RELOC_ACK); extAck == nil || binary.LittleEndian.Uint16(extAck.Raw[2:4]) != 0x0241 {
		t.Fatalf("Extracción Reward Box 2 -> 9 falló: %v", extAck)
	}

	// Probar SPLIT_REQ (0x300D, 8B: from, to, uint32 lot)
	splitDestSlot := (destSlot + 3) % 24
	splitReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_REQ)
	splitReq.WriteUint16((9 << 10) | destSlot)
	splitReq.WriteUint16((9 << 10) | splitDestSlot)
	splitReq.WriteUint32(2)
	_ = zClient1.sendPacket(splitReq)
	splitPkts := collectPackets(zConn1, 150*time.Millisecond)
	if spAck := findPkt(splitPkts, protocol.PROTO_NC_ITEM_SPLIT_ACK); spAck == nil || binary.LittleEndian.Uint16(spAck.Raw[2:4]) != 0x0441 {
		t.Fatalf("SPLIT_REQ (8B uint32 lot) falló: %v", spAck)
	}

	// Probar EQUIP/UNEQUIP y recálculo de stats
	eqItem := &database.Item{
		ItemKey:     ts*100000 + 102,
		StorageType: 9,
		Owner:       charID1,
		Storage:     int16((destSlot + 4) % 24),
		ItemID:      swordItemID,
		Flags:       1,
		Date:        time.Now(),
	}
	if err := charRepo.CreateItem(eqItem); err != nil {
		t.Fatalf("CreateItem arma Box 9 falló: %v", err)
	}
	eqReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_REQ)
	eqReq.WriteByte(byte(eqItem.Storage))
	_ = zClient1.sendPacket(eqReq)
	eqPkts := collectPackets(zConn1, 150*time.Millisecond)
	if findPkt(eqPkts, protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD) == nil {
		t.Fatalf("EQUIP_REQ no emitió PROTO_NC_ITEM_EQUIPCHANGE_CMD (0x3002)")
	}
	if findPkt(eqPkts, protocol.PROTO_NC_ITEM_EQUIP_ACK) == nil {
		t.Fatalf("EQUIP_REQ no emitió PROTO_NC_ITEM_EQUIP_ACK (0x3011)")
	}

	// Probar DROP_REQ (14B) + PICK_REQ (0x0341 Cen)
	dropReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_REQ)
	dropReq.WriteUint16((9 << 10) | splitDestSlot)
	dropReq.WriteUint32(1)
	dropReq.WriteUint32(6888)
	dropReq.WriteUint32(7045)
	_ = zClient1.sendPacket(dropReq)
	dropPkts := collectPackets(zConn1, 150*time.Millisecond)
	if dropAck := findPkt(dropPkts, protocol.PROTO_NC_ITEM_DROP_ACK); dropAck == nil || binary.LittleEndian.Uint16(dropAck.Raw[2:4]) != 0x0301 {
		t.Fatalf("DROP_REQ (14B) falló: %v", dropAck)
	}
	if dropBrief := findPkt(dropPkts, protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD); dropBrief != nil && len(dropBrief.Raw)-2 >= 2 {
		dropHandle := binary.LittleEndian.Uint16(dropBrief.Raw[2:4])
		pickReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
		pickReq.WriteUint16(dropHandle)
		_ = zClient1.sendPacket(pickReq)
		pickPkts := collectPackets(zConn1, 150*time.Millisecond)
		if pickAck := findPkt(pickPkts, protocol.PROTO_NC_ITEM_PICK_ACK); pickAck == nil || len(pickAck.Raw)-2 != 10 {
			t.Fatalf("PICK_ACK inválido o != 10B: %v", pickAck)
		}
	}

	// Probar AUTO_ARRANGE_INVEN_REQ (0x304A) -> 0x304B (err=0x30A8)
	_ = zClient1.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ))
	arrPkts := collectPackets(zConn1, 150*time.Millisecond)
	if arrAck := findPkt(arrPkts, protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK); arrAck == nil || len(arrAck.Raw)-2 != 2 || binary.LittleEndian.Uint16(arrAck.Raw[2:4]) != 0x30A8 {
		t.Fatalf("AUTO_ARRANGE_INVEN_ACK inválido (esperado 2B 0x30A8): %v", arrAck)
	}

	// Probar CHARSAVE_UI_STATE_SAVE_REQ (0x4025 -> 0x4026, 2B)
	uiReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ)
	uiReq.WriteByte(1)
	_ = zClient1.sendPacket(uiReq)
	uiPkts := collectPackets(zConn1, 120*time.Millisecond)
	if uiAck := findPkt(uiPkts, protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK); uiAck == nil || len(uiAck.Raw)-2 != 2 {
		t.Fatalf("CHARSAVE_UI_STATE_SAVE_ACK inválido: %v", uiAck)
	}

	// Probar Storage open (0x3028 -> 0x3C08) y SELL/REPURCHASE/BOOTH/TRADE códigos oficiales
	stReq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_OPENSTORAGEPAGE_REQ)
	stReq.WriteByte(0)
	_ = zClient1.sendPacket(stReq)
	stPkts := collectPackets(zConn1, 150*time.Millisecond)
	if findPkt(stPkts, protocol.PROTO_NC_MENU_OPENSTORAGE_CMD) == nil {
		t.Fatalf("OPENSTORAGEPAGE_REQ no emitió PROTO_NC_MENU_OPENSTORAGE_CMD (0x3C08)")
	}

	// Probar Charged Withdraw rechazo 0x1045 cuando Box 9 está completamente lleno
	if _, err := accountRepo.InsertChargedItem(&database.ChargedItem{
		UserNo:  userNo1,
		OrderNo: uint32(ts*100 + 2),
		GoodsNo: uint32(potionItemID),
		Amount:  1,
	}); err != nil {
		t.Fatalf("InsertChargedItem #2 falló: %v", err)
	}
	existing9, _ := charRepo.GetItemsByStorage(charID1, 9)
	occupied9 := make(map[int16]bool)
	for _, it := range existing9 {
		occupied9[it.Storage] = true
	}
	for s := int16(0); s < 96; s++ {
		if !occupied9[s] {
			_ = charRepo.CreateItem(&database.Item{
				ItemKey:     ts*100000 + 200 + int64(s),
				StorageType: 9,
				Owner:       charID1,
				Storage:     s,
				ItemID:      potionItemID,
				Flags:       1,
				Date:        time.Now(),
			})
		}
	}
	ciOpen2 := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ)
	ciOpen2.WriteUint16(0)
	_ = zClient1.sendPacket(ciOpen2)
	ciOpen2Pkts := collectPackets(zConn1, 150*time.Millisecond)
	if ci2Ack := findPkt(ciOpen2Pkts, protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK); ci2Ack != nil && len(ci2Ack.Raw)-2 >= 9 {
		order2 := binary.LittleEndian.Uint32(ci2Ack.Raw[7:11])
		cwFull := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ)
		cwFull.WriteUint32(order2)
		_ = zClient1.sendPacket(cwFull)
		cwFullPkts := collectPackets(zConn1, 150*time.Millisecond)
		if fullAck := findPkt(cwFullPkts, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK); fullAck == nil || binary.LittleEndian.Uint16(fullAck.Raw[6:8]) != 0x1045 {
			t.Fatalf("CHARGED_WITHDRAW con Box 9 lleno debería devolver 0x1045, obtenido: %v", fullAck)
		}
	}

	// Probar BOOTH OPEN + ENTRY + SEARCH con 2 jugadores TCP en ZoneService
	uName2 := fmt.Sprintf("f19u2_%d", ts)
	cName2 := fmt.Sprintf("F19HeroB%d", ts%1000)
	zConn2, zClient2, charID2, _ := connectPlayerToZone(userNo2, uName2, cName2)
	defer zConn2.Close()
	_ = collectPackets(zConn2, 200*time.Millisecond)
	_ = charID2

	// Abrir kiosco SELL en Player 2 (Signboard[30] + Flag(1B) + Items)
	bOpen := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_OPEN_REQ)
	bOpen.WriteString("F19Shop", 30)
	bOpen.WriteByte(0) // isSell=0 -> Flag=0
	bOpen.WriteByte(0) // 0 items
	_ = zClient2.sendPacket(bOpen)
	bOpenPkts := collectPackets(zConn2, 150*time.Millisecond)
	p1AOIPkts := collectPackets(zConn1, 150*time.Millisecond)
	if bOpenAck := findPkt(bOpenPkts, protocol.PROTO_NC_BOOTH_OPEN_ACK); bOpenAck != nil && len(bOpenAck.Raw)-2 >= 2 {
		if soOpen := findPkt(p1AOIPkts, protocol.PROTO_NC_BOOTH_SOMEONEOPEN_CMD); soOpen != nil && len(soOpen.Raw)-2 >= 2 {
			boothHandle := binary.LittleEndian.Uint16(soOpen.Raw[2:4])
			bEnt := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ENTRY_REQ)
			bEnt.WriteUint16(boothHandle)
			_ = zClient1.sendPacket(bEnt)
			bEntPkts := collectPackets(zConn1, 150*time.Millisecond)
			if bEntAck := findPkt(bEntPkts, protocol.PROTO_NC_BOOTH_ENTRY_SELL_ACK); bEntAck != nil {
				if errCode := binary.LittleEndian.Uint16(bEntAck.Raw[2:4]); errCode != 0x0A41 {
					t.Fatalf("BOOTH_ENTRY_SELL_ACK err=0x%04X, esperado=0x0A41", errCode)
				}
			}
		}
	}

	// Buscar kioscos categorizados (0x6813: uint32 category + uint8 onlyMyClass + uint8 isOpeningReq -> 0x6814 err=0x31C8)
	bSearch := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ)
	bSearch.WriteUint32(0)
	bSearch.WriteByte(0)
	bSearch.WriteByte(1)
	_ = zClient1.sendPacket(bSearch)
	bSearchPkts := collectPackets(zConn1, 150*time.Millisecond)
	if bsAck := findPkt(bSearchPkts, protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK); bsAck == nil || len(bsAck.Raw)-2 < 3 || binary.LittleEndian.Uint16(bsAck.Raw[2:4]) != 0x31C8 {
		t.Fatalf("BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK inválido (esperado err=0x31C8): %v", bsAck)
	}

	bPosMissing := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ)
	bPosMissing.WriteUint32(999999)
	bPosMissing.WriteUint16(65000)
	_ = zClient1.sendPacket(bPosMissing)
	bPosPkts := collectPackets(zConn1, 150*time.Millisecond)
	if bpAck := findPkt(bPosPkts, protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK); bpAck == nil || len(bpAck.Raw)-2 < 2 || binary.LittleEndian.Uint16(bpAck.Raw[2:4]) != 0x31D3 {
		t.Fatalf("BOOTH_SEARCH_BOOTH_POSITION_ACK (not found) inválido (esperado err=0x31D3): %v", bpAck)
	}
}
