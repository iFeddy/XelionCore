package zone

import (
	"strings"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Códigos de respuesta oficiales de Zone.exe para NC_ITEM Buy/Sell/Repurchase/Storage
const (
	errItemBuySuccess uint16 = 0x0201
	errItemBuyFull    uint16 = 0x0203
	errItemBuyState   uint16 = 0x0204
	errItemBuyNoCen   uint16 = 0x020E
	errItemBuyFail    uint16 = 0x0211

	errItemSellSuccess uint16 = 0x0381
	errItemSellState   uint16 = 0x0382
	errItemSellMissing uint16 = 0x0383
	errItemSellNoSell  uint16 = 0x0384

	errItemRepurchaseSuccess uint16 = 0x3668
	errItemRepurchaseFail    uint16 = 0x366A
	errItemRepurchaseNoCen   uint16 = 0x366B
	errItemRepurchaseFull    uint16 = 0x3671

	errItemStorageSuccess uint16 = 0x0F80
	errItemStorageState   uint16 = 0x0F81
	errItemStorageFail    uint16 = 0x0F82
	errItemStorageNoCen   uint16 = 0x0F83

	errGuildStorageSuccess uint16 = 0x1480
	errGuildStorageState   uint16 = 0x1481
	errGuildStorageNoGuild uint16 = 0x1482
	errGuildStorageNoCen   uint16 = 0x1483
	errGuildStorageDenied  uint16 = 0x1484

	errAcademyWithdrawSuccess uint16 = 0x17E0
	errAcademyWithdrawFail    uint16 = 0x17E1
	errAcademyWithdrawState   uint16 = 0x17E2
	errAcademyWithdrawNoGuild uint16 = 0x17E3
	errAcademyWithdrawNoCen   uint16 = 0x17E5

	errAcademyDepositSuccess uint16 = 0x19A0
	errAcademyDepositNoCen   uint16 = 0x19A1
	errAcademyDepositState   uint16 = 0x19A2
	errAcademyDepositNoGuild uint16 = 0x19A3
	errAcademyDepositFail    uint16 = 0x19A4

	errAcademyStorageSuccess uint16 = 0x17D8
)

func (h *Handler) findNPCMerchantName(mobID uint16) string {
	return ""
}

// handleItemBuyReq implementa sp_NC_ITEM_BUY_REQ (0x3003) contra NPCItemList/*.json, ItemShop.shn e ItemMoney.shn.
func (h *Handler) handleItemBuyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	itemID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	lot, err := pkt.Reader.ReadUint32()
	if err != nil {
		return
	}

	sendBuyAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendBuyAck(errItemBuyState)
		return
	}
	if lot == 0 || h.itemStore == nil {
		sendBuyAck(errItemBuyFail)
		return
	}
	info, ok := h.itemStore.GetItemInfo(itemID)
	if !ok || info == nil {
		sendBuyAck(errItemBuyFail)
		return
	}

	// Validar catálogo del NPC mercader actual (si tiene tienda abierta y RoleArg en NPCItemList)
	if player.OpenNPCMobID != 0 && player.CurrentMap != nil && h.dataStore != nil && h.dataStore.NPCItem != nil {
		var npcRoleArg string
		for _, n := range player.CurrentMap.GetAllNPCs() {
			if n != nil && n.MobID == player.OpenNPCMobID {
				npcRoleArg = n.RoleArg
				if npcRoleArg == "" {
					npcRoleArg = n.Name
				}
				break
			}
		}
		if npcRoleArg != "" {
			if merchant, exists := h.dataStore.NPCItem.GetMerchant(npcRoleArg); exists && merchant != nil {
				if !merchant.HasItem(info.InxName) {
					sendBuyAck(errItemBuyState)
					return
				}
			}
		}
	}

	// Revisar paquete/lote en ItemShop.shn (si lot == 1 y ItemShop define paquete > 1)
	actualLot := lot
	bundleLot := h.getShopBundleLot(info.InxName)
	if bundleLot > 1 && actualLot == 1 {
		actualLot = bundleLot
	}

	// Verificar ItemMoney.shn (mercader de fichas / tokens) primero
	var tokenItem *database.Item
	var tokenQty uint32
	if tbl := h.getSHN("ItemMoney"); tbl != nil {
		if row := tbl.FindRow("IM_ItemIDX", info.InxName); row != nil {
			moneyIdx := strings.TrimSpace(row.GetString("IM_MoneyIDX"))
			moneyLot := uint32(row.GetUint16("IM_MoneyLot")) * lot
			if moneyIdx != "" && moneyIdx != "-" && moneyLot > 0 {
				if moneyInfo, mOk := h.itemStore.GetItemInfoByName(moneyIdx); mOk && moneyInfo != nil {
					items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
					for _, it := range items {
						if it != nil && it.StorageType == 9 && it.ItemID == moneyInfo.ID && uint32(it.Flags) >= moneyLot {
							tokenItem = it
							tokenQty = moneyLot
							break
						}
					}
					if tokenItem == nil {
						sendBuyAck(errItemBuyNoCen)
						return
					}
				}
			}
		}
	}

	totalPrice := uint64(info.BuyPrice) * uint64(lot)
	if tokenItem != nil {
		totalPrice = 0
	}
	if uint64(player.Char.TCharacterState.Money) < totalPrice {
		sendBuyAck(errItemBuyNoCen)
		return
	}

	// Verificar espacio / apilamiento en inventario (Box 9)
	maxLot := info.MaxLot
	if maxLot == 0 {
		maxLot = 1
	}
	var targetItem *database.Item
	var stacked bool
	items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
	if maxLot > 1 {
		for _, it := range items {
			if it != nil && it.StorageType == 9 && it.ItemID == itemID {
				cur := uint32(it.Flags)
				if cur == 0 {
					cur = 1
				}
				if cur+actualLot <= maxLot {
					targetItem = it
					stacked = true
					break
				}
			}
		}
	}

	if stacked && targetItem != nil {
		newFlags := targetItem.Flags + int32(actualLot)
		if err := h.charRepo.BuyItemLot(targetItem.ItemKey, player.Char.ID, 0, 0, int64(totalPrice), 0, 0); err != nil {
			sendBuyAck(errItemBuyFail)
			return
		}
		targetItem.Flags = newFlags
		_ = h.charRepo.UpdateItemFlags(targetItem.ItemKey, int(newFlags))
	} else {
		freeSlot, ok := h.findFreeInvenSlot(player.Char.ID)
		if !ok {
			sendBuyAck(errItemBuyFull)
			return
		}
		newIt := &database.Item{
			StorageType: 9,
			Owner:       player.Char.ID,
			Storage:     freeSlot,
			ItemID:      itemID,
			Flags:       int32(actualLot),
		}
		if err := h.charRepo.BuyItemAll(newIt.ItemKey, player.Char.ID, 9, freeSlot, int(itemID), int(actualLot), int64(totalPrice), 0, 0); err != nil {
			sendBuyAck(errItemBuyFail)
			return
		}
		// Recargar llave desde inventario o crear si Memory/SP asignó una
		reloaded, _ := h.charRepo.GetCharacterItems(player.Char.ID)
		for _, r := range reloaded {
			if r.StorageType == 9 && r.Storage == freeSlot && r.ItemID == itemID {
				targetItem = r
				break
			}
		}
		if targetItem == nil {
			targetItem = newIt
		}
	}

	if tokenItem != nil && tokenQty > 0 {
		h.consumeItemUnits(player, tokenItem, tokenQty)
	}
	player.Char.TCharacterState.Money -= int64(totalPrice)
	h.syncPlayerMoney(player)

	if targetItem != nil {
		h.sendCellChangeFull(client, 9, targetItem.Storage, targetItem, 0xFFFF)
	}
	sendBuyAck(errItemBuySuccess)

	// Notificar acción de compra exitosa (0x3047 PROTO_NC_ITEM_BUY_SUC_ACTION_CMD)
	sucAct := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_SUC_ACTION_CMD)
	_ = client.SendPacket(sucAct)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType3(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(itemID), int(player.OpenNPCMobID), int(totalPrice))
	}
}

// handleItemSellReq implementa sp_NC_ITEM_SELL_REQ (0x3006) con soporte de Recompra (0x305C).
func (h *Handler) handleItemSellReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	lot, err := pkt.Reader.ReadUint32()
	if err != nil {
		return
	}

	sendSellAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SELL_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendSellAck(errItemSellState)
		return
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(slot))
	if srcItem == nil || lot == 0 {
		sendSellAck(errItemSellMissing)
		return
	}
	rule := h.getBelongRule(srcItem.ItemID)
	if rule.NoSell {
		sendSellAck(errItemSellNoSell)
		return
	}

	curLot := uint32(srcItem.Flags)
	if curLot == 0 {
		curLot = 1
	}
	if lot > curLot {
		sendSellAck(errItemSellMissing)
		return
	}

	var unitSell uint64 = 1
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(srcItem.ItemID); ok && info != nil {
			unitSell = uint64(info.SellPrice)
		}
	}
	sellTotal := unitSell * uint64(lot)
	itemOpts, _ := h.charRepo.GetItemOptions(srcItem.ItemKey)

	if lot == curLot {
		if err := h.charRepo.SellItemAll(srcItem.ItemKey, player.Char.ID, int64(sellTotal)); err != nil {
			sendSellAck(errItemSellMissing)
			return
		}
		h.sendCellChangeFull(client, 9, int16(slot), nil, 0xFFFF)
	} else {
		rem := int32(curLot - lot)
		if err := h.charRepo.SellItemLot(srcItem.ItemKey, player.Char.ID, 0, 0, int64(sellTotal)); err != nil {
			sendSellAck(errItemSellMissing)
			return
		}
		srcItem.Flags = rem
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(rem))
		h.sendCellChangeFull(client, 9, int16(slot), srcItem, 0xFFFF)
	}

	player.Char.TCharacterState.Money += int64(sellTotal)
	h.syncPlayerMoney(player)

	// Guardar entrada en historial de Recompra y notificar PROTO_NC_ITEM_SELL_ITEM_INSERT_CMD (0x305C, 105B)
	var maxH uint16
	for _, rp := range player.RepurchaseItems {
		if rp.Handle > maxH {
			maxH = rp.Handle
		}
	}
	repHandle := maxH + 1
	repEntry := &RepurchaseEntry{
		Handle:  repHandle,
		ItemID:  srcItem.ItemID,
		Lot:     lot,
		Price:   sellTotal,
		Options: itemOpts,
	}
	player.RepurchaseItems = append(player.RepurchaseItems, repEntry)
	if len(player.RepurchaseItems) > 20 {
		player.RepurchaseItems = player.RepurchaseItems[len(player.RepurchaseItems)-20:]
	}

	snapItem := &database.Item{
		ItemID: srcItem.ItemID,
		Flags:  int32(lot),
	}
	infoBytes := h.buildFixed103ItemInfo(snapItem)
	ins := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SELL_ITEM_INSERT_CMD)
	ins.WriteUint16(repHandle)
	ins.WriteBytes(infoBytes[:])
	_ = client.SendPacket(ins)

	sendSellAck(errItemSellSuccess)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType3(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(srcItem.ItemID), int(player.OpenNPCMobID), int(sellTotal))
	}
}

// handleItemRepurchaseReq implementa sp_NC_ITEM_REPURCHASE_REQ (0x3059).
func (h *Handler) handleItemRepurchaseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	handle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REPURCHASE_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendAck(errItemRepurchaseFail)
		return
	}

	idx := -1
	var found *RepurchaseEntry
	for i, rp := range player.RepurchaseItems {
		if rp != nil && rp.Handle == handle {
			idx = i
			found = rp
			break
		}
	}
	if found == nil {
		sendAck(errItemRepurchaseFail)
		return
	}
	if uint64(player.Char.TCharacterState.Money) < found.Price {
		sendAck(errItemRepurchaseNoCen)
		return
	}

	freeSlot, ok := h.findFreeInvenSlot(player.Char.ID)
	if !ok {
		sendAck(errItemRepurchaseFull)
		return
	}

	if err := h.charRepo.BuyItemAll(0, player.Char.ID, 9, freeSlot, int(found.ItemID), int(found.Lot), int64(found.Price), 0, 0); err != nil {
		sendAck(errItemRepurchaseFail)
		return
	}

	var restored *database.Item
	items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
	for _, it := range items {
		if it != nil && it.StorageType == 9 && it.Storage == freeSlot && it.ItemID == found.ItemID {
			restored = it
			break
		}
	}
	if restored != nil {
		for k, v := range found.Options {
			_ = h.charRepo.SetItemOption(restored.ItemKey, k, v)
		}
		h.sendCellChangeFull(client, 9, freeSlot, restored, 0xFFFF)
	}

	player.RepurchaseItems = append(player.RepurchaseItems[:idx], player.RepurchaseItems[idx+1:]...)
	player.Char.TCharacterState.Money -= int64(found.Price)
	h.syncPlayerMoney(player)

	delCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SELL_ITEM_DELETE_CMD)
	delCmd.WriteUint16(handle)
	_ = client.SendPacket(delCmd)

	sendAck(errItemRepurchaseSuccess)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType3(3, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(found.ItemID), int(handle), int(found.Price))
	}
}

// handleItemDepositReq implementa sp_NC_ITEM_DEPOSIT_REQ (0x301C -> p_User_Deposit).
func (h *Handler) handleItemDepositReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}

	sendAck := func(remStorage uint64, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DEPOSIT_ACK)
		ack.WriteUint64(remStorage)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageState)
		return
	}
	if cen == 0 || uint64(player.Char.TCharacterState.Money) < cen {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageNoCen)
		return
	}

	storedCen, err := h.charRepo.UserDepositMoney(player.Char.AccountID, player.Char.ID, int64(cen))
	if err != nil {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageFail)
		return
	}
	player.Char.TCharacterState.Money -= int64(cen)
	h.syncPlayerMoney(player)
	sendAck(uint64(storedCen), errItemStorageSuccess)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueType5(1, player.Char.ID, int(cen), int(storedCen))
	}
}

// handleItemWithdrawReq implementa sp_NC_ITEM_WITHDRAW_REQ (0x301E -> p_User_Withdraw).
func (h *Handler) handleItemWithdrawReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}

	sendAck := func(remStorage uint64, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_WITHDRAW_ACK)
		ack.WriteUint64(remStorage)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageState)
		return
	}
	if cen == 0 {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageNoCen)
		return
	}

	storedCen, err := h.charRepo.UserWithdrawMoney(player.Char.AccountID, player.Char.ID, int64(cen))
	if err != nil {
		cur, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
		sendAck(uint64(cur), errItemStorageNoCen)
		return
	}
	player.Char.TCharacterState.Money += int64(cen)
	h.syncPlayerMoney(player)
	sendAck(uint64(storedCen), errItemStorageSuccess)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueType5(2, player.Char.ID, int(cen), int(storedCen))
	}
}

// sendStoragePage envia PROTO_NC_MENU_OPENSTORAGE_CMD (0x3C08) para almacén personal (box 10) o de cuenta (box 6).
func (h *Handler) sendStoragePage(player *ZonePlayer, box uint8, page uint8, openType uint8) {
	if player == nil || player.Client == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	storedCen, _ := h.charRepo.UserGetMoney(player.Char.AccountID)
	items, _ := h.charRepo.GetCharacterItems(player.Char.ID)
	var list []*database.Item
	for _, it := range items {
		if it != nil && it.StorageType == box {
			list = append(list, it)
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MENU_OPENSTORAGE_CMD)
	pw.WriteUint64(uint64(storedCen))
	pw.WriteByte(4)        // maxpage
	pw.WriteByte(page)     // curpage
	pw.WriteByte(openType) // nOpenType
	if len(list) > 250 {
		list = list[:250]
	}
	pw.WriteByte(byte(len(list)))
	for _, it := range list {
		pw.WriteBytes(h.buildItemPacketInformWithLoc(box, it.Storage, it))
	}
	_ = player.Client.SendPacket(pw)
}

func (h *Handler) handleItemOpenStoragePageReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	page, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	box := uint8(10)
	if player.AccountStorageOpen {
		box = 6
	}
	h.sendStoragePage(player, box, page, 0)
}

func (h *Handler) handleAccountStorageOpenCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.AccountStorageOpen = true
	h.sendStoragePage(player, 6, 0, 1)
}

func (h *Handler) handleAccountStorageCloseCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	player.AccountStorageOpen = false
}

// handleItemGuildStorageWithdrawReq implementa sp_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ (0x3026 -> p_Guild_MoneySub).
func (h *Handler) handleItemGuildStorageWithdrawReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}

	sendAck := func(rem uint64, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK)
		ack.WriteUint64(rem)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendAck(0, errGuildStorageState)
		return
	}
	gid := h.resolvePlayerGuildID(player)
	if gid <= 0 {
		sendAck(0, errGuildStorageNoGuild)
		return
	}
	if cen == 0 {
		sendAck(0, errGuildStorageNoCen)
		return
	}
	rem, _, err := h.charRepo.SubGuildMoney(gid, player.Char.ID, int64(cen))
	if err != nil {
		sendAck(0, errGuildStorageDenied)
		return
	}
	player.Char.TCharacterState.Money += int64(cen)
	h.syncPlayerMoney(player)
	sendAck(uint64(rem), errGuildStorageSuccess)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeB(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, gid, 0, int64(cen), int(rem), 0, 0)
	}
}

func (h *Handler) resolvePlayerGuildID(player *ZonePlayer) int {
	if player == nil || player.Char == nil {
		return 0
	}
	if h.charRepo != nil {
		if g, _, err := h.charRepo.GetGuildByCharID(player.Char.ID); err == nil && g != nil {
			return g.ID
		}
		if ac, _, err := h.charRepo.GetAcademyByCharID(player.Char.ID); err == nil && ac != nil {
			return ac.GuildID
		}
	}
	return 0
}

func (h *Handler) handleItemAcademyStorageWithdrawReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}
	sendAck := func(rem uint64, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK)
		ack.WriteUint64(rem)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}
	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendAck(0, errAcademyWithdrawState)
		return
	}
	gid := h.resolvePlayerGuildID(player)
	if gid <= 0 {
		sendAck(0, errAcademyWithdrawNoGuild)
		return
	}
	if cen == 0 {
		cur, _ := h.charRepo.GetAcademyMoney(gid)
		sendAck(uint64(cur), errAcademyWithdrawFail)
		return
	}
	if err := h.charRepo.WithdrawAcademyMoney(gid, player.Char.ID, int64(cen)); err != nil {
		cur, _ := h.charRepo.GetAcademyMoney(gid)
		sendAck(uint64(cur), errAcademyWithdrawNoCen)
		return
	}
	rem, _ := h.charRepo.GetAcademyMoney(gid)
	player.Char.TCharacterState.Money += int64(cen)
	h.syncPlayerMoney(player)
	sendAck(uint64(rem), errAcademyWithdrawSuccess)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeB(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, gid, 0, int64(cen), int(rem), 0, 0)
	}
}

func (h *Handler) handleItemAcademyStorageDepositReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}
	sendAck := func(rem uint64, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK)
		ack.WriteUint64(rem)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}
	if player.State == 3 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendAck(0, errAcademyDepositState)
		return
	}
	gid := h.resolvePlayerGuildID(player)
	if gid <= 0 {
		sendAck(0, errAcademyDepositNoGuild)
		return
	}
	if cen == 0 || uint64(player.Char.TCharacterState.Money) < cen {
		cur, _ := h.charRepo.GetAcademyMoney(gid)
		sendAck(uint64(cur), errAcademyDepositNoCen)
		return
	}
	if err := h.charRepo.DepositAcademyMoney(gid, player.Char.ID, int64(cen)); err != nil {
		cur, _ := h.charRepo.GetAcademyMoney(gid)
		sendAck(uint64(cur), errAcademyDepositFail)
		return
	}
	rem, _ := h.charRepo.GetAcademyMoney(gid)
	player.Char.TCharacterState.Money -= int64(cen)
	h.syncPlayerMoney(player)
	sendAck(uint64(rem), errAcademyDepositSuccess)

	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeB(3, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, gid, 0, int64(cen), int(rem), 0, 0)
	}
}

func (h *Handler) handleGuildAcademyRewardStorageOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	openType, _ := pkt.Reader.ReadByte()
	gid := h.resolvePlayerGuildID(player)
	if gid <= 0 {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK)
		ack.WriteUint16(errItemStorageFail)
		ack.WriteUint64(0)
		ack.WriteByte(0)
		_ = client.SendPacket(ack)
		return
	}
	player.AcademyStorageOpen = true
	cen, _ := h.charRepo.GetAcademyMoney(gid)
	rewards, _ := h.charRepo.GetAcademyRewardsAll(gid, int(openType))

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK)
	ack.WriteUint16(errAcademyStorageSuccess)
	ack.WriteUint64(uint64(cen))
	if len(rewards) > 250 {
		rewards = rewards[:250]
	}
	ack.WriteByte(byte(len(rewards)))
	for idx, rw := range rewards {
		cnt := int32(rw.Value)
		if cnt <= 0 {
			cnt = 1
		}
		it := &database.Item{ItemID: uint16(rw.ItemID), Flags: cnt}
		ack.WriteBytes(h.buildItemPacketInformWithLoc(0, int16(idx), it))
	}
	_ = client.SendPacket(ack)
}

func (h *Handler) handleGuildAcademyRewardStorageCloseCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	_, _ = pkt.Reader.ReadByte()
	player.AcademyStorageOpen = false
}
