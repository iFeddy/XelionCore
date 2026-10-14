package login

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler processes incoming packets for the Login service.
type Handler struct {
	repo         database.AccountRepository
	worldManager *WorldManager
}

// NewHandler creates a new Login packet handler.
func NewHandler(repo database.AccountRepository, wm *WorldManager) *Handler {
	return &Handler{
		repo:         repo,
		worldManager: wm,
	}
}

// HandlePacket dispatches the incoming packet to its dedicated routine based on OpCode.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.CH3Type_Version:
		h.handleVersion(client, pkt)
	case protocol.CH3Type_Login:
		h.handleLogin(client, pkt)
	case protocol.CH3Type_FileHash:
		h.handleFileHash(client, pkt)
	case protocol.CH3Type_WorldReRequest:
		h.handleWorldReRequest(client, pkt)
	case protocol.CH3Type_WorldSelect:
		h.handleWorldSelect(client, pkt)
	default:
		log.Printf("[%s] Unhandled Login OpCode: 0x%04X (Header=%d, Type=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

func (h *Handler) handleVersion(client *network.Client, pkt *protocol.Packet) {
	year, err := pkt.Reader.ReadUint16()
	if err != nil {
		log.Printf("[%s] failed to read client version year: %v", client.RemoteAddr(), err)
		client.Close()
		return
	}
	version, err := pkt.Reader.ReadUint16()
	if err != nil {
		log.Printf("[%s] failed to read client version number: %v", client.RemoteAddr(), err)
		client.Close()
		return
	}

	client.ClientYear = year
	client.ClientVersion = version

	log.Printf("[%s] Client authenticated version: Year=%d, Version=%d", client.RemoteAddr(), year, version)

	// Send version allowed response
	resp := protocol.NewPacketWriter(protocol.SH3Type_VersionAllowed)
	resp.WriteInt16(1)
	if err := client.SendPacket(resp); err != nil {
		log.Printf("[%s] error sending VersionAllowed: %v", client.RemoteAddr(), err)
	}
}

func (h *Handler) handleLogin(client *network.Client, pkt *protocol.Packet) {
	// Payload structure:
	// - Username: 260 bytes (null-terminated)
	// - Password: 32 bytes (MD5 hex string)
	// - Padding: 4 bytes
	// - ClientTag: 8 bytes ("Original")
	username, err := pkt.Reader.ReadString(260)
	if err != nil {
		log.Printf("[%s] failed to read username: %v", client.RemoteAddr(), err)
		h.sendError(client, protocol.ErrException)
		return
	}

	clientPassword, err := pkt.Reader.ReadString(32)
	if err != nil {
		log.Printf("[%s] failed to read password: %v", client.RemoteAddr(), err)
		h.sendError(client, protocol.ErrException)
		return
	}

	_ = pkt.Reader.Skip(4) // Skip padding

	clientTag, _ := pkt.Reader.ReadString(8)

	username = strings.TrimSpace(username)
	log.Printf("[%s] Login attempt from user '%s' (tag: '%s')", client.RemoteAddr(), username, clientTag)

	if len(username) == 0 || len(username) > 32 {
		log.Printf("[%s] rejected invalid username length", client.RemoteAddr())
		h.sendError(client, protocol.ErrBlocked)
		return
	}

	// Database lookup / auto account creation
	acc, err := h.repo.GetAccountByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrAccountNotFound) {
			// Auto create account
			log.Printf("[%s] Account '%s' not found, auto-creating...", client.RemoteAddr(), username)
			acc, err = h.repo.CreateAccount(username, clientPassword)
			if err != nil {
				log.Printf("[%s] failed to auto-create account: %v", client.RemoteAddr(), err)
				h.sendError(client, protocol.ErrDatabase)
				return
			}
		} else {
			log.Printf("[%s] database error on login: %v", client.RemoteAddr(), err)
			h.sendError(client, protocol.ErrDatabase)
			return
		}
	}

	// Validate credentials: compare MD5 hex string or plaintext
	computedMD5 := md5Hash(acc.Password)
	if clientPassword != acc.Password && clientPassword != computedMD5 {
		log.Printf("[%s] Invalid password for user '%s'", client.RemoteAddr(), username)
		h.sendError(client, protocol.ErrInvalidCredentials)
		return
	}

	if acc.Banned {
		log.Printf("[%s] Account '%s' is banned", client.RemoteAddr(), username)
		h.sendError(client, protocol.ErrBlocked)
		return
	}

	client.AccountID = acc.ID
	client.Username = acc.Username
	client.Admin = acc.Admin
	client.IsAuthenticated = true

	log.Printf("[%s] Successfully authenticated as '%s' (ID: %d)", client.RemoteAddr(), client.Username, client.AccountID)

	// Send initial World List
	h.sendWorldList(client, false)
}

func (h *Handler) handleFileHash(client *network.Client, pkt *protocol.Packet) {
	_, err := pkt.Reader.ReadString(pkt.Reader.Remaining())
	if err != nil {
		h.sendError(client, protocol.ErrException)
		return
	}

	resp := protocol.NewPacketWriter(protocol.SH3Type_FilecheckAllow)
	resp.WriteBool(true)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleWorldReRequest(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[%s] unauthenticated world re-request rejected", client.RemoteAddr())
		return
	}
	h.sendWorldList(client, true)
}

func (h *Handler) handleWorldSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.IsTransferring {
		h.sendError(client, protocol.ErrException)
		return
	}

	worldID, err := pkt.Reader.ReadByte()
	if err != nil {
		log.Printf("[%s] invalid world select packet", client.RemoteAddr())
		return
	}

	world, exists := h.worldManager.GetWorld(worldID)
	if !exists {
		log.Printf("[%s] user '%s' selected non-existent world %d", client.RemoteAddr(), client.Username, worldID)
		return
	}

	if world.Status == protocol.WorldStatus_Maintenance || world.Status == protocol.WorldStatus_Offline {
		log.Printf("[%s] user '%s' attempted to join offline/maintenance world %d", client.RemoteAddr(), client.Username, worldID)
		h.sendError(client, protocol.ErrServerMaintenance)
		return
	}

	token := GenerateTransferToken()
	client.IsTransferring = true

	log.Printf("[%s] Transferring '%s' to World '%s' (%s:%d)",
		client.RemoteAddr(), client.Username, world.Name, world.IP, world.Port)

	// Send WorldServerIP packet
	resp := protocol.NewPacketWriter(protocol.SH3Type_WorldServerIP)
	resp.WriteByte(byte(world.Status))
	resp.WriteString(world.IP, 16)
	resp.WriteUint16(world.Port)
	resp.WriteString(token, 32)
	resp.Fill(32, 0) // 32 bytes padding

	_ = client.SendPacket(resp)
}

func (h *Handler) sendWorldList(client *network.Client, isPing bool) {
	opCode := protocol.SH3Type_WorldlistNew
	if isPing {
		opCode = protocol.SH3Type_WorldistResend
	}

	pw := protocol.NewPacketWriter(opCode)
	const maxWorldSlots = 11
	pw.WriteByte(maxWorldSlots)

	worlds := h.worldManager.WorldsList()
	for _, w := range worlds {
		pw.WriteByte(w.ID)
		pw.WriteString(w.Name, 16)
		pw.WriteByte(byte(w.Status))
	}

	// Fill remaining dummy slots up to maxWorldSlots
	for i := len(worlds); i < maxWorldSlots; i++ {
		pw.WriteByte(byte(i))
		pw.WriteString(fmt.Sprintf("DUMMY%d", i), 16)
		pw.WriteByte(byte(protocol.WorldStatus_Offline))
	}

	_ = client.SendPacket(pw)
}

func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.NewPacketWriter(protocol.SH3Type_Error)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
