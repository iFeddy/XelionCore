# Fase 07 — Integración `WorldService` + `ZoneService`, 49 SHNs y Auditoría Forense Final 100%

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 1, 5.1, 5.2, 6.1, 8-Fase 5, 8-Fase 6 y 9 completa**)
- **Dependencias previas:** Fases 01 a 06 completadas.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 07 final (`docs/plan/FASE_07_INTEGRACION_WORLD_ZONE_Y_AUDITORIA_FINAL.md`) leyendo las Secciones 1, 5.1, 5.2, 6.1, 8 (Fases 5 y 6) y la Sección 9 completa de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md`.
Completá la integración transversal con `WorldService` y `ZoneService`, conectá los 49 SHNs de `ChecksumEngine` con `tHackScanLog`, creá `pkg/audit/login_account_completeness_test.go` y `scripts/verify_login_account_100.sh`, y no te detengas hasta que los 7 controles de la Sección 9 pasen al 100% (`exit 0`). Al terminar, actualizá `docs/plan/README.md` y hacé el `git commit` final en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Integración en [`pkg/server/world/`](../../pkg/server/world/handler.go) y [`cmd/WorldService/main.go`](../../cmd/WorldService/main.go):**
   - **Bug #14:** en `PROTO_NC_USER_WORLD_LOGIN_REQ` (`0x0C0F`), leer y validar los **`64 bytes` íntegros** (`uint16 validate_new[32]`) invocando `CALL usp_Session_ConsumeWorldToken` en MySQL (eliminando el truncamiento a 32B y el fallback legacy `18+32`), y si existía una sesión previa activa del mismo `nUserNo`, expulsarla limpiamente (`KICKOFFFROMWORLD_CMD` / `CONNECTCUT_CMD`).
   - **Bug #3:** en `PROTO_NC_USER_WILL_WORLD_SELECT_REQ` (`0x0C33`), invocar `CALL usp_Session_CreateOTP` (`0x0C35` -> `0x0C36`) generando un `Name8 sOTP` criptográfico real de `32 bytes` y enviarlo en `PROTO_NC_USER_WILL_WORLD_SELECT_ACK` (`0x0C34`).
   - Conectar los flujos `Account` / `AccountLog` consumidos desde `WorldService`:
     * `PROTO_NC_USER_RETURN_CHECK_REQ` (`0x0C60`) -> `PROTO_NC_USER_RETURN_CHECK_ACK` (`0x0C61`) y `PROTO_NC_USER_SET_RETURN_CMD` (`0x0C64` -> `CALL usp_ReturnUser_SetUserItem`).
     * `PROTO_NC_USER_LOGOUT_LAST_DAY_REQ` (`0x0C62`) -> `PROTO_NC_USER_LOGOUT_LAST_DAY_ACK` (`0x0C63` -> `CALL usp_AccountLog_GetDateDiffLastlogout`).
     * `PROTO_NC_MISC_WEB_DB_KEY_REQ` (`0x081D` -> `0x081E`), `PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ` (`0x082F`), `PROTO_NC_MISC_EVENT_HIT3_ADD_CASH` (`0x0820`), `PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC` (`0x0838`), `PROTO_NC_MISC_EVENT_L20_DB_REQ` (`0x082C` -> `0x082D`), `PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD` (`0x0843` -> `CALL usp_EventiPad_setTicket`).
   - Actualizar periódicamente `CALL usp_WorldServer_UpdateCCU` y `CALL p_SameTimeChar_Add` con el CCU real de `WorldService` + `CALL p_Syslog_Add` (`0xC811`).
2. **Integración en [`pkg/server/zone/`](../../pkg/server/zone/handler.go) y [`cmd/ZoneService/main.go`](../../cmd/ZoneService/main.go):**
   - **Checksums de los 49 `.shn` (`shndatafilechecksum.cpp`, Sección 5.2):** al fallar `VerifyDataFiles` en `ChecksumEngine`, persistir el incidente en `tHackScanLog` mediante `CALL p_HackScan_Add` (`PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD` `0x0835`).
   - **Bug #12 (Logout completo):** al desconectar un jugador, pasar el `worldNo` real (no `0` hardcodeado) a `CALL usp_Log_setLogout`, y cuando finaliza la sesión de la cuenta ejecutar `CALL usp_Session_CloseAndGetPlayMin` + `CALL p_AccountLog_Add` (`nType = 2`, `nPlayMin` real) + `CALL p_Syslog_Add` (`0xC803`).
   - **Bug #8 (Telemetría por clase y nivel):** corregir el cron de `cmd/ZoneService/main.go` para que `CALL p_SameTimeCharClass_Add` agrupe por `(worldNo, Job, Level)` reales de los jugadores online en vez de pasar `level = 1` fijo.
   - Conectar retiro de `tChargeItem` (`CALL usp_Charge_ItemSelectList` / `CALL usp_Charge_ItemDraw`) y consulta/uso de `tCash` (`CALL usp_Cash_getCash` / `CALL usp_Cash_useCash` con `status = 0`).
3. **Suite Final de Auditoría 100% (`pkg/audit/login_account_completeness_test.go` + `scripts/verify_login_account_100.sh`):**
   - Unificar y ejecutar los **7 Controles Obligatorios de la Sección 9** (`9.1` Tablas/Vistas/Columnas, `9.2` 100% Routines SQL, `9.3` AST Cero SQL Inline/JSON, `9.4` OpCodes y `sizeof` contra los 3 `.md` de PDBs, `9.5` Cobertura total de Dispatchers en `Login`/`World`/`Zone`, `9.6` Checksums reales de los 49 `.shn` + `tHackScanLog`, y `9.7` E2E forense completo Login -> World -> OTP -> Re-Login -> ReturnUser -> Cash -> Report -> Logout).

---

## Criterio de Cierre Definitivo del Plan

1. `go build ./...` -> OK.
2. `go test -v ./...` -> `PASS` en todos los paquetes del repositorio.
3. `./scripts/verify_login_account_100.sh` -> imprime el reporte detallado de los **7 controles en `PASS`** y finaliza con `exit 0`.
4. Marcar Fase 7 como `Completado` en `docs/plan/README.md` y ejecutar `git commit -m "feat(core): Fase 07 - integracion completa Login/Account/AccountLog con World/Zone y auditoria 100% PASS"`.
