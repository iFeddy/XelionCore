package zone

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// faultInjectingCharRepo wraps MemoryCharacterRepository to test Character.exe fc_NC_KQ_REWARD_REQ
// 3-stage rollback semantics (0x1101 / 0x1102 / 0x1103, exe: 0x004321D0).
type faultInjectingCharRepo struct {
	*database.MemoryCharacterRepository
	failMoneyAdd bool
	failFameAdd  bool
	subMoneyCall int
	deletedItems []int64
}

func (f *faultInjectingCharRepo) AddCharMoney(charNo int, amount int64) error {
	if f.failMoneyAdd {
		return errors.New("injected p_Char_MoneyAdd failure")
	}
	return f.MemoryCharacterRepository.AddCharMoney(charNo, amount)
}

func (f *faultInjectingCharRepo) SubCharMoney(charNo int, amount int64) (bool, error) {
	f.subMoneyCall++
	return f.MemoryCharacterRepository.SubCharMoney(charNo, amount)
}

func (f *faultInjectingCharRepo) AddCharFame(charNo int, amount int) error {
	if f.failFameAdd {
		return errors.New("injected p_Char_FameAdd failure")
	}
	return f.MemoryCharacterRepository.AddCharFame(charNo, amount)
}

func (f *faultInjectingCharRepo) DeleteItem(itemKey int64) error {
	f.deletedItems = append(f.deletedItems, itemKey)
	return f.MemoryCharacterRepository.DeleteItem(itemKey)
}

func TestKQRewards_AxialListAndShineRewardPipeline(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	if ds.KQStore == nil || len(ds.KQStore.ShineRewards) != 435 {
		t.Fatalf("expected 435 ShineReward rows loaded in DataStore.KQStore")
	}

	// Pick a KingdomQuestRew row that grants items/exp/cen/fame (row 0 has 8 ShineReward slots)
	rewRow := ds.KQStore.Rewards[0]

	pLowKills, connLow := newTestClientAndPlayer(t, repo, ds, zMap, h, 8101, "LowKills", 30, 1500, 1500)
	pHighKills, connHigh := newTestClientAndPlayer(t, repo, ds, zMap, h, 8102, "HighKills", 30, 1600, 1600)

	// Setup active KQ via W2Z_MAKE_REQ with DemandMobKill = 3, RewardIndex = rewRow.ID
	info := makeTestKQInfo(3001, 1, "KDEddyHill", "KDEddyHill00", "KQ/KingSlime/KingSlime")
	info.DemandMobKill = 3
	info.RewardIndex = uint16(rewRow.ID)
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	kqMap := h.mapManager.GetMap("KDEddyHill00")
	if kqMap == nil {
		t.Fatalf("expected KDEddyHill00 created by W2Z_MAKE_REQ")
	}

	// Move both players into kqMap and register them in KQPlayerInfoList
	zMap.RemovePlayer(pLowKills)
	zMap.RemovePlayer(pHighKills)
	kqMap.AddPlayer(pLowKills)
	kqMap.AddPlayer(pHighKills)
	pLowKills.CurrentMap = kqMap
	pHighKills.CurrentMap = kqMap
	if elem, _, ok := h.kqState().container.kqc_SearchHandle(info.Handle); ok && elem != nil {
		elem.PlayerInfoList.kqpil_InsertPlayerInfo(uint32(pLowKills.Char.ID), 0, 1)
		elem.PlayerInfoList.kqpil_InsertPlayerInfo(uint32(pHighKills.Char.ID), 1, 1)
	}

	// Record mob kills via KQContributeList::kqcl_SendMobKill (exe: 0x00499FF0)
	h.RecordKQMobKill(kqMap, pLowKills) // 1 < 3 -> ineligible
	for i := 0; i < 3; i++ {
		h.RecordKQMobKill(kqMap, pHighKills) // 3 >= 3 -> eligible
	}
	if got := h.kqState().contribute.kqcl_GetMobKill(info.Handle, uint32(pHighKills.Char.ID)); got != 3 {
		t.Fatalf("expected 3 mob kills for HighKills, got %d", got)
	}

	// Deterministic RNG (always roll 0 < RewardRate so every configured ShineReward triggers)
	h.SetKQRewardRand(func(n int) int { return 0 })

	initExp := pHighKills.Char.TCharacterState.Exp
	initMoney := pHighKills.Char.TCharacterState.Money
	initFame := pHighKills.Char.TCharacterState.Fame
	initInv, _ := repo.GetCharacterItems(pHighKills.Char.ID)

	connLow.clear()
	connHigh.clear()

	// Execute cReward(Field, "KQ") -> AxialListKQRewardLua (exe: 0x005DF060 -> 0x004288C0 / 0x00428490)
	if !kqMap.bridgeRewardAll("KQ") {
		t.Fatalf("bridgeRewardAll('KQ') returned false")
	}

	// 1. LowKills (< DemandMobKill) MUST receive PROTO_NC_KQ_NOREWARD_CMD (0x5823) with 0x1104
	noRewPkts := connLow.findByOpcode(protocol.PROTO_NC_KQ_NOREWARD_CMD)
	if len(noRewPkts) != 1 {
		t.Fatalf("expected 1 PROTO_NC_KQ_NOREWARD_CMD for LowKills, got %d", len(noRewPkts))
	}
	if len(noRewPkts[0].Payload) != 2 {
		t.Fatalf("expected 2-byte payload for PROTO_NC_KQ_NOREWARD_CMD, got %d", len(noRewPkts[0].Payload))
	}
	if errCode := binary.LittleEndian.Uint16(noRewPkts[0].Payload[0:2]); errCode != kqErrRewardNoMobKill {
		t.Fatalf("expected 0x%04X (KQ_NOREWARD), got 0x%04X", kqErrRewardNoMobKill, errCode)
	}

	// 2. HighKills (>= DemandMobKill) MUST NOT get NOREWARD and MUST receive real rewards from ShineReward.shn
	if len(connHigh.findByOpcode(protocol.PROTO_NC_KQ_NOREWARD_CMD)) != 0 {
		t.Fatalf("eligible player HighKills must not receive PROTO_NC_KQ_NOREWARD_CMD")
	}
	afterInv, _ := repo.GetCharacterItems(pHighKills.Char.ID)
	rewardedSomething := pHighKills.Char.TCharacterState.Exp > initExp ||
		pHighKills.Char.TCharacterState.Money > initMoney ||
		pHighKills.Char.TCharacterState.Fame > initFame ||
		len(afterInv) > len(initInv)
	if !rewardedSomething {
		t.Fatalf("expected HighKills to receive EXP/Money/Fame/Items from KingdomQuestRew[0] -> ShineReward.shn")
	}

	// Verify pending reward order was cleared by gds_NC_KQ_REWARDSUC_ACK (exe: 0x0052E570)
	if len(h.kqRewardTracker().pending) != 0 {
		t.Fatalf("expected pending reward orders drained after gds_NC_KQ_REWARDSUC_ACK, got %d", len(h.kqRewardTracker().pending))
	}

	// Verify cKQRewardIndex (exe: 0x005E4000 -> so_ply_KQRewardIndex 0x0052EA80) resolves by IndexName
	namedIdx := ""
	for i := range ds.KQStore.Rewards {
		if n := ds.KQStore.Rewards[i].IndexName(); n != "" {
			namedIdx = n
			break
		}
	}
	if namedIdx == "" {
		t.Fatalf("expected at least one named row in KingdomQuestRew.shn")
	}
	beforeMoney := pLowKills.Char.TCharacterState.Money
	beforeFame := pLowKills.Char.TCharacterState.Fame
	beforeExp := pLowKills.Char.TCharacterState.Exp
	beforeInvLow, _ := repo.GetCharacterItems(pLowKills.Char.ID)
	kqMap.bridgeKQRewardIndex(uint32(pLowKills.MapObjectID), namedIdx)
	afterInvLow, _ := repo.GetCharacterItems(pLowKills.Char.ID)
	if pLowKills.Char.TCharacterState.Exp == beforeExp &&
		pLowKills.Char.TCharacterState.Money == beforeMoney &&
		pLowKills.Char.TCharacterState.Fame == beforeFame &&
		len(afterInvLow) == len(beforeInvLow) {
		t.Fatalf("expected cKQRewardIndex(%q) to grant reward via sp_KQReward", namedIdx)
	}

	// If money or fame was granted, GameLog should have recorded 0x0519 / 0x051C
	if pHighKills.Char.TCharacterState.Money > initMoney {
		if gl := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
			return e.CharNo == pHighKills.Char.ID && e.Type == 0x0519
		}); gl == nil {
			t.Fatalf("expected GameLog 0x0519 for KQ reward money")
		}
	}
}

func TestKQRewards_CharacterExeRollbackStages(t *testing.T) {
	h, ds, baseRepo, zMap, _ := setupFase14Handler(t)
	faultRepo := &faultInjectingCharRepo{MemoryCharacterRepository: baseRepo}
	h.charRepo = faultRepo

	player, _ := newTestClientAndPlayer(t, baseRepo, ds, zMap, h, 8201, "RollbackHero", 40, 1500, 1500)
	player.Char.TCharacterState.Money = 1000
	player.Char.TCharacterState.Fame = 50
	_ = baseRepo.UpdateCharacter(player.Char)

	// Resolve a valid item ID from ItemInfo.shn
	var validItemID uint16 = 1
	if ds.ItemStore != nil {
		if info, ok := ds.ItemStore.GetItemInfoByName("Elrue1"); ok && info != nil {
			validItemID = info.ID
		}
	}

	// Stage 2 failure: p_Char_MoneyAdd fails -> error 0x1102 -> rollback created items
	faultRepo.failMoneyAdd = true
	faultRepo.failFameAdd = false
	faultRepo.deletedItems = nil
	order1 := &KQRewardOrder{
		Handle:     player.MapObjectID,
		CharNo:     uint32(player.Char.ID),
		OrderIndex: 1,
		Money:      500,
		Fame:       10,
		Items:      []KQRewardItemSpec{{ItemID: validItemID, Quantity: 1}},
	}
	h.kqRewardTracker().mu.Lock()
	h.kqRewardTracker().pending[kqRewardOrderKey(order1.CharNo, order1.OrderIndex)] = order1
	h.kqRewardTracker().mu.Unlock()

	h.fc_NC_KQ_REWARD_REQ(player, order1)
	if len(faultRepo.deletedItems) != 1 {
		t.Fatalf("Stage 2 (0x1102): expected 1 created item rolled back, got %d", len(faultRepo.deletedItems))
	}
	if player.Char.TCharacterState.Money != 1000 || player.Char.TCharacterState.Fame != 50 {
		t.Fatalf("Stage 2 (0x1102): money/fame should remain unchanged (1000/50), got %d/%d",
			player.Char.TCharacterState.Money, player.Char.TCharacterState.Fame)
	}
	if len(h.kqRewardTracker().pending) != 0 {
		t.Fatalf("Stage 2 (0x1102): expected gds_NC_KQ_REWARDFAIL_ACK to remove pending order")
	}

	// Stage 3 failure: p_Char_FameAdd fails -> error 0x1103 -> rollback money (SubCharMoney) AND items
	faultRepo.failMoneyAdd = false
	faultRepo.failFameAdd = true
	faultRepo.subMoneyCall = 0
	faultRepo.deletedItems = nil
	order2 := &KQRewardOrder{
		Handle:     player.MapObjectID,
		CharNo:     uint32(player.Char.ID),
		OrderIndex: 2,
		Money:      750,
		Fame:       25,
		Items:      []KQRewardItemSpec{{ItemID: validItemID, Quantity: 1}},
	}
	h.kqRewardTracker().mu.Lock()
	h.kqRewardTracker().pending[kqRewardOrderKey(order2.CharNo, order2.OrderIndex)] = order2
	h.kqRewardTracker().mu.Unlock()

	h.fc_NC_KQ_REWARD_REQ(player, order2)
	if faultRepo.subMoneyCall != 1 {
		t.Fatalf("Stage 3 (0x1103): expected p_Char_MoneySub rollback called once, got %d", faultRepo.subMoneyCall)
	}
	if len(faultRepo.deletedItems) != 1 {
		t.Fatalf("Stage 3 (0x1103): expected 1 created item rolled back, got %d", len(faultRepo.deletedItems))
	}
	if ch, _ := baseRepo.GetCharacterByID(player.Char.ID); ch != nil && ch.TCharacterState.Money != 1000 {
		t.Fatalf("Stage 3 (0x1103): expected repo money restored to 1000, got %d", ch.TCharacterState.Money)
	}
	if len(h.kqRewardTracker().pending) != 0 {
		t.Fatalf("Stage 3 (0x1103): expected gds_NC_KQ_REWARDFAIL_ACK to remove pending order")
	}
}

func TestKQEnd_AxialListEndAndShineEndOfKingdomQuest(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	p1, conn1 := newTestClientAndPlayer(t, repo, ds, zMap, h, 8301, "EndHero1", 35, 1500, 1500)
	info := makeTestKQInfo(3101, 5, "KDEddyHill", "KDEddyHill01", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	kqMap := h.mapManager.GetMap("KDEddyHill01")
	if kqMap == nil {
		t.Fatalf("expected KDEddyHill01 created")
	}
	zMap.RemovePlayer(p1)
	kqMap.AddPlayer(p1)
	p1.CurrentMap = kqMap
	_ = kqMap.SpawnMobFromTemplate(&MobTemplate{ID: 1, Name: "Slime", Level: 5, MaxHP: 100}, 1500, 1500, 0, "", 0, true)

	// 1. AxialListKQEndLua("Success") -> 0x5812 (empty payload) + GameLog 0x0442 (Leave_Status=1) + fm_ClearObject(0x1B0)
	conn1.clear()
	if !kqMap.bridgeQuestResult(kqMap.MapID, "Success") {
		t.Fatalf("bridgeQuestResult(Success) failed")
	}
	compPkts := conn1.findByOpcode(protocol.PROTO_NC_KQ_COMPLETE_CMD)
	if len(compPkts) != 1 || len(compPkts[0].Payload) != 0 {
		t.Fatalf("expected 1 empty-payload PROTO_NC_KQ_COMPLETE_CMD (0x5812), got %+v", compPkts)
	}
	if gl := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == p1.Char.ID && e.Type == int(protocol.PROTO_NC_LOG_GAME_KQ_LEAVE) && e.Int3 == 1
	}); gl == nil {
		t.Fatalf("expected GameLog PROTO_NC_LOG_GAME_KQ_LEAVE (0x0442, Leave_Status=1)")
	}
	if len(kqMap.mobs) != 0 {
		t.Fatalf("expected fm_ClearObject(0x1B0) to clear mobs on cQuestResult, remaining=%d", len(kqMap.mobs))
	}

	// 2. AxialListKQEndLua("Fail") -> 0x5813 (empty payload) + GameLog 0x0442 (Leave_Status=0, exe: 0x00429EA7 sete al)
	conn1.clear()
	if !kqMap.bridgeQuestResult(kqMap.MapID, "Fail") {
		t.Fatalf("bridgeQuestResult(Fail) failed")
	}
	failPkts := conn1.findByOpcode(protocol.PROTO_NC_KQ_FAIL_CMD)
	if len(failPkts) != 1 || len(failPkts[0].Payload) != 0 {
		t.Fatalf("expected 1 empty-payload PROTO_NC_KQ_FAIL_CMD (0x5813), got %+v", failPkts)
	}
	if gl := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == p1.Char.ID && e.Type == int(protocol.PROTO_NC_LOG_GAME_KQ_LEAVE) && e.Int3 == 0
	}); gl == nil {
		t.Fatalf("expected GameLog PROTO_NC_LOG_GAME_KQ_LEAVE (0x0442, Leave_Status=0)")
	}

	// 3. ShineEndOfKingdomQuest / cEndOfKingdomQuest (exe: 0x004F5FC0 / 0x005DF1F0)
	clearWMSBridge(bridge)
	if !kqMap.bridgeEndOfKingdomQuest(true) {
		t.Fatalf("bridgeEndOfKingdomQuest returned false")
	}
	ev, ok := lastWMSEventBySymbol(bridge, "wms_NC_KQ_Z2W_END_CMD")
	if !ok {
		t.Fatalf("expected wms_NC_KQ_Z2W_END_CMD (0x5810) emitted by wms_EndOfKQPacket")
	}
	if gotH, _ := payloadInt(ev.Payload, "handle"); uint32(gotH) != info.Handle {
		t.Fatalf("wms_NC_KQ_Z2W_END_CMD handle = %d, want %d", gotH, info.Handle)
	}
}

func TestKQLua_TimersTeamBattleScoreBoardAndRealScriptPlayFilm(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	p1, conn1 := newTestClientAndPlayer(t, repo, ds, zMap, h, 8401, "ArenaHero", 50, 1500, 1500)
	info := makeTestKQInfo(3201, 2, "KDEddyHill", "KDEddyHill02", "KQ/KingSlime/KingSlime")
	info.LimitTime = 25
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))

	kqMap := h.mapManager.GetMap("KDEddyHill02")
	zMap.RemovePlayer(p1)
	kqMap.AddPlayer(p1)
	p1.CurrentMap = kqMap
	p1.KQTeam = 1 // BLUE

	wb := NewZoneWorldBridge(kqMap)

	// cGetKQLimitSecond -> 25 * 60 = 1500
	if sec, ok := wb.GetKQLimitSecondExact("KDEddyHill02"); !ok || sec != 1500 {
		t.Fatalf("GetKQLimitSecondExact = (%d, %v), want (1500, true)", sec, ok)
	}
	// cGetKQTeamType -> 1 for p1, 2 (KQTT_MAX) for unknown
	if got := wb.GetKQTeamTypeExact(uint32(p1.MapObjectID)); got != 1 {
		t.Fatalf("GetKQTeamTypeExact(p1) = %d, want 1", got)
	}
	if got := wb.GetKQTeamTypeExact(0); got != 2 {
		t.Fatalf("GetKQTeamTypeExact(0) = %d, want 2 (KQTT_MAX)", got)
	}

	// cShowKQTimerWithLife / cVanishTimer -> PROTO_NC_MISC_RESTMINUTE_CMD (0x0810: flag=1, uint16 sec)
	conn1.clear()
	wb.ShowKQTimerWithLifeMap("KDEddyHill02", 600)
	wb.VanishTimerMap("KDEddyHill02")
	restPkts := conn1.findByOpcode(protocol.PROTO_NC_MISC_RESTMINUTE_CMD)
	if len(restPkts) != 2 {
		t.Fatalf("expected 2 PROTO_NC_MISC_RESTMINUTE_CMD (0x0810) packets, got %d", len(restPkts))
	}
	if len(restPkts[0].Payload) != 3 || restPkts[0].Payload[0] != 1 || binary.LittleEndian.Uint16(restPkts[0].Payload[1:3]) != 600 {
		t.Fatalf("ShowKQTimerWithLife payload mismatch: %x", restPkts[0].Payload)
	}
	if len(restPkts[1].Payload) != 3 || restPkts[1].Payload[0] != 1 || binary.LittleEndian.Uint16(restPkts[1].Payload[1:3]) != 0 {
		t.Fatalf("VanishTimer payload mismatch: %x", restPkts[1].Payload)
	}

	// cTimer -> PROTO_NC_SCENARIO_TIMER_CMD (0x6C16: uint16 sec)
	conn1.clear()
	wb.TimerMap("KDEddyHill02", 90)
	tPkts := conn1.findByOpcode(protocol.PROTO_NC_SCENARIO_TIMER_CMD)
	if len(tPkts) != 1 || len(tPkts[0].Payload) != 2 || binary.LittleEndian.Uint16(tPkts[0].Payload[0:2]) != 90 {
		t.Fatalf("TimerMap payload mismatch: %+v", tPkts)
	}

	// cSetTeamBattle -> 0x1828 (ON) / 0x1829 (OFF) + elem.Info.IsTeamPVP
	conn1.clear()
	wb.SetTeamBattleMap("KDEddyHill02", true)
	elem, _, _ := h.kqState().container.kqc_SearchHandle(info.Handle)
	if elem == nil || elem.Info.IsTeamPVP != 1 {
		t.Fatalf("expected elem.Info.IsTeamPVP == 1 after SetTeamBattle(true)")
	}
	if len(conn1.findByOpcode(protocol.PROTO_NC_MAP_KQTEAMBATTLEZONE_ON_CMD)) != 1 {
		t.Fatalf("expected PROTO_NC_MAP_KQTEAMBATTLEZONE_ON_CMD (0x1828)")
	}
	wb.SetTeamBattleMap("KDEddyHill02", false)
	if elem.Info.IsTeamPVP != 0 || len(conn1.findByOpcode(protocol.PROTO_NC_MAP_KQTEAMBATTLEZONE_OFF_CMD)) != 1 {
		t.Fatalf("expected IsTeamPVP=0 and PROTO_NC_MAP_KQTEAMBATTLEZONE_OFF_CMD (0x1829)")
	}

	// cScoreBoard_AllInMap -> 0x583C (6 bytes payload, clamped win<=5, score<=99)
	conn1.clear()
	wb.ScoreBoardAllInMap("KDEddyHill02", true, 3, 9, 150, 8, 200)
	sbPkts := conn1.findByOpcode(protocol.PROTO_NC_KQ_SCORE_BOARD_INFO_CMD)
	if len(sbPkts) != 1 || len(sbPkts[0].Payload) != 6 {
		t.Fatalf("expected 1 6-byte PROTO_NC_KQ_SCORE_BOARD_INFO_CMD, got %+v", sbPkts)
	}
	wantSB := []byte{1, 3, 5, 99, 5, 99}
	for i := range wantSB {
		if sbPkts[0].Payload[i] != wantSB[i] {
			t.Fatalf("ScoreBoard byte[%d] = %d, want %d", i, sbPkts[0].Payload[i], wantSB[i])
		}
	}

	// cScoreInfo_AllInMap -> 0x5836 (8 bytes payload: 2 x uint32)
	conn1.clear()
	wb.ScoreInfoAllInMap("KDEddyHill02", []uint32{42, 77})
	siPkts := conn1.findByOpcode(protocol.PROTO_NC_KQ_SCORE_INFO_CMD)
	if len(siPkts) != 1 || len(siPkts[0].Payload) != 8 {
		t.Fatalf("expected 1 8-byte PROTO_NC_KQ_SCORE_INFO_CMD, got %+v", siPkts)
	}
	if binary.LittleEndian.Uint32(siPkts[0].Payload[0:4]) != 42 || binary.LittleEndian.Uint32(siPkts[0].Payload[4:8]) != 77 {
		t.Fatalf("ScoreInfo payload mismatch: %x", siPkts[0].Payload)
	}

	// Real cc_PlayFilm -> StartKQScript + TickKQScript + cc_DropFilm for standard and PineEventScript-migrated KQ
	now := time.Now()
	kqMap.cc_PlayFilm("KQ/KingSlime/KingSlime", "1", now)
	if !kqMap.KQScriptActive || kqMap.luaEngine == nil {
		t.Fatalf("expected KingSlime Lua script active after cc_PlayFilm")
	}
	kqMap.TickKQScript()
	kqMap.cc_DropFilm("KQ/KingSlime/KingSlime")
	if kqMap.KQScriptActive || kqMap.luaEngine != nil {
		t.Fatalf("expected KingSlime Lua script stopped after cc_DropFilm")
	}

	kqMap.cc_PlayFilm("KQ/ConditionOfHero/ConditionOfHero", "0", now)
	if !kqMap.KQScriptActive || kqMap.luaEngine == nil {
		t.Fatalf("expected ConditionOfHero Lua script active after cc_PlayFilm")
	}
	kqMap.cc_DropFilm("KQ/ConditionOfHero/ConditionOfHero")
}

func TestKQSaveLocation_PreKQReturnPosAndGameLogEnter(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	player, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 8501, "SaveLocHero", 30, 5100, 6200)

	// Make a KQ map and move player into it, saving pre-KQ return position
	info := makeTestKQInfo(3301, 7, "KDEddyHill", "KDEddyHill03", "KQ/KingSlime/KingSlime")
	h.ApplyIncomingWMS("wms_NC_KQ_W2Z_MAKE_REQ", 0, 0, makeW2ZMakePayload(t, info))
	kqMap := h.mapManager.GetMap("KDEddyHill03")
	if kqMap == nil {
		t.Fatalf("expected KDEddyHill03 created")
	}

	h.SavePreKQReturnPos(player, "Rou", 5100, 6200)
	zMap.RemovePlayer(player)
	kqMap.AddPlayer(player)
	player.CurrentMap = kqMap
	player.Char.XPos = 2400
	player.Char.YPos = 3500

	// 1. so_SaveLocation (exe: 0x004542F0) persists pre-KQ loginMap ("Rou", 5100, 6200) + KQ handle/map/coords
	loginMap, lx, ly, kqH, kqM, kqX, kqY, ok := h.so_SaveLocation(player)
	if !ok {
		t.Fatalf("so_SaveLocation returned ok=false")
	}
	if loginMap != "Rou" || lx != 5100 || ly != 6200 {
		t.Fatalf("expected pre-KQ loginMap=(Rou, 5100, 6200), got (%s, %d, %d)", loginMap, lx, ly)
	}
	if kqH != info.Handle || kqM != "KDEddyHill03" || kqX != 2400 || kqY != 3500 {
		t.Fatalf("expected KQ saved=(%d, KDEddyHill03, 2400, 3500), got (%d, %s, %d, %d)",
			info.Handle, kqH, kqM, kqX, kqY)
	}

	// 2. Verify p_Char_GetKQMap (exe: Character.exe @ 0x0045A130) reads back exact KQMap fields
	kqInfo, err := repo.GetCharKQMap(player.Char.ID)
	if err != nil || kqInfo == nil {
		t.Fatalf("GetCharKQMap failed: %v", err)
	}
	if kqInfo.KQHandle != info.Handle || kqInfo.KQMapName != "KDEddyHill03" || kqInfo.KQX != 2400 || kqInfo.KQY != 3500 {
		t.Fatalf("GetCharKQMap mismatch: %+v", kqInfo)
	}
	if kqInfo.KQDatePacked == 0 {
		t.Fatalf("expected non-zero packed dKQDate from SaveCharLocationWithKQ")
	}

	// 3. EmitKQEnterGameLog (exe: 0x00496740 -> PROTO_NC_LOG_GAME_KQ_ENTER 0x0441)
	h.EmitKQEnterGameLog(player, "Rou", 5100, 6200, info.ID, info.Handle)
	if gl := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == player.Char.ID &&
			e.Type == int(protocol.PROTO_NC_LOG_GAME_KQ_ENTER) &&
			e.Int1 == int(info.ID) &&
			e.Int2 == int(info.Handle)
	}); gl == nil {
		t.Fatalf("expected GameLog PROTO_NC_LOG_GAME_KQ_ENTER (0x0441)")
	}
}
