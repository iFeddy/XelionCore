package lua

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	glua "github.com/yuin/gopher-lua"
)

type recordedKQBridge struct {
	DefaultBridge

	limitMap      map[string]int
	teamByHandle  map[uint32]int
	joinerHandles map[uint32]bool

	endCalls        []string
	rewardCalls     [][2]string
	rewardIdxCalls  []struct {
		handle uint32
		idx    string
	}
	timerLifeMap    []struct {
		mapName string
		sec     int
	}
	timerLifeObj    []struct {
		handle uint32
		sec    int
	}
	vanishTimerMaps []string
	teamBattleCalls []struct {
		mapName string
		enabled bool
	}
	questResCalls   [][2]string
	questResInd     []struct {
		handle uint32
		result string
	}
	scoreBoardCalls []struct {
		mapName   string
		show      bool
		round     int
		redWin    int
		redScore  int
		blueWin   int
		blueScore int
	}
	scoreInfoMap    []struct {
		mapName string
		scores  []uint32
	}
	scoreInfoObj    []struct {
		handle uint32
		scores []uint32
	}
	timerMapCalls   []struct {
		mapName string
		sec     int
	}
	timerObjCalls   []struct {
		handle uint32
		sec    int
	}
}

var _ LuaKQRuntimeBridge = (*recordedKQBridge)(nil)

func (b *recordedKQBridge) EndOfKingdomQuestMap(mapName string) bool {
	b.endCalls = append(b.endCalls, mapName)
	return mapName != ""
}

func (b *recordedKQBridge) GetKQLimitSecondExact(mapName string) (int, bool) {
	if sec, ok := b.limitMap[mapName]; ok {
		return sec, true
	}
	return 0, false
}

func (b *recordedKQBridge) ShowKQTimerWithLifeMap(mapName string, seconds int) bool {
	b.timerLifeMap = append(b.timerLifeMap, struct {
		mapName string
		sec     int
	}{mapName, seconds})
	return mapName != ""
}

func (b *recordedKQBridge) ShowKQTimerWithLifeObjExact(playerHandle uint32, seconds int) bool {
	b.timerLifeObj = append(b.timerLifeObj, struct {
		handle uint32
		sec    int
	}{playerHandle, seconds})
	return playerHandle != 0
}

func (b *recordedKQBridge) VanishTimerMap(mapName string) bool {
	b.vanishTimerMaps = append(b.vanishTimerMaps, mapName)
	return mapName != ""
}

func (b *recordedKQBridge) GetKQTeamTypeExact(handle uint32) int {
	if t, ok := b.teamByHandle[handle]; ok {
		return t
	}
	return 2 // KQTT_MAX (exe: 0x005E823C)
}

func (b *recordedKQBridge) IsKQJoinerExact(handle uint32) bool {
	return b.joinerHandles[handle]
}

func (b *recordedKQBridge) KQRewardIndexExact(playerHandle uint32, rewardIndex string) {
	b.rewardIdxCalls = append(b.rewardIdxCalls, struct {
		handle uint32
		idx    string
	}{playerHandle, rewardIndex})
}

func (b *recordedKQBridge) RewardKQMap(mapName, rewardType string) bool {
	b.rewardCalls = append(b.rewardCalls, [2]string{mapName, rewardType})
	return mapName != ""
}

func (b *recordedKQBridge) SetTeamBattleMap(mapName string, enabled bool) {
	b.teamBattleCalls = append(b.teamBattleCalls, struct {
		mapName string
		enabled bool
	}{mapName, enabled})
}

func (b *recordedKQBridge) QuestResultMap(mapName, result string) bool {
	b.questResCalls = append(b.questResCalls, [2]string{mapName, result})
	return mapName != "" && (result == "Success" || result == "Fail")
}

func (b *recordedKQBridge) QuestResultIndividualExact(playerHandle uint32, result string) bool {
	b.questResInd = append(b.questResInd, struct {
		handle uint32
		result string
	}{playerHandle, result})
	return playerHandle != 0 && (result == "Success" || result == "Fail")
}

func (b *recordedKQBridge) ScoreBoardAllInMap(mapName string, show bool, round, redWin, redScore, blueWin, blueScore int) {
	b.scoreBoardCalls = append(b.scoreBoardCalls, struct {
		mapName   string
		show      bool
		round     int
		redWin    int
		redScore  int
		blueWin   int
		blueScore int
	}{mapName, show, round, redWin, redScore, blueWin, blueScore})
}

func (b *recordedKQBridge) ScoreInfoAllInMap(mapName string, scores []uint32) {
	cp := append([]uint32(nil), scores...)
	b.scoreInfoMap = append(b.scoreInfoMap, struct {
		mapName string
		scores  []uint32
	}{mapName, cp})
}

func (b *recordedKQBridge) ScoreInfoObj(playerHandle uint32, scores []uint32) {
	cp := append([]uint32(nil), scores...)
	b.scoreInfoObj = append(b.scoreInfoObj, struct {
		handle uint32
		scores []uint32
	}{playerHandle, cp})
}

func (b *recordedKQBridge) TimerMap(mapName string, seconds int) bool {
	b.timerMapCalls = append(b.timerMapCalls, struct {
		mapName string
		sec     int
	}{mapName, seconds})
	return mapName != ""
}

func (b *recordedKQBridge) TimerObj(playerHandle uint32, seconds int) bool {
	b.timerObjCalls = append(b.timerObjCalls, struct {
		handle uint32
		sec    int
	}{playerHandle, seconds})
	return playerHandle != 0
}

// TestKQBindings_DisassembledParity verifies every disassembled Zone.exe c*KQ* binding 1:1.
func TestKQBindings_DisassembledParity(t *testing.T) {
	bridge := &recordedKQBridge{
		DefaultBridge: DefaultBridge{MapID: "KDEddyHill00"},
		limitMap:      map[string]int{"KDEddyHill00": 1800},
		teamByHandle:  map[uint32]int{101: 0, 102: 1},
		joinerHandles: map[uint32]bool{101: true},
	}

	eng, err := NewEngine(Config{Bridge: bridge})
	if err != nil {
		t.Fatalf("NewEngine failed: %v", err)
	}
	defer eng.Close()

	script := `
		-- cGetKQLimitSecond (exe: 0x005E5390)
		res_limit_ok = cGetKQLimitSecond("KDEddyHill00")
		res_limit_missing = cGetKQLimitSecond("UnknownMap")
		res_limit_nil = cGetKQLimitSecond(nil)

		-- cGetKQTeamType (exe: 0x005E81C0): 0=RED, 1=BLUE, 2=KQTT_MAX on nil/missing
		res_team_red = cGetKQTeamType(101)
		res_team_blue = cGetKQTeamType(102)
		res_team_unknown = cGetKQTeamType(999)
		res_team_nil = cGetKQTeamType(nil)

		-- cIsKQJoiner (exe: 0x005E82D0)
		res_joiner_yes = cIsKQJoiner(101)
		res_joiner_no = cIsKQJoiner(102)
		res_joiner_nil = cIsKQJoiner(nil)

		-- cShowKQTimerWithLife / _Obj / cVanishTimer (exe: 0x005E54F0 / 0x005E5690 / 0x005E5840)
		res_tlife_map = cShowKQTimerWithLife("KDEddyHill00", 600)
		res_tlife_obj = cShowKQTimerWithLife_Obj(101, 300)
		res_vanish = cVanishTimer("KDEddyHill00")

		-- cTimer / cTimer_Obj (exe: 0x005DD310 / 0x005DD4B0)
		res_timer_map = cTimer("KDEddyHill00", 120)
		res_timer_obj = cTimer_Obj(101, 45)

		-- cSetTeamBattle (exe: 0x005E83E0)
		cSetTeamBattle("KDEddyHill00", true)
		cSetTeamBattle("KDEddyHill00", false)

		-- cScoreBoard_AllInMap (exe: 0x005EDAD0): clamps win<=5, score<=99
		cScoreBoard_AllInMap("KDEddyHill00", 1, 2, 9, 150, 7, 120)

		-- cScoreInfo_AllInMap / cScoreInfo (exe: 0x005E3970 / 0x005E3790)
		cScoreInfo_AllInMap("KDEddyHill00", 2, 11, 22)
		cScoreInfo(101, 2, 33, 44)

		-- cKQRewardIndex (exe: 0x005E4000) & cReward (exe: 0x005DF060)
		cKQRewardIndex(101, "KQ_SLIME_01")
		res_rew = cReward("KDEddyHill00", "KQ")

		-- cQuestResult / cQuestResult_Individual (exe: 0x005E30D0 / 0x005ECCA0)
		res_q_ok = cQuestResult("KDEddyHill00", "Success")
		res_q_bad = cQuestResult("KDEddyHill00", "Invalid")
		res_q_ind = cQuestResult_Individual(101, "Fail")

		-- cEndOfKingdomQuest (exe: 0x005DF1F0)
		res_end = cEndOfKingdomQuest("KDEddyHill00")
		res_end_nil = cEndOfKingdomQuest(nil)
	`
	if err := eng.state.DoString(script); err != nil {
		t.Fatalf("DoString failed: %v", err)
	}

	assertNum := func(name string, want float64) {
		t.Helper()
		v := eng.GetGlobal(name)
		if v.Type() != glua.LTNumber || float64(v.(glua.LNumber)) != want {
			t.Fatalf("%s = %v (%s), want %v", name, v, v.Type(), want)
		}
	}
	assertNil := func(name string) {
		t.Helper()
		if v := eng.GetGlobal(name); v != glua.LNil {
			t.Fatalf("%s = %v, want nil", name, v)
		}
	}

	assertNum("res_limit_ok", 1800)
	assertNil("res_limit_missing")
	assertNil("res_limit_nil")

	assertNum("res_team_red", 0)
	assertNum("res_team_blue", 1)
	assertNum("res_team_unknown", 2)
	assertNum("res_team_nil", 2) // KQTT_MAX (exe: 0x005E823C)

	if eng.GetGlobal("res_joiner_yes") != glua.LTrue {
		t.Fatalf("res_joiner_yes = %v, want true", eng.GetGlobal("res_joiner_yes"))
	}
	if eng.GetGlobal("res_joiner_no") != glua.LFalse {
		t.Fatalf("res_joiner_no = %v, want false", eng.GetGlobal("res_joiner_no"))
	}
	if eng.GetGlobal("res_joiner_nil") != glua.LFalse {
		t.Fatalf("res_joiner_nil = %v, want false (exe: 0x005E8387)", eng.GetGlobal("res_joiner_nil"))
	}

	assertNum("res_tlife_map", 1)
	assertNum("res_tlife_obj", 1)
	assertNum("res_vanish", 1)
	assertNum("res_timer_map", 1)
	assertNum("res_timer_obj", 1)
	assertNum("res_rew", 1)
	assertNum("res_q_ok", 1)
	assertNil("res_q_bad")
	assertNum("res_q_ind", 1)
	assertNum("res_end", 1)
	assertNil("res_end_nil")

	if len(bridge.teamBattleCalls) != 2 || !bridge.teamBattleCalls[0].enabled || bridge.teamBattleCalls[1].enabled {
		t.Fatalf("unexpected teamBattleCalls: %+v", bridge.teamBattleCalls)
	}
	if len(bridge.scoreBoardCalls) != 1 {
		t.Fatalf("expected 1 scoreBoard call, got %d", len(bridge.scoreBoardCalls))
	}
	sb := bridge.scoreBoardCalls[0]
	if !sb.show || sb.round != 2 || sb.redWin != 5 || sb.redScore != 99 || sb.blueWin != 5 || sb.blueScore != 99 {
		t.Fatalf("scoreBoard clamping mismatch: %+v", sb)
	}
	if len(bridge.scoreInfoMap) != 1 || len(bridge.scoreInfoMap[0].scores) != 2 || bridge.scoreInfoMap[0].scores[0] != 11 || bridge.scoreInfoMap[0].scores[1] != 22 {
		t.Fatalf("scoreInfoMap mismatch: %+v", bridge.scoreInfoMap)
	}
	if len(bridge.scoreInfoObj) != 1 || len(bridge.scoreInfoObj[0].scores) != 2 || bridge.scoreInfoObj[0].scores[0] != 33 || bridge.scoreInfoObj[0].scores[1] != 44 {
		t.Fatalf("scoreInfoObj mismatch: %+v", bridge.scoreInfoObj)
	}
	if len(bridge.rewardIdxCalls) != 1 || bridge.rewardIdxCalls[0].handle != 101 || bridge.rewardIdxCalls[0].idx != "KQ_SLIME_01" {
		t.Fatalf("rewardIdxCalls mismatch: %+v", bridge.rewardIdxCalls)
	}
	if len(bridge.rewardCalls) != 1 || bridge.rewardCalls[0] != [2]string{"KDEddyHill00", "KQ"} {
		t.Fatalf("rewardCalls mismatch: %+v", bridge.rewardCalls)
	}
	if len(bridge.endCalls) != 1 || bridge.endCalls[0] != "KDEddyHill00" {
		t.Fatalf("endCalls mismatch: %+v", bridge.endCalls)
	}
}

func findProjectDataDir(t *testing.T) string {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "..", "data"),
		"data",
	}
	for _, c := range candidates {
		if info, err := os.Stat(filepath.Join(c, "LuaScript", "KQ")); err == nil && info.IsDir() {
			return c
		}
	}
	t.Fatalf("data/LuaScript/KQ not found")
	return ""
}

// TestKQScripts_All23DirectoriesAndAllSHNRowsExecute loads and executes the entry point
// of all 23 directories under data/LuaScript/KQ/ as well as all 57 rows of KingdomQuest.shn.
func TestKQScripts_All23DirectoriesAndAllSHNRowsExecute(t *testing.T) {
	dataDir := findProjectDataDir(t)
	luaRoot := filepath.Join(dataDir, "LuaScript")
	kqRoot := filepath.Join(luaRoot, "KQ")

	entries, err := os.ReadDir(kqRoot)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", kqRoot, err)
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	if len(dirs) != 23 {
		t.Fatalf("expected 23 KQ directories in data/LuaScript/KQ/, got %d: %v", len(dirs), dirs)
	}

	tableEntries := map[string]string{
		"conditionofhero": "KQConditionOfHero",
		"gordonmaster":    "KQGordonMaster",
		"honeying":        "KQHoneying",
		"kqhbat":          "KQHBat",
		"underhall":       "KQUnderHall",
	}

	// 1. Execute every one of the 23 KQ directories directly
	for _, dir := range dirs {
		scriptPath := filepath.Join(kqRoot, dir, dir+".lua")
		// KDArena usa wrappers de entrada KDArena1.lua..KDArena6.lua que hacen require("KQ/KDArena/KDArena")
		if dir == "KDArena" {
			scriptPath = filepath.Join(kqRoot, dir, "KDArena1.lua")
		}
		if _, err := os.Stat(scriptPath); err != nil {
			t.Fatalf("missing main script for KQ dir %s: %s", dir, scriptPath)
		}
		bridge := &recordedKQBridge{
			DefaultBridge: DefaultBridge{MapID: "KQTestMap00"},
			limitMap:      map[string]int{"KQTestMap00": 1800},
		}
		eng, err := NewEngine(Config{BasePath: luaRoot, Bridge: bridge})
		if err != nil {
			t.Fatalf("[%s] NewEngine failed: %v", dir, err)
		}
		if err := eng.DoFile(scriptPath); err != nil {
			eng.Close()
			t.Fatalf("[%s] DoFile(%s) failed: %v", dir, scriptPath, err)
		}
		if tbl, isTable := tableEntries[strings.ToLower(dir)]; isTable {
			if _, err := eng.CallTableFunction(tbl, "Main", glua.LString("KQTestMap00"), glua.LNumber(0)); err != nil {
				eng.Close()
				t.Fatalf("[%s] %s.Main failed: %v", dir, tbl, err)
			}
		} else {
			if eng.GetGlobal("Main").Type() != glua.LTFunction {
				eng.Close()
				t.Fatalf("[%s] expected global Main(Field) function", dir)
			}
			if _, err := eng.CallFunction("Main", glua.LString("KQTestMap00")); err != nil {
				eng.Close()
				t.Fatalf("[%s] Main(Field) failed: %v", dir, err)
			}
		}
		eng.Close()
	}

	// 2. Verify all 57 rows in KingdomQuest.shn resolve their ScriptLanguage and execute with InitValue
	store, err := kq.LoadDir(dataDir)
	if err != nil {
		t.Fatalf("kq.LoadDir(%s): %v", dataDir, err)
	}
	if len(store.Quests) == 0 {
		t.Fatalf("expected non-empty KingdomQuest.shn")
	}

	executedSHN := 0
	for _, q := range store.Quests {
		lang := strings.TrimSpace(strings.ReplaceAll(q.ScriptLanguageString(), "\\", "/"))
		initVal := strings.TrimSpace(q.InitValueString())
		if lang == "" {
			continue
		}
		parts := strings.Split(lang, "/")
		base := parts[len(parts)-1]
		baseNoDigits := strings.TrimRight(base, "0123456789")
		parent := ""
		if len(parts) > 1 {
			parent = strings.Join(parts[:len(parts)-1], "/")
		}
		candidates := []string{
			filepath.Join(luaRoot, filepath.FromSlash(lang)+".lua"),
			filepath.Join(luaRoot, filepath.FromSlash(lang), base+".lua"),
		}
		if baseNoDigits != "" && baseNoDigits != base {
			candidates = append(candidates,
				filepath.Join(luaRoot, filepath.FromSlash(parent), baseNoDigits, baseNoDigits+".lua"),
				filepath.Join(luaRoot, "KQ", baseNoDigits, baseNoDigits+".lua"),
			)
		}
		var resolved string
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && !info.IsDir() {
				resolved = c
				break
			}
		}
		if resolved == "" {
			t.Fatalf("KingdomQuest.shn ID=%d ScriptLanguage=%q could not resolve .lua script", q.ID, lang)
		}

		bridge := &recordedKQBridge{
			DefaultBridge: DefaultBridge{MapID: "KQSHNMap00"},
			limitMap:      map[string]int{"KQSHNMap00": int(q.LimitTime) * 60},
		}
		eng, err := NewEngine(Config{BasePath: luaRoot, Bridge: bridge})
		if err != nil {
			t.Fatalf("SHN ID=%d NewEngine failed: %v", q.ID, err)
		}
		eng.SetGlobal("InitValue", glua.LString(initVal))
		if err := eng.DoFile(resolved); err != nil {
			eng.Close()
			t.Fatalf("SHN ID=%d (%s) DoFile(%s) failed: %v", q.ID, lang, resolved, err)
		}
		dirKey := strings.ToLower(strings.TrimSuffix(filepath.Base(resolved), filepath.Ext(resolved)))
		if tbl, isTable := tableEntries[dirKey]; isTable {
			var arg glua.LValue = glua.LString(initVal)
			if n, convErr := strconv.Atoi(initVal); convErr == nil {
				arg = glua.LNumber(n)
			}
			if _, err := eng.CallTableFunction(tbl, "Main", glua.LString("KQSHNMap00"), arg); err != nil {
				eng.Close()
				t.Fatalf("SHN ID=%d (%s) %s.Main(KQSHNMap00, %q) failed: %v", q.ID, lang, tbl, initVal, err)
			}
		} else if eng.GetGlobal("Main").Type() == glua.LTFunction {
			if _, err := eng.CallFunction("Main", glua.LString("KQSHNMap00")); err != nil {
				eng.Close()
				t.Fatalf("SHN ID=%d (%s) Main(KQSHNMap00) failed: %v", q.ID, lang, err)
			}
		}
		eng.Close()
		executedSHN++
	}
	if executedSHN != len(store.Quests) {
		t.Fatalf("executed %d SHN rows, want all %d", executedSHN, len(store.Quests))
	}
}
