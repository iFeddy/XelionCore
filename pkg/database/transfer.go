package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ClientTransfer representa el traspaso de un cliente autenticado desde LoginServer a WorldServer.
type ClientTransfer struct {
	Token     string    `json:"token"`
	AccountID int       `json:"account_id"`
	Username  string    `json:"username"`
	Admin     byte      `json:"admin"`
	HostIP    string    `json:"host_ip"`
	CreatedAt time.Time `json:"created_at"`
}

// ZoneTransfer representa el traspaso de un jugador desde WorldServer a ZoneServer.
type ZoneTransfer struct {
	AccountID int       `json:"account_id"`
	Username  string    `json:"username"`
	CharName  string    `json:"char_name"`
	CharID    int       `json:"char_id"`
	RandomID  uint16    `json:"random_id"`
	Admin     byte      `json:"admin"`
	HostIP    string    `json:"host_ip"`
	CreatedAt time.Time `json:"created_at"`
}

// TransferManager gestiona las sesiones de transferencia 100% en memoria RAM
// (utilizando tabla MySQL ENGINE=MEMORY y maps sincronizados, 0 escrituras a disco).
type TransferManager struct {
	mu            sync.RWMutex
	db            *sql.DB
	transfers     map[string]*ClientTransfer
	zoneTransfers map[string]*ZoneTransfer
}

// NewTransferManager inicializa el gestor de transferencias. Si se provee una conexión MySQL,
// sincroniza en la tabla volátil tTransferToken (ENGINE=MEMORY).
func NewTransferManager(db ...*sql.DB) *TransferManager {
	tm := &TransferManager{
		transfers:     make(map[string]*ClientTransfer),
		zoneTransfers: make(map[string]*ZoneTransfer),
	}
	if len(db) > 0 && db[0] != nil {
		tm.db = db[0]
	}
	return tm
}

// CreateTransfer registra un nuevo token de transferencia para el cliente (Login -> World).
func (tm *TransferManager) CreateTransfer(accountID int, username string, admin byte, hostIP string) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	now := time.Now()
	transfer := &ClientTransfer{
		Token:     token,
		AccountID: accountID,
		Username:  username,
		Admin:     admin,
		HostIP:    hostIP,
		CreatedAt: now,
	}

	tm.mu.Lock()
	tm.transfers[token] = transfer
	tm.mu.Unlock()

	// Persistencia volátil en MySQL MEMORY (sin tocar disco físico)
	if tm.db != nil {
		query := `
			INSERT INTO tTransferToken (sToken, nUserNo, sUserID, nAuthID, sIP, dCreatedAt)
			VALUES (?, ?, ?, ?, ?, NOW());
		`
		_, _ = tm.db.Exec(query, token, accountID, username, admin, hostIP)

		// Limpieza de tokens viejos expirados (> 120 segundos)
		go func(db *sql.DB) {
			_, _ = db.Exec("DELETE FROM tTransferToken WHERE dCreatedAt < NOW() - INTERVAL 120 SECOND")
		}(tm.db)
	}

	return token
}

// ConsumeTransfer valida y consume un token de transferencia (Login -> World).
// Al ser de un solo uso, se remueve inmediatamente tras su validación estricta por token.
func (tm *TransferManager) ConsumeTransfer(token string, raw []byte, username ...string) (*ClientTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, false
	}

	// 1. Búsqueda directa en memoria local por token oficial
	if t, exists := tm.transfers[token]; exists {
		delete(tm.transfers, token)
		if tm.db != nil {
			_, _ = tm.db.Exec("DELETE FROM tTransferToken WHERE sToken = ?", token)
		}
		return t, true
	}

	// 2. Consulta en MySQL tTransferToken (inter-proceso LoginService -> WorldService)
	if tm.db != nil {
		query := `
			SELECT sToken, nUserNo, sUserID, nAuthID, sIP, dCreatedAt 
			FROM tTransferToken 
			WHERE sToken = ? AND dCreatedAt >= NOW() - INTERVAL 120 SECOND
			LIMIT 1;
		`
		var (
			sTok      string
			userNo    int
			userID    string
			authID    byte
			ip        string
			createdAt time.Time
		)

		err := tm.db.QueryRow(query, token).Scan(&sTok, &userNo, &userID, &authID, &ip, &createdAt)
		if err == nil {
			// Borrado atómico de un solo uso
			_, _ = tm.db.Exec("DELETE FROM tTransferToken WHERE sToken = ?", sTok)
			return &ClientTransfer{
				Token:     sTok,
				AccountID: userNo,
				Username:  userID,
				Admin:     authID,
				HostIP:    ip,
				CreatedAt: createdAt,
			}, true
		}
	}

	return nil, false
}

// CreateZoneTransfer registra un nuevo traspaso desde WorldServer a ZoneServer.
func (tm *TransferManager) CreateZoneTransfer(accountID int, username, charName string, charID int, randomID uint16, admin byte, hostIP string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(charName))
	tm.zoneTransfers[key] = &ZoneTransfer{
		AccountID: accountID,
		Username:  username,
		CharName:  charName,
		CharID:    charID,
		RandomID:  randomID,
		Admin:     admin,
		HostIP:    hostIP,
		CreatedAt: time.Now(),
	}

	if tm.db != nil {
		tok := fmt.Sprintf("ZONE:%s", key)
		payload := fmt.Sprintf("%s|%d|%d", username, charID, randomID)
		query := `
			INSERT INTO tTransferToken (sToken, nUserNo, sUserID, nAuthID, sIP, dCreatedAt)
			VALUES (?, ?, ?, ?, ?, NOW())
			ON DUPLICATE KEY UPDATE nUserNo = VALUES(nUserNo), sUserID = VALUES(sUserID), nAuthID = VALUES(nAuthID), sIP = VALUES(sIP), dCreatedAt = NOW();
		`
		_, _ = tm.db.Exec(query, tok, accountID, payload, admin, hostIP)
	}
}

// ConsumeZoneTransfer valida y consume un token de traspaso a Zona (World -> Zone).
func (tm *TransferManager) ConsumeZoneTransfer(charName string, randomID uint16, raw []byte) (*ZoneTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(charName))

	// 1. Búsqueda directa en memoria local por nombre de personaje
	if t, exists := tm.zoneTransfers[key]; exists {
		delete(tm.zoneTransfers, key)
		if tm.db != nil {
			_, _ = tm.db.Exec("DELETE FROM tTransferToken WHERE sToken = ?", fmt.Sprintf("ZONE:%s", key))
		}
		return t, true
	}

	// 2. Búsqueda en memoria local por RandomID
	for k, t := range tm.zoneTransfers {
		if t.RandomID == randomID {
			delete(tm.zoneTransfers, k)
			if tm.db != nil {
				_, _ = tm.db.Exec("DELETE FROM tTransferToken WHERE sToken = ?", fmt.Sprintf("ZONE:%s", k))
			}
			return t, true
		}
	}

	// 3. Consulta en MySQL tTransferToken (inter-proceso WorldService -> ZoneService)
	if tm.db != nil {
		tok := fmt.Sprintf("ZONE:%s", key)
		query := `
			SELECT sToken, nUserNo, sUserID, nAuthID, sIP, dCreatedAt
			FROM tTransferToken
			WHERE sToken = ? AND dCreatedAt >= NOW() - INTERVAL 120 SECOND
			LIMIT 1;
		`
		var (
			sTok      string
			userNo    int
			userID    string
			authID    byte
			ip        string
			createdAt time.Time
		)
		err := tm.db.QueryRow(query, tok).Scan(&sTok, &userNo, &userID, &authID, &ip, &createdAt)
		if err == nil {
			_, _ = tm.db.Exec("DELETE FROM tTransferToken WHERE sToken = ?", sTok)
			parts := strings.Split(userID, "|")
			var uName string
			var cID int
			var rID uint16
			if len(parts) >= 3 {
				uName = parts[0]
				cID, _ = strconv.Atoi(parts[1])
				rVal, _ := strconv.ParseUint(parts[2], 10, 16)
				rID = uint16(rVal)
			} else {
				uName = userID
			}
			return &ZoneTransfer{
				AccountID: userNo,
				Username:  uName,
				CharName:  charName,
				CharID:    cID,
				RandomID:  rID,
				Admin:     authID,
				HostIP:    ip,
				CreatedAt: createdAt,
			}, true
		}
	}

	return nil, false
}
