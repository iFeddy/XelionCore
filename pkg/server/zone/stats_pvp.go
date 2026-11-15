package zone

import (
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// handleCharStatIncPointReq implementa PROTO_NC_CHAR_STAT_INCPOINT_REQ (0x105C).
// Incrementa el atributo solicitado (0=Str, 1=Con, 2=Dex, 3=Int, 4=Men), persiste mediante CALL p_Char_StatPoint_Inc y audita en GameLog Tipo 2.
func (h *Handler) handleCharStatIncPointReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	stat, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}

	if player.Char.TCharacter.RedistributePoint <= 0 {
		failAck := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_INCPOINTFAIL_ACK)
		failAck.WriteUint8(stat)
		_ = client.SendPacket(failAck)
		return
	}

	player.Char.TCharacter.RedistributePoint--
	switch stat {
	case 0:
		player.Char.TCharacter.Strength++
		player.Str = byte(player.Char.TCharacter.Strength)
	case 1:
		player.Char.TCharacter.Constitute++
		player.End = byte(player.Char.TCharacter.Constitute)
	case 2:
		player.Char.TCharacter.Dexterity++
		player.Dex = byte(player.Char.TCharacter.Dexterity)
	case 3:
		player.Char.TCharacter.Intelligence++
		player.Int = byte(player.Char.TCharacter.Intelligence)
	case 4:
		player.Char.TCharacter.MentalPower++
		player.Spr = byte(player.Char.TCharacter.MentalPower)
	}

	_ = h.charRepo.IncStatPoint(player.Char.ID, stat)

	// Recalcular HP/SP máximos con FreeStatTable y ClassParam
	if h.dataStore != nil && h.dataStore.CombatStore != nil {
		freeBonus := h.dataStore.CombatStore.GetFreeStatBonus(
			player.Char.TCharacter.Strength,
			player.Char.TCharacter.Intelligence,
			player.Char.TCharacter.Dexterity,
			player.Char.TCharacter.Constitute,
			player.Char.TCharacter.MentalPower,
		)
		player.MaxHP = player.LevelParam.MaxHP + freeBonus.MaxHP
		player.MaxSP = player.LevelParam.MaxSP + freeBonus.MaxSP
	}

	// Auditoría oficial en GameLog Tipo 2
	if h.gameLogDispatcher != nil {
		mapName := ""
		if player.CurrentMap != nil {
			mapName = player.CurrentMap.MapID
		}
		x, y := player.GetPosition()
		h.gameLogDispatcher.EnqueueType2(2, player.Char.ID, mapName, x, y, 0, int(stat), player.Char.TCharacter.RedistributePoint)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_INCPOINTSUC_ACK)
	sucAck.WriteUint8(stat)
	_ = client.SendPacket(sucAck)

	remainCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_REMAINPOINT_CMD)
	remainCmd.WriteUint8(uint8(player.Char.TCharacter.RedistributePoint))
	_ = client.SendPacket(remainCmd)
}

// handleCharStatDecPointReq implementa PROTO_NC_CHAR_STAT_DECPOINT_REQ (0x1062).
// Decrementa el atributo indicado, persiste mediante CALL p_Char_StatPoint_Dec y audita en GameLog Tipo 2.
func (h *Handler) handleCharStatDecPointReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	stat, err := pkt.Reader.ReadUint8()
	if err != nil {
		return
	}

	switch stat {
	case 0:
		if player.Char.TCharacter.Strength <= 0 {
			return
		}
		player.Char.TCharacter.Strength--
		player.Str = byte(player.Char.TCharacter.Strength)
	case 1:
		if player.Char.TCharacter.Constitute <= 0 {
			return
		}
		player.Char.TCharacter.Constitute--
		player.End = byte(player.Char.TCharacter.Constitute)
	case 2:
		if player.Char.TCharacter.Dexterity <= 0 {
			return
		}
		player.Char.TCharacter.Dexterity--
		player.Dex = byte(player.Char.TCharacter.Dexterity)
	case 3:
		if player.Char.TCharacter.Intelligence <= 0 {
			return
		}
		player.Char.TCharacter.Intelligence--
		player.Int = byte(player.Char.TCharacter.Intelligence)
	case 4:
		if player.Char.TCharacter.MentalPower <= 0 {
			return
		}
		player.Char.TCharacter.MentalPower--
		player.Spr = byte(player.Char.TCharacter.MentalPower)
	}

	player.Char.TCharacter.RedistributePoint++
	_ = h.charRepo.DecStatPoint(player.Char.ID, stat)

	if h.dataStore != nil && h.dataStore.CombatStore != nil {
		freeBonus := h.dataStore.CombatStore.GetFreeStatBonus(
			player.Char.TCharacter.Strength,
			player.Char.TCharacter.Intelligence,
			player.Char.TCharacter.Dexterity,
			player.Char.TCharacter.Constitute,
			player.Char.TCharacter.MentalPower,
		)
		player.MaxHP = player.LevelParam.MaxHP + freeBonus.MaxHP
		player.MaxSP = player.LevelParam.MaxSP + freeBonus.MaxSP
	}

	if h.gameLogDispatcher != nil {
		mapName := ""
		if player.CurrentMap != nil {
			mapName = player.CurrentMap.MapID
		}
		x, y := player.GetPosition()
		h.gameLogDispatcher.EnqueueType2(2, player.Char.ID, mapName, x, y, 0, int(stat), player.Char.TCharacter.RedistributePoint)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_DECPOINTSUC_ACK)
	sucAck.WriteUint8(stat)
	_ = client.SendPacket(sucAck)

	remainCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_REMAINPOINT_CMD)
	remainCmd.WriteUint8(uint8(player.Char.TCharacter.RedistributePoint))
	_ = client.SendPacket(remainCmd)
}

// handleCharReviveReq implementa PROTO_NC_CHAR_REVIVE_REQ (0x104E).
// Verificado en Zone.exe sp_NC_CHAR_REVIVE_REQ (0x00452980) -> sp_ReviveReqProcess (0x00451BD0):
// restaura (MaxHP * 0x133) >> 10 (~30%), limpia debuffs y emite REVIVESAME_CMD (0x104F) o REVIVEOTHER_CMD (0x1050).
func (h *Handler) handleCharReviveReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil {
		return
	}

	// Restaurar estado de vida: (MaxHP * 0x133) >> 10 (~30% en Zone.exe 0x00452046)
	player.State = 1 // Normal
	player.Char.CurHP = (player.MaxHP * 0x133) >> 10
	if player.Char.CurHP <= 0 {
		player.Char.CurHP = 50
	}
	// ab_Recover / as_Remove: limpiar debuffs/estados incapacitantes al revivir
	for id, st := range player.ActiveStates {
		if st != nil && (st.IsDebuff || st.IsStun || st.IsRoot || st.IsSilence) {
			player.RemoveAbState(id)
		}
	}

	// Coordenadas y mapa de respawn desde FieldList.json si están configurados
	targetMapName := ""
	if player.CurrentMap.FieldEntry != nil {
		if player.CurrentMap.FieldEntry.Respawn.City != "" {
			targetMapName = player.CurrentMap.FieldEntry.Respawn.City
		}
		if len(player.CurrentMap.FieldEntry.Respawn.Spots) > 0 {
			spot := player.CurrentMap.FieldEntry.Respawn.Spots[0]
			if spot.X > 0 && spot.Y > 0 {
				player.Char.XPos = spot.X
				player.Char.YPos = spot.Y
			}
		}
	}

	var targetMap *ZoneMap = player.CurrentMap
	if targetMapName != "" && player.CurrentMap.MapID != "" && targetMapName != player.CurrentMap.MapID && h.mapManager != nil {
		if rm := h.mapManager.GetMap(targetMapName); rm != nil {
			targetMap = rm
		}
	}

	if targetMap != nil && targetMap != player.CurrentMap {
		// PROTO_NC_CHAR_REVIVEOTHER_CMD (0x1050, 30B): mapid(2), x(4), y(4), ip[16], port(2), linkflag/id(2)
		otherCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVEOTHER_CMD)
		otherCmd.WriteUint16(targetMap.ID)
		otherCmd.WriteInt32(int32(player.Char.XPos))
		otherCmd.WriteInt32(int32(player.Char.YPos))
		writeFixedStr(otherCmd, "127.0.0.1", 16)
		otherCmd.WriteUint16(9013)
		otherCmd.WriteUint16(player.MapObjectID)
		_ = client.SendPacket(otherCmd)
	} else {
		// 1. Notificar al jugador que resucitó en el mismo mapa (PROTO_NC_CHAR_REVIVESAME_CMD: 10 bytes = mapid[2], location[8])
		sameCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVESAME_CMD)
		sameCmd.WriteUint16(player.CurrentMap.ID)
		sameCmd.WriteInt32(int32(player.Char.XPos))
		sameCmd.WriteInt32(int32(player.Char.YPos))
		_ = client.SendPacket(sameCmd)
	}

	// 2. Actualizar barra de vida
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0)
	_ = client.SendPacket(hpCmd)
}

// handleCharRebirthCmd implementa PROTO_NC_CHAR_REBIRTH_CMD (0x1102).
func (h *Handler) handleCharRebirthCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	player.State = 1
	player.Char.CurHP = player.MaxHP / 2
	if player.Char.CurHP <= 0 {
		player.Char.CurHP = 50
	}

	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0)
	_ = client.SendPacket(hpCmd)
}

// handleCharRebirthRejectCmd implementa PROTO_NC_CHAR_REBIRTH_REJECT_CMD (0x1103).
func (h *Handler) handleCharRebirthRejectCmd(client *network.Client, pkt *protocol.Packet) {
	// Rechazo silencioso de resurrección de clérigo
}

// handleItemReviveItemUseCmd implementa PROTO_NC_ITEM_REVIVEITEMUSE_CMD (0x3039).
func (h *Handler) handleItemReviveItemUseCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	player.State = 1
	player.Char.CurHP = player.MaxHP
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0)
	_ = client.SendPacket(hpCmd)
}
