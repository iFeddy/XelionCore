package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	TargetMap string    `json:"target_map"`
	TargetX   int       `json:"target_x"`
	TargetY   int       `json:"target_y"`
	WorldNo   byte      `json:"world_no"`
	ZoneNo    uint16    `json:"zone_no"`
	CreatedAt time.Time `json:"created_at"`
}

// TransferManager gestiona las sesiones de transferencia de forma atómica y distribuida
// interactuando 100% mediante Stored Procedures de MySQL (usp_Session_*).
// Regla Mandatoria: CERO mapas de transferencia en memoria (transfers/zoneTransfers erradicados).
type TransferManager struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewTransferManager inicializa el gestor de transferencias respaldado en MySQL.
func NewTransferManager(db ...*sql.DB) *TransferManager {
	tm := &TransferManager{}
	if len(db) > 0 && db[0] != nil {
		tm.db = db[0]
	} else {
		// Conexión por defecto a MySQL xelion_online para interoperabilidad distribuida transparente
		if defaultDB, err := sql.Open("mysql", "root:lerd@tcp(127.0.0.1:3306)/xelion_online?parseTime=true"); err == nil {
			if err := defaultDB.Ping(); err == nil {
				tm.db = defaultDB
			}
		}
	}
	return tm
}

// SetDB asigna la conexión de base de datos MySQL al gestor de transferencias.
func (tm *TransferManager) SetDB(db *sql.DB) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.db = db
}

// DB retorna la conexión a base de datos activa.
func (tm *TransferManager) DB() *sql.DB {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.db
}

// CreateTransfer registra un token criptográfico para LoginServer -> WorldServer
// invocando usp_Session_IssueWorldToken en MySQL.
func (tm *TransferManager) CreateTransfer(accountID int, username string, admin byte, hostIP string, opts ...interface{}) string {
	var token64 [64]byte
	var token string
	worldNo := byte(0)

	for _, opt := range opts {
		switch v := opt.(type) {
		case []byte:
			if len(v) > 0 {
				copy(token64[:], v)
				token = hex.EncodeToString(token64[:16])
			}
		case byte:
			worldNo = v
		case int:
			worldNo = byte(v)
		}
	}

	if token == "" {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			copy(b, []byte(time.Now().Format(time.RFC3339Nano)))
		}
		token = hex.EncodeToString(b)
		copy(token64[:], []byte(token))
	}

	tm.mu.RLock()
	db := tm.db
	tm.mu.RUnlock()

	if db != nil {
		var regRet int
		_ = db.QueryRow("CALL usp_Session_RegisterLogin(?, ?, ?)", accountID, username, hostIP).Scan(&regRet)
		var ret int
		_ = db.QueryRow("CALL usp_Session_IssueWorldToken(?, ?, ?)", accountID, worldNo, token64[:]).Scan(&ret)
	}

	return token
}

// ConsumeTransfer valida y consume atómicamente el token de WorldServer invocando usp_Session_ConsumeWorldToken.
func (tm *TransferManager) ConsumeTransfer(token string, raw []byte, opts ...interface{}) (*ClientTransfer, bool) {
	token = strings.TrimSpace(token)
	if token == "" && len(raw) == 0 {
		return nil, false
	}

	worldNo := byte(0)
	for _, opt := range opts {
		switch v := opt.(type) {
		case byte:
			worldNo = v
		case int:
			worldNo = byte(v)
		}
	}

	tm.mu.RLock()
	db := tm.db
	tm.mu.RUnlock()

	if db == nil {
		return nil, false
	}

	token64 := make([]byte, 64)
	if len(raw) >= 64 {
		copy(token64, raw[:64])
	} else {
		copy(token64, []byte(token))
	}

	hostIP := ""
	var ret, userNo int
	var sUserID string
	err := db.QueryRow("CALL usp_Session_ConsumeWorldToken(?, ?, ?)", worldNo, token64, hostIP).Scan(&ret, &userNo, &sUserID)
	if (err != nil || ret != 0 || userNo <= 0) && worldNo == 0 {
		// Fallback intentando world 1
		_ = db.QueryRow("CALL usp_Session_ConsumeWorldToken(1, ?, ?)", token64, hostIP).Scan(&ret, &userNo, &sUserID)
	}
	if ret != 0 || userNo <= 0 {
		return nil, false
	}

	return &ClientTransfer{
		Token:     token,
		AccountID: userNo,
		Username:  sUserID,
		Admin:     0,
		HostIP:    hostIP,
		CreatedAt: time.Now(),
	}, true
}

// CreateZoneTransfer registra un traspaso de WorldServer a ZoneServer invocando usp_Session_IssueZoneTransfer.
func (tm *TransferManager) CreateZoneTransfer(accountID int, username, charName string, charID int, randomID uint16, admin byte, hostIP string) {
	tm.mu.RLock()
	db := tm.db
	tm.mu.RUnlock()

	if db == nil {
		return
	}

	var ret int
	_ = db.QueryRow(
		"CALL usp_Session_IssueZoneTransfer(1, 1, ?, ?, ?, ?, ?, ?, ?, '', 0, 0)",
		accountID, username, charID, charName, randomID, admin, hostIP,
	).Scan(&ret)
}

// ConsumeZoneTransfer valida, consume y destruye de forma atómica y transaccional el traspaso a Zona
// invocando usp_Session_ConsumeZoneTransfer en MySQL con protección anti-replay.
func (tm *TransferManager) ConsumeZoneTransfer(charName string, randomID uint16, raw []byte) (*ZoneTransfer, bool) {
	tm.mu.RLock()
	db := tm.db
	tm.mu.RUnlock()

	if db == nil {
		return nil, false
	}

	var (
		worldNo, zoneNo, userNo, charNo, adminLevel, targetX, targetY, ret int
		userID, cID, targetMap                                             string
	)

	row := db.QueryRow("CALL usp_Session_ConsumeZoneTransfer(?, ?, '')", strings.TrimSpace(charName), randomID)
	err := row.Scan(&worldNo, &zoneNo, &userNo, &userID, &charNo, &cID, &adminLevel, &targetMap, &targetX, &targetY, &ret)
	if err != nil || ret != 0 || charNo <= 0 {
		return nil, false
	}

	return &ZoneTransfer{
		AccountID: userNo,
		Username:  userID,
		CharName:  cID,
		CharID:    charNo,
		RandomID:  randomID,
		Admin:     byte(adminLevel),
		HostIP:    "",
		TargetMap: targetMap,
		TargetX:   targetX,
		TargetY:   targetY,
		WorldNo:   byte(worldNo),
		ZoneNo:    uint16(zoneNo),
		CreatedAt: time.Now(),
	}, true
}
