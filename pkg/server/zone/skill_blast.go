package zone

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// calcOfficialSkillAttackPower implementa RulesOfEngagementPhisycalSkill::roe_AttackPower (0x00506B40)
// y RulesOfEngagementMagicalSkill::roe_AttackPower (0x00506ED0):
//
//	minAtk = minBase + (minBase * rateMin) / 1000.0 + addMin + empDmgBonus
//	maxAtk = maxBase + (maxBase * rateMax) / 1000.0 + addMax + empDmgBonus
//	result = minAtk + rb_largerandom(int(maxAtk - minAtk))
func calcOfficialSkillAttackPower(player *ZonePlayer, skill *SkillDef, empDmg uint8) int {
	if skill == nil {
		return 35
	}
	empBonus := float64(skill.EmpowerDamageBonus(empDmg))
	var minBase, maxBase, minRate, maxRate, minAdd, maxAdd float64

	if skill.IsMagical() {
		if player != nil {
			minBase = float64(clampStatUint32(player.EquipMinMA, player.AbMA)) + float64(player.Int)
			maxBase = float64(clampStatUint32(player.EquipMaxMA, player.AbMA)) + float64(player.Int)
		}
		minRate = float64(skill.MinMARate)
		maxRate = float64(skill.MaxMARate)
		minAdd = float64(skill.MinMA)
		maxAdd = float64(skill.MaxMA)
	} else {
		if player != nil {
			minBase = float64(clampStatUint32(player.EquipMinWC, player.AbWC)) + float64(player.Str)
			maxBase = float64(clampStatUint32(player.EquipMaxWC, player.AbWC)) + float64(player.Str)
		}
		minRate = float64(skill.MinWCRate)
		maxRate = float64(skill.MaxWCRate)
		minAdd = float64(skill.MinWC)
		maxAdd = float64(skill.MaxWC)
	}

	minAtk := minBase + (minBase*minRate)/1000.0 + minAdd + empBonus
	maxAtk := maxBase + (maxBase*maxRate)/1000.0 + maxAdd + empBonus
	if maxAtk < minAtk {
		maxAtk = minAtk
	}

	span := int(maxAtk - minAtk)
	val := int(minAtk)
	if span > 0 {
		val += rand.Intn(span + 1)
	}
	if val < 1 {
		val = 1
	}
	return val
}

// calcOfficialSkillHeal implementa smo_SkillBlast heal (Zone.exe 0x005817CE..0x00581980):
// heal = baseHeal + empowerDamageBonus(nibble 0).
func calcOfficialSkillHeal(skill *SkillDef, empDmg uint8) int {
	if skill == nil {
		return 0
	}
	base := skill.HealAmount()
	if base <= 0 {
		return 0
	}
	return base + skill.EmpowerDamageBonus(empDmg)
}

// selectAoEMonsters implementa la selección oficial de objetivos AoE de sbe_BlastObject (0x00583BC0)
// y sbe_BlastLocate (0x00584160):
//   - radio = Area / 2 (shr edx, 1 en 0x00583DD5 / 0x0058430F)
//   - sector cónico si 0 < SkillDegree < 360
//   - ordenamiento por distancia ascendente (qsort 0x00583E93 / 0x00584338)
//   - tope de objetivos TargetNumber ([eax+0x107])
func selectAoEMonsters(zMap *ZoneMap, casterX, casterY, centerX, centerY int, skill *SkillDef) []*ZoneMonster {
	if zMap == nil {
		return nil
	}
	radius := 200.0
	if skill != nil && skill.Area > 0 {
		radius = float64(skill.Area) / 2.0
	}
	radiusSq := radius * radius

	useCone := skill != nil && skill.SkillDegree > 0 && skill.SkillDegree < 360
	halfConeRad := 0.0
	dirAngle := 0.0
	if useCone {
		halfConeRad = (float64(skill.SkillDegree) / 2.0) * math.Pi / 180.0
		dirAngle = math.Atan2(float64(centerY-casterY), float64(centerX-casterX))
		if skill.DirectionRotate != 0 {
			dirAngle += float64(skill.DirectionRotate) * math.Pi / 180.0
		}
	}

	type mobDist struct {
		mob    *ZoneMonster
		distSq float64
	}
	var candidates []mobDist
	for _, mob := range zMap.GetAllMonsters() {
		if mob == nil {
			continue
		}
		mob.mu.RLock()
		dead := mob.IsDead || mob.CurHP <= 0
		mob.mu.RUnlock()
		if dead {
			continue
		}
		mx, my := mob.GetPosition()
		dx := float64(mx - centerX)
		dy := float64(my - centerY)
		dSq := dx*dx + dy*dy
		if dSq > radiusSq {
			continue
		}
		if useCone && dSq > 0 {
			ang := math.Atan2(float64(my-casterY), float64(mx-casterX))
			diff := math.Abs(ang - dirAngle)
			for diff > math.Pi {
				diff = math.Abs(2*math.Pi - diff)
			}
			if diff > halfConeRad {
				continue
			}
		}
		candidates = append(candidates, mobDist{mob: mob, distSq: dSq})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].distSq < candidates[j].distSq
	})

	limit := len(candidates)
	if skill != nil && skill.TargetNumber > 0 && skill.TargetNumber < limit {
		limit = skill.TargetNumber
	}

	out := make([]*ZoneMonster, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, candidates[i].mob)
	}
	return out
}

// resolveSkillBlastOnEntity resuelve un impacto individual de smo_SkillBlast (0x00581160):
// soporta curación (isHeal=1), daño físico/mágico (RulesOfEngagement), AbState 0x16D (1-dmg),
// absorción 0xC9/0xCA y aplicación de hasta 4 AbStates (StaNameA..D) con empower de duración.
func (h *Handler) resolveSkillBlastOnEntity(caster *ZonePlayer, target GridEntity, skill *SkillDef, empDmg, empKeep uint8) (skillDamageTarget, bool) {
	if caster == nil || target == nil {
		return skillDamageTarget{}, false
	}

	// 1. Curación (Zone.exe 0x005817CE..0x00581980 + bono PassiveSkill HealUPRate)
	if healAmt := calcOfficialSkillHeal(skill, empDmg); healAmt > 0 {
		if h != nil && h.dataStore != nil && h.dataStore.SkillCatalog != nil {
			if pb := h.dataStore.SkillCatalog.ComputePassiveBonuses(caster); pb.HealUpRate > 0 {
				healAmt += (healAmt * pb.HealUpRate) / 1000
			}
		}
		if tp, ok := target.(*ZonePlayer); ok && tp.Char != nil {
			oldHP := tp.Char.CurHP
			tp.Char.CurHP += healAmt
			if tp.Char.CurHP > tp.MaxHP {
				tp.Char.CurHP = tp.MaxHP
			}
			actualHealed := tp.Char.CurHP - oldHP
			if actualHealed < 0 {
				actualHealed = 0
			}
			restHP := uint32(tp.Char.CurHP)
			if tp.Client != nil {
				hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
				hpCmd.WriteUint32(restHP)
				hpCmd.WriteUint16(tp.NextHPChangeOrder())
				_ = tp.Client.SendPacket(hpCmd)
			}
			appliedAb := h.applySkillAbStatesToPlayer(tp, skill, empKeep)
			return skillDamageTarget{
				Handle:    tp.MapObjectID,
				Damage:    uint32(actualHealed),
				RestHP:    restHP,
				IsHeal:    true,
				IsAbState: appliedAb,
			}, false
		}
	}

	// 1b. Buff / Party-Buff / Debuff puro (sin daño ni heal directo, ej. Resist 1 / MentalShout01)
	if skill != nil && skill.IsPureBuff() {
		switch te := target.(type) {
		case *ZonePlayer:
			appliedAb := h.applySkillAbStatesToPlayer(te, skill, empKeep)
			var restHP uint32
			if te.Char != nil {
				restHP = uint32(te.Char.CurHP)
			}
			return skillDamageTarget{
				Handle:    te.MapObjectID,
				Damage:    0,
				RestHP:    restHP,
				IsAbState: appliedAb,
			}, false
		case *ZoneMonster:
			appliedAb := h.applySkillAbStatesToMonster(caster.CurrentMap, te, skill, empKeep)
			return skillDamageTarget{
				Handle:    te.Handle,
				Damage:    0,
				RestHP:    uint32(te.CurHP),
				IsAbState: appliedAb,
			}, false
		}
	}

	// 2. Daño físico o mágico (RulesOfEngagementPhisycalSkill / MagicalSkill + bonos de PassiveSkill.shn)
	atk := calcOfficialSkillAttackPower(caster, skill, empDmg)
	isMag := skill != nil && skill.IsMagical()
	if h != nil && h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		pb := h.dataStore.SkillCatalog.ComputePassiveBonuses(caster)
		if pb.WeaponFlatPlus > 0 {
			atk += pb.WeaponFlatPlus
		}
		rateBoost := pb.WeaponRatePermille
		if isMag {
			rateBoost += pb.MARateUp
		} else {
			rateBoost += pb.WCRateUp
		}
		if rateBoost > 0 {
			atk += (atk * rateBoost) / 1000
		}
	}

	var def int
	var isOneDmgShield bool
	var absorbToHeal bool
	switch te := target.(type) {
	case *ZoneMonster:
		if isMag {
			def = int(te.MR)
		} else {
			def = int(te.AC)
		}
		if te.ActiveStates != nil {
			if _, ok := te.ActiveStates[0x16D]; ok {
				isOneDmgShield = true
			}
			if (!isMag && te.ActiveStates[0xC9] != nil) || (isMag && te.ActiveStates[0xCA] != nil) {
				absorbToHeal = true
			}
		}
	case *ZonePlayer:
		if isMag {
			def = int(clampStatUint32(te.EquipMR, te.AbMR)) + int(te.Spr)
		} else {
			def = int(clampStatUint32(te.EquipAC, te.AbAC)) + int(te.End)
		}
		if te.HasAbState(0x16D) {
			isOneDmgShield = true
		}
		if (!isMag && te.HasAbState(0xC9)) || (isMag && te.HasAbState(0xCA)) {
			absorbToHeal = true
		}
	}

	dmg := atk - def
	if tp, ok := target.(*ZonePlayer); ok && h != nil && h.dataStore != nil && h.dataStore.SkillCatalog != nil {
		if pb := h.dataStore.SkillCatalog.ComputePassiveBonuses(tp); pb.DMGMinusRate > 0 && pb.DMGMinusRate < 1000 {
			dmg -= (dmg * pb.DMGMinusRate) / 1000
		}
	}
	if dmg < 1 {
		dmg = 1
	}

	// Crítico (Zone.exe 0x00505D7C..0x00505E07)
	isCrit := false
	if skill != nil && !isOneDmgShield {
		if CheckCritPermille(ROECalcCritRatePermille(50, 0, 0)) {
			isCrit = true
			dmg = (dmg * 1500) / 1000
		}
	}
	if isOneDmgShield {
		dmg = 1
	}

	// Absorción daño->heal por AbState 0xC9 (phys) / 0xCA (mag) (Zone.exe 0x00581F26..0x00581FB0)
	if absorbToHeal {
		if tp, ok := target.(*ZonePlayer); ok && tp.Char != nil {
			tp.Char.CurHP += dmg
			if tp.Char.CurHP > tp.MaxHP {
				tp.Char.CurHP = tp.MaxHP
			}
			restHP := uint32(tp.Char.CurHP)
			if tp.Client != nil {
				hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
				hpCmd.WriteUint32(restHP)
				hpCmd.WriteUint16(tp.NextHPChangeOrder())
				_ = tp.Client.SendPacket(hpCmd)
			}
			return skillDamageTarget{
				Handle: tp.MapObjectID,
				Damage: uint32(dmg),
				RestHP: restHP,
				IsHeal: true,
			}, false
		}
	}

	// Aplicar daño
	var restHP uint32
	var isDead bool
	var appliedAb bool
	switch te := target.(type) {
	case *ZoneMonster:
		restHP, isDead = caster.CurrentMap.ApplyDamageToMonster(te, caster, dmg)
		if !isDead {
			appliedAb = h.applySkillAbStatesToMonster(caster.CurrentMap, te, skill, empKeep)
		}
		return skillDamageTarget{
			Handle:    te.Handle,
			Damage:    uint32(dmg),
			RestHP:    restHP,
			IsCrit:    isCrit,
			IsDead:    isDead,
			IsAbState: appliedAb,
		}, isDead
	case *ZonePlayer:
		if te.Char != nil {
			dmg = te.AbsorbDamageWithShield(caster.CurrentMap, dmg)
			if dmg > 0 {
				if te.Char.CurHP <= dmg {
					te.Char.CurHP = 0
					te.State = 3
					isDead = true
				} else {
					te.Char.CurHP -= dmg
				}
			}
			restHP = uint32(te.Char.CurHP)
			if te.Client != nil {
				hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
				hpCmd.WriteUint32(restHP)
				hpCmd.WriteUint16(te.NextHPChangeOrder())
				_ = te.Client.SendPacket(hpCmd)
			}
			if !isDead {
				appliedAb = h.applySkillAbStatesToPlayer(te, skill, empKeep)
			}
		}
		return skillDamageTarget{
			Handle:    te.MapObjectID,
			Damage:    uint32(dmg),
			RestHP:    restHP,
			IsCrit:    isCrit,
			IsDead:    isDead,
			IsAbState: appliedAb,
		}, isDead
	}

	return skillDamageTarget{}, false
}

// resolveAbStateInfo obtiene ID, duración base, metadatos y acciones desde Abstate.shn + SubAbstate.shn
// (con fallback al repositorio JSON data/AbState para los 10 estados legacy).
func (h *Handler) resolveAbStateInfo(name string, strength int) (id uint16, dur time.Duration, def *AbStateDef, sub *SubAbStateDef, ok bool) {
	if h == nil || h.dataStore == nil {
		return 0, 0, nil, nil, false
	}
	if h.dataStore.SkillCatalog != nil {
		if d := h.dataStore.SkillCatalog.LookupAbState(name); d != nil {
			ms := 10000
			s := h.dataStore.SkillCatalog.LookupSubAbState(d.SubAbState, strength)
			if s != nil && s.KeepTimeMs > 0 {
				ms = s.KeepTimeMs
			}
			return d.ID, time.Duration(ms) * time.Millisecond, d, s, true
		}
	}
	if h.dataStore.AbState != nil {
		if cfg, found := h.dataStore.AbState.GetByName(name); found && cfg != nil {
			d := 10 * time.Second
			if lvlCfg, err := cfg.GetLevelConfig(strength); err == nil && lvlCfg != nil && lvlCfg.DurationMs > 0 {
				d = time.Duration(lvlCfg.DurationMs) * time.Millisecond
			}
			return cfg.ID, d, nil, nil, true
		}
	}
	return 0, 0, nil, nil, false
}

// applyOneAbStateToPlayer agrega un AbState al jugador, aplica los modificadores de SubAbstate.shn,
// sincroniza stats en cliente y envía los paquetes oficiales S2C:
// PROTO_NC_BAT_ABSTATEINFORM_CMD (0x2429, 8B) al jugador y PROTO_NC_BAT_ABSTATESET_CMD (0x2426, 6B) al AOI.
func (h *Handler) applyOneAbStateToPlayer(target *ZonePlayer, abID uint16, name string, strength int, dur time.Duration, isBuff bool, sub *SubAbStateDef) {
	if target == nil {
		return
	}
	st := &ActiveAbState{
		ID:        abID,
		Name:      name,
		Level:     strength,
		Duration:  dur,
		ExpiresAt: time.Now().Add(dur),
		IsBuff:    isBuff,
		IsDebuff:  !isBuff,
	}
	if sub != nil {
		applySubAbStateActions(st, sub.Actions)
	}
	target.AddAbState(st)
	h.syncPlayerAbStateParams(target)

	if target.Client != nil {
		inform := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD)
		inform.WriteUint32(uint32(abID))
		inform.WriteUint32(uint32(dur.Milliseconds()))
		_ = target.Client.SendPacket(inform)
	}
	if target.CurrentMap != nil {
		// Zone.exe so_AbnormalState_BroadcastSet (VA 0x00405738): difunde 0x2427 excluyendo al objetivo,
		// porque el objetivo ya recibió 0x2429 con el keepTimeMs exacto (enviarle 0x2427 pisa y borra el timer).
		setCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATESET_CMD)
		setCmd.WriteUint16(target.MapObjectID)
		setCmd.WriteUint32(uint32(abID))
		target.CurrentMap.broadcastSkillAOI(setCmd, target.MapObjectID, target.MapObjectID)
	}
}

// applySkillAbStatesToPlayer aplica los hasta 4 AbStates (StaNameA..D) definidos en la skill
// sumando el bono de duración de empower track 2 (Zone.exe 0x005824A7..0x00582791),
// enviando 0x2429 + 0x2427 y propagando PartyState1 a miembros de party en rango sin duplicar el aura Main.
func (h *Handler) applySkillAbStatesToPlayer(target *ZonePlayer, skill *SkillDef, empKeep uint8) bool {
	if target == nil || skill == nil || h.dataStore == nil {
		return false
	}
	applied := false
	extraDur := time.Duration(skill.EmpowerKeepTimeBonus(empKeep)) * time.Millisecond
	for _, entry := range skill.AbStates {
		if entry.Name == "" || entry.Name == "-" || entry.SucRate <= 0 {
			continue
		}
		if entry.SucRate < 1000 && rand.Intn(1000) >= entry.SucRate {
			continue
		}
		abID, dur, def, sub, ok := h.resolveAbStateInfo(entry.Name, entry.Strength)
		if !ok {
			continue
		}
		dur += extraDur

		// Si el estado es un portador de aura de party (ej. StaMentalShoutMain con PartyState1 = StaMentalShout),
		// aplicar únicamente PartyState1 visible al lanzador y a los miembros de la party en PartyRange
		// para evitar duplicar el ícono con descripción "-".
		if def != nil && def.PartyStates[0] != "" && def.PartyStates[0] != "-" {
			if pAbID, _, pDef, pSub, pOK := h.resolveAbStateInfo(def.PartyStates[0], entry.Strength); pOK {
				pIsBuff := pDef == nil || pDef.DispelIndex != 2
				h.applyOneAbStateToPlayer(target, pAbID, def.PartyStates[0], entry.Strength, dur, pIsBuff, pSub)
				applied = true

				if h.charRepo != nil && target.Char != nil && target.CurrentMap != nil {
					if partyNo, _, _, err := h.charRepo.GetPartyByChar(target.Char.ID); err == nil && partyNo > 0 {
						if members, err := h.charRepo.GetPartyMembers(partyNo); err == nil {
							tx, ty := target.GetPosition()
							maxRange := float64(def.PartyRange)
							if maxRange <= 0 {
								maxRange = 1600
							}
							for _, m := range members {
								if m == nil || m.CharNo == target.Char.ID {
									continue
								}
								mp := h.findPlayerInMapByCharID(target.CurrentMap, m.CharNo)
								if mp == nil || mp.Char == nil || mp.State == 3 {
									continue
								}
								mx, my := mp.GetPosition()
								if math.Hypot(float64(mx-tx), float64(my-ty)) <= maxRange {
									h.applyOneAbStateToPlayer(mp, pAbID, def.PartyStates[0], entry.Strength, dur, pIsBuff, pSub)
								}
							}
						}
					}
				}
				continue
			}
		}

		isBuff := def == nil || def.DispelIndex != 2
		h.applyOneAbStateToPlayer(target, abID, entry.Name, entry.Strength, dur, isBuff, sub)
		applied = true
	}
	return applied
}

// applySkillAbStatesToMonster aplica los hasta 4 AbStates de la skill a un monstruo vivo
// y difunde PROTO_NC_BAT_ABSTATESET_CMD (0x2426, 6B) al AOI.
func (h *Handler) applySkillAbStatesToMonster(zMap *ZoneMap, mob *ZoneMonster, skill *SkillDef, empKeep uint8) bool {
	if mob == nil || skill == nil || h.dataStore == nil {
		return false
	}
	applied := false
	extraDur := time.Duration(skill.EmpowerKeepTimeBonus(empKeep)) * time.Millisecond
	for _, entry := range skill.AbStates {
		if entry.Name == "" || entry.Name == "-" || entry.SucRate <= 0 {
			continue
		}
		if entry.SucRate < 1000 && rand.Intn(1000) >= entry.SucRate {
			continue
		}
		abID, dur, _, sub, ok := h.resolveAbStateInfo(entry.Name, entry.Strength)
		if !ok {
			continue
		}
		dur += extraDur
		st := &ActiveAbState{
			ID:        abID,
			Name:      entry.Name,
			Level:     entry.Strength,
			Duration:  dur,
			ExpiresAt: time.Now().Add(dur),
			IsDebuff:  true,
		}
		if sub != nil {
			applySubAbStateActions(st, sub.Actions)
		}
		mob.mu.Lock()
		if mob.ActiveStates == nil {
			mob.ActiveStates = make(map[uint16]*ActiveAbState)
		}
		mob.ActiveStates[abID] = st
		mob.mu.Unlock()
		if zMap != nil {
			setCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATESET_CMD)
			setCmd.WriteUint16(mob.Handle)
			setCmd.WriteUint32(uint32(abID))
			zMap.Broadcast(setCmd, 0)
		}
		applied = true
	}
	return applied
}
