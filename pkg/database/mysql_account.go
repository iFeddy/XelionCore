package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// MySQLAccountRepository implementa AccountRepository, AccountShopRepository,
// AccountLogRepository y LoginSessionRepository invocando 100% Stored Procedures oficiales.
type MySQLAccountRepository struct {
	db *sql.DB
}

// NewMySQLAccountRepository crea una nueva instancia del repositorio MySQL oficial.
func NewMySQLAccountRepository(db *sql.DB) *MySQLAccountRepository {
	return &MySQLAccountRepository{db: db}
}

// LoginGameUS autentica una cuenta mediante usp_User_USloginGame (0x0C5A).
// Soporta login por sUserID, fallback por sEmail y alta automática con prefijo 'r_'.
func (r *MySQLAccountRepository) LoginGameUS(userID, userPW, ip string) (*Account, error) {
	pwHash := strings.ToLower(userPW)
	if len(pwHash) != 32 {
		h := md5.Sum([]byte(userPW))
		pwHash = hex.EncodeToString(h[:])
	}

	var nResult, outUserNo, outAuthID int
	err := r.db.QueryRow("CALL usp_User_USloginGame(?, ?, ?, @u, @a, @r)", userID, pwHash, ip).
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		return nil, fmt.Errorf("usp_User_USloginGame failed: %w", err)
	}

	switch nResult {
	case 0:
		acc := &Account{
			ID:       outUserNo,
			Username: strings.TrimPrefix(userID, "r_"),
			Password: pwHash,
			Admin:    byte(outAuthID),
			Banned:   false,
		}
		acc.Sync()
		return acc, nil
	case -1:
		return nil, ErrAccountNotFound
	case -2:
		return nil, ErrInvalidPassword
	case -3:
		return nil, ErrAccountBlocked
	default:
		return nil, fmt.Errorf("login failed with code %d", nResult)
	}
}

// LoginGame autentica una cuenta estándar mediante usp_User_loginGame (0x0C06).
func (r *MySQLAccountRepository) LoginGame(userID, userPW, ip string) (*Account, error) {
	pwHash := strings.ToLower(userPW)
	if len(pwHash) != 32 {
		h := md5.Sum([]byte(userPW))
		pwHash = hex.EncodeToString(h[:])
	}

	var nResult, outUserNo, outAuthID int
	err := r.db.QueryRow("CALL usp_User_loginGame(?, ?, ?, @u, @a, @r)", userID, pwHash, ip).
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		return nil, fmt.Errorf("usp_User_loginGame failed: %w", err)
	}

	switch nResult {
	case 0:
		acc := &Account{
			ID:       outUserNo,
			Username: userID,
			Password: pwHash,
			Admin:    byte(outAuthID),
			Banned:   false,
		}
		acc.Sync()
		return acc, nil
	case -1:
		return nil, ErrAccountNotFound
	case -2:
		return nil, ErrInvalidPassword
	case -3:
		return nil, ErrAccountBlocked
	default:
		return nil, fmt.Errorf("login failed with code %d", nResult)
	}
}

// LoginNetmarble autentica una cuenta de plataforma Netmarble mediante usp_User_LoginNetmarble (0x0C1D).
func (r *MySQLAccountRepository) LoginNetmarble(uniID, userID string) (*Account, error) {
	var nResult, outUserNo, outAuthID int
	err := r.db.QueryRow("CALL usp_User_LoginNetmarble(?, ?, @u, @a, @r)", uniID, userID).
		Scan(&nResult, &outUserNo, &outAuthID)
	if err != nil {
		return nil, fmt.Errorf("usp_User_LoginNetmarble failed: %w", err)
	}

	if nResult != 0 {
		return nil, fmt.Errorf("netmarble login failed with code %d", nResult)
	}

	acc := &Account{
		ID:       outUserNo,
		Username: userID,
		Admin:    byte(outAuthID),
		Banned:   false,
	}
	acc.Sync()
	return acc, nil
}

// CheckBlock consulta el estado de sanción/bloqueo de una cuenta mediante usp_User_chkBlock.
func (r *MySQLAccountRepository) CheckBlock(userNo int) (*BlockInfo, error) {
	var bIsBlock int
	var sReason sql.NullString
	var dBlockDate sql.NullTime

	err := r.db.QueryRow("CALL usp_User_chkBlock(?)", userNo).
		Scan(&bIsBlock, &sReason, &dBlockDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return &BlockInfo{UserNo: userNo, IsBlock: 0}, nil
		}
		return nil, fmt.Errorf("usp_User_chkBlock failed: %w", err)
	}

	info := &BlockInfo{
		UserNo:  userNo,
		IsBlock: bIsBlock,
	}
	if sReason.Valid {
		info.Reason = sReason.String
	}
	if dBlockDate.Valid {
		info.BlockDate = dBlockDate.Time
	}
	return info, nil
}

// GetAccountByUsername recupera los datos completos del perfil de tUser mediante usp_User_getUser.
func (r *MySQLAccountRepository) GetAccountByUsername(username string) (*Account, error) {
	var (
		userNo     int
		userID     string
		userPW     string
		userName   string
		isBlock    int
		isDelete   int
		authID     byte
		userIP     string
		createDate time.Time
	)

	err := r.db.QueryRow("CALL usp_User_getUser(?)", username).
		Scan(&userNo, &userID, &userPW, &userName, &isBlock, &isDelete, &authID, &userIP, &createDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("usp_User_getUser failed: %w", err)
	}

	acc := &Account{
		ID:       userNo,
		Username: userID,
		Password: userPW,
		Admin:    authID,
		Banned:   isBlock != 0 || isDelete != 0,
		TUser: TUser{
			UserNo:     userNo,
			UserID:     userID,
			UserPW:     userPW,
			UserName:   userName,
			IsBlock:    isBlock,
			IsDelete:   isDelete,
			AuthID:     authID,
			UserIP:     userIP,
			CreateDate: createDate,
		},
	}
	acc.Sync()
	return acc, nil
}

// CreateAccount registra una nueva cuenta en tUser y tUserOption mediante Stored Procedures.
func (r *MySQLAccountRepository) CreateAccount(username, password string) (*Account, error) {
	pwHash := strings.ToLower(password)
	if len(pwHash) != 32 {
		h := md5.Sum([]byte(password))
		pwHash = hex.EncodeToString(h[:])
	}

	var ret, newUserNo int
	err := r.db.QueryRow("CALL usp_User_insert(?, ?, '', '', '', ?, '127.0.0.1', '', '', '', 0, '', 1, 'ONSON', 1, @out_userNo)", username, pwHash, username).
		Scan(&ret, &newUserNo)
	if err != nil {
		return nil, fmt.Errorf("usp_User_insert failed: %w", err)
	}
	if ret != 0 {
		return nil, ErrAccountExists
	}

	// Inicializar registro de opciones de usuario
	_, _ = r.db.Exec("CALL usp_User_insertOption(?, '', '', '', '', '')", newUserNo)

	acc := &Account{
		ID:       newUserNo,
		Username: username,
		Password: pwHash,
		Admin:    1,
		Banned:   false,
	}
	acc.Sync()
	return acc, nil
}

// UpdatePassword actualiza la contraseña de un usuario mediante usp_User_updatePassword.
func (r *MySQLAccountRepository) UpdatePassword(userNo int, newPassword string) error {
	pwHash := strings.ToLower(newPassword)
	if len(pwHash) != 32 {
		h := md5.Sum([]byte(newPassword))
		pwHash = hex.EncodeToString(h[:])
	}

	_, err := r.db.Exec("CALL usp_User_updatePassword(?, ?)", userNo, pwHash)
	if err != nil {
		return fmt.Errorf("usp_User_updatePassword failed: %w", err)
	}
	return nil
}

// DeleteAccount efectúa la baja lógica de una cuenta mediante usp_User_delete.
func (r *MySQLAccountRepository) DeleteAccount(userNo int) error {
	_, err := r.db.Exec("CALL usp_User_delete(?)", userNo)
	if err != nil {
		return fmt.Errorf("usp_User_delete failed: %w", err)
	}
	return nil
}

// SetDrawComment registra el motivo de baja en tDrawComment mediante usp_User_setDrawComment.
func (r *MySQLAccountRepository) SetDrawComment(userNo int, comment string) error {
	_, err := r.db.Exec("CALL usp_User_setDrawComment(?, ?)", userNo, comment)
	if err != nil {
		return fmt.Errorf("usp_User_setDrawComment failed: %w", err)
	}
	return nil
}

// SetChannelLog registra el canal seleccionado en tChannelLog mediante usp_User_setChannelLog.
func (r *MySQLAccountRepository) SetChannelLog(userNo int, channelID int) error {
	_, err := r.db.Exec("CALL usp_User_setChannelLog(?, ?)", userNo, channelID)
	if err != nil {
		return fmt.Errorf("usp_User_setChannelLog failed: %w", err)
	}
	return nil
}

// CheckPcbangIP verifica si una dirección IP pertenece a un cibercafé/PC Bang mediante usp_Pcbang_chkIP.
func (r *MySQLAccountRepository) CheckPcbangIP(ip string) (bool, error) {
	var isPcbang int
	err := r.db.QueryRow("CALL usp_Pcbang_chkIP(?)", ip).Scan(&isPcbang)
	if err != nil {
		return false, fmt.Errorf("usp_Pcbang_chkIP failed: %w", err)
	}
	return isPcbang == 1, nil
}

// GetWebKey recupera la clave de autenticación web de tUserKey mediante p_UserKey_Get_FromAccount.
func (r *MySQLAccountRepository) GetWebKey(userNo int) (string, error) {
	var ret int
	var webKey sql.NullString
	var expire sql.NullTime

	err := r.db.QueryRow("CALL p_UserKey_Get_FromAccount(?, @k, @e, @r)", userNo).
		Scan(&ret, &webKey, &expire)
	if err != nil {
		return "", fmt.Errorf("p_UserKey_Get_FromAccount failed: %w", err)
	}
	if ret != 0 || !webKey.Valid {
		return "", fmt.Errorf("web key not found or expired (ret=%d)", ret)
	}
	return webKey.String, nil
}

// GetPlayRemainTime recupera el saldo de minutos restantes para menores mediante p_UserPlayRemainTime_Get_FromAccount.
func (r *MySQLAccountRepository) GetPlayRemainTime(userNo int) (int, error) {
	var remainTime int
	err := r.db.QueryRow("CALL p_UserPlayRemainTime_Get_FromAccount(?, @r)", userNo).Scan(&remainTime)
	if err != nil {
		return 0, fmt.Errorf("p_UserPlayRemainTime_Get_FromAccount failed: %w", err)
	}
	return remainTime, nil
}

// SetPlayRemainTime establece el saldo de minutos para menores mediante p_UserPlayRemainTime_Set_FromAccount.
func (r *MySQLAccountRepository) SetPlayRemainTime(userNo int, remainSec int) error {
	_, err := r.db.Exec("CALL p_UserPlayRemainTime_Set_FromAccount(?, ?)", userNo, remainSec)
	if err != nil {
		return fmt.Errorf("p_UserPlayRemainTime_Set_FromAccount failed: %w", err)
	}
	return nil
}

// GetUserDateDiffCreated consulta los segundos desde la creación de la cuenta mediante usp_User_GetDatediffCreated.
func (r *MySQLAccountRepository) GetUserDateDiffCreated(userNo int) (int64, error) {
	var diffSec int64
	err := r.db.QueryRow("CALL usp_User_GetDatediffCreated(?, @diff)", userNo).Scan(&diffSec)
	if err != nil {
		return 0, fmt.Errorf("usp_User_GetDatediffCreated failed: %w", err)
	}
	return diffSec, nil
}

// GetReturnUserDateDiffLastLog consulta los segundos desde el último retorno mediante usp_ReturnUser_GetDatediffLastlog.
func (r *MySQLAccountRepository) GetReturnUserDateDiffLastLog(userNo int) (int64, error) {
	var diffSec int64
	err := r.db.QueryRow("CALL usp_ReturnUser_GetDatediffLastlog(?, @diff)", userNo).Scan(&diffSec)
	if err != nil {
		return 0, fmt.Errorf("usp_ReturnUser_GetDatediffLastlog failed: %w", err)
	}
	return diffSec, nil
}

// SetReturnUserItem entrega el paquete de retorno mediante usp_ReturnUser_SetUserItem.
func (r *MySQLAccountRepository) SetReturnUserItem(userNo int, itemNo int) error {
	_, err := r.db.Exec("CALL usp_ReturnUser_SetUserItem(?, ?, @ret)", userNo, itemNo)
	if err != nil {
		return fmt.Errorf("usp_ReturnUser_SetUserItem failed: %w", err)
	}
	return nil
}

// GetTutorialIsSkip consulta si el usuario puede omitir el tutorial mediante usp_Tutorial_getIsSkip.
func (r *MySQLAccountRepository) GetTutorialIsSkip(userNo int) (bool, error) {
	var canSkip int
	err := r.db.QueryRow("CALL usp_Tutorial_getIsSkip(?, @ret)", userNo).Scan(&canSkip)
	if err != nil {
		return false, fmt.Errorf("usp_Tutorial_getIsSkip failed: %w", err)
	}
	return canSkip == 1, nil
}

// GetTutorialCharInfo consulta el paso actual del tutorial mediante usp_Tutorial_getCharInfo.
func (r *MySQLAccountRepository) GetTutorialCharInfo(userNo int) (int, error) {
	var nState, nStep int
	err := r.db.QueryRow("CALL usp_Tutorial_getCharInfo(?, @ret)", userNo).Scan(&nState, &nStep)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("usp_Tutorial_getCharInfo failed: %w", err)
	}
	return nStep, nil
}

// SetTutorialCharInfo actualiza el paso del tutorial mediante usp_Tutorial_setCharInfo.
func (r *MySQLAccountRepository) SetTutorialCharInfo(userNo int, step int) error {
	_, err := r.db.Exec("CALL usp_Tutorial_setCharInfo(?, 0, 1, ?, @ret)", userNo, step)
	if err != nil {
		return fmt.Errorf("usp_Tutorial_setCharInfo failed: %w", err)
	}
	return nil
}
