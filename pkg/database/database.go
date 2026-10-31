package database

import (
	"errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountExists   = errors.New("account already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrAccountBlocked  = errors.New("account blocked")
)

// AccountRepository define las operaciones oficiales de autenticación, perfil y ciclo de vida de cuentas.
type AccountRepository interface {
	// Autenticación oficial multirregión
	LoginGameUS(userID, userPW, ip string) (*Account, error)
	LoginGame(userID, userPW, ip string) (*Account, error)
	LoginNetmarble(uniID, userID string) (*Account, error)
	CheckBlock(userNo int) (*BlockInfo, error)

	// Gestión de cuenta y credenciales
	GetAccountByUsername(username string) (*Account, error)
	CreateAccount(username, password string) (*Account, error)
	UpdatePassword(userNo int, newPassword string) error
	DeleteAccount(userNo int) error
	SetDrawComment(userNo int, comment string) error
	SetChannelLog(userNo int, channelID int) error

	// Auditoría y seguridad
	IsIPBlocked(userNo int, ip string) (bool, error)
	LogAccount(userID string, logType uint8, userNo int, playMin uint16, userPW string, worldNo uint8, errorCode uint16, ip string) error

	// Extensiones y subsistemas de cuenta
	CheckPcbangIP(ip string) (bool, error)
	GetWebKey(userNo int) (string, error)
	GetOrIssueWebKey(userNo int, candidateKey string) (string, error)
	GetPlayRemainTime(userNo int) (int, error)
	SetPlayRemainTime(userNo int, remainSec int) error
	GetUserDateDiffCreated(userNo int) (int64, error)
	GetReturnUserDateDiffLastLog(userNo int) (int64, error)
	SetReturnUserItem(userNo int, itemNo ...int) error
	GetTutorialIsSkip(userNo int) (bool, error)
	GetTutorialCharInfo(userNo int) (int, error)
	SetTutorialCharInfo(userNo int, step int) error
	LoginOutspark(token string, ip string) (*Account, error)
	CheckPasswordS2S(userID, userPW string, allowEmail bool) (int, int, string, byte, byte, byte, error)
	GetAccountByUserNo(userNo int) (*Account, error)
	CreateAccountWithIP(username, password, ip string) (*Account, error)
}

// AccountShopRepository define las operaciones de Cash, ItemShop y cupones.
type AccountShopRepository interface {
	GetCash(userNo int, cashType uint8) (int, error)
	AddCash(userNo int, cashType uint8, cash int, reference string) error
	UseCash(userNo int, cashType uint8, cash int, reference string) error
	AddCashEventL20(userNo int, cash int) error
	AddEventCashByUserID(userID string, cash int64) error
	RedeemMunsangEvent(userNo int, questID uint16) (uint16, error)
	SetSerialKey(userNo int, serialKey int64) error
	GetChargedItems(userNo int) ([]*ChargedItem, error)
	InsertChargedItem(item *ChargedItem) (int64, error)
	DrawChargedItem(orderNo int64) error
	DeleteChargedItem(orderNo int64) error
	GetItemSell(itemNo int) (*ItemSell, error)
}

// AccountLogRepository define las operaciones de auditoría, logs de sesión, reportes y telemetría.
type AccountLogRepository interface {
	IsIPBlocked(userNo int, ip string) (bool, error)
	LogAccount(userID string, logType uint8, userNo int, playMin uint16, userPW string, worldNo uint8, errorCode uint16, ip string) error
	LogLogout(worldNo uint8, userNo, charNo int, classVal, level uint8, exp int64, playTime int16) error
	GetDateDiffLastLogout(userNo int) (int64, error)
	AddReport(report *Report) (int, error)
	GetReports(reporterUserNo int, worldNo uint8) ([]*ReportSummary, error)
	DeleteReport(regNumber int) (int16, error)
	AddSameTimeChar(worldNo uint8, count int) (int, error)
	AddSameTimeCharClass(worldNo uint8, classVal uint8, level uint16, count int) error
	EventL20Get(userNo int) (int, error)
	EventL20Set(userNo int, charNo int, worldNo uint8, eventNo uint8, eventPoint int) error
	EventL20SetFull(userNo int, userID string, eventNo uint8, worldNo uint8, charNo int, eventPoint int) error
	EventiPadSetTicket(userNo int, playTime ...uint16) error
	AddHackScanLog(userNo int, worldNo uint8, charNo int, ip, charID, desc string) error
	AddHackScanLogFull(userNo int, worldNo uint8, charNo int, hackType int16, ip, charID, desc string) error
	AddSyslog(opCode uint16, userNo int, userID string, ip string, worldNo uint8, ccu int, errorCode uint16) error
}

// LoginSessionRepository define las operaciones de tokens de mundo, versión, OTP y topología de mundos.
type LoginSessionRepository interface {
	SyncClientVersionKey(key string) error
	GetClientVersionKey() (string, error)
	GetWorldServerList() ([]*WorldServerInfo, error)
	GetWorldServerByNo(worldNo uint8) (*WorldServerInfo, error)
	UpdateWorldServerCCU(worldNo uint8, ccu int) error
	UpdateWorldServerCCUWithMax(worldNo uint8, ccu int, maxUser int) error
	WorldServerHeartbeat(worldNo uint8) error
	SetMaintenanceAll(maintenance bool) error
	RegisterLoginSession(userNo int, userID, ip string) (bool, error)
	IssueWorldToken(userNo int, worldNo uint8, token [64]byte) error
	ConsumeWorldToken(worldNo uint8, token [64]byte, ip string) (int, string, error)
	CreateOTP(userNo int, otp [32]byte) error
	ConsumeOTP(otp [32]byte, ip string) (int, string, error)
	CloseSessionAndGetPlayMin(userNo int) (uint16, error)
}
