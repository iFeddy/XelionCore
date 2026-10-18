package questparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Store holds the mapping between class names and their official numeric IDs (gQuestParserData).
type Store struct {
	mu       sync.RWMutex
	nameToID map[string]int
	idToName map[int]string
}

// NewStore initializes an empty QuestParser store.
func NewStore() *Store {
	return &Store{
		nameToID: make(map[string]int),
		idToName: make(map[int]string),
	}
}

// LoadFile reads the QuestParser.json file and populates the bidirectional class ID lookups.
func (s *Store) LoadFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	var rawMap map[string]int
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("unmarshal %s: %w", filePath, err)
	}

	normNameToID := make(map[string]int, len(rawMap))
	idToName := make(map[int]string, len(rawMap))

	for name, id := range rawMap {
		normNameToID[strings.ToLower(name)] = id
		idToName[id] = name
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.nameToID = normNameToID
	s.idToName = idToName

	return nil
}

// GetID returns the numeric class ID for a given class name (case-insensitive).
func (s *Store) GetID(className string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.nameToID[strings.ToLower(className)]
	return id, ok
}

// GetName returns the canonical class name for a given numeric ID.
func (s *Store) GetName(classID int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	name, ok := s.idToName[classID]
	return name, ok
}

// Count returns the number of registered class mappings.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.idToName)
}
