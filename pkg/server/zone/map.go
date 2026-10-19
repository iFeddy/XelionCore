package zone

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ZoneMap manages the players and entities inside a single map instance.
type ZoneMap struct {
	mu        sync.RWMutex
	ID        uint16
	MapID     string
	Name      string
	ShortName string
	Width     int
	Height    int
	nextObjID uint16
	players   map[uint16]*ZonePlayer
}

// NewZoneMap creates a new map instance with real metadata.
func NewZoneMap(id uint16, mapID, name string) *ZoneMap {
	return &ZoneMap{
		ID:        id,
		MapID:     mapID,
		Name:      name,
		ShortName: mapID,
		nextObjID: 1,
		players:   make(map[uint16]*ZonePlayer),
	}
}

// MapManager manages active ZoneMap instances dynamically mapped from data/World/Map definitions.
type MapManager struct {
	mu         sync.RWMutex
	fieldStore *field.Store
	maps       map[string]*ZoneMap
}

// NewMapManager initializes a new dynamic MapManager.
func NewMapManager(fieldStore *field.Store) *MapManager {
	return &MapManager{
		fieldStore: fieldStore,
		maps:       make(map[string]*ZoneMap),
	}
}

// GetOrCreateMap retrieves an existing active map or instantiates one from field definitions.
func (mm *MapManager) GetOrCreateMap(mapID string) *ZoneMap {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	normID := strings.TrimSpace(mapID)
	if normID == "" {
		normID = "RouN"
	}
	key := strings.ToLower(normID)

	if existing, ok := mm.maps[key]; ok {
		return existing
	}

	var (
		serial uint16 = 0
		name   string = normID
		width  int    = 256
		height int    = 256
	)

	if mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetField(normID); ok && entry != nil {
			serial = uint16(entry.Serial)
			name = entry.Name
			width = entry.Width
			height = entry.Height
			normID = entry.MapID
		}
	}

	zMap := &ZoneMap{
		ID:        serial,
		MapID:     normID,
		Name:      name,
		ShortName: normID,
		Width:     width,
		Height:    height,
		nextObjID: 1,
		players:   make(map[uint16]*ZonePlayer),
	}
	mm.maps[key] = zMap
	log.Printf("[Zone] Initialized dynamic map '%s' (Serial: %d, Name: '%s', Size: %dx%d)",
		normID, serial, name, width, height)
	return zMap
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

	// 1. Send all existing players to the joining player (PROTO_NC_BRIEFINFO_CHARACTER_ACK)
	// En Zone.exe oficial (bixc_Send 0x43D5D8), si no hay otros jugadores en el mapa (count == 0),
	// el servidor NO envía el paquete 0x1C07 al cliente.
	if len(others) > 0 {
		spawnOthers := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHARACTER_ACK)
		spawnOthers.WriteByte(byte(len(others)))
		for _, other := range others {
			other.WriteCharacterDisplay(spawnOthers)
		}
		_ = p.Client.SendPacket(spawnOthers)
	}

	// 2. Broadcast the joining player to all other players (PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
	spawnSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
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

	removePkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LEAVE_CMD)
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

	op := protocol.PROTO_NC_ACT_RUN_ACK
	if walk {
		op = protocol.PROTO_NC_ACT_WALK_ACK
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

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_STOP_ACK)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteInt32(int32(x))
	pw.WriteInt32(int32(y))
	m.Broadcast(pw, p.MapObjectID)
}

// JumpPlayer broadcasts player jump.
func (m *ZoneMap) JumpPlayer(p *ZonePlayer) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_JUMP_ACK)
	pw.WriteUint16(p.MapObjectID)
	m.Broadcast(pw, p.MapObjectID)
}

// NormalChat broadcasts chat to all players in the map (including sender).
func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_CHAT_ACK)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(byte(len(text)))
	pw.WriteByte(color)
	pw.WriteString(text, len(text))
	m.Broadcast(pw, 0)
}
