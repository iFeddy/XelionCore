# Plan Maestro Exhaustivo — Etapa 3 (100% Fidelidad PDB + Binarios + `.bak` + SHN/JSON/BlockInfo/Lua + MySQL): `ZoneService` y `GameLog`

> [!IMPORTANT]
> **Alcance y Reglas Mandatorias de este Plan Maestro (Auto-Contenido para Ejecución Directa):**
> - **Cero omisiones, cero resúmenes y cero datos hardcodeados:** Este documento inventaría al **100% y sin abreviaturas** las **2 tablas oficiales** de `World00_GameLog.bak` (cerrando las **176/176 tablas oficiales** de [`docs/TablasLegacy.md`](./TablasLegacy.md)), los **27 Stored Procedures** de `World00_GameLog` (`23` ODBC de `GameLog.exe` $\cup$ `26` de `World00_GameLog.bak`), las **44 estructuras `PROTO_NC_LOG_*`** de [`docs/Servicios/GameLog.md`](./Servicios/GameLog.md), las **245 entradas oficiales C2S** del dispatcher `protocolstore(ShinePlayer)` (`0x004D3510`) de `Zone.exe`, los **62 métodos `iac_*` (`ItemAttributeClass`)**, los **138 handlers S2S `wms_NC_*`**, los **196 handlers S2S `gds_NC_*`**, las **1.343 estructuras C++** de [`docs/Servicios/Zone.md`](./Servicios/Zone.md) (`1.293` structs `PROTO_NC_*` de cable + `9` `ShineSendType<PROTO_NC_*>` + `1` `OpTool` excluida + `40` artefactos STL TPI), las **199 referencias `.shn`** de `Zone.exe` (`194` rutas presentes en `data/` y `data/View/` = `180` basenames únicos en disco + `5` opcionales ausentes en `AstraServer` = `185` basenames únicos), los **31 archivos `.txt`**, los **301 archivos binarios de colisión/geometría** en `data/BlockInfo/`, las **44 máscaras `.bmp`** en `data/AreaBMP/` y las **274 funciones C-API `lua_State*`** de `Zone.pdb`.
> - **100% MySQL (`xelion_online`) + `data/` como únicas fuentes de verdad:** Prohibido terminantemente ejecutar SQL inline (`SELECT`/`INSERT`/`UPDATE`/`DELETE`) en Go, mantener tablas fantasma duplicadas (`database/zone_tables.sql`), hardcodear fórmulas de daño/stats (`* 3`, `NextExp = 280`), hardcodear mapas (`"RouN"`) o privilegios GM, o dejar stubs mudos en el motor Lua.
> - **Únicas exclusiones autorizadas (Sección 7):** `GamigoZR` (mock DRM HTTP 2015 descartado por [`AGENTS.md`](../AGENTS.md)), `OpTool` (`PROTO_NC_OPTOOL_CONNECT_BRIF_ACK` y las 3 tablas de `OperatorTool.bak`), y los `40` artefactos internos de plantillas C++ STL (`std::vector` / `std::allocator` / `_Iterator`) emitidos por el volcado TPI de `Zone.pdb`.

---

## 0. Decisiones de Arquitectura Cerradas y Obligatorias

Quien ejecute cualquiera de las fases de la Etapa 3 (`FASE_08` a `FASE_16`) debe cumplir **sin excepción ni reinterpretación** las siguientes 6 decisiones arquitectónicas:

1. **Cero Tablas SQL Propias en `ZoneService` y Cierre `176/176` en `xelion_online`:**
   - En la arquitectura oficial (`Zone.exe`), el servidor de zona **no importa `ODBC32.dll` ni posee base de datos propia**: delega toda persistencia de personajes/ítems/misiones/gremios a `World00_Character` (`425` SPs ya compilados y verificados en la Etapa 2 / Fase 02A) y toda auditoría de eventos de juego a `World00_GameLog` (**`27` SPs** compilados en `database/procedures_gamelog.sql` en la **Fase 08**).
   - El archivo obsoleto [`database/zone_tables.sql`](../database/zone_tables.sql) (borrador temprano con `86` tablas duplicadas y divergentes respecto a `database/schema_character.sql`) **debe eliminarse definitivamente en la Fase 08**.
   - Con `tGameLog` (presente en [`database/schema.sql`](../database/schema.sql)) y `TEMP_QuestComplete` (presente en [`database/schema_character.sql`](../database/schema_character.sql)), la base única `xelion_online` cubre el **100% de las 176 tablas oficiales** de [`docs/TablasLegacy.md`](./TablasLegacy.md) (`34` `Account` + `10` `AccountLog` + `13` `StatisticsData` + `114` `World00_Character` + `2` `World00_GameLog` + `3` `OperatorTool`).
2. **`GameLog` Asíncrono No Bloqueante y 100% Pool-Safe (`CALL p_GameLog_*`):**
   - En el legado, `Zone.exe` y `WorldManager.exe` enviaban paquetes `PROTO_NC_LOG_*` por TCP al puerto `9034` de `GameLog.exe`, quien insertaba en `tGameLog` vía ODBC.
   - En `XelionCore`, `GameLogRepository` (`pkg/database/mysql_gamelog.go`) encapsula el **100% de los `27` Stored Procedures** mediante `CALL <SP>` con la convención pool-safe (`proc_logic: BEGIN /* validaciones e inserción */ END proc_logic; SELECT out_nLogNo AS nLogNo;`), alimentado por un `GameLogDispatcher` asíncrono con canal bufferizado en Go y `Flush/Close` determinístico al apagar el nodo, garantizando **cero bloqueo de E/S en el Game Loop de `ZoneService`**.
3. **Topología Multi-Instancia `Zone00`..`Zone04` Gobernada por `FieldList.json` y `tWorldZoneNode`:**
   - `ZoneService` soporta ejecución tanto en nodo único (desarrollo) como desplegado en las 5 instancias oficiales (`Zone00` puerto `9016`, `Zone01` puerto `9019`, `Zone02` puerto `9022`, `Zone03` puerto `9025`, `Zone04` puerto `9028`), filtrando qué mapas aloja según el `zoneNumber` declarado en `data/World/Map/FieldList.json` y anunciándose en `tWorldZoneNode` (`CALL usp_WorldZone_UpsertNode` / `CALL usp_WorldZone_Heartbeat`).
4. **Game Loop Determinístico por Mapa + Particionamiento Espacial AOI (Sector Grid):**
   - Está prohibido difundir paquetes de movimiento, combate o aparición a *todos* los jugadores del mapa (`O(N^2)`). Cada `ZoneMap` mantiene una grilla espacial de sectores AOI acorde al rango visual oficial de `Zone.exe`, calculando altas y bajas incrementales de visibilidad (`PROTO_NC_BRIEFINFO_*_CMD` y `PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD`) en cada cambio de celda.
5. **Geometría, Colisiones y Fórmulas 100% Gobernadas por `data/` (Cero Magia en Código Go):**
   - Todo movimiento jugador/monstruo y línea de visión de habilidad valida contra los archivos binarios reales de `data/BlockInfo/` (`.shbd` walkability, `.sbi` puertas dinámicas `DoorBlock`, `.aid` áreas de evento) y `data/AreaBMP/` (`.bmp` máscaras de área).
   - Todo cálculo de HP/SP máximo, daño físico/mágico, crítico, evasión, ángulo (`DamageByAngle.json`), diferencia de nivel (`DamageLvGap*.json`), experiencia (`ExpTable.json`), drops (`ItemDropTable.json` / `ItemDropGroup.json`) y mejoras (`UpgradeInfo.shn`) proviene estrictamente de `data/`, prohibiendo constantes inventadas en Go.
6. **Nomenclatura Oficial 100% Idéntica a `Zone.pdb` y `GameLog.md` en `pkg/protocol/`:**
   - Todo paquete C2S y S2S usa su nombre exacto `PROTO_NC_*` y verifica mediante tests automáticos su tamaño exacto `sizeof` contra [`docs/Servicios/Zone.md`](./Servicios/Zone.md) y [`docs/Servicios/GameLog.md`](./Servicios/GameLog.md).

---

## 1. Diagnóstico Forense de Brechas y 18 Errores Críticos Detectados en el Estado Actual

La auditoría forense del código actual en `pkg/server/zone/`, `database/`, `pkg/data/`, `pkg/script/lua/` y `pkg/protocol/` frente a `Zone.exe`, `Zone.pdb`, `GameLog.exe` y `World00_GameLog.bak` identificó **18 brechas estructurales críticas** que esta Etapa 3 resuelve al 100%:

| # | Severidad | Archivo / Componente Actual | Bug / Brecha Detectada | Corrección Oficial en Etapa 3 |
| :---: | :---: | :--- | :--- | :--- |
| **1** | **CRÍTICO** | [`database/zone_tables.sql`](../database/zone_tables.sql) | Archivo huérfano con `86` tablas duplicadas y divergentes respecto a `database/schema_character.sql`; contradice `Zone.exe` (que tiene `0` tablas SQL propias). | **Fase 08:** Eliminar `database/zone_tables.sql` y validar en auditoría que `xelion_online` contiene exactamente el universo oficial sin DDL duplicado. |
| **2** | **CRÍTICO** | `database/` (ausencia de `procedures_gamelog.sql`) | Ninguno de los **`27` Stored Procedures** de `World00_GameLog` (`p_GameLog_Add*`, `usp_Login_setChannel`, `p_GameLog_Get*`) está transpilado ni cargado en MySQL, y no existe `GameLogRepository` en `pkg/database/`. | **Fase 08:** Transpilar los `27` SPs pool-safe en `database/procedures_gamelog.sql`, integrarlos en `scripts/apply_mysql_procedures.sh` y crear `pkg/database/mysql_gamelog.go` + worker asíncrono. |
| **3** | **CRÍTICO** | `pkg/protocol/` (ausencia de `PROTO_NC_LOG_*` y brecha en `Zone.md`) | Faltan las **`44` estructuras `PROTO_NC_LOG_*`** de `GameLog.md` y cientos de las **`1.293` estructuras `PROTO_NC_*`** y **`245` OpCodes C2S** de `Zone.md`. | **Fase 09:** Implementar en `pkg/protocol/` el 100% de OpCodes C2S/S2S y las `1.293 + 44` estructuras con validación automática de `sizeof` byte a byte. |
| **4** | **CRÍTICO** | [`pkg/server/zone/handler.go`](../pkg/server/zone/handler.go#L53-L143) | El `switch pkt.OpCode` solo contempla `40` casos de los **`245` OpCodes C2S** de `Zone.exe`, incluye alias privados (`PROTO_NC_CHAR_INGAME_REQ`, `PROTO_NC_ACT_REST_START_REQ` en vez de `PITCHTENT_REQ`/`FOLDTENT_REQ`) e intercepta `PROTO_NC_DATA_REPORT_*` (propios de `WorldManager`). | **Fases 09 y 11..16:** Purgar alias no oficiales y alcanzar cobertura real **`245/245` C2S** en `protocolstore(ShinePlayer)`. |
| **5** | **CRÍTICO** | [`pkg/server/zone/data_store.go`](../pkg/server/zone/data_store.go#L26-L233) | `DataStore` solo carga **9 archivos `.shn`** (`ItemInfo`, `ItemInfoServer`, `ItemSort`, `UseClassTypeInfo`, `ClassName`, `ItemUseEffect`, `ActiveSkill`, `PassiveSkill`, `MobInfo`) de las **`194` rutas `.shn` reales** en disco (`180` basenames únicos en disco + `5` opcionales = `185` basenames en `Zone.exe`) que lee `Zone.exe` (`UpgradeInfo`, `QuestData`, `Produce`, `ItemMix`, `ItemDismantle`, `ItemSocket`, `MiniHouse*`, `Pet*`, `Mover*`, `DiceGame`, `SlotMachine`, `AdminLvSet`, etc.). | **Fase 10:** Cargar e indexar en `DataStore` la totalidad de las **`194` rutas `.shn` en disco (`180` basenames únicos)** + verificar los `49` `.shn` de `ChecksumEngine`. |
| **6** | **CRÍTICO** | `pkg/data/` (ausencia de `blockinfo` y `areabmp`) | Existen **`301` archivos binarios** en `data/BlockInfo/` (`158` `.shbd`, `33` `.sbi`, `76` `.bdt`, `23` `.aid`, `8` `.shab`, `2` `.shmd`, `1` `.shad`) y **`44` `.bmp`** en `data/AreaBMP/`, pero **no existe ningún cargador en Go**: los jugadores atraviesan paredes, puertas cerradas y obstáculos sin validación alguna. | **Fase 10 y Fase 11:** Crear `pkg/data/blockinfo` y `pkg/data/areabmp` fieles a `MapBlock::MapBlockInformation`, `MapDoorArray`, `AreaInfoData` y `CFsAreaBMP` de `Zone.pdb`, y aplicarlos en movimiento, IA y combate. |
| **7** | **CRÍTICO** | [`pkg/script/lua/engine.go`](../pkg/script/lua/engine.go#L241-L287) | `lua.Engine` nunca se instancia en `ZoneService` y solo implementa `9` funciones básicas más `32` stubs mudos (`stubNoOp`/`stubReturnZero`), mientras `Zone.pdb` define **`274` funciones C-API `lua_State*`** (`246` `c*` de gameplay) de las cuales `222` son invocadas activamente en `data/LuaScript/`. | **Fase 10 y Fase 15:** Implementar las **`274` funciones C-API `lua_State*`** enlazadas a `ZoneService` para KQ, Instance Dungeons, IA de jefes, misiones y tutorial. |
| **8** | **CRÍTICO** | [`pkg/server/zone/map.go`](../pkg/server/zone/map.go#L44-L72) | `ZoneMap` solo tiene `players`, `npcs`, `gates` y `droppedItems`: **no tiene monstruos vivos (`mobs`), ni Game Loop por ticks, ni grilla espacial AOI** (difunde todo al mapa entero e ignora `MobRegen`, `MobRoam`, `MobSetting` y `MobAttackSequence` pese a cargarlos en RAM). | **Fases 11 y 13:** Incorporar Game Loop determinístico, grilla espacial AOI por sectores y ciclo de vida completo de monstruos (`MobRegen` -> `MobRoam` -> `MobBehavior` -> `MobAttack` -> `Drops`). |
| **9** | **ALTO** | [`pkg/server/zone/character.go`](../pkg/server/zone/character.go#L39-L56) | `JobToClassName` solo contempla 6 clases base y cae en `"Fighter"` para las otras **21 subclases** (CleverFighter, Warrior, Gladiator, Knight, Paladin, HighCleric, Guardian, HolyKnight, HawkArcher, Scout, SharpShooter, Ranger, WizMage, Enchanter, Warlock, Wizard, Gambbler, Renegade, Spectre, Reaper, Sentinel), corrompiendo los stats de `ClassParam` a partir de nivel 20. | **Fase 11:** Mapear las **27 clases oficiales** desde `ClassName.shn` / `SingleData.shn` y `data/World/ClassParam/<Class>.json` sin fallbacks falsos. |
| **10** | **ALTO** | [`pkg/server/zone/character.go`](../pkg/server/zone/character.go#L106-L148) y [`L241-L300`](../pkg/server/zone/character.go#L241-L300) | Fórmulas inventadas y valores hardcodeados: `Constitute * 3` (con comentario `// VER DE DONDE SALE`), `maxHP = 52`, `maxStonesHP = 15`, `PrevExp = 0`, `NextExp = 280`, `RedistributePoint = 0`, `adminlevel = 0`, `upgrade = {0,0,0}`. | **Fases 11 y 12:** Calcular stats exactamente con `FreeStatTable.json`, `Common.json`, `ExpTable.json`, mejoras reales de ítems equipados y campos reales de `tCharacter`. |
| **11** | **ALTO** | [`pkg/server/zone/map.go`](../pkg/server/zone/map.go#L176-L190) y [`character.go`](../pkg/server/zone/character.go#L263-L265) | Persisten fallbacks hardcodeados `normID = "RouN"`, `shortName = "RouN"` y dimensiones ficticias `width = 256, height = 256`. | **Fase 11:** Resolver siempre mapa y dimensiones reales desde `FieldList.json`, `MapInfo.shn` y cabecera `.shbd`, rechazando mapas inexistentes con código de error oficial. |
| **12** | **ALTO** | `pkg/server/zone/` (ausencia total de Combate, Skills y `AbState`) | **0 de 9** handlers C2S de `NC_BAT` (`Dept 9`) y **0 de 4** de `NC_SKILL` (`Dept 18`): no hay auto-ataque, casteo de skills de objetivo/campo, empoderamiento, buffs/debuffs `AbState`, ni PvP/PK. | **Fase 12:** Implementar motor completo de Combate, `ActiveSkill`/`PassiveSkill`, `AbState`, `SoulStone` y PvP/PK con fórmulas exactas de `data/World/Combat/`. |
| **13** | **ALTO** | `pkg/server/zone/` (brecha de `36/48` en `NC_ITEM` + `0/10` `TRADE` + `0/8` `BOOTH`) | Faltan compra/venta/recompra NPC, almacén personal/cuenta/gremio/academia, `SafeTrade` entre jugadores, tiendas `Booth` y búsqueda categorizada, `Upgrade` clásico/nuevo/brazalete/montura, `Enchant` gemas/sockets, `RandomOption`, `Mix`, `Dismantle`, `Produce` y `MysteryVault`. | **Fase 14:** Implementar los `48/48` handlers de `NC_ITEM`, `62` métodos `iac_*`, `10/10` `NC_TRADE` y `8/8` `NC_BOOTH` con transacciones MySQL atómicas y `GameLog`. |
| **14** | **ALTO** | `pkg/server/zone/` (ausencia de `QUEST`, `MENU` completo, `SCENARIO` e Instancias) | **0 de 7** handlers C2S de `NC_QUEST` (`Dept 17`), **0 de 2** de `NC_SCENARIO` (`Dept 27`), **0 de 2** de `NC_INSTANCE_DUNGEON` (`Dept 41`), **0 de 1** de `NC_KQ` en mapa (`Dept 22`) y **0 de 6** de `GUILD_TOURNAMENT` en mapa (`Dept 29`). | **Fase 15:** Implementar motor completo de Quests (`QuestData.shn` + `QuestParser.json` + `PineScript`), diálogos NPC, escenarios e instancias `ID`/`KQ`/`GT`/`Wedding` con Lua. |
| **15** | **ALTO** | `pkg/server/zone/` (ausencia de `MINIHOUSE`, `PET`, `MOVER`, `COLLECT`, `WT`) | **0 de 23** C2S de `NC_MINIHOUSE` (`Dept 35`), **0 de 3** C2S de `NC_PET` (`Dept 53`), **0 de 2** C2S de `NC_MOVER` (`Dept 51`), **0 de 4** C2S de `NC_COLLECT` (`Dept 49`) y **0 de 1** C2S de `NC_WT` (`Dept 23`). | **Fases 11, 13 y 15:** Implementar Monturas (`Fase 11`), Colección de Cartas y Títulos (`Fase 13`), y MiniHouse, Mascotas y WeaponTitle (`Fase 15`). |
| **16** | **ALTO** | `pkg/server/zone/` (ausencia de Casino: `DICE` y `GAMBLE`) | **0 de 9** C2S de `NC_DICE` (`Dept 43`) y **0 de 18** C2S de `NC_GAMBLE` (`Dept 47`: GambleHouse, Coin, Taisai y SlotMachine). | **Fase 15:** Implementar Taisai Dice y GambleHouse/SlotMachine leyendo `DiceGame.shn` y `SlotMachine.shn` y persistiendo en `tDice*` / `tGamble*` + `GameLog` Tipo `I`. |
| **17** | **ALTO** | `pkg/server/zone/` (ausencia de los `138` flujos S2S `wms_NC_*` y `AdminLvSet.shn`) | `ZoneService` no procesa los `138` flujos de coordinación con `WorldService` (`wms_NC_*`: avisos globales, sincronización de Party/Guild/War/Academy/Marriage/Prison) ni los comandos GM gobernados por `AdminLvSet.shn`. | **Fase 16:** Implementar la coordinación completa `Zone <-> World` (`138` flujos `wms_NC_*`) y el despachador de comandos GM validado contra `AdminLvSet.shn` + `tGMItemReward*`. |
| **18** | **ALTO** | [`cmd/ZoneService/main.go`](../cmd/ZoneService/main.go#L48-L62) | Arranca siempre con `MapID: 0` y puerto único `9210` sin filtrar mapas por `zoneNumber` (`Zone00`..`Zone04` -> `9016, 9019, 9022, 9025, 9028`) ni enlazar `GameLogDispatcher` ni `LuaEngine`. | **Fases 08, 10 y 16:** Parametrizar `ZONE_ID` (`0..4` o modo unificado), puertos oficiales `9016..9028`, `GameLogDispatcher`, `LuaEngine` y drenaje limpio en `SIGTERM`. |

---

## 2. Base de Datos `World00_GameLog.bak`, Cierre `176/176` Tablas, `27` Stored Procedures y `44` Estructuras `PROTO_NC_LOG_*`

### 2.1 Cierre Definitivo de las 176 Tablas Oficiales del Clúster (`docs/TablasLegacy.md`)

Con la culminación de la Etapa 1 (`Account`, `AccountLog`, `StatisticsData`), la Etapa 2 (`World00_Character`) y esta Etapa 3 (`World00_GameLog`), la base de datos única `xelion_online` cubre el **100% de las 176 tablas oficiales** extraídas de los `.bak` originales:

| Base de Datos Legacy (`.bak`) | Total Tablas Oficiales | Ubicación en `XelionCore` | Estado |
| :--- | :---: | :--- | :---: |
| `Account.bak` | `34` | [`database/schema.sql`](../database/schema.sql) (`tItem` de tienda resuelta como `tAccountItem`) | Cubierto (Etapa 1) |
| `AccountLog.bak` | `10` | [`database/schema.sql`](../database/schema.sql) | Cubierto (Etapa 1) |
| `StatisticsData.bak` | `13` | [`database/schema.sql`](../database/schema.sql) | Cubierto (Etapa 1) |
| `World00_Character.bak` | `114` | [`database/schema_character.sql`](../database/schema_character.sql) (`109` tablas) + [`database/schema.sql`](../database/schema.sql) (`5` compartidas) + `2` S2S (`tActiveZoneTransfer`, `tWorldZoneNode`) | Cubierto (Etapa 2 / Fase 02A) |
| `World00_GameLog.bak` | **`2`** | `tGameLog` en [`database/schema.sql`](../database/schema.sql) + `TEMP_QuestComplete` en [`database/schema_character.sql`](../database/schema_character.sql) | **Validado y cableado en Fase 08** |
| `OperatorTool.bak` | `3` | `tMenuAuth`, `tOpLog`, `tOperator` en [`database/schema.sql`](../database/schema.sql) (herramienta externa `OpTool` excluida en runtime) | DDL presente / Runtime excluido |
| **TOTAL CLÚSTER OFICIAL** | **`176`** | **`xelion_online` (MySQL 8.x)** | **`176 / 176` (100%)** |

Detalle exacto de las **2 tablas oficiales** de `World00_GameLog.bak` según [`docs/TablasLegacy.md`](./TablasLegacy.md):

| # | Tabla Oficial (`World00_GameLog.bak`) | Columnas y Tipos Oficiales Exactos | Índices Requeridos para Escritura/Consulta Rápida |
| :---: | :--- | :--- | :--- |
| 1 | **`tGameLog`** | `nLogNo` (`bigint` PK `AUTO_INCREMENT`), `dDate` (`datetime`), `nType` (`int`), `nCharNo` (`int`), `sMap` (`varchar(16)`), `nMapX` (`int`), `nMapY` (`int`), `nMapZ` (`int`), `nTargetCharNo` (`int`), `nTargetID` (`int`), `nItemKey` (`bigint`), `nInt1` (`int`), `nInt2` (`int`), `nInt3` (`int`), `nBigint1` (`bigint`) | PK `(nLogNo)`, `idx_gamelog_date_type (dDate, nType)`, `idx_gamelog_char (nCharNo, dDate)`, `idx_gamelog_itemkey (nItemKey)` |
| 2 | **`TEMP_QuestComplete`** | `nCharNo` (`int`), `nQuestCompleteCnt` (`int`) | `idx_temp_quest_char (nCharNo)` |

### 2.2 Retiro Obligatorio de `database/zone_tables.sql` (86 Tablas Duplicadas)

En la **Fase 08** se debe eliminar el archivo [`database/zone_tables.sql`](../database/zone_tables.sql), el cual contiene las siguientes **86 tablas duplicadas** de un borrador antiguo ya superadas por `database/schema_character.sql` y `database/schema.sql`:

`tArenaChallenge`, `tArenaFBZ`, `tArenaGuild`, `tArenaPvP`, `tAttendance`, `tAttendanceReward`, `tAutoMon`, `tAutoPick`, `tCharacterChangeID`, `tCharacterChargedBuffer`, `tCharacterGameOption`, `tCharacterKeyMap`, `tCharacterLoading`, `tCharacterSingleOption`, `tCharacterState`, `tCharacterTitle`, `tChatBlock`, `tChatColor`, `tCollectionCard`, `tCollectionFavorites`, `tCollectionReward`, `tCollectionSerialCount`, `tDice`, `tDiceLog`, `tDicePayLog`, `tDiceResult`, `tExp`, `tFriend`, `tGMItemReward`, `tGMItemRewardLog`, `tGambleCoin`, `tGambleCoinExchangeLog`, `tGambleDice`, `tGambleDiceLog`, `tGambleDiceResult`, `tGambleJackpot`, `tGambleJackpotLog`, `tGambleSlotMachineLog`, `tGameOptionInit`, `tGuild`, `tGuildAcademy`, `tGuildAcademyGraduate`, `tGuildAcademyHistory`, `tGuildAcademyMember`, `tGuildAcademyRanking`, `tGuildAcademyReward`, `tGuildAcademyRewardType`, `tGuildEmblem`, `tGuildMember`, `tGuildToken`, `tGuildTournament`, `tGuildTournamentRegist`, `tHolyPromise`, `tHolyPromiseReward`, `tInGameGuide`, `tItemActionCoolTime`, `tItemEnchant`, `tItemMobList`, `tItemOptions`, `tItemStorageFindEmpty`, `tItemTitle`, `tKeyMapInit`, `tKingdomQuest`, `tLog_GuildMasterChange`, `tMarriage`, `tMasterPupils`, `tMinihouse`, `tMinihouseGate`, `tMinihouseGuest`, `tParty`, `tPartyMember`, `tPet`, `tPetNew`, `tPetSkill`, `tPrison`, `tQuest`, `tQuestTimes`, `tRestExp`, `tSkill`, `tSpamBlockTime`, `tSpamFilter`, `tStatCharLevel`, `tStatGameMoney`, `tStatisticsGameMoney`, `tTutorial`, `tUserMoney`.

### 2.3 Catálogo 100% Exhaustivo de los 27 Stored Procedures de `World00_GameLog` (`database/procedures_gamelog.sql`)

La unión matemática entre las **`23` invocaciones ODBC de `GameLog.exe`** ([`docs/Servicios/GameLog.md`](./Servicios/GameLog.md)) y los **`26` procedimientos recuperados de `World00_GameLog.bak`** ([`docs/sql/procedures/World00_GameLog_all_procedures.sql`](./sql/procedures/World00_GameLog_all_procedures.sql)) arroja exactamente **`27` Stored Procedures únicos** (`22` compartidos + `1` exclusivo de `GameLog.exe` + `4` de consulta administrativa en `.bak`).

> [!NOTE]
> **Verificación forense binaria de `usp_Login_setChannel` en `GameLog.exe` (`VA 0x41B690..0x41B720`, string `{CALL usp_Login_setChannel( ?, ?, ?, ? )}` en `0x43996C`):**
> El desensamblado de las 4 llamadas `SQLBindParameter` confirma exactamente:
> - Parámetro 1 (`IN`): `SQL_C_ULONG` / `SQL_INTEGER` (`[edi + 0x00]`) -> `in_nUserNo INT`
> - Parámetro 2 (`IN`): `SQL_C_ULONG` / `SQL_INTEGER` (`[edi + 0x04]`) -> `in_nChannelNo INT`
> - Parámetro 3 (`IN`): `SQL_C_CHAR` / `SQL_WVARCHAR(20)` (`[edi + 0x08]`, longitud `0x14 = 20`) -> `in_sChannelID VARCHAR(20)`
> - Parámetro 4 (`OUT`): `SQL_PARAM_OUTPUT` / `SQL_C_SLONG` (`[edi + 0x1C]`, inicializado en `0` y retornado en `edx:eax`) -> `OUT out_nResult INT`

Listado íntegro 1:1 de los **27 Stored Procedures** que deben compilarse en `database/procedures_gamelog.sql` con el patrón pool-safe (`proc_logic` + `SELECT` final incondicional):

| # | Stored Procedure Oficial | Origen | Firma Oficial de Parámetros (`.bak` / `GameLog.exe` ODBC) | Evento Auditado en `tGameLog` |
| :---: | :--- | :---: | :--- | :--- |
| 1 | `p_GameLog_Add` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int, @nBigint1 bigint` | Inserción completa de las 13 columnas de evento en `tGameLog` (`PROTO_NC_LOG_GAME_DATA_TYPE_0`) |
| 2 | `p_GameLog_Add_Type1` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int` | Tipo 1: Entrada/salida de mapa y coordenadas (`nCharNo, sMap, X, Y, Z`) |
| 3 | `p_GameLog_Add_Type2` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nInt1 int, @nInt2 int` | Tipo 2: Subida de nivel y progresión de estadísticas (`nInt1, nInt2`) |
| 4 | `p_GameLog_Add_Type3` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int, @nTargetID int, @nInt1 int` | Tipo 3: Transacciones de Cen / compra-venta con NPC (`nTargetCharNo, nTargetID, nInt1`) |
| 5 | `p_GameLog_Add_Type4` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nInt1 int, @nInt2 int, @nInt3 int` | Tipo 4: Eventos de instancia/script Lua y comercio (`nInt1, nInt2, nInt3`) |
| 6 | `p_GameLog_Add_Type5` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @nInt1 int, @nInt2 int` | Tipo 5: Eventos compactos de almacén/script Lua (`nCharNo, nInt1, nInt2`) |
| 7 | `p_GameLog_Add_Type6` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nInt1 int, @nBigint1 bigint` | Tipo 6: Obtención/descarte de ítems y montos Cen (`sMap, X, Y, Z, nInt1, nBigint1`) |
| 8 | `p_GameLog_Add_Type7` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int, @nInt1 int, @nInt2 int` | Tipo 7: Mejora, encantamiento y alteración de ítem (`nItemKey, nInt1, nBigint1`) |
| 9 | `p_GameLog_Add_Type8` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nInt1 int` | Tipo 8: Muerte de personaje PvE/PvP (`nTargetID, nBigint1`) |
| 10 | `p_GameLog_Add_Type9` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nItemKey bigint, @nInt1 int, @nInt2 int` | Tipo 9: Inicio, progreso y entrega de misiones (`nTargetID, nInt1, nBigint1`) |
| 11 | `p_GameLog_Add_TypeA` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @nTargetCharNo int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int` | Tipo A: Eventos de gremio con ítem/moneda (`nItemKey, nInt1, nBigint1`) |
| 12 | `p_GameLog_Add_TypeB` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int` | Tipo B: Almacén de gremio y transferencias completas (`nTargetCharNo..nBigint1`) |
| 13 | `p_GameLog_Add_TypeC` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int` | Tipo C: Kioscos / tiendas personales `Booth` (`nTargetCharNo, nItemKey..nBigint1`) |
| 14 | `p_GameLog_Add_TypeD` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @nItemKey bigint, @nInt1 int, @nInt2 int` | Tipo D: Producción, mezcla y desarme de ítems (`nTargetID, nItemKey, nInt1`) |
| 15 | `p_GameLog_Add_TypeE` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int` | Tipo E: Habilidades, empoderamiento y reinicio (`nTargetCharNo, nTargetID, nInt1, nInt2`) |
| 16 | `p_GameLog_Add_TypeF` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int, @nBigint1 bigint` | Tipo F: Operaciones en `MiniHouse` y mobiliario (`nTargetID, nItemKey..nBigint1`) |
| 17 | `p_GameLog_Add_TypeG` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nBigint1 bigint` | Tipo G: Mascotas (`Pet`) y monturas (`Mover`) (`nTargetCharNo, nTargetID, nInt1..nBigint1`) |
| 18 | `p_GameLog_Add_TypeH` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int` | Tipo H: `Kingdom Quest` y eventos comunitarios (`nTargetCharNo, nTargetID, nInt1, nInt2`) |
| 19 | `p_GameLog_Add_TypeI` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetCharNo int, @nTargetID int` | Tipo I: Casino, dados `Taisai` y `SlotMachine` (`sMap, X, Y, Z, nTargetID, nInt1`) |
| 20 | `p_GameLog_Add_TypeJ` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int` | Tipo J: Bodas, anillos y divorcios (`nTargetCharNo, nItemKey, nInt1, nInt2`) |
| 21 | `p_GameLog_Add_TypeK` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @sMap varchar(16), @nMapX int, @nMapY int, @nMapZ int, @nInt1 int, @nInt2 int, @nBigint1 bigint` | Tipo K: Academia de gremio y graduaciones (`nTargetCharNo, nInt1, nBigint1`) |
| 22 | `p_GameLog_Add_TypeL` | Ambos (`GameLog.exe` + `.bak`) | `@nLogNo bigint OUTPUT, @nType int, @nCharNo int, @nTargetCharNo int, @nTargetID int, @nItemKey bigint, @nInt1 int, @nInt2 int, @nInt3 int, @nBigint1 bigint` | Tipo L: Reportes, sanciones GM y recompensas (`nTargetCharNo, nTargetID, nItemKey, nInt1..nBigint1`) |
| 23 | `p_GameLog_Get` | Solo `World00_GameLog.bak` | `@WhereExpr nvarchar(256)` | Consulta administrativa filtrada sobre `tGameLog` (sentencia preparada segura en MySQL) |
| 24 | `p_GameLog_GetCount` | Solo `World00_GameLog.bak` | `@WhereExpr nvarchar(256), @Count int OUTPUT` | Conteo administrativo de registros filtrados en `tGameLog` |
| 25 | `p_GameLog_GetLog` | Solo `World00_GameLog.bak` | `@QueryExpr nvarchar(256)` | Ejecución de consulta parametrizada de auditoría sobre `tGameLog` |
| 26 | `p_GameLog_GetLogCount` | Solo `World00_GameLog.bak` | `@QueryExpr nvarchar(256), @Count int OUTPUT` | Ejecución de conteo con `OUT @Count INT` sobre `tGameLog` |
| 27 | `usp_Login_setChannel` | Solo `GameLog.exe` (`0x41B707`) | `@nUserNo int, @nChannelNo int, @sChannelID varchar(20), @nResult int OUTPUT` | Registro de canal de conexión y trazabilidad en `tGameLog` |

### 2.4 Catálogo 100% Exhaustivo de las 44 Estructuras S2S `PROTO_NC_LOG_*` (`docs/Servicios/GameLog.md`)

Todas las estructuras de auditoría definidas en `Zone.pdb` y documentadas en [`docs/Servicios/GameLog.md`](./Servicios/GameLog.md) que `pkg/protocol/` y `GameLogDispatcher` deben soportar con paridad exacta de tamaño:

| # | Estructura Oficial (`Zone.pdb` / `GameLog.md`) | `sizeof` Exacto | Procedimiento Destino en `GameLogRepository` |
| :---: | :--- | :---: | :--- |
| 1 | `PROTO_NC_LOG_GAME_DATA_TYPE_0` | `68 bytes` | `CALL p_GameLog_Add` |
| 2 | `PROTO_NC_LOG_GAME_DATA_TYPE_1` | `32 bytes` | `CALL p_GameLog_Add_Type1` |
| 3 | `PROTO_NC_LOG_GAME_DATA_TYPE_2` | `40 bytes` | `CALL p_GameLog_Add_Type2` |
| 4 | `PROTO_NC_LOG_GAME_DATA_TYPE_3` | `44 bytes` | `CALL p_GameLog_Add_Type3` |
| 5 | `PROTO_NC_LOG_GAME_DATA_TYPE_4` | `44 bytes` | `CALL p_GameLog_Add_Type4` |
| 6 | `PROTO_NC_LOG_GAME_DATA_TYPE_5` | `16 bytes` | `CALL p_GameLog_Add_Type5` |
| 7 | `PROTO_NC_LOG_GAME_DATA_TYPE_6` | `44 bytes` | `CALL p_GameLog_Add_Type6` |
| 8 | `PROTO_NC_LOG_GAME_DATA_TYPE_7` | `44 bytes` | `CALL p_GameLog_Add_Type7` |
| 9 | `PROTO_NC_LOG_GAME_DATA_TYPE_8` | `36 bytes` | `CALL p_GameLog_Add_Type8` |
| 10 | `PROTO_NC_LOG_GAME_DATA_TYPE_9` | `48 bytes` | `CALL p_GameLog_Add_Type9` |
| 11 | `PROTO_NC_LOG_GAME_DATA_TYPE_A` | `36 bytes` | `CALL p_GameLog_Add_TypeA` |
| 12 | `PROTO_NC_LOG_GAME_DATA_TYPE_B` | `60 bytes` | `CALL p_GameLog_Add_TypeB` |
| 13 | `PROTO_NC_LOG_GAME_DATA_TYPE_C` | `56 bytes` | `CALL p_GameLog_Add_TypeC` |
| 14 | `PROTO_NC_LOG_GAME_DATA_TYPE_D` | `24 bytes` | `CALL p_GameLog_Add_TypeD` |
| 15 | `PROTO_NC_LOG_GAME_DATA_TYPE_E` | `48 bytes` | `CALL p_GameLog_Add_TypeE` |
| 16 | `PROTO_NC_LOG_GAME_DATA_TYPE_F` | `64 bytes` | `CALL p_GameLog_Add_TypeF` |
| 17 | `PROTO_NC_LOG_GAME_DATA_TYPE_G` | `60 bytes` | `CALL p_GameLog_Add_TypeG` |
| 18 | `PROTO_NC_LOG_GAME_DATA_TYPE_H` | `36 bytes` | `CALL p_GameLog_Add_TypeH` |
| 19 | `PROTO_NC_LOG_GAME_DATA_TYPE_I` | `40 bytes` | `CALL p_GameLog_Add_TypeI` |
| 20 | `PROTO_NC_LOG_GAME_DATA_TYPE_J` | `52 bytes` | `CALL p_GameLog_Add_TypeJ` |
| 21 | `PROTO_NC_LOG_GAME_DATA_TYPE_K` | `48 bytes` | `CALL p_GameLog_Add_TypeK` |
| 22 | `PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO` | `32 bytes` | `CALL p_GameLog_Add_Type7` |
| 23 | `PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED` | `16 bytes` | `CALL p_GameLog_Add_Type7` |
| 24 | `PROTO_NC_LOG_GAME_ITEM_UPGRADE` | `43 bytes` | `CALL p_GameLog_Add_Type7` |
| 25 | `PROTO_NC_LOG_GAME_KQ_ENTER` | `30 bytes` | `CALL p_GameLog_Add_TypeH` |
| 26 | `PROTO_NC_LOG_GAME_KQ_LEAVE` | `31 bytes` | `CALL p_GameLog_Add_TypeH` |
| 27 | `PROTO_NC_LOG_GAME_LEVEL_UP` | `33 bytes` | `CALL p_GameLog_Add_Type2` |
| 28 | `PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING` | `24 bytes` | `CALL p_GameLog_Add_TypeF` |
| 29 | `PROTO_NC_LOG_GAME_MINIHOUSE_VISIT` | `28 bytes` | `CALL p_GameLog_Add_TypeF` |
| 30 | `PROTO_NC_LOG_GAME_MK2_DEAD` | `26 bytes` | `CALL p_GameLog_Add_Type8` |
| 31 | `PROTO_NC_LOG_GAME_MK2_START` | `28 bytes` | `CALL p_GameLog_Add_Type8` |
| 32 | `PROTO_NC_LOG_GAME_MKED` | `34 bytes` | `CALL p_GameLog_Add_Type8` |
| 33 | `PROTO_NC_LOG_GAME_MK_DROP_ITEM` | `34 bytes` | `CALL p_GameLog_Add_Type6` |
| 34 | `PROTO_NC_LOG_GAME_MK_GETEXPINFIELD` | `12 bytes` | `CALL p_GameLog_Add_Type2` |
| 35 | `PROTO_NC_LOG_GAME_PKED` | `33 bytes` | `CALL p_GameLog_Add_Type8` |
| 36 | `PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO` | `41 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |
| 37 | `PROTO_NC_LOG_REGENLOCATESAVE_CMD` | `38 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |
| 38 | `PROTO_NC_LOG_WEDDING_HALL_CANCEL` | `8 bytes` | `CALL p_GameLog_Add_TypeJ` |
| 39 | `PROTO_NC_LOG_WEDDING_HALL_RESERVE` | `8 bytes` | `CALL p_GameLog_Add_TypeJ` |
| 40 | `PROTO_NC_LOG_WEDDING_HALL_START` | `8 bytes` | `CALL p_GameLog_Add_TypeJ` |
| 41 | `PROTO_NC_SYSLOG_CHAR_DEATH` | `22 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |
| 42 | `PROTO_NC_SYSLOG_CHAR_QUEST_STARTED` | `13 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |
| 43 | `PROTO_NC_SYSLOG_CHAR_VICTORY` | `22 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |
| 44 | `PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION` | `11 bytes` | `CALL p_GameLog_Add_Type1` / `CALL p_GameLog_Add` |

---

## 3. Universo Completo de Datos Estáticos, Geometría Binaria y Runtime Lua de `ZoneService`

### 3.1 Las 199 Referencias `.shn` de `Zone.exe` (`194` Rutas en Disco = `180` Basenames Reales + `5` Opcionales = `185` Basenames Únicos)

El análisis binario de `Zone.exe` identifica **`199` referencias relativas a archivos `.shn`** (`185` nombres base únicos):
- **`194` referencias existen físicamente en `data/` y `data/View/`** (correspondientes a **`180` nombres base `.shn` únicos en disco**, dado que 14 referencias cargan o validan variantes entre raíz `data/` y `data/View/`), además de las **`49` tablas verificadas criptográficamente** por `shndatafilechecksum.cpp` (`ChecksumEngine`) en `Zone.exe` y `GameLog.exe`.
- Exactamente **`5` cadenas `.shn` en `Zone.exe` son referencias opcionales/legacy que tampoco existen en `../AstraServer/9Data/Shine/`**: `BMPIndex.shn`, `DiecGame.shn` (typo legacy reemplazado por `DiceGame.shn`), `MoveMain.shn` (reemplazado por `MoverMain.shn`), `MoverUpgrade.shn` (absorbido en `MoverSub.shn` / `UpgradeInfo.shn`) y `Quest.shn` (reemplazado por `QuestData.shn`). `DataStore` debe manejar esas `5` como opcionales sin error y exigir obligatoriamente las **`180` tablas `.shn` reales en disco** (`194` rutas).

Listado completo 1:1 de las **`180` tablas `.shn` reales en disco** requeridas en `pkg/server/zone/data_store.go`:

| # | Archivo `.shn` Oficial | Ubicación en `data/` | Estado Actual en `DataStore` |
| :---: | :--- | :--- | :---: |
| 1 | `AbState.shn` | `data/Abstate.shn` | **Cargar e indexar en Fase 10** |
| 2 | `AbStateSaveTypeInfo.shn` | `data/AbStateSaveTypeInfo.shn` | **Cargar e indexar en Fase 10** |
| 3 | `AbStateView.shn` | `data/AbStateView.shn`, `data/View/AbStateView.shn` | **Cargar e indexar en Fase 10** |
| 4 | `AccUpGrade.shn` | `data/AccUpgrade.shn` | **Cargar e indexar en Fase 10** |
| 5 | `AccUpGradeInfo.shn` | `data/AccUpGradeInfo.shn` | **Cargar e indexar en Fase 10** |
| 6 | `ActionEffectAbState.shn` | `data/ActionEffectAbState.shn` | **Cargar e indexar en Fase 10** |
| 7 | `ActionEffectItem.shn` | `data/ActionEffectItem.shn` | **Cargar e indexar en Fase 10** |
| 8 | `ActionRangeFactor.shn` | `data/ActionRangeFactor.shn` | **Cargar e indexar en Fase 10** |
| 9 | `ActionViewInfo.shn` | `data/ActionViewInfo.shn`, `data/View/ActionViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 10 | `ActiveSkill.shn` | `data/ActiveSkill.shn` | Cargado hoy (`9/180`) |
| 11 | `ActiveSkillGroup.shn` | `data/ActiveSkillGroup.shn` | **Cargar e indexar en Fase 10** |
| 12 | `ActiveSkillInfoServer.shn` | `data/ActiveSkillInfoServer.shn` | **Cargar e indexar en Fase 10** |
| 13 | `ActiveSkillView.shn` | `data/ActiveSkillView.shn`, `data/View/ActiveSkillView.shn` | **Cargar e indexar en Fase 10** |
| 14 | `AdminLvSet.shn` | `data/AdminLvSet.shn` | **Cargar e indexar en Fase 10** |
| 15 | `AnnounceData.shn` | `data/AnnounceData.shn` | **Cargar e indexar en Fase 10** |
| 16 | `AreaSkill.shn` | `data/AreaSkill.shn` | **Cargar e indexar en Fase 10** |
| 17 | `BadNameFilter.shn` | `data/BadNameFilter.shn` | **Cargar e indexar en Fase 10** |
| 18 | `BelongDice.shn` | `data/BelongDice.shn` | **Cargar e indexar en Fase 10** |
| 19 | `BelongTypeInfo.shn` | `data/BelongTypeInfo.shn` | **Cargar e indexar en Fase 10** |
| 20 | `BMP.shn` | `data/BMP.shn` | **Cargar e indexar en Fase 10** |
| 21 | `BRAccUpGrade.shn` | `data/BRAccUpgrade.shn` | **Cargar e indexar en Fase 10** |
| 22 | `BRAccUpgradeInfo.shn` | `data/BRAccUpgradeInfo.shn` | **Cargar e indexar en Fase 10** |
| 23 | `CharacterTitleData.shn` | `data/CharacterTitleData.shn` | **Cargar e indexar en Fase 10** |
| 24 | `CharacterTitleStateServer.shn` | `data/CharacterTitleStateServer.shn` | **Cargar e indexar en Fase 10** |
| 25 | `CharacterTitleStateView.shn` | `data/CharacterTitleStateView.shn`, `data/View/CharacterTitleStateView.shn` | **Cargar e indexar en Fase 10** |
| 26 | `ChargedDeletableBuff.shn` | `data/ChargedDeletableBuff.shn` | **Cargar e indexar en Fase 10** |
| 27 | `ChargedEffect.shn` | `data/ChargedEffect.shn` | **Cargar e indexar en Fase 10** |
| 28 | `ChatColor.shn` | `data/ChatColor.shn` | **Cargar e indexar en Fase 10** |
| 29 | `ClassName.shn` | `data/ClassName.shn` | Cargado hoy (`9/180`) |
| 30 | `CollectCard.shn` | `data/CollectCard.shn` | **Cargar e indexar en Fase 10** |
| 31 | `CollectCardDropRate.shn` | `data/CollectCardDropRate.shn` | **Cargar e indexar en Fase 10** |
| 32 | `CollectCardMobGroup.shn` | `data/CollectCardMobGroup.shn` | **Cargar e indexar en Fase 10** |
| 33 | `CollectCardReward.shn` | `data/CollectCardReward.shn` | **Cargar e indexar en Fase 10** |
| 34 | `CollectCardStarRate.shn` | `data/CollectCardStarRate.shn` | **Cargar e indexar en Fase 10** |
| 35 | `CollectCardTitle.shn` | `data/CollectCardTitle.shn` | **Cargar e indexar en Fase 10** |
| 36 | `CollectCardView.shn` | `data/CollectCardView.shn`, `data/View/CollectCardView.shn` | **Cargar e indexar en Fase 10** |
| 37 | `DamageLvGapEVP.shn` | `data/DamageLvGapEVP.shn` | **Cargar e indexar en Fase 10** |
| 38 | `DamageLvGapPVE.shn` | `data/DamageLvGapPVE.shn` | **Cargar e indexar en Fase 10** |
| 39 | `DamageLvGapPVP.shn` | `data/DamageLvGapPVP.shn` | **Cargar e indexar en Fase 10** |
| 40 | `DiceDividind.shn` | `data/DiceDividind.shn` | **Cargar e indexar en Fase 10** |
| 41 | `DiceGame.shn` | `data/DiceGame.shn` | **Cargar e indexar en Fase 10** |
| 42 | `DiceRate.shn` | `data/DiceRate.shn` | **Cargar e indexar en Fase 10** |
| 43 | `EffectViewInfo.shn` | `data/EffectViewInfo.shn`, `data/View/EffectViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 44 | `EnchantSocketRate.shn` | `data/EnchantSocketRate.shn` | **Cargar e indexar en Fase 10** |
| 45 | `FieldLvCondition.shn` | `data/FieldLvCondition.shn` | **Cargar e indexar en Fase 10** |
| 46 | `FriendPointReward.shn` | `data/FriendPointReward.shn` | **Cargar e indexar en Fase 10** |
| 47 | `Gather.shn` | `data/Gather.shn` | **Cargar e indexar en Fase 10** |
| 48 | `GBBanTime.shn` | `data/GBBanTime.shn` | **Cargar e indexar en Fase 10** |
| 49 | `GBDiceDividind.shn` | `data/GBDiceDividind.shn` | **Cargar e indexar en Fase 10** |
| 50 | `GBDiceGame.shn` | `data/GBDiceGame.shn` | **Cargar e indexar en Fase 10** |
| 51 | `GBDiceRate.shn` | `data/GBDiceRate.shn` | **Cargar e indexar en Fase 10** |
| 52 | `GBEventCode.shn` | `data/GBEventCode.shn` | **Cargar e indexar en Fase 10** |
| 53 | `GBExchangeMaxCoin.shn` | `data/GBExchangeMaxCoin.shn` | **Cargar e indexar en Fase 10** |
| 54 | `GBHouse.shn` | `data/GBHouse.shn` | **Cargar e indexar en Fase 10** |
| 55 | `GBJoinGameMember.shn` | `data/GBJoinGameMember.shn` | **Cargar e indexar en Fase 10** |
| 56 | `GBReward.shn` | `data/GBReward.shn` | **Cargar e indexar en Fase 10** |
| 57 | `GBSMAll.shn` | `data/GBSMAll.shn` | **Cargar e indexar en Fase 10** |
| 58 | `GBSMBetCoin.shn` | `data/GBSMBetCoin.shn` | **Cargar e indexar en Fase 10** |
| 59 | `GBSMCardRate.shn` | `data/GBSMCardRate.shn` | **Cargar e indexar en Fase 10** |
| 60 | `GBSMCenter.shn` | `data/GBSMCenter.shn` | **Cargar e indexar en Fase 10** |
| 61 | `GBSMGroup.shn` | `data/GBSMGroup.shn` | **Cargar e indexar en Fase 10** |
| 62 | `GBSMJPRate.shn` | `data/GBSMJPRate.shn` | **Cargar e indexar en Fase 10** |
| 63 | `GBSMLine.shn` | `data/GBSMLine.shn` | **Cargar e indexar en Fase 10** |
| 64 | `GBSMNPC.shn` | `data/GBSMNPC.shn` | **Cargar e indexar en Fase 10** |
| 65 | `GBTaxRate.shn` | `data/GBTaxRate.shn` | **Cargar e indexar en Fase 10** |
| 66 | `GradeItemOption.shn` | `data/GradeItemOption.shn` | **Cargar e indexar en Fase 10** |
| 67 | `GroupAbState.shn` | `data/GroupAbState.shn` | **Cargar e indexar en Fase 10** |
| 68 | `GTIBreedSubject.shn` | `data/GTIBreedSubject.shn` | **Cargar e indexar en Fase 10** |
| 69 | `GTIGetRate.shn` | `data/GTIGetRate.shn` | **Cargar e indexar en Fase 10** |
| 70 | `GTIGetRateGap.shn` | `data/GTIGetRateGap.shn` | **Cargar e indexar en Fase 10** |
| 71 | `GTIServer.shn` | `data/GTIServer.shn` | **Cargar e indexar en Fase 10** |
| 72 | `GTIView.shn` | `data/GTIView.shn`, `data/View/GTIView.shn` | **Cargar e indexar en Fase 10** |
| 73 | `GTWinScore.shn` | `data/GTWinScore.shn` | **Cargar e indexar en Fase 10** |
| 74 | `GuildTournamentLvGap.shn` | `data/GuildTournamentLvGap.shn` | **Cargar e indexar en Fase 10** |
| 75 | `GuildTournamentMasterBuff.shn` | `data/GuildTournamentMasterBuff.shn` | **Cargar e indexar en Fase 10** |
| 76 | `GuildTournamentOccupy.shn` | `data/GuildTournamentOccupy.shn` | **Cargar e indexar en Fase 10** |
| 77 | `GuildTournamentScore.shn` | `data/GuildTournamentScore.shn` | **Cargar e indexar en Fase 10** |
| 78 | `GuildTournamentSkill.shn` | `data/GuildTournamentSkill.shn` | **Cargar e indexar en Fase 10** |
| 79 | `ItemAction.shn` | `data/ItemAction.shn` | **Cargar e indexar en Fase 10** |
| 80 | `ItemActionCondition.shn` | `data/ItemActionCondition.shn` | **Cargar e indexar en Fase 10** |
| 81 | `ItemActionEffect.shn` | `data/ItemActionEffect.shn` | **Cargar e indexar en Fase 10** |
| 82 | `ItemDismantle.shn` | `data/ItemDismantle.shn` | **Cargar e indexar en Fase 10** |
| 83 | `ItemDropLog.shn` | `data/ItemDropLog.shn` | **Cargar e indexar en Fase 10** |
| 84 | `ItemInfo.shn` | `data/ItemInfo.shn` | Cargado hoy (`9/180`) |
| 85 | `ItemInfoServer.shn` | `data/ItemInfoServer.shn` | Cargado hoy (`9/180`) |
| 86 | `ItemInvenDel.shn` | `data/ItemInvenDel.shn` | **Cargar e indexar en Fase 10** |
| 87 | `ItemMerchantInfo.shn` | `data/ItemMerchantInfo.shn` | **Cargar e indexar en Fase 10** |
| 88 | `ItemMix.shn` | `data/ItemMix.shn` | **Cargar e indexar en Fase 10** |
| 89 | `ItemMoney.shn` | `data/ItemMoney.shn` | **Cargar e indexar en Fase 10** |
| 90 | `ItemOptions.shn` | `data/ItemOptions.shn` | **Cargar e indexar en Fase 10** |
| 91 | `ItemServerEquipTypeInfo.shn` | `data/ItemServerEquipTypeInfo.shn` | **Cargar e indexar en Fase 10** |
| 92 | `ItemShop.shn` | `data/ItemShop.shn` | **Cargar e indexar en Fase 10** |
| 93 | `ItemShopView.shn` | `data/ItemShopView.shn`, `data/View/ItemShopView.shn` | **Cargar e indexar en Fase 10** |
| 94 | `ItemSort.shn` | `data/ItemSort.shn` | Cargado hoy (`9/180`) |
| 95 | `ItemUpgrade.shn` | `data/ItemUpgrade.shn` | **Cargar e indexar en Fase 10** |
| 96 | `ItemUseEffect.shn` | `data/ItemUseEffect.shn` | Cargado hoy (`9/180`) |
| 97 | `ItemViewEquipTypeInfo.shn` | `data/ItemViewEquipTypeInfo.shn`, `data/View/ItemViewEquipTypeInfo.shn` | **Cargar e indexar en Fase 10** |
| 98 | `ItemViewInfo.shn` | `data/ItemViewInfo.shn`, `data/View/ItemViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 99 | `JobEquipInfo.shn` | `data/JobEquipInfo.shn` | **Cargar e indexar en Fase 10** |
| 100 | `KingdomQuestRew.shn` | `data/KingdomQuestRew.shn` | **Cargar e indexar en Fase 10** |
| 101 | `KQItem.shn` | `data/KQItem.shn` | **Cargar e indexar en Fase 10** |
| 102 | `LCGroupRate.shn` | `data/LCGroupRate.shn` | **Cargar e indexar en Fase 10** |
| 103 | `LCReward.shn` | `data/LCReward.shn` | **Cargar e indexar en Fase 10** |
| 104 | `MapBuff.shn` | `data/MapBuff.shn` | **Cargar e indexar en Fase 10** |
| 105 | `MapInfo.shn` | `data/MapInfo.shn` | **Cargar e indexar en Fase 10** |
| 106 | `MapLinkPoint.shn` | `data/MapLinkPoint.shn` | **Cargar e indexar en Fase 10** |
| 107 | `MapViewInfo.shn` | `data/MapViewInfo.shn`, `data/View/MapViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 108 | `MapWayPoint.shn` | `data/MapWayPoint.shn` | **Cargar e indexar en Fase 10** |
| 109 | `MarketSearchInfo.shn` | `data/MarketSearchInfo.shn` | **Cargar e indexar en Fase 10** |
| 110 | `MiniHouse.shn` | `data/MiniHouse.shn` | **Cargar e indexar en Fase 10** |
| 111 | `MiniHouseEndure.shn` | `data/MiniHouseEndure.shn` | **Cargar e indexar en Fase 10** |
| 112 | `MiniHouseFurniture.shn` | `data/MiniHouseFurniture.shn` | **Cargar e indexar en Fase 10** |
| 113 | `MiniHouseFurnitureObjEffect.shn` | `data/MiniHouseFurnitureObjEffect.shn` | **Cargar e indexar en Fase 10** |
| 114 | `MiniHouseObjAni.shn` | `data/MiniHouseObjAni.shn` | **Cargar e indexar en Fase 10** |
| 115 | `MinimonAutoUseItem.shn` | `data/MinimonAutoUseItem.shn` | **Cargar e indexar en Fase 10** |
| 116 | `MinimonInfo.shn` | `data/MinimonInfo.shn` | **Cargar e indexar en Fase 10** |
| 117 | `MobAbStateDropSetting.shn` | `data/MobAbStateDropSetting.shn` | **Cargar e indexar en Fase 10** |
| 118 | `MobAutoAction.shn` | `data/MobAutoAction.shn` | **Cargar e indexar en Fase 10** |
| 119 | `MobConditionServer.shn` | `data/MobConditionServer.shn` | **Cargar e indexar en Fase 10** |
| 120 | `MobInfo.shn` | `data/MobInfo.shn` | Cargado hoy (`9/180`) |
| 121 | `MobInfoServer.shn` | `data/MobInfoServer.shn` | **Cargar e indexar en Fase 10** |
| 122 | `MobKillAble.shn` | `data/MobKillAble.shn` | **Cargar e indexar en Fase 10** |
| 123 | `MobKillAnnounce.shn` | `data/MobKillAnnounce.shn` | **Cargar e indexar en Fase 10** |
| 124 | `MobKillLog.shn` | `data/MobKillLog.shn` | **Cargar e indexar en Fase 10** |
| 125 | `MobLifeTime.shn` | `data/MobLifeTime.shn` | **Cargar e indexar en Fase 10** |
| 126 | `MobRegenAni.shn` | `data/MobRegenAni.shn` | **Cargar e indexar en Fase 10** |
| 127 | `MobResist.shn` | `data/MobResist.shn` | **Cargar e indexar en Fase 10** |
| 128 | `MobSpecies.shn` | `data/MobSpecies.shn` | **Cargar e indexar en Fase 10** |
| 129 | `MobViewInfo.shn` | `data/MobViewInfo.shn`, `data/View/MobViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 130 | `MobWeapon.shn` | `data/MobWeapon.shn` | **Cargar e indexar en Fase 10** |
| 131 | `MoverAbility.shn` | `data/MoverAbility.shn` | **Cargar e indexar en Fase 10** |
| 132 | `MoverHG.shn` | `data/MoverHG.shn` | **Cargar e indexar en Fase 10** |
| 133 | `MoverItem.shn` | `data/MoverItem.shn` | **Cargar e indexar en Fase 10** |
| 134 | `MoverMain.shn` | `data/MoverMain.shn` | **Cargar e indexar en Fase 10** |
| 135 | `MoverUpgradeEffect.shn` | `data/MoverUpgradeEffect.shn` | **Cargar e indexar en Fase 10** |
| 136 | `MoverUseSkill.shn` | `data/MoverUseSkill.shn` | **Cargar e indexar en Fase 10** |
| 137 | `MultiHitType.shn` | `data/MultiHitType.shn` | **Cargar e indexar en Fase 10** |
| 138 | `MysteryVaultServer.shn` | `data/MysteryVaultServer.shn` | **Cargar e indexar en Fase 10** |
| 139 | `NpcSchedule.shn` | `data/NpcSchedule.shn` | **Cargar e indexar en Fase 10** |
| 140 | `NPCViewInfo.shn` | `data/NPCViewInfo.shn`, `data/View/NPCViewInfo.shn` | **Cargar e indexar en Fase 10** |
| 141 | `PartyBonusByLvDiff.shn` | `data/PartyBonusByLvDiff.shn` | **Cargar e indexar en Fase 10** |
| 142 | `PartyBonusByMember.shn` | `data/PartyBonusByMember.shn` | **Cargar e indexar en Fase 10** |
| 143 | `PartyBonusLimit.shn` | `data/PartyBonusLimit.shn` | **Cargar e indexar en Fase 10** |
| 144 | `PassiveSkill.shn` | `data/PassiveSkill.shn` | Cargado hoy (`9/180`) |
| 145 | `PassiveSkillView.shn` | `data/PassiveSkillView.shn`, `data/View/PassiveSkillView.shn` | **Cargar e indexar en Fase 10** |
| 146 | `Produce.shn` | `data/Produce.shn` | **Cargar e indexar en Fase 10** |
| 147 | `ProduceView.shn` | `data/ProduceView.shn`, `data/View/ProduceView.shn` | **Cargar e indexar en Fase 10** |
| 148 | `PSkillSetAbstate.shn` | `data/PSkillSetAbstate.shn` | **Cargar e indexar en Fase 10** |
| 149 | `PupCase.shn` | `data/PupCase.shn` | **Cargar e indexar en Fase 10** |
| 150 | `PupCaseDesc.shn` | `data/PupCaseDesc.shn` | **Cargar e indexar en Fase 10** |
| 151 | `PupFactorCondition.shn` | `data/PupFactorCondition.shn` | **Cargar e indexar en Fase 10** |
| 152 | `PupMain.shn` | `data/PupMain.shn` | **Cargar e indexar en Fase 10** |
| 153 | `PupMind.shn` | `data/PupMind.shn` | **Cargar e indexar en Fase 10** |
| 154 | `PupPriority.shn` | `data/PupPriority.shn` | **Cargar e indexar en Fase 10** |
| 155 | `PupServer.shn` | `data/PupServer.shn` | **Cargar e indexar en Fase 10** |
| 156 | `QuestData.shn` | `data/QuestData.shn` | **Cargar e indexar en Fase 10** |
| 157 | `QuestSpecies.shn` | `data/QuestSpecies.shn` | **Cargar e indexar en Fase 10** |
| 158 | `RandomOption.shn` | `data/RandomOption.shn` | **Cargar e indexar en Fase 10** |
| 159 | `RandomOptionCount.shn` | `data/RandomOptionCount.shn` | **Cargar e indexar en Fase 10** |
| 160 | `RareMoverEachRate.shn` | `data/RareMoverEachRate.shn` | **Cargar e indexar en Fase 10** |
| 161 | `RareMoverRate.shn` | `data/RareMoverRate.shn` | **Cargar e indexar en Fase 10** |
| 162 | `RareMoverSubRate.shn` | `data/RareMoverSubRate.shn` | **Cargar e indexar en Fase 10** |
| 163 | `ReactionType.shn` | `data/ReactionType.shn` | **Cargar e indexar en Fase 10** |
| 164 | `Riding.shn` | `data/Riding.shn` | **Cargar e indexar en Fase 10** |
| 165 | `SetEffect.shn` | `data/SetEffect.shn` | **Cargar e indexar en Fase 10** |
| 166 | `ShineReward.shn` | `data/ShineReward.shn` | **Cargar e indexar en Fase 10** |
| 167 | `SingleData.shn` | `data/SingleData.shn` | **Cargar e indexar en Fase 10** |
| 168 | `SpamerPenalty.shn` | `data/SpamerPenalty.shn` | **Cargar e indexar en Fase 10** |
| 169 | `SpamerPenaltyRule.shn` | `data/SpamerPenaltyRule.shn` | **Cargar e indexar en Fase 10** |
| 170 | `StateField.shn` | `data/StateField.shn` | **Cargar e indexar en Fase 10** |
| 171 | `StateItem.shn` | `data/StateItem.shn` | **Cargar e indexar en Fase 10** |
| 172 | `StateMob.shn` | `data/StateMob.shn` | **Cargar e indexar en Fase 10** |
| 173 | `SubAbState.shn` | `data/SubAbstate.shn` | **Cargar e indexar en Fase 10** |
| 174 | `TermExtendMatch.shn` | `data/TermExtendMatch.shn` | **Cargar e indexar en Fase 10** |
| 175 | `ToggleSkill.shn` | `data/ToggleSkill.shn` | **Cargar e indexar en Fase 10** |
| 176 | `TownPortal.shn` | `data/TownPortal.shn` | **Cargar e indexar en Fase 10** |
| 177 | `UpGradeInfo.shn` | `data/UpgradeInfo.shn` | **Cargar e indexar en Fase 10** |
| 178 | `UseClassTypeInfo.shn` | `data/UseClassTypeInfo.shn` | Cargado hoy (`9/180`) |
| 179 | `WeaponAttrib.shn` | `data/WeaponAttrib.shn` | **Cargar e indexar en Fase 10** |
| 180 | `WeaponTitleData.shn` | `data/WeaponTitleData.shn` | **Cargar e indexar en Fase 10** |

Registro explícito de las **5 referencias `.shn` opcionales** presentes como strings en `Zone.exe` pero ausentes tanto en `AstraServer` como en `data/` (`180 + 5 = 185` basenames únicos en `Zone.exe`):
1. `BMPIndex.shn` (opcional: las máscaras `.bmp` se resuelven directamente por nombre desde `ActiveSkill.shn` / `data/AreaBMP/`).
2. `DiecGame.shn` (errata tipográfica legacy reemplazada por `DiceGame.shn`, incluida en la tabla superior).
3. `MoveMain.shn` (reemplazado por `MoverMain.shn`, incluida en la tabla superior).
4. `MoverUpgrade.shn` (absorbido por `MoverSub.shn` y `UpgradeInfo.shn`).
5. `Quest.shn` (reemplazado por `QuestData.shn`, incluida en la tabla superior).

### 3.2 Los 31 Archivos `.txt` Referenciados por `Zone.exe` y su Mapeo 1:1 en `XelionCore`

| # | Archivo `.txt` en `Zone.exe` | Equivalente Modernizado / Fuente Oficial en `XelionCore` | Subsistema Go |
| :---: | :--- | :--- | :--- |
| 1 | `02d.txt` | Patrón de formato de logs diarios de sesión (`%02d.txt`) -> `pkg/logger` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 2 | `_ServerGroup.txt` | Topología de nodos y puertos -> `.env` + tabla `tWorldZoneNode` (`xelion_online`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 3 | `CallStack.txt` | Volcado de diagnóstico ante crash C++ -> `runtime/debug.Stack()` + `pkg/logger` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 4 | `ChrCommon.txt` | `data/World/Combat/Common.json`, `ExpTable.json`, `FreeStatTable.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 5 | `DamageByAngle.txt` | `data/World/Combat/DamageByAngle.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 6 | `DamageBySoul.txt` | `data/World/Combat/Common.json` (`soulStoneMultiplier`) + `SingleData.shn` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 7 | `Dbg.txt` | Trazas de depuración -> `pkg/logger` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 8 | `Field.txt` | `data/World/Map/FieldList.json`, `InstanceDungeon.json`, `GuildTournament.json`, `TutorialMap.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 9 | `FunctionPerformance.txt` | Telemetría de ticks de función -> métricas de latencia del Game Loop en `pkg/server/zone` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 10 | `ItemDropGroup.txt` | `data/World/Drop/ItemDropGroup.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 11 | `ItemDropTable.txt` | `data/World/Drop/ItemDropTable.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 12 | `ItemUseFunction.txt` | `data/World/Item/ItemUseFunction.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 13 | `MapLoad.txt` | `data/World/Map/FieldList.json` (`zoneNumber` por instancia `Zone00`..`Zone04`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 14 | `Message.txt` | `data/Script/*.json` + `TextData.shn` / `ErrorCodeTable.shn` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 15 | `MiscDataTable.txt` | `data/World/Misc/SkillMechanics.json` + `StealthMatrix.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 16 | `MobChat.txt` | `data/World/Misc/MobChat.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 17 | `NPC.txt` | `data/World/NPC/Maps/<Map>.json` + `data/World/NPC/LinkTable.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 18 | `NPCAction.txt` | `data/World/NPC/NPCAction.json` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 19 | `PartBlock.txt` | Subcapas de bloqueo parcial de mapa -> `data/BlockInfo/*.sbi` (`pkg/data/blockinfo`) | **Nuevo `pkg/data/blockinfo` (Fase 10)** |
| 20 | `PineScript.txt` | `data/Script/*.json` + `data/LuaScript/` (`pkg/data/script` + `pkg/script/lua`) | **`pkg/data/script` + `pkg/script/lua` (Fase 10)** |
| 21 | `QuestParser.txt` | `data/World/Quest/QuestParser.json` (`pkg/data/questparser`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 22 | `RecallCoord.txt` | `data/World/Misc/RecallCoord.json` (`pkg/data/misc`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 23 | `s.txt` | Patrón de formato `%s.txt` de carga de tablas por nombre (`<Class>.txt`, `<Mob>.txt`, `<NPC>.txt`) -> JSONs modulares en `data/` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 24 | `ServerGroup.txt` | Topología del clúster -> `.env` + `tWorldServer` + `tWorldZoneNode` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 25 | `Size.txt` | Dimensiones de estructuras/tablas -> validadas en `pkg/audit` contra PDB `sizeof` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 26 | `StorageError.txt` | Auditoría de errores de almacenamiento -> `pkg/logger` + `tSyslog` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 27 | `SubAbStateClass.txt` | `data/AbState/*.json` (`pkg/data/abstate`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 28 | `SubLayerInteract.txt` | `data/BlockInfo/*.sbi` (`SubLayerInteract` en `pkg/data/blockinfo`) | **Nuevo `pkg/data/blockinfo` (Fase 10)** |
| 29 | `SubLayerInteractLoader.txt` | `data/BlockInfo/*.sbi` (`SubLayerInteractLoader` en `pkg/data/blockinfo`) | **Nuevo `pkg/data/blockinfo` (Fase 10)** |
| 30 | `TreasureReward.txt` | `data/World/Item/TreasureReward.json` (`pkg/data/item`) | Cargador activo en `pkg/data/*` / `pkg/server/zone` |
| 31 | `ZoneServerInfo.txt` | Configuración de instancia `Zone00`..`Zone04` -> `pkg/config` + `tWorldZoneNode` | Cargador activo en `pkg/data/*` / `pkg/server/zone` |

### 3.3 Geometría Binaria de Mapas y Colisiones (`pkg/data/blockinfo` y `pkg/data/areabmp` — Fase 10)

`Zone.exe` utiliza los cargadores C++ `MapBlock::MapBlockInformation`, `MapBlock::MapDoorArray`, `MapBlock::MapDoorBlockElement`, `SubLayerInteractLoader`, `AreaInfoData` y `CFsAreaBMP` documentados en los símbolos de `Zone.pdb`. En la **Fase 10** se deben crear dos paquetes oficiales nuevos en `pkg/data/`:

| Paquete Go Nuevo | Extensión en `data/` | Cantidad en Disco | Formato Binario Verificado Forensemente | Responsabilidad en `ZoneService` |
| :--- | :---: | :---: | :--- | :--- |
| **`pkg/data/blockinfo`** | `.shbd` | `158` | Header `8 bytes`: `uint32 rowBytes` (`widthTiles / 8`), `uint32 heightTiles`, seguido exactamente de `rowBytes * heightTiles` bytes (`1 bit` por celda: `1` = bloqueado, `0` = transitable; ej. `Adl.shbd`: `950 * 7600 + 8 = 7.220.008 bytes`). | Validación `IsMoveBlocked(x, y)`, `IsAttackBlocked(x1, y1, x2, y2)` (Bresenham LoS) y `FindNearestWalkable(x, y)` para jugadores, mascotas y monstruos. |
| **`pkg/data/blockinfo`** | `.sbi` | `33` | Header `uint32 count` seguido de registros `MapDoorBlockElement` (`char doorName[32]` + rectángulo y bitmap local de bloqueo; ej. `AdlF.sbi`: `2` puertas `DoorBlock01..02`, `3.220 bytes`). | Apertura/cierre dinámico de puertas e instancias (`cDoorBuild`, `cDoorAction`, `mda_CloseAllDoor`) en KQ e Instance Dungeons. |
| **`pkg/data/blockinfo`** | `.aid` | `23` | Header `uint32 count` seguido de entradas `AreaInfoData` (`AE_Circle` / `AE_Square`, ej. `Adl.aid`: `MH_Zone1`, `60 bytes`). | Detección de zonas especiales (`MiniHouse`, `ScenarioArea`, zonas seguras, triggers `cIsInArea`). |
| **`pkg/data/blockinfo`** | `.bdt` | `76` | Tablas binarias de detalle de bloques y objetos de puerta/puente (`BFGate01.bdt`, etc.). | Geometría auxiliar de colisión de estructuras interactivas. |
| **`pkg/data/blockinfo`** | `.shab` | `8` | Header `8 bytes` (`uint32 rowBytes, uint32 heightTiles`) + matriz `rowBytes * heightTiles` de atributos de arena (`KDArena.shab`: `128 * 1024 + 8 = 131.080 bytes`). | Atributos zonales de arenas de combate y torneos. |
| **`pkg/data/blockinfo`** | `.shmd` / `.shad` | `2` / `1` | Descriptores de escena (`Rou.shmd` texto `shmd0_5`) y atributo extendido (`UrgDragon.shad`: `256 * 2048 + 8 = 524.296 bytes`). | Altura/capas especiales de encuentro de jefe (`UrgDragon`). |
| **`pkg/data/areabmp`** | `.bmp` | `44` | Imágenes Windows BMP indexadas de `8 bpp` (`BM` header, ej. `ArkAstanica_Skill01_W.bmp`) cargadas por `CFsAreaBMP` en `Zone.exe`. | Máscaras espaciales rotables de impacto AoE para habilidades especiales de monstruos/jefes (`ActiveSkill.shn`). |

### 3.4 Catálogo 100% Exhaustivo de las 274 Funciones C-API `lua_State*` de `Zone.pdb` (`pkg/script/lua/`)

El flujo de símbolos públicos de `Zone.pdb` contiene exactamente **`274` funciones con firma `int __cdecl <name>(lua_State *)`** (`?*@@YAHPAUlua_State@@@Z`):
- **`246` funciones `c*` de Gameplay/IA/KQ/Instancias/Quests/Pets**.
- **`22` callbacks de configuración de arranque `ZoneServer_zs_*`**.
- **`6` utilidades base del motor** (`Parameter_RunSpeed`, `Parameter_WalkSpeed`, `cpp_CurrentTime`, `cpp_Print`, `cpp_RegCheck`, `tonumber`), más las funciones auxiliares invocadas por scripts (`cDistanceSquar`, `cRandomInt`, `CallNamedRoutine`).

Listado completo 1:1 de las **`274` funciones C-API `lua_State*` oficiales de `Zone.pdb`** que deben quedar registradas y operativas en `pkg/script/lua/`:

| # | Símbolo C-API en `Zone.pdb` | Categoría Funcional | # | Símbolo C-API en `Zone.pdb` | Categoría Funcional |
| :---: | :--- | :--- | :---: | :--- | :--- |
| 1 | `Parameter_RunSpeed` | Objetos / Movimiento / Mensajes | 138 | `cMobDialog_FileName` | IA Monstruos / NPCs / Spawns |
| 2 | `Parameter_WalkSpeed` | Objetos / Movimiento / Mensajes | 139 | `cMobDialog_Obj` | IA Monstruos / NPCs / Spawns |
| 3 | `ZoneServer_zs_service_acceptor` | Configuración `ZoneServer` | 140 | `cMobDialog_Range` | IA Monstruos / NPCs / Spawns |
| 4 | `ZoneServer_zs_service_certification` | Configuración `ZoneServer` | 141 | `cMobDialog_String` | IA Monstruos / NPCs / Spawns |
| 5 | `ZoneServer_zs_service_iocp` | Configuración `ZoneServer` | 142 | `cMobIDFind` | IA Monstruos / NPCs / Spawns |
| 6 | `ZoneServer_zs_service_mainthread` | Configuración `ZoneServer` | 143 | `cMobRegen_Circle` | IA Monstruos / NPCs / Spawns |
| 7 | `ZoneServer_zs_service_multisocket` | Configuración `ZoneServer` | 144 | `cMobRegen_Obj` | IA Monstruos / NPCs / Spawns |
| 8 | `ZoneServer_zs_service_stackchecker` | Configuración `ZoneServer` | 145 | `cMobRegen_Rectangle` | IA Monstruos / NPCs / Spawns |
| 9 | `ZoneServer_zs_start_black` | Configuración `ZoneServer` | 146 | `cMobRegen_XY` | IA Monstruos / NPCs / Spawns |
| 10 | `ZoneServer_zs_start_earth` | Configuración `ZoneServer` | 147 | `cMobRegen_XY_Layer` | Mapa / Puertas `.sbi` / Capas |
| 11 | `ZoneServer_zs_start_full` | Configuración `ZoneServer` | 148 | `cMobShout` | IA Monstruos / NPCs / Spawns |
| 12 | `ZoneServer_zs_start_home` | Configuración `ZoneServer` | 149 | `cMobSuicide` | IA Monstruos / NPCs / Spawns |
| 13 | `ZoneServer_zs_start_house` | Configuración `ZoneServer` | 150 | `cMobWeaponRate` | IA Monstruos / NPCs / Spawns |
| 14 | `ZoneServer_zs_start_long` | Configuración `ZoneServer` | 151 | `cMove2Where` | Objetos / Movimiento / Mensajes |
| 15 | `ZoneServer_zs_start_moon` | Configuración `ZoneServer` | 152 | `cMoveLayer` | Mapa / Puertas `.sbi` / Capas |
| 16 | `ZoneServer_zs_start_open` | Configuración `ZoneServer` | 153 | `cMoveStop` | Objetos / Movimiento / Mensajes |
| 17 | `ZoneServer_zs_start_sink` | Configuración `ZoneServer` | 154 | `cNPCChat` | IA Monstruos / NPCs / Spawns |
| 18 | `ZoneServer_zs_start_sky` | Configuración `ZoneServer` | 155 | `cNPCChatTest` | IA Monstruos / NPCs / Spawns |
| 19 | `ZoneServer_zs_start_sleep` | Configuración `ZoneServer` | 156 | `cNPCMenuOpen` | IA Monstruos / NPCs / Spawns |
| 20 | `ZoneServer_zs_start_star` | Configuración `ZoneServer` | 157 | `cNPCRegen` | IA Monstruos / NPCs / Spawns |
| 21 | `ZoneServer_zs_start_sun` | Configuración `ZoneServer` | 158 | `cNPCSkillUse` | IA Monstruos / NPCs / Spawns |
| 22 | `ZoneServer_zs_start_tough` | Configuración `ZoneServer` | 159 | `cNPCVanish` | IA Monstruos / NPCs / Spawns |
| 23 | `ZoneServer_zs_start_wide` | Configuración `ZoneServer` | 160 | `cNearObjectList` | Objetos / Movimiento / Mensajes |
| 24 | `ZoneServer_zs_start_yellow` | Configuración `ZoneServer` | 161 | `cNotice` | Objetos / Movimiento / Mensajes |
| 25 | `cAIScriptFunc` | IA Monstruos / NPCs / Spawns | 162 | `cNoticeRedWarningCode` | Objetos / Movimiento / Mensajes |
| 26 | `cAIScriptSet` | IA Monstruos / NPCs / Spawns | 163 | `cNoticeString` | Objetos / Movimiento / Mensajes |
| 27 | `cAbstateRestTime` | Combate / Skills / `AbState` | 164 | `cNotice_Obj` | Objetos / Movimiento / Mensajes |
| 28 | `cActByEventCode` | Objetos / Movimiento / Mensajes | 165 | `cObjectCount` | Objetos / Movimiento / Mensajes |
| 29 | `cAddCharacterTitle` | Objetos / Movimiento / Mensajes | 166 | `cObjectEffect` | Objetos / Movimiento / Mensajes |
| 30 | `cAggroList` | IA Monstruos / NPCs / Spawns | 167 | `cObjectFind` | Objetos / Movimiento / Mensajes |
| 31 | `cAggroListSize` | IA Monstruos / NPCs / Spawns | 168 | `cObjectHP` | Objetos / Movimiento / Mensajes |
| 32 | `cAggroReset` | IA Monstruos / NPCs / Spawns | 169 | `cObjectLocate` | Objetos / Movimiento / Mensajes |
| 33 | `cAggroSet` | IA Monstruos / NPCs / Spawns | 170 | `cObjectLocateDirection` | Objetos / Movimiento / Mensajes |
| 34 | `cAnimate` | Objetos / Movimiento / Mensajes | 171 | `cObjectSound` | Objetos / Movimiento / Mensajes |
| 35 | `cAnimateForcedly` | Objetos / Movimiento / Mensajes | 172 | `cObjectType` | Objetos / Movimiento / Mensajes |
| 36 | `cAssertLog` | Objetos / Movimiento / Mensajes | 173 | `cPartyJoin` | Objetos / Movimiento / Mensajes |
| 37 | `cCameraMove` | Objetos / Movimiento / Mensajes | 174 | `cPartyLeave` | Objetos / Movimiento / Mensajes |
| 38 | `cCameraMove_Obj` | Objetos / Movimiento / Mensajes | 175 | `cPermileRate` | Objetos / Movimiento / Mensajes |
| 39 | `cCastTeleport` | Objetos / Movimiento / Mensajes | 176 | `cPet_ChangeMind` | Mascotas (`Pet`) |
| 40 | `cCastingBar` | Objetos / Movimiento / Mensajes | 177 | `cPet_GetActionRecord` | Mascotas (`Pet`) |
| 41 | `cCharTitleAddValue` | Objetos / Movimiento / Mensajes | 178 | `cPet_GetMind` | Mascotas (`Pet`) |
| 42 | `cClassChangeOpen` | Objetos / Movimiento / Mensajes | 179 | `cPet_GetRunSpeed` | Mascotas (`Pet`) |
| 43 | `cCurSec` | Objetos / Movimiento / Mensajes | 180 | `cPet_GetStress` | Mascotas (`Pet`) |
| 44 | `cCurrentSecond` | Objetos / Movimiento / Mensajes | 181 | `cPet_GetWalkSpeed` | Mascotas (`Pet`) |
| 45 | `cDamaged` | Combate / Skills / `AbState` | 182 | `cPet_IsMasterCalling` | Mascotas (`Pet`) |
| 46 | `cDePolymorph` | Objetos / Movimiento / Mensajes | 183 | `cPet_SaveTendency` | Mascotas (`Pet`) |
| 47 | `cDePolymorph_Area` | Objetos / Movimiento / Mensajes | 184 | `cPet_Unsummon` | Mascotas (`Pet`) |
| 48 | `cDebugLog` | Objetos / Movimiento / Mensajes | 185 | `cPlaySound` | Objetos / Movimiento / Mensajes |
| 49 | `cDelDirectionalArrow` | Objetos / Movimiento / Mensajes | 186 | `cPlayerExist` | Objetos / Movimiento / Mensajes |
| 50 | `cDirectionalArrow` | Objetos / Movimiento / Mensajes | 187 | `cProgressTutorial` | Misiones / Tutorial / Menús |
| 51 | `cDistanceSquar` | Objetos / Movimiento / Mensajes | 188 | `cQuestMobKill` | IA Monstruos / NPCs / Spawns |
| 52 | `cDoorAction` | Mapa / Puertas `.sbi` / Capas | 189 | `cQuestMobKill_AllInMap` | Mapa / Puertas `.sbi` / Capas |
| 53 | `cDoorBuild` | Mapa / Puertas `.sbi` / Capas | 190 | `cQuestResult` | Misiones / Tutorial / Menús |
| 54 | `cDropFilm` | Ítems / Botín / Recompensas | 191 | `cQuestResult_Individual` | Misiones / Tutorial / Menús |
| 55 | `cDropItem` | Ítems / Botín / Recompensas | 192 | `cRandom` | Objetos / Movimiento / Mensajes |
| 56 | `cEffectMsg` | Objetos / Movimiento / Mensajes | 193 | `cRandomInt` | Objetos / Movimiento / Mensajes |
| 57 | `cEffectMsg_AllInMap` | Mapa / Puertas `.sbi` / Capas | 194 | `cRegenGroupActiv` | IA Monstruos / NPCs / Spawns |
| 58 | `cEffectRegen_Object` | IA Monstruos / NPCs / Spawns | 195 | `cResetAbstate` | Combate / Skills / `AbState` |
| 59 | `cEffectRegen_XY` | IA Monstruos / NPCs / Spawns | 196 | `cRevivalAll` | Combate / Skills / `AbState` |
| 60 | `cEffectTimer` | Objetos / Movimiento / Mensajes | 197 | `cReward` | Ítems / Botín / Recompensas |
| 61 | `cEmotion` | Objetos / Movimiento / Mensajes | 198 | `cRewardItem` | Ítems / Botín / Recompensas |
| 62 | `cEndCinematicText` | Misiones / Tutorial / Menús | 199 | `cRewardItem_AllInMap` | Mapa / Puertas `.sbi` / Capas |
| 63 | `cEndOfKingdomQuest` | Kingdom Quest / Ranking | 200 | `cRewardItem_CharInven` | Ítems / Botín / Recompensas |
| 64 | `cExecCheck` | Objetos / Movimiento / Mensajes | 201 | `cRunTo` | Objetos / Movimiento / Mensajes |
| 65 | `cFieldScriptFunc` | Objetos / Movimiento / Mensajes | 202 | `cRunToUntilBlock` | Objetos / Movimiento / Mensajes |
| 66 | `cFindAttackBlockLocate` | Objetos / Movimiento / Mensajes | 203 | `cSaveTutorialInfo` | Misiones / Tutorial / Menús |
| 67 | `cFindNearPlayer` | Objetos / Movimiento / Mensajes | 204 | `cScoreBoard` | Kingdom Quest / Ranking |
| 68 | `cFindNearestMobList` | IA Monstruos / NPCs / Spawns | 205 | `cScoreBoard_AllInMap` | Kingdom Quest / Ranking |
| 69 | `cFollow` | Objetos / Movimiento / Mensajes | 206 | `cScoreInfo` | Kingdom Quest / Ranking |
| 70 | `cFreeStatInit` | Objetos / Movimiento / Mensajes | 207 | `cScoreInfo_AllInMap` | Kingdom Quest / Ranking |
| 71 | `cGetAbstate` | Combate / Skills / `AbState` | 208 | `cScoreResultList` | Kingdom Quest / Ranking |
| 72 | `cGetAdminLevel` | Objetos / Movimiento / Mensajes | 209 | `cScoreTopList` | Kingdom Quest / Ranking |
| 73 | `cGetAreaObject` | Objetos / Movimiento / Mensajes | 210 | `cScriptMessage` | Objetos / Movimiento / Mensajes |
| 74 | `cGetAreaObjectList` | Objetos / Movimiento / Mensajes | 211 | `cScriptMessage_Obj` | Objetos / Movimiento / Mensajes |
| 75 | `cGetAroundCoord` | Objetos / Movimiento / Mensajes | 212 | `cScriptMessage_Range` | Objetos / Movimiento / Mensajes |
| 76 | `cGetBaseClass` | Objetos / Movimiento / Mensajes | 213 | `cScriptMsg` | Objetos / Movimiento / Mensajes |
| 77 | `cGetCharNo` | Objetos / Movimiento / Mensajes | 214 | `cScriptMsg_World` | Objetos / Movimiento / Mensajes |
| 78 | `cGetCoord_Circle` | Objetos / Movimiento / Mensajes | 215 | `cSendGameLogDataType_4` | Auditoría `GameLog` |
| 79 | `cGetCurMapIndex` | Mapa / Puertas `.sbi` / Capas | 216 | `cSendGameLogDataType_5` | Auditoría `GameLog` |
| 80 | `cGetCurrentSkillInfo` | Combate / Skills / `AbState` | 217 | `cSendRankingList` | Kingdom Quest / Ranking |
| 81 | `cGetDirect` | Objetos / Movimiento / Mensajes | 218 | `cServerMenu` | Misiones / Tutorial / Menús |
| 82 | `cGetItemIndex` | Ítems / Botín / Recompensas | 219 | `cSetAIScript` | IA Monstruos / NPCs / Spawns |
| 83 | `cGetItemLot` | Ítems / Botín / Recompensas | 220 | `cSetAbstate` | Combate / Skills / `AbState` |
| 84 | `cGetKQLimitSecond` | Kingdom Quest / Ranking | 221 | `cSetAbstateInArea` | Combate / Skills / `AbState` |
| 85 | `cGetKQTeamType` | Kingdom Quest / Ranking | 222 | `cSetAbstate_Range` | Combate / Skills / `AbState` |
| 86 | `cGetLevel` | Objetos / Movimiento / Mensajes | 223 | `cSetCanUseReviveItem` | Combate / Skills / `AbState` |
| 87 | `cGetMapName` | Mapa / Puertas `.sbi` / Capas | 224 | `cSetCanUseReviveSkill` | Combate / Skills / `AbState` |
| 88 | `cGetMaster` | Objetos / Movimiento / Mensajes | 225 | `cSetDeadDelayTime` | Objetos / Movimiento / Mensajes |
| 89 | `cGetMobID` | IA Monstruos / NPCs / Spawns | 226 | `cSetFieldScript` | Objetos / Movimiento / Mensajes |
| 90 | `cGetMoveState` | Objetos / Movimiento / Mensajes | 227 | `cSetFreeBattle` | Objetos / Movimiento / Mensajes |
| 91 | `cGetNPCHandle` | IA Monstruos / NPCs / Spawns | 228 | `cSetItemDropMobID` | IA Monstruos / NPCs / Spawns |
| 92 | `cGetNearObjListByCoord` | Objetos / Movimiento / Mensajes | 229 | `cSetMobAttr` | IA Monstruos / NPCs / Spawns |
| 93 | `cGetNearestObjByCoord` | Objetos / Movimiento / Mensajes | 230 | `cSetNPCIsItemDrop` | IA Monstruos / NPCs / Spawns |
| 94 | `cGetObjectMode` | Objetos / Movimiento / Mensajes | 231 | `cSetNPCParam` | IA Monstruos / NPCs / Spawns |
| 95 | `cGetPlayerList` | Objetos / Movimiento / Mensajes | 232 | `cSetNPCResist` | IA Monstruos / NPCs / Spawns |
| 96 | `cGetPlayerName` | Objetos / Movimiento / Mensajes | 233 | `cSetObjectDirect` | Objetos / Movimiento / Mensajes |
| 97 | `cGetQuestHero_ItemUse` | Misiones / Tutorial / Menús | 234 | `cSetReviveDelayTime` | Combate / Skills / `AbState` |
| 98 | `cGetQuestHero_NPC` | IA Monstruos / NPCs / Spawns | 235 | `cSetServantFlag` | IA Monstruos / NPCs / Spawns |
| 99 | `cGetRegistNumber` | Objetos / Movimiento / Mensajes | 236 | `cSetSightState` | Mapa / Puertas `.sbi` / Capas |
| 100 | `cGetRunSpeed` | Objetos / Movimiento / Mensajes | 237 | `cSetTeamBattle` | Objetos / Movimiento / Mensajes |
| 101 | `cGetScriptString` | Objetos / Movimiento / Mensajes | 238 | `cShowCinematicText` | Misiones / Tutorial / Menús |
| 102 | `cGetTargetList` | Objetos / Movimiento / Mensajes | 239 | `cShowKQTimerWithLife` | Kingdom Quest / Ranking |
| 103 | `cGetTutorialInfo` | Misiones / Tutorial / Menús | 240 | `cShowKQTimerWithLife_Obj` | Kingdom Quest / Ranking |
| 104 | `cGetWalkSpeed` | Objetos / Movimiento / Mensajes | 241 | `cSimpleChatScriptMsg` | Objetos / Movimiento / Mensajes |
| 105 | `cGetWhoKillMe` | Combate / Skills / `AbState` | 242 | `cSkillBlast` | Combate / Skills / `AbState` |
| 106 | `cGroupRegen` | IA Monstruos / NPCs / Spawns | 243 | `cSkillDamageCalculate` | Combate / Skills / `AbState` |
| 107 | `cGroupRegenInstance` | IA Monstruos / NPCs / Spawns | 244 | `cSomebodyShout` | Objetos / Movimiento / Mensajes |
| 108 | `cGroupRegenInstance_XY` | IA Monstruos / NPCs / Spawns | 245 | `cStartMsg_AllInMap` | Mapa / Puertas `.sbi` / Capas |
| 109 | `cHeal` | Combate / Skills / `AbState` | 246 | `cStaticDamage` | Combate / Skills / `AbState` |
| 110 | `cHideOtherPlayer` | Objetos / Movimiento / Mensajes | 247 | `cStaticDamage_smo` | Combate / Skills / `AbState` |
| 111 | `cIndunRankResult` | Objetos / Movimiento / Mensajes | 248 | `cStaticMoverSpeed` | Objetos / Movimiento / Mensajes |
| 112 | `cInvenItemDestroy` | Ítems / Botín / Recompensas | 249 | `cStaticRunSpeed` | Objetos / Movimiento / Mensajes |
| 113 | `cIsEquipItem` | Ítems / Botín / Recompensas | 250 | `cStaticSpeed` | Objetos / Movimiento / Mensajes |
| 114 | `cIsInArea` | Objetos / Movimiento / Mensajes | 251 | `cStaticWalkSpeed` | Objetos / Movimiento / Mensajes |
| 115 | `cIsInMap` | Mapa / Puertas `.sbi` / Capas | 252 | `cSystemMessage_Obj` | Objetos / Movimiento / Mensajes |
| 116 | `cIsKQJoiner` | Kingdom Quest / Ranking | 253 | `cTargetChangeNull` | Objetos / Movimiento / Mensajes |
| 117 | `cIsMovable` | Objetos / Movimiento / Mensajes | 254 | `cTargetHandle` | Objetos / Movimiento / Mensajes |
| 118 | `cIsNoAttacOrNoMove` | Objetos / Movimiento / Mensajes | 255 | `cTimer` | Objetos / Movimiento / Mensajes |
| 119 | `cIsObjectAlreadyDead` | Objetos / Movimiento / Mensajes | 256 | `cTimerEnd` | Objetos / Movimiento / Mensajes |
| 120 | `cIsObjectDead` | Objetos / Movimiento / Mensajes | 257 | `cTimerStart` | Objetos / Movimiento / Mensajes |
| 121 | `cItemErase` | Ítems / Botín / Recompensas | 258 | `cTimer_Obj` | Objetos / Movimiento / Mensajes |
| 122 | `cKQRewardIndex` | Kingdom Quest / Ranking | 259 | `cTopView` | Objetos / Movimiento / Mensajes |
| 123 | `cKillObject` | Combate / Skills / `AbState` | 260 | `cTriggerRegist` | Objetos / Movimiento / Mensajes |
| 124 | `cLevelUp` | Objetos / Movimiento / Mensajes | 261 | `cUseMoney` | Ítems / Botín / Recompensas |
| 125 | `cLinkTo` | Objetos / Movimiento / Mensajes | 262 | `cVanishAll` | Objetos / Movimiento / Mensajes |
| 126 | `cLinkToAll` | Objetos / Movimiento / Mensajes | 263 | `cVanishReserv` | Objetos / Movimiento / Mensajes |
| 127 | `cLocationRotate` | Objetos / Movimiento / Mensajes | 264 | `cVanishTimer` | Objetos / Movimiento / Mensajes |
| 128 | `cMagicFieldSpread` | Objetos / Movimiento / Mensajes | 265 | `cViewSlotEquip` | Ítems / Botín / Recompensas |
| 129 | `cMapFog` | Mapa / Puertas `.sbi` / Capas | 266 | `cViewSlotUnEquipAll` | Ítems / Botín / Recompensas |
| 130 | `cMapMark` | Mapa / Puertas `.sbi` / Capas | 267 | `cWaitIdle` | Objetos / Movimiento / Mensajes |
| 131 | `cMapMark_FieldSight` | Mapa / Puertas `.sbi` / Capas | 268 | `cWalkTo` | Objetos / Movimiento / Mensajes |
| 132 | `cMapMark_Obj` | Mapa / Puertas `.sbi` / Capas | 269 | `cWillMovement` | Objetos / Movimiento / Mensajes |
| 133 | `cMapObjectControl` | Mapa / Puertas `.sbi` / Capas | 270 | `cWinter_Event_ScoreBoard_AllInMap` | Kingdom Quest / Ranking |
| 134 | `cMessage` | Objetos / Movimiento / Mensajes | 271 | `cpp_CurrentTime` | Objetos / Movimiento / Mensajes |
| 135 | `cMobChat` | IA Monstruos / NPCs / Spawns | 272 | `cpp_Print` | Objetos / Movimiento / Mensajes |
| 136 | `cMobDetectRange` | IA Monstruos / NPCs / Spawns | 273 | `cpp_RegCheck` | Objetos / Movimiento / Mensajes |
| 137 | `cMobDialog` | IA Monstruos / NPCs / Spawns | 274 | `tonumber` | Objetos / Movimiento / Mensajes |

---

## 4. Catálogo 100% Exhaustivo de los Dispatchers C2S y S2S (`Zone.exe` y `Zone.pdb`)

### 4.1 Las 245 Entradas Oficiales del Dispatcher C2S (`protocolstore(ShinePlayer)` en `0x004D3510`)

Cada una de las **245 entradas C2S** registradas en `Zone.exe` debe estar enrutada en `pkg/server/zone/handler.go` con su constante oficial `PROTO_NC_*` y su lógica real:

| # | OpCode (Hex) | OpCode (Dec) | Dept | Cmd | Callback en `Zone.exe` | Estructura Oficial `Zone.pdb` | Estado Hoy en `handler.go` | Fase Asignada |
| :---: | :---: | :---: | :---: | :---: | :--- | :--- | :---: | :--- |
| 1 | `0x0804` | `2052` | `2` (NC_MISC) | `4` | `iac_FoodEat` | `iac_FoodEat` | **Faltante** | Fase 14 (Consumibles iac_FoodEat) |
| 2 | `0x0805` | `2053` | `2` (NC_MISC) | `5` | `iac_FoodEat` | `iac_FoodEat` | **Faltante** | Fase 14 (Consumibles iac_FoodEat) |
| 3 | `0x0806` | `2054` | `2` (NC_MISC) | `6` | `sp_NC_MISC_SEED_REQ` | `PROTO_NC_MISC_SEED_REQ` | **Faltante** | Fase 11 (Semilla Cripto & Visibilidad AOI) |
| 4 | `0x0812` | `2066` | `2` (NC_MISC) | `18` | `sp_NC_MISC_PINGTEST_CLIENT_ZONE_DB` | `PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB` | **Faltante** | Fase 16 (Diagnóstico PingTest & ServerParameter) |
| 5 | `0x0813` | `2067` | `2` (NC_MISC) | `19` | `sp_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB` | `PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB` | **Faltante** | Fase 16 (Diagnóstico PingTest & ServerParameter) |
| 6 | `0x0817` | `2071` | `2` (NC_MISC) | `23` | `sp_NC_MISC_HIDE_EXCEPT_ME_ON_CMD` | `PROTO_NC_MISC_HIDE_EXCEPT_ME_ON_CMD` | **Faltante** | Fase 11 (Semilla Cripto & Visibilidad AOI) |
| 7 | `0x0818` | `2072` | `2` (NC_MISC) | `24` | `sp_NC_MISC_HIDE_EXCEPT_ME_OFF_CMD` | `PROTO_NC_MISC_HIDE_EXCEPT_ME_OFF_CMD` | **Faltante** | Fase 11 (Semilla Cripto & Visibilidad AOI) |
| 8 | `0x082F` | `2095` | `2` (NC_MISC) | `47` | `sp_NC_MISC_SERVERPARAMETER_REQ` | `PROTO_NC_MISC_SERVERPARAMETER_REQ` | **Faltante** | Fase 16 (Diagnóstico PingTest & ServerParameter) |
| 9 | `0x0C18` | `3096` | `3` (NC_USER) | `24` | `sp_NC_USER_NORMALLOGOUT_CMD` | `PROTO_NC_USER_NORMALLOGOUT_CMD` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 10 | `0x0C23` | `3107` | `3` (NC_USER) | `35` | `sp_NC_USER_REGISENUMBER_REQ` | `PROTO_NC_USER_REGISENUMBER_REQ` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 11 | `0x0C32` | `3122` | `3` (NC_USER) | `50` | `sp_NC_USER_USE_BEAUTY_SHOP_CMD` | `PROTO_NC_USER_USE_BEAUTY_SHOP_CMD` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 12 | `0x104E` | `4174` | `4` (NC_CHAR) | `78` | `sp_NC_CHAR_REVIVE_REQ` | `PROTO_NC_CHAR_REVIVE_REQ` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 13 | `0x105C` | `4188` | `4` (NC_CHAR) | `92` | `sp_NC_CHAR_STAT_INCPOINT_REQ` | `PROTO_NC_CHAR_STAT_INCPOINT_REQ` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 14 | `0x1062` | `4194` | `4` (NC_CHAR) | `98` | `sp_NC_CHAR_STAT_DECPOINT_REQ` | `PROTO_NC_CHAR_STAT_DECPOINT_REQ` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 15 | `0x1071` | `4209` | `4` (NC_CHAR) | `113` | `sp_NC_CHAR_LOGOUTREADY_CMD` | `PROTO_NC_CHAR_LOGOUTREADY_CMD` | **Faltante** | Fase 11 / Fase 16 (Sesión en Mapa & Logout) |
| 16 | `0x1072` | `4210` | `4` (NC_CHAR) | `114` | `sp_NC_CHAR_LOGOUTCANCEL_CMD` | `PROTO_NC_CHAR_LOGOUTCANCEL_CMD` | **Faltante** | Fase 11 / Fase 16 (Sesión en Mapa & Logout) |
| 17 | `0x1093` | `4243` | `4` (NC_CHAR) | `147` | `sp_NC_CHAR_WEDDING_PARTNER_INFO_REQ` | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ` | **Faltante** | Fase 15 (Tutorial & Bodas en Mapa) |
| 18 | `0x10B0` | `4272` | `4` (NC_CHAR) | `176` | `sp_NC_CHAR_DEPOLYMORPH_CMD` | `PROTO_NC_CHAR_DEPOLYMORPH_CMD` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 19 | `0x10B9` | `4281` | `4` (NC_CHAR) | `185` | `sp_NC_CHAR_SAVE_LINK_REQ` | `PROTO_NC_CHAR_SAVE_LINK_REQ` | **Faltante** | Fase 11 / Fase 16 (Sesión en Mapa & Logout) |
| 20 | `0x10BC` | `4284` | `4` (NC_CHAR) | `188` | `sp_NC_CHAR_CLIENT_AUTO_PICK_REQ` | `PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ` | **Faltante** | Fase 13 (Drops & AutoPick) |
| 21 | `0x1102` | `4354` | `4` (NC_CHAR) | `258` | `sp_NC_CHAR_REBIRTH_CMD` | `PROTO_NC_CHAR_REBIRTH_CMD` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 22 | `0x1103` | `4355` | `4` (NC_CHAR) | `259` | `sp_NC_CHAR_REBIRTH_REJECT_CMD` | `PROTO_NC_CHAR_REBIRTH_REJECT_CMD` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 23 | `0x1117` | `4375` | `4` (NC_CHAR) | `279` | `sp_NC_CHAR_TUTORIAL_STEP_ACK` | `PROTO_NC_CHAR_TUTORIAL_STEP_ACK` | **Faltante** | Fase 15 (Tutorial & Bodas en Mapa) |
| 24 | `0x1124` | `4388` | `4` (NC_CHAR) | `292` | `sp_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ` | `PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ` | **Faltante** | Fase 14 (Equipamiento & Minimon) |
| 25 | `0x1126` | `4390` | `4` (NC_CHAR) | `294` | `sp_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ` | `PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ` | **Faltante** | Fase 14 (Equipamiento & Minimon) |
| 26 | `0x1128` | `4392` | `4` (NC_CHAR) | `296` | `sp_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ` | `PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ` | **Faltante** | Fase 14 (Equipamiento & Minimon) |
| 27 | `0x112A` | `4394` | `4` (NC_CHAR) | `298` | `sp_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ` | `PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ` | **Faltante** | Fase 14 (Equipamiento & Minimon) |
| 28 | `0x112E` | `4398` | `4` (NC_CHAR) | `302` | `sp_NC_CHAR_CHARGEDBUFF_ERASE_REQ` | `PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ` | **Faltante** | Fase 12 (Combate, Stats, AbState & Resurrección) |
| 29 | `0x1801` | `6145` | `6` (NC_MAP) | `1` | `sp_NC_MAP_LOGIN_REQ` | `PROTO_NC_MAP_LOGIN_REQ` | Parcial / Base | Fase 11 (Mapas, Links y Portales) |
| 30 | `0x1803` | `6147` | `6` (NC_MAP) | `3` | `sp_NC_MAP_LOGINCOMPLETE_CMD` | `PROTO_NC_MAP_LOGINCOMPLETE_CMD` | Parcial / Base | Fase 11 (Mapas, Links y Portales) |
| 31 | `0x1814` | `6164` | `6` (NC_MAP) | `20` | `sp_NC_MAP_WING_SAVE_REQ` | `PROTO_NC_MAP_WING_SAVE_REQ` | **Faltante** | Fase 11 (Mapas, Links y Portales) |
| 32 | `0x1816` | `6166` | `6` (NC_MAP) | `22` | `sp_NC_MAP_WING_FLY_REQ` | `PROTO_NC_MAP_WING_FLY_REQ` | Parcial / Base | Fase 11 (Mapas, Links y Portales) |
| 33 | `0x181A` | `6170` | `6` (NC_MAP) | `26` | `sp_NC_MAP_TOWNPORTAL_REQ` | `PROTO_NC_MAP_TOWNPORTAL_REQ` | Parcial / Base | Fase 11 (Mapas, Links y Portales) |
| 34 | `0x181C` | `6172` | `6` (NC_MAP) | `28` | `sp_NC_MAP_TONORMALCOORD_CMD` | `PROTO_NC_MAP_TONORMALCOORD_CMD` | **Faltante** | Fase 11 (Mapas, Links y Portales) |
| 35 | `0x181F` | `6175` | `6` (NC_MAP) | `31` | `sp_NC_MAP_MULTY_LINK_SELECT_REQ` | `PROTO_NC_MAP_MULTY_LINK_SELECT_REQ` | Parcial / Base | Fase 11 (Mapas, Links y Portales) |
| 36 | `0x1C01` | `7169` | `7` (NC_BRIEFINFO) | `1` | `sp_NC_BRIEFINFO_INFORM_CMD` | `PROTO_NC_BRIEFINFO_INFORM_CMD` | **Faltante** | Fase 11 (AOI & BRIEFINFO) |
| 37 | `0x2001` | `8193` | `8` (NC_ACT) | `1` | `sp_NC_ACT_CHAT_REQ` | `PROTO_NC_ACT_CHAT_REQ` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 38 | `0x2003` | `8195` | `8` (NC_ACT) | `3` | `sp_NC_ACT_WALK_REQ` | `PROTO_NC_ACT_WALK_REQ` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 39 | `0x2005` | `8197` | `8` (NC_ACT) | `5` | `sp_NC_ACT_RUN_REQ` | `PROTO_NC_ACT_RUN_REQ` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 40 | `0x2008` | `8200` | `8` (NC_ACT) | `8` | `sp_NC_ACT_CHANGEMODE_REQ` | `PROTO_NC_ACT_CHANGEMODE_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 41 | `0x200A` | `8202` | `8` (NC_ACT) | `10` | `sp_NC_ACT_NPCCLICK_CMD` | `PROTO_NC_ACT_NPCCLICK_CMD` | Parcial / Base | Fase 15 (NPCs, Menús & Quests) |
| 42 | `0x200B` | `8203` | `8` (NC_ACT) | `11` | `sp_NC_ACT_ENDOFTRADE_CMD` | `PROTO_NC_ACT_ENDOFTRADE_CMD` | Parcial / Base | Fase 14 (Economía, Trade & Produce) |
| 43 | `0x2012` | `8210` | `8` (NC_ACT) | `18` | `sp_NC_ACT_STOP_REQ` | `PROTO_NC_ACT_STOP_REQ` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 44 | `0x2017` | `8215` | `8` (NC_ACT) | `23` | `sp_NC_ACT_MOVEWALK_CMD` | `PROTO_NC_ACT_MOVEWALK_CMD` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 45 | `0x2019` | `8217` | `8` (NC_ACT) | `25` | `sp_NC_ACT_MOVERUN_CMD` | `PROTO_NC_ACT_MOVERUN_CMD` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 46 | `0x201D` | `8221` | `8` (NC_ACT) | `29` | `sp_NC_ACT_NPCMENUOPEN_ACK` | `PROTO_NC_ACT_NPCMENUOPEN_ACK` | Parcial / Base | Fase 15 (NPCs, Menús & Quests) |
| 47 | `0x201E` | `8222` | `8` (NC_ACT) | `30` | `sp_NC_ACT_SHOUT_CMD` | `PROTO_NC_ACT_SHOUT_CMD` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 48 | `0x2020` | `8224` | `8` (NC_ACT) | `32` | `sp_NC_ACT_EMOTICON_CMD` | `PROTO_NC_ACT_EMOTICON_CMD` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 49 | `0x2022` | `8226` | `8` (NC_ACT) | `34` | `sp_NC_ACT_EMOTICONSTOP_CMD` | `PROTO_NC_ACT_EMOTICONSTOP_CMD` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 50 | `0x2024` | `8228` | `8` (NC_ACT) | `36` | `sp_NC_ACT_JUMP_CMD` | `PROTO_NC_ACT_JUMP_CMD` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 51 | `0x2027` | `8231` | `8` (NC_ACT) | `39` | `sp_NC_ACT_PITCHTENT_REQ` | `PROTO_NC_ACT_PITCHTENT_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 52 | `0x202A` | `8234` | `8` (NC_ACT) | `42` | `sp_NC_ACT_FOLDTENT_REQ` | `PROTO_NC_ACT_FOLDTENT_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 53 | `0x202D` | `8237` | `8` (NC_ACT) | `45` | `sp_NC_ACT_GATHERSTART_REQ` | `PROTO_NC_ACT_GATHERSTART_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 54 | `0x2030` | `8240` | `8` (NC_ACT) | `48` | `sp_NC_ACT_GATHERCANCEL_CMD` | `PROTO_NC_ACT_GATHERCANCEL_CMD` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 55 | `0x2032` | `8242` | `8` (NC_ACT) | `50` | `sp_NC_ACT_GATHERCOMPLETE_REQ` | `PROTO_NC_ACT_GATHERCOMPLETE_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 56 | `0x2035` | `8245` | `8` (NC_ACT) | `53` | `sp_NC_ACT_PRODUCE_CAST_REQ` | `PROTO_NC_ACT_PRODUCE_CAST_REQ` | **Faltante** | Fase 14 (Economía, Trade & Produce) |
| 57 | `0x2039` | `8249` | `8` (NC_ACT) | `57` | `sp_NC_ACT_PRODUCE_CASTABORT_CMD` | `PROTO_NC_ACT_PRODUCE_CASTABORT_CMD` | **Faltante** | Fase 14 (Economía, Trade & Produce) |
| 58 | `0x2044` | `8260` | `8` (NC_ACT) | `68` | `sp_NC_ACT_RIDE_FEEDING_REQ` | `PROTO_NC_ACT_RIDE_FEEDING_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 59 | `0x204B` | `8267` | `8` (NC_ACT) | `75` | `sp_NC_ACT_ROAR_REQ` | `PROTO_NC_ACT_ROAR_REQ` | **Faltante** | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 60 | `0x204F` | `8271` | `8` (NC_ACT) | `79` | `sp_NC_ACT_WEDDING_PROPOSEREQ_REQ` | `PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 61 | `0x2051` | `8273` | `8` (NC_ACT) | `81` | `sp_NC_ACT_WEDDING_PROPOSEACK_ACK` | `PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 62 | `0x2053` | `8275` | `8` (NC_ACT) | `83` | `sp_NC_ACT_WEDDING_HALL_RESERV_REQ` | `PROTO_NC_ACT_WEDDING_HALL_RESERV_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 63 | `0x2056` | `8278` | `8` (NC_ACT) | `86` | `sp_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_REQ` | `PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 64 | `0x2058` | `8280` | `8` (NC_ACT) | `88` | `sp_NC_ACT_WEDDING_HALL_GUEST_ENTER_REQ` | `PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 65 | `0x205B` | `8283` | `8` (NC_ACT) | `91` | `sp_NC_ACT_WEDDING_AGREEMENT_DIVORCE_REQ` | `PROTO_NC_ACT_WEDDING_AGREEMENT_DIVORCE_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 66 | `0x205C` | `8284` | `8` (NC_ACT) | `92` | `sp_NC_ACT_WEDDING_COMPULSORY_DIVORCE_REQ` | `PROTO_NC_ACT_WEDDING_COMPULSORY_DIVORCE_REQ` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 67 | `0x205E` | `8286` | `8` (NC_ACT) | `94` | `sp_NC_ACT_WEDDING_DIVORCE_REFUSE_CMD` | `PROTO_NC_ACT_WEDDING_DIVORCE_REFUSE_CMD` | **Faltante** | Fase 15 (Bodas & Instancias) |
| 68 | `0x2067` | `8295` | `8` (NC_ACT) | `103` | `sp_NC_ACT_ACTIONBYITEM_REQ` | `PROTO_NC_ACT_ACTIONBYITEM_REQ` | **Faltante** | Fase 14 (Economía, Trade & Produce) |
| 69 | `0x2070` | `8304` | `8` (NC_ACT) | `112` | `sp_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ` | `PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ` | Parcial / Base | Fase 11 (Movimiento, Tent, Emotes & Gather) |
| 70 | `0x2073` | `8307` | `8` (NC_ACT) | `115` | `sp_NC_ACT_NPC_MENU_CMD` | `PROTO_NC_ACT_NPC_MENU_CMD` | Parcial / Base | Fase 15 (NPCs, Menús & Quests) |
| 71 | `0x2401` | `9217` | `9` (NC_BAT) | `1` | `sp_NC_BAT_TARGETTING_REQ` | `PROTO_NC_BAT_TARGETTING_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 72 | `0x2408` | `9224` | `9` (NC_BAT) | `8` | `sp_NC_BAT_UNTARGET_REQ` | `PROTO_NC_BAT_UNTARGET_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 73 | `0x242B` | `9259` | `9` (NC_BAT) | `43` | `sp_NC_BAT_BASHSTART_CMD` | `PROTO_NC_BAT_BASHSTART_CMD` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 74 | `0x2432` | `9266` | `9` (NC_BAT) | `50` | `sp_NC_BAT_BASHSTOP_CMD` | `PROTO_NC_BAT_BASHSTOP_CMD` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 75 | `0x243E` | `9278` | `9` (NC_BAT) | `62` | `sp_NC_BAT_ASSIST_REQ` | `PROTO_NC_BAT_ASSIST_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 76 | `0x2440` | `9280` | `9` (NC_BAT) | `64` | `sp_NC_BAT_SKILLBASH_OBJ_CAST_REQ` | `PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 77 | `0x2441` | `9281` | `9` (NC_BAT) | `65` | `sp_NC_BAT_SKILLBASH_FLD_CAST_REQ` | `PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 78 | `0x2444` | `9284` | `9` (NC_BAT) | `68` | `sp_NC_BAT_SKILLBASH_CASTABORT_REQ` | `PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 79 | `0x2454` | `9300` | `9` (NC_BAT) | `84` | `sp_NC_BAT_ABSTATE_ERASE_REQ` | `PROTO_NC_BAT_ABSTATE_ERASE_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 80 | `0x3003` | `12291` | `12` (NC_ITEM) | `3` | `sp_NC_ITEM_BUY_REQ` | `PROTO_NC_ITEM_BUY_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 81 | `0x3006` | `12294` | `12` (NC_ITEM) | `6` | `sp_NC_ITEM_SELL_REQ` | `PROTO_NC_ITEM_SELL_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 82 | `0x3007` | `12295` | `12` (NC_ITEM) | `7` | `sp_NC_ITEM_DROP_REQ` | `PROTO_NC_ITEM_DROP_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 83 | `0x3009` | `12297` | `12` (NC_ITEM) | `9` | `sp_NC_ITEM_PICK_REQ` | `PROTO_NC_ITEM_PICK_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 84 | `0x300B` | `12299` | `12` (NC_ITEM) | `11` | `sp_NC_ITEM_RELOC_REQ` | `PROTO_NC_ITEM_RELOC_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 85 | `0x300D` | `12301` | `12` (NC_ITEM) | `13` | `sp_NC_ITEM_SPLIT_REQ` | `PROTO_NC_ITEM_SPLIT_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 86 | `0x300F` | `12303` | `12` (NC_ITEM) | `15` | `sp_NC_ITEM_EQUIP_REQ` | `PROTO_NC_ITEM_EQUIP_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 87 | `0x3010` | `12304` | `12` (NC_ITEM) | `16` | `sp_NC_ITEM_RINGEQUIP_REQ` | `PROTO_NC_ITEM_RINGEQUIP_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 88 | `0x3012` | `12306` | `12` (NC_ITEM) | `18` | `sp_NC_ITEM_UNEQUIP_REQ` | `PROTO_NC_ITEM_UNEQUIP_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 89 | `0x3015` | `12309` | `12` (NC_ITEM) | `21` | `sp_NC_ITEM_USE_REQ` | `PROTO_NC_ITEM_USE_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 90 | `0x3017` | `12311` | `12` (NC_ITEM) | `23` | `sp_NC_ITEM_UPGRADE_REQ` | `PROTO_NC_ITEM_UPGRADE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 91 | `0x3019` | `12313` | `12` (NC_ITEM) | `25` | `sp_NC_ITEM_USEABORT_CMD` | `PROTO_NC_ITEM_USEABORT_CMD` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 92 | `0x301C` | `12316` | `12` (NC_ITEM) | `28` | `sp_NC_ITEM_DEPOSIT_REQ` | `PROTO_NC_ITEM_DEPOSIT_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 93 | `0x301E` | `12318` | `12` (NC_ITEM) | `30` | `sp_NC_ITEM_WITHDRAW_REQ` | `PROTO_NC_ITEM_WITHDRAW_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 94 | `0x3020` | `12320` | `12` (NC_ITEM) | `32` | `sp_NC_ITEM_CHARGEDINVENOPEN_REQ` | `PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 95 | `0x3022` | `12322` | `12` (NC_ITEM) | `34` | `sp_NC_ITEM_CHARGED_WITHDRAW_REQ` | `PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 96 | `0x3026` | `12326` | `12` (NC_ITEM) | `38` | `sp_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ` | `PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 97 | `0x3028` | `12328` | `12` (NC_ITEM) | `40` | `sp_NC_ITEM_OPENSTORAGEPAGE_REQ` | `PROTO_NC_ITEM_OPENSTORAGEPAGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 98 | `0x302A` | `12330` | `12` (NC_ITEM) | `42` | `sp_NC_ITEM_DISMANTLE_REQ` | `PROTO_NC_ITEM_DISMANTLE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 99 | `0x302C` | `12332` | `12` (NC_ITEM) | `44` | `sp_NC_ITEM_REWARDINVENOPEN_REQ` | `PROTO_NC_ITEM_REWARDINVENOPEN_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 100 | `0x302F` | `12335` | `12` (NC_ITEM) | `47` | `sp_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ` | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 101 | `0x3031` | `12337` | `12` (NC_ITEM) | `49` | `sp_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ` | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 102 | `0x3033` | `12339` | `12` (NC_ITEM) | `51` | `sp_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_REQ` | `PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 103 | `0x3035` | `12341` | `12` (NC_ITEM) | `53` | `sp_NC_ITEM_WEAPONENDURE_CHARGE_REQ` | `PROTO_NC_ITEM_WEAPONENDURE_CHARGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 104 | `0x3039` | `12345` | `12` (NC_ITEM) | `57` | `sp_NC_ITEM_REVIVEITEMUSE_CMD` | `PROTO_NC_ITEM_REVIVEITEMUSE_CMD` | **Faltante** | Fase 12 (Combate, Skills & Resurrección) |
| 105 | `0x303C` | `12348` | `12` (NC_ITEM) | `60` | `sp_NC_ITEM_DICE_GAME_START_REQ` | `PROTO_NC_ITEM_DICE_GAME_START_REQ` | **Faltante** | Fase 15 (Casino, Dados & Tragamonedas) |
| 106 | `0x3040` | `12352` | `12` (NC_ITEM) | `64` | `sp_NC_ITEM_ENCHANT_ADD_GEM_REQ` | `PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 107 | `0x3042` | `12354` | `12` (NC_ITEM) | `66` | `sp_NC_ITEM_ENCHANT_REMOVE_GEM_REQ` | `PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 108 | `0x3044` | `12356` | `12` (NC_ITEM) | `68` | `sp_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ` | `PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 109 | `0x3048` | `12360` | `12` (NC_ITEM) | `72` | `sp_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_REQ` | `PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 110 | `0x304A` | `12362` | `12` (NC_ITEM) | `74` | `sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ` | `PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ` | Parcial / Base | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 111 | `0x304C` | `12364` | `12` (NC_ITEM) | `76` | `sp_NC_ITEM_ACCOUNT_STORAGE_OPEN_CMD` | `PROTO_NC_ITEM_ACCOUNT_STORAGE_OPEN_CMD` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 112 | `0x304D` | `12365` | `12` (NC_ITEM) | `77` | `sp_NC_ITEM_ACCOUNT_STORAGE_CLOSE_CMD` | `PROTO_NC_ITEM_ACCOUNT_STORAGE_CLOSE_CMD` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 113 | `0x304E` | `12366` | `12` (NC_ITEM) | `78` | `sp_NC_ITEM_USE_ACTIVESKILL_REQ` | `PROTO_NC_ITEM_USE_ACTIVESKILL_REQ` | **Faltante** | Fase 12 (Combate, Skills & Resurrección) |
| 114 | `0x3050` | `12368` | `12` (NC_ITEM) | `80` | `sp_NC_ITEM_MINIMON_EQUIP_REQ` | `PROTO_NC_ITEM_MINIMON_EQUIP_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 115 | `0x3051` | `12369` | `12` (NC_ITEM) | `81` | `sp_NC_ITEM_MOVER_UPGRADE_REQ` | `PROTO_NC_ITEM_MOVER_UPGRADE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 116 | `0x3053` | `12371` | `12` (NC_ITEM) | `83` | `sp_NC_ITEM_RANDOMOPTION_CHANGE_REQ` | `PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 117 | `0x3055` | `12373` | `12` (NC_ITEM) | `85` | `sp_NC_ITEM_CHAT_COLOR_CHANGE_REQ` | `PROTO_NC_ITEM_CHAT_COLOR_CHANGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 118 | `0x3057` | `12375` | `12` (NC_ITEM) | `87` | `sp_NC_ITEM_TERMEXTEND_REQ` | `PROTO_NC_ITEM_TERMEXTEND_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 119 | `0x3059` | `12377` | `12` (NC_ITEM) | `89` | `sp_NC_ITEM_REPURCHASE_REQ` | `PROTO_NC_ITEM_REPURCHASE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 120 | `0x305E` | `12382` | `12` (NC_ITEM) | `94` | `sp_NC_ITEM_CLASS_CHANGE_REQ` | `PROTO_NC_ITEM_CLASS_CHANGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 121 | `0x3061` | `12385` | `12` (NC_ITEM) | `97` | `sp_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_REQ` | `PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 122 | `0x3064` | `12388` | `12` (NC_ITEM) | `100` | `sp_NC_ITEM_SHIELDENDURE_CHARGE_REQ` | `PROTO_NC_ITEM_SHIELDENDURE_CHARGE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 123 | `0x3067` | `12391` | `12` (NC_ITEM) | `103` | `sp_NC_ITEM_MAPLINK_SCROLL_REQ` | `PROTO_NC_ITEM_MAPLINK_SCROLL_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 124 | `0x3069` | `12393` | `12` (NC_ITEM) | `105` | `sp_NC_ITEM_MIX_ITEM_REQ` | `PROTO_NC_ITEM_MIX_ITEM_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 125 | `0x306B` | `12395` | `12` (NC_ITEM) | `107` | `sp_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ` | `PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 126 | `0x306D` | `12397` | `12` (NC_ITEM) | `109` | `sp_NC_ITEM_NEW_UPGRADE_REQ` | `PROTO_NC_ITEM_NEW_UPGRADE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 127 | `0x306F` | `12399` | `12` (NC_ITEM) | `111` | `sp_NC_ITEM_BRACELET_UPGRADE_REQ` | `PROTO_NC_ITEM_BRACELET_UPGRADE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 128 | `0x3848` | `14408` | `14` (NC_PARTY) | `72` | `sp_NC_PARTY_MEMBERINFOREQ_CMD` | `PROTO_NC_PARTY_MEMBERINFOREQ_CMD` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 129 | `0x3C02` | `15362` | `15` (NC_MENU) | `2` | `sp_NC_MENU_SERVERMENU_ACK` | `PROTO_NC_MENU_SERVERMENU_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 130 | `0x4025` | `16421` | `16` (NC_CHARSAVE) | `37` | `sp_NC_CHARSAVE_UI_STATE_SAVE_REQ` | `PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 131 | `0x4402` | `17410` | `17` (NC_QUEST) | `2` | `sp_NC_QUEST_SCRIPT_CMD_ACK` | `PROTO_NC_QUEST_SCRIPT_CMD_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 132 | `0x4407` | `17415` | `17` (NC_QUEST) | `7` | `sp_NC_QUEST_GIVE_UP_REQ` | `PROTO_NC_QUEST_GIVE_UP_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 133 | `0x440B` | `17419` | `17` (NC_QUEST) | `11` | `sp_NC_QUEST_CLIENT_SCENARIO_DONE_REQ` | `PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 134 | `0x440F` | `17423` | `17` (NC_QUEST) | `15` | `sp_NC_QUEST_SELECT_START_REQ` | `PROTO_NC_QUEST_SELECT_START_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 135 | `0x4411` | `17425` | `17` (NC_QUEST) | `17` | `sp_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD` | `PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 136 | `0x4414` | `17428` | `17` (NC_QUEST) | `20` | `sp_NC_QUEST_START_REQ` | `PROTO_NC_QUEST_START_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 137 | `0x4416` | `17430` | `17` (NC_QUEST) | `22` | `sp_NC_QUEST_READ_REQ` | `PROTO_NC_QUEST_READ_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 138 | `0x4811` | `18449` | `18` (NC_SKILL) | `17` | `sp_NC_SKILL_EMPOWALLOC_REQ` | `PROTO_NC_SKILL_EMPOWALLOC_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 139 | `0x4822` | `18466` | `18` (NC_SKILL) | `34` | `sp_NC_SKILL_PRODUCTFIELD_REQ` | `PROTO_NC_SKILL_PRODUCTFIELD_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 140 | `0x4824` | `18468` | `18` (NC_SKILL) | `36` | `sp_NC_SKILL_UNLEARN_REQ` | `PROTO_NC_SKILL_UNLEARN_REQ` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 141 | `0x482A` | `18474` | `18` (NC_SKILL) | `42` | `sp_NC_SKILL_REPLYREVIVE_CMD` | `PROTO_NC_SKILL_REPLYREVIVE_CMD` | **Faltante** | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 142 | `0x4C01` | `19457` | `19` (NC_TRADE) | `1` | `sp_NC_TRADE_PROPOSE_REQ` | `PROTO_NC_TRADE_PROPOSE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 143 | `0x4C03` | `19459` | `19` (NC_TRADE) | `3` | `sp_NC_TRADE_PROPOSE_ASKNO_ACK` | `PROTO_NC_TRADE_PROPOSE_ASKNO_ACK` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 144 | `0x4C06` | `19462` | `19` (NC_TRADE) | `6` | `sp_NC_TRADE_PROPOSEYES_ACK` | `PROTO_NC_TRADE_PROPOSEYES_ACK` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 145 | `0x4C07` | `19463` | `19` (NC_TRADE) | `7` | `sp_NC_TRADE_PROPOSE_CANCEL_CMD` | `PROTO_NC_TRADE_PROPOSE_CANCEL_CMD` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 146 | `0x4C0A` | `19466` | `19` (NC_TRADE) | `10` | `sp_NC_TRADE_CANCEL_REQ` | `PROTO_NC_TRADE_CANCEL_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 147 | `0x4C0D` | `19469` | `19` (NC_TRADE) | `13` | `sp_NC_TRADE_UPBOARD_REQ` | `PROTO_NC_TRADE_UPBOARD_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 148 | `0x4C11` | `19473` | `19` (NC_TRADE) | `17` | `sp_NC_TRADE_DOWNBOARD_REQ` | `PROTO_NC_TRADE_DOWNBOARD_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 149 | `0x4C15` | `19477` | `19` (NC_TRADE) | `21` | `sp_NC_TRADE_CENBOARDING_REQ` | `PROTO_NC_TRADE_CENBOARDING_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 150 | `0x4C19` | `19481` | `19` (NC_TRADE) | `25` | `sp_NC_TRADE_BOARDLOCK_REQ` | `PROTO_NC_TRADE_BOARDLOCK_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 151 | `0x4C1F` | `19487` | `19` (NC_TRADE) | `31` | `sp_NC_TRADE_DECIDE_REQ` | `PROTO_NC_TRADE_DECIDE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 152 | `0x5001` | `20481` | `20` (NC_SOULSTONE) | `1` | `sp_NC_SOULSTONE_HP_BUY_REQ` | `PROTO_NC_SOULSTONE_HP_BUY_REQ` | Parcial / Base | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 153 | `0x5002` | `20482` | `20` (NC_SOULSTONE) | `2` | `sp_NC_SOULSTONE_SP_BUY_REQ` | `PROTO_NC_SOULSTONE_SP_BUY_REQ` | Parcial / Base | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 154 | `0x5007` | `20487` | `20` (NC_SOULSTONE) | `7` | `sp_NC_SOULSTONE_HP_USE_REQ` | `PROTO_NC_SOULSTONE_HP_USE_REQ` | Parcial / Base | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 155 | `0x5009` | `20489` | `20` (NC_SOULSTONE) | `9` | `sp_NC_SOULSTONE_SP_USE_REQ` | `PROTO_NC_SOULSTONE_SP_USE_REQ` | Parcial / Base | Fase 12 (Combate, Skills, AbState & SoulStone) |
| 156 | `0x5421` | `21537` | `21` (NC_FRIEND) | `33` | `sp_NC_FRIEND_UES_FRIEND_POINT_REQ` | `PROTO_NC_FRIEND_UES_FRIEND_POINT_REQ` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 157 | `0x581A` | `22554` | `22` (NC_KQ) | `26` | `sp_NC_KQ_ENTRYRESPONCE_ACK` | `PROTO_NC_KQ_ENTRYRESPONCE_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 158 | `0x5C01` | `23553` | `23` (NC_WT) | `1` | `sp_NC_WT_LICENSE_REQ` | `PROTO_NC_WT_LICENSE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 159 | `0x6001` | `24577` | `24` (NC_CT) | `1` | `sp_NC_CT_SET_CURRENT_REQ` | `PROTO_NC_CT_SET_CURRENT_REQ` | **Faltante** | Fase 13 (Drops, AutoPick, Títulos & Cartas) |
| 160 | `0x6801` | `26625` | `26` (NC_BOOTH) | `1` | `sp_NC_BOOTH_OPEN_REQ` | `PROTO_NC_BOOTH_OPEN_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 161 | `0x6804` | `26628` | `26` (NC_BOOTH) | `4` | `sp_NC_BOOTH_CLOSE_REQ` | `PROTO_NC_BOOTH_CLOSE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 162 | `0x6807` | `26631` | `26` (NC_BOOTH) | `7` | `sp_NC_BOOTH_ENTRY_REQ` | `PROTO_NC_BOOTH_ENTRY_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 163 | `0x680A` | `26634` | `26` (NC_BOOTH) | `10` | `sp_NC_BOOTH_REFRESH_REQ` | `PROTO_NC_BOOTH_REFRESH_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 164 | `0x680D` | `26637` | `26` (NC_BOOTH) | `13` | `sp_NC_BOOTH_ITEMTRADE_REQ` | `PROTO_NC_BOOTH_ITEMTRADE_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 165 | `0x6810` | `26640` | `26` (NC_BOOTH) | `16` | `sp_NC_BOOTH_INTERIORSTART_REQ` | `PROTO_NC_BOOTH_INTERIORSTART_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 166 | `0x6813` | `26643` | `26` (NC_BOOTH) | `19` | `sp_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ` | `PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 167 | `0x6815` | `26645` | `26` (NC_BOOTH) | `21` | `sp_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ` | `PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 168 | `0x6C06` | `27654` | `27` (NC_SCENARIO) | `6` | `sp_NC_SCENARIO_AREAENTRY_ACK` | `PROTO_NC_SCENARIO_AREAENTRY_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 169 | `0x6C08` | `27656` | `27` (NC_SCENARIO) | `8` | `sp_NC_SCENARIO_AREALEAVE_ACK` | `PROTO_NC_SCENARIO_AREALEAVE_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 170 | `0x7491` | `29841` | `29` (NC_GUILD) | `145` | `sp_NC_GUILD_STORAGEOPEN_REQ` | `PROTO_NC_GUILD_STORAGEOPEN_REQ` | **Faltante** | Fase 16 (Coordinación S2S Gremios & Almacén) |
| 171 | `0x749C` | `29852` | `29` (NC_GUILD) | `156` | `sp_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ` | `PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 172 | `0x74A0` | `29856` | `29` (NC_GUILD) | `160` | `sp_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD` | `PROTO_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD` | **Faltante** | Fase 16 (Coordinación S2S Gremios & Almacén) |
| 173 | `0x74AC` | `29868` | `29` (NC_GUILD) | `172` | `sp_NC_GUILD_RENAME_REQ` | `PROTO_NC_GUILD_RENAME_REQ` | **Faltante** | Fase 16 (Coordinación S2S Gremios & Almacén) |
| 174 | `0x74B4` | `29876` | `29` (NC_GUILD) | `180` | `sp_NC_GUILD_RETYPE_REQ` | `PROTO_NC_GUILD_RETYPE_REQ` | **Faltante** | Fase 16 (Coordinación S2S Gremios & Almacén) |
| 175 | `0x74CA` | `29898` | `29` (NC_GUILD) | `202` | `sp_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ` | `PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 176 | `0x74CF` | `29903` | `29` (NC_GUILD) | `207` | `sp_NC_GUILD_TOURNAMENT_USESKILL_REQ` | `PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 177 | `0x74D3` | `29907` | `29` (NC_GUILD) | `211` | `sp_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ` | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 178 | `0x74D7` | `29911` | `29` (NC_GUILD) | `215` | `sp_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ` | `PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 179 | `0x74DB` | `29915` | `29` (NC_GUILD) | `219` | `sp_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ` | `PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ` | **Faltante** | Fase 15 (Guild Tournament en Mapa) |
| 180 | `0x8C01` | `35841` | `35` (NC_MINIHOUSE) | `1` | `sp_NC_MINIHOUSE_ACTIV_REQ` | `PROTO_NC_MINIHOUSE_ACTIV_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 181 | `0x8C03` | `35843` | `35` (NC_MINIHOUSE) | `3` | `sp_NC_MINIHOUSE_BUILDING_REQ` | `PROTO_NC_MINIHOUSE_BUILDING_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 182 | `0x8C0A` | `35850` | `35` (NC_MINIHOUSE) | `10` | `sp_NC_MINIHOUSE_LOGINCOMPLETE_REQ` | `PROTO_NC_MINIHOUSE_LOGINCOMPLETE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 183 | `0x8C0E` | `35854` | `35` (NC_MINIHOUSE) | `14` | `sp_NC_MINIHOUSE_VISITREADY_REQ` | `PROTO_NC_MINIHOUSE_VISITREADY_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 184 | `0x8C10` | `35856` | `35` (NC_MINIHOUSE) | `16` | `sp_NC_MINIHOUSE_VISIT_REQ` | `PROTO_NC_MINIHOUSE_VISIT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 185 | `0x8C12` | `35858` | `35` (NC_MINIHOUSE) | `18` | `sp_NC_MINIHOUSE_KICKOUT_REQ` | `PROTO_NC_MINIHOUSE_KICKOUT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 186 | `0x8C18` | `35864` | `35` (NC_MINIHOUSE) | `24` | `sp_NC_MINIHOUSE_EXIT_REQ` | `PROTO_NC_MINIHOUSE_EXIT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 187 | `0x8C1D` | `35869` | `35` (NC_MINIHOUSE) | `29` | `sp_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ` | `PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 188 | `0x8C1F` | `35871` | `35` (NC_MINIHOUSE) | `31` | `sp_NC_MINIHOUSE_REARRANGE_REQ` | `PROTO_NC_MINIHOUSE_REARRANGE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 189 | `0x8C22` | `35874` | `35` (NC_MINIHOUSE) | `34` | `sp_NC_MINIHOUSE_CREATE_FURNITURE_REQ` | `PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 190 | `0x8C25` | `35877` | `35` (NC_MINIHOUSE) | `37` | `sp_NC_MINIHOUSE_DELETE_FURNITURE_REQ` | `PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 191 | `0x8C2B` | `35883` | `35` (NC_MINIHOUSE) | `43` | `sp_NC_MINIHOUSE_MODIFY_PASSWORD_REQ` | `PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 192 | `0x8C2E` | `35886` | `35` (NC_MINIHOUSE) | `46` | `sp_NC_MINIHOUSE_MODIFY_TITLE_REQ` | `PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 193 | `0x8C32` | `35890` | `35` (NC_MINIHOUSE) | `50` | `sp_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ` | `PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 194 | `0x8C3D` | `35901` | `35` (NC_MINIHOUSE) | `61` | `sp_NC_MINIHOUSE_CHAR_ACTION_REQ` | `PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 195 | `0x8C40` | `35904` | `35` (NC_MINIHOUSE) | `64` | `sp_NC_MINIHOUSE_FURNITURE_EFFECT_REQ` | `PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 196 | `0x8C43` | `35907` | `35` (NC_MINIHOUSE) | `67` | `sp_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ` | `PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 197 | `0x8C48` | `35912` | `35` (NC_MINIHOUSE) | `72` | `sp_NC_MINIHOUSE_MODIFY_NOTIFY_REQ` | `PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 198 | `0x8C50` | `35920` | `35` (NC_MINIHOUSE) | `80` | `sp_NC_MINIHOUSE_PORTAL_OPEN_REQ` | `PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 199 | `0x8C53` | `35923` | `35` (NC_MINIHOUSE) | `83` | `sp_NC_MINIHOUSE_PORTAL_CLOSE_REQ` | `PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 200 | `0x8C56` | `35926` | `35` (NC_MINIHOUSE) | `86` | `sp_NC_MINIHOUSE_PORTAL_EFFECT_REQ` | `PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 201 | `0x8C62` | `35938` | `35` (NC_MINIHOUSE) | `98` | `sp_NC_MINIHOUSE_FURNITURE_EMOTION_REQ` | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 202 | `0x8C65` | `35941` | `35` (NC_MINIHOUSE) | `101` | `sp_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ` | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 203 | `0x943C` | `37948` | `37` (NC_HOLY_PROMISE) | `60` | `sp_NC_HOLY_PROMISE_GET_CEN_REWARD_REQ` | `PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 204 | `0x9440` | `37952` | `37` (NC_HOLY_PROMISE) | `64` | `sp_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_REQ` | `PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 205 | `0x986D` | `39021` | `38` (NC_GUILD_ACADEMY) | `109` | `sp_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ` | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 206 | `0x986F` | `39023` | `38` (NC_GUILD_ACADEMY) | `111` | `sp_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD` | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD` | **Faltante** | Fase 14 (Inventario, Mercaderes, Trade, Booth & Forja) |
| 207 | `0xA40D` | `41997` | `41` (NC_INSTANCE_DUNGEON) | `13` | `sp_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK` | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 208 | `0xA417` | `42007` | `41` (NC_INSTANCE_DUNGEON) | `23` | `sp_NC_INSTANCE_DUNGEON_RANK_LIST_REQ` | `PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 209 | `0xAC01` | `44033` | `43` (NC_DICE) | `1` | `sp_NC_DICE_TAISAI_GAME_JOIN_REQ` | `PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 210 | `0xAC05` | `44037` | `43` (NC_DICE) | `5` | `sp_NC_DICE_TAISAI_GAME_LEAVE_REQ` | `PROTO_NC_DICE_TAISAI_GAME_LEAVE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 211 | `0xAC0B` | `44043` | `43` (NC_DICE) | `11` | `sp_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ` | `PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 212 | `0xAC0F` | `44047` | `43` (NC_DICE) | `15` | `sp_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ` | `PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 213 | `0xAC13` | `44051` | `43` (NC_DICE) | `19` | `sp_NC_DICE_TAISAI_GAME_MODE_CHANGE_REQ` | `PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 214 | `0xAC16` | `44054` | `43` (NC_DICE) | `22` | `sp_NC_DICE_TAISAI_RANK_REQ` | `PROTO_NC_DICE_TAISAI_RANK_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 215 | `0xAC1A` | `44058` | `43` (NC_DICE) | `26` | `sp_NC_DICE_TAISAI_BETTING_REQ` | `PROTO_NC_DICE_TAISAI_BETTING_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 216 | `0xAC1D` | `44061` | `43` (NC_DICE) | `29` | `sp_NC_DICE_TAISAI_BETTING_CANCEL_REQ` | `PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 217 | `0xAC28` | `44072` | `43` (NC_DICE) | `40` | `sp_NC_DICE_TAISAI_LARGE_AMOUNT_REQ` | `PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 218 | `0xB410` | `46096` | `45` (NC_USER_CONNECTION) | `16` | `sp_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ` | `PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ` | **Faltante** | Fase 16 (Coordinación S2S, Cuenta & Cierre 245/245) |
| 219 | `0xBC02` | `48130` | `47` (NC_GAMBLE) | `2` | `sp_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ` | `PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 220 | `0xBC0D` | `48141` | `47` (NC_GAMBLE) | `13` | `sp_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ` | `PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 221 | `0xBC0F` | `48143` | `47` (NC_GAMBLE) | `15` | `sp_NC_GAMBLE_COIN_BUY_REQ` | `PROTO_NC_GAMBLE_COIN_BUY_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 222 | `0xBC13` | `48147` | `47` (NC_GAMBLE) | `19` | `sp_NC_GAMBLE_COIN_SELL_REQ` | `PROTO_NC_GAMBLE_COIN_SELL_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 223 | `0xBC17` | `48151` | `47` (NC_GAMBLE) | `23` | `sp_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ` | `PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 224 | `0xBC64` | `48228` | `47` (NC_GAMBLE) | `100` | `sp_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 225 | `0xBC68` | `48232` | `47` (NC_GAMBLE) | `104` | `sp_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 226 | `0xBC6B` | `48235` | `47` (NC_GAMBLE) | `107` | `sp_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 227 | `0xBC6F` | `48239` | `47` (NC_GAMBLE) | `111` | `sp_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 228 | `0xBC71` | `48241` | `47` (NC_GAMBLE) | `113` | `sp_NC_GAMBLE_DICE_TAISAI_BETTING_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 229 | `0xBC74` | `48244` | `47` (NC_GAMBLE) | `116` | `sp_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ` | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 230 | `0xBCC8` | `48328` | `47` (NC_GAMBLE) | `200` | `sp_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 231 | `0xBCCA` | `48330` | `47` (NC_GAMBLE) | `202` | `sp_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 232 | `0xBCCC` | `48332` | `47` (NC_GAMBLE) | `204` | `sp_NC_GAMBLE_SLOTMACHINE_START_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 233 | `0xBCCE` | `48334` | `47` (NC_GAMBLE) | `206` | `sp_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 234 | `0xBCD0` | `48336` | `47` (NC_GAMBLE) | `208` | `sp_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 235 | `0xBCD9` | `48345` | `47` (NC_GAMBLE) | `217` | `sp_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 236 | `0xBCDB` | `48347` | `47` (NC_GAMBLE) | `219` | `sp_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ` | `PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 237 | `0xC401` | `50177` | `49` (NC_COLLECT) | `1` | `sp_NC_COLLECT_PEEPING_REQ` | `PROTO_NC_COLLECT_PEEPING_REQ` | **Faltante** | Fase 13 (Drops, AutoPick, Títulos & Cartas) |
| 238 | `0xC408` | `50184` | `49` (NC_COLLECT) | `8` | `sp_NC_COLLECT_CARDREGIST_REQ` | `PROTO_NC_COLLECT_CARDREGIST_REQ` | **Faltante** | Fase 13 (Drops, AutoPick, Títulos & Cartas) |
| 239 | `0xC40C` | `50188` | `49` (NC_COLLECT) | `12` | `sp_NC_COLLECT_BOOKMARK_REGIST_REQ` | `PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ` | **Faltante** | Fase 13 (Drops, AutoPick, Títulos & Cartas) |
| 240 | `0xC410` | `50192` | `49` (NC_COLLECT) | `16` | `sp_NC_COLLECT_REWARD_REQ` | `PROTO_NC_COLLECT_REWARD_REQ` | **Faltante** | Fase 13 (Drops, AutoPick, Títulos & Cartas) |
| 241 | `0xCC01` | `52225` | `51` (NC_MOVER) | `1` | `sp_NC_MOVER_RIDE_ON_REQ` | `PROTO_NC_MOVER_RIDE_ON_REQ` | **Faltante** | Fase 11 (Monturas MOVER) |
| 242 | `0xCC05` | `52229` | `51` (NC_MOVER) | `5` | `sp_NC_MOVER_RIDE_OFF_REQ` | `PROTO_NC_MOVER_RIDE_OFF_REQ` | **Faltante** | Fase 11 (Monturas MOVER) |
| 243 | `0xD406` | `54278` | `53` (NC_PET) | `6` | `sp_NC_PET_ASK_NEW_NAME_ACK` | `PROTO_NC_PET_ASK_NEW_NAME_ACK` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 244 | `0xD407` | `54279` | `53` (NC_PET) | `7` | `sp_NC_PET_SET_NAME_REQ` | `PROTO_NC_PET_SET_NAME_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |
| 245 | `0xD409` | `54281` | `53` (NC_PET) | `9` | `sp_NC_PET_SET_NAME_CANCEL_REQ` | `PROTO_NC_PET_SET_NAME_CANCEL_REQ` | **Faltante** | Fase 15 (Quests, Lua, Instancias, MiniHouse, Pets & Casino) |

### 4.2 Los 62 Métodos `iac_*` (`ItemAttributeClass`) y los 18 Símbolos `sp_NC_*` Auxiliares en `Zone.pdb`

Además de las 245 entradas directas de `protocolstore(ShinePlayer)`, `Zone.pdb` expone:

1. **Los `62` métodos `iac_*` (`ItemAttributeClass`)** encargados de codificar/decodificar atributos binarios de ítems, durabilidad, mejoras, mascotas invocadas, monturas, cofres, títulos de arma (`wptitle`) y consumibles (`iac_FoodEat`):

| # | Método `ItemAttributeClass` (`Zone.pdb`) | # | Método `ItemAttributeClass` (`Zone.pdb`) |
| :---: | :--- | :---: | :--- |
| 1 | `iac_ClearItemOptionUpgrade` | 32 | `iac_SetDeteleDateTime` |
| 2 | `iac_DropItemMake` | 33 | `iac_SetDuringPetSummon` |
| 3 | `iac_EndOfOption` | 34 | `iac_SetDuringRide` |
| 4 | `iac_FindOption` | 35 | `iac_SetEndTime` |
| 5 | `iac_FoodEat` | 36 | `iac_SetEndureEndTime` |
| 6 | `iac_GetAmount` | 37 | `iac_SetEndureGrade` |
| 7 | `iac_GetAttrSize` | 38 | `iac_SetFurnitureKey` |
| 8 | `iac_GetChargedNumber` | 39 | `iac_SetGrade` |
| 9 | `iac_GetDeteleDateTime` | 40 | `iac_SetHP` |
| 10 | `iac_GetDuringPetSummon` | 41 | `iac_SetIsPutOnBelonged` |
| 11 | `iac_GetDuringRide` | 42 | `iac_SetItemOptionUpgrade` |
| 12 | `iac_GetEndureEndTime` | 43 | `iac_SetLot` |
| 13 | `iac_GetGrade` | 44 | `iac_SetPetID` |
| 14 | `iac_GetHP` | 45 | `iac_SetPetName` |
| 15 | `iac_GetHungry` | 46 | `iac_SetPetRegistNumber` |
| 16 | `iac_GetIsPutOnBelonged` | 47 | `iac_SetRandomOptionChangeCount` |
| 17 | `iac_GetItemOptionStruct` | 48 | `iac_SetRareFailCount` |
| 18 | `iac_GetItemOptionUpgrade` | 49 | `iac_SetUseableDateTime` |
| 19 | `iac_GetLot` | 50 | `iac_Setoption` |
| 20 | `iac_GetPetID` | 51 | `iac_StepGet` |
| 21 | `iac_GetPetName` | 52 | `iac_charge` |
| 22 | `iac_GetPetRegistNumber` | 53 | `iac_equipbriefinfo` |
| 23 | `iac_GetRandomOptionChangeCount` | 54 | `iac_fill2buffer` |
| 24 | `iac_GetRareFailCount` | 55 | `iac_fill2struct` |
| 25 | `iac_GetUseableDateTime` | 56 | `iac_islotitem` |
| 26 | `iac_Getitemchest` | 57 | `iac_itemcreate` |
| 27 | `iac_Getupgrade` | 58 | `iac_wptitle_Getlicencemob` |
| 28 | `iac_InvenStore` | 59 | `iac_wptitle_Gettitle` |
| 29 | `iac_ItemClass` | 60 | `iac_wptitle_Setlicense` |
| 30 | `iac_SetAmount` | 61 | `iac_wptitle_Settitle` |
| 31 | `iac_SetChargedNumber` | 62 | `iac_wptitle_mobkilladd` |

2. **Los `18` handlers `sp_NC_*` auxiliares presentes en los símbolos de `Zone.pdb`** (latidos de conexión, casteo desde montura, reinicio pago de stats/empow, promoción de clase y subcomandos de MiniHouse/Torneo):
   - `sp_NC_MISC_HEARTBEAT_REQ`, `sp_NC_MISC_HEARTBEAT_ACK`
   - `sp_NC_BAT_SKILLBASH_CASTABORT_CMD`, `sp_NC_MOVER_SKILLBASH_FLD_CAST_REQ`, `sp_NC_MOVER_SKILLBASH_OBJ_CAST_REQ`
   - `sp_NC_CHAR_PROMOTE_REQ`, `sp_NC_CHARGED_STAT_INITIALIZE_REQ`, `sp_NC_CHARGED_SKILLEMPOW_INITIALIZE_REQ`, `sp_NC_SKILL_EMPOW_RESET_REQ`
   - `sp_NC_DICE_TAISAI_DICE_ROLL_REQ`
   - `sp_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK`, `sp_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK`, `sp_NC_GUILD_TOURNAMENT_FLAGCAPTURE_CMD`, `sp_NC_GUILD_TOURNAMENT_USESKILL_ACK`
   - `sp_NC_MINIHOUSE_ARRANGEMODE_REQ`, `sp_NC_MINIHOUSE_COMPULSIONMOVETO_REQ`, `sp_NC_MINIHOUSE_MODIFY_OWNERBLOG_REQ`, `sp_NC_MINIHOUSE_OWNERBLOG_REQ`

### 4.3 Catálogo 100% Exhaustivo de los 138 Handlers S2S `wms_NC_*` (`WorldManagerSession` — `Zone <-> WorldService`)

Listado completo 1:1 de los **138 símbolos `wms_NC_*`** de `Zone.pdb` mediante los cuales `ZoneService` coordina eventos inter-zona y globales con `WorldService` (implementados en Go sobre canales tipados en modo unificado y bus/tabla MySQL en modo multi-proceso en **Fase 16**):

| # | Handler S2S `wms_NC_*` (`Zone.pdb`) | # | Handler S2S `wms_NC_*` (`Zone.pdb`) |
| :---: | :--- | :---: | :--- |
| 1 | `wms_NC_CHAR_GET_CHAT_BLOCK_SPAME` | 70 | `wms_NC_GUILD_ZONE_MAKE_CMD` |
| 2 | `wms_NC_CHAR_GET_CHAT_BLOCK_SPAMER_CMD` | 71 | `wms_NC_GUILD_ZONE_MEMBER_JOIN_CMD` |
| 3 | `wms_NC_CHAR_GUILD_ACADEMY_ZONE_CMD` | 72 | `wms_NC_GUILD_ZONE_MEMBER_LEAVE_CMD` |
| 4 | `wms_NC_CHAR_ITEMACTIONCOOLTIME_ZONE_CMD` | 73 | `wms_NC_GUILD_ZONE_WAR_END_CMD` |
| 5 | `wms_NC_CHAR_REGISTNUMBER_ACK` | 74 | `wms_NC_GUILD_ZONE_WAR_START_CMD` |
| 6 | `wms_NC_CHAR_SELL_ITEM_INFO_ZONE_CMD` | 75 | `wms_NC_HOLY_PROMISE_MYUPPER_ACK` |
| 7 | `wms_NC_CHAR_SUPPORT_REWARD_CMD` | 76 | `wms_NC_HOLY_PROMISE_MY_UP_ZONE` |
| 8 | `wms_NC_CHAR_TUTORIAL_INFO_ZONE_CMD` | 77 | `wms_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD` |
| 9 | `wms_NC_CHAR_USEITEM_MINIMON_INFO_ZONE_CMD` | 78 | `wms_NC_INSTANCE_DUNGEON_LEVEL_SE` |
| 10 | `wms_NC_CHAR_WEDDING_CANCEL_WEDDING` | 79 | `wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK` |
| 11 | `wms_NC_CHAR_WORLD_AUTO_PICK_CMD` | 80 | `wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ` |
| 12 | `wms_NC_CHAR_WORLD_REST_EXP_LAST_EXEC_TIME_CMD` | 81 | `wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK` |
| 13 | `wms_NC_CHAR_ZONECHARABSTATE_CMD` | 82 | `wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ` |
| 14 | `wms_NC_CHAR_ZONECHARDATA_ACK` | 83 | `wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD` |
| 15 | `wms_NC_CHAR_ZONECHARQUEST_CMD` | 84 | `wms_NC_KQ_PLAYER_DISJOIN_CMD` |
| 16 | `wms_NC_CHAR_ZONECHARSKILL_CMD` | 85 | `wms_NC_KQ_W2Z_DESTROY_CMD` |
| 17 | `wms_NC_CHAR_ZONE_ABSTATE_CMD` | 86 | `wms_NC_KQ_W2Z_MAKE_REQ` |
| 18 | `wms_NC_CHAR_ZONE_BASE_CMD` | 87 | `wms_NC_KQ_W2Z_START_CMD` |
| 19 | `wms_NC_CHAR_ZONE_CARDCOLLECT_BOOKMARK_CMD` | 88 | `wms_NC_MAP_LINKRESERVE_ACK` |
| 20 | `wms_NC_CHAR_ZONE_CARDCOLLECT_CMD` | 89 | `wms_NC_MISC_CONNECTFROMWHERE_ACK` |
| 21 | `wms_NC_CHAR_ZONE_CARDCOLLECT_REWARD_CMD` | 90 | `wms_NC_MISC_EVENTNPC_STANDEND_ZO` |
| 22 | `wms_NC_CHAR_ZONE_CHARDATAFAIL_ACK` | 91 | `wms_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD` |
| 23 | `wms_NC_CHAR_ZONE_CHARDATA_ACK` | 92 | `wms_NC_MISC_EVENTNPC_STANDSTART_` |
| 24 | `wms_NC_CHAR_ZONE_CHARGEDBUFF_CMD` | 93 | `wms_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD` |
| 25 | `wms_NC_CHAR_ZONE_CHARTITLE_CMD` | 94 | `wms_NC_MISC_HEARTBEAT_ACK` |
| 26 | `wms_NC_CHAR_ZONE_CHAT_COLOR_CMD` | 95 | `wms_NC_MISC_HEARTBEAT_REQ` |
| 27 | `wms_NC_CHAR_ZONE_CHESTINFO_CMD` | 96 | `wms_NC_MISC_PINGTEST_TOOL_WM_ZONE` |
| 28 | `wms_NC_CHAR_ZONE_COININFO_CMD` | 97 | `wms_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB` |
| 29 | `wms_NC_CHAR_ZONE_GAME_CMD` | 98 | `wms_NC_MISC_S2SCONNECTION_ACK` |
| 30 | `wms_NC_CHAR_ZONE_GUILD_CMD` | 99 | `wms_NC_MISC_S2SCONNECTION_RDY` |
| 31 | `wms_NC_CHAR_ZONE_ITEM_CMD` | 100 | `wms_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD` |
| 32 | `wms_NC_CHAR_ZONE_LINK_FROM_CMD` | 101 | `wms_NC_MISC_START_THE_BOOM_CMD` |
| 33 | `wms_NC_CHAR_ZONE_OPTION_CMD` | 102 | `wms_NC_MISC_TIMEFROMWORLD_CMD` |
| 34 | `wms_NC_CHAR_ZONE_PASSIVE_CMD` | 103 | `wms_NC_MISC_WHSHANDLEREPAIR_CMD` |
| 35 | `wms_NC_CHAR_ZONE_QUEST_DOING_CMD` | 104 | `wms_NC_MOVER_LINKDATA_ZONE_CMD` |
| 36 | `wms_NC_CHAR_ZONE_QUEST_DONE_CMD` | 105 | `wms_NC_PARTY_CHANGEMASTER_CMD` |
| 37 | `wms_NC_CHAR_ZONE_QUEST_READ_CMD` | 106 | `wms_NC_PARTY_DISMISS_CMD` |
| 38 | `wms_NC_CHAR_ZONE_QUEST_REPEAT_CMD` | 107 | `wms_NC_PARTY_FUNDAMENTAL_CMD` |
| 39 | `wms_NC_CHAR_ZONE_SHAPE_CMD` | 108 | `wms_NC_PARTY_ITEM_LOOTING_ZONE_C` |
| 40 | `wms_NC_CHAR_ZONE_SINGLE_OPTION_CMD` | 109 | `wms_NC_PARTY_ITEM_LOOTING_ZONE_CMD` |
| 41 | `wms_NC_CHAR_ZONE_SKILL_CMD` | 110 | `wms_NC_PARTY_JOIN_CMD` |
| 42 | `wms_NC_CT_ADD_FRIEND_CMD` | 111 | `wms_NC_PARTY_KICKOFF_CMD` |
| 43 | `wms_NC_CT_CHARTTING_CMD` | 112 | `wms_NC_PARTY_LEAVE_CMD` |
| 44 | `wms_NC_CT_LUASCRIPT_SET_ZONE_CMD` | 113 | `wms_NC_PARTY_LOGIN_CMD` |
| 45 | `wms_NC_EVENT_ADD_UPDATE_EVENT_CMD` | 114 | `wms_NC_PARTY_LOGOUT_CMD` |
| 46 | `wms_NC_EVENT_DEL_EVENT_CMD` | 115 | `wms_NC_PARTY_MEMBERINFORM_ACK` |
| 47 | `wms_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD` | 116 | `wms_NC_PARTY_ZONE_SET_LOOTER_CMD` |
| 48 | `wms_NC_GAMBLE_ZONE_PREVMAPNAME_CMD` | 117 | `wms_NC_PET_LINK_RESUMMON_CMD` |
| 49 | `wms_NC_GUILD_ACADEMY_REWARDSTORAGE_GRADE_INFO_CMD` | 118 | `wms_NC_QUEST_RESET_TIME_CMD` |
| 50 | `wms_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_ACK` | 119 | `wms_NC_QUEST_RESET_TIME_ZONE_CMD` |
| 51 | `wms_NC_GUILD_ACADEMY_ZONE_GUILD_ACADEMY_MASTER_BUFF_CMD` | 120 | `wms_NC_RAID_ZONE_DEL_SUB_LEADER_CMD` |
| 52 | `wms_NC_GUILD_ACADEMY_ZONE_GUILD_BUFF_CMD` | 121 | `wms_NC_RAID_ZONE_DISMISS_CMD` |
| 53 | `wms_NC_GUILD_ACADEMY_ZONE_MASTER_TELEPORT_CMD` | 122 | `wms_NC_RAID_ZONE_LEAVE_CMD` |
| 54 | `wms_NC_GUILD_ACADEMY_ZONE_MEMBER_GUILD_JOIN_CMD` | 123 | `wms_NC_RAID_ZONE_LOGIN_CMD` |
| 55 | `wms_NC_GUILD_ACADEMY_ZONE_MEMBER_JOIN_CMD` | 124 | `wms_NC_RAID_ZONE_LOGOUT_CMD` |
| 56 | `wms_NC_GUILD_ACADEMY_ZONE_MEMBER_LEAVE_CMD` | 125 | `wms_NC_RAID_ZONE_LOOTING_TYPE_CMD` |
| 57 | `wms_NC_GUILD_GUILDWARCONFIRM_ACK` | 126 | `wms_NC_RAID_ZONE_MAKE_CMD` |
| 58 | `wms_NC_GUILD_GUILDWARSTATUS_ACK` | 127 | `wms_NC_RAID_ZONE_MEMBER_JOIN_CMD` |
| 59 | `wms_NC_GUILD_STORAGE_WITHDRAW_GRADE_ACK` | 128 | `wms_NC_RAID_ZONE_SET_LEADER_CMD` |
| 60 | `wms_NC_GUILD_TOURNAMENT_GOLD_REFUND_ZONE_CMD` | 129 | `wms_NC_RAID_ZONE_SET_LOOTER_CMD` |
| 61 | `wms_NC_GUILD_TOURNAMENT_START_CMD` | 130 | `wms_NC_RAID_ZONE_SET_SUB_LEADER_CMD` |
| 62 | `wms_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_REQ` | 131 | `wms_NC_RAID_ZONE_SLOT_MOVE_CMD` |
| 63 | `wms_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_ACK` | 132 | `wms_NC_RAID_ZONE_VANISH_CMD` |
| 64 | `wms_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_ACK` | 133 | `wms_NC_USER_CONNECTCUT2ZONE_CMD` |
| 65 | `wms_NC_GUILD_WORLD_RENAME_ACK` | 134 | `wms_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD` |
| 66 | `wms_NC_GUILD_WORLD_RETYPE_ACK` | 135 | `wms_NC_USER_GET_PLAY_MIN_ACK` |
| 67 | `wms_NC_GUILD_WORLD_USE_GUILD_TOKEN_CMD` | 136 | `wms_NC_USER_KICKOFFFROMWORLD_CMD` |
| 68 | `wms_NC_GUILD_ZONE_DELETE_CMD` | 137 | `wms_NC_ZONE_PRISON_END_CMD` |
| 69 | `wms_NC_GUILD_ZONE_LIST_ACK` | 138 | `wms_NC_ZONE_PRISON_GO_REQ` |

### 4.4 Catálogo 100% Exhaustivo de los 196 Handlers S2S `gds_NC_*` (`GameDBSession` — `Zone <-> Character DB`)

En el legado, `Zone.exe` recibía las respuestas asíncronas de `Character.exe` mediante estos **196 handlers `gds_NC_*`**. En `XelionCore`, cada uno de estos 196 flujos se resuelve invocando directamente los Stored Procedures correspondientes de `CharacterRepository` (`425` SPs de `World00_Character` en `xelion_online`) sin latencia de socket intermedio:

| # | Handler S2S `gds_NC_*` (`Zone.pdb`) | # | Handler S2S `gds_NC_*` (`Zone.pdb`) |
| :---: | :--- | :---: | :--- |
| 1 | `gds_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_FAIL_ACK` | 99 | `gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK` |
| 2 | `gds_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_SUC_ACK` | 100 | `gds_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_ACK` |
| 3 | `gds_NC_CHARGED_STAT_INITIALIZE_DB_FAIL_ACK` | 101 | `gds_NC_ITEMDB_GUILD_TOKEN_BUYALLSUC_ACK` |
| 4 | `gds_NC_CHARGED_STAT_INITIALIZE_DB_SUC_ACK` | 102 | `gds_NC_ITEMDB_GUILD_TOKEN_BUYFAIL_ACK` |
| 5 | `gds_NC_CHARSAVE_ABSTATE_ACK` | 103 | `gds_NC_ITEMDB_GUILD_TOKEN_BUYLOTSUC_ACK` |
| 6 | `gds_NC_CHARSAVE_ALL_ACK` | 104 | `gds_NC_ITEMDB_INC_DEC_MONEYFAIL_ACK` |
| 7 | `gds_NC_CHARSAVE_AUTO_PICK_SAVE_ACK` | 105 | `gds_NC_ITEMDB_INC_DEC_MONEYSUC_ACK` |
| 8 | `gds_NC_CHARSAVE_DB_UI_STATE_SAVE_ACK` | 106 | `gds_NC_ITEMDB_ITEMBREAKFAIL_ACK` |
| 9 | `gds_NC_CHARSAVE_REST_EXP_LAST_EXEC_TIME_SAVE_ACK` | 107 | `gds_NC_ITEMDB_ITEMBREAKSUC_ACK` |
| 10 | `gds_NC_CHARSAVE_SKILL_ACK` | 108 | `gds_NC_ITEMDB_ITEMCHANGE_ACK` |
| 11 | `gds_NC_CHAR_CLASSCHANGE_ACK` | 109 | `gds_NC_ITEMDB_ITEMINFO_ACK` |
| 12 | `gds_NC_CHAR_FREESTAT_SET_DB_ACK` | 110 | `gds_NC_ITEMDB_ITEMMONEY_BUY_ACK` |
| 13 | `gds_NC_CHAR_GET_ITEMLIST_BY_TYPE_ACK` | 111 | `gds_NC_ITEMDB_ITEMREBUILD_ACK` |
| 14 | `gds_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_ACK` | 112 | `gds_NC_ITEMDB_ITEMTOTALINFORM_ACK` |
| 15 | `gds_NC_CHAR_REGNUM_VARIFICATION_ACK` | 113 | `gds_NC_ITEMDB_ITEMTRADEFAIL_ACK` |
| 16 | `gds_NC_CHAR_STAT_DECPOINTFAIL_DB_ACK` | 114 | `gds_NC_ITEMDB_ITEMTRADESUC_ACK` |
| 17 | `gds_NC_CHAR_STAT_DECPOINTSUC_DB_ACK` | 115 | `gds_NC_ITEMDB_MAPLINK_ITEM_CONSUME_ACK` |
| 18 | `gds_NC_CHAR_STAT_INCPOINTFAIL_DB_ACK` | 116 | `gds_NC_ITEMDB_MERGEFAIL_ACK` |
| 19 | `gds_NC_CHAR_STAT_INCPOINTSUC_DB_ACK` | 117 | `gds_NC_ITEMDB_MERGESUC_ACK` |
| 20 | `gds_NC_CHAR_TUTORIAL_FREESTAT_INIT_ACK` | 118 | `gds_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_ACK` |
| 21 | `gds_NC_CHAR_TUTORIAL_STEP_SAVE_ACK` | 119 | `gds_NC_ITEMDB_MIX_ITEM_ACK` |
| 22 | `gds_NC_CHAR_WEDDINGDATA_ACK` | 120 | `gds_NC_ITEMDB_MOVER_RAREMOVER_ACK` |
| 23 | `gds_NC_CHAR_WEDDING_DIVORCE_BY_AGREE_ACK` | 121 | `gds_NC_ITEMDB_MOVER_UPGRADE_ACK` |
| 24 | `gds_NC_CHAR_WEDDING_DIVORCE_BY_FORCE_ACK` | 122 | `gds_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_ACK` |
| 25 | `gds_NC_CHAR_WEDDING_DIVORCE_CANCEL_ACK` | 123 | `gds_NC_ITEMDB_NEW_UPGRADE_AC` |
| 26 | `gds_NC_CHAR_WEDDING_DIVORCE_DO_ACK` | 124 | `gds_NC_ITEMDB_NEW_UPGRADE_ACK` |
| 27 | `gds_NC_CHAR_WEDDING_DO_ACK` | 125 | `gds_NC_ITEMDB_OPENSTORAGE_ACK` |
| 28 | `gds_NC_CHAR_WEDDING_ESCAPE_DIVORCE_ACK` | 126 | `gds_NC_ITEMDB_OPENSTORAGE_FAIL_ACK` |
| 29 | `gds_NC_CHAR_WEDDING_PARTNER_INFO_ACK` | 127 | `gds_NC_ITEMDB_OPEN_GUILD_STORAGE_ACK` |
| 30 | `gds_NC_CHAR_WEDDING_PROPOSE_ACK` | 128 | `gds_NC_ITEMDB_OPEN_GUILD_STORAGE_FAIL_ACK` |
| 31 | `gds_NC_CHAR_WEDDING_PROPOSE_CANCEL_ACK` | 129 | `gds_NC_ITEMDB_PICKALLFAIL_ACK` |
| 32 | `gds_NC_COLLECT_BOOKMARK_REGIST_DB_ACK` | 130 | `gds_NC_ITEMDB_PICKALLSUC_ACK` |
| 33 | `gds_NC_COLLECT_PEEPING_DB_COLLECT_ACK` | 131 | `gds_NC_ITEMDB_PICKMERGEFAIL_ACK` |
| 34 | `gds_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK` | 132 | `gds_NC_ITEMDB_PICKMERGESUC_ACK` |
| 35 | `gds_NC_COLLECT_REGIST_ACK` | 133 | `gds_NC_ITEMDB_PICKMONEYFAIL_ACK` |
| 36 | `gds_NC_COLLECT_REWARD_DB_ACK` | 134 | `gds_NC_ITEMDB_PICKMONEYSUC_ACK` |
| 37 | `gds_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK` | 135 | `gds_NC_ITEMDB_PRODUCE_ACK` |
| 38 | `gds_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK` | 136 | `gds_NC_ITEMDB_QUESTREWARD_ACK` |
| 39 | `gds_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK` | 137 | `gds_NC_ITEMDB_RANDOMOPTION_CHANGE_ACK` |
| 40 | `gds_NC_DICE_TAISAI_DB_RANK_ACK` | 138 | `gds_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_ACK` |
| 41 | `gds_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK` | 139 | `gds_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK` |
| 42 | `gds_NC_GAMBLE_COIN_DB_BUY_ACK` | 140 | `gds_NC_ITEMDB_REINFORCEUNEQUIPFAIL_ACK` |
| 43 | `gds_NC_GAMBLE_COIN_DB_SELL_ACK` | 141 | `gds_NC_ITEMDB_REINFORCEUNEQUIPSUC_ACK` |
| 44 | `gds_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK` | 142 | `gds_NC_ITEMDB_RELOCFAIL_ACK` |
| 45 | `gds_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK` | 143 | `gds_NC_ITEMDB_RELOCSUC_ACK` |
| 46 | `gds_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK` | 144 | `gds_NC_ITEMDB_REPURCHASE_ALL_ACK` |
| 47 | `gds_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK` | 145 | `gds_NC_ITEMDB_REPURCHASE_LOT_ACK` |
| 48 | `gds_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK` | 146 | `gds_NC_ITEMDB_SE` |
| 49 | `gds_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK` | 147 | `gds_NC_ITEMDB_SELLALLSUC_ACK` |
| 50 | `gds_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_ACK` | 148 | `gds_NC_ITEMDB_SELLFAIL_ACK` |
| 51 | `gds_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_ACK` | 149 | `gds_NC_ITEMDB_SELLLOTSUC_ACK` |
| 52 | `gds_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK` | 150 | `gds_NC_ITEMDB_SHIELDENDURE_CHARGE_ACK` |
| 53 | `gds_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK` | 151 | `gds_NC_ITEMDB_SOULSTONEBUYFAIL_ACK` |
| 54 | `gds_NC_ITEMDB_ADMINCREATEFAIL_ACK` | 152 | `gds_NC_ITEMDB_SOULSTONEBUYSUC_ACK` |
| 55 | `gds_NC_ITEMDB_ADMINCREATESUC_ACK` | 153 | `gds_NC_ITEMDB_SPLITFAIL_ACK` |
| 56 | `gds_NC_ITEMDB_BOOTHTRADE_ACK` | 154 | `gds_NC_ITEMDB_SPLITSUC_ACK` |
| 57 | `gds_NC_ITEMDB_BRACELET_UPGRADE_ACK` | 155 | `gds_NC_ITEMDB_SPLIT_N_MERGEFAIL_ACK` |
| 58 | `gds_NC_ITEMDB_BUYALLSUC_ACK` | 156 | `gds_NC_ITEMDB_SPLIT_N_MERGESUC_ACK` |
| 59 | `gds_NC_ITEMDB_BUYCAPSULE_ACK` | 157 | `gds_NC_ITEMDB_TERMEXTEND_ACK` |
| 60 | `gds_NC_ITEMDB_BUYFAIL_ACK` | 158 | `gds_NC_ITEMDB_UES_FRIEND_POINT_ACK` |
| 61 | `gds_NC_ITEMDB_BUYLOTSUC_ACK` | 159 | `gds_NC_ITEMDB_UNEQUIPFAIL_ACK` |
| 62 | `gds_NC_ITEMDB_CAPSULEITEM_ACK` | 160 | `gds_NC_ITEMDB_UNEQUIPSUC_ACK` |
| 63 | `gds_NC_ITEMDB_CHARGED_LIST_ACK` | 161 | `gds_NC_ITEMDB_UPGRADE_ACK` |
| 64 | `gds_NC_ITEMDB_CHARGED_WITHDRAW_ACK` | 162 | `gds_NC_ITEMDB_USE_ACK` |
| 65 | `gds_NC_ITEMDB_CHAT_COLOR_CHANGE_ACK` | 163 | `gds_NC_ITEMDB_WEAPONENDURE_CHARGE_ACK` |
| 66 | `gds_NC_ITEMDB_CHESTITEM_ACK` | 164 | `gds_NC_ITEMDB_WITHDRAW_ACK` |
| 67 | `gds_NC_ITEMDB_CLASS_CHANGE_ACK` | 165 | `gds_NC_KQ_REWARDFAIL_ACK` |
| 68 | `gds_NC_ITEMDB_CREATEITEMLISTFAIL_ACK` | 166 | `gds_NC_KQ_REWARDSUC_ACK` |
| 69 | `gds_NC_ITEMDB_CREATEITEMLISTSUC_ACK` | 167 | `gds_NC_MINIHOUSE_DB_PORTAL_ADD_ACK` |
| 70 | `gds_NC_ITEMDB_CREATEMUSHROOMFAIL_ACK` | 168 | `gds_NC_MINIHOUSE_DB_PORTAL_DEL_ACK` |
| 71 | `gds_NC_ITEMDB_CREATEMUSHROOMSUC_ACK` | 169 | `gds_NC_MINIHOUSE_DB_PORTAL_LIST_ACK` |
| 72 | `gds_NC_ITEMDB_DEPOSIT_ACK` | 170 | `gds_NC_MINIHOUSE_DB_VISITER_COUNT_ACK` |
| 73 | `gds_NC_ITEMDB_DESTROY_ACK` | 171 | `gds_NC_MISC_HEARTBEAT_ACK` |
| 74 | `gds_NC_ITEMDB_DISMANTLE_ACK` | 172 | `gds_NC_MISC_HEARTBEAT_REQ` |
| 75 | `gds_NC_ITEMDB_DROPALLFAIL_ACK` | 173 | `gds_NC_MISC_PINGTEST_CLIENT_ZONE_DB` |
| 76 | `gds_NC_ITEMDB_DROPALLLOT_ACK` | 174 | `gds_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB` |
| 77 | `gds_NC_ITEMDB_DROPALLSUC_ACK` | 175 | `gds_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB` |
| 78 | `gds_NC_ITEMDB_DROPLOTFAIL_ACK` | 176 | `gds_NC_MISC_S2SCONNECTION_ACK` |
| 79 | `gds_NC_ITEMDB_ENCHANT_ADD_GEM_AC` | 177 | `gds_NC_MISC_S2SCONNECTION_RDY` |
| 80 | `gds_NC_ITEMDB_ENCHANT_ADD_GEM_ACK` | 178 | `gds_NC_PET_CREATE_DB_ACK` |
| 81 | `gds_NC_ITEMDB_ENCHANT_ADD_NEW_SO` | 179 | `gds_NC_PET_LOAD_INFO_DB_ACK` |
| 82 | `gds_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK` | 180 | `gds_NC_PET_SET_NAME_DB_ACK` |
| 83 | `gds_NC_ITEMDB_ENCHANT_REMOVE_GEM` | 181 | `gds_NC_QUEST_DB_CLEAR_ACK` |
| 84 | `gds_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK` | 182 | `gds_NC_QUEST_DB_DONE_ACK` |
| 85 | `gds_NC_ITEMDB_ENCHANT_SET_GEM_LOT_ACK` | 183 | `gds_NC_QUEST_DB_GIVE_UP_ACK` |
| 86 | `gds_NC_ITEMDB_EQUIPFAIL_ACK` | 184 | `gds_NC_QUEST_DB_READ_ACK` |
| 87 | `gds_NC_ITEMDB_EQUIPSUC_ACK` | 185 | `gds_NC_QUEST_DB_SET_INFO_ACK` |
| 88 | `gds_NC_ITEMDB_EQUIP_BELONGED_CAN` | 186 | `gds_NC_SKILL_EMPOWALLOC_DB_ACK` |
| 89 | `gds_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_ACK` | 187 | `gds_NC_SKILL_EMPOW_RESET_DB_FAIL_ACK` |
| 90 | `gds_NC_ITEMDB_EXCHANGEFAIL_ACK` | 188 | `gds_NC_SKILL_EMPOW_RESET_DB_SUC_ACK` |
| 91 | `gds_NC_ITEMDB_EXCHANGESUC_ACK` | 189 | `gds_NC_SKILL_ERASE_ACK` |
| 92 | `gds_NC_ITEMDB_FURNITURE_ENDURE_ACK` | 190 | `gds_NC_SKILL_SKILLTEACHFAIL_ACK` |
| 93 | `gds_NC_ITEMDB_GETFROMCAPSULE_ACK` | 191 | `gds_NC_SKILL_SKILLTEACHSUC_ACK` |
| 94 | `gds_NC_ITEMDB_GETFROMCHESTFAIL_ACK` | 192 | `gds_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK` |
| 95 | `gds_NC_ITEMDB_GETFROMCHESTSUC_ACK` | 193 | `gds_NC_WT_LICENSE_CLR_ACK` |
| 96 | `gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK` | 194 | `gds_NC_WT_LICENSE_SET_ACK` |
| 97 | `gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK` | 195 | `gds_NC_WT_MOB_KILLCOUNT_SET_ACK` |
| 98 | `gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_FAIL_ACK` | 196 | `gds_NC_WT_TITLE_SET_ACK` |

---

## 5. Catálogo 100% Exhaustivo de las 1.343 Estructuras C++ de `Zone.pdb` (`docs/Servicios/Zone.md`)

Desglose exacto de las **1.343 estructuras** documentadas en [`docs/Servicios/Zone.md`](./Servicios/Zone.md):
- **`1.293` estructuras oficiales de protocolo `PROTO_NC_*`** (Sección 5.1 — implementación obligatoria 100% en `pkg/protocol/` con validación `sizeof` en **Fase 09**).
- **`9` plantillas `ShineSendType<PROTO_NC_*>`** (Sección 5.2 — helpers de serialización de `1 byte` sin payload adicional).
- **`1` estructura `OpTool` excluida** (`PROTO_NC_OPTOOL_CONNECT_BRIF_ACK`, `170 bytes`, OpCode `0x2808` — Sección 7).
- **`40` artefactos internos C++ STL** (`std::vector`, `std::allocator`, `std::_Vector_iterator` — Sección 5.3).

### 5.1 Las 1.293 Estructuras Oficiales `PROTO_NC_*` de `Zone.pdb` y su `sizeof` Exacto

| # | Estructura Oficial `PROTO_NC_*` | `sizeof` | # | Estructura Oficial `PROTO_NC_*` | `sizeof` |
| :---: | :--- | :---: | :---: | :--- | :---: |
| 1 | `PROTO_NC_ACT_ACTIONBYITEM_ACK` | `3 bytes` | 648 | `PROTO_NC_ITEMDB_BRACELET_UPGRADE_REQ` | `33 bytes` |
| 2 | `PROTO_NC_ACT_ACTIONBYITEM_REQ` | `1 bytes` | 649 | `PROTO_NC_ITEMDB_BUYALLSUC_ACK` | `2 bytes` |
| 3 | `PROTO_NC_ACT_ANIMATION_LEVEL_CHANGE_CMD` | `3 bytes` | 650 | `PROTO_NC_ITEMDB_BUYALL_REQ` | `62 bytes` |
| 4 | `PROTO_NC_ACT_ANIMATION_START_CMD` | `34 bytes` | 651 | `PROTO_NC_ITEMDB_BUYCAPSULE_ACK` | `12 bytes` |
| 5 | `PROTO_NC_ACT_ANIMATION_STOP_CMD` | `34 bytes` | 652 | `PROTO_NC_ITEMDB_BUYCAPSULE_REQ` | `252 bytes` |
| 6 | `PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK` | `2 bytes` | 653 | `PROTO_NC_ITEMDB_BUYLOTFAIL_ACK` | `4 bytes` |
| 7 | `PROTO_NC_ACT_CHANGEMODE_REQ` | `1 bytes` | 654 | `PROTO_NC_ITEMDB_BUYLOTSUC_ACK` | `2 bytes` |
| 8 | `PROTO_NC_ACT_CHAT_REQ` | `2 bytes` | 655 | `PROTO_NC_ITEMDB_BUYLOT_REQ` | `38 bytes` |
| 9 | `PROTO_NC_ACT_CREATECASTBAR` | `2 bytes` | 656 | `PROTO_NC_ITEMDB_CAPSULEITEM_ACK` | `129 bytes` |
| 10 | `PROTO_NC_ACT_EFFECT_MESSAGE_CMD` | `4 bytes` | 657 | `PROTO_NC_ITEMDB_CAPSULEITEM_REQ` | `16 bytes` |
| 11 | `PROTO_NC_ACT_EMOTICON_CMD` | `1 bytes` | 658 | `PROTO_NC_ITEMDB_CHARGED_LIST_ACK` | `11 bytes` |
| 12 | `PROTO_NC_ACT_EVENT_CODE_ACTION_CMD` | `6 bytes` | 659 | `PROTO_NC_ITEMDB_CHARGED_LIST_REQ` | `14 bytes` |
| 13 | `PROTO_NC_ACT_GATHERCOMPLETE_ACK` | `4 bytes` | 660 | `PROTO_NC_ITEMDB_CHARGED_WITHDRAW_ACK` | `26 bytes` |
| 14 | `PROTO_NC_ACT_GATHERSTART_ACK` | `6 bytes` | 661 | `PROTO_NC_ITEMDB_CHARGED_WITHDRAW_REQ` | `29 bytes` |
| 15 | `PROTO_NC_ACT_GATHERSTART_REQ` | `2 bytes` | 662 | `PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_ACK` | `23 bytes` |
| 16 | `PROTO_NC_ACT_MOVEFAIL_CMD` | `8 bytes` | 663 | `PROTO_NC_ITEMDB_CHAT_COLOR_CHANGE_REQ` | `25 bytes` |
| 17 | `PROTO_NC_ACT_MOVESPEED_CMD` | `4 bytes` | 664 | `PROTO_NC_ITEMDB_CHESTITEM_ACK` | `19 bytes` |
| 18 | `PROTO_NC_ACT_MOVEWALK_CMD` | `16 bytes` | 665 | `PROTO_NC_ITEMDB_CHESTITEM_REQ` | `16 bytes` |
| 19 | `PROTO_NC_ACT_NPCCLICK_CMD` | `2 bytes` | 666 | `PROTO_NC_ITEMDB_CLASS_CHANGE_ACK` | `10 bytes` |
| 20 | `PROTO_NC_ACT_NPCMENUOPEN_ACK` | `1 bytes` | 667 | `PROTO_NC_ITEMDB_CLASS_CHANGE_REQ` | `41 bytes` |
| 21 | `PROTO_NC_ACT_NPCMENUOPEN_REQ` | `2 bytes` | 668 | `PROTO_NC_ITEMDB_CREATEITEMLISTFAIL_ACK` | `10 bytes` |
| 22 | `PROTO_NC_ACT_NPC_ACTION_CMD` | `7 bytes` | 669 | `PROTO_NC_ITEMDB_CREATEITEMLISTSUC_ACK` | `8 bytes` |
| 23 | `PROTO_NC_ACT_NPC_MENU_CMD` | `1 bytes` | 670 | `PROTO_NC_ITEMDB_CREATEITEMLIST_REQ` | `23 bytes` |
| 24 | `PROTO_NC_ACT_OBJECT_EFFECT_CMD` | `6 bytes` | 671 | `PROTO_NC_ITEMDB_DEPOSIT_ACK` | `18 bytes` |
| 25 | `PROTO_NC_ACT_OBJECT_SOUND_CMD` | `34 bytes` | 672 | `PROTO_NC_ITEMDB_DEPOSIT_REQ` | `20 bytes` |
| 26 | `PROTO_NC_ACT_PLAY_SOUND_CMD` | `2 bytes` | 673 | `PROTO_NC_ITEMDB_DESTROY_ACK` | `10 bytes` |
| 27 | `PROTO_NC_ACT_PRODUCE_CAST_REQ` | `2 bytes` | 674 | `PROTO_NC_ITEMDB_DISMANTLE_ACK` | `10 bytes` |
| 28 | `PROTO_NC_ACT_PRODUCE_MAKE_CMD` | `4 bytes` | 675 | `PROTO_NC_ITEMDB_DISMANTLE_REQ` | `26 bytes` |
| 29 | `PROTO_NC_ACT_REINFORCEMOVEBYPATH_CMD` | `6 bytes` | 676 | `PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-produce>` | `43 bytes` |
| 30 | `PROTO_NC_ACT_REINFORCE_RUN_CMD` | `8 bytes` | 677 | `PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-produce>::<unnamed-type-iteminfo>` | `38 bytes` |
| 31 | `PROTO_NC_ACT_REINFORCE_STOP_CMD` | `8 bytes` | 678 | `PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-raw>` | `17 bytes` |
| 32 | `PROTO_NC_ACT_RIDE_FEEDING_REQ` | `1 bytes` | 679 | `PROTO_NC_ITEMDB_DISMANTLE_REQ::<unnamed-type-raw>::<unnamed-type-iteminfo>` | `12 bytes` |
| 33 | `PROTO_NC_ACT_RIDE_HUNGRY_CMD` | `2 bytes` | 680 | `PROTO_NC_ITEMDB_DROPALLSUC_ACK` | `111 bytes` |
| 34 | `PROTO_NC_ACT_RIDE_ON_CMD` | `2 bytes` | 681 | `PROTO_NC_ITEMDB_DROPALL_REQ` | `123 bytes` |
| 35 | `PROTO_NC_ACT_ROAR_REQ` | `2 bytes` | 682 | `PROTO_NC_ITEMDB_DROPLOT_REQ` | `30 bytes` |
| 36 | `PROTO_NC_ACT_SCRIPT_MSG_CMD` | `39 bytes` | 683 | `PROTO_NC_ITEMDB_DROPLOT_REQ::<unnamed-type-dropitem>` | `12 bytes` |
| 37 | `PROTO_NC_ACT_SCRIPT_MSG_WORLD_CMD` | `39 bytes` | 684 | `PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK` | `18 bytes` |
| 38 | `PROTO_NC_ACT_SETITEMHEALEFFECT` | `5 bytes` | 685 | `PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK::<unnamed-type-equipment>` | `6 bytes` |
| 39 | `PROTO_NC_ACT_SETITEMHEALEFFECT::<unnamed-type-flag>` | `1 bytes` | 686 | `PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_ACK::<unnamed-type-gem>` | `2 bytes` |
| 40 | `PROTO_NC_ACT_SHOUT_CMD` | `2 bytes` | 687 | `PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ` | `34 bytes` |
| 41 | `PROTO_NC_ACT_SHOW_CINEMATIC_CMD` | `2 bytes` | 688 | `PROTO_NC_ITEMDB_ENCHANT_ADD_GEM_REQ::<unnamed-type-equipment>` | `12 bytes` |
| 42 | `PROTO_NC_ACT_SOMEEONEJUMP_CMD` | `2 bytes` | 689 | `PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK` | `23 bytes` |
| 43 | `PROTO_NC_ACT_SOMEONECHANGEMODE_CMD` | `3 bytes` | 690 | `PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK::<unnamed-type-drill>` | `10 bytes` |
| 44 | `PROTO_NC_ACT_SOMEONECHAT_CMD` | `7 bytes` | 691 | `PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK::<unnamed-type-equipment>` | `3 bytes` |
| 45 | `PROTO_NC_ACT_SOMEONECHAT_CMD::<unnamed-type-flag>` | `1 bytes` | 692 | `PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ` | `32 bytes` |
| 46 | `PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD` | `5 bytes` | 693 | `PROTO_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_REQ::<unnamed-type-equipment>` | `12 bytes` |
| 47 | `PROTO_NC_ACT_SOMEONEEMOTICON_CMD` | `3 bytes` | 694 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK` | `26 bytes` |
| 48 | `PROTO_NC_ACT_SOMEONEFOLDTENT_CMD` | `45 bytes` | 695 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-equipment>` | `2 bytes` |
| 49 | `PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD` | `4 bytes` | 696 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-gemSlot>` | `1 bytes` |
| 50 | `PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD` | `4 bytes` | 697 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK::<unnamed-type-remover>` | `10 bytes` |
| 51 | `PROTO_NC_ACT_SOMEONEGATHERSTART_CMD` | `6 bytes` | 698 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ` | `36 bytes` |
| 52 | `PROTO_NC_ACT_SOMEONEMOVEWALK_CMD` | `22 bytes` | 699 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ::<unnamed-type-equipment>` | `12 bytes` |
| 53 | `PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr>` | `2 bytes` | 700 | `PROTO_NC_ITEMDB_ENCHANT_REMOVE_GEM_REQ::<unnamed-type-gemSlot>` | `1 bytes` |
| 54 | `PROTO_NC_ACT_SOMEONEPITCHTENT_CMD` | `14 bytes` | 701 | `PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ` | `24 bytes` |
| 55 | `PROTO_NC_ACT_SOMEONEPRODUCE_CASTCUT_CMD` | `2 bytes` | 702 | `PROTO_NC_ITEMDB_ENCHANT_SET_GEM_LOT_REQ::<unnamed-type-sockets>` | `3 bytes` |
| 56 | `PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD` | `4 bytes` | 703 | `PROTO_NC_ITEMDB_EQUIPFAIL_ACK` | `4 bytes` |
| 57 | `PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD` | `4 bytes` | 704 | `PROTO_NC_ITEMDB_EQUIPSUC_ACK` | `2 bytes` |
| 58 | `PROTO_NC_ACT_SOMEONERIDE_OFF_CMD` | `2 bytes` | 705 | `PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_ACK` | `34 bytes` |
| 59 | `PROTO_NC_ACT_SOMEONERIDE_ON_CMD` | `4 bytes` | 706 | `PROTO_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_REQ` | `32 bytes` |
| 60 | `PROTO_NC_ACT_SOMEONESHOUT_CMD` | `23 bytes` | 707 | `PROTO_NC_ITEMDB_EQUIP_REQ` | `29 bytes` |
| 61 | `PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag>` | `1 bytes` | 708 | `PROTO_NC_ITEMDB_EQUIP_REQ::<unnamed-type-unequip>` | `9 bytes` |
| 62 | `PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker>` | `20 bytes` | 709 | `PROTO_NC_ITEMDB_EXCHANGEFAIL_ACK` | `4 bytes` |
| 63 | `PROTO_NC_ACT_SOMEONESTOP_CMD` | `10 bytes` | 710 | `PROTO_NC_ITEMDB_EXCHANGESUC_ACK` | `2 bytes` |
| 64 | `PROTO_NC_ACT_SOMEONEWALK_CMD` | `18 bytes` | 711 | `PROTO_NC_ITEMDB_EXCHANGE_REQ` | `42 bytes` |
| 65 | `PROTO_NC_ACT_STOP_REQ` | `8 bytes` | 712 | `PROTO_NC_ITEMDB_FAIL_ACK` | `2 bytes` |
| 66 | `PROTO_NC_ACT_WALK_REQ` | `16 bytes` | 713 | `PROTO_NC_ITEMDB_FURNITURE_ENDURE_ACK` | `18 bytes` |
| 67 | `PROTO_NC_ACT_WEDDING_COUPLE_ENTRANCE_RNG` | `16 bytes` | 714 | `PROTO_NC_ITEMDB_FURNITURE_ENDURE_REQ` | `35 bytes` |
| 68 | `PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK` | `2 bytes` | 715 | `PROTO_NC_ITEMDB_GETFROMCAPSULE_ACK` | `10 bytes` |
| 69 | `PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_ACK` | `174 bytes` | 716 | `PROTO_NC_ITEMDB_GETFROMCAPSULE_REQ` | `30 bytes` |
| 70 | `PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK` | `46 bytes` | 717 | `PROTO_NC_ITEMDB_GETFROMCHESTFAIL_ACK` | `10 bytes` |
| 71 | `PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK` | `40 bytes` | 718 | `PROTO_NC_ITEMDB_GETFROMCHESTSUC_ACK` | `8 bytes` |
| 72 | `PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ` | `39 bytes` | 719 | `PROTO_NC_ITEMDB_GETFROMCHEST_REQ` | `21 bytes` |
| 73 | `PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK` | `41 bytes` | 720 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK` | `18 bytes` |
| 74 | `PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ` | `39 bytes` | 721 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ` | `20 bytes` |
| 75 | `PROTO_NC_ANNOUNCE_Z2W_CMD` | `2 bytes` | 722 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK` | `16 bytes` |
| 76 | `PROTO_NC_BAT_ABSTATEINFORM_CMD` | `8 bytes` | 723 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ` | `11 bytes` |
| 77 | `PROTO_NC_BAT_ABSTATERESET_CMD` | `6 bytes` | 724 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK` | `18 bytes` |
| 78 | `PROTO_NC_BAT_ABSTATESET_CMD` | `6 bytes` | 725 | `PROTO_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ` | `20 bytes` |
| 79 | `PROTO_NC_BAT_ABSTATE_ERASE_ACK` | `2 bytes` | 726 | `PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_ACK` | `18 bytes` |
| 80 | `PROTO_NC_BAT_ABSTATE_ERASE_REQ` | `2 bytes` | 727 | `PROTO_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_REQ` | `20 bytes` |
| 81 | `PROTO_NC_BAT_AREADOTDAMAGE_CMD` | `5 bytes` | 728 | `PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALLSUC_ACK` | `10 bytes` |
| 82 | `PROTO_NC_BAT_AREADOTDAMAGE_CMD::targetinfo` | `12 bytes` | 729 | `PROTO_NC_ITEMDB_GUILD_TOKEN_BUYALL_REQ` | `48 bytes` |
| 83 | `PROTO_NC_BAT_ASSIST_ACK` | `2 bytes` | 730 | `PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOTSUC_ACK` | `10 bytes` |
| 84 | `PROTO_NC_BAT_ASSIST_REQ` | `2 bytes` | 731 | `PROTO_NC_ITEMDB_GUILD_TOKEN_BUYLOT_REQ` | `24 bytes` |
| 85 | `PROTO_NC_BAT_BASH_HITTED_CMD` | `13 bytes` | 732 | `PROTO_NC_ITEMDB_INC_DEC_MONEY_REQ` | `9 bytes` |
| 86 | `PROTO_NC_BAT_BASH_HITTED_CMD::<unnamed-type-damage>` | `2 bytes` | 733 | `PROTO_NC_ITEMDB_ITEMBREAKFAIL_ACK` | `10 bytes` |
| 87 | `PROTO_NC_BAT_BASH_HITTED_CMD::<unnamed-type-flag>` | `1 bytes` | 734 | `PROTO_NC_ITEMDB_ITEMBREAKSUC_ACK` | `8 bytes` |
| 88 | `PROTO_NC_BAT_BASH_HIT_CMD` | `13 bytes` | 735 | `PROTO_NC_ITEMDB_ITEMBREAK_REQ` | `25 bytes` |
| 89 | `PROTO_NC_BAT_BASH_HIT_CMD::<unnamed-type-damage>` | `2 bytes` | 736 | `PROTO_NC_ITEMDB_ITEMCHANGE_ACK` | `10 bytes` |
| 90 | `PROTO_NC_BAT_BASH_HIT_CMD::<unnamed-type-flag>` | `1 bytes` | 737 | `PROTO_NC_ITEMDB_ITEMCHANGE_REQ` | `16 bytes` |
| 91 | `PROTO_NC_BAT_BASH_MISSED_CMD` | `5 bytes` | 738 | `PROTO_NC_ITEMDB_ITEMINFO_ACK` | `121 bytes` |
| 92 | `PROTO_NC_BAT_BASH_MISSED_CMD::<unnamed-type-flag>` | `1 bytes` | 739 | `PROTO_NC_ITEMDB_ITEMINFO_REQ` | `16 bytes` |
| 93 | `PROTO_NC_BAT_BASH_MISS_CMD` | `5 bytes` | 740 | `PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ` | `34 bytes` |
| 94 | `PROTO_NC_BAT_BASH_MISS_CMD::<unnamed-type-flag>` | `1 bytes` | 741 | `PROTO_NC_ITEMDB_ITEMMONEY_BUYALL_REQ::CREATE_ITEM_INFO` | `113 bytes` |
| 95 | `PROTO_NC_BAT_CEASE_FIRE_CMD` | `2 bytes` | 742 | `PROTO_NC_ITEMDB_ITEMMONEY_BUYLOT_REQ` | `47 bytes` |
| 96 | `PROTO_NC_BAT_DOTDAMAGE_CMD` | `13 bytes` | 743 | `PROTO_NC_ITEMDB_ITEMMONEY_BUY_ACK` | `10 bytes` |
| 97 | `PROTO_NC_BAT_EXPGAIN_CMD` | `6 bytes` | 744 | `PROTO_NC_ITEMDB_ITEMREBUILD_ACK` | `10 bytes` |
| 98 | `PROTO_NC_BAT_EXPLOST_CMD` | `4 bytes` | 745 | `PROTO_NC_ITEMDB_ITEMREBUILD_REQ` | `54 bytes` |
| 99 | `PROTO_NC_BAT_FAMEGAIN_CMD` | `4 bytes` | 746 | `PROTO_NC_ITEMDB_ITEMTOTALINFORM_ACK` | `111 bytes` |
| 100 | `PROTO_NC_BAT_HPCHANGE_CMD` | `6 bytes` | 747 | `PROTO_NC_ITEMDB_ITEMTOTALINFORM_REQ` | `8 bytes` |
| 101 | `PROTO_NC_BAT_LEVELUP_CMD` | `235 bytes` | 748 | `PROTO_NC_ITEMDB_ITEMTRADEFAIL_ACK` | `10 bytes` |
| 102 | `PROTO_NC_BAT_LPCHANGE_CMD` | `4 bytes` | 749 | `PROTO_NC_ITEMDB_ITEMTRADESUC_ACK` | `8 bytes` |
| 103 | `PROTO_NC_BAT_MOBSLAYER_CMD` | `3 bytes` | 750 | `PROTO_NC_ITEMDB_ITEMTRADE_REQ` | `42 bytes` |
| 104 | `PROTO_NC_BAT_PKINPKFIELD_CLIENT_CMD` | `8 bytes` | 751 | `PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-cen>` | `25 bytes` |
| 105 | `PROTO_NC_BAT_PKINPKFIELD_WMS_CMD` | `10 bytes` | 752 | `PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-tradeitem>` | `26 bytes` |
| 106 | `PROTO_NC_BAT_REALLYKILL_CMD` | `4 bytes` | 753 | `PROTO_NC_ITEMDB_ITEMTRADE_REQ::<unnamed-type-tradeitem>::<unnamed-type-invenA>` | `13 bytes` |
| 107 | `PROTO_NC_BAT_REFLECTIONDAMAGE_CMD` | `14 bytes` | 754 | `PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_ACK` | `28 bytes` |
| 108 | `PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ` | `10 bytes` | 755 | `PROTO_NC_ITEMDB_MAPLINK_ITEM_CONSUME_REQ` | `40 bytes` |
| 109 | `PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD` | `6 bytes` | 756 | `PROTO_NC_ITEMDB_MERGEFAIL_ACK` | `4 bytes` |
| 110 | `PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD` | `5 bytes` | 757 | `PROTO_NC_ITEMDB_MERGESUC_ACK` | `12 bytes` |
| 111 | `PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage` | `14 bytes` | 758 | `PROTO_NC_ITEMDB_MERGE_REQ` | `36 bytes` |
| 112 | `PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage::<unnamed-type-flag>` | `2 bytes` | 759 | `PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_ACK` | `12 bytes` |
| 113 | `PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD` | `12 bytes` | 760 | `PROTO_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_REQ` | `32 bytes` |
| 114 | `PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD` | `6 bytes` | 761 | `PROTO_NC_ITEMDB_MINIHOUSE_PORTAL_EFFECT_DEMANDGOOD_REQ` | `32 bytes` |
| 115 | `PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ` | `4 bytes` | 762 | `PROTO_NC_ITEMDB_MIX_ITEM_ACK` | `10 bytes` |
| 116 | `PROTO_NC_BAT_SKILLBLAST_LIGHTNINGWAVE_CMD` | `8 bytes` | 763 | `PROTO_NC_ITEMDB_MIX_ITEM_REQ` | `163 bytes` |
| 117 | `PROTO_NC_BAT_SKILLENCHANT_REQ` | `4 bytes` | 764 | `PROTO_NC_ITEMDB_MOVER_RAREMOVER_ACK` | `114 bytes` |
| 118 | `PROTO_NC_BAT_SKILLSMASH_CMD` | `5 bytes` | 765 | `PROTO_NC_ITEMDB_MOVER_RAREMOVER_REQ` | `343 bytes` |
| 119 | `PROTO_NC_BAT_SKILLSMASH_HITTED_CMD` | `9 bytes` | 766 | `PROTO_NC_ITEMDB_MOVER_UPGRADE_ACK` | `115 bytes` |
| 120 | `PROTO_NC_BAT_SKILLSMASH_HIT_CMD` | `7 bytes` | 767 | `PROTO_NC_ITEMDB_MOVER_UPGRADE_REQ` | `233 bytes` |
| 121 | `PROTO_NC_BAT_SKILLSMASH_HIT_CMD::<unnamed-type-damage>` | `2 bytes` | 768 | `PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_ACK` | `10 bytes` |
| 122 | `PROTO_NC_BAT_SKILLSMASH_MISS_CMD` | `3 bytes` | 769 | `PROTO_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_REQ` | `19 bytes` |
| 123 | `PROTO_NC_BAT_SMASH_HITTED_ACK` | `5 bytes` | 770 | `PROTO_NC_ITEMDB_NEW_UPGRADE_ACK` | `10 bytes` |
| 124 | `PROTO_NC_BAT_SMASH_HITTED_ACK::<unnamed-type-hitted>` | `3 bytes` | 771 | `PROTO_NC_ITEMDB_NEW_UPGRADE_REQ` | `36 bytes` |
| 125 | `PROTO_NC_BAT_SMASH_HITTED_ACK::<unnamed-type-hitted>::<unnamed-type-damage>` | `2 bytes` | 772 | `PROTO_NC_ITEMDB_OPENSTORAGE_ACK` | `17 bytes` |
| 126 | `PROTO_NC_BAT_SMASH_HIT_ACK` | `5 bytes` | 773 | `PROTO_NC_ITEMDB_OPENSTORAGE_FAIL_ACK` | `8 bytes` |
| 127 | `PROTO_NC_BAT_SMASH_HIT_ACK::<unnamed-type-damage>` | `2 bytes` | 774 | `PROTO_NC_ITEMDB_OPENSTORAGE_REQ` | `12 bytes` |
| 128 | `PROTO_NC_BAT_SMASH_MISS_ACK` | `3 bytes` | 775 | `PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_ACK` | `23 bytes` |
| 129 | `PROTO_NC_BAT_SOMEONEBASH_HIT_CMD` | `13 bytes` | 776 | `PROTO_NC_ITEMDB_OPEN_GUILD_STORAGE_REQ` | `10 bytes` |
| 130 | `PROTO_NC_BAT_SOMEONEBASH_MISS_CMD` | `7 bytes` | 777 | `PROTO_NC_ITEMDB_PICKMERGE_REQ` | `26 bytes` |
| 131 | `PROTO_NC_BAT_SOMEONEDAMAGED_LARGE_CMD` | `7 bytes` | 778 | `PROTO_NC_ITEMDB_PICKMONEY_REQ` | `16 bytes` |
| 132 | `PROTO_NC_BAT_SOMEONEDAMAGED_SMALL_CMD` | `6 bytes` | 779 | `PROTO_NC_ITEMDB_PRODUCE_ACK` | `12 bytes` |
| 133 | `PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD` | `2 bytes` | 780 | `PROTO_NC_ITEMDB_PRODUCE_REQ` | `50 bytes` |
| 134 | `PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD` | `12 bytes` | 781 | `PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-produce>` | `39 bytes` |
| 135 | `PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD` | `14 bytes` | 782 | `PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-produce>::<unnamed-type-iteminfo>` | `38 bytes` |
| 136 | `PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD` | `8 bytes` | 783 | `PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-raw>` | `17 bytes` |
| 137 | `PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD` | `6 bytes` | 784 | `PROTO_NC_ITEMDB_PRODUCE_REQ::<unnamed-type-raw>::<unnamed-type-iteminfo>` | `12 bytes` |
| 138 | `PROTO_NC_BAT_SOMEONESKILLSMASH_DEAD_CMD` | `6 bytes` | 785 | `PROTO_NC_ITEMDB_QUESTREWARD_ACK` | `12 bytes` |
| 139 | `PROTO_NC_BAT_SOMEONESMASH_DAMAGED_CMD` | `5 bytes` | 786 | `PROTO_NC_ITEMDB_QUESTREWARD_REQ` | `45 bytes` |
| 140 | `PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD` | `13 bytes` | 787 | `PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_ACK` | `10 bytes` |
| 141 | `PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD::<unnamed-type-flag>` | `1 bytes` | 788 | `PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_ACK` | `21 bytes` |
| 142 | `PROTO_NC_BAT_SOULCOLLECT_CMD` | `3 bytes` | 789 | `PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_REQ` | `134 bytes` |
| 143 | `PROTO_NC_BAT_SPCHANGE_CMD` | `4 bytes` | 790 | `PROTO_NC_ITEMDB_RANDOMOPTION_CHANGE_REQ` | `220 bytes` |
| 144 | `PROTO_NC_BAT_SUMEONELEVELUP_CMD` | `4 bytes` | 791 | `PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK` | `10 bytes` |
| 145 | `PROTO_NC_BAT_SUMEONESKILLCUT_CMD` | `2 bytes` | 792 | `PROTO_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ` | `146 bytes` |
| 146 | `PROTO_NC_BAT_SWING_DAMAGE_CMD` | `16 bytes` | 793 | `PROTO_NC_ITEMDB_REINFORCEUNEQUIPSUC_ACK` | `2 bytes` |
| 147 | `PROTO_NC_BAT_SWING_DAMAGE_CMD::<unnamed-type-flag>` | `2 bytes` | 794 | `PROTO_NC_ITEMDB_REINFORCEUNEQUIP_REQ` | `14 bytes` |
| 148 | `PROTO_NC_BAT_SWING_START_CMD` | `9 bytes` | 795 | `PROTO_NC_ITEMDB_RELOCFAIL_ACK` | `4 bytes` |
| 149 | `PROTO_NC_BAT_TARGETCHANGE_CMD` | `2 bytes` | 796 | `PROTO_NC_ITEMDB_RELOCSUC_ACK` | `2 bytes` |
| 150 | `PROTO_NC_BAT_TARGETINFO_CMD` | `30 bytes` | 797 | `PROTO_NC_ITEMDB_RELOC_REQ` | `28 bytes` |
| 151 | `PROTO_NC_BAT_TARGET_REQ` | `2 bytes` | 798 | `PROTO_NC_ITEMDB_REPURCHASE_ALL_ACK` | `10 bytes` |
| 152 | `PROTO_NC_BAT_TOGGLESKILL_OFF_CMD` | `2 bytes` | 799 | `PROTO_NC_ITEMDB_REPURCHASE_ALL_REQ` | `64 bytes` |
| 153 | `PROTO_NC_BAT_TOGGLESKILL_ON_CMD` | `2 bytes` | 800 | `PROTO_NC_ITEMDB_REPURCHASE_LOT_ACK` | `10 bytes` |
| 154 | `PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD` | `4 bytes` | 801 | `PROTO_NC_ITEMDB_REPURCHASE_LOT_REQ` | `40 bytes` |
| 155 | `PROTO_NC_BOOTH_BUYREFRESH_CMD` | `3 bytes` | 802 | `PROTO_NC_ITEMDB_RESET_SCROLL_LINK_MAP_INFO_CMD` | `4 bytes` |
| 156 | `PROTO_NC_BOOTH_ENTRY_BUY_ACK` | `5 bytes` | 803 | `PROTO_NC_ITEMDB_SELLALLSUC_ACK` | `2 bytes` |
| 157 | `PROTO_NC_BOOTH_ENTRY_BUY_ACK::BoothItemList` | `13 bytes` | 804 | `PROTO_NC_ITEMDB_SELLALL_REQ` | `33 bytes` |
| 158 | `PROTO_NC_BOOTH_ENTRY_REQ` | `2 bytes` | 805 | `PROTO_NC_ITEMDB_SELLLOTFAIL_ACK` | `4 bytes` |
| 159 | `PROTO_NC_BOOTH_ENTRY_SELL_ACK` | `5 bytes` | 806 | `PROTO_NC_ITEMDB_SELLLOTSUC_ACK` | `2 bytes` |
| 160 | `PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList` | `113 bytes` | 807 | `PROTO_NC_ITEMDB_SELLLOT_REQ` | `26 bytes` |
| 161 | `PROTO_NC_BOOTH_ITEMTRADE_REQ` | `5 bytes` | 808 | `PROTO_NC_ITEMDB_SHIELDENDURESET_CMD` | `12 bytes` |
| 162 | `PROTO_NC_BOOTH_OPEN_REQ` | `31 bytes` | 809 | `PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_ACK` | `24 bytes` |
| 163 | `PROTO_NC_BOOTH_OPEN_REQ::<unnamed-type-flag>` | `1 bytes` | 810 | `PROTO_NC_ITEMDB_SHIELDENDURE_CHARGE_REQ` | `35 bytes` |
| 164 | `PROTO_NC_BOOTH_OPEN_REQ::BoothItem` | `12 bytes` | 811 | `PROTO_NC_ITEMDB_SOULSTONEBUYFAIL_ACK` | `4 bytes` |
| 165 | `PROTO_NC_BOOTH_REFRESH_REQ` | `2 bytes` | 812 | `PROTO_NC_ITEMDB_SOULSTONEBUYSUC_ACK` | `2 bytes` |
| 166 | `PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD` | `2 bytes` | 813 | `PROTO_NC_ITEMDB_SOULSTONEBUY_REQ` | `17 bytes` |
| 167 | `PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK` | `12 bytes` | 814 | `PROTO_NC_ITEMDB_SPLITFAIL_ACK` | `4 bytes` |
| 168 | `PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ` | `6 bytes` | 815 | `PROTO_NC_ITEMDB_SPLITSUC_ACK` | `12 bytes` |
| 169 | `PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK` | `3 bytes` | 816 | `PROTO_NC_ITEMDB_SPLIT_N_MERGEFAIL_ACK` | `4 bytes` |
| 170 | `PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ` | `6 bytes` | 817 | `PROTO_NC_ITEMDB_SPLIT_N_MERGESUC_ACK` | `24 bytes` |
| 171 | `PROTO_NC_BOOTH_SOMEONECLOSE_CMD` | `45 bytes` | 818 | `PROTO_NC_ITEMDB_SPLIT_N_MERGE_REQ` | `38 bytes` |
| 172 | `PROTO_NC_BOOTH_SOMEONEINTERIORSTART_CMD` | `32 bytes` | 819 | `PROTO_NC_ITEMDB_SPLIT_REQ` | `62 bytes` |
| 173 | `PROTO_NC_BOOTH_SOMEONEOPEN_CMD` | `45 bytes` | 820 | `PROTO_NC_ITEMDB_TERMEXTEND_ACK` | `10 bytes` |
| 174 | `PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD` | `14 bytes` | 821 | `PROTO_NC_ITEMDB_TERMEXTEND_REQ` | `141 bytes` |
| 175 | `PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD` | `3 bytes` | 822 | `PROTO_NC_ITEMDB_UES_FRIEND_POINT_ACK` | `15 bytes` |
| 176 | `PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD` | `2 bytes` | 823 | `PROTO_NC_ITEMDB_UES_FRIEND_POINT_REQ` | `123 bytes` |
| 177 | `PROTO_NC_BRIEFINFO_BUILDDOOR_CMD` | `48 bytes` | 824 | `PROTO_NC_ITEMDB_UNEQUIPFAIL_ACK` | `4 bytes` |
| 178 | `PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD` | `5 bytes` | 825 | `PROTO_NC_ITEMDB_UNEQUIPSUC_ACK` | `2 bytes` |
| 179 | `PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD` | `6 bytes` | 826 | `PROTO_NC_ITEMDB_UNEQUIP_REQ` | `14 bytes` |
| 180 | `PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD` | `9 bytes` | 827 | `PROTO_NC_ITEMDB_UPGRADE_ACK` | `10 bytes` |
| 181 | `PROTO_NC_BRIEFINFO_CHARACTER_CMD` | `1 bytes` | 828 | `PROTO_NC_ITEMDB_UPGRADE_REQ` | `78 bytes` |
| 182 | `PROTO_NC_BRIEFINFO_DOOR_CMD` | `1 bytes` | 829 | `PROTO_NC_ITEMDB_USEALL_REQ` | `12 bytes` |
| 183 | `PROTO_NC_BRIEFINFO_DROPEDITEM_CMD` | `15 bytes` | 830 | `PROTO_NC_ITEMDB_USELOT_REQ` | `14 bytes` |
| 184 | `PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr>` | `1 bytes` | 831 | `PROTO_NC_ITEMDB_USE_ACK` | `4 bytes` |
| 185 | `PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD` | `48 bytes` | 832 | `PROTO_NC_ITEMDB_USE_VARIATION_ITEM_REQ` | `21 bytes` |
| 186 | `PROTO_NC_BRIEFINFO_EFFECTBLAST_CMD::<unnamed-type-flag>` | `1 bytes` | 833 | `PROTO_NC_ITEMDB_WEAPONENDURESET_CMD` | `12 bytes` |
| 187 | `PROTO_NC_BRIEFINFO_EFFECT_CMD` | `1 bytes` | 834 | `PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_ACK` | `24 bytes` |
| 188 | `PROTO_NC_BRIEFINFO_INFORM_CMD` | `6 bytes` | 835 | `PROTO_NC_ITEMDB_WEAPONENDURE_CHARGE_REQ` | `35 bytes` |
| 189 | `PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD` | `1 bytes` | 836 | `PROTO_NC_ITEMDB_WITHDRAW_ACK` | `18 bytes` |
| 190 | `PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD` | `235 bytes` | 837 | `PROTO_NC_ITEMDB_WITHDRAW_REQ` | `20 bytes` |
| 191 | `PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD::<unnamed-type-shapedata>` | `45 bytes` | 838 | `PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK` | `2 bytes` |
| 192 | `PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD` | `1 bytes` | 839 | `PROTO_NC_ITEM_BRACELET_UPGRADE_ACK` | `2 bytes` |
| 193 | `PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD` | `16 bytes` | 840 | `PROTO_NC_ITEM_BRACELET_UPGRADE_REQ` | `2 bytes` |
| 194 | `PROTO_NC_BRIEFINFO_MINIHOUSEBUILD_CMD` | `64 bytes` | 841 | `PROTO_NC_ITEM_BUY_REQ` | `6 bytes` |
| 195 | `PROTO_NC_BRIEFINFO_MINIHOUSE_CMD` | `1 bytes` | 842 | `PROTO_NC_ITEM_CELLCHANGE_CMD` | `6 bytes` |
| 196 | `PROTO_NC_BRIEFINFO_MOB_CMD` | `1 bytes` | 843 | `PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK` | `5 bytes` |
| 197 | `PROTO_NC_BRIEFINFO_MOVER_CMD` | `1 bytes` | 844 | `PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ` | `2 bytes` |
| 198 | `PROTO_NC_BRIEFINFO_PET_CMD` | `1 bytes` | 845 | `PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK` | `6 bytes` |
| 199 | `PROTO_NC_BRIEFINFO_PLAYER_INFO_APPEAR_CMD` | `30 bytes` | 846 | `PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ` | `4 bytes` |
| 200 | `PROTO_NC_BRIEFINFO_PLAYER_LIST_INFO_APPEAR_CMD` | `1 bytes` | 847 | `PROTO_NC_ITEM_CHAT_COLOR_CHANGE_ACK` | `3 bytes` |
| 201 | `PROTO_NC_BRIEFINFO_REGENMOB_CMD` | `149 bytes` | 848 | `PROTO_NC_ITEM_CHAT_COLOR_CHANGE_REQ` | `1 bytes` |
| 202 | `PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag>` | `99 bytes` | 849 | `PROTO_NC_ITEM_CLASS_CHANGE_REQ` | `1 bytes` |
| 203 | `PROTO_NC_BRIEFINFO_REGENMOVER_CMD` | `139 bytes` | 850 | `PROTO_NC_ITEM_DEPOSIT_ACK` | `10 bytes` |
| 204 | `PROTO_NC_BRIEFINFO_REGENPET_CMD` | `63 bytes` | 851 | `PROTO_NC_ITEM_DEPOSIT_REQ` | `8 bytes` |
| 205 | `PROTO_NC_BRIEFINFO_UNEQUIP_CMD` | `3 bytes` | 852 | `PROTO_NC_ITEM_DICE_GAME_CMD` | `6 bytes` |
| 206 | `PROTO_NC_CHARGED_BOOTHSLOTSIZE_CMD` | `1 bytes` | 853 | `PROTO_NC_ITEM_DICE_GAME_RESULT_CMD` | `28 bytes` |
| 207 | `PROTO_NC_CHARGED_BUFFTERMINATE_CMD` | `4 bytes` | 854 | `PROTO_NC_ITEM_DICE_GAME_START_ACK` | `2 bytes` |
| 208 | `PROTO_NC_CHARGED_DELETEWEAPONTITLE_CMD` | `2 bytes` | 855 | `PROTO_NC_ITEM_DICE_GAME_START_CMD` | `6 bytes` |
| 209 | `PROTO_NC_CHARGED_RESETBUFF_CMD` | `14 bytes` | 856 | `PROTO_NC_ITEM_DICE_GAME_START_REQ` | `3 bytes` |
| 210 | `PROTO_NC_CHARGED_SETBUFF_CMD` | `32 bytes` | 857 | `PROTO_NC_ITEM_DISMANTLE_REQ` | `1 bytes` |
| 211 | `PROTO_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_REQ` | `21 bytes` | 858 | `PROTO_NC_ITEM_DROP_REQ` | `14 bytes` |
| 212 | `PROTO_NC_CHARGED_STAT_INITIALIZE_DB_FAIL_ACK` | `10 bytes` | 859 | `PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK` | `7 bytes` |
| 213 | `PROTO_NC_CHARGED_STAT_INITIALIZE_DB_SUC_ACK` | `8 bytes` | 860 | `PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ` | `4 bytes` |
| 214 | `PROTO_NC_CHARGED_STAT_INITIALIZE_SUC_CMD` | `1 bytes` | 861 | `PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK` | `5 bytes` |
| 215 | `PROTO_NC_CHARSAVE_2WLDMAN_ABSTATE_CMD` | `10 bytes` | 862 | `PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ` | `4 bytes` |
| 216 | `PROTO_NC_CHARSAVE_2WLDMAN_ALL_CMD` | `107 bytes` | 863 | `PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK` | `5 bytes` |
| 217 | `PROTO_NC_CHARSAVE_2WLDMAN_CHAT_COLOR_CMD` | `6 bytes` | 864 | `PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ` | `5 bytes` |
| 218 | `PROTO_NC_CHARSAVE_2WLDMAN_CHESTINFO_CMD` | `10 bytes` | 865 | `PROTO_NC_ITEM_ENCHANT_SET_GEM_LOT_CMD` | `4 bytes` |
| 219 | `PROTO_NC_CHARSAVE_2WLDMAN_COININFO_CMD` | `22 bytes` | 866 | `PROTO_NC_ITEM_EQUIPCHANGE_CMD` | `5 bytes` |
| 220 | `PROTO_NC_CHARSAVE_2WLDMAN_GAME_CMD` | `11 bytes` | 867 | `PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_ACK` | `2 bytes` |
| 221 | `PROTO_NC_CHARSAVE_2WLDMAN_ITEMACTIONCOOLTIME_CMD` | `6 bytes` | 868 | `PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_REQ` | `2 bytes` |
| 222 | `PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD` | `9 bytes` | 869 | `PROTO_NC_ITEM_EQUIP_REQ` | `1 bytes` |
| 223 | `PROTO_NC_CHARSAVE_2WLDMAN_ITEM_CMD::<unnamed-type-PacketOrder>` | `1 bytes` | 870 | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK` | `10 bytes` |
| 224 | `PROTO_NC_CHARSAVE_2WLDMAN_LINK_FROM_CMD` | `26 bytes` | 871 | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ` | `8 bytes` |
| 225 | `PROTO_NC_CHARSAVE_2WLDMAN_MISC_CMD` | `18 bytes` | 872 | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK` | `10 bytes` |
| 226 | `PROTO_NC_CHARSAVE_2WLDMAN_PET_LINK_RESUMMON_CMD` | `23 bytes` | 873 | `PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ` | `8 bytes` |
| 227 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD` | `8 bytes` | 874 | `PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK` | `10 bytes` |
| 228 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DOING_CMD::tagQUESTDOING` | `6 bytes` | 875 | `PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ` | `8 bytes` |
| 229 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD` | `14 bytes` | 876 | `PROTO_NC_ITEM_ITEMBREAK_CMD` | `2 bytes` |
| 230 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_DONE_CMD::tagQUESTDONE` | `12 bytes` | 877 | `PROTO_NC_ITEM_MAPLINK_SCROLL_ACK` | `4 bytes` |
| 231 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD` | `8 bytes` | 878 | `PROTO_NC_ITEM_MAPLINK_SCROLL_REQ` | `4 bytes` |
| 232 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_READ_CMD::tagQUESTREAD` | `6 bytes` | 879 | `PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_ACK` | `6 bytes` |
| 233 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD` | `8 bytes` | 880 | `PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_REQ` | `3 bytes` |
| 234 | `PROTO_NC_CHARSAVE_2WLDMAN_QUEST_REPEAT_CMD::tagQUESTREPEAT` | `6 bytes` | 881 | `PROTO_NC_ITEM_MINIMON_EQUIP_REQ` | `2 bytes` |
| 235 | `PROTO_NC_CHARSAVE_2WLDMAN_SINGLE_OPTION_CMD` | `17 bytes` | 882 | `PROTO_NC_ITEM_MIX_ITEM_ACK` | `2 bytes` |
| 236 | `PROTO_NC_CHARSAVE_2WLDMAN_SKILL_CMD` | `13 bytes` | 883 | `PROTO_NC_ITEM_MIX_ITEM_REQ` | `3 bytes` |
| 237 | `PROTO_NC_CHARSAVE_2WLDMAN_TITLE_CMD` | `1037 bytes` | 884 | `PROTO_NC_ITEM_MOVER_UPGRADE_ACK` | `108 bytes` |
| 238 | `PROTO_NC_CHARSAVE_ABSTATE_REQ` | `8 bytes` | 885 | `PROTO_NC_ITEM_MOVER_UPGRADE_REQ` | `2 bytes` |
| 239 | `PROTO_NC_CHARSAVE_ALL_ACK` | `6 bytes` | 886 | `PROTO_NC_ITEM_NEW_UPGRADE_ACK` | `2 bytes` |
| 240 | `PROTO_NC_CHARSAVE_ALL_REQ` | `105 bytes` | 887 | `PROTO_NC_ITEM_NEW_UPGRADE_REQ` | `2 bytes` |
| 241 | `PROTO_NC_CHARSAVE_CHARSTAT_CMD` | `34 bytes` | 888 | `PROTO_NC_ITEM_OPENCLASSCHANGEMENU_CMD` | `11 bytes` |
| 242 | `PROTO_NC_CHARSAVE_CHAT_COLOR_CMD` | `6 bytes` | 889 | `PROTO_NC_ITEM_OPENSTORAGEPAGE_REQ` | `1 bytes` |
| 243 | `PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_ACK` | `8 bytes` | 890 | `PROTO_NC_ITEM_PICKOTHER_ACK` | `4 bytes` |
| 244 | `PROTO_NC_CHARSAVE_DB_UI_STATE_SAVE_REQ` | `7 bytes` | 891 | `PROTO_NC_ITEM_PICK_ACK` | `10 bytes` |
| 245 | `PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK` | `2 bytes` | 892 | `PROTO_NC_ITEM_PICK_REQ` | `2 bytes` |
| 246 | `PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_REQ` | `2 bytes` | 893 | `PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_ACK` | `2 bytes` |
| 247 | `PROTO_NC_CHARSAVE_LEVEL_CMD` | `17 bytes` | 894 | `PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK` | `27 bytes` |
| 248 | `PROTO_NC_CHARSAVE_LOCATION_CMD` | `48 bytes` | 895 | `PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ` | `2 bytes` |
| 249 | `PROTO_NC_CHARSAVE_PKCOUNT_CMD` | `8 bytes` | 896 | `PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK` | `2 bytes` |
| 250 | `PROTO_NC_CHARSAVE_QUEST_DOING_REQ` | `8 bytes` | 897 | `PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ` | `2 bytes` |
| 251 | `PROTO_NC_CHARSAVE_REST_EXP_LAST_EXEC_TIME_SAVE_REQ` | `6 bytes` | 898 | `PROTO_NC_ITEM_RELOC_REQ` | `4 bytes` |
| 252 | `PROTO_NC_CHARSAVE_SELL_ITEM_INFO_CMD` | `7 bytes` | 899 | `PROTO_NC_ITEM_REPURCHASE_REQ` | `2 bytes` |
| 253 | `PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_DB_CMD` | `12 bytes` | 900 | `PROTO_NC_ITEM_REWARDINVENOPEN_ACK` | `1 bytes` |
| 254 | `PROTO_NC_CHARSAVE_SET_CHAT_BLOCK_SPAMER_WM_CMD` | `14 bytes` | 901 | `PROTO_NC_ITEM_REWARDINVENOPEN_REQ` | `2 bytes` |
| 255 | `PROTO_NC_CHARSAVE_SKILL_REQ` | `8 bytes` | 902 | `PROTO_NC_ITEM_RINGEQUIP_REQ` | `2 bytes` |
| 256 | `PROTO_NC_CHARSAVE_TITLE_REQ` | `8 bytes` | 903 | `PROTO_NC_ITEM_SELL_ITEM_DELETE_CMD` | `2 bytes` |
| 257 | `PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK` | `2 bytes` | 904 | `PROTO_NC_ITEM_SELL_ITEM_INSERT_CMD` | `105 bytes` |
| 258 | `PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ` | `1 bytes` | 905 | `PROTO_NC_ITEM_SELL_ITEM_LIST_CMD` | `1 bytes` |
| 259 | `PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_DB_CMD` | `78 bytes` | 906 | `PROTO_NC_ITEM_SELL_REQ` | `5 bytes` |
| 260 | `PROTO_NC_CHARSAVE_USEITEM_MINIMON_INFO_WORLD_CMD` | `56 bytes` | 907 | `PROTO_NC_ITEM_SHIELDENDURESET_CMD` | `4 bytes` |
| 261 | `PROTO_NC_CHAR_ABSTATE_CMD` | `6 bytes` | 908 | `PROTO_NC_ITEM_SHIELDENDURE_CHARGE_ACK` | `8 bytes` |
| 262 | `PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD` | `1 bytes` | 909 | `PROTO_NC_ITEM_SHIELDENDURE_CHARGE_REQ` | `3 bytes` |
| 263 | `PROTO_NC_CHAR_BASEPARAMCHANGE_CMD` | `1 bytes` | 910 | `PROTO_NC_ITEM_SOMEONEPICK_CMD` | `26 bytes` |
| 264 | `PROTO_NC_CHAR_BASE_CMD` | `105 bytes` | 911 | `PROTO_NC_ITEM_SOMEONEUSE_CMD` | `12 bytes` |
| 265 | `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>` | `4 bytes` | 912 | `PROTO_NC_ITEM_SPLIT_REQ` | `8 bytes` |
| 266 | `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str>` | `4 bytes` | 913 | `PROTO_NC_ITEM_TERMEXTEND_ACK` | `2 bytes` |
| 267 | `PROTO_NC_CHAR_BASE_CMD::LoginLocation` | `21 bytes` | 914 | `PROTO_NC_ITEM_TERMEXTEND_REQ` | `2 bytes` |
| 268 | `PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD` | `6 bytes` | 915 | `PROTO_NC_ITEM_UNEQUIP_REQ` | `2 bytes` |
| 269 | `PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo` | `4 bytes` | 916 | `PROTO_NC_ITEM_UPGRADE_ACK` | `2 bytes` |
| 270 | `PROTO_NC_CHAR_CARDCOLLECT_CMD` | `6 bytes` | 917 | `PROTO_NC_ITEM_UPGRADE_REQ` | `9 bytes` |
| 271 | `PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform` | `11 bytes` | 918 | `PROTO_NC_ITEM_USE_ACK` | `5 bytes` |
| 272 | `PROTO_NC_CHAR_CARDCOLLECT_REWARD_CMD` | `6 bytes` | 919 | `PROTO_NC_ITEM_USE_ACTIVESKILL_REQ` | `11 bytes` |
| 273 | `PROTO_NC_CHAR_CENCHANGE_CMD` | `8 bytes` | 920 | `PROTO_NC_ITEM_USE_REQ` | `2 bytes` |
| 274 | `PROTO_NC_CHAR_CHANGEBYCONDITION_PARAM_CMD` | `6 bytes` | 921 | `PROTO_NC_ITEM_WEAPONENDURESET_CMD` | `4 bytes` |
| 275 | `PROTO_NC_CHAR_CHARGEDBUFF_CMD` | `2 bytes` | 922 | `PROTO_NC_ITEM_WEAPONENDURE_CHARGE_ACK` | `8 bytes` |
| 276 | `PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ` | `4 bytes` | 923 | `PROTO_NC_ITEM_WEAPONENDURE_CHARGE_REQ` | `3 bytes` |
| 277 | `PROTO_NC_CHAR_CHESTINFO_CMD` | `2 bytes` | 924 | `PROTO_NC_ITEM_WITHDRAW_ACK` | `10 bytes` |
| 278 | `PROTO_NC_CHAR_CHESTINFO_REQ` | `8 bytes` | 925 | `PROTO_NC_ITEM_WITHDRAW_REQ` | `8 bytes` |
| 279 | `PROTO_NC_CHAR_CLASSCHANGE_ACK` | `7 bytes` | 926 | `PROTO_NC_KQ_ENTRYRESPONCE_ACK` | `5 bytes` |
| 280 | `PROTO_NC_CHAR_CLASSCHANGE_CMD` | `9 bytes` | 927 | `PROTO_NC_KQ_JOINER` | `5 bytes` |
| 281 | `PROTO_NC_CHAR_CLASSCHANGE_REQ` | `8 bytes` | 928 | `PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD` | `52 bytes` |
| 282 | `PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK` | `3 bytes` | 929 | `PROTO_NC_KQ_MOBKILLNUMBER_CMD` | `4 bytes` |
| 283 | `PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD` | `3 bytes` | 930 | `PROTO_NC_KQ_PLAYER_DISJOIN_CMD` | `8 bytes` |
| 284 | `PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ` | `1 bytes` | 931 | `PROTO_NC_KQ_RESTDEADNUM_CMD` | `1 bytes` |
| 285 | `PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD` | `6 bytes` | 932 | `PROTO_NC_KQ_REWARDFAIL_ACK` | `10 bytes` |
| 286 | `PROTO_NC_CHAR_CLIENT_ITEM_CMD` | `3 bytes` | 933 | `PROTO_NC_KQ_REWARDSUC_ACK` | `8 bytes` |
| 287 | `PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag>` | `1 bytes` | 934 | `PROTO_NC_KQ_REWARD_REQ` | `35 bytes` |
| 288 | `PROTO_NC_CHAR_CLIENT_SKILL_CMD` | `10 bytes` | 935 | `PROTO_NC_KQ_SCORE_BOARD_INFO_CMD` | `6 bytes` |
| 289 | `PROTO_NC_CHAR_COININFO_CMD` | `16 bytes` | 936 | `PROTO_NC_KQ_SCORE_CMD` | `5 bytes` |
| 290 | `PROTO_NC_CHAR_DEADMENU_CMD` | `9 bytes` | 937 | `PROTO_NC_KQ_SCORE_CMD::KQScore` | `38 bytes` |
| 291 | `PROTO_NC_CHAR_EMPTY_INSTANCE_DUNGEON_RNG` | `11 bytes` | 938 | `PROTO_NC_KQ_SCORE_CMD::KQScore::<unnamed-type-flag>` | `2 bytes` |
| 292 | `PROTO_NC_CHAR_EXP_CHANGED_CMD` | `14 bytes` | 939 | `PROTO_NC_KQ_SCORE_INFO_CMD` | `8 bytes` |
| 293 | `PROTO_NC_CHAR_FAMECHANGE_CMD` | `4 bytes` | 940 | `PROTO_NC_KQ_SCORE_SIMPLE_CMD` | `18 bytes` |
| 294 | `PROTO_NC_CHAR_FAMESAVE_CMD` | `8 bytes` | 941 | `PROTO_NC_KQ_SCORE_SIMPLE_CMD::<unnamed-type-self>` | `13 bytes` |
| 295 | `PROTO_NC_CHAR_FREESTAT_SET_DB_ACK` | `15 bytes` | 942 | `PROTO_NC_KQ_SCORE_SIMPLE_CMD::KQScore` | `32 bytes` |
| 296 | `PROTO_NC_CHAR_FREESTAT_SET_DB_REQ` | `15 bytes` | 943 | `PROTO_NC_KQ_TEAM_TYPE_CMD` | `1 bytes` |
| 297 | `PROTO_NC_CHAR_FREESTAT_SET_DB_REQ::<unnamed-type-DecItemInfo>` | `12 bytes` | 944 | `PROTO_NC_KQ_W2Z_DESTROY_CMD` | `4 bytes` |
| 298 | `PROTO_NC_CHAR_GET_CHAT_BLOCK_SPAMER_CMD` | `8 bytes` | 945 | `PROTO_NC_KQ_W2Z_MAKE_REQ` | `377 bytes` |
| 299 | `PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_ACK` | `16 bytes` | 946 | `PROTO_NC_KQ_W2Z_START_CMD` | `379 bytes` |
| 300 | `PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_ACK` | `16 bytes` | 947 | `PROTO_NC_KQ_WINTER_EVENT_2014_SCORE_CMD` | `4 bytes` |
| 301 | `PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_REQ` | `16 bytes` | 948 | `PROTO_NC_KQ_Z2W_END_CMD` | `4 bytes` |
| 302 | `PROTO_NC_CHAR_GET_ITEMLIST_BY_TYPE_REQ` | `11 bytes` | 949 | `PROTO_NC_KQ_Z2W_MAKE_ACK` | `6 bytes` |
| 303 | `PROTO_NC_CHAR_GUILD_ACADEMY_ZONE_CMD` | `4 bytes` | 950 | `PROTO_NC_LOG_GAME_DATA_TYPE_0` | `68 bytes` |
| 304 | `PROTO_NC_CHAR_ITEM_CMD` | `3 bytes` | 951 | `PROTO_NC_LOG_GAME_DATA_TYPE_1` | `32 bytes` |
| 305 | `PROTO_NC_CHAR_ITEM_CMD::<unnamed-type-PacketOrder>` | `1 bytes` | 952 | `PROTO_NC_LOG_GAME_DATA_TYPE_2` | `40 bytes` |
| 306 | `PROTO_NC_CHAR_KICKPLAYEROUT_RNG` | `28 bytes` | 953 | `PROTO_NC_LOG_GAME_DATA_TYPE_3` | `44 bytes` |
| 307 | `PROTO_NC_CHAR_LEVEL_CHANGED_CMD` | `7 bytes` | 954 | `PROTO_NC_LOG_GAME_DATA_TYPE_4` | `44 bytes` |
| 308 | `PROTO_NC_CHAR_LINKFROM_CMD` | `24 bytes` | 955 | `PROTO_NC_LOG_GAME_DATA_TYPE_5` | `16 bytes` |
| 309 | `PROTO_NC_CHAR_LOGIN_ACK` | `18 bytes` | 956 | `PROTO_NC_LOG_GAME_DATA_TYPE_6` | `44 bytes` |
| 310 | `PROTO_NC_CHAR_MAPLOGIN_ACK` | `242 bytes` | 957 | `PROTO_NC_LOG_GAME_DATA_TYPE_7` | `44 bytes` |
| 311 | `PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD` | `1 bytes` | 958 | `PROTO_NC_LOG_GAME_DATA_TYPE_8` | `36 bytes` |
| 312 | `PROTO_NC_CHAR_PLAYERBANNED_RNG` | `27 bytes` | 959 | `PROTO_NC_LOG_GAME_DATA_TYPE_9` | `48 bytes` |
| 313 | `PROTO_NC_CHAR_PLAYERFOUND_BY_NORMAL_USER_RNG` | `49 bytes` | 960 | `PROTO_NC_LOG_GAME_DATA_TYPE_A` | `36 bytes` |
| 314 | `PROTO_NC_CHAR_PLAYERFOUND_RNG` | `47 bytes` | 961 | `PROTO_NC_LOG_GAME_DATA_TYPE_B` | `60 bytes` |
| 315 | `PROTO_NC_CHAR_PLAYERSEARCH_BY_NORMAL_USER_RNG` | `27 bytes` | 962 | `PROTO_NC_LOG_GAME_DATA_TYPE_C` | `56 bytes` |
| 316 | `PROTO_NC_CHAR_PLAYERSEARCH_RNG` | `27 bytes` | 963 | `PROTO_NC_LOG_GAME_DATA_TYPE_D` | `24 bytes` |
| 317 | `PROTO_NC_CHAR_PLAYERSUMMON_RNG` | `47 bytes` | 964 | `PROTO_NC_LOG_GAME_DATA_TYPE_E` | `48 bytes` |
| 318 | `PROTO_NC_CHAR_POLYMORPH_CMD` | `8 bytes` | 965 | `PROTO_NC_LOG_GAME_DATA_TYPE_F` | `64 bytes` |
| 319 | `PROTO_NC_CHAR_PROMOTE_ACK` | `1 bytes` | 966 | `PROTO_NC_LOG_GAME_DATA_TYPE_G` | `60 bytes` |
| 320 | `PROTO_NC_CHAR_QUEST_DOING_CMD` | `6 bytes` | 967 | `PROTO_NC_LOG_GAME_DATA_TYPE_H` | `36 bytes` |
| 321 | `PROTO_NC_CHAR_QUEST_DONE_CMD` | `12 bytes` | 968 | `PROTO_NC_LOG_GAME_DATA_TYPE_I` | `40 bytes` |
| 322 | `PROTO_NC_CHAR_QUEST_READ_CMD` | `6 bytes` | 969 | `PROTO_NC_LOG_GAME_DATA_TYPE_J` | `52 bytes` |
| 323 | `PROTO_NC_CHAR_QUEST_REPEAT_CMD` | `6 bytes` | 970 | `PROTO_NC_LOG_GAME_DATA_TYPE_K` | `48 bytes` |
| 324 | `PROTO_NC_CHAR_REGISTNUMBER_ACK` | `30 bytes` | 971 | `PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO` | `32 bytes` |
| 325 | `PROTO_NC_CHAR_REGISTNUMBER_REQ` | `22 bytes` | 972 | `PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED` | `16 bytes` |
| 326 | `PROTO_NC_CHAR_REGNUM_VARIFICATION_ACK` | `7 bytes` | 973 | `PROTO_NC_LOG_GAME_ITEM_UPGRADE` | `43 bytes` |
| 327 | `PROTO_NC_CHAR_REGNUM_VARIFICATION_REQ` | `26 bytes` | 974 | `PROTO_NC_LOG_GAME_KQ_ENTER` | `30 bytes` |
| 328 | `PROTO_NC_CHAR_REVIVEOTHER_CMD` | `30 bytes` | 975 | `PROTO_NC_LOG_GAME_KQ_LEAVE` | `31 bytes` |
| 329 | `PROTO_NC_CHAR_REVIVESAME_CMD` | `10 bytes` | 976 | `PROTO_NC_LOG_GAME_LEVEL_UP` | `33 bytes` |
| 330 | `PROTO_NC_CHAR_SELL_ITEM_INFO_ZONE_CMD` | `1 bytes` | 977 | `PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING` | `24 bytes` |
| 331 | `PROTO_NC_CHAR_SINGLE_OPTION_CMD` | `15 bytes` | 978 | `PROTO_NC_LOG_GAME_MINIHOUSE_VISIT` | `28 bytes` |
| 332 | `PROTO_NC_CHAR_SKILLCLIENT_CMD` | `6 bytes` | 979 | `PROTO_NC_LOG_GAME_MK2_DEAD` | `26 bytes` |
| 333 | `PROTO_NC_CHAR_SKILL_CMD` | `9 bytes` | 980 | `PROTO_NC_LOG_GAME_MK2_START` | `28 bytes` |
| 334 | `PROTO_NC_CHAR_SKILL_PASSIVE_CMD` | `2 bytes` | 981 | `PROTO_NC_LOG_GAME_MKED` | `34 bytes` |
| 335 | `PROTO_NC_CHAR_SOMEONEGUILDACADEMYCHANGE_CMD` | `6 bytes` | 982 | `PROTO_NC_LOG_GAME_MK_DROP_ITEM` | `34 bytes` |
| 336 | `PROTO_NC_CHAR_SOMEONEGUILDCHANGE_CMD` | `6 bytes` | 983 | `PROTO_NC_LOG_GAME_MK_GETEXPINFIELD` | `12 bytes` |
| 337 | `PROTO_NC_CHAR_SOMEONEPROMOTE_CMD` | `3 bytes` | 984 | `PROTO_NC_LOG_GAME_PKED` | `33 bytes` |
| 338 | `PROTO_NC_CHAR_STAT_DECPOINTFAIL_DB_ACK` | `10 bytes` | 985 | `PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO` | `41 bytes` |
| 339 | `PROTO_NC_CHAR_STAT_DECPOINTSUC_ACK` | `1 bytes` | 986 | `PROTO_NC_LOG_REGENLOCATESAVE_CMD` | `38 bytes` |
| 340 | `PROTO_NC_CHAR_STAT_DECPOINTSUC_DB_ACK` | `9 bytes` | 987 | `PROTO_NC_LOG_WEDDING_HALL_CANCEL` | `8 bytes` |
| 341 | `PROTO_NC_CHAR_STAT_DECPOINT_REQ` | `1 bytes` | 988 | `PROTO_NC_LOG_WEDDING_HALL_RESERVE` | `8 bytes` |
| 342 | `PROTO_NC_CHAR_STAT_INCPOINTFAIL_DB_ACK` | `8 bytes` | 989 | `PROTO_NC_LOG_WEDDING_HALL_START` | `8 bytes` |
| 343 | `PROTO_NC_CHAR_STAT_INCPOINTSUC_ACK` | `1 bytes` | 990 | `PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD` | `1 bytes` |
| 344 | `PROTO_NC_CHAR_STAT_INCPOINTSUC_DB_ACK` | `7 bytes` | 991 | `PROTO_NC_MAP_EXPBONUS_RNG` | `9 bytes` |
| 345 | `PROTO_NC_CHAR_STAT_INCPOINT_REQ` | `1 bytes` | 992 | `PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD` | `4 bytes` |
| 346 | `PROTO_NC_CHAR_STAT_REMAINPOINT_CMD` | `1 bytes` | 993 | `PROTO_NC_MAP_INDUN_LEVEL_VIEW_CMD` | `4 bytes` |
| 347 | `PROTO_NC_CHAR_SUPPORT_REWARD_CMD` | `8 bytes` | 994 | `PROTO_NC_MAP_ITEMBONUS_RNG` | `9 bytes` |
| 348 | `PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_ACK` | `6 bytes` | 995 | `PROTO_NC_MAP_LINKEND_CMD` | `47 bytes` |
| 349 | `PROTO_NC_CHAR_TUTORIAL_FREESTAT_INIT_REQ` | `10 bytes` | 996 | `PROTO_NC_MAP_LINKRESERVE_ACK` | `36 bytes` |
| 350 | `PROTO_NC_CHAR_TUTORIAL_INFO_ZONE_CMD` | `11 bytes` | 997 | `PROTO_NC_MAP_LINKRESERVE_REQ` | `59 bytes` |
| 351 | `PROTO_NC_CHAR_TUTORIAL_STEP_ACK` | `1 bytes` | 998 | `PROTO_NC_MAP_LINKRESERVE_REQ::<unnamed-type-linkto>` | `20 bytes` |
| 352 | `PROTO_NC_CHAR_TUTORIAL_STEP_REQ` | `1 bytes` | 999 | `PROTO_NC_MAP_LINKSTART_CMD` | `6 bytes` |
| 353 | `PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK` | `9 bytes` | 1000 | `PROTO_NC_MAP_LOGINFAIL_ACK` | `3 bytes` |
| 354 | `PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_CMD` | `9 bytes` | 1001 | `PROTO_NC_MAP_LOGIN_REQ` | `1590 bytes` |
| 355 | `PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_REQ` | `13 bytes` | 1002 | `PROTO_NC_MAP_LOGOUT_CMD` | `2 bytes` |
| 356 | `PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ` | `1 bytes` | 1003 | `PROTO_NC_MAP_MULTY_LINK_CMD` | `73 bytes` |
| 357 | `PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ` | `24 bytes` | 1004 | `PROTO_NC_MAP_MULTY_LINK_SELECT_ACK` | `2 bytes` |
| 358 | `PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD` | `50 bytes` | 1005 | `PROTO_NC_MAP_MULTY_LINK_SELECT_REQ` | `12 bytes` |
| 359 | `PROTO_NC_CHAR_USEITEM_MINIMON_INFO_ZONE_CMD` | `56 bytes` | 1006 | `PROTO_NC_MAP_REGIST_CMD` | `2 bytes` |
| 360 | `PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ` | `1 bytes` | 1007 | `PROTO_NC_MAP_REGIST_CMD::<unnamed-type-maparray>` | `14 bytes` |
| 361 | `PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ` | `24 bytes` | 1008 | `PROTO_NC_MAP_TOWNPORTAL_REQ` | `1 bytes` |
| 362 | `PROTO_NC_CHAR_USEITEM_MINIMON_NOTICE_CMD` | `4 bytes` | 1009 | `PROTO_NC_MAP_WING_FLY_REQ` | `1 bytes` |
| 363 | `PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD` | `3 bytes` | 1010 | `PROTO_NC_MAP_WING_SAVE_REQ` | `1 bytes` |
| 364 | `PROTO_NC_CHAR_WEDDINGDATA_ACK` | `43 bytes` | 1011 | `PROTO_NC_MENU_OPENSTORAGE_CMD` | `12 bytes` |
| 365 | `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner>` | `24 bytes` | 1012 | `PROTO_NC_MENU_SERVERMENU_ACK` | `1 bytes` |
| 366 | `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status>` | `1 bytes` | 1013 | `PROTO_NC_MENU_SERVERMENU_REQ` | `142 bytes` |
| 367 | `PROTO_NC_CHAR_WEDDINGDATA_REQ` | `6 bytes` | 1014 | `PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD` | `24 bytes` |
| 368 | `PROTO_NC_CHAR_WEDDING_CANCEL_WEDDING` | `4 bytes` | 1015 | `PROTO_NC_MENU_SHOPOPENSOULSTONE_CMD::SOULSTONEMENU` | `12 bytes` |
| 369 | `PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_ACK` | `16 bytes` | 1016 | `PROTO_NC_MENU_SHOPOPENTABLE_CMD` | `4 bytes` |
| 370 | `PROTO_NC_CHAR_WEDDING_DIVORCE_CANCEL_REQ` | `11 bytes` | 1017 | `PROTO_NC_MENU_SHOPOPENTABLE_CMD::MENUITEM` | `3 bytes` |
| 371 | `PROTO_NC_CHAR_WEDDING_DIVORCE_DO_REQ` | `11 bytes` | 1018 | `PROTO_NC_MENU_SHOPOPEN_CMD` | `4 bytes` |
| 372 | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK` | `14 bytes` | 1019 | `PROTO_NC_MENU_SHOPOPEN_CMD::MENUITEM` | `2 bytes` |
| 373 | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo>` | `8 bytes` | 1020 | `PROTO_NC_MINIHOUSE_ACTIV_REQ` | `1 bytes` |
| 374 | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ` | `10 bytes` | 1021 | `PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ` | `1 bytes` |
| 375 | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG` | `15 bytes` | 1022 | `PROTO_NC_MINIHOUSE_BUILDING_CMD` | `66 bytes` |
| 376 | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_RNG::<unnamed-type-PartnerInfo>` | `8 bytes` | 1023 | `PROTO_NC_MINIHOUSE_BUILDING_REQ` | `133 bytes` |
| 377 | `PROTO_NC_CHAR_WEDDING_PARTNER_SUMMON_RNG` | `35 bytes` | 1024 | `PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK` | `2 bytes` |
| 378 | `PROTO_NC_CHAR_WEDDING_PROPOSE_ACK` | `14 bytes` | 1025 | `PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD` | `10 bytes` |
| 379 | `PROTO_NC_CHAR_WEDDING_PROPOSE_REQ` | `10 bytes` | 1026 | `PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ` | `4 bytes` |
| 380 | `PROTO_NC_CHAR_WEDDING_REFRESH_INFO_RNG` | `7 bytes` | 1027 | `PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ` | `10 bytes` |
| 381 | `PROTO_NC_CHAR_ZONE_CHARDATA_ACK` | `2 bytes` | 1028 | `PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD` | `1 bytes` |
| 382 | `PROTO_NC_CHAR_ZONE_CHARDATA_REQ` | `22 bytes` | 1029 | `PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD::<unnamed-type-rearrange>` | `21 bytes` |
| 383 | `PROTO_NC_CHAR_ZONE_CHAT_COLOR_CMD` | `8 bytes` | 1030 | `PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ` | `1 bytes` |
| 384 | `PROTO_NC_CHAR_ZONE_GUILD_CMD` | `4 bytes` | 1031 | `PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ::<unnamed-type-rearrange>` | `18 bytes` |
| 385 | `PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK` | `6 bytes` | 1032 | `PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_ACK` | `17 bytes` |
| 386 | `PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_ACK` | `12 bytes` | 1033 | `PROTO_NC_MINIHOUSE_DB_PORTAL_ADD_REQ` | `25 bytes` |
| 387 | `PROTO_NC_COLLECT_BOOKMARK_REGIST_DB_REQ` | `10 bytes` | 1034 | `PROTO_NC_MINIHOUSE_DB_PORTAL_DEL_ACK` | `8 bytes` |
| 388 | `PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ` | `4 bytes` | 1035 | `PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_ACK` | `10 bytes` |
| 389 | `PROTO_NC_COLLECT_CARDOPEN_CMD` | `106 bytes` | 1036 | `PROTO_NC_MINIHOUSE_DB_PORTAL_LIST_REQ` | `6 bytes` |
| 390 | `PROTO_NC_COLLECT_CARDREGIST_ACK` | `10 bytes` | 1037 | `PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_ACK` | `8 bytes` |
| 391 | `PROTO_NC_COLLECT_CARDREGIST_REQ` | `1 bytes` | 1038 | `PROTO_NC_MINIHOUSE_DB_VISITER_COUNT_REQ` | `10 bytes` |
| 392 | `PROTO_NC_COLLECT_PEEPING_COLLECT_ACK` | `28 bytes` | 1039 | `PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ` | `1 bytes` |
| 393 | `PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK` | `28 bytes` | 1040 | `PROTO_NC_MINIHOUSE_FUNICHER_CMD` | `2 bytes` |
| 394 | `PROTO_NC_COLLECT_PEEPING_DB_COLLECT_ACK` | `34 bytes` | 1041 | `PROTO_NC_MINIHOUSE_FUNICHER_CMD::<unnamed-type-array>` | `21 bytes` |
| 395 | `PROTO_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK` | `34 bytes` | 1042 | `PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ` | `1 bytes` |
| 396 | `PROTO_NC_COLLECT_PEEPING_DB_REQ` | `26 bytes` | 1043 | `PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK` | `4 bytes` |
| 397 | `PROTO_NC_COLLECT_PEEPING_REQ` | `20 bytes` | 1044 | `PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD` | `4 bytes` |
| 398 | `PROTO_NC_COLLECT_REGIST_ACK` | `18 bytes` | 1045 | `PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ` | `4 bytes` |
| 399 | `PROTO_NC_COLLECT_REGIST_REQ` | `26 bytes` | 1046 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK` | `6 bytes` |
| 400 | `PROTO_NC_COLLECT_REWARD_ACK` | `4 bytes` | 1047 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK` | `2 bytes` |
| 401 | `PROTO_NC_COLLECT_REWARD_DB_ACK` | `10 bytes` | 1048 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD` | `2 bytes` |
| 402 | `PROTO_NC_COLLECT_REWARD_DB_REQ` | `119 bytes` | 1049 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ` | `1 bytes` |
| 403 | `PROTO_NC_COLLECT_REWARD_REQ` | `4 bytes` | 1050 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD` | `2 bytes` |
| 404 | `PROTO_NC_CT_ADD_FRIEND_CMD` | `8 bytes` | 1051 | `PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ` | `4 bytes` |
| 405 | `PROTO_NC_CT_CHARTTING_CMD` | `6 bytes` | 1052 | `PROTO_NC_MINIHOUSE_FURNITURE_ENDURE_CMD` | `1 bytes` |
| 406 | `PROTO_NC_CT_DB_SET_CMD` | `8 bytes` | 1053 | `PROTO_NC_MINIHOUSE_KICKOUT_CMD` | `1 bytes` |
| 407 | `PROTO_NC_CT_LUASCRIPT_SET_WORLD_CMD` | `6 bytes` | 1054 | `PROTO_NC_MINIHOUSE_KICKOUT_REQ` | `2 bytes` |
| 408 | `PROTO_NC_CT_LUASCRIPT_SET_ZONE_CMD` | `6 bytes` | 1055 | `PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK` | `140 bytes` |
| 409 | `PROTO_NC_CT_SET_CMD` | `2 bytes` | 1056 | `PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK::<unnamed-type-player>` | `3 bytes` |
| 410 | `PROTO_NC_CT_SET_CURRENT_ACK` | `6 bytes` | 1057 | `PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD` | `3 bytes` |
| 411 | `PROTO_NC_CT_SET_CURRENT_DB_CMD` | `14 bytes` | 1058 | `PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD` | `2 bytes` |
| 412 | `PROTO_NC_CT_SET_CURRENT_REQ` | `4 bytes` | 1059 | `PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD` | `1 bytes` |
| 413 | `PROTO_NC_CT_SET_SOMEONECHANGE_CMD` | `6 bytes` | 1060 | `PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ` | `1 bytes` |
| 414 | `PROTO_NC_DICE_TAISAI_BETTING_ACK` | `4 bytes` | 1061 | `PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ` | `1 bytes` |
| 415 | `PROTO_NC_DICE_TAISAI_BETTING_CANCEL_ACK` | `4 bytes` | 1062 | `PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD` | `1 bytes` |
| 416 | `PROTO_NC_DICE_TAISAI_BETTING_CANCEL_CMD` | `10 bytes` | 1063 | `PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ` | `1 bytes` |
| 417 | `PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ` | `10 bytes` | 1064 | `PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD` | `23 bytes` |
| 418 | `PROTO_NC_DICE_TAISAI_BETTING_CMD` | `10 bytes` | 1065 | `PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD` | `1 bytes` |
| 419 | `PROTO_NC_DICE_TAISAI_BETTING_REQ` | `10 bytes` | 1066 | `PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ` | `9 bytes` |
| 420 | `PROTO_NC_DICE_TAISAI_BET_START_CMD` | `7 bytes` | 1067 | `PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ` | `21 bytes` |
| 421 | `PROTO_NC_DICE_TAISAI_BOARD_FOLD_CMD` | `1 bytes` | 1068 | `PROTO_NC_MINIHOUSE_PORTAL_ADD_CMD` | `6 bytes` |
| 422 | `PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK` | `22 bytes` | 1069 | `PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ` | `1 bytes` |
| 423 | `PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD` | `20 bytes` | 1070 | `PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD` | `2 bytes` |
| 424 | `PROTO_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ` | `19 bytes` | 1071 | `PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK` | `2 bytes` |
| 425 | `PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK` | `30 bytes` | 1072 | `PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD` | `4 bytes` |
| 426 | `PROTO_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_REQ` | `26 bytes` | 1073 | `PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ` | `1 bytes` |
| 427 | `PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK` | `32 bytes` | 1074 | `PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD` | `2 bytes` |
| 428 | `PROTO_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_REQ` | `32 bytes` | 1075 | `PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ` | `2 bytes` |
| 429 | `PROTO_NC_DICE_TAISAI_DB_RANK_ACK` | `13 bytes` | 1076 | `PROTO_NC_MINIHOUSE_REARRANGE_REQ` | `1 bytes` |
| 430 | `PROTO_NC_DICE_TAISAI_DB_RANK_REQ` | `18 bytes` | 1077 | `PROTO_NC_MINIHOUSE_REARRANGE_REQ::<unnamed-type-rearrange>` | `18 bytes` |
| 431 | `PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK` | `28 bytes` | 1078 | `PROTO_NC_MINIHOUSE_VISITREADY_ACK` | `7 bytes` |
| 432 | `PROTO_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_REQ` | `18 bytes` | 1079 | `PROTO_NC_MINIHOUSE_VISITREADY_REQ` | `2 bytes` |
| 433 | `PROTO_NC_DICE_TAISAI_DICE_ROLL_ACK` | `2 bytes` | 1080 | `PROTO_NC_MINIHOUSE_VISIT_REQ` | `11 bytes` |
| 434 | `PROTO_NC_DICE_TAISAI_DICE_ROLL_CMD` | `2 bytes` | 1081 | `PROTO_NC_MISC_CHAT_BLOCK_SPAMER_BLOCKTIME_CMD` | `4 bytes` |
| 435 | `PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD` | `5 bytes` | 1082 | `PROTO_NC_MISC_CONNECTFROMWHERE_ACK` | `7 bytes` |
| 436 | `PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD` | `1 bytes` | 1083 | `PROTO_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD` | `32 bytes` |
| 437 | `PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_ACK` | `16 bytes` | 1084 | `PROTO_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD` | `32 bytes` |
| 438 | `PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ` | `2 bytes` | 1085 | `PROTO_NC_MISC_GM_CHAT_COLOR_REQ` | `4 bytes` |
| 439 | `PROTO_NC_DICE_TAISAI_GAME_BETTING_INFO_CMD` | `3 bytes` | 1086 | `PROTO_NC_MISC_HEARTBEAT_REQ_SEND` | `3 bytes` |
| 440 | `PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK` | `50 bytes` | 1087 | `PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB` | `90 bytes` |
| 441 | `PROTO_NC_DICE_TAISAI_GAME_JOIN_CMD` | `2 bytes` | 1088 | `PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB` | `130 bytes` |
| 442 | `PROTO_NC_DICE_TAISAI_GAME_LEAVE_ACK` | `2 bytes` | 1089 | `PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE` | `34 bytes` |
| 443 | `PROTO_NC_DICE_TAISAI_GAME_LEAVE_CMD` | `3 bytes` | 1090 | `PROTO_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB` | `58 bytes` |
| 444 | `PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_ACK` | `6 bytes` | 1091 | `PROTO_NC_MISC_RESTMINUTE_CMD` | `3 bytes` |
| 445 | `PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_CMD` | `4 bytes` | 1092 | `PROTO_NC_MISC_S2SCONNECTION_ACK` | `3 bytes` |
| 446 | `PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_CMD` | `1 bytes` | 1093 | `PROTO_NC_MISC_S2SCONNECTION_REQ` | `7 bytes` |
| 447 | `PROTO_NC_DICE_TAISAI_RANK_ACK` | `484 bytes` | 1094 | `PROTO_NC_MISC_SEED_ACK` | `2 bytes` |
| 448 | `PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_ACK` | `12 bytes` | 1095 | `PROTO_NC_MISC_SERVERPARAMETER_ACK` | `3 bytes` |
| 449 | `PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ` | `2 bytes` | 1096 | `PROTO_NC_MISC_SPAMMER_RELEASE_CHAT_BAN_REQ` | `24 bytes` |
| 450 | `PROTO_NC_DICE_TAISAI_TIMER_CMD` | `7 bytes` | 1097 | `PROTO_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD` | `24 bytes` |
| 451 | `PROTO_NC_EVENT_ADD_UPDATE_EVENT_CMD` | `1363 bytes` | 1098 | `PROTO_NC_MISC_TIMEFROMWORLD_CMD` | `36 bytes` |
| 452 | `PROTO_NC_EVENT_DEL_EVENT_CMD` | `2 bytes` | 1099 | `PROTO_NC_MISC_WHSHANDLEREPAIR_CMD` | `11 bytes` |
| 453 | `PROTO_NC_FRIEND_SOMEONE_GET_SPECIALITEM_ZONE_CMD` | `20 bytes` | 1100 | `PROTO_NC_MISC_ZONERINGLINKTEST_RNG` | `8 bytes` |
| 454 | `PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK` | `7 bytes` | 1101 | `PROTO_NC_MOVER_HUNGRY_CMD` | `2 bytes` |
| 455 | `PROTO_NC_GAMBLE_COIN_BUY_REQ` | `8 bytes` | 1102 | `PROTO_NC_MOVER_MOVESPEED_CMD` | `6 bytes` |
| 456 | `PROTO_NC_GAMBLE_COIN_CHANGE_CMD` | `8 bytes` | 1103 | `PROTO_NC_MOVER_RIDE_ON_CMD` | `4 bytes` |
| 457 | `PROTO_NC_GAMBLE_COIN_DB_BUY_ACK` | `10 bytes` | 1104 | `PROTO_NC_MOVER_RIDE_ON_REQ` | `3 bytes` |
| 458 | `PROTO_NC_GAMBLE_COIN_DB_BUY_REQ` | `32 bytes` | 1105 | `PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ` | `10 bytes` |
| 459 | `PROTO_NC_GAMBLE_COIN_DB_SELL_ACK` | `10 bytes` | 1106 | `PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ` | `4 bytes` |
| 460 | `PROTO_NC_GAMBLE_COIN_DB_SELL_REQ` | `24 bytes` | 1107 | `PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD` | `2 bytes` |
| 461 | `PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK` | `18 bytes` | 1108 | `PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD` | `6 bytes` |
| 462 | `PROTO_NC_GAMBLE_COIN_DB_USE_COINITEM_REQ` | `24 bytes` | 1109 | `PROTO_NC_PARTY_CHANGEMASTER_CMD` | `6 bytes` |
| 463 | `PROTO_NC_GAMBLE_COIN_USE_COINITEM_MESSAGE_CMD` | `8 bytes` | 1110 | `PROTO_NC_PARTY_DISMISS_CMD` | `2 bytes` |
| 464 | `PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ` | `1 bytes` | 1111 | `PROTO_NC_PARTY_FUNDAMENTAL_CMD` | `66 bytes` |
| 465 | `PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK` | `3 bytes` | 1112 | `PROTO_NC_PARTY_ITEM_LOOTING_ZONE_CMD` | `8 bytes` |
| 466 | `PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ` | `1 bytes` | 1113 | `PROTO_NC_PARTY_JOIN_CMD` | `34 bytes` |
| 467 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK` | `4 bytes` | 1114 | `PROTO_NC_PARTY_KICKOFF_CMD` | `6 bytes` |
| 468 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK` | `4 bytes` | 1115 | `PROTO_NC_PARTY_LEAVE_CMD` | `33 bytes` |
| 469 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD` | `10 bytes` | 1116 | `PROTO_NC_PARTY_LOGIN_CMD` | `6 bytes` |
| 470 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ` | `10 bytes` | 1117 | `PROTO_NC_PARTY_LOGOUT_CMD` | `6 bytes` |
| 471 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD` | `10 bytes` | 1118 | `PROTO_NC_PARTY_MEMBERCLASS_CMD` | `1 bytes` |
| 472 | `PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ` | `10 bytes` | 1119 | `PROTO_NC_PARTY_MEMBERINFORM_ACK` | `3 bytes` |
| 473 | `PROTO_NC_GAMBLE_DICE_TAISAI_BET_START_CMD` | `7 bytes` | 1120 | `PROTO_NC_PARTY_MEMBERINFORM_CMD` | `1 bytes` |
| 474 | `PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK` | `3 bytes` | 1121 | `PROTO_NC_PARTY_MEMBERINFORM_REQ` | `2 bytes` |
| 475 | `PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ` | `1 bytes` | 1122 | `PROTO_NC_PARTY_MEMBERLOCATION_CMD` | `1 bytes` |
| 476 | `PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK` | `9 bytes` | 1123 | `PROTO_NC_PARTY_MEMBERLOCATION_CMD::<unnamed-type-member>` | `28 bytes` |
| 477 | `PROTO_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_REQ` | `6 bytes` | 1124 | `PROTO_NC_PARTY_MEMBERMAPOUT` | `20 bytes` |
| 478 | `PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK` | `14 bytes` | 1125 | `PROTO_NC_PARTY_ZONE_JOIN_CMD` | `8 bytes` |
| 479 | `PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_LOG_CMD` | `12 bytes` | 1126 | `PROTO_NC_PARTY_ZONE_LEAVE_CMD` | `4 bytes` |
| 480 | `PROTO_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_REQ` | `11 bytes` | 1127 | `PROTO_NC_PARTY_ZONE_SET_LOOTER_CMD` | `6 bytes` |
| 481 | `PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_CMD` | `2 bytes` | 1128 | `PROTO_NC_PET_ASK_NEW_NAME_ACK` | `4 bytes` |
| 482 | `PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD` | `5 bytes` | 1129 | `PROTO_NC_PET_ASK_NEW_NAME_REQ` | `7 bytes` |
| 483 | `PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_EMOTION_CMD` | `1 bytes` | 1130 | `PROTO_NC_PET_CREATE_DB_ACK` | `8 bytes` |
| 484 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_BETTING_INFO_CMD` | `3 bytes` | 1131 | `PROTO_NC_PET_CREATE_DB_REQ` | `24 bytes` |
| 485 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK` | `50 bytes` | 1132 | `PROTO_NC_PET_LINK_RESUMMON_CMD` | `23 bytes` |
| 486 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD` | `2 bytes` | 1133 | `PROTO_NC_PET_LOAD_INFO_DB_ACK` | `26 bytes` |
| 487 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ` | `4 bytes` | 1134 | `PROTO_NC_PET_LOAD_INFO_DB_REQ` | `6 bytes` |
| 488 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK` | `2 bytes` | 1135 | `PROTO_NC_PET_SET_NAME_ACK` | `6 bytes` |
| 489 | `PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD` | `3 bytes` | 1136 | `PROTO_NC_PET_SET_NAME_CANCEL_ACK` | `6 bytes` |
| 490 | `PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD` | `1 bytes` | 1137 | `PROTO_NC_PET_SET_NAME_CANCEL_REQ` | `2 bytes` |
| 491 | `PROTO_NC_GAMBLE_DICE_TAISAI_TIMER_CMD` | `7 bytes` | 1138 | `PROTO_NC_PET_SET_NAME_CMD` | `18 bytes` |
| 492 | `PROTO_NC_GAMBLE_ENTER_PLAYER_DIRECT_CMD` | `3 bytes` | 1139 | `PROTO_NC_PET_SET_NAME_DB_ACK` | `8 bytes` |
| 493 | `PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD` | `8 bytes` | 1140 | `PROTO_NC_PET_SET_NAME_DB_REQ` | `22 bytes` |
| 494 | `PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK` | `12 bytes` | 1141 | `PROTO_NC_PET_SET_NAME_REQ` | `18 bytes` |
| 495 | `PROTO_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_REQ` | `8 bytes` | 1142 | `PROTO_NC_PET_SET_TENDENCY_DB_REQ` | `8 bytes` |
| 496 | `PROTO_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD` | `4 bytes` | 1143 | `PROTO_NC_PET_USE_ITEM_FAIL_ACK` | `2 bytes` |
| 497 | `PROTO_NC_GAMBLE_PLAYERACT_CMD` | `6 bytes` | 1144 | `PROTO_NC_PRISON_END_ACK` | `6 bytes` |
| 498 | `PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK` | `8 bytes` | 1145 | `PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK` | `2 bytes` |
| 499 | `PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ` | `1 bytes` | 1146 | `PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ` | `2 bytes` |
| 500 | `PROTO_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_REQ::<unnamed-type-MachineIndex>` | `3 bytes` | 1147 | `PROTO_NC_QUEST_DB_CLEAR_ACK` | `10 bytes` |
| 501 | `PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK` | `10 bytes` | 1148 | `PROTO_NC_QUEST_DB_CLEAR_REQ` | `8 bytes` |
| 502 | `PROTO_NC_GAMBLE_SLOTMACHINE_DB_RESULT_REQ` | `36 bytes` | 1149 | `PROTO_NC_QUEST_DB_DONE_ACK` | `43 bytes` |
| 503 | `PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK` | `4 bytes` | 1150 | `PROTO_NC_QUEST_DB_DONE_REQ` | `83 bytes` |
| 504 | `PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ` | `4 bytes` | 1151 | `PROTO_NC_QUEST_DB_GIVE_UP_ACK` | `10 bytes` |
| 505 | `PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK` | `2 bytes` | 1152 | `PROTO_NC_QUEST_DB_GIVE_UP_REQ` | `8 bytes` |
| 506 | `PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD` | `1 bytes` | 1153 | `PROTO_NC_QUEST_DB_READ_ACK` | `10 bytes` |
| 507 | `PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD::<unnamed-type-JackPotList>` | `6 bytes` | 1154 | `PROTO_NC_QUEST_DB_READ_REQ` | `8 bytes` |
| 508 | `PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK` | `1 bytes` | 1155 | `PROTO_NC_QUEST_DB_SET_INFO_ACK` | `41 bytes` |
| 509 | `PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD` | `28 bytes` | 1156 | `PROTO_NC_QUEST_DB_SET_INFO_REQ` | `39 bytes` |
| 510 | `PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK` | `2 bytes` | 1157 | `PROTO_NC_QUEST_GIVE_UP_ACK` | `4 bytes` |
| 511 | `PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ` | `4 bytes` | 1158 | `PROTO_NC_QUEST_GIVE_UP_REQ` | `2 bytes` |
| 512 | `PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK` | `11 bytes` | 1159 | `PROTO_NC_QUEST_JOBDUNGEON_FIND_RNG` | `115 bytes` |
| 513 | `PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK` | `10 bytes` | 1160 | `PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD` | `1 bytes` |
| 514 | `PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK` | `1 bytes` | 1161 | `PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest` | `4 bytes` |
| 515 | `PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK` | `9 bytes` | 1162 | `PROTO_NC_QUEST_READ_ACK` | `4 bytes` |
| 516 | `PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ` | `2 bytes` | 1163 | `PROTO_NC_QUEST_READ_REQ` | `2 bytes` |
| 517 | `PROTO_NC_GAMBLE_WORLD_PREVMAPNAME_CMD` | `18 bytes` | 1164 | `PROTO_NC_QUEST_RESET_TIME_CMD` | `16 bytes` |
| 518 | `PROTO_NC_GAMBLE_ZONE_PREVMAPNAME_CMD` | `12 bytes` | 1165 | `PROTO_NC_QUEST_RESET_TIME_ZONE_CMD` | `20 bytes` |
| 519 | `PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_ACK` | `61 bytes` | 1166 | `PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD` | `2 bytes` |
| 520 | `PROTO_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_REQ` | `15 bytes` | 1167 | `PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD` | `6 bytes` |
| 521 | `PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_GRADE_INFO_CMD` | `9 bytes` | 1168 | `PROTO_NC_QUEST_SCENARIO_RUN_CMD` | `2 bytes` |
| 522 | `PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_ITEM_INFO_ZONE_RNG` | `176 bytes` | 1169 | `PROTO_NC_QUEST_SCRIPT_CMD_ACK` | `7 bytes` |
| 523 | `PROTO_NC_GUILD_ACADEMY_REWARDSTORAGE_REWARD_ZONE_RNG` | `48 bytes` | 1170 | `PROTO_NC_QUEST_SCRIPT_CMD_REQ` | `103 bytes` |
| 524 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD` | `1 bytes` | 1171 | `PROTO_NC_QUEST_SELECT_START_ACK` | `6 bytes` |
| 525 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_CMD` | `41 bytes` | 1172 | `PROTO_NC_QUEST_SELECT_START_REQ` | `4 bytes` |
| 526 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_RNG` | `53 bytes` | 1173 | `PROTO_NC_QUEST_START_ACK` | `2 bytes` |
| 527 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK` | `11 bytes` | 1174 | `PROTO_NC_QUEST_START_REQ` | `2 bytes` |
| 528 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ` | `1 bytes` | 1175 | `PROTO_NC_RAID_LOCATION_CMD` | `1 bytes` |
| 529 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_CMD` | `41 bytes` | 1176 | `PROTO_NC_RAID_LOCATION_CMD::Memberinfo` | `12 bytes` |
| 530 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_ACK` | `9 bytes` | 1177 | `PROTO_NC_RAID_MAPOUT_CMD` | `20 bytes` |
| 531 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_REQ` | `6 bytes` | 1178 | `PROTO_NC_RAID_MEMBERCLASS_CMD` | `1 bytes` |
| 532 | `PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_RNG` | `53 bytes` | 1179 | `PROTO_NC_RAID_MEMBERINFORM_CMD` | `16 bytes` |
| 533 | `PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_ACADEMY_MASTER_BUFF_CMD` | `4 bytes` | 1180 | `PROTO_NC_RAID_ZONE_DEL_SUB_LEADER_CMD` | `10 bytes` |
| 534 | `PROTO_NC_GUILD_ACADEMY_ZONE_GUILD_BUFF_CMD` | `36 bytes` | 1181 | `PROTO_NC_RAID_ZONE_DISMISS_CMD` | `2 bytes` |
| 535 | `PROTO_NC_GUILD_ACADEMY_ZONE_MASTER_TELEPORT_CMD` | `24 bytes` | 1182 | `PROTO_NC_RAID_ZONE_LEAVE_CMD` | `10 bytes` |
| 536 | `PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_GUILD_JOIN_CMD` | `8 bytes` | 1183 | `PROTO_NC_RAID_ZONE_LOOTING_TYPE_CMD` | `4 bytes` |
| 537 | `PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_JOIN_CMD` | `8 bytes` | 1184 | `PROTO_NC_RAID_ZONE_MAKE_CMD` | `32 bytes` |
| 538 | `PROTO_NC_GUILD_ACADEMY_ZONE_MEMBER_LEAVE_CMD` | `4 bytes` | 1185 | `PROTO_NC_RAID_ZONE_MEMBER_JOIN_CMD` | `42 bytes` |
| 539 | `PROTO_NC_GUILD_EMBLEM_LEVELUP_CMD` | `4 bytes` | 1186 | `PROTO_NC_RAID_ZONE_SET_LEADER_CMD` | `10 bytes` |
| 540 | `PROTO_NC_GUILD_GUILDWARCONFIRM_ACK` | `16 bytes` | 1187 | `PROTO_NC_RAID_ZONE_SET_LOOTER_CMD` | `10 bytes` |
| 541 | `PROTO_NC_GUILD_GUILDWARCONFIRM_REQ` | `8 bytes` | 1188 | `PROTO_NC_RAID_ZONE_SET_SUB_LEADER_CMD` | `10 bytes` |
| 542 | `PROTO_NC_GUILD_GUILDWARSTATUS_ACK` | `2 bytes` | 1189 | `PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD` | `50 bytes` |
| 543 | `PROTO_NC_GUILD_MOBGUILD_CMD` | `6 bytes` | 1190 | `PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD::<unnamed-type-PostPosition>` | `24 bytes` |
| 544 | `PROTO_NC_GUILD_RENAME_REQ` | `17 bytes` | 1191 | `PROTO_NC_RAID_ZONE_SLOT_MOVE_CMD::<unnamed-type-PrePosition>` | `24 bytes` |
| 545 | `PROTO_NC_GUILD_RETYPE_REQ` | `2 bytes` | 1192 | `PROTO_NC_RAID_ZONE_VANISH_CMD` | `10 bytes` |
| 546 | `PROTO_NC_GUILD_STORAGEOPEN_ACK` | `19 bytes` | 1193 | `PROTO_NC_SCENARIO_ANIMATE_CMD` | `34 bytes` |
| 547 | `PROTO_NC_GUILD_STORAGEWITHDRAW_CMD` | `33 bytes` | 1194 | `PROTO_NC_SCENARIO_AREAENTRY_ACK` | `32 bytes` |
| 548 | `PROTO_NC_GUILD_STORAGEWITHDRAW_RNG` | `44 bytes` | 1195 | `PROTO_NC_SCENARIO_AREAENTRY_REQ` | `32 bytes` |
| 549 | `PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_ACK` | `9 bytes` | 1196 | `PROTO_NC_SCENARIO_AREALEAVE_ACK` | `32 bytes` |
| 550 | `PROTO_NC_GUILD_STORAGE_WITHDRAW_GRADE_REQ` | `6 bytes` | 1197 | `PROTO_NC_SCENARIO_AREALEAVE_REQ` | `32 bytes` |
| 551 | `PROTO_NC_GUILD_TOURNAMENTSTART_CMD` | `14 bytes` | 1198 | `PROTO_NC_SCENARIO_BRIEFEFFECT_CMD` | `1 bytes` |
| 552 | `PROTO_NC_GUILD_TOURNAMENTSTOPMSG_CMD` | `20 bytes` | 1199 | `PROTO_NC_SCENARIO_CAMERA_MOVE_CMD` | `15 bytes` |
| 553 | `PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_ACK` | `35 bytes` | 1200 | `PROTO_NC_SCENARIO_CHATWIN_CMD` | `3 bytes` |
| 554 | `PROTO_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_REQ` | `33 bytes` | 1201 | `PROTO_NC_SCENARIO_DEL_DIRECTIONAL_ARROW_CMD` | `1 bytes` |
| 555 | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_BEFORE_END_TIME_MSG_CMD` | `4 bytes` | 1202 | `PROTO_NC_SCENARIO_DIALOG_CMD` | `34 bytes` |
| 556 | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_START_CMD` | `52 bytes` | 1203 | `PROTO_NC_SCENARIO_DIRECTIONAL_ARROW_CMD` | `8 bytes` |
| 557 | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK` | `2 bytes` | 1204 | `PROTO_NC_SCENARIO_DOORSTATE_CMD` | `3 bytes` |
| 558 | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_CMD` | `8 bytes` | 1205 | `PROTO_NC_SCENARIO_EFFECT_TIMER_CMD` | `3 bytes` |
| 559 | `PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ` | `1 bytes` | 1206 | `PROTO_NC_SCENARIO_FOG_CMD` | `3 bytes` |
| 560 | `PROTO_NC_GUILD_TOURNAMENT_END_CMD` | `12 bytes` | 1207 | `PROTO_NC_SCENARIO_HIDE_OTHER_PLAYER_CMD` | `1 bytes` |
| 561 | `PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK` | `2 bytes` | 1208 | `PROTO_NC_SCENARIO_LIGHT_CMD` | `3 bytes` |
| 562 | `PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_CMD` | `4 bytes` | 1209 | `PROTO_NC_SCENARIO_MAPMARK_CMD` | `2 bytes` |
| 563 | `PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ` | `2 bytes` | 1210 | `PROTO_NC_SCENARIO_MAPMARK_CMD::IconInfo` | `15 bytes` |
| 564 | `PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_RESULT_CMD` | `48 bytes` | 1211 | `PROTO_NC_SCENARIO_MAP_OBJECT_CONTROL_CMD` | `37 bytes` |
| 565 | `PROTO_NC_GUILD_TOURNAMENT_GOLD_REFUND_ZONE_CMD` | `3 bytes` | 1212 | `PROTO_NC_SCENARIO_MESSAGENOTICE_CMD` | `4 bytes` |
| 566 | `PROTO_NC_GUILD_TOURNAMENT_ITEM_EFFECT_CMD` | `2 bytes` | 1213 | `PROTO_NC_SCENARIO_MESSAGE_CMD` | `17 bytes` |
| 567 | `PROTO_NC_GUILD_TOURNAMENT_ITEM_FLAG_CMD` | `5 bytes` | 1214 | `PROTO_NC_SCENARIO_NPCCHAT_CMD` | `35 bytes` |
| 568 | `PROTO_NC_GUILD_TOURNAMENT_ITEM_PICK_CMD` | `2 bytes` | 1215 | `PROTO_NC_SCENARIO_OBJTYPECHANGE_CMD` | `3 bytes` |
| 569 | `PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD` | `2 bytes` | 1216 | `PROTO_NC_SCENARIO_RANKING_LIST_CMD` | `1000 bytes` |
| 570 | `PROTO_NC_GUILD_TOURNAMENT_ITEM_SCAN_CMD::Memberinfo` | `8 bytes` | 1217 | `PROTO_NC_SCENARIO_SCRIPTMESSAGE_CMD` | `37 bytes` |
| 571 | `PROTO_NC_GUILD_TOURNAMENT_MANAGERUSERMSG_CMD` | `20 bytes` | 1218 | `PROTO_NC_SCENARIO_SYSTEM_MESSAGE_CMD` | `1 bytes` |
| 572 | `PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK` | `11 bytes` | 1219 | `PROTO_NC_SCENARIO_TIMER_CMD` | `2 bytes` |
| 573 | `PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ` | `1 bytes` | 1220 | `PROTO_NC_SCENARIO_TIMER_END_CMD` | `4 bytes` |
| 574 | `PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ` | `1 bytes` | 1221 | `PROTO_NC_SCENARIO_TIMER_START_CMD` | `1 bytes` |
| 575 | `PROTO_NC_GUILD_TOURNAMENT_PLAYERDIEMSG_CMD` | `44 bytes` | 1222 | `PROTO_NC_SCENARIO_TOPVIEW_CMD` | `13 bytes` |
| 576 | `PROTO_NC_GUILD_TOURNAMENT_PLAYERKILLMSG_CMD` | `44 bytes` | 1223 | `PROTO_NC_SKILL_COOLTIME_CMD` | `6 bytes` |
| 577 | `PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK` | `2 bytes` | 1224 | `PROTO_NC_SKILL_EMPOWALLOC_ACK` | `4 bytes` |
| 578 | `PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_CMD` | `1 bytes` | 1225 | `PROTO_NC_SKILL_EMPOWALLOC_DB_ACK` | `13 bytes` |
| 579 | `PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ` | `1 bytes` | 1226 | `PROTO_NC_SKILL_EMPOWALLOC_DB_REQ` | `40 bytes` |
| 580 | `PROTO_NC_GUILD_TOURNAMENT_SCORE_CMD` | `17 bytes` | 1227 | `PROTO_NC_SKILL_EMPOWALLOC_REQ` | `6 bytes` |
| 581 | `PROTO_NC_GUILD_TOURNAMENT_SKILLPOINT_CMD` | `2 bytes` | 1228 | `PROTO_NC_SKILL_EMPOWPOINT_CMD` | `1 bytes` |
| 582 | `PROTO_NC_GUILD_TOURNAMENT_STARTMSG_CMD` | `4 bytes` | 1229 | `PROTO_NC_SKILL_EMPOW_RESET_DB_REQ` | `6 bytes` |
| 583 | `PROTO_NC_GUILD_TOURNAMENT_START_CMD` | `17 bytes` | 1230 | `PROTO_NC_SKILL_EMPOW_RESET_SUC_ACK` | `1 bytes` |
| 584 | `PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK` | `4 bytes` | 1231 | `PROTO_NC_SKILL_ERASE_ACK` | `10 bytes` |
| 585 | `PROTO_NC_GUILD_TOURNAMENT_USESKILL_CMD` | `6 bytes` | 1232 | `PROTO_NC_SKILL_ERASE_REQ` | `8 bytes` |
| 586 | `PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ` | `2 bytes` | 1233 | `PROTO_NC_SKILL_JUMP_CMD` | `21 bytes` |
| 587 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_ENTER_CMD` | `9 bytes` | 1234 | `PROTO_NC_SKILL_PASSIVESKILL_LEARN_CMD` | `8 bytes` |
| 588 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_FIGHTER_OUT_CMD` | `9 bytes` | 1235 | `PROTO_NC_SKILL_PRODUCTFIELD_ACK` | `4 bytes` |
| 589 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_ACK` | `27 bytes` | 1236 | `PROTO_NC_SKILL_PRODUCTFIELD_REQ` | `2 bytes` |
| 590 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_REQ` | `25 bytes` | 1237 | `PROTO_NC_SKILL_REVIVE_CMD` | `12 bytes` |
| 591 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_ACK` | `9 bytes` | 1238 | `PROTO_NC_SKILL_SKILLEXP_CLIENT_CMD` | `6 bytes` |
| 592 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_REQ` | `5 bytes` | 1239 | `PROTO_NC_SKILL_SKILLTEACHSUC_ACK` | `3 bytes` |
| 593 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_ACK` | `5 bytes` | 1240 | `PROTO_NC_SKILL_SKILLTEACH_REQ` | `9 bytes` |
| 594 | `PROTO_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_REQ` | `3 bytes` | 1241 | `PROTO_NC_SKILL_SOMEONEREVIVE_CMD` | `14 bytes` |
| 595 | `PROTO_NC_GUILD_WORLD_RENAME_ACK` | `25 bytes` | 1242 | `PROTO_NC_SKILL_UNLEARN_ACK` | `4 bytes` |
| 596 | `PROTO_NC_GUILD_WORLD_RENAME_CMD` | `22 bytes` | 1243 | `PROTO_NC_SKILL_UNLEARN_REQ` | `3 bytes` |
| 597 | `PROTO_NC_GUILD_WORLD_RENAME_REQ` | `23 bytes` | 1244 | `PROTO_NC_SKILL_WARP_CMD` | `20 bytes` |
| 598 | `PROTO_NC_GUILD_WORLD_RETYPE_ACK` | `10 bytes` | 1245 | `PROTO_NC_SOULSTONE_HP_BUY_REQ` | `2 bytes` |
| 599 | `PROTO_NC_GUILD_WORLD_RETYPE_CMD` | `7 bytes` | 1246 | `PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD` | `2 bytes` |
| 600 | `PROTO_NC_GUILD_WORLD_RETYPE_REQ` | `8 bytes` | 1247 | `PROTO_NC_SOULSTONE_SP_BUY_ACK` | `2 bytes` |
| 601 | `PROTO_NC_GUILD_WORLD_USE_GUILD_TOKEN_CMD` | `20 bytes` | 1248 | `PROTO_NC_SOULSTONE_SP_BUY_REQ` | `2 bytes` |
| 602 | `PROTO_NC_GUILD_ZONE_DELETE_CMD` | `4 bytes` | 1249 | `PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD` | `2 bytes` |
| 603 | `PROTO_NC_GUILD_ZONE_LIST_ACK` | `3 bytes` | 1250 | `PROTO_NC_SYSLOG_CHAR_DEATH` | `22 bytes` |
| 604 | `PROTO_NC_GUILD_ZONE_LIST_REQ` | `1 bytes` | 1251 | `PROTO_NC_SYSLOG_CHAR_QUEST_STARTED` | `13 bytes` |
| 605 | `PROTO_NC_GUILD_ZONE_MAKE_CMD` | `53 bytes` | 1252 | `PROTO_NC_SYSLOG_CHAR_VICTORY` | `22 bytes` |
| 606 | `PROTO_NC_GUILD_ZONE_MEMBER_JOIN_CMD` | `8 bytes` | 1253 | `PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION` | `11 bytes` |
| 607 | `PROTO_NC_GUILD_ZONE_MEMBER_LEAVE_CMD` | `4 bytes` | 1254 | `PROTO_NC_TRADE_CENBOARDING_ACK` | `8 bytes` |
| 608 | `PROTO_NC_GUILD_ZONE_USE_GUILD_TOKEN_CMD` | `12 bytes` | 1255 | `PROTO_NC_TRADE_DOWNBOARD_ACK` | `1 bytes` |
| 609 | `PROTO_NC_GUILD_ZONE_WAR_END_CMD` | `8 bytes` | 1256 | `PROTO_NC_TRADE_DOWNBOARD_REQ` | `1 bytes` |
| 610 | `PROTO_NC_GUILD_ZONE_WAR_KILL_CMD` | `4 bytes` | 1257 | `PROTO_NC_TRADE_OPPOSITCENBOARDING_CMD` | `8 bytes` |
| 611 | `PROTO_NC_GUILD_ZONE_WAR_START_CMD` | `8 bytes` | 1258 | `PROTO_NC_TRADE_OPPOSITDOWNBOARD_CMD` | `1 bytes` |
| 612 | `PROTO_NC_HOLY_PROMISE_CENTRANSFER_RNG` | `35 bytes` | 1259 | `PROTO_NC_TRADE_OPPOSITUPBOARD_CMD` | `104 bytes` |
| 613 | `PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK` | `16 bytes` | 1260 | `PROTO_NC_TRADE_PROPOSE_ASK_REQ` | `2 bytes` |
| 614 | `PROTO_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_REQ` | `6 bytes` | 1261 | `PROTO_NC_TRADE_PROPOSE_REQ` | `2 bytes` |
| 615 | `PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK` | `16 bytes` | 1262 | `PROTO_NC_TRADE_START_CMD` | `2 bytes` |
| 616 | `PROTO_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_REQ` | `6 bytes` | 1263 | `PROTO_NC_TRADE_UPBOARD_ACK` | `2 bytes` |
| 617 | `PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK` | `10 bytes` | 1264 | `PROTO_NC_TRADE_UPBOARD_REQ` | `1 bytes` |
| 618 | `PROTO_NC_HOLY_PROMISE_MYUPPER_ACK` | `10 bytes` | 1265 | `PROTO_NC_USER_CONNECTCUT2WORLDMANAGER_CMD` | `8 bytes` |
| 619 | `PROTO_NC_HOLY_PROMISE_MY_UP_ZONE` | `8 bytes` | 1266 | `PROTO_NC_USER_CONNECTCUT2ZONE_CMD` | `4 bytes` |
| 620 | `PROTO_NC_HOLY_PROMISE_REWARD_MONEY_CMD` | `28 bytes` | 1267 | `PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK` | `4 bytes` |
| 621 | `PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK` | `10 bytes` | 1268 | `PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ` | `1 bytes` |
| 622 | `PROTO_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD` | `9 bytes` | 1269 | `PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK` | `12 bytes` |
| 623 | `PROTO_NC_INSTANCE_DUNGEON_EMPTY_DUNGEON_CMD` | `5 bytes` | 1270 | `PROTO_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_REQ` | `38 bytes` |
| 624 | `PROTO_NC_INSTANCE_DUNGEON_FIND_RNG` | `70 bytes` | 1271 | `PROTO_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD` | `45 bytes` |
| 625 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK` | `12 bytes` | 1272 | `PROTO_NC_USER_NORMALLOGOUT_CMD` | `1 bytes` |
| 626 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_ACK` | `13 bytes` | 1273 | `PROTO_NC_USER_POSSIBLE_NEW_CONNECT_CMD` | `6 bytes` |
| 627 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ` | `35 bytes` | 1274 | `PROTO_NC_USER_REGISENUMBER_ACK` | `4 bytes` |
| 628 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_REQ` | `36 bytes` | 1275 | `PROTO_NC_USER_USE_BEAUTY_SHOP_CMD` | `1 bytes` |
| 629 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK` | `24 bytes` | 1276 | `PROTO_NC_WT_GRADE_CMD` | `3 bytes` |
| 630 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_ACK` | `25 bytes` | 1277 | `PROTO_NC_WT_LICENSE_ACK` | `7 bytes` |
| 631 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ` | `39 bytes` | 1278 | `PROTO_NC_WT_LICENSE_CLR_ACK` | `24 bytes` |
| 632 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_REQ` | `40 bytes` | 1279 | `PROTO_NC_WT_LICENSE_CLR_ACK::echo_data` | `4 bytes` |
| 633 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK` | `6 bytes` | 1280 | `PROTO_NC_WT_LICENSE_CLR_REQ` | `35 bytes` |
| 634 | `PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_REQ` | `15 bytes` | 1281 | `PROTO_NC_WT_LICENSE_CLR_REQ::echo_data` | `4 bytes` |
| 635 | `PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD` | `2 bytes` | 1282 | `PROTO_NC_WT_LICENSE_REQ` | `7 bytes` |
| 636 | `PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK` | `2 bytes` | 1283 | `PROTO_NC_WT_LICENSE_SET_ACK` | `14 bytes` |
| 637 | `PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD` | `5 bytes` | 1284 | `PROTO_NC_WT_LICENSE_SET_ACK::echo_data` | `4 bytes` |
| 638 | `PROTO_NC_INSTANCE_DUNGEON_RANK_SAVE_DB_CMD` | `19 bytes` | 1285 | `PROTO_NC_WT_LICENSE_SET_REQ` | `54 bytes` |
| 639 | `PROTO_NC_INSTANCE_DUNGEON_RESET_COUNTDOWN_CMD` | `5 bytes` | 1286 | `PROTO_NC_WT_LICENSE_SET_REQ::echo_data` | `4 bytes` |
| 640 | `PROTO_NC_ITEMDB_ADMINCREATEFAIL_ACK` | `6 bytes` | 1287 | `PROTO_NC_WT_MOBINC_CMD` | `2 bytes` |
| 641 | `PROTO_NC_ITEMDB_ADMINCREATESUC_ACK` | `2 bytes` | 1288 | `PROTO_NC_WT_MOB_KILLCOUNT_SET_ACK` | `22 bytes` |
| 642 | `PROTO_NC_ITEMDB_ADMINCREATE_REQ` | `40 bytes` | 1289 | `PROTO_NC_WT_MOB_KILLCOUNT_SET_REQ` | `20 bytes` |
| 643 | `PROTO_NC_ITEMDB_BOOTHTRADE_ACK` | `32 bytes` | 1290 | `PROTO_NC_WT_TITLE_SET_ACK` | `8 bytes` |
| 644 | `PROTO_NC_ITEMDB_BOOTHTRADE_ALL_REQ` | `68 bytes` | 1291 | `PROTO_NC_ZONE_PRISON_END_CMD` | `6 bytes` |
| 645 | `PROTO_NC_ITEMDB_BOOTHTRADE_LOT_REQ` | `103 bytes` | 1292 | `PROTO_NC_ZONE_PRISON_GO_ACK` | `374 bytes` |
| 646 | `PROTO_NC_ITEMDB_BOOTHTRADE_MERGE_REQ` | `74 bytes` | 1293 | `PROTO_NC_ZONE_PRISON_GO_REQ` | `372 bytes` |
| 647 | `PROTO_NC_ITEMDB_BRACELET_UPGRADE_ACK` | `10 bytes` | — | — | — |

### 5.2 Las 9 Plantillas `ShineSendType<PROTO_NC_*>` de `Zone.pdb`

| # | Plantilla `ShineSendType` (`Zone.pdb`) | `sizeof` |
| :---: | :--- | :---: |
| 1 | `ShineSendType<PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK>` | `1 bytes` |
| 2 | `ShineSendType<PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK>` | `1 bytes` |
| 3 | `ShineSendType<PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK>` | `1 bytes` |
| 4 | `ShineSendType<PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK>` | `1 bytes` |
| 5 | `ShineSendType<PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK>` | `1 bytes` |
| 6 | `ShineSendType<PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK>` | `1 bytes` |
| 7 | `ShineSendType<PROTO_NC_MAP_MULTY_LINK_SELECT_ACK>` | `1 bytes` |
| 8 | `ShineSendType<PROTO_NC_QUEST_READ_ACK>` | `1 bytes` |
| 9 | `ShineSendType<PROTO_NC_QUEST_START_ACK>` | `1 bytes` |

### 5.3 Registro Explícito de los 40 Artefactos Internos C++ STL del Flujo TPI de `Zone.pdb`

Los siguientes **40 tipos** en [`docs/Servicios/Zone.md`](./Servicios/Zone.md) son instanciaciones internas de memoria de `MSVC STL` (`std::vector` / `std::allocator` / `std::_Vector_iterator`) para campos dinámicos de 4 paquetes variables (`SkillDamage`, `CardBookmarkInfo`, `CardInform`, `MobOfQuest`). No son paquetes de red independientes; en Go se representan como slices tipados `[]SkillDamage`, `[]CardBookmarkInfo`, `[]CardInform` y `[]MobOfQuest`:

| # | Tipo Interno C++ STL en TPI (`Zone.pdb`) | `sizeof` |
| :---: | :--- | :---: |
| 1 | `std::_Allocator_base<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage>` | `1 bytes` |
| 2 | `std::_Allocator_base<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo>` | `1 bytes` |
| 3 | `std::_Allocator_base<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform>` | `1 bytes` |
| 4 | `std::_Allocator_base<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest>` | `1 bytes` |
| 5 | `std::_Iterator012<std::random_access_iterator_tag,PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,int,PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage const *,PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage const &,std::_Iterator_base0>` | `1 bytes` |
| 6 | `std::_Iterator012<std::random_access_iterator_tag,PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,int,PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo const *,PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo const &,std::_Iterator_base0>` | `1 bytes` |
| 7 | `std::_Iterator012<std::random_access_iterator_tag,PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,int,PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform const *,PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform const &,std::_Iterator_base0>` | `1 bytes` |
| 8 | `std::_Iterator012<std::random_access_iterator_tag,PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,int,PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest const *,PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest const &,std::_Iterator_base0>` | `1 bytes` |
| 9 | `std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> > >` | `4 bytes` |
| 10 | `std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> > >` | `4 bytes` |
| 11 | `std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> > >` | `4 bytes` |
| 12 | `std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> > >` | `4 bytes` |
| 13 | `std::_Vector_iterator<std::_Vector_val<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> > >` | `4 bytes` |
| 14 | `std::_Vector_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> > >` | `4 bytes` |
| 15 | `std::_Vector_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> > >` | `4 bytes` |
| 16 | `std::_Vector_iterator<std::_Vector_val<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> > >` | `4 bytes` |
| 17 | `std::_Vector_val<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> >` | `16 bytes` |
| 18 | `std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> >` | `16 bytes` |
| 19 | `std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> >` | `16 bytes` |
| 20 | `std::_Vector_val<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> >` | `16 bytes` |
| 21 | `std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage>` | `1 bytes` |
| 22 | `std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage>::rebind<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage>` | `1 bytes` |
| 23 | `std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo>` | `1 bytes` |
| 24 | `std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo>::rebind<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo>` | `1 bytes` |
| 25 | `std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform>` | `1 bytes` |
| 26 | `std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform>::rebind<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform>` | `1 bytes` |
| 27 | `std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest>` | `1 bytes` |
| 28 | `std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest>::rebind<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest>` | `1 bytes` |
| 29 | `std::reverse_iterator<std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> > > >` | `0 bytes` |
| 30 | `std::reverse_iterator<std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> > > >` | `0 bytes` |
| 31 | `std::reverse_iterator<std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> > > >` | `0 bytes` |
| 32 | `std::reverse_iterator<std::_Vector_const_iterator<std::_Vector_val<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> > > >` | `0 bytes` |
| 33 | `std::reverse_iterator<std::_Vector_iterator<std::_Vector_val<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> > > >` | `0 bytes` |
| 34 | `std::reverse_iterator<std::_Vector_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> > > >` | `0 bytes` |
| 35 | `std::reverse_iterator<std::_Vector_iterator<std::_Vector_val<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> > > >` | `0 bytes` |
| 36 | `std::reverse_iterator<std::_Vector_iterator<std::_Vector_val<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> > > >` | `0 bytes` |
| 37 | `std::vector<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage,std::allocator<PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage> >` | `16 bytes` |
| 38 | `std::vector<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD::CardBookmarkInfo> >` | `16 bytes` |
| 39 | `std::vector<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform,std::allocator<PROTO_NC_CHAR_CARDCOLLECT_CMD::CardInform> >` | `16 bytes` |
| 40 | `std::vector<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest,std::allocator<PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD::MobOfQuest> >` | `16 bytes` |

---

## 6. Especificación Técnica de los Motores de Gameplay de `ZoneService`

### 6.1 Game Loop Determinístico, Grilla Espacial AOI, Colisión `.shbd`/`.sbi`, Movimiento y Monturas (`Fase 11`)
1. **Tick del Mapa (`ZoneMap.Tick`):** Cada mapa activo ejecuta un bucle determinístico (`100 ms` para movimiento/combate/casteos/IA, `1 s` para `AbState`/regeneración/temporizadores Lua/expiración de propiedad y desaparición de drops en suelo, `30 s` para latido de nodo y autoguardado periódico equivalente a `gds_NC_CHARSAVE_ALL_REQ`).
2. **Grilla Espacial AOI (Sector Grid):** El mapa se divide en celdas cuadradas de visión; cada entidad (`ZonePlayer`, `ZoneMob`, `ZoneNPC`, `ZoneDroppedItem`, `ZoneCamp`, `ZonePet`, `ZoneMagicField`, `ZoneBooth`) se registra en su celda actual. Al cruzar de celda:
   - Se envían al jugador los paquetes `PROTO_NC_BRIEFINFO_*_CMD` de las entidades en las celdas nuevas.
   - Se envía `PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD` (`2 bytes`: `handle`) para las entidades que salen del radio de visión.
   - Se respeta `HIDE_EXCEPT_ME_ON_CMD` / `OFF_CMD` y la matriz de sigilo `StealthMatrix.json` (`Detect` vs `Hide`).
3. **Validación de Movimiento y Puertas:** Todo paquete `WALK_REQ`, `RUN_REQ`, `MOVEWALK_CMD`, `MOVERUN_CMD`, `STOP_REQ`, `JUMP_CMD` y `MOVER_RIDE_ON/OFF` valida transitabilidad en `pkg/data/blockinfo` (`.shbd` + estado vivo de puertas `.sbi`) y cota de velocidad según `ClassParam` + equipamiento + `AbState` + `MoverMain.shn`.
4. **Stats de las 27 Clases, Descanso (`Tent`) y Recolección (`Gather`):** Se mapean las 27 clases reales de `ClassName.shn`, se calculan `MaxHP`/`MaxSP`/`MaxLP` y bonos de stats desde `ClassParam` + `FreeStatTable.json` + `Common.json` + `ExpTable.json`, se gestionan tiendas de descanso (`PITCHTENT_REQ` / `FOLDTENT_REQ` -> `PROTO_NC_BRIEFINFO_CAMP_CMD`) y nodos de recolección (`GATHERSTART_REQ` / `CANCEL` / `COMPLETE` + `CollectInfo.shn`).

### 6.2 Motor de Combate (`NC_BAT`), Habilidades (`NC_SKILL`), `AbState`, `SoulStone` y PvP/PK (`Fase 12`)
1. **Auto-Ataque y Fórmulas Oficiales:** `TARGETTING_REQ`, `UNTARGET_REQ`, `BASHSTART_CMD`, `BASHSTOP_CMD` y `ASSIST_REQ` ejecutan cadencia de golpe según velocidad de arma (`ItemInfo.shn`), línea de visión (`.shbd` `IsAttackBlocked`), precisión/evasión (`TH` vs `TB`), crítico, multiplicador angular posicional (`DamageByAngle.json` — espalda/flanco/frente), penalización/bono por brecha de nivel (`DamageLvGap*.json`), `SwingDamageFactor.shn`, `ArmorMaterialFactor.shn` y consumo de `SoulStone` (`DamageBySoul` / `NC_SOULSTONE` `4` C2S).
2. **Habilidades Activas, de Campo y Pasivas:** `SKILLBASH_OBJ_CAST_REQ`, `SKILLBASH_FLD_CAST_REQ` y `SKILLBASH_CASTABORT_REQ` validan tiempos de casteo, enfriamiento (`CoolTime`), costo de SP/LP, rango, forma de área (`ActionRangeFactor.shn` + `pkg/data/areabmp`), objetivos máximos, consumo de reactivos e invocación de campos mágicos (`PRODUCTFIELD_REQ` / `SkillMechanics.json`: tótems, trampas y escudos). Incluye empoderamiento (`SKILL_EMPOWALLOC_REQ`), olvido (`SKILL_UNLEARN_REQ`) y resurrección por habilidad (`SKILL_REPLYREVIVE_CMD`).
3. **Máquina de Estados `AbState`:** Aplicación, apilamiento, inmunidades/resistencias (`SubAbState.shn` + `AbState.shn` + `data/AbState/*.json`), ticks periódicos DoT/HoT, alteraciones de velocidad/stats, polimorfismo (`CHAR_DEPOLYMORPH_CMD`), borrado voluntario (`BAT_ABSTATE_ERASE_REQ`, `CHAR_CHARGEDBUFF_ERASE_REQ`) y persistencia en `tCharacterState` vía `CALL p_Char_SetStateWithRestTime`.
4. **Muerte, Resurrección y Reglas PvP/PK:** Pérdida/exención de EXP al morir según nivel y buffs pagos, resurrección en punto de guardado (`CHAR_REVIVE_REQ`), resurrección por clérigo/ítem (`CHAR_REBIRTH_CMD`/`REJECT_CMD`, `ITEM_REVIVEITEMUSE_CMD`), arenas libres (`FreeBattleField`), banderas PK (`PKCount`, `PKYellowTime`) y `GameLog` Tipos `2`, `8` y `E`.

### 6.3 IA de Monstruos (`MobBehavior`), Spawns (`MobRegen`), Drops Oficiales, Títulos y Cartas (`Fase 13`)
1. **Ciclo de Vida de Monstruos:** Spawns iniciales y colas de reaparición por mapa/grupo desde `data/MobRegen/*.json` + `MobInfo.shn` + `MobInfoServer.shn` + `MobSpecies.shn`.
2. **IA y Patrullaje:** Rutas `data/MobRoam/*.json`, tabla de amenaza (`AggroList`), persecución, leash de retorno, llamadas de auxilio de familia, diálogos `MobChat.json`, secuencias de golpe `data/MobAttackSequence/*.json`, transiciones `data/MobBehaviorDescript.json` y hooks Lua `data/LuaScript/AIScript/`.
3. **Recompensas, Drops y `AutoPick`:** Cálculo de EXP individual/Party/Raid (`ExpTable.json`, `RestExp`), generación de ítems y Cen desde `ItemDropTable.json` + `ItemDropGroup.json` + `WorldDrop.shn` con opciones aleatorias (`RandomOptionTable.json`), reserva temporal de dueño, recolección automática (`CHAR_CLIENT_AUTO_PICK_REQ` -> `CALL p_AutoPick_*`), anuncios de jefe (`MobKillAnnounce.shn`) y `GameLog` Tipo `6`.
4. **Títulos y Colección de Cartas:** Actualización automática de progreso de títulos (`CharacterTitleData.shn` -> `CALL p_Char_SetTitle*`), cambio de título activo (`CT_SET_CURRENT_REQ`), registro/favoritos/recompensas de cartas (`COLLECT_*` `4` C2S -> `CollectCard*.shn` + `CALL p_Collection_*`).

### 6.4 Economía In-Game, `48` C2S `NC_ITEM`, `SafeTrade`, `Booth`, `Upgrade`, `Enchant` y Crafting (`Fase 14`)
1. **Inventario, Almacenes y Mercaderes NPC:** Cobertura íntegra de los `48` OpCodes C2S de `NC_ITEM` y los `62` métodos `iac_*`: compra (`ITEM_BUY_REQ`), venta (`SELL_REQ`), recompra (`REPURCHASE_REQ`) validando `NPCItemList/*.json` e `ItemShop.shn`; almacén personal (`DEPOSIT`/`WITHDRAW`/`OPENSTORAGEPAGE`), almacén de cuenta (`ACCOUNT_STORAGE_OPEN/CLOSE`), almacén de gremio y academia (`GUILD_STORAGE_WITHDRAW`, `GUILD_ACADEMY_REWARD_STORAGE_*`).
2. **Comercio Seguro (`NC_TRADE`) y Kioscos (`NC_BOOTH`):** Máquina de estados transaccional de 10 pasos para `SafeTrade` (bloqueo de tablero, validación de Cen/slots y transferencia atómica en MySQL) y sistema completo de tiendas personales (`BOOTH_OPEN`, `CLOSE`, `ENTRY`, `REFRESH`, `ITEMTRADE`, `INTERIORSTART`, `SEARCH_ITEM_LIST_CATEGORIZED`, `SEARCH_BOOTH_POSITION` contra `BoothSearchItemList.shn`).
3. **Forja, Gemas, Opciones, Producción y Cofres:** Mejora de armas/armaduras/accesorios/brazaletes/monturas (`UpgradeInfo.shn`, `AccUpGradeInfo.shn`, `BraceletUpgrade.shn`, protección contra rotura/degradación), engaste y perforación de sockets (`ItemSocket.shn`, `ENCHANT_*`), reroll de opciones (`RANDOMOPTION_*`), desarme (`ItemDismantle.shn`), mezcla (`ItemMix.shn`), producción (`Produce.shn`, `RawInfo.shn`), cofres (`MYSTERYVAULT` + `TreasureReward.json`) y `GameLog` Tipos `3, 4, 5, 6, 7, B, C, D`.

### 6.5 Misiones (`QUEST`), Diálogos NPC, Runtime `LuaScript`, Instancias (`ID`/`KQ`/`GT`/`Wedding`), `MiniHouse`, `Pet` y Casino (`Fase 15`)
1. **Quests, NPCs y Escenarios:** Los `7` C2S de `NC_QUEST` (`QuestData.shn`, `QuestParser.json`, `QuestDialog.shn`, `QuestRewardOrder.shn`), menús dinámicos `MENU_SERVERMENU_ACK`, áreas de escenario (`SCENARIO_AREAENTRY_ACK`/`LEAVE_ACK`), `NPCAction.json` y pasos de tutorial (`CHAR_TUTORIAL_STEP_ACK`).
2. **Instancias Gobernadas por Lua:** Instanciación aislada de mapas para Mazmorras (`InstanceDungeon.json` + `NC_INSTANCE_DUNGEON` `2` C2S + `data/LuaScript/ID/`), `Kingdom Quests` en mapa (`KQ_ENTRYRESPONCE_ACK` + `data/LuaScript/KQ/` + `KQReward`), `Guild Tournament` (`6` C2S de captura de bandera, skills, dados y observadores en `GuildTournament.json`) y ceremonias de `Wedding` (`8` C2S en `NC_ACT` + `WeddingData`) más `HolyPromise` (`2` C2S).
3. **`MiniHouse`, `Pet`, `WeaponTitle` y Casino (`Dice`/`Gamble`):** Los `23` C2S de `NC_MINIHOUSE` (`MiniHouse*.shn` + `tMinihouse*`), el ciclo completo de mascotas evolucionables (`NC_PET` + `Pet*.shn` + `tPet`/`tPetNew`/`tPetSkill`), licencias de `WeaponTitle` (`WT_LICENSE_REQ` + `WeaponTitleData.shn`), dados `Taisai` (`NC_DICE` `9` C2S + `DiceGame.shn`) y `GambleHouse` / `SlotMachine` (`NC_GAMBLE` `18` C2S + `SlotMachine.shn`), con `GameLog` Tipos `9, F, G, H, I, J, K`.

### 6.6 Coordinación S2S `Zone <-> World` (`138` `wms_NC_*`), Comandos GM (`AdminLvSet.shn`) y Cierre E2E (`Fase 16`)
1. **Coordinación `wms_NC_*`:** Sincronización en tiempo real entre nodos `Zone00`..`Zone04` y `WorldService` para miembros de grupo (`PARTY_MEMBERINFOREQ_CMD`), puntos de amistad (`FRIEND_UES_FRIEND_POINT_REQ`), gremios (`GUILD_STORAGEOPEN/CLOSE`, `RENAME`, `RETYPE`), guerras de gremio, avisos globales, expulsiones, traslados a prisión y cambio de nombre (`USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ`).
2. **Sistema Oficial de Comandos GM:** Validación estricta del nivel de administrador (`tCharacter.nAdminLevel`) contra la tabla oficial `AdminLvSet.shn`, registrando cada recompensa o sanción en `tGMItemReward`, `tGMItemRewardLog` y `GameLog` Tipo `L` (`CALL p_GameLog_Add_TypeL`).
3. **Diagnóstico, Cuenta y Cierre `245/245`:** `PINGTEST_*`, `SERVERPARAMETER_REQ`, `NORMALLOGOUT_CMD`, `REGISENUMBER_REQ`, `USE_BEAUTY_SHOP_CMD`, `CHARSAVE_UI_STATE_SAVE_REQ` y auditoría final del clúster completo.

---

## 7. Registro Explícito de las Únicas Exclusiones Autorizadas

Conforme a [`AGENTS.md`](../AGENTS.md), **únicamente** quedan fuera del runtime de `XelionCore` los siguientes 3 elementos documentados:
1. **`GamigoZR`:** Mock HTTP local externo creado en 2015 únicamente para esquivar un chequeo DRM de Gamigo en el binario parcheado `Zone.exe`; no forma parte del juego ni del protocolo del cliente.
2. **`OpTool` (`OperatorTool.exe`):** La estructura S2S `PROTO_NC_OPTOOL_CONNECT_BRIF_ACK` (`0x2808`, Dept `10`) en `Zone.md` y las `3` tablas `tMenuAuth`, `tOpLog` y `tOperator` de `OperatorTool.bak` (cuyo DDL ya reside en `database/schema.sql` por completitud documental).
3. **Los `40` tipos internos C++ STL del volcado TPI de `Zone.pdb`:** Listados exhaustivamente en la Sección 5.3 (`std::vector`, `std::allocator`, `std::_Vector_iterator`), que son detalles de implementación de memoria de MSVC y no estructuras independientes de protocolo.

---

## 8. Subdivisión Estratégica de Ejecución por Fases (`FASE_08` a `FASE_16`)

La Etapa 3 se divide estratégicamente en **9 fases modulares, autocontenidas y seguras para la cuota de tokens**, arrancando en **`FASE_08`** inmediatamente a continuación de la Etapa 2 (`FASE_01`..`FASE_07`) y garantizando compilación y tests en verde al cierre de cada fase:

| Fase | Archivo en `docs/plan/` | Alcance Técnico Exacto | Nivel Modelo | Controles que Activa |
| :---: | :--- | :--- | :---: | :--- |
| **`08`** | [`FASE_08_ESQUEMA_PROCEDURES_Y_REPOSITORIO_GAMELOG.md`](./plan/FASE_08_ESQUEMA_PROCEDURES_Y_REPOSITORIO_GAMELOG.md) | Eliminación de `database/zone_tables.sql`, verificación `176/176` tablas en `xelion_online`, transpilación pool-safe de los **`27` SPs** en `database/procedures_gamelog.sql`, `GameLogRepository` (`100% CALL <SP>`) y worker asíncrono `GameLogDispatcher` | **`Medium` / `High`** | **Controles 1, 2 y 3 (Etapa 3)** |
| **`09`** | [`FASE_09_PROTOCOLO_OPCODES_Y_STRUCTS_ZONE_GAMELOG.md`](./plan/FASE_09_PROTOCOLO_OPCODES_Y_STRUCTS_ZONE_GAMELOG.md) | Registro 100% de los **`245` OpCodes C2S** + S2S de `Zone.md`, las **`44` structs `PROTO_NC_LOG_*`** de `GameLog.md` y las **`1.293` structs `PROTO_NC_*`** de `Zone.md` en `pkg/protocol/` con paridad `sizeof` al byte y cero alias privados | **`High`** | **Control 4 (Etapa 3)** |
| **`10`** | [`FASE_10_UNIVERSO_DATOS_SHN_BLOCKINFO_AREABMP_Y_LUA_API.md`](./plan/FASE_10_UNIVERSO_DATOS_SHN_BLOCKINFO_AREABMP_Y_LUA_API.md) | Carga en `DataStore` de las **`180` tablas `.shn` reales en disco** (`194` rutas) + `49` checksums, creación de **`pkg/data/blockinfo`** (`301` archivos `.shbd`/`.sbi`/`.aid`/`.bdt`/`.shab`) y **`pkg/data/areabmp`** (`44` `.bmp`), e implementación de las **`274` funciones C-API `lua_State*`** en `pkg/script/lua/` | **`High`** | **Control 6 (Etapa 3)** |
| **`11`** | [`FASE_11_MOTOR_ESPACIAL_AOI_COLISIONES_BRIEFINFO_Y_MOVIMIENTO.md`](./plan/FASE_11_MOTOR_ESPACIAL_AOI_COLISIONES_BRIEFINFO_Y_MOVIMIENTO.md) | Game Loop por ticks, grilla espacial AOI, `32` paquetes `BRIEFINFO`, colisión real `.shbd`/`.sbi`, las **27 clases** en `ClassParam`, movimiento, salto, emotes, tiendas `Tent`, recolección `Gather` y monturas `MOVER` (`Dept 51`) | **`High`** | Tests AOI, Colisión y Movimiento |
| **`12`** | [`FASE_12_MOTOR_COMBATE_SKILLS_ABSTATE_SOULSTONE_Y_PVP.md`](./plan/FASE_12_MOTOR_COMBATE_SKILLS_ABSTATE_SOULSTONE_Y_PVP.md) | Motor de combate físico/mágico (`Dept 9`), `ActiveSkill`/`PassiveSkill`/campos mágicos (`Dept 18`), máquina `AbState`, `SoulStone` (`Dept 20`), distribución de stats, muerte/revivir, PvP/PK y `GameLog` Tipos `2, 8, E` | **`High`** | Tests Fórmulas Combate, Skills y `AbState` |
| **`13`** | [`FASE_13_MOTOR_MONSTRUOS_IA_MOBBEHAVIOR_DROPS_Y_COLECCION.md`](./plan/FASE_13_MOTOR_MONSTRUOS_IA_MOBBEHAVIOR_DROPS_Y_COLECCION.md) | Spawns `MobRegen`, patrullas `MobRoam`, IA `MobBehavior` + `AIScript` Lua, `MobAttackSequence`, reparto de EXP Party/Raid, tablas de `Drops`, `AutoPick`, Títulos (`Dept 24`), Cartas (`Dept 49`) y `GameLog` Tipo `6` | **`High`** | Tests IA Monstruos, Drops y Colección |
| **`14`** | [`FASE_14_ECONOMIA_INVENTARIO_COMERCIO_BOOTH_UPGRADE_Y_CRAFTING.md`](./plan/FASE_14_ECONOMIA_INVENTARIO_COMERCIO_BOOTH_UPGRADE_Y_CRAFTING.md) | **`48/48` C2S `NC_ITEM`** + **`62` métodos `iac_*`**, Mercaderes NPC, almacenes, `SafeTrade` (`10` C2S), tiendas `BOOTH` (`8` C2S), `Upgrade`, `Enchant`, `RandomOption`, `Dismantle`, `Mix`, `Produce`, `MysteryVault` y `GameLog` Tipos `3, 4, 5, 7, B, C, D` | **`High`** | Tests Transaccionales Economía y Forja |
| **`15`** | [`FASE_15_QUESTS_NPCS_LUA_INSTANCIAS_KQ_MINIHOUSE_PETS_Y_CASINO.md`](./plan/FASE_15_QUESTS_NPCS_LUA_INSTANCIAS_KQ_MINIHOUSE_PETS_Y_CASINO.md) | Motor `QUEST` (`7` C2S), menús NPC y escenarios, instancias Lua (`InstanceDungeon`, `KQ` en mapa, `GuildTournament` en mapa, `Wedding`), `HolyPromise`, `MiniHouse` (`23` C2S), `Pet` (`3` C2S + ciclo completo), `WeaponTitle` y Casino (`DICE` `9` C2S + `GAMBLE` `18` C2S) + `GameLog` Tipos `9, F, G, H, I, J, K` | **`High`** | Tests Quests, Lua KQ/ID, MiniHouse, Pets y Casino |
| **`16`** | [`FASE_16_COORDINACION_S2S_ADMIN_E2E_CLUSTER_Y_AUDITORIA_FINAL.md`](./plan/FASE_16_COORDINACION_S2S_ADMIN_E2E_CLUSTER_Y_AUDITORIA_FINAL.md) | Coordinación S2S `Zone <-> World` (**`138` flujos `wms_NC_*`**), topología multi-nodo `Zone00`..`Zone04`, comandos GM (`AdminLvSet.shn` + `GameLog` Tipo `L`), cierre **`245/245` C2S** y suite forense `scripts/verify_zone_gamelog_100.sh` (`8/8 PASS`) | **`High`** | **Controles 5, 7 y 8 (`8/8 PASS`, `exit 0`)** |

---

## 9. Los 8 Controles Obligatorios de Auditoría Automatizada (`scripts/verify_zone_gamelog_100.sh`)

La Etapa 3 (y con ella la emulación íntegra del clúster oficial `XelionCore`) **solo se considera 100% terminada** cuando el script ejecutable `scripts/verify_zone_gamelog_100.sh` ejecuta y aprueba en verde (`[PASS]`, `exit 0`) los siguientes **8 controles automatizados**, manteniendo simultáneamente en `exit 0` las suites de Etapa 1 (`verify_login_account_100.sh`) y Etapa 2 (`verify_world_character_100.sh`):

1. **Control 1 — Paridad 100% de Esquema `176/176` Tablas Oficiales y Retiro de DDL Fantasma:**
   - Verifica contra `information_schema.TABLES` en `xelion_online` la presencia exacta de `tGameLog` (con sus 15 columnas y tipos exactos) y `TEMP_QuestComplete`, completando las **`176/176` tablas oficiales** de [`docs/TablasLegacy.md`](./TablasLegacy.md), y comprueba que `database/zone_tables.sql` ya no existe en el repositorio.
2. **Control 2 — Paridad `27/27` Stored Procedures de `World00_GameLog` Pool-Safe:**
   - Verifica en `information_schema.ROUTINES` la existencia y compilación válida de los **`27` Stored Procedures** de `database/procedures_gamelog.sql` (`p_GameLog_Add`, `Type1..9`, `TypeA..L`, `usp_Login_setChannel`, `p_GameLog_Get*`) y prueba su invocación concurrente desde un pool `*sql.DB` sin `sql.ErrNoRows` ni contaminación de variables de sesión.
3. **Control 3 — Escaneo AST Go: Cero SQL Inline en `GameLogRepository` y `ZoneService`:**
   - Analiza con `go/ast` todos los archivos de `pkg/database/mysql_gamelog*.go` y `pkg/server/zone/` garantizando **cero sentencias `SELECT`/`INSERT`/`UPDATE`/`DELETE` embebidas** (100% `CALL <SP>`).
4. **Control 4 — Paridad Binaria 100% de OpCodes y `sizeof` (`1.293` Structs `Zone.md` + `44` Structs `GameLog.md`):**
   - Parsea directamente [`docs/Servicios/Zone.md`](./Servicios/Zone.md) y [`docs/Servicios/GameLog.md`](./Servicios/GameLog.md) y verifica que los **`245` OpCodes C2S**, las **`1.293` estructuras `PROTO_NC_*`** de `Zone.md` y las **`44` estructuras `PROTO_NC_LOG_*`** de `GameLog.md` existen con su nombre oficial exacto en `pkg/protocol/` y producen exactamente el `sizeof` en bytes del PDB TPI.
5. **Control 5 — Cobertura `245/245` del Dispatcher C2S de `Zone.exe` y Coordinación S2S (`138` `wms_NC_*` + `196` `gds_NC_*`):**
   - Verifica por AST y tabla de enrutamiento que `pkg/server/zone/` despacha el **100% de las `245` entradas C2S** de `protocolstore(ShinePlayer)` (`0x004D3510`) sin caer en `Unhandled OpCode`, erradicando todo alias privado y cubriendo los flujos S2S.
6. **Control 6 — Universo Completo de Datos en RAM (`180` SHNs Reales / `194` Rutas + `49` Checksums + `301` `BlockInfo` + `44` `AreaBMP` + `274` Lua C-APIs):**
   - Instancia `zone.LoadAllData("data")`, `blockinfo`, `areabmp` y `lua.NewEngine`, verificando carga sin error de los **`180` basenames `.shn` reales en disco** (`194` archivos), los **`49` checksums `.shn`**, los **`158` mapas `.shbd`**, los **`33` `.sbi`**, los **`23` `.aid`**, los **`44` `.bmp`** y el registro real de las **`274` funciones C-API `lua_State*`** de `Zone.pdb` (cero stubs mudos no enlazados).
7. **Control 7 — Prueba E2E Forense Integral del Clúster Completo (`Login -> World -> Zone -> GameLog -> MySQL`):**
   - Levanta `LoginService`, `WorldService`, `ZoneService` y `GameLogDispatcher` contra `xelion_online` y ejecuta un viaje real de cliente: handshake y login -> creación/selección de personaje -> entrada a mapa validando 49 checksums -> recepción del In-Game Chunk -> movimiento validado contra `.shbd` (aceptando celda libre y rechazando celda bloqueada) -> combate y muerte de monstruo con drop oficial + `AutoPick` -> compra NPC, mejora `Upgrade` e intercambio `SafeTrade` -> aceptación y entrega de misión `Quest` -> ejecución de script Lua de instancia/puerta `.sbi` -> desconexión limpia -> verificación directa en MySQL de los cambios persistidos en `tCharacter`, `tItem`, `tQuest` y las filas exactas insertadas en **`tGameLog`** (`Tipos 1, 2, 3, 4, 6, 7, 8, 9`).
8. **Control 8 — Escaneo AST de Cero Datos Hardcodeados en `pkg/server/zone/` y `pkg/script/lua/`:**
   - Verifica por `go/ast` la ausencia total de fallbacks `"RouN"`, constantes mágicas de stats/experiencia (`* 3`, `280`, `JobToClassName` truncado a 6 clases), privilegios GM hardcodeados fuera de `AdminLvSet.shn` y `stubNoOp` mudos en `pkg/script/lua/`.
