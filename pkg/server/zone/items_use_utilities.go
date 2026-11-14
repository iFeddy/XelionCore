package zone

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// === Dept 2: iac_FoodEat (`0x0804`, `0x0805`) & Heartbeat compatibilidad ===

func (h *Handler) handleFoodEatCmd(client *network.Client, pkt *protocol.Packet) {
	// En Dept 2, 0x0804 y 0x0805 también transportan HEARTBEAT cuando no hay payload.
	if pkt == nil || pkt.Reader == nil || pkt.Reader.Remaining() == 0 {
		if pkt != nil && pkt.OpCode == protocol.PROTO_NC_MISC_HEARTBEAT_REQ {
			ack := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HEARTBEAT_ACK)
			_ = client.SendPacket(ack)
		}
		return
	}
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	feedItem := h.findItemAt(player.Char.ID, 9, int16(slot))
	if feedItem == nil {
		return
	}
	// Buscar montura o mascota activa para aplicar iac_FoodEat
	var iac ItemAttributeClass
	iac.iac_SetPetID(feedItem.ItemID)
	restored := iac.iac_FoodEat(200)
	_ = h.charRepo.SetItemOption(feedItem.ItemKey, iotPetHungry, int64(restored))
}

// === Dept 12: UseAbort (`0x3019`), ActionByItem (`0x2067`), Minimon (`0x3050`, `0x1124..0x112A`) ===

func (h *Handler) handleItemUseAbortCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.ActiveProduce = nil
}

func (h *Handler) handleActionByItemReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	it := h.findItemAt(player.Char.ID, 9, int16(slot))
	sendAck := func(errCode uint16, emoticon uint8) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_ACTIONBYITEM_ACK)
		ack.WriteUint16(errCode)
		ack.WriteByte(emoticon)
		_ = client.SendPacket(ack)
	}
	if it == nil || h.itemStore == nil {
		sendAck(0x0701, 0)
		return
	}
	info, ok := h.itemStore.GetItemInfo(it.ItemID)
	if !ok || info == nil {
		sendAck(0x0701, 0)
		return
	}
	emoticonID := uint8(1)
	if tbl := h.getSHN("ActionEffectItem"); tbl != nil {
		if row := tbl.FindRow("InxName", info.InxName); row != nil {
			if e := row.GetUint16("EmoticonID"); e > 0 {
				emoticonID = uint8(e)
			}
		}
	}
	sendAck(0x0700, emoticonID)
}

func (h *Handler) handleMinimonEquipReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slotNo, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	it := h.findItemAt(player.Char.ID, 9, int16(slotNo))
	if it == nil {
		return
	}
	if tbl := h.getSHN("MinimonInfo"); tbl != nil && h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(it.ItemID); ok && info != nil {
			if tbl.FindRow("ItemIDX", info.InxName) == nil {
				return
			}
		}
	}
	// Colocar en slot de equipamiento (Box 8) e informar
	prevEquip := h.findItemAt(player.Char.ID, 8, int16(equipSlot))
	if prevEquip != nil {
		_ = h.charRepo.ExchangeItemStorage(it.ItemKey, 8, int16(equipSlot), prevEquip.ItemKey, 9, int16(slotNo))
		prevEquip.StorageType = 9
		prevEquip.Storage = int16(slotNo)
		h.sendCellChangeFull(client, 9, int16(slotNo), prevEquip, 0xFFFF)
	} else {
		_ = h.charRepo.SetItemStorage(it.ItemKey, 8, int16(equipSlot))
	}
	it.StorageType = 8
	it.Storage = int16(equipSlot)
	h.sendCellChangeFull(client, 8, int16(equipSlot), it, 0xFFFF)
	h.syncPlayerEquipsAndStats(player)
}

func (h *Handler) readMinimonSlots(pkt *protocol.Packet) [12]uint16 {
	var slots [12]uint16
	for i := 0; i < 12; i++ {
		v, err := pkt.Reader.ReadUint16()
		if err != nil {
			break
		}
		slots[i] = v
	}
	return slots
}

func (h *Handler) handleMinimonNormalOnReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slots := h.readMinimonSlots(pkt)
	player.MinimonNormalSlots = slots
	for i, itemID := range slots {
		if itemID != 0 && itemID != 0xFFFF {
			_ = h.charRepo.SetMinimonAutoItem(player.Char.ID, 0, i, int(itemID))
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleMinimonNormalOffReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	player.MinimonNormalSlots = [12]uint16{}
	for i := 0; i < 12; i++ {
		_ = h.charRepo.SetMinimonAutoItem(player.Char.ID, 0, i, 0)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleMinimonChargedOnReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slots := h.readMinimonSlots(pkt)
	player.MinimonChargedSlots = slots
	for i, itemID := range slots {
		if itemID != 0 && itemID != 0xFFFF {
			_ = h.charRepo.SetMinimonAutoItem(player.Char.ID, 1, i, int(itemID))
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleMinimonChargedOffReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	player.MinimonChargedSlots = [12]uint16{}
	for i := 0; i < 12; i++ {
		_ = h.charRepo.SetMinimonAutoItem(player.Char.ID, 1, i, 0)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// === Dept 12: Dice Game, Durabilidad, Desligar, Chat Color, TermExtend, ClassChange, MapLink, Cofres ===

func (h *Handler) handleDiceGameStartReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	itemHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	isYes, _ := pkt.Reader.ReadByte()

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_START_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	roll := int32(rand.Intn(100) + 1)
	if isYes == 0 {
		roll = 0
	}
	if player.CurrentMap != nil {
		startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_START_CMD)
		startCmd.WriteUint16(player.MapObjectID)
		startCmd.WriteInt32(roll)
		player.CurrentMap.Broadcast(startCmd, 0)

		resCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_RESULT_CMD)
		resCmd.WriteUint16(0)
		var nameBytes [20]byte
		copy(nameBytes[:], player.Char.Name)
		resCmd.WriteBytes(nameBytes[:])
		resCmd.WriteInt32(roll)
		resCmd.WriteUint16(itemHandle)
		player.CurrentMap.Broadcast(resCmd, 0)
	}
}

func (h *Handler) handleMHFurnitureEndureKitUseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	furnHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	kit := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if kit == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_ACK)
		ack.WriteUint32(0)
		ack.WriteUint16(0x0701)
		_ = client.SendPacket(ack)
		return
	}
	days := uint32(30)
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(kit.ItemID); ok && info != nil {
			if eff, ok := h.itemStore.GetItemUseEffect(info.InxName); ok && eff != nil && eff.UseValueA > 0 {
				days = uint32(eff.UseValueA)
			}
		}
	}
	newEndure := encodeShineEndTimeDays(uint32(time.Now().Unix()), days)
	_ = h.charRepo.FurnitureEndureUpdate(player.Char.ID, kit.ItemKey, 0, int16(kit.Flags), int64(furnHandle), iotEndureEndTime, int64(newEndure))
	h.consumeItemUnits(player, kit, 1)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_ACK)
	ack.WriteUint32(newEndure)
	ack.WriteUint16(0x0700)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleEndureChargeCommon(player *ZonePlayer, pkt *protocol.Packet, isShield bool) {
	chargingSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	targetLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	box := uint8(targetLoc >> 10)
	if box == 0 {
		box = 9
	}
	slot := int16(targetLoc & 0x03FF)

	sendAck := func(charged uint32, errCode uint16) {
		op := protocol.PROTO_NC_ITEM_WEAPONENDURE_CHARGE_ACK
		if isShield {
			op = protocol.PROTO_NC_ITEM_SHIELDENDURE_CHARGE_ACK
		}
		ack := protocol.NewPacketWriter(op)
		ack.WriteUint32(charged)
		ack.WriteUint16(targetLoc)
		ack.WriteUint16(errCode)
		_ = player.Client.SendPacket(ack)
	}

	kit := h.findItemAt(player.Char.ID, 9, int16(chargingSlot))
	target := h.findItemAt(player.Char.ID, box, slot)
	if kit == nil || target == nil {
		sendAck(0, 0x0701)
		return
	}
	addEndure := uint32(1000)
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(kit.ItemID); ok && info != nil {
			if eff, ok := h.itemStore.GetItemUseEffect(info.InxName); ok && eff != nil && eff.UseValueA > 0 {
				addEndure = uint32(eff.UseValueA)
			}
		}
	}
	opts, _ := h.charRepo.GetItemOptions(target.ItemKey)
	cur := uint32(opts[iotEndureEndTime]) + addEndure
	_ = h.charRepo.CostumeWeaponEndureUpdate(player.Char.ID, kit.ItemKey, 0, int16(kit.Flags), target.ItemKey, iotEndureEndTime, int64(cur))
	h.consumeItemUnits(player, kit, 1)
	h.sendCellChangeFull(player.Client, box, slot, target, 0xFFFF)

	sendAck(cur, 0x0700)
	setOp := protocol.PROTO_NC_ITEM_WEAPONENDURESET_CMD
	if isShield {
		setOp = protocol.PROTO_NC_ITEM_SHIELDENDURESET_CMD
	}
	setCmd := protocol.NewPacketWriter(setOp)
	setCmd.WriteUint32(cur)
	_ = player.Client.SendPacket(setCmd)
}

func (h *Handler) handleWeaponEndureChargeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	h.handleEndureChargeCommon(player, pkt, false)
}

func (h *Handler) handleShieldEndureChargeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	h.handleEndureChargeCommon(player, pkt, true)
}

func (h *Handler) handleEquipBelongedCancelUseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cancelSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	belongSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	cancelIt := h.findItemAt(player.Char.ID, 9, int16(cancelSlot))
	targetIt := h.findItemAt(player.Char.ID, 9, int16(belongSlot))
	if cancelIt == nil || targetIt == nil {
		sendAck(0x0701)
		return
	}
	if err := h.charRepo.PutOnBelongedCancel(player.Char.ID, cancelIt.ItemKey, 0, 0, targetIt.ItemKey, iotIsPutOnBelonged, 0); err != nil {
		sendAck(0x0701)
		return
	}
	h.consumeItemUnits(player, cancelIt, 1)
	h.sendCellChangeFull(client, 9, int16(belongSlot), targetIt, 0xFFFF)
	sendAck(0x0700)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType6(10, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(targetIt.ItemID), targetIt.ItemKey)
	}
}

func (h *Handler) handleChatColorChangeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	it := h.findItemAt(player.Char.ID, 9, int16(slot))
	sendAck := func(errCode uint16, colorID uint8) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHAT_COLOR_CHANGE_ACK)
		ack.WriteUint16(errCode)
		ack.WriteByte(colorID)
		_ = client.SendPacket(ack)
	}
	if it == nil {
		sendAck(0x0701, 0)
		return
	}
	colorID := uint8(1)
	if tbl := h.getSHN("ChatColor"); tbl != nil && len(tbl.Rows) > 0 {
		colorID = uint8(tbl.Rows[0].GetUint16("ChatColorID"))
		if h.itemStore != nil {
			if info, ok := h.itemStore.GetItemInfo(it.ItemID); ok && info != nil {
				for _, r := range tbl.Rows {
					if strings.Contains(strings.ToLower(info.InxName), strings.ToLower(r.GetString("InxName"))) {
						colorID = uint8(r.GetUint16("ChatColorID"))
						break
					}
				}
			}
		}
	}
	h.consumeItemUnits(player, it, 1)
	sendAck(0x0700, colorID)
}

func (h *Handler) handleTermExtendReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	extSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	termSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_TERMEXTEND_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}
	extIt := h.findItemAt(player.Char.ID, 9, int16(extSlot))
	termIt := h.findItemAt(player.Char.ID, 9, int16(termSlot))
	if extIt == nil || termIt == nil || h.itemStore == nil {
		sendAck(0x0701)
		return
	}
	extInfo, _ := h.itemStore.GetItemInfo(extIt.ItemID)
	termInfo, _ := h.itemStore.GetItemInfo(termIt.ItemID)
	addDays := uint32(7)
	if tbl := h.getSHN("TermExtendMatch"); tbl != nil && extInfo != nil && termInfo != nil {
		matched := false
		for _, r := range tbl.Rows {
			if strings.EqualFold(r.GetString("ExtendItemIDX"), extInfo.InxName) &&
				strings.EqualFold(r.GetString("TermItemIDX"), termInfo.InxName) {
				matched = true
				if d := r.GetUint32("ExtendDay"); d > 0 {
					addDays = d
				}
				break
			}
		}
		if !matched && len(tbl.Rows) > 0 {
			// Permitir match flexible si la tabla tiene sufijo numérico
			for _, r := range tbl.Rows {
				if strings.EqualFold(r.GetString("ExtendItemIDX"), extInfo.InxName) {
					matched = true
					break
				}
			}
			if !matched {
				sendAck(0x0701)
				return
			}
		}
	}

	opts, _ := h.charRepo.GetItemOptions(termIt.ItemKey)
	newEnd := encodeShineEndTimeDays(uint32(opts[iotUseableDateTime]), addDays)
	_ = h.charRepo.SetItemOption(termIt.ItemKey, iotUseableDateTime, int64(newEnd))
	h.consumeItemUnits(player, extIt, 1)
	h.sendCellChangeFull(client, 9, int16(termSlot), termIt, 0xFFFF)
	sendAck(0x0700)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType6(11, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(termIt.ItemID), int64(newEnd))
	}
}

func (h *Handler) handleClassChangeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	stuffSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CLASS_CHANGE_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	stuff := h.findItemAt(player.Char.ID, 9, int16(stuffSlot))
	if stuff == nil || player.Char.Job >= 25 {
		sendAck(0x0701)
		return
	}
	newClass := player.Char.Job + 1
	player.Char.Job = newClass
	player.Char.Sync()
	_ = h.charRepo.UpdateCharacter(player.Char)
	h.consumeItemUnits(player, stuff, 1)

	openCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_OPENCLASSCHANGEMENU_CMD)
	openCmd.WriteByte(newClass)
	openCmd.WriteUint64(0)
	openCmd.WriteUint16(0x0700)
	_ = client.SendPacket(openCmd)

	sendAck(0x0700)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType6(12, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(newClass), stuff.ItemKey)
	}
}

func (h *Handler) handleMapLinkScrollReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	targetMapID, err := pkt.Reader.ReadUint32()
	if err != nil {
		return
	}

	sendAck := func(errCode uint16, lvFrom, lvTo uint8) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MAPLINK_SCROLL_ACK)
		ack.WriteUint16(errCode)
		ack.WriteByte(lvFrom)
		ack.WriteByte(lvTo)
		_ = client.SendPacket(ack)
	}

	// Verificar si el mapa destino existe en FieldList
	if h.fieldStore == nil {
		sendAck(0x0701, 0, 0)
		return
	}
	fld, ok := h.fieldStore.GetFieldBySerial(int(targetMapID))
	if !ok || fld == nil {
		sendAck(0x0701, 0, 0)
		return
	}
	lvFrom := uint8(0)
	lvTo := uint8(150)
	if tbl := h.getSHN("MapInfo"); tbl != nil {
		if row := tbl.FindRow("MapName", fld.MapID); row != nil {
			lvFrom = uint8(row.GetInt("LvFrom"))
			if lt := row.GetInt("LvTo"); lt > 0 {
				lvTo = uint8(lt)
			}
		}
	}
	if player.Char.Level < lvFrom {
		sendAck(0x0702, lvFrom, lvTo)
		return
	}
	rx, ry := 0, 0
	if len(fld.Respawn.Spots) > 0 {
		rx = fld.Respawn.Spots[0].X
		ry = fld.Respawn.Spots[0].Y
	}
	sendAck(0x0700, lvFrom, lvTo)
	h.triggerGateWarp(player, fmt.Sprintf("Map:%s,%d,%d", fld.MapID, rx, ry))
}

// openTreasureChestItem abre un cofre o MysteryVault según TreasureReward.json o MysteryVaultServer.shn.
func (h *Handler) openTreasureChestItem(player *ZonePlayer, chestItem *database.Item) ([]*database.Item, bool) {
	if player == nil || chestItem == nil || h.itemStore == nil {
		return nil, false
	}
	info, ok := h.itemStore.GetItemInfo(chestItem.ItemID)
	if !ok || info == nil {
		return nil, false
	}

	var granted []*database.Item
	if container, ok := h.itemStore.GetContainer(info.InxName); ok && container != nil && len(container.Rewards) > 0 {
		totalRate := 0
		for _, r := range container.Rewards {
			if r.Rate > 0 {
				totalRate += r.Rate
			}
		}
		var picked *item.TreasureRewardItem
		if totalRate > 0 {
			roll := rand.Intn(totalRate)
			acc := 0
			for i := range container.Rewards {
				r := &container.Rewards[i]
				if r.Rate <= 0 {
					continue
				}
				acc += r.Rate
				if roll < acc {
					picked = r
					break
				}
			}
		} else {
			picked = &container.Rewards[0]
		}
		if picked != nil {
			if rwInfo, ok := h.itemStore.GetItemInfoByName(picked.ItemID); ok && rwInfo != nil {
				lot := uint32(picked.Lot)
				if lot == 0 {
					lot = 1
				}
				if created, added := h.addItemToInventory(player, rwInfo.ID, lot, nil); added {
					granted = append(granted, created)
				}
			}
		}
	} else if mvTbl := h.getSHN("MysteryVaultServer"); mvTbl != nil {
		for _, r := range mvTbl.Rows {
			if strings.EqualFold(r.GetString("MV_GroupItemInx"), info.InxName) {
				rwIdx := strings.TrimSpace(r.GetString("ItemInx"))
				if rwInfo, ok := h.itemStore.GetItemInfoByName(rwIdx); ok && rwInfo != nil {
					lot := r.GetUint32("Lot")
					if lot == 0 {
						lot = 1
					}
					if created, added := h.addItemToInventory(player, rwInfo.ID, lot, nil); added {
						granted = append(granted, created)
						break
					}
				}
			}
		}
	}

	if len(granted) == 0 {
		return nil, false
	}
	h.consumeItemUnits(player, chestItem, 1)
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType6(13, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(granted[0].ItemID), granted[0].ItemKey)
	}
	return granted, true
}
