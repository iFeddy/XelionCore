package world

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/s2s"
)

func sessionCharName(client *network.Client) string {
	if client == nil {
		return ""
	}
	if client.CharName != "" {
		return client.CharName
	}
	return client.Username
}

// HandleList implementa CParserClient::fc_NC_KQ_LIST_REQ (WorldManager.exe @ 0x00413980).
func (km *KQManager) HandleList(client *network.Client) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	now := km.nowTime()
	filtered := make([]protocol.ProtoKQInfoClient, 0, len(km.rooms))
	startHandle := uint32(kqUnsetHandle) // exe: 0x00413A5C (-1 cuando lista vacía)
	endHandle := uint32(0)

	for _, r := range km.rooms {
		if r.Info.Status > kq.KQStatusRunning { // exe: 0x004139E2 (Status <= 4)
			continue
		}
		if len(filtered) == 0 {
			startHandle = r.Info.Handle // exe: 0x00413A93
		}
		endHandle = r.Info.Handle // exe: 0x00413A9D
		filtered = append(filtered, km.toClientInfo(r))
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_ACK(now, startHandle, endHandle, filtered))
}

// HandleSchedule implementa CParserClient::fc_NC_KQ_SCHEDULE_REQ (WorldManager.exe @ 0x00413E00).
func (km *KQManager) HandleSchedule(client *network.Client) {
	if client == nil {
		return
	}
	km.mu.RLock()
	defer km.mu.RUnlock()

	if len(km.rooms) == 0 { // exe: 0x00413E57
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_SCHEDULE_ACK(kqUnsetHandle, 0, nil))
		return
	}
	all := make([]protocol.ProtoKQInfoClient, 0, len(km.rooms))
	for _, r := range km.rooms {
		all = append(all, km.toClientInfo(r))
	}
	startHandle := km.rooms[0].Info.Handle            // exe: 0x00413EE0
	endHandle := km.rooms[len(km.rooms)-1].Info.Handle // exe: 0x00413EEE
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_SCHEDULE_ACK(startHandle, endHandle, all))
}

// HandleRefresh implementa CParserClient::fc_NC_KQ_LIST_REFRESH_REQ (WorldManager.exe @ 0x00413F80).
func (km *KQManager) HandleRefresh(client *network.Client) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st := km.getOrCreateSessionLocked(client)
	st.AutoRefreshActive = true     // exe: 0x00413FA0 ([esi+0x1DF68] = 1)
	st.LastRefreshAt = km.nowTime() // exe: 0x00413FB8 ([esi+0x1DF6C] = now)
	km.ackListRefreshLocked(st, st.LastRefreshAt)
}

// HasSessionTimeAck reporta si la sesión ya recibió LIST_TIME_ACK (primer 0x581B).
func (km *KQManager) HasSessionTimeAck(client *network.Client) bool {
	if client == nil {
		return false
	}
	km.mu.RLock()
	defer km.mu.RUnlock()
	st, ok := km.sessions[client]
	return ok && st.TimeAckSent
}

// ackListRefreshLocked porta CKQServer::Ack_NC_KQ_LIST_REFRESH (WorldManager.exe @ 0x004559E0).
func (km *KQManager) ackListRefreshLocked(st *kqSessionState, now time.Time) {
	if st == nil || st.Client == nil {
		return
	}
	// 1. Enviar LIST_TIME_ACK solo una vez por sesión (exe: 0x00455A2B)
	if !st.TimeAckSent {
		st.TimeAckSent = true
		_ = st.Client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_TIME_ACK(now))
	}

	// 2. Snapshot actual con Status <= 4 (exe: 0x00455A90)
	curSnap := make([]kq.PROTO_KQ_INFO, 0, len(km.rooms))
	curByHandle := make(map[uint32]kq.PROTO_KQ_INFO, len(km.rooms))
	for _, r := range km.rooms {
		if r.Info.Status <= kq.KQStatusRunning {
			curSnap = append(curSnap, r.Info)
			curByHandle[r.Info.Handle] = r.Info
		}
	}

	oldByHandle := make(map[uint32]kq.PROTO_KQ_INFO, len(st.Snapshot))
	for _, old := range st.Snapshot {
		oldByHandle[old.Handle] = old
	}

	// 3. Borrados + Actualizados comparando cache previo vs snapshot actual (exe: 0x00455ADE..0x00455CA0)
	var delHandles []uint32
	var updates []protocol.KQUpdateItem
	for _, old := range st.Snapshot {
		cur, exists := curByHandle[old.Handle]
		if !exists {
			delHandles = append(delHandles, old.Handle)
			continue
		}
		if cur.Status != old.Status || cur.NumOfJoiner != old.NumOfJoiner {
			updates = append(updates, protocol.KQUpdateItem{
				Handle:      cur.Handle,
				Status:      byte(cur.Status),
				NumOfJoiner: cur.NumOfJoiner,
			})
		}
	}
	if len(delHandles) > 0 {
		_ = st.Client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_DELETE_ACK(delHandles))
	}
	if len(updates) > 0 {
		_ = st.Client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_UPDATE_ACK(updates))
	}

	// 4. Añadidos nuevos con batching cada vez que (2 + n*141) >= 0x1D26 (n=53) (exe: 0x00455D5B)
	const protoKQClientSize = 141 // exe: 0x00455D4F (imul ecx, 0x8D)
	const addHeaderSize = 2       // exe: 0x00455D58
	var addBatch []protocol.ProtoKQInfoClient
	for _, cur := range curSnap {
		if _, existed := oldByHandle[cur.Handle]; existed {
			continue
		}
		r, _ := km.getRoomByHandle(cur.Handle)
		if r == nil {
			continue
		}
		addBatch = append(addBatch, km.toClientInfo(r))
		if addHeaderSize+len(addBatch)*protoKQClientSize >= kqListAddFlushBytes {
			_ = st.Client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_ADD_ACK(addBatch))
			addBatch = addBatch[:0]
		}
	}
	if len(addBatch) > 0 {
		_ = st.Client.SendPacket(protocol.New_PROTO_NC_KQ_LIST_ADD_ACK(addBatch))
	}

	// 5. Guardar snapshot en sesión (exe: 0x00455D90)
	st.Snapshot = curSnap
}

// autoRefreshListLocked porta CKQServer::AutoRefreshList (WorldManager.exe @ 0x00455DD0).
// Selecciona la única sesión activa con LastRefreshAt más antiguo y < now.
func (km *KQManager) autoRefreshListLocked(now time.Time) {
	var oldest *kqSessionState
	for _, st := range km.sessions {
		if st == nil || !st.AutoRefreshActive {
			continue
		}
		if !st.LastRefreshAt.Before(now) { // exe: 0x00455DF0
			continue
		}
		if oldest == nil || st.LastRefreshAt.Before(oldest.LastRefreshAt) {
			oldest = st
		}
	}
	if oldest != nil {
		oldest.LastRefreshAt = now // exe: 0x00455E29
		km.ackListRefreshLocked(oldest, now)
	}
}

// HandleStatus implementa CParserClient::fc_NC_KQ_STATUS_REQ (WorldManager.exe @ 0x00413B20).
func (km *KQManager) HandleStatus(client *network.Client, handle uint32) {
	if client == nil {
		return
	}
	km.mu.RLock()
	defer km.mu.RUnlock()

	r, _ := km.getRoomByHandle(handle)
	if r == nil { // exe: 0x00413BA4 (Status = 0x0A NotFound, NumOfJoiner = 0)
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_STATUS_ACK(handle, byte(kq.KQStatusNotFound), nil))
		return
	}
	names := make([]string, len(r.Joiners))
	for i, j := range r.Joiners {
		names[i] = j.CharName
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_STATUS_ACK(handle, byte(r.Info.Status), names))
}

// HandleJoin implementa CParserClient::fc_NC_KQ_JOIN_REQ (WorldManager.exe @ 0x00413CC0).
func (km *KQManager) HandleJoin(client *network.Client, handle uint32, level, class byte) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st := km.getOrCreateSessionLocked(client)
	if level > 0 {
		st.Level = level
	}
	if class > 0 {
		st.Class = class
	}

	if st.PrisonMin > 0 { // exe: 0x00413CF5 (word [charData+0x62] != 0 -> EldPri blocked)
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, kqErrJoinPrisonBlocked))
		return
	}
	if km.isPlayerAlreadyJoinedLocked(client, handle) { // exe: 0x00413D10
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, kqErrJoinAlreadyJoined))
		return
	}
	km.playerDisjoinLocked(client)                 // exe: 0x00413D2A
	ret := km.playerJoinLocked(client, handle, st) // exe: 0x00413D39
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_ACK(handle, kqErrJoinOK+uint16(ret)))
}

func (km *KQManager) isPlayerAlreadyJoinedLocked(client *network.Client, handle uint32) bool {
	// exe: 0x004524D0 (CKQServer::IsPlayerAlreadyJoin)
	st, ok := km.sessions[client]
	if !ok || st.HKQ != handle {
		return false
	}
	r, _ := km.getRoomByHandle(handle)
	if r == nil {
		return false
	}
	return km.getJoinerIndex(r, sessionCharName(client)) >= 0
}

// playerJoinLocked porta CKQServer::PlayerJoin (WorldManager.exe @ 0x004542F0).
// Retorna códigos 0..7 para sumarse a 0x0991 en PROTO_NC_KQ_JOIN_ACK.
func (km *KQManager) playerJoinLocked(client *network.Client, handle uint32, st *kqSessionState) int {
	r, _ := km.getRoomByHandle(handle)
	if r == nil {
		return 1 // exe: 0x0045430F -> 0x0992
	}
	if len(r.Joiners) >= kqMaxJoiners || len(r.Joiners) >= int(r.Info.MaxPlayers) { // exe: 0x00454334..0x0045434D
		return 2 // exe: 0x0045433B -> 0x0993
	}

	// CKQ::IsJoinable (exe: 0x0045CED0)
	if r.Info.Status != kq.KQStatusJoining { // exe: 0x0045CF20 -> 2 -> PlayerJoin returns 3 (0x0994)
		return 3
	}
	if st.Level < r.Info.MinLevel || st.Level > r.Info.MaxLevel { // exe: 0x0045CF28..0x0045CF30 -> 3 -> returns 4 (0x0995)
		return 4
	}
	classBit := int64(1) << st.Class
	if r.Info.DemandClass != 0 && (r.Info.DemandClass&classBit) == 0 { // exe: 0x0045CF42 -> 4 -> returns 5 (0x0996)
		return 5
	}
	genderBit := byte(1) << st.Gender
	if r.Info.DemandGender != 0 && (r.Info.DemandGender&genderBit) == 0 { // exe: 0x0045CF65 -> 5 -> returns 6 (0x0997)
		return 6
	}

	teamType := byte(kqTeamInvalid) // exe: 0x00454469 (default=2 cuando KQTeam==nil o DivideType!=2)
	if km.store != nil {
		if kt, ok := km.store.GetKQTeamData(r.Info.ID); ok && kt.KQTeamDivideType == 2 { // exe: 0x00454429 (KQTD_USERSELECT=2)
			if r.TeamCounts[0] <= r.TeamCounts[1] {
				teamType = 0
				r.TeamCounts[0]++
			} else {
				teamType = 1
				r.TeamCounts[1]++
			}
		}
	}

	joiner := &kqJoiner{
		CharRegNo: uint32(client.CharacterID),
		Slot:      uint16(client.Slot),
		Level:     st.Level,
		Class:     st.Class,
		CharName:  sessionCharName(client),
		TeamType:  teamType,
	}
	r.Joiners = append(r.Joiners, joiner)
	r.Info.NumOfJoiner = uint16(len(r.Joiners))
	st.HKQ = handle

	km.sendBroadcastJoinListLocked(r) // exe: 0x00454496
	return 0                          // exe: 0x0045449E -> 0x0991
}

// HandleJoinCancel implementa CParserClient::fc_NC_KQ_JOIN_CANCEL_REQ (WorldManager.exe @ 0x00413D70).
func (km *KQManager) HandleJoinCancel(client *network.Client, handle uint32) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.playerDisjoinLocked(client) { // exe: 0x00413DB5
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle, kqErrJoinCancelOK))
		return
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle, kqErrJoinCancelFail))
}

// playerDisjoinLocked porta CKQServer::PlayerDisjoin (WorldManager.exe @ 0x00454510).
func (km *KQManager) playerDisjoinLocked(client *network.Client) bool {
	if client == nil {
		return false
	}
	st, ok := km.sessions[client]
	if !ok || st.HKQ == kqUnsetHandle {
		return false
	}
	handle := st.HKQ
	r, _ := km.getRoomByHandle(handle)
	if r == nil {
		st.HKQ = kqUnsetHandle
		return false
	}
	idx := -1
	charRegNo := uint32(client.CharacterID)
	charName := sessionCharName(client)
	for i, j := range r.Joiners {
		if j.CharRegNo == charRegNo || (charName != "" && j.CharName == charName) {
			idx = i
			break
		}
	}
	if idx < 0 {
		st.HKQ = kqUnsetHandle
		return false
	}

	leaving := r.Joiners[idx]
	// Si era objetivo de votación activa y aún no fue baneado -> cancelar votación (exe: 0x00454610..0x00454779)
	if idx == r.Vote.Target && !leaving.IsBan {
		cancelCmd := protocol.New_PROTO_NC_KQ_VOTE_CANCEL_CMD(leaving.CharName)
		for i, j := range r.Joiners {
			j.CanVote = false
			if i != idx {
				km.sendToJoinerLocked(j, cancelCmd)
			}
		}
		km.clearVoteInfo(r)
	}

	if leaving.TeamType < 2 && r.TeamCounts[leaving.TeamType] > 0 { // exe: 0x00454796
		r.TeamCounts[leaving.TeamType]--
	}

	// Compactar array y ajustar índices de votación (exe: 0x004547BA..0x00454816)
	r.Joiners = append(r.Joiners[:idx], r.Joiners[idx+1:]...)
	r.Info.NumOfJoiner = uint16(len(r.Joiners))
	if r.Vote.Starter == idx {
		r.Vote.Starter = -1
	} else if r.Vote.Starter > idx {
		r.Vote.Starter--
	}
	if r.Vote.Target == idx {
		r.Vote.Target = -1
	} else if r.Vote.Target > idx {
		r.Vote.Target--
	}
	st.HKQ = kqUnsetHandle

	km.sendPlayerDisjoinCmdLocked(handle, charRegNo)
	km.sendBroadcastJoinListLocked(r)
	return true
}

// HandleJoinList implementa CKQServer::Recv_NC_KQ_JOIN_LIST_REQ (WorldManager.exe @ 0x004568C0).
func (km *KQManager) HandleJoinList(client *network.Client, handle uint32) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	r, _ := km.getRoomByHandle(handle)
	if r == nil { // exe: 0x0045696F
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(kqErrJoinListNotFound, nil))
		return
	}
	st := km.getOrCreateSessionLocked(client)
	now := km.nowTime()
	if now.Before(st.JoinListNextAt) { // exe: 0x004569E9
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(kqErrJoinListCooldown, nil))
		return
	}
	coolSec := km.getSingleDataValue("KQPlayerList_ResetListCoolTime", kqJoinListResetCoolTimeDefault)
	st.JoinListNextAt = now.Add(time.Duration(coolSec) * time.Second)

	chars := km.buildJoinListChars(r)
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(kqErrJoinListOK, chars))
}

func (km *KQManager) buildJoinListChars(r *kqRoom) []protocol.KQJoinCharInfo {
	out := make([]protocol.KQJoinCharInfo, len(r.Joiners))
	for i, j := range r.Joiners {
		out[i] = protocol.KQJoinCharInfo{
			Level:    j.Level,
			Class:    j.Class,
			CharName: j.CharName,
			TeamType: j.TeamType,
		}
	}
	return out
}

// sendBroadcastJoinListLocked porta CKQServer::SendBroadcast_NC_KQ_JOIN_LIST_ACK (exe: 0x00453BD0).
func (km *KQManager) sendBroadcastJoinListLocked(r *kqRoom) {
	pkt := protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(kqErrJoinListOK, km.buildJoinListChars(r))
	for _, j := range r.Joiners {
		km.sendToJoinerLocked(j, pkt)
	}
}

// HandleTeamSelect implementa CKQServer::Recv_NC_KQ_TEAM_SELECT_REQ (WorldManager.exe @ 0x004536C0).
func (km *KQManager) HandleTeamSelect(client *network.Client, teamType byte) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st, ok := km.sessions[client]
	if !ok || st.HKQ == kqUnsetHandle { // exe: 0x0045371E
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectNotJoined, teamType))
		return
	}
	r, _ := km.getRoomByHandle(st.HKQ)
	if r == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectNotJoined, teamType))
		return
	}
	charName := sessionCharName(client)
	jIdx := km.getJoinerIndex(r, charName)
	if jIdx < 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectNotJoined, teamType))
		return
	}
	if r.Info.Status != kq.KQStatusJoining { // exe: 0x00453753
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectBadStatus, teamType))
		return
	}
	if km.store == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectNoTeamCfg, teamType))
		return
	}
	kt, found := km.store.GetKQTeamData(r.Info.ID)
	if !found { // exe: 0x0045381A
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectNoTeamCfg, teamType))
		return
	}
	if kt.KQTeamDivideType != 2 { // exe: 0x00453831 (KQTD_USERSELECT)
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectBadDivide, teamType))
		return
	}
	joiner := r.Joiners[jIdx]
	if teamType >= 2 {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectBadDivide, teamType))
		return
	}
	if joiner.TeamType == teamType { // exe: 0x0045378C
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectSameTeam, teamType))
		return
	}
	maxPerTeam := int(r.Info.MaxPlayers) / 2 // exe: 0x004537AC
	if int(r.TeamCounts[teamType]) >= maxPerTeam {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectFull, teamType))
		return
	}
	// Simular cambio de bando y validar MaxMemberGap (exe: 0x004537E1)
	newTarget := int(r.TeamCounts[teamType]) + 1
	newOther := int(r.TeamCounts[1-teamType])
	if joiner.TeamType < 2 && joiner.TeamType != teamType && newOther > 0 {
		newOther--
	}
	gap := newTarget - newOther
	if gap > int(kt.MaxMemberGap) {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectGapOver, teamType))
		return
	}

	if joiner.TeamType < 2 && r.TeamCounts[joiner.TeamType] > 0 {
		r.TeamCounts[joiner.TeamType]--
	}
	joiner.TeamType = teamType
	r.TeamCounts[teamType]++

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(kqErrTeamSelectOK, teamType))

	// exe: 0x004538B0: difundir PROTO_NC_KQ_TEAM_SELECT_CMD al resto de la sala
	cmd := protocol.New_PROTO_NC_KQ_TEAM_SELECT_CMD(joiner.CharName, teamType)
	for i, j := range r.Joiners {
		if i != jIdx {
			km.sendToJoinerLocked(j, cmd)
		}
	}
}

// kqTeamCanKQStartLocked porta CKQServer::KQTeam_CanKQStart (WorldManager.exe @ 0x004553C0).
func (km *KQManager) kqTeamCanKQStartLocked(roomIdx int) bool {
	if km.store == nil || roomIdx < 0 || roomIdx >= len(km.rooms) {
		return true
	}
	r := km.rooms[roomIdx]
	kt, ok := km.store.GetKQTeamData(r.Info.ID)
	if !ok || kt.KQTeamDivideType != 2 { // exe: 0x004553E7
		return true
	}
	c0, c1 := r.TeamCounts[0], r.TeamCounts[1]
	if c0 == 0 || c1 == 0 { // exe: 0x004553F3..0x00455412
		km.setDoneSkipLocked(roomIdx, 2)
		return false
	}
	diff := int(c0) - int(c1)
	if diff < 0 {
		diff = -diff
	}
	if diff > int(kt.MaxMemberGap) { // exe: 0x0045543A
		km.setDoneSkipLocked(roomIdx, 3)
		return false
	}
	return true
}

// kqTeamDivideRandomLocked porta CKQServer::KQTeam_DivideRandom (WorldManager.exe @ 0x00453540).
func (km *KQManager) kqTeamDivideRandomLocked(r *kqRoom) {
	if km.store == nil || r == nil {
		return
	}
	kt, ok := km.store.GetKQTeamData(r.Info.ID)
	if !ok || kt.KQTeamDivideType != 1 { // exe: 0x0045356B (KQTD_RANDOM=1)
		return
	}
	r.TeamCounts = [2]byte{0, 0}
	perm := rand.Perm(len(r.Joiners))
	for n, jIdx := range perm {
		team := byte(n % 2)
		r.Joiners[jIdx].TeamType = team
		r.TeamCounts[team]++
	}
}

// kqTeamLeavePartyLocked porta CKQServer::KQTeam_LeaveParty (WorldManager.exe @ 0x00453490).
func (km *KQManager) kqTeamLeavePartyLocked(r *kqRoom) {
	if km.store == nil || km.onLeaveParty == nil {
		return
	}
	kt, ok := km.store.GetKQTeamData(r.Info.ID)
	if !ok || kt.IsTeamPVP == 0 {
		return
	}
	for _, j := range r.Joiners {
		if j.CharRegNo > 0 {
			km.onLeaveParty(int(j.CharRegNo))
		}
	}
}

// CanParty porta CWMClientSession::KQTeam_IsCanParty (WorldManager.exe @ 0x00433740).
func (km *KQManager) CanParty(a, b *network.Client) bool {
	if a == nil || b == nil {
		return false
	}
	km.mu.RLock()
	defer km.mu.RUnlock()

	sa, okA := km.sessions[a]
	sb, okB := km.sessions[b]
	var aRunning, bRunning bool
	var ra, rb *kqRoom
	if okA && sa.HKQ != kqUnsetHandle {
		ra, _ = km.getRoomByHandle(sa.HKQ)
		aRunning = ra != nil && ra.Info.Status == kq.KQStatusRunning
	}
	if okB && sb.HKQ != kqUnsetHandle {
		rb, _ = km.getRoomByHandle(sb.HKQ)
		bRunning = rb != nil && rb.Info.Status == kq.KQStatusRunning
	}
	if !aRunning && !bRunning { // exe: 0x00433792
		return true
	}
	ta := byte(kqTeamInvalid)
	if ra != nil {
		if idx := km.getJoinerIndex(ra, sessionCharName(a)); idx >= 0 {
			ta = ra.Joiners[idx].TeamType
		}
	}
	tb := byte(kqTeamInvalid)
	if rb != nil {
		if idx := km.getJoinerIndex(rb, sessionCharName(b)); idx >= 0 {
			tb = rb.Joiners[idx].TeamType
		}
	}
	return ta == tb // exe: 0x004337E5
}

// HandleVoteStartCheck implementa CKQServer::Recv_NC_KQ_VOTE_START_CHECK_REQ (WorldManager.exe @ 0x00453340).
func (km *KQManager) HandleVoteStartCheck(client *network.Client) {
	if client == nil {
		return
	}
	km.mu.RLock()
	defer km.mu.RUnlock()

	st, ok := km.sessions[client]
	if !ok || st.HKQ == kqUnsetHandle { // exe: 0x0045340C
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_CHECK_ACK(kqErrVoteCheckNotJoined))
		return
	}
	r, _ := km.getRoomByHandle(st.HKQ)
	if r == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_CHECK_ACK(kqErrVoteCheckNotJoined))
		return
	}
	if km.store == nil || !km.store.GetIsVote(r.Info.ID) { // exe: 0x0045346C
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_CHECK_ACK(kqErrVoteCheckDisabled))
		return
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_CHECK_ACK(kqErrVoteCheckOK))
}

// HandleVoteStart inicia una votación (compat wrapper sin tipo/motivo explícito -> voteType=0, reason="Kick").
func (km *KQManager) HandleVoteStart(client *network.Client, targetName string) {
	km.HandleVoteStartFull(client, targetName, 0, "Kick")
}

// HandleVoteStartFull implementa CKQServer::Recv_NC_KQ_VOTE_START_REQ (WorldManager.exe @ 0x00456310).
func (km *KQManager) HandleVoteStartFull(client *network.Client, targetName string, voteType byte, contents string) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st, ok := km.sessions[client]
	if !ok || st.HKQ == kqUnsetHandle { // exe: 0x004563AE
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartNotJoined))
		return
	}
	r, _ := km.getRoomByHandle(st.HKQ)
	if r == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartNotJoined))
		return
	}
	if r.Info.Status != kq.KQStatusRunning { // exe: 0x00456420
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartBadStatus))
		return
	}
	if km.inVoting(r) { // exe: 0x00456442
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartInVoting))
		return
	}
	if km.store == nil || !km.store.GetIsVote(r.Info.ID) { // exe: 0x0045646A
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartDisabled))
		return
	}
	starterName := sessionCharName(client)
	starterIdx := km.getJoinerIndex(r, starterName)
	targetIdx := km.getJoinerIndex(r, strings.TrimRight(targetName, "\x00"))
	if starterIdx < 0 { // exe: 0x004564C6
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartNotJoined))
		return
	}
	if targetIdx < 0 { // exe: 0x004564E3
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartBadTarget))
		return
	}
	starter := r.Joiners[starterIdx]
	target := r.Joiners[targetIdx]
	if starter.TeamType != target.TeamType || target.IsBan { // exe: 0x00456512..0x004565A5
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartBadTarget))
		return
	}
	tClient := km.onlineByCharID[int(target.CharRegNo)]
	if tClient == nil { // exe: 0x00456568
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartBadTarget))
		return
	}
	if starterIdx == targetIdx { // exe: 0x004565C7
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartSelf))
		return
	}
	if len(contents) == 0 { // exe: 0x004565EB
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartNoReason))
		return
	}
	now := km.nowTime()
	if now.Before(st.VoteSuggestNextAt) { // exe: 0x0045661B
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartCooldown))
		return
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_START_ACK(kqErrVoteStartOK))

	suggestCool := km.getSingleDataValue("KQVote_SuggestCoolTime", kqVoteSuggestCoolTimeDefault)
	voteLimit := km.getSingleDataValue("KQVote_VoteLimitTime", kqVoteLimitTimeDefaultSec)
	st.VoteSuggestNextAt = now.Add(time.Duration(suggestCool) * time.Second)
	endTime := now.Add(time.Duration(voteLimit) * time.Second)

	km.clearVoteInfo(r) // exe: 0x004566D5
	r.Vote = kqVoteState{
		TeamType: starter.TeamType,
		Starter:  starterIdx,
		Target:   targetIdx,
		EndTime:  endTime,
	}
	cmd := protocol.New_PROTO_NC_KQ_VOTE_VOTING_CMD(starter.CharName, target.CharName, voteType, endTime, contents)
	for i, j := range r.Joiners {
		if i == starterIdx || i == targetIdx || j.TeamType != r.Vote.TeamType {
			continue
		}
		if km.sendToJoinerLocked(j, cmd) {
			j.CanVote = true     // exe: 0x0045685D
			r.Vote.CancelNum++   // exe: 0x00456864
		}
	}
}

// HandleVoteCast implementa CKQServer::Recv_NC_KQ_VOTE_VOTING_REQ (WorldManager.exe @ 0x00452EE0).
func (km *KQManager) HandleVoteCast(client *network.Client, agree bool) {
	if agree {
		km.HandleVoteCastType(client, 1) // KQ_VOTING_YES
	} else {
		km.HandleVoteCastType(client, 2) // KQ_VOTING_NO
	}
}

// HandleVoteCastType recibe eVoteType exacto (1=YES, 2=NO).
func (km *KQManager) HandleVoteCastType(client *network.Client, voteType uint32) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st, ok := km.sessions[client]
	if !ok || st.HKQ == kqUnsetHandle { // exe: 0x00452F7D
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastNotJoined))
		return
	}
	r, _ := km.getRoomByHandle(st.HKQ)
	if r == nil {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastNotJoined))
		return
	}
	jIdx := km.getJoinerIndex(r, sessionCharName(client))
	if jIdx < 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastNotJoined))
		return
	}
	joiner := r.Joiners[jIdx]
	if !km.inVoting(r) || joiner.TeamType != r.Vote.TeamType { // exe: 0x00453001
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastNotVoting))
		return
	}
	if !joiner.CanVote { // exe: 0x00453044
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastCannotVote))
		return
	}

	joiner.CanVote = false // exe: 0x00453065
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_VOTING_ACK(kqErrVoteCastOK))

	switch voteType {
	case 1: // exe: 0x0045308D (YES)
		if r.Vote.CancelNum > 0 {
			r.Vote.CancelNum--
		}
		r.Vote.YesNum++
	case 2: // exe: 0x0045307F (NO)
		if r.Vote.CancelNum > 0 {
			r.Vote.CancelNum--
		}
		r.Vote.NoNum++
	}
}

// voteProcessingLocked porta CKQServer::VoteProcessing (WorldManager.exe @ 0x00454F60).
func (km *KQManager) voteProcessingLocked(now time.Time) {
	for _, r := range km.rooms {
		for _, j := range r.Joiners {
			if j.BanForce { // exe: 0x00454FF5
				km.sendLinkToForceByBanLocked(r, j)
			}
		}
		if !km.inVoting(r) {
			continue
		}
		if now.Before(r.Vote.EndTime) || now.Equal(r.Vote.EndTime) { // exe: 0x0045506A..0x0045506D (EndTime < now)
			continue
		}
		if r.Vote.Target < 0 || r.Vote.Target >= len(r.Joiners) {
			km.clearVoteInfo(r)
			continue
		}
		target := r.Joiners[r.Vote.Target]
		reqRate := byte(0)
		if km.store != nil {
			if rate, ok := km.store.GetVoteAgreeRate(int(target.VoteCount)); ok {
				reqRate = rate
			}
		}
		agreePct := byte(0)
		totalCast := int(r.Vote.YesNum) + int(r.Vote.NoNum)
		if r.Vote.YesNum > 0 && totalCast > 0 { // exe: 0x004550C7..0x004550D8
			agreePct = byte((int(r.Vote.YesNum) * kqPercentBase) / totalCast)
		}
		target.VoteCount++ // exe: 0x004550E8

		passed := agreePct >= reqRate // exe: 0x004550DA
		var resultPkt *protocol.PacketWriter
		if passed {
			resultPkt = protocol.New_PROTO_NC_KQ_VOTE_RESULT_SUC_CMD(target.CharName, reqRate, r.Vote.YesNum, r.Vote.NoNum, r.Vote.CancelNum)
		} else {
			resultPkt = protocol.New_PROTO_NC_KQ_VOTE_RESULT_FAIL_CMD(target.CharName, r.Vote.YesNum, r.Vote.NoNum, r.Vote.CancelNum)
		}
		for _, j := range r.Joiners {
			j.CanVote = false
			if j.TeamType == r.Vote.TeamType { // exe: 0x004551FA
				km.sendToJoinerLocked(j, resultPkt)
			}
		}
		if passed { // exe: 0x00455266..0x0045532B
			target.IsBan = true
			target.BanForce = true
			target.BanAt = now
			banMsg := protocol.New_PROTO_NC_KQ_VOTE_BAN_MSG_CMD(reqRate, r.Vote.YesNum, r.Vote.NoNum, r.Vote.CancelNum)
			km.sendToJoinerLocked(target, banMsg)
		}
		km.clearVoteInfo(r)
	}
}

// HandleZ2WMakeAck implementa CParserZone::fc_NC_KQ_Z2W_MAKE_ACK (WorldManager.exe @ 0x0042A9C0).
func (km *KQManager) HandleZ2WMakeAck(handle uint32, errCode uint16) {
	km.mu.Lock()
	defer km.mu.Unlock()
	if errCode == kqErrZ2WMakeSuccess { // exe: 0x0042A9FA (0x0981)
		km.setJoiningLocked(handle) // exe: 0x0042AA0E
	} else {
		km.setNoMapBFLocked(handle) // exe: 0x0042AA22
	}
}

// HandleZ2WEndCmd implementa CParserZone::fc_NC_KQ_Z2W_END_CMD (WorldManager.exe @ 0x0042AAC0).
func (km *KQManager) HandleZ2WEndCmd(handle uint32) bool {
	km.mu.Lock()
	defer km.mu.Unlock()
	return km.setDoneLocked(handle)
}

// SendJoiningAlarmList porta CKQServer::Send_NC_KQ_JOINING_ALARM_LIST (WorldManager.exe @ 0x00452B60).
func (km *KQManager) SendJoiningAlarmList(client *network.Client) {
	if client == nil {
		return
	}
	km.mu.RLock()
	defer km.mu.RUnlock()

	var items []protocol.KQJoiningAlarmItem
	payloadSize := 2
	for _, r := range km.rooms {
		if r.Info.Status != kq.KQStatusJoining { // exe: 0x00452BDB
			continue
		}
		if payloadSize >= kqJoiningAlarmListMaxByte { // exe: 0x00452C38
			break
		}
		items = append(items, protocol.KQJoiningAlarmItem{
			Handle:   r.Info.Handle,
			ID:       r.Info.ID,
			MinLevel: r.Info.MinLevel,
			MaxLevel: r.Info.MaxLevel,
		})
		payloadSize += 8
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOINING_ALARM_LIST(items))
}

// IsExisted porta CKQServer::IsExisted (WorldManager.exe @ 0x00452580).
func (km *KQManager) IsExisted(handle uint32, charName string, scheduleTime uint32) bool {
	km.mu.RLock()
	defer km.mu.RUnlock()
	r, _ := km.getRoomByHandle(handle)
	if r == nil || uint32(r.Info.ScheduleTime) != scheduleTime {
		return false
	}
	return km.getJoinerIndex(r, charName) >= 0
}

// OnLinkEnd porta CParserZone::fc_NC_MAP_LINKEND_CMD KQ hooks (exe: 0x0042DCD4 + 0x0042DD27 + 0x00456B20).
func (km *KQManager) OnLinkEnd(client *network.Client, mapName string) {
	if client == nil {
		return
	}
	km.mu.Lock()
	defer km.mu.Unlock()

	st := km.getOrCreateSessionLocked(client)
	// Send_NC_KQ_JOINING_ALARM_LIST (exe: 0x0042DCD4)
	var items []protocol.KQJoiningAlarmItem
	payloadSize := 2
	for _, r := range km.rooms {
		if r.Info.Status != kq.KQStatusJoining {
			continue
		}
		if payloadSize >= kqJoiningAlarmListMaxByte {
			break
		}
		items = append(items, protocol.KQJoiningAlarmItem{
			Handle:   r.Info.Handle,
			ID:       r.Info.ID,
			MinLevel: r.Info.MinLevel,
			MaxLevel: r.Info.MaxLevel,
		})
		payloadSize += 8
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOINING_ALARM_LIST(items))

	// LinkEndProcessing (exe: 0x00456B20)
	if st.BannedInLogin { // exe: 0x00456B3F
		_ = client.SendPacket(protocol.New_PROTO_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD())
		st.BannedInLogin = false
	}
	if st.HKQ == kqUnsetHandle || mapName == "" {
		return
	}
	r, _ := km.getRoomByHandle(st.HKQ)
	if r == nil {
		km.playerDisjoinLocked(client)
		return
	}
	// LinkOutKQ (exe: 0x00454DC0): si está Running y mapName no coincide con ningún MapName, sale de KQ
	// LinkInKQ (exe: 0x00454E50): si está Running y mapName coincide con algún MapName, envía JOIN_LIST_ACK
	if r.Info.Status != kq.KQStatusRunning {
		return
	}
	inKQMap := false
	for m := 0; m < kqMaxMapsPerRoom; m++ {
		if r.Info.MapLink[m].MapNameString() == mapName {
			inKQMap = true
			break
		}
	}
	if !inKQMap { // exe: 0x00456B66 (LinkOutKQ == 1)
		km.playerDisjoinLocked(client)
		return
	}
	// exe: 0x00456B82 (LinkInKQ == 1 -> Recv_NC_KQ_JOIN_LIST_REQ)
	_ = client.SendPacket(protocol.New_PROTO_NC_KQ_JOIN_LIST_ACK(kqErrJoinListOK, km.buildJoinListChars(r)))
}

func (km *KQManager) sendToJoinerLocked(j *kqJoiner, pkt *protocol.PacketWriter) bool {
	if j == nil || pkt == nil {
		return false
	}
	c := km.onlineByCharID[int(j.CharRegNo)]
	if c == nil {
		return false
	}
	return c.SendPacket(pkt) == nil
}

func (km *KQManager) sendNotifyToJoinersLocked(r *kqRoom, msg string) {
	// exe: 0x004525D0 / 0x00452750 (Send_NC_KQ_NOTIFY_CMD_ToJoiner)
	if msg == "" {
		return
	}
	pkt := protocol.New_PROTO_NC_KQ_NOTIFY_CMD(msg)
	for _, j := range r.Joiners {
		km.sendToJoinerLocked(j, pkt)
	}
}

func (km *KQManager) sendJoiningAlarmEndLocked(r *kqRoom) {
	// exe: 0x00452A70 (Send_NC_KQ_JOINING_ALARM_END_CMD)
	pkt := protocol.New_PROTO_NC_KQ_JOINING_ALARM_END_CMD(r.Info.Handle, r.Info.ID)
	for _, j := range r.Joiners {
		km.sendToJoinerLocked(j, pkt)
	}
}

func (km *KQManager) sendW2ZMakeReqLocked(r *kqRoom) {
	// exe: 0x00454B50 (Send_NC_KQ_W2Z_MAKE_REQ)
	if km.s2sSink == nil || r == nil {
		return
	}
	raw, err := r.Info.MarshalBinary()
	if err != nil {
		return
	}
	km.s2sSink.Broadcast(s2s.Event{
		Symbol:  "wms_NC_KQ_W2Z_MAKE_REQ",
		WorldID: km.worldID,
		Payload: map[string]interface{}{
			"handle":       r.Info.Handle,
			"id":           r.Info.ID,
			"kq_info_b64":  base64.StdEncoding.EncodeToString(raw),
		},
	})
}

func (km *KQManager) sendW2ZStartCmdLocked(r *kqRoom) {
	// exe: 0x00454A10 (Send_NC_KQ_W2Z_START_CMD)
	if km.s2sSink == nil || r == nil {
		return
	}
	raw, err := r.Info.MarshalBinary()
	if err != nil {
		return
	}
	joiners := make([]map[string]interface{}, len(r.Joiners))
	for i, j := range r.Joiners {
		joiners[i] = map[string]interface{}{
			"char_reg_num": j.CharRegNo,
			"char_name":    j.CharName,
			"team_type":    j.TeamType,
		}
	}
	km.s2sSink.Broadcast(s2s.Event{
		Symbol:  "wms_NC_KQ_W2Z_START_CMD",
		WorldID: km.worldID,
		Payload: map[string]interface{}{
			"handle":      r.Info.Handle,
			"id":          r.Info.ID,
			"kq_info_b64": base64.StdEncoding.EncodeToString(raw),
			"joiners":     joiners,
		},
	})
}

func (km *KQManager) sendW2ZDestroyCmdLocked(handle uint32) {
	// exe: 0x00452880 (Send_NC_KQ_W2Z_DESTROY_CMD)
	if km.s2sSink == nil {
		return
	}
	km.s2sSink.Broadcast(s2s.Event{
		Symbol:  "wms_NC_KQ_W2Z_DESTROY_CMD",
		WorldID: km.worldID,
		Payload: map[string]interface{}{
			"handle": handle,
		},
	})
}

func (km *KQManager) sendPlayerDisjoinCmdLocked(handle, charRegNo uint32) {
	// exe: 0x00453B10 (send_NC_KQ_PLAYER_DISJOIN_CMD)
	if km.s2sSink == nil {
		return
	}
	km.s2sSink.Broadcast(s2s.Event{
		Symbol:  "wms_NC_KQ_PLAYER_DISJOIN_CMD",
		WorldID: km.worldID,
		Payload: map[string]interface{}{
			"handle":       handle,
			"char_reg_num": charRegNo,
		},
	})
}

func (km *KQManager) sendLinkToForceByBanLocked(r *kqRoom, j *kqJoiner) {
	// exe: 0x004531A0 (Send_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD)
	if km.s2sSink == nil || r == nil || j == nil {
		return
	}
	outMaps := [4]string{
		r.Info.MapLink[0].MapNameString(),
		r.Info.MapLink[1].MapNameString(),
		r.Info.MapLink[2].MapNameString(),
		r.Info.MapLink[3].MapNameString(),
	}
	km.s2sSink.Broadcast(s2s.Event{
		Symbol:  "wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD",
		WorldID: km.worldID,
		Payload: map[string]interface{}{
			"handle":       r.Info.Handle,
			"char_reg_num": j.CharRegNo,
			"char_name":    j.CharName,
			"out_maps":     outMaps,
		},
	})
}

// ============================================================================
// Sub-fase 27D: Restauración de posición previa / KQMap en Login
// (Character.exe @ 0x0045A130 p_Char_GetKQMap, 0x0043903B PROTO_NC_CHAR_KQMAP_CMD,
//  WorldManager.exe @ 0x0040AA10 fc_NC_CHAR_KQMAP_CMD)
// ============================================================================

// BuildProtoCharKQMapCmd serializa PROTO_NC_CHAR_KQMAP_CMD (0x101A, 28 bytes payload)
// idéntico a Character.exe @ 0x0043903B..0x004391A9.
func BuildProtoCharKQMapCmd(info *database.CharKQMapInfo) (*protocol.PacketWriter, protocol.PROTO_NC_CHAR_KQMAP_CMD_DATA) {
	var data protocol.PROTO_NC_CHAR_KQMAP_CMD_DATA
	if info != nil {
		data.NKQHandle = info.KQHandle
		copy(data.SKQMapName[:], info.KQMapName)
		data.NKQCoord = uint64(info.KQX) | (uint64(info.KQY) << 32)
		data.DKQDate = info.KQDatePacked
		if data.DKQDate == 0 && !info.KQDate.IsZero() {
			data.DKQDate = database.PackKQDateTime(info.KQDate) // exe: Character.exe @ 0x0043911F..0x00439199
		}
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_KQMAP_CMD) // 0x101A
	pw.WriteUint32(data.NKQHandle)
	pw.WriteBytes(data.SKQMapName[:])
	pw.WriteUint64(data.NKQCoord)
	pw.WriteUint32(data.DKQDate)
	return pw, data
}

// InKQStatusRunning verifica si una instancia KQ por handle está actualmente en estado RUNNING
// (exe: WorldManager.exe @ 0x00452D98).
func (km *KQManager) InKQStatusRunning(handle uint32) bool {
	if km == nil || handle == 0 || handle == kqUnsetHandle {
		return false
	}
	km.mu.RLock()
	defer km.mu.RUnlock()
	r, _ := km.getRoomByHandle(handle)
	return r != nil && r.Info.Status == kq.KQStatusRunning
}

// IsCharActiveKQJoiner verifica que el personaje siga siendo participante activo y no expulsado
// de una KQ actualmente en ejecución.
func (km *KQManager) IsCharActiveKQJoiner(handle uint32, charName string) bool {
	if km == nil || handle == 0 || handle == kqUnsetHandle || charName == "" {
		return false
	}
	km.mu.RLock()
	defer km.mu.RUnlock()
	r, _ := km.getRoomByHandle(handle)
	if r == nil || r.Info.Status != kq.KQStatusRunning {
		return false
	}
	idx := km.getJoinerIndex(r, charName)
	if idx < 0 {
		return false
	}
	return !r.Joiners[idx].IsBan
}

// ResolveCharKQLoginLocation ejecuta p_Char_GetKQMap (Character.exe @ 0x0045A130) y aplica
// la lógica de CParserCharacter::fc_NC_CHAR_KQMAP_CMD (WorldManager.exe @ 0x0040AA10):
//   - Si el personaje tiene nKQHandle/sKQMap activos en una KQ RUNNING y no fue expulsado,
//     restaura su entrada en sKQMap (nKQX, nKQY).
//   - Si la KQ ya terminó o fue expulsado, restaura la posición normal previa (sLoginZone, fallbackX, fallbackY).
func (km *KQManager) ResolveCharKQLoginLocation(
	charRepo database.CharacterAvatarRepository,
	charID int,
	charName string,
	fallbackMap string,
	fallbackX, fallbackY int,
) (destMap string, destX, destY int, kqMapPkt *protocol.PacketWriter, restoredToKQ bool) {
	destMap = fallbackMap
	destX = fallbackX
	destY = fallbackY
	if charRepo == nil || charID <= 0 {
		return destMap, destX, destY, nil, false
	}
	info, err := charRepo.GetCharKQMap(charID) // exe: Character.exe @ 0x0045A130 (p_Char_GetKQMap)
	if err != nil || info == nil {
		return destMap, destX, destY, nil, false
	}
	kqMapPkt, _ = BuildProtoCharKQMapCmd(info) // exe: Character.exe @ 0x0043903B (0x101A)

	kqMap := strings.TrimSpace(info.KQMapName)
	if km != nil && info.KQHandle != 0 && info.KQHandle != kqUnsetHandle && kqMap != "" {
		if km.IsCharActiveKQJoiner(info.KQHandle, charName) {
			return kqMap, int(info.KQX), int(info.KQY), kqMapPkt, true
		}
	}
	return destMap, destX, destY, kqMapPkt, false
}
