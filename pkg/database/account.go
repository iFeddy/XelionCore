package database

// Account represents a user account in XelionCore.
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    byte   `json:"admin"`
	Banned   bool   `json:"banned"`
}
