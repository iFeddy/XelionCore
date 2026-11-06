# Fase 07 — Integración Transversal `WorldService <-> ZoneService`, In-Game Chunk 100% MySQL y Auditoría Forense Final (`8/8 PASS`)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 1-Bugs #7, #15 y #16, 5.1, 6 y 9 completa**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 01, 02, 02A, 03, 04, 05 y 06 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 07 final (`docs/plan/FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md`) leyendo obligatoriamente las Secciones 1 (Bugs #7, #15 y #16), 5.1, 6 y la Sección 9 completa de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados (`"RouN"` prohibido como fallback en sesión), CERO listas vacías artificiales en el In-Game Chunk y 100% persistencia entre procesos `WorldService` y `ZoneService` a través de MySQL `xelion_online`.
1) En `cmd/ZoneService/main.go` y `pkg/server/zone/`, registrá el nodo de zona al iniciar (`CALL usp_WorldZone_UpsertNode` + latido `CALL usp_WorldZone_Heartbeat`), consumí el traspaso de un solo uso con `CALL usp_Session_ConsumeZoneTransfer` en `PROTO_NC_MAP_LOGIN_REQ`, cargá desde MySQL + `data/World/ClassParam/` los 16 paquetes reales del In-Game Chunk (base, shape, quests doing/done/read/repeat, passives, charged buffs, active skills, game options, items en boxes 8/9/12/15 con options y enchants, titles, collection cards, states, cooltimes y restExp), conectá los portales inter-zona (`LinkTable.json` / `MapLinkPoint.shn`) emitiendo nueva transferencia en `tActiveZoneTransfer`, y persistí el estado íntegro al desconectar (`CALL p_Char_Logout` + states + cooltimes + restExp).
2) Creá `pkg/audit/world_character_completeness_test.go` (Controles 6, 7 y 8 integrados con 1..5) y el script ejecutable `scripts/verify_world_character_100.sh`, y no te detengas hasta que los 8 Controles Obligatorios de la Sección 9 pasen al 100% (`[PASS]`, `exit 0`).
Al terminar, marcá Fase 07 como "Completado" en `docs/plan/README.md` y hacé el `git commit` final en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Registro de Topología y Traspaso Multi-Proceso `WorldService <-> ZoneService` (`cmd/ZoneService/main.go` + `pkg/server/zone/handler.go`):**
   - Al arrancar `ZoneService`, registrar su `(worldNo, zoneNo, zoneIP, zonePort)` en `tWorldZoneNode` (`CALL usp_WorldZone_UpsertNode`) y actualizar cada 30 segundos su latido y recuento de jugadores (`CALL usp_WorldZone_Heartbeat`).
   - En `PROTO_NC_MAP_LOGIN_REQ` (`0x1801`):
     * Mantener la validación criptográfica de los `49` `.shn` (`ChecksumEngine.VerifyDataFiles`) y persistencia en `tHackScanLog` ante mismatch (`0x0147`).
     * Consumir el token de transferencia desde MySQL (`CALL usp_Session_ConsumeZoneTransfer`), garantizando que `cmd/WorldService` y `cmd/ZoneService` funcionen en binarios separados sin compartir memoria.
     * Eliminar el fallback hardcodeado `if mapID == "" { mapID = "RouN" }`: el mapa proviene obligatoriamente de `sTargetMap` / `tCharacter.sLoginZone` validado contra `data/World/Map/FieldList.json` (`MapInfo.shn`).

2. **Carga Real 100% MySQL + `data/` de los 16 Paquetes del In-Game Chunk y Portales Inter-Zona (`Bug #16`):**
   - Alimentar cada paquete de `sendInGameChunk` desde los Stored Procedures reales de `xelion_online` y `data/`:
     * `0x1038` (`PROTO_NC_CHAR_CLIENT_BASE_CMD`) + `0x1039` (`SHAPE_CMD`): `CALL p_Char_GetLoginData` + `CALL p_Char_GetShape` + cálculo exacto de stats de `data/World/ClassParam/<Class>.json` + puntos redistribuidos.
     * `0x103A` (`QUEST_DOING_CMD`), `0x103B` (`QUEST_DONE_CMD`), `0x10CE` (`QUEST_READ_CMD`), `0x10D7` (`QUEST_REPEAT_CMD`): `CALL p_Quest_GetAll`, `CALL p_Quest_GetAllDone` y `CALL p_Quest_GetRepeat` (`tQuest` y `tQuestTimes`, incluyendo `starter_quests` creadas en Fase 05).
     * `0x103E` (`PASSIVE_CMD`) y `0x103D` (`SKILL_CMD`): `CALL p_Skill_GetAll` clasificando pasivas (`PassiveSkill.shn`) y activas (`ActiveSkill.shn`).
     * `0x104A` (`CHARGEDBUFF_CMD`): `CALL p_Char_GetChargedBufferAll` (`tCharacterChargedBuffer` + `ChargedEffect.shn`).
     * `0x1048` (`GAME_CMD`): `CALL p_Char_GetOptions`, `CALL p_Char_GetStateAll`, `CALL p_Character_getItemActionCoolTime`, `CALL usp_Char_GetRestExpExecTime`.
     * `0x1047` (`ITEM_CMD` — Boxes `8` Equip, `9` Inven, `12` House, `15` ActionItem): `CALL p_Item_GetList` + `CALL p_Item_GetOptions` + `CALL p_Item_GetEnchant` mapeados contra `ItemInfo.shn`.
     * `0x1049` (`CHARTITLE_CMD`): `CALL p_Char_GetTitleAll` (`tCharacterTitle` + `CharacterTitleData.shn`).
     * `Collection Cards`: `CALL p_Collection_Get` (`tCollectionCard` + `CollectCardView.shn`).
   - **Saltos Inter-Mapa / Inter-Zona (`useGate`, `TownPortal`, `WingFly`, `RecallCoord.json`):** si el `zoneNumber` destino en `FieldList.json` pertenece a otra instancia de zona, resolver endpoint con `CALL usp_WorldZone_ResolveMapEndpoint`, persistir posición con `CALL p_Char_SaveLink`, emitir token con `CALL usp_Session_IssueZoneTransfer` y enviar paquete de redirección de zona al cliente.
   - **Logout Transaccional Completo (`OnClientDisconnect`):** persistir `HP`, `SP`, `HPS`, `SPS`, `EXP`, `Money`, `sLoginZone`, `X`, `Y`, `D` y `nPlayMin` mediante `CALL p_Char_Logout`, junto con `CALL p_Char_SetStateWithRestTime`, `CALL p_Character_setItemActionCoolTime`, `CALL usp_Char_SetRestExpExecTime`, más cierre de sesión de cuenta (`usp_Log_setLogout` / `usp_Session_CloseAndGetPlayMin`).

3. **Suite Forense Final (`pkg/audit/world_character_completeness_test.go`) y Script Maestro (`scripts/verify_world_character_100.sh`):**
   - Implementar los **Controles 6, 7 y 8** e integrar los **8 Controles Obligatorios de la Sección 9**:
     * **Control 1:** `114/114` tablas `.bak` + `2` tablas operativas S2S en `information_schema.TABLES`, cero filas `*_Template` en `tCharacter`, y JSONs de inicio completos.
     * **Control 2:** `425/425` Stored Procedures en `information_schema.ROUTINES` + verificación pool-safe sin `sql.ErrNoRows`.
     * **Control 3:** Escaneo AST de Go — cero SQL inline en repositorios de `Character`/`ZoneTopology` y cero mapas de transferencia en memoria.
     * **Control 4:** Paridad 100% de OpCodes y `sizeof` contra `WorldManager.md` y `Character.md`, cero alias prohibidos y cero carpeta `pkg/protocol/pdb/`.
     * **Control 5:** Cobertura `167/167` handlers C2S de `WorldManager.pdb` + S2S/persistencia de `Character.pdb`.
     * **Control 6:** Carga y verificación criptográfica real de los **54 `.shn` de `Character.exe`**, los **71 `.shn` de `WorldManager.exe`** (`74` únicos) y los **49 `.shn` de `ChecksumEngine`**.
     * **Control 7:** Prueba E2E integral forense punta a punta: `LoginService -> WorldService (LOGINWORLD 64B -> AVATAR_CREATE -> Opciones/Guild/Friend/Party -> CHAR_LOGIN) -> ZoneService (MAP_LOGIN con 49 checksums -> recepción verificada de los 16 paquetes del In-Game Chunk -> movimiento de ítem -> subida de nivel/EXP -> desconexión -> re-verificación de columnas exactas en MySQL xelion_online)`.
     * **Control 8:** Escaneo AST de **cero datos hardcodeados** en `pkg/server/world/` y `pkg/server/zone/` (prohibiendo slices de `BadNameFilter`, `loadDefaults()` en KQ y fallbacks `"RouN"`).
   - Crear `scripts/verify_world_character_100.sh` con permisos de ejecución (`chmod +x`) que corra los 8 controles secuencialmente, imprima el resumen `[PASS]` `8/8` y devuelva `exit 0`.

---

## Criterio de Cierre Definitivo de la Etapa 2 (`WorldService` + `Character`)

1. `go build ./...` -> `OK`.
2. `go test -count=1 ./...` -> `PASS` en el 100% de los paquetes del repositorio (incluyendo Etapa 1: `./scripts/verify_login_account_100.sh` sigue en `exit 0`).
3. `./scripts/verify_world_character_100.sh` -> **`8/8` controles en `[PASS]`, `exit 0`**.
4. Marcar Fase 07 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(core): Fase 07 - integracion completa WorldService/Character con ZoneService y auditoria 8/8 PASS"`
