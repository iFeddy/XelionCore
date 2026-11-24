package zone

import (
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
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

// WriteRegenMob serializa la entrada oficial de 149 bytes de PROTO_NC_BRIEFINFO_REGENMOB_CMD / MOB_CMD para un NPC o Gate.
// Verificado en Zone.exe ShineNPC::so_SetGateBriefInform (0x00557BC0) y ShineMob::so_mob_Regenerate (0x004B41B0).
func (n *ZoneNPC) WriteRegenMob(pw *protocol.PacketWriter) {
	pw.WriteUint16(n.Handle)               // 0x00: handle (2B)
	pw.WriteByte(0)                        // 0x02: mode (1B, 0=vivo/normal en so_mob_Regenerate)
	pw.WriteUint16(n.MobID)                // 0x03: mobid (2B)
	pw.WriteInt32(n.X)                     // 0x05: coord.x (4B)
	pw.WriteInt32(n.Y)                     // 0x09: coord.y (4B)
	pw.WriteByte(encodeCoordDir(n.Direct)) // 0x0D: coord.dir (1B, direct/2)
	if n.IsGate {
		pw.WriteByte(1)                   // 0x0E: flagstate = 1 (gate)
		pw.WriteString(n.GateToWhere, 12) // 0x0F: gate2where (12B)
		pw.Fill(87, 0)                    // 0x1B: resto de unión flag[99]
	} else {
		pw.WriteByte(0) // 0x0E: flagstate = 0
		pw.Fill(99, 0)  // 0x0F: abstatebit (99B)
	}
	pw.WriteString("", 32) // 0x72: sAnimation (32B)
	pw.WriteByte(0)        // 0x92: nAnimationLevel (1B)
	pw.WriteByte(0)        // 0x93: nKQTeamType (1B)
	pw.WriteByte(0)        // 0x94: bRegenAni (1B)
}

// WriteDropEntry serializa las 15 bytes oficiales de PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (0x1C0A) y entradas de ITEMONFIELD_CMD (0x1C0B).
func (d *ZoneDroppedItem) WriteDropEntry(pw *protocol.PacketWriter) {
	pw.WriteUint16(d.Handle)        // 0x00: handle (2B)
	pw.WriteUint16(d.ItemID)        // 0x02: itemid (2B)
	pw.WriteInt32(d.X)              // 0x04: location.x (4B)
	pw.WriteInt32(d.Y)              // 0x08: location.y (4B)
	pw.WriteUint16(d.DropMobHandle) // 0x0C: dropmobhandle (2B)
	pw.WriteByte(d.Attr)            // 0x0E: attr (1B)
}

// ZoneMoverBrief modela una entrada PROTO_NC_BRIEFINFO_REGENMOVER_CMD (139B) para lotes PROTO_NC_BRIEFINFO_MOVER_CMD (0x1C1B).
type ZoneMoverBrief struct {
	Handle     uint16
	ID         uint32
	HP         uint32
	X          int32
	Y          int32
	Direct     int16
	AbstateBit [99]byte
	Grade      uint8
	SlotHandle [10]uint16
}

// WriteRegenMover serializa 139 bytes exactos de PROTO_NC_BRIEFINFO_REGENMOVER_CMD (Zone.exe bixmv_Send 0x0043E0D0).
func (mv *ZoneMoverBrief) WriteRegenMover(pw *protocol.PacketWriter) {
	pw.WriteUint16(mv.Handle)
	pw.WriteUint32(mv.ID)
	pw.WriteUint32(mv.HP)
	pw.WriteInt32(mv.X)
	pw.WriteInt32(mv.Y)
	pw.WriteByte(encodeCoordDir(mv.Direct))
	pw.WriteBytes(mv.AbstateBit[:])
	pw.WriteByte(mv.Grade)
	for i := 0; i < 10; i++ {
		pw.WriteUint16(mv.SlotHandle[i])
	}
}

// ZoneMagicFieldBrief modela una entrada PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD (16B) para MAGICFIELDINFO_CMD (0x1C0D).
type ZoneMagicFieldBrief struct {
	Handle  uint16
	Caster  uint16
	SkillID uint16
	X       int32
	Y       int32
	Radius  uint16
}

// WriteMagicFieldSpread serializa 16 bytes exactos de PROTO_NC_BRIEFINFO_MAGICFIELDSPREAD_CMD.
func (mf *ZoneMagicFieldBrief) WriteMagicFieldSpread(pw *protocol.PacketWriter) {
	pw.WriteUint16(mf.Handle)
	pw.WriteUint16(mf.Caster)
	pw.WriteUint16(mf.SkillID)
	pw.WriteInt32(mf.X)
	pw.WriteInt32(mf.Y)
	pw.WriteUint16(mf.Radius)
}

// ZoneDoorBrief modela una entrada PROTO_NC_BRIEFINFO_BUILDDOOR_CMD (48B) para PROTO_NC_BRIEFINFO_DOOR_CMD (0x1C10).
type ZoneDoorBrief struct {
	Handle     uint16
	MobID      uint16
	X          int32
	Y          int32
	Direct     int16
	DoorState  byte
	BlockIndex string
	Scale      uint16
}

// WriteBuildDoor serializa 48 bytes exactos de PROTO_NC_BRIEFINFO_BUILDDOOR_CMD (Zone.exe 0BriefInformationDoor 0x00549260).
func (d *ZoneDoorBrief) WriteBuildDoor(pw *protocol.PacketWriter) {
	scale := d.Scale
	if scale == 0 {
		scale = 1000 // 0x3e8 default en 0BriefInformationDoor (0x005492B6)
	}
	pw.WriteUint16(d.Handle)
	pw.WriteUint16(d.MobID)
	pw.WriteInt32(d.X)
	pw.WriteInt32(d.Y)
	pw.WriteByte(encodeCoordDir(d.Direct))
	pw.WriteByte(d.DoorState)
	pw.WriteString(d.BlockIndex, 32)
	pw.WriteUint16(scale)
}

// ZoneDroppedItem representa un ítem tirado en el suelo del mapa según PROTO_NC_BRIEFINFO_DROPEDITEM_CMD (0x1C0A).
type ZoneDroppedItem struct {
	Handle         uint16
	ItemID         uint16
	Lot            uint32
	Flags          int
	X              int32
	Y              int32
	DropMobHandle  uint16
	Attr           uint8
	OwnerCharID    int
	OwnerExpiresAt time.Time
	AllCanLoot     bool
	IsCenDrop      bool
	StatOptions    []item.ItemStatOption
	CreatedAt      time.Time
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
	players          map[uint16]*ZonePlayer
	npcs             map[uint16]*ZoneNPC
	mobs             map[uint16]*ZoneMonster
	monstersSnapshot []*ZoneMonster
	monstersDirty    bool
	gates            []*ZoneNPC
	droppedItems     map[uint16]*ZoneDroppedItem
	tents            map[uint16]uint16 // MapObjectID -> tentID
	grid             *SpatialGrid
	blockInfo        *blockinfo.MapBlockInformation
	doors            *blockinfo.MapDoorArray
	gameLoop         *MapGameLoop
	FieldEntry       *field.FieldEntry
	dataStore        *DataStore
	handler          *Handler
	mobCatalog       *MobCatalog
	regenGroups      map[string]*MobRegenGroup
	luaEngine        *lua.Engine
	loadedAIScripts  map[string]bool
	worldBridge      *ZoneWorldBridge
	aiTimersSec      map[int]int64

	// Campos de subsistemas de la Fase 15 (Instancias Lua, KQ, Guild Tournament y Bodas)
	IsInstance         bool
	InstanceKind       string // "ID", "KQ", "GT", "Wedding"
	InstanceDefID      uint16
	InstanceLimitSec   int
	InstanceStartedAt  time.Time
	ScenarioInterrupts map[string]string // EventName -> Lua callback
	FieldScriptFuncs   map[string]string // EventName ("MapLogin", etc.) -> Lua func name
	AIScriptFuncs      map[uint32]map[string]string
	VanishTimers       map[string]int64
	NPCInteractCB      map[uint32]string // npcHandle -> Lua callback
	KQScoreRed         int
	KQScoreBlue        int
	KQWinRed           int
	KQWinBlue          int
	KQTeamByChar       map[int]int
	KQHandle           uint32 // FieldMap+0x3C (exe: 0x004980C3 / 0x00497383, 0xFFFFFFFF = unbound)
	KQID               uint16 // FieldMap+0x40 (exe: 0x004980C6 / 0x00497386, 0xFFFF = unbound)
	KQScriptActive     bool   // CinemaComplex active film count > 0 (exe: 0x0043FC35)
	KQScriptName       string // Active KQ script name for cc_DropFilm / cc_PlayFilm
	KQInitValue        string // Active KQ script init value
	kqScriptTable      string // Tabla global Lua para scripts KQ migrados (ej. KQConditionOfHero)
	GTState            *GuildTournamentMatch
	WeddingState       *WeddingCeremonyState
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
		mobs:         make(map[uint16]*ZoneMonster),
		droppedItems: make(map[uint16]*ZoneDroppedItem),
		tents:        make(map[uint16]uint16),
		grid:         NewSpatialGrid(DefaultCellSize),
		KQHandle:     kqInvalidHandle,
		KQID:         kqInvalidID,
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

// GameLoop retorna el bucle de simulación del mapa.
func (m *ZoneMap) GameLoop() *MapGameLoop {
	return m.gameLoop
}

// Close finaliza el bucle de juego y libera recursos del mapa.
func (m *ZoneMap) Close() {
	if m.gameLoop != nil {
		m.gameLoop.Stop()
	}
	m.mu.Lock()
	if m.luaEngine != nil {
		m.luaEngine.Close()
		m.luaEngine = nil
	}
	m.mu.Unlock()
}

// MapManager manages active ZoneMap instances dynamically mapped from data/World/Map definitions.
type MapManager struct {
	mu         sync.RWMutex
	fieldStore *field.Store
	npcStore   *npc.Store
	blockStore *blockinfo.Store
	dataStore  *DataStore
	handler    *Handler
	mobMap     map[string]uint16 // InxName -> MobID
	maps       map[string]*ZoneMap
	zoneFilter int // -1 carga todos los mapas; >=0 filtra mapas pertenecientes a Zone0X
}

// NewMapManager initializes a new dynamic MapManager.
func NewMapManager(fieldStore *field.Store) *MapManager {
	return &MapManager{
		fieldStore: fieldStore,
		mobMap:     make(map[string]uint16),
		maps:       make(map[string]*ZoneMap),
		zoneFilter: -1,
	}
}

// SetZoneFilter asigna el número de nodo Zone00..Zone04 (-1 = todos los mapas en modo single-node).
func (mm *MapManager) SetZoneFilter(zoneFilter int) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.zoneFilter = zoneFilter
}

// ZoneFilter devuelve el filtro de nodo Zone activo (-1 en modo single-node).
func (mm *MapManager) ZoneFilter() int {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	return mm.zoneFilter
}

// OwnsMap indica si mapID pertenece a este proceso Zone según FieldList.Fiesta y InstanceDungeon.ZoneNumber.
func (mm *MapManager) OwnsMap(mapID string) bool {
	mm.mu.RLock()
	zf := mm.zoneFilter
	fs := mm.fieldStore
	mm.mu.RUnlock()
	if zf < 0 || fs == nil {
		return true
	}
	norm := strings.TrimSpace(mapID)
	if norm == "" {
		return false
	}
	if primary, ok := fs.GetField(norm); ok && primary != nil {
		fiestaNode := primary.Fiesta
		if fiestaNode == 99 {
			fiestaNode = 0
		}
		if fiestaNode == zf {
			return true
		}
		return false
	}
	for _, d := range fs.GetDungeonsByZone(zf) {
		if strings.EqualFold(d.MapID, norm) {
			return true
		}
	}
	return false
}

// PreloadZoneMaps pre-inicializa en RAM los mapas asignados al filtro de nodo activo.
func (mm *MapManager) PreloadZoneMaps() int {
	mm.mu.RLock()
	fs := mm.fieldStore
	zf := mm.zoneFilter
	mm.mu.RUnlock()
	if fs == nil {
		return 0
	}
	seen := make(map[string]struct{})
	loaded := 0
	for _, f := range fs.GetAllFields() {
		fiestaNode := f.Fiesta
		if fiestaNode == 99 {
			fiestaNode = 0
		}
		if zf >= 0 && fiestaNode != zf {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(f.MapID))
		if k == "" {
			continue
		}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		mm.GetOrCreateMap(f.MapID)
		loaded++
	}
	if zf >= 0 {
		for _, d := range fs.GetDungeonsByZone(zf) {
			k := strings.ToLower(strings.TrimSpace(d.MapID))
			if k == "" {
				continue
			}
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			mm.GetOrCreateMap(d.MapID)
			loaded++
		}
	}
	return loaded
}

// GetAllMaps retorna un snapshot de todos los mapas activos administrados por este MapManager.
func (mm *MapManager) GetAllMaps() []*ZoneMap {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	res := make([]*ZoneMap, 0, len(mm.maps))
	for _, m := range mm.maps {
		if m != nil {
			res = append(res, m)
		}
	}
	return res
}

// StopAll detiene todos los MapGameLoop activos y libera los mapas administrados.
func (mm *MapManager) StopAll() {
	mm.mu.Lock()
	maps := make([]*ZoneMap, 0, len(mm.maps))
	for _, m := range mm.maps {
		if m != nil {
			maps = append(maps, m)
		}
	}
	mm.maps = make(map[string]*ZoneMap)
	mm.mu.Unlock()

	for _, m := range maps {
		if gl := m.GameLoop(); gl != nil {
			gl.Stop()
		}
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
	if mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetField(normID); ok {
			zMap.FieldEntry = entry
		}
	}

	if mm.blockStore != nil {
		mbi := mm.blockStore.GetSHBD(normID)
		doors := mm.blockStore.NewInstanceDoors(normID)
		zMap.SetCollisionData(mbi, doors)
	}

	mm.maps[key] = zMap
	ds := mm.dataStore
	handler := mm.handler
	log.Printf("[Zone] Initialized dynamic map '%s' (Serial: %d, Name: '%s', Size: %dx%d)",
		normID, serial, name, zMap.Width, zMap.Height)

	mm.loadMapNPCs(zMap)
	if handler != nil {
		zMap.SetHandler(handler)
	}
	if ds != nil {
		zMap.SetDataStore(ds)
	}
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
	var mobsInAOI []*ZoneMonster
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
		case GridEntityMob:
			if zm, ok := ent.(*ZoneMonster); ok && !zm.IsDead {
				mobsInAOI = append(mobsInAOI, zm)
			}
		case GridEntityItem:
			if zi, ok := ent.(*ZoneDroppedItem); ok {
				itemsInAOI = append(itemsInAOI, zi)
			}
		}
	}

	// 1. Enviar jugadores en AOI al jugador entrante (lotes de máx 28 según bixc_Exchange 0x0043FD00)
	if p.Client != nil {
		for i := 0; i < len(otherPlayers); i += BriefInfoMaxBatch {
			end := i + BriefInfoMaxBatch
			if end > len(otherPlayers) {
				end = len(otherPlayers)
			}
			chunk := otherPlayers[i:end]
			spawnOthers := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHARACTER_CMD)
			spawnOthers.WriteByte(byte(len(chunk)))
			for _, other := range chunk {
				other.WriteCharacterDisplay(spawnOthers)
			}
			_ = p.Client.SendPacket(spawnOthers)
		}
	}

	// 2. Enviar jugador entrante a los otros jugadores en AOI
	spawnSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
	p.WriteCharacterDisplay(spawnSelf)
	for _, other := range otherPlayers {
		if other.Client != nil && !other.HideExceptMe {
			_ = other.Client.SendPacket(spawnSelf)
		}
	}

	// 3. Enviar NPCs, Gates y Monstruos vivos del mapa en AOI
	m.SendNPCs(p)
	if len(mobsInAOI) > 0 {
		m.SendMobs(p, mobsInAOI)
	}

	// 4. Enviar ítems en AOI al jugador (lotes de máx 28 según bixi_Exchange 0x004401A0)
	for i := 0; i < len(itemsInAOI); i += BriefInfoMaxBatch {
		end := i + BriefInfoMaxBatch
		if end > len(itemsInAOI) {
			end = len(itemsInAOI)
		}
		chunk := itemsInAOI[i:end]
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD)
		pw.WriteByte(byte(len(chunk)))
		for _, di := range chunk {
			di.WriteDropEntry(pw)
		}
		_ = p.Client.SendPacket(pw)
	}

	// 5. Enviar sincronización de HP
	if p.Client != nil {
		hpUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpUpdate.WriteUint32(uint32(p.Char.CurHP))
		hpUpdate.WriteUint16(p.NextHPChangeOrder())
		_ = p.Client.SendPacket(hpUpdate)
	}

	// 6. Hook oficial FieldScript "MapLogin" en instancias Lua (Zone.exe cFieldScriptFunc)
	m.NotifyPlayerMapLogin(p)
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
	di.WriteDropEntry(pw)
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

// CleanupExpiredDrops libera la propiedad temporal vencida y elimina ítems del suelo cuya antigüedad supera maxAge.
func (m *ZoneMap) CleanupExpiredDrops(maxAge time.Duration) {
	m.mu.Lock()
	var expired []*ZoneDroppedItem
	now := time.Now()
	for h, di := range m.droppedItems {
		if di.OwnerCharID != 0 && !di.OwnerExpiresAt.IsZero() && now.After(di.OwnerExpiresAt) {
			di.OwnerCharID = 0
			di.AllCanLoot = true
		}
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

// AddMonster registra un monstruo en el mapa y en la grilla espacial AOI.
func (m *ZoneMap) AddMonster(mob *ZoneMonster) {
	if mob == nil {
		return
	}
	m.mu.Lock()
	if mob.Handle == 0 {
		mob.Handle = m.nextObjID
		m.nextObjID++
	}
	m.mobs[mob.Handle] = mob
	m.monstersDirty = true
	m.mu.Unlock()

	if m.grid != nil {
		m.grid.AddEntity(mob)
	}
}

// RemoveMonster remueve un monstruo del mapa y de la grilla AOI.
func (m *ZoneMap) RemoveMonster(handle uint16) {
	m.mu.Lock()
	_, ok := m.mobs[handle]
	if ok {
		delete(m.mobs, handle)
		m.monstersDirty = true
	}
	m.mu.Unlock()

	if ok && m.grid != nil {
		m.grid.RemoveEntity(handle)
	}
}

// GetPlayer recupera un jugador por su MapObjectID.
func (m *ZoneMap) GetPlayer(handle uint16) *ZonePlayer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.players[handle]
}

// GetMonster recupera un monstruo por su handle de mapa.
func (m *ZoneMap) GetMonster(handle uint16) *ZoneMonster {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.mobs[handle]
}

// GetAllMonsters retorna una instantánea inmutable cacheada de todos los monstruos en el mapa.
// En estado estable (sin spawns/despawns) opera en O(1) con 0 alocaciones por tick.
func (m *ZoneMap) GetAllMonsters() []*ZoneMonster {
	m.mu.RLock()
	if !m.monstersDirty {
		snap := m.monstersSnapshot
		m.mu.RUnlock()
		return snap
	}
	m.mu.RUnlock()

	m.mu.Lock()
	if m.monstersDirty {
		if len(m.mobs) == 0 {
			m.monstersSnapshot = nil
		} else {
			snap := make([]*ZoneMonster, 0, len(m.mobs))
			for _, mob := range m.mobs {
				snap = append(snap, mob)
			}
			m.monstersSnapshot = snap
		}
		m.monstersDirty = false
	}
	snap := m.monstersSnapshot
	m.mu.Unlock()
	return snap
}

// GetEntity recupera cualquier entidad activa (jugador, monstruo, NPC) por su handle.
func (m *ZoneMap) GetEntity(handle uint16) GridEntity {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if p, ok := m.players[handle]; ok {
		return p
	}
	if mob, ok := m.mobs[handle]; ok {
		return mob
	}
	if n, ok := m.npcs[handle]; ok {
		return n
	}
	if di, ok := m.droppedItems[handle]; ok {
		return di
	}
	return nil
}

// SendDroppedItems transmits all dropped items on the field to a joining player using PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD (0x1C0B), máx 28 por paquete.
func (m *ZoneMap) SendDroppedItems(p *ZonePlayer) {
	m.mu.RLock()
	items := make([]*ZoneDroppedItem, 0, len(m.droppedItems))
	for _, di := range m.droppedItems {
		items = append(items, di)
	}
	m.mu.RUnlock()

	if len(items) == 0 || p == nil || p.Client == nil {
		return
	}

	for i := 0; i < len(items); i += BriefInfoMaxBatch {
		end := i + BriefInfoMaxBatch
		if end > len(items) {
			end = len(items)
		}
		chunk := items[i:end]
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD)
		pw.WriteByte(byte(len(chunk)))
		for _, di := range chunk {
			di.WriteDropEntry(pw)
		}
		_ = p.Client.SendPacket(pw)
	}
}

// SendNPCs transmits all map NPCs and Gates to a player using official PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09), máx 28 por paquete.
func (m *ZoneMap) SendNPCs(p *ZonePlayer) {
	m.mu.RLock()
	npcList := make([]*ZoneNPC, 0, len(m.npcs))
	for _, npc := range m.npcs {
		npcList = append(npcList, npc)
	}
	m.mu.RUnlock()

	if len(npcList) == 0 || p == nil || p.Client == nil {
		return
	}

	for i := 0; i < len(npcList); i += BriefInfoMaxBatch {
		end := i + BriefInfoMaxBatch
		if end > len(npcList) {
			end = len(npcList)
		}
		chunk := npcList[i:end]

		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_MOB_CMD)
		pw.WriteByte(byte(len(chunk))) // mobnum (1 byte)
		for _, npc := range chunk {
			npc.WriteRegenMob(pw)
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

func (m *ZoneMap) sendMoveFail(p *ZonePlayer) {
	if p == nil || p.Client == nil || p.Char == nil {
		return
	}
	failPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_MOVEFAIL_CMD)
	failPkt.WriteInt32(int32(p.Char.XPos))
	failPkt.WriteInt32(int32(p.Char.YPos))
	_ = p.Client.SendPacket(failPkt)
}

func (m *ZoneMap) validatePlayerMove(p *ZonePlayer, fromX, fromY, toX, toY int) bool {
	if p == nil || p.Char == nil {
		return false
	}
	dx := p.Char.XPos - fromX
	dy := p.Char.YPos - fromY
	if dx*dx+dy*dy > 0x9c40 {
		return false
	}
	if m.blockInfo != nil {
		if !m.blockInfo.OccupCellCheck(float64(toX), float64(toY), 1) {
			return false
		}
		if !m.blockInfo.OccupCellCheck(float64(p.Char.XPos), float64(p.Char.YPos), 1) {
			return false
		}
		if m.blockInfo.IsMovePathBlocked(p.Char.XPos, p.Char.YPos, toX, toY, 2) {
			return false
		}
	}
	if m.doors != nil {
		tx, ty := blockinfo.WorldToTile(float64(toX), float64(toY))
		if m.doors.IsDoorMoveBlocked(tx, ty) {
			return false
		}
	}
	return true
}

func (m *ZoneMap) syncMoveAOI(p *ZonePlayer, toX, toY int, movePw *protocol.PacketWriter) {
	crossCell, entered, left := m.grid.UpdateEntityPosition(p, toX, toY)

	if crossCell {
		// A. Nuevas entidades que entraron al rango de visión
		enteredEnts := m.grid.GetEntitiesInCells(entered)
		var enteredMobs []*ZoneMonster
		for _, e := range enteredEnts {
			if e.GetHandle() == p.MapObjectID {
				continue
			}
			if other, ok := e.(*ZonePlayer); ok && p.Client != nil {
				spawnOther := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
				other.WriteCharacterDisplay(spawnOther)
				_ = p.Client.SendPacket(spawnOther)

				if other.Client != nil && !other.HideExceptMe {
					spawnSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
					p.WriteCharacterDisplay(spawnSelf)
					_ = other.Client.SendPacket(spawnSelf)
				}
			} else if mob, ok := e.(*ZoneMonster); ok && !mob.IsDead {
				enteredMobs = append(enteredMobs, mob)
			} else if di, ok := e.(*ZoneDroppedItem); ok && p.Client != nil {
				pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
				di.WriteDropEntry(pw)
				_ = p.Client.SendPacket(pw)
			}
		}
		if len(enteredMobs) > 0 {
			m.SendMobs(p, enteredMobs)
		}

		// B. Entidades que salieron del rango de visión
		leftEnts := m.grid.GetEntitiesInCells(left)
		for _, e := range leftEnts {
			if e.GetHandle() == p.MapObjectID {
				continue
			}
			if other, ok := e.(*ZonePlayer); ok {
				if p.Client != nil {
					delOther := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
					delOther.WriteUint16(other.MapObjectID)
					_ = p.Client.SendPacket(delOther)
				}
				if other.Client != nil {
					delSelf := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
					delSelf.WriteUint16(p.MapObjectID)
					_ = other.Client.SendPacket(delSelf)
				}
			} else if p.Client != nil {
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

// MovePlayer actualiza la posición del jugador para sp_NC_ACT_WALK_REQ (0x2003) / sp_NC_ACT_RUN_REQ (0x2005)
// y emite PROTO_NC_ACT_SOMEONEWALK_CMD (0x2004) / PROTO_NC_ACT_SOMEONERUN_CMD (0x2006).
func (m *ZoneMap) MovePlayer(p *ZonePlayer, oldX, oldY, newX, newY int, walk bool) {
	if !m.validatePlayerMove(p, oldX, oldY, newX, newY) {
		m.sendMoveFail(p)
		return
	}

	srvFromX, srvFromY := p.Char.XPos, p.Char.YPos
	p.Char.XPos = newX
	p.Char.YPos = newY

	op := protocol.PROTO_NC_ACT_SOMEONERUN_CMD
	if walk {
		op = protocol.PROTO_NC_ACT_SOMEONEWALK_CMD
	}
	movePw := protocol.NewPacketWriter(op)
	movePw.WriteUint16(p.MapObjectID)
	movePw.WriteInt32(int32(srvFromX))
	movePw.WriteInt32(int32(srvFromY))
	movePw.WriteInt32(int32(newX))
	movePw.WriteInt32(int32(newY))
	m.syncMoveAOI(p, newX, newY, movePw)
}

// MovePlayerCmd valida colisiones contra .shbd y .sbi según Zone.exe (0x00457C60 / 0x00457FD0),
// actualiza la grilla AOI y transmite PROTO_NC_ACT_SOMEONEMOVEWALK_CMD / PROTO_NC_ACT_SOMEONEMOVERUN_CMD.
func (m *ZoneMap) MovePlayerCmd(p *ZonePlayer, fromX, fromY, toX, toY int, walk bool) {
	if !m.validatePlayerMove(p, fromX, fromY, toX, toY) {
		m.sendMoveFail(p)
		return
	}

	srvFromX, srvFromY := p.Char.XPos, p.Char.YPos
	p.Char.XPos = toX
	p.Char.YPos = toY

	op := protocol.PROTO_NC_ACT_SOMEONEMOVERUN_CMD
	if walk {
		op = protocol.PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
	}
	speed := p.GetMoveSpeed(walk)

	movePw := protocol.NewPacketWriter(op)
	movePw.WriteUint16(p.MapObjectID)  // 0x00: handle (2 bytes)
	movePw.WriteInt32(int32(srvFromX)) // 0x02: from.x (4 bytes)
	movePw.WriteInt32(int32(srvFromY)) // 0x06: from.y (4 bytes)
	movePw.WriteInt32(int32(toX))      // 0x0A: to.x (4 bytes)
	movePw.WriteInt32(int32(toY))      // 0x0E: to.y (4 bytes)
	movePw.WriteUint16(speed)          // 0x12: speed (2 bytes)
	movePw.WriteUint16(0)              // 0x14: moveattr (2 bytes, bit0=0)
	m.syncMoveAOI(p, toX, toY, movePw)
}

// StopPlayer updates position and broadcasts stop to AOI according to sp_NC_ACT_STOP_REQ (0x00456FF0).
func (m *ZoneMap) StopPlayer(p *ZonePlayer, x, y int) {
	if p == nil || p.Char == nil {
		return
	}
	dx := p.Char.XPos - x
	dy := p.Char.YPos - y
	if dx*dx+dy*dy > 0x9c40 {
		m.sendMoveFail(p)
		return
	}
	if m.blockInfo != nil && !m.blockInfo.OccupCellCheck(float64(x), float64(y), 1) {
		m.sendMoveFail(p)
		return
	}
	if m.doors != nil {
		tx, ty := blockinfo.WorldToTile(float64(x), float64(y))
		if m.doors.IsDoorMoveBlocked(tx, ty) {
			m.sendMoveFail(p)
			return
		}
	}

	p.Char.XPos = x
	p.Char.YPos = y

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESTOP_CMD)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteInt32(int32(x))
	pw.WriteInt32(int32(y))
	m.syncMoveAOI(p, x, y, pw)
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
// NormalChat emite PROTO_NC_ACT_SOMEONECHAT_CMD sin enlaces de ítems.
func (m *ZoneMap) NormalChat(p *ZonePlayer, text string, color byte) {
	m.NormalChatWithLinks(p, text, color, 0, nil)
}

// NormalChatWithLinks emite PROTO_NC_ACT_SOMEONECHAT_CMD (0x2002) con cabecera oficial de 7 bytes
// verificada en Zone.exe (so_Chat @ 0x00458C90) + content + itemLinkData (count * 103B).
func (m *ZoneMap) NormalChatWithLinks(p *ZonePlayer, text string, color byte, linkCount byte, linkData []byte) {
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONECHAT_CMD)
	var flag byte = 0x02 // bit 1 = chatwin
	if color == 0x03 || isZoneGM(p) {
		flag |= 0x01 // bit 0 = GMColor
	}
	pw.WriteByte(linkCount)
	pw.WriteUint16(p.MapObjectID)
	pw.WriteByte(byte(len(text)))
	pw.WriteByte(flag)
	pw.WriteByte(0) // nChatFontColorID
	pw.WriteByte(0) // nChatBalloonColorID
	pw.WriteString(text, len(text))
	if linkCount > 0 && len(linkData) > 0 {
		pw.WriteBytes(linkData)
	}
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
