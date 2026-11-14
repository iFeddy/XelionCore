# Fase 04 — Protocolo Oficial en `pkg/protocol`: OpCodes, Estructuras PDB (`WorldManager` + `Character`), `sizeof` Exacto y Cero Alias

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 0.4, 1-Bug #8, 4.2, 4.3, 5.2, 7 y 9-Control 4**), [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md) y [`docs/Servicios/Character.md`](../Servicios/Character.md).
- **Dependencias previas:** Fases 01, 02, 02A y 03 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 04 (`docs/plan/FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md`) leyendo obligatoriamente las Secciones 0.4, 1 (Bug #8), 4.2, 4.3, 5.2, 7 y 9 (Control 4) de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, junto con las definiciones C++ exactas en `docs/Servicios/WorldManager.md` y `docs/Servicios/Character.md`, respetando estrictamente `AGENTS.md`.
Implementá en `pkg/protocol/` (NUNCA crees un subpaquete `pkg/protocol/pdb/`) el 100% de los OpCodes `(department << 10) | command` y las estructuras C++ no-OpTool de `WorldManager.md` y `Character.md` con sus nombres oficiales exactos (`PROTO_NC_*`) y `sizeof` idéntico al byte. Eliminá los 19 alias de emuladores privados detectados en `pkg/protocol/opcodes.go` y actualizá sus llamadores en `pkg/server/world/` y `pkg/server/zone/` para que todo compile y pase tests. Creá `pkg/audit/fase04_world_character_sizeof_test.go` (Control 4 completo parseando `WorldManager.md` y `Character.md`). Al terminar, marcá Fase 04 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 05.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Completar y sanear [`pkg/protocol/opcodes.go`](../../pkg/protocol/opcodes.go) (Reglas de Oro #1 y #3 de `AGENTS.md`):**
   - Definir mediante `protocol = (department << 10) | command` el **100%** de los OpCodes no-`OpTool` de [`WorldManager.md`](../Servicios/WorldManager.md) (`463` entradas) y [`Character.md`](../Servicios/Character.md) (`317` entradas), respetando los departamentos reales del PDB:
     * Dept `02` (`NC_MISC` — `0x0800+`): heartbeat, gametime, pingtest, apex, web key, hack scan, loading bug, spammer, itemshop URL.
     * Dept `03` (`NC_USER` — `0x0C00+`): `PROTO_NC_USER_LOGINWORLD_REQ` (`0x0C0F`), `PROTO_NC_USER_LOGINWORLD_ACK` (`0x0C10`), `PROTO_NC_USER_NORMALLOGOUT_CMD` (`0x0C18`), `PROTO_NC_USER_AVATAR_LIST_REQ` (`0x0C1F`), `PROTO_NC_USER_WILL_WORLD_SELECT_REQ` (`0x0C33`), S2S `WILLLOGIN`, etc.
     * Dept `04` (`NC_CHAR` — `0x1000+`): `PROTO_NC_CHAR_LOGIN_REQ` (`0x1001`), `PROTO_NC_CHAR_LOGIN_ACK` (`0x1002`), `PROTO_NC_CHAR_CHARDATA_REQ` (`0x1004`), `PROTO_NC_CHAR_LOGIN_DB` (`0x1036`), `PROTO_NC_CHAR_LOGOUT_DB` (`0x1037`), `PROTO_NC_CHAR_STAT_INCPOINT_DB_REQ` (`0x105D`), `PROTO_NC_CHAR_SET_STYLE_REQ` (`0x1078`), `PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ` (`0x107C`), `PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD` (`0x10DA`), `PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ` (`0x10F1`), `PROTO_NC_CHAR_TUTORIAL_POPUP_ACK` (`0x1111`), y paquetes `CHARSAVE_*` / `CLIENT_*`.
     * Dept `05` (`NC_AVATAR` — `0x1400+`): `PROTO_NC_AVATAR_CREATE_REQ` (`0x1401`), `PROTO_NC_AVATAR_CREATEFAIL_ACK` (`0x1402`), `PROTO_NC_AVATAR_CREATESUCC_ACK` (`0x1403`), `PROTO_NC_AVATAR_ERASE_REQ` (`0x1407`), `PROTO_NC_AVATAR_ERASEDATA_REQ` (`0x1408`), `PROTO_NC_AVATAR_ERASESUCC_ACK`, `PROTO_NC_AVATAR_ERASEFAIL_ACK`, `PROTO_NC_AVATAR_RENAME_REQ` (`0x140F`), `PROTO_NC_AVATAR_GUILD_MEMBER_DATA_REQ` (`0x1415`).
     * Dept `06` (`NC_MAP` — `0x1800+`), Dept `08` (`NC_ACT` — `0x2000+`: `WHISPER`, `NOTICE`, `PARTYCHAT`), Dept `12` (`NC_ITEM` / `NC_ITEMDB` — `0x3000+`), Dept `14` (`NC_PARTY` — `0x3800+`), Dept `17` (`NC_QUEST` — `0x4400+`), Dept `18` (`NC_SKILL` — `0x4800+`), Dept `21` (`NC_FRIEND` — `0x5400+`), Dept `22` (`NC_KQ` — `0x5800+`), Dept `28` (`NC_CHAR_OPTION` — `0x7000+`), Dept `29` (`NC_GUILD` — `0x7400+`), Dept `31` (`NC_PRISON` — `0x7C00+`), Dept `32` (`NC_REPORT` — `0x8000+`), Dept `37` (`NC_HOLY_PROMISE` — `0x9400+`), Dept `38` (`NC_GUILD_ACADEMY` — `0x9800+`), Dept `41` (`NC_INSTANCE_DUNGEON` — `0xA400+`), Dept `42` (`NC_CHAT_RESTRICT` — `0xA800+`), Dept `44` (`NC_RAID` — `0xB000+`), Dept `45` (`NC_USER_CONNECTION` — `0xB400+`).
   - **Erradicación total de los 19 alias prohibidos (Sección 4.3):** eliminar de `pkg/protocol/` y reemplazar en `pkg/server/world/` y `pkg/server/zone/` todos los nombres no oficiales (`PROTO_NC_MISC_WORLD_CLIENT_KEY_REQ`, `PROTO_NC_CHAR_CLIENT_CHARLIST_ACK`, `PROTO_NC_CHAR_CREATE_REQ`, `PROTO_NC_CHAR_DELETE_REQ`, `PROTO_NC_CHAR_GET_QUICKBAR_REQ`, `PROTO_NC_CHAR_GET_QUICKBAR_STATE_REQ`, `PROTO_NC_CHAR_GET_GAME_SETTINGS_REQ`, `PROTO_NC_CHAR_GET_CLIENT_SETTINGS_REQ`, `PROTO_NC_CHAR_GET_SHORTCUTS_REQ`, `PROTO_NC_CHAR_SAVE_QUICKBAR_REQ`, `PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ`, `PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ`, `PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ`, `PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ`, `PROTO_NC_CHAR_INGAME_REQ`, `PROTO_NC_MISC_TIME_SYNC_ACK`).

2. **Implementar en `pkg/protocol/` las Estructuras Oficiales C++ con `sizeof` Exacto (`#pragma pack(push, 1)`):**
   - Organizar limpiamente en `pkg/protocol/` (`proto_nc_avatar.go`, `proto_nc_char.go`, `proto_nc_char_option.go`, `proto_nc_act_chat.go`, `proto_nc_party.go`, `proto_nc_raid.go`, `proto_nc_friend.go`, `proto_nc_kq.go`, `proto_nc_guild.go`, `proto_nc_guild_academy.go`, `proto_nc_holy_promise.go`, `proto_nc_prison.go`, `proto_nc_itemdb.go`, `proto_nc_quest_skill.go`, `proto_nc_minihouse_pet.go`, `proto_nc_instance_rank.go`, `proto_nc_chat_restrict.go`) cubriendo todas las structs C++ no-`OpTool` de `WorldManager.md` y `Character.md`, junto con sus métodos `MarshalBinary` / `UnmarshalBinary` (o constructores/lectores binarios exactos al byte para payloads con arrays flexibles `[0]` / conteo dinámico).

3. **Crear `pkg/audit/fase04_world_character_sizeof_test.go` (`Control 4` completo):**
   - Parsea automáticamente [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md) y [`docs/Servicios/Character.md`](../Servicios/Character.md):
     * Verifica que el 100% de los OpCodes no-`OpTool` de ambos dispatchers estén definidos en `pkg/protocol` con el valor exacto `(department << 10) | command`.
     * Verifica que el `sizeof` serializado de cada estructura fija registrada en `pkg/protocol` coincida exactamente al byte con el comentario `// size = N bytes` del PDB.
     * Falla automáticamente si reaparece cualquiera de los 19 alias prohibidos en el repositorio o si existe algún directorio `pkg/protocol/pdb`.

---

## Criterio de Cierre de la Fase 04

1. `go build ./...` -> `OK`.
2. `go test -v ./pkg/protocol/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04"` -> `PASS`.
3. `go test ./pkg/server/...` -> `PASS`.
4. Marcar Fase 04 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(protocol): Fase 04 - 100% opcodes y estructuras oficiales de WorldManager y Character PDB sin alias"`
