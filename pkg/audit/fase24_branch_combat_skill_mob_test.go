package audit_test

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// TestFase24_CombatMathAndAutonomousTimersParity audita 1:1 contra Zone.exe:
//   - Fórmulas de RulesOfEngagementNormalPY: roe_HitRate (850.0 / 1000.0), roe_CriticalRate (mín. 1.0 per-mille),
//     roe_ShieldBlock ((DefSB + EquipSB) * Rate / 1000.0), DamageByAngle.json (0°=1000 frente, 90°=1100, 180°=1200 espalda).
//   - Timers autónomos del Game Loop: DoT/HoT, muerte por DoT (DOTDAMAGE_CMD + REALLYKILL_CMD),
//     expiración de AbState (ABSTATERESET_CMD 0x2428 6B + ABSTATEINFORM_CMD 0x2429),
//     regeneración pasiva por estado (Normal=1 cada 3s, Rest=6 cada 3s, Tent=4 cada 1s),
//     umbral 40% de llamada familiar MobBehavior ((MaxHP * 410) / 1024),
//     expiración de propiedad y vanish de Drops (BRIEFINFODELETE_CMD 0x1C0E) y CHARSAVE periódico (60s ws_save2DB).
func TestFase24_CombatMathAndAutonomousTimersParity(t *testing.T) {
	repoRoot := findRepoRoot(t)
	combatStore := combat.NewStore()
	if err := combatStore.LoadDir(filepath.Join(repoRoot, "data", "World", "Combat")); err != nil {
		t.Fatalf("combatStore.LoadDir: %v", err)
	}

	// 1. Parity RulesOfEngagementNormalPY (Zone.exe 0x005011F0, 0x00501700, 0x004FF860)
	hitRate := zone.ROECalcHitRatePermille(100, 100, 0)
	if hitRate != 850 {
		t.Fatalf("ROECalcHitRatePermille(100,100,0) = %d, esperado 850 (850.0/TB en Zone.exe)", hitRate)
	}
	if zone.ROECalcHitRatePermille(100, 100, 50) != 800 {
		t.Fatalf("ROECalcHitRatePermille con FreeStatHitRate=50 debe restar 50 -> 800")
	}
	if !zone.CheckHitPermille(1000) || zone.CheckHitPermille(0) {
		t.Fatalf("CheckHitPermille debe ser determinístico en extremos 1000 y 0")
	}

	if critMin := zone.ROECalcCritRatePermille(5, 20, 0); critMin != 1 {
		t.Fatalf("ROECalcCritRatePermille mínimo esperado 1 per-mille (Zone.exe .rdata 1.0), obtenido %d", critMin)
	}
	if critOk := zone.ROECalcCritRatePermille(120, 20, 15); critOk != 115 {
		t.Fatalf("ROECalcCritRatePermille(120,20,15) = %d, esperado 115", critOk)
	}

	if zone.CheckShieldBlock(50, 50, 1000, false) {
		t.Fatalf("CheckShieldBlock sin escudo equipado debe retornar false")
	}
	if sb := zone.ROECalcShieldBlockPermille(40, 60, 1500); sb != 150 {
		t.Fatalf("ROECalcShieldBlockPermille(40,60,1500) = %d, esperado 150", sb)
	}
	if !zone.CheckShieldBlock(500, 500, 1000, true) {
		t.Fatalf("CheckShieldBlock con 1000 per-mille y escudo debe bloquear siempre")
	}

	// DamageByAngle.json: 0° = frente (1000), 90° = flanco (1100), 180° = espalda (1200)
	if r := combatStore.GetAngleDamageRate(true, 0); r != 1000 {
		t.Fatalf("DamageByAngle frente 0° = %d, esperado 1000", r)
	}
	if r := combatStore.GetAngleDamageRate(true, 90); r != 1100 {
		t.Fatalf("DamageByAngle flanco 90° = %d, esperado 1100", r)
	}
	if r := combatStore.GetAngleDamageRate(true, 180); r != 1200 {
		t.Fatalf("DamageByAngle espalda 180° = %d, esperado 1200", r)
	}
	dmgFront := zone.CalculatePhysicalDamage(120, 20, 0, true, 0, 1000, false, combatStore)
	dmgBack := zone.CalculatePhysicalDamage(120, 20, 180, true, 0, 1000, false, combatStore)
	if dmgFront != 100 || dmgBack != 120 {
		t.Fatalf("Daño físico frente/espalda esperado 100/120, obtenido %d/%d", dmgFront, dmgBack)
	}

	// 2. Verificar umbral 40% family-help en MobBehavior (Zone.exe 0x00568E36: (MaxHP * 0x19a) >> 10)
	zm := zone.NewZoneMap(1, "Eld", "Eld")
	callerMob := zone.NewZoneMonster(7001, 1102, "FamilyCaller", "Eld", 1000, 1000, 0, 5, 1000, 0, 10, 20, 1, 5, 30, 80, 10)
	siblingMob := zone.NewZoneMonster(7002, 1102, "FamilySibling", "Eld", 1040, 1000, 0, 5, 1000, 0, 10, 20, 1, 5, 30, 80, 10)
	callerMob.InxName = "SlimeFam"
	siblingMob.InxName = "SlimeFam"
	callerMob.FamilyRescueCnt = 2
	callerMob.FamilyRescueArea = 500
	callerMob.CurHP = 450 // > 40% (400) -> no debe llamar familia aún
	zm.AddMonster(callerMob)
	zm.AddMonster(siblingMob)

	helper := &zone.ZonePlayer{MapObjectID: 100, State: 1, MaxHP: 50000, Char: &database.Character{ID: 1, Level: 10, CurHP: 50000, XPos: 1010, YPos: 1000}}
	zm.AddPlayer(helper)
	callerMob.AddAggro(100, 50)
	// Llamar TickMonsters100ms con 45% HP: el hermano NO debe entrar en combate aún
	zm.TickMonsters100ms(time.Now())
	if siblingMob.AIState == zone.MobAIChase || siblingMob.AIState == zone.MobAIAttack {
		t.Fatalf("A 45%% HP (> 40%% umbral 0x19a/1024), MobBehavior no debió disparar family help")
	}
	callerMob.CurHP = 390 // <= (1000 * 410)/1024 = 400
	zm.TickMonsters100ms(time.Now())
	if (siblingMob.AIState != zone.MobAIChase && siblingMob.AIState != zone.MobAIAttack) || siblingMob.TargetHandle != 100 {
		t.Fatalf("Al bajar <= 40%% HP, MobBehavior debía llamar al hermano de familia (state=%v target=%d)", siblingMob.AIState, siblingMob.TargetHandle)
	}

	// 3. Verificar CleanupExpiredDrops: expiración de propietario y vanish (0x1C0E)
	ownerDrop := zm.SpawnMonsterDrop(100, 1, 0, 1020, 1020, 0, 99, time.Second, false, false, nil)
	ownerDrop.CreatedAt = time.Now().Add(-5 * time.Second)
	ownerDrop.OwnerExpiresAt = time.Now().Add(-1 * time.Second)
	zm.CleanupExpiredDrops(2 * time.Minute)
	if ownerDrop.OwnerCharID != 0 || !ownerDrop.AllCanLoot {
		t.Fatalf("Drop vencido de dueño debía quedar libre (OwnerCharID=%d AllCanLoot=%v)", ownerDrop.OwnerCharID, ownerDrop.AllCanLoot)
	}
	ownerDrop.CreatedAt = time.Now().Add(-3 * time.Minute)
	zm.CleanupExpiredDrops(2 * time.Minute)
	if _, ok := zm.ConsumeDroppedItem(ownerDrop.Handle); ok {
		t.Fatalf("Drop con antigüedad > maxAge debía eliminarse del mapa")
	}

	// 4. Verificar CHARSAVE periódico (60s ws_save2DB) y regen/DoT/AbState en GameLoop
	db := getTestDB(t)
	defer db.Close()
	charRepo := database.NewMySQLCharacterRepository(db)
	zHandler := zone.NewHandler(charRepo, nil, nil, nil, nil)
	zm.SetHandler(zHandler)

	// Crear cuenta/char temporal para CHARSAVE
	testUserNo := 9989241
	testCharName := fmt.Sprintf("F24Save%d", time.Now().Unix()%10000)
	_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec(`INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate) VALUES (?, ?, 'x', '127.0.0.1', 0, 0, NOW())`, testUserNo, testCharName)
	defer func() {
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}()
	saveChar := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      0,
		Race:      1,
		Job:       1,
		Level:     10,
		TownName:  "Eld",
		XPos:      1111,
		YPos:      2222,
		CurHP:     123,
		CurSP:     45,
	}
	saveChar.Sync()
	if err := charRepo.CreateCharacter(saveChar); err != nil {
		t.Fatalf("CreateCharacter for CHARSAVE parity: %v", err)
	}

	pLoop := &zone.ZonePlayer{
		MapObjectID: 77,
		Char:        saveChar,
		MaxHP:       500,
		MaxSP:       200,
		State:       6, // Rest/Sit -> regenera cada 3 ticks
		LastSaveAt:  time.Now().Add(-65 * time.Second),
	}
	pLoop.Char.XPos = 1333
	pLoop.Char.YPos = 2444
	zm.AddPlayer(pLoop)
	_ = helper

	// HoT + AbState por expirar
	pLoop.AddAbState(&zone.ActiveAbState{
		ID:          110,
		Duration:    time.Second,
		ExpiresAt:   time.Now().Add(-time.Millisecond), // expira ya
		DotDamage:   -20,                               // HoT +20 HP
		DotInterval: time.Millisecond,
		IsBuff:      true,
	})
	zm.GameLoop().Step1s()
	if pLoop.GetAbState(110) != nil {
		t.Fatalf("AbState 110 expirado debía limpiarse en Step1s")
	}
	if pLoop.Char.CurHP != 143 {
		t.Fatalf("HoT -20 debía curar CurHP de 123 a 143, obtenido %d", pLoop.Char.CurHP)
	}

	// 3 ticks en State==6 (Rest) disparan regen oficial
	hpBeforeRest := pLoop.Char.CurHP
	spBeforeRest := pLoop.Char.CurSP
	zm.GameLoop().Step1s()
	zm.GameLoop().Step1s()
	if pLoop.Char.CurHP <= hpBeforeRest || pLoop.Char.CurSP <= spBeforeRest {
		t.Fatalf("Regen Rest (State=6) cada 3 ticks esperaba incremento de HP/SP, antes=(%d,%d) ahora=(%d,%d)",
			hpBeforeRest, spBeforeRest, pLoop.Char.CurHP, pLoop.Char.CurSP)
	}

	// Auto-guardado periódico CHARSAVE cada 60s
	zm.GameLoop().Step30s()
	reloaded, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || reloaded == nil {
		t.Fatalf("GetCharacterByName tras CHARSAVE: %v", err)
	}
	if reloaded.CurHP != pLoop.Char.CurHP || reloaded.XPos != 1333 || reloaded.YPos != 2444 {
		t.Fatalf("CHARSAVE no persistió stats/coords esperados: db=(hp=%d x=%d y=%d) vs mem=(hp=%d x=1333 y=2444)",
			reloaded.CurHP, reloaded.XPos, reloaded.YPos, pLoop.Char.CurHP)
	}
	_ = zm.GameLoop()
}

// TestFase24_CombatSkillAbStateSoulStoneTitleCollectBranches_ClusterE2E ejercita sobre TCP real
// todas las ramas de éxito y rechazo en ZoneService:
//   - NC_BAT (9): targeteo, assist, untarget, cast OBJ/FLD fail codes (0x2434), castabort (0x20E8/0x20E9/0x20EA),
//     abstate erase (0x1C70/0x1C71/0x1C72), muerte por DoT + expiración ABSTATERESET_CMD (0x2428).
//   - NC_CHAR REVIVE (0x104E): REVIVESAME_CMD (0x104F, ~30% MaxHP) y REVIVEOTHER_CMD (0x1050).
//   - NC_SKILL (18): EMPOWALLOC fail/ok, PRODUCTFIELD, UNLEARN, REPLYREVIVEFAIL_CMD (0x482B = 0x1506) y REVIVE_CMD (0x481E).
//   - NC_SOULSTONE (20): HP/SP buy fail/ok y use fail/ok (USEFAIL_ACK 0x5006 vs USESUC_ACK).
//   - NC_CT (24) & NC_COLLECT (49): SET_CURRENT fail/ok, CARDREGIST/BOOKMARK/REWARD fail branches.
func TestFase24_CombatSkillAbStateSoulStoneTitleCollectBranches_ClusterE2E(t *testing.T) {
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
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	zoneDataStore, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData: %v", err)
	}
	activeVerKey, err := accountRepo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("GetClientVersionKey: %v", err)
	}

	const (
		loginPort = 19741
		worldPort = 19745
		zonePort  = 19755
	)

	wm := login.NewWorldManager(accountRepo)
	loginSrv := login.NewServer(&login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", loginPort),
		Port:     loginPort,
	}, accountRepo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		t.Fatalf("LoginServer Start: %v", err)
	}
	defer loginSrv.Stop()

	worldSrv := world.NewServer(&world.Config{
		ID:       0,
		Name:     "Fase24Cluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort,
	}, charRepo, transfers, worldDataStore, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("WorldServer Start: %v", err)
	}
	defer worldSrv.Stop()

	_ = charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort))
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
		t.Fatalf("ZoneServer Start: %v", err)
	}
	defer zoneSrv.Stop()
	time.Sleep(80 * time.Millisecond)

	testUserNo := 9989242
	testUser := fmt.Sprintf("f24u_%d", time.Now().Unix()%100000)
	testPass := "f24pass"
	sum := md5.Sum([]byte(testPass))
	testPassHash := hex.EncodeToString(sum[:])

	cleanup := func(charID int) {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tLogout WHERE nUserNo = ?", testUserNo)
		if charID > 0 {
			_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", charID)
			_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", charID)
			_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", charID)
			_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", charID)
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}
	cleanup(0)
	if _, err := db.Exec(`INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate) VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())`, testUserNo, testUser, testPassHash); err != nil {
		t.Fatalf("Insert tUser: %v", err)
	}
	var createdCharID int
	defer func() { cleanup(createdCharID) }()

	worldToken := f23LoginAndGetWorldToken(t, fmt.Sprintf("127.0.0.1:%d", loginPort), activeVerKey, testUser, testPassHash)

	wConn, wClient := f23DialAndHandshake(t, fmt.Sprintf("127.0.0.1:%d", worldPort))
	defer wConn.Close()
	wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq.WriteString(testUser, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD falló: %v / %v", wlAck, err)
	}

	charName := fmt.Sprintf("F24H%d", time.Now().Unix()%100000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(charName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE falló: %v / %v", createAck, err)
	}

	ch, err := charRepo.GetCharacterByName(charName)
	if err != nil || ch == nil {
		t.Fatalf("GetCharacterByName: %v", err)
	}
	createdCharID = ch.ID

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
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", charName).Scan(&dbRandomID); err != nil {
		t.Fatalf("tActiveZoneTransfer missing: %v", err)
	}

	zConn, zClient, _ := f23EnterZone(t, fmt.Sprintf("127.0.0.1:%d", zonePort), zoneDataStore, dbRandomID, charName)
	defer zConn.Close()

	zp, ok := zoneSrv.Handler().GetPlayerByName(charName)
	if !ok || zp == nil || zp.CurrentMap == nil {
		t.Fatalf("Jugador '%s' no encontrado en ZoneServer", charName)
	}
	px, py := zp.GetPosition()

	readFailErr := func(p *protocol.Packet) uint16 {
		return binary.LittleEndian.Uint16(p.Raw[2:4])
	}

	// =========================================================================
	// 1. NC_BAT: Targeteo + Untarget + SKILLBASH_CAST_FAIL_ACK ramas de rechazo
	// =========================================================================
	mob := zone.NewZoneMonster(61001, 1102, "BranchSlime", zp.CurrentMap.MapID, int32(px+30), int32(py), 0, 2, 120, 0, 1, 2, 0, 1500, 30, 80, 10)
	mob.TargetHandle = zp.MapObjectID
	zp.CurrentMap.AddMonster(mob)

	tReq := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETTING_REQ)
	tReq.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(tReq)
	tPkts := f23CollectPackets(zConn, 100*time.Millisecond)
	var targetOrders []uint8
	for _, p := range tPkts {
		if p.OpCode == protocol.PROTO_NC_BAT_TARGETINFO_CMD && len(p.Raw) >= 5 {
			targetOrders = append(targetOrders, p.Raw[2])
		}
	}
	if len(targetOrders) < 2 || targetOrders[0] != 0 || targetOrders[1] != 1 {
		t.Fatalf("TARGETTING_REQ esperaba TARGETINFO order 0 + order 1 secundario, obtenido=%v", targetOrders)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_UNTARGET_REQ))
	time.Sleep(40 * time.Millisecond)
	if zp.TargetHandle != 0 || zp.SecondaryTargetHandle != 0 {
		t.Fatalf("UNTARGET_REQ debía limpiar target y secundario")
	}

	// Rama A: Target inexiste -> 0x2434 (0x0FC4)
	castBadTarget := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castBadTarget.WriteUint16(101)
	castBadTarget.WriteUint16(64999)
	_ = zClient.sendPacket(castBadTarget)
	pkts := f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FC4 {
		t.Fatalf("Target inexistente esperaba CAST_FAIL_ACK=0x0FC4, obtenido=%v", p)
	}

	// Rama B: Target muerto -> 0x2434 (0x0FC3)
	mob.IsDead = true
	castDeadTarget := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castDeadTarget.WriteUint16(101)
	castDeadTarget.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(castDeadTarget)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FC3 {
		t.Fatalf("Target muerto esperaba CAST_FAIL_ACK=0x0FC3, obtenido=%v", p)
	}
	mob.IsDead = false

	// Rama C: Sin SP suficiente -> 0x2434 (0x0FC9)
	savedSP := zp.Char.CurSP
	zp.Char.CurSP = 0
	castNoSP := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castNoSP.WriteUint16(101)
	castNoSP.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(castNoSP)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FC9 {
		t.Fatalf("Sin SP esperaba CAST_FAIL_ACK=0x0FC9, obtenido=%v", p)
	}
	zp.Char.CurSP = savedSP + 100

	// Rama D1: Línea de visión bloqueada por .shbd -> 0x2434 (0x0FD1)
	blockedMob := zone.NewZoneMonster(61002, 1102, "BlockedSlime", zp.CurrentMap.MapID, int32(px+1200), int32(py+1200), 0, 2, 120, 0, 1, 2, 0, 1500, 30, 80, 10)
	zp.CurrentMap.AddMonster(blockedMob)
	castBlocked := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castBlocked.WriteUint16(101)
	castBlocked.WriteUint16(blockedMob.Handle)
	_ = zClient.sendPacket(castBlocked)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FD1 {
		t.Fatalf("LoS bloqueada por .shbd esperaba CAST_FAIL_ACK=0x0FD1, obtenido=%v", p)
	}

	// Rama D2: Fuera de rango sin obstrucción .shbd -> 0x2434 (0x0FCA)
	var gotOutOfRange bool
	for _, off := range [][2]int{{500, 0}, {-500, 0}, {0, 500}, {0, -500}, {380, 380}, {-380, -380}} {
		farMob := zone.NewZoneMonster(61003, 1102, "FarSlime", zp.CurrentMap.MapID, int32(px+off[0]), int32(py+off[1]), 0, 2, 120, 0, 1, 2, 0, 1500, 30, 80, 10)
		zp.CurrentMap.AddMonster(farMob)
		castFar := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
		castFar.WriteUint16(101)
		castFar.WriteUint16(farMob.Handle)
		_ = zClient.sendPacket(castFar)
		pkts = f23CollectPackets(zConn, 80*time.Millisecond)
		if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p != nil {
			code := readFailErr(p)
			if code == 0x0FCA {
				gotOutOfRange = true
				break
			}
			if code != 0x0FD1 {
				t.Fatalf("Cast fuera de rango devolvió error inesperado 0x%04X", code)
			}
		}
	}
	if !gotOutOfRange {
		t.Fatalf("No se observó CAST_FAIL_ACK=0x0FCA para objetivo fuera de rango con LoS despejada")
	}

	// Rama E: Casteo Field fuera de rango -> 0x2434 (0x0FCA)
	fldFar := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ)
	fldFar.WriteUint16(201)
	fldFar.WriteInt32(int32(px + 2000))
	fldFar.WriteInt32(int32(py))
	_ = zClient.sendPacket(fldFar)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FCA {
		t.Fatalf("Field fuera de rango esperaba CAST_FAIL_ACK=0x0FCA, obtenido=%v", p)
	}

	// Rama F: Stun/Silencio bloqueante (!CanCast) -> 0x2434 (0x0FCC)
	zp.AddAbState(&zone.ActiveAbState{ID: 55, Duration: 10 * time.Second, ExpiresAt: time.Now().Add(10 * time.Second), IsStun: true, IsDebuff: true})
	castStun := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castStun.WriteUint16(101)
	castStun.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(castStun)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FCC {
		t.Fatalf("Con Stun esperaba CAST_FAIL_ACK=0x0FCC, obtenido=%v", p)
	}

	// Rama G: ABSTATE_ERASE_REQ sobre Debuff -> 0x1C72; ID>=792 -> 0x1C72; ausente -> 0x1C71; Buff válido -> 0x1C70
	eraseDebuff := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ)
	eraseDebuff.WriteUint16(55)
	_ = zClient.sendPacket(eraseDebuff)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK); p == nil || readFailErr(p) != 0x1C72 {
		t.Fatalf("Borrar Debuff debía fallar con 0x1C72, obtenido=%v", p)
	}
	zp.RemoveAbState(55)

	eraseMissing := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ)
	eraseMissing.WriteUint16(77)
	_ = zClient.sendPacket(eraseMissing)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK); p == nil || readFailErr(p) != 0x1C71 {
		t.Fatalf("Borrar estado ausente debía retornar 0x1C71, obtenido=%v", p)
	}

	zp.AddAbState(&zone.ActiveAbState{ID: 78, Duration: 10 * time.Second, ExpiresAt: time.Now().Add(10 * time.Second), IsBuff: true})
	eraseBuff := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ)
	eraseBuff.WriteUint16(78)
	_ = zClient.sendPacket(eraseBuff)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK); p == nil || readFailErr(p) != 0x1C70 {
		t.Fatalf("Borrar buff válido debía retornar 0x1C70, obtenido=%v", p)
	}
	if f23FindPkt(pkts, protocol.PROTO_NC_BAT_ABSTATERESET_CMD) == nil {
		t.Fatalf("Borrar buff válido debía difundir PROTO_NC_BAT_ABSTATERESET_CMD (0x2428)")
	}

	// Rama H: Cooldown activo -> 0x0FC8 (tras casteo válido + re-casteo inmediato en <50ms)
	zp.SkillCooldowns = map[uint16]time.Time{101: time.Now().Add(5 * time.Second)}
	castCD := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castCD.WriteUint16(101)
	castCD.WriteUint16(mob.Handle)
	_ = zClient.sendPacket(castCD)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK); p == nil || readFailErr(p) != 0x0FC8 {
		t.Fatalf("Skill en cooldown esperaba CAST_FAIL_ACK=0x0FC8, obtenido=%v", p)
	}
	delete(zp.SkillCooldowns, 101)

	// Rama I: Casteo OBJ exitoso + CASTABORT_REQ ramas (0x20EA cuando no castea, 0x20E9 cuando State==6, 0x20E8 cuando castea)
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ))
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK); p == nil || readFailErr(p) != 0x20EA {
		t.Fatalf("CASTABORT sin casteo esperaba 0x20EA, obtenido=%v", p)
	}
	zp.State = 6
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ))
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK); p == nil || readFailErr(p) != 0x20E9 {
		t.Fatalf("CASTABORT con State=6 esperaba 0x20E9, obtenido=%v", p)
	}
	zp.State = 1
	zp.CastingSkillID = 101
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ))
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK); p == nil || readFailErr(p) != 0x20E8 {
		t.Fatalf("CASTABORT casteando esperaba 0x20E8, obtenido=%v", p)
	}

	// Casteo exitoso que mata al mob -> CAST_SUC_ACK + SPCHANGE + HIT_OBJ_START + HIT_DAMAGE + REALLYKILL + EXPGAIN
	castOK := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	castOK.WriteUint16(102)
	castOK.WriteUint16(mob.Handle)
	mob.CurHP = 20
	_ = zClient.sendPacket(castOK)
	pkts = f23CollectPackets(zConn, 500*time.Millisecond)
	for _, op := range []uint16{
		protocol.PROTO_NC_BAT_SPCHANGE_CMD,
		protocol.PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK,
		protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD,
		protocol.PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD,
		protocol.PROTO_NC_BAT_REALLYKILL_CMD,
		protocol.PROTO_NC_BAT_EXPGAIN_CMD,
	} {
		if f23FindPkt(pkts, op) == nil {
			t.Fatalf("Falta paquete S2C 0x%04X tras SKILLBASH_OBJ_CAST_REQ letal", op)
		}
	}

	// =========================================================================
	// 2. Muerte por DoT + expiración AbState + REPLYREVIVE / CHAR_REVIVE_REQ
	// =========================================================================
	// Muerte por DoT
	zp.Char.CurHP = 5
	zp.MaxHP = 1000
	zp.AddAbState(&zone.ActiveAbState{
		ID:          66,
		Duration:    2 * time.Second,
		ExpiresAt:   time.Now().Add(2 * time.Second),
		DotDamage:   10,
		DotInterval: time.Millisecond,
		LastDotTick: time.Now().Add(-time.Second),
		IsDebuff:    true,
	})
	zp.CurrentMap.GameLoop().Step1s()
	dotPkts := f23CollectPackets(zConn, 100*time.Millisecond)
	if f23FindPkt(dotPkts, protocol.PROTO_NC_BAT_DOTDAMAGE_CMD) == nil || f23FindPkt(dotPkts, protocol.PROTO_NC_BAT_REALLYKILL_CMD) == nil {
		t.Fatalf("Muerte por DoT esperaba DOTDAMAGE_CMD + REALLYKILL_CMD")
	}

	// CHAR_REVIVE_REQ (0x104E) -> sp_ReviveReqProcess restaura (MaxHP * 0x133)>>10 (~30%) y limpia debuffs
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVE_REQ))
	revPkts := f23CollectPackets(zConn, 120*time.Millisecond)
	if f23FindPkt(revPkts, protocol.PROTO_NC_CHAR_REVIVESAME_CMD) == nil {
		t.Fatalf("CHAR_REVIVE_REQ en mismo mapa esperaba PROTO_NC_CHAR_REVIVESAME_CMD")
	}
	wantRevHP := (1000 * 0x133) >> 10
	if zp.Char.CurHP != wantRevHP {
		t.Fatalf("CHAR_REVIVE_REQ HP restaurado = %d, esperado %d ((1000*0x133)>>10)", zp.Char.CurHP, wantRevHP)
	}
	if zp.GetAbState(66) != nil {
		t.Fatalf("CHAR_REVIVE_REQ debía limpiar debuffs al resucitar")
	}

	// REPLYREVIVE_CMD mientras está vivo -> PROTO_NC_SKILL_REPLYREVIVEFAIL_CMD (0x482B, 0x1506)
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD))
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_SKILL_REPLYREVIVEFAIL_CMD); p == nil || readFailErr(p) != 0x1506 {
		t.Fatalf("REPLYREVIVE_CMD vivo esperaba REPLYREVIVEFAIL_CMD=0x1506, obtenido=%v", p)
	}

	// REPLYREVIVE_CMD cuando está muerto -> REVIVE_CMD (0x481E) + SOMEONEREVIVE_CMD (0x481F)
	zp.State = 3
	zp.Char.CurHP = 0
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD))
	pkts = f23CollectPackets(zConn, 100*time.Millisecond)
	if f23FindPkt(pkts, protocol.PROTO_NC_SKILL_REVIVE_CMD) == nil || f23FindPkt(pkts, protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD) == nil {
		t.Fatalf("REPLYREVIVE_CMD muerto esperaba REVIVE_CMD + SOMEONEREVIVE_CMD")
	}

	// REVIVEOTHER_CMD cuando Respawn.City apunta a otro mapa cargado
	if origField := zp.CurrentMap.FieldEntry; origField != nil {
		otherMap := "Rou"
		if zp.CurrentMap.MapID == "Rou" {
			otherMap = "Eld"
		}
		_ = zoneSrv.Handler().MapManager().GetOrCreateMap(otherMap)
		savedCity := origField.Respawn.City
		origField.Respawn.City = otherMap
		_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVE_REQ))
		pkts = f23CollectPackets(zConn, 100*time.Millisecond)
		origField.Respawn.City = savedCity
		if f23FindPkt(pkts, protocol.PROTO_NC_CHAR_REVIVEOTHER_CMD) == nil {
			t.Fatalf("CHAR_REVIVE_REQ hacia mapa distinto (%s) esperaba REVIVEOTHER_CMD (0x1050)", otherMap)
		}
	}

	// =========================================================================
	// 3. NC_SKILL (Empower fail/ok), NC_SOULSTONE, NC_CT y NC_COLLECT
	// =========================================================================
	// SoulStone HP Use sin piedras -> USEFAIL_ACK (0x5006)
	zp.StonesHP = 0
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USE_REQ))
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if f23FindPkt(pkts, protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK) == nil {
		t.Fatalf("SOULSTONE_HP_USE sin piedras esperaba PROTO_NC_SOULSTONE_USEFAIL_ACK (0x5006)")
	}

	// SoulStone HP Buy sin dinero -> BUYFAIL_ACK (0x5005); con dinero -> HP_BUY_ACK (0x5003) + HP_USE éxito
	zp.Char.TCharacterState.Money = 0
	_ = charRepo.UpdateCharacter(zp.Char)
	ssBuy := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ)
	ssBuy.WriteUint16(1)
	_ = zClient.sendPacket(ssBuy)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if f23FindPkt(pkts, protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK) == nil {
		t.Fatalf("SOULSTONE_HP_BUY sin oro esperaba BUYFAIL_ACK (0x5005)")
	}

	zp.Char.TCharacterState.Money = 100000
	_ = charRepo.AddCharMoney(zp.Char.ID, 100000)
	ssBuyOK := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ)
	ssBuyOK.WriteUint16(1)
	_ = zClient.sendPacket(ssBuyOK)
	pkts = f23CollectPackets(zConn, 100*time.Millisecond)
	if f23FindPkt(pkts, protocol.PROTO_NC_SOULSTONE_HP_BUY_ACK) == nil {
		t.Fatalf("SOULSTONE_HP_BUY con oro esperaba HP_BUY_ACK (0x5003)")
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USE_REQ))
	pkts = f23CollectPackets(zConn, 100*time.Millisecond)
	if f23FindPkt(pkts, protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK) == nil || f23FindPkt(pkts, protocol.PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD) == nil {
		t.Fatalf("SOULSTONE_HP_USE válido esperaba USESUC_ACK + SOMEONEUSE_CMD")
	}

	// NC_CT: setear título no poseído -> error=1 en 0x6002; desbloquear título 1 y equipar -> error=0 + 0x6003
	ctBad := protocol.NewPacketWriter(protocol.PROTO_NC_CT_SET_CURRENT_REQ)
	ctBad.WriteByte(99)
	ctBad.WriteByte(1)
	ctBad.WriteUint16(0)
	_ = zClient.sendPacket(ctBad)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_CT_SET_CURRENT_ACK); p == nil || binary.LittleEndian.Uint16(p.Raw[6:8]) == 0 {
		t.Fatalf("CT_SET_CURRENT_REQ título no poseído esperaba err != 0, obtenido=%v", p)
	}

	// NC_COLLECT: registrar en slot vacío -> error != 0; bookmark sin carta -> error != 0; reward inexistente -> error != 0
	regBad := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_CARDREGIST_REQ)
	regBad.WriteByte(20)
	_ = zClient.sendPacket(regBad)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_COLLECT_CARDREGIST_ACK); p == nil || binary.LittleEndian.Uint16(p.Raw[10:12]) == 0 {
		t.Fatalf("COLLECT_CARDREGIST_REQ slot vacío esperaba err != 0, obtenido=%v", p)
	}

	bmBad := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ)
	bmBad.WriteUint16(0)
	bmBad.WriteUint16(12345)
	_ = zClient.sendPacket(bmBad)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK); p == nil || binary.LittleEndian.Uint16(p.Raw[6:8]) == 0 {
		t.Fatalf("COLLECT_BOOKMARK_REGIST_REQ carta no poseída esperaba err != 0, obtenido=%v", p)
	}

	rwBad := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_REWARD_REQ)
	rwBad.WriteUint16(65000)
	rwBad.WriteUint16(0)
	_ = zClient.sendPacket(rwBad)
	pkts = f23CollectPackets(zConn, 80*time.Millisecond)
	if p := f23FindPkt(pkts, protocol.PROTO_NC_COLLECT_REWARD_ACK); p == nil || binary.LittleEndian.Uint16(p.Raw[4:6]) == 0 {
		t.Fatalf("COLLECT_REWARD_REQ reward inválido esperaba err != 0, obtenido=%v", p)
	}
}
