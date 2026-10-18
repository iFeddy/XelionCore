package script

import (
	"path/filepath"
	"testing"
)

func TestLoadScriptDatabase(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "Script")

	db, err := LoadFromDir(dirPath)
	if err != nil {
		t.Fatalf("Error al cargar scripts: %v", err)
	}

	if db.TotalFiles() != 51 {
		t.Errorf("Esperado 51 archivos de script, obtenido %d", db.TotalFiles())
	}

	if db.TotalEntries() != 1380 {
		t.Errorf("Esperado 1380 entradas totales, obtenido %d", db.TotalEntries())
	}

	// 1. Probar JobChange1
	val, ok := db.Get("JobChange1", "StartDear0")
	if !ok || val != "Who am I? I am who you think I am." {
		t.Errorf("JobChange1.StartDear0 mismatch: %v (encontrado: %v)", val, ok)
	}

	valKQ, ok := db.Get("JobChange1", "KQReturn30")
	if !ok || valKQ != "Move to Grandmaster Sean in 30 seconds." {
		t.Errorf("JobChange1.KQReturn30 mismatch: %v", valKQ)
	}

	// 2. Probar Wedding
	valWed, ok := db.Get("Wedding", "Wedding1")
	if !ok || valWed == "" {
		t.Errorf("Wedding.Wedding1 no encontrado")
	}

	// 3. Probar clave inexistente
	if _, ok := db.Get("JobChange1", "NonExistentKey"); ok {
		t.Errorf("NonExistentKey no debería existir")
	}
}
