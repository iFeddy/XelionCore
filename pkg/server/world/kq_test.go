package world_test

import (
	"net"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

func TestKQManager_FullFlow(t *testing.T) {
	km := world.NewKQManager(nil) // Usa KQs por defecto

	c1Server, c1Client := net.Pipe()
	defer c1Server.Close()
	defer c1Client.Close()

	c2Server, c2Client := net.Pipe()
	defer c2Server.Close()
	defer c2Client.Close()

	drain := func(conn net.Conn) {
		buf := make([]byte, 4096)
		for {
			if _, err := conn.Read(buf); err != nil {
				return
			}
		}
	}
	go drain(c1Client)
	go drain(c2Client)

	client1, err := network.NewClient(c1Server, nil)
	if err != nil {
		t.Fatalf("failed to create client1: %v", err)
	}
	client1.CharacterID = 501
	client1.Username = "HeroOne"

	client2, err := network.NewClient(c2Server, nil)
	if err != nil {
		t.Fatalf("failed to create client2: %v", err)
	}
	client2.CharacterID = 502
	client2.Username = "HeroTwo"

	km.RegisterOnline(client1, 501)
	km.RegisterOnline(client2, 502)

	// 1. Solicitar listado de KQs
	km.HandleList(client1)

	// 2. Jugador 1 intenta anotarse en KQ Handle 1 (Mara Pirates, Niveles 20-30) con nivel 10 -> Debe fallar
	km.HandleJoin(client1, 1, 10, 1)

	// 3. Jugador 1 se anota con nivel válido 25 -> Éxito
	km.HandleJoin(client1, 1, 25, 1)

	// 4. Verificar estado e inscritos
	km.HandleStatus(client1, 1)
	km.HandleJoinList(client1, 1)

	// 5. Jugador 1 selecciona equipo (Team 1)
	km.HandleTeamSelect(client1, 1)

	// 6. Jugador 2 se anota con nivel válido 28
	km.HandleJoin(client2, 1, 28, 2)

	// 7. Jugador 1 cancela su registro
	km.HandleJoinCancel(client1, 1)

	// 8. Desconectar jugadores
	km.UnregisterOnline(501)
	km.UnregisterOnline(502)
}

func TestKQProtocol_PacketStructures(t *testing.T) {
	// 1. PROTO_NC_KQ_LIST_TIME_ACK: 40 bytes payload + 3 bytes header = 43 bytes
	timeAck := protocol.New_PROTO_NC_KQ_LIST_TIME_ACK(time.Now())
	if timeAck.OpCode() != protocol.PROTO_NC_KQ_LIST_TIME_ACK {
		t.Errorf("expected opcode 0x%04X, got 0x%04X", protocol.PROTO_NC_KQ_LIST_TIME_ACK, timeAck.OpCode())
	}
	wireTime := timeAck.BuildWirePacket()
	if len(wireTime) != 43 {
		t.Errorf("expected 43 bytes for PROTO_NC_KQ_LIST_TIME_ACK, got %d", len(wireTime))
	}

	// 2. PROTO_NC_KQ_JOIN_ACK: 6 bytes payload + 3 bytes header = 9 bytes
	joinAck := protocol.New_PROTO_NC_KQ_JOIN_ACK(101, 0)
	wireJoin := joinAck.BuildWirePacket()
	if len(wireJoin) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOIN_ACK, got %d", len(wireJoin))
	}

	// 3. PROTO_NC_KQ_JOIN_CANCEL_ACK: 6 bytes payload + 3 bytes header = 9 bytes
	cancelAck := protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(101, 0)
	wireCancel := cancelAck.BuildWirePacket()
	if len(wireCancel) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOIN_CANCEL_ACK, got %d", len(wireCancel))
	}

	// 4. PROTO_NC_KQ_TEAM_SELECT_ACK: 3 bytes payload + 3 bytes header = 6 bytes
	teamAck := protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(0, 1)
	wireTeam := teamAck.BuildWirePacket()
	if len(wireTeam) != 6 {
		t.Errorf("expected 6 bytes for PROTO_NC_KQ_TEAM_SELECT_ACK, got %d", len(wireTeam))
	}

	// 5. PROTO_NC_KQ_TEAM_SELECT_CMD: 21 bytes payload + 3 bytes header = 24 bytes
	teamCmd := protocol.New_PROTO_NC_KQ_TEAM_SELECT_CMD("TestChar", 2)
	wireTeamCmd := teamCmd.BuildWirePacket()
	if len(wireTeamCmd) != 24 {
		t.Errorf("expected 24 bytes for PROTO_NC_KQ_TEAM_SELECT_CMD, got %d", len(wireTeamCmd))
	}

	// 6. PROTO_NC_KQ_JOINING_ALARM_END_CMD: 6 bytes payload + 3 bytes header = 9 bytes
	alarmEnd := protocol.New_PROTO_NC_KQ_JOINING_ALARM_END_CMD(55, 2)
	wireAlarmEnd := alarmEnd.BuildWirePacket()
	if len(wireAlarmEnd) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOINING_ALARM_END_CMD, got %d", len(wireAlarmEnd))
	}
}
