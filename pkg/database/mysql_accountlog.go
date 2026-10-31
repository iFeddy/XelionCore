package database

import (
	"database/sql"
	"fmt"
)

// IsIPBlocked consulta si una dirección IP está bloqueada en tIPBlock mediante p_IPBlock_IsBlock.
// Aplica whitelist para cuentas veteranas con logins exitosos previos.
func (r *MySQLAccountRepository) IsIPBlocked(userNo int, ip string) (bool, error) {
	var isBlock int
	err := r.db.QueryRow("CALL p_IPBlock_IsBlock(?, ?, @b)", userNo, ip).Scan(&isBlock)
	if err != nil {
		return false, fmt.Errorf("p_IPBlock_IsBlock failed: %w", err)
	}
	return isBlock == 1, nil
}

// LogAccount audita un evento de login o fallo en tAccountLog mediante p_AccountLog_Add.
// Cuando logType = 1 (Login exitoso), actualiza tUser.sUserIP preservando dDate e inserta en tUserLog.
func (r *MySQLAccountRepository) LogAccount(userID string, logType uint8, userNo int, playMin uint16, userPW string, worldNo uint8, errorCode uint16, ip string) error {
	var logNo int64
	err := r.db.QueryRow("CALL p_AccountLog_Add(?, ?, ?, ?, ?, ?, ?, ?, @out_log)",
		logType, userNo, worldNo, playMin, errorCode, ip, userID, userPW).Scan(&logNo)
	if err != nil {
		return fmt.Errorf("p_AccountLog_Add failed: %w", err)
	}
	return nil
}

// LogLogout registra la desconexión de un personaje en tLogout mediante usp_Log_setLogout.
func (r *MySQLAccountRepository) LogLogout(worldNo uint8, userNo, charNo int, classVal, level uint8, exp int64, playTime int16) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Log_setLogout(?, ?, ?, ?, ?, ?, ?, @r)",
		worldNo, userNo, charNo, classVal, level, exp, playTime).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_Log_setLogout failed: %w", err)
	}
	return nil
}

// GetDateDiffLastLogout consulta los segundos transcurridos desde el último logout en tLogout mediante usp_AccountLog_GetDateDiffLastlogout.
func (r *MySQLAccountRepository) GetDateDiffLastLogout(userNo int) (int64, error) {
	var diffSec int64
	err := r.db.QueryRow("CALL usp_AccountLog_GetDateDiffLastlogout(?, @diff)", userNo).Scan(&diffSec)
	if err != nil {
		return 0, fmt.Errorf("usp_AccountLog_GetDateDiffLastlogout failed: %w", err)
	}
	return diffSec, nil
}

// AddReport registra una denuncia GM en tReport mediante p_Report_Add validando cooldowns 1h/24h.
func (r *MySQLAccountRepository) AddReport(report *Report) (int, error) {
	var regNo int
	err := r.db.QueryRow("CALL p_Report_Add(?, ?, ?, ?, ?, ?, ?, @regNo)",
		report.Reporter, report.ReporterCharID, report.Harmer, report.WorldNo,
		report.ReportType, report.Reason, report.ChatLog).Scan(&regNo)
	if err != nil {
		return 0, fmt.Errorf("p_Report_Add failed: %w", err)
	}
	if regNo < 0 {
		return regNo, nil
	}
	return regNo, nil
}

// GetReports consulta las denuncias registradas por un jugador mediante p_Report_Get.
func (r *MySQLAccountRepository) GetReports(reporterUserNo int, worldNo uint8) ([]*ReportSummary, error) {
	rows, err := r.db.Query("CALL p_Report_Get(?, ?)", reporterUserNo, worldNo)
	if err != nil {
		return nil, fmt.Errorf("p_Report_Get failed: %w", err)
	}
	defer rows.Close()

	var list []*ReportSummary
	for rows.Next() {
		var regNo int
		var result uint8
		var remark sql.NullString
		if err := rows.Scan(&regNo, &result, &remark); err != nil {
			return nil, fmt.Errorf("error scanning report: %w", err)
		}
		s := &ReportSummary{
			RegNumber: regNo,
			Result:    result,
		}
		if remark.Valid {
			s.Remark = remark.String
		}
		list = append(list, s)
	}
	return list, nil
}

// DeleteReport cancela una denuncia dentro de la hora posterior a su creación mediante p_Report_Del.
func (r *MySQLAccountRepository) DeleteReport(regNumber int) (int16, error) {
	var ret int16
	err := r.db.QueryRow("CALL p_Report_Del(?, @r)", regNumber).Scan(&ret)
	if err != nil {
		return 0, fmt.Errorf("p_Report_Del failed: %w", err)
	}
	return ret, nil
}

// AddSameTimeChar registra usuarios concurrentes por mundo en tSameTimeChar mediante p_SameTimeChar_Add.
func (r *MySQLAccountRepository) AddSameTimeChar(worldNo uint8, count int) (int, error) {
	var idx int
	err := r.db.QueryRow("CALL p_SameTimeChar_Add(?, ?, @idx)", worldNo, count).Scan(&idx)
	if err != nil {
		return 0, fmt.Errorf("p_SameTimeChar_Add failed: %w", err)
	}
	return idx, nil
}

// AddSameTimeCharClass registra usuarios concurrentes por clase y nivel en tDayClassLevel mediante p_SameTimeCharClass_Add.
func (r *MySQLAccountRepository) AddSameTimeCharClass(worldNo uint8, classVal uint8, level uint16, count int) error {
	var idx int
	err := r.db.QueryRow("CALL p_SameTimeCharClass_Add(?, ?, ?, ?, @idx)", worldNo, classVal, level, count).Scan(&idx)
	if err != nil {
		return fmt.Errorf("p_SameTimeCharClass_Add failed: %w", err)
	}
	return nil
}

// EventL20Get consulta el estado del evento nivel 20 en tEvent_L20 mediante p_Event_L20_Get.
func (r *MySQLAccountRepository) EventL20Get(userNo int) (int, error) {
	var bRewarded int
	var nRet int
	err := r.db.QueryRow("CALL p_Event_L20_Get(?, @b, @r)", userNo).Scan(&bRewarded, &nRet)
	if err != nil {
		return 0, fmt.Errorf("p_Event_L20_Get failed: %w", err)
	}
	return bRewarded, nil
}

// EventL20Set registra la entrega del premio nivel 20 en tEvent_L20 mediante p_Event_L20_Set.
func (r *MySQLAccountRepository) EventL20Set(userNo int, charNo int, worldNo uint8, eventNo uint8, eventPoint int) error {
	var nRet int
	err := r.db.QueryRow("CALL p_Event_L20_Set(?, '', ?, ?, @r)", userNo, worldNo, charNo).Scan(&nRet)
	if err != nil {
		return fmt.Errorf("p_Event_L20_Set failed: %w", err)
	}
	return nil
}

// EventiPadSetTicket acumula boletos de sorteo según tiempo de juego mediante usp_EventiPad_setTicket.
func (r *MySQLAccountRepository) EventiPadSetTicket(userNo int) error {
	var nRet int
	err := r.db.QueryRow("CALL usp_EventiPad_setTicket(?, 60, @r)", userNo).Scan(&nRet)
	if err != nil {
		return fmt.Errorf("usp_EventiPad_setTicket failed: %w", err)
	}
	return nil
}

// AddHackScanLog audita alteraciones binarias o detección de trampas en tHackScanLog mediante p_HackScan_Add.
func (r *MySQLAccountRepository) AddHackScanLog(userNo int, worldNo uint8, charNo int, ip, charID, desc string) error {
	_, err := r.db.Exec("CALL p_HackScan_Add(?, ?, ?, 1, ?)", userNo, charNo, worldNo, desc)
	if err != nil {
		return fmt.Errorf("p_HackScan_Add failed: %w", err)
	}
	return nil
}

// AddSyslog registra eventos de infraestructura en tSyslog mediante p_Syslog_Add.
func (r *MySQLAccountRepository) AddSyslog(serverID uint8, logType uint8, desc string, ip string) error {
	_, err := r.db.Exec("CALL p_Syslog_Add(?, 0, ?, ?, ?, 0, 0)", logType, desc, ip, serverID)
	if err != nil {
		return fmt.Errorf("p_Syslog_Add failed: %w", err)
	}
	return nil
}
