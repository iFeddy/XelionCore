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

// Report representa un reporte/denuncia en la tabla tReport (AccountLog).
type Report struct {
	RegNumber      int       `json:"nRegNumber"`
	Reporter       int       `json:"nReporter"`
	ReporterCharID string    `json:"sReporterCharID"`
	Harmer         string    `json:"sHarmer"`
	WorldNo        uint8     `json:"nWorldNo"`
	ReportType     string    `json:"sReportType"`
	Reason         string    `json:"sReason"`
	ChatLog        string    `json:"sChatLog"`
	Result         uint8     `json:"nResult"`
	Remark         string    `json:"sRemark"`
	ReportDate     time.Time `json:"dReportDate"`
}

// ReportSummary representa el resumen devuelto por p_Report_Get.
type ReportSummary struct {
	RegNumber int    `json:"nRegNumber"`
	Result    uint8  `json:"nResult"`
	Remark    string `json:"sRemark"`
}

// LogoutLog representa un registro en la tabla oficial tLogout (AccountLog).
type LogoutLog struct {
	RowNo      int64     `json:"rowNo"`
	LogDate    time.Time `json:"logDate"`
	WorldNo    uint8     `json:"worldNo"`
	UserNo     int       `json:"userNo"`
	CharNo     int       `json:"charNo"`
	Class      uint8     `json:"class"`
	Level      uint8     `json:"level"`
	Experience int64     `json:"experience"`
	PlayTime   int16     `json:"playTime"`
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
