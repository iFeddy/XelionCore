# Fase 16 — Coordinación S2S `Zone <-> World` (`138` `wms_NC_*`), Multi-Nodo `Zone00`..`Zone04`, Comandos GM (`AdminLvSet.shn`), Cierre `245/245` C2S y Auditoría Forense Final (`8/8 PASS`)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0-Regla #3, 1-Bugs #4, #17 y #18, 4.1, 4.3, 4.4, 6.6 y Sección 9 completa**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md), [`docs/Servicios/GameLog.md`](../Servicios/GameLog.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 a 15 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 16 final (`docs/plan/FASE_16_COORDINACION_S2S_ADMIN_E2E_CLUSTER_Y_AUDITORIA_FINAL.md`) leyendo obligatoriamente las Secciones 0 (Regla #3), 1 (Bugs #4, #17 y #18), 4.1, 4.3, 4.4, 6.6 y la Sección 9 completa de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md`, `docs/Servicios/GameLog.md` y `AGENTS.md`.
Regla Mandatoria: Cobertura 245/245 OpCodes C2S en `ZoneService` con CERO `Unhandled OpCode`, CERO SQL inline, CERO datos o privilegios GM hardcodeados y 100% compatibilidad multi-nodo (`Zone00`..`Zone04`).
1) Completá los handlers C2S restantes para cerrar el 100% de las 245 entradas de `protocolstore(ShinePlayer)` (`PINGTEST_*`, `SERVERPARAMETER_REQ`, `USER_NORMALLOGOUT_CMD`, `USER_REGISENUMBER_REQ`, `USER_USE_BEAUTY_SHOP_CMD`, `CHAR_LOGOUTREADY_CMD`, `CHAR_LOGOUTCANCEL_CMD`, `PARTY_MEMBERINFOREQ_CMD`, `CHARSAVE_UI_STATE_SAVE_REQ`, `FRIEND_UES_FRIEND_POINT_REQ`, `GUILD_STORAGEOPEN_REQ`, `ITEMDB_CLOSE_GUILD_STORAGE_CMD`, `GUILD_RENAME_REQ`, `GUILD_RETYPE_REQ`, `USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ`), implementá la coordinación S2S completa de los 138 flujos `wms_NC_*` (`Zone <-> WorldService`), el modo multi-instancia `Zone00`..`Zone04` (`9016, 9019, 9022, 9025, 9028` según `zoneNumber` de `FieldList.json`) y el despachador oficial de comandos GM gobernado estrictamente por `AdminLvSet.shn` + `tGMItemReward*` + `GameLog` Tipo `L`.
2) Creá `pkg/audit/zone_gamelog_completeness_test.go` y el script ejecutable `scripts/verify_zone_gamelog_100.sh`, y no te detengas hasta que los 8 Controles Obligatorios de la Sección 9 pasen al 100% (`8/8 [PASS]`, `exit 0`), verificando también que `./scripts/verify_login_account_100.sh` (Etapa 1) y `./scripts/verify_world_character_100.sh` (Etapa 2) sigan pasando en `exit 0`.
Al terminar, marcá Fase 16 como "Completado" en `docs/plan/README.md` y hacé el `git commit` final en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Cierre `245/245` C2S, Coordinación S2S `138` `wms_NC_*` y Topología `Zone00`..`Zone04` (`cmd/ZoneService/main.go` + `pkg/server/zone/`):**
   - Implementar los handlers C2S restantes de sesión, diagnóstico y enlace social:
     * `Dept 2`: `PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB` (`0x0812`), `PINGTEST_TOOL_WM_CLIENT_ZONE_DB` (`0x0813`), `SERVERPARAMETER_REQ` (`0x082F`).
     * `Dept 3` y `Dept 4`: `USER_NORMALLOGOUT_CMD` (`0x0C18`), `USER_REGISENUMBER_REQ` (`0x0C23`), `USER_USE_BEAUTY_SHOP_CMD` (`0x0C32`), `CHAR_LOGOUTREADY_CMD` (`0x1071`), `CHAR_LOGOUTCANCEL_CMD` (`0x1072`).
     * `Dept 14`, `16`, `21`, `29` y `45`: `PARTY_MEMBERINFOREQ_CMD` (`0x3848`), `CHARSAVE_UI_STATE_SAVE_REQ` (`0x4025`), `FRIEND_UES_FRIEND_POINT_REQ` (`0x5421`), `GUILD_STORAGEOPEN_REQ` (`0x7491`), `ITEMDB_CLOSE_GUILD_STORAGE_CMD` (`0x74A0`), `GUILD_RENAME_REQ` (`0x74AC`), `GUILD_RETYPE_REQ` (`0x74B4`) y `USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ` (`0xB405`).
   - Implementar la capa de coordinación de los **`138` flujos `wms_NC_*`** (`Sección 4.3`) entre `ZoneService` y `WorldService` (sincronización en vivo de Party, Raid, Guild, Guerra de Gremios, Academia, Bodas, Prisión, Avisos Globales y KQ) y parametrizar `cmd/ZoneService/main.go` para operar por `ZONE_ID` (`Zone00`..`Zone04` en puertos `9016, 9019, 9022, 9025, 9028` filtrando `FieldList.json` por `zoneNumber`, o modo unificado de desarrollo).

2. **Sistema Oficial de Comandos GM Gobernado 100% por `AdminLvSet.shn`:**
   - Implementar el despachador de comandos de administración validando cada comando contra el nivel exigido en `AdminLvSet.shn` y `tCharacter.nAdminLevel` (cero nombres o cuentas hardcodeadas), persistiendo entregas en `tGMItemReward` / `tGMItemRewardLog` y registrando auditoría en `GameLog` Tipo `L` (`CALL p_GameLog_Add_TypeL`).

3. **Suite Forense Final (`pkg/audit/zone_gamelog_completeness_test.go`) y Script Maestro (`scripts/verify_zone_gamelog_100.sh`):**
   - Integrar y certificar en verde los **8 Controles Obligatorios de la Sección 9**:
     * **Control 1:** `176/176` tablas oficiales en `xelion_online` (`tGameLog` y `TEMP_QuestComplete` verificadas) y ausencia total de `database/zone_tables.sql`.
     * **Control 2:** `27/27` Stored Procedures de `World00_GameLog` activos en `information_schema.ROUTINES` y verificados concurrentemente sobre `*sql.DB`.
     * **Control 3:** Escaneo AST de Go — cero SQL inline en `pkg/database/mysql_gamelog*.go` y `pkg/server/zone/`.
     * **Control 4:** Paridad 100% de OpCodes y `sizeof` contra `Zone.md` (`1.293` structs `PROTO_NC_*` + `9` `ShineSendType`) y `GameLog.md` (`44` structs `PROTO_NC_LOG_*`).
     * **Control 5:** Cobertura `245/245` handlers C2S en el dispatcher de `ZoneService` + `138` flujos `wms_NC_*` y `196` mapeos `gds_NC_*`.
     * **Control 6:** Carga íntegra en RAM de `180/180` `.shn` reales (`194` rutas) + `49` checksums + `301` archivos `BlockInfo` + `44` `AreaBMP` + `274/274` funciones C-API `lua_State*` de `Zone.pdb`.
     * **Control 7:** Prueba E2E forense integral del clúster completo (`LoginService -> WorldService -> ZoneService -> GameLogDispatcher -> MySQL xelion_online`) cubriendo login, colisión `.shbd`, combate, drop + `AutoPick`, NPC merchant, `Upgrade`, `SafeTrade`, `Quest`, instancia Lua + puerta `.sbi`, logout y verificación directa de filas en `tCharacter`, `tItem`, `tQuest` y `tGameLog`.
     * **Control 8:** Escaneo AST de cero datos hardcodeados en `pkg/server/zone/` y `pkg/script/lua/`.
   - Crear `scripts/verify_zone_gamelog_100.sh` (`chmod +x`) que ejecute los 8 controles más la verificación de no regresión de Etapa 1 y Etapa 2, imprimiendo `[PASS] 8/8` y devolviendo `exit 0`.

---

## Criterio de Cierre Definitivo de la Etapa 3 y del Clúster Completo `XelionCore`

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 ./...` -> `PASS` en el 100% de los paquetes del repositorio.
3. `./scripts/verify_login_account_100.sh` -> **`exit 0`** (Etapa 1 intacta).
4. `./scripts/verify_world_character_100.sh` -> **`exit 0`** (Etapa 2 intacta).
5. `./scripts/verify_zone_gamelog_100.sh` -> **`8/8` controles en `[PASS]`, `exit 0`**.
6. Marcar Fase 16 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(core): Fase 16 - cierre 100% ZoneService + GameLog (245/245 C2S, S2S, GM AdminLvSet) y auditoria 8/8 PASS"`
