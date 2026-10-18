package script

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScriptTable almacena el mapeo clave-valor de diálogos para un escenario específico (ScrIndex -> ScrString).
type ScriptTable map[string]string

// Database contiene todos los catálogos de texto de escenarios indexados por nombre de archivo.
// Corresponde a la especificación oficial de ShineScript::ss_Load de Zone.exe.
type Database struct {
	Tables map[string]ScriptTable
}

// LoadFromDir carga todos los archivos .json de textos de escenarios desde el directorio especificado.
func LoadFromDir(dirPath string) (*Database, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de scripts (%s): %w", dirPath, err)
	}

	tables := make(map[string]ScriptTable)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error leyendo archivo (%s): %w", fullPath, err)
		}

		var table ScriptTable
		if err := json.Unmarshal(data, &table); err != nil {
			return nil, fmt.Errorf("error parseando JSON de script (%s): %w", fullPath, err)
		}

		scriptName := strings.TrimSuffix(entry.Name(), ".json")
		tables[scriptName] = table
	}

	return &Database{
		Tables: tables,
	}, nil
}

// Get obtiene el texto correspondiente a una clave en un archivo de script.
func (db *Database) Get(scriptName, key string) (string, bool) {
	table, ok := db.Tables[scriptName]
	if !ok {
		return "", false
	}
	text, ok := table[key]
	return text, ok
}

// TotalFiles devuelve la cantidad de catálogos de script cargados.
func (db *Database) TotalFiles() int {
	return len(db.Tables)
}

// TotalEntries devuelve la cantidad total de diálogos y textos cargados en memoria.
func (db *Database) TotalEntries() int {
	total := 0
	for _, t := range db.Tables {
		total += len(t)
	}
	return total
}
