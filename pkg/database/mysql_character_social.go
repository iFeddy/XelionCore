package database

import (
	"database/sql"
	"fmt"
)

// GetFriends recupera la lista de amistades del personaje invocando p_Friend_GetAll2.
func (r *MySQLCharacterRepository) GetFriends(charNo int) ([]*FriendInfo, error) {
	rows, err := r.db.Query("CALL p_Friend_GetAll2(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Friend_GetAll2(%d): %w", charNo, err)
	}
	defer rows.Close()

	var friends []*FriendInfo
	for rows.Next() {
		f := &FriendInfo{}
		var flags int
		var charID, loginZone string
		var class, level byte
		var lastLogin sql.NullTime

		err := rows.Scan(&f.CharNo, &flags, &charID, &class, &level, &loginZone, &lastLogin)
		if err == nil {
			f.CharName = charID
			f.Class = class
			f.Level = level
			f.Map = loginZone
			if lastLogin.Valid {
				f.CreatedAt = lastLogin.Time
			}
			friends = append(friends, f)
		}
	}
	return friends, nil
}

// AddFriend registra un nuevo amigo en tFriend invocando p_Friend_Set.
func (r *MySQLCharacterRepository) AddFriend(charNo, friendCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Friend_Set(50, ?, ?, 0, @out_nRet)", charNo, friendCharNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Friend_Set: %w", err)
	}
	return nil
}

// RemoveFriend elimina una relación de amistad invocando p_Friend_Del.
func (r *MySQLCharacterRepository) RemoveFriend(charNo, friendCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Friend_Del(?, ?, @out_nRet)", charNo, friendCharNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Friend_Del: %w", err)
	}
	return nil
}

// GetFriendPoints consulta los puntos de amistad acumulados invocando usp_Character_getFriendpoint.
func (r *MySQLCharacterRepository) GetFriendPoints(charNo int) (int, error) {
	rows, err := r.db.Query("CALL usp_Character_getFriendpoint(?, @out_nRet)", charNo)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_Character_getFriendpoint(%d): %w", charNo, err)
	}
	defer rows.Close()

	var pts int
	if rows.Next() {
		_ = rows.Scan(&pts)
	}
	return pts, nil
}

// SetFriendPoints actualiza los puntos de amistad acumulados invocando usp_Character_setFriendpoint.
func (r *MySQLCharacterRepository) SetFriendPoints(charNo int, points int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setFriendpoint(?, ?, @out_nRet)", charNo, points).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setFriendpoint: %w", err)
	}
	return nil
}

// GetKingdomQuests recupera todas las definiciones de Kingdom Quests invocando p_KingdomQuest_GetAll.
func (r *MySQLCharacterRepository) GetKingdomQuests() ([]*KingdomQuestData, error) {
	rows, err := r.db.Query("CALL p_KingdomQuest_GetAll()")
	if err != nil {
		return nil, fmt.Errorf("error calling p_KingdomQuest_GetAll: %w", err)
	}
	defer rows.Close()

	var list []*KingdomQuestData
	for rows.Next() {
		kq := &KingdomQuestData{}
		if err := rows.Scan(&kq.QuestNo, &kq.Data); err == nil {
			list = append(list, kq)
		}
	}
	return list, nil
}

// GetKingdomQuest consulta la configuración de una Kingdom Quest específica invocando p_KingdomQuest_Get.
func (r *MySQLCharacterRepository) GetKingdomQuest(questNo int) (*KingdomQuestData, error) {
	var data string
	var ret int
	err := r.db.QueryRow("CALL p_KingdomQuest_Get(?, @out_sData, @out_nRet)", questNo).Scan(&data, &ret)
	if err != nil {
		return nil, fmt.Errorf("error calling p_KingdomQuest_Get(%d): %w", questNo, err)
	}
	return &KingdomQuestData{QuestNo: questNo, Data: data}, nil
}

// SetKingdomQuest guarda la configuración de una Kingdom Quest invocando p_KingdomQuest_Set.
func (r *MySQLCharacterRepository) SetKingdomQuest(questNo int, data string) error {
	var ret int
	err := r.db.QueryRow("CALL p_KingdomQuest_Set(?, ?, @out_nRet)", questNo, data).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_KingdomQuest_Set: %w", err)
	}
	return nil
}

// DelKingdomQuest remueve una Kingdom Quest invocando p_KingdomQuest_Del.
func (r *MySQLCharacterRepository) DelKingdomQuest(questNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_KingdomQuest_Del(?, @out_nRet)", questNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_KingdomQuest_Del: %w", err)
	}
	return nil
}

// CheckHolyPromiseRing valida la relación de compromiso sagrado invocando p_HolyPromise_CheckRing.
func (r *MySQLCharacterRepository) CheckHolyPromiseRing(charNo, upCharNo int) (bool, error) {
	var ret int
	err := r.db.QueryRow("CALL p_HolyPromise_CheckRing(?, ?, @out_nRet)", charNo, upCharNo).Scan(&ret)
	if err != nil {
		return false, fmt.Errorf("error calling p_HolyPromise_CheckRing: %w", err)
	}
	return ret == 1, nil
}

// GetMasterPupils recupera relaciones de maestro y alumnos invocando p_MasPup_GetAll.
func (r *MySQLCharacterRepository) GetMasterPupils(charNo int) ([]*MasterPupilData, error) {
	rows, err := r.db.Query("CALL p_MasPup_GetAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_MasPup_GetAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*MasterPupilData
	for rows.Next() {
		var master, pupil, flags int
		if err := rows.Scan(&master, &pupil, &flags); err == nil {
			if master == charNo {
				list = append(list, &MasterPupilData{CharNo: pupil, IsMaster: false})
			} else {
				list = append(list, &MasterPupilData{CharNo: master, IsMaster: true})
			}
		}
	}
	return list, nil
}

// SetMaster vincula un maestro al personaje invocando p_MasPup_SetMaster.
func (r *MySQLCharacterRepository) SetMaster(charNo int, masterCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_MasPup_SetMaster(?, ?, 0, @out_nRet)", charNo, masterCharNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_MasPup_SetMaster: %w", err)
	}
	return nil
}

// SetPupil vincula un alumno al personaje invocando p_MasPup_SetPupils.
func (r *MySQLCharacterRepository) SetPupil(charNo int, pupilCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_MasPup_SetPupils(?, ?, 0, @out_nRet)", charNo, pupilCharNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_MasPup_SetPupils: %w", err)
	}
	return nil
}

// DelMaster desvincula el maestro invocando p_MasPup_DelMaster.
func (r *MySQLCharacterRepository) DelMaster(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_MasPup_DelMaster(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_MasPup_DelMaster: %w", err)
	}
	return nil
}

// DelPupils desvincula a un alumno invocando p_MasPup_DelPupils.
func (r *MySQLCharacterRepository) DelPupils(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_MasPup_DelPupils(?, 0, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_MasPup_DelPupils: %w", err)
	}
	return nil
}

// DelAllMasterPupils remueve todas las relaciones de maestro y alumnos invocando p_MasPup_DelAll.
func (r *MySQLCharacterRepository) DelAllMasterPupils(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_MasPup_DelAll(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_MasPup_DelAll: %w", err)
	}
	return nil
}

// GetPetInfo consulta los datos de la mascota invocando usp_Pet_getInfo.
func (r *MySQLCharacterRepository) GetPetInfo(ownerNo int) (*PetData, error) {
	rows, err := r.db.Query("CALL usp_Pet_getInfo(?, @out_nRet)", ownerNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Pet_getInfo(%d): %w", ownerNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		p := &PetData{OwnerNo: ownerNo}
		var sName string
		var mind, stress byte
		if err := rows.Scan(&p.OwnerNo, &sName, &mind, &stress); err == nil {
			p.PetName = sName
			p.Mind = mind
			p.Stress = stress
			return p, nil
		}
	}
	return nil, nil
}

// SetPetName renombra a una mascota invocando usp_Pet_setName.
func (r *MySQLCharacterRepository) SetPetName(petNo int, ownerNo int, name string) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Pet_setName(?, ?, @out_nRet)", petNo, name).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Pet_setName: %w", err)
	}
	return nil
}

// SetPetTendency actualiza el temperamento de la mascota invocando usp_Pet_setTendency.
func (r *MySQLCharacterRepository) SetPetTendency(petNo int, tendency byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Pet_setTendency(?, ?, 0, @out_nRet)", petNo, tendency).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Pet_setTendency: %w", err)
	}
	return nil
}

// GetMiniHouseURL recupera la URL o blog de la minihouse invocando p_MiniHouseURL_Get.
func (r *MySQLCharacterRepository) GetMiniHouseURL(charNo int) (string, error) {
	var url string
	row := r.db.QueryRow("CALL p_MiniHouseURL_Get(?, @out_sURL)", charNo)
	if err := row.Scan(&url); err != nil {
		return "", nil
	}
	return url, nil
}

// SetMiniHouseURL actualiza la URL o blog de la minihouse invocando p_MiniHouseURL_Set.
func (r *MySQLCharacterRepository) SetMiniHouseURL(charNo int, url string) error {
	_, err := r.db.Exec("CALL p_MiniHouseURL_Set(?, ?)", charNo, url)
	if err != nil {
		return fmt.Errorf("error calling p_MiniHouseURL_Set: %w", err)
	}
	return nil
}

// GetMiniHouseVisitors consulta el número de visitantes a la minihouse invocando usp_Minihouse_getVisitorCount.
func (r *MySQLCharacterRepository) GetMiniHouseVisitors(charNo int) (int, error) {
	rows, err := r.db.Query("CALL usp_Minihouse_getVisitorCount(?, @out_nRet)", charNo)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_Minihouse_getVisitorCount(%d): %w", charNo, err)
	}
	defer rows.Close()

	var visitors int
	if rows.Next() {
		_ = rows.Scan(&visitors)
	}
	return visitors, nil
}

// SetMiniHouseVisitors incrementa el contador de visitantes invocando usp_Minihouse_setVisitorCount.
func (r *MySQLCharacterRepository) SetMiniHouseVisitors(charNo int, visitors int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Minihouse_setVisitorCount(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Minihouse_setVisitorCount: %w", err)
	}
	return nil
}
