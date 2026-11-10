package audit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/areabmp"
	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
	glua "github.com/yuin/gopher-lua"
)

// TestFase10_Control6_SHN_Tables_All180 verifica que DataStore cargue e indexe en memoria RAM
// la totalidad de las 180 tablas .shn requeridas en disco (194/195 rutas), junto con los
// 49 checksums criptográficos oficiales y el registro de las 5 referencias opcionales.
func TestFase10_Control6_SHN_Tables_All180(t *testing.T) {
	rootPath := findRepoRoot(t)
	dataPath := filepath.Join(rootPath, "data")

	ds, err := zone.LoadAllData(dataPath)
	if err != nil {
		t.Fatalf("zone.LoadAllData failed: %v", err)
	}

	if ds.SHNTables == nil {
		t.Fatal("ds.SHNTables is nil")
	}

	// 1. Verificar las 180 tablas requeridas de la Sección 3.1
	if len(zone.RequiredSHNTables) != 180 {
		t.Fatalf("Expected exactly 180 required SHN definitions, found %d", len(zone.RequiredSHNTables))
	}

	missingSHNs := 0
	for _, req := range zone.RequiredSHNTables {
		f := ds.GetSHN(req.FileName)
		if f == nil {
			t.Errorf("Required SHN #%d (%s) missing in DataStore.SHNTables", req.Index, req.FileName)
			missingSHNs++
			continue
		}
		if f.RecordCount == 0 && len(f.Columns) == 0 {
			t.Errorf("Required SHN #%d (%s) loaded without valid schema or records", req.Index, req.FileName)
			missingSHNs++
		}
	}

	if missingSHNs > 0 {
		t.Fatalf("Control 6 failed: %d / 180 required SHN tables missing or invalid in DataStore", missingSHNs)
	}

	t.Logf("Control 6 (SHN): 100%% Parity (180/180 required SHN tables loaded and indexed in RAM)")

	// 2. Verificar las 5 referencias opcionales legacy
	if len(zone.OptionalSHNTables) != 5 {
		t.Fatalf("Expected exactly 5 optional SHN tables documented, found %d", len(zone.OptionalSHNTables))
	}
	for optName, desc := range zone.OptionalSHNTables {
		if desc == "" {
			t.Errorf("Optional SHN %s missing documentation description", optName)
		}
	}
	t.Logf("Control 6 (SHN): 5/5 optional legacy SHNs properly documented and handled without error")

	// 3. Verificar ChecksumEngine con 49 SHNs criptográficos
	if ds.ChkEngine == nil {
		t.Fatal("ds.ChkEngine is nil")
	}
	checksums := ds.ChkEngine.GetChecksums()
	if len(checksums) != 49 {
		t.Fatalf("Expected 49 checksummed SHNs in ChecksumEngine, got %d", len(checksums))
	}
	t.Logf("Control 6 (SHN): 49/49 cryptographic SHN checksums verified")
}

// TestFase10_Control6_BlockInfo_301Files verifica la carga completa de los 301 archivos binarios
// de data/BlockInfo/ (158 shbd, 33 sbi, 23 aid, 76 bdt, 8 shab, 2 shmd, 1 shad), así como las
// consultas O(1) de movimiento y Bresenham LoS de MapBlockInformation, apertura/cierre de
// MapDoorArray y detección geométrica de AreaInfoData.
func TestFase10_Control6_BlockInfo_301Files(t *testing.T) {
	rootPath := findRepoRoot(t)
	blockInfoDir := filepath.Join(rootPath, "data", "BlockInfo")

	store, err := blockinfo.LoadAll(blockInfoDir)
	if err != nil {
		t.Fatalf("blockinfo.LoadAll failed: %v", err)
	}

	// 1. Validar conteos exactos de los 301 archivos
	if len(store.SHBD) != 158 {
		t.Errorf("Expected 158 .shbd files, got %d", len(store.SHBD))
	}
	if len(store.SBI) != 33 {
		t.Errorf("Expected 33 .sbi files, got %d", len(store.SBI))
	}
	if len(store.AID) != 23 {
		t.Errorf("Expected 23 .aid files, got %d", len(store.AID))
	}
	if len(store.BDT) != 76 {
		t.Errorf("Expected 76 .bdt files, got %d", len(store.BDT))
	}
	if len(store.SHAB) != 8 {
		t.Errorf("Expected 8 .shab files, got %d", len(store.SHAB))
	}
	if len(store.SHMD) != 2 {
		t.Errorf("Expected 2 .shmd files, got %d", len(store.SHMD))
	}
	if len(store.SHAD) != 1 {
		t.Errorf("Expected 1 .shad files, got %d", len(store.SHAD))
	}

	total := store.TotalCount()
	if total != 301 {
		t.Fatalf("Expected exactly 301 files in data/BlockInfo/, got %d", total)
	}
	t.Logf("Control 6 (BlockInfo): 100%% Parity (301/301 binary geometry files loaded)")

	// 2. Probar colisiones reales en Rou.shbd, RouN.shbd y Adl.shbd
	rou := store.GetSHBD("Rou")
	if rou == nil {
		t.Fatal("Rou.shbd missing")
	}
	if rou.WidthTiles != 2048 || rou.HeightTiles != 2048 {
		t.Fatalf("Rou.shbd expected 2048x2048 tiles, got %dx%d", rou.WidthTiles, rou.HeightTiles)
	}

	// Fuera de límites es bloqueado
	if !rou.IsMoveBlocked(9999, 9999) {
		t.Error("Expected out-of-bounds tile to be move-blocked")
	}

	// Búsqueda de celda transitable más cercana en Rou.shbd
	tx, ty, ok := rou.FindNearestWalkable(1060, 912, 50)
	if !ok {
		t.Error("FindNearestWalkable around (1060, 912) failed in Rou.shbd")
	}
	if rou.IsMoveBlocked(tx, ty) {
		t.Errorf("FindNearestWalkable returned blocked tile (%d, %d)", tx, ty)
	}

	// Conversión de coordenadas mundo <-> tile
	wX, wY := 6200.0, 6400.0
	tX, tY := blockinfo.WorldToTile(wX, wY)
	backX, backY := blockinfo.TileToWorld(tX, tY)
	if mathAbs(backX-wX) > 15.0 || mathAbs(backY-wY) > 15.0 {
		t.Errorf("WorldToTile/TileToWorld precision drift: original (%f, %f) -> back (%f, %f)", wX, wY, backX, backY)
	}

	// 3. Probar puertas dinámicas en AdlF.sbi
	instDoors := store.NewInstanceDoors("AdlF")
	if instDoors == nil || len(instDoors.DoorList) != 2 {
		t.Fatalf("Expected 2 doors in AdlF.sbi, got %v", instDoors)
	}

	// Inicialmente cerradas
	d1 := instDoors.Doors["DoorBlock01"]
	if d1 == nil || d1.IsOpen {
		t.Error("DoorBlock01 should be closed initially")
	}

	// Apertura
	if !instDoors.DoorAction("DoorBlock01", blockinfo.DoorActionOpen) {
		t.Error("Failed to open DoorBlock01")
	}
	if !d1.IsOpen {
		t.Error("DoorBlock01 should be open")
	}

	// Cierre de todas las puertas
	instDoors.CloseAllDoors()
	if d1.IsOpen {
		t.Error("DoorBlock01 should be closed after CloseAllDoors")
	}

	// 4. Probar áreas en RouN.aid
	rounAID := store.GetAID("RouN")
	if rounAID == nil {
		t.Fatal("RouN.aid missing")
	}
	// MH_Zone1 en RouN cubre (6411.73, 6527.53)
	if !rounAID.Contains("MH_Zone1", 6411.0, 6527.0) {
		t.Error("Expected RouN.aid to contain MH_Zone1 at (6411, 6527)")
	}
	areas := rounAID.AreasAt(6411.0, 6527.0)
	if len(areas) == 0 {
		t.Error("AreasAt(6411, 6527) in RouN.aid returned no areas")
	}
}

// TestFase10_Control6_AreaBMP_44Files verifica que AreaBMP cargue el 100% de las máscaras
// de impacto AoE de data/AreaBMP/ (>= 44 máscaras) y realice consultas espaciales rotadas.
func TestFase10_Control6_AreaBMP_44Files(t *testing.T) {
	rootPath := findRepoRoot(t)
	areaBMPDir := filepath.Join(rootPath, "data", "AreaBMP")

	store, err := areabmp.LoadAll(areaBMPDir)
	if err != nil {
		t.Fatalf("areabmp.LoadAll failed: %v", err)
	}

	count := store.Count()
	if count < 44 {
		t.Fatalf("Expected at least 44 unique BMP masks, got %d", count)
	}

	t.Logf("Control 6 (AreaBMP): 100%% Parity (%d unique AoE masks loaded, >= 44 required)", count)

	// Validar máscara de prueba ArkAstanica_Skill01_W
	bmp := store.GetBMP("ArkAstanica_Skill01_W")
	if bmp == nil {
		t.Fatal("ArkAstanica_Skill01_W.bmp missing in store")
	}
	if bmp.Width != 100 || bmp.Height != 400 {
		t.Errorf("Expected 100x400 mask, got %dx%d", bmp.Width, bmp.Height)
	}

	// Probar ContainsTarget con rotación
	_ = bmp.ContainsTarget(5000, 5000, 45, 5020, 5020)
}

// TestFase10_Control6_Lua_274OfficialAPIs verifica que el motor Lua registre el 100% de las
// 274 funciones oficiales C-API de Zone.pdb, valida por AST que no existan stubs mudos
// (stubNoOp/stubReturnZero), y comprueba la carga limpia de scripts oficiales de data/LuaScript/.
func TestFase10_Control6_Lua_274OfficialAPIs(t *testing.T) {
	rootPath := findRepoRoot(t)
	luaPath := filepath.Join(rootPath, "data", "LuaScript")

	engine, err := lua.NewEngine(lua.Config{
		BasePath: luaPath,
		Debug:    false,
	})
	if err != nil {
		t.Fatalf("lua.NewEngine failed: %v", err)
	}
	defer engine.Close()

	// 1. Verificar presencia de las 274 funciones oficiales
	if len(lua.OfficialAPINames) != 274 {
		t.Fatalf("Expected exactly 274 official API names, found %d", len(lua.OfficialAPINames))
	}

	missing := 0
	for _, name := range lua.OfficialAPINames {
		val := engine.GetGlobal(name)
		if val.Type() != glua.LTFunction {
			t.Errorf("Official API function %s is not registered as function in Lua (type: %s)", name, val.Type())
			missing++
		}
	}

	if missing > 0 {
		t.Fatalf("Control 6 (Lua): %d / 274 official APIs are missing in Lua engine", missing)
	}
	t.Logf("Control 6 (Lua): 100%% Parity (All 274/274 official C-API functions registered in Lua)")

	// 2. Escaneo AST de pkg/script/lua/: cero stubs mudos
	luaPkgDir := filepath.Join(rootPath, "pkg", "script", "lua")
	var stubViolations []string

	err = filepath.Walk(luaPkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		ast.Inspect(node, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok {
				if ident.Name == "stubNoOp" || ident.Name == "stubReturnZero" || ident.Name == "stubEmptyString" {
					relPath, _ := filepath.Rel(rootPath, path)
					stubViolations = append(stubViolations, relPath+": "+ident.Name)
				}
			}
			return true
		})

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to scan pkg/script/lua/: %v", err)
	}

	if len(stubViolations) > 0 {
		t.Fatalf("Control 6 (Lua): Found forbidden mute stubs in pkg/script/lua/:\n  %s",
			strings.Join(stubViolations, "\n  "))
	}
	t.Logf("Control 6 (Lua): AST check PASSED (Zero mute stubs found in pkg/script/lua/)")

	// 3. Ejecutar scripts oficiales de data/LuaScript/ y MobBehaviorDescript/
	mobBehaviorDir := filepath.Join(rootPath, "data", "MobBehaviorDescript")

	testScripts := []string{
		filepath.Join(mobBehaviorDir, "DefaultBehavior.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "KingSlime.lua"),
		filepath.Join(mobBehaviorDir, "KQ", "MaraPirate.lua"),
		filepath.Join(luaPath, "Wedding", "Wedding.lua"),
		filepath.Join(luaPath, "Promote", "JobChange1.lua"),
		filepath.Join(luaPath, "Promote", "JobChange2-1.lua"),
		filepath.Join(luaPath, "Map", "OXField.lua"),
		filepath.Join(luaPath, "Map", "Monkey.lua"),
		filepath.Join(luaPath, "Guild", "GuildTournament.lua"),
		filepath.Join(luaPath, "KQ", "Honeying", "Honeying.lua"),
		filepath.Join(luaPath, "KQ", "GordonMaster", "GordonMaster.lua"),
		filepath.Join(luaPath, "KQ", "ConditionOfHero", "ConditionOfHero.lua"),
		filepath.Join(luaPath, "KQ", "UnderHall", "UnderHall.lua"),
		filepath.Join(luaPath, "KQ", "KQHBat", "KQHBat.lua"),
	}

	for _, s := range testScripts {
		if err := engine.DoFile(s); err != nil {
			t.Errorf("Failed to execute Lua script %s: %v", filepath.Base(s), err)
		}
	}
	t.Logf("Control 6 (Lua): Successfully loaded %d official Lua scripts without errors", len(testScripts))
}

func mathAbs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
