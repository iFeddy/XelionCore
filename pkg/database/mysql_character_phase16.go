package database

import (
	"fmt"
	"strings"
)

// RenameGuild renombra un gremio existente invocando CALL p_Guild_setName.
func (r *MySQLCharacterRepository) RenameGuild(guildID int, newName string) (int, error) {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_setName(?, ?, @out_nRet)", guildID, strings.TrimSpace(newName)).Scan(&ret)
	if err != nil {
		return -1393, fmt.Errorf("error calling p_Guild_setName(%d): %w", guildID, err)
	}
	return ret, nil
}

// RetypeGuild actualiza el tipo oficial del gremio invocando CALL p_Guild_SetType.
func (r *MySQLCharacterRepository) RetypeGuild(guildID int, newType byte) (int, error) {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_SetType(?, ?, @out_nRet)", guildID, newType).Scan(&ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_Guild_SetType(%d): %w", guildID, err)
	}
	return ret, nil
}

// SaveMysteryBolt persiste bMysteryBolt en tCharacterSingleOption invocando CALL usp_Character_setMysteryBolt.
func (r *MySQLCharacterRepository) SaveMysteryBolt(charNo int, mysteryType byte) (int, error) {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setMysteryBolt(?, ?, @out_nRet)", charNo, mysteryType).Scan(&ret)
	if err != nil {
		return -3193, fmt.Errorf("error calling usp_Character_setMysteryBolt(%d): %w", charNo, err)
	}
	return ret, nil
}

// SaveLinkScrollDestMap persiste el mapa destino de LinkScroll en tCharacter.sKQMap invocando CALL usp_LinkScroll_SetDestMap.
func (r *MySQLCharacterRepository) SaveLinkScrollDestMap(charNo int, destMap string) error {
	var ret int
	err := r.db.QueryRow("CALL usp_LinkScroll_SetDestMap(?, ?, @out_nRet)", charNo, strings.TrimSpace(destMap)).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_LinkScroll_SetDestMap(%d): %w", charNo, err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_LinkScroll_SetDestMap failed: %d", ret)
	}
	return nil
}

// UseCharacterChangeIDItem registra la habilitación de cambio de nombre en tCharacterChangeID invocando CALL usp_Character_setChangeIDItemUsed.
func (r *MySQLCharacterRepository) UseCharacterChangeIDItem(charNo int, bInit byte, userNo int, itemKey int64, itemOption int16, itemLot int16) (int, error) {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Character_setChangeIDItemUsed(?, ?, ?, ?, ?, ?, @out_nRet)",
		charNo, bInit, userNo, itemKey, itemOption, itemLot,
	).Scan(&ret)
	if err != nil {
		return -3036, fmt.Errorf("error calling usp_Character_setChangeIDItemUsed(%d): %w", charNo, err)
	}
	return ret, nil
}

// SetAdminLevel actualiza tCharacter.nAdminLevel invocando CALL p_Char_SetAdminLevel.
func (r *MySQLCharacterRepository) SetAdminLevel(charNo int, level byte) (int, error) {
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetAdminLevel(?, ?, @out_nRet)", charNo, level).Scan(&ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_Char_SetAdminLevel(%d): %w", charNo, err)
	}
	return ret, nil
}

// CreateGMItemRewardGroup registra un lote de entrega GM en tGMItemReward invocando CALL usp_GM_getItemRewardGroupNo.
func (r *MySQLCharacterRepository) CreateGMItemRewardGroup(operID, ip, cause string) (int64, error) {
	var groupNo int64
	err := r.db.QueryRow(
		"CALL usp_GM_getItemRewardGroupNo(?, ?, ?, @out_nRet)",
		strings.TrimSpace(operID), strings.TrimSpace(ip), strings.TrimSpace(cause),
	).Scan(&groupNo)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_GM_getItemRewardGroupNo: %w", err)
	}
	if groupNo <= 0 {
		return groupNo, fmt.Errorf("usp_GM_getItemRewardGroupNo failed with code %d", groupNo)
	}
	return groupNo, nil
}

// GrantGMItemReward entrega un ítem GM en tItem y registra auditoría en tGMItemRewardLog invocando CALL usp_GM_setItemReward.
func (r *MySQLCharacterRepository) GrantGMItemReward(groupNo int64, charNo int, storageType byte, itemID int, itemKey int64) (int, error) {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_GM_setItemReward(?, ?, ?, ?, ?, @out_nRet)",
		groupNo, charNo, storageType, itemID, itemKey,
	).Scan(&ret)
	if err != nil {
		return -1111, fmt.Errorf("error calling usp_GM_setItemReward: %w", err)
	}
	if ret != 0 {
		return ret, fmt.Errorf("usp_GM_setItemReward failed with code %d", ret)
	}
	return ret, nil
}

