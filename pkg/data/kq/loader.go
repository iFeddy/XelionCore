package kq

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// Store holds all official Kingdom Quest SHN tables (7 KQ SHNs + ShineReward.shn)
// loaded with positional column access.
type Store struct {
	Quests            []KINGDOM_QUEST
	Maps              []KINGDOM_QUEST_MAP
	Teams             []KQTeam
	IsVotes           []KQIsVote
	VoteMajorityRates []KQVoteMajorityRate
	Rewards           []KINGDOM_QUEST_REW
	Items             []KQItem
	ShineRewards      []SHINE_REWARD
}

// LoadDir loads all Kingdom Quest SHN files (+ ShineReward.shn if present) from dataDir.
func LoadDir(dataDir string) (*Store, error) {
	quests, err := LoadKingdomQuest(filepath.Join(dataDir, "KingdomQuest.shn"))
	if err != nil {
		return nil, err
	}
	maps, err := LoadKingdomQuestMap(filepath.Join(dataDir, "KingdomQuestMap.shn"))
	if err != nil {
		return nil, err
	}
	teams, err := LoadKQTeam(filepath.Join(dataDir, "KQTeam.shn"))
	if err != nil {
		return nil, err
	}
	isVotes, err := LoadKQIsVote(filepath.Join(dataDir, "KQIsVote.shn"))
	if err != nil {
		return nil, err
	}
	voteRates, err := LoadKQVoteMajorityRate(filepath.Join(dataDir, "KQVoteMajorityRate.shn"))
	if err != nil {
		return nil, err
	}
	rewards, err := LoadKingdomQuestRew(filepath.Join(dataDir, "KingdomQuestRew.shn"))
	if err != nil {
		return nil, err
	}
	items, err := LoadKQItem(filepath.Join(dataDir, "KQItem.shn"))
	if err != nil {
		return nil, err
	}

	var shineRewards []SHINE_REWARD
	shineRewPath := filepath.Join(dataDir, "ShineReward.shn")
	if _, statErr := os.Stat(shineRewPath); statErr == nil {
		shineRewards, err = LoadShineReward(shineRewPath)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		Quests:            quests,
		Maps:              maps,
		Teams:             teams,
		IsVotes:           isVotes,
		VoteMajorityRates: voteRates,
		Rewards:           rewards,
		Items:             items,
		ShineRewards:      shineRewards,
	}, nil
}

// LoadKingdomQuest loads KingdomQuest.shn using positional column indices (35 columns, 180-byte record).
// Matches CKQServer::Read (WorldManager.exe @ 0x00451B40, DataVersionKey=0x2535, RowSize=0xb4).
func LoadKingdomQuest(path string) ([]KINGDOM_QUEST, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 35, SizeOfKingdomQuest); err != nil {
		return nil, err
	}
	out := make([]KINGDOM_QUEST, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		q := KINGDOM_QUEST{
			ID:                 r.GetUint16At(0),
			LimitTime:          r.GetUint16At(2),
			STYear:             r.GetByteAt(3),
			STMonth:            r.GetByteAt(4),
			STDay:              r.GetByteAt(5),
			STHour:             r.GetByteAt(6),
			STMinute:           r.GetByteAt(7),
			STSecond:           r.GetByteAt(8),
			StartWaitTime:      r.GetUint16At(9),
			NextStartMode:      r.GetByteAt(10),
			NextStartDeleyMin:  r.GetUint16At(11),
			RepeatMode:         r.GetByteAt(12),
			RepeatCount:        r.GetUint16At(13),
			MinLevel:           r.GetByteAt(14),
			MaxLevel:           r.GetByteAt(15),
			MinPlayers:         r.GetUint16At(16),
			MaxPlayers:         r.GetUint16At(17),
			PlayerRepeatMode:   r.GetByteAt(18),
			PlayerRepeatCount:  r.GetUint16At(19),
			PlayerRevivalMode:  r.GetByteAt(20),
			PlayerRevivalCount: r.GetByteAt(21),
			DemandQuest:        r.GetUint16At(22),
			DemandItem:         r.GetUint16At(23),
			DemandMobKill:      r.GetByteAt(24),
			RewardIndex:        r.GetUint32At(25),
			UseClass:           r.GetUint32At(32),
		}
		copy(q.Title[:], r.GetBytesAt(1))
		for m := 0; m < 4; m++ {
			q.MapLink[m] = r.GetUint16At(26 + m)
		}
		copy(q.ScriptLanguage[:], r.GetBytesAt(30))
		copy(q.InitValue[:], r.GetBytesAt(31))
		q.DemandGender[0] = r.GetByteAt(33)
		q.DemandGender[1] = r.GetByteAt(34)
		out[i] = q
	}
	return out, nil
}

// LoadKingdomQuestMap loads KingdomQuestMap.shn using positional column indices (22 columns, 363-byte record).
// Matches CKQServer::ReadMap (WorldManager.exe @ 0x00453E70, DataVersionKey=0x2085, RowSize=0x16b).
func LoadKingdomQuestMap(path string) ([]KINGDOM_QUEST_MAP, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 22, SizeOfKingdomQuestMap); err != nil {
		return nil, err
	}
	out := make([]KINGDOM_QUEST_MAP, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		m := KINGDOM_QUEST_MAP{
			NumOfMap: r.GetByteAt(0),
		}
		copy(m.BaseMap[:], r.GetBytesAt(1))
		for j := 0; j < 10; j++ {
			copy(m.Map[j][:], r.GetBytesAt(2+j))
			m.Clear[j] = r.GetByteAt(12 + j)
		}
		out[i] = m
	}
	return out, nil
}

// LoadKQTeam loads KQTeam.shn using positional column indices (8 columns, 22-byte record).
// Matches CKQServer::ReadTeam (WorldManager.exe @ 0x00451D80, DataVersionKey=0xf1, RowSize=0x16).
func LoadKQTeam(path string) ([]KQTeam, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 8, SizeOfKQTeam); err != nil {
		return nil, err
	}
	out := make([]KQTeam, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		out[i] = KQTeam{
			ID:               r.GetUint16At(0),
			MaxMemberGap:     r.GetByteAt(1),
			IsTeamPVP:        r.GetByteAt(2),
			KQTeamDivideType: r.GetUint16At(3),
			RegenXRed:        r.GetUint32At(4),
			RegenYRed:        r.GetUint32At(5),
			RegenXBlue:       r.GetUint32At(6),
			RegenYBlue:       r.GetUint32At(7),
		}
	}
	return out, nil
}

// LoadKQIsVote loads KQIsVote.shn using positional column indices (2 columns, 3-byte record).
// Matches CKQServer::ReadIsVote (WorldManager.exe @ 0x00451B80, DataVersionKey=0x5, RowSize=0x3).
func LoadKQIsVote(path string) ([]KQIsVote, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 2, SizeOfKQIsVote); err != nil {
		return nil, err
	}
	out := make([]KQIsVote, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		out[i] = KQIsVote{
			ID:     r.GetUint16At(0),
			IsVote: r.GetByteAt(1),
		}
	}
	return out, nil
}

// LoadKQVoteMajorityRate loads KQVoteMajorityRate.shn using positional column indices (1 column, 1-byte record).
// Matches CKQServer::ReadVoteRate (WorldManager.exe @ 0x00451B60, DataVersionKey=0x1, RowSize=0x1).
func LoadKQVoteMajorityRate(path string) ([]KQVoteMajorityRate, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 1, SizeOfKQVoteMajorityRate); err != nil {
		return nil, err
	}
	out := make([]KQVoteMajorityRate, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		out[i] = KQVoteMajorityRate{
			VoteAgreeRate: r.GetByteAt(0),
		}
	}
	return out, nil
}

// LoadKingdomQuestRew loads KingdomQuestRew.shn using positional column indices (33 columns, 128-byte record).
// Matches KQRewardDataBox::kqrdb_Load (Zone.exe @ 0x0049A940, DataVersionKey=0x32a9, RowSize=0x80).
func LoadKingdomQuestRew(path string) ([]KINGDOM_QUEST_REW, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 33, SizeOfKingdomQuestRew); err != nil {
		return nil, err
	}
	out := make([]KINGDOM_QUEST_REW, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		rew := KINGDOM_QUEST_REW{
			ID: r.GetUint32At(0),
		}
		copy(rew.IndexString[:], r.GetBytesAt(1))
		copy(rew.KQBoxItemIDX[:], r.GetBytesAt(2))
		for j := 0; j < 15; j++ {
			rew.Reward[j] = r.GetUint16At(3 + j)
			rew.RewardRate[j] = r.GetUint16At(18 + j)
		}
		out[i] = rew
	}
	return out, nil
}

// LoadKQItem loads KQItem.shn using positional column indices (4 columns, 38-byte record).
// Matches KQItemDataBox::kqidb_Load (Zone.exe @ 0x0064BD80, DataVersionKey=0x22, RowSize=0x26).
func LoadKQItem(path string) ([]KQItem, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 4, SizeOfKQItem); err != nil {
		return nil, err
	}
	out := make([]KQItem, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		item := KQItem{
			MoveSpdRate:    r.GetUint16At(1),
			AbsoluteAttack: r.GetUint16At(2),
			PickupLimit:    r.GetUint16At(3),
		}
		copy(item.ItemIndex[:], r.GetBytesAt(0))
		out[i] = item
	}
	return out, nil
}

func validateSHNLayout(path string, f *shn.File, wantCols, wantPayloadSize int) error {
	if len(f.Columns) != wantCols {
		return fmt.Errorf("%s: column count mismatch: got %d, want %d", path, len(f.Columns), wantCols)
	}
	if f.RowSize-2 != wantPayloadSize {
		return fmt.Errorf("%s: payload size mismatch: got %d, want %d", path, f.RowSize-2, wantPayloadSize)
	}
	return nil
}

// GetIsVote mirrors CKQServer::GetIsVote (WorldManager.exe @ 0x00451DA0):
// returns true iff KQ id exists in Quests and KQIsVote has IsVote == 1 for that ID.
func (s *Store) GetIsVote(kqID uint16) bool {
	if s == nil {
		return false
	}
	hasQuest := false
	for i := range s.Quests {
		if s.Quests[i].ID == kqID {
			hasQuest = true
			break
		}
	}
	if !hasQuest {
		return false
	}
	for i := range s.IsVotes {
		if s.IsVotes[i].ID == kqID {
			return s.IsVotes[i].IsVote == 1
		}
	}
	return false
}

// GetKQTeamData mirrors CKQServer::GetKQTeamData (WorldManager.exe @ 0x00451E50):
// returns the KQTeam entry matching kqID if kqID exists in Quests.
func (s *Store) GetKQTeamData(kqID uint16) (*KQTeam, bool) {
	if s == nil {
		return nil, false
	}
	hasQuest := false
	for i := range s.Quests {
		if s.Quests[i].ID == kqID {
			hasQuest = true
			break
		}
	}
	if !hasQuest {
		return nil, false
	}
	for i := range s.Teams {
		if s.Teams[i].ID == kqID {
			t := s.Teams[i]
			return &t, true
		}
	}
	return nil, false
}

// GetMap returns the KINGDOM_QUEST_MAP at row index mapIdx.
func (s *Store) GetMap(mapIdx uint16) (*KINGDOM_QUEST_MAP, bool) {
	if s == nil || int(mapIdx) >= len(s.Maps) {
		return nil, false
	}
	m := s.Maps[mapIdx]
	return &m, true
}

// LoadShineReward loads ShineReward.shn using positional column indices (16 columns, 66-byte record).
// Matches RewardData::rd_Load (Zone.exe @ 0x0064C830, "../9Data/Shine/ShineReward.shn").
func LoadShineReward(path string) ([]SHINE_REWARD, error) {
	f, err := shn.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateSHNLayout(path, f, 16, SizeOfShineReward); err != nil {
		return nil, err
	}
	out := make([]SHINE_REWARD, len(f.Rows))
	for i := range f.Rows {
		r := &f.Rows[i]
		sr := SHINE_REWARD{
			RewardHandle: r.GetUint16At(0),
			RewardType:   ShineRewardType(r.GetByteAt(1)),
			Quantity:     r.GetUint32At(3),
			OptionDegree: r.GetUint16At(14),
			TitleDegree:  r.GetUint32At(15),
		}
		copy(sr.Argument[:], r.GetBytesAt(2))
		for j := 0; j < 10; j++ {
			sr.Upgrade[j] = r.GetUint16At(4 + j)
		}
		out[i] = sr
	}
	return out, nil
}

// GetReward returns the KINGDOM_QUEST_REW matching rewardID.
// Mirrors KQRewardDataBox::operator[](uint16) (Zone.exe @ 0x0049A950).
func (s *Store) GetReward(rewardID uint32) (*KINGDOM_QUEST_REW, bool) {
	if s == nil {
		return nil, false
	}
	for i := range s.Rewards {
		if s.Rewards[i].ID == rewardID {
			r := s.Rewards[i]
			return &r, true
		}
	}
	return nil, false
}

// GetRewardByName returns the KINGDOM_QUEST_REW matching IndexString exactly (case-sensitive).
// Mirrors KQRewardDataBox::operator[](uint8*) (Zone.exe @ 0x0049A9A0).
func (s *Store) GetRewardByName(inxName string) (*KINGDOM_QUEST_REW, bool) {
	if s == nil || inxName == "" {
		return nil, false
	}
	for i := range s.Rewards {
		if s.Rewards[i].IndexName() == inxName {
			r := s.Rewards[i]
			return &r, true
		}
	}
	return nil, false
}

// GetShineReward returns the SHINE_REWARD matching RewardHandle.
// Mirrors RewardData::rd_FindHandle (Zone.exe @ 0x0064C850).
func (s *Store) GetShineReward(handle uint16) (*SHINE_REWARD, bool) {
	if s == nil {
		return nil, false
	}
	for i := range s.ShineRewards {
		if s.ShineRewards[i].RewardHandle == handle {
			sr := s.ShineRewards[i]
			return &sr, true
		}
	}
	return nil, false
}

// GetVoteAgreeRate returns the vote majority rate at row index idx (0 = normal, 1 = tie/alt).
func (s *Store) GetVoteAgreeRate(idx int) (uint8, bool) {
	if s == nil || idx < 0 || idx >= len(s.VoteMajorityRates) {
		return 0, false
	}
	return s.VoteMajorityRates[idx].VoteAgreeRate, true
}

