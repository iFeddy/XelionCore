# Hoja de Ruta Modular Completa por Fases (`XelionCore` — 100% Fidelidad PDB + `.exe` + `.bak` + `data/` + MySQL)

Esta carpeta organiza la ejecución modular, autónoma y segura para la cuota de tokens de las etapas restantes del clúster oficial de Fiesta Online sobre `main`:

- **Etapa 1 (`LoginService`, `Account`, `AccountLog`, `StatisticsData`):** **Completada y auditada al 100%** (archivada en [fases_terminadas/](./fases_terminadas/)).
- **Etapa 2 (`WorldService` y `Character` — Fases `01` a `07`):** Gobernada por [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md).
- **Etapa 3 (`ZoneService` y `GameLog` — Fases `08` a `16`):** Gobernada por [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md).

---

## Reglas de Seguridad ante Corte de Tokens

1. **Una sesión limpia de chat por cada Fase (o ejecución secuencial con `/goal`):** Cada archivo `FASE_XX_*.md` contiene su prompt exacto listo para copiar y pegar y criterios objetivos de cierre.
2. **Compilación garantizada entre fases:** El orden secuencial (`01` -> `02` -> `02A` -> `03`..`07` -> `08`..`16`) asegura que al terminar cada fase `go build ./...` y todos los tests acumulados (`go test -p 1 -count=1 ./...`) sigan pasando en verde.
3. **Un `git commit` al cerrar cada Fase:** Cada fase termina haciendo commit en `main` **únicamente** cuando pasan al 100% sus controles automáticos.
4. **Cero Datos Hardcodeados y Cero Omisiones:** Todo nombre de tabla (`176/176`), Stored Procedure, OpCode y struct C++ proviene 1:1 de los PDBs (`Login.md`, `Account.md`, `AccountLog.md`, `WorldManager.md`, `Character.md`, `Zone.md`, `GameLog.md`), `.bak` (`TablasLegacy.md`, `docs/sql/procedures/`), binarios `.exe` y recursos `data/` (`.shn`, `.json`, `.txt`, `BlockInfo`, `AreaBMP`, `LuaScript`).

---

## Etapa 2: `WorldService` y `Character` (Fases `01` a `07`)

Plan Maestro: [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md)

| Orden | Archivo de Fase | Objetivo Principal | Nivel Recomendado | Control de Auditoría que Activa | Estado |
| :---: | :--- | :--- | :---: | :--- | :---: |
| **1** | [FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md](./FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md) | Creación inicial de 91 tablas de `World00_Character` y semillas base en `xelion_online` | **`Medium`** | **Control 1 (base)** | Completado |
| **2** | [FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md](./FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md) | Transpilación y validación de los 311 Stored Procedures en `database/procedures_character.sql` | **`Medium`** *(o `High`)* | **Control 2 (base)** | Completado |
| **2A** | [FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md](./FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md) | **Saneamiento 100% Fases 1 y 2:** `114/114` tablas `.bak` + `2` S2S (`tActiveZoneTransfer`, `tWorldZoneNode`), limpieza de `*_Template` falsos en `tCharacter`, regeneración completa de `default/tutorial_character_data.json` y los **`425` SPs** pool-safe (`proc_logic` + `SELECT` final) | **`High`** | **Controles 1 y 2 al 100%** | Completado |
| **3** | [FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md](./FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md) | Migración de los 50 métodos inline + todos los dominios a **`CALL <SP>`** en `pkg/database` y eliminación de `TransferManager` en memoria | **`High`** | **Control 3** (AST Go: cero SQL inline / cero memoria S2S) | Completado |
| **4** | [FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md](./FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md) | 100% de OpCodes y structs C++ no-`OpTool` de `WorldManager.md` y `Character.md` en `pkg/protocol/`, erradicación de los 19 alias privados y paridad `sizeof` | **`High`** | **Control 4** (Paridad OpCodes + `sizeof` PDB) | Completado |
| **5** | [FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md](./FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md) | Ciclo completo de cuenta/avatares, estilo, opciones `IMPROVE`, tutorial, asistencia y carga de los **`74` `.shn` únicos** (`54 + 71`) sin hardcode | **`High`** | **Control 5A** (Dispatcher Avatar/Char + 74 SHNs) | Completado |
| **6** | [FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md](./FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md) | Cobertura **100% de los `167` handlers C2S** de `WorldManager.pdb`: Chat, Party, Raid, Friends, KQ, Guild/War/Tournament/Emblem, Academy, HolyPromise, Prison, Report y Topología `tWorldZoneNode` | **`High`** | **Control 5B** (`167/167` C2S + S2S) | Completado |
| **7** | [FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md](./FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md) | Traspaso multi-proceso `World <-> Zone`, carga real de los 16 paquetes del In-Game Chunk, portales inter-zona, logout completo y suite final `verify_world_character_100.sh` | **`High`** | **Controles 6, 7 y 8** (`8/8 PASS`, `exit 0`) | **Completado** |

---

## Etapa 3: `ZoneService` y `GameLog` (Fases `08` a `16`)

Plan Maestro: [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md)

| Orden | Archivo de Fase | Objetivo Principal | Nivel Recomendado | Control de Auditoría que Activa | Estado |
| :---: | :--- | :--- | :---: | :--- | :---: |
| **8** | [FASE_08_ESQUEMA_PROCEDURES_Y_REPOSITORIO_GAMELOG.md](./FASE_08_ESQUEMA_PROCEDURES_Y_REPOSITORIO_GAMELOG.md) | Retiro de `database/zone_tables.sql`, cierre **`176/176` tablas** en `xelion_online`, **`27` SPs pool-safe** en `database/procedures_gamelog.sql`, `GameLogRepository` (`100% CALL`) y `GameLogDispatcher` asíncrono | **`Medium`** *(o `High`)* | **Controles 1, 2 y 3 (Etapa 3)** | **Completado** |
| **9** | [FASE_09_PROTOCOLO_OPCODES_Y_STRUCTS_ZONE_GAMELOG.md](./FASE_09_PROTOCOLO_OPCODES_Y_STRUCTS_ZONE_GAMELOG.md) | Registro 100% de los **`245` OpCodes C2S** + S2S, las **`44` structs `PROTO_NC_LOG_*`** (`GameLog.md`) y las **`1.293` structs `PROTO_NC_*`** (`Zone.md`) en `pkg/protocol/` con paridad `sizeof` | **`High`** | **Control 4 (Etapa 3)** | **Completado** |
| **10** | [FASE_10_UNIVERSO_DATOS_SHN_BLOCKINFO_AREABMP_Y_LUA_API.md](./FASE_10_UNIVERSO_DATOS_SHN_BLOCKINFO_AREABMP_Y_LUA_API.md) | Carga de **`180` `.shn` reales** (`194` rutas) + `49` checksums, paquetes **`pkg/data/blockinfo`** (`301` archivos `.shbd`/`.sbi`/`.aid`/`.bdt`/`.shab`) y **`pkg/data/areabmp`** (`44` `.bmp`), y **`274` funciones C-API `lua_State*`** en `pkg/script/lua/` | **`High`** | **Control 6 (Etapa 3)** | **Completado** |
| **11** | [FASE_11_MOTOR_ESPACIAL_AOI_COLISIONES_BRIEFINFO_Y_MOVIMIENTO.md](./FASE_11_MOTOR_ESPACIAL_AOI_COLISIONES_BRIEFINFO_Y_MOVIMIENTO.md) | Game Loop determinístico, grilla espacial AOI, `BRIEFINFO`, colisiones `.shbd`/`.sbi`, las **27 clases** en `ClassParam`, movimiento, tiendas `Tent`, recolección `Gather` y monturas `MOVER` | **`High`** | Suite AOI, Colisiones y Movimiento | Pendiente |
| **12** | [FASE_12_MOTOR_COMBATE_SKILLS_ABSTATE_SOULSTONE_Y_PVP.md](./FASE_12_MOTOR_COMBATE_SKILLS_ABSTATE_SOULSTONE_Y_PVP.md) | Motor de combate físico/mágico (`NC_BAT`), `ActiveSkill`/`PassiveSkill`/campos mágicos (`NC_SKILL`), máquina `AbState`, `SoulStone`, stats libres, muerte/resurrección, PvP/PK y `GameLog` Tipos `2, 8, E` | **`High`** | Suite Combate, Skills y `AbState` | Pendiente |
| **13** | [FASE_13_MOTOR_MONSTRUOS_IA_MOBBEHAVIOR_DROPS_Y_COLECCION.md](./FASE_13_MOTOR_MONSTRUOS_IA_MOBBEHAVIOR_DROPS_Y_COLECCION.md) | Spawns `MobRegen`, patrullas `MobRoam`, IA `MobBehavior` + Lua `AIScript`, `MobAttackSequence`, EXP Party/Raid, tablas de `Drops`, `AutoPick`, Títulos (`NC_CT`), Cartas (`NC_COLLECT`) y `GameLog` Tipo `6` | **`High`** | Suite IA Monstruos, Drops y Colección | Pendiente |
| **14** | [FASE_14_ECONOMIA_INVENTARIO_COMERCIO_BOOTH_UPGRADE_Y_CRAFTING.md](./FASE_14_ECONOMIA_INVENTARIO_COMERCIO_BOOTH_UPGRADE_Y_CRAFTING.md) | **`48/48` C2S `NC_ITEM`** + **`62` métodos `iac_*`**, Mercaderes NPC, almacenes, `SafeTrade` (`10` C2S), Kioscos `BOOTH` (`8` C2S), `Upgrade`, `Enchant`, `RandomOption`, `Dismantle`, `Mix`, `Produce` y `GameLog` Tipos `3, 4, 5, 7, B, C, D` | **`High`** | Suite Transaccional Economía y Forja | Pendiente |
| **15** | [FASE_15_QUESTS_NPCS_LUA_INSTANCIAS_KQ_MINIHOUSE_PETS_Y_CASINO.md](./FASE_15_QUESTS_NPCS_LUA_INSTANCIAS_KQ_MINIHOUSE_PETS_Y_CASINO.md) | Motor `QUEST` (`7` C2S), NPCs/escenarios, instancias Lua (`InstanceDungeon`, `KQ` en mapa, `GuildTournament` en mapa, `Wedding`), `HolyPromise`, `MiniHouse` (`23` C2S), `Pet`, `WeaponTitle` y Casino (`DICE` `9` C2S + `GAMBLE` `18` C2S) | **`High`** | Suite Quests, Lua, MiniHouse, Pets y Casino | Pendiente |
| **16** | [FASE_16_COORDINACION_S2S_ADMIN_E2E_CLUSTER_Y_AUDITORIA_FINAL.md](./FASE_16_COORDINACION_S2S_ADMIN_E2E_CLUSTER_Y_AUDITORIA_FINAL.md) | Coordinación S2S `Zone <-> World` (**`138` flujos `wms_NC_*`**), multi-nodo `Zone00`..`Zone04`, comandos GM (`AdminLvSet.shn`), cierre **`245/245` C2S** y suite final `scripts/verify_zone_gamelog_100.sh` | **`High`** | **Controles 5, 7 y 8** (`8/8 PASS`, `exit 0`) | Pendiente |

---

## Archivo Histórico de Etapas Anteriores

- [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./fases_terminadas/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md): Plan maestro de la Etapa 1 completada al 100% (`LoginService`, `Account`, `AccountLog`, `StatisticsData`).
- [fases_terminadas/](./fases_terminadas/): Especificaciones históricas de las Fases 01 a 07 de la Etapa 1.
