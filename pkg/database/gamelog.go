package database

import (
	"context"
	"time"
)

// GameLogEntry representa una fila de auditoría de gameplay para tGameLog.
type GameLogEntry struct {
	Type         int
	CharNo       int
	Map          string
	MapX         int
	MapY         int
	MapZ         int
	TargetCharNo int
	TargetID     int
	ItemKey      int64
	Int1         int
	Int2         int
	Int3         int
	Bigint1      int64
}

// GameLogRecord representa una fila leída desde tGameLog por p_GameLog_Get.
type GameLogRecord struct {
	LogNo        int64
	Date         time.Time
	CharNo       int
	Type         int
	Map          string
	MapX         int
	MapY         int
	MapZ         int
	TargetCharNo int
	TargetID     int
	ItemKey      int64
	Int1         int
	Int2         int
	Int3         int
	Bigint1      int64
}

// GameLogCharRecord representa una fila leída desde tGameLog unida con tCharacter por p_GameLog_GetLog.
type GameLogCharRecord struct {
	LogNo        int64
	Date         time.Time
	CharNo       int
	Type         int
	Map          string
	MapX         int
	MapY         int
	TargetCharNo int
	TargetID     int
	ItemKey      int64
	Int1         int
	Int2         int
	Int3         int
	Bigint1      int64
	CharID       string
}

// GameLogRepository define el contrato oficial para interactuar con los 27 Stored Procedures de World00_GameLog.
// Regla mandatoria: CERO SQL inline, 100% CALL <SP>.
type GameLogRepository interface {
	// Inserciones (p_GameLog_Add y p_GameLog_Add_Type1..TypeL)
	AddLog(ctx context.Context, entry *GameLogEntry) (int64, error)
	AddType1(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z int) (int64, error)
	AddType2(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2 int) (int64, error)
	AddType3(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID, int1 int) (int64, error)
	AddType4(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2, int3 int) (int64, error)
	AddType5(ctx context.Context, nType int, nCharNo int, int1, int2 int) (int64, error)
	AddType6(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1 int, bigint1 int64) (int64, error)
	AddType7(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, int1, int2 int) (int64, error)
	AddType8(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1 int) (int64, error)
	AddType9(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z int, itemKey int64, int1, int2 int) (int64, error)
	AddTypeA(ctx context.Context, nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) (int64, error)
	AddTypeB(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) (int64, error)
	AddTypeC(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int) (int64, error)
	AddTypeD(ctx context.Context, nType int, nCharNo int, itemKey int64, int1, int2 int) (int64, error)
	AddTypeE(ctx context.Context, nType int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) (int64, error)
	AddTypeF(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) (int64, error)
	AddTypeG(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int, bigint1 int64) (int64, error)
	AddTypeH(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo int) (int64, error)
	AddTypeI(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int) (int64, error)
	AddTypeJ(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) (int64, error)
	AddTypeK(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2 int, bigint1 int64) (int64, error)
	AddTypeL(ctx context.Context, nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) (int64, error)

	// Consultas administrativas (World00_GameLog.bak)
	GetLogs(ctx context.Context, whereExpr string) ([]GameLogRecord, error)
	GetLogCount(ctx context.Context, whereExpr string) (int, error)
	GetLogsWithChar(ctx context.Context, queryExpr string) ([]GameLogCharRecord, error)
	GetLogsWithCharCount(ctx context.Context, queryExpr string) (int, error)

	// Exclusivo GameLog.exe (usp_Login_setChannel)
	SetLoginChannel(ctx context.Context, userNo int, channelNo int, channelID string) (int, error)
}
