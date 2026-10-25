package world

import (
	"fmt"
	"log"
	"net"

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
	config      *Config
	tcpServer   *network.TCPServer
	auxListener net.Listener
	handler     *Handler
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
			handler.friendMgr.BroadcastLogout(c.CharacterID, c.Username)
			handler.guildMgr.UnregisterOnline(c.CharacterID)
			handler.kqMgr.UnregisterOnline(c.CharacterID)
		}
	})
	return s
}

// Start listens and begins accepting World client connections.
func (s *Server) Start() error {
	log.Printf("[WorldServer] Starting World Server '%s' (ID: %d) on %s (Zone -> %s:%d)...",
		s.config.Name, s.config.ID, s.config.BindAddr, s.config.ZoneIP, s.config.ZonePort)

	// Iniciar listener para el puerto auxiliar 9350 requerido por el cliente NA2026
	s.startAuxPort9350()

	return s.tcpServer.Start()
}

// Stop gracefully shuts down the World Server.
func (s *Server) Stop() {
	log.Printf("[WorldServer] Shutting down '%s'...", s.config.Name)
	if s.auxListener != nil {
		_ = s.auxListener.Close()
	}
	s.tcpServer.Stop()
}

func (s *Server) startAuxPort9350() {
	ln, err := net.Listen("tcp", "0.0.0.0:9350")
	if err != nil {
		log.Printf("[WorldServer] Port 9350 aux listener not started (may be in use): %v", err)
		return
	}
	s.auxListener = ln
	log.Printf("[WorldServer] Listening for NA2026 aux connection on 0.0.0.0:9350")

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go s.handleAuxClient(conn)
		}
	}()
}

func (s *Server) handleAuxClient(conn net.Conn) {
	defer conn.Close()
	// Handshake 0x0804 capturado en puerto 9350 (wire: 02 04 08)
	handshake := []byte{0x02, 0x04, 0x08}
	if _, err := conn.Write(handshake); err != nil {
		return
	}
	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
	}
}
