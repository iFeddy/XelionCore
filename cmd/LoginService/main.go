package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/logger"
	"github.com/iFeddy/XelionCore/pkg/server/login"
)

func main() {
	logFile, err := logger.InitSessionLogger("LoginService")
	if err != nil {
		log.Printf("[LoginService] Warning: failed to initialize session log file: %v", err)
	} else {
		defer logFile.Close()
	}

	log.Println("==================================================")
	log.Println("           XELION ONLINE - LOGIN SERVICE          ")
	log.Println("==================================================")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[LoginService] Failed to load configuration: %v", err)
	}

	log.Printf("[LoginService] Binding on %s:%d", cfg.LoginHost, cfg.LoginPort)

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("[LoginService] Failed to open MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[LoginService] Failed to ping MySQL at %s:%d: %v", cfg.DBHost, cfg.DBPort, err)
	}
	log.Printf("[LoginService] Connected to MySQL database '%s' at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	repo := database.NewMySQLAccountRepository(db)
	transfers := database.NewTransferManager(db)

	// Validar versión obligatoria desde MySQL (FOUS.GameVersionKey)
	vKey, err := repo.GetClientVersionKey()
	if err != nil || strings.TrimSpace(vKey) == "" {
		log.Fatalf("[LoginService] FOUS.GameVersionKey es obligatoria en MySQL y no está disponible: %v", err)
	}
	log.Printf("[LoginService] Client Version Key activa en MySQL (FOUS): '%s'", strings.TrimSpace(vKey))

	// WorldManager gobernado 100% desde tWorldServer en MySQL
	wm := login.NewWorldManager(repo)
	worlds, err := repo.GetWorldServerList()
	if err != nil {
		log.Fatalf("[LoginService] Error leyendo tWorldServer desde MySQL: %v", err)
	}
	for _, w := range worlds {
		log.Printf("[LoginService] Mundo en tWorldServer -> ID %d: '%s' (%s:%d, Maintenance=%v)",
			w.WorldNo, w.WorldName, w.WorldIP, w.WorldPort, w.IsMaintenance)
	}

	loginCfg := &login.Config{
		BindAddr: fmt.Sprintf("%s:%d", cfg.LoginHost, cfg.LoginPort),
		Port:     int(cfg.LoginPort),
	}

	srv := login.NewServer(loginCfg, repo, wm, transfers)
	if err := srv.Start(); err != nil {
		log.Fatalf("[LoginService] Fatal error starting server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n[LoginService] Shutting down...")
	srv.Stop()
	log.Println("[LoginService] Shutdown complete.")
}
