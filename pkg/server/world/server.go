package world

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/s2s"
)

// Config holds runtime options for the World Server.
type Config struct {
	ID          byte   `json:"id"`
	Name        string `json:"name"`
	BindAddr    string `json:"bind_addr"`
	Port        int    `json:"port"`
	ZoneIP      string `json:"zone_ip"`
	ZonePort    uint16 `json:"zone_port"`
	S2SBindAddr string `json:"s2s_bind_addr"`
}

// Server represents a Fiesta World Server instance.
type Server struct {
	config    *Config
	tcpServer *network.TCPServer
	handler   *Handler
	s2sServer *s2s.Server
	tickMu    sync.Mutex
	stopTick  chan struct{}
}

// NewServer initializes a new World Server.
func NewServer(
	cfg *Config,
	charRepo database.CharacterRepository,
	transfers *database.TransferManager,
	dataStore *WorldDataStore,
	accountRepo ...database.AccountRepository,
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
	handler.SetWorldID(cfg.ID)
	if len(accountRepo) > 0 && accountRepo[0] != nil {
		handler.SetAccountRepo(accountRepo[0])
	}

	s := &Server{
		config:  cfg,
		handler: handler,
	}

	s.tcpServer = network.NewTCPServer(cfg.BindAddr, nil, handler.HandlePacket)
	s.tcpServer.SetOnDisconnect(func(c *network.Client) {
		handler.OnClientDisconnect(c)
		if c.CharacterID > 0 {
			handler.partyMgr.UnregisterOnline(c.CharacterID)
			handler.friendMgr.BroadcastLogout(c.CharacterID, handler.activeCharName(c))
			handler.guildMgr.UnregisterOnline(c.CharacterID)
			handler.kqMgr.UnregisterOnline(c.CharacterID)
		}
	})
	if cfg.S2SBindAddr != "" {
		s.s2sServer = s2s.NewServer(cfg.S2SBindAddr, handler.HandleZoneEvent)
		handler.SetS2SSink(s.s2sServer)
	}
	return s
}

// SetAccountRepo configura el repositorio oficial de Account/AccountLog/Session.
func (s *Server) SetAccountRepo(repo database.AccountRepository) {
	s.handler.SetAccountRepo(repo)
}

// SetGameLogDispatcher inyecta el despachador de logs de auditoría en WorldService.
func (s *Server) SetGameLogDispatcher(disp *database.GameLogDispatcher) {
	if s.handler != nil {
		s.handler.SetGameLogDispatcher(disp)
	}
}

// Handler retorna el despachador de red de WorldServer.
func (s *Server) Handler() *Handler {
	return s.handler
}

// Config retorna la configuración del servidor de mundo.
func (s *Server) Config() *Config {
	return s.config
}

// S2SAddr retorna la dirección del coordinador S2S vivo si está habilitado.
func (s *Server) S2SAddr() string {
	if s.s2sServer == nil {
		return ""
	}
	return s.s2sServer.Addr()
}

// S2SConnectedZones retorna la cantidad de nodos ZoneService enlazados al coordinador S2S.
func (s *Server) S2SConnectedZones() int {
	if s.s2sServer == nil {
		return 0
	}
	return s.s2sServer.ConnectedZones()
}

// Start listens and begins accepting World client connections.
func (s *Server) Start() error {
	log.Printf("[WorldServer] Starting World Server '%s' (ID: %d) on %s (Zone -> %s:%d)...",
		s.config.Name, s.config.ID, s.config.BindAddr, s.config.ZoneIP, s.config.ZonePort)
	if s.s2sServer != nil {
		if err := s.s2sServer.Start(); err != nil {
			return fmt.Errorf("failed to start S2S coordinator on %s: %w", s.config.S2SBindAddr, err)
		}
		log.Printf("[WorldServer] Live S2S coordinator listening on %s", s.s2sServer.Addr())
	}
	s.startWorldTicker()
	return s.tcpServer.Start()
}

func (s *Server) startWorldTicker() {
	s.tickMu.Lock()
	defer s.tickMu.Unlock()
	if s.stopTick != nil {
		return
	}
	ch := make(chan struct{})
	s.stopTick = ch
	go func() {
		ticker := time.NewTicker(1 * time.Second) // exe: 0x00406769 (WorldManagerServer timer -> DoKingdomQuest)
		defer ticker.Stop()
		for {
			select {
			case <-ch:
				return
			case now := <-ticker.C:
				if s.handler != nil {
					s.handler.RunWorldTick(now)
				}
			}
		}
	}()
}

// Stop gracefully shuts down the World Server.
func (s *Server) Stop() {
	log.Printf("[WorldServer] Shutting down '%s'...", s.config.Name)
	s.tickMu.Lock()
	if s.stopTick != nil {
		close(s.stopTick)
		s.stopTick = nil
	}
	s.tickMu.Unlock()
	if s.s2sServer != nil {
		s.s2sServer.Stop()
	}
	s.tcpServer.Stop()
}
