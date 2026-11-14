package database

import (
	"fmt"
	"time"
)

// GetCharacterItems recupera todos los ítems asociados a un personaje invocando p_Item_GetList.
func (r *MySQLCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error) {
	rows, err := r.db.Query("CALL p_Item_GetList(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Item_GetList(%d): %w", charNo, err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		it := &Item{Owner: charNo}
		var storageType int
		var itemID, flags int
		if err := rows.Scan(&it.ItemKey, &storageType, &it.Storage, &itemID, &flags, &it.Date); err != nil {
			return nil, fmt.Errorf("error scanning item row: %w", err)
		}
		it.StorageType = byte(storageType)
		it.ItemID = uint16(itemID)
		it.Flags = int32(flags)
		items = append(items, it)
	}
	return items, nil
}

// GetItemsByStorage recupera los ítems de un tipo de almacén específico invocando p_Item_GetListType.
func (r *MySQLCharacterRepository) GetItemsByStorage(charNo int, storageType byte) ([]*Item, error) {
	rows, err := r.db.Query("CALL p_Item_GetListType(?, ?)", charNo, storageType)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Item_GetListType(%d, %d): %w", charNo, storageType, err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		it := &Item{
			Owner:       charNo,
			StorageType: storageType,
		}
		var itemID, flags int
		if err := rows.Scan(&it.ItemKey, &it.Storage, &itemID, &flags, &it.Date); err != nil {
			return nil, fmt.Errorf("error scanning storage item row: %w", err)
		}
		it.ItemID = uint16(itemID)
		it.Flags = int32(flags)
		items = append(items, it)
	}
	return items, nil
}

// CreateItem inserta un ítem en tItem invocando p_Item_Create.
func (r *MySQLCharacterRepository) CreateItem(item *Item) error {
	if item.ItemKey == 0 {
		item.ItemKey = time.Now().UnixNano()
	}
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Create(?, ?, ?, ?, ?, ?, @out_nRet)",
		item.ItemKey, item.Owner, item.StorageType, item.Storage, item.ItemID, item.Flags,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Create: %w", err)
	}
	return nil
}

// CreateItemUseEmptyStorage ubica el ítem en la primera ranura libre invocando p_Item_Create_UseEmptyStorage.
func (r *MySQLCharacterRepository) CreateItemUseEmptyStorage(itemKey int64, owner int, storageType byte, itemID int, flags int) (int16, error) {
	var storage int16
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Create_UseEmptyStorage(?, ?, ?, @out_nStorage, ?, ?, @out_nRet)",
		itemKey, owner, storageType, itemID, flags,
	).Scan(&storage, &ret)
	if err != nil {
		return -1, fmt.Errorf("error calling p_Item_Create_UseEmptyStorage: %w", err)
	}
	if ret != 0 {
		return -1, fmt.Errorf("p_Item_Create_UseEmptyStorage failed with ret=%d", ret)
	}
	return storage, nil
}

// SaveCharacterItem actualiza un ítem en tItem invocando p_Item_SetData.
func (r *MySQLCharacterRepository) SaveCharacterItem(item *Item) error {
	dDate := item.Date
	if dDate.IsZero() {
		dDate = time.Now()
	}
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_SetData(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		item.ItemKey, item.Owner, item.StorageType, item.Storage, item.ItemID, item.Flags, dDate,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SetData: %w", err)
	}
	return nil
}

// DeleteItem elimina un ítem y registra su respaldo en tItem_deleted invocando p_Item_Delete.
func (r *MySQLCharacterRepository) DeleteItem(itemKey int64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Item_Delete(?, @out_nRet)", itemKey).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Delete(%d): %w", itemKey, err)
	}
	return nil
}

// SetItemStorage actualiza el tipo y posición de almacén de un ítem invocando p_Item_SetStorage.
func (r *MySQLCharacterRepository) SetItemStorage(itemKey int64, storageType byte, storageSlot int16) error {
	var ret int
	err := r.db.QueryRow("CALL p_Item_SetStorage(?, ?, ?, @out_nRet)", itemKey, storageType, storageSlot).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SetStorage: %w", err)
	}
	return nil
}

// ExchangeItemStorage intercambia de forma atómica y transaccional dos ítems invocando p_Item_ExchangeStorage.
func (r *MySQLCharacterRepository) ExchangeItemStorage(itemKey1 int64, storageType1 byte, storageSlot1 int16, itemKey2 int64, storageType2 byte, storageSlot2 int16) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_ExchangeStorage(?, ?, ?, ?, ?, ?, @out_nRet)",
		itemKey1, storageType1, storageSlot1, itemKey2, storageType2, storageSlot2,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_ExchangeStorage: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_ExchangeStorage failed with ret=%d", ret)
	}
	return nil
}

// ArrangeInventory reordena secuencialmente las ranuras de inventario invocando p_Item_SetStorage.
func (r *MySQLCharacterRepository) ArrangeInventory(charNo int, itemKeys []int64) error {
	for i, key := range itemKeys {
		var ret int
		err := r.db.QueryRow("CALL p_Item_SetStorage(?, ?, ?, @out_nRet)", key, 9, int16(i)).Scan(&ret)
		if err != nil {
			return fmt.Errorf("error arranging item %d: %w", key, err)
		}
	}
	return nil
}

// UpdateItemFlags actualiza los flags de un ítem mediante p_Item_GetData + p_Item_SetData.
func (r *MySQLCharacterRepository) UpdateItemFlags(itemKey int64, flags int) error {
	var owner, st, slot, itemID, oldFlags, ret int
	var dDate time.Time
	err := r.db.QueryRow(
		"CALL p_Item_GetData(?, @out_nOwner, @out_nStorageType, @out_nStorage, @out_nItemID, @out_nFlags, @out_dDate, @out_nRet)",
		itemKey,
	).Scan(&owner, &st, &slot, &itemID, &oldFlags, &dDate, &ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_GetData(%d): %w", itemKey, err)
	}
	if ret == 0 {
		return fmt.Errorf("p_Item_GetData(%d) not found", itemKey)
	}
	return r.SetItemData(itemKey, owner, byte(st), int16(slot), itemID, flags, dDate)
}

// GetItemOptions recupera las opciones aleatorias y atributos de un ítem invocando p_Item_GetOptionAll.
func (r *MySQLCharacterRepository) GetItemOptions(itemKey int64) (map[int16]int64, error) {
	rows, err := r.db.Query("CALL p_Item_GetOptionAll(?)", itemKey)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Item_GetOptionAll(%d): %w", itemKey, err)
	}
	defer rows.Close()

	opts := make(map[int16]int64)
	for rows.Next() {
		var optType int16
		var optData int64
		if err := rows.Scan(&optType, &optData); err == nil {
			opts[optType] = optData
		}
	}
	return opts, nil
}

// SetItemOption persiste una opción sobre un ítem invocando p_Item_SetOption.
func (r *MySQLCharacterRepository) SetItemOption(itemKey int64, optionType int16, optionData int64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Item_SetOption(?, ?, ?, @out_nRet)", itemKey, optionType, optionData).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SetOption: %w", err)
	}
	return nil
}

// GetItemEnchant recupera las ranuras de gemas asociadas al equipamiento invocando usp_Item_Get_EnchantAll.
func (r *MySQLCharacterRepository) GetItemEnchant(itemKey int64) ([]*ItemEnchantData, error) {
	rows, err := r.db.Query("CALL usp_Item_Get_EnchantAll(?, @out_nRet)", itemKey)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Item_Get_EnchantAll(%d): %w", itemKey, err)
	}
	defer rows.Close()

	var enchants []*ItemEnchantData
	for rows.Next() {
		var slot, restCount byte
		var gemItemID int
		if err := rows.Scan(&slot, &gemItemID, &restCount); err == nil {
			enchants = append(enchants, &ItemEnchantData{
				Slot:      slot,
				GemItemID: gemItemID,
				RestCount: restCount,
			})
		}
	}
	return enchants, nil
}

// SetItemData actualiza el registro completo de un ítem invocando p_Item_SetData.
func (r *MySQLCharacterRepository) SetItemData(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, dDate time.Time) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_SetData(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		itemKey, owner, storageType, storage, itemID, flags, dDate,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SetData: %w", err)
	}
	return nil
}

// MergeItems fusiona dos stacks de ítems invocando p_Item_Merge.
func (r *MySQLCharacterRepository) MergeItems(ownerFrom, ownerTo int, keyFrom, keyTo int64, optType int16, optData int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Merge(?, ?, ?, ?, ?, ?, @out_nRet)",
		ownerFrom, ownerTo, keyFrom, keyTo, optType, optData,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Merge: %w", err)
	}
	return nil
}

// SplitAndMergeItems divide y combina stacks de ítems invocando p_Item_Split_N_Merge.
func (r *MySQLCharacterRepository) SplitAndMergeItems(keyFrom int64, ownerFrom int, optTypeFrom int16, optDataFrom int64, keyTo int64, ownerTo int, optTypeTo int16, optDataTo int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Split_N_Merge(?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		keyFrom, ownerFrom, optTypeFrom, optDataFrom, keyTo, ownerTo, optTypeTo, optDataTo,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Split_N_Merge: %w", err)
	}
	return nil
}

// UpgradeItem ejecuta la mejora o destrucción de un ítem invocando p_Item_Upgrade.
func (r *MySQLCharacterRepository) UpgradeItem(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64) error {
	return r.UpgradeItemFull(owner, itemKey, itemDestroy, upOptType, upOptValue, upFailCountType, upFailCountValue, 0, 0, 0, 0, 0)
}

// UpgradeItemFull ejecuta la mejora completa incluyendo catalizadores y costo en Cen vía p_Item_Upgrade.
func (r *MySQLCharacterRepository) UpgradeItemFull(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64, rawKey, rawLeftKey, rawRightKey, rawMiddleKey, giftMoney int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Upgrade(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		owner, itemKey, itemDestroy, upOptType, upOptValue, upFailCountType, upFailCountValue,
		rawKey, rawLeftKey, rawRightKey, rawMiddleKey, giftMoney,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Upgrade: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_Upgrade failed with code %d", ret)
	}
	return nil
}

// BuyItemAll compra un ítem nuevo atómicamente vía p_Item_BuyAll.
func (r *MySQLCharacterRepository) BuyItemAll(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, price int64, fame int, coin int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_BuyAll(?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		itemKey, owner, storageType, storage, itemID, flags, price, fame, coin,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_BuyAll: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_BuyAll failed with code %d", ret)
	}
	return nil
}

// BuyItemLot añade cantidad a un stack existente atómicamente vía p_Item_BuyLot.
func (r *MySQLCharacterRepository) BuyItemLot(itemKey int64, owner int, optType int16, optData int64, price int64, fame int, coin int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_BuyLot(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		itemKey, owner, optType, optData, price, fame, coin,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_BuyLot: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_BuyLot failed with code %d", ret)
	}
	return nil
}

// SellItemAll vende un ítem completo atómicamente vía p_Item_SellAll.
func (r *MySQLCharacterRepository) SellItemAll(itemKey int64, owner int, price int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_SellAll(?, ?, ?, @out_nRet)",
		itemKey, owner, price,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SellAll: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_SellAll failed with code %d", ret)
	}
	return nil
}

// SellItemLot vende parte de un stack atómicamente vía p_Item_SellLot.
func (r *MySQLCharacterRepository) SellItemLot(itemKey int64, owner int, optType int16, optData int64, price int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_SellLot(?, ?, ?, ?, ?, @out_nRet)",
		itemKey, owner, optType, optData, price,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_SellLot: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_SellLot failed with code %d", ret)
	}
	return nil
}

// UserDepositMoney deposita Cen del personaje al almacén de cuenta vía p_User_Deposit.
func (r *MySQLCharacterRepository) UserDepositMoney(userNo, charNo int, amount int64) (int64, error) {
	var totalMoney int64
	var ret int
	err := r.db.QueryRow(
		"CALL p_User_Deposit(?, ?, ?, @out_nResultUserMoney, @out_nRet)",
		userNo, charNo, amount,
	).Scan(&totalMoney, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_User_Deposit: %w", err)
	}
	if ret != 0 {
		return 0, fmt.Errorf("p_User_Deposit failed with code %d", ret)
	}
	return totalMoney, nil
}

// UserWithdrawMoney retira Cen del almacén de cuenta al personaje vía p_User_Withdraw.
func (r *MySQLCharacterRepository) UserWithdrawMoney(userNo, charNo int, amount int64) (int64, error) {
	var totalMoney int64
	var ret int
	err := r.db.QueryRow(
		"CALL p_User_Withdraw(?, ?, ?, @out_nResultUserMoney, @out_nRet)",
		userNo, charNo, amount,
	).Scan(&totalMoney, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_User_Withdraw: %w", err)
	}
	if ret != 0 {
		return 0, fmt.Errorf("p_User_Withdraw failed with code %d", ret)
	}
	return totalMoney, nil
}

// UserGetMoney consulta el saldo del almacén de cuenta vía p_User_MoneyGet.
func (r *MySQLCharacterRepository) UserGetMoney(userNo int) (int64, error) {
	var totalMoney int64
	var ret int
	err := r.db.QueryRow(
		"CALL p_User_MoneyGet(?, @out_nUserMoney, @out_nRet)",
		userNo,
	).Scan(&totalMoney, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_User_MoneyGet: %w", err)
	}
	return totalMoney, nil
}

// EnchantAddGem engarza una gema en un socket del equipamiento vía usp_Item_Add_EnchantGem.
func (r *MySQLCharacterRepository) EnchantAddGem(charNo int, equipKey, gemKey int64, socketIdx byte, gemID int, restGemCount byte) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_Add_EnchantGem(?, ?, ?, ?, ?, ?, @out_nRet)",
		charNo, equipKey, gemKey, socketIdx, gemID, restGemCount,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Item_Add_EnchantGem: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Item_Add_EnchantGem failed with code %d", ret)
	}
	return nil
}

// EnchantAddSocket abre un nuevo socket vía usp_Item_Add_EnchantSocket.
func (r *MySQLCharacterRepository) EnchantAddSocket(charNo int, equipKey, drillKey int64, drillLotType int16) (byte, int64, error) {
	var sockCount byte
	var drillLot int64
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_Add_EnchantSocket(?, ?, ?, @out_nAllSocketCount, ?, @out_nDrillLotValue, @out_nRet)",
		charNo, equipKey, drillKey, drillLotType,
	).Scan(&sockCount, &drillLot, &ret)
	if err != nil {
		return 0, 0, fmt.Errorf("error calling usp_Item_Add_EnchantSocket: %w", err)
	}
	if ret != 0 {
		return 0, 0, fmt.Errorf("usp_Item_Add_EnchantSocket failed with code %d", ret)
	}
	return sockCount, drillLot, nil
}

// EnchantRemoveGem extrae o limpia una gema de un socket vía usp_Item_Del_EnchantGem.
func (r *MySQLCharacterRepository) EnchantRemoveGem(charNo int, equipKey, removerKey int64, socketIdx, needItem byte, removerLotType int16) (int64, error) {
	var remLot int64
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_Del_EnchantGem(?, ?, ?, ?, ?, ?, @out_nRemoverLotValue, @out_nRet)",
		charNo, equipKey, removerKey, socketIdx, needItem, removerLotType,
	).Scan(&remLot, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling usp_Item_Del_EnchantGem: %w", err)
	}
	if ret != 0 {
		return 0, fmt.Errorf("usp_Item_Del_EnchantGem failed with code %d", ret)
	}
	return remLot, nil
}

// SetEnchantGemLot actualiza las cargas restantes de la gema engarzada vía usp_Item_Set_EnchantGemLot.
func (r *MySQLCharacterRepository) SetEnchantGemLot(charNo int, equipKey int64, slot, restGemCount byte) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_Set_EnchantGemLot(?, ?, ?, ?, @out_nRet)",
		charNo, equipKey, slot, restGemCount,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Item_Set_EnchantGemLot: %w", err)
	}
	return nil
}

// PutOnBelongedCancel desvincula un equipamiento ligado consumiendo el pergamino vía usp_Item_setPutOnBelongedCancel.
func (r *MySQLCharacterRepository) PutOnBelongedCancel(charNo int, cancelKey int64, cancelOptType int16, cancelOptData int64, targetKey int64, targetOptType int16, targetOptData int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_setPutOnBelongedCancel(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		charNo, cancelKey, cancelOptType, cancelOptData, targetKey, targetOptType, targetOptData,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Item_setPutOnBelongedCancel: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Item_setPutOnBelongedCancel failed with code %d", ret)
	}
	return nil
}

// SetMinimonAutoItem configura los ítems automáticos de Minimon vía usp_Item_setAutoMon.
func (r *MySQLCharacterRepository) SetMinimonAutoItem(charNo, monType, slot, itemID int) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Item_setAutoMon(?, ?, ?, ?, @out_nRet)",
		charNo, monType, slot, itemID,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Item_setAutoMon: %w", err)
	}
	return nil
}

// FurnitureEndureUpdate repara durabilidad de mueble de MiniHouse vía p_Item_Furniture_Endure_Update.
func (r *MySQLCharacterRepository) FurnitureEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, furnitureKey int64, endureOptType int16, endureOptData int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_Furniture_Endure_Update(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		owner, kitKey, kitOptType, kitLot, furnitureKey, endureOptType, endureOptData,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_Furniture_Endure_Update: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_Furniture_Endure_Update failed with code %d", ret)
	}
	return nil
}

// CostumeWeaponEndureUpdate repara durabilidad/tiempo de arma/traje vía p_Item_CostumWeapon_Endure_Update.
func (r *MySQLCharacterRepository) CostumeWeaponEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, costumeKey int64, endureOptType int16, endureOptData int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Item_CostumWeapon_Endure_Update(?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		owner, kitKey, kitOptType, kitLot, costumeKey, endureOptType, endureOptData,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Item_CostumWeapon_Endure_Update: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("p_Item_CostumWeapon_Endure_Update failed with code %d", ret)
	}
	return nil
}

// SafeTradeAll transfiere ítem completo y/o Cen atómicamente vía usp_SafeTrade_setAll.
func (r *MySQLCharacterRepository) SafeTradeAll(buyer, seller int, moveMoney, itemKey int64, storageType byte, storage int16) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_SafeTrade_setAll(?, ?, ?, ?, ?, ?, @out_nRet)",
		buyer, seller, moveMoney, itemKey, storageType, storage,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_SafeTrade_setAll: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_SafeTrade_setAll failed with code %d", ret)
	}
	return nil
}

// SafeTradeLot divide parte de un stack hacia el comprador vía usp_SafeTrade_setLot.
func (r *MySQLCharacterRepository) SafeTradeLot(buyer int, buyerItemKey int64, buyerStorageType byte, buyerStorage int16, buyerItemID, buyerFlags int, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_SafeTrade_setLot(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		buyer, buyerItemKey, buyerStorageType, buyerStorage, buyerItemID, buyerFlags, buyerOptType,
		seller, sellerItemKey, sellerOptType, moveMoney, splitLot,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_SafeTrade_setLot: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_SafeTrade_setLot failed with code %d", ret)
	}
	return nil
}

// SafeTradeMerge fusiona parte de un stack con un stack existente del comprador vía usp_SafeTrade_setMerge.
func (r *MySQLCharacterRepository) SafeTradeMerge(buyer int, buyerItemKey int64, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64, isDelete byte) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_SafeTrade_setMerge(?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		buyer, buyerItemKey, buyerOptType, seller, sellerItemKey, sellerOptType, moveMoney, splitLot, isDelete,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_SafeTrade_setMerge: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_SafeTrade_setMerge failed with code %d", ret)
	}
	return nil
}

// GetItemActionCoolTime consulta los tiempos de recarga de grupos de ítems invocando usp_Character_getItemActionCoolTime.
func (r *MySQLCharacterRepository) GetItemActionCoolTime(charNo int) ([]*ItemActionCoolTimeData, error) {
	rows, err := r.db.Query("CALL usp_Character_getItemActionCoolTime(?, @out_nRet)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Character_getItemActionCoolTime(%d): %w", charNo, err)
	}
	defer rows.Close()

	var list []*ItemActionCoolTimeData
	for rows.Next() {
		var cNo int
		var group int16
		var endTime time.Time
		if err := rows.Scan(&cNo, &group, &endTime); err == nil {
			list = append(list, &ItemActionCoolTimeData{
				CharNo:      cNo,
				ActionGroup: group,
				EndTime:     endTime,
			})
		}
	}
	return list, nil
}

// SetItemActionCoolTime registra el tiempo de expiración para un grupo de acciones de ítem invocando usp_Character_setItemActionCoolTime.
func (r *MySQLCharacterRepository) SetItemActionCoolTime(charNo int, actionGroup int16, endTime time.Time) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Character_setItemActionCoolTime(?, ?, ?, @out_nRet)",
		charNo, actionGroup, endTime,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setItemActionCoolTime: %w", err)
	}
	return nil
}
