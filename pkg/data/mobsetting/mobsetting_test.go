package mobsetting

import (
	"path/filepath"
	"testing"
)

func TestLoadMobSetting(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "MobSetting")

	db, err := LoadFromDir(dirPath)
	if err != nil {
		t.Fatalf("Error al cargar MobSetting desde dir: %v", err)
	}

	if db.TotalEntities() != 56 {
		t.Errorf("Esperado 56 entidades, obtenido %d", db.TotalEntities())
	}

	if db.TotalTriggers() != 106 {
		t.Errorf("Esperado 106 triggers totales, obtenido %d", db.TotalTriggers())
	}

	// 1. Probar Gate_Lab
	if !db.HasTriggers("Gate_Lab") {
		t.Errorf("Gate_Lab debería tener triggers registrados")
	}
	gateLab := db.GetTriggers("Gate_Lab")
	if len(gateLab) != 1 {
		t.Fatalf("Gate_Lab esperado 1 trigger, obtenido %d", len(gateLab))
	}
	if gateLab[0].ConditionID != 1 || gateLab[0].ActionID != 1 {
		t.Errorf("Gate_Lab trigger inesperado: %+v", gateLab[0])
	}

	// 2. Probar DT_RadionOre
	if !db.HasTriggers("DT_RadionOre") {
		t.Errorf("DT_RadionOre debería tener triggers registrados")
	}
	radion := db.GetTriggers("DT_RadionOre")
	if len(radion) != 3 {
		t.Fatalf("DT_RadionOre esperado 3 triggers, obtenido %d", len(radion))
	}
	for _, tr := range radion {
		if tr.ConditionID != 20 {
			t.Errorf("DT_RadionOre ConditionID esperado 20, obtenido %d", tr.ConditionID)
		}
	}

	// 3. Probar ChristmasTree
	tree := db.GetTriggers("ChristmasTree")
	if len(tree) != 20 {
		t.Errorf("ChristmasTree esperado 20 triggers, obtenido %d", len(tree))
	}

	// 4. Probar entidad inexistente
	if db.HasTriggers("NonExistentMob") {
		t.Errorf("NonExistentMob no debería tener triggers")
	}
	if len(db.GetTriggers("NonExistentMob")) != 0 {
		t.Errorf("GetTriggers para mob inexistente debería retornar slice vacío")
	}
}
