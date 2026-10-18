package npcitem

import (
	"path/filepath"
	"testing"
)

func TestLoadNPCItemList(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "NPCItemList")

	db, err := LoadFromDir(dirPath)
	if err != nil {
		t.Fatalf("Error cargando directorio NPCItemList: %v", err)
	}

	if db.TotalMerchants() != 72 {
		t.Errorf("Esperado 72 comerciantes, obtenido %d", db.TotalMerchants())
	}

	if db.TotalTabs() != 238 {
		t.Errorf("Esperado 238 pestañas totales, obtenido %d", db.TotalTabs())
	}

	if db.TotalValidItems() != 4997 {
		t.Errorf("Esperado 4997 ítems válidos en venta, obtenido %d", db.TotalValidItems())
	}

	// 1. Probar EldSmithKarls
	karls, ok := db.GetMerchant("EldSmithKarls")
	if !ok {
		t.Fatalf("No se encontró comerciante EldSmithKarls")
	}
	if len(karls.Tabs) < 2 {
		t.Fatalf("EldSmithKarls esperado al menos 2 tabs, obtenido %d", len(karls.Tabs))
	}
	if !karls.HasItem("RoumenSword") {
		t.Errorf("EldSmithKarls debería vender RoumenSword")
	}

	// Slot testing: Tab00, Row 0, Col 0 -> Slot 0 ("RoumenSword")
	tab0Items := karls.GetItems(0)
	if len(tab0Items) == 0 {
		t.Fatalf("Tab 0 de Karls no debería estar vacío")
	}
	if tab0Items[0].Slot != 0 || tab0Items[0].Item != "RoumenSword" {
		t.Errorf("Esperado slot 0 RoumenSword, obtenido slot %d item %s", tab0Items[0].Slot, tab0Items[0].Item)
	}

	// Slot testing: Tab00, Row 1, Col 0 -> Slot 6 ("StealSword")
	var stealSword *ShopItem
	for i := range tab0Items {
		if tab0Items[i].Item == "StealSword" {
			stealSword = &tab0Items[i]
			break
		}
	}
	if stealSword == nil {
		t.Fatalf("No se encontró StealSword en tab 0 de Karls")
	}
	if stealSword.Slot != 6 {
		t.Errorf("StealSword esperado slot 6 (row 1, col 0), obtenido %d", stealSword.Slot)
	}

	// 2. Probar AdlSkillEdwina (13 pestañas)
	edwina, ok := db.GetMerchant("AdlSkillEdwina")
	if !ok {
		t.Fatalf("No se encontró comerciante AdlSkillEdwina")
	}
	if len(edwina.Tabs) != 13 {
		t.Errorf("AdlSkillEdwina esperado 13 pestañas, obtenido %d", len(edwina.Tabs))
	}

	// 3. Probar comerciante inexistente
	if _, ok := db.GetMerchant("NonExistentNPC"); ok {
		t.Errorf("NonExistentNPC no debería existir")
	}
}
