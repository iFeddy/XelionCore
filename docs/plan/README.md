# Hoja de Ruta Modular por Fases (`Login`, `Account` y `AccountLog` — 100% MySQL)

Esta carpeta divide la ejecución del plan maestro [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) en **7 fases cortas, autónomas y seguras para tu cuota de tokens**, trabajando directamente sobre `main`.

## Reglas de Seguridad ante Corte de Tokens

1. **Una sesión limpia de chat por cada Fase:** abrí un chat nuevo para cada archivo `FASE_0X_*.md`. Así el contexto arranca vacío y consumís una fracción mínima de tokens.
2. **Compilación garantizada entre fases:** el orden (`01` -> `07`) está diseñado para que al terminar cada fase `go build ./...` y los tests existentes sigan pasando en verde.
3. **Un `git commit` al cerrar cada Fase:** cada fase termina haciendo commit en `main` **únicamente** cuando pasan sus controles automáticos.
4. **¿Qué hacer si te quedás sin tokens en mitad de una Fase?**
   - Como la fase anterior ya quedó commiteada en `main` y funcionando, cuando retomes (en una semana o con otro modelo) tenés dos opciones simples:
     * **Opción 1 (Continuar donde quedó):** abrir un chat nuevo y pegar el prompt de esa misma `FASE_0X` agregando *"Revisá `git status` y `git diff` y terminá lo que falte de esta fase"*.
     * **Opción 2 (Volver al último punto 100% limpio):** ejecutar `git reset --hard HEAD && git clean -fd` para volver instantáneamente al cierre exacto de la fase anterior.

---

## Orden de Ejecución y Nivel de Modelo Recomendado

| Orden | Archivo de Fase | Objetivo Principal | Nivel Recomendado | Control de Auditoría que Activa | Estado |
| :---: | :--- | :--- | :---: | :--- | :---: |
| **1** | [FASE_01_ESQUEMA_TABLAS_Y_VISTAS.md](./FASE_01_ESQUEMA_TABLAS_Y_VISTAS.md) | Saneamiento 100% de tablas, tipos, FKs, vistas y semillas en `xelion_online` | **`Medium`** | **Control 1** (`information_schema.TABLES/VIEWS/COLUMNS`) | Completado |
| **2** | [FASE_02_PROCEDURES_ACCOUNT.md](./FASE_02_PROCEDURES_ACCOUNT.md) | Creación y testeo de los **75 Stored Procedures** de `Account` (`Account.pdb` + `Account.bak`) | **`Medium`** *(o `High`)* | **Control 2A** (`information_schema.ROUTINES` — `Account`) | Completado |
| **3** | [FASE_03_PROCEDURES_ACCOUNTLOG_STATS_SESSION.md](./FASE_03_PROCEDURES_ACCOUNTLOG_STATS_SESSION.md) | Creación y testeo de los **59 Stored Procedures** de `AccountLog`, `StatisticsData` y `LoginSession` | **`Medium`** | **Control 2B** (100% `information_schema.ROUTINES`) | Completado |
| **4** | [FASE_04_PROTOCOLO_OPCODES_Y_STRUCTS_PDB.md](./FASE_04_PROTOCOLO_OPCODES_Y_STRUCTS_PDB.md) | Corrección de OpCodes (`0x0C02`), las estructuras C++ de los 3 PDBs y limpieza de nombres deprecados | **`Medium`** | **Control 4** (Paridad OpCodes + `sizeof` contra `.md` PDB) | Pendiente |
| **5** | [FASE_05_REPOSITORIOS_GO_100_CALL_MYSQL.md](./FASE_05_REPOSITORIOS_GO_100_CALL_MYSQL.md) | Reemplazo total de SQL inline por `CALL` en `pkg/database`, borrado de `MemoryAccountRepository` y `accounts.json` | **`High`** | **Control 3** (Escaneo AST Go: cero SQL inline / cero JSON) | Pendiente |
| **6** | [FASE_06_LOGINSERVICE_MULTIRREGION_Y_MUNDOS.md](./FASE_06_LOGINSERVICE_MULTIRREGION_Y_MUNDOS.md) | `LoginService` completo: `FOUS`, `XTrap`, todas las regiones (`US/GER/TW/JP/CH/Netmarble/Outspark/OTP/Teenager`), `tWorldServer` | **`High`** | **Control 5** (Cobertura 100% Dispatcher `LoginService`) | Pendiente |
| **7** | [FASE_07_INTEGRACION_WORLD_ZONE_Y_AUDITORIA_FINAL.md](./FASE_07_INTEGRACION_WORLD_ZONE_Y_AUDITORIA_FINAL.md) | Integración `WorldService` + `ZoneService` (`64B` token, `sOTP`, `ReturnUser`, `49 SHNs`, Logout, CCU real) + E2E forense final | **`High`** | **Controles 6 y 7** + `scripts/verify_login_account_100.sh` | Pendiente |

> [!TIP]
> **Ahorro de cuota:** Las Fases **1, 2, 3 y 4** pueden correrse perfectamente en **`Medium`**, reservando **`High`** únicamente para las Fases **5, 6 y 7**.
