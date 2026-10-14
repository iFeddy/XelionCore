package network

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// TCPServer manages incoming TCP client connections for a Xelion service.
type TCPServer struct {
	addr          string
	listener      net.Listener
	clients       sync.Map
	onConnect     func(*Client)
	packetHandler func(*Client, *protocol.Packet)
	wg            sync.WaitGroup
	quit          chan struct{}
}

// NewTCPServer creates a new TCPServer instance.
func NewTCPServer(addr string, onConnect func(*Client), packetHandler func(*Client, *protocol.Packet)) *TCPServer {
	return &TCPServer{
		addr:          addr,
		onConnect:     onConnect,
		packetHandler: packetHandler,
		quit:          make(chan struct{}),
	}
}

// Start begins listening and serving connections.
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	s.listener = listener
	log.Printf("[TCPServer] Listening on %s", s.addr)

	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

func (s *TCPServer) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Printf("[TCPServer] accept error: %v", err)
				continue
			}
		}

		client, err := NewClient(conn, func(c *Client) {
			s.clients.Delete(c.RemoteAddr())
			log.Printf("[%s] client disconnected", c.RemoteAddr())
		})
		if err != nil {
			log.Printf("[TCPServer] failed to initialize client crypto: %v", err)
			_ = conn.Close()
			continue
		}

		s.clients.Store(client.RemoteAddr(), client)
		log.Printf("[%s] client connected, xor offset: %d", client.RemoteAddr(), client.Crypto().XorPos())

		// Send initial Fiesta handshake packet: SH2Type_SetXorKeyPosition
		handshake := protocol.NewPacketWriter(protocol.SH2Type_SetXorKeyPosition)
		handshake.WriteInt16(client.Crypto().XorPos())
		if err := client.SendPacket(handshake); err != nil {
			log.Printf("[%s] failed to send initial handshake: %v", client.RemoteAddr(), err)
			client.Close()
			continue
		}

		if s.onConnect != nil {
			s.onConnect(client)
		}

		go client.StartReadLoop(s.packetHandler)
	}
}

// Stop shuts down the listener and disconnects all clients.
func (s *TCPServer) Stop() {
	close(s.quit)
	if s.listener != nil {
		_ = s.listener.Close()
	}

	s.clients.Range(func(key, value any) bool {
		if client, ok := value.(*Client); ok {
			client.Close()
		}
		return true
	})

	s.wg.Wait()
	log.Printf("[TCPServer] Stopped listener on %s", s.addr)
}
