package world

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/s2s"
)

func TestLiveS2SCoordinator_WorldAndZonesTCP(t *testing.T) {
	charRepo := database.NewMemoryCharacterRepository()
	_ = charRepo.CreateCharacter(&database.Character{
		ID:        501,
		AccountID: 77,
		Slot:      0,
		Name:      "S2SHero",
		Job:       1,
		Level:     20,
		Map:       0,
	})
	_ = charRepo.CreateCharacter(&database.Character{
		ID:        502,
		AccountID: 88,
		Slot:      0,
		Name:      "S2STarget",
		Job:       6,
		Level:     25,
		Map:       1,
	})

	transfers := database.NewTransferManager(nil)
	srv := NewServer(&Config{
		ID:          0,
		Name:        "S2STestWorld",
		BindAddr:    "127.0.0.1:0",
		ZoneIP:      "127.0.0.1",
		ZonePort:    9210,
		S2SBindAddr: "127.0.0.1:0",
	}, charRepo, transfers, nil)

	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start world server with S2S: %v", err)
	}
	defer srv.Stop()

	s2sAddr := srv.S2SAddr()
	if s2sAddr == "" {
		t.Fatalf("expected non-empty S2SAddr")
	}

	var mu sync.Mutex
	var zone1Events, zone2Events []s2s.Event
	recvCh := make(chan string, 16)

	z1 := s2s.NewClient(s2sAddr, 0, 1, func(ev s2s.Event) {
		mu.Lock()
		zone1Events = append(zone1Events, ev)
		mu.Unlock()
		recvCh <- "z1:" + ev.Symbol
	})
	z2 := s2s.NewClient(s2sAddr, 0, 2, func(ev s2s.Event) {
		mu.Lock()
		zone2Events = append(zone2Events, ev)
		mu.Unlock()
		recvCh <- "z2:" + ev.Symbol
	})
	if err := z1.Start(); err != nil {
		t.Fatalf("z1 start failed: %v", err)
	}
	defer z1.Stop()
	if err := z2.Start(); err != nil {
		t.Fatalf("z2 start failed: %v", err)
	}
	defer z2.Stop()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if z1.Connected() && z2.Connected() && srv.S2SConnectedZones() == 2 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if srv.S2SConnectedZones() != 2 {
		t.Fatalf("expected 2 connected zones, got %d", srv.S2SConnectedZones())
	}

	// Registrar clientes en WorldService (uno ya transferido a Zone)
	c1Conn, c1Peer := net.Pipe()
	defer c1Conn.Close()
	defer c1Peer.Close()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, err := c1Peer.Read(buf); err != nil {
				return
			}
		}
	}()

	c2Conn, c2Peer := net.Pipe()
	defer c2Conn.Close()
	defer c2Peer.Close()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, err := c2Peer.Read(buf); err != nil {
				return
			}
		}
	}()

	c1, _ := network.NewClient(c1Conn, nil)
	c1.IsAuthenticated = true
	c1.AccountID = 77
	c1.Username = "heroacc"
	c1.CharacterID = 501
	c1.CharName = "S2SHero"
	c1.Admin = 1
	c1.IsTransferring = true

	c2, _ := network.NewClient(c2Conn, nil)
	c2.IsAuthenticated = true
	c2.AccountID = 88
	c2.Username = "targetacc"
	c2.CharacterID = 502
	c2.CharName = "S2STarget"
	c2.IsTransferring = true

	h := srv.Handler()
	h.mu.Lock()
	h.activeClients[77] = c1
	h.activeClients[88] = c2
	h.mu.Unlock()
	h.partyMgr.RegisterOnline(c1, 501, "S2SHero", 1, 20)
	h.partyMgr.RegisterOnline(c2, 502, "S2STarget", 6, 25)

	// 1. Roar desde Zone 1 debe propagarse en vivo hacia Zone 2 (y no duplicarse a Zone 1)
	if err := z1.Publish(s2s.Event{
		Symbol:  "wms_NC_CT_CHARTTING_CMD",
		WorldID: 0,
		ZoneNo:  1,
		Payload: map[string]interface{}{
			"type":    "roar",
			"sender":  "S2SHero",
			"message": "Hola mundo multi-zona!",
		},
	}); err != nil {
		t.Fatalf("publish roar failed: %v", err)
	}

	waitForEvent := func(list *[]s2s.Event, sym string) bool {
		end := time.Now().Add(2 * time.Second)
		for time.Now().Before(end) {
			mu.Lock()
			for _, ev := range *list {
				if ev.Symbol == sym {
					mu.Unlock()
					return true
				}
			}
			mu.Unlock()
			time.Sleep(15 * time.Millisecond)
		}
		return false
	}
	if !waitForEvent(&zone2Events, "wms_NC_CT_CHARTTING_CMD") {
		t.Fatalf("expected Zone 2 to receive cross-zone roar wms_NC_CT_CHARTTING_CMD")
	}

	// 2. Whisper iniciado en WorldService debe enviarse a ZoneService vía S2S vivo
	wPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPER_REQ)
	wPkt.WriteByte(0)
	wPkt.WriteString("S2STarget", 20)
	msg := "Ping desde World"
	wPkt.WriteByte(byte(len(msg)))
	wPkt.WriteBytes([]byte(msg))
	h.HandlePacket(c1, &protocol.Packet{OpCode: protocol.PROTO_NC_ACT_WHISPER_REQ, Reader: protocol.NewPacketReader(wPkt.Data()[2:])})

	if !waitForEvent(&zone1Events, "wms_NC_CT_CHARTTING_CMD") {
		t.Fatalf("expected Zone 1 to receive whisper relay wms_NC_CT_CHARTTING_CMD")
	}

	// 3. Prison GM desde WorldService debe enviar wms_NC_ZONE_PRISON_GO_REQ hacia ZoneService
	pPkt := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_ADD_GM_REQ)
	pPkt.WriteString("S2STarget", 20)
	pPkt.WriteInt16(15)
	pPkt.WriteString("RuleCheck", 16)
	pPkt.WriteBytes(make([]byte, 64))
	h.HandlePacket(c1, &protocol.Packet{OpCode: protocol.PROTO_NC_PRISON_ADD_GM_REQ, Reader: protocol.NewPacketReader(pPkt.Data()[2:])})

	if !waitForEvent(&zone1Events, "wms_NC_ZONE_PRISON_GO_REQ") {
		t.Fatalf("expected Zone 1 to receive wms_NC_ZONE_PRISON_GO_REQ")
	}

	// 4. Retorno a Selección de Personaje en ZoneService (LogoutType == 1) debe limpiar sesión de personaje pero mantener TCP abierto
	if err := z1.Publish(s2s.Event{
		Symbol:  "wms_NC_USER_CONNECTCUT2ZONE_CMD",
		WorldID: 0,
		ZoneNo:  1,
		Payload: map[string]interface{}{
			"account_id":  77,
			"char_id":     501,
			"logout_type": 1,
		},
	}); err != nil {
		t.Fatalf("publish logout sync failed: %v", err)
	}

	end := time.Now().Add(2 * time.Second)
	for time.Now().Before(end) {
		if c1.CharacterID == 0 && !c1.IsTransferring {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if c1.CharacterID != 0 || c1.IsTransferring || h.ActivePlayersCount() < 2 {
		t.Fatalf("expected Return-to-CharSelect to reset CharacterID/IsTransferring while keeping World activeClient: charID=%d transf=%v count=%d",
			c1.CharacterID, c1.IsTransferring, h.ActivePlayersCount())
	}
}
