package world

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// KQParticipant representa a un jugador registrado en una Kingdom Quest.
type KQParticipant struct {
	CharID   int
	CharName string
	Level    byte
	Class    byte
	TeamType byte
	Client   *network.Client
}

// KQInstance representa el estado en memoria de una Kingdom Quest activa o programada.
type KQInstance struct {
	Handle        uint32
	ID            uint16
	Title         string
	Status        byte // 0=Scheduled, 1=Recruiting, 2=Running
	StartTime     time.Time
	LimitTime     uint16 // seconds
	StartWaitTime uint16 // seconds
	MinLevel      byte
	MaxLevel      byte
	MinPlayers    uint16
	MaxPlayers    uint16
	Participants  map[int]*KQParticipant // CharID -> Participant
}

// KQManager administra el ciclo de vida, inscripción y difusión de Kingdom Quests (Dept 22).
type KQManager struct {
	mu            sync.RWMutex
	instances     map[uint32]*KQInstance
	charToKQ      map[int]uint32 // CharID -> Handle
	onlinePlayers map[int]*network.Client
	nextHandle    uint32
}

// NewKQManager inicializa el gestor de Kingdom Quests a partir de KingdomQuest.shn si está disponible.
func NewKQManager(kqSHN *shn.File) *KQManager {
	km := &KQManager{
		instances:     make(map[uint32]*KQInstance),
		charToKQ:      make(map[int]uint32),
		onlinePlayers: make(map[int]*network.Client),
		nextHandle:    1,
	}

	if kqSHN != nil && kqSHN.RecordCount > 0 {
		km.loadFromSHN(kqSHN)
	} else {
		km.loadDefaults()
	}

	return km
}

func (km *KQManager) loadFromSHN(f *shn.File) {
	now := time.Now()
	for i := 0; i < f.RecordCount; i++ {
		row := f.Rows[i]
		id := row.GetUint16("ID")
		title := strings.TrimSpace(row.GetString("Title"))
		if title == "" {
			continue
		}

		minLev := row.GetByte("MinLevel")
		maxLev := row.GetByte("MaxLevel")
		minPly := row.GetUint16("MinPlayers")
		maxPly := row.GetUint16("MaxPlayers")
		limitTime := row.GetUint16("LimitTime")
		startWait := row.GetUint16("StartWaitTime")

		if minPly == 0 {
			minPly = 5
		}
		if maxPly == 0 {
			maxPly = 15
		}
		if limitTime == 0 {
			limitTime = 3600
		}
		if startWait == 0 {
			startWait = 300
		}

		handle := km.nextHandle
		km.nextHandle++

		km.instances[handle] = &KQInstance{
			Handle:        handle,
			ID:            id,
			Title:         title,
			Status:        1, // Reclutando inicialmente para disponibilidad en vivo
			StartTime:     now.Add(time.Duration(10+i*5) * time.Minute),
			LimitTime:     limitTime,
			StartWaitTime: startWait,
			MinLevel:      minLev,
			MaxLevel:      maxLev,
			MinPlayers:    minPly,
			MaxPlayers:    maxPly,
			Participants:  make(map[int]*KQParticipant),
		}
	}
	log.Printf("[KQManager] Loaded %d Kingdom Quests from KingdomQuest.shn", len(km.instances))
}

func (km *KQManager) loadDefaults() {
	now := time.Now()
	defaults := []struct {
		id     uint16
		title  string
		minLev byte
		maxLev byte
		maxPly uint16
	}{
		{1, "Mara Pirates' Ambition", 20, 30, 15},
		{2, "King Slime Counterattack", 30, 40, 15},
		{3, "Gold Garden", 40, 50, 20},
		{4, "Lost Temple", 50, 60, 20},
	}

	for _, d := range defaults {
		handle := km.nextHandle
		km.nextHandle++

		km.instances[handle] = &KQInstance{
			Handle:        handle,
			ID:            d.id,
			Title:         d.title,
			Status:        1, // Reclutando
			StartTime:     now.Add(15 * time.Minute),
			LimitTime:     3600,
			StartWaitTime: 300,
			MinLevel:      d.minLev,
			MaxLevel:      d.maxLev,
			MinPlayers:    5,
			MaxPlayers:    d.maxPly,
			Participants:  make(map[int]*KQParticipant),
		}
	}
	log.Printf("[KQManager] Initialized %d default Kingdom Quests", len(km.instances))
}

// RegisterOnline registra al jugador para recibir alertas de Kingdom Quest.
func (km *KQManager) RegisterOnline(client *network.Client, charID int) {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.onlinePlayers[charID] = client
}

// UnregisterOnline desregistra al jugador y lo remueve de la KQ si estaba anotado.
func (km *KQManager) UnregisterOnline(charID int) {
	km.mu.Lock()
	defer km.mu.Unlock()

	delete(km.onlinePlayers, charID)
	handle, ok := km.charToKQ[charID]
	if ok {
		delete(km.charToKQ, charID)
		if inst, found := km.instances[handle]; found {
			delete(inst.Participants, charID)
		}
	}
}

// BroadcastToOnline envía un paquete a todos los clientes conectados.
func (km *KQManager) BroadcastToOnline(packet *protocol.PacketWriter) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	for _, c := range km.onlinePlayers {
		if c != nil {
			_ = c.SendPacket(packet)
		}
	}
}

// HandleList entrega la lista completa de KQs activas al cliente (PROTO_NC_KQ_LIST_ACK 0x5802).
func (km *KQManager) HandleList(client *network.Client) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	now := time.Now()
	// 1. Reloj de sincronización
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_TIME_ACK(now))

	var list []protocol.ProtoKQInfoClient
	var startHandle, endHandle uint32

	for _, inst := range km.instances {
		if startHandle == 0 || inst.Handle < startHandle {
			startHandle = inst.Handle
		}
		if inst.Handle > endHandle {
			endHandle = inst.Handle
		}

		list = append(list, protocol.ProtoKQInfoClient{
			Handle:             inst.Handle,
			Status:             inst.Status,
			NumOfJoiner:        uint16(len(inst.Participants)),
			ID:                 inst.ID,
			Title:              inst.Title,
			LimitTime:          inst.LimitTime,
			StartTime:          inst.StartTime,
			StartWaitTime:      inst.StartWaitTime,
			MinLevel:           inst.MinLevel,
			MaxLevel:           inst.MaxLevel,
			MinPlayers:         inst.MinPlayers,
			MaxPlayers:         inst.MaxPlayers,
			PlayerRepeatMode:   0,
			PlayerRepeatCount:  0,
			PlayerRevivalMode:  0,
			PlayerRevivalCount: 3,
			DemandQuest:        0,
			DemandItem:         0,
			DemandClass:        0,
			DemandGender:       0,
		})
	}

	resp := protocol.New_PROTO_NC_KQ_LIST_ACK(now, startHandle, endHandle, list)
	_ = client.SendPacket(resp)
}

// HandleStatus responde el estado detallado y lista de inscritos en una KQ (OpCode 0x5804).
func (km *KQManager) HandleStatus(client *network.Client, handle uint32) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	inst, ok := km.instances[handle]
	if !ok {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_STATUS_ACK(handle, 0, nil))
		return
	}

	var names []string
	for _, p := range inst.Participants {
		names = append(names, p.CharName)
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_STATUS_ACK(handle, inst.Status, names))
}

// HandleJoin procesa la inscripción de un jugador en la Kingdom Quest (OpCode 0x5806).
func (km *KQManager) HandleJoin(client *network.Client, handle uint32, charLevel, charClass byte) {
	charID := client.CharacterID
	if charID <= 0 {
		return
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	// Chequear si ya está registrado en alguna KQ
	if _, registered := km.charToKQ[charID]; registered {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, 1)) // Error: ya inscrito
		return
	}

	inst, ok := km.instances[handle]
	if !ok {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, 2)) // Error: KQ no encontrada
		return
	}

	// Validar rango de nivel
	if charLevel < inst.MinLevel || charLevel > inst.MaxLevel {
		log.Printf("[KQ] Player '%s' (Lev: %d) out of level range (%d-%d) for '%s'",
			client.Username, charLevel, inst.MinLevel, inst.MaxLevel, inst.Title)
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, 3)) // Error: Nivel inválido
		return
	}

	// Validar cupos disponibles
	if uint16(len(inst.Participants)) >= inst.MaxPlayers {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, 4)) // Error: Cupo lleno
		return
	}

	// Inscribir al jugador
	inst.Participants[charID] = &KQParticipant{
		CharID:   charID,
		CharName: client.Username,
		Level:    charLevel,
		Class:    charClass,
		TeamType: 0,
		Client:   client,
	}
	km.charToKQ[charID] = handle

	log.Printf("[KQ] Player '%s' successfully joined KQ '%s' (Handle: %d, Total: %d/%d)",
		client.Username, inst.Title, handle, len(inst.Participants), inst.MaxPlayers)

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, 0))

	// Difundir actualización de cupos
	updatePkt := protocol.New_PROTO_NC_KQ_LIST_UPDATE_ACK([]protocol.KQUpdateItem{
		{
			Handle:      handle,
			Status:      inst.Status,
			NumOfJoiner: uint16(len(inst.Participants)),
		},
	})
	for _, c := range km.onlinePlayers {
		if c != nil {
			_ = c.SendPacket(updatePkt)
		}
	}
}

// HandleJoinCancel cancela la inscripción del jugador (OpCode 0x5808).
func (km *KQManager) HandleJoinCancel(client *network.Client, handle uint32) {
	charID := client.CharacterID
	km.mu.Lock()
	defer km.mu.Unlock()

	inst, ok := km.instances[handle]
	if !ok {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle, 1))
		return
	}

	delete(inst.Participants, charID)
	delete(km.charToKQ, charID)

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle, 0))

	// Difundir actualización
	updatePkt := protocol.New_PROTO_NC_KQ_LIST_UPDATE_ACK([]protocol.KQUpdateItem{
		{
			Handle:      handle,
			Status:      inst.Status,
			NumOfJoiner: uint16(len(inst.Participants)),
		},
	})
	for _, c := range km.onlinePlayers {
		if c != nil {
			_ = c.SendPacket(updatePkt)
		}
	}
}

// HandleJoinList retorna los detalles de todos los participantes anotados (OpCode 0x5832).
func (km *KQManager) HandleJoinList(client *network.Client, handle uint32) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	inst, ok := km.instances[handle]
	if !ok {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(1, nil))
		return
	}

	var chars []protocol.KQJoinCharInfo
	for _, p := range inst.Participants {
		chars = append(chars, protocol.KQJoinCharInfo{
			Level:    p.Level,
			Class:    p.Class,
			CharName: p.CharName,
			TeamType: p.TeamType,
		})
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(0, chars))
}

// HandleTeamSelect asigna el equipo del jugador en KQs con bandos (OpCode 0x5838).
func (km *KQManager) HandleTeamSelect(client *network.Client, teamType byte) {
	charID := client.CharacterID
	km.mu.Lock()
	defer km.mu.Unlock()

	handle, ok := km.charToKQ[charID]
	if !ok {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(1, teamType))
		return
	}

	inst, found := km.instances[handle]
	if !found {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(1, teamType))
		return
	}

	participant, exists := inst.Participants[charID]
	if !exists {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(1, teamType))
		return
	}

	participant.TeamType = teamType
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(0, teamType))

	// Notificar a todos los miembros de la sala
	cmdPkt := protocol.New_PROTO_NC_KQ_TEAM_SELECT_CMD(participant.CharName, teamType)
	for _, p := range inst.Participants {
		if p.Client != nil {
			_ = p.Client.SendPacket(cmdPkt)
		}
	}
}
