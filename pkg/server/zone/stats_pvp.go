package zone

import (
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/database"
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

	//Arreglo HP/SP en cliente
	h.syncPlayerVitals(player)

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

// sendCharDeadMenu envía PROTO_NC_CHAR_DEADMENU_CMD (0x104D, 9 bytes) para abrir la ventana de resurrección.
func sendCharDeadMenu(p *ZonePlayer) {
	if p == nil || p.Client == nil {
		return
	}
	pkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_DEADMENU_CMD)
	pkt.WriteUint32(300) // Second (300s = 5 min antes de respawn automático)
	pkt.WriteUint8(1)    // priority
	pkt.WriteUint32(0)   // eMenuType (0 = menú normal de resurrección)
	_ = p.Client.SendPacket(pkt)
}

// handleCharReviveReq implementa PROTO_NC_CHAR_REVIVE_REQ (0x104E).
// Verificado en Zone.exe sp_NC_CHAR_REVIVE_REQ (0x00452980) -> sp_ReviveReqProcess (0x00451BD0):
// restaura (MaxHP * 0x133) >> 10 (~30%), limpia debuffs y emite REVIVESAME_CMD (0x104F) o REVIVEOTHER_CMD (0x1050).
func (h *Handler) handleCharReviveReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil {
		return
	}

	// Idempotencia: solo un jugador muerto (State == 3) puede revivir. Evita que REVIVE_REQ duplicados
	// (varios DEADMENU recibidos, doble click) disparen warps repetidos mientras el cliente aún carga el mapa.
	if player.State != 3 {
		return
	}

	// 1. Restaurar estado de vida: (MaxHP * 0x133) >> 10 (~30% en Zone.exe 0x00452046)
	player.State = 1 // Normal
	player.Char.CurHP = (player.MaxHP * 0x133) >> 10
	if player.Char.CurHP <= 0 {
		player.Char.CurHP = 50
	}
	player.IsAutoAttacking = false
	player.TargetHandle = 0

	// Limpiar debuffs / estados incapacitantes al revivir
	for id, st := range player.ActiveStates {
		if st != nil && (st.IsDebuff || st.IsStun || st.IsRoot || st.IsSilence) {
			player.RemoveAbState(id)
		}
	}

	// 2. Resolver mapa y coordenadas de respawn desde FieldList.json (nunca vacío)
	targetMapName, spawnX, spawnY := h.resolveReviveDestination(player)

	// 3. Traslado: si es el mismo mapa, revivir in-place con REVIVESAME_CMD (0x104F) + broadcast AOI
	// sin disparar sendInGameChunk; si es otro mapa, usar WarpPlayer (y fallback a REVIVESAME_CMD si se rechaza).
	sendReviveSame := func(x, y int) {
		sameCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVESAME_CMD)
		sameCmd.WriteUint16(h.resolveClientMapID(player.CurrentMap.MapID, player.CurrentMap.ID))
		sameCmd.WriteInt32(int32(x))
		sameCmd.WriteInt32(int32(y))
		_ = client.SendPacket(sameCmd)
	}

	if strings.EqualFold(targetMapName, player.CurrentMap.MapID) {
		player.Char.XPos = spawnX
		player.Char.YPos = spawnY
		player.Char.TCharacterState.XPos = spawnX
		player.Char.TCharacterState.YPos = spawnY
		if h.charRepo != nil {
			_ = h.charRepo.SaveCharLocation(player.Char.ID, player.CurrentMap.MapID, spawnX, spawnY, player.Rotation)
		}
		sendReviveSame(spawnX, spawnY)

		stopPw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESTOP_CMD)
		stopPw.WriteUint16(player.MapObjectID)
		stopPw.WriteInt32(int32(spawnX))
		stopPw.WriteInt32(int32(spawnY))
		player.CurrentMap.syncMoveAOI(player, spawnX, spawnY, stopPw)
	} else if !h.WarpPlayer(player, targetMapName, spawnX, spawnY, player.Rotation) {
		// Warp rechazado (p. ej. condición de nivel): revivir en el lugar con REVIVESAME_CMD (0x104F, 10 bytes).
		sendReviveSame(player.Char.XPos, player.Char.YPos)
	}

	// 4. Sincronizar barra de vida al cliente
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(player.NextHPChangeOrder())
	_ = client.SendPacket(hpCmd)
}

// resolveReviveDestination determina el mapa y coordenadas de resurrección según FieldList.json.
// Orden de resolución del mapa: Respawn.City del field actual -> CurrentMap.MapID -> Char.TownName.
// Cada candidato se valida contra el FieldStore (si existe) para no devolver jamás un MapID vacío
// o inexistente, que en WarpPlayer/GetOrCreateMap caería al field de Serial 0.
// Las coordenadas salen de Respawn.Spots[0] del field cuya configuración de respawn se usó;
// si no hay spot válido y el destino es el mismo mapa, se conservan las coordenadas actuales.
func (h *Handler) resolveReviveDestination(player *ZonePlayer) (string, int, int) {
	curX, curY := player.Char.XPos, player.Char.YPos
	curMapID := ""
	if player.CurrentMap != nil {
		curMapID = strings.TrimSpace(player.CurrentMap.MapID)
	}
	townName := strings.TrimSpace(player.Char.TownName)

	valid := func(id string) bool {
		if id == "" {
			return false
		}
		if h.fieldStore == nil {
			return true
		}
		f, ok := h.fieldStore.GetField(id)
		return ok && f != nil
	}

	// Field de origen: el FieldEntry cacheado en el mapa o, si falta, el del FieldStore.
	var entry *field.FieldEntry
	if player.CurrentMap != nil {
		entry = player.CurrentMap.FieldEntry
	}
	if entry == nil && h.fieldStore != nil {
		for _, id := range []string{curMapID, townName} {
			if id == "" {
				continue
			}
			if f, ok := h.fieldStore.GetField(id); ok && f != nil {
				entry = f
				break
			}
		}
	}

	target := ""
	x, y := curX, curY
	if entry != nil {
		city := strings.TrimSpace(entry.Respawn.City)
		if valid(city) {
			target = city
			for _, spot := range entry.Respawn.Spots {
				if spot.X > 0 && spot.Y > 0 {
					x, y = spot.X, spot.Y
					break
				}
			}
		}
	}
	if target == "" {
		for _, id := range []string{curMapID, townName} {
			if valid(id) {
				target = id
				break
			}
		}
	}
	if target == "" {
		// Último recurso: mapa actual tal cual (aunque no esté en FieldList) para no caer en Serial 0.
		target = curMapID
		if target == "" {
			target = townName
		}
	}
	// Si el destino es otro mapa y no hubo spot, usar el spot de respawn propio del mapa destino.
	if (x == curX && y == curY) && !strings.EqualFold(target, curMapID) && h.fieldStore != nil {
		if f, ok := h.fieldStore.GetField(target); ok && f != nil {
			for _, spot := range f.Respawn.Spots {
				if spot.X > 0 && spot.Y > 0 {
					x, y = spot.X, spot.Y
					break
				}
			}
		}
	}
	return target, x, y
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
	hpCmd.WriteUint16(player.NextHPChangeOrder()) // hpchangeorder
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

	sendFail := func(errCode uint16) {
		failCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REVIVEITEMUSEFAIL_CMD)
		failCmd.WriteUint16(errCode)
		_ = client.SendPacket(failCmd)
	}

	// Verificado en Zone.exe sp_NC_ITEM_REVIVEITEMUSE_CMD (0x00417580):
	// debe estar muerto (State == 3 -> 0x0701 si vivo); si existe itemStore/charRepo requiere ítem de resurrección en Box 9 (0x071E).
	if player.State != 3 {
		sendFail(0x0701)
		return
	}

	if h.itemStore != nil && h.charRepo != nil {
		items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
		var reviveItem *database.Item
		for _, it := range items {
			if it == nil || it.StorageType != 9 {
				continue
			}
			if info, ok := h.itemStore.GetItemInfo(it.ItemID); ok && info != nil {
				low := strings.ToLower(info.InxName)
				if strings.Contains(low, "revive") || strings.Contains(low, "resurrect") || strings.Contains(low, "rebirth") {
					reviveItem = it
					break
				}
			}
		}
		if reviveItem == nil || !h.consumeItemUnits(player, reviveItem, 1) {
			sendFail(0x071E)
			return
		}
	}

	player.State = 1
	player.Char.CurHP = player.MaxHP
	if player.Char.CurHP <= 0 {
		player.Char.CurHP = 100
	}
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(player.NextHPChangeOrder()) // hpchangeorder
	_ = client.SendPacket(hpCmd)
}
