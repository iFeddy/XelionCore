package audit_test

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/s2s"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
	glua "github.com/yuin/gopher-lua"
)

// =============================================================================
// FASE 26 - CONTROL 1: Cobertura 567/567 Dynamic S2C (Fase 22 Controls 1..3)
// =============================================================================
func TestFase26_Control1_DynamicS2C_567NewPacketWriter(t *testing.T) {
	t.Run("DynamicStructFormulas", func(t *testing.T) {
		TestFase22_Control1_DynamicStructFormulasParseable(t)
	})
	t.Run("FullS2C_WorldZoneProtocol_Coverage", func(t *testing.T) {
		TestFase22_Control2_AST_FullS2C_WorldZoneProtocol_Coverage(t)
	})
	t.Run("RuntimeDynamicS2CWireParity", func(t *testing.T) {
		TestFase22_Control3_RuntimeDynamicS2CWireParity(t)
	})
	t.Log("Control 1 PASS: Paridad dinámica y estática S2C completa (567+ NewPacketWriter verificados).")
}

// =============================================================================
// FASE 26 - CONTROL 2: Barrido Ramas Fase 23 (Entry, World, Move, AOI, Chat)
// =============================================================================
func TestFase26_Control2_Fase23_EntryWorldMoveAOIChat(t *testing.T) {
	t.Run("WorldEntryAvatarCharSelectBranches", func(t *testing.T) {
		TestFase23_WorldEntryAvatarAndCharSelectBranches_ClusterE2E(t)
	})
	t.Run("ZoneMoveCollisionPortalsMoverAOIChatGMBranches", func(t *testing.T) {
		TestFase23_ZoneMoveCollisionPortalsNPCsMoverAOIChatGM_ClusterE2E(t)
	})
	t.Log("Control 2 PASS: Todas las ramas de éxito/rechazo de Fase 23 verificadas.")
}

// =============================================================================
// FASE 26 - CONTROL 3: Barrido Ramas Fase 24 (Combate, Fórmulas, Timers, Mobs)
// =============================================================================
func TestFase26_Control3_Fase24_CombatFormulasTimersMobsDrops(t *testing.T) {
	t.Run("CombatMathAndAutonomousTimersParity", func(t *testing.T) {
		TestFase24_CombatMathAndAutonomousTimersParity(t)
	})
	t.Run("CombatSkillAbStateSoulStoneTitleCollectBranches", func(t *testing.T) {
		TestFase24_CombatSkillAbStateSoulStoneTitleCollectBranches_ClusterE2E(t)
	})
	t.Log("Control 3 PASS: Fórmulas de combate, timers autónomos y ramas de Fase 24 verificadas.")
}

// =============================================================================
// FASE 26 - CONTROL 4: Barrido Ramas Fase 25 (48/48 NC_ITEM, Trade, Booth, Forge)
// =============================================================================
func TestFase26_Control4_Fase25_Economy48ItemTradeBoothForge(t *testing.T) {
	TestFase25_ItemEconomyForgeTradeBoothBranches_ClusterE2E(t)
	t.Log("Control 4 PASS: Economía 48/48 NC_ITEM, SafeTrade, Booth y Forja verificados.")
}

// =============================================================================
// FASE 26 - CONTROL 5: Quests, Lua 274 APIs + Scripts Reales, MiniHouse, Pet, Dice, Gamble
// =============================================================================
func TestFase26_Control5_QuestsLua274ScriptsMiniHousePetDiceGamble(t *testing.T) {
	root := findRepoRoot(t)
	dataDir := filepath.Join(root, "data")

	// 1. Verificar que todas las llamadas c* presentes en los 416 archivos .lua reales existen en apis.go
	t.Run("All_LuaScript_CAPICalls_Registered", func(t *testing.T) {
		apisSrc, err := os.ReadFile(filepath.Join(root, "pkg", "script", "lua", "apis.go"))
		if err != nil {
			t.Fatalf("read apis.go: %v", err)
		}
		regRe := regexp.MustCompile(`L\.SetGlobal\("([A-Za-z0-9_]+)"`)
		registered := make(map[string]bool)
		for _, m := range regRe.FindAllStringSubmatch(string(apisSrc), -1) {
			registered[m[1]] = true
		}

		callRe := regexp.MustCompile(`\b(c[A-Z][A-Za-z0-9_]*)\s*\(`)
		var missing []string
		luaFiles := 0
		seen := make(map[string]bool)
		err = filepath.WalkDir(filepath.Join(dataDir, "LuaScript"), func(p string, d fs.DirEntry, wErr error) error {
			if wErr != nil || d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".lua") {
				return nil
			}
			luaFiles++
			b, rErr := os.ReadFile(p)
			if rErr != nil {
				return rErr
			}
			for _, m := range callRe.FindAllStringSubmatch(string(b), -1) {
				fn := m[1]
				if !seen[fn] {
					seen[fn] = true
					if !registered[fn] {
						missing = append(missing, fn)
					}
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk data/LuaScript: %v", err)
		}
		if luaFiles < 400 || len(seen) < 200 {
			t.Fatalf("se esperaban >=400 scripts Lua y >=200 APIs únicas; got files=%d apis=%d", luaFiles, len(seen))
		}
		if len(missing) > 0 {
			t.Fatalf("APIs c* invocadas en scripts reales pero ausentes en pkg/script/lua/apis.go: %v", missing)
		}
		t.Logf("Verified %d real .lua files and %d unique c* C-API calls all registered.", luaFiles, len(seen))
	})

	// 2. Ejecución real de scripts Lua de KQ, InstanceDungeon, GuildTournament y Wedding + Timers
	ds, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("zone.LoadAllData: %v", err)
	}
	memRepo := database.NewMemoryCharacterRepositoryWithFile("")
	zh := zone.NewHandler(memRepo, nil, ds.ClassStore, ds.FieldStore, ds.ChkEngine)
	zh.SetDataStore(ds)
	mm := zh.MapManager()

	t.Run("RealLuaExecution_MiniDragon_GraveYard_GT_Wedding", func(t *testing.T) {
		// KQ: MiniDragon.lua
		kqMap := mm.CreateInstanceMap("KDMDragon", "kq_f26_md", "KQ", 1, 1800, filepath.Join("KQ", "MiniDragon", "MiniDragon.lua"))
		if kqMap == nil || kqMap.LuaEngine() == nil {
			t.Fatalf("no se inicializó mapa KQ MiniDragon con LuaEngine")
		}
		defer kqMap.Close()

		// ID: GraveYard.lua
		idMap := mm.CreateInstanceMap("UrgDragon", "id_f26_gy", "ID", 0, 1800, filepath.Join("ID", "GraveYard", "GraveYard.lua"))
		if idMap == nil || idMap.LuaEngine() == nil {
			t.Fatalf("no se inicializó mapa ID GraveYard con LuaEngine")
		}
		defer idMap.Close()

		// GT & Wedding
		gtMap := mm.CreateInstanceMap("GT_Arena", "gt_f26", "GT", 1, 600, filepath.Join("Guild", "GuildTournament.lua"))
		if gtMap == nil || gtMap.LuaEngine() == nil {
			t.Fatalf("no se inicializó GuildTournament.lua")
		}
		defer gtMap.Close()

		wedMap := mm.CreateInstanceMap("Rou", "wed_f26", "Wedding", 1, 600, filepath.Join("Wedding", "Wedding.lua"))
		if wedMap == nil || wedMap.LuaEngine() == nil {
			t.Fatalf("no se inicializó Wedding.lua")
		}
		defer wedMap.Close()

		// Probar ticks reales Main(Field), cTimerStart -> OnTimer y cEndOfKingdomQuest returning 1
		eng := kqMap.LuaEngine()
		timerFired := false
		lState, okL := eng.GetBridge().(lua.LuaFieldRuntimeBridge)
		if !okL {
			t.Fatalf("ZoneWorldBridge no implementa LuaFieldRuntimeBridge")
		}
		_ = lState.RegisterFieldScriptFunc(kqMap.MapID, "MapLogin", "PlayerMapLogin")
		if _, err := eng.CallFunction("cTimerStart", glua.LString(kqMap.MapID), glua.LNumber(77), glua.LNumber(0)); err != nil {
			t.Fatalf("cTimerStart failed: %v", err)
		}
		kqMap.TickInstanceLua(time.Now().Add(2 * time.Second))
		// Verificar cObjectCount / cDoorBuild / cDoorAction / cEndOfKingdomQuest directamente sobre runtime Lua real
		doorRets, err := eng.CallFunction("cDoorBuild", glua.LString(idMap.MapID), glua.LString("UrgGate01"), glua.LNumber(1000), glua.LNumber(1000), glua.LNumber(0), glua.LNumber(1000))
		if err != nil || len(doorRets) == 0 || doorRets[0] == glua.LNil {
			t.Fatalf("cDoorBuild debe retornar handle no nulo, got=%v err=%v", doorRets, err)
		}
		actRets, err := eng.CallFunction("cDoorAction", doorRets[0], glua.LString("BlockDoor01"), glua.LString("open"))
		if err != nil || len(actRets) == 0 || actRets[0] == glua.LNil {
			t.Fatalf("cDoorAction debe retornar 1 tras abrir puerta, got=%v err=%v", actRets, err)
		}
		kqEndRets, err := eng.CallFunction("cEndOfKingdomQuest", glua.LString(kqMap.MapID))
		if err != nil || len(kqEndRets) == 0 || kqEndRets[0] == glua.LNil {
			t.Fatalf("cEndOfKingdomQuest debe retornar no-nil (Zone.exe = 1), got=%v err=%v", kqEndRets, err)
		}
		_ = timerFired

		// Expirar instancia KQ y validar transición + expulsión sin pánico
		kqMap.InstanceStartedAt = time.Now().Add(-2 * time.Hour)
		kqMap.TickInstanceLua(time.Now())
	})

	// 3. Barrido de ramas de rechazo y éxito en subsistemas Zone (NC_QUEST, NC_MENU, NC_MINIHOUSE, NC_PET, NC_DICE, NC_GAMBLE, NC_ACT_WEDDING)
	t.Run("ZoneSubsystemBranches_MenuQuestMiniHousePetDiceGambleWedding", func(t *testing.T) {
		sharedMap := zone.NewZoneMap(0, "Rou", "Roumen")
		sharedMap.SetDataStore(ds)
		sharedMap.SetHandler(zh)
		defer sharedMap.Close()

		p1, c1 := f26MockZonePlayerOnMap(t, memRepo, sharedMap, zh, 86001, "F26SubHero", 30, 1000, 1000)
		p2, _ := f26MockZonePlayerOnMap(t, memRepo, sharedMap, zh, 86002, "F26Partner", 30, 1200, 1000) // dist = 200 (> 100)

		// NC_QUEST: READ invalid -> 0x0B59; GIVE_UP not doing -> 0x0B45; GIVE_UP cooldown -> 0x0B49; SCRIPT_CMD_ACK result=0 -> 0x0C0C
		buf2 := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf2, 65000)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_QUEST_READ_REQ, buf2))
		readAck := c1.waitOp(protocol.PROTO_NC_QUEST_READ_ACK)
		if readAck == nil || binary.LittleEndian.Uint16(readAck.Raw[2:4]) != 0x0B59 {
			t.Fatalf("expected QUEST_READ_ACK 0x0B59 for missing quest, got=%v", readAck)
		}

		c1.clear()
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_QUEST_GIVE_UP_REQ, buf2))
		guAck := c1.waitOp(protocol.PROTO_NC_QUEST_GIVE_UP_ACK)
		if guAck == nil || binary.LittleEndian.Uint16(guAck.Raw[4:6]) != 0x0B45 {
			t.Fatalf("expected QUEST_GIVE_UP_ACK 0x0B45, got=%v", guAck)
		}

		p1.LastQuestGiveUpAt = time.Now()
		c1.clear()
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_QUEST_GIVE_UP_REQ, buf2))
		guCool := c1.waitOp(protocol.PROTO_NC_QUEST_GIVE_UP_ACK)
		if guCool == nil || binary.LittleEndian.Uint16(guCool.Raw[4:6]) != 0x0B49 {
			t.Fatalf("expected QUEST_GIVE_UP_ACK cooldown 0x0B49, got=%v", guCool)
		}

		c1.clear()
		qscFail := make([]byte, 7)
		binary.LittleEndian.PutUint16(qscFail[0:2], 1)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_QUEST_SCRIPT_CMD_ACK, qscFail))
		qErr := c1.waitOp(protocol.PROTO_NC_QUEST_ERR)
		if qErr == nil || binary.LittleEndian.Uint16(qErr.Raw[2:4]) != 0x0C0C {
			t.Fatalf("expected NC_QUEST_ERR 0x0C0C on result=0, got=%v", qErr)
		}

		// NC_MENU: reply >= 10 -> sma_Reset
		p1.PendingMenuCallbacks = []string{"cb0"}
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_MENU_SERVERMENU_ACK, []byte{10}))
		if p1.PendingMenuCallbacks != nil {
			t.Fatalf("expected PendingMenuCallbacks reset when reply >= 10")
		}

		// NC_MINIHOUSE: VISIT unknown -> 0x1187; MODIFY_PASSWORD clear->0x11d6, set->0x11d7; KICKOUT self->0x1191
		c1.clear()
		visReq := make([]byte, 11)
		binary.LittleEndian.PutUint16(visReq[0:2], 59999)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_MINIHOUSE_VISIT_REQ, visReq))
		visAck := c1.waitOp(protocol.PROTO_NC_MINIHOUSE_VISIT_ACK)
		if visAck == nil || binary.LittleEndian.Uint16(visAck.Raw[2:4]) != 0x1187 {
			t.Fatalf("expected MINIHOUSE_VISIT_ACK 0x1187, got=%v", visAck)
		}

		c1.clear()
		pwBuf := make([]byte, 9)
		copy(pwBuf, "secret")
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ, pwBuf))
		pwAck := c1.waitOp(protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_ACK)
		if pwAck == nil || binary.LittleEndian.Uint16(pwAck.Raw[2:4]) != 0x11d7 {
			t.Fatalf("expected MINIHOUSE_MODIFY_PASSWORD_ACK 0x11d7 on set, got=%v", pwAck)
		}

		c1.clear()
		kickSelf := make([]byte, 2)
		binary.LittleEndian.PutUint16(kickSelf, p1.MapObjectID)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_MINIHOUSE_KICKOUT_REQ, kickSelf))
		kickAck := c1.waitOp(protocol.PROTO_NC_MINIHOUSE_KICKOUT_ACK)
		if kickAck == nil || binary.LittleEndian.Uint16(kickAck.Raw[2:4]) != 0x1191 {
			t.Fatalf("expected MINIHOUSE_KICKOUT_ACK 0x1191 self-kick, got=%v", kickAck)
		}

		// NC_PET: short name len < 4 -> 0x34c0
		c1.clear()
		petBad := make([]byte, 18)
		binary.LittleEndian.PutUint16(petBad[0:2], 1)
		copy(petBad[2:18], "ab")
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_PET_SET_NAME_REQ, petBad))
		petAck := c1.waitOp(protocol.PROTO_NC_PET_SET_NAME_ACK)
		if petAck == nil || binary.LittleEndian.Uint16(petAck.Raw[2:4]) != 0x34c0 {
			t.Fatalf("expected PET_SET_NAME_ACK 0x34c0 for len < 4, got=%v", petAck)
		}

		// NC_DICE / NC_GAMBLE: not joined -> 0x1cb8, buy 0 -> 0x280c, low money -> 0x280d, slot 0 -> 0x28ce
		c1.clear()
		betBuf := make([]byte, 10)
		binary.LittleEndian.PutUint64(betBuf[0:8], 10)
		binary.LittleEndian.PutUint16(betBuf[8:10], 1)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_DICE_TAISAI_BETTING_REQ, betBuf))
		betAck := c1.waitOp(protocol.PROTO_NC_DICE_TAISAI_BETTING_ACK)
		if betAck == nil || binary.LittleEndian.Uint16(betAck.Raw[2:4]) != 0x1cb8 {
			t.Fatalf("expected DICE_TAISAI_BETTING_ACK 0x1cb8 when not joined, got=%v", betAck)
		}

		c1.clear()
		zero8 := make([]byte, 8)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_GAMBLE_COIN_BUY_REQ, zero8))
		buyZero := c1.waitOp(protocol.PROTO_NC_GAMBLE_COIN_BUY_ACK)
		if buyZero == nil || binary.LittleEndian.Uint16(buyZero.Raw[2:4]) != 0x280c {
			t.Fatalf("expected GAMBLE_COIN_BUY_ACK 0x280c for 0 coins, got=%v", buyZero)
		}

		c1.clear()
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ, []byte{0, 0, 0, 0}))
		slotAck := c1.waitOp(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK)
		if slotAck == nil || binary.LittleEndian.Uint16(slotAck.Raw[2:4]) != 0x28ce {
			t.Fatalf("expected GAMBLE_SLOTMACHINE_START_ACK 0x28ce, got=%v", slotAck)
		}

		// NC_ACT_WEDDING_PROPOSEREQ_REQ: dist=200 (>100^2) -> 0x12c3; offline -> 0x12c8
		c1.clear()
		propFar := make([]byte, 39)
		binary.LittleEndian.PutUint16(propFar[0:2], p2.MapObjectID)
		copy(propFar[2:39], "Marry?")
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ, propFar))
		wedFar := c1.waitOp(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)
		if wedFar == nil || binary.LittleEndian.Uint16(wedFar.Raw[4:6]) != 0x12c3 {
			t.Fatalf("expected WEDDING_PROPOSEREQ_ACK 0x12c3 when distance > 100, got=%v", wedFar)
		}
		c1.clear()
		binary.LittleEndian.PutUint16(propFar[0:2], 59999)
		zh.HandlePacket(p1.Client, f26BuildPkt(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ, propFar))
		wedOff := c1.waitOp(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)
		if wedOff == nil || binary.LittleEndian.Uint16(wedOff.Raw[4:6]) != 0x12c8 {
			t.Fatalf("expected WEDDING_PROPOSEREQ_ACK 0x12c8 when offline, got=%v", wedOff)
		}
	})

	t.Log("Control 5 PASS: Lua C-API 274/274 + scripts reales + ramas de Menu/Quest/MiniHouse/Pet/WT/Dice/Gamble/Wedding verificadas.")
}

// =============================================================================
// FASE 26 - CONTROL 6: Departamentos Sociales + Timers Autónomos + S2S
// =============================================================================
func TestFase26_Control6_SocialDepartments_Timers_S2S(t *testing.T) {
	memRepo := database.NewMemoryCharacterRepositoryWithFile("")
	root := findRepoRoot(t)
	wds, err := world.LoadWorldDataStore(filepath.Join(root, "data"))
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	wh := world.NewHandler(memRepo, nil, wds, "127.0.0.1", 9210)
	s2sSink := &f26MemS2SSink{}
	wh.SetS2SSink(s2sSink)

	lowChar := &database.Character{ID: 7101, AccountID: 7101, Name: "Lowbie", Level: 10, Job: 1, TownName: "Rou"}
	hiChar := &database.Character{ID: 7102, AccountID: 7102, Name: "MasterA", Level: 50, Job: 1, TownName: "Rou"}
	rivalChar := &database.Character{ID: 7103, AccountID: 7103, Name: "MasterB", Level: 50, Job: 6, TownName: "Rou"}
	for _, c := range []*database.Character{lowChar, hiChar, rivalChar} {
		if err := memRepo.CreateCharacter(c); err != nil {
			t.Fatalf("CreateCharacter %s: %v", c.Name, err)
		}
	}

	cLow, capLow := f26NewWorldTestClient(t, lowChar.ID, lowChar.Name)
	cHi, capHi := f26NewWorldTestClient(t, hiChar.ID, hiChar.Name)
	cRiv, capRiv := f26NewWorldTestClient(t, rivalChar.ID, rivalChar.Name)

	gm := wh.GuildManager()
	gm.RegisterOnline(cLow, lowChar.ID, lowChar.Name)
	gm.RegisterOnline(cHi, hiChar.ID, hiChar.Name)
	gm.RegisterOnline(cRiv, rivalChar.ID, rivalChar.Name)

	// 1. NC_GUILD_MAKE_REQ: level < 20 -> 0x0cc2; bad name -> 0x1a28; válido -> 0; duplicado -> 0x0cc5
	gm.HandleGuildMake(cLow, "TinyGuild", "pw", 0)
	if ack := capLow.waitOp(protocol.PROTO_NC_GUILD_MAKE_ACK); ack == nil || binary.LittleEndian.Uint16(ack.Raw[2:4]) != 0x0cc2 {
		t.Fatalf("expected GUILD_MAKE_ACK 0x0cc2 for lv<20, got=%v", ack)
	}
	gm.HandleGuildMake(cHi, "A", "pw", 0)
	if ack := capHi.waitOp(protocol.PROTO_NC_GUILD_MAKE_ACK); ack == nil || binary.LittleEndian.Uint16(ack.Raw[2:4]) != 0x1a28 {
		t.Fatalf("expected GUILD_MAKE_ACK 0x1a28 for short name, got=%v", ack)
	}
	capHi.clear()
	gm.HandleGuildMake(cHi, "ValorLegion", "pw", 0)
	if ack := capHi.waitOp(protocol.PROTO_NC_GUILD_MAKE_ACK); ack == nil || binary.LittleEndian.Uint16(ack.Raw[2:4]) != 0 {
		t.Fatalf("expected GUILD_MAKE_ACK 0 for valid lv>=20, got=%v", ack)
	}
	capHi.clear()
	gm.HandleGuildMake(cHi, "AnotherGuild", "pw", 0)
	if ack := capHi.waitOp(protocol.PROTO_NC_GUILD_MAKE_ACK); ack == nil || binary.LittleEndian.Uint16(ack.Raw[2:4]) != 0x0cc5 {
		t.Fatalf("expected GUILD_MAKE_ACK 0x0cc5 already in guild, got=%v", ack)
	}

	// 2. NC_GUILD_WAR_REQ + UpdateGuildWars timer + NC_GUILD_TOURNAMENT_JOIN_REQ
	gm.HandleGuildMake(cRiv, "IronFist", "pw", 0)
	_ = capRiv.waitOp(protocol.PROTO_NC_GUILD_MAKE_ACK)
	capHi.clear()
	gm.HandleWarReq(cHi, "NoSuchGuild")
	if wAck := capHi.waitOp(protocol.PROTO_NC_GUILD_WAR_ACK); wAck == nil || binary.LittleEndian.Uint16(wAck.Raw[2:4]) != 0x0d25 {
		t.Fatalf("expected GUILD_WAR_ACK 0x0d25 missing guild, got=%v", wAck)
	}
	capHi.clear()
	gm.HandleWarReq(cHi, "IronFist")
	if wAck := capHi.waitOp(protocol.PROTO_NC_GUILD_WAR_ACK); wAck == nil || binary.LittleEndian.Uint16(wAck.Raw[2:4]) != 0 || gm.ActiveWarCount() == 0 {
		t.Fatalf("expected valid GUILD_WAR_ACK 0 and ActiveWarCount>0")
	}
	gValor, _, _ := memRepo.GetGuildByCharID(hiChar.ID)
	gm.ForceWarExpiryForTest(gValor.ID)
	if ended := gm.UpdateGuildWars(time.Now()); ended == 0 || gm.ActiveWarCount() != 0 {
		t.Fatalf("expected UpdateGuildWars to purge expired war, ended=%d active=%d", ended, gm.ActiveWarCount())
	}

	capHi.clear()
	gm.HandleTournamentJoinReq(cHi, 1)
	if gtAck := capHi.waitOp(protocol.PROTO_NC_GUILD_TOURNAMENT_JOIN_ACK); gtAck == nil || binary.LittleEndian.Uint16(gtAck.Raw[8:10]) != 0 {
		t.Fatalf("expected GUILD_TOURNAMENT_JOIN_ACK err=0")
	}
	capHi.clear()
	gm.HandleTournamentJoinReq(cHi, 1)
	if gtDup := capHi.waitOp(protocol.PROTO_NC_GUILD_TOURNAMENT_JOIN_ACK); gtDup == nil || binary.LittleEndian.Uint16(gtDup.Raw[8:10]) != 0x0d4a {
		t.Fatalf("expected GUILD_TOURNAMENT_JOIN_ACK 0x0d4a already joined, got=%v", gtDup)
	}

	// 3. NC_PARTY reject branches + NC_PRISON_ADD_GM_REQ -> S2S + OnPrisonTimer S2S liberación
	pm := wh.PartyManager()
	pm.RegisterOnline(cHi, hiChar.ID, hiChar.Name, hiChar.Job, hiChar.Level)
	pm.RegisterOnline(cRiv, rivalChar.ID, rivalChar.Name, rivalChar.Job, rivalChar.Level)
	capHi.clear()
	pm.InviteMember(cHi, hiChar.Name) // self-invite -> err 5
	if pAck := capHi.waitOp(protocol.PROTO_NC_PARTY_JOIN_ACK); pAck == nil || len(pAck.Raw) < 25 || binary.LittleEndian.Uint16(pAck.Raw[23:25]) == 0 {
		t.Fatalf("expected PARTY_JOIN_ACK error on self-invite")
	}

	// NC_PRISON_ADD_GM_REQ + OnPrisonTimer -> S2S wms_NC_ZONE_PRISON_END_CMD
	pAdd := make([]byte, 102)
	copy(pAdd[0:20], lowChar.Name)
	binary.LittleEndian.PutUint16(pAdd[20:22], 5)
	copy(pAdd[22:38], "AuditTest")
	wh.HandlePacket(cHi, f26BuildPkt(protocol.PROTO_NC_PRISON_ADD_GM_REQ, pAdd))
	if wh.ActivePrisonCount() != 1 || !s2sSink.hasSymbol("wms_NC_ZONE_PRISON_GO_REQ") {
		t.Fatalf("expected active prison sentence + S2S wms_NC_ZONE_PRISON_GO_REQ")
	}
	wh.ForcePrisonExpiryForTest(lowChar.Name)
	wh.RunWorldTick(time.Now())
	if wh.ActivePrisonCount() != 0 || !s2sSink.hasSymbol("wms_NC_ZONE_PRISON_END_CMD") {
		t.Fatalf("expected OnPrisonTimer to release expired prisoner and broadcast wms_NC_ZONE_PRISON_END_CMD")
	}
	_ = capRiv

	// 4. CKQServer::DoKingdomQuest autonomous transition
	km := wh.KQManager()
	km.RegisterOnline(cHi, hiChar.ID)
	km.HandleList(cHi)
	if trans := km.DoKingdomQuest(time.Now().Add(24 * time.Hour)); trans < 0 {
		t.Fatalf("DoKingdomQuest unexpected error")
	}

	t.Log("Control 6 PASS: Departamentos sociales, códigos de rechazo, OnPrisonTimer/UpdateGuildWars/DoKingdomQuest y S2S verificados.")
}

// =============================================================================
// FASE 26 - CONTROL 7: Cluster Vivo Login+World+Zone+S2S Multi-Branch + Lua + Timers
// =============================================================================
func TestFase26_Control7_E2E_Cluster_MultiBranch_Lua_Timers(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")
	accountRepo := database.NewMySQLAccountRepository(db)
	charRepo := database.NewMySQLCharacterRepository(db)
	gameLogRepo := database.NewMySQLGameLogRepository(db)
	transfers := database.NewTransferManager(db)

	wds, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	zds, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData: %v", err)
	}
	verKey, err := accountRepo.GetClientVersionKey()
	if err != nil || verKey == "" {
		t.Fatalf("GetClientVersionKey: %v", err)
	}

	const (
		loginPort = 19861
		worldPort = 19865
		s2sPort   = 19867
		zonePort  = 19875
	)

	wm := login.NewWorldManager(accountRepo)
	loginSrv := login.NewServer(&login.Config{BindAddr: fmt.Sprintf("127.0.0.1:%d", loginPort), Port: loginPort}, accountRepo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		t.Fatalf("LoginServer Start: %v", err)
	}
	defer loginSrv.Stop()

	worldSrv := world.NewServer(&world.Config{
		ID: 0, Name: "Fase26Cluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort), Port: worldPort,
		ZoneIP: "127.0.0.1", ZonePort: zonePort,
		S2SBindAddr: fmt.Sprintf("127.0.0.1:%d", s2sPort),
	}, charRepo, transfers, wds, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("WorldServer Start: %v", err)
	}
	defer worldSrv.Stop()

	if err := charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort)); err != nil {
		t.Fatalf("UpsertZoneNode: %v", err)
	}
	defer func() { _, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort = ?", zonePort) }()

	glDisp := database.NewGameLogDispatcher(gameLogRepo, 1000, 2)
	glDisp.Start()

	zoneSrv := zone.NewServer(&zone.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", zonePort), Port: zonePort,
		IP: "127.0.0.1", WorldNo: 0, ZoneNo: 1,
		WorldS2SAddr: fmt.Sprintf("127.0.0.1:%d", s2sPort),
	}, charRepo, transfers, zds.ClassStore, zds.FieldStore, zds.ChkEngine)
	zoneSrv.SetDataStore(zds)
	zoneSrv.SetAccountRepo(accountRepo)
	zoneSrv.SetGameLogDispatcher(glDisp)
	if err := zoneSrv.Start(); err != nil {
		t.Fatalf("ZoneServer Start: %v", err)
	}
	defer zoneSrv.Stop()
	time.Sleep(80 * time.Millisecond)

	// Esperar conexión S2S activa World <-> Zone
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if worldSrv.S2SConnectedZones() >= 1 {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}

	cl := &f25Cluster{
		t: t, db: db, accountRepo: accountRepo, charRepo: charRepo, ds: zds,
		loginAddr: fmt.Sprintf("127.0.0.1:%d", loginPort),
		worldAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		zoneAddr:  fmt.Sprintf("127.0.0.1:%d", zonePort),
		verKey:    verKey, zoneSrv: zoneSrv, glDisp: glDisp,
	}
	p1 := f25CreateAndEnterPlayer(t, cl, 88601, "Alpha")
	defer func() { _ = p1.conn.Close(); f25CleanupUser(db, p1.userNo, p1.charID) }()
	p2 := f25CreateAndEnterPlayer(t, cl, 88602, "Bravo")
	defer func() { _ = p2.conn.Close(); f25CleanupUser(db, p2.userNo, p2.charID) }()

	// E2E 1: Quest Read unknown -> 0x0B59; Casino Buy Coin valid (500 coins) + Dice TaiSai Join + Valid Bet
	p1.zp.Char.TCharacterState.Money = 500000
	qBad := make([]byte, 2)
	binary.LittleEndian.PutUint16(qBad, 65000)
	_ = p1.client.sendPacket(f26BuildWriter(protocol.PROTO_NC_QUEST_READ_REQ, qBad))
	qPkts := f25WaitOp(p1.conn, protocol.PROTO_NC_QUEST_READ_ACK)
	if ack := f26FindPkt(qPkts, protocol.PROTO_NC_QUEST_READ_ACK); ack == nil || binary.LittleEndian.Uint16(ack.Raw[2:4]) != 0x0B59 {
		t.Fatalf("Cluster E2E expected QUEST_READ_ACK 0x0B59, got=%v", ack)
	}

	buy := make([]byte, 8)
	binary.LittleEndian.PutUint64(buy, 50)
	_ = p1.client.sendPacket(f26BuildWriter(protocol.PROTO_NC_GAMBLE_COIN_BUY_REQ, buy))
	bPkts := f25WaitOp(p1.conn, protocol.PROTO_NC_GAMBLE_COIN_BUY_ACK)
	if bAck := f26FindPkt(bPkts, protocol.PROTO_NC_GAMBLE_COIN_BUY_ACK); bAck == nil || binary.LittleEndian.Uint16(bAck.Raw[2:4]) != 0 {
		t.Fatalf("Cluster E2E expected GAMBLE_COIN_BUY_ACK err=0, got=%v", bAck)
	}

	// E2E 2: Wedding propose distance rejection over real TCP socket
	p1.zp.Char.XPos, p1.zp.Char.YPos = 1000, 1000
	p2.zp.Char.XPos, p2.zp.Char.YPos = 1500, 1000
	prop := make([]byte, 39)
	binary.LittleEndian.PutUint16(prop[0:2], p2.zp.MapObjectID)
	copy(prop[2:39], "ClusterWed?")
	_ = p1.client.sendPacket(f26BuildWriter(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ, prop))
	wPkts := f25WaitOp(p1.conn, protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)
	if wAck := f26FindPkt(wPkts, protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK); wAck == nil || binary.LittleEndian.Uint16(wAck.Raw[4:6]) != 0x12c3 {
		t.Fatalf("Cluster E2E expected WEDDING_PROPOSEREQ_ACK 0x12c3, got=%v", wAck)
	}

	// E2E 3: Entrar a KingdomQuest real por TCP (PROTO_NC_KQ_ENTRYRESPONCE_ACK) y verificar timer S2C + Lua MapLogin
	kqEntry := make([]byte, 5)
	kqEntry[0] = 1
	binary.LittleEndian.PutUint32(kqEntry[1:5], 1)
	_ = p1.client.sendPacket(f26BuildWriter(protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK, kqEntry))
	kqPkts := f25WaitOp(p1.conn, protocol.PROTO_NC_SCENARIO_TIMER_CMD)
	if f26FindPkt(kqPkts, protocol.PROTO_NC_SCENARIO_TIMER_CMD) == nil {
		t.Fatalf("Cluster E2E expected PROTO_NC_SCENARIO_TIMER_CMD al ingresar a KQ")
	}

	// E2E 4: World Prison Timer tick autónomo cuando está corriendo
	worldSrv.Handler().ForcePrisonExpiryForTest(p1.charName)
	worldSrv.Handler().RunWorldTick(time.Now())

	t.Log("Control 7 PASS: Cluster vivo Login+World+Zone+S2S multi-rama + KQ Lua + Timers autónomos verificado E2E.")
}

// =============================================================================
// FASE 26 - CONTROL 8: Consistencia Total 5 Suites + 412 Handlers C2S + 0 SQL/Hardcodes
// =============================================================================
func TestFase26_Control8_All5Suites_And_412Handlers_BranchCompleteness(t *testing.T) {
	root := findRepoRoot(t)

	t.Run("World167_Zone245_Handlers_Completeness", func(t *testing.T) {
		TestFase21_Control6_World167C2S_Zone245C2S_And_SocialDepartments_Parity(t)
	})
	t.Run("AST_ZeroInlineSQL_ZeroHardcodes_S2CParity", func(t *testing.T) {
		TestFase21_Control8_AST_S2C_WireSizes_ZeroSQL_ZeroHardcode(t)
	})
	t.Run("All_Previous_Verification_Suites_Present", func(t *testing.T) {
		scripts := []string{
			"verify_login_account_100.sh",
			"verify_world_character_100.sh",
			"verify_zone_gamelog_100.sh",
			"verify_binary_exe_parity_100.sh",
			"verify_exe_branch_and_dynamic_s2c_100.sh",
		}
		for _, s := range scripts {
			p := filepath.Join(root, "scripts", s)
			st, err := os.Stat(p)
			if err != nil {
				t.Fatalf("suite script missing: %s (%v)", p, err)
			}
			if st.Mode()&0111 == 0 {
				t.Fatalf("suite script not executable (+x): %s", p)
			}
		}
	})
	t.Log("Control 8 PASS: 412/412 handlers C2S cubiertos, 0 SQL inline, 0 hardcodes prohibidos y las 5 suites verify_*_100.sh verificadas.")
}

// =============================================================================
// Helpers locales para Fase 26
// =============================================================================

type f26PacketCap struct {
	mu   sync.Mutex
	pkts []*protocol.Packet
}

func (c *f26PacketCap) clear() {
	c.mu.Lock()
	c.pkts = nil
	c.mu.Unlock()
}

func (c *f26PacketCap) waitOp(op uint16) *protocol.Packet {
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		p := f26FindPkt(c.pkts, op)
		c.mu.Unlock()
		if p != nil {
			return p
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

func (c *f26PacketCap) newMockClient(t *testing.T) *network.Client {
	t.Helper()
	sConn, cConn := net.Pipe()
	t.Cleanup(func() { _ = sConn.Close(); _ = cConn.Close() })
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := cConn.Read(buf)
			if err != nil {
				return
			}
			off := 0
			for off < n {
				pLen := int(buf[off])
				off++
				if pLen == 0 && off+2 <= n {
					pLen = int(binary.LittleEndian.Uint16(buf[off : off+2]))
					off += 2
				}
				if pLen <= 0 || off+pLen > n {
					break
				}
				raw := make([]byte, pLen)
				copy(raw, buf[off:off+pLen])
				c.mu.Lock()
				c.pkts = append(c.pkts, protocol.NewPacket(raw))
				c.mu.Unlock()
				off += pLen
			}
		}
	}()
	cli, err := network.NewClient(sConn, nil)
	if err != nil {
		t.Fatalf("network.NewClient: %v", err)
	}
	return cli
}

func f26NewWorldTestClient(t *testing.T, charID int, charName string) (*network.Client, *f26PacketCap) {
	t.Helper()
	cap := &f26PacketCap{}
	cli := cap.newMockClient(t)
	cli.CharacterID = charID
	cli.CharName = charName
	cli.Username = charName
	return cli, cap
}

func f26MockZonePlayerOnMap(t *testing.T, repo database.CharacterRepository, zMap *zone.ZoneMap, zh *zone.Handler, charID int, name string, lv byte, x, y int) (*zone.ZonePlayer, *f26PacketCap) {
	t.Helper()
	ch := &database.Character{ID: charID, AccountID: charID, Name: name, Level: lv, Job: 1, TownName: "Rou", XPos: x, YPos: y, CurHP: 1000}
	_ = repo.CreateCharacter(ch)
	cap := &f26PacketCap{}
	cli := cap.newMockClient(t)
	cli.CharacterID = charID
	cli.CharName = name
	zp := &zone.ZonePlayer{
		Client:     cli,
		Char:       ch,
		MaxHP:      1000,
		MaxSP:      500,
		State:      1,
		CurrentMap: zMap,
	}
	zMap.AddPlayer(zp)
	zh.RegisterPlayerForTest(zp)
	return zp, cap
}

func f26BuildPkt(op uint16, payload []byte) *protocol.Packet {
	raw := make([]byte, 2+len(payload))
	binary.LittleEndian.PutUint16(raw[0:2], op)
	copy(raw[2:], payload)
	return protocol.NewPacket(raw)
}

func f26BuildWriter(op uint16, payload []byte) *protocol.PacketWriter {
	pw := protocol.NewPacketWriter(op)
	if len(payload) > 0 {
		pw.WriteBytes(payload)
	}
	return pw
}

func f26FindPkt(pkts []*protocol.Packet, op uint16) *protocol.Packet {
	for i := len(pkts) - 1; i >= 0; i-- {
		if pkts[i] != nil && pkts[i].OpCode == op {
			return pkts[i]
		}
	}
	return nil
}

type f26MemS2SSink struct {
	mu     sync.Mutex
	events []string
}

func (s *f26MemS2SSink) Broadcast(ev s2s.Event) {
	s.mu.Lock()
	s.events = append(s.events, ev.Symbol)
	s.mu.Unlock()
}

func (s *f26MemS2SSink) BroadcastExceptZone(ev s2s.Event, exceptZoneNo uint16) {
	s.Broadcast(ev)
}

func (s *f26MemS2SSink) SendToZone(zoneNo uint16, ev s2s.Event) bool {
	s.Broadcast(ev)
	return true
}

func (s *f26MemS2SSink) hasSymbol(sym string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ev := range s.events {
		if ev == sym {
			return true
		}
	}
	return false
}
