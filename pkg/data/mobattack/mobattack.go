package mobattack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AttackOrder representa un paso en la rotación fija de combate de un monstruo.
type AttackOrder struct {
	Order  int    `json:"order"`
	Attack string `json:"attack"` // "-" para ataque básico normal, o nombre de habilidad (ej: "KarenWingAtk")
}

// SkillChange define una regla condicional que altera la rotación de ataque.
type SkillChange struct {
	Condition string `json:"condition"` // "HPLow", "HPLow_ChangeOrder", "OutOfRange", "NextSkill"
	Value     any    `json:"value"`     // Permile de HP (ej: 700 = 70%, 400 = 40%) o "-"
	From      any    `json:"from"`      // Paso u orden original a sustituir
	To        string `json:"to"`        // Habilidad sustituta a ejecutar
}

// MobAI almacena la secuencia de turnos y las reglas dinámicas de combate de un mob.
type MobAI struct {
	Sequence     []AttackOrder `json:"sequence"`
	SkillChanges []SkillChange `json:"skill_changes,omitempty"`
}

// Database almacena en memoria la base de datos completa de secuencias de combate de todos los mobs.
type Database struct {
	Mobs map[string]MobAI
}

// LoadFromFile carga y deserializa el catálogo de secuencias de ataque desde mob_attack_sequence.json o directorio.
func LoadFromFile(path string) (*Database, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		path = filepath.Join(path, "mob_attack_sequence.json")
	} else if err != nil {
		alt := filepath.Join(filepath.Dir(path), "MobAttackSequence", filepath.Base(path))
		if _, errAlt := os.Stat(alt); errAlt == nil {
			path = alt
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo mob_attack_sequence (%s): %w", path, err)
	}

	mobs := make(map[string]MobAI)
	if err := json.Unmarshal(data, &mobs); err != nil {
		return nil, fmt.Errorf("error parseando JSON mob_attack_sequence (%s): %w", path, err)
	}

	return &Database{
		Mobs: mobs,
	}, nil
}

// GetMobAI busca la secuencia y reglas de combate de un monstruo por su nombre técnico oficial.
func (db *Database) GetMobAI(name string) (*MobAI, bool) {
	ai, ok := db.Mobs[name]
	if !ok {
		return nil, false
	}
	return &ai, true
}
