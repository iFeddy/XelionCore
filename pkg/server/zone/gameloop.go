package zone

import (
	"math"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// MapGameLoop gobierna los ciclos determinísticos del mapa (100ms, 1s, 30s) según la arquitectura oficial.
type MapGameLoop struct {
	zMap     *ZoneMap
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewMapGameLoop instancia un bucle de juego para el mapa indicado.
func NewMapGameLoop(zMap *ZoneMap) *MapGameLoop {
	return &MapGameLoop{
		zMap:     zMap,
		stopChan: make(chan struct{}),
	}
}

// Start inicia las rutinas periódicas de simulación del mapa.
func (gl *MapGameLoop) Start() {
	gl.mu.Lock()
	if gl.running {
		gl.mu.Unlock()
		return
	}
	gl.running = true
	gl.mu.Unlock()

	gl.wg.Add(1)
	go gl.run()
}

// Stop finaliza limpiamente el bucle de juego del mapa.
func (gl *MapGameLoop) Stop() {
	gl.mu.Lock()
	if !gl.running {
		gl.mu.Unlock()
		return
	}
	gl.running = false
	close(gl.stopChan)
	gl.mu.Unlock()

	gl.wg.Wait()
}

func (gl *MapGameLoop) run() {
	defer gl.wg.Done()

	tick100ms := time.NewTicker(100 * time.Millisecond)
	defer tick100ms.Stop()

	tick1s := time.NewTicker(1 * time.Second)
	defer tick1s.Stop()

	tick30s := time.NewTicker(30 * time.Second)
	defer tick30s.Stop()

	for {
		select {
		case <-gl.stopChan:
			return
		case <-tick100ms.C:
			gl.tick100ms()
		case <-tick1s.C:
			gl.tick1s()
		case <-tick30s.C:
			gl.tick30s()
		}
	}
}

// tick100ms ejecuta la lógica de alta frecuencia (cadencia de auto-ataque y habilidades).
func (gl *MapGameLoop) tick100ms() {
	gl.zMap.mu.RLock()
	players := make([]*ZonePlayer, 0, len(gl.zMap.players))
	for _, p := range gl.zMap.players {
		players = append(players, p)
	}
	gl.zMap.mu.RUnlock()

	now := time.Now()
	for _, p := range players {
		if p.IsAutoAttacking && p.TargetHandle != 0 && p.CanAttack() {
			if now.Sub(p.LastAttackTime) >= 1000*time.Millisecond {
				gl.executeAutoAttack(p, now)
			}
		}
	}
}

func (gl *MapGameLoop) executeAutoAttack(player *ZonePlayer, now time.Time) {
	if player == nil || gl.zMap == nil || player.TargetHandle == 0 || !player.IsAutoAttacking || !player.CanAttack() {
		return
	}

	targetEntity := gl.zMap.GetEntity(player.TargetHandle)
	if targetEntity == nil {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
		return
	}

	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	dist := math.Hypot(float64(px-tx), float64(py-ty))
	if dist > 250.0 {
		return
	}

	if gl.zMap.blockInfo != nil && gl.zMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(tx), float64(ty)) {
		return
	}

	player.LastAttackTime = now

	// 1. Notificar inicio de animación
	swingStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_START_CMD)
	swingStart.WriteUint16(player.MapObjectID)
	swingStart.WriteUint16(player.TargetHandle)
	swingStart.WriteUint8(1)
	swingStart.WriteUint16(1000)
	swingStart.WriteUint8(1)
	swingStart.WriteUint8(1)
	gl.zMap.Broadcast(swingStart, 0)

	// 2. Verificar hit/miss
	hit := CheckHit(int(player.EquipTH)+30, 10)
	if !hit {
		return
	}

	isCrit := CheckCrit(15, 5)
	angle := CalculateAngleDeg(float64(px), float64(py), float64(tx), float64(ty), 0)
	dmg := CalculatePhysicalDamage(int(player.EquipMaxWC)+15, 5, angle, true, 0, 1000, isCrit, nil)

	var restHP uint32
	var isDead bool

	switch e := targetEntity.(type) {
	case *ZoneMonster:
		if e.CurHP <= dmg {
			e.CurHP = 0
			e.IsDead = true
			isDead = true
		} else {
			e.CurHP -= dmg
		}
		restHP = uint32(e.CurHP)
	case *ZonePlayer:
		if e.Char.CurHP <= dmg {
			e.Char.CurHP = 0
			e.State = 3
			isDead = true
		} else {
			e.Char.CurHP -= dmg
		}
		restHP = uint32(e.Char.CurHP)
		if e.Client != nil {
			hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpCmd.WriteUint32(restHP)
			hpCmd.WriteUint16(0)
			_ = e.Client.SendPacket(hpCmd)
		}
	}

	swingDmg := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_DAMAGE_CMD)
	swingDmg.WriteUint16(player.MapObjectID)
	swingDmg.WriteUint16(player.TargetHandle)
	flag := uint16(0)
	if isCrit {
		flag = 1
	}
	swingDmg.WriteUint16(flag)
	swingDmg.WriteUint16(uint16(dmg))
	swingDmg.WriteUint32(restHP)
	swingDmg.WriteUint16(0)
	swingDmg.WriteUint8(1)
	swingDmg.WriteUint8(1)
	gl.zMap.Broadcast(swingDmg, 0)

	if isDead {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
		deadCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONEDEAD_LARGE_CMD)
		deadCmd.WriteUint16(targetEntity.GetHandle())
		gl.zMap.Broadcast(deadCmd, 0)
	}
}

// tick1s ejecuta la regeneración natural y en tienda (Tent) de HP y SP, DoTs y ticks de estados.
func (gl *MapGameLoop) tick1s() {
	gl.zMap.mu.RLock()
	players := make([]*ZonePlayer, 0, len(gl.zMap.players))
	for _, p := range gl.zMap.players {
		players = append(players, p)
	}
	gl.zMap.mu.RUnlock()

	now := time.Now()
	for _, p := range players {
		if p.Char == nil || p.Char.CurHP <= 0 {
			continue
		}

		hpChanged := false
		spChanged := false

		// Decrementar PKYellowTime si corresponde
		if p.PKYellowTime > 0 {
			p.PKYellowTime--
		}

		// Procesar ticks de DoT / HoT y expiración de estados anormales
		for id, st := range p.ActiveStates {
			if st.DotDamage > 0 && (st.LastDotTick.IsZero() || now.Sub(st.LastDotTick) >= st.DotInterval) {
				st.LastDotTick = now
				if p.Char.CurHP <= st.DotDamage {
					p.Char.CurHP = 0
					p.State = 3 // Dead
				} else {
					p.Char.CurHP -= st.DotDamage
				}
				hpChanged = true

				dotCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_DOTDAMAGE_CMD)
				dotCmd.WriteUint16(p.MapObjectID)
				dotCmd.WriteUint32(uint32(p.Char.CurHP))
				dotCmd.WriteUint16(uint16(st.DotDamage))
				dotCmd.WriteUint16(st.ID)
				dotCmd.WriteUint16(0)
				dotCmd.WriteUint8(0) // hit
				gl.zMap.Broadcast(dotCmd, 0)
			}

			if now.After(st.ExpiresAt) {
				p.RemoveAbState(id)
				resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_RESETABSTATE_CMD)
				resetCmd.WriteUint16(id)
				if p.Client != nil {
					_ = p.Client.SendPacket(resetCmd)
				}
			}
		}

		// Estado 4 = Resting en Tent: regeneración rápida de HP y SP (10% por segundo)
		if p.State == 4 {
			regenHP := p.MaxHP / 10
			if regenHP < 5 {
				regenHP = 5
			}
			if p.Char.CurHP < p.MaxHP {
				p.Char.CurHP += regenHP
				if p.Char.CurHP > p.MaxHP {
					p.Char.CurHP = p.MaxHP
				}
				hpChanged = true
			}

			regenSP := p.MaxSP / 10
			if regenSP < 5 {
				regenSP = 5
			}
			if p.Char.CurSP < p.MaxSP {
				p.Char.CurSP += regenSP
				if p.Char.CurSP > p.MaxSP {
					p.Char.CurSP = p.MaxSP
				}
				spChanged = true
			}
		}

		if hpChanged && p.Client != nil {
			hpUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpUpdate.WriteUint32(uint32(p.Char.CurHP))
			hpUpdate.WriteUint16(0)
			_ = p.Client.SendPacket(hpUpdate)
		}

		if spChanged && p.Client != nil {
			spUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
			spUpdate.WriteUint32(uint32(p.Char.CurSP))
			_ = p.Client.SendPacket(spUpdate)
		}
	}
}

// tick30s limpia ítems del suelo expirados y estados periódicos lentos.
func (gl *MapGameLoop) tick30s() {
	gl.zMap.CleanupExpiredDrops(2 * time.Minute)
}
