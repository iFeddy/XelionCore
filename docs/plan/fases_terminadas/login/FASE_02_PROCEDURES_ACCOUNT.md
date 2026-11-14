# Fase 02 — Los 75 Stored Procedures Nativos de `Account` en MySQL

- **Nivel recomendado:** `Medium` *(o `High`)*
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md](./PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md) (**Secciones 0, 3.1 y 3.2 completas**)
- **Fuente SQL original obligatoria:** [`docs/sql/procedures/Account_all_procedures.sql`](../sql/procedures/Account_all_procedures.sql) y [`docs/Servicios/Account.md`](../Servicios/Account.md)
- **Dependencias previas:** Fase 01 completada.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 02 (`docs/plan/FASE_02_PROCEDURES_ACCOUNT.md`) leyendo las Secciones 0, 3.1 y 3.2 de `docs/plan/PLAN_ACCION_LOGIN_ACCOUNT_ACCOUNTLOG.md` y traduciendo fielmente el código SQL oficial de `docs/sql/procedures/Account_all_procedures.sql` a MySQL sobre `xelion_online`.
Cero omisiones: debés crear y cargar los 75 Stored Procedures no-OpTool de Account (13 ODBC + 62 restantes de Account.bak) respetando la convención de `SELECT` final para compatibilidad con el pool `*sql.DB` de Go. Al pasar todos los tests, actualizá `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Crear `database/procedures_account.sql` (`DROP PROCEDURE IF EXISTS` + `CREATE PROCEDURE`):**
   - **Los 13 SPs invocados por ODBC en `Account.pdb` (Sección 3.1):**
     1. `usp_User_USLoginGame` *(busca `tUser.sUserID`, fallback `tUserOption.sEmail` con `bIsDelete=0`, alta automática `r_` gobernada por `Auto_Account_Creation.Allow_Creation`, invoca lógica `usp_User_chkBlock`, valida `tUserAuth.bIsLoginAble`, retorna `nResult`, `nUserNo`, `nAuthID`)*
     2. `usp_User_loginGame`
     3. `usp_User_LoginGame` *(variante de 5 argumentos; Nota: en MySQL los nombres de rutinas son case-insensitive, por lo que `usp_User_loginGame` y `usp_User_LoginGame` comparten nombre en MySQL -> documentar y exponer `usp_User_loginGame` + `usp_User_LoginGame_5` si difiere aridad, verificando ambos)*
     4. `usp_User_LoginNetmarble`
     5. `usp_Pcbang_chkIP`
     6. `usp_ReturnUser_GetDatediffLastlog`
     7. `usp_ReturnUser_SetUserItem` *(inserta en `tReturnUser` y acredita ítem en `tChargeItem` con `isDraw = 0`)*
     8. `usp_User_GetDatediffCreated`
     9. `usp_Event_addCash` *(con `tCash.status = 0` + `tCashLog`)*
     10. `usp_Item_setSerialKey` *(rechaza `serialKey` ya utilizado sin sobrescribir)*
     11. `p_UserKey_Get_FromAccount`
     12. `p_UserPlayRemainTime_Get_FromAccount`
     13. `p_UserPlayRemainTime_Set_FromAccount`
   - **Los 62 SPs restantes no-OpTool de `Account.bak` (Sección 3.2, totalizando los 68 no-OpTool del `.bak`):**
     * Permisos/Bloqueos (5): `p_UserGetAuthList`, `p_UserSetAuthLoginable`, `usp_User_getAuthList`, `usp_User_setLoginable`, `usp_User_chkBlock` *(solo limpia `tUser.bIsBlock = 0` cuando no queda ningún otro bloqueo activo `bStatus = 1` en `tUserBlock`)*.
     * Cuenta y Opciones (18): `usp_User_chkRegister`, `usp_User_chkUserID`, `usp_User_delete`, `usp_User_getAge`, `usp_User_getLogin`, `usp_User_getOption`, `usp_User_getUser`, `usp_User_insert`, `usp_User_insert_Old`, `usp_User_insertOption`, `usp_User_insertOption_Old`, `usp_User_search`, `usp_User_searchZipcode`, `usp_User_setChannelLog`, `usp_User_setDrawComment`, `usp_User_update`, `usp_User_updatePassword`, `usp_User_USloginGame_O`.
     * Cash (4): `usp_Cash_getCash`, `usp_Cash_setCash`, `usp_Cash_setCash_EVENT_L20`, `usp_Cash_useCash` *(transacción con `FOR UPDATE` y `status = 0`)*.
     * ChargeItem (4): `usp_Charge_ItemDel` *(copia a `tChargeItemDelLog` con el mismo `rowNo` y borra)*, `usp_Charge_ItemDraw`, `usp_Charge_ItemInsert`, `usp_Charge_ItemSelectList`.
     * Item Shop sobre `tAccountItem` y vistas (20): `usp_Item_getBuyCount`, `usp_Item_getBuyList`, `usp_Item_getCategoryName`, `usp_Item_getEventList`, `usp_Item_getGroup`, `usp_Item_getGroupList`, `usp_Item_getGroupPrice`, `usp_Item_getItem`, `usp_Item_getOrder`, `usp_Item_getOrderList`, `usp_Item_getSubCategory`, `usp_Item_getTotalCount`, `usp_Item_getUser`, `usp_Item_setBuyOrder`, `usp_Item_setCategory`, `usp_Item_setItem`, `usp_Item_setOrder`, `usp_Stat_sellitem`, `usp_User_getBuyCount`, `usp_User_getBuyList`.
     * Pagos (8): `usp_Pay_chkLimit`, `usp_Pay_getAgreement`, `usp_Pay_getOrder`, `usp_Pay_setAgreement`, `usp_Pay_setOrder`, `usp_Pay_setResult`, `usp_User_getPayCount`, `usp_User_getPayList`.
     * Tutorial (3): `usp_Tutorial_getCharInfo`, `usp_Tutorial_getIsSkip`, `usp_Tutorial_setCharInfo`.
2. **Cargar `database/procedures_account.sql` en MySQL `xelion_online`** verificando cero errores de sintaxis.
3. **Crear `pkg/audit/fase02_account_sp_test.go`:**
   - **Control 2A:** parsea `docs/ProcedimientosLegacy.md` (sección `Account`) + `docs/Servicios/Account.md` y valida contra `information_schema.ROUTINES` que el 100% de los SPs no-`usp_OPT_*` existen en `xelion_online`.
   - Ejecuta pruebas reales de negocio con `CALL` para login por ID, login por email, prefijo `r_` (`Allow_Creation` `1` vs `0`), expiración selectiva en `usp_User_chkBlock`, `usp_Cash_setCash`/`getCash`/`useCash` (`status = 0`), `usp_Charge_*`, `usp_Item_setBuyOrder`, `usp_ReturnUser_*` y `usp_Tutorial_getIsSkip`.

---

## Criterio de Cierre de la Fase 02

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "TestFase01|TestFase02"` -> `PASS`.
3. Actualizar `docs/plan/README.md` y ejecutar `git commit -m "feat(db): Fase 02 - 100% stored procedures oficiales de Account en MySQL"`.
