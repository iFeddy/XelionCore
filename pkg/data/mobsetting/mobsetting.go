package mobsetting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ActionTrigger representa un vínculo condicional entre una condición del jugador/mundo y una acción ejecutada por ZoneService.
// Corresponde a la especificación oficial de NPCAction::Trigger / NPCAction::Util::loadTriggers de Zone.exe.
type ActionTrigger struct {
	ConditionID uint32 `json:"condition_id"`
	ActionID    uint32 `json:"action_id"`
}

// Database mantiene en memoria el mapeo de triggers de interacción indexados por Mob o NPC.
type Database struct {
	Actions map[string][]ActionTrigger
}

// LoadFromFile carga y deserializa el registro de acciones desde mob_actions.json o la carpeta data/MobSetting.
func LoadFromFile(path string) (*Database, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		path = filepath.Join(path, "mob_actions.json")
	} else if err != nil {
		alt := filepath.Join(filepath.Dir(path), "MobSetting", filepath.Base(path))
		if _, errAlt := os.Stat(alt); errAlt == nil {
			path = alt
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo mob_actions (%s): %w", path, err)
	}

	actions := make(map[string][]ActionTrigger)
	if err := json.Unmarshal(data, &actions); err != nil {
		return nil, fmt.Errorf("error parseando JSON mob_actions (%s): %w", path, err)
	}

	return &Database{
		Actions: actions,
	}, nil
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
