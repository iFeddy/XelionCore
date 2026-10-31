# Fase 01 — Esquema Único `xelion_online`: Tablas, Columnas, Vistas y Semillas

- **Nivel recomendado:** `Medium`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 0, 1 y 2 completas**)
- **Dependencias previas:** Ninguna (punto de partida sobre `main`).

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 01 (`docs/plan/FASE_01_ESQUEMA_TABLAS_Y_VISTAS.md`) leyendo como especificación obligatoria las Secciones 0, 1 y 2 completas de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md` y respetando `AGENTS.md`.
No avances a la Fase 02 ni modifiques todavía código de servidores Go. Al terminar con todos los tests en verde, actualizá el estado de la Fase 1 a "Completado" en `docs/plan/README.md` y hacé un `git commit` limpio en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Actualizar [`database/schema.sql`](../../database/schema.sql) en la base única `xelion_online`:**
   - **Resolver las 3 colisiones de `.bak` (Sección 2):**
     * Crear tabla `tAccountItem` con las **22 columnas exactas** de `Account..tItem` y FK `category -> tItemCategory(category)`.
     * Crear tabla `tAccountDemographics` (`userNo INT PK, gender INT, age INT, registerDate DATETIME`) para `AccountLog..tUser`.
     * Eliminar la segunda definición duplicada de `tTutorial` dejando una única tabla `tTutorial` con PK `(nUserNo, nCharNo)` e índice `(nUserNo, nState)`.
   - **Corregir todos los bugs estructurales de columnas/PKs/defaults detectados en las Secciones 2.1, 2.2 y 2.3:**
     * `tCash`: `status TINYINT(1) NOT NULL DEFAULT 0`, PK `(userNo, cashtype)`.
     * `tCashLog`: unificar columna `vaildDate DATETIME DEFAULT NULL`.
     * `tChargeItem`: `orderNo BIGINT NOT NULL DEFAULT 0`, índice `(userNo, isDraw)`.
     * `tChargeItemDelLog`: quitar `AUTO_INCREMENT` de `rowNo`, cambiar `orderNo` a `BIGINT`.
     * `tDrawComment`: ampliar `sComment` a `VARCHAR(1000)`.
     * `tOrder`: permitir `pcbangNo INT DEFAULT NULL`.
     * `tPasswordLog`: ampliar `sPrevPW` y `sCurrPW` a `VARCHAR(32)`.
     * `tPay`: activar `AUTO_INCREMENT` en `orderNo`.
     * `tPayResult`: permitir `pgTID VARCHAR(50) DEFAULT NULL`.
     * `tSerialKey`: PK `(serialKey)` e índice `(userNo)`.
     * `tUser`: alinear las **12 columnas exactas** de `Account.bak` (`nUserNo`, `sUserID`, `sUserPW`, `sUserName`, `bIsBlock`, `bIsDelete`, `nAuthID`, `sUserIP`, `dDate`, `sUserNo`, `blsBlock`, `blsDelete`) — **Nota temporal de compatibilidad con Fase 01..04:** antes de borrar `bIsLoginAble` verificá si `pkg/database/mysql_account.go` aún la consulta; si la consulta, podés dejar `mysql_account.go` leyendo `COALESCE(a.bIsLoginAble, 1)` sin `u.bIsLoginAble` en esta misma fase para eliminar desde ya la columna inventada `tUser.bIsLoginAble` sin romper tests existentes.
     * `tUserOption`: ampliar `sEmail` a `VARCHAR(256)` con índice `(sEmail)` y permitir `NULL` en los campos opcionales.
     * `tbSellitem`: corregir `baseDate VARCHAR(10) NOT NULL` con PK `(baseDate, goodsNo)`.
     * `tAccountLog`: ampliar `sUserID VARCHAR(60)`, `sUserPassword VARCHAR(32)`, `sIP VARCHAR(30)` e índices `(nUserNo, nType, dDate)` y `(dDate, nType)`.
     * `tEventWorldCup`: PK `(nEventCharNo, nLevel, nWorldNo)`.
     * `tEvent_Friend`: agregar `nIndex INT AUTO_INCREMENT PK`, `dRegDate DATETIME`, `dCompleteDate DATETIME NULL`.
     * `tIPBlock`: PK `(sIPStart, sIPEnd)` + columnas generadas/indexadas `nIPStartNum` y `nIPEndNum` con `INET_ATON` para rangos IPv4.
     * `tLogout`: `level SMALLINT UNSIGNED` e índice `(userNo, logDate DESC)`.
     * `tReport`: ampliar `sReporterCharID`, `sHarmer`, `sReportType` a `VARCHAR(20)` y añadir `dDealDate DATETIME NULL`, `sOperID VARCHAR(20) NULL`.
   - **Agregar todas las tablas faltantes listadas en las Secciones 2.1, 2.2 y 2.3:**
     * `Auto_Account_Creation`, `tPcbang`, `tUserKey`, `tUserPlayRemainTime`, `tZipcode`, `tbl_user`, `tbl_user_test`, `tUserTemp`, `tUserTemp2`, `tUser_ori`, `temp_newUser`.
     * `tEvent_L20`, `tHackScanLog`, `tSyslog`.
     * Verificar las 13 tablas de `StatisticsData.bak` + crear `tWorldServer` y `tActiveAccountSession`.
   - **Crear las 3 vistas oficiales:** `VT_ITEM_SELL`, `VT_ITEM_GROUP` y `VT_PAY`.
2. **Crear `database/seed_login_world.sql` (idempotente con `INSERT ... ON DUPLICATE KEY UPDATE`):**
   - `FOUS`: `GameVersionKey = '10022024000000'` (sincronizado con `ClientVersionKeyInfo.txt`).
   - `tUserAuth`: niveles `1` (Player, `bIsLoginAble=1`), `2` (Tester, `1`), `3` (GM, `1`), `9`/`100` (Admin, `1`), `0` (BlockedAuth, `0`).
   - `Auto_Account_Creation`: `Allow_Creation = 1`.
   - `tWorldServer`: mundo `0` / `1` configurado con IP, puerto, `nMaxUser = 1500`, `nUserRatable = 100`, `bIsMaintenance = 0`.
   - Cuenta `admin` inicial en `tUser` + `tUserOption` (sin pisar cuentas existentes).
3. **Aplicar migración no destructiva sobre la instancia MySQL local `xelion_online`** para que tanto una instalación limpia (`schema.sql` + `seed_login_world.sql`) como la base actual queden 100% alineadas.
4. **Crear `pkg/audit/fase01_schema_test.go` (Control 1 de la Sección 9.1):**
   - Verifica contra `information_schema.TABLES`, `VIEWS` y `COLUMNS` de `xelion_online` el 100% de las tablas, vistas y tipos/defaults críticos.

---

## Criterio de Cierre de la Fase 01

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run TestFase01` -> `PASS`.
3. `go test ./pkg/database/...` -> `PASS`.
4. Marcar Fase 1 como `Completado` en `docs/plan/README.md` y ejecutar `git commit -m "feat(db): Fase 01 - esquema 100% oficial Account, AccountLog y StatisticsData en xelion_online"`.
