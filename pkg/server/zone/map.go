package zone

import (
	"log"
	"math"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ZoneNPC representa un NPC o Gate activo en una instancia de mapa según la definición oficial.
type ZoneNPC struct {
	Handle      uint16
	MobID       uint16
	Name        string
	Map         string
	X           int32
	Y           int32
	Direct      int16
	Menu        uint8
	Role        string
	RoleArg     string
	IsGate      bool
	GateToWhere string
}

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
	npcs      map[uint16]*ZoneNPC
	gates     []*ZoneNPC
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
		npcs:      make(map[uint16]*ZoneNPC),
	}
}

// MapManager manages active ZoneMap instances dynamically mapped from data/World/Map definitions.
type MapManager struct {
	mu         sync.RWMutex
	fieldStore *field.Store
	npcStore   *npc.Store
	mobMap     map[string]uint16 // InxName -> MobID
	maps       map[string]*ZoneMap
}

// NewMapManager initializes a new dynamic MapManager.
func NewMapManager(fieldStore *field.Store) *MapManager {
	return &MapManager{
		fieldStore: fieldStore,
		mobMap:     make(map[string]uint16),
		maps:       make(map[string]*ZoneMap),
	}
}

// SetNPCStore configures the NPC store and builds the MobID mapping from MobInfo.shn.
func (mm *MapManager) SetNPCStore(npcStore *npc.Store, mobInfoSHN *shn.File) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.npcStore = npcStore
	if mobInfoSHN != nil {
		for _, row := range mobInfoSHN.Rows {
			name := strings.TrimSpace(row.GetString("InxName"))
			id := row.GetUint16("ID")
			if name != "" && id > 0 {
				mm.mobMap[strings.ToLower(name)] = id
			}
		}
		log.Printf("[Zone] Indexed %d monsters/NPCs from MobInfo.shn.", len(mm.mobMap))
	}

	for _, zMap := range mm.maps {
		mm.loadMapNPCs(zMap)
	}
}

func (mm *MapManager) resolveMobID(mobName string) uint16 {
	k := strings.ToLower(strings.TrimSpace(mobName))
	if id, ok := mm.mobMap[k]; ok {
		return id
	}
	return 1
}

func (mm *MapManager) loadMapNPCs(zMap *ZoneMap) {
	if mm.npcStore == nil {
		return
	}
	zMap.mu.Lock()
	defer zMap.mu.Unlock()

	npcs := mm.npcStore.GetNPCs(zMap.MapID, true)
	if len(npcs) == 0 {
		return
	}

	zMap.npcs = make(map[uint16]*ZoneNPC)
	zMap.gates = nil

	for _, n := range npcs {
		mobID := mm.resolveMobID(n.MobName)
		handle := zMap.nextObjID
		zMap.nextObjID++

		isGate := n.Role == "Gate" || n.Role == "IDGate" || strings.Contains(strings.ToLower(n.MobName), "gate")
		gateToWhere := ""
		if isGate {
			gateToWhere = n.RoleArg
		}

		zNPC := &ZoneNPC{
			Handle:      handle,
			MobID:       mobID,
			Name:        n.MobName,
			Map:         zMap.MapID,
			X:           n.X,
			Y:           n.Y,
			Direct:      n.Direct,
			Menu:        n.Menu,
			Role:        n.Role,
			RoleArg:     n.RoleArg,
			IsGate:      isGate,
			GateToWhere: gateToWhere,
		}
		zMap.npcs[handle] = zNPC
		if isGate {
			zMap.gates = append(zMap.gates, zNPC)
		}
	}
	log.Printf("[Zone][Map %s] Loaded %d NPCs (%d gates) into world.", zMap.MapID, len(zMap.npcs), len(zMap.gates))
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
		npcs:      make(map[uint16]*ZoneNPC),
	}
	mm.maps[key] = zMap
	log.Printf("[Zone] Initialized dynamic map '%s' (Serial: %d, Name: '%s', Size: %dx%d)",
		normID, serial, name, width, height)

	mm.loadMapNPCs(zMap)
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

	// 3. Send all NPCs and Gates on the map (PROTO_NC_BRIEFINFO_MOB_CMD, 0x1C09)
	m.SendNPCs(p)
}

// SendNPCs transmits all map NPCs and Gates to a player using official PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09).
func (m *ZoneMap) SendNPCs(p *ZonePlayer) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.npcs) == 0 {
		return
	}

	chunkSize := 100
	npcList := make([]*ZoneNPC, 0, len(m.npcs))
	for _, npc := range m.npcs {
		npcList = append(npcList, npc)
	}

	for i := 0; i < len(npcList); i += chunkSize {
		end := i + chunkSize
		if end > len(npcList) {
			end = len(npcList)
		}
		chunk := npcList[i:end]

		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_MOB_CMD)
		pw.WriteByte(byte(len(chunk))) // mobnum (1 byte)
		for _, npc := range chunk {
			pw.WriteUint16(npc.Handle)     // 0x00: handle (2 bytes)
			pw.WriteByte(1)                // 0x02: mode (1 byte, 1=idle/normal)
			pw.WriteUint16(npc.MobID)      // 0x03: mobid (2 bytes)
			pw.WriteInt32(npc.X)           // 0x05: coord.x (4 bytes)
			pw.WriteInt32(npc.Y)           // 0x09: coord.y (4 bytes)
			pw.WriteByte(byte(npc.Direct)) // 0x0D: coord.dir (1 byte)
			if npc.IsGate {
				pw.WriteByte(1)                     // 0x0E: flagstate = 1 (gate)
				pw.WriteString(npc.GateToWhere, 12) // 0x0F: gate2where (12 bytes)
				pw.Fill(87, 0)                      // 0x1B: remaining of 99 bytes flag
			} else {
				pw.WriteByte(0) // 0x0E: flagstate = 0
				pw.Fill(99, 0)  // 0x0F: abstatebit (99 bytes)
			}
			pw.WriteString("", 32) // 0x72: sAnimation (32 bytes)
			pw.WriteByte(0)        // 0x92: nAnimationLevel (1 byte)
			pw.WriteByte(0)        // 0x93: nKQTeamType (1 byte)
			pw.WriteByte(0)        // 0x94: bRegenAni (1 byte)
		}
		_ = p.Client.SendPacket(pw)
	}
}

// GetNPC retrieves an active NPC on this map by its handle.
func (m *ZoneMap) GetNPC(handle uint16) (*ZoneNPC, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	npc, ok := m.npcs[handle]
	return npc, ok
}

// GetAllNPCs returns a snapshot of all active NPCs on this map.
func (m *ZoneMap) GetAllNPCs() []*ZoneNPC {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*ZoneNPC, 0, len(m.npcs))
	for _, n := range m.npcs {
		list = append(list, n)
	}
	return list
}

// GetGates returns all gate NPCs configured on this map.
func (m *ZoneMap) GetGates() []*ZoneNPC {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*ZoneNPC, len(m.gates))
	copy(list, m.gates)
	return list
}

// FindGateNear searches for a gate within maxDist of the specified coordinates.
func (m *ZoneMap) FindGateNear(x, y int, maxDist float64) (*ZoneNPC, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, g := range m.gates {
		dx := float64(int(g.X) - x)
		dy := float64(int(g.Y) - y)
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= maxDist {
			return g, true
		}
	}
	return nil, false
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

// MovePlayerCmd actualiza la posición del jugador y transmite PROTO_NC_ACT_SOMEONEMOVEWALK_CMD (0x2018)
// o PROTO_NC_ACT_SOMEONEMOVERUN_CMD (0x201A) de 22 bytes según Zone.pdb (sp_NC_ACT_MOVEWALK_CMD / sp_NC_ACT_MOVERUN_CMD).
func (m *ZoneMap) MovePlayerCmd(p *ZonePlayer, fromX, fromY, toX, toY int, walk bool) {
	p.Char.XPos = toX
	p.Char.YPos = toY

	op := protocol.PROTO_NC_ACT_SOMEONEMOVERUN_CMD
	speed := uint16(115)
	if walk {
		op = protocol.PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
		speed = 50
	}

	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(p.MapObjectID) // 0x00: handle (2 bytes)
	pw.WriteInt32(int32(fromX))   // 0x02: from.x (4 bytes)
	pw.WriteInt32(int32(fromY))   // 0x06: from.y (4 bytes)
	pw.WriteInt32(int32(toX))     // 0x0A: to.x (4 bytes)
	pw.WriteInt32(int32(toY))     // 0x0E: to.y (4 bytes)
	pw.WriteUint16(speed)         // 0x12: speed (2 bytes)
	pw.WriteUint16(0)             // 0x14: moveattr (2 bytes, direct=0)
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

// EmoticonPlayer transmite PROTO_NC_ACT_SOMEONEEMOTICON_CMD (0x2021) de 3 bytes a todos en el mapa (incluido el emisor)
// según sp_DoEmotion en Zone.pdb.
func (m *ZoneMap) EmoticonPlayer(p *ZonePlayer, emoticonID byte) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEEMOTICON_CMD)
	pw.WriteUint16(p.MapObjectID) // 0x00: handle (2 bytes)
	pw.WriteByte(emoticonID)      // 0x02: emoticonid (1 byte)
	m.Broadcast(pw, 0)
}

// EmoticonStopPlayer transmite PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD (0x2023) de 5 bytes a todos en el mapa (incluido el emisor)
// según sp_NC_ACT_EMOTICONSTOP_CMD en Zone.pdb.
func (m *ZoneMap) EmoticonStopPlayer(p *ZonePlayer, emoticonID byte, emoticonFrame uint16) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD)
	pw.WriteUint16(p.MapObjectID) // 0x00: handle (2 bytes)
	pw.WriteByte(emoticonID)      // 0x02: emoticon.emoticonid (1 byte)
	pw.WriteUint16(emoticonFrame) // 0x03: emoticon.emoticonframe (2 bytes)
	m.Broadcast(pw, 0)
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
