package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ClientTransfer represents a handoff from Login Server to World Server.
type ClientTransfer struct {
	Token     string    `json:"token"`
	AccountID int       `json:"account_id"`
	Username  string    `json:"username"`
	Admin     byte      `json:"admin"`
	HostIP    string    `json:"host_ip"`
	CreatedAt time.Time `json:"created_at"`
}

// ZoneTransfer represents a handoff from World Server to Zone Server.
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

type transferStore struct {
	Transfers     map[string]*ClientTransfer `json:"transfers"`
	ZoneTransfers map[string]*ZoneTransfer   `json:"zone_transfers"`
}

// TransferManager handles tokens for inter-service client handoffs.
type TransferManager struct {
	mu            sync.RWMutex
	filePath      string
	transfers     map[string]*ClientTransfer
	zoneTransfers map[string]*ZoneTransfer
}

// NewTransferManager creates a new TransferManager, persisting to data/transfers.json if possible.
func NewTransferManager() *TransferManager {
	filePath := filepath.Join("data", "transfers.json")
	tm := &TransferManager{
		filePath:      filePath,
		transfers:     make(map[string]*ClientTransfer),
		zoneTransfers: make(map[string]*ZoneTransfer),
	}
	tm.load()
	return tm
}

func (tm *TransferManager) load() {
	if tm.filePath == "" {
		return
	}
	data, err := os.ReadFile(tm.filePath)
	if err != nil {
		return
	}
	var store transferStore
	if err := json.Unmarshal(data, &store); err == nil {
		if store.Transfers != nil {
			tm.transfers = store.Transfers
		}
		if store.ZoneTransfers != nil {
			tm.zoneTransfers = store.ZoneTransfers
		}
	}
}

func (tm *TransferManager) save() {
	if tm.filePath == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(tm.filePath), 0755)
	store := transferStore{
		Transfers:     tm.transfers,
		ZoneTransfers: tm.zoneTransfers,
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(tm.filePath, data, 0644)
	}
}

// CreateTransfer registers a new transfer token for a client (Login -> World).
func (tm *TransferManager) CreateTransfer(accountID int, username string, admin byte, hostIP string) string {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.load()

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	tm.transfers[token] = &ClientTransfer{
		Token:     token,
		AccountID: accountID,
		Username:  username,
		Admin:     admin,
		HostIP:    hostIP,
		CreatedAt: time.Now(),
	}
	tm.save()

	return token
}

// ConsumeTransfer verifies and removes a transfer token (Login -> World).
func (tm *TransferManager) ConsumeTransfer(token string, raw []byte) (*ClientTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.load()

	// 1. Direct match
	if token != "" {
		if t, exists := tm.transfers[token]; exists {
			delete(tm.transfers, token)
			tm.save()
			return t, true
		}
	}

	// 2. Substring match in raw payload
	rawStr := string(raw)
	for key, t := range tm.transfers {
		if strings.Contains(rawStr, key) {
			delete(tm.transfers, key)
			tm.save()
			return t, true
		}
	}

	// 3. Fallback: match most recent pending transfer (within 60 seconds)
	now := time.Now()
	for key, t := range tm.transfers {
		if now.Sub(t.CreatedAt) < 60*time.Second {
			delete(tm.transfers, key)
			tm.save()
			return t, true
		}
	}

	return nil, false
}

// CreateZoneTransfer registers a new handoff from World to Zone.
func (tm *TransferManager) CreateZoneTransfer(accountID int, username, charName string, charID int, randomID uint16, admin byte, hostIP string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.load()

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
	tm.save()
}

// ConsumeZoneTransfer verifies and removes a transfer token (World -> Zone).
func (tm *TransferManager) ConsumeZoneTransfer(charName string, randomID uint16, raw []byte) (*ZoneTransfer, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.load()

	key := strings.ToLower(strings.TrimSpace(charName))

	// 1. Direct match by charName
	if t, exists := tm.zoneTransfers[key]; exists {
		delete(tm.zoneTransfers, key)
		tm.save()
		return t, true
	}

	// 2. Match by randomID
	for k, t := range tm.zoneTransfers {
		if t.RandomID == randomID {
			delete(tm.zoneTransfers, k)
			tm.save()
			return t, true
		}
	}

	// 3. Fallback: most recent transfer within 60 seconds
	now := time.Now()
	for k, t := range tm.zoneTransfers {
		if now.Sub(t.CreatedAt) < 60*time.Second {
			delete(tm.zoneTransfers, k)
			tm.save()
			return t, true
		}
	}

	return nil, false
}
