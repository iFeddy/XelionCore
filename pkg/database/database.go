package database

import (
	"database/sql"
	"errors"
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

// MemoryAccountRepository is a thread-safe in-memory store for testing without MySQL.
type MemoryAccountRepository struct {
	mu       sync.RWMutex
	accounts map[string]*Account
	nextID   int
}

// NewMemoryAccountRepository creates an in-memory repository with an optional pre-seeded admin.
func NewMemoryAccountRepository() *MemoryAccountRepository {
	repo := &MemoryAccountRepository{
		accounts: make(map[string]*Account),
		nextID:   1,
	}
	// Default admin account (password "admin" md5 is 21232f297a57a5a743894a0e4a801fc3)
	repo.accounts["admin"] = &Account{
		ID:       1,
		Username: "admin",
		Password: "admin",
		Admin:    2,
		Banned:   false,
	}
	repo.nextID = 2
	return repo
}

func (r *MemoryAccountRepository) GetAccountByUsername(username string) (*Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
