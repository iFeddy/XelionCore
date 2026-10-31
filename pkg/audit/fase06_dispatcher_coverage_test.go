package audit_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net"
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
)

// TestFase06_Control5_DispatcherCoverage verifica que:
// 1. Los 37 puntos de protocolo oficiales no-OpTool de Login.pdb (Sección 4.1 de PLAN_ACCION)
//    y los handlers integrados de Account/AccountLog están explícitamente contemplados en el Dispatcher.
// 2. Mediante inspección AST de pkg/server/login/handler.go, se comprueba que ningún OpCode
//    caiga en la rama default / unhandled.
// 3. Ejecución en runtime de HandlePacket con cliente de prueba para validar que todos los casos
//    se despachan limpiamente y sin panics.
func TestFase06_Control5_DispatcherCoverage(t *testing.T) {
	rootPath := findRepoRoot(t)
	handlerFile := filepath.Join(rootPath, "pkg", "server", "login", "handler.go")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, handlerFile, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("Error al parsear AST de '%s': %v", handlerFile, err)
	}

	// 1. Catálogo obligatorio de OpCodes no-OpTool de Login.pdb y Account/AccountLog
	expectedOpCodes := map[uint16]string{
		// Topología y S2S (Login.pdb)
		protocol.PROTO_NC_MISC_S2SCONNECTION_RDY:         "PROTO_NC_MISC_S2SCONNECTION_RDY",         // 0x0801
		protocol.PROTO_NC_MISC_S2SCONNECTION_REQ:         "PROTO_NC_MISC_S2SCONNECTION_REQ",         // 0x0802
		protocol.PROTO_NC_MISC_S2SCONNECTION_ACK:         "PROTO_NC_MISC_S2SCONNECTION_ACK",         // 0x0803
		protocol.PROTO_NC_MISC_HEARTBEAT_REQ:             "PROTO_NC_MISC_HEARTBEAT_REQ",             // 0x0804
		protocol.PROTO_NC_MISC_HEARTBEAT_ACK:             "PROTO_NC_MISC_HEARTBEAT_ACK",             // 0x0805
		protocol.PROTO_NC_MISC_USER_COUNT_CMD:            "PROTO_NC_MISC_USER_COUNT_CMD",            // 0x0808
		protocol.PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD: "PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD", // 0x080B
		protocol.PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD: "PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD", // 0x080C
		protocol.PROTO_NC_MISC_START_THE_BOOM_CMD:        "PROTO_NC_MISC_START_THE_BOOM_CMD",        // 0x0823

		// Protocolo de Usuario C2S y respuestas oficiales (Login.pdb)
		protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ: "PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ", // 0x0C65
		protocol.PROTO_NC_USER_XTRAP_REQ:                "PROTO_NC_USER_XTRAP_REQ",                // 0x0C04
		protocol.PROTO_NC_USER_LOGIN_REQ:                "PROTO_NC_USER_LOGIN_REQ",                // 0x0C06
		protocol.PROTO_NC_USER_US_LOGIN_REQ:             "PROTO_NC_USER_US_LOGIN_REQ",             // 0x0C5A
		protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ:       "PROTO_NC_USER_LOGIN_WITH_OTP_REQ",       // 0x0C37
		protocol.PROTO_NC_USER_LOGIN_NETMARBLE_REQ:       "PROTO_NC_USER_LOGIN_NETMARBLE_REQ",       // 0x0C1D
		protocol.PROTO_NC_USER_LOGIN_OUTSPARK_REQ:         "PROTO_NC_USER_LOGIN_OUTSPARK_REQ",         // 0x0C20
		protocol.PROTO_NC_USER_WORLD_STATUS_REQ:          "PROTO_NC_USER_WORLD_STATUS_REQ",          // 0x0C1B
		protocol.PROTO_NC_USER_WORLDSELECT_REQ:           "PROTO_NC_USER_WORLDSELECT_REQ",           // 0x0C0B
		protocol.PROTO_NC_USER_NORMALLOGOUT_CMD:          "PROTO_NC_USER_NORMALLOGOUT_CMD",          // 0x0C18

		// Handlers regionales y de seguridad S2S (Account.pdb / Login.pdb)
		protocol.PROTO_NC_USER_PASSWORD_CHECK_REQ:     "PROTO_NC_USER_PASSWORD_CHECK_REQ",     // 0x0C07
		protocol.PROTO_NC_USER_PASSWORD_CHECK_ACK:     "PROTO_NC_USER_PASSWORD_CHECK_ACK",     // 0x0C08
		protocol.PROTO_NC_USER_WILLLOGIN_ACK:          "PROTO_NC_USER_WILLLOGIN_ACK",          // 0x0C0E
		protocol.PROTO_NC_USER_IS_IP_BLOCK_REQ:        "PROTO_NC_USER_IS_IP_BLOCK_REQ",        // 0x0C2A
		protocol.PROTO_NC_USER_IS_IP_BLOCK_ACK:        "PROTO_NC_USER_IS_IP_BLOCK_ACK",        // 0x0C2B
		protocol.PROTO_NC_USER_TEENAGER_REQ:           "PROTO_NC_USER_TEENAGER_REQ",           // 0x0C26
		protocol.PROTO_NC_USER_TEENAGER_ACK:           "PROTO_NC_USER_TEENAGER_ACK",           // 0x0C27
		protocol.PROTO_NC_USER_TEENAGER_SET_CMD:       "PROTO_NC_USER_TEENAGER_SET_CMD",       // 0x0C28
		protocol.PROTO_NC_USER_CREATE_OTP_REQ:         "PROTO_NC_USER_CREATE_OTP_REQ",         // 0x0C35
		protocol.PROTO_NC_USER_GER_PASSWORD_CHECK_REQ: "PROTO_NC_USER_GER_PASSWORD_CHECK_REQ", // 0x0C39
		protocol.PROTO_NC_USER_GER_PASSWORD_CHECK_ACK: "PROTO_NC_USER_GER_PASSWORD_CHECK_ACK", // 0x0C3A
		protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_REQ:    "PROTO_NC_USER_GER_IS_IP_BLOCK_REQ",    // 0x0C3B
		protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_ACK:    "PROTO_NC_USER_GER_IS_IP_BLOCK_ACK",    // 0x0C3C
		protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_REQ:  "PROTO_NC_USER_TW_PASSWORD_CHECK_REQ",  // 0x0C3E
		protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_ACK:  "PROTO_NC_USER_TW_PASSWORD_CHECK_ACK",  // 0x0C3F
		protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_REQ:     "PROTO_NC_USER_TW_IS_IP_BLOCK_REQ",     // 0x0C40
		protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_ACK:     "PROTO_NC_USER_TW_IS_IP_BLOCK_ACK",     // 0x0C41
		protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_REQ:  "PROTO_NC_USER_JP_PASSWORD_CHECK_REQ",  // 0x0C4B
		protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_ACK:  "PROTO_NC_USER_JP_PASSWORD_CHECK_ACK",  // 0x0C4C
		protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_REQ:     "PROTO_NC_USER_JP_IS_IP_BLOCK_REQ",     // 0x0C4D
		protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_ACK:     "PROTO_NC_USER_JP_IS_IP_BLOCK_ACK",     // 0x0C4E
		protocol.PROTO_NC_USER_CH_PASSWORD_CHECK_REQ:  "PROTO_NC_USER_CH_PASSWORD_CHECK_REQ",  // 0x0C53
		protocol.PROTO_NC_USER_CH_PASSWORD_CHECK_ACK:  "PROTO_NC_USER_CH_PASSWORD_CHECK_ACK",  // 0x0C54
		protocol.PROTO_NC_USER_US_PASSWORD_CHECK_REQ:  "PROTO_NC_USER_US_PASSWORD_CHECK_REQ",  // 0x0C5B
		protocol.PROTO_NC_USER_US_PASSWORD_CHECK_ACK:  "PROTO_NC_USER_US_PASSWORD_CHECK_ACK",  // 0x0C5C
		protocol.PROTO_NC_USER_US_IS_IP_BLOCK_REQ:     "PROTO_NC_USER_US_IS_IP_BLOCK_REQ",     // 0x0C5D
		protocol.PROTO_NC_USER_US_IS_IP_BLOCK_ACK:     "PROTO_NC_USER_US_IS_IP_BLOCK_ACK",     // 0x0C5E
		protocol.PROTO_NC_USER_RETURN_CHECK_REQ:       "PROTO_NC_USER_RETURN_CHECK_REQ",       // 0x0C60
		protocol.PROTO_NC_USER_LOGOUT_LAST_DAY_REQ:    "PROTO_NC_USER_LOGOUT_LAST_DAY_REQ",    // 0x0C62
		protocol.PROTO_NC_USER_SET_RETURN_CMD:         "PROTO_NC_USER_SET_RETURN_CMD",         // 0x0C64
	}

	// 2. Extraer los identificadores presentes en los case del switch en HandlePacket
	switchedCases := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "HandlePacket" {
			return true
		}

		ast.Inspect(fn.Body, func(bodyNode ast.Node) bool {
			sw, ok := bodyNode.(*ast.SwitchStmt)
			if !ok {
				return true
			}

			for _, stmt := range sw.Body.List {
				caseClause, ok := stmt.(*ast.CaseClause)
				if !ok {
					continue
				}
				for _, expr := range caseClause.List {
					if sel, ok := expr.(*ast.SelectorExpr); ok {
						switchedCases[sel.Sel.Name] = true
					}
				}
			}
			return false
		})

		return false
	})

	missingCases := 0
	for opCode, name := range expectedOpCodes {
		if !switchedCases[name] {
			t.Errorf("OPCODE NO RUTEO EN DISPATCHER: 0x%04X (%s)", opCode, name)
			missingCases++
		}
	}

	if missingCases > 0 {
		t.Fatalf("Control 5 FALLIDO: %d OpCodes oficiales sin handler en LoginService", missingCases)
	}

	t.Logf("Control 5 AST: Verificados %d/%d OpCodes oficiales en el Dispatcher de LoginService",
		len(expectedOpCodes), len(expectedOpCodes))

	// 3. Prueba runtime con HandlePacket invocando cada OpCode sin pánico
	wm := login.NewWorldManager()
	h := login.NewHandler(nil, wm, nil)

	serverPipe, clientPipe := net.Pipe()
	defer serverPipe.Close()
	defer clientPipe.Close()

	// Drenar clientPipe en segundo plano para evitar bloqueos síncronos en net.Pipe
	go func() {
		buf := make([]byte, 4096)
		for {
			_, err := clientPipe.Read(buf)
			if err != nil {
				return
			}
		}
	}()

	client, err := network.NewClient(serverPipe, nil)
	if err != nil {
		t.Fatalf("Error al instanciar mock network.Client: %v", err)
	}

	for opCode, name := range expectedOpCodes {
		t.Run(fmt.Sprintf("OpCode_0x%04X_%s", opCode, name), func(t *testing.T) {
			pw := protocol.NewPacketWriter(opCode)
			// Llenar payload dummy suficiente para evitar slice out of bounds
			pw.Fill(320, 0)
			pkt := protocol.NewPacket(pw.Data())

			// Debería ejecutarse sin panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Panic al despachar OpCode 0x%04X (%s): %v", opCode, name, r)
					}
				}()
				h.HandlePacket(client, pkt)
			}()
		})
	}
}
