package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/logger"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

func main() {
	logFile, err := logger.InitSessionLogger("WorldService")
	if err != nil {
		log.Printf("[WorldService] Warning: failed to initialize session log file: %v", err)
	} else {
		defer logFile.Close()
	}

	log.Println("==================================================")
	log.Println("           XELION ONLINE - WORLD SERVICE          ")
	log.Println("==================================================")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[WorldService] Failed to load configuration: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[WorldService] Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[WorldService] Failed to ping MySQL at %s:%d: %v", cfg.DBHost, cfg.DBPort, err)
	}
	log.Printf("[WorldService] Connected to MySQL database '%s' at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	charRepo := database.NewMySQLCharacterRepository(db)
	transfers := database.NewTransferManager(db)

	// Cargar todos los archivos SHN y configuraciones estáticas oficiales (22 WorldManager + 4 Character + Checksums)
	dataStore, err := world.LoadWorldDataStore("data")
	if err != nil {
		log.Fatalf("[WorldService] Failed to load WorldDataStore: %v", err)
	}
	log.Printf("[WorldService] WorldDataStore loaded successfully (%d SHN files loaded in memory)", len(dataStore.LoadedFiles))

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

		srv := world.NewServer(worldCfg, charRepo, transfers, dataStore)
		if err := srv.Start(); err != nil {
			log.Fatalf("[WorldService] Fatal error starting world '%s': %v", w.Name, err)
		}
		worldServers = append(worldServers, srv)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[WorldService] Shutting down...")
	for _, ws := range worldServers {
		ws.Stop()
	}
	log.Println("[WorldService] Shutdown complete.")
}
