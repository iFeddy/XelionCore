package audit_test

import (
	"testing"
	"time"
)

// TestFase03_Control2A_ExhaustiveCatalog valida contra information_schema que
// el 100% de los Stored Procedures de Account, AccountLog, StatisticsData y LoginSession
// existen en la base de datos xelion_online sin omisiones.
func TestFase03_Control2A_ExhaustiveCatalog(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. Catálogo esperado de AccountLog (28 SPs: 11 ODBC + 17 restantes no-OPT)
	expectedAccountLog := []string{
		"p_accountlog_add",
		"p_ipblock_isblock",
		"p_ipblock_add",
		"p_ipblock_delete",
		"p_report_add",
		"p_report_del",
		"p_sametimecharclass_add",
		"p_event_friend_delrequest",
		"p_event_friend_sethidemsg",
		"usp_accountlog_getdatedifflastlogout",
		"p_event_l20_get",
		"p_event_l20_set",
		"p_accountlog_get",
		"p_accountlog_getcount",
		"p_accountlog_getlog",
		"p_accountlog_getlogcount",
		"p_event_friend_eventcomplete",
		"p_event_friend_get",
		"p_event_friend_getcount",
		"p_event_friend_giveup",
		"p_event_worldcup_enter",
		"p_ipblock_get",
		"p_ipblock_getlist",
		"p_report_get",
		"usp_eventipad_getlist",
		"usp_eventipad_getticket",
		"usp_eventipad_setticket",
		"usp_log_setlogout",
	}

	// 2. Catálogo esperado de StatisticsData (18 SPs)
	expectedStatistics := []string{
		"pfindccu",
		"p_dayclasslevelplay_add",
		"p_dayclasslevelplay_query",
		"p_dayclasslevel_add",
		"p_dayclasslevel_query",
		"p_dayconnectiontime",
		"p_dayuser",
		"p_dayuser_add",
		"p_dayuser_query",
		"p_sametimechar_add",
		"p_sametimechar_query",
		"p_sametimeworld",
		"p_useragereport",
		"p_userage_query",
		"p_usergendertest",
		"p_usergender_query",
		"p_userlocation_query",
		"p_getquesttimes",
	}

	// 3. Catálogo esperado de LoginSession (13 SPs)
	expectedLoginSession := []string{
		"usp_version_getclientkey",
		"usp_version_syncclientkey",
		"usp_worldserver_getlist",
		"usp_worldserver_getbyno",
		"usp_worldserver_updateccu",
		"usp_session_registerlogin",
		"usp_session_issueworldtoken",
		"usp_session_consumeworldtoken",
		"usp_session_createotp",
		"usp_session_consumeotp",
		"usp_session_closeandgetplaymin",
		"p_hackscan_add",
		"p_syslog_add",
	}

	// Consultar information_schema.ROUTINES en MySQL
	rows, err := db.Query(`
		SELECT LOWER(ROUTINE_NAME) 
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = 'xelion_online' AND ROUTINE_TYPE = 'PROCEDURE'
	`)
	if err != nil {
		t.Fatalf("Error consultando information_schema.ROUTINES: %v", err)
	}
	defer rows.Close()

	actualRoutines := make(map[string]bool)
	for rows.Next() {
		var rName string
		if err := rows.Scan(&rName); err != nil {
			t.Fatalf("Error escaneando ROUTINE_NAME: %v", err)
		}
		actualRoutines[rName] = true
	}

	// Validar AccountLog
	for _, proc := range expectedAccountLog {
		if !actualRoutines[proc] {
			t.Errorf("[AccountLog] Procedimiento faltante en MySQL: %s", proc)
		}
	}

	// Validar StatisticsData
	for _, proc := range expectedStatistics {
		if !actualRoutines[proc] {
			t.Errorf("[StatisticsData] Procedimiento faltante en MySQL: %s", proc)
		}
	}

	// Validar LoginSession
	for _, proc := range expectedLoginSession {
		if !actualRoutines[proc] {
			t.Errorf("[LoginSession] Procedimiento faltante en MySQL: %s", proc)
		}
	}

	totalRequired := len(expectedAccountLog) + len(expectedStatistics) + len(expectedLoginSession)
	t.Logf("Control 2A (Fase 03): Verificados %d/%d procedimientos (AccountLog=%d, Statistics=%d, LoginSession=%d)",
		totalRequired, totalRequired, len(expectedAccountLog), len(expectedStatistics), len(expectedLoginSession))
}

// TestFase03_AccountLog_Add verifica la lógica oficial de registro de log de cuenta:
// cuando nType = 1 (Login exitoso), actualiza tUser.sUserIP preservando dDate y
// registra el login diario en tUserLog.
func TestFase03_AccountLog_Add(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. Crear usuario de prueba con fecha de registro fija en el pasado
	pastDate := time.Now().Add(-48 * time.Hour).Truncate(time.Second)
	_, err := db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES ('test_actlog_user', 'pass_hash', 'ActLogTest', 0, 0, 1, '127.0.0.1', ?)
		ON DUPLICATE KEY UPDATE sUserIP='127.0.0.1', dDate=VALUES(dDate)
	`, pastDate)
	if err != nil {
		t.Fatalf("Error insertando usuario de prueba: %v", err)
	}
	var userNo int
	err = db.QueryRow("SELECT nUserNo FROM tUser WHERE sUserID = 'test_actlog_user'").Scan(&userNo)
	if err != nil {
		t.Fatalf("Error obteniendo nUserNo: %v", err)
	}

	// Limpiar datos previos
	db.Exec("DELETE FROM tAccountLog WHERE nUserNo = ?", userNo)
	db.Exec("DELETE FROM tUserLog WHERE userNo = ?", userNo)

	// 2. Ejecutar p_AccountLog_Add con nType = 1 (Login exitoso) desde IP '10.20.30.40'
	newIP := "10.20.30.40"
	var logNo int64
	err = db.QueryRow("CALL p_AccountLog_Add(1, ?, 1, 0, 0, ?, 'test_actlog_user', 'pass_hash', @out_log)", userNo, newIP).Scan(&logNo)
	if err != nil {
		t.Fatalf("Error ejecutando p_AccountLog_Add: %v", err)
	}
	if logNo <= 0 {
		t.Errorf("Se esperaba logNo > 0, se obtuvo %d", logNo)
	}

	// 3. Verificar que tAccountLog tiene el registro
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tAccountLog WHERE nUserNo = ? AND sIP = ? AND nType = 1", userNo, newIP).Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("No se encontro el registro en tAccountLog: count=%d, err=%v", count, err)
	}

	// 4. Verificar que tUser.sUserIP fue actualizado a la nueva IP
	var updatedIP string
	var currentDate time.Time
	err = db.QueryRow("SELECT sUserIP, dDate FROM tUser WHERE nUserNo = ?", userNo).Scan(&updatedIP, &currentDate)
	if err != nil {
		t.Fatalf("Error consultando tUser: %v", err)
	}
	if updatedIP != newIP {
		t.Errorf("sUserIP no fue actualizado: esperado=%s, obtenido=%s", newIP, updatedIP)
	}
	// 5. Verificar que tUser.dDate NO fue alterado (preserva la fecha de creación original)
	if currentDate.Sub(pastDate).Abs() > 5*time.Second {
		t.Errorf("dDate fue modificado indebidamente! Original: %v, Actual: %v", pastDate, currentDate)
	}

	// 6. Verificar que tUserLog tiene el registro diario
	var logCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tUserLog WHERE userNo = ? AND joinDate = CURDATE()", userNo).Scan(&logCount)
	if err != nil || logCount != 1 {
		t.Errorf("tUserLog no contiene el registro diario: count=%d, err=%v", logCount, err)
	}
}

// TestFase03_IPBlock_IsBlock verifica la regla oficial de firewall de IP:
// Bloqueo de rangos con bypass para cuentas veteranas (que ya tienen nType=1 en tAccountLog).
func TestFase03_IPBlock_IsBlock(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Crear usuarios: uno nuevo sin logins previos, y uno veterano con login exitoso previo
	db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES ('user_novice_ip', 'pwd', 'NoviceUser', 0, 0, 1, '127.0.0.1', NOW())
		ON DUPLICATE KEY UPDATE sUserIP='127.0.0.1'
	`)
	db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES ('user_veteran_ip', 'pwd', 'VeteranUser', 0, 0, 1, '127.0.0.1', NOW())
		ON DUPLICATE KEY UPDATE sUserIP='127.0.0.1'
	`)

	var noviceNo, veteranNo int
	db.QueryRow("SELECT nUserNo FROM tUser WHERE sUserID = 'user_novice_ip'").Scan(&noviceNo)
	db.QueryRow("SELECT nUserNo FROM tUser WHERE sUserID = 'user_veteran_ip'").Scan(&veteranNo)

	// Limpiar logs y registrar login previo SOLO para el veterano
	db.Exec("DELETE FROM tAccountLog WHERE nUserNo IN (?, ?)", noviceNo, veteranNo)
	_, err := db.Exec("INSERT INTO tAccountLog (sUserID, nUserNo, nWorldNo, nType, sIP, dDate) VALUES ('user_veteran_ip', ?, 1, 1, '200.10.10.5', NOW() - INTERVAL 1 DAY)", veteranNo)
	if err != nil {
		t.Fatalf("Error insertando tAccountLog previo: %v", err)
	}

	// Agregar un rango de bloqueo en tIPBlock: 198.51.100.1 a 198.51.100.254
	db.Exec("DELETE FROM tIPBlock WHERE sIPStart = '198.51.100.1'")
	var addRet int
	err = db.QueryRow("CALL p_IPBlock_Add('198.51.100.1', '198.51.100.254', 'Test Block Range', @out_ret)").Scan(&addRet)
	if err != nil {
		t.Fatalf("Error en p_IPBlock_Add: %v", err)
	}

	testIP := "198.51.100.50"

	// 1. Novato intentando loguear desde la IP bloqueada -> DEBE SER BLOQUEADO (bIsBlock = 1)
	var isBlockNovice int
	err = db.QueryRow("CALL p_IPBlock_IsBlock(?, ?, @out_block)", noviceNo, testIP).Scan(&isBlockNovice)
	if err != nil {
		t.Fatalf("Error en p_IPBlock_IsBlock (novice): %v", err)
	}
	if isBlockNovice != 1 {
		t.Errorf("Se esperaba bIsBlock=1 para usuario novato en IP bloqueada, se obtuvo %d", isBlockNovice)
	}

	// 2. Veterano intentando loguear desde la misma IP bloqueada -> BYPASS AUTORIZADO (bIsBlock = 0)
	var isBlockVeteran int
	err = db.QueryRow("CALL p_IPBlock_IsBlock(?, ?, @out_block)", veteranNo, testIP).Scan(&isBlockVeteran)
	if err != nil {
		t.Fatalf("Error en p_IPBlock_IsBlock (veteran): %v", err)
	}
	if isBlockVeteran != 0 {
		t.Errorf("Se esperaba bIsBlock=0 (whitelist bypass) para usuario veterano, se obtuvo %d", isBlockVeteran)
	}

	// 3. Novato desde IP no bloqueada -> NO BLOQUEADO (bIsBlock = 0)
	var isBlockFree int
	err = db.QueryRow("CALL p_IPBlock_IsBlock(?, '8.8.8.8', @out_block)", noviceNo).Scan(&isBlockFree)
	if err != nil || isBlockFree != 0 {
		t.Errorf("Se esperaba bIsBlock=0 en IP libre, se obtuvo %d", isBlockFree)
	}

	// 4. Limpiar bloqueo con p_IPBlock_Delete
	var delRet int
	err = db.QueryRow("CALL p_IPBlock_Delete('198.51.100.1', @out_del)").Scan(&delRet)
	if err != nil {
		t.Fatalf("Error en p_IPBlock_Delete: %v", err)
	}
}

// TestFase03_Report_Add_Del valida los cooldowns oficiales de denuncias de jugadores (tReport):
// - Cooldown de 1 hora para el mismo reportero en el mismo mundo (-1).
// - Reglas de cancelación con p_Report_Del.
func TestFase03_Report_Add_Del(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	reporterNo := 9991
	worldNo := 1
	reporterChar := "ReporterChar"
	harmer1 := "CheaterBoy1"
	harmer2 := "CheaterBoy2"

	db.Exec("DELETE FROM tReport WHERE nReporter = ?", reporterNo)

	// 1. Primera denuncia válida -> devuelve nReportRegNo > 0
	var reportRegNo int
	err := db.QueryRow("CALL p_Report_Add(?, ?, ?, ?, 'BOT', 'Uso de bot', 'Chat log evidence...', @out_reg)",
		reporterNo, reporterChar, harmer1, worldNo).Scan(&reportRegNo)
	if err != nil || reportRegNo <= 0 {
		t.Fatalf("Primera denuncia debió ser aceptada (>0), obtenido: %d, err: %v", reportRegNo, err)
	}

	// 2. Inmediatamente intentar reportar en el mismo mundo a otro infractor -> -1 (cooldown 1h por mundo)
	var res2 int
	err = db.QueryRow("CALL p_Report_Add(?, ?, ?, ?, 'ABUSE', 'Insultos', 'Chat log 2...', @out_reg)",
		reporterNo, reporterChar, harmer2, worldNo).Scan(&res2)
	if err != nil || res2 != -1 {
		t.Errorf("Se esperaba -1 por cooldown de 1h en el mismo mundo, obtenido: %d, err: %v", res2, err)
	}

	// 3. Cancelar reporte recién creado -> 0 (éxito dentro de la hora)
	var delRet int
	err = db.QueryRow("CALL p_Report_Del(?, @out_del)", reportRegNo).Scan(&delRet)
	if err != nil || delRet != 0 {
		t.Errorf("Se esperaba 0 al cancelar reporte recién creado, obtenido: %d, err: %v", delRet, err)
	}

	// 4. Intentar cancelar reporte que ya no existe -> 1
	err = db.QueryRow("CALL p_Report_Del(?, @out_del)", reportRegNo).Scan(&delRet)
	if err != nil || delRet != 1 {
		t.Errorf("Se esperaba 1 al cancelar reporte inexistente, obtenido: %d", delRet)
	}
}

// TestFase03_Session_Lifecycle valida el ciclo completo de vida de sesión, token de mundo
// y OTP en tActiveAccountSession sin queries inline en Go.
func TestFase03_Session_Lifecycle(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Crear usuario de prueba
	db.Exec(`
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES ('session_test_user', 'pwd', 'SessionTest', 0, 0, 1, '127.0.0.1', NOW())
		ON DUPLICATE KEY UPDATE sUserIP='127.0.0.1'
	`)
	var userNo int
	db.QueryRow("SELECT nUserNo FROM tUser WHERE sUserID = 'session_test_user'").Scan(&userNo)

	// 1. usp_Session_RegisterLogin
	clientIP := "192.168.1.100"
	var regResult int
	err := db.QueryRow("CALL usp_Session_RegisterLogin(?, 'session_test_user', ?)", userNo, clientIP).Scan(&regResult)
	if err != nil || regResult != 0 {
		t.Fatalf("Error en usp_Session_RegisterLogin: res=%d, err=%v", regResult, err)
	}

	// 2. usp_Session_IssueWorldToken (Generar token binario de 64 bytes para World 1)
	rawToken := make([]byte, 64)
	for i := range rawToken {
		rawToken[i] = byte(i + 1)
	}
	var issueResult int
	err = db.QueryRow("CALL usp_Session_IssueWorldToken(?, 1, ?)", userNo, rawToken).Scan(&issueResult)
	if err != nil || issueResult != 0 {
		t.Fatalf("Error en usp_Session_IssueWorldToken: res=%d, err=%v", issueResult, err)
	}

	// 3. usp_Session_ConsumeWorldToken (WorldServer valida y consume el token de 64 bytes)
	var retVal, consumedUserNo int
	var consumedUserID string
	err = db.QueryRow("CALL usp_Session_ConsumeWorldToken(1, ?, ?)", rawToken, clientIP).Scan(&retVal, &consumedUserNo, &consumedUserID)
	if err != nil {
		t.Fatalf("Error en usp_Session_ConsumeWorldToken: %v", err)
	}
	if retVal != 0 || consumedUserNo != userNo || consumedUserID != "session_test_user" {
		t.Errorf("Token consumido devolvió valores incorrectos: ret=%d, userNo=%d, userID=%s", retVal, consumedUserNo, consumedUserID)
	}

	// 4. Intentar reusar el mismo token -> DEBE FALLAR (single-use, devuelve ret = -1)
	var reuseRet, reuseUserNo int
	var reuseUserID string
	err = db.QueryRow("CALL usp_Session_ConsumeWorldToken(1, ?, ?)", rawToken, clientIP).Scan(&reuseRet, &reuseUserNo, &reuseUserID)
	if err != nil {
		t.Fatalf("Error consultando reutilizacion de token: %v", err)
	}
	if reuseRet == 0 {
		t.Errorf("Se esperaba que el token ya no fuese válido tras ser consumido, pero devolvió ret=%d", reuseRet)
	}

	// 5. Test de OTP (32 caracteres)
	otp := "A1B2C3D4E5F678901234567890ABCDEF"
	var otpCreateRes int
	err = db.QueryRow("CALL usp_Session_CreateOTP(?, ?)", userNo, otp).Scan(&otpCreateRes)
	if err != nil || otpCreateRes != 0 {
		t.Fatalf("Error en usp_Session_CreateOTP: res=%d, err=%v", otpCreateRes, err)
	}

	var otpRet, consumedOTPUserNo int
	var consumedOTPUserID string
	err = db.QueryRow("CALL usp_Session_ConsumeOTP(?, ?)", otp, clientIP).Scan(&otpRet, &consumedOTPUserNo, &consumedOTPUserID)
	if err != nil || otpRet != 0 || consumedOTPUserNo != userNo {
		t.Errorf("Error consumiendo OTP: ret=%d, esperado=%d, obtenido=%d, err=%v", otpRet, userNo, consumedOTPUserNo, err)
	}

	// OTP consumido ya no existe
	var replayRet, replayOTPUserNo int
	var replayOTPUserID string
	_ = db.QueryRow("CALL usp_Session_ConsumeOTP(?, ?)", otp, clientIP).Scan(&replayRet, &replayOTPUserNo, &replayOTPUserID)
	if replayRet == 0 {
		t.Errorf("OTP consumido no debe ser reutilizable, obtenido: ret=%d", replayRet)
	}

	// 6. usp_Session_CloseAndGetPlayMin
	var playMin int
	err = db.QueryRow("CALL usp_Session_CloseAndGetPlayMin(?)", userNo).Scan(&playMin)
	if err != nil {
		t.Fatalf("Error en usp_Session_CloseAndGetPlayMin: %v", err)
	}
	if playMin < 0 {
		t.Errorf("playMin no puede ser negativo: %d", playMin)
	}
}

// TestFase03_Event_L20 valida el registro de reclamo de recompensa por nivel 20 (tEvent_L20).
func TestFase03_Event_L20(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userNo := 8881
	userID := "EventL20Hero"

	db.Exec("DELETE FROM tEvent_L20 WHERE nUserNo = ?", userNo)

	// 1. Inicialmente no tiene premio (out_nRet = -1)
	var bRewarded, nRet int
	err := db.QueryRow("CALL p_Event_L20_Get(?, @bRewarded, @nRet)", userNo).Scan(&nRet, &bRewarded)
	if err != nil || nRet != -1 {
		t.Errorf("Se esperaba nRet=-1 para usuario sin premio, obtenido nRet=%d, err=%v", nRet, err)
	}

	// 2. Otorgar premio
	var setResult int
	err = db.QueryRow("CALL p_Event_L20_Set(?, ?, 1, 1001, @out_set)", userNo, userID).Scan(&setResult)
	if err != nil || setResult != 0 {
		t.Fatalf("Error en p_Event_L20_Set: res=%d, err=%v", setResult, err)
	}

	// 3. Ahora debe tener premio (bRewarded = 1, nRet = 0)
	err = db.QueryRow("CALL p_Event_L20_Get(?, @bRewarded, @nRet)", userNo).Scan(&nRet, &bRewarded)
	if err != nil || bRewarded != 1 || nRet != 0 {
		t.Errorf("Se esperaba bRewarded=1, nRet=0 tras otorgar premio, obtenido: bRewarded=%d, nRet=%d", bRewarded, nRet)
	}
}

// TestFase03_AccountLog_GetDateDiffLastlogout valida el cálculo de diferencia de segundos desde el último logout.
func TestFase03_AccountLog_GetDateDiffLastlogout(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userNo := 7771
	db.Exec("DELETE FROM tLogout WHERE userNo = ?", userNo)

	// Insertar logout de hace 5 horas
	pastLogout := time.Now().Add(-5 * time.Hour)
	_, err := db.Exec("INSERT INTO tLogout (userNo, logDate, playTime) VALUES (?, ?, 120)", userNo, pastLogout)
	if err != nil {
		t.Fatalf("Error insertando tLogout: %v", err)
	}

	var diffSec int
	err = db.QueryRow("CALL usp_AccountLog_GetDateDiffLastlogout(?, @out_sec)", userNo).Scan(&diffSec)
	if err != nil {
		t.Fatalf("Error en usp_AccountLog_GetDateDiffLastlogout: %v", err)
	}
	expectedSec := 5 * 3600
	if diffSec < expectedSec-60 || diffSec > expectedSec+60 {
		t.Errorf("Se esperaba diffSec alrededor de %d, se obtuvo %d", expectedSec, diffSec)
	}
}

// TestFase03_WorldServer_And_Version valida los stored procedures operacionales de configuración.
func TestFase03_WorldServer_And_Version(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Guardar clave original de FOUS para restaurarla al finalizar el test
	var originalKey string
	_ = db.QueryRow("SELECT GameVersionKey FROM FOUS LIMIT 1").Scan(&originalKey)
	defer func() {
		if originalKey != "" {
			_, _ = db.Exec("UPDATE FOUS SET GameVersionKey = ?", originalKey)
		}
	}()

	// 1. Sincronizar versión de cliente
	keyData := "XelionClientVersionHashKey1234567890ABCDEF"
	var syncKey string
	err := db.QueryRow("CALL usp_Version_SyncClientKey(?)", keyData).Scan(&syncKey)
	if err != nil || syncKey != keyData {
		t.Fatalf("Error en usp_Version_SyncClientKey: res=%s, err=%v", syncKey, err)
	}

	// 2. Obtener clave de versión
	var fetchedKey string
	err = db.QueryRow("CALL usp_Version_GetClientKey()").Scan(&fetchedKey)
	if err != nil {
		t.Fatalf("Error en usp_Version_GetClientKey: %v", err)
	}
	if fetchedKey != keyData {
		t.Errorf("Clave de versión discordante: esperado=%s, obtenido=%s", keyData, fetchedKey)
	}

	// 3. Actualizar CCU de un mundo
	worldNo := 1
	var ccuRes int
	err = db.QueryRow("CALL usp_WorldServer_UpdateCCU(?, 150)", worldNo).Scan(&ccuRes)
	if err != nil || ccuRes != 0 {
		t.Fatalf("Error en usp_WorldServer_UpdateCCU: res=%d, err=%v", ccuRes, err)
	}

	// 4. Obtener mundo por número
	var wNo int
	var wName, wIP string
	var wPort, wMaxUser, wRatable, wMaint, wCCU int
	var wHeartbeat time.Time
	var wStatus int
	err = db.QueryRow("CALL usp_WorldServer_GetByNo(?)", worldNo).Scan(
		&wNo, &wName, &wIP, &wPort, &wMaxUser, &wRatable, &wMaint, &wCCU, &wHeartbeat, &wStatus,
	)
	if err != nil {
		t.Fatalf("Error en usp_WorldServer_GetByNo: %v", err)
	}
	if wCCU != 150 {
		t.Errorf("CCU no actualizado: esperado=150, obtenido=%d", wCCU)
	}
}

// TestFase03_HackScan_And_Syslog valida el registro de auditoría de seguridad y logs del sistema.
func TestFase03_HackScan_And_Syslog(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. p_HackScan_Add
	userNo := 6661
	var logNo int
	err := db.QueryRow("CALL p_HackScan_Add(?, 1001, 1, 101, 'SpeedHack.exe detectado')", userNo).Scan(&logNo)
	if err != nil || logNo <= 0 {
		t.Fatalf("Error en p_HackScan_Add: logNo=%d, err=%v", logNo, err)
	}

	var hackCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tHackScanLog WHERE nUserNo = ? AND nHackType = 101", userNo).Scan(&hackCount)
	if err != nil || hackCount < 1 {
		t.Errorf("No se encontro el registro de hack scan: count=%d", hackCount)
	}

	// 2. p_Syslog_Add
	var syslogNo int64
	err = db.QueryRow("CALL p_Syslog_Add(1001, ?, 'session_test_user', '192.168.1.100', 1, 50, 0)", userNo).Scan(&syslogNo)
	if err != nil || syslogNo <= 0 {
		t.Fatalf("Error en p_Syslog_Add: res=%d, err=%v", syslogNo, err)
	}
}

// TestFase03_Statistics_Procedures valida procedimientos de métricas y estadísticas.
func TestFase03_Statistics_Procedures(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. p_SameTimeChar_Add
	var idx int
	err := db.QueryRow("CALL p_SameTimeChar_Add(1, 250, @out_idx)").Scan(&idx)
	if err != nil {
		t.Fatalf("Error en p_SameTimeChar_Add: %v", err)
	}

	// 2. pFindCCU
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	rows, err := db.Query("CALL pFindCCU(?, ?)", nowStr, nowStr)
	if err != nil {
		t.Fatalf("Error en pFindCCU: %v", err)
	}
	rows.Close()

	// 3. p_SameTimeWorld
	rows, err = db.Query("CALL p_SameTimeWorld()")
	if err != nil {
		t.Fatalf("Error en p_SameTimeWorld: %v", err)
	}
	rows.Close()
}
