package database_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/database"
)

func setupTestCharDB(t *testing.T) (*sql.DB, *database.MySQLCharacterRepository) {
	cfg, err := config.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping MySQL: %v", err)
	}

	repo := database.NewMySQLCharacterRepository(db)
	return db, repo
}

func TestMySQLCharacterRepository_Lifecycle(t *testing.T) {
	db, repo := setupTestCharDB(t)
	defer db.Close()

	testUserNo := 1 // Admin user
	testCharName := "TestKnight"

	// Limpieza previa
	_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestKnight%' OR sID = ? OR (nUserNo = ? AND nSlotNo IN (4, 5))", testCharName, testUserNo)

	// 1. Crear personaje
	ch := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      4,
		Job:       1, // Fighter
		Level:     1,
		Male:      true,
		Hair:      2,
		HairColor: 3,
		Face:      1,
		TownName:  "RouN",
		CurHP:     120,
		CurSP:     80,
		Equips:    map[byte]uint16{0: 250}, // Wooden Sword
	}

	err := repo.CreateCharacter(ch)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	if ch.ID == 0 {
		t.Fatalf("Expected non-zero character ID, got %d", ch.ID)
	}

	// Validar que se creó tCharacterOptions
	var optionsCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tCharacterOptions WHERE nCharNo = ?", ch.ID).Scan(&optionsCount)
	if err != nil || optionsCount != 1 {
		t.Fatalf("Expected tCharacterOptions row for char %d, count=%d, err=%v", ch.ID, optionsCount, err)
	}

	// 2. Error por slot en uso
	chDupSlot := &database.Character{
		AccountID: testUserNo,
		Name:      "AnotherKnight",
		Slot:      ch.Slot, // Mismo slot
		Job:       2,
	}
	err = repo.CreateCharacter(chDupSlot)
	if !errors.Is(err, database.ErrSlotInUse) {
		t.Fatalf("Expected ErrSlotInUse, got %v", err)
	}

	// 3. Error por nombre duplicado
	chDupName := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      5,
		Job:       2,
	}
	err = repo.CreateCharacter(chDupName)
	if !errors.Is(err, database.ErrCharacterExists) {
		t.Fatalf("Expected ErrCharacterExists, got %v", err)
	}

	// 4. Buscar por cuenta
	list, err := repo.GetCharactersByAccount(testUserNo)
	if err != nil {
		t.Fatalf("Failed to get characters by account: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("Expected at least 1 character for account %d", testUserNo)
	}

	found := false
	for _, c := range list {
		if c.ID == ch.ID && c.Name == testCharName {
			found = true
			if c.Equips[0] != 250 {
				t.Errorf("Expected weapon 250 in slot 0, got %d", c.Equips[0])
			}
			break
		}
	}
	if !found {
		t.Errorf("Created character was not found in account list")
	}

	// 5. Buscar por nombre
	byName, err := repo.GetCharacterByName(testCharName)
	if err != nil {
		t.Fatalf("Failed to get character by name: %v", err)
	}
	if byName.ID != ch.ID {
		t.Errorf("Expected ID %d, got %d", ch.ID, byName.ID)
	}

	// 6. Buscar por ID
	byID, err := repo.GetCharacterByID(ch.ID)
	if err != nil {
		t.Fatalf("Failed to get character by ID: %v", err)
	}
	if byID.Name != testCharName {
		t.Errorf("Expected Name %s, got %s", testCharName, byID.Name)
	}

	// 7. Borrado suave oficial (soft delete)
	err = repo.DeleteCharacter(testUserNo, ch.Slot)
	if err != nil {
		t.Fatalf("Failed to soft delete character: %v", err)
	}

	// Verificar que ya no aparece en GetCharactersByAccount
	listAfterDelete, err := repo.GetCharactersByAccount(testUserNo)
	if err != nil {
		t.Fatalf("Failed to get characters after delete: %v", err)
	}
	for _, c := range listAfterDelete {
		if c.ID == ch.ID {
			t.Errorf("Character should not appear in active account list after delete")
		}
	}

	// Verificar que el nombre quedó libre (se renombró a nombre~YYYYMMDD#0)
	var deletedSID string
	var bDeleted int
	err = db.QueryRow("SELECT sID, bDeleted FROM tCharacter WHERE nCharNo = ?", ch.ID).Scan(&deletedSID, &bDeleted)
	if err != nil {
		t.Fatalf("Failed to query deleted char row: %v", err)
	}
	if bDeleted != 1 {
		t.Errorf("Expected bDeleted=1, got %d", bDeleted)
	}
	if deletedSID == testCharName {
		t.Errorf("Expected sID to be renamed with tilde suffix, got %s", deletedSID)
	}

	// Ahora el nombre original puede volverse a crear
	chRecreate := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      ch.Slot,
		Job:       1,
		Level:     1,
		Male:      false,
		TownName:  "RouN",
	}
	err = repo.CreateCharacter(chRecreate)
	if err != nil {
		t.Fatalf("Expected to be able to reuse character name after soft delete, got: %v", err)
	}

	// Limpieza final
	_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestKnight%' OR sID = ?", testCharName)
}

func TestMySQLCharacterRepository_StarterItemsAndShortcuts(t *testing.T) {
	db, repo := setupTestCharDB(t)
	defer db.Close()

	testUserNo := 1
	testCharName := "TestStarterChar"

	// Limpieza previa
	_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestStarterChar%' OR sID = ? OR (nUserNo = ? AND nSlotNo = 3)", testCharName, testUserNo)

	ch := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      3,
		Job:       21, // Joker
		Level:     1,
		Male:      true,
		TownName:  "RouN",
		XPos:      6888,
		YPos:      7045,
		Equips:    map[byte]uint16{12: 57363},
		Items: []*database.Item{
			{StorageType: 9, Storage: 0, ItemID: 46003, Flags: 1},
		},
		Shortcuts: []*database.Shortcut{
			{SlotNo: 0, CodeType: 1, CodeNo: 4, Value: 0},
			{SlotNo: 10, CodeType: 1, CodeNo: 1, Value: 3500},
		},
		Skills: []*database.CharacterSkill{
			{SkillNo: 29200, SkillLevel: 1},
		},
	}

	err := repo.CreateCharacter(ch)
	if err != nil {
		t.Fatalf("Failed to create character with starters: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestStarterChar%' OR sID = ?", testCharName)
		_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", ch.ID)
		_, _ = db.Exec("DELETE FROM tCharacterShortCut WHERE nCharNo = ?", ch.ID)
		_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", ch.ID)
	}()

	// 1. Verificar lectura de items
	items, err := repo.GetCharacterItems(ch.ID)
	if err != nil {
		t.Fatalf("Failed to get character items: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("Expected 2 items (1 equip + 1 inven), got %d", len(items))
	}

	// 2. Verificar lectura de shortcuts
	shortcuts, err := repo.GetCharacterShortcuts(ch.ID)
	if err != nil {
		t.Fatalf("Failed to get character shortcuts: %v", err)
	}
	if len(shortcuts) != 2 {
		t.Fatalf("Expected 2 shortcuts, got %d", len(shortcuts))
	}

	// 3. Verificar lectura de skills
	skills, err := repo.GetCharacterSkills(ch.ID)
	if err != nil {
		t.Fatalf("Failed to get character skills: %v", err)
	}
	if len(skills) != 1 || skills[0].SkillNo != 29200 {
		t.Fatalf("Expected 1 skill with ID 29200, got %+v", skills)
	}

	// 4. Verificar tCharacterOptions (p_Char_SetOptShortCutData / p_Char_GetOptShortCutData)
	testSCData := make([]byte, 1024)
	testSCData[0] = 0xAA
	testSCData[1023] = 0xBB
	if err := repo.SaveCharacterOption(ch.ID, "shortcut_data", testSCData); err != nil {
		t.Fatalf("Failed to save shortcut_data option: %v", err)
	}

	testGameOpt := []byte{1, 2, 3, 4, 5}
	if err := repo.SaveCharacterOption(ch.ID, "game", testGameOpt); err != nil {
		t.Fatalf("Failed to save game option: %v", err)
	}

	opts, err := repo.GetCharacterOptions(ch.ID)
	if err != nil {
		t.Fatalf("Failed to get character options: %v", err)
	}
	if len(opts.ShortCutData) == 0 || opts.ShortCutData[0] != 0xAA || opts.ShortCutData[1023] != 0xBB {
		t.Fatalf("Character options shortcut_data mismatch: got len=%d", len(opts.ShortCutData))
	}
	if len(opts.Game) == 0 || opts.Game[0] != 1 {
		t.Fatalf("Character options game mismatch: got %+v", opts.Game)
	}
}

func TestMySQLCharacterRepository_ZonePersistence(t *testing.T) {
	db, repo := setupTestCharDB(t)
	defer db.Close()

	testUserNo := 1
	testCharName := "ZonePersistenceKnight"

	// Limpieza previa
	_, _ = db.Exec("DELETE FROM tCharacter WHERE sID = ?", testCharName)

	ch := &database.Character{
		AccountID: testUserNo,
		Name:      testCharName,
		Slot:      5,
		Job:       1,
		Level:     1,
		Male:      true,
		TownName:  "RouN",
		CurHP:     100,
		CurSP:     50,
	}

	if err := repo.CreateCharacter(ch); err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nCharNo = ?", ch.ID)
	}()

	// 1. Test CreateItem & GetItemsByStorage
	item1 := &database.Item{
		StorageType: 9, // Inventory
		Owner:       ch.ID,
		Storage:     0,
		ItemID:      1001,
		Flags:       5,
	}
	if err := repo.CreateItem(item1); err != nil {
		t.Fatalf("Failed to create item1: %v", err)
	}
	if item1.ItemKey == 0 {
		t.Fatalf("Expected non-zero ItemKey")
	}

	item2 := &database.Item{
		StorageType: 9,
		Owner:       ch.ID,
		Storage:     1,
		ItemID:      1002,
		Flags:       1,
	}
	if err := repo.CreateItem(item2); err != nil {
		t.Fatalf("Failed to create item2: %v", err)
	}

	invenItems, err := repo.GetItemsByStorage(ch.ID, 9)
	if err != nil {
		t.Fatalf("Failed to get items by storage: %v", err)
	}
	if len(invenItems) != 2 {
		t.Fatalf("Expected 2 inventory items, got %d", len(invenItems))
	}

	// 2. Test SetItemStorage (e.g. equip item to box 8 slot 1)
	if err := repo.SetItemStorage(item1.ItemKey, 8, 1); err != nil {
		t.Fatalf("Failed to set item storage: %v", err)
	}
	equipItems, err := repo.GetItemsByStorage(ch.ID, 8)
	if err != nil || len(equipItems) != 1 || equipItems[0].Storage != 1 {
		t.Fatalf("Expected 1 equip item in slot 1, got %+v (err=%v)", equipItems, err)
	}

	// 3. Test ExchangeItemStorage
	// Move item1 back to box 9 slot 0, then exchange item1 (slot 0) and item2 (slot 1)
	_ = repo.SetItemStorage(item1.ItemKey, 9, 0)
	if err := repo.ExchangeItemStorage(item1.ItemKey, 9, 1, item2.ItemKey, 9, 0); err != nil {
		t.Fatalf("Failed to exchange items: %v", err)
	}
	updatedInven, _ := repo.GetItemsByStorage(ch.ID, 9)
	for _, it := range updatedInven {
		if it.ItemKey == item1.ItemKey && it.Storage != 1 {
			t.Fatalf("Expected item1 in slot 1, got %d", it.Storage)
		}
		if it.ItemKey == item2.ItemKey && it.Storage != 0 {
			t.Fatalf("Expected item2 in slot 0, got %d", it.Storage)
		}
	}

	// 4. Test DeleteItem
	if err := repo.DeleteItem(item1.ItemKey); err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}
	afterDelete, _ := repo.GetItemsByStorage(ch.ID, 9)
	if len(afterDelete) != 1 || afterDelete[0].ItemKey != item2.ItemKey {
		t.Fatalf("Expected only item2 remaining, got %d items", len(afterDelete))
	}

	// 5. Test AddCharMoney & SubCharMoney
	if err := repo.AddCharMoney(ch.ID, 5000); err != nil {
		t.Fatalf("Failed to add money: %v", err)
	}
	// Verify current money
	chFetched, err := repo.GetCharacterByID(ch.ID)
	if err != nil || chFetched.TCharacterState.Money != 5000 {
		t.Fatalf("Expected 5000 money, got %d (err=%v)", chFetched.TCharacterState.Money, err)
	}

	// Overdraft deduction should be rejected
	ok, err := repo.SubCharMoney(ch.ID, 10000)
	if err != nil || ok {
		t.Fatalf("Expected overdraft deduction to return ok=false, got ok=%v, err=%v", ok, err)
	}

	// Valid deduction
	ok, err = repo.SubCharMoney(ch.ID, 2000)
	if err != nil || !ok {
		t.Fatalf("Expected successful deduction, got ok=%v, err=%v", ok, err)
	}
	chFetched2, _ := repo.GetCharacterByID(ch.ID)
	if chFetched2.TCharacterState.Money != 3000 {
		t.Fatalf("Expected 3000 money remaining, got %d", chFetched2.TCharacterState.Money)
	}

	// 6. Test SaveCharLocation
	if err := repo.SaveCharLocation(ch.ID, "Eld", 1500, 2500, 0x88); err != nil {
		t.Fatalf("Failed to save location: %v", err)
	}
	chLoc, _ := repo.GetCharacterByID(ch.ID)
	if chLoc.TownName != "Eld" || chLoc.XPos != 1500 || chLoc.YPos != 2500 {
		t.Fatalf("Expected location Eld (1500, 2500), got %s (%d, %d)", chLoc.TownName, chLoc.XPos, chLoc.YPos)
	}

	// 7. Test SaveCharStats
	if err := repo.SaveCharStats(ch.ID, 250, 180, 15, 7); err != nil {
		t.Fatalf("Failed to save stats: %v", err)
	}
	chStats, _ := repo.GetCharacterByID(ch.ID)
	if chStats.CurHP != 250 || chStats.CurSP != 180 {
		t.Fatalf("Expected HP=250 SP=180, got HP=%d SP=%d", chStats.CurHP, chStats.CurSP)
	}

	// 8. Test SaveCharLevelExp
	if err := repo.SaveCharLevelExp(ch.ID, 10, 50000); err != nil {
		t.Fatalf("Failed to save level exp: %v", err)
	}
	chLvl, _ := repo.GetCharacterByID(ch.ID)
	if chLvl.Level != 10 || chLvl.TCharacterState.Exp != 50000 {
		t.Fatalf("Expected Level 10, Exp 50000, got Lvl=%d Exp=%d", chLvl.Level, chLvl.TCharacterState.Exp)
	}

	// 9. Test SaveCharFreeStats
	if err := repo.SaveCharFreeStats(ch.ID, 15, 12, 10, 5, 8, 3); err != nil {
		t.Fatalf("Failed to save free stats: %v", err)
	}
}
