package database_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
)

func openTestMySQL(t *testing.T) *sql.DB {
	t.Helper()
	cfg, err := config.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open MySQL: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("Failed to ping MySQL at %s: %v", dsn, err)
	}
	return db
}

func TestMySQLConnection(t *testing.T) {
	db := openTestMySQL(t)
	defer db.Close()

	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM tUser").Scan(&userCount)
	if err != nil {
		t.Fatalf("Failed to query tUser: %v", err)
	}

	if userCount < 1 {
		t.Fatalf("Expected at least 1 seed user in tUser, got %d", userCount)
	}
	t.Logf("Successfully connected to MySQL! Found %d users in tUser.", userCount)
}

func cleanupTestUserAndChar(db *sql.DB, userNo int, charName string) {
	_, _ = db.Exec("DELETE FROM tKingdomQuest WHERE nQuestNo = 31999")
	_, _ = db.Exec("DELETE FROM tItemOptions WHERE nItemKey IN (SELECT k FROM (SELECT i.nItemKey AS k FROM tItem i JOIN tCharacter c ON i.nOwner = c.nCharNo WHERE c.nUserNo = ? OR c.sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tItemTitle WHERE nItemKey IN (SELECT k FROM (SELECT i.nItemKey AS k FROM tItem i JOIN tCharacter c ON i.nOwner = c.nCharNo WHERE c.nUserNo = ? OR c.sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tItem WHERE nOwner IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterOptions WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterShape WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterShortCut WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterKeyMap WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterGameOption WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterSingleOption WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterState WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterTitle WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tMinihouseGuest WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tMinihouse WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tCharacterLoading WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tSpamBlockTime WHERE nCharNo IN (SELECT id FROM (SELECT nCharNo AS id FROM tCharacter WHERE nUserNo = ? OR sID = ?) AS tmp)", userNo, charName)
	_, _ = db.Exec("DELETE FROM tSpamFilter WHERE nUserNo = ?", userNo)
	_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ? OR sID = ?", userNo, charName)
	_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", userNo)
}

// TestMySQLSPIdempotenceAndVarbinaryWidths verifica el Punto 1 del Plan:
// Invocar dos veces consecutivas con exactamente los mismos argumentos cada método
// de persistencia/upsert del repositorio sobre MySQL 8.4 real nunca debe devolver
// Error 1062 (Duplicate entry), Error 1406 (Data too long) ni fallar por ROW_COUNT() = 0.
func TestMySQLSPIdempotenceAndVarbinaryWidths(t *testing.T) {
	db := openTestMySQL(t)
	defer db.Close()

	charRepo := database.NewMySQLCharacterRepository(db)

	const testUserNo = 99401
	const testCharName = "IdemHero"
	cleanupTestUserAndChar(db, testUserNo, testCharName)
	defer cleanupTestUserAndChar(db, testUserNo, testCharName)
	_, err := db.Exec(
		"INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP) VALUES (?, 'idem_user_99401', 'pass', 'IdemTest', 0, 0, 1, '127.0.0.1')",
		testUserNo,
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	ch := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      0,
		Race:      1,
		Job:       1,
		Level:     15,
		Male:      true,
		Hair:      2,
		HairColor: 3,
		Face:      1,
		TownName:  "RouN",
		XPos:      5000,
		YPos:      6000,
		CurHP:     250,
		CurSP:     120,
	}
	if err := charRepo.CreateCharacter(ch); err != nil {
		t.Fatalf("CreateCharacter failed: %v", err)
	}

	now := time.Now().Truncate(time.Second)

	// 1. UpdateCharacter (p_Char_SetAllData + p_Char_SetShape) x2 idéntico
	for i := 0; i < 2; i++ {
		if err := charRepo.UpdateCharacter(ch); err != nil {
			t.Fatalf("UpdateCharacter pass %d failed: %v", i+1, err)
		}
	}

	// 2. SetCharacterShape (p_Char_SetShape) x2 idéntico
	for i := 0; i < 2; i++ {
		if err := charRepo.SetCharacterShape(ch.ID, ch.Race, ch.Job, 1, ch.Hair, ch.HairColor, ch.Face); err != nil {
			t.Fatalf("SetCharacterShape pass %d failed: %v", i+1, err)
		}
	}

	// 3. SaveCharLevelExp (p_Char_SetLevelExp), SaveCharMoney (p_Char_MoneySet), SaveCharFreeStats (usp_Character_setStatPoint) x2 idéntico
	for i := 0; i < 2; i++ {
		if err := charRepo.SaveCharLevelExp(ch.ID, int16(ch.Level), 12345); err != nil {
			t.Fatalf("SaveCharLevelExp pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SaveCharMoney(ch.ID, 50000); err != nil {
			t.Fatalf("SaveCharMoney pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SaveCharFreeStats(ch.ID, 10, 10, 10, 10, 10, 5); err != nil {
			t.Fatalf("SaveCharFreeStats pass %d failed: %v", i+1, err)
		}
	}

	// 4. SetCharacterTitleProgress (p_Char_SetTitle) y SetCharacterState (p_Char_SetStateWithRestTime) x2 idéntico
	for i := 0; i < 2; i++ {
		if err := charRepo.SetCharacterTitleProgress(ch.ID, 1, 0, 5); err != nil {
			t.Fatalf("SetCharacterTitleProgress pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetCharacterState(ch.ID, 100, 1, 60000, 30000, 1, 0, 0, 0, 0); err != nil {
			t.Fatalf("SetCharacterState pass %d failed: %v", i+1, err)
		}
	}

	// 5. SaveCharacterSkill (p_Skill_Set), SetSkillCoolTime (p_Skill_SetCoolTime), SetSkillPower (p_Skill_SetPowerAll), SetSingleSkillPower (p_Skill_SetPower) x2 idéntico
	sk := &database.CharacterSkill{
		CharNo:     ch.ID,
		SkillNo:    101,
		SkillLevel: 1,
		SkillExp:   50,
	}
	for i := 0; i < 2; i++ {
		if err := charRepo.SaveCharacterSkill(ch.ID, sk); err != nil {
			t.Fatalf("SaveCharacterSkill pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetSkillCoolTime(ch.ID, 101, 1000, 500); err != nil {
			t.Fatalf("SetSkillCoolTime pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetSkillPower(ch.ID, 101, 1, 1, 1, 1); err != nil {
			t.Fatalf("SetSkillPower pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetSingleSkillPower(ch.ID, 101, 0, 2); err != nil {
			t.Fatalf("SetSingleSkillPower pass %d failed: %v", i+1, err)
		}
	}

	// 6. SetCharacterQuest (p_Quest_Set) con payload binario de 100 bytes x2 idéntico
	questPayload := bytes.Repeat([]byte{0xAB}, 100)
	for i := 0; i < 2; i++ {
		if err := charRepo.SetCharacterQuest(ch.ID, 501, 1, questPayload); err != nil {
			t.Fatalf("SetCharacterQuest pass %d failed: %v", i+1, err)
		}
	}

	// 7. SaveCharacterOption / GetCharacterOptions con anchos completos NA2016 (1024B ShortCutData, 392B WindowsPos, 308B KeyMapping, 64B Game, 60B Video, 24B ShortCutSize, 1B Sound) x2 idéntico
	optPayloads := map[string][]byte{
		"ShortCutData": bytes.Repeat([]byte{0x11}, 1024),
		"ShortCutSize": bytes.Repeat([]byte{0x22}, 24),
		"Video":        bytes.Repeat([]byte{0x33}, 60),
		"Sound":        {0x01},
		"Game":         bytes.Repeat([]byte{0x44}, 64),
		"WindowsPos":   bytes.Repeat([]byte{0x55}, 392),
		"KeyMapping":   bytes.Repeat([]byte{0x66}, 308),
	}
	for i := 0; i < 2; i++ {
		for optType, payload := range optPayloads {
			if err := charRepo.SaveCharacterOption(ch.ID, optType, payload); err != nil {
				t.Fatalf("SaveCharacterOption(%s) pass %d failed: %v", optType, i+1, err)
			}
		}
	}
	loadedOpts, err := charRepo.GetCharacterOptions(ch.ID)
	if err != nil {
		t.Fatalf("GetCharacterOptions failed: %v", err)
	}
	if len(loadedOpts.WindowsPos) != 392 || len(loadedOpts.KeyMapping) != 308 || len(loadedOpts.ShortCutData) != 1024 {
		t.Fatalf("GetCharacterOptions returned unexpected lengths: win=%d key=%d sc=%d",
			len(loadedOpts.WindowsPos), len(loadedOpts.KeyMapping), len(loadedOpts.ShortCutData))
	}

	// 8. Improve Options (usp_Character_setKeyMap, usp_Character_setGameOption), SetItemActionCoolTime, SetAutoPick, SetRestExpExecTime, SetMiniHouseURL, SetSpamBlockTime, SetSpamFilter, SetKingdomQuest x2 idéntico
	for i := 0; i < 2; i++ {
		if err := charRepo.SetImproveKeyMap(ch.ID, 1, 0, 65); err != nil {
			t.Fatalf("SetImproveKeyMap pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetImproveGameOption(ch.ID, 1, 1); err != nil {
			t.Fatalf("SetImproveGameOption pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetItemActionCoolTime(ch.ID, 1, now.Add(time.Minute)); err != nil {
			t.Fatalf("SetItemActionCoolTime pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetAutoPick(ch.ID, true); err != nil {
			t.Fatalf("SetAutoPick pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetRestExpExecTime(ch.ID, now); err != nil {
			t.Fatalf("SetRestExpExecTime pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetMiniHouseURL(ch.ID, "https://xelion.online/home/idem"); err != nil {
			t.Fatalf("SetMiniHouseURL pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetSpamBlockTime(ch.ID, 60, 30); err != nil {
			t.Fatalf("SetSpamBlockTime pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetSpamFilter(testUserNo, 1, 1, 50); err != nil {
			t.Fatalf("SetSpamFilter pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetKingdomQuest(31999, bytes.NewBuffer(bytes.Repeat([]byte("K"), 256)).String()); err != nil {
			t.Fatalf("SetKingdomQuest(256B) pass %d failed: %v", i+1, err)
		}
	}

	// 9. Item persistence & options idempotence (p_Item_SetData, p_Item_SetStorage, p_Item_SetOption) x2 idéntico
	itemKey := time.Now().UnixNano()
	it := &database.Item{
		ItemKey:     itemKey,
		Owner:       ch.ID,
		StorageType: 9,
		Storage:     0,
		ItemID:      1001,
		Flags:       1,
		Date:        now,
	}
	if err := charRepo.CreateItem(it); err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := charRepo.SetItemData(itemKey, ch.ID, 9, 0, 1001, 1, now); err != nil {
			t.Fatalf("SetItemData pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetItemStorage(itemKey, 9, 0); err != nil {
			t.Fatalf("SetItemStorage pass %d failed: %v", i+1, err)
		}
		if err := charRepo.SetItemOption(itemKey, 0, 5); err != nil {
			t.Fatalf("SetItemOption pass %d failed: %v", i+1, err)
		}
	}
}

// TestCharacterSyncDoesNotOverwriteMutatedFlatFields verifica el Punto 2 del Plan:
// Cargar un Character desde MySQL (TCharacter.CharNo != 0), modificar únicamente sus
// campos planos (Char.XPos, Char.YPos, Char.Hair, Char.HairColor, Char.Face, Char.Level,
// Char.CurHP, Char.CurSP), llamar a UpdateCharacter(ch) y volver a leerlo desde MySQL con
// GetCharacterByID(ch.ID), comprobando que todos los cambios se guardaron sin ser pisados
// por los valores viejos de TCharacter / TCharacterShape / TCharacterState.
func TestCharacterSyncDoesNotOverwriteMutatedFlatFields(t *testing.T) {
	db := openTestMySQL(t)
	defer db.Close()

	charRepo := database.NewMySQLCharacterRepository(db)

	const testUserNo = 99402
	const testCharName = "SyncHero"
	cleanupTestUserAndChar(db, testUserNo, testCharName)
	defer cleanupTestUserAndChar(db, testUserNo, testCharName)
	_, err := db.Exec(
		"INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP) VALUES (?, 'sync_user_99402', 'pass', 'SyncTest', 0, 0, 1, '127.0.0.1')",
		testUserNo,
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	initial := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      1,
		Race:      1,
		Job:       6,
		Level:     10,
		Male:      false,
		Hair:      1,
		HairColor: 2,
		Face:      3,
		TownName:  "RouN",
		XPos:      1000,
		YPos:      2000,
		CurHP:     150,
		CurSP:     80,
	}
	if err := charRepo.CreateCharacter(initial); err != nil {
		t.Fatalf("CreateCharacter failed: %v", err)
	}

	loaded, err := charRepo.GetCharacterByID(initial.ID)
	if err != nil {
		t.Fatalf("GetCharacterByID failed: %v", err)
	}
	if loaded.TCharacter.CharNo == 0 {
		t.Fatalf("Expected loaded.TCharacter.CharNo != 0, got 0")
	}

	// Modificar ÚNICAMENTE los campos planos (simulando un handler de ZoneService / WorldService)
	loaded.XPos = 7777
	loaded.YPos = 8888
	loaded.Hair = 9
	loaded.HairColor = 11
	loaded.Face = 5
	loaded.Level = 42
	loaded.CurHP = 999
	loaded.CurSP = 444
	loaded.TownName = "Eld"

	if err := charRepo.UpdateCharacter(loaded); err != nil {
		t.Fatalf("UpdateCharacter failed: %v", err)
	}

	reloaded, err := charRepo.GetCharacterByID(initial.ID)
	if err != nil {
		t.Fatalf("GetCharacterByID after UpdateCharacter failed: %v", err)
	}

	if reloaded.XPos != 7777 || reloaded.YPos != 8888 {
		t.Fatalf("Position was overwritten by stale TCharacterState! got (%d, %d), want (7777, 8888)", reloaded.XPos, reloaded.YPos)
	}
	if reloaded.Hair != 9 || reloaded.HairColor != 11 || reloaded.Face != 5 {
		t.Fatalf("Appearance was overwritten by stale TCharacterShape! got (hair=%d, color=%d, face=%d), want (9, 11, 5)",
			reloaded.Hair, reloaded.HairColor, reloaded.Face)
	}
	if reloaded.Level != 42 {
		t.Fatalf("Level was overwritten by stale TCharacter! got %d, want 42", reloaded.Level)
	}
	if reloaded.CurHP != 999 || reloaded.CurSP != 444 {
		t.Fatalf("Vitals were overwritten by stale TCharacterState! got (HP=%d, SP=%d), want (999, 444)", reloaded.CurHP, reloaded.CurSP)
	}
	if reloaded.TownName != "Eld" {
		t.Fatalf("TownName was not persisted! got %q, want %q", reloaded.TownName, "Eld")
	}
}
