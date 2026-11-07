# Fase 08 — Saneamiento de Esquema `World00_GameLog`, `27` Stored Procedures Pool-Safe y Worker Asíncrono `GameLogRepository`

- **Nivel recomendado:** `High` *(o `Medium`)*
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0, 1-Bugs #1 y #2, 2.1, 2.2, 2.3, 2.4 y 9-Controles 1, 2 y 3**), [`docs/Servicios/GameLog.md`](../Servicios/GameLog.md), [`docs/TablasLegacy.md`](../TablasLegacy.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Etapas 1 y 2 (`FASE_01` a `FASE_07`) completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 08 (`docs/plan/FASE_08_ESQUEMA_PROCEDURES_Y_REPOSITORIO_GAMELOG.md`) leyendo obligatoriamente las Secciones 0, 1 (Bugs #1 y #2), 2.1, 2.2, 2.3, 2.4 y 9 (Controles 1, 2 y 3) de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/GameLog.md`, `docs/sql/procedures/World00_GameLog_all_procedures.sql` y `AGENTS.md`.
Regla Mandatoria: CERO SQL inline (`SELECT`/`INSERT`/`UPDATE`/`DELETE`), CERO omisiones y 100% Stored Procedures pool-safe (`proc_logic` + `SELECT` final incondicional).
1) Eliminá el archivo obsoleto `database/zone_tables.sql` (86 tablas duplicadas) y verificá las 2 tablas oficiales de `World00_GameLog.bak` (`tGameLog` con sus 15 columnas exactas e índices de auditoría en `database/schema.sql`, y `TEMP_QuestComplete` en `database/schema_character.sql`), cerrando el 100% de las 176/176 tablas oficiales de `docs/TablasLegacy.md` en `xelion_online`.
2) Creá `scripts/build_procedures_gamelog.py` y generá `database/procedures_gamelog.sql` con los 27 Stored Procedures exactos de `World00_GameLog` (23 ODBC de `GameLog.exe` U 26 de `World00_GameLog.bak`, incluyendo `usp_Login_setChannel`), integrándolos en `scripts/apply_mysql_procedures.sh`.
3) Implementá `pkg/database/gamelog.go`, `pkg/database/mysql_gamelog.go` (100% `CALL <SP>`) y el despachador asíncrono no bloqueante `GameLogDispatcher` con canal bufferizado y `Flush`/`Close` limpio en apagado, inyectándolo en `cmd/ZoneService/main.go` y `cmd/WorldService/main.go`.
4) Creá `pkg/audit/fase08_gamelog_db_sp_test.go` cubriendo los Controles 1, 2 y 3 de la Etapa 3 al 100%. Al terminar, marcá Fase 08 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 09.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Retiro de DDL Fantasma y Cierre `176/176` Tablas Oficiales en `xelion_online` (`Secciones 2.1 y 2.2`):**
   - Eliminar definitivamente [`database/zone_tables.sql`](../../database/zone_tables.sql) y cualquier referencia a dicho archivo en scripts o tests.
   - Verificar que `tGameLog` (`nLogNo` bigint PK AUTO_INCREMENT, `dDate`, `nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`, `nBigint1`) posea los índices compuestos `idx_gamelog_date_type (dDate, nType)`, `idx_gamelog_char (nCharNo, dDate)` e `idx_gamelog_itemkey (nItemKey)`, y que junto a `TEMP_QuestComplete` completen las **`176/176` tablas oficiales** de [`docs/TablasLegacy.md`](../TablasLegacy.md).

2. **Transpilación Pool-Safe de los `27` Stored Procedures de `World00_GameLog` (`Sección 2.3`):**
   - Crear `scripts/build_procedures_gamelog.py` y generar `database/procedures_gamelog.sql` compilando en MySQL 8.x el **100% de los `27` SPs**:
     * **Los `22` SPs de inserción compartidos (`GameLog.exe` + `.bak`):** `p_GameLog_Add`, `p_GameLog_Add_Type1`, `p_GameLog_Add_Type2`, `p_GameLog_Add_Type3`, `p_GameLog_Add_Type4`, `p_GameLog_Add_Type5`, `p_GameLog_Add_Type6`, `p_GameLog_Add_Type7`, `p_GameLog_Add_Type8`, `p_GameLog_Add_Type9`, `p_GameLog_Add_TypeA`, `p_GameLog_Add_TypeB`, `p_GameLog_Add_TypeC`, `p_GameLog_Add_TypeD`, `p_GameLog_Add_TypeE`, `p_GameLog_Add_TypeF`, `p_GameLog_Add_TypeG`, `p_GameLog_Add_TypeH`, `p_GameLog_Add_TypeI`, `p_GameLog_Add_TypeJ`, `p_GameLog_Add_TypeK`, `p_GameLog_Add_TypeL` (todos con `@nLogNo BIGINT OUTPUT` asignado desde `LAST_INSERT_ID()` y retornado atómicamente por `SELECT out_nLogNo AS nLogNo;`).
     * **El SP exclusivo de `GameLog.exe` (`0x41B707`):** `usp_Login_setChannel(IN in_nUserNo INT, IN in_nChannelNo INT, IN in_sChannelID VARCHAR(20), OUT out_nResult INT)` con `SELECT out_nResult AS nResult;` final.
     * **Los `4` SPs de consulta de `World00_GameLog.bak`:** `p_GameLog_Get`, `p_GameLog_GetCount`, `p_GameLog_GetLog` y `p_GameLog_GetLogCount` migrados a `PREPARE` / `EXECUTE` / `DEALLOCATE PREPARE` seguros sobre `tGameLog`.
   - Actualizar [`scripts/apply_mysql_procedures.sh`](../../scripts/apply_mysql_procedures.sh) para cargar automáticamente `database/procedures_gamelog.sql` junto con `procedures.sql` y `procedures_character.sql`.

3. **Repositorio Go `GameLogRepository` y Worker Asíncrono `GameLogDispatcher` (`pkg/database/mysql_gamelog.go`):**
   - Definir la interfaz `GameLogRepository` en `pkg/database/gamelog.go` cubriendo los `27` SPs y helpers tipados para los 44 eventos de [`docs/Servicios/GameLog.md`](../Servicios/GameLog.md) (`AddLog`, `AddType1`..`AddTypeL`, `SetLoginChannel`, `GetLogs`, `GetLogCount`), ejecutando exclusivamente `CALL <SP>(...)`.
   - Implementar `GameLogDispatcher` (cola asíncrona con canal bufferizado configurable, worker pool ligero, métricas de encolado/errores y `Close(ctx)` con drenaje completo antes de cerrar `*sql.DB`) e inyectarlo en `ZoneService` y `WorldService`, además de validar los `49` `.shn` de `ChecksumEngine` requeridos por `GameLog.exe`.

4. **Auditoría Automatizada (`pkg/audit/fase08_gamelog_db_sp_test.go` — Controles 1, 2 y 3):**
   - Verifica `176/176` tablas oficiales en `xelion_online`, ausencia física de `database/zone_tables.sql`, existencia de los `27/27` SPs de `World00_GameLog` en `information_schema.ROUTINES`, ejecución concurrente en paralelo sobre `*sql.DB` verificando inserción real en `tGameLog` y `nLogNo > 0`, y escaneo `go/ast` garantizando cero SQL inline en `pkg/database/mysql_gamelog*.go`.

---

## Criterio de Cierre de la Fase 08

1. `bash scripts/apply_mysql_procedures.sh` -> `exit 0`.
2. `go build ./...` -> `OK`.
3. `go test -p 1 -count=1 -v ./pkg/database/... ./pkg/audit -run "TestFase08"` -> `PASS`.
4. `go test -p 1 -count=1 ./...` -> `PASS` (manteniendo intactas Etapas 1 y 2).
5. Marcar Fase 08 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(database,gamelog): Fase 08 - 27 SPs pool-safe World00_GameLog, retiro zone_tables.sql y GameLogDispatcher asincrono"`
