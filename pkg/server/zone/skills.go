package zone

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// SpecialSkillConfig permite configurar o resolver los IDs especiales de Zone.exe
// (las 5 skills maestras de producción [0x13286CFC..0x13286D04], el ítem Empower Reset [0x14D50496]
// y el ítem Skill Erase Scroll [0x14D504AE]).
type SpecialSkillConfig struct {
	ProductionSkills  [5]uint16
	EmpowerResetItem  uint16
	SkillEraseItem    uint16
}

func (h *Handler) getSpecialSkillConfig() *SpecialSkillConfig {
	if h != nil && h.dataStore != nil && h.dataStore.SpecialSkills != nil {
		return h.dataStore.SpecialSkills
	}
	return nil
}

// isPlayerStateValidForEmpower valida la bitmask oficial 0x0076 (Zone.exe 0x0044F1F4):
// estados permitidos 1, 2, 4, 5, 6 (rechaza estado 3 = Dead -> 0x0882).
func isPlayerStateValidForEmpower(state byte) bool {
	return state == 1 || state == 2 || state == 4 || state == 5 || state == 6
}

// isPlayerStateValidForProductField valida la bitmask oficial 0x0016 (Zone.exe 0x0046C8CA):
// estados permitidos 1, 2, 4 (rechaza Dead=3, Vendor=5, Mount=6 -> 0x0B04).
func isPlayerStateValidForProductField(state byte) bool {
	return state == 1 || state == 2 || state == 4
}

// handleSkillEmpowAllocReq implementa PROTO_NC_SKILL_EMPOWALLOC_REQ (0x4811, VA 0x0044F150).
// Formato oficial Zone.md (6B): skill (uint16) + plus (SKILL_EMPOWER uint16) + minus (SKILL_EMPOWER uint16).
// Responde PROTO_NC_SKILL_EMPOWALLOC_ACK (0x4814, 4B: skill + err) y PROTO_NC_SKILL_EMPOWPOINT_CMD (0x4810, 1B: restempow).
func (h *Handler) handleSkillEmpowAllocReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	skillIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	plusMask, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	minusMask, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	sendAck := func(sk uint16, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK)
		ack.WriteUint16(sk)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	// 1. Estado del jugador (Zone.exe 0x0044F1F4: bitmask 0x0076 -> skill=0xFFFF, err=0x0882)
	if !isPlayerStateValidForEmpower(player.State) {
		sendAck(0xFFFF, 0x0882)
		return
	}

	// 2. Conflicto de nibbles (Zone.exe 0x0044F233..0x0044F2A6 -> err 0x0889)
	var plusNibbles [4]uint8
	var minusNibbles [4]uint8
	totalPlus := 0
	totalMinus := 0
	for i := uint(0); i < 4; i++ {
		pNibble := uint8((plusMask >> (4 * i)) & 0xF)
		mNibble := uint8((minusMask >> (4 * i)) & 0xF)
		if pNibble > 0 && mNibble > 0 {
			sendAck(skillIndex, 0x0889)
			return
		}
		plusNibbles[i] = pNibble
		minusNibbles[i] = mNibble
		totalPlus += int(pNibble)
		totalMinus += int(mNibble)
	}

	// 3. Si minusMask > 0, verificar ítems de Empower Reset en inventario box 9 (Zone.exe 0x0044F2A8..0x0044F3CE -> 0x0886)
	var resetItemsToConsume []*database.Item
	if totalMinus > 0 {
		cfg := h.getSpecialSkillConfig()
		var resetItemID uint16
		if cfg != nil {
			resetItemID = cfg.EmpowerResetItem
		}
		if resetItemID == 0 {
			sendAck(skillIndex, 0x0886)
			return
		}
		items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
		avail := 0
		for _, it := range items {
			if it != nil && it.StorageType == 9 && it.ItemID == resetItemID {
				lot := int(it.Flags)
				if lot <= 0 {
					lot = 1
				}
				avail += lot
				resetItemsToConsume = append(resetItemsToConsume, it)
			}
		}
		if avail < totalMinus {
			sendAck(skillIndex, 0x0886)
			return
		}
	}

	// 4. Validar tope global de puntos de empower (Zone.exe 0x0044F3D5..0x0044F470 -> 0x0883)
	// y límites por track 0..5 / underflow / skill aprendida (csl_EmpowerAlloc 0x00448570 -> 0x0884)
	var curTracks [4]uint8
	if player.LearnedSkills != nil {
		sk := player.GetLearnedSkill(skillIndex)
		if sk == nil {
			sendAck(skillIndex, 0x0884)
			return
		}
		curTracks = [4]uint8{sk.SkillPowerDemage, sk.SkillPowerSP, sk.SkillPowerKeepTime, sk.SkillPowerCoolTime}

		netDelta := totalPlus - totalMinus
		if player.TotalEmpowerUsed()+netDelta > player.MaxEmpowerPoints() {
			sendAck(skillIndex, 0x0883)
			return
		}
		for i := 0; i < 4; i++ {
			if int(curTracks[i])+int(plusNibbles[i]) > 5 || curTracks[i] < minusNibbles[i] {
				sendAck(skillIndex, 0x0884)
				return
			}
		}
	} else {
		for i := 0; i < 4; i++ {
			if plusNibbles[i] > 5 || minusNibbles[i] > 0 {
				sendAck(skillIndex, 0x0884)
				return
			}
		}
	}

	// 5. Consumir ítems de Empower Reset si totalMinus > 0 (Zone.exe 0x0044F508..0x0044F724)
	if totalMinus > 0 {
		rem := uint32(totalMinus)
		for _, it := range resetItemsToConsume {
			if rem == 0 {
				break
			}
			take := uint32(it.Flags)
			if take == 0 {
				take = 1
			}
			if take > rem {
				take = rem
			}
			if !h.consumeItemUnits(player, it, take) {
				sendAck(skillIndex, 0x0886)
				return
			}
			rem -= take
		}
	}

	// 6. Aplicar cambios a los 4 tracks en memoria + DB y auditar en GameLog Tipo E
	var newTracks [4]uint8
	for i := 0; i < 4; i++ {
		val := int(curTracks[i]) + int(plusNibbles[i]) - int(minusNibbles[i])
		if val < 0 {
			val = 0
		}
		if val > 5 {
			val = 5
		}
		newTracks[i] = uint8(val)
		if plusNibbles[i] > 0 || minusNibbles[i] > 0 {
			_ = h.charRepo.SetSingleSkillPower(player.Char.ID, int16(skillIndex), uint8(i), newTracks[i])
			if h.gameLogDispatcher != nil {
				mapName := ""
				if player.CurrentMap != nil {
					mapName = player.CurrentMap.MapID
				}
				x, y := player.GetPosition()
				h.gameLogDispatcher.EnqueueTypeE(14, mapName, x, y, 0, int(skillIndex), 0, i, int(newTracks[i]))
			}
		}
	}
	if sk := player.GetLearnedSkill(skillIndex); sk != nil {
		sk.SkillPowerDemage = newTracks[0]
		sk.SkillPowerSP = newTracks[1]
		sk.SkillPowerKeepTime = newTracks[2]
		sk.SkillPowerCoolTime = newTracks[3]
	}

	sendAck(skillIndex, 0x0881) // Success (gds_NC_SKILL_EMPOWALLOC_DB_ACK VA 0x0044F87E)

	// PROTO_NC_SKILL_EMPOWPOINT_CMD (0x4810, sizeof = 1 byte: restempow unsigned char)
	restEmpow := player.MaxEmpowerPoints() - player.TotalEmpowerUsed()
	if restEmpow < 0 {
		restEmpow = 0
	}
	if restEmpow > 255 {
		restEmpow = 255
	}
	empowPointCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD)
	empowPointCmd.WriteUint8(uint8(restEmpow))
	_ = client.SendPacket(empowPointCmd)
}

// handleSkillProductFieldReq implementa sp_NC_SKILL_PRODUCTFIELD_REQ (0x4822, VA 0x0046C820).
// Selecciona y aprende una profesión de producción (máximo 2 de las 5 maestras).
// Formato oficial S2C PROTO_NC_SKILL_PRODUCTFIELD_ACK (0x4823, 4B): mainskill (uint16) + err (uint16).
func (h *Handler) handleSkillProductFieldReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	skillID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_ACK)
		ack.WriteUint16(skillID)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	// 1. Bitmask de estado 0x0016 (Zone.exe 0x0046C8CA -> 0x0B04)
	if !isPlayerStateValidForProductField(player.State) {
		sendAck(0x0B04)
		return
	}

	// 2. Whitelist de las 5 skills maestras de producción (Zone.exe 0x0046C936..0x0046C98B -> 0x0B12)
	cfg := h.getSpecialSkillConfig()
	if cfg != nil && cfg.ProductionSkills != ([5]uint16{}) {
		validProd := false
		for _, pid := range cfg.ProductionSkills {
			if pid != 0 && pid == skillID {
				validProd = true
				break
			}
		}
		if !validProd {
			sendAck(0x0B12)
			return
		}
	}

	// 3. Ya aprendida (Zone.exe 0x0046C990: csl_IsLearn -> 0x0B13)
	if player.HasSkill(skillID) {
		sendAck(0x0B13)
		return
	}

	// 4. Máximo 2 profesiones de producción aprendidas (Zone.exe 0x0046C9AA..0x0046CA38 -> 0x0B14)
	if cfg != nil && cfg.ProductionSkills != ([5]uint16{}) {
		learnedCount := 0
		for _, pid := range cfg.ProductionSkills {
			if pid != 0 && player.HasSkill(pid) {
				learnedCount++
			}
		}
		if learnedCount >= 2 {
			sendAck(0x0B14)
			return
		}
	}

	// 5. Aprender y persistir la profesión de producción
	if player.Char != nil {
		newSk := &database.CharacterSkill{
			CharNo:     player.Char.ID,
			SkillNo:    skillID,
			SkillLevel: 1,
		}
		player.SetLearnedSkill(newSk)
		if h.charRepo != nil {
			_ = h.charRepo.SaveCharacterSkill(player.Char.ID, newSk)
		}
	}

	// Auditoría oficial en GameLog Tipo E
	if h.gameLogDispatcher != nil {
		x, y := player.GetPosition()
		h.gameLogDispatcher.EnqueueTypeE(15, player.CurrentMap.MapID, x, y, 0, int(skillID), 0, 1, 0)
	}

	sendAck(0x0B01) // Success (Zone.exe VA 0x0046CA3D)
}

// handleSkillUnlearnReq implementa sp_NC_SKILL_UNLEARN_REQ (0x4824, VA 0x0057AD80).
// Formato oficial C2S (3B): slot (uint8) + skill (uint16).
// Respuesta S2C PROTO_NC_SKILL_UNLEARN_ACK (0x4827, 4B): skill (uint16) + err (uint16).
func (h *Handler) handleSkillUnlearnReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}
	skillIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_ACK)
		ack.WriteUint16(skillIndex)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	// 1. Verificar si la skill está aprendida (Zone.exe 0x0057AEA3: csl_IsLearn -> 0x1503)
	if player.LearnedSkills != nil && !player.HasSkill(skillIndex) {
		sendAck(0x1503)
		return
	}

	// 2. Validar slot e ítem Skill Erase Scroll en inventario box 9 (Zone.exe 0x0057AEE8..0x0057B000)
	cfg := h.getSpecialSkillConfig()
	if cfg != nil && cfg.SkillEraseItem != 0 {
		if slot >= 96 {
			sendAck(0x0702)
			return
		}
		it := h.findItemAt(player.Char.ID, 9, int16(slot))
		if it == nil || it.ItemID != cfg.SkillEraseItem {
			sendAck(0x0702)
			return
		}
		if it.Flags < 0 {
			sendAck(0x0703)
			return
		}
		if !h.consumeItemUnits(player, it, 1) {
			sendAck(0x0703)
			return
		}
	}

	player.RemoveLearnedSkill(skillIndex)
	_ = h.charRepo.DeleteCharacterSkill(player.Char.ID, int16(skillIndex))

	errCode := uint16(0x1501) // Active skill unlearn (Zone.exe VA 0x0057F269)
	if (h.dataStore != nil && h.dataStore.IsPassiveSkill(skillIndex)) || skillIndex >= 0x7530 {
		errCode = 0x1505 // Passive skill unlearn (Zone.exe VA 0x0057F2A6)
	}

	sendAck(errCode)
}

// handleSkillReplyReviveCmd implementa sp_NC_SKILL_REPLYREVIVE_CMD (0x482A, VA 0x418A30) -> sp_ReviveNow (VA 0x57D290).
// Envía PROTO_NC_SKILL_REVIVE_CMD (0x481E, 12B: hp, sp, lp), SOMEONEREVIVE_CMD (0x481F, 14B: reviver, hp, sp, lp)
// y los cambios de barras HPCHANGE/SPCHANGE.
func (h *Handler) handleSkillReplyReviveCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	// Verificado en Zone.exe sp_NC_SKILL_REPLYREVIVE_CMD (0x00418A30):
	// si el personaje está vivo (sin estado muerto/pendiente de resurrección), responde PROTO_NC_SKILL_REPLYREVIVEFAIL_CMD (0x482B, 0x1506).
	if player.State != 3 && player.Char.CurHP > 0 {
		failCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVEFAIL_CMD)
		failCmd.WriteUint16(0x1506)
		_ = client.SendPacket(failCmd)
		return
	}

	// Si el jugador estaba muerto, revivir con 50% de HP y limpiar estados negativos
	// exe: 0x0057D354 (sp_ReviveNow -> sp_KQElementFunction(kqpi_KnockDownMinus))
	h.sp_KQElementFunction(player, (*KQPlayerInform).kqpi_KnockDownMinus)
	player.State = 1 // Normal
	player.Char.CurHP = player.MaxHP / 2
	if player.Char.CurHP <= 0 {
		player.Char.CurHP = 50
	}
	for id, st := range player.ActiveStates {
		if st != nil && (st.IsDebuff || st.IsStun || st.IsRoot || st.IsSilence) {
			player.RemoveAbState(id)
		}
	}

	hp := uint32(player.Char.CurHP)
	sp := uint32(player.Char.CurSP)
	lp := uint32(0)

	reviveCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REVIVE_CMD)
	reviveCmd.WriteUint32(hp)
	reviveCmd.WriteUint32(sp)
	reviveCmd.WriteUint32(lp)
	_ = client.SendPacket(reviveCmd)

	if player.CurrentMap != nil {
		someRevive := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD)
		someRevive.WriteUint16(player.MapObjectID)
		someRevive.WriteUint32(hp)
		someRevive.WriteUint32(sp)
		someRevive.WriteUint32(lp)
		player.CurrentMap.Broadcast(someRevive, 0)
	}

	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(hp)
	hpCmd.WriteUint16(player.NextHPChangeOrder())
	_ = client.SendPacket(hpCmd)

	spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
	spCmd.WriteUint32(sp)
	_ = client.SendPacket(spCmd)
}

// handleItemUseActiveSkillReq implementa PROTO_NC_ITEM_USE_ACTIVESKILL_REQ (0x304E).
// Struct oficial PDB: PROTO_NC_ITEM_USE_ACTIVESKILL_REQ (11 bytes: nInvenSlot uint8, nTargetHnd uint16, TargetLoc SHINE_XY_TYPE 8B).
// Aprende una habilidad activa a partir del uso de un pergamino en inventario.
func (h *Handler) handleItemUseActiveSkillReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}
	targetHnd, _ := pkt.Reader.ReadUint16()
	targetX, _ := pkt.Reader.ReadInt32()
	targetY, _ := pkt.Reader.ReadInt32()
	_, _, _ = targetHnd, targetX, targetY

	sendUseSkillAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	// Verificado en Zone.exe sp_NC_ITEM_USE_ACTIVESKILL_REQ (0x00540B90):
	// estado inválido (muerto/vendor) -> 0x3129; ítem inexistente -> 0x312C; éxito -> 0x3128.
	if player.State == 3 || player.State == 5 {
		sendUseSkillAck(0x3129)
		return
	}

	it := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if it == nil {
		sendUseSkillAck(0x312C)
		return
	}

	// Consumir pergamino de inventario (Box 9)
	if !h.consumeItemUnits(player, it, 1) {
		sendUseSkillAck(0x312C)
		return
	}

	// Persistir la nueva habilidad aprendida en base de datos mediante CALL p_Skill_Set
	newSkill := &database.CharacterSkill{
		CharNo:     player.Char.ID,
		SkillNo:    it.ItemID,
		SkillLevel: 1,
	}
	_ = h.charRepo.SaveCharacterSkill(player.Char.ID, newSkill)

	sucCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SKILL_LEARNSUC_CMD)
	sucCmd.WriteUint16(it.ItemID)
	_ = client.SendPacket(sucCmd)
	sendUseSkillAck(0x3128)
}
