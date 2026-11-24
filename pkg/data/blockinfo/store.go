package blockinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store alberga el universo completo de los 301 archivos de geometría, colisión y puertas de data/BlockInfo/.
type Store struct {
	mu   sync.RWMutex
	SHBD map[string]*MapBlockInformation // 158 mapas (.shbd)
	SBI  map[string]*MapDoorArray        // 33 mapas (.sbi)
	AID  map[string]*AreaInfoData        // 23 mapas (.aid)
	BDT  map[string][]byte               // 76 archivos (.bdt) cargados perezosamente bajo demanda
	SHAB map[string]*ArenaAttribute      // 8 arenas (.shab)
	SHMD map[string]string               // 2 descriptores (.shmd)
	SHAD map[string]*ArenaAttribute      // 1 atributo extendido (.shad)

	// Índices O(1) insensibles a mayúsculas/minúsculas y registro lazy de .bdt
	shbdLower map[string]*MapBlockInformation
	sbiLower  map[string]*MapDoorArray
	aidLower  map[string]*AreaInfoData
	bdtPaths  map[string]string // base -> fullPath
	bdtLower  map[string]string // lower(base) -> base
}

// NewStore crea una instancia vacía del repositorio de geometría.
func NewStore() *Store {
	return &Store{
		SHBD:      make(map[string]*MapBlockInformation),
		SBI:       make(map[string]*MapDoorArray),
		AID:       make(map[string]*AreaInfoData),
		BDT:       make(map[string][]byte),
		SHAB:      make(map[string]*ArenaAttribute),
		SHMD:      make(map[string]string),
		SHAD:      make(map[string]*ArenaAttribute),
		shbdLower: make(map[string]*MapBlockInformation),
		sbiLower:  make(map[string]*MapDoorArray),
		aidLower:  make(map[string]*AreaInfoData),
		bdtPaths:  make(map[string]string),
		bdtLower:  make(map[string]string),
	}
}

// LoadAll escanea e indexa la totalidad de los 301 archivos binarios de data/BlockInfo/.
// Los 76 archivos .bdt se indexan por ruta y se cargan perezosamente en GetBDT (ahorro de ~20,6 MB RAM).
func LoadAll(blockInfoDir string) (*Store, error) {
	entries, err := os.ReadDir(blockInfoDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio BlockInfo '%s': %w", blockInfoDir, err)
	}

	store := NewStore()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		base := strings.TrimSuffix(name, filepath.Ext(name))
		lowerBase := strings.ToLower(base)
		fullPath := filepath.Join(blockInfoDir, name)

		switch ext {
		case ".shbd":
			mbi, err := LoadSHBDFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SHBD[base] = mbi
			store.shbdLower[lowerBase] = mbi

		case ".sbi":
			mda, err := LoadSBIFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SBI[base] = mda
			store.sbiLower[lowerBase] = mda

		case ".aid":
			aid, err := LoadAIDFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.AID[base] = aid
			store.aidLower[lowerBase] = aid

		case ".bdt":
			if _, err := os.Stat(fullPath); err != nil {
				return nil, fmt.Errorf("error verificando %s: %w", name, err)
			}
			// Registrar la clave con slice nil (0 bytes de payload en RAM; carga perezosa en GetBDT)
			store.BDT[base] = nil
			store.bdtPaths[base] = fullPath
			store.bdtLower[lowerBase] = base

		case ".shab":
			aa, err := LoadArenaAttributeFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SHAB[base] = aa

		case ".shmd":
			txt, err := LoadSHMDFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SHMD[base] = txt

		case ".shad":
			shad, err := LoadArenaAttributeFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SHAD[base] = shad
		}
	}

	// Paridad con Zone.exe mbi_Load (0x0049E3B0): los atributos .shab se combinan byte a byte (OR) en la matriz .shbd del mapa.
	for base, aa := range store.SHAB {
		mbi := store.SHBD[base]
		if mbi == nil {
			mbi = store.shbdLower[strings.ToLower(base)]
		}
		if mbi == nil {
			continue
		}
		if len(aa.Data) == len(mbi.Data) {
			for i := range mbi.Data {
				mbi.Data[i] |= aa.Data[i]
			}
		}
	}

	return store, nil
}

// BDTCount devuelve la cantidad total de archivos .bdt indexados (lazy o cargados en RAM).
func (s *Store) BDTCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := len(s.bdtPaths)
	for k := range s.BDT {
		if _, ok := s.bdtPaths[k]; !ok {
			count++
		}
	}
	return count
}

// TotalCount devuelve la cantidad total de archivos binarios indexados (incluyendo .bdt lazy).
func (s *Store) TotalCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bdtCount := len(s.bdtPaths)
	for k := range s.BDT {
		if _, ok := s.bdtPaths[k]; !ok {
			bdtCount++
		}
	}
	return len(s.SHBD) + len(s.SBI) + len(s.AID) + bdtCount + len(s.SHAB) + len(s.SHMD) + len(s.SHAD)
}

// GetBDT retorna los bytes crudos de un archivo .bdt, cargándolo perezosamente desde disco en la primera consulta.
func (s *Store) GetBDT(mapName string) []byte {
	s.mu.RLock()
	if data := s.BDT[mapName]; len(data) > 0 {
		s.mu.RUnlock()
		return data
	}
	base, hasLower := s.bdtLower[strings.ToLower(mapName)]
	if hasLower {
		if data := s.BDT[base]; len(data) > 0 {
			s.mu.RUnlock()
			return data
		}
	}
	path := ""
	if hasLower {
		path = s.bdtPaths[base]
	} else if p, ok := s.bdtPaths[mapName]; ok {
		base = mapName
		path = p
	}
	s.mu.RUnlock()

	if path == "" {
		return nil
	}

	data, err := LoadBDTFile(path)
	if err != nil {
		return nil
	}

	s.mu.Lock()
	s.BDT[base] = data
	s.mu.Unlock()
	return data
}

// GetSHBD busca la matriz de bloqueo de un mapa en O(1) (búsqueda insensible a mayúsculas/minúsculas).
func (s *Store) GetSHBD(mapName string) *MapBlockInformation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if mbi, ok := s.SHBD[mapName]; ok {
		return mbi
	}
	if s.shbdLower != nil {
		if mbi, ok := s.shbdLower[strings.ToLower(mapName)]; ok {
			return mbi
		}
		if len(s.shbdLower) == len(s.SHBD) {
			return nil
		}
	}
	for k, v := range s.SHBD {
		if strings.EqualFold(k, mapName) {
			return v
		}
	}
	return nil
}

// GetSBI busca el conjunto de puertas base de un mapa en O(1).
func (s *Store) GetSBI(mapName string) *MapDoorArray {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if mda, ok := s.SBI[mapName]; ok {
		return mda
	}
	if s.sbiLower != nil {
		if mda, ok := s.sbiLower[strings.ToLower(mapName)]; ok {
			return mda
		}
		if len(s.sbiLower) == len(s.SBI) {
			return nil
		}
	}
	for k, v := range s.SBI {
		if strings.EqualFold(k, mapName) {
			return v
		}
	}
	return nil
}

// NewInstanceDoors produce una copia clonada e independiente de las puertas para una instancia de mapa en ejecución.
func (s *Store) NewInstanceDoors(mapName string) *MapDoorArray {
	base := s.GetSBI(mapName)
	if base == nil {
		return nil
	}
	return base.Clone()
}

// GetAID busca las áreas de evento de un mapa en O(1).
func (s *Store) GetAID(mapName string) *AreaInfoData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if aid, ok := s.AID[mapName]; ok {
		return aid
	}
	if s.aidLower != nil {
		if aid, ok := s.aidLower[strings.ToLower(mapName)]; ok {
			return aid
		}
		if len(s.aidLower) == len(s.AID) {
			return nil
		}
	}
	for k, v := range s.AID {
		if strings.EqualFold(k, mapName) {
			return v
		}
	}
	return nil
}
