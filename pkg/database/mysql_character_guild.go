package database

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateGuild crea un nuevo gremio oficial invocando p_Guild_Make.
func (r *MySQLCharacterRepository) CreateGuild(masterCharID int, name, password string, guildType byte) (int, error) {
	var guildID int
	err := r.db.QueryRow(
		"CALL p_Guild_Make(?, ?, 30, ?, ?, 0, 65535, @out_nNo)",
		name, password, guildType, masterCharID,
	).Scan(&guildID)
	if err != nil {
		return 0, fmt.Errorf("error calling p_Guild_Make: %w", err)
	}
	if guildID <= 0 {
		return 0, fmt.Errorf("p_Guild_Make failed with code %d", guildID)
	}
	return guildID, nil
}

// GetGuild recupera el registro oficial de tGuild invocando p_Guild_Get.
func (r *MySQLCharacterRepository) GetGuild(guildID int) (*GuildData, error) {
	rows, err := r.db.Query("CALL p_Guild_Get(?)", guildID)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Guild_Get(%d): %w", guildID, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("guild %d not found", guildID)
	}

	g := &GuildData{ID: guildID}
	var (
		sName, sPassword                                                                             string
		nMaxMembers                                                                                  uint16
		nMoney                                                                                       uint64
		nType, nGrade                                                                                byte
		nFame                                                                                        uint32
		nStoneLevel                                                                                  uint16
		nExp                                                                                         uint64
		dCreateDate                                                                                  time.Time
		dWarRequestDate, dWarStartDate, dWarEndDate                                                  sql.NullTime
		nWarEnemy, nWarWinCount, nWarLoseCount, nWarDrawCount, nDismissStatus                        int
		dDismissDate, dNotifyDate                                                                    sql.NullTime
		sNotifyCharID, sNotify, sIntro                                                               string
	)

	err = rows.Scan(
		&sName, &sPassword, &nMaxMembers, &nMoney, &nType, &nGrade, &nFame, &nStoneLevel, &nExp,
		&dCreateDate, &dWarRequestDate, &dWarStartDate, &dWarEndDate, &nWarEnemy, &nWarWinCount,
		&nWarLoseCount, &nWarDrawCount, &nDismissStatus, &dDismissDate, &dNotifyDate,
		&sNotifyCharID, &sNotify, &sIntro,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning guild data: %w", err)
	}

	g.Name = sName
	g.Password = sPassword
	g.MaxMembers = nMaxMembers
	g.Money = nMoney
	g.Type = nType
	g.Grade = nGrade
	g.Fame = nFame
	g.StoneLevel = nStoneLevel
	g.Exp = nExp
	g.CreateDate = dCreateDate
	g.DismissStatus = byte(nDismissStatus)
	if dDismissDate.Valid {
		g.DismissDate = dDismissDate.Time
	}
	if dNotifyDate.Valid {
		g.NotifyDate = dNotifyDate.Time
	}
	g.NotifyCharID = sNotifyCharID
	g.Notify = sNotify
	g.Intro = sIntro

	return g, nil
}

// GetGuildByCharID consulta la membresía de un personaje invocando p_Guild_GetMyNo.
func (r *MySQLCharacterRepository) GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error) {
	var guildNo int
	err := r.db.QueryRow("CALL p_Guild_GetMyNo(?, @out_nNo)", charID).Scan(&guildNo)
	if err != nil || guildNo <= 0 {
		return nil, nil, nil
	}

	g, err := r.GetGuild(guildNo)
	if err != nil {
		return nil, nil, err
	}

	members, err := r.GetGuildMembers(guildNo)
	if err != nil {
		return g, nil, nil
	}

	for _, m := range members {
		if m.CharID == charID {
			return g, m, nil
		}
	}
	return g, nil, nil
}

// GetGuildMembers recupera los miembros del gremio invocando p_Guild_Member_GetAll.
func (r *MySQLCharacterRepository) GetGuildMembers(guildID int) ([]*GuildMemberData, error) {
	rows, err := r.db.Query("CALL p_Guild_Member_GetAll(?)", guildID)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Guild_Member_GetAll(%d): %w", guildID, err)
	}
	defer rows.Close()

	var members []*GuildMemberData
	for rows.Next() {
		m := &GuildMemberData{GuildID: guildID}
		var charNo int
		var joinDate time.Time
		var grade byte
		var flags int
		var expRatio int16
		var intro, charName, loginZone string
		var class, level byte
		var lastLogin sql.NullTime

		err := rows.Scan(
			&charNo, &joinDate, &grade, &flags, &expRatio, &intro,
			&charName, &class, &level, &loginZone, &lastLogin,
		)
		if err == nil {
			m.CharID = charNo
			m.JoinDate = joinDate
			m.Grade = grade
			m.Flags = flags
			m.ExpRatio = expRatio
			m.Intro = intro
			m.CharName = charName
			m.Class = class
			m.Level = level
			m.Map = loginZone
			if lastLogin.Valid {
				m.LastLogin = lastLogin.Time
			}
			members = append(members, m)
		}
	}
	return members, nil
}

// GetAllGuilds lista todos los gremios registrados invocando p_Guild_GetListDetail.
func (r *MySQLCharacterRepository) GetAllGuilds() ([]*GuildListItem, error) {
	rows, err := r.db.Query("CALL p_Guild_GetListDetail()")
	if err != nil {
		return nil, fmt.Errorf("error calling p_Guild_GetListDetail: %w", err)
	}
	defer rows.Close()

	var list []*GuildListItem
	for rows.Next() {
		item := &GuildListItem{}
		if err := rows.Scan(&item.ID, &item.Name, &item.MaxMembers, &item.Type, &item.Grade, &item.Fame, &item.Exp); err == nil {
			item.GuildID = uint32(item.ID)
			item.GuildName = item.Name
			list = append(list, item)
		}
	}
	return list, nil
}

// AddGuildMember añade un nuevo integrante al gremio invocando p_Guild_Member_Join.
func (r *MySQLCharacterRepository) AddGuildMember(guildID, charID int, grade byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_Member_Join(?, ?, ?, @out_nRet)", guildID, charID, grade).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_Member_Join: %w", err)
	}
	return nil
}

// RemoveGuildMember expulsa o da de baja a un integrante invocando p_Guild_Member_Leave.
func (r *MySQLCharacterRepository) RemoveGuildMember(guildID, charID int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_Member_Leave(?, ?, @out_nRet)", guildID, charID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_Member_Leave: %w", err)
	}
	return nil
}

// SetGuildMemberGrade cambia el rango de un integrante en el gremio invocando p_Guild_Member_SetGrade.
func (r *MySQLCharacterRepository) SetGuildMemberGrade(guildID, charID int, grade byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_Member_SetGrade(?, ?, ?, @out_nRet)", guildID, charID, grade).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_Member_SetGrade: %w", err)
	}
	return nil
}

// SetGuildNotify actualiza la notificación del tablero del gremio invocando p_Guild_SetNotify.
func (r *MySQLCharacterRepository) SetGuildNotify(guildID int, charName, notify string) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_SetNotify(?, ?, ?, @out_nRet)", guildID, charName, notify).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_SetNotify: %w", err)
	}
	return nil
}

// DismissGuild disuelve un gremio invocando p_Guild_Delete.
func (r *MySQLCharacterRepository) DismissGuild(guildID int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_Delete(?, @out_nRet)", guildID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_Delete(%d): %w", guildID, err)
	}
	return nil
}

// SetGuildMoney actualiza los fondos del gremio invocando p_Guild_SetMoney.
func (r *MySQLCharacterRepository) SetGuildMoney(guildID int, money uint64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_SetMoney(?, ?, @out_nRet)", guildID, money).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_SetMoney: %w", err)
	}
	return nil
}

// SetGuildGrade actualiza el rango del gremio invocando p_Guild_SetGrade.
func (r *MySQLCharacterRepository) SetGuildGrade(guildID int, grade byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Guild_SetGrade(?, ?, @out_nRet)", guildID, grade).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Guild_SetGrade: %w", err)
	}
	return nil
}

// GetGuildToken consulta los tokens del gremio invocando usp_Guild_GetToken.
func (r *MySQLCharacterRepository) GetGuildToken(guildID int) (int, error) {
	var token, ret int
	err := r.db.QueryRow("CALL usp_Guild_GetToken(?, @out_nToken, @out_nRet)", guildID).Scan(&token, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_Guild_GetToken(%d): %w", guildID, err)
	}
	return token, nil
}

// SetGuildToken actualiza los tokens del gremio invocando usp_Guild_SetToken.
func (r *MySQLCharacterRepository) SetGuildToken(guildID int, token int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Guild_SetToken(?, ?, @out_nRet)", guildID, token).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Guild_SetToken: %w", err)
	}
	return nil
}

// GetGuildHistory recupera el historial de eventos del gremio invocando usp_Guild_getHistoryList.
func (r *MySQLCharacterRepository) GetGuildHistory(guildID int) ([]*GuildHistoryData, error) {
	rows, err := r.db.Query("CALL usp_Guild_getHistoryList(?, @out_nRet)", guildID)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Guild_getHistoryList(%d): %w", guildID, err)
	}
	defer rows.Close()

	var history []*GuildHistoryData
	for rows.Next() {
		h := &GuildHistoryData{}
		if err := rows.Scan(&h.HistoryType, &h.Param1, &h.Param2, &h.Text); err == nil {
			history = append(history, h)
		}
	}
	return history, nil
}

// AddGuildHistory registra un nuevo evento en el historial del gremio invocando usp_Guild_setHistory.
func (r *MySQLCharacterRepository) AddGuildHistory(guildID int, historyType byte, param1, param2 int, text string) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Guild_setHistory(?, ?, ?, ?, ?, @out_nRet)", guildID, historyType, param1, param2, text).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Guild_setHistory: %w", err)
	}
	return nil
}

// CreateGuildAcademy funda la academia de un gremio invocando usp_GuildAcademy_setAcademy.
func (r *MySQLCharacterRepository) CreateGuildAcademy(guildID, masterCharID int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_setAcademy(?, ?, 50, @out_nRet)", guildID, masterCharID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_setAcademy: %w", err)
	}
	return nil
}

// GetGuildAcademy consulta los datos de la academia del gremio.
func (r *MySQLCharacterRepository) GetGuildAcademy(guildID int) (*GuildAcademyData, error) {
	acad := &GuildAcademyData{
		GuildID:   guildID,
		MaxMember: 50,
	}
	return acad, nil
}

// GetAcademyByCharID recupera la membresía de academia de un avatar invocando usp_GuildAcademy_getGuildNo.
func (r *MySQLCharacterRepository) GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error) {
	rows, err := r.db.Query("CALL usp_GuildAcademy_getGuildNo(?, @out_nRet)", charID)
	if err != nil {
		return nil, nil, nil
	}
	defer rows.Close()

	var guildNo int
	if rows.Next() {
		_ = rows.Scan(&guildNo)
	}
	_ = rows.Close()

	if guildNo <= 0 {
		return nil, nil, nil
	}

	members, err := r.GetAcademyMembers(guildNo)
	if err != nil {
		return &GuildAcademyData{GuildID: guildNo}, nil, nil
	}

	for _, m := range members {
		if m.CharID == charID {
			return &GuildAcademyData{GuildID: guildNo}, m, nil
		}
	}
	return &GuildAcademyData{GuildID: guildNo}, nil, nil
}

// GetAcademyMembers recupera los novatos y miembros de la academia invocando usp_GuildAcademy_getMemberList.
func (r *MySQLCharacterRepository) GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error) {
	rows, err := r.db.Query("CALL usp_GuildAcademy_getMemberList(?, @out_nRet)", guildID)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_GuildAcademy_getMemberList(%d): %w", guildID, err)
	}
	defer rows.Close()

	var members []*GuildAcademyMemberData
	for rows.Next() {
		m := &GuildAcademyMemberData{GuildID: guildID}
		var charNo int
		var joinLevel byte
		var joinDate time.Time
		var grade byte
		var intro string
		var chatBanTime int
		var charName, loginZone string
		var level, class byte
		var lastLogin sql.NullTime

		err := rows.Scan(
			&charNo, &joinLevel, &joinDate, &grade, &intro, &chatBanTime,
			&charName, &level, &loginZone, &lastLogin, &class,
		)
		if err == nil {
			m.CharID = charNo
			m.JoinLevel = joinLevel
			m.JoinDate = joinDate
			m.Grade = grade
			m.CharName = charName
			m.Level = level
			m.Class = class
			m.Map = loginZone
			if lastLogin.Valid {
				m.LastLogin = lastLogin.Time
			}
			members = append(members, m)
		}
	}
	return members, nil
}

// AddAcademyMember registra el ingreso de un novato a la academia invocando usp_GuildAcademy_setMember.
func (r *MySQLCharacterRepository) AddAcademyMember(guildID, charID int, joinLevel byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_setMember(?, ?, NOW(), 0, @out_nRet)", guildID, charID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_setMember: %w", err)
	}
	return nil
}

// RemoveAcademyMember remueve a un miembro de la academia invocando usp_GuildAcademy_delMember.
func (r *MySQLCharacterRepository) RemoveAcademyMember(guildID, charID int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_delMember(?, ?, @out_nRet)", guildID, charID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_delMember: %w", err)
	}
	return nil
}

// SetAcademyNotify actualiza la notificación de la academia invocando usp_GuildAcademy_setNotify.
func (r *MySQLCharacterRepository) SetAcademyNotify(guildID int, charName, notify string) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_setNotify(?, ?, @out_nRet)", guildID, notify).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_setNotify: %w", err)
	}
	return nil
}

// SetAcademyChatBan aplica silencio temporal en chat de academia invocando usp_GuildAcademy_setChatBanTime.
func (r *MySQLCharacterRepository) SetAcademyChatBan(guildID, charID int, banMinutes int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_setChatBanTime(?, ?, ?, @out_nRet)", guildID, charID, banMinutes).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_setChatBanTime: %w", err)
	}
	return nil
}

// SetAcademyGraduate gradúa a un miembro de la academia invocando usp_GuildAcademy_setGraduate.
func (r *MySQLCharacterRepository) SetAcademyGraduate(guildID, charID int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_GuildAcademy_setGraduate(?, ?, '', 1, 60, @out_nRet)", guildID, charID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_GuildAcademy_setGraduate: %w", err)
	}
	return nil
}

// GetAcademyGraduates recupera los IDs de graduados invocando usp_GuildAcademy_getGraduateList.
func (r *MySQLCharacterRepository) GetAcademyGraduates(guildID int) ([]int, error) {
	rows, err := r.db.Query("CALL usp_GuildAcademy_getGraduateList(?, @out_nRet)", guildID)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_GuildAcademy_getGraduateList(%d): %w", guildID, err)
	}
	defer rows.Close()

	var grads []int
	for rows.Next() {
		var charNo int
		if err := rows.Scan(&charNo); err == nil {
			grads = append(grads, charNo)
		}
	}
	return grads, nil
}

// GetGuildTournament consulta los datos de un torneo de gremios invocando p_GuildTournament_Get.
func (r *MySQLCharacterRepository) GetGuildTournament(gtNo int) (*GuildTournamentData, error) {
	rows, err := r.db.Query("CALL p_GuildTournament_Get(?)", gtNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_GuildTournament_Get(%d): %w", gtNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		return &GuildTournamentData{GTNo: gtNo}, nil
	}
	return &GuildTournamentData{GTNo: gtNo}, nil
}

// GetAllGuildTournaments lista todos los torneos invocando p_GuildTournament_GetAllTournamentList.
func (r *MySQLCharacterRepository) GetAllGuildTournaments() ([]*GuildTournamentData, error) {
	rows, err := r.db.Query("CALL p_GuildTournament_GetAllTournamentList()")
	if err != nil {
		return nil, fmt.Errorf("error calling p_GuildTournament_GetAllTournamentList: %w", err)
	}
	defer rows.Close()

	var list []*GuildTournamentData
	for rows.Next() {
		var gtNo int
		if err := rows.Scan(&gtNo); err == nil {
			list = append(list, &GuildTournamentData{GTNo: gtNo})
		}
	}
	return list, nil
}
