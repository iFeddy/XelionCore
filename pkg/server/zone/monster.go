package zone

import (
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/mobattack"
	"github.com/iFeddy/XelionCore/pkg/data/mobroam"
	"github.com/iFeddy/XelionCore/pkg/data/mobsetting"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// MobAIState define los estados oficiales de la máquina de comportamiento del monstruo.
type MobAIState int

const (
	MobAIIdle   MobAIState = iota // En reposo esperando estímulo o patrullaje
	MobAIRoam                    // Patrullaje siguiendo MobRoam o desplazamiento aleatorio
	MobAIChase                   // Persiguiendo al objetivo en tabla de aggro
	MobAIAttack                  // Ejecutando secuencia de ataque / skill
	MobAIReturn                  // Retornando (leash) al origen de respawn
	MobAIDead                    // Muerto esperando timer de respawn
)

// ZoneMonster representa un monstruo activo instanciado en un ZoneMap con todo su ciclo de vida, IA y combate.
type ZoneMonster struct {
	mu           sync.RWMutex
	Handle       uint16
	MobID        uint16
	Name         string // Nombre legible (MobInfo.Name)
	InxName      string // Índice oficial (MobInfo.InxName)
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

	// Origen de respawn y parámetros de MobRegen / MobInfoServer
	SpawnX           int32
	SpawnY           int32
	SpawnDirect      int16
	RegenGroupID     string
	RegenIntervalSec int
	EXPRange         int
	WalkSpeed        uint16
	RunSpeed         uint16
	DetectRange      int
	FollowRange      int
	AttackRange      int
	AttackDelayMs    int
	EnemyDetectType  int // > 0 => agresivo (auto-detecta jugadores)
	ReturnToRegen    bool
	AllCanLoot       bool
	Visible          bool
	IsKillableByUser bool
	Rank             uint8
	GradeType        int
	MobType          int
	RaceType         int
	FamilyArea       int
	FamilyRescueArea int
	FamilyRescueCnt  int
	PeaceHPRecovDly  int
	PeaceHPRecov     int
	IsRoaming        bool
	RoamDistance     int
	RoamRestTimeMs   int

	// Estado de IA, patrullaje y combate
	AIState        MobAIState
	Route          *mobroam.Route
	RoamIndex      int
	RoamReverse    bool
	RoamWaitUntil  time.Time
	DestX          int32
	DestY          int32
	IsMoving       bool
	LastMoveTick   time.Time
	LastAttackAt   time.Time
	LastDamagedAt  time.Time
	LastPeaceHPAt  time.Time
	AttackSeqIndex int
	MobAI          *mobattack.MobAI
	Triggers       []mobsetting.ActionTrigger
	FiredTriggers  map[int]bool
	Aggro          map[uint16]int // TargetHandle -> threat score
	DamageContrib  map[int]int    // CharNo -> total damage dealt
	LastAttacker   uint16
	LastKillerChar int
	HelpCalled     bool
	Animation      string
	AnimationLevel uint8
	KQTeamType     uint8
	AIScriptFunc   string
	BehaviorFunc   string
}

// NewZoneMonster inicializa un monstruo con stats de combate compatibles con llamadas previas.
func NewZoneMonster(handle uint16, mobID uint16, name, mapID string, x, y int32, direct int16, level uint8, hp, sp int, ac, mr, minWC, maxWC, th, tb uint32, exp uint32) *ZoneMonster {
	if hp <= 0 {
		hp = 100
	}
	if maxWC <= 0 {
		maxWC = 10
	}
	return &ZoneMonster{
		Handle:           handle,
		MobID:            mobID,
		Name:             name,
		InxName:          name,
		Map:              mapID,
		X:                x,
		Y:                y,
		SpawnX:           x,
		SpawnY:           y,
		SpawnDirect:      direct,
		Direct:           direct,
		Level:            level,
		CurHP:            hp,
		MaxHP:            hp,
		CurSP:            sp,
		MaxSP:            sp,
		AC:               ac,
		MR:               mr,
		MinWC:            minWC,
		MaxWC:            maxWC,
		TH:               th,
		TB:               tb,
		Exp:              exp,
		WalkSpeed:        35,
		RunSpeed:         115,
		DetectRange:      264,
		FollowRange:      1500,
		AttackRange:      120,
		AttackDelayMs:    1300,
		EXPRange:         1500,
		RegenIntervalSec: 25,
		ReturnToRegen:    true,
		Visible:          true,
		IsKillableByUser: true,
		AIState:          MobAIIdle,
		ActiveStates:     make(map[uint16]*ActiveAbState),
		Aggro:            make(map[uint16]int),
		DamageContrib:    make(map[int]int),
		FiredTriggers:    make(map[int]bool),
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
	if m.ActiveStates == nil {
		m.ActiveStates = make(map[uint16]*ActiveAbState)
	}
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

// AddAggro incrementa la amenaza hacia un objetivo y actualiza el TargetHandle principal.
func (m *ZoneMonster) AddAggro(targetHandle uint16, amount int) {
	if targetHandle == 0 || amount <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.IsDead {
		return
	}
	if m.Aggro == nil {
		m.Aggro = make(map[uint16]int)
	}
	m.Aggro[targetHandle] += amount
	m.LastAttacker = targetHandle
	if m.AIState == MobAIIdle || m.AIState == MobAIRoam || m.AIState == MobAIReturn {
		m.AIState = MobAIChase
	}
	m.reselectTopTargetLocked()
}

// RecordDamage registra el daño infligido por un personaje y añade threat proporcional.
func (m *ZoneMonster) RecordDamage(attackerHandle uint16, charNo int, dmg int) {
	if dmg < 1 {
		dmg = 1
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.DamageContrib == nil {
		m.DamageContrib = make(map[int]int)
	}
	if charNo > 0 {
		m.DamageContrib[charNo] += dmg
		m.LastKillerChar = charNo
	}
	if attackerHandle != 0 {
		if m.Aggro == nil {
			m.Aggro = make(map[uint16]int)
		}
		m.Aggro[attackerHandle] += dmg
		m.LastAttacker = attackerHandle
		m.LastDamagedAt = time.Now()
		if m.AIState == MobAIIdle || m.AIState == MobAIRoam {
			m.AIState = MobAIChase
		}
		m.reselectTopTargetLocked()
	}
}

func (m *ZoneMonster) reselectTopTargetLocked() {
	var bestH uint16
	bestScore := -1
	for h, score := range m.Aggro {
		if score > bestScore {
			bestScore = score
			bestH = h
		}
	}
	m.TargetHandle = bestH
}

// RemoveAggroTarget elimina un objetivo de la tabla de aggro y recalcula el objetivo principal.
func (m *ZoneMonster) RemoveAggroTarget(targetHandle uint16) uint16 {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Aggro, targetHandle)
	m.reselectTopTargetLocked()
	return m.TargetHandle
}

// ResetAggro limpia completamente la tabla de aggro del monstruo.
func (m *ZoneMonster) ResetAggro() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Aggro = make(map[uint16]int)
	m.TargetHandle = 0
}

// TopDamageChar retorna el CharNo que más daño aportó o el último atacante.
func (m *ZoneMonster) TopDamageChar() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	bestChar := m.LastKillerChar
	bestDmg := -1
	for cid, dmg := range m.DamageContrib {
		if dmg > bestDmg {
			bestDmg = dmg
			bestChar = cid
		}
	}
	return bestChar
}

// BriefInfoMaxBatch define el límite máximo oficial de entradas por paquete de lote BriefInfo en Zone.exe (0x1C = 28).
const BriefInfoMaxBatch = 28

// encodeCoordDir convierte el ángulo en grados (0..360 / negativo) al byte de dirección oficial (direct / 2, +180 si negativo),
// verificado en ShineMob::so_mob_Regenerate (Zone.exe VA 0x004B4279).
func encodeCoordDir(direct int16) byte {
	d := direct / 2
	if d < 0 {
		d += 180
	}
	return byte(d)
}

// WriteRegenMob serializa la estructura oficial PROTO_NC_BRIEFINFO_REGENMOB_CMD (149 bytes exactos en PDB).
// Alineado con Zone.exe ShineMob::so_mob_Regenerate (0x004B41B0): mode=0 (vivo/idle), dir=direct/2.
func (m *ZoneMonster) WriteRegenMob(pw *protocol.PacketWriter, regenAni byte) {
	m.mu.RLock()
	handle := m.Handle
	mobID := m.MobID
	x := m.X
	y := m.Y
	dir := m.Direct
	isDead := m.IsDead
	anim := m.Animation
	animLv := m.AnimationLevel
	kqTeam := m.KQTeamType
	var abBit [99]byte
	for stateID := range m.ActiveStates {
		byteIdx := int(stateID) / 8
		bitIdx := uint(stateID) % 8
		if byteIdx >= 0 && byteIdx < len(abBit) {
			abBit[byteIdx] |= 1 << bitIdx
		}
	}
	m.mu.RUnlock()

	var mode byte
	if isDead {
		mode = 3
	}

	pw.WriteUint16(handle)            // 0x00: handle (2B)
	pw.WriteByte(mode)                // 0x02: mode (1B, 0=vivo/idle en ShineMob::so_mob_Regenerate)
	pw.WriteUint16(mobID)             // 0x03: mobid (2B)
	pw.WriteInt32(x)                  // 0x05: coord.x (4B)
	pw.WriteInt32(y)                  // 0x09: coord.y (4B)
	pw.WriteByte(encodeCoordDir(dir)) // 0x0D: coord.dir (1B, direct/2)
	pw.WriteByte(0)                   // 0x0E: flagstate (1B, 0=mob abstatebit)
	pw.WriteBytes(abBit[:])           // 0x0F: ABNORMAL_STATE_BIT (99B)
	pw.WriteString(anim, 32)          // 0x72: sAnimation (32B)
	pw.WriteByte(animLv)              // 0x92: nAnimationLevel (1B)
	pw.WriteByte(kqTeam)              // 0x93: nKQTeamType (1B)
	pw.WriteByte(regenAni)            // 0x94: bRegenAni (1B)
}
