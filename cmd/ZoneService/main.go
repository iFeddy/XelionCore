package main

import (
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
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

func main() {
	logFile, err := logger.InitSessionLogger("ZoneService")
	if err != nil {
		log.Printf("[ZoneService] Warning: failed to initialize session log file: %v", err)
	} else {
		defer logFile.Close()
	}

	log.Println("==================================================")
	log.Println("           XELION ONLINE - ZONE SERVER            ")
	log.Println("==================================================")

	// 1. Cargar configuración de entorno (.env)
	cfg, err := config.Load()
	if err != nil {
		log.Printf("[Zone] Warning loading config: %v, using defaults", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[Zone] Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[Zone] Failed to ping MySQL at %s:%d: %v", cfg.DBHost, cfg.DBPort, err)
	}
	log.Printf("[Zone] Connected to MySQL database '%s' at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	zonePort := 9210
	zoneIP := "127.0.0.1"
	if len(cfg.Worlds) > 0 {
		if cfg.Worlds[0].ZonePort > 0 {
			zonePort = int(cfg.Worlds[0].ZonePort)
		}
		if cfg.Worlds[0].ZoneIP != "" {
			zoneIP = cfg.Worlds[0].ZoneIP
		}
	}

	zoneCfg := &zone.Config{
		BindAddr: fmt.Sprintf("0.0.0.0:%d", zonePort),
		Port:     zonePort,
		IP:       zoneIP,
		MapID:    0,
	}

	// 2. Cargar la totalidad de subsistemas y datos en memoria RAM (Zero-Disk-I/O durante runtime)
	log.Println("[Zone] Loading complete official game data into RAM (ClassParam, Field, NPC, MobRegen, Items, Combat, Skills)...")
	dataStore, err := zone.LoadAllData("data")
	if err != nil {
		log.Fatalf("[Zone] Fatal error loading data store: %v", err)
	}

	classStore := dataStore.ClassStore
	fieldStore := dataStore.FieldStore
	chkEngine := dataStore.ChkEngine

	// 3. Inicializar repositorios y gestor de transferencias
	charRepo := database.NewMySQLCharacterRepository(db)
	accountRepo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)

	// 4. Iniciar servidor e inyectar almacén global en memoria RAM
	srv := zone.NewServer(zoneCfg, charRepo, transfers, classStore, fieldStore, chkEngine)
	srv.SetAccountRepo(accountRepo)
	srv.SetDataStore(dataStore)

	// 6. Activar telemetría periódica oficial (p_SameTimeChar_Add y p_SameTimeCharClass_Add en tSameTimeChar)
	telemetryTicker := time.NewTicker(60 * time.Second)
	defer telemetryTicker.Stop()
	go func() {
		for range telemetryTicker.C {
			count := srv.Handler().ActivePlayersCount()
			_, _ = accountRepo.AddSameTimeChar(0, count)

			players := srv.Handler().ActivePlayersSnapshot()
			classCounts := make(map[byte]int)
			for _, p := range players {
				if p != nil && p.Char != nil {
					classCounts[p.Char.Job]++
				}
			}
			for classVal, c := range classCounts {
				_ = accountRepo.AddSameTimeCharClass(0, classVal, 1, c)
			}
		}
	}()
	if err := srv.Start(); err != nil {
		log.Fatalf("[Zone] Fatal error starting server: %v", err)
	}

	log.Printf("[Zone] Server listening on %s (IP: %s)", zoneCfg.BindAddr, zoneCfg.IP)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[Zone] Shutting down...")
	srv.Stop()
	log.Println("[Zone] Shutdown complete.")
}
