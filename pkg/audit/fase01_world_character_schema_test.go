package audit_test

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/server/world"
)

// TestFase01_WorldCharacter_Control1 verifica exhaustivamente contra information_schema el 100% de las 114 tablas
// oficiales de World00_Character.bak, las 2 tablas operativas S2S, invariantes de columnas, semillas auxiliares,
// ausencia total de filas falsas *_Template / Tutorial_* en tCharacter y carga íntegra de los JSONs oficiales de inicio.
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

	// 3. Inventario Adicional Fase 01 (30 tablas oficiales de World00_Character.bak)
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

	// 3.1 Tablas Oficiales Compartidas en schema.sql presentes también en World00_Character.bak (5 tablas)
	shared5BakTables := []string{
		"QuestItemFix", "deleted_item", "removeskill", "skilldata", "tTEMP_ItemClass",
	}

	// 3.2 Las 18 Tablas Oficiales Restantes de World00_Character.bak completadas en Fase 02A
	remaining18BakTables := []string{
		"num", "tItem_deleted", "auditCharacterData", "TEMP_ItemCnt",
		"TEMP_QuestComplete", "Temp_Guar", "Temp_Holy", "Waits_Sysprocesses",
		"chr_a", "chr_b", "chr_bak", "tBugReport_20091204",
		"tCharacterState_Backup", "tCharacterState_isSaveZero_BackUp", "tGuildAcademy_Backup", "tMimic_100610",
		"tQuest_BackUp", "temp_GA",
	}

	all114BakTables := make([]string, 0, 114)
	all114BakTables = append(all114BakTables, core61Tables...)
	all114BakTables = append(all114BakTables, additional30Tables...)
	all114BakTables = append(all114BakTables, shared5BakTables...)
	all114BakTables = append(all114BakTables, remaining18BakTables...)

	if len(all114BakTables) != 114 {
		t.Fatalf("El catálogo oficial de World00_Character.bak debe sumar exactamente 114 tablas, suma %d", len(all114BakTables))
	}

	seenBak := make(map[string]bool, 114)
	for _, tbl := range all114BakTables {
		lw := strings.ToLower(tbl)
		if seenBak[lw] {
			t.Fatalf("Tabla duplicada en inventario de auditoría de 114 tablas: %s", tbl)
		}
		seenBak[lw] = true
		if _, ok := existingTables[lw]; !ok {
			t.Errorf("[Control 1] Falta tabla oficial de World00_Character.bak (114/114): %s", tbl)
		}
	}

	// 3.3 Las 2 Tablas Operativas S2S (WorldService <-> ZoneService)
	operationalS2STables := []string{"tActiveZoneTransfer", "tWorldZoneNode"}
	for _, tbl := range operationalS2STables {
		if _, ok := existingTables[strings.ToLower(tbl)]; !ok {
			t.Errorf("[Control 1] Falta tabla operativa S2S requerida: %s", tbl)
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

	// Invariante 4.11: num, tItem_deleted y tablas S2S
	checkColumns(t, db, schema, "num", 1, []string{"num"})
	checkColumns(t, db, schema, "tItem_deleted", 7, []string{
		"nitemkey", "nstoragetype", "nowner", "nstorage", "nitemid", "nflags", "ddate",
	})
	checkColumns(t, db, schema, "tActiveZoneTransfer", 14, []string{
		"scharid", "nworldno", "nzoneno", "nuserno", "suserid", "ncharno",
		"nrandomid", "nadminlevel", "sclientip", "stargetmap", "ntargetx", "ntargety",
		"dissuedat", "dexpiresat",
	})
	checkColumns(t, db, schema, "tWorldZoneNode", 7, []string{
		"nworldno", "nzoneno", "szoneip", "nzoneport", "ncurrentplayers", "bisonline", "dlastheartbeat",
	})

	// 5. Verificación de Semillas Oficiales e Higiene de tCharacter

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

	// 5.4 Verificar secuencia auxiliar num (0..1000 -> >= 1001 filas) y tItemStorageFindEmpty (0..255 -> >= 256 filas)
	var numCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM num").Scan(&numCount); err != nil {
		t.Fatalf("Error consultando num: %v", err)
	}
	if numCount < 1001 {
		t.Errorf("[Control 1 Semillas] Tabla num debe contener al menos 1001 enteros (0..1000), tiene %d", numCount)
	}

	var emptyStorageCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tItemStorageFindEmpty").Scan(&emptyStorageCount); err != nil {
		t.Fatalf("Error consultando tItemStorageFindEmpty: %v", err)
	}
	if emptyStorageCount < 256 {
		t.Errorf("[Control 1 Semillas] Tabla tItemStorageFindEmpty debe contener al menos 256 slots (0..255), tiene %d", emptyStorageCount)
	}

	// 5.5 PROHIBICIÓN ESTRICTA de filas falsas *_Template / Tutorial_* en tCharacter (Bug Crítico #1 resuelto)
	var fakeTemplateCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM tCharacter
		WHERE sID IN (
			'Fighter_Template', 'Cleric_Template', 'Archer_Template',
			'Mage_Template', 'Trickster_Template', 'Sentinel_Template',
			'Tutorial_Fighter', 'Tutorial_Cleric', 'Tutorial_Archer',
			'Tutorial_Mage', 'Tutorial_Trickster'
		)
	`).Scan(&fakeTemplateCount)
	if err != nil {
		t.Fatalf("Error verificando ausencia de templates falsos en tCharacter: %v", err)
	}
	if fakeTemplateCount != 0 {
		t.Errorf("[Control 1 Higiene] Prohibido contaminar tCharacter con filas *_Template o Tutorial_*: encontradas %d filas", fakeTemplateCount)
	}

	// 6. Verificación de Archivos Oficiales de Inicio (data/default_character_data.json y data/tutorial_character_data.json)
	dataDir := filepath.Join("..", "..", "data")

	defStore, err := world.LoadDefaultCharacterData(filepath.Join(dataDir, "default_character_data.json"))
	if err != nil {
		t.Fatalf("[Control 1 Data] Error cargando default_character_data.json: %v", err)
	}
	expectedDefaultClasses := []byte{1, 6, 11, 16, 21, 26}
	for _, clsID := range expectedDefaultClasses {
		clsData, ok := defStore.GetClassData(clsID)
		if !ok {
			t.Errorf("[Control 1 Data] Falta clase oficial %d en default_character_data.json", clsID)
			continue
		}
		if len(clsData.StarterItems) == 0 {
			t.Errorf("[Control 1 Data] Clase %d sin starter_items en default_character_data.json", clsID)
		}
	}
	if sentinel, ok := defStore.GetClassData(26); ok {
		if sentinel.StartMap != "Eld" || sentinel.StartX != 12259 || sentinel.StartY != 13434 || sentinel.InitLevel != 60 || sentinel.InitExp != 24346145 {
			t.Errorf("[Control 1 Data] Sentinel (26) con parámetros incorrectos: %+v", sentinel)
		}
	}
	if defStore.StarterQuests == nil {
		t.Errorf("[Control 1 Data] default_character_data.json debe incluir la clave starter_quests (#DEFINE QUEST)")
	}

	tutStore, err := world.LoadDefaultCharacterData(filepath.Join(dataDir, "tutorial_character_data.json"))
	if err != nil {
		t.Fatalf("[Control 1 Data] Error cargando tutorial_character_data.json: %v", err)
	}
	expectedTutorialClasses := []byte{1, 6, 11, 16, 21}
	for _, clsID := range expectedTutorialClasses {
		tutData, ok := tutStore.GetClassData(clsID)
		if !ok {
			t.Errorf("[Control 1 Data] Falta clase de tutorial %d en tutorial_character_data.json", clsID)
			continue
		}
		if tutData.StartMap != "RouN" || tutData.StartX != 7178 || tutData.StartY != 7011 || tutData.Money != 300 {
			t.Errorf("[Control 1 Data] Clase tutorial %d con StartMap/Coords/Money inesperado: map=%s (%d,%d) money=%d",
				clsID, tutData.StartMap, tutData.StartX, tutData.StartY, tutData.Money)
		}
		if len(tutData.StarterItems) == 0 || tutData.TutorialItems == nil {
			t.Errorf("[Control 1 Data] Clase tutorial %d debe tener starter_items (%d) y clave tutorial_items no nula",
				clsID, len(tutData.StarterItems))
		}
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
