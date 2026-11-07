package audit_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

// TestFase05_Control5A_SHNFiles_And_DataLoading verifica la carga de los 74+ archivos SHN únicos,
// las 27 clases de ClassParam, FieldList.json + TutorialMap.json, y default/tutorial character data.
func TestFase05_Control5A_SHNFiles_And_DataLoading(t *testing.T) {
	rootPath := findRepoRoot(t)
	dataDir := filepath.Join(rootPath, "data")

	store, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("LoadWorldDataStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil WorldDataStore")
	}

	// 1. Verificar carga de los 74+ archivos SHN únicos (54 de Character U 71 de WorldManager)
	if len(store.LoadedFiles) < 74 {
		t.Errorf("CRITERIO FALLIDO: Se esperaban al menos 74 archivos SHN cargados en LoadedFiles, se obtuvieron %d",
			len(store.LoadedFiles))
	}

	for _, name := range world.CharacterSHNFiles {
		if store.LoadedFiles[name] == nil {
			t.Errorf("CRITERIO FALLIDO: Archivo SHN de Character.exe '%s' no está cargado", name)
		}
	}
	for _, name := range world.WorldManagerSHNFiles {
		if store.LoadedFiles[name] == nil {
			t.Errorf("CRITERIO FALLIDO: Archivo SHN de WorldManager.exe '%s' no está cargado", name)
		}
	}

	// 2. Verificar ClassParam con las 27 clases oficiales
	if store.ClassParams == nil {
		t.Fatal("CRITERIO FALLIDO: ClassParams store es nil")
	}
	if count := store.ClassParams.ClassCount(); count != 27 {
		t.Errorf("CRITERIO FALLIDO: Se esperaban 27 clases en ClassParams, se obtuvieron %d", count)
	}

	// 3. Verificar Fields y TutorialMap
	if store.Fields == nil {
		t.Fatal("CRITERIO FALLIDO: Fields store es nil")
	}
	if store.Fields.FieldCount() == 0 {
		t.Errorf("CRITERIO FALLIDO: No se cargaron mapas en Fields store")
	}
	if store.Fields.TutorialMap() == nil {
		t.Errorf("CRITERIO FALLIDO: TutorialMap es nil")
	}

	// 4. Verificar default_character_data.json (6 clases base incluyendo Sentinel Lv 60 en Eld)
	if store.DefaultData == nil {
		t.Fatal("CRITERIO FALLIDO: DefaultData es nil")
	}
	for _, classID := range []byte{1, 6, 11, 16, 21, 26} {
		cd, ok := store.DefaultData.GetClassData(classID)
		if !ok {
			t.Errorf("CRITERIO FALLIDO: Clase base %d no encontrada en DefaultData", classID)
		}
		if cd.HP <= 0 || cd.SP <= 0 {
			t.Errorf("CRITERIO FALLIDO: Stats inválidos para clase %d: HP=%d SP=%d", classID, cd.HP, cd.SP)
		}
	}
	sentinel, ok := store.DefaultData.GetClassData(26)
	if !ok || sentinel.StartMap != "Eld" || sentinel.InitLevel != 60 || sentinel.InitExp != 24346145 {
		t.Errorf("CRITERIO FALLIDO: Parámetros de Sentinel (26) incorrectos: %+v", sentinel)
	}

	// 5. Verificar tutorial_character_data.json
	if store.TutorialData == nil {
		t.Fatal("CRITERIO FALLIDO: TutorialData es nil")
	}
	for _, classID := range []byte{1, 6, 11, 16, 21} {
		td, ok := store.TutorialData.GetTutorialData(classID)
		if !ok {
			t.Errorf("CRITERIO FALLIDO: Clase tutorial %d no encontrada en TutorialData", classID)
		}
		if td.StartMap != "RouN" || td.StartX != 7178 || td.StartY != 7011 || td.Money != 300 {
			t.Errorf("CRITERIO FALLIDO: Parámetros de tutorial incorrectos para clase %d: %+v", classID, td)
		}
	}

	// 6. Verificar resolución dinámica de cupones de belleza desde ItemInfo.shn
	if c1 := store.GetStyleCouponID(1); c1 != 30490 {
		t.Errorf("CRITERIO FALLIDO: HairShop01 cupón esperado 30490, obtenido %d", c1)
	}
	if c5 := store.GetStyleCouponID(5); c5 != 30494 {
		t.Errorf("CRITERIO FALLIDO: HairShop05 cupón esperado 30494, obtenido %d", c5)
	}

	// 7. Verificar nombres de clase desde ClassName.shn
	if name := store.GetClassName(1); name != "Fighter" {
		t.Errorf("CRITERIO FALLIDO: Nombre de clase 1 esperado 'Fighter', obtenido '%s'", name)
	}
	if name := store.GetClassName(26); name != "Sentinel" {
		t.Errorf("CRITERIO FALLIDO: Nombre de clase 26 esperado 'Sentinel', obtenido '%s'", name)
	}

	// 8. Verificar recompensas de asistencia desde AttendReward.shn
	if rew, ok := store.GetAttendReward(0); !ok || rew.ItemID == 0 {
		t.Errorf("CRITERIO FALLIDO: Recompensa de asistencia en índice 0 no válida: %+v", rew)
	}
}

// TestFase05_Control5A_BadNameFilter_ZeroHardcode verifica que BadNameFilter no contenga
// slices hardcodeados en código Go y filtre exclusivamente por BadNameFilter.shn.
func TestFase05_Control5A_BadNameFilter_ZeroHardcode(t *testing.T) {
	rootPath := findRepoRoot(t)
	staticDataFile := filepath.Join(rootPath, "pkg", "server", "world", "static_data.go")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, staticDataFile, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("Error parseando static_data.go: %v", err)
	}

	forbiddenWords := []string{"admin", "gm", "master", "support", "fiesta", "gamigo", "onson"}

	ast.Inspect(node, func(n ast.Node) bool {
		// A. Prohibir campo subStrings en struct BadNameFilter
		if typeSpec, ok := n.(*ast.TypeSpec); ok && typeSpec.Name.Name == "BadNameFilter" {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				for _, f := range structType.Fields.List {
					for _, name := range f.Names {
						if name.Name == "subStrings" {
							t.Errorf("CRITERIO FALLIDO: BadNameFilter contiene campo legado hardcodeado 'subStrings' en línea %d",
								fset.Position(f.Pos()).Line)
						}
					}
				}
			}
		}

		// B. Prohibir literales de string de palabras prohibidas hardcodeadas en LoadBadNameFilter
		if cl, ok := n.(*ast.CompositeLit); ok {
			for _, elt := range cl.Elts {
				if bl, ok := elt.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					val := strings.Trim(strings.ToLower(bl.Value), "\"")
					for _, fw := range forbiddenWords {
						if val == fw {
							t.Errorf("CRITERIO FALLIDO: Palabra prohibida hardcodeada '%s' detectada en slice en línea %d",
								val, fset.Position(bl.Pos()).Line)
						}
					}
				}
			}
		}

		return true
	})

	// C. Validar comportamiento funcional real con BadNameFilter.shn
	dataDir := filepath.Join(rootPath, "data")
	filter, err := world.LoadBadNameFilter(filepath.Join(dataDir, "BadNameFilter.shn"))
	if err != nil {
		t.Fatalf("Error cargando BadNameFilter.shn: %v", err)
	}

	if !filter.IsBadName("admin") {
		t.Errorf("CRITERIO FALLIDO: 'admin' debe ser bloqueado por BadNameFilter.shn")
	}
	if !filter.IsBadName("AdminMaster") {
		t.Errorf("CRITERIO FALLIDO: 'AdminMaster' debe ser bloqueado por subcadena de BadNameFilter.shn")
	}
	if !filter.IsBadName("gm") {
		t.Errorf("CRITERIO FALLIDO: 'gm' debe ser bloqueado por BadNameFilter.shn")
	}
	if filter.IsBadName("Legolas") {
		t.Errorf("CRITERIO FALLIDO: 'Legolas' no debe ser bloqueado")
	}
}

// TestFase05_Control5A_DispatcherCoverage_And_ZeroInlineSQL valida AST de pkg/server/world/:
// 1. Cero SQL inline en todo el paquete world.
// 2. Cobertura del 100% de los OpCodes oficiales de la Fase 05 en HandlePacket.
func TestFase05_Control5A_DispatcherCoverage_And_ZeroInlineSQL(t *testing.T) {
	rootPath := findRepoRoot(t)
	worldDir := filepath.Join(rootPath, "pkg", "server", "world")

	files := []string{
		filepath.Join(worldDir, "handler.go"),
		filepath.Join(worldDir, "handler_avatar_options.go"),
		filepath.Join(worldDir, "packet_helper.go"),
		filepath.Join(worldDir, "static_data.go"),
	}

	fset := token.NewFileSet()
	handledOpCodes := make(map[string]bool)

	for _, file := range files {
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("Error parseando AST de %s: %v", file, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// A. Cero SQL inline en strings
			if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
				val := strings.ToUpper(bl.Value)
				sqlVerbs := []string{"SELECT ", "INSERT INTO ", "UPDATE ", "DELETE FROM "}
				for _, verb := range sqlVerbs {
					if strings.Contains(val, verb) && !strings.Contains(val, "CALL ") {
						t.Errorf("CRITERIO FALLIDO: Consulta SQL inline detectada en %s:%d: %s",
							file, fset.Position(bl.Pos()).Line, bl.Value)
					}
				}
			}

			// B. Registrar case tags en HandlePacket
			if sw, ok := n.(*ast.SwitchStmt); ok {
				for _, stmt := range sw.Body.List {
					if cc, ok := stmt.(*ast.CaseClause); ok {
						for _, expr := range cc.List {
							if sel, ok := expr.(*ast.SelectorExpr); ok {
								handledOpCodes[sel.Sel.Name] = true
							}
						}
					}
				}
			}

			return true
		})
	}

	requiredOpCodes := []string{
		// Dept 3
		"PROTO_NC_USER_LOGINWORLD_REQ",
		"PROTO_NC_USER_NORMALLOGOUT_CMD",
		"PROTO_NC_USER_AVATAR_LIST_REQ",
		"PROTO_NC_USER_WILL_WORLD_SELECT_REQ",
		// Dept 4
		"PROTO_NC_CHAR_LOGIN_REQ",
		"PROTO_NC_CHAR_SET_STYLE_REQ",
		"PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ",
		"PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD",
		"PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ",
		"PROTO_NC_CHAR_TUTORIAL_POPUP_ACK",
		// Dept 5
		"PROTO_NC_AVATAR_CREATE_REQ",
		"PROTO_NC_AVATAR_ERASE_REQ",
		"PROTO_NC_AVATAR_RENAME_REQ",
		// Dept 28
		"PROTO_NC_CHAR_OPTION_GET_ALL_REQ",
		"PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ",
		"PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ",
		"PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ",
		"PROTO_NC_CHAR_OPTION_GET_SOUND_REQ",
		"PROTO_NC_CHAR_OPTION_GET_GAME_REQ",
		"PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ",
		"PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ",
		"PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD",
		"PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD",
		"PROTO_NC_CHAR_OPTION_SET_VIDEO_CMD",
		"PROTO_NC_CHAR_OPTION_SET_SOUND_CMD",
		"PROTO_NC_CHAR_OPTION_SET_GAME_CMD",
		"PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD",
		"PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD",
		"PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ",
		"PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ",
		"PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ",
		"PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ",
		"PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ",
		"PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR",
		// Dept 45
		"PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ",
	}

	for _, op := range requiredOpCodes {
		if !handledOpCodes[op] {
			t.Errorf("CRITERIO FALLIDO: OpCode oficial requerido '%s' no está manejado en HandlePacket", op)
		}
	}
}

// TestFase05_Control5A_FullMySQLAvatarCycle prueba el ciclo de vida completo de avatares contra MySQL xelion_online.
func TestFase05_Control5A_FullMySQLAvatarCycle(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repo := database.NewMySQLCharacterRepository(db)

	const (
		testAccountID = 96666
		testCharName1 = "Fase5Fighter"
		testCharName2 = "Fase5Sentinel"
	)

	// Limpieza defensiva previa
	cleanup := func() {
		_, _ = db.Exec("DELETE FROM tItem WHERE nOwner IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterShape WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterState WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterShortCut WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterKeyMap WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterGameOption WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterOptions WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacterChangeID WHERE nUserNo = ?", testAccountID)
		_, _ = db.Exec("DELETE FROM tInGameGuide WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tAttendance WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tGuildMember WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testAccountID)
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", testAccountID)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testAccountID)
	}

	cleanup()
	defer cleanup()

	// Crear usuario de prueba
	_, err := db.Exec("INSERT INTO tUser (nUserNo, sUserID, sUserPW, nAuthID) VALUES (?, 'Fase5User', 'pass123', 1)", testAccountID)
	if err != nil {
		t.Fatalf("Error insertando tUser: %v", err)
	}

	// 1. Creación de Personaje Base (Fighter, Slot 0)
	char1 := &database.Character{
		AccountID: testAccountID,
		Name:      testCharName1,
		Slot:      0,
		Job:       1,
		Level:     1,
		Male:      true,
		Hair:      1,
		HairColor: 1,
		Face:      1,
		TownName:  "RouN",
		XPos:      6888,
		YPos:      7045,
		CurHP:     100,
		CurSP:     100,
		Equips:    map[byte]uint16{12: 250}, // Slot 12 = Weapon (Shortsword)
		Items: []*database.Item{
			{StorageType: 9, Storage: 0, ItemID: 46003, Flags: 5},
		},
		Shortcuts: []*database.Shortcut{
			{SlotNo: 0, CodeType: 1, CodeNo: 1, Value: 0},
		},
	}
	char1.TCharacterState.Money = 500
	char1.TCharacterState.HPStones = 15
	char1.TCharacterState.SPStones = 7

	if err := repo.CreateCharacter(char1); err != nil {
		t.Fatalf("Error en CreateCharacter Fighter: %v", err)
	}
	if char1.ID <= 0 {
		t.Fatalf("Expected valid CharID, got %d", char1.ID)
	}

	// 2. Creación de Sentinel (Class 26, Slot 1) en Eld a Nivel 60 con Exp 24346145
	char2 := &database.Character{
		AccountID: testAccountID,
		Name:      testCharName2,
		Slot:      1,
		Job:       26,
		Race:      2,
		Level:     60,
		Male:      false,
		Hair:      2,
		HairColor: 2,
		Face:      2,
		TownName:  "Eld",
		XPos:      12259,
		YPos:      13434,
		CurHP:     2500,
		CurSP:     1200,
		Equips:    map[byte]uint16{12: 62923},
	}
	char2.TCharacterState.Money = 50000
	char2.TCharacterState.Exp = 24346145

	if err := repo.CreateCharacter(char2); err != nil {
		t.Fatalf("Error en CreateCharacter Sentinel: %v", err)
	}

	// Verificar Sentinel en MySQL
	var dbLevel int
	var dbExp int64
	var dbMap string
	err = db.QueryRow("SELECT nLevel, nExp, sLoginZone FROM tCharacter WHERE nCharNo = ?", char2.ID).Scan(&dbLevel, &dbExp, &dbMap)
	if err != nil {
		t.Fatalf("Error consultando Sentinel en MySQL: %v", err)
	}
	if dbLevel != 60 || dbExp != 24346145 || dbMap != "Eld" {
		t.Errorf("CRITERIO FALLIDO: Sentinel no persistido correctamente: Level=%d, Exp=%d, Map=%s", dbLevel, dbExp, dbMap)
	}

	// 3. Renombramiento de personaje (ChangeCharacterName)
	const renamedName = "Fase5Renamed"
	if err := repo.ChangeCharacterName(char1.ID, renamedName); err != nil {
		t.Fatalf("Error en ChangeCharacterName: %v", err)
	}
	var currentName string
	_ = db.QueryRow("SELECT sID FROM tCharacter WHERE nCharNo = ?", char1.ID).Scan(&currentName)
	if currentName != renamedName {
		t.Errorf("CRITERIO FALLIDO: Nombre no actualizado en MySQL, esperado '%s', obtenido '%s'", renamedName, currentName)
	}

	// 4. Cambio de Estilo (SetCharacterShape)
	if err := repo.SetCharacterShape(char1.ID, 1, 1, 1, 3, 4, 2); err != nil {
		t.Fatalf("Error en SetCharacterShape: %v", err)
	}
	var dbHair, dbColor, dbFace byte
	_ = db.QueryRow("SELECT nHairType, nHairColor, nFaceShape FROM tCharacterShape WHERE nCharNo = ?", char1.ID).Scan(&dbHair, &dbColor, &dbFace)
	if dbHair != 3 || dbColor != 4 || dbFace != 2 {
		t.Errorf("CRITERIO FALLIDO: Shape no actualizado en MySQL: Hair=%d Color=%d Face=%d", dbHair, dbColor, dbFace)
	}

	// 5. Opciones Mejoradas (IMPROVE)
	if err := repo.SetImproveShortCut(char1.ID, 0, 101, 555); err != nil {
		t.Fatalf("Error en SetImproveShortCut: %v", err)
	}
	scList, err := repo.GetImproveShortCutList(char1.ID)
	foundSC := false
	for _, sc := range scList {
		if sc.SlotNo == 0 && sc.CodeNo == 101 && sc.Value == 555 {
			foundSC = true
			break
		}
	}
	if err != nil || !foundSC {
		t.Errorf("CRITERIO FALLIDO: GetImproveShortCutList no retornó valor guardado: %+v", scList)
	}

	if err := repo.SetImproveKeyMap(char1.ID, 5, 1, 'Q'); err != nil {
		t.Fatalf("Error en SetImproveKeyMap: %v", err)
	}
	kmList, err := repo.GetImproveKeyMap(char1.ID)
	foundKM := false
	for _, km := range kmList {
		if km.FunctionNo == 5 && km.ASCIICode == 'Q' && km.ExtendKey == 1 {
			foundKM = true
			break
		}
	}
	if err != nil || !foundKM {
		t.Errorf("CRITERIO FALLIDO: GetImproveKeyMap no retornó valor guardado para funcNo 5: %+v", kmList)
	}

	if err := repo.SetImproveGameOption(char1.ID, 10, 1); err != nil {
		t.Fatalf("Error en SetImproveGameOption: %v", err)
	}
	goList, err := repo.GetImproveGameOptions(char1.ID)
	foundGO := false
	for _, opt := range goList {
		if opt.OptionNo == 10 && opt.Value == 1 {
			foundGO = true
			break
		}
	}
	if err != nil || !foundGO {
		t.Errorf("CRITERIO FALLIDO: GetImproveGameOptions no retornó valor guardado para optionNo 10: %+v", goList)
	}

	// 6. Guía de Novato y Asistencia
	if err := repo.SetInGameGuideView(char1.ID, 7); err != nil {
		t.Fatalf("Error en SetInGameGuideView: %v", err)
	}
	guides, err := repo.GetInGameGuideViewList(char1.ID)
	if err != nil || len(guides) == 0 || guides[0] != 7 {
		t.Errorf("CRITERIO FALLIDO: GetInGameGuideViewList no retornó guía vista: %v", guides)
	}

	if _, err := repo.SetAttendanceRewardAt(char1.ID, time.Now(), 1); err != nil {
		t.Fatalf("Error en SetAttendanceRewardAt: %v", err)
	}

	// 7. Verificación de rechazo de borrado si es Maestro de Gremio
	// Asignar personaje 1 como Guild Master en tGuildMember
	_, _ = db.Exec("INSERT INTO tGuild (nNo, sID, nMasterCharNo) VALUES (99999, 'AuditGuild', ?)", char1.ID)
	_, _ = db.Exec("INSERT INTO tGuildMember (nNo, nCharNo, nGrade) VALUES (99999, ?, 0)", char1.ID)

	isGM, err := repo.IsActiveGuildMaster(char1.ID)
	if err != nil || !isGM {
		t.Errorf("CRITERIO FALLIDO: IsActiveGuildMaster debió retornar true para el Guild Master")
	}

	// Limpiar membresía de gremio para permitir borrado
	_, _ = db.Exec("DELETE FROM tGuildMember WHERE nNo = 99999")
	_, _ = db.Exec("DELETE FROM tGuild WHERE nNo = 99999")

	// 8. Soft-delete oficial (DeleteCharacter)
	if err := repo.DeleteCharacter(testAccountID, 0); err != nil {
		t.Fatalf("Error en DeleteCharacter slot 0: %v", err)
	}
	var bDeleted byte
	_ = db.QueryRow("SELECT bDeleted FROM tCharacter WHERE nCharNo = ?", char1.ID).Scan(&bDeleted)
	if bDeleted != 1 {
		t.Errorf("CRITERIO FALLIDO: Personaje no marcado como soft-deleted (bDeleted = %d)", bDeleted)
	}

	// Verificar que el slot 0 quedó libre para nuevo personaje
	charsAfter, err := repo.GetCharactersByAccount(testAccountID)
	if err != nil {
		t.Fatalf("Error obteniendo personajes tras borrado: %v", err)
	}
	for _, ch := range charsAfter {
		if ch.Slot == 0 {
			t.Errorf("CRITERIO FALLIDO: Slot 0 sigue ocupado tras DeleteCharacter: %+v", ch)
		}
	}
}
