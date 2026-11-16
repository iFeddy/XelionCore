package audit_test

import (
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

type f25Cluster struct {
	t           *testing.T
	db          *sql.DB
	accountRepo *database.MySQLAccountRepository
	charRepo    database.CharacterRepository
	ds        *zone.DataStore
	loginAddr string
	worldAddr string
	zoneAddr  string
	verKey    string
	zoneSrv   *zone.Server
	glDisp    *database.GameLogDispatcher
}

type f25PlayerConn struct {
	userNo   int
	userName string
	charName string
	charID   int
	conn     net.Conn
	client   *auditMockClient
	zp       *zone.ZonePlayer
}

func f25PackInven(box uint8, slot uint8) uint16 {
	return (uint16(box) << 10) | uint16(slot)
}

func f25ReadU16(p *protocol.Packet) uint16 {
	if p == nil || len(p.Raw) < 4 {
		return 0xFFFF
	}
	return binary.LittleEndian.Uint16(p.Raw[2:4])
}

func f25FindErr(pkts []*protocol.Packet, op uint16) (uint16, bool) {
	errOffset := 2 // default payload offset 0 -> Raw[2:4]
	switch op {
	case protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK:
		errOffset = 4 // uint16 ProduceID + uint16 ErrorCode
	case protocol.PROTO_NC_ITEM_PICK_ACK:
		errOffset = 8 // uint16 Itemid + uint32 Lot + uint16 Error
	case protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK:
		errOffset = 6 // uint32 OrderNo + uint16 Err
	case protocol.PROTO_NC_ITEM_DEPOSIT_ACK,
		protocol.PROTO_NC_ITEM_WITHDRAW_ACK,
		protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK,
		protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK,
		protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK:
		errOffset = 10 // uint64 Cen + uint16 Err
	}
	for _, p := range pkts {
		if p != nil && p.OpCode == op && len(p.Raw) >= errOffset+2 {
			return binary.LittleEndian.Uint16(p.Raw[errOffset : errOffset+2]), true
		}
	}
	return 0, false
}

func f25ReadPacketUntil(conn net.Conn, deadline time.Time) (*protocol.Packet, error) {
	_ = conn.SetReadDeadline(deadline)
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}
	payloadLen := int(lenBuf[0])
	if payloadLen == 0 {
		extBuf := make([]byte, 2)
		if _, err := io.ReadFull(conn, extBuf); err != nil {
			return nil, err
		}
		payloadLen = int(binary.LittleEndian.Uint16(extBuf))
	}
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return protocol.NewPacket(payload), nil
}

func f25WaitOp(conn net.Conn, waitOps ...uint16) []*protocol.Packet {
	var pkts []*protocol.Packet
	want := make(map[uint16]bool, len(waitOps))
	for _, op := range waitOps {
		want[op] = true
	}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		p, err := f25ReadPacketUntil(conn, deadline)
		if err != nil {
			break
		}
		if p != nil {
			pkts = append(pkts, p)
			if want[p.OpCode] {
				if tail := f23CollectPackets(conn, 25*time.Millisecond); len(tail) > 0 {
					pkts = append(pkts, tail...)
				}
				return pkts
			}
		}
	}
	return pkts
}

func f25CleanupUser(db *sql.DB, userNo int, charID int) {
	_, _ = db.Exec("DELETE FROM tChargeItem WHERE userNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tLogout WHERE nUserNo = ?", userNo)
	if charID > 0 {
		_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", charID)
		_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", charID)
		_, _ = db.Exec("DELETE FROM tItemOption WHERE nItemKey IN (SELECT nItemKey FROM tItem WHERE nOwner = ?)", charID)
		_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", charID)
		_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", charID)
	}
	_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", userNo)
}

func f25CreateAndEnterPlayer(t *testing.T, cl *f25Cluster, userNo int, suffix string) *f25PlayerConn {
	t.Helper()
	ts := time.Now().Unix() % 10000
	userName := fmt.Sprintf("f25u_%s_%d", suffix, ts)
	charName := fmt.Sprintf("F25%s%d", suffix, ts)
	if len(charName) > 15 {
		charName = charName[:15]
	}
	pass := "f25pass"
	sum := md5.Sum([]byte(pass))
	passHash := hex.EncodeToString(sum[:])

	f25CleanupUser(cl.db, userNo, 0)
	if _, err := cl.db.Exec(`INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate) VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())`, userNo, userName, passHash); err != nil {
		t.Fatalf("Insert tUser %d: %v", userNo, err)
	}

	worldToken := f23LoginAndGetWorldToken(t, cl.loginAddr, cl.verKey, userName, passHash)
	wConn, wClient := f23DialAndHandshake(t, cl.worldAddr)
	defer wConn.Close()

	wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq.WriteString(userName, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD %s failed: %v / %v", suffix, wlAck, err)
	}

	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(charName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE %s failed: %v / %v", suffix, createAck, err)
	}

	ch, err := cl.charRepo.GetCharacterByName(charName)
	if err != nil || ch == nil {
		t.Fatalf("GetCharacterByName(%s): %v", charName, err)
	}

	cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	cSelReq.WriteByte(0)
	_ = wClient.sendPacket(cSelReq)
	if cSelAck, err := auditReadServerPacket(wConn); err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CHAR_LOGIN %s failed: %v / %v", suffix, cSelAck, err)
	}
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)

	var dbRandomID uint16
	if err := cl.db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", charName).Scan(&dbRandomID); err != nil {
		t.Fatalf("tActiveZoneTransfer missing for %s: %v", charName, err)
	}

	zConn, zClient, _ := f23EnterZone(t, cl.zoneAddr, cl.ds, dbRandomID, charName)
	zp, ok := cl.zoneSrv.Handler().GetPlayerByName(charName)
	if !ok || zp == nil || zp.CurrentMap == nil {
		_ = zConn.Close()
		t.Fatalf("Jugador '%s' no encontrado en ZoneServer", charName)
	}
	return &f25PlayerConn{
		userNo:   userNo,
		userName: userName,
		charName: charName,
		charID:   ch.ID,
		conn:     zConn,
		client:   zClient,
		zp:       zp,
	}
}

func f25SetPos(p *f25PlayerConn, x, y int) {
	if p != nil && p.zp != nil && p.zp.Char != nil {
		p.zp.Char.XPos = x
		p.zp.Char.YPos = y
	}
}

func f25SeedItem(t *testing.T, repo database.CharacterRepository, charID int, box uint8, slot int16, itemID uint16, lot int) *database.Item {
	t.Helper()
	it := &database.Item{
		Owner:       charID,
		ItemID:      itemID,
		StorageType: box,
		Storage:     slot,
		Flags:       int32(lot),
	}
	if err := repo.CreateItem(it); err != nil {
		t.Fatalf("CreateItem(char=%d, box=%d, slot=%d, id=%d): %v", charID, box, slot, itemID, err)
	}
	return it
}

func f25ClearInventory(t *testing.T, repo database.CharacterRepository, charID int, box uint8) {
	t.Helper()
	items, err := repo.GetCharacterItems(charID)
	if err != nil {
		t.Fatalf("GetCharacterItems: %v", err)
	}
	for _, it := range items {
		if it != nil && it.StorageType == box {
			_ = repo.DeleteItem(it.ItemKey)
		}
	}
}

type f25MatNeed struct {
	id  uint16
	qty int
}

type f25SampleItems struct {
	lowWeapon      *item.ItemInfo
	highWeapon     *item.ItemInfo
	bracelet       *item.ItemInfo
	stackPotion    *item.ItemInfo
	upgradeStone   *item.ItemInfo
	luckyProtect   *item.ItemInfo
	reviveScroll   *item.ItemInfo
	skillScroll    *item.ItemInfo
	activeSkillID  uint16
	produceID      uint16
	produceMats    []f25MatNeed
	mixMatA        uint16
	mixMatB        uint16
	noTradeItem    *item.ItemInfo
}

func f25PickCatalogItems(t *testing.T, cl *f25Cluster) f25SampleItems {
	t.Helper()
	var s f25SampleItems
	noTradeByInx := make(map[string]bool)
	belongByBT := make(map[int]bool)
	if btTbl := cl.ds.GetSHN("BelongTypeInfo"); btTbl != nil {
		for _, r := range btTbl.Rows {
			if r.GetInt("NoTrade") != 0 {
				belongByBT[r.GetInt("BT_Inx")] = true
			}
		}
	}
	if iiTbl := cl.ds.GetSHN("ItemInfo"); iiTbl != nil && len(belongByBT) > 0 {
		for _, r := range iiTbl.Rows {
			bt := r.GetInt("BT_Inx")
			if bt != 0 && belongByBT[bt] {
				noTradeByInx[strings.ToLower(strings.TrimSpace(r.GetString("InxName")))] = true
			}
		}
	}
	for id := uint16(1); id < 60000; id++ {
		info, ok := cl.ds.ItemStore.GetItemInfo(id)
		if !ok || info == nil {
			continue
		}
		low := strings.ToLower(info.InxName)
		if s.lowWeapon == nil && (info.Equip == 11 || info.Equip == 10) && info.DemandLv <= 5 && info.UpLimit >= 9 && (info.UseClass&(1<<1)) != 0 {
			if can, know := cl.ds.ItemStore.CanUseItemByClassID(id, 1); !know || can {
				s.lowWeapon = info
			}
		}
		if s.highWeapon == nil && info.Class == 5 && info.DemandLv >= 60 {
			s.highWeapon = info
		}
		if s.bracelet == nil && (strings.Contains(low, "bracelet") || strings.Contains(low, "bangle")) && info.Equip != 0 {
			s.bracelet = info
		}
		if s.stackPotion == nil && (info.Class == 0 || info.Class == 1 || info.Class == 2) && info.MaxLot >= 10 && info.BuyPrice > 0 {
			s.stackPotion = info
		}
		if s.upgradeStone == nil && (strings.Contains(low, "elrue") || strings.Contains(low, "lix")) && info.MaxLot >= 1 {
			s.upgradeStone = info
		}
		if s.luckyProtect == nil && (strings.Contains(low, "lucky") || strings.Contains(low, "karis")) {
			s.luckyProtect = info
		}
		if s.reviveScroll == nil && (strings.Contains(low, "revive") || strings.Contains(low, "resurrect")) {
			s.reviveScroll = info
		}
		if s.skillScroll == nil && strings.Contains(low, "scroll") && info.MaxLot >= 1 {
			s.skillScroll = info
		}
		if s.noTradeItem == nil && noTradeByInx[low] {
			s.noTradeItem = info
		}
	}
	if tbl := cl.ds.GetSHN("ItemMixInfo"); tbl != nil {
		for _, row := range tbl.Rows {
			m1 := row.GetUint16("Material01")
			m2 := row.GetUint16("Material02")
			if m1 > 0 && m2 > 0 {
				s.mixMatA, s.mixMatB = m1, m2
				break
			}
		}
	}
	if tbl := cl.ds.GetSHN("Produce"); tbl != nil {
		for _, row := range tbl.Rows {
			pid := row.GetUint16("ProductID")
			prodName := strings.TrimSpace(row.GetString("Product"))
			if pid == 0 || prodName == "" {
				continue
			}
			if _, ok := cl.ds.ItemStore.GetItemInfoByName(prodName); !ok {
				continue
			}
			var mats []f25MatNeed
			valid := true
			for i := 0; i < 8; i++ {
				rn := strings.TrimSpace(row.GetString(fmt.Sprintf("Raw%d", i)))
				q := row.GetInt(fmt.Sprintf("Quantity%d", i))
				if q <= 0 || rn == "" || rn == "-" {
					continue
				}
				if ri, ok := cl.ds.ItemStore.GetItemInfoByName(rn); ok && ri != nil {
					mats = append(mats, f25MatNeed{id: ri.ID, qty: q})
				} else {
					valid = false
					break
				}
			}
			if valid && len(mats) > 0 {
				s.produceID = pid
				s.produceMats = mats
				break
			}
		}
	}
	if s.lowWeapon == nil || s.stackPotion == nil {
		t.Fatalf("No se encontraron ítems base necesarios en ItemInfo.shn (lowWeapon=%v stackPotion=%v)", s.lowWeapon, s.stackPotion)
	}
	if s.upgradeStone == nil {
		s.upgradeStone = s.stackPotion
	}
	if s.luckyProtect == nil {
		s.luckyProtect = s.upgradeStone
	}
	if s.bracelet == nil {
		s.bracelet = s.lowWeapon
	}
	if s.skillScroll == nil {
		s.skillScroll = s.stackPotion
	}
	if s.activeSkillID == 0 {
		s.activeSkillID = s.skillScroll.ID
	}
	return s
}

// TestFase25_ItemEconomyForgeTradeBoothBranches_ClusterE2E audita rama por rama contra Zone.exe
// todos los códigos de error oficiales, cambios de celda, actualizaciones de Cen, protección de forja,
// SafeTrade, Booth y GameLog para NC_ITEM, NC_ITEMDB, NC_BOOTH, NC_CHARSAVE y NC_ACT_PRODUCE.
func TestFase25_ItemEconomyForgeTradeBoothBranches_ClusterE2E(t *testing.T) {
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
		loginPort = 19761
		worldPort = 19765
		zonePort  = 19775
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
		Name:     "Fase25Cluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort,
	}, charRepo, transfers, worldDataStore, accountRepo)
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

	cl := &f25Cluster{
		t:           t,
		db:          db,
		accountRepo: accountRepo,
		charRepo:    charRepo,
		ds:        zoneDataStore,
		loginAddr: fmt.Sprintf("127.0.0.1:%d", loginPort),
		worldAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		zoneAddr:  fmt.Sprintf("127.0.0.1:%d", zonePort),
		verKey:    activeVerKey,
		zoneSrv:   zoneSrv,
		glDisp:    glDispatcher,
	}

	pA := f25CreateAndEnterPlayer(t, cl, 9989251, "A")
	defer func() {
		_ = pA.conn.Close()
		f25CleanupUser(db, pA.userNo, pA.charID)
	}()
	pB := f25CreateAndEnterPlayer(t, cl, 9989252, "B")
	defer func() {
		_ = pB.conn.Close()
		f25CleanupUser(db, pB.userNo, pB.charID)
	}()

	samples := f25PickCatalogItems(t, cl)
	_ = f23CollectPackets(pA.conn, 60*time.Millisecond)
	_ = f23CollectPackets(pB.conn, 60*time.Millisecond)

	runFase25InventoryBranches(t, cl, pA, pB, samples)
	runFase25EconomyStorageBranches(t, cl, pA, samples)
	runFase25ForgeCraftBranches(t, cl, pA, samples)
	runFase25UtilitiesAndTradeBoothBranches(t, cl, pA, pB, samples)

	// Verificar GameLog en BD tras operaciones de economía/forja/trade
	glDispatcher.Close(t.Context())
	var logCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tGameLog WHERE nCharNo = ?", pA.charID).Scan(&logCount); err != nil {
		t.Fatalf("Query tGameLog: %v", err)
	}
	if logCount == 0 {
		t.Fatalf("Se esperaban registros GameLog para nCharNo=%d tras drop/buy/charged/upgrade/trade", pA.charID)
	}
}

func runFase25InventoryBranches(t *testing.T, cl *f25Cluster, pA, pB *f25PlayerConn, s f25SampleItems) {
	t.Helper()
	f25ClearInventory(t, cl.charRepo, pA.charID, 9)
	f25ClearInventory(t, cl.charRepo, pA.charID, 8)

	// 1. RELOC_REQ: estado muerto (0x0242), destino inválido Box 2 (0x0243), éxito (0x0241)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 0, s.lowWeapon.ID, 1)

	pA.zp.State = 3
	relocDead := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocDead.WriteUint16(f25PackInven(9, 0))
	relocDead.WriteUint16(f25PackInven(9, 1))
	_ = pA.client.sendPacket(relocDead)
	if errCode, ok := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RELOC_ACK), protocol.PROTO_NC_ITEM_RELOC_ACK); !ok || errCode != 0x0242 {
		t.Fatalf("RELOC estado muerto esperaba 0x0242, got ok=%v err=0x%04X", ok, errCode)
	}
	pA.zp.State = 1

	relocBox2 := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocBox2.WriteUint16(f25PackInven(9, 0))
	relocBox2.WriteUint16(f25PackInven(2, 0))
	_ = pA.client.sendPacket(relocBox2)
	if errCode, ok := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RELOC_ACK), protocol.PROTO_NC_ITEM_RELOC_ACK); !ok || errCode != 0x0243 {
		t.Fatalf("RELOC hacia Box 2 esperaba 0x0243, got ok=%v err=0x%04X", ok, errCode)
	}

	relocOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	relocOk.WriteUint16(f25PackInven(9, 0))
	relocOk.WriteUint16(f25PackInven(9, 3))
	_ = pA.client.sendPacket(relocOk)
	pkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RELOC_ACK)
	if errCode, ok := f25FindErr(pkts, protocol.PROTO_NC_ITEM_RELOC_ACK); !ok || errCode != 0x0241 || f23FindPkt(pkts, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD) == nil {
		t.Fatalf("RELOC válido esperaba 0x0241 + CELLCHANGE_CMD, got ok=%v err=0x%04X", ok, errCode)
	}

	// 2. SPLIT_REQ: no Box 9 (0x0443), mismo slot (0x0445), faltante (0x0448), ocupado (0x0446), lot inválido (0x0449), éxito (0x0441)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 4, s.stackPotion.ID, 5)
	sendSplit := func(fromBox, fromSlot, toBox, toSlot uint8, lot uint32) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_REQ)
		pw.WriteUint16(f25PackInven(fromBox, fromSlot))
		pw.WriteUint16(f25PackInven(toBox, toSlot))
		pw.WriteUint32(lot)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_SPLIT_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_SPLIT_ACK)
		return e, ps
	}
	if e, _ := sendSplit(12, 4, 9, 5, 1); e != 0x0443 {
		t.Fatalf("SPLIT no Box9 esperaba 0x0443, got 0x%04X", e)
	}
	if e, _ := sendSplit(9, 4, 9, 4, 1); e != 0x0445 {
		t.Fatalf("SPLIT mismo slot esperaba 0x0445, got 0x%04X", e)
	}
	if e, _ := sendSplit(9, 25, 9, 5, 1); e != 0x0448 {
		t.Fatalf("SPLIT origen faltante esperaba 0x0448, got 0x%04X", e)
	}
	if e, _ := sendSplit(9, 4, 9, 3, 1); e != 0x0446 {
		t.Fatalf("SPLIT destino ocupado esperaba 0x0446, got 0x%04X", e)
	}
	if e, _ := sendSplit(9, 4, 9, 5, 99); e != 0x0449 {
		t.Fatalf("SPLIT lot>=actual esperaba 0x0449, got 0x%04X", e)
	}
	if e, ps := sendSplit(9, 4, 9, 5, 2); e != 0x0441 || f23FindPkt(ps, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD) == nil {
		t.Fatalf("SPLIT válido esperaba 0x0441 + CELLCHANGE_CMD, got 0x%04X", e)
	}

	// 3. EQUIP / RINGEQUIP / UNEQUIP
	eqMissing := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_REQ)
	eqMissing.WriteByte(40)
	_ = pA.client.sendPacket(eqMissing)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_EQUIP_ACK), protocol.PROTO_NC_ITEM_EQUIP_ACK); e != 0x0283 {
		t.Fatalf("EQUIP slot vacío esperaba 0x0283, got 0x%04X", e)
	}

	if s.highWeapon != nil {
		f25SeedItem(t, cl.charRepo, pA.charID, 9, 6, s.highWeapon.ID, 1)
		eqHigh := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_REQ)
		eqHigh.WriteByte(6)
		_ = pA.client.sendPacket(eqHigh)
		if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_EQUIP_ACK), protocol.PROTO_NC_ITEM_EQUIP_ACK); e != 0x0285 {
			t.Fatalf("EQUIP nivel insuficiente esperaba 0x0285, got 0x%04X", e)
		}
	}

	ringBad := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RINGEQUIP_REQ)
	ringBad.WriteByte(3)
	ringBad.WriteByte(0x05)
	_ = pA.client.sendPacket(ringBad)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_EQUIP_ACK), protocol.PROTO_NC_ITEM_EQUIP_ACK); e != 0x0283 {
		t.Fatalf("RINGEQUIP slot inválido esperaba 0x0283, got 0x%04X", e)
	}

	eqOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_REQ)
	eqOk.WriteByte(3)
	_ = pA.client.sendPacket(eqOk)
	eqPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_EQUIP_ACK)
	if e, ok := f25FindErr(eqPkts, protocol.PROTO_NC_ITEM_EQUIP_ACK); !ok || e != 0x0281 || f23FindPkt(eqPkts, protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD) == nil {
		t.Fatalf("EQUIP válido esperaba 0x0281 + EQUIPCHANGE_CMD, got ok=%v err=0x%04X", ok, e)
	}

	// Unequip slot vacío -> 0x02C2; válido -> 0x3011 (0x0281)
	uneqMissing := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_UNEQUIP_REQ)
	uneqMissing.WriteByte(0x1E)
	uneqMissing.WriteByte(7)
	_ = pA.client.sendPacket(uneqMissing)
	if e, ok := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_UNEQUIP_ACK), protocol.PROTO_NC_ITEM_UNEQUIP_ACK); !ok || e != 0x02C2 {
		t.Fatalf("UNEQUIP slot vacío esperaba UNEQUIP_ACK=0x02C2, got ok=%v err=0x%04X", ok, e)
	}

	equipSlot := byte(s.lowWeapon.Equip)
	if items, _ := cl.charRepo.GetCharacterItems(pA.charID); len(items) > 0 {
		for _, it := range items {
			if it != nil && it.StorageType == 8 && it.ItemID == s.lowWeapon.ID {
				equipSlot = byte(it.Storage)
				break
			}
		}
	}
	uneqOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_UNEQUIP_REQ)
	uneqOk.WriteByte(equipSlot)
	uneqOk.WriteByte(7)
	_ = pA.client.sendPacket(uneqOk)
	if e, ok := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_EQUIP_ACK), protocol.PROTO_NC_ITEM_EQUIP_ACK); !ok || e != 0x0281 {
		t.Fatalf("UNEQUIP válido esperaba EQUIP_ACK=0x0281, got ok=%v err=0x%04X", ok, e)
	}

	// 4. DROP_REQ / PICK_REQ: no Box 9 (0x030D), estado vendor (0x030F), faltante (0x0303), éxito (0x0301)
	sendDrop := func(box, slot uint8, lot uint32) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_REQ)
		pw.WriteUint16(f25PackInven(box, slot))
		pw.WriteUint32(lot)
		pw.WriteUint32(uint32(pA.zp.Char.XPos))
		pw.WriteUint32(uint32(pA.zp.Char.YPos))
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_DROP_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_DROP_ACK)
		return e, ps
	}
	if e, _ := sendDrop(8, 1, 1); e != 0x030D {
		t.Fatalf("DROP no Box9 esperaba 0x030D, got 0x%04X", e)
	}
	pA.zp.State = 3
	if e, _ := sendDrop(9, 7, 1); e != 0x0302 {
		t.Fatalf("DROP muerto esperaba 0x0302, got 0x%04X", e)
	}
	pA.zp.State = 1
	pA.zp.ActiveBooth = &zone.PlayerBooth{}
	if e, _ := sendDrop(9, 7, 1); e != 0x030F {
		t.Fatalf("DROP en booth/vendor esperaba 0x030F, got 0x%04X", e)
	}
	pA.zp.ActiveBooth = nil
	if e, _ := sendDrop(9, 45, 1); e != 0x0303 {
		t.Fatalf("DROP ítem inexistente esperaba 0x0303, got 0x%04X", e)
	}
	if e, _ := sendDrop(9, 7, 1); e != 0x0301 {
		t.Fatalf("DROP válido esperaba 0x0301, got 0x%04X", e)
	}

	// PICK_REQ: inexistente (0x034B), lejos >=125u (0x0344), dueño ajeno (0x0345), lleno (0x0346), Cen cuando inventario lleno (0x0341)
	pickMissing := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pickMissing.WriteUint16(64900)
	_ = pA.client.sendPacket(pickMissing)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_PICK_ACK), protocol.PROTO_NC_ITEM_PICK_ACK); e != 0x034B {
		t.Fatalf("PICK handle inexistente esperaba 0x034B, got 0x%04X", e)
	}

	px, py := pA.zp.GetPosition()
	ownedDrop := pA.zp.CurrentMap.SpawnMonsterDrop(s.lowWeapon.ID, 1, 0, int32(px+5), int32(py+5), 0, pA.charID, 30*time.Second, false, false, nil)
	f25SetPos(pB, px+5, py+5)
	pickOwned := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pickOwned.WriteUint16(ownedDrop.Handle)
	_ = pB.client.sendPacket(pickOwned)
	if e, _ := f25FindErr(f25WaitOp(pB.conn, protocol.PROTO_NC_ITEM_PICK_ACK), protocol.PROTO_NC_ITEM_PICK_ACK); e != 0x0345 {
		t.Fatalf("PICK con dueño activo esperaba 0x0345, got 0x%04X", e)
	}

	f25SetPos(pA, px+250, py)
	pickFar := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pickFar.WriteUint16(ownedDrop.Handle)
	_ = pA.client.sendPacket(pickFar)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_PICK_ACK), protocol.PROTO_NC_ITEM_PICK_ACK); e != 0x0344 {
		t.Fatalf("PICK lejos (>=125u) esperaba 0x0344, got 0x%04X", e)
	}
	f25SetPos(pA, px, py)

	// Llenar las 96 celdas de Box 9 de pA -> PICK ítem debe responder 0x0346; PICK Cen debe dar 0x0341
	f25ClearInventory(t, cl.charRepo, pA.charID, 9)
	for slot := int16(0); slot < 96; slot++ {
		f25SeedItem(t, cl.charRepo, pA.charID, 9, slot, s.lowWeapon.ID, 1)
	}
	pickFull := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pickFull.WriteUint16(ownedDrop.Handle)
	_ = pA.client.sendPacket(pickFull)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_PICK_ACK), protocol.PROTO_NC_ITEM_PICK_ACK); e != 0x0346 {
		t.Fatalf("PICK con inventario lleno esperaba 0x0346, got 0x%04X", e)
	}

	f25SeedItem(t, cl.charRepo, pA.charID, 8, int16(equipSlot), s.lowWeapon.ID, 1)
	uneqFull := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_UNEQUIP_REQ)
	uneqFull.WriteByte(equipSlot)
	uneqFull.WriteByte(0)
	_ = pA.client.sendPacket(uneqFull)
	if e, ok := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_UNEQUIP_ACK), protocol.PROTO_NC_ITEM_UNEQUIP_ACK); !ok || e != 0x02C4 {
		t.Fatalf("UNEQUIP con inventario lleno esperaba UNEQUIP_ACK=0x02C4, got ok=%v err=0x%04X", ok, e)
	}
	f25ClearInventory(t, cl.charRepo, pA.charID, 8)

	cenDrop := pA.zp.CurrentMap.SpawnMonsterDrop(0, 777, 0, int32(px+5), int32(py+5), 0, pA.charID, 30*time.Second, false, true, nil)
	moneyBefore := pA.zp.Char.TCharacterState.Money
	pickCen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pickCen.WriteUint16(cenDrop.Handle)
	_ = pA.client.sendPacket(pickCen)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_PICK_ACK), protocol.PROTO_NC_ITEM_PICK_ACK); e != 0x0341 || pA.zp.Char.TCharacterState.Money != moneyBefore+777 {
		t.Fatalf("PICK Cen con inventario lleno esperaba 0x0341 y +777 Cen, got err=0x%04X money=%d->%d", e, moneyBefore, pA.zp.Char.TCharacterState.Money)
	}

	// 5. AUTO_ARRANGE_INVEN_REQ: muerto -> 0x30AA, trade -> 0x30AC, ok -> 0x30A8
	pA.zp.State = 3
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK), protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK); e != 0x30AA {
		t.Fatalf("AUTO_ARRANGE muerto esperaba 0x30AA, got 0x%04X", e)
	}
	pA.zp.State = 1
	pA.zp.ActiveTrade = &zone.SafeTradeSession{}
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK), protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK); e != 0x30AC {
		t.Fatalf("AUTO_ARRANGE en trade esperaba 0x30AC, got 0x%04X", e)
	}
	pA.zp.ActiveTrade = nil
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK), protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK); e != 0x30A8 {
		t.Fatalf("AUTO_ARRANGE válido esperaba 0x30A8, got 0x%04X", e)
	}
}

func runFase25EconomyStorageBranches(t *testing.T, cl *f25Cluster, pA *f25PlayerConn, s f25SampleItems) {
	t.Helper()
	// 1. BUY_REQ: inventario lleno (0x0203), Cen insuficiente (0x020E), estado muerto (0x0204), compra válida (0x0201)
	pA.zp.Char.TCharacterState.Money = 1_000_000
	sendBuy := func(id uint16, lot uint32) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_REQ)
		pw.WriteUint16(id)
		pw.WriteUint32(lot)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_BUY_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_BUY_ACK)
		return e, ps
	}
	// Nota: pA tiene 96 slots ocupados del paso anterior!
	if e, _ := sendBuy(s.stackPotion.ID, 1); e != 0x0203 {
		t.Fatalf("BUY con inventario lleno esperaba 0x0203, got 0x%04X", e)
	}
	f25ClearInventory(t, cl.charRepo, pA.charID, 9)

	pA.zp.LastSaveAt = time.Now()
	pA.zp.Char.TCharacterState.Money = 0
	if e, _ := sendBuy(s.stackPotion.ID, 5); e != 0x020E {
		t.Fatalf("BUY sin Cen esperaba 0x020E, got 0x%04X", e)
	}
	pA.zp.State = 3
	if e, _ := sendBuy(s.stackPotion.ID, 1); e != 0x0204 {
		t.Fatalf("BUY muerto esperaba 0x0204, got 0x%04X", e)
	}
	pA.zp.State = 1
	pA.zp.Char.TCharacterState.Money = 500_000
	_ = cl.charRepo.SaveCharMoney(pA.charID, 500_000)
	if e, ps := sendBuy(s.stackPotion.ID, 5); e != 0x0201 || f23FindPkt(ps, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD) == nil {
		t.Fatalf("BUY válido esperaba 0x0201 + CELLCHANGE_CMD, got 0x%04X", e)
	}

	// 2. SELL_REQ: estado (0x0382), faltante (0x0383), válido (0x0381 + SELL_ITEM_INSERT_CMD)
	items, _ := cl.charRepo.GetCharacterItems(pA.charID)
	var buySlot byte
	for _, it := range items {
		if it != nil && it.StorageType == 9 && it.ItemID == s.stackPotion.ID {
			buySlot = byte(it.Storage)
			break
		}
	}
	sendSell := func(slot uint8, lot uint32) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SELL_REQ)
		pw.WriteByte(slot)
		pw.WriteUint32(lot)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_SELL_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_SELL_ACK)
		return e, ps
	}
	pA.zp.State = 3
	if e, _ := sendSell(buySlot, 1); e != 0x0382 {
		t.Fatalf("SELL muerto esperaba 0x0382, got 0x%04X", e)
	}
	pA.zp.State = 1
	if e, _ := sendSell(80, 1); e != 0x0383 {
		t.Fatalf("SELL slot vacío esperaba 0x0383, got 0x%04X", e)
	}
	if e, ps := sendSell(buySlot, 5); e != 0x0381 || f23FindPkt(ps, protocol.PROTO_NC_ITEM_SELL_ITEM_INSERT_CMD) == nil {
		t.Fatalf("SELL válido esperaba 0x0381 + INSERT_CMD, got 0x%04X", e)
	}

	// 3. REPURCHASE_REQ: faltante (0x366A), sin Cen (0x366B), válido (0x3668 + SELL_ITEM_DELETE_CMD)
	repMiss := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REPURCHASE_REQ)
	repMiss.WriteUint16(64999)
	_ = pA.client.sendPacket(repMiss)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_REPURCHASE_ACK), protocol.PROTO_NC_ITEM_REPURCHASE_ACK); e != 0x366A {
		t.Fatalf("REPURCHASE handle inexistente esperaba 0x366A, got 0x%04X", e)
	}

	if len(pA.zp.RepurchaseItems) > 0 {
		hnd := pA.zp.RepurchaseItems[0].Handle
		savedMoney := pA.zp.Char.TCharacterState.Money
		pA.zp.LastSaveAt = time.Now()
		pA.zp.Char.TCharacterState.Money = 0
		repNoCen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REPURCHASE_REQ)
		repNoCen.WriteUint16(hnd)
		_ = pA.client.sendPacket(repNoCen)
		if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_REPURCHASE_ACK), protocol.PROTO_NC_ITEM_REPURCHASE_ACK); e != 0x366B {
			t.Fatalf("REPURCHASE sin Cen esperaba 0x366B, got 0x%04X", e)
		}
		pA.zp.Char.TCharacterState.Money = savedMoney

		repOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REPURCHASE_REQ)
		repOk.WriteUint16(hnd)
		_ = pA.client.sendPacket(repOk)
		repPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_REPURCHASE_ACK)
		if e, ok := f25FindErr(repPkts, protocol.PROTO_NC_ITEM_REPURCHASE_ACK); !ok || e != 0x3668 || f23FindPkt(repPkts, protocol.PROTO_NC_ITEM_SELL_ITEM_DELETE_CMD) == nil {
			t.Fatalf("REPURCHASE válido esperaba 0x3668 + DELETE_CMD, got ok=%v err=0x%04X", ok, e)
		}
	}

	// 4. STORAGE DEPOSIT / WITHDRAW + GUILD + ACADEMY
	pA.zp.State = 3
	depDead := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DEPOSIT_REQ)
	depDead.WriteUint64(100)
	_ = pA.client.sendPacket(depDead)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_DEPOSIT_ACK), protocol.PROTO_NC_ITEM_DEPOSIT_ACK); e != 0x0F81 {
		t.Fatalf("DEPOSIT muerto esperaba 0x0F81, got 0x%04X", e)
	}
	pA.zp.State = 1

	depNoCen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DEPOSIT_REQ)
	depNoCen.WriteUint64(999_999_999)
	_ = pA.client.sendPacket(depNoCen)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_DEPOSIT_ACK), protocol.PROTO_NC_ITEM_DEPOSIT_ACK); e != 0x0F83 {
		t.Fatalf("DEPOSIT sin Cen esperaba 0x0F83, got 0x%04X", e)
	}

	depOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DEPOSIT_REQ)
	depOk.WriteUint64(1000)
	_ = pA.client.sendPacket(depOk)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_DEPOSIT_ACK), protocol.PROTO_NC_ITEM_DEPOSIT_ACK); e != 0x0F80 {
		t.Fatalf("DEPOSIT válido esperaba 0x0F80, got 0x%04X", e)
	}

	wdNoCen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_WITHDRAW_REQ)
	wdNoCen.WriteUint64(999_999)
	_ = pA.client.sendPacket(wdNoCen)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_WITHDRAW_ACK), protocol.PROTO_NC_ITEM_WITHDRAW_ACK); e != 0x0F83 {
		t.Fatalf("WITHDRAW sin saldo en banco esperaba 0x0F83, got 0x%04X", e)
	}

	wdOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_WITHDRAW_REQ)
	wdOk.WriteUint64(500)
	_ = pA.client.sendPacket(wdOk)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_WITHDRAW_ACK), protocol.PROTO_NC_ITEM_WITHDRAW_ACK); e != 0x0F80 {
		t.Fatalf("WITHDRAW válido esperaba 0x0F80, got 0x%04X", e)
	}

	gWdNoGuild := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ)
	gWdNoGuild.WriteUint64(100)
	_ = pA.client.sendPacket(gWdNoGuild)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK), protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK); e != 0x1482 {
		t.Fatalf("GUILD_STORAGE_WITHDRAW sin guild esperaba 0x1482, got 0x%04X", e)
	}

	pA.zp.State = 3
	acadDepDead := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ)
	acadDepDead.WriteUint64(100)
	_ = pA.client.sendPacket(acadDepDead)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK), protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK); e != 0x19A2 {
		t.Fatalf("ACADEMY_REWARD_STORAGE_DEPOSIT muerto esperaba 0x19A2, got 0x%04X", e)
	}
	pA.zp.State = 1

	acadDepNoGuild := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ)
	acadDepNoGuild.WriteUint64(100)
	_ = pA.client.sendPacket(acadDepNoGuild)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK), protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK); e != 0x19A3 {
		t.Fatalf("ACADEMY_REWARD_STORAGE_DEPOSIT sin guild esperaba 0x19A3, got 0x%04X", e)
	}

	acadWdFail := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ)
	acadWdFail.WriteUint64(999_999)
	_ = pA.client.sendPacket(acadWdFail)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK), protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK); e != 0x17E3 {
		t.Fatalf("ACADEMY_REWARD_STORAGE_WITHDRAW sin saldo esperaba 0x17E3, got 0x%04X", e)
	}

	// 5. CHARGED_WITHDRAW_REQ: sin abrir -> 0x1043; abrir + faltante -> 0x1044; éxito -> 0x1041
	findChargedErr := func(ps []*protocol.Packet) (uint16, bool) {
		p := f23FindPkt(ps, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)
		if p == nil || len(p.Raw) < 8 {
			return 0, false
		}
		return binary.LittleEndian.Uint16(p.Raw[6:8]), true
	}
	pA.zp.ChargedInvenOpened = false
	cwUnopen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ)
	cwUnopen.WriteUint32(1)
	_ = pA.client.sendPacket(cwUnopen)
	if e, ok := findChargedErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)); !ok || e != 0x1043 {
		t.Fatalf("CHARGED_WITHDRAW sin abrir esperaba 0x1043, got ok=%v err=0x%04X", ok, e)
	}

	orderNo := uint32(250000 + time.Now().Unix()%10000)
	if _, err := cl.accountRepo.InsertChargedItem(&database.ChargedItem{
		UserNo:  pA.userNo,
		OrderNo: orderNo,
		GoodsNo: uint32(s.stackPotion.ID),
		Amount:  1,
	}); err != nil {
		t.Fatalf("InsertChargedItem: %v", err)
	}

	chOpen := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ)
	chOpen.WriteUint16(0)
	_ = pA.client.sendPacket(chOpen)
	chOpenPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK)
	chOpenAck := f23FindPkt(chOpenPkts, protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK)
	if chOpenAck == nil || len(chOpenAck.Raw) < 11 {
		t.Fatalf("CHARGEDINVENOPEN_ACK inválido o sin items: %v", chOpenAck)
	}
	chargedRowNo := binary.LittleEndian.Uint32(chOpenAck.Raw[7:11])

	cwMissing := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ)
	cwMissing.WriteUint32(99999999)
	_ = pA.client.sendPacket(cwMissing)
	if e, ok := findChargedErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)); !ok || e != 0x1044 {
		t.Fatalf("CHARGED_WITHDRAW orderNo inexistente esperaba 0x1044, got ok=%v err=0x%04X", ok, e)
	}

	cwOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ)
	cwOk.WriteUint32(chargedRowNo)
	_ = pA.client.sendPacket(cwOk)
	if e, ok := findChargedErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)); !ok || e != 0x1041 {
		t.Fatalf("CHARGED_WITHDRAW válido esperaba 0x1041, got ok=%v err=0x%04X", ok, e)
	}
}

func runFase25ForgeCraftBranches(t *testing.T, cl *f25Cluster, pA *f25PlayerConn, s f25SampleItems) {
	t.Helper()
	f25ClearInventory(t, cl.charRepo, pA.charID, 9)

	eqItem := f25SeedItem(t, cl.charRepo, pA.charID, 9, 0, s.lowWeapon.ID, 1)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 1, s.upgradeStone.ID, 20)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 2, s.luckyProtect.ID, 10)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 3, s.bracelet.ID, 1)

	// 1. UPGRADE_REQ ramas de rechazo: estado (0x08D6), target faltante (0x08C8), piedra faltante (0x08CA), sin Cen (0x08CB)
	sendUpg := func(eqS, rawS, leftS, rightS, midS byte, gift uint32) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_UPGRADE_REQ)
		pw.WriteByte(eqS)
		pw.WriteByte(rawS)
		pw.WriteByte(leftS)
		pw.WriteByte(rightS)
		pw.WriteUint32(gift)
		pw.WriteByte(midS)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_UPGRADE_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_UPGRADE_ACK)
		return e, ps
	}
	pA.zp.State = 3
	if e, _ := sendUpg(0, 1, 0xFF, 0xFF, 0xFF, 0); e != 0x08D6 {
		t.Fatalf("UPGRADE muerto esperaba 0x08D6, got 0x%04X", e)
	}
	pA.zp.State = 1
	if e, _ := sendUpg(50, 1, 0xFF, 0xFF, 0xFF, 0); e != 0x08C8 {
		t.Fatalf("UPGRADE target faltante esperaba 0x08C8, got 0x%04X", e)
	}
	if e, _ := sendUpg(0, 50, 0xFF, 0xFF, 0xFF, 0); e != 0x08CA {
		t.Fatalf("UPGRADE piedra faltante esperaba 0x08CA, got 0x%04X", e)
	}
	savedMoney := pA.zp.Char.TCharacterState.Money
	pA.zp.LastSaveAt = time.Now()
	pA.zp.Char.TCharacterState.Money = 0
	if e, _ := sendUpg(0, 1, 0xFF, 0xFF, 0xFF, 5000); e != 0x08CB {
		t.Fatalf("UPGRADE sin Cen esperaba 0x08CB, got 0x%04X", e)
	}
	pA.zp.Char.TCharacterState.Money = savedMoney

	// Con piedra protectora (leftSlot=2) -> jamás debe destruir (0x08C7) aunque esté +8
	_ = cl.charRepo.SetItemOption(eqItem.ItemKey, 0, 8)
	if e, ps := sendUpg(0, 1, 2, 0xFF, 0xFF, 0); e != 0x08C3 && e != 0x08C4 && e != 0x08C5 {
		var ops []string
		for _, pk := range ps {
			ops = append(ops, fmt.Sprintf("0x%04X", pk.OpCode))
		}
		t.Fatalf("UPGRADE con piedra protectora esperaba 0x08C3/0x08C4/0x08C5 (sin rotura/bajada), got 0x%04X (pkts=%v)", e, ops)
	}

	// 2. NEW_UPGRADE_REQ & BRACELET_UPGRADE_REQ (usar arma +0 fresca en slot 8)
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 8, s.lowWeapon.ID, 1)
	newUpgBad := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_NEW_UPGRADE_REQ)
	newUpgBad.WriteByte(55)
	newUpgBad.WriteByte(1)
	_ = pA.client.sendPacket(newUpgBad)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_NEW_UPGRADE_ACK), protocol.PROTO_NC_ITEM_NEW_UPGRADE_ACK); e != 0x3801 && e != 0x3806 {
		t.Fatalf("NEW_UPGRADE target inválido esperaba 0x3801/0x3806, got 0x%04X", e)
	}

	newUpgOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_NEW_UPGRADE_REQ)
	newUpgOk.WriteByte(8)
	newUpgOk.WriteByte(1)
	_ = pA.client.sendPacket(newUpgOk)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_NEW_UPGRADE_ACK), protocol.PROTO_NC_ITEM_NEW_UPGRADE_ACK); e != 0x3803 && e != 0x3804 && e != 0x3805 {
		t.Fatalf("NEW_UPGRADE resultado esperaba 0x3803..0x3805, got 0x%04X", e)
	}

	// Asegurar que exista ítem en slot 0 para gemas/randomopt
	var hasSlot0 bool
	if allItems, _ := cl.charRepo.GetCharacterItems(pA.charID); len(allItems) > 0 {
		for _, it := range allItems {
			if it != nil && it.StorageType == 9 && it.Storage == 0 {
				eqItem = it
				hasSlot0 = true
				break
			}
		}
	}
	if !hasSlot0 {
		eqItem = f25SeedItem(t, cl.charRepo, pA.charID, 9, 0, s.lowWeapon.ID, 1)
	}
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 4, s.stackPotion.ID, 1)
	bracNot := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BRACELET_UPGRADE_REQ)
	bracNot.WriteByte(4)
	bracNot.WriteByte(1)
	_ = pA.client.sendPacket(bracNot)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_BRACELET_UPGRADE_ACK), protocol.PROTO_NC_ITEM_BRACELET_UPGRADE_ACK); e != 0x3C05 {
		t.Fatalf("BRACELET_UPGRADE sobre poción esperaba 0x3C05, got 0x%04X", e)
	}

	// 3. MOVER_UPGRADE_REQ
	mvMiss := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MOVER_UPGRADE_REQ)
	mvMiss.WriteByte(70)
	mvMiss.WriteByte(71)
	_ = pA.client.sendPacket(mvMiss)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK), protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK); e != 0x3229 {
		t.Fatalf("MOVER_UPGRADE faltante esperaba 0x3229, got 0x%04X", e)
	}

	// 4. ENCHANT ADD_NEW_SOCKET / ADD_GEM / REMOVE_GEM
	sockAdd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ)
	sockAdd.WriteUint16(f25PackInven(9, 0))
	sockAdd.WriteUint16(f25PackInven(9, 4))
	_ = pA.client.sendPacket(sockAdd)
	sockPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK)
	sockErr, ok := f25FindErr(sockPkts, protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK)
	if !ok || (sockErr != 0x20A0 && sockErr != 0x20A5) {
		t.Fatalf("ENCHANT_ADD_NEW_SOCKET esperaba 0x20A0/0x20A5, got ok=%v err=0x%04X", ok, sockErr)
	}
	_ = cl.charRepo.SetItemOption(eqItem.ItemKey, 524, 1) // garantizar 1 socket abierto (iotCreatedSocket=524) para probar ADD_GEM/REMOVE_GEM

	f25SeedItem(t, cl.charRepo, pA.charID, 9, 5, s.upgradeStone.ID, 5)
	gemMissing := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ)
	gemMissing.WriteUint16(f25PackInven(9, 77))
	gemMissing.WriteUint16(f25PackInven(9, 5))
	_ = pA.client.sendPacket(gemMissing)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK), protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK); e != 0x2085 {
		t.Fatalf("ENCHANT_ADD_GEM faltante esperaba 0x2085, got 0x%04X", e)
	}

	gemAdd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ)
	gemAdd.WriteUint16(f25PackInven(9, 0))
	gemAdd.WriteUint16(f25PackInven(9, 5))
	_ = pA.client.sendPacket(gemAdd)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK), protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK); e != 0x2080 {
		t.Fatalf("ENCHANT_ADD_GEM válido esperaba 0x2080, got 0x%04X", e)
	}

	f25SeedItem(t, cl.charRepo, pA.charID, 9, 6, s.stackPotion.ID, 3)
	gemRem := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ)
	gemRem.WriteUint16(f25PackInven(9, 0))
	gemRem.WriteUint16(f25PackInven(9, 6))
	gemRem.WriteByte(0)
	_ = pA.client.sendPacket(gemRem)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK), protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK); e != 0x2090 {
		t.Fatalf("ENCHANT_REMOVE_GEM válido esperaba 0x2090, got 0x%04X", e)
	}
	gemRemEmpty := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ)
	gemRemEmpty.WriteUint16(f25PackInven(9, 0))
	gemRemEmpty.WriteUint16(f25PackInven(9, 6))
	gemRemEmpty.WriteByte(0)
	_ = pA.client.sendPacket(gemRemEmpty)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK), protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK); e != 0x2093 {
		t.Fatalf("ENCHANT_REMOVE_GEM slot vacío esperaba 0x2093, got 0x%04X", e)
	}

	// 5. RANDOMOPTION CHANGE / ACCEPT / RECOVER_COUNT_LIMIT
	roMiss := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ)
	roMiss.WriteByte(79)
	roMiss.WriteByte(6)
	_ = pA.client.sendPacket(roMiss)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK), protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK); e != 0x3249 {
		t.Fatalf("RANDOMOPTION_CHANGE faltante esperaba 0x3249, got 0x%04X", e)
	}

	roOk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ)
	roOk.WriteByte(0)
	roOk.WriteByte(6)
	_ = pA.client.sendPacket(roOk)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK), protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK); e != 0x3248 {
		t.Fatalf("RANDOMOPTION_CHANGE válido esperaba 0x3248, got 0x%04X", e)
	}

	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_REQ))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_ACK), protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_ACK); e != 0x3248 {
		t.Fatalf("RANDOMOPTION_ACCEPT válido esperaba 0x3248, got 0x%04X", e)
	}

	f25SeedItem(t, cl.charRepo, pA.charID, 9, 7, s.stackPotion.ID, 5)
	roRec := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ)
	roRec.WriteByte(0)
	roRec.WriteByte(7)
	_ = pA.client.sendPacket(roRec)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK), protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK); e != 0x3251 {
		t.Fatalf("RANDOMOPTION_RECOVER válido esperaba 0x3251, got 0x%04X", e)
	}

	// 6. DISMANTLE / MIX_ITEM / ACT_PRODUCE_CAST + ABORT
	dismMiss := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DISMANTLE_REQ)
	dismMiss.WriteByte(81)
	_ = pA.client.sendPacket(dismMiss)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_DISMANTLE_ACK), protocol.PROTO_NC_ITEM_DISMANTLE_ACK); e != 0x1682 {
		t.Fatalf("DISMANTLE faltante esperaba 0x1682, got 0x%04X", e)
	}

	mixBad := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MIX_ITEM_REQ)
	mixBad.WriteByte(0)
	mixBad.WriteByte(81)
	mixBad.WriteByte(82)
	_ = pA.client.sendPacket(mixBad)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_MIX_ITEM_ACK), protocol.PROTO_NC_ITEM_MIX_ITEM_ACK); e != 0x37C2 && e != 0x37C6 {
		t.Fatalf("MIX_ITEM faltante esperaba 0x37C2/0x37C6, got 0x%04X", e)
	}

	prodBad := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CAST_REQ)
	prodBad.WriteUint16(64999)
	_ = pA.client.sendPacket(prodBad)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK), protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK); e != 0x0B04 {
		t.Fatalf("ACT_PRODUCE receta inválida esperaba CAST_FAIL_ACK=0x0B04, got 0x%04X", e)
	}

	if s.produceID > 0 {
		prodNoMat := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CAST_REQ)
		prodNoMat.WriteUint16(s.produceID)
		_ = pA.client.sendPacket(prodNoMat)
		if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK), protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK); e != 0x0B08 {
			t.Fatalf("ACT_PRODUCE sin materiales esperaba CAST_FAIL_ACK=0x0B08, got 0x%04X", e)
		}

		for i, mat := range s.produceMats {
			f25SeedItem(t, cl.charRepo, pA.charID, 9, int16(10+i), mat.id, mat.qty+5)
		}
		prodOk := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CAST_REQ)
		prodOk.WriteUint16(s.produceID)
		_ = pA.client.sendPacket(prodOk)
		prodPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_ACT_PRODUCE_CAST_SUC_ACK, protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK)
		if f23FindPkt(prodPkts, protocol.PROTO_NC_ACT_PRODUCE_CAST_SUC_ACK) == nil {
			t.Fatalf("ACT_PRODUCE válido esperaba CAST_SUC_ACK")
		}
		_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CASTABORT_CMD))
		_ = f23CollectPackets(pA.conn, 50*time.Millisecond)
	}
}

func runFase25UtilitiesAndTradeBoothBranches(t *testing.T, cl *f25Cluster, pA, pB *f25PlayerConn, s f25SampleItems) {
	t.Helper()
	f25ClearInventory(t, cl.charRepo, pA.charID, 9)
	f25ClearInventory(t, cl.charRepo, pB.charID, 9)

	// 1. USE_REQ: muerto -> 0x0701, box inválido -> 0x0702, faltante -> 0x0701, poción válida -> 0x0700
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 0, s.stackPotion.ID, 3)
	sendUse := func(slot, box uint8) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_REQ)
		pw.WriteByte(slot)
		pw.WriteByte(box)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_USE_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_USE_ACK)
		return e, ps
	}
	pA.zp.State = 3
	if e, _ := sendUse(0, 9); e != 0x0701 {
		t.Fatalf("USE muerto esperaba 0x0701, got 0x%04X", e)
	}
	pA.zp.State = 1
	if e, _ := sendUse(0, 12); e != 0x0702 {
		t.Fatalf("USE box inválido esperaba 0x0702, got 0x%04X", e)
	}
	if e, _ := sendUse(40, 9); e != 0x0701 {
		t.Fatalf("USE faltante esperaba 0x0701, got 0x%04X", e)
	}
	if e, _ := sendUse(0, 9); e != 0x0700 {
		t.Fatalf("USE válido esperaba 0x0700, got 0x%04X", e)
	}

	// 2. USE_ACTIVESKILL_REQ (0x304F): muerto -> 0x3129, faltante -> 0x312C, válido -> 0x3128
	f25SeedItem(t, cl.charRepo, pA.charID, 9, 1, s.skillScroll.ID, 2)
	sendUseSkill := func(slot uint8) (uint16, []*protocol.Packet) {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_REQ)
		pw.WriteByte(slot)
		pw.WriteUint16(0)
		pw.WriteUint32(0)
		pw.WriteUint32(0)
		_ = pA.client.sendPacket(pw)
		ps := f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_ACK)
		e, _ := f25FindErr(ps, protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_ACK)
		return e, ps
	}
	pA.zp.State = 3
	if e, _ := sendUseSkill(1); e != 0x3129 {
		t.Fatalf("USE_ACTIVESKILL muerto esperaba 0x3129, got 0x%04X", e)
	}
	pA.zp.State = 1
	if e, _ := sendUseSkill(42); e != 0x312C {
		t.Fatalf("USE_ACTIVESKILL faltante esperaba 0x312C, got 0x%04X", e)
	}
	if e, ps := sendUseSkill(1); e != 0x3128 || f23FindPkt(ps, protocol.PROTO_NC_SKILL_SKILL_LEARNSUC_CMD) == nil {
		t.Fatalf("USE_ACTIVESKILL válido esperaba 0x3128 + SKILL_LEARNSUC_CMD, got 0x%04X", e)
	}

	// 3. REVIVEITEMUSE_CMD: vivo -> 0x303A (0x0701); muerto sin ítem -> 0x303A (0x071E); con ítem revive -> State=1
	pA.zp.State = 1
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD), protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD); e != 0x0701 {
		t.Fatalf("REVIVEITEMUSE vivo esperaba 0x0701, got 0x%04X", e)
	}
	pA.zp.State = 3
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD), protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD); e != 0x071E {
		t.Fatalf("REVIVEITEMUSE sin ítem esperaba 0x071E, got 0x%04X", e)
	}
	if s.reviveScroll != nil {
		f25SeedItem(t, cl.charRepo, pA.charID, 9, 2, s.reviveScroll.ID, 1)
		_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD))
		_ = f25WaitOp(pA.conn, protocol.PROTO_NC_ITEM_CELLCHANGE_CMD, protocol.PROTO_NC_CHAR_REVIVESAME_CMD, protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD)
		if pA.zp.State != 1 {
			t.Fatalf("REVIVEITEMUSE con ítem debía revivir a State=1, got %d", pA.zp.State)
		}
	}
	pA.zp.State = 1

	// 4. SafeTrade ramas: distancia >1000 (0x06C3), NoTrade (0x0851), Cen insuficiente (0x0848), Decide sin Lock (0x084A), break por distancia
	px, py := pA.zp.GetPosition()
	f25SetPos(pB, px+2000, py)
	propFar := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSE_REQ)
	propFar.WriteUint16(pB.zp.MapObjectID)
	_ = pA.client.sendPacket(propFar)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_PROPOSENO_ACK), protocol.PROTO_NC_TRADE_PROPOSENO_ACK); e != 0x06C3 {
		t.Fatalf("TRADE_PROPOSE lejos esperaba 0x06C3, got 0x%04X", e)
	}

	f25SetPos(pB, px+10, py+10)
	propOk := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSE_REQ)
	propOk.WriteUint16(pB.zp.MapObjectID)
	_ = pA.client.sendPacket(propOk)
	askPkts := f25WaitOp(pB.conn, protocol.PROTO_NC_TRADE_PROPOSE_ASK_REQ)
	if f23FindPkt(askPkts, protocol.PROTO_NC_TRADE_PROPOSE_ASK_REQ) == nil {
		t.Fatalf("TRADE_PROPOSE_ASK_REQ no llegó a pB (ops=%v, pB.Pending=%d)", len(askPkts), pB.zp.PendingTradeFrom)
	}
	_ = pB.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSEYES_ACK))
	startPA := f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_START_CMD)
	if f23FindPkt(startPA, protocol.PROTO_NC_TRADE_START_CMD) == nil {
		t.Fatalf("TRADE_START_CMD no llegó a pA (activeTrade=%v)", pA.zp.ActiveTrade != nil)
	}
	startPB := f25WaitOp(pB.conn, protocol.PROTO_NC_TRADE_START_CMD)
	if f23FindPkt(startPB, protocol.PROTO_NC_TRADE_START_CMD) == nil {
		t.Fatalf("TRADE_START_CMD no llegó a pB (activeTrade=%v)", pB.zp.ActiveTrade != nil)
	}

	if s.noTradeItem != nil {
		f25SeedItem(t, cl.charRepo, pA.charID, 9, 3, s.noTradeItem.ID, 1)
		upNoTr := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_UPBOARD_REQ)
		upNoTr.WriteByte(3)
		_ = pA.client.sendPacket(upNoTr)
		upPkts := f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, protocol.PROTO_NC_TRADE_UPBOARD_ACK, protocol.PROTO_NC_TRADE_CANCEL_ACK)
		if e, ok := f25FindErr(upPkts, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK); !ok || e != 0x0851 {
			var ops []string
			for _, pk := range upPkts {
				ops = append(ops, fmt.Sprintf("0x%04X", pk.OpCode))
			}
			t.Fatalf("TRADE_UPBOARD NoTrade esperaba 0x0851, got ok=%v err=0x%04X (itemID=%d inx=%q activeTrade=%v pkts=%v)",
				ok, e, s.noTradeItem.ID, s.noTradeItem.InxName, pA.zp.ActiveTrade != nil, ops)
		}
	}

	pA.zp.LastSaveAt = time.Now()
	pA.zp.Char.TCharacterState.Money = 100
	cenBad := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_CENBOARDING_REQ)
	cenBad.WriteUint64(999999)
	_ = pA.client.sendPacket(cenBad)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_CENBOARDINGFAIL_ACK), protocol.PROTO_NC_TRADE_CENBOARDINGFAIL_ACK); e != 0x0848 {
		t.Fatalf("TRADE_CENBOARDING sin Cen esperaba 0x0848, got 0x%04X", e)
	}

	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_DECIDE_REQ))
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_DECIDEFAIL_ACK), protocol.PROTO_NC_TRADE_DECIDEFAIL_ACK); e != 0x084A {
		t.Fatalf("TRADE_DECIDE sin Lock esperaba 0x084A, got 0x%04X", e)
	}

	// Alejar a pB > 1000 durante trade activo -> CANCEL_ACK / CANCEL_CMD automático
	f25SetPos(pB, px+2000, py)
	_ = pA.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_BOARDLOCK_REQ))
	pktsBreakA := f25WaitOp(pA.conn, protocol.PROTO_NC_TRADE_CANCEL_ACK)
	if f23FindPkt(pktsBreakA, protocol.PROTO_NC_TRADE_CANCEL_ACK) == nil || pA.zp.ActiveTrade != nil {
		t.Fatalf("TRADE lejos durante sesión debía cancelar vía sp_tradebreak")
	}
	_ = f23CollectPackets(pB.conn, 50*time.Millisecond)

	// 5. BOOTH ramas: ENTRY sobre cerrado (0x0A5B), distancia >=300 (0x0A4C), sin Cen (0x020E), éxito (0x0A41) + CLOSE (0x0A41)
	bEntryClosed := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ENTRY_REQ)
	bEntryClosed.WriteUint16(pB.zp.MapObjectID)
	_ = pA.client.sendPacket(bEntryClosed)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_BOOTH_ENTRY_SELL_ACK), protocol.PROTO_NC_BOOTH_ENTRY_SELL_ACK); e != 0x0A5B {
		t.Fatalf("BOOTH_ENTRY sobre tienda cerrada esperaba 0x0A5B, got 0x%04X", e)
	}

	f25SetPos(pB, px+20, py+20)
	f25SeedItem(t, cl.charRepo, pB.charID, 9, 0, s.lowWeapon.ID, 1)
	bOpen := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_OPEN_REQ)
	bOpen.WriteString("F25Booth", 30)
	bOpen.WriteByte(1 | (1 << 1))
	bOpen.WriteByte(0)
	bOpen.WriteByte(0)
	bOpen.WriteUint64(5000)
	bOpen.WriteUint16(1)
	_ = pB.client.sendPacket(bOpen)
	if e, _ := f25FindErr(f25WaitOp(pB.conn, protocol.PROTO_NC_BOOTH_OPEN_ACK), protocol.PROTO_NC_BOOTH_OPEN_ACK); e != 0x0A41 {
		t.Fatalf("BOOTH_OPEN válido esperaba 0x0A41, got 0x%04X", e)
	}
	_ = f23CollectPackets(pA.conn, 50*time.Millisecond)

	// Distancia >= 300 -> 0x0A4C
	f25SetPos(pA, px+500, py)
	bBuyFar := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ITEMTRADE_REQ)
	bBuyFar.WriteUint16(pB.zp.MapObjectID)
	bBuyFar.WriteByte(0)
	bBuyFar.WriteUint16(1)
	_ = pA.client.sendPacket(bBuyFar)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK), protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK); e != 0x0A4C {
		t.Fatalf("BOOTH_ITEMTRADE lejos (>=300u) esperaba 0x0A4C, got 0x%04X", e)
	}

	// En rango sin Cen -> 0x020E
	f25SetPos(pA, px+20, py+20)
	pA.zp.LastSaveAt = time.Now()
	pA.zp.Char.TCharacterState.Money = 100
	bBuyNoCen := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ITEMTRADE_REQ)
	bBuyNoCen.WriteUint16(pB.zp.MapObjectID)
	bBuyNoCen.WriteByte(0)
	bBuyNoCen.WriteUint16(1)
	_ = pA.client.sendPacket(bBuyNoCen)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK), protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK); e != 0x020E {
		t.Fatalf("BOOTH_ITEMTRADE sin Cen esperaba 0x020E, got 0x%04X", e)
	}

	// Con Cen -> 0x0A41
	pA.zp.LastSaveAt = time.Now()
	pA.zp.Char.TCharacterState.Money = 20000
	bBuyOk := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ITEMTRADE_REQ)
	bBuyOk.WriteUint16(pB.zp.MapObjectID)
	bBuyOk.WriteByte(0)
	bBuyOk.WriteUint16(1)
	_ = pA.client.sendPacket(bBuyOk)
	if e, _ := f25FindErr(f25WaitOp(pA.conn, protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK), protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK); e != 0x0A41 {
		t.Fatalf("BOOTH_ITEMTRADE válido esperaba 0x0A41, got 0x%04X", e)
	}

	_ = pB.client.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_CLOSE_REQ))
	if e, _ := f25FindErr(f25WaitOp(pB.conn, protocol.PROTO_NC_BOOTH_CLOSE_ACK), protocol.PROTO_NC_BOOTH_CLOSE_ACK); e != 0x0A41 || pB.zp.ActiveBooth != nil {
		t.Fatalf("BOOTH_CLOSE esperaba 0x0A41 y ActiveBooth=nil, got 0x%04X", e)
	}
}
