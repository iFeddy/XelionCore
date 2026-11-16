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

// Step1s ejecuta un ciclo de 1 segundo (DoT/HoT, expiración de AbState y regeneración pasiva).
func (gl *MapGameLoop) Step1s() {
	gl.tick1s()
}

// Step30s ejecuta un ciclo de 30 segundos (expiración de drops, pets y CHARSAVE cada 60s).
func (gl *MapGameLoop) Step30s() {
	gl.tick30s()
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
	gl.zMap.TickMonsters100ms(now)

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
		var curHP uint32
		switch e := targetEntity.(type) {
		case *ZoneMonster:
			curHP = uint32(e.CurHP)
		case *ZonePlayer:
			if e.Char != nil {
				curHP = uint32(e.Char.CurHP)
			}
		}
		emitSwingDamage(gl.zMap, player, targetEntity, player.MapObjectID, player.TargetHandle, 0, curHP, false, true, false, 1)
		return
	}

	isCrit := CheckCrit(15, 5)
	angle := CalculateAngleDeg(float64(px), float64(py), float64(tx), float64(ty), 0)
	dmg := CalculatePhysicalDamage(int(player.EquipMaxWC)+15, 5, angle, true, 0, 1000, isCrit, nil)

	var restHP uint32
	var isDead bool

	switch e := targetEntity.(type) {
	case *ZoneMonster:
		restHP, isDead = gl.zMap.ApplyDamageToMonster(e, player, dmg)
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

	emitSwingDamage(gl.zMap, player, targetEntity, player.MapObjectID, player.TargetHandle, uint16(dmg), restHP, isCrit, false, isDead, 1)

	if isDead {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
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
	gl.zMap.TickMonsters1s(now)
	gl.zMap.TickInstanceLua(now)
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
				killedByDot := false
				if p.Char.CurHP <= st.DotDamage {
					p.Char.CurHP = 0
					p.State = 3 // Dead
					p.IsAutoAttacking = false
					p.TargetHandle = 0
					killedByDot = true
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

				if killedByDot {
					emitReallyKill(gl.zMap, p.MapObjectID, 0)
				}
			} else if st.DotDamage < 0 && (st.LastDotTick.IsZero() || now.Sub(st.LastDotTick) >= st.DotInterval) {
				st.LastDotTick = now
				if p.Char.CurHP > 0 && p.Char.CurHP < p.MaxHP {
					p.Char.CurHP += -st.DotDamage
					if p.Char.CurHP > p.MaxHP {
						p.Char.CurHP = p.MaxHP
					}
					hpChanged = true
				}
			}

			if !st.ExpiresAt.IsZero() && now.After(st.ExpiresAt) {
				p.RemoveAbState(id)
				if p.Client != nil {
					informCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD)
					informCmd.WriteUint32(uint32(id))
					informCmd.WriteUint32(0)
					_ = p.Client.SendPacket(informCmd)
				}
				// Verificado en Zone.exe so_AbnormalState_BroadcastReset (0x00405800):
				// broadcast PROTO_NC_BAT_ABSTATERESET_CMD (0x2428, Handle[2] + Abstate[4]).
				resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATERESET_CMD)
				resetCmd.WriteUint16(p.MapObjectID)
				resetCmd.WriteUint32(uint32(id))
				gl.zMap.Broadcast(resetCmd, 0)
			}
		}

		// Regeneración oficial por estado (Zone.exe sp_RecoverPoint_House / sp_RecoverPoint_Rest):
		// Estado 4 = Tent (10%/s); Estado 6 = Rest/Sit (cada 3s); Estado 1 = Normal pasivo (cada 3s).
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
		} else if p.State == 6 || p.State == 1 {
			p.RestTickCount++
			if p.RestTickCount >= 3 {
				p.RestTickCount = 0
				var regenHP, regenSP int
				if p.State == 6 {
					regenHP = (p.MaxHP * 3) / 100
					if regenHP < 5 {
						regenHP = 5
					}
					regenSP = (p.MaxSP * 5) / 100
					if regenSP < 5 {
						regenSP = 5
					}
				} else {
					regenHP = p.MaxHP / 100
					if regenHP < 1 {
						regenHP = 1
					}
					regenSP = p.MaxSP / 50
					if regenSP < 1 {
						regenSP = 1
					}
				}
				if p.Char.CurHP > 0 && p.Char.CurHP < p.MaxHP {
					p.Char.CurHP += regenHP
					if p.Char.CurHP > p.MaxHP {
						p.Char.CurHP = p.MaxHP
					}
					hpChanged = true
				}
				if p.Char.CurSP < p.MaxSP {
					p.Char.CurSP += regenSP
					if p.Char.CurSP > p.MaxSP {
						p.Char.CurSP = p.MaxSP
					}
					spChanged = true
				}
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

// tick30s limpia ítems del suelo expirados, ejecuta auto-guardado periódico CHARSAVE (60s ws_save2DB) y estados lentos (Pets GTI).
func (gl *MapGameLoop) tick30s() {
	gl.zMap.CleanupExpiredDrops(2 * time.Minute)
	now := time.Now()
	if gl.zMap.handler != nil {
		for _, p := range gl.zMap.GetPlayers() {
			if p == nil {
				continue
			}
			if p.ActivePet != nil {
				gl.zMap.handler.TickActivePet(p)
			}
			// Auto-guardado periódico CHARSAVE (Zone.exe ws_save2DB 0x00455A10: cada 0x258=600 ticks = 60s)
			if p.Char != nil && gl.zMap.handler.charRepo != nil {
				if p.LastSaveAt.IsZero() {
					p.LastSaveAt = now
				} else if now.Sub(p.LastSaveAt) >= 60*time.Second {
					px, py := p.GetPosition()
					p.Char.XPos = int(px)
					p.Char.YPos = int(py)
					if gl.zMap.MapID != "" {
						p.Char.TownName = gl.zMap.MapID
					}
					p.Char.TCharacterState.XPos = p.Char.XPos
					p.Char.TCharacterState.YPos = p.Char.YPos
					p.Char.TCharacterState.HP = p.Char.CurHP
					p.Char.TCharacterState.SP = p.Char.CurSP
					p.Char.TCharacterState.HPStones = p.StonesHP
					p.Char.TCharacterState.SPStones = p.StonesSP
					p.Char.TCharacterState.Rotation = p.Rotation
					_ = gl.zMap.handler.charRepo.SaveCharLocation(p.Char.ID, p.Char.TownName, p.Char.XPos, p.Char.YPos, p.Rotation)
					_ = gl.zMap.handler.charRepo.SaveCharStats(p.Char.ID, int16(p.Char.CurHP), int16(p.Char.CurSP), p.StonesHP, p.StonesSP)
					_ = gl.zMap.handler.charRepo.UpdateCharacter(p.Char)
					p.LastSaveAt = now
				}
			}
		}
	}
}
