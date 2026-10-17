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

func md5HashString(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
