package world

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ItemSlot identifiers matching Fiesta client equip slots.
const (
	SlotHelm          byte = 1
	SlotCostumeHelm   byte = 2
	SlotArmor         byte = 7
	SlotCostumeArmor  byte = 8
	SlotWing          byte = 9
	SlotWeapon2       byte = 10
	SlotCostumeShield byte = 11
	SlotWeapon        byte = 12
	SlotCostumeWeapon byte = 13
	SlotTail          byte = 18
	SlotPants         byte = 19
	SlotCostumePants  byte = 20
	SlotBoots         byte = 21
	SlotCostumeBoots  byte = 22
	SlotGlasses       byte = 26
	SlotPet           byte = 28
)

// WriteBasicCharInfo encodes character information into the packet buffer.
// Matches NextGen.World PacketHelper.WriteBasicCharInfo binary layout.
func WriteBasicCharInfo(ch *database.Character, pw *protocol.PacketWriter) {
	pw.WriteInt32(int32(ch.ID))
	pw.WriteString(ch.Name, 16)
	pw.WriteInt32(0) // unk
	pw.WriteInt16(int16(ch.Level))
	pw.WriteByte(ch.Slot)

	mapShortName := "Rou"
	switch ch.Map {
	case 0:
		mapShortName = "Rou"
	case 1:
		mapShortName = "RouCos01"
	case 2:
		mapShortName = "RouVal01"
	case 9:
		mapShortName = "Eld"
	default:
		mapShortName = "Rou"
	}
	pw.WriteString(mapShortName, 13)
	pw.WriteInt32(0) // random seed

	// Look: 4 bytes
	maleBit := byte(0)
	if ch.Male {
		maleBit = 1
	}
	lookByte := byte(0x01 | (ch.Job << 2) | (maleBit << 7))
	pw.WriteByte(lookByte)
	pw.WriteByte(ch.Hair)
	pw.WriteByte(ch.HairColor)
	pw.WriteByte(ch.Face)

	// Equipment: 42 bytes
	pw.WriteUint16(ch.GetEquip(SlotHelm))
	pw.WriteUint16(ch.GetEquip(SlotWeapon))
	pw.WriteUint16(ch.GetEquip(SlotArmor))
	pw.WriteUint16(ch.GetEquip(SlotWeapon2))
	pw.WriteUint16(ch.GetEquip(SlotPants))
	pw.WriteUint16(ch.GetEquip(SlotBoots))
	pw.WriteUint16(ch.GetEquip(SlotCostumeBoots))
	pw.WriteUint16(ch.GetEquip(SlotCostumePants))
	pw.WriteUint16(ch.GetEquip(SlotCostumeArmor))
	pw.Fill(6, 0xFF)
	pw.WriteUint16(ch.GetEquip(SlotGlasses))
	pw.WriteUint16(ch.GetEquip(SlotCostumeHelm))
	pw.Fill(2, 0xFF)
	pw.WriteUint16(ch.GetEquip(SlotCostumeWeapon))
	pw.WriteUint16(ch.GetEquip(SlotWing))
	pw.Fill(2, 0xFF)
	pw.WriteUint16(ch.GetEquip(SlotTail))
	pw.WriteUint16(ch.GetEquip(SlotPet))

	// Refinement: 3 bytes
	pw.Fill(3, 0xFF)

	// Trailer
	pw.WriteByte(0)
	pw.WriteByte(0xF0)
	pw.WriteByte(0xFF)
	pw.WriteByte(0xFF)

	pw.WriteString(mapShortName, 12)
	pw.WriteInt32(0) // pos X
	pw.WriteInt32(0) // pos Y
	pw.WriteUint16(0xDB78)
	pw.WriteUint16(4910)
	pw.WriteUint16(25600)
	pw.Fill(4, 0)
}

// GetBeginnerWeapon returns the starter weapon ID for a given job.
func GetBeginnerWeapon(job protocol.Job) (slot byte, itemID uint16) {
	switch job {
	case protocol.Job_Archer:
		return SlotWeapon2, 1250 // Bow
	case protocol.Job_Fighter:
		return SlotWeapon, 250 // Sword
	case protocol.Job_Cleric:
		return SlotWeapon, 750 // Mace
	case protocol.Job_Mage:
		return SlotWeapon, 1750 // Staff
	case protocol.Job_Trickster:
		return SlotWeapon, 57363 // Claws
	default:
		return SlotWeapon, 250
	}
}
