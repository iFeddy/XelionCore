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
	BDT  map[string][]byte               // 76 archivos (.bdt)
	SHAB map[string]*ArenaAttribute      // 8 arenas (.shab)
	SHMD map[string]string               // 2 descriptores (.shmd)
	SHAD map[string]*ArenaAttribute      // 1 atributo extendido (.shad)
}

// NewStore crea una instancia vacía del repositorio de geometría.
func NewStore() *Store {
	return &Store{
		SHBD: make(map[string]*MapBlockInformation),
		SBI:  make(map[string]*MapDoorArray),
		AID:  make(map[string]*AreaInfoData),
		BDT:  make(map[string][]byte),
		SHAB: make(map[string]*ArenaAttribute),
		SHMD: make(map[string]string),
		SHAD: make(map[string]*ArenaAttribute),
	}
}

// LoadAll escanea y carga la totalidad de los 301 archivos binarios de data/BlockInfo/.
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
		fullPath := filepath.Join(blockInfoDir, name)

		switch ext {
		case ".shbd":
			mbi, err := LoadSHBDFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SHBD[base] = mbi

		case ".sbi":
			mda, err := LoadSBIFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.SBI[base] = mda

		case ".aid":
			aid, err := LoadAIDFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.AID[base] = aid

		case ".bdt":
			data, err := LoadBDTFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error cargando %s: %w", name, err)
			}
			store.BDT[base] = data

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

	return store, nil
}

// TotalCount devuelve la cantidad total de archivos binarios indexados.
func (s *Store) TotalCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.SHBD) + len(s.SBI) + len(s.AID) + len(s.BDT) + len(s.SHAB) + len(s.SHMD) + len(s.SHAD)
}

// GetSHBD busca la matriz de bloqueo de un mapa (búsqueda insensible a mayúsculas/minúsculas).
func (s *Store) GetSHBD(mapName string) *MapBlockInformation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if mbi, ok := s.SHBD[mapName]; ok {
		return mbi
	}
	for k, v := range s.SHBD {
		if strings.EqualFold(k, mapName) {
			return v
		}
	}
	return nil
}

// GetSBI busca el conjunto de puertas base de un mapa.
func (s *Store) GetSBI(mapName string) *MapDoorArray {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if mda, ok := s.SBI[mapName]; ok {
		return mda
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

// GetAID busca las áreas de evento de un mapa.
func (s *Store) GetAID(mapName string) *AreaInfoData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if aid, ok := s.AID[mapName]; ok {
		return aid
	}
	for k, v := range s.AID {
		if strings.EqualFold(k, mapName) {
			return v
		}
	}
	return nil
}

