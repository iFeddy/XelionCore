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

// CreateParty registra un nuevo grupo oficial en tParty invocando p_Party_Create.
func (r *MySQLCharacterRepository) CreateParty(masterCharNo int, flags int) (int, error) {
	var partyNo, ret int
	err := r.db.QueryRow("CALL p_Party_Create(?, ?, @out_nPartyNo, @out_nRet)", masterCharNo, flags).Scan(&partyNo, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_Party_Create(%d): %w", masterCharNo, err)
	}
	if ret != 0 {
		return 0, fmt.Errorf("p_Party_Create failed with code %d", ret)
	}
	return partyNo, nil
}

// AddPartyMember añade un personaje a un grupo en tPartyMember invocando p_Party_AddMember.
func (r *MySQLCharacterRepository) AddPartyMember(partyNo int, charNo int, flag byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Party_AddMember(?, ?, ?, @out_nRet)", partyNo, charNo, flag).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Party_AddMember: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Party_AddMember failed with code %d", ret)
	}
	return nil
}

// RemovePartyMember remueve un miembro de tPartyMember invocando p_Party_RemoveMember.
func (r *MySQLCharacterRepository) RemovePartyMember(partyNo int, charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Party_RemoveMember(?, ?, @out_nRet)", partyNo, charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Party_RemoveMember: %w", err)
	}
	return nil
}

// SetPartyMaster designa un nuevo líder en tParty y tPartyMember invocando p_Party_SetMaster.
func (r *MySQLCharacterRepository) SetPartyMaster(partyNo int, masterCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Party_SetMaster(?, ?, @out_nRet)", partyNo, masterCharNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Party_SetMaster: %w", err)
	}
	return nil
}

// SetPartyFlags actualiza las banderas y modo de loteo en tParty invocando p_Party_SetFlags.
func (r *MySQLCharacterRepository) SetPartyFlags(partyNo int, flags int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Party_SetFlags(?, ?, @out_nRet)", partyNo, flags).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Party_SetFlags: %w", err)
	}
	return nil
}

// DeleteParty disuelve el grupo eliminando registros de tParty/tPartyMember invocando p_Party_Delete.
func (r *MySQLCharacterRepository) DeleteParty(partyNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Party_Delete(?, @out_nRet)", partyNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Party_Delete: %w", err)
	}
	return nil
}

// GetParty consulta los datos base del grupo invocando p_Party_Get.
func (r *MySQLCharacterRepository) GetParty(partyNo int) (*PartyData, error) {
	rows, err := r.db.Query("CALL p_Party_Get(?)", partyNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Party_Get(%d): %w", partyNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		p := &PartyData{}
		if err := rows.Scan(&p.PartyNo, &p.MasterCharNo, &p.Flags); err == nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("party %d not found", partyNo)
}

// GetPartyMembers lista los miembros de un grupo invocando p_Party_GetMembers.
func (r *MySQLCharacterRepository) GetPartyMembers(partyNo int) ([]*PartyMemberData, error) {
	rows, err := r.db.Query("CALL p_Party_GetMembers(?)", partyNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Party_GetMembers(%d): %w", partyNo, err)
	}
	defer rows.Close()

	var members []*PartyMemberData
	for rows.Next() {
		m := &PartyMemberData{}
		if err := rows.Scan(&m.CharNo, &m.Flag, &m.CharID, &m.Class, &m.Level, &m.LoginZone, &m.HP, &m.SP); err == nil {
			members = append(members, m)
		}
	}
	return members, nil
}

// GetPartyByChar consulta si un personaje pertenece a un grupo invocando p_Party_GetByChar.
func (r *MySQLCharacterRepository) GetPartyByChar(charNo int) (int, int, int, error) {
	rows, err := r.db.Query("CALL p_Party_GetByChar(?)", charNo)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error calling p_Party_GetByChar(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var partyNo, masterCharNo, flags int
		if err := rows.Scan(&partyNo, &masterCharNo, &flags); err == nil {
			return partyNo, masterCharNo, flags, nil
		}
	}
	return 0, 0, 0, nil
}

// AddChatRestrict registra un bloqueo de chat en tChatBlock invocando p_Char_Chat_Restrict_Add.
func (r *MySQLCharacterRepository) AddChatRestrict(userNo int, charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_Chat_Restrict_Add(?, ?, @out_nRet)", userNo, charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Chat_Restrict_Add: %w", err)
	}
	return nil
}

// DelChatRestrict elimina un bloqueo de chat de tChatBlock invocando p_Char_Chat_Restrict_Del.
func (r *MySQLCharacterRepository) DelChatRestrict(userNo int, charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_Chat_Restrict_Del(?, ?, @out_nRet)", userNo, charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Chat_Restrict_Del: %w", err)
	}
	return nil
}

// DelAllChatRestrict elimina todos los bloqueos de chat de un usuario invocando p_Char_Chat_Restrict_DelAll.
func (r *MySQLCharacterRepository) DelAllChatRestrict(userNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_Chat_Restrict_DelAll(?, @out_nRet)", userNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Chat_Restrict_DelAll: %w", err)
	}
	return nil
}

// AddPrison ingresa a un personaje en tPrison invocando p_Prison_Add2.
func (r *MySQLCharacterRepository) AddPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error {
	var errCode int16
	err := r.db.QueryRow("CALL p_Prison_Add2(?, ?, ?, ?, ?, @out_err)", gmUserID, targetCharName, minutes, reason, remark).Scan(&errCode)
	if err != nil {
		return fmt.Errorf("error calling p_Prison_Add2: %w", err)
	}
	return nil
}

// EndPrison libera a un personaje encarcelado invocando p_Prison_End.
func (r *MySQLCharacterRepository) EndPrison(charNo int) error {
	_, err := r.db.Exec("CALL p_Prison_End(?)", charNo)
	if err != nil {
		return fmt.Errorf("error calling p_Prison_End(%d): %w", charNo, err)
	}
	return nil
}

// AlterPrison modifica el tiempo de sanción carcelaria invocando p_Prison_Alter2.
func (r *MySQLCharacterRepository) AlterPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error {
	var errCode int16
	err := r.db.QueryRow("CALL p_Prison_Alter2(?, ?, ?, ?, ?, @out_err)", gmUserID, targetCharName, minutes, reason, remark).Scan(&errCode)
	if err != nil {
		return fmt.Errorf("error calling p_Prison_Alter2: %w", err)
	}
	return nil
}

// SetClientLoadingBug registra fallos de carga del cliente en tCharacterLoading invocando p_Set_Client_Loading_Bug_Detect.
func (r *MySQLCharacterRepository) SetClientLoadingBug(charNo int, mapName string, failType int16, failCnt int16) error {
	var ret int
	err := r.db.QueryRow("CALL p_Set_Client_Loading_Bug_Detect(?, ?, ?, ?, @out_nRet)", charNo, mapName, failType, failCnt).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Set_Client_Loading_Bug_Detect: %w", err)
	}
	return nil
}

// SetSpamFilter configura el filtro anti-spam en tSpamFilter invocando usp_SpamFilter_set.
func (r *MySQLCharacterRepository) SetSpamFilter(userNo int, isSpamSpread, isSpamBlock byte, rate byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_SpamFilter_set(?, ?, ?, ?, @out_nRet)", userNo, isSpamSpread, isSpamBlock, rate).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_SpamFilter_set: %w", err)
	}
	return nil
}

// GetSpamFilter consulta la configuración anti-spam en tSpamFilter invocando usp_SpamFilter_get.
func (r *MySQLCharacterRepository) GetSpamFilter(userNo int) (bool, bool, byte, error) {
	rows, err := r.db.Query("CALL usp_SpamFilter_get(?, @out_nRet)", userNo)
	if err != nil {
		return false, false, 0, fmt.Errorf("error calling usp_SpamFilter_get(%d): %w", userNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var uNo int
		var spread, block, rate byte
		if err := rows.Scan(&uNo, &spread, &block, &rate); err == nil {
			return spread != 0, block != 0, rate, nil
		}
	}
	return false, false, 0, nil
}

// SetSpamBlockTime configura el tiempo de bloqueo en tSpamBlockTime invocando usp_SpamBlock_setTime.
func (r *MySQLCharacterRepository) SetSpamBlockTime(charNo int, blockTime, standbyTime int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_SpamBlock_setTime(?, ?, ?, @out_nRet)", charNo, blockTime, standbyTime).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_SpamBlock_setTime: %w", err)
	}
	return nil
}

// GetSpamBlockTime consulta el tiempo de bloqueo en tSpamBlockTime invocando usp_SpamBlock_getTime.
func (r *MySQLCharacterRepository) GetSpamBlockTime(charNo int) (int, int, error) {
	rows, err := r.db.Query("CALL usp_SpamBlock_getTime(?, @out_nRet)", charNo)
	if err != nil {
		return 0, 0, fmt.Errorf("error calling usp_SpamBlock_getTime(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var cNo, bTime, sTime int
		if err := rows.Scan(&cNo, &bTime, &sTime); err == nil {
			return bTime, sTime, nil
		}
	}
	return 0, 0, nil
}

// SetHolyPromise vincula o actualiza una promesa sagrada invocando p_HolyPromise_Set.
func (r *MySQLCharacterRepository) SetHolyPromise(charNo, upCharNo int, flags byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_HolyPromise_Set(30, ?, ?, ?, @out_nRet)", charNo, upCharNo, flags).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_HolyPromise_Set: %w", err)
	}
	return nil
}

// DelHolyPromise elimina una relación de promesa sagrada invocando p_HolyPromise_Del.
func (r *MySQLCharacterRepository) DelHolyPromise(charNo, upCharNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_HolyPromise_Del(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_HolyPromise_Del: %w", err)
	}
	return nil
}

// DelHolyPromiseChar remueve todas las promesas sagradas del personaje invocando p_HolyPromise_DelChar.
func (r *MySQLCharacterRepository) DelHolyPromiseChar(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_HolyPromise_DelChar(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_HolyPromise_DelChar: %w", err)
	}
	return nil
}

// GetHolyPromise consulta los datos de promesa sagrada invocando p_HolyPromise_Get.
func (r *MySQLCharacterRepository) GetHolyPromise(charNo int) (*HolyPromiseData, error) {
	rows, err := r.db.Query("CALL p_HolyPromise_Get(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_HolyPromise_Get(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		hp := &HolyPromiseData{CharNo: charNo}
		var regDate, rejoinDate, acceptDate sql.NullTime
		var lastLogin sql.NullTime
		var sID sql.NullString
		var level sql.NullInt64
		if err := rows.Scan(&hp.UpCharNo, &hp.Flags, &regDate, &rejoinDate, &acceptDate, &sID, &level, &lastLogin); err != nil {
			return nil, fmt.Errorf("error scanning p_HolyPromise_Get: %w", err)
		}
		return hp, nil
	}
	return nil, nil
}
