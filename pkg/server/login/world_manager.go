package login

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

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

// WorldManager manages worlds and transfer tokens.
type WorldManager struct {
	mu     sync.RWMutex
	worlds map[byte]*World
}

// NewWorldManager creates a WorldManager with default world entries.
func NewWorldManager() *WorldManager {
	wm := &WorldManager{
		worlds: make(map[byte]*World),
	}

	// Default development World (World 0: INITIO)
	wm.worlds[0] = &World{
		ID:     0,
		Name:   "INITIO",
		Status: protocol.WorldStatus_Low,
		IP:     "127.0.0.1",
		Port:   9110,
	}

	return wm
}

// AddWorld registers or updates a world.
func (wm *WorldManager) AddWorld(w *World) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.worlds[w.ID] = w
}

// GetWorld returns a world by its ID.
func (wm *WorldManager) GetWorld(id byte) (*World, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	w, ok := wm.worlds[id]
	return w, ok
}

// WorldsList returns a slice of all registered worlds.
func (wm *WorldManager) WorldsList() []*World {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	list := make([]*World, 0, len(wm.worlds))
	for _, w := range wm.worlds {
		list = append(list, w)
	}
	return list
}

// GenerateTransferToken creates a random 32-character hex transfer token.
func GenerateTransferToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
