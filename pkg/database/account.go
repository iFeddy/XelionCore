package database

import "time"

// TUser represents the official tUser table from Account.bak.
type TUser struct {
	UserNo     int       `json:"nUserNo"`
	UserID     string    `json:"sUserID"`
	UserPW     string    `json:"sUserPW"`
	UserName   string    `json:"sUserName"`
	IsBlock    int       `json:"bIsBlock"`  // 0 = Active, 1 = Blocked
	IsDelete   int       `json:"bIsDelete"` // 0 = Normal, 1 = Deleted
	AuthID     byte      `json:"nAuthID"`   // 1 = Player, 100 = Admin/GM
	UserIP     string    `json:"sUserIP"`
	CreateDate time.Time `json:"dCreateDate"`
	LoginDate  time.Time `json:"dLoginDate"`
}

// TUserAuth represents the privilege tier from tUserAuth in Account.bak.
type TUserAuth struct {
	AuthID      byte   `json:"nAuthID"`
	IsLoginAble int    `json:"bIsLoginAble"`
	AuthName    string `json:"sAuthName"`
}

// TUserBlock tracks block reasons matching tUserBlock in Account.bak.
type TUserBlock struct {
	UserNo    int       `json:"nUserNo"`
	Reason    string    `json:"sReason"`
	BlockDate time.Time `json:"dBlockDate"`
}

// Account represents a user account in XelionCore, backed by official tUser schema.
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    byte   `json:"admin"`
	Banned   bool   `json:"banned"`

	// Relational record matching tUser
	TUser TUser `json:"t_user"`
}

// Sync ensures bidirectional alignment between flat fields and TUser record.
func (a *Account) Sync() {
	if a.TUser.UserNo == 0 && a.ID != 0 {
		a.TUser.UserNo = a.ID
		a.TUser.UserID = a.Username
		a.TUser.UserPW = a.Password
		a.TUser.UserName = a.Username
		if a.Admin > 0 {
			a.TUser.AuthID = 100
		} else {
			a.TUser.AuthID = 1
		}
		if a.Banned {
			a.TUser.IsBlock = 1
		} else {
			a.TUser.IsBlock = 0
		}
		if a.TUser.CreateDate.IsZero() {
			a.TUser.CreateDate = time.Now()
		}
	} else if a.TUser.UserNo != 0 {
		a.ID = a.TUser.UserNo
		a.Username = a.TUser.UserID
		a.Password = a.TUser.UserPW
		if a.TUser.AuthID >= 100 {
			a.Admin = 2
		} else {
			a.Admin = 0
		}
		a.Banned = (a.TUser.IsBlock != 0)
	}
}
