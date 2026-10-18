package mobsetting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ActionTrigger representa un vínculo condicional entre una condición del jugador/mundo y una acción ejecutada por ZoneService.
type ActionTrigger struct {
	ConditionID uint32 `json:"condition_id"`
	ActionID    uint32 `json:"action_id"`
}

// EntitySetting almacena la configuración de triggers para una entidad individual (Mob o NPC).
type EntitySetting struct {
	Entity   string          `json:"entity"`
	Triggers []ActionTrigger `json:"triggers"`
}

// Database mantiene en memoria el mapeo de triggers de interacción indexados por Mob o NPC.
type Database struct {
	Actions map[string][]ActionTrigger
}

// LoadFromDir carga todos los archivos .json individuales desde el directorio especificado.
func LoadFromDir(dirPath string) (*Database, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de mobsetting (%s): %w", dirPath, err)
	}

	actions := make(map[string][]ActionTrigger)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error leyendo archivo (%s): %w", fullPath, err)
		}

		var setting EntitySetting
		if err := json.Unmarshal(data, &setting); err != nil {
			return nil, fmt.Errorf("error parseando JSON (%s): %w", fullPath, err)
		}

		actions[setting.Entity] = setting.Triggers
	}

	return &Database{
		Actions: actions,
	}, nil
}

// LoadFromFile carga la base de datos ya sea desde un directorio o un archivo individual.
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
		return nil, fmt.Errorf("error leyendo archivo (%s): %w", path, err)
	}

	var setting EntitySetting
	if err := json.Unmarshal(data, &setting); err == nil && setting.Entity != "" {
		return &Database{
			Actions: map[string][]ActionTrigger{
				setting.Entity: setting.Triggers,
			},
		}, nil
	}

	// Fallback para mapa directo si fuera unificado
	var direct map[string][]ActionTrigger
	if err := json.Unmarshal(data, &direct); err == nil {
		return &Database{Actions: direct}, nil
	}

	return nil, fmt.Errorf("formato no reconocido en (%s)", path)
}

// GetTriggers devuelve la lista de triggers de interacción asociados al mob o NPC indicado.
func (db *Database) GetTriggers(name string) []ActionTrigger {
	return db.Actions[name]
}

// HasTriggers indica si el mob o NPC posee al menos un trigger de acción registrado.
func (db *Database) HasTriggers(name string) bool {
	triggers, ok := db.Actions[name]
	return ok && len(triggers) > 0
}

// TotalEntities devuelve la cantidad total de mobs y NPCs con triggers configurados.
func (db *Database) TotalEntities() int {
	return len(db.Actions)
}

// TotalTriggers devuelve la cantidad total de registros Condition->Action en todo el juego.
func (db *Database) TotalTriggers() int {
	count := 0
	for _, triggers := range db.Actions {
		count += len(triggers)
	}
	return count
}
