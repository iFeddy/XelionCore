package network

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

var packetDebug atomic.Bool

func init() {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG"))) {
	case "true", "1", "yes", "on":
		packetDebug.Store(true)
	}
}

// SetPacketDebug enables or disables verbose C2S/S2C packet logging.
func SetPacketDebug(enabled bool) {
	packetDebug.Store(enabled)
}

// IsPacketDebug reports whether verbose C2S/S2C packet logging is enabled.
func IsPacketDebug() bool {
	return packetDebug.Load()
}

func formatPacketHex(raw []byte) string {
	const maxDumpBytes = 512
	if len(raw) == 0 {
		return ""
	}
	if len(raw) <= maxDumpBytes {
		return hex.EncodeToString(raw)
	}
	return fmt.Sprintf("%s...(+%d bytes)", hex.EncodeToString(raw[:maxDumpBytes]), len(raw)-maxDumpBytes)
}

func logPacket(dir, remoteAddr string, opCode uint16, raw []byte) {
	log.Printf("[PKT][%s][%s] op=0x%04X (%s) len=%d hex=%s",
		dir, remoteAddr, opCode, protocol.OpCodeName(opCode), len(raw), formatPacketHex(raw))
}

// Client represents an active client connection to a Xelion server.
type Client struct {
	conn       net.Conn
	remoteAddr string
	crypto     *crypto.NetCrypto

	// Client metadata
	AccountID       int
	Username        string
	CharacterID     int
	Slot            byte
	IsAuthenticated bool
	ClientYear      uint16
	ClientVersion   uint16
	IsTransferring  bool
	Admin           byte
	RandomID        uint16
	TutorialChosen  bool
	TutorialSkip    bool

	closed atomic.Bool
	sendMu sync.Mutex
	onClose func(*Client)
}

// NewClient initializes a new client connection and assigns a random NetCrypto offset.
func NewClient(conn net.Conn, onClose func(*Client)) (*Client, error) {
	xorPos := crypto.RandomXorOffset()
	nc, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:       conn,
		remoteAddr: conn.RemoteAddr().String(),
		crypto:     nc,
		onClose:    onClose,
	}

	return c, nil
}

// RemoteAddr returns the client's network address string.
func (c *Client) RemoteAddr() string {
	return c.remoteAddr
}

// RemoteIP returns only the IP address portion of the client's network address.
func (c *Client) RemoteIP() string {
	host, _, err := net.SplitHostPort(c.remoteAddr)
	if err != nil {
		return c.remoteAddr
	}
	return host
}

// XorOffset returns the initial/current XOR table position.
func (c *Client) XorOffset() int16 {
	if c.crypto != nil {
		return c.crypto.XorPos()
	}
	return 0
}

// Crypto returns the client's crypto instance.
func (c *Client) Crypto() *crypto.NetCrypto {
	return c.crypto
}

// SendPacket writes a PacketWriter buffer over the wire with Fiesta framing.
func (c *Client) SendPacket(pw *protocol.PacketWriter) error {
	if c.closed.Load() {
		return errors.New("client connection closed")
	}

	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	wire := pw.BuildWirePacket()
	_, err := c.conn.Write(wire)
	if err == nil && packetDebug.Load() {
		logPacket("S2C", c.remoteAddr, pw.OpCode(), pw.Data())
	}
	return err
}

// Close gracefully closes the client connection.
func (c *Client) Close() {
	if c.closed.CompareAndSwap(false, true) {
		_ = c.conn.Close()
		if c.onClose != nil {
			c.onClose(c)
		}
	}
}

// StartReadLoop reads and dispatches incoming Fiesta packets.
func (c *Client) StartReadLoop(packetHandler func(*Client, *protocol.Packet)) {
	defer c.Close()

	headerBuf := make([]byte, 1)

	for {
		// 1. Read first byte of framing
		if _, err := io.ReadFull(c.conn, headerBuf); err != nil {
			if !errors.Is(err, io.EOF) && !c.closed.Load() {
				log.Printf("[%s] connection read error: %v", c.remoteAddr, err)
			}
			return
		}

		var payloadLen int
		firstByte := headerBuf[0]

		if firstByte == 0 {
			// Extended length (2 bytes Little-Endian)
			lenBuf := make([]byte, 2)
			if _, err := io.ReadFull(c.conn, lenBuf); err != nil {
				return
			}
			payloadLen = int(binary.LittleEndian.Uint16(lenBuf))
		} else {
			payloadLen = int(firstByte)
		}

		if payloadLen <= 0 || payloadLen > 16384 {
			log.Printf("[%s] invalid packet payload length: %d", c.remoteAddr, payloadLen)
			return
		}

		// 2. Read full packet payload
		payload := make([]byte, payloadLen)
		if _, err := io.ReadFull(c.conn, payload); err != nil {
			return
		}

		// 3. Decrypt in-place using NetCrypto
		c.crypto.Crypt(payload)

		// 4. Construct packet and dispatch
		pkt := protocol.NewPacket(payload)
		if pkt == nil {
			log.Printf("[%s] packet shorter than OpCode header (2 bytes)", c.remoteAddr)
			return
		}
		if packetDebug.Load() {
			logPacket("C2S", c.remoteAddr, pkt.OpCode, payload)
		}

		packetHandler(c, pkt)
	}
}
