package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

type Config struct {
	ZoneServer zone.Config `json:"zone_server"`
}

func main() {
	log.Println("==================================================")
	log.Println("           XELION ONLINE - ZONE SERVER            ")
	log.Println("                  Port: 9210                      ")
	log.Println("==================================================")

	cfg := loadConfig("configs/config.json")

	charRepo := database.NewMemoryCharacterRepository()
	transfers := database.NewTransferManager()

	srv := zone.NewServer(&cfg.ZoneServer, charRepo, transfers)
	if err := srv.Start(); err != nil {
		log.Fatalf("[Zone] Fatal error starting server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[Zone] Shutting down...")
	srv.Stop()
	log.Println("[Zone] Shutdown complete.")
}

func loadConfig(path string) *Config {
	defaultCfg := &Config{
		ZoneServer: zone.Config{
			BindAddr: "0.0.0.0:9210",
			Port:     9210,
			IP:       "127.0.0.1",
			MapID:    0,
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultCfg
	}

	var parsed Config
	if err := json.Unmarshal(data, &parsed); err != nil {
		return defaultCfg
	}
	return &parsed
}
