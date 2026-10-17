package database_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
)

func TestMySQLConnection(t *testing.T) {
	cfg, err := config.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping MySQL at %s: %v", dsn, err)
	}

	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tUser").Scan(&userCount)
	if err != nil {
		t.Fatalf("Failed to query tUser: %v", err)
	}

	if userCount < 2 {
		t.Fatalf("Expected at least 2 seed users in tUser, got %d", userCount)
	}
	t.Logf("Successfully connected to MySQL! Found %d users in tUser.", userCount)
}
