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

// WriteBasicCharInfo serializa la información del personaje según el struct oficial
// PROTO_AVATARINFORMATION (130 bytes exactos) extraído de WorldManager.pdb / Character.pdb.
func WriteBasicCharInfo(ch *database.Character, pw *protocol.PacketWriter) {
	// 0x000 (4 bytes): unsigned int chrregnum
	pw.WriteUint32(uint32(ch.ID))

	// 0x004 (20 bytes): Name5 name
	pw.WriteString(ch.Name, 20)

	// 0x018 (2 bytes): unsigned short level
	pw.WriteUint16(uint16(ch.Level))

	// 0x01a (1 byte): unsigned char slot
	pw.WriteByte(ch.Slot)

	// 0x01b (12 bytes): Name3 loginmap
	mapShortName := ch.TownName
	if mapShortName == "" {
		mapShortName = "RouN"
	}
	pw.WriteString(mapShortName, 12)

	// 0x027 (5 bytes): PROTO_AVATAR_DELETE_INFO delinfo (year, month, day, hour, min)
	pw.WriteBytes([]byte{0, 0, 0, 0, 0})

	// 0x02c (4 bytes): PROTO_AVATAR_SHAPE_INFO shape (race:2, chrclass:5, gender:1, hairtype, haircolor, faceshape)
	genderBit := byte(0)
	if ch.Male {
		genderBit = 1
	}
	shapeByte0 := byte((1 & 0x03) | ((ch.Job & 0x1F) << 2) | ((genderBit & 0x01) << 7))
	pw.WriteByte(shapeByte0)
	pw.WriteByte(ch.Hair)
	pw.WriteByte(ch.HairColor)
	pw.WriteByte(ch.Face)

	// 0x030 (43 bytes): PROTO_EQUIPMENT equip (20 slots * 2 bytes + 3 bytes upgrade)
	pw.WriteUint16(ch.GetEquip(SlotHelm))          // Equ_Head
	pw.WriteUint16(0)                             // Equ_Mouth
	pw.WriteUint16(ch.GetEquip(SlotWeapon))        // Equ_RightHand
	pw.WriteUint16(ch.GetEquip(SlotArmor))         // Equ_Body
	pw.WriteUint16(ch.GetEquip(SlotWeapon2))       // Equ_LeftHand
	pw.WriteUint16(ch.GetEquip(SlotPants))         // Equ_Pant
	pw.WriteUint16(ch.GetEquip(SlotBoots))         // Equ_Boot
	pw.WriteUint16(ch.GetEquip(SlotCostumeBoots))  // Equ_AccBoot
	pw.WriteUint16(ch.GetEquip(SlotCostumePants))  // Equ_AccPant
	pw.WriteUint16(ch.GetEquip(SlotCostumeArmor))  // Equ_AccBody
	pw.WriteUint16(ch.GetEquip(SlotCostumeHelm))   // Equ_AccHeadA
	pw.WriteUint16(0)                             // Equ_MiniMon_R
	pw.WriteUint16(ch.GetEquip(SlotGlasses))       // Equ_Eye
	pw.WriteUint16(ch.GetEquip(SlotCostumeShield)) // Equ_AccLeftHand
	pw.WriteUint16(ch.GetEquip(SlotCostumeWeapon)) // Equ_AccRightHand
	pw.WriteUint16(ch.GetEquip(SlotWing))          // Equ_AccBack
	pw.WriteUint16(0)                             // Equ_CosEff
	pw.WriteUint16(ch.GetEquip(SlotTail))          // Equ_AccHip
	pw.WriteUint16(ch.GetEquip(SlotPet))           // Equ_Minimon
	pw.WriteUint16(0)                             // Equ_AccShield
	pw.WriteBytes([]byte{0, 0, 0})                // upgrade (3 bytes: lefthand, righthand, body, leg, shoes)

	// 0x05b (4 bytes): unsigned int nKQHandle
	pw.WriteUint32(0)

	// 0x05f (12 bytes): Name3 sKQMapName
	pw.WriteString("", 12)

	// 0x06b (8 bytes): SHINE_XY_TYPE nKQCoord (x, y)
	pw.WriteUint32(0)
	pw.WriteUint32(0)

	// 0x073 (4 bytes): SHINE_DATETIME dKQDate
	pw.WriteUint32(0)

	// 0x077 (6 bytes): CHAR_ID_CHANGE_DATA CharIDChangeData (bNeedChangeID, bInit, nRowNo)
	pw.WriteByte(0)
	pw.WriteByte(0)
	pw.WriteUint32(0)

	// 0x07d (5 bytes): PROTO_TUTORIAL_INFO TutorialInfo (nTutorialState, nTutorialStep)
	pw.WriteUint32(0)
	pw.WriteByte(0)
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
	case protocol.Job_Crusader:
		return SlotWeapon, 62923 // Novice Crusader's Blade
	default:
		return SlotWeapon, 250
	}
}
