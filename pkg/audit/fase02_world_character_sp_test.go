package audit_test

import (
	"database/sql"
	"strings"
	"testing"
)

// TestFase02_WorldCharacter_Control2 valida exhaustivamente el Control 2 de WorldService / Character (Fase 02A):
// 1. Presencia obligatoria del catálogo completo de 425 Stored Procedures (419 oficiales + 6 S2S) en information_schema.ROUTINES.
// 2. Garantía estructural Pool-Safe (proc_logic + SELECT final) en los procedimientos con parámetros OUT.
func TestFase02_WorldCharacter_Control2(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const schema = "xelion_online"

	// 1. Consultar todos los procedimientos almacenados existentes en information_schema
	rows, err := db.Query(`
		SELECT ROUTINE_NAME, IFNULL(ROUTINE_DEFINITION, '')
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE'
	`, schema)
	if err != nil {
		t.Fatalf("Error consultando information_schema.ROUTINES: %v", err)
	}
	defer rows.Close()

	existingProcs := make(map[string]string)
	for rows.Next() {
		var name, def string
		if err := rows.Scan(&name, &def); err != nil {
			t.Fatalf("Error escaneando ROUTINE_NAME: %v", err)
		}
		existingProcs[strings.ToLower(name)] = def
	}

	// 2. Catálogo Exhaustivo de 425 Stored Procedures de World00_Character (419 Oficiales + 6 S2S)
	target425 := []string{
		"p_Char_AddExp", "p_Char_ChangeClass", "p_Char_ChangeID", "p_Char_CheckOptionsExistance",
		"p_Char_Create", "p_Char_CreateSetDefaultData", "p_Char_DelChargedBuff", "p_Char_DelCharTitle",
		"p_Char_Delete", "p_Char_DelOptions", "p_Char_DelState", "p_Char_DelStateAll",
		"p_Char_FameAdd", "p_Char_FameGet", "p_Char_FameSet", "p_Char_FameSub",
		"p_Char_Find", "p_Char_GetAllData", "p_Char_GetChargedBufferAll", "p_Char_GetChargedBufferCount",
		"p_Char_GetCharID", "p_Char_GetCharNo", "p_Char_GetIsDeleted", "p_Char_GetKQMap",
		"p_Char_GetListOfUserChar", "p_Char_GetListOfUserCharAll", "p_Char_GetLoginData", "p_Char_GetNumOfUserChar",
		"p_Char_GetNumOfUserCharAll", "p_Char_GetOptGame", "p_Char_GetOptKeyMapping", "p_Char_GetOptShortCutData",
		"p_Char_GetOptShortCutSize", "p_Char_GetOptSound", "p_Char_GetOptVideo", "p_Char_GetOptWindowsPos",
		"p_Char_GetShape", "p_Char_GetState", "p_Char_GetStateAll", "p_Char_GetStateAllWithRestTime",
		"p_Char_GetStateAllWithSetItem", "p_Char_GetStateAllWithSetItemAndRestTime", "p_Char_GetStateWithRestTime", "p_Char_GetStateWithSetItem",
		"p_Char_GetStateWithSetItemAndRestTime", "p_Char_GetTitleAll", "p_Char_GetTitleCount", "p_Char_GetUserNo",
		"p_Char_Login", "p_Char_Logout", "p_Char_MoneyAdd", "p_Char_MoneyGet",
		"p_Char_MoneyGetCharsByUserNo", "p_Char_MoneyGetGreaterThen", "p_Char_MoneyMove", "p_Char_MoneySet",
		"p_Char_MoneySub", "p_Char_Rename", "p_Char_SaveLevelExpFrame", "p_Char_SaveLocation",
		"p_Char_SavePKCount", "p_Char_SaveStat", "p_Char_SetAdminLevel", "p_Char_SetAllData",
		"p_Char_SetChargedBuffer", "p_Char_SetLevelExp", "p_Char_SetOptGame", "p_Char_SetOptKeyMapping",
		"p_Char_SetOptShortCutData", "p_Char_SetOptShortCutSize", "p_Char_SetOptSound", "p_Char_SetOptVideo",
		"p_Char_SetOptWindowsPos", "p_Char_SetShape", "p_Char_SetState", "p_Char_SetStateWithRestTime",
		"p_Char_SetStateWithSetItem", "p_Char_SetStateWithSetItemAndRestTime", "p_Char_SetTitle", "p_Char_StatPoint_Dec",
		"p_Char_StatPoint_Inc", "p_Char_SubExp", "p_Char_Wedding_Cancel_Propose", "p_Char_Wedding_Divorce_by_Agree",
		"p_Char_Wedding_Divorce_by_Force", "p_Char_Wedding_Divorce_Cancel", "p_Char_Wedding_Divorce_Do", "p_Char_Wedding_Divorce_Do_by_Item",
		"p_Char_Wedding_Do", "p_Char_Wedding_Escape_Divorce_State", "p_Char_Wedding_Get_Partner_Info", "p_Char_Wedding_Propose",
		"p_Char_WeddingStatus_Get", "p_Char_WeddingStatus_Set", "p_Char_Wipe", "p_Event_Friend_Reward",
		"p_Event_WorldCup_Init", "p_Friend_Del", "p_Friend_Del_All", "p_Friend_Get",
		"p_Friend_GetAll", "p_Friend_GetAll2", "p_Friend_GetNumOfFriend", "p_Friend_Set",
		"p_Guild_Delete", "p_Guild_Get", "p_Guild_GetAll", "p_Guild_GetCount",
		"p_Guild_GetIntro", "p_Guild_GetList", "p_Guild_GetListDetail", "p_Guild_GetMyNo",
		"p_Guild_GetName", "p_Guild_GetNotify", "p_Guild_Make", "p_Guild_Make2",
		"p_Guild_Member_Get", "p_Guild_Member_GetAll", "p_Guild_Member_GetCount", "p_Guild_Member_Join",
		"p_Guild_Member_Leave", "p_Guild_Member_SetExpRatio", "p_Guild_Member_SetFlags", "p_Guild_Member_SetGrade",
		"p_Guild_Member_SetIntro", "p_Guild_MoneyAdd", "p_Guild_MoneySub", "p_Guild_Set",
		"p_Guild_SetDismiss", "p_Guild_SetExp", "p_Guild_SetFame", "p_Guild_SetGrade",
		"p_Guild_SetIntro", "p_Guild_SetMoney", "p_Guild_setName", "p_Guild_SetNotify",
		"p_Guild_SetStoneLevel", "p_Guild_SetType", "p_Guild_SetWar", "p_Guild_SetWarResult",
		"p_GuildTournament_Get", "p_GuildTournament_GetAllTournamentList", "p_GuildTournament_GetLastMatch", "p_GuildTournament_Set",
		"p_HolyPromise_CheckRing", "p_HolyPromise_Del", "p_HolyPromise_DelChar", "p_HolyPromise_Get",
		"p_HolyPromise_GetMember", "p_HolyPromise_GetMemberCount", "p_HolyPromise_Set", "p_HolyPromise_SetDate",
		"p_IsPromotionUser", "p_Item_BuyAll", "p_Item_BuyLot", "p_Item_ClrWeaponTitle",
		"p_Item_CostumWeapon_Endure_Set", "p_Item_CostumWeapon_Endure_Update", "p_Item_Create", "p_Item_Create_UseEmptyStorage",
		"p_Item_Create_UseEmptyStorage_tmp", "p_Item_Delete", "p_Item_ExchangeStorage", "p_Item_ExchangeStorageAndOwner",
		"p_Item_Furniture_Endure_Update", "p_Item_GetData", "p_Item_GetList", "p_Item_GetListType",
		"p_Item_GetMobList", "p_Item_GetMobListAll", "p_Item_GetNumOfItem", "p_Item_GetOption",
		"p_Item_GetOptionAll", "p_Item_GetOptionCount", "p_Item_GetStorage", "p_Item_GetStorageAndOwner",
		"p_Item_GetStorageEmpty", "p_Item_GetStorageEmptyList", "p_Item_GetStorageListType", "p_Item_GetStorageTypeCount",
		"p_Item_GetTitle", "p_Item_IsValidOwner", "p_Item_Merge", "p_Item_Merge_New",
		"p_Item_MiniHouse_Effect_Resource_Del", "p_Item_SellAll", "p_Item_SellAll_LotCheck", "p_Item_SellLot",
		"p_Item_SellLot_LotCheck", "p_Item_SetData", "p_Item_SetMobList", "p_Item_SetOption",
		"p_Item_SetStorage", "p_Item_SetStorageAndOwner", "p_Item_SetStorageAndOwner_UseEmptyStorage", "p_Item_SetTitle",
		"p_Item_SetWeaponTitle", "p_Item_Split_N_Merge", "p_Item_Split_N_Merge_New", "p_Item_Upgrade",
		"p_Item_Upgrade_CelllInAll", "p_Item_Upgrade_GOLD_NINE", "p_Item_Upgrade_OLD", "p_KingdomQuest_Del",
		"p_KingdomQuest_Get", "p_KingdomQuest_GetAll", "p_KingdomQuest_GetNumOfKQ", "p_KingdomQuest_Set",
		"p_MasPup_DelAll", "p_MasPup_DelMaster", "p_MasPup_DelPupils", "p_MasPup_GetAll",
		"p_MasPup_SetMaster", "p_MasPup_SetPupils", "p_MiniHouseURL_Get", "p_MiniHouseURL_Set",
		"p_Prison_Add", "p_Prison_Add2", "p_Prison_Alter", "p_Prison_Alter2",
		"p_Prison_End", "p_Prison_Get", "p_Prison_UpdateCharPrisonMin", "p_Promotion_Get",
		"p_Promotion_Set", "p_Quest_Add", "p_Quest_Del", "p_Quest_Get",
		"p_Quest_GetAll", "p_Quest_GetAllDoing", "p_Quest_GetAllDone", "p_Quest_Set",
		"p_Skill_Del", "p_Skill_Get", "p_Skill_GetAll", "p_Skill_GetNumOfSkill",
		"p_Skill_ResetPower", "p_Skill_Set", "p_Skill_SetCoolTime", "p_Skill_SetPassive",
		"p_Skill_SetPower", "p_Skill_SetPowerAll", "p_User_Deposit", "p_User_MoneyAdd",
		"p_User_MoneyGet", "p_User_MoneyGetGreaterThen", "p_User_MoneySet", "p_User_MoneySub",
		"p_User_Withdraw", "p_View_GetData", "usp_Attendance_getCallDateList", "usp_Attendance_getCount",
		"usp_Attendance_getRewardList", "usp_Attendance_setCall", "usp_Attendance_setReward", "usp_Char_chkCharNo",
		"usp_Char_GetAutoPick", "usp_Char_GetRestExpExecTime", "usp_Char_SetAutoPick", "usp_Char_SetRestExpExecTime",
		"usp_Character_chkKeyMap", "usp_Character_getAutoScroll", "usp_Character_getChangeIDList", "usp_Character_GetDiffFriendpoint",
		"usp_Character_getFriendpoint", "usp_Character_getGameOption", "usp_Character_getItemActionCoolTime", "usp_Character_getKeyMap",
		"usp_Character_getShortCutList", "usp_Character_getSingleOption", "usp_Character_initKeyMap", "usp_Character_setAutoPick",
		"usp_Character_setBuffScroll", "usp_Character_setChangeID", "usp_Character_setChangeIDItemUsed", "usp_Character_setChargedbuffScroll",
		"usp_Character_setChatBan", "usp_Character_setFriendpoint", "usp_Character_setGameOption", "usp_Character_setGameOptionInit",
		"usp_Character_setItemActionCoolTime", "usp_Character_setKeyMap", "usp_Character_setMysteryBolt", "usp_Character_setRestExpExecTime",
		"usp_Character_setShortCut", "usp_Character_setStatPoint", "usp_CharacterLoading", "usp_CharacterLoadingDetailInfo",
		"usp_ChargeItem_Draw", "usp_ChargeItem_GetList", "usp_ChargeItem_UseSkillPowerInitDel", "usp_ChargeItem_UseSkillPowerInitLotChange",
		"usp_ChargeItem_UseStatPointDecDel", "usp_ChargeItem_UseStatPointDecLotChange", "usp_ChargeItem_UseStatPointInitDel", "usp_ChargeItem_UseStatPointInitLotChange",
		"usp_ChatBlock_delBlock", "usp_ChatBlock_delBlockALL", "usp_ChatBlock_getBlockList", "usp_ChatBlock_setBlock",
		"usp_ChatColor_get", "usp_ChatColor_setBalloon", "usp_ChatColor_setFont", "usp_Collection_getCardList",
		"usp_Collection_getFavoritesList", "usp_Collection_getRewardList", "usp_Collection_getSerialList", "usp_Collection_setCard",
		"usp_Collection_setFavorites", "usp_Collection_setReward", "usp_Collection_setSerial", "usp_Dice_chgMachineToMoneycard",
		"usp_Dice_getDiceExpect", "usp_Dice_getDiceRank", "usp_Dice_setDice", "usp_Dice_setDiceLog",
		"usp_Dice_setDicePay", "usp_Dice_setDiceResult", "usp_Gamble_addCoin", "usp_Gamble_getCoin",
		"usp_Gamble_getDiceRank", "usp_Gamble_getExchangedCoin", "usp_Gamble_getJackpotCoin", "usp_Gamble_getJackpotRanking",
		"usp_Gamble_setCoinToMoney", "usp_Gamble_setDice", "usp_Gamble_setDiceLog", "usp_Gamble_setDiceResult",
		"usp_Gamble_setMoneyToCoin", "usp_Gamble_setSlotMachineResult", "usp_Gamble_subCoin", "usp_GM_getItemRewardGroupNo",
		"usp_GM_getItemRewardList", "usp_GM_setItemReward", "usp_Guild_getHistoryList", "usp_Guild_GetToken",
		"usp_Guild_GetTokenAll", "usp_Guild_setHistory", "usp_Guild_SetToken", "usp_Guild_TokenSub",
		"usp_GuildAcademy_delMember", "usp_GuildAcademy_delReward", "usp_GuildAcademy_delRewardAll", "usp_GuildAcademy_depositMoney",
		"usp_GuildAcademy_getGraduateCount", "usp_GuildAcademy_getGraduateList", "usp_GuildAcademy_getGuildAcademyList", "usp_GuildAcademy_getGuildNo",
		"usp_GuildAcademy_getHistoryList", "usp_GuildAcademy_getMemberList", "usp_GuildAcademy_getMoney", "usp_GuildAcademy_getRanking",
		"usp_GuildAcademy_getReward", "usp_GuildAcademy_getRewardAll", "usp_GuildAcademy_raiseMember", "usp_GuildAcademy_setAcademy",
		"usp_GuildAcademy_setChatBanTime", "usp_GuildAcademy_setGraduate", "usp_GuildAcademy_setHistory", "usp_GuildAcademy_setMaster",
		"usp_GuildAcademy_setMember", "usp_GuildAcademy_setMemberLevelup", "usp_GuildAcademy_setNotify", "usp_GuildAcademy_setRanking",
		"usp_GuildAcademy_setReward", "usp_GuildAcademy_withdrawMoney", "usp_GuildEmblem_GetInfo", "usp_GuildEmblem_GetInfoAll",
		"usp_GuildEmblem_SetInfo", "usp_GuildTournament_getChoiceList", "usp_GuildTournament_getMoney", "usp_GuildTournament_getRegistList",
		"usp_GuildTournament_setChoice", "usp_GuildTournament_setRefundNotice", "usp_GuildTournament_setRegist", "usp_GuildTournament_SetResult",
		"usp_HolyPromise_getRemainMoney", "usp_HolyPromise_setMoney", "usp_HolyPromise_withDrawRemainMoney", "usp_InGameGuide_getViewList",
		"usp_InGameGuide_setGuideNo", "usp_Item_Add_EnchantGem", "usp_Item_Add_EnchantSocket", "usp_Item_Del_EnchantGem",
		"usp_Item_Get_EnchantAll", "usp_Item_getAutoMonList", "usp_Item_getEmptyStorage", "usp_Item_GuildToken_BuyAll",
		"usp_Item_GuildToken_BuyLot", "usp_Item_Set_EnchantGemLot", "usp_Item_setAutoMon", "usp_Item_setPutOnBelonged",
		"usp_Item_setPutOnBelongedCancel", "usp_Item_setPutOnBelongedExchange", "usp_LinkScroll_SetDestMap", "usp_Minihouse_delGate",
		"usp_Minihouse_getGateList", "usp_Minihouse_getVisitorCount", "usp_Minihouse_setGate", "usp_Minihouse_setVisitorCount",
		"usp_Pet_getInfo", "usp_Pet_setCreate", "usp_Pet_setDelete", "usp_Pet_setName",
		"usp_Pet_setTendency", "usp_SafeTrade_setAll", "usp_SafeTrade_setLot", "usp_SafeTrade_setMerge",
		"usp_Session_ConsumeZoneTransfer", "usp_Session_IssueZoneTransfer", "usp_SpamBlock_getTime", "usp_SpamBlock_setTime",
		"usp_SpamFilter_get", "usp_SpamFilter_set", "usp_Statistics_getCharacterLevel", "usp_Statistics_setCharacterLevel",
		"usp_Statistics_setGameMoney", "usp_Statistics_setTotalMoney", "usp_Tutorial_getCharInfo", "usp_Tutorial_getIsSkip",
		"usp_Tutorial_setCharInfo", "usp_WorldZone_GetList", "usp_WorldZone_Heartbeat", "usp_WorldZone_ResolveMapEndpoint",
		"usp_WorldZone_UpsertNode",
	}

	if len(target425) != 425 {
		t.Fatalf("El catálogo exhaustivo de World00_Character debe contener exactamente 425 procedimientos, tiene %d", len(target425))
	}

	targetSet := make(map[string]string, 425)
	var missing []string
	for _, sp := range target425 {
		lw := strings.ToLower(sp)
		if _, dup := targetSet[lw]; dup {
			t.Fatalf("Procedimiento duplicado en catálogo target425: %s", sp)
		}
		targetSet[lw] = sp
		if _, ok := existingProcs[lw]; !ok {
			missing = append(missing, sp)
		}
	}

	if len(missing) > 0 {
		t.Fatalf("[Control 2] Faltan %d Stored Procedures en la base de datos: %v", len(missing), missing)
	}

	// 3. Auditar que todo SP de World00_Character con parámetros OUT esté blindado con proc_logic + SELECT final
	outRows, err := db.Query(`
		SELECT DISTINCT SPECIFIC_NAME
		FROM information_schema.PARAMETERS
		WHERE SPECIFIC_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' AND PARAMETER_MODE IN ('OUT', 'INOUT')
	`, schema)
	if err != nil {
		t.Fatalf("Error consultando information_schema.PARAMETERS: %v", err)
	}
	defer outRows.Close()

	for outRows.Next() {
		var spName string
		if err := outRows.Scan(&spName); err != nil {
			t.Fatalf("Error escaneando SPECIFIC_NAME: %v", err)
		}
		lw := strings.ToLower(spName)
		if _, isWorldCharSP := targetSet[lw]; !isWorldCharSP {
			continue
		}
		def := existingProcs[lw]
		if !strings.Contains(def, "proc_logic:") || !strings.Contains(strings.ToUpper(def), "END PROC_LOGIC") {
			t.Errorf("[Control 2 Pool-Safe] El procedimiento %s tiene parámetros OUT pero carece de bloque proc_logic", spName)
		}
		if strings.Contains(strings.ToUpper(def), "LEAVE SP_BODY") {
			t.Errorf("[Control 2 Pool-Safe] El procedimiento %s contiene LEAVE sp_body que salta el SELECT final", spName)
		}
	}

	t.Logf("Control 2 PASS: 425/425 Stored Procedures de World00_Character (419 oficiales + 6 S2S) verificados y pool-safe.")
}

// TestFase02_WorldCharacter_Functional valida la ejecución funcional de los Stored Procedures críticos
// sobre un pool *sql.DB, cubriendo tanto ramas de éxito como ramas negativas/tempranas (cero sql.ErrNoRows),
// asignación automática de GM (nAuthID=9 -> nAdminLevel=100) y coordinación S2S (ZoneTransfer y WorldZoneNode).
func TestFase02_WorldCharacter_Functional(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const (
		testUserNo      = 88888
		testAdminUserNo = 88889
		testCharID      = "TestCharF2Audit"
		testAdminCharID = "TestGMF2Audit"
		testWorldNo     = 99
		testZoneNo      = 7
	)

	// 1. Limpieza y preparación de usuarios de prueba
	cleanup := func() {
		for _, uNo := range []int{testUserNo, testAdminUserNo} {
			_, _ = db.Exec("DELETE FROM tItem_deleted WHERE nOwner IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", uNo)
			_, _ = db.Exec("DELETE FROM tItem WHERE nOwner IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", uNo)
			_, _ = db.Exec("DELETE FROM tCharacterKeyMap WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", uNo)
			_, _ = db.Exec("DELETE FROM tCharacterShape WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", uNo)
			_, _ = db.Exec("DELETE FROM tCharacterOptions WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", uNo)
			_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", uNo)
			_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", uNo)
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestCharF2Audit%' OR sID LIKE 'TestGMF2Audit%'")
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE sCharID = ?", testCharID)
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nWorldNo = ?", testWorldNo)
	}
	cleanup()
	defer cleanup()

	// Asegurar existencia de usuario normal (nAuthID=1) y usuario GM (nAuthID=9) en tUser
	_, err := db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, sUserIP, nAuthID)
		VALUES
			(?, 'testuser88888', 'testpass123', 'Tester', '127.0.0.1', 1),
			(?, 'testgm88889', 'testpass123', 'TesterGM', '127.0.0.1', 9)
		ON DUPLICATE KEY UPDATE nAuthID = VALUES(nAuthID)
	`, testUserNo, testAdminUserNo)
	if err != nil {
		t.Fatalf("Error creando tUser para test funcional: %v", err)
	}

	// 2. Probar p_Char_Create (Éxito)
	var createdCharNo int
	row := db.QueryRow(`
		CALL p_Char_Create(?, 1, 0, 0, ?, 1, 1, 1, 1, 1, 1, @out_charNo)
	`, testUserNo, testCharID)
	if err := row.Scan(&createdCharNo); err != nil {
		t.Fatalf("Error ejecutando p_Char_Create: %v", err)
	}
	if createdCharNo <= 0 {
		t.Fatalf("p_Char_Create retornó nCharNo inválido: %d", createdCharNo)
	}

	// 2.1 Probar rama negativa Pool-Safe: Slot ocupado en p_Char_Create -> retorna -1 vía Scan (JAMÁS sql.ErrNoRows)
	var dupSlotRet int
	row = db.QueryRow(`
		CALL p_Char_Create(?, 1, 0, 0, 'OtherCharSlot0', 1, 1, 1, 1, 1, 1, @out_charNo)
	`, testUserNo)
	if err := row.Scan(&dupSlotRet); err != nil {
		t.Fatalf("[Pool-Safe Bug] p_Char_Create en slot ocupado devolvió error (%v) en vez de resultset con -1", err)
	}
	if dupSlotRet != -1 {
		t.Errorf("p_Char_Create en slot ocupado retornó %d, esperado -1", dupSlotRet)
	}

	// 2.2 Probar rama negativa Pool-Safe: sID duplicado en p_Char_Create -> retorna -2 vía Scan (JAMÁS sql.ErrNoRows)
	var dupNameRet int
	row = db.QueryRow(`
		CALL p_Char_Create(?, 1, 0, 1, ?, 1, 1, 1, 1, 1, 1, @out_charNo)
	`, testUserNo, testCharID)
	if err := row.Scan(&dupNameRet); err != nil {
		t.Fatalf("[Pool-Safe Bug] p_Char_Create con sID duplicado devolvió error (%v) en vez de resultset con -2", err)
	}
	if dupNameRet != -2 {
		t.Errorf("p_Char_Create con sID duplicado retornó %d, esperado -2", dupNameRet)
	}

	// 2.3 Probar resolución Cross-DB tUser.nAuthID = 9 -> nAdminLevel = 100 en tCharacter
	var gmCharNo int
	if err := db.QueryRow(`
		CALL p_Char_Create(?, 1, 0, 0, ?, 1, 1, 1, 1, 1, 1, @out_charNo)
	`, testAdminUserNo, testAdminCharID).Scan(&gmCharNo); err != nil {
		t.Fatalf("Error creando personaje GM con nAuthID=9: %v", err)
	}
	var gmAdminLevel int
	if err := db.QueryRow("SELECT nAdminLevel FROM tCharacter WHERE nCharNo = ?", gmCharNo).Scan(&gmAdminLevel); err != nil {
		t.Fatalf("Error consultando nAdminLevel del personaje GM: %v", err)
	}
	if gmAdminLevel != 100 {
		t.Errorf("p_Char_Create para usuario con nAuthID=9 asignó nAdminLevel=%d, esperado 100", gmAdminLevel)
	}

	// 3. Probar p_Char_Find y p_Char_GetCharNo
	var foundCharNo int
	if err := db.QueryRow(`CALL p_Char_Find(?, @out_charNo)`, testCharID).Scan(&foundCharNo); err != nil {
		t.Fatalf("Error ejecutando p_Char_Find: %v", err)
	}
	if foundCharNo != createdCharNo {
		t.Errorf("p_Char_Find retornó %d, esperado %d", foundCharNo, createdCharNo)
	}

	var foundCharNo2 int
	if err := db.QueryRow(`CALL p_Char_GetCharNo(?, @out_charNo)`, testCharID).Scan(&foundCharNo2); err != nil {
		t.Fatalf("Error ejecutando p_Char_GetCharNo: %v", err)
	}
	if foundCharNo2 != createdCharNo {
		t.Errorf("p_Char_GetCharNo retornó %d, esperado %d", foundCharNo2, createdCharNo)
	}

	// 4. Probar p_Char_GetListOfUserChar
	rows, err := db.Query(`CALL p_Char_GetListOfUserChar(?)`, testUserNo)
	if err != nil {
		t.Fatalf("Error ejecutando p_Char_GetListOfUserChar: %v", err)
	}
	var charList []int
	for rows.Next() {
		var cNo int
		if err := rows.Scan(&cNo); err != nil {
			rows.Close()
			t.Fatalf("Error escaneando lista de personajes: %v", err)
		}
		charList = append(charList, cNo)
	}
	rows.Close()
	if len(charList) != 1 || charList[0] != createdCharNo {
		t.Errorf("p_Char_GetListOfUserChar retornó %v, esperado [%d]", charList, createdCharNo)
	}

	// 5. Probar p_Char_GetLoginData
	var (
		sID, sLoginZone string
		nLevel, nSlotNo uint8
		bDeleted        uint8
		dDeletedDate    sql.NullString
		nRace, nClass   uint8
		nGender         uint8
		nHairType       uint8
		nHairColor      uint8
		nFaceShape      uint8
		nKQHandle       uint16
		sKQMap          string
		nKQX, nKQY      uint32
		nKQD            uint8
		dKQDate         sql.NullString
	)

	if err := db.QueryRow(`CALL p_Char_GetLoginData(?)`, createdCharNo).Scan(
		&sID, &nLevel, &nSlotNo, &sLoginZone, &bDeleted, &dDeletedDate,
		&nRace, &nClass, &nGender, &nHairType, &nHairColor, &nFaceShape,
		&nKQHandle, &sKQMap, &nKQX, &nKQY, &nKQD, &dKQDate,
	); err != nil {
		t.Fatalf("Error ejecutando p_Char_GetLoginData: %v", err)
	}
	if sID != testCharID || nClass != 1 || bDeleted != 0 {
		t.Errorf("p_Char_GetLoginData retornó datos inesperados: sID=%s, nClass=%d, bDeleted=%d", sID, nClass, bDeleted)
	}

	// 6. Probar p_Char_Login y p_Char_Logout
	var retLogin int
	if err := db.QueryRow(`CALL p_Char_Login(?, @out_ret)`, createdCharNo).Scan(&retLogin); err != nil {
		t.Fatalf("Error ejecutando p_Char_Login: %v", err)
	}

	var retLogout int
	if err := db.QueryRow(`CALL p_Char_Logout(?, 20, @out_ret)`, createdCharNo).Scan(&retLogout); err != nil {
		t.Fatalf("Error ejecutando p_Char_Logout: %v", err)
	}
	if retLogout != 1 {
		t.Errorf("p_Char_Logout retornó %d, esperado 1", retLogout)
	}

	// 7. Probar p_Item_Create, p_Item_GetStorageEmpty, p_Item_GetListType y p_Item_Delete (con auditoría tItem_deleted)
	const testItemKey = int64(9876543210)
	var retItemCreate int
	if err := db.QueryRow(`CALL p_Item_Create(?, ?, 1, 0, 1005, 0, @out_ret)`, testItemKey, createdCharNo).Scan(&retItemCreate); err != nil {
		t.Fatalf("Error ejecutando p_Item_Create: %v", err)
	}
	if retItemCreate != 0 {
		t.Errorf("p_Item_Create retornó %d, esperado 0", retItemCreate)
	}

	var nextEmptySlot int
	if err := db.QueryRow(`CALL p_Item_GetStorageEmpty(?, 1, @out_ret)`, createdCharNo).Scan(&nextEmptySlot); err != nil {
		t.Fatalf("Error ejecutando p_Item_GetStorageEmpty: %v", err)
	}
	if nextEmptySlot != 1 {
		t.Errorf("p_Item_GetStorageEmpty con slot 0 ocupado retornó %d, esperado 1", nextEmptySlot)
	}

	itemRows, err := db.Query(`CALL p_Item_GetListType(?, 1)`, createdCharNo)
	if err != nil {
		t.Fatalf("Error ejecutando p_Item_GetListType: %v", err)
	}
	var itemCount int
	for itemRows.Next() {
		var iKey int64
		var storage int16
		var itemID, flags int
		var dDate string
		if err := itemRows.Scan(&iKey, &storage, &itemID, &flags, &dDate); err != nil {
			itemRows.Close()
			t.Fatalf("Error escaneando fila de p_Item_GetListType: %v", err)
		}
		if iKey == testItemKey && itemID == 1005 {
			itemCount++
		}
	}
	itemRows.Close()
	if itemCount != 1 {
		t.Errorf("p_Item_GetListType no encontró el item creado (%d)", testItemKey)
	}

	var retItemDel int
	if err := db.QueryRow(`CALL p_Item_Delete(?, @out_ret)`, testItemKey).Scan(&retItemDel); err != nil {
		t.Fatalf("Error ejecutando p_Item_Delete: %v", err)
	}
	var archivedItemCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tItem_deleted WHERE nItemKey = ?", testItemKey).Scan(&archivedItemCount); err != nil {
		t.Fatalf("Error verificando tItem_deleted: %v", err)
	}
	if archivedItemCount != 1 {
		t.Errorf("p_Item_Delete debe archivar el ítem en tItem_deleted, encontrado %d", archivedItemCount)
	}

	// 8. Probar coordinación S2S (usp_Session_IssueZoneTransfer / ConsumeZoneTransfer y usp_WorldZone_*)
	var issueRet int
	if err := db.QueryRow(`
		CALL usp_Session_IssueZoneTransfer(?, ?, ?, 'testuser88888', ?, ?, 4321, 0, '127.0.0.1', 'Rou', 5625, 4505)
	`, testWorldNo, testZoneNo, testUserNo, createdCharNo, testCharID).Scan(&issueRet); err != nil {
		t.Fatalf("Error ejecutando usp_Session_IssueZoneTransfer: %v", err)
	}
	if issueRet != 0 {
		t.Errorf("usp_Session_IssueZoneTransfer retornó %d, esperado 0", issueRet)
	}

	var (
		ztWorldNo   uint8
		ztZoneNo    uint16
		ztUserNo    int
		ztUserID    string
		ztCharNo    int
		ztCharID    string
		ztAdminLv   uint8
		ztTargetMap string
		ztTargetX   int
		ztTargetY   int
		ztRet       int
	)
	// Consumo con RandomID inválido -> -2 sin destruir el token
	if err := db.QueryRow(`CALL usp_Session_ConsumeZoneTransfer(?, 9999, '127.0.0.1')`, testCharID).Scan(
		&ztWorldNo, &ztZoneNo, &ztUserNo, &ztUserID, &ztCharNo, &ztCharID, &ztAdminLv, &ztTargetMap, &ztTargetX, &ztTargetY, &ztRet,
	); err != nil {
		t.Fatalf("Error consumiendo ZoneTransfer con RandomID inválido: %v", err)
	}
	if ztRet != -2 {
		t.Errorf("usp_Session_ConsumeZoneTransfer con RandomID erróneo retornó %d, esperado -2", ztRet)
	}

	// Consumo válido de un solo uso -> 0
	if err := db.QueryRow(`CALL usp_Session_ConsumeZoneTransfer(?, 4321, '127.0.0.1')`, testCharID).Scan(
		&ztWorldNo, &ztZoneNo, &ztUserNo, &ztUserID, &ztCharNo, &ztCharID, &ztAdminLv, &ztTargetMap, &ztTargetX, &ztTargetY, &ztRet,
	); err != nil {
		t.Fatalf("Error consumiendo ZoneTransfer válido: %v", err)
	}
	if ztRet != 0 || ztCharNo != createdCharNo || ztTargetMap != "Rou" || ztTargetX != 5625 || ztTargetY != 4505 {
		t.Errorf("Payload inesperado en usp_Session_ConsumeZoneTransfer: ret=%d charNo=%d map=%s (%d,%d)", ztRet, ztCharNo, ztTargetMap, ztTargetX, ztTargetY)
	}

	// Segundo consumo del mismo token (replay) -> -1 sin sql.ErrNoRows
	if err := db.QueryRow(`CALL usp_Session_ConsumeZoneTransfer(?, 4321, '127.0.0.1')`, testCharID).Scan(
		&ztWorldNo, &ztZoneNo, &ztUserNo, &ztUserID, &ztCharNo, &ztCharID, &ztAdminLv, &ztTargetMap, &ztTargetX, &ztTargetY, &ztRet,
	); err != nil {
		t.Fatalf("Segundo usp_Session_ConsumeZoneTransfer devolvió error (%v) en lugar de -1", err)
	}
	if ztRet != -1 {
		t.Errorf("Replay de usp_Session_ConsumeZoneTransfer retornó %d, esperado -1", ztRet)
	}

	// Registro y resolución de nodo de zona (usp_WorldZone_*)
	var wzRet int
	if err := db.QueryRow(`CALL usp_WorldZone_UpsertNode(?, ?, '127.0.0.1', 9123)`, testWorldNo, testZoneNo).Scan(&wzRet); err != nil {
		t.Fatalf("Error ejecutando usp_WorldZone_UpsertNode: %v", err)
	}
	if err := db.QueryRow(`CALL usp_WorldZone_Heartbeat(?, ?, 42)`, testWorldNo, testZoneNo).Scan(&wzRet); err != nil {
		t.Fatalf("Error ejecutando usp_WorldZone_Heartbeat: %v", err)
	}
	var resIP string
	var resPort uint16
	var resOnline int
	if err := db.QueryRow(`CALL usp_WorldZone_ResolveMapEndpoint(?, ?)`, testWorldNo, testZoneNo).Scan(&resIP, &resPort, &resOnline); err != nil {
		t.Fatalf("Error ejecutando usp_WorldZone_ResolveMapEndpoint: %v", err)
	}
	if resIP != "127.0.0.1" || resPort != 9123 || resOnline != 1 {
		t.Errorf("usp_WorldZone_ResolveMapEndpoint retornó (%s, %d, %d), esperado (127.0.0.1, 9123, 1)", resIP, resPort, resOnline)
	}

	// 9. Probar p_Char_Delete (Soft delete oficial + ramas negativas Pool-Safe)
	var retCharDel int
	if err := db.QueryRow(`CALL p_Char_Delete(?, @out_ret)`, createdCharNo).Scan(&retCharDel); err != nil {
		t.Fatalf("Error ejecutando p_Char_Delete: %v", err)
	}
	if retCharDel != 0 {
		t.Errorf("p_Char_Delete retornó %d, esperado 0", retCharDel)
	}

	// Rama negativa Pool-Safe 1: Borrar personaje ya borrado -> retorna 1 vía Scan (JAMÁS sql.ErrNoRows)
	var retAlreadyDel int
	if err := db.QueryRow(`CALL p_Char_Delete(?, @out_ret)`, createdCharNo).Scan(&retAlreadyDel); err != nil {
		t.Fatalf("[Pool-Safe Bug] p_Char_Delete sobre personaje ya eliminado devolvió error (%v)", err)
	}
	if retAlreadyDel != 1 {
		t.Errorf("p_Char_Delete sobre personaje ya eliminado retornó %d, esperado 1", retAlreadyDel)
	}

	// Rama negativa Pool-Safe 2: Borrar personaje inexistente -> retorna 2 vía Scan (JAMÁS sql.ErrNoRows)
	var retNotFoundDel int
	if err := db.QueryRow(`CALL p_Char_Delete(-999999, @out_ret)`).Scan(&retNotFoundDel); err != nil {
		t.Fatalf("[Pool-Safe Bug] p_Char_Delete sobre personaje inexistente devolvió error (%v)", err)
	}
	if retNotFoundDel != 2 {
		t.Errorf("p_Char_Delete sobre personaje inexistente retornó %d, esperado 2", retNotFoundDel)
	}

	// Verificar en tCharacter que bDeleted = 1 y sID tiene el tag ~#0
	var delFlag uint8
	var deletedSID string
	if err := db.QueryRow("SELECT bDeleted, sID FROM tCharacter WHERE nCharNo = ?", createdCharNo).Scan(&delFlag, &deletedSID); err != nil {
		t.Fatalf("Error consultando tCharacter tras delete: %v", err)
	}
	if delFlag != 1 || !strings.Contains(deletedSID, "~") {
		t.Errorf("Soft-delete incorrecto en tCharacter: bDeleted=%d, sID=%s", delFlag, deletedSID)
	}

	// Verificar que p_Char_GetListOfUserChar ya no retorne el personaje eliminado
	rowsAfterDel, err := db.Query(`CALL p_Char_GetListOfUserChar(?)`, testUserNo)
	if err != nil {
		t.Fatalf("Error ejecutando p_Char_GetListOfUserChar tras delete: %v", err)
	}
	defer rowsAfterDel.Close()

	var listAfterDel []int
	for rowsAfterDel.Next() {
		var cNo int
		if err := rowsAfterDel.Scan(&cNo); err != nil {
			t.Fatalf("Error escaneando lista tras delete: %v", err)
		}
		listAfterDel = append(listAfterDel, cNo)
	}
	if len(listAfterDel) != 0 {
		t.Errorf("p_Char_GetListOfUserChar retornó personajes tras soft-delete: %v", listAfterDel)
	}

	t.Logf("Flujos funcionales Pool-Safe (positivos, negativos, GM y S2S) ejecutados con 100%% de éxito.")
}
