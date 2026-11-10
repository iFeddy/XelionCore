package zone

import (
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/network"
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

// GetHandle implementa GridEntity para ZoneNPC.
func (n *ZoneNPC) GetHandle() uint16 {
	return n.Handle
}

// GetPosition implementa GridEntity para ZoneNPC.
func (n *ZoneNPC) GetPosition() (int, int) {
	return int(n.X), int(n.Y)
}

// GetEntityType implementa GridEntity para ZoneNPC.
func (n *ZoneNPC) GetEntityType() GridEntityType {
	return GridEntityNPC
}

// GetClient implementa GridEntity para ZoneNPC.
func (n *ZoneNPC) GetClient() *network.Client {
	return nil
}

// ZoneDroppedItem representa un ítem tirado en el suelo del mapa según PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (0x1C0A).
type ZoneDroppedItem struct {
	Handle        uint16
	ItemID        uint16
	Lot           uint32
	Flags         int
	X             int32
	Y             int32
	DropMobHandle uint16
	Attr          uint8
	OwnerCharID   int
	CreatedAt     time.Time
}

// GetHandle implementa GridEntity para ZoneDroppedItem.
func (d *ZoneDroppedItem) GetHandle() uint16 {
	return d.Handle
}

// GetPosition implementa GridEntity para ZoneDroppedItem.
func (d *ZoneDroppedItem) GetPosition() (int, int) {
	return int(d.X), int(d.Y)
}

// GetEntityType implementa GridEntity para ZoneDroppedItem.
func (d *ZoneDroppedItem) GetEntityType() GridEntityType {
	return GridEntityItem
}

// GetClient implementa GridEntity para ZoneDroppedItem.
func (d *ZoneDroppedItem) GetClient() *network.Client {
	return nil
}

// ZoneMap manages the players and entities inside a single map instance with spatial grid AOI.
type ZoneMap struct {
	mu           sync.RWMutex
	ID           uint16
	MapID        string
	Name         string
	ShortName    string
	Width        int
	Height       int
	nextObjID    uint16
	players      map[uint16]*ZonePlayer
	npcs         map[uint16]*ZoneNPC
	gates        []*ZoneNPC
	droppedItems map[uint16]*ZoneDroppedItem
	tents        map[uint16]uint16 // MapObjectID -> tentID
	grid         *SpatialGrid
	blockInfo    *blockinfo.MapBlockInformation
	doors        *blockinfo.MapDoorArray
	gameLoop     *MapGameLoop
}

// NewZoneMap creates a new map instance with real metadata and starts its GameLoop.
func NewZoneMap(id uint16, mapID, name string) *ZoneMap {
	zm := &ZoneMap{
		ID:           id,
		MapID:        mapID,
		Name:         name,
		ShortName:    mapID,
		nextObjID:    1,
		players:      make(map[uint16]*ZonePlayer),
		npcs:         make(map[uint16]*ZoneNPC),
		droppedItems: make(map[uint16]*ZoneDroppedItem),
		tents:        make(map[uint16]uint16),
		grid:         NewSpatialGrid(DefaultCellSize),
	}
	zm.gameLoop = NewMapGameLoop(zm)
	zm.gameLoop.Start()
	return zm
}

// SetCollisionData configura la geometría binaria (.shbd) y puertas (.sbi) para validación de movimiento.
func (m *ZoneMap) SetCollisionData(mbi *blockinfo.MapBlockInformation, doors *blockinfo.MapDoorArray) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blockInfo = mbi
	m.doors = doors
	if mbi != nil && mbi.WidthTiles > 0 && mbi.HeightTiles > 0 {
		m.Width = int(mbi.WidthTiles * 100 / 8)
		m.Height = int(mbi.HeightTiles * 100 / 8)
	}
}

// Close finaliza el bucle de juego y libera recursos del mapa.
func (m *ZoneMap) Close() {
	if m.gameLoop != nil {
		m.gameLoop.Stop()
	}
}

// MapManager manages active ZoneMap instances dynamically mapped from data/World/Map definitions.
type MapManager struct {
	mu         sync.RWMutex
	fieldStore *field.Store
	npcStore   *npc.Store
	blockStore *blockinfo.Store
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

// SetBlockStore configura el repositorio de colisiones .shbd y puertas .sbi.
func (mm *MapManager) SetBlockStore(blockStore *blockinfo.Store) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.blockStore = blockStore

	for _, zMap := range mm.maps {
		if mm.blockStore != nil {
			mbi := mm.blockStore.GetSHBD(zMap.MapID)
			doors := mm.blockStore.NewInstanceDoors(zMap.MapID)
			zMap.SetCollisionData(mbi, doors)
		}
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
		zMap.grid.AddEntity(zNPC)
	}
	log.Printf("[Zone][Map %s] Loaded %d NPCs (%d gates) into world and AOI.", zMap.MapID, len(zMap.npcs), len(zMap.gates))
}

// GetOrCreateMap retrieves an existing active map or instantiates one from field definitions.
func (mm *MapManager) GetOrCreateMap(mapID string) *ZoneMap {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	normID := strings.TrimSpace(mapID)
	if normID == "" && mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetFieldBySerial(0); ok && entry != nil {
			normID = entry.MapID
		}
	}
	key := strings.ToLower(normID)

	if existing, ok := mm.maps[key]; ok {
		return existing
	}

	var (
		serial uint16 = 0
		name   string = normID
		width  int    = 20000
		height int    = 20000
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

	zMap := NewZoneMap(serial, normID, name)
	zMap.Width = width
	zMap.Height = height

	if mm.blockStore != nil {
		mbi := mm.blockStore.GetSHBD(normID)
		doors := mm.blockStore.NewInstanceDoors(normID)
		zMap.SetCollisionData(mbi, doors)
	}

	mm.maps[key] = zMap
	log.Printf("[Zone] Initialized dynamic map '%s' (Serial: %d, Name: '%s', Size: %dx%d)",
		normID, serial, name, zMap.Width, zMap.Height)

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

// AddPlayer registers a player into the map and handles mutual visibility via AOI.
func (m *ZoneMap) AddPlayer(p *ZonePlayer) {
	m.mu.Lock()
	if p.MapObjectID == 0 {
		p.MapObjectID = m.nextObjID
		m.nextObjID++
	}
	m.players[p.MapObjectID] = p
	m.mu.Unlock()

	p.CurrentMap = m
	m.grid.AddEntity(p)

	log.Printf("[Zone][Map %d] Added '%s' (MapObjectID: %d, Pos: %d,%d).",
		m.ID, p.Char.Name, p.MapObjectID, p.Char.XPos, p.Char.YPos)

	// Obtener entidades dentro del radio AOI (vecindad 3x3)
	surrounding := m.grid.GetSurroundingEntities(p.Char.XPos, p.Char.YPos)
	var otherPlayers []*ZonePlayer
	var itemsInAOI []*ZoneDroppedItem

	for _, ent := range surrounding {
		if ent.GetHandle() == p.MapObjectID {
			continue
		}
		switch ent.GetEntityType() {
		case GridEntityPlayer:
			if zp, ok := ent.(*ZonePlayer); ok {
				otherPlayers = append(otherPlayers, zp)
			}
		case GridEntityItem:
			if zi, ok := ent.(*ZoneDroppedItem); ok {
				itemsInAOI = append(itemsInAOI, zi)
			}
		}
	}

	// 1. Enviar jugadores en AOI al jugador entrante
	if len(otherPlayers) > 0 {
		spawnOthers := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHARACTER_CMD)
		spawnOthers.WriteByte(byte(len(otherPlayers)))
		for _, other := range otherPlayers {
			other.WriteCharacterDisplay(spawnOthers)
		}
		_ = p.Client.SendPacket(spawnOthers)
	}

	// 2. Enviar jugador entrante a los otros jugadores en AOI
	spawnSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
	p.WriteCharacterDisplay(spawnSelf)
	for _, other := range otherPlayers {
		if other.Client != nil && !other.HideExceptMe {
			_ = other.Client.SendPacket(spawnSelf)
		}
	}

	// 3. Enviar NPCs y Gates del mapa para visualización y minimapa
	m.SendNPCs(p)

	// 4. Enviar ítems en AOI al jugador
	if len(itemsInAOI) > 0 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD)
		pw.WriteByte(byte(len(itemsInAOI)))
		for _, di := range itemsInAOI {
			pw.WriteUint16(di.Handle)
			pw.WriteUint16(di.ItemID)
			pw.WriteInt32(di.X)
			pw.WriteInt32(di.Y)
			pw.WriteUint16(di.DropMobHandle)
			pw.WriteByte(di.Attr)
		}
		_ = p.Client.SendPacket(pw)
	}

	// 5. Enviar sincronización de HP
	if p.Client != nil {
		hpUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpUpdate.WriteUint32(uint32(p.Char.CurHP))
		hpUpdate.WriteUint16(0)
		_ = p.Client.SendPacket(hpUpdate)
	}
}

// SpawnDroppedItem spawns an item on the ground and broadcasts PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (0x1C0A, 15 bytes) to AOI.
func (m *ZoneMap) SpawnDroppedItem(itemID uint16, lot uint32, flags int, x, y int32, dropMobHandle uint16, ownerCharID int) *ZoneDroppedItem {
	m.mu.Lock()
	handle := m.nextObjID
	m.nextObjID++
	di := &ZoneDroppedItem{
		Handle:        handle,
		ItemID:        itemID,
		Lot:           lot,
		Flags:         flags,
		X:             x,
		Y:             y,
		DropMobHandle: dropMobHandle,
		Attr:          0x01,
		OwnerCharID:   ownerCharID,
		CreatedAt:     time.Now(),
	}
	if m.droppedItems == nil {
		m.droppedItems = make(map[uint16]*ZoneDroppedItem)
	}
	m.droppedItems[handle] = di
	m.mu.Unlock()

	m.grid.AddEntity(di)

	// PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (0x1C0A, 15 bytes)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
	pw.WriteUint16(di.Handle)        // 0x00: handle (2B)
	pw.WriteUint16(di.ItemID)        // 0x02: itemid (2B)
	pw.WriteInt32(di.X)              // 0x04: location.x (4B)
	pw.WriteInt32(di.Y)              // 0x08: location.y (4B)
	pw.WriteUint16(di.DropMobHandle) // 0x0C: dropmobhandle (2B)
	pw.WriteByte(di.Attr)            // 0x0E: attr (1B)
	m.grid.BroadcastAOI(pw, int(x), int(y), 0)

	return di
}

// ConsumeDroppedItem removes a dropped item from the map and broadcasts PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD (0x1C0E) to AOI.
func (m *ZoneMap) ConsumeDroppedItem(handle uint16) (*ZoneDroppedItem, bool) {
	m.mu.Lock()
	di, ok := m.droppedItems[handle]
	if ok {
		delete(m.droppedItems, handle)
	}
	m.mu.Unlock()

	if !ok || di == nil {
		return nil, false
	}

	m.grid.RemoveEntity(handle)

	delPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
	delPkt.WriteUint16(handle)
	m.grid.BroadcastAOI(delPkt, int(di.X), int(di.Y), 0)

	return di, true
}

// CleanupExpiredDrops elimina ítems del suelo cuya antigüedad supera maxAge.
func (m *ZoneMap) CleanupExpiredDrops(maxAge time.Duration) {
	m.mu.Lock()
	var expired []*ZoneDroppedItem
	now := time.Now()
	for h, di := range m.droppedItems {
		if now.Sub(di.CreatedAt) >= maxAge {
			expired = append(expired, di)
			delete(m.droppedItems, h)
		}
	}
	m.mu.Unlock()

	for _, di := range expired {
		m.grid.RemoveEntity(di.Handle)
		delPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
		delPkt.WriteUint16(di.Handle)
		m.grid.BroadcastAOI(delPkt, int(di.X), int(di.Y), 0)
	}
}

// SendDroppedItems transmits all dropped items on the field to a joining player using PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD (0x1C0B).
func (m *ZoneMap) SendDroppedItems(p *ZonePlayer) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.droppedItems) == 0 {
		return
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD)
	pw.WriteByte(byte(len(m.droppedItems)))
	for _, di := range m.droppedItems {
		pw.WriteUint16(di.Handle)
		pw.WriteUint16(di.ItemID)
		pw.WriteInt32(di.X)
		pw.WriteInt32(di.Y)
		pw.WriteUint16(di.DropMobHandle)
		pw.WriteByte(di.Attr)
	}
	_ = p.Client.SendPacket(pw)
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

// RemovePlayer removes a player from the map and informs surrounding AOI.
func (m *ZoneMap) RemovePlayer(p *ZonePlayer) {
	m.mu.Lock()
	delete(m.players, p.MapObjectID)
	delete(m.tents, p.MapObjectID)
	objID := p.MapObjectID
	posX := p.Char.XPos
	posY := p.Char.YPos
	m.mu.Unlock()

	m.grid.RemoveEntity(objID)

	log.Printf("[Zone][Map %d] Player '%s' (MapObjectID: %d) left the map.", m.ID, p.Char.Name, objID)

	removePkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
	removePkt.WriteUint16(objID)
	m.grid.BroadcastAOI(removePkt, posX, posY, objID)
}

// Broadcast sends a packet to all players on the map, optionally skipping excludeID.
func (m *ZoneMap) Broadcast(pw *protocol.PacketWriter, excludeID uint16) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, player := range m.players {
		if id == excludeID || player.Client == nil || player.HideExceptMe {
			continue
		}
		_ = player.Client.SendPacket(pw)
	}
}

// MovePlayer actualiza la posición del jugador y transmite el movimiento con validación de colisiones.
func (m *ZoneMap) MovePlayer(p *ZonePlayer, oldX, oldY, newX, newY int, walk bool) {
	m.MovePlayerCmd(p, oldX, oldY, newX, newY, walk)
}

// MovePlayerCmd valida colisiones contra .shbd y .sbi, actualiza la grilla AOI y transmite movimiento.
func (m *ZoneMap) MovePlayerCmd(p *ZonePlayer, fromX, fromY, toX, toY int, walk bool) {
	// 1. Verificación de colisiones contra .shbd (geometría de terreno)
	if m.blockInfo != nil && m.blockInfo.IsWorldMoveBlocked(float64(toX), float64(toY)) {
		if p.Client != nil {
			failPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVEFAIL_CMD)
			failPkt.WriteInt32(int32(fromX))
			failPkt.WriteInt32(int32(fromY))
			_ = p.Client.SendPacket(failPkt)
		}
		return
	}

	// 2. Verificación de colisiones contra .sbi (puertas cerradas dinámicas)
	tx, ty := blockinfo.WorldToTile(float64(toX), float64(toY))
	if m.doors != nil && m.doors.IsDoorMoveBlocked(tx, ty) {
		if p.Client != nil {
			failPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVEFAIL_CMD)
			failPkt.WriteInt32(int32(fromX))
			failPkt.WriteInt32(int32(fromY))
			_ = p.Client.SendPacket(failPkt)
		}
		return
	}

	p.Char.XPos = toX
	p.Char.YPos = toY

	op := protocol.PROTO_NC_ACT_SOMEONEMOVERUN_CMD
	speed := uint16(115)
	if p.MoverHandle > 0 && p.MoverSpeed > 0 {
		speed = p.MoverSpeed
	} else if walk {
		op = protocol.PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
		speed = 50
	}

	movePw := protocol.NewPacketWriter(op)
	movePw.WriteUint16(p.MapObjectID) // 0x00: handle (2 bytes)
	movePw.WriteInt32(int32(fromX))   // 0x02: from.x (4 bytes)
	movePw.WriteInt32(int32(fromY))   // 0x06: from.y (4 bytes)
	movePw.WriteInt32(int32(toX))     // 0x0A: to.x (4 bytes)
	movePw.WriteInt32(int32(toY))     // 0x0E: to.y (4 bytes)
	movePw.WriteUint16(speed)         // 0x12: speed (2 bytes)
	movePw.WriteUint16(0)             // 0x14: moveattr (2 bytes, direct=0)

	// 3. Actualización de posición en SpatialGrid y detección de cruce de sectores AOI
	crossCell, entered, left := m.grid.UpdateEntityPosition(p, toX, toY)

	if crossCell {
		// A. Nuevas entidades que entraron al rango de visión
		enteredEnts := m.grid.GetEntitiesInCells(entered)
		for _, e := range enteredEnts {
			if e.GetHandle() == p.MapObjectID {
				continue
			}
			if other, ok := e.(*ZonePlayer); ok {
				spawnOther := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
				other.WriteCharacterDisplay(spawnOther)
				_ = p.Client.SendPacket(spawnOther)

				if other.Client != nil && !other.HideExceptMe {
					spawnSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
					p.WriteCharacterDisplay(spawnSelf)
					_ = other.Client.SendPacket(spawnSelf)
				}
			} else if di, ok := e.(*ZoneDroppedItem); ok {
				pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
				pw.WriteUint16(di.Handle)
				pw.WriteUint16(di.ItemID)
				pw.WriteInt32(di.X)
				pw.WriteInt32(di.Y)
				pw.WriteUint16(di.DropMobHandle)
				pw.WriteByte(di.Attr)
				_ = p.Client.SendPacket(pw)
			}
		}

		// B. Entidades que salieron del rango de visión
		leftEnts := m.grid.GetEntitiesInCells(left)
		for _, e := range leftEnts {
			if e.GetHandle() == p.MapObjectID {
				continue
			}
			if other, ok := e.(*ZonePlayer); ok {
				delOther := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
				delOther.WriteUint16(other.MapObjectID)
				_ = p.Client.SendPacket(delOther)

				if other.Client != nil {
					delSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
					delSelf.WriteUint16(p.MapObjectID)
					_ = other.Client.SendPacket(delSelf)
				}
			} else {
				delPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
				delPkt.WriteUint16(e.GetHandle())
				_ = p.Client.SendPacket(delPkt)
			}
		}

		// C. Jugadores en celdas superpuestas reciben el movimiento
		surrounding := m.grid.GetSurroundingPlayers(toX, toY)
		for _, other := range surrounding {
			if other.MapObjectID == p.MapObjectID || other.Client == nil || other.HideExceptMe {
				continue
			}
			_ = other.Client.SendPacket(movePw)
		}
	} else {
		// Difusión AOI a los jugadores en la vecindad 3x3
		m.grid.BroadcastAOI(movePw, toX, toY, p.MapObjectID)
	}
}

// StopPlayer updates position and broadcasts stop to AOI.
func (m *ZoneMap) StopPlayer(p *ZonePlayer, x, y int) {
	p.Char.XPos = x
	p.Char.YPos = y

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESTOP_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteInt32(int32(x))
	pw.WriteInt32(int32(y))
	m.grid.BroadcastAOI(pw, x, y, p.MapObjectID)
}

// JumpPlayer broadcasts player jump to AOI.
func (m *ZoneMap) JumpPlayer(p *ZonePlayer) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEEONEJUMP_CMD)
	pw.WriteUint16(p.MapObjectID)
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, p.MapObjectID)
}

// EmoticonPlayer transmite PROTO_NC_ACT_SOMEONEEMOTICON_CMD (0x2021) de 3 bytes a la vecindad AOI.
func (m *ZoneMap) EmoticonPlayer(p *ZonePlayer, emoticonID byte) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEEMOTICON_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(emoticonID)
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, 0)
}

// EmoticonStopPlayer transmite PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD (0x2023) de 5 bytes a la vecindad AOI.
func (m *ZoneMap) EmoticonStopPlayer(p *ZonePlayer, emoticonID byte, emoticonFrame uint16) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(emoticonID)
	pw.WriteUint16(emoticonFrame)
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, 0)
}

// NormalChat broadcasts chat to players in the AOI neighborhood.
func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONECHAT_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(byte(len(text)))
	pw.WriteByte(color)
	pw.WriteString(text, len(text))
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, 0)
}

// PitchTent activa el modo de descanso en carpa/tienda para el jugador y difunde a la vecindad AOI.
func (m *ZoneMap) PitchTent(p *ZonePlayer, tentID uint16) {
	p.State = 4 // Resting
	p.TentID = tentID

	m.mu.Lock()
	if m.tents == nil {
		m.tents = make(map[uint16]uint16)
	}
	m.tents[p.MapObjectID] = tentID
	m.mu.Unlock()

	// PROTO_NC_ACT_SOMEONEPITCHTENT_CMD (0x2029, 14 bytes)
	// uint16 handle, CHARBRIEFINFO_CAMP camp (12 bytes)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEPITCHTENT_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteUint16(tentID)
	pw.WriteUint32(0xFFFFFFFF)
	pw.WriteUint32(0xFFFFFFFF)
	pw.WriteUint16(0xFFFF)
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, 0)
}

// FoldTent desarma la carpa, restablece el estado normal y difunde la apariencia visual a la vecindad AOI.
func (m *ZoneMap) FoldTent(p *ZonePlayer) {
	p.State = 1 // Normal
	p.TentID = 0

	m.mu.Lock()
	delete(m.tents, p.MapObjectID)
	m.mu.Unlock()

	// PROTO_NC_ACT_SOMEONEFOLDTENT_CMD (0x202C, 45 bytes)
	// uint16 handle, CHARBRIEFINFO_NOTCAMP shape (43 bytes = WriteEquipment)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEFOLDTENT_CMD)
	pw.WriteUint16(p.MapObjectID)
	p.WriteEquipment(pw)
	m.grid.BroadcastAOI(pw, p.Char.XPos, p.Char.YPos, 0)
}
