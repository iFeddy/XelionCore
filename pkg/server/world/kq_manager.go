package world

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Constantes oficiales de CKQServer verificadas en WorldManager.exe + SingleData.shn.
const (
	kqMaxRooms                = 300        // exe: 0x00455895 (CKQServer::AddNewScheduleList cap 0x12c)
	kqMaxJoiners              = 100        // exe: 0x00454334 (CKQServer::PlayerJoin cap 0x64)
	kqMaxMapsPerRoom          = 4          // exe: 0x00453FF3 (CKQServer::AllocMapLink cmp ecx, 4)
	kqScheduleSlotsPerQuest   = 2          // exe: 0x004521F0 (CKQServer::GetNextScheduleTime cmp ebx, 2)
	kqCountdownSeconds        = 10         // exe: 0x0045608F (CKQServer::DoSetStart mov [ebx+0x46], 0x0a)
	kqMakeRoomTimeoutMin      = 1          // exe: 0x004522C7 (CKQServer::DoSetNoMapBF inc tm_min -> +60s)
	kqListAddFlushBytes       = 0x1D26     // exe: 0x00455D5B (CKQServer::Ack_NC_KQ_LIST_REFRESH batch flush >= 0x1D26)
	kqJoiningAlarmListMaxByte = 0x1F97     // exe: 0x00452C38 (CKQServer::Send_NC_KQ_JOINING_ALARM_LIST cap 0x1F97)
	kqTeamInvalid             = 2          // exe: 0x004543AF (KQTD_USERSELECT / fallback KQTeam=2)
	kqPercentBase             = 100        // exe: 0x004550D4 (VoteProcessing imul eax, 0x64)
	kqUnsetHandle             = ^uint32(0) // exe: 0x0045473C (HANDLE -1 = 0xFFFFFFFF)

	// Valores por defecto de SingleData.shn cuando no hay store externo (CSingleDataMap):
	kqVoteLimitTimeDefaultSec      = 60  // exe: 0x004562B9 (SingleData.shn KQVote_VoteLimitTime = 60)
	kqVoteSuggestCoolTimeDefault   = 300 // exe: 0x0045666F (SingleData.shn KQVote_SuggestCoolTime = 300)
	kqVoteLoginCoolTimeDefault     = 300 // exe: 0x00454F31 (SingleData.shn KQVote_LoginCoolTime = 300)
	kqJoinListResetCoolTimeDefault = 5   // exe: 0x004569BE (SingleData.shn KQPlayerList_ResetListCoolTime = 5)
)

// Códigos de error oficiales KQ (WorldManager.exe).
const (
	kqErrJoinOK              uint16 = 0x0991 // exe: 0x00413D3E (0x991 + 0)
	kqErrJoinFailedGeneric   uint16 = 0x0998 // exe: 0x00454385 (ret > 5 -> 7 -> 0x998)
	kqErrJoinPrisonBlocked   uint16 = 0x0999 // exe: 0x00413CFE (word [charData+0x62] != 0)
	kqErrJoinAlreadyJoined   uint16 = 0x099A // exe: 0x00413D19 (IsPlayerAlreadyJoin != 0)
	kqErrJoinCancelOK        uint16 = 0x09A1 // exe: 0x00413DBE (PlayerDisjoin != 0)
	kqErrJoinCancelFail      uint16 = 0x09A2 // exe: 0x00413DC9 (PlayerDisjoin == 0)
	kqErrZ2WMakeSuccess      uint16 = 0x0981 // exe: 0x0042AA0E (fc_NC_KQ_Z2W_MAKE_ACK)
	kqErrTeamSelectOK        uint16 = 0x31F0 // exe: 0x004536CC
	kqErrTeamSelectNotJoined uint16 = 0x31F1 // exe: 0x0045371E
	kqErrTeamSelectBadStatus uint16 = 0x31F2 // exe: 0x00453753
	kqErrTeamSelectSameTeam  uint16 = 0x31F3 // exe: 0x0045378C
	kqErrTeamSelectFull      uint16 = 0x31F4 // exe: 0x004537AC
	kqErrTeamSelectGapOver   uint16 = 0x31F5 // exe: 0x004537E1
	kqErrTeamSelectNoTeamCfg uint16 = 0x31F6 // exe: 0x0045381A
	kqErrTeamSelectBadDivide uint16 = 0x31F7 // exe: 0x00453831
	kqErrVoteStartOK         uint16 = 0x3100 // exe: 0x00456397
	kqErrVoteStartNotJoined  uint16 = 0x3101 // exe: 0x004563AE
	kqErrVoteStartBadStatus  uint16 = 0x3102 // exe: 0x00456420
	kqErrVoteStartInVoting   uint16 = 0x3103 // exe: 0x00456442
	kqErrVoteStartBadTarget  uint16 = 0x3105 // exe: 0x004564E3
	kqErrVoteStartSelf       uint16 = 0x3106 // exe: 0x004565C7
	kqErrVoteStartNoReason   uint16 = 0x3107 // exe: 0x004565EB
	kqErrVoteStartCooldown   uint16 = 0x3108 // exe: 0x0045661B
	kqErrVoteStartDisabled   uint16 = 0x3109 // exe: 0x0045646A
	kqErrVoteCastOK          uint16 = 0x3110 // exe: 0x00452F66
	kqErrVoteCastNotJoined   uint16 = 0x3111 // exe: 0x00452F7D
	kqErrVoteCastNotVoting   uint16 = 0x3112 // exe: 0x00453001
	kqErrVoteCastCannotVote  uint16 = 0x3113 // exe: 0x00453044
	kqErrJoinListOK          uint16 = 0x3118 // exe: 0x0045692C
	kqErrJoinListNotFound    uint16 = 0x3119 // exe: 0x0045696F
	kqErrJoinListCooldown    uint16 = 0x311A // exe: 0x004569E9
	kqErrVoteCheckOK         uint16 = 0x3120 // exe: 0x004533F3
	kqErrVoteCheckNotJoined  uint16 = 0x3121 // exe: 0x0045340C
	kqErrVoteCheckDisabled   uint16 = 0x3122 // exe: 0x0045346C
)

// kqJoiner representa un slot KQ_JOINER_BF (stride 0x2C bytes) en WorldManager.exe.
type kqJoiner struct {
	CharRegNo uint32    // +0x00
	Slot      uint16    // +0x04
	Level     byte      // +0x06
	Class     byte      // +0x07
	CharName  string    // +0x08 (Name5[20])
	TeamType  byte      // +0x1C (0=Red, 1=Blue, 2=None)
	CanVote   bool      // +0x1E (int32)
	BanForce  bool      // +0x20 (int32)
	IsBan     bool      // +0x24 (int32)
	VoteCount byte      // +0x28
	BanAt     time.Time // para CheckCharBannedInLogin
}

// kqVoteState modela el estado de votación activa de una sala (offset +0x15E694).
type kqVoteState struct {
	TeamType  byte      // +0x00
	Starter   int       // +0x04 (-1 si inactivo)
	Target    int       // +0x08 (-1 si inactivo)
	EndTime   time.Time // +0x0C
	YesNum    byte      // +0x10
	NoNum     byte      // +0x11
	CancelNum byte      // +0x12
}

// kqRoom representa una entrada CKQ + datos asociados en CKQServer.
type kqRoom struct {
	Info       kq.PROTO_KQ_INFO
	Joiners    []*kqJoiner
	TeamCounts [2]byte
	Vote       kqVoteState
}

// kqSessionState mantiene los campos KQ de CWMClientSession en WorldManager.exe.
type kqSessionState struct {
	Client            *network.Client
	HKQ               uint32 // +0x258C (-1 = kqUnsetHandle)
	TimeAckSent       bool   // +0x2590
	Snapshot          []kq.PROTO_KQ_INFO
	AutoRefreshActive bool      // +0x1DF68
	LastRefreshAt     time.Time // +0x1DF6C
	VoteSuggestNextAt time.Time // +0x1DF70
	BannedInLogin     bool      // +0x1DF74
	JoinListNextAt    time.Time
	Level             byte
	Class             byte
	Gender            byte
	PrisonMin         uint16
}

// KQManager porta CKQServer (WorldManager.exe @ 0x00451B40..0x00456C00) 1:1 en Go.
type KQManager struct {
	mu           sync.RWMutex
	charRepo     database.CharacterRepository
	store        *kq.Store
	s2sSink      S2SSink
	worldID      byte
	onLeaveParty func(charID int) // hook para KQTeam_LeaveParty

	// Datos auxiliares cargados de SHNs oficiales (ReadGlobalGameData @ 0x00404100)
	singleData       map[string]uint16
	msgWorldManager  []string
	useClassTypeBits map[uint32]int64

	// Estado de CKQServer
	rooms         []*kqRoom
	nextHandle    uint32    // exe: 0x004556EB (handle allocation empieza en 0)
	mapLinkOwners [][]int64 // [mapRowIdx][slotIdx] -> room Handle o -1

	lastDoTimeSec   int64 // exe: 0x004561B2 (m_NextDoTime)
	lastScheduleMin int64 // exe: 0x00455E8A ([ebx+0x1bb48])

	onlineByCharID map[int]*network.Client
	sessions       map[*network.Client]*kqSessionState
	nowFunc        func() time.Time
}

// NewKQManager inicializa CKQServer a partir de los 5 SHNs KQ y tablas globales de WorldDataStore.
// WorldManager.exe @ 0x00404100 (ReadGlobalGameData) -> 0x00451B40 (CKQServer::Read).
// Nota: WorldManager.exe NO invoca p_KingdomQuest_Set al arrancar; carga SHNs en memoria.
func NewKQManager(dataStore *WorldDataStore, charRepo database.CharacterRepository) *KQManager {
	km := &KQManager{
		charRepo:         charRepo,
		singleData:       make(map[string]uint16),
		useClassTypeBits: make(map[uint32]int64),
		onlineByCharID:   make(map[int]*network.Client),
		sessions:         make(map[*network.Client]*kqSessionState),
		nowFunc:          time.Now,
	}

	if dataStore != nil {
		km.loadKQData(dataStore)
		km.loadAuxData(dataStore)
	}
	if km.store != nil {
		km.initMapTable()
		// Primera pasada equivalente al primer segundo de MainWork (exe: 0x004070F0)
		km.DoSchedule(time.Now())
	}
	return km
}

// SetClock configura el proveedor de reloj para pruebas determinísticas con reloj simulado.
func (km *KQManager) SetClock(fn func() time.Time) {
	km.mu.Lock()
	defer km.mu.Unlock()
	if fn == nil {
		km.nowFunc = time.Now
		return
	}
	km.nowFunc = fn
}

func (km *KQManager) nowTime() time.Time {
	if km.nowFunc != nil {
		return km.nowFunc()
	}
	return time.Now()
}

// SetS2SSink conecta el bus S2S para W2Z_MAKE_REQ / W2Z_START_CMD / W2Z_DESTROY_CMD / BAN.
func (km *KQManager) SetS2SSink(sink S2SSink, worldID byte) {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.s2sSink = sink
	km.worldID = worldID
}

// SetLeavePartyHook conecta KQTeam_LeaveParty (exe: 0x00453490) con PartyManager.
func (km *KQManager) SetLeavePartyHook(fn func(charID int)) {
	km.mu.Lock()
	defer km.mu.Unlock()
	km.onLeaveParty = fn
}

func (km *KQManager) loadKQData(ds *WorldDataStore) {
	if ds.DataDir != "" {
		st, err := kq.LoadDir(ds.DataDir)
		if err == nil {
			km.store = st
			return
		}
		log.Printf("[CKQServer] kq.LoadDir(%s) failed: %v", ds.DataDir, err)
	}
}

func (km *KQManager) loadAuxData(ds *WorldDataStore) {
	// CSingleDataMap::InitSingleDataArray (exe: 0x00459DE0) sobre SingleData.shn
	if ds.SingleData != nil {
		for _, r := range ds.SingleData.Rows {
			key := strings.TrimRight(string(r.GetBytesAt(0)), "\x00")
			val := r.GetUint16At(1)
			if key != "" {
				km.singleData[key] = val
			}
		}
	}
	// WorldManagerServer::GetMsg (exe: 0x00404570) sobre MsgWorldManager.shn
	if ds.MsgWorldManager != nil {
		km.msgWorldManager = make([]string, ds.MsgWorldManager.RecordCount)
		for i, r := range ds.MsgWorldManager.Rows {
			km.msgWorldManager[i] = strings.TrimRight(string(r.GetBytesAt(1)), "\x00")
		}
	}
	// CharClassDataBox::ccdb_UseClassTypeToBit (exe: 0x00465A90) sobre UseClassTypeInfo.shn
	if ds.UseClassTypeInfo != nil {
		for _, r := range ds.UseClassTypeInfo.Rows {
			uc := uint32(r.GetIntAt(0))
			var mask int64
			// Columnas 1..27 (Fig..Sav); shl final deja Fig en bit 1 (exe: 0x00465C6C..0x00465C70)
			for col := 1; col <= 27 && col < len(ds.UseClassTypeInfo.Columns); col++ {
				if r.GetByteAt(col) != 0 {
					mask |= 1 << col
				}
			}
			km.useClassTypeBits[uc] = mask
		}
	}
}

func (km *KQManager) getSingleDataValue(key string, fallback uint16) uint16 {
	if v, ok := km.singleData[key]; ok {
		return v
	}
	return fallback
}

func (km *KQManager) getMsg(idx int) string {
	// exe: 0x00404570 (WorldManagerServer::GetMsg)
	if idx < 0 || idx >= len(km.msgWorldManager) {
		return ""
	}
	return km.msgWorldManager[idx]
}

func (km *KQManager) useClassTypeToBit(uc uint32) int64 {
	// exe: 0x00465A90 (CharClassDataBox::ccdb_UseClassTypeToBit)
	return km.useClassTypeBits[uc]
}

// initMapTable inicializa el pool de instancias de mapas en -1 (exe: 0x00451E90 CKQServer::InitMapTable).
func (km *KQManager) initMapTable() {
	if km.store == nil {
		return
	}
	km.mapLinkOwners = make([][]int64, len(km.store.Maps))
	for i := range km.store.Maps {
		cnt := int(km.store.Maps[i].NumOfMap)
		if cnt < 0 {
			cnt = 0
		}
		if cnt > 10 { // exe: 0x00451F00 (10 slots máx por fila de KingdomQuestMap.shn)
			cnt = 10
		}
		km.mapLinkOwners[i] = make([]int64, cnt)
		for j := 0; j < cnt; j++ {
			km.mapLinkOwners[i][j] = -1 // exe: 0x00451F30
		}
	}
}

// getEmptyMapLink busca un slot libre en la fila mapIdx (exe: 0x00451FD0).
func (km *KQManager) getEmptyMapLink(mapIdx int) int {
	if km.store == nil || mapIdx < 0 || mapIdx >= len(km.store.Maps) {
		return -1
	}
	m := &km.store.Maps[mapIdx]
	cnt := int(m.NumOfMap)
	if cnt > len(km.mapLinkOwners[mapIdx]) {
		cnt = len(km.mapLinkOwners[mapIdx])
	}
	for j := 0; j < cnt; j++ {
		if m.Clear[j] == 0 { // exe: 0x00452013
			return j
		}
		if km.mapLinkOwners[mapIdx][j] == -1 { // exe: 0x00452029
			return j
		}
	}
	return -1
}

// freeMapLink libera todos los slots asignados a handle (exe: 0x00452040).
func (km *KQManager) freeMapLink(handle uint32) {
	for i := range km.mapLinkOwners {
		for j := range km.mapLinkOwners[i] {
			if km.mapLinkOwners[i][j] == int64(handle) {
				km.mapLinkOwners[i][j] = -1 // exe: 0x004520BD
			}
		}
	}
}

// allocMapLink reserva los MapLink[0..3] de KINGDOM_QUEST para la sala (exe: 0x00453EB0).
func (km *KQManager) allocMapLink(room *kqRoom) bool {
	if km.store == nil {
		return false
	}
	var tpl *kq.KINGDOM_QUEST
	for i := range km.store.Quests {
		if km.store.Quests[i].ID == room.Info.ID { // exe: 0x00453EEB
			tpl = &km.store.Quests[i]
			break
		}
	}
	if tpl == nil {
		return false
	}
	for slot := 0; slot < kqMaxMapsPerRoom; slot++ {
		mapIdx := tpl.MapLink[slot]
		if mapIdx == 0xFFFF { // exe: 0x00453F23
			continue
		}
		if int(mapIdx) >= len(km.store.Maps) { // exe: 0x00453F46
			km.freeMapLink(room.Info.Handle)
			return false
		}
		linkIdx := km.getEmptyMapLink(int(mapIdx)) // exe: 0x00453F5B
		if linkIdx == -1 {
			km.freeMapLink(room.Info.Handle) // exe: 0x00454013
			return false
		}
		m := &km.store.Maps[mapIdx]
		if m.Clear[linkIdx] != 0 { // exe: 0x00453F6B
			km.mapLinkOwners[mapIdx][linkIdx] = int64(room.Info.Handle)
		}
		mi := kq.PROTO_KQ_MAP_INFO{
			MapIndex: byte(linkIdx),
			MapClear: m.Clear[linkIdx],
		}
		mi.SetMapBase(m.BaseMapString())
		copy(mi.MapName[:], []byte(m.MapString(linkIdx)))
		room.Info.MapLink[slot] = mi
	}
	return true
}

// QuestCount retorna la cantidad de Kingdom Quests actualmente programadas/activas en CKQ.
func (km *KQManager) QuestCount() int {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return len(km.rooms)
}

// TemplateCount retorna el total de definiciones cargadas de KingdomQuest.shn.
func (km *KQManager) TemplateCount() int {
	km.mu.RLock()
	defer km.mu.RUnlock()
	if km.store == nil {
		return 0
	}
	return len(km.store.Quests)
}

func (km *KQManager) getRoomIndexByHandle(handle uint32) int {
	// exe: 0x0045CD80 (CKQ::GetBFIndex)
	for i, r := range km.rooms {
		if r.Info.Handle == handle {
			return i
		}
	}
	return -1
}

func (km *KQManager) getRoomByHandle(handle uint32) (*kqRoom, int) {
	idx := km.getRoomIndexByHandle(handle)
	if idx < 0 {
		return nil, -1
	}
	return km.rooms[idx], idx
}

func (km *KQManager) getJoinerIndex(room *kqRoom, charName string) int {
	// exe: 0x00452DB0 (CKQServer::GetJoinerIndex)
	for i, j := range room.Joiners {
		if j.CharName == charName {
			return i
		}
	}
	return -1
}

func (km *KQManager) clearVoteInfo(room *kqRoom) {
	// exe: 0x00452CF0 (CKQServer::ClearVoteInfo)
	room.Vote = kqVoteState{
		TeamType: kqTeamInvalid, // exe: 0x00452D14
		Starter:  -1,            // exe: 0x00452D08
		Target:   -1,            // exe: 0x00452D0E
	}
	for _, j := range room.Joiners {
		j.CanVote = false // exe: 0x00452D31
	}
}

func (km *KQManager) inVoting(room *kqRoom) bool {
	// exe: 0x00452D40 (CKQServer::InVoting)
	return room.Vote.Starter != -1 && room.Vote.Target != -1
}

func (km *KQManager) getOrCreateSessionLocked(client *network.Client) *kqSessionState {
	st, ok := km.sessions[client]
	if !ok {
		st = &kqSessionState{
			Client: client,
			HKQ:    kqUnsetHandle,
		}
		km.sessions[client] = st
	}
	return st
}

// RegisterOnline registra una sesión activa y sincroniza su nivel/clase/género/prisión si existe DB.
func (km *KQManager) RegisterOnline(client *network.Client, charID int) {
	if client == nil || charID <= 0 {
		return
	}
	km.mu.Lock()
	km.onlineByCharID[charID] = client
	st := km.getOrCreateSessionLocked(client)
	km.mu.Unlock()

	if km.charRepo != nil {
		if ch, err := km.charRepo.GetCharacterByID(charID); err == nil && ch != nil {
			km.mu.Lock()
			st.Level = ch.Level
			st.Class = ch.Job
			st.Gender = ch.TCharacterShape.Gender
			km.mu.Unlock()
		}
	}
}

// SetSessionCharacterStats actualiza metadata de personaje usada en IsJoinable y fc_NC_KQ_JOIN_REQ.
func (km *KQManager) SetSessionCharacterStats(client *network.Client, level, class, gender byte, prisonMin uint16) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()
	st := km.getOrCreateSessionLocked(client)
	st.Level = level
	st.Class = class
	st.Gender = gender
	st.PrisonMin = prisonMin
}

// UnregisterOnline desregistra la sesión del jugador al desconectarse y ejecuta PlayerDisjoin si no está Running.
// exe: 0x00436DF0..0x00436E09 (CWMClientSession::Logout)
func (km *KQManager) UnregisterOnline(charID int) {
	km.mu.Lock()
	defer km.mu.Unlock()
	client, ok := km.onlineByCharID[charID]
	if ok && client != nil {
		if st, exists := km.sessions[client]; exists {
			if st.HKQ != kqUnsetHandle {
				if room, _ := km.getRoomByHandle(st.HKQ); room == nil || room.Info.Status != kq.KQStatusRunning {
					km.playerDisjoinLocked(client)
				}
			}
			delete(km.sessions, client)
		}
		delete(km.onlineByCharID, charID)
	}
}

// RemovePlayer remueve explícitamente a un personaje de cualquier sala no iniciada.
func (km *KQManager) RemovePlayer(charID int) {
	if charID <= 0 {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()
	if client, ok := km.onlineByCharID[charID]; ok && client != nil {
		km.playerDisjoinLocked(client)
		delete(km.sessions, client)
		delete(km.onlineByCharID, charID)
	}
}

// DoKingdomQuest ejecuta el tick principal de CKQServer::DoKingdomQuest (WorldManager.exe @ 0x00456190).
// Gateado a 1 vez por segundo (exe: 0x004561B2). Retorna número de transiciones realizadas.
func (km *KQManager) DoKingdomQuest(now time.Time) int {
	km.mu.Lock()
	defer km.mu.Unlock()

	nowSec := now.Unix()
	if nowSec < km.lastDoTimeSec { // exe: 0x004561BA
		return 0
	}
	km.lastDoTimeSec = nowSec + 1 // exe: 0x004561DE (tm_sec + 1)

	transitions := 0
	transitions += km.doSetNoMapBFLocked(now)    // exe: 0x004561FB
	transitions += km.doScheduleLocked(now)      // exe: 0x00456202
	transitions += km.doSetMakeRoomLocked(now)   // exe: 0x00456209
	transitions += km.doSetStartLocked(now)      // exe: 0x00456210
	transitions += km.delOldScheduleListLocked() // exe: 0x00456217
	km.autoRefreshListLocked(now)                // exe: 0x0045621E
	km.voteProcessingLocked(now)                 // exe: 0x00456225
	return transitions
}

// DoSchedule ejecuta el planificador minuto-a-minuto de CKQServer (exe: 0x00455E40).
func (km *KQManager) DoSchedule(now time.Time) int {
	km.mu.Lock()
	defer km.mu.Unlock()
	return km.doScheduleLocked(now)
}

func (km *KQManager) doScheduleLocked(now time.Time) int {
	if km.store == nil || len(km.store.Quests) == 0 {
		return 0
	}
	// Truncar now al inicio del minuto actual (exe: 0x00455E7C tm_sec = 0)
	loc := now.Location()
	nowMin := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)
	nowMinUnix := nowMin.Unix()
	if nowMinUnix == km.lastScheduleMin { // exe: 0x00455E8A
		return 0
	}
	km.lastScheduleMin = nowMinUnix // exe: 0x00455EA5

	added := 0
	for i := range km.store.Quests {
		q := &km.store.Quests[i]
		// exe: 0x00455EDE..0x00455F0A: usa Year/Month de now, Day/Hour/Minute de KINGDOM_QUEST
		base := time.Date(nowMin.Year(), nowMin.Month(), int(q.STDay), int(q.STHour), int(q.STMinute), 0, 0, loc)
		times := getNextScheduleTimes(base, q.NextStartDeleyMin, nowMin) // exe: 0x00455F3E
		for _, t := range times {
			if km.addNewScheduleLocked(q, t) { // exe: 0x00455F4E
				added++
			}
		}
	}
	return added
}

// getNextScheduleTimes porta CKQServer::GetNextScheduleTime (exe: 0x00452160).
func getNextScheduleTimes(base time.Time, delayMin uint16, nowMin time.Time) []time.Time {
	if delayMin == 0 {
		if !base.Before(nowMin) {
			return []time.Time{base}
		}
		return nil
	}
	step := time.Duration(delayMin) * time.Minute
	cur := base
	for cur.Before(nowMin) { // exe: 0x0045219D..0x004521BA
		cur = cur.Add(step)
	}
	out := make([]time.Time, 0, kqScheduleSlotsPerQuest)
	for i := 0; i < kqScheduleSlotsPerQuest; i++ { // exe: 0x004521F0
		out = append(out, cur)
		cur = cur.Add(step)
	}
	return out
}

// addNewScheduleLocked porta CKQServer::AddNewScheduleList (exe: 0x004555D0).
func (km *KQManager) addNewScheduleLocked(q *kq.KINGDOM_QUEST, schedTime time.Time) bool {
	if len(km.rooms) >= kqMaxRooms { // exe: 0x00455895
		return false
	}
	schedUnix := schedTime.Unix()
	for _, r := range km.rooms {
		if r.Info.ID == q.ID && int64(r.Info.ScheduleTime) == schedUnix { // exe: 0x00455676..0x00455685
			return false
		}
	}

	var handle uint32
	for {
		handle = km.nextHandle // exe: 0x004556EB
		km.nextHandle++
		if km.getRoomIndexByHandle(handle) < 0 {
			break
		}
	}

	var isTeamPVP uint8
	var regenXY [2]kq.SHINE_XY_TYPE
	if kt, ok := km.store.GetKQTeamData(q.ID); ok { // exe: 0x00455874..0x00455890
		isTeamPVP = kt.IsTeamPVP
		regenXY[0] = kq.SHINE_XY_TYPE{X: kt.RegenXRed, Y: kt.RegenYRed}
		regenXY[1] = kq.SHINE_XY_TYPE{X: kt.RegenXBlue, Y: kt.RegenYBlue}
	}
	tmVal := kq.CTimeTMFromTime(schedTime)
	room := &kqRoom{
		Info: kq.PROTO_KQ_INFO{
			PROTO_KQ_INFO_CLIENT: kq.PROTO_KQ_INFO_CLIENT{
				Handle:             handle,
				Status:             kq.KQStatusScheduled, // exe: 0x00455704
				NumOfJoiner:        0,
				ID:                 q.ID,
				Title:              q.Title,
				LimitTime:          q.LimitTime,
				StartTime:          int32(schedUnix),
				TmStartTime:        tmVal,
				StartWaitTime:      q.StartWaitTime,
				MinLevel:           q.MinLevel,
				MaxLevel:           q.MaxLevel,
				MinPlayers:         q.MinPlayers,
				MaxPlayers:         q.MaxPlayers,
				PlayerRepeatMode:   q.PlayerRepeatMode,
				PlayerRepeatCount:  q.PlayerRepeatCount,
				PlayerRevivalMode:  q.PlayerRevivalMode,
				PlayerRevivalCount: q.PlayerRevivalCount,
				DemandQuest:        q.DemandQuest,
				DemandItem:         q.DemandItem,
				DemandClass:        km.useClassTypeToBit(q.UseClass),        // exe: 0x00455804
				DemandGender:       q.DemandGender[1]*2 + q.DemandGender[0], // exe: 0x0045580F
			},
			NextStartMode:     q.NextStartMode,
			NextStartDelayMin: q.NextStartDeleyMin,
			RepeatMode:        q.RepeatMode,
			RepeatCount:       q.RepeatCount,
			RewardIndex:       uint16(q.RewardIndex),
			DemandMobKill:     uint16(q.DemandMobKill),
			ScheduleTime:      int32(schedUnix),
			TmScheduleTime:    tmVal,
			RunCounter:        0,
			ScriptLanguage:    q.ScriptLanguage,
			ScriptInitValue:   q.InitValue,
			IsTeamPVP:         isTeamPVP,
			TeamRegenXY:       regenXY,
		},
		Vote: kqVoteState{
			TeamType: kqTeamInvalid,
			Starter:  -1,
			Target:   -1,
		},
	}
	km.rooms = append(km.rooms, room)
	return true
}

// doSetNoMapBFLocked porta CKQServer::DoSetNoMapBF (exe: 0x00452270).
func (km *KQManager) doSetNoMapBFLocked(now time.Time) int {
	nowUnix := now.Unix()
	changed := 0
	for _, r := range km.rooms {
		if r.Info.Status != kq.KQStatusMakingRoom { // exe: 0x004522A0
			continue
		}
		schedTime := time.Unix(int64(r.Info.ScheduleTime), 0)
		deadline := schedTime.Add(kqMakeRoomTimeoutMin * time.Minute).Unix() // exe: 0x004522C7
		if deadline <= nowUnix {                                              // exe: 0x004522EE
			r.Info.Status = kq.KQStatusMakeTimeout // exe: 0x00452300
			km.freeMapLink(r.Info.Handle)          // exe: 0x00452340
			changed++
		}
	}
	return changed
}

// setNoMapBFLocked porta CKQServer::SetNoMapBF (exe: 0x00452440).
func (km *KQManager) setNoMapBFLocked(handle uint32) {
	r, _ := km.getRoomByHandle(handle)
	if r == nil {
		return
	}
	r.Info.Status = kq.KQStatusNoMap // exe: 0x00452465
	km.freeMapLink(handle)
}

// doSetMakeRoomLocked porta CKQServer::DoSetMakeRoom (exe: 0x00455980).
func (km *KQManager) doSetMakeRoomLocked(now time.Time) int {
	nowUnix := now.Unix()
	changed := 0
	for _, r := range km.rooms {
		if r.Info.Status != kq.KQStatusScheduled { // exe: 0x00455998
			continue
		}
		if int64(r.Info.ScheduleTime) > nowUnix { // exe: 0x004559A3
			continue
		}
		if km.allocMapLink(r) { // exe: 0x004559AE
			r.Info.Status = kq.KQStatusMakingRoom // exe: 0x004559BA
			km.sendW2ZMakeReqLocked(r)            // exe: 0x004559BD
		} else {
			r.Info.Status = kq.KQStatusNoMap // exe: 0x004559C4
		}
		changed++
	}
	return changed
}

// setJoiningLocked porta CKQServer::SetJoining (exe: 0x00454170).
func (km *KQManager) setJoiningLocked(handle uint32) {
	r, _ := km.getRoomByHandle(handle)
	if r == nil {
		return
	}
	r.Info.Status = kq.KQStatusJoining // exe: 0x00454192
	msgFmt := km.getMsg(0)             // exe: 0x004541B0 ("Recruitment for Kingdom Quest - '%s' has begun.")
	msg := ""
	if msgFmt != "" {
		msg = fmt.Sprintf(msgFmt, r.Info.TitleString())
	}
	alarm := protocol.New_PROTO_NC_KQ_JOINING_ALARM_CMD(r.Info.Handle, r.Info.ID, r.Info.MinLevel, r.Info.MaxLevel, msg)
	for _, c := range km.onlineByCharID {
		if c != nil && c.CharacterID > 0 {
			_ = c.SendPacket(alarm)
		}
	}
}

// doSetStartLocked porta CKQServer::DoSetStart (exe: 0x00455F90).
func (km *KQManager) doSetStartLocked(now time.Time) int {
	nowUnix := now.Unix()
	changed := 0
	for i, r := range km.rooms {
		switch r.Info.Status {
		case kq.KQStatusCountdown: // exe: 0x00455FD3
			if r.Info.RunCounter == 0 { // exe: 0x00455FDA
				r.Info.Status = kq.KQStatusRunning // exe: 0x00455FE1
				km.kqTeamDivideRandomLocked(r)     // exe: 0x00455FE5
				km.kqTeamLeavePartyLocked(r)       // exe: 0x00455FED
				km.sendW2ZStartCmdLocked(r)        // exe: 0x00455FF5
				changed++
			} else {
				if r.Info.RunCounter == kqCountdownSeconds { // exe: 0x00455FFF
					msgFmt := km.getMsg(1)
					if msgFmt != "" {
						msg := fmt.Sprintf(msgFmt, r.Info.TitleString(), int(r.Info.RunCounter))
						km.sendNotifyToJoinersLocked(r, msg) // exe: 0x00456031
					}
				}
				km.sendJoiningAlarmEndLocked(r) // exe: 0x00456039
				r.Info.RunCounter--             // exe: 0x0045603E
				changed++
			}

		case kq.KQStatusJoining: // exe: 0x00456046
			if len(r.Joiners) >= int(r.Info.MaxPlayers) { // exe: 0x00456059
				r.Info.Status = kq.KQStatusCountdown // exe: 0x0045608B
				r.Info.RunCounter = kqCountdownSeconds
				changed++
				continue
			}
			// exe: 0x004560B4 (StartWaitTime se suma a tm_sec sobre ScheduleTime)
			joinDeadline := int64(r.Info.ScheduleTime) + int64(r.Info.StartWaitTime)
			if nowUnix < joinDeadline { // exe: 0x004560CF..0x004560D5
				continue
			}
			if len(r.Joiners) < int(r.Info.MinPlayers) { // exe: 0x004560E3
				km.setDoneSkipLocked(i, 2) // exe: 0x00456124
				changed++
			} else if km.kqTeamCanKQStartLocked(i) { // exe: 0x00456138
				r.Info.Status = kq.KQStatusCountdown // exe: 0x00456142
				r.Info.RunCounter = kqCountdownSeconds
				changed++
			}
		}
	}
	return changed
}

// setDoneLocked porta CKQServer::SetDone (exe: 0x00454860).
func (km *KQManager) setDoneLocked(handle uint32) bool {
	r, idx := km.getRoomByHandle(handle)
	if r == nil {
		return false
	}
	r.Info.Status = kq.KQStatusDone // exe: 0x00454892
	km.sendW2ZDestroyCmdLocked(handle)
	km.freeMapLink(handle)
	km.freeJoinerLocked(idx)
	return true
}

// setDoneSkipLocked porta CKQServer::SetDoneSkip (exe: 0x004548C0).
func (km *KQManager) setDoneSkipLocked(idx int, reason int) {
	if idx < 0 || idx >= len(km.rooms) {
		return
	}
	r := km.rooms[idx]
	r.Info.Status = kq.KQStatusDoneSkip // exe: 0x004548F8
	km.sendW2ZDestroyCmdLocked(r.Info.Handle)
	km.freeMapLink(r.Info.Handle)

	switch reason {
	case 2: // exe: 0x004549AA ("...canceled due to lack of participants(%d/%d).")
		if msgFmt := km.getMsg(2); msgFmt != "" {
			msg := fmt.Sprintf(msgFmt, r.Info.TitleString(), r.Info.NumOfJoiner, r.Info.MinPlayers)
			km.sendNotifyToJoinersLocked(r, msg)
		}
	case 3: // exe: 0x00454956
		if m3 := km.getMsg(3); m3 != "" {
			km.sendNotifyToJoinersLocked(r, m3)
		}
		if m4 := km.getMsg(4); m4 != "" {
			km.sendNotifyToJoinersLocked(r, m4)
		}
	}
	km.freeJoinerLocked(idx)
}

// freeJoinerLocked porta CKQServer::FreeJoiner (exe: 0x004526E0).
func (km *KQManager) freeJoinerLocked(idx int) {
	if idx < 0 || idx >= len(km.rooms) {
		return
	}
	r := km.rooms[idx]
	for _, j := range r.Joiners {
		for c, st := range km.sessions {
			if c != nil && uint32(c.CharacterID) == j.CharRegNo && st.HKQ == r.Info.Handle {
				st.HKQ = kqUnsetHandle
			}
		}
	}
	r.Joiners = nil
	r.Info.NumOfJoiner = 0
	r.TeamCounts = [2]byte{0, 0}
	km.clearVoteInfo(r)
}

// delOldScheduleListLocked porta CKQServer::DelOldShceduleList (exe: 0x00454080).
func (km *KQManager) delOldScheduleListLocked() int {
	deletedAny := false
	for _, r := range km.rooms {
		if r.Info.Status < kq.KQStatusDone || r.Info.Status > kq.KQStatusNotFound { // exe: 0x004540B2..0x004540B6
			continue
		}
		for _, other := range km.rooms {
			if other.Info.ID != r.Info.ID {
				continue
			}
			if other.Info.Status < kq.KQStatusDone || other.Info.Status > kq.KQStatusNotFound {
				continue
			}
			if r.Info.ScheduleTime < other.Info.ScheduleTime { // exe: 0x004540E3
				r.Info.Status = kq.KQStatusDeleting // exe: 0x004540FD
				deletedAny = true
				break
			}
		}
	}
	if !deletedAny {
		return 0
	}
	removed := 0
	kept := km.rooms[:0]
	for i, r := range km.rooms {
		if r.Info.Status == kq.KQStatusDeleting { // exe: 0x00454126
			km.freeMapLink(r.Info.Handle)
			km.freeJoinerLocked(i)
			removed++
			continue
		}
		kept = append(kept, r)
	}
	km.rooms = kept
	return removed
}

func (km *KQManager) toClientInfo(r *kqRoom) protocol.ProtoKQInfoClient {
	return protocol.ProtoKQInfoClient{
		Handle:             r.Info.Handle,
		Status:             byte(r.Info.Status),
		NumOfJoiner:        r.Info.NumOfJoiner,
		ID:                 r.Info.ID,
		Title:              r.Info.TitleString(),
		LimitTime:          r.Info.LimitTime,
		StartTime:          time.Unix(int64(r.Info.ScheduleTime), 0),
		StartWaitTime:      r.Info.StartWaitTime,
		MinLevel:           r.Info.MinLevel,
		MaxLevel:           r.Info.MaxLevel,
		MinPlayers:         r.Info.MinPlayers,
		MaxPlayers:         r.Info.MaxPlayers,
		PlayerRepeatMode:   r.Info.PlayerRepeatMode,
		PlayerRepeatCount:  r.Info.PlayerRepeatCount,
		PlayerRevivalMode:  r.Info.PlayerRevivalMode,
		PlayerRevivalCount: r.Info.PlayerRevivalCount,
		DemandQuest:        r.Info.DemandQuest,
		DemandItem:         r.Info.DemandItem,
		DemandClass:        r.Info.DemandClass,
		DemandGender:       r.Info.DemandGender,
	}
}
