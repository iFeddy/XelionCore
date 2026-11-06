# Hoja de Ruta Modular por Fases: `WorldService` y `Character` (100% MySQL)

Esta carpeta divide la ejecución del plan maestro [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) en **7 fases cortas, autónomas y seguras para tu cuota de tokens**, trabajando directamente sobre `main`.

> [!NOTE]
> La **Etapa 1 (`LoginService`, `Account` y `AccountLog`)** fue completada exitosamente al 100% y se encuentra archivada con sus 7 fases en [fases_terminadas/](./fases_terminadas/).

---

## Reglas de Seguridad ante Corte de Tokens

1. **Una sesión limpia de chat por cada Fase:** Abrí un chat nuevo para cada archivo `FASE_0X_*.md`. Así el contexto arranca vacío y consumís una fracción mínima de tokens.
2. **Compilación garantizada entre fases:** El orden (`01` -> `07`) está diseñado para que al terminar cada fase `go build ./...` y los tests acumulados sigan pasando en verde.
3. **Un `git commit` al cerrar cada Fase:** Cada fase termina haciendo commit en `main` **únicamente** cuando pasan sus controles automáticos.
4. **¿Qué hacer si te quedás sin tokens en mitad de una Fase?**
   - Como la fase anterior ya quedó commiteada en `main` y funcionando, cuando retomes tenés dos opciones simples:
     * **Opción 1 (Continuar donde quedó):** Abrir un chat nuevo y pegar el prompt de esa misma `FASE_0X` agregando *"Revisá `git status` y `git diff` y terminá lo que falte de esta fase"*.
     * **Opción 2 (Volver al último punto 100% limpio):** Ejecutar `git reset --hard HEAD && git clean -fd` para volver instantáneamente al cierre exacto de la fase anterior.

---

## Orden de Ejecución y Nivel de Modelo Recomendado

| Orden | Archivo de Fase | Objetivo Principal | Nivel Recomendado | Control de Auditoría que Activa | Estado |
| :---: | :--- | :--- | :---: | :--- | :---: |
| **1** | [FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md](./FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md) | Creación y verificación de las **61 tablas** oficiales de `World00_Character` en MySQL y datos de inicialización | **`Medium`** | **Control 1** (`information_schema.TABLES/COLUMNS`) | Completado |
| **2** | [FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md](./FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md) | Implementación y validación de los **311 Stored Procedures** oficiales de `Character.exe` y `WorldManager.exe` | **`Medium`** *(o `High`)* | **Control 2** (`information_schema.ROUTINES` — 311 SPs) | Completado |
| **2A** | [FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md](./FASE_02A_SANEAMIENTO_DB_SEMILLAS_Y_425_PROCEDURES.md) | Saneamiento de esquema (`114/114`), semillas (`data/`) y los `425` Stored Procedures Pool-Safe | **`High`** | **Controles 1 y 2 al 100%** | Pendiente |
| **3** | [FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md](./FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md) | Reemplazo total de SQL inline por `CALL` en `pkg/database/mysql_character.go` e interfaces tipadas | **`High`** | **Control 3** (Escaneo AST Go: cero SQL inline) | Pendiente |
| **4** | [FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md](./FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md) | Definición estricta de OpCodes y estructuras C++ de los PDBs de `Character` y `WorldManager` | **`Medium`** | **Control 4** (Paridad OpCodes + `sizeof` al byte) | Pendiente |
| **5** | [FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md](./FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md) | Gestión completa de avatares en `WorldService`: lista, creación, borrado y selección 100% MySQL | **`High`** | **Control 5** (Cobertura Dispatcher Character) | Pendiente |
| **6** | [FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md](./FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md) | Coordinación de chat global, gremios, academia, amigos, party y topología S2S con zonas | **`High`** | **Control 6** (Subsistemas y Topología S2S) | Pendiente |
| **7** | [FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md](./FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md) | Integración transversal con `ZoneService`, 54+71 SHNs, prueba E2E forense y suite de 8 controles | **`High`** | **Controles 7 y 8** + `verify_world_character_100.sh` | Pendiente |

---

## Archivo Histórico de Etapas Anteriores

- [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./fases_terminadas/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md): Plan maestro de la Etapa 1 completada al 100% (LoginService, Account, AccountLog, StatisticsData, PDB 1:1).
- [fases_terminadas/](./fases_terminadas/): Registro histórico de las Fases 01 a 07 de la Etapa 1 con sus commits en verde.
