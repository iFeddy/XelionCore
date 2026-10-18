package mobroam_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/mobroam"
)

func TestLoadFromFile(t *testing.T) {
	dirPath := filepath.Join("..", "..", "..", "data", "MobRoam")

	db, err := mobroam.LoadFromFile(dirPath)
	if err != nil {
		t.Fatalf("LoadFromFile failed from dir: %v", err)
	}

	if len(db.Routes) != 28 {
		t.Errorf("Expected 28 routes, got %d", len(db.Routes))
	}

	// 1. Probar ruta de Golden Beach Waitress 01 (modo rotate)
	waitress, ok := db.GetRoute("GB_Waitress01")
	if !ok {
		t.Fatalf("Route GB_Waitress01 not found")
	}
	if waitress.LoopMode != "rotate" {
		t.Errorf("Expected loop_mode 'rotate', got '%s'", waitress.LoopMode)
	}
	if len(waitress.Points) != 25 {
		t.Errorf("Expected 25 waypoints in GB_Waitress01, got %d", len(waitress.Points))
	}
	// Primer waypoint: (2399, 2298)
	if waitress.Points[0].X != 2399 || waitress.Points[0].Y != 2298 {
		t.Errorf("Waypoint 0 coords mismatch: (%d, %d)", waitress.Points[0].X, waitress.Points[0].Y)
	}

	// 2. Probar ruta de Garuda en Bera Valley (modo return)
	garuda, ok := db.GetRoute("BerValeGaruda")
	if !ok {
		t.Fatalf("Route BerValeGaruda not found")
	}
	if garuda.LoopMode != "return" {
		t.Errorf("Expected loop_mode 'return', got '%s'", garuda.LoopMode)
	}
	if len(garuda.Points) != 38 {
		t.Errorf("Expected 38 waypoints in BerValeGaruda, got %d", len(garuda.Points))
	}

	// 3. Probar ruta inexistente
	_, ok = db.GetRoute("NonExistentRoute")
	if ok {
		t.Errorf("Expected false for non existent route")
	}
}
