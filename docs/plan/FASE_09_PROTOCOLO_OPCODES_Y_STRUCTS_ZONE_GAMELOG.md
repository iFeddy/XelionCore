# Fase 09 — Protocolo Oficial 100% `Zone.pdb` y `GameLog.md`: `245` OpCodes C2S, S2S y `1.293 + 44` Estructuras C++ con Paridad `sizeof`

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0-Regla #6, 1-Bugs #3 y #4, 2.4, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3, 7 y 9-Control 4**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md), [`docs/Servicios/GameLog.md`](../Servicios/GameLog.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fase 08 completada sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 09 (`docs/plan/FASE_09_PROTOCOLO_OPCODES_Y_STRUCTS_ZONE_GAMELOG.md`) leyendo obligatoriamente las Secciones 0 (Regla #6), 1 (Bugs #3 y #4), 2.4, 4.1..4.4, 5.1..5.3, 7 y 9 (Control 4) de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md`, `docs/Servicios/GameLog.md` y `AGENTS.md`.
Regla Mandatoria: Nomenclatura 100% idéntica a los PDBs (`PROTO_NC_*`) directamente en `pkg/protocol/` (prohibido crear subcarpetas `pkg/protocol/pdb/`), CERO omisiones y paridad binaria exacta `sizeof` al byte.
1) Registrá en `pkg/protocol/` el 100% de los 245 OpCodes C2S de `Zone.exe` (`protocolstore(ShinePlayer)` en `0x004D3510`, respetando `protocol = (department << 10) | command`) y todos los OpCodes S2S asociados a `wms_NC_*` (138), `gds_NC_*` (196) y `GameLog`, erradicando cualquier alias no oficial restante en `pkg/protocol/` y `pkg/server/zone/handler.go` (como `PROTO_NC_ACT_REST_START_REQ` / `PROTO_NC_ACT_REST_STOP_REQ` -> `PROTO_NC_ACT_PITCHTENT_REQ` / `PROTO_NC_ACT_FOLDTENT_REQ`).
2) Implementá y serializá en `pkg/protocol/` las 44 estructuras `PROTO_NC_LOG_*` de `docs/Servicios/GameLog.md` y las 1.293 estructuras oficiales `PROTO_NC_*` (+ las 9 plantillas `ShineSendType<PROTO_NC_*>`) de `docs/Servicios/Zone.md` (excluyendo únicamente `PROTO_NC_OPTOOL_CONNECT_BRIF_ACK` y los 40 tipos internos C++ STL de la Sección 5.3).
3) Creá `pkg/audit/fase09_zone_gamelog_protocol_test.go` (Control 4 de Etapa 3) que parsee directamente `docs/Servicios/Zone.md` y `docs/Servicios/GameLog.md` y valide automáticamente el 100% de los nombres y `sizeof` en bytes. Al terminar, marcá Fase 09 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 10.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Cobertura 100% de OpCodes Oficiales de `Zone.exe` y Limpieza de Alias (`pkg/protocol/opcodes*.go` y `pkg/server/zone/handler.go`):**
   - Verificar y declarar el **100% de las `245` constantes de protocolo C2S** listadas en la Sección 4.1 de [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (incluyendo los callbacks especiales `0x0804` y `0x0805` de `iac_FoodEat` con constantes oficiales explícitas documentadas) más todas las constantes S2S de `Zone.md` y `GameLog.md`, calculadas estrictamente mediante `(department << 10) | command`.
   - Eliminar de `pkg/protocol/` y `pkg/server/zone/handler.go` cualquier alias no oficial de emuladores privados (`PROTO_NC_ACT_REST_START_REQ`, `PROTO_NC_ACT_REST_STOP_REQ`, `PROTO_NC_CHAR_INGAME_REQ`, etc.) reemplazándolos por las constantes exactas del PDB (`PROTO_NC_ACT_PITCHTENT_REQ`, `PROTO_NC_ACT_FOLDTENT_REQ`, `PROTO_NC_MAP_LOGINCOMPLETE_CMD`, etc.).

2. **Las `44` Estructuras `PROTO_NC_LOG_*` de `docs/Servicios/GameLog.md` (`pkg/protocol/gamelog_structs.go`):**
   - Implementar las **`44` estructuras exactas** listadas en la Sección 2.4 (`PROTO_NC_LOG_GAME_DATA_TYPE_0` a `TYPE_9`, `TYPE_A` a `TYPE_O` y las `19` estructuras especializadas `PROTO_NC_LOG_GAME_*` / `PROTO_NC_LOG_WEDDING_*` / `PROTO_NC_LOG_REGENLOCATESAVE_CMD`) con sus lectores/escritores binarios Little-Endian `#pragma pack(1)` y conectarlas con `GameLogDispatcher` de la Fase 08.

3. **Las `1.293` Estructuras `PROTO_NC_*` y `9` `ShineSendType<...>` de `docs/Servicios/Zone.md` (`pkg/protocol/zone_structs_*.go`):**
   - Organizar modularmente por departamento en `pkg/protocol/` (sin crear subpaquetes paralelos) el **100% de las `1.293` estructuras `PROTO_NC_*`** de la Sección 5.1 (incluyendo sub-estructuras anidadas `unnamed-type-*` y tipos compuestos oficiales `SHINE_XY_TYPE`, `MAPPOS`, `CHAR_PARAMETER_DATA`, `PROTO_EQUIPMENT`, `ITEM_INVEN`, `ItemOptionStorage::Element`, etc., reutilizando sin colisión las estructuras ya creadas en Etapa 1 y Etapa 2 cuando comparten definición idéntica en el PDB).
   - Para paquetes con arrays flexibles (`zero-length` o trailing arrays en C++), proveer métodos exactos de `FixedWireSize()` (idéntico al `sizeof` estático de `Zone.pdb`) más serialización/deserialización completa de sus slices dinámicos (`SkillDamage`, `CardBookmarkInfo`, `CardInform`, `MobOfQuest`, etc.).

4. **Auditoría Automatizada (`pkg/audit/fase09_zone_gamelog_protocol_test.go` — Control 4):**
   - Parsea en tiempo de test [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`docs/Servicios/GameLog.md`](../Servicios/GameLog.md) verificando:
     * `245/245` OpCodes C2S registrados con valor `(dept << 10) | cmd` exacto.
     * `44/44` estructuras de `GameLog.md` presentes y con `sizeof` idéntico al byte.
     * `1.293/1.293` estructuras `PROTO_NC_*` (+ `9/9` `ShineSendType`) de `Zone.md` presentes y con `sizeof` idéntico al byte.
     * Ausencia total de alias privados y de cualquier carpeta `pkg/protocol/pdb/`.

---

## Criterio de Cierre de la Fase 09

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/protocol/... ./pkg/audit -run "TestFase09"` -> `PASS`.
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 09 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(protocol): Fase 09 - 100% OpCodes y 1293+44 structs oficiales de Zone.pdb y GameLog con paridad sizeof"`
