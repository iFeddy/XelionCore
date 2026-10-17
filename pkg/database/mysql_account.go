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
	return &MySQLAccountRepository{db: db}
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
