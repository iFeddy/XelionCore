package audit_test

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"path/filepath"
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

// TestFase18_BinaryProtocolSizes_Dept7_9_18_20 verifica al byte las estructuras
// y códigos OpCode oficiales ((dept << 10) | command) para los departamentos:
//   - Dept 7 (NC_BRIEFINFO): Players (235B), Mobs/NPCs/Gates (149B, máx. 28/lote),
//     Drops (15B), MagicField (16B), Mover (139B), Door (48B).
//   - Dept 9 (NC_BAT): TARGETINFO (30B), ASSIST_ACK (2B), SWING_DAMAGE (16B),
//     SOMEONESWING_DAMAGE (13B), REALLYKILL (4B), EXPGAIN (6B), LEVELUP (235B),
//     SKILLBASH_HIT_OBJ/FLD/DAMAGE, ABSTATESET/RESET (6B), ABSTATEINFORM (8B), ERASE_ACK (2B).
//   - Dept 18 (NC_SKILL): EMPOWALLOC_REQ (6B), EMPOWALLOC_ACK (4B), EMPOWPOINT_CMD (1B),
//     PRODUCTFIELD_ACK (4B), UNLEARN_REQ (3B), UNLEARN_ACK (4B), REVIVE_CMD (12B), SOMEONEREVIVE_CMD (14B).
//   - Dept 20 (NC_SOULSTONE — Dept 20 en Zone.pdb): HP/SP_BUY_REQ/ACK (2B), BUYFAIL_ACK (2B), USEFAIL/SUC (0B), SOMEONEUSE (2B).
func TestFase18_BinaryProtocolSizes_Dept7_9_18_20(t *testing.T) {
	if zone.BriefInfoMaxBatch != 28 {
		t.Fatalf("BriefInfoMaxBatch esperado=28 (0x1C en Zone.exe), obtenido=%d", zone.BriefInfoMaxBatch)
	}

	// 1. Verificar fórmula exacta OpCode = (dept << 10) | command
	opcodeChecks := []struct {
		name string
		got  uint16
		dept uint16
		cmd  uint16
	}{
		{"PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD", protocol.PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD, 7, 2},
		{"PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD", protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD, 7, 6},
		{"PROTO_NC_BRIEFINFO_CHARACTER_CMD", protocol.PROTO_NC_BRIEFINFO_CHARACTER_CMD, 7, 7},
		{"PROTO_NC_BRIEFINFO_REGENMOB_CMD", protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD, 7, 8},
		{"PROTO_NC_BRIEFINFO_MOB_CMD", protocol.PROTO_NC_BRIEFINFO_MOB_CMD, 7, 9},
		{"PROTO_NC_BRIEFINFO_DROPEDITEM_CMD", protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD, 7, 10},
		{"PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD", protocol.PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD, 7, 11},
		{"PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD", protocol.PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD, 7, 12},
		{"PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD", protocol.PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD, 7, 13},
		{"PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD", protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD, 7, 14},
		{"PROTO_NC_BRIEFINFO_BUILDDOOR_CMD", protocol.PROTO_NC_BRIEFINFO_BUILDDOOR_CMD, 7, 15},
		{"PROTO_NC_BRIEFINFO_DOOR_CMD", protocol.PROTO_NC_BRIEFINFO_DOOR_CMD, 7, 16},
		{"PROTO_NC_BRIEFINFO_REGENMOVER_CMD", protocol.PROTO_NC_BRIEFINFO_REGENMOVER_CMD, 7, 26},
		{"PROTO_NC_BRIEFINFO_MOVER_CMD", protocol.PROTO_NC_BRIEFINFO_MOVER_CMD, 7, 27},
		{"PROTO_NC_BAT_TARGETTING_REQ", protocol.PROTO_NC_BAT_TARGETTING_REQ, 9, 1},
		{"PROTO_NC_BAT_TARGETINFO_CMD", protocol.PROTO_NC_BAT_TARGETINFO_CMD, 9, 2},
		{"PROTO_NC_BAT_UNTARGET_REQ", protocol.PROTO_NC_BAT_UNTARGET_REQ, 9, 8},
		{"PROTO_NC_BAT_EXPGAIN_CMD", protocol.PROTO_NC_BAT_EXPGAIN_CMD, 9, 11},
		{"PROTO_NC_BAT_LEVELUP_CMD", protocol.PROTO_NC_BAT_LEVELUP_CMD, 9, 12},
		{"PROTO_NC_BAT_SUMEONELEVELUP_CMD", protocol.PROTO_NC_BAT_SUMEONELEVELUP_CMD, 9, 13},
		{"PROTO_NC_BAT_HPCHANGE_CMD", protocol.PROTO_NC_BAT_HPCHANGE_CMD, 9, 14},
		{"PROTO_NC_BAT_SPCHANGE_CMD", protocol.PROTO_NC_BAT_SPCHANGE_CMD, 9, 15},
		{"PROTO_NC_BAT_ABSTATESET_CMD", protocol.PROTO_NC_BAT_ABSTATESET_CMD, 9, 39},
		{"PROTO_NC_BAT_ABSTATERESET_CMD", protocol.PROTO_NC_BAT_ABSTATERESET_CMD, 9, 40},
		{"PROTO_NC_BAT_ABSTATEINFORM_CMD", protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD, 9, 41},
		{"PROTO_NC_BAT_BASHSTART_CMD", protocol.PROTO_NC_BAT_BASHSTART_CMD, 9, 43},
		{"PROTO_NC_BAT_BASHSTOP_CMD", protocol.PROTO_NC_BAT_BASHSTOP_CMD, 9, 50},
		{"PROTO_NC_BAT_ASSIST_REQ", protocol.PROTO_NC_BAT_ASSIST_REQ, 9, 62},
		{"PROTO_NC_BAT_ASSIST_ACK", protocol.PROTO_NC_BAT_ASSIST_ACK, 9, 63},
		{"PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ", protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ, 9, 64},
		{"PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ", protocol.PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ, 9, 65},
		{"PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ", protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ, 9, 68},
		{"PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK", protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK, 9, 69},
		{"PROTO_NC_BAT_SWING_START_CMD", protocol.PROTO_NC_BAT_SWING_START_CMD, 9, 71},
		{"PROTO_NC_BAT_SWING_DAMAGE_CMD", protocol.PROTO_NC_BAT_SWING_DAMAGE_CMD, 9, 72},
		{"PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD", protocol.PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD, 9, 73},
		{"PROTO_NC_BAT_REALLYKILL_CMD", protocol.PROTO_NC_BAT_REALLYKILL_CMD, 9, 74},
		{"PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD", protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD, 9, 78},
		{"PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD", protocol.PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD, 9, 82},
		{"PROTO_NC_BAT_ABSTATE_ERASE_REQ", protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ, 9, 84},
		{"PROTO_NC_BAT_ABSTATE_ERASE_ACK", protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK, 9, 85},
		{"PROTO_NC_BAT_LPCHANGE_CMD", protocol.PROTO_NC_BAT_LPCHANGE_CMD, 9, 95},
		{"PROTO_NC_SKILL_EMPOWPOINT_CMD", protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD, 18, 16},
		{"PROTO_NC_SKILL_EMPOWALLOC_REQ", protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ, 18, 17},
		{"PROTO_NC_SKILL_EMPOWALLOC_ACK", protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK, 18, 20},
		{"PROTO_NC_SKILL_REVIVE_CMD", protocol.PROTO_NC_SKILL_REVIVE_CMD, 18, 30},
		{"PROTO_NC_SKILL_SOMEONEREVIVE_CMD", protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD, 18, 31},
		{"PROTO_NC_SKILL_PRODUCTFIELD_REQ", protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ, 18, 34},
		{"PROTO_NC_SKILL_PRODUCTFIELD_ACK", protocol.PROTO_NC_SKILL_PRODUCTFIELD_ACK, 18, 35},
		{"PROTO_NC_SKILL_UNLEARN_REQ", protocol.PROTO_NC_SKILL_UNLEARN_REQ, 18, 36},
		{"PROTO_NC_SKILL_UNLEARN_ACK", protocol.PROTO_NC_SKILL_UNLEARN_ACK, 18, 39},
		{"PROTO_NC_SKILL_REPLYREVIVE_CMD", protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD, 18, 42},
		{"PROTO_NC_SKILL_ITEMACTIONCOOLTIME_CMD", protocol.PROTO_NC_SKILL_ITEMACTIONCOOLTIME_CMD, 18, 44},
		{"PROTO_NC_SOULSTONE_HP_BUY_REQ", protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ, 20, 1},
		{"PROTO_NC_SOULSTONE_SP_BUY_REQ", protocol.PROTO_NC_SOULSTONE_SP_BUY_REQ, 20, 2},
		{"PROTO_NC_SOULSTONE_HP_BUY_ACK", protocol.PROTO_NC_SOULSTONE_HP_BUY_ACK, 20, 3},
		{"PROTO_NC_SOULSTONE_SP_BUY_ACK", protocol.PROTO_NC_SOULSTONE_SP_BUY_ACK, 20, 4},
		{"PROTO_NC_SOULSTONE_BUYFAIL_ACK", protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK, 20, 5},
		{"PROTO_NC_SOULSTONE_USEFAIL_ACK", protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK, 20, 6},
		{"PROTO_NC_SOULSTONE_HP_USE_REQ", protocol.PROTO_NC_SOULSTONE_HP_USE_REQ, 20, 7},
		{"PROTO_NC_SOULSTONE_HP_USESUC_ACK", protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK, 20, 8},
		{"PROTO_NC_SOULSTONE_SP_USE_REQ", protocol.PROTO_NC_SOULSTONE_SP_USE_REQ, 20, 9},
		{"PROTO_NC_SOULSTONE_SP_USESUC_ACK", protocol.PROTO_NC_SOULSTONE_SP_USESUC_ACK, 20, 10},
		{"PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD", protocol.PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD, 20, 11},
		{"PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD", protocol.PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD, 20, 12},
	}
	for _, tc := range opcodeChecks {
		want := (tc.dept << 10) | tc.cmd
		if tc.got != want {
			t.Errorf("%s = 0x%04X, esperado (dept=%d<<10)|%d = 0x%04X", tc.name, tc.got, tc.dept, tc.cmd, want)
		}
	}

	// 2. Verificar tamaños exactos de entradas NC_BRIEFINFO (Dept 7)
	player := zone.NewZonePlayer(nil, &database.Character{
		Name: "BriefHero", Job: 1, Level: 12, XPos: 6888, YPos: 7045,
		TCharacter: database.TCharacter{CharNo: 100},
	}, nil, nil)
	player.MapObjectID = 100
	player.State = 1
	charPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
	player.WriteCharacterDisplay(charPW)
	if got := len(charPW.Data()) - 2; got != 235 {
		t.Fatalf("PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD payload=%dB, esperado=235B (0BriefInformationCharacter)", got)
	}
	if typeVal := charPW.Data()[2+194]; typeVal != 2 {
		t.Fatalf("CHARBRIEFINFO.type en offset 194 = %d, esperado=2 (Zone.exe 0x005490F5)", typeVal)
	}

	// Variante Tent/Camp (State == 4 && TentID > 0) debe mantener 235B exactos
	player.State = 4
	player.TentID = 7001
	tentPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
	player.WriteCharacterDisplay(tentPW)
	if got := len(tentPW.Data()) - 2; got != 235 {
		t.Fatalf("LOGINCHARACTER Tent variant payload=%dB, esperado=235B", got)
	}

	mob := zone.NewZoneMonster(200, 1102, "Slime", "RouN", 6900, 7060, 180, 3, 120, 0, 2, 5, 0, 1500, 30, 80, 25)
	mobPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD)
	mob.WriteRegenMob(mobPW, 1)
	if got := len(mobPW.Data()) - 2; got != 149 {
		t.Fatalf("PROTO_NC_BRIEFINFO_REGENMOB_CMD Mob payload=%dB, esperado=149B", got)
	}
	if mobPW.Data()[2+2] != 0 || mobPW.Data()[2+13] != 90 {
		t.Fatalf("Mob mode=%d dir=%d, esperado mode=0 dir=90 (180/2)", mobPW.Data()[2+2], mobPW.Data()[2+13])
	}

	gateNPC := &zone.ZoneNPC{
		Handle:      201,
		MobID:       4501,
		X:           6950,
		Y:           7100,
		Direct:      90,
		IsGate:      true,
		GateToWhere: "Rou",
	}
	gatePW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD)
	gateNPC.WriteRegenMob(gatePW)
	gatePayload := gatePW.Data()[2:]
	if len(gatePayload) != 149 {
		t.Fatalf("PROTO_NC_BRIEFINFO_REGENMOB_CMD Gate payload=%dB, esperado=149B", len(gatePayload))
	}
	if gatePayload[14] != 1 || !strings.HasPrefix(string(gatePayload[15:27]), "Rou") {
		t.Fatalf("Gate flagstate=%d gate2where=%q, esperado flagstate=1 prefix=\"Rou\"", gatePayload[14], string(gatePayload[15:27]))
	}

	drop := &zone.ZoneDroppedItem{Handle: 300, ItemID: 500, X: 6888, Y: 7045, DropMobHandle: 200, Attr: 1}
	dropPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
	drop.WriteDropEntry(dropPW)
	if got := len(dropPW.Data()) - 2; got != 15 {
		t.Fatalf("PROTO_NC_BRIEFINFO_DROPEDITEM_CMD payload=%dB, esperado=15B", got)
	}

	mf := &zone.ZoneMagicFieldBrief{Handle: 301, Caster: 100, SkillID: 4001, X: 6888, Y: 7045, Radius: 200}
	mfPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD)
	mf.WriteMagicFieldSpread(mfPW)
	if got := len(mfPW.Data()) - 2; got != 16 {
		t.Fatalf("PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD payload=%dB, esperado=16B", got)
	}

	mv := &zone.ZoneMoverBrief{Handle: 302, ID: 5500, HP: 500, X: 6888, Y: 7045, Direct: 120, Grade: 1}
	mvPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOVER_CMD)
	mv.WriteRegenMover(mvPW)
	if got := len(mvPW.Data()) - 2; got != 139 {
		t.Fatalf("PROTO_NC_BRIEFINFO_REGENMOVER_CMD payload=%dB, esperado=139B", got)
	}

	door := &zone.ZoneDoorBrief{Handle: 303, MobID: 9001, X: 6888, Y: 7045, Direct: 0, DoorState: 0, BlockIndex: "DoorA"}
	doorPW := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BUILDDOOR_CMD)
	door.WriteBuildDoor(doorPW)
	if got := len(doorPW.Data()) - 2; got != 48 {
		t.Fatalf("PROTO_NC_BRIEFINFO_BUILDDOOR_CMD payload=%dB, esperado=48B", got)
	}

	// 3. Verificar tamaño binario empaquetado (binary.Size) de structs PDB Dept 18 en pkg/protocol
	skillSizes := []struct {
		name string
		got  int
		want int
	}{
		{"PROTO_NC_SKILL_EMPOWALLOC_REQ", binary.Size(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ_DATA{}), 6},
		{"PROTO_NC_SKILL_EMPOWALLOC_ACK", binary.Size(protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK_DATA{}), 4},
		{"PROTO_NC_SKILL_EMPOWPOINT_CMD", binary.Size(protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD_DATA{}), 1},
		{"PROTO_NC_SKILL_PRODUCTFIELD_REQ", binary.Size(protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ_DATA{}), 2},
		{"PROTO_NC_SKILL_PRODUCTFIELD_ACK", binary.Size(protocol.PROTO_NC_SKILL_PRODUCTFIELD_ACK_DATA{}), 4},
		{"PROTO_NC_SKILL_UNLEARN_REQ", binary.Size(protocol.PROTO_NC_SKILL_UNLEARN_REQ_DATA{}), 3},
		{"PROTO_NC_SKILL_UNLEARN_ACK", binary.Size(protocol.PROTO_NC_SKILL_UNLEARN_ACK_DATA{}), 4},
		{"PROTO_NC_SKILL_REVIVE_CMD", binary.Size(protocol.PROTO_NC_SKILL_REVIVE_CMD_DATA{}), 12},
		{"PROTO_NC_SKILL_SOMEONEREVIVE_CMD", binary.Size(protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD_DATA{}), 14},
	}
	for _, sc := range skillSizes {
		if sc.got != sc.want {
			t.Errorf("%s binary.Size = %d, esperado=%d", sc.name, sc.got, sc.want)
		}
	}
}

// TestFase18_BriefInfo_Combat_Skills_SoulStone_ClusterE2E ejecuta el flujo TCP
// real LoginService -> WorldService -> ZoneService -> MySQL validando al byte:
//  1) Lotes NC_BRIEFINFO_MOB_CMD (1 + N*149B, máx. 28 entradas/lote) y BRIEFINFO_DROPEDITEM_CMD (15B).
//  2) Dept 9 (NC_BAT): TARGETTING_REQ -> TARGETINFO_CMD (30B), ASSIST_REQ -> ASSIST_ACK (2B),
//     BASHSTART/BASHSTOP, SKILLBASH_OBJ_CAST_REQ -> SKILLBASH_HIT_OBJ_START_CMD (6B) +
//     SKILLBASH_HIT_DAMAGE_CMD (5+14*N B), SWING_DAMAGE_CMD (16B), REALLYKILL_CMD (4B),
//     EXPGAIN_CMD (6B) y ABSTATE_ERASE_REQ -> ABSTATE_ERASE_ACK (2B = 0x1C70 / 0x1C71).
//  3) Dept 18 (NC_SKILL): EMPOWALLOC_REQ (6B) -> EMPOWALLOC_ACK (4B, 0x0881 / 0x0889) +
//     EMPOWPOINT_CMD (1B), PRODUCTFIELD_REQ -> PRODUCTFIELD_ACK (4B, 0x0B01),
//     UNLEARN_REQ (3B: slot+skill) -> UNLEARN_ACK (4B, 0x1501),
//     REPLYREVIVE_CMD -> REVIVE_CMD (12B) + SOMEONEREVIVE_CMD (14B) + HPCHANGE (6B) + SPCHANGE (4B).
//  4) Dept 20 (NC_SOULSTONE): HP/SP_BUY_REQ -> HP/SP_BUY_ACK (2B) + BUYFAIL_ACK (2B = 0x0742),
//     HP/SP_USE_REQ -> SOMEONEUSE_CMD (2B) + USESUC_ACK (0B) + HPCHANGE (6B) / SPCHANGE (4B).
func TestFase18_BriefInfo_Combat_Skills_SoulStone_ClusterE2E(t *testing.T) {
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
		loginPort = 19511
		worldPort = 19515
		zonePort  = 19525
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
		Name:     "Fase18Cluster",
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

	testUserNo := 998918
	testUser := fmt.Sprintf("f18cl_%d", time.Now().Unix()%100000)
	testPass := "f18pass"
	sum := md5.Sum([]byte(testPass))
	testPassHash := hex.EncodeToString(sum[:])

	cleanupClusterUser := func(charID int) {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tLogout WHERE nUserNo = ?", testUserNo)
		var charIDs []int
		if charID > 0 {
			charIDs = append(charIDs, charID)
		} else if rows, qErr := db.Query("SELECT nCharNo FROM tCharacter WHERE nUserNo = ?", testUserNo); qErr == nil {
			for rows.Next() {
				var cid int
				if err := rows.Scan(&cid); err == nil {
					charIDs = append(charIDs, cid)
				}
			}
			_ = rows.Close()
		}
		for _, cid := range charIDs {
			_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", cid)
			_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", cid)
			_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", cid)
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}
	cleanupClusterUser(0)

	if _, err := db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
	`, testUserNo, testUser, testPassHash); err != nil {
		t.Fatalf("Error insertando tUser: %v", err)
	}

	var createdCharID int
	defer func() {
		cleanupClusterUser(createdCharID)
	}()

	// 1. Handshake y selección de mundo -> WorldService
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
	_, _ = auditReadServerPacket(lConn) // 0x0C67
	_, _ = auditReadServerPacket(lConn) // 0x0807

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
	wlReq.WriteString(testUser, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD falló: %v / %v", wlAck, err)
	}

	testCharName := fmt.Sprintf("F18Hero%d", time.Now().Unix()%10000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(testCharName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE falló: %v / %v", createAck, err)
	}

	createdChar, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || createdChar == nil {
		t.Fatalf("Personaje '%s' no encontrado: %v", testCharName, err)
	}
	createdCharID = createdChar.ID
	_ = charRepo.AddCharMoney(createdCharID, 50000)

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
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName).Scan(&dbRandomID); err != nil {
		t.Fatalf("Token tActiveZoneTransfer no encontrado: %v", err)
	}

	// 2. Entrar a ZoneService y capturar lotes NC_BRIEFINFO_MOB_CMD
	zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
	if err != nil {
		t.Fatalf("Dial ZoneService: %v", err)
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
			t.Fatalf("Error leyendo MAP_LOGIN: %v", err)
		}
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			break
		}
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))

	collectPackets := func(wait time.Duration) []*protocol.Packet {
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

	findPkt := func(pkts []*protocol.Packet, op uint16) *protocol.Packet {
		for _, p := range pkts {
			if p.OpCode == op {
				return p
			}
		}
		return nil
	}

	initPkts := collectPackets(350 * time.Millisecond)
	mobBatchCount := 0
	for _, p := range initPkts {
		if p.OpCode == protocol.PROTO_NC_BRIEFINFO_MOB_CMD {
			mobBatchCount++
			payload := p.Raw[2:]
			if len(payload) < 1 {
				t.Fatalf("PROTO_NC_BRIEFINFO_MOB_CMD vacío")
			}
			cnt := int(payload[0])
			if cnt == 0 || cnt > zone.BriefInfoMaxBatch {
				t.Fatalf("PROTO_NC_BRIEFINFO_MOB_CMD count=%d fuera de [1..%d]", cnt, zone.BriefInfoMaxBatch)
			}
			if len(payload) != 1+cnt*149 {
				t.Fatalf("PROTO_NC_BRIEFINFO_MOB_CMD payload=%dB, esperado 1+%d*149=%dB", len(payload), cnt, 1+cnt*149)
			}
		}
	}
	if mobBatchCount == 0 {
		t.Fatalf("Se esperaba al menos un lote PROTO_NC_BRIEFINFO_MOB_CMD tras MAP_LOGINCOMPLETE")
	}

	zp, ok := zoneSrv.Handler().GetPlayerByName(testCharName)
	if !ok || zp == nil || zp.CurrentMap == nil {
		t.Fatalf("Jugador '%s' no activo en mapa", testCharName)
	}
	zp.Char.TCharacterState.Money = 50000

	// Spawnear mob adyacente al jugador para auditar TARGETINFO, ASSIST, SKILLBASH, SWING_DAMAGE, REALLYKILL y EXPGAIN
	px, py := zp.GetPosition()
	mob := zone.NewZoneMonster(60990, 1102, "AuditSlime", zp.CurrentMap.MapID, int32(px+20), int32(py), 0, 2, 10, 0, 1, 2, 0, 1500, 30, 80, 35)
	zp.CurrentMap.AddMonster(mob)

	// 3. Dept 9: TARGETTING_REQ (0x2401) + ASSIST_REQ (0x243E)
	tgtReq := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETTING_REQ)
	tgtReq.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(tgtReq)
	tgtPkts := collectPackets(150 * time.Millisecond)
	tgtInfo := findPkt(tgtPkts, protocol.PROTO_NC_BAT_TARGETINFO_CMD)
	if tgtInfo == nil {
		t.Fatalf("Falta PROTO_NC_BAT_TARGETINFO_CMD tras TARGETTING_REQ")
	}
	if got := len(tgtInfo.Raw) - 2; got != 30 {
		t.Fatalf("PROTO_NC_BAT_TARGETINFO_CMD payload=%dB, esperado=30B", got)
	}
	tiPayload := tgtInfo.Raw[2:]
	if gotMaxHP := binary.LittleEndian.Uint32(tiPayload[7:11]); gotMaxHP != 10 {
		t.Fatalf("TARGETINFO MaxHP=%d, esperado=10", gotMaxHP)
	}

	assistReq := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ASSIST_REQ)
	assistReq.WriteUint16(zp.MapObjectID)
	_ = zClient.sendPacket(assistReq)
	assistPkts := collectPackets(120 * time.Millisecond)
	assistAck := findPkt(assistPkts, protocol.PROTO_NC_BAT_ASSIST_ACK)
	if assistAck == nil || len(assistAck.Raw)-2 != 2 {
		t.Fatalf("Falta PROTO_NC_BAT_ASSIST_ACK (2B) tras ASSIST_REQ")
	}

	// 4. Dept 9: SKILLBASH_OBJ_CAST_REQ (0x2440) -> SKILLBASH_HIT_OBJ_START_CMD (6B) + HIT_DAMAGE + REALLYKILL (4B) + EXPGAIN (6B)
	castReq := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castReq.WriteUint16(101)
	castReq.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(castReq)
	castPkts := collectPackets(250 * time.Millisecond)

	if p := findPkt(castPkts, protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD); p == nil || len(p.Raw)-2 != 6 {
		t.Fatalf("SKILLBASH_HIT_OBJ_START_CMD ausente o tamaño != 6B: %v", p)
	}
	if p := findPkt(castPkts, protocol.PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD); p == nil || len(p.Raw)-2 < 5+14 {
		t.Fatalf("SKILLBASH_HIT_DAMAGE_CMD ausente o tamaño < 19B: %v", p)
	}
	if p := findPkt(castPkts, protocol.PROTO_NC_BAT_REALLYKILL_CMD); p == nil || len(p.Raw)-2 != 4 {
		t.Fatalf("REALLYKILL_CMD ausente o tamaño != 4B tras muerte del mob: %v", p)
	}
	if p := findPkt(castPkts, protocol.PROTO_NC_BAT_EXPGAIN_CMD); p == nil || len(p.Raw)-2 != 6 {
		t.Fatalf("EXPGAIN_CMD ausente o tamaño != 6B tras matar mob: %v", p)
	}

	// BASHSTART (0x242B) sobre un segundo mob -> SWING_DAMAGE_CMD (16B) + BASHSTOP (0x2432)
	mob2 := zone.NewZoneMonster(60991, 1102, "AuditSlime2", zp.CurrentMap.MapID, int32(px+20), int32(py), 0, 2, 500, 0, 1, 2, 0, 1500, 30, 80, 10)
	zp.CurrentMap.AddMonster(mob2)
	bashStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_BASHSTART_CMD)
	bashStart.WriteUint16(mob2.Handle)
	_ = zClient.sendPacket(bashStart)
	bashPkts := collectPackets(200 * time.Millisecond)
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_BASHSTOP_CMD))
	if p := findPkt(bashPkts, protocol.PROTO_NC_BAT_SWING_DAMAGE_CMD); p == nil || len(p.Raw)-2 != 16 {
		t.Fatalf("SWING_DAMAGE_CMD ausente o tamaño != 16B: %v", p)
	}

	// SKILLBASH_CASTABORT_REQ (0x2444) -> SKILLBASH_CASTABORT_ACK (2B = 0x20EA cuando no castea)
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ))
	abortPkts := collectPackets(120 * time.Millisecond)
	if p := findPkt(abortPkts, protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("SKILLBASH_CASTABORT_ACK ausente o tamaño != 2B: %v", p)
	}

	// ABSTATE_ERASE_REQ (0x2454) -> ABSTATE_ERASE_ACK (2B = 0x1C70 cuando existe / 0x1C71 si ausente)
	eraseReq := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ)
	eraseReq.WriteUint16(999)
	_ = zClient.sendPacket(eraseReq)
	erasePkts := collectPackets(120 * time.Millisecond)
	if p := findPkt(erasePkts, protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("ABSTATE_ERASE_ACK ausente o tamaño != 2B: %v", p)
	}

	// 5. Dept 18 (NC_SKILL): EMPOWALLOC (6B -> 4B + 1B), PRODUCTFIELD (2B -> 4B), UNLEARN (3B -> 4B)
	empReq := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empReq.WriteUint16(101)
	empReq.WriteUint16(0x0002)
	empReq.WriteUint16(0)
	_ = zClient.sendPacket(empReq)
	empPkts := collectPackets(150 * time.Millisecond)
	empAck := findPkt(empPkts, protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK)
	if empAck == nil || len(empAck.Raw)-2 != 4 {
		t.Fatalf("EMPOWALLOC_ACK ausente o tamaño != 4B: %v", empAck)
	}
	if errCode := binary.LittleEndian.Uint16(empAck.Raw[4:6]); errCode != 0x0881 {
		t.Fatalf("EMPOWALLOC_ACK err=0x%04X, esperado=0x0881", errCode)
	}
	empPt := findPkt(empPkts, protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD)
	if empPt == nil || len(empPt.Raw)-2 != 1 {
		t.Fatalf("EMPOWPOINT_CMD ausente o tamaño != 1B: %v", empPt)
	}

	prodReq := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ)
	prodReq.WriteUint16(301)
	_ = zClient.sendPacket(prodReq)
	prodPkts := collectPackets(120 * time.Millisecond)
	prodAck := findPkt(prodPkts, protocol.PROTO_NC_SKILL_PRODUCTFIELD_ACK)
	if prodAck == nil || len(prodAck.Raw)-2 != 4 || binary.LittleEndian.Uint16(prodAck.Raw[4:6]) != 0x0B01 {
		t.Fatalf("PRODUCTFIELD_ACK inválido (esperado 4B err=0x0B01): %v", prodAck)
	}

	unlearnReq := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_REQ)
	unlearnReq.WriteByte(0)
	unlearnReq.WriteUint16(101)
	_ = zClient.sendPacket(unlearnReq)
	unlearnPkts := collectPackets(120 * time.Millisecond)
	unlearnAck := findPkt(unlearnPkts, protocol.PROTO_NC_SKILL_UNLEARN_ACK)
	if unlearnAck == nil || len(unlearnAck.Raw)-2 != 4 || binary.LittleEndian.Uint16(unlearnAck.Raw[4:6]) != 0x1501 {
		t.Fatalf("UNLEARN_ACK inválido (esperado 4B err=0x1501): %v", unlearnAck)
	}

	// REPLYREVIVE_CMD (0x482A) -> REVIVE_CMD (12B) + SOMEONEREVIVE_CMD (14B) + HPCHANGE (6B) + SPCHANGE (4B)
	zp.State = 3
	zp.Char.CurHP = 0
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD))
	revPkts := collectPackets(150 * time.Millisecond)
	if p := findPkt(revPkts, protocol.PROTO_NC_SKILL_REVIVE_CMD); p == nil || len(p.Raw)-2 != 12 {
		t.Fatalf("REVIVE_CMD ausente o tamaño != 12B: %v", p)
	}
	if p := findPkt(revPkts, protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD); p == nil || len(p.Raw)-2 != 14 {
		t.Fatalf("SOMEONEREVIVE_CMD ausente o tamaño != 14B: %v", p)
	}
	if p := findPkt(revPkts, protocol.PROTO_NC_BAT_HPCHANGE_CMD); p == nil || len(p.Raw)-2 != 6 {
		t.Fatalf("HPCHANGE_CMD ausente o tamaño != 6B tras REPLYREVIVE: %v", p)
	}

	// 6. Dept 20 (NC_SOULSTONE): HP_BUY / SP_BUY / BUYFAIL / HP_USE / SP_USE
	zp.StonesHP = 1
	zp.StonesSP = 1
	buyHP := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ)
	buyHP.WriteUint16(2)
	_ = zClient.sendPacket(buyHP)
	buyHPPkts := collectPackets(250 * time.Millisecond)
	if p := findPkt(buyHPPkts, protocol.PROTO_NC_SOULSTONE_HP_BUY_ACK); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("SOULSTONE_HP_BUY_ACK ausente o tamaño != 2B: %v", buyHPPkts)
	}

	buySP := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_BUY_REQ)
	buySP.WriteUint16(2)
	_ = zClient.sendPacket(buySP)
	buySPPkts := collectPackets(250 * time.Millisecond)
	if p := findPkt(buySPPkts, protocol.PROTO_NC_SOULSTONE_SP_BUY_ACK); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("SOULSTONE_SP_BUY_ACK ausente o tamaño != 2B: %v", buySPPkts)
	}

	buyOverflow := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ)
	buyOverflow.WriteUint16(500)
	_ = zClient.sendPacket(buyOverflow)
	failBuyPkts := collectPackets(120 * time.Millisecond)
	buyFail := findPkt(failBuyPkts, protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
	if buyFail == nil || len(buyFail.Raw)-2 != 2 || binary.LittleEndian.Uint16(buyFail.Raw[2:4]) != 0x0742 {
		t.Fatalf("SOULSTONE_BUYFAIL_ACK inválido (esperado 2B 0x0742): %v", buyFail)
	}

	zp.Char.CurHP = 10
	zp.Char.CurSP = 5
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USE_REQ))
	useHPPkts := collectPackets(150 * time.Millisecond)
	if p := findPkt(useHPPkts, protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK); p == nil || len(p.Raw)-2 != 0 {
		t.Fatalf("SOULSTONE_HP_USESUC_ACK ausente o tamaño != 0B: %v", p)
	}
	if p := findPkt(useHPPkts, protocol.PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("SOULSTONE_HP_SOMEONEUSE_CMD ausente o tamaño != 2B: %v", p)
	}
	if p := findPkt(useHPPkts, protocol.PROTO_NC_BAT_HPCHANGE_CMD); p == nil || len(p.Raw)-2 != 6 {
		t.Fatalf("HPCHANGE_CMD ausente o tamaño != 6B tras SOULSTONE_HP_USE: %v", p)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_USE_REQ))
	useSPPkts := collectPackets(150 * time.Millisecond)
	if p := findPkt(useSPPkts, protocol.PROTO_NC_SOULSTONE_SP_USESUC_ACK); p == nil || len(p.Raw)-2 != 0 {
		t.Fatalf("SOULSTONE_SP_USESUC_ACK ausente o tamaño != 0B: %v", p)
	}
	if p := findPkt(useSPPkts, protocol.PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD); p == nil || len(p.Raw)-2 != 2 {
		t.Fatalf("SOULSTONE_SP_SOMEONEUSE_CMD ausente o tamaño != 2B: %v", p)
	}
	if p := findPkt(useSPPkts, protocol.PROTO_NC_BAT_SPCHANGE_CMD); p == nil || len(p.Raw)-2 != 4 {
		t.Fatalf("SPCHANGE_CMD ausente o tamaño != 4B tras SOULSTONE_SP_USE: %v", p)
	}
}
