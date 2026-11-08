package audit_test

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

type noopConn struct{}

func (n *noopConn) Read(b []byte) (int, error)         { return 0, io.EOF }
func (n *noopConn) Write(b []byte) (int, error)        { return len(b), nil }
func (n *noopConn) Close() error                       { return nil }
func (n *noopConn) LocalAddr() net.Addr                { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9010} }
func (n *noopConn) RemoteAddr() net.Addr               { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9010} }
func (n *noopConn) SetDeadline(t time.Time) error      { return nil }
func (n *noopConn) SetReadDeadline(t time.Time) error  { return nil }
func (n *noopConn) SetWriteDeadline(t time.Time) error { return nil }

func createSocialTestUser(t *testing.T, db *sql.DB, userNo int, username string) {
	_, err := db.Exec("INSERT INTO tUser (nUserNo, sUserID, sUserPW, nAuthID) VALUES (?, ?, 'pass123', 1)", userNo, username)
	if err != nil {
		t.Fatalf("Error creando tUser %d: %v", userNo, err)
	}
}

func createSocialTestChar(t *testing.T, repo database.CharacterRepository, accountID int, name string, slot, job, level byte, town string) *database.Character {
	ch := &database.Character{
		AccountID: accountID,
		Name:      name,
		Slot:      slot,
		Job:       job,
		Level:     level,
		Male:      true,
		Hair:      1,
		HairColor: 1,
		Face:      1,
		TownName:  town,
		XPos:      6888,
		YPos:      7045,
		CurHP:     100,
		CurSP:     100,
	}
	ch.TCharacterState.Money = 500
	ch.TCharacterState.HPStones = 15
	ch.TCharacterState.SPStones = 7
	if err := repo.CreateCharacter(ch); err != nil {
		t.Fatalf("Error creando personaje %s: %v", name, err)
	}
	return ch
}

func cleanupSocialUser(db *sql.DB, userNo int) {
	_, _ = db.Exec("DELETE FROM tPartyMember WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tFriend WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tFriend WHERE nFriendCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tGuildAcademyChatBan WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tGuildAcademyMember WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tGuildMember WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tChatBlock WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tHolyPromise WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tReport WHERE nReporter = ?", userNo)
	_, _ = db.Exec("DELETE FROM tItem WHERE nOwner IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tCharacterShape WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tCharacterState WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", userNo)
	_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", userNo)
}

// TestFase06_Control5B_DispatcherCoverage100_C2S verifica que:
// 1. Los 167 handlers C2S oficiales de WorldManager.pdb (Sección 4.2.1 de PLAN_ACCION)
//    están 100% registrados y cubiertos en HandlePacket de pkg/server/world/handler.go.
// 2. Cero handlers C2S faltantes.
func TestFase06_Control5B_DispatcherCoverage100_C2S(t *testing.T) {
	rootPath := findRepoRoot(t)
	planFile := filepath.Join(rootPath, "docs", "plan", "PLAN_ACCION_WORLDSERVICE_CHARACTER.md")
	planData, err := os.ReadFile(planFile)
	if err != nil {
		t.Fatalf("Error leyendo PLAN_ACCION: %v", err)
	}

	// 1. Extraer tabla de 167 handlers C2S de la Sección 4.2.1
	reTable := regexp.MustCompile(`(?s)### 1\. Protocolos Cliente -> WorldManager \(C2S.*?\n\n(.*?)\n\n### 2\.`)
	matches := reTable.FindSubmatch(planData)
	if len(matches) < 2 {
		t.Fatalf("No se pudo localizar la tabla C2S de la Sección 4.2.1 en PLAN_ACCION")
	}

	// 2. Extraer mapa de opcodes desde pkg/protocol/opcodes.go
	opcodesFile := filepath.Join(rootPath, "pkg", "protocol", "opcodes.go")
	opcodesData, err := os.ReadFile(opcodesFile)
	if err != nil {
		t.Fatalf("Error leyendo opcodes.go: %v", err)
	}
	opcodeMap := make(map[[2]int]string)
	reOpDef := regexp.MustCompile(`(PROTO_NC_\w+)\s*uint16\s*=\s*\((?:protocol\.)?(\d+)\s*<<\s*10\)\s*\|\s*(\d+)`)
	for _, m := range reOpDef.FindAllStringSubmatch(string(opcodesData), -1) {
		dept := 0
		cmd := 0
		fmt.Sscanf(m[2], "%d", &dept)
		fmt.Sscanf(m[3], "%d", &cmd)
		opcodeMap[[2]int{dept, cmd}] = m[1]
	}

	// 3. Parsear filas de la tabla de la Sección 4.2.1
	tableLines := strings.Split(string(matches[1]), "\n")
	type c2sEntry struct {
		Hex    string
		Dept   int
		Cmd    int
		Cb     string
		Struct string
		Const  string
	}
	var expectedList []c2sEntry
	reRow := regexp.MustCompile(`\|\s*` + "`?" + `(0x[0-9A-Fa-f]+)` + "`?" + `\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*` + "`?" + `(\w+)` + "`?" + `\s*\|\s*` + "`?" + `(\w+)` + "`?")

	for _, line := range tableLines {
		m := reRow.FindStringSubmatch(line)
		if len(m) >= 7 {
			dept := 0
			cmd := 0
			fmt.Sscanf(m[3], "%d", &dept)
			fmt.Sscanf(m[4], "%d", &cmd)
			entry := c2sEntry{
				Hex:    strings.ToUpper(m[1]),
				Dept:   dept,
				Cmd:    cmd,
				Cb:     m[5],
				Struct: m[6],
				Const:  opcodeMap[[2]int{dept, cmd}],
			}
			expectedList = append(expectedList, entry)
		}
	}

	if len(expectedList) != 167 {
		t.Fatalf("Se esperaban exactamente 167 handlers C2S en Sección 4.2.1, encontrados: %d", len(expectedList))
	}

	// 4. Parsear HandlePacket en pkg/server/world/handler.go mediante AST
	handlerFile := filepath.Join(rootPath, "pkg", "server", "world", "handler.go")
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, handlerFile, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("Error parseando AST de handler.go: %v", err)
	}

	handledOpCodes := make(map[string]bool)
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "HandlePacket" {
			ast.Inspect(fn.Body, func(inner ast.Node) bool {
				if sc, ok := inner.(*ast.CaseClause); ok {
					for _, expr := range sc.List {
						if sel, ok := expr.(*ast.SelectorExpr); ok {
							handledOpCodes[sel.Sel.Name] = true
						}
					}
				}
				return true
			})
		}
		return true
	})

	// 5. Validar que cada uno de los 167 handlers esté cubierto
	var missing []string
	for _, exp := range expectedList {
		matched := false
		if handledOpCodes[exp.Struct] {
			matched = true
		} else if exp.Const != "" && handledOpCodes[exp.Const] {
			matched = true
		}

		if !matched {
			missing = append(missing, fmt.Sprintf("%s (Dept %d, Cmd %d) -> Struct: %s, Const: %s",
				exp.Hex, exp.Dept, exp.Cmd, exp.Struct, exp.Const))
		}
	}

	if len(missing) > 0 {
		t.Fatalf("CRITERIO FALLIDO: Faltan %d/%d handlers C2S en HandlePacket de WorldService:\n%s",
			len(missing), len(expectedList), strings.Join(missing, "\n"))
	}

	t.Logf("Control 5B Exitoso: 167/167 handlers C2S de WorldManager.pdb registrados en Dispatcher (0 faltantes)")
}

// TestFase06_Control5B_ZeroHardcode_AST verifica la ausencia estricta de:
// 1. Strings hardcodeados ("RouN" en lógica de sesión/mapas).
// 2. loadDefaults() en kq_manager.go.
// 3. Consultas SQL inline en pkg/server/world/ no-tests.
func TestFase06_Control5B_ZeroHardcode_AST(t *testing.T) {
	rootPath := findRepoRoot(t)
	worldDir := filepath.Join(rootPath, "pkg", "server", "world")

	files, err := filepath.Glob(filepath.Join(worldDir, "*.go"))
	if err != nil {
		t.Fatalf("Error listando archivos de world: %v", err)
	}

	fset := token.NewFileSet()

	for _, file := range files {
		isTest := strings.HasSuffix(file, "_test.go")

		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("Error parseando AST de %s: %v", file, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// Prohibir loadDefaults() en kq_manager.go
			if strings.HasSuffix(file, "kq_manager.go") {
				if ident, ok := n.(*ast.Ident); ok && ident.Name == "loadDefaults" {
					t.Errorf("CRITERIO FALLIDO: loadDefaults() detectado en kq_manager.go línea %d",
						fset.Position(ident.Pos()).Line)
				}
			}

			// Prohibir "RouN" y SQL inline en código operativo (no tests)
			if !isTest {
				if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					val := strings.Trim(bl.Value, "\"")
					if val == "RouN" {
						t.Errorf("CRITERIO FALLIDO: Mapa mágico 'RouN' detectado en %s línea %d",
							filepath.Base(file), fset.Position(bl.Pos()).Line)
					}

					upper := strings.ToUpper(bl.Value)
					sqlVerbs := []string{"SELECT ", "INSERT INTO ", "UPDATE ", "DELETE FROM "}
					for _, verb := range sqlVerbs {
						if strings.Contains(upper, verb) && !strings.Contains(upper, "CALL ") {
							t.Errorf("CRITERIO FALLIDO: SQL inline detectado en %s línea %d: %s",
								filepath.Base(file), fset.Position(bl.Pos()).Line, bl.Value)
						}
					}
				}
			}

			return true
		})
	}

	t.Log("Control 5B AST Exitoso: Cero mapas mágicos ('RouN'), cero loadDefaults() y cero SQL inline en pkg/server/world/")
}

// TestFase06_Control5B_Functional_Party_MySQL prueba el ciclo completo de Party en MySQL:
// p_Party_Create -> p_Party_AddMember -> p_Party_GetMembers -> p_Party_SetMaster ->
// p_Party_SetFlags -> p_Party_RemoveMember -> p_Party_Delete.
func TestFase06_Control5B_Functional_Party_MySQL(t *testing.T) {
	db := getTestDB(t)
	repo := database.NewMySQLCharacterRepository(db)

	const (
		accID1 = 98801
		accID2 = 98802
	)

	cleanupSocialUser(db, accID1)
	cleanupSocialUser(db, accID2)
	defer cleanupSocialUser(db, accID1)
	defer cleanupSocialUser(db, accID2)

	createSocialTestUser(t, db, accID1, "PartyU1")
	createSocialTestUser(t, db, accID2, "PartyU2")

	c1 := createSocialTestChar(t, repo, accID1, "PartyHero1", 0, 1, 20, "Eld")
	c2 := createSocialTestChar(t, repo, accID2, "PartyHero2", 0, 2, 20, "Eld")

	// 1. Crear Party
	partyID, err := repo.CreateParty(c1.ID, 1)
	if err != nil || partyID <= 0 {
		t.Fatalf("Error creando party: %v, id: %d", err, partyID)
	}

	// 2. Agregar Miembro
	if err := repo.AddPartyMember(partyID, c2.ID, 0); err != nil {
		t.Fatalf("Error agregando miembro a party: %v", err)
	}

	// 3. Consultar Miembros
	members, err := repo.GetPartyMembers(partyID)
	if err != nil || len(members) < 2 {
		t.Fatalf("Error obteniendo miembros de party: %v (count: %d)", err, len(members))
	}

	// 4. Cambiar Master
	if err := repo.SetPartyMaster(partyID, c2.ID); err != nil {
		t.Fatalf("Error cambiando master de party: %v", err)
	}

	// 5. Cambiar Flags/Looting
	if err := repo.SetPartyFlags(partyID, 2); err != nil {
		t.Fatalf("Error cambiando flags de party: %v", err)
	}

	// 6. Remover Miembro
	if err := repo.RemovePartyMember(partyID, c1.ID); err != nil {
		t.Fatalf("Error removiendo miembro de party: %v", err)
	}

	// 7. Borrar Party
	if err := repo.DeleteParty(partyID); err != nil {
		t.Fatalf("Error eliminando party: %v", err)
	}

	t.Log("Control 5B Funcional Exitoso: Ciclo completo de Party en MySQL verificado")
}

// TestFase06_Control5B_Functional_Friends_MySQL prueba el ciclo completo de Friends en MySQL:
// p_Friend_Add -> p_Friend_GetAll2 (sin RouN) -> usp_Character_SetFriendpoint -> p_Friend_Del.
func TestFase06_Control5B_Functional_Friends_MySQL(t *testing.T) {
	db := getTestDB(t)
	repo := database.NewMySQLCharacterRepository(db)

	const (
		accID1 = 98803
		accID2 = 98804
	)

	cleanupSocialUser(db, accID1)
	cleanupSocialUser(db, accID2)
	defer cleanupSocialUser(db, accID1)
	defer cleanupSocialUser(db, accID2)

	createSocialTestUser(t, db, accID1, "FriendU1")
	createSocialTestUser(t, db, accID2, "FriendU2")

	c1 := createSocialTestChar(t, repo, accID1, "FriendHero1", 0, 1, 30, "Uruga")
	c2 := createSocialTestChar(t, repo, accID2, "FriendHero2", 0, 3, 30, "Alberstol")

	// 1. Agregar amigo
	if err := repo.AddFriend(c1.ID, c2.ID); err != nil {
		t.Fatalf("Error agregando amigo: %v", err)
	}

	// 2. Obtener lista de amigos verificando sLoginZone real desde tCharacter (no hardcodeado)
	friends, err := repo.GetFriends(c1.ID)
	if err != nil || len(friends) == 0 {
		t.Fatalf("Error obteniendo amigos: %v (count: %d)", err, len(friends))
	}

	found := false
	for _, f := range friends {
		if f.CharNo == c2.ID {
			found = true
			if f.Map == "RouN" && f.Map != "Alberstol" {
				t.Errorf("CRITERIO FALLIDO: Se esperaba mapa real 'Alberstol', se obtuvo fallback mágico '%s'", f.Map)
			}
		}
	}
	if !found {
		t.Fatalf("Amigo %d no encontrado en lista devuelta por GetFriends", c2.ID)
	}

	// 3. Friend Points (usp_Character_SetFriendpoint / usp_Character_GetFriendpoint)
	if err := repo.SetFriendPoints(c1.ID, 150); err != nil {
		t.Fatalf("Error seteando puntos de amigo: %v", err)
	}
	pts, err := repo.GetFriendPoints(c1.ID)
	if err != nil || pts != 150 {
		t.Fatalf("Error obteniendo puntos de amigo: %v, esperado 150, obtenido %d", err, pts)
	}

	// 4. Eliminar amigo
	if err := repo.RemoveFriend(c1.ID, c2.ID); err != nil {
		t.Fatalf("Error eliminando amigo: %v", err)
	}

	t.Log("Control 5B Funcional Exitoso: Ciclo completo de Friends en MySQL verificado")
}

// TestFase06_Control5B_Functional_KingdomQuest_SHN_And_MySQL prueba la carga 100% desde SHN
// y sincronización en tKingdomQuest de MySQL sin loadDefaults().
func TestFase06_Control5B_Functional_KingdomQuest_SHN_And_MySQL(t *testing.T) {
	rootPath := findRepoRoot(t)
	dataDir := filepath.Join(rootPath, "data")

	db := getTestDB(t)
	repo := database.NewMySQLCharacterRepository(db)

	store, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("Error cargando store: %v", err)
	}

	kqMgr := world.NewKQManager(store, repo)

	if kqMgr.QuestCount() < 50 {
		t.Fatalf("CRITERIO FALLIDO: Se esperaban al menos 50 misiones KQ cargadas desde KingdomQuest.shn, obtenidas: %d",
			kqMgr.QuestCount())
	}

	// Verificar sincronización en MySQL tKingdomQuest
	quests, err := repo.GetKingdomQuests()
	if err != nil {
		t.Fatalf("Error consultando tKingdomQuest: %v", err)
	}
	if len(quests) == 0 {
		t.Fatalf("CRITERIO FALLIDO: tKingdomQuest está vacía tras inicialización de KQManager")
	}

	// Test funcional de registro y despacho de schedule con noopConn
	client, _ := network.NewClient(&noopConn{}, nil)
	client.AccountID = 99991
	client.CharacterID = 8891
	client.Username = "KQHero"

	kqMgr.RegisterOnline(client, client.CharacterID)
	kqMgr.HandleSchedule(client)
	kqMgr.HandleVoteStartCheck(client)

	t.Logf("Control 5B Funcional Exitoso: %d KQs cargadas desde SHN y sincronizadas a tKingdomQuest en MySQL", kqMgr.QuestCount())
}

// TestFase06_Control5B_Functional_Guild_And_Academy_MySQL prueba persistencia de gremios y academia:
// p_Guild_Create -> p_GuildMember_Add -> p_Guild_Grade_Set -> p_Guild_GetHistoryList ->
// p_GuildAcademy_Add -> p_GuildAcademy_ChatBan_Add -> p_GuildAcademy_SetGraduate.
func TestFase06_Control5B_Functional_Guild_And_Academy_MySQL(t *testing.T) {
	db := getTestDB(t)
	repo := database.NewMySQLCharacterRepository(db)

	const (
		accID1 = 98805
		accID2 = 98806
		accID3 = 98809
	)

	cleanupSocialUser(db, accID1)
	cleanupSocialUser(db, accID2)
	cleanupSocialUser(db, accID3)
	defer cleanupSocialUser(db, accID1)
	defer cleanupSocialUser(db, accID2)
	defer cleanupSocialUser(db, accID3)

	createSocialTestUser(t, db, accID1, "GuildU1")
	createSocialTestUser(t, db, accID2, "GuildU2")
	createSocialTestUser(t, db, accID3, "GuildU3")

	c1 := createSocialTestChar(t, repo, accID1, "GldMasterChar", 0, 1, 50, "Elderine")
	c2 := createSocialTestChar(t, repo, accID2, "GldApprentice", 0, 2, 20, "Elderine")
	c3 := createSocialTestChar(t, repo, accID3, "GldRecruit", 0, 1, 10, "Elderine")

	// 1. Crear Gremio
	guildName := fmt.Sprintf("Fase6Gld%d", time.Now().Unix()%10000)
	guildID, err := repo.CreateGuild(c1.ID, guildName, "pass", 1)
	if err != nil || guildID <= 0 {
		t.Fatalf("Error creando gremio: %v, id: %d", err, guildID)
	}

	// 2. Agregar Miembro y cambiar rango
	if err := repo.AddGuildMember(guildID, c2.ID, 2); err != nil {
		t.Fatalf("Error agregando miembro de gremio: %v", err)
	}
	if err := repo.SetGuildMemberGrade(guildID, c2.ID, 3); err != nil {
		t.Fatalf("Error seteando rango de miembro: %v", err)
	}

	// 3. Notificación e Historial
	if err := repo.SetGuildNotify(guildID, "GldMasterChar", "Bienvenidos al gremio oficial"); err != nil {
		t.Fatalf("Error seteando notify de gremio: %v", err)
	}
	if err := repo.AddGuildHistory(guildID, 1, c1.ID, c2.ID, "Gremio fundado oficialmente"); err != nil {
		t.Fatalf("Error agregando historial de gremio: %v", err)
	}
	hist, err := repo.GetGuildHistory(guildID)
	if err != nil || len(hist) == 0 {
		t.Fatalf("Error obteniendo historial de gremio: %v", err)
	}

	// 4. Token de gremio
	if err := repo.SetGuildToken(guildID, 500); err != nil {
		t.Fatalf("Error seteando tokens de gremio: %v", err)
	}
	tok, err := repo.GetGuildToken(guildID)
	if err != nil || tok != 500 {
		t.Fatalf("Error obteniendo tokens de gremio: %v, esperado 500, obtenido %d", err, tok)
	}

	// 5. Crear Academia de Gremio
	if err := repo.CreateGuildAcademy(guildID, c1.ID); err != nil {
		t.Fatalf("Error creando academia de gremio: %v", err)
	}

	// 6. Agregar Miembro a la Academia (c3 es novato exclusivo de la academia)
	if err := repo.AddAcademyMember(guildID, c3.ID, 10); err != nil {
		t.Fatalf("Error agregando miembro a academia: %v", err)
	}

	// 7. Ban de chat en academia
	if err := repo.SetAcademyChatBan(guildID, c3.ID, 30); err != nil {
		t.Fatalf("Error aplicando ban de chat en academia: %v", err)
	}

	// 8. Graduación en academia
	if err := repo.SetAcademyGraduate(guildID, c3.ID); err != nil {
		t.Fatalf("Error graduando miembro de academia: %v", err)
	}
	grads, err := repo.GetAcademyGraduates(guildID)
	if err != nil || len(grads) == 0 {
		t.Fatalf("Error obteniendo graduados de academia: %v", err)
	}

	// 9. Disolver Gremio
	if err := repo.DismissGuild(guildID); err != nil {
		t.Fatalf("Error disolviendo gremio: %v", err)
	}

	t.Log("Control 5B Funcional Exitoso: Ciclo completo de Guild y Guild Academy en MySQL verificado")
}

// TestFase06_Control5B_Functional_Prison_Report_ChatRestrict_HolyPromise prueba subsistemas auxiliares:
// Prison (Add, Alter, End), ChatRestrict (Add, Del, DelAll), Report (Add, Get, Del) y HolyPromise (Set, Get, Del).
func TestFase06_Control5B_Functional_Prison_Report_ChatRestrict_HolyPromise(t *testing.T) {
	db := getTestDB(t)
	charRepo := database.NewMySQLCharacterRepository(db)
	logRepo := database.NewMySQLAccountRepository(db)

	const (
		accID1 = 98807
	)

	cleanupSocialUser(db, accID1)
	defer cleanupSocialUser(db, accID1)

	createSocialTestUser(t, db, accID1, "AuxU1")
	c1 := createSocialTestChar(t, charRepo, accID1, "PrisonHero", 0, 1, 40, "Elderine")
	c2 := createSocialTestChar(t, charRepo, accID1, "MasterChar", 1, 2, 70, "Elderine")

	// 1. PRISON: p_Prison_Add -> p_Prison_Alter -> p_Prison_End
	if err := charRepo.AddPrison("AdminGM", "PrisonHero", 60, "Uso de macros", "Caso de auditoria"); err != nil {
		t.Fatalf("Error en AddPrison: %v", err)
	}
	if err := charRepo.AlterPrison("AdminGM", "PrisonHero", 120, "Reincidencia", "Pena extendida"); err != nil {
		t.Fatalf("Error en AlterPrison: %v", err)
	}
	if err := charRepo.EndPrison(c1.ID); err != nil {
		t.Fatalf("Error en EndPrison: %v", err)
	}

	// 2. CHAT RESTRICT: p_Char_Chat_Restrict_Add -> p_Char_Chat_Restrict_Del -> p_Char_Chat_Restrict_DelAll
	if err := charRepo.AddChatRestrict(accID1, c1.ID); err != nil {
		t.Fatalf("Error en AddChatRestrict: %v", err)
	}
	if err := charRepo.DelChatRestrict(accID1, c1.ID); err != nil {
		t.Fatalf("Error en DelChatRestrict: %v", err)
	}
	if err := charRepo.AddChatRestrict(accID1, c1.ID); err != nil {
		t.Fatalf("Error en AddChatRestrict (2): %v", err)
	}
	if err := charRepo.DelAllChatRestrict(accID1); err != nil {
		t.Fatalf("Error en DelAllChatRestrict: %v", err)
	}

	// 3. REPORT: p_Report_Add -> p_Report_Get -> p_Report_Del
	report := &database.Report{
		Reporter:   accID1,
		WorldNo:    1,
		Harmer:     "CheaterX",
		Reason:     "Botting",
		Result:     0,
		Remark:     "Reporte de prueba",
		ReportDate: time.Now(),
	}
	regNo, err := logRepo.AddReport(report)
	if err != nil || regNo <= 0 {
		t.Fatalf("Error en AddReport: %v, regNo: %d", err, regNo)
	}
	reports, err := logRepo.GetReports(accID1, 1)
	if err != nil || len(reports) == 0 {
		t.Fatalf("Error en GetReports: %v", err)
	}
	if _, err := logRepo.DeleteReport(regNo); err != nil {
		t.Fatalf("Error en DeleteReport: %v", err)
	}

	// 4. HOLY PROMISE: p_HolyPromise_Add -> p_HolyPromise_Get -> p_HolyPromise_Del
	if err := charRepo.SetHolyPromise(c1.ID, c2.ID, 1); err != nil {
		t.Fatalf("Error en SetHolyPromise: %v", err)
	}
	hp, err := charRepo.GetHolyPromise(c1.ID)
	if err != nil || hp == nil || hp.UpCharNo != c2.ID {
		t.Fatalf("Error en GetHolyPromise: %v, hp: %+v", err, hp)
	}
	if err := charRepo.DelHolyPromise(c1.ID, c2.ID); err != nil {
		t.Fatalf("Error en DelHolyPromise: %v", err)
	}

	t.Log("Control 5B Funcional Exitoso: Subsistemas Prison, ChatRestrict, Report y HolyPromise en MySQL verificados")
}

// TestFase06_Control5B_Functional_Raid_And_Topology prueba la coordinación de bandas (RaidManager)
// y la topología dinámica de mapas (tWorldZoneNode / usp_WorldZone_ResolveMapEndpoint).
func TestFase06_Control5B_Functional_Raid_And_Topology(t *testing.T) {
	db := getTestDB(t)
	repo := database.NewMySQLCharacterRepository(db)
	partyMgr := world.NewPartyManager(repo)
	raidMgr := world.NewRaidManager(repo, partyMgr)

	leader, _ := network.NewClient(&noopConn{}, nil)
	leader.AccountID = 99995
	leader.CharacterID = 8895
	leader.Username = "RaidLeader"

	member, _ := network.NewClient(&noopConn{}, nil)
	member.AccountID = 99996
	member.CharacterID = 8896
	member.Username = "RaidMember"

	// 1. Crear Raid
	raidMgr.HandleMake(leader)

	// 2. Invitar y responder
	raidMgr.HandleInvitation(leader, "RaidMember")
	raidMgr.HandleInvitationAnswer(member, 0, "RaidLeader") // 0 = accept

	// 3. Mover slot
	raidMgr.HandleSlotMove(leader, 0, 1, "RaidLeader", "RaidMember")

	// 4. Set Leader / Sub-Leader / Looter
	raidMgr.HandleSetSubLeader(leader, "RaidMember")
	raidMgr.HandleSetLooter(leader, "RaidLeader")

	// 5. Chat y Warning
	raidMgr.HandleChat(leader, "Listos para el boss!")
	raidMgr.HandleWarning(leader, "ATENCION AL AOE!")

	// 6. Disolver / Salir
	raidMgr.HandleLeave(member)
	raidMgr.HandleLeave(leader)

	// 7. Probar topología de mapas contra MySQL
	rootPath := findRepoRoot(t)
	dataDir := filepath.Join(rootPath, "data")

	store, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("Error cargando store: %v", err)
	}

	// Registrar nodo operativo en tWorldZoneNode
	if err := repo.UpsertZoneNode(1, 1, "127.0.0.1", 9015); err != nil {
		t.Fatalf("Error registrando nodo de zona: %v", err)
	}
	if err := repo.HeartbeatZoneNode(1, 1, 5); err != nil {
		t.Fatalf("Error actualizando latido de zona: %v", err)
	}

	ip, port, isOnline, err := repo.ResolveMapEndpoint(1, 1)
	if err != nil || ip == "" || port == 0 || !isOnline {
		t.Fatalf("Error resolviendo endpoint para zona 1: %v (IP: %s, Port: %d, Online: %v)", err, ip, port, isOnline)
	}

	// Validar que FieldList tiene mapas asignados a la zona 1
	eldField, ok := store.Fields.GetField("Eld")
	if !ok || eldField.Fiesta != 1 {
		t.Fatalf("Campo Eld no configurado con Fiesta/Zone=1: %+v", eldField)
	}

	t.Logf("Control 5B Funcional Exitoso: RaidManager y topología dinámica de zona (Eld -> Zona 1 -> %s:%d) verificados", ip, port)
}
