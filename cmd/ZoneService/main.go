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
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

func main() {
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

	// 2. Cargar datos del juego modernizados
	log.Println("[Zone] Loading ClassParam and Field definitions...")
	classStore, err := classparam.Load("data/World/ClassParam")
	if err != nil {
		log.Printf("[Zone] Warning: ClassParam load error: %v", err)
	} else {
		log.Printf("[Zone] Loaded %d character classes into memory.", classStore.ClassCount())
	}

	fieldStore, err := field.Load("data/World/Map")
	if err != nil {
		log.Printf("[Zone] Warning: Field load error: %v", err)
	} else {
		log.Printf("[Zone] Loaded %d field definitions.", fieldStore.FieldCount())
	}

	// 3. Inicializar repositorios y gestor de transferencias
	charRepo := database.NewMySQLCharacterRepository(db)
	transfers := database.NewTransferManager(db)

	// 4. Iniciar servidor
	srv := zone.NewServer(zoneCfg, charRepo, transfers, classStore, fieldStore)
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
