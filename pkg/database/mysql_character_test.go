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
	_, _ = db.Exec("DELETE FROM tCharacter WHERE sID LIKE 'TestKnight%' OR sID = ?", testCharName)

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
		Slot:      1,
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
