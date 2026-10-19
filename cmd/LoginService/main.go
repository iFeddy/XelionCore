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
	"github.com/iFeddy/XelionCore/pkg/protocol"
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

	log.Printf("[LoginService] Binding on %s:%d (Client Version: %s)", cfg.LoginHost, cfg.LoginPort, cfg.ClientVersion)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
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

	wm := login.NewWorldManager()
	for _, w := range cfg.Worlds {
		log.Printf("[LoginService] Registering World ID %d: '%s' -> %s:%d", w.ID, w.Name, w.IP, w.Port)
		wm.AddWorld(&login.World{
			ID:     w.ID,
			Name:   w.Name,
			Status: protocol.WorldStatus_Low, // Estado inicial Online / Activo
			IP:     w.IP,
			Port:   w.Port,
		})
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
