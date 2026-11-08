package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("[WorldService] Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[WorldService] Failed to ping MySQL at %s:%d: %v", cfg.DBHost, cfg.DBPort, err)
	}
	log.Printf("[WorldService] Connected to MySQL database '%s' at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	charRepo := database.NewMySQLCharacterRepository(db)
	accountRepo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)
	gameLogRepo := database.NewMySQLGameLogRepository(db)
	gameLogDispatcher := database.NewGameLogDispatcher(gameLogRepo, 10000, 4)
	gameLogDispatcher.Start()
	log.Println("[WorldService] GameLogDispatcher started (4 async workers)")

	// Cargar todos los archivos SHN y configuraciones estáticas oficiales (22 WorldManager + 4 Character + Checksums)
	dataStore, err := world.LoadWorldDataStore("data")
	if err != nil {
		log.Fatalf("[WorldService] Failed to load WorldDataStore: %v", err)
	}
	log.Printf("[WorldService] WorldDataStore loaded successfully (%d SHN files loaded in memory)", len(dataStore.LoadedFiles))

	dbWorlds, err := accountRepo.GetWorldServerList()
	if err != nil || len(dbWorlds) == 0 {
		log.Fatalf("[WorldService] No se encontraron mundos configurados en tWorldServer (MySQL): %v", err)
	}

	zoneIP := cfg.ZoneIP
	if zoneIP == "" {
		zoneIP = "127.0.0.1"
	}
	zonePort := cfg.ZonePort
	if zonePort == 0 {
		zonePort = 9210
	}

	var worldServers []*world.Server
	for _, w := range dbWorlds {
		worldCfg := &world.Config{
			ID:       w.WorldNo,
			Name:     w.WorldName,
			BindAddr: fmt.Sprintf("0.0.0.0:%d", w.WorldPort),
			Port:     int(w.WorldPort),
			ZoneIP:   zoneIP,
			ZonePort: zonePort,
		}

		srv := world.NewServer(worldCfg, charRepo, transfers, dataStore, accountRepo)
		srv.SetGameLogDispatcher(gameLogDispatcher)
		if err := srv.Start(); err != nil {
			log.Fatalf("[WorldService] Fatal error starting world '%s': %v", w.WorldName, err)
		}
		worldServers = append(worldServers, srv)
	}

	// Telemetría periódica oficial y actualización de CCU real en tWorldServer / tSameTimeChar
	ccuTicker := time.NewTicker(30 * time.Second)
	defer ccuTicker.Stop()
	go func() {
		for range ccuTicker.C {
			for _, ws := range worldServers {
				ccu := ws.Handler().ActivePlayersCount()
				worldID := ws.Config().ID
				_ = accountRepo.UpdateWorldServerCCU(worldID, ccu)
				_, _ = accountRepo.AddSameTimeChar(worldID, ccu)
				_ = accountRepo.AddSyslog(0xC811, 0, "", "", worldID, ccu, 0)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[WorldService] Shutting down...")
	for _, ws := range worldServers {
		ws.Stop()
	}

	drainCtx, drainCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer drainCancel()
	if err := gameLogDispatcher.Close(drainCtx); err != nil {
		log.Printf("[WorldService] GameLogDispatcher drain timeout: %v", err)
	} else {
		log.Println("[WorldService] GameLogDispatcher flushed and stopped cleanly.")
	}

	log.Println("[WorldService] Shutdown complete.")
}
