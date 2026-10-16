package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

type Config struct {
	Worlds []WorldEntry `json:"worlds"`
}

type WorldEntry struct {
	ID       byte   `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	ZoneIP   string `json:"zone_ip"`
	ZonePort uint16 `json:"zone_port"`
}

func main() {
	log.Println("==================================================")
	log.Println("      XELION ONLINE - WORLD / CHAR SERVER         ")
	log.Println("                  Port: 9110                      ")
	log.Println("==================================================")

	cfg := loadConfig("configs/config.json")

	charRepo := database.NewMemoryCharacterRepository()
	transfers := database.NewTransferManager()

	var worldServers []*world.Server
	for _, w := range cfg.Worlds {
		zoneIP := w.ZoneIP
		if zoneIP == "" {
			zoneIP = "192.168.0.4"
		}
		zonePort := w.ZonePort
		if zonePort == 0 {
			zonePort = 9210
		}

		worldCfg := &world.Config{
			ID:       w.ID,
			Name:     w.Name,
			BindAddr: fmt.Sprintf("0.0.0.0:%d", w.Port),
			Port:     int(w.Port),
			ZoneIP:   zoneIP,
			ZonePort: zonePort,
		}

		srv := world.NewServer(worldCfg, charRepo, transfers)
		if err := srv.Start(); err != nil {
			log.Fatalf("[World] Fatal error starting world '%s': %v", w.Name, err)
		}
		worldServers = append(worldServers, srv)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[World] Shutting down...")
	for _, ws := range worldServers {
		ws.Stop()
	}
	log.Println("[World] Shutdown complete.")
}

func loadConfig(path string) *Config {
	defaultCfg := &Config{
		Worlds: []WorldEntry{
			{
				ID:       0,
				Name:     "INITIO",
				IP:       "192.168.0.4",
				Port:     9110,
				ZoneIP:   "192.168.0.4",
				ZonePort: 9210,
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
