package world

import (
	"fmt"
	"log"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
)

// Config holds runtime options for the World Server.
type Config struct {
	ID       byte   `json:"id"`
	Name     string `json:"name"`
	BindAddr string `json:"bind_addr"`
	Port     int    `json:"port"`
	ZoneIP   string `json:"zone_ip"`
	ZonePort uint16 `json:"zone_port"`
}

// Server represents a Fiesta World Server instance.
type Server struct {
	config    *Config
	tcpServer *network.TCPServer
	handler   *Handler
}

// NewServer initializes a new World Server.
func NewServer(
	cfg *Config,
	charRepo database.CharacterRepository,
	transfers *database.TransferManager,
	dataStore *WorldDataStore,
) *Server {
	if cfg.BindAddr == "" {
		cfg.BindAddr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	}
	if cfg.ZoneIP == "" {
		cfg.ZoneIP = "127.0.0.1"
	}
	if cfg.ZonePort == 0 {
		cfg.ZonePort = 9210
	}

	handler := NewHandler(charRepo, transfers, dataStore, cfg.ZoneIP, cfg.ZonePort)

	s := &Server{
		config:  cfg,
		handler: handler,
	}

	s.tcpServer = network.NewTCPServer(cfg.BindAddr, nil, handler.HandlePacket)
	s.tcpServer.SetOnDisconnect(func(c *network.Client) {
		if c.CharacterID > 0 {
			handler.partyMgr.UnregisterOnline(c.CharacterID)
		}
	})
	return s
}

// Start listens and begins accepting World client connections.
func (s *Server) Start() error {
	log.Printf("[WorldServer] Starting World Server '%s' (ID: %d) on %s (Zone -> %s:%d)...",
		s.config.Name, s.config.ID, s.config.BindAddr, s.config.ZoneIP, s.config.ZonePort)
	return s.tcpServer.Start()
}

// Stop gracefully shuts down the World Server.
func (s *Server) Stop() {
	log.Printf("[WorldServer] Shutting down '%s'...", s.config.Name)
	s.tcpServer.Stop()
}
