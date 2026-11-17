# FASE 27: Port 1:1 de Kingdom Quest — `CKQServer` (`WorldManager.exe`), `KingdomQuest::*` (`Zone.exe`), `CPFsCharacter` KQ (`Character.exe`), Lua `data/LuaScript/KQ/` y verificación por secuencia real de paquetes del cliente

> **Nivel de Modelo Recomendado:** **`High`**
> **Alcance:** 5 sub-fases (`27A` → `27E`), una sesión de chat por sub-fase (o `/goal` secuencial). Cada sub-fase compila, pasa tests y hace commit propio.
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/WorldManager/WorldManager.exe` + `.pdb`, `../AstraServer/Server/Zone00/Zone.exe` + `.pdb`, `../AstraServer/Server/Character/Character.exe`, `data/KingdomQuest.shn`, `data/KingdomQuestMap.shn`, `data/KQTeam.shn`, `data/KQIsVote.shn`, `data/KQVoteMajorityRate.shn`, `data/KingdomQuestRew.shn`, `data/KQItem.shn`, `data/LuaScript/KQ/`, `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md`, `docs/Servicios/Character.md`.
> **Herramientas:** `python3 scripts/disasm_pdb_handler.py {world|zone|character} <símbolo> [bytes]`, `scripts/parse_pdb_types.py`, `scripts/extract_pdb_enums.py`.

---

## 0. Por qué existe esta fase (diagnóstico del 2026-11-16)

Las Fases 15, 21 y 26 marcaron KQ como "Completado" con suites `8/8 PASS`, pero esas suites validaban **tamaños de struct, opcodes y existencia de handlers**, no el comportamiento. `pkg/server/world/kq_manager.go` es una implementación **inventada**:

| # | Comportamiento actual (inventado) | Comportamiento oficial (evidencia) |
| :-: | :--- | :--- |
| 1 | `LIST_REFRESH_REQ` (`0x581B`) responde `LIST_TIME_ACK` + `LIST_ACK` (`0x5802`) con **las 57 KQ en un paquete de 8089 bytes** (ver `logs/WorldService_20261116_1614.log`). | `CKQServer::Ack_NC_KQ_LIST_REFRESH` @ `0x004559E0`: `LIST_TIME_ACK` (`0x581C`, `PROTO_NC_DATA_SEND1`, 0x2B bytes) **solo la 1ª vez por sesión** (flag `session+0x2590`); snapshot `CKQ` de KQ con `Status <= 4`; diff contra el `CKQ` cacheado en la sesión (`session+0x2594`) → `LIST_DELETE_ACK` (`0x581E`), `LIST_UPDATE_ACK` (`0x581F`, si cambió `Status` o `NumOfJoiner`) y `LIST_ADD_ACK` (`0x581D`) **en chunks**: flush cuando `n*141+2 >= 0x1D26`; al final copia el snapshot a la sesión (`CKQ::operator=`). |
| 2 | `SCHEDULE_REQ` responde `LIST_ACK` (`0x5802`). | `fc_NC_KQ_SCHEDULE_REQ` @ `0x00413E00` responde `SCHEDULE_ACK` (`0x580A`). |
| 3 | `LIST_REQ` devuelve la lista completa. | `fc_NC_KQ_LIST_REQ` @ `0x00413980` responde `0x5802` con su propia lógica (desensamblar). |
| 4 | Las 57 filas se cargan como `Status=1` con `StartTime = boot + 10 + i*5 min`, defaults inventados (`MaxPlayers=15`, `LimitTime=3600`, `StartWaitTime=300`, `VoteAgreeRate=70`). | La agenda sale de `KINGDOM_QUEST` (180 bytes): `ST_Year..ST_Second`, `StartWaitTime`, `NextStartMode`, `NextStartDeleyMin`, `RepeatMode`, `RepeatCount` → `CKQServer::AddNewScheduleList` @ `0x004555D0`, `GetNextScheduleTime` @ `0x00452160`, `DoSchedule` @ `0x00455E40`, `DelOldShceduleList` @ `0x00454080`. Estados reales hasta `0x0B`. |
| 5 | Nunca se emite `JOINING_ALARM_CMD` (`0x5824`) ni `JOINING_ALARM_LIST` (`0x5826`). | `SetJoining` @ `0x00454170` → `Send_NC_KQ_JOINING_ALARM_CMD` @ `0x00452930` (broadcast `GetFirstCharSession/GetNextCharSession`); `Send_NC_KQ_JOINING_ALARM_LIST` @ `0x00452B60` al loguear. |
| 6 | No existe armado de mapa ni inicio. `wms_NC_KQ_W2Z_*` en Zone solo están **listados por nombre**. | `DoSetMakeRoom` @ `0x00455980` → `AllocMapLink` → `Send_NC_KQ_W2Z_MAKE_REQ` (377 bytes `PROTO_KQ_INFO`); `DoSetStart` @ `0x00455F90` → `KQTeam_DivideRandom`, `KQTeam_LeaveParty`, `Send_NC_KQ_W2Z_START_CMD`, `Send_NC_KQ_NOTIFY_CMD_ToJoiner`, `Send_NC_KQ_JOINING_ALARM_END_CMD`. Zone: `wms_NC_KQ_W2Z_MAKE_REQ` @ `0x00497790`, `wms_NC_KQ_W2Z_DESTROY_CMD` @ `0x004972C0`, `KingdomQuestEntrance::operator()` @ `0x00496150`. |
| 7 | Al vencer, la KQ pasa a `Status=0` para siempre. | Ciclo `DoKingdomQuest` @ `0x00456190` (llama a `DoSchedule`, `DoSetMakeRoom`, `DoSetStart`, `VoteProcessing`, `DoSetNoMapBF`, `AutoRefreshList`, etc.). |
| 8 | Errores de `JOIN_ACK` inventados `1/2/3/4`; nombre de participante = `client.Username` (cuenta, no personaje). | `CKQServer::PlayerJoin` @ `0x004542F0`, `CKQ::IsJoinable` @ `0x0045CED0`, `IsPlayerAlreadyJoin` @ `0x004524D0`. |
| 9 | En arranque se llama `p_KingdomQuest_Set(id, title)` por fila. | Ningún `.exe` del clúster llama `p_KingdomQuest_*` (strings vacías en los 12 binarios); es del OPTool. |
| 10 | Zone: `handleKQEntryResponseAck` hace `if kqID == 0 { kqID = 1 }`. | `ShinePlayer::sp_NC_KQ_ENTRYRESPONCE_ACK` @ `0x00496990`. |
| 11 | `KingdomQuest.shn` tiene columnas repetidas sin nombre (`MapLink[4]`, `DemandGender[2]`; `KingdomQuestMap.Map[10]`/`Clear[10]`; `KingdomQuestRew.Reward[15]`/`RewardRate[15]`). `Row.GetX(name)` solo devuelve la primera. | Usar acceso posicional `Row.GetXAt(idx)` para los arrays. |

---

## 1. Estructuras oficiales de referencia (extraídas del PDB)

| Struct | `sizeof` | Uso |
| :--- | :-: | :--- |
| `KINGDOM_QUEST` | 180 | Fila de `KingdomQuest.shn` (incluye `ST_*`, `NextStartMode`, `NextStartDeleyMin`, `RepeatMode`, `RepeatCount`, `RewardIndex`, `MapLink[4]`, `ScriptLanguage[32]`, `InitValue[32]`, `UseClass`, `DemandGender[2]`). |
| `KINGDOM_QUEST_MAP` | 363 | `NumOfMap`, `BaseMap[32]`, `Map[10][32]`, `Clear[10]`. |
| `KQTeam` | 22 | `MaxMemberGap`, `IsTeamPVP`, `KQTeamDivideType` (`KQTD_NONE/RANDOM/USERSELECT`), regen XY rojo/azul. |
| `PROTO_KQ_INFO_CLIENT` | 141 | Lo que ve el cliente (ya correcto en `pkg/protocol/proto_nc_kq.go`). |
| `PROTO_KQ_INFO` | 377 | `PROTO_KQ_INFO_CLIENT` + `NextStartMode`, `NextStartDelayMin`, `RepeatMode`, `RepeatCount`, `RewardIndex`, `DemandMobKill`, `ScheduleTime`+`tm`, `RunCounter`, `PROTO_KQ_MAP_INFO MapLink[4]`, `ScriptLanguage[32]`, `ScriptInitValue[32]`, `IsTeamPVP`, `SHINE_XY_TYPE TeamRegenXY[2]`. Es el payload de `W2Z_MAKE_REQ`/`W2Z_START_CMD`. |
| `PROTO_KQ_MAP_INFO` | 26 | `MapIndex`, `Name3 MapBase`, `Name3 MapName`, `MapClear`. |
| `KQ_MAP_ALLOC_TABLE` | 40 | `AllocatedQuest[10]`. |
| `CKQServer::KQ_ONE_SCHEDULE_LIST` | 76 | `NumOfList`, `tm TMList[2]`. |
| `CKQServer::KQ_JOINER_BF` | 44 | `nCharNo`, `WMHandle`, `KQ_JOIN_CHAR_INFO CharInfo`, `bInVote`, `bBan`, `nVotingCount`. |
| `CKQServer::KQ_VOTE_INFO` | 20 | `nTeamType`, `nStarterIndex`, `nTargetIndex`, `tVoteEndTime`, `nYesNum`, `nNoNum`, `nCancelNum`. |
| `CKQ` / `CKQServer` | 113108 / 1441284 | Contenedores (array de `PROTO_KQ_INFO` de stride `0x179`). |
| Enums | — | `KINGDOMQUEST_ID`, `KINGDOMQUEST_REW_ID`, `KQ_TEAM_TYPE`, `KQTeamDivideType`, `KQ_VOTING_TYPE`. |

---

## 2. Sub-fases

### 27A — Datos y modelo (`pkg/data/kq` nuevo) — **Completado**
1. Loader tipado de los 7 SHN KQ con acceso **posicional** para columnas repetidas; tests que comparen fila a fila contra el dump binario (57 filas `KingdomQuest`, 38 `KingdomQuestMap`, 8 `KQTeam`, 30 `KQIsVote`, 2 `KQVoteMajorityRate`, 64 `KingdomQuestRew`, 2 `KQItem`).
2. Tipos Go `KINGDOM_QUEST`, `KINGDOM_QUEST_MAP`, `KQTeam`, `PROTO_KQ_INFO`, `PROTO_KQ_MAP_INFO` con serializadores y tests de `sizeof` (377 / 26).
3. Extraer el **enum de estados** de KQ leyendo cada `mov byte ptr [reg+4], N` en `CKQServer::*` y `CKQ::*` (documentarlo en este archivo con dirección de cada transición).

#### Enum oficial de estados de `PROTO_KQ_INFO::Status` (`+0x04` / `CKQ+0x08`, extraído de `WorldManager.exe`)

| Valor | Constante Go (`pkg/data/kq`) | Escrituras (`mov`) en `WorldManager.exe` | Lecturas / transiciones (`cmp`) en `WorldManager.exe` |
| :---: | :--- | :--- | :--- |
| `0x00` | `KQStatusScheduled` | `CKQServer::AddNewScheduleList` @ `0x004556DE` (`mov byte ptr [ebp-0x188], 0x0`) | `CKQServer::DoSetMakeRoom` @ `0x00455998` (`cmp byte ptr [esi], 0x0`, busca `ScheduleTime <= now`) |
| `0x01` | `KQStatusMakingRoom` | `CKQServer::DoSetMakeRoom` @ `0x004559BA` (`mov byte ptr [esi], 0x1`, tras `AllocMapLink` + `Send_NC_KQ_W2Z_MAKE_REQ`) | `CKQServer::DoSetNoMapBF` @ `0x004522A0` (`cmp byte ptr [esi], 0x1`), `CKQServer::SetJoining` @ `0x004541EE` (`cmp al, 0x1`) |
| `0x02` | `KQStatusJoining` | `CKQServer::SetJoining` @ `0x004542A4` (`mov byte ptr [ebx+0x8], 0x2`, llamada por `fc_NC_KQ_Z2W_MAKE_ACK` @ `0x0042AA0E` cuando `err == 0x981`) | `CKQ::IsJoinable` @ `0x0045CF20`, `DoSetStart` @ `0x00456046`, `Send_NC_KQ_JOINING_ALARM_CMD` @ `0x0045296F`, `Send_NC_KQ_JOINING_ALARM_END_CMD` @ `0x00452AAB`, `Send_NC_KQ_JOINING_ALARM_LIST` @ `0x00452BDB`, `Recv_NC_KQ_TEAM_SELECT_REQ` @ `0x00453822` |
| `0x03` | `KQStatusCountdown` | `CKQServer::DoSetStart` @ `0x0045608B` (`NumOfJoiner >= MaxPlayers`) y @ `0x00456142` (`StartTime + StartWaitTime <= now` y `NumOfJoiner >= MinPlayers` y `KQTeam_CanKQStart == 1`); inicializa `RunCounter = 0x0A` @ `0x0045608F` / `0x00456146` | `CKQServer::DoSetStart` @ `0x00455FD3` (decrementa `RunCounter`, emite `NOTIFY_CMD` en `10` y `JOINING_ALARM_END_CMD`) |
| `0x04` | `KQStatusRunning` | `CKQServer::DoSetStart` @ `0x00455FE1` (`mov byte ptr [ebx-0x75], 0x4` cuando `Status == 3 && RunCounter == 0`, seguido de `KQTeam_DivideRandom`, `KQTeam_LeaveParty`, `Send_NC_KQ_W2Z_START_CMD`) | `InKQStatusRunning` @ `0x00452D98`, `LinkInKQ` @ `0x00454E88`, `LinkOutKQ` @ `0x00454E29`, `Recv_NC_KQ_VOTE_START_REQ` @ `0x00456415`, `Recv_NC_KQ_VOTE_START_CHECK_REQ` @ `0x004533C2`; cota visible `Status <= 4` en `fc_NC_KQ_LIST_REQ` (@ `0x004139E2`, `0x00413A89`) y `Ack_NC_KQ_LIST_REFRESH` (@ `0x00455A90`) |
| `0x05` | `KQStatusDone` | `CKQServer::SetDone` @ `0x00454892` (`mov byte ptr [ebx+0x4], 0x5`, llamada por `fc_NC_KQ_Z2W_END_CMD` @ `0x0042AB33`) | `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` (`sub al, 5; cmp al, 5` -> rango terminado `0x05..0x0A`) |
| `0x06` | `KQStatusDoneSkip` | `CKQServer::SetDoneSkip` @ `0x004548F8` (`mov byte ptr [esi+0x4], 0x6`, llamada por `DoSetStart` @ `0x00456124` reason=`2` y `KQTeam_CanKQStart` @ `0x004553FC` reason=`2` / `0x00455443` reason=`3`) | `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` (rango `0x05..0x0A`) |
| `0x07` | `KQStatusReserved7` | *(Ninguna instrucción `mov` escribe `0x07` en `WorldManager.exe`)* | Cubierto por el rango `0x05..0x0A` en `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` |
| `0x08` | `KQStatusNoMap` | `CKQServer::DoSetMakeRoom` @ `0x004559C4` (`AllocMapLink == 0`); `CKQServer::SetNoMapBF` @ `0x004524AE` (llamada por `fc_NC_KQ_Z2W_MAKE_ACK` @ `0x0042AA22` cuando `err != 0x981`) | `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` (rango `0x05..0x0A`) |
| `0x09` | `KQStatusMakeTimeout` | `CKQServer::DoSetNoMapBF` @ `0x00452300` (`mov byte ptr [eax], 0x9` cuando `Status == 1` y `ScheduleTime + 60s <= now` sin `Z2W_MAKE_ACK`) | `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` (rango `0x05..0x0A`) |
| `0x0A` | `KQStatusNotFound` | `CParserClient::fc_NC_KQ_STATUS_REQ` @ `0x00413C6A` (`mov byte ptr [ebp-0x200b], 0xa` en respuesta `0x5804` cuando `CKQ::Get(handle) == NULL`) | `CKQServer::DelOldShceduleList` @ `0x004540B2-0x004540B6` (rango `0x05..0x0A`) |
| `0x0B` | `KQStatusDeleting` | `CKQServer::DelOldShceduleList` @ `0x004540FD` (schedule viejo con status en `5..10` superado por uno más reciente del mismo `ID`); `CKQServer::SetJoining` @ `0x00454227` (status previo distinto de `1`/`2`); `fc_NC_OPTOOL_KQ_DELETE_CMD` @ `0x00425DF5` | `CKQServer::DelOldShceduleList` @ `0x00454126` (`cmp byte ptr [esi+0x4], 0xb` -> `FreeMapLink`, `FreeJoiner`, `Del`) |

### 27B — `CKQServer` en `WorldService` (reemplaza `kq_manager.go`)
Portar 1:1, función por función, con la dirección del exe en el comentario de cada método:
`Read`, `ReadMap`, `ReadTeam`, `ReadIsVote`, `ReadVoteRate`, `InitMapTable`, `AddNewScheduleList`, `GetNextScheduleTime`, `Add`/`Del`, `DoSchedule`, `DelOldShceduleList`, `SetJoining`, `PlayerJoin`, `PlayerDisjoin`, `IsPlayerAlreadyJoin`, `GetJoinerIndex`, `JoinerInfoUpdateByLogin`, `CheckCharBannedInLogin`, `DoSetMakeRoom`, `AllocMapLink`, `GetEmptyMapLink`, `FreeMapLink`, `DoSetNoMapBF`/`SetNoMapBF`, `DoSetStart`, `KQTeam_CanKQStart`, `KQTeam_DivideRandom`, `KQTeam_LeaveParty`, `GetJoinerTeamType`, `SetDone`, `SetDoneSkip`, `FreeJoiner`, `LinkInKQ`, `LinkOutKQ`, `LinkEndProcessing`, `AutoRefreshList`, `Ack_NC_KQ_LIST_REFRESH`, `DoKingdomQuest`, votación completa (`InitAllVoteInfo`, `ClearVoteInfo`, `InVoting`, `VoteProcessing`, `GetVoingEndTime`, `Recv_NC_KQ_VOTE_*`, `Send_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD`, `Send_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD`), `SendBroadcast_NC_KQ_JOIN_LIST_ACK`, `send_NC_KQ_PLAYER_DISJOIN_CMD`, todos los `Send_NC_KQ_*`.
- Handlers `CParserClient::fc_NC_KQ_*` (`LIST`, `STATUS`, `JOIN`, `JOIN_CANCEL`, `SCHEDULE`, `LIST_REFRESH`, `VOTE_*`, `JOIN_LIST`, `TEAM_SELECT`) respetando `CheckConnectionValidation` y los campos de sesión que escriben (p. ej. `session+0x1DF68`/`0x1DF6C` en `LIST_REFRESH`).
- Handlers `CParserZone::fc_NC_KQ_Z2W_MAKE_ACK` @ `0x0042A9C0` y `fc_NC_KQ_Z2W_END_CMD` @ `0x0042AAC0`.
- Estado por sesión: flag de `TIME_ACK` enviado y `CKQ` cacheado para el diff.
- Timer: el período con que `WorldManagerServer` llama a `DoKingdomQuest` (desensamblar el caller) en lugar del tick genérico de `RunWorldTick`.
- Eliminar `p_KingdomQuest_Set` en el arranque.

### 27C — `KingdomQuest::*` en `ZoneService`
`wms_NC_KQ_W2Z_MAKE_REQ` → `KingdomQuestContainer` (`kqc_MapUseCheck`, `kqc_SearchHandle`), creación de los mapas de `MapLink[4]` / `KingdomQuestMap`, `Z2W_MAKE_ACK`; `wms_NC_KQ_W2Z_START_CMD` → `KingdomQuestEntrance` (`kqe_newentrance`, `kqe_entry`, `kqe_allentry`, `kqe_sort`, `kqe_search`), `KQ_ENTRYRESPONCE` (`sp_NC_KQ_ENTRYRESPONCE_ACK`), `KingdomQuestGateKeeper::erb_happen`/`erb_CheckLinkMap`; `KQElement` (`kqe_QuestStart`, `kqe_ReviveNumber`, `kqe_GetKQLimitMin`), `KQPlayerInfoList` (knockdowns, team type), `sp_EscapeFromEmptyKQ`, `sp_KQItemDrop` (`KQItem.shn`), `wms_NC_KQ_W2Z_DESTROY_CMD`, `wms_NC_KQ_PLAYER_DISJOIN_CMD`, `wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD`, `Z2W_END_CMD`, comandos GM `&kingdomquest` / `&Reroadkingdomquest` / `ac_EndOfKingdomQuest` / `ac_KQReward`. Quitar `if kqID == 0 { kqID = 1 }`.

### 27D — Lua y recompensas
Bindings C-API de KQ en `pkg/script/lua/` desensamblados uno por uno (`cEndOfKingdomQuest` @ `0x005DF1F0` y el resto de `c*KQ*`), `AxialListKQReward`/`KQRewardLua`/`KQEnd`/`KQEndLua`, `ShinePlayer::so_ply_KQRewardStruct`, `PineEventScriptNode::ShineEndOfKingdomQuest`. Ejecución real de los 23 directorios de `data/LuaScript/KQ/` con el `ScriptLanguage`/`InitValue` de cada fila. Recompensa vía `Character.exe` (`CPFsCharacter::fc_NC_KQ_REWARD_REQ`, `p_Char_GetKQMap`) y `KingdomQuestRew.shn`; guardado/restauración de posición previa (`nKQHandle`, `sKQMap`, `nKQX/Y/D`, `dKQDate`, `PROTO_NC_CHAR_KQMAP_CMD`). GameLog `KQ_ENTER/LEAVE`.

### 27E — Verificación por comportamiento (no por tamaño)
1. **Tests de secuencia de paquetes** reproduciendo la sesión real del cliente: login de personaje → `0x581B` → se esperan exactamente `0x581C` (1ª vez) + `0x581D` en chunks de ≤ `0x1D26` bytes; segundo `0x581B` → sin `0x581C` y solo diffs.
2. Reloj simulado: una fila con `ST_Hour:ST_Minute` y `NextStartDeleyMin` debe pasar por todos los estados oficiales y emitir `JOINING_ALARM_CMD`, `W2Z_MAKE_REQ`, `W2Z_START_CMD`, `NOTIFY_CMD`, `JOINING_ALARM_END_CMD` y volver a agendarse según `RepeatMode`/`RepeatCount`.
3. **Regla anti-invención** (nuevo `pkg/audit/fase27_kq_no_invention_test.go`): falla si `pkg/server/world/kq*`/`pkg/server/zone/*kq*` contienen literales de tiempo o defaults sin un comentario `// exe: 0x...` que los respalde.
4. **Prueba manual con el cliente** (obligatoria, no se marca "Completado" sin ella): la lista de KQ aparece, el aviso de inscripción aparece, unirse, entrar al mapa, completar y recibir recompensa. Adjuntar log.

---

## 3. Prompt listo para copiar y pegar (por sub-fase)

```text
/goal Ejecutá la sub-fase 27<X> de docs/plan/FASE_27_KINGDOM_QUEST_CKQSERVER_ZONE_LUA_PORT_1A1_EXE.md siguiendo AGENTS.md.
Regla Mandatoria: cada función se porta desde su desensamblado (python3 scripts/disasm_pdb_handler.py world|zone|character <símbolo>), con la dirección del exe en un comentario. Prohibido inventar estados, timers, defaults, códigos de error, chunking u opcodes de respuesta; si algo no puede confirmarse en el exe/pdb/data, detenete y reportalo en vez de suponer.
Los tests deben validar COMPORTAMIENTO (secuencia y contenido de paquetes, transiciones de estado con reloj simulado), no solo sizeof/presencia.
Al terminar: go build ./... && go test -p 1 -count=1 ./..., actualizá el estado en docs/plan/README.md y hacé commit en main.
```

---

## 4. Criterios de cierre

1. `kq_manager.go` reemplazado; cero constantes sin `// exe: 0x...`.
2. Tests 27E 1–3 en verde; suites `verify_*_100.sh` previas siguen en `exit 0`.
3. Prueba manual 27E-4 con el cliente real documentada (log adjunto en `logs/`).
4. Commit: `feat(fase27): port 1:1 de Kingdom Quest (CKQServer + KingdomQuest Zone + Lua + rewards)`.


## 5. Prompts 

/goal Cerrá la sub-fase 27B de docs/plan/FASE_27_KINGDOM_QUEST_CKQSERVER_ZONE_LUA_PORT_1A1_EXE.md siguiendo AGENTS.md.
1) Corregí el fallo en TestCKQServer_JoinDisjoinTeamSelectAndVoteKick_Behavior (pkg/server/world/kq_test.go:502: expected wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD after passed vote) desensamblando Send_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD y VoteProcessing en WorldManager.exe.
2) Verificá que todos los requisitos de 27B estén completos y con comentario // exe: 0x....
3) Corré go test -v ./pkg/server/world -run TestCKQServer, marcá 27B como Completado en docs/plan/README.md y hacé commit en main. No toques 27C todavía.


/goal Ejecutá la sub-fase 27C de docs/plan/FASE_27_KINGDOM_QUEST_CKQSERVER_ZONE_LUA_PORT_1A1_EXE.md siguiendo AGENTS.md.
Regla Mandatoria: portá cada función desde su desensamblado (python3 scripts/disasm_pdb_handler.py zone <símbolo>) con // exe: 0x..., sin inventar defaults ni opcodes. Quitá el hack if kqID == 0 { kqID = 1 }.
Al terminar: go build ./..., corré los tests nuevos de KQ en pkg/server/zone, marcá 27C como Completado en docs/plan/README.md y hacé commit en main. No toques 27D todavía.
PROHIBIDO correr go test ./... ni suites completas de zone/world (yo las corro afuera); ejecutá únicamente el comando go test puntual indicado arriba.

/goal Ejecutá la sub-fase 27D de docs/plan/FASE_27_KINGDOM_QUEST_CKQSERVER_ZONE_LUA_PORT_1A1_EXE.md siguiendo AGENTS.md.
Regla Mandatoria: desensamblá cada binding c*KQ* y los handlers de Character.exe/Zone.exe con python3 scripts/disasm_pdb_handler.py {zone|character} <símbolo>, dejando // exe: 0x... en cada uno.
Al terminar: go build ./..., corré los tests de pkg/script/lua y recompensas KQ, marcá 27D como Completado en docs/plan/README.md y hacé commit en main. No toques 27E todavía.
PROHIBIDO correr go test ./... ni suites completas de zone/world (yo las corro afuera); ejecutá únicamente el comando go test puntual indicado arriba.

/goal Ejecutá los puntos 1, 2 y 3 de la sub-fase 27E de docs/plan/FASE_27_KINGDOM_QUEST_CKQSERVER_ZONE_LUA_PORT_1A1_EXE.md siguiendo AGENTS.md.
Creá pkg/audit/fase27_kq_no_invention_test.go y los tests de secuencia real de paquetes y reloj simulado.
Al terminar: go build ./... && go test -v ./pkg/audit -run TestFase27 && go test -v ./pkg/server/world -run TestCKQServer, actualizá docs/plan/README.md (dejando pendiente solo la prueba manual 27E-4 con el cliente) y hacé commit en main.
PROHIBIDO correr go test ./... ni suites completas de zone/world (yo las corro afuera); ejecutá únicamente el comando go test puntual indicado arriba.