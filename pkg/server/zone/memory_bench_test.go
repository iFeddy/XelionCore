package zone

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func TestMemory_ZoneMonster_LazyMaps(t *testing.T) {
	mob := NewZoneMonster(100, 1, "Slime", "RouN", 1000, 1000, 0, 1, 100, 50, 10, 10, 5, 10, 10, 10, 20)
	if mob.Aggro != nil || mob.DamageContrib != nil || mob.ActiveStates != nil || mob.FiredTriggers != nil {
		t.Fatalf("expected all dynamic maps to be nil on newly spawned ZoneMonster")
	}

	// Mutadores deben inicializar perezosamente solo el mapa requerido
	mob.AddAggro(1, 50)
	if mob.Aggro == nil || mob.Aggro[1] != 50 {
		t.Fatalf("AddAggro did not lazily initialize Aggro map")
	}
	if mob.DamageContrib != nil || mob.ActiveStates != nil || mob.FiredTriggers != nil {
		t.Fatalf("AddAggro unexpectedly allocated unrelated maps")
	}

	mob.RecordDamage(1, 42, 25)
	if mob.DamageContrib == nil || mob.DamageContrib[42] != 25 {
		t.Fatalf("RecordDamage did not lazily initialize DamageContrib map")
	}

	mob.AddAbState(&ActiveAbState{ID: 7, Level: 1})
	if mob.ActiveStates == nil || mob.GetAbState(7) == nil {
		t.Fatalf("AddAbState did not lazily initialize ActiveStates map")
	}

	mob.ResetAggro()
	if mob.Aggro != nil {
		t.Fatalf("ResetAggro must release Aggro map to nil")
	}

	mob.RemoveAbState(7)
	if mob.ActiveStates != nil {
		t.Fatalf("RemoveAbState must release ActiveStates map to nil when empty")
	}

	// Verificar liberación en respawnMonster y tickMobReturn
	zm := NewZoneMap(1, "RouN", "Roumen")
	defer zm.Close()
	mob.AddAggro(2, 10)
	mob.RecordDamage(2, 99, 15)
	mob.AddAbState(&ActiveAbState{ID: 9, Level: 1})
	mob.FiredTriggers = map[int]bool{0: true}

	zm.respawnMonster(mob)
	if mob.Aggro != nil || mob.DamageContrib != nil || mob.ActiveStates != nil || mob.FiredTriggers != nil {
		t.Fatalf("respawnMonster must release all dynamic maps back to nil")
	}

	mob.AddAggro(3, 20)
	mob.RecordDamage(3, 77, 10)
	mob.AddAbState(&ActiveAbState{ID: 11, Level: 1})
	mob.FiredTriggers = map[int]bool{1: true}
	mob.X, mob.Y = mob.SpawnX, mob.SpawnY
	mob.LastMoveTick = time.Time{}

	zm.tickMobReturn(mob, time.Now())
	if mob.Aggro != nil || mob.DamageContrib != nil || mob.ActiveStates != nil || mob.FiredTriggers != nil {
		t.Fatalf("tickMobReturn at spawn must release all dynamic maps back to nil")
	}
}

func TestMemory_ZoneMap_GetAllMonsters_ZeroAllocs(t *testing.T) {
	zm := NewZoneMap(1, "RouN", "Roumen")
	defer zm.Close()

	for i := 0; i < 50; i++ {
		zm.AddMonster(NewZoneMonster(uint16(100+i), 1, "Slime", "RouN", 1000, 1000, 0, 1, 100, 50, 10, 10, 5, 10, 10, 10, 20))
	}

	// Primera llamada construye el snapshot cacheado
	snap := zm.GetAllMonsters()
	if len(snap) != 50 {
		t.Fatalf("expected 50 monsters, got %d", len(snap))
	}

	allocs := testing.AllocsPerRun(100, func() {
		s := zm.GetAllMonsters()
		if len(s) != 50 {
			t.Fatalf("unexpected monster count %d", len(s))
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs/op on cached GetAllMonsters(), got %.2f", allocs)
	}

	// Verificar invalidación al remover un monstruo
	zm.RemoveMonster(100)
	snapAfter := zm.GetAllMonsters()
	if len(snapAfter) != 49 {
		t.Fatalf("expected 49 monsters after removal, got %d", len(snapAfter))
	}
}

func TestMemory_LazyLuaVM_And_PrecompiledProto(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	ds, err := LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData failed: %v", err)
	}

	zm := NewZoneMap(1, "RouCos01", "Sand Beach")
	defer zm.Close()
	spawned := zm.LoadMobsFromRegen(ds)
	if spawned == 0 {
		t.Fatalf("expected RouCos01 to spawn monsters from MobRegen")
	}

	// Un mapa estándar sin script KQ ni AIScriptFunc NO debe crear VM Lua ansiosamente
	zm.mu.RLock()
	hasEngine := zm.luaEngine != nil
	zm.mu.RUnlock()
	if hasEngine {
		t.Fatalf("expected RouCos01 luaEngine to remain nil until lazily invoked")
	}

	// Al invocar RunMobBehaviorLua, la VM debe inicializarse bajo demanda usando bytecode precompilado
	mobs := zm.GetAllMonsters()
	if len(mobs) == 0 || mobs[0].Template == nil {
		t.Fatalf("expected spawned monster to carry shared *MobTemplate reference")
	}
	_ = zm.RunMobBehaviorLua(mobs[0], "DefaultBehavior")

	zm.mu.RLock()
	eng := zm.luaEngine
	zm.mu.RUnlock()
	if eng == nil {
		t.Fatalf("expected luaEngine to be lazily initialized by RunMobBehaviorLua")
	}
}

func TestMemory_DataStore_Deduplication_And_Profile(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	ds, err := LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData failed: %v", err)
	}

	// Verificar reutilización exacta por puntero (0 lecturas duplicadas)
	if ds.ActiveSkillSHN != ds.GetSHN("ActiveSkill") {
		t.Fatalf("ActiveSkillSHN pointer differs from ds.SHNTables")
	}
	if ds.PassiveSkillSHN != ds.GetSHN("PassiveSkill") {
		t.Fatalf("PassiveSkillSHN pointer differs from ds.SHNTables")
	}
	if ds.MobInfoSHN != ds.GetSHN("MobInfo") {
		t.Fatalf("MobInfoSHN pointer differs from ds.SHNTables")
	}

	// Verificar paridad 49/49 de ChecksumEngine generado desde ds.SHNTables
	refChk, err := shn.NewChecksumEngine(dataDir, filepath.Join(dataDir, "View"))
	if err != nil {
		t.Fatalf("NewChecksumEngine reference failed: %v", err)
	}
	if ok, badIdx := ds.ChkEngine.VerifyDataFiles(refChk.GetChecksums()); !ok {
		t.Fatalf("DataStore ChkEngine mismatch at index %d (%s)", badIdx, shn.ChecksumFileList[badIdx])
	}

	// Verificar BlockInfo: 301 archivos indexados y .bdt lazy en reposo (0 bytes cargados)
	if ds.BlockInfo.TotalCount() != 301 {
		t.Fatalf("expected BlockInfo TotalCount=301, got %d", ds.BlockInfo.TotalCount())
	}
	if len(ds.BlockInfo.BDT) != 76 {
		t.Fatalf("expected 76 indexed BDT keys, got %d", len(ds.BlockInfo.BDT))
	}
	var sampleBDTKey string
	for k, v := range ds.BlockInfo.BDT {
		if len(v) != 0 {
			t.Fatalf("expected lazy BDT %q to hold 0 bytes before GetBDT, got %d", k, len(v))
		}
		if sampleBDTKey == "" {
			sampleBDTKey = k
		}
	}
	if bdt := ds.BlockInfo.GetBDT(sampleBDTKey); len(bdt) == 0 {
		t.Fatalf("expected lazy GetBDT(%q) to load non-empty BDT bytes", sampleBDTKey)
	}

	// Verificar AreaBMP 1D contiguo
	if ds.AreaBMP.Count() == 0 {
		t.Fatalf("expected non-empty AreaBMP store")
	}
	for _, bmp := range ds.AreaBMP.BMPs {
		if len(bmp.Data) != bmp.Width*bmp.Height {
			t.Fatalf("BMP %s Data len=%d, want %d", bmp.Name, len(bmp.Data), bmp.Width*bmp.Height)
		}
		if len(bmp.Pixels) != bmp.Height || (bmp.Height > 0 && len(bmp.Pixels[0]) != bmp.Width) {
			t.Fatalf("BMP %s Pixels dimensions mismatch", bmp.Name)
		}
		break
	}

	mm := NewMapManager(ds.FieldStore)
	mm.SetBlockStore(ds.BlockInfo)
	mm.SetNPCStore(ds.NPCStore, ds.MobInfoSHN)
	mm.SetDataStore(ds)
	loadedMaps := mm.PreloadZoneMaps()
	defer mm.StopAll()

	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("Profile after LoadAllData + PreloadZoneMaps(%d maps): HeapAlloc=%.2f MB, HeapInuse=%.2f MB, HeapObjects=%d",
		loadedMaps,
		float64(ms.HeapAlloc)/(1024*1024),
		float64(ms.HeapInuse)/(1024*1024),
		ms.HeapObjects,
	)
}

func BenchmarkLoadAllDataAndPreloadMaps(b *testing.B) {
	dataDir := filepath.Join("..", "..", "..", "data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shn.ClearReadCache()
		ds, err := LoadAllData(dataDir)
		if err != nil {
			b.Fatalf("LoadAllData failed: %v", err)
		}
		mm := NewMapManager(ds.FieldStore)
		mm.SetBlockStore(ds.BlockInfo)
		mm.SetNPCStore(ds.NPCStore, ds.MobInfoSHN)
		mm.SetDataStore(ds)
		mm.PreloadZoneMaps()

		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		b.ReportMetric(float64(ms.HeapAlloc)/(1024*1024), "HeapAlloc_MB")
		b.ReportMetric(float64(ms.HeapInuse)/(1024*1024), "HeapInuse_MB")
		b.ReportMetric(float64(ms.HeapObjects), "HeapObjects")
		mm.StopAll()
	}
}

func BenchmarkTickMonsters100ms_Idle(b *testing.B) {
	zm := NewZoneMap(1, "RouN", "Roumen")
	defer zm.Close()
	for i := 0; i < 200; i++ {
		zm.AddMonster(NewZoneMonster(uint16(100+i), 1, "Slime", "RouN", 1000, 1000, 0, 1, 100, 50, 10, 10, 5, 10, 10, 10, 20))
	}
	_ = zm.GetAllMonsters()
	now := time.Now()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zm.TickMonsters100ms(now)
	}
}
