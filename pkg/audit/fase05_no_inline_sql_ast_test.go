package audit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
)

// TestFase05_Control3_AST_NoInlineSQL verifica mediante AST que:
// 1. En los repositorios de cuenta y sesión (mysql_account*.go, mysql_world_session.go),
//    así como en LoginService, el 100% de las invocaciones a base de datos utilizan Stored Procedures (CALL <SP>).
// 2. Ninguna consulta contiene SELECT, INSERT, UPDATE, DELETE o CREATE TABLE inline fuera de Stored Procedures.
// 3. No existen remanentes de persistencia en JSON o memoria (MemoryAccountRepository, accounts.json).
// 4. No existen funciones auxiliares obsoletas (UpdateLoginSession, initAuxiliaryTables).
// 5. El DSN de configuración incluye obligatoriamente parseTime=true, multiStatements=true y charset=utf8mb4.
// 6. Las interfaces de repositorio están completamente implementadas por los adaptadores MySQL.
func TestFase05_Control3_AST_NoInlineSQL(t *testing.T) {
	rootPath := findRepoRoot(t)

	// 1. Validar ausencia de archivos JSON y memoria
	jsonAccounts := filepath.Join(rootPath, "pkg", "server", "login", "data", "accounts.json")
	if _, err := os.Stat(jsonAccounts); !os.IsNotExist(err) {
		t.Errorf("CRITERIO FALLIDO: 'accounts.json' no debe existir en el repositorio: %s", jsonAccounts)
	}

	// 2. Validar DSN con parámetros requeridos
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load falló: %v", err)
	}
	dsn := cfg.DSN()
	if !strings.Contains(dsn, "parseTime=true") {
		t.Errorf("CRITERIO FALLIDO: DSN debe incluir 'parseTime=true', obtenido: %s", dsn)
	}
	if !strings.Contains(dsn, "multiStatements=true") {
		t.Errorf("CRITERIO FALLIDO: DSN debe incluir 'multiStatements=true', obtenido: %s", dsn)
	}
	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Errorf("CRITERIO FALLIDO: DSN debe incluir 'charset=utf8mb4', obtenido: %s", dsn)
	}

	// 3. Inspección AST de repositorios y LoginService
	filesToCheck := []string{
		filepath.Join(rootPath, "pkg", "database", "mysql_account.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_account_shop.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_accountlog.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_world_session.go"),
		filepath.Join(rootPath, "pkg", "server", "login", "handler.go"),
		filepath.Join(rootPath, "cmd", "LoginService", "main.go"),
	}

	forbiddenSymbols := []string{
		"MemoryAccountRepository",
		"UpdateLoginSession",
		"initAuxiliaryTables",
	}

	totalDBCalls := 0
	callCount := 0

	for _, file := range filesToCheck {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("Error al parsear archivo AST '%s': %v", file, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// A. Validar símbolos prohibidos en identificadores
			if ident, ok := n.(*ast.Ident); ok {
				for _, forbidden := range forbiddenSymbols {
					if ident.Name == forbidden {
						t.Errorf("SIMBOLO PROHIBIDO ENCONTRADO '%s' en %s:%d",
							forbidden, file, fset.Position(ident.Pos()).Line)
					}
				}
			}

			// B. Validar llamadas a base de datos (Query, QueryRow, Exec)
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					methodName := sel.Sel.Name
					if methodName == "Query" || methodName == "QueryRow" || methodName == "Exec" {
						totalDBCalls++
						if len(call.Args) > 0 {
							if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								queryVal := strings.Trim(lit.Value, "\"`'")
								queryUpper := strings.ToUpper(strings.TrimSpace(queryVal))

								if strings.HasPrefix(queryUpper, "CALL ") {
									callCount++
								} else {
									t.Errorf("CONSULTA SQL NO-CALL ENCONTRADA en %s:%d: %s",
										file, fset.Position(call.Pos()).Line, queryVal)
								}

								// Asegurar que no hay sentencias SQL inline desprotegidas
								if strings.HasPrefix(queryUpper, "SELECT ") ||
									strings.HasPrefix(queryUpper, "INSERT ") ||
									strings.HasPrefix(queryUpper, "UPDATE ") ||
									strings.HasPrefix(queryUpper, "DELETE ") ||
									strings.HasPrefix(queryUpper, "CREATE TABLE ") {
									t.Errorf("INLINE SQL PROHIBIDO en %s:%d: %s",
										file, fset.Position(call.Pos()).Line, queryVal)
								}
							}
						}
					}
				}
			}
			return true
		})
	}

	if totalDBCalls == 0 {
		t.Fatalf("No se detectaron llamadas a base de datos en los archivos analizados")
	}

	if callCount != totalDBCalls {
		t.Errorf("Discrepancia en llamadas: esperadas %d consultas con 'CALL ', encontradas %d",
			totalDBCalls, callCount)
	}

	t.Logf("Control 3 AST Exitoso: %d/%d llamadas a DB utilizan estrictamente Stored Procedures (CALL <SP>)",
		callCount, totalDBCalls)

	// 4. Validar interfaces tipadas en pkg/database
	var (
		_ database.AccountRepository      = (*database.MySQLAccountRepository)(nil)
		_ database.AccountShopRepository  = (*database.MySQLAccountRepository)(nil)
		_ database.AccountLogRepository   = (*database.MySQLAccountRepository)(nil)
		_ database.LoginSessionRepository = (*database.MySQLAccountRepository)(nil)
	)
}

func findRepoRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd falló: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("No se encontró la raíz del repositorio (go.mod)")
		}
		dir = parent
	}
}
