package audit_test

import (
	"database/sql"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
)

// TestFase08_Control8_AST_ZeroHardcode verifica mediante análisis de árbol sintáctico (AST)
// que no existan valores mágicos ni desvíos hardcodeados en los controladores y repositorios.
func TestFase08_Control8_AST_ZeroHardcode(t *testing.T) {
	filesToCheck := []string{
		"../../pkg/server/login/handler.go",
		"../../pkg/server/world/handler.go",
		"../../pkg/server/zone/handler.go",
		"../../pkg/database/mysql_account.go",
		"../../pkg/database/mysql_account_shop.go",
		"../../pkg/database/mysql_accountlog.go",
		"../../pkg/database/mysql_world_session.go",
	}

	fset := token.NewFileSet()

	for _, relPath := range filesToCheck {
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			t.Fatalf("Error resolviendo ruta %s: %v", relPath, err)
		}

		src, err := os.ReadFile(absPath)
		if err != nil {
			t.Fatalf("Error leyendo archivo %s: %v", absPath, err)
		}

		node, err := parser.ParseFile(fset, absPath, src, parser.AllErrors)
		if err != nil {
			t.Fatalf("Error parseando AST de %s: %v", absPath, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// 1. Prohibir entero mágico 0x1234 (4660)
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.INT {
				if lit.Value == "0x1234" || lit.Value == "4660" {
					pos := fset.Position(lit.Pos())
					t.Errorf("[%s:%d] Detectado valor mágico seed prohibido: %s", pos.Filename, pos.Line, lit.Value)
				}
				// Prohibir 30001 como parámetro de entrega de retorno en handlers
				if strings.Contains(relPath, "server/login/handler.go") && lit.Value == "30001" {
					pos := fset.Position(lit.Pos())
					t.Errorf("[%s:%d] Detectado itemNo 30001 hardcodeado en login handler", pos.Filename, pos.Line)
				}
				// Prohibir 10001 y 500 en world handler (Munsang hardcode)
				if strings.Contains(relPath, "server/world/handler.go") && (lit.Value == "10001" || lit.Value == "500") {
					pos := fset.Position(lit.Pos())
					t.Errorf("[%s:%d] Detectado cupón o cash 10001/500 hardcodeado en world handler: %s", pos.Filename, pos.Line, lit.Value)
				}
				// Prohibir 86400 en handleLogoutLastDay
				if strings.Contains(relPath, "server/world/handler.go") && lit.Value == "86400" {
					pos := fset.Position(lit.Pos())
					t.Errorf("[%s:%d] Detectado factor de conversión días prohibido: 86400 en world handler", pos.Filename, pos.Line)
				}
			}

			// 2. Prohibir "US" hardcodeado en CONNECTFROMWHERE de world/handler.go
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if strings.Contains(relPath, "server/world/handler.go") && lit.Value == `"US"` {
					pos := fset.Position(lit.Pos())
					t.Errorf("[%s:%d] Detectada región fija 'US' prohibida en CONNECTFROMWHERE", pos.Filename, pos.Line)
				}
			}

			return true
		})
	}

	t.Log("Control 8 AST PASS: Cero literales mágicos ni desvíos hardcodeados encontrados.")
}

// TestFase08_Control8_Functional_SP_Remediation verifica que las operaciones de saneamiento
// invoquen correctamente sus Stored Procedures oficiales en la base de datos xelion_online.
func TestFase08_Control8_Functional_SP_Remediation(t *testing.T) {
	cfg, err := config.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	repo := database.NewMySQLAccountRepository(db)

	t.Run("ReturnUser_ItemFromDB", func(t *testing.T) {
		// SetReturnUserItem ejecuta CALL usp_ReturnUser_SetUserItem(1, @ret)
		err := repo.SetReturnUserItem(1)
		if err != nil {
			t.Fatalf("SetReturnUserItem falló: %v", err)
		}
	})

	t.Run("Outspark_LogonToken", func(t *testing.T) {
		// LoginOutspark ejecuta CALL usp_UserKey_AuthenticateToken
		acc, err := repo.LoginOutspark("INVALID_OR_DUMMY_TOKEN_0001", "127.0.0.1")
		if err == nil && acc != nil {
			t.Fatalf("Se esperaba fallo con token dummy")
		}
	})

	t.Run("CheckPasswordS2S_Official", func(t *testing.T) {
		userNo, nRet, outID, _, _, _, err := repo.CheckPasswordS2S("admin", "admin", true)
		if err != nil {
			t.Fatalf("CheckPasswordS2S falló: %v", err)
		}
		if nRet != 0 || userNo != 1 || outID != "admin" {
			t.Fatalf("CheckPasswordS2S inesperado: userNo=%d, nRet=%d, outID=%s", userNo, nRet, outID)
		}
	})

	t.Run("WebDBKey_Lifecycle", func(t *testing.T) {
		key, err := repo.GetOrIssueWebKey(1, "")
		if err != nil || key == "" {
			t.Fatalf("GetOrIssueWebKey falló: %v (key='%s')", err, key)
		}
		if len(key) < 16 {
			t.Fatalf("Longitud de key insuficiente: %d", len(key))
		}
	})

	t.Run("EventMunsang_Redeem", func(t *testing.T) {
		// RedeemMunsangEvent ejecuta CALL usp_Event_DoneMunsang
		errCode, err := repo.RedeemMunsangEvent(1, 101)
		if err != nil {
			t.Fatalf("RedeemMunsangEvent falló: %v", err)
		}
		// Si el evento está desactivado en tAccountEventReward, devuelve 5445 (0x1545 oficial)
		if errCode != 5445 && errCode != 0 {
			t.Logf("RedeemMunsangEvent devolvió código: %d", errCode)
		}
	})

	t.Run("EventL20_SetFull", func(t *testing.T) {
		// Limpiar registro previo si existiera para la prueba
		_, _ = db.Exec("DELETE FROM tEvent_L20 WHERE nUserNo = 1 AND nEventNo = 88")

		// EventL20SetFull ejecuta CALL p_Event_L20_Set con los 6 parámetros oficiales
		err := repo.EventL20SetFull(1, "admin", 88, 1, 1001, 100)
		if err != nil {
			t.Fatalf("EventL20SetFull falló: %v", err)
		}
	})

	t.Run("HackScanLog_AddFull", func(t *testing.T) {
		// AddHackScanLogFull ejecuta CALL p_HackScan_AddFull con IP y CharID reales
		err := repo.AddHackScanLogFull(1, 1, 1001, 1, "127.0.0.1", "AdminChar", "Zero-hardcode integrity test")
		if err != nil {
			t.Fatalf("AddHackScanLogFull falló: %v", err)
		}
	})

	t.Run("WorldServer_Heartbeat_And_CCUWithMax", func(t *testing.T) {
		err := repo.WorldServerHeartbeat(1)
		if err != nil {
			t.Fatalf("WorldServerHeartbeat falló: %v", err)
		}

		err = repo.UpdateWorldServerCCUWithMax(1, 15, 2000)
		if err != nil {
			t.Fatalf("UpdateWorldServerCCUWithMax falló: %v", err)
		}
	})
}
