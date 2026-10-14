package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
)

// AppConfig represents configuration for the Login daemon.
type AppConfig struct {
	LoginServer login.Config `json:"login_server"`
	Worlds      []WorldEntry `json:"worlds"`
}

type WorldEntry struct {
	ID     byte   `json:"id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Status uint8  `json:"status"`
}

func main() {
	log.Println("==================================================")
	log.Println("           XELION ONLINE - CORE SERVER            ")
	log.Println("           Modern Go Implementation v0.1          ")
	log.Println("==================================================")

	cfg := loadConfig("configs/config.json")

	// Initialize Account Repository (In-Memory for zero-dependency standalone start)
	repo := database.NewMemoryAccountRepository()
	log.Println("[Database] Using high-performance MemoryAccountRepository (Pre-seeded with 'admin'/'admin')")

	// Initialize WorldManager
	wm := login.NewWorldManager()
	if len(cfg.Worlds) > 0 {
		for _, w := range cfg.Worlds {
			wm.AddWorld(&login.World{
				ID:     w.ID,
				Name:   w.Name,
				Status: protocol.WorldStatus(w.Status),
				IP:     w.IP,
				Port:   w.Port,
			})
		}
	}

	// Create and Start Login Server
	srv := login.NewServer(&cfg.LoginServer, repo, wm)
	if err := srv.Start(); err != nil {
		log.Fatalf("[Login] Fatal error starting server: %v", err)
	}

	// Wait for OS shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[Login] Interrupted, initiating graceful shutdown...")
	srv.Stop()
	log.Println("[Login] Shutdown complete. Goodbye!")
}

func loadConfig(path string) *AppConfig {
	defaultCfg := &AppConfig{
		LoginServer: login.Config{
			BindAddr: "0.0.0.0:9010",
			Port:     9010,
		},
		Worlds: []WorldEntry{
			{
				ID:     0,
				Name:   "INITIO",
				IP:     "127.0.0.1",
				Port:   9110,
				Status: uint8(protocol.WorldStatus_Low),
			},
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[Config] Could not read %s (using defaults): %v", path, err)
		return defaultCfg
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("[Config] Failed to parse %s (using defaults): %v", path, err)
		return defaultCfg
	}

	return &cfg
}
