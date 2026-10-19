package zone

import (
	"fmt"
	"log"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
)

// Server represents a Fiesta Zone Server instance.
type Server struct {
	config    *Config
	tcpServer *network.TCPServer
	handler   *Handler
}

// NewServer initializes a new Zone Server.
func NewServer(cfg *Config, charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.Store, chkEngine *shn.ChecksumEngine) *Server {
	if cfg.BindAddr == "" {
		cfg.BindAddr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	}

	handler := NewHandler(charRepo, transfers, classStore, fieldStore, chkEngine)

	s := &Server{
		config:  cfg,
		handler: handler,
	}

	onConnect := func(client *network.Client) {
		log.Printf("[ZoneServer] Client connected from %s (XOR Position: %d)",
			client.RemoteAddr(), client.XorOffset())
	}

	s.tcpServer = network.NewTCPServer(cfg.BindAddr, onConnect, handler.HandlePacket)
	return s
}

// Start begins accepting Zone client connections.
func (s *Server) Start() error {
	log.Printf("[ZoneServer] Starting Zone Server on %s (Dynamic Map Engine active)...",
		s.config.BindAddr)
	return s.tcpServer.Start()
}

// Stop gracefully shuts down the Zone Server.
func (s *Server) Stop() {
	log.Printf("[ZoneServer] Shutting down Zone Server...")
	s.tcpServer.Stop()
}
