package network

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestPacketDebugLogging_S2C_And_C2S(t *testing.T) {
	var buf bytes.Buffer
	origOut := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origOut)

	origDebug := IsPacketDebug()
	defer SetPacketDebug(origDebug)

	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	client, err := NewClient(c1, nil)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// 1. Disabled -> silent
	SetPacketDebug(false)
	pwOff := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SEED_ACK)
	pwOff.WriteInt16(123)
	go func() {
		tmp := make([]byte, 64)
		_, _ = c2.Read(tmp)
	}()
	if err := client.SendPacket(pwOff); err != nil {
		t.Fatalf("SendPacket(off) failed: %v", err)
	}
	if strings.Contains(buf.String(), "[PKT]") {
		t.Fatalf("expected no [PKT] log when DEBUG is disabled, got: %s", buf.String())
	}

	// 2. Enabled -> S2C & C2S logged with OpCode name + hex
	buf.Reset()
	SetPacketDebug(true)

	pwOn := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SEED_ACK)
	pwOn.WriteInt16(456)
	go func() {
		tmp := make([]byte, 64)
		_, _ = c2.Read(tmp)
	}()
	if err := client.SendPacket(pwOn); err != nil {
		t.Fatalf("SendPacket(on) failed: %v", err)
	}
	outS2C := buf.String()
	if !strings.Contains(outS2C, "[PKT][S2C]") || !strings.Contains(outS2C, "PROTO_NC_MISC_SEED_ACK") || !strings.Contains(outS2C, "0x0807") {
		t.Fatalf("expected S2C packet debug log, got: %s", outS2C)
	}

	// Test C2S via StartReadLoop
	buf.Reset()
	done := make(chan struct{})
	go client.StartReadLoop(func(_ *Client, pkt *protocol.Packet) {
		if pkt.OpCode == protocol.PROTO_NC_MISC_HEARTBEAT_ACK {
			close(done)
		}
	})

	rawPayload := make([]byte, 2)
	binary.LittleEndian.PutUint16(rawPayload, protocol.PROTO_NC_MISC_HEARTBEAT_ACK)
	peerCrypto, err := crypto.NewNetCrypto(client.XorOffset())
	if err != nil {
		t.Fatalf("NewNetCrypto failed: %v", err)
	}
	peerCrypto.Crypt(rawPayload)
	wire := append([]byte{byte(len(rawPayload))}, rawPayload...)
	if _, err := c2.Write(wire); err != nil {
		t.Fatalf("c2.Write failed: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for C2S packet dispatch")
	}

	outC2S := buf.String()
	if !strings.Contains(outC2S, "[PKT][C2S]") || !strings.Contains(outC2S, "PROTO_NC_MISC_HEARTBEAT_ACK") {
		t.Fatalf("expected C2S packet debug log, got: %s", outC2S)
	}
}
