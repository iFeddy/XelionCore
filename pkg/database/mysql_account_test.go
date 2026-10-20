package database_test

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
)

func setupTestDB(t *testing.T) (*sql.DB, *database.MySQLAccountRepository) {
	cfg, err := config.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping MySQL: %v", err)
	}

	repo := database.NewMySQLAccountRepository(db)
	return db, repo
}

func TestMySQLAccountRepository_GetAccountByUsername(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	// 1. Validar usuario 'admin' existente
	admin, err := repo.GetAccountByUsername("admin")
	if err != nil {
		t.Fatalf("Failed to get admin account: %v", err)
	}
	if admin.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", admin.Username)
	}
	if admin.Admin != 4 {
		t.Errorf("Expected Admin auth tier 4, got %d", admin.Admin)
	}
	expectedPW := md5HashString("admin")
	if admin.Password != expectedPW {
		t.Errorf("Expected password hash %s, got %s", expectedPW, admin.Password)
	}

	// 2. Validar usuario inexistente
	_, err = repo.GetAccountByUsername("non_existent_user_99999")
	if !errors.Is(err, database.ErrAccountNotFound) {
		t.Errorf("Expected ErrAccountNotFound, got %v", err)
	}
}

func TestMySQLAccountRepository_CreateAccount(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	testUser := "test_created_user"
	// Limpiar si ya existía de pruebas previas
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", testUser)

	acc, err := repo.CreateAccount(testUser, "secret123")
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
	if acc.Username != testUser {
		t.Errorf("Expected username '%s', got '%s'", testUser, acc.Username)
	}
	expectedHash := md5HashString("secret123")
	if acc.Password != expectedHash {
		t.Errorf("Expected MD5 hash %s, got %s", expectedHash, acc.Password)
	}

	// Limpieza posterior
	_, _ = db.Exec("DELETE FROM tUser WHERE sUserID = ?", testUser)
}

func TestMySQLAccountRepository_LogAccountAndIPBlock(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	// 1. Test LogAccount
	err := repo.LogAccount("admin", 1, 1, 0, 0, "127.0.0.1")
	if err != nil {
		t.Fatalf("Failed to log account: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tAccountLog WHERE sUserID = 'admin' AND nType = 1").Scan(&count)
	if err != nil || count == 0 {
		t.Errorf("Expected log row in tAccountLog, got count=%d, err=%v", count, err)
	}

	// 2. Test IsIPBlocked
	blocked, err := repo.IsIPBlocked("192.168.1.100")
	if err != nil {
		t.Fatalf("Failed to check IP block: %v", err)
	}
	if blocked {
		t.Errorf("Did not expect 192.168.1.100 to be blocked")
	}

	// Insertar IP temporalmente para validar
	_, _ = db.Exec("INSERT INTO tIPBlock (sIPStart, sIPEnd, sDesc, dDate) VALUES ('10.0.0.1', '10.0.0.255', 'Test Block', NOW())")
	blocked, err = repo.IsIPBlocked("10.0.0.50")
	if err != nil || !blocked {
		t.Errorf("Expected 10.0.0.50 to be blocked, got blocked=%v, err=%v", blocked, err)
	}
	_, _ = db.Exec("DELETE FROM tIPBlock WHERE sDesc = 'Test Block'")
}

func TestMySQLAccountRepository_AdditionalProcedures(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	// 1. GetUserDateDiffCreated (usp_User_GetDatediffCreated)
	diffCreated, err := repo.GetUserDateDiffCreated(1)
	if err != nil {
		t.Fatalf("GetUserDateDiffCreated failed: %v", err)
	}
	if diffCreated < 0 {
		t.Errorf("Expected non-negative seconds for created diff, got %d", diffCreated)
	}

	// 2. ReturnUser Procedures (usp_ReturnUser_SetUserItem & usp_ReturnUser_GetDatediffLastlog)
	_, _ = db.Exec("DELETE FROM tReturnUser WHERE nUserNo = 1")
	diffBefore, err := repo.GetReturnUserDateDiffLastLog(1)
	if err != nil || diffBefore != -1 {
		t.Errorf("Expected -1 for user without return log, got %d, err=%v", diffBefore, err)
	}

	err = repo.SetReturnUserItem(1, 30001)
	if err != nil {
		t.Fatalf("SetReturnUserItem failed: %v", err)
	}
	diffAfter, err := repo.GetReturnUserDateDiffLastLog(1)
	if err != nil || diffAfter < 0 {
		t.Errorf("Expected non-negative diff after SetReturnUserItem, got %d, err=%v", diffAfter, err)
	}
	_, _ = db.Exec("DELETE FROM tReturnUser WHERE nUserNo = 1")

	// 3. Cash Procedures (usp_Cash_setCash / usp_Event_addCash & usp_Cash_getCash)
	_, _ = db.Exec("DELETE FROM tCash WHERE userNo = 1")
	_, _ = db.Exec("DELETE FROM tCashLog WHERE userNo = 1")

	err = repo.AddCash(1, 0, 500, "EventTest")
	if err != nil {
		t.Fatalf("AddCash failed: %v", err)
	}
	cash, err := repo.GetCash(1, 0)
	if err != nil || cash != 500 {
		t.Errorf("Expected cash=500, got %d, err=%v", cash, err)
	}

	// Add more cash
	err = repo.AddCash(1, 0, 250, "EventTest2")
	if err != nil {
		t.Fatalf("AddCash incremental failed: %v", err)
	}
	cash, err = repo.GetCash(1, 0)
	if err != nil || cash != 750 {
		t.Errorf("Expected cash=750, got %d, err=%v", cash, err)
	}
	_, _ = db.Exec("DELETE FROM tCash WHERE userNo = 1")
	_, _ = db.Exec("DELETE FROM tCashLog WHERE userNo = 1")

	// 4. SerialKey (usp_Item_setSerialKey)
	_, _ = db.Exec("DELETE FROM tSerialKey WHERE userNo = 1")
	err = repo.SetSerialKey(1, 9876543210)
	if err != nil {
		t.Fatalf("SetSerialKey failed: %v", err)
	}
	var readSerial int64
	err = db.QueryRow("SELECT serialKey FROM tSerialKey WHERE userNo = 1").Scan(&readSerial)
	if err != nil || readSerial != 9876543210 {
		t.Errorf("Expected serialKey 9876543210, got %d, err=%v", readSerial, err)
	}
	_, _ = db.Exec("DELETE FROM tSerialKey WHERE userNo = 1")

	// 5. Logout Procedures (usp_Log_setLogout & usp_AccountLog_GetDateDiffLastlogout)
	_, _ = db.Exec("DELETE FROM tLogout WHERE userNo = 1")
	err = repo.LogLogout(0, 1, 100, 1, 20, 15000, 45)
	if err != nil {
		t.Fatalf("LogLogout failed: %v", err)
	}
	diffLogout, err := repo.GetDateDiffLastLogout(1)
	if err != nil || diffLogout < 0 {
		t.Errorf("Expected valid diff for last logout, got %d, err=%v", diffLogout, err)
	}
	_, _ = db.Exec("DELETE FROM tLogout WHERE userNo = 1")

	// 6. Report Procedures (p_Report_Add, p_Report_Get, p_Report_Del)
	_, _ = db.Exec("DELETE FROM tReport WHERE nReporter = 1")
	report := &database.Report{
		Reporter:       1,
		ReporterCharID: "TestReporter",
		Harmer:         "BadPlayer",
		WorldNo:        0,
		ReportType:     "Spam",
		Reason:         "Spamming town chat",
		ChatLog:        "[Chat] spam spam spam",
	}
	regNo, err := repo.AddReport(report)
	if err != nil || regNo <= 0 {
		t.Fatalf("AddReport failed, regNo=%d, err=%v", regNo, err)
	}

	// Test 1 hour cooldown against same reporter
	dupReg, err := repo.AddReport(report)
	if err != nil || dupReg != -1 {
		t.Errorf("Expected -1 (1h cooldown), got %d, err=%v", dupReg, err)
	}

	// Test GetReports
	reports, err := repo.GetReports(1, 0)
	if err != nil || len(reports) != 1 || reports[0].RegNumber != regNo {
		t.Fatalf("GetReports failed or mismatched: len=%d, err=%v", len(reports), err)
	}

	// Test DeleteReport
	delCode, err := repo.DeleteReport(regNo)
	if err != nil || delCode != 0 {
		t.Errorf("Expected delCode 0, got %d, err=%v", delCode, err)
	}

	// Delete again should return 1 (not found)
	delCodeNotFound, err := repo.DeleteReport(regNo)
	if err != nil || delCodeNotFound != 1 {
		t.Errorf("Expected delCode 1 (not found), got %d, err=%v", delCodeNotFound, err)
	}
	_, _ = db.Exec("DELETE FROM tReport WHERE nReporter = 1")

	// 7. SameTimeChar Telemetry (p_SameTimeChar_Add & p_SameTimeCharClass_Add)
	stcID, err := repo.AddSameTimeChar(0, 42)
	if err != nil || stcID <= 0 {
		t.Fatalf("AddSameTimeChar failed: id=%d, err=%v", stcID, err)
	}
	_, _ = db.Exec("DELETE FROM tSameTimeChar WHERE nIndex = ?", stcID)

	err = repo.AddSameTimeCharClass(0, 1, 10, 5)
	if err != nil {
		t.Fatalf("AddSameTimeCharClass failed: %v", err)
	}
	_, _ = db.Exec("DELETE FROM tDayClassLevel WHERE nWorld = 0 AND nClass = 1 AND nLevel = 10")
}

func md5HashString(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
