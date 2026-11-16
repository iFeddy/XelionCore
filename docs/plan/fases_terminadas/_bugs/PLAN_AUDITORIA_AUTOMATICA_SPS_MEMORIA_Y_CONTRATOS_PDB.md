# Plan de Auditoría y Corrección Automática de SPs MySQL, Sincronización de Memoria y Contratos C2S/S2C

## 1. Objetivo
Erradicar de forma **100% automatizada** (sin depender de pruebas manuales en el cliente `Fiesta.bin`) las tres familias mecánicas de bugs detectadas en pruebas en vivo:

1. **Punto 1 — Stored Procedures T-SQL $\rightarrow$ MySQL 8.4 (`ROW_COUNT()`, `ON DUPLICATE KEY UPDATE` y anchos `VARBINARY`):**
   - En SQL Server (`T-SQL`), un `UPDATE` sobre una fila existente devuelve `@@ROWCOUNT = 1` incluso si los valores nuevos son idénticos a los almacenados.
   - En MySQL 8.4, `ROW_COUNT()` devuelve **`0`** cuando la fila existe pero ningún byte cambió. Esto provoca que patrones traducidos como `UPDATE ...; IF ROW_COUNT() = 0 THEN INSERT ...` arrojen `Error 1062 (23000): Duplicate entry`, y que patrones `UPDATE ...; IF ROW_COUNT() <> 1 THEN ROLLBACK;` aborten transacciones válidas.
   - Adicionalmente, columnas binarias o de texto cuyas longitudes crecieron en el cliente NA2016 arrojan `Error 1406 (22001): Data too long for column`.

2. **Punto 2 — Eliminación de Desincronización de Estado en Memoria (`Character.Sync()` y campos duplicados en `database.Character`):**
   - `database.Character` mantiene campos planos (`ID`, `Name`, `Job`, `Level`, `Male`, `Hair`, `HairColor`, `Face`, `Map`, `XPos`, `YPos`, `CurHP`, `CurSP`) en paralelo con los sub-structs relacionales (`TCharacter`, `TCharacterShape`, `TCharacterState`).
   - En `Character.Sync()` (`pkg/database/character.go`), la rama `else if c.TCharacter.CharNo != 0` sobrescribe los campos planos con los valores de `TCharacter` / `TCharacterShape` / `TCharacterState`. Si un handler modifica `player.Char.XPos`, `player.Char.Hair` o `player.Char.Level` sin actualizar al mismo tiempo el sub-struct interno, cualquier llamada posterior a `UpdateCharacter` revierte el cambio y persiste datos viejos en MySQL.

3. **Punto 3 — Verificador Automático de Contratos Binarios C2S/S2C (`sizeof` del PDB y `.exe` vs. `Read*` / `Write*` en Go):**
   - Cruzar automáticamente las especificaciones TPI de los PDBs oficiales (`docs/Servicios/*.md`) y los binarios (`WorldManager.exe`, `Character.exe`, `Zone.exe`) contra las lecturas `pkt.Reader.Read*` y escrituras `pw.Write*` de cada handler en `pkg/server/login`, `pkg/server/world` y `pkg/server/zone`.
   - Detectar y corregir cualquier handler C2S que lea bytes no declarados en el struct PDB (p. ej. structs vacíos de `1 byte` en C++ sin campos), o cualquier paquete S2C cuyos tipos/anchos de campo (`uint8` vs `uint16` vs `uint32`) o preservación de estado de sesión (`CharacterID` / `StyleCharID` en transiciones `Zone <-> World`) difieran del binario oficial.

---

## 2. Alcance Detallado por Punto

### Punto 1: Barrido y Saneamiento Automático de los 459 Stored Procedures en MySQL 8.4

#### 1.1. Archivos afectados
- [database/procedures_character.sql](../../database/procedures_character.sql) (425 SPs)
- [database/procedures_account.sql](../../database/procedures_account.sql) (7 SPs)
- [database/procedures_gamelog.sql](../../database/procedures_gamelog.sql) (27 SPs)
- [database/schema.sql](../../database/schema.sql) (176 tablas)

#### 1.2. Reglas de transformación automática en SPs
1. **Patrón Upsert (`UPDATE` seguido de `IF ROW_COUNT() = 0 THEN INSERT`):**
   - Auditar todos los SPs que hacen `UPDATE <tabla> SET ... WHERE <pk>; IF ROW_COUNT() = 0 THEN INSERT INTO <tabla> ...`.
   - Reemplazarlos por `INSERT INTO <tabla>(...) VALUES (...) ON DUPLICATE KEY UPDATE ...` o por verificación explícita `IF EXISTS (SELECT 1 FROM <tabla> WHERE <pk>) THEN UPDATE ... ELSE INSERT ... END IF;` para que una actualización con valores idénticos **jamás** intente hacer `INSERT` de una clave primaria existente (`Error 1062`).
2. **Patrón Validación de Existencia en Transacciones (`UPDATE` seguido de `IF ROW_COUNT() <> 1 THEN ROLLBACK`):**
   - Auditar todos los SPs transaccionales donde se actualiza una fila (`tCharacter`, `tItem`, `tItemOptions`, `tGuild`, `tGuildToken`, `tGambleCoin`, `tGambleJackpot`, etc.) y luego se aborta si `ROW_COUNT() <> 1` o `ROW_COUNT() = 0`.
   - Asegurar que cuando la fila existe pero el valor asignado es igual al actual (por ejemplo, actualizar stats, dinero con delta 0, opciones o equipamiento al mismo valor), la condición valide la existencia de la fila (`IF NOT EXISTS (...)`) en lugar de abortar por `ROW_COUNT() = 0`.
3. **Auditoría de Tamaños de Columnas y Parámetros (`VARBINARY` / `VARCHAR`) contra PDBs NA2016:**
   - Comparar todos los parámetros `IN in_s* VARBINARY(N)` / `VARCHAR(N)` de los SPs y las columnas de `database/schema.sql` contra los tamaños máximos de los structs del PDB de `WorldManager.md`, `Character.md` y `Zone.md`.
   - Aplicar los ajustes tanto en `database/schema.sql` y `database/procedures_*.sql` como en vivo sobre la base de datos `xelion_online`.
4. **Test de Idempotencia en MySQL Real:**
   - Añadir un test en `pkg/database/mysql_test.go` que invoque **dos veces consecutivas con exactamente los mismos argumentos** cada método de persistencia/upsert del repositorio sobre MySQL real, verificando que la segunda llamada nunca devuelva `Error 1062`, `Error 1406` ni falle por `ROW_COUNT() = 0`.

---

### Punto 2: Saneamiento de `Character.Sync()` y Unificación de Estado en Memoria (`database.Character`)

#### 2.1. Archivos afectados
- [pkg/database/character.go](../../pkg/database/character.go)
- [pkg/database/mysql_character.go](../../pkg/database/mysql_character.go)
- [pkg/server/zone/](../../pkg/server/zone/) (todos los lugares que mutan `player.Char.*` e invocan `UpdateCharacter` o `SaveChar*`)
- [pkg/server/world/](../../pkg/server/world/) (todos los lugares que mutan `ch.*` e invocan `UpdateCharacter` o `SetCharacterShape`)

#### 2.2. Reglas de corrección
1. **Inversión de prioridad en `Character.Sync()` (`pkg/database/character.go`):**
   - Separar claramente:
     - `PopulateFromCanonical()` (usado **únicamente** al cargar una fila recién leída desde la base de datos o al inicializar un struct donde los campos planos están vacíos).
     - `SyncToCanonical()` / `Sync()` (usado antes de persistir con `UpdateCharacter()`): los campos planos de gameplay (`ID`, `AccountID`, `Name`, `Slot`, `Job`, `Level`, `Male`, `Hair`, `HairColor`, `Face`, `Map`, `XPos`, `YPos`, `CurHP`, `CurSP`) que fueron modificados por los handlers de `ZoneService` o `WorldService` deben **escribirse hacia** `TCharacter`, `TCharacterShape` y `TCharacterState`, y **nunca al revés**.
2. **Barrido AST de mutaciones de `Character`:**
   - Escanear en `pkg/server/zone` y `pkg/server/world` todas las escrituras a `Char.XPos`, `Char.YPos`, `Char.Map`, `Char.TownName`, `Char.CurHP`, `Char.CurSP`, `Char.Level`, `Char.Hair`, `Char.HairColor`, `Char.Face`, `Char.TCharacterState.Money`, `Char.TCharacterState.Exp`, `Char.TCharacterState.Fame` y `Char.TCharacterState.HPStones`/`SPStones`.
   - Garantizar que tanto `UpdateCharacter` como los métodos específicos (`SaveCharLocation`, `SaveCharStats`, `SaveCharLevelExp`, `SetCharacterShape`) mantengan sincronizados en memoria el campo plano y el campo de `TCharacter*`.
3. **Test Unitario de No-Regresión de `Sync()`:**
   - Añadir un test en `pkg/database` que cargue un `Character` desde MySQL (`TCharacter.CharNo != 0`), modifique únicamente sus campos planos (`Char.XPos`, `Char.YPos`, `Char.Hair`, `Char.Level`, `Char.CurHP`), llame a `UpdateCharacter(ch)` y vuelva a leerlo desde MySQL con `GetCharacterByID(ch.ID)`, comprobando que todos los cambios se guardaron sin ser pisados por los valores viejos de `TCharacterState`/`TCharacterShape`.

---

### Punto 3: Verificador Automático de Contratos Binarios C2S/S2C (`PDB` + `.exe` vs. Go)

#### 3.1. Archivos afectados
- Nuevo script verificador: `scripts/verify_c2s_s2c_pdb_contracts.py`
- `pkg/protocol/*.go`
- `pkg/server/login/*.go`
- `pkg/server/world/*.go`
- `pkg/server/zone/*.go`

#### 3.2. Controles automáticos del script verificador
1. **Auditoría de Lecturas C2S (`pkt.Reader.Read*`) contra `sizeof` del PDB:**
   - Parsear todas las estructuras `PROTO_NC_*` de `docs/Servicios/Login.md`, `WorldManager.md` y `Zone.md`.
   - Para cada `case protocol.PROTO_NC_*:` en los dispatchers C2S de `LoginService`, `WorldService` y `ZoneService`:
     - Si el struct en el PDB es una **estructura vacía** (`0 bytes` de campos / `1 byte` en C++ sin campos), verificar que el handler en Go **no consuma bytes inexistentes** que puedan pertenecer al siguiente paquete en el stream TCP o interpretar basura/padding (`0xFF`).
     - Si el struct en el PDB tiene campos de tamaño fijo $N$ bytes, verificar que la suma de `ReadByte` (1), `ReadUint16`/`ReadInt16` (2), `ReadUint32`/`ReadInt32` (4), `ReadUint64`/`ReadInt64` (8) y `ReadBytes(k)`/`ReadString(k)` en el camino principal coincida exactamente con los tipos y offsets `[0..N)` del PDB.
2. **Auditoría de Escrituras S2C (`NewPacketWriter` / `New_PROTO_NC_*`) contra el PDB y `.exe`:**
   - Para cada paquete S2C construido en `pkg/protocol`, `pkg/server/login`, `pkg/server/world` y `pkg/server/zone`:
     - Verificar que el tamaño fijo (o cabecera fija + múltiplo del elemento repetido en paquetes dinámicos) coincida byte a byte con la tabla de offsets del PDB oficial y con el desensamblado en `Zone.exe` / `WorldManager.exe` (por ejemplo: campos de longitud de string de 1 byte vs 2 bytes, códigos de resultado `uint16`, tamaños de arrays fijos).
3. **Auditoría de Preservación de Estado de Sesión en Transiciones `World <-> Zone`:**
   - Verificar todos los flujos C2S de `WorldService` que ocurren inmediatamente después de `PROTO_NC_USER_NORMALLOGOUT_CMD` (`0x0C18` con `logoutType == 1`: retorno a selección de personaje o transición a Beauty Shop) o antes de `PROTO_NC_CHAR_LOGIN_REQ` (`0x1001`), asegurando que `client.CharacterID`, `client.StyleCharID` y `client.AccountID` se preserven tal como lo hace `WorldManager.exe` en la sesión del usuario.

---

## 3. Criterios de Aceptación y Verificación
1. `python3 scripts/verify_c2s_s2c_pdb_contracts.py` termina con `0` discrepancias (`exit 0`).
2. `go test -p 1 -count=1 ./pkg/database/... ./pkg/server/login/... ./pkg/server/world/... ./pkg/server/zone/...` pasa al 100% (`exit 0`), incluyendo los nuevos tests de idempotencia de SPs en MySQL real y el test de `Character.Sync()`.
3. La base de datos `xelion_online` queda limpia (conservando únicamente la cuenta `nUserNo = 1` y el personaje `nCharNo = 1`).

---

## Prompt Listo para Ejecutar en Nueva Sesión

```text
Ejecutá al 100% el plan documentado en docs/plan/PLAN_AUDITORIA_AUTOMATICA_SPS_MEMORIA_Y_CONTRATOS_PDB.md cubriendo los Puntos 1, 2 y 3 (Stored Procedures MySQL 8.4 idempotentes, Character.Sync() sin sobrescritura de campos planos y verificador automático de contratos C2S/S2C contra PDB/EXE), verificando con tests sobre MySQL real y dejando xelion_online limpio al finalizar.
```

