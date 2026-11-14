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
	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// TestFase17_SHBD_CoordinateConversionAndMovement valida la conversión exacta de
// coordenadas de mundo a celda/bit .shbd de Zone.exe ((world * 8) / 50) y elimina
// los falsos positivos PROTO_NC_ACT_MOVEFAIL_CMD en mapas oficiales (RouCos02, RouN).
func TestFase17_SHBD_CoordinateConversionAndMovement(t *testing.T) {
	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	// Fórmula oficial Zone.exe: tile = (world * 8) / 50
	tx, ty := blockinfo.WorldToTile(9647, 9002)
	wantTX, wantTY := uint32((9647*8)/50), uint32((9002*8)/50)
	if tx != wantTX || ty != wantTY {
		t.Fatalf("WorldToTile(9647,9002) = (%d,%d), esperado Zone.exe (%d,%d)", tx, ty, wantTX, wantTY)
	}

	rouCos02, err := blockinfo.LoadSHBDFile(filepath.Join(dataDir, "BlockInfo", "RouCos02.shbd"))
	if err != nil {
		t.Fatalf("Error cargando RouCos02.shbd: %v", err)
	}
	if rouCos02.IsWorldMoveBlocked(9647, 9002) {
		t.Fatalf("RouCos02 (9647,9002) debería ser caminable")
	}
	if rouCos02.IsWorldMoveBlocked(9588, 9200) {
		t.Fatalf("RouCos02 (9588,9200) debería ser caminable")
	}
	if rouCos02.IsMovePathBlocked(9647, 9002, 9588, 9200, 1) {
		t.Fatalf("RouCos02 trayectoria (9647,9002 -> 9588,9200) no debe bloquearse")
	}

	rouN, err := blockinfo.LoadSHBDFile(filepath.Join(dataDir, "BlockInfo", "RouN.shbd"))
	if err != nil {
		t.Fatalf("Error cargando RouN.shbd: %v", err)
	}
	if rouN.IsWorldMoveBlocked(6888, 7045) {
		t.Fatalf("RouN (6888,7045) debería ser caminable")
	}
	if !rouN.IsWorldMoveBlocked(1e9, 1e9) {
		t.Fatalf("Coordenadas fuera de límites deben reportarse bloqueadas")
	}
}

// TestFase17_InGameChunk_ChatGM_And_Logout_ClusterE2E audita el flujo binario real
// en clúster TCP (LoginService + WorldService + ZoneService + MySQL):
//  1) In-Game Chunk completo (0x1038..0x104A, 0x10CE, 0x10D7 -> 0x1802 PROTO_NC_MAP_LOGIN_ACK 242B).
//  2) Movimiento sin MOVEFAIL falso en terreno caminable (PROTO_NC_ACT_MOVEWALK/MOVERUN/STOP).
//  3) Chat normal con/sin itemLinkData (103B) -> PROTO_NC_ACT_SOMEONECHAT_CMD (7B header),
//     Shout -> PROTO_NC_ACT_SOMEONESHOUT_CMD (23B header), y comandos GM AdminLvSet.shn (&adminlevel).
//  4) Logout completo: 0x1071 LOGOUTREADY -> 0x1072 LOGOUTCANCEL -> 0x0C18 LogoutType=1 (retorno
//     a Selección de Personaje en WorldService preservando tActiveAccountSession) -> 0x0C33/0x0C34 OTP
//     -> 0x0C37 LoginService OTP -> 0x0C18 LogoutType=0 (cierre total de tActiveAccountSession).
func TestFase17_InGameChunk_ChatGM_And_Logout_ClusterE2E(t *testing.T) {
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
		loginPort = 19510
		worldPort = 19514
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
		ID:       0,
		Name:     "Fase17Cluster",
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

	testUserNo := 998917
	testUser := fmt.Sprintf("f17cl_%d", time.Now().Unix()%100000)
	testPass := "f17pass"
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

	// 1. LoginService -> WorldSelect
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
	verAck, err := auditReadServerPacket(lConn)
	if err != nil || verAck.OpCode != protocol.PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK {
		t.Fatalf("VersionCheck falló: %v / %v", verAck, err)
	}
	_, _ = auditReadServerPacket(lConn) // 0x0807 seed ack

	xtPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_XTRAP_REQ)
	xtPkt.WriteByte(1)
	_ = lClient.sendPacket(xtPkt)
	_, _ = auditReadServerPacket(lConn)

	lReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	lReq.WriteString(testUser, 260)
	lReq.WriteString(testPassHash, 36)
	lReq.Fill(20, 0)
	_ = lClient.sendPacket(lReq)
	lAck, err := auditReadServerPacket(lConn)
	if err != nil || lAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
		t.Fatalf("US_LOGIN falló: %v / %v", lAck, err)
	}

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

	// 2. WorldService (manteniendo wConn abierto en segundo plano mientras el jugador entra a Zone)
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
	wlAck, err := auditReadServerPacket(wConn)
	if err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD falló: %v / %v", wlAck, err)
	}

	testCharName := fmt.Sprintf("F17Hero%d", time.Now().Unix()%10000)
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
		t.Fatalf("AVATAR_CREATE falló: %v / %v", createAck, err)
	}

	createdChar, err := charRepo.GetCharacterByName(testCharName)
	if err != nil || createdChar == nil {
		t.Fatalf("Personaje '%s' no encontrado: %v", testCharName, err)
	}
	createdCharID = createdChar.ID

	// Elevar AdminLevel conforme AdminLvSet.shn para verificar '&adminlevel'
	reqAdmLv, hasAdmLv := zoneDataStore.AdminLvCatalog.RequiredLevel("adminlevel")
	if !hasAdmLv {
		reqAdmLv = 100
	}
	if _, err := charRepo.SetAdminLevel(createdCharID, reqAdmLv); err != nil {
		t.Fatalf("SetAdminLevel falló: %v", err)
	}

	cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	cSelReq.WriteByte(0)
	_ = wClient.sendPacket(cSelReq)
	cSelAck, err := auditReadServerPacket(wConn)
	if err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CHAR_LOGIN falló: %v / %v", cSelAck, err)
	}
	_, _ = auditReadServerPacket(wConn) // Option packets
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)

	var dbRandomID uint16
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName).Scan(&dbRandomID); err != nil {
		t.Fatalf("Token no emitido en tActiveZoneTransfer: %v", err)
	}

	// 3. ZoneService: In-Game Chunk completo
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
	mapReq.WriteString(testCharName, 20)
	sums := zoneDataStore.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		mapReq.WriteBytes(sums[i][:])
	}
	_ = zClient.sendPacket(mapReq)

	chunkSeen := make(map[uint16]int)
	var mapAckPkt *protocol.Packet
	var baseCmdPkt *protocol.Packet
	for {
		pkt, err := auditReadServerPacket(zConn)
		if err != nil {
			t.Fatalf("Error leyendo In-Game Chunk: %v", err)
		}
		chunkSeen[pkt.OpCode]++
		if pkt.OpCode == protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD {
			baseCmdPkt = pkt
		}
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			mapAckPkt = pkt
			break
		}
	}

	requiredChunk := []uint16{
		protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_SKILL_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD,
		protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD,
		protocol.PROTO_NC_MAP_LOGIN_ACK,
	}
	for _, op := range requiredChunk {
		if chunkSeen[op] == 0 {
			t.Fatalf("Falta paquete obligatorio In-Game Chunk 0x%04X (%s)", op, protocol.OpCodeName(op))
		}
	}
	if baseCmdPkt == nil || len(baseCmdPkt.Raw)-2 != 105 {
		t.Fatalf("PROTO_NC_CHAR_CLIENT_BASE_CMD payload esperado=105B, obtenido=%d", len(baseCmdPkt.Raw)-2)
	}
	if mapAckPkt == nil || len(mapAckPkt.Raw)-2 != 242 {
		t.Fatalf("PROTO_NC_MAP_LOGIN_ACK payload esperado=242B, obtenido=%d", len(mapAckPkt.Raw)-2)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	hbAck := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HEARTBEAT_ACK)
	_ = zClient.sendPacket(hbAck)
	time.Sleep(40 * time.Millisecond)

	zp, ok := zoneSrv.Handler().GetPlayerByName(testCharName)
	if !ok || zp == nil || zp.CurrentMap == nil {
		t.Fatalf("Jugador '%s' no activo en ZoneService tras LOGINCOMPLETE", testCharName)
	}
	if zp.MaxHP == 0 || zp.Char.CurHP == 0 || zp.MaxSP == 0 || zp.WalkSpeed == 0 || zp.RunSpeed == 0 {
		t.Fatalf("Stats HP/SP/Speed inválidos en ZonePlayer: HP=%d/%d SP=%d/%d Walk=%d Run=%d",
			zp.Char.CurHP, zp.MaxHP, zp.Char.CurSP, zp.MaxSP, zp.WalkSpeed, zp.RunSpeed)
	}

	// 4. Movimiento caminable sin PROTO_NC_ACT_MOVEFAIL_CMD falso
	walkPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVEWALK_CMD)
	walkPkt.WriteUint32(uint32(zp.Char.XPos))
	walkPkt.WriteUint32(uint32(zp.Char.YPos))
	walkPkt.WriteUint32(uint32(zp.Char.XPos + 20))
	walkPkt.WriteUint32(uint32(zp.Char.YPos + 20))
	_ = zClient.sendPacket(walkPkt)

	runPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVERUN_CMD)
	runPkt.WriteUint32(uint32(zp.Char.XPos + 20))
	runPkt.WriteUint32(uint32(zp.Char.YPos + 20))
	runPkt.WriteUint32(uint32(zp.Char.XPos + 50))
	runPkt.WriteUint32(uint32(zp.Char.YPos + 50))
	_ = zClient.sendPacket(runPkt)

	stopPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_STOP_REQ)
	stopPkt.WriteUint32(uint32(zp.Char.XPos + 50))
	stopPkt.WriteUint32(uint32(zp.Char.YPos + 50))
	_ = zClient.sendPacket(stopPkt)

	// 5. Chat oficial: &adminlevel (0120 00 0b "&adminlevel") + Chat con itemLinkDataCount=1 (103B)
	gmChat := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_CHAT_REQ)
	gmChat.WriteByte(0) // itemLinkDataCount
	gmChat.WriteByte(byte(len("&adminlevel")))
	gmChat.WriteBytes([]byte("&adminlevel"))
	_ = zClient.sendPacket(gmChat)

	linkChat := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_CHAT_REQ)
	linkChat.WriteByte(1) // itemLinkDataCount = 1
	msgTxt := "Mira mi espada!"
	linkChat.WriteByte(byte(len(msgTxt)))
	linkChat.WriteBytes([]byte(msgTxt))
	linkData := make([]byte, 103)
	binary.LittleEndian.PutUint16(linkData[0:2], 500)
	linkChat.WriteBytes(linkData)
	_ = zClient.sendPacket(linkChat)

	shoutTxt := "Aviso en mapa!"
	shoutPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SHOUT_CMD)
	shoutPkt.WriteByte(0) // itemLinkDataCount
	shoutPkt.WriteByte(byte(len(shoutTxt)))
	shoutPkt.WriteBytes([]byte(shoutTxt))
	_ = zClient.sendPacket(shoutPkt)

	var gotAdminNotice, gotLinkedChat, gotShout, gotFalseMoveFail bool
	_ = zConn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
	for !gotAdminNotice || !gotLinkedChat || !gotShout {
		pkt, rErr := auditReadServerPacket(zConn)
		if rErr != nil {
			break
		}
		payload := pkt.Raw[2:]
		switch pkt.OpCode {
		case protocol.PROTO_NC_ACT_MOVEFAIL_CMD:
			gotFalseMoveFail = true
		case protocol.PROTO_NC_ACT_NOTICE_CMD:
			if strings.Contains(string(payload), "Admin level is") {
				gotAdminNotice = true
			}
		case protocol.PROTO_NC_ACT_SOMEONECHAT_CMD:
			// Cabecera 7B: itemLinkDataCount(1) + handle(2) + len(1) + flag(1) + fontColor(1) + balloonColor(1)
			if len(payload) == 7+len(msgTxt)+103 && payload[0] == 1 && int(payload[3]) == len(msgTxt) {
				gotLinkedChat = true
			}
		case protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD:
			// Cabecera 23B: linkCount(1) + charID[20] + flag(1) + len(1)
			if len(payload) == 23+len(shoutTxt) && int(payload[22]) == len(shoutTxt) {
				gotShout = true
			}
		}
	}
	_ = zConn.SetReadDeadline(time.Time{})
	if gotFalseMoveFail {
		t.Fatalf("No debe emitirse PROTO_NC_ACT_MOVEFAIL_CMD en movimiento caminable")
	}
	if !gotAdminNotice {
		t.Fatalf("No se recibió PROTO_NC_ACT_NOTICE_CMD 'Admin level is' tras '&adminlevel'")
	}
	if !gotLinkedChat {
		t.Fatalf("No se recibió PROTO_NC_ACT_SOMEONECHAT_CMD con header 7B + itemLinkData (103B)")
	}
	if !gotShout {
		t.Fatalf("No se recibió PROTO_NC_ACT_SOMEONESHOUT_CMD con header 23B")
	}

	// Verificar catálogo AdminLvSet.shn en ZoneDataStore
	if zoneDataStore.AdminLvCatalog == nil || len(zoneDataStore.AdminLvCatalog.Levels) == 0 {
		t.Fatalf("AdminLvCatalog vacío tras LoadAllData")
	}

	// 6. Logout flow: 0x1071 LOGOUTREADY -> 0x1072 LOGOUTCANCEL -> 0x0C18 LogoutType=1 (Return Char Select)
	waitCond := func(timeout time.Duration, cond func() bool) bool {
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			if cond() {
				return true
			}
			time.Sleep(10 * time.Millisecond)
		}
		return cond()
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGOUTREADY_CMD))
	if !waitCond(500*time.Millisecond, func() bool { return zp.LogoutPending }) {
		t.Fatalf("LogoutPending esperado=true tras PROTO_NC_CHAR_LOGOUTREADY_CMD")
	}
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGOUTCANCEL_CMD))
	if !waitCond(500*time.Millisecond, func() bool { return !zp.LogoutPending }) {
		t.Fatalf("LogoutPending esperado=false tras PROTO_NC_CHAR_LOGOUTCANCEL_CMD")
	}

	zLogout := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	zLogout.WriteByte(1) // LogoutType = 1 (volver a Selección de Personaje)
	_ = zClient.sendPacket(zLogout)
	if !waitCond(500*time.Millisecond, func() bool {
		_, stillInZone := zoneSrv.Handler().GetPlayerByName(testCharName)
		return !stillInZone
	}) {
		t.Fatalf("Jugador '%s' no se desconectó de ZoneService tras NORMALLOGOUT", testCharName)
	}
	_ = zConn.Close()

	// CRÍTICO: tActiveAccountSession debe seguir viva tras desconectar de Zone
	var sessionCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo).Scan(&sessionCount); err != nil || sessionCount != 1 {
		t.Fatalf("tActiveAccountSession debe permanecer activa (count=1) tras salir de Zone hacia CharSelect, obtenido=%d err=%v", sessionCount, err)
	}

	// Enviar LogoutType=1 a WorldService: socket debe seguir abierto + sesión viva + responder AVATAR_LIST_REQ (0x0C1F)
	wLogoutCharSel := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	wLogoutCharSel.WriteByte(1)
	_ = wClient.sendPacket(wLogoutCharSel)
	time.Sleep(40 * time.Millisecond)

	if err := db.QueryRow("SELECT COUNT(*) FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo).Scan(&sessionCount); err != nil || sessionCount != 1 {
		t.Fatalf("tActiveAccountSession debe permanecer activa (count=1) en WorldService tras LogoutType=1, obtenido=%d err=%v", sessionCount, err)
	}

	avatarListReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_AVATAR_LIST_REQ)
	_ = wClient.sendPacket(avatarListReq)
	avatarListAck, err := auditReadServerPacket(wConn)
	if err != nil || avatarListAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("Tras LogoutType=1 WorldService debe permanecer abierto y responder PROTO_NC_USER_LOGINWORLD_ACK (0x0C14), got=%v err=%v", avatarListAck, err)
	}

	// Retorno a Selección de Mundos (0x0C33 -> 0x0C34 OTP -> LoginService 0x0C37)
	willWorldReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WILL_WORLD_SELECT_REQ)
	_ = wClient.sendPacket(willWorldReq)
	willWorldAck, err := auditReadServerPacket(wConn)
	if err != nil || willWorldAck.OpCode != protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK {
		t.Fatalf("PROTO_NC_USER_WILL_WORLD_SELECT_REQ falló: %v / %v", willWorldAck, err)
	}
	nErr, _ := willWorldAck.Reader.ReadUint16()
	otpBytes, _ := willWorldAck.Reader.ReadBytes(32)
	if nErr != 0 || len(otpBytes) != 32 {
		t.Fatalf("OTP inválido en PROTO_NC_USER_WILL_WORLD_SELECT_ACK: nError=%d len=%d", nErr, len(otpBytes))
	}
	_ = wConn.Close()

	lConn2, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
	if err != nil {
		t.Fatalf("Reconexión a LoginService falló: %v", err)
	}
	defer lConn2.Close()
	lHS2, _ := auditReadServerPacket(lConn2)
	lSeed2, _ := lHS2.Reader.ReadInt16()
	lCrypto2, _ := crypto.NewNetCrypto(lSeed2)
	lClient2 := &auditMockClient{conn: lConn2, crypto: lCrypto2}

	verPkt2 := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt2.WriteString(activeVerKey, 64)
	_ = lClient2.sendPacket(verPkt2)
	_, _ = auditReadServerPacket(lConn2) // 0x0C67
	_, _ = auditReadServerPacket(lConn2) // 0x0807

	otpReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ)
	otpReq.WriteBytes(otpBytes)
	_ = lClient2.sendPacket(otpReq)
	otpLoginAck, err := auditReadServerPacket(lConn2)
	if err != nil || otpLoginAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
		t.Fatalf("LoginService PROTO_NC_USER_LOGIN_WITH_OTP_REQ (0x0C37) falló: %v / %v", otpLoginAck, err)
	}

	// Salida total final: re-entrar a WorldService y emitir LogoutType=0 para verificar cierre de sesión MySQL + socket
	wSelReq2 := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	wSelReq2.WriteByte(0)
	_ = lClient2.sendPacket(wSelReq2)
	wSelAck2, err := auditReadServerPacket(lConn2)
	if err != nil || wSelAck2.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("WORLDSELECT tras OTP falló: %v / %v", wSelAck2, err)
	}
	_, _ = wSelAck2.Reader.ReadByte()
	_, _ = wSelAck2.Reader.ReadString(16)
	_, _ = wSelAck2.Reader.ReadUint16()
	tokenBytes2, _ := wSelAck2.Reader.ReadBytes(64)
	worldToken2 := strings.TrimRight(string(tokenBytes2), "\x00")
	_ = lConn2.Close()

	wConn2, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
	if err != nil {
		t.Fatalf("Dial WorldService final falló: %v", err)
	}
	defer wConn2.Close()
	wHS2, _ := auditReadServerPacket(wConn2)
	wSeed2, _ := wHS2.Reader.ReadInt16()
	wCrypto2, _ := crypto.NewNetCrypto(wSeed2)
	wClient2 := &auditMockClient{conn: wConn2, crypto: wCrypto2}

	wlReq2 := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq2.WriteString(testUser, 256)
	wlReq2.WriteString(worldToken2, 64)
	_ = wClient2.sendPacket(wlReq2)
	wlAck2, err := auditReadServerPacket(wConn2)
	if err != nil || wlAck2.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD final falló: %v / %v", wlAck2, err)
	}

	fullLogout := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	fullLogout.WriteByte(0) // LogoutType = 0
	_ = wClient2.sendPacket(fullLogout)
	time.Sleep(60 * time.Millisecond)

	if err := db.QueryRow("SELECT COUNT(*) FROM tActiveAccountSession WHERE nUserNo = ?", testUserNo).Scan(&sessionCount); err != nil || sessionCount != 0 {
		t.Fatalf("tActiveAccountSession debía cerrarse (count=0) tras LogoutType=0, obtenido=%d err=%v", sessionCount, err)
	}
}
