package zone

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func waitForGameLogType(t *testing.T, glRepo *memGameLogRepo, match func(*database.GameLogEntry) bool) *database.GameLogEntry {
	t.Helper()
	deadline := time.Now().Add(400 * time.Millisecond)
	for time.Now().Before(deadline) {
		glRepo.mu.Lock()
		for i := len(glRepo.entries) - 1; i >= 0; i-- {
			if match(glRepo.entries[i]) {
				e := glRepo.entries[i]
				glRepo.mu.Unlock()
				return e
			}
		}
		glRepo.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func TestFase15_CatalogsLoadedFromOfficialData(t *testing.T) {
	_, ds, _, _, _ := setupFase14Handler(t)

	if ds.QuestCatalog == nil || len(ds.QuestCatalog.ByID) == 0 {
		t.Fatalf("expected non-empty QuestCatalog loaded from QuestData.shn")
	}
	if len(ds.QuestCatalog.Dialogs) == 0 {
		t.Fatalf("expected non-empty QuestDialog.shn entries")
	}
	if len(ds.QuestCatalog.SpeciesMobs) == 0 {
		t.Fatalf("expected non-empty QuestSpecies.shn entries")
	}
	if ds.KQCatalog == nil || len(ds.KQCatalog.ByID) == 0 {
		t.Fatalf("expected non-empty KQCatalog")
	}
	if ds.GTCatalog == nil || len(ds.GTCatalog.Skills) == 0 {
		t.Fatalf("expected non-empty GTCatalog from GuildTournament.json")
	}
	if ds.MiniHouseCatalog == nil || len(ds.MiniHouseCatalog.ByItemID) == 0 || len(ds.MiniHouseCatalog.Furniture) == 0 {
		t.Fatalf("expected non-empty MiniHouseCatalog from MiniHouse*.shn")
	}
	if ds.PetCatalog == nil || len(ds.PetCatalog.ActionsByInx) == 0 {
		t.Fatalf("expected non-empty PetCatalog from GTIServer.shn")
	}
	if ds.WeaponTitleCatalog == nil || len(ds.WeaponTitleCatalog.Rows) == 0 {
		t.Fatalf("expected non-empty WeaponTitleCatalog from WeaponTitleData.shn")
	}
	if ds.CasinoCatalog == nil || (len(ds.CasinoCatalog.SlotBetCoins) == 0 && len(ds.CasinoCatalog.DiceGames) == 0) {
		t.Fatalf("expected non-empty CasinoCatalog from casino SHNs")
	}
}

func TestFase15_QuestsNPCMenusScenariosAndTutorial(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9501, "QuestHero", 30, 1000, 1000)

	// Buscar determinísticamente la menor quest real en QuestData.shn que cumpla nivel 30 y tenga objetivo de caza válido
	var chosenQuest *QuestDef
	for id := uint16(1); id < 65000; id++ {
		q := ds.QuestCatalog.ByID[id]
		if q == nil || q.CheckPrevQuest {
			continue
		}
		if q.CheckLevel && (p.Char.Level < q.MinLevel || (q.MaxLevel > 0 && p.Char.Level > q.MaxLevel)) {
			continue
		}
		hasKill := false
		for _, mr := range q.MobKills {
			if mr.Enabled && mr.MobOrSpeciesID > 0 {
				if !mr.IsSpecies || len(ds.QuestCatalog.SpeciesMobs[mr.MobOrSpeciesID]) > 0 {
					hasKill = true
					break
				}
			}
		}
		if hasKill {
			chosenQuest = q
			break
		}
		if chosenQuest == nil {
			chosenQuest = q
		}
	}
	if chosenQuest == nil {
		for _, q := range ds.QuestCatalog.ByID {
			if q != nil {
				chosenQuest = q
				chosenQuest.CheckPrevQuest = false
				chosenQuest.CheckLevel = false
				break
			}
		}
	}
	if chosenQuest == nil {
		t.Fatalf("no quest definition found in QuestCatalog")
	}

	// 1. PROTO_NC_QUEST_START_REQ (Start quest)
	startBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(startBuf, chosenQuest.QuestID)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_START_REQ, startBuf))

	startAcks := conn.findByOpcode(protocol.PROTO_NC_QUEST_START_ACK)
	if len(startAcks) == 0 {
		t.Fatalf("expected PROTO_NC_QUEST_START_ACK")
	}
	if len(startAcks[0].Payload) < 2 || binary.LittleEndian.Uint16(startAcks[0].Payload[0:2]) != 0 {
		t.Fatalf("expected START_ACK errCode=0 for Quest %d, got payload=%x", chosenQuest.QuestID, startAcks[0].Payload)
	}
	if len(conn.findByOpcode(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD after start")
	}

	// 2. PROTO_NC_QUEST_READ_REQ
	conn.clear()
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_READ_REQ, startBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_QUEST_READ_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_QUEST_READ_ACK")
	}

	// 3. OnQuestMobKilled (con Status=1 Doing) -> actualiza contador de caza y emite NOTIFY_MOB_KILL_CMD
	for _, mr := range chosenQuest.MobKills {
		if mr.Enabled && mr.MobOrSpeciesID > 0 {
			mobID := mr.MobOrSpeciesID
			if mr.IsSpecies {
				if len(ds.QuestCatalog.SpeciesMobs[mr.MobOrSpeciesID]) == 0 {
					continue
				}
				mobID = ds.QuestCatalog.SpeciesMobs[mr.MobOrSpeciesID][0]
			}
			conn.clear()
			h.OnQuestMobKilled(p, mobID)
			if len(conn.findByOpcode(protocol.PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD)) == 0 {
				t.Fatalf("expected PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD when killing target mob %d", mobID)
			}
			break
		}
	}

	// 4. PROTO_NC_QUEST_SCRIPT_CMD_ACK, REWARD_SELECT_ITEM_INDEX_CMD & CLIENT_SCENARIO_DONE_REQ (entrega y GameLog Type 9)
	conn.clear()
	scriptAckBuf := make([]byte, 7)
	binary.LittleEndian.PutUint16(scriptAckBuf[0:2], chosenQuest.QuestID)
	scriptAckBuf[2] = 0
	binary.LittleEndian.PutUint32(scriptAckBuf[3:7], 1)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_SCRIPT_CMD_ACK, scriptAckBuf))

	selRewardBuf := make([]byte, 6)
	binary.LittleEndian.PutUint16(selRewardBuf[0:2], chosenQuest.QuestID)
	binary.LittleEndian.PutUint32(selRewardBuf[2:6], 0)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD, selRewardBuf))

	scenDoneBuf := make([]byte, 3)
	binary.LittleEndian.PutUint16(scenDoneBuf[0:2], chosenQuest.QuestID)
	scenDoneBuf[2] = 1
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ, scenDoneBuf))

	if len(conn.findByOpcode(protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK")
	}

	doneQuests, _ := repo.GetCharacterDoneQuests(p.Char.ID)
	foundDone := false
	for _, dq := range doneQuests {
		if uint16(dq) == chosenQuest.QuestID {
			foundDone = true
			break
		}
	}
	if !foundDone {
		t.Fatalf("expected Quest %d persisted in CharacterDoneQuests", chosenQuest.QuestID)
	}

	if gl := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == p.Char.ID && e.Int1 == int(chosenQuest.QuestID)
	}); gl == nil {
		t.Fatalf("expected GameLog Type 9 recorded for Quest %d", chosenQuest.QuestID)
	}

	// 5. SELECT_START_REQ + GIVE_UP_REQ
	conn.clear()
	selStart := make([]byte, 4)
	binary.LittleEndian.PutUint16(selStart[0:2], chosenQuest.StartNPCID)
	binary.LittleEndian.PutUint16(selStart[2:4], chosenQuest.QuestID)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_SELECT_START_REQ, selStart))
	if len(conn.findByOpcode(protocol.PROTO_NC_QUEST_SELECT_START_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_QUEST_SELECT_START_ACK")
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_QUEST_GIVE_UP_REQ, startBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_QUEST_GIVE_UP_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_QUEST_GIVE_UP_ACK")
	}

	// 6. NPC Click, Menu Ack, ServerMenu Ack, Scenarios y Tutorial Step
	conn.clear()
	npc := &ZoneNPC{MobID: 1, Name: "RouHelper", X: 1000, Y: 1000, Menu: 1}
	zMap.AddNPC(npc)
	clickBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(clickBuf, npc.Handle)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_NPCCLICK_CMD, clickBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ)) == 0 {
		t.Fatalf("expected PROTO_NC_ACT_NPCMENUOPEN_REQ on NPC click")
	}

	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK, []byte{1}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_MENU_SERVERMENU_ACK, []byte{0}))

	areaBuf := make([]byte, 32)
	copy(areaBuf, "Zone_Tutorial_01")
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_SCENARIO_AREAENTRY_ACK, areaBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_SCENARIO_AREALEAVE_ACK, areaBuf))

	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_TUTORIAL_STEP_ACK, []byte{3}))
	if p.Char.TutorialStep != 3 {
		t.Fatalf("expected TutorialStep=3 persisted, got %d", p.Char.TutorialStep)
	}
	if len(conn.findByOpcode(protocol.PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK")
	}
}

func TestFase15_LuaInstancesIDKQGTWeddingAndHolyPromise(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	groom, connGroom := newTestClientAndPlayer(t, repo, ds, zMap, h, 9511, "GroomKnight", 50, 1000, 1000)
	bride, connBride := newTestClientAndPlayer(t, repo, ds, zMap, h, 9512, "BrideMage", 50, 1020, 1000)
	groom.Char.TCharacterState.Money = 500000

	// 1. Dept 41: Instance Dungeon (Level Select Menu Ack & Rank List Req)
	idBuf := make([]byte, 6)
	binary.LittleEndian.PutUint16(idBuf[0:2], 1)
	binary.LittleEndian.PutUint32(idBuf[2:6], 1)
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK, idBuf))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD")
	}
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_REQ, nil))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK")
	}

	// 2. Dept 22: Kingdom Quest Execution in Map + KQReward + GameLog Type H
	connGroom.clear()
	kqBuf := make([]byte, 5)
	kqBuf[0] = 1
	var firstKQID uint32 = 1
	for k := range ds.KQCatalog.ByID {
		firstKQID = uint32(k)
		break
	}
	binary.LittleEndian.PutUint32(kqBuf[1:5], firstKQID)
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, kqBuf))
	if groom.CurrentMap == nil || groom.CurrentMap.InstanceKind != "KQ" {
		t.Fatalf("expected Groom warped into KQ instance map, got %+v", groom.CurrentMap)
	}
	if len(connGroom.findByOpcode(protocol.PROTO_NC_SCENARIO_TIMER_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_SCENARIO_TIMER_CMD from KQ timer")
	}
	groom.CurrentMap.bridgeKQRewardIndex(uint32(groom.MapObjectID), "1")
	groom.CurrentMap.bridgeRewardAll("")
	groom.CurrentMap.bridgeEndOfKingdomQuest(true)
	if len(connGroom.findByOpcode(protocol.PROTO_NC_KQ_COMPLETE_CMD)) == 0 || len(connGroom.findByOpcode(protocol.PROTO_NC_KQ_REWARD_REQ)) == 0 {
		t.Fatalf("expected KQ_COMPLETE_CMD and KQ_REWARD_REQ from KQ completion")
	}
	if glH := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == groom.Char.ID && e.TargetCharNo == int(firstKQID)
	}); glH == nil {
		t.Fatalf("expected GameLog Type H recorded for KQ")
	}

	// 3. Dept 29: Guild Tournament 6 C2S en mapa real con Lua script
	connGroom.clear()
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ, []byte{1}))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK")
	}
	if groom.CurrentMap != nil {
		groom.CurrentMap.AddPlayer(groom)
	}
	flagBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(flagBuf, 10)
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ, flagBuf))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK)) == 0 || len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_RESULT_CMD)) == 0 {
		t.Fatalf("expected GT FlagCapture ACK and RESULT_CMD")
	}
	skBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(skBuf, 1)
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ, skBuf))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK")
	}
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ, nil))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK")
	}
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ, []byte{1}))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK")
	}
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ, []byte{1}))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_ACK")
	}

	// 4. Dept 8 Wedding
	h.WarpPlayer(groom, zMap.MapID, 1000, 1000, 0)
	h.WarpPlayer(bride, zMap.MapID, 1020, 1000, 0)
	zMap.AddPlayer(groom)
	zMap.AddPlayer(bride)
	connGroom.clear()
	connBride.clear()

	propBuf := make([]byte, 39)
	binary.LittleEndian.PutUint16(propBuf[0:2], bride.MapObjectID)
	copy(propBuf[2:39], "Will you marry me?")
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ, propBuf))
	if len(connBride.findByOpcode(protocol.PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ)) == 0 {
		t.Fatalf("expected bride to receive PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ")
	}

	accBuf := make([]byte, 40)
	binary.LittleEndian.PutUint16(accBuf[0:2], groom.MapObjectID)
	accBuf[2] = 1
	copy(accBuf[3:40], "Yes!")
	h.HandlePacket(bride.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ, accBuf))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)) == 0 {
		t.Fatalf("expected groom to receive PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK")
	}

	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_HALL_RESERV_REQ, nil))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK")
	}
	h.HandlePacket(bride.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_REQ, nil))
	h.HandlePacket(bride.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_REQ, nil))
	if len(connBride.findByOpcode(protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK")
	}

	_ = repo.CompleteWedding(groom.Char.ID, bride.Char.ID)
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ, nil))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK")
	}

	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_AGREEMENT_DIVORCE_REQ, nil))
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_COMPULSORY_DIVORCE_REQ, nil))
	h.HandlePacket(bride.Client, buildTestPacket(protocol.PROTO_NC_ACT_WEDDING_DIVORCE_REFUSE_CMD, nil))
	if len(connGroom.findByOpcode(protocol.PROTO_NC_ACT_WEDDING_DIVORCE_REQ_ACK)) < 2 {
		t.Fatalf("expected PROTO_NC_ACT_WEDDING_DIVORCE_REQ_ACK for agreement & compulsory divorce")
	}
	if glJ := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == groom.Char.ID && e.Int1 == 2
	}); glJ == nil {
		t.Fatalf("expected GameLog Type J recorded for wedding/divorce")
	}

	// 5. Dept 37: HolyPromise Cen Rewards + GameLog Type K
	_ = repo.SetHolyPromiseMoney(groom.Char.ID, 12345, 5)
	connGroom.clear()
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_REQ, nil))
	getCenAcks := connGroom.findByOpcode(protocol.PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK)
	if len(getCenAcks) == 0 || len(getCenAcks[0].Payload) < 10 || binary.LittleEndian.Uint64(getCenAcks[0].Payload[2:10]) != 12345 {
		t.Fatalf("expected HOLY_PROMISE_GET_CEN_REWARD_ACK with 12345 cen")
	}
	moneyBefore := groom.Char.TCharacterState.Money
	h.HandlePacket(groom.Client, buildTestPacket(protocol.PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_REQ, nil))
	if groom.Char.TCharacterState.Money != moneyBefore+12345 {
		t.Fatalf("expected Groom money +12345 after HolyPromise withdraw, got %d", groom.Char.TCharacterState.Money)
	}
	if glK := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == groom.Char.ID && e.Int1 == 12345
	}); glK == nil {
		t.Fatalf("expected GameLog Type K recorded for HolyPromise withdraw")
	}
}

func TestFase15_MiniHouseFullLifecycleFurnitureAndGates(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	owner, connOwner := newTestClientAndPlayer(t, repo, ds, zMap, h, 9521, "HouseOwner", 40, 1200, 1200)
	guest, connGuest := newTestClientAndPlayer(t, repo, ds, zMap, h, 9522, "HouseGuest", 40, 1210, 1200)

	var mhItemID, furItemID uint16
	for k := range ds.MiniHouseCatalog.ByItemID {
		mhItemID = k
		break
	}
	for k := range ds.MiniHouseCatalog.Furniture {
		furItemID = k
		break
	}
	if mhItemID == 0 {
		mhItemID = 100
	}
	if furItemID == 0 {
		furItemID = 101
	}
	h.addItemToInventory(owner, mhItemID, 1, nil)
	h.addItemToInventory(owner, furItemID, 1, nil)
	invItems, _ := repo.GetItemsByStorage(owner.Char.ID, 9)
	if len(invItems) < 2 {
		t.Fatalf("expected owner to have MiniHouse & furniture items in Inven")
	}
	mhSlot := byte(invItems[0].Storage)
	furItemKey := invItems[1].ItemKey
	furHandle := uint16(furItemKey & 0xFFFF)
	if furHandle == 0 {
		furHandle = 77
	}

	// 1. Activate + Building
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_ACTIV_REQ, []byte{mhSlot}))
	if len(connOwner.findByOpcode(protocol.PROTO_NC_MINIHOUSE_ACTIV_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_ACTIV_ACK")
	}

	buildBuf := make([]byte, 32+101)
	buildBuf[0] = 8
	copy(buildBuf[10:31], "XelionHouse")
	copy(buildBuf[32:], "Welcome")
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_BUILDING_REQ, buildBuf))
	if len(connOwner.findByOpcode(protocol.PROTO_NC_MINIHOUSE_BUILDING_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_BUILDING_ACK")
	}

	// 2. Enter + ArrangeMode + CreateFurniture (Box 12) + Rearrange + Delete + Complete
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_LOGINCOMPLETE_REQ, nil))
	if len(connOwner.findByOpcode(protocol.PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK")
	}
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ, []byte{1}))
	if len(connOwner.findByOpcode(protocol.PROTO_NC_MINIHOUSE_ARRANGEMODE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_ARRANGEMODE_ACK")
	}

	createFur := make([]byte, 1+18)
	createFur[0] = 1
	binary.LittleEndian.PutUint16(createFur[1:3], uint16(invItems[1].ItemKey&0xFFFF))
	binary.LittleEndian.PutUint32(createFur[3:7], 500)
	binary.LittleEndian.PutUint32(createFur[7:11], 600)
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ, createFur))
	if len(connOwner.findByOpcode(protocol.PROTO_NC_MINIHOUSE_CREATE_FURNITURE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_CREATE_FURNITURE_ACK")
	}
	box12Items, _ := repo.GetItemsByStorage(owner.Char.ID, 12)
	if len(box12Items) == 0 {
		t.Fatalf("expected placed furniture persisted in StorageType 12")
	}

	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ, nil))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_REARRANGE_REQ, createFur))

	// 3. Modificar Pass, Title, MaxEnter, Notify, OwnerBlog
	passBuf := make([]byte, 16)
	copy(passBuf, "sec123")
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ, passBuf))
	titleBuf := make([]byte, 21)
	copy(titleBuf, "CozyDen")
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ, titleBuf))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ, []byte{12}))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ, []byte("Stay awhile")))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_OWNERBLOG_REQ, []byte("My Fiesta Diary")))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_OWNERBLOG_REQ, nil))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ, []byte{0}))

	// 4. Guest VisitReady + Visit + Emotion + Portal Gates (tMinihouseGate)
	visitHndBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(visitHndBuf, owner.MapObjectID)
	h.HandlePacket(guest.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_VISITREADY_REQ, visitHndBuf))
	if len(connGuest.findByOpcode(protocol.PROTO_NC_MINIHOUSE_VISITREADY_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_VISITREADY_ACK")
	}

	visitReqBuf := make([]byte, 2+16)
	binary.LittleEndian.PutUint16(visitReqBuf[0:2], owner.MapObjectID)
	copy(visitReqBuf[2:18], "sec123")
	h.HandlePacket(guest.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_VISIT_REQ, visitReqBuf))
	if len(connGuest.findByOpcode(protocol.PROTO_NC_MINIHOUSE_VISIT_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_VISIT_ACK")
	}

	emBuf := make([]byte, 5)
	binary.LittleEndian.PutUint16(emBuf[0:2], furHandle)
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ, emBuf))
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ, nil))

	portalOpenBuf := make([]byte, 12)
	binary.LittleEndian.PutUint16(portalOpenBuf[0:2], 9001)
	portalOpenBuf[2] = 1
	copy(portalOpenBuf[3:], "Rou")
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ, portalOpenBuf))
	gates, _ := repo.GetMiniHouseGates(owner.Char.ID)
	if len(gates) == 0 {
		t.Fatalf("expected tMinihouseGate persisted after PORTAL_OPEN_REQ")
	}
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ, nil))

	h.HandlePacket(guest.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_EXIT_REQ, nil))
	if len(connGuest.findByOpcode(protocol.PROTO_NC_MINIHOUSE_EXIT_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MINIHOUSE_EXIT_ACK")
	}

	delBuf := []byte{1, byte(furItemKey & 0xFF), byte((furItemKey >> 8) & 0xFF)}
	h.HandlePacket(owner.Client, buildTestPacket(protocol.PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ, delBuf))

	expectedHouseID := int(ds.MiniHouseCatalog.ByItemID[mhItemID].Handle)
	if glF := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == owner.Char.ID && e.TargetID == expectedHouseID && e.Type == 2
	}); glF == nil {
		t.Fatalf("expected GameLog Type F recorded for MiniHouse lifecycle")
	}
}

func TestFase15_PetLifecycleLuaGTIAndWeaponTitle(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9531, "PetTamer", 50, 1000, 1000)

	var petTypeID uint32 = 1
	for _, act := range ds.PetCatalog.ActionsByInx {
		if act != nil && act.ID > 0 {
			petTypeID = act.ID
			break
		}
	}

	// 1. Invocar Mascota (ASK_NEW_NAME_REQ con firstSummon=1 -> Crea en tPet/tPetNew + GameLog G)
	askBuf := make([]byte, 7)
	binary.LittleEndian.PutUint16(askBuf[0:2], 501)
	binary.LittleEndian.PutUint32(askBuf[2:6], petTypeID)
	askBuf[6] = 1
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_PET_ASK_NEW_NAME_REQ, askBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_PET_ASK_NEW_NAME_ACK")
	}
	if p.ActivePet == nil || p.ActivePet.Handle != 501 {
		t.Fatalf("expected ActivePet summoned with handle 501")
	}

	// 2. Cambiar Nombre (SET_NAME_REQ + CMD) + Cancelar (SET_NAME_CANCEL_REQ)
	setBuf := make([]byte, 18)
	binary.LittleEndian.PutUint16(setBuf[0:2], 501)
	copy(setBuf[2:18], "DracoMini")
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_PET_SET_NAME_REQ, setBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_PET_SET_NAME_ACK)) == 0 || len(conn.findByOpcode(protocol.PROTO_NC_PET_SET_NAME_CMD)) == 0 {
		t.Fatalf("expected PET_SET_NAME_ACK and CMD")
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_PET_SET_NAME_CANCEL_REQ, setBuf[:2]))
	if len(conn.findByOpcode(protocol.PROTO_NC_PET_SET_NAME_CANCEL_ACK)) == 0 {
		t.Fatalf("expected PET_SET_NAME_CANCEL_ACK")
	}

	// 3. Alimentación, Tick periódico y los 6 puentes Lua Pet_*
	zMap.bridgePetChangeMind(501, 40)
	if zMap.bridgePetGetMind(501) != 40 {
		t.Fatalf("expected Pet Mind=40, got %d", zMap.bridgePetGetMind(501))
	}
	if !h.FeedActivePet(p, 25, 10) || p.ActivePet.Mind != 65 || p.ActivePet.Stress != 10 {
		t.Fatalf("unexpected FeedActivePet state: %+v", p.ActivePet)
	}
	if zMap.bridgePetGetStress(501) != 10 {
		t.Fatalf("expected Pet Stress=10")
	}
	if !zMap.bridgePetSaveTendency(501) {
		t.Fatalf("expected bridgePetSaveTendency=true")
	}
	for act := range ds.PetCatalog.ActionsByInx {
		if inx, _, _, _, ok := zMap.bridgePetGetActionRecord(501, act); !ok || inx == "" {
			t.Fatalf("expected valid action record for %q", act)
		}
		break
	}
	h.TickActivePet(p)
	if p.ActivePet.Mind != 64 {
		t.Fatalf("expected TickActivePet to decrement Mind by 1 (to 64), got %d", p.ActivePet.Mind)
	}
	zMap.bridgePetUnsummon(501)
	if p.ActivePet != nil {
		t.Fatalf("expected ActivePet nil after bridgePetUnsummon")
	}
	if glG := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == p.Char.ID && e.TargetID == int(petTypeID)
	}); glG == nil {
		t.Fatalf("expected GameLog Type G recorded for Pet")
	}

	// 4. WeaponTitle (Dept 23): WT_LICENSE_REQ + OnWeaponTitleMobKilled
	conn.clear()
	var wtMobID uint16 = 100
	if len(ds.WeaponTitleCatalog.Rows) > 0 && ds.WeaponTitleCatalog.Rows[0].MobID > 0 {
		wtMobID = ds.WeaponTitleCatalog.Rows[0].MobID
	}
	var wepID, licID uint16 = 1001, 1002
	for id := uint16(1); id < 60000; id++ {
		if info, ok := ds.ItemStore.GetItemInfo(id); ok && info != nil {
			if info.Class == 5 && wepID == 1001 {
				wepID = info.ID
			} else if info.Class != 5 && licID == 1002 {
				licID = info.ID
			}
			if wepID != 1001 && licID != 1002 {
				break
			}
		}
	}
	wepItem := &database.Item{Owner: p.Char.ID, StorageType: 0, Storage: 10, ItemID: wepID, Flags: 1}
	licItem := &database.Item{Owner: p.Char.ID, StorageType: 9, Storage: 5, ItemID: licID, Flags: 1}
	_ = repo.CreateItem(wepItem)
	_ = repo.CreateItem(licItem)

	wtReq := make([]byte, 7)
	binary.LittleEndian.PutUint16(wtReq[0:2], (9<<8)|5)
	binary.LittleEndian.PutUint16(wtReq[2:4], (0<<8)|10)
	binary.LittleEndian.PutUint16(wtReq[4:6], wtMobID)
	wtReq[6] = 1
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_WT_LICENSE_REQ, wtReq))
	if len(conn.findByOpcode(protocol.PROTO_NC_WT_LICENSE_ACK)) == 0 || len(conn.findByOpcode(protocol.PROTO_NC_WT_TITLE_CMD)) == 0 {
		t.Fatalf("expected WT_LICENSE_ACK and WT_TITLE_CMD")
	}

	h.OnWeaponTitleMobKilled(p, wtMobID)
	if len(conn.findByOpcode(protocol.PROTO_NC_WT_MOBINC_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_WT_MOBINC_CMD after killing licensed mob %d", wtMobID)
	}
}

func TestFase15_CasinoItemDiceTaisaiAndGambleSlotMachine(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9541, "HighRoller", 60, 1000, 1000)
	p.Char.TCharacterState.Money = 1000000

	// 1. Dept 12: PROTO_NC_ITEM_DICE_GAME_START_REQ
	diceBuf := []byte{11, 0, 1}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_DICE_GAME_START_REQ, diceBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_DICE_GAME_START_ACK)) == 0 || len(conn.findByOpcode(protocol.PROTO_NC_ITEM_DICE_GAME_RESULT_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_ITEM_DICE_GAME_START_ACK and RESULT_CMD")
	}

	// 2. Dept 47: Casino Coins Buy / VIP Card UI / Sell
	conn.clear()
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_ACK)) == 0 {
		t.Fatalf("expected GAMBLE_COIN_VIPCARD_UI_OPEN_ACK")
	}
	buyBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buyBuf, 500)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_COIN_BUY_REQ, buyBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_COIN_BUY_ACK)) == 0 || p.CasinoCoins != 500 {
		t.Fatalf("expected 500 CasinoCoins after COIN_BUY_REQ, got %d", p.CasinoCoins)
	}

	// 3. Dept 43: DICE_TAISAI
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_RANK_REQ, nil))
	modeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(modeBuf, 1)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_REQ, modeBuf))

	betBuf := make([]byte, 19)
	binary.LittleEndian.PutUint64(betBuf[0:8], 10)
	binary.LittleEndian.PutUint16(betBuf[8:10], 1)
	betBuf[10] = 0
	binary.LittleEndian.PutUint64(betBuf[11:19], 20)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_BETTING_REQ, betBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ, []byte{0, 0}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ, []byte{0, 0}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ, betBuf[:10]))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_DICE_TAISAI_GAME_LEAVE_REQ, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD)) == 0 {
		t.Fatalf("expected Dept 43 DICE_TAISAI join & roll result packets")
	}

	// 4. Dept 47: GAMBLE_DICE_TAISAI & GAMBLE_SLOTMACHINE
	conn.clear()
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ, []byte{1, 0}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ, betBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ, betBuf[:10]))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_REQ, nil))

	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ, nil))
	slotStart := make([]byte, 4)
	binary.LittleEndian.PutUint32(slotStart, 10)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ, slotStart))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_REQ, nil))

	sellBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sellBuf, 50)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_COIN_SELL_REQ, sellBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ, nil))

	if len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_ACK)) == 0 {
		t.Fatalf("expected SlotMachine start/stop/wheel/exit ACKs")
	}

	if glI := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool {
		return e.CharNo == p.Char.ID && e.TargetCharNo == 500
	}); glI == nil {
		t.Fatalf("expected GameLog Type I recorded for Casino Coins & Slots")
	}
}
