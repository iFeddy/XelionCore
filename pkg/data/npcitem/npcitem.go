package npcitem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShopTab representa una pestaña dentro de la tienda de un comerciante.
type ShopTab struct {
	Tab   string     `json:"tab"`
	Title string     `json:"title"`
	Rows  [][]string `json:"rows"`
}

// ShopItem representa un ítem individual colocado en una casilla de la tienda,
// con su slot calculado para el paquete de red PROTO_NC_MENU_SHOPOPENTABLE_CMD.
type ShopItem struct {
	Tab     int    `json:"tab"`
	TabName string `json:"tab_name"`
	Row     int    `json:"row"`
	Col     int    `json:"col"`
	Slot    uint8  `json:"slot"` // (Row * 6) + Col
	Item    string `json:"item"` // ItemInCode
}

// Merchant almacena toda la información de venta y pestañas de un comerciante.
type Merchant struct {
	Name string    `json:"merchant"`
	Tabs []ShopTab `json:"tabs"`
}

// Database mantiene en memoria todas las tiendas de NPCs del servidor.
type Database struct {
	Merchants map[string]*Merchant
}

// LoadFromDir carga todos los archivos .json de tiendas desde el directorio especificado.
func LoadFromDir(dirPath string) (*Database, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de comerciantes (%s): %w", dirPath, err)
	}

	merchants := make(map[string]*Merchant)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		mct, err := LoadMerchant(fullPath)
		if err != nil {
			return nil, err
		}
		merchants[mct.Name] = mct
	}

	return &Database{
		Merchants: merchants,
	}, nil
}

// LoadMerchant carga una tienda de un comerciante individual desde su archivo .json.
func LoadMerchant(filePath string) (*Merchant, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de comerciante (%s): %w", filePath, err)
	}

	var mct Merchant
	if err := json.Unmarshal(data, &mct); err != nil {
		return nil, fmt.Errorf("error parseando JSON de comerciante (%s): %w", filePath, err)
	}

	return &mct, nil
}

// GetMerchant obtiene un comerciante por su MobName técnico (ej: "EldSmithKarls", "RouSmithJames").
func (db *Database) GetMerchant(name string) (*Merchant, bool) {
	m, ok := db.Merchants[name]
	return m, ok
}

// TotalMerchants devuelve la cantidad de comerciantes cargados en memoria.
func (db *Database) TotalMerchants() int {
	return len(db.Merchants)
}

// TotalTabs devuelve la cantidad total de pestañas a lo largo de todos los comercios.
func (db *Database) TotalTabs() int {
	total := 0
	for _, m := range db.Merchants {
		total += len(m.Tabs)
	}
	return total
}

// TotalValidItems devuelve la cantidad total de ítems en venta (excluyendo casillas vacías "-").
func (db *Database) TotalValidItems() int {
	total := 0
	for _, m := range db.Merchants {
		total += len(m.GetAllItems())
	}
	return total
}

// GetTab obtiene una pestaña específica por su índice numérico (0, 1, ...).
func (m *Merchant) GetTab(index int) (*ShopTab, bool) {
	if index < 0 || index >= len(m.Tabs) {
		return nil, false
	}
	return &m.Tabs[index], true
}

// GetItems obtiene todos los ítems válidos de una pestaña determinada con su slot calculado.
func (m *Merchant) GetItems(tabIndex int) []ShopItem {
	tab, ok := m.GetTab(tabIndex)
	if !ok {
		return nil
	}

	var items []ShopItem
	for rIdx, row := range tab.Rows {
		for cIdx, item := range row {
			if item != "-" && item != "" {
				slot := uint8((rIdx * 6) + cIdx)
				items = append(items, ShopItem{
					Tab:     tabIndex,
					TabName: tab.Tab,
					Row:     rIdx,
					Col:     cIdx,
					Slot:    slot,
					Item:    item,
				})
			}
		}
	}
	return items
}

// GetAllItems obtiene todos los ítems en venta de todas las pestañas de este comerciante.
func (m *Merchant) GetAllItems() []ShopItem {
	var allItems []ShopItem
	for i := range m.Tabs {
		allItems = append(allItems, m.GetItems(i)...)
	}
	return allItems
}

// HasItem consulta si el comerciante vende un ítem determinado en cualquiera de sus pestañas.
func (m *Merchant) HasItem(itemInCode string) bool {
	for _, tab := range m.Tabs {
		for _, row := range tab.Rows {
			for _, item := range row {
				if item == itemInCode {
					return true
				}
			}
		}
	}
	return false
}

// GetItemAt obtiene el nombre del ítem en una posición exacta (pestaña, fila, columna).
func (m *Merchant) GetItemAt(tabIndex, row, col int) (string, bool) {
	tab, ok := m.GetTab(tabIndex)
	if !ok {
		return "", false
	}
	if row < 0 || row >= len(tab.Rows) {
		return "", false
	}
	if col < 0 || col >= len(tab.Rows[row]) {
		return "", false
	}
	item := tab.Rows[row][col]
	if item == "-" || item == "" {
		return "", false
	}
	return item, true
}
