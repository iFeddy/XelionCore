package mobroam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Waypoint representa una coordenada geográfica en la ruta de patrullaje.
type Waypoint struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

// Route almacena la secuencia ordenada de waypoints y el comportamiento de ciclo ("rotate" o "return").
type Route struct {
	LoopMode    string     `json:"loop_mode"`
	TotalPoints int        `json:"total_points"`
	Points      []Waypoint `json:"points"`
}

// Database mantiene en memoria todas las rutas de patrullaje indexadas para ZoneService.
type Database struct {
	Routes map[string]Route
}

// LoadFromFile carga y deserializa la base de datos de rutas desde mob_roam.json o la carpeta data/MobRoam.
func LoadFromFile(path string) (*Database, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		path = filepath.Join(path, "mob_roam.json")
	} else if err != nil {
		alt := filepath.Join(filepath.Dir(path), "MobRoam", filepath.Base(path))
		if _, errAlt := os.Stat(alt); errAlt == nil {
			path = alt
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo mob_roam (%s): %w", path, err)
	}

	routes := make(map[string]Route)
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("error parseando JSON mob_roam (%s): %w", path, err)
	}

	return &Database{
		Routes: routes,
	}, nil
}

// GetRoute busca una ruta de patrullaje por su nombre técnico (ej: "AdlThornR01", "GB_Waitress01").
func (db *Database) GetRoute(name string) (*Route, bool) {
	r, ok := db.Routes[name]
	if !ok {
		return nil, false
	}
	return &r, true
}
