package database

import (
	"fmt"
	"sort"
	"time"
)

func (r *MemoryCharacterRepository) SetWeaponTitle(weaponItemKey, licenseItemKey int64, mobID int16, titleCharNo int, title string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.weaponTitles[weaponItemKey] == nil {
		r.weaponTitles[weaponItemKey] = make(map[int16]string)
	}
	r.weaponTitles[weaponItemKey][mobID] = title
	return nil
}

func (r *MemoryCharacterRepository) ClearWeaponTitle(weaponItemKey int64, mobID int16, clearItemKey int64, delFlag byte, optType int16, optData int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if m, ok := r.weaponTitles[weaponItemKey]; ok {
		delete(m, mobID)
	}
	return nil
}

func (r *MemoryCharacterRepository) DeleteMiniHouseEffectResource(owner int, itemKey, demandMoney int64, optType int16, remainLot int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if demandMoney > 0 {
		if ch, ok := r.characters[owner]; ok {
			ch.TCharacterState.Money -= demandMoney
			if ch.TCharacterState.Money < 0 {
				ch.TCharacterState.Money = 0
			}
		}
	}
	if itemKey != 0 {
		if items, ok := r.items[owner]; ok {
			filtered := items[:0]
			for _, it := range items {
				if it.ItemKey != itemKey {
					filtered = append(filtered, it)
				}
			}
			r.items[owner] = filtered
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetGuildTournamentResult(gtNo int, typ int16, data int, guildNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) SetGuildTournamentChoice(gtNo int, limitLevel int, limitMember int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetMiniHouseGates(charNo int) ([]*MiniHousePortalData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*MiniHousePortalData
	for _, g := range r.miniHouseGates[charNo] {
		cp := *g
		res = append(res, &cp)
	}
	return res, nil
}

func (r *MemoryCharacterRepository) SetMiniHouseGate(charNo, itemID int, bType byte, endDate time.Time, itemKey int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, g := range r.miniHouseGates[charNo] {
		if g.ItemID == itemID {
			return fmt.Errorf("duplicate minihouse gate item %d", itemID)
		}
	}
	r.miniHouseGates[charNo] = append(r.miniHouseGates[charNo], &MiniHousePortalData{
		CharNo:  charNo,
		ItemID:  itemID,
		Type:    bType,
		EndDate: endDate,
	})
	return nil
}

func (r *MemoryCharacterRepository) DelMiniHouseGate(charNo, itemID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.miniHouseGates[charNo]
	out := list[:0]
	for _, g := range list {
		if g.ItemID != itemID {
			out = append(out, g)
		}
	}
	r.miniHouseGates[charNo] = out
	return nil
}

func (r *MemoryCharacterRepository) CreatePet(ownerNo int, name string, mind, stress byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	petNo := r.nextPetNo
	r.nextPetNo++
	r.pets[petNo] = &PetData{
		PetNo:   petNo,
		OwnerNo: ownerNo,
		PetName: name,
		Mind:    mind,
		Stress:  stress,
	}
	return petNo, nil
}

func (r *MemoryCharacterRepository) DeletePet(ownerNo, petNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.pets, petNo)
	return nil
}

func (r *MemoryCharacterRepository) GetWeddingStatus(charNo int) (*WeddingStatusData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if w, ok := r.weddings[charNo]; ok {
		cp := *w
		return &cp, nil
	}
	return nil, nil
}

func (r *MemoryCharacterRepository) SetWeddingStatus(charNo, partnerCharNo int, status, subStatus int16, when, subWhen time.Time, marriageCount, divorceCancel byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.weddings[charNo] = &WeddingStatusData{
		CharNo:        charNo,
		PartnerCharNo: partnerCharNo,
		Status:        status,
		SubStatus:     subStatus,
		When:          when,
		SubWhen:       subWhen,
		MarriageCount: marriageCount,
		DivorceCancel: divorceCancel,
	}
	return nil
}

func (r *MemoryCharacterRepository) ProposeWedding(charANo, charBNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.weddings[charANo] = &WeddingStatusData{CharNo: charANo, PartnerCharNo: charBNo, Status: 1, When: now, SubWhen: now}
	r.weddings[charBNo] = &WeddingStatusData{CharNo: charBNo, PartnerCharNo: charANo, Status: 1, When: now, SubWhen: now}
	return nil
}

func (r *MemoryCharacterRepository) CancelWeddingPropose(charANo, charBNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.weddings, charANo)
	delete(r.weddings, charBNo)
	return nil
}

func (r *MemoryCharacterRepository) CompleteWedding(charANo, charBNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for _, pair := range [][2]int{{charANo, charBNo}, {charBNo, charANo}} {
		w := r.weddings[pair[0]]
		if w == nil {
			w = &WeddingStatusData{CharNo: pair[0], PartnerCharNo: pair[1]}
			r.weddings[pair[0]] = w
		}
		w.PartnerCharNo = pair[1]
		w.Status = 2
		w.SubStatus = 0
		w.When = now
		w.MarriageCount++
	}
	return nil
}

func (r *MemoryCharacterRepository) RequestAgreeDivorce(charANo, charBNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	exp := time.Now().Add(10 * time.Minute)
	if w := r.weddings[charANo]; w != nil {
		w.SubStatus = 1
		w.SubWhen = exp
	}
	if w := r.weddings[charBNo]; w != nil {
		w.SubStatus = 6
		w.SubWhen = exp
	}
	return nil
}

func (r *MemoryCharacterRepository) RequestForceDivorce(charANo, charBNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	exp := time.Now().Add(72 * time.Hour)
	if w := r.weddings[charANo]; w != nil {
		w.SubStatus = 3
		w.SubWhen = exp
	}
	if w := r.weddings[charBNo]; w != nil {
		w.SubStatus = 4
		w.SubWhen = exp
	}
	return nil
}

func (r *MemoryCharacterRepository) CancelDivorce(charANo, charBNo int, incCount byte) (byte, byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var subA, subB byte
	if w := r.weddings[charANo]; w != nil {
		subA = byte(w.SubStatus)
		if w.SubStatus == 1 {
			w.SubStatus = 5
		} else if w.SubStatus == 3 {
			w.SubStatus = 7
		}
		w.DivorceCancel += incCount
		w.SubWhen = w.When
	}
	if w := r.weddings[charBNo]; w != nil {
		subB = byte(w.SubStatus)
		w.SubStatus = 0
		w.SubWhen = w.When
	}
	return subA, subB, nil
}

func (r *MemoryCharacterRepository) DoDivorce(charANo, charBNo int, byItem bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.weddings, charANo)
	delete(r.weddings, charBNo)
	return nil
}

func (r *MemoryCharacterRepository) EscapeDivorceState(charNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w := r.weddings[charNo]; w != nil {
		w.SubStatus = 0
		w.SubWhen = w.When
	}
	return nil
}

func (r *MemoryCharacterRepository) GetWeddingPartnerInfo(partnerCharNo int) (byte, time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ch, ok := r.characters[partnerCharNo]; ok {
		return ch.Level, time.Now(), nil
	}
	return 0, time.Time{}, nil
}

func (r *MemoryCharacterRepository) GetHolyPromiseRemainMoney(charNo int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.holyPromiseMoney[charNo], nil
}

func (r *MemoryCharacterRepository) SetHolyPromiseMoney(charNo int, rewardMoney int64, levelGap int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.holyPromiseMoney[charNo] += rewardMoney
	return nil
}

func (r *MemoryCharacterRepository) WithdrawHolyPromiseRemainMoney(charNo int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	amt := r.holyPromiseMoney[charNo]
	r.holyPromiseMoney[charNo] = 0
	return amt, nil
}

func (r *MemoryCharacterRepository) GetDiceExpect(charNo int, itemKey int64) (int, int64, error) {
	return 0, 0, nil
}

func (r *MemoryCharacterRepository) GetDiceRank(charNo int, itemKey int64) ([]*DiceRankEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*DiceRankEntry(nil), r.diceRank[charNo]...), nil
}

func (r *MemoryCharacterRepository) SetDiceGame(charNo int, itemKey int64, dice string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextDiceGameNo
	r.nextDiceGameNo++
	return id, nil
}

func (r *MemoryCharacterRepository) SetDiceLog(gameNo, customCharNo, slotNo int, betMoney int64, betRate int, status byte, taxRate int, dividend int64) error {
	return nil
}

func (r *MemoryCharacterRepository) SetDicePay(charNo int, itemKey int64, totalDividend int64) error {
	return nil
}

func (r *MemoryCharacterRepository) SetDiceResult(gameNo, customCharNo int, totalDividend int64) error {
	return nil
}

func (r *MemoryCharacterRepository) GetGambleCoin(charNo int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.gambleCoins[charNo], nil
}

func (r *MemoryCharacterRepository) AddGambleCoin(charNo int, itemKey, coin int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gambleCoins[charNo] += coin
	return nil
}

func (r *MemoryCharacterRepository) SubGambleCoin(charNo int, coin int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.gambleCoins[charNo] < coin {
		return fmt.Errorf("insufficient gamble coins")
	}
	r.gambleCoins[charNo] -= coin
	return nil
}

func (r *MemoryCharacterRepository) GetGambleExchangedCoin(charNo int, since time.Time) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.gambleExchanged[charNo], nil
}

func (r *MemoryCharacterRepository) ExchangeMoneyToCoin(charNo int, money, coin, limit int64, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		if ch.TCharacterState.Money < money {
			return fmt.Errorf("insufficient money")
		}
		if limit > 0 && r.gambleExchanged[charNo]+coin > limit {
			return fmt.Errorf("daily exchange limit reached")
		}
		ch.TCharacterState.Money -= money
	}
	r.gambleCoins[charNo] += coin
	r.gambleExchanged[charNo] += coin
	return nil
}

func (r *MemoryCharacterRepository) ExchangeCoinToMoney(charNo int, money, coin int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.gambleCoins[charNo] < coin {
		return fmt.Errorf("insufficient coins")
	}
	r.gambleCoins[charNo] -= coin
	if ch, ok := r.characters[charNo]; ok {
		ch.TCharacterState.Money += money
	}
	return nil
}

func (r *MemoryCharacterRepository) SetGambleDice(dice string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextDiceGameNo
	r.nextDiceGameNo++
	return id, nil
}

func (r *MemoryCharacterRepository) SetGambleDiceLog(gameNo, charNo, slotNo int, betMoney int64, betRate int, status byte, taxRate int, dividend int64) error {
	if status == 1 {
		r.mu.Lock()
		if ch, ok := r.characters[charNo]; ok {
			found := false
			for _, e := range r.gambleDiceRank {
				if e.CharName == ch.Name {
					e.WinCount++
					found = true
					break
				}
			}
			if !found {
				r.gambleDiceRank = append(r.gambleDiceRank, &DiceRankEntry{CharName: ch.Name, WinCount: 1})
			}
		}
		r.mu.Unlock()
	}
	return nil
}

func (r *MemoryCharacterRepository) SetGambleDiceResult(gameNo, charNo int, totalDividend int64) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if totalDividend > 0 {
		r.gambleCoins[charNo] += totalDividend
	}
	return r.gambleCoins[charNo], nil
}

func (r *MemoryCharacterRepository) GetGambleDiceRank() ([]*DiceRankEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*DiceRankEntry, 0, len(r.gambleDiceRank))
	for _, e := range r.gambleDiceRank {
		cp := *e
		res = append(res, &cp)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].WinCount > res[j].WinCount })
	return res, nil
}

func (r *MemoryCharacterRepository) GetJackpotCoin(gameName string) ([]*JackpotCoinEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := r.gambleJackpot[gameName]
	var res []*JackpotCoinEntry
	for u, c := range m {
		res = append(res, &JackpotCoinEntry{BattingUnit: u, JackpotCoin: c})
	}
	return res, nil
}

func (r *MemoryCharacterRepository) GetJackpotRanking(gameName string, battingUnit int64) ([]*JackpotRankEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*JackpotRankEntry(nil), r.gambleJackpotRank[gameName]...), nil
}

func (r *MemoryCharacterRepository) SetSlotMachineResult(charNo int, battingUnit int64, batting byte, gameResult string, isWin bool, prizeCoin int64, isJackpot bool, jackpotRate int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	betCoins := battingUnit * int64(batting)
	if r.gambleCoins[charNo] < betCoins {
		return fmt.Errorf("insufficient coins")
	}
	r.gambleCoins[charNo] -= betCoins
	if isWin && prizeCoin > 0 {
		r.gambleCoins[charNo] += prizeCoin
	}
	if r.gambleJackpot["SlotMachine"] == nil {
		r.gambleJackpot["SlotMachine"] = make(map[int64]int64)
	}
	r.gambleJackpot["SlotMachine"][battingUnit] += betCoins
	if isJackpot {
		if ch, ok := r.characters[charNo]; ok {
			r.gambleJackpotRank["SlotMachine"] = append(r.gambleJackpotRank["SlotMachine"], &JackpotRankEntry{
				CharName:    ch.Name,
				JackpotCoin: prizeCoin,
			})
		}
		r.gambleJackpot["SlotMachine"][battingUnit] = 0
	}
	return nil
}

