package zone

import (
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
)

// ZoneMonster representa un monstruo activo en una instancia de ZoneMap.
type ZoneMonster struct {
	mu           sync.RWMutex
	Handle       uint16
	MobID        uint16
	Name         string
	Map          string
	X            int32
	Y            int32
	Direct       int16
	Level        uint8
	CurHP        int
	MaxHP        int
	CurSP        int
	MaxSP        int
	AC           uint32
	MR           uint32
	MinWC        uint32
	MaxWC        uint32
	MinMA        uint32
	MaxMA        uint32
	TH           uint32
	TB           uint32
	Exp          uint32
	IsDead       bool
	RespawnTime  time.Time
	TargetHandle uint16
	ActiveStates map[uint16]*ActiveAbState
}

// NewZoneMonster inicializa un monstruo con stats de combate.
func NewZoneMonster(handle uint16, mobID uint16, name, mapID string, x, y int32, direct int16, level uint8, hp, sp int, ac, mr, minWC, maxWC, th, tb uint32, exp uint32) *ZoneMonster {
	if hp <= 0 {
		hp = 100
	}
	if maxWC <= 0 {
		maxWC = 10
	}
	return &ZoneMonster{
		Handle:       handle,
		MobID:        mobID,
		Name:         name,
		Map:          mapID,
		X:            x,
		Y:            y,
		Direct:       direct,
		Level:        level,
		CurHP:        hp,
		MaxHP:        hp,
		CurSP:        sp,
		MaxSP:        sp,
		AC:           ac,
		MR:           mr,
		MinWC:        minWC,
		MaxWC:        maxWC,
		TH:           th,
		TB:           tb,
		Exp:          exp,
		ActiveStates: make(map[uint16]*ActiveAbState),
	}
}

// GetHandle implementa GridEntity para ZoneMonster.
func (m *ZoneMonster) GetHandle() uint16 {
	return m.Handle
}

// GetPosition implementa GridEntity para ZoneMonster.
func (m *ZoneMonster) GetPosition() (int, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int(m.X), int(m.Y)
}

// GetEntityType implementa GridEntity para ZoneMonster.
func (m *ZoneMonster) GetEntityType() GridEntityType {
	return GridEntityMob
}

// GetClient implementa GridEntity para ZoneMonster.
func (m *ZoneMonster) GetClient() *network.Client {
	return nil
}

// AddAbState aplica un estado anormal al monstruo.
func (m *ZoneMonster) AddAbState(st *ActiveAbState) {
	if st == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveStates[st.ID] = st
}

// RemoveAbState remueve un estado anormal del monstruo.
func (m *ZoneMonster) RemoveAbState(id uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ActiveStates, id)
}

// GetAbState obtiene un estado anormal por su ID.
func (m *ZoneMonster) GetAbState(id uint16) *ActiveAbState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ActiveStates[id]
}
