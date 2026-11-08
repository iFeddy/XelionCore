package world

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
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

// KQVoteSession administra una sesión activa de votación para expulsión.
type KQVoteSession struct {
	TargetCharID   int
	TargetCharName string
	StartedAt      time.Time
	VotesAgree     int
	VotesDisagree  int
	VotedPlayers   map[int]bool
}

// KQInstance representa el estado en memoria y persistente de una Kingdom Quest oficial.
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
	IsTeamPVP     bool
	CanVote       bool
	VoteAgreeRate byte
	Participants  map[int]*KQParticipant // CharID -> Participant
	ActiveVote    *KQVoteSession
}

// KQManager administra el ciclo de vida, inscripción, votación y difusión de Kingdom Quests (Dept 22).
type KQManager struct {
	mu            sync.RWMutex
	charRepo      database.CharacterRepository
	instances     map[uint32]*KQInstance
	charToKQ      map[int]uint32 // CharID -> Handle
	onlinePlayers map[int]*network.Client
	nextHandle    uint32
}

// QuestCount retorna el total de Kingdom Quests cargadas desde archivos SHN.
func (km *KQManager) QuestCount() int {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return len(km.instances)
}

// NewKQManager inicializa el gestor de Kingdom Quests a partir de archivos SHN oficiales y persiste en tKingdomQuest.
func NewKQManager(dataStore *WorldDataStore, charRepo database.CharacterRepository) *KQManager {
	km := &KQManager{
		charRepo:      charRepo,
		instances:     make(map[uint32]*KQInstance),
		charToKQ:      make(map[int]uint32),
		onlinePlayers: make(map[int]*network.Client),
		nextHandle:    1,
	}

	if dataStore != nil && dataStore.LoadedFiles != nil {
		km.loadFromDataStore(dataStore)
	}

	return km
}

func (km *KQManager) loadFromDataStore(store *WorldDataStore) {
	kqFile := store.LoadedFiles["KingdomQuest.shn"]
	if kqFile == nil || kqFile.RecordCount == 0 {
		log.Printf("[KQManager] KingdomQuest.shn not present or empty")
		return
	}

	// Cargar metadatos complementarios oficiales de SHN
	teamFile := store.LoadedFiles["KQTeam.shn"]
	isVoteFile := store.LoadedFiles["KQIsVote.shn"]
	voteRateFile := store.LoadedFiles["KQVoteMajorityRate.shn"]

	teamPvPMap := make(map[uint16]bool)
	if teamFile != nil {
		for i := 0; i < teamFile.RecordCount; i++ {
			row := teamFile.Rows[i]
			id := row.GetUint16("ID")
			isPVP := row.GetByte("IsTeamPVP")
			if isPVP != 0 {
				teamPvPMap[id] = true
			}
		}
	}

	canVoteMap := make(map[uint16]bool)
	if isVoteFile != nil {
		for i := 0; i < isVoteFile.RecordCount; i++ {
			row := isVoteFile.Rows[i]
			id := row.GetUint16("ID")
			isVote := row.GetByte("IsVote")
			canVoteMap[id] = (isVote != 0)
		}
	}

	voteMajorityRate := byte(70)
	if voteRateFile != nil && voteRateFile.RecordCount > 0 {
		rate := voteRateFile.Rows[0].GetByte("VoteAgreeRate")
		if rate > 0 {
			voteMajorityRate = rate
		}
	}

	now := time.Now()
	for i := 0; i < kqFile.RecordCount; i++ {
		row := kqFile.Rows[i]
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
			minPly = 1
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

		inst := &KQInstance{
			Handle:        handle,
			ID:            id,
			Title:         title,
			Status:        1, // Reclutando
			StartTime:     now.Add(time.Duration(10+i*5) * time.Minute),
			LimitTime:     limitTime,
			StartWaitTime: startWait,
			MinLevel:      minLev,
			MaxLevel:      maxLev,
			MinPlayers:    minPly,
			MaxPlayers:    maxPly,
			IsTeamPVP:     teamPvPMap[id],
			CanVote:       canVoteMap[id],
			VoteAgreeRate: voteMajorityRate,
			Participants:  make(map[int]*KQParticipant),
		}
		km.instances[handle] = inst

		// Sincronizar estado oficial en MySQL tKingdomQuest (CALL p_KingdomQuest_Set)
		if km.charRepo != nil {
			_ = km.charRepo.SetKingdomQuest(int(id), title)
		}
	}
	log.Printf("[KQManager] Loaded %d Kingdom Quests from KingdomQuest.shn (synced to tKingdomQuest)", len(km.instances))
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

// HandleSchedule responde el calendario de KQs programadas (0x5809 PROTO_NC_KQ_SCHEDULE_REQ -> 0x580A).
func (km *KQManager) HandleSchedule(client *network.Client) {
	km.HandleList(client)
}

// HandleRefresh refresca el listado de KQs para el cliente (0x581B PROTO_NC_KQ_LIST_REFRESH_REQ).
func (km *KQManager) HandleRefresh(client *network.Client) {
	km.HandleList(client)
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

// HandleVoteStartCheck verifica si una KQ admite votación (0x5834 PROTO_NC_KQ_VOTE_START_CHECK_REQ -> 0x5835).
func (km *KQManager) HandleVoteStartCheck(client *network.Client) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	handle, ok := km.charToKQ[client.CharacterID]
	canVote := uint16(0)
	if ok {
		if inst, found := km.instances[handle]; found && inst.CanVote {
			canVote = 1
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_START_CHECK_ACK)
	pw.WriteUint16(canVote)
	_ = client.SendPacket(pw)
}

// HandleVoteStart inicia una votación para expulsar a un jugador (0x5827 PROTO_NC_KQ_VOTE_START_REQ -> 0x5828).
func (km *KQManager) HandleVoteStart(client *network.Client, targetName string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	handle, ok := km.charToKQ[client.CharacterID]
	if !ok {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_START_ACK)
		pw.WriteUint16(1) // Not in KQ
		_ = client.SendPacket(pw)
		return
	}

	inst, found := km.instances[handle]
	if !found || !inst.CanVote || inst.ActiveVote != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_START_ACK)
		pw.WriteUint16(2) // Cannot vote or vote in progress
		_ = client.SendPacket(pw)
		return
	}

	var targetID int
	for _, p := range inst.Participants {
		if strings.EqualFold(p.CharName, strings.TrimSpace(targetName)) {
			targetID = p.CharID
			break
		}
	}

	if targetID == 0 || targetID == client.CharacterID {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_START_ACK)
		pw.WriteUint16(3) // Target invalid
		_ = client.SendPacket(pw)
		return
	}

	inst.ActiveVote = &KQVoteSession{
		TargetCharID:   targetID,
		TargetCharName: targetName,
		StartedAt:      time.Now(),
		VotesAgree:     1, // Requester automatically votes agree
		VotesDisagree:  0,
		VotedPlayers:   map[int]bool{client.CharacterID: true},
	}

	ackPw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_START_ACK)
	ackPw.WriteUint16(0)
	_ = client.SendPacket(ackPw)

	// Broadcast voting cmd (0x5829)
	votingCmd := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_VOTING_CMD)
	votingCmd.WriteString(targetName, 16)
	votingCmd.WriteString(client.Username, 16)
	for _, p := range inst.Participants {
		if p.Client != nil {
			_ = p.Client.SendPacket(votingCmd)
		}
	}
}

// HandleVoteCast registra el voto de un participante (0x582A PROTO_NC_KQ_VOTE_VOTING_REQ -> 0x582B).
func (km *KQManager) HandleVoteCast(client *network.Client, agree bool) {
	km.mu.Lock()
	defer km.mu.Unlock()

	handle, ok := km.charToKQ[client.CharacterID]
	if !ok {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)
		pw.WriteUint16(1)
		_ = client.SendPacket(pw)
		return
	}

	inst, found := km.instances[handle]
	if !found || inst.ActiveVote == nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)
		pw.WriteUint16(2)
		_ = client.SendPacket(pw)
		return
	}

	vote := inst.ActiveVote
	if vote.VotedPlayers[client.CharacterID] {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)
		pw.WriteUint16(3) // Already voted
		_ = client.SendPacket(pw)
		return
	}

	vote.VotedPlayers[client.CharacterID] = true
	if agree {
		vote.VotesAgree++
	} else {
		vote.VotesDisagree++
	}

	ackPw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)
	ackPw.WriteUint16(0)
	_ = client.SendPacket(ackPw)

	// Verificar si se alcanzó la mayoría
	totalParticipants := len(inst.Participants)
	threshold := int(float64(totalParticipants) * float64(inst.VoteAgreeRate) / 100.0)
	if threshold < 1 {
		threshold = 1
	}

	if vote.VotesAgree >= threshold {
		// Voto exitoso: expulsar objetivo
		targetID := vote.TargetCharID
		targetName := vote.TargetCharName
		delete(inst.Participants, targetID)
		delete(km.charToKQ, targetID)
		inst.ActiveVote = nil

		sucCmd := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_RESULT_SUC_CMD)
		sucCmd.WriteString(targetName, 16)
		for _, p := range inst.Participants {
			if p.Client != nil {
				_ = p.Client.SendPacket(sucCmd)
			}
		}
	} else if vote.VotesDisagree > totalParticipants-threshold {
		// Voto fallido
		inst.ActiveVote = nil
		failCmd := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_VOTE_RESULT_FAIL_CMD)
		for _, p := range inst.Participants {
			if p.Client != nil {
				_ = p.Client.SendPacket(failCmd)
			}
		}
	}
}
