package database

import (
	"context"
	"database/sql"
	"time"
)

// MySQLGameLogRepository implementa GameLogRepository ejecutando estrictamente CALL <SP>.
type MySQLGameLogRepository struct {
	db *sql.DB
}

// NewMySQLGameLogRepository construye una nueva instancia de MySQLGameLogRepository.
func NewMySQLGameLogRepository(db *sql.DB) *MySQLGameLogRepository {
	return &MySQLGameLogRepository{db: db}
}

// AddLog inserta un registro genérico completo en tGameLog invocando p_GameLog_Add.
func (r *MySQLGameLogRepository) AddLog(ctx context.Context, entry *GameLogEntry) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		entry.Type,
		entry.CharNo,
		entry.Map,
		entry.MapX,
		entry.MapY,
		entry.MapZ,
		entry.TargetCharNo,
		entry.TargetID,
		entry.ItemKey,
		entry.Int1,
		entry.Int2,
		entry.Int3,
		entry.Bigint1,
	).Scan(&logNo)
	return logNo, err
}

// AddType1 inserta login/logout de personaje y coordenadas en mapa (p_GameLog_Add_Type1).
func (r *MySQLGameLogRepository) AddType1(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type1(@out_nLogNo, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z,
	).Scan(&logNo)
	return logNo, err
}

// AddType2 inserta subida de nivel y progresión de estadísticas (p_GameLog_Add_Type2).
func (r *MySQLGameLogRepository) AddType2(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type2(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddType3 inserta transacciones de dinero y compra/venta con NPC (p_GameLog_Add_Type3).
func (r *MySQLGameLogRepository) AddType3(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID, int1 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type3(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetCharNo, targetID, int1,
	).Scan(&logNo)
	return logNo, err
}

// AddType4 inserta comercio directo entre jugadores o script Lua (p_GameLog_Add_Type4).
func (r *MySQLGameLogRepository) AddType4(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2, int3 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type4(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, int1, int2, int3,
	).Scan(&logNo)
	return logNo, err
}

// AddType5 inserta interacción compacta con almacén (p_GameLog_Add_Type5).
func (r *MySQLGameLogRepository) AddType5(ctx context.Context, nType int, nCharNo int, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type5(@out_nLogNo, ?, ?, ?, ?)",
		nType, nCharNo, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddType6 inserta obtención/descarte de ítems y montos Cen (p_GameLog_Add_Type6).
func (r *MySQLGameLogRepository) AddType6(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1 int, bigint1 int64) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type6(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, int1, bigint1,
	).Scan(&logNo)
	return logNo, err
}

// AddType7 inserta mejora, encantamiento o alteración de ítem (p_GameLog_Add_Type7).
func (r *MySQLGameLogRepository) AddType7(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type7(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetCharNo, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddType8 inserta muerte de personaje PvE/PvP (p_GameLog_Add_Type8).
func (r *MySQLGameLogRepository) AddType8(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type8(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, int1,
	).Scan(&logNo)
	return logNo, err
}

// AddType9 inserta inicio, progreso y entrega de misiones (p_GameLog_Add_Type9).
func (r *MySQLGameLogRepository) AddType9(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z int, itemKey int64, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_Type9(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, itemKey, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeA inserta eventos de gremio (p_GameLog_Add_TypeA).
func (r *MySQLGameLogRepository) AddTypeA(ctx context.Context, nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeA(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, targetCharNo, targetID, itemKey, int1, int2, int3,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeB inserta almacén de gremio y transferencias (p_GameLog_Add_TypeB).
func (r *MySQLGameLogRepository) AddTypeB(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeB(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetCharNo, targetID, itemKey, int1, int2, int3,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeC inserta quioscos o tiendas personales Booth (p_GameLog_Add_TypeC).
func (r *MySQLGameLogRepository) AddTypeC(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeC(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetID, itemKey, int1, int2, int3,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeD inserta producción, mezcla y desarme de ítems (p_GameLog_Add_TypeD).
func (r *MySQLGameLogRepository) AddTypeD(ctx context.Context, nType int, nCharNo int, itemKey int64, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeD(@out_nLogNo, ?, ?, ?, ?, ?)",
		nType, nCharNo, itemKey, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeE inserta aprendizaje y reinicio de habilidades (p_GameLog_Add_TypeE).
func (r *MySQLGameLogRepository) AddTypeE(ctx context.Context, nType int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeE(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, sMap, x, y, z, targetID, itemKey, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeF inserta operaciones en MiniHouse y mobiliario (p_GameLog_Add_TypeF).
func (r *MySQLGameLogRepository) AddTypeF(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeF(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetID, itemKey, int1, int2, int3, bigint1,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeG inserta eventos de mascotas y monturas (p_GameLog_Add_TypeG).
func (r *MySQLGameLogRepository) AddTypeG(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int, bigint1 int64) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeG(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetID, itemKey, int1, int2, bigint1,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeH inserta Kingdom Quest y eventos comunitarios (p_GameLog_Add_TypeH).
func (r *MySQLGameLogRepository) AddTypeH(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeH(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetCharNo,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeI inserta eventos de Casino, Dados y SlotMachine (p_GameLog_Add_TypeI).
func (r *MySQLGameLogRepository) AddTypeI(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeI(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetCharNo, targetID,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeJ inserta bodas, anillos y divorcios (p_GameLog_Add_TypeJ).
func (r *MySQLGameLogRepository) AddTypeJ(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeJ(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, targetID, itemKey, int1, int2,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeK inserta eventos de Academia de Gremio (p_GameLog_Add_TypeK).
func (r *MySQLGameLogRepository) AddTypeK(ctx context.Context, nType int, nCharNo int, sMap string, x, y, z, int1, int2 int, bigint1 int64) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeK(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, sMap, x, y, z, int1, int2, bigint1,
	).Scan(&logNo)
	return logNo, err
}

// AddTypeL inserta reportes, sanciones GM y recompensas (p_GameLog_Add_TypeL).
func (r *MySQLGameLogRepository) AddTypeL(ctx context.Context, nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) (int64, error) {
	var logNo int64
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_Add_TypeL(@out_nLogNo, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nType, nCharNo, targetCharNo, targetID, itemKey, int1, int2, int3, bigint1,
	).Scan(&logNo)
	return logNo, err
}

// GetLogs realiza una consulta administrativa sobre tGameLog (p_GameLog_Get).
func (r *MySQLGameLogRepository) GetLogs(ctx context.Context, whereExpr string) ([]GameLogRecord, error) {
	rows, err := r.db.QueryContext(ctx, "CALL p_GameLog_Get(?)", whereExpr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []GameLogRecord
	for rows.Next() {
		var rec GameLogRecord
		var dVal any
		if err := rows.Scan(
			&rec.LogNo,
			&dVal,
			&rec.CharNo,
			&rec.Type,
			&rec.Map,
			&rec.MapX,
			&rec.MapY,
			&rec.MapZ,
			&rec.TargetCharNo,
			&rec.TargetID,
			&rec.ItemKey,
			&rec.Int1,
			&rec.Int2,
			&rec.Int3,
			&rec.Bigint1,
		); err != nil {
			return nil, err
		}
		rec.Date = parseDateTime(dVal)
		records = append(records, rec)
	}
	return records, rows.Err()
}

// GetLogCount retorna el conteo de registros para un filtro en tGameLog (p_GameLog_GetCount).
func (r *MySQLGameLogRepository) GetLogCount(ctx context.Context, whereExpr string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_GetCount(?, @out_nCount)", whereExpr).Scan(&count)
	return count, err
}

// GetLogsWithChar realiza una consulta uniendo tGameLog y tCharacter (p_GameLog_GetLog).
func (r *MySQLGameLogRepository) GetLogsWithChar(ctx context.Context, queryExpr string) ([]GameLogCharRecord, error) {
	rows, err := r.db.QueryContext(ctx, "CALL p_GameLog_GetLog(?)", queryExpr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []GameLogCharRecord
	for rows.Next() {
		var rec GameLogCharRecord
		var dVal any
		if err := rows.Scan(
			&rec.LogNo,
			&dVal,
			&rec.CharNo,
			&rec.Type,
			&rec.Map,
			&rec.MapX,
			&rec.MapY,
			&rec.TargetCharNo,
			&rec.TargetID,
			&rec.ItemKey,
			&rec.Int1,
			&rec.Int2,
			&rec.Int3,
			&rec.Bigint1,
			&rec.CharID,
		); err != nil {
			return nil, err
		}
		rec.Date = parseDateTime(dVal)
		records = append(records, rec)
	}
	return records, rows.Err()
}

// GetLogsWithCharCount retorna el conteo de registros uniendo tGameLog y tCharacter (p_GameLog_GetLogCount).
func (r *MySQLGameLogRepository) GetLogsWithCharCount(ctx context.Context, queryExpr string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "CALL p_GameLog_GetLogCount(?, @out_nCount)", queryExpr).Scan(&count)
	return count, err
}

// SetLoginChannel registra el canal de conexión del usuario invocando usp_Login_setChannel.
func (r *MySQLGameLogRepository) SetLoginChannel(ctx context.Context, userNo int, channelNo int, channelID string) (int, error) {
	var result int
	err := r.db.QueryRowContext(ctx, "CALL usp_Login_setChannel(?, ?, ?, @out_res)", userNo, channelNo, channelID).Scan(&result)
	return result, err
}

// parseDateTime convierte el valor datetime escaneado de MySQL a time.Time.
func parseDateTime(val any) time.Time {
	switch v := val.(type) {
	case time.Time:
		return v
	case []byte:
		t, _ := time.Parse("2006-01-02 15:04:05", string(v))
		return t
	case string:
		t, _ := time.Parse("2006-01-02 15:04:05", v)
		return t
	default:
		return time.Time{}
	}
}

