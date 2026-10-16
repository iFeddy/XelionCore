package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

type AppConfig struct {
	LoginServer login.Config `json:"login_server"`
	Worlds      []WorldEntry `json:"worlds"`
	ZoneServer  zone.Config  `json:"zone_server"`
}

type WorldEntry struct {
	ID       byte   `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	Status   uint8  `json:"status"`
	ZoneIP   string `json:"zone_ip"`
	ZonePort uint16 `json:"zone_port"`
}

func main() {
	log.Println("==================================================")
	log.Println("           XELION ONLINE - ALL-IN-ONE             ")
	log.Println("       Login (9010) + World (9110) + Zone (9210)  ")
	log.Println("==================================================")

	cfg := loadConfig("configs/config.json")

	repo := database.NewMemoryAccountRepository()
	charRepo := database.NewMemoryCharacterRepository()
	transfers := database.NewTransferManager()

	// 1. Login Server
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
	loginSrv := login.NewServer(&cfg.LoginServer, repo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		log.Fatalf("[Login] Fatal error starting server: %v", err)
	}

	// 2. World Server
	var worldServers []*world.Server
	for _, w := range cfg.Worlds {
		zoneIP := w.ZoneIP
		if zoneIP == "" {
			zoneIP = "127.0.0.1"
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
		worldSrv := world.NewServer(worldCfg, charRepo, transfers)
		if err := worldSrv.Start(); err != nil {
			log.Fatalf("[World] Fatal error starting world '%s': %v", w.Name, err)
		}
		worldServers = append(worldServers, worldSrv)
	}

	// 3. Zone Server
	zoneSrv := zone.NewServer(&cfg.ZoneServer, charRepo, transfers)
	if err := zoneSrv.Start(); err != nil {
		log.Fatalf("[Zone] Fatal error starting server: %v", err)
	}

	log.Println("[Server] All 3 services (Login, World, Zone) are online!")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[Server] Shutting down all services...")
	loginSrv.Stop()
	for _, ws := range worldServers {
		ws.Stop()
	}
	zoneSrv.Stop()
	log.Println("[Server] Shutdown complete.")
}

func loadConfig(path string) *AppConfig {
	defaultCfg := &AppConfig{
		LoginServer: login.Config{
			BindAddr: "0.0.0.0:9010",
			Port:     9010,
		},
		Worlds: []WorldEntry{
			{
				ID:       0,
				Name:     "INITIO",
				IP:       "127.0.0.1",
				Port:     9110,
				Status:   6,
				ZoneIP:   "127.0.0.1",
				ZonePort: 9210,
			},
		},
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

	var parsed AppConfig
	if err := json.Unmarshal(data, &parsed); err != nil {
		return defaultCfg
	}
	return &parsed
}
