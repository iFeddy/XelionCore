package world

import (
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestWriteBasicCharInfo(t *testing.T) {
	ch := &database.Character{
		ID:        1,
		AccountID: 10,
		Name:      "PlayerOne",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
		Hair:      1,
		HairColor: 2,
		Face:      3,
		TownName:  "Roumen",
		Equips:    map[byte]uint16{SlotWeapon: 250},
	}

	pw := protocol.NewPacketWriter(protocol.SH5Type_CharCreationOK)
	pw.WriteByte(1)
	WriteBasicCharInfo(ch, pw)

	wire := pw.BuildWirePacket()
	if len(wire) < 50 {
		t.Fatalf("expected packet to be at least 50 bytes, got %d", len(wire))
	}
}

func TestCharacterRepositoryOperations(t *testing.T) {
	repo := database.NewMemoryCharacterRepository()

	ch1 := &database.Character{
		AccountID: 1,
		Name:      "Warrior",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
	}

	err := repo.CreateCharacter(ch1)
	if err != nil {
		t.Fatalf("unexpected error creating char: %v", err)
	}

	// Duplicate name error
	chDup := &database.Character{
		AccountID: 2,
		Name:      "Warrior",
		Slot:      0,
		Job:       byte(protocol.Job_Archer),
	}
	if err := repo.CreateCharacter(chDup); err != database.ErrCharacterExists {
		t.Fatalf("expected ErrCharacterExists, got: %v", err)
	}

	// Slot collision error
	chSlot := &database.Character{
		AccountID: 1,
		Name:      "MageHero",
		Slot:      0,
		Job:       byte(protocol.Job_Mage),
	}
	if err := repo.CreateCharacter(chSlot); err != database.ErrSlotInUse {
		t.Fatalf("expected ErrSlotInUse, got: %v", err)
	}

	// Retrieve
	chars, err := repo.GetCharactersByAccount(1)
	if err != nil || len(chars) != 1 {
		t.Fatalf("expected 1 character, got %d (err: %v)", len(chars), err)
	}

	// Delete
	err = repo.DeleteCharacter(1, 0)
	if err != nil {
		t.Fatalf("unexpected error deleting character: %v", err)
	}

	charsAfter, _ := repo.GetCharactersByAccount(1)
	if len(charsAfter) != 0 {
		t.Fatalf("expected 0 characters after deletion, got %d", len(charsAfter))
	}
}
