package login

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// World represents a World Server known to the Login Server.
type World struct {
	ID     byte
	Name   string
	Status protocol.WorldStatus
	IP     string
	Port   uint16
}

// WorldManager manages worlds dynamically from MySQL (tWorldServer) or fallback in-memory.
type WorldManager struct {
	mu          sync.RWMutex
	repo        database.LoginSessionRepository
	localWorlds map[byte]*World
}

// CalculateWorldStatus calcula dinámicamente el estado de un mundo según sus parámetros oficiales de MySQL:
// bIsMaintenance, nCurrentCCU, nMaxUser y nUserRatable (LOGIN_USER_RATABLE).
func CalculateWorldStatus(isMaintenance bool, currentCCU, maxUser, userRatable int) protocol.WorldStatus {
	if isMaintenance {
		return protocol.WorldStatus_Maintenance
	}
	if userRatable <= 0 {
		userRatable = 100
	}
	effectiveMax := maxUser * userRatable / 100
	if effectiveMax <= 0 {
		effectiveMax = 1000
	}
	if currentCCU >= effectiveMax {
		return protocol.WorldStatus_Full
	}
	pct := currentCCU * 100 / effectiveMax
	if pct >= 75 {
		return protocol.WorldStatus_High
	}
	if pct >= 40 {
		return protocol.WorldStatus_Medium
	}
	return protocol.WorldStatus_Low
}

// NewWorldManager creates a WorldManager. Si se provee LoginSessionRepository,
// los mundos se consultarán dinámicamente desde tWorldServer en MySQL.
func NewWorldManager(repos ...database.LoginSessionRepository) *WorldManager {
	wm := &WorldManager{
		localWorlds: make(map[byte]*World),
	}
	if len(repos) > 0 && repos[0] != nil {
		wm.repo = repos[0]
	} else {
		// Servidor por defecto oficial en modo standalone/mock
		wm.localWorlds[0] = &World{
			ID:     0,
			Name:   "Apoline",
			Status: protocol.WorldStatus_Low,
			IP:     "127.0.0.1",
			Port:   9110,
		}
	}
	return wm
}

// AddWorld registers or updates a world (localmente en memoria).
func (wm *WorldManager) AddWorld(w *World) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.localWorlds[w.ID] = w
}

// GetWorld returns a world by its ID (desde MySQL si hay repo disponible, sino memoria).
func (wm *WorldManager) GetWorld(id byte) (*World, bool) {
	if wm.repo != nil {
		info, err := wm.repo.GetWorldServerByNo(id)
		if err == nil && info != nil {
			st := CalculateWorldStatus(info.IsMaintenance, info.CurrentCCU, info.MaxUser, info.UserRatable)
			return &World{
				ID:     info.WorldNo,
				Name:   info.WorldName,
				Status: st,
				IP:     info.WorldIP,
				Port:   info.WorldPort,
			}, true
		}
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()
	w, ok := wm.localWorlds[id]
	return w, ok
}

// WorldsList returns a slice of all registered worlds.
func (wm *WorldManager) WorldsList() []*World {
	if wm.repo != nil {
		list, err := wm.repo.GetWorldServerList()
		if err == nil && len(list) > 0 {
			res := make([]*World, 0, len(list))
			for _, info := range list {
				st := CalculateWorldStatus(info.IsMaintenance, info.CurrentCCU, info.MaxUser, info.UserRatable)
				res = append(res, &World{
					ID:     info.WorldNo,
					Name:   info.WorldName,
					Status: st,
					IP:     info.WorldIP,
					Port:   info.WorldPort,
				})
			}
			return res
		}
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()
	res := make([]*World, 0, len(wm.localWorlds))
	for _, w := range wm.localWorlds {
		res = append(res, w)
	}
	return res
}

// GenerateTransferToken creates a random 32-character hex transfer token.
func GenerateTransferToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
