package zone

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// handleSkillEmpowAllocReq implementa PROTO_NC_SKILL_EMPOWALLOC_REQ (0x4811).
// Persiste el punto de empoderamiento mediante CALL p_Skill_SetPower y audita en GameLog Tipo E.
func (h *Handler) handleSkillEmpowAllocReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	skillIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	empowIndex, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}
	point, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}

	// Validar y persistir usando el procedimiento oficial p_Skill_SetPower
	_ = h.charRepo.SetSingleSkillPower(player.Char.ID, int16(skillIndex), empowIndex, point)

	// Auditoría oficial en GameLog Tipo E
	if h.gameLogDispatcher != nil {
		mapName := ""
		if player.CurrentMap != nil {
			mapName = player.CurrentMap.MapID
		}
		x, y := player.GetPosition()
		h.gameLogDispatcher.EnqueueTypeE(14, mapName, x, y, 0, int(skillIndex), 0, int(empowIndex), int(point))
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_ACK)
	ack.WriteUint16(skillIndex)
	ack.WriteUint8(empowIndex)
	ack.WriteUint8(point)
	_ = client.SendPacket(ack)

	empowPointCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWPOINT_CMD)
	empowPointCmd.WriteUint16(0) // Remaining empow points
	_ = client.SendPacket(empowPointCmd)
}

// handleSkillProductFieldReq implementa PROTO_NC_SKILL_PRODUCTFIELD_REQ (0x4822).
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
	_ = client.SendPacket(ack)
}

// handleSkillUnlearnReq implementa PROTO_NC_SKILL_UNLEARN_REQ (0x4824).
// Elimina la habilidad invocando el procedimiento oficial CALL p_Skill_Del.
func (h *Handler) handleSkillUnlearnReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	skillIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	_ = h.charRepo.DeleteCharacterSkill(player.Char.ID, int16(skillIndex))

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_ACK)
	ack.WriteUint16(skillIndex)
	_ = client.SendPacket(ack)
}

// handleSkillReplyReviveCmd implementa PROTO_NC_SKILL_REPLYREVIVE_CMD (0x482A).
// Procesa la aceptación de resurrección ofrecida por una habilidad aliada.
func (h *Handler) handleSkillReplyReviveCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	// Si el jugador estaba muerto, revivir con 50% de HP
	if player.State == 3 || player.Char.CurHP <= 0 {
		player.State = 1 // Normal
		player.Char.CurHP = player.MaxHP / 2
		if player.Char.CurHP <= 0 {
			player.Char.CurHP = 50
		}

		reviveCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REVIVE_CMD)
		_ = client.SendPacket(reviveCmd)

		if player.CurrentMap != nil {
			someRevive := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SOMEONEREVIVE_CMD)
			someRevive.WriteUint16(player.MapObjectID)
			player.CurrentMap.Broadcast(someRevive, 0)
		}

		hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpCmd.WriteUint32(uint32(player.Char.CurHP))
		hpCmd.WriteUint16(0)
		_ = client.SendPacket(hpCmd)
	}
}

// handleItemUseActiveSkillReq implementa PROTO_NC_ITEM_USE_ACTIVESKILL_REQ (0x304E).
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
	_ = targetHnd

	// Persistir la nueva habilidad aprendida en base de datos mediante CALL p_Skill_Set
	newSkill := &database.CharacterSkill{
		CharNo:     player.Char.ID,
		SkillNo:    uint16(invenSlot),
		SkillLevel: 1,
	}
	_ = h.charRepo.SaveCharacterSkill(player.Char.ID, newSkill)

	sucCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SKILL_LEARNSUC_CMD)
	sucCmd.WriteUint16(uint16(invenSlot))
	_ = client.SendPacket(sucCmd)
}
