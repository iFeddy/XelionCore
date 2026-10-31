package database

import (
	"database/sql"
	"fmt"
)

// GetClientVersionKey obtiene la versión oficial esperada desde FOUS mediante usp_Version_GetClientKey.
func (r *MySQLAccountRepository) GetClientVersionKey() (string, error) {
	var key string
	err := r.db.QueryRow("CALL usp_Version_GetClientKey()").Scan(&key)
	if err != nil {
		return "", fmt.Errorf("usp_Version_GetClientKey failed: %w", err)
	}
	return key, nil
}

// GetWorldServerList obtiene el listado de mundos registrados en tWorldServer mediante usp_WorldServer_GetList.
func (r *MySQLAccountRepository) GetWorldServerList() ([]*WorldServerInfo, error) {
	rows, err := r.db.Query("CALL usp_WorldServer_GetList()")
	if err != nil {
		return nil, fmt.Errorf("usp_WorldServer_GetList failed: %w", err)
	}
	defer rows.Close()

	var list []*WorldServerInfo
	for rows.Next() {
		var (
			worldNo       uint8
			worldName     string
			worldIP       string
			worldPort     uint16
			maxUser       int
			userRatable   int
			isMaintenance bool
			currentCCU    int
			lastHeartbeat sql.NullTime
			worldStatus   uint8
		)
		err := rows.Scan(&worldNo, &worldName, &worldIP, &worldPort,
			&maxUser, &userRatable, &isMaintenance, &currentCCU, &lastHeartbeat, &worldStatus)
		if err != nil {
			return nil, fmt.Errorf("error scanning world server: %w", err)
		}
		list = append(list, &WorldServerInfo{
			WorldNo:       worldNo,
			WorldName:     worldName,
			WorldIP:       worldIP,
			WorldPort:     worldPort,
			WorldStatus:   worldStatus,
			CurrentCCU:    currentCCU,
			MaxUser:       maxUser,
			IsMaintenance: isMaintenance,
		})
	}
	return list, nil
}

// GetWorldServerByNo obtiene los datos operativos de un mundo específico mediante usp_WorldServer_GetByNo.
func (r *MySQLAccountRepository) GetWorldServerByNo(worldNo uint8) (*WorldServerInfo, error) {
	var (
		no            uint8
		worldName     string
		worldIP       string
		worldPort     uint16
		maxUser       int
		userRatable   int
		isMaintenance bool
		currentCCU    int
		lastHeartbeat sql.NullTime
		worldStatus   uint8
	)
	err := r.db.QueryRow("CALL usp_WorldServer_GetByNo(?)", worldNo).Scan(
		&no, &worldName, &worldIP, &worldPort,
		&maxUser, &userRatable, &isMaintenance, &currentCCU, &lastHeartbeat, &worldStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("world %d not found", worldNo)
		}
		return nil, fmt.Errorf("usp_WorldServer_GetByNo failed: %w", err)
	}

	return &WorldServerInfo{
		WorldNo:       no,
		WorldName:     worldName,
		WorldIP:       worldIP,
		WorldPort:     worldPort,
		WorldStatus:   worldStatus,
		CurrentCCU:    currentCCU,
		MaxUser:       maxUser,
		IsMaintenance: isMaintenance,
	}, nil
}

// UpdateWorldServerCCU actualiza la concurrencia y latido del servidor de mundo mediante usp_WorldServer_UpdateCCU.
func (r *MySQLAccountRepository) UpdateWorldServerCCU(worldNo uint8, ccu int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_WorldServer_UpdateCCU(?, ?)", worldNo, ccu).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_WorldServer_UpdateCCU failed: %w", err)
	}
	return nil
}

// RegisterLoginSession registra la sesión en tActiveAccountSession mediante usp_Session_RegisterLogin.
// Devuelve true si existía una sesión previa activa que debe ser expulsada (KICKOFFFROMWORLD_CMD).
func (r *MySQLAccountRepository) RegisterLoginSession(userNo int, userID, ip string) (bool, error) {
	var wasActive int
	err := r.db.QueryRow("CALL usp_Session_RegisterLogin(?, ?, ?)", userNo, userID, ip).Scan(&wasActive)
	if err != nil {
		return false, fmt.Errorf("usp_Session_RegisterLogin failed: %w", err)
	}
	return wasActive == 1, nil
}

// IssueWorldToken almacena el token criptográfico de transferencia de 64 bytes mediante usp_Session_IssueWorldToken.
func (r *MySQLAccountRepository) IssueWorldToken(userNo int, worldNo uint8, token [64]byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Session_IssueWorldToken(?, ?, ?)", userNo, worldNo, token[:]).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_Session_IssueWorldToken failed: %w", err)
	}
	return nil
}

// ConsumeWorldToken valida y consume de un solo uso el token de 64 bytes mediante usp_Session_ConsumeWorldToken.
func (r *MySQLAccountRepository) ConsumeWorldToken(worldNo uint8, token [64]byte, ip string) (int, string, error) {
	var (
		ret    int
		userNo int
		userID string
	)
	err := r.db.QueryRow("CALL usp_Session_ConsumeWorldToken(?, ?, ?)", worldNo, token[:], ip).
		Scan(&ret, &userNo, &userID)
	if err != nil {
		return 0, "", fmt.Errorf("usp_Session_ConsumeWorldToken failed: %w", err)
	}
	if ret != 0 || userNo <= 0 {
		return 0, "", fmt.Errorf("invalid or expired transfer token (ret=%d)", ret)
	}
	return userNo, userID, nil
}

// CreateOTP genera o persiste un One-Time Password de 32 bytes mediante usp_Session_CreateOTP.
func (r *MySQLAccountRepository) CreateOTP(userNo int, otp [32]byte) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Session_CreateOTP(?, ?)", userNo, string(otp[:])).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_Session_CreateOTP failed: %w", err)
	}
	return nil
}

// ConsumeOTP valida y consume el One-Time Password de 32 bytes mediante usp_Session_ConsumeOTP.
func (r *MySQLAccountRepository) ConsumeOTP(otp [32]byte, ip string) (int, string, error) {
	var (
		ret    int
		userNo int
		userID string
	)
	err := r.db.QueryRow("CALL usp_Session_ConsumeOTP(?, ?)", string(otp[:]), ip).
		Scan(&ret, &userNo, &userID)
	if err != nil {
		return 0, "", fmt.Errorf("usp_Session_ConsumeOTP failed: %w", err)
	}
	if ret != 0 || userNo <= 0 {
		return 0, "", fmt.Errorf("invalid or expired OTP token (ret=%d)", ret)
	}
	return userNo, userID, nil
}

// CloseSessionAndGetPlayMin finaliza la sesión activa y calcula los minutos jugados mediante usp_Session_CloseAndGetPlayMin.
func (r *MySQLAccountRepository) CloseSessionAndGetPlayMin(userNo int) (uint16, error) {
	var playMin int
	err := r.db.QueryRow("CALL usp_Session_CloseAndGetPlayMin(?)", userNo).Scan(&playMin)
	if err != nil {
		return 0, fmt.Errorf("usp_Session_CloseAndGetPlayMin failed: %w", err)
	}
	if playMin < 0 {
		playMin = 0
	}
	return uint16(playMin), nil
}
