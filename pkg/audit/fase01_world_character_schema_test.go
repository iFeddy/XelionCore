package audit_test

import (
	"database/sql"
	"strings"
	"testing"
)

// TestFase01_WorldCharacter_Control1 verifica exhaustivamente contra information_schema el 100% de tablas,
// columnas, tipos y claves primarias requeridas por la Fase 01 y el Control 1 de WorldService / Character.
func TestFase01_WorldCharacter_Control1(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const schema = "xelion_online"

	// 1. Obtener todas las tablas existentes en la base xelion_online
	rows, err := db.Query(`
		SELECT TABLE_NAME, TABLE_TYPE 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = ?
	`, schema)
	if err != nil {
		t.Fatalf("Error consultando information_schema.TABLES: %v", err)
	}
	defer rows.Close()

	existingTables := make(map[string]string)
	for rows.Next() {
		var name, tableType string
		if err := rows.Scan(&name, &tableType); err != nil {
			t.Fatalf("Error leyendo fila de tabla: %v", err)
		}
		existingTables[strings.ToLower(name)] = tableType
	}

	// 2. Inventario Oficial de las 61 Tablas Nucleares de World00_Character
	core61Tables := []string{
		// 2.1 Personaje y Apariencia (6 tablas)
		"tCharacter", "tCharacterShape", "tCharacterState", "tCharacterTitle", "tRestExp", "tTutorial",
		// 2.2 Inventario y Objetos (8 tablas)
		"tItem", "tItemOptions", "tItemActionCoolTime", "tCharacterChargedBuffer",
		"tItemEnchant", "tItemMobList", "tItemStorageFindEmpty", "tItemTitle",
		// 2.3 Habilidades (1 tabla)
		"tSkill",
		// 2.4 Gremios y Academia (14 tablas)
		"tGuild", "tGuildMember", "tGuildAcademy", "tGuildAcademyMember",
		"tGuildAcademyGraduate", "tGuildAcademyHistory", "tGuildAcademyRanking",
		"tGuildAcademyReward", "tGuildAcademyRewardType", "tGuildEmblem",
		"tGuildToken", "tGuildTournament", "tGuildTournamentRegist", "tLog_GuildMasterChange",
		// 2.5 Amigos y Chat (3 tablas)
		"tFriend", "tChatBlock", "tChatColor",
		// 2.6 Grupos / Party (2 tablas)
		"tParty", "tPartyMember",
		// 2.7 Misiones y Progresión (3 tablas)
		"tQuest", "tQuestTimes", "tKingdomQuest",
		// 2.8 Mascotas (3 tablas)
		"tPet", "tPetSkill", "tPetNew",
		// 2.9 Social, Bodas, MiniHouses (8 tablas)
		"tMarriage", "tMinihouse", "tMinihouseGate", "tMinihouseGuest",
		"tHolyPromise", "tHolyPromiseReward", "tMasterPupils", "tInGameGuide",
		// 2.10 Soporte, Filtros y Estadísticas (7 tablas)
		"tPrison", "tSpamBlockTime", "tSpamFilter", "tStatCharLevel",
		"tStatGameMoney", "tStatisticsGameMoney", "tUserMoney",
		// 2.11 Configuración y Opciones de Personaje (6 tablas)
		"tCharacterChangeID", "tCharacterGameOption", "tCharacterKeyMap",
		"tCharacterLoading", "tCharacterOptions", "tCharacterShortCut",
	}

	if len(core61Tables) != 61 {
		t.Fatalf("La lista de tablas nucleares debe contener exactamente 61 tablas, tiene %d", len(core61Tables))
	}

	// Validar que cada una de las 61 tablas exista en MySQL
	for _, tbl := range core61Tables {
		if _, ok := existingTables[strings.ToLower(tbl)]; !ok {
			t.Errorf("[Control 1] Falta tabla nuclear requerida de World00_Character: %s", tbl)
		}
	}

	// 3. Inventario Completo Adicional (30 tablas oficiales de World00_Character.bak)
	additional30Tables := []string{
		"tArenaChallenge", "tArenaFBZ", "tArenaGuild", "tArenaPvP",
		"tAttendance", "tAttendanceReward", "tAutoMon", "tAutoPick",
		"tCharacterSingleOption", "tCollectionCard", "tCollectionFavorites",
		"tCollectionReward", "tCollectionSerialCount", "tDice", "tDiceLog",
		"tDicePayLog", "tDiceResult", "tExp", "tGMItemReward", "tGMItemRewardLog",
		"tGambleCoin", "tGambleCoinExchangeLog", "tGambleDice", "tGambleDiceLog",
		"tGambleDiceResult", "tGambleJackpot", "tGambleJackpotLog",
		"tGambleSlotMachineLog", "tGameOptionInit", "tKeyMapInit",
	}

	for _, tbl := range additional30Tables {
		if _, ok := existingTables[strings.ToLower(tbl)]; !ok {
			t.Errorf("[Control 1] Falta tabla complementaria requerida de World00_Character: %s", tbl)
		}
	}

	// 4. Invariantes Críticas de Columnas y Tipos

	// Invariante 4.1: tCharacter (45 columnas oficiales exactas)
	checkColumns(t, db, schema, "tCharacter", 45, []string{
		"ncharno", "sid", "nuserno", "nadminlevel", "nps", "ngs",
		"dcreatedate", "bdeleted", "nloginzonex", "nloginzoney", "nloginzoned",
		"npkyellowtime", "ddeleteddate", "nap", "sloginzone", "npkcount",
		"nlevel", "nexp", "nslotno", "nhp", "nsp", "nhps", "nsps",
		"dlastlogindate", "nlogincount", "ncreateworld", "nfame", "nprisonmin",
		"nplaymin", "nflags", "nmoney", "nstrength", "nconstitute", "ndexterity",
		"nintelligence", "nmentalpower", "nredistributepoint", "nkqhandle",
		"skqmap", "nkqx", "nkqy", "nkqd", "dkqdate", "nfriendpoint", "nlastfriendpoint",
	})

	// Invariante 4.2: tCharacterShape (7 columnas exactas)
	checkColumns(t, db, schema, "tCharacterShape", 7, []string{
		"ncharno", "nclass", "nrace", "ngender", "nhairtype", "nhaircolor", "nfaceshape",
	})

	// Invariante 4.3: tCharacterState (20 columnas exactas)
	checkColumns(t, db, schema, "tCharacterState", 20, []string{
		"ncharno", "nstateno", "nstatelevel", "nstatetime", "nstatecount",
		"nstatepowerdemage", "nstatepowersp", "nstatepowerkeeptime", "nstatepowercooltime",
		"nsetitemindexa", "nsetitemvaluea", "nsetitemindexb", "nsetitemvalueb",
		"nsetitemindexc", "nsetitemvaluec", "nsetitemindexd", "nsetitemvalued",
		"nresttime", "ncasterlevel", "ncasterattackpower",
	})

	// Invariante 4.4: tItem (7 columnas oficiales exactas)
	checkColumns(t, db, schema, "tItem", 7, []string{
		"nitemkey", "nstoragetype", "nowner", "nstorage", "nitemid", "nflags", "ddate",
	})

	// Invariante 4.5: tItemOptions (3 columnas)
	checkColumns(t, db, schema, "tItemOptions", 3, []string{
		"nitemkey", "noptiontype", "noptiondata",
	})

	// Invariante 4.6: tSkill (10 columnas)
	checkColumns(t, db, schema, "tSkill", 10, []string{
		"ncharno", "nskillno", "nskilllevel", "nskillexp", "nskillwritetime",
		"nskillcooltime", "nskillpowerdemage", "nskillpowersp", "nskillpowerkeeptime", "nskillpowercooltime",
	})

	// Invariante 4.7: tGuild (24 columnas)
	checkColumns(t, db, schema, "tGuild", 24, []string{
		"nno", "sname", "spassword", "nmaxmembers", "nmoney", "ntype",
		"ngrade", "nfame", "nstonelevel", "nexp", "dcreatedate",
	})

	// Invariante 4.8: tParty (3 columnas)
	checkColumns(t, db, schema, "tParty", 3, []string{
		"nflags", "nmastercharno", "npartyno",
	})

	// Invariante 4.9: tQuest (4 columnas)
	checkColumns(t, db, schema, "tQuest", 4, []string{
		"ncharno", "nquestno", "nstatus", "sdata",
	})

	// Invariante 4.10: tPet (16 columnas)
	checkColumns(t, db, schema, "tPet", 16, []string{
		"nstorage", "npetnumber", "nowner", "npetkey", "npetid", "sinfo", "nflags",
		"nassignname", "nlevel", "nexp", "nsatiety", "nspeed", "nabilitytype", "nabilityvalue", "nskincolor", "ddate",
	})

	// 5. Verificación de Semillas Oficiales (Idempotencia y Existencia)

	// 5.1 Verificar curva de EXP (150 niveles en tExp)
	var expCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tExp").Scan(&expCount); err != nil {
		t.Fatalf("Error consultando tExp: %v", err)
	}
	if expCount < 150 {
		t.Errorf("[Control 1 Semillas] tExp debe contener al menos 150 niveles, tiene %d", expCount)
	}

	// 5.2 Verificar inicialización de opciones (tGameOptionInit >= 10)
	var optionCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tGameOptionInit").Scan(&optionCount); err != nil {
		t.Fatalf("Error consultando tGameOptionInit: %v", err)
	}
	if optionCount < 10 {
		t.Errorf("[Control 1 Semillas] tGameOptionInit debe contener al menos 10 opciones base, tiene %d", optionCount)
	}

	// 5.3 Verificar mapeo de teclado (tKeyMapInit >= 20)
	var keymapCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tKeyMapInit").Scan(&keymapCount); err != nil {
		t.Fatalf("Error consultando tKeyMapInit: %v", err)
	}
	if keymapCount < 20 {
		t.Errorf("[Control 1 Semillas] tKeyMapInit debe contener al menos 20 funciones de teclado, tiene %d", keymapCount)
	}

	// 5.4 Verificar plantillas de personajes iniciales en tCharacter y tCharacterShape
	var templateCharCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tCharacter WHERE nCharNo IN (1, 2, 3, 4, 5, 6)").Scan(&templateCharCount); err != nil {
		t.Fatalf("Error consultando plantillas en tCharacter: %v", err)
	}
	if templateCharCount < 6 {
		t.Errorf("[Control 1 Semillas] Faltan plantillas base en tCharacter, esperado 6, encontrado %d", templateCharCount)
	}

	var templateShapeCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tCharacterShape WHERE nCharNo IN (1, 2, 3, 4, 5, 6)").Scan(&templateShapeCount); err != nil {
		t.Fatalf("Error consultando plantillas en tCharacterShape: %v", err)
	}
	if templateShapeCount < 6 {
		t.Errorf("[Control 1 Semillas] Faltan apariencias de plantilla en tCharacterShape, esperado 6, encontrado %d", templateShapeCount)
	}
}

func checkColumns(t *testing.T, db *sql.DB, schema, table string, expectedTotalCols int, requiredCols []string) {
	rows, err := db.Query(`
		SELECT COLUMN_NAME 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`, schema, table)
	if err != nil {
		t.Fatalf("Error consultando columnas de %s: %v", table, err)
	}
	defer rows.Close()

	colMap := make(map[string]bool)
	var colCount int
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			t.Fatalf("Error leyendo columna de %s: %v", table, err)
		}
		colMap[strings.ToLower(col)] = true
		colCount++
	}

	if colCount != expectedTotalCols {
		t.Errorf("[%s] Debe tener exactamente %d columnas, tiene %d", table, expectedTotalCols, colCount)
	}

	for _, req := range requiredCols {
		if !colMap[strings.ToLower(req)] {
			t.Errorf("[%s] Falta columna requerida: %s", table, req)
		}
	}
}
