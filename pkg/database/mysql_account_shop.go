package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetCash consulta el saldo de Cash activo (status = 0) mediante usp_Cash_getCash.
func (r *MySQLAccountRepository) GetCash(userNo int, cashType uint8) (int, error) {
	var cash int
	err := r.db.QueryRow("CALL usp_Cash_getCash(?, ?, @c)", userNo, cashType).Scan(&cash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("usp_Cash_getCash failed: %w", err)
	}
	return cash, nil
}

// AddCash acredita saldo de Cash activo en tCash y genera registro en tCashLog mediante usp_Cash_setCash.
func (r *MySQLAccountRepository) AddCash(userNo int, cashType uint8, cash int, reference string) error {
	_, err := r.db.Exec("CALL usp_Cash_setCash(?, ?, ?, ?)", userNo, cashType, cash, reference)
	if err != nil {
		return fmt.Errorf("usp_Cash_setCash failed: %w", err)
	}
	return nil
}

// UseCash deduce saldo de Cash activo mediante usp_Cash_useCash.
func (r *MySQLAccountRepository) UseCash(userNo int, cashType uint8, cash int, reference string) error {
	_, err := r.db.Exec("CALL usp_Cash_useCash(?, ?, ?, ?)", userNo, cashType, cash, reference)
	if err != nil {
		return fmt.Errorf("usp_Cash_useCash failed: %w", err)
	}
	return nil
}

// AddCashEventL20 acredita el bono cultural o evento nivel 20 mediante usp_Event_addCash.
func (r *MySQLAccountRepository) AddCashEventL20(userNo int, cash int) error {
	_, err := r.db.Exec("CALL usp_Event_addCash(?, 0, ?)", userNo, cash)
	if err != nil {
		return fmt.Errorf("usp_Event_addCash failed: %w", err)
	}
	return nil
}

// SetSerialKey asocia un cupón serial o clave a una cuenta mediante usp_Item_setSerialKey.
func (r *MySQLAccountRepository) SetSerialKey(userNo int, serialKey int64) error {
	var nRet int
	err := r.db.QueryRow("CALL usp_Item_setSerialKey(?, ?, @r)", userNo, serialKey).Scan(&nRet)
	if err != nil {
		return fmt.Errorf("usp_Item_setSerialKey failed: %w", err)
	}
	if nRet != 0 {
		return fmt.Errorf("usp_Item_setSerialKey rejected key (ret=%d)", nRet)
	}
	return nil
}

// GetChargedItems lista todos los ítems comprados pendientes de retiro (isDraw = 0) mediante usp_Charge_ItemSelectList.
func (r *MySQLAccountRepository) GetChargedItems(userNo int) ([]*ChargedItem, error) {
	rows, err := r.db.Query("CALL usp_Charge_ItemSelectList(?)", userNo)
	if err != nil {
		return nil, fmt.Errorf("usp_Charge_ItemSelectList failed: %w", err)
	}
	defer rows.Close()

	var items []*ChargedItem
	for rows.Next() {
		var (
			rowNo   uint32
			goodsNo uint32
			amount  uint32
			regDate time.Time
		)
		if err := rows.Scan(&rowNo, &goodsNo, &amount, &regDate); err != nil {
			return nil, fmt.Errorf("error scanning charged item: %w", err)
		}
		items = append(items, &ChargedItem{
			RowNo:        rowNo,
			UserNo:       userNo,
			GoodsNo:      goodsNo,
			Amount:       amount,
			RegisterDate: regDate,
			IsDraw:       0,
		})
	}
	return items, nil
}

// InsertChargedItem inserta un ítem en el inventario premium mediante usp_Charge_ItemInsert.
func (r *MySQLAccountRepository) InsertChargedItem(item *ChargedItem) (int64, error) {
	var ret int
	err := r.db.QueryRow("CALL usp_Charge_ItemInsert(?, ?, ?, ?, @r)",
		item.UserNo, item.OrderNo, item.GoodsNo, item.Amount).Scan(&ret)
	if err != nil {
		return 0, fmt.Errorf("usp_Charge_ItemInsert failed: %w", err)
	}
	if ret != 0 {
		return 0, fmt.Errorf("usp_Charge_ItemInsert failed with code %d", ret)
	}
	return int64(item.OrderNo), nil
}

// DrawChargedItem marca un ítem premium como retirado hacia el personaje mediante usp_Charge_ItemDraw.
func (r *MySQLAccountRepository) DrawChargedItem(rowNo int64) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Charge_ItemDraw(?, 0, @r)", rowNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_Charge_ItemDraw failed: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Charge_ItemDraw failed with code %d", ret)
	}
	return nil
}

// DeleteChargedItem elimina o cancela un ítem premium mediante usp_Charge_ItemDel.
func (r *MySQLAccountRepository) DeleteChargedItem(rowNo int64) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Charge_ItemDel(?, 0, @r)", rowNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("usp_Charge_ItemDel failed: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Charge_ItemDel failed with code %d", ret)
	}
	return nil
}

// GetItemSell consulta el precio y detalles de un ítem de tienda mediante usp_Item_getItem.
func (r *MySQLAccountRepository) GetItemSell(itemNo int) (*ItemSell, error) {
	var (
		goodsNo     int
		billingNo   int
		category    string
		groupNo     int
		name        string
		imageName   string
		itemType    string
		unit        int
		unitDesc    string
		price       int
		discount    int
		sellPrice   int
		isSell      int
		description sql.NullString
	)

	err := r.db.QueryRow("CALL usp_Item_getItem(?)", itemNo).Scan(
		&goodsNo, &billingNo, &category, &groupNo, &name, &imageName,
		&itemType, &unit, &unitDesc, &price, &discount, &sellPrice, &isSell, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item %d not found", itemNo)
		}
		return nil, fmt.Errorf("usp_Item_getItem failed: %w", err)
	}

	return &ItemSell{
		ItemNo:   goodsNo,
		ItemName: name,
		Price:    sellPrice,
	}, nil
}
