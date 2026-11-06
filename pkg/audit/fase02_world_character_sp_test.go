package audit_test

import (
	"database/sql"
	"strings"
	"testing"
)

// TestFase02_WorldCharacter_Control2 valida exhaustivamente el Control 2 de WorldService / Character:
// 1. Presencia obligatoria del catálogo completo de 311 Stored Procedures oficiales en information_schema.ROUTINES.
// 2. Fidelidad funcional de ejecución y lógica relacional sobre MySQL 8.4 para los flujos nucleares.
func TestFase02_WorldCharacter_Control2(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const schema = "xelion_online"

	// 1. Consultar todos los procedimientos almacenados existentes en information_schema
	rows, err := db.Query(`
		SELECT ROUTINE_NAME 
		FROM information_schema.ROUTINES 
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE'
	`, schema)
	if err != nil {
		t.Fatalf("Error consultando information_schema.ROUTINES: %v", err)
	}
	defer rows.Close()

	existingProcs := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Error escaneando ROUTINE_NAME: %v", err)
		}
		existingProcs[strings.ToLower(name)] = true
	}

	// 2. Catálogo Oficial Exhaustivo de 311 Stored Procedures de World00_Character
	target311 := []string{
		"p_Char_AddExp", "p_Char_ChangeClass", "p_Char_ChangeID", "p_Char_CheckOptionsExistance",
		"p_Char_Create", "p_Char_CreateSetDefaultData", "p_Char_DelChargedBuff", "p_Char_DelOptions",
		"p_Char_DelState", "p_Char_DelStateAll", "p_Char_Delete", "p_Char_FameAdd",
		"p_Char_FameGet", "p_Char_FameSet", "p_Char_FameSub", "p_Char_Find",
		"p_Char_GetAllData", "p_Char_GetCharID", "p_Char_GetCharNo", "p_Char_GetChargedBufferAll",
		"p_Char_GetChargedBufferCount", "p_Char_GetIsDeleted", "p_Char_GetKQMap", "p_Char_GetListOfUserChar",
		"p_Char_GetListOfUserCharAll", "p_Char_GetLoginData", "p_Char_GetNumOfUserChar", "p_Char_GetNumOfUserCharAll",
		"p_Char_GetOptGame", "p_Char_GetOptKeyMapping", "p_Char_GetOptShortCutData", "p_Char_GetOptShortCutSize",
		"p_Char_GetOptSound", "p_Char_GetOptVideo", "p_Char_GetOptWindowsPos", "p_Char_GetShape",
		"p_Char_GetState", "p_Char_GetStateAll", "p_Char_GetStateAllWithSetItemAndRestTime", "p_Char_GetTitleAll",
		"p_Char_GetTitleCount", "p_Char_GetUserNo", "p_Char_Login", "p_Char_Logout",
		"p_Char_MoneyAdd", "p_Char_MoneyGet", "p_Char_MoneyMove", "p_Char_MoneySub",
		"p_Char_Rename", "p_Char_SaveLevelExpFrame", "p_Char_SaveLocation", "p_Char_SavePKCount",
		"p_Char_SaveStat", "p_Char_SetChargedBuffer", "p_Char_SetLevelExp", "p_Char_SetOptGame",
		"p_Char_SetOptKeyMapping", "p_Char_SetOptShortCutData", "p_Char_SetOptShortCutSize", "p_Char_SetOptSound",
		"p_Char_SetOptVideo", "p_Char_SetOptWindowsPos", "p_Char_SetShape", "p_Char_SetState",
		"p_Char_SetStateWithSetItem", "p_Char_SetStateWithSetItemAndRestTime", "p_Char_SetTitle", "p_Char_StatPoint_Dec",
		"p_Char_StatPoint_Inc", "p_Char_SubExp", "p_Char_WeddingStatus_Get", "p_Char_WeddingStatus_Set",
		"p_Char_Wedding_Cancel_Propose", "p_Char_Wedding_Divorce_Cancel", "p_Char_Wedding_Divorce_Do", "p_Char_Wedding_Divorce_Do_by_Item",
		"p_Char_Wedding_Divorce_by_Agree", "p_Char_Wedding_Divorce_by_Force", "p_Char_Wedding_Do", "p_Char_Wedding_Escape_Divorce_State",
		"p_Char_Wedding_Get_Partner_Info", "p_Char_Wedding_Propose", "p_Char_Wipe", "p_Friend_Del",
		"p_Friend_Del_All", "p_Friend_Get", "p_Friend_GetAll", "p_Friend_GetNumOfFriend",
		"p_Friend_Set", "p_GuildTournament_GetLastMatch", "p_GuildTournament_Set", "p_Guild_Delete",
		"p_Guild_Get", "p_Guild_GetAll", "p_Guild_GetMyNo", "p_Guild_Make2",
		"p_Guild_Member_GetAll", "p_Guild_Member_Join", "p_Guild_Member_Leave", "p_Guild_Member_SetGrade",
		"p_Guild_Member_SetIntro", "p_Guild_MoneyAdd", "p_Guild_MoneySub", "p_Guild_SetDismiss",
		"p_Guild_SetExp", "p_Guild_SetFame", "p_Guild_SetIntro", "p_Guild_SetNotify",
		"p_Guild_SetType", "p_Guild_SetWar", "p_Guild_SetWarResult", "p_Guild_setName",
		"p_HolyPromise_Del", "p_HolyPromise_DelChar", "p_HolyPromise_Get", "p_HolyPromise_GetMember",
		"p_HolyPromise_Set", "p_HolyPromise_SetDate", "p_Item_BuyAll", "p_Item_BuyLot",
		"p_Item_ClrWeaponTitle", "p_Item_CostumWeapon_Endure_Set", "p_Item_CostumWeapon_Endure_Update", "p_Item_Create",
		"p_Item_Create_UseEmptyStorage", "p_Item_Delete", "p_Item_ExchangeStorage", "p_Item_ExchangeStorageAndOwner",
		"p_Item_Furniture_Endure_Update", "p_Item_GetData", "p_Item_GetListType", "p_Item_GetMobListAll",
		"p_Item_GetOption", "p_Item_GetOptionAll", "p_Item_GetStorageAndOwner", "p_Item_GetStorageEmptyList",
		"p_Item_GetStorageListType", "p_Item_GetTitle", "p_Item_IsValidOwner", "p_Item_Merge_New",
		"p_Item_SellAll_LotCheck", "p_Item_SellLot_LotCheck", "p_Item_SetMobList", "p_Item_SetOption",
		"p_Item_SetStorage", "p_Item_SetStorageAndOwner", "p_Item_SetStorageAndOwner_UseEmptyStorage", "p_Item_SetTitle",
		"p_Item_SetWeaponTitle", "p_Item_Split_N_Merge_New", "p_Item_Upgrade_CelllInAll", "p_Prison_Add",
		"p_Prison_Add2", "p_Prison_Alter", "p_Prison_Alter2", "p_Prison_Get",
		"p_Prison_UpdateCharPrisonMin", "p_Quest_Del", "p_Quest_GetAllDoing", "p_Quest_Set",
		"p_Skill_Del", "p_Skill_Get", "p_Skill_GetAll", "p_Skill_ResetPower",
		"p_Skill_Set", "p_Skill_SetPassive", "p_Skill_SetPowerAll", "p_User_Deposit",
		"p_User_MoneyGet", "p_User_Withdraw", "usp_Attendance_getCallDateList", "usp_Attendance_getCount",
		"usp_Attendance_getRewardList", "usp_Attendance_setCall", "usp_Attendance_setReward", "usp_Char_chkCharNo",
		"usp_CharacterLoading", "usp_Character_chkKeyMap", "usp_Character_getAutoScroll", "usp_Character_getChangeIDList",
		"usp_Character_getGameOption", "usp_Character_getItemActionCoolTime", "usp_Character_getKeyMap", "usp_Character_getShortCutList",
		"usp_Character_getSingleOption", "usp_Character_initKeyMap", "usp_Character_setAutoPick", "usp_Character_setBuffScroll",
		"usp_Character_setChangeID", "usp_Character_setChangeIDItemUsed", "usp_Character_setChargedbuffScroll", "usp_Character_setChatBan",
		"usp_Character_setGameOption", "usp_Character_setGameOptionInit", "usp_Character_setItemActionCoolTime", "usp_Character_setKeyMap",
		"usp_Character_setMysteryBolt", "usp_Character_setRestExpExecTime", "usp_Character_setShortCut", "usp_Character_setStatPoint",
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
		"usp_Gamble_setMoneyToCoin", "usp_Gamble_setSlotMachineResult", "usp_GuildAcademy_delMember", "usp_GuildAcademy_delRewardAll",
		"usp_GuildAcademy_depositMoney", "usp_GuildAcademy_getGraduateCount", "usp_GuildAcademy_getGraduateList", "usp_GuildAcademy_getGuildAcademyList",
		"usp_GuildAcademy_getGuildNo", "usp_GuildAcademy_getHistoryList", "usp_GuildAcademy_getMemberList", "usp_GuildAcademy_getMoney",
		"usp_GuildAcademy_getRanking", "usp_GuildAcademy_getReward", "usp_GuildAcademy_getRewardAll", "usp_GuildAcademy_raiseMember",
		"usp_GuildAcademy_setChatBanTime", "usp_GuildAcademy_setGraduate", "usp_GuildAcademy_setHistory", "usp_GuildAcademy_setMaster",
		"usp_GuildAcademy_setMember", "usp_GuildAcademy_setMemberLevelup", "usp_GuildAcademy_setNotify", "usp_GuildAcademy_setRanking",
		"usp_GuildAcademy_setReward", "usp_GuildAcademy_withdrawMoney", "usp_GuildEmblem_GetInfo", "usp_GuildEmblem_GetInfoAll",
		"usp_GuildEmblem_SetInfo", "usp_GuildTournament_getChoiceList", "usp_GuildTournament_getMoney", "usp_GuildTournament_getRegistList",
		"usp_GuildTournament_setChoice", "usp_GuildTournament_setRegist", "usp_Guild_GetToken", "usp_Guild_GetTokenAll",
		"usp_Guild_SetToken", "usp_HolyPromise_getRemainMoney", "usp_HolyPromise_setMoney", "usp_HolyPromise_withDrawRemainMoney",
		"usp_InGameGuide_getViewList", "usp_InGameGuide_setGuideNo", "usp_Item_Add_EnchantGem", "usp_Item_Add_EnchantSocket",
		"usp_Item_Del_EnchantGem", "usp_Item_Get_EnchantAll", "usp_Item_GuildToken_BuyAll", "usp_Item_GuildToken_BuyLot",
		"usp_Item_Set_EnchantGemLot", "usp_Item_getAutoMonList", "usp_Item_setAutoMon", "usp_Item_setPutOnBelonged",
		"usp_Item_setPutOnBelongedCancel", "usp_Item_setPutOnBelongedExchange", "usp_Minihouse_delGate", "usp_Minihouse_getGateList",
		"usp_Minihouse_setGate", "usp_Minihouse_setVisitorCount", "usp_Pet_getInfo", "usp_Pet_setCreate",
		"usp_Pet_setName", "usp_SafeTrade_setAll", "usp_SafeTrade_setLot", "usp_SafeTrade_setMerge",
		"usp_SpamBlock_getTime", "usp_SpamBlock_setTime", "usp_SpamFilter_get", "usp_SpamFilter_set",
		"usp_Tutorial_getCharInfo", "usp_Tutorial_getIsSkip", "usp_Tutorial_setCharInfo",
	}

	if len(target311) != 311 {
		t.Fatalf("El catálogo oficial de World00_Character debe contener exactamente 311 procedimientos, tiene %d", len(target311))
	}

	var missing []string
	for _, sp := range target311 {
		if !existingProcs[strings.ToLower(sp)] {
			missing = append(missing, sp)
		}
	}

	if len(missing) > 0 {
		t.Fatalf("[Control 2] Faltan %d Stored Procedures oficiales en la base de datos: %v", len(missing), missing)
	}

	t.Logf("Control 2 PASS: 311/311 Stored Procedures oficiales de World00_Character verificados en information_schema.")
}

// TestFase02_WorldCharacter_Functional valida la ejecución funcional de los Stored Procedures críticos
// asegurando el ciclo de vida de avatares, inventario y persistencia relacional en MySQL 8.4.
func TestFase02_WorldCharacter_Functional(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	const (
		testUserNo = 88888
		testCharID = "TestCharF2Audit"
	)

	// 1. Limpieza y preparación de usuario de prueba
	cleanup := func() {
		_, _ = db.Exec("DELETE FROM tItem WHERE nOwner IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testUserNo)
		_, _ = db.Exec("DELETE FROM tCharacterShape WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testUserNo)
		_, _ = db.Exec("DELETE FROM tCharacterOptions WHERE nCharNo IN (SELECT nCharNo FROM tCharacter WHERE nUserNo = ?)", testUserNo)
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ? OR sID LIKE 'TestCharF2Audit%'", testUserNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", testUserNo)
	}
	cleanup()
	defer cleanup()

	// Asegurar existencia de usuario en tUser para satisfacer Foreign Key
	_, err := db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, sUserIP)
		VALUES (?, 'testuser88888', 'testpass123', 'Tester', '127.0.0.1')
		ON DUPLICATE KEY UPDATE sUserID = VALUES(sUserID)
	`, testUserNo)
	if err != nil {
		t.Fatalf("Error creando tUser para test funcional: %v", err)
	}

	// 2. Probar p_Char_Create
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
	t.Logf("p_Char_Create OK -> nCharNo: %d", createdCharNo)

	// 3. Probar p_Char_Find y p_Char_GetCharNo
	var foundCharNo int
	row = db.QueryRow(`CALL p_Char_Find(?, @out_charNo)`, testCharID)
	if err := row.Scan(&foundCharNo); err != nil {
		t.Fatalf("Error ejecutando p_Char_Find: %v", err)
	}
	if foundCharNo != createdCharNo {
		t.Errorf("p_Char_Find retornó %d, esperado %d", foundCharNo, createdCharNo)
	}

	var foundCharNo2 int
	row = db.QueryRow(`CALL p_Char_GetCharNo(?, @out_charNo)`, testCharID)
	if err := row.Scan(&foundCharNo2); err != nil {
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
	defer rows.Close()

	var charList []int
	for rows.Next() {
		var cNo int
		if err := rows.Scan(&cNo); err != nil {
			t.Fatalf("Error escaneando lista de personajes: %v", err)
		}
		charList = append(charList, cNo)
	}
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

	row = db.QueryRow(`CALL p_Char_GetLoginData(?)`, createdCharNo)
	if err := row.Scan(
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
	row = db.QueryRow(`CALL p_Char_Login(?, @out_ret)`, createdCharNo)
	if err := row.Scan(&retLogin); err != nil {
		t.Fatalf("Error ejecutando p_Char_Login: %v", err)
	}

	var retLogout int
	row = db.QueryRow(`CALL p_Char_Logout(?, 20, @out_ret)`, createdCharNo)
	if err := row.Scan(&retLogout); err != nil {
		t.Fatalf("Error ejecutando p_Char_Logout: %v", err)
	}
	if retLogout != 1 {
		t.Errorf("p_Char_Logout retornó %d, esperado 1", retLogout)
	}

	// 7. Probar p_Item_Create, p_Item_GetListType y p_Item_Delete
	const testItemKey = int64(9876543210)
	var retItemCreate int
	row = db.QueryRow(`CALL p_Item_Create(?, ?, 1, 0, 1005, 0, @out_ret)`, testItemKey, createdCharNo)
	if err := row.Scan(&retItemCreate); err != nil {
		t.Fatalf("Error ejecutando p_Item_Create: %v", err)
	}
	if retItemCreate != 0 {
		t.Errorf("p_Item_Create retornó %d, esperado 0", retItemCreate)
	}

	itemRows, err := db.Query(`CALL p_Item_GetListType(?, 1)`, createdCharNo)
	if err != nil {
		t.Fatalf("Error ejecutando p_Item_GetListType: %v", err)
	}
	defer itemRows.Close()

	var itemCount int
	for itemRows.Next() {
		var iKey int64
		var storage int16
		var itemID, flags int
		var dDate string
		if err := itemRows.Scan(&iKey, &storage, &itemID, &flags, &dDate); err != nil {
			t.Fatalf("Error escaneando fila de p_Item_GetListType: %v", err)
		}
		if iKey == testItemKey && itemID == 1005 {
			itemCount++
		}
	}
	if itemCount != 1 {
		t.Errorf("p_Item_GetListType no encontró el item creado (%d)", testItemKey)
	}

	var retItemDel int
	row = db.QueryRow(`CALL p_Item_Delete(?, @out_ret)`, testItemKey)
	if err := row.Scan(&retItemDel); err != nil {
		t.Fatalf("Error ejecutando p_Item_Delete: %v", err)
	}

	// 8. Probar p_Char_Delete (Soft delete oficial: bDeleted=1, sID=ID~YYYYMMDD#0)
	var retCharDel int
	row = db.QueryRow(`CALL p_Char_Delete(?, @out_ret)`, createdCharNo)
	if err := row.Scan(&retCharDel); err != nil {
		t.Fatalf("Error ejecutando p_Char_Delete: %v", err)
	}
	if retCharDel != 0 {
		t.Errorf("p_Char_Delete retornó %d, esperado 0", retCharDel)
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

	t.Logf("Flujos funcionales de ciclo de vida de avatares e ítems ejecutados con 100%% de éxito.")
}
