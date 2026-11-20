package zone

import (
	"encoding/base64"
	"encoding/binary"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func makeTestKQInfo(handle uint32, id uint16, baseMap, subMap, scriptLang string) kq.PROTO_KQ_INFO {
	var info kq.PROTO_KQ_INFO
	info.Handle = handle
	info.ID = id
	info.LimitTime = 30
	info.PlayerRevivalMode = 2
	info.PlayerRevivalCount = 3
	info.SetTitle("Test KQ")
	info.SetScriptLanguage(scriptLang)
	info.SetScriptInitValue("1")
	info.MapLink[0].MapIndex = 0
	info.MapLink[0].SetMapBase(baseMap)
	info.MapLink[0].SetMapName(subMap)
	info.MapLink[0].MapClear = 1
	return info
}

func makeW2ZMakePayload(t *testing.T, info kq.PROTO_KQ_INFO) map[string]interface{} {
	t.Helper()
	raw, err := info.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary(PROTO_KQ_INFO): %v", err)
	}
	return map[string]interface{}{
		"handle":      info.Handle,
		"id":          info.ID,
		"kq_info_b64": base64.StdEncoding.EncodeToString(raw),
	}
}

func lastWMSEventBySymbol(b *memWMSBridge, sym string) (WMSEvent, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := len(b.events) - 1; i >= 0; i-- {
		if b.events[i].Symbol == sym {
			return b.events[i], true
		}
	}
	return WMSEvent{}, false
}

func clearWMSBridge(b *memWMSBridge) {
	b.mu.Lock()
	b.events = nil
	b.mu.Unlock()
}

func TestKQZone_MakeReq_MapAllocAndAllErrorCodes(t *testing.T) {
	h, ds, _, _, _ := setupFase14Handler(t)
	if ds.KQStore == nil || len(ds.KQStore.Quests) != 57 {
		t.Fatalf("expected 57 KingdomQuest rows loaded in DataStore.KQStore, got %v", ds.KQStore)
	}
	if ds.KQItemBox == nil {
		t.Fatalf("expected KQItemBox initialized from KQItem.shn")
	}
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	// Pre-populate a mob on KDEddyHill00 to verify MapClear=1 triggers fm_ClearObject(0x3B6)
	preMap := h.mapManager.EnsureKQMap("KDEddyHill00", "KDEddyHill")
	if preMap == nil {
		t.Fatalf("expected KDEddyHill00 to be hostable")
	}
	_ = preMap.SpawnMobFromTemplate(&MobTemplate{ID: 1, Name: "PreMob", Level: 1, MaxHP: 100}, 1000, 1000, 0, "", 0, true)

	// 1. Valid W2Z_MAKE_REQ -> Z2W_MAKE_ACK(0x0981)
	info1 := makeTestKQInfo(100, 1, "KDEddyHill", "KDEddyHill00", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info1))
	ev, ok := lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if !ok {
		t.Fatalf("expected wms_NC_KQ_Z2W_MAKE_ACK")
	}
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeNoErr {
		t.Fatalf("expected err=0x%04X, got 0x%04X", kqErrMakeNoErr, errCode)
	}
	if preMap.fm_GetKQhandle() != 100 || preMap.KQID != 1 {
		t.Fatalf("expected KDEddyHill00 bound to handle=100, id=1; got handle=%d id=%d", preMap.fm_GetKQhandle(), preMap.KQID)
	}
	if len(preMap.mobs) != 0 {
		t.Fatalf("expected MapClear=1 to clear pre-existing mobs via fm_ClearObject(0x3B6), remaining=%d", len(preMap.mobs))
	}
	if h.kqState().container.kqc_MapUseCheck("KDEddyHill00") {
		t.Fatalf("expected kqc_MapUseCheck(KDEddyHill00) == false while in use")
	}
	if !h.kqState().container.kqc_MapUseCheck("KDEddyHill01") {
		t.Fatalf("expected kqc_MapUseCheck(KDEddyHill01) == true while free")
	}

	// 2. Duplicate handle -> 0x0982
	clearWMSBridge(bridge)
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info1))
	ev, ok = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if !ok {
		t.Fatalf("expected Z2W_MAKE_ACK on duplicate handle")
	}
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeDuplicate {
		t.Fatalf("expected duplicate err 0x%04X, got 0x%04X", kqErrMakeDuplicate, errCode)
	}

	// 3. Missing script -> 0x098C
	clearWMSBridge(bridge)
	infoBadScript := makeTestKQInfo(101, 1, "KDEddyHill", "KDEddyHill01", "KQ/NoSuchScript/Missing")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoBadScript))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeScriptNotFound {
		t.Fatalf("expected script-not-found err 0x%04X, got 0x%04X", kqErrMakeScriptNotFound, errCode)
	}

	// 4. Empty MapLink[0].MapName -> 0x0984
	clearWMSBridge(bridge)
	infoNoMap0 := makeTestKQInfo(102, 1, "KDEddyHill", "", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoNoMap0))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeNoMap[0] {
		t.Fatalf("expected no-map0 err 0x%04X, got 0x%04X", kqErrMakeNoMap[0], errCode)
	}

	// 5. Unknown MapBase in slot 1 -> 0x0985
	clearWMSBridge(bridge)
	infoNoMap1 := makeTestKQInfo(103, 1, "KDEddyHill", "KDEddyHill02", "KQ/KingSlime/KingSlime")
	infoNoMap1.MapLink[1].SetMapBase("NoSuchMap")
	infoNoMap1.MapLink[1].SetMapName("NoSuchMap00")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoNoMap1))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeNoMap[1] {
		t.Fatalf("expected no-map1 err 0x%04X, got 0x%04X", kqErrMakeNoMap[1], errCode)
	}
	// Verify slot 0 rollback after slot 1 failed
	if m2 := h.mapManager.GetMap("KDEddyHill02"); m2 != nil && m2.fm_GetKQhandle() != kqInvalidHandle {
		t.Fatalf("expected KDEddyHill02 rolled back to unbound on slot 1 failure")
	}

	// 6. Submap already in use in slot 0 -> 0x0988, and in slot 1 -> 0x0989
	clearWMSBridge(bridge)
	infoInUse0 := makeTestKQInfo(104, 1, "KDEddyHill", "KDEddyHill00", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoInUse0))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeMapInUse[0] {
		t.Fatalf("expected in-use0 err 0x%04X, got 0x%04X", kqErrMakeMapInUse[0], errCode)
	}

	clearWMSBridge(bridge)
	infoInUse1 := makeTestKQInfo(105, 1, "KDEddyHill", "KDEddyHill03", "KQ/KingSlime/KingSlime")
	infoInUse1.MapLink[1].SetMapBase("KDEddyHill")
	infoInUse1.MapLink[1].SetMapName("KDEddyHill00")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoInUse1))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeMapInUse[1] {
		t.Fatalf("expected in-use1 err 0x%04X, got 0x%04X", kqErrMakeMapInUse[1], errCode)
	}

	// 7. Container full (0x12C = 300) -> 0x0983 (ERR_KINGDOMQUEST_MAKE_TOOMANY)
	st := h.kqState()
	st.container.mu.Lock()
	for i := 0; i < kqContainerMaxElements; i++ {
		if st.container.slots[i] == nil {
			el := &KQElement{}
			el.Info.Handle = uint32(9000 + i)
			st.container.slots[i] = el
		}
	}
	st.container.mu.Unlock()
	clearWMSBridge(bridge)
	infoFull := makeTestKQInfo(106, 1, "KDEddyHill", "KDEddyHill04", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, infoFull))
	ev, _ = lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_MAKE_ACK")
	if errCode, _ := payloadInt(ev.Payload, "err"); uint16(errCode) != kqErrMakeTooMany {
		t.Fatalf("expected too-many err 0x%04X, got 0x%04X", kqErrMakeTooMany, errCode)
	}
}

func TestKQZone_StartCmd_Entrance_GateKeeperWaves_And_TeamRegen(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	baseTime := time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC)
	curTime := baseTime
	h.SetKQClock(func() time.Time { return curTime })

	// Create 4 players on zMap
	pRed, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 7001, "RedPlayer", 30, 1500, 1500)
	pBlue, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 7002, "BluePlayer", 30, 1500, 1500)
	pBooth, connBooth := newTestClientAndPlayer(t, repo, ds, zMap, h, 7003, "BoothPlayer", 30, 1500, 1500)
	pBlocked, connBlocked := newTestClientAndPlayer(t, repo, ds, zMap, h, 7004, "BlockedPlayer", 30, 1500, 1500)

	// Setup states: pBlue dead at t+0, pBooth in booth at t+0..t+2, pBlocked on LinkOut=false map
	pRed.State = 1
	pBlue.State = 3 // Dead at t+0
	pBooth.State = 5 // Booth at t+0
	lockedMap := NewZoneMap(999, "LockedArena", "LockedArena")
	feCopy := *zMap.FieldEntry
	feCopy.Teleport.LinkOut = false
	lockedMap.FieldEntry = &feCopy
	pBlocked.CurrentMap = lockedMap
	pBlocked.State = 1

	// 1. W2Z_MAKE_REQ with TeamRegenXY for Red(0) and Blue(1)
	info := makeTestKQInfo(200, 2, "KDEddyHill", "KDEddyHill05", "KQ/KingSlime/KingSlime")
	info.LimitTime = 20
	info.IsTeamPVP = 1
	info.TeamRegenXY[0] = kq.SHINE_XY_TYPE{X: 1111, Y: 2222}
	info.TeamRegenXY[1] = kq.SHINE_XY_TYPE{X: 3333, Y: 4444}
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	// 2. W2Z_START_CMD -> triggers wave 0 (t+0s)
	raw, _ := info.MarshalBinary()
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_START_CMD", 0, 0, map[string]interface{}{
		"handle":      info.Handle,
		"id":          info.ID,
		"kq_info_b64": base64.StdEncoding.EncodeToString(raw),
		"joiners": []map[string]interface{}{
			{"char_reg_num": uint32(pRed.Char.ID), "char_name": pRed.Char.Name, "team_type": uint8(0)},
			{"char_reg_num": uint32(pBlue.Char.ID), "char_name": pBlue.Char.Name, "team_type": uint8(1)},
			{"char_reg_num": uint32(pBooth.Char.ID), "char_name": pBooth.Char.Name, "team_type": uint8(2)},
			{"char_reg_num": uint32(pBlocked.Char.ID), "char_name": pBlocked.Char.Name, "team_type": uint8(2)},
		},
	})

	// Verify kqe_QuestStart: EndTick = curTick + 20 * 0x258 = 12000, LimitMin = 20
	elem, _, ok := h.kqState().container.kqc_SearchHandle(200)
	if !ok || elem == nil {
		t.Fatalf("expected KQElement handle 200")
	}
	if elem.EndTick != uint64(20)*kqTicksPerMinute || elem.kqe_GetKQLimitMin() != 20 {
		t.Fatalf("unexpected EndTick=%d LimitMin=%d", elem.EndTick, elem.kqe_GetKQLimitMin())
	}
	if !elem.Maps[0].FieldMap.KQScriptActive {
		t.Fatalf("expected cc_PlayFilm to mark KQScriptActive=true")
	}

	// Wave 0 (t+0s): pRed entered at Red TeamRegenXY (1111, 2222); pBlue dead; pBooth in booth; pBlocked rejected
	if pRed.CurrentMap == nil || pRed.CurrentMap.MapID != "KDEddyHill05" || pRed.Char.XPos != 1111 || pRed.Char.YPos != 2222 {
		t.Fatalf("expected pRed warped to KDEddyHill05 (1111,2222), got map=%v (%d,%d)", pRed.CurrentMap, pRed.Char.XPos, pRed.Char.YPos)
	}
	if pBlue.CurrentMap.MapID == "KDEddyHill05" {
		t.Fatalf("expected dead pBlue to be skipped at wave 0")
	}
	if len(connBooth.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) < 2 {
		t.Fatalf("expected KQEntryInCamp0/1 notices for booth player")
	}
	if len(connBlocked.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) < 2 {
		t.Fatalf("expected KQEntryMapErr0/1 notices for LinkOut=false map")
	}

	// Wave 1 (t+2s): pBlue revives -> enters at Blue TeamRegenXY (3333, 4444)
	pBlue.State = 1
	curTime = baseTime.Add(2 * time.Second)
	entered1 := h.ProcessKQGateKeepers(curTime)
	if entered1 != 1 {
		t.Fatalf("expected 1 player entered at wave t+2s, got %d", entered1)
	}
	if pBlue.CurrentMap.MapID != "KDEddyHill05" || pBlue.Char.XPos != 3333 || pBlue.Char.YPos != 4444 {
		t.Fatalf("expected pBlue warped to KDEddyHill05 (3333,4444), got map=%s (%d,%d)", pBlue.CurrentMap.MapID, pBlue.Char.XPos, pBlue.Char.YPos)
	}

	// Wave 2 (t+5s): pBooth closes booth -> enters at default MapInfo RegenX/RegenY (TeamType=2)
	pBooth.State = 1
	curTime = baseTime.Add(5 * time.Second)
	entered2 := h.ProcessKQGateKeepers(curTime)
	if entered2 != 1 {
		t.Fatalf("expected 1 player entered at wave t+5s, got %d", entered2)
	}
	if pBooth.CurrentMap.MapID != "KDEddyHill05" || uint32(pBooth.Char.XPos) != elem.Maps[0].RegenX || uint32(pBooth.Char.YPos) != elem.Maps[0].RegenY {
		t.Fatalf("expected pBooth at default regen (%d,%d), got (%d,%d)", elem.Maps[0].RegenX, elem.Maps[0].RegenY, pBooth.Char.XPos, pBooth.Char.YPos)
	}

	// Wave 3 (t+7s): pBlocked moves to normal map -> enters -> kqe_allentry becomes true
	pBlocked.CurrentMap = zMap
	curTime = baseTime.Add(7 * time.Second)
	entered3 := h.ProcessKQGateKeepers(curTime)
	if entered3 != 1 {
		t.Fatalf("expected 1 player entered at wave t+7s, got %d", entered3)
	}

	// Wave 4 (t+10s): all entered -> kqe_allentry short-circuits, 0 entered
	curTime = baseTime.Add(10 * time.Second)
	if entered4 := h.ProcessKQGateKeepers(curTime); entered4 != 0 {
		t.Fatalf("expected 0 entered at wave t+10s after kqe_allentry, got %d", entered4)
	}

	// Verify KingdomQuestEntrance cap at 0x64 (100) and descending kqe_sort/kqe_search
	many := make([]KQJoinerEntry, 120)
	for i := range many {
		many[i] = KQJoinerEntry{CharRegNum: uint32(i + 1), TeamType: uint8(i % 2)}
	}
	ent := &KingdomQuestEntrance{}
	idx := ent.kqe_newentrance(many)
	if ent.slots[idx].Total != kqEntranceMaxJoiners {
		t.Fatalf("expected entrance joiners capped at %d, got %d", kqEntranceMaxJoiners, ent.slots[idx].Total)
	}
	ent.kqe_sort(idx)
	if ent.slots[idx].Joiners[0].CharRegNum != 100 || ent.slots[idx].Joiners[99].CharRegNum != 1 {
		t.Fatalf("expected descending kqe_Compare order [100..1], got first=%d last=%d",
			ent.slots[idx].Joiners[0].CharRegNum, ent.slots[idx].Joiners[99].CharRegNum)
	}
	if hit, ok := ent.kqe_search(57, idx); !ok || hit == nil || hit.CharRegNum != 57 {
		t.Fatalf("expected kqe_search(57) to succeed")
	}
}

func TestKQZone_EntryResponseAck_XorDecode_And_NoZeroHack(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	p, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 7101, "EntryHero", 30, 1500, 1500)
	p.State = 1

	// 1. Reply=1, EncHandle=0 with no KQ present -> MUST NOT fabricate kqID=1!
	bufZero := make([]byte, 5)
	bufZero[0] = 1
	binary.LittleEndian.PutUint32(bufZero[1:5], 0)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, bufZero))
	if p.CurrentMap.MapID != zMap.MapID {
		t.Fatalf("hack regression: EncHandle=0 warped player to %s instead of ignoring", p.CurrentMap.MapID)
	}

	// 2. Make real KQ handle=300
	info := makeTestKQInfo(300, 1, "KDEddyHill", "KDEddyHill06", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	// 3. Reply=0 -> rejected
	bufNo := make([]byte, 5)
	bufNo[0] = 0
	binary.LittleEndian.PutUint32(bufNo[1:5], uint32(300)^uint32(p.MapObjectID))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, bufNo))
	if p.CurrentMap.MapID != zMap.MapID {
		t.Fatalf("expected Reply=0 to be ignored")
	}

	// 4. Dead state (3 not in 0x56 mask) -> rejected
	p.State = 3
	bufYes := make([]byte, 5)
	bufYes[0] = 1
	binary.LittleEndian.PutUint32(bufYes[1:5], uint32(300)^uint32(p.MapObjectID))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, bufYes))
	if p.CurrentMap.MapID != zMap.MapID {
		t.Fatalf("expected State=3 (outside 0x56 bitmask) to be rejected")
	}

	// 5. State=4 (Resting, inside 0x56 mask) + valid XOR EncHandle -> stands up & warps to KDEddyHill06
	p.State = 4
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, bufYes))
	if p.State != 1 {
		t.Fatalf("expected Resting state (4) reset to Normal (1), got %d", p.State)
	}
	if p.CurrentMap.MapID != "KDEddyHill06" {
		t.Fatalf("expected warp to KDEddyHill06 via XOR EncHandle, got %s", p.CurrentMap.MapID)
	}
}

func TestKQZone_ReviveModes_KnockDowns_TeamPVP_And_EscapeFromEmptyKQ(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)

	// 1. Verify all kqe_ReviveNumber modes (exe: 0x00496B90, jmptbl 0x00496C34)
	cases := []struct {
		mode     uint8
		cnt      uint8
		wantMax  int
		wantFlag uint8
	}{
		{0, 5, 0, 0},
		{1, 5, kqReviveUnlimited, 0},
		{2, 3, 3, 0},
		{3, 5, kqReviveUnlimited, 1},
		{4, 2, 2, 1},
		{9, 5, 0, 0},
	}
	for _, tc := range cases {
		el := &KQElement{}
		el.Info.PlayerRevivalMode = tc.mode
		el.Info.PlayerRevivalCount = tc.cnt
		gotMax, gotFlag := el.kqe_ReviveNumber()
		if gotMax != tc.wantMax || gotFlag != tc.wantFlag {
			t.Fatalf("kqe_ReviveNumber(mode=%d,cnt=%d) = (%d,%d), want (%d,%d)",
				tc.mode, tc.cnt, gotMax, gotFlag, tc.wantMax, tc.wantFlag)
		}
	}

	// 2. Make + Start TeamPVP KQ (mode=2, maxRevives=2)
	info := makeTestKQInfo(400, 1, "KDEddyHill", "KDEddyHill07", "KQ/KingSlime/KingSlime")
	info.IsTeamPVP = 1
	info.PlayerRevivalMode = 2
	info.PlayerRevivalCount = 2
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	kqMap := h.mapManager.GetMap("KDEddyHill07")
	pJoiner, connJoiner := newTestClientAndPlayer(t, repo, ds, kqMap, h, 7201, "JoinerA", 30, 1500, 1500)
	pEnemy, _ := newTestClientAndPlayer(t, repo, ds, kqMap, h, 7202, "EnemyB", 30, 1500, 1500)
	pIntruder, _ := newTestClientAndPlayer(t, repo, ds, kqMap, h, 7203, "Intruder", 30, 1500, 1500)
	pGM, connGM := newTestClientAndPlayer(t, repo, ds, kqMap, h, 7204, "GMObserver", 30, 1500, 1500)
	pGM.Char.AdminLevel = 100

	elem, _, _ := h.kqState().container.kqc_SearchHandle(400)
	elem.PlayerInfoList.kqpil_InsertPlayerInfo(uint32(pJoiner.Char.ID), 0, 1)
	elem.PlayerInfoList.kqpil_InsertPlayerInfo(uint32(pEnemy.Char.ID), 1, 1)

	// Missing player returns 0x1869F (99999) knockdown count
	if kd := elem.PlayerInfoList.kqpil_StorePlayer_GetKnockDownNumber(uint32(pIntruder.Char.ID)); kd != kqReviveUnlimited {
		t.Fatalf("expected missing player knockdown=0x1869F, got %d", kd)
	}

	// 3. MapLoginComplete KQ hooks for Joiner (kd=0 < 2 -> stays, gets TEAM_TYPE_CMD=0 and RESTDEADNUM_CMD=2)
	connJoiner.clear()
	h.onPlayerMapLoginCompleteKQ(pJoiner)
	if pJoiner.CurrentMap.MapID != "KDEddyHill07" || pJoiner.KQTeam != 0 {
		t.Fatalf("expected JoinerA to stay on KDEddyHill07 with team 0, got map=%s team=%d", pJoiner.CurrentMap.MapID, pJoiner.KQTeam)
	}
	teamPkts := connJoiner.findByOpcode(protocol.PROTO_NC_KQ_TEAM_TYPE_CMD)
	if len(teamPkts) == 0 || teamPkts[0].Payload[0] != 0 {
		t.Fatalf("expected PROTO_NC_KQ_TEAM_TYPE_CMD(0)")
	}
	restPkts := connJoiner.findByOpcode(protocol.PROTO_NC_KQ_RESTDEADNUM_CMD)
	if len(restPkts) == 0 || restPkts[0].Payload[0] != 2 {
		t.Fatalf("expected PROTO_NC_KQ_RESTDEADNUM_CMD(2), got %+v", restPkts)
	}

	// TeamPVP enemy checks
	pEnemy.KQTeam = 1
	if !h.so_IsTeamPVPKQ(pJoiner) || !h.so_smo_TargetIsEnemyKQTeam(pJoiner, pEnemy) {
		t.Fatalf("expected TeamPVP enemy true for Red(0) vs Blue(1)")
	}
	if h.so_smo_TargetIsEnemyKQTeam(pJoiner, pJoiner) {
		t.Fatalf("expected same-team enemy check == false")
	}

	// KnockDown++ twice -> kd == 2 >= maxRevives(2) -> sp_EscapeFromEmptyKQ warps Joiner out via sp_2Roumen (RouN)
	h.onPlayerKQKnockDown(pJoiner)
	h.onPlayerKQKnockDown(pJoiner)
	if kd := h.sp_KQElementFunction(pJoiner, (*KQPlayerInform).kqpi_KnockDownGet); kd != 2 {
		t.Fatalf("expected knockdown count 2, got %d", kd)
	}
	h.sp_EscapeFromEmptyKQ(pJoiner)
	if pJoiner.CurrentMap.MapID != "RouN" {
		t.Fatalf("expected exhausted joiner escaped to RouN, got %s", pJoiner.CurrentMap.MapID)
	}

	// Non-joiner normal player -> escaped to RouN
	h.sp_EscapeFromEmptyKQ(pIntruder)
	if pIntruder.CurrentMap.MapID != "RouN" {
		t.Fatalf("expected non-joiner escaped to RouN, got %s", pIntruder.CurrentMap.MapID)
	}

	// Non-joiner GM -> stays and receives "Doing KQ" notice
	connGM.clear()
	h.sp_EscapeFromEmptyKQ(pGM)
	if pGM.CurrentMap.MapID != "KDEddyHill07" {
		t.Fatalf("expected GM to remain on KQ map, got %s", pGM.CurrentMap.MapID)
	}
	if len(connGM.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) == 0 {
		t.Fatalf("expected 'Doing KQ' notice for GM observer")
	}

	// Unbound KQ map with no active film -> escapes player
	emptyKQ := h.mapManager.EnsureKQMap("KDEddyHill08", "KDEddyHill")
	emptyKQ.KQHandle = kqInvalidHandle
	emptyKQ.KQScriptActive = false
	pEmpty, _ := newTestClientAndPlayer(t, repo, ds, emptyKQ, h, 7205, "LostInKQ", 30, 1500, 1500)
	h.sp_EscapeFromEmptyKQ(pEmpty)
	if pEmpty.CurrentMap.MapID != "RouN" {
		t.Fatalf("expected player in empty KQ without film escaped to RouN, got %s", pEmpty.CurrentMap.MapID)
	}
	_ = zMap
}

func TestKQZone_KQItemDrop_RemoveKQItem_And_PickupLimit(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 7301, "ItemHero", 30, 1500, 1500)

	// 1. Find a real item from KQItem.shn in DataStore.KQItemBox
	var kqDropItemID uint16
	var pickupLimit uint16
	for _, kqIt := range ds.KQStore.Items {
		if info, ok := ds.ItemStore.GetItemInfoByName(kqIt.ItemIndexString()); ok && info != nil {
			kqDropItemID = info.ID
			pickupLimit = kqIt.PickupLimit
			// Ensure flags reflect ItemInfoServer.shn or verify them
			break
		}
	}
	if kqDropItemID == 0 {
		t.Fatalf("expected at least one KQItem.shn entry resolved in ItemInfo.shn")
	}
	if _, ok := ds.KQItemBox.kqidb_GetData(kqDropItemID); !ok {
		t.Fatalf("expected kqidb_GetData(%d) to succeed", kqDropItemID)
	}

	// Test kqidb_CanPickUp boundary
	if pickupLimit > 0 {
		if !ds.KQItemBox.kqidb_CanPickUp(kqDropItemID, uint64(pickupLimit)) {
			t.Fatalf("expected CanPickUp true at exact PickupLimit=%d", pickupLimit)
		}
		if ds.KQItemBox.kqidb_CanPickUp(kqDropItemID, uint64(pickupLimit)+1) {
			t.Fatalf("expected CanPickUp false above PickupLimit=%d", pickupLimit)
		}
		// Give player pickupLimit copies and try picking up 1 more -> PROTO_NC_ITEM_PICK_ACK err=0x034A
		_ = repo.CreateItem(&database.Item{
			Owner:       p.Char.ID,
			ItemID:      kqDropItemID,
			StorageType: 9,
			Storage:     0,
			Flags:       int32(pickupLimit),
		})
		drop := zMap.SpawnDroppedItem(kqDropItemID, 1, 1, 1505, 1505, 0, 0)
		conn.clear()
		pickBuf := make([]byte, 2)
		binary.LittleEndian.PutUint16(pickBuf, drop.Handle)
		h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_PICK_REQ, pickBuf))
		acks := conn.findByOpcode(protocol.PROTO_NC_ITEM_PICK_ACK)
		if len(acks) == 0 {
			t.Fatalf("expected PROTO_NC_ITEM_PICK_ACK")
		}
		errCode := binary.LittleEndian.Uint16(acks[0].Payload[6:8])
		if errCode != kqErrPickupLimit {
			t.Fatalf("expected KQ pickup limit err 0x%04X, got 0x%04X", kqErrPickupLimit, errCode)
		}
	}

	// 2. Test sp_KQItemDrop (drops items with ItemInfoServer.KQ_Item_Drop != 0)
	var dropFlagItemID uint16
	var kqOnlyItemID uint16
	for id := uint16(1); id < 65000; id++ {
		info, ok := ds.ItemStore.GetItemInfo(id)
		if !ok || info == nil {
			continue
		}
		if dropFlagItemID == 0 && info.KQItemDrop {
			dropFlagItemID = id
		}
		if kqOnlyItemID == 0 && info.KQItem && !info.KQItemDrop {
			kqOnlyItemID = id
		}
		if dropFlagItemID != 0 && kqOnlyItemID != 0 {
			break
		}
	}
	if dropFlagItemID == 0 || kqOnlyItemID == 0 {
		t.Fatalf("expected real ItemInfoServer.shn items with KQ_Item_Drop=%d and KQItemOnly=%d", dropFlagItemID, kqOnlyItemID)
	}

	p2, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 7302, "DropHero", 30, 1600, 1600)
	_ = repo.CreateItem(&database.Item{Owner: p2.Char.ID, ItemID: dropFlagItemID, StorageType: 9, Storage: 1, Flags: 1})
	_ = repo.CreateItem(&database.Item{Owner: p2.Char.ID, ItemID: kqOnlyItemID, StorageType: 9, Storage: 2, Flags: 1})

	droppedCnt := h.sp_KQItemDrop(p2)
	if droppedCnt != 1 {
		t.Fatalf("expected sp_KQItemDrop to drop 1 KQ_Item_Drop item, dropped=%d", droppedCnt)
	}
	// Remaining item on non-KQ map is kqOnlyItemID -> sp_RemoveKQItem(9) removes it on non-KQ map!
	removedCnt := h.sp_RemoveKQItem(p2, kqRemoveBoxInven)
	if removedCnt != 1 {
		t.Fatalf("expected sp_RemoveKQItem(9) on non-KQ map to remove KQItem, removed=%d", removedCnt)
	}
}

func TestKQZone_Destroy_Disjoin_Ban_And_GMCommands(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	// 1. Make + Start KQ handle=500 on KDEddyHill09
	info := makeTestKQInfo(500, 1, "KDEddyHill", "KDEddyHill09", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	kqMap := h.mapManager.GetMap("KDEddyHill09")
	pBanned, _ := newTestClientAndPlayer(t, repo, ds, kqMap, h, 7401, "BannedGuy", 30, 1500, 1500)
	pOutside, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 7402, "OutsideGuy", 30, 1500, 1500)
	pGM, connGM := newTestClientAndPlayer(t, repo, ds, zMap, h, 7403, "KQGM", 99, 1500, 1500)
	pGM.Char.AdminLevel = 100

	elem, _, _ := h.kqState().container.kqc_SearchHandle(500)
	elem.PlayerInfoList.kqpil_InsertPlayerInfo(uint32(pBanned.Char.ID), 0, 1)

	// 2. wms_NC_KQ_PLAYER_DISJOIN_CMD -> deletes player from KQPlayerInfoList
	h.ApplyIncomingWMS("wms_NC_KQ_PLAYER_DISJOIN_CMD", 0, 0, map[string]interface{}{
		"handle":       uint32(500),
		"char_reg_num": uint32(pBanned.Char.ID),
	})
	if elem.PlayerInfoList.kqpil_IsKQJoiner(uint32(pBanned.Char.ID)) != 0 {
		t.Fatalf("expected player removed by wms_NC_KQ_PLAYER_DISJOIN_CMD")
	}

	// 3. wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD -> warps player on KDEddyHill09 to RouN (6445, 8630), ignores player outside
	h.ApplyIncomingWMS("wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD", 0, 0, map[string]interface{}{
		"handle":       uint32(500),
		"char_reg_num": uint32(pBanned.Char.ID),
		"out_maps":     [4]string{"KDEddyHill09", "", "", ""},
	})
	if pBanned.CurrentMap.MapID != kqRoumenFallbackMap || pBanned.Char.XPos != kqRoumenFallbackX || pBanned.Char.YPos != kqRoumenFallbackY {
		t.Fatalf("expected banned player force-warped to RouN (6445,8630), got %s (%d,%d)",
			pBanned.CurrentMap.MapID, pBanned.Char.XPos, pBanned.Char.YPos)
	}
	h.ApplyIncomingWMS("wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD", 0, 0, map[string]interface{}{
		"handle":       uint32(500),
		"char_reg_num": uint32(pOutside.Char.ID),
		"out_maps":     [4]string{"KDEddyHill09", "", "", ""},
	})
	if pOutside.CurrentMap.MapID != zMap.MapID {
		t.Fatalf("expected player outside banned KQ map NOT to be warped")
	}

	// 4. wms_NC_KQ_W2Z_DESTROY_CMD -> unbinds map, drops film, frees container slot
	kqMap.cc_PlayFilm("KQ/KingSlime/KingSlime", "1", time.Now())
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_DESTROY_CMD", 0, 0, map[string]interface{}{
		"handle": uint32(500),
	})
	if kqMap.fm_GetKQhandle() != kqInvalidHandle || kqMap.KQID != kqInvalidID || kqMap.KQScriptActive {
		t.Fatalf("expected KDEddyHill09 unbound & film dropped after DESTROY_CMD")
	}
	if _, _, stillThere := h.kqState().container.kqc_SearchHandle(500); stillThere {
		t.Fatalf("expected KQElement 500 freed from container")
	}
	if !h.kqState().container.kqc_MapUseCheck("KDEddyHill09") {
		t.Fatalf("expected KDEddyHill09 free again after DESTROY_CMD")
	}

	// 5. GM commands: &kingdomquest, &Reroadkingdomquest, &endofkq, &kqreward
	if !h.handleGMCommand(pGM, "&kingdomquest KDEddyHill KQ/KingSlime/KingSlime 1") {
		t.Fatalf("expected &kingdomquest handled")
	}
	if pGM.CurrentMap.MapID != "KDEddyHill" || !pGM.CurrentMap.KQScriptActive {
		t.Fatalf("expected &kingdomquest to warp GM to KDEddyHill and start film, got map=%s active=%v",
			pGM.CurrentMap.MapID, pGM.CurrentMap.KQScriptActive)
	}

	if !h.handleGMCommand(pGM, "&Reroadkingdomquest KQ/KingSlime/KingSlime") {
		t.Fatalf("expected &Reroadkingdomquest handled")
	}
	if pGM.CurrentMap.KQScriptActive {
		t.Fatalf("expected &Reroadkingdomquest to drop running film via cm_AllDropFilm")
	}

	clearWMSBridge(bridge)
	if !h.handleGMCommand(pGM, "&endofkq 777") {
		t.Fatalf("expected &endofkq handled")
	}
	evEnd, ok := lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_END_CMD")
	if !ok {
		t.Fatalf("expected &endofkq to publish wms_NC_KQ_Z2W_END_CMD")
	}
	if endH, _ := payloadInt(evEnd.Payload, "handle"); endH != 777 {
		t.Fatalf("expected Z2W_END_CMD handle=777, got %d", endH)
	}

	if !h.handleGMCommand(pGM, "&kqreward 1") {
		t.Fatalf("expected &kqreward handled")
	}
	if pGM.KQRewardIdx == "" {
		t.Fatalf("expected &kqreward to invoke so_ply_KQRewardStruct and populate KQRewardIdx")
	}

	// Invalid map name in &kingdomquest sends "Invalid map name" notice
	connGM.clear()
	h.handleGMCommand(pGM, "&kingdomquest NoSuchMap999 KQ/KingSlime/KingSlime 1")
	if len(connGM.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) == 0 {
		t.Fatalf("expected 'Invalid map name' notice for invalid &kingdomquest map")
	}
}
