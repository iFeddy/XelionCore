package abstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LevelConfig almacena los parámetros exactos de un estado anormal para un nivel específico.
type LevelConfig struct {
	Level         uint32   `json:"level"`
	DurationMs    uint32   `json:"duration_ms"`
	MinDef        uint32   `json:"min_def,omitempty"`
	MaxDef        uint32   `json:"max_def,omitempty"`
	MinShield     uint32   `json:"min_shield,omitempty"`
	MaxShield     uint32   `json:"max_shield,omitempty"`
	MinDodge      uint32   `json:"min_dodge,omitempty"`
	MaxDodge      uint32   `json:"max_dodge,omitempty"`
	EvadeRate     uint32   `json:"evade_rate,omitempty"`
	DotDamage     uint32   `json:"dot_damage,omitempty"`
	DotIntervalMs uint32   `json:"dot_interval_ms,omitempty"`
	SpeedDownRate uint32   `json:"speed_down_rate,omitempty"`
	Raw           []uint32 `json:"raw"`
}

// StateConfig define la metadata y progresión completa de un estado anormal (151 niveles).
type StateConfig struct {
	ID          uint16        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MaxLevels   uint32        `json:"max_levels"`
	RecordSize  int           `json:"record_size_bytes"`
	Levels      []LevelConfig `json:"levels"`
}

// Database almacena en memoria todos los estados anormales cargados desde archivos individuales.
type Database struct {
	States map[string]StateConfig `json:"states"`
	byID   map[uint16]*StateConfig
}

// LoadFromDir carga y deserializa todos los estados anormales desde archivos individuales en data/AbState.
func LoadFromDir(dirPath string) (*Database, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio abstate (%s): %w", dirPath, err)
	}

	states := make(map[string]StateConfig)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error leyendo archivo (%s): %w", fullPath, err)
		}

		var s StateConfig
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("error parseando JSON (%s): %w", fullPath, err)
		}

		stateName := s.Name
		if stateName == "" {
			stateName = strings.TrimSuffix(entry.Name(), ".json")
		}
		states[stateName] = s
	}

	db := &Database{
		States: states,
		byID:   make(map[uint16]*StateConfig, len(states)),
	}

	for _, s := range states {
		cfg := s
		db.byID[s.ID] = &cfg
	}

	return db, nil
}

// LoadFromFile carga y deserializa la base de datos de estados anormales desde un directorio o archivo.
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
		return nil, fmt.Errorf("error leyendo archivo abstate (%s): %w", path, err)
	}

	var single StateConfig
	if err := json.Unmarshal(data, &single); err == nil && single.Name != "" {
		cfg := single
		return &Database{
			States: map[string]StateConfig{single.Name: single},
			byID:   map[uint16]*StateConfig{single.ID: &cfg},
		}, nil
	}

	states := make(map[string]StateConfig)
	if err := json.Unmarshal(data, &states); err != nil {
		return nil, fmt.Errorf("error parseando JSON abstate (%s): %w", path, err)
	}

	db := &Database{
		States: states,
		byID:   make(map[uint16]*StateConfig, len(states)),
	}

	for _, s := range states {
		cfg := s
		db.byID[s.ID] = &cfg
	}

	return db, nil
}

// GetByName busca un estado por su nombre técnico (ej: "StaHardSkin", "StaShield").
func (db *Database) GetByName(name string) (*StateConfig, bool) {
	cfg, ok := db.States[name]
	if !ok {
		return nil, false
	}
	return &cfg, true
}

// GetByID busca un estado por su ID numérico oficial (coincidente con AbState.shn).
func (db *Database) GetByID(id uint16) (*StateConfig, bool) {
	cfg, ok := db.byID[id]
	return cfg, ok
}

// GetLevelConfig obtiene los parámetros para un nivel específico de un estado.
func (s *StateConfig) GetLevelConfig(level int) (*LevelConfig, error) {
	if level < 0 || level >= len(s.Levels) {
		return nil, fmt.Errorf("nivel %d fuera de rango [0, %d] para estado %s", level, len(s.Levels)-1, s.Name)
	}
	return &s.Levels[level], nil
}
