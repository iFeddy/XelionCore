package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountExists   = errors.New("account already exists")
)

// AccountRepository defines operations for user authentication.
type AccountRepository interface {
	GetAccountByUsername(username string) (*Account, error)
	CreateAccount(username, password string) (*Account, error)
}

type accountStore struct {
	NextID   int        `json:"next_id"`
	Accounts []*Account `json:"accounts"`
}

// MemoryAccountRepository is a thread-safe store with JSON persistence for standalone and multi-process runs.
type MemoryAccountRepository struct {
	mu       sync.RWMutex
	filePath string
	accounts map[string]*Account
	nextID   int
}

// NewMemoryAccountRepository creates an account repository with JSON persistence at data/accounts.json.
func NewMemoryAccountRepository() *MemoryAccountRepository {
	filePath := filepath.Join("data", "accounts.json")
	repo := &MemoryAccountRepository{
		filePath: filePath,
		accounts: make(map[string]*Account),
		nextID:   1,
	}

	repo.load()

	// Ensure default admin account exists
	if _, exists := repo.accounts["admin"]; !exists {
		repo.accounts["admin"] = &Account{
			ID:       1,
			Username: "admin",
			Password: "admin",
			Admin:    2,
			Banned:   false,
		}
		if repo.nextID <= 1 {
			repo.nextID = 2
		}
		repo.save()
	}

	return repo
}

func (r *MemoryAccountRepository) load() {
	if r.filePath == "" {
		return
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}
	var store accountStore
	if err := json.Unmarshal(data, &store); err == nil {
		if store.NextID > 0 {
			r.nextID = store.NextID
		}
		for _, acc := range store.Accounts {
			r.accounts[strings.ToLower(acc.Username)] = acc
			if acc.ID >= r.nextID {
				r.nextID = acc.ID + 1
			}
		}
	}
}

func (r *MemoryAccountRepository) save() {
	if r.filePath == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(r.filePath), 0755)
	var list []*Account
	for _, acc := range r.accounts {
		list = append(list, acc)
	}
	store := accountStore{
		NextID:   r.nextID,
		Accounts: list,
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(r.filePath, data, 0644)
	}
}

func (r *MemoryAccountRepository) GetAccountByUsername(username string) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	acc, exists := r.accounts[strings.ToLower(username)]
	if !exists {
		return nil, ErrAccountNotFound
	}
	copyAcc := *acc
	return &copyAcc, nil
}

func (r *MemoryAccountRepository) CreateAccount(username, password string) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	key := strings.ToLower(username)
	if _, exists := r.accounts[key]; exists {
		return nil, ErrAccountExists
	}

	acc := &Account{
		ID:       r.nextID,
		Username: username,
		Password: password,
		Admin:    0,
		Banned:   false,
	}
	r.nextID++
	r.accounts[key] = acc
	r.save()

	copyAcc := *acc
	return &copyAcc, nil
}

// MySQLAccountRepository executes queries against a MySQL / MariaDB database.
type MySQLAccountRepository struct {
	db *sql.DB
}

// NewMySQLAccountRepository wraps an existing sql.DB.
func NewMySQLAccountRepository(db *sql.DB) *MySQLAccountRepository {
	return &MySQLAccountRepository{db: db}
}

func (r *MySQLAccountRepository) GetAccountByUsername(username string) (*Account, error) {
	row := r.db.QueryRow("SELECT id, username, password, admin, banned FROM accounts WHERE username = ?", username)
	acc := &Account{}
	var bannedInt int
	err := row.Scan(&acc.ID, &acc.Username, &acc.Password, &acc.Admin, &bannedInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	acc.Banned = (bannedInt != 0)
	return acc, nil
}

func (r *MySQLAccountRepository) CreateAccount(username, password string) (*Account, error) {
	res, err := r.db.Exec("INSERT INTO accounts (username, password, admin, banned) VALUES (?, ?, 0, 0)", username, password)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:       int(id),
		Username: username,
		Password: password,
		Admin:    0,
		Banned:   false,
	}, nil
}
