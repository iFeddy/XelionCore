package lua_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/script/lua"
	glua "github.com/yuin/gopher-lua"
)

func TestLuaEngine_InitAndCommon(t *testing.T) {
	luaPath := filepath.Join("..", "..", "..", "data", "LuaScript")

	engine, err := lua.NewEngine(lua.Config{
		BasePath: luaPath,
		Debug:    false,
	})
	if err != nil {
		t.Fatalf("Failed to initialize Lua engine: %v", err)
	}
	defer engine.Close()

	// 1. Verificar variables globales de common.lua
	// ObjectType.Player == 2
	val, err := engine.GetGlobalTableField("ObjectType", "Player")
	if err != nil {
		t.Fatalf("Failed to get ObjectType.Player: %v", err)
	}
	if num, ok := val.(glua.LNumber); !ok || int(num) != 2 {
		t.Errorf("Expected ObjectType.Player to be 2, got %v", val)
	}

	// BasicClass.Joker == 21
	valJoker, err := engine.GetGlobalTableField("BasicClass", "Joker")
	if err != nil {
		t.Fatalf("Failed to get BasicClass.Joker: %v", err)
	}
	if num, ok := valJoker.(glua.LNumber); !ok || int(num) != 21 {
		t.Errorf("Expected BasicClass.Joker to be 21, got %v", valJoker)
	}

	// 2. Test cDistanceSquar(0, 0, 3, 4) -> 3^2 + 4^2 = 25
	rets, err := engine.CallFunction("cDistanceSquar", glua.LNumber(0), glua.LNumber(0), glua.LNumber(3), glua.LNumber(4))
	if err != nil {
		t.Fatalf("Failed to call cDistanceSquar: %v", err)
	}
	if len(rets) != 1 {
		t.Fatalf("Expected 1 return value, got %d", len(rets))
	}
	if num, ok := rets[0].(glua.LNumber); !ok || float64(num) != 25.0 {
		t.Errorf("Expected distance squared 25, got %v", rets[0])
	}

	// 3. Test cRandom(10, 20)
	rets, err = engine.CallFunction("cRandom", glua.LNumber(10), glua.LNumber(20))
	if err != nil {
		t.Fatalf("Failed to call cRandom: %v", err)
	}
	randVal := int(rets[0].(glua.LNumber))
	if randVal < 10 || randVal > 20 {
		t.Errorf("cRandom(10, 20) returned out-of-range value: %d", randVal)
	}
}

func TestMobBehavior_LoadAllLua(t *testing.T) {
	luaPath := filepath.Join("..", "..", "..", "data", "LuaScript")

	engine, err := lua.NewEngine(lua.Config{
		BasePath: luaPath,
		Debug:    false,
	})
	if err != nil {
		t.Fatalf("Failed to initialize Lua engine: %v", err)
	}
	defer engine.Close()

	mobBehaviorDir := filepath.Join("..", "..", "..", "data", "MobBehaviorDescript")

	scripts := []string{
		filepath.Join(mobBehaviorDir, "DefaultBehavior.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "KingSlime.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "MaraPirate.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "GordonMaster.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "BossRobo.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "EndlessMaze.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "TravelerDungeon.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "UnderHall.lua"),
	}

	for _, s := range scripts {
		if err := engine.DoFile(s); err != nil {
			t.Errorf("Failed to compile/execute Lua script %s: %v", s, err)
		}
	}
}
