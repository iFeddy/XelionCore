package zone

import (
	"log"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ============================================================================
// Dept 2: Diagnóstico y Parámetros de Servidor (PINGTEST / SERVERPARAMETER)
// ============================================================================

// handleMiscPingTestClientZoneDB implementa sp_NC_MISC_PINGTEST_CLIENT_ZONE_DB (0x0812).
func (h *Handler) handleMiscPingTestClientZoneDB(client *network.Client, pkt *protocol.Packet) {
	if client == nil || pkt == nil {
		return
	}
	raw := make([]byte, len(pkt.Raw))
	copy(raw, pkt.Raw)

	nowTick := uint32(time.Now().UnixMilli() & 0xFFFFFFFF)
	// PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB (sizeof = 90 bytes)
	// Stamps ZoneRecvFromClient (offset 26), ZoneSendToDB (14), DBRecvFromZone (30),
	// DBSendToZone (34), ZoneRecvFromDB (18), ZoneSendToClient (22)
	if len(raw) >= 46 {
		putUint32LE(raw[26:30], nowTick)
		putUint32LE(raw[14:18], nowTick)
		putUint32LE(raw[30:34], nowTick)
		putUint32LE(raw[34:38], nowTick)
		putUint32LE(raw[18:22], nowTick)
		putUint32LE(raw[22:26], uint32(time.Now().UnixMilli()&0xFFFFFFFF))
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB)
	ack.WriteBytes(raw)
	_ = client.SendPacket(ack)
}

// handleMiscPingTestToolWMClientZoneDB implementa sp_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB (0x0813).
func (h *Handler) handleMiscPingTestToolWMClientZoneDB(client *network.Client, pkt *protocol.Packet) {
	if client == nil || pkt == nil {
		return
	}
	raw := make([]byte, len(pkt.Raw))
	copy(raw, pkt.Raw)

	nowTick := uint32(time.Now().UnixMilli() & 0xFFFFFFFF)
	// PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB (sizeof = 130 bytes)
	if len(raw) >= 86 {
		putUint32LE(raw[66:70], nowTick) // ZoneRecvFromClient
		putUint32LE(raw[54:58], nowTick) // ZoneSendToDB
		putUint32LE(raw[70:74], nowTick) // DBRecvFromZone
		putUint32LE(raw[74:78], nowTick) // DBSendToZone
		putUint32LE(raw[58:62], nowTick) // ZoneRecvFromDB
		putUint32LE(raw[62:66], nowTick) // ZoneSendToClient
	}
	h.DispatchWMS("wms_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB", map[string]interface{}{
		"tick": nowTick,
	})

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB)
	ack.WriteBytes(raw)
	_ = client.SendPacket(ack)
}

// handleMiscServerParameterReq implementa sp_NC_MISC_SERVERPARAMETER_REQ (0x082F).
func (h *Handler) handleMiscServerParameterReq(client *network.Client, pkt *protocol.Packet) {
	if client == nil {
		return
	}
	// PROTO_NC_MISC_SERVERPARAMETER_ACK (3B): Shoutlevel (1B) + Shoutdelay (2B)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SERVERPARAMETER_ACK)
	ack.WriteByte(1)
	ack.WriteUint16(30)
	_ = client.SendPacket(ack)
}

// ============================================================================
// Dept 3 & 4: Sesión de Cuenta, Registro, Beauty Shop y Logout
// ============================================================================

// handleUserNormalLogoutCmd implementa sp_NC_USER_NORMALLOGOUT_CMD (0x0C18) en Zone.exe (0x00558690).
// LogoutType: 0 = salida completa del cliente, 1 = retorno a Selección de Personaje en WorldService.
func (h *Handler) handleUserNormalLogoutCmd(client *network.Client, pkt *protocol.Packet) {
	if client == nil {
		return
	}
	logoutType, _ := pkt.Reader.ReadByte()
	if player, ok := h.getPlayer(client); ok && player != nil && player.Char != nil {
		log.Printf("[Zone][%s] PROTO_NC_USER_NORMALLOGOUT_CMD (type=%d) for '%s'",
			client.RemoteAddr(), logoutType, player.Char.Name)
		player.LogoutPending = true
		player.LogoutType = logoutType
		h.DispatchWMS("wms_NC_USER_CONNECTCUT2ZONE_CMD", map[string]interface{}{
			"charID":     player.Char.ID,
			"userNo":     player.Char.AccountID,
			"charName":   player.Char.Name,
			"type":       logoutType,
			"logoutType": logoutType,
		})
	}
	h.OnClientDisconnect(client)
	client.Close()
}

// handleUserRegiseNumberReq implementa sp_NC_USER_REGISENUMBER_REQ (0x0C23).
func (h *Handler) handleUserRegiseNumberReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	var regNum uint32 = uint32(player.Char.AccountID)
	if h.accountRepo != nil {
		if acc, err := h.accountRepo.GetAccountByUserNo(player.Char.AccountID); err == nil && acc != nil && acc.ID > 0 {
			regNum = uint32(acc.ID)
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_USER_REGISENUMBER_ACK)
	ack.WriteUint32(regNum)
	_ = client.SendPacket(ack)
}

// handleUserUseBeautyShopCmd implementa sp_NC_USER_USE_BEAUTY_SHOP_CMD (0x0C32).
func (h *Handler) handleUserUseBeautyShopCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	player.BeautyShopOpen = true
	if pkt != nil && pkt.Reader.Remaining() >= 3 {
		hair, _ := pkt.Reader.ReadByte()
		hairColor, _ := pkt.Reader.ReadByte()
		face, _ := pkt.Reader.ReadByte()
		player.Char.Hair = hair
		player.Char.HairColor = hairColor
		player.Char.Face = face
		player.Char.TCharacterShape.CharNo = player.Char.ID
		player.Char.TCharacterShape.Hair = hair
		player.Char.TCharacterShape.HairColor = hairColor
		player.Char.TCharacterShape.Face = face
		if h.charRepo != nil {
			_ = h.charRepo.UpdateCharacter(player.Char)
		}
		if h.gameLogDispatcher != nil && player.CurrentMap != nil {
			h.gameLogDispatcher.EnqueueType4(6, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(hair), int(hairColor), int(face))
		}
		h.DispatchWMS("wms_NC_CHAR_ZONE_SHAPE_CMD", map[string]interface{}{
			"charID":    player.Char.ID,
			"hair":      hair,
			"hairColor": hairColor,
			"face":      face,
		})
	}
}

// handleCharLogoutReadyCmd implementa sp_NC_CHAR_LOGOUTREADY_CMD (0x1071) en Zone.exe (0x00558450).
func (h *Handler) handleCharLogoutReadyCmd(client *network.Client, pkt *protocol.Packet) {
	if client == nil {
		return
	}
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	player.LogoutPending = true
	if player.CurrentMap != nil && player.CurrentMap.MapID != "" {
		player.Char.TownName = player.CurrentMap.MapID
	}
	if h.charRepo != nil {
		_ = h.charRepo.SaveCharLocation(player.Char.ID, player.Char.TownName, player.Char.XPos, player.Char.YPos, player.Rotation)
		_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)
		_ = h.charRepo.UpdateCharacter(player.Char)
	}
	h.DispatchWMS("wms_NC_USER_GET_PLAY_MIN_ACK", map[string]interface{}{
		"charID": player.Char.ID,
		"userNo": player.Char.AccountID,
	})
}

// handleCharLogoutCancelCmd implementa sp_NC_CHAR_LOGOUTCANCEL_CMD (0x1072).
func (h *Handler) handleCharLogoutCancelCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.LogoutPending = false
}

// handleCharSaveLinkReq implementa sp_NC_CHAR_SAVE_LINK_REQ (0x10B9 -> usp_LinkScroll_SetDestMap).
func (h *Handler) handleCharSaveLinkReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	destMap := ""
	if player.CurrentMap != nil {
		destMap = player.CurrentMap.MapID
	}
	if pkt != nil && pkt.Reader.Remaining() > 0 {
		rawLen := pkt.Reader.Remaining()
		if rawLen > 12 {
			rawLen = 12
		}
		if s := strings.TrimSpace(readFixedStr(pkt.Reader, rawLen)); s != "" {
			destMap = s
		}
	}
	if destMap != "" {
		_ = h.charRepo.SaveLinkScrollDestMap(player.Char.ID, destMap)
		player.Char.KQMapName = destMap
		h.DispatchWMS("wms_NC_CHAR_ZONE_LINK_FROM_CMD", map[string]interface{}{
			"charID":  player.Char.ID,
			"destMap": destMap,
		})
	}
}

// ============================================================================
// Dept 8: Modos de Acción, Shout y Roar Mundial
// ============================================================================

// handleActChangeModeReq implementa sp_NC_ACT_CHANGEMODE_REQ (0x2008).
func (h *Handler) handleActChangeModeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	mode, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	if mode != 1 && mode != 2 {
		return
	}
	player.ActionMode = mode
	if player.CurrentMap != nil {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONECHANGEMODE_CMD)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteByte(mode)
		player.CurrentMap.Broadcast(cmd, 0)
	}
}

// handleActShoutCmd implementa sp_NC_ACT_SHOUT_CMD (0x201E -> 0x201F PROTO_NC_ACT_SOMEONESHOUT_CMD).
func (h *Handler) handleActShoutCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	linkCount, _ := pkt.Reader.ReadByte()
	msgLen, _ := pkt.Reader.ReadByte()
	var content string
	if msgLen > 0 && pkt.Reader.Remaining() > 0 {
		readN := int(msgLen)
		if readN > pkt.Reader.Remaining() {
			readN = pkt.Reader.Remaining()
		}
		content, _ = pkt.Reader.ReadString(readN)
	}
	var linkData []byte
	if linkCount > 0 {
		expectedLinks := int(linkCount) * 103
		if pkt.Reader.Remaining() < expectedLinks {
			return
		}
		linkData, _ = pkt.Reader.ReadBytes(expectedLinks)
	}
	trimmed := strings.TrimSpace(content)
	if trimmed == "" && linkCount == 0 {
		return
	}
	if linkCount == 0 {
		content = trimmed
	}

	var flag byte
	if (player.Client != nil && player.Client.Admin > 0) || player.Char.AdminLevel > 0 {
		flag |= 0x01 // bit 0 = GMColor
	}

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD)
	cmd.WriteByte(linkCount)
	cmd.WriteString(player.Char.Name, 20)
	cmd.WriteByte(flag)
	cmd.WriteByte(byte(len(content)))
	cmd.WriteBytes([]byte(content))
	if linkCount > 0 && len(linkData) > 0 {
		cmd.WriteBytes(linkData)
	}
	if player.CurrentMap != nil {
		player.CurrentMap.Broadcast(cmd, 0)
	} else {
		_ = client.SendPacket(cmd)
	}
	h.DispatchWMS("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
		"type":    "shout",
		"sender":  player.Char.Name,
		"message": trimmed,
	})
}

// handleActRoarReq implementa sp_NC_ACT_ROAR_REQ (0x204B -> 0x204C PROTO_NC_ACT_ROAR_ACK).
func (h *Handler) handleActRoarReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	slot, _ := pkt.Reader.ReadByte()
	msgLen, _ := pkt.Reader.ReadByte()
	var content string
	if msgLen > 0 && pkt.Reader.Remaining() > 0 {
		readN := int(msgLen)
		if readN > pkt.Reader.Remaining() {
			readN = pkt.Reader.Remaining()
		}
		content, _ = pkt.Reader.ReadString(readN)
	}

	var errCode uint16 = 0
	if h.charRepo != nil {
		if it := h.findItemAt(player.Char.ID, 9, int16(slot)); it != nil {
			if it.Flags > 1 {
				it.Flags--
				_ = h.charRepo.SaveCharacterItem(it)
				h.sendCellChangeFull(client, 9, int16(slot), it, 0xFFFF)
			} else {
				_ = h.charRepo.DeleteItem(it.ItemKey)
				h.sendCellChangeFull(client, 9, int16(slot), nil, 0xFFFF)
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_ROAR_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	content = strings.TrimSpace(content)
	if content != "" {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD)
		cmd.WriteByte(0)
		cmd.WriteString(player.Char.Name, 20)
		cmd.WriteByte(1)
		cmd.WriteByte(byte(len(content)))
		cmd.WriteBytes([]byte(content))
		if player.CurrentMap != nil {
			player.CurrentMap.Broadcast(cmd, 0)
		}
		h.DispatchWMS("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
			"type":    "roar",
			"sender":  player.Char.Name,
			"message": content,
		})
	}
}

// ============================================================================
// Dept 14, 16 & 21: Party MemberInfo, UI State (MysteryBolt) & Friend Point
// ============================================================================

// handlePartyMemberInfoReqCmd implementa sp_NC_PARTY_MEMBERINFOREQ_CMD (0x3848 -> 0x3832/0x384F).
func (h *Handler) handlePartyMemberInfoReqCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	var partyNo int
	if h.charRepo != nil {
		pNo, _, _, err := h.charRepo.GetPartyByChar(player.Char.ID)
		if err == nil && pNo > 0 {
			partyNo = pNo
		}
	}
	var memberCount byte
	if h.charRepo != nil && partyNo > 0 {
		if members, err := h.charRepo.GetPartyMembers(partyNo); err == nil {
			memberCount = byte(len(members))
		}
	}
	h.DispatchWMS("wms_NC_PARTY_MEMBERINFORM_ACK", map[string]interface{}{
		"charID":      player.Char.ID,
		"partyNo":     partyNo,
		"memberCount": memberCount,
	})

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_PARTY_MEMBERINFORM_CMD)
	cmd.WriteByte(memberCount)
	_ = client.SendPacket(cmd)
}

// handleCharSaveUIStateSaveReq implementa sp_NC_CHARSAVE_UI_STATE_SAVE_REQ (0x4025 -> usp_Character_setMysteryBolt).
func (h *Handler) handleCharSaveUIStateSaveReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	mysteryType, _ := pkt.Reader.ReadByte()
	var errCode uint16 = 0
	if h.charRepo != nil {
		if ret, err := h.charRepo.SaveMysteryBolt(player.Char.ID, mysteryType); err != nil || ret != 0 {
			errCode = 1
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleFriendUseFriendPointReq implementa sp_NC_FRIEND_UES_FRIEND_POINT_REQ (0x5421 -> usp_Character_getFriendpoint/setFriendpoint).
func (h *Handler) handleFriendUseFriendPointReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	costPts := 10
	if pkt != nil && pkt.Reader.Remaining() >= 2 {
		if reqPts, err := pkt.Reader.ReadUint16(); err == nil && reqPts > 0 {
			costPts = int(reqPts)
		}
	}

	curPts, err := h.charRepo.GetFriendPoints(player.Char.ID)
	var (
		errCode    uint16 = 0
		remPts     uint16 = 0
		rewardItem uint16 = 0
		rewardNum  uint8  = 0
	)
	if err != nil || curPts < costPts {
		errCode = 1
		if curPts > 0 {
			remPts = uint16(curPts)
		}
	} else {
		newPts := curPts - costPts
		if setErr := h.charRepo.SetFriendPoints(player.Char.ID, newPts); setErr != nil {
			errCode = 1
			remPts = uint16(curPts)
		} else {
			remPts = uint16(newPts)
			rewardItem = 8004
			rewardNum = 1
			itemKey := gmItemKeySeq.Add(1)
			if slot, err := h.charRepo.CreateItemUseEmptyStorage(itemKey, player.Char.ID, 9, int(rewardItem), int(rewardNum)); err == nil && slot >= 0 {
				it := h.findItemAt(player.Char.ID, 9, slot)
				h.sendCellChangeFull(client, 9, slot, it, 0xFFFF)
			}
			if h.gameLogDispatcher != nil && player.CurrentMap != nil {
				h.gameLogDispatcher.EnqueueType7(4, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, player.Char.ID, int(rewardItem), costPts)
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(remPts)
	ack.WriteUint16(rewardItem)
	ack.WriteByte(rewardNum)
	_ = client.SendPacket(ack)
}

// ============================================================================
// Dept 29, 45 & 53: Guild Storage / Rename / Retype, ChangeCharID & Pet Ask
// ============================================================================

// handleGuildStorageOpenReq implementa sp_NC_GUILD_STORAGEOPEN_REQ (0x7491).
func (h *Handler) handleGuildStorageOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_STORAGEOPEN_ACK)
	g, _, err := h.charRepo.GetGuildByCharID(player.Char.ID)
	if err != nil || g == nil {
		ack.WriteUint16(errItemStorageFail)
		ack.WriteUint64(0)
		ack.WriteUint64(0)
		ack.WriteByte(0)
		_ = client.SendPacket(ack)
		return
	}
	token, _ := h.charRepo.GetGuildToken(g.ID)
	items, _ := h.charRepo.GetItemsByStorage(g.ID, 16)
	if len(items) > 250 {
		items = items[:250]
	}

	player.GuildStorageOpen = true
	ack.WriteUint16(errItemStorageSuccess)
	ack.WriteUint64(g.Money)
	ack.WriteUint64(uint64(token))
	ack.WriteByte(byte(len(items)))
	for _, it := range items {
		ack.WriteBytes(h.buildItemPacketInform(it))
	}
	_ = client.SendPacket(ack)
}

// handleItemDBCloseGuildStorageCmd implementa sp_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD (0x74A0).
func (h *Handler) handleItemDBCloseGuildStorageCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.GuildStorageOpen = false
}

// handleGuildRenameReq implementa sp_NC_GUILD_RENAME_REQ (0x74AC -> p_Guild_setName).
func (h *Handler) handleGuildRenameReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil || pkt.Reader.Remaining() < 17 {
		return
	}
	slot, _ := pkt.Reader.ReadByte()
	newName := strings.TrimSpace(readFixedStr(pkt.Reader, 16))

	var errCode uint16 = 0
	var guildID int
	g, mbr, err := h.charRepo.GetGuildByCharID(player.Char.ID)
	if err != nil || g == nil || mbr == nil || mbr.Grade != 0 || newName == "" {
		errCode = 1
	} else {
		guildID = g.ID
		ret, rErr := h.charRepo.RenameGuild(guildID, newName)
		if rErr != nil || ret != 0 {
			errCode = 1
		} else {
			if it := h.findItemAt(player.Char.ID, 9, int16(slot)); it != nil {
				_ = h.charRepo.DeleteItem(it.ItemKey)
				h.sendCellChangeFull(client, 9, int16(slot), nil, 0xFFFF)
			}
			cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_RENAME_CMD)
			cmd.WriteUint32(uint32(guildID))
			cmd.WriteString(newName, 16)
			if player.CurrentMap != nil {
				player.CurrentMap.Broadcast(cmd, 0)
			} else {
				_ = client.SendPacket(cmd)
			}
			h.DispatchWMS("wms_NC_GUILD_WORLD_RENAME_ACK", map[string]interface{}{
				"guildID": guildID,
				"newName": newName,
			})
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_RENAME_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleGuildRetypeReq implementa sp_NC_GUILD_RETYPE_REQ (0x74B4 -> p_Guild_SetType).
func (h *Handler) handleGuildRetypeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	slot, _ := pkt.Reader.ReadByte()
	newType, _ := pkt.Reader.ReadByte()

	var errCode uint16 = 0
	g, mbr, err := h.charRepo.GetGuildByCharID(player.Char.ID)
	if err != nil || g == nil || mbr == nil || mbr.Grade != 0 {
		errCode = 1
	} else {
		ret, rErr := h.charRepo.RetypeGuild(g.ID, newType)
		if rErr != nil || ret <= 0 {
			errCode = 1
		} else {
			if it := h.findItemAt(player.Char.ID, 9, int16(slot)); it != nil {
				_ = h.charRepo.DeleteItem(it.ItemKey)
				h.sendCellChangeFull(client, 9, int16(slot), nil, 0xFFFF)
			}
			cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_RETYPE_CMD)
			cmd.WriteUint32(uint32(g.ID))
			cmd.WriteByte(newType)
			if player.CurrentMap != nil {
				player.CurrentMap.Broadcast(cmd, 0)
			} else {
				_ = client.SendPacket(cmd)
			}
			h.DispatchWMS("wms_NC_GUILD_WORLD_RETYPE_ACK", map[string]interface{}{
				"guildID": g.ID,
				"newType": newType,
			})
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_RETYPE_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleUserConnectionChangeCharIDItemUseReq implementa sp_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ (0xB410 -> usp_Character_setChangeIDItemUsed).
func (h *Handler) handleUserConnectionChangeCharIDItemUseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slot, _ := pkt.Reader.ReadByte()
	var (
		errCode    uint16 = 1
		usedItemID uint16 = 0
	)
	if it := h.findItemAt(player.Char.ID, 9, int16(slot)); it != nil {
		usedItemID = it.ItemID
		remainLot := int16(0)
		if it.Flags > 1 {
			remainLot = int16(it.Flags - 1)
		}
		ret, err := h.charRepo.UseCharacterChangeIDItem(player.Char.ID, 0, player.Char.AccountID, it.ItemKey, 0, remainLot)
		if err == nil && ret == 0 {
			errCode = 0
			if remainLot == 0 {
				h.sendCellChangeFull(client, 9, int16(slot), nil, 0xFFFF)
			} else {
				it.Flags = int32(remainLot)
				h.sendCellChangeFull(client, 9, int16(slot), it, 0xFFFF)
			}
			h.DispatchWMS("wms_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD", map[string]interface{}{
				"charID": player.Char.ID,
				"userNo": player.Char.AccountID,
			})
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(usedItemID)
	_ = client.SendPacket(ack)
}

// handlePetAskNewNameAck implementa sp_NC_PET_ASK_NEW_NAME_ACK (0xD406).
func (h *Handler) handlePetAskNewNameAck(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || pkt.Reader.Remaining() < 4 {
		return
	}
	errVal, _ := pkt.Reader.ReadUint16()
	petHnd, _ := pkt.Reader.ReadUint16()
	if errVal == 0 && player.ActivePet == nil {
		player.ActivePet = &ActivePetState{
			Handle:     petHnd,
			PetID:      1,
			Mind:       100,
			Hunger:     100,
			SummonedAt: time.Now(),
		}
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeG(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(errVal), int64(petHnd), 0, 0, 0)
	}
}

func putUint32LE(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
