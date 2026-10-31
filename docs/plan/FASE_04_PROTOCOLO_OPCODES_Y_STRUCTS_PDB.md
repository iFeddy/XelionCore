# Fase 04 — Protocolo Oficial: OpCodes, las 292 Estructuras C++ de los 3 PDBs y `sizeof` Exacto

- **Nivel recomendado:** `Medium`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 1, 4.2, 5.3, 6.2, 7 y 9.4**)
- **Fuentes PDB obligatorias:** [`docs/Servicios/Login.md`](../Servicios/Login.md), [`docs/Servicios/Account.md`](../Servicios/Account.md) y [`docs/Servicios/AccountLog.md`](../Servicios/AccountLog.md)
- **Dependencias previas:** Fases 01, 02 y 03 completadas.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 04 (`docs/plan/FASE_04_PROTOCOLO_OPCODES_Y_STRUCTS_PDB.md`) leyendo las Secciones 1, 4.2, 5.3, 6.2, 7 y 9.4 de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md` y las definiciones C++ exactas en `docs/Servicios/Login.md`, `docs/Servicios/Account.md` y `docs/Servicios/AccountLog.md`.
Implementá todos los OpCodes y estructuras no-OpTool con sus nombres oficiales exactos y tamaños en bytes idénticos al PDB, actualizá las referencias de nombres viejos en el código existente para que todo compile limpio, y dejá pasando el Control 4 al 100%. Al terminar, actualizá `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Corregir y completar [`pkg/protocol/opcodes.go`](../../pkg/protocol/opcodes.go):**
   - **Bug Crítico #1:** cambiar `PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK` de `(3 << 10) | 102` (`0x0C66`) a **`(3 << 10) | 2` (`0x0C02`)**.
   - Definir el **100%** de los OpCodes no-`OpTool` de `Login.md`, `Account.md` y `AccountLog.md` usando estrictamente `protocol = (department << 10) | command` (incluyendo todas las regiones `US`, `GER`, `TW`, `JP`, `CH`, `Netmarble`, `Outspark`, `Teenager`, `ReturnUser`, `LogoutLastDay`, `Tutorial` `0x0C5F`, `Log` `0x0720..0x0722`, `Misc` `0x0801..0x0843`, `Report` `0x8007..0x800C` y `Syslog` `0xC801..0xC811`).
   - **Bug #13 (Regla de Oro #1):** eliminar los nombres no oficiales `PROTO_NC_USER_FILE_CHECK_REQ/ACK`, `PROTO_NC_USER_WORLDSTATUS_REQ/ACK` y `PROTO_NC_MISC_SEED_CMD`, reemplazándolos en `pkg/protocol` y en sus llamadores actuales (`pkg/server/login/handler.go` y `pkg/server/login/e2e_login_test.go`) por los nombres oficiales exactos:
     * `PROTO_NC_USER_XTRAP_REQ` (`0x0C04`), `PROTO_NC_USER_XTRAP_ACK` (`0x0C05`)
     * `PROTO_NC_USER_WORLD_STATUS_REQ` (`0x0C1B`), `PROTO_NC_USER_WORLD_STATUS_ACK` (`0x0C1C`)
     * `PROTO_NC_MISC_SEED_ACK` (`0x0807`)
2. **Implementar y verificar al byte todas las estructuras C++ no-`OpTool` en `pkg/protocol/`:**
   - [`pkg/protocol/proto_nc_user.go`](../../pkg/protocol/proto_nc_user.go): todas las structs `PROTO_NC_USER_*` y `PROTO_NC_USP_*` de `Login.md`, `Account.md` y `AccountLog.md` (incluyendo `PROTO_NC_USER_LOGIN_WITH_OTP_REQ` `32B`, `PROTO_NC_USER_CREATE_OTP_REQ` `266B`, `PROTO_NC_USER_CREATE_OTP_ACK` `36B`, `PROTO_NC_USER_LOGIN_OUTSPARK_REQ` `85B`, `PROTO_NC_USER_LOGIN_NETMARBLE_REQ`, `PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ` `278B`, `PROTO_NC_USER_RETURN_CHECK_REQ` `6B`, `PROTO_NC_USER_RETURN_CHECK_ACK` `11B`, `PROTO_NC_USER_SET_RETURN_CMD` `4B`, `PROTO_NC_USER_LOGOUT_LAST_DAY_REQ` `6B`, `PROTO_NC_USER_LOGOUT_LAST_DAY_ACK` `10B`, `PROTO_NC_USER_TEENAGER_REQ/ACK/SET_CMD/CMD`, `PROTO_NC_USER_KICKOFFFROMWORLD_CMD` `4B`, `PROTO_NC_USER_CONNECTCUT_CMD` `2B`, y todas las structs `US`, `GER`, `TW`, `JP`, `CH`).
   - `pkg/protocol/proto_nc_log.go`: `PROTO_NC_LOG_USER_LOGINFAIL` (`278B`), `PROTO_NC_LOG_USER_LOGIN` (`9B`), `PROTO_NC_LOG_USER_LOGOUT` (`9B`).
   - `pkg/protocol/proto_nc_misc.go`: `PROTO_NC_MISC_SEED_ACK` (`2B`), `PROTO_NC_MISC_S2SCONNECTION_*`, `PROTO_NC_MISC_USER_COUNT_CMD` (`4B`), `PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD` (`2B`), `PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD` (`2B`), `PROTO_NC_MISC_WEB_DB_KEY_REQ` (`6B`), `PROTO_NC_MISC_WEB_DB_KEY_ACK` (`40B`), `PROTO_NC_MISC_EVENT_HIT3_ADD_CASH` (`28B`), `PROTO_NC_MISC_CONNECTFROMWHERE_*`, `PROTO_NC_MISC_EVENT_DONE_MUNSANG_*`, `PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS` (`21B`), `PROTO_NC_MISC_EVENT_L20_DB_REQ` (`272B`), `PROTO_NC_MISC_EVENT_L20_DB_ACK` (`13B`), `PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD` (`51B`), `PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD` (`6B`), `PROTO_NC_MISC_CONNECTER_REQ` (`1B`), `PROTO_NC_MISC_CONNECTER_ACK` (`131B`), `PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ/ACK` (`8B`).
   - `pkg/protocol/proto_nc_syslog.go`: `PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS` (`40B`), `PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE` (`296B`), `PROTO_NC_SYSLOG_ACCOUNT_LOGOUT` (`8B`), `PROTO_NC_SYSLOG_SERVER_CCU` (`5B`).
   - [`pkg/protocol/proto_nc_report.go`](../../pkg/protocol/proto_nc_report.go): verificar exactitud de `2369B`, `28B`, `6B`, `8B`, `7B`, `137B`.
   - [`pkg/protocol/proto_nc_char.go`](../../pkg/protocol/proto_nc_char.go): actualizar `New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK(otp [32]byte)` para recibir el `Name8 sOTP` real de 32 bytes en vez de rellenar ceros (manteniendo compatibilidad con su llamador en `pkg/server/world/handler.go`).
3. **Crear `pkg/audit/fase04_protocol_sizeof_test.go` (Control 4 completo de la Sección 9.4):**
   - Parsea automáticamente `docs/Servicios/Login.md`, `docs/Servicios/Account.md` y `docs/Servicios/AccountLog.md`, comprueba cada OpCode y cada `sizeof` serializado, y falla si aparece cualquier alias prohibido.

---

## Criterio de Cierre de la Fase 04

1. `go build ./...` -> OK.
2. `go test -v ./pkg/protocol/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04"` -> `PASS`.
3. `go test ./pkg/server/...` -> `PASS`.
4. Actualizar `docs/plan/README.md` y ejecutar `git commit -m "feat(protocol): Fase 04 - 100% estructuras y opcodes oficiales de Login, Account y AccountLog PDB"`.
