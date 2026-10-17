package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	envContent := `
LOGIN_HOST=127.0.0.1
LOGIN_PORT=9010
CLIENT_VERSION=10022024000000

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=xelion_online

WORLD_0_NAME=Xelion
WORLD_0_IP=127.0.0.1
WORLD_0_PORT=9014

WORLD_1_NAME=Bijou
WORLD_1_IP=192.168.1.150
WORLD_1_PORT=9015
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

	if cfg.ClientVersion != "10022024000000" {
		t.Errorf("Unexpected client version: %s", cfg.ClientVersion)
	}

	if cfg.DBPassword != "secret" {
		t.Errorf("Unexpected db password: %s", cfg.DBPassword)
	}

	if cfg.DBName != "xelion_online" {
		t.Errorf("Unexpected db name: %s", cfg.DBName)
	}

	if len(cfg.Worlds) != 2 {
		t.Fatalf("Expected 2 worlds, got %d", len(cfg.Worlds))
	}

	if cfg.Worlds[0].Name != "Xelion" || cfg.Worlds[0].Port != 9014 {
		t.Errorf("Unexpected world 0: %+v", cfg.Worlds[0])
	}

	if cfg.Worlds[1].Name != "Bijou" || cfg.Worlds[1].IP != "192.168.1.150" || cfg.Worlds[1].Port != 9015 {
		t.Errorf("Unexpected world 1: %+v", cfg.Worlds[1])
	}
}
