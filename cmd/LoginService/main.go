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

type Config struct {
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
	log.Println("           XELION ONLINE - LOGIN SERVER           ")
	log.Println("                  Port: 9010                      ")
	log.Println("==================================================")

	cfg := loadConfig("configs/config.json")

	repo := database.NewMemoryAccountRepository()
	transfers := database.NewTransferManager()

	wm := login.NewWorldManager()
	for _, w := range cfg.Worlds {
		wm.AddWorld(&login.World{
			ID:     w.ID,
			Name:   w.Name,
			Status: protocol.WorldStatus(w.Status),
			IP:     w.IP,
			Port:   w.Port,
		})
	}

	srv := login.NewServer(&cfg.LoginServer, repo, wm, transfers)
	if err := srv.Start(); err != nil {
		log.Fatalf("[Login] Fatal error starting server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[Login] Shutting down...")
	srv.Stop()
	log.Println("[Login] Shutdown complete.")
}

func loadConfig(path string) *Config {
	defaultCfg := &Config{
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
				Status: 6,
			},
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
