# Fase 05 — `WorldService` Ciclo Completo de Cuenta, Avatares, Apariencia, Opciones, Tutorial y `74` Archivos `.shn` (Cero Hardcode)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 0.3, 1-Bugs #9 y #10, 4.1, 4.2, 5.1, 6 y 9-Control 5A**), [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md) y [`docs/Servicios/Character.md`](../Servicios/Character.md).
- **Dependencias previas:** Fases 01, 02, 02A, 03 y 04 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 05 (`docs/plan/FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md`) leyendo obligatoriamente las Secciones 0.3, 1 (Bugs #9 y #10), 4.1, 4.2, 5.1, 6 y 9 (Control 5A) de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y respetando estrictamente `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados, CERO alias de OpCodes y 100% persistencia vía `CALL` en `xelion_online`.
1) Actualizá `WorldDataStore` (`pkg/server/world/static_data.go`) para cargar y verificar los 74 archivos `.shn` únicos (`54` de `Character.exe` U `71` de `WorldManager.exe`), eliminá el slice hardcodeado `subStrings` de `BadNameFilter` filtrando exclusivamente con `BadNameFilter.shn`, y cargá `data/default_character_data.json`, `data/tutorial_character_data.json` y `data/World/ClassParam/<Class>.json` (27 clases).
2) Implementá al 100% con estructuras PDB exactas todos los handlers de sesión de mundo, avatares, estilo, opciones clásicas/mejoradas, tutorial, asistencia y cambio de nombre en `pkg/server/world/` (`Dept 3, 4, 5, 28 y 45`).
3) Creá `pkg/audit/fase05_world_character_dispatcher_test.go` (Control 5A) y ampliá `pkg/server/world/e2e_world_test.go` contra MySQL `xelion_online`. Al terminar, marcá Fase 05 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 06.
```

---

## Alcance Exacto (Cero Omisiones)

1. **`WorldDataStore` 100% Fiel a Binarios (`pkg/server/world/static_data.go` — `Bug #10`):**
   - Cargar y validar en memoria los **74 archivos `.shn` únicos** (`54` de `Character.exe` listados en la Sección 5.1 $\cup$ `71` de `WorldManager.exe` listados en la Sección 4.1).
   - Eliminar `subStrings: []string{"admin", "gm", ...}` de `BadNameFilter`: la validación de nombres restringidos se rige **únicamente** por las filas reales de `BadNameFilter.shn`.
   - Cargar `data/default_character_data.json` (6 clases base incluyendo `Class 26` en `"Eld"` Lv `60` y `starter_quests`), `data/tutorial_character_data.json` (5 clases de tutorial), `data/World/Map/FieldList.json` + `TutorialMap.json` y `data/World/ClassParam/*.json` (las 27 clases).

2. **Implementar al 100% los Handlers C2S/S2S de Cuenta, Avatar, Opciones y Progresión en `pkg/server/world/`:**
   - **Sesión de Mundo (`Dept 3` y `Dept 45`):**
     * `0x0C0F` `fc_NC_USER_LOGINWORLD_REQ` (`PROTO_NC_USER_LOGINWORLD_REQ`, `320B`: `user[256]` + `validate_new[32]` `64B`): consume token en MySQL (`CALL usp_Session_ConsumeWorldToken`, sin fallbacks en memoria), expulsa sesión duplicada (`KICKOFFFROMWORLD_CMD` + `CONNECTCUT_CMD`) y responde `0x0C10` `PROTO_NC_USER_LOGINWORLD_ACK` con los avatares serializados al byte (`CALL p_Char_GetListOfUserChar` + `CALL p_Item_GetUpgradeEquip`).
     * `0x0C1F` `fc_NC_USER_AVATAR_LIST_REQ` -> recarga y envía la lista actualizada de avatares desde MySQL.
     * `0x0C18` `fc_NC_USER_NORMALLOGOUT_CMD` -> cierra sesión limpia en MySQL (`CALL usp_Session_CloseAndGetPlayMin` + auditoría `AccountLog`).
     * `0x0C33` `fc_NC_USER_WILL_WORLD_SELECT_REQ` -> emite `sOTP` de `32B` (`CALL usp_Session_CreateOTP`) y envía `0x0C34` `PROTO_NC_USER_WILL_WORLD_SELECT_ACK`.
     * `0xB409` `fc_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ` -> cambio de nombre vinculado a cupón/conexión (`CALL p_Char_ChangeID` + registro en `tCharacterChangeID`).
   - **Creación, Borrado y Renombrado de Avatares (`Dept 5` — `PROTO_NC_AVATAR_*`, `Bug #9`):**
     * `0x1401` `fc_NC_AVATAR_CREATE_REQ`: deserializa `PROTO_NC_AVATAR_CREATE_REQ` exacto del PDB; valida contra `BadNameFilter.shn`, `FaceInfo.shn`, `HairInfo.shn`, `HairColorInfo.shn` y `max_shape`; consulta `CALL usp_Tutorial_getIsSkip(userNo)` y el flag de tutorial del paquete para elegir entre `tutorial_character_data.json` + `TutorialMap.json` o `default_character_data.json` + `ClassParam/<Class>.json`; ejecuta `CALL p_Char_Create`, `CALL p_Char_CreateSetDefaultData`, stats base de `ClassParam` en nivel inicial (`1` o `60` para `Sentinel`), ítems iniciales (`CALL p_Item_Create`), skills iniciales (`CALL p_Skill_Create`), quests iniciales (`CALL p_Quest_Start`), shortcuts (`CALL p_Char_SetShortCut`) y opciones iniciales (`CALL p_Init_Char_Option_Improve_Game_Option_Data` + `Key_Map_Data`); responde `0x1403` `PROTO_NC_AVATAR_CREATESUCC_ACK` (o `0x1402` `PROTO_NC_AVATAR_CREATEFAIL_ACK` con el código oficial de error).
     * `0x1407` `fc_NC_AVATAR_ERASE_REQ`: verifica propiedad del personaje, que no sea maestro de un gremio activo (`CALL p_Guild_Member_Get`), ejecuta soft-delete oficial (`CALL p_Char_Delete`) y responde `PROTO_NC_AVATAR_ERASESUCC_ACK` / `PROTO_NC_AVATAR_ERASEFAIL_ACK`.
     * `0x140F` `fc_NC_AVATAR_RENAME_REQ`: valida disponibilidad y filtro `BadNameFilter.shn`, ejecuta `CALL p_Char_ChangeID` y responde con el ACK oficial.
   - **Selección de Avatar, Estilo, Guía, Tutorial y Asistencia (`Dept 4` — `PROTO_NC_CHAR_*`):**
     * `0x1001` `fc_NC_CHAR_LOGIN_REQ`: carga el personaje (`CALL p_Char_GetLoginData`), verifica prisión (`tPrison`), resuelve `zoneNumber` del mapa en `FieldList.json`, obtiene `(sZoneIP, nZonePort)` vía `CALL usp_WorldZone_ResolveMapEndpoint`, registra `CALL p_Char_Login` y `CALL usp_Session_IssueZoneTransfer`, y envía `0x1002` `PROTO_NC_CHAR_LOGIN_ACK`.
     * `0x1078` `fc_NC_CHAR_SET_STYLE_REQ` y `0x107C` `fc_NC_CHAR_SET_STYLE_GET_INFO_REQ` (`CALL p_Char_SetShape` / `p_Char_GetShape` validados contra `FaceInfo.shn`, `HairInfo.shn`, `HairColorInfo.shn`).
     * `0x10DA` `fc_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD` (`tInGameGuide`), `0x10F1` `fc_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ` (`AttendSchedule.shn` + `AttendReward.shn` + `CALL p_Attendance_*`) y `0x1111` `fc_NC_CHAR_TUTORIAL_POPUP_ACK` (`tTutorial`).
   - **Los 21 Handlers Oficiales de Opciones de Cliente y Accesos Rápidos (`Dept 28` — `PROTO_NC_CHAR_OPTION_*`):**
     * Clásicos (`0x7001`, `0x7002`, `0x7004`, `0x7006`, `0x7008`, `0x700A`, `0x700C`, `0x700E`, `0x7010`..`0x7016`) + Mejorados `IMPROVE` (`0x7037`, `0x7039`, `0x703B`, `0x7043`, `0x7045`) + `0x705F` (`LOGIN_BLOCKDATA_ERR`), persistiendo 100% en `tCharacterOptions`, `tCharacterShortCut`, `tCharacterGameOption` y `tCharacterKeyMap` vía `CALL`.

3. **Crear `pkg/audit/fase05_world_character_dispatcher_test.go` (`Control 5A`) y ampliar Pruebas E2E:**
   - Verifica que los 74 archivos `.shn` cargan sin errores, que `BadNameFilter` no tiene strings hardcodeados, y prueba contra MySQL la creación de las 6 clases (incluyendo `Class 26` en `Eld` Lv `60` y modo Tutorial), renombrado, cambio de estilo, opciones `IMPROVE`, borrado lógico y emisión de transferencia en `tActiveZoneTransfer`.

---

## Criterio de Cierre de la Fase 05

1. `go build ./...` -> `OK`.
2. `go test -v ./pkg/server/world/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04|TestFase05"` -> `PASS`.
3. `go test ./...` -> `PASS`.
4. Marcar Fase 05 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(world): Fase 05 - ciclo completo de avatares, estilo, opciones, tutorial y 74 SHNs 100% MySQL"`
