#!/usr/bin/env python3
"""
generate_shn_catalog.py
Genera pkg/server/zone/shn_catalog.go con las 180 tablas .shn oficiales requeridas en disco
y las 5 referencias opcionales documentadas.
"""

with open("docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md") as f:
    text = f.read()

start = text.find("### 3.1 Las 199 Referencias")
end = text.find("### 3.2 Los 31 Archivos")
sec31 = text[start:end]

lines = sec31.splitlines()
shn_list = []
for line in lines:
    line = line.strip()
    if line.startswith("|") and not line.startswith("| #") and not line.startswith("| :-"):
        parts = [p.strip() for p in line.split("|")]
        if len(parts) >= 5:
            num = int(parts[1])
            name = parts[2].replace("`", "").strip()
            locs = [l.strip().replace("`", "") for l in parts[3].split(",")]
            shn_list.append((num, name, locs))

print(f"Extracted {len(shn_list)} SHNs from Section 3.1")
assert len(shn_list) == 180

go_code = """package zone

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// SHNTableRef define un archivo .shn oficial de Zone.exe y sus ubicaciones en data/.
type SHNTableRef struct {
	Index    int
	FileName string
	Paths    []string
}

// RequiredSHNTables contiene las 180 tablas .shn oficiales en disco requeridas por Zone.exe (Sección 3.1).
var RequiredSHNTables = []SHNTableRef{
"""

for num, name, locs in shn_list:
    paths_str = ", ".join(f'"{l}"' for l in locs)
    go_code += f'\t{{Index: {num}, FileName: "{name}", Paths: []string{{{paths_str}}}}},\n'

go_code += """}

// OptionalSHNTables documenta las 5 referencias .shn presentes en binarios pero ausentes en disco/AstraServer.
var OptionalSHNTables = map[string]string{
	"BMPIndex.shn":     "Opcional: las máscaras .bmp se resuelven directamente por nombre desde data/AreaBMP/",
	"DiecGame.shn":     "Errata tipográfica legacy en Zone.exe reemplazada por DiceGame.shn",
	"MoveMain.shn":     "Referencia legacy reemplazada por MoverMain.shn",
	"MoverUpgrade.shn": "Referencia legacy absorbida por MoverSub.shn y UpgradeInfo.shn",
	"Quest.shn":        "Referencia legacy reemplazada por QuestData.shn",
}

// LoadAllSHNTables carga e indexa en memoria RAM las 180 tablas .shn requeridas en disco.
func LoadAllSHNTables(dataDir string) (map[string]*shn.File, error) {
	tables := make(map[string]*shn.File, len(RequiredSHNTables)*2)

	for _, ref := range RequiredSHNTables {
		var loaded *shn.File
		var err error

		// Intentar cargar por cada ruta declarada
		for _, relPath := range ref.Paths {
			// Normalizar ruta: si relPath empieza con "data/", quitar el prefijo
			p := relPath
			if strings.HasPrefix(p, "data/") {
				p = strings.TrimPrefix(p, "data/")
			}
			fullPath := filepath.Join(dataDir, p)

			if _, statErr := os.Stat(fullPath); statErr == nil {
				loaded, err = shn.ReadFile(fullPath)
				if err == nil {
					break
				}
			}
		}

		// Si falló por las rutas declaradas, buscar directamente en dataDir o dataDir/View
		if loaded == nil {
			directPath := filepath.Join(dataDir, ref.FileName)
			if _, statErr := os.Stat(directPath); statErr == nil {
				loaded, err = shn.ReadFile(directPath)
			}
		}
		if loaded == nil {
			viewPath := filepath.Join(dataDir, "View", ref.FileName)
			if _, statErr := os.Stat(viewPath); statErr == nil {
				loaded, err = shn.ReadFile(viewPath)
			}
		}

		if loaded == nil {
			return nil, fmt.Errorf("error crítico cargando .shn requerido #%d (%s): %v", ref.Index, ref.FileName, err)
		}

		base := strings.TrimSuffix(ref.FileName, filepath.Ext(ref.FileName))
		tables[base] = loaded
		tables[strings.ToLower(base)] = loaded
	}

	// Comprobar las 5 referencias opcionales (documentar en log sin fallar si faltan)
	for optName, desc := range OptionalSHNTables {
		optPath := filepath.Join(dataDir, optName)
		if _, err := os.Stat(optPath); err == nil {
			f, err := shn.ReadFile(optPath)
			if err == nil {
				base := strings.TrimSuffix(optName, filepath.Ext(optName))
				tables[base] = f
				tables[strings.ToLower(base)] = f
			}
		} else {
			log.Printf("[Zone][Data] Nota: SHN opcional '%s' ausente en disco (%s)", optName, desc)
		}
	}

	return tables, nil
}
"""

with open("pkg/server/zone/shn_catalog.go", "w") as f:
    f.write(go_code)

print("pkg/server/zone/shn_catalog.go written successfully.")

