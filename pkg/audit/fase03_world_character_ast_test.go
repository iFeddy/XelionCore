package audit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
)

// TestFase03_Control3_AST_WorldCharacter verifica exhaustivamente el Control 3 de WorldService / Character (Fase 03):
// 1. Escaneo AST de pkg/database/mysql_character*.go, mysql_zone_topology.go y transfer.go.
// 2. 100% de las invocaciones a base de datos (Query, QueryRow, Exec, Context variants) usan estrictamente CALL <SP>.
// 3. Cero SQL inline (SELECT, INSERT, UPDATE, DELETE, CREATE TABLE).
// 4. Cero mapas en memoria para transferencias (TransferManager no contiene map en su struct).
// 5. Garantía de implementación de interfaces modulares tipadas por dominio.
func TestFase03_Control3_AST_WorldCharacter(t *testing.T) {
	rootPath := findRepoRoot(t)

	filesToCheck := []string{
		filepath.Join(rootPath, "pkg", "database", "mysql_character.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_items.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_progression.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_guild.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_social.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_zone_topology.go"),
		filepath.Join(rootPath, "pkg", "database", "transfer.go"),
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
			// A. Inspeccionar TransferManager para asegurar cero campos de tipo mapa
			if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == "TransferManager" {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if _, isMap := field.Type.(*ast.MapType); isMap {
							var fieldName string
							if len(field.Names) > 0 {
								fieldName = field.Names[0].Name
							}
							t.Errorf("CRITERIO FALLIDO: TransferManager contiene campo tipo MAP en memoria '%s' en %s:%d",
								fieldName, file, fset.Position(field.Pos()).Line)
						}
					}
				}
			}

			// B. Validar llamadas a métodos de ejecución de consultas SQL
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					methodName := sel.Sel.Name
					isDBMethod := methodName == "Query" || methodName == "QueryRow" || methodName == "Exec" ||
						methodName == "QueryContext" || methodName == "QueryRowContext" || methodName == "ExecContext"

					if isDBMethod {
						totalDBCalls++
						// En métodos de context, la query es el argumento 1; en directos es el argumento 0
						argIndex := 0
						if strings.HasSuffix(methodName, "Context") {
							argIndex = 1
						}

						if len(call.Args) > argIndex {
							if lit, ok := call.Args[argIndex].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								queryVal := strings.Trim(lit.Value, "\"`'")
								queryUpper := strings.ToUpper(strings.TrimSpace(queryVal))

								if strings.HasPrefix(queryUpper, "CALL ") {
									callCount++
								} else {
									t.Errorf("CONSULTA SQL NO-CALL ENCONTRADA en %s:%d: %s",
										file, fset.Position(call.Pos()).Line, queryVal)
								}

								// Asegurar que no hay sentencias SQL inline
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
		t.Fatalf("No se detectaron llamadas a base de datos en los archivos de repositorios analizados")
	}

	if callCount != totalDBCalls {
		t.Errorf("Discrepancia en llamadas: esperadas %d consultas con 'CALL ', encontradas %d",
			totalDBCalls, callCount)
	}

	t.Logf("Control 3 AST Exitoso: %d/%d llamadas a DB utilizan estrictamente Stored Procedures (CALL <SP>) y TransferManager es 100%% sin mapas en memoria",
		callCount, totalDBCalls)

	// C. Validar interfaces tipadas por dominio implementadas por MySQLCharacterRepository
	var (
		_ database.CharacterAvatarRepository      = (*database.MySQLCharacterRepository)(nil)
		_ database.CharacterItemRepository        = (*database.MySQLCharacterRepository)(nil)
		_ database.CharacterProgressionRepository = (*database.MySQLCharacterRepository)(nil)
		_ database.GuildRepository                = (*database.MySQLCharacterRepository)(nil)
		_ database.WorldSocialRepository          = (*database.MySQLCharacterRepository)(nil)
		_ database.ZoneTopologyRepository         = (*database.MySQLCharacterRepository)(nil)
		_ database.CharacterRepository            = (*database.MySQLCharacterRepository)(nil)
	)
}

// TestFase03_Functional_MySQL_ZoneTransfer_And_Topology valida funcionalmente contra xelion_online:
// 1. Emisión y consumo de ZoneTransfer vía Stored Procedures (usp_Session_IssueZoneTransfer / usp_Session_ConsumeZoneTransfer).
// 2. Protección anti-replay (un token consumido no puede reutilizarse).
// 3. Operaciones de topología de nodos de zona (usp_WorldZone_*).
// 4. Operaciones atómicas de TransferManager 100% sobre MySQL.
func TestFase03_Functional_MySQL_ZoneTransfer_And_Topology(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repo := database.NewMySQLCharacterRepository(db)

	const (
		testAccountID = 88888
		testUsername  = "AuditUser03"
		testCharName  = "AuditChar03"
		testCharNo    = 77777
		testRandomID  = 43210
		testWorldNo   = 1
		testZoneNo    = 88
	)

	// Limpieza previa defensiva
	_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE sCharID = ?", testCharName)
	_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nWorldNo = ? AND nZoneNo = ?", testWorldNo, testZoneNo)

	t.Run("ZoneTransfer_Issue_And_Consume_Flow", func(t *testing.T) {
		// 1. Issue ZoneTransfer
		err := repo.IssueZoneTransfer(
			uint8(testWorldNo), uint16(testZoneNo), testAccountID, testUsername,
			testCharNo, testCharName, uint16(testRandomID), 0,
			"127.0.0.1", "Eld", 1234, 6789,
		)
		if err != nil {
			t.Fatalf("IssueZoneTransfer falló: %v", err)
		}

		// 2. Consume ZoneTransfer (primer intento: debe ser exitoso)
		res, err := repo.ConsumeZoneTransfer(testCharName, uint16(testRandomID), "127.0.0.1")
		if err != nil {
			t.Fatalf("ConsumeZoneTransfer intento 1 falló: %v", err)
		}
		if res == nil {
			t.Fatalf("ConsumeZoneTransfer retornó nil")
		}
		if res.AccountID != testAccountID {
			t.Errorf("AccountID incorrecto: esperado %d, obtenido %d", testAccountID, res.AccountID)
		}
		if res.CharID != testCharNo {
			t.Errorf("CharID incorrecto: esperado %d, obtenido %d", testCharNo, res.CharID)
		}
		if res.CharName != testCharName {
			t.Errorf("CharName incorrecto: esperado %s, obtenido %s", testCharName, res.CharName)
		}

		// 3. Protección Anti-Replay: el segundo intento debe fallar porque la fila ya fue consumida
		replayRes, err := repo.ConsumeZoneTransfer(testCharName, uint16(testRandomID), "127.0.0.1")
		if err == nil && replayRes != nil {
			t.Errorf("CRITERIO FALLIDO: Protección anti-replay violada. Retorno esperado error de consumo, obtenido: %+v", replayRes)
		}
	})

	t.Run("ZoneTopology_NodeManagement", func(t *testing.T) {
		// 1. UpsertZoneNode
		err := repo.UpsertZoneNode(uint8(testWorldNo), uint16(testZoneNo), "127.0.0.1", 9955)
		if err != nil {
			t.Fatalf("UpsertZoneNode falló: %v", err)
		}

		// 2. HeartbeatZoneNode
		err = repo.HeartbeatZoneNode(uint8(testWorldNo), uint16(testZoneNo), 15)
		if err != nil {
			t.Fatalf("HeartbeatZoneNode falló: %v", err)
		}

		// 3. GetZoneNodes
		nodes, err := repo.GetZoneNodes(uint8(testWorldNo))
		if err != nil {
			t.Fatalf("GetZoneNodes falló: %v", err)
		}
		found := false
		for _, n := range nodes {
			if n.ZoneNo == uint16(testZoneNo) {
				found = true
				if n.ZoneIP != "127.0.0.1" || n.ZonePort != 9955 {
					t.Errorf("Endpoint de nodo inesperado: %s:%d", n.ZoneIP, n.ZonePort)
				}
				if n.CurrentPlayers != 15 {
					t.Errorf("CurrentPlayers inesperado: esperado 15, obtenido %d", n.CurrentPlayers)
				}
				if !n.IsOnline {
					t.Errorf("Nodo debería estar marcado como IsOnline = true")
				}
				break
			}
		}
		if !found {
			t.Errorf("Nodo testZoneNo=%d no encontrado en lista de zonas de WorldNo=%d", testZoneNo, testWorldNo)
		}

		// 4. ResolveMapEndpoint
		ip, port, isOnline, err := repo.ResolveMapEndpoint(uint8(testWorldNo), uint16(testZoneNo))
		if err != nil {
			t.Fatalf("ResolveMapEndpoint falló: %v", err)
		}
		if ip != "127.0.0.1" || port != 9955 || !isOnline {
			t.Errorf("ResolveMapEndpoint resultado erróneo: ip=%s, port=%d, isOnline=%v", ip, port, isOnline)
		}
	})

	t.Run("TransferManager_MySQL_Flow", func(t *testing.T) {
		tm := database.NewTransferManager(db)

		// 1. World Client Token
		token := tm.CreateTransfer(testAccountID, testUsername, 0, "127.0.0.1")
		if token == "" {
			t.Fatalf("CreateTransfer retornó token vacío")
		}

		ct, ok := tm.ConsumeTransfer(token, nil)
		if !ok || ct == nil {
			t.Fatalf("ConsumeTransfer intento 1 falló (ok=%v, ct=%v)", ok, ct)
		}
		if ct.AccountID != testAccountID {
			t.Errorf("AccountID inesperado: %d vs %d", ct.AccountID, testAccountID)
		}

		// Anti-replay World Token
		replayCT, replayOk := tm.ConsumeTransfer(token, nil)
		if replayOk || replayCT != nil {
			t.Errorf("CRITERIO FALLIDO: World token consumido dos veces (ok=%v, ct=%v)", replayOk, replayCT)
		}

		// 2. Zone Transfer
		tm.CreateZoneTransfer(testAccountID, testUsername, "AuditCharTM", 77778, 55555, 0, "127.0.0.1")

		zt, ztOk := tm.ConsumeZoneTransfer("AuditCharTM", 55555, nil)
		if !ztOk || zt == nil {
			t.Fatalf("ConsumeZoneTransfer intento 1 falló (ok=%v, zt=%v)", ztOk, zt)
		}
		if zt.CharID != 77778 || zt.CharName != "AuditCharTM" {
			t.Errorf("ZoneTransfer datos inconsistentes: %+v", zt)
		}

		// Anti-replay Zone Transfer
		replayZT, replayZTOk := tm.ConsumeZoneTransfer("AuditCharTM", 55555, nil)
		if replayZTOk || replayZT != nil {
			t.Errorf("CRITERIO FALLIDO: Zone transfer consumido dos veces (ok=%v, zt=%v)", replayZTOk, replayZT)
		}
	})

	// Limpieza final defensiva
	_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE sCharID IN (?, ?)", testCharName, "AuditCharTM")
	_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nWorldNo = ? AND nZoneNo = ?", testWorldNo, testZoneNo)
}
