# Fase 05 — Refactorización de `pkg/database`: Cero SQL Inline, Cero JSON y 100% `CALL <SP>`

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 0, 1, 3.5, 8-Fase 2 y 9.3**)
- **Dependencias previas:** Fases 01, 02, 03 y 04 completadas.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 05 (`docs/plan/FASE_05_REPOSITORIOS_GO_100_CALL_MYSQL.md`) leyendo las Secciones 0, 1, 3.5, 8 (Fase 2) y 9.3 de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md`.
Eliminá `MemoryAccountRepository` y `pkg/server/login/data/accounts.json`, y reemplazá el 100% del SQL inline en `pkg/database` (incluyendo los 20 métodos actuales de `mysql_account.go` y `transfer.go`) por invocaciones `CALL` a los Stored Procedures de `xelion_online`. Adaptá los llamadores actuales para que todo el workspace compile y pase `go test ./...` y el Control 3 (AST). Al terminar, actualizá `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Configurar DSN MySQL seguro para Stored Procedures ([`pkg/config/config.go`](../../pkg/config/config.go) / inicialización DB):**
   - Garantizar siempre `parseTime=true&multiStatements=true&charset=utf8mb4` en el DSN de conexión a `xelion_online`.
2. **Eliminar todo residuo en memoria / JSON:**
   - Borrar [`pkg/server/login/data/accounts.json`](../../pkg/server/login/data/accounts.json) (y el directorio `pkg/server/login/data` si queda vacío).
   - Eliminar `MemoryAccountRepository` y `NewMemoryAccountRepository` de [`pkg/database/database.go`](../../pkg/database/database.go), definiendo interfaces ricas y tipadas para `AccountRepository` / `AccountLogRepository` / `LoginSessionRepository` (cero type-assertions `repo.(*database.MySQLAccountRepository)` en los servidores).
   - Actualizar `README.md` si menciona `MemoryAccountRepository`.
3. **Reemplazar los 20 métodos actuales de [`pkg/database/mysql_account.go`](../../pkg/database/mysql_account.go) + agregar todos los métodos faltantes usando exclusivamente `CALL <SP>(...)`:**
   - Eliminar `initAuxiliaryTables()` y `UpdateLoginSession()`.
   - Dividir limpiamente en:
     * `pkg/database/mysql_account.go`: `LoginGameUS` (`CALL usp_User_USLoginGame`), `LoginGame` (`CALL usp_User_loginGame`), `LoginNetmarble` (`CALL usp_User_LoginNetmarble`), `CheckBlock` (`CALL usp_User_chkBlock`), `GetAccountByUsername` (`CALL usp_User_getUser`), `CreateAccount` (`CALL usp_User_insert` + `usp_User_insertOption`), `UpdatePassword` (`CALL usp_User_updatePassword`), `DeleteAccount` (`CALL usp_User_delete`), `SetDrawComment` (`CALL usp_User_setDrawComment`), `SetChannelLog` (`CALL usp_User_setChannelLog`), `CheckPcbangIP` (`CALL usp_Pcbang_chkIP`), `GetWebKey` (`CALL p_UserKey_Get_FromAccount`), `GetPlayRemainTime` / `SetPlayRemainTime`, `GetUserDateDiffCreated`, `GetReturnUserDateDiffLastLog`, `SetReturnUserItem`, `GetTutorialIsSkip` (`CALL usp_Tutorial_getIsSkip`), `GetTutorialCharInfo`, `SetTutorialCharInfo`.
     * `pkg/database/mysql_account_shop.go`: `GetCash`, `AddCash` (`usp_Cash_setCash` / `usp_Event_addCash`), `UseCash` (`usp_Cash_useCash`), `AddCashEventL20`, `SetSerialKey`, `GetChargedItems`, `InsertChargedItem`, `DrawChargedItem`, `DeleteChargedItem`, y métodos de catálogo/órdenes/pagos (`usp_Item_*`, `usp_Pay_*`).
     * `pkg/database/mysql_accountlog.go`: `IsIPBlocked(userNo, ip)` (`CALL p_IPBlock_IsBlock`), `LogAccount` (`CALL p_AccountLog_Add`), `LogLogout` (`CALL usp_Log_setLogout`), `GetDateDiffLastLogout` (`CALL usp_AccountLog_GetDateDiffLastlogout`), `AddReport`, `GetReports`, `DeleteReport`, `AddSameTimeChar`, `AddSameTimeCharClass`, `EventL20Get` / `EventL20Set`, `EventiPadSetTicket`, `AddHackScanLog` (`CALL p_HackScan_Add`), `AddSyslog` (`CALL p_Syslog_Add`).
     * `pkg/database/mysql_world_session.go` + actualización de [`pkg/database/transfer.go`](../../pkg/database/transfer.go): persistir y consumir tokens de mundo (`64 bytes` `uint16[32]`), `sOTP` (`32 bytes`), sesiones activas (`tActiveAccountSession`), versión (`FOUS`) y mundos (`tWorldServer`) vía `CALL usp_Version_*`, `CALL usp_WorldServer_*` y `CALL usp_Session_*` (manteniendo `CreateZoneTransfer`/`ConsumeZoneTransfer` para saltos de mapa intra-mundo `World <-> Zone`).
4. **Adaptar llamadores existentes** (`pkg/server/login/`, `pkg/server/world/`, `pkg/server/zone/`, `cmd/` y `*_test.go`) a las nuevas firmas (como `IsIPBlocked(userNo, ip)` y `LogAccount` con `playMin` y `userPW`) para que todo el repositorio compile y pase tests.
5. **Crear `pkg/audit/fase05_no_inline_sql_ast_test.go` (Control 3 completo de la Sección 9.3):**
   - Escanea el AST de Go verificando que **ningún** repositorio de `Account` / `AccountLog` / `WorldSession` contenga `SELECT`, `INSERT`, `UPDATE`, `DELETE` o `CREATE TABLE` inline (100% `CALL`), ni existan `MemoryAccountRepository`, `UpdateLoginSession`, `initAuxiliaryTables` o `accounts.json`.

---

## Criterio de Cierre de la Fase 05

1. `go build ./...` -> OK.
2. `go test -v ./pkg/database/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04|TestFase05"` -> `PASS`.
3. `go test ./pkg/server/...` -> `PASS`.
4. Actualizar `docs/plan/README.md` y ejecutar `git commit -m "refactor(database): Fase 05 - migracion 100% a CALL stored procedures y eliminacion de JSON/memoria"`.
