# Fase 06 — Coordinación Global `WorldService`: 100% de los `167` Handlers C2S, Subsistemas Sociales y Topología S2S

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 1-Bugs #11 a #15, 4.1, 4.2, 5.2, 6 y 9-Control 5B**) y [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md).
- **Dependencias previas:** Fases 01 a 05 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 06 (`docs/plan/FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md`) leyendo obligatoriamente las Secciones 1 (Bugs #11 a #15), 4.1, 4.2, 5.2, 6 y 9 (Control 5B) de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y `docs/Servicios/WorldManager.md`, respetando estrictamente `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados (`loadDefaults()` en KQ y `"RouN"` en Friends prohibidos), CERO estado huérfano solo en memoria cuando existe tabla en `World00_Character` y 100% persistencia mediante `CALL <SP>` en `xelion_online`.
Completá el 100% de los 167 handlers C2S del Dispatcher oficial de `WorldManager.pdb` (Sección 4.2.1) y la coordinación S2S en `pkg/server/world/`: Misc/Seguridad/URLs (`Dept 2`), Chat/Susurros/Avisos (`Dept 8`), Party persistida en `tParty`/`tPartyMember` (`Dept 14`), Friends sin `"RouN"` (`Dept 21`), Kingdom Quest leída 100% de `.shn` y persistida en `tKingdomQuest` con votaciones y calendario (`Dept 22`), Guild + War + Tournament + Token + History + Emblem (`Dept 29`), Prison (`Dept 31`), Report (`Dept 32`), Holy Promise / Master-Pupil (`Dept 37`), Guild Academy completa (`Dept 38`), Instance Dungeon Rankings (`Dept 41`), Chat Restrict (`Dept 42`), Raid (`Dept 44`) y Topología viva de Zonas (`tWorldZoneNode`).
Creá `pkg/audit/fase06_world_social_coverage_test.go` (Control 5B verificando 167/167 handlers C2S + S2S). Al terminar, marcá Fase 06 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 07.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Misc, Seguridad, Anti-Spam, WebKey e ItemShop URL (`Dept 2` — 17 handlers C2S):**
   - `0x0804`, `0x0805`, `0x0819` (`HEARTBEAT_REQ/ACK`), `0x080D` (`GAMETIME_REQ`), `0x080F` (`CLIENT_DEBUG_MSG_CMD`), `0x0813` (`PINGTEST_TOOL_WM_CLIENT_ZONE_DB`), `0x081A` / `0x0836` (`APEX_CLIENT_*`), `0x081B` (`WEB_KEY_REQ` -> `CALL usp_UserKey_GetOrIssue`), `0x0827` / `0x0832` (`CS_ACK`), `0x0834` (`HACK_SCAN_STORE_CMD` -> `CALL p_HackScan_AddFull`), `0x0840` (`CLIENT_LOADING_BUG_DETECT_CMD` -> `CALL p_Set_Client_Loading_Bug_Detect`), `0x0846` (`SET_CHAT_BLOCK_SPAM_FILTER_CMD` -> `tSpamFilter`), `0x085F` (`SPAMMER_REPORT_REQ` -> `SpamerReport.shn` + `tSpamBlockTime`), `0x0861` (`SPAMMER_CHAT_BAN_REQ`), `0x086E` (`ITEMSHOP_URL_REQ`).

2. **Chat Global, Restricciones de Chat, Prisión y Reportes (`Dept 8`, `Dept 31`, `Dept 32`, `Dept 42`):**
   - **Chat (`Dept 8`):** `0x200C` (`fc_NC_ACT_WHISPER_REQ` validando bloqueo mutuo y `tSpamBlockTime`), `0x2010` (`fc_NC_ACT_NOTICE_REQ` validando `nAdminLevel`), `0x2014` (`fc_NC_ACT_PARTYCHAT_REQ`), respetando `ChatColor.shn`.
   - **Chat Restrict (`Dept 42`):** `0xA803` (`CHAT_RESTRICT_ADD_REQ`), `0xA807` (`CHAT_RESTRICT_DEL_REQ`), `0xA80B` (`CHAT_RESTRICT_DEL_ALL_REQ`) contra `tChatBlock` (`CALL p_Char_Chat_Restrict*`).
   - **Prison (`Dept 31`):** `0x7C01` (`PRISON_ADD_GM_REQ`), `0x7C04` (`PRISON_END_REQ`), `0x7C06` (`PRISON_GET_REQ`), `0x7C09` (`PRISON_ALTER_GM_REQ`) contra `tPrison` (`CALL p_Prison_*`).
   - **Report (`Dept 32`):** `0x8001` (`REPORT_ADD_REQ`), `0x8003` (`REPORT_CANCEL_REQ`), `0x8005` (`REPORT_GET_REQ`) conectados a `AccountLogRepository` (`CALL p_Report_Add/Del/Get`).

3. **Party (`Dept 14` — 11 handlers) y Raid (`Dept 44` — 13 handlers) (`Bug #12`):**
   - **Party (`0x3802`, `0x3804`, `0x3805`, `0x380A`, `0x3814`, `0x3828`, `0x384B`, `0x3850`, `0x3852`, `0x3854`, `0x385A`):** sincronizar todas las mutaciones en MySQL (`tParty` y `tPartyMember` vía `CALL p_Party_*`), incluyendo `SET_LOOTER_REQ` (`0x385A`) y buscador de grupos (`PARTY_FINDER_*`).
   - **Raid (`0xB001`, `0xB008`, `0xB00C`, `0xB00F`, `0xB014`, `0xB018`, `0xB027`, `0xB02B`, `0xB030`, `0xB035`, `0xB03A`, `0xB03F`, `0xB042`):** creación de banda, reglas de botín, invitaciones, movimiento de slots, líder/sub-líder, maestro despojador, chat de banda y alerta `RAID_WARNING_REQ`.

4. **Friends (`Dept 21` — 5 handlers, `Bug #13`) y Holy Promise / Master-Pupil (`Dept 37` — 4 handlers):**
   - **Friends (`0x5401`, `0x5404`, `0x5405`, `0x541D`, `0x541F`):** eliminar `currentMap = "RouN"` en `friend_manager.go` leyendo `sLoginZone` real de `tCharacter` (`CALL p_Friend_GetAll2`), e implementar `0x541F` (`PROTO_NC_FRIEND_FIND_FRIENDS_REQ`) y puntos de amistad (`usp_Character_GetFriendpoint`/`SetFriendpoint`).
   - **Holy Promise / Master-Pupil (`0x9401`, `0x9405`, `0x9406`, `0x940A`):** implementar vínculo maestro-discípulo y promesa sagrada leyendo `HolyPromiseReward.shn` + `PupMain.shn` y persistiendo en `tHolyPromise` / `tMasterPupils` vía `CALL p_HolyPromise_*` y `CALL p_MasPup_*`.

5. **Kingdom Quests (`Dept 22` — 11 handlers, `Bug #11`) e Instance Dungeon Rankings (`Dept 41` — 4 handlers):**
   - **Kingdom Quests (`0x5801`, `0x5803`, `0x5805`, `0x5807`, `0x5809`, `0x581B`, `0x5827`, `0x582A`, `0x5831`, `0x5834`, `0x5837`):** eliminar por completo `loadDefaults()` en `kq_manager.go`, leer columnas exactas de `KingdomQuest.shn`, `KingdomQuestMap.shn`, `KQTeam.shn`, `KQIsVote.shn` y `KQVoteMajorityRate.shn`, persistir estado en `tKingdomQuest` (`CALL p_KingdomQuest_*`) e implementar calendario (`0x5809` `KQ_SCHEDULE_REQ`), refresco (`0x581B`) y sistema de votación de expulsión (`0x5827`, `0x582A`, `0x5834`).
   - **Instance Rankings (`0xA41B`..`0xA41E`):** pestañas `DAMAGEDEALT`, `DAMAGETAKEN`, `HEALINGDONE` y `MYRANK` consultando `InstanceDungeon.json` y estadísticas de instancia.

6. **Gremios, Guerras, Torneos, Emblemas (`Dept 29` — 30 entradas C2S) y Academia (`Dept 38` — 24 handlers C2S):**
   - **Guild (`Dept 29`):** además de los 10 handlers básicos existentes, implementar guerras (`0x7422` `GUILD_WAR_REQ`, `0x7424` `GUILD_WAR_ABLE_LIST_REQ`, `0x7426` `GUILD_WAR_LIST_REQ`), consulta de nombre (`0x7476`), torneos (`0x747A`, `0x747E`, `0x74C3`, `0x74E9` usando `GuildTournament*.shn` + `data/World/Map/GuildTournament.json` + `CALL usp_GuildTournament_*`), tokens (`0x74BE` -> `CALL p_Guild_GetToken`), historial (`0x74FD` -> `CALL usp_Guild_getHistoryList`), emblemas (`0x7501`, `0x7505` -> `CALL p_GuildEmblem_*`) y parámetros de nivel/grado desde `GuildGradeData.shn`, `GuildGradeScoreData.shn`, `GuildLevelScoreData.shn`.
   - **Guild Academy (`Dept 38`):** completar los **24 handlers** (`0x9807`..`0x988F`: info, rank info, list, member list/join/vanish/leave, master teleport, chat ban + cancel, notify, guild invite/join, get/set/clear reward money & item, graduate list, ranking list, set master, chat, history list) usando `GuildAcademy*.shn` y los SPs reales de `tGuildAcademy*`.

7. **Topología Dinámica de Zonas (`pkg/server/world/zone_manager.go` + `cmd/WorldService/main.go`, `Bug #15`):**
   - Resolver el endpoint `(sZoneIP, nZonePort)` de cualquier mapa a través de `FieldList.json` (`zoneNumber`) + `CALL usp_WorldZone_ResolveMapEndpoint`, y gestionar traspasos inter-zona (`PROTO_NC_CHAR_MAPLOGIN_REQ`) vía `tActiveZoneTransfer`.

8. **Crear `pkg/audit/fase06_world_social_coverage_test.go` (`Control 5B` completo):**
   - Parsea la tabla **Sección 4.2.1 (`167` handlers C2S de `WorldManager.md`)** y verifica automáticamente que **167/167 OpCodes C2S** están registrados en el Dispatcher de `WorldService`, además de ejecutar pruebas funcionales contra MySQL para Chat, Party, Raid, Friends, KQ, Guild, GuildTournament, GuildAcademy, HolyPromise, Prison, Report y ChatRestrict.

---

## Criterio de Cierre de la Fase 06

1. `go build ./...` -> `OK`.
2. `go test -v ./pkg/server/world/... ./pkg/audit -run "TestFase01|TestFase02|TestFase03|TestFase04|TestFase05|TestFase06"` -> `PASS`.
3. `go test ./...` -> `PASS`.
4. Marcar Fase 06 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(world): Fase 06 - 167 handlers C2S de WorldManager, subsistemas sociales y topologia 100% MySQL"`
