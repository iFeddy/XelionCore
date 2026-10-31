package audit_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func getTestDB(t *testing.T) *sql.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "lerd"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "xelion_online"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open DB connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping DB: %v", err)
	}
	return db
}

// TestFase01_Control1 verifica exhaustivamente contra information_schema el 100% de tablas,
// vistas, columnas y tipos requeridos por la Fase 01 y el Control 1 (Sección 9.1).
func TestFase01_Control1(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const schema = "xelion_online"

	// 1. Obtener todas las tablas y vistas existentes
	rows, err := db.Query(`
		SELECT TABLE_NAME, TABLE_TYPE 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = ?
	`, schema)
	if err != nil {
		t.Fatalf("Error querying information_schema.TABLES: %v", err)
	}
	defer rows.Close()

	tables := make(map[string]string)
	for rows.Next() {
		var name, tableType string
		if err := rows.Scan(&name, &tableType); err != nil {
			t.Fatalf("Scan error: %v", err)
		}
		tables[strings.ToLower(name)] = tableType
	}

	// 2. Inventario exhaustivo de las 76 entidades requeridas
	expectedTables := []string{
		// 34 tablas de Account.bak
		"Event_Christmas_Item", "Event_Christmas_Log", "FOUS",
		"last2_onson_users_1", "last2_onson_users_2", "login_issue_check",
		"onson_0419", "onson_users_0314", "tCash", "tCashLog",
		"tChannelLog", "tCharacterMoveLog", "tChargeItem", "tChargeItemDelLog",
		"tDrawComment", "tAccountItem", "tItemCategory", "tLv120",
		"tOrder", "tPasswordLog", "tPay", "tPayAgreement",
		"tPayResult", "tPayType", "tReturnUser", "tSerialKey",
		"tTutorial", "tUser", "tUserAuth", "tUserBlock",
		"tUserMergeYet", "tUserOption", "tbSellitem", "tuseroption_2",

		// 11 tablas auxiliares / ODBC de Account
		"Auto_Account_Creation", "tPcbang", "tUserKey", "tUserPlayRemainTime",
		"tZipcode", "tbl_user", "tbl_user_test", "tUserTemp",
		"tUserTemp2", "tUser_ori", "temp_newUser",

		// 10 tablas de AccountLog.bak
		"tAccountLog", "tEventWorldCup", "tEvent_Friend", "tEvent_Friend_HideMsg",
		"tEventiPad", "tIPBlock", "tLogout", "tReport",
		"tAccountDemographics", "tUserLog",

		// 3 tablas operativas de AccountLog
		"tEvent_L20", "tHackScanLog", "tSyslog",

		// 13 tablas de StatisticsData.bak
		"TopCharacter", "TopCharacter_MAY", "tAge", "tDayActiveUser",
		"tDayClassLevel", "tDayClassLevelPlay", "tDayUser", "tGender",
		"tLocation", "tQuestTimesQ20016", "tQuestTimesQ2079", "tQuestTimesQ63011",
		"tSameTimeChar",

		// 2 tablas operativas de LoginService
		"tWorldServer", "tActiveAccountSession",
	}

	for _, tbl := range expectedTables {
		if _, ok := tables[strings.ToLower(tbl)]; !ok {
			t.Errorf("[Fase 01] Falta tabla requerida: %s", tbl)
		}
	}

	// 3. Verificar las 3 Vistas Oficiales
	expectedViews := []string{"VT_ITEM_SELL", "VT_ITEM_GROUP", "VT_PAY"}
	for _, v := range expectedViews {
		ttype, ok := tables[strings.ToLower(v)]
		if !ok {
			t.Errorf("[Fase 01] Falta vista requerida: %s", v)
		} else if ttype != "VIEW" {
			t.Errorf("[Fase 01] %s existe pero no es de tipo VIEW (tipo actual: %s)", v, ttype)
		}
	}

	// 4. Invariantes de columnas críticas (Control 1 / Secciones 2.1, 2.2, 2.3)

	// Invariante 4.1: tAccountItem tiene las 22 columnas exactas
	colRows, err := db.Query(`
		SELECT COLUMN_NAME 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tAccountItem'
		ORDER BY ORDINAL_POSITION
	`, schema)
	if err != nil {
		t.Fatalf("Error querying tAccountItem columns: %v", err)
	}
	defer colRows.Close()

	var accountItemCols []string
	for colRows.Next() {
		var col string
		if err := colRows.Scan(&col); err != nil {
			t.Fatalf("Scan error: %v", err)
		}
		accountItemCols = append(accountItemCols, strings.ToLower(col))
	}
	if len(accountItemCols) != 22 {
		t.Errorf("[Fase 01] tAccountItem debe tener exactamente 22 columnas, tiene %d: %v", len(accountItemCols), accountItemCols)
	}

	expectedAccountItemCols := []string{
		"goodsno", "billingno", "category", "groupno", "name", "imagename",
		"itemtype", "unit", "unitdescription", "price", "discount",
		"discountstartdate", "discountenddate", "isnew", "ispopular",
		"isrecommand", "issell", "sellstartdate", "sellenddate", "sellcount",
		"registerdate", "description",
	}
	for _, expectedCol := range expectedAccountItemCols {
		found := false
		for _, col := range accountItemCols {
			if col == expectedCol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("[Fase 01] tAccountItem falta columna obligatoria: %s", expectedCol)
		}
	}

	// Invariante 4.2: tCash.status con default '0'
	var cashStatusDefault sql.NullString
	err = db.QueryRow(`
		SELECT COLUMN_DEFAULT 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tCash' AND COLUMN_NAME = 'status'
	`, schema).Scan(&cashStatusDefault)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tCash.status: %v", err)
	} else if !cashStatusDefault.Valid || (cashStatusDefault.String != "0" && cashStatusDefault.String != "'0'") {
		t.Errorf("[Fase 01] tCash.status debe tener DEFAULT '0', actual: %s", cashStatusDefault.String)
	}

	// Invariante 4.3: tChargeItem.orderNo bigint
	var chargeItemOrderType string
	err = db.QueryRow(`
		SELECT DATA_TYPE 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tChargeItem' AND COLUMN_NAME = 'orderNo'
	`, schema).Scan(&chargeItemOrderType)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tChargeItem.orderNo: %v", err)
	} else if chargeItemOrderType != "bigint" {
		t.Errorf("[Fase 01] tChargeItem.orderNo debe ser bigint, actual: %s", chargeItemOrderType)
	}

	// Invariante 4.4: tChargeItemDelLog.orderNo bigint y rowNo sin auto_increment
	var chargeDelOrderType, chargeDelRowExtra string
	err = db.QueryRow(`
		SELECT DATA_TYPE 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tChargeItemDelLog' AND COLUMN_NAME = 'orderNo'
	`, schema).Scan(&chargeDelOrderType)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tChargeItemDelLog.orderNo: %v", err)
	} else if chargeDelOrderType != "bigint" {
		t.Errorf("[Fase 01] tChargeItemDelLog.orderNo debe ser bigint, actual: %s", chargeDelOrderType)
	}

	err = db.QueryRow(`
		SELECT EXTRA 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tChargeItemDelLog' AND COLUMN_NAME = 'rowNo'
	`, schema).Scan(&chargeDelRowExtra)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tChargeItemDelLog.rowNo: %v", err)
	} else if strings.Contains(strings.ToLower(chargeDelRowExtra), "auto_increment") {
		t.Errorf("[Fase 01] tChargeItemDelLog.rowNo NO debe ser auto_increment, actual: %s", chargeDelRowExtra)
	}

	// Invariante 4.5: tPay.orderNo con auto_increment
	var payOrderExtra string
	err = db.QueryRow(`
		SELECT EXTRA 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tPay' AND COLUMN_NAME = 'orderNo'
	`, schema).Scan(&payOrderExtra)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tPay.orderNo: %v", err)
	} else if !strings.Contains(strings.ToLower(payOrderExtra), "auto_increment") {
		t.Errorf("[Fase 01] tPay.orderNo debe ser auto_increment, actual: %s", payOrderExtra)
	}

	// Invariante 4.6: tUser NO tiene bIsLoginAble y tiene las 12 columnas exactas
	tUserColRows, err := db.Query(`
		SELECT COLUMN_NAME 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tUser'
	`, schema)
	if err != nil {
		t.Fatalf("Error querying tUser columns: %v", err)
	}
	defer tUserColRows.Close()

	var tUserCols []string
	for tUserColRows.Next() {
		var col string
		if err := tUserColRows.Scan(&col); err != nil {
			t.Fatalf("Scan error: %v", err)
		}
		tUserCols = append(tUserCols, strings.ToLower(col))
	}

	for _, col := range tUserCols {
		if col == "bisloginable" {
			t.Errorf("[Fase 01] tUser NO debe tener la columna inventada bIsLoginAble")
		}
	}
	if len(tUserCols) != 12 {
		t.Errorf("[Fase 01] tUser debe tener exactamente las 12 columnas oficiales de Account.bak, tiene %d: %v", len(tUserCols), tUserCols)
	}

	expectedTUserCols := []string{
		"nuserno", "suserid", "suserpw", "susername", "bisblock",
		"bisdelete", "nauthid", "suserip", "ddate", "suserno",
		"blsblock", "blsdelete",
	}
	for _, expected := range expectedTUserCols {
		found := false
		for _, col := range tUserCols {
			if col == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("[Fase 01] tUser falta columna oficial: %s", expected)
		}
	}

	// Invariante 4.7: FOUS poblado con GameVersionKey = '10022024000000'
	var gameVerKey string
	err = db.QueryRow(`
		SELECT GameVersionKey 
		FROM FOUS 
		WHERE GameAccountTitle = 'FiestaOnline'
	`).Scan(&gameVerKey)
	if err != nil {
		t.Errorf("[Fase 01] Error leyendo FOUS.GameVersionKey: %v", err)
	} else if gameVerKey != "10022024000000" {
		t.Errorf("[Fase 01] FOUS.GameVersionKey esperado '10022024000000', actual: '%s'", gameVerKey)
	}

	// Invariante 4.8: tIPBlock con columnas generadas nIPStartNum y nIPEndNum
	var startNumExtra, endNumExtra string
	err = db.QueryRow(`
		SELECT EXTRA 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tIPBlock' AND COLUMN_NAME = 'nIPStartNum'
	`, schema).Scan(&startNumExtra)
	if err != nil {
		t.Errorf("[Fase 01] tIPBlock falta columna generada nIPStartNum: %v", err)
	} else if !strings.Contains(strings.ToLower(startNumExtra), "generated") {
		t.Errorf("[Fase 01] tIPBlock.nIPStartNum debe ser GENERATED STORED, actual extra: %s", startNumExtra)
	}

	err = db.QueryRow(`
		SELECT EXTRA 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tIPBlock' AND COLUMN_NAME = 'nIPEndNum'
	`, schema).Scan(&endNumExtra)
	if err != nil {
		t.Errorf("[Fase 01] tIPBlock falta columna generada nIPEndNum: %v", err)
	} else if !strings.Contains(strings.ToLower(endNumExtra), "generated") {
		t.Errorf("[Fase 01] tIPBlock.nIPEndNum debe ser GENERATED STORED, actual extra: %s", endNumExtra)
	}

	// Invariante 4.9: tLogout.level es smallint
	var logoutLevelType string
	err = db.QueryRow(`
		SELECT DATA_TYPE 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tLogout' AND COLUMN_NAME = 'level'
	`, schema).Scan(&logoutLevelType)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tLogout.level: %v", err)
	} else if logoutLevelType != "smallint" {
		t.Errorf("[Fase 01] tLogout.level debe ser smallint, actual: %s", logoutLevelType)
	}

	// Invariante 4.10: tDrawComment.sComment >= 1000
	var commentLen int
	err = db.QueryRow(`
		SELECT CHARACTER_MAXIMUM_LENGTH 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tDrawComment' AND COLUMN_NAME = 'sComment'
	`, schema).Scan(&commentLen)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tDrawComment.sComment: %v", err)
	} else if commentLen < 1000 {
		t.Errorf("[Fase 01] tDrawComment.sComment debe ser al menos VARCHAR(1000), actual: %d", commentLen)
	}

	// Invariante 4.11: tUserOption.sEmail >= 256
	var emailLen int
	err = db.QueryRow(`
		SELECT CHARACTER_MAXIMUM_LENGTH 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tUserOption' AND COLUMN_NAME = 'sEmail'
	`, schema).Scan(&emailLen)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tUserOption.sEmail: %v", err)
	} else if emailLen < 256 {
		t.Errorf("[Fase 01] tUserOption.sEmail debe ser al menos VARCHAR(256), actual: %d", emailLen)
	}

	// Invariante 4.12: tbSellitem.baseDate varchar
	var sellBaseDateType string
	err = db.QueryRow(`
		SELECT DATA_TYPE 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tbSellitem' AND COLUMN_NAME = 'baseDate'
	`, schema).Scan(&sellBaseDateType)
	if err != nil {
		t.Errorf("[Fase 01] Error verificando tbSellitem.baseDate: %v", err)
	} else if sellBaseDateType != "varchar" {
		t.Errorf("[Fase 01] tbSellitem.baseDate debe ser varchar, actual: %s", sellBaseDateType)
	}

	// Invariante 4.13: tTutorial existe como tabla única
	var tutorialCount int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'tTutorial'
	`, schema).Scan(&tutorialCount)
	if err != nil || tutorialCount != 1 {
		t.Errorf("[Fase 01] tTutorial debe existir exactamente una vez en information_schema, encontrados: %d", tutorialCount)
	}
}
