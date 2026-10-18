package mobroam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Waypoint representa una coordenada geográfica en la ruta de patrullaje.
type Waypoint struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

// Route almacena la secuencia ordenada de waypoints y el comportamiento de ciclo ("rotate" o "return").
type Route struct {
	Route       string     `json:"route"`
	LoopMode    string     `json:"loop_mode"`
	TotalPoints int        `json:"total_points"`
	Points      []Waypoint `json:"points"`
}

// Database mantiene en memoria todas las rutas de patrullaje indexadas para ZoneService.
type Database struct {
	Routes map[string]Route
}

// LoadFromDir carga y deserializa todas las rutas desde archivos JSON individuales en data/MobRoam.
func LoadFromDir(dirPath string) (*Database, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de rutas (%s): %w", dirPath, err)
	}

	routes := make(map[string]Route)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error leyendo archivo de ruta (%s): %w", fullPath, err)
		}

		var r Route
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error parseando JSON de ruta (%s): %w", fullPath, err)
		}

		routeName := r.Route
		if routeName == "" {
			routeName = strings.TrimSuffix(entry.Name(), ".json")
		}
		routes[routeName] = r
	}

	return &Database{
		Routes: routes,
	}, nil
}

// LoadFromFile carga y deserializa la base de datos de rutas desde un directorio o archivo.
func LoadFromFile(path string) (*Database, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("ruta no encontrada (%s): %w", path, err)
	}
	if fi.IsDir() {
		return LoadFromDir(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo mob_roam (%s): %w", path, err)
	}

	var single Route
	if err := json.Unmarshal(data, &single); err == nil && single.Route != "" {
		return &Database{
			Routes: map[string]Route{single.Route: single},
		}, nil
	}

	// Fallback por si fuera un mapa unificado
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
