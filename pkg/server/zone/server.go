package zone



import (
	"fmt"
	"log"
	"time"



	"github.com/iFeddy/XelionCore/pkg/data/classparam"

	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"

	"github.com/iFeddy/XelionCore/pkg/data/shn"

	"github.com/iFeddy/XelionCore/pkg/database"

	"github.com/iFeddy/XelionCore/pkg/network"

	"github.com/iFeddy/XelionCore/pkg/s2s"

)



// Server represents a Fiesta Zone Server instance.

type liveZoneWMSBridge struct {
	client *s2s.Client
}

func (b *liveZoneWMSBridge) PublishWMS(ev WMSEvent) {
	if b == nil || b.client == nil {
		return
	}
	_ = b.client.Publish(s2s.Event{
		Symbol:  ev.Symbol,
		WorldID: ev.WorldID,
		ZoneNo:  ev.ZoneNo,
		Payload: ev.Payload,
	})
}

type Server struct {
	config    *Config
	tcpServer *network.TCPServer
	handler   *Handler
	s2sClient *s2s.Client
	kqStopCh  chan struct{}
}



// NewServer initializes a new Zone Server.

func NewServer(cfg *Config, charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.Store, chkEngine *shn.ChecksumEngine) *Server {

	if cfg.BindAddr == "" {

		cfg.BindAddr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	}



	handler := NewHandler(charRepo, transfers, classStore, fieldStore, chkEngine)
	if cfg != nil {
		if cfg.ZoneID != "" {
			if zNo, zPort, ok := ResolveZoneNodeConfig(cfg.ZoneID); ok {
				cfg.MultiNode = true
				cfg.ZoneNo = zNo
				if cfg.Port == 0 {
					cfg.Port = zPort
					cfg.BindAddr = fmt.Sprintf("0.0.0.0:%d", cfg.Port)
				}
			}
		}
		handler.SetWorldID(cfg.WorldNo)
		handler.SetZoneNo(cfg.ZoneNo)
		if cfg.MultiNode {
			handler.SetMultiNode(true)
			handler.MapManager().SetZoneFilter(int(cfg.ZoneNo))
		}
	}



	s := &Server{

		config:  cfg,

		handler: handler,

	}

	if cfg != nil && cfg.WorldS2SAddr != "" {
		s2sCli := s2s.NewClient(cfg.WorldS2SAddr, cfg.WorldNo, cfg.ZoneNo, func(ev s2s.Event) {
			handler.ApplyIncomingWMS(ev.Symbol, ev.WorldID, ev.ZoneNo, ev.Payload)
		})
		s.s2sClient = s2sCli
		handler.SetWMSBridge(&liveZoneWMSBridge{client: s2sCli})
	}



	onConnect := func(client *network.Client) {

		log.Printf("[ZoneServer] Client connected from %s (XOR Position: %d)",

			client.RemoteAddr(), client.XorOffset())

	}



	s.tcpServer = network.NewTCPServer(cfg.BindAddr, onConnect, handler.HandlePacket)

	s.tcpServer.SetOnDisconnect(handler.OnClientDisconnect)

	return s

}



// Start begins accepting Zone client connections.

func (s *Server) Start() error {

	log.Printf("[ZoneServer] Starting Zone Server on %s (Dynamic Map Engine active)...",

		s.config.BindAddr)

	if s.s2sClient != nil {
		_ = s.s2sClient.Start()
		log.Printf("[ZoneServer] S2S WMS bridge targeting WorldService at %s (World=%d, Zone=%d)",
			s.config.WorldS2SAddr, s.config.WorldNo, s.config.ZoneNo)
	}

	// Zone.exe 10Hz event tick (100ms) for KingdomQuestGateKeeper waves (exe: 0x00499699)
	s.kqStopCh = make(chan struct{})
	go func(stopCh <-chan struct{}) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case now := <-ticker.C:
				if s.handler != nil {
					s.handler.ProcessKQGateKeepers(now)
				}
			}
		}
	}(s.kqStopCh)

	return s.tcpServer.Start()

}

// S2SConnected indica si el enlace S2S vivo hacia WorldService está conectado.
func (s *Server) S2SConnected() bool {
	return s.s2sClient != nil && s.s2sClient.Connected()
}



// Stop gracefully shuts down the Zone Server.

func (s *Server) Stop() {

	log.Printf("[ZoneServer] Shutting down Zone Server...")

	if s.kqStopCh != nil {
		close(s.kqStopCh)
		s.kqStopCh = nil
	}

	if s.s2sClient != nil {
		s.s2sClient.Stop()
	}

	s.tcpServer.Stop()

	if s.handler != nil {
		if mm := s.handler.MapManager(); mm != nil {
			mm.StopAll()
		}
	}

}



// SetAccountRepo sets the MySQLAccountRepository for session logging and telemetry.

// SetItemStore sets the item data store for the handler.
func (s *Server) SetDataStore(ds *DataStore) {
	s.handler.SetDataStore(ds)
	if s.config != nil && (s.config.MultiNode || s.config.PreloadAll) && s.handler != nil && s.handler.MapManager() != nil {
		loaded := s.handler.MapManager().PreloadZoneMaps()
		if s.config.MultiNode {
			log.Printf("[ZoneServer] Multi-node %s (ZoneNo=%d) preloaded %d owned maps", s.config.ZoneID, s.config.ZoneNo, loaded)
		} else {
			log.Printf("[ZoneServer] Single-node preload enabled (ZONE_PRELOAD_ALL=true): preloaded %d maps", loaded)
		}
	}
}

func (s *Server) SetItemStore(store *item.Store) {
	s.handler.SetItemStore(store)
}

func (s *Server) SetAccountRepo(repo *database.MySQLAccountRepository) {
	s.handler.SetAccountRepo(repo)
}

// SetGameLogDispatcher inyecta el despachador de logs de auditoría en el handler de zona.
func (s *Server) SetGameLogDispatcher(disp *database.GameLogDispatcher) {
	if s.handler != nil {
		s.handler.SetGameLogDispatcher(disp)
	}
}



// Handler returns the zone packet handler.

func (s *Server) Handler() *Handler {

	return s.handler

}

