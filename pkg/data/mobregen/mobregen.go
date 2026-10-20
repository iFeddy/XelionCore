package mobregen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MobSpawnEntry define los parámetros de spawn y tiempos de regeneración de un monstruo en un área.
type MobSpawnEntry struct {
	Mob         string `json:"mob"`
	Count       int    `json:"count"`
	KillNum     int    `json:"kill_num"`
	RegStandard int    `json:"reg_standard"`
	RegMin      int    `json:"reg_min"`
	RegMax      int    `json:"reg_max"`
	Deltas      []int  `json:"deltas"`
	DeltaSecs   []int  `json:"delta_secs"`
	Enabled     bool   `json:"enabled"`
}

// Group define un área de respawn en el mapa y la lista de monstruos asignados a ella.
type Group struct {
	Enabled     bool            `json:"enabled"`
	IsFamily    bool            `json:"is_family"`
	CenterX     int             `json:"center_x"`
	CenterY     int             `json:"center_y"`
	Width       int             `json:"width"`
	Height      int             `json:"height"`
	RangeDegree int             `json:"range_degree"`
	Mobs        []MobSpawnEntry `json:"mobs"`
}

// MapRegen almacena en memoria todos los grupos de respawn indexados para un mapa.
type MapRegen struct {
	MapName string
	Groups  map[string]Group
}

// LoadMapRegen carga y parsea la base de datos de respawn para un mapa determinado.
// Soporta rutas directas a JSON, o nombre de mapa buscando en data/MobRegen, Instant/ y KingdomQuest/.
func LoadMapRegen(baseDir, mapName string) (*MapRegen, error) {
	cleanName := strings.TrimSuffix(mapName, ".json")
	cleanName = strings.TrimSuffix(cleanName, ".txt")

	candidates := []string{
		filepath.Join(baseDir, cleanName+".json"),
		filepath.Join(baseDir, "KingdomQuest", cleanName+".json"),
		filepath.Join(baseDir, "Instant", cleanName+".json"),
	}

	var foundPath string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		return nil, fmt.Errorf("mapa de regen no encontrado para '%s' en %s", mapName, baseDir)
	}

	data, err := os.ReadFile(foundPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo %s: %w", foundPath, err)
	}

	groups := make(map[string]Group)
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("error parseando JSON %s: %w", foundPath, err)
	}

	return &MapRegen{
		MapName: cleanName,
		Groups:  groups,
	}, nil
}

// LoadAll carga y parsea todas las configuraciones de respawn en baseDir y sus subdirectorios.
func LoadAll(baseDir string) (map[string]*MapRegen, error) {
	result := make(map[string]*MapRegen)
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var groups map[string]Group
		if err := json.Unmarshal(data, &groups); err != nil {
			return err
		}
		mapName := strings.TrimSuffix(info.Name(), ".json")
		result[mapName] = &MapRegen{
			MapName: mapName,
			Groups:  groups,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error cargando MobRegen desde %s: %w", baseDir, err)
	}
	return result, nil
}

// GetActiveGroups devuelve únicamente las áreas de respawn activas (descartando eventos apagados).
func (m *MapRegen) GetActiveGroups() map[string]Group {
	active := make(map[string]Group)
	for k, g := range m.Groups {
		if g.Enabled {
			active[k] = g
		}
	}
	return active
}
