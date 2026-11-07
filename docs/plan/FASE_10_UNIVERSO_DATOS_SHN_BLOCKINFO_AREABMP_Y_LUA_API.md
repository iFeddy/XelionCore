# Fase 10 — Universo Completo de Datos Estáticos de `ZoneService`: `180` SHNs Reales, Geometría `BlockInfo` (`.shbd`/`.sbi`/`.aid`), `AreaBMP` y Runtime `Lua` C-API (`274` Símbolos)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0-Regla #5, 1-Bugs #5, #6 y #7, 3.1, 3.2, 3.3, 3.4 y 9-Control 6**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 y 09 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 10 (`docs/plan/FASE_10_UNIVERSO_DATOS_SHN_BLOCKINFO_AREABMP_Y_LUA_API.md`) leyendo obligatoriamente las Secciones 0 (Regla #5), 1 (Bugs #5, #6 y #7), 3.1, 3.2, 3.3, 3.4 y 9 (Control 6) de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md` y respetando estrictamente `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados, CERO omisiones de archivos en `data/` y CERO stubs mudos en el runtime Lua.
1) Expandí `pkg/server/zone/data_store.go` para cargar en RAM e indexar las 180 tablas `.shn` reales en disco (`194` rutas en `data/` y `data/View/`, manejando explícitamente como opcionales las 5 referencias ausentes en `AstraServer`: `BMPIndex.shn`, `DiecGame.shn`, `MoveMain.shn`, `MoverUpgrade.shn`, `Quest.shn`) junto con la validación criptográfica de los 49 `.shn` de `ChecksumEngine`.
2) Creá el paquete oficial `pkg/data/blockinfo` para leer los 301 archivos binarios de `data/BlockInfo/` (`158` `.shbd` con walkability/LoS Bresenham `MapBlockInformation`, `33` `.sbi` `MapDoorArray`/`MapDoorBlockElement` con apertura y cierre dinámico de puertas `DoorAction`, `23` `.aid` `AreaInfoData` círculos y rectángulos, `76` `.bdt`, `8` `.shab`, `2` `.shmd`, `1` `.shad`) y creá `pkg/data/areabmp` para leer las `44` máscaras `.bmp` (`CFsAreaBMP`) de `data/AreaBMP/`, integrando ambos en `DataStore`.
3) Refactorizá `pkg/script/lua/` definiendo la interfaz `WorldBridge` e implementando el registro completo de las 274 funciones C-API `lua_State*` exactas de `Zone.pdb` (Sección 3.4), eliminando los 32 `stubNoOp`/`stubReturnZero` mudos y validando que todos los scripts de `data/LuaScript/` carguen limpiamente.
4) Creá `pkg/audit/fase10_zone_datastore_geometry_lua_test.go` (Control 6 completo de Etapa 3). Al terminar, marcá Fase 10 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 11.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Carga Completa de las `180` Tablas `.shn` Reales (`194` Rutas) en `pkg/server/zone/data_store.go` (`Sección 3.1`):**
   - Extender `DataStore` con un catálogo indexado (`SHNTables map[string]*shn.File` + vistas tipadas rápidas para tablas críticas de gameplay: `UpgradeInfo`, `AccUpGradeInfo`, `BraceletUpgrade`, `ItemSocket`, `ItemMix`, `ItemDismantle`, `Produce`, `RawInfo`, `CollectInfo`, `QuestData`, `QuestDialog`, `SubAbState`, `AbState`, `MoverMain`, `MoverSub`, `Pet*`, `MiniHouse*`, `DiceGame`, `SlotMachine`, `CharacterTitleData`, `CollectCard*`, `WeaponTitleData`, `BoothSearchItemList`, `FieldLvCondition`, `WorldDrop`, `AdminLvSet`, `ErrorCodeTable`, etc.).
   - Verificar obligatoriamente la carga exitosa de los **`180` basenames `.shn` reales** (`194` archivos entre `data/` y `data/View/`) y registrar como opcionales documentadas únicamente las `5` cadenas inexistentes en `AstraServer` (`BMPIndex.shn`, `DiecGame.shn`, `MoveMain.shn`, `MoverUpgrade.shn`, `Quest.shn`).

2. **Nuevo Cargador Oficial de Geometría y Colisión `pkg/data/blockinfo` (`301` Archivos en `data/BlockInfo/` — `Sección 3.3`):**
   - Implementar lectura binaria fiel a `Zone.pdb`:
     * **`.shbd` (`158` mapas — `MapBlock::MapBlockInformation`):** lectura de `rowBytes (uint32)` y `heightTiles (uint32)` (`widthTiles = rowBytes * 8`), validación `len(payload) == int(rowBytes*heightTiles)`, métodos O(1) `IsMoveBlocked(x, y)`, `IsAttackBlocked(x1, y1, x2, y2)` (trazado de línea Bresenham), `FindNearestWalkable(x, y, maxRadius)` y conversión de coordenadas mundo <-> celda de bloque según escala oficial de `Zone.exe`.
     * **`.sbi` (`33` mapas — `MapBlock::MapDoorArray` / `MapDoorBlockElement`):** carga de bloques dinámicos con nombre (`DoorBlock01`, etc.), estado por instancia de mapa (`Open`/`Close` vía `DoorAction` y `CloseAllDoors`) superpuesto sin mutar el `.shbd` compartido en RAM.
     * **`.aid` (`23` mapas — `AreaInfoData`):** carga de zonas nombradas (`AE_Circle` y `AE_Square`, ej. `MH_Zone1`) con `Contains(areaName, x, y)` y `AreasAt(x, y)`.
     * **`.bdt` (`76`), `.shab` (`8`), `.shmd` (`2`), `.shad` (`1`):** parsers y catálogo completo accesibles desde `blockinfo.Store`.

3. **Nuevo Cargador Oficial de Máscaras AoE `pkg/data/areabmp` (`44` Archivos en `data/AreaBMP/` — `Sección 3.3`):**
   - Parsear los `44` archivos Windows BMP (`8 bpp`) de `data/AreaBMP/` (`CFsAreaBMP`), exponiendo consulta espacial rotada según el ángulo/dirección del lanzador (`ContainsTarget(casterX, casterY, dirDeg, targetX, targetY)`).

4. **Motor Oficial `pkg/script/lua/` con las `274` Funciones C-API `lua_State*` de `Zone.pdb` (`Sección 3.4`):**
   - Diseñar la interfaz limpia `HostBridge` (desacoplada de ciclos de importación) que enlaza cada VM Lua por instancia/mapa con objetos, monstruos, puertas `.sbi`, `AbState`, misiones, KQ, mascotas y `GameLog`.
   - Registrar e implementar las **`274` funciones exactas de `Zone.pdb`** listadas en la Sección 3.4 (`246` `c*` + `22` `ZoneServer_zs_*` + `6` base + helpers de `common.lua`), eliminando los 32 stubs mudos de [`pkg/script/lua/engine.go`](../../pkg/script/lua/engine.go).

5. **Auditoría Automatizada (`pkg/audit/fase10_zone_datastore_geometry_lua_test.go` — Control 6):**
   - Verifica carga real de `180/180` `.shn` (`194` rutas) + `49/49` checksums, carga íntegra de los `301` archivos de `data/BlockInfo/` y `44` `.bmp` de `data/AreaBMP/` (probando colisión real en `Rou.shbd`/`Adl.shbd` y apertura/cierre en `AdlF.sbi`), y presencia activa en Lua de las `274/274` funciones C-API de `Zone.pdb` + carga de `data/LuaScript/`.

---

## Criterio de Cierre de la Fase 10

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/data/... ./pkg/script/... ./pkg/server/zone/... ./pkg/audit -run "TestFase10"` -> `PASS`.
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 10 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(data,lua,zone): Fase 10 - carga 180 SHNs reales, paquetes blockinfo/areabmp y 274 C-APIs Lua de Zone.pdb"`
