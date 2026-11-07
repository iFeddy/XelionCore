# Fase 15 — Misiones (`QUEST`), Diálogos NPC, Instancias `LuaScript` (`ID`/`KQ`/`GT`/`Wedding`), `MiniHouse`, Mascotas (`PET`) y Casino (`DICE`/`GAMBLE`)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 1-Bugs #7, #14, #15 y #16, 4.1-Depts 4/8/12/15/17/22/23/27/29/35/37/41/43/47/53 y 6.5**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 a 14 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 15 (`docs/plan/FASE_15_QUESTS_NPCS_LUA_INSTANCIAS_KQ_MINIHOUSE_PETS_Y_CASINO.md`) leyendo obligatoriamente las Secciones 1 (Bugs #7, #14, #15 y #16), 4.1 (Depts 4, 8, 12-`DICE_GAME_START`, 15, 17, 22, 23, 27, 29-`TOURNAMENT`, 35, 37, 41, 43, 47 y 53) y 6.5 de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: 100% de las misiones, diálogos, instancias, KQs, muebles, mascotas y juegos de casino se rigen por `data/LuaScript/`, `data/Script/`, `data/World/` y las tablas `.shn` oficiales, persistiendo en MySQL con `CALL <SP>` y auditando en `GameLog`.
1) Implementá el motor completo de Misiones (`Dept 17` 7 C2S: `QUEST_SCRIPT_CMD_ACK`, `GIVE_UP_REQ`, `CLIENT_SCENARIO_DONE_REQ`, `SELECT_START_REQ`, `REWARD_SELECT_ITEM_INDEX_CMD`, `START_REQ`, `READ_REQ` con `QuestData.shn`, `QuestParser.json`, `QuestDialog.shn` y `CALL p_Quest_*`), Diálogos/Menús NPC (`MENU_SERVERMENU_ACK`, `ACT_NPCCLICK_CMD`, `NPCAction.json`), Escenarios (`Dept 27` 2 C2S) y Tutorial en mapa (`CHAR_TUTORIAL_STEP_ACK`).
2) Enlazá `pkg/script/lua/` con mapas instanciados reales: Mazmorras Instanciadas (`Dept 41` 2 C2S + `InstanceDungeon.json` + `data/LuaScript/ID/`), ejecución en mapa de Kingdom Quests (`Dept 22` `KQ_ENTRYRESPONCE_ACK` + `data/LuaScript/KQ/` + recompensas `KQReward`), Guild Tournament en mapa (`Dept 29` 6 C2S de captura de bandera, skills, dados y observadores + `GuildTournament.json`), Bodas (`Dept 8` `WEDDING_*` 8 C2S + `CHAR_WEDDING_PARTNER_INFO_REQ`) y recompensas `HolyPromise` (`Dept 37` 2 C2S).
3) Implementá `MiniHouse` completa (`Dept 35`: los 23 C2S de construcción, mobiliario, visitas, contraseñas, portales y emociones en `tMinihouse*` + `MiniHouse*.shn`), Mascotas (`Dept 53` `PET_*` 3 C2S + invocación, hambre, evolución, tendencias y skills en `tPet`/`tPetNew`/`tPetSkill` + `Pet*.shn`), `WeaponTitle` (`Dept 23` `WT_LICENSE_REQ`) y Casino (`Dept 43` `DICE_TAISAI_*` 9 C2S + `ITEM_DICE_GAME_START_REQ` + `Dept 47` `GAMBLE_*` 18 C2S contra `DiceGame.shn`, `SlotMachine.shn` y `tDice*`/`tGamble*`), emitiendo `GameLog` Tipos `9, F, G, H, I, J, K`.
4) Creá tests integrales para Quests, instancias Lua (`ID` y `KQ`), `MiniHouse`, `Pet` y Casino (`Dice`/`Gamble`). Al terminar, marcá Fase 15 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 16.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Misiones (`NC_QUEST` — `7` C2S), Diálogos NPC (`NC_MENU`), Escenarios (`NC_SCENARIO`) y Tutorial:**
   - Implementar `PROTO_NC_QUEST_START_REQ` (`0x4414`), `SELECT_START_REQ` (`0x440F`), `READ_REQ` (`0x4416`), `GIVE_UP_REQ` (`0x4407`), `REWARD_SELECT_ITEM_INDEX_CMD` (`0x4411`), `SCRIPT_CMD_ACK` (`0x4402`) y `CLIENT_SCENARIO_DONE_REQ` (`0x440B`) validando pre-requisitos de nivel/clase/misión previa/ítem en `QuestData.shn` y `QuestParser.json`, actualizando automáticamente contadores al matar monstruos (`PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD`), recolectar ítems o hablar con NPCs, entregando EXP/Cen/Fama/Ítems vía `CALL p_Quest_*` y auditando en `GameLog` Tipo `9`.
   - Completar `MENU_SERVERMENU_ACK` (`0x3C02`), `ACT_NPCCLICK_CMD` (`0x200A`), `ACT_NPCMENUOPEN_ACK` (`0x201D`), `ACT_NPC_MENU_CMD` (`0x2073`), `SCENARIO_AREAENTRY_ACK` (`0x6C06`), `SCENARIO_AREALEAVE_ACK` (`0x6C08`) y `CHAR_TUTORIAL_STEP_ACK` (`0x1117`).

2. **Instancias Dinámicas con Runtime `LuaScript` (`InstanceDungeon`, `KingdomQuest`, `GuildTournament`, `Wedding`, `HolyPromise`):**
   - **Instance Dungeons (`Dept 41` `2` C2S):** `INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK` (`0xA40A`) e `INSTANCE_DUNGEON_RANK_LIST_REQ` (`0xA413`), creando copias aisladas de `ZoneMap` por grupo/raid gobernadas por `InstanceDungeon.json` y `data/LuaScript/ID/` (incluyendo apertura/cierre real de puertas `.sbi`).
   - **Kingdom Quest en Mapa (`Dept 22` `1` C2S):** `KQ_ENTRYRESPONCE_ACK` (`0x581A`), ciclo completo de scripts `data/LuaScript/KQ/`, marcador/temporizador, entrega de recompensas (`KQReward`) y `GameLog` Tipo `H` + `PROTO_NC_LOG_GAME_KQ_ENTER`/`LEAVE`.
   - **Guild Tournament en Mapa (`Dept 29` `6` C2S):** `GUILD_TOURNAMENT_RECALL_ENTER_REQ` (`0x749C`), `FLAGCAPTURE_REQ` (`0x74CA`), `USESKILL_REQ` (`0x74CF`), `DICEGAME_THROW_REQ` (`0x74D3`), `OBSERVER_ENTER_REQ` (`0x74D7`) y `OBSERVER_OUT_REQ` (`0x74DB`) con `GuildTournament.json`.
   - **Bodas (`Dept 8` `8` C2S + `Dept 4`) y `HolyPromise` (`Dept 37` `2` C2S):** `ACT_WEDDING_PROPOSEREQ_REQ`..`DIVORCE_REFUSE_CMD`, `CHAR_WEDDING_PARTNER_INFO_REQ` (`0x1093`), `HOLY_PROMISE_GET_CEN_REWARD_REQ` (`0x943C`) y `WITHDRAW_CEN_REWARD_REQ` (`0x943E`) (`GameLog` Tipos `J` y `K`).

3. **`MiniHouse` (`23` C2S), Mascotas (`PET` + `WT`) y Casino (`DICE` `9` C2S + `GAMBLE` `18` C2S):**
   - **`MiniHouse` (`Dept 35` — los `23` C2S completos):** construcción, entrada/salida, visitas, expulsión, contraseña, título, cupo máximo, colocación/reubicación/borrado de muebles (`Box 12`), efectos/emociones de mobiliario y portales (`MiniHouse*.shn` + `CALL p_Minihouse_*` + `GameLog` Tipo `F`).
   - **Mascotas (`Dept 53` `3` C2S: `PET_ASK_NEW_NAME_ACK`, `PET_SET_NAME_REQ`, `PET_SET_NAME_CANCEL_REQ` + ciclo en mapa) y `WeaponTitle` (`Dept 23` `WT_LICENSE_REQ`):** invocación, saciedad, ánimo/estrés, evolución, habilidades de mascota (`Pet*.shn` + `CALL p_Pet_*` / `usp_Pet_*` + `GameLog` Tipo `G`) y licencias de título de arma (`WeaponTitleData.shn` + `iac_wptitle_*`).
   - **Casino (`Dept 43` `9` C2S `DICE_TAISAI_*` + `0x303C` `ITEM_DICE_GAME_START_REQ` + `Dept 47` `18` C2S `GAMBLE_*`):** compra/venta de fichas `tGambleCoin`, mesas `Taisai` y `SlotMachine` con jackpots y rankings (`DiceGame.shn`, `SlotMachine.shn`, `tDice*`, `tGamble*` y `GameLog` Tipo `I`).

---

## Criterio de Cierre de la Fase 15

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/server/zone/... ./pkg/script/...` -> `PASS` (incluyendo tests de Quests, instancias Lua `ID`/`KQ`, `MiniHouse`, `Pet`, `WeaponTitle` y `Dice`/`Gamble`).
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 15 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(zone,lua): Fase 15 - Quests, NPCs, instancias Lua ID/KQ/GT/Wedding, MiniHouse, Pets y Casino con GameLog"`
