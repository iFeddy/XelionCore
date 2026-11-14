# Fase 03 — Stored Procedures Nativos de `AccountLog`, `StatisticsData` y `LoginSession`

- **Nivel recomendado:** `Medium`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 0, 3.3, 3.4 y 3.5 completas**)
- **Fuentes SQL originales obligatorias:** [`AccountLog_all_procedures.sql`](../sql/procedures/AccountLog_all_procedures.sql), [`StatisticsData_all_procedures.sql`](../sql/procedures/StatisticsData_all_procedures.sql) y [`docs/Servicios/AccountLog.md`](../Servicios/AccountLog.md)
- **Dependencias previas:** Fases 01 y 02 completadas.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 03 (`docs/plan/FASE_03_PROCEDURES_ACCOUNTLOG_STATS_SESSION.md`) leyendo las Secciones 0, 3.3, 3.4 y 3.5 de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md` y traduciendo fielmente `AccountLog_all_procedures.sql` y `StatisticsData_all_procedures.sql` a MySQL sobre `xelion_online`.
Creá `database/procedures_accountlog.sql`, `database/procedures_statistics.sql` y `database/procedures_login_session.sql`, aplicálos en `xelion_online` y completá el Control 2 al 100%. Al pasar todos los tests, actualizá `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Crear `database/procedures_accountlog.sql`:**
   - **Los 12 SPs ODBC de `AccountLog.pdb` + los 17 SPs restantes no-OpTool de `AccountLog.bak` (Sección 3.3 — 29 SPs únicos, excluyendo solo los 3 `p_Report_*_OPT`):**
     * `p_AccountLog_Add` *(inserta en `tAccountLog`; cuando `nType = 1` actualiza **únicamente** `tUser.sUserIP` —**prohibido tocar `tUser.dDate`**— e inserta primer ingreso diario en `tUserLog`; cuando `nType = 2` guarda `nPlayMin`; cuando `nType = 3` guarda `sUserPassword` y `nErrorCode`)*.
     * `p_IPBlock_IsBlock` *(compara numéricamente IPv4 con `INET_ATON` contra `tIPBlock`; si coincide, permite pasar si la cuenta ya posee al menos un registro `nType = 1` en `tAccountLog`, y bloquea con `1` si es cuenta sin login exitoso previo)*.
     * `usp_Log_setLogout`, `p_Report_Add`, `p_Report_Del`, `p_Report_Get`, `p_SameTimeChar_Add`, `p_SameTimeCharClass_Add`, `p_Event_L20_Get`, `p_Event_L20_Set`, `usp_AccountLog_GetDateDiffLastlogout`, `usp_EventiPad_setTicket`.
     * `p_AccountLog_Get`, `p_AccountLog_GetCount`, `p_AccountLog_GetLog`, `p_AccountLog_GetLogCount`.
     * `p_IPBlock_Add`, `p_IPBlock_Delete`, `p_IPBlock_Get`, `p_IPBlock_GetList`.
     * `p_Event_Friend_DelRequest`, `p_Event_Friend_EventComplete`, `p_Event_Friend_Get`, `p_Event_Friend_GetCount`, `p_Event_Friend_Giveup`, `p_Event_Friend_SetHideMsg`.
     * `p_Event_WorldCup_Enter`, `usp_EventiPad_getList`, `usp_EventiPad_getTicket`.
2. **Crear `database/procedures_statistics.sql` (los 18 SPs de la Sección 3.4):**
   - `pFindCCU`, `p_DayClassLevelPlay_Add`, `p_DayClassLevelPlay_Query`, `p_DayClassLevel_Add`, `p_DayClassLevel_Query`, `p_DayConnectionTime`, `p_DayUser`, `p_DayUser_Add`, `p_DayUser_Query`, `p_SameTimeChar_Add` *(compartido con AccountLog)*, `p_SameTimeChar_Query`, `p_SameTimeWorld`, `p_UserAgeReport` *(sobre `tAccountDemographics`)*, `p_UserAge_Query`, `p_UserGenderTest` *(sobre `tAccountDemographics`)*, `p_UserGender_Query`, `p_UserLocation_Query`, `p_getQuestTimes`.
3. **Crear `database/procedures_login_session.sql` (los 12 SPs operativos de la Sección 3.5 para cero SQL inline en Go):**
   - `usp_Version_GetClientKey`, `usp_Version_SyncClientKey`, `usp_WorldServer_GetList`, `usp_WorldServer_GetByNo`, `usp_WorldServer_UpdateCCU`, `usp_Session_RegisterLogin`, `usp_Session_IssueWorldToken` (`BINARY(64)`), `usp_Session_ConsumeWorldToken` (`BINARY(64)`), `usp_Session_CreateOTP` (`CHAR(32)`), `usp_Session_ConsumeOTP` (`CHAR(32)`), `usp_Session_CloseAndGetPlayMin`, `p_HackScan_Add` y `p_Syslog_Add`.
4. **Aplicar los 3 archivos SQL en `xelion_online`** y crear un script `scripts/apply_mysql_procedures.sh` para recargar todos los SPs de forma repetible.
5. **Crear `pkg/audit/fase03_accountlog_sp_test.go` (Control 2 completo de la Sección 9.2):**
   - Verifica en `information_schema.ROUTINES` el **100%** de los SPs de `Account`, `AccountLog`, `StatisticsData` y `LoginSession`.
   - Prueba funcionalmente `p_AccountLog_Add` (comprobando que `tUser.dDate` queda intacto al segundo), `p_IPBlock_IsBlock` (cuenta nueva bloqueada vs cuenta veterana permitida), `p_Report_Add`/`Get`/`Del` (cooldowns 1h y 24h), `usp_Session_IssueWorldToken`/`ConsumeWorldToken` (64 bytes exactos) y `usp_Session_CreateOTP`/`ConsumeOTP` (32 bytes single-use).

---

## Criterio de Cierre de la Fase 03

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "TestFase01|TestFase02|TestFase03"` -> `PASS`.
3. Actualizar `docs/plan/README.md` y ejecutar `git commit -m "feat(db): Fase 03 - 100% stored procedures de AccountLog, StatisticsData y LoginSession"`.
