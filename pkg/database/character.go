package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrCharacterExists   = errors.New("character name already exists")
	ErrSlotInUse         = errors.New("character slot already in use")
	ErrSlotOutOfRange    = errors.New("character slot out of range (0-5)")
)

// Character represents an avatar in XelionCore.
type Character struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id"`
	Name      string          `json:"name"`
	Slot      byte            `json:"slot"`
	Job       byte            `json:"job"`
	Level     byte            `json:"level"`
	Male      bool            `json:"male"`
	Hair      byte            `json:"hair"`
	HairColor byte            `json:"hair_color"`
	Face      byte            `json:"face"`
	Map       uint16          `json:"map"`
	XPos      int             `json:"x_pos"`
	YPos      int             `json:"y_pos"`
	CurHP     int             `json:"cur_hp"`
	CurSP     int             `json:"cur_sp"`
	TownName  string          `json:"town_name"`
	Equips    map[byte]uint16 `json:"equips"` // Slot index -> item ID
}

// GetEquip returns the equipped item ID for a slot, or 0xFFFF if empty.
func (c *Character) GetEquip(slot byte) uint16 {
	if c.Equips != nil {
		if val, exists := c.Equips[slot]; exists {
			return val
		}
	}
	return 0xFFFF
}

// CharacterRepository defines persistence operations for player characters.
type CharacterRepository interface {
	GetCharactersByAccount(accountID int) ([]*Character, error)
	GetCharacterByName(name string) (*Character, error)
	GetCharacterByID(id int) (*Character, error)
	CreateCharacter(ch *Character) error
	UpdateCharacter(ch *Character) error
	DeleteCharacter(accountID int, slot byte) error
}

type charStore struct {
	NextID     int          `json:"next_id"`
	Characters []*Character `json:"characters"`
}

// MemoryCharacterRepository provides thread-safe character storage with JSON persistence.
type MemoryCharacterRepository struct {
	mu         sync.RWMutex
	filePath   string
	characters map[int]*Character // Key: Character.ID
	nextID     int
}

// NewMemoryCharacterRepository initializes a character repository persisted at data/characters.json.
func NewMemoryCharacterRepository() *MemoryCharacterRepository {
	filePath := filepath.Join("data", "characters.json")
	repo := &MemoryCharacterRepository{
		filePath:   filePath,
		characters: make(map[int]*Character),
		nextID:     1,
	}
	repo.load()
	return repo
}

func (r *MemoryCharacterRepository) load() {
	if r.filePath == "" {
		return
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}
	var store charStore
	if err := json.Unmarshal(data, &store); err == nil {
		if store.NextID > 0 {
			r.nextID = store.NextID
		}
		for _, ch := range store.Characters {
			r.characters[ch.ID] = ch
			if ch.ID >= r.nextID {
				r.nextID = ch.ID + 1
			}
		}
	}
}

func (r *MemoryCharacterRepository) save() {
	if r.filePath == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(r.filePath), 0755)
	var list []*Character
	for _, ch := range r.characters {
		list = append(list, ch)
	}
	store := charStore{
		NextID:     r.nextID,
		Characters: list,
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(r.filePath, data, 0644)
	}
}

func (r *MemoryCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	var list []*Character
	for _, ch := range r.characters {
		if ch.AccountID == accountID {
			copied := *ch
			list = append(list, &copied)
		}
	}
	return list, nil
}

func (r *MemoryCharacterRepository) GetCharacterByName(name string) (*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	target := strings.ToLower(strings.TrimSpace(name))
	for _, ch := range r.characters {
		if strings.ToLower(ch.Name) == target {
			copied := *ch
			return &copied, nil
		}
	}
	return nil, ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) GetCharacterByID(id int) (*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	ch, exists := r.characters[id]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	copied := *ch
	return &copied, nil
}

func (r *MemoryCharacterRepository) CreateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	if ch.Slot > 5 {
		return ErrSlotOutOfRange
	}

	targetName := strings.ToLower(strings.TrimSpace(ch.Name))
	for _, existing := range r.characters {
		if strings.ToLower(existing.Name) == targetName {
			return ErrCharacterExists
		}
		if existing.AccountID == ch.AccountID && existing.Slot == ch.Slot {
			return ErrSlotInUse
		}
	}

	ch.ID = r.nextID
	r.nextID++

	if ch.TownName == "" {
		ch.TownName = "Roumen"
	}
	if ch.Level == 0 {
		ch.Level = 1
	}

	// Deep copy
	stored := *ch
	if ch.Equips != nil {
		stored.Equips = make(map[byte]uint16)
		for k, v := range ch.Equips {
			stored.Equips[k] = v
		}
	}

	r.characters[ch.ID] = &stored
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) UpdateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	if _, exists := r.characters[ch.ID]; !exists {
		return ErrCharacterNotFound
	}

	stored := *ch
	r.characters[ch.ID] = &stored
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacter(accountID int, slot byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	for id, ch := range r.characters {
		if ch.AccountID == accountID && ch.Slot == slot {
			delete(r.characters, id)
			r.save()
			return nil
		}
	}
	return ErrCharacterNotFound
}
