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

// handleTargettingReq implementa PROTO_NC_BAT_TARGETTING_REQ (0x2401).
func (h *Handler) handleTargettingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	targetHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	player.TargetHandle = targetHandle

	// Buscar entidad objetivo para responder con PROTO_NC_BAT_TARGETINFO_CMD
	entity := player.CurrentMap.GetEntity(targetHandle)
	if entity == nil {
		return
	}

	info := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETINFO_CMD)
	info.WriteUint8(1) // Order
	info.WriteUint16(targetHandle)

	switch e := entity.(type) {
	case *ZonePlayer:
		info.WriteUint32(uint32(e.Char.CurHP))
		info.WriteUint32(uint32(e.MaxHP))
		info.WriteUint32(uint32(e.Char.CurSP))
		info.WriteUint32(uint32(e.MaxSP))
		info.WriteUint32(0) // LP
		info.WriteUint32(0) // MaxLP
		info.WriteUint8(uint8(e.Char.Level))
	case *ZoneMonster:
		info.WriteUint32(uint32(e.CurHP))
		info.WriteUint32(uint32(e.MaxHP))
		info.WriteUint32(uint32(e.CurSP))
		info.WriteUint32(uint32(e.MaxSP))
		info.WriteUint32(0)
		info.WriteUint32(0)
		info.WriteUint8(e.Level)
	default:
		info.WriteUint32(100)
		info.WriteUint32(100)
		info.WriteUint32(50)
		info.WriteUint32(50)
		info.WriteUint32(0)
		info.WriteUint32(0)
		info.WriteUint8(1)
	}

	info.WriteUint16(0) // hpchangeorder
	_ = client.SendPacket(info)
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
func (h *Handler) handleAssistReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	tankerHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	tanker := player.CurrentMap.GetPlayer(tankerHandle)
	if tanker == nil || tanker.TargetHandle == 0 {
		return
	}

	player.TargetHandle = tanker.TargetHandle

	// Responder con TARGETINFO del objetivo asistido
	entity := player.CurrentMap.GetEntity(player.TargetHandle)
	if entity == nil {
		return
	}

	info := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETINFO_CMD)
	info.WriteUint8(1)
	info.WriteUint16(player.TargetHandle)
	if mob, ok := entity.(*ZoneMonster); ok {
		info.WriteUint32(uint32(mob.CurHP))
		info.WriteUint32(uint32(mob.MaxHP))
		info.WriteUint32(uint32(mob.CurSP))
		info.WriteUint32(uint32(mob.MaxSP))
		info.WriteUint32(0)
		info.WriteUint32(0)
		info.WriteUint8(mob.Level)
	} else if p, ok := entity.(*ZonePlayer); ok {
		info.WriteUint32(uint32(p.Char.CurHP))
		info.WriteUint32(uint32(p.MaxHP))
		info.WriteUint32(uint32(p.Char.CurSP))
		info.WriteUint32(uint32(p.MaxSP))
		info.WriteUint32(0)
		info.WriteUint32(0)
		info.WriteUint8(uint8(p.Char.Level))
	}
	info.WriteUint16(0)
	_ = client.SendPacket(info)
}

// handleSkillBashObjCastReq implementa PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x2440).
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

	// Costo de SP base si tiene maná suficiente
	spCost := 10
	if player.Char.CurSP >= spCost {
		player.Char.CurSP -= spCost
		spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
		spCmd.WriteUint32(uint32(player.Char.CurSP))
		_ = client.SendPacket(spCmd)
	}

	// Broadcast inicio del golpe de habilidad
	startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD)
	startCmd.WriteUint16(skillID)
	startCmd.WriteUint16(targetObj)
	startCmd.WriteUint16(1) // Index
	player.CurrentMap.Broadcast(startCmd, 0)

	// Calcular daño de habilidad
	dmg := 35
	if mob, ok := targetEntity.(*ZoneMonster); ok {
		mob.CurHP -= dmg
		if mob.CurHP < 0 {
			mob.CurHP = 0
		}
		if mob.CurHP == 0 {
			mob.IsDead = true
		}
	} else if defP, ok := targetEntity.(*ZonePlayer); ok {
		defP.Char.CurHP -= dmg
		if defP.Char.CurHP < 0 {
			defP.Char.CurHP = 0
		}
		hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpCmd.WriteUint32(uint32(defP.Char.CurHP))
		hpCmd.WriteUint16(0)
		_ = defP.Client.SendPacket(hpCmd)
	}
}

// handleSkillBashFldCastReq implementa PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x2441).
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

	// Broadcast inicio de casteo en área
	fldCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD)
	fldCmd.WriteUint16(skillID)
	fldCmd.WriteInt32(fldX)
	fldCmd.WriteInt32(fldY)
	fldCmd.WriteUint16(1)
	player.CurrentMap.Broadcast(fldCmd, 0)

	// Aplicar daño a entidades dentro del radio AoE usando areabmp o distancia euclidiana
	radius := 200.0
	for _, mob := range player.CurrentMap.GetAllMonsters() {
		mx, my := mob.GetPosition()
		dist := math.Hypot(float64(mx)-float64(fldX), float64(my)-float64(fldY))
		if dist <= radius && !mob.IsDead {
			mob.CurHP -= 25
			if mob.CurHP <= 0 {
				mob.CurHP = 0
				mob.IsDead = true
			}
		}
	}
}

// handleSkillBashCastAbortReq implementa PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ (0x2444).
func (h *Handler) handleSkillBashCastAbortReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.CastingSkillID = 0
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

	// 1. Notificar inicio de animación de swing
	swingStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_START_CMD)
	swingStart.WriteUint16(player.MapObjectID)
	swingStart.WriteUint16(player.TargetHandle)
	swingStart.WriteUint8(1) // Actioncode
	swingStart.WriteUint16(1000) // Attackspeed
	swingStart.WriteUint8(1)
	swingStart.WriteUint8(1)
	player.CurrentMap.Broadcast(swingStart, 0)

	// 2. Verificar impacto (Hit / Miss)
	hit := CheckHit(int(player.EquipTH)+30, 10)
	if !hit {
		// Enviar miss
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
			e.State = 3 // Dead
			isDead = true
		} else {
			e.Char.CurHP -= dmg
		}
		restHP = uint32(e.Char.CurHP)
		hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpCmd.WriteUint32(restHP)
		hpCmd.WriteUint16(0)
		_ = e.Client.SendPacket(hpCmd)
	}

	// 4. Notificar impacto con daño
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
	swingDmg.WriteUint16(0) // hpchangeorder
	swingDmg.WriteUint8(1)
	swingDmg.WriteUint8(1)
	player.CurrentMap.Broadcast(swingDmg, 0)

	if isDead {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
		deadCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONEDEAD_LARGE_CMD)
		deadCmd.WriteUint16(targetEntity.GetHandle())
		player.CurrentMap.Broadcast(deadCmd, 0)
	}
}
