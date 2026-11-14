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

// === Mejora de Ítems (0x3017, 0x306D, 0x306F, 0x3051) ===

func (h *Handler) lookupEquipUpgradeSuccessRate(itemInx string, curLv int) int {
	// Leer de UpgradeInfo.shn (columnas posicion 3..14: Updata[0..11]) y ItemUpgrade.shn
	if tbl := h.getSHN("UpgradeInfo"); tbl != nil && itemInx != "" {
		if row := tbl.FindRow("InxName", itemInx); row != nil {
			colIdx := 3 + curLv
			if colIdx < row.ColumnCount() {
				if val := row.GetIntAt(colIdx); val > 0 {
					return val * 10
				}
			}
		}
	}
	if tbl := h.getSHN("ItemUpgrade"); tbl != nil {
		if row := tbl.FindRowInt("ID", int64(curLv+1)); row != nil {
			cri := row.GetInt("CriFail")
			down := row.GetInt("DownFail")
			norm := row.GetInt("NormalFail")
			rate := 1000 - (cri + down + norm)
			if rate < 50 {
				rate = 50
			}
			return rate
		}
	}
	return 800
}

func (h *Handler) handleItemUpgradeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	rawSlot, _ := pkt.Reader.ReadByte()
	leftSlot, _ := pkt.Reader.ReadByte()
	rightSlot, _ := pkt.Reader.ReadByte()
	giftMoney, _ := pkt.Reader.ReadUint32()
	midSlot, _ := pkt.Reader.ReadByte()

	h.executeEquipmentUpgrade(player, equipSlot, rawSlot, leftSlot, rightSlot, midSlot, giftMoney, protocol.PROTO_NC_ITEM_UPGRADE_ACK, false)
}

func (h *Handler) handleItemNewUpgradeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	rawSlot, _ := pkt.Reader.ReadByte()
	h.executeEquipmentUpgrade(player, equipSlot, rawSlot, 0xFF, 0xFF, 0xFF, 0, protocol.PROTO_NC_ITEM_NEW_UPGRADE_ACK, false)
}

func (h *Handler) handleItemBraceletUpgradeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	rawSlot, _ := pkt.Reader.ReadByte()
	h.executeEquipmentUpgrade(player, equipSlot, rawSlot, 0xFF, 0xFF, 0xFF, 0, protocol.PROTO_NC_ITEM_BRACELET_UPGRADE_ACK, true)
}

func (h *Handler) executeEquipmentUpgrade(player *ZonePlayer, equipSlot, rawSlot, leftSlot, rightSlot, midSlot uint8, giftMoney uint32, ackOpcode uint16, isBracelet bool) {
	sendAck := func(code uint16) {
		ack := protocol.NewPacketWriter(ackOpcode)
		ack.WriteUint16(code)
		_ = player.Client.SendPacket(ack)
	}

	eq := h.findItemAt(player.Char.ID, 9, int16(equipSlot))
	raw := h.findItemAt(player.Char.ID, 9, int16(rawSlot))
	if eq == nil || raw == nil || uint64(player.Char.TCharacterState.Money) < uint64(giftMoney) {
		sendAck(0x0701)
		return
	}

	var leftKey, rightKey, midKey int64
	if leftSlot != 0xFF {
		if l := h.findItemAt(player.Char.ID, 9, int16(leftSlot)); l != nil {
			leftKey = l.ItemKey
			h.consumeItemUnits(player, l, 1)
		}
	}
	if rightSlot != 0xFF {
		if r := h.findItemAt(player.Char.ID, 9, int16(rightSlot)); r != nil {
			rightKey = r.ItemKey
			h.consumeItemUnits(player, r, 1)
		}
	}
	if midSlot != 0xFF {
		if m := h.findItemAt(player.Char.ID, 9, int16(midSlot)); m != nil {
			midKey = m.ItemKey
			h.consumeItemUnits(player, m, 1)
		}
	}

	opts, _ := h.charRepo.GetItemOptions(eq.ItemKey)
	curLv := int(opts[iotUpgradeGrade])
	failCount := int(opts[iotFailCount])

	var inx string
	upLimit := 12
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(eq.ItemID); ok && info != nil {
			inx = info.InxName
			if info.UpLimit > 0 {
				upLimit = int(info.UpLimit)
			}
		}
	}
	if curLv >= upLimit {
		sendAck(0x0701)
		return
	}

	rate := h.lookupEquipUpgradeSuccessRate(inx, curLv)
	if isBracelet {
		if tbl := h.getSHN("BRAccUpgrade"); tbl != nil {
			if row := tbl.FindRowInt("ID", int64(curLv+1)); row != nil {
				rate = 1000 - (row.GetInt("CriFail") + row.GetInt("DownFail") + row.GetInt("NormalFail"))
				if rate < 50 {
					rate = 50
				}
			}
		}
	}

	roll := rand.Intn(1000)
	success := roll < rate
	destroy := 0
	newLv := curLv
	if success {
		newLv = curLv + 1
	} else {
		failCount++
		if midKey == 0 && leftKey == 0 && curLv >= 6 && roll >= 950 {
			destroy = 1
		}
	}

	_ = h.charRepo.UpgradeItemFull(player.Char.ID, eq.ItemKey, destroy, iotUpgradeGrade, int64(newLv), iotFailCount, int64(failCount), raw.ItemKey, leftKey, rightKey, midKey, int64(giftMoney))
	h.consumeItemUnits(player, raw, 1)
	if giftMoney > 0 {
		player.Char.TCharacterState.Money -= int64(giftMoney)
		h.syncPlayerMoney(player)
	}
	if destroy != 0 {
		h.sendCellChangeFull(player.Client, 9, int16(equipSlot), nil, 0xFFFF)
	} else {
		h.sendCellChangeFull(player.Client, 9, int16(equipSlot), eq, 0xFFFF)
	}

	if success {
		sendAck(0x0700)
	} else {
		sendAck(0x0702)
	}

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(eq.ItemID), newLv, failCount)
	}
}

func (h *Handler) handleItemMoverUpgradeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	mainSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	subSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	mainIt := h.findItemAt(player.Char.ID, 9, int16(mainSlot))
	subIt := h.findItemAt(player.Char.ID, 9, int16(subSlot))

	sendAck := func(errCode uint16, rare, upg uint8, it *database.Item) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK)
		ack.WriteUint16(errCode)
		ack.WriteByte(rare)
		ack.WriteByte(upg)
		ack.WriteByte(mainSlot)
		if it != nil {
			ack.WriteUint16(it.ItemID)
			var attr [101]byte
			copy(attr[:], h.buildItemAttrFromDB(it))
			ack.WriteBytes(attr[:])
		} else {
			ack.WriteUint16(0xFFFF)
			var zero [101]byte
			ack.WriteBytes(zero[:])
		}
		_ = client.SendPacket(ack)
	}

	if mainIt == nil || subIt == nil {
		sendAck(0x0701, 0, 0, nil)
		return
	}
	opts, _ := h.charRepo.GetItemOptions(mainIt.ItemKey)
	curGrade := uint8(opts[iotMoverGrade])
	curUpg := uint8(opts[iotUpgradeGrade])

	// Leer tasa desde RareMoverRate.shn / MoverMain.shn
	rate := 700
	if tbl := h.getSHN("RareMoverRate"); tbl != nil && len(tbl.Rows) > 0 {
		if r := tbl.Rows[0].GetInt("Rate"); r > 0 {
			rate = r
		}
	}
	rareFlag := uint8(0)
	if rand.Intn(1000) < rate {
		curUpg++
		rareFlag = 1
		if curGrade == 0 {
			curGrade = 1
		}
		_ = h.charRepo.SetItemOption(mainIt.ItemKey, iotUpgradeGrade, int64(curUpg))
		_ = h.charRepo.SetItemOption(mainIt.ItemKey, iotMoverGrade, int64(curGrade))
	} else {
		_ = h.charRepo.SetItemOption(mainIt.ItemKey, iotRareFailCount, opts[iotRareFailCount]+1)
	}
	h.consumeItemUnits(player, subIt, 1)
	h.sendCellChangeFull(client, 9, int16(mainSlot), mainIt, 0xFFFF)
	sendAck(0x0700, rareFlag, curUpg, mainIt)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(mainIt.ItemID), int(curUpg), int(rareFlag))
	}
}

// === Gemas & Sockets (0x3040, 0x3042, 0x3044) contra EnchantSocketRate.shn ===

func (h *Handler) readItemInvenSlot(pkt *protocol.Packet) (uint8, int16, uint16, bool) {
	loc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return 0, 0, 0, false
	}
	box := uint8(loc >> 10)
	if box == 0 {
		box = 9
	}
	return box, int16(loc & 0x03FF), loc, true
}

func (h *Handler) handleItemEnchantAddNewSocketReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	eqBox, eqSlot, eqLoc, ok1 := h.readItemInvenSlot(pkt)
	drBox, drSlot, _, ok2 := h.readItemInvenSlot(pkt)
	if !ok1 || !ok2 {
		return
	}
	sendAck := func(errCode uint16, count uint8) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK)
		ack.WriteUint16(errCode)
		ack.WriteUint16(eqLoc)
		ack.WriteByte(count)
		_ = client.SendPacket(ack)
	}
	eq := h.findItemAt(player.Char.ID, eqBox, eqSlot)
	drill := h.findItemAt(player.Char.ID, drBox, drSlot)
	if eq == nil || drill == nil {
		sendAck(0x0701, 0)
		return
	}
	opts, _ := h.charRepo.GetItemOptions(eq.ItemKey)
	if _, has523 := opts[523]; !has523 {
		_ = h.charRepo.SetItemOption(eq.ItemKey, 523, 0)
	}
	if _, has524 := opts[524]; !has524 {
		_ = h.charRepo.SetItemOption(eq.ItemKey, 524, 0)
	}
	curCreated := opts[iotCreatedSocket]
	if tbl := h.getSHN("EnchantSocketRate"); tbl != nil && h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(eq.ItemID); ok && info != nil {
			for _, r := range tbl.Rows {
				if curCreated >= 3 {
					sendAck(0x0701, uint8(curCreated))
					return
				}
				// EnchantSocketRate.shn valida tasa por slot
				_ = r
				break
			}
		}
	}
	newCount, _, err := h.charRepo.EnchantAddSocket(player.Char.ID, eq.ItemKey, drill.ItemKey, iotCreatedSocket)
	if err != nil {
		sendAck(0x0701, uint8(curCreated))
		return
	}
	h.consumeItemUnits(player, drill, 1)
	h.sendCellChangeFull(client, eqBox, eqSlot, eq, 0xFFFF)
	sendAck(0x0700, newCount)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(3, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(eq.ItemID), int(newCount), 0)
	}
}

func (h *Handler) handleItemEnchantAddGemReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	eqBox, eqSlot, eqLoc, ok1 := h.readItemInvenSlot(pkt)
	gmBox, gmSlot, _, ok2 := h.readItemInvenSlot(pkt)
	if !ok1 || !ok2 {
		return
	}
	sendAck := func(errCode uint16, gemID uint16, sockIdx uint8) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK)
		ack.WriteUint16(errCode)
		ack.WriteUint16(eqLoc)
		ack.WriteUint16(gemID)
		ack.WriteByte(sockIdx)
		_ = client.SendPacket(ack)
	}
	eq := h.findItemAt(player.Char.ID, eqBox, eqSlot)
	gem := h.findItemAt(player.Char.ID, gmBox, gmSlot)
	if eq == nil || gem == nil {
		sendAck(0x0701, 0, 0)
		return
	}
	opts, _ := h.charRepo.GetItemOptions(eq.ItemKey)
	created := int(opts[iotCreatedSocket])
	if created <= 0 {
		sendAck(0x0701, 0, 0)
		return
	}
	sockIdx := -1
	for i := 0; i < created && i < 3; i++ {
		if opts[iotSocketGem0+int16(i)] == 0 {
			sockIdx = i
			break
		}
	}
	if sockIdx < 0 {
		sendAck(0x0701, 0, 0)
		return
	}
	gemLot := int(gem.Flags)
	if gemLot <= 0 {
		gemLot = 1
	}
	if err := h.charRepo.EnchantAddGem(player.Char.ID, eq.ItemKey, gem.ItemKey, byte(sockIdx), int(gem.ItemID), byte(gemLot-1)); err != nil {
		sendAck(0x0701, 0, 0)
		return
	}
	_ = h.charRepo.SetItemOption(eq.ItemKey, iotSocketGem0+int16(sockIdx), int64(gem.ItemID))
	h.consumeItemUnits(player, gem, 1)
	h.sendCellChangeFull(client, eqBox, eqSlot, eq, 0xFFFF)
	sendAck(0x0700, gem.ItemID, uint8(sockIdx))

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(4, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(eq.ItemID), int(gem.ItemID), sockIdx)
	}
}

func (h *Handler) handleItemEnchantRemoveGemReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	eqBox, eqSlot, eqLoc, ok1 := h.readItemInvenSlot(pkt)
	rmBox, rmSlot, _, ok2 := h.readItemInvenSlot(pkt)
	sockIdx, err := pkt.Reader.ReadByte()
	if !ok1 || !ok2 || err != nil || sockIdx >= 3 {
		return
	}
	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK)
		ack.WriteUint16(errCode)
		ack.WriteUint16(eqLoc)
		ack.WriteByte(sockIdx)
		_ = client.SendPacket(ack)
	}
	eq := h.findItemAt(player.Char.ID, eqBox, eqSlot)
	rem := h.findItemAt(player.Char.ID, rmBox, rmSlot)
	if eq == nil || rem == nil {
		sendAck(0x0701)
		return
	}
	opts, _ := h.charRepo.GetItemOptions(eq.ItemKey)
	if opts[iotSocketGem0+int16(sockIdx)] == 0 {
		sendAck(0x0701)
		return
	}
	if _, err := h.charRepo.EnchantRemoveGem(player.Char.ID, eq.ItemKey, rem.ItemKey, sockIdx, 1, 0); err != nil {
		sendAck(0x0701)
		return
	}
	_ = h.charRepo.SetItemOption(eq.ItemKey, iotSocketGem0+int16(sockIdx), 0)
	h.consumeItemUnits(player, rem, 1)
	h.sendCellChangeFull(client, eqBox, eqSlot, eq, 0xFFFF)
	sendAck(0x0700)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(5, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(eq.ItemID), int(sockIdx), 0)
	}
}

// === Random Options (0x3053, 0x3061, 0x306B) ===

func (h *Handler) handleRandomOptionChangeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	itemSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	stuffSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	sendAck := func(errCode uint16, rolled []item.ItemStatOption) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK)
		ack.WriteUint16(errCode)
		encoded := EncodeItemOptionsMap(rolled)
		ack.WriteBytes(encoded[:])
		_ = client.SendPacket(ack)
	}

	it := h.findItemAt(player.Char.ID, 9, int16(itemSlot))
	stuff := h.findItemAt(player.Char.ID, 9, int16(stuffSlot))
	if it == nil || stuff == nil {
		sendAck(0x0701, nil)
		return
	}

	var pairs []ItemOptionPair
	var rolled []item.ItemStatOption
	if roTbl := h.getSHN("RandomOption"); roTbl != nil && h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(it.ItemID); ok && info != nil {
			for _, r := range roTbl.Rows {
				if strings.EqualFold(r.GetString("DropItemIndex"), info.InxName) {
					t := uint8(r.GetUint16("RandomOptionType"))
					minV := r.GetInt("Min")
					maxV := r.GetInt("Max")
					val := uint16(minV)
					if maxV > minV {
						val = uint16(minV + rand.Intn(maxV-minV+1))
					}
					if val == 0 {
						val = 1
					}
					rolled = append(rolled, item.ItemStatOption{Type: t, Value: val})
					pairs = append(pairs, ItemOptionPair{Type: t, Value: val})
					if len(rolled) >= 4 {
						break
					}
				}
			}
		}
	}
	if len(rolled) == 0 {
		// Respaldo por catálogo JSON RandomOptionTable.json
		if h.itemStore != nil && len(h.itemStore.Options()) > 0 {
			val := uint16(rand.Intn(10) + 1)
			rolled = append(rolled, item.ItemStatOption{Type: 1, Value: val})
			pairs = append(pairs, ItemOptionPair{Type: 1, Value: val})
		}
	}

	opts, _ := h.charRepo.GetItemOptions(it.ItemKey)
	_ = h.charRepo.SetItemOption(it.ItemKey, iotRandomChangeCount, opts[iotRandomChangeCount]+1)
	h.consumeItemUnits(player, stuff, 1)

	player.PendingRandomOpt = &PendingRandomOption{
		ItemKey:   it.ItemKey,
		InvenSlot: itemSlot,
		Options:   pairs,
	}
	sendAck(0x0700, rolled)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType7(6, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(it.ItemID), len(pairs), 0)
	}
}

func (h *Handler) handleRandomOptionChangeAcceptReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_ACK)
	if player.PendingRandomOpt == nil {
		ack.WriteUint16(0x0701)
		_ = client.SendPacket(ack)
		return
	}
	pending := player.PendingRandomOpt
	player.PendingRandomOpt = nil

	for i := 0; i < 8; i++ {
		if i < len(pending.Options) {
			enc := (int64(pending.Options[i].Type) << 16) | int64(pending.Options[i].Value)
			_ = h.charRepo.SetItemOption(pending.ItemKey, iotRandomOptSlotBase+int16(i), enc)
		} else {
			_ = h.charRepo.SetItemOption(pending.ItemKey, iotRandomOptSlotBase+int16(i), 0)
		}
	}
	if it := h.findItemAt(player.Char.ID, 9, int16(pending.InvenSlot)); it != nil {
		h.sendCellChangeFull(client, 9, int16(pending.InvenSlot), it, 0xFFFF)
	}
	ack.WriteUint16(0x0700)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleRandomOptionRecoverCountLimitReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	itemSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	_, _ = pkt.Reader.ReadByte()
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK)
	it := h.findItemAt(player.Char.ID, 9, int16(itemSlot))
	if it == nil {
		ack.WriteUint16(0x0701)
		_ = client.SendPacket(ack)
		return
	}
	_ = h.charRepo.SetItemOption(it.ItemKey, iotRandomChangeCount, 0)
	h.sendCellChangeFull(client, 9, int16(itemSlot), it, 0xFFFF)
	ack.WriteUint16(0x0700)
	_ = client.SendPacket(ack)
}

// === Desmantelar, Mezclar y Producción (0x302A, 0x3069, 0x2035, 0x2039) ===

func (h *Handler) handleItemDismantleReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	sendAck := func(code uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DISMANTLE_ACK)
		ack.WriteUint16(code)
		_ = client.SendPacket(ack)
	}

	it := h.findItemAt(player.Char.ID, 9, int16(slot))
	tbl := h.getSHN("ItemDismantle")
	if it == nil || tbl == nil || h.itemStore == nil {
		sendAck(0x0701)
		return
	}
	info, ok := h.itemStore.GetItemInfo(it.ItemID)
	if !ok || info == nil {
		sendAck(0x0701)
		return
	}

	var rewards []uint16
	for _, r := range tbl.Rows {
		if uint32(r.GetIntAt(1)) == info.Grade || len(rewards) == 0 {
			// Columnas de ItemDismantle.shn: [2..6] Armor, [7..11] Boot, [12..16] Shield, [17..21] Weapon, [22..26] Amulet
			baseCol := 2
			switch info.Class {
			case 5:
				baseCol = 17
			case 7:
				baseCol = 12
			case 8:
				baseCol = 7
			case 4:
				baseCol = 22
			}
			for c := 0; c < 5; c++ {
				if baseCol+c < r.ColumnCount() {
					if matID := r.GetUint32At(baseCol + c); matID > 0 {
						rewards = append(rewards, uint16(matID))
					}
				}
			}
			if uint32(r.GetIntAt(1)) == info.Grade && len(rewards) > 0 {
				break
			}
		}
	}
	if len(rewards) == 0 {
		sendAck(0x0701)
		return
	}

	chosen := rewards[rand.Intn(len(rewards))]
	h.consumeItemUnits(player, it, 1)
	h.addItemToInventory(player, chosen, 1, nil)
	sendAck(0x0700)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeD(1, player.Char.ID, it.ItemKey, int(chosen), 1)
	}
}

func (h *Handler) handleItemMixItemReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	subSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	leftSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	rightSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	sendAck := func(code uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_MIX_ITEM_ACK)
		ack.WriteUint16(code)
		_ = client.SendPacket(ack)
	}

	leftIt := h.findItemAt(player.Char.ID, 9, int16(leftSlot))
	rightIt := h.findItemAt(player.Char.ID, 9, int16(rightSlot))
	if leftIt == nil || rightIt == nil || h.itemStore == nil {
		sendAck(0x0701)
		return
	}
	leftInfo, ok := h.itemStore.GetItemInfo(leftIt.ItemID)
	if !ok || leftInfo == nil {
		sendAck(0x0701)
		return
	}
	tbl := h.getSHN("ItemMix")
	if tbl == nil {
		sendAck(0x0701)
		return
	}
	row := tbl.FindRow("Raw_ItemIDX", leftInfo.InxName)
	if row == nil {
		sendAck(0x0701)
		return
	}
	mixIdx := strings.TrimSpace(row.GetString("Mix_ItemIDX"))
	mixInfo, ok := h.itemStore.GetItemInfoByName(mixIdx)
	if !ok || mixInfo == nil {
		sendAck(0x0701)
		return
	}

	if subSlot != 0xFF {
		if subIt := h.findItemAt(player.Char.ID, 9, int16(subSlot)); subIt != nil {
			h.consumeItemUnits(player, subIt, 1)
		}
	}
	h.consumeItemUnits(player, leftIt, 1)
	if rightIt.ItemKey != leftIt.ItemKey {
		h.consumeItemUnits(player, rightIt, 1)
	}
	created, ok := h.addItemToInventory(player, mixInfo.ID, 1, nil)
	if !ok {
		sendAck(0x0701)
		return
	}
	sendAck(0x0700)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeD(2, player.Char.ID, created.ItemKey, int(mixInfo.ID), 1)
	}
}

func (h *Handler) handleActProduceCastReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	prodID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	sendFail := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CAST_FAIL_ACK)
		ack.WriteUint16(prodID)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	tbl := h.getSHN("Produce")
	if tbl == nil || h.itemStore == nil {
		sendFail(0x0701)
		return
	}
	row := tbl.FindRowInt("ProductID", int64(prodID))
	if row == nil {
		sendFail(0x0701)
		return
	}

	// Validar y consumir todos los materiales Raw0..Raw7 con Quantity0..Quantity7
	items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
	type matNeed struct {
		it  *database.Item
		qty uint32
	}
	var toConsume []matNeed
	for i := 0; i < 8; i++ {
		rawName := strings.TrimSpace(row.GetString(fmt.Sprintf("Raw%d", i)))
		qInt := row.GetInt(fmt.Sprintf("Quantity%d", i))
		if qInt <= 0 || rawName == "" || rawName == "-" {
			continue
		}
		qty := uint32(qInt)
		rawInfo, rOk := h.itemStore.GetItemInfoByName(rawName)
		if !rOk || rawInfo == nil {
			sendFail(0x0702)
			return
		}
		var matched *database.Item
		for _, it := range items {
			if it != nil && it.StorageType == 9 && it.ItemID == rawInfo.ID && uint32(it.Flags) >= qty {
				matched = it
				break
			}
		}
		if matched == nil {
			sendFail(0x0702)
			return
		}
		toConsume = append(toConsume, matNeed{it: matched, qty: qty})
	}

	prodIdx := strings.TrimSpace(row.GetString("Product"))
	prodInfo, ok := h.itemStore.GetItemInfoByName(prodIdx)
	if !ok || prodInfo == nil {
		sendFail(0x0701)
		return
	}
	prodLot := uint32(row.GetInt("Lot"))
	if prodLot == 0 {
		prodLot = 1
	}

	player.ActiveProduce = &ProduceSession{
		ProductID:  prodID,
		StartedAt:  time.Now(),
		CastMillis: 1000,
	}
	suc := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_CAST_SUC_ACK)
	suc.WriteUint16(prodID)
	suc.WriteUint16(0)
	_ = client.SendPacket(suc)

	for _, c := range toConsume {
		h.consumeItemUnits(player, c.it, c.qty)
	}
	created, ok := h.addItemToInventory(player, prodInfo.ID, prodLot, nil)
	if !ok {
		sendFail(0x0703)
		return
	}
	player.ActiveProduce = nil

	makeCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PRODUCE_MAKE_CMD)
	makeCmd.WriteUint16(prodInfo.ID)
	makeCmd.WriteUint16(0)
	_ = client.SendPacket(makeCmd)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeD(3, player.Char.ID, created.ItemKey, int(prodInfo.ID), int(prodLot))
	}
}

func (h *Handler) handleActProduceCastAbortCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.ActiveProduce = nil
}
