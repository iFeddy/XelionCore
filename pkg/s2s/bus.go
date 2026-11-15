package s2s

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// Event representa un evento S2S oficial wms_NC_* intercambiado en vivo entre procesos.
type Event struct {
	Symbol  string                 `json:"symbol"`
	WorldID byte                   `json:"world_id"`
	ZoneNo  uint16                 `json:"zone_no"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type peerConn struct {
	conn    net.Conn
	writeMu sync.Mutex
	zoneNo  uint16
	worldID byte
}

func (p *peerConn) send(ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	_ = p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = p.conn.Write(data)
	return err
}

// Server coordina en WorldService las conexiones S2S vivas desde nodos ZoneService.
type Server struct {
	bindAddr string
	ln       net.Listener
	onEvent  func(zoneNo uint16, ev Event)

	mu     sync.RWMutex
	peers  map[net.Conn]*peerConn
	closed bool
	wg     sync.WaitGroup
}

// NewServer crea un servidor coordinador S2S para WorldService.
func NewServer(bindAddr string, onEvent func(zoneNo uint16, ev Event)) *Server {
	return &Server{
		bindAddr: bindAddr,
		onEvent:  onEvent,
		peers:    make(map[net.Conn]*peerConn),
	}
}

// Start inicia el listener TCP S2S.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.ln = ln
	s.closed = false
	s.mu.Unlock()

	s.wg.Add(1)
	go s.acceptLoop(ln)
	return nil
}

// Addr devuelve la dirección efectiva escuchada (útil cuando bindAddr usa puerto :0).
func (s *Server) Addr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ln == nil {
		return s.bindAddr
	}
	return s.ln.Addr().String()
}

// Stop cierra el listener y todas las conexiones de nodos ZoneService.
func (s *Server) Stop() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	if s.ln != nil {
		_ = s.ln.Close()
	}
	for c := range s.peers {
		_ = c.Close()
	}
	s.peers = make(map[net.Conn]*peerConn)
	s.mu.Unlock()
	s.wg.Wait()
}

// ConnectedZones devuelve cuántos nodos ZoneService mantienen enlace activo.
func (s *Server) ConnectedZones() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peers)
}

// Broadcast envía un evento wms_NC_* a todos los nodos ZoneService conectados.
func (s *Server) Broadcast(ev Event) {
	s.BroadcastExceptZone(ev, 0xFFFF)
}

// BroadcastExceptZone envía un evento wms_NC_* a todos los nodos excepto exceptZoneNo.
func (s *Server) BroadcastExceptZone(ev Event, exceptZoneNo uint16) {
	s.mu.RLock()
	list := make([]*peerConn, 0, len(s.peers))
	for _, p := range s.peers {
		if p.zoneNo != exceptZoneNo {
			list = append(list, p)
		}
	}
	s.mu.RUnlock()

	for _, p := range list {
		if err := p.send(ev); err != nil {
			_ = p.conn.Close()
		}
	}
}

// SendToZone envía un evento puntual a un nodo ZoneService específico.
func (s *Server) SendToZone(zoneNo uint16, ev Event) bool {
	s.mu.RLock()
	var targets []*peerConn
	for _, p := range s.peers {
		if p.zoneNo == zoneNo {
			targets = append(targets, p)
		}
	}
	s.mu.RUnlock()
	sent := false
	for _, p := range targets {
		if err := p.send(ev); err == nil {
			sent = true
		}
	}
	return sent
}

func (s *Server) acceptLoop(ln net.Listener) {
	defer s.wg.Done()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		pc := &peerConn{conn: conn}
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			_ = conn.Close()
			return
		}
		s.peers[conn] = pc
		s.mu.Unlock()

		s.wg.Add(1)
		go s.readPeer(pc)
	}
}

func (s *Server) readPeer(pc *peerConn) {
	defer s.wg.Done()
	defer func() {
		s.mu.Lock()
		delete(s.peers, pc.conn)
		s.mu.Unlock()
		_ = pc.conn.Close()
	}()

	scanner := bufio.NewScanner(pc.conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var ev Event
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}
		if ev.ZoneNo > 0 {
			pc.zoneNo = ev.ZoneNo
			pc.worldID = ev.WorldID
		}
		if ev.Symbol == "wms_NC_MISC_S2SCONNECTION_RDY" {
			_ = pc.send(Event{
				Symbol:  "wms_NC_MISC_S2SCONNECTION_ACK",
				WorldID: ev.WorldID,
				ZoneNo:  ev.ZoneNo,
				Payload: map[string]interface{}{"status": "ready"},
			})
		}
		if s.onEvent != nil {
			s.onEvent(pc.zoneNo, ev)
		}
	}
}

// Client conecta un proceso ZoneService con el coordinador S2S de WorldService.
type Client struct {
	worldAddr string
	worldID   byte
	zoneNo    uint16
	onEvent   func(ev Event)

	mu      sync.Mutex
	conn    net.Conn
	writeMu sync.Mutex
	stopCh  chan struct{}
	stopped bool
	wg      sync.WaitGroup
}

// NewClient crea un cliente S2S persistente para ZoneService con reconexión automática.
func NewClient(worldAddr string, worldID byte, zoneNo uint16, onEvent func(ev Event)) *Client {
	return &Client{
		worldAddr: worldAddr,
		worldID:   worldID,
		zoneNo:    zoneNo,
		onEvent:   onEvent,
		stopCh:    make(chan struct{}),
	}
}

// Start inicia la conexión inicial (si World ya está escuchando) y el loop en segundo plano.
func (c *Client) Start() error {
	_ = c.tryDial()
	c.wg.Add(1)
	go c.maintainLoop()
	return nil
}

// Connected indica si el enlace TCP hacia WorldService está activo.
func (c *Client) Connected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn != nil
}

// Publish envía un evento wms_NC_* hacia WorldService.
func (c *Client) Publish(ev Event) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		if err := c.tryDial(); err != nil {
			return err
		}
		c.mu.Lock()
		conn = c.conn
		c.mu.Unlock()
		if conn == nil {
			return errors.New("s2s: world connection unavailable")
		}
	}
	return c.writeConn(conn, ev)
}

func (c *Client) writeConn(conn net.Conn, ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write(data)
	if err != nil {
		c.clearConn(conn)
	}
	return err
}

// Stop cierra limpiamente el cliente S2S.
func (c *Client) Stop() {
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		return
	}
	c.stopped = true
	close(c.stopCh)
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()
	c.wg.Wait()
}

func (c *Client) clearConn(target net.Conn) {
	c.mu.Lock()
	if c.conn == target {
		c.conn = nil
	}
	c.mu.Unlock()
	_ = target.Close()
}

func (c *Client) tryDial() error {
	if c.worldAddr == "" {
		return errors.New("s2s: empty world address")
	}
	conn, err := net.DialTimeout("tcp", c.worldAddr, 1500*time.Millisecond)
	if err != nil {
		return err
	}
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		_ = conn.Close()
		return errors.New("s2s: client stopped")
	}
	old := c.conn
	c.conn = conn
	c.mu.Unlock()
	if old != nil {
		_ = old.Close()
	}

	_ = c.writeConn(conn, Event{
		Symbol:  "wms_NC_MISC_S2SCONNECTION_RDY",
		WorldID: c.worldID,
		ZoneNo:  c.zoneNo,
		Payload: map[string]interface{}{"zoneNo": c.zoneNo, "worldID": c.worldID},
	})

	c.wg.Add(1)
	go c.readLoop(conn)
	return nil
}

func (c *Client) maintainLoop() {
	defer c.wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			if !c.Connected() {
				if err := c.tryDial(); err != nil {
					log.Printf("[Zone][S2S] Waiting for WorldService S2S at %s...", c.worldAddr)
				}
			}
		}
	}
}

func (c *Client) readLoop(conn net.Conn) {
	defer c.wg.Done()
	defer c.clearConn(conn)

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var ev Event
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}
		if c.onEvent != nil {
			c.onEvent(ev)
		}
	}
}

