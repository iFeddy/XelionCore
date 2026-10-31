package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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

// TransferManager gestiona las sesiones de transferencia de forma segura.
type TransferManager struct {
	mu            sync.RWMutex
	db            *sql.DB
	transfers     map[string]*ClientTransfer
	zoneTransfers map[string]*ZoneTransfer
}

// NewTransferManager inicializa el gestor de transferencias.
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

	return token
}

// ConsumeTransfer valida y consume un token de transferencia (Login -> World).
func (tm *TransferManager) ConsumeTransfer(token string, raw []byte, username ...string) (*ClientTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, false
	}

	if t, exists := tm.transfers[token]; exists {
		delete(tm.transfers, token)
		return t, true
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
}

// ConsumeZoneTransfer valida y consume un token de traspaso a Zona (World -> Zone).
func (tm *TransferManager) ConsumeZoneTransfer(charName string, randomID uint16, raw []byte) (*ZoneTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(charName))

	if t, exists := tm.zoneTransfers[key]; exists {
		delete(tm.zoneTransfers, key)
		return t, true
	}

	for k, t := range tm.zoneTransfers {
		if t.RandomID == randomID {
			delete(tm.zoneTransfers, k)
			return t, true
		}
	}

	return nil, false
}
