package audit_test

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// md5Hash helper for test passwords
func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// TestFase02_Control2A_ExhaustiveCatalog valida contra information_schema que
// el 100% de los 75 Stored Procedures no-OpTool de Account existen en xelion_online.
func TestFase02_Control2A_ExhaustiveCatalog(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. Parsear docs/ProcedimientosLegacy.md para obtener los no-OpTool de Account
	repoRoot := "../.."
	legacyDocPath := filepath.Join(repoRoot, "docs", "ProcedimientosLegacy.md")
	content, err := os.ReadFile(legacyDocPath)
	if err != nil {
		t.Fatalf("No se pudo leer %s: %v", legacyDocPath, err)
	}

	docText := string(content)
	accountSecIdx := strings.Index(docText, "## Account (`Account.bak`)")
	if accountSecIdx == -1 {
		t.Fatalf("No se encontro seccion de Account en %s", legacyDocPath)
	}
	nextSecIdx := strings.Index(docText[accountSecIdx:], "\n## World00_Character")
	if nextSecIdx == -1 {
		nextSecIdx = len(docText) - accountSecIdx
	}
	accountSection := docText[accountSecIdx : accountSecIdx+nextSecIdx]

	// Extraer nombres de procedimientos de las tablas Markdown: | `nombre` |
	reProc := regexp.MustCompile("`([a-zA-Z0-9_]+)`")
	matches := reProc.FindAllStringSubmatch(accountSection, -1)

	legacyNonOpt := make(map[string]bool)
	for _, m := range matches {
		name := m[1]
		if strings.HasPrefix(name, "usp_OPT") || name == "Procedimiento" {
			continue
		}
		legacyNonOpt[strings.ToLower(name)] = true
	}

	if len(legacyNonOpt) != 68 {
		t.Fatalf("Se esperaban exactamente 68 procedimientos no-OpTool de Account.bak, pero se encontraron %d", len(legacyNonOpt))
	}

	// 2. Agregar los 7 procedimientos ODBC exclusivos documentados en Account.md y Seccion 3.1
	exclusiveOdbc := []string{
		"usp_user_logingame_5",
		"usp_user_loginnetmarble",
		"usp_pcbang_chkip",
		"usp_event_addcash",
		"p_userkey_get_fromaccount",
		"p_userplayremaintime_get_fromaccount",
		"p_userplayremaintime_set_fromaccount",
	}
	allExpected := make(map[string]bool)
	for k := range legacyNonOpt {
		allExpected[k] = true
	}
	for _, p := range exclusiveOdbc {
		allExpected[p] = true
	}

	if len(allExpected) != 75 {
		t.Fatalf("El universo total de procedimientos esperados debe ser 75, pero es %d", len(allExpected))
	}

	// 3. Consultar information_schema.ROUTINES en MySQL
	rows, err := db.Query(`
		SELECT ROUTINE_NAME 
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = 'xelion_online' AND ROUTINE_TYPE = 'PROCEDURE'
	`)
	if err != nil {
		t.Fatalf("Error consultando information_schema.ROUTINES: %v", err)
	}
	defer rows.Close()

	actualRoutines := make(map[string]bool)
	for rows.Next() {
		var rName string
		if err := rows.Scan(&rName); err != nil {
			t.Fatalf("Error escaneando ROUTINE_NAME: %v", err)
		}
		actualRoutines[strings.ToLower(rName)] = true
	}

	// 4. Validar cobertura del 100%
	var missing []string
	for p := range allExpected {
		if !actualRoutines[p] {
			missing = append(missing, p)
		}
	}

	if len(missing) > 0 {
		t.Fatalf("Control 2A FAILED: Faltan %d Stored Procedures en xelion_online:\n%s",
			len(missing), strings.Join(missing, "\n"))
	}

	t.Logf("Control 2A PASS: Los 75 Stored Procedures oficiales no-OpTool de Account existen en xelion_online")
}

// TestFase02_LoginUS_ByID_And_Email prueba la autenticacion oficial de US
// soportando tanto login por sUserID como fallback por sEmail con tUserOption.
func TestFase02_LoginUS_ByID_And_Email(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUser = "test_fase02_login"
	const testEmail = "test_fase02@xelion.test"
	passPlain := "supersecret123"
	passHash := md5Hash(passPlain)

	// Limpiar datos previos
	_, _ = db.Exec("DELETE FROM tUserOption WHERE sEmail = ?", testEmail)
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", testUser)

	// Crear usuario con email en tUserOption
	res, err := db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, ?, ?, 0, 0, 1, '127.0.0.1', NOW())
	`, testUser, passHash, testUser)
	if err != nil {
		t.Fatalf("Error creando usuario de prueba: %v", err)
	}
	userNo, _ := res.LastInsertId()

	_, err = db.Exec(`
		INSERT INTO tUserOption (nUserNo, sEmail)
		VALUES (?, ?)
	`, userNo, testEmail)
	if err != nil {
		t.Fatalf("Error creando tUserOption: %v", err)
	}

	// 1. Probar login por sUserID directo
	var nResult, outUserNo, outAuthID int
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", testUser, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL usp_User_USloginGame por ID fallo: %v", err)
	}
	if nResult != 0 || outUserNo != int(userNo) || outAuthID != 1 {
		t.Fatalf("Login por ID esperado (0, %d, 1), obtenido (%d, %d, %d)", userNo, nResult, outUserNo, outAuthID)
	}

	// 2. Probar fallback de login por sEmail
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", testEmail, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL usp_User_USloginGame por Email fallo: %v", err)
	}
	if nResult != 0 || outUserNo != int(userNo) || outAuthID != 1 {
		t.Fatalf("Login por Email esperado (0, %d, 1), obtenido (%d, %d, %d)", userNo, nResult, outUserNo, outAuthID)
	}

	// 3. Probar clave incorrecta -> -2
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", testUser, "badpasshash", "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL con clave incorrecta fallo: %v", err)
	}
	if nResult != -2 {
		t.Fatalf("Se esperaba nResult = -2 para clave incorrecta, obtenido: %d", nResult)
	}

	// 4. Probar usuario no existente -> -1
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", "non_existent_user_999", passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL con usuario inexistente fallo: %v", err)
	}
	if nResult != -1 {
		t.Fatalf("Se esperaba nResult = -1 para usuario inexistente, obtenido: %d", nResult)
	}
}

// TestFase02_LoginUS_AutoCreationPrefixR prueba el alta automatica con prefijo 'r_'
// gobernada por Auto_Account_Creation.Allow_Creation.
func TestFase02_LoginUS_AutoCreationPrefixR(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const baseUser = "fase02_auto_r"
	prefixedUser := "r_" + baseUser
	passHash := md5Hash("testpass")

	// Limpiar
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", baseUser)
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", prefixedUser)

	// 1. Caso Allow_Creation = 1 -> debe crear la cuenta sin 'r_' y retornar 0 (exito)
	_, err := db.Exec("UPDATE Auto_Account_Creation SET Allow_Creation = 1")
	if err != nil {
		t.Fatalf("Error habilitando Auto_Account_Creation: %v", err)
	}

	var nResult, outUserNo, outAuthID int
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", prefixedUser, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL usp_User_USloginGame r_ con Allow_Creation=1 fallo: %v", err)
	}
	if nResult != 0 {
		t.Fatalf("Se esperaba nResult = 0 para auto-creacion r_, obtenido: %d", nResult)
	}

	// Verificar que en tUser se guardo con baseUser (sin 'r_')
	var savedUser string
	err = db.QueryRow("SELECT sUserID FROM tUser WHERE nUserNo = ?", outUserNo).Scan(&savedUser)
	if err != nil || savedUser != baseUser {
		t.Fatalf("La cuenta debio guardarse como '%s', pero se obtuvo '%s' (err: %v)", baseUser, savedUser, err)
	}

	// 2. Caso Allow_Creation = 0 con cuenta nueva -> debe rechazar (nResult = -1) y NO crear cuenta
	const baseUserDisabled = "fase02_r_disabled"
	prefixedUserDisabled := "r_" + baseUserDisabled
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", baseUserDisabled)

	_, err = db.Exec("UPDATE Auto_Account_Creation SET Allow_Creation = 0")
	if err != nil {
		t.Fatalf("Error deshabilitando Auto_Account_Creation: %v", err)
	}

	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", prefixedUserDisabled, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL con Allow_Creation=0 fallo: %v", err)
	}
	if nResult != -1 {
		t.Fatalf("Se esperaba nResult = -1 cuando Allow_Creation=0, obtenido: %d", nResult)
	}

	var countDisabled int
	_ = db.QueryRow("SELECT COUNT(*) FROM tUser WHERE sUserID = ?", baseUserDisabled).Scan(&countDisabled)
	if countDisabled > 0 {
		t.Fatalf("No debio haberse creado el usuario con Allow_Creation=0")
	}

	// Restaurar Allow_Creation = 1
	_, _ = db.Exec("UPDATE Auto_Account_Creation SET Allow_Creation = 1")
}

// TestFase02_UserChkBlock prueba la expiracion selectiva y desbloqueo de usp_User_chkBlock.
func TestFase02_UserChkBlock(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const blockUser = "test_fase02_block"
	passHash := md5Hash("blockpass")

	_, _ = db.Exec("DELETE FROM tUserBlock WHERE nUserNo IN (SELECT nUserNo FROM tUser WHERE sUserID = ?)", blockUser)
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", blockUser)

	res, err := db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, ?, ?, 1, 0, 1, '127.0.0.1', NOW())
	`, blockUser, passHash, blockUser)
	if err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}
	userNo, _ := res.LastInsertId()

	// 1. Insertar bloqueo temporal ya expirado (dEndDate en el pasado, nType = 1)
	_, err = db.Exec(`
		INSERT INTO tUserBlock (nUserNo, nType, sComment, dEndDate, bStatus)
		VALUES (?, 1, 'Temp block expired', NOW() - INTERVAL 1 HOUR, 1)
	`, userNo)
	if err != nil {
		t.Fatalf("Error insertando tUserBlock expirado: %v", err)
	}

	// Ejecutar usp_User_chkBlock -> debe expirar el bloqueo y desbloquear al usuario
	var ret int
	err = db.QueryRow("CALL usp_User_chkBlock(?, @ret)", userNo).Scan(&ret)
	if err != nil {
		t.Fatalf("CALL usp_User_chkBlock fallo: %v", err)
	}
	if ret != 1 {
		t.Fatalf("Se esperaba ret = 1 (desbloqueado), obtenido: %d", ret)
	}

	// Verificar tUser.bIsBlock = 0
	var isBlock int
	_ = db.QueryRow("SELECT bIsBlock FROM tUser WHERE nUserNo = ?", userNo).Scan(&isBlock)
	if isBlock != 0 {
		t.Fatalf("tUser.bIsBlock debio actualizarse a 0")
	}

	// 2. Insertar bloqueo activo permanente (nType = 99, bStatus = 1, dEndDate NULL o futuro)
	_, err = db.Exec(`
		INSERT INTO tUserBlock (nUserNo, nType, sComment, dEndDate, bStatus)
		VALUES (?, 99, 'Permanent block', NULL, 1)
	`, userNo)
	if err != nil {
		t.Fatalf("Error insertando bloqueo permanente: %v", err)
	}

	err = db.QueryRow("CALL usp_User_chkBlock(?, @ret)", userNo).Scan(&ret)
	if err != nil {
		t.Fatalf("CALL usp_User_chkBlock permanente fallo: %v", err)
	}
	if ret != 0 {
		t.Fatalf("Se esperaba ret = 0 (sigue bloqueado), obtenido: %d", ret)
	}

	// Verificar login US -> debe retornar -3 (bloqueado)
	var nResult, outUserNo, outAuthID int
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", blockUser, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL usp_User_USloginGame fallo: %v", err)
	}
	if nResult != -3 {
		t.Fatalf("Se esperaba nResult = -3 para cuenta bloqueada, obtenido: %d", nResult)
	}
}

// TestFase02_UserAuth_Permission prueba el bloqueo a nivel autorizacion con tUserAuth.bIsLoginAble.
func TestFase02_UserAuth_Permission(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const authUser = "test_fase02_auth"
	passHash := md5Hash("authpass")

	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", authUser)

	// Crear usuario con nAuthID = 2
	res, err := db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, ?, ?, 0, 0, 2, '127.0.0.1', NOW())
	`, authUser, passHash, authUser)
	if err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}
	_ = res

	// Deshabilitar login para AuthID = 2
	var ret int
	err = db.QueryRow("CALL usp_User_setLoginable(2, 0, @ret)").Scan(&ret)
	if err != nil {
		t.Fatalf("CALL usp_User_setLoginable fallo: %v", err)
	}

	// Login debe fallar con nResult = -4
	var nResult, outUserNo, outAuthID int
	err = db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", authUser, passHash, "127.0.0.1").
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		t.Fatalf("CALL usp_User_USloginGame fallo: %v", err)
	}
	if nResult != -4 {
		t.Fatalf("Se esperaba nResult = -4 cuando bIsLoginAble=0, obtenido: %d", nResult)
	}

	// Restaurar AuthID = 2 a bIsLoginAble = 1
	_ = db.QueryRow("CALL usp_User_setLoginable(2, 1, @ret)").Scan(&ret)
}

// TestFase02_Cash_Operations prueba el ciclo completo de Cash (status = 0)
func TestFase02_Cash_Operations(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUserNo = 888881
	_, _ = db.Exec("DELETE FROM tCash WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tCashLog WHERE userNo = ?", testUserNo)

	// 1. SetCash 500
	var ret int
	err := db.QueryRow("CALL usp_Cash_setCash(?, 1, 500, 'TEST_CREDIT')", testUserNo).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_Cash_setCash fallo: %v, ret=%d", err, ret)
	}

	// 2. GetCash -> debe devolver 500
	var cash int
	err = db.QueryRow("CALL usp_Cash_getCash(?, 1, @cash)", testUserNo).Scan(&cash)
	if err != nil || cash != 500 {
		t.Fatalf("usp_Cash_getCash esperado 500, obtenido %d (err: %v)", cash, err)
	}

	// 3. UseCash 200 -> exito (ret = 0)
	err = db.QueryRow("CALL usp_Cash_useCash(?, 1, 200, 'TEST_BUY')", testUserNo).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_Cash_useCash esperado 0, obtenido %d (err: %v)", ret, err)
	}

	// GetCash -> debe devolver 300
	_ = db.QueryRow("CALL usp_Cash_getCash(?, 1, @cash)", testUserNo).Scan(&cash)
	if cash != 300 {
		t.Fatalf("Saldo esperado 300, obtenido %d", cash)
	}

	// 4. UseCash excesivo 9999 -> fallo (ret = -1)
	err = db.QueryRow("CALL usp_Cash_useCash(?, 1, 9999, 'TEST_FAIL')", testUserNo).Scan(&ret)
	if err != nil || ret != -1 {
		t.Fatalf("UseCash excesivo debio retornar -1, obtenido %d", ret)
	}

	// 5. usp_Event_addCash 150 (Firma oficial ODBC de 3 parámetros en Account.md: {CALL usp_Event_addCash(?,?,?)})
	err = db.QueryRow("CALL usp_Event_addCash(?, 150, @ret)", testUserNo).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_Event_addCash fallo: %v, ret=%d", err, ret)
	}

	_ = db.QueryRow("CALL usp_Cash_getCash(?, 1, @cash)", testUserNo).Scan(&cash)
	if cash != 450 {
		t.Fatalf("Saldo tras evento esperado 450, obtenido %d", cash)
	}
}

// TestFase02_ChargeItem_Lifecycle prueba entrega, listado, marcado de uso y borrado con log.
func TestFase02_ChargeItem_Lifecycle(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUserNo = 888882
	_, _ = db.Exec("DELETE FROM tChargeItem WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tChargeItemDelLog WHERE userNo = ?", testUserNo)

	// 1. Insertar item
	var ret int
	err := db.QueryRow("CALL usp_Charge_ItemInsert(?, 0, 10050, 1, @ret)", testUserNo).Scan(&ret)
	if err != nil || ret != 1 {
		t.Fatalf("usp_Charge_ItemInsert esperado ret=1, obtenido %d (err: %v)", ret, err)
	}

	// 2. SelectList -> debe contener el item con isDraw = 0
	rows, err := db.Query("CALL usp_Charge_ItemSelectList(?)", testUserNo)
	if err != nil {
		t.Fatalf("usp_Charge_ItemSelectList fallo: %v", err)
	}
	var rowNo, goodsNo, amount int
	var regDate string
	found := false
	for rows.Next() {
		if err := rows.Scan(&rowNo, &goodsNo, &amount, &regDate); err != nil {
			t.Fatalf("Scan error: %v", err)
		}
		if goodsNo == 10050 {
			found = true
			break
		}
	}
	rows.Close()

	if !found {
		t.Fatalf("No se encontro el item 10050 en usp_Charge_ItemSelectList")
	}

	// 3. Draw item -> exito (ret = 1)
	err = db.QueryRow("CALL usp_Charge_ItemDraw(?, ?, @ret)", rowNo, testUserNo).Scan(&ret)
	if err != nil || ret != 1 {
		t.Fatalf("usp_Charge_ItemDraw esperado ret=1, obtenido %d (err: %v)", ret, err)
	}

	// 4. Draw item repetido -> debe rechazar con ret = -1 (ya usado)
	err = db.QueryRow("CALL usp_Charge_ItemDraw(?, ?, @ret)", rowNo, testUserNo).Scan(&ret)
	if err != nil || ret != -1 {
		t.Fatalf("usp_Charge_ItemDraw repetido debio retornar -1, obtenido %d", ret)
	}

	// 5. Insertar otro item y probar borrado con copia a tChargeItemDelLog
	_ = db.QueryRow("CALL usp_Charge_ItemInsert(?, 0, 10051, 1, @ret)", testUserNo).Scan(&ret)
	var rowNo2 int
	_ = db.QueryRow("SELECT rowNo FROM tChargeItem WHERE userNo = ? AND goodsNo = 10051", testUserNo).Scan(&rowNo2)

	err = db.QueryRow("CALL usp_Charge_ItemDel(?, ?, @ret)", rowNo2, testUserNo).Scan(&ret)
	if err != nil || ret != 1 {
		t.Fatalf("usp_Charge_ItemDel esperado ret=1, obtenido %d (err: %v)", ret, err)
	}

	// Verificar que se copio a tChargeItemDelLog con el mismo rowNo y se borro de tChargeItem
	var logCount, activeCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM tChargeItemDelLog WHERE rowNo = ?", rowNo2).Scan(&logCount)
	_ = db.QueryRow("SELECT COUNT(*) FROM tChargeItem WHERE rowNo = ?", rowNo2).Scan(&activeCount)

	if logCount != 1 || activeCount != 0 {
		t.Fatalf("Borrado incorrecto: logCount=%d (esperado 1), activeCount=%d (esperado 0)", logCount, activeCount)
	}
}

// TestFase02_ItemShop_BuyOrder prueba la transaccion de compra en Item Shop
func TestFase02_ItemShop_BuyOrder(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUserNo = 888883
	const testGoodsNo = 999901
	passHash := md5Hash("itempass")

	// Crear categoria, usuario y producto
	_, _ = db.Exec("DELETE FROM tOrder WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tCash WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tAccountItem WHERE goodsNo = ?", testGoodsNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)

	_, _ = db.Exec(`
		INSERT INTO tItemCategory (category, name, parent, status)
		VALUES ('WEAPON', 'Weapon', '', 1)
		ON DUPLICATE KEY UPDATE name = VALUES(name)
	`)

	_, _ = db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, 'item_buyer', ?, 'buyer', 0, 0, 1, '127.0.0.1', NOW())
	`, testUserNo, passHash)

	_, _ = db.Exec(`
		INSERT INTO tAccountItem (goodsNo, billingNo, category, groupNo, name, itemType, unit, price, discount, isSell, sellCount)
		VALUES (?, 101, 'WEAPON', 10, 'Super Sword', 'ITEM', 1, 100, 0, 1, 0)
	`, testGoodsNo)

	// Acreditar 250 de Cash
	var ret int
	_ = db.QueryRow("CALL usp_Cash_setCash(?, 1, 250, 'SHOP_TEST')", testUserNo).Scan(&ret)

	// Crear orden de compra por 100 cash (10 params)
	var orderNo int
	err := db.QueryRow("CALL usp_Item_setOrder(0, ?, ?, 100, 1, 0, 100, 1, '127.0.0.1', @orderNo)", testUserNo, testGoodsNo).
		Scan(&orderNo)
	if err != nil || orderNo <= 0 {
		t.Fatalf("usp_Item_setOrder fallo: err=%v, orderNo=%d", err, orderNo)
	}

	// Ejecutar compra
	var buyRet int
	err = db.QueryRow("CALL usp_Item_setBuyOrder(?, ?, 0, @ret)", orderNo, testUserNo).Scan(&buyRet)
	if err != nil || buyRet != 0 {
		t.Fatalf("usp_Item_setBuyOrder fallo: err=%v, ret=%d", err, buyRet)
	}

	// Verificar:
	// a) Saldo reducido a 150
	var currentCash int
	_ = db.QueryRow("CALL usp_Cash_getCash(?, 1, @cash)", testUserNo).Scan(&currentCash)
	if currentCash != 150 {
		t.Fatalf("Saldo esperado 150, obtenido %d", currentCash)
	}

	// b) Item entregado en tChargeItem
	var chargeItemCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM tChargeItem WHERE userNo = ? AND goodsNo = ? AND orderNo = ?", testUserNo, testGoodsNo, orderNo).
		Scan(&chargeItemCount)
	if chargeItemCount != 1 {
		t.Fatalf("Item no acreditado en tChargeItem")
	}

	// c) tAccountItem.sellCount incrementado a 1
	var sellCount int
	_ = db.QueryRow("SELECT sellCount FROM tAccountItem WHERE goodsNo = ?", testGoodsNo).Scan(&sellCount)
	if sellCount != 1 {
		t.Fatalf("sellCount esperado 1, obtenido %d", sellCount)
	}
}

// TestFase02_ReturnUser_And_SerialKey prueba cupones unicos y premios a usuarios retornantes.
func TestFase02_ReturnUser_And_SerialKey(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUserNo = 888884
	passHash := md5Hash("retpass")

	_, _ = db.Exec("DELETE FROM tChargeItem WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tReturnUser WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tSerialKey WHERE userNo = ?", testUserNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)

	_, _ = db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, 'return_user', ?, 'retuser', 0, 0, 1, '127.0.0.1', NOW())
	`, testUserNo, passHash)

	// 1. usp_ReturnUser_SetUserItem (Firma oficial de 2 parámetros en Account.bak / Account.md: {CALL usp_ReturnUser_SetUserItem(?,?)})
	var ret int
	err := db.QueryRow("CALL usp_ReturnUser_SetUserItem(?, @ret)", testUserNo).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_ReturnUser_SetUserItem fallo: %v, ret=%d", err, ret)
	}

	// Verificar registro en tReturnUser y tChargeItem
	var retCount, chargeCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM tReturnUser WHERE nUserNo = ? AND nItemNo = 30001", testUserNo).Scan(&retCount)
	_ = db.QueryRow("SELECT COUNT(*) FROM tChargeItem WHERE userNo = ? AND goodsNo = 30001", testUserNo).Scan(&chargeCount)

	if retCount != 1 || chargeCount != 1 {
		t.Fatalf("Error acreditando return user item: retCount=%d, chargeCount=%d", retCount, chargeCount)
	}

	// usp_ReturnUser_GetDatediffLastlog
	var diffSec int
	err = db.QueryRow("CALL usp_ReturnUser_GetDatediffLastlog(?, @diff)", testUserNo).Scan(&diffSec)
	if err != nil || diffSec < 0 {
		t.Fatalf("usp_ReturnUser_GetDatediffLastlog fallo: %v, diffSec=%d", err, diffSec)
	}

	// 2. usp_Item_setSerialKey
	const serialKey int64 = 9876543210123
	_, _ = db.Exec("DELETE FROM tSerialKey WHERE serialKey = ?", serialKey)

	err = db.QueryRow("CALL usp_Item_setSerialKey(?, ?, @ret)", testUserNo, serialKey).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_Item_setSerialKey insercion inicial fallo: %v, ret=%d", err, ret)
	}

	// Intentar reutilizar la misma serialKey -> debe rechazar con -1 sin sobrescribir
	err = db.QueryRow("CALL usp_Item_setSerialKey(?, ?, @ret)", testUserNo, serialKey).Scan(&ret)
	if err != nil || ret != -1 {
		t.Fatalf("usp_Item_setSerialKey duplicada debio retornar -1, obtenido: %d", ret)
	}
}

// TestFase02_Tutorial prueba usp_Tutorial_getIsSkip y setCharInfo
func TestFase02_Tutorial(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const testUserNo = 888885
	const testCharNo = 7701
	_, _ = db.Exec("DELETE FROM tTutorial WHERE nUserNo = ?", testUserNo)

	// Inicialmente no puede saltar tutorial -> 0
	var isSkip int
	err := db.QueryRow("CALL usp_Tutorial_getIsSkip(?, @ret)", testUserNo).Scan(&isSkip)
	if err != nil || isSkip != 0 {
		t.Fatalf("Tutorial isSkip inicial esperado 0, obtenido %d (err: %v)", isSkip, err)
	}

	// Completar tutorial (nState = 1)
	var ret int
	err = db.QueryRow("CALL usp_Tutorial_setCharInfo(?, ?, 1, 5, @ret)", testUserNo, testCharNo).Scan(&ret)
	if err != nil || ret != 0 {
		t.Fatalf("usp_Tutorial_setCharInfo fallo: %v, ret=%d", err, ret)
	}

	// Ahora debe poder saltar tutorial -> 1
	err = db.QueryRow("CALL usp_Tutorial_getIsSkip(?, @ret)", testUserNo).Scan(&isSkip)
	if err != nil || isSkip != 1 {
		t.Fatalf("Tutorial isSkip tras completar esperado 1, obtenido %d", isSkip)
	}
}

// TestFase02_Pcbang_And_Key prueba verificacion de IP de cibercafe y claves web
func TestFase02_Pcbang_And_Key(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. Pcbang
	_, _ = db.Exec("DELETE FROM tPcbang WHERE sName = 'TEST_CAFE'")
	_, _ = db.Exec(`
		INSERT INTO tPcbang (sIPStart, sIPEnd, sName, bIsActive)
		VALUES ('192.168.10.1', '192.168.10.254', 'TEST_CAFE', 1)
	`)

	var pcbangNo int
	// IP dentro del rango
	err := db.QueryRow("CALL usp_Pcbang_chkIP('192.168.10.50', @no)").Scan(&pcbangNo)
	if err != nil || pcbangNo <= 0 {
		t.Fatalf("usp_Pcbang_chkIP IP en rango fallo: %v, pcbangNo=%d", err, pcbangNo)
	}

	// IP fuera del rango -> 0
	err = db.QueryRow("CALL usp_Pcbang_chkIP('10.0.0.1', @no)").Scan(&pcbangNo)
	if err != nil || pcbangNo != 0 {
		t.Fatalf("usp_Pcbang_chkIP IP fuera de rango esperado 0, obtenido %d", pcbangNo)
	}

	// 2. tUserKey
	const testUserNo = 888886
	_, _ = db.Exec("DELETE FROM tUserKey WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec(`
		INSERT INTO tUserKey (nUserNo, sWebKey, dIssueDate, dExpireDate)
		VALUES (?, 'VALID_WEB_KEY_1234567890ABCDEF', NOW(), NOW() + INTERVAL 1 HOUR)
	`, testUserNo)

	var nRet int
	var sWebKey string
	var dExpire sql.NullTime
	err = db.QueryRow("CALL p_UserKey_Get_FromAccount(?, @k, @e, @r)", testUserNo).Scan(&nRet, &sWebKey, &dExpire)
	if err != nil || nRet != 0 || sWebKey != "VALID_WEB_KEY_1234567890ABCDEF" {
		t.Fatalf("p_UserKey_Get_FromAccount fallo: %v, nRet=%d, key=%s", err, nRet, sWebKey)
	}

	// 3. Teenager play time
	_, _ = db.Exec("DELETE FROM tUserPlayRemainTime WHERE nUserNo = ?", testUserNo)
	_, _ = db.Exec(`
		INSERT INTO tUserPlayRemainTime (nUserNo, nRemainSec, bIsTeenager, dLastUpdate)
		VALUES (?, 7200, 1, NOW())
	`, testUserNo)

	var remainSec int
	var isTeenager int
	err = db.QueryRow("CALL p_UserPlayRemainTime_Get_FromAccount(?, @rem, @teen, @ret)", testUserNo).
		Scan(&nRet, &remainSec, &isTeenager)
	if err != nil || nRet != 0 || remainSec != 7200 || isTeenager != 1 {
		t.Fatalf("p_UserPlayRemainTime_Get fallo: %v, nRet=%d, rem=%d, teen=%d", err, nRet, remainSec, isTeenager)
	}

	// Descontar 1200 segundos
	err = db.QueryRow("CALL p_UserPlayRemainTime_Set_FromAccount(?, 1200, @rem, @ret)", testUserNo).
		Scan(&nRet, &remainSec)
	if err != nil || nRet != 0 || remainSec != 6000 {
		t.Fatalf("p_UserPlayRemainTime_Set fallo: %v, nRet=%d, rem=%d", err, nRet, remainSec)
	}
}
