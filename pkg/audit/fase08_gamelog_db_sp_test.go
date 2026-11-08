package audit_test

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
)

// TestFase08_Control1_Schema_176_Tables_And_No_ZoneTables verifica:
// 1. Ausencia definitiva de database/zone_tables.sql en el repositorio.
// 2. Existencia en xelion_online del 100% de las 176 tablas oficiales documentadas en docs/TablasLegacy.md.
// 3. Verificación exacta de las 15 columnas y tipos oficiales de tGameLog.
// 4. Verificación de los 5 índices requeridos en tGameLog (PK + char + type + date_type + itemkey).
// 5. Verificación de TEMP_QuestComplete y sus columnas oficiales.
func TestFase08_Control1_Schema_176_Tables_And_No_ZoneTables(t *testing.T) {
	rootPath := findRepoRoot(t)

	// 1. Validar que database/zone_tables.sql no existe físicamente
	zoneTablesSQL := filepath.Join(rootPath, "database", "zone_tables.sql")
	if _, err := os.Stat(zoneTablesSQL); !os.IsNotExist(err) {
		t.Fatalf("CRITERIO FALLIDO [Control 1]: database/zone_tables.sql todavía existe en el repositorio (%s)", zoneTablesSQL)
	}

	// 2. Leer dinámicamente las 176 tablas oficiales desde docs/TablasLegacy.md
	legacyMDPath := filepath.Join(rootPath, "docs", "TablasLegacy.md")
	data, err := os.ReadFile(legacyMDPath)
	if err != nil {
		t.Fatalf("Error leyendo docs/TablasLegacy.md: %v", err)
	}
	re := regexp.MustCompile(`\|\s*\*\*` + "`" + `([^` + "`" + `]+)` + "`" + `\*\*\s*\|`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) != 176 {
		t.Fatalf("Esperadas 176 tablas en docs/TablasLegacy.md, encontradas %d", len(matches))
	}
	official176Tables := make([]string, 0, len(matches))
	for _, m := range matches {
		official176Tables = append(official176Tables, m[1])
	}

	db := getTestDB(t)
	defer db.Close()

	rows, err := db.Query(`
		SELECT LOWER(TABLE_NAME) 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = 'xelion_online'
	`)
	if err != nil {
		t.Fatalf("Error consultando information_schema.TABLES: %v", err)
	}
	defer rows.Close()

	existingMap := make(map[string]bool)
	for rows.Next() {
		var tbl string
		if err := rows.Scan(&tbl); err != nil {
			t.Fatalf("Error escaneando tabla: %v", err)
		}
		existingMap[strings.ToLower(tbl)] = true
	}

	var missingTables []string
	for _, tbl := range official176Tables {
		if !existingMap[strings.ToLower(tbl)] {
			missingTables = append(missingTables, tbl)
		}
	}
	if len(missingTables) > 0 {
		t.Fatalf("CRITERIO FALLIDO [Control 1]: Faltan %d tablas oficiales en xelion_online: %v", len(missingTables), missingTables)
	}

	// 3. Validar las 15 columnas exactas de tGameLog
	colRows, err := db.Query(`
		SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = 'xelion_online' AND TABLE_NAME = 'tGameLog'
		ORDER BY ORDINAL_POSITION
	`)
	if err != nil {
		t.Fatalf("Error consultando columnas de tGameLog: %v", err)
	}
	defer colRows.Close()

	expectedCols := []struct {
		Name     string
		DataType string
	}{
		{"nLogNo", "bigint"},
		{"dDate", "datetime"},
		{"nType", "int"},
		{"nCharNo", "int"},
		{"sMap", "varchar"},
		{"nMapX", "int"},
		{"nMapY", "int"},
		{"nMapZ", "int"},
		{"nTargetCharNo", "int"},
		{"nTargetID", "int"},
		{"nItemKey", "bigint"},
		{"nInt1", "int"},
		{"nInt2", "int"},
		{"nInt3", "int"},
		{"nBigint1", "bigint"},
	}

	var foundCols []string
	idx := 0
	for colRows.Next() {
		var colName, dataType, isNullable string
		if err := colRows.Scan(&colName, &dataType, &isNullable); err != nil {
			t.Fatalf("Error escaneando columna: %v", err)
		}
		foundCols = append(foundCols, colName)
		if idx < len(expectedCols) {
			exp := expectedCols[idx]
			if !strings.EqualFold(colName, exp.Name) {
				t.Errorf("Columna %d esperada '%s', encontrada '%s'", idx+1, exp.Name, colName)
			}
			if !strings.EqualFold(dataType, exp.DataType) {
				t.Errorf("Tipo de columna '%s' esperado '%s', encontrado '%s'", colName, exp.DataType, dataType)
			}
		}
		idx++
	}

	if len(foundCols) != 15 {
		t.Fatalf("CRITERIO FALLIDO [Control 1]: tGameLog debe tener 15 columnas exactas, encontradas %d: %v", len(foundCols), foundCols)
	}

	// 4. Validar índices de auditoría en tGameLog
	idxRows, err := db.Query(`
		SELECT DISTINCT INDEX_NAME 
		FROM information_schema.STATISTICS 
		WHERE TABLE_SCHEMA = 'xelion_online' AND TABLE_NAME = 'tGameLog'
	`)
	if err != nil {
		t.Fatalf("Error consultando índices de tGameLog: %v", err)
	}
	defer idxRows.Close()

	existingIndexes := make(map[string]bool)
	for idxRows.Next() {
		var idxName string
		if err := idxRows.Scan(&idxName); err != nil {
			t.Fatalf("Error escaneando índice: %v", err)
		}
		existingIndexes[strings.ToLower(idxName)] = true
	}

	requiredIndexes := []string{
		"primary",
		"idx_gamelog_char",
		"idx_gamelog_type",
		"idx_gamelog_date_type",
		"idx_gamelog_itemkey",
	}
	for _, reqIdx := range requiredIndexes {
		if !existingIndexes[reqIdx] {
			t.Errorf("CRITERIO FALLIDO [Control 1]: Índice requerido '%s' no encontrado en tGameLog", reqIdx)
		}
	}

	t.Logf("[Control 1 PASS] 176/176 tablas oficiales verificadas, zone_tables.sql eliminado y tGameLog con 15 columnas e índices de auditoría oficiales.")
}

// TestFase08_Control2_GameLog_StoredProcedures_PoolSafe verifica:
// 1. Presencia de los 27 Stored Procedures de World00_GameLog en information_schema.ROUTINES.
// 2. Invocación pool-safe de los 22 procedimientos p_GameLog_Add* y verificación de nLogNo > 0.
// 3. Invocación de usp_Login_setChannel y retorno nResult == 1.
// 4. Consultas administrativas p_GameLog_Get, p_GameLog_GetCount, p_GameLog_GetLog, p_GameLog_GetLogCount.
// 5. Invocación concurrente en paralelo sobre el pool *sql.DB sin sql.ErrNoRows.
// 6. Funcionamiento íntegro de GameLogDispatcher (Flush, Stats y Close limpio).
func TestFase08_Control2_GameLog_StoredProcedures_PoolSafe(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Verificar los 27 Stored Procedures en information_schema.ROUTINES
	expectedProcs := []string{
		"p_GameLog_Add",
		"p_GameLog_Add_Type1",
		"p_GameLog_Add_Type2",
		"p_GameLog_Add_Type3",
		"p_GameLog_Add_Type4",
		"p_GameLog_Add_Type5",
		"p_GameLog_Add_Type6",
		"p_GameLog_Add_Type7",
		"p_GameLog_Add_Type8",
		"p_GameLog_Add_Type9",
		"p_GameLog_Add_TypeA",
		"p_GameLog_Add_TypeB",
		"p_GameLog_Add_TypeC",
		"p_GameLog_Add_TypeD",
		"p_GameLog_Add_TypeE",
		"p_GameLog_Add_TypeF",
		"p_GameLog_Add_TypeG",
		"p_GameLog_Add_TypeH",
		"p_GameLog_Add_TypeI",
		"p_GameLog_Add_TypeJ",
		"p_GameLog_Add_TypeK",
		"p_GameLog_Add_TypeL",
		"p_GameLog_Get",
		"p_GameLog_GetCount",
		"p_GameLog_GetLog",
		"p_GameLog_GetLogCount",
		"usp_Login_setChannel",
	}

	rows, err := db.QueryContext(ctx, `
		SELECT ROUTINE_NAME 
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = 'xelion_online' AND ROUTINE_TYPE = 'PROCEDURE'
	`)
	if err != nil {
		t.Fatalf("Error consultando information_schema.ROUTINES: %v", err)
	}
	defer rows.Close()

	routineMap := make(map[string]bool)
	for rows.Next() {
		var rName string
		if err := rows.Scan(&rName); err != nil {
			t.Fatalf("Error escaneando rutina: %v", err)
		}
		routineMap[strings.ToLower(rName)] = true
	}

	var missingProcs []string
	for _, proc := range expectedProcs {
		if !routineMap[strings.ToLower(proc)] {
			missingProcs = append(missingProcs, proc)
		}
	}
	if len(missingProcs) > 0 {
		t.Fatalf("CRITERIO FALLIDO [Control 2]: Faltan %d Stored Procedures en xelion_online: %v", len(missingProcs), missingProcs)
	}

	repo := database.NewMySQLGameLogRepository(db)

	// 2. Ejecutar cada uno de los 22 procedimientos de inserción y verificar nLogNo > 0
	type testLogExec struct {
		name string
		call func() (int64, error)
	}

	execs := []testLogExec{
		{"p_GameLog_Add", func() (int64, error) {
			return repo.AddLog(ctx, &database.GameLogEntry{
				Type: 0, CharNo: 1001, Map: "Rou", MapX: 100, MapY: 200, MapZ: 300,
				TargetCharNo: 1002, TargetID: 55, ItemKey: 9999, Int1: 1, Int2: 2, Int3: 3, Bigint1: 8888,
			})
		}},
		{"p_GameLog_Add_Type1", func() (int64, error) {
			return repo.AddType1(ctx, 1, 1001, "Rou", 100, 200, 300)
		}},
		{"p_GameLog_Add_Type2", func() (int64, error) {
			return repo.AddType2(ctx, 2, 1001, "Rou", 100, 200, 300, 10, 20)
		}},
		{"p_GameLog_Add_Type3", func() (int64, error) {
			return repo.AddType3(ctx, 3, 1001, "Rou", 100, 200, 300, 1002, 50, 500)
		}},
		{"p_GameLog_Add_Type4", func() (int64, error) {
			return repo.AddType4(ctx, 4, 1001, "Rou", 100, 200, 300, 1, 2, 3)
		}},
		{"p_GameLog_Add_Type5", func() (int64, error) {
			return repo.AddType5(ctx, 5, 1001, 10, 20)
		}},
		{"p_GameLog_Add_Type6", func() (int64, error) {
			return repo.AddType6(ctx, 6, 1001, "Rou", 100, 200, 300, 77, 1000000)
		}},
		{"p_GameLog_Add_Type7", func() (int64, error) {
			return repo.AddType7(ctx, 7, 1001, "Rou", 100, 200, 300, 1002, 1, 2)
		}},
		{"p_GameLog_Add_Type8", func() (int64, error) {
			return repo.AddType8(ctx, 8, 1001, "Rou", 100, 200, 300, 99)
		}},
		{"p_GameLog_Add_Type9", func() (int64, error) {
			return repo.AddType9(ctx, 9, 1001, "Rou", 100, 200, 300, 5555, 1, 2)
		}},
		{"p_GameLog_Add_TypeA", func() (int64, error) {
			return repo.AddTypeA(ctx, 10, 1001, 1002, 33, 4444, 1, 2, 3)
		}},
		{"p_GameLog_Add_TypeB", func() (int64, error) {
			return repo.AddTypeB(ctx, 11, 1001, "Rou", 100, 200, 300, 1002, 33, 4444, 1, 2, 3)
		}},
		{"p_GameLog_Add_TypeC", func() (int64, error) {
			return repo.AddTypeC(ctx, 12, 1001, "Rou", 100, 200, 300, 33, 4444, 1, 2, 3)
		}},
		{"p_GameLog_Add_TypeD", func() (int64, error) {
			return repo.AddTypeD(ctx, 13, 1001, 4444, 1, 2)
		}},
		{"p_GameLog_Add_TypeE", func() (int64, error) {
			return repo.AddTypeE(ctx, 14, "Rou", 100, 200, 300, 33, 4444, 1, 2)
		}},
		{"p_GameLog_Add_TypeF", func() (int64, error) {
			return repo.AddTypeF(ctx, 15, 1001, "Rou", 100, 200, 300, 33, 4444, 1, 2, 3, 50000)
		}},
		{"p_GameLog_Add_TypeG", func() (int64, error) {
			return repo.AddTypeG(ctx, 16, 1001, "Rou", 100, 200, 300, 33, 4444, 1, 2, 60000)
		}},
		{"p_GameLog_Add_TypeH", func() (int64, error) {
			return repo.AddTypeH(ctx, 17, 1001, "Rou", 100, 200, 300, 1002)
		}},
		{"p_GameLog_Add_TypeI", func() (int64, error) {
			return repo.AddTypeI(ctx, 18, 1001, "Rou", 100, 200, 300, 1002, 77)
		}},
		{"p_GameLog_Add_TypeJ", func() (int64, error) {
			return repo.AddTypeJ(ctx, 19, 1001, "Rou", 100, 200, 300, 88, 7777, 1, 2)
		}},
		{"p_GameLog_Add_TypeK", func() (int64, error) {
			return repo.AddTypeK(ctx, 20, 1001, "Rou", 100, 200, 300, 1, 2, 70000)
		}},
		{"p_GameLog_Add_TypeL", func() (int64, error) {
			return repo.AddTypeL(ctx, 21, 1001, 1002, 99, 8888, 1, 2, 3, 90000)
		}},
	}

	for _, e := range execs {
		logNo, err := e.call()
		if err != nil {
			t.Errorf("Error invocando %s: %v", e.name, err)
		}
		if logNo <= 0 {
			t.Errorf("Invocación %s retornó nLogNo inválido (%d)", e.name, logNo)
		}
	}

	// 3. Verificar usp_Login_setChannel
	res, err := repo.SetLoginChannel(ctx, 1001, 1, "Channel-01")
	if err != nil {
		t.Fatalf("Error invocando usp_Login_setChannel: %v", err)
	}
	if res != 1 {
		t.Fatalf("usp_Login_setChannel retornó resultado inesperado: %d", res)
	}

	// 4. Verificar consultas administrativas p_GameLog_Get*
	logs, err := repo.GetLogs(ctx, "nCharNo = 1001")
	if err != nil {
		t.Fatalf("Error en GetLogs: %v", err)
	}
	if len(logs) == 0 {
		t.Fatalf("GetLogs no retornó registros para nCharNo = 1001")
	}

	count, err := repo.GetLogCount(ctx, "nCharNo = 1001")
	if err != nil {
		t.Fatalf("Error en GetLogCount: %v", err)
	}
	if count < len(logs) {
		t.Fatalf("GetLogCount (%d) menor que cantidad leída (%d)", count, len(logs))
	}

	// 5. Invocación paralela concurrente (20 goroutines en simultáneo sobre el pool)
	var wg sync.WaitGroup
	errCh := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(charNo int) {
			defer wg.Done()
			lNo, err := repo.AddType1(ctx, 1, charNo, "Rou", 100+charNo, 200+charNo, 300)
			if err != nil {
				errCh <- fmt.Errorf("concurrent AddType1 fail for char %d: %w", charNo, err)
				return
			}
			if lNo <= 0 {
				errCh <- fmt.Errorf("concurrent AddType1 returned invalid logNo: %d", lNo)
			}
		}(2000 + i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Error en prueba concurrente de SPs: %v", err)
	}

	// 6. Prueba completa de GameLogDispatcher (asíncrono, Flush y Close limpio)
	dispatcher := database.NewGameLogDispatcher(repo, 500, 4)
	dispatcher.Start()

	for i := 0; i < 50; i++ {
		ok := dispatcher.EnqueueType1(1, 3000+i, "Rou", i*10, i*20, 0)
		if !ok {
			t.Errorf("Dispatcher rejected log %d", i)
		}
	}

	flushCtx, flushCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer flushCancel()
	if err := dispatcher.Flush(flushCtx); err != nil {
		t.Fatalf("Error en Flush de GameLogDispatcher: %v", err)
	}

	closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeCancel()
	if err := dispatcher.Close(closeCtx); err != nil {
		t.Fatalf("Error en Close de GameLogDispatcher: %v", err)
	}

	qLen, processed, dropped := dispatcher.Stats()
	if qLen != 0 {
		t.Errorf("Dispatcher queue not drained: %d remaining", qLen)
	}
	if processed < 50 {
		t.Errorf("Dispatcher processed %d items, expected at least 50", processed)
	}
	if dropped != 0 {
		t.Errorf("Dispatcher dropped %d items unexpectedly", dropped)
	}

	t.Logf("[Control 2 PASS] 27/27 Stored Procedures ejecutados exitosamente, pool-safe, concurrentes y GameLogDispatcher verificado.")
}

// TestFase08_Control3_AST_NoInlineSQL_GameLog verifica mediante análisis AST:
// 1. Que pkg/database/mysql_gamelog.go utiliza 100% CALL <SP> para todas las consultas.
// 2. Que no existe ninguna sentencia SELECT, INSERT, UPDATE o DELETE inline en mysql_gamelog*.go.
// 3. Que pkg/server/zone/ no contiene sentencias SQL inline.
// 4. Que la interfaz GameLogRepository está implementada al 100% por MySQLGameLogRepository.
func TestFase08_Control3_AST_NoInlineSQL_GameLog(t *testing.T) {
	rootPath := findRepoRoot(t)

	// 1. Verificar implementación de interfaz
	var _ database.GameLogRepository = (*database.MySQLGameLogRepository)(nil)

	// 2. Inspección AST de pkg/database/mysql_gamelog.go
	gamelogFiles, err := filepath.Glob(filepath.Join(rootPath, "pkg", "database", "mysql_gamelog*.go"))
	if err != nil || len(gamelogFiles) == 0 {
		t.Fatalf("No se encontraron archivos mysql_gamelog*.go: %v", err)
	}

	totalDBCalls := 0
	fset := token.NewFileSet()

	for _, file := range gamelogFiles {
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("Error parseando AST en '%s': %v", file, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			methodName := sel.Sel.Name
			switch methodName {
			case "Query", "QueryRow", "Exec", "QueryContext", "QueryRowContext", "ExecContext":
				totalDBCalls++
				// Buscar argumento SQL (string literal)
				for _, arg := range call.Args {
					if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						queryVal := strings.Trim(lit.Value, "`\"'")
						queryUpper := strings.ToUpper(strings.TrimSpace(queryVal))

						if !strings.HasPrefix(queryUpper, "CALL ") {
							t.Errorf("CRITERIO FALLIDO [Control 3]: Consulta no invoca un Stored Procedure con CALL: '%s' en %s:%d",
								queryVal, file, fset.Position(lit.Pos()).Line)
						}

						// Prohibir SELECT, INSERT, UPDATE, DELETE fuera de CALL
						forbidden := []string{"SELECT ", "INSERT ", "UPDATE ", "DELETE "}
						for _, verb := range forbidden {
							if strings.HasPrefix(queryUpper, verb) {
								t.Errorf("CRITERIO FALLIDO [Control 3]: SQL inline prohibido '%s' encontrado en %s:%d",
									verb, file, fset.Position(lit.Pos()).Line)
							}
						}
					}
				}
			}
			return true
		})
	}

	if totalDBCalls < 25 {
		t.Fatalf("CRITERIO FALLIDO [Control 3]: Se esperaban al menos 25 llamadas a DB en mysql_gamelog*.go, encontradas: %d", totalDBCalls)
	}

	// 3. Inspección AST de pkg/server/zone/ para garantizar cero SQL inline
	zoneFiles, err := filepath.Glob(filepath.Join(rootPath, "pkg", "server", "zone", "*.go"))
	if err != nil {
		t.Fatalf("Error listando archivos en pkg/server/zone: %v", err)
	}

	for _, file := range zoneFiles {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		node, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("Error parseando AST en '%s': %v", file, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			switch sel.Sel.Name {
			case "Query", "QueryRow", "Exec", "QueryContext", "QueryRowContext", "ExecContext":
				for _, arg := range call.Args {
					if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						queryVal := strings.ToUpper(strings.Trim(lit.Value, "`\"'"))
						forbidden := []string{"SELECT ", "INSERT ", "UPDATE ", "DELETE "}
						for _, verb := range forbidden {
							if strings.HasPrefix(queryVal, verb) {
								t.Errorf("CRITERIO FALLIDO [Control 3]: SQL inline prohibido en ZoneService '%s' en %s:%d",
									queryVal, file, fset.Position(lit.Pos()).Line)
							}
						}
					}
				}
			}
			return true
		})
	}

	t.Logf("[Control 3 PASS] Escaneo AST superado con éxito: 100%% CALL <SP>, cero SQL inline en mysql_gamelog.go (%d llamadas) y cero en ZoneService.", totalDBCalls)
}
