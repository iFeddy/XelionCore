package world_test

import (
	"encoding/base64"
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/s2s"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

type recordingS2SSink struct {
	mu     sync.Mutex
	events []s2s.Event
}

func (s *recordingS2SSink) Broadcast(ev s2s.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, ev)
}

func (s *recordingS2SSink) BroadcastExceptZone(ev s2s.Event, _ uint16) { s.Broadcast(ev) }
func (s *recordingS2SSink) SendToZone(_ uint16, ev s2s.Event) bool    { s.Broadcast(ev); return true }

func (s *recordingS2SSink) popAll() []s2s.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]s2s.Event, len(s.events))
	copy(out, s.events)
	s.events = nil
	return out
}

type capturedClient struct {
	client *network.Client
	peer   net.Conn
	mu     sync.Mutex
	pkts   []*protocol.Packet
}

func newCapturedClient(t *testing.T, charID int, charName string) *capturedClient {
	t.Helper()
	srvConn, cliConn := net.Pipe()
	c, err := network.NewClient(srvConn, nil)
	if err != nil {
		t.Fatalf("network.NewClient: %v", err)
	}
	c.CharacterID = charID
	c.CharName = charName
	c.Username = charName

	cc := &capturedClient{client: c, peer: cliConn}
	go func() {
		for {
			pkt, err := readServerPacket(cliConn)
			if err != nil {
				return
			}
			cc.mu.Lock()
			cc.pkts = append(cc.pkts, pkt)
			cc.mu.Unlock()
		}
	}()
	t.Cleanup(func() {
		c.Close()
		_ = cliConn.Close()
	})
	return cc
}

func (cc *capturedClient) drain(t *testing.T) []*protocol.Packet {
	t.Helper()
	deadline := time.Now().Add(40 * time.Millisecond)
	for time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	out := make([]*protocol.Packet, len(cc.pkts))
	copy(out, cc.pkts)
	cc.pkts = nil
	return out
}

func findPacket(pkts []*protocol.Packet, op uint16) []*protocol.Packet {
	var out []*protocol.Packet
	for _, p := range pkts {
		if p.OpCode == op {
			out = append(out, p)
		}
	}
	return out
}

func TestCKQServer_ScheduleAndDoSchedule(t *testing.T) {
	ds, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	km := world.NewKQManager(ds, nil)
	if got := km.TemplateCount(); got != 57 {
		t.Fatalf("expected 57 rows from KingdomQuest.shn, got %d", got)
	}
	// Each template generates 2 upcoming slots => 114 rooms (> 50, <= 300)
	if got := km.QuestCount(); got < 50 || got > 300 {
		t.Fatalf("unexpected scheduled KQ count: %d", got)
	}
	// Re-running DoSchedule inside same minute must add 0 duplicates
	now := time.Now()
	if added := km.DoSchedule(now); added != 0 {
		t.Fatalf("expected 0 duplicates within same minute, got %d", added)
	}
}

func TestCKQServer_ListRefreshDiffAndBatch53(t *testing.T) {
	ds, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	km := world.NewKQManager(ds, nil)
	c1 := newCapturedClient(t, 501, "HeroOne")
	km.RegisterOnline(c1.client, 501)

	// 1. HandleList (0x5801 -> 0x5802)
	km.HandleList(c1.client)
	listPkts := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_LIST_ACK)
	if len(listPkts) != 1 {
		t.Fatalf("expected 1 PROTO_NC_KQ_LIST_ACK, got %d", len(listPkts))
	}
	r := protocol.NewPacketReader(listPkts[0].Raw)
	_, _ = r.ReadBytes(2) // opcode
	_, _ = r.ReadBytes(40) // serverTime + tm
	startH, _ := r.ReadUint32()
	endH, _ := r.ReadUint32()
	cnt, _ := r.ReadUint16()
	if cnt < 50 || startH == 0xFFFFFFFF || endH < startH {
		t.Fatalf("unexpected LIST_ACK header: start=%d end=%d cnt=%d", startH, endH, cnt)
	}
	if len(listPkts[0].Raw)-2 != 50+int(cnt)*141 {
		t.Fatalf("LIST_ACK payload length mismatch: got %d, want %d", len(listPkts[0].Raw)-2, 50+int(cnt)*141)
	}

	// 2. First HandleRefresh -> must emit 1x LIST_TIME_ACK (0x581C) + chunked LIST_ADD_ACK (0x581D, batch size 53 => 7475B)
	km.HandleRefresh(c1.client)
	pkts1 := c1.drain(t)
	timeAcks := findPacket(pkts1, protocol.PROTO_NC_KQ_LIST_TIME_ACK)
	if len(timeAcks) != 1 || len(timeAcks[0].Raw)-2 != 40 {
		t.Fatalf("expected 1 LIST_TIME_ACK (40B payload), got %d", len(timeAcks))
	}
	addAcks := findPacket(pkts1, protocol.PROTO_NC_KQ_LIST_ADD_ACK)
	if len(addAcks) < 2 {
		t.Fatalf("expected >=2 chunked LIST_ADD_ACK for %d quests (53 per chunk), got %d", cnt, len(addAcks))
	}
	firstAddCnt := binary.LittleEndian.Uint16(addAcks[0].Raw[2:4])
	if firstAddCnt != 53 {
		t.Fatalf("expected first LIST_ADD_ACK batch of 53 (2 + 53*141 = 7475 >= 0x1D26), got %d", firstAddCnt)
	}

	// 3. Second HandleRefresh -> MUST NOT resend LIST_TIME_ACK; if unchanged, 0 ADD/UPDATE/DELETE
	km.HandleRefresh(c1.client)
	pkts2 := c1.drain(t)
	if len(findPacket(pkts2, protocol.PROTO_NC_KQ_LIST_TIME_ACK)) != 0 {
		t.Fatalf("LIST_TIME_ACK must only be sent once per session")
	}
	if len(findPacket(pkts2, protocol.PROTO_NC_KQ_LIST_ADD_ACK)) != 0 {
		t.Fatalf("unexpected LIST_ADD_ACK on unchanged snapshot")
	}

	// 4. Transition a room and re-refresh -> must emit UPDATE (0x581F) or DELETE (0x581E)
	sink := &recordingS2SSink{}
	km.SetS2SSink(sink, 0)
	// Fast-forward clock to trigger MakeRoom on at least one scheduled room, then Z2W_MAKE_ACK -> Joining
	km.HandleZ2WMakeAck(startH, 0x0981)
	km.HandleRefresh(c1.client)
	pkts3 := c1.drain(t)
	updPkts := findPacket(pkts3, protocol.PROTO_NC_KQ_LIST_UPDATE_ACK)
	if len(updPkts) != 1 {
		t.Fatalf("expected 1 LIST_UPDATE_ACK after Status transition, got %d", len(updPkts))
	}
	updCnt := binary.LittleEndian.Uint16(updPkts[0].Raw[2:4])
	updH := binary.LittleEndian.Uint32(updPkts[0].Raw[4:8])
	updSt := updPkts[0].Raw[8]
	if updCnt != 1 || updH != startH || updSt != byte(kq.KQStatusJoining) {
		t.Fatalf("unexpected UPDATE item: cnt=%d h=%d status=%d", updCnt, updH, updSt)
	}

	// Complete the room -> Status becomes Done (5) -> hidden from snapshot -> refresh emits LIST_DELETE_ACK (0x581E)
	if !km.HandleZ2WEndCmd(startH) {
		t.Fatalf("HandleZ2WEndCmd returned false")
	}
	km.HandleRefresh(c1.client)
	pkts4 := c1.drain(t)
	delPkts := findPacket(pkts4, protocol.PROTO_NC_KQ_LIST_DELETE_ACK)
	if len(delPkts) != 1 {
		t.Fatalf("expected 1 LIST_DELETE_ACK after Done transition, got %d", len(delPkts))
	}
	if gotH := binary.LittleEndian.Uint32(delPkts[0].Raw[4:8]); gotH != startH {
		t.Fatalf("expected deleted handle %d, got %d", startH, gotH)
	}
}

func TestCKQServer_LifecycleWithSimulatedClock(t *testing.T) {
	ds, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	km := world.NewKQManager(ds, nil)
	sink := &recordingS2SSink{}
	km.SetS2SSink(sink, 0)

	// Find earliest scheduled room
	c1 := newCapturedClient(t, 601, "Alpha")
	c2 := newCapturedClient(t, 602, "Bravo")
	km.RegisterOnline(c1.client, 601)
	km.RegisterOnline(c2.client, 602)

	// Read HandleSchedule (0x580A) to pick a real room
	km.HandleSchedule(c1.client)
	schedPkts := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_SCHEDULE_ACK)
	if len(schedPkts) != 1 {
		t.Fatalf("expected 1 SCHEDULE_ACK")
	}
	// First entry at payload offset 10 (raw[12:])
	targetHandle := binary.LittleEndian.Uint32(schedPkts[0].Raw[12:16])
	startUnix := int64(int32(binary.LittleEndian.Uint32(schedPkts[0].Raw[12+0x4B : 12+0x4F])))
	schedT := time.Unix(startUnix, 0)

	// 1. Tick at schedT => DoSetMakeRoom emits wms_NC_KQ_W2Z_MAKE_REQ with 377B PROTO_KQ_INFO
	km.DoKingdomQuest(schedT)
	s2sEvs := sink.popAll()
	var makeFound bool
	for _, ev := range s2sEvs {
		if ev.Symbol == "wms_NC_KQ_W2Z_MAKE_REQ" {
			if h, ok := ev.Payload["handle"].(uint32); ok && h == targetHandle {
				b64, _ := ev.Payload["kq_info_b64"].(string)
				rawInfo, err := base64.StdEncoding.DecodeString(b64)
				if err != nil || len(rawInfo) != kq.SizeOfProtoKQInfo {
					t.Fatalf("invalid W2Z_MAKE_REQ kq_info_b64 len=%d err=%v", len(rawInfo), err)
				}
				var decoded kq.PROTO_KQ_INFO
				if err := decoded.UnmarshalBinary(rawInfo); err != nil {
					t.Fatalf("UnmarshalBinary(377B): %v", err)
				}
				if decoded.Status != kq.KQStatusMakingRoom {
					t.Fatalf("expected KQStatusMakingRoom, got %d", decoded.Status)
				}
				makeFound = true
			}
		}
	}
	if !makeFound {
		t.Fatalf("expected wms_NC_KQ_W2Z_MAKE_REQ for handle %d", targetHandle)
	}

	// 2. Simulate Z2W_MAKE_ACK(0x0981) -> Joining + broadcasts 0x5824 to online players
	km.HandleZ2WMakeAck(targetHandle, 0x0981)
	alarmPkts := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_JOINING_ALARM_CMD)
	if len(alarmPkts) == 0 {
		t.Fatalf("expected PROTO_NC_KQ_JOINING_ALARM_CMD after Z2W_MAKE_ACK(0x0981)")
	}
	_ = c2.drain(t)

	// 3. Test DoneSkip when MinPlayers not met after StartWaitTime:
	// First verify Status is Joining (2)
	km.HandleStatus(c1.client, targetHandle)
	stPkts := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_STATUS_ACK)
	if len(stPkts) != 1 || stPkts[0].Raw[6] != byte(kq.KQStatusJoining) {
		t.Fatalf("expected Status=Joining (2)")
	}

	// Extract StartWaitTime from SCHEDULE_ACK entry (offset 0x73)
	startWaitSec := binary.LittleEndian.Uint16(schedPkts[0].Raw[12+0x73 : 12+0x75])
	minPlayers := binary.LittleEndian.Uint16(schedPkts[0].Raw[12+0x77 : 12+0x79])

	// Join one player (most KQs require >= 2..5 min players) so player gets cancel notify when DoneSkip fires
	// Use level=20, class=1 (will adjust if needed by inspecting Min/Max)
	minLv := schedPkts[0].Raw[12+0x75]
	maxLv := schedPkts[0].Raw[12+0x76]
	testLv := minLv
	if testLv < 1 {
		testLv = 1
	}
	if maxLv > 0 && testLv > maxLv {
		testLv = maxLv
	}
	km.SetSessionCharacterStats(c1.client, testLv, 1, 1, 0)
	km.HandleJoin(c1.client, targetHandle, testLv, 1)
	joinAck := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_JOIN_ACK)
	if len(joinAck) == 1 && minPlayers > 1 {
		errCode := binary.LittleEndian.Uint16(joinAck[0].Raw[6:8])
		if errCode == 0x0991 {
			km.DoKingdomQuest(schedT.Add(time.Duration(startWaitSec+1) * time.Second))
			c1Pkts := c1.drain(t)
			if len(findPacket(c1Pkts, protocol.PROTO_NC_KQ_NOTIFY_CMD)) == 0 {
				t.Fatalf("expected PROTO_NC_KQ_NOTIFY_CMD cancellation notice on DoneSkip")
			}
			km.HandleStatus(c1.client, targetHandle)
			st2 := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_STATUS_ACK)
			if len(st2) != 1 || st2[0].Raw[6] != byte(kq.KQStatusDoneSkip) {
				t.Fatalf("expected DoneSkip (6), got %d", st2[0].Raw[6])
			}
		}
	}
}

func TestCKQServer_JoinDisjoinTeamSelectAndVoteKick_Behavior(t *testing.T) {
	ds, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	km := world.NewKQManager(ds, nil)
	sink := &recordingS2SSink{}
	km.SetS2SSink(sink, 0)

	// Pick a KQ with KQTeam.DivideType == 2 (UserSelect) or Vote enabled so we test TeamSelect & Vote end-to-end.
	// Inspect ds.KingdomQuest rows via HandleSchedule
	p1 := newCapturedClient(t, 701, "Starter")
	p2 := newCapturedClient(t, 702, "Voter")
	p3 := newCapturedClient(t, 703, "Target")
	km.RegisterOnline(p1.client, 701)
	km.RegisterOnline(p2.client, 702)
	km.RegisterOnline(p3.client, 703)

	km.HandleSchedule(p1.client)
	schedAck := findPacket(p1.drain(t), protocol.PROTO_NC_KQ_SCHEDULE_ACK)[0]
	cnt := int(binary.LittleEndian.Uint16(schedAck.Raw[10:12]))

	kqDataStore, err := kq.LoadDir("../../../data")
	if err != nil {
		t.Fatalf("kq.LoadDir: %v", err)
	}
	var kqHandle uint32
	var testLv byte = 20
	var maxPlayers uint16
	var teamKQH uint32
	var teamLv byte = 20
	for i := 0; i < cnt; i++ {
		off := 12 + i*141
		id := binary.LittleEndian.Uint16(schedAck.Raw[off+7 : off+9])
		if kqHandle == 0 && kqDataStore.GetIsVote(id) {
			kqHandle = binary.LittleEndian.Uint32(schedAck.Raw[off : off+4])
			testLv = schedAck.Raw[off+0x75]
			maxPlayers = binary.LittleEndian.Uint16(schedAck.Raw[off+0x79 : off+0x7B])
		}
		if teamKQH == 0 {
			if kt, ok := kqDataStore.GetKQTeamData(id); ok && kt.KQTeamDivideType == 2 {
				teamKQH = binary.LittleEndian.Uint32(schedAck.Raw[off : off+4])
				teamLv = schedAck.Raw[off+0x75]
			}
		}
	}
	if maxPlayers == 0 || kqHandle == 0 {
		t.Fatalf("failed to find vote-enabled test KQ handle")
	}

	// Activate Joining
	km.HandleZ2WMakeAck(kqHandle, 0x0981)
	_ = p1.drain(t)
	_ = p2.drain(t)
	_ = p3.drain(t)

	// Prison block -> 0x0999
	km.SetSessionCharacterStats(p1.client, testLv, 1, 1, 15)
	km.HandleJoin(p1.client, kqHandle, testLv, 1)
	if errCode := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_JOIN_ACK)[0].Raw[6:8]); errCode != 0x0999 {
		t.Fatalf("expected prison error 0x0999, got 0x%04X", errCode)
	}

	// Too low level -> 0x0995 (if minLv > 1)
	if testLv > 1 {
		km.SetSessionCharacterStats(p1.client, 1, 1, 1, 0)
		km.HandleJoin(p1.client, kqHandle, 1, 1)
		if errCode := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_JOIN_ACK)[0].Raw[6:8]); errCode != 0x0995 {
			t.Fatalf("expected level error 0x0995, got 0x%04X", errCode)
		}
	}

	// Valid joins for p1, p2, p3
	for _, cc := range []*capturedClient{p1, p2, p3} {
		km.SetSessionCharacterStats(cc.client, testLv, 1, 1, 0)
		km.HandleJoin(cc.client, kqHandle, testLv, 1)
		ack := findPacket(cc.drain(t), protocol.PROTO_NC_KQ_JOIN_ACK)[0]
		if errCode := binary.LittleEndian.Uint16(ack.Raw[6:8]); errCode != 0x0991 {
			t.Fatalf("expected join OK 0x0991 for %s, got 0x%04X", cc.client.CharName, errCode)
		}
	}

	// Already joined -> 0x099A
	km.HandleJoin(p1.client, kqHandle, testLv, 1)
	if errCode := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_JOIN_ACK)[0].Raw[6:8]); errCode != 0x099A {
		t.Fatalf("expected already-joined error 0x099A, got 0x%04X", errCode)
	}

	// JoinList OK -> Cooldown (5s from SingleData.shn)
	km.HandleJoinList(p1.client, kqHandle)
	jl1 := findPacket(p1.drain(t), protocol.PROTO_NC_KQ_JOIN_LIST_ACK)
	if len(jl1) == 0 || binary.LittleEndian.Uint16(jl1[len(jl1)-1].Raw[2:4]) != 0x3118 {
		t.Fatalf("expected JOIN_LIST_ACK OK 0x3118")
	}
	km.HandleJoinList(p1.client, kqHandle)
	jl2 := findPacket(p1.drain(t), protocol.PROTO_NC_KQ_JOIN_LIST_ACK)
	if len(jl2) != 1 || binary.LittleEndian.Uint16(jl2[0].Raw[2:4]) != 0x311A {
		t.Fatalf("expected JOIN_LIST_ACK cooldown 0x311A, got %v", jl2)
	}

	// Cannot vote while Joining -> 0x3102
	km.HandleVoteStartFull(p1.client, "Target", 1, "AFK")
	if errCode := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_START_ACK)[0].Raw[2:4]); errCode != 0x3102 {
		t.Fatalf("expected 0x3102 during Joining, got 0x%04X", errCode)
	}

	// Simulate MaxPlayers fill + 11s countdown -> Running (4)
	off := -1
	for i := 0; i < cnt; i++ {
		if binary.LittleEndian.Uint32(schedAck.Raw[12+i*141:12+i*141+4]) == kqHandle {
			off = 12 + i*141
			break
		}
	}
	startUnix := int64(int32(binary.LittleEndian.Uint32(schedAck.Raw[off+0x4B : off+0x4F])))
	startWait := binary.LittleEndian.Uint16(schedAck.Raw[off+0x73 : off+0x75])
	simNow := time.Unix(startUnix, 0).Add(time.Duration(startWait+1) * time.Second)
	km.SetClock(func() time.Time { return simNow })
	for sec := 0; sec <= 12; sec++ {
		simNow = simNow.Add(time.Second)
		km.DoKingdomQuest(simNow)
	}
	_ = p1.drain(t)
	_ = p2.drain(t)
	_ = p3.drain(t)

	// Confirm W2Z_START_CMD emitted
	evs := sink.popAll()
	var started bool
	for _, ev := range evs {
		if ev.Symbol == "wms_NC_KQ_W2Z_START_CMD" {
			started = true
		}
	}
	if !started {
		t.Fatalf("expected wms_NC_KQ_W2Z_START_CMD after countdown")
	}

	// VoteStartCheck -> MUST be 0x3120 on GetIsVote=true room
	km.HandleVoteStartCheck(p1.client)
	chkAck := findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_START_CHECK_ACK)[0]
	if chkErr := binary.LittleEndian.Uint16(chkAck.Raw[2:4]); chkErr != 0x3120 {
		t.Fatalf("expected 0x3120 on vote-enabled KQ, got 0x%04X", chkErr)
	}
	// Vote self -> 0x3106; empty reason -> 0x3107
	km.HandleVoteStartFull(p1.client, "Starter", 1, "AFK")
	if e := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_START_ACK)[0].Raw[2:4]); e != 0x3106 {
		t.Fatalf("expected self vote 0x3106, got 0x%04X", e)
	}
	km.HandleVoteStartFull(p1.client, "Target", 1, "")
	if e := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_START_ACK)[0].Raw[2:4]); e != 0x3107 {
		t.Fatalf("expected empty reason 0x3107, got 0x%04X", e)
	}

	// Valid VoteStart
	km.HandleVoteStartFull(p1.client, "Target", 1, "AFK")
	if e := binary.LittleEndian.Uint16(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_START_ACK)[0].Raw[2:4]); e != 0x3100 {
		t.Fatalf("expected vote start 0x3100, got 0x%04X", e)
	}
	if len(findPacket(p2.drain(t), protocol.PROTO_NC_KQ_VOTE_VOTING_CMD)) != 1 {
		t.Fatalf("expected voter p2 to receive PROTO_NC_KQ_VOTE_VOTING_CMD")
	}
	// p2 casts YES -> 0x3110; duplicate cast -> 0x3113
	km.HandleVoteCastType(p2.client, 1)
	if e := binary.LittleEndian.Uint16(findPacket(p2.drain(t), protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)[0].Raw[2:4]); e != 0x3110 {
		t.Fatalf("expected vote cast OK 0x3110, got 0x%04X", e)
	}
	km.HandleVoteCastType(p2.client, 1)
	if e := binary.LittleEndian.Uint16(findPacket(p2.drain(t), protocol.PROTO_NC_KQ_VOTE_VOTING_ACK)[0].Raw[2:4]); e != 0x3113 {
		t.Fatalf("expected duplicate vote error 0x3113, got 0x%04X", e)
	}

	// Advance simulated clock past KQVote_VoteLimitTime (60s) -> VoteProcessing emits SUC + BAN_MSG_CMD
	// and sets bBan=1 (0x004552A5); the next VoteProcessing tick (0x00454FF5) emits S2S BAN.
	simNow = simNow.Add(65 * time.Second)
	km.DoKingdomQuest(simNow)
	if len(findPacket(p1.drain(t), protocol.PROTO_NC_KQ_VOTE_RESULT_SUC_CMD)) != 1 {
		t.Fatalf("expected PROTO_NC_KQ_VOTE_RESULT_SUC_CMD")
	}
	if len(findPacket(p3.drain(t), protocol.PROTO_NC_KQ_VOTE_BAN_MSG_CMD)) != 1 {
		t.Fatalf("expected PROTO_NC_KQ_VOTE_BAN_MSG_CMD on target")
	}
	simNow = simNow.Add(1 * time.Second)
	km.DoKingdomQuest(simNow)
	evsBan := sink.popAll()
	var banSent bool
	for _, ev := range evsBan {
		if ev.Symbol == "wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD" {
			banSent = true
		}
	}
	if !banSent {
		t.Fatalf("expected wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD after passed vote")
	}

	// TeamSelect (0x5837 -> 0x5838 + 0x5839) if dataset contains DivideType==2
	if teamKQH != 0 {
		t1 := newCapturedClient(t, 801, "RedPlayer")
		t2 := newCapturedClient(t, 802, "BluePlayer")
		km.RegisterOnline(t1.client, 801)
		km.RegisterOnline(t2.client, 802)
		km.HandleZ2WMakeAck(teamKQH, 0x0981)
		km.SetSessionCharacterStats(t1.client, teamLv, 1, 1, 0)
		km.SetSessionCharacterStats(t2.client, teamLv, 1, 1, 0)
		km.HandleJoin(t1.client, teamKQH, teamLv, 1) // auto team 0
		km.HandleJoin(t2.client, teamKQH, teamLv, 1) // auto team 1
		_ = t1.drain(t)
		_ = t2.drain(t)
		// Same team -> 0x31F3
		km.HandleTeamSelect(t1.client, 0)
		if e := binary.LittleEndian.Uint16(findPacket(t1.drain(t), protocol.PROTO_NC_KQ_TEAM_SELECT_ACK)[0].Raw[2:4]); e != 0x31F3 {
			t.Fatalf("expected same team 0x31F3, got 0x%04X", e)
		}
		// Valid cancel join -> 0x09A1
		km.HandleJoinCancel(t2.client, teamKQH)
		if e := binary.LittleEndian.Uint16(findPacket(t2.drain(t), protocol.PROTO_NC_KQ_JOIN_CANCEL_ACK)[0].Raw[6:8]); e != 0x09A1 {
			t.Fatalf("expected join cancel 0x09A1, got 0x%04X", e)
		}
		km.HandleTeamSelect(t1.client, 1)
		_ = t1.drain(t)
	}
}

func TestKQProtocol_PacketStructures(t *testing.T) {
	timeAck := protocol.New_PROTO_NC_KQ_LIST_TIME_ACK(time.Now())
	if wire := timeAck.BuildWirePacket(); len(wire) != 43 {
		t.Errorf("expected 43 bytes for PROTO_NC_KQ_LIST_TIME_ACK, got %d", len(wire))
	}
	joinAck := protocol.New_PROTO_NC_KQ_JOIN_ACK(101, 0)
	if wire := joinAck.BuildWirePacket(); len(wire) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOIN_ACK, got %d", len(wire))
	}
	cancelAck := protocol.New_PROTO_NC_KQ_JOIN_CANCEL_ACK(101, 0)
	if wire := cancelAck.BuildWirePacket(); len(wire) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOIN_CANCEL_ACK, got %d", len(wire))
	}
	teamAck := protocol.New_PROTO_NC_KQ_TEAM_SELECT_ACK(0, 1)
	if wire := teamAck.BuildWirePacket(); len(wire) != 6 {
		t.Errorf("expected 6 bytes for PROTO_NC_KQ_TEAM_SELECT_ACK, got %d", len(wire))
	}
	teamCmd := protocol.New_PROTO_NC_KQ_TEAM_SELECT_CMD("TestChar", 2)
	if wire := teamCmd.BuildWirePacket(); len(wire) != 24 {
		t.Errorf("expected 24 bytes for PROTO_NC_KQ_TEAM_SELECT_CMD, got %d", len(wire))
	}
	alarmEnd := protocol.New_PROTO_NC_KQ_JOINING_ALARM_END_CMD(55, 2)
	if wire := alarmEnd.BuildWirePacket(); len(wire) != 9 {
		t.Errorf("expected 9 bytes for PROTO_NC_KQ_JOINING_ALARM_END_CMD, got %d", len(wire))
	}
	if wire := protocol.New_PROTO_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD().BuildWirePacket(); len(wire) != 4 {
		t.Errorf("expected 4 bytes for PROTO_NC_KQ_VOTE_BAN_MSG_LOGOFF_CMD, got %d", len(wire))
	}
	var outMaps [4][12]byte
	if wire := protocol.New_PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD(1, outMaps).BuildWirePacket(); len(wire) != 55 {
		t.Errorf("expected 55 bytes for PROTO_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD, got %d", len(wire))
	}
}

func TestCKQServer_KQMapSaveRestoreAndProtoCharKQMapCmd(t *testing.T) {
	// 1. Verify PackKQDateTime exact bitfields (Character.exe @ 0x0043911F..0x00439199)
	ts := time.Date(2026, 11, 20, 15, 42, 33, 0, time.UTC)
	packed := database.PackKQDateTime(ts)
	if (packed & 0x0F) != uint32(ts.Year()&0x0F) {
		t.Fatalf("year nibble mismatch: got %d want %d", packed&0x0F, ts.Year()&0x0F)
	}
	if ((packed >> 4) & 0x0F) != uint32(ts.Month()) {
		t.Fatalf("month nibble mismatch: got %d want %d", (packed>>4)&0x0F, ts.Month())
	}
	if ((packed >> 8) & 0x1F) != uint32(ts.Day()) {
		t.Fatalf("day bits mismatch: got %d want %d", (packed>>8)&0x1F, ts.Day())
	}
	if ((packed >> 13) & 0x1F) != uint32(ts.Hour()) {
		t.Fatalf("hour bits mismatch: got %d want %d", (packed>>13)&0x1F, ts.Hour())
	}
	if ((packed >> 18) & 0x3F) != uint32(ts.Minute()) {
		t.Fatalf("minute bits mismatch: got %d want %d", (packed>>18)&0x3F, ts.Minute())
	}
	if ((packed >> 24) & 0x3F) != uint32(ts.Second()) {
		t.Fatalf("second bits mismatch: got %d want %d", (packed>>24)&0x3F, ts.Second())
	}

	// 2. Verify PROTO_NC_CHAR_KQMAP_CMD (0x101A, 28 bytes payload -> 31 bytes wire)
	info := &database.CharKQMapInfo{
		KQHandle:  0x00001234,
		KQMapName: "KDEddyHill00",
		KQX:       2400,
		KQY:       3500,
		KQDate:    ts,
	}
	pw, data := world.BuildProtoCharKQMapCmd(info)
	wire := pw.BuildWirePacket()
	// 1-byte length prefix (30) + 2-byte opcode (0x101A) + 28-byte payload = 31 bytes
	if len(wire) != 31 {
		t.Fatalf("PROTO_NC_CHAR_KQMAP_CMD wire len = %d, want 31", len(wire))
	}
	if op := binary.LittleEndian.Uint16(wire[1:3]); op != protocol.PROTO_NC_CHAR_KQMAP_CMD {
		t.Fatalf("opcode = 0x%04X, want 0x%04X", op, protocol.PROTO_NC_CHAR_KQMAP_CMD)
	}
	if data.NKQHandle != 0x1234 || data.DKQDate != packed {
		t.Fatalf("PROTO_NC_CHAR_KQMAP_CMD data mismatch: %+v", data)
	}
	if uint32(data.NKQCoord&0xFFFFFFFF) != 2400 || uint32(data.NKQCoord>>32) != 3500 {
		t.Fatalf("NKQCoord mismatch: 0x%016X", data.NKQCoord)
	}

	// 3. Verify ResolveCharKQLoginLocation (p_Char_GetKQMap + fc_NC_CHAR_KQMAP_CMD)
	ds, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("LoadWorldDataStore: %v", err)
	}
	store, err := kq.LoadDir("../../../data")
	if err != nil {
		t.Fatalf("kq.LoadDir: %v", err)
	}
	km := world.NewKQManager(ds, nil)
	sink := &recordingS2SSink{}
	km.SetS2SSink(sink, 1)

	repo := database.NewMemoryCharacterRepositoryWithFile("")
	ch := &database.Character{
		AccountID: 99,
		Slot:      1,
		Name:      "RestoreHero",
		TownName:  "Rou",
		XPos:      5100,
		YPos:      6200,
		Level:     15,
	}
	if err := repo.CreateCharacter(ch); err != nil {
		t.Fatalf("CreateCharacter: %v", err)
	}

	// Start a KQ room and join RestoreHero using its exact StartTime/StartWait from HandleSchedule
	c1 := newCapturedClient(t, ch.ID, ch.Name)
	km.RegisterOnline(c1.client, ch.ID)
	km.HandleSchedule(c1.client)
	schedAck := findPacket(c1.drain(t), protocol.PROTO_NC_KQ_SCHEDULE_ACK)[0]
	cnt := int(binary.LittleEndian.Uint16(schedAck.Raw[10:12]))

	var kqHandle uint32
	var minLv byte = 10
	var startUnix int64
	var startWait uint16
	for i := 0; i < cnt; i++ {
		off := 12 + i*141
		id := binary.LittleEndian.Uint16(schedAck.Raw[off+7 : off+9])
		for _, q := range store.Quests {
			if q.ID == id && q.MinPlayers <= 1 {
				kqHandle = binary.LittleEndian.Uint32(schedAck.Raw[off : off+4])
				minLv = q.MinLevel
				startUnix = int64(int32(binary.LittleEndian.Uint32(schedAck.Raw[off+0x4B : off+0x4F])))
				startWait = binary.LittleEndian.Uint16(schedAck.Raw[off+0x73 : off+0x75])
				break
			}
		}
		if kqHandle != 0 {
			break
		}
	}
	if kqHandle == 0 {
		t.Fatalf("no KQ room scheduled with MinPlayers <= 1")
	}
	km.HandleZ2WMakeAck(kqHandle, 0x0981)
	_ = c1.drain(t)
	km.SetSessionCharacterStats(c1.client, minLv, 1, 1, 0)
	km.HandleJoin(c1.client, kqHandle, minLv, 1)
	_ = c1.drain(t)

	// Advance clock to StartTime + StartWait + countdown -> KQStatusRunning
	simNow := time.Unix(startUnix, 0).Add(time.Duration(startWait+1) * time.Second)
	km.SetClock(func() time.Time { return simNow })
	for sec := 0; sec <= 12; sec++ {
		simNow = simNow.Add(time.Second)
		km.DoKingdomQuest(simNow)
	}
	if !km.InKQStatusRunning(kqHandle) {
		t.Fatalf("expected KQ handle %d to be KQStatusRunning", kqHandle)
	}

	// Persist KQ position via p_Char_SaveLocation (SaveCharLocationWithKQ)
	if err := repo.SaveCharLocationWithKQ(ch.ID, "Rou", 5100, 6200, kqHandle, "KDEddyHill00", 2400, 3500); err != nil {
		t.Fatalf("SaveCharLocationWithKQ: %v", err)
	}

	// Case A: KQ is RUNNING and player is joiner -> restores sKQMap (KDEddyHill00, 2400, 3500)
	destMap, dx, dy, kqPkt, restored := km.ResolveCharKQLoginLocation(repo, ch.ID, ch.Name, "Rou", 5100, 6200)
	if !restored || destMap != "KDEddyHill00" || dx != 2400 || dy != 3500 || kqPkt == nil {
		t.Fatalf("expected KQ restore (KDEddyHill00, 2400, 3500), got (%s, %d, %d, restored=%v)", destMap, dx, dy, restored)
	}

	// Case B: KQ ends via Z2W_END_CMD -> restores pre-KQ position (Rou, 5100, 6200)
	km.HandleZ2WEndCmd(kqHandle)
	destMap2, dx2, dy2, _, restored2 := km.ResolveCharKQLoginLocation(repo, ch.ID, ch.Name, "Rou", 5100, 6200)
	if restored2 || destMap2 != "Rou" || dx2 != 5100 || dy2 != 6200 {
		t.Fatalf("expected fallback to pre-KQ (Rou, 5100, 6200) after KQ end, got (%s, %d, %d, restored=%v)", destMap2, dx2, dy2, restored2)
	}
}
