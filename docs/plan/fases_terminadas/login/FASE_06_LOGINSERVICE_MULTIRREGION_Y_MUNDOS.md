# Fase 06 — `LoginService` Unificado: Todas las Regiones, `FOUS`, `tWorldServer` y OTP en MySQL

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 0, 1, 4.1, 5.1, 6.1, 8-Fase 4 y 9.5**)
- **Dependencias previas:** Fases 01 a 05 completadas.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 06 (`docs/plan/FASE_06_LOGINSERVICE_MULTIRREGION_Y_MUNDOS.md`) leyendo las Secciones 0, 1, 4.1, 5.1, 6.1, 8 (Fase 4) y 9.5 de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md`.
Completá el 100% de `LoginService` (`Login` + `Account` + `AccountLog` unificados contra `CALL` en MySQL `xelion_online`), activando todas las variantes regionales y de plataforma, lectura dinámica de `FOUS` y `tWorldServer`, OTP real de 32 bytes y token de mundo de 64 bytes, sin ningún estado hardcodeado. Al pasar los tests y el Control 5 de LoginService, actualizá `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Inicialización 100% desde MySQL en [`cmd/LoginService/main.go`](../../cmd/LoginService/main.go) y [`pkg/server/login/world_manager.go`](../../pkg/server/login/world_manager.go):**
   - Al arrancar, si `ClientVersionKeyInfo.txt` está disponible y `FOUS` estuviera vacía, sincronizar con `CALL usp_Version_SyncClientKey`; en cada handshake leer `GameVersionKey` desde MySQL vía `CALL usp_Version_GetClientKey()`.
   - Reemplazar la lista estática de mundos en memoria por lectura dinámica desde `tWorldServer` (`CALL usp_WorldServer_GetList()` y `CALL usp_WorldServer_GetByNo()`), calculando `WorldStatus` en función de `bIsMaintenance`, `nCurrentCCU`, `nMaxUser` y `nUserRatable` (`LOGIN_USER_RATABLE`).
   - **Bug #8:** eliminar el goroutine de `cmd/LoginService/main.go` que insertaba `AddSameTimeChar(w.ID, 0)` con `0` hardcodeado cada 60 segundos.
2. **Completar el 100% de los handlers del Dispatcher en [`pkg/server/login/handler.go`](../../pkg/server/login/handler.go) y [`pkg/server/login/server.go`](../../pkg/server/login/server.go):**
   - `0x0C65` `fc_NC_USER_CLIENT_VERSION_CHECK_REQ`: valida los `64B` contra `FOUS.GameVersionKey`. Si coincide -> envía `0x0C67` (`PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK`) + `0x0807` (`PROTO_NC_MISC_SEED_ACK`). Si falla -> envía **`0x0C02`** (`PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK`) y cierra conexión.
   - `0x0C04` `fc_NC_USER_XTRAP_REQ` -> responde `0x0C05` (`PROTO_NC_USER_XTRAP_ACK`).
   - **Pipeline unificado oficial para TODAS las variantes de Login:**
     * C2S directos:
       - `0x0C5A` `fc_NC_USER_US_LOGIN_REQ` (`316B`: `sUserID`/`sEmail`/`r_` vía `CALL usp_User_USLoginGame`)
       - `0x0C06` `fc_NC_USER_LOGIN_REQ` (`272B`: vía `CALL usp_User_loginGame`)
       - `0x0C37` `fc_NC_USER_LOGIN_WITH_OTP_REQ` (`32B`: vía `CALL usp_Session_ConsumeOTP`, **corrigiendo Bug #3**)
       - `0x0C1D` `fc_NC_USER_LOGIN_NETMARBLE_REQ` (vía `CALL usp_User_LoginNetmarble`)
       - `0x0C20` `fc_NC_USER_LOGIN_OUTSPARK_REQ` (`85B`: vía `CALL p_UserKey_Get_FromAccount` + `CALL usp_User_chkBlock`)
     * Handlers regionales y de soporte S2S/internos completos:
       - `GER`: `0x0C39`..`0x0C3C`
       - `TW`: `0x0C3E`..`0x0C41`
       - `JP`: `0x0C4B`..`0x0C4E`
       - `CH`: `0x0C53`..`0x0C54`
       - `Teenager`: `0x0C26`..`0x0C28` (`CALL p_UserPlayRemainTime_Get_FromAccount` / `Set_FromAccount` + envío de `PROTO_NC_USER_TEENAGER_CMD` si aplica)
     * **Orden exacto de negocio (corrigiendo Bugs #2, #4, #11, #12):**
       1. Autenticar credenciales primero para obtener `nUserNo` y `nAuthID`. Si falla -> `CALL p_AccountLog_Add` (`nType = 3`, guardando `sUserPassword` y `nErrorCode`) + `CALL p_Syslog_Add` (`0xC802`) + enviar `0x0C09` (`PROTO_NC_USER_LOGINFAIL_ACK`).
       2. Ejecutar `CALL p_IPBlock_IsBlock(nUserNo, clientIP)`. Si retorna bloqueado -> auditar fallo `nError = 73` y enviar `0x0C09`.
       3. Evaluar `CALL usp_Pcbang_chkIP(clientIP)` y `CALL p_UserPlayRemainTime_Get_FromAccount(nUserNo)`.
       4. Registrar sesión activa con `CALL usp_Session_RegisterLogin` y auditoría exitosa con `CALL p_AccountLog_Add` (`nType = 1`, **sin modificar `tUser.dDate`**) + `CALL p_Syslog_Add` (`0xC801`).
       5. **No** auto-insertar `SetReturnUserItem` durante el login (eso ocurre en `0x0C60`/`0x0C64`).
       6. Consultar `CALL usp_Tutorial_getIsSkip(nUserNo)` y si aplica enviar `PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD` (`0x0C5F`).
       7. Enviar `PROTO_NC_USER_LOGIN_ACK` (`0x0C0A`) con la lista de mundos leída en vivo desde `tWorldServer`.
   - `0x0C1B` `fc_NC_USER_WORLD_STATUS_REQ` -> responde `0x0C1C` leyendo `tWorldServer` desde MySQL.
   - `0x0C0B` `fc_NC_USER_WORLDSELECT_REQ` -> valida mantenimiento y `nMaxUser * nUserRatable / 100` en `tWorldServer`, emite `validate[32]` (`64 bytes` criptográficos) en `tActiveAccountSession` (`CALL usp_Session_IssueWorldToken`) y envía `0x0C0C` (`83 bytes`).
   - `0x0C18` `fc_NC_USER_NORMALLOGOUT_CMD` y desconexión de socket -> `CALL usp_Session_CloseAndGetPlayMin` + `CALL p_AccountLog_Add` (`nType = 2`, `nPlayMin`) + `CALL p_Syslog_Add` (`0xC803`).
   - Handlers S2S / puente de `Login.pdb`: `0x0801..0x0805`, `0x0808` (`fc_NC_MISC_USER_COUNT_CMD` -> `CALL usp_WorldServer_UpdateCCU`), `0x080B` (`DELIVER_WM_LOGIN_ACDB_CMD`), `0x080C` (`DELIVER_WM_LOGIN_ALDB_CMD`), `0x0823` (`START_THE_BOOM_CMD`) y `0x0C35` (`CREATE_OTP_REQ` -> `0x0C36`).
3. **Crear `pkg/audit/fase06_dispatcher_coverage_test.go` (Control 5 de `LoginService`) y ampliar [`pkg/server/login/e2e_login_test.go`](../../pkg/server/login/e2e_login_test.go):**
   - Verifica que los 37 puntos de protocolo no-`OpTool` de `Login.pdb` están registrados y probados contra `xelion_online`.

---

## Criterio de Cierre de la Fase 06

1. `go build ./...` -> OK.
2. `go test -v ./pkg/server/login/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04|TestFase05|TestFase06"` -> `PASS`.
3. Actualizar `docs/plan/README.md` y ejecutar `git commit -m "feat(login): Fase 06 - LoginService 100% fiel a PDB con multirregion, OTP y mundos desde MySQL"`.
