package database

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// GetCharacterSkills recupera el árbol de habilidades aprendidas invocando p_Skill_GetAll.
func (r *MySQLCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error) {
	rows, err := r.db.Query("CALL p_Skill_GetAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Skill_GetAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var skills []*CharacterSkill
	for rows.Next() {
		var skillNo uint16
		var skillLevel byte
		var exp int
		var writeTime int64
		var coolTime int
		var pd, ps, keepTime, ct byte
		if err := rows.Scan(&skillNo, &skillLevel, &exp, &writeTime, &coolTime, &pd, &ps, &keepTime, &ct); err == nil {
			skills = append(skills, &CharacterSkill{
				CharNo:     charNo,
				SkillNo:    skillNo,
				SkillLevel: skillLevel,
			})
		}
	}
	return skills, nil
}

// SaveCharacterSkill almacena o actualiza una habilidad invocando p_Skill_Set.
func (r *MySQLCharacterRepository) SaveCharacterSkill(charNo int, skill *CharacterSkill) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Skill_Set(?, ?, ?, 0, 0, 0, 0, 0, 0, 0, @out_nRet)",
		charNo, skill.SkillNo, skill.SkillLevel,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Skill_Set: %w", err)
	}
	return nil
}

// SaveCharacterSkills guarda un conjunto de habilidades invocando p_Skill_Set.
func (r *MySQLCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error {
	for _, sk := range skills {
		if err := r.SaveCharacterSkill(charNo, sk); err != nil {
			return err
		}
	}
	return nil
}

// SetSkillCoolTime registra el tiempo de recarga de una habilidad invocando p_Skill_SetCoolTime.
func (r *MySQLCharacterRepository) SetSkillCoolTime(charNo int, skillNo int16, writeTime int64, coolTime int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Skill_SetCoolTime(?, ?, ?, ?, @out_nRet)", charNo, skillNo, writeTime, coolTime).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Skill_SetCoolTime: %w", err)
	}
	return nil
}

// SetSkillPower asigna la potencia y refuerzo de una habilidad invocando p_Skill_SetPower.
func (r *MySQLCharacterRepository) SetSkillPower(charNo int, skillNo int16, pd, ps, keepTime, coolTime byte) error {
	var ret int
	_ = r.db.QueryRow("CALL p_Skill_SetPower(?, ?, 0, ?, @out_nRet)", charNo, skillNo, pd).Scan(&ret)
	_ = r.db.QueryRow("CALL p_Skill_SetPower(?, ?, 1, ?, @out_nRet)", charNo, skillNo, ps).Scan(&ret)
	_ = r.db.QueryRow("CALL p_Skill_SetPower(?, ?, 2, ?, @out_nRet)", charNo, skillNo, keepTime).Scan(&ret)
	_ = r.db.QueryRow("CALL p_Skill_SetPower(?, ?, 3, ?, @out_nRet)", charNo, skillNo, coolTime).Scan(&ret)
	return nil
}

// GetCharacterShortcuts lee los atajos serializados de la barra rápida invocando p_Char_GetOptShortCutData.
func (r *MySQLCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error) {
	var ret int
	var data []byte
	row := r.db.QueryRow("CALL p_Char_GetOptShortCutData(?, @out_nRet, @out_sData)", charNo)
	if err := row.Scan(&ret, &data); err != nil || ret != 1 || len(data) == 0 {
		return nil, nil
	}

	var shortcuts []*Shortcut
	entrySize := 8
	count := len(data) / entrySize
	for i := 0; i < count; i++ {
		offset := i * entrySize
		slot := byte(i)
		codeType := int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
		codeNo := int16(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))
		value := int32(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		if codeType != 0 || codeNo != 0 {
			shortcuts = append(shortcuts, &Shortcut{
				CharNo:   charNo,
				SlotNo:   slot,
				CodeType: codeType,
				CodeNo:   codeNo,
				Value:    value,
			})
		}
	}
	return shortcuts, nil
}

// SaveCharacterShortcuts serializa y persiste la barra de accesos directos invocando p_Char_SetOptShortCutData.
func (r *MySQLCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error {
	buf := make([]byte, 1024)
	for _, sc := range shortcuts {
		offset := int(sc.SlotNo) * 8
		if offset+8 <= len(buf) {
			binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(sc.CodeType))
			binary.LittleEndian.PutUint16(buf[offset+2:offset+4], uint16(sc.CodeNo))
			binary.LittleEndian.PutUint32(buf[offset+4:offset+8], uint32(sc.Value))
		}
	}
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetOptShortCutData(?, ?, @out_nRet)", charNo, buf).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetOptShortCutData: %w", err)
	}
	return nil
}

// GetCharacterOptions recupera todos los bloques binarios de configuración del cliente en tCharacterOptions.
func (r *MySQLCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error) {
	opt := &CharacterOptions{CharNo: charNo}

	var ret int
	var raw []byte

	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptGame(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.Game = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptShortCutData(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.ShortCutData = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptShortCutSize(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.ShortCutSize = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptSound(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.Sound = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptVideo(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.Video = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptWindowsPos(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.WindowsPos = raw
	}
	ret, raw = 0, nil
	if err := r.db.QueryRow("CALL p_Char_GetOptKeyMapping(?, @out_nRet, @out_sData)", charNo).Scan(&ret, &raw); err == nil && ret == 1 {
		opt.KeyMapping = raw
	}

	return opt, nil
}

// SaveCharacterOption persiste un bloque específico de opciones del cliente invocando p_Char_SetOpt*.
func (r *MySQLCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error {
	var ret int
	var err error
	switch strings.ToLower(optType) {
	case "game":
		err = r.db.QueryRow("CALL p_Char_SetOptGame(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "shortcut", "shortcutdata", "shortcut_data":
		err = r.db.QueryRow("CALL p_Char_SetOptShortCutData(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "shortcutsize", "shortcut_size":
		err = r.db.QueryRow("CALL p_Char_SetOptShortCutSize(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "sound":
		err = r.db.QueryRow("CALL p_Char_SetOptSound(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "video":
		err = r.db.QueryRow("CALL p_Char_SetOptVideo(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "windowspos", "window_pos", "windows_pos":
		err = r.db.QueryRow("CALL p_Char_SetOptWindowsPos(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	case "keymapping", "key_mapping", "key_map":
		err = r.db.QueryRow("CALL p_Char_SetOptKeyMapping(?, ?, @out_nRet)", charNo, data).Scan(&ret)
	default:
		return fmt.Errorf("unknown option type: %s", optType)
	}
	if err != nil {
		return fmt.Errorf("error persisting option %s: %w", optType, err)
	}
	return nil
}

// DeleteCharacterOptions borra el bloque completo de opciones clásicas invocando p_Char_DelOptions.
func (r *MySQLCharacterRepository) DeleteCharacterOptions(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_DelOptions(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_DelOptions(%d): %w", charNo, err)
	}
	return nil
}

// GetImproveShortCutList recupera la lista de atajos mejorados NA2016 invocando usp_Character_getShortCutList.
func (r *MySQLCharacterRepository) GetImproveShortCutList(charNo int) ([]*ImproveShortcutEntry, error) {
	rows, err := r.db.Query("CALL usp_Character_getShortCutList(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Character_getShortCutList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*ImproveShortcutEntry
	for rows.Next() {
		var slotNo byte
		var codeNo uint16
		var value int32
		if scanErr := rows.Scan(&slotNo, &codeNo, &value); scanErr == nil {
			list = append(list, &ImproveShortcutEntry{
				SlotNo: slotNo,
				CodeNo: codeNo,
				Value:  value,
			})
		}
	}
	return list, nil
}

// SetImproveShortCut persiste una ranura de atajo mejorado NA2016 invocando usp_Character_setShortCut.
func (r *MySQLCharacterRepository) SetImproveShortCut(charNo int, slotNo byte, codeNo uint16, value int32) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setShortCut(?, ?, ?, ?, @out_nRet)", charNo, slotNo, codeNo, value).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setShortCut: %w", err)
	}
	return nil
}

// GetImproveKeyMap recupera el mapeo de teclas mejorado NA2016 invocando usp_Character_getKeyMap.
func (r *MySQLCharacterRepository) GetImproveKeyMap(charNo int) ([]*ImproveKeyMapEntry, error) {
	rows, err := r.db.Query("CALL usp_Character_getKeyMap(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Character_getKeyMap(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*ImproveKeyMapEntry
	for rows.Next() {
		var funcNo uint16
		var extKey, ascii byte
		if scanErr := rows.Scan(&funcNo, &extKey, &ascii); scanErr == nil {
			list = append(list, &ImproveKeyMapEntry{
				FunctionNo: funcNo,
				ExtendKey:  extKey,
				ASCIICode:  ascii,
			})
		}
	}
	return list, nil
}

// SetImproveKeyMap persiste una entrada de mapeo de tecla mejorado NA2016 invocando usp_Character_setKeyMap.
func (r *MySQLCharacterRepository) SetImproveKeyMap(charNo int, funcNo uint16, extendKey, asciiCode byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setKeyMap(?, ?, ?, ?, @out_nRet)", charNo, funcNo, extendKey, asciiCode).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setKeyMap: %w", err)
	}
	return nil
}

// InitImproveKeyMap restablece el mapeo de teclas mejorado desde tKeyMapInit invocando usp_Character_initKeyMap.
func (r *MySQLCharacterRepository) InitImproveKeyMap(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_initKeyMap(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_initKeyMap(%d): %w", charNo, err)
	}
	return nil
}

// GetImproveGameOptions recupera las opciones de juego mejoradas NA2016 invocando usp_Character_getGameOption.
func (r *MySQLCharacterRepository) GetImproveGameOptions(charNo int) ([]*ImproveGameOptionEntry, error) {
	rows, err := r.db.Query("CALL usp_Character_getGameOption(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Character_getGameOption(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*ImproveGameOptionEntry
	for rows.Next() {
		var optNo uint16
		var val byte
		if scanErr := rows.Scan(&optNo, &val); scanErr == nil {
			list = append(list, &ImproveGameOptionEntry{
				OptionNo: optNo,
				Value:    val,
			})
		}
	}
	return list, nil
}

// SetImproveGameOption persiste una opción de juego mejorada NA2016 invocando usp_Character_setGameOption.
func (r *MySQLCharacterRepository) SetImproveGameOption(charNo int, optionNo uint16, value byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setGameOption(?, ?, ?, @out_nRet)", charNo, optionNo, value).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setGameOption: %w", err)
	}
	return nil
}

// InitImproveGameOptions restablece las opciones de juego mejoradas invocando usp_Character_setGameOptionInit.
func (r *MySQLCharacterRepository) InitImproveGameOptions(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Character_setGameOptionInit(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setGameOptionInit(%d): %w", charNo, err)
	}
	return nil
}

// GetInGameGuideViewList recupera los IDs de guías vistas invocando usp_InGameGuide_getViewList.
func (r *MySQLCharacterRepository) GetInGameGuideViewList(charNo int) ([]byte, error) {
	rows, err := r.db.Query("CALL usp_InGameGuide_getViewList(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_InGameGuide_getViewList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []byte
	for rows.Next() {
		var g int
		if scanErr := rows.Scan(&g); scanErr == nil {
			list = append(list, byte(g))
		}
	}
	return list, nil
}

// SetInGameGuideView registra una guía vista invocando usp_InGameGuide_setGuideNo.
func (r *MySQLCharacterRepository) SetInGameGuideView(charNo int, guideNo byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_InGameGuide_setGuideNo(?, ?, @out_nRet)", charNo, int(guideNo)).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_InGameGuide_setGuideNo: %w", err)
	}
	return nil
}

// GetCharacterQuests recupera todas las misiones registradas de un personaje invocando p_Quest_GetAll.
func (r *MySQLCharacterRepository) GetCharacterQuests(charNo int) ([]*QuestData, error) {
	rows, err := r.db.Query("CALL p_Quest_GetAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Quest_GetAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var quests []*QuestData
	for rows.Next() {
		q := &QuestData{CharNo: charNo}
		if err := rows.Scan(&q.QuestNo, &q.Status, &q.Data); err == nil {
			quests = append(quests, q)
		}
	}
	return quests, nil
}

// GetCharacterDoneQuests recupera los IDs de misiones completadas invocando p_Quest_GetAllDone (status = 2).
func (r *MySQLCharacterRepository) GetCharacterDoneQuests(charNo int) ([]int, error) {
	rows, err := r.db.Query("CALL p_Quest_GetAllDone(?, 2)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Quest_GetAllDone(%d): %w", charNo, err)
	}
	defer rows.Close()

	var done []int
	for rows.Next() {
		var qNo int
		if err := rows.Scan(&qNo); err == nil {
			done = append(done, qNo)
		}
	}
	return done, nil
}

// AddCharacterQuest registra una nueva misión activa en tQuest invocando p_Quest_Add.
func (r *MySQLCharacterRepository) AddCharacterQuest(charNo int, questNo int, status byte, data []byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Quest_Add(?, ?, ?, ?, @out_nRet)", charNo, questNo, status, data).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Quest_Add: %w", err)
	}
	return nil
}

// SetCharacterQuest actualiza el progreso o estado de una misión invocando p_Quest_Set.
func (r *MySQLCharacterRepository) SetCharacterQuest(charNo int, questNo int, status byte, data []byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Quest_Set(?, ?, ?, ?, @out_nRet)", charNo, questNo, status, data).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Quest_Set: %w", err)
	}
	return nil
}

// DeleteCharacterQuest cancela y remueve una misión activa invocando p_Quest_Del.
func (r *MySQLCharacterRepository) DeleteCharacterQuest(charNo int, questNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Quest_Del(?, ?, @out_nRet)", charNo, questNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Quest_Del: %w", err)
	}
	return nil
}

// GetAttendance obtiene la lista de fechas de asistencia del personaje invocando usp_Attendance_getCallDateList.
func (r *MySQLCharacterRepository) GetAttendance(charNo int) ([]time.Time, error) {
	rows, err := r.db.Query("CALL usp_Attendance_getCallDateList(?, '2000-01-01', '2099-12-31', @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Attendance_getCallDateList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err == nil {
			dates = append(dates, d)
		}
	}
	return dates, nil
}

// CheckAttendance registra la marca de asistencia diaria invocando usp_Attendance_setCall.
func (r *MySQLCharacterRepository) CheckAttendance(charNo int, dDate time.Time) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Attendance_setCall(?, ?, @out_nRet)", charNo, dDate).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Attendance_setCall: %w", err)
	}
	return nil
}

// GetAttendanceRewards consulta las recompensas de asistencia reclamadas invocando usp_Attendance_getRewardList.
func (r *MySQLCharacterRepository) GetAttendanceRewards(charNo int) ([]int, error) {
	rows, err := r.db.Query("CALL usp_Attendance_getRewardList(?, '2000-01-01', @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Attendance_getRewardList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var rewards []int
	for rows.Next() {
		var d time.Time
		var rID int
		if err := rows.Scan(&d, &rID); err == nil {
			rewards = append(rewards, rID)
		}
	}
	return rewards, nil
}

// SetAttendanceReward otorga y registra una recompensa de asistencia invocando usp_Attendance_setReward.
func (r *MySQLCharacterRepository) SetAttendanceReward(charNo int, rewardID int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Attendance_setReward(?, NOW(), ?, @out_nRet)", charNo, rewardID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Attendance_setReward: %w", err)
	}
	return nil
}

// SetAttendanceRewardAt registra una recompensa de asistencia para la fecha de evento indicada invocando usp_Attendance_setReward.
func (r *MySQLCharacterRepository) SetAttendanceRewardAt(charNo int, eventDate time.Time, rewardID int) (int, error) {
	var ret int
	err := r.db.QueryRow("CALL usp_Attendance_setReward(?, ?, ?, @out_nRet)", charNo, eventDate, rewardID).Scan(&ret)
	if err != nil {
		return -9999, fmt.Errorf("error calling usp_Attendance_setReward: %w", err)
	}
	return ret, nil
}

// GetCollectionCards recupera las cartas de colección del personaje invocando usp_Collection_getCardList.
func (r *MySQLCharacterRepository) GetCollectionCards(charNo int) ([]*CollectionCardData, error) {
	rows, err := r.db.Query("CALL usp_Collection_getCardList(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Collection_getCardList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var cards []*CollectionCardData
	for rows.Next() {
		var cardID int16
		var serialNo int64
		var starNumber byte
		var d time.Time
		if err := rows.Scan(&cardID, &serialNo, &starNumber, &d); err == nil {
			cards = append(cards, &CollectionCardData{
				CardID: cardID,
				Count:  int16(starNumber),
			})
		}
	}
	return cards, nil
}

// SetCollectionCard registra una carta de colección invocando usp_Collection_setCard.
func (r *MySQLCharacterRepository) SetCollectionCard(charNo int, cardID int16, count int16) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Collection_setCard(?, ?, 1, 1, 0, @out_nRet)", charNo, cardID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Collection_setCard: %w", err)
	}
	return nil
}

// GetTutorialCharInfo consulta el progreso del tutorial del personaje invocando usp_Tutorial_getCharInfo.
func (r *MySQLCharacterRepository) GetTutorialCharInfo(charNo int) (byte, byte, error) {
	rows, err := r.db.Query("CALL usp_Tutorial_getCharInfo(?, @out_nRet)", charNo)
	if err != nil {
		return 0, 0, fmt.Errorf("error calling usp_Tutorial_getCharInfo(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var state, step byte
		if err := rows.Scan(&state, &step); err == nil {
			return state, step, nil
		}
	}
	return 0, 0, nil
}

// SetTutorialCharInfo actualiza el progreso del tutorial invocando usp_Tutorial_setCharInfo.
func (r *MySQLCharacterRepository) SetTutorialCharInfo(userNo, charNo int, state, step byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Tutorial_setCharInfo(?, ?, ?, ?, @out_nRet)", userNo, charNo, state, step).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Tutorial_setCharInfo: %w", err)
	}
	return nil
}

// GetTutorialIsSkip consulta si la cuenta puede saltear el tutorial invocando usp_Tutorial_getIsSkip.
func (r *MySQLCharacterRepository) GetTutorialIsSkip(userNo int) (bool, error) {
	var ret int
	err := r.db.QueryRow("CALL usp_Tutorial_getIsSkip(?, @out_nRet)", userNo).Scan(&ret)
	if err != nil {
		return false, fmt.Errorf("error calling usp_Tutorial_getIsSkip(%d): %w", userNo, err)
	}
	return ret == 1, nil
}
