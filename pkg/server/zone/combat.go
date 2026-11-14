package zone

import (
	"math"
	"math/rand"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// CalculateAngleDeg calcula la diferencia angular en grados (0..180) entre el encaramiento del defensor y la posición del atacante.
// 0 grados = ataque por la espalda, 90 grados = ataque por el flanco, 180 grados = ataque frontal.
func CalculateAngleDeg(attX, attY, defX, defY float64, defDirect int16) int {
	// En Fiesta, Direct va de 0 a 360 o 0 a 255. Normalizamos a radianes.
	defAngleRad := float64(defDirect) * math.Pi / 180.0
	dx := attX - defX
	dy := attY - defY
	if dx == 0 && dy == 0 {
		return 180
	}
	attackAngleRad := math.Atan2(dy, dx)
	diffRad := math.Abs(attackAngleRad - defAngleRad)
	for diffRad > math.Pi {
		diffRad = 2*math.Pi - diffRad
	}
	return int(diffRad * 180.0 / math.Pi)
}

// CalculatePhysicalDamage evalúa el daño físico usando fórmulas de combate oficiales gobernadas por data/.
func CalculatePhysicalDamage(attackerAtk, defenderDef, angleDeg int, isMonster bool, lvGap int, soulRate int, isCrit bool, combatStore interface{ GetAngleDamageRate(bool, int) int }) int {
	baseDmg := attackerAtk - defenderDef
	if baseDmg < 1 {
		baseDmg = 1
	}

	angleMult := 1000
	if combatStore != nil {
		angleMult = combatStore.GetAngleDamageRate(isMonster, angleDeg)
	}

	// Multiplicador por diferencia de nivel (DamageLvGap)
	lvGapMult := 1000
	if lvGap < -5 {
		lvGapMult = 1500
	} else if lvGap < 0 {
		lvGapMult = 1000 + (-lvGap)*100
	}

	if soulRate <= 0 {
		soulRate = 1000
	}

	critMult := 1000
	if isCrit {
		critMult = 1500
	}

	total := float64(baseDmg) * (float64(angleMult) / 1000.0) * (float64(lvGapMult) / 1000.0) * (float64(soulRate) / 1000.0) * (float64(critMult) / 1000.0)
	if total < 1.0 {
		return 1
	}
	return int(math.Round(total))
}

// CheckHit evalúa precisión (AimRate) frente a evasión (EvasionRate).
func CheckHit(aimRate, evasionRate int) bool {
	if aimRate <= 0 {
		aimRate = 10
	}
	if evasionRate <= 0 {
		evasionRate = 5
	}
	chance := float64(aimRate) / float64(aimRate+evasionRate)
	return rand.Float64() <= chance
}

// CheckCrit evalúa probabilidad de golpe crítico.
func CheckCrit(critRate, critResist int) bool {
	if critRate <= 0 {
		return false
	}
	rate := float64(critRate-critResist) / 100.0
	if rate < 0.05 {
		rate = 0.05
	}
	if rate > 0.80 {
		rate = 0.80
	}
	return rand.Float64() <= rate
}

// writeTargetInfo construye PROTO_NC_BAT_TARGETINFO_CMD (0x2402, 30B payload = 32B wire).
// Verificado en Zone.exe to_SendInform (0x0042B3F0): order=0 para objetivo directo, order=1 para objetivo del objetivo.
func buildTargetInfoPacket(order uint8, targetHandle uint16, entity GridEntity) *protocol.PacketWriter {
	info := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETINFO_CMD)
	info.WriteUint8(order)
	info.WriteUint16(targetHandle)

	switch e := entity.(type) {
	case *ZonePlayer:
		var hp, sp uint32
		var lv uint8 = 1
		if e.Char != nil {
			hp = uint32(e.Char.CurHP)
			sp = uint32(e.Char.CurSP)
			lv = uint8(e.Char.Level)
		}
		info.WriteUint32(hp)
		info.WriteUint32(uint32(e.MaxHP))
		info.WriteUint32(sp)
		info.WriteUint32(uint32(e.MaxSP))
		info.WriteUint32(0) // targetlp
		info.WriteUint32(0) // targetmaxlp
		info.WriteUint8(lv)
	case *ZoneMonster:
		e.mu.RLock()
		hp := uint32(e.CurHP)
		maxHP := uint32(e.MaxHP)
		sp := uint32(e.CurSP)
		maxSP := uint32(e.MaxSP)
		lv := e.Level
		e.mu.RUnlock()
		info.WriteUint32(hp)
		info.WriteUint32(maxHP)
		info.WriteUint32(sp)
		info.WriteUint32(maxSP)
		info.WriteUint32(0) // targetlp
		info.WriteUint32(0) // targetmaxlp
		info.WriteUint8(lv)
	default:
		info.WriteUint32(100)
		info.WriteUint32(100)
		info.WriteUint32(50)
		info.WriteUint32(50)
		info.WriteUint32(0)
		info.WriteUint32(0)
		info.WriteUint8(1)
	}

	info.WriteUint16(0) // hpchangeorder (vtable+0x7e8)
	return info
}

// sendTargetInfoWithSecondary envía TARGETINFO_CMD (order=0) y, si el objetivo tiene su propio objetivo activo,
// envía el TARGETINFO_CMD secundario (order=1) alineado con to_SendInform (0x0042B562).
func sendTargetInfoWithSecondary(client *network.Client, zMap *ZoneMap, targetHandle uint16) {
	if client == nil || zMap == nil || targetHandle == 0 {
		return
	}
	entity := zMap.GetEntity(targetHandle)
	if entity == nil {
		return
	}
	_ = client.SendPacket(buildTargetInfoPacket(0, targetHandle, entity))

	var secondaryHandle uint16
	switch e := entity.(type) {
	case *ZonePlayer:
		secondaryHandle = e.TargetHandle
	case *ZoneMonster:
		e.mu.RLock()
		secondaryHandle = e.TargetHandle
		e.mu.RUnlock()
	}
	if secondaryHandle != 0 {
		if secEntity := zMap.GetEntity(secondaryHandle); secEntity != nil {
			_ = client.SendPacket(buildTargetInfoPacket(1, secondaryHandle, secEntity))
		}
	}
}

// handleTargettingReq implementa PROTO_NC_BAT_TARGETTING_REQ (0x2401).
// Verificado en Zone.exe sp_NC_BAT_TARGETTING_REQ (0x0042FFE0).
func (h *Handler) handleTargettingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	targetHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	entity := player.CurrentMap.GetEntity(targetHandle)
	if entity == nil {
		return
	}

	player.TargetHandle = targetHandle
	sendTargetInfoWithSecondary(client, player.CurrentMap, targetHandle)
}

// handleUntargetReq implementa PROTO_NC_BAT_UNTARGET_REQ (0x2408).
func (h *Handler) handleUntargetReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.TargetHandle = 0
	player.IsAutoAttacking = false
}

// handleBashStartCmd implementa PROTO_NC_BAT_BASHSTART_CMD (0x242B).
func (h *Handler) handleBashStartCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.TargetHandle == 0 || !player.CanAttack() {
		return
	}
	player.IsAutoAttacking = true
}

// handleBashStopCmd implementa PROTO_NC_BAT_BASHSTOP_CMD (0x2432).
func (h *Handler) handleBashStopCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.IsAutoAttacking = false
}

// handleAssistReq implementa PROTO_NC_BAT_ASSIST_REQ (0x243E).
// Verificado en Zone.exe sp_NC_BAT_ASSIST_REQ (0x0042B7A0): responde PROTO_NC_BAT_ASSIST_ACK (0x243F, 2B target).
func (h *Handler) handleAssistReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	tankerHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	entity := player.CurrentMap.GetEntity(tankerHandle)
	if entity == nil {
		return
	}

	var assistedTarget uint16 = 0xFFFF
	switch t := entity.(type) {
	case *ZonePlayer:
		if t.TargetHandle != 0 {
			assistedTarget = t.TargetHandle
			player.TargetHandle = assistedTarget
		}
	case *ZoneMonster:
		t.mu.RLock()
		if t.TargetHandle != 0 {
			assistedTarget = t.TargetHandle
			player.TargetHandle = assistedTarget
		}
		t.mu.RUnlock()
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ASSIST_ACK)
	ack.WriteUint16(assistedTarget)
	_ = client.SendPacket(ack)
}

// skillDamageTarget encapsula un objetivo impactado para PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD (0x2452).
type skillDamageTarget struct {
	Handle   uint16
	Damage   uint32
	RestHP   uint32
	IsCrit   bool
	IsMissed bool
	IsDead   bool
}

// emitSkillBashHitDamage construye y difunde PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD (0x2452, 5 + 14*N bytes).
func emitSkillBashHitDamage(zMap *ZoneMap, index, caster uint16, targets []skillDamageTarget) {
	if zMap == nil || len(targets) == 0 {
		return
	}
	// En Zone.exe PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD envía targetnum (uint8) + SkillDamage[targetnum].
	if len(targets) > 255 {
		targets = targets[:255]
	}
	pkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD)
	pkt.WriteUint16(index)
	pkt.WriteUint16(caster)
	pkt.WriteUint8(uint8(len(targets)))
	for _, t := range targets {
		pkt.WriteUint16(t.Handle)
		var flag uint16 = 0x0001 // isdamage
		if t.IsCrit {
			flag |= 0x0002 // iscritical
		}
		if t.IsMissed {
			flag &^= 0x0001
			flag |= 0x0004 // ismissed
		}
		if t.IsDead {
			flag |= 0x0100 // isDead
		}
		pkt.WriteUint16(flag)
		pkt.WriteUint32(t.Damage)
		pkt.WriteUint32(t.RestHP)
		pkt.WriteUint16(0) // hpchangeorder
	}
	zMap.Broadcast(pkt, 0)
}

// emitReallyKill difunde PROTO_NC_BAT_REALLYKILL_CMD (0x244A, 4B payload = 6B wire).
// Verificado en Zone.exe so_mob_SendReallyKill (0x00553760).
func emitReallyKill(zMap *ZoneMap, deadHandle, attackerHandle uint16) {
	if zMap == nil || deadHandle == 0 {
		return
	}
	killCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_REALLYKILL_CMD)
	killCmd.WriteUint16(deadHandle)
	killCmd.WriteUint16(attackerHandle)
	zMap.Broadcast(killCmd, 0)
}

// emitSwingDamage emite PROTO_NC_BAT_SWING_DAMAGE_CMD (0x2448, 16B payload) al atacante y defensor,
// y PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD (0x2449, 13B payload) al resto de jugadores en el mapa.
// Verificado 1:1 en Zone.exe smo_SwingDamage (0x004334FD..0x00433810).
func emitSwingDamage(zMap *ZoneMap, attacker, defender GridEntity, attackerHandle, defenderHandle uint16, dmg uint16, restHP uint32, isCrit, isMissed, isDead bool, atkSeq uint8) {
	if zMap == nil {
		return
	}

	var dmgFlag uint16
	if isCrit {
		dmgFlag |= 0x0001 // iscritical
	}
	if isMissed {
		dmgFlag |= 0x0004 // ismissed
	}
	if isDead {
		dmgFlag |= 0x0020 // isDead
	}

	swingDmg := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_DAMAGE_CMD)
	swingDmg.WriteUint16(attackerHandle)
	swingDmg.WriteUint16(defenderHandle)
	swingDmg.WriteUint16(dmgFlag)
	swingDmg.WriteUint16(dmg)
	swingDmg.WriteUint32(restHP)
	swingDmg.WriteUint16(0)      // hpchangeorder
	swingDmg.WriteUint8(1)       // damageindex
	swingDmg.WriteUint8(atkSeq)  // attacksequence

	if atkP, ok := attacker.(*ZonePlayer); ok && atkP != nil && atkP.Client != nil {
		_ = atkP.Client.SendPacket(swingDmg)
	}
	if defP, ok := defender.(*ZonePlayer); ok && defP != nil && defP.Client != nil && defP.MapObjectID != attackerHandle {
		_ = defP.Client.SendPacket(swingDmg)
	}

	var someoneFlag uint8
	if isMissed {
		someoneFlag |= 0x01 // bit 0: isMissed
	}
	someoneDmg := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESWING_DAMAGE_CMD)
	someoneDmg.WriteUint16(attackerHandle)
	someoneDmg.WriteUint16(defenderHandle)
	someoneDmg.WriteUint8(someoneFlag)
	someoneDmg.WriteUint32(restHP)
	someoneDmg.WriteUint16(0)      // hpchangeorder
	someoneDmg.WriteUint8(1)       // damageindex
	someoneDmg.WriteUint8(atkSeq)  // attacksequence

	zMap.mu.RLock()
	for id, p := range zMap.players {
		if p == nil || p.Client == nil || p.HideExceptMe {
			continue
		}
		if id == attackerHandle || id == defenderHandle {
			continue
		}
		_ = p.Client.SendPacket(someoneDmg)
	}
	zMap.mu.RUnlock()

	if isDead {
		emitReallyKill(zMap, defenderHandle, attackerHandle)
	}
}

// handleSkillBashObjCastReq implementa PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x2440).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x0057FA50) + smo_SkillBash_Casting (0x0057BEF0).
func (h *Handler) handleSkillBashObjCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || !player.CanCast() {
		return
	}

	skillID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	targetObj, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	// Validar objetivo en mapa
	targetEntity := player.CurrentMap.GetEntity(targetObj)
	if targetEntity == nil {
		return
	}

	// Validar línea de visión (LoS) con .shbd
	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(tx), float64(ty)) {
		return
	}

	player.CastingSkillID = skillID

	// Costo de SP base si tiene maná suficiente
	spCost := 10
	if player.Char.CurSP >= spCost {
		player.Char.CurSP -= spCost
		spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
		spCmd.WriteUint32(uint32(player.Char.CurSP))
		_ = client.SendPacket(spCmd)
	}

	const castIndex uint16 = 1

	// Lanzador recibe PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD (0x244E, 6B payload = 8B wire)
	startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD)
	startCmd.WriteUint16(skillID)
	startCmd.WriteUint16(targetObj)
	startCmd.WriteUint16(castIndex)
	_ = client.SendPacket(startCmd)

	// Circundantes reciben PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD (0x244F, 8B payload = 10B wire)
	someoneStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD)
	someoneStart.WriteUint16(player.MapObjectID)
	someoneStart.WriteUint16(skillID)
	someoneStart.WriteUint16(targetObj)
	someoneStart.WriteUint16(castIndex)
	player.CurrentMap.Broadcast(someoneStart, player.MapObjectID)

	// Calcular daño de habilidad
	dmg := 35
	var restHP uint32
	var isDead bool
	if mob, ok := targetEntity.(*ZoneMonster); ok {
		restHP, isDead = player.CurrentMap.ApplyDamageToMonster(mob, player, dmg)
	} else if defP, ok := targetEntity.(*ZonePlayer); ok {
		if defP.Char.CurHP <= dmg {
			defP.Char.CurHP = 0
			defP.State = 3
			isDead = true
		} else {
			defP.Char.CurHP -= dmg
		}
		restHP = uint32(defP.Char.CurHP)
		if defP.Client != nil {
			hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpCmd.WriteUint32(restHP)
			hpCmd.WriteUint16(0)
			_ = defP.Client.SendPacket(hpCmd)
		}
	}

	emitSkillBashHitDamage(player.CurrentMap, castIndex, player.MapObjectID, []skillDamageTarget{{
		Handle: targetObj,
		Damage: uint32(dmg),
		RestHP: restHP,
		IsDead: isDead,
	}})

	if isDead {
		emitReallyKill(player.CurrentMap, targetObj, player.MapObjectID)
	}
	player.CastingSkillID = 0
}

// handleSkillBashFldCastReq implementa PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x2441).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x0057F400).
func (h *Handler) handleSkillBashFldCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || !player.CanCast() {
		return
	}

	skillID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	fldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	fldY, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}

	player.CastingSkillID = skillID
	const castIndex uint16 = 1

	// Lanzador recibe PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD (0x2450, 12B payload)
	fldCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD)
	fldCmd.WriteUint16(skillID)
	fldCmd.WriteInt32(fldX)
	fldCmd.WriteInt32(fldY)
	fldCmd.WriteUint16(castIndex)
	_ = client.SendPacket(fldCmd)

	// Circundantes reciben PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD (0x2451, 14B payload)
	someoneFld := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD)
	someoneFld.WriteUint16(player.MapObjectID)
	someoneFld.WriteUint16(skillID)
	someoneFld.WriteInt32(fldX)
	someoneFld.WriteInt32(fldY)
	someoneFld.WriteUint16(castIndex)
	player.CurrentMap.Broadcast(someoneFld, player.MapObjectID)

	// Aplicar daño a entidades dentro del radio AoE usando distancia euclidiana
	radius := 200.0
	const aoeDmg = 25
	var hitTargets []skillDamageTarget
	var deadHandles []uint16
	for _, mob := range player.CurrentMap.GetAllMonsters() {
		mx, my := mob.GetPosition()
		dist := math.Hypot(float64(mx)-float64(fldX), float64(my)-float64(fldY))
		if dist <= radius && !mob.IsDead {
			restHP, isDead := player.CurrentMap.ApplyDamageToMonster(mob, player, aoeDmg)
			hitTargets = append(hitTargets, skillDamageTarget{
				Handle: mob.Handle,
				Damage: aoeDmg,
				RestHP: restHP,
				IsDead: isDead,
			})
			if isDead {
				deadHandles = append(deadHandles, mob.Handle)
			}
		}
	}

	emitSkillBashHitDamage(player.CurrentMap, castIndex, player.MapObjectID, hitTargets)
	for _, deadH := range deadHandles {
		emitReallyKill(player.CurrentMap, deadH, player.MapObjectID)
	}
	player.CastingSkillID = 0
}

// handleSkillBashCastAbortReq implementa PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ (0x2444).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_CASTABORT_REQ (0x0057C480).
func (h *Handler) handleSkillBashCastAbortReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	var errCode uint16 = 0x20EA // no está casteando / modo inválido
	if player.State == 6 {
		errCode = 0x20E9
	} else if player.CastingSkillID != 0 {
		errCode = 0x20E8 // abortado con éxito
		player.CastingSkillID = 0
		if player.CurrentMap != nil {
			cutCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD)
			cutCmd.WriteUint16(player.MapObjectID)
			player.CurrentMap.Broadcast(cutCmd, player.MapObjectID)
		}
	} else {
		// Permitir también abortar limpiamente si el cliente solicita abortar tras casteo reciente
		errCode = 0x20E8
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleMoverSkillBashObjCastReq implementa PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ (0xCC0E).
func (h *Handler) handleMoverSkillBashObjCastReq(client *network.Client, pkt *protocol.Packet) {
	h.handleSkillBashObjCastReq(client, pkt)
}

// handleMoverSkillBashFldCastReq implementa PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ (0xCC0F).
func (h *Handler) handleMoverSkillBashFldCastReq(client *network.Client, pkt *protocol.Packet) {
	h.handleSkillBashFldCastReq(client, pkt)
}

// ExecuteAutoAttack realiza un golpe del bucle de auto-ataque contra el objetivo actual.
func (h *Handler) ExecuteAutoAttack(player *ZonePlayer, now time.Time) {
	if player == nil || player.CurrentMap == nil || player.TargetHandle == 0 || !player.IsAutoAttacking || !player.CanAttack() {
		return
	}

	targetEntity := player.CurrentMap.GetEntity(player.TargetHandle)
	if targetEntity == nil {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
		return
	}

	// Validar distancia de ataque (rango cuerpo a cuerpo: ~150 unidades)
	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	dist := math.Hypot(float64(px-tx), float64(py-ty))
	if dist > 250.0 {
		return // Fuera de rango
	}

	// Validar línea de visión (LoS)
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(tx), float64(ty)) {
		return // Ataque bloqueado por obstáculo
	}

	player.LastAttackTime = now

	// 1. Notificar inicio de animación de swing: PROTO_NC_BAT_SWING_START_CMD (0x2447, 9B payload)
	swingStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_START_CMD)
	swingStart.WriteUint16(player.MapObjectID)
	swingStart.WriteUint16(player.TargetHandle)
	swingStart.WriteUint8(1)     // actioncode
	swingStart.WriteUint16(1000) // attackspeed
	swingStart.WriteUint8(1)     // damageindex
	swingStart.WriteUint8(1)     // attacksequence
	player.CurrentMap.Broadcast(swingStart, 0)

	// 2. Verificar impacto (Hit / Miss)
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
		emitSwingDamage(player.CurrentMap, player, targetEntity, player.MapObjectID, player.TargetHandle, 0, curHP, false, true, false, 1)
		return
	}

	// 3. Evaluar crítico y daño
	isCrit := CheckCrit(15, 5)
	angle := CalculateAngleDeg(float64(px), float64(py), float64(tx), float64(ty), 0)
	dmg := CalculatePhysicalDamage(int(player.EquipMaxWC)+15, 5, angle, true, 0, 1000, isCrit, h.dataStore.CombatStore)

	var restHP uint32
	var isDead bool

	switch e := targetEntity.(type) {
	case *ZoneMonster:
		restHP, isDead = player.CurrentMap.ApplyDamageToMonster(e, player, dmg)
	case *ZonePlayer:
		if e.Char.CurHP <= dmg {
			e.Char.CurHP = 0
			e.State = 3 // Dead
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

	// 4. Notificar impacto con daño (0x2448 atacante/defensor + 0x2449 entorno + 0x244A si muere)
	emitSwingDamage(player.CurrentMap, player, targetEntity, player.MapObjectID, player.TargetHandle, uint16(dmg), restHP, isCrit, false, isDead, 1)

	if isDead {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
	}
}
