package audit_test

import (
	"crypto/md5"
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

func f23DialAndHandshake(t *testing.T, addr string) (net.Conn, *auditMockClient) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		t.Fatalf("Dial %s failed: %v", addr, err)
	}
	hs, err := auditReadServerPacket(conn)
	if err != nil || hs.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		_ = conn.Close()
		t.Fatalf("Handshake SEED_ACK on %s failed: pkt=%v err=%v", addr, hs, err)
	}
	seed, _ := hs.Reader.ReadInt16()
	cr, err := crypto.NewNetCrypto(seed)
	if err != nil {
		_ = conn.Close()
		t.Fatalf("NewNetCrypto failed: %v", err)
	}
	return conn, &auditMockClient{conn: conn, crypto: cr}
}

func f23LoginAndGetWorldToken(t *testing.T, loginAddr, verKey, uName, passHash string) string {
	t.Helper()
	lConn, lClient := f23DialAndHandshake(t, loginAddr)
	defer lConn.Close()

	verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt.WriteString(verKey, 64)
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
	lAck, err := auditReadServerPacket(lConn)
	if err != nil || lAck.OpCode != protocol.PROTO_NC_USER_LOGIN_ACK {
		t.Fatalf("US_LOGIN_REQ failed: pkt=%v err=%v", lAck, err)
	}

	wSel := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	wSel.WriteByte(0)
	_ = lClient.sendPacket(wSel)
	wSelAck, err := auditReadServerPacket(lConn)
	if err != nil || wSelAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("WORLDSELECT_REQ failed: pkt=%v err=%v", wSelAck, err)
	}
	_, _ = wSelAck.Reader.ReadByte()
	_, _ = wSelAck.Reader.ReadString(16)
	_, _ = wSelAck.Reader.ReadUint16()
	tokenBytes, _ := wSelAck.Reader.ReadBytes(64)
	return strings.TrimRight(string(tokenBytes), "\x00")
}

func f23CollectPackets(conn net.Conn, wait time.Duration) []*protocol.Packet {
	var out []*protocol.Packet
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		rem := time.Until(deadline)
		if rem > 120*time.Millisecond {
			rem = 120 * time.Millisecond
		}
		if rem < 20*time.Millisecond {
			rem = 20 * time.Millisecond
		}
		_ = conn.SetReadDeadline(time.Now().Add(rem))
		p, err := auditReadServerPacket(conn)
		if err != nil || p == nil {
			continue
		}
		out = append(out, p)
	}
	_ = conn.SetReadDeadline(time.Time{})
	return out
}

func f23FindPkt(pkts []*protocol.Packet, op uint16) *protocol.Packet {
	for _, p := range pkts {
		if p != nil && p.OpCode == op {
			return p
		}
	}
	return nil
}

func f23EnterZone(t *testing.T, zoneAddr string, ds *zone.DataStore, randomID uint16, charName string) (net.Conn, *auditMockClient, []*protocol.Packet) {
	t.Helper()
	zConn, zClient := f23DialAndHandshake(t, zoneAddr)
	mapReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	mapReq.WriteUint16(randomID)
	mapReq.WriteString(charName, 20)
	sums := ds.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		mapReq.WriteBytes(sums[i][:])
	}
	_ = zClient.sendPacket(mapReq)

	var pkts []*protocol.Packet
	loginDeadline := time.Now().Add(2500 * time.Millisecond)
	gotAck := false
	for time.Now().Before(loginDeadline) {
		_ = zConn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		p, err := auditReadServerPacket(zConn)
		if err != nil || p == nil {
			continue
		}
		pkts = append(pkts, p)
		if p.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			gotAck = true
			break
		}
	}
	_ = zConn.SetReadDeadline(time.Time{})
	if !gotAck {
		_ = zConn.Close()
		t.Fatalf("Zone MAP_LOGIN_REQ para '%s' no emitió PROTO_NC_MAP_LOGIN_ACK", charName)
	}
	// Drenar cualquier paquete adicional del In-Game Chunk ya en vuelo antes de LOGINCOMPLETE
	pkts = append(pkts, f23CollectPackets(zConn, 200*time.Millisecond)...)
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	readyPkts := f23CollectPackets(zConn, 300*time.Millisecond)
	pkts = append(pkts, readyPkts...)
	return zConn, zClient, pkts
}

// TestFase23_WorldEntryAvatarAndCharSelectBranches_ClusterE2E audita rama por rama
// en WorldManager.exe todos los caminos de éxito y rechazo de NC_USER, NC_AVATAR y NC_CHAR:
//  - LOGINWORLD_REQ token inválido -> 0x0C15 (0x0044) vs válido -> 0x0C14.
//  - AVATAR_CREATE_REQ BadName/espacios -> 0x1404 (132), Crusader sin Lv>=ChrLevel_CanCreateSen -> 0x1404 (134),
//    nombre duplicado -> 0x1404 (132), slot ocupado -> 0x1404 (133), válido -> 0x1406.
//  - AVATAR_RENAME_REQ BadName/corto -> 0x1410 (0x00D1), válido -> 0x1410 (0x00D0).
//  - AVATAR_ERASE_REQ slot vacío/inválido -> 0x140A (0x00C6), válido -> 0x140C.
//  - CHAR_LOGIN_REQ slot>5 -> 0x1002 (0x0143), slot vacío -> 0x1002 (0x0144), mapa inválido -> 0x1002 (0x0145), válido -> 0x1003.
func TestFase23_WorldEntryAvatarAndCharSelectBranches_ClusterE2E(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	accountRepo := database.NewMySQLAccountRepository(db)
	charRepo := database.NewMySQLCharacterRepository(db)
	transfers := database.NewTransferManager(db)

	worldDS, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	activeVerKey, err := accountRepo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("GetClientVersionKey: %v", err)
	}

	const (
		loginPort = 19702
		worldPort = 19704
		zonePort  = 19706
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
		Name:     "Fase23World",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort,
	}, charRepo, transfers, worldDS, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("WorldServer Start: %v", err)
	}
	defer worldSrv.Stop()

	if err := charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort)); err != nil {
		t.Fatalf("UpsertZoneNode: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort = ?", zonePort)
	}()

	time.Sleep(80 * time.Millisecond)

	const userNo = 998923
	ts := time.Now().Unix() % 100000
	pass := "f23pass"
	sum := md5.Sum([]byte(pass))
	passHash := hex.EncodeToString(sum[:])

	cleanupUser := func() {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", userNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", userNo)
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", userNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", userNo)
	}
	cleanupUser()
	defer cleanupUser()

	uName := fmt.Sprintf("f23w_%05d", ts)
	if _, err := db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
	`, userNo, uName, passHash); err != nil {
		t.Fatalf("INSERT tUser: %v", err)
	}

	// 1) Token de World inválido -> PROTO_NC_USER_LOGINWORLDFAIL_ACK (0x0C15, 0x0044)
	wConnBad, wClientBad := f23DialAndHandshake(t, fmt.Sprintf("127.0.0.1:%d", worldPort))
	badWL := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	badWL.WriteString(uName, 256)
	badWL.WriteString("INVALID_WORLD_TOKEN_XYZ", 64)
	_ = wClientBad.sendPacket(badWL)
	badWLAck, err := auditReadServerPacket(wConnBad)
	_ = wConnBad.Close()
	if err != nil || badWLAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLDFAIL_ACK {
		t.Fatalf("Esperado 0x0C15 LOGINWORLDFAIL_ACK con token inválido, obtenido: %v / %v", badWLAck, err)
	}
	errCode, _ := badWLAck.Reader.ReadUint16()
	if errCode != uint16(protocol.ErrWrongRegion) {
		t.Fatalf("Código esperado en 0x0C15=%#x, obtenido=%#x", protocol.ErrWrongRegion, errCode)
	}

	// 2) Login válido en LoginService -> WorldService -> 0x0C14
	worldToken := f23LoginAndGetWorldToken(t, fmt.Sprintf("127.0.0.1:%d", loginPort), activeVerKey, uName, passHash)
	wConn, wClient := f23DialAndHandshake(t, fmt.Sprintf("127.0.0.1:%d", worldPort))
	defer wConn.Close()

	okWL := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	okWL.WriteString(uName, 256)
	okWL.WriteString(worldToken, 64)
	_ = wClient.sendPacket(okWL)
	okWLAck, err := auditReadServerPacket(wConn)
	if err != nil || okWLAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("Esperado 0x0C14 LOGINWORLD_ACK, obtenido: %v / %v", okWLAck, err)
	}

	sendCreate := func(slot byte, name string, classID byte) (uint16, uint16) {
		req := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
		req.WriteByte(slot)
		req.WriteString(name, 20)
		req.WriteByte(byte((classID << 2) | (1 << 7)))
		req.WriteByte(1)
		req.WriteByte(1)
		req.WriteByte(4)
		_ = wClient.sendPacket(req)
		ack, err := auditReadServerPacket(wConn)
		if err != nil {
			t.Fatalf("read AVATAR_CREATE response: %v", err)
		}
		code, _ := ack.Reader.ReadUint16()
		return ack.OpCode, code
	}

	// 3) Ramas de rechazo de NC_AVATAR_CREATE_REQ
	if op, code := sendCreate(0, "AdminGuy", 1); op != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK || code != uint16(protocol.CreateCharError_NameTaken) {
		t.Fatalf("BadName esperaba (0x1404, 132), obtenido=(%#x, %d)", op, code)
	}
	if op, code := sendCreate(0, "Bad Space", 1); op != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK || code != uint16(protocol.CreateCharError_NameTaken) {
		t.Fatalf("Espacio en nombre esperaba (0x1404, 132), obtenido=(%#x, %d)", op, code)
	}
	if op, code := sendCreate(0, fmt.Sprintf("Cru%05d", ts), byte(protocol.Job_Crusader)); op != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK || code != 0x0086 {
		t.Fatalf("Crusader sin nivel 60 esperaba (0x1404, 0x0086), obtenido=(%#x, %#x)", op, code)
	}

	// Crear personaje válido en slot 0
	cMain := fmt.Sprintf("F23A%04d", ts%10000)
	if op, _ := sendCreate(0, cMain, 1); op != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("Creación válida en slot 0 esperaba 0x1406, obtenido=%#x", op)
	}
	// Nombre duplicado -> 132
	if op, code := sendCreate(1, cMain, 1); op != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK || code != uint16(protocol.CreateCharError_NameTaken) {
		t.Fatalf("Nombre duplicado esperaba (0x1404, 132), obtenido=(%#x, %d)", op, code)
	}
	// Slot ocupado -> 133
	if op, code := sendCreate(0, fmt.Sprintf("F23B%04d", ts%10000), 1); op != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK || code != uint16(protocol.CreateCharError_ErrorInMaxSlot) {
		t.Fatalf("Slot ocupado esperaba (0x1404, 133), obtenido=(%#x, %d)", op, code)
	}

	// 4) Ramas de RENAME_REQ (0x00D1 en BadName vs 0x00D0 en éxito) y ERASE_REQ (0x00C6 en slot vacío/inválido vs 0x140C)
	sendRename := func(slot byte, newName string) uint16 {
		req := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_RENAME_REQ)
		req.WriteByte(slot)
		req.WriteString(newName, 20)
		_ = wClient.sendPacket(req)
		ack, err := auditReadServerPacket(wConn)
		if err != nil || ack.OpCode != protocol.PROTO_NC_AVATAR_RENAME_ACK {
			t.Fatalf("RENAME_ACK esperado: pkt=%v err=%v", ack, err)
		}
		_, _ = ack.Reader.ReadByte()
		_, _ = ack.Reader.ReadString(20)
		errC, _ := ack.Reader.ReadUint16()
		return errC
	}
	if errC := sendRename(0, "AdminBad"); errC != 0x00D1 {
		t.Fatalf("Rename con BadName esperaba 0x00D1, obtenido=%#x", errC)
	}
	cRenamed := fmt.Sprintf("F23R%04d", ts%10000)
	if errC := sendRename(0, cRenamed); errC != 0x00D0 {
		t.Fatalf("Rename válido esperaba 0x00D0, obtenido=%#x", errC)
	}

	// Erase en slot vacío (slot 2) -> 0x140A con 0x00C6
	eraseEmpty := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	eraseEmpty.WriteByte(2)
	_ = wClient.sendPacket(eraseEmpty)
	eraseFailAck, err := auditReadServerPacket(wConn)
	if err != nil || eraseFailAck.OpCode != protocol.PROTO_NC_AVATAR_ERASEFAIL_ACK {
		t.Fatalf("Esperado 0x140A ERASEFAIL_ACK en slot vacío, obtenido=%v err=%v", eraseFailAck, err)
	}
	if eraseErr, _ := eraseFailAck.Reader.ReadUint16(); eraseErr != 0x00C6 {
		t.Fatalf("ERASEFAIL_ACK esperaba 0x00C6, obtenido=%#x", eraseErr)
	}

	// Crear y borrar slot 1 con éxito -> 0x140C
	cDel := fmt.Sprintf("F23D%04d", ts%10000)
	if op, _ := sendCreate(1, cDel, 1); op != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("Create slot 1 esperaba 0x1406, obtenido=%#x", op)
	}
	eraseOK := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	eraseOK.WriteByte(1)
	_ = wClient.sendPacket(eraseOK)
	if eraseSucc, err := auditReadServerPacket(wConn); err != nil || eraseSucc.OpCode != protocol.PROTO_NC_AVATAR_ERASESUCC_ACK {
		t.Fatalf("Erase válido esperaba 0x140C, obtenido=%v err=%v", eraseSucc, err)
	}

	// 5) Ramas de CHAR_LOGIN_REQ (0x1001): slot > 5 (0x0143), slot vacío (0x0144), mapa inválido (0x0145), válido (0x1003)
	sendCharSel := func(slot byte) (*protocol.Packet, uint16) {
		req := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
		req.WriteByte(slot)
		_ = wClient.sendPacket(req)
		ack, err := auditReadServerPacket(wConn)
		if err != nil {
			t.Fatalf("CHAR_LOGIN read failed: %v", err)
		}
		var code uint16
		if ack.OpCode == protocol.PROTO_NC_CHAR_LOGINFAIL_ACK {
			code, _ = ack.Reader.ReadUint16()
		}
		return ack, code
	}

	if ack, code := sendCharSel(6); ack.OpCode != protocol.PROTO_NC_CHAR_LOGINFAIL_ACK || code != 0x0143 {
		t.Fatalf("Slot 6 esperaba (0x1002, 0x0143), obtenido=(%#x, %#x)", ack.OpCode, code)
	}
	if ack, code := sendCharSel(3); ack.OpCode != protocol.PROTO_NC_CHAR_LOGINFAIL_ACK || code != 0x0144 {
		t.Fatalf("Slot vacío 3 esperaba (0x1002, 0x0144), obtenido=(%#x, %#x)", ack.OpCode, code)
	}

	// Mapa inválido en DB -> 0x0145
	chA, chErr := charRepo.GetCharacterByName(cRenamed)
	if chErr != nil || chA == nil {
		t.Fatalf("GetCharacterByName(%s): %v", cRenamed, chErr)
	}
	if err := charRepo.SaveCharLocation(chA.ID, "NoMap99", chA.XPos, chA.YPos, 0); err != nil {
		t.Fatalf("SaveCharLocation NoMap99: %v", err)
	}
	if ack, code := sendCharSel(0); ack.OpCode != protocol.PROTO_NC_CHAR_LOGINFAIL_ACK || code != 0x0145 {
		t.Fatalf("Mapa inválido esperaba (0x1002, 0x0145), obtenido=(%#x, %#x)", ack.OpCode, code)
	}
	if err := charRepo.SaveCharLocation(chA.ID, "RouN", chA.XPos, chA.YPos, 0); err != nil {
		t.Fatalf("SaveCharLocation RouN: %v", err)
	}

	// Selección válida -> 0x1003 + 0x7032..0x7034
	if ack, _ := sendCharSel(0); ack.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CharSelect válido esperaba 0x1003, obtenido=%#x", ack.OpCode)
	}
	for _, optOp := range []uint16{
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD,
	} {
		if optPkt, err := auditReadServerPacket(wConn); err != nil || optPkt.OpCode != optOp {
			t.Fatalf("Esperado %#x tras 0x1003, obtenido=%v err=%v", optOp, optPkt, err)
		}
	}

	// 6) Whisper a destinatario desconectado -> PROTO_NC_ACT_WHISPERFAIL_ACK (0x200E) con 0x0F69 (3945)
	wsp := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPER_REQ)
	wsp.WriteByte(0)
	wsp.WriteString("OfflineTargetXYZ", 20)
	msg := "hola"
	wsp.WriteByte(byte(len(msg)))
	wsp.WriteBytes([]byte(msg))
	_ = wClient.sendPacket(wsp)
	wspAck, err := auditReadServerPacket(wConn)
	if err != nil || wspAck.OpCode != protocol.PROTO_NC_ACT_WHISPERFAIL_ACK {
		t.Fatalf("Whisper offline esperaba 0x200E, obtenido=%v err=%v", wspAck, err)
	}
	if wErr, _ := wspAck.Reader.ReadUint16(); wErr != 0x0F69 {
		t.Fatalf("Whisper offline esperaba error 0x0F69 (3945), obtenido=%#x", wErr)
	}
}

// TestFase23_ZoneMoveCollisionPortalsNPCsMoverAOIChatGM_ClusterE2E audita rama por rama en Zone.exe:
//  - TransferKey inválido -> 0x1807 PROTO_NC_MAP_LOGINFAIL_ACK.
//  - Movimiento caminable (MOVEWALK/MOVERUN/STOP/JUMP) + colisión .shbd / paso > 0x9C40 -> 0x201B MOVEFAIL_CMD.
//  - AOI BRIEFINFO_INFORM_CMD + MISC_HIDE_EXCEPT_ME_ON/OFF_CMD (0x0817/0x0818).
//  - Portales: LINKSAME (0x1809), LINKOTHER inter-zona (0x180A), MULTY_LINK_SELECT (0x2029/0x202B/0),
//    WING_FLY (0x1143 en !LinkOut vs 0x1141), TOWNPORTAL (0x19C2 muerto / 0x19C3 índice / 0x19C1 ok).
//  - Interacción NPC: NPCCLICK fuera de rango (>0xF4240) / inválido vs RouSmithJames (0x201C -> ACK 1 -> 0x3C09),
//    StoreManager (0x3C08).
//  - Monturas NC_MOVER: RIDE_ON -> 0xCC02, ya montado -> 0xCC03 (0x306A), hnd=0 -> 0xCC03 (0x306B),
//    RIDE_OFF -> 0xCC06, desmontado -> 0xCC07 (0x3081).
//  - Chat (0x2A normal vs 0x03 GM), Shout, Roar y rechazo/éxito AdminLvSet.shn en &adminlevel.
func TestFase23_ZoneMoveCollisionPortalsNPCsMoverAOIChatGM_ClusterE2E(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	accountRepo := database.NewMySQLAccountRepository(db)
	charRepo := database.NewMySQLCharacterRepository(db)
	gameLogRepo := database.NewMySQLGameLogRepository(db)
	transfers := database.NewTransferManager(db)

	worldDS, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	zoneDS, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData: %v", err)
	}
	activeVerKey, err := accountRepo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("GetClientVersionKey: %v", err)
	}

	const (
		loginPort = 19712
		worldPort = 19714
		zonePort1 = 19724
		zonePort2 = 19726
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
		Name:     "Fase23ZoneCluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort1,
	}, charRepo, transfers, worldDS, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("WorldServer Start: %v", err)
	}
	defer worldSrv.Stop()

	if err := charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort1)); err != nil {
		t.Fatalf("UpsertZoneNode 1: %v", err)
	}
	if err := charRepo.UpsertZoneNode(0, 2, "127.0.0.1", uint16(zonePort2)); err != nil {
		t.Fatalf("UpsertZoneNode 2: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort IN (?, ?)", zonePort1, zonePort2)
	}()

	glDisp := database.NewGameLogDispatcher(gameLogRepo, 1000, 2)
	glDisp.Start()
	defer glDisp.Close(t.Context())

	zoneSrv1 := zone.NewServer(&zone.Config{
		BindAddr:  fmt.Sprintf("127.0.0.1:%d", zonePort1),
		Port:      zonePort1,
		IP:        "127.0.0.1",
		WorldNo:   0,
		ZoneNo:    1,
		MultiNode: true,
	}, charRepo, transfers, zoneDS.ClassStore, zoneDS.FieldStore, zoneDS.ChkEngine)
	zoneSrv1.SetDataStore(zoneDS)
	zoneSrv1.SetAccountRepo(accountRepo)
	zoneSrv1.SetGameLogDispatcher(glDisp)
	if err := zoneSrv1.Start(); err != nil {
		t.Fatalf("Zone1 Start: %v", err)
	}
	defer zoneSrv1.Stop()

	zoneSrv2 := zone.NewServer(&zone.Config{
		BindAddr:  fmt.Sprintf("127.0.0.1:%d", zonePort2),
		Port:      zonePort2,
		IP:        "127.0.0.1",
		WorldNo:   0,
		ZoneNo:    2,
		MultiNode: true,
	}, charRepo, transfers, zoneDS.ClassStore, zoneDS.FieldStore, zoneDS.ChkEngine)
	zoneSrv2.SetDataStore(zoneDS)
	zoneSrv2.SetAccountRepo(accountRepo)
	zoneSrv2.SetGameLogDispatcher(glDisp)
	if err := zoneSrv2.Start(); err != nil {
		t.Fatalf("Zone2 Start: %v", err)
	}
	defer zoneSrv2.Stop()

	time.Sleep(100 * time.Millisecond)

	const (
		userA = 998924
		userB = 998925
	)
	ts := time.Now().Unix() % 100000
	pass := "f23zpass"
	sum := md5.Sum([]byte(pass))
	passHash := hex.EncodeToString(sum[:])

	cleanupUsers := func() {
		for _, u := range []int{userA, userB} {
			_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", u)
			_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", u)
			if rows, qErr := db.Query("SELECT nCharNo FROM tCharacter WHERE nUserNo = ?", u); qErr == nil {
				var cids []int
				for rows.Next() {
					var cid int
					if rows.Scan(&cid) == nil {
						cids = append(cids, cid)
					}
				}
				_ = rows.Close()
				for _, cid := range cids {
					_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", cid)
					_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", cid)
					_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", cid)
				}
			}
			_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", u)
			_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", u)
		}
	}
	cleanupUsers()
	defer cleanupUsers()

	uNameA := fmt.Sprintf("f23a_%05d", ts)
	uNameB := fmt.Sprintf("f23b_%05d", ts)
	cNameA := fmt.Sprintf("Z23A%04d", ts%10000)
	cNameB := fmt.Sprintf("Z23B%04d", ts%10000)

	_, err = db.Exec(`INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate, nAuthID) VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW(), 1)`, userA, uNameA, passHash)
	if err != nil {
		t.Fatalf("insert userA: %v", err)
	}
	_, err = db.Exec(`INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate, nAuthID) VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW(), 1)`, userB, uNameB, passHash)
	if err != nil {
		t.Fatalf("insert userB: %v", err)
	}

	spawnIntoZone := func(uName, cName string, adminLv int) (uint16, int) {
		tok := f23LoginAndGetWorldToken(t, fmt.Sprintf("127.0.0.1:%d", loginPort), activeVerKey, uName, passHash)
		wConn, wClient := f23DialAndHandshake(t, fmt.Sprintf("127.0.0.1:%d", worldPort))
		defer wConn.Close()

		wl := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
		wl.WriteString(uName, 256)
		wl.WriteString(tok, 64)
		_ = wClient.sendPacket(wl)
		if ack, err := auditReadServerPacket(wConn); err != nil || ack.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
			t.Fatalf("LOGINWORLD failed: %v / %v", ack, err)
		}

		cr := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
		cr.WriteByte(0)
		cr.WriteString(cName, 20)
		cr.WriteByte(byte((1 << 2) | (1 << 7)))
		cr.WriteByte(1)
		cr.WriteByte(1)
		cr.WriteByte(4)
		_ = wClient.sendPacket(cr)
		if ack, err := auditReadServerPacket(wConn); err != nil || ack.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
			t.Fatalf("AVATAR_CREATE failed: %v / %v", ack, err)
		}

		created, gErr := charRepo.GetCharacterByName(cName)
		if gErr != nil || created == nil {
			t.Fatalf("GetCharacterByName(%s): %v", cName, gErr)
		}
		if err := charRepo.SaveCharLocation(created.ID, "RouN", 6888, 7045, 0); err != nil {
			t.Fatalf("SaveCharLocation(%s): %v", cName, err)
		}
		if adminLv > 0 {
			if _, err := charRepo.SetAdminLevel(created.ID, byte(adminLv)); err != nil {
				t.Fatalf("SetAdminLevel(%s): %v", cName, err)
			}
		}

		sel := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
		sel.WriteByte(0)
		_ = wClient.sendPacket(sel)
		if ack, err := auditReadServerPacket(wConn); err != nil || ack.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
			t.Fatalf("CHAR_LOGIN failed: %v / %v", ack, err)
		}
		_, _ = auditReadServerPacket(wConn)
		_, _ = auditReadServerPacket(wConn)
		_, _ = auditReadServerPacket(wConn)

		var rnd uint16
		if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", cName).Scan(&rnd); err != nil {
			t.Fatalf("read randomID: %v", err)
		}
		return rnd, created.ID
	}

	rndA, charIDA := spawnIntoZone(uNameA, cNameA, 0)
	rndB, _ := spawnIntoZone(uNameB, cNameB, 100)

	// 1) Token de Zone inválido -> PROTO_NC_MAP_LOGINFAIL_ACK (0x1807)
	zBadConn, zBadClient := f23DialAndHandshake(t, fmt.Sprintf("127.0.0.1:%d", zonePort1))
	badZReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	badZReq.WriteUint16(64999)
	badZReq.WriteString(cNameA, 20)
	sums := zoneDS.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		badZReq.WriteBytes(sums[i][:])
	}
	_ = zBadClient.sendPacket(badZReq)
	zBadAck, err := auditReadServerPacket(zBadConn)
	_ = zBadConn.Close()
	if err != nil || zBadAck.OpCode != protocol.PROTO_NC_MAP_LOGINFAIL_ACK {
		t.Fatalf("Zone transferKey inválido esperaba 0x1807, obtenido=%v err=%v", zBadAck, err)
	}

	// 2) Entrar a Zone1 con Player A (normal) y Player B (GM)
	zoneAddr1 := fmt.Sprintf("127.0.0.1:%d", zonePort1)
	zConnA, zClientA, initA := f23EnterZone(t, zoneAddr1, zoneDS, rndA, cNameA)
	defer zConnA.Close()
	if f23FindPkt(initA, protocol.PROTO_NC_BRIEFINFO_MOB_CMD) == nil {
		t.Fatalf("Player A no recibió PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09) tras LOGINCOMPLETE")
	}

	zConnB, zClientB, _ := f23EnterZone(t, zoneAddr1, zoneDS, rndB, cNameB)
	defer zConnB.Close()

	// Consumir spawn mutuo inicial
	_ = f23CollectPackets(zConnA, 100*time.Millisecond)
	_ = f23CollectPackets(zConnB, 100*time.Millisecond)

	pA, okA := zoneSrv1.Handler().GetPlayerByName(cNameA)
	pB, okB := zoneSrv1.Handler().GetPlayerByName(cNameB)
	if !okA || pA == nil || !okB || pB == nil || pA.CurrentMap == nil {
		t.Fatalf("No se pudieron recuperar Player A/B en Zone1")
	}

	// 3) AOI: BRIEFINFO_INFORM_CMD + HIDE_EXCEPT_ME_ON/OFF
	informReq := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_INFORM_CMD)
	informReq.WriteUint16(pA.MapObjectID)
	informReq.WriteUint16(0)
	informReq.WriteUint16(pB.MapObjectID)
	_ = zClientA.sendPacket(informReq)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD) == nil {
		t.Fatalf("BRIEFINFO_INFORM_CMD no devolvió LOGINCHARACTER_CMD de Player B")
	}

	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HIDE_EXCEPT_ME_ON_CMD))
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD) == nil {
		t.Fatalf("HIDE_EXCEPT_ME_ON no emitió BRIEFINFODELETE_CMD")
	}
	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HIDE_EXCEPT_ME_OFF_CMD))
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD) == nil {
		t.Fatalf("HIDE_EXCEPT_ME_OFF no re-emitió LOGINCHARACTER_CMD")
	}

	// 4) Movimiento & Colisiones .shbd (MOVEWALK/MOVERUN/STOP/JUMP vs MOVEFAIL_CMD 0x201B)
	walk := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVEWALK_CMD)
	walk.WriteInt32(6888)
	walk.WriteInt32(7045)
	walk.WriteInt32(6920)
	walk.WriteInt32(7045)
	_ = zClientA.sendPacket(walk)
	pktsB := f23CollectPackets(zConnB, 100*time.Millisecond)
	if f23FindPkt(pktsB, protocol.PROTO_NC_ACT_SOMEONEMOVEWALK_CMD) == nil {
		t.Fatalf("MOVEWALK caminable no emitió SOMEONEMOVEWALK_CMD en Player B")
	}
	pktsA := f23CollectPackets(zConnA, 60*time.Millisecond)
	if f23FindPkt(pktsA, protocol.PROTO_NC_ACT_MOVEFAIL_CMD) != nil {
		t.Fatalf("MOVEWALK caminable no debía emitir MOVEFAIL_CMD")
	}

	run := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVERUN_CMD)
	run.WriteInt32(6920)
	run.WriteInt32(7045)
	run.WriteInt32(6960)
	run.WriteInt32(7045)
	_ = zClientA.sendPacket(run)
	if pkts := f23CollectPackets(zConnB, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_SOMEONEMOVERUN_CMD) == nil {
		t.Fatalf("MOVERUN caminable no emitió SOMEONEMOVERUN_CMD")
	}

	stop := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_STOP_REQ)
	stop.WriteInt32(6960)
	stop.WriteInt32(7045)
	_ = zClientA.sendPacket(stop)
	if pkts := f23CollectPackets(zConnB, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_SOMEONESTOP_CMD) == nil {
		t.Fatalf("STOP válido no emitió SOMEONESTOP_CMD")
	}

	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ACT_JUMP_CMD))
	if pkts := f23CollectPackets(zConnB, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_SOMEEONEJUMP_CMD) == nil {
		t.Fatalf("JUMP no emitió SOMEEONEJUMP_CMD")
	}

	// Buscar un punto bloqueado real en RouN.shbd + probar salto excesivo (> 0x9C40) -> MOVEFAIL_CMD (0x201B, 8B)
	rouSHBD, err := blockinfo.LoadSHBDFile(filepath.Join(dataDir, "BlockInfo", "RouN.shbd"))
	if err != nil {
		t.Fatalf("LoadSHBDFile RouN: %v", err)
	}
	var blockedX, blockedY int32 = 999999, 999999
	for candX := 1000; candX <= 20000; candX += 200 {
		for candY := 1000; candY <= 20000; candY += 200 {
			if rouSHBD.IsWorldMoveBlocked(float64(candX), float64(candY)) {
				blockedX, blockedY = int32(candX), int32(candY)
				break
			}
		}
		if blockedX != 999999 {
			break
		}
	}
	badMove := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVERUN_CMD)
	badMove.WriteInt32(int32(pA.Char.XPos))
	badMove.WriteInt32(int32(pA.Char.YPos))
	badMove.WriteInt32(blockedX)
	badMove.WriteInt32(blockedY)
	_ = zClientA.sendPacket(badMove)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		mf := f23FindPkt(pkts, protocol.PROTO_NC_ACT_MOVEFAIL_CMD)
		return mf == nil || len(mf.Raw)-2 != 8
	}() {
		t.Fatalf("Movimiento bloqueado/excesivo debía emitir MOVEFAIL_CMD (0x201B) de 8B payload")
	}

	// 5) Chat (flag 0x02 normal vs 0x03 GM en PROTO_NC_ACT_SOMEONECHAT_CMD), Shout, Roar y AdminLvSet.shn en &adminlevel
	sendChat := func(cl *auditMockClient, msg string) {
		p := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_CHAT_REQ)
		p.WriteByte(0)
		p.WriteByte(byte(len(msg)))
		p.WriteBytes([]byte(msg))
		_ = cl.sendPacket(p)
	}
	sendChat(zClientA, "Hola Mundo")
	_ = f23CollectPackets(zConnA, 80*time.Millisecond)
	pktsB = f23CollectPackets(zConnB, 100*time.Millisecond)
	if chatB := f23FindPkt(pktsB, protocol.PROTO_NC_ACT_SOMEONECHAT_CMD); chatB == nil {
		t.Fatalf("Normal chat no recibido por B")
	} else {
		_, _ = chatB.Reader.ReadByte()   // linkCount
		_, _ = chatB.Reader.ReadUint16() // MapObjectID
		_, _ = chatB.Reader.ReadByte()   // len(text)
		cFlag, _ := chatB.Reader.ReadByte()
		if cFlag != 0x02 {
			t.Fatalf("Flag chat normal esperado 0x02, obtenido=%#x", cFlag)
		}
	}

	sendChat(zClientB, "Aviso GM")
	_ = f23CollectPackets(zConnB, 80*time.Millisecond)
	pktsA = f23CollectPackets(zConnA, 100*time.Millisecond)
	if chatA := f23FindPkt(pktsA, protocol.PROTO_NC_ACT_SOMEONECHAT_CMD); chatA == nil {
		t.Fatalf("Chat GM no recibido por A")
	} else {
		_, _ = chatA.Reader.ReadByte()   // linkCount
		_, _ = chatA.Reader.ReadUint16() // MapObjectID
		_, _ = chatA.Reader.ReadByte()   // len(text)
		cFlag, _ := chatA.Reader.ReadByte()
		if cFlag != 0x03 {
			t.Fatalf("Flag chat GM esperado 0x03, obtenido=%#x", cFlag)
		}
	}

	// Shout & Roar
	shoutMsg := "Grito23"
	shout := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SHOUT_CMD)
	shout.WriteByte(0) // linkCount
	shout.WriteByte(byte(len(shoutMsg)))
	shout.WriteBytes([]byte(shoutMsg))
	_ = zClientA.sendPacket(shout)
	if pkts := f23CollectPackets(zConnB, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD) == nil {
		t.Fatalf("SHOUT_CMD no emitió SOMEONESHOUT_CMD")
	}

	roarMsg := "Rugido23"
	roar := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_ROAR_REQ)
	roar.WriteByte(0) // slot
	roar.WriteByte(byte(len(roarMsg)))
	roar.WriteBytes([]byte(roarMsg))
	_ = zClientA.sendPacket(roar)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_ROAR_ACK) == nil {
		t.Fatalf("ROAR_REQ no emitió ROAR_ACK")
	}

	// Comando GM sin permisos (Player A Admin=0) -> rechazado con notice de nivel insuficiente
	sendChat(zClientA, "&adminlevel")
	pktsA = f23CollectPackets(zConnA, 120*time.Millisecond)
	if f23FindPkt(pktsA, protocol.PROTO_NC_ACT_NOTICE_CMD) != nil {
		t.Fatalf("Player A ejecutando &adminlevel debía recibir NOTICE_CMD de rechazo")
	}
	sendChat(zClientB, "&adminlevel")
	pktsB = f23CollectPackets(zConnB, 120*time.Millisecond)
	if f23FindPkt(pktsB, protocol.PROTO_NC_ACT_NOTICE_CMD) == nil {
		t.Fatalf("Player B (GM 100) ejecutando &adminlevel debía recibir NOTICE_CMD")
	}

	// 6) NC_MOVER: inválido hnd=0 (0xCC03 0x306B), RIDE_ON (0xCC02), doble RIDE_ON (0xCC03 0x306A), RIDE_OFF (0xCC06), doble RIDE_OFF (0xCC07 0x3081)
	rideZero := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_ON_REQ)
	rideZero.WriteUint16(0)
	rideZero.WriteByte(0)
	_ = zClientA.sendPacket(rideZero)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		f := f23FindPkt(pkts, protocol.PROTO_NC_MOVER_RIDE_ON_FAIL_CMD)
		if f == nil {
			return true
		}
		code, _ := f.Reader.ReadUint16()
		return code != 0x306B
	}() {
		t.Fatalf("RIDE_ON hnd=0 debía emitir 0xCC03 con 0x306B")
	}

	rideOn := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_ON_REQ)
	rideOn.WriteUint16(77)
	rideOn.WriteByte(2)
	_ = zClientA.sendPacket(rideOn)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_MOVER_RIDE_ON_CMD) == nil {
		t.Fatalf("RIDE_ON válido debía emitir 0xCC02")
	}
	// Doble RIDE_ON estando montado -> 0xCC03 (0x306A)
	_ = zClientA.sendPacket(rideOn)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		f := f23FindPkt(pkts, protocol.PROTO_NC_MOVER_RIDE_ON_FAIL_CMD)
		if f == nil {
			return true
		}
		code, _ := f.Reader.ReadUint16()
		return code != 0x306A
	}() {
		t.Fatalf("RIDE_ON ya montado debía emitir 0xCC03 con 0x306A")
	}

	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_OFF_REQ))
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_MOVER_RIDE_OFF_CMD) == nil {
		t.Fatalf("RIDE_OFF montado debía emitir 0xCC06")
	}
	// RIDE_OFF desmontado -> 0xCC07 (0x3081)
	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_OFF_REQ))
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		f := f23FindPkt(pkts, protocol.PROTO_NC_MOVER_RIDE_OFF_FAIL_CMD)
		if f == nil {
			return true
		}
		code, _ := f.Reader.ReadUint16()
		return code != 0x3081
	}() {
		t.Fatalf("RIDE_OFF sin montura debía emitir 0xCC07 con 0x3081")
	}

	// 7) Interacción NPC: fuera de rango + inválido vs RouSmithJames (Weapon 0x3C09) y StoreManager (0x3C08)
	var smithNPC, storeNPC *zone.ZoneNPC
	for _, n := range pA.CurrentMap.GetAllNPCs() {
		if n.Name == "RouSmithJames" {
			smithNPC = n
		}
		if n.Role == "StoreManager" && storeNPC == nil {
			storeNPC = n
		}
	}
	if smithNPC == nil || storeNPC == nil {
		t.Fatalf("Faltan NPCs oficiales en RouN (smith=%v, store=%v)", smithNPC, storeNPC)
	}

	// Lejos (> 1000) -> no abre menú 0x201C
	pA.Char.XPos = 100
	pA.Char.YPos = 100
	farClick := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCCLICK_CMD)
	farClick.WriteUint16(smithNPC.Handle)
	_ = zClientA.sendPacket(farClick)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ) != nil {
		t.Fatalf("Click NPC fuera de rango (0xF4240) no debía emitir 0x201C")
	}

	// Cerca de RouSmithJames -> 0x201C -> ACK(1) -> 0x3C09 SHOPOPENTABLE_WEAPON_CMD
	pA.Char.XPos = int(smithNPC.X)
	pA.Char.YPos = int(smithNPC.Y)
	_ = zClientA.sendPacket(farClick)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ) == nil {
		t.Fatalf("Click cerca de RouSmithJames debía emitir 0x201C")
	}
	menuAck := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK)
	menuAck.WriteByte(1)
	_ = zClientA.sendPacket(menuAck)
	if pkts := f23CollectPackets(zConnA, 120*time.Millisecond); f23FindPkt(pkts, protocol.PROTO_NC_MENU_SHOPOPENTABLE_WEAPON_CMD) == nil {
		t.Fatalf("NPCMENUOPEN_ACK(1) en RouSmithJames debía emitir PROTO_NC_MENU_SHOPOPENTABLE_WEAPON_CMD (0x3C09)")
	}
	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ACT_ENDOFTRADE_CMD))

	// Cerca de StoreManager -> 0x3C08 OPENSTORAGE_CMD
	pA.Char.XPos = int(storeNPC.X)
	pA.Char.YPos = int(storeNPC.Y)
	storeClick := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCCLICK_CMD)
	storeClick.WriteUint16(storeNPC.Handle)
	_ = zClientA.sendPacket(storeClick)
	pktsA = f23CollectPackets(zConnA, 100*time.Millisecond)
	if storeNPC.Menu != 0 {
		if f23FindPkt(pktsA, protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ) == nil {
			t.Fatalf("Click en StoreManager (Menu=%d) debía emitir 0x201C", storeNPC.Menu)
		}
		_ = zClientA.sendPacket(menuAck)
		pktsA = f23CollectPackets(zConnA, 120*time.Millisecond)
	}
	if f23FindPkt(pktsA, protocol.PROTO_NC_MENU_OPENSTORAGE_CMD) == nil {
		t.Fatalf("StoreManager debía emitir PROTO_NC_MENU_OPENSTORAGE_CMD (0x3C08)")
	}
	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ACT_ENDOFTRADE_CMD))

	// 8) Portales y Links: TOWNPORTAL (muerto/inválido/ok), WING_FLY (!LinkOut/ok), MULTY_LINK_SELECT (0x2029/0x202B/LINKSAME/LINKOTHER)
	pA.State = 4 // Dead -> TOWNPORTAL 0x19C2 & WING_FLY 0x1146
	tpReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_TOWNPORTAL_REQ)
	tpReq.WriteByte(0)
	_ = zClientA.sendPacket(tpReq)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		p := f23FindPkt(pkts, protocol.PROTO_NC_MAP_TOWNPORTAL_ACK)
		if p == nil {
			return true
		}
		c, _ := p.Reader.ReadUint16()
		return c != 0x19C2
	}() {
		t.Fatalf("TOWNPORTAL en estado muerto debía emitir 0x181B con 0x19C2")
	}
	pA.State = 1

	tpBadIdx := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_TOWNPORTAL_REQ)
	tpBadIdx.WriteByte(9)
	_ = zClientA.sendPacket(tpBadIdx)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		p := f23FindPkt(pkts, protocol.PROTO_NC_MAP_TOWNPORTAL_ACK)
		if p == nil {
			return true
		}
		c, _ := p.Reader.ReadUint16()
		return c != 0x19C3
	}() {
		t.Fatalf("TOWNPORTAL índice inválido debía emitir 0x181B con 0x19C3")
	}

	// MULTY_LINK_SELECT inválido (mapa inexistente / FieldLvCondition) -> 0x1820 (0x2029)
	mlBad := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ)
	mlBad.WriteString("NoSuchMap99", 12)
	_ = zClientA.sendPacket(mlBad)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		p := f23FindPkt(pkts, protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK)
		if p == nil {
			return true
		}
		c, _ := p.Reader.ReadUint16()
		return c != 0x2029
	}() {
		t.Fatalf("MULTY_LINK mapa inválido debía emitir 0x1820 con 0x2029")
	}

	// Registrar MultiLink Gate temporal para probar distancia (> 0xF4240 -> 0x202B)
	mlNPC := &zone.ZoneNPC{MobID: 1, Name: "TestMultiGate", X: 6888, Y: 7045, IsGate: true, GateToWhere: "Rou,RouN"}
	pA.CurrentMap.AddNPC(mlNPC)
	pA.Char.XPos = 6888
	pA.Char.YPos = 7045
	mlClick := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCCLICK_CMD)
	mlClick.WriteUint16(mlNPC.Handle)
	_ = zClientA.sendPacket(mlClick)
	_ = f23CollectPackets(zConnA, 100*time.Millisecond)
	pA.Char.XPos = 100 // alejar > 1000 unidades
	pA.Char.YPos = 100
	mlSame := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ)
	mlSame.WriteString("Rou", 12)
	_ = zClientA.sendPacket(mlSame)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		p := f23FindPkt(pkts, protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK)
		if p == nil {
			return true
		}
		c, _ := p.Reader.ReadUint16()
		return c != 0x202B
	}() {
		t.Fatalf("MULTY_LINK fuera de rango (>0xF4240) debía emitir 0x1820 con 0x202B")
	}

	// LINKSAME válido ("Rou", mismo zoneNo=1)
	_ = zClientA.sendPacket(mlSame)
	pktsA = f23CollectPackets(zConnA, 180*time.Millisecond)
	if f23FindPkt(pktsA, protocol.PROTO_NC_MAP_LINKSAME_CMD) == nil {
		t.Fatalf("MULTY_LINK a 'Rou' debía emitir PROTO_NC_MAP_LINKSAME_CMD (0x1809)")
	}
	_ = zClientA.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	_ = f23CollectPackets(zConnA, 120*time.Millisecond)

	// WING_FLY bloqueado (!LinkOut en 'EldPri') -> 0x1817 (0x1143)
	prevMap := pA.CurrentMap
	pA.CurrentMap = zoneSrv1.Handler().MapManager().GetOrCreateMap("EldPri")
	wfReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_WING_FLY_REQ)
	wfReq.WriteByte(0)
	_ = zClientA.sendPacket(wfReq)
	if pkts := f23CollectPackets(zConnA, 100*time.Millisecond); func() bool {
		p := f23FindPkt(pkts, protocol.PROTO_NC_MAP_WING_FLY_ACK)
		if p == nil {
			return true
		}
		c, _ := p.Reader.ReadUint16()
		return c != 0x1143
	}() {
		t.Fatalf("WING_FLY en mapa !LinkOut ('EldPri') debía emitir 0x1817 con 0x1143")
	}
	pA.CurrentMap = prevMap

	// LINKOTHER inter-zona real ('RouCos02' tiene fiesta=2 -> Zone2): verificar 0x180A + 30B payload
	mlOther := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ)
	mlOther.WriteString("RouCos02", 12)
	_ = zClientA.sendPacket(mlOther)
	pktsA = f23CollectPackets(zConnA, 200*time.Millisecond)
	linkOther := f23FindPkt(pktsA, protocol.PROTO_NC_MAP_LINKOTHER_CMD)
	if linkOther == nil {
		t.Fatalf("MULTY_LINK hacia 'RouCos02' (Zone 2) debía emitir PROTO_NC_MAP_LINKOTHER_CMD (0x180A)")
	}
	if len(linkOther.Raw)-2 != 30 {
		t.Fatalf("PROTO_NC_MAP_LINKOTHER_CMD debía medir 30B payload, obtenido=%d", len(linkOther.Raw)-2)
	}
	_ = linkOther.Reader.Skip(10) // mapSerial(2) + x(4) + y(4)
	dstIP, _ := linkOther.Reader.ReadString(16)
	dstPort, _ := linkOther.Reader.ReadUint16()
	if dstIP != "127.0.0.1" || int(dstPort) != zonePort2 {
		t.Fatalf("LINKOTHER endpoint esperado 127.0.0.1:%d, obtenido %s:%d", zonePort2, dstIP, dstPort)
	}

	// Verificar GameLog registrado para charIDA (e.g. TypeG de MOVER / TypeL de GM)
	time.Sleep(150 * time.Millisecond)
	var logCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM tGameLog WHERE nCharNo = ?", charIDA).Scan(&logCount)
	if logCount == 0 {
		t.Fatalf("Se esperaban registros en tGameLog para char %d", charIDA)
	}
}
