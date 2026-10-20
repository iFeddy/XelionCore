package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// MySQLAccountRepository implementa AccountRepository interactuando directamente
// con las tablas oficiales tUser, tUserAuth, tUserBlock, tAccountLog y tIPBlock.
type MySQLAccountRepository struct {
	db *sql.DB
}

// NewMySQLAccountRepository crea una nueva instancia del repositorio MySQL oficial.
func NewMySQLAccountRepository(db *sql.DB) *MySQLAccountRepository {
	repo := &MySQLAccountRepository{db: db}
	repo.initAuxiliaryTables()
	return repo
}

// initAuxiliaryTables asegura la existencia de las tablas canónicas de Account y AccountLog.
func (r *MySQLAccountRepository) initAuxiliaryTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tReturnUser (
			nRowNo INT NOT NULL AUTO_INCREMENT,
			dDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			nUserNo INT NOT NULL,
			nItemNo INT NOT NULL DEFAULT 30001,
			PRIMARY KEY (nRowNo),
			KEY idx_treturnuser_userno_date (nUserNo, dDate)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tSerialKey (
			userNo INT NOT NULL,
			serialKey BIGINT NOT NULL,
			registerDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (userNo)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tCash (
			userNo INT NOT NULL,
			cashType TINYINT UNSIGNED NOT NULL DEFAULT 0,
			cash INT NOT NULL DEFAULT 0,
			status TINYINT UNSIGNED NOT NULL DEFAULT 1,
			PRIMARY KEY (userNo, cashType)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tCashLog (
			rowNo INT NOT NULL AUTO_INCREMENT,
			userNo INT NOT NULL,
			logType TINYINT UNSIGNED NOT NULL DEFAULT 1,
			cashType TINYINT UNSIGNED NOT NULL DEFAULT 0,
			cash INT NOT NULL,
			reference VARCHAR(50) NOT NULL DEFAULT '',
			validDate DATETIME DEFAULT NULL,
			issueDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (rowNo),
			KEY idx_tcashlog_userno (userNo)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tLogout (
			rowNo BIGINT NOT NULL AUTO_INCREMENT,
			logDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			worldNo TINYINT UNSIGNED NOT NULL DEFAULT 0,
			userNo INT NOT NULL,
			charNo INT NOT NULL DEFAULT 0,
			class TINYINT UNSIGNED NOT NULL DEFAULT 0,
			level TINYINT UNSIGNED NOT NULL DEFAULT 1,
			experience BIGINT NOT NULL DEFAULT 0,
			playTime SMALLINT NOT NULL DEFAULT 0,
			PRIMARY KEY (rowNo),
			KEY idx_tlogout_userno (userNo)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tReport (
			nRegNumber INT NOT NULL AUTO_INCREMENT,
			sReportType VARCHAR(16) NOT NULL DEFAULT '',
			nReporter INT NOT NULL,
			sReporterCharID VARCHAR(16) NOT NULL DEFAULT '',
			sHarmer VARCHAR(16) NOT NULL DEFAULT '',
			nWorldNo TINYINT UNSIGNED NOT NULL DEFAULT 0,
			sReason VARCHAR(512) NOT NULL DEFAULT '',
			sChatLog TEXT,
			nResult TINYINT UNSIGNED NOT NULL DEFAULT 0,
			sRemark VARCHAR(512) NOT NULL DEFAULT '',
			dReportDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (nRegNumber),
			KEY idx_treport_reporter (nReporter, nWorldNo, dReportDate),
			KEY idx_treport_harmer (sHarmer, dReportDate)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tSameTimeChar (
			nIndex INT NOT NULL AUTO_INCREMENT,
			dDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			nWorld TINYINT UNSIGNED NOT NULL DEFAULT 0,
			nCount INT NOT NULL DEFAULT 0,
			PRIMARY KEY (nIndex),
			KEY idx_tsametimechar_date (dDate)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tDayClassLevel (
			nIndex INT NOT NULL AUTO_INCREMENT,
			dDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			nWorld TINYINT UNSIGNED NOT NULL DEFAULT 0,
			nClass TINYINT UNSIGNED NOT NULL DEFAULT 0,
			nLevel SMALLINT UNSIGNED NOT NULL DEFAULT 1,
			nCount INT NOT NULL DEFAULT 0,
			PRIMARY KEY (nIndex),
			KEY idx_tdayclasslevel_date (dDate, nWorld)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS tChargeItem (
			rowNo INT NOT NULL AUTO_INCREMENT,
			userNo INT NOT NULL,
			orderNo INT NOT NULL DEFAULT 0,
			goodsNo INT NOT NULL,
			amount INT NOT NULL DEFAULT 1,
			registerDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			isDraw TINYINT UNSIGNED NOT NULL DEFAULT 0,
			drawDate DATETIME DEFAULT NULL,
			PRIMARY KEY (rowNo),
			KEY idx_tchargeitem_userno_draw (userNo, isDraw)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, q := range queries {
		_, _ = r.db.Exec(q)
	}
}

// GetAccountByUsername busca un usuario por su nombre de usuario en tUser.
func (r *MySQLAccountRepository) GetAccountByUsername(username string) (*Account, error) {
	username = strings.TrimSpace(username)
	query := `
		SELECT 
			u.nUserNo, u.sUserID, u.sUserPW, u.sUserName, 
			u.bIsBlock, u.bIsLoginAble, u.bIsDelete, u.nAuthID, 
			u.sUserIP, u.dDate,
			COALESCE(a.bIsLoginAble, 1) as authLoginAble
		FROM tUser u
		LEFT JOIN tUserAuth a ON u.nAuthID = a.nAuthID
		WHERE LOWER(u.sUserID) = LOWER(?)
		LIMIT 1;
	`

	var (
		userNo        int
		userID        string
		userPW        string
		userName      string
		isBlock       int
		isLoginAble   int
		isDelete      int
		authID        byte
		userIP        string
		dDate         time.Time
		authLoginAble int
	)

	err := r.db.QueryRow(query, username).Scan(
		&userNo, &userID, &userPW, &userName,
		&isBlock, &isLoginAble, &isDelete, &authID,
		&userIP, &dDate, &authLoginAble,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	// Comprobar si tiene algún bloqueo activo en tUserBlock
	if isBlock == 0 {
		var activeBlockCount int
		blockQuery := `
			SELECT COUNT(*) 
			FROM tUserBlock 
			WHERE nUserNo = ? AND bStatus = 1 
			  AND (dEndDate IS NULL OR dEndDate > NOW());
		`
		if err := r.db.QueryRow(blockQuery, userNo).Scan(&activeBlockCount); err == nil && activeBlockCount > 0 {
			isBlock = 1
		}
	}

	acc := &Account{
		ID:       userNo,
		Username: userID,
		Password: userPW,
		Admin:    authID,
		Banned:   isBlock != 0 || isDelete != 0 || isLoginAble == 0 || authLoginAble == 0,
		TUser: TUser{
			UserNo:     userNo,
			UserID:     userID,
			UserPW:     userPW,
			UserName:   userName,
			IsBlock:    isBlock,
			IsDelete:   isDelete,
			AuthID:     authID,
			UserIP:     userIP,
			CreateDate: dDate,
			LoginDate:  dDate,
		},
	}

	return acc, nil
}

// CreateAccount registra una nueva cuenta en tUser con contraseña cifrada en MD5.
func (r *MySQLAccountRepository) CreateAccount(username, password string) (*Account, error) {
	username = strings.TrimSpace(username)
	pwHash := password
	if len(password) != 32 {
		h := md5.Sum([]byte(password))
		pwHash = hex.EncodeToString(h[:])
	}

	query := `
		INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsLoginAble, bIsDelete, nAuthID, sUserIP, dDate)
		VALUES (?, ?, ?, 0, 1, 0, 1, '', NOW());
	`

	res, err := r.db.Exec(query, username, pwHash, username)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ErrAccountExists
		}
		return nil, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetAccountByUsername(username)
}

// IsIPBlocked verifica si la IP está registrada en tIPBlock.
func (r *MySQLAccountRepository) IsIPBlocked(ip string) (bool, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return false, nil
	}

	query := `
		SELECT COUNT(*) 
		FROM tIPBlock 
		WHERE (
			(INET_ATON(?) IS NOT NULL AND INET_ATON(sIPStart) IS NOT NULL AND INET_ATON(sIPEnd) IS NOT NULL
			 AND INET_ATON(?) >= INET_ATON(sIPStart) AND INET_ATON(?) <= INET_ATON(sIPEnd))
			OR (? >= sIPStart AND ? <= sIPEnd)
		);
	`
	var count int
	err := r.db.QueryRow(query, ip, ip, ip, ip, ip).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// LogAccount registra un intento de login en la tabla oficial tAccountLog.
// nType: 1 = Login Success, 2 = Login Fail
func (r *MySQLAccountRepository) LogAccount(userID string, logType uint8, userNo int, worldNo uint8, errorCode uint16, ip string) error {
	query := `
		INSERT INTO tAccountLog (sUserID, nType, nUserNo, nPlayMin, sUserPassword, nWorldNo, nErrorCode, sIP, dDate)
		VALUES (?, ?, ?, 0, '', ?, ?, ?, NOW());
	`
	_, err := r.db.Exec(query, userID, logType, userNo, worldNo, errorCode, ip)
	return err
}

// UpdateLoginSession actualiza la IP del cliente y la fecha del último login en tUser.
func (r *MySQLAccountRepository) UpdateLoginSession(userNo int, ip string) error {
	query := `
		UPDATE tUser 
		SET sUserIP = ?, dDate = NOW() 
		WHERE nUserNo = ?;
	`
	_, err := r.db.Exec(query, ip, userNo)
	return err
}

// GetUserDateDiffCreated implementa la lógica de usp_User_GetDatediffCreated:
// retorna los segundos transcurridos desde que se creó la cuenta en tUser.
func (r *MySQLAccountRepository) GetUserDateDiffCreated(userNo int) (int64, error) {
	query := `SELECT TIMESTAMPDIFF(SECOND, dDate, NOW()) FROM tUser WHERE nUserNo = ?;`
	var diff int64
	err := r.db.QueryRow(query, userNo).Scan(&diff)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrAccountNotFound
		}
		return -1, err
	}
	return diff, nil
}

// GetReturnUserDateDiffLastLog implementa la lógica de usp_ReturnUser_GetDatediffLastlog:
// retorna los segundos transcurridos desde el último registro de retorno en tReturnUser.
func (r *MySQLAccountRepository) GetReturnUserDateDiffLastLog(userNo int) (int64, error) {
	query := `SELECT TIMESTAMPDIFF(SECOND, dDate, NOW()) FROM tReturnUser WHERE nUserNo = ? ORDER BY dDate DESC LIMIT 1;`
	var diff int64
	err := r.db.QueryRow(query, userNo).Scan(&diff)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil // Sin registros previos
		}
		return -1, err
	}
	return diff, nil
}

// SetReturnUserItem implementa la lógica de usp_ReturnUser_SetUserItem:
// asigna el paquete de ítems de bienvenida para un jugador retornante en tReturnUser.
func (r *MySQLAccountRepository) SetReturnUserItem(userNo int, itemNo int) error {
	if itemNo <= 0 {
		itemNo = 30001 // Valor default oficial en usp_ReturnUser_SetUserItem
	}
	query := `INSERT INTO tReturnUser (dDate, nUserNo, nItemNo) VALUES (NOW(), ?, ?);`
	_, err := r.db.Exec(query, userNo, itemNo)
	return err
}

// AddCash implementa la lógica oficial de usp_Cash_setCash / usp_Event_addCash:
// acredita saldo Cash/Spark para el tipo indicado y registra la auditoría en tCashLog.
func (r *MySQLAccountRepository) AddCash(userNo int, cashType uint8, cash int, reference string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	upsertQuery := `
		INSERT INTO tCash (userNo, cashType, cash, status)
		VALUES (?, ?, ?, 1)
		ON DUPLICATE KEY UPDATE cash = cash + VALUES(cash);
	`
	if _, err := tx.Exec(upsertQuery, userNo, cashType, cash); err != nil {
		return err
	}

	logQuery := `
		INSERT INTO tCashLog (userNo, logType, cashType, cash, reference, issueDate)
		VALUES (?, 1, ?, ?, ?, NOW());
	`
	if _, err := tx.Exec(logQuery, userNo, cashType, cash, reference); err != nil {
		return err
	}

	return tx.Commit()
}

// GetCash implementa usp_Cash_getCash: consulta el saldo de moneda premium por tipo.
func (r *MySQLAccountRepository) GetCash(userNo int, cashType uint8) (int, error) {
	query := `SELECT COALESCE(cash, 0) FROM tCash WHERE userNo = ? AND cashType = ?;`
	var cash int
	err := r.db.QueryRow(query, userNo, cashType).Scan(&cash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return cash, nil
}

// SetSerialKey implementa usp_Item_setSerialKey: registra o canjea la clave serial de un usuario.
func (r *MySQLAccountRepository) SetSerialKey(userNo int, serialKey int64) error {
	query := `
		INSERT INTO tSerialKey (userNo, serialKey, registerDate)
		VALUES (?, ?, NOW())
		ON DUPLICATE KEY UPDATE serialKey = VALUES(serialKey), registerDate = NOW();
	`
	_, err := r.db.Exec(query, userNo, serialKey)
	return err
}

// LogLogout implementa usp_Log_setLogout: registra en tLogout las estadísticas de la sesión finalizada.
func (r *MySQLAccountRepository) LogLogout(worldNo uint8, userNo, charNo int, classVal, level uint8, exp int64, playTime int16) error {
	query := `
		INSERT INTO tLogout (worldNo, userNo, charNo, class, level, experience, playTime, logDate)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW());
	`
	_, err := r.db.Exec(query, worldNo, userNo, charNo, classVal, level, exp, playTime)
	return err
}

// GetDateDiffLastLogout implementa usp_AccountLog_GetDateDiffLastlogout:
// devuelve los segundos transcurridos desde el último logout registrado en tLogout.
func (r *MySQLAccountRepository) GetDateDiffLastLogout(userNo int) (int64, error) {
	query := `SELECT TIMESTAMPDIFF(SECOND, MAX(logDate), NOW()) FROM tLogout WHERE userNo = ?;`
	var diff sql.NullInt64
	err := r.db.QueryRow(query, userNo).Scan(&diff)
	if err != nil {
		return -1, err
	}
	if !diff.Valid {
		return -1, nil // Sin registros de logout previo
	}
	return diff.Int64, nil
}

// AddReport implementa p_Report_Add:
// Registra una denuncia con las validaciones oficiales de tiempo:
// - Retorna -1 si ya realizó un reporte en la última hora (cooldown 1h).
// - Retorna -2 si ya reportó al mismo personaje en las últimas 24 horas (cooldown 1d).
// - Retorna el nRegNumber generado si se guardó con éxito.
func (r *MySQLAccountRepository) AddReport(report *Report) (int, error) {
	if report == nil {
		return 0, errors.New("report cannot be nil")
	}

	// 1. Cooldown de 1 hora por reportero en el mismo mundo
	var recentReports int
	cooldown1hQuery := `
		SELECT COUNT(*) FROM tReport 
		WHERE nReporter = ? AND nWorldNo = ? AND dReportDate > DATE_SUB(NOW(), INTERVAL 1 HOUR);
	`
	if err := r.db.QueryRow(cooldown1hQuery, report.Reporter, report.WorldNo).Scan(&recentReports); err == nil && recentReports > 0 {
		return -1, nil
	}

	// 2. Cooldown de 24 horas contra el mismo infractor (sHarmer)
	var recentSameHarmer int
	cooldown24hQuery := `
		SELECT COUNT(*) FROM tReport 
		WHERE nReporter = ? AND nWorldNo = ? AND sHarmer = ? AND dReportDate > DATE_SUB(NOW(), INTERVAL 1 DAY);
	`
	if err := r.db.QueryRow(cooldown24hQuery, report.Reporter, report.WorldNo, report.Harmer).Scan(&recentSameHarmer); err == nil && recentSameHarmer > 0 {
		return -2, nil
	}

	// 3. Inserción oficial
	insertQuery := `
		INSERT INTO tReport (sReportType, nReporter, sReporterCharID, sHarmer, nWorldNo, sReason, sChatLog, nResult, sRemark, dReportDate)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, '', NOW());
	`
	res, err := r.db.Exec(insertQuery,
		report.ReportType, report.Reporter, report.ReporterCharID,
		report.Harmer, report.WorldNo, report.Reason, report.ChatLog,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	report.RegNumber = int(id)
	return int(id), nil
}

// GetReports implementa p_Report_Get: consulta las denuncias generadas en la última hora.
func (r *MySQLAccountRepository) GetReports(reporterUserNo int, worldNo uint8) ([]*ReportSummary, error) {
	query := `
		SELECT nRegNumber, nResult, sRemark 
		FROM tReport 
		WHERE nReporter = ? AND nWorldNo = ? AND dReportDate > DATE_SUB(NOW(), INTERVAL 1 HOUR);
	`
	rows, err := r.db.Query(query, reporterUserNo, worldNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*ReportSummary
	for rows.Next() {
		var s ReportSummary
		if err := rows.Scan(&s.RegNumber, &s.Result, &s.Remark); err != nil {
			return nil, err
		}
		summaries = append(summaries, &s)
	}
	return summaries, nil
}

// DeleteReport implementa p_Report_Del con los códigos de retorno oficiales:
// - Retorna 0: Éxito (reporte eliminado).
// - Retorna 1: Reporte no encontrado.
// - Retorna 2: Reporte con más de 1 hora de antigüedad (no se puede cancelar).
// - Retorna 3: Reporte ya procesado por un GM (nResult != 0).
func (r *MySQLAccountRepository) DeleteReport(regNumber int) (int16, error) {
	query := `
		SELECT TIMESTAMPDIFF(SECOND, dReportDate, NOW()), nResult 
		FROM tReport 
		WHERE nRegNumber = ?;
	`
	var ageSec int64
	var result uint8
	err := r.db.QueryRow(query, regNumber).Scan(&ageSec, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 1, nil // No existe
		}
		return 1, err
	}

	if ageSec > 3600 {
		return 2, nil // Mayor a 1 hora
	}
	if result != 0 {
		return 3, nil // Ya procesado
	}

	delQuery := `DELETE FROM tReport WHERE nRegNumber = ?;`
	_, err = r.db.Exec(delQuery, regNumber)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

// AddSameTimeChar implementa p_SameTimeChar_Add: inserta la métrica de concurrentes globales (CCU).
func (r *MySQLAccountRepository) AddSameTimeChar(worldNo uint8, count int) (int, error) {
	query := `INSERT INTO tSameTimeChar (nWorld, nCount, dDate) VALUES (?, ?, NOW());`
	res, err := r.db.Exec(query, worldNo, count)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// AddSameTimeCharClass implementa p_SameTimeCharClass_Add: guarda el desglose de concurrentes por clase y nivel.
func (r *MySQLAccountRepository) AddSameTimeCharClass(worldNo uint8, classVal uint8, level uint16, count int) error {
	query := `
		INSERT INTO tDayClassLevel (dDate, nWorld, nClass, nLevel, nCount)
		VALUES (NOW(), ?, ?, ?, ?);
	`
	_, err := r.db.Exec(query, worldNo, classVal, level, count)
	return err
}

// GetChargedItems implementa usp_Charge_ItemSelectList / usp_ChargeItem_GetList:
// consulta los ítems del inventario premium pendientes de retirar (isDraw = 0) en tChargeItem.
func (r *MySQLAccountRepository) GetChargedItems(userNo int) ([]*ChargedItem, error) {
	query := `
		SELECT rowNo, goodsNo, amount, registerDate
		FROM tChargeItem
		WHERE userNo = ? AND isDraw = 0
		ORDER BY rowNo ASC;
	`
	rows, err := r.db.Query(query, userNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ChargedItem
	for rows.Next() {
		var it ChargedItem
		it.UserNo = userNo
		if err := rows.Scan(&it.RowNo, &it.GoodsNo, &it.Amount, &it.RegisterDate); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}
