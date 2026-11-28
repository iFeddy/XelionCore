package zone

import (
	"math"
	"math/rand"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// CalculateAngleDeg calcula la diferencia angular en grados (0..180) entre el encaramiento del defensor y la posición del atacante.
// En DamageByAngle.json / Zone.exe roe_Damage: 0 grados = frontal (1000), 90 grados = flanco (1100), 180 grados = espalda (1200).
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

// CheckHit evalúa precisión (TH) frente a evasión (TB) usando la fórmula oficial de Zone.exe roe_HitRate:
// permille = int(TH * 850.0 / TB), comparado contra well512_GetRandom(1000).
func CheckHit(aimRate, evasionRate int) bool {
	return CheckHitPermille(ROECalcHitRatePermille(aimRate, evasionRate, 0))
}

// ROECalcHitRatePermille replica RulesOfEngagementNormalPY::roe_HitRate (0x005011F0, const 850.0 / 1000.0).
func ROECalcHitRatePermille(th, tb, freeStatHitRate int) int {
	if th <= 0 {
		th = 10
	}
	if tb <= 0 {
		tb = 1
	}
	permille := int((float64(th)*850.0)/float64(tb)) - freeStatHitRate
	if permille < 50 {
		permille = 50
	}
	return permille
}

// CheckHitPermille tira contra 1000 per-mille (Zone.exe well512_GetRandom(1000)).
func CheckHitPermille(permille int) bool {
	if permille >= 1000 {
		return true
	}
	if permille <= 0 {
		return false
	}
	return rand.Intn(1000) < permille
}

// ROECalcCritRatePermille replica RulesOfEngagementNormalPY::roe_CriticalRate (0x00501700, mínimo 1.0 per-mille).
func ROECalcCritRatePermille(criSum, targetCritResist, freeStatCritRate int) int {
	rate := criSum - targetCritResist + freeStatCritRate
	if rate < 1 {
		rate = 1
	}
	return rate
}

// CheckCritPermille evalúa un crítico en escala oficial per-mille (0..1000).
func CheckCritPermille(permille int) bool {
	if permille < 1 {
		permille = 1
	}
	if permille >= 1000 {
		return true
	}
	return rand.Intn(1000) < permille
}

// ROECalcShieldBlockPermille replica RulesOfEngagementNormalPY::roe_ShieldBlock (0x004FF860):
// ((DefSB + EquipSB) * sbRate) / 1000 en escala per-mille.
func ROECalcShieldBlockPermille(defSB, equipSB, sbRate int) int {
	if sbRate <= 0 {
		sbRate = 1000
	}
	block := ((defSB + equipSB) * sbRate) / 1000
	if block < 0 {
		return 0
	}
	return block
}

// CheckShieldBlock evalúa bloqueo con escudo cuando el defensor lleva escudo equipado (Zone.exe 0x004FF860).
func CheckShieldBlock(defSB, equipSB, sbRate int, shieldEquipped bool) bool {
	if !shieldEquipped {
		return false
	}
	permille := ROECalcShieldBlockPermille(defSB, equipSB, sbRate)
	if permille <= 0 {
		return false
	}
	if permille >= 1000 {
		return true
	}
	return rand.Intn(1000) < permille
}

// CheckCrit evalúa probabilidad de golpe crítico (porcentaje 0..100 escalado a per-mille oficial).
func CheckCrit(critRate, critResist int) bool {
	if critRate <= 0 {
		return false
	}
	permille := (critRate - critResist) * 10
	if permille < 50 {
		permille = 50
	}
	if permille > 800 {
		permille = 800
	}
	return rand.Intn(1000) < permille
}

// writeTargetInfo construye PROTO_NC_BAT_TARGETINFO_CMD (0x2402, 30B payload = 32B wire).
// Verificado en Zone.exe to_SendInform (0x0042B3F0): order=0 para objetivo directo, order=1 para objetivo del objetivo.
func buildTargetInfoPacket(order uint8, targetHandle uint16, entity GridEntity) *protocol.PacketWriter {
	info := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETINFO_CMD)
	info.WriteUint8(order)
	info.WriteUint16(targetHandle)

	var hpOrder uint16

	switch e := entity.(type) {
	case *ZonePlayer:
		var hp, sp uint32
		var lv uint8 = 1
		if e.Char != nil {
			hp = uint32(e.Char.CurHP)
			sp = uint32(e.Char.CurSP)
			lv = uint8(e.Char.Level)
		}
		hpOrder = e.HPChangeOrder
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
		hpOrder = e.HPChangeOrder
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

	info.WriteUint16(hpOrder) // hpchangeorder
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
	// sp_NC_BAT_TARGETTING_REQ: jugador muerto o con estado bloqueante (0x188=392) no puede targetear
	if player.State == 3 || (player.Char != nil && player.Char.CurHP <= 0) || player.GetAbState(392) != nil {
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
	switch e := entity.(type) {
	case *ZonePlayer:
		player.SecondaryTargetHandle = e.TargetHandle
	case *ZoneMonster:
		e.mu.RLock()
		player.SecondaryTargetHandle = e.TargetHandle
		e.mu.RUnlock()
	default:
		player.SecondaryTargetHandle = 0
	}
	sendTargetInfoWithSecondary(client, player.CurrentMap, targetHandle)
}

// handleUntargetReq implementa PROTO_NC_BAT_UNTARGET_REQ (0x2408).
// Verificado en Zone.exe sp_NC_BAT_UNTARGET_REQ (0x0042B6A0).
func (h *Handler) handleUntargetReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.TargetHandle = 0
	player.SecondaryTargetHandle = 0
	player.IsAutoAttacking = false
}

// handleBashStartCmd implementa PROTO_NC_BAT_BASHSTART_CMD (0x242B).
func (h *Handler) handleBashStartCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.State == 3 || (player.Char != nil && player.Char.CurHP <= 0) || player.TargetHandle == 0 || !player.CanAttack() {
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
	Handle    uint16
	Damage    uint32
	RestHP    uint32
	IsCrit    bool
	IsMissed  bool
	IsHeal    bool
	IsAbState bool
	IsResist  bool
	IsDead    bool
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
		if t.IsHeal {
			flag &^= 0x0001
			flag |= 0x0010 // isheal (Zone.exe SkillResultBuffer +0x26)
		}
		if t.IsAbState {
			flag |= 0x0020 // isabstate
		}
		if t.IsResist {
			flag |= 0x0040 // isresist
		}
		if t.IsDead {
			flag |= 0x0100 // isDead
		}
		var hpOrder uint16
		if m := zMap.GetMonster(t.Handle); m != nil {
			hpOrder = m.NextHPChangeOrder()
		} else if p := zMap.GetPlayer(t.Handle); p != nil {
			hpOrder = p.NextHPChangeOrder()
		}
		pkt.WriteUint16(flag)
		pkt.WriteUint32(t.Damage)
		pkt.WriteUint32(t.RestHP)
		pkt.WriteUint16(hpOrder) // hpchangeorder
	}
	zMap.broadcastSkillAOI(pkt, caster, 0)
}

// broadcastSkillAOI difunde un paquete de skill por AOI alrededor de la entidad origen (Zone.exe so_AllOfRange).
// Si el mapa no tiene SpatialGrid inicializado, hace fallback a Broadcast de mapa.
func (m *ZoneMap) broadcastSkillAOI(pw *protocol.PacketWriter, originHandle, excludeHandle uint16) {
	if m == nil || pw == nil {
		return
	}
	if m.grid != nil {
		if ent := m.GetEntity(originHandle); ent != nil {
			ox, oy := ent.GetPosition()
			m.grid.BroadcastAOI(pw, ox, oy, excludeHandle)
			return
		}
	}
	m.Broadcast(pw, excludeHandle)
}

// calcOfficialSkillMaxRange implementa la fórmula exacta de rango de Zone.exe
// (sp_NC_BAT_SKILLBASH_OBJ_CAST_REQ 0x0058016E..0x005801A8 / FLD 0x0057F8EE..0x0057F92B):
//
//	baseRange = objSize + skillRange
//	maxRange  = min((baseRange * 10) >> 3, baseRange + 60)
func calcOfficialSkillMaxRange(objSize, skillRange int) int {
	base := objSize + skillRange
	if base < 0 {
		base = 0
	}
	scaled := (base * 10) >> 3
	if scaled-base > 0x3C {
		return base + 0x3C
	}
	return scaled
}

// calcOfficialSkillCosts implementa las fórmulas de csl_SPCheck (Zone.exe 0x0055B920):
//   - HPCost = baseHP + (MaxHP * hpRate) / 1000
//   - LPCost = baseLP
//   - SPCost = max(0, baseSP + (MaxSP * spRate) / 1000 - empowerSPReduction)
func calcOfficialSkillCosts(skill *SkillDef, maxHP, maxSP int, empowerSPLevel uint8) (hpCost, spCost, lpCost int) {
	if skill == nil {
		return 0, 10, 0
	}
	hpCost = skill.HP + (maxHP*skill.HPRate)/1000
	lpCost = skill.LP
	spCost = skill.SP + (maxSP*skill.SPRate)/1000 - skill.EmpowerSPReduction(empowerSPLevel)
	if spCost < 0 {
		spCost = 0
	}
	return hpCost, spCost, lpCost
}

// calcOfficialSkillCooldown implementa csl_CoolTimeCheck / smo_ply_SkillCoolTime (Zone.exe 0x0046A420 / 0x0055B6A0).
func calcOfficialSkillCooldown(skill *SkillDef, empowerCoolLevel uint8) time.Duration {
	if skill == nil {
		return 50 * time.Millisecond
	}
	dlyMs := skill.DlyTimeMs - skill.EmpowerCoolTimeReduction(empowerCoolLevel)
	if dlyMs < 0 {
		dlyMs = 0
	}
	return time.Duration(dlyMs) * time.Millisecond
}

// mapCanUseSkill verifica FieldInfo.CanSkill (Zone.exe [field+0x9F] != 0, 0x0057FB0E / 0x0057F515).
func (h *Handler) mapCanUseSkill(zMap *ZoneMap) bool {
	if zMap == nil {
		return false
	}
	if h.dataStore != nil && h.dataStore.FieldStore != nil {
		if f, ok := h.dataStore.FieldStore.GetField(zMap.MapID); ok && f != nil {
			return f.Permissions.CanSkill
		}
	}
	return true
}

// isPlayerStateValidForSkill valida el estado del jugador (Zone.exe (1<<state)&0x0004: solo Normal/InGame).
// En XelionCore State==1 (Normal) y State==2 (Player) representan el estado InGame normal.
func isPlayerStateValidForSkill(state byte) bool {
	return state == 1 || state == 2
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

	var hpOrder uint16
	if defP, ok := defender.(*ZonePlayer); ok && defP != nil {
		hpOrder = defP.NextHPChangeOrder()
	}

	swingDmg := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_DAMAGE_CMD)
	swingDmg.WriteUint16(attackerHandle)
	swingDmg.WriteUint16(defenderHandle)
	swingDmg.WriteUint16(dmgFlag)
	swingDmg.WriteUint16(dmg)
	swingDmg.WriteUint32(restHP)
	swingDmg.WriteUint16(hpOrder) // hpchangeorder
	swingDmg.WriteUint8(1)        // damageindex
	swingDmg.WriteUint8(atkSeq)   // attacksequence

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
	someoneDmg.WriteUint16(0)     // hpchangeorder
	someoneDmg.WriteUint8(1)      // damageindex
	someoneDmg.WriteUint8(atkSeq) // attacksequence

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

// sendSkillBashCastFail envía PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK (0x2434, 2B uint16 errCode)
// alineado con Zone.exe so_SendErrorCode / smo_ply_SkillCastFail (0x00579480).
func sendSkillBashCastFail(client *network.Client, errCode uint16) {
	if client == nil {
		return
	}
	failAck := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CAST_FAIL_ACK)
	failAck.WriteUint16(errCode)
	_ = client.SendPacket(failAck)
}

// emitSkillBashHitBlast difunde PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD (0x2457, 6B payload)
// a todo el AOI incluyendo al lanzador (paridad Zone.exe smo_SkillBash_Blast VA 0x0057CC56..0x0057CC9F).
func emitSkillBashHitBlast(zMap *ZoneMap, castIndex, casterHandle, skillID uint16) {
	if zMap == nil {
		return
	}
	blastCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD)
	blastCmd.WriteUint16(castIndex)
	blastCmd.WriteUint16(casterHandle)
	blastCmd.WriteUint16(skillID)
	zMap.broadcastSkillAOI(blastCmd, casterHandle, 0)
}

// sendSkillCoolTime envía PROTO_NC_SKILL_COOLTIME_CMD (0x4821, 6B payload: uint16 skill + uint32 sec)
// sincronizando el temporizador de recarga en el cliente (Zone.exe so_ply_NewSkillCoolTime VA 0x00569260).
func sendSkillCoolTime(client *network.Client, skillID uint16, cd time.Duration) {
	if client == nil || cd < time.Second {
		return
	}
	coolSec := uint32(cd.Milliseconds() / 1000)
	if coolSec == 0 {
		return
	}
	cdCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_COOLTIME_CMD)
	cdCmd.WriteUint16(skillID)
	cdCmd.WriteUint32(coolSec)
	_ = client.SendPacket(cdCmd)
}

// handleSkillBashObjCastReq implementa PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x2440).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x0057FA50) + smo_SkillBash_Casting (0x0057BEF0).
func (h *Handler) handleSkillBashObjCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || player.Char == nil {
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

	// 1. Validación de estado del jugador (Zone.exe 0x0057FAEA: (1<<state)&0x0004 -> 0x0FC0)
	if !isPlayerStateValidForSkill(player.State) || player.Char.CurHP <= 0 {
		sendSkillBashCastFail(client, 0x0FC0)
		return
	}
	// 2. Validación de mapa canUseSkill (Zone.exe 0x0057FB0E -> 0x0FC3)
	if !h.mapCanUseSkill(player.CurrentMap) {
		sendSkillBashCastFail(client, 0x0FC3)
		return
	}
	// 3. Ya casteando o aturdido (Zone.exe 0x0057FB5C: [esi+0x1E78]&0x04 -> 0x0FCC)
	if player.CastingSkillID != 0 || player.IsStunned() {
		sendSkillBashCastFail(client, 0x0FCC)
		return
	}
	// 4. AbStates 0x123 (291) y 0x1F3 (499) en caster (Zone.exe 0x0057FC6A..0x0057FCBF -> 0x0FC1)
	if player.HasAbState(0x123) || player.HasAbState(0x1F3) {
		sendSkillBashCastFail(client, 0x0FC1)
		return
	}
	// 5. Silence (Zone.exe 0x0058001E: smo_SkillBlastOption() == 2 -> 0x0FC7)
	if player.IsSilenced() {
		sendSkillBashCastFail(client, 0x0FC7)
		return
	}

	// 6. Validar objetivo en mapa y si sigue vivo
	targetEntity := player.CurrentMap.GetEntity(targetObj)
	if targetEntity == nil {
		sendSkillBashCastFail(client, 0x0FC4)
		return
	}
	targetObjSize := 0
	switch te := targetEntity.(type) {
	case *ZoneMonster:
		te.mu.RLock()
		dead := te.IsDead || te.CurHP <= 0
		te.mu.RUnlock()
		if dead {
			sendSkillBashCastFail(client, 0x0FC3)
			return
		}
		targetObjSize = 30
	case *ZonePlayer:
		isRevive := h.dataStore != nil && h.dataStore.SkillCatalog != nil && h.dataStore.SkillCatalog.Get(skillID).IsReviveSkill()
		isDeadPlayer := te.State == 3 || (te.Char != nil && te.Char.CurHP <= 0)
		if isRevive {
			if !isDeadPlayer {
				sendSkillBashCastFail(client, 0x0FC3)
				return
			}
		} else if isDeadPlayer {
			sendSkillBashCastFail(client, 0x0FC3)
			return
		}
		// Zone.exe 0x0057FCC1: target con AbState 0x1F3 -> abort silencioso; 0x123 -> 0x0FC2
		if te.HasAbState(0x1F3) {
			return
		}
		if te.HasAbState(0x123) {
			sendSkillBashCastFail(client, 0x0FC2)
			return
		}
	}

	// 7. Validar línea de visión (LoS) con .shbd (Zone.exe 0x0057FD52: so_FindBlockLocate -> 0x0FD1)
	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(tx), float64(ty)) {
		sendSkillBashCastFail(client, 0x0FD1)
		return
	}

	// 8. Lookup en SkillCatalog y validación de existencia / no-pasiva (Zone.exe 0x0057FD92..0x0057FDF6 -> 0x0FC6)
	var skillDef *SkillDef
	if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		if h.dataStore.IsPassiveSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		skillDef = h.dataStore.SkillCatalog.Get(skillID)
		if skillDef == nil {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		// Verificar que el jugador haya aprendido la skill (csl_CoolTimeCheck 0x0046A420 -> 0x0FC8)
		if player.LearnedSkills != nil && !player.HasSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	// 9. Validar cooldown activo individual y de grupo DlyGroupNum (Zone.exe 0x0046A420 -> 0x0FC8)
	now := time.Now()
	if player.SkillCooldowns != nil {
		if cdUntil, onCD := player.SkillCooldowns[skillID]; onCD && now.Before(cdUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}
	if skillDef != nil && skillDef.DlyGroupNum > 0 && player.SkillGroupCooldowns != nil {
		if grpUntil, onGrpCD := player.SkillGroupCooldowns[skillDef.DlyGroupNum]; onGrpCD && now.Before(grpUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	// 10. Validar costos HP / LP / SP (Zone.exe csl_SPCheck 0x0055B920 + 0x005800D5..0x00580169)
	_, empSP, _, empCool := player.GetSkillEmpower(skillID)
	hpCost, spCost, lpCost := calcOfficialSkillCosts(skillDef, player.MaxHP, player.MaxSP, empSP)
	if hpCost > 0 && player.Char.CurHP <= hpCost {
		sendSkillBashCastFail(client, 0x0FD5)
		return
	}
	if lpCost > 0 {
		// CurLP no implementado en tCharacter aún (0): si skill exige LP > 0 -> 0x0FDC
		sendSkillBashCastFail(client, 0x0FDC)
		return
	}
	if player.Char.CurSP < spCost {
		sendSkillBashCastFail(client, 0x0FC9)
		return
	}

	// 11. Validar distancia máxima oficial (Zone.exe 0x0058016E..0x005801A8 -> 0x0FCA)
	maxRange := 450
	if skillDef != nil && skillDef.Range > 0 {
		maxRange = calcOfficialSkillMaxRange(targetObjSize, skillDef.Range)
	}
	dx := int64(px - tx)
	dy := int64(py - ty)
	if dx*dx+dy*dy > int64(maxRange)*int64(maxRange) {
		sendSkillBashCastFail(client, 0x0FCA)
		return
	}

	player.CastingSkillID = skillID
	player.CastingTargetHandle = targetObj
	empDmg, _, empKeep, _ := player.GetSkillEmpower(skillID)

	castMs := 0
	if skillDef != nil && skillDef.CastTimeMs > 0 {
		castMs = skillDef.CastTimeMs
		if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
			if pb := h.dataStore.SkillCatalog.ComputePassiveBonuses(player); pb.CastingTimeRate > 0 && pb.CastingTimeRate < 1000 {
				castMs = (castMs * pb.CastingTimeRate) / 1000
			}
		}
	}

	// Confirmar inicio de casteo exitoso: PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK (0x2435, vacío)
	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK))

	if castMs > 0 {
		// Difundir inicio de barra de casteo a AOI: PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD (0x2442, 6B)
		someoneCast := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_OBJ_CAST_CMD)
		someoneCast.WriteUint16(skillID)
		someoneCast.WriteUint16(player.MapObjectID)
		someoneCast.WriteUint16(targetObj)
		player.CurrentMap.broadcastSkillAOI(someoneCast, player.MapObjectID, player.MapObjectID)

		player.PendingCast = &PendingSkillCast{
			SkillID:   skillID,
			IsField:   false,
			TargetObj: targetObj,
			ExecuteAt: now.Add(time.Duration(castMs) * time.Millisecond),
			IsMoving:  skillDef.IsMovingSkill,
			HPCost:    hpCost,
			SPCost:    spCost,
			EmpDmg:    empDmg,
			EmpKeep:   empKeep,
			EmpCool:   empCool,
			SkillDef:  skillDef,
		}
		return
	}

	h.finishObjSkillCast(player, skillID, targetObj, skillDef, hpCost, spCost, empDmg, empKeep, empCool, now)
}

func (h *Handler) finishObjSkillCast(player *ZonePlayer, skillID, targetObj uint16, skillDef *SkillDef, hpCost, spCost int, empDmg, empKeep, empCool uint8, now time.Time) {
	if player == nil || player.CurrentMap == nil || player.Char == nil {
		return
	}
	defer func() {
		player.PendingCast = nil
		player.CastingSkillID = 0
		player.CastingTargetHandle = 0
	}()

	targetEntity := player.CurrentMap.GetEntity(targetObj)
	if targetEntity == nil {
		return
	}
	if mob, isMob := targetEntity.(*ZoneMonster); isMob && mob.IsDead {
		return
	}
	if tp, isPlayer := targetEntity.(*ZonePlayer); isPlayer && tp.State == 3 {
		if skillDef == nil || !skillDef.IsReviveSkill() {
			return
		}
	}

	if player.SkillCooldowns == nil {
		player.SkillCooldowns = make(map[uint16]time.Time)
	}
	skillCD := calcOfficialSkillCooldown(skillDef, empCool)
	player.SkillCooldowns[skillID] = now.Add(skillCD)
	if skillDef != nil && skillDef.DlyGroupNum > 0 && skillDef.DlyTimeGroupMs > 0 {
		if player.SkillGroupCooldowns == nil {
			player.SkillGroupCooldowns = make(map[int]time.Time)
		}
		player.SkillGroupCooldowns[skillDef.DlyGroupNum] = now.Add(time.Duration(skillDef.DlyTimeGroupMs) * time.Millisecond)
	}

	// Consumir HP y SP según csl_SPCheck y notificar cambios
	if hpCost > 0 {
		player.Char.CurHP -= hpCost
		if player.Char.CurHP < 1 {
			player.Char.CurHP = 1
		}
		if player.Client != nil {
			hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpCmd.WriteUint32(uint32(player.Char.CurHP))
			hpCmd.WriteUint16(player.NextHPChangeOrder())
			_ = player.Client.SendPacket(hpCmd)
		}
	}
	if spCost > 0 {
		player.Char.CurSP -= spCost
		if player.Char.CurSP < 0 {
			player.Char.CurSP = 0
		}
		if player.Client != nil {
			spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
			spCmd.WriteUint32(uint32(player.Char.CurSP))
			_ = player.Client.SendPacket(spCmd)
		}
	}

	castIndex := player.NextSkillCastIndex()

	// Lanzador recibe PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD (0x244E, 6B payload = 8B wire)
	if player.Client != nil {
		startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD)
		startCmd.WriteUint16(skillID)
		startCmd.WriteUint16(targetObj)
		startCmd.WriteUint16(castIndex)
		_ = player.Client.SendPacket(startCmd)
	}

	// Circundantes en AOI reciben PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD (0x244F, 8B payload = 10B wire)
	someoneStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD)
	someoneStart.WriteUint16(player.MapObjectID)
	someoneStart.WriteUint16(skillID)
	someoneStart.WriteUint16(targetObj)
	someoneStart.WriteUint16(castIndex)
	player.CurrentMap.broadcastSkillAOI(someoneStart, player.MapObjectID, player.MapObjectID)

	// Resolver impacto: si es Revive (SkillTargetState == 1) sobre jugador muerto, ofrecer resurrección (0x4829)
	if skillDef != nil && skillDef.IsReviveSkill() {
		if tp, ok := targetEntity.(*ZonePlayer); ok && tp != nil {
			tp.PendingRevive = true
			tp.PendingReviveRate = 1000 // 100% HP/SP (sp_ReviveNow 0x3E8)
			if tp.Client != nil {
				_ = tp.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SOMEONEREVAVALTOME_CMD))
			}
		}
		emitSkillBashHitBlast(player.CurrentMap, castIndex, player.MapObjectID, skillID)
		sendSkillCoolTime(player.Client, skillID, skillCD)
		return
	}

	// Resolver impacto usando motor oficial smo_SkillBlast (daño físico/mágico, heal, AoE, AbStates)
	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	var hitTargets []skillDamageTarget
	var deadEntities []GridEntity

	if skillDef != nil && skillDef.Area > 0 && skillDef.TargetNumber > 1 && skillDef.HealAmount() == 0 && !skillDef.IsPureBuff() {
		// Skill objetivo con AoE alrededor del target (sbe_BlastObject 0x00583BC0)
		mobs := selectAoEMonsters(player.CurrentMap, px, py, tx, ty, skillDef)
		for _, m := range mobs {
			res, dead := h.resolveSkillBlastOnEntity(player, m, skillDef, empDmg, empKeep)
			hitTargets = append(hitTargets, res)
			if dead {
				deadEntities = append(deadEntities, m)
			}
		}
	} else {
		res, dead := h.resolveSkillBlastOnEntity(player, targetEntity, skillDef, empDmg, empKeep)
		hitTargets = append(hitTargets, res)
		if dead {
			deadEntities = append(deadEntities, targetEntity)
		}
	}

	if skillDef == nil || !skillDef.IsPureBuff() {
		emitSkillBashHitDamage(player.CurrentMap, castIndex, player.MapObjectID, hitTargets)
	}
	emitSkillBashHitBlast(player.CurrentMap, castIndex, player.MapObjectID, skillID)
	sendSkillCoolTime(player.Client, skillID, skillCD)

	for _, de := range deadEntities {
		emitReallyKill(player.CurrentMap, de.GetHandle(), player.MapObjectID)
		if defP, ok := de.(*ZonePlayer); ok {
			sendCharDeadMenu(defP) // el menú va a la víctima, no al atacante
		}
	}

	if durMs := skillDef.MagicFieldDurationMs(); durMs > 0 && !skillDef.IsPureBuff() {
		player.CurrentMap.AddMagicField(&ActiveMagicField{
			Caster:       player,
			SkillDef:     skillDef,
			X:            int32(tx),
			Y:            int32(ty),
			EmpDmg:       empDmg,
			EmpKeep:      empKeep,
			CastIndex:    castIndex,
			ExpiresAt:    now.Add(time.Duration(durMs) * time.Millisecond),
			LastTickAt:   now,
			TickInterval: time.Second,
		})
	}
}

// handleSkillBashFldCastReq implementa PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x2441).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x0057F400).
func (h *Handler) handleSkillBashFldCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || player.Char == nil {
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

	// 1. Estado del jugador (Zone.exe 0x0057F4A2 -> 0x0FC0)
	if !isPlayerStateValidForSkill(player.State) || player.Char.CurHP <= 0 {
		sendSkillBashCastFail(client, 0x0FC0)
		return
	}
	// 2. AbStates 0x123 / 0x1F3 en caster (Zone.exe 0x0057F4D8..0x0057F504 -> 0x0FC1)
	if player.HasAbState(0x123) || player.HasAbState(0x1F3) {
		sendSkillBashCastFail(client, 0x0FC1)
		return
	}
	// 3. Mapa canUseSkill (Zone.exe 0x0057F515 -> 0x0FC3)
	if !h.mapCanUseSkill(player.CurrentMap) {
		sendSkillBashCastFail(client, 0x0FC3)
		return
	}
	// 4. Ya casteando o aturdido (Zone.exe 0x0057F594 -> 0x0FCC)
	if player.CastingSkillID != 0 || player.IsStunned() {
		sendSkillBashCastFail(client, 0x0FCC)
		return
	}
	// 5. Silence (Zone.exe 0x0057F7BF -> 0x0FC7)
	if player.IsSilenced() {
		sendSkillBashCastFail(client, 0x0FC7)
		return
	}

	// 6. LoS en FLD (Zone.exe 0x0057F605..0x0057F638 -> 0x0FD1)
	px, py := player.GetPosition()
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(fldX), float64(fldY)) {
		sendSkillBashCastFail(client, 0x0FD1)
		return
	}

	// 7. Lookup en SkillCatalog y validación de existencia / aprendizaje (Zone.exe 0x0057F648..0x0057F69A -> 0x0FC6 / 0x0FC8)
	var skillDef *SkillDef
	if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		if h.dataStore.IsPassiveSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		skillDef = h.dataStore.SkillCatalog.Get(skillID)
		if skillDef == nil {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		if player.LearnedSkills != nil && !player.HasSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	// 8. Si es skill de teleport/dash, validar move-block (Zone.exe 0x0057F6F8..0x0057F74C -> 0x0FD1)
	if skillDef != nil && skillDef.IsTeleportSkill() && player.CurrentMap.blockInfo != nil {
		if player.CurrentMap.blockInfo.IsWorldMoveBlocked(float64(fldX), float64(fldY)) {
			sendSkillBashCastFail(client, 0x0FD1)
			return
		}
	}

	// 9. Cooldown activo (Zone.exe 0x0046A420 -> 0x0FC8)
	now := time.Now()
	if player.SkillCooldowns != nil {
		if cdUntil, onCD := player.SkillCooldowns[skillID]; onCD && now.Before(cdUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}
	if skillDef != nil && skillDef.DlyGroupNum > 0 && player.SkillGroupCooldowns != nil {
		if grpUntil, onGrpCD := player.SkillGroupCooldowns[skillDef.DlyGroupNum]; onGrpCD && now.Before(grpUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	// 10. Costos HP / LP / SP (Zone.exe 0x0057F869..0x0057F8E9)
	_, empSP, _, empCool := player.GetSkillEmpower(skillID)
	hpCost, spCost, lpCost := calcOfficialSkillCosts(skillDef, player.MaxHP, player.MaxSP, empSP)
	if skillDef == nil {
		spCost = 15
	}
	if hpCost > 0 && player.Char.CurHP <= hpCost {
		sendSkillBashCastFail(client, 0x0FD5)
		return
	}
	if lpCost > 0 {
		sendSkillBashCastFail(client, 0x0FDC)
		return
	}
	if player.Char.CurSP < spCost {
		sendSkillBashCastFail(client, 0x0FC9)
		return
	}

	// 11. Rango máximo oficial FLD (Zone.exe 0x0057F8EE..0x0057F92B -> 0x0FCA)
	maxRange := 600
	if skillDef != nil && skillDef.Range > 0 {
		maxRange = calcOfficialSkillMaxRange(0, skillDef.Range)
	}
	dx := int64(px) - int64(fldX)
	dy := int64(py) - int64(fldY)
	if dx*dx+dy*dy > int64(maxRange)*int64(maxRange) {
		sendSkillBashCastFail(client, 0x0FCA)
		return
	}

	player.CastingSkillID = skillID
	empDmg, _, empKeep, _ := player.GetSkillEmpower(skillID)

	castMs := 0
	if skillDef != nil && skillDef.CastTimeMs > 0 {
		castMs = skillDef.CastTimeMs
		if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
			if pb := h.dataStore.SkillCatalog.ComputePassiveBonuses(player); pb.CastingTimeRate > 0 && pb.CastingTimeRate < 1000 {
				castMs = (castMs * pb.CastingTimeRate) / 1000
			}
		}
	}

	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK))

	if castMs > 0 {
		// Difundir inicio de barra de casteo a AOI: PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD (0x2443, 12B)
		someoneCast := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_FLD_CAST_CMD)
		someoneCast.WriteUint16(skillID)
		someoneCast.WriteUint16(player.MapObjectID)
		someoneCast.WriteInt32(fldX)
		someoneCast.WriteInt32(fldY)
		player.CurrentMap.broadcastSkillAOI(someoneCast, player.MapObjectID, player.MapObjectID)

		player.PendingCast = &PendingSkillCast{
			SkillID:   skillID,
			IsField:   true,
			FldX:      fldX,
			FldY:      fldY,
			ExecuteAt: now.Add(time.Duration(castMs) * time.Millisecond),
			IsMoving:  skillDef.IsMovingSkill,
			HPCost:    hpCost,
			SPCost:    spCost,
			EmpDmg:    empDmg,
			EmpKeep:   empKeep,
			EmpCool:   empCool,
			SkillDef:  skillDef,
		}
		return
	}

	h.finishFldSkillCast(player, skillID, fldX, fldY, skillDef, hpCost, spCost, empDmg, empKeep, empCool, now)
}

func (h *Handler) finishFldSkillCast(player *ZonePlayer, skillID uint16, fldX, fldY int32, skillDef *SkillDef, hpCost, spCost int, empDmg, empKeep, empCool uint8, now time.Time) {
	if player == nil || player.CurrentMap == nil || player.Char == nil {
		return
	}
	defer func() {
		player.PendingCast = nil
		player.CastingSkillID = 0
		player.CastingTargetHandle = 0
	}()

	if player.SkillCooldowns == nil {
		player.SkillCooldowns = make(map[uint16]time.Time)
	}
	skillCD := calcOfficialSkillCooldown(skillDef, empCool)
	player.SkillCooldowns[skillID] = now.Add(skillCD)
	if skillDef != nil && skillDef.DlyGroupNum > 0 && skillDef.DlyTimeGroupMs > 0 {
		if player.SkillGroupCooldowns == nil {
			player.SkillGroupCooldowns = make(map[int]time.Time)
		}
		player.SkillGroupCooldowns[skillDef.DlyGroupNum] = now.Add(time.Duration(skillDef.DlyTimeGroupMs) * time.Millisecond)
	}

	if hpCost > 0 {
		player.Char.CurHP -= hpCost
		if player.Char.CurHP < 1 {
			player.Char.CurHP = 1
		}
		if player.Client != nil {
			hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpCmd.WriteUint32(uint32(player.Char.CurHP))
			hpCmd.WriteUint16(player.NextHPChangeOrder())
			_ = player.Client.SendPacket(hpCmd)
		}
	}
	if spCost > 0 {
		player.Char.CurSP -= spCost
		if player.Char.CurSP < 0 {
			player.Char.CurSP = 0
		}
		if player.Client != nil {
			spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
			spCmd.WriteUint32(uint32(player.Char.CurSP))
			_ = player.Client.SendPacket(spCmd)
		}
	}

	castIndex := player.NextSkillCastIndex()

	// Lanzador recibe PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD (0x2450, 12B payload)
	if player.Client != nil {
		fldCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD)
		fldCmd.WriteUint16(skillID)
		fldCmd.WriteInt32(fldX)
		fldCmd.WriteInt32(fldY)
		fldCmd.WriteUint16(castIndex)
		_ = player.Client.SendPacket(fldCmd)
	}

	// Circundantes en AOI reciben PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD (0x2451, 14B payload)
	someoneFld := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD)
	someoneFld.WriteUint16(player.MapObjectID)
	someoneFld.WriteUint16(skillID)
	someoneFld.WriteInt32(fldX)
	someoneFld.WriteInt32(fldY)
	someoneFld.WriteUint16(castIndex)
	player.CurrentMap.broadcastSkillAOI(someoneFld, player.MapObjectID, player.MapObjectID)

	// Selección oficial AoE (sbe_BlastLocate 0x00584160): Area/2, cono SkillDegree, qsort por distancia, cap TargetNumber
	px, py := player.GetPosition()
	var hitTargets []skillDamageTarget
	var deadHandles []uint16
	if skillDef != nil && skillDef.IsPureBuff() {
		res, _ := h.resolveSkillBlastOnEntity(player, player, skillDef, empDmg, empKeep)
		hitTargets = append(hitTargets, res)
	} else {
		mobs := selectAoEMonsters(player.CurrentMap, px, py, int(fldX), int(fldY), skillDef)
		for _, mob := range mobs {
			res, isDead := h.resolveSkillBlastOnEntity(player, mob, skillDef, empDmg, empKeep)
			hitTargets = append(hitTargets, res)
			if isDead {
				deadHandles = append(deadHandles, mob.Handle)
			}
		}
	}

	if skillDef == nil || !skillDef.IsPureBuff() {
		emitSkillBashHitDamage(player.CurrentMap, castIndex, player.MapObjectID, hitTargets)
	}
	emitSkillBashHitBlast(player.CurrentMap, castIndex, player.MapObjectID, skillID)
	sendSkillCoolTime(player.Client, skillID, skillCD)
	for _, deadH := range deadHandles {
		emitReallyKill(player.CurrentMap, deadH, player.MapObjectID)
	}

	if durMs := skillDef.MagicFieldDurationMs(); durMs > 0 && !skillDef.IsPureBuff() {
		player.CurrentMap.AddMagicField(&ActiveMagicField{
			Caster:       player,
			SkillDef:     skillDef,
			X:            fldX,
			Y:            fldY,
			EmpDmg:       empDmg,
			EmpKeep:      empKeep,
			CastIndex:    castIndex,
			ExpiresAt:    now.Add(time.Duration(durMs) * time.Millisecond),
			LastTickAt:   now,
			TickInterval: time.Second,
		})
	}
}

// FinishPendingSkillCast completa una habilidad con CastTime > 0 cuando llega ExecuteAt en tick100ms.
func (h *Handler) FinishPendingSkillCast(player *ZonePlayer, now time.Time) {
	if player == nil || player.PendingCast == nil {
		return
	}
	pc := player.PendingCast
	if now.Before(pc.ExecuteAt) {
		return
	}
	if player.State == 3 || player.IsStunned() || player.IsSilenced() {
		h.abortCastAndBroadcast(player)
		return
	}
	if pc.IsField {
		h.finishFldSkillCast(player, pc.SkillID, pc.FldX, pc.FldY, pc.SkillDef, pc.HPCost, pc.SPCost, pc.EmpDmg, pc.EmpKeep, pc.EmpCool, now)
	} else {
		h.finishObjSkillCast(player, pc.SkillID, pc.TargetObj, pc.SkillDef, pc.HPCost, pc.SPCost, pc.EmpDmg, pc.EmpKeep, pc.EmpCool, now)
	}
}

// abortCastAndBroadcast cancela un casteo activo y difunde PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD (0x2438).
func (h *Handler) abortCastAndBroadcast(player *ZonePlayer) {
	if player == nil {
		return
	}
	player.PendingCast = nil
	player.CastingSkillID = 0
	player.CastingTargetHandle = 0
	if player.CurrentMap != nil {
		cutCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD)
		cutCmd.WriteUint16(player.MapObjectID)
		player.CurrentMap.broadcastSkillAOI(cutCmd, player.MapObjectID, 0)
	}
}

// abortNonMovingCastOnMove cancela la canalización si el jugador se mueve con una skill !IsMovingSkill.
func (h *Handler) abortNonMovingCastOnMove(player *ZonePlayer) {
	if player == nil || player.PendingCast == nil || player.PendingCast.IsMoving {
		return
	}
	h.abortCastAndBroadcast(player)
}

// handleSkillBashCastAbortReq implementa PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ (0x2444).
// Verificado en Zone.exe sp_NC_BAT_SKILLBASH_CASTABORT_REQ (0x0057C480) + sp_NC_BAT_SKILLBASH_CASTABORT_CMD (0x00579F10).
func (h *Handler) handleSkillBashCastAbortReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	var errCode uint16
	if player.State == 6 {
		// Zone.exe 0x0057C4EE: estado 6 -> 0x20E9
		errCode = 0x20E9
	} else if !isPlayerStateValidForSkill(player.State) {
		// Zone.exe 0x0057C538: bitmask (1<<state)&0x0004 falla -> 0x20EA
		errCode = 0x20EA
	} else {
		// Zone.exe 0x00579F10 (sp_NC_BAT_SKILLBASH_CASTABORT_CMD): emite 0x2438 (4B)
		errCode = 0x20E8
		cutCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_CASTCUT_CMD)
		cutCmd.WriteUint16(player.MapObjectID)
		if player.CastingSkillID != 0 {
			player.PendingCast = nil
			player.CastingSkillID = 0
			player.CastingTargetHandle = 0
			if player.CurrentMap != nil {
				// Broadcast por AOI (so_AllOfRange 0x0057A0AC..0x0057A0DD)
				player.CurrentMap.broadcastSkillAOI(cutCmd, player.MapObjectID, 0)
			}
		} else {
			// Si no estaba casteando, envía 0x2438 solo al propio jugador (0x0057A016..0x0057A09E)
			_ = client.SendPacket(cutCmd)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// isMoverStateValidForSkill valida el estado del jinete/montura (Zone.exe so_smv_SkillBash_*: bitmask 0x0006 -> estados 1, 2 o 6 montado).
func isMoverStateValidForSkill(player *ZonePlayer) bool {
	if player == nil || player.Char == nil || player.Char.CurHP <= 0 {
		return false
	}
	if player.State == 3 {
		return false
	}
	return player.State == 1 || player.State == 2 || player.State == 6
}

// handleMoverSkillBashObjCastReq implementa PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ (0xCC0E).
// Verificado en Zone.exe ShineMover::so_smv_SkillBash_Object (0x00580800):
// bitmask 0x0006, sin consumo de SP del jugador, cooldown en MoverSkillList (msl_CoolTimeCheck),
// rango oficial min(1.25x, +60), LoS (0x0FD1) y error 0x0FD6 para objetivo inválido.
func (h *Handler) handleMoverSkillBashObjCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || player.Char == nil {
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

	if !isMoverStateValidForSkill(player) {
		sendSkillBashCastFail(client, 0x0FC0)
		return
	}
	targetEntity := player.CurrentMap.GetEntity(targetObj)
	if targetEntity == nil {
		sendSkillBashCastFail(client, 0x0FC4)
		return
	}
	if !h.mapCanUseSkill(player.CurrentMap) {
		sendSkillBashCastFail(client, 0x0FC3)
		return
	}
	if player.CastingSkillID != 0 || player.IsStunned() {
		sendSkillBashCastFail(client, 0x0FCC)
		return
	}
	if player.IsSilenced() {
		sendSkillBashCastFail(client, 0x0FC7)
		return
	}
	if player.HasAbState(0x123) || player.HasAbState(0x1F3) {
		sendSkillBashCastFail(client, 0x0FC1)
		return
	}

	targetObjSize := 0
	switch te := targetEntity.(type) {
	case *ZoneMonster:
		te.mu.RLock()
		dead := te.IsDead || te.CurHP <= 0
		killable := te.IsKillableByUser
		te.mu.RUnlock()
		if dead || !killable {
			// Zone.exe 0x00580974 / 0x00580BD3: target no matable -> 0x0FD6
			sendSkillBashCastFail(client, 0x0FD6)
			return
		}
		targetObjSize = 30
	case *ZonePlayer:
		if te.HasAbState(0x1F3) || te.State == 3 || (te.Char != nil && te.Char.CurHP <= 0) {
			sendSkillBashCastFail(client, 0x0FD6)
			return
		}
		if te.HasAbState(0x123) {
			sendSkillBashCastFail(client, 0x0FC2)
			return
		}
	}

	px, py := player.GetPosition()
	tx, ty := targetEntity.GetPosition()
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(tx), float64(ty)) {
		sendSkillBashCastFail(client, 0x0FD1)
		return
	}

	var skillDef *SkillDef
	if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		if h.dataStore.IsPassiveSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		skillDef = h.dataStore.SkillCatalog.Get(skillID)
		if skillDef == nil {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
	}

	// MoverSkillList::msl_CoolTimeCheck (0x00580C52) — sin consumo de SP del jugador
	now := time.Now()
	if player.MoverSkillCooldowns != nil {
		if cdUntil, onCD := player.MoverSkillCooldowns[skillID]; onCD && now.Before(cdUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	maxRange := 450
	if skillDef != nil && skillDef.Range > 0 {
		maxRange = calcOfficialSkillMaxRange(targetObjSize, skillDef.Range)
	}
	dx := int64(px - tx)
	dy := int64(py - ty)
	if dx*dx+dy*dy > int64(maxRange)*int64(maxRange) {
		sendSkillBashCastFail(client, 0x0FCA)
		return
	}

	player.CastingSkillID = skillID
	if player.MoverSkillCooldowns == nil {
		player.MoverSkillCooldowns = make(map[uint16]time.Time)
	}
	player.MoverSkillCooldowns[skillID] = now.Add(calcOfficialSkillCooldown(skillDef, 0))

	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK))
	castIndex := player.NextSkillCastIndex()

	casterHandle := player.MapObjectID
	if player.MoverHandle != 0 {
		casterHandle = player.MoverHandle
	}

	startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD)
	startCmd.WriteUint16(skillID)
	startCmd.WriteUint16(targetObj)
	startCmd.WriteUint16(castIndex)
	_ = client.SendPacket(startCmd)

	someoneStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD)
	someoneStart.WriteUint16(casterHandle)
	someoneStart.WriteUint16(skillID)
	someoneStart.WriteUint16(targetObj)
	someoneStart.WriteUint16(castIndex)
	player.CurrentMap.broadcastSkillAOI(someoneStart, player.MapObjectID, player.MapObjectID)

	res, dead := h.resolveSkillBlastOnEntity(player, targetEntity, skillDef, 0, 0)
	emitSkillBashHitDamage(player.CurrentMap, castIndex, casterHandle, []skillDamageTarget{res})
	emitSkillBashHitBlast(player.CurrentMap, castIndex, casterHandle, skillID)
	if dead {
		emitReallyKill(player.CurrentMap, targetObj, player.MapObjectID)
		if defP, ok := targetEntity.(*ZonePlayer); ok {
			sendCharDeadMenu(defP)
		}
	}
	player.CastingSkillID = 0
}

// handleMoverSkillBashFldCastReq implementa PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ (0xCC0F).
// Verificado en Zone.exe ShineMover::so_smv_SkillBash_Location (0x00580CF0).
func (h *Handler) handleMoverSkillBashFldCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || player.Char == nil {
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

	if !isMoverStateValidForSkill(player) {
		sendSkillBashCastFail(client, 0x0FC0)
		return
	}
	if !h.mapCanUseSkill(player.CurrentMap) {
		sendSkillBashCastFail(client, 0x0FC3)
		return
	}
	if player.CastingSkillID != 0 || player.IsStunned() {
		sendSkillBashCastFail(client, 0x0FCC)
		return
	}
	if player.IsSilenced() {
		sendSkillBashCastFail(client, 0x0FC7)
		return
	}
	if player.HasAbState(0x123) || player.HasAbState(0x1F3) {
		sendSkillBashCastFail(client, 0x0FC1)
		return
	}

	px, py := player.GetPosition()
	if player.CurrentMap.blockInfo != nil && player.CurrentMap.blockInfo.IsWorldAttackBlocked(float64(px), float64(py), float64(fldX), float64(fldY)) {
		sendSkillBashCastFail(client, 0x0FD1)
		return
	}

	var skillDef *SkillDef
	if h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		if h.dataStore.IsPassiveSkill(skillID) {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
		skillDef = h.dataStore.SkillCatalog.Get(skillID)
		if skillDef == nil {
			sendSkillBashCastFail(client, 0x0FC6)
			return
		}
	}

	if skillDef != nil && skillDef.IsTeleportSkill() && player.CurrentMap.blockInfo != nil {
		if player.CurrentMap.blockInfo.IsWorldMoveBlocked(float64(fldX), float64(fldY)) {
			sendSkillBashCastFail(client, 0x0FD1)
			return
		}
	}

	now := time.Now()
	if player.MoverSkillCooldowns != nil {
		if cdUntil, onCD := player.MoverSkillCooldowns[skillID]; onCD && now.Before(cdUntil) {
			sendSkillBashCastFail(client, 0x0FC8)
			return
		}
	}

	maxRange := 600
	if skillDef != nil && skillDef.Range > 0 {
		maxRange = calcOfficialSkillMaxRange(0, skillDef.Range)
	}
	dx := int64(px) - int64(fldX)
	dy := int64(py) - int64(fldY)
	if dx*dx+dy*dy > int64(maxRange)*int64(maxRange) {
		sendSkillBashCastFail(client, 0x0FCA)
		return
	}

	player.CastingSkillID = skillID
	if player.MoverSkillCooldowns == nil {
		player.MoverSkillCooldowns = make(map[uint16]time.Time)
	}
	player.MoverSkillCooldowns[skillID] = now.Add(calcOfficialSkillCooldown(skillDef, 0))

	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CAST_SUC_ACK))
	castIndex := player.NextSkillCastIndex()

	casterHandle := player.MapObjectID
	if player.MoverHandle != 0 {
		casterHandle = player.MoverHandle
	}

	fldCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_HIT_FLD_START_CMD)
	fldCmd.WriteUint16(skillID)
	fldCmd.WriteInt32(fldX)
	fldCmd.WriteInt32(fldY)
	fldCmd.WriteUint16(castIndex)
	_ = client.SendPacket(fldCmd)

	someoneFld := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SOMEONESKILLBASH_HIT_FLD_START_CMD)
	someoneFld.WriteUint16(casterHandle)
	someoneFld.WriteUint16(skillID)
	someoneFld.WriteInt32(fldX)
	someoneFld.WriteInt32(fldY)
	someoneFld.WriteUint16(castIndex)
	player.CurrentMap.broadcastSkillAOI(someoneFld, player.MapObjectID, player.MapObjectID)

	mobs := selectAoEMonsters(player.CurrentMap, px, py, int(fldX), int(fldY), skillDef)
	var hitTargets []skillDamageTarget
	var deadHandles []uint16
	for _, mob := range mobs {
		res, isDead := h.resolveSkillBlastOnEntity(player, mob, skillDef, 0, 0)
		hitTargets = append(hitTargets, res)
		if isDead {
			deadHandles = append(deadHandles, mob.Handle)
		}
	}

	emitSkillBashHitDamage(player.CurrentMap, castIndex, casterHandle, hitTargets)
	emitSkillBashHitBlast(player.CurrentMap, castIndex, casterHandle, skillID)
	for _, deadH := range deadHandles {
		emitReallyKill(player.CurrentMap, deadH, player.MapObjectID)
	}
	player.CastingSkillID = 0
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
			hpCmd.WriteUint16(e.NextHPChangeOrder())
			_ = e.Client.SendPacket(hpCmd)
		}
	}

	// 4. Notificar impacto con daño (0x2448 atacante/defensor + 0x2449 entorno + 0x244A si muere)
	emitSwingDamage(player.CurrentMap, player, targetEntity, player.MapObjectID, player.TargetHandle, uint16(dmg), restHP, isCrit, false, isDead, 1)

	if isDead {
		player.IsAutoAttacking = false
		player.TargetHandle = 0
		if defP, ok := targetEntity.(*ZonePlayer); ok {
			sendCharDeadMenu(defP) // el menú va a la víctima, no al atacante
		}
	}
}
