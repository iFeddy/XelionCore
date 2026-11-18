package zone

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// handleSkillEmpowAllocReq implementa PROTO_NC_SKILL_EMPOWALLOC_REQ (0x4811, VA 0x44F150).
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

	// Zone.exe 0x44F233..0x44F2A6: ningún nibble (4 bits) puede incrementarse y decrementarse a la vez (err 0x0889).
	for i := uint(0); i < 4; i++ {
		pNibble := (plusMask >> (4 * i)) & 0xF
		mNibble := (minusMask >> (4 * i)) & 0xF
		if pNibble > 0 && mNibble > 0 {
			ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK)
			ack.WriteUint16(skillIndex)
			ack.WriteUint16(0x0889)
			_ = client.SendPacket(ack)
			return
		}
	}

	// Aplicar cada track de empoderamiento (0..3) mediante p_Skill_SetPower y auditar en GameLog Tipo E.
	for i := uint(0); i < 4; i++ {
		pNibble := uint8((plusMask >> (4 * i)) & 0xF)
		if pNibble == 0 {
			continue
		}
		_ = h.charRepo.SetSingleSkillPower(player.Char.ID, int16(skillIndex), uint8(i), pNibble)
		if h.gameLogDispatcher != nil {
			mapName := ""
			if player.CurrentMap != nil {
				mapName = player.CurrentMap.MapID
			}
			x, y := player.GetPosition()
			h.gameLogDispatcher.EnqueueTypeE(14, mapName, x, y, 0, int(skillIndex), 0, int(i), int(pNibble))
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK)
	ack.WriteUint16(skillIndex)
	ack.WriteUint16(0x0881) // Success (gds_NC_SKILL_EMPOWALLOC_DB_ACK VA 0x44F87E)
	_ = client.SendPacket(ack)

	// PROTO_NC_SKILL_EMPOWPOINT_CMD (0x4810, sizeof = 1 byte: restempow unsigned char)
	empowPointCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD)
	empowPointCmd.WriteUint8(0)
	_ = client.SendPacket(empowPointCmd)
}

// handleSkillProductFieldReq implementa sp_NC_SKILL_PRODUCTFIELD_REQ (0x4822, VA 0x46C820).
// Formato oficial S2C PROTO_NC_SKILL_PRODUCTFIELD_ACK (0x4823, 4B): mainskill (uint16) + err (uint16 = 0x0B01).
func (h *Handler) handleSkillProductFieldReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	skillID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	// Auditoría oficial en GameLog Tipo E (invocación de campo mágico/tótem)
	if h.gameLogDispatcher != nil {
		x, y := player.GetPosition()
		h.gameLogDispatcher.EnqueueTypeE(15, player.CurrentMap.MapID, x, y, 0, int(skillID), 0, 1, 0)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_ACK)
	ack.WriteUint16(skillID)
	ack.WriteUint16(0x0B01) // Success (Zone.exe VA 0x46CA3D)
	_ = client.SendPacket(ack)
}

// handleSkillUnlearnReq implementa sp_NC_SKILL_UNLEARN_REQ (0x4824, VA 0x57AD80).
// Formato oficial C2S (3B): slot (uint8) + skill (uint16).
// Respuesta S2C PROTO_NC_SKILL_UNLEARN_ACK (0x4827, 4B): skill (uint16) + err (uint16 = 0x1501 / 0x1505).
func (h *Handler) handleSkillUnlearnReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}
	_ = slot
	skillIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	_ = h.charRepo.DeleteCharacterSkill(player.Char.ID, int16(skillIndex))

	errCode := uint16(0x1501) // Active skill unlearn (Zone.exe VA 0x57F269)
	if skillIndex >= 0x7530 {
		errCode = 0x1505 // Passive skill unlearn (Zone.exe VA 0x57F2A6)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_ACK)
	ack.WriteUint16(skillIndex)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
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
