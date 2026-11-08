package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("[Zone] Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[Zone] Failed to ping MySQL at %s:%d: %v", cfg.DBHost, cfg.DBPort, err)
	}
	log.Printf("[Zone] Connected to MySQL database '%s' at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	zonePort := 9210
	if cfg.ZonePort > 0 {
		zonePort = int(cfg.ZonePort)
	}
	zoneIP := "127.0.0.1"
	if cfg.ZoneIP != "" {
		zoneIP = cfg.ZoneIP
	}

	worldNo := byte(0)
	if envW := os.Getenv("ZONE_WORLD_NO"); envW != "" {
		if w, err := strconv.Atoi(envW); err == nil {
			worldNo = byte(w)
		}
	}
	zoneNo := uint16(1)
	if envZ := os.Getenv("ZONE_NO"); envZ != "" {
		if z, err := strconv.Atoi(envZ); err == nil {
			zoneNo = uint16(z)
		}
	}

	zoneCfg := &zone.Config{
		BindAddr: fmt.Sprintf("0.0.0.0:%d", zonePort),
		Port:     zonePort,
		IP:       zoneIP,
		MapID:    0,
		WorldNo:  worldNo,
		ZoneNo:   zoneNo,
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

	// Registrar nodo en tWorldZoneNode mediante CALL usp_WorldZone_UpsertNode oficial
	if err := charRepo.UpsertZoneNode(worldNo, zoneNo, zoneIP, uint16(zonePort)); err != nil {
		log.Printf("[Zone] Warning: failed to register zone node: %v", err)
	} else {
		log.Printf("[Zone] Node registered in tWorldZoneNode (World: %d, Zone: %d, IP: %s, Port: %d)",
			worldNo, zoneNo, zoneIP, zonePort)
	}
	if worldNo == 0 {
		_ = charRepo.UpsertZoneNode(1, zoneNo, zoneIP, uint16(zonePort))
	}

	// Latido periódico de vida en tWorldZoneNode (cada 30 segundos)
	zoneHeartbeatTicker := time.NewTicker(30 * time.Second)
	defer zoneHeartbeatTicker.Stop()

	// 4. Iniciar servidor e inyectar almacén global en memoria RAM
	srv := zone.NewServer(zoneCfg, charRepo, transfers, classStore, fieldStore, chkEngine)
	srv.SetAccountRepo(accountRepo)
	srv.SetDataStore(dataStore)

	go func() {
		for range zoneHeartbeatTicker.C {
			count := srv.Handler().ActivePlayersCount()
			_ = charRepo.HeartbeatZoneNode(worldNo, zoneNo, count)
			if worldNo == 0 {
				_ = charRepo.HeartbeatZoneNode(1, zoneNo, count)
			}
		}
	}()

	// 6. Activar telemetría periódica oficial (p_SameTimeChar_Add y p_SameTimeCharClass_Add en tSameTimeChar)
	telemetryTicker := time.NewTicker(60 * time.Second)
	defer telemetryTicker.Stop()
	go func() {
		for range telemetryTicker.C {
			count := srv.Handler().ActivePlayersCount()
			_, _ = accountRepo.AddSameTimeChar(0, count)

			players := srv.Handler().ActivePlayersSnapshot()
			type classLevelKey struct {
				worldNo byte
				job     byte
				level   byte
			}
			classLevelCounts := make(map[classLevelKey]int)
			for _, p := range players {
				if p != nil && p.Char != nil {
					k := classLevelKey{
						worldNo: p.WorldNo,
						job:     p.Char.Job,
						level:   p.Char.Level,
					}
					classLevelCounts[k]++
				}
			}
			for k, c := range classLevelCounts {
				_ = accountRepo.AddSameTimeCharClass(k.worldNo, k.job, uint16(k.level), c)
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
