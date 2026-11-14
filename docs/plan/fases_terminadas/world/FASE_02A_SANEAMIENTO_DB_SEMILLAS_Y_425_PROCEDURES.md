# Fase 02A — Saneamiento Integral de Esquema (`114/114`), Semillas (`data/`) y los `425` Stored Procedures Pool-Safe

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 0, 1, 2, 3.1, 3.2 y 3.3 completas**), [`docs/TablasLegacy.md`](../TablasLegacy.md), [`docs/sql/procedures/World00_Character_all_procedures.sql`](../sql/procedures/World00_Character_all_procedures.sql), [`data/DefaultCharacterData.txt`](../../data/DefaultCharacterData.txt) y [`data/TutorialCharacterData.txt`](../../data/TutorialCharacterData.txt).
- **Dependencias previas:** Fases 01 y 02 finalizadas en `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 02A (`docs/plan/FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md`) leyendo las Secciones 0, 1 (Bugs #1 a #5), 2, 3.1, 3.2 y 3.3 de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y respetando estrictamente `AGENTS.md`.
Saneá al 100% lo entregado en Fase 01 y Fase 02:
1) Completá las 18 tablas restantes del `.bak` para llegar a las 114/114 tablas oficiales de `World00_Character.bak` en `xelion_online` más las 2 tablas operativas S2S (`tActiveZoneTransfer` y `tWorldZoneNode`), y poblá la tabla auxiliar `num`.
2) Eliminá las 11 filas falsas `*_Template` / `Tutorial_*` insertadas bajo `nUserNo = 1` en `tCharacter`, `tCharacterShape` y `tItem` en `database/seed_character_defaults.sql` (y en MySQL en vivo), y regenerá al 100% `data/default_character_data.json` (incluyendo `starter_quests` de `#DEFINE QUEST` y las 6 clases base `1, 6, 11, 16, 21, 26`) y `data/tutorial_character_data.json` (decodificando `data/TutorialCharacterData.txt` en UTF-16LE con `#DEFINE TUTORIAL_CHARACTER`, `DEFAULT_ITEM` y `TUTORIAL_ITEM`).
3) Completá en `database/procedures_character.sql` los 425 Stored Procedures exactos (`311` ODBC de `Character.exe` U `405` no-DBA de `World00_Character.bak` = `419` SPs oficiales + los `6` SPs operativos S2S de la Sección 3.3), corrigiendo las referencias cross-DB (`Account..tUser` -> `tUser`) y encapsulando todo SP con parámetros `OUT` en `proc_logic: BEGIN ... END proc_logic;` seguido incondicionalmente de `SELECT out_...` para que jamás devuelva `sql.ErrNoRows` en el pool `*sql.DB` de Go.
Actualizá los tests de auditoría Control 1 y Control 2 en `pkg/audit/`, verificá que todo pase al 100%, marcá Fase 02A como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 03.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Completar el 100% del Esquema `World00_Character.bak` (`114/114` tablas) + `2` Tablas Operativas S2S (`database/schema_character.sql`):**
   - Agregar las **18 tablas oficiales restantes** de [`docs/TablasLegacy.md`](../TablasLegacy.md) listadas en la **Sección 2** (`num`, `tItem_deleted`, `auditCharacterData`, `TEMP_ItemCnt`, `TEMP_QuestComplete`, `Temp_Guar`, `Temp_Holy`, `Waits_Sysprocesses`, `chr_a`, `chr_b`, `chr_bak`, `tBugReport_20091204`, `tCharacterState_Backup`, `tCharacterState_isSaveZero_BackUp`, `tGuildAcademy_Backup`, `tMimic_100610`, `tQuest_BackUp`, `temp_GA`) con sus columnas y tipos exactos.
   - Crear las **2 tablas operativas MySQL** de la **Sección 2.1**:
     * `tActiveZoneTransfer` (reemplazo MySQL de un solo uso para `TransferManager.zoneTransfers` entre procesos `WorldService` y `ZoneService`).
     * `tWorldZoneNode` (registro vivo de IPs, puertos y jugadores por `(nWorldNo, nZoneNo)` para eliminar endpoints hardcodeados).
   - Poblar idempotentemente la tabla `num` (`0..1000`) en `database/seed_character_defaults.sql`.

2. **Saneamiento de Semillas y Archivos Oficiales de Inicio (`Bug Crítico #1` y `Bug Crítico #2`):**
   - **Eliminar de [`database/seed_character_defaults.sql`](../../database/seed_character_defaults.sql) y de la base `xelion_online` en vivo** los 11 registros falsos `Fighter_Template`..`Sentinel_Template` y `Tutorial_Fighter`..`Tutorial_Trickster` insertados en `tCharacter`, `tCharacterShape` y `tItem` bajo `nUserNo = 1` (conservando `tExp` con los 150 niveles, `tGameOptionInit` y `tKeyMapInit`).
   - Actualizar [`scripts/convert_character_data.py`](../../scripts/convert_character_data.py) y regenerar:
     * [`data/default_character_data.json`](../../data/default_character_data.json) desde [`data/DefaultCharacterData.txt`](../../data/DefaultCharacterData.txt) incluyendo los 6 bloques `CHARACTER` (`1, 6, 11, 16, 21, 26`, donde `26` inicia en `"Eld"` `12259, 13434`, Lv `60`, Exp `24346145`), `ITEM`, `SKILL`, **`QUEST` (`starter_quests`)**, `MAXSHAPE` y `SHORTCUT`.
     * [`data/tutorial_character_data.json`](../../data/tutorial_character_data.json) decodificando [`data/TutorialCharacterData.txt`](../../data/TutorialCharacterData.txt) (`UTF-16LE`) con sus bloques `TUTORIAL_CHARACTER`, `DEFAULT_ITEM` y `TUTORIAL_ITEM` (cero clases vacías).

3. **Completar y Blindar los `425` Stored Procedures en `database/procedures_character.sql` (`Bug Crítico #4` y `Bug Crítico #5`):**
   - Incluir el **100% de la unión exacta** documentada en la **Sección 3**:
     * Los **311 SPs ODBC de `Character.exe`** (`297` compartidos con `.bak` usando matching case-insensitive + los **14 SPs exclusivos del binario `Character.exe`**: `p_Guild_SetMoney`, `p_IsPromotionUser`, `p_Item_MiniHouse_Effect_Resource_Del`, `p_MiniHouseURL_Get`, `p_MiniHouseURL_Set`, `p_Promotion_Get`, `p_Promotion_Set`, `usp_CharacterLoadingDetailInfo`, `usp_GuildTournament_setRefundNotice`, `usp_Guild_getHistoryList`, `usp_Guild_setHistory`, `usp_LinkScroll_SetDestMap`, `usp_Pet_setDelete`, `usp_Pet_setTendency`).
     * Los **108 SPs adicionales no-DBA de `World00_Character.bak`** (`p_Guild_Make`, `p_Guild_GetList`, `p_Guild_SetGrade`, `p_Item_GetList`, `p_Item_Merge`, `p_Item_Split_N_Merge`, `p_Item_Upgrade`, `p_KingdomQuest_GetAll`, `p_KingdomQuest_Set`, `p_MasPup_GetAll`, `p_MasPup_SetMaster`, `p_Quest_Add`, `p_Quest_GetAll`, `p_Quest_GetAllDone`, `p_Prison_End`, `p_HolyPromise_CheckRing`, etc.), excluyendo únicamente `Adam_ListProc`, `CommandExecute` e `IndexOptimize`.
     * Los **6 SPs operativos S2S** de la **Sección 3.3** (`usp_Session_IssueZoneTransfer`, `usp_Session_ConsumeZoneTransfer`, `usp_WorldZone_UpsertNode`, `usp_WorldZone_Heartbeat`, `usp_WorldZone_GetList`, `usp_WorldZone_ResolveMapEndpoint`).
   - **Corrección Pool-Safe Obligatoria (`proc_logic`):** En los procedimientos con parámetros `OUT`, ningún `LEAVE` temprano por validación o error puede saltar por encima del `SELECT out_...` final. Encapsular el cuerpo en `proc_logic: BEGIN ... END proc_logic;` y emitir `SELECT out_...` siempre al final.
   - **Resolución Cross-DB:** Reemplazar toda referencia legacy `Account..tUser` / `Account.dbo.tUser` por `tUser` en `xelion_online` (y en `p_Char_Create`, si `tUser.nAuthID = 9`, asignar `nAdminLevel = 100` sin consultar una base inexistente `Options..Auto_GM`).
   - Actualizar [`scripts/apply_mysql_procedures.sh`](../../scripts/apply_mysql_procedures.sh) y `scripts/apply_mysql_schema_character.sh` para aplicar todo de forma idempotente.

4. **Actualizar Suite de Auditoría DB/Data (`Control 1` y `Control 2` al 100%):**
   - Actualizar [`pkg/audit/fase01_world_character_schema_test.go`](../../pkg/audit/fase01_world_character_schema_test.go) para verificar las **114 tablas del `.bak` + 2 tablas operativas**, prohibir filas `*_Template` en `tCharacter`, y verificar que `default_character_data.json` (6 clases con `starter_quests`) y `tutorial_character_data.json` (5 clases de tutorial) carguen completos.
   - Actualizar/crear `pkg/audit/fase02_world_character_sp_test.go` validando la existencia exacta de los **425 SPs** en `information_schema.ROUTINES` y ejecutando pruebas funcionales en `*sql.DB` tanto en ramas de éxito como de fallo (ej. crear personaje en slot ocupado o con nombre duplicado en `p_Char_Create` debe retornar código `<= 0` vía `Scan`, jamás `sql.ErrNoRows`, más ciclo `usp_Session_IssueZoneTransfer` / `ConsumeZoneTransfer`).

---

## Criterio de Cierre de la Fase 02A

1. `bash scripts/apply_mysql_schema_character.sh && bash scripts/apply_mysql_procedures.sh` -> `OK`.
2. `go build ./...` -> `OK`.
3. `go test -v ./pkg/audit -run "TestFase01|TestFase02"` -> `PASS` (incluyendo los controles de Etapa 1 y los 425 SPs + 116 tablas de Etapa 2).
4. `go test ./...` -> `PASS`.
5. Marcar Fase 02A como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "fix(database,data): Fase 02A - 114 tablas bak, 425 SPs pool-safe y semillas Default/Tutorial sin templates falsos"`
