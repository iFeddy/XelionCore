package login

import (
	"fmt"
	"log"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
)

// Config contains runtime settings for the Login Server.
type Config struct {
	BindAddr      string `json:"bind_addr"` // e.g. "0.0.0.0:9010"
	Port          int    `json:"port"`      // e.g. 9010
	ClientVersion string `json:"client_version"`
}

// Server encapsulates the Login Server components.
type Server struct {
	config       *Config
	tcpServer    *network.TCPServer
	handler      *Handler
	worldManager *WorldManager
	repo         database.AccountRepository
}

// NewServer creates a new instance of the Login Server.
func NewServer(cfg *Config, repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager) *Server {
	if cfg.BindAddr == "" {
		cfg.BindAddr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	}

	handler := NewHandler(repo, wm, tm)
	if cfg.ClientVersion != "" {
		handler.SetClientVersion(cfg.ClientVersion)
	}

	s := &Server{
		config:       cfg,
		handler:      handler,
		worldManager: wm,
		repo:         repo,
	}

	s.tcpServer = network.NewTCPServer(cfg.BindAddr, nil, handler.HandlePacket)

	return s
}

// Start launches the Login TCP listener.
func (s *Server) Start() error {
	log.Printf("[LoginServer] Starting XelionCore Login Server on %s...", s.config.BindAddr)
	return s.tcpServer.Start()
}

// Stop shuts down the Login Server.
func (s *Server) Stop() {
	log.Printf("[LoginServer] Shutting down...")
	s.tcpServer.Stop()
}
