package zone

import (
	"log"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ZoneMap manages the players and entities inside a single map instance (e.g. Roumen).
type ZoneMap struct {
	mu          sync.RWMutex
	ID          uint16
	Name        string
	ShortName   string
	nextObjID   uint16
	players     map[uint16]*ZonePlayer
}

// NewZoneMap creates a new map instance for Roumen (Map 0).
func NewZoneMap(id uint16, name, shortName string) *ZoneMap {
	return &ZoneMap{
		ID:        id,
		Name:      name,
		ShortName: shortName,
		nextObjID: 1,
		players:   make(map[uint16]*ZonePlayer),
	}
}

// AssignObjectID assigns a unique MapObjectID to a player before in-game chunk transmission.
func (m *ZoneMap) AssignObjectID(p *ZonePlayer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.MapObjectID == 0 {
		p.MapObjectID = m.nextObjID
		m.nextObjID++
	}
}

// AddPlayer registers a player into the map and handles mutual visibility.
func (m *ZoneMap) AddPlayer(p *ZonePlayer) {
	m.mu.Lock()
	if p.MapObjectID == 0 {
		p.MapObjectID = m.nextObjID
		m.nextObjID++
	}
	m.players[p.MapObjectID] = p

	// Collect other players
	var others []*ZonePlayer
	for id, other := range m.players {
		if id != p.MapObjectID {
			others = append(others, other)
		}
	}
	m.mu.Unlock()

	log.Printf("[Zone][Map %d] Added '%s' (MapObjectID: %d, Pos: %d,%d). Current map population: %d",
		m.ID, p.Char.Name, p.MapObjectID, p.Char.XPos, p.Char.YPos, len(others)+1)

	// 1. Send all existing players to the joining player (SH7Type_SpawnMultiPlayer)
	spawnOthers := protocol.NewPacketWriter(protocol.SH7Type_SpawnMultiPlayer)
	spawnOthers.WriteByte(byte(len(others)))
	for _, other := range others {
		other.WriteCharacterDisplay(spawnOthers)
	}
	_ = p.Client.SendPacket(spawnOthers)

	// 2. Broadcast the joining player to all other players (SH7Type_SpawnSinglePlayer)
	spawnSelf := protocol.NewPacketWriter(protocol.SH7Type_SpawnSinglePlayer)
	p.WriteCharacterDisplay(spawnSelf)
	m.Broadcast(spawnSelf, p.MapObjectID)
}

// RemovePlayer removes a player from the map and informs others.
func (m *ZoneMap) RemovePlayer(p *ZonePlayer) {
	m.mu.Lock()
	delete(m.players, p.MapObjectID)
	objID := p.MapObjectID
	m.mu.Unlock()

	log.Printf("[Zone][Map %d] Player '%s' (MapObjectID: %d) left the map.", m.ID, p.Char.Name, objID)

	removePkt := protocol.NewPacketWriter(protocol.SH7Type_RemoveObject)
	removePkt.WriteUint16(objID)
	m.Broadcast(removePkt, objID)
}

// Broadcast sends a packet to all players on the map, optionally skipping excludeID.
func (m *ZoneMap) Broadcast(pw *protocol.PacketWriter, excludeID uint16) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, player := range m.players {
		if id == excludeID || player.Client == nil {
			continue
		}
		_ = player.Client.SendPacket(pw)
	}
}

// MovePlayer updates player position and broadcasts movement to others.
func (m *ZoneMap) MovePlayer(p *ZonePlayer, oldX, oldY, newX, newY int, walk bool) {
	p.Char.XPos = newX
	p.Char.YPos = newY

	op := protocol.SH8Type_Move
	if walk {
		op = protocol.SH8Type_Walk
	}

	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteInt32(int32(oldX))
	pw.WriteInt32(int32(oldY))
	pw.WriteInt32(int32(newX))
	pw.WriteInt32(int32(newY))
	pw.WriteUint16(115) // Move speed
	m.Broadcast(pw, p.MapObjectID)
}

// StopPlayer updates position and broadcasts stop to others.
func (m *ZoneMap) StopPlayer(p *ZonePlayer, x, y int) {
	p.Char.XPos = x
	p.Char.YPos = y

	pw := protocol.NewPacketWriter(protocol.SH8Type_StopTele)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteInt32(int32(x))
	pw.WriteInt32(int32(y))
	m.Broadcast(pw, p.MapObjectID)
}

// JumpPlayer broadcasts player jump.
func (m *ZoneMap) JumpPlayer(p *ZonePlayer) {
	pw := protocol.NewPacketWriter(protocol.SH8Type_Jump)
	pw.WriteUint16(p.MapObjectID)
	m.Broadcast(pw, p.MapObjectID)
}

// NormalChat broadcasts chat to all players in the map (including sender).
func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte) {
	pw := protocol.NewPacketWriter(protocol.SH8Type_ChatNormal)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(byte(len(text)))
	pw.WriteByte(color)
	pw.WriteString(text, len(text))
	m.Broadcast(pw, 0)
}
