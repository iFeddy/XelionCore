package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (r *MySQLCharacterRepository) SetWeaponTitle(weaponItemKey, licenseItemKey int64, mobID int16, titleCharNo int, title string) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL p_Item_SetWeaponTitle(?, ?, ?, ?, ?, @out_nRet)",
		weaponItemKey, licenseItemKey, mobID, titleCharNo, title,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Item_SetWeaponTitle: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_SetWeaponTitle failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) ClearWeaponTitle(weaponItemKey int64, mobID int16, clearItemKey int64, delFlag byte, optType int16, optData int64) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL p_Item_ClrWeaponTitle(?, ?, ?, ?, ?, ?, @out_nRet)",
		weaponItemKey, mobID, clearItemKey, delFlag, optType, optData,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Item_ClrWeaponTitle: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_ClrWeaponTitle failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) DeleteMiniHouseEffectResource(owner int, itemKey, demandMoney int64, optType int16, remainLot int) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL p_Item_MiniHouse_Effect_Resource_Del(?, ?, ?, ?, ?, @out_nRet)",
		owner, itemKey, demandMoney, optType, remainLot,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Item_MiniHouse_Effect_Resource_Del: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_MiniHouse_Effect_Resource_Del failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetGuildTournamentResult(gtNo int, typ int16, data int, guildNo int) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_GuildTournament_SetResult(?, ?, ?, ?, @out_nRet)",
		gtNo, typ, data, guildNo,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_GuildTournament_SetResult: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_GuildTournament_SetResult failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetGuildTournamentChoice(gtNo int, limitLevel int, limitMember int) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_GuildTournament_setChoice(?, ?, ?, @out_nRet)",
		gtNo, limitLevel, limitMember,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_GuildTournament_setChoice: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_GuildTournament_setChoice failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) GetMiniHouseGates(charNo int) ([]*MiniHousePortalData, error) {
	rows, err := r.db.Query("CALL usp_Minihouse_getGateList(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Minihouse_getGateList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*MiniHousePortalData
	for rows.Next() {
		g := &MiniHousePortalData{}
		var dEnd sql.NullTime
		if err := rows.Scan(&g.CharNo, &g.ItemID, &g.Type, &dEnd); err == nil {
			if dEnd.Valid {
				g.EndDate = dEnd.Time
			}
			list = append(list, g)
		}
	}
	return list, nil
}

func (r *MySQLCharacterRepository) SetMiniHouseGate(charNo, itemID int, bType byte, endDate time.Time, itemKey int64) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_Minihouse_setGate(?, ?, ?, ?, ?, @out_nRet)",
		charNo, itemID, bType, endDate, itemKey,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Minihouse_setGate: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Minihouse_setGate failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) DelMiniHouseGate(charNo, itemID int) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_Minihouse_delGate(?, ?, @out_nRet)",
		charNo, itemID,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Minihouse_delGate: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Minihouse_delGate failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) CreatePet(ownerNo int, name string, mind, stress byte) (int, error) {
	var petNo int
	if err := r.db.QueryRow(
		"CALL usp_Pet_setCreate(?, ?, ?, ?, @out_nRet)",
		ownerNo, name, mind, stress,
	).Scan(&petNo); err != nil {
		return 0, fmt.Errorf("error calling usp_Pet_setCreate: %w", err)
	}
	if petNo <= 0 {
		return 0, fmt.Errorf("usp_Pet_setCreate failed: %d", petNo)
	}
	return petNo, nil
}

func (r *MySQLCharacterRepository) DeletePet(ownerNo, petNo int) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Pet_setDelete(?, ?, @out_nRet)", ownerNo, petNo).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Pet_setDelete: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) GetWeddingStatus(charNo int) (*WeddingStatusData, error) {
	rows, err := r.db.Query("CALL p_Char_WeddingStatus_Get(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_WeddingStatus_Get(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		w := &WeddingStatusData{CharNo: charNo}
		var dWhen, dSubWhen sql.NullTime
		if err := rows.Scan(&w.PartnerCharNo, &w.Status, &w.SubStatus, &dWhen, &dSubWhen, &w.MarriageCount, &w.DivorceCancel); err == nil {
			if dWhen.Valid {
				w.When = dWhen.Time
			}
			if dSubWhen.Valid {
				w.SubWhen = dSubWhen.Time
			}
			return w, nil
		}
	}
	return nil, nil
}

func (r *MySQLCharacterRepository) SetWeddingStatus(charNo, partnerCharNo int, status, subStatus int16, when, subWhen time.Time, marriageCount, divorceCancel byte) error {
	_, err := r.db.Exec(
		"CALL p_Char_WeddingStatus_Set(?, ?, ?, ?, ?, ?, ?, ?)",
		charNo, partnerCharNo, status, subStatus, when, subWhen, marriageCount, divorceCancel,
	)
	return err
}

func (r *MySQLCharacterRepository) ProposeWedding(charANo, charBNo int) error {
	var ret int
	if err := r.db.QueryRow("CALL p_Char_Wedding_Propose(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Char_Wedding_Propose: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Char_Wedding_Propose failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) CancelWeddingPropose(charANo, charBNo int) error {
	var ret int
	if err := r.db.QueryRow("CALL p_Char_Wedding_Cancel_Propose(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Char_Wedding_Cancel_Propose: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) CompleteWedding(charANo, charBNo int) error {
	var ret int
	if err := r.db.QueryRow("CALL p_Char_Wedding_Do(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret); err != nil {
		return fmt.Errorf("error calling p_Char_Wedding_Do: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) RequestAgreeDivorce(charANo, charBNo int) error {
	var ret int
	_ = r.db.QueryRow("CALL p_Char_Wedding_Divorce_by_Agree(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret)
	return nil
}

func (r *MySQLCharacterRepository) RequestForceDivorce(charANo, charBNo int) error {
	var ret int
	_ = r.db.QueryRow("CALL p_Char_Wedding_Divorce_by_Force(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret)
	return nil
}

func (r *MySQLCharacterRepository) CancelDivorce(charANo, charBNo int, incCount byte) (byte, byte, error) {
	var subA, subB byte
	var ret int
	rows, err := r.db.Query("CALL p_Char_Wedding_Divorce_Cancel(?, ?, ?, @out_SubA, @out_SubB, @out_nRet)", charANo, charBNo, incCount)
	if err != nil {
		return 0, 0, fmt.Errorf("error calling p_Char_Wedding_Divorce_Cancel: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		_ = rows.Scan(&subA, &subB, &ret)
	}
	return subA, subB, nil
}

func (r *MySQLCharacterRepository) DoDivorce(charANo, charBNo int, byItem bool) error {
	var ret int
	proc := "p_Char_Wedding_Divorce_Do"
	if byItem {
		proc = "p_Char_Wedding_Divorce_Do_by_Item"
	}
	if err := r.db.QueryRow("CALL "+proc+"(?, ?, @out_nRet)", charANo, charBNo).Scan(&ret); err != nil {
		return fmt.Errorf("error calling %s: %w", proc, err)
	}
	return nil
}

func (r *MySQLCharacterRepository) EscapeDivorceState(charNo int) error {
	var ret int
	_ = r.db.QueryRow("CALL p_Char_Wedding_Escape_Divorce_State(?, @out_nRet)", charNo).Scan(&ret)
	return nil
}

func (r *MySQLCharacterRepository) GetWeddingPartnerInfo(partnerCharNo int) (byte, time.Time, error) {
	rows, err := r.db.Query("CALL p_Char_Wedding_Get_Partner_Info(?)", partnerCharNo)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("error calling p_Char_Wedding_Get_Partner_Info: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var lv byte
		var lastLogin sql.NullTime
		if err := rows.Scan(&lv, &lastLogin); err == nil {
			var tm time.Time
			if lastLogin.Valid {
				tm = lastLogin.Time
			}
			return lv, tm, nil
		}
	}
	return 0, time.Time{}, nil
}

func (r *MySQLCharacterRepository) GetHolyPromiseRemainMoney(charNo int) (int64, error) {
	var money int64
	if err := r.db.QueryRow("CALL usp_HolyPromise_getRemainMoney(?, @out_nRet)", charNo).Scan(&money); err != nil {
		return 0, fmt.Errorf("error calling usp_HolyPromise_getRemainMoney: %w", err)
	}
	return money, nil
}

func (r *MySQLCharacterRepository) SetHolyPromiseMoney(charNo int, rewardMoney int64, levelGap int) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_HolyPromise_setMoney(?, ?, ?, @out_nRet)", charNo, rewardMoney, levelGap).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_HolyPromise_setMoney: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) WithdrawHolyPromiseRemainMoney(charNo int) (int64, error) {
	var money int64
	if err := r.db.QueryRow("CALL usp_HolyPromise_withDrawRemainMoney(?, @out_nRet)", charNo).Scan(&money); err != nil {
		return 0, fmt.Errorf("error calling usp_HolyPromise_withDrawRemainMoney: %w", err)
	}
	return money, nil
}

func (r *MySQLCharacterRepository) GetDiceExpect(charNo int, itemKey int64) (int, int64, error) {
	rows, err := r.db.Query("CALL usp_Dice_getDiceExpect(?, ?, @out_nRet)", charNo, itemKey)
	if err != nil {
		return 0, 0, fmt.Errorf("error calling usp_Dice_getDiceExpect: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var cnt int
		var div int64
		if err := rows.Scan(&cnt, &div); err == nil {
			return cnt, div, nil
		}
	}
	return 0, 0, nil
}

func (r *MySQLCharacterRepository) GetDiceRank(charNo int, itemKey int64) ([]*DiceRankEntry, error) {
	rows, err := r.db.Query("CALL usp_Dice_getDiceRank(?, ?, @out_nRet)", charNo, itemKey)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Dice_getDiceRank: %w", err)
	}
	defer rows.Close()
	var list []*DiceRankEntry
	for rows.Next() {
		e := &DiceRankEntry{}
		if err := rows.Scan(&e.CharName, &e.WinCount); err == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

func (r *MySQLCharacterRepository) SetDiceGame(charNo int, itemKey int64, dice string) (int, error) {
	var gameNo int
	if err := r.db.QueryRow("CALL usp_Dice_setDice(?, ?, ?, @out_nRet)", charNo, itemKey, dice).Scan(&gameNo); err != nil {
		return 0, fmt.Errorf("error calling usp_Dice_setDice: %w", err)
	}
	return gameNo, nil
}

func (r *MySQLCharacterRepository) SetDiceLog(gameNo, customCharNo, slotNo int, betMoney int64, betRate int, status byte, taxRate int, dividend int64) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_Dice_setDiceLog(?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		gameNo, customCharNo, slotNo, betMoney, betRate, status, taxRate, dividend,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Dice_setDiceLog: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetDicePay(charNo int, itemKey int64, totalDividend int64) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Dice_setDicePay(?, ?, ?, @out_nRet)", charNo, itemKey, totalDividend).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Dice_setDicePay: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetDiceResult(gameNo, customCharNo int, totalDividend int64) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Dice_setDiceResult(?, ?, ?, @out_nRet)", gameNo, customCharNo, totalDividend).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Dice_setDiceResult: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) GetGambleCoin(charNo int) (int64, error) {
	var coin int64
	if err := r.db.QueryRow("CALL usp_Gamble_getCoin(?, @out_nRet)", charNo).Scan(&coin); err != nil {
		return 0, fmt.Errorf("error calling usp_Gamble_getCoin: %w", err)
	}
	return coin, nil
}

func (r *MySQLCharacterRepository) AddGambleCoin(charNo int, itemKey, coin int64) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Gamble_addCoin(?, ?, ?, @out_nRet)", charNo, itemKey, coin).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_addCoin: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Gamble_addCoin failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) SubGambleCoin(charNo int, coin int64) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Gamble_subCoin(?, ?, @out_nRet)", charNo, coin).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_subCoin: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Gamble_subCoin failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) GetGambleExchangedCoin(charNo int, since time.Time) (int64, error) {
	var coin int64
	if err := r.db.QueryRow("CALL usp_Gamble_getExchangedCoin(?, ?, @out_nRet)", charNo, since).Scan(&coin); err != nil {
		return 0, fmt.Errorf("error calling usp_Gamble_getExchangedCoin: %w", err)
	}
	return coin, nil
}

func (r *MySQLCharacterRepository) ExchangeMoneyToCoin(charNo int, money, coin, limit int64, now time.Time) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Gamble_setMoneyToCoin(?, ?, ?, ?, ?, @out_nRet)", charNo, money, coin, limit, now).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_setMoneyToCoin: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Gamble_setMoneyToCoin failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) ExchangeCoinToMoney(charNo int, money, coin int64) error {
	var ret int
	if err := r.db.QueryRow("CALL usp_Gamble_setCoinToMoney(?, ?, ?, @out_nRet)", charNo, money, coin).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_setCoinToMoney: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Gamble_setCoinToMoney failed: %d", ret)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetGambleDice(dice string) (int, error) {
	var gameNo int
	if err := r.db.QueryRow("CALL usp_Gamble_setDice(?, @out_nRet)", dice).Scan(&gameNo); err != nil {
		return 0, fmt.Errorf("error calling usp_Gamble_setDice: %w", err)
	}
	return gameNo, nil
}

func (r *MySQLCharacterRepository) SetGambleDiceLog(gameNo, charNo, slotNo int, betMoney int64, betRate int, status byte, taxRate int, dividend int64) error {
	var ret int
	if err := r.db.QueryRow(
		"CALL usp_Gamble_setDiceLog(?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		gameNo, charNo, slotNo, betMoney, betRate, status, taxRate, dividend,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_setDiceLog: %w", err)
	}
	return nil
}

func (r *MySQLCharacterRepository) SetGambleDiceResult(gameNo, charNo int, totalDividend int64) (int64, error) {
	rows, err := r.db.Query("CALL usp_Gamble_setDiceResult(?, ?, ?, @out_nRet)", gameNo, charNo, totalDividend)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_Gamble_setDiceResult: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var cNo int
		var coin int64
		if err := rows.Scan(&cNo, &coin); err == nil {
			return coin, nil
		}
	}
	return 0, nil
}

func (r *MySQLCharacterRepository) GetGambleDiceRank() ([]*DiceRankEntry, error) {
	rows, err := r.db.Query("CALL usp_Gamble_getDiceRank(@out_nRet)")
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Gamble_getDiceRank: %w", err)
	}
	defer rows.Close()
	var list []*DiceRankEntry
	for rows.Next() {
		e := &DiceRankEntry{}
		if err := rows.Scan(&e.CharName, &e.WinCount); err == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

func (r *MySQLCharacterRepository) GetJackpotCoin(gameName string) ([]*JackpotCoinEntry, error) {
	rows, err := r.db.Query("CALL usp_Gamble_getJackpotCoin(?, @out_nRet)", gameName)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Gamble_getJackpotCoin: %w", err)
	}
	defer rows.Close()
	var list []*JackpotCoinEntry
	for rows.Next() {
		e := &JackpotCoinEntry{}
		if err := rows.Scan(&e.BattingUnit, &e.JackpotCoin); err == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

func (r *MySQLCharacterRepository) GetJackpotRanking(gameName string, battingUnit int64) ([]*JackpotRankEntry, error) {
	rows, err := r.db.Query("CALL usp_Gamble_getJackpotRanking(?, ?, @out_nRet)", gameName, battingUnit)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Gamble_getJackpotRanking: %w", err)
	}
	defer rows.Close()
	var list []*JackpotRankEntry
	for rows.Next() {
		e := &JackpotRankEntry{}
		if err := rows.Scan(&e.CharName, &e.JackpotCoin); err == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

func (r *MySQLCharacterRepository) SetSlotMachineResult(charNo int, battingUnit int64, batting byte, gameResult string, isWin bool, prizeCoin int64, isJackpot bool, jackpotRate int) error {
	var ret int
	winFlag, jpFlag := 0, 0
	if isWin {
		winFlag = 1
	}
	if isJackpot {
		jpFlag = 1
	}
	if err := r.db.QueryRow(
		"CALL usp_Gamble_setSlotMachineResult(?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		charNo, battingUnit, batting, gameResult, winFlag, prizeCoin, jpFlag, jackpotRate,
	).Scan(&ret); err != nil {
		return fmt.Errorf("error calling usp_Gamble_setSlotMachineResult: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Gamble_setSlotMachineResult failed: %d", ret)
	}
	return nil
}

