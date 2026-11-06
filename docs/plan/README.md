# Hoja de Ruta Modular por Fases: `WorldService` y `Character` (100% MySQL + PDB + `data/`)

Esta carpeta divide la ejecución del plan maestro exhaustivo [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) en **8 fases cortas, autónomas y seguras para tu cuota de tokens** (`01`, `02`, **`02A`**, `03`, `04`, `05`, `06` y `07`), trabajando directamente sobre `main`.

> [!NOTE]
> La **Etapa 1 (`LoginService`, `Account` y `AccountLog`)** fue completada y auditada al 100% (`7/7` + cero hardcode) y se encuentra archivada en [fases_terminadas/](./fases_terminadas/).

---

## Reglas de Seguridad ante Corte de Tokens

1. **Una sesión limpia de chat por cada Fase (o ejecución secuencial con `/goal`):** Cada archivo `FASE_0X_*.md` contiene su prompt exacto y criterios objetivos de cierre.
2. **Compilación garantizada entre fases:** El orden (`01` -> `02` -> `02A` -> `03` -> `04` -> `05` -> `06` -> `07`) asegura que al terminar cada fase `go build ./...` y todos los tests acumulados (Etapa 1 + Etapa 2) sigan pasando en verde.
3. **Un `git commit` al cerrar cada Fase:** Cada fase termina haciendo commit en `main` **únicamente** cuando pasan al 100% sus controles automáticos.
4. **Cero Datos Hardcodeados y Cero Omisiones:** Todo nombre de tabla, Stored Procedure, OpCode y struct C++ proviene 1:1 de los PDBs (`WorldManager.md`, `Character.md`), `.bak` (`TablasLegacy.md`, `World00_Character_all_procedures.sql`), binarios `.exe` y recursos `data/` (`.shn`, `.json`, `.txt`).

---

## Orden de Ejecución y Nivel de Modelo Recomendado

| Orden | Archivo de Fase | Objetivo Principal | Nivel Recomendado | Control de Auditoría que Activa | Estado |
| :---: | :--- | :--- | :---: | :--- | :---: |
| **1** | [FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md](./FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md) | Creación inicial de 91 tablas de `World00_Character` y semillas base en `xelion_online` | **`Medium`** | **Control 1 (base)** | Completado |
| **2** | [FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md](./FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md) | Transpilación y validación de los 311 Stored Procedures en `database/procedures_character.sql` | **`Medium`** *(o `High`)* | **Control 2 (base)** | Completado |
| **2A** | [FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md](./FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md) | **Saneamiento 100% Fases 1 y 2:** `114/114` tablas `.bak` + `2` S2S (`tActiveZoneTransfer`, `tWorldZoneNode`), limpieza de `*_Template` falsos en `tCharacter`, regeneración completa de `default/tutorial_character_data.json` y los **`425` SPs** pool-safe (`proc_logic` + `SELECT` final) | **`High`** | **Controles 1 y 2 al 100%** | Pendiente |
| **3** | [FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md](./FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md) | Migración de los 50 métodos inline + todos los dominios a **`CALL <SP>`** en `pkg/database` y eliminación de `TransferManager` en memoria | **`High`** | **Control 3** (AST Go: cero SQL inline / cero memoria S2S) | Pendiente |
| **4** | [FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md](./FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md) | 100% de OpCodes y structs C++ no-`OpTool` de `WorldManager.md` y `Character.md` en `pkg/protocol/`, erradicación de los 19 alias privados y paridad `sizeof` | **`High`** | **Control 4** (Paridad OpCodes + `sizeof` PDB) | Pendiente |
| **5** | [FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md](./FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md) | Ciclo completo de cuenta/avatares, estilo, opciones `IMPROVE`, tutorial, asistencia y carga de los **`74` `.shn` únicos** (`54 + 71`) sin hardcode | **`High`** | **Control 5A** (Dispatcher Avatar/Char + 74 SHNs) | Pendiente |
| **6** | [FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md](./FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md) | Cobertura **100% de los `167` handlers C2S** de `WorldManager.pdb`: Chat, Party, Raid, Friends, KQ, Guild/War/Tournament/Emblem, Academy, HolyPromise, Prison, Report y Topología `tWorldZoneNode` | **`High`** | **Control 5B** (`167/167` C2S + S2S) | Pendiente |
| **7** | [FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md](./FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md) | Traspaso multi-proceso `World <-> Zone`, carga real de los 16 paquetes del In-Game Chunk, portales inter-zona, logout completo y suite final `verify_world_character_100.sh` | **`High`** | **Controles 6, 7 y 8** (`8/8 PASS`, `exit 0`) | Pendiente |

---

## Archivo Histórico de Etapas Anteriores

- [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./fases_terminadas/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md): Plan maestro de la Etapa 1 completada al 100% (`LoginService`, `Account`, `AccountLog`, `StatisticsData`).
- [fases_terminadas/](./fases_terminadas/): Especificaciones históricas de las Fases 01 a 07 de la Etapa 1.
