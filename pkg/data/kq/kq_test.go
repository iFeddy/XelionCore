package kq

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func dataDir(t *testing.T) string {
	t.Helper()
	p := filepath.Join("..", "..", "..", "data")
	if _, err := os.Stat(filepath.Join(p, "KingdomQuest.shn")); err != nil {
		t.Fatalf("cannot locate data directory at %s: %v", p, err)
	}
	return p
}

// decryptSHNRawRows decrypts an SHN file independently and returns each row's raw payload (after the 2-byte row_len).
func decryptSHNRawRows(t *testing.T, path string) (recordCount int, rowSize int, payloads [][]byte) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	if len(raw) < 52 {
		t.Fatalf("%s too short: %d", path, len(raw))
	}
	dataMode := binary.LittleEndian.Uint32(raw[28:32])
	full := raw
	if dataMode == 1 {
		dec := shn.Decrypt(raw[36:])
		full = make([]byte, 36+len(dec))
		copy(full[:36], raw[:36])
		copy(full[36:], dec)
	}
	recordCount = int(binary.LittleEndian.Uint32(full[40:44]))
	rowSize = int(binary.LittleEndian.Uint32(full[44:48]))
	colCount := int(binary.LittleEndian.Uint32(full[48:52]))
	off := 52 + colCount*56
	payloads = make([][]byte, recordCount)
	for i := 0; i < recordCount; i++ {
		row := full[off : off+rowSize]
		recLen := int(binary.LittleEndian.Uint16(row[0:2]))
		if recLen != rowSize {
			t.Fatalf("%s row %d: row_len=%d != rowSize=%d", path, i, recLen, rowSize)
		}
		p := make([]byte, rowSize-2)
		copy(p, row[2:])
		payloads[i] = p
		off += rowSize
	}
	return recordCount, rowSize, payloads
}

type binaryRecord interface {
	MarshalBinary() ([]byte, error)
	Size() int
}

func assertSHNBinaryParity[T binaryRecord](t *testing.T, path string, wantCount, wantSize int, items []T, unmarshal func([]byte) (T, error)) {
	t.Helper()
	recCount, rowSize, rawPayloads := decryptSHNRawRows(t, path)
	if recCount != wantCount {
		t.Fatalf("%s: recordCount=%d, want %d", path, recCount, wantCount)
	}
	if rowSize-2 != wantSize {
		t.Fatalf("%s: payloadSize=%d, want %d", path, rowSize-2, wantSize)
	}
	if len(items) != wantCount {
		t.Fatalf("%s: loaded %d rows, want %d", path, len(items), wantCount)
	}
	for i, item := range items {
		if item.Size() != wantSize {
			t.Fatalf("%s row %d: Size()=%d, want %d", path, i, item.Size(), wantSize)
		}
		got, err := item.MarshalBinary()
		if err != nil {
			t.Fatalf("%s row %d: MarshalBinary error: %v", path, i, err)
		}
		if len(got) != wantSize {
			t.Fatalf("%s row %d: len(MarshalBinary)=%d, want %d", path, i, len(got), wantSize)
		}
		if !bytes.Equal(got, rawPayloads[i]) {
			t.Fatalf("%s row %d: binary mismatch against decrypted SHN\ngot:  %x\nwant: %x", path, i, got, rawPayloads[i])
		}
		rt, err := unmarshal(rawPayloads[i])
		if err != nil {
			t.Fatalf("%s row %d: UnmarshalBinary error: %v", path, i, err)
		}
		if !reflect.DeepEqual(rt, item) {
			t.Fatalf("%s row %d: round-trip mismatch\ngot:  %+v\nwant: %+v", path, i, rt, item)
		}
	}
}

func TestLoadAll7SHNs_RowByRowBinaryParity(t *testing.T) {
	dir := dataDir(t)
	store, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir(%s): %v", dir, err)
	}

	assertSHNBinaryParity(t, filepath.Join(dir, "KingdomQuest.shn"), 57, SizeOfKingdomQuest, store.Quests, func(b []byte) (KINGDOM_QUEST, error) {
		var v KINGDOM_QUEST
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KingdomQuestMap.shn"), 38, SizeOfKingdomQuestMap, store.Maps, func(b []byte) (KINGDOM_QUEST_MAP, error) {
		var v KINGDOM_QUEST_MAP
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KQTeam.shn"), 8, SizeOfKQTeam, store.Teams, func(b []byte) (KQTeam, error) {
		var v KQTeam
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KQIsVote.shn"), 30, SizeOfKQIsVote, store.IsVotes, func(b []byte) (KQIsVote, error) {
		var v KQIsVote
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KQVoteMajorityRate.shn"), 2, SizeOfKQVoteMajorityRate, store.VoteMajorityRates, func(b []byte) (KQVoteMajorityRate, error) {
		var v KQVoteMajorityRate
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KingdomQuestRew.shn"), 64, SizeOfKingdomQuestRew, store.Rewards, func(b []byte) (KINGDOM_QUEST_REW, error) {
		var v KINGDOM_QUEST_REW
		err := v.UnmarshalBinary(b)
		return v, err
	})
	assertSHNBinaryParity(t, filepath.Join(dir, "KQItem.shn"), 2, SizeOfKQItem, store.Items, func(b []byte) (KQItem, error) {
		var v KQItem
		err := v.UnmarshalBinary(b)
		return v, err
	})
}

func TestPositionalRepeatedColumns(t *testing.T) {
	dir := dataDir(t)
	store, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	// KingdomQuest.shn row 0: MapLink = [0, 0xFFFF, 0xFFFF, 0xFFFF], DemandGender = [1, 1]
	q0 := store.Quests[0]
	if q0.MapLink != [4]uint16{0, 0xFFFF, 0xFFFF, 0xFFFF} {
		t.Fatalf("row 0 MapLink = %v, want [0 65535 65535 65535]", q0.MapLink)
	}
	if q0.DemandGender != [2]uint8{1, 1} {
		t.Fatalf("row 0 DemandGender = %v, want [1 1]", q0.DemandGender)
	}

	// KingdomQuestMap.shn row 0: KDEddyHill with Map[0..9] = KDEddyHill00..09 and Clear[0..9] = 1
	m0 := store.Maps[0]
	if m0.NumOfMap != 10 || m0.BaseMapString() != "KDEddyHill" {
		t.Fatalf("map 0 header unexpected: num=%d base=%q", m0.NumOfMap, m0.BaseMapString())
	}
	for i := 0; i < 10; i++ {
		wantName := "KDEddyHill0" + string(rune('0'+i))
		if m0.MapString(i) != wantName {
			t.Fatalf("map 0 Map[%d] = %q, want %q", i, m0.MapString(i), wantName)
		}
		if m0.Clear[i] != 1 {
			t.Fatalf("map 0 Clear[%d] = %d, want 1", i, m0.Clear[i])
		}
	}

	// KingdomQuestRew.shn row 0: Reward[0..7] = 1..8, RewardRate[0..7] = [150, 150, 100, 500, 300, 100, 1000, 1000]
	r0 := store.Rewards[0]
	wantRew := [15]uint16{1, 2, 3, 4, 5, 6, 7, 8}
	wantRate := [15]uint16{150, 150, 100, 500, 300, 100, 1000, 1000}
	if r0.Reward != wantRew || r0.RewardRate != wantRate {
		t.Fatalf("reward 0 mismatch: rew=%v rate=%v", r0.Reward, r0.RewardRate)
	}

	// Check Store lookups matching CKQServer::GetIsVote (0x00451DA0) and GetKQTeamData (0x00451E50)
	if !store.GetIsVote(0) || store.GetIsVote(17) || store.GetIsVote(9999) {
		t.Fatalf("GetIsVote unexpected results: 0=%v 17=%v 9999=%v", store.GetIsVote(0), store.GetIsVote(17), store.GetIsVote(9999))
	}
	if tm, ok := store.GetKQTeamData(26); !ok || tm.KQTeamDivideType != uint16(KQTeamDivideRandom) {
		t.Fatalf("GetKQTeamData(26) = (%+v, %v)", tm, ok)
	}
	if _, ok := store.GetKQTeamData(0); ok {
		t.Fatalf("GetKQTeamData(0) should be false (not in KQTeam.shn)")
	}
	if rate, ok := store.GetVoteAgreeRate(0); !ok || rate != 70 {
		t.Fatalf("GetVoteAgreeRate(0) = (%d, %v), want (70, true)", rate, ok)
	}
}

func TestProtoKQInfoAndMapInfo_SizeofAndOffsets(t *testing.T) {
	var mapInfo PROTO_KQ_MAP_INFO
	if mapInfo.Size() != 26 {
		t.Fatalf("PROTO_KQ_MAP_INFO.Size() = %d, want 26", mapInfo.Size())
	}
	mapInfo.MapIndex = 3
	mapInfo.SetMapBase("KDEddyHill")
	mapInfo.SetMapName("KDEddyHill03")
	mapInfo.MapClear = 1

	mb, err := mapInfo.MarshalBinary()
	if err != nil || len(mb) != 26 {
		t.Fatalf("PROTO_KQ_MAP_INFO MarshalBinary len=%d err=%v", len(mb), err)
	}
	if mb[0x00] != 3 || trimCString(mb[0x01:0x0d]) != "KDEddyHill" || trimCString(mb[0x0d:0x19]) != "KDEddyHill03" || mb[0x19] != 1 {
		t.Fatalf("PROTO_KQ_MAP_INFO field offsets mismatch: %x", mb)
	}
	var mapRT PROTO_KQ_MAP_INFO
	if err := mapRT.UnmarshalBinary(mb); err != nil || mapRT != mapInfo {
		t.Fatalf("PROTO_KQ_MAP_INFO round-trip failed: %+v err=%v", mapRT, err)
	}

	var clientInfo PROTO_KQ_INFO_CLIENT
	if clientInfo.Size() != 141 {
		t.Fatalf("PROTO_KQ_INFO_CLIENT.Size() = %d, want 141", clientInfo.Size())
	}

	startT := time.Date(2026, 11, 16, 17, 30, 0, 0, time.Local)
	schedT := time.Date(2026, 11, 16, 17, 29, 0, 0, time.Local)

	var info PROTO_KQ_INFO
	if info.Size() != 377 {
		t.Fatalf("PROTO_KQ_INFO.Size() = %d, want 377", info.Size())
	}
	info.Handle = 0x11223344
	info.Status = KQStatusJoining
	info.NumOfJoiner = 7
	info.ID = 5
	info.SetTitle("Lost Mini Dragon[A]")
	info.LimitTime = 50
	info.StartTime = int32(startT.Unix())
	info.TmStartTime = CTimeTMFromTime(startT)
	info.StartWaitTime = 30
	info.MinLevel = 46
	info.MaxLevel = 60
	info.MinPlayers = 1
	info.MaxPlayers = 15
	info.PlayerRepeatMode = 0
	info.PlayerRepeatCount = 2
	info.PlayerRevivalMode = 4
	info.PlayerRevivalCount = 3
	info.DemandQuest = 101
	info.DemandItem = 202
	info.DemandClass = 0x0102030405060708
	info.DemandGender = 3
	info.NextStartMode = 1
	info.NextStartDelayMin = 120
	info.RepeatMode = 0
	info.RepeatCount = 3
	info.RewardIndex = 5
	info.DemandMobKill = 5
	info.ScheduleTime = int32(schedT.Unix())
	info.TmScheduleTime = CTimeTMFromTime(schedT)
	info.RunCounter = 10
	info.MapLink[0] = mapInfo
	info.SetScriptLanguage("KQ/MiniDragon/MiniDragon")
	info.SetScriptInitValue("10")
	info.IsTeamPVP = 1
	info.TeamRegenXY[0] = SHINE_XY_TYPE{X: 111, Y: 222}
	info.TeamRegenXY[1] = SHINE_XY_TYPE{X: 333, Y: 444}

	raw, err := info.MarshalBinary()
	if err != nil || len(raw) != 377 {
		t.Fatalf("PROTO_KQ_INFO MarshalBinary len=%d err=%v", len(raw), err)
	}

	// Verify exact PDB byte offsets
	if binary.LittleEndian.Uint32(raw[0x00:0x04]) != 0x11223344 {
		t.Fatalf("offset 0x00 Handle mismatch")
	}
	if raw[0x04] != byte(KQStatusJoining) {
		t.Fatalf("offset 0x04 Status mismatch: %d", raw[0x04])
	}
	if binary.LittleEndian.Uint16(raw[0x05:0x07]) != 7 || binary.LittleEndian.Uint16(raw[0x07:0x09]) != 5 {
		t.Fatalf("offset 0x05/0x07 mismatch")
	}
	if raw[0x8d] != 1 || binary.LittleEndian.Uint16(raw[0x8e:0x90]) != 120 {
		t.Fatalf("offset 0x8d/0x8e mismatch")
	}
	if binary.LittleEndian.Uint32(raw[0x97:0x9b]) != uint32(schedT.Unix()) {
		t.Fatalf("offset 0x97 ScheduleTime mismatch")
	}
	if raw[0xbf] != 10 {
		t.Fatalf("offset 0xbf RunCounter mismatch")
	}
	if !bytes.Equal(raw[0xc0:0xc0+26], mb) {
		t.Fatalf("offset 0xc0 MapLink[0] mismatch")
	}
	if trimCString(raw[0x128:0x148]) != "KQ/MiniDragon/MiniDragon" || trimCString(raw[0x148:0x168]) != "10" {
		t.Fatalf("offset 0x128/0x148 script fields mismatch")
	}
	if raw[0x168] != 1 || binary.LittleEndian.Uint32(raw[0x169:0x16d]) != 111 || binary.LittleEndian.Uint32(raw[0x175:0x179]) != 444 {
		t.Fatalf("offset 0x168/0x169 team pvp/regen mismatch")
	}

	var infoRT PROTO_KQ_INFO
	if err := infoRT.UnmarshalBinary(raw); err != nil {
		t.Fatalf("PROTO_KQ_INFO UnmarshalBinary: %v", err)
	}
	if !reflect.DeepEqual(infoRT, info) {
		t.Fatalf("PROTO_KQ_INFO round-trip mismatch:\ngot:  %+v\nwant: %+v", infoRT, info)
	}
	if !infoRT.TmStartTime.ToTime(time.Local).Equal(startT) {
		t.Fatalf("TmStartTime round-trip mismatch: got %v want %v", infoRT.TmStartTime.ToTime(time.Local), startT)
	}
}

func TestLoadShineReward_BinaryParityAndLookups(t *testing.T) {
	dir := dataDir(t)
	store, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir(%s) failed: %v", dir, err)
	}
	if len(store.ShineRewards) != 435 {
		t.Fatalf("ShineReward.shn row count = %d, want 435", len(store.ShineRewards))
	}

	// Verify 66-byte binary round-trip on every row of ShineReward.shn (Zone.exe @ 0x0064C830)
	for i, row := range store.ShineRewards {
		raw, err := row.MarshalBinary()
		if err != nil {
			t.Fatalf("ShineReward[%d] MarshalBinary: %v", i, err)
		}
		if len(raw) != SizeOfShineReward {
			t.Fatalf("ShineReward[%d] size = %d, want %d", i, len(raw), SizeOfShineReward)
		}
		var rt SHINE_REWARD
		if err := rt.UnmarshalBinary(raw); err != nil {
			t.Fatalf("ShineReward[%d] UnmarshalBinary: %v", i, err)
		}
		if rt != row {
			t.Fatalf("ShineReward[%d] round-trip mismatch: got %+v want %+v", i, rt, row)
		}
	}

	// Verify KQRewardDataBox::operator[](uint8*) (exe: 0x0049A9A0) and RewardData::rd_FindHandle (exe: 0x0064C850)
	var namedRow *KINGDOM_QUEST_REW
	for i := range store.Rewards {
		if store.Rewards[i].IndexName() != "" {
			namedRow = &store.Rewards[i]
			break
		}
	}
	if namedRow == nil {
		t.Fatalf("expected at least one KingdomQuestRew row with non-empty IndexName")
	}
	rewNamed, ok := store.GetRewardByName(namedRow.IndexName())
	if !ok || rewNamed.ID != namedRow.ID {
		t.Fatalf("GetRewardByName(%q) = (%+v, %v), want ID=%d", namedRow.IndexName(), rewNamed, ok, namedRow.ID)
	}
	if _, ok := store.GetRewardByName("NonExistentKQReward"); ok {
		t.Fatalf("GetRewardByName should return false for unknown index")
	}

	// Every non-zero Reward slot in KingdomQuestRew.shn row 0 must resolve in ShineReward.shn
	rew0 := store.Rewards[0]
	resolved := 0
	for slot := 0; slot < 15; slot++ {
		srID := rew0.Reward[slot]
		if srID == 0 {
			continue
		}
		sr, ok := store.GetShineReward(srID)
		if !ok {
			t.Fatalf("KingdomQuestRew[0].Reward[%d] = %d not found in ShineReward.shn", slot, srID)
		}
		if sr.RewardHandle != srID {
			t.Fatalf("GetShineReward(%d).RewardHandle = %d", srID, sr.RewardHandle)
		}
		resolved++
	}
	if resolved == 0 {
		t.Fatalf("expected at least one resolved ShineReward from KingdomQuestRew[0]")
	}
}


