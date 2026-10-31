package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	envContent := `
LOGIN_HOST=127.0.0.1
LOGIN_PORT=9010

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=xelion_online

ZONE_IP=192.168.0.4
ZONE_PORT=9210
`
	tmpFile := "test.env"
	if err := os.WriteFile(tmpFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to write test.env: %v", err)
	}
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LoginHost != "127.0.0.1" || cfg.LoginPort != 9010 {
		t.Errorf("Unexpected login host/port: %s:%d", cfg.LoginHost, cfg.LoginPort)
	}

	if cfg.DBPassword != "secret" {
		t.Errorf("Unexpected db password: %s", cfg.DBPassword)
	}

	if cfg.DBName != "xelion_online" {
		t.Errorf("Unexpected db name: %s", cfg.DBName)
	}

	if cfg.ZoneIP != "192.168.0.4" || cfg.ZonePort != 9210 {
		t.Errorf("Unexpected zone IP/port: %s:%d", cfg.ZoneIP, cfg.ZonePort)
	}
}
